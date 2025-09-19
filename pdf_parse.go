package pdft

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// PDFParse parse pdf
func PDFParse(file io.Reader, outPdf *PDFData) error {

	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = parseXref(&raw, outPdf)
	if err != nil {
		return err
	}

	err = parseTrailer(&raw, outPdf)
	if err != nil {
		return err
	}

	err = parseObjOnlyKeywordN(&raw, outPdf)
	if err != nil {
		return err
	}

	err = setPagesObj(outPdf)
	if err != nil {
		return err
	}

	err = parsePageSize(outPdf)
	if err != nil {
		return err
	}

	return nil
}

func parsePageSize(pdf *PDFData) error {
	pdf.pageSizes = make(map[int][]float64)
	pageIDs, err := pdf.getPageObjIDs()
	if err != nil {
		return err
	}

	for n, pageID := range pageIDs {
		data := pdf.getObjByID(pageID).data
		propContentsObj, err := readProperty(&data, "MediaBox")
		if err != nil {
			return err
		}
		valType := propContentsObj.valType()
		if valType != array {
			fmt.Printf("MediaBox type is not array, but %s", valType)
			continue
		}
		parseMediaBoxd, err := parseFloatSlice(propContentsObj.rawVal)
		if err != nil {
			return err
		}
		pdf.pageSizes[n] = parseMediaBoxd
	}
	return nil
}

func parseFloatSlice(s string) ([]float64, error) {
	// ตัด [] ออก
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	// แยกด้วยช่องว่าง (จัดการ space เกินด้วย Fields)
	parts := strings.Fields(s)

	var result []float64
	for _, p := range parts {
		f, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}
	return result, nil
}

func setPagesObj(pdf *PDFData) error {
	pagesProps, err := pdf.getObjByID(pdf.trailer.rootObjID).readProperties()
	if err != nil {
		return err
	}
	pagesProp := pagesProps.getPropByKey("Pages")
	if pagesProp == nil {
		return errors.New("/Pages not found")
	}
	pagesID, _, err := pagesProp.asDictionary()
	if err != nil {
		return err
	}

	pagesObj := pdf.getObjByID(pagesID)
	if err != nil {
		return err
	}
	pdf.pagesObj = pagesObj
	return nil
}

func parseTrailer(raw *[]byte, pdf *PDFData) error {

	indexTrailers := regexpTrailer.FindAllIndex(*raw, 1)
	if len(indexTrailers) <= 0 {
		return errors.New("trailer not found")
	}

	indexStartxrefs := regexpStartxref.FindAllIndex(*raw, 1)
	if len(indexStartxrefs) <= 0 {
		return errors.New("startxref not found")
	}

	tmp := (*raw)[indexTrailers[0][1]:indexStartxrefs[0][0]]
	var props PDFObjPropertiesData
	err := readProperties(&tmp, &props)
	if err != nil {
		return err
	}

	propRoot := props.getPropByKey("Root")
	if propRoot == nil {
		return errors.New("/Root not found")
	}

	rootID, _, err := propRoot.asDictionary()
	if err != nil {
		return err
	}
	pdf.trailer.rootObjID = rootID

	return nil
}

func parseXref(raw *[]byte, pdf *PDFData) error {

	indexXrefs := regexpXref.FindAllIndex(*raw, 1)
	if len(indexXrefs) <= 0 {
		return errors.New("xref not found")
	}
	startXrefAt := indexXrefs[0][1]
	xrefLines := regexpXrefLine.FindAllIndex(*raw, -1)
	for _, xrefLine := range xrefLines {
		if startXrefAt > xrefLine[0] { //xref must below "xref" string
			continue
		}
		var xref XrefData
		err := xref.parse(string((*raw)[xrefLine[0]:xrefLine[1]]))
		if err != nil {
			return err
		}
		pdf.xrefs = append(pdf.xrefs, xref)
	}

	return nil
}

func parseObjOnlyKeywordN(raw *[]byte, pdf *PDFData) error {
	for _, xref := range pdf.xrefs {
		//fmt.Printf("%d\n", i)
		if xref.Keyword == "n" {
			var pdfobj PDFObjData
			err := pdfobj.parse(raw, xref.N10Digit)
			if err != nil {
				return err
			}
			pdf.put(pdfobj)
		}
	}
	//fmt.Printf("parseObjOnlyKeywordN = %d\n", pdf.Len())
	return nil
}
