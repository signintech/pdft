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
	trailer  TrailerData
	xrefs    []XrefData
	objIDs   []int
	objs     []PDFObjData
	pagesObj *PDFObjData
}

//TrailerData trailer
type TrailerData struct {
	rootObjID int
}

//Len count
func (p *PDFData) Len() int {
	return len(p.objIDs)
}

func (p *PDFData) put(pdfobj PDFObjData) {
	p.objIDs = append(p.objIDs, pdfobj.objID)
	p.objs = append(p.objs, pdfobj)
}

func (p *PDFData) putNewObject(pdfobj PDFObjData) int {
	newObjID := p.maxID() + 1
	pdfobj.objID = newObjID
	p.put(pdfobj)
	return newObjID
}

func (p *PDFData) removeObjByID(objID int) error {
	for i, id := range p.objIDs {
		if id == objID {
			p.objIDs = append(p.objIDs[:i], p.objIDs[i+1:]...)
			p.objs = append(p.objs[:i], p.objs[i+1:]...)
			return nil
		}
	}
	return errors.New("Not Found")
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

// getPagesObjID return number of page of the pdf
func (p *PDFData) getPageObjIDs() ([]int, error) {
	pagesProp, err := p.pagesObj.readProperties()
	if err != nil {
		return nil, err
	}
	pagesKids := pagesProp.getPropByKey("Kids")
	if pagesKids == nil {
		return nil, errors.New("Not found Kids property in this object")
	}
	listPagesObjID, _, err := pagesKids.asDictionaryArr()
	if err != nil {
		return nil, err
	}
	return listPagesObjID, err
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

func (p *PDFData) injectImgsToPDF(pdfImgs []*PDFImageData) error {

	var err error
	isEmbedResources := false
	rootOfXObjectID := -1
	resourcesContent := ""
	cp := crawlPages{}
	cwRes, _ := cp.getPageCrawl(p, p.trailer.rootObjID, "Kids", "Resources")
	if err != nil {
		return err
	}
	foundRes := false
	for resID, r := range cwRes.results {
		resources, err := r.valOf("Resources")
		if err == ErrCrawlResultValOfNotFound {
			continue
		} else if err != nil {
			return err
		} else {
			foundRes = true
			resourcesID, _, err := readObjIDFromDictionary(resources)
			if err == ErrorObjectIDNotFound {
				rootOfXObjectID = resID
				resourcesContent = resources
				isEmbedResources = true
			} else if err != nil {
				return err
			} else {
				rootOfXObjectID = resourcesID
				data := p.getObjByID(resourcesID)
				if data != nil {
					resourcesContent = string(data.data)
				}
				isEmbedResources = false
			}
			break
		}
	}

	if !foundRes {
		return errors.New("not found /Resources in /Type/Pages")
	}

	var cw crawl
	//cw.set(p, p.trailer.rootObjID, "Pages", "Kids", "Resources", "XObject")
	cw.set(p, rootOfXObjectID, "XObject")
	err = cw.run()
	if err != nil {
		return err
	}

	found := false
	xObjectVals := make(map[int]string)
	for objID, r := range cw.results {
		xobject, err := r.valOf("XObject")
		if err == ErrCrawlResultValOfNotFound {
			continue
		} else if err != nil {
			return err
		} else {
			xObjectVals[objID] = xobject
			found = true
		}
	}

	if !found { //ถ้ายังไม่เจออีก
		cp2 := crawlPages{}
		cw2, _ := cp2.getPageCrawl(p, p.trailer.rootObjID, "Kids", "Resources", "XObject")
		cw = *cw2
		if err != nil {
			return err
		}
		for objID, r := range cw.results {
			xobject, err := r.valOf("XObject")
			if err == ErrCrawlResultValOfNotFound {
				continue
			} else if err != nil {
				return err
			} else {
				xObjectVals[objID] = xobject
				found = true
			}
		}
	}

	var xobjs crawlResultXObjects
	var xObjIndex int
	xObjChar := "I"
	if found {
		for _, xObjectVal := range xObjectVals {
			propVal := []byte(xObjectVal)
			xobjs.parse(&propVal)
			if len(xobjs) > 0 {
				xObjChar = xobjs[len(xobjs)-1].xObjChar
				if xobjs[len(xobjs)-1].xObjIndex > xObjIndex {
					xObjIndex = xobjs[len(xobjs)-1].xObjIndex
				}
			}
		}
	}

	i := 0
	max := len(pdfImgs)
	for i < max {
		objID := pdfImgs[i].objID
		pdfImgs[i].xObjChar = xObjChar
		pdfImgs[i].xObjIndex = xObjIndex + i + 1

		var xobj crawlResultXObject
		xobj.xObjChar = xObjChar
		xobj.xObjIndex = xObjIndex + i + 1
		xobj.xObjObjID = objID
		xobjs = append(xobjs, xobj)
		i++
	}

	objMustReplaces := make(map[int]string)
	if found {
		for objID, r := range cw.results {
			_, err = r.valOf("XObject")
			if err == ErrCrawlResultValOfNotFound {
				continue
			} else if err != nil {
				return err
			}
			r.setValOf("XObject", fmt.Sprintf("<<%s>>\n", xobjs.String()))
			objMustReplaces[objID] = r.String()
		}
	} else {
		if isEmbedResources {
			var cw01 crawl
			cw01.set(p, p.trailer.rootObjID, "Pages", "Kids", "Resources")
			err = cw01.run()
			if err != nil {
				return err
			}
			for objID, r := range cw01.results {
				res, err := r.valOf("Resources")
				if err == ErrCrawlResultValOfNotFound {
					continue
				} else if err != nil {
					return err
				} else {
					res = strings.TrimSpace(res)
					res = fmt.Sprintf("%s /XObject <<%s>>", res[2:len(res)-2], xobjs.String())
					r.setValOf("Resources", fmt.Sprintf("<<%s>>\n", res))
					objMustReplaces[objID] = r.String()
				}
			}
		} else {
			for objID, r := range cw.results {
				res := strings.TrimSpace(resourcesContent)
				res = fmt.Sprintf("<<%s>>\n", xobjs.String())
				r.add("XObject", res)
				objMustReplaces[objID] = r.String()
				//fmt.Printf("%s\n", r.String())
			}
		}
	}

	for objID := range objMustReplaces {
		p.getObjByID(objID).data = []byte("<<\n" + objMustReplaces[objID] + ">>\n")
	}

	return nil
}

func (p *PDFData) injectFontsToPDF(fontDatas map[string]*PDFFontData) error {
	var err error
	cp := crawlPages{}
	cw, _ := cp.getPageCrawl(p, p.trailer.rootObjID, "Kids", "Resources", "Font")
	if err != nil {
		return err
	}

	maxFontIndex, err := findMaxFontIndex(cw, p)
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

		//fmt.Printf("buff=%s\n\n", buff.String())

		_, err = buff.WriteTo(pageBuffs[pageNum])
		if err != nil {
			return err
		}
	}
	cp := crawlPages{}
	cwt, _ := cp.getPageCrawl(p, p.trailer.rootObjID, "Kids", "Parent")
	pageObjIDs, _ := cp.getPageObjIDs(cwt)
	objMustReplaces := make(map[int]string)
	for pageIndex, pageObjID := range pageObjIDs {

		var cw2Content crawl
		cw2Content.set(p, pageObjID, "Contents")
		err = cw2Content.run()
		if err != nil {
			return err
		}

		for _, r := range cw2Content.results {

			//fmt.Printf("%s\n\n", r.String())

			var propContentsVal string
			// fmt.Printf("id=%d\n", id)
			propContentsVal, err = r.valOf("Contents")
			// fmt.Printf("%d propContentsVal=%s\n\n", 0, r.String())
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
			//fmt.Printf("stm=%s\n\n", stm.String())

			if _, ok := pageBuffs[pageIndex+1]; ok {
				stm.WriteString("\n")
				pageBuffs[pageIndex+1].WriteTo(stm)
				objMustReplaces[contentsObjID] = fmt.Sprintf("<<\n/Length %d\n>>\nstream\n%sendstream", stm.Len(), stm.String())
			}

		}
	}

	for objID := range objMustReplaces {
		//_ = objID
		//fmt.Printf("objID=%d\n", objID)
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
		return strconv.Atoi(strings.TrimSpace(prop.rawVal))
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
	id, err := strconv.Atoi(strings.TrimSpace(tokens[0]))
	if err != nil {
		return 0, err
	}
	return id, nil
}
