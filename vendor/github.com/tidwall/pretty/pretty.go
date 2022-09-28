package pretty

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
)

// Options is Pretty options
type Options struct ***REMOVED***
	// Width is an max column width for single line arrays
	// Default is 80
	Width int
	// Prefix is a prefix for all lines
	// Default is an empty string
	Prefix string
	// Indent is the nested indentation
	// Default is two spaces
	Indent string
	// SortKeys will sort the keys alphabetically
	// Default is false
	SortKeys bool
***REMOVED***

// DefaultOptions is the default options for pretty formats.
var DefaultOptions = &Options***REMOVED***Width: 80, Prefix: "", Indent: "  ", SortKeys: false***REMOVED***

// Pretty converts the input json into a more human readable format where each
// element is on it's own line with clear indentation.
func Pretty(json []byte) []byte ***REMOVED*** return PrettyOptions(json, nil) ***REMOVED***

// PrettyOptions is like Pretty but with customized options.
func PrettyOptions(json []byte, opts *Options) []byte ***REMOVED***
	if opts == nil ***REMOVED***
		opts = DefaultOptions
	***REMOVED***
	buf := make([]byte, 0, len(json))
	if len(opts.Prefix) != 0 ***REMOVED***
		buf = append(buf, opts.Prefix...)
	***REMOVED***
	buf, _, _, _ = appendPrettyAny(buf, json, 0, true,
		opts.Width, opts.Prefix, opts.Indent, opts.SortKeys,
		0, 0, -1)
	if len(buf) > 0 ***REMOVED***
		buf = append(buf, '\n')
	***REMOVED***
	return buf
***REMOVED***

// Ugly removes insignificant space characters from the input json byte slice
// and returns the compacted result.
func Ugly(json []byte) []byte ***REMOVED***
	buf := make([]byte, 0, len(json))
	return ugly(buf, json)
***REMOVED***

// UglyInPlace removes insignificant space characters from the input json
// byte slice and returns the compacted result. This method reuses the
// input json buffer to avoid allocations. Do not use the original bytes
// slice upon return.
func UglyInPlace(json []byte) []byte ***REMOVED*** return ugly(json, json) ***REMOVED***

func ugly(dst, src []byte) []byte ***REMOVED***
	dst = dst[:0]
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] > ' ' ***REMOVED***
			dst = append(dst, src[i])
			if src[i] == '"' ***REMOVED***
				for i = i + 1; i < len(src); i++ ***REMOVED***
					dst = append(dst, src[i])
					if src[i] == '"' ***REMOVED***
						j := i - 1
						for ; ; j-- ***REMOVED***
							if src[j] != '\\' ***REMOVED***
								break
							***REMOVED***
						***REMOVED***
						if (j-i)%2 != 0 ***REMOVED***
							break
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

func isNaNOrInf(src []byte) bool ***REMOVED***
	return src[0] == 'i' || //Inf
		src[0] == 'I' || // inf
		src[0] == '+' || // +Inf
		src[0] == 'N' || // Nan
		(src[0] == 'n' && len(src) > 1 && src[1] != 'u') // nan
***REMOVED***

func appendPrettyAny(buf, json []byte, i int, pretty bool, width int, prefix, indent string, sortkeys bool, tabs, nl, max int) ([]byte, int, int, bool) ***REMOVED***
	for ; i < len(json); i++ ***REMOVED***
		if json[i] <= ' ' ***REMOVED***
			continue
		***REMOVED***
		if json[i] == '"' ***REMOVED***
			return appendPrettyString(buf, json, i, nl)
		***REMOVED***

		if (json[i] >= '0' && json[i] <= '9') || json[i] == '-' || isNaNOrInf(json[i:]) ***REMOVED***
			return appendPrettyNumber(buf, json, i, nl)
		***REMOVED***
		if json[i] == '***REMOVED***' ***REMOVED***
			return appendPrettyObject(buf, json, i, '***REMOVED***', '***REMOVED***', pretty, width, prefix, indent, sortkeys, tabs, nl, max)
		***REMOVED***
		if json[i] == '[' ***REMOVED***
			return appendPrettyObject(buf, json, i, '[', ']', pretty, width, prefix, indent, sortkeys, tabs, nl, max)
		***REMOVED***
		switch json[i] ***REMOVED***
		case 't':
			return append(buf, 't', 'r', 'u', 'e'), i + 4, nl, true
		case 'f':
			return append(buf, 'f', 'a', 'l', 's', 'e'), i + 5, nl, true
		case 'n':
			return append(buf, 'n', 'u', 'l', 'l'), i + 4, nl, true
		***REMOVED***
	***REMOVED***
	return buf, i, nl, true
***REMOVED***

type pair struct ***REMOVED***
	kstart, kend int
	vstart, vend int
***REMOVED***

type byKeyVal struct ***REMOVED***
	sorted bool
	json   []byte
	buf    []byte
	pairs  []pair
***REMOVED***

func (arr *byKeyVal) Len() int ***REMOVED***
	return len(arr.pairs)
***REMOVED***
func (arr *byKeyVal) Less(i, j int) bool ***REMOVED***
	if arr.isLess(i, j, byKey) ***REMOVED***
		return true
	***REMOVED***
	if arr.isLess(j, i, byKey) ***REMOVED***
		return false
	***REMOVED***
	return arr.isLess(i, j, byVal)
***REMOVED***
func (arr *byKeyVal) Swap(i, j int) ***REMOVED***
	arr.pairs[i], arr.pairs[j] = arr.pairs[j], arr.pairs[i]
	arr.sorted = true
***REMOVED***

type byKind int

const (
	byKey byKind = 0
	byVal byKind = 1
)

type jtype int

const (
	jnull jtype = iota
	jfalse
	jnumber
	jstring
	jtrue
	jjson
)

func getjtype(v []byte) jtype ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return jnull
	***REMOVED***
	switch v[0] ***REMOVED***
	case '"':
		return jstring
	case 'f':
		return jfalse
	case 't':
		return jtrue
	case 'n':
		return jnull
	case '[', '***REMOVED***':
		return jjson
	default:
		return jnumber
	***REMOVED***
***REMOVED***

func (arr *byKeyVal) isLess(i, j int, kind byKind) bool ***REMOVED***
	k1 := arr.json[arr.pairs[i].kstart:arr.pairs[i].kend]
	k2 := arr.json[arr.pairs[j].kstart:arr.pairs[j].kend]
	var v1, v2 []byte
	if kind == byKey ***REMOVED***
		v1 = k1
		v2 = k2
	***REMOVED*** else ***REMOVED***
		v1 = bytes.TrimSpace(arr.buf[arr.pairs[i].vstart:arr.pairs[i].vend])
		v2 = bytes.TrimSpace(arr.buf[arr.pairs[j].vstart:arr.pairs[j].vend])
		if len(v1) >= len(k1)+1 ***REMOVED***
			v1 = bytes.TrimSpace(v1[len(k1)+1:])
		***REMOVED***
		if len(v2) >= len(k2)+1 ***REMOVED***
			v2 = bytes.TrimSpace(v2[len(k2)+1:])
		***REMOVED***
	***REMOVED***
	t1 := getjtype(v1)
	t2 := getjtype(v2)
	if t1 < t2 ***REMOVED***
		return true
	***REMOVED***
	if t1 > t2 ***REMOVED***
		return false
	***REMOVED***
	if t1 == jstring ***REMOVED***
		s1 := parsestr(v1)
		s2 := parsestr(v2)
		return string(s1) < string(s2)
	***REMOVED***
	if t1 == jnumber ***REMOVED***
		n1, _ := strconv.ParseFloat(string(v1), 64)
		n2, _ := strconv.ParseFloat(string(v2), 64)
		return n1 < n2
	***REMOVED***
	return string(v1) < string(v2)

***REMOVED***

func parsestr(s []byte) []byte ***REMOVED***
	for i := 1; i < len(s); i++ ***REMOVED***
		if s[i] == '\\' ***REMOVED***
			var str string
			json.Unmarshal(s, &str)
			return []byte(str)
		***REMOVED***
		if s[i] == '"' ***REMOVED***
			return s[1:i]
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func appendPrettyObject(buf, json []byte, i int, open, close byte, pretty bool, width int, prefix, indent string, sortkeys bool, tabs, nl, max int) ([]byte, int, int, bool) ***REMOVED***
	var ok bool
	if width > 0 ***REMOVED***
		if pretty && open == '[' && max == -1 ***REMOVED***
			// here we try to create a single line array
			max := width - (len(buf) - nl)
			if max > 3 ***REMOVED***
				s1, s2 := len(buf), i
				buf, i, _, ok = appendPrettyObject(buf, json, i, '[', ']', false, width, prefix, "", sortkeys, 0, 0, max)
				if ok && len(buf)-s1 <= max ***REMOVED***
					return buf, i, nl, true
				***REMOVED***
				buf = buf[:s1]
				i = s2
			***REMOVED***
		***REMOVED*** else if max != -1 && open == '***REMOVED***' ***REMOVED***
			return buf, i, nl, false
		***REMOVED***
	***REMOVED***
	buf = append(buf, open)
	i++
	var pairs []pair
	if open == '***REMOVED***' && sortkeys ***REMOVED***
		pairs = make([]pair, 0, 8)
	***REMOVED***
	var n int
	for ; i < len(json); i++ ***REMOVED***
		if json[i] <= ' ' ***REMOVED***
			continue
		***REMOVED***
		if json[i] == close ***REMOVED***
			if pretty ***REMOVED***
				if open == '***REMOVED***' && sortkeys ***REMOVED***
					buf = sortPairs(json, buf, pairs)
				***REMOVED***
				if n > 0 ***REMOVED***
					nl = len(buf)
					if buf[nl-1] == ' ' ***REMOVED***
						buf[nl-1] = '\n'
					***REMOVED*** else ***REMOVED***
						buf = append(buf, '\n')
					***REMOVED***
				***REMOVED***
				if buf[len(buf)-1] != open ***REMOVED***
					buf = appendTabs(buf, prefix, indent, tabs)
				***REMOVED***
			***REMOVED***
			buf = append(buf, close)
			return buf, i + 1, nl, open != '***REMOVED***'
		***REMOVED***
		if open == '[' || json[i] == '"' ***REMOVED***
			if n > 0 ***REMOVED***
				buf = append(buf, ',')
				if width != -1 && open == '[' ***REMOVED***
					buf = append(buf, ' ')
				***REMOVED***
			***REMOVED***
			var p pair
			if pretty ***REMOVED***
				nl = len(buf)
				if buf[nl-1] == ' ' ***REMOVED***
					buf[nl-1] = '\n'
				***REMOVED*** else ***REMOVED***
					buf = append(buf, '\n')
				***REMOVED***
				if open == '***REMOVED***' && sortkeys ***REMOVED***
					p.kstart = i
					p.vstart = len(buf)
				***REMOVED***
				buf = appendTabs(buf, prefix, indent, tabs+1)
			***REMOVED***
			if open == '***REMOVED***' ***REMOVED***
				buf, i, nl, _ = appendPrettyString(buf, json, i, nl)
				if sortkeys ***REMOVED***
					p.kend = i
				***REMOVED***
				buf = append(buf, ':')
				if pretty ***REMOVED***
					buf = append(buf, ' ')
				***REMOVED***
			***REMOVED***
			buf, i, nl, ok = appendPrettyAny(buf, json, i, pretty, width, prefix, indent, sortkeys, tabs+1, nl, max)
			if max != -1 && !ok ***REMOVED***
				return buf, i, nl, false
			***REMOVED***
			if pretty && open == '***REMOVED***' && sortkeys ***REMOVED***
				p.vend = len(buf)
				if p.kstart > p.kend || p.vstart > p.vend ***REMOVED***
					// bad data. disable sorting
					sortkeys = false
				***REMOVED*** else ***REMOVED***
					pairs = append(pairs, p)
				***REMOVED***
			***REMOVED***
			i--
			n++
		***REMOVED***
	***REMOVED***
	return buf, i, nl, open != '***REMOVED***'
***REMOVED***
func sortPairs(json, buf []byte, pairs []pair) []byte ***REMOVED***
	if len(pairs) == 0 ***REMOVED***
		return buf
	***REMOVED***
	vstart := pairs[0].vstart
	vend := pairs[len(pairs)-1].vend
	arr := byKeyVal***REMOVED***false, json, buf, pairs***REMOVED***
	sort.Stable(&arr)
	if !arr.sorted ***REMOVED***
		return buf
	***REMOVED***
	nbuf := make([]byte, 0, vend-vstart)
	for i, p := range pairs ***REMOVED***
		nbuf = append(nbuf, buf[p.vstart:p.vend]...)
		if i < len(pairs)-1 ***REMOVED***
			nbuf = append(nbuf, ',')
			nbuf = append(nbuf, '\n')
		***REMOVED***
	***REMOVED***
	return append(buf[:vstart], nbuf...)
***REMOVED***

func appendPrettyString(buf, json []byte, i, nl int) ([]byte, int, int, bool) ***REMOVED***
	s := i
	i++
	for ; i < len(json); i++ ***REMOVED***
		if json[i] == '"' ***REMOVED***
			var sc int
			for j := i - 1; j > s; j-- ***REMOVED***
				if json[j] == '\\' ***REMOVED***
					sc++
				***REMOVED*** else ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			if sc%2 == 1 ***REMOVED***
				continue
			***REMOVED***
			i++
			break
		***REMOVED***
	***REMOVED***
	return append(buf, json[s:i]...), i, nl, true
***REMOVED***

func appendPrettyNumber(buf, json []byte, i, nl int) ([]byte, int, int, bool) ***REMOVED***
	s := i
	i++
	for ; i < len(json); i++ ***REMOVED***
		if json[i] <= ' ' || json[i] == ',' || json[i] == ':' || json[i] == ']' || json[i] == '***REMOVED***' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return append(buf, json[s:i]...), i, nl, true
***REMOVED***

func appendTabs(buf []byte, prefix, indent string, tabs int) []byte ***REMOVED***
	if len(prefix) != 0 ***REMOVED***
		buf = append(buf, prefix...)
	***REMOVED***
	if len(indent) == 2 && indent[0] == ' ' && indent[1] == ' ' ***REMOVED***
		for i := 0; i < tabs; i++ ***REMOVED***
			buf = append(buf, ' ', ' ')
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := 0; i < tabs; i++ ***REMOVED***
			buf = append(buf, indent...)
		***REMOVED***
	***REMOVED***
	return buf
***REMOVED***

// Style is the color style
type Style struct ***REMOVED***
	Key, String, Number [2]string
	True, False, Null   [2]string
	Escape              [2]string
	Append              func(dst []byte, c byte) []byte
***REMOVED***

func hexp(p byte) byte ***REMOVED***
	switch ***REMOVED***
	case p < 10:
		return p + '0'
	default:
		return (p - 10) + 'a'
	***REMOVED***
***REMOVED***

// TerminalStyle is for terminals
var TerminalStyle *Style

func init() ***REMOVED***
	TerminalStyle = &Style***REMOVED***
		Key:    [2]string***REMOVED***"\x1B[94m", "\x1B[0m"***REMOVED***,
		String: [2]string***REMOVED***"\x1B[92m", "\x1B[0m"***REMOVED***,
		Number: [2]string***REMOVED***"\x1B[93m", "\x1B[0m"***REMOVED***,
		True:   [2]string***REMOVED***"\x1B[96m", "\x1B[0m"***REMOVED***,
		False:  [2]string***REMOVED***"\x1B[96m", "\x1B[0m"***REMOVED***,
		Null:   [2]string***REMOVED***"\x1B[91m", "\x1B[0m"***REMOVED***,
		Escape: [2]string***REMOVED***"\x1B[35m", "\x1B[0m"***REMOVED***,
		Append: func(dst []byte, c byte) []byte ***REMOVED***
			if c < ' ' && (c != '\r' && c != '\n' && c != '\t' && c != '\v') ***REMOVED***
				dst = append(dst, "\\u00"...)
				dst = append(dst, hexp((c>>4)&0xF))
				return append(dst, hexp((c)&0xF))
			***REMOVED***
			return append(dst, c)
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Color will colorize the json. The style parma is used for customizing
// the colors. Passing nil to the style param will use the default
// TerminalStyle.
func Color(src []byte, style *Style) []byte ***REMOVED***
	if style == nil ***REMOVED***
		style = TerminalStyle
	***REMOVED***
	apnd := style.Append
	if apnd == nil ***REMOVED***
		apnd = func(dst []byte, c byte) []byte ***REMOVED***
			return append(dst, c)
		***REMOVED***
	***REMOVED***
	type stackt struct ***REMOVED***
		kind byte
		key  bool
	***REMOVED***
	var dst []byte
	var stack []stackt
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] == '"' ***REMOVED***
			key := len(stack) > 0 && stack[len(stack)-1].key
			if key ***REMOVED***
				dst = append(dst, style.Key[0]...)
			***REMOVED*** else ***REMOVED***
				dst = append(dst, style.String[0]...)
			***REMOVED***
			dst = apnd(dst, '"')
			esc := false
			uesc := 0
			for i = i + 1; i < len(src); i++ ***REMOVED***
				if src[i] == '\\' ***REMOVED***
					if key ***REMOVED***
						dst = append(dst, style.Key[1]...)
					***REMOVED*** else ***REMOVED***
						dst = append(dst, style.String[1]...)
					***REMOVED***
					dst = append(dst, style.Escape[0]...)
					dst = apnd(dst, src[i])
					esc = true
					if i+1 < len(src) && src[i+1] == 'u' ***REMOVED***
						uesc = 5
					***REMOVED*** else ***REMOVED***
						uesc = 1
					***REMOVED***
				***REMOVED*** else if esc ***REMOVED***
					dst = apnd(dst, src[i])
					if uesc == 1 ***REMOVED***
						esc = false
						dst = append(dst, style.Escape[1]...)
						if key ***REMOVED***
							dst = append(dst, style.Key[0]...)
						***REMOVED*** else ***REMOVED***
							dst = append(dst, style.String[0]...)
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						uesc--
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					dst = apnd(dst, src[i])
				***REMOVED***
				if src[i] == '"' ***REMOVED***
					j := i - 1
					for ; ; j-- ***REMOVED***
						if src[j] != '\\' ***REMOVED***
							break
						***REMOVED***
					***REMOVED***
					if (j-i)%2 != 0 ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if esc ***REMOVED***
				dst = append(dst, style.Escape[1]...)
			***REMOVED*** else if key ***REMOVED***
				dst = append(dst, style.Key[1]...)
			***REMOVED*** else ***REMOVED***
				dst = append(dst, style.String[1]...)
			***REMOVED***
		***REMOVED*** else if src[i] == '***REMOVED***' || src[i] == '[' ***REMOVED***
			stack = append(stack, stackt***REMOVED***src[i], src[i] == '***REMOVED***'***REMOVED***)
			dst = apnd(dst, src[i])
		***REMOVED*** else if (src[i] == '***REMOVED***' || src[i] == ']') && len(stack) > 0 ***REMOVED***
			stack = stack[:len(stack)-1]
			dst = apnd(dst, src[i])
		***REMOVED*** else if (src[i] == ':' || src[i] == ',') && len(stack) > 0 && stack[len(stack)-1].kind == '***REMOVED***' ***REMOVED***
			stack[len(stack)-1].key = !stack[len(stack)-1].key
			dst = apnd(dst, src[i])
		***REMOVED*** else ***REMOVED***
			var kind byte
			if (src[i] >= '0' && src[i] <= '9') || src[i] == '-' || isNaNOrInf(src[i:]) ***REMOVED***
				kind = '0'
				dst = append(dst, style.Number[0]...)
			***REMOVED*** else if src[i] == 't' ***REMOVED***
				kind = 't'
				dst = append(dst, style.True[0]...)
			***REMOVED*** else if src[i] == 'f' ***REMOVED***
				kind = 'f'
				dst = append(dst, style.False[0]...)
			***REMOVED*** else if src[i] == 'n' ***REMOVED***
				kind = 'n'
				dst = append(dst, style.Null[0]...)
			***REMOVED*** else ***REMOVED***
				dst = apnd(dst, src[i])
			***REMOVED***
			if kind != 0 ***REMOVED***
				for ; i < len(src); i++ ***REMOVED***
					if src[i] <= ' ' || src[i] == ',' || src[i] == ':' || src[i] == ']' || src[i] == '***REMOVED***' ***REMOVED***
						i--
						break
					***REMOVED***
					dst = apnd(dst, src[i])
				***REMOVED***
				if kind == '0' ***REMOVED***
					dst = append(dst, style.Number[1]...)
				***REMOVED*** else if kind == 't' ***REMOVED***
					dst = append(dst, style.True[1]...)
				***REMOVED*** else if kind == 'f' ***REMOVED***
					dst = append(dst, style.False[1]...)
				***REMOVED*** else if kind == 'n' ***REMOVED***
					dst = append(dst, style.Null[1]...)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// Spec strips out comments and trailing commas and convert the input to a
// valid JSON per the official spec: https://tools.ietf.org/html/rfc8259
//
// The resulting JSON will always be the same length as the input and it will
// include all of the same line breaks at matching offsets. This is to ensure
// the result can be later processed by a external parser and that that
// parser will report messages or errors with the correct offsets.
func Spec(src []byte) []byte ***REMOVED***
	return spec(src, nil)
***REMOVED***

// SpecInPlace is the same as Spec, but this method reuses the input json
// buffer to avoid allocations. Do not use the original bytes slice upon return.
func SpecInPlace(src []byte) []byte ***REMOVED***
	return spec(src, src)
***REMOVED***

func spec(src, dst []byte) []byte ***REMOVED***
	dst = dst[:0]
	for i := 0; i < len(src); i++ ***REMOVED***
		if src[i] == '/' ***REMOVED***
			if i < len(src)-1 ***REMOVED***
				if src[i+1] == '/' ***REMOVED***
					dst = append(dst, ' ', ' ')
					i += 2
					for ; i < len(src); i++ ***REMOVED***
						if src[i] == '\n' ***REMOVED***
							dst = append(dst, '\n')
							break
						***REMOVED*** else if src[i] == '\t' || src[i] == '\r' ***REMOVED***
							dst = append(dst, src[i])
						***REMOVED*** else ***REMOVED***
							dst = append(dst, ' ')
						***REMOVED***
					***REMOVED***
					continue
				***REMOVED***
				if src[i+1] == '*' ***REMOVED***
					dst = append(dst, ' ', ' ')
					i += 2
					for ; i < len(src)-1; i++ ***REMOVED***
						if src[i] == '*' && src[i+1] == '/' ***REMOVED***
							dst = append(dst, ' ', ' ')
							i++
							break
						***REMOVED*** else if src[i] == '\n' || src[i] == '\t' ||
							src[i] == '\r' ***REMOVED***
							dst = append(dst, src[i])
						***REMOVED*** else ***REMOVED***
							dst = append(dst, ' ')
						***REMOVED***
					***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***
		dst = append(dst, src[i])
		if src[i] == '"' ***REMOVED***
			for i = i + 1; i < len(src); i++ ***REMOVED***
				dst = append(dst, src[i])
				if src[i] == '"' ***REMOVED***
					j := i - 1
					for ; ; j-- ***REMOVED***
						if src[j] != '\\' ***REMOVED***
							break
						***REMOVED***
					***REMOVED***
					if (j-i)%2 != 0 ***REMOVED***
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if src[i] == '***REMOVED***' || src[i] == ']' ***REMOVED***
			for j := len(dst) - 2; j >= 0; j-- ***REMOVED***
				if dst[j] <= ' ' ***REMOVED***
					continue
				***REMOVED***
				if dst[j] == ',' ***REMOVED***
					dst[j] = ' '
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***
