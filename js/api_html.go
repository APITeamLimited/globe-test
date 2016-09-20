package js

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func (a JSAPI) HTMLParse(src string) *goquery.Selection ***REMOVED***
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(src))
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***
	return doc.Selection
***REMOVED***
