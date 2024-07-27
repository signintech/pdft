package render

import (
	"fmt"
	"log"

	"github.com/signintech/pdft"
	"github.com/signintech/pdft/textbreak"
)

// ErrNotFoundKey key not found
// var ErrNotFoundKey = errors.New("not found key")
//var errNotFoundKey = "not found key %s"

// NewRender create new render
func NewRender(pdfTmpl string, finfos FieldInfos) (*Render, error) {
	var rd Render
	rd.finfoMap = finfos.toMap()
	err := rd.Open(pdfTmpl)
	if err != nil {
		return nil, err
	}
	return &rd, nil
}

// Render pdf render
type Render struct {
	pdft.PDFt
	finfoMap    map[string]FieldInfo
	textBreaker textbreak.TextBreaker
}

func (r *Render) SetTextBreaker(tb textbreak.TextBreaker) {
	r.textBreaker = tb
}

// Text write text to pdf
func (r *Render) Text(key string, text string) error {

	finfo, ok := r.finfoMap[key]
	if !ok {
		log.Printf("Warr: Not found key %s", key)
		return nil
	}
	err := r.SetFont(finfo.Font, "", finfo.Size)
	if err != nil {
		return fmt.Errorf("setFont %s : %w", finfo.Font, err)
	}
	if finfo.IsWrapText {
		err = r.InsertWrapText(text, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H, finfo.Align, nil, r.textBreaker)
		if err != nil {
			return fmt.Errorf("insertWrapText %s : %w", text, err)
		}
		return nil
	}
	err = r.Insert(text, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H, finfo.Align, nil)
	if err != nil {
		return fmt.Errorf("insert %s : %w", text, err)
	}
	return nil

}

// ImgBase64 image base 64
func (r *Render) ImgBase64(key string, base64 string) error {
	//fmt.Printf("ImgBase64 %s\n\n", base64)
	if finfo, ok := r.finfoMap[key]; ok {
		err := r.InsertImgBase64(base64, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H)
		if err != nil {
			return fmt.Errorf("insertImgBase64 %s : %w", base64, err)
		}
		return nil
	}
	log.Printf("Warr: Not found key %s", key)
	return nil
}

// ShowDesignView for debug
func (r *Render) ShowDesignView() {
	r.ShowCellBorder(true)
	for key, finfo := range r.finfoMap {
		r.SetFont(finfo.Font, "", finfo.Size)
		r.Insert("#"+key, finfo.PageNum, finfo.X, finfo.Y, finfo.W, finfo.H, finfo.Align, nil)
	}
}
