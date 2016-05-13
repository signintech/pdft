package pdft

import "bytes"

//Contenter create content stream
type Contenter interface {
	page() int
	toSteram() (*bytes.Buffer, error)
}
