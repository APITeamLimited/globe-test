// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package yaml

import (
	"bytes"
	"encoding"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// if it encounters an Unmarshaler, indirect stops and returns that.
// if decodingNull is true, indirect stops at the last pointer so it can be set to nil.
func indirect(v reflect.Value, decodingNull bool) (json.Unmarshaler, encoding.TextUnmarshaler, reflect.Value) ***REMOVED***
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() ***REMOVED***
		v = v.Addr()
	***REMOVED***
	for ***REMOVED***
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() ***REMOVED***
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) ***REMOVED***
				v = e
				continue
			***REMOVED***
		***REMOVED***

		if v.Kind() != reflect.Ptr ***REMOVED***
			break
		***REMOVED***

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() ***REMOVED***
			break
		***REMOVED***
		if v.IsNil() ***REMOVED***
			if v.CanSet() ***REMOVED***
				v.Set(reflect.New(v.Type().Elem()))
			***REMOVED*** else ***REMOVED***
				v = reflect.New(v.Type().Elem())
			***REMOVED***
		***REMOVED***
		if v.Type().NumMethod() > 0 ***REMOVED***
			if u, ok := v.Interface().(json.Unmarshaler); ok ***REMOVED***
				return u, nil, reflect.Value***REMOVED******REMOVED***
			***REMOVED***
			if u, ok := v.Interface().(encoding.TextUnmarshaler); ok ***REMOVED***
				return nil, u, reflect.Value***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
		v = v.Elem()
	***REMOVED***
	return nil, nil, v
***REMOVED***

// A field represents a single field found in a struct.
type field struct ***REMOVED***
	name      string
	nameBytes []byte                 // []byte(name)
	equalFold func(s, t []byte) bool // bytes.EqualFold or equivalent

	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool
	quoted    bool
***REMOVED***

func fillField(f field) field ***REMOVED***
	f.nameBytes = []byte(f.name)
	f.equalFold = foldFunc(f.nameBytes)
	return f
***REMOVED***

// byName sorts field by name, breaking ties with depth,
// then breaking ties with "name came from json tag", then
// breaking ties with index sequence.
type byName []field

func (x byName) Len() int ***REMOVED*** return len(x) ***REMOVED***

func (x byName) Swap(i, j int) ***REMOVED*** x[i], x[j] = x[j], x[i] ***REMOVED***

func (x byName) Less(i, j int) bool ***REMOVED***
	if x[i].name != x[j].name ***REMOVED***
		return x[i].name < x[j].name
	***REMOVED***
	if len(x[i].index) != len(x[j].index) ***REMOVED***
		return len(x[i].index) < len(x[j].index)
	***REMOVED***
	if x[i].tag != x[j].tag ***REMOVED***
		return x[i].tag
	***REMOVED***
	return byIndex(x).Less(i, j)
***REMOVED***

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int ***REMOVED*** return len(x) ***REMOVED***

func (x byIndex) Swap(i, j int) ***REMOVED*** x[i], x[j] = x[j], x[i] ***REMOVED***

func (x byIndex) Less(i, j int) bool ***REMOVED***
	for k, xik := range x[i].index ***REMOVED***
		if k >= len(x[j].index) ***REMOVED***
			return false
		***REMOVED***
		if xik != x[j].index[k] ***REMOVED***
			return xik < x[j].index[k]
		***REMOVED***
	***REMOVED***
	return len(x[i].index) < len(x[j].index)
***REMOVED***

// typeFields returns a list of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typeFields(t reflect.Type) []field ***REMOVED***
	// Anonymous fields to explore at the current level and the next.
	current := []field***REMOVED******REMOVED***
	next := []field***REMOVED******REMOVED***typ: t***REMOVED******REMOVED***

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int***REMOVED******REMOVED***
	nextCount := map[reflect.Type]int***REMOVED******REMOVED***

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool***REMOVED******REMOVED***

	// Fields found.
	var fields []field

	for len(next) > 0 ***REMOVED***
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int***REMOVED******REMOVED***

		for _, f := range current ***REMOVED***
			if visited[f.typ] ***REMOVED***
				continue
			***REMOVED***
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ ***REMOVED***
				sf := f.typ.Field(i)
				if sf.PkgPath != "" ***REMOVED*** // unexported
					continue
				***REMOVED***
				tag := sf.Tag.Get("json")
				if tag == "-" ***REMOVED***
					continue
				***REMOVED***
				name, opts := parseTag(tag)
				if !isValidTag(name) ***REMOVED***
					name = ""
				***REMOVED***
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr ***REMOVED***
					// Follow pointer.
					ft = ft.Elem()
				***REMOVED***

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct ***REMOVED***
					tagged := name != ""
					if name == "" ***REMOVED***
						name = sf.Name
					***REMOVED***
					fields = append(fields, fillField(field***REMOVED***
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       ft,
						omitEmpty: opts.Contains("omitempty"),
						quoted:    opts.Contains("string"),
					***REMOVED***))
					if count[f.typ] > 1 ***REMOVED***
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					***REMOVED***
					continue
				***REMOVED***

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 ***REMOVED***
					next = append(next, fillField(field***REMOVED***name: ft.Name(), index: index, typ: ft***REMOVED***))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	sort.Sort(byName(fields))

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance ***REMOVED***
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ ***REMOVED***
			fj := fields[i+advance]
			if fj.name != name ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if advance == 1 ***REMOVED*** // Only one field with this name
			out = append(out, fi)
			continue
		***REMOVED***
		dominant, ok := dominantField(fields[i : i+advance])
		if ok ***REMOVED***
			out = append(out, dominant)
		***REMOVED***
	***REMOVED***

	fields = out
	sort.Sort(byIndex(fields))

	return fields
***REMOVED***

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's embedding rules, modified by the presence of
// JSON tags. If there are multiple top-level fields, the boolean
// will be false: This condition is an error in Go and we skip all
// the fields.
func dominantField(fields []field) (field, bool) ***REMOVED***
	// The fields are sorted in increasing index-length order. The winner
	// must therefore be one with the shortest index length. Drop all
	// longer entries, which is easy: just truncate the slice.
	length := len(fields[0].index)
	tagged := -1 // Index of first tagged field.
	for i, f := range fields ***REMOVED***
		if len(f.index) > length ***REMOVED***
			fields = fields[:i]
			break
		***REMOVED***
		if f.tag ***REMOVED***
			if tagged >= 0 ***REMOVED***
				// Multiple tagged fields at the same level: conflict.
				// Return no field.
				return field***REMOVED******REMOVED***, false
			***REMOVED***
			tagged = i
		***REMOVED***
	***REMOVED***
	if tagged >= 0 ***REMOVED***
		return fields[tagged], true
	***REMOVED***
	// All remaining fields have the same length. If there's more than one,
	// we have a conflict (two fields named "X" at the same level) and we
	// return no field.
	if len(fields) > 1 ***REMOVED***
		return field***REMOVED******REMOVED***, false
	***REMOVED***
	return fields[0], true
***REMOVED***

var fieldCache struct ***REMOVED***
	sync.RWMutex
	m map[reflect.Type][]field
***REMOVED***

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field ***REMOVED***
	fieldCache.RLock()
	f := fieldCache.m[t]
	fieldCache.RUnlock()
	if f != nil ***REMOVED***
		return f
	***REMOVED***

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = typeFields(t)
	if f == nil ***REMOVED***
		f = []field***REMOVED******REMOVED***
	***REMOVED***

	fieldCache.Lock()
	if fieldCache.m == nil ***REMOVED***
		fieldCache.m = map[reflect.Type][]field***REMOVED******REMOVED***
	***REMOVED***
	fieldCache.m[t] = f
	fieldCache.Unlock()
	return f
***REMOVED***

func isValidTag(s string) bool ***REMOVED***
	if s == "" ***REMOVED***
		return false
	***REMOVED***
	for _, c := range s ***REMOVED***
		switch ***REMOVED***
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_***REMOVED***|***REMOVED***~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

const (
	caseMask     = ^byte(0x20) // Mask to ignore case in ASCII.
	kelvin       = '\u212a'
	smallLongEss = '\u017f'
)

// foldFunc returns one of four different case folding equivalence
// functions, from most general (and slow) to fastest:
//
// 1) bytes.EqualFold, if the key s contains any non-ASCII UTF-8
// 2) equalFoldRight, if s contains special folding ASCII ('k', 'K', 's', 'S')
// 3) asciiEqualFold, no special, but includes non-letters (including _)
// 4) simpleLetterEqualFold, no specials, no non-letters.
//
// The letters S and K are special because they map to 3 runes, not just 2:
//  * S maps to s and to U+017F 'ſ' Latin small letter long s
//  * k maps to K and to U+212A 'K' Kelvin sign
// See http://play.golang.org/p/tTxjOc0OGo
//
// The returned function is specialized for matching against s and
// should only be given s. It's not curried for performance reasons.
func foldFunc(s []byte) func(s, t []byte) bool ***REMOVED***
	nonLetter := false
	special := false // special letter
	for _, b := range s ***REMOVED***
		if b >= utf8.RuneSelf ***REMOVED***
			return bytes.EqualFold
		***REMOVED***
		upper := b & caseMask
		if upper < 'A' || upper > 'Z' ***REMOVED***
			nonLetter = true
		***REMOVED*** else if upper == 'K' || upper == 'S' ***REMOVED***
			// See above for why these letters are special.
			special = true
		***REMOVED***
	***REMOVED***
	if special ***REMOVED***
		return equalFoldRight
	***REMOVED***
	if nonLetter ***REMOVED***
		return asciiEqualFold
	***REMOVED***
	return simpleLetterEqualFold
***REMOVED***

// equalFoldRight is a specialization of bytes.EqualFold when s is
// known to be all ASCII (including punctuation), but contains an 's',
// 'S', 'k', or 'K', requiring a Unicode fold on the bytes in t.
// See comments on foldFunc.
func equalFoldRight(s, t []byte) bool ***REMOVED***
	for _, sb := range s ***REMOVED***
		if len(t) == 0 ***REMOVED***
			return false
		***REMOVED***
		tb := t[0]
		if tb < utf8.RuneSelf ***REMOVED***
			if sb != tb ***REMOVED***
				sbUpper := sb & caseMask
				if 'A' <= sbUpper && sbUpper <= 'Z' ***REMOVED***
					if sbUpper != tb&caseMask ***REMOVED***
						return false
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			t = t[1:]
			continue
		***REMOVED***
		// sb is ASCII and t is not. t must be either kelvin
		// sign or long s; sb must be s, S, k, or K.
		tr, size := utf8.DecodeRune(t)
		switch sb ***REMOVED***
		case 's', 'S':
			if tr != smallLongEss ***REMOVED***
				return false
			***REMOVED***
		case 'k', 'K':
			if tr != kelvin ***REMOVED***
				return false
			***REMOVED***
		default:
			return false
		***REMOVED***
		t = t[size:]

	***REMOVED***
	if len(t) > 0 ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// asciiEqualFold is a specialization of bytes.EqualFold for use when
// s is all ASCII (but may contain non-letters) and contains no
// special-folding letters.
// See comments on foldFunc.
func asciiEqualFold(s, t []byte) bool ***REMOVED***
	if len(s) != len(t) ***REMOVED***
		return false
	***REMOVED***
	for i, sb := range s ***REMOVED***
		tb := t[i]
		if sb == tb ***REMOVED***
			continue
		***REMOVED***
		if ('a' <= sb && sb <= 'z') || ('A' <= sb && sb <= 'Z') ***REMOVED***
			if sb&caseMask != tb&caseMask ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// simpleLetterEqualFold is a specialization of bytes.EqualFold for
// use when s is all ASCII letters (no underscores, etc) and also
// doesn't contain 'k', 'K', 's', or 'S'.
// See comments on foldFunc.
func simpleLetterEqualFold(s, t []byte) bool ***REMOVED***
	if len(s) != len(t) ***REMOVED***
		return false
	***REMOVED***
	for i, b := range s ***REMOVED***
		if b&caseMask != t[i]&caseMask ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) ***REMOVED***
	if idx := strings.Index(tag, ","); idx != -1 ***REMOVED***
		return tag[:idx], tagOptions(tag[idx+1:])
	***REMOVED***
	return tag, tagOptions("")
***REMOVED***

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool ***REMOVED***
	if len(o) == 0 ***REMOVED***
		return false
	***REMOVED***
	s := string(o)
	for s != "" ***REMOVED***
		var next string
		i := strings.Index(s, ",")
		if i >= 0 ***REMOVED***
			s, next = s[:i], s[i+1:]
		***REMOVED***
		if s == optionName ***REMOVED***
			return true
		***REMOVED***
		s = next
	***REMOVED***
	return false
***REMOVED***
