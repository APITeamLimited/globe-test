/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package log

import "fmt"

type token struct ***REMOVED***
	key, value string
	inside     rune // shows whether it's inside a given collection, currently [ means it's an array
***REMOVED***

type tokenizer struct ***REMOVED***
	i          int
	s          string
	currentKey string
***REMOVED***

func (t *tokenizer) readKey() (string, error) ***REMOVED***
	start := t.i
	for ; t.i < len(t.s); t.i++ ***REMOVED***
		if t.s[t.i] == '=' && t.i != len(t.s)-1 ***REMOVED***
			t.i++

			return t.s[start : t.i-1], nil
		***REMOVED***
		if t.s[t.i] == ',' ***REMOVED***
			k := t.s[start:t.i]

			return k, fmt.Errorf("key `%s` with no value", k)
		***REMOVED***
	***REMOVED***

	s := t.s[start:]

	return s, fmt.Errorf("key `%s` with no value", s)
***REMOVED***

func (t *tokenizer) readValue() string ***REMOVED***
	start := t.i
	for ; t.i < len(t.s); t.i++ ***REMOVED***
		if t.s[t.i] == ',' ***REMOVED***
			t.i++

			return t.s[start : t.i-1]
		***REMOVED***
	***REMOVED***

	return t.s[start:]
***REMOVED***

func (t *tokenizer) readArray() (string, error) ***REMOVED***
	start := t.i
	for ; t.i < len(t.s); t.i++ ***REMOVED***
		if t.s[t.i] == ']' ***REMOVED***
			if t.i+1 == len(t.s) || t.s[t.i+1] == ',' ***REMOVED***
				t.i += 2

				return t.s[start : t.i-2], nil
			***REMOVED***
			t.i++

			return t.s[start : t.i-1], fmt.Errorf("there was no ',' after an array with key '%s'", t.currentKey)
		***REMOVED***
	***REMOVED***

	return t.s[start:], fmt.Errorf("array value for key `%s` didn't end", t.currentKey)
***REMOVED***

func tokenize(s string) ([]token, error) ***REMOVED***
	result := []token***REMOVED******REMOVED***
	t := &tokenizer***REMOVED***s: s***REMOVED***

	var err error
	var value string
	for t.i < len(s) ***REMOVED***
		t.currentKey, err = t.readKey()
		if err != nil ***REMOVED***
			return result, err
		***REMOVED***
		if t.s[t.i] == '[' ***REMOVED***
			t.i++
			value, err = t.readArray()

			result = append(result, token***REMOVED***
				key:    t.currentKey,
				value:  value,
				inside: '[',
			***REMOVED***)
			if err != nil ***REMOVED***
				return result, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			value = t.readValue()
			result = append(result, token***REMOVED***
				key:   t.currentKey,
				value: value,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return result, nil
***REMOVED***
