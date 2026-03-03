package htmlutil

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func ExtractInputs(body []byte, formID string) (map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)

	form := doc.Find("form#" + formID)
	if form.Length() == 0 {
		form = doc.Selection
	}

	form.Find("input").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name != "" {
			val, _ := s.Attr("value")
			values[name] = val
		}
	})

	return values, nil
}
