package pdft

import "regexp"

var regexpXref = regexp.MustCompile("xref")
var regexpXrefLine = regexp.MustCompile("[0-9]{10}[\\t ]+[0-9]{5}[\\t ][f,n]")
var regexpStartObj = regexp.MustCompile("[0-9]+[\\t ]0[\\t ]obj")
var regexpEndObj = regexp.MustCompile("endobj")
var regexpTrailer = regexp.MustCompile("trailer")
var regexpStartxref = regexp.MustCompile("startxref")
