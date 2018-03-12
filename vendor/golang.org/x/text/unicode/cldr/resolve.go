// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

// This file implements the various inheritance constructs defined by LDML.
// See http://www.unicode.org/reports/tr35/#Inheritance_and_Validity
// for more details.

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// fieldIter iterates over fields in a struct. It includes
// fields of embedded structs.
type fieldIter struct ***REMOVED***
	v        reflect.Value
	index, n []int
***REMOVED***

func iter(v reflect.Value) fieldIter ***REMOVED***
	if v.Kind() != reflect.Struct ***REMOVED***
		log.Panicf("value %v must be a struct", v)
	***REMOVED***
	i := fieldIter***REMOVED***
		v:     v,
		index: []int***REMOVED***0***REMOVED***,
		n:     []int***REMOVED***v.NumField()***REMOVED***,
	***REMOVED***
	i.descent()
	return i
***REMOVED***

func (i *fieldIter) descent() ***REMOVED***
	for f := i.field(); f.Anonymous && f.Type.NumField() > 0; f = i.field() ***REMOVED***
		i.index = append(i.index, 0)
		i.n = append(i.n, f.Type.NumField())
	***REMOVED***
***REMOVED***

func (i *fieldIter) done() bool ***REMOVED***
	return len(i.index) == 1 && i.index[0] >= i.n[0]
***REMOVED***

func skip(f reflect.StructField) bool ***REMOVED***
	return !f.Anonymous && (f.Name[0] < 'A' || f.Name[0] > 'Z')
***REMOVED***

func (i *fieldIter) next() ***REMOVED***
	for ***REMOVED***
		k := len(i.index) - 1
		i.index[k]++
		if i.index[k] < i.n[k] ***REMOVED***
			if !skip(i.field()) ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if k == 0 ***REMOVED***
				return
			***REMOVED***
			i.index = i.index[:k]
			i.n = i.n[:k]
		***REMOVED***
	***REMOVED***
	i.descent()
***REMOVED***

func (i *fieldIter) value() reflect.Value ***REMOVED***
	return i.v.FieldByIndex(i.index)
***REMOVED***

func (i *fieldIter) field() reflect.StructField ***REMOVED***
	return i.v.Type().FieldByIndex(i.index)
***REMOVED***

type visitor func(v reflect.Value) error

var stopDescent = fmt.Errorf("do not recurse")

func (f visitor) visit(x interface***REMOVED******REMOVED***) error ***REMOVED***
	return f.visitRec(reflect.ValueOf(x))
***REMOVED***

// visit recursively calls f on all nodes in v.
func (f visitor) visitRec(v reflect.Value) error ***REMOVED***
	if v.Kind() == reflect.Ptr ***REMOVED***
		if v.IsNil() ***REMOVED***
			return nil
		***REMOVED***
		return f.visitRec(v.Elem())
	***REMOVED***
	if err := f(v); err != nil ***REMOVED***
		if err == stopDescent ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Struct:
		for i := iter(v); !i.done(); i.next() ***REMOVED***
			if err := f.visitRec(i.value()); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ ***REMOVED***
			if err := f.visitRec(v.Index(i)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// getPath is used for error reporting purposes only.
func getPath(e Elem) string ***REMOVED***
	if e == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	if e.enclosing() == nil ***REMOVED***
		return e.GetCommon().name
	***REMOVED***
	if e.GetCommon().Type == "" ***REMOVED***
		return fmt.Sprintf("%s.%s", getPath(e.enclosing()), e.GetCommon().name)
	***REMOVED***
	return fmt.Sprintf("%s.%s[type=%s]", getPath(e.enclosing()), e.GetCommon().name, e.GetCommon().Type)
***REMOVED***

// xmlName returns the xml name of the element or attribute
func xmlName(f reflect.StructField) (name string, attr bool) ***REMOVED***
	tags := strings.Split(f.Tag.Get("xml"), ",")
	for _, s := range tags ***REMOVED***
		attr = attr || s == "attr"
	***REMOVED***
	return tags[0], attr
***REMOVED***

func findField(v reflect.Value, key string) (reflect.Value, error) ***REMOVED***
	v = reflect.Indirect(v)
	for i := iter(v); !i.done(); i.next() ***REMOVED***
		if n, _ := xmlName(i.field()); n == key ***REMOVED***
			return i.value(), nil
		***REMOVED***
	***REMOVED***
	return reflect.Value***REMOVED******REMOVED***, fmt.Errorf("cldr: no field %q in element %#v", key, v.Interface())
***REMOVED***

var xpathPart = regexp.MustCompile(`(\pL+)(?:\[@(\pL+)='([\w-]+)'\])?`)

func walkXPath(e Elem, path string) (res Elem, err error) ***REMOVED***
	for _, c := range strings.Split(path, "/") ***REMOVED***
		if c == ".." ***REMOVED***
			if e = e.enclosing(); e == nil ***REMOVED***
				panic("path ..")
				return nil, fmt.Errorf(`cldr: ".." moves past root in path %q`, path)
			***REMOVED***
			continue
		***REMOVED*** else if c == "" ***REMOVED***
			continue
		***REMOVED***
		m := xpathPart.FindStringSubmatch(c)
		if len(m) == 0 || len(m[0]) != len(c) ***REMOVED***
			return nil, fmt.Errorf("cldr: syntax error in path component %q", c)
		***REMOVED***
		v, err := findField(reflect.ValueOf(e), m[1])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		switch v.Kind() ***REMOVED***
		case reflect.Slice:
			i := 0
			if m[2] != "" || v.Len() > 1 ***REMOVED***
				if m[2] == "" ***REMOVED***
					m[2] = "type"
					if m[3] = e.GetCommon().Default(); m[3] == "" ***REMOVED***
						return nil, fmt.Errorf("cldr: type selector or default value needed for element %s", m[1])
					***REMOVED***
				***REMOVED***
				for ; i < v.Len(); i++ ***REMOVED***
					vi := v.Index(i)
					key, err := findField(vi.Elem(), m[2])
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					key = reflect.Indirect(key)
					if key.Kind() == reflect.String && key.String() == m[3] ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if i == v.Len() || v.Index(i).IsNil() ***REMOVED***
				return nil, fmt.Errorf("no %s found with %s==%s", m[1], m[2], m[3])
			***REMOVED***
			e = v.Index(i).Interface().(Elem)
		case reflect.Ptr:
			if v.IsNil() ***REMOVED***
				return nil, fmt.Errorf("cldr: element %q not found within element %q", m[1], e.GetCommon().name)
			***REMOVED***
			var ok bool
			if e, ok = v.Interface().(Elem); !ok ***REMOVED***
				return nil, fmt.Errorf("cldr: %q is not an XML element", m[1])
			***REMOVED*** else if m[2] != "" || m[3] != "" ***REMOVED***
				return nil, fmt.Errorf("cldr: no type selector allowed for element %s", m[1])
			***REMOVED***
		default:
			return nil, fmt.Errorf("cldr: %q is not an XML element", m[1])
		***REMOVED***
	***REMOVED***
	return e, nil
***REMOVED***

const absPrefix = "//ldml/"

func (cldr *CLDR) resolveAlias(e Elem, src, path string) (res Elem, err error) ***REMOVED***
	if src != "locale" ***REMOVED***
		if !strings.HasPrefix(path, absPrefix) ***REMOVED***
			return nil, fmt.Errorf("cldr: expected absolute path, found %q", path)
		***REMOVED***
		path = path[len(absPrefix):]
		if e, err = cldr.resolve(src); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return walkXPath(e, path)
***REMOVED***

func (cldr *CLDR) resolveAndMergeAlias(e Elem) error ***REMOVED***
	alias := e.GetCommon().Alias
	if alias == nil ***REMOVED***
		return nil
	***REMOVED***
	a, err := cldr.resolveAlias(e, alias.Source, alias.Path)
	if err != nil ***REMOVED***
		return fmt.Errorf("%v: error evaluating path %q: %v", getPath(e), alias.Path, err)
	***REMOVED***
	// Ensure alias node was already evaluated. TODO: avoid double evaluation.
	err = cldr.resolveAndMergeAlias(a)
	v := reflect.ValueOf(e).Elem()
	for i := iter(reflect.ValueOf(a).Elem()); !i.done(); i.next() ***REMOVED***
		if vv := i.value(); vv.Kind() != reflect.Ptr || !vv.IsNil() ***REMOVED***
			if _, attr := xmlName(i.field()); !attr ***REMOVED***
				v.FieldByIndex(i.index).Set(vv)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (cldr *CLDR) aliasResolver() visitor ***REMOVED***
	return func(v reflect.Value) (err error) ***REMOVED***
		if e, ok := v.Addr().Interface().(Elem); ok ***REMOVED***
			err = cldr.resolveAndMergeAlias(e)
			if err == nil && blocking[e.GetCommon().name] ***REMOVED***
				return stopDescent
			***REMOVED***
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

// elements within blocking elements do not inherit.
// Taken from CLDR's supplementalMetaData.xml.
var blocking = map[string]bool***REMOVED***
	"identity":         true,
	"supplementalData": true,
	"cldrTest":         true,
	"collation":        true,
	"transform":        true,
***REMOVED***

// Distinguishing attributes affect inheritance; two elements with different
// distinguishing attributes are treated as different for purposes of inheritance,
// except when such attributes occur in the indicated elements.
// Taken from CLDR's supplementalMetaData.xml.
var distinguishing = map[string][]string***REMOVED***
	"key":        nil,
	"request_id": nil,
	"id":         nil,
	"registry":   nil,
	"alt":        nil,
	"iso4217":    nil,
	"iso3166":    nil,
	"mzone":      nil,
	"from":       nil,
	"to":         nil,
	"type": []string***REMOVED***
		"abbreviationFallback",
		"default",
		"mapping",
		"measurementSystem",
		"preferenceOrdering",
	***REMOVED***,
	"numberSystem": nil,
***REMOVED***

func in(set []string, s string) bool ***REMOVED***
	for _, v := range set ***REMOVED***
		if v == s ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// attrKey computes a key based on the distinguishable attributes of
// an element and it's values.
func attrKey(v reflect.Value, exclude ...string) string ***REMOVED***
	parts := []string***REMOVED******REMOVED***
	ename := v.Interface().(Elem).GetCommon().name
	v = v.Elem()
	for i := iter(v); !i.done(); i.next() ***REMOVED***
		if name, attr := xmlName(i.field()); attr ***REMOVED***
			if except, ok := distinguishing[name]; ok && !in(exclude, name) && !in(except, ename) ***REMOVED***
				v := i.value()
				if v.Kind() == reflect.Ptr ***REMOVED***
					v = v.Elem()
				***REMOVED***
				if v.IsValid() ***REMOVED***
					parts = append(parts, fmt.Sprintf("%s=%s", name, v.String()))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	sort.Strings(parts)
	return strings.Join(parts, ";")
***REMOVED***

// Key returns a key for e derived from all distinguishing attributes
// except those specified by exclude.
func Key(e Elem, exclude ...string) string ***REMOVED***
	return attrKey(reflect.ValueOf(e), exclude...)
***REMOVED***

// linkEnclosing sets the enclosing element as well as the name
// for all sub-elements of child, recursively.
func linkEnclosing(parent, child Elem) ***REMOVED***
	child.setEnclosing(parent)
	v := reflect.ValueOf(child).Elem()
	for i := iter(v); !i.done(); i.next() ***REMOVED***
		vf := i.value()
		if vf.Kind() == reflect.Slice ***REMOVED***
			for j := 0; j < vf.Len(); j++ ***REMOVED***
				linkEnclosing(child, vf.Index(j).Interface().(Elem))
			***REMOVED***
		***REMOVED*** else if vf.Kind() == reflect.Ptr && !vf.IsNil() && vf.Elem().Kind() == reflect.Struct ***REMOVED***
			linkEnclosing(child, vf.Interface().(Elem))
		***REMOVED***
	***REMOVED***
***REMOVED***

func setNames(e Elem, name string) ***REMOVED***
	e.setName(name)
	v := reflect.ValueOf(e).Elem()
	for i := iter(v); !i.done(); i.next() ***REMOVED***
		vf := i.value()
		name, _ = xmlName(i.field())
		if vf.Kind() == reflect.Slice ***REMOVED***
			for j := 0; j < vf.Len(); j++ ***REMOVED***
				setNames(vf.Index(j).Interface().(Elem), name)
			***REMOVED***
		***REMOVED*** else if vf.Kind() == reflect.Ptr && !vf.IsNil() && vf.Elem().Kind() == reflect.Struct ***REMOVED***
			setNames(vf.Interface().(Elem), name)
		***REMOVED***
	***REMOVED***
***REMOVED***

// deepCopy copies elements of v recursively.  All elements of v that may
// be modified by inheritance are explicitly copied.
func deepCopy(v reflect.Value) reflect.Value ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Ptr:
		if v.IsNil() || v.Elem().Kind() != reflect.Struct ***REMOVED***
			return v
		***REMOVED***
		nv := reflect.New(v.Elem().Type())
		nv.Elem().Set(v.Elem())
		deepCopyRec(nv.Elem(), v.Elem())
		return nv
	case reflect.Slice:
		nv := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ ***REMOVED***
			deepCopyRec(nv.Index(i), v.Index(i))
		***REMOVED***
		return nv
	***REMOVED***
	panic("deepCopy: must be called with pointer or slice")
***REMOVED***

// deepCopyRec is only called by deepCopy.
func deepCopyRec(nv, v reflect.Value) ***REMOVED***
	if v.Kind() == reflect.Struct ***REMOVED***
		t := v.Type()
		for i := 0; i < v.NumField(); i++ ***REMOVED***
			if name, attr := xmlName(t.Field(i)); name != "" && !attr ***REMOVED***
				deepCopyRec(nv.Field(i), v.Field(i))
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		nv.Set(deepCopy(v))
	***REMOVED***
***REMOVED***

// newNode is used to insert a missing node during inheritance.
func (cldr *CLDR) newNode(v, enc reflect.Value) reflect.Value ***REMOVED***
	n := reflect.New(v.Type())
	for i := iter(v); !i.done(); i.next() ***REMOVED***
		if name, attr := xmlName(i.field()); name == "" || attr ***REMOVED***
			n.Elem().FieldByIndex(i.index).Set(i.value())
		***REMOVED***
	***REMOVED***
	n.Interface().(Elem).GetCommon().setEnclosing(enc.Addr().Interface().(Elem))
	return n
***REMOVED***

// v, parent must be pointers to struct
func (cldr *CLDR) inheritFields(v, parent reflect.Value) (res reflect.Value, err error) ***REMOVED***
	t := v.Type()
	nv := reflect.New(t)
	nv.Elem().Set(v)
	for i := iter(v); !i.done(); i.next() ***REMOVED***
		vf := i.value()
		f := i.field()
		name, attr := xmlName(f)
		if name == "" || attr ***REMOVED***
			continue
		***REMOVED***
		pf := parent.FieldByIndex(i.index)
		if blocking[name] ***REMOVED***
			if vf.IsNil() ***REMOVED***
				vf = pf
			***REMOVED***
			nv.Elem().FieldByIndex(i.index).Set(deepCopy(vf))
			continue
		***REMOVED***
		switch f.Type.Kind() ***REMOVED***
		case reflect.Ptr:
			if f.Type.Elem().Kind() == reflect.Struct ***REMOVED***
				if !vf.IsNil() ***REMOVED***
					if vf, err = cldr.inheritStructPtr(vf, pf); err != nil ***REMOVED***
						return reflect.Value***REMOVED******REMOVED***, err
					***REMOVED***
					vf.Interface().(Elem).setEnclosing(nv.Interface().(Elem))
					nv.Elem().FieldByIndex(i.index).Set(vf)
				***REMOVED*** else if !pf.IsNil() ***REMOVED***
					n := cldr.newNode(pf.Elem(), v)
					if vf, err = cldr.inheritStructPtr(n, pf); err != nil ***REMOVED***
						return reflect.Value***REMOVED******REMOVED***, err
					***REMOVED***
					vf.Interface().(Elem).setEnclosing(nv.Interface().(Elem))
					nv.Elem().FieldByIndex(i.index).Set(vf)
				***REMOVED***
			***REMOVED***
		case reflect.Slice:
			vf, err := cldr.inheritSlice(nv.Elem(), vf, pf)
			if err != nil ***REMOVED***
				return reflect.Zero(t), err
			***REMOVED***
			nv.Elem().FieldByIndex(i.index).Set(vf)
		***REMOVED***
	***REMOVED***
	return nv, nil
***REMOVED***

func root(e Elem) *LDML ***REMOVED***
	for ; e.enclosing() != nil; e = e.enclosing() ***REMOVED***
	***REMOVED***
	return e.(*LDML)
***REMOVED***

// inheritStructPtr first merges possible aliases in with v and then inherits
// any underspecified elements from parent.
func (cldr *CLDR) inheritStructPtr(v, parent reflect.Value) (r reflect.Value, err error) ***REMOVED***
	if !v.IsNil() ***REMOVED***
		e := v.Interface().(Elem).GetCommon()
		alias := e.Alias
		if alias == nil && !parent.IsNil() ***REMOVED***
			alias = parent.Interface().(Elem).GetCommon().Alias
		***REMOVED***
		if alias != nil ***REMOVED***
			a, err := cldr.resolveAlias(v.Interface().(Elem), alias.Source, alias.Path)
			if a != nil ***REMOVED***
				if v, err = cldr.inheritFields(v.Elem(), reflect.ValueOf(a).Elem()); err != nil ***REMOVED***
					return reflect.Value***REMOVED******REMOVED***, err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if !parent.IsNil() ***REMOVED***
			return cldr.inheritFields(v.Elem(), parent.Elem())
		***REMOVED***
	***REMOVED*** else if parent.IsNil() ***REMOVED***
		panic("should not reach here")
	***REMOVED***
	return v, nil
***REMOVED***

// Must be slice of struct pointers.
func (cldr *CLDR) inheritSlice(enc, v, parent reflect.Value) (res reflect.Value, err error) ***REMOVED***
	t := v.Type()
	index := make(map[string]reflect.Value)
	if !v.IsNil() ***REMOVED***
		for i := 0; i < v.Len(); i++ ***REMOVED***
			vi := v.Index(i)
			key := attrKey(vi)
			index[key] = vi
		***REMOVED***
	***REMOVED***
	if !parent.IsNil() ***REMOVED***
		for i := 0; i < parent.Len(); i++ ***REMOVED***
			vi := parent.Index(i)
			key := attrKey(vi)
			if w, ok := index[key]; ok ***REMOVED***
				index[key], err = cldr.inheritStructPtr(w, vi)
			***REMOVED*** else ***REMOVED***
				n := cldr.newNode(vi.Elem(), enc)
				index[key], err = cldr.inheritStructPtr(n, vi)
			***REMOVED***
			index[key].Interface().(Elem).setEnclosing(enc.Addr().Interface().(Elem))
			if err != nil ***REMOVED***
				return v, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	keys := make([]string, 0, len(index))
	for k, _ := range index ***REMOVED***
		keys = append(keys, k)
	***REMOVED***
	sort.Strings(keys)
	sl := reflect.MakeSlice(t, len(index), len(index))
	for i, k := range keys ***REMOVED***
		sl.Index(i).Set(index[k])
	***REMOVED***
	return sl, nil
***REMOVED***

func parentLocale(loc string) string ***REMOVED***
	parts := strings.Split(loc, "_")
	if len(parts) == 1 ***REMOVED***
		return "root"
	***REMOVED***
	parts = parts[:len(parts)-1]
	key := strings.Join(parts, "_")
	return key
***REMOVED***

func (cldr *CLDR) resolve(loc string) (res *LDML, err error) ***REMOVED***
	if r := cldr.resolved[loc]; r != nil ***REMOVED***
		return r, nil
	***REMOVED***
	x := cldr.RawLDML(loc)
	if x == nil ***REMOVED***
		return nil, fmt.Errorf("cldr: unknown locale %q", loc)
	***REMOVED***
	var v reflect.Value
	if loc == "root" ***REMOVED***
		x = deepCopy(reflect.ValueOf(x)).Interface().(*LDML)
		linkEnclosing(nil, x)
		err = cldr.aliasResolver().visit(x)
	***REMOVED*** else ***REMOVED***
		key := parentLocale(loc)
		var parent *LDML
		for ; cldr.locale[key] == nil; key = parentLocale(key) ***REMOVED***
		***REMOVED***
		if parent, err = cldr.resolve(key); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		v, err = cldr.inheritFields(reflect.ValueOf(x).Elem(), reflect.ValueOf(parent).Elem())
		x = v.Interface().(*LDML)
		linkEnclosing(nil, x)
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cldr.resolved[loc] = x
	return x, err
***REMOVED***

// finalize finalizes the initialization of the raw LDML structs.  It also
// removed unwanted fields, as specified by filter, so that they will not
// be unnecessarily evaluated.
func (cldr *CLDR) finalize(filter []string) ***REMOVED***
	for _, x := range cldr.locale ***REMOVED***
		if filter != nil ***REMOVED***
			v := reflect.ValueOf(x).Elem()
			t := v.Type()
			for i := 0; i < v.NumField(); i++ ***REMOVED***
				f := t.Field(i)
				name, _ := xmlName(f)
				if name != "" && name != "identity" && !in(filter, name) ***REMOVED***
					v.Field(i).Set(reflect.Zero(f.Type))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		linkEnclosing(nil, x) // for resolving aliases and paths
		setNames(x, "ldml")
	***REMOVED***
***REMOVED***
