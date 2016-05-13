package pdft

//AnalyzePDF analyze pdf file
/*func AnalyzePDF(pdf *PDFData, outAnalyze *AnalyzeData) error {

	var cw crawl
	cw.set(pdf, pdf.trailer.rootObjID, "Pages", "Kids", "Resources", "Font")
	cw.run()
	cw.lastRouteObjID()
	return nil
}*/

/*
//AnalyzePDF analyze pdf file
func AnalyzePDF(pdf *PDFData, outAnalyze *AnalyzeData) error {

	rootObjID := pdf.trailer.rootObjID
	pagesProp, err := getPropFromObj(pdf, rootObjID, "Pages")
	if err != nil {
		return err
	}

	pagesID, _, err := pagesProp.asDictionary()
	if err != nil {
		return err
	}

	kidsProp, err := getPropFromObj(pdf, pagesID, "Kids")
	if err != nil {
		return err
	}

	kidIDs, _, err := kidsProp.asDictionaryArr()
	if err != nil {
		return err
	}

	for i, kidID := range kidIDs {

		resourcesProp, err := getPropFromObj(pdf, kidID, "Resources")
		if err != nil {
			return err
		}
		contentsProp, err := getPropFromObj(pdf, kidID, "Contents")
		if err != nil {
			return err
		}

		resourcesID, _, err := resourcesProp.asDictionary()
		if err != nil {
			return err
		}

		var resourcesObjRaw []byte
		if resourcesProp.valType() == subobject {
			resourcesObjRaw = []byte(resourcesProp.rawVal)
		} else if resourcesProp.valType() == dictionary {
			if err != nil {
				return err
			}
			resourcesObjRaw = pdf.getObjByID(resourcesID).data
		}

		contentsID, _, err := contentsProp.asDictionary()
		if err != nil {
			return err
		}

		//create pageAnalyze
		pageAnalyze := PageAnalyzeData{
			pageNum:    i + 1,
			contentsID: contentsID,
		}
		pageAnalyze.resources.resourcesID = resourcesID
		err = analyzeResources(pdf, &resourcesObjRaw, &pageAnalyze.resources)
		if err != nil {
			return err
		}
		outAnalyze.pages = append(outAnalyze.pages, pageAnalyze)
	}

	return nil
}
*/
/*
func analyzeResources(pdf *PDFData, resourcesObjRaw *[]byte, outresourcesAn *ResourcesAnalyzeData) error {

	fontprop, err := readProperty(resourcesObjRaw, "Font")
	if err != nil {
		return err
	}

	var rawSubObj []byte
	if fontprop.valType() == dictionary {
		var fontObjID int
		fontObjID, _, err = fontprop.asDictionary()
		if err != nil {
			return err
		}
		obj := pdf.getObjByID(fontObjID)
		rawSubObj = obj.data
	} else if fontprop.valType() == subobject {
		rawSubObj = []byte(fontprop.rawVal)
	} else {
		return errors.New("cannot read resources")
	}

	var props PDFObjPropertiesData
	err = readProperties(&rawSubObj, &props)
	if err != nil {
		return err
	}

	for _, prop := range props {
		var fontData FontAnalyzeData
		fontIndex, err := strconv.Atoi(strings.Replace(prop.key, "F", "", -1))
		if err != nil {
			return err
		}
		objID, _, err := prop.asDictionary()
		if err != nil {
			return err
		}
		fontData.fontIndex = fontIndex
		fontData.objID = objID
		outresourcesAn.fonts = append(outresourcesAn.fonts, fontData)
	}

	return nil
}

func getPropFromObj(pdf *PDFData, objID int, propKey string) (*PDFObjPropertyData, error) {

	obj := pdf.getObjByID(objID)
	props, err := obj.readProperties()
	if err != nil {
		return nil, err
	}
	prop := props.getPropByKey(propKey)
	if prop == nil {
		return nil, errors.New("not found prop " + propKey)
	}

	return prop, nil
}*/
