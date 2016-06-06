package pdft

import "fmt"

type crawlResultXObjects []crawlResultXObject

func (c *crawlResultXObjects) parse(propVal *[]byte) error {
	var props PDFObjPropertiesData
	err := readProperties(propVal, &props)
	if err != nil {
		return err
	}

	for _, prop := range props {
		fmt.Printf("\t %#v\n", prop)
	}

	return nil
}

type crawlResultXObject struct {
	xObjChar  string
	xObjIndex int
	xObjObjID int
}
