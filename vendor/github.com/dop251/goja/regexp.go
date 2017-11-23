package goja

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"regexp"
	"unicode/utf16"
	"unicode/utf8"
)

type regexpPattern interface ***REMOVED***
	FindSubmatchIndex(valueString, int) []int
	FindAllSubmatchIndex(valueString, int) [][]int
	FindAllSubmatchIndexUTF8(string, int) [][]int
	FindAllSubmatchIndexASCII(string, int) [][]int
	MatchString(valueString) bool
***REMOVED***

type regexp2Wrapper regexp2.Regexp
type regexpWrapper regexp.Regexp

type regexpObject struct ***REMOVED***
	baseObject
	pattern regexpPattern
	source  valueString

	global, multiline, ignoreCase bool
***REMOVED***

func (r *regexp2Wrapper) FindSubmatchIndex(s valueString, start int) (result []int) ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	var match *regexp2.Match
	var err error
	switch s := s.(type) ***REMOVED***
	case asciiString:
		match, err = wrapped.FindStringMatch(string(s)[start:])
	case unicodeString:
		match, err = wrapped.FindRunesMatch(utf16.Decode(s[start:]))
	default:
		panic(fmt.Errorf("Unknown string type: %T", s))
	***REMOVED***
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

func (r *regexp2Wrapper) FindAllSubmatchIndexUTF8(s string, n int) [][]int ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	if n < 0 ***REMOVED***
		n = len(s) + 1
	***REMOVED***
	results := make([][]int, 0, n)

	idxMap := make([]int, 0, len(s))
	runes := make([]rune, 0, len(s))
	for pos, rr := range s ***REMOVED***
		runes = append(runes, rr)
		idxMap = append(idxMap, pos)
	***REMOVED***
	idxMap = append(idxMap, len(s))

	match, err := wrapped.FindRunesMatch(runes)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	i := 0
	for match != nil && i < n ***REMOVED***
		groups := match.Groups()

		result := make([]int, 0, len(groups)<<1)

		for _, group := range groups ***REMOVED***
			if len(group.Captures) > 0 ***REMOVED***
				result = append(result, idxMap[group.Index], idxMap[group.Index+group.Length])
			***REMOVED*** else ***REMOVED***
				result = append(result, -1, 0)
			***REMOVED***
		***REMOVED***

		results = append(results, result)
		match, err = wrapped.FindNextMatch(match)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		i++
	***REMOVED***
	return results
***REMOVED***

func (r *regexp2Wrapper) FindAllSubmatchIndexASCII(s string, n int) [][]int ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	if n < 0 ***REMOVED***
		n = len(s) + 1
	***REMOVED***
	results := make([][]int, 0, n)

	match, err := wrapped.FindStringMatch(s)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	i := 0
	for match != nil && i < n ***REMOVED***
		groups := match.Groups()

		result := make([]int, 0, len(groups)<<1)

		for _, group := range groups ***REMOVED***
			if len(group.Captures) > 0 ***REMOVED***
				result = append(result, group.Index, group.Index+group.Length)
			***REMOVED*** else ***REMOVED***
				result = append(result, -1, 0)
			***REMOVED***
		***REMOVED***

		results = append(results, result)
		match, err = wrapped.FindNextMatch(match)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
		i++
	***REMOVED***
	return results
***REMOVED***

func (r *regexp2Wrapper) findAllSubmatchIndexUTF16(s unicodeString, n int) [][]int ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)
	if n < 0 ***REMOVED***
		n = len(s) + 1
	***REMOVED***
	results := make([][]int, 0, n)

	rd := runeReaderReplace***REMOVED***s.reader(0)***REMOVED***
	posMap := make([]int, s.length()+1)
	curPos := 0
	curRuneIdx := 0
	runes := make([]rune, 0, s.length())
	for ***REMOVED***
		rn, size, err := rd.ReadRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		runes = append(runes, rn)
		posMap[curRuneIdx] = curPos
		curRuneIdx++
		curPos += size
	***REMOVED***
	posMap[curRuneIdx] = curPos

	match, err := wrapped.FindRunesMatch(runes)
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

		results = append(results, result)
		match, err = wrapped.FindNextMatch(match)
		if err != nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return results
***REMOVED***

func (r *regexp2Wrapper) FindAllSubmatchIndex(s valueString, n int) [][]int ***REMOVED***
	switch s := s.(type) ***REMOVED***
	case asciiString:
		return r.FindAllSubmatchIndexASCII(string(s), n)
	case unicodeString:
		return r.findAllSubmatchIndexUTF16(s, n)
	default:
		panic("Unsupported string type")
	***REMOVED***
***REMOVED***

func (r *regexp2Wrapper) MatchString(s valueString) bool ***REMOVED***
	wrapped := (*regexp2.Regexp)(r)

	switch s := s.(type) ***REMOVED***
	case asciiString:
		matched, _ := wrapped.MatchString(string(s))
		return matched
	case unicodeString:
		matched, _ := wrapped.MatchRunes(utf16.Decode(s))
		return matched
	default:
		panic(fmt.Errorf("Unknown string type: %T", s))
	***REMOVED***
***REMOVED***

func (r *regexpWrapper) FindSubmatchIndex(s valueString, start int) (result []int) ***REMOVED***
	wrapped := (*regexp.Regexp)(r)
	return wrapped.FindReaderSubmatchIndex(runeReaderReplace***REMOVED***s.reader(start)***REMOVED***)
***REMOVED***

func (r *regexpWrapper) MatchString(s valueString) bool ***REMOVED***
	wrapped := (*regexp.Regexp)(r)
	return wrapped.MatchReader(runeReaderReplace***REMOVED***s.reader(0)***REMOVED***)
***REMOVED***

func (r *regexpWrapper) FindAllSubmatchIndex(s valueString, n int) [][]int ***REMOVED***
	wrapped := (*regexp.Regexp)(r)
	switch s := s.(type) ***REMOVED***
	case asciiString:
		return wrapped.FindAllStringSubmatchIndex(string(s), n)
	case unicodeString:
		return r.findAllSubmatchIndexUTF16(s, n)
	default:
		panic("Unsupported string type")
	***REMOVED***
***REMOVED***

func (r *regexpWrapper) FindAllSubmatchIndexUTF8(s string, n int) [][]int ***REMOVED***
	wrapped := (*regexp.Regexp)(r)
	return wrapped.FindAllStringSubmatchIndex(s, n)
***REMOVED***

func (r *regexpWrapper) FindAllSubmatchIndexASCII(s string, n int) [][]int ***REMOVED***
	return r.FindAllSubmatchIndexUTF8(s, n)
***REMOVED***

func (r *regexpWrapper) findAllSubmatchIndexUTF16(s unicodeString, n int) [][]int ***REMOVED***
	wrapped := (*regexp.Regexp)(r)
	utf8Bytes := make([]byte, 0, len(s)*2)
	posMap := make(map[int]int)
	curPos := 0
	rd := runeReaderReplace***REMOVED***s.reader(0)***REMOVED***
	for ***REMOVED***
		rn, size, err := rd.ReadRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		l := len(utf8Bytes)
		utf8Bytes = append(utf8Bytes, 0, 0, 0, 0)
		n := utf8.EncodeRune(utf8Bytes[l:], rn)
		utf8Bytes = utf8Bytes[:l+n]
		posMap[l] = curPos
		curPos += size
	***REMOVED***
	posMap[len(utf8Bytes)] = curPos

	rr := wrapped.FindAllSubmatchIndex(utf8Bytes, n)
	for _, res := range rr ***REMOVED***
		for j, pos := range res ***REMOVED***
			mapped, exists := posMap[pos]
			if !exists ***REMOVED***
				panic("Unicode match is not on rune boundary")
			***REMOVED***
			res[j] = mapped
		***REMOVED***
	***REMOVED***
	return rr
***REMOVED***

func (r *regexpObject) execResultToArray(target valueString, result []int) Value ***REMOVED***
	captureCount := len(result) >> 1
	valueArray := make([]Value, captureCount)
	matchIndex := result[0]
	lowerBound := matchIndex
	for index := 0; index < captureCount; index++ ***REMOVED***
		offset := index << 1
		if result[offset] >= lowerBound ***REMOVED***
			valueArray[index] = target.substring(int64(result[offset]), int64(result[offset+1]))
			lowerBound = result[offset]
		***REMOVED*** else ***REMOVED***
			valueArray[index] = _undefined
		***REMOVED***
	***REMOVED***
	match := r.val.runtime.newArrayValues(valueArray)
	match.self.putStr("input", target, false)
	match.self.putStr("index", intToValue(int64(matchIndex)), false)
	return match
***REMOVED***

func (r *regexpObject) execRegexp(target valueString) (match bool, result []int) ***REMOVED***
	lastIndex := int64(0)
	if p := r.getStr("lastIndex"); p != nil ***REMOVED***
		lastIndex = p.ToInteger()
		if lastIndex < 0 ***REMOVED***
			lastIndex = 0
		***REMOVED***
	***REMOVED***
	index := lastIndex
	if !r.global ***REMOVED***
		index = 0
	***REMOVED***
	if index >= 0 && index <= target.length() ***REMOVED***
		result = r.pattern.FindSubmatchIndex(target, int(index))
	***REMOVED***
	if result == nil ***REMOVED***
		r.putStr("lastIndex", intToValue(0), true)
		return
	***REMOVED***
	match = true
	startIndex := index
	endIndex := int(lastIndex) + result[1]
	// We do this shift here because the .FindStringSubmatchIndex above
	// was done on a local subordinate slice of the string, not the whole string
	for index, _ := range result ***REMOVED***
		result[index] += int(startIndex)
	***REMOVED***
	if r.global ***REMOVED***
		r.putStr("lastIndex", intToValue(int64(endIndex)), true)
	***REMOVED***
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
	r1.global = r.global
	r1.ignoreCase = r.ignoreCase
	r1.multiline = r.multiline
	return r1.val
***REMOVED***

func (r *regexpObject) init() ***REMOVED***
	r.baseObject.init()
	r._putProp("lastIndex", intToValue(0), true, false, false)
***REMOVED***
