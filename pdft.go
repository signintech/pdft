package pdft

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	gopdf "github.com/signintech/pdft/minigopdf"
)

//ErrAddSameFontName add same font name
var ErrAddSameFontName = errors.New("add same font name")

//ErrFontNameNotFound font name not found
var ErrFontNameNotFound = errors.New("font name not found")

//Left left
const Left = gopdf.Left //001000
//Top top
const Top = gopdf.Top //000100
//Right right
const Right = gopdf.Right //000010
//Bottom bottom
const Bottom = gopdf.Bottom //000001
//Center center
const Center = gopdf.Center //010000
//Middle middle
const Middle = gopdf.Middle //100000

//PDFt inject text to pdf
type PDFt struct {
	pdf           PDFData
	fontDatas     map[string]*PDFFontData
	pdfImgs       []*PDFImageData
	pdfImgsMd5Map map[string]*PDFImageData
	curr          current
	contenters    []Contenter
	pdfProtection *gopdf.PDFProtection
}

type current struct {
	fontName  string
	fontStyle int
	fontSize  int
	lineWidth float64
}

func pageHeight() float64 {
	return 841.89
}

func (i *PDFt) protection() *gopdf.PDFProtection {
	return i.pdfProtection
}

//ShowCellBorder  show cell of border
func (i *PDFt) ShowCellBorder(isShow bool) {
	var clw ContentLineStyle
	if isShow {
		clw.width = 0.1
		clw.lineType = "dotted"
		i.curr.lineWidth = 0.1
	} else {
		clw.width = 0.0
		clw.lineType = ""
		i.curr.lineWidth = 0.0
	}
	i.contenters = append(i.contenters, &clw)
}

//Open open pdf file
func (i *PDFt) Open(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	return i.OpenFrom(f)
}

//OpenFrom open pdf from io.Reader
func (i *PDFt) OpenFrom(r io.Reader) error {
	//init
	i.fontDatas = make(map[string]*PDFFontData)
	i.curr.lineWidth = 1.0
	i.pdfImgsMd5Map = make(map[string]*PDFImageData)
	//open
	err := PDFParse(r, &i.pdf)
	if err != nil {
		return err
	}

	i.ShowCellBorder(false)

	return nil
}

// DuplicatePageAfter ...
func (i *PDFt) DuplicatePageAfter(targetPageNumber, position int) error {
	pageObjIds, err := i.pdf.getPageObjIDs()
	if err != nil {
		return err
	}
	if targetPageNumber > 0 && len(pageObjIds) < targetPageNumber {
		return errors.New("No desired page to copy")
	}

	pageObj := *(i.pdf.getObjByID(pageObjIds[targetPageNumber-1])) //copy object value

	props, err := pageObj.readProperties()
	if err != nil {
		return err
	}
	pageContent := props.getPropByKey("Contents")
	if pageContent == nil {
		return errors.New("No Contents property in this object")
	}
	contentID, _, err := pageContent.asDictionary()
	if err != nil {
		return err
	}
	contentID = i.pdf.putNewObject(*(i.pdf.getObjByID(contentID))) //copy object value
	pageContent.setAsDictionary(contentID, 0)
	pageObj.setProperties(props)

	pageID := i.pdf.putNewObject(pageObj)

	if position < 0 {
		position = len(pageObjIds) + position // like python
	}
	pageObjIds = append(pageObjIds, 0)
	copy(pageObjIds[position+1:], pageObjIds[position:])
	pageObjIds[position+1] = pageID

	return i.setPages(pageObjIds)
}

// RemovePage remove page at targetPageNumber
func (i *PDFt) RemovePage(targetPageNumber int) error {
	pageObjIds, err := i.pdf.getPageObjIDs()
	if err != nil {
		return err
	}
	if targetPageNumber > 0 && len(pageObjIds) < targetPageNumber {
		return errors.New("No desired page to remove")
	}
	copy(pageObjIds[targetPageNumber-1:], pageObjIds[targetPageNumber:])
	pageObjIds = pageObjIds[:len(pageObjIds)-1]

	return i.setPages(pageObjIds)
}

// RemoveOtherPages remove all pages, but not targetPageNumber
func (i *PDFt) RemoveOtherPages(targetPageNumber int) error {
	pageObjIds, err := i.pdf.getPageObjIDs()
	if err != nil {
		return err
	}
	if targetPageNumber > 0 && len(pageObjIds) < targetPageNumber {
		return errors.New("No desired page to keep")
	}

	return i.setPages([]int{pageObjIds[targetPageNumber]})
}

func (i *PDFt) setPages(pageObjIds []int) error {
	props, err := i.pdf.pagesObj.readProperties()
	if err != nil {
		return err
	}
	nPage := len(pageObjIds)
	props.getPropByKey("Count").rawVal = strconv.Itoa(nPage)
	props.getPropByKey("Kids").setAsDictionaryArr(pageObjIds, nil)
	i.pdf.pagesObj.setProperties(props)
	return nil
}

//GetNumberOfPage get number of page
func (i *PDFt) GetNumberOfPage() int {
	pageObjIds, err := i.pdf.getPageObjIDs()
	if err != nil {
		return 0
	}
	return len(pageObjIds)
}

//Insert insert text in to pdf
func (i *PDFt) Insert(text string, pageNum int, x float64, y float64, w float64, h float64, align int) error {
	var ct ContentText
	ct.text = text
	ct.fontName = i.curr.fontName
	ct.fontStyle = i.curr.fontStyle
	ct.fontSize = i.curr.fontSize
	ct.pageNum = pageNum
	ct.x = x
	ct.y = y
	ct.w = w
	ct.h = h
	ct.align = align
	ct.lineWidth = i.curr.lineWidth
	ct.setProtection(i.protection())
	if _, have := i.fontDatas[ct.fontName]; !have {
		return ErrFontNameNotFound
	}
	ct.pdfFontData = i.fontDatas[ct.fontName]
	i.contenters = append(i.contenters, &ct)
	return i.fontDatas[ct.fontName].addChars(text)
}

// MeasureTextWidth measure text width
func (i *PDFt) MeasureTextWidth(text string) (float64, error) {
	i.fontDatas[i.curr.fontName].addChars(text)
	var ct ContentText
	ct.text = text
	ct.fontName = i.curr.fontName
	ct.fontStyle = i.curr.fontStyle
	ct.fontSize = i.curr.fontSize
	ct.lineWidth = i.curr.lineWidth
	if _, have := i.fontDatas[ct.fontName]; !have {
		return 0, ErrFontNameNotFound
	}
	ct.pdfFontData = i.fontDatas[ct.fontName]
	width, err := ct.measureTextWidth()
	return width, err
}

//InsertImgBase64 insert img base 64
func (i *PDFt) InsertImgBase64(base64str string, pageNum int, x float64, y float64, w float64, h float64) error {

	var pdfimg PDFImageData
	imgObj, smask, err := createImgObjFromImgBase64(base64str)
	if err != nil {
		return err
	}

	if smask != nil {
		buff, err := smask.BytesBuffer(0) //ใส่ id ไปมั่วๆทำให้ไม่ support password protect
		if err != nil {
			return err
		}
		var pdfObj PDFObjData
		pdfObj.data = buff.Bytes()
		smaskObjID := i.pdf.putNewObject(pdfObj)
		imgObj.SetSMaskObjID(smaskObjID)
	}

	err = pdfimg.setImgObj(imgObj)
	if err != nil {
		return err
	}
	/*err := pdfimg.setImgBase64(base64str)
	if err != nil {
		return err
	}*/
	i.pdfImgs = append(i.pdfImgs, &pdfimg)
	//fmt.Printf("-->%d\n", len(i.pdfImgs))

	var ct contentImgBase64
	ct.pageNum = pageNum
	ct.x = x
	ct.y = y
	ct.h = h
	ct.w = w
	ct.refPdfimg = &pdfimg //i.pdfImgs[len(i.pdfImgs)-1]
	i.contenters = append(i.contenters, &ct)
	return nil
}

//InsertImg insert img
func (i *PDFt) InsertImg(img []byte, pageNum int, x float64, y float64, w float64, h float64) error {

	var pdfimg PDFImageData
	/*err := pdfimg.setImg(img)
	if err != nil {
		return err
	}*/
	imgObj, smask, err := createImgObjFromBytes(img)
	if err != nil {
		return err
	}
	if smask != nil {
		buff, err := smask.BytesBuffer(0) //ใส่ id ไปมั่วๆทำให้ไม่ support password protect
		if err != nil {
			return err
		}
		var pdfObj PDFObjData
		pdfObj.data = buff.Bytes()
		smaskObjID := i.pdf.putNewObject(pdfObj)
		imgObj.SetSMaskObjID(smaskObjID)
	}

	pdfimg.setImgObj(imgObj)

	i.pdfImgs = append(i.pdfImgs, &pdfimg)
	//fmt.Printf("-->%d\n", len(i.pdfImgs))

	var ct contentImgBase64
	ct.pageNum = pageNum
	ct.x = x
	ct.y = y
	ct.h = h
	ct.w = w
	ct.refPdfimg = &pdfimg //i.pdfImgs[len(i.pdfImgs)-1]
	i.contenters = append(i.contenters, &ct)
	//fmt.Printf("append(i.contenters, &ct) %d\n", len(i.contenters))
	//i.insertContenters(0, &ct)
	return nil
}

//InsertImgWithCache insert img with cache
func (i *PDFt) InsertImgWithCache(img []byte, pageNum int, x float64, y float64, w float64, h float64) error {
	md5Str := fmt.Sprintf("%x", md5.Sum(img))
	var pdfimg *PDFImageData
	if val, ok := i.pdfImgsMd5Map[md5Str]; ok {
		pdfimg = val
	} else {
		pdfimg = &PDFImageData{}
		/*err := pdfimg.setImg(img)
		if err != nil {
			return err
		}*/
		imgObj, smask, err := createImgObjFromBytes(img)
		if err != nil {
			return err
		}
		if smask != nil {
			buff, err := smask.BytesBuffer(0) //ใส่ id ไปมั่วๆทำให้ไม่ support password protect
			if err != nil {
				return err
			}
			var pdfObj PDFObjData
			pdfObj.data = buff.Bytes()
			smaskObjID := i.pdf.putNewObject(pdfObj)
			imgObj.SetSMaskObjID(smaskObjID)
		}
		pdfimg.setImgObj(imgObj)

		i.pdfImgs = append(i.pdfImgs, pdfimg)
		i.pdfImgsMd5Map[md5Str] = pdfimg
	}
	var ct contentImgBase64
	ct.pageNum = pageNum
	ct.x = x
	ct.y = y
	ct.h = h
	ct.w = w
	ct.refPdfimg = pdfimg
	i.contenters = append(i.contenters, &ct)
	return nil
}

/*
func (i *PDFt) insertContenters(index int, src Contenter) {
	i.contenters = append(i.contenters, nil)
	copy(i.contenters[index+1:], i.contenters[index:])
	i.contenters[index] = src
}*/

//AddFont add ttf font
func (i *PDFt) AddFont(name string, ttfpath string) error {

	if _, have := i.fontDatas[name]; have {
		return ErrAddSameFontName
	}

	fontData, err := PDFParseFont(ttfpath, name)
	if err != nil {
		return err
	}

	i.fontDatas[name] = fontData
	return nil
}

//TextriseOverride override text rise
//Text rise, Trise , specifies the distance, in unscaled text space units,
//to move the baseline up or down from its default location.
//Positive values of text rise move the baseline up.
//Adjustments to the baseline are useful for drawing superscripts or subscripts.
//The default location of the baseline can be restored by setting the text rise to 0.
func (i *PDFt) TextriseOverride(name string, fn FuncTextriseOverride) error {
	if _, have := i.fontDatas[name]; !have {
		return ErrFontNameNotFound
	}
	i.fontDatas[name].font.SetFuncTextriseOverride(func(
		leftRune rune,
		rightRune rune,
		fontsize int,
		allText string,
		currTextIndex int,
	) float32 {
		return fn(leftRune, rightRune, fontsize, allText, currTextIndex)
	})
	return nil
}

//KernOverride override kerning
func (i *PDFt) KernOverride(name string, fn FuncKernOverride) error {
	if _, have := i.fontDatas[name]; !have {
		return ErrFontNameNotFound
	}
	i.fontDatas[name].font.SetFuncKernOverride(func(
		leftRune rune,
		rightRune rune,
		leftPair uint,
		rightPair uint,
		pairVal int16,
	) int16 {
		return fn(leftRune, rightRune, leftPair, rightPair, pairVal)
	})
	return nil
}

//SetFont set font
func (i *PDFt) SetFont(name string, style string, size int) error {

	if _, have := i.fontDatas[name]; !have {
		return ErrFontNameNotFound
	}
	i.curr.fontName = name
	i.curr.fontStyle = getConvertedStyle(style)
	i.curr.fontSize = size
	return nil
}

//Save save output pdf
func (i *PDFt) Save(filepath string) error {
	var buff bytes.Buffer
	err := i.SaveTo(&buff)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, buff.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

//SaveTo save pdf to io.Writer
func (i *PDFt) SaveTo(w io.Writer) error {

	newpdf, lastID, err := i.build()
	if err != nil {
		return err
	}

	buff, err := i.toStream(newpdf, lastID)
	if err != nil {
		return err
	}
	_, err = buff.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

func (i *PDFt) build() (*PDFData, int, error) {

	var err error
	nextID := i.pdf.maxID()
	for _, fontData := range i.fontDatas {
		nextID++
		fontData.setStartID(nextID)
		nextID, err = fontData.build()
		if err != nil {
			return nil, 0, err
		}
	}

	newpdf := i.pdf //copy

	err = newpdf.injectFontsToPDF(i.fontDatas)
	if err != nil {
		return nil, 0, err
	}

	//ยัด subsetfont obj กลับไป
	for _, fontData := range i.fontDatas {

		var fontobj, cidObj, unicodeMapObj, fontDescObj, dictionaryObj PDFObjData

		fontobj.objID = fontData.fontID
		fontobj.data = fontData.fontStream.Bytes()

		cidObj.objID = fontData.cidID
		cidObj.data = fontData.cidStream.Bytes()

		unicodeMapObj.objID = fontData.unicodeMapID
		unicodeMapObj.data = fontData.unicodeMapStream.Bytes()

		fontDescObj.objID = fontData.fontDescID
		fontDescObj.data = fontData.fontDescStream.Bytes()

		dictionaryObj.objID = fontData.dictionaryID
		dictionaryObj.data = fontData.dictionaryStream.Bytes()

		newpdf.put(fontobj)
		newpdf.put(cidObj)
		newpdf.put(unicodeMapObj)
		newpdf.put(fontDescObj)
		newpdf.put(dictionaryObj)
	}

	for j, pdfImg := range i.pdfImgs {
		nextID++
		var obj PDFObjData
		obj.objID = nextID
		obj.data = pdfImg.imgObj.GetObjBuff().Bytes()
		i.pdfImgs[j].objID = obj.objID
		//fmt.Printf("---->%d\n", obj.objID)
		newpdf.put(obj)
	}

	err = newpdf.injectImgsToPDF(i.pdfImgs)
	if err != nil {
		return nil, 0, err
	}

	err = newpdf.injectContentToPDF(&i.contenters)
	if err != nil {
		return nil, 0, err
	}

	//set for protection
	if i.protection() != nil {
		max := newpdf.Len()
		x := 0
		for x < max {
			newpdf.objs[x].encrypt(i.protection())
			x++
		}
	}

	return &newpdf, nextID, nil
}

func (i *PDFt) toStream(newpdf *PDFData, lastID int) (*bytes.Buffer, error) {

	//set for protection
	encryptionObjID := -1
	if i.protection() != nil {
		lastID++
		encryptionObjID = lastID
		enObj := i.protection().EncryptionObj()
		err := enObj.Build(lastID)
		if err != nil {
			return nil, err
		}
		buff := enObj.GetObjBuff()
		var enPDFObjData PDFObjData
		enPDFObjData.data = buff.Bytes()
		enPDFObjData.objID = lastID
		newpdf.put(enPDFObjData)
	}

	var buff bytes.Buffer
	buff.WriteString("%PDF-1.7\n\n")
	xrefs := make(map[int]int)
	for _, obj := range newpdf.objs {
		//xrefs = append(xrefs, buff.Len())
		//fmt.Printf("%d\n", obj.objID)
		xrefs[obj.objID] = buff.Len()
		buff.WriteString(fmt.Sprintf("\n%d 0 obj\n", obj.objID))
		buff.WriteString(strings.TrimSpace(string(obj.data)))
		buff.WriteString("\nendobj\n")
	}
	i.xref(xrefs, &buff, lastID, newpdf.trailer.rootObjID, encryptionObjID)

	return &buff, nil
}

type xrefrow struct {
	offset int
	gen    string
	flag   string
}

func (i *PDFt) xref(linelens map[int]int, buff *bytes.Buffer, size int, rootID int, encryptionObjID int) {
	xrefbyteoffset := buff.Len()

	//start xref
	buff.WriteString("\nxref\n")
	buff.WriteString(fmt.Sprintf("0 %d\r\n", size+1))
	var xrefrows []xrefrow
	xrefrows = append(xrefrows, xrefrow{offset: 0, flag: "f", gen: "65535"})
	lastIndexOfF := 0
	j := 1
	//fmt.Printf("size:%d\n", size)
	for j <= size {
		if linelen, ok := linelens[j]; ok {
			xrefrows = append(xrefrows, xrefrow{offset: linelen, flag: "n", gen: "00000"})
		} else {
			xrefrows = append(xrefrows, xrefrow{offset: 0, flag: "f", gen: "65535"})
			offset := len(xrefrows) - 1
			xrefrows[lastIndexOfF].offset = offset
			lastIndexOfF = offset
		}
		j++
	}

	for _, xrefrow := range xrefrows {
		buff.WriteString(i.formatXrefline(xrefrow.offset) + " " + xrefrow.gen + " " + xrefrow.flag + " \n")
	}
	//end xref

	buff.WriteString("trailer\n")
	buff.WriteString("<<\n")
	buff.WriteString(fmt.Sprintf("/Size %d\n", size+1))
	buff.WriteString(fmt.Sprintf("/Root %d 0 R\n", rootID))
	if i.protection() != nil {
		buff.WriteString(fmt.Sprintf("/Encrypt %d 0 R\n", encryptionObjID))
		buff.WriteString("/ID [()()]\n")
	}
	buff.WriteString(">>\n")

	buff.WriteString("startxref\n")
	buff.WriteString(strconv.Itoa(xrefbyteoffset))
	buff.WriteString("\n%%EOF\n")
}

func (i *PDFt) formatXrefline(n int) string {
	str := strconv.Itoa(n)
	for len(str) < 10 {
		str = "0" + str
	}
	return str
}

//SetProtection set pdf protection
func (i *PDFt) SetProtection(
	permissions int,
	userPass []byte,
	ownerPass []byte,
) error {
	var p gopdf.PDFProtection
	err := p.SetProtection(permissions, userPass, ownerPass)
	if err != nil {
		return err
	}
	i.pdfProtection = &p
	return nil
}

//Regular - font style regular
const Regular = 0 //000000
//Italic - font style italic
const Italic = 1 //000001
//Bold - font style bold
const Bold = 2 //000010
//Underline - font style underline
const Underline = 4 //000100

func getConvertedStyle(fontStyle string) (style int) {
	fontStyle = strings.ToUpper(fontStyle)
	if strings.Contains(fontStyle, "B") {
		style = style | Bold
	}
	if strings.Contains(fontStyle, "I") {
		style = style | Italic
	}
	if strings.Contains(fontStyle, "U") {
		style = style | Underline
	}
	return
}
