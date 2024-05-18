package pdft

import (
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

func pdf(t *testing.T, filename string) {
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

	err = ipdf.Insert("ที่ กั้น", 1, 10, 10, 100, 100, gopdf.Center|gopdf.Bottom, nil)
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

func TestMeasureTextWidth(t *testing.T) {

	filename := "pdf_from_docx_with_f.pdf"
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

	width1, err := ipdf.MeasureTextWidth("การปั้นโต้ง")
	if err != nil {
		t.Error(err)
		return
	}

	width2, err := ipdf.MeasureTextWidth("การปั้นโต้งการปั้นโต้ง")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.SetFont("arial", "", 15)
	if err != nil {
		t.Error(err)
		return
	}

	width3, err := ipdf.MeasureTextWidth("การปั้นโต้ง")
	if err != nil {
		t.Error(err)
		return
	}

	if width1*2 != width2 {
		t.Error("width1 * 2 != width2")
	}

	if width1 == width3 {
		t.Error("width1 == width3")
	}

	//gitfmt.Printf("%f %f %f\n", width1, width2, width3)
}

/*
func TestSlice(t *testing.T) {
	src := []byte("ABCDEFGHIJ")
	dest := src[9:10]
	fmt.Printf("src=%d dest=%d\n", cap(src), cap(dest))
	dest = append(dest, []byte("1234567890123456789012345678901")...)
	fmt.Println(string(src))
}*/

func TestRemoveOtherPages(t *testing.T) {
	filename := "pdf_from_docx_with_f.pdf"
	var ipdf PDFt
	err := ipdf.Open("test/pdf/" + filename)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.RemoveOtherPages(1)
	if err != nil {
		t.Error(err)
		return
	}
	err = ipdf.Save("test/out/RemoveOtherPages_" + filename)
	if err != nil {
		t.Error(err)
		return
	}
}
