package pdft

import (
	"fmt"
	"testing"

	gopdf "github.com/signintech/pdft/minigopdf"
)

func TestPDF(t *testing.T) {
	pdf(t, "pdf_from_gopdf.pdf")
	//pdf(t, "pdf_from_docx.pdf")
	pdf(t, "pdf_from_docx_with_f.pdf")
	//pdf(t, "pdf_from_iia.pdf")
	//pdf(t, "pdf_from_delphi.pdf")
	//pdf(t, "pdf_from_word2010.pdf")
	//pdf(t, "pdf_from_word2010_b.pdf")
	//pdf(t, "pdf_from_chrome_50_win10.pdf")
	//pdf(t, "pdf_from_chrome_50_linux64.pdf")
	pdf(t, "pdf_from_word2013.pdf")
	//pdf(t, "pdf_from_word2013_b.pdf")
	//pdf(t, "pdf_from_rdlc.pdf")
}

func TestPDFLeluo(t *testing.T) {
	path := "/Users/oneplus/test/1.pdf"
	var ipdf PDFt
	//err := ipdf.Open("test/pdf/pdf_from_docx_with_f.pdf")
	err := ipdf.Open(path)
	//err := ipdf.Open("test/pdf/pdf_from_docx.pdf")
	if err != nil {
		t.Error(err)
		return
	}
	/*
		err = ipdf.AddFont("songti", "/Users/oneplus/test/songti.ttf")
		if err != nil {
			t.Error(err)
			return
		}

		err = ipdf.SetFont("songti", "", 14)
		if err != nil {
			t.Error(err)
			return
		}*/
	err = ipdf.AddFont("arial", "test/ttf/angsa.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.SetFont("arial", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.Insert("ABC", 1, 0, 0, 100, 100, gopdf.Left|gopdf.Bottom)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.Save("test/out/leluo.pdf")
	if err != nil {
		t.Error(err)
		return
	}
	//pdf(t, "pdf_from_docx.pdf")
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

func TestGetFontIndex(t *testing.T) {
	//crawlResultFonts
	var c crawlResultFonts
	_, v, err := c.getFontIndex("C2_12")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if v != 12 {
		t.Fatalf("%d", v)
	}
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

	err = ipdf.AddFont("arial", "test/ttf/angsa.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.SetFont("arial", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.Insert("ที่ กั้น", 1, 10, 10, 100, 100, gopdf.Center|gopdf.Bottom)
	if err != nil {
		t.Error(err)
		return
	}
	/*
		ipdf.SetProtection(
			gopdf.PermissionsPrint|gopdf.PermissionsCopy|gopdf.PermissionsModify,
			[]byte("1234"),
			[]byte("5555"),
		)*/
	err = ipdf.Save("test/out/" + filename)
	if err != nil {
		t.Error(err)
		return
	}
}

/*
func TestSlice(t *testing.T) {
	src := []byte("ABCDEFGHIJ")
	dest := src[9:10]
	fmt.Printf("src=%d dest=%d\n", cap(src), cap(dest))
	dest = append(dest, []byte("1234567890123456789012345678901")...)
	fmt.Println(string(src))
}*/
