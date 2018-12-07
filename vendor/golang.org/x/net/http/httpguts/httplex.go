// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpguts

import (
	"net"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

var isTokenTable = [127]bool***REMOVED***
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

func IsTokenRune(r rune) bool ***REMOVED***
	i := int(r)
	return i < len(isTokenTable) && isTokenTable[i]
***REMOVED***

func isNotToken(r rune) bool ***REMOVED***
	return !IsTokenRune(r)
***REMOVED***

// HeaderValuesContainsToken reports whether any string in values
// contains the provided token, ASCII case-insensitively.
func HeaderValuesContainsToken(values []string, token string) bool ***REMOVED***
	for _, v := range values ***REMOVED***
		if headerValueContainsToken(v, token) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// isOWS reports whether b is an optional whitespace byte, as defined
// by RFC 7230 section 3.2.3.
func isOWS(b byte) bool ***REMOVED*** return b == ' ' || b == '\t' ***REMOVED***

// trimOWS returns x with all optional whitespace removes from the
// beginning and end.
func trimOWS(x string) string ***REMOVED***
	// TODO: consider using strings.Trim(x, " \t") instead,
	// if and when it's fast enough. See issue 10292.
	// But this ASCII-only code will probably always beat UTF-8
	// aware code.
	for len(x) > 0 && isOWS(x[0]) ***REMOVED***
		x = x[1:]
	***REMOVED***
	for len(x) > 0 && isOWS(x[len(x)-1]) ***REMOVED***
		x = x[:len(x)-1]
	***REMOVED***
	return x
***REMOVED***

// headerValueContainsToken reports whether v (assumed to be a
// 0#element, in the ABNF extension described in RFC 7230 section 7)
// contains token amongst its comma-separated tokens, ASCII
// case-insensitively.
func headerValueContainsToken(v string, token string) bool ***REMOVED***
	v = trimOWS(v)
	if comma := strings.IndexByte(v, ','); comma != -1 ***REMOVED***
		return tokenEqual(trimOWS(v[:comma]), token) || headerValueContainsToken(v[comma+1:], token)
	***REMOVED***
	return tokenEqual(v, token)
***REMOVED***

// lowerASCII returns the ASCII lowercase version of b.
func lowerASCII(b byte) byte ***REMOVED***
	if 'A' <= b && b <= 'Z' ***REMOVED***
		return b + ('a' - 'A')
	***REMOVED***
	return b
***REMOVED***

// tokenEqual reports whether t1 and t2 are equal, ASCII case-insensitively.
func tokenEqual(t1, t2 string) bool ***REMOVED***
	if len(t1) != len(t2) ***REMOVED***
		return false
	***REMOVED***
	for i, b := range t1 ***REMOVED***
		if b >= utf8.RuneSelf ***REMOVED***
			// No UTF-8 or non-ASCII allowed in tokens.
			return false
		***REMOVED***
		if lowerASCII(byte(b)) != lowerASCII(t2[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// isLWS reports whether b is linear white space, according
// to http://www.w3.org/Protocols/rfc2616/rfc2616-sec2.html#sec2.2
//      LWS            = [CRLF] 1*( SP | HT )
func isLWS(b byte) bool ***REMOVED*** return b == ' ' || b == '\t' ***REMOVED***

// isCTL reports whether b is a control byte, according
// to http://www.w3.org/Protocols/rfc2616/rfc2616-sec2.html#sec2.2
//      CTL            = <any US-ASCII control character
//                       (octets 0 - 31) and DEL (127)>
func isCTL(b byte) bool ***REMOVED***
	const del = 0x7f // a CTL
	return b < ' ' || b == del
***REMOVED***

// ValidHeaderFieldName reports whether v is a valid HTTP/1.x header name.
// HTTP/2 imposes the additional restriction that uppercase ASCII
// letters are not allowed.
//
//  RFC 7230 says:
//   header-field   = field-name ":" OWS field-value OWS
//   field-name     = token
//   token          = 1*tchar
//   tchar = "!" / "#" / "$" / "%" / "&" / "'" / "*" / "+" / "-" / "." /
//           "^" / "_" / "`" / "|" / "~" / DIGIT / ALPHA
func ValidHeaderFieldName(v string) bool ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return false
	***REMOVED***
	for _, r := range v ***REMOVED***
		if !IsTokenRune(r) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// ValidHostHeader reports whether h is a valid host header.
func ValidHostHeader(h string) bool ***REMOVED***
	// The latest spec is actually this:
	//
	// http://tools.ietf.org/html/rfc7230#section-5.4
	//     Host = uri-host [ ":" port ]
	//
	// Where uri-host is:
	//     http://tools.ietf.org/html/rfc3986#section-3.2.2
	//
	// But we're going to be much more lenient for now and just
	// search for any byte that's not a valid byte in any of those
	// expressions.
	for i := 0; i < len(h); i++ ***REMOVED***
		if !validHostByte[h[i]] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// See the validHostHeader comment.
var validHostByte = [256]bool***REMOVED***
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true,
	'8': true, '9': true,

	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true,
	'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true,
	'q': true, 'r': true, 's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true,
	'y': true, 'z': true,

	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true, 'H': true,
	'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true, 'O': true, 'P': true,
	'Q': true, 'R': true, 'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true,
	'Y': true, 'Z': true,

	'!':  true, // sub-delims
	'$':  true, // sub-delims
	'%':  true, // pct-encoded (and used in IPv6 zones)
	'&':  true, // sub-delims
	'(':  true, // sub-delims
	')':  true, // sub-delims
	'*':  true, // sub-delims
	'+':  true, // sub-delims
	',':  true, // sub-delims
	'-':  true, // unreserved
	'.':  true, // unreserved
	':':  true, // IPv6address + Host expression's optional port
	';':  true, // sub-delims
	'=':  true, // sub-delims
	'[':  true,
	'\'': true, // sub-delims
	']':  true,
	'_':  true, // unreserved
	'~':  true, // unreserved
***REMOVED***

// ValidHeaderFieldValue reports whether v is a valid "field-value" according to
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2 :
//
//        message-header = field-name ":" [ field-value ]
//        field-value    = *( field-content | LWS )
//        field-content  = <the OCTETs making up the field-value
//                         and consisting of either *TEXT or combinations
//                         of token, separators, and quoted-string>
//
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec2.html#sec2.2 :
//
//        TEXT           = <any OCTET except CTLs,
//                          but including LWS>
//        LWS            = [CRLF] 1*( SP | HT )
//        CTL            = <any US-ASCII control character
//                         (octets 0 - 31) and DEL (127)>
//
// RFC 7230 says:
//  field-value    = *( field-content / obs-fold )
//  obj-fold       =  N/A to http2, and deprecated
//  field-content  = field-vchar [ 1*( SP / HTAB ) field-vchar ]
//  field-vchar    = VCHAR / obs-text
//  obs-text       = %x80-FF
//  VCHAR          = "any visible [USASCII] character"
//
// http2 further says: "Similarly, HTTP/2 allows header field values
// that are not valid. While most of the values that can be encoded
// will not alter header field parsing, carriage return (CR, ASCII
// 0xd), line feed (LF, ASCII 0xa), and the zero character (NUL, ASCII
// 0x0) might be exploited by an attacker if they are translated
// verbatim. Any request or response that contains a character not
// permitted in a header field value MUST be treated as malformed
// (Section 8.1.2.6). Valid characters are defined by the
// field-content ABNF rule in Section 3.2 of [RFC7230]."
//
// This function does not (yet?) properly handle the rejection of
// strings that begin or end with SP or HTAB.
func ValidHeaderFieldValue(v string) bool ***REMOVED***
	for i := 0; i < len(v); i++ ***REMOVED***
		b := v[i]
		if isCTL(b) && !isLWS(b) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func isASCII(s string) bool ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] >= utf8.RuneSelf ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// PunycodeHostPort returns the IDNA Punycode version
// of the provided "host" or "host:port" string.
func PunycodeHostPort(v string) (string, error) ***REMOVED***
	if isASCII(v) ***REMOVED***
		return v, nil
	***REMOVED***

	host, port, err := net.SplitHostPort(v)
	if err != nil ***REMOVED***
		// The input 'v' argument was just a "host" argument,
		// without a port. This error should not be returned
		// to the caller.
		host = v
		port = ""
	***REMOVED***
	host, err = idna.ToASCII(host)
	if err != nil ***REMOVED***
		// Non-UTF-8? Not representable in Punycode, in any
		// case.
		return "", err
	***REMOVED***
	if port == "" ***REMOVED***
		return host, nil
	***REMOVED***
	return net.JoinHostPort(host, port), nil
***REMOVED***
