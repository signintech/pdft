package render

import (
	"errors"

	"github.com/signintech/pdft"
)

//ErrNotFoundKey key not found
var ErrNotFoundKey = errors.New("not found key")

//NewRender create new render
func NewRender(pdfTmpl string, finfos FieldInfos) (*Render, error) {
	var rd Render
	rd.finfoMap = finfos.toMap()
	err := rd.open(pdfTmpl)
	if err != nil {
		return nil, err
	}
	return &rd, nil
}

//Render pdf render
type Render struct {
	pt       pdft.PDFt
	finfoMap map[string]FieldInfo
}

func (r *Render) open(filepath string) error {
	err := r.pt.Open(filepath)
	if err != nil {
		return err
	}
	return err
}

//WriteText write text to pdf
func (r *Render) WriteText(key string, text string) error {

	return ErrNotFoundKey
}
