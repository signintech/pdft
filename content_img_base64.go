package pdft

import (
	"bytes"
	"fmt"
)

type contentImgBase64 struct {
	base64    string
	pageNum   int
	x         float64
	y         float64
	w         float64
	h         float64
	refPdfimg *PDFImageData
}

func (c *contentImgBase64) page() int {
	return c.pageNum
}

func (c *contentImgBase64) toSteram() (*bytes.Buffer, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /%s%d Do Q\n", c.w, c.h, c.x, pageHeight()-(c.y+c.h), c.refPdfimg.xObjChar, c.refPdfimg.xObjIndex))
	return &buff, nil
}
