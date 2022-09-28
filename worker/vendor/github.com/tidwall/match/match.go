// Package match provides a simple pattern matcher with unicode support.
package match

import (
	"unicode/utf8"
)

// Match returns true if str matches pattern. This is a very
// simple wildcard match where '*' matches on any number characters
// and '?' matches on any one character.
//
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
	return match(str, pattern, 0, nil, -1) == rMatch
***REMOVED***

// MatchLimit is the same as Match but will limit the complexity of the match
// operation. This is to avoid long running matches, specifically to avoid ReDos
// attacks from arbritary inputs.
//
// How it works:
// The underlying match routine is recursive and may call itself when it
// encounters a sandwiched wildcard pattern, such as: `user:*:name`.
// Everytime it calls itself a counter is incremented.
// The operation is stopped when counter > maxcomp*len(str).
func MatchLimit(str, pattern string, maxcomp int) (matched, stopped bool) ***REMOVED***
	if pattern == "*" ***REMOVED***
		return true, false
	***REMOVED***
	counter := 0
	r := match(str, pattern, len(str), &counter, maxcomp)
	if r == rStop ***REMOVED***
		return false, true
	***REMOVED***
	return r == rMatch, false
***REMOVED***

type result int

const (
	rNoMatch result = iota
	rMatch
	rStop
)

func match(str, pat string, slen int, counter *int, maxcomp int) result ***REMOVED***
	// check complexity limit
	if maxcomp > -1 ***REMOVED***
		if *counter > slen*maxcomp ***REMOVED***
			return rStop
		***REMOVED***
		*counter++
	***REMOVED***

	for len(pat) > 0 ***REMOVED***
		var wild bool
		pc, ps := rune(pat[0]), 1
		if pc > 0x7f ***REMOVED***
			pc, ps = utf8.DecodeRuneInString(pat)
		***REMOVED***
		var sc rune
		var ss int
		if len(str) > 0 ***REMOVED***
			sc, ss = rune(str[0]), 1
			if sc > 0x7f ***REMOVED***
				sc, ss = utf8.DecodeRuneInString(str)
			***REMOVED***
		***REMOVED***
		switch pc ***REMOVED***
		case '?':
			if ss == 0 ***REMOVED***
				return rNoMatch
			***REMOVED***
		case '*':
			// Ignore repeating stars.
			for len(pat) > 1 && pat[1] == '*' ***REMOVED***
				pat = pat[1:]
			***REMOVED***

			// If this star is the last character then it must be a match.
			if len(pat) == 1 ***REMOVED***
				return rMatch
			***REMOVED***

			// Match and trim any non-wildcard suffix characters.
			var ok bool
			str, pat, ok = matchTrimSuffix(str, pat)
			if !ok ***REMOVED***
				return rNoMatch
			***REMOVED***

			// Check for single star again.
			if len(pat) == 1 ***REMOVED***
				return rMatch
			***REMOVED***

			// Perform recursive wildcard search.
			r := match(str, pat[1:], slen, counter, maxcomp)
			if r != rNoMatch ***REMOVED***
				return r
			***REMOVED***
			if len(str) == 0 ***REMOVED***
				return rNoMatch
			***REMOVED***
			wild = true
		default:
			if ss == 0 ***REMOVED***
				return rNoMatch
			***REMOVED***
			if pc == '\\' ***REMOVED***
				pat = pat[ps:]
				pc, ps = utf8.DecodeRuneInString(pat)
				if ps == 0 ***REMOVED***
					return rNoMatch
				***REMOVED***
			***REMOVED***
			if sc != pc ***REMOVED***
				return rNoMatch
			***REMOVED***
		***REMOVED***
		str = str[ss:]
		if !wild ***REMOVED***
			pat = pat[ps:]
		***REMOVED***
	***REMOVED***
	if len(str) == 0 ***REMOVED***
		return rMatch
	***REMOVED***
	return rNoMatch
***REMOVED***

// matchTrimSuffix matches and trims any non-wildcard suffix characters.
// Returns the trimed string and pattern.
//
// This is called because the pattern contains extra data after the wildcard
// star. Here we compare any suffix characters in the pattern to the suffix of
// the target string. Basically a reverse match that stops when a wildcard
// character is reached. This is a little trickier than a forward match because
// we need to evaluate an escaped character in reverse.
//
// Any matched characters will be trimmed from both the target
// string and the pattern.
func matchTrimSuffix(str, pat string) (string, string, bool) ***REMOVED***
	// It's expected that the pattern has at least two bytes and the first byte
	// is a wildcard star '*'
	match := true
	for len(str) > 0 && len(pat) > 1 ***REMOVED***
		pc, ps := utf8.DecodeLastRuneInString(pat)
		var esc bool
		for i := 0; ; i++ ***REMOVED***
			if pat[len(pat)-ps-i-1] != '\\' ***REMOVED***
				if i&1 == 1 ***REMOVED***
					esc = true
					ps++
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if pc == '*' && !esc ***REMOVED***
			match = true
			break
		***REMOVED***
		sc, ss := utf8.DecodeLastRuneInString(str)
		if !((pc == '?' && !esc) || pc == sc) ***REMOVED***
			match = false
			break
		***REMOVED***
		str = str[:len(str)-ss]
		pat = pat[:len(pat)-ps]
	***REMOVED***
	return str, pat, match
***REMOVED***

var maxRuneBytes = [...]byte***REMOVED***244, 143, 191, 191***REMOVED***

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
			maxb = append(maxb, maxRuneBytes[:]...)
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
***REMOVED***

// IsPattern returns true if the string is a pattern.
func IsPattern(str string) bool ***REMOVED***
	for i := 0; i < len(str); i++ ***REMOVED***
		if str[i] == '*' || str[i] == '?' ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
