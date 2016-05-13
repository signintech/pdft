package pdft

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//PDFData pdf file data
type PDFData struct {
	trailer TrailerData
	xrefs   []XrefData
	objIDs  []int
	objs    []PDFObjData
}

//TrailerData trailer
type TrailerData struct {
	rootObjID int
}

func (p *PDFData) put(pdfobj PDFObjData) {
	p.objIDs = append(p.objIDs, pdfobj.objID)
	p.objs = append(p.objs, pdfobj)
}

//GetObjByID get obj by objid
func (p *PDFData) getObjByID(objID int) *PDFObjData {
	for i, id := range p.objIDs {
		if id == objID {
			return &p.objs[i]
		}
	}
	return nil
}

func (p *PDFData) maxID() int {
	max := 0
	for _, id := range p.objIDs {
		if id > max {
			max = id
		}
	}
	return max
}

func (p *PDFData) injectFontsToPDF(fontDatas map[string]*PDFFontData) error {

	var err error
	var cw crawl
	cw.set(p, p.trailer.rootObjID, "Pages", "Kids", "Resources", "Font")
	err = cw.run()
	if err != nil {
		return err
	}

	maxFontIndex, err := findMaxFontIndex(&cw, p)
	if err != nil {
		return err
	}

	var newCrFonts crawlResultFonts //font ใหม่ที่จะยัดเข้าไป
	for _, pdffontdata := range fontDatas {
		maxFontIndex++
		newCrFonts.append(maxFontIndex, pdffontdata.fontID)
		pdffontdata.setFontIndex(maxFontIndex)
	}

	objMustReplaces := make(map[int]string)
	//หา obj ที่ต้องยัด font ใหม่ลงไป
	for objID, r := range cw.results { //วน แต่ละ ojb
		fontPropVal, err := r.valOf("Font")
		if err == ErrCrawlResultValOfNotFound {
			continue
		} else if err != nil {
			return err
		}

		fontPropValType := propertyType(fontPropVal)
		if fontPropValType == object {
			var crFonts crawlResultFonts
			tmp := []byte(fontPropVal)
			err = crFonts.parse(&tmp)
			if err != nil {
				return err
			}
			crFonts = append(crFonts, newCrFonts...)
			r.setValOf("Font", "<<\n"+crFonts.String()+">>\n")
			objMustReplaces[objID] = r.String()
		} else if fontPropValType == dictionary {
			var fontObjID int
			fontObjID, _, err = readObjIDFromDictionary(fontPropVal)
			if err != nil {
				return err
			}
			var crFonts crawlResultFonts
			fontObj := p.getObjByID(fontObjID)
			err = crFonts.parse(&fontObj.data)
			if err != nil {
				return err
			}
			crFonts = append(crFonts, newCrFonts...)
			objMustReplaces[fontObjID] = crFonts.String()
		}
	}

	for objID := range objMustReplaces {

		p.getObjByID(objID).data = []byte("<<\n" + objMustReplaces[objID] + ">>\n")
	}

	return nil
}

func (p *PDFData) injectContentToPDF(contenters *[]Contenter) error {
	var err error
	var cw crawl
	cw.set(p, p.trailer.rootObjID, "Pages", "Kids")
	err = cw.run()
	if err != nil {
		return err
	}

	pageBuffs := make(map[int]*bytes.Buffer)
	for _, ctn := range *contenters {
		pageNum := ctn.page()
		if _, ok := pageBuffs[pageNum]; !ok {
			pageBuffs[pageNum] = new(bytes.Buffer)
		}
		var buff *bytes.Buffer
		buff, err = ctn.toSteram()
		if err != nil {
			return err
		}
		_, err = buff.WriteTo(pageBuffs[pageNum])
		if err != nil {
			return err
		}
	}

	var pageObjIDs []int
	for _, r := range cw.results { //วน แต่ละ obj
		var propKidsVal string
		propKidsVal, err = r.valOf("Kids")
		if err == ErrCrawlResultValOfNotFound {
			continue
		}

		propKidsValType := propertyType(propKidsVal)
		if propKidsValType != array {
			return errors.New("not support /Kids type not array yet")
		}

		pageObjIDs, _, err = readObjIDFromDictionaryArr(propKidsVal)
		fmt.Printf("pageObjIDs = %#v\n%s\\n\n", pageObjIDs, propKidsVal)
		if err != nil {
			return err
		}

	}

	objMustReplaces := make(map[int]string)
	for pageIndex, pageObjID := range pageObjIDs {

		var cw2Content crawl
		//fmt.Printf("cw2Content.set = %d\n\n", pageObjID)
		cw2Content.set(p, pageObjID, "Contents")
		err = cw2Content.run()
		if err != nil {
			return err
		}

		for _, r := range cw2Content.results {

			fmt.Printf("%s\n\n", r.String())

			var propContentsVal string
			//fmt.Printf("id=%d\n", id)
			propContentsVal, err = r.valOf("Contents")
			//fmt.Printf("%d propContentsVal=%s\n\n", id, propContentsVal)
			if err == ErrCrawlResultValOfNotFound {
				continue
			}

			propContentsValType := propertyType(propContentsVal)
			/*if propContentsValType != dictionary {
				return errors.New("not support /Contents type not dictionary yet")
			}*/
			var contentsObjID int
			if propContentsValType == dictionary {
				contentsObjID, _, err = readObjIDFromDictionary(propContentsVal)
				if err != nil {
					return err
				}
			} else if propContentsValType == array {
				contentsObjIDs, _, err := readObjIDFromDictionaryArr(propContentsVal)
				if err != nil || len(contentsObjIDs) <= 0 {
					return err
				}
				contentsObjID = contentsObjIDs[0]
			} else {
				return errors.New("not support /Contents type not dictionary,array yet")
			}

			data := &p.getObjByID(contentsObjID).data
			zip := true
			propContentsObj, err := readProperty(data, "FlateDecode")
			if err != nil {
				return err
			}
			if propContentsObj == nil {
				zip = false
			}

			var stm *bytes.Buffer
			//fmt.Printf("\n-------------------%d-----------------------\n%s\n\n", contentsObjID, string(*data))
			stmLen, err := streamLength(p, data)
			if err != nil {
				return err
			}

			stm, err = extractStream(data, stmLen, zip)
			if err != nil {
				return err
			}

			if _, ok := pageBuffs[pageIndex+1]; ok {
				stm.WriteString("\n")
				pageBuffs[pageIndex+1].WriteTo(stm)
				objMustReplaces[contentsObjID] = fmt.Sprintf("<<\n/Length %d\n>>\nstream\n%sendstream", stm.Len(), stm.String())
			}

		}
	}

	for objID := range objMustReplaces {
		//_ = objID
		fmt.Printf("objID=%d\n", objID)
		p.getObjByID(objID).data = []byte("" + objMustReplaces[objID] + "")
		//fmt.Printf("objId=%d %s\n", objID, string(p.getObjByID(objID).data))
	}

	return nil
}

func streamLength(p *PDFData, data *[]byte) (int, error) {

	prop, err := readProperty(data, "Length")
	if err != nil {
		return 0, err
	}
	if prop == nil {
		prop, err = readProperty(data, "Length1")
		if err != nil {
			return 0, err
		}
		if prop == nil {
			return 0, errors.New("/Length or /Length1 not found")
		}
	}

	propType := prop.valType()
	if propType == number {
		return strconv.Atoi(prop.rawVal)
	} else if propType == dictionary {
		objID, _, err := prop.asDictionary()
		if err != nil {
			return 0, err
		}
		fontlengthObj := p.getObjByID(objID)
		return strconv.Atoi(strings.TrimSpace(string(fontlengthObj.data)))
	} else {
		return 0, errors.New("/Length or /Length1  wrong type")
	}

}

var extractStreamBytes = []byte{0x73, 0x74, 0x72, 0x65, 0x61, 0x6D}

func extractStream(b *[]byte, length int, zip bool) (*bytes.Buffer, error) {

	index := bytes.Index(*b, extractStreamBytes)
	offset := len(extractStreamBytes)
	tmp := (*b)[index+offset:]
	tmp = bytes.TrimSpace(tmp)
	tmp = tmp[0:length]
	var buff bytes.Buffer
	buff.Write(tmp)
	if !zip {
		return &buff, nil
	}
	r, err := zlib.NewReader(&buff)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func findMaxFontIndex(cw *crawl, p *PDFData) (int, error) {
	//find max font index
	max := 0
	for _, item := range cw.results {
		fontPropVal, err := item.valOf("Font")
		if err == ErrCrawlResultValOfNotFound {
			continue
		} else if err != nil {
			return 0, err
		}

		var crFonts crawlResultFonts
		fontPropValType := propertyType(fontPropVal)
		if fontPropValType == object {
			tmp := []byte(fontPropVal)
			err = crFonts.parse(&tmp)
			if err != nil {
				return 0, err
			}
			//fmt.Printf("%#v\n", crFonts)
		} else if fontPropValType == dictionary {
			var fontObjID int
			fontObjID, _, err = readObjIDFromDictionary(fontPropVal)
			if err != nil {
				return 0, err
			}
			fontObj := p.getObjByID(fontObjID)
			err = crFonts.parse(&fontObj.data)
			if err != nil {
				return 0, err
			}
			//fmt.Printf("%#v\n", crFonts)
		} else {
			return 0, errors.New("not support /Font type array yet")
		}

		maxFontIndex := crFonts.maxFontIndex()
		if maxFontIndex > max {
			max = maxFontIndex
		}
	}

	return max, nil
}

func objIDFromStartObjLine(line string) (int, error) {
	tokens := strings.Split(line, " ")
	if len(tokens) < 3 {
		return 0, errors.New("bad start obj")
	}
	id, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, err
	}
	return id, nil
}
