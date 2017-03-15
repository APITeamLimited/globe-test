/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package js

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
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
