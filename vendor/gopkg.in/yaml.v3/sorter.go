//
// Copyright (c) 2011-2019 Canonical Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

import (
	"reflect"
	"unicode"
)

type keyList []reflect.Value

func (l keyList) Len() int      ***REMOVED*** return len(l) ***REMOVED***
func (l keyList) Swap(i, j int) ***REMOVED*** l[i], l[j] = l[j], l[i] ***REMOVED***
func (l keyList) Less(i, j int) bool ***REMOVED***
	a := l[i]
	b := l[j]
	ak := a.Kind()
	bk := b.Kind()
	for (ak == reflect.Interface || ak == reflect.Ptr) && !a.IsNil() ***REMOVED***
		a = a.Elem()
		ak = a.Kind()
	***REMOVED***
	for (bk == reflect.Interface || bk == reflect.Ptr) && !b.IsNil() ***REMOVED***
		b = b.Elem()
		bk = b.Kind()
	***REMOVED***
	af, aok := keyFloat(a)
	bf, bok := keyFloat(b)
	if aok && bok ***REMOVED***
		if af != bf ***REMOVED***
			return af < bf
		***REMOVED***
		if ak != bk ***REMOVED***
			return ak < bk
		***REMOVED***
		return numLess(a, b)
	***REMOVED***
	if ak != reflect.String || bk != reflect.String ***REMOVED***
		return ak < bk
	***REMOVED***
	ar, br := []rune(a.String()), []rune(b.String())
	digits := false
	for i := 0; i < len(ar) && i < len(br); i++ ***REMOVED***
		if ar[i] == br[i] ***REMOVED***
			digits = unicode.IsDigit(ar[i])
			continue
		***REMOVED***
		al := unicode.IsLetter(ar[i])
		bl := unicode.IsLetter(br[i])
		if al && bl ***REMOVED***
			return ar[i] < br[i]
		***REMOVED***
		if al || bl ***REMOVED***
			if digits ***REMOVED***
				return al
			***REMOVED*** else ***REMOVED***
				return bl
			***REMOVED***
		***REMOVED***
		var ai, bi int
		var an, bn int64
		if ar[i] == '0' || br[i] == '0' ***REMOVED***
			for j := i - 1; j >= 0 && unicode.IsDigit(ar[j]); j-- ***REMOVED***
				if ar[j] != '0' ***REMOVED***
					an = 1
					bn = 1
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for ai = i; ai < len(ar) && unicode.IsDigit(ar[ai]); ai++ ***REMOVED***
			an = an*10 + int64(ar[ai]-'0')
		***REMOVED***
		for bi = i; bi < len(br) && unicode.IsDigit(br[bi]); bi++ ***REMOVED***
			bn = bn*10 + int64(br[bi]-'0')
		***REMOVED***
		if an != bn ***REMOVED***
			return an < bn
		***REMOVED***
		if ai != bi ***REMOVED***
			return ai < bi
		***REMOVED***
		return ar[i] < br[i]
	***REMOVED***
	return len(ar) < len(br)
***REMOVED***

// keyFloat returns a float value for v if it is a number/bool
// and whether it is a number/bool or not.
func keyFloat(v reflect.Value) (f float64, ok bool) ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint()), true
	case reflect.Bool:
		if v.Bool() ***REMOVED***
			return 1, true
		***REMOVED***
		return 0, true
	***REMOVED***
	return 0, false
***REMOVED***

// numLess returns whether a < b.
// a and b must necessarily have the same kind.
func numLess(a, b reflect.Value) bool ***REMOVED***
	switch a.Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int()
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() < b.Uint()
	case reflect.Bool:
		return !a.Bool() && b.Bool()
	***REMOVED***
	panic("not a number")
***REMOVED***
