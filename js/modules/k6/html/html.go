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

package html

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
)

type HTML struct***REMOVED******REMOVED***

func New() *HTML ***REMOVED***
	return &HTML***REMOVED******REMOVED***
***REMOVED***

func (HTML) ParseHTML(ctx context.Context, src string) (Selection, error) ***REMOVED***
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(src))
	if err != nil ***REMOVED***
		return Selection***REMOVED******REMOVED***, err
	***REMOVED***
	return Selection***REMOVED***common.GetRuntime(ctx), doc.Selection***REMOVED***, nil
***REMOVED***

type Selection struct ***REMOVED***
	rt  *goja.Runtime
	sel *goquery.Selection
***REMOVED***

func (s Selection) Add(arg goja.Value) Selection ***REMOVED***
	switch val := arg.Export().(type) ***REMOVED***
	case Selection:
		return Selection***REMOVED***s.rt, s.sel.AddSelection(val.sel)***REMOVED***
	default:
		return Selection***REMOVED***s.rt, s.sel.Add(arg.String())***REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Find(sel string) Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Find(sel)***REMOVED***
***REMOVED***

func (s Selection) Text() string ***REMOVED***
	return s.sel.Text()
***REMOVED***

func (s Selection) Attr(name string, def ...goja.Value) goja.Value ***REMOVED***
	val, exists := s.sel.Attr(name)
	if !exists ***REMOVED***
		if len(def) > 0 ***REMOVED***
			return def[0]
		***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return s.rt.ToValue(val)
***REMOVED***
