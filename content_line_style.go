package pdft

import (
	"bytes"
	"fmt"
)

//ContentLineStyle set line width
type ContentLineStyle struct {
	width    float64
	lineType string
}

func (c *ContentLineStyle) toSteram() (*bytes.Buffer, error) {
	var buff bytes.Buffer
	switch c.lineType {
	case "dashed":
		buff.WriteString(fmt.Sprint("[5] 2 d\n"))
	case "dotted":
		buff.WriteString(fmt.Sprint("[2 3] 11 d\n"))
	default:
		buff.WriteString(fmt.Sprint("[] 0 d\n"))
	}
	buff.WriteString(fmt.Sprintf("%.2f w\n", c.width))
	return &buff, nil
}

func (c *ContentLineStyle) page() int {
	return 1
}
