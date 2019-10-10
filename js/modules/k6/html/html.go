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
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/pkg/errors"
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
	return Selection***REMOVED***rt: common.GetRuntime(ctx), sel: doc.Selection***REMOVED***, nil
***REMOVED***

type Selection struct ***REMOVED***
	rt  *goja.Runtime
	sel *goquery.Selection
	URL string `json:"url"`
***REMOVED***

func (s Selection) emptySelection() Selection ***REMOVED***
	// Goquery has no direct way to return an empty selection apart from asking for an out of bounds item.
	return s.Eq(s.Size())
***REMOVED***

func (s Selection) buildMatcher(v goja.Value, gojaFn goja.Callable) func(int, *goquery.Selection) bool ***REMOVED***
	return func(idx int, sel *goquery.Selection) bool ***REMOVED***
		fnRes, fnErr := gojaFn(v, s.rt.ToValue(idx), s.rt.ToValue(sel))
		if fnErr != nil ***REMOVED***
			common.Throw(s.rt, fnErr)
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
		return Selection***REMOVED***s.rt, selFilter(v.sel), s.URL***REMOVED***

	case string:
		return Selection***REMOVED***s.rt, strFilter(v), s.URL***REMOVED***

	case Element:
		return Selection***REMOVED***s.rt, nodeFilter(v.node), s.URL***REMOVED***

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
		return Selection***REMOVED***s.rt, filtered(def[0]), s.URL***REMOVED***
	***REMOVED***

	return Selection***REMOVED***s.rt, unfiltered(), s.URL***REMOVED***
***REMOVED***

func (s Selection) adjacentUntil(until func(string) *goquery.Selection,
	untilSelection func(*goquery.Selection) *goquery.Selection,
	filteredUntil func(string, string) *goquery.Selection,
	filteredUntilSelection func(string, *goquery.Selection) *goquery.Selection,
	def ...goja.Value) Selection ***REMOVED***

	switch len(def) ***REMOVED***
	case 0:
		return Selection***REMOVED***s.rt, until(""), s.URL***REMOVED***
	case 1:
		switch selector := def[0].Export().(type) ***REMOVED***
		case string:
			return Selection***REMOVED***s.rt, until(selector), s.URL***REMOVED***

		case Selection:
			return Selection***REMOVED***s.rt, untilSelection(selector.sel), s.URL***REMOVED***

		case nil:
			return Selection***REMOVED***s.rt, until(""), s.URL***REMOVED***
		***REMOVED***
	case 2:
		filter := def[1].String()

		switch selector := def[0].Export().(type) ***REMOVED***
		case string:
			return Selection***REMOVED***s.rt, filteredUntil(filter, selector), s.URL***REMOVED***

		case Selection:
			return Selection***REMOVED***s.rt, filteredUntilSelection(filter, selector.sel), s.URL***REMOVED***

		case nil:
			return Selection***REMOVED***s.rt, filteredUntil(filter, ""), s.URL***REMOVED***
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

	return Selection***REMOVED***s.rt, s.sel.NotFunction(s.buildMatcher(v, gojaFn)), s.URL***REMOVED***
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
	return Selection***REMOVED***s.rt, s.sel.End(), s.URL***REMOVED***
***REMOVED***

func (s Selection) Eq(idx int) Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Eq(idx), s.URL***REMOVED***
***REMOVED***

func (s Selection) First() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.First(), s.URL***REMOVED***
***REMOVED***

func (s Selection) Last() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Last(), s.URL***REMOVED***
***REMOVED***

func (s Selection) Contents() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Contents(), s.URL***REMOVED***
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

// nolint: goconst
func (s Selection) Val() goja.Value ***REMOVED***
	switch goquery.NodeName(s.sel) ***REMOVED***
	case InputTagName:
		val, exists := s.sel.Attr("value")
		if !exists ***REMOVED***
			inputType, _ := s.sel.Attr("type")
			if inputType == "radio" || inputType == "checkbox" ***REMOVED***
				val = "on"
			***REMOVED*** else ***REMOVED***
				val = ""
			***REMOVED***
		***REMOVED***
		return s.rt.ToValue(val)

	case ButtonTagName:
		val, exists := s.sel.Attr("value")
		if !exists ***REMOVED***
			val = ""
		***REMOVED***
		return s.rt.ToValue(val)

	case TextAreaTagName:
		return s.Html()

	case OptionTagName:
		return s.rt.ToValue(valueOrHTML(s.sel))

	case SelectTagName:
		selected := s.sel.First().Find("option[selected]")
		if _, exists := s.sel.Attr("multiple"); exists ***REMOVED***
			return s.rt.ToValue(selected.Map(func(idx int, opt *goquery.Selection) string ***REMOVED*** return valueOrHTML(opt) ***REMOVED***))
		***REMOVED***

		return s.rt.ToValue(valueOrHTML(selected))

	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (s Selection) Children(def ...string) Selection ***REMOVED***
	if len(def) == 0 ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Children(), s.URL***REMOVED***
	***REMOVED***

	return Selection***REMOVED***s.rt, s.sel.ChildrenFiltered(def[0]), s.URL***REMOVED***
***REMOVED***

func (s Selection) Each(v goja.Value) Selection ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if !isFn ***REMOVED***
		common.Throw(s.rt, errors.New("Argument to each() must be a function."))
	***REMOVED***

	fn := func(idx int, sel *goquery.Selection) ***REMOVED***
		if _, err := gojaFn(v, s.rt.ToValue(idx), selToElement(Selection***REMOVED***s.rt, s.sel.Eq(idx), s.URL***REMOVED***)); err != nil ***REMOVED***
			common.Throw(s.rt, errors.Wrap(err, "Function passed to each() failed."))
		***REMOVED***
	***REMOVED***

	return Selection***REMOVED***s.rt, s.sel.Each(fn), s.URL***REMOVED***
***REMOVED***

func (s Selection) Filter(v goja.Value) Selection ***REMOVED***
	switch val := v.Export().(type) ***REMOVED***
	case string:
		return Selection***REMOVED***s.rt, s.sel.Filter(val), s.URL***REMOVED***

	case Selection:
		return Selection***REMOVED***s.rt, s.sel.FilterSelection(val.sel), s.URL***REMOVED***
	***REMOVED***

	gojaFn, isFn := goja.AssertFunction(v)
	if !isFn ***REMOVED***
		common.Throw(s.rt, errors.New("Argument to filter() must be a function, a selector or a selection"))
	***REMOVED***

	return Selection***REMOVED***s.rt, s.sel.FilterFunction(s.buildMatcher(v, gojaFn)), s.URL***REMOVED***
***REMOVED***

func (s Selection) Is(v goja.Value) bool ***REMOVED***
	switch val := v.Export().(type) ***REMOVED***
	case string:
		return s.sel.Is(val)

	case Selection:
		return s.sel.IsSelection(val.sel)

	default:
		gojaFn, isFn := goja.AssertFunction(v)
		if !isFn ***REMOVED***
			common.Throw(s.rt, errors.New("Argument to is() must be a function, a selector or a selection"))
		***REMOVED***

		return s.sel.IsFunction(s.buildMatcher(v, gojaFn))
	***REMOVED***
***REMOVED***

func (s Selection) Map(v goja.Value) []string ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if !isFn ***REMOVED***
		common.Throw(s.rt, errors.New("Argument to map() must be a function"))
	***REMOVED***

	fn := func(idx int, sel *Selection) string ***REMOVED***
		if fnRes, fnErr := gojaFn(v, s.rt.ToValue(idx), s.rt.ToValue(sel)); fnErr == nil ***REMOVED***
			return fnRes.String()
		***REMOVED***
		return ""
	***REMOVED***

	// Mimics goquery.Selector.Map function.
	// We can not use s.sel.Map directly, otherwise, goja will see the function body elements as type
	// *gohtml.Selection instead of our *Selection wrapper, so goja runtime will call wrong methods.
	// See issue #1195
	result := make([]string, len(s.sel.Nodes))
	for i, n := range s.sel.Nodes ***REMOVED***
		result[i] = fn(i, &Selection***REMOVED***rt: s.rt, sel: &goquery.Selection***REMOVED***Nodes: []*gohtml.Node***REMOVED***n***REMOVED******REMOVED******REMOVED***)
	***REMOVED***
	return result
***REMOVED***

func (s Selection) Slice(start int, def ...int) Selection ***REMOVED***
	if len(def) > 0 ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Slice(start, def[0]), s.URL***REMOVED***
	***REMOVED***

	return Selection***REMOVED***s.rt, s.sel.Slice(start, s.sel.Length()), s.URL***REMOVED***
***REMOVED***

func (s Selection) Get(def ...int) goja.Value ***REMOVED***
	switch ***REMOVED***
	case len(def) == 0:
		var items []goja.Value
		for i := 0; i < len(s.sel.Nodes); i++ ***REMOVED***
			items = append(items, selToElement(s.Eq(i)))
		***REMOVED***
		return s.rt.ToValue(items)

	case def[0] < s.sel.Length() && def[0] > -s.sel.Length():
		return selToElement(s.Eq(def[0]))

	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (s Selection) ToArray() []Selection ***REMOVED***
	items := make([]Selection, len(s.sel.Nodes))
	for i := range s.sel.Nodes ***REMOVED***
		items[i] = Selection***REMOVED***s.rt, s.sel.Eq(i), s.URL***REMOVED***
	***REMOVED***
	return items
***REMOVED***

func (s Selection) Index(def ...goja.Value) int ***REMOVED***
	if len(def) == 0 ***REMOVED***
		return s.sel.Index()
	***REMOVED***

	switch v := def[0].Export().(type) ***REMOVED***
	case Selection:
		return s.sel.IndexOfSelection(v.sel)

	case string:
		return s.sel.IndexSelector(v)

	case Element:
		return s.sel.IndexOfNode(v.node)

	default:
		return -1
	***REMOVED***
***REMOVED***

// When 0 arguments: Read all data from attributes beginning with "data-".
// When 1 argument: Append argument to "data-" then find for a matching attribute
func (s Selection) Data(def ...string) goja.Value ***REMOVED***
	if s.sel.Length() == 0 || len(s.sel.Nodes[0].Attr) == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	if len(def) > 0 ***REMOVED***
		val, exists := s.sel.Attr("data-" + propertyToAttr(def[0]))
		if exists ***REMOVED***
			return s.rt.ToValue(convertDataAttrVal(val))
		***REMOVED*** else ***REMOVED***
			return goja.Undefined()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		data := make(map[string]interface***REMOVED******REMOVED***)
		for _, attr := range s.sel.Nodes[0].Attr ***REMOVED***
			if strings.HasPrefix(attr.Key, "data-") && len(attr.Key) > 5 ***REMOVED***
				data[attrToProperty(attr.Key[5:])] = convertDataAttrVal(attr.Val)
			***REMOVED***
		***REMOVED***
		return s.rt.ToValue(data)
	***REMOVED***
***REMOVED***
