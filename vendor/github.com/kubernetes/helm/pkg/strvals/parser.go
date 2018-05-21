/*
Copyright 2016 The Kubernetes Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package strvals

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
)

// ErrNotList indicates that a non-list was treated as a list.
var ErrNotList = errors.New("not a list")

// ToYAML takes a string of arguments and converts to a YAML document.
func ToYAML(s string) (string, error) ***REMOVED***
	m, err := Parse(s)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	d, err := yaml.Marshal(m)
	return strings.TrimSuffix(string(d), "\n"), err
***REMOVED***

// Parse parses a set line.
//
// A set line is of the form name1=value1,name2=value2
func Parse(s string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	vals := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	scanner := bytes.NewBufferString(s)
	t := newParser(scanner, vals, false)
	err := t.parse()
	return vals, err
***REMOVED***

// ParseString parses a set line and forces a string value.
//
// A set line is of the form name1=value1,name2=value2
func ParseString(s string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	vals := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	scanner := bytes.NewBufferString(s)
	t := newParser(scanner, vals, true)
	err := t.parse()
	return vals, err
***REMOVED***

// ParseInto parses a strvals line and merges the result into dest.
//
// If the strval string has a key that exists in dest, it overwrites the
// dest version.
func ParseInto(s string, dest map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	scanner := bytes.NewBufferString(s)
	t := newParser(scanner, dest, false)
	return t.parse()
***REMOVED***

// ParseIntoString parses a strvals line nad merges the result into dest.
//
// This method always returns a string as the value.
func ParseIntoString(s string, dest map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	scanner := bytes.NewBufferString(s)
	t := newParser(scanner, dest, true)
	return t.parse()
***REMOVED***

// parser is a simple parser that takes a strvals line and parses it into a
// map representation.
type parser struct ***REMOVED***
	sc   *bytes.Buffer
	data map[string]interface***REMOVED******REMOVED***
	st   bool
***REMOVED***

func newParser(sc *bytes.Buffer, data map[string]interface***REMOVED******REMOVED***, stringBool bool) *parser ***REMOVED***
	return &parser***REMOVED***sc: sc, data: data, st: stringBool***REMOVED***
***REMOVED***

func (t *parser) parse() error ***REMOVED***
	for ***REMOVED***
		err := t.key(t.data)
		if err == nil ***REMOVED***
			continue
		***REMOVED***
		if err == io.EOF ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
***REMOVED***

func runeSet(r []rune) map[rune]bool ***REMOVED***
	s := make(map[rune]bool, len(r))
	for _, rr := range r ***REMOVED***
		s[rr] = true
	***REMOVED***
	return s
***REMOVED***

func (t *parser) key(data map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	stop := runeSet([]rune***REMOVED***'=', '[', ',', '.'***REMOVED***)
	for ***REMOVED***
		switch k, last, err := runesUntil(t.sc, stop); ***REMOVED***
		case err != nil:
			if len(k) == 0 ***REMOVED***
				return err
			***REMOVED***
			return fmt.Errorf("key %q has no value", string(k))
			//set(data, string(k), "")
			//return err
		case last == '[':
			// We are in a list index context, so we need to set an index.
			i, err := t.keyIndex()
			if err != nil ***REMOVED***
				return fmt.Errorf("error parsing index: %s", err)
			***REMOVED***
			kk := string(k)
			// Find or create target list
			list := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
			if _, ok := data[kk]; ok ***REMOVED***
				list = data[kk].([]interface***REMOVED******REMOVED***)
			***REMOVED***

			// Now we need to get the value after the ].
			list, err = t.listItem(list, i)
			set(data, kk, list)
			return err
		case last == '=':
			//End of key. Consume =, Get value.
			// FIXME: Get value list first
			vl, e := t.valList()
			switch e ***REMOVED***
			case nil:
				set(data, string(k), vl)
				return nil
			case io.EOF:
				set(data, string(k), "")
				return e
			case ErrNotList:
				v, e := t.val()
				set(data, string(k), typedVal(v, t.st))
				return e
			default:
				return e
			***REMOVED***

		case last == ',':
			// No value given. Set the value to empty string. Return error.
			set(data, string(k), "")
			return fmt.Errorf("key %q has no value (cannot end with ,)", string(k))
		case last == '.':
			// First, create or find the target map.
			inner := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
			if _, ok := data[string(k)]; ok ***REMOVED***
				inner = data[string(k)].(map[string]interface***REMOVED******REMOVED***)
			***REMOVED***

			// Recurse
			e := t.key(inner)
			if len(inner) == 0 ***REMOVED***
				return fmt.Errorf("key map %q has no value", string(k))
			***REMOVED***
			set(data, string(k), inner)
			return e
		***REMOVED***
	***REMOVED***
***REMOVED***

func set(data map[string]interface***REMOVED******REMOVED***, key string, val interface***REMOVED******REMOVED***) ***REMOVED***
	// If key is empty, don't set it.
	if len(key) == 0 ***REMOVED***
		return
	***REMOVED***
	data[key] = val
***REMOVED***

func setIndex(list []interface***REMOVED******REMOVED***, index int, val interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	if len(list) <= index ***REMOVED***
		newlist := make([]interface***REMOVED******REMOVED***, index+1)
		copy(newlist, list)
		list = newlist
	***REMOVED***
	list[index] = val
	return list
***REMOVED***

func (t *parser) keyIndex() (int, error) ***REMOVED***
	// First, get the key.
	stop := runeSet([]rune***REMOVED***']'***REMOVED***)
	v, _, err := runesUntil(t.sc, stop)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	// v should be the index
	return strconv.Atoi(string(v))

***REMOVED***
func (t *parser) listItem(list []interface***REMOVED******REMOVED***, i int) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	stop := runeSet([]rune***REMOVED***'[', '.', '='***REMOVED***)
	switch k, last, err := runesUntil(t.sc, stop); ***REMOVED***
	case len(k) > 0:
		return list, fmt.Errorf("unexpected data at end of array index: %q", k)
	case err != nil:
		return list, err
	case last == '=':
		vl, e := t.valList()
		switch e ***REMOVED***
		case nil:
			return setIndex(list, i, vl), nil
		case io.EOF:
			return setIndex(list, i, ""), err
		case ErrNotList:
			v, e := t.val()
			return setIndex(list, i, typedVal(v, t.st)), e
		default:
			return list, e
		***REMOVED***
	case last == '[':
		// now we have a nested list. Read the index and handle.
		i, err := t.keyIndex()
		if err != nil ***REMOVED***
			return list, fmt.Errorf("error parsing index: %s", err)
		***REMOVED***
		// Now we need to get the value after the ].
		list2, err := t.listItem(list, i)
		return setIndex(list, i, list2), err
	case last == '.':
		// We have a nested object. Send to t.key
		inner := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
		if len(list) > i ***REMOVED***
			inner = list[i].(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***

		// Recurse
		e := t.key(inner)
		return setIndex(list, i, inner), e
	default:
		return nil, fmt.Errorf("parse error: unexpected token %v", last)
	***REMOVED***
***REMOVED***

func (t *parser) val() ([]rune, error) ***REMOVED***
	stop := runeSet([]rune***REMOVED***','***REMOVED***)
	v, _, err := runesUntil(t.sc, stop)
	return v, err
***REMOVED***

func (t *parser) valList() ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	r, _, e := t.sc.ReadRune()
	if e != nil ***REMOVED***
		return []interface***REMOVED******REMOVED******REMOVED******REMOVED***, e
	***REMOVED***

	if r != '***REMOVED***' ***REMOVED***
		t.sc.UnreadRune()
		return []interface***REMOVED******REMOVED******REMOVED******REMOVED***, ErrNotList
	***REMOVED***

	list := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	stop := runeSet([]rune***REMOVED***',', '***REMOVED***'***REMOVED***)
	for ***REMOVED***
		switch v, last, err := runesUntil(t.sc, stop); ***REMOVED***
		case err != nil:
			if err == io.EOF ***REMOVED***
				err = errors.New("list must terminate with '***REMOVED***'")
			***REMOVED***
			return list, err
		case last == '***REMOVED***':
			// If this is followed by ',', consume it.
			if r, _, e := t.sc.ReadRune(); e == nil && r != ',' ***REMOVED***
				t.sc.UnreadRune()
			***REMOVED***
			list = append(list, typedVal(v, t.st))
			return list, nil
		case last == ',':
			list = append(list, typedVal(v, t.st))
		***REMOVED***
	***REMOVED***
***REMOVED***

func runesUntil(in io.RuneReader, stop map[rune]bool) ([]rune, rune, error) ***REMOVED***
	v := []rune***REMOVED******REMOVED***
	for ***REMOVED***
		switch r, _, e := in.ReadRune(); ***REMOVED***
		case e != nil:
			return v, r, e
		case inMap(r, stop):
			return v, r, nil
		case r == '\\':
			next, _, e := in.ReadRune()
			if e != nil ***REMOVED***
				return v, next, e
			***REMOVED***
			v = append(v, next)
		default:
			v = append(v, r)
		***REMOVED***
	***REMOVED***
***REMOVED***

func inMap(k rune, m map[rune]bool) bool ***REMOVED***
	_, ok := m[k]
	return ok
***REMOVED***

func typedVal(v []rune, st bool) interface***REMOVED******REMOVED*** ***REMOVED***
	val := string(v)
	if strings.EqualFold(val, "true") ***REMOVED***
		return true
	***REMOVED***

	if strings.EqualFold(val, "false") ***REMOVED***
		return false
	***REMOVED***

	// If this value does not start with zero, and not returnString, try parsing it to an int
	if !st && len(val) != 0 && val[0] != '0' ***REMOVED***
		if iv, err := strconv.ParseInt(val, 10, 64); err == nil ***REMOVED***
			return iv
		***REMOVED***
	***REMOVED***

	return val
***REMOVED***
