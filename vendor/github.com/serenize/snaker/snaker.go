// Package snaker provides methods to convert CamelCase names to snake_case and back.
// It considers the list of allowed initialsms used by github.com/golang/lint/golint (e.g. ID or HTTP)
package snaker

import (
	"strings"
	"unicode"
)

// CamelToSnake converts a given string to snake case
func CamelToSnake(s string) string ***REMOVED***
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ ***REMOVED***
		if i > 0 && unicode.IsUpper(rs[i]) ***REMOVED***
			if initialism := startsWithInitialism(s[lastPos:]); initialism != "" ***REMOVED***
				words = append(words, initialism)

				i += len(initialism) - 1
				lastPos = i
				continue
			***REMOVED***

			words = append(words, s[lastPos:i])
			lastPos = i
		***REMOVED***
	***REMOVED***

	// append the last word
	if s[lastPos:] != "" ***REMOVED***
		words = append(words, s[lastPos:])
	***REMOVED***

	for k, word := range words ***REMOVED***
		if k > 0 ***REMOVED***
			result += "_"
		***REMOVED***

		result += strings.ToLower(word)
	***REMOVED***

	return result
***REMOVED***

func snakeToCamel(s string, upperCase bool) string ***REMOVED***
	var result string

	words := strings.Split(s, "_")

	for i, word := range words ***REMOVED***
		if exception := snakeToCamelExceptions[word]; len(exception) > 0 ***REMOVED***
			result += exception
			continue
		***REMOVED***

		if upperCase || i > 0 ***REMOVED***
			if upper := strings.ToUpper(word); commonInitialisms[upper] ***REMOVED***
				result += upper
				continue
			***REMOVED***
		***REMOVED***

		if (upperCase || i > 0) && len(word) > 0 ***REMOVED***
			w := []rune(word)
			w[0] = unicode.ToUpper(w[0])
			result += string(w)
		***REMOVED*** else ***REMOVED***
			result += word
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***

// SnakeToCamel returns a string converted from snake case to uppercase
func SnakeToCamel(s string) string ***REMOVED***
	return snakeToCamel(s, true)
***REMOVED***

// SnakeToCamelLower returns a string converted from snake case to lowercase
func SnakeToCamelLower(s string) string ***REMOVED***
	return snakeToCamel(s, false)
***REMOVED***

// startsWithInitialism returns the initialism if the given string begins with it
func startsWithInitialism(s string) string ***REMOVED***
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	for i := 1; i <= 5; i++ ***REMOVED***
		if len(s) > i-1 && commonInitialisms[s[:i]] ***REMOVED***
			initialism = s[:i]
		***REMOVED***
	***REMOVED***
	return initialism
***REMOVED***

// commonInitialisms, taken from
// https://github.com/golang/lint/blob/206c0f020eba0f7fbcfbc467a5eb808037df2ed6/lint.go#L731
var commonInitialisms = map[string]bool***REMOVED***
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"ETA":   true,
	"GPU":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"OS":    true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
	"OAuth": true,
***REMOVED***

// add exceptions here for things that are not automatically convertable
var snakeToCamelExceptions = map[string]string***REMOVED***
	"oauth": "OAuth",
***REMOVED***
