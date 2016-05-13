package pdft

type crawl struct {
	//setup
	pdf   *PDFData
	objID int
	paths []string

	//result
	index int

	//onCrawl     funcOnCrawl
	results map[int]*crawlResult
}

func (c *crawl) resultByObjID(id int) *crawlResult {
	if r, ok := c.results[id]; ok {
		return r
	}
	c.results[id] = new(crawlResult)
	return c.results[id]
}

func (c *crawl) set(pdf *PDFData, objID int, p ...string) {

	c.pdf = pdf
	c.objID = objID
	c.paths = p
	//init
	c.results = make(map[int]*crawlResult)
}

func (c *crawl) run() error {
	objdata := c.pdf.getObjByID(c.objID)
	err := c.next(&objdata.data, 0, c.objID, c.resultByObjID(c.objID))
	return err
}

func (c *crawl) next(data *[]byte, i int, id int, cr *crawlResult) error {

	lenPath := len(c.paths)
	if lenPath <= i {

		return nil
	}

	var err error
	var props PDFObjPropertiesData
	err = readProperties(data, &props)
	if err != nil {
		return err
	}

	//var cr crawlResult
	for _, prop := range props {
		var item crawlResultItem
		item.key = prop.key
		if prop.key != c.paths[i] {
			item.setValStr(prop.rawVal)
		} else {
			propType := prop.valType()
			if propType == dictionary {
				var objID int
				item.setValStr(prop.rawVal)
				objID, _, err = prop.asDictionary()
				if err != nil {
					return err
				}
				objdata := c.pdf.getObjByID(objID)
				err = c.next(&objdata.data, i+1, objID, c.resultByObjID(objID))
				if err != nil {
					return err
				}
			} else if propType == array {
				var objIDs []int
				item.setValStr(prop.rawVal)
				objIDs, _, err = prop.asDictionaryArr()
				if err != nil {
					return err
				}
				for _, objID := range objIDs {
					objdata := c.pdf.getObjByID(objID)
					c.next(&objdata.data, i+1, objID, c.resultByObjID(objID))
				}
			} else if propType == object {
				if lenPath <= i+1 {
					item.setValStr(prop.rawVal)
				} else {
					var subCr crawlResult
					item.setValCr(&subCr)
					tmp := []byte(prop.rawVal)
					err = c.next(&tmp, i+1, -1, &subCr)
					if err != nil {
						return err
					}
				}
			}
		}
		cr.items = append(cr.items, item)
	}

	return nil
}
