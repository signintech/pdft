package test

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	_ "image/png"

	"github.com/signintech/pdft"
)

func TestInsertImgDocx(t *testing.T) {
	testWriteImg(t, "./pdf/pdf_from_docx.pdf", "./out/pdf_from_docx_out.pdf")
}

func TestInsertImgGopdf(t *testing.T) {
	testWriteImg(t, "./pdf/pdf_from_gopdf.pdf", "./out/pdf_from_gopdf_out.pdf")
}

func _TestInsertImgChromeLinuxPdf(t *testing.T) {
	//FIXME: test file
	testWriteImg(t, "./pdf/pdf_from_chrome_50_linux64.pdf", "./out/pdf_from_chrome_50_linux64_out.pdf")
}

func _TestInsertImgChromeWin10Pdf(t *testing.T) {
	//FIXME: test file
	testWriteImg(t, "./pdf/pdf_from_chrome_50_win10.pdf", "./out/pdf_from_chrome_50_win10_out.pdf")
}

func TestInsertImgWord2013Pdf(t *testing.T) {
	testWriteImg(t, "./pdf/pdf_from_word2013.pdf", "./out/pdf_from_word2013_out.pdf")
}

func TestInsertImgWord2010Pdf(t *testing.T) {
	testWriteImg(t, "./pdf/pdf_from_word2010.pdf", "./out/pdf_from_word2010_out.pdf")
}

func TestInsertImgRdlcPdf(t *testing.T) {
	testWriteImg(t, "./pdf/pdf_from_rdlc.pdf", "./out/pdf_from_rdlc_out.pdf")
}

func testWriteImg(t *testing.T, source string, target string) {

	signature := "./img/gopher2.jpg"

	var ipdf pdft.PDFt
	err := ipdf.Open(source) //open source PDF file
	if err != nil {
		t.Error("Couldn't open pdf.")
		return
	}
	ipdf.AddFont("arial", "./ttf/arial.ttf")

	_, rawData, err := readImg(signature)
	if err != nil {
		t.Error("Couldn't read image")
		return
	}

	//insert image (support only jpg)
	err = ipdf.InsertImg(rawData, 1, 100.0, 100.0, 100, 100)
	if err != nil {
		t.Errorf("Couldn't insert image %+v", err)
		return
	}

	//insert text
	ipdf.SetFont("arial", "", 14)
	err = ipdf.Insert("Hello PDF", 1, 10.0, 10.0, 100, 100, pdft.Left|pdft.Top)
	if err != nil {
		t.Errorf("Couldn't insert text %+v", err)
		return
	}

	//save to target PDF file
	err = ipdf.Save(target)
	if err != nil {
		t.Errorf("Couldn't save pdf. %+v", err)
		return
	}
}

func readImg(path string) (string, []byte, error) {

	f, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", nil, err
	}

	encoded64 := base64.StdEncoding.EncodeToString(data)
	return encoded64, data, nil
}
