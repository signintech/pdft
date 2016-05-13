package pdft

import (
	"fmt"
	"testing"
)

func TestPDF(t *testing.T) {
	//pdf(t, "pdf_from_gopdf.pdf")
	//pdf(t, "pdf_from_docx.pdf")
	//pdf(t, "pdf_from_docx_with_f.pdf")
	//pdf(t, "pdf_from_iia.pdf")
	//pdf(t, "pdf_from_delphi.pdf")
	//pdf(t, "pdf_from_word2010.pdf")
	//pdf(t, "pdf_from_word2010_b.pdf")
	//pdf(t, "pdf_from_chrome_50_win10.pdf")
	//pdf(t, "pdf_from_word2013.pdf")
	//pdf(t, "pdf_from_word2013_b.pdf")
	//pdf(t, "pdf_from_rdlc.pdf")
}

func pdf(t *testing.T, filename string) {
	fmt.Printf("\n\n\n####Open %s ####\n\n", filename)
	var ipdf InjPDF
	//err := ipdf.Open("test/pdf/pdf_from_docx_with_f.pdf")
	err := ipdf.Open("test/pdf/" + filename)
	//err := ipdf.Open("test/pdf/pdf_from_docx.pdf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.AddFont("arial", "test/ttf/arial.ttf")
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
	}

	err = ipdf.Save("test/out/" + filename)
	if err != nil {
		t.Error(err)
		return
	}
}
