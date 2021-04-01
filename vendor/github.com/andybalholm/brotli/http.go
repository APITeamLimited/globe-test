package brotli

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// HTTPCompressor chooses a compression method (brotli, gzip, or none) based on
// the Accept-Encoding header, sets the Content-Encoding header, and returns a
// WriteCloser that implements that compression. The Close method must be called
// before the current HTTP handler returns.
//
// Due to https://github.com/golang/go/issues/31753, the response will not be
// compressed unless you set a Content-Type header before you call
// HTTPCompressor.
func HTTPCompressor(w http.ResponseWriter, r *http.Request) io.WriteCloser ***REMOVED***
	if w.Header().Get("Content-Type") == "" ***REMOVED***
		return nopCloser***REMOVED***w***REMOVED***
	***REMOVED***

	if w.Header().Get("Vary") == "" ***REMOVED***
		w.Header().Set("Vary", "Accept-Encoding")
	***REMOVED***

	encoding := negotiateContentEncoding(r, []string***REMOVED***"br", "gzip"***REMOVED***)
	switch encoding ***REMOVED***
	case "br":
		w.Header().Set("Content-Encoding", "br")
		return NewWriter(w)
	case "gzip":
		w.Header().Set("Content-Encoding", "gzip")
		return gzip.NewWriter(w)
	***REMOVED***
	return nopCloser***REMOVED***w***REMOVED***
***REMOVED***

// negotiateContentEncoding returns the best offered content encoding for the
// request's Accept-Encoding header. If two offers match with equal weight and
// then the offer earlier in the list is preferred. If no offers are
// acceptable, then "" is returned.
func negotiateContentEncoding(r *http.Request, offers []string) string ***REMOVED***
	bestOffer := "identity"
	bestQ := -1.0
	specs := parseAccept(r.Header, "Accept-Encoding")
	for _, offer := range offers ***REMOVED***
		for _, spec := range specs ***REMOVED***
			if spec.Q > bestQ &&
				(spec.Value == "*" || spec.Value == offer) ***REMOVED***
				bestQ = spec.Q
				bestOffer = offer
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if bestQ == 0 ***REMOVED***
		bestOffer = ""
	***REMOVED***
	return bestOffer
***REMOVED***

// acceptSpec describes an Accept* header.
type acceptSpec struct ***REMOVED***
	Value string
	Q     float64
***REMOVED***

// parseAccept parses Accept* headers.
func parseAccept(header http.Header, key string) (specs []acceptSpec) ***REMOVED***
loop:
	for _, s := range header[key] ***REMOVED***
		for ***REMOVED***
			var spec acceptSpec
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
