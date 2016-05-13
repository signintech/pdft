package pdft

import (
	"bytes"

	"github.com/signintech/gopdf"
)

//ContentText text in pdf
type ContentText struct {
	text        string
	fontName    string
	fontStyle   string
	fontSize    int
	pageNum     int
	x           float64
	y           float64
	pdfFontData *PDFFontData
}

func (c *ContentText) toSteram() (*bytes.Buffer, error) {

	var cc gopdf.CacheContent
	cc.Setup(nil,
		gopdf.Rgb{},
		1.0,
		c.pdfFontData.fontIndex(),
		c.fontSize,
		c.fontStyle,
		0,
		c.x,
		c.y,
		&c.pdfFontData.font,
		841.89,
		gopdf.ContentTypeText,
		gopdf.CellOption{},
		0.0,
	)
	cc.WriteTextToContent(c.text)
	buff, err := cc.ToStream()
	if err != nil {
		return nil, err
	}
	buff.Write([]byte("\r\n"))

	return buff, nil
}

func (c *ContentText) page() int {
	return c.pageNum
}
