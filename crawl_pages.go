package pdft

import (
	"strings"
)

type crawlPages struct{}

func (c *crawlPages) getPageCrawl(pdf *PDFData, objID int, p ...string) (*crawl, error) {
	var cw crawl
	pagePath := append([]string{"Pages"}, p...)
	cw.set(pdf, objID, pagePath...)
	cw.run()
	checkedQueue := []int{}
	for k := range cw.results {
		checkedQueue = append(checkedQueue, k)
	}
	for len(checkedQueue) > 0 {
		key := checkedQueue[0]
		if s := cw.results[key].String(); strings.Contains(s, "/Pages") && strings.Contains(s, "/Parent") {
			var subCw crawl
			subCw.set(pdf, key, p...)
			subCw.run()
			for k, v := range subCw.results {
				cw.results[k] = v
				if _, ok := cw.results[k]; !ok {
					checkedQueue = append(checkedQueue, k)
				}
			}
		}
		checkedQueue = checkedQueue[1:]
	}
	return &cw, nil
}

func (c *crawlPages) getPageObjIDs(cw *crawl) ([]int, error) {
	results := []int{}
	for k, v := range cw.results {
		if s := v.String(); !strings.Contains(s, "/Pages") && strings.Contains(s, "/Page") && strings.Contains(s, "/Parent") {
			results = append(results, k)
		}
	}
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j] < results[i] {
				t := results[i]
				results[i] = results[j]
				results[j] = t
			}
		}
	}
	return results, nil
}
