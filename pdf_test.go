package pdft

import "testing"

func TestPDF(t *testing.T) {
	var ipdf InjPDF
	//err := ipdf.Open("test/pdf/pdf_from_docx_with_f.pdf")
	err := ipdf.Open("test/pdf/pdf_from_gopdf.pdf")
	//err := ipdf.Open("test/pdf/pdf_from_docx.pdf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.AddFont("arial", "test/ttf/arial.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.SetFont("arial", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.InsertText("hi", 1, 100, 100)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.Save("test/out/pdf_from_gopdf.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPDF02(t *testing.T) {
	var ipdf InjPDF
	//err := ipdf.Open("test/pdf/pdf_from_gopdf.pdf")
	err := ipdf.Open("test/pdf/pdf_from_docx.pdf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.AddFont("arial", "test/ttf/arial.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.SetFont("arial", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	err = ipdf.InsertText("hi", 1, 100, 100)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.Save("test/out/pdf_from_docx.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPDF03(t *testing.T) {

	var ipdf InjPDF
	err := ipdf.Open("test/pdf/pdf_from_docx_with_f.pdf")
	//err := ipdf.Open("test/pdf/pdf_from_gopdf.pdf")
	//err := ipdf.Open("test/pdf/pdf_from_docx.pdf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.AddFont("arial", "test/ttf/arial.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.SetFont("arial", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.InsertText("hi", 1, 100, 100)
	if err != nil {
		t.Error(err)
		return
	}

	err = ipdf.Save("test/out/pdf_from_docx_with_f.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}
