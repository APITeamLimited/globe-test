package goja

import (
	"bytes"
	"github.com/dop251/goja/parser"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
	"math"
	"strings"
	"unicode/utf8"
)

func (r *Runtime) collator() *collate.Collator ***REMOVED***
	collator := r._collator
	if collator == nil ***REMOVED***
		collator = collate.New(language.Und)
		r._collator = collator
	***REMOVED***
	return collator
***REMOVED***

func (r *Runtime) builtin_String(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) > 0 ***REMOVED***
		arg := call.Arguments[0]
		if _, ok := arg.assertString(); ok ***REMOVED***
			return arg
		***REMOVED***
		return arg.ToString()
	***REMOVED*** else ***REMOVED***
		return newStringValue("")
	***REMOVED***
***REMOVED***

func (r *Runtime) _newString(s valueString) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	o := &stringObject***REMOVED******REMOVED***
	o.class = classString
	o.val = v
	o.extensible = true
	v.self = o
	o.prototype = r.global.StringPrototype
	if s != nil ***REMOVED***
		o.value = s
	***REMOVED***
	o.init()
	return v
***REMOVED***

func (r *Runtime) builtin_newString(args []Value) *Object ***REMOVED***
	var s valueString
	if len(args) > 0 ***REMOVED***
		s = args[0].ToString()
	***REMOVED*** else ***REMOVED***
		s = stringEmpty
	***REMOVED***
	return r._newString(s)
***REMOVED***

func searchSubstringUTF8(str, search string) (ret [][]int) ***REMOVED***
	searchPos := 0
	l := len(str)
	if searchPos < l ***REMOVED***
		p := strings.Index(str[searchPos:], search)
		if p != -1 ***REMOVED***
			p += searchPos
			searchPos = p + len(search)
			ret = append(ret, []int***REMOVED***p, searchPos***REMOVED***)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) stringproto_toStringValueOf(this Value, funcName string) Value ***REMOVED***
	if str, ok := this.assertString(); ok ***REMOVED***
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

func (r *Runtime) string_fromcharcode(call FunctionCall) Value ***REMOVED***
	b := make([]byte, len(call.Arguments))
	for i, arg := range call.Arguments ***REMOVED***
		chr := toUInt16(arg)
		if chr >= utf8.RuneSelf ***REMOVED***
			bb := make([]uint16, len(call.Arguments))
			for j := 0; j < i; j++ ***REMOVED***
				bb[j] = uint16(b[j])
			***REMOVED***
			bb[i] = chr
			i++
			for j, arg := range call.Arguments[i:] ***REMOVED***
				bb[i+j] = toUInt16(arg)
			***REMOVED***
			return unicodeString(bb)
		***REMOVED***
		b[i] = byte(chr)
	***REMOVED***

	return asciiString(b)
***REMOVED***

func (r *Runtime) stringproto_charAt(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()
	pos := call.Argument(0).ToInteger()
	if pos < 0 || pos >= s.length() ***REMOVED***
		return stringEmpty
	***REMOVED***
	return newStringValue(string(s.charAt(pos)))
***REMOVED***

func (r *Runtime) stringproto_charCodeAt(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()
	pos := call.Argument(0).ToInteger()
	if pos < 0 || pos >= s.length() ***REMOVED***
		return _NaN
	***REMOVED***
	return intToValue(int64(s.charAt(pos) & 0xFFFF))
***REMOVED***

func (r *Runtime) stringproto_concat(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	strs := make([]valueString, len(call.Arguments)+1)
	strs[0] = call.This.ToString()
	_, allAscii := strs[0].(asciiString)
	totalLen := strs[0].length()
	for i, arg := range call.Arguments ***REMOVED***
		s := arg.ToString()
		if allAscii ***REMOVED***
			_, allAscii = s.(asciiString)
		***REMOVED***
		strs[i+1] = s
		totalLen += s.length()
	***REMOVED***

	if allAscii ***REMOVED***
		buf := bytes.NewBuffer(make([]byte, 0, totalLen))
		for _, s := range strs ***REMOVED***
			buf.WriteString(s.String())
		***REMOVED***
		return asciiString(buf.String())
	***REMOVED*** else ***REMOVED***
		buf := make([]uint16, totalLen)
		pos := int64(0)
		for _, s := range strs ***REMOVED***
			switch s := s.(type) ***REMOVED***
			case asciiString:
				for i := 0; i < len(s); i++ ***REMOVED***
					buf[pos] = uint16(s[i])
					pos++
				***REMOVED***
			case unicodeString:
				copy(buf[pos:], s)
				pos += s.length()
			***REMOVED***
		***REMOVED***
		return unicodeString(buf)
	***REMOVED***
***REMOVED***

func (r *Runtime) stringproto_indexOf(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	value := call.This.ToString()
	target := call.Argument(0).ToString()
	pos := call.Argument(1).ToInteger()

	if pos < 0 ***REMOVED***
		pos = 0
	***REMOVED*** else ***REMOVED***
		l := value.length()
		if pos > l ***REMOVED***
			pos = l
		***REMOVED***
	***REMOVED***

	return intToValue(value.index(target, pos))
***REMOVED***

func (r *Runtime) stringproto_lastIndexOf(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	value := call.This.ToString()
	target := call.Argument(0).ToString()
	numPos := call.Argument(1).ToNumber()

	var pos int64
	if f, ok := numPos.assertFloat(); ok && math.IsNaN(f) ***REMOVED***
		pos = value.length()
	***REMOVED*** else ***REMOVED***
		pos = numPos.ToInteger()
		if pos < 0 ***REMOVED***
			pos = 0
		***REMOVED*** else ***REMOVED***
			l := value.length()
			if pos > l ***REMOVED***
				pos = l
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return intToValue(value.lastIndex(target, pos))
***REMOVED***

func (r *Runtime) stringproto_localeCompare(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	this := norm.NFD.String(call.This.String())
	that := norm.NFD.String(call.Argument(0).String())
	return intToValue(int64(r.collator().CompareString(this, that)))
***REMOVED***

func (r *Runtime) stringproto_match(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()
	regexp := call.Argument(0)
	var rx *regexpObject
	if regexp, ok := regexp.(*Object); ok ***REMOVED***
		rx, _ = regexp.self.(*regexpObject)
	***REMOVED***

	if rx == nil ***REMOVED***
		rx = r.builtin_newRegExp([]Value***REMOVED***regexp***REMOVED***).self.(*regexpObject)
	***REMOVED***

	if rx.global ***REMOVED***
		rx.putStr("lastIndex", intToValue(0), false)
		var a []Value
		var previousLastIndex int64
		for ***REMOVED***
			match, result := rx.execRegexp(s)
			if !match ***REMOVED***
				break
			***REMOVED***
			thisIndex := rx.getStr("lastIndex").ToInteger()
			if thisIndex == previousLastIndex ***REMOVED***
				previousLastIndex++
				rx.putStr("lastIndex", intToValue(previousLastIndex), false)
			***REMOVED*** else ***REMOVED***
				previousLastIndex = thisIndex
			***REMOVED***
			a = append(a, s.substring(int64(result[0]), int64(result[1])))
		***REMOVED***
		if len(a) == 0 ***REMOVED***
			return _null
		***REMOVED***
		return r.newArrayValues(a)
	***REMOVED*** else ***REMOVED***
		return rx.exec(s)
	***REMOVED***
***REMOVED***

func (r *Runtime) stringproto_replace(call FunctionCall) Value ***REMOVED***
	s := call.This.ToString()
	var str string
	var isASCII bool
	if astr, ok := s.(asciiString); ok ***REMOVED***
		str = string(astr)
		isASCII = true
	***REMOVED*** else ***REMOVED***
		str = s.String()
	***REMOVED***
	searchValue := call.Argument(0)
	replaceValue := call.Argument(1)

	var found [][]int

	if searchValue, ok := searchValue.(*Object); ok ***REMOVED***
		if regexp, ok := searchValue.self.(*regexpObject); ok ***REMOVED***
			find := 1
			if regexp.global ***REMOVED***
				find = -1
			***REMOVED***
			if isASCII ***REMOVED***
				found = regexp.pattern.FindAllSubmatchIndexASCII(str, find)
			***REMOVED*** else ***REMOVED***
				found = regexp.pattern.FindAllSubmatchIndexUTF8(str, find)
			***REMOVED***
			if found == nil ***REMOVED***
				return s
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if found == nil ***REMOVED***
		found = searchSubstringUTF8(str, searchValue.String())
	***REMOVED***

	if len(found) == 0 ***REMOVED***
		return s
	***REMOVED***

	var buf bytes.Buffer
	lastIndex := 0

	var rcall func(FunctionCall) Value

	if replaceValue, ok := replaceValue.(*Object); ok ***REMOVED***
		if c, ok := replaceValue.self.assertCallable(); ok ***REMOVED***
			rcall = c
		***REMOVED***
	***REMOVED***

	if rcall != nil ***REMOVED***
		for _, item := range found ***REMOVED***
			if item[0] != lastIndex ***REMOVED***
				buf.WriteString(str[lastIndex:item[0]])
			***REMOVED***
			matchCount := len(item) / 2
			argumentList := make([]Value, matchCount+2)
			for index := 0; index < matchCount; index++ ***REMOVED***
				offset := 2 * index
				if item[offset] != -1 ***REMOVED***
					if isASCII ***REMOVED***
						argumentList[index] = asciiString(str[item[offset]:item[offset+1]])
					***REMOVED*** else ***REMOVED***
						argumentList[index] = newStringValue(str[item[offset]:item[offset+1]])
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
			***REMOVED***).String()
			buf.WriteString(replacement)
			lastIndex = item[1]
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		newstring := replaceValue.String()

		for _, item := range found ***REMOVED***
			if item[0] != lastIndex ***REMOVED***
				buf.WriteString(str[lastIndex:item[0]])
			***REMOVED***
			matches := len(item) / 2
			for i := 0; i < len(newstring); i++ ***REMOVED***
				if newstring[i] == '$' && i < len(newstring)-1 ***REMOVED***
					ch := newstring[i+1]
					switch ch ***REMOVED***
					case '$':
						buf.WriteByte('$')
					case '`':
						buf.WriteString(str[0:item[0]])
					case '\'':
						buf.WriteString(str[item[1]:])
					case '&':
						buf.WriteString(str[item[0]:item[1]])
					default:
						matchNumber := 0
						l := 0
						for _, ch := range newstring[i+1:] ***REMOVED***
							if ch >= '0' && ch <= '9' ***REMOVED***
								m := matchNumber*10 + int(ch-'0')
								if m >= matches ***REMOVED***
									break
								***REMOVED***
								matchNumber = m
								l++
							***REMOVED*** else ***REMOVED***
								break
							***REMOVED***
						***REMOVED***
						if l > 0 ***REMOVED***
							offset := 2 * matchNumber
							if offset < len(item) && item[offset] != -1 ***REMOVED***
								buf.WriteString(str[item[offset]:item[offset+1]])
							***REMOVED***
							i += l - 1
						***REMOVED*** else ***REMOVED***
							buf.WriteByte('$')
							buf.WriteByte(ch)
						***REMOVED***

					***REMOVED***
					i++
				***REMOVED*** else ***REMOVED***
					buf.WriteByte(newstring[i])
				***REMOVED***
			***REMOVED***
			lastIndex = item[1]
		***REMOVED***
	***REMOVED***

	if lastIndex != len(str) ***REMOVED***
		buf.WriteString(str[lastIndex:])
	***REMOVED***

	return newStringValue(buf.String())
***REMOVED***

func (r *Runtime) stringproto_search(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()
	regexp := call.Argument(0)
	var rx *regexpObject
	if regexp, ok := regexp.(*Object); ok ***REMOVED***
		rx, _ = regexp.self.(*regexpObject)
	***REMOVED***

	if rx == nil ***REMOVED***
		rx = r.builtin_newRegExp([]Value***REMOVED***regexp***REMOVED***).self.(*regexpObject)
	***REMOVED***

	match, result := rx.execRegexp(s)
	if !match ***REMOVED***
		return intToValue(-1)
	***REMOVED***
	return intToValue(int64(result[0]))
***REMOVED***

func (r *Runtime) stringproto_slice(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()

	l := s.length()
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
		return s.substring(start, end)
	***REMOVED***
	return stringEmpty
***REMOVED***

func (r *Runtime) stringproto_split(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()

	separatorValue := call.Argument(0)
	limitValue := call.Argument(1)
	limit := -1
	if limitValue != _undefined ***REMOVED***
		limit = int(toUInt32(limitValue))
	***REMOVED***

	if limit == 0 ***REMOVED***
		return r.newArrayValues(nil)
	***REMOVED***

	if separatorValue == _undefined ***REMOVED***
		return r.newArrayValues([]Value***REMOVED***s***REMOVED***)
	***REMOVED***

	var search *regexpObject
	if o, ok := separatorValue.(*Object); ok ***REMOVED***
		search, _ = o.self.(*regexpObject)
	***REMOVED***

	if search != nil ***REMOVED***
		targetLength := s.length()
		valueArray := []Value***REMOVED******REMOVED***
		result := search.pattern.FindAllSubmatchIndex(s, -1)
		lastIndex := 0
		found := 0

		for _, match := range result ***REMOVED***
			if match[0] == match[1] ***REMOVED***
				// FIXME Ugh, this is a hack
				if match[0] == 0 || int64(match[0]) == targetLength ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***

			if lastIndex != match[0] ***REMOVED***
				valueArray = append(valueArray, s.substring(int64(lastIndex), int64(match[0])))
				found++
			***REMOVED*** else if lastIndex == match[0] ***REMOVED***
				if lastIndex != -1 ***REMOVED***
					valueArray = append(valueArray, stringEmpty)
					found++
				***REMOVED***
			***REMOVED***

			lastIndex = match[1]
			if found == limit ***REMOVED***
				goto RETURN
			***REMOVED***

			captureCount := len(match) / 2
			for index := 1; index < captureCount; index++ ***REMOVED***
				offset := index * 2
				var value Value
				if match[offset] != -1 ***REMOVED***
					value = s.substring(int64(match[offset]), int64(match[offset+1]))
				***REMOVED*** else ***REMOVED***
					value = _undefined
				***REMOVED***
				valueArray = append(valueArray, value)
				found++
				if found == limit ***REMOVED***
					goto RETURN
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if found != limit ***REMOVED***
			if int64(lastIndex) != targetLength ***REMOVED***
				valueArray = append(valueArray, s.substring(int64(lastIndex), targetLength))
			***REMOVED*** else ***REMOVED***
				valueArray = append(valueArray, stringEmpty)
			***REMOVED***
		***REMOVED***

	RETURN:
		return r.newArrayValues(valueArray)

	***REMOVED*** else ***REMOVED***
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

***REMOVED***

func (r *Runtime) stringproto_substring(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()

	l := s.length()
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

	return s.substring(intStart, intEnd)
***REMOVED***

func (r *Runtime) stringproto_toLowerCase(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()

	return s.toLower()
***REMOVED***

func (r *Runtime) stringproto_toUpperCase(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()

	return s.toUpper()
***REMOVED***

func (r *Runtime) stringproto_trim(call FunctionCall) Value ***REMOVED***
	r.checkObjectCoercible(call.This)
	s := call.This.ToString()

	return newStringValue(strings.Trim(s.String(), parser.WhitespaceChars))
***REMOVED***

func (r *Runtime) stringproto_substr(call FunctionCall) Value ***REMOVED***
	s := call.This.ToString()
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

	return s.substring(start, start+length)
***REMOVED***

func (r *Runtime) initString() ***REMOVED***
	r.global.StringPrototype = r.builtin_newString([]Value***REMOVED***stringEmpty***REMOVED***)

	o := r.global.StringPrototype.self
	o.(*stringObject).prototype = r.global.ObjectPrototype
	o._putProp("toString", r.newNativeFunc(r.stringproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.stringproto_valueOf, nil, "valueOf", nil, 0), true, false, true)
	o._putProp("charAt", r.newNativeFunc(r.stringproto_charAt, nil, "charAt", nil, 1), true, false, true)
	o._putProp("charCodeAt", r.newNativeFunc(r.stringproto_charCodeAt, nil, "charCodeAt", nil, 1), true, false, true)
	o._putProp("concat", r.newNativeFunc(r.stringproto_concat, nil, "concat", nil, 1), true, false, true)
	o._putProp("indexOf", r.newNativeFunc(r.stringproto_indexOf, nil, "indexOf", nil, 1), true, false, true)
	o._putProp("lastIndexOf", r.newNativeFunc(r.stringproto_lastIndexOf, nil, "lastIndexOf", nil, 1), true, false, true)
	o._putProp("localeCompare", r.newNativeFunc(r.stringproto_localeCompare, nil, "localeCompare", nil, 1), true, false, true)
	o._putProp("match", r.newNativeFunc(r.stringproto_match, nil, "match", nil, 1), true, false, true)
	o._putProp("replace", r.newNativeFunc(r.stringproto_replace, nil, "replace", nil, 2), true, false, true)
	o._putProp("search", r.newNativeFunc(r.stringproto_search, nil, "search", nil, 1), true, false, true)
	o._putProp("slice", r.newNativeFunc(r.stringproto_slice, nil, "slice", nil, 2), true, false, true)
	o._putProp("split", r.newNativeFunc(r.stringproto_split, nil, "split", nil, 2), true, false, true)
	o._putProp("substring", r.newNativeFunc(r.stringproto_substring, nil, "substring", nil, 2), true, false, true)
	o._putProp("toLowerCase", r.newNativeFunc(r.stringproto_toLowerCase, nil, "toLowerCase", nil, 0), true, false, true)
	o._putProp("toLocaleLowerCase", r.newNativeFunc(r.stringproto_toLowerCase, nil, "toLocaleLowerCase", nil, 0), true, false, true)
	o._putProp("toUpperCase", r.newNativeFunc(r.stringproto_toUpperCase, nil, "toUpperCase", nil, 0), true, false, true)
	o._putProp("toLocaleUpperCase", r.newNativeFunc(r.stringproto_toUpperCase, nil, "toLocaleUpperCase", nil, 0), true, false, true)
	o._putProp("trim", r.newNativeFunc(r.stringproto_trim, nil, "trim", nil, 0), true, false, true)

	// Annex B
	o._putProp("substr", r.newNativeFunc(r.stringproto_substr, nil, "substr", nil, 2), true, false, true)

	r.global.String = r.newNativeFunc(r.builtin_String, r.builtin_newString, "String", r.global.StringPrototype, 1)
	o = r.global.String.self
	o._putProp("fromCharCode", r.newNativeFunc(r.string_fromcharcode, nil, "fromCharCode", nil, 1), true, false, true)

	r.addToGlobal("String", r.global.String)

	r.stringSingleton = r.builtin_new(r.global.String, nil).self.(*stringObject)
***REMOVED***
