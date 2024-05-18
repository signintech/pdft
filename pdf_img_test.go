package pdft

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"
)

func TestInsertImage(t *testing.T) {

	imgBase64 := "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"
	//imgpath := "./test/img/gopher.png"

	/*_, rawData, err := readImg(imgpath)
	if err != nil {
		t.Error("Couldn't read image")
		return
	}*/

	//filename := "tmpl_empty.pdf"
	//var ipdf PDFt
	//err := ipdf.Open("test/pdf/pdf_from_docx_with_f.pdf")

	filename := "pdf_from_docx_with_f.pdf"
	var ipdf PDFt
	//err := ipdf.Open("test/pdf/pdf_from_docx_with_f.pdf")
	err := ipdf.Open("test/pdf/" + filename)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	//ipdf.InsertImg(rawData, 1, 100.0, 100.0, 100, 100)
	ipdf.InsertImgBase64(imgBase64, 1, 100.0, 100.0, 100, 100)
	err = ipdf.Save("test/out/out3_" + filename)
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
