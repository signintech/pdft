package pdft

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type crawlResultFonts []crawlResultFont

var ErrNoFontIndex = fmt.Errorf("no font index")

// Split a font string like '/F1' or '/TT2' into the prefix and index
func splitFont(s string) (string, int, error) {
	prefix := ""
	for i, r := range s {
		if unicode.IsDigit(r) {
			s := strings.ReplaceAll(s, "_", "")
			val, err := strconv.Atoi(s[i:])
			if err != nil {
				return "", 0, err
			}
			return prefix, val, nil
		}
		prefix += string(r)
	}
	return "", 0, ErrNoFontIndex
}

func (c *crawlResultFonts) parse(propVal *[]byte) error {
	var props PDFObjPropertiesData
	err := readProperties(propVal, &props)
	if err != nil {
		return err
	}

	for _, prop := range props {
		var crFont crawlResultFont

		prefix, fontIndex, err := splitFont(strings.TrimSpace(prop.key))
		if err != nil {
			return err
		}
		objID, _, err := prop.asDictionary()
		if err != nil {
			return err
		}
		crFont.prefix = prefix
		crFont.fontIndex = fontIndex
		crFont.fontObjID = objID
		*c = append(*c, crFont)
	}
	return nil
}

func (c *crawlResultFonts) String() string {
	var buff bytes.Buffer
	for _, f := range *c {
		buff.WriteString(fmt.Sprintf("/%s%d %d 0 R\n", f.prefix, f.fontIndex, f.fontObjID))
	}
	return buff.String()
}

func (c *crawlResultFonts) append(fontIndex int, fontObjID int) {
	var crFont crawlResultFont
	crFont.fontIndex = fontIndex
	crFont.fontObjID = fontObjID
	crFont.prefix = "F"
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
	prefix    string
}
