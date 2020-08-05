package goja

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/dop251/goja/unistring"
	"io"
	"regexp"
	"sort"
	"strings"
	"unicode/utf16"
)

type regexp2Wrapper regexp2.Regexp
type regexpWrapper regexp.Regexp

type positionMapItem struct ***REMOVED***
	src, dst int
***REMOVED***
type positionMap []positionMapItem

func (m positionMap) get(src int) int ***REMOVED***
	if src == 0 ***REMOVED***
		return 0
	***REMOVED***
	res := sort.Search(len(m), func(n int) bool ***REMOVED*** return m[n].src >= src ***REMOVED***)
	if res >= len(m) || m[res].src != src ***REMOVED***
		panic("index not found")
	***REMOVED***
	return m[res].dst
***REMOVED***

type arrayRuneReader struct ***REMOVED***
	runes []rune
	pos   int
***REMOVED***

func (rd *arrayRuneReader) ReadRune() (r rune, size int, err error) ***REMOVED***
	if rd.pos < len(rd.runes) ***REMOVED***
		r = rd.runes[rd.pos]
		size = 1
		rd.pos++
	***REMOVED*** else ***REMOVED***
		err = io.EOF
	***REMOVED***
	return
***REMOVED***

type regexpPattern struct ***REMOVED***
	src string

	global, ignoreCase, multiline, sticky, unicode bool

	regexpWrapper  *regexpWrapper
	regexp2Wrapper *regexp2Wrapper
***REMOVED***

func compileRegexp2(src string, multiline, ignoreCase bool) (*regexp2Wrapper, error) ***REMOVED***
	var opts regexp2.RegexOptions = regexp2.ECMAScript
	if multiline ***REMOVED***
		opts |= regexp2.Multiline
	***REMOVED***
	if ignoreCase ***REMOVED***
		opts |= regexp2.IgnoreCase
	***REMOVED***
	regexp2Pattern, err1 := regexp2.Compile(src, opts)
	if err1 != nil ***REMOVED***
		return nil, fmt.Errorf("Invalid regular expression (regexp2): %s (%v)", src, err1)
	***REMOVED***

	return (*regexp2Wrapper)(regexp2Pattern), nil
***REMOVED***

func (p *regexpPattern) createRegexp2() ***REMOVED***
	if p.regexp2Wrapper != nil ***REMOVED***
		return
	***REMOVED***
	rx, err := compileRegexp2(p.src, p.multiline, p.ignoreCase)
	if err != nil ***REMOVED***
		// At this point the regexp should have been successfully converted to re2, if it fails now, it's a bug.
		panic(err)
	***REMOVED***
	p.regexp2Wrapper = rx
***REMOVED***

func buildUTF8PosMap(s valueString) (positionMap, string) ***REMOVED***
	pm := make(positionMap, 0, s.length())
	rd := s.reader(0)
	sPos, utf8Pos := 0, 0
	var sb strings.Builder
	for ***REMOVED***
		r, size, err := rd.ReadRune()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			// the string contains invalid UTF-16, bailing out
			return nil, ""
		***REMOVED***
		utf8Size, _ := sb.WriteRune(r)
		sPos += size
		utf8Pos += utf8Size
		pm = append(pm, positionMapItem***REMOVED***src: utf8Pos, dst: sPos***REMOVED***)
	***REMOVED***
	return pm, sb.String()
***REMOVED***

func (p *regexpPattern) findSubmatchIndex(s valueString, start int) []int ***REMOVED***
	if p.regexpWrapper == nil ***REMOVED***
		return p.regexp2Wrapper.findSubmatchIndex(s, start, p.unicode)
	***REMOVED***
	if start != 0 ***REMOVED***
		// Unfortunately Go's regexp library does not allow starting from an arbitrary position.
		// If we just drop the first _start_ characters of the string the assertions (^, $, \b and \B) will not
		// work correctly.
		p.createRegexp2()
		return p.regexp2Wrapper.findSubmatchIndex(s, start, p.unicode)
	***REMOVED***
	return p.regexpWrapper.findSubmatchIndex(s, p.unicode)
***REMOVED***

func (p *regexpPattern) findAllSubmatchIndex(s valueString, start int, limit int, sticky bool) [][]int ***REMOVED***
	if p.regexpWrapper == nil ***REMOVED***
		return p.regexp2Wrapper.findAllSubmatchIndex(s, start, limit, sticky, p.unicode)
	***REMOVED***
	if start == 0 ***REMOVED***
		if s, ok := s.(asciiString); ok ***REMOVED***
			return p.regexpWrapper.findAllSubmatchIndex(s.String(), limit, sticky)
		***REMOVED***
		if limit == 1 ***REMOVED***
			result := p.regexpWrapper.findSubmatchIndex(s, p.unicode)
			if result == nil ***REMOVED***
				return nil
			***REMOVED***
			return [][]int***REMOVED***result***REMOVED***
		***REMOVED***
		// Unfortunately Go's regexp library lacks FindAllReaderSubmatchIndex(), so we have to use a UTF-8 string as an
		// input.
		if p.unicode ***REMOVED***
			// Try to convert s to UTF-8. If it does not contain any invalid UTF-16 we can do the matching in UTF-8.
			pm, str := buildUTF8PosMap(s)
			if pm != nil ***REMOVED***
				res := p.regexpWrapper.findAllSubmatchIndex(str, limit, sticky)
				for _, result := range res ***REMOVED***
					for i, idx := range result ***REMOVED***
						result[i] = pm.get(idx)
					***REMOVED***
				***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
	***REMOVED***

	p.createRegexp2()
	return p.regexp2Wrapper.findAllSubmatchIndex(s, start, limit, sticky, p.unicode)
***REMOVED***

type regexpObject struct ***REMOVED***
	baseObject
	pattern *regexpPattern
	source  valueString

	standard bool
***REMOVED***

func (r *regexp2Wrapper) findSubmatchIndex(s valueString, start int, fullUnicode bool) (result []int) ***REMOVED***
	if fullUnicode ***REMOVED***
		return r.findSubmatchIndexUnicode(s, start)
	***REMOVED***
	return r.findSubmatchIndexUTF16(s, start)
***REMOVED***

func (r *regexp2Wrapper) findSubmatchIndexUTF16(s valueString, start int) (result []int) ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	match, err := wrapped.FindRunesMatchStartingAt(s.utf16Runes(), start)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if match == nil ***REMOVED***
		return
	***REMOVED***
	groups := match.Groups()

	result = make([]int, 0, len(groups)<<1)
	for _, group := range groups ***REMOVED***
		if len(group.Captures) > 0 ***REMOVED***
			result = append(result, group.Index, group.Index+group.Length)
		***REMOVED*** else ***REMOVED***
			result = append(result, -1, 0)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (r *regexp2Wrapper) findSubmatchIndexUnicode(s valueString, start int) (result []int) ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	posMap, runes, mappedStart := buildPosMap(&lenientUtf16Decoder***REMOVED***utf16Reader: s.utf16Reader(0)***REMOVED***, s.length(), start)
	match, err := wrapped.FindRunesMatchStartingAt(runes, mappedStart)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if match == nil ***REMOVED***
		return
	***REMOVED***
	groups := match.Groups()

	result = make([]int, 0, len(groups)<<1)
	for _, group := range groups ***REMOVED***
		if len(group.Captures) > 0 ***REMOVED***
			result = append(result, posMap[group.Index], posMap[group.Index+group.Length])
		***REMOVED*** else ***REMOVED***
			result = append(result, -1, 0)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (r *regexp2Wrapper) findAllSubmatchIndexUTF16(s valueString, start, limit int, sticky bool) [][]int ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	runes := s.utf16Runes()
	match, err := wrapped.FindRunesMatchStartingAt(runes, start)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	if limit < 0 ***REMOVED***
		limit = len(runes) + 1
	***REMOVED***
	results := make([][]int, 0, limit)
	for match != nil ***REMOVED***
		groups := match.Groups()

		result := make([]int, 0, len(groups)<<1)

		for _, group := range groups ***REMOVED***
			if len(group.Captures) > 0 ***REMOVED***
				startPos := group.Index
				endPos := group.Index + group.Length
				result = append(result, startPos, endPos)
			***REMOVED*** else ***REMOVED***
				result = append(result, -1, 0)
			***REMOVED***
		***REMOVED***

		if sticky && len(result) > 1 ***REMOVED***
			if result[0] != start ***REMOVED***
				break
			***REMOVED***
			start = result[1]
		***REMOVED***

		results = append(results, result)
		limit--
		if limit <= 0 ***REMOVED***
			break
		***REMOVED***
		match, err = wrapped.FindNextMatch(match)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return results
***REMOVED***

func buildPosMap(rd io.RuneReader, l, start int) (posMap []int, runes []rune, mappedStart int) ***REMOVED***
	posMap = make([]int, 0, l+1)
	curPos := 0
	runes = make([]rune, 0, l)
	startFound := false
	for ***REMOVED***
		if !startFound ***REMOVED***
			if curPos == start ***REMOVED***
				mappedStart = len(runes)
				startFound = true
			***REMOVED***
			if curPos > start ***REMOVED***
				// start position splits a surrogate pair
				mappedStart = len(runes) - 1
				_, second := utf16.EncodeRune(runes[mappedStart])
				runes[mappedStart] = second
				startFound = true
			***REMOVED***
		***REMOVED***
		rn, size, err := rd.ReadRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		runes = append(runes, rn)
		posMap = append(posMap, curPos)
		curPos += size
	***REMOVED***
	posMap = append(posMap, curPos)
	return
***REMOVED***

func (r *regexp2Wrapper) findAllSubmatchIndexUnicode(s unicodeString, start, limit int, sticky bool) [][]int ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	if limit < 0 ***REMOVED***
		limit = len(s) + 1
	***REMOVED***
	results := make([][]int, 0, limit)
	posMap, runes, mappedStart := buildPosMap(&lenientUtf16Decoder***REMOVED***utf16Reader: s.utf16Reader(0)***REMOVED***, s.length(), start)

	match, err := wrapped.FindRunesMatchStartingAt(runes, mappedStart)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	for match != nil ***REMOVED***
		groups := match.Groups()

		result := make([]int, 0, len(groups)<<1)

		for _, group := range groups ***REMOVED***
			if len(group.Captures) > 0 ***REMOVED***
				start := posMap[group.Index]
				end := posMap[group.Index+group.Length]
				result = append(result, start, end)
			***REMOVED*** else ***REMOVED***
				result = append(result, -1, 0)
			***REMOVED***
		***REMOVED***

		if sticky && len(result) > 1 ***REMOVED***
			if result[0] != start ***REMOVED***
				break
			***REMOVED***
			start = result[1]
		***REMOVED***

		results = append(results, result)
		match, err = wrapped.FindNextMatch(match)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return results
***REMOVED***

func (r *regexp2Wrapper) findAllSubmatchIndex(s valueString, start, limit int, sticky, fullUnicode bool) [][]int ***REMOVED***
	switch s := s.(type) ***REMOVED***
	case asciiString:
		return r.findAllSubmatchIndexUTF16(s, start, limit, sticky)
	case unicodeString:
		if fullUnicode ***REMOVED***
			return r.findAllSubmatchIndexUnicode(s, start, limit, sticky)
		***REMOVED***
		return r.findAllSubmatchIndexUTF16(s, start, limit, sticky)
	default:
		panic("Unsupported string type")
	***REMOVED***
***REMOVED***

func (r *regexpWrapper) findAllSubmatchIndex(s string, limit int, sticky bool) (results [][]int) ***REMOVED***
	wrapped := (*regexp.Regexp)(r)
	results = wrapped.FindAllStringSubmatchIndex(s, limit)
	pos := 0
	if sticky ***REMOVED***
		for i, result := range results ***REMOVED***
			if len(result) > 1 ***REMOVED***
				if result[0] != pos ***REMOVED***
					return results[:i]
				***REMOVED***
				pos = result[1]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (r *regexpWrapper) findSubmatchIndex(s valueString, fullUnicode bool) (result []int) ***REMOVED***
	wrapped := (*regexp.Regexp)(r)
	if fullUnicode ***REMOVED***
		posMap, runes, _ := buildPosMap(&lenientUtf16Decoder***REMOVED***utf16Reader: s.utf16Reader(0)***REMOVED***, s.length(), 0)
		res := wrapped.FindReaderSubmatchIndex(&arrayRuneReader***REMOVED***runes: runes***REMOVED***)
		for i, item := range res ***REMOVED***
			res[i] = posMap[item]
		***REMOVED***
		return res
	***REMOVED***
	return wrapped.FindReaderSubmatchIndex(s.utf16Reader(0))
***REMOVED***

func (r *regexpObject) execResultToArray(target valueString, result []int) Value ***REMOVED***
	captureCount := len(result) >> 1
	valueArray := make([]Value, captureCount)
	matchIndex := result[0]
	lowerBound := matchIndex
	for index := 0; index < captureCount; index++ ***REMOVED***
		offset := index << 1
		if result[offset] >= lowerBound ***REMOVED***
			valueArray[index] = target.substring(result[offset], result[offset+1])
			lowerBound = result[offset]
		***REMOVED*** else ***REMOVED***
			valueArray[index] = _undefined
		***REMOVED***
	***REMOVED***
	match := r.val.runtime.newArrayValues(valueArray)
	match.self.setOwnStr("input", target, false)
	match.self.setOwnStr("index", intToValue(int64(matchIndex)), false)
	return match
***REMOVED***

func (r *regexpObject) getLastIndex() int64 ***REMOVED***
	lastIndex := toLength(r.getStr("lastIndex", nil))
	if !r.pattern.global && !r.pattern.sticky ***REMOVED***
		return 0
	***REMOVED***
	return lastIndex
***REMOVED***

func (r *regexpObject) updateLastIndex(index int64, firstResult, lastResult []int) bool ***REMOVED***
	if r.pattern.sticky ***REMOVED***
		if firstResult == nil || int64(firstResult[0]) != index ***REMOVED***
			r.setOwnStr("lastIndex", intToValue(0), true)
			return false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if firstResult == nil ***REMOVED***
			if r.pattern.global ***REMOVED***
				r.setOwnStr("lastIndex", intToValue(0), true)
			***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	if r.pattern.global || r.pattern.sticky ***REMOVED***
		r.setOwnStr("lastIndex", intToValue(int64(lastResult[1])), true)
	***REMOVED***
	return true
***REMOVED***

func (r *regexpObject) execRegexp(target valueString) (match bool, result []int) ***REMOVED***
	index := r.getLastIndex()
	if index >= 0 && index <= int64(target.length()) ***REMOVED***
		result = r.pattern.findSubmatchIndex(target, int(index))
	***REMOVED***
	match = r.updateLastIndex(index, result, result)
	return
***REMOVED***

func (r *regexpObject) exec(target valueString) Value ***REMOVED***
	match, result := r.execRegexp(target)
	if match ***REMOVED***
		return r.execResultToArray(target, result)
	***REMOVED***
	return _null
***REMOVED***

func (r *regexpObject) test(target valueString) bool ***REMOVED***
	match, _ := r.execRegexp(target)
	return match
***REMOVED***

func (r *regexpObject) clone() *Object ***REMOVED***
	r1 := r.val.runtime.newRegexpObject(r.prototype)
	r1.source = r.source
	r1.pattern = r.pattern

	return r1.val
***REMOVED***

func (r *regexpObject) init() ***REMOVED***
	r.baseObject.init()
	r.standard = true
	r._putProp("lastIndex", intToValue(0), true, false, false)
***REMOVED***

func (r *regexpObject) setProto(proto *Object, throw bool) bool ***REMOVED***
	res := r.baseObject.setProto(proto, throw)
	if res ***REMOVED***
		r.standard = false
	***REMOVED***
	return res
***REMOVED***

func (r *regexpObject) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	res := r.baseObject.defineOwnPropertyStr(name, desc, throw)
	if res ***REMOVED***
		r.standard = false
	***REMOVED***
	return res
***REMOVED***

func (r *regexpObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	res := r.baseObject.deleteStr(name, throw)
	if res ***REMOVED***
		r.standard = false
	***REMOVED***
	return res
***REMOVED***

func (r *regexpObject) setOwnStr(name unistring.String, value Value, throw bool) bool ***REMOVED***
	if r.standard ***REMOVED***
		if name == "exec" ***REMOVED***
			res := r.baseObject.setOwnStr(name, value, throw)
			if res ***REMOVED***
				r.standard = false
			***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return r.baseObject.setOwnStr(name, value, throw)
***REMOVED***
