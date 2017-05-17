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
	"encoding/json"
	"errors"
	"strconv"
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
	return s.Eq(s.Size()) //ask for out of bounds item to empty selection
***REMOVED***

func (s Selection) varargFnCall(arg goja.Value,
	strFilter func(string) *goquery.Selection,
	selFilter func(*goquery.Selection) *goquery.Selection,
	nodeFilter func(...*gohtml.Node) *goquery.Selection) Selection ***REMOVED***

	val := arg.Export()
	switch val.(type) ***REMOVED***
	case Selection:
		return Selection***REMOVED***s.rt, selFilter(val.(Selection).sel)***REMOVED***

	case string:
		return Selection***REMOVED***s.rt, strFilter(val.(string))***REMOVED***

	case map[string]interface***REMOVED******REMOVED***:
		if elem, ok := valToElement(arg); ok ***REMOVED***
			return Selection***REMOVED***s.rt, nodeFilter(elem.node)***REMOVED***
		***REMOVED*** else ***REMOVED***
			return Selection***REMOVED***s.rt, s.emptySelection().sel***REMOVED***
		***REMOVED***

	default:
		return Selection***REMOVED***s.rt, s.emptySelection().sel***REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Add(arg goja.Value) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Add, s.sel.AddSelection, s.sel.AddNodes)
***REMOVED***

func (s Selection) Find(arg goja.Value) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Find, s.sel.FindSelection, s.sel.FindNodes)
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
	if val, exists := s.Attr("value"); exists ***REMOVED***
		return val
	***REMOVED***

	if val, err := s.Html(); err == nil ***REMOVED***
		return val
	***REMOVED***

	return ""
***REMOVED***

func (s Selection) Val() goja.Value ***REMOVED***
	switch goquery.NodeName(s.sel) ***REMOVED***
	case "input":
		return s.Attr("value")

	case "textarea":
		return s.Html()

	case "button":
		return s.Attr("value")

	case "option":
		return s.rt.ToValue(optionVal(s.sel))

	case "select":
		selected := s.sel.First().Find("option[selected]")

		if _, exists := s.sel.Attr("multiple"); exists ***REMOVED***
			return s.rt.ToValue(selected.Map(func(idx int, opt *goquery.Selection) string ***REMOVED*** return optionVal(opt) ***REMOVED***))
		***REMOVED*** else ***REMOVED***
			return s.rt.ToValue(optionVal(selected))
		***REMOVED***

	default:
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (s Selection) Closest(arg goja.Value) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Closest, s.sel.ClosestSelection, s.sel.ClosestNodes)
***REMOVED***

func (s Selection) Children(def ...string) Selection ***REMOVED***
	if len(def) == 0 ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Children()***REMOVED***
	***REMOVED*** else ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.ChildrenFiltered(def[0])***REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Contents() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.Contents()***REMOVED***
***REMOVED***

func (s Selection) Each(v goja.Value) Selection ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if isFn ***REMOVED***
		fn := func(idx int, sel *goquery.Selection) ***REMOVED***
			gojaFn(v, s.rt.ToValue(idx), selToElement(s))
		***REMOVED***

		return Selection***REMOVED***s.rt, s.sel.Each(fn)***REMOVED***
	***REMOVED*** else ***REMOVED***
		panic(s.rt.NewGoError(errors.New("Argument to each() must be a function")))
		return s
	***REMOVED***
***REMOVED***

func (s Selection) End() Selection ***REMOVED***
	return Selection***REMOVED***s.rt, s.sel.End()***REMOVED***
***REMOVED***

func (s Selection) buildMatcher(v goja.Value, gojaFn goja.Callable) func(int, *goquery.Selection) bool ***REMOVED***
	return func(idx int, sel *goquery.Selection) bool ***REMOVED***
		fnRes, fnErr := gojaFn(v, s.rt.ToValue(idx), s.rt.ToValue(sel))
		return fnErr == nil && fnRes.ToBoolean()
	***REMOVED***
***REMOVED***

func (s Selection) Filter(v goja.Value) Selection ***REMOVED***
	if gojaFn, isFn := goja.AssertFunction(v); isFn ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.FilterFunction(s.buildMatcher(v, gojaFn))***REMOVED***
	***REMOVED*** else if cmp, isSel := v.Export().(Selection); isSel ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.FilterSelection(cmp.sel)***REMOVED***
	***REMOVED*** else if str, isStr := v.Export().(string); isStr ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Filter(str)***REMOVED***
	***REMOVED*** else ***REMOVED***
		panic(s.rt.NewGoError(errors.New("Argument to filter() must be a function, a selector or a query object")))
		return Selection***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Is(v goja.Value) bool ***REMOVED***
	if gojaFn, isFn := goja.AssertFunction(v); isFn ***REMOVED***
		return s.sel.IsFunction(s.buildMatcher(v, gojaFn))
	***REMOVED*** else if cmp, isSel := v.Export().(Selection); isSel ***REMOVED***
		return s.sel.IsSelection(cmp.sel)
	***REMOVED*** else if str, isStr := v.Export().(string); isStr ***REMOVED***
		return s.sel.Is(str)
	***REMOVED*** else ***REMOVED***
		panic(s.rt.NewGoError(errors.New("Argument to is() must be a function, a selector or a query object")))
		return false
	***REMOVED***
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

func (s Selection) Has(arg goja.Value) Selection ***REMOVED***
	return s.varargFnCall(arg, s.sel.Has, s.sel.HasSelection, s.sel.HasNodes)
***REMOVED***

func (s Selection) Map(v goja.Value) (result []string) ***REMOVED***
	gojaFn, isFn := goja.AssertFunction(v)
	if isFn ***REMOVED***
		fn := func(idx int, sel *goquery.Selection) string ***REMOVED***
			if fnRes, fnErr := gojaFn(v, s.rt.ToValue(idx), s.rt.ToValue(sel)); fnErr == nil ***REMOVED***
				return fnRes.String()
			***REMOVED*** else ***REMOVED***
				return ""
			***REMOVED***
		***REMOVED***
		return s.sel.Map(fn)
	***REMOVED*** else ***REMOVED***
		panic(s.rt.NewGoError(errors.New("Argument to map() must be a function")))
		return nil
	***REMOVED***
***REMOVED***

func (s Selection) Not(v goja.Value) Selection ***REMOVED***
	if gojaFn, isFn := goja.AssertFunction(v); isFn ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.NotFunction(s.buildMatcher(v, gojaFn))***REMOVED***
	***REMOVED*** else ***REMOVED***
		return s.varargFnCall(v, s.sel.Not, s.sel.NotSelection, s.sel.NotNodes)
	***REMOVED***
***REMOVED***

func (s Selection) adjacent(unfiltered func() *goquery.Selection,
	filtered func(string) *goquery.Selection,
	def ...string) Selection ***REMOVED***
	if len(def) == 0 ***REMOVED***
		return Selection***REMOVED***s.rt, unfiltered()***REMOVED***
	***REMOVED*** else ***REMOVED***
		return Selection***REMOVED***s.rt, filtered(def[0])***REMOVED***
	***REMOVED***
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

func (s Selection) adjacentUntil(until func(string) *goquery.Selection,
	untilSelection func(*goquery.Selection) *goquery.Selection,
	filteredUntil func(string, string) *goquery.Selection,
	filteredUntilSelection func(string, *goquery.Selection) *goquery.Selection,
	def ...goja.Value) Selection ***REMOVED***
	// empty selector to nextuntil and prevuntil acts like revAll and nextAll
	// relies on goquery.compileMatcher returning a matcher which fails all matches when the selector being compiled is invalid
	if len(def) == 0 ***REMOVED***
		return Selection***REMOVED***s.rt, until("")***REMOVED***
	***REMOVED***

	selector := def[0].Export()
	if len(def) == 1 ***REMOVED***
		switch selector.(type) ***REMOVED***
		case string:
			return Selection***REMOVED***s.rt, until(selector.(string))***REMOVED***

		case Selection:
			return Selection***REMOVED***s.rt, untilSelection(selector.(Selection).sel)***REMOVED***

		case nil:
			return Selection***REMOVED***s.rt, until("")***REMOVED***

		default:
			panic(s.rt.NewGoError(errors.New("Invalid argument. The selector must be a string or query object")))
			return Selection***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		filter := def[1].String()
		switch selector.(type) ***REMOVED***
		case string:
			return Selection***REMOVED***s.rt, filteredUntil(filter, selector.(string))***REMOVED***

		case Selection:
			return Selection***REMOVED***s.rt, filteredUntilSelection(filter, selector.(Selection).sel)***REMOVED***

		case nil:
			return Selection***REMOVED***s.rt, filteredUntil(filter, "")***REMOVED***

		default:
			panic(s.rt.NewGoError(errors.New("Invalid argument. The selector must be a string or query object")))
			return Selection***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// prevUntil, nextUntil and parentsUntil support two args
// 1st arg is either a selector string, goquery selection object, or nil
// 2nd arg is filter selector string or nil/undefined
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

func (s Selection) Slice(start int, def ...int) Selection ***REMOVED***
	if len(def) > 0 ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Slice(start, def[0])***REMOVED***
	***REMOVED*** else ***REMOVED***
		return Selection***REMOVED***s.rt, s.sel.Slice(start, s.sel.Length())***REMOVED***
	***REMOVED***
***REMOVED***

func (s Selection) Get(def ...int) goja.Value ***REMOVED***
	if len(def) == 0 ***REMOVED***
		var items []goja.Value

		for i := 0; i < len(s.sel.Nodes); i++ ***REMOVED***
			items = append(items, selToElement(s.Eq(i)))
		***REMOVED***

		return s.rt.ToValue(items)
	***REMOVED*** else if def[0] < s.sel.Length() && def[0] > -s.sel.Length() ***REMOVED***
		return selToElement(s.Eq(def[0]))
	***REMOVED*** else ***REMOVED***
		return goja.Undefined()
	***REMOVED***
***REMOVED***

func (s Selection) ToArray() (items []Selection) ***REMOVED***
	for i := range s.sel.Nodes ***REMOVED***
		items = append(items, Selection***REMOVED***s.rt, s.sel.Eq(i)***REMOVED***)
	***REMOVED***
	return
***REMOVED***

func (s Selection) Size() int ***REMOVED***
	return s.sel.Length()
***REMOVED***

func (s Selection) Index(def ...goja.Value) int ***REMOVED***
	if len(def) == 0 ***REMOVED***
		return s.sel.Index()
	***REMOVED***

	v := def[0].Export()
	switch v.(type) ***REMOVED***
	case Selection:
		return s.sel.IndexOfSelection(v.(Selection).sel)

	case string:
		return s.sel.IndexSelector(v.(string))

	case map[string]interface***REMOVED******REMOVED***:
		if elem, ok := valToElement(def[0]); ok ***REMOVED***
			return s.sel.IndexOfNode(elem.node)
		***REMOVED*** else ***REMOVED***
			return -1
		***REMOVED***
	default:
		return -1
	***REMOVED***
***REMOVED***

// end result of the following is two strings.Replacer objects
// Replacer("-a", "A", "-b", "B"..., "-z", "Z") and Replacer("A", "-a",...)
//to translate to "data-attr-name" to "attrName" and back
const (
	lowAlpha  = "abcdefghijklmnopqrstuvwxyz"
	highAlpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func makeReplacerArray(prefixFrom, from, prefixTo, to string) (vals []string) ***REMOVED***
	for idx, _ := range from ***REMOVED***
		vals = append(vals, prefixFrom+string(from[idx]), prefixTo+string(to[idx]))
	***REMOVED***
	return
***REMOVED***

func makeNameReplacer(prefixFrom, from, prefixTo, to string) *strings.Replacer ***REMOVED***
	return strings.NewReplacer(makeReplacerArray(prefixFrom, from, prefixTo, to)...)
***REMOVED***

var attrToDataName = makeNameReplacer("-", lowAlpha, "", highAlpha)
var dataToAttrName = makeNameReplacer("", highAlpha, "-", lowAlpha)

func toAttrName(dataName string) string ***REMOVED***
	return dataToAttrName.Replace(dataName)
***REMOVED***

func toDataName(attrName string) string ***REMOVED***
	return attrToDataName.Replace(attrName)
***REMOVED***

// return numeric value when the representation is unchanged by conversion to float and back
// other numeric values (ie "101.00" "1E02") are left as strings
func toNumeric(val string) (float64, bool) ***REMOVED***
	if fltVal, err := strconv.ParseFloat(val, 64); err != nil ***REMOVED***
		return 0, false
	***REMOVED*** else if repr := strconv.FormatFloat(fltVal, 'f', -1, 64); repr == val ***REMOVED***
		return fltVal, true
	***REMOVED*** else ***REMOVED***
		return 0, false
	***REMOVED***
***REMOVED***

func convert(val string) interface***REMOVED******REMOVED*** ***REMOVED***
	if len(val) == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED*** else if val[0] == '***REMOVED***' || val[0] == '[' ***REMOVED***
		var subdata interface***REMOVED******REMOVED***

		err := json.Unmarshal([]byte(val), &subdata)
		if err == nil ***REMOVED***
			return subdata
		***REMOVED*** else ***REMOVED***
			return val
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		switch val ***REMOVED***
		case "true":
			return true

		case "false":
			return false

		case "null":
			return goja.Undefined()

		case "undefined":
			return goja.Undefined()

		default:
			if fltVal, isOk := toNumeric(val); isOk ***REMOVED***
				return fltVal
			***REMOVED*** else ***REMOVED***
				return val
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

//when 0 args, read all data from attributes beginning with "data-".
//when 1 arg, read requested data attr
func (s Selection) Data(def ...string) goja.Value ***REMOVED***
	if s.sel.Length() == 0 || len(s.sel.Nodes[0].Attr) == 0 ***REMOVED***
		return goja.Undefined()
	***REMOVED***

	if len(def) > 0 ***REMOVED***
		val, exists := s.sel.Attr("data-" + toAttrName(def[0]))
		if exists ***REMOVED***
			return s.rt.ToValue(convert(val))
		***REMOVED*** else ***REMOVED***
			return goja.Undefined()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		data := make(map[string]interface***REMOVED******REMOVED***)
		for _, attr := range s.sel.Nodes[0].Attr ***REMOVED***
			if strings.HasPrefix(attr.Key, "data-") && len(attr.Key) > 5 ***REMOVED***
				data[toDataName(attr.Key[5:])] = convert(attr.Val)
			***REMOVED***
		***REMOVED***
		return s.rt.ToValue(data)
	***REMOVED***
***REMOVED***
