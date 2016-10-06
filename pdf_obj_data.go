package pdft

import (
	"bytes"
	"crypto/rc4"
	"errors"
	"fmt"

	"github.com/signintech/gopdf"
)

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

func (p *PDFObjData) encrypt(protection *gopdf.PDFProtection) error {
	//fmt.Printf("=====\n%s\n\n", string(p.data))
	stream := []byte("stream\n")
	endstream := []byte("endstream")
	sIdx := bytes.Index(p.data, stream)
	if sIdx != -1 {
		eIdx := bytes.LastIndex(p.data, endstream)
		head := p.data[0:sIdx]
		body := p.data[sIdx+len(stream) : eIdx]
		body, err := rc4Cip(protection.Objectkey(p.objID), body)
		if err != nil {
			return err
		}
		_ = body
		_ = head

		fmt.Printf("objID=%d\n%s\n\n", p.objID, string(head))

		var data []byte
		data = append(head, stream...)
		data = append(data, body...)
		data = append(data, []byte("\n")...)
		data = append(data, endstream...)
		p.data = data
	}
	return nil
}

func rc4Cip(key []byte, src []byte) ([]byte, error) {
	cip, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dest := make([]byte, len(src))
	cip.XORKeyStream(dest, src)
	return dest, nil
}
