package test

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	_ "image/png"

	"github.com/signintech/pdft"
)

func TestInsertToDocx(t *testing.T) {

	if err := initTest(); err != nil {
		t.Errorf("%+v", err)
		return
	}

	writePdf(t, "./pdf/pdf_from_docx.pdf", "./out/pdf_from_docx_out.pdf")
}

func TestInsertToGopdf(t *testing.T) {

	if err := initTest(); err != nil {
		t.Errorf("%+v", err)
		return
	}

	writePdf(t, "./pdf/pdf_from_gopdf.pdf", "./out/pdf_from_gopdf_out.pdf")
}

func _TestInsertToChromeLinuxPdf(t *testing.T) {

	if err := initTest(); err != nil {
		t.Errorf("%+v", err)
		return
	}

	//FIXME: test file
	writePdf(t, "./pdf/pdf_from_chrome_50_linux64.pdf", "./out/pdf_from_chrome_50_linux64_out.pdf")
}

func _TestInsertToChromeWin10Pdf(t *testing.T) {

	if err := initTest(); err != nil {
		t.Errorf("%+v", err)
		return
	}
	//FIXME: test file
	writePdf(t, "./pdf/pdf_from_chrome_50_win10.pdf", "./out/pdf_from_chrome_50_win10_out.pdf")
}

func TestInsertToWord2013Pdf(t *testing.T) {

	if err := initTest(); err != nil {
		t.Errorf("%+v", err)
		return
	}

	writePdf(t, "./pdf/pdf_from_word2013.pdf", "./out/pdf_from_word2013_out.pdf")
}

func TestInsertToWord2010Pdf(t *testing.T) {
	writePdf(t, "./pdf/pdf_from_word2010.pdf", "./out/pdf_from_word2010_out.pdf")
}

func TestInsertToRdlcPdf(t *testing.T) {
	writePdf(t, "./pdf/pdf_from_rdlc.pdf", "./out/pdf_from_rdlc_out.pdf")
}

func TestInsertToWithImgPdf(t *testing.T) {

	if err := initTest(); err != nil {
		t.Errorf("%+v", err)
		return
	}

	writePdf(t, "./pdf/pdf_with_img.pdf", "./out/pdf_with_img_out.pdf")
}

func writePdf(t *testing.T, source string, target string) {

	signature := "./img/gopher2.jpg"

	var ipdf pdft.PDFt
	err := ipdf.Open(source) //open source PDF file
	if err != nil {
		t.Error("Couldn't open pdf.")
		return
	}
	ipdf.AddFont("arial", "./ttf/angsa.ttf")

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
	err = ipdf.Insert("Hello PDF  กั้น ชั้น ที่", 1, 10.0, 10.0, 100, 100, pdft.Left|pdft.Top)
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

func initTest() error {
	isExists, err := exists("./out/")
	if err != nil {
		return errors.New("can not create out/ for test")
	}
	if !isExists {
		err := os.MkdirAll("./out/", 0777)
		if err != nil {
			return errors.New("can not create out/ for test")
		}
	}
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
