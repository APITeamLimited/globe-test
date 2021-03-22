package goja

import (
	"github.com/dop251/goja/unistring"
	"math"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/dop251/goja/parser"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

func (r *Runtime) collator() *collate.Collator ***REMOVED***
	collator := r._collator
	if collator == nil ***REMOVED***
		collator = collate.New(language.Und)
		r._collator = collator
	***REMOVED***
	return collator
***REMOVED***

func toString(arg Value) valueString ***REMOVED***
	if s, ok := arg.(valueString); ok ***REMOVED***
		return s
	***REMOVED***
	if s, ok := arg.(*Symbol); ok ***REMOVED***
		return s.descriptiveString()
	***REMOVED***
	return arg.toString()
***REMOVED***

func (r *Runtime) builtin_String(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) > 0 ***REMOVED***
		return toString(call.Arguments[0])
	***REMOVED*** else ***REMOVED***
		return stringEmpty
	***REMOVED***
***REMOVED***

func (r *Runtime) _newString(s valueString, proto *Object) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	o := &stringObject***REMOVED******REMOVED***
	o.class = classString
	o.val = v
	o.extensible = true
	v.self = o
	o.prototype = proto
	if s != nil ***REMOVED***
		o.value = s
	***REMOVED***
	o.init()
	return v
***REMOVED***

func (r *Runtime) builtin_newString(args []Value, proto *Object) *Object ***REMOVED***
	var s valueString
	if len(args) > 0 ***REMOVED***
		s = args[0].toString()
	***REMOVED*** else ***REMOVED***
		s = stringEmpty
	***REMOVED***
	return r._newString(s, proto)
***REMOVED***

func (r *Runtime) stringproto_toStringValueOf(this Value, funcName string) Value ***REMOVED***
	if str, ok := this.(valueString); ok ***REMOVED***
		return str
	***REMOVED***
	if obj, ok := this.(*Object); ok ***REMOVED***
		if strObj, ok := obj.self.(*stringObject); ok ***REMOVED***
			return strObj.value
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "String.prototype.%s is called on incompatible receiver", funcName)
	return nil
***REMOVED***

func (r *Runtime) stringproto_toString(call FunctionCall) Value ***REMOVED***
	return r.stringproto_toStringValueOf(call.This, "toString")
***REMOVED***

func (r *Runtime) stringproto_valueOf(call FunctionCall) Value ***REMOVED***
	return r.stringproto_toStringValueOf(call.This, "valueOf")
***REMOVED***

func (r *Runtime) stringproto_iterator(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	return r.createStringIterator(call.This.toString())
***REMOVED***

func (r *Runtime) string_fromcharcode(call FunctionCall) Value ***REMOVED***
	b := make([]byte, len(call.Arguments))
	for i, arg := range call.Arguments ***REMOVED***
		chr := toUint16(arg)
		if chr >= utf8.RuneSelf ***REMOVED***
			bb := make([]uint16, len(call.Arguments)+1)
			bb[0] = unistring.BOM
			bb1 := bb[1:]
			for j := 0; j < i; j++ ***REMOVED***
				bb1[j] = uint16(b[j])
			***REMOVED***
			bb1[i] = chr
			i++
			for j, arg := range call.Arguments[i:] ***REMOVED***
				bb1[i+j] = toUint16(arg)
			***REMOVED***
			return unicodeString(bb)
		***REMOVED***
		b[i] = byte(chr)
	***REMOVED***

	return asciiString(b)
***REMOVED***

func (r *Runtime) string_fromcodepoint(call FunctionCall) Value ***REMOVED***
	var sb valueStringBuilder
	for _, arg := range call.Arguments ***REMOVED***
		num := arg.ToNumber()
		var c rune
		if numInt, ok := num.(valueInt); ok ***REMOVED***
			if numInt < 0 || numInt > utf8.MaxRune ***REMOVED***
				panic(r.newError(r.global.RangeError, "Invalid code point %d", numInt))
			***REMOVED***
			c = rune(numInt)
		***REMOVED*** else ***REMOVED***
			panic(r.newError(r.global.RangeError, "Invalid code point %s", num))
		***REMOVED***
		sb.WriteRune(c)
	***REMOVED***
	return sb.String()
***REMOVED***

func (r *Runtime) string_raw(call FunctionCall) Value ***REMOVED***
	cooked := call.Argument(0).ToObject(r)
	raw := nilSafe(cooked.self.getStr("raw", nil)).ToObject(r)
	literalSegments := toLength(raw.self.getStr("length", nil))
	if literalSegments <= 0 ***REMOVED***
		return stringEmpty
	***REMOVED***
	var stringElements valueStringBuilder
	nextIndex := int64(0)
	numberOfSubstitutions := int64(len(call.Arguments) - 1)
	for ***REMOVED***
		nextSeg := nilSafe(raw.self.getIdx(valueInt(nextIndex), nil)).toString()
		stringElements.WriteString(nextSeg)
		if nextIndex+1 == literalSegments ***REMOVED***
			return stringElements.String()
		***REMOVED***
		if nextIndex < numberOfSubstitutions ***REMOVED***
			stringElements.WriteString(nilSafe(call.Arguments[nextIndex+1]).toString())
		***REMOVED***
		nextIndex++
	***REMOVED***
***REMOVED***

func (r *Runtime) stringproto_charAt(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	pos := call.Argument(0).ToInteger()
	if pos < 0 || pos >= int64(s.length()) ***REMOVED***
		return stringEmpty
	***REMOVED***
	return newStringValue(string(s.charAt(toIntStrict(pos))))
***REMOVED***

func (r *Runtime) stringproto_charCodeAt(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	pos := call.Argument(0).ToInteger()
	if pos < 0 || pos >= int64(s.length()) ***REMOVED***
		return _NaN
	***REMOVED***
	return intToValue(int64(s.charAt(toIntStrict(pos)) & 0xFFFF))
***REMOVED***

func (r *Runtime) stringproto_codePointAt(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	p := call.Argument(0).ToInteger()
	size := s.length()
	if p < 0 || p >= int64(size) ***REMOVED***
		return _undefined
	***REMOVED***
	pos := toIntStrict(p)
	first := s.charAt(pos)
	if isUTF16FirstSurrogate(first) ***REMOVED***
		pos++
		if pos < size ***REMOVED***
			second := s.charAt(pos)
			if isUTF16SecondSurrogate(second) ***REMOVED***
				return intToValue(int64(utf16.DecodeRune(first, second)))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return intToValue(int64(first & 0xFFFF))
***REMOVED***

func (r *Runtime) stringproto_concat(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	strs := make([]valueString, len(call.Arguments)+1)
	strs[0] = call.This.toString()
	_, allAscii := strs[0].(asciiString)
	totalLen := strs[0].length()
	for i, arg := range call.Arguments ***REMOVED***
		s := arg.toString()
		if allAscii ***REMOVED***
			_, allAscii = s.(asciiString)
		***REMOVED***
		strs[i+1] = s
		totalLen += s.length()
	***REMOVED***

	if allAscii ***REMOVED***
		var buf strings.Builder
		buf.Grow(totalLen)
		for _, s := range strs ***REMOVED***
			buf.WriteString(s.String())
		***REMOVED***
		return asciiString(buf.String())
	***REMOVED*** else ***REMOVED***
		buf := make([]uint16, totalLen+1)
		buf[0] = unistring.BOM
		pos := 1
		for _, s := range strs ***REMOVED***
			switch s := s.(type) ***REMOVED***
			case asciiString:
				for i := 0; i < len(s); i++ ***REMOVED***
					buf[pos] = uint16(s[i])
					pos++
				***REMOVED***
			case unicodeString:
				copy(buf[pos:], s[1:])
				pos += s.length()
			***REMOVED***
		***REMOVED***
		return unicodeString(buf)
	***REMOVED***
***REMOVED***

func (r *Runtime) stringproto_endsWith(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	searchString := call.Argument(0)
	if isRegexp(searchString) ***REMOVED***
		panic(r.NewTypeError("First argument to String.prototype.endsWith must not be a regular expression"))
	***REMOVED***
	searchStr := searchString.toString()
	l := int64(s.length())
	var pos int64
	if posArg := call.Argument(1); posArg != _undefined ***REMOVED***
		pos = posArg.ToInteger()
	***REMOVED*** else ***REMOVED***
		pos = l
	***REMOVED***
	end := toIntStrict(min(max(pos, 0), l))
	searchLength := searchStr.length()
	start := end - searchLength
	if start < 0 ***REMOVED***
		return valueFalse
	***REMOVED***
	for i := 0; i < searchLength; i++ ***REMOVED***
		if s.charAt(start+i) != searchStr.charAt(i) ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) stringproto_includes(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	searchString := call.Argument(0)
	if isRegexp(searchString) ***REMOVED***
		panic(r.NewTypeError("First argument to String.prototype.includes must not be a regular expression"))
	***REMOVED***
	searchStr := searchString.toString()
	var pos int64
	if posArg := call.Argument(1); posArg != _undefined ***REMOVED***
		pos = posArg.ToInteger()
	***REMOVED*** else ***REMOVED***
		pos = 0
	***REMOVED***
	start := toIntStrict(min(max(pos, 0), int64(s.length())))
	if s.index(searchStr, start) != -1 ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) stringproto_indexOf(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	value := call.This.toString()
	target := call.Argument(0).toString()
	pos := call.Argument(1).ToInteger()

	if pos < 0 ***REMOVED***
		pos = 0
	***REMOVED*** else ***REMOVED***
		l := int64(value.length())
		if pos > l ***REMOVED***
			pos = l
		***REMOVED***
	***REMOVED***

	return intToValue(int64(value.index(target, toIntStrict(pos))))
***REMOVED***

func (r *Runtime) stringproto_lastIndexOf(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	value := call.This.toString()
	target := call.Argument(0).toString()
	numPos := call.Argument(1).ToNumber()

	var pos int64
	if f, ok := numPos.(valueFloat); ok && math.IsNaN(float64(f)) ***REMOVED***
		pos = int64(value.length())
	***REMOVED*** else ***REMOVED***
		pos = numPos.ToInteger()
		if pos < 0 ***REMOVED***
			pos = 0
		***REMOVED*** else ***REMOVED***
			l := int64(value.length())
			if pos > l ***REMOVED***
				pos = l
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return intToValue(int64(value.lastIndex(target, toIntStrict(pos))))
***REMOVED***

func (r *Runtime) stringproto_localeCompare(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	this := norm.NFD.String(call.This.toString().String())
	that := norm.NFD.String(call.Argument(0).toString().String())
	return intToValue(int64(r.collator().CompareString(this, that)))
***REMOVED***

func (r *Runtime) stringproto_match(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	regexp := call.Argument(0)
	if regexp != _undefined && regexp != _null ***REMOVED***
		if matcher := toMethod(r.getV(regexp, SymMatch)); matcher != nil ***REMOVED***
			return matcher(FunctionCall***REMOVED***
				This:      regexp,
				Arguments: []Value***REMOVED***call.This***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	var rx *regexpObject
	if regexp, ok := regexp.(*Object); ok ***REMOVED***
		rx, _ = regexp.self.(*regexpObject)
	***REMOVED***

	if rx == nil ***REMOVED***
		rx = r.newRegExp(regexp, nil, r.global.RegExpPrototype)
	***REMOVED***

	if matcher, ok := r.toObject(rx.getSym(SymMatch, nil)).self.assertCallable(); ok ***REMOVED***
		return matcher(FunctionCall***REMOVED***
			This:      rx.val,
			Arguments: []Value***REMOVED***call.This.toString()***REMOVED***,
		***REMOVED***)
	***REMOVED***

	panic(r.NewTypeError("RegExp matcher is not a function"))
***REMOVED***

func (r *Runtime) stringproto_matchAll(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	regexp := call.Argument(0)
	if regexp != _undefined && regexp != _null ***REMOVED***
		if isRegexp(regexp) ***REMOVED***
			if o, ok := regexp.(*Object); ok ***REMOVED***
				flags := o.self.getStr("flags", nil)
				r.checkObjectCoercible(flags)
				if !strings.Contains(flags.toString().String(), "g") ***REMOVED***
					panic(r.NewTypeError("RegExp doesn't have global flag set"))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if matcher := toMethod(r.getV(regexp, SymMatchAll)); matcher != nil ***REMOVED***
			return matcher(FunctionCall***REMOVED***
				This:      regexp,
				Arguments: []Value***REMOVED***call.This***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	rx := r.newRegExp(regexp, asciiString("g"), r.global.RegExpPrototype)

	if matcher, ok := r.toObject(rx.getSym(SymMatchAll, nil)).self.assertCallable(); ok ***REMOVED***
		return matcher(FunctionCall***REMOVED***
			This:      rx.val,
			Arguments: []Value***REMOVED***call.This.toString()***REMOVED***,
		***REMOVED***)
	***REMOVED***

	panic(r.NewTypeError("RegExp matcher is not a function"))
***REMOVED***

func (r *Runtime) stringproto_normalize(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	var form string
	if formArg := call.Argument(0); formArg != _undefined ***REMOVED***
		form = formArg.toString().toString().String()
	***REMOVED*** else ***REMOVED***
		form = "NFC"
	***REMOVED***
	var f norm.Form
	switch form ***REMOVED***
	case "NFC":
		f = norm.NFC
	case "NFD":
		f = norm.NFD
	case "NFKC":
		f = norm.NFKC
	case "NFKD":
		f = norm.NFKD
	default:
		panic(r.newError(r.global.RangeError, "The normalization form should be one of NFC, NFD, NFKC, NFKD"))
	***REMOVED***

	if s, ok := s.(unicodeString); ok ***REMOVED***
		ss := s.String()
		return newStringValue(f.String(ss))
	***REMOVED***

	return s
***REMOVED***

func (r *Runtime) stringproto_padEnd(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	maxLength := toLength(call.Argument(0))
	stringLength := int64(s.length())
	if maxLength <= stringLength ***REMOVED***
		return s
	***REMOVED***
	var filler valueString
	var fillerASCII bool
	if fillString := call.Argument(1); fillString != _undefined ***REMOVED***
		filler = fillString.toString()
		if filler.length() == 0 ***REMOVED***
			return s
		***REMOVED***
		_, fillerASCII = filler.(asciiString)
	***REMOVED*** else ***REMOVED***
		filler = asciiString(" ")
		fillerASCII = true
	***REMOVED***
	remaining := toIntStrict(maxLength - stringLength)
	_, stringASCII := s.(asciiString)
	if fillerASCII && stringASCII ***REMOVED***
		fl := filler.length()
		var sb strings.Builder
		sb.Grow(toIntStrict(maxLength))
		sb.WriteString(s.String())
		fs := filler.String()
		for remaining >= fl ***REMOVED***
			sb.WriteString(fs)
			remaining -= fl
		***REMOVED***
		if remaining > 0 ***REMOVED***
			sb.WriteString(fs[:remaining])
		***REMOVED***
		return asciiString(sb.String())
	***REMOVED***
	var sb unicodeStringBuilder
	sb.Grow(toIntStrict(maxLength))
	sb.WriteString(s)
	fl := filler.length()
	for remaining >= fl ***REMOVED***
		sb.WriteString(filler)
		remaining -= fl
	***REMOVED***
	if remaining > 0 ***REMOVED***
		sb.WriteString(filler.substring(0, remaining))
	***REMOVED***

	return sb.String()
***REMOVED***

func (r *Runtime) stringproto_padStart(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	maxLength := toLength(call.Argument(0))
	stringLength := int64(s.length())
	if maxLength <= stringLength ***REMOVED***
		return s
	***REMOVED***
	var filler valueString
	var fillerASCII bool
	if fillString := call.Argument(1); fillString != _undefined ***REMOVED***
		filler = fillString.toString()
		if filler.length() == 0 ***REMOVED***
			return s
		***REMOVED***
		_, fillerASCII = filler.(asciiString)
	***REMOVED*** else ***REMOVED***
		filler = asciiString(" ")
		fillerASCII = true
	***REMOVED***
	remaining := toIntStrict(maxLength - stringLength)
	_, stringASCII := s.(asciiString)
	if fillerASCII && stringASCII ***REMOVED***
		fl := filler.length()
		var sb strings.Builder
		sb.Grow(toIntStrict(maxLength))
		fs := filler.String()
		for remaining >= fl ***REMOVED***
			sb.WriteString(fs)
			remaining -= fl
		***REMOVED***
		if remaining > 0 ***REMOVED***
			sb.WriteString(fs[:remaining])
		***REMOVED***
		sb.WriteString(s.String())
		return asciiString(sb.String())
	***REMOVED***
	var sb unicodeStringBuilder
	sb.Grow(toIntStrict(maxLength))
	fl := filler.length()
	for remaining >= fl ***REMOVED***
		sb.WriteString(filler)
		remaining -= fl
	***REMOVED***
	if remaining > 0 ***REMOVED***
		sb.WriteString(filler.substring(0, remaining))
	***REMOVED***
	sb.WriteString(s)

	return sb.String()
***REMOVED***

func (r *Runtime) stringproto_repeat(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	n := call.Argument(0).ToNumber()
	if n == _positiveInf ***REMOVED***
		panic(r.newError(r.global.RangeError, "Invalid count value"))
	***REMOVED***
	numInt := n.ToInteger()
	if numInt < 0 ***REMOVED***
		panic(r.newError(r.global.RangeError, "Invalid count value"))
	***REMOVED***
	if numInt == 0 || s.length() == 0 ***REMOVED***
		return stringEmpty
	***REMOVED***
	num := toIntStrict(numInt)
	if s, ok := s.(asciiString); ok ***REMOVED***
		var sb strings.Builder
		sb.Grow(len(s) * num)
		for i := 0; i < num; i++ ***REMOVED***
			sb.WriteString(string(s))
		***REMOVED***
		return asciiString(sb.String())
	***REMOVED***

	var sb unicodeStringBuilder
	sb.Grow(s.length() * num)
	for i := 0; i < num; i++ ***REMOVED***
		sb.WriteString(s)
	***REMOVED***
	return sb.String()
***REMOVED***

func getReplaceValue(replaceValue Value) (str valueString, rcall func(FunctionCall) Value) ***REMOVED***
	if replaceValue, ok := replaceValue.(*Object); ok ***REMOVED***
		if c, ok := replaceValue.self.assertCallable(); ok ***REMOVED***
			rcall = c
			return
		***REMOVED***
	***REMOVED***
	str = replaceValue.toString()
	return
***REMOVED***

func stringReplace(s valueString, found [][]int, newstring valueString, rcall func(FunctionCall) Value) Value ***REMOVED***
	if len(found) == 0 ***REMOVED***
		return s
	***REMOVED***

	var str string
	var isASCII bool
	if astr, ok := s.(asciiString); ok ***REMOVED***
		str = string(astr)
		isASCII = true
	***REMOVED***

	var buf valueStringBuilder

	lastIndex := 0
	lengthS := s.length()
	if rcall != nil ***REMOVED***
		for _, item := range found ***REMOVED***
			if item[0] != lastIndex ***REMOVED***
				buf.WriteSubstring(s, lastIndex, item[0])
			***REMOVED***
			matchCount := len(item) / 2
			argumentList := make([]Value, matchCount+2)
			for index := 0; index < matchCount; index++ ***REMOVED***
				offset := 2 * index
				if item[offset] != -1 ***REMOVED***
					if isASCII ***REMOVED***
						argumentList[index] = asciiString(str[item[offset]:item[offset+1]])
					***REMOVED*** else ***REMOVED***
						argumentList[index] = s.substring(item[offset], item[offset+1])
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					argumentList[index] = _undefined
				***REMOVED***
			***REMOVED***
			argumentList[matchCount] = valueInt(item[0])
			argumentList[matchCount+1] = s
			replacement := rcall(FunctionCall***REMOVED***
				This:      _undefined,
				Arguments: argumentList,
			***REMOVED***).toString()
			buf.WriteString(replacement)
			lastIndex = item[1]
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, item := range found ***REMOVED***
			if item[0] != lastIndex ***REMOVED***
				buf.WriteString(s.substring(lastIndex, item[0]))
			***REMOVED***
			matchCount := len(item) / 2
			writeSubstitution(s, item[0], matchCount, func(idx int) valueString ***REMOVED***
				if item[idx*2] != -1 ***REMOVED***
					if isASCII ***REMOVED***
						return asciiString(str[item[idx*2]:item[idx*2+1]])
					***REMOVED***
					return s.substring(item[idx*2], item[idx*2+1])
				***REMOVED***
				return stringEmpty
			***REMOVED***, newstring, &buf)
			lastIndex = item[1]
		***REMOVED***
	***REMOVED***

	if lastIndex != lengthS ***REMOVED***
		buf.WriteString(s.substring(lastIndex, lengthS))
	***REMOVED***

	return buf.String()
***REMOVED***

func (r *Runtime) stringproto_replace(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	searchValue := call.Argument(0)
	replaceValue := call.Argument(1)
	if searchValue != _undefined && searchValue != _null ***REMOVED***
		if replacer := toMethod(r.getV(searchValue, SymReplace)); replacer != nil ***REMOVED***
			return replacer(FunctionCall***REMOVED***
				This:      searchValue,
				Arguments: []Value***REMOVED***call.This, replaceValue***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	s := call.This.toString()
	var found [][]int
	searchStr := searchValue.toString()
	pos := s.index(searchStr, 0)
	if pos != -1 ***REMOVED***
		found = append(found, []int***REMOVED***pos, pos + searchStr.length()***REMOVED***)
	***REMOVED***

	str, rcall := getReplaceValue(replaceValue)
	return stringReplace(s, found, str, rcall)
***REMOVED***

func (r *Runtime) stringproto_search(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	regexp := call.Argument(0)
	if regexp != _undefined && regexp != _null ***REMOVED***
		if searcher := toMethod(r.getV(regexp, SymSearch)); searcher != nil ***REMOVED***
			return searcher(FunctionCall***REMOVED***
				This:      regexp,
				Arguments: []Value***REMOVED***call.This***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	var rx *regexpObject
	if regexp, ok := regexp.(*Object); ok ***REMOVED***
		rx, _ = regexp.self.(*regexpObject)
	***REMOVED***

	if rx == nil ***REMOVED***
		rx = r.newRegExp(regexp, nil, r.global.RegExpPrototype)
	***REMOVED***

	if searcher, ok := r.toObject(rx.getSym(SymSearch, nil)).self.assertCallable(); ok ***REMOVED***
		return searcher(FunctionCall***REMOVED***
			This:      rx.val,
			Arguments: []Value***REMOVED***call.This.toString()***REMOVED***,
		***REMOVED***)
	***REMOVED***

	panic(r.NewTypeError("RegExp searcher is not a function"))
***REMOVED***

func (r *Runtime) stringproto_slice(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()

	l := int64(s.length())
	start := call.Argument(0).ToInteger()
	var end int64
	if arg1 := call.Argument(1); arg1 != _undefined ***REMOVED***
		end = arg1.ToInteger()
	***REMOVED*** else ***REMOVED***
		end = l
	***REMOVED***

	if start < 0 ***REMOVED***
		start += l
		if start < 0 ***REMOVED***
			start = 0
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if start > l ***REMOVED***
			start = l
		***REMOVED***
	***REMOVED***

	if end < 0 ***REMOVED***
		end += l
		if end < 0 ***REMOVED***
			end = 0
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if end > l ***REMOVED***
			end = l
		***REMOVED***
	***REMOVED***

	if end > start ***REMOVED***
		return s.substring(int(start), int(end))
	***REMOVED***
	return stringEmpty
***REMOVED***

func (r *Runtime) stringproto_split(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	separatorValue := call.Argument(0)
	limitValue := call.Argument(1)
	if separatorValue != _undefined && separatorValue != _null ***REMOVED***
		if splitter := toMethod(r.getV(separatorValue, SymSplit)); splitter != nil ***REMOVED***
			return splitter(FunctionCall***REMOVED***
				This:      separatorValue,
				Arguments: []Value***REMOVED***call.This, limitValue***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	s := call.This.toString()

	limit := -1
	if limitValue != _undefined ***REMOVED***
		limit = int(toUint32(limitValue))
	***REMOVED***

	separatorValue = separatorValue.ToString()

	if limit == 0 ***REMOVED***
		return r.newArrayValues(nil)
	***REMOVED***

	if separatorValue == _undefined ***REMOVED***
		return r.newArrayValues([]Value***REMOVED***s***REMOVED***)
	***REMOVED***

	separator := separatorValue.String()

	excess := false
	str := s.String()
	if limit > len(str) ***REMOVED***
		limit = len(str)
	***REMOVED***
	splitLimit := limit
	if limit > 0 ***REMOVED***
		splitLimit = limit + 1
		excess = true
	***REMOVED***

	// TODO handle invalid UTF-16
	split := strings.SplitN(str, separator, splitLimit)

	if excess && len(split) > limit ***REMOVED***
		split = split[:limit]
	***REMOVED***

	valueArray := make([]Value, len(split))
	for index, value := range split ***REMOVED***
		valueArray[index] = newStringValue(value)
	***REMOVED***

	return r.newArrayValues(valueArray)
***REMOVED***

func (r *Runtime) stringproto_startsWith(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	searchString := call.Argument(0)
	if isRegexp(searchString) ***REMOVED***
		panic(r.NewTypeError("First argument to String.prototype.startsWith must not be a regular expression"))
	***REMOVED***
	searchStr := searchString.toString()
	l := int64(s.length())
	var pos int64
	if posArg := call.Argument(1); posArg != _undefined ***REMOVED***
		pos = posArg.ToInteger()
	***REMOVED***
	start := toIntStrict(min(max(pos, 0), l))
	searchLength := searchStr.length()
	if int64(searchLength+start) > l ***REMOVED***
		return valueFalse
	***REMOVED***
	for i := 0; i < searchLength; i++ ***REMOVED***
		if s.charAt(start+i) != searchStr.charAt(i) ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) stringproto_substring(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()

	l := int64(s.length())
	intStart := call.Argument(0).ToInteger()
	var intEnd int64
	if end := call.Argument(1); end != _undefined ***REMOVED***
		intEnd = end.ToInteger()
	***REMOVED*** else ***REMOVED***
		intEnd = l
	***REMOVED***
	if intStart < 0 ***REMOVED***
		intStart = 0
	***REMOVED*** else if intStart > l ***REMOVED***
		intStart = l
	***REMOVED***

	if intEnd < 0 ***REMOVED***
		intEnd = 0
	***REMOVED*** else if intEnd > l ***REMOVED***
		intEnd = l
	***REMOVED***

	if intStart > intEnd ***REMOVED***
		intStart, intEnd = intEnd, intStart
	***REMOVED***

	return s.substring(int(intStart), int(intEnd))
***REMOVED***

func (r *Runtime) stringproto_toLowerCase(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()

	return s.toLower()
***REMOVED***

func (r *Runtime) stringproto_toUpperCase(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()

	return s.toUpper()
***REMOVED***

func (r *Runtime) stringproto_trim(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()

	// TODO handle invalid UTF-16
	return newStringValue(strings.Trim(s.String(), parser.WhitespaceChars))
***REMOVED***

func (r *Runtime) stringproto_trimEnd(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()

	// TODO handle invalid UTF-16
	return newStringValue(strings.TrimRight(s.String(), parser.WhitespaceChars))
***REMOVED***

func (r *Runtime) stringproto_trimStart(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()

	// TODO handle invalid UTF-16
	return newStringValue(strings.TrimLeft(s.String(), parser.WhitespaceChars))
***REMOVED***

func (r *Runtime) stringproto_substr(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.toString()
	start := call.Argument(0).ToInteger()
	var length int64
	sl := int64(s.length())
	if arg := call.Argument(1); arg != _undefined ***REMOVED***
		length = arg.ToInteger()
	***REMOVED*** else ***REMOVED***
		length = sl
	***REMOVED***

	if start < 0 ***REMOVED***
		start = max(sl+start, 0)
	***REMOVED***

	length = min(max(length, 0), sl-start)
	if length <= 0 ***REMOVED***
		return stringEmpty
	***REMOVED***

	return s.substring(int(start), int(start+length))
***REMOVED***

func (r *Runtime) stringIterProto_next(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	if iter, ok := thisObj.self.(*stringIterObject); ok ***REMOVED***
		return iter.next()
	***REMOVED***
	panic(r.NewTypeError("Method String Iterator.prototype.next called on incompatible receiver %s", thisObj.String()))
***REMOVED***

func (r *Runtime) createStringIterProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.IteratorPrototype, classObject)

	o._putProp("next", r.newNativeFunc(r.stringIterProto_next, nil, "next", nil, 0), true, false, true)
	o._putSym(SymToStringTag, valueProp(asciiString(classStringIterator), false, false, true))

	return o
***REMOVED***

func (r *Runtime) initString() ***REMOVED***
	r.global.StringIteratorPrototype = r.newLazyObject(r.createStringIterProto)
	r.global.StringPrototype = r.builtin_newString([]Value***REMOVED***stringEmpty***REMOVED***, r.global.ObjectPrototype)

	o := r.global.StringPrototype.self
	o._putProp("charAt", r.newNativeFunc(r.stringproto_charAt, nil, "charAt", nil, 1), true, false, true)
	o._putProp("charCodeAt", r.newNativeFunc(r.stringproto_charCodeAt, nil, "charCodeAt", nil, 1), true, false, true)
	o._putProp("codePointAt", r.newNativeFunc(r.stringproto_codePointAt, nil, "codePointAt", nil, 1), true, false, true)
	o._putProp("concat", r.newNativeFunc(r.stringproto_concat, nil, "concat", nil, 1), true, false, true)
	o._putProp("endsWith", r.newNativeFunc(r.stringproto_endsWith, nil, "endsWith", nil, 1), true, false, true)
	o._putProp("includes", r.newNativeFunc(r.stringproto_includes, nil, "includes", nil, 1), true, false, true)
	o._putProp("indexOf", r.newNativeFunc(r.stringproto_indexOf, nil, "indexOf", nil, 1), true, false, true)
	o._putProp("lastIndexOf", r.newNativeFunc(r.stringproto_lastIndexOf, nil, "lastIndexOf", nil, 1), true, false, true)
	o._putProp("localeCompare", r.newNativeFunc(r.stringproto_localeCompare, nil, "localeCompare", nil, 1), true, false, true)
	o._putProp("match", r.newNativeFunc(r.stringproto_match, nil, "match", nil, 1), true, false, true)
	o._putProp("matchAll", r.newNativeFunc(r.stringproto_matchAll, nil, "matchAll", nil, 1), true, false, true)
	o._putProp("normalize", r.newNativeFunc(r.stringproto_normalize, nil, "normalize", nil, 0), true, false, true)
	o._putProp("padEnd", r.newNativeFunc(r.stringproto_padEnd, nil, "padEnd", nil, 1), true, false, true)
	o._putProp("padStart", r.newNativeFunc(r.stringproto_padStart, nil, "padStart", nil, 1), true, false, true)
	o._putProp("repeat", r.newNativeFunc(r.stringproto_repeat, nil, "repeat", nil, 1), true, false, true)
	o._putProp("replace", r.newNativeFunc(r.stringproto_replace, nil, "replace", nil, 2), true, false, true)
	o._putProp("search", r.newNativeFunc(r.stringproto_search, nil, "search", nil, 1), true, false, true)
	o._putProp("slice", r.newNativeFunc(r.stringproto_slice, nil, "slice", nil, 2), true, false, true)
	o._putProp("split", r.newNativeFunc(r.stringproto_split, nil, "split", nil, 2), true, false, true)
	o._putProp("startsWith", r.newNativeFunc(r.stringproto_startsWith, nil, "startsWith", nil, 1), true, false, true)
	o._putProp("substring", r.newNativeFunc(r.stringproto_substring, nil, "substring", nil, 2), true, false, true)
	o._putProp("toLocaleLowerCase", r.newNativeFunc(r.stringproto_toLowerCase, nil, "toLocaleLowerCase", nil, 0), true, false, true)
	o._putProp("toLocaleUpperCase", r.newNativeFunc(r.stringproto_toUpperCase, nil, "toLocaleUpperCase", nil, 0), true, false, true)
	o._putProp("toLowerCase", r.newNativeFunc(r.stringproto_toLowerCase, nil, "toLowerCase", nil, 0), true, false, true)
	o._putProp("toString", r.newNativeFunc(r.stringproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("toUpperCase", r.newNativeFunc(r.stringproto_toUpperCase, nil, "toUpperCase", nil, 0), true, false, true)
	o._putProp("trim", r.newNativeFunc(r.stringproto_trim, nil, "trim", nil, 0), true, false, true)
	trimEnd := r.newNativeFunc(r.stringproto_trimEnd, nil, "trimEnd", nil, 0)
	trimStart := r.newNativeFunc(r.stringproto_trimStart, nil, "trimStart", nil, 0)
	o._putProp("trimEnd", trimEnd, true, false, true)
	o._putProp("trimStart", trimStart, true, false, true)
	o._putProp("trimRight", trimEnd, true, false, true)
	o._putProp("trimLeft", trimStart, true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.stringproto_valueOf, nil, "valueOf", nil, 0), true, false, true)

	o._putSym(SymIterator, valueProp(r.newNativeFunc(r.stringproto_iterator, nil, "[Symbol.iterator]", nil, 0), true, false, true))

	// Annex B
	o._putProp("substr", r.newNativeFunc(r.stringproto_substr, nil, "substr", nil, 2), true, false, true)

	r.global.String = r.newNativeFunc(r.builtin_String, r.builtin_newString, "String", r.global.StringPrototype, 1)
	o = r.global.String.self
	o._putProp("fromCharCode", r.newNativeFunc(r.string_fromcharcode, nil, "fromCharCode", nil, 1), true, false, true)
	o._putProp("fromCodePoint", r.newNativeFunc(r.string_fromcodepoint, nil, "fromCodePoint", nil, 1), true, false, true)
	o._putProp("raw", r.newNativeFunc(r.string_raw, nil, "raw", nil, 1), true, false, true)

	r.addToGlobal("String", r.global.String)

	r.stringSingleton = r.builtin_new(r.global.String, nil).self.(*stringObject)
***REMOVED***
