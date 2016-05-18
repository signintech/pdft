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
	err := rd.Open(pdfTmpl)
	if err != nil {
		return nil, err
	}
	return &rd, nil
}

//Render pdf render
type Render struct {
	pdft.PDFt
	finfoMap map[string]FieldInfo
}

//Text write text to pdf
func (r *Render) Text(key string, text string) error {
	if finfo, ok := r.finfoMap[key]; ok {

		err := r.SetFont(finfo.Font, "", finfo.Size)
		if err != nil {
			return err
		}

		err = r.Insert(text, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H, finfo.Align)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrNotFoundKey
}

//ShowDesignView for debug
func (r *Render) ShowDesignView() {
	r.ShowCellBorder(true)
	for key, finfo := range r.finfoMap {
		r.SetFont(finfo.Font, "", finfo.Size)
		r.Insert("#"+key, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H, finfo.Align)
	}
}
