package goja

import (
	"errors"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

const hexUpper = "0123456789ABCDEF"

var (
	parseFloatRegexp = regexp.MustCompile(`^([+-]?(?:Infinity|[0-9]*\.?[0-9]*(?:[eE][+-]?[0-9]+)?))`)
)

func (r *Runtime) builtin_isNaN(call FunctionCall) Value ***REMOVED***
	if math.IsNaN(call.Argument(0).ToFloat()) ***REMOVED***
		return valueTrue
	***REMOVED*** else ***REMOVED***
		return valueFalse
	***REMOVED***
***REMOVED***

func (r *Runtime) builtin_parseInt(call FunctionCall) Value ***REMOVED***
	str := call.Argument(0).ToString().toTrimmedUTF8()
	radix := int(toInt32(call.Argument(1)))
	v, _ := parseInt(str, radix)
	return v
***REMOVED***

func (r *Runtime) builtin_parseFloat(call FunctionCall) Value ***REMOVED***
	m := parseFloatRegexp.FindStringSubmatch(call.Argument(0).ToString().toTrimmedUTF8())
	if len(m) == 2 ***REMOVED***
		if s := m[1]; s != "" && s != "+" && s != "-" ***REMOVED***
			switch s ***REMOVED***
			case "+", "-":
			case "Infinity", "+Infinity":
				return _positiveInf
			case "-Infinity":
				return _negativeInf
			default:
				f, err := strconv.ParseFloat(s, 64)
				if err == nil || isRangeErr(err) ***REMOVED***
					return floatToValue(f)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return _NaN
***REMOVED***

func (r *Runtime) builtin_isFinite(call FunctionCall) Value ***REMOVED***
	f := call.Argument(0).ToFloat()
	if math.IsNaN(f) || math.IsInf(f, 0) ***REMOVED***
		return valueFalse
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) _encode(uriString valueString, unescaped *[256]bool) valueString ***REMOVED***
	reader := uriString.reader(0)
	utf8Buf := make([]byte, utf8.UTFMax)
	needed := false
	l := 0
	for ***REMOVED***
		rn, _, err := reader.ReadRune()
		if err != nil ***REMOVED***
			if err != io.EOF ***REMOVED***
				panic(r.newError(r.global.URIError, "Malformed URI"))
			***REMOVED***
			break
		***REMOVED***

		if rn >= utf8.RuneSelf ***REMOVED***
			needed = true
			l += utf8.EncodeRune(utf8Buf, rn) * 3
		***REMOVED*** else if !unescaped[rn] ***REMOVED***
			needed = true
			l += 3
		***REMOVED*** else ***REMOVED***
			l++
		***REMOVED***
	***REMOVED***

	if !needed ***REMOVED***
		return uriString
	***REMOVED***

	buf := make([]byte, l)
	i := 0
	reader = uriString.reader(0)
	for ***REMOVED***
		rn, _, err := reader.ReadRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***

		if rn >= utf8.RuneSelf ***REMOVED***
			n := utf8.EncodeRune(utf8Buf, rn)
			for _, b := range utf8Buf[:n] ***REMOVED***
				buf[i] = '%'
				buf[i+1] = hexUpper[b>>4]
				buf[i+2] = hexUpper[b&15]
				i += 3
			***REMOVED***
		***REMOVED*** else if !unescaped[rn] ***REMOVED***
			buf[i] = '%'
			buf[i+1] = hexUpper[rn>>4]
			buf[i+2] = hexUpper[rn&15]
			i += 3
		***REMOVED*** else ***REMOVED***
			buf[i] = byte(rn)
			i++
		***REMOVED***
	***REMOVED***
	return asciiString(string(buf))
***REMOVED***

func (r *Runtime) _decode(sv valueString, reservedSet *[256]bool) valueString ***REMOVED***
	s := sv.String()
	hexCount := 0
	for i := 0; i < len(s); ***REMOVED***
		switch s[i] ***REMOVED***
		case '%':
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) ***REMOVED***
				panic(r.newError(r.global.URIError, "Malformed URI"))
			***REMOVED***
			c := unhex(s[i+1])<<4 | unhex(s[i+2])
			if !reservedSet[c] ***REMOVED***
				hexCount++
			***REMOVED***
			i += 3
		default:
			i++
		***REMOVED***
	***REMOVED***

	if hexCount == 0 ***REMOVED***
		return sv
	***REMOVED***

	t := make([]byte, len(s)-hexCount*2)
	j := 0
	isUnicode := false
	for i := 0; i < len(s); ***REMOVED***
		ch := s[i]
		switch ch ***REMOVED***
		case '%':
			c := unhex(s[i+1])<<4 | unhex(s[i+2])
			if reservedSet[c] ***REMOVED***
				t[j] = s[i]
				t[j+1] = s[i+1]
				t[j+2] = s[i+2]
				j += 3
			***REMOVED*** else ***REMOVED***
				t[j] = c
				if c >= utf8.RuneSelf ***REMOVED***
					isUnicode = true
				***REMOVED***
				j++
			***REMOVED***
			i += 3
		default:
			if ch >= utf8.RuneSelf ***REMOVED***
				isUnicode = true
			***REMOVED***
			t[j] = ch
			j++
			i++
		***REMOVED***
	***REMOVED***

	if !isUnicode ***REMOVED***
		return asciiString(t)
	***REMOVED***

	us := make([]rune, 0, len(s))
	for len(t) > 0 ***REMOVED***
		rn, size := utf8.DecodeRune(t)
		if rn == utf8.RuneError ***REMOVED***
			if size != 3 || t[0] != 0xef || t[1] != 0xbf || t[2] != 0xbd ***REMOVED***
				panic(r.newError(r.global.URIError, "Malformed URI"))
			***REMOVED***
		***REMOVED***
		us = append(us, rn)
		t = t[size:]
	***REMOVED***
	return unicodeString(utf16.Encode(us))
***REMOVED***

func ishex(c byte) bool ***REMOVED***
	switch ***REMOVED***
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	***REMOVED***
	return false
***REMOVED***

func unhex(c byte) byte ***REMOVED***
	switch ***REMOVED***
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	***REMOVED***
	return 0
***REMOVED***

func (r *Runtime) builtin_decodeURI(call FunctionCall) Value ***REMOVED***
	uriString := call.Argument(0).ToString()
	return r._decode(uriString, &uriReservedHash)
***REMOVED***

func (r *Runtime) builtin_decodeURIComponent(call FunctionCall) Value ***REMOVED***
	uriString := call.Argument(0).ToString()
	return r._decode(uriString, &emptyEscapeSet)
***REMOVED***

func (r *Runtime) builtin_encodeURI(call FunctionCall) Value ***REMOVED***
	uriString := call.Argument(0).ToString()
	return r._encode(uriString, &uriReservedUnescapedHash)
***REMOVED***

func (r *Runtime) builtin_encodeURIComponent(call FunctionCall) Value ***REMOVED***
	uriString := call.Argument(0).ToString()
	return r._encode(uriString, &uriUnescaped)
***REMOVED***

func (r *Runtime) builtin_escape(call FunctionCall) Value ***REMOVED***
	s := call.Argument(0).ToString()
	var sb strings.Builder
	l := s.length()
	for i := int64(0); i < l; i++ ***REMOVED***
		r := uint16(s.charAt(i))
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' ||
			r == '@' || r == '*' || r == '_' || r == '+' || r == '-' || r == '.' || r == '/' ***REMOVED***
			sb.WriteByte(byte(r))
		***REMOVED*** else if r <= 0xff ***REMOVED***
			sb.WriteByte('%')
			sb.WriteByte(hexUpper[r>>4])
			sb.WriteByte(hexUpper[r&0xf])
		***REMOVED*** else ***REMOVED***
			sb.WriteString("%u")
			sb.WriteByte(hexUpper[r>>12])
			sb.WriteByte(hexUpper[(r>>8)&0xf])
			sb.WriteByte(hexUpper[(r>>4)&0xf])
			sb.WriteByte(hexUpper[r&0xf])
		***REMOVED***
	***REMOVED***
	return asciiString(sb.String())
***REMOVED***

func (r *Runtime) builtin_unescape(call FunctionCall) Value ***REMOVED***
	s := call.Argument(0).ToString()
	l := s.length()
	_, unicode := s.(unicodeString)
	var asciiBuf []byte
	var unicodeBuf []uint16
	if unicode ***REMOVED***
		unicodeBuf = make([]uint16, 0, l)
	***REMOVED*** else ***REMOVED***
		asciiBuf = make([]byte, 0, l)
	***REMOVED***
	for i := int64(0); i < l; ***REMOVED***
		r := s.charAt(i)
		if r == '%' ***REMOVED***
			if i <= l-6 && s.charAt(i+1) == 'u' ***REMOVED***
				c0 := s.charAt(i + 2)
				c1 := s.charAt(i + 3)
				c2 := s.charAt(i + 4)
				c3 := s.charAt(i + 5)
				if c0 <= 0xff && ishex(byte(c0)) &&
					c1 <= 0xff && ishex(byte(c1)) &&
					c2 <= 0xff && ishex(byte(c2)) &&
					c3 <= 0xff && ishex(byte(c3)) ***REMOVED***
					r = rune(unhex(byte(c0)))<<12 |
						rune(unhex(byte(c1)))<<8 |
						rune(unhex(byte(c2)))<<4 |
						rune(unhex(byte(c3)))
					i += 5
					goto out
				***REMOVED***
			***REMOVED***
			if i <= l-3 ***REMOVED***
				c0 := s.charAt(i + 1)
				c1 := s.charAt(i + 2)
				if c0 <= 0xff && ishex(byte(c0)) &&
					c1 <= 0xff && ishex(byte(c1)) ***REMOVED***
					r = rune(unhex(byte(c0))<<4 | unhex(byte(c1)))
					i += 2
				***REMOVED***
			***REMOVED***
		***REMOVED***
	out:
		if r >= utf8.RuneSelf && !unicode ***REMOVED***
			unicodeBuf = make([]uint16, 0, l)
			for _, b := range asciiBuf ***REMOVED***
				unicodeBuf = append(unicodeBuf, uint16(b))
			***REMOVED***
			asciiBuf = nil
			unicode = true
		***REMOVED***
		if unicode ***REMOVED***
			unicodeBuf = append(unicodeBuf, uint16(r))
		***REMOVED*** else ***REMOVED***
			asciiBuf = append(asciiBuf, byte(r))
		***REMOVED***
		i++
	***REMOVED***
	if unicode ***REMOVED***
		return unicodeString(unicodeBuf)
	***REMOVED***

	return asciiString(asciiBuf)
***REMOVED***

func (r *Runtime) initGlobalObject() ***REMOVED***
	o := r.globalObject.self
	o._putProp("NaN", _NaN, false, false, false)
	o._putProp("undefined", _undefined, false, false, false)
	o._putProp("Infinity", _positiveInf, false, false, false)

	o._putProp("isNaN", r.newNativeFunc(r.builtin_isNaN, nil, "isNaN", nil, 1), true, false, true)
	o._putProp("parseInt", r.newNativeFunc(r.builtin_parseInt, nil, "parseInt", nil, 2), true, false, true)
	o._putProp("parseFloat", r.newNativeFunc(r.builtin_parseFloat, nil, "parseFloat", nil, 1), true, false, true)
	o._putProp("isFinite", r.newNativeFunc(r.builtin_isFinite, nil, "isFinite", nil, 1), true, false, true)
	o._putProp("decodeURI", r.newNativeFunc(r.builtin_decodeURI, nil, "decodeURI", nil, 1), true, false, true)
	o._putProp("decodeURIComponent", r.newNativeFunc(r.builtin_decodeURIComponent, nil, "decodeURIComponent", nil, 1), true, false, true)
	o._putProp("encodeURI", r.newNativeFunc(r.builtin_encodeURI, nil, "encodeURI", nil, 1), true, false, true)
	o._putProp("encodeURIComponent", r.newNativeFunc(r.builtin_encodeURIComponent, nil, "encodeURIComponent", nil, 1), true, false, true)
	o._putProp("escape", r.newNativeFunc(r.builtin_escape, nil, "escape", nil, 1), true, false, true)
	o._putProp("unescape", r.newNativeFunc(r.builtin_unescape, nil, "unescape", nil, 1), true, false, true)

	o._putProp("toString", r.newNativeFunc(func(FunctionCall) Value ***REMOVED***
		return stringGlobalObject
	***REMOVED***, nil, "toString", nil, 0), false, false, false)

	// TODO: Annex B

***REMOVED***

func digitVal(d byte) int ***REMOVED***
	var v byte
	switch ***REMOVED***
	case '0' <= d && d <= '9':
		v = d - '0'
	case 'a' <= d && d <= 'z':
		v = d - 'a' + 10
	case 'A' <= d && d <= 'Z':
		v = d - 'A' + 10
	default:
		return 36
	***REMOVED***
	return int(v)
***REMOVED***

// ECMAScript compatible version of strconv.ParseInt
func parseInt(s string, base int) (Value, error) ***REMOVED***
	var n int64
	var err error
	var cutoff, maxVal int64
	var sign bool
	i := 0

	if len(s) < 1 ***REMOVED***
		err = strconv.ErrSyntax
		goto Error
	***REMOVED***

	switch s[0] ***REMOVED***
	case '-':
		sign = true
		s = s[1:]
	case '+':
		s = s[1:]
	***REMOVED***

	if len(s) < 1 ***REMOVED***
		err = strconv.ErrSyntax
		goto Error
	***REMOVED***

	// Look for hex prefix.
	if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') ***REMOVED***
		if base == 0 || base == 16 ***REMOVED***
			base = 16
			s = s[2:]
		***REMOVED***
	***REMOVED***

	switch ***REMOVED***
	case len(s) < 1:
		err = strconv.ErrSyntax
		goto Error

	case 2 <= base && base <= 36:
	// valid base; nothing to do

	case base == 0:
		// Look for hex prefix.
		switch ***REMOVED***
		case s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X'):
			if len(s) < 3 ***REMOVED***
				err = strconv.ErrSyntax
				goto Error
			***REMOVED***
			base = 16
			s = s[2:]
		default:
			base = 10
		***REMOVED***

	default:
		err = errors.New("invalid base " + strconv.Itoa(base))
		goto Error
	***REMOVED***

	// Cutoff is the smallest number such that cutoff*base > maxInt64.
	// Use compile-time constants for common cases.
	switch base ***REMOVED***
	case 10:
		cutoff = math.MaxInt64/10 + 1
	case 16:
		cutoff = math.MaxInt64/16 + 1
	default:
		cutoff = math.MaxInt64/int64(base) + 1
	***REMOVED***

	maxVal = math.MaxInt64
	for ; i < len(s); i++ ***REMOVED***
		if n >= cutoff ***REMOVED***
			// n*base overflows
			return parseLargeInt(float64(n), s[i:], base, sign)
		***REMOVED***
		v := digitVal(s[i])
		if v >= base ***REMOVED***
			break
		***REMOVED***
		n *= int64(base)

		n1 := n + int64(v)
		if n1 < n || n1 > maxVal ***REMOVED***
			// n+v overflows
			return parseLargeInt(float64(n)+float64(v), s[i+1:], base, sign)
		***REMOVED***
		n = n1
	***REMOVED***

	if i == 0 ***REMOVED***
		err = strconv.ErrSyntax
		goto Error
	***REMOVED***

	if sign ***REMOVED***
		n = -n
	***REMOVED***
	return intToValue(n), nil

Error:
	return _NaN, err
***REMOVED***

func parseLargeInt(n float64, s string, base int, sign bool) (Value, error) ***REMOVED***
	i := 0
	b := float64(base)
	for ; i < len(s); i++ ***REMOVED***
		v := digitVal(s[i])
		if v >= base ***REMOVED***
			break
		***REMOVED***
		n = n*b + float64(v)
	***REMOVED***
	if sign ***REMOVED***
		n = -n
	***REMOVED***
	// We know it can't be represented as int, so use valueFloat instead of floatToValue
	return valueFloat(n), nil
***REMOVED***

var (
	uriUnescaped             [256]bool
	uriReserved              [256]bool
	uriReservedHash          [256]bool
	uriReservedUnescapedHash [256]bool
	emptyEscapeSet           [256]bool
)

func init() ***REMOVED***
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.!~*'()" ***REMOVED***
		uriUnescaped[c] = true
	***REMOVED***

	for _, c := range ";/?:@&=+$," ***REMOVED***
		uriReserved[c] = true
	***REMOVED***

	for i := 0; i < 256; i++ ***REMOVED***
		if uriUnescaped[i] || uriReserved[i] ***REMOVED***
			uriReservedUnescapedHash[i] = true
		***REMOVED***
		uriReservedHash[i] = uriReserved[i]
	***REMOVED***
	uriReservedUnescapedHash['#'] = true
	uriReservedHash['#'] = true
***REMOVED***
