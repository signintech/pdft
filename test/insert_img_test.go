package test

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/signintech/pdft"
)

func TestInsertImg(t *testing.T) {
	source := "./pdf/pdf_from_docx.pdf"
	signature := "./img/avatar.jpg"
	target := "./out/pdf_from_docx_out.pdf"
	payload, err := ioutil.ReadFile(signature)
	if err != nil {
		t.Error("Couldn't read signature.")
	}

	var ipdf pdft.PDFt
	err = ipdf.Open(source)
	if err != nil {
		t.Error("Couldn't open pdf.")
	}

	encoded := base64.StdEncoding.EncodeToString(payload)
	//fmt.Println(encoded)

	err = ipdf.InsertImgBase64(encoded, 1, 182.0, 165.0, 172.0, 49.0)
	if err != nil {
		t.Error("Couldn't insert image")
	}

	err = ipdf.Save(target)
	if err != nil {
		t.Error("Couldn't save pdf.")
	}
}
