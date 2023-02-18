package pdft

import "io"

// PDFParseFont  parse font
func PDFParseFont(path string, name string) (*PDFFontData, error) {

	var fontData PDFFontData
	fontData.init()
	err := pdfParseFont(path, name, &fontData)
	if err != nil {
		return nil, err
	}
	return &fontData, nil
}

func pdfParseFont(path string, name string, outFontData *PDFFontData) error {
	err := outFontData.setTTFPath(path)
	if err != nil {
		return err
	}
	outFontData.setFontName(name)
	return nil
}

func PDFParseFontReader(reader io.Reader, name string) (*PDFFontData, error) {

	var fontData PDFFontData
	fontData.init()
	err := pdfParseFontReader(reader, name, &fontData)
	if err != nil {
		return nil, err
	}
	return &fontData, nil
}

func pdfParseFontReader(reader io.Reader, name string, outFontData *PDFFontData) error {
	err := outFontData.setTTFReader(reader)
	if err != nil {
		return err
	}
	outFontData.setFontName(name)
	return nil
}
