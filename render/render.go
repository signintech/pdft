package render

import (
	"log"

	"github.com/signintech/pdft"
)

//ErrNotFoundKey key not found
//var ErrNotFoundKey = errors.New("not found key")
var errNotFoundKey = "Not found key %s"

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

	log.Printf("Warr: Not found key %s", key)
	return nil
}

//ImgBase64 image base 64
func (r *Render) ImgBase64(key string, base64 string) error {
	//fmt.Printf("ImgBase64 %s\n\n", base64)
	if finfo, ok := r.finfoMap[key]; ok {
		err := r.InsertImgBase64(base64, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H)
		if err != nil {
			return err
		}
		return nil
	}
	log.Printf("Warr: Not found key %s", key)
	return nil
}

//ShowDesignView for debug
func (r *Render) ShowDesignView() {
	r.ShowCellBorder(true)
	for key, finfo := range r.finfoMap {
		r.SetFont(finfo.Font, "", finfo.Size)
		r.Insert("#"+key, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H, finfo.Align)
	}
}
