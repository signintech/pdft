package pdft

import "testing"

func TestPDFParse(t *testing.T) {
	//var err error
	/*err = parsePDF("pdf_from_docx.pdf")
	if err != nil {
	t.Error(err)
	return
	}*/
	/*
		err = parsePDF("test/pdf/pdf_from_docx_with_f.pdf")
		if err != nil {
		t.Error(err)
		return
		}
	*/
	/*
		err = analyzePDF("pdf_from_docx.pdf")
		if err != nil {
			t.Error(err)
			return
		}*/
	/*err = analyzePDF("pdf_from_gopdf.pdf")
	if err != nil {
		t.Error(err)
		return
	}*/
	/*
		err = analyzePDF("pdf_from_docx_with_f.pdf")
		if err != nil {
			t.Error(err)
			return
		}*/

}

/*
func _TestReadProperties(t *testing.T) {
	err := testReadProperties("pdf_from_gopdf.pdf", 3)
	if err != nil {
		t.Error(err)
		return
	}

}

func testReadProperties(file string, id int) error {

	f, err := os.Open("test/pdf/" + file)
	if err != nil {
		return err
	}
	defer f.Close()

	var pdf PDFData
	err = PDFParse(f, &pdf)
	if err != nil {
		return err
	}

	obj := pdf.getObjByID(id)
	prop, err := obj.readProperties()
	if err != nil {
		return err
	}
	_ = prop
	//fmt.Printf("prop=%#v", prop)

	return nil

}

func parsePDF(file string) error {

	f, err := os.Open("test/pdf/" + file)
	if err != nil {
		return err
	}
	defer f.Close()

	var pdf PDFData
	err = PDFParse(f, &pdf)
	if err != nil {
		return err
	}

	//fmt.Printf("%#v\n", pdf)

	return nil
}

func analyzePDF(file string) error {
	f, err := os.Open("test/pdf/" + file)
	if err != nil {
		return err
	}
	defer f.Close()

	var pdf PDFData
	err = PDFParse(f, &pdf)
	if err != nil {
		return err
	}
	var analyze AnalyzeData
	err = AnalyzePDF(&pdf, &analyze)
	if err != nil {
		return err
	}
	//fmt.Printf("%#v\n", analyze)

	return nil
}*/
