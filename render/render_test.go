package render

import (
	"testing"

	"github.com/signintech/pdft"
)

func TestRender(t *testing.T) {
	var fis FieldInfos
	fis = append(fis, FieldInfo{
		Key:   "fname",
		X:     10,
		Y:     11,
		W:     100,
		H:     100,
		Align: pdft.Middle,
	})

	fis = append(fis, FieldInfo{
		Key:   "lname",
		X:     10,
		Y:     11,
		W:     100,
		H:     100,
		Align: pdft.Middle,
	})

	rd, err := NewRender("../pdf/pdf_from_docx.pdf", fis)
	if err != nil {
		t.Error(err)
		return
	}
	rd.WriteText("fname", "วันพรรษ")
	rd.WriteText("lname", "อนันตพันธ์")
}
