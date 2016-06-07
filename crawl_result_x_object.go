package pdft

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

type crawlResultXObjects []crawlResultXObject

var regSplit = regexp.MustCompile("([a-z A-Z]*)([0-9]*)")

func (c *crawlResultXObjects) parse(propVal *[]byte) error {
	var props PDFObjPropertiesData
	err := readProperties(propVal, &props)
	if err != nil {
		return err
	}

	for _, prop := range props {
		//fmt.Printf("\t %#v\n", prop)
		tokens := regSplit.FindStringSubmatch(prop.key)
		if len(tokens) < 3 {
			continue
		}

		var xObj crawlResultXObject
		xObj.xObjChar = tokens[1]
		xObjIndex, err := strconv.Atoi(tokens[2])
		if err != nil {
			return err
		}
		xObj.xObjIndex = xObjIndex

		xObjObjID, _, err := readObjIDFromDictionary(prop.rawVal)
		if err != nil {
			return err
		}

		xObj.xObjObjID = xObjObjID
		*c = append(*c, xObj)
	}

	return nil
}

func (c *crawlResultXObjects) String() string {
	var buff bytes.Buffer
	for _, xObj := range *c {
		buff.WriteString(fmt.Sprintf("/%s%d %d 0 R", xObj.xObjChar, xObj.xObjIndex, xObj.xObjObjID))
	}
	return buff.String()
}

type crawlResultXObject struct {
	xObjChar  string
	xObjIndex int
	xObjObjID int
}
