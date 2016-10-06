package pdft

import (
	"fmt"
	"testing"

	"github.com/signintech/gopdf"
)

func TestPDF(t *testing.T) {
	//pdf(t, "pdf_from_gopdf.pdf")
	pdf(t, "pdf_from_docx.pdf")
	//pdf(t, "pdf_from_docx_with_f.pdf")
	//pdf(t, "pdf_from_iia.pdf")
	//pdf(t, "pdf_from_delphi.pdf")
	//pdf(t, "pdf_from_word2010.pdf")
	//pdf(t, "pdf_from_word2010_b.pdf")
	//pdf(t, "pdf_from_chrome_50_win10.pdf")
	//pdf(t, "pdf_from_chrome_50_linux64.pdf")
	//pdf(t, "pdf_from_word2013.pdf")
	//pdf(t, "pdf_from_word2013_b.pdf")
	//pdf(t, "pdf_from_rdlc.pdf")
}

func pdf(t *testing.T, filename string) {
	fmt.Printf("\n\n\n####Open %s ####\n\n", filename)
	var ipdf PDFt
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

	err = ipdf.Insert("hi", 1, 10, 10, 100, 100, gopdf.Center|gopdf.Bottom)
	if err != nil {
		t.Error(err)
		return
	}

	ipdf.SetProtection(
		gopdf.PermissionsPrint|gopdf.PermissionsCopy|gopdf.PermissionsModify,
		[]byte("1234"),
		[]byte("5555"),
	)
	err = ipdf.Save("test/out/" + filename)
	if err != nil {
		t.Error(err)
		return
	}
}
