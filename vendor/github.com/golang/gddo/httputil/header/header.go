// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// Package header provides functions for parsing HTTP headers.
package header

import (
	"net/http"
	"strings"
	"time"
)

// Octet types from RFC 2616.
var octetTypes [256]octetType

type octetType byte

const (
	isToken octetType = 1 << iota
	isSpace
)

func init() ***REMOVED***
	// OCTET      = <any 8-bit sequence of data>
	// CHAR       = <any US-ASCII character (octets 0 - 127)>
	// CTL        = <any US-ASCII control character (octets 0 - 31) and DEL (127)>
	// CR         = <US-ASCII CR, carriage return (13)>
	// LF         = <US-ASCII LF, linefeed (10)>
	// SP         = <US-ASCII SP, space (32)>
	// HT         = <US-ASCII HT, horizontal-tab (9)>
	// <">        = <US-ASCII double-quote mark (34)>
	// CRLF       = CR LF
	// LWS        = [CRLF] 1*( SP | HT )
	// TEXT       = <any OCTET except CTLs, but including LWS>
	// separators = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | <">
	//              | "/" | "[" | "]" | "?" | "=" | "***REMOVED***" | "***REMOVED***" | SP | HT
	// token      = 1*<any CHAR except CTLs or separators>
	// qdtext     = <any TEXT except <">>

	for c := 0; c < 256; c++ ***REMOVED***
		var t octetType
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.IndexRune(" \t\"(),/:;<=>?@[]\\***REMOVED******REMOVED***", rune(c)) >= 0
		if strings.IndexRune(" \t\r\n", rune(c)) >= 0 ***REMOVED***
			t |= isSpace
		***REMOVED***
		if isChar && !isCtl && !isSeparator ***REMOVED***
			t |= isToken
		***REMOVED***
		octetTypes[c] = t
	***REMOVED***
***REMOVED***

// Copy returns a shallow copy of the header.
func Copy(header http.Header) http.Header ***REMOVED***
	h := make(http.Header)
	for k, vs := range header ***REMOVED***
		h[k] = vs
	***REMOVED***
	return h
***REMOVED***

var timeLayouts = []string***REMOVED***"Mon, 02 Jan 2006 15:04:05 GMT", time.RFC850, time.ANSIC***REMOVED***

// ParseTime parses the header as time. The zero value is returned if the
// header is not present or there is an error parsing the
// header.
func ParseTime(header http.Header, key string) time.Time ***REMOVED***
	if s := header.Get(key); s != "" ***REMOVED***
		for _, layout := range timeLayouts ***REMOVED***
			if t, err := time.Parse(layout, s); err == nil ***REMOVED***
				return t.UTC()
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return time.Time***REMOVED******REMOVED***
***REMOVED***

// ParseList parses a comma separated list of values. Commas are ignored in
// quoted strings. Quoted values are not unescaped or unquoted. Whitespace is
// trimmed.
func ParseList(header http.Header, key string) []string ***REMOVED***
	var result []string
	for _, s := range header[http.CanonicalHeaderKey(key)] ***REMOVED***
		begin := 0
		end := 0
		escape := false
		quote := false
		for i := 0; i < len(s); i++ ***REMOVED***
			b := s[i]
			switch ***REMOVED***
			case escape:
				escape = false
				end = i + 1
			case quote:
				switch b ***REMOVED***
				case '\\':
					escape = true
				case '"':
					quote = false
				***REMOVED***
				end = i + 1
			case b == '"':
				quote = true
				end = i + 1
			case octetTypes[b]&isSpace != 0:
				if begin == end ***REMOVED***
					begin = i + 1
					end = begin
				***REMOVED***
			case b == ',':
				if begin < end ***REMOVED***
					result = append(result, s[begin:end])
				***REMOVED***
				begin = i + 1
				end = begin
			default:
				end = i + 1
			***REMOVED***
		***REMOVED***
		if begin < end ***REMOVED***
			result = append(result, s[begin:end])
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// ParseValueAndParams parses a comma separated list of values with optional
// semicolon separated name-value pairs. Content-Type and Content-Disposition
// headers are in this format.
func ParseValueAndParams(header http.Header, key string) (value string, params map[string]string) ***REMOVED***
	params = make(map[string]string)
	s := header.Get(key)
	value, s = expectTokenSlash(s)
	if value == "" ***REMOVED***
		return
	***REMOVED***
	value = strings.ToLower(value)
	s = skipSpace(s)
	for strings.HasPrefix(s, ";") ***REMOVED***
		var pkey string
		pkey, s = expectToken(skipSpace(s[1:]))
		if pkey == "" ***REMOVED***
			return
		***REMOVED***
		if !strings.HasPrefix(s, "=") ***REMOVED***
			return
		***REMOVED***
		var pvalue string
		pvalue, s = expectTokenOrQuoted(s[1:])
		if pvalue == "" ***REMOVED***
			return
		***REMOVED***
		pkey = strings.ToLower(pkey)
		params[pkey] = pvalue
		s = skipSpace(s)
	***REMOVED***
	return
***REMOVED***

// AcceptSpec describes an Accept* header.
type AcceptSpec struct ***REMOVED***
	Value string
	Q     float64
***REMOVED***

// ParseAccept parses Accept* headers.
func ParseAccept(header http.Header, key string) (specs []AcceptSpec) ***REMOVED***
loop:
	for _, s := range header[key] ***REMOVED***
		for ***REMOVED***
			var spec AcceptSpec
			spec.Value, s = expectTokenSlash(s)
			if spec.Value == "" ***REMOVED***
				continue loop
			***REMOVED***
			spec.Q = 1.0
			s = skipSpace(s)
			if strings.HasPrefix(s, ";") ***REMOVED***
				s = skipSpace(s[1:])
				if !strings.HasPrefix(s, "q=") ***REMOVED***
					continue loop
				***REMOVED***
				spec.Q, s = expectQuality(s[2:])
				if spec.Q < 0.0 ***REMOVED***
					continue loop
				***REMOVED***
			***REMOVED***
			specs = append(specs, spec)
			s = skipSpace(s)
			if !strings.HasPrefix(s, ",") ***REMOVED***
				continue loop
			***REMOVED***
			s = skipSpace(s[1:])
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func skipSpace(s string) (rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if octetTypes[s[i]]&isSpace == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[i:]
***REMOVED***

func expectToken(s string) (token, rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if octetTypes[s[i]]&isToken == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[:i], s[i:]
***REMOVED***

func expectTokenSlash(s string) (token, rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		b := s[i]
		if (octetTypes[b]&isToken == 0) && b != '/' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[:i], s[i:]
***REMOVED***

func expectQuality(s string) (q float64, rest string) ***REMOVED***
	switch ***REMOVED***
	case len(s) == 0:
		return -1, ""
	case s[0] == '0':
		q = 0
	case s[0] == '1':
		q = 1
	default:
		return -1, ""
	***REMOVED***
	s = s[1:]
	if !strings.HasPrefix(s, ".") ***REMOVED***
		return q, s
	***REMOVED***
	s = s[1:]
	i := 0
	n := 0
	d := 1
	for ; i < len(s); i++ ***REMOVED***
		b := s[i]
		if b < '0' || b > '9' ***REMOVED***
			break
		***REMOVED***
		n = n*10 + int(b) - '0'
		d *= 10
	***REMOVED***
	return q + float64(n)/float64(d), s[i:]
***REMOVED***

func expectTokenOrQuoted(s string) (value string, rest string) ***REMOVED***
	if !strings.HasPrefix(s, "\"") ***REMOVED***
		return expectToken(s)
	***REMOVED***
	s = s[1:]
	for i := 0; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			p := make([]byte, len(s)-1)
			j := copy(p, s[:i])
			escape := true
			for i = i + 1; i < len(s); i++ ***REMOVED***
				b := s[i]
				switch ***REMOVED***
				case escape:
					escape = false
					p[j] = b
					j++
				case b == '\\':
					escape = true
				case b == '"':
					return string(p[:j]), s[i+1:]
				default:
					p[j] = b
					j++
				***REMOVED***
			***REMOVED***
			return "", ""
		***REMOVED***
	***REMOVED***
	return "", ""
***REMOVED***
