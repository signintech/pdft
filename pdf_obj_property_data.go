package pdft

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const dictionary = "dictionary"
const object = "object"
const array = "array"
const number = "number"

//PDFObjPropertyData property of pdf obj
type PDFObjPropertyData struct {
	key    string
	rawVal string
}

func (p *PDFObjPropertyData) setRaw(rawVal string) {
	p.rawVal = rawVal
}

func (p *PDFObjPropertyData) setAsDictionary(value, revision int) {
	p.rawVal = fmt.Sprintf("%d %d R", value, revision)
}

func (p *PDFObjPropertyData) setAsDictionaryArr(values, revisions []int) {
	if revisions == nil {
		revisions = make([]int, len(values))
	}
	var data bytes.Buffer
	data.WriteString("[")
	for i := range values {
		data.WriteString(fmt.Sprintf("%d %d R ", values[i], revisions[i]))
	}
	data.WriteString("]")
	p.rawVal = data.String()
}

func (p *PDFObjPropertyData) asDictionary() (int, int, error) {
	return readObjIDFromDictionary(p.rawVal)
}

func (p *PDFObjPropertyData) asDictionaryArr() ([]int, []int, error) {
	return readObjIDFromDictionaryArr(p.rawVal)
}

func (p *PDFObjPropertyData) valType() string {
	return propertyType(p.rawVal)
}

func propertyType(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) > len("<<") && raw[0:len("<<")] == "<<" {
		return object
	} else if len(raw) > len("[") && raw[0:len("[")] == "[" {
		return array
	} else if _, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil {
		return number
	}
	//fmt.Printf("raw=%s\n", raw)
	return dictionary
}

func readProperty(rawObj *[]byte, key string) (*PDFObjPropertyData, error) {
	var outProps PDFObjPropertiesData
	err := readProperties(rawObj, &outProps)
	if err != nil {
		return nil, err
	}
	return outProps.getPropByKey(key), nil
}

func readProperties(rawObj *[]byte, outProps *PDFObjPropertiesData) error {

	tmp0 := *rawObj
	index := bytes.Index(*rawObj, extractStreamBytes)
	if index != -1 { //เป็น stream
		tmp0 = tmp0[0:index]
	}
	startObjInx := strings.Index(string(tmp0), "<<")
	endObjInx := strings.LastIndex(string(tmp0), ">>")
	if startObjInx > endObjInx || startObjInx == -1 || endObjInx == -1 {
		return errors.New("bad obj properties")
	}

	//fmt.Printf("\n\n%s\n\n%d\n", string(*rawObj), endObjInx)

	var regexpSlash = regexp.MustCompile("[\\n\\t ]+\\/")
	var regexpOpenB = regexp.MustCompile("[\\n\\t ]+\\[")
	var regexpCloseB = regexp.MustCompile("[\\n\\t ]+\\]")
	var regexpOpen = regexp.MustCompile("[\\n\\t ]+\\<\\<")
	var regexpClose = regexp.MustCompile("[\\n\\t ]+\\>\\>")
	var regexpLine = regexp.MustCompile("[\\n\\t ]+")

	tmp := strings.TrimSpace(string((*rawObj)[startObjInx+len("<<") : endObjInx]))
	tmp = regexpLine.ReplaceAllString(tmp, " ")
	tmp = regexpSlash.ReplaceAllString(tmp, "/")
	tmp = regexpOpenB.ReplaceAllString(tmp, "[")
	tmp = regexpCloseB.ReplaceAllString(tmp, "]")
	tmp = regexpOpen.ReplaceAllString(tmp, "<<")
	tmp = regexpClose.ReplaceAllString(tmp, ">>")

	var pp parseProps

	pp.set(tmp, outProps)
	return nil
}

type parseProps struct {
	str        string
	max        int
	propsIndex int
	props      *PDFObjPropertiesData
}

func (p *parseProps) set(str string, props *PDFObjPropertiesData) {
	p.str = str
	p.max = len(str)
	p.propsIndex = -1
	p.props = props
	p.loop(0, "")
}

func (p *parseProps) loop(i int, status string) (int, string) {
	count01 := 0
	count02 := 0
	for i < p.max {
		r := string(p.str[i])
		if status == "" && r == "/" {
			p.propsIndex++
			p.props.append(PDFObjPropertyData{})
			i, status = p.loop(i+1, "key")
		} else if status == "key" {
			if r == " " {
				i, status = p.loop(i+1, "val")
			} else if r == "<" || r == "[" {
				i, status = p.loop(i, "val")
			} else if r == "/" {
				return i - 1, ""
			} else {
				p.props.at(p.propsIndex).key += r
			}
		} else if status == "val" {

			if r == "<" {
				count01++
			} else if r == "[" {
				count02++
			} else if r == ">" {
				count01--
			} else if r == "]" {
				count02--
			}

			if (r == "]" || r == ">") && (count01 == 0 && count02 == 0) {
				p.props.at(p.propsIndex).rawVal += r
				return i, ""
			} else if r == "/" && (count01 == 0 && count02 == 0) {
				return i - 1, ""
			}
			p.props.at(p.propsIndex).rawVal += r

		}
		i++
	}
	return i, status
}

func readObjIDFromDictionaryArr(str string) ([]int, []int, error) {

	str = strings.Replace(str, "[", "", -1)
	str = strings.Replace(str, "]", "", -1)
	str = strings.TrimSpace(str)
	tokens := strings.Split(str, " ")
	var objIDs []int
	var revisions []int

	i := 0
	max := len(tokens)
	for i < max {
		objID, err := strconv.Atoi(strings.TrimSpace(tokens[i]))
		if err != nil {
			return nil, nil, err
		}
		revision, err := strconv.Atoi(strings.TrimSpace(tokens[i+1]))
		if err != nil {
			return nil, nil, err
		}
		objIDs = append(objIDs, objID)
		revisions = append(revisions, revision)
		i += 3
	}

	return objIDs, revisions, nil
}

//ErrorObjectIDNotFound Object ID not found
var ErrorObjectIDNotFound = errors.New("Object ID not found")

func readObjIDFromDictionary(str string) (int, int, error) {

	str = strings.TrimSpace(str)
	if str == "" {
		return 0, 0, ErrorObjectIDNotFound
	}

	tokens := strings.Split(str, " ")
	if len(tokens) != 3 {
		return 0, 0, ErrorObjectIDNotFound
	}

	id, err := strconv.Atoi(strings.TrimSpace(tokens[0]))
	if err != nil {
		return 0, 0, err
	}
	return id, 0, nil
}
