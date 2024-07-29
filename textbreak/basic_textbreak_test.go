package textbreak

import "testing"

func TestBasicTextBreak(t *testing.T) {
	text := "Hello World 123"
	tbk := BasicTextbreak{}
	token, err := tbk.BreakTextToToken(text)
	if err != nil {
		t.Error(err)
		return
	}

	if len(token) != 5 {
		t.Errorf("expect 5 but got %d", len(token))
		return
	}

}

func TestBasicTextBreak2(t *testing.T) {
	text := "HelloWorld123"
	tbk := BasicTextbreak{}
	token, err := tbk.BreakTextToToken(text)
	if err != nil {
		t.Error(err)
		return
	}

	if len(token) != 1 {
		t.Errorf("expect 1 but got %d", len(token))
		return
	}
	if token[0] != text {
		t.Errorf("expect %s but got %s", text, token[0])
		return
	}
}
