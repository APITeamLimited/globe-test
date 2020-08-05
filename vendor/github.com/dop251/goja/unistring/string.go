// Package unistring contains an implementation of a hybrid ASCII/UTF-16 string.
// For ASCII strings the underlying representation is equivalent to a normal Go string.
// For unicode strings the underlying representation is UTF-16 as []uint16 with 0th element set to 0xFEFF.
// unicode.String allows representing malformed UTF-16 values (e.g. stand-alone parts of surrogate pairs)
// which cannot be represented in UTF-8.
// At the same time it is possible to use unicode.String as property keys just as efficiently as simple strings,
// (the leading 0xFEFF ensures there is no clash with ASCII string), and it is possible to convert it
// to valueString without extra allocations.
package unistring

import (
	"reflect"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

const (
	BOM = 0xFEFF
)

type String string

func NewFromString(s string) String ***REMOVED***
	ascii := true
	size := 0
	for _, c := range s ***REMOVED***
		if c >= utf8.RuneSelf ***REMOVED***
			ascii = false
			if c > 0xFFFF ***REMOVED***
				size++
			***REMOVED***
		***REMOVED***
		size++
	***REMOVED***
	if ascii ***REMOVED***
		return String(s)
	***REMOVED***
	b := make([]uint16, size+1)
	b[0] = BOM
	i := 1
	for _, c := range s ***REMOVED***
		if c <= 0xFFFF ***REMOVED***
			b[i] = uint16(c)
		***REMOVED*** else ***REMOVED***
			first, second := utf16.EncodeRune(c)
			b[i] = uint16(first)
			i++
			b[i] = uint16(second)
		***REMOVED***
		i++
	***REMOVED***
	return FromUtf16(b)
***REMOVED***

func NewFromRunes(s []rune) String ***REMOVED***
	ascii := true
	size := 0
	for _, c := range s ***REMOVED***
		if c >= utf8.RuneSelf ***REMOVED***
			ascii = false
			if c > 0xFFFF ***REMOVED***
				size++
			***REMOVED***
		***REMOVED***
		size++
	***REMOVED***
	if ascii ***REMOVED***
		return String(s)
	***REMOVED***
	b := make([]uint16, size+1)
	b[0] = BOM
	i := 1
	for _, c := range s ***REMOVED***
		if c <= 0xFFFF ***REMOVED***
			b[i] = uint16(c)
		***REMOVED*** else ***REMOVED***
			first, second := utf16.EncodeRune(c)
			b[i] = uint16(first)
			i++
			b[i] = uint16(second)
		***REMOVED***
		i++
	***REMOVED***
	return FromUtf16(b)
***REMOVED***

func FromUtf16(b []uint16) String ***REMOVED***
	var str string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	hdr.Data = uintptr(unsafe.Pointer(&b[0]))
	hdr.Len = len(b) * 2

	return String(str)
***REMOVED***

func (s String) String() string ***REMOVED***
	if b := s.AsUtf16(); b != nil ***REMOVED***
		return string(utf16.Decode(b[1:]))
	***REMOVED***

	return string(s)
***REMOVED***

func (s String) AsUtf16() []uint16 ***REMOVED***
	if len(s) < 4 || len(s)&1 != 0 ***REMOVED***
		return nil
	***REMOVED***
	l := len(s) / 2
	raw := string(s)
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&raw))
	a := *(*[]uint16)(unsafe.Pointer(&reflect.SliceHeader***REMOVED***
		Data: hdr.Data,
		Len:  l,
		Cap:  l,
	***REMOVED***))
	if a[0] == BOM ***REMOVED***
		return a
	***REMOVED***

	return nil
***REMOVED***
