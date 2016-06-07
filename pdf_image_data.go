package pdft

import (
	"bytes"
	"encoding/base64"

	"github.com/signintech/gopdf"
)

//PDFImageData pdf img
type PDFImageData struct {
	objID     int
	imgObj    gopdf.ImageObj
	xObjChar  string
	xObjIndex int
}

func (p *PDFImageData) setImgBase64(base64str string) error {

	r := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(base64str)))
	err := p.imgObj.SetImage(r)
	if err != nil {
		return err
	}

	err = p.imgObj.Build()
	if err != nil {
		return err
	}

	return nil
}
