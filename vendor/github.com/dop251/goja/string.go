package goja

import (
	"io"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
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
	stringNumber       valueString = asciiString("number")
	stringNaN          valueString = asciiString("NaN")
	stringInfinity                 = asciiString("Infinity")
	stringPlusInfinity             = asciiString("+Infinity")
	stringNegInfinity              = asciiString("-Infinity")
	stringEmpty        valueString = asciiString("")
	string__proto__    valueString = asciiString(__proto__)

	stringError          valueString = asciiString("Error")
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
	stringGlobalObject    valueString = asciiString("Global Object")
	stringInvalidDate     valueString = asciiString("Invalid Date")
)

type valueString interface ***REMOVED***
	Value
	charAt(int64) rune
	length() int64
	concat(valueString) valueString
	substring(start, end int64) valueString
	compareTo(valueString) int
	reader(start int) io.RuneReader
	index(valueString, int64) int64
	lastIndex(valueString, int64) int64
	toLower() valueString
	toUpper() valueString
	toTrimmedUTF8() string
***REMOVED***

type stringObject struct ***REMOVED***
	baseObject
	value      valueString
	length     int64
	lengthProp valueProperty
***REMOVED***

func newUnicodeString(s string) valueString ***REMOVED***
	return unicodeString(utf16.Encode([]rune(s)))
***REMOVED***

func newStringValue(s string) valueString ***REMOVED***
	for _, chr := range s ***REMOVED***
		if chr >= utf8.RuneSelf ***REMOVED***
			return newUnicodeString(s)
		***REMOVED***
	***REMOVED***
	return asciiString(s)
***REMOVED***

func (s *stringObject) init() ***REMOVED***
	s.baseObject.init()
	s.setLength()
***REMOVED***

func (s *stringObject) setLength() ***REMOVED***
	if s.value != nil ***REMOVED***
		s.length = s.value.length()
	***REMOVED***
	s.lengthProp.value = intToValue(s.length)
	s._put("length", &s.lengthProp)
***REMOVED***

func (s *stringObject) get(n Value) Value ***REMOVED***
	if idx := toIdx(n); idx >= 0 && idx < s.length ***REMOVED***
		return s.getIdx(idx)
	***REMOVED***
	return s.baseObject.get(n)
***REMOVED***

func (s *stringObject) getStr(name string) Value ***REMOVED***
	if i := strToIdx(name); i >= 0 && i < s.length ***REMOVED***
		return s.getIdx(i)
	***REMOVED***
	return s.baseObject.getStr(name)
***REMOVED***

func (s *stringObject) getPropStr(name string) Value ***REMOVED***
	if i := strToIdx(name); i >= 0 && i < s.length ***REMOVED***
		return s.getIdx(i)
	***REMOVED***
	return s.baseObject.getPropStr(name)
***REMOVED***

func (s *stringObject) getProp(n Value) Value ***REMOVED***
	if i := toIdx(n); i >= 0 && i < s.length ***REMOVED***
		return s.getIdx(i)
	***REMOVED***
	return s.baseObject.getProp(n)
***REMOVED***

func (s *stringObject) getOwnProp(name string) Value ***REMOVED***
	if i := strToIdx(name); i >= 0 && i < s.length ***REMOVED***
		val := s.getIdx(i)
		return &valueProperty***REMOVED***
			value:      val,
			enumerable: true,
		***REMOVED***
	***REMOVED***

	return s.baseObject.getOwnProp(name)
***REMOVED***

func (s *stringObject) getIdx(idx int64) Value ***REMOVED***
	return s.value.substring(idx, idx+1)
***REMOVED***

func (s *stringObject) put(n Value, val Value, throw bool) ***REMOVED***
	if i := toIdx(n); i >= 0 && i < s.length ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%d' of a String", i)
		return
	***REMOVED***

	s.baseObject.put(n, val, throw)
***REMOVED***

func (s *stringObject) putStr(name string, val Value, throw bool) ***REMOVED***
	if i := strToIdx(name); i >= 0 && i < s.length ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%d' of a String", i)
		return
	***REMOVED***

	s.baseObject.putStr(name, val, throw)
***REMOVED***

func (s *stringObject) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool ***REMOVED***
	if i := toIdx(n); i >= 0 && i < s.length ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot redefine property: %d", i)
		return false
	***REMOVED***

	return s.baseObject.defineOwnProperty(n, descr, throw)
***REMOVED***

type stringPropIter struct ***REMOVED***
	str         valueString // separate, because obj can be the singleton
	obj         *stringObject
	idx, length int64
	recursive   bool
***REMOVED***

func (i *stringPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.length ***REMOVED***
		name := strconv.FormatInt(i.idx, 10)
		i.idx++
		return propIterItem***REMOVED***name: name, enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	return i.obj.baseObject._enumerate(i.recursive)()
***REMOVED***

func (s *stringObject) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	return (&stringPropIter***REMOVED***
		str:       s.value,
		obj:       s,
		length:    s.length,
		recursive: recursive,
	***REMOVED***).next
***REMOVED***

func (s *stringObject) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: s._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next
***REMOVED***

func (s *stringObject) deleteStr(name string, throw bool) bool ***REMOVED***
	if i := strToIdx(name); i >= 0 && i < s.length ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of a String", i)
		return false
	***REMOVED***

	return s.baseObject.deleteStr(name, throw)
***REMOVED***

func (s *stringObject) delete(n Value, throw bool) bool ***REMOVED***
	if i := toIdx(n); i >= 0 && i < s.length ***REMOVED***
		s.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of a String", i)
		return false
	***REMOVED***

	return s.baseObject.delete(n, throw)
***REMOVED***

func (s *stringObject) hasOwnProperty(n Value) bool ***REMOVED***
	if i := toIdx(n); i >= 0 && i < s.length ***REMOVED***
		return true
	***REMOVED***
	return s.baseObject.hasOwnProperty(n)
***REMOVED***

func (s *stringObject) hasOwnPropertyStr(name string) bool ***REMOVED***
	if i := strToIdx(name); i >= 0 && i < s.length ***REMOVED***
		return true
	***REMOVED***
	return s.baseObject.hasOwnPropertyStr(name)
***REMOVED***
