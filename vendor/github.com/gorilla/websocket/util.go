// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

func computeAcceptKey(challengeKey string) string ***REMOVED***
	h := sha1.New()
	h.Write([]byte(challengeKey))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
***REMOVED***

func generateChallengeKey() (string, error) ***REMOVED***
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return base64.StdEncoding.EncodeToString(p), nil
***REMOVED***

// Octet types from RFC 2616.
var octetTypes [256]byte

const (
	isTokenOctet = 1 << iota
	isSpaceOctet
)

func init() ***REMOVED***
	// From RFC 2616
	//
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
		var t byte
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.IndexRune(" \t\"(),/:;<=>?@[]\\***REMOVED******REMOVED***", rune(c)) >= 0
		if strings.IndexRune(" \t\r\n", rune(c)) >= 0 ***REMOVED***
			t |= isSpaceOctet
		***REMOVED***
		if isChar && !isCtl && !isSeparator ***REMOVED***
			t |= isTokenOctet
		***REMOVED***
		octetTypes[c] = t
	***REMOVED***
***REMOVED***

func skipSpace(s string) (rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if octetTypes[s[i]]&isSpaceOctet == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[i:]
***REMOVED***

func nextToken(s string) (token, rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if octetTypes[s[i]]&isTokenOctet == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[:i], s[i:]
***REMOVED***

func nextTokenOrQuoted(s string) (value string, rest string) ***REMOVED***
	if !strings.HasPrefix(s, "\"") ***REMOVED***
		return nextToken(s)
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

// equalASCIIFold returns true if s is equal to t with ASCII case folding.
func equalASCIIFold(s, t string) bool ***REMOVED***
	for s != "" && t != "" ***REMOVED***
		sr, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		tr, size := utf8.DecodeRuneInString(t)
		t = t[size:]
		if sr == tr ***REMOVED***
			continue
		***REMOVED***
		if 'A' <= sr && sr <= 'Z' ***REMOVED***
			sr = sr + 'a' - 'A'
		***REMOVED***
		if 'A' <= tr && tr <= 'Z' ***REMOVED***
			tr = tr + 'a' - 'A'
		***REMOVED***
		if sr != tr ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return s == t
***REMOVED***

// tokenListContainsValue returns true if the 1#token header with the given
// name contains a token equal to value with ASCII case folding.
func tokenListContainsValue(header http.Header, name string, value string) bool ***REMOVED***
headers:
	for _, s := range header[name] ***REMOVED***
		for ***REMOVED***
			var t string
			t, s = nextToken(skipSpace(s))
			if t == "" ***REMOVED***
				continue headers
			***REMOVED***
			s = skipSpace(s)
			if s != "" && s[0] != ',' ***REMOVED***
				continue headers
			***REMOVED***
			if equalASCIIFold(t, value) ***REMOVED***
				return true
			***REMOVED***
			if s == "" ***REMOVED***
				continue headers
			***REMOVED***
			s = s[1:]
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// parseExtensiosn parses WebSocket extensions from a header.
func parseExtensions(header http.Header) []map[string]string ***REMOVED***
	// From RFC 6455:
	//
	//  Sec-WebSocket-Extensions = extension-list
	//  extension-list = 1#extension
	//  extension = extension-token *( ";" extension-param )
	//  extension-token = registered-token
	//  registered-token = token
	//  extension-param = token [ "=" (token | quoted-string) ]
	//     ;When using the quoted-string syntax variant, the value
	//     ;after quoted-string unescaping MUST conform to the
	//     ;'token' ABNF.

	var result []map[string]string
headers:
	for _, s := range header["Sec-Websocket-Extensions"] ***REMOVED***
		for ***REMOVED***
			var t string
			t, s = nextToken(skipSpace(s))
			if t == "" ***REMOVED***
				continue headers
			***REMOVED***
			ext := map[string]string***REMOVED***"": t***REMOVED***
			for ***REMOVED***
				s = skipSpace(s)
				if !strings.HasPrefix(s, ";") ***REMOVED***
					break
				***REMOVED***
				var k string
				k, s = nextToken(skipSpace(s[1:]))
				if k == "" ***REMOVED***
					continue headers
				***REMOVED***
				s = skipSpace(s)
				var v string
				if strings.HasPrefix(s, "=") ***REMOVED***
					v, s = nextTokenOrQuoted(skipSpace(s[1:]))
					s = skipSpace(s)
				***REMOVED***
				if s != "" && s[0] != ',' && s[0] != ';' ***REMOVED***
					continue headers
				***REMOVED***
				ext[k] = v
			***REMOVED***
			if s != "" && s[0] != ',' ***REMOVED***
				continue headers
			***REMOVED***
			result = append(result, ext)
			if s == "" ***REMOVED***
				continue headers
			***REMOVED***
			s = s[1:]
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***
