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
	ipdf.AddFont("angsa", "./ttf/angsa.ttf")

	_, rawData, err := readImg(signature)
	if err != nil {
		t.Error("Couldn't read image")
		return
	}

	angsatextrisefn := func(
		leftRune rune,
		rightRune rune,
		fontSize int,
		allText string,
		currTextIndex int,
	) float32 {
		gap := float32(0)
		if rightRune == '่' || rightRune == '้' || rightRune == '๊' || rightRune == '๋' {
			if leftRune == 'ิ' || leftRune == 'ี' ||
				leftRune == 'ึ' || leftRune == 'ื' ||
				leftRune == 'ั' {
				gap = 9.5
			} else if leftRune == 'ป' || leftRune == 'ฬ' || leftRune == 'ฝ' || leftRune == 'ฟ' {
				gap = 3
			} else {
				runes := []rune(allText)
				if len(runes) > currTextIndex+1 {
					nextRune := runes[currTextIndex+1]
					if nextRune == 'ำ' {
						gap = 9.5
					}
				}
			}

		}

		return gap * (float32(fontSize) / 40.0)
	}

	ipdf.TextriseOverride("angsa", angsatextrisefn)

	//insert image (support only jpg)
	err = ipdf.InsertImg(rawData, 1, 100.0, 100.0, 100, 100)
	if err != nil {
		t.Errorf("Couldn't insert image %+v", err)
		return
	}

	//insert text
	ipdf.SetFont("angsa", "", 14)
	err = ipdf.Insert("น้ำอยู่บนฟ้าเป็นละอองเป็นไอฉันจะบินไปหาบนนภาทันใดแล้วโปรยเกลือผงหรือโซเดียมคลอไรด์มันจะดูดความชื้นเป็นหยดน้ำเป็นก้อนเมฆ", 1, 10.0, 10.0, 100, 100, pdft.Left|pdft.Top, nil)
	if err != nil {
		t.Errorf("Couldn't insert text %+v", err)
		return
	}

	ipdf.SetFont("angsa", "", 14)
	err = ipdf.Insert("นาคีมีพิษเพี้ยง สุริโย เลื้อยบ่ทำเดโช แช่มช้า พิษน้อยหยิ่งโยโส แมงป่อง ชูแต่หางเองอ้า อวดอ้าง ฤทธี", 1, 10.0, 40.0, 100, 100, pdft.Left|pdft.Top, nil)
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
