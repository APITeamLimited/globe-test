package js

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"strings"
)

func (a JSAPI) HTMLParse(src string) *goquery.Selection ***REMOVED***
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(src))
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***
	return doc.Selection
***REMOVED***

func (a JSAPI) HTMLSelectionAddSelection(vA, vB otto.Value) *goquery.Selection ***REMOVED***
	iA, err := vA.Export()
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***
	selA, ok := iA.(*goquery.Selection)
	if !ok ***REMOVED***
		panic(a.vu.vm.MakeTypeError("HTMLSelectionAddSelection argument A is not a *goquery.Selection"))
	***REMOVED***

	iB, err := vB.Export()
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***
	selB, ok := iB.(*goquery.Selection)
	if !ok ***REMOVED***
		panic(a.vu.vm.MakeTypeError("HTMLSelectionAddSelection argument B is not a *goquery.Selection"))
	***REMOVED***

	return selA.AddSelection(selB)
***REMOVED***
