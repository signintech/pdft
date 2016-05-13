package pdft

import (
	"fmt"
	"testing"
)

func TestPDFParse(t *testing.T) {
	pdfParseChrome50Win10(t)
}

func pdfParseChrome50Win10(t *testing.T) {
	filename := "pdf_from_chrome_50_win10.pdf"
	fmt.Printf("\n\n\n####Open %s ####\n\n", filename)
	var ipdf InjPDF
	//err := ipdf.Open("test/pdf/pdf_from_docx_with_f.pdf")
	err := ipdf.Open("test/pdf/" + filename)
	//err := ipdf.Open("test/pdf/pdf_from_docx.pdf")
	if err != nil {
		t.Error(err)
		return
	}

	var props PDFObjPropertiesData
	pdfObj := ipdf.pdf.getObjByID(84)
	err = readProperties(&pdfObj.data, &props)
	if err != nil {
		t.Error(err)
		return
	}

	//fmt.Printf("props=%#v\n", props)

	/*err = ipdf.AddFont("arial", "test/ttf/arial.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.SetFont("arial", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.InsertText("hi", 1, 100, 100)
	if err != nil {
		t.Error(err)
		return
	}*/

	/*err = ipdf.Save("test/out/" + filename)
	if err != nil {
		t.Error(err)
		return
	}*/

}
