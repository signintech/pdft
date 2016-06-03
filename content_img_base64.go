package pdft

import "bytes"

type contentImgBase64 struct {
	base64  string
	pageNum int
	x       float64
	y       float64
	w       float64
	h       float64
}

func (c *contentImgBase64) page() int {
	return c.pageNum
}

func (c *contentImgBase64) toSteram() (*bytes.Buffer, error) {

	/*r := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(c.base64)))

	var imgObj gopdf.ImageObj
	err := imgObj.SetImage(r)
	if err != nil {
		return nil, err
	}

	err = imgObj.Build()
	if err != nil {
		return nil, err
	}

	return imgObj.GetObjBuff(), nil*/
	return nil, nil
}
