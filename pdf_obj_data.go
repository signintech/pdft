package pdft

import (
	"bytes"
	"crypto/rc4"
	"errors"

	gopdf "github.com/signintech/pdft/minigopdf"
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

	stream := []byte("stream")
	endstream := []byte("endstream")
	start := bytes.Index(p.data, stream)
	if start != -1 {
		end := bytes.LastIndex(p.data, endstream)

		var head, body bytes.Buffer
		head.Write(p.data[0:start])
		body.Write(p.data[start+len(stream) : end])
		tmp := bytes.Trim(body.Bytes(), "\r\n")
		tmp = bytes.Trim(tmp, "\n")
		bodyRc4, err := rc4Cip(protection.Objectkey(p.objID), tmp)
		if err != nil {
			return err
		}

		var data bytes.Buffer
		data.Write(head.Bytes())
		data.Write(stream)
		data.WriteString("\n")
		data.Write(bodyRc4)
		data.WriteString("\n")
		data.Write(endstream)
		p.data = data.Bytes()
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
