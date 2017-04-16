package test

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	_ "image/png"

	"github.com/signintech/pdft"
)

func TestInsertImg(t *testing.T) {

	source := "./pdf/pdf_from_docx_with_f.pdf"
	//source := "./pdf/pdf_from_docx.pdf"
	//source := "./pdf/pdf_from_gopdf.pdf"
	//signature := "./img/gopher.png"
	signature := "./img/gopher2.jpg"
	target := "./out/pdf_from_docx_out.pdf"

	var ipdf pdft.PDFt
	err := ipdf.Open(source)
	if err != nil {
		t.Error("Couldn't open pdf.")
	}

	encoded, data, err := readImg(signature)
	if err != nil {
		t.Error("Couldn't read image")
	}

	/*err = ipdf.InsertImgBase64(encoded, 1, 100, 200, 100, 100)
	if err != nil {
		t.Error("Couldn't insert image base64")
	}

	*/
	_ = data
	_ = encoded
	err = ipdf.InsertImg(data, 1, 100.0, 100.0, 100, 100)
	if err != nil {
		t.Error("Couldn't insert image")
	}

	err = ipdf.Save(target)
	if err != nil {
		t.Error("Couldn't save pdf.")
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

	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, data, nil
}
