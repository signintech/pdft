package pdft

import (
	"bytes"
	"io"

	gopdf "github.com/kelvinsantos/pdft/minigopdf"
)

// PDFFontData font data
type PDFFontData struct {
	fontname  string
	fontindex int
	//gopdf obj
	font       gopdf.SubsetFontObj
	cid        gopdf.CIDFontObj
	unicodeMap gopdf.UnicodeMap
	fontDesc   gopdf.SubfontDescriptorObj
	dictionary gopdf.PdfDictionaryObj
	//IDs
	startID                                               int
	fontID, cidID, unicodeMapID, fontDescID, dictionaryID int
	//stream
	fontStream       *bytes.Buffer
	cidStream        *bytes.Buffer
	unicodeMapStream *bytes.Buffer
	fontDescStream   *bytes.Buffer
	dictionaryStream *bytes.Buffer
}

func (p *PDFFontData) init() {
	p.font.CharacterToGlyphIndex = gopdf.NewMapOfCharacterToGlyphIndex() //make(map[rune]uint)
}

func (p *PDFFontData) setFontName(name string) {
	p.font.SetFamily(name)
	p.fontname = name
}

func (p *PDFFontData) setFontIndex(fontindex int) {
	p.fontindex = fontindex
}

func (p *PDFFontData) fontIndex() int {
	return p.fontindex
}

func (p *PDFFontData) fontName() string {
	return p.fontname
}

func (p *PDFFontData) setTTFPath(path string) error {
	return p.font.SetTTFByPath(path)
}

func (p *PDFFontData) setTTFReader(reader io.Reader) error {
	return p.font.SetTTFByReader(reader)
}

func (p *PDFFontData) addChars(text string) error {
	return p.font.AddChars(text)
}

func (p *PDFFontData) build() (int, error) {

	p.fontID = p.startID + 1
	p.cidID = p.startID + 2
	p.unicodeMapID = p.startID + 3
	p.fontDescID = p.startID + 4
	p.dictionaryID = p.startID + 5
	newMaxID := p.dictionaryID

	//font
	p.font.SetIndexObjCIDFont(p.cidID - 1)
	p.font.SetIndexObjUnicodeMap(p.unicodeMapID - 1)

	//cid
	p.cid.SetPtrToSubsetFontObj(&p.font)
	p.cid.SetIndexObjSubfontDescriptor(p.fontDescID - 1)

	//unicode
	p.unicodeMap.SetPtrToSubsetFontObj(&p.font)

	//font descriptor
	p.fontDesc.SetPtrToSubsetFontObj(&p.font)
	p.fontDesc.SetIndexObjPdfDictionary(p.dictionaryID - 1)

	//dictionary
	p.dictionary.SetPtrToSubsetFontObj(&p.font)

	var err error
	err = p.font.Build(p.fontID)
	if err != nil {
		return 0, err
	}

	err = p.cid.Build(p.cidID)
	if err != nil {
		return 0, err
	}

	err = p.unicodeMap.Build(p.unicodeMapID)
	if err != nil {
		return 0, err
	}

	err = p.fontDesc.Build(p.fontDescID)
	if err != nil {
		return 0, err
	}

	err = p.dictionary.Build(p.dictionaryID)
	if err != nil {
		return 0, err
	}

	p.fontStream = p.font.GetObjBuff()
	p.cidStream = p.cid.GetObjBuff()
	p.unicodeMapStream = p.unicodeMap.GetObjBuff()
	p.fontDescStream = p.fontDesc.GetObjBuff()
	p.dictionaryStream = p.dictionary.GetObjBuff()

	return newMaxID, nil
}

func (p *PDFFontData) setStartID(id int) {
	p.startID = id
}
