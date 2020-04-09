package goja

import (
	"errors"
	"fmt"
	"github.com/dop251/goja/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"math"
	"reflect"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

type unicodeString []uint16

type unicodeRuneReader struct ***REMOVED***
	s   unicodeString
	pos int
***REMOVED***

type runeReaderReplace struct ***REMOVED***
	wrapped io.RuneReader
***REMOVED***

var (
	InvalidRuneError = errors.New("Invalid rune")
)

func (rr runeReaderReplace) ReadRune() (r rune, size int, err error) ***REMOVED***
	r, size, err = rr.wrapped.ReadRune()
	if err == InvalidRuneError ***REMOVED***
		err = nil
		r = utf8.RuneError
	***REMOVED***
	return
***REMOVED***

func (rr *unicodeRuneReader) ReadRune() (r rune, size int, err error) ***REMOVED***
	if rr.pos < len(rr.s) ***REMOVED***
		r = rune(rr.s[rr.pos])
		if r != utf8.RuneError ***REMOVED***
			if utf16.IsSurrogate(r) ***REMOVED***
				if rr.pos+1 < len(rr.s) ***REMOVED***
					r1 := utf16.DecodeRune(r, rune(rr.s[rr.pos+1]))
					size++
					rr.pos++
					if r1 == utf8.RuneError ***REMOVED***
						err = InvalidRuneError
					***REMOVED*** else ***REMOVED***
						r = r1
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					err = InvalidRuneError
				***REMOVED***
			***REMOVED***
		***REMOVED***
		size++
		rr.pos++
	***REMOVED*** else ***REMOVED***
		err = io.EOF
	***REMOVED***
	return
***REMOVED***

func (s unicodeString) reader(start int) io.RuneReader ***REMOVED***
	return &unicodeRuneReader***REMOVED***
		s: s[start:],
	***REMOVED***
***REMOVED***

func (s unicodeString) ToInteger() int64 ***REMOVED***
	return 0
***REMOVED***

func (s unicodeString) ToString() valueString ***REMOVED***
	return s
***REMOVED***

func (s unicodeString) ToFloat() float64 ***REMOVED***
	return math.NaN()
***REMOVED***

func (s unicodeString) ToBoolean() bool ***REMOVED***
	return len(s) > 0
***REMOVED***

func (s unicodeString) toTrimmedUTF8() string ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return strings.Trim(s.String(), parser.WhitespaceChars)
***REMOVED***

func (s unicodeString) ToNumber() Value ***REMOVED***
	return asciiString(s.toTrimmedUTF8()).ToNumber()
***REMOVED***

func (s unicodeString) ToObject(r *Runtime) *Object ***REMOVED***
	return r._newString(s)
***REMOVED***

func (s unicodeString) equals(other unicodeString) bool ***REMOVED***
	if len(s) != len(other) ***REMOVED***
		return false
	***REMOVED***
	for i, r := range s ***REMOVED***
		if r != other[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (s unicodeString) SameAs(other Value) bool ***REMOVED***
	if otherStr, ok := other.(unicodeString); ok ***REMOVED***
		return s.equals(otherStr)
	***REMOVED***

	return false
***REMOVED***

func (s unicodeString) Equals(other Value) bool ***REMOVED***
	if s.SameAs(other) ***REMOVED***
		return true
	***REMOVED***

	if _, ok := other.assertInt(); ok ***REMOVED***
		return false
	***REMOVED***

	if _, ok := other.assertFloat(); ok ***REMOVED***
		return false
	***REMOVED***

	if _, ok := other.(valueBool); ok ***REMOVED***
		return false
	***REMOVED***

	if o, ok := other.(*Object); ok ***REMOVED***
		return s.Equals(o.self.toPrimitive())
	***REMOVED***
	return false
***REMOVED***

func (s unicodeString) StrictEquals(other Value) bool ***REMOVED***
	return s.SameAs(other)
***REMOVED***

func (s unicodeString) assertInt() (int64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (s unicodeString) assertFloat() (float64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (s unicodeString) assertString() (valueString, bool) ***REMOVED***
	return s, true
***REMOVED***

func (s unicodeString) baseObject(r *Runtime) *Object ***REMOVED***
	ss := r.stringSingleton
	ss.value = s
	ss.setLength()
	return ss.val
***REMOVED***

func (s unicodeString) charAt(idx int64) rune ***REMOVED***
	return rune(s[idx])
***REMOVED***

func (s unicodeString) length() int64 ***REMOVED***
	return int64(len(s))
***REMOVED***

func (s unicodeString) concat(other valueString) valueString ***REMOVED***
	switch other := other.(type) ***REMOVED***
	case unicodeString:
		return unicodeString(append(s, other...))
	case asciiString:
		b := make([]uint16, len(s)+len(other))
		copy(b, s)
		b1 := b[len(s):]
		for i := 0; i < len(other); i++ ***REMOVED***
			b1[i] = uint16(other[i])
		***REMOVED***
		return unicodeString(b)
	default:
		panic(fmt.Errorf("Unknown string type: %T", other))
	***REMOVED***
***REMOVED***

func (s unicodeString) substring(start, end int64) valueString ***REMOVED***
	ss := s[start:end]
	for _, c := range ss ***REMOVED***
		if c >= utf8.RuneSelf ***REMOVED***
			return unicodeString(ss)
		***REMOVED***
	***REMOVED***
	as := make([]byte, end-start)
	for i, c := range ss ***REMOVED***
		as[i] = byte(c)
	***REMOVED***
	return asciiString(as)
***REMOVED***

func (s unicodeString) String() string ***REMOVED***
	return string(utf16.Decode(s))
***REMOVED***

func (s unicodeString) compareTo(other valueString) int ***REMOVED***
	return strings.Compare(s.String(), other.String())
***REMOVED***

func (s unicodeString) index(substr valueString, start int64) int64 ***REMOVED***
	var ss []uint16
	switch substr := substr.(type) ***REMOVED***
	case unicodeString:
		ss = substr
	case asciiString:
		ss = make([]uint16, len(substr))
		for i := 0; i < len(substr); i++ ***REMOVED***
			ss[i] = uint16(substr[i])
		***REMOVED***
	default:
		panic(fmt.Errorf("Unknown string type: %T", substr))
	***REMOVED***

	// TODO: optimise
	end := int64(len(s) - len(ss))
	for start <= end ***REMOVED***
		for i := int64(0); i < int64(len(ss)); i++ ***REMOVED***
			if s[start+i] != ss[i] ***REMOVED***
				goto nomatch
			***REMOVED***
		***REMOVED***

		return start
	nomatch:
		start++
	***REMOVED***
	return -1
***REMOVED***

func (s unicodeString) lastIndex(substr valueString, start int64) int64 ***REMOVED***
	var ss []uint16
	switch substr := substr.(type) ***REMOVED***
	case unicodeString:
		ss = substr
	case asciiString:
		ss = make([]uint16, len(substr))
		for i := 0; i < len(substr); i++ ***REMOVED***
			ss[i] = uint16(substr[i])
		***REMOVED***
	default:
		panic(fmt.Errorf("Unknown string type: %T", substr))
	***REMOVED***

	if maxStart := int64(len(s) - len(ss)); start > maxStart ***REMOVED***
		start = maxStart
	***REMOVED***
	// TODO: optimise
	for start >= 0 ***REMOVED***
		for i := int64(0); i < int64(len(ss)); i++ ***REMOVED***
			if s[start+i] != ss[i] ***REMOVED***
				goto nomatch
			***REMOVED***
		***REMOVED***

		return start
	nomatch:
		start--
	***REMOVED***
	return -1
***REMOVED***

func (s unicodeString) toLower() valueString ***REMOVED***
	caser := cases.Lower(language.Und)
	r := []rune(caser.String(s.String()))
	// Workaround
	ascii := true
	for i := 0; i < len(r)-1; i++ ***REMOVED***
		if (i == 0 || r[i-1] != 0x3b1) && r[i] == 0x345 && r[i+1] == 0x3c2 ***REMOVED***
			i++
			r[i] = 0x3c3
		***REMOVED***
		if r[i] >= utf8.RuneSelf ***REMOVED***
			ascii = false
		***REMOVED***
	***REMOVED***
	if ascii ***REMOVED***
		ascii = r[len(r)-1] < utf8.RuneSelf
	***REMOVED***
	if ascii ***REMOVED***
		return asciiString(r)
	***REMOVED***
	return unicodeString(utf16.Encode(r))
***REMOVED***

func (s unicodeString) toUpper() valueString ***REMOVED***
	caser := cases.Upper(language.Und)
	return newStringValue(caser.String(s.String()))
***REMOVED***

func (s unicodeString) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return s.String()
***REMOVED***

func (s unicodeString) ExportType() reflect.Type ***REMOVED***
	return reflectTypeString
***REMOVED***
