package pdft

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type crawlResult struct {
	items []crawlResultItem
}

func (c *crawlResult) String() string {
	var buff bytes.Buffer
	loopCrawlResultString(&buff, c)
	return buff.String()
}

//ErrCrawlResultValOfNotFound CrawlResult Val Of Not Found
var ErrCrawlResultValOfNotFound = errors.New("CrawlResult Val Of Not Found")

func (c *crawlResult) valOf(key string) (string, error) {
	return loopValOf(c, key)
}

func (c *crawlResult) setValOf(key string, val string) error {
	return loopSetValOf(c, key, val)
}

func (c *crawlResult) add(key string, val string) error {
	var item crawlResultItem
	item.key = key
	item.valStr = val
	item.typeOfVal = typeOfValStr
	//fmt.Printf("key===%s\n", key)
	c.items = append(c.items, item)
	return nil
}

func loopSetValOf(cr *crawlResult, key string, val string) error {
	for i := range cr.items {
		item := &cr.items[i]
		if item.typeOfVal == typeOfValStr {
			if item.key == key {
				item.setValStr(val)
				return nil
			}
		} else if item.typeOfVal == typeOfValCr {
			v, _ := item.getValCr()
			return loopSetValOf(v, key, val)
		}
	}
	return ErrCrawlResultValOfNotFound
}

func loopValOf(cr *crawlResult, key string) (string, error) {
	//fmt.Printf("\n#####%s\n", key)
	for _, item := range cr.items {
		if item.typeOfVal == typeOfValStr {
			//fmt.Printf("key=%s\n", item.key)
			if item.key == key {
				v, err := item.getValStr()
				if err != nil {
					return "", err
				}
				return v, nil
			}
		} else if item.typeOfVal == typeOfValCr {
			v, err := item.getValCr()
			if err != nil {
				return "", err
			}
			return loopValOf(v, key)
		}
	}
	return "", ErrCrawlResultValOfNotFound
}

func loopCrawlResultString(buff *bytes.Buffer, cr *crawlResult) {
	for _, item := range cr.items {
		buff.WriteString(fmt.Sprintf("/%s ", item.key))
		if item.typeOfVal == typeOfValStr {
			if strings.TrimSpace(item.valStr) != "" {
				buff.WriteString(fmt.Sprintf("%s\n", item.valStr))
			}
		} else if item.typeOfVal == typeOfValCr {
			buff.WriteString("\n\t<<")
			loopCrawlResultString(buff, item.valCr)
			buff.WriteString(">>\n")
		}
	}
}

var typeOfValStr = "str"
var typeOfValCr = "cr" //crawlResult

type crawlResultItem struct {
	key       string
	typeOfVal string
	valCr     *crawlResult
	valStr    string
}

//ErrWrongTypeOfVal wrong type of val
var ErrWrongTypeOfVal = errors.New("wrong type of val")

func (c *crawlResultItem) getValCr() (*crawlResult, error) {
	if c.typeOfVal != typeOfValCr {
		return nil, ErrWrongTypeOfVal
	}
	return c.valCr, nil
}

func (c *crawlResultItem) getValStr() (string, error) {
	if c.typeOfVal != typeOfValStr {
		return "", ErrWrongTypeOfVal
	}
	return c.valStr, nil
}

func (c *crawlResultItem) setValStr(v string) {
	c.typeOfVal = typeOfValStr
	c.valStr = v
}

func (c *crawlResultItem) setValCr(v *crawlResult) {
	c.typeOfVal = typeOfValCr
	c.valCr = v
}
