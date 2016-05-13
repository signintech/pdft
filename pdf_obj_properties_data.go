package pdft

//PDFObjPropertiesData array of PDFObjPropertyData
type PDFObjPropertiesData []PDFObjPropertyData

func (p *PDFObjPropertiesData) append(prop PDFObjPropertyData) {
	(*p) = append((*p), prop)
}

func (p *PDFObjPropertiesData) size() int {
	return len(*p)
}

func (p *PDFObjPropertiesData) getPropByKey(key string) *PDFObjPropertyData {
	for _, prop := range *p {
		if prop.key == key {
			return &prop
		}
	}
	return nil
}

func (p *PDFObjPropertiesData) at(i int) *PDFObjPropertyData {
	return &(*p)[i]
}
