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

// Token octets per RFC 2616.
var isTokenOctet = [256]bool***REMOVED***
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'W':  true,
	'V':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'|':  true,
	'~':  true,
***REMOVED***

// skipSpace returns a slice of the string s with all leading RFC 2616 linear
// whitespace removed.
func skipSpace(s string) (rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if b := s[i]; b != ' ' && b != '\t' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[i:]
***REMOVED***

// nextToken returns the leading RFC 2616 token of s and the string following
// the token.
func nextToken(s string) (token, rest string) ***REMOVED***
	i := 0
	for ; i < len(s); i++ ***REMOVED***
		if !isTokenOctet[s[i]] ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return s[:i], s[i:]
***REMOVED***

// nextTokenOrQuoted returns the leading token or quoted string per RFC 2616
// and the string following the token or quoted string.
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

// equalASCIIFold returns true if s is equal to t with ASCII case folding as
// defined in RFC 4790.
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

// parseExtensions parses WebSocket extensions from a header.
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
