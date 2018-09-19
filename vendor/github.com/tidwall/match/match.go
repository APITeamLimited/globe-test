// Match provides a simple pattern matcher with unicode support.
package match

import "unicode/utf8"

// Match returns true if str matches pattern. This is a very
// simple wildcard match where '*' matches on any number characters
// and '?' matches on any one character.

// pattern:
// 	***REMOVED*** term ***REMOVED***
// term:
// 	'*'         matches any sequence of non-Separator characters
// 	'?'         matches any single non-Separator character
// 	c           matches character c (c != '*', '?', '\\')
// 	'\\' c      matches character c
//
func Match(str, pattern string) bool ***REMOVED***
	if pattern == "*" ***REMOVED***
		return true
	***REMOVED***
	return deepMatch(str, pattern)
***REMOVED***
func deepMatch(str, pattern string) bool ***REMOVED***
	for len(pattern) > 0 ***REMOVED***
		if pattern[0] > 0x7f ***REMOVED***
			return deepMatchRune(str, pattern)
		***REMOVED***
		switch pattern[0] ***REMOVED***
		default:
			if len(str) == 0 ***REMOVED***
				return false
			***REMOVED***
			if str[0] > 0x7f ***REMOVED***
				return deepMatchRune(str, pattern)
			***REMOVED***
			if str[0] != pattern[0] ***REMOVED***
				return false
			***REMOVED***
		case '?':
			if len(str) == 0 ***REMOVED***
				return false
			***REMOVED***
		case '*':
			return deepMatch(str, pattern[1:]) ||
				(len(str) > 0 && deepMatch(str[1:], pattern))
		***REMOVED***
		str = str[1:]
		pattern = pattern[1:]
	***REMOVED***
	return len(str) == 0 && len(pattern) == 0
***REMOVED***

func deepMatchRune(str, pattern string) bool ***REMOVED***
	var sr, pr rune
	var srsz, prsz int

	// read the first rune ahead of time
	if len(str) > 0 ***REMOVED***
		if str[0] > 0x7f ***REMOVED***
			sr, srsz = utf8.DecodeRuneInString(str)
		***REMOVED*** else ***REMOVED***
			sr, srsz = rune(str[0]), 1
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		sr, srsz = utf8.RuneError, 0
	***REMOVED***
	if len(pattern) > 0 ***REMOVED***
		if pattern[0] > 0x7f ***REMOVED***
			pr, prsz = utf8.DecodeRuneInString(pattern)
		***REMOVED*** else ***REMOVED***
			pr, prsz = rune(pattern[0]), 1
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		pr, prsz = utf8.RuneError, 0
	***REMOVED***
	// done reading
	for pr != utf8.RuneError ***REMOVED***
		switch pr ***REMOVED***
		default:
			if srsz == utf8.RuneError ***REMOVED***
				return false
			***REMOVED***
			if sr != pr ***REMOVED***
				return false
			***REMOVED***
		case '?':
			if srsz == utf8.RuneError ***REMOVED***
				return false
			***REMOVED***
		case '*':
			return deepMatchRune(str, pattern[prsz:]) ||
				(srsz > 0 && deepMatchRune(str[srsz:], pattern))
		***REMOVED***
		str = str[srsz:]
		pattern = pattern[prsz:]
		// read the next runes
		if len(str) > 0 ***REMOVED***
			if str[0] > 0x7f ***REMOVED***
				sr, srsz = utf8.DecodeRuneInString(str)
			***REMOVED*** else ***REMOVED***
				sr, srsz = rune(str[0]), 1
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			sr, srsz = utf8.RuneError, 0
		***REMOVED***
		if len(pattern) > 0 ***REMOVED***
			if pattern[0] > 0x7f ***REMOVED***
				pr, prsz = utf8.DecodeRuneInString(pattern)
			***REMOVED*** else ***REMOVED***
				pr, prsz = rune(pattern[0]), 1
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			pr, prsz = utf8.RuneError, 0
		***REMOVED***
		// done reading
	***REMOVED***

	return srsz == 0 && prsz == 0
***REMOVED***

var maxRuneBytes = func() []byte ***REMOVED***
	b := make([]byte, 4)
	if utf8.EncodeRune(b, '\U0010FFFF') != 4 ***REMOVED***
		panic("invalid rune encoding")
	***REMOVED***
	return b
***REMOVED***()

// Allowable parses the pattern and determines the minimum and maximum allowable
// values that the pattern can represent.
// When the max cannot be determined, 'true' will be returned
// for infinite.
func Allowable(pattern string) (min, max string) ***REMOVED***
	if pattern == "" || pattern[0] == '*' ***REMOVED***
		return "", ""
	***REMOVED***

	minb := make([]byte, 0, len(pattern))
	maxb := make([]byte, 0, len(pattern))
	var wild bool
	for i := 0; i < len(pattern); i++ ***REMOVED***
		if pattern[i] == '*' ***REMOVED***
			wild = true
			break
		***REMOVED***
		if pattern[i] == '?' ***REMOVED***
			minb = append(minb, 0)
			maxb = append(maxb, maxRuneBytes...)
		***REMOVED*** else ***REMOVED***
			minb = append(minb, pattern[i])
			maxb = append(maxb, pattern[i])
		***REMOVED***
	***REMOVED***
	if wild ***REMOVED***
		r, n := utf8.DecodeLastRune(maxb)
		if r != utf8.RuneError ***REMOVED***
			if r < utf8.MaxRune ***REMOVED***
				r++
				if r > 0x7f ***REMOVED***
					b := make([]byte, 4)
					nn := utf8.EncodeRune(b, r)
					maxb = append(maxb[:len(maxb)-n], b[:nn]...)
				***REMOVED*** else ***REMOVED***
					maxb = append(maxb[:len(maxb)-n], byte(r))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return string(minb), string(maxb)
	/*
		return
		if wild ***REMOVED***
			r, n := utf8.DecodeLastRune(maxb)
			if r != utf8.RuneError ***REMOVED***
				if r < utf8.MaxRune ***REMOVED***
					infinite = true
				***REMOVED*** else ***REMOVED***
					r++
					if r > 0x7f ***REMOVED***
						b := make([]byte, 4)
						nn := utf8.EncodeRune(b, r)
						maxb = append(maxb[:len(maxb)-n], b[:nn]...)
					***REMOVED*** else ***REMOVED***
						maxb = append(maxb[:len(maxb)-n], byte(r))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return string(minb), string(maxb), infinite
	*/
***REMOVED***
