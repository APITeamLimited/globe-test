// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package strs provides string manipulation functionality specific to protobuf.
package strs

import (
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// EnforceUTF8 reports whether to enforce strict UTF-8 validation.
func EnforceUTF8(fd protoreflect.FieldDescriptor) bool ***REMOVED***
	if flags.ProtoLegacy ***REMOVED***
		if fd, ok := fd.(interface***REMOVED*** EnforceUTF8() bool ***REMOVED***); ok ***REMOVED***
			return fd.EnforceUTF8()
		***REMOVED***
	***REMOVED***
	return fd.Syntax() == protoreflect.Proto3
***REMOVED***

// GoCamelCase camel-cases a protobuf name for use as a Go identifier.
//
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
func GoCamelCase(s string) string ***REMOVED***
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	var b []byte
	for i := 0; i < len(s); i++ ***REMOVED***
		c := s[i]
		switch ***REMOVED***
		case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '.' in ".***REMOVED******REMOVED***lowercase***REMOVED******REMOVED***".
		case c == '.':
			b = append(b, '_') // convert '.' to '_'
		case c == '_' && (i == 0 || s[i-1] == '.'):
			// Convert initial '_' to ensure we start with a capital letter.
			// Do the same for '_' after '.' to match historic behavior.
			b = append(b, 'X') // convert '_' to 'X'
		case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '_' in "_***REMOVED******REMOVED***lowercase***REMOVED******REMOVED***".
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			// Assume we have a letter now - if not, it's a bogus identifier.
			// The next word is a sequence of characters that must start upper case.
			if isASCIILower(c) ***REMOVED***
				c -= 'a' - 'A' // convert lowercase to uppercase
			***REMOVED***
			b = append(b, c)

			// Accept lower case sequence that follows.
			for ; i+1 < len(s) && isASCIILower(s[i+1]); i++ ***REMOVED***
				b = append(b, s[i+1])
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return string(b)
***REMOVED***

// GoSanitized converts a string to a valid Go identifier.
func GoSanitized(s string) string ***REMOVED***
	// Sanitize the input to the set of valid characters,
	// which must be '_' or be in the Unicode L or N categories.
	s = strings.Map(func(r rune) rune ***REMOVED***
		if unicode.IsLetter(r) || unicode.IsDigit(r) ***REMOVED***
			return r
		***REMOVED***
		return '_'
	***REMOVED***, s)

	// Prepend '_' in the event of a Go keyword conflict or if
	// the identifier is invalid (does not start in the Unicode L category).
	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) ***REMOVED***
		return "_" + s
	***REMOVED***
	return s
***REMOVED***

// JSONCamelCase converts a snake_case identifier to a camelCase identifier,
// according to the protobuf JSON specification.
func JSONCamelCase(s string) string ***REMOVED***
	var b []byte
	var wasUnderscore bool
	for i := 0; i < len(s); i++ ***REMOVED*** // proto identifiers are always ASCII
		c := s[i]
		if c != '_' ***REMOVED***
			if wasUnderscore && isASCIILower(c) ***REMOVED***
				c -= 'a' - 'A' // convert to uppercase
			***REMOVED***
			b = append(b, c)
		***REMOVED***
		wasUnderscore = c == '_'
	***REMOVED***
	return string(b)
***REMOVED***

// JSONSnakeCase converts a camelCase identifier to a snake_case identifier,
// according to the protobuf JSON specification.
func JSONSnakeCase(s string) string ***REMOVED***
	var b []byte
	for i := 0; i < len(s); i++ ***REMOVED*** // proto identifiers are always ASCII
		c := s[i]
		if isASCIIUpper(c) ***REMOVED***
			b = append(b, '_')
			c += 'a' - 'A' // convert to lowercase
		***REMOVED***
		b = append(b, c)
	***REMOVED***
	return string(b)
***REMOVED***

// MapEntryName derives the name of the map entry message given the field name.
// See protoc v3.8.0: src/google/protobuf/descriptor.cc:254-276,6057
func MapEntryName(s string) string ***REMOVED***
	var b []byte
	upperNext := true
	for _, c := range s ***REMOVED***
		switch ***REMOVED***
		case c == '_':
			upperNext = true
		case upperNext:
			b = append(b, byte(unicode.ToUpper(c)))
			upperNext = false
		default:
			b = append(b, byte(c))
		***REMOVED***
	***REMOVED***
	b = append(b, "Entry"...)
	return string(b)
***REMOVED***

// EnumValueName derives the camel-cased enum value name.
// See protoc v3.8.0: src/google/protobuf/descriptor.cc:297-313
func EnumValueName(s string) string ***REMOVED***
	var b []byte
	upperNext := true
	for _, c := range s ***REMOVED***
		switch ***REMOVED***
		case c == '_':
			upperNext = true
		case upperNext:
			b = append(b, byte(unicode.ToUpper(c)))
			upperNext = false
		default:
			b = append(b, byte(unicode.ToLower(c)))
			upperNext = false
		***REMOVED***
	***REMOVED***
	return string(b)
***REMOVED***

// TrimEnumPrefix trims the enum name prefix from an enum value name,
// where the prefix is all lowercase without underscores.
// See protoc v3.8.0: src/google/protobuf/descriptor.cc:330-375
func TrimEnumPrefix(s, prefix string) string ***REMOVED***
	s0 := s // original input
	for len(s) > 0 && len(prefix) > 0 ***REMOVED***
		if s[0] == '_' ***REMOVED***
			s = s[1:]
			continue
		***REMOVED***
		if unicode.ToLower(rune(s[0])) != rune(prefix[0]) ***REMOVED***
			return s0 // no prefix match
		***REMOVED***
		s, prefix = s[1:], prefix[1:]
	***REMOVED***
	if len(prefix) > 0 ***REMOVED***
		return s0 // no prefix match
	***REMOVED***
	s = strings.TrimLeft(s, "_")
	if len(s) == 0 ***REMOVED***
		return s0 // avoid returning empty string
	***REMOVED***
	return s
***REMOVED***

func isASCIILower(c byte) bool ***REMOVED***
	return 'a' <= c && c <= 'z'
***REMOVED***
func isASCIIUpper(c byte) bool ***REMOVED***
	return 'A' <= c && c <= 'Z'
***REMOVED***
func isASCIIDigit(c byte) bool ***REMOVED***
	return '0' <= c && c <= '9'
***REMOVED***
