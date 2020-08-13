package goja

import (
	"errors"
	"fmt"
	"hash/maphash"
	"io"
	"math"
	"reflect"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/dop251/goja/parser"
	"github.com/dop251/goja/unistring"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type unicodeString []uint16

type unicodeRuneReader struct ***REMOVED***
	s   unicodeString
	pos int
***REMOVED***

type utf16RuneReader struct ***REMOVED***
	s   unicodeString
	pos int
***REMOVED***

// passes through invalid surrogate pairs
type lenientUtf16Decoder struct ***REMOVED***
	utf16Reader io.RuneReader
	prev        rune
	prevSet     bool
***REMOVED***

type valueStringBuilder struct ***REMOVED***
	asciiBuilder   strings.Builder
	unicodeBuilder unicodeStringBuilder
***REMOVED***

type unicodeStringBuilder struct ***REMOVED***
	buf     []uint16
	unicode bool
***REMOVED***

var (
	InvalidRuneError = errors.New("invalid rune")
)

func (rr *utf16RuneReader) ReadRune() (r rune, size int, err error) ***REMOVED***
	if rr.pos < len(rr.s) ***REMOVED***
		r = rune(rr.s[rr.pos])
		size++
		rr.pos++
		return
	***REMOVED***
	err = io.EOF
	return
***REMOVED***

func (rr *lenientUtf16Decoder) ReadRune() (r rune, size int, err error) ***REMOVED***
	if rr.prevSet ***REMOVED***
		r = rr.prev
		size = 1
		rr.prevSet = false
	***REMOVED*** else ***REMOVED***
		r, size, err = rr.utf16Reader.ReadRune()
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if isUTF16FirstSurrogate(r) ***REMOVED***
		second, _, err1 := rr.utf16Reader.ReadRune()
		if err1 != nil ***REMOVED***
			if err1 != io.EOF ***REMOVED***
				err = err1
			***REMOVED***
			return
		***REMOVED***
		if isUTF16SecondSurrogate(second) ***REMOVED***
			r = utf16.DecodeRune(r, second)
			size++
		***REMOVED*** else ***REMOVED***
			rr.prev = second
			rr.prevSet = true
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (rr *unicodeRuneReader) ReadRune() (r rune, size int, err error) ***REMOVED***
	if rr.pos < len(rr.s) ***REMOVED***
		r = rune(rr.s[rr.pos])
		size++
		rr.pos++
		if isUTF16FirstSurrogate(r) ***REMOVED***
			if rr.pos < len(rr.s) ***REMOVED***
				second := rune(rr.s[rr.pos])
				if isUTF16SecondSurrogate(second) ***REMOVED***
					r = utf16.DecodeRune(r, second)
					size++
					rr.pos++
				***REMOVED*** else ***REMOVED***
					err = InvalidRuneError
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				err = InvalidRuneError
			***REMOVED***
		***REMOVED*** else if isUTF16SecondSurrogate(r) ***REMOVED***
			err = InvalidRuneError
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		err = io.EOF
	***REMOVED***
	return
***REMOVED***

func (b *unicodeStringBuilder) grow(n int) ***REMOVED***
	if cap(b.buf)-len(b.buf) < n ***REMOVED***
		buf := make([]uint16, len(b.buf), 2*cap(b.buf)+n)
		copy(buf, b.buf)
		b.buf = buf
	***REMOVED***
***REMOVED***

func (b *unicodeStringBuilder) Grow(n int) ***REMOVED***
	b.grow(n + 1)
***REMOVED***

func (b *unicodeStringBuilder) ensureStarted(initialSize int) ***REMOVED***
	b.grow(len(b.buf) + initialSize + 1)
	if len(b.buf) == 0 ***REMOVED***
		b.buf = append(b.buf, unistring.BOM)
	***REMOVED***
***REMOVED***

func (b *unicodeStringBuilder) WriteString(s valueString) ***REMOVED***
	b.ensureStarted(s.length())
	switch s := s.(type) ***REMOVED***
	case unicodeString:
		b.buf = append(b.buf, s[1:]...)
		b.unicode = true
	case asciiString:
		for i := 0; i < len(s); i++ ***REMOVED***
			b.buf = append(b.buf, uint16(s[i]))
		***REMOVED***
	default:
		panic(fmt.Errorf("unsupported string type: %T", s))
	***REMOVED***
***REMOVED***

func (b *unicodeStringBuilder) String() valueString ***REMOVED***
	if b.unicode ***REMOVED***
		return unicodeString(b.buf)
	***REMOVED***
	if len(b.buf) == 0 ***REMOVED***
		return stringEmpty
	***REMOVED***
	buf := make([]byte, 0, len(b.buf)-1)
	for _, c := range b.buf[1:] ***REMOVED***
		buf = append(buf, byte(c))
	***REMOVED***
	return asciiString(buf)
***REMOVED***

func (b *unicodeStringBuilder) WriteRune(r rune) ***REMOVED***
	if r <= 0xFFFF ***REMOVED***
		b.ensureStarted(1)
		b.buf = append(b.buf, uint16(r))
		if !b.unicode && r >= utf8.RuneSelf ***REMOVED***
			b.unicode = true
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		b.ensureStarted(2)
		first, second := utf16.EncodeRune(r)
		b.buf = append(b.buf, uint16(first), uint16(second))
		b.unicode = true
	***REMOVED***
***REMOVED***

func (b *unicodeStringBuilder) writeASCIIString(bytes string) ***REMOVED***
	b.ensureStarted(len(bytes))
	for _, c := range bytes ***REMOVED***
		b.buf = append(b.buf, uint16(c))
	***REMOVED***
***REMOVED***

func (b *valueStringBuilder) ascii() bool ***REMOVED***
	return len(b.unicodeBuilder.buf) == 0
***REMOVED***

func (b *valueStringBuilder) WriteString(s valueString) ***REMOVED***
	if ascii, ok := s.(asciiString); ok ***REMOVED***
		if b.ascii() ***REMOVED***
			b.asciiBuilder.WriteString(string(ascii))
		***REMOVED*** else ***REMOVED***
			b.unicodeBuilder.writeASCIIString(string(ascii))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		b.switchToUnicode(s.length())
		b.unicodeBuilder.WriteString(s)
	***REMOVED***
***REMOVED***

func (b *valueStringBuilder) WriteRune(r rune) ***REMOVED***
	if r < utf8.RuneSelf ***REMOVED***
		if b.ascii() ***REMOVED***
			b.asciiBuilder.WriteByte(byte(r))
		***REMOVED*** else ***REMOVED***
			b.unicodeBuilder.WriteRune(r)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var extraLen int
		if r <= 0xFFFF ***REMOVED***
			extraLen = 1
		***REMOVED*** else ***REMOVED***
			extraLen = 2
		***REMOVED***
		b.switchToUnicode(extraLen)
		b.unicodeBuilder.WriteRune(r)
	***REMOVED***
***REMOVED***

func (b *valueStringBuilder) String() valueString ***REMOVED***
	if b.ascii() ***REMOVED***
		return asciiString(b.asciiBuilder.String())
	***REMOVED***
	return b.unicodeBuilder.String()
***REMOVED***

func (b *valueStringBuilder) Grow(n int) ***REMOVED***
	if b.ascii() ***REMOVED***
		b.asciiBuilder.Grow(n)
	***REMOVED*** else ***REMOVED***
		b.unicodeBuilder.Grow(n)
	***REMOVED***
***REMOVED***

func (b *valueStringBuilder) switchToUnicode(extraLen int) ***REMOVED***
	if b.ascii() ***REMOVED***
		b.unicodeBuilder.ensureStarted(b.asciiBuilder.Len() + extraLen)
		b.unicodeBuilder.writeASCIIString(b.asciiBuilder.String())
		b.asciiBuilder.Reset()
	***REMOVED***
***REMOVED***

func (b *valueStringBuilder) WriteSubstring(source valueString, start int, end int) ***REMOVED***
	if ascii, ok := source.(asciiString); ok ***REMOVED***
		if b.ascii() ***REMOVED***
			b.asciiBuilder.WriteString(string(ascii[start:end]))
			return
		***REMOVED***
	***REMOVED***
	us := source.(unicodeString)
	if b.ascii() ***REMOVED***
		uc := false
		for i := start; i < end; i++ ***REMOVED***
			if us.charAt(i) >= utf8.RuneSelf ***REMOVED***
				uc = true
				break
			***REMOVED***
		***REMOVED***
		if uc ***REMOVED***
			b.switchToUnicode(end - start + 1)
		***REMOVED*** else ***REMOVED***
			b.asciiBuilder.Grow(end - start + 1)
			for i := start; i < end; i++ ***REMOVED***
				b.asciiBuilder.WriteByte(byte(us.charAt(i)))
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	b.unicodeBuilder.buf = append(b.unicodeBuilder.buf, us[start+1:end+1]...)
	b.unicodeBuilder.unicode = true
***REMOVED***

func (s unicodeString) reader(start int) io.RuneReader ***REMOVED***
	return &unicodeRuneReader***REMOVED***
		s: s[start+1:],
	***REMOVED***
***REMOVED***

func (s unicodeString) utf16Reader(start int) io.RuneReader ***REMOVED***
	return &utf16RuneReader***REMOVED***
		s: s[start+1:],
	***REMOVED***
***REMOVED***

func (s unicodeString) runes() []rune ***REMOVED***
	return utf16.Decode(s[1:])
***REMOVED***

func (s unicodeString) utf16Runes() []rune ***REMOVED***
	runes := make([]rune, len(s)-1)
	for i, ch := range s[1:] ***REMOVED***
		runes[i] = rune(ch)
	***REMOVED***
	return runes
***REMOVED***

func (s unicodeString) ToInteger() int64 ***REMOVED***
	return 0
***REMOVED***

func (s unicodeString) toString() valueString ***REMOVED***
	return s
***REMOVED***

func (s unicodeString) ToString() Value ***REMOVED***
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
	return r._newString(s, r.global.StringPrototype)
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

	if o, ok := other.(*Object); ok ***REMOVED***
		return s.Equals(o.toPrimitive())
	***REMOVED***
	return false
***REMOVED***

func (s unicodeString) StrictEquals(other Value) bool ***REMOVED***
	return s.SameAs(other)
***REMOVED***

func (s unicodeString) baseObject(r *Runtime) *Object ***REMOVED***
	ss := r.stringSingleton
	ss.value = s
	ss.setLength()
	return ss.val
***REMOVED***

func (s unicodeString) charAt(idx int) rune ***REMOVED***
	return rune(s[idx+1])
***REMOVED***

func (s unicodeString) length() int ***REMOVED***
	return len(s) - 1
***REMOVED***

func (s unicodeString) concat(other valueString) valueString ***REMOVED***
	switch other := other.(type) ***REMOVED***
	case unicodeString:
		b := make(unicodeString, len(s)+len(other)-1)
		copy(b, s)
		copy(b[len(s):], other[1:])
		return b
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

func (s unicodeString) substring(start, end int) valueString ***REMOVED***
	ss := s[start+1 : end+1]
	for _, c := range ss ***REMOVED***
		if c >= utf8.RuneSelf ***REMOVED***
			b := make(unicodeString, end-start+1)
			b[0] = unistring.BOM
			copy(b[1:], ss)
			return b
		***REMOVED***
	***REMOVED***
	as := make([]byte, end-start)
	for i, c := range ss ***REMOVED***
		as[i] = byte(c)
	***REMOVED***
	return asciiString(as)
***REMOVED***

func (s unicodeString) String() string ***REMOVED***
	return string(utf16.Decode(s[1:]))
***REMOVED***

func (s unicodeString) compareTo(other valueString) int ***REMOVED***
	// TODO handle invalid UTF-16
	return strings.Compare(s.String(), other.String())
***REMOVED***

func (s unicodeString) index(substr valueString, start int) int ***REMOVED***
	var ss []uint16
	switch substr := substr.(type) ***REMOVED***
	case unicodeString:
		ss = substr[1:]
	case asciiString:
		ss = make([]uint16, len(substr))
		for i := 0; i < len(substr); i++ ***REMOVED***
			ss[i] = uint16(substr[i])
		***REMOVED***
	default:
		panic(fmt.Errorf("unknown string type: %T", substr))
	***REMOVED***
	s1 := s[1:]
	// TODO: optimise
	end := len(s1) - len(ss)
	for start <= end ***REMOVED***
		for i := 0; i < len(ss); i++ ***REMOVED***
			if s1[start+i] != ss[i] ***REMOVED***
				goto nomatch
			***REMOVED***
		***REMOVED***

		return start
	nomatch:
		start++
	***REMOVED***
	return -1
***REMOVED***

func (s unicodeString) lastIndex(substr valueString, start int) int ***REMOVED***
	var ss []uint16
	switch substr := substr.(type) ***REMOVED***
	case unicodeString:
		ss = substr[1:]
	case asciiString:
		ss = make([]uint16, len(substr))
		for i := 0; i < len(substr); i++ ***REMOVED***
			ss[i] = uint16(substr[i])
		***REMOVED***
	default:
		panic(fmt.Errorf("Unknown string type: %T", substr))
	***REMOVED***

	s1 := s[1:]
	if maxStart := len(s1) - len(ss); start > maxStart ***REMOVED***
		start = maxStart
	***REMOVED***
	// TODO: optimise
	for start >= 0 ***REMOVED***
		for i := 0; i < len(ss); i++ ***REMOVED***
			if s1[start+i] != ss[i] ***REMOVED***
				goto nomatch
			***REMOVED***
		***REMOVED***

		return start
	nomatch:
		start--
	***REMOVED***
	return -1
***REMOVED***

func unicodeStringFromRunes(r []rune) unicodeString ***REMOVED***
	return unistring.NewFromRunes(r).AsUtf16()
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
	return unicodeStringFromRunes(r)
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

func (s unicodeString) hash(hash *maphash.Hash) uint64 ***REMOVED***
	_, _ = hash.WriteString(string(unistring.FromUtf16(s)))
	h := hash.Sum64()
	hash.Reset()
	return h
***REMOVED***

func (s unicodeString) string() unistring.String ***REMOVED***
	return unistring.FromUtf16(s)
***REMOVED***
