package render

import (
	"testing"

	"github.com/signintech/pdft"
)

func TestRender(t *testing.T) {
	var fis FieldInfos
	fis = append(fis, FieldInfo{
		PageNum: 1,
		Key:     "fname",
		X:       40,
		Y:       212,
		W:       65,
		H:       10,
		Align:   pdft.Bottom | pdft.Center,
	})

	fis = append(fis, FieldInfo{
		PageNum: 1,
		Key:     "lname",
		X:       115,
		Y:       212,
		W:       80,
		H:       10,
		Align:   pdft.Bottom | pdft.Left,
	})

	rd, err := NewRender("../test/pdf/pdf_from_word2013_b.pdf", fis)
	if err != nil {
		t.Error(err)
		return
	}
	//rd.ShowCellBorder(true)
	err = rd.AddFont("arial", "../test/ttf/Loma.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = rd.SetFont("arial", "", 9)
	if err != nil {
		t.Error(err)
		return
	}

	rd.ShowDesignView()
	//rd.Text("fname", "วันพรรษ")
	//rd.Text("lname", "อนันตพันธ์ xxxxxxxxxxxxxxxxx")
	err = rd.Save("../test/out/render01.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}
