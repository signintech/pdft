package pdft

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type crawlResultFonts []crawlResultFont

func (c *crawlResultFonts) parse(propVal *[]byte) error {
	var props PDFObjPropertiesData
	err := readProperties(propVal, &props)
	if err != nil {
		return err
	}

	for _, prop := range props {
		var crFont crawlResultFont
		fontIndex, err := strconv.Atoi(strings.TrimSpace(strings.Replace(prop.key, "F", "", -1)))
		if err != nil {
			return err
		}
		objID, _, err := prop.asDictionary()
		if err != nil {
			return err
		}
		crFont.fontIndex = fontIndex
		crFont.fontObjID = objID
		*c = append(*c, crFont)
	}
	return nil
}

func (c *crawlResultFonts) String() string {
	var buff bytes.Buffer
	for _, f := range *c {
		buff.WriteString(fmt.Sprintf("/F%d %d 0 R\n", f.fontIndex, f.fontObjID))
	}
	return buff.String()
}

func (c *crawlResultFonts) append(fontIndex int, fontObjID int) {
	var crFont crawlResultFont
	crFont.fontIndex = fontIndex
	crFont.fontObjID = fontObjID
	*c = append(*c, crFont)
}

func (c *crawlResultFonts) maxFontIndex() int {
	max := 0
	for _, f := range *c {
		if f.fontIndex > max {
			max = f.fontIndex
		}
	}
	return max
}

type crawlResultFont struct {
	fontIndex int
	fontObjID int
}
