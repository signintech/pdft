package pdft

import (
	"errors"
	"io"
	"io/ioutil"
)

//PDFParse parse pdf
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

	return nil
}

func parseTrailer(raw *[]byte, pdf *PDFData) error {
	indexTrailers := regexpTrailer.FindAllIndex(*raw, 1)
	indexStartxrefs := regexpStartxref.FindAllIndex(*raw, 1)
	if len(indexTrailers) <= 0 {
		return errors.New("trailer not found")
	}
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
		if xref.Keyword == "n" {
			var pdfobj PDFObjData
			err := pdfobj.parse(raw, xref.N10Digit)
			if err != nil {
				return err
			}
			pdf.put(pdfobj)
		}
	}

	return nil
}
