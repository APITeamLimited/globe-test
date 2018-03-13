// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

import (
	"fmt"
	"reflect"
	"sort"
)

// Slice provides utilities for modifying slices of elements.
// It can be wrapped around any slice of which the element type implements
// interface Elem.
type Slice struct ***REMOVED***
	ptr reflect.Value
	typ reflect.Type
***REMOVED***

// Value returns the reflect.Value of the underlying slice.
func (s *Slice) Value() reflect.Value ***REMOVED***
	return s.ptr.Elem()
***REMOVED***

// MakeSlice wraps a pointer to a slice of Elems.
// It replaces the array pointed to by the slice so that subsequent modifications
// do not alter the data in a CLDR type.
// It panics if an incorrect type is passed.
func MakeSlice(slicePtr interface***REMOVED******REMOVED***) Slice ***REMOVED***
	ptr := reflect.ValueOf(slicePtr)
	if ptr.Kind() != reflect.Ptr ***REMOVED***
		panic(fmt.Sprintf("MakeSlice: argument must be pointer to slice, found %v", ptr.Type()))
	***REMOVED***
	sl := ptr.Elem()
	if sl.Kind() != reflect.Slice ***REMOVED***
		panic(fmt.Sprintf("MakeSlice: argument must point to a slice, found %v", sl.Type()))
	***REMOVED***
	intf := reflect.TypeOf((*Elem)(nil)).Elem()
	if !sl.Type().Elem().Implements(intf) ***REMOVED***
		panic(fmt.Sprintf("MakeSlice: element type of slice (%v) does not implement Elem", sl.Type().Elem()))
	***REMOVED***
	nsl := reflect.MakeSlice(sl.Type(), sl.Len(), sl.Len())
	reflect.Copy(nsl, sl)
	sl.Set(nsl)
	return Slice***REMOVED***
		ptr: ptr,
		typ: sl.Type().Elem().Elem(),
	***REMOVED***
***REMOVED***

func (s Slice) indexForAttr(a string) []int ***REMOVED***
	for i := iter(reflect.Zero(s.typ)); !i.done(); i.next() ***REMOVED***
		if n, _ := xmlName(i.field()); n == a ***REMOVED***
			return i.index
		***REMOVED***
	***REMOVED***
	panic(fmt.Sprintf("MakeSlice: no attribute %q for type %v", a, s.typ))
***REMOVED***

// Filter filters s to only include elements for which fn returns true.
func (s Slice) Filter(fn func(e Elem) bool) ***REMOVED***
	k := 0
	sl := s.Value()
	for i := 0; i < sl.Len(); i++ ***REMOVED***
		vi := sl.Index(i)
		if fn(vi.Interface().(Elem)) ***REMOVED***
			sl.Index(k).Set(vi)
			k++
		***REMOVED***
	***REMOVED***
	sl.Set(sl.Slice(0, k))
***REMOVED***

// Group finds elements in s for which fn returns the same value and groups
// them in a new Slice.
func (s Slice) Group(fn func(e Elem) string) []Slice ***REMOVED***
	m := make(map[string][]reflect.Value)
	sl := s.Value()
	for i := 0; i < sl.Len(); i++ ***REMOVED***
		vi := sl.Index(i)
		key := fn(vi.Interface().(Elem))
		m[key] = append(m[key], vi)
	***REMOVED***
	keys := []string***REMOVED******REMOVED***
	for k, _ := range m ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	sort.Strings(keys)
	res := []Slice***REMOVED******REMOVED***
	for _, k := range keys ***REMOVED***
		nsl := reflect.New(sl.Type())
		nsl.Elem().Set(reflect.Append(nsl.Elem(), m[k]...))
		res = append(res, MakeSlice(nsl.Interface()))
	***REMOVED***
	return res
***REMOVED***

// SelectAnyOf filters s to contain only elements for which attr matches
// any of the values.
func (s Slice) SelectAnyOf(attr string, values ...string) ***REMOVED***
	index := s.indexForAttr(attr)
	s.Filter(func(e Elem) bool ***REMOVED***
		vf := reflect.ValueOf(e).Elem().FieldByIndex(index)
		return in(values, vf.String())
	***REMOVED***)
***REMOVED***

// SelectOnePerGroup filters s to include at most one element e per group of
// elements matching Key(attr), where e has an attribute a that matches any
// the values in v.
// If more than one element in a group matches a value in v preference
// is given to the element that matches the first value in v.
func (s Slice) SelectOnePerGroup(a string, v []string) ***REMOVED***
	index := s.indexForAttr(a)
	grouped := s.Group(func(e Elem) string ***REMOVED*** return Key(e, a) ***REMOVED***)
	sl := s.Value()
	sl.Set(sl.Slice(0, 0))
	for _, g := range grouped ***REMOVED***
		e := reflect.Value***REMOVED******REMOVED***
		found := len(v)
		gsl := g.Value()
		for i := 0; i < gsl.Len(); i++ ***REMOVED***
			vi := gsl.Index(i).Elem().FieldByIndex(index)
			j := 0
			for ; j < len(v) && v[j] != vi.String(); j++ ***REMOVED***
			***REMOVED***
			if j < found ***REMOVED***
				found = j
				e = gsl.Index(i)
			***REMOVED***
		***REMOVED***
		if found < len(v) ***REMOVED***
			sl.Set(reflect.Append(sl, e))
		***REMOVED***
	***REMOVED***
***REMOVED***

// SelectDraft drops all elements from the list with a draft level smaller than d
// and selects the highest draft level of the remaining.
// This method assumes that the input CLDR is canonicalized.
func (s Slice) SelectDraft(d Draft) ***REMOVED***
	s.SelectOnePerGroup("draft", drafts[len(drafts)-2-int(d):])
***REMOVED***
