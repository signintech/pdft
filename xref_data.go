package pdft

import (
	"errors"
	"strconv"
	"strings"
)

// XrefData nnnnnnnnnn ggggg x eol
type XrefData struct {
	N10Digit int    //nnnnnnnnnn
	G5Digit  int    //ggggg
	Keyword  string //x
}

func (x *XrefData) parse(xrefline string) error {
	xrefline = strings.TrimSpace(xrefline)
	tokens := strings.Split(xrefline, " ")
	if len(tokens) < 3 {
		return errors.New("bad xref format")
	}
	var err error
	x.N10Digit, err = strconv.Atoi(strings.TrimSpace(tokens[0]))
	if err != nil {
		return err
	}

	x.G5Digit, err = strconv.Atoi(strings.TrimSpace(tokens[1]))
	if err != nil {
		return err
	}

	x.Keyword = strings.TrimSpace(tokens[2])
	if x.Keyword != "n" && x.Keyword != "f" {
		return errors.New("unkonw xref keyword")
	}

	return nil
}
