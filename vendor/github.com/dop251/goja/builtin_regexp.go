package goja

import (
	"fmt"
	"github.com/dop251/goja/parser"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func (r *Runtime) newRegexpObject(proto *Object) *regexpObject ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	o := &regexpObject***REMOVED******REMOVED***
	o.class = classRegExp
	o.val = v
	o.extensible = true
	v.self = o
	o.prototype = proto
	o.init()
	return o
***REMOVED***

func (r *Runtime) newRegExpp(pattern *regexpPattern, patternStr valueString, proto *Object) *Object ***REMOVED***
	o := r.newRegexpObject(proto)

	o.pattern = pattern
	o.source = patternStr

	return o.val
***REMOVED***

func decodeHex(s string) (int, bool) ***REMOVED***
	var hex int
	for i := 0; i < len(s); i++ ***REMOVED***
		var n byte
		chr := s[i]
		switch ***REMOVED***
		case '0' <= chr && chr <= '9':
			n = chr - '0'
		case 'a' <= chr && chr <= 'f':
			n = chr - 'a' + 10
		case 'A' <= chr && chr <= 'F':
			n = chr - 'A' + 10
		default:
			return 0, false
		***REMOVED***
		hex = hex*16 + int(n)
	***REMOVED***
	return hex, true
***REMOVED***

func writeHex4(b *strings.Builder, i int) ***REMOVED***
	b.WriteByte(hex[i>>12])
	b.WriteByte(hex[(i>>8)&0xF])
	b.WriteByte(hex[(i>>4)&0xF])
	b.WriteByte(hex[i&0xF])
***REMOVED***

// Convert any valid surrogate pairs in the form of \uXXXX\uXXXX to unicode characters
func convertRegexpToUnicode(patternStr string) string ***REMOVED***
	var sb strings.Builder
	pos := 0
	for i := 0; i < len(patternStr)-11; ***REMOVED***
		r, size := utf8.DecodeRuneInString(patternStr[i:])
		if r == '\\' ***REMOVED***
			i++
			if patternStr[i] == 'u' && patternStr[i+5] == '\\' && patternStr[i+6] == 'u' ***REMOVED***
				if first, ok := decodeHex(patternStr[i+1 : i+5]); ok ***REMOVED***
					if isUTF16FirstSurrogate(rune(first)) ***REMOVED***
						if second, ok := decodeHex(patternStr[i+7 : i+11]); ok ***REMOVED***
							if isUTF16SecondSurrogate(rune(second)) ***REMOVED***
								r = utf16.DecodeRune(rune(first), rune(second))
								sb.WriteString(patternStr[pos : i-1])
								sb.WriteRune(r)
								i += 11
								pos = i
								continue
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
			i++
		***REMOVED*** else ***REMOVED***
			i += size
		***REMOVED***
	***REMOVED***
	if pos > 0 ***REMOVED***
		sb.WriteString(patternStr[pos:])
		return sb.String()
	***REMOVED***
	return patternStr
***REMOVED***

// Convert any extended unicode characters to UTF-16 in the form of \uXXXX\uXXXX
func convertRegexpToUtf16(patternStr string) string ***REMOVED***
	var sb strings.Builder
	pos := 0
	var prevRune rune
	for i := 0; i < len(patternStr); ***REMOVED***
		r, size := utf8.DecodeRuneInString(patternStr[i:])
		if r > 0xFFFF ***REMOVED***
			sb.WriteString(patternStr[pos:i])
			if prevRune == '\\' ***REMOVED***
				sb.WriteRune('\\')
			***REMOVED***
			first, second := utf16.EncodeRune(r)
			sb.WriteString(`\u`)
			writeHex4(&sb, int(first))
			sb.WriteString(`\u`)
			writeHex4(&sb, int(second))
			pos = i + size
		***REMOVED***
		i += size
		prevRune = r
	***REMOVED***
	if pos > 0 ***REMOVED***
		sb.WriteString(patternStr[pos:])
		return sb.String()
	***REMOVED***
	return patternStr
***REMOVED***

// convert any broken UTF-16 surrogate pairs to \uXXXX
func escapeInvalidUtf16(s valueString) string ***REMOVED***
	if ascii, ok := s.(asciiString); ok ***REMOVED***
		return ascii.String()
	***REMOVED***
	var sb strings.Builder
	rd := &lenientUtf16Decoder***REMOVED***utf16Reader: s.utf16Reader(0)***REMOVED***
	pos := 0
	utf8Size := 0
	var utf8Buf [utf8.UTFMax]byte
	for ***REMOVED***
		c, size, err := rd.ReadRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if utf16.IsSurrogate(c) ***REMOVED***
			if sb.Len() == 0 ***REMOVED***
				sb.Grow(utf8Size + 7)
				hrd := s.reader(0)
				var c rune
				for p := 0; p < pos; ***REMOVED***
					var size int
					var err error
					c, size, err = hrd.ReadRune()
					if err != nil ***REMOVED***
						// will not happen
						panic(fmt.Errorf("error while reading string head %q, pos: %d: %w", s.String(), pos, err))
					***REMOVED***
					sb.WriteRune(c)
					p += size
				***REMOVED***
				if c == '\\' ***REMOVED***
					sb.WriteRune(c)
				***REMOVED***
			***REMOVED***
			sb.WriteString(`\u`)
			writeHex4(&sb, int(c))
		***REMOVED*** else ***REMOVED***
			if sb.Len() > 0 ***REMOVED***
				sb.WriteRune(c)
			***REMOVED*** else ***REMOVED***
				utf8Size += utf8.EncodeRune(utf8Buf[:], c)
				pos += size
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if sb.Len() > 0 ***REMOVED***
		return sb.String()
	***REMOVED***
	return s.String()
***REMOVED***

func compileRegexpFromValueString(patternStr valueString, flags string) (*regexpPattern, error) ***REMOVED***
	return compileRegexp(escapeInvalidUtf16(patternStr), flags)
***REMOVED***

func compileRegexp(patternStr, flags string) (p *regexpPattern, err error) ***REMOVED***
	var global, ignoreCase, multiline, sticky, unicode bool
	var wrapper *regexpWrapper
	var wrapper2 *regexp2Wrapper

	if flags != "" ***REMOVED***
		invalidFlags := func() ***REMOVED***
			err = fmt.Errorf("Invalid flags supplied to RegExp constructor '%s'", flags)
		***REMOVED***
		for _, chr := range flags ***REMOVED***
			switch chr ***REMOVED***
			case 'g':
				if global ***REMOVED***
					invalidFlags()
					return
				***REMOVED***
				global = true
			case 'm':
				if multiline ***REMOVED***
					invalidFlags()
					return
				***REMOVED***
				multiline = true
			case 'i':
				if ignoreCase ***REMOVED***
					invalidFlags()
					return
				***REMOVED***
				ignoreCase = true
			case 'y':
				if sticky ***REMOVED***
					invalidFlags()
					return
				***REMOVED***
				sticky = true
			case 'u':
				if unicode ***REMOVED***
					invalidFlags()
				***REMOVED***
				unicode = true
			default:
				invalidFlags()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if unicode ***REMOVED***
		patternStr = convertRegexpToUnicode(patternStr)
	***REMOVED*** else ***REMOVED***
		patternStr = convertRegexpToUtf16(patternStr)
	***REMOVED***

	re2Str, err1 := parser.TransformRegExp(patternStr)
	if err1 == nil ***REMOVED***
		re2flags := ""
		if multiline ***REMOVED***
			re2flags += "m"
		***REMOVED***
		if ignoreCase ***REMOVED***
			re2flags += "i"
		***REMOVED***
		if len(re2flags) > 0 ***REMOVED***
			re2Str = fmt.Sprintf("(?%s:%s)", re2flags, re2Str)
		***REMOVED***

		pattern, err1 := regexp.Compile(re2Str)
		if err1 != nil ***REMOVED***
			err = fmt.Errorf("Invalid regular expression (re2): %s (%v)", re2Str, err1)
			return
		***REMOVED***
		wrapper = (*regexpWrapper)(pattern)
	***REMOVED*** else ***REMOVED***
		if _, incompat := err1.(parser.RegexpErrorIncompatible); !incompat ***REMOVED***
			err = err1
			return
		***REMOVED***
		wrapper2, err = compileRegexp2(patternStr, multiline, ignoreCase)
		if err != nil ***REMOVED***
			err = fmt.Errorf("Invalid regular expression (regexp2): %s (%v)", patternStr, err)
			return
		***REMOVED***
	***REMOVED***

	p = &regexpPattern***REMOVED***
		src:            patternStr,
		regexpWrapper:  wrapper,
		regexp2Wrapper: wrapper2,
		global:         global,
		ignoreCase:     ignoreCase,
		multiline:      multiline,
		sticky:         sticky,
		unicode:        unicode,
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) _newRegExp(patternStr valueString, flags string, proto *Object) *Object ***REMOVED***
	pattern, err := compileRegexpFromValueString(patternStr, flags)
	if err != nil ***REMOVED***
		panic(r.newSyntaxError(err.Error(), -1))
	***REMOVED***
	return r.newRegExpp(pattern, patternStr, proto)
***REMOVED***

func (r *Runtime) builtin_newRegExp(args []Value, proto *Object) *Object ***REMOVED***
	var patternVal, flagsVal Value
	if len(args) > 0 ***REMOVED***
		patternVal = args[0]
	***REMOVED***
	if len(args) > 1 ***REMOVED***
		flagsVal = args[1]
	***REMOVED***
	return r.newRegExp(patternVal, flagsVal, proto)
***REMOVED***

func (r *Runtime) newRegExp(patternVal, flagsVal Value, proto *Object) *Object ***REMOVED***
	var pattern valueString
	var flags string
	if isRegexp(patternVal) ***REMOVED*** // this may have side effects so need to call it anyway
		if obj, ok := patternVal.(*Object); ok ***REMOVED***
			if rx, ok := obj.self.(*regexpObject); ok ***REMOVED***
				if flagsVal == nil || flagsVal == _undefined ***REMOVED***
					return rx.clone()
				***REMOVED*** else ***REMOVED***
					return r._newRegExp(rx.source, flagsVal.toString().String(), proto)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				pattern = nilSafe(obj.self.getStr("source", nil)).toString()
				if flagsVal == nil || flagsVal == _undefined ***REMOVED***
					flags = nilSafe(obj.self.getStr("flags", nil)).toString().String()
				***REMOVED*** else ***REMOVED***
					flags = flagsVal.toString().String()
				***REMOVED***
				goto exit
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if patternVal != nil && patternVal != _undefined ***REMOVED***
		pattern = patternVal.toString()
	***REMOVED***
	if flagsVal != nil && flagsVal != _undefined ***REMOVED***
		flags = flagsVal.toString().String()
	***REMOVED***

	if pattern == nil ***REMOVED***
		pattern = stringEmpty
	***REMOVED***
exit:
	return r._newRegExp(pattern, flags, proto)
***REMOVED***

func (r *Runtime) builtin_RegExp(call FunctionCall) Value ***REMOVED***
	pattern := call.Argument(0)
	patternIsRegExp := isRegexp(pattern)
	flags := call.Argument(1)
	if patternIsRegExp && flags == _undefined ***REMOVED***
		if obj, ok := call.Argument(0).(*Object); ok ***REMOVED***
			patternConstructor := obj.self.getStr("constructor", nil)
			if patternConstructor == r.global.RegExp ***REMOVED***
				return pattern
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return r.newRegExp(pattern, flags, r.global.RegExpPrototype)
***REMOVED***

func (r *Runtime) regexpproto_compile(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		var (
			pattern *regexpPattern
			source  valueString
			flags   string
			err     error
		)
		patternVal := call.Argument(0)
		flagsVal := call.Argument(1)
		if o, ok := patternVal.(*Object); ok ***REMOVED***
			if p, ok := o.self.(*regexpObject); ok ***REMOVED***
				if flagsVal != _undefined ***REMOVED***
					panic(r.NewTypeError("Cannot supply flags when constructing one RegExp from another"))
				***REMOVED***
				this.pattern = p.pattern
				this.source = p.source
				goto exit
			***REMOVED***
		***REMOVED***
		if patternVal != _undefined ***REMOVED***
			source = patternVal.toString()
		***REMOVED*** else ***REMOVED***
			source = stringEmpty
		***REMOVED***
		if flagsVal != _undefined ***REMOVED***
			flags = flagsVal.toString().String()
		***REMOVED***
		pattern, err = compileRegexpFromValueString(source, flags)
		if err != nil ***REMOVED***
			panic(r.newSyntaxError(err.Error(), -1))
		***REMOVED***
		this.pattern = pattern
		this.source = source
	exit:
		this.setOwnStr("lastIndex", intToValue(0), true)
		return call.This
	***REMOVED***

	panic(r.NewTypeError("Method RegExp.prototype.compile called on incompatible receiver %s", call.This.toString()))
***REMOVED***

func (r *Runtime) regexpproto_exec(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		return this.exec(call.Argument(0).toString())
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.exec called on incompatible receiver %s", call.This.toString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_test(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.test(call.Argument(0).toString()) ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.test called on incompatible receiver %s", call.This.toString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_toString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	if this := r.checkStdRegexp(obj); this != nil ***REMOVED***
		var sb valueStringBuilder
		sb.WriteRune('/')
		if !this.writeEscapedSource(&sb) ***REMOVED***
			sb.WriteString(this.source)
		***REMOVED***
		sb.WriteRune('/')
		if this.pattern.global ***REMOVED***
			sb.WriteRune('g')
		***REMOVED***
		if this.pattern.ignoreCase ***REMOVED***
			sb.WriteRune('i')
		***REMOVED***
		if this.pattern.multiline ***REMOVED***
			sb.WriteRune('m')
		***REMOVED***
		if this.pattern.unicode ***REMOVED***
			sb.WriteRune('u')
		***REMOVED***
		if this.pattern.sticky ***REMOVED***
			sb.WriteRune('y')
		***REMOVED***
		return sb.String()
	***REMOVED***
	pattern := nilSafe(obj.self.getStr("source", nil)).toString()
	flags := nilSafe(obj.self.getStr("flags", nil)).toString()
	var sb valueStringBuilder
	sb.WriteRune('/')
	sb.WriteString(pattern)
	sb.WriteRune('/')
	sb.WriteString(flags)
	return sb.String()
***REMOVED***

func (r *regexpObject) writeEscapedSource(sb *valueStringBuilder) bool ***REMOVED***
	if r.source.length() == 0 ***REMOVED***
		sb.WriteString(asciiString("(?:)"))
		return true
	***REMOVED***
	pos := 0
	lastPos := 0
	rd := &lenientUtf16Decoder***REMOVED***utf16Reader: r.source.utf16Reader(0)***REMOVED***
L:
	for ***REMOVED***
		c, size, err := rd.ReadRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		switch c ***REMOVED***
		case '\\':
			pos++
			_, size, err = rd.ReadRune()
			if err != nil ***REMOVED***
				break L
			***REMOVED***
		case '/', '\u000a', '\u000d', '\u2028', '\u2029':
			sb.WriteSubstring(r.source, lastPos, pos)
			sb.WriteRune('\\')
			switch c ***REMOVED***
			case '\u000a':
				sb.WriteRune('n')
			case '\u000d':
				sb.WriteRune('r')
			default:
				sb.WriteRune('u')
				sb.WriteRune(rune(hex[c>>12]))
				sb.WriteRune(rune(hex[(c>>8)&0xF]))
				sb.WriteRune(rune(hex[(c>>4)&0xF]))
				sb.WriteRune(rune(hex[c&0xF]))
			***REMOVED***
			lastPos = pos + size
		***REMOVED***
		pos += size
	***REMOVED***
	if lastPos > 0 ***REMOVED***
		sb.WriteSubstring(r.source, lastPos, r.source.length())
		return true
	***REMOVED***
	return false
***REMOVED***

func (r *Runtime) regexpproto_getSource(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		var sb valueStringBuilder
		if this.writeEscapedSource(&sb) ***REMOVED***
			return sb.String()
		***REMOVED***
		return this.source
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.source getter called on incompatible receiver")
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getGlobal(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.pattern.global ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.global getter called on incompatible receiver %s", call.This.toString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getMultiline(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.pattern.multiline ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.multiline getter called on incompatible receiver %s", call.This.toString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getIgnoreCase(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.pattern.ignoreCase ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.ignoreCase getter called on incompatible receiver %s", call.This.toString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getUnicode(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.pattern.unicode ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.unicode getter called on incompatible receiver %s", call.This.toString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getSticky(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.pattern.sticky ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.sticky getter called on incompatible receiver %s", call.This.toString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getFlags(call FunctionCall) Value ***REMOVED***
	var global, ignoreCase, multiline, sticky, unicode bool

	thisObj := r.toObject(call.This)
	size := 0
	if v := thisObj.self.getStr("global", nil); v != nil ***REMOVED***
		global = v.ToBoolean()
		if global ***REMOVED***
			size++
		***REMOVED***
	***REMOVED***
	if v := thisObj.self.getStr("ignoreCase", nil); v != nil ***REMOVED***
		ignoreCase = v.ToBoolean()
		if ignoreCase ***REMOVED***
			size++
		***REMOVED***
	***REMOVED***
	if v := thisObj.self.getStr("multiline", nil); v != nil ***REMOVED***
		multiline = v.ToBoolean()
		if multiline ***REMOVED***
			size++
		***REMOVED***
	***REMOVED***
	if v := thisObj.self.getStr("sticky", nil); v != nil ***REMOVED***
		sticky = v.ToBoolean()
		if sticky ***REMOVED***
			size++
		***REMOVED***
	***REMOVED***
	if v := thisObj.self.getStr("unicode", nil); v != nil ***REMOVED***
		unicode = v.ToBoolean()
		if unicode ***REMOVED***
			size++
		***REMOVED***
	***REMOVED***

	var sb strings.Builder
	sb.Grow(size)
	if global ***REMOVED***
		sb.WriteByte('g')
	***REMOVED***
	if ignoreCase ***REMOVED***
		sb.WriteByte('i')
	***REMOVED***
	if multiline ***REMOVED***
		sb.WriteByte('m')
	***REMOVED***
	if unicode ***REMOVED***
		sb.WriteByte('u')
	***REMOVED***
	if sticky ***REMOVED***
		sb.WriteByte('y')
	***REMOVED***

	return asciiString(sb.String())
***REMOVED***

func (r *Runtime) regExpExec(execFn func(FunctionCall) Value, rxObj *Object, arg Value) Value ***REMOVED***
	res := execFn(FunctionCall***REMOVED***
		This:      rxObj,
		Arguments: []Value***REMOVED***arg***REMOVED***,
	***REMOVED***)

	if res != _null ***REMOVED***
		if _, ok := res.(*Object); !ok ***REMOVED***
			panic(r.NewTypeError("RegExp exec method returned something other than an Object or null"))
		***REMOVED***
	***REMOVED***

	return res
***REMOVED***

func (r *Runtime) getGlobalRegexpMatches(rxObj *Object, s valueString) []Value ***REMOVED***
	fullUnicode := nilSafe(rxObj.self.getStr("unicode", nil)).ToBoolean()
	rxObj.self.setOwnStr("lastIndex", intToValue(0), true)
	execFn, ok := r.toObject(rxObj.self.getStr("exec", nil)).self.assertCallable()
	if !ok ***REMOVED***
		panic(r.NewTypeError("exec is not a function"))
	***REMOVED***
	var a []Value
	for ***REMOVED***
		res := r.regExpExec(execFn, rxObj, s)
		if res == _null ***REMOVED***
			break
		***REMOVED***
		a = append(a, res)
		matchStr := nilSafe(r.toObject(res).self.getIdx(valueInt(0), nil)).toString()
		if matchStr.length() == 0 ***REMOVED***
			thisIndex := toLength(rxObj.self.getStr("lastIndex", nil))
			rxObj.self.setOwnStr("lastIndex", valueInt(advanceStringIndex64(s, thisIndex, fullUnicode)), true)
		***REMOVED***
	***REMOVED***

	return a
***REMOVED***

func (r *Runtime) regexpproto_stdMatcherGeneric(rxObj *Object, s valueString) Value ***REMOVED***
	rx := rxObj.self
	global := rx.getStr("global", nil)
	if global != nil && global.ToBoolean() ***REMOVED***
		a := r.getGlobalRegexpMatches(rxObj, s)
		if len(a) == 0 ***REMOVED***
			return _null
		***REMOVED***
		ar := make([]Value, 0, len(a))
		for _, result := range a ***REMOVED***
			obj := r.toObject(result)
			matchStr := nilSafe(obj.self.getIdx(valueInt(0), nil)).ToString()
			ar = append(ar, matchStr)
		***REMOVED***
		return r.newArrayValues(ar)
	***REMOVED***

	execFn, ok := r.toObject(rx.getStr("exec", nil)).self.assertCallable()
	if !ok ***REMOVED***
		panic(r.NewTypeError("exec is not a function"))
	***REMOVED***

	return r.regExpExec(execFn, rxObj, s)
***REMOVED***

func (r *Runtime) checkStdRegexp(rxObj *Object) *regexpObject ***REMOVED***
	if deoptimiseRegexp ***REMOVED***
		return nil
	***REMOVED***

	rx, ok := rxObj.self.(*regexpObject)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	if !rx.standard || rx.prototype == nil || rx.prototype.self != r.global.stdRegexpProto ***REMOVED***
		return nil
	***REMOVED***

	return rx
***REMOVED***

func (r *Runtime) regexpproto_stdMatcher(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	s := call.Argument(0).toString()
	rx := r.checkStdRegexp(thisObj)
	if rx == nil ***REMOVED***
		return r.regexpproto_stdMatcherGeneric(thisObj, s)
	***REMOVED***
	if rx.pattern.global ***REMOVED***
		res := rx.pattern.findAllSubmatchIndex(s, 0, -1, rx.pattern.sticky)
		if len(res) == 0 ***REMOVED***
			rx.setOwnStr("lastIndex", intToValue(0), true)
			return _null
		***REMOVED***
		a := make([]Value, 0, len(res))
		for _, result := range res ***REMOVED***
			a = append(a, s.substring(result[0], result[1]))
		***REMOVED***
		rx.setOwnStr("lastIndex", intToValue(int64(res[len(res)-1][1])), true)
		return r.newArrayValues(a)
	***REMOVED*** else ***REMOVED***
		return rx.exec(s)
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_stdSearchGeneric(rxObj *Object, arg valueString) Value ***REMOVED***
	rx := rxObj.self
	previousLastIndex := nilSafe(rx.getStr("lastIndex", nil))
	zero := intToValue(0)
	if !previousLastIndex.SameAs(zero) ***REMOVED***
		rx.setOwnStr("lastIndex", zero, true)
	***REMOVED***
	execFn, ok := r.toObject(rx.getStr("exec", nil)).self.assertCallable()
	if !ok ***REMOVED***
		panic(r.NewTypeError("exec is not a function"))
	***REMOVED***

	result := r.regExpExec(execFn, rxObj, arg)
	currentLastIndex := nilSafe(rx.getStr("lastIndex", nil))
	if !currentLastIndex.SameAs(previousLastIndex) ***REMOVED***
		rx.setOwnStr("lastIndex", previousLastIndex, true)
	***REMOVED***

	if result == _null ***REMOVED***
		return intToValue(-1)
	***REMOVED***

	return r.toObject(result).self.getStr("index", nil)
***REMOVED***

func (r *Runtime) regexpproto_stdSearch(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	s := call.Argument(0).toString()
	rx := r.checkStdRegexp(thisObj)
	if rx == nil ***REMOVED***
		return r.regexpproto_stdSearchGeneric(thisObj, s)
	***REMOVED***

	previousLastIndex := rx.getStr("lastIndex", nil)
	rx.setOwnStr("lastIndex", intToValue(0), true)

	match, result := rx.execRegexp(s)
	rx.setOwnStr("lastIndex", previousLastIndex, true)

	if !match ***REMOVED***
		return intToValue(-1)
	***REMOVED***
	return intToValue(int64(result[0]))
***REMOVED***

func (r *Runtime) regexpproto_stdSplitterGeneric(splitter *Object, s valueString, limit Value, unicodeMatching bool) Value ***REMOVED***
	var a []Value
	var lim int64
	if limit == nil || limit == _undefined ***REMOVED***
		lim = maxInt - 1
	***REMOVED*** else ***REMOVED***
		lim = toLength(limit)
	***REMOVED***
	if lim == 0 ***REMOVED***
		return r.newArrayValues(a)
	***REMOVED***
	size := s.length()
	p := 0
	execFn := toMethod(splitter.ToObject(r).self.getStr("exec", nil)) // must be non-nil

	if size == 0 ***REMOVED***
		if r.regExpExec(execFn, splitter, s) == _null ***REMOVED***
			a = append(a, s)
		***REMOVED***
		return r.newArrayValues(a)
	***REMOVED***

	q := p
	for q < size ***REMOVED***
		splitter.self.setOwnStr("lastIndex", intToValue(int64(q)), true)
		z := r.regExpExec(execFn, splitter, s)
		if z == _null ***REMOVED***
			q = advanceStringIndex(s, q, unicodeMatching)
		***REMOVED*** else ***REMOVED***
			z := r.toObject(z)
			e := toLength(splitter.self.getStr("lastIndex", nil))
			if e == int64(p) ***REMOVED***
				q = advanceStringIndex(s, q, unicodeMatching)
			***REMOVED*** else ***REMOVED***
				a = append(a, s.substring(p, q))
				if int64(len(a)) == lim ***REMOVED***
					return r.newArrayValues(a)
				***REMOVED***
				if e > int64(size) ***REMOVED***
					p = size
				***REMOVED*** else ***REMOVED***
					p = int(e)
				***REMOVED***
				numberOfCaptures := max(toLength(z.self.getStr("length", nil))-1, 0)
				for i := int64(1); i <= numberOfCaptures; i++ ***REMOVED***
					a = append(a, z.self.getIdx(valueInt(i), nil))
					if int64(len(a)) == lim ***REMOVED***
						return r.newArrayValues(a)
					***REMOVED***
				***REMOVED***
				q = p
			***REMOVED***
		***REMOVED***
	***REMOVED***
	a = append(a, s.substring(p, size))
	return r.newArrayValues(a)
***REMOVED***

func advanceStringIndex(s valueString, pos int, unicode bool) int ***REMOVED***
	next := pos + 1
	if !unicode ***REMOVED***
		return next
	***REMOVED***
	l := s.length()
	if next >= l ***REMOVED***
		return next
	***REMOVED***
	if !isUTF16FirstSurrogate(s.charAt(pos)) ***REMOVED***
		return next
	***REMOVED***
	if !isUTF16SecondSurrogate(s.charAt(next)) ***REMOVED***
		return next
	***REMOVED***
	return next + 1
***REMOVED***

func advanceStringIndex64(s valueString, pos int64, unicode bool) int64 ***REMOVED***
	next := pos + 1
	if !unicode ***REMOVED***
		return next
	***REMOVED***
	l := int64(s.length())
	if next >= l ***REMOVED***
		return next
	***REMOVED***
	if !isUTF16FirstSurrogate(s.charAt(int(pos))) ***REMOVED***
		return next
	***REMOVED***
	if !isUTF16SecondSurrogate(s.charAt(int(next))) ***REMOVED***
		return next
	***REMOVED***
	return next + 1
***REMOVED***

func (r *Runtime) regexpproto_stdSplitter(call FunctionCall) Value ***REMOVED***
	rxObj := r.toObject(call.This)
	s := call.Argument(0).toString()
	limitValue := call.Argument(1)
	var splitter *Object
	search := r.checkStdRegexp(rxObj)
	c := r.speciesConstructorObj(rxObj, r.global.RegExp)
	if search == nil || c != r.global.RegExp ***REMOVED***
		flags := nilSafe(rxObj.self.getStr("flags", nil)).toString()
		flagsStr := flags.String()

		// Add 'y' flag if missing
		if !strings.Contains(flagsStr, "y") ***REMOVED***
			flags = flags.concat(asciiString("y"))
		***REMOVED***
		splitter = r.toConstructor(c)([]Value***REMOVED***rxObj, flags***REMOVED***, nil)
		search = r.checkStdRegexp(splitter)
		if search == nil ***REMOVED***
			return r.regexpproto_stdSplitterGeneric(splitter, s, limitValue, strings.Contains(flagsStr, "u"))
		***REMOVED***
	***REMOVED***

	pattern := search.pattern // toUint32() may recompile the pattern, but we still need to use the original
	limit := -1
	if limitValue != _undefined ***REMOVED***
		limit = int(toUint32(limitValue))
	***REMOVED***

	if limit == 0 ***REMOVED***
		return r.newArrayValues(nil)
	***REMOVED***

	targetLength := s.length()
	var valueArray []Value
	lastIndex := 0
	found := 0

	result := pattern.findAllSubmatchIndex(s, 0, -1, false)
	if targetLength == 0 ***REMOVED***
		if result == nil ***REMOVED***
			valueArray = append(valueArray, s)
		***REMOVED***
		goto RETURN
	***REMOVED***

	for _, match := range result ***REMOVED***
		if match[0] == match[1] ***REMOVED***
			// FIXME Ugh, this is a hack
			if match[0] == 0 || match[0] == targetLength ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if lastIndex != match[0] ***REMOVED***
			valueArray = append(valueArray, s.substring(lastIndex, match[0]))
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
				value = s.substring(match[offset], match[offset+1])
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
		if lastIndex != targetLength ***REMOVED***
			valueArray = append(valueArray, s.substring(lastIndex, targetLength))
		***REMOVED*** else ***REMOVED***
			valueArray = append(valueArray, stringEmpty)
		***REMOVED***
	***REMOVED***

RETURN:
	return r.newArrayValues(valueArray)
***REMOVED***

func (r *Runtime) regexpproto_stdReplacerGeneric(rxObj *Object, s, replaceStr valueString, rcall func(FunctionCall) Value) Value ***REMOVED***
	var results []Value
	if nilSafe(rxObj.self.getStr("global", nil)).ToBoolean() ***REMOVED***
		results = r.getGlobalRegexpMatches(rxObj, s)
	***REMOVED*** else ***REMOVED***
		execFn := toMethod(rxObj.self.getStr("exec", nil)) // must be non-nil
		result := r.regExpExec(execFn, rxObj, s)
		if result != _null ***REMOVED***
			results = append(results, result)
		***REMOVED***
	***REMOVED***
	lengthS := s.length()
	nextSourcePosition := 0
	var resultBuf valueStringBuilder
	for _, result := range results ***REMOVED***
		obj := r.toObject(result)
		nCaptures := max(toLength(obj.self.getStr("length", nil))-1, 0)
		matched := nilSafe(obj.self.getIdx(valueInt(0), nil)).toString()
		matchLength := matched.length()
		position := toIntStrict(max(min(nilSafe(obj.self.getStr("index", nil)).ToInteger(), int64(lengthS)), 0))
		var captures []Value
		if rcall != nil ***REMOVED***
			captures = make([]Value, 0, nCaptures+3)
		***REMOVED*** else ***REMOVED***
			captures = make([]Value, 0, nCaptures+1)
		***REMOVED***
		captures = append(captures, matched)
		for n := int64(1); n <= nCaptures; n++ ***REMOVED***
			capN := nilSafe(obj.self.getIdx(valueInt(n), nil))
			if capN != _undefined ***REMOVED***
				capN = capN.ToString()
			***REMOVED***
			captures = append(captures, capN)
		***REMOVED***
		var replacement valueString
		if rcall != nil ***REMOVED***
			captures = append(captures, intToValue(int64(position)), s)
			replacement = rcall(FunctionCall***REMOVED***
				This:      _undefined,
				Arguments: captures,
			***REMOVED***).toString()
			if position >= nextSourcePosition ***REMOVED***
				resultBuf.WriteString(s.substring(nextSourcePosition, position))
				resultBuf.WriteString(replacement)
				nextSourcePosition = position + matchLength
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if position >= nextSourcePosition ***REMOVED***
				resultBuf.WriteString(s.substring(nextSourcePosition, position))
				writeSubstitution(s, position, len(captures), func(idx int) valueString ***REMOVED***
					capture := captures[idx]
					if capture != _undefined ***REMOVED***
						return capture.toString()
					***REMOVED***
					return stringEmpty
				***REMOVED***, replaceStr, &resultBuf)
				nextSourcePosition = position + matchLength
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if nextSourcePosition < lengthS ***REMOVED***
		resultBuf.WriteString(s.substring(nextSourcePosition, lengthS))
	***REMOVED***
	return resultBuf.String()
***REMOVED***

func writeSubstitution(s valueString, position int, numCaptures int, getCapture func(int) valueString, replaceStr valueString, buf *valueStringBuilder) ***REMOVED***
	l := s.length()
	rl := replaceStr.length()
	matched := getCapture(0)
	tailPos := position + matched.length()

	for i := 0; i < rl; i++ ***REMOVED***
		c := replaceStr.charAt(i)
		if c == '$' && i < rl-1 ***REMOVED***
			ch := replaceStr.charAt(i + 1)
			switch ch ***REMOVED***
			case '$':
				buf.WriteRune('$')
			case '`':
				buf.WriteString(s.substring(0, position))
			case '\'':
				if tailPos < l ***REMOVED***
					buf.WriteString(s.substring(tailPos, l))
				***REMOVED***
			case '&':
				buf.WriteString(matched)
			default:
				matchNumber := 0
				j := i + 1
				for j < rl ***REMOVED***
					ch := replaceStr.charAt(j)
					if ch >= '0' && ch <= '9' ***REMOVED***
						m := matchNumber*10 + int(ch-'0')
						if m >= numCaptures ***REMOVED***
							break
						***REMOVED***
						matchNumber = m
						j++
					***REMOVED*** else ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
				if matchNumber > 0 ***REMOVED***
					buf.WriteString(getCapture(matchNumber))
					i = j - 1
					continue
				***REMOVED*** else ***REMOVED***
					buf.WriteRune('$')
					buf.WriteRune(ch)
				***REMOVED***
			***REMOVED***
			i++
		***REMOVED*** else ***REMOVED***
			buf.WriteRune(c)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_stdReplacer(call FunctionCall) Value ***REMOVED***
	rxObj := r.toObject(call.This)
	s := call.Argument(0).toString()
	replaceStr, rcall := getReplaceValue(call.Argument(1))

	rx := r.checkStdRegexp(rxObj)
	if rx == nil ***REMOVED***
		return r.regexpproto_stdReplacerGeneric(rxObj, s, replaceStr, rcall)
	***REMOVED***

	var index int64
	find := 1
	if rx.pattern.global ***REMOVED***
		find = -1
		rx.setOwnStr("lastIndex", intToValue(0), true)
	***REMOVED*** else ***REMOVED***
		index = rx.getLastIndex()
	***REMOVED***
	found := rx.pattern.findAllSubmatchIndex(s, toIntStrict(index), find, rx.pattern.sticky)
	if len(found) > 0 ***REMOVED***
		if !rx.updateLastIndex(index, found[0], found[len(found)-1]) ***REMOVED***
			found = nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rx.updateLastIndex(index, nil, nil)
	***REMOVED***

	return stringReplace(s, found, replaceStr, rcall)
***REMOVED***

func (r *Runtime) initRegExp() ***REMOVED***
	o := r.newGuardedObject(r.global.ObjectPrototype, classObject)
	r.global.RegExpPrototype = o.val
	r.global.stdRegexpProto = o
	o._putProp("compile", r.newNativeFunc(r.regexpproto_compile, nil, "compile", nil, 2), true, false, true)
	o._putProp("exec", r.newNativeFunc(r.regexpproto_exec, nil, "exec", nil, 1), true, false, true)
	o._putProp("test", r.newNativeFunc(r.regexpproto_test, nil, "test", nil, 1), true, false, true)
	o._putProp("toString", r.newNativeFunc(r.regexpproto_toString, nil, "toString", nil, 0), true, false, true)
	o.setOwnStr("source", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getSource, nil, "get source", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.setOwnStr("global", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getGlobal, nil, "get global", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.setOwnStr("multiline", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getMultiline, nil, "get multiline", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.setOwnStr("ignoreCase", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getIgnoreCase, nil, "get ignoreCase", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.setOwnStr("unicode", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getUnicode, nil, "get unicode", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.setOwnStr("sticky", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getSticky, nil, "get sticky", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.setOwnStr("flags", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getFlags, nil, "get flags", nil, 0),
		accessor:     true,
	***REMOVED***, false)

	o._putSym(SymMatch, valueProp(r.newNativeFunc(r.regexpproto_stdMatcher, nil, "[Symbol.match]", nil, 1), true, false, true))
	o._putSym(SymSearch, valueProp(r.newNativeFunc(r.regexpproto_stdSearch, nil, "[Symbol.search]", nil, 1), true, false, true))
	o._putSym(SymSplit, valueProp(r.newNativeFunc(r.regexpproto_stdSplitter, nil, "[Symbol.split]", nil, 2), true, false, true))
	o._putSym(SymReplace, valueProp(r.newNativeFunc(r.regexpproto_stdReplacer, nil, "[Symbol.replace]", nil, 2), true, false, true))
	o.guard("exec", "global", "multiline", "ignoreCase", "unicode", "sticky")

	r.global.RegExp = r.newNativeFunc(r.builtin_RegExp, r.builtin_newRegExp, "RegExp", r.global.RegExpPrototype, 2)
	rx := r.global.RegExp.self
	rx._putSym(SymSpecies, &valueProperty***REMOVED***
		getterFunc:   r.newNativeFunc(r.returnThis, nil, "get [Symbol.species]", nil, 0),
		accessor:     true,
		configurable: true,
	***REMOVED***)
	r.addToGlobal("RegExp", r.global.RegExp)
***REMOVED***
