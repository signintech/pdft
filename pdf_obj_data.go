package pdft

import "errors"

//PDFObjData byte of obj
type PDFObjData struct {
	objID int
	data  []byte
}

func (p *PDFObjData) parse(raw *[]byte, stratoffset int) error {

	tmp := (*raw)[stratoffset:]
	endObjIndex := regexpEndObj.FindAllIndex(tmp, 1)
	if len(endObjIndex) <= 0 {
		return errors.New("bad endobj")
	}
	endObjOffsetBefore := endObjIndex[0][0] + stratoffset

	tmp = (*raw)[stratoffset:endObjOffsetBefore]
	startObjIndex := regexpStartObj.FindAllIndex(tmp, 1)
	if len(startObjIndex) <= 0 {
		return errors.New("bad start obj")
	}

	startObjOffsetAfter := startObjIndex[0][1]
	startObjLine := tmp[0:startObjOffsetAfter]
	data := tmp[startObjOffsetAfter:]
	objID, err := objIDFromStartObjLine(string(startObjLine))
	if err != nil {
		return err
	}

	p.objID = objID
	p.data = data
	//fmt.Printf("%s\n", string(data))

	return nil
}

//ReadProperties read all obj Properties
func (p *PDFObjData) readProperties() (*PDFObjPropertiesData, error) {
	var props PDFObjPropertiesData
	err := readProperties(&p.data, &props)
	if err != nil {
		return nil, err
	}
	return &props, nil
}
