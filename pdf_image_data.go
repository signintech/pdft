package pdft

import (
	"bytes"
	"encoding/base64"

	gopdf "github.com/kelvinsantos/pdft/minigopdf"
)

// PDFImageData pdf img
type PDFImageData struct {
	objID     int
	imgObj    gopdf.ImageObj
	xObjChar  string
	xObjIndex int
}

/*
func (p *PDFImageData) setImgBase64(base64str string) error {

	r := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(base64str)))
	err := p.imgObj.SetImage(r)
	if err != nil {
		return err
	}

	err = p.imgObj.Parse()
	if err != nil {
		return err
	}

	err = p.imgObj.Build(p.objID)
	if err != nil {
		return err
	}

	return nil
}*/

func createImgObjFromImgBase64(base64str string) (gopdf.ImageObj, *gopdf.SMask, error) {
	r := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(base64str)))
	var imgObj gopdf.ImageObj
	err := imgObj.SetImage(r)
	if err != nil {
		return gopdf.ImageObj{}, nil, err
	}
	err = imgObj.Parse()
	if err != nil {
		return gopdf.ImageObj{}, nil, err
	}
	smask, err := imgObj.CreateSMask()
	if err != nil {
		return gopdf.ImageObj{}, nil, err
	}
	return imgObj, smask, nil
}

func createImgObjFromBytes(img []byte) (gopdf.ImageObj, *gopdf.SMask, error) {
	buff := bytes.NewBuffer(img)
	var imgObj gopdf.ImageObj
	err := imgObj.SetImage(buff)
	if err != nil {
		return gopdf.ImageObj{}, nil, err
	}
	err = imgObj.Parse()
	if err != nil {
		return gopdf.ImageObj{}, nil, err
	}
	smask, err := imgObj.CreateSMask()
	if err != nil {
		return gopdf.ImageObj{}, nil, err
	}
	return imgObj, smask, nil
}

func (p *PDFImageData) setImgObj(imgObj gopdf.ImageObj) error {
	p.imgObj = imgObj
	err := p.imgObj.Build(p.objID)
	if err != nil {
		return err
	}
	return nil
}

/*
func (p *PDFImageData) setImg(img []byte) error {

	buff := bytes.NewBuffer(img)
	err := p.imgObj.SetImage(buff)
	if err != nil {
		return err
	}

	err = p.imgObj.Parse()
	if err != nil {
		return err
	}

	err = p.imgObj.Build(p.objID)
	if err != nil {
		return err
	}

	return nil
}
*/
