package goja

import (
	"io"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/dop251/goja/unistring"
)

const (
	__proto__ = "__proto__"
)

var (
	stringTrue         valueString = asciiString("true")
	stringFalse        valueString = asciiString("false")
	stringNull         valueString = asciiString("null")
	stringUndefined    valueString = asciiString("undefined")
	stringObjectC      valueString = asciiString("object")
	stringFunction     valueString = asciiString("function")
	stringBoolean      valueString = asciiString("boolean")
	stringString       valueString = asciiString("string")
	stringSymbol       valueString = asciiString("symbol")
	stringNumber       valueString = asciiString("number")
	stringNaN          valueString = asciiString("NaN")
	stringInfinity                 = asciiString("Infinity")
	stringPlusInfinity             = asciiString("+Infinity")
	stringNegInfinity              = asciiString("-Infinity")
	stringBound_       valueString = asciiString("bound ")
	stringEmpty        valueString = asciiString("")

	stringError          valueString = asciiString("Error")
	stringAggregateError valueString = asciiString("AggregateError")
	stringTypeError      valueString = asciiString("TypeError")
	stringReferenceError valueString = asciiString("ReferenceError")
	stringSyntaxError    valueString = asciiString("SyntaxError")
	stringRangeError     valueString = asciiString("RangeError")
	stringEvalError      valueString = asciiString("EvalError")
	stringURIError       valueString = asciiString("URIError")
	stringGoError        valueString = asciiString("GoError")

	stringObjectNull      valueString = asciiString("[object Null]")
	stringObjectObject    valueString = asciiString("[object Object]")
	stringObjectUndefined valueString = asciiString("[object Undefined]")
	stringInvalidDate     valueString = asciiString("Invalid Date")
)

type valueString interface ***REMOVED***
	Value
	charAt(int) rune
	length() int
	concat(valueString) valueString
	substring(start, end int) valueString
	compareTo(valueString) int
	reader(start int) io.RuneReader
	utf16Reader(start int) io.RuneReader
	utf16Runes() []rune
	index(valueString, int) int
	lastIndex(valueString, int) int
	toLower() valueString
	toUpper() valueString
	toTrimmedUTF8() string
***REMOVED***

type stringIterObject struct ***REMOVED***
	baseObject
	reader io.RuneReader
***REMOVED***

func isUTF16FirstSurrogate(r rune) bool ***REMOVED***
	return r >= 0xD800 && r <= 0xDBFF
***REMOVED***

func isUTF16SecondSurrogate(r rune) bool ***REMOVED***
	return r >= 0xDC00 && r <= 0xDFFF
***REMOVED***

func (si *stringIterObject) next() Value ***REMOVED***
	if si.reader == nil ***REMOVED***
		return si.val.runtime.createIterResultObject(_undefined, true)
	***REMOVED***
	r, _, err := si.reader.ReadRune()
	if err == io.EOF ***REMOVED***
		si.reader = nil
		return si.val.runtime.createIterResultObject(_undefined, true)
	***REMOVED***
	return si.val.runtime.createIterResultObject(stringFromRune(r), false)
***REMOVED***

func stringFromRune(r rune) valueString ***REMOVED***
	if r < utf8.RuneSelf ***REMOVED***
		var sb strings.Builder
		sb.Grow(1)
		sb.WriteByte(byte(r))
		return asciiString(sb.String())
	***REMOVED***
	var sb unicodeStringBuilder
	if r <= 0xFFFF ***REMOVED***
		sb.Grow(1)
	***REMOVED*** else ***REMOVED***
		sb.Grow(2)
	***REMOVED***
	sb.WriteRune(r)
	return sb.String()
***REMOVED***

func (r *Runtime) createStringIterator(s valueString) Value ***REMOVED***
	o := &Object***REMOVED***runtime: r***REMOVED***

	si := &stringIterObject***REMOVED***
		reader: &lenientUtf16Decoder***REMOVED***utf16Reader: s.utf16Reader(0)***REMOVED***,
	***REMOVED***
	si.class = classStringIterator
	si.val = o
	si.extensible = true
	o.self = si
	si.prototype = r.global.StringIteratorPrototype
	si.init()

	return o
***REMOVED***

type stringObject struct ***REMOVED***
	baseObject
	value      valueString
	length     int
	lengthProp valueProperty
***REMOVED***

func newStringValue(s string) valueString ***REMOVED***
	utf16Size := 0
	ascii := true
	for _, chr := range s ***REMOVED***
		utf16Size++
		if chr >= utf8.RuneSelf ***REMOVED***
			ascii = false
			if chr > 0xFFFF ***REMOVED***
				utf16Size++
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if ascii ***REMOVED***
		return asciiString(s)
	***REMOVED***
	buf := make([]uint16, utf16Size+1)
	buf[0] = unistring.BOM
	c := 1
	for _, chr := range s ***REMOVED***
		if chr <= 0xFFFF ***REMOVED***
			buf[c] = uint16(chr)
		***REMOVED*** else ***REMOVED***
			first, second := utf16.EncodeRune(chr)
			buf[c] = uint16(first)
			c++
			buf[c] = uint16(second)
		***REMOVED***
		c++
	***REMOVED***
	return unicodeString(buf)
***REMOVED***

func stringValueFromRaw(raw unistring.String) valueString ***REMOVED***
	if b := raw.AsUtf16(); b != nil ***REMOVED***
		return unicodeString(b)
	***REMOVED***
	return asciiString(raw)
***REMOVED***

func (s *stringObject) init() ***REMOVED***
	s.baseObject.init()
	s.setLength()
***REMOVED***

func (s *stringObject) setLength() ***REMOVED***
	if s.value != nil ***REMOVED***
		s.length = s.value.length()
	***REMOVED***
	s.lengthProp.value = intToValue(int64(s.length))
	s._put("length", &s.lengthProp)
***REMOVED***

func (s *stringObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	if i := strToGoIdx(name); i >= 0 && i < s.length ***REMOVED***
		return s._getIdx(i)
	***REMOVED***
	return s.baseObject.getStr(name, receiver)
***REMOVED***

func (s *stringObject) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	i := int(idx)
	if i >= 0 && i < s.length ***REMOVED***
		return s._getIdx(i)
	***REMOVED***
	return s.baseObject.getStr(idx.string(), receiver)
***REMOVED***

func (s *stringObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if i := strToGoIdx(name); i >= 0 && i < s.length ***REMOVED***
		val := s._getIdx(i)
		return &valueProperty***REMOVED***
			value:      val,
			enumerable: true,
		***REMOVED***
	***REMOVED***

	return s.baseObject.getOwnPropStr(name)
***REMOVED***

func (s *stringObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	i := int64(idx)
	if i >= 0 ***REMOVED***
		if i < int64(s.length) ***REMOVED***
			val := s._getIdx(int(i))
			return &valueProperty***REMOVED***
				value:      val,
				enumerable: true,
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	return s.baseObject.getOwnPropStr(idx.string())
***REMOVED***

func (s *stringObject) _getIdx(idx int) Value ***REMOVED***
	return s.value.substring(idx, idx+1)
***REMOVED***

func (s *stringObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if i := strToGoIdx(name); i >= 0 && i < s.length ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%d' of a String", i)
		return false
	***REMOVED***

	return s.baseObject.setOwnStr(name, val, throw)
***REMOVED***

func (s *stringObject) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	i := int64(idx)
	if i >= 0 && i < int64(s.length) ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%d' of a String", i)
		return false
	***REMOVED***

	return s.baseObject.setOwnStr(idx.string(), val, throw)
***REMOVED***

func (s *stringObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return s._setForeignStr(name, s.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

func (s *stringObject) setForeignIdx(idx valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return s._setForeignIdx(idx, s.getOwnPropIdx(idx), val, receiver, throw)
***REMOVED***

func (s *stringObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if i := strToGoIdx(name); i >= 0 && i < s.length ***REMOVED***
		_, ok := s._defineOwnProperty(name, &valueProperty***REMOVED***enumerable: true***REMOVED***, descr, throw)
		return ok
	***REMOVED***

	return s.baseObject.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (s *stringObject) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	i := int64(idx)
	if i >= 0 && i < int64(s.length) ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot redefine property: %d", i)
		return false
	***REMOVED***

	return s.baseObject.defineOwnPropertyStr(idx.string(), descr, throw)
***REMOVED***

type stringPropIter struct ***REMOVED***
	str         valueString // separate, because obj can be the singleton
	obj         *stringObject
	idx, length int
***REMOVED***

func (i *stringPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.length ***REMOVED***
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: asciiString(name), enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	return i.obj.baseObject.iterateStringKeys()()
***REMOVED***

func (s *stringObject) iterateStringKeys() iterNextFunc ***REMOVED***
	return (&stringPropIter***REMOVED***
		str:    s.value,
		obj:    s,
		length: s.length,
	***REMOVED***).next
***REMOVED***

func (s *stringObject) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	for i := 0; i < s.length; i++ ***REMOVED***
		accum = append(accum, asciiString(strconv.Itoa(i)))
	***REMOVED***

	return s.baseObject.stringKeys(all, accum)
***REMOVED***

func (s *stringObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if i := strToGoIdx(name); i >= 0 && i < s.length ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of a String", i)
		return false
	***REMOVED***

	return s.baseObject.deleteStr(name, throw)
***REMOVED***

func (s *stringObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	i := int64(idx)
	if i >= 0 && i < int64(s.length) ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of a String", i)
		return false
	***REMOVED***

	return s.baseObject.deleteStr(idx.string(), throw)
***REMOVED***

func (s *stringObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if i := strToGoIdx(name); i >= 0 && i < s.length ***REMOVED***
		return true
	***REMOVED***
	return s.baseObject.hasOwnPropertyStr(name)
***REMOVED***

func (s *stringObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	i := int64(idx)
	if i >= 0 && i < int64(s.length) ***REMOVED***
		return true
	***REMOVED***
	return s.baseObject.hasOwnPropertyStr(idx.string())
***REMOVED***
