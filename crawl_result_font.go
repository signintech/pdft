package pdft

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	gopdf "github.com/signintech/pdft/minigopdf"
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
		//fontIndex, err := strconv.Atoi(strings.TrimSpace(strings.Replace(prop.key, "F", "", -1)))
		prefix, fontIndex, err := c.getFontIndex(prop.key)
		if err != nil {
			return err
		}
		objID, _, err := prop.asDictionary()
		if err != nil {
			return err
		}
		crFont.fontPrefix = prefix
		crFont.fontIndex = fontIndex
		crFont.fontObjID = objID
		*c = append(*c, crFont)
	}
	return nil
}

func (c crawlResultFonts) getFontIndex(key string) (string, int, error) {
	//fmt.Printf("key=%s\n", key)
	ss := []rune(key)
	ssSize := len(ss)
	var nn []string
	for i := ssSize - 1; i >= 0; i-- {
		//fmt.Printf("%s\n", string(rs[i]))
		s := string(ss[i])
		_, err := strconv.Atoi(s)
		if err != nil {
			break
		}
		nn = append(nn, s)
	}

	var buff bytes.Buffer
	nnSize := len(nn)
	for i := nnSize - 1; i >= 0; i-- {
		buff.WriteString(nn[i])
	}

	prefix := string(ss[0 : ssSize-nnSize])

	n, err := strconv.Atoi(buff.String())
	if err != nil {
		return "", 0, err
	}
	//fmt.Printf("%s  %d   %d\n", prefix, n, nnSize)
	return prefix, n, nil
}

func (c *crawlResultFonts) String() string {
	var buff bytes.Buffer
	for _, f := range *c {
		prefix := strings.TrimSpace(f.fontPrefix)
		if prefix == "" {
			prefix = gopdf.FontPrefixDefault
		}
		buff.WriteString(fmt.Sprintf("/%s%d %d 0 R\n", prefix, f.fontIndex, f.fontObjID))
	}
	fmt.Printf("%s\n", buff.String())
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
	fontPrefix string
	fontIndex  int
	fontObjID  int
}
