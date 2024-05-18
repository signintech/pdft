PDFT
====

PDFT is a GO library for creating PDF documents using existing PDFs as template.
This library depend on [gopdf](https://github.com/signintech/gopdf). 

Tested with PDF template files created from Libre office, Google Docs, Microsoft Word.

 
### SAMPLE
```go
var pt pdft.PDFt
err := pt.Open(pdfsource)
if err != nil {
	panic("Couldn't open pdf.")
}


err = pt.AddFont("arial", "./arial.ttf")
if err != nil {
    t.Error(err)
    return
}

err = pt.SetFont("arial", "", 14)
if err != nil {
    panic(err) 
}

//insert text to pdf
err = pt.Insert("Hi", 1, 10, 10, 100, 100, gopdf.Center|gopdf.Bottom, nil)
if err != nil {
    panic(err) 
}

// insert text to pdf with color
err = pt.Insert("Hi", 1, 10, 10, 100, 100, gopdf.Center|gopdf.Bottom, &FontColor{R: 255, G: 255, B: 255})
if err != nil {
    panic(err) 
}

// measure text width
var textWidth float64
textWidth, err = pt.MeasureTextWidth("Hi")

// read image file
pic, err := ioutil.ReadFile(picPath)
if err != nil {
	panic("Couldn't read pic.")
}

 
// insert image to pdf
err = pt.InsertImg(pic, 1, 182.0, 165.0, 172.0, 49.0)
if err != nil {
	panic("Couldn't insert image")
}

// insert image to pdf with cache, avoiding redundant data when inserting same images many times
err = pt.InsertImgWithCache(pic, 1, 182.0, 165.0, 172.0, 49.0)
if err != nil {
	panic("Couldn't insert image")
}

// Duplicate first page to last page
err = pt.DuplicatePageAfter(1, -1)
if err != nil {
	panic("Couldn't duplicate first page")
}

// Remove second page
err = pt.RemovePage(2)
if err != nil {
	panic("Couldn't remove targetPage")
}

err = pt.Save(target)
if err != nil {
	panic("Couldn't save pdf.")
}
```


