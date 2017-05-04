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
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"

	gohtml "golang.org/x/net/html"
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

func (s Selection) emptySelection() Selection ***REMOVED***
	// Ask for out of bounds item for an empty selection.
	return s.Eq(s.Size())
***REMOVED***

func (s Selection) buildMatcher(v goja.Value, gojaFn goja.Callable) func(int, *goquery.Selection) bool ***REMOVED***
	return func(idx int, sel *goquery.Selection) bool ***REMOVED***
		fnRes, fnErr := gojaFn(v, s.rt.ToValue(idx), s.rt.ToValue(sel))

		if fnErr != nil ***REMOVED***
			panic(fnErr)
		***REMOVED***

		return fnRes.ToBoolean()
	***REMOVED***
***REMOVED***

func (s Selection) varargFnCall(arg interface***REMOVED******REMOVED***,
	strFilter func(string) *goquery.Selection,
	selFilter func(*goquery.Selection) *goquery.Selection,
	nodeFilter func(...*gohtml.Node) *goquery.Selection) Selection ***REMOVED***

	switch v := arg.(type) ***REMOVED***
	case Selection:
		return Selection***REMOVED***s.rt, selFilter(v.sel)***REMOVED***

	case string:
		return Selection***REMOVED***s.rt, strFilter(v)***REMOVED***

	case Element:
		return Selection***REMOVED***s.rt, nodeFilter(v.node)***REMOVED***

	case goja.Value:
		return s.varargFnCall(v.Export(), strFilter, selFilter, nodeFilter)

	default:
		errmsg := fmt.Sprintf("Invalid argument: Cannot use a %T as a selector", arg)
		panic(s.rt.NewGoError(errors.New(errmsg)))
	***REMOVED***
***REMOVED***

func (s Selection) adjacent(unfiltered func() *goquery.Selection,
	filtered func(string) *goquery.Selection,
	def ...string) Selection ***REMOVED***
	if len(def) > 0 ***REMOVED***
		return Selection***REMOVED***s.rt, filtered(def[0])***REMOVED***
	***REMOVED***

	return Selection***REMOVED***s.rt, unfiltered()***REMOVED***
***REMOVED***

func (s Selection) adjacentUntil(until func(string) *goquery.Selection,
	untilSelection func(*goquery.Selection) *goquery.Selection,
	filteredUntil func(string, string) *goquery.Selection,
	filteredUntilSelection func(string, *goquery.Selection) *goquery.Selection,
	def ...goja.Value) Selection ***REMOVED***

	switch len(def) ***REMOVED***
	case 0:
		return Selection***REMOVED***s.rt, until("")***REMOVED***
	case 1:
		switch selector := def[0].Export().(type) ***REMOVED***
		case string:
			return Selection***REMOVED***s.rt, until(selector)***REMOVED***

		case Selection:
			return Selection***REMOVED***s.rt, untilSelection(selector.sel)***REMOVED***

		case nil:
			return Selection***REMOVED***s.rt, until("")***REMOVED***
		***REMOVED***
	case 2:
		filter := def[1].String()

		switch selector := def[0].Export().(type) ***REMOVED***
		case string:
			return Selection***REMOVED***s.rt, filteredUntil(filter, selector)***REMOVED***

		case Selection:
			return Selection***REMOVED***s.rt, filteredUntilSelection(filter, selector.sel)***REMOVED***

		case nil:
			return Selection***REMOVED***s.rt, filteredUntil(filter, "")***REMOVED***
		***REMOVED***
	***REMOVED***

	errmsg := fmt.Sprintf("Invalid argument: Cannot use a %T as a selector", def[0].Export())
	panic(s.rt.NewGoError(errors.New(errmsg)))
***REMOVED***

func (s Selection) Add(arg interface***REMOVED******REMOVED***) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Add, s.sel.AddSelection, s.sel.AddNodes)
***REMOVED***

func (s Selection) Find(arg interface***REMOVED******REMOVED***) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Find, s.sel.FindSelection, s.sel.FindNodes)
***REMOVED***

func (s Selection) Closest(arg interface***REMOVED******REMOVED***) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Closest, s.sel.ClosestSelection, s.sel.ClosestNodes)
***REMOVED***

func (s Selection) Has(arg interface***REMOVED******REMOVED***) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Has, s.sel.HasSelection, s.sel.HasNodes)
***REMOVED***

func (s Selection) Not(v goja.Value) Selection ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if !isFn ***REMOVED***
		return s.varargFnCall(v, s.sel.Not, s.sel.NotSelection, s.sel.NotNodes)
	***REMOVED***

	return Selection***REMOVED***s.rt, s.sel.NotFunction(s.buildMatcher(v, gojaFn))***REMOVED***
***REMOVED***

func (s Selection) Next(def ...string) Selection ***REMOVED***
	return s.adjacent(s.sel.Next, s.sel.NextFiltered, def...)
***REMOVED***

func (s Selection) NextAll(def ...string) Selection ***REMOVED***
	return s.adjacent(s.sel.NextAll, s.sel.NextAllFiltered, def...)
***REMOVED***

func (s Selection) Prev(def ...string) Selection ***REMOVED***
	return s.adjacent(s.sel.Prev, s.sel.PrevFiltered, def...)
***REMOVED***

func (s Selection) PrevAll(def ...string) Selection ***REMOVED***
	return s.adjacent(s.sel.PrevAll, s.sel.PrevAllFiltered, def...)
***REMOVED***

func (s Selection) Parent(def ...string) Selection ***REMOVED***
	return s.adjacent(s.sel.Parent, s.sel.ParentFiltered, def...)
***REMOVED***

func (s Selection) Parents(def ...string) Selection ***REMOVED***
	return s.adjacent(s.sel.Parents, s.sel.ParentsFiltered, def...)
***REMOVED***

func (s Selection) Siblings(def ...string) Selection ***REMOVED***
	return s.adjacent(s.sel.Siblings, s.sel.SiblingsFiltered, def...)
***REMOVED***

// prevUntil, nextUntil and parentsUntil support two arguments with mutable type.
// 1st argument is the selector. Either a selector string, a Selection object, or nil
// 2nd argument is the filter. Either a selector string or nil/undefined
func (s Selection) PrevUntil(def ...goja.Value) Selection ***REMOVED***
	return s.adjacentUntil(
		s.sel.PrevUntil,
		s.sel.PrevUntilSelection,
		s.sel.PrevFilteredUntil,
		s.sel.PrevFilteredUntilSelection,
		def...,
	)
***REMOVED***

func (s Selection) NextUntil(def ...goja.Value) Selection ***REMOVED***
	return s.adjacentUntil(
		s.sel.NextUntil,
		s.sel.NextUntilSelection,
		s.sel.NextFilteredUntil,
		s.sel.NextFilteredUntilSelection,
		def...,
	)
***REMOVED***

func (s Selection) ParentsUntil(def ...goja.Value) Selection ***REMOVED***
	return s.adjacentUntil(
		s.sel.ParentsUntil,
		s.sel.ParentsUntilSelection,
		s.sel.ParentsFilteredUntil,
		s.sel.ParentsFilteredUntilSelection,
		def...,
	)
***REMOVED***

func (s Selection) Size() int ***REMOVED***
	return s.sel.Length()
***REMOVED***

func (s Selection) End() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.End()***REMOVED***
***REMOVED***

func (s Selection) Eq(idx int) Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Eq(idx)***REMOVED***
***REMOVED***

func (s Selection) First() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.First()***REMOVED***
***REMOVED***

func (s Selection) Last() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Last()***REMOVED***
***REMOVED***

func (s Selection) Contents() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Contents()***REMOVED***
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

func (s Selection) Html() goja.Value ***REMOVED***
	val, err := s.sel.Html()
	if err != nil ***REMOVED***
		return goja.Undefined()
	***REMOVED***
	return s.rt.ToValue(val)
***REMOVED***

func optionVal(s *goquery.Selection) string ***REMOVED***
	val, exists := s.Attr("value")

	if exists ***REMOVED***
		return val
	***REMOVED***

	val, err := s.Html()

	if err != nil ***REMOVED***
		return ""
	***REMOVED***

	return val
***REMOVED***

func(s Selection) Val() goja.Value ***REMOVED***
	switch goquery.NodeName(s.sel) ***REMOVED***
		case "input":
			return s.Attr("value")

		case "textarea":
			return s.Html()

		case "button":
			return s.Attr("value")

		case "select":
			selected := s.sel.First().Find("option[selected]")

			_, exists := s.sel.Attr("multiple")

			if exists ***REMOVED***
				return s.rt.ToValue(selected.Map(func(idx int, opt *goquery.Selection) string ***REMOVED*** return optionVal(opt) ***REMOVED***))
			***REMOVED*** else ***REMOVED***
				return s.rt.ToValue(optionVal(selected))
			***REMOVED***

		case "":
			return goja.Undefined()
		default:
			return goja.Undefined()
	***REMOVED***
***REMOVED***

func (s Selection) Closest(name string) Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Closest(name)***REMOVED***
***REMOVED***
