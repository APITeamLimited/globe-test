// Package jlexer contains a JSON lexer implementation.
//
// It is expected that it is mostly used with generated parser code, so the interface is tuned
// for a parser that knows what kind of data is expected.
package jlexer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/josharian/intern"
)

// tokenKind determines type of a token.
type tokenKind byte

const (
	tokenUndef  tokenKind = iota // No token.
	tokenDelim                   // Delimiter: one of '***REMOVED***', '***REMOVED***', '[' or ']'.
	tokenString                  // A string literal, e.g. "abc\u1234"
	tokenNumber                  // Number literal, e.g. 1.5e5
	tokenBool                    // Boolean literal: true or false.
	tokenNull                    // null keyword.
)

// token describes a single token: type, position in the input and value.
type token struct ***REMOVED***
	kind tokenKind // Type of a token.

	boolValue       bool   // Value if a boolean literal token.
	byteValueCloned bool   // true if byteValue was allocated and does not refer to original json body
	byteValue       []byte // Raw value of a token.
	delimValue      byte
***REMOVED***

// Lexer is a JSON lexer: it iterates over JSON tokens in a byte slice.
type Lexer struct ***REMOVED***
	Data []byte // Input data given to the lexer.

	start int   // Start of the current token.
	pos   int   // Current unscanned position in the input stream.
	token token // Last scanned token, if token.kind != tokenUndef.

	firstElement bool // Whether current element is the first in array or an object.
	wantSep      byte // A comma or a colon character, which need to occur before a token.

	UseMultipleErrors bool          // If we want to use multiple errors.
	fatalError        error         // Fatal error occurred during lexing. It is usually a syntax error.
	multipleErrors    []*LexerError // Semantic errors occurred during lexing. Marshalling will be continued after finding this errors.
***REMOVED***

// FetchToken scans the input for the next token.
func (r *Lexer) FetchToken() ***REMOVED***
	r.token.kind = tokenUndef
	r.start = r.pos

	// Check if r.Data has r.pos element
	// If it doesn't, it mean corrupted input data
	if len(r.Data) < r.pos ***REMOVED***
		r.errParse("Unexpected end of data")
		return
	***REMOVED***
	// Determine the type of a token by skipping whitespace and reading the
	// first character.
	for _, c := range r.Data[r.pos:] ***REMOVED***
		switch c ***REMOVED***
		case ':', ',':
			if r.wantSep == c ***REMOVED***
				r.pos++
				r.start++
				r.wantSep = 0
			***REMOVED*** else ***REMOVED***
				r.errSyntax()
			***REMOVED***

		case ' ', '\t', '\r', '\n':
			r.pos++
			r.start++

		case '"':
			if r.wantSep != 0 ***REMOVED***
				r.errSyntax()
			***REMOVED***

			r.token.kind = tokenString
			r.fetchString()
			return

		case '***REMOVED***', '[':
			if r.wantSep != 0 ***REMOVED***
				r.errSyntax()
			***REMOVED***
			r.firstElement = true
			r.token.kind = tokenDelim
			r.token.delimValue = r.Data[r.pos]
			r.pos++
			return

		case '***REMOVED***', ']':
			if !r.firstElement && (r.wantSep != ',') ***REMOVED***
				r.errSyntax()
			***REMOVED***
			r.wantSep = 0
			r.token.kind = tokenDelim
			r.token.delimValue = r.Data[r.pos]
			r.pos++
			return

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			if r.wantSep != 0 ***REMOVED***
				r.errSyntax()
			***REMOVED***
			r.token.kind = tokenNumber
			r.fetchNumber()
			return

		case 'n':
			if r.wantSep != 0 ***REMOVED***
				r.errSyntax()
			***REMOVED***

			r.token.kind = tokenNull
			r.fetchNull()
			return

		case 't':
			if r.wantSep != 0 ***REMOVED***
				r.errSyntax()
			***REMOVED***

			r.token.kind = tokenBool
			r.token.boolValue = true
			r.fetchTrue()
			return

		case 'f':
			if r.wantSep != 0 ***REMOVED***
				r.errSyntax()
			***REMOVED***

			r.token.kind = tokenBool
			r.token.boolValue = false
			r.fetchFalse()
			return

		default:
			r.errSyntax()
			return
		***REMOVED***
	***REMOVED***
	r.fatalError = io.EOF
	return
***REMOVED***

// isTokenEnd returns true if the char can follow a non-delimiter token
func isTokenEnd(c byte) bool ***REMOVED***
	return c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '[' || c == ']' || c == '***REMOVED***' || c == '***REMOVED***' || c == ',' || c == ':'
***REMOVED***

// fetchNull fetches and checks remaining bytes of null keyword.
func (r *Lexer) fetchNull() ***REMOVED***
	r.pos += 4
	if r.pos > len(r.Data) ||
		r.Data[r.pos-3] != 'u' ||
		r.Data[r.pos-2] != 'l' ||
		r.Data[r.pos-1] != 'l' ||
		(r.pos != len(r.Data) && !isTokenEnd(r.Data[r.pos])) ***REMOVED***

		r.pos -= 4
		r.errSyntax()
	***REMOVED***
***REMOVED***

// fetchTrue fetches and checks remaining bytes of true keyword.
func (r *Lexer) fetchTrue() ***REMOVED***
	r.pos += 4
	if r.pos > len(r.Data) ||
		r.Data[r.pos-3] != 'r' ||
		r.Data[r.pos-2] != 'u' ||
		r.Data[r.pos-1] != 'e' ||
		(r.pos != len(r.Data) && !isTokenEnd(r.Data[r.pos])) ***REMOVED***

		r.pos -= 4
		r.errSyntax()
	***REMOVED***
***REMOVED***

// fetchFalse fetches and checks remaining bytes of false keyword.
func (r *Lexer) fetchFalse() ***REMOVED***
	r.pos += 5
	if r.pos > len(r.Data) ||
		r.Data[r.pos-4] != 'a' ||
		r.Data[r.pos-3] != 'l' ||
		r.Data[r.pos-2] != 's' ||
		r.Data[r.pos-1] != 'e' ||
		(r.pos != len(r.Data) && !isTokenEnd(r.Data[r.pos])) ***REMOVED***

		r.pos -= 5
		r.errSyntax()
	***REMOVED***
***REMOVED***

// fetchNumber scans a number literal token.
func (r *Lexer) fetchNumber() ***REMOVED***
	hasE := false
	afterE := false
	hasDot := false

	r.pos++
	for i, c := range r.Data[r.pos:] ***REMOVED***
		switch ***REMOVED***
		case c >= '0' && c <= '9':
			afterE = false
		case c == '.' && !hasDot:
			hasDot = true
		case (c == 'e' || c == 'E') && !hasE:
			hasE = true
			hasDot = true
			afterE = true
		case (c == '+' || c == '-') && afterE:
			afterE = false
		default:
			r.pos += i
			if !isTokenEnd(c) ***REMOVED***
				r.errSyntax()
			***REMOVED*** else ***REMOVED***
				r.token.byteValue = r.Data[r.start:r.pos]
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	r.pos = len(r.Data)
	r.token.byteValue = r.Data[r.start:]
***REMOVED***

// findStringLen tries to scan into the string literal for ending quote char to determine required size.
// The size will be exact if no escapes are present and may be inexact if there are escaped chars.
func findStringLen(data []byte) (isValid bool, length int) ***REMOVED***
	for ***REMOVED***
		idx := bytes.IndexByte(data, '"')
		if idx == -1 ***REMOVED***
			return false, len(data)
		***REMOVED***
		if idx == 0 || (idx > 0 && data[idx-1] != '\\') ***REMOVED***
			return true, length + idx
		***REMOVED***

		// count \\\\\\\ sequences. even number of slashes means quote is not really escaped
		cnt := 1
		for idx-cnt-1 >= 0 && data[idx-cnt-1] == '\\' ***REMOVED***
			cnt++
		***REMOVED***
		if cnt%2 == 0 ***REMOVED***
			return true, length + idx
		***REMOVED***

		length += idx + 1
		data = data[idx+1:]
	***REMOVED***
***REMOVED***

// unescapeStringToken performs unescaping of string token.
// if no escaping is needed, original string is returned, otherwise - a new one allocated
func (r *Lexer) unescapeStringToken() (err error) ***REMOVED***
	data := r.token.byteValue
	var unescapedData []byte

	for ***REMOVED***
		i := bytes.IndexByte(data, '\\')
		if i == -1 ***REMOVED***
			break
		***REMOVED***

		escapedRune, escapedBytes, err := decodeEscape(data[i:])
		if err != nil ***REMOVED***
			r.errParse(err.Error())
			return err
		***REMOVED***

		if unescapedData == nil ***REMOVED***
			unescapedData = make([]byte, 0, len(r.token.byteValue))
		***REMOVED***

		var d [4]byte
		s := utf8.EncodeRune(d[:], escapedRune)
		unescapedData = append(unescapedData, data[:i]...)
		unescapedData = append(unescapedData, d[:s]...)

		data = data[i+escapedBytes:]
	***REMOVED***

	if unescapedData != nil ***REMOVED***
		r.token.byteValue = append(unescapedData, data...)
		r.token.byteValueCloned = true
	***REMOVED***
	return
***REMOVED***

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune ***REMOVED***
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' ***REMOVED***
		return -1
	***REMOVED***
	var val rune
	for i := 2; i < len(s) && i < 6; i++ ***REMOVED***
		var v byte
		c := s[i]
		switch c ***REMOVED***
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			v = c - '0'
		case 'a', 'b', 'c', 'd', 'e', 'f':
			v = c - 'a' + 10
		case 'A', 'B', 'C', 'D', 'E', 'F':
			v = c - 'A' + 10
		default:
			return -1
		***REMOVED***

		val <<= 4
		val |= rune(v)
	***REMOVED***
	return val
***REMOVED***

// decodeEscape processes a single escape sequence and returns number of bytes processed.
func decodeEscape(data []byte) (decoded rune, bytesProcessed int, err error) ***REMOVED***
	if len(data) < 2 ***REMOVED***
		return 0, 0, errors.New("incorrect escape symbol \\ at the end of token")
	***REMOVED***

	c := data[1]
	switch c ***REMOVED***
	case '"', '/', '\\':
		return rune(c), 2, nil
	case 'b':
		return '\b', 2, nil
	case 'f':
		return '\f', 2, nil
	case 'n':
		return '\n', 2, nil
	case 'r':
		return '\r', 2, nil
	case 't':
		return '\t', 2, nil
	case 'u':
		rr := getu4(data)
		if rr < 0 ***REMOVED***
			return 0, 0, errors.New("incorrectly escaped \\uXXXX sequence")
		***REMOVED***

		read := 6
		if utf16.IsSurrogate(rr) ***REMOVED***
			rr1 := getu4(data[read:])
			if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar ***REMOVED***
				read += 6
				rr = dec
			***REMOVED*** else ***REMOVED***
				rr = unicode.ReplacementChar
			***REMOVED***
		***REMOVED***
		return rr, read, nil
	***REMOVED***

	return 0, 0, errors.New("incorrectly escaped bytes")
***REMOVED***

// fetchString scans a string literal token.
func (r *Lexer) fetchString() ***REMOVED***
	r.pos++
	data := r.Data[r.pos:]

	isValid, length := findStringLen(data)
	if !isValid ***REMOVED***
		r.pos += length
		r.errParse("unterminated string literal")
		return
	***REMOVED***
	r.token.byteValue = data[:length]
	r.pos += length + 1 // skip closing '"' as well
***REMOVED***

// scanToken scans the next token if no token is currently available in the lexer.
func (r *Lexer) scanToken() ***REMOVED***
	if r.token.kind != tokenUndef || r.fatalError != nil ***REMOVED***
		return
	***REMOVED***

	r.FetchToken()
***REMOVED***

// consume resets the current token to allow scanning the next one.
func (r *Lexer) consume() ***REMOVED***
	r.token.kind = tokenUndef
	r.token.delimValue = 0
***REMOVED***

// Ok returns true if no error (including io.EOF) was encountered during scanning.
func (r *Lexer) Ok() bool ***REMOVED***
	return r.fatalError == nil
***REMOVED***

const maxErrorContextLen = 13

func (r *Lexer) errParse(what string) ***REMOVED***
	if r.fatalError == nil ***REMOVED***
		var str string
		if len(r.Data)-r.pos <= maxErrorContextLen ***REMOVED***
			str = string(r.Data)
		***REMOVED*** else ***REMOVED***
			str = string(r.Data[r.pos:r.pos+maxErrorContextLen-3]) + "..."
		***REMOVED***
		r.fatalError = &LexerError***REMOVED***
			Reason: what,
			Offset: r.pos,
			Data:   str,
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Lexer) errSyntax() ***REMOVED***
	r.errParse("syntax error")
***REMOVED***

func (r *Lexer) errInvalidToken(expected string) ***REMOVED***
	if r.fatalError != nil ***REMOVED***
		return
	***REMOVED***
	if r.UseMultipleErrors ***REMOVED***
		r.pos = r.start
		r.consume()
		r.SkipRecursive()
		switch expected ***REMOVED***
		case "[":
			r.token.delimValue = ']'
			r.token.kind = tokenDelim
		case "***REMOVED***":
			r.token.delimValue = '***REMOVED***'
			r.token.kind = tokenDelim
		***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Reason: fmt.Sprintf("expected %s", expected),
			Offset: r.start,
			Data:   string(r.Data[r.start:r.pos]),
		***REMOVED***)
		return
	***REMOVED***

	var str string
	if len(r.token.byteValue) <= maxErrorContextLen ***REMOVED***
		str = string(r.token.byteValue)
	***REMOVED*** else ***REMOVED***
		str = string(r.token.byteValue[:maxErrorContextLen-3]) + "..."
	***REMOVED***
	r.fatalError = &LexerError***REMOVED***
		Reason: fmt.Sprintf("expected %s", expected),
		Offset: r.pos,
		Data:   str,
	***REMOVED***
***REMOVED***

func (r *Lexer) GetPos() int ***REMOVED***
	return r.pos
***REMOVED***

// Delim consumes a token and verifies that it is the given delimiter.
func (r *Lexer) Delim(c byte) ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***

	if !r.Ok() || r.token.delimValue != c ***REMOVED***
		r.consume() // errInvalidToken can change token if UseMultipleErrors is enabled.
		r.errInvalidToken(string([]byte***REMOVED***c***REMOVED***))
	***REMOVED*** else ***REMOVED***
		r.consume()
	***REMOVED***
***REMOVED***

// IsDelim returns true if there was no scanning error and next token is the given delimiter.
func (r *Lexer) IsDelim(c byte) bool ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	return !r.Ok() || r.token.delimValue == c
***REMOVED***

// Null verifies that the next token is null and consumes it.
func (r *Lexer) Null() ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() || r.token.kind != tokenNull ***REMOVED***
		r.errInvalidToken("null")
	***REMOVED***
	r.consume()
***REMOVED***

// IsNull returns true if the next token is a null keyword.
func (r *Lexer) IsNull() bool ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	return r.Ok() && r.token.kind == tokenNull
***REMOVED***

// Skip skips a single token.
func (r *Lexer) Skip() ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	r.consume()
***REMOVED***

// SkipRecursive skips next array or object completely, or just skips a single token if not
// an array/object.
//
// Note: no syntax validation is performed on the skipped data.
func (r *Lexer) SkipRecursive() ***REMOVED***
	r.scanToken()
	var start, end byte

	switch r.token.delimValue ***REMOVED***
	case '***REMOVED***':
		start, end = '***REMOVED***', '***REMOVED***'
	case '[':
		start, end = '[', ']'
	default:
		r.consume()
		return
	***REMOVED***

	r.consume()

	level := 1
	inQuotes := false
	wasEscape := false

	for i, c := range r.Data[r.pos:] ***REMOVED***
		switch ***REMOVED***
		case c == start && !inQuotes:
			level++
		case c == end && !inQuotes:
			level--
			if level == 0 ***REMOVED***
				r.pos += i + 1
				return
			***REMOVED***
		case c == '\\' && inQuotes:
			wasEscape = !wasEscape
			continue
		case c == '"' && inQuotes:
			inQuotes = wasEscape
		case c == '"':
			inQuotes = true
		***REMOVED***
		wasEscape = false
	***REMOVED***
	r.pos = len(r.Data)
	r.fatalError = &LexerError***REMOVED***
		Reason: "EOF reached while skipping array/object or token",
		Offset: r.pos,
		Data:   string(r.Data[r.pos:]),
	***REMOVED***
***REMOVED***

// Raw fetches the next item recursively as a data slice
func (r *Lexer) Raw() []byte ***REMOVED***
	r.SkipRecursive()
	if !r.Ok() ***REMOVED***
		return nil
	***REMOVED***
	return r.Data[r.start:r.pos]
***REMOVED***

// IsStart returns whether the lexer is positioned at the start
// of an input string.
func (r *Lexer) IsStart() bool ***REMOVED***
	return r.pos == 0
***REMOVED***

// Consumed reads all remaining bytes from the input, publishing an error if
// there is anything but whitespace remaining.
func (r *Lexer) Consumed() ***REMOVED***
	if r.pos > len(r.Data) || !r.Ok() ***REMOVED***
		return
	***REMOVED***

	for _, c := range r.Data[r.pos:] ***REMOVED***
		if c != ' ' && c != '\t' && c != '\r' && c != '\n' ***REMOVED***
			r.AddError(&LexerError***REMOVED***
				Reason: "invalid character '" + string(c) + "' after top-level value",
				Offset: r.pos,
				Data:   string(r.Data[r.pos:]),
			***REMOVED***)
			return
		***REMOVED***

		r.pos++
		r.start++
	***REMOVED***
***REMOVED***

func (r *Lexer) unsafeString(skipUnescape bool) (string, []byte) ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() || r.token.kind != tokenString ***REMOVED***
		r.errInvalidToken("string")
		return "", nil
	***REMOVED***
	if !skipUnescape ***REMOVED***
		if err := r.unescapeStringToken(); err != nil ***REMOVED***
			r.errInvalidToken("string")
			return "", nil
		***REMOVED***
	***REMOVED***

	bytes := r.token.byteValue
	ret := bytesToStr(r.token.byteValue)
	r.consume()
	return ret, bytes
***REMOVED***

// UnsafeString returns the string value if the token is a string literal.
//
// Warning: returned string may point to the input buffer, so the string should not outlive
// the input buffer. Intended pattern of usage is as an argument to a switch statement.
func (r *Lexer) UnsafeString() string ***REMOVED***
	ret, _ := r.unsafeString(false)
	return ret
***REMOVED***

// UnsafeBytes returns the byte slice if the token is a string literal.
func (r *Lexer) UnsafeBytes() []byte ***REMOVED***
	_, ret := r.unsafeString(false)
	return ret
***REMOVED***

// UnsafeFieldName returns current member name string token
func (r *Lexer) UnsafeFieldName(skipUnescape bool) string ***REMOVED***
	ret, _ := r.unsafeString(skipUnescape)
	return ret
***REMOVED***

// String reads a string literal.
func (r *Lexer) String() string ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() || r.token.kind != tokenString ***REMOVED***
		r.errInvalidToken("string")
		return ""
	***REMOVED***
	if err := r.unescapeStringToken(); err != nil ***REMOVED***
		r.errInvalidToken("string")
		return ""
	***REMOVED***
	var ret string
	if r.token.byteValueCloned ***REMOVED***
		ret = bytesToStr(r.token.byteValue)
	***REMOVED*** else ***REMOVED***
		ret = string(r.token.byteValue)
	***REMOVED***
	r.consume()
	return ret
***REMOVED***

// StringIntern reads a string literal, and performs string interning on it.
func (r *Lexer) StringIntern() string ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() || r.token.kind != tokenString ***REMOVED***
		r.errInvalidToken("string")
		return ""
	***REMOVED***
	if err := r.unescapeStringToken(); err != nil ***REMOVED***
		r.errInvalidToken("string")
		return ""
	***REMOVED***
	ret := intern.Bytes(r.token.byteValue)
	r.consume()
	return ret
***REMOVED***

// Bytes reads a string literal and base64 decodes it into a byte slice.
func (r *Lexer) Bytes() []byte ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() || r.token.kind != tokenString ***REMOVED***
		r.errInvalidToken("string")
		return nil
	***REMOVED***
	ret := make([]byte, base64.StdEncoding.DecodedLen(len(r.token.byteValue)))
	n, err := base64.StdEncoding.Decode(ret, r.token.byteValue)
	if err != nil ***REMOVED***
		r.fatalError = &LexerError***REMOVED***
			Reason: err.Error(),
		***REMOVED***
		return nil
	***REMOVED***

	r.consume()
	return ret[:n]
***REMOVED***

// Bool reads a true or false boolean keyword.
func (r *Lexer) Bool() bool ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() || r.token.kind != tokenBool ***REMOVED***
		r.errInvalidToken("bool")
		return false
	***REMOVED***
	ret := r.token.boolValue
	r.consume()
	return ret
***REMOVED***

func (r *Lexer) number() string ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() || r.token.kind != tokenNumber ***REMOVED***
		r.errInvalidToken("number")
		return ""
	***REMOVED***
	ret := bytesToStr(r.token.byteValue)
	r.consume()
	return ret
***REMOVED***

func (r *Lexer) Uint8() uint8 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 8)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return uint8(n)
***REMOVED***

func (r *Lexer) Uint16() uint16 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return uint16(n)
***REMOVED***

func (r *Lexer) Uint32() uint32 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return uint32(n)
***REMOVED***

func (r *Lexer) Uint64() uint64 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return n
***REMOVED***

func (r *Lexer) Uint() uint ***REMOVED***
	return uint(r.Uint64())
***REMOVED***

func (r *Lexer) Int8() int8 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 8)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return int8(n)
***REMOVED***

func (r *Lexer) Int16() int16 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 16)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return int16(n)
***REMOVED***

func (r *Lexer) Int32() int32 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return int32(n)
***REMOVED***

func (r *Lexer) Int64() int64 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return n
***REMOVED***

func (r *Lexer) Int() int ***REMOVED***
	return int(r.Int64())
***REMOVED***

func (r *Lexer) Uint8Str() uint8 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 8)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return uint8(n)
***REMOVED***

func (r *Lexer) Uint16Str() uint16 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return uint16(n)
***REMOVED***

func (r *Lexer) Uint32Str() uint32 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return uint32(n)
***REMOVED***

func (r *Lexer) Uint64Str() uint64 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return n
***REMOVED***

func (r *Lexer) UintStr() uint ***REMOVED***
	return uint(r.Uint64Str())
***REMOVED***

func (r *Lexer) UintptrStr() uintptr ***REMOVED***
	return uintptr(r.Uint64Str())
***REMOVED***

func (r *Lexer) Int8Str() int8 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 8)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return int8(n)
***REMOVED***

func (r *Lexer) Int16Str() int16 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 16)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return int16(n)
***REMOVED***

func (r *Lexer) Int32Str() int32 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return int32(n)
***REMOVED***

func (r *Lexer) Int64Str() int64 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return n
***REMOVED***

func (r *Lexer) IntStr() int ***REMOVED***
	return int(r.Int64Str())
***REMOVED***

func (r *Lexer) Float32() float32 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseFloat(s, 32)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return float32(n)
***REMOVED***

func (r *Lexer) Float32Str() float32 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***
	n, err := strconv.ParseFloat(s, 32)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return float32(n)
***REMOVED***

func (r *Lexer) Float64() float64 ***REMOVED***
	s := r.number()
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***

	n, err := strconv.ParseFloat(s, 64)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   s,
		***REMOVED***)
	***REMOVED***
	return n
***REMOVED***

func (r *Lexer) Float64Str() float64 ***REMOVED***
	s, b := r.unsafeString(false)
	if !r.Ok() ***REMOVED***
		return 0
	***REMOVED***
	n, err := strconv.ParseFloat(s, 64)
	if err != nil ***REMOVED***
		r.addNonfatalError(&LexerError***REMOVED***
			Offset: r.start,
			Reason: err.Error(),
			Data:   string(b),
		***REMOVED***)
	***REMOVED***
	return n
***REMOVED***

func (r *Lexer) Error() error ***REMOVED***
	return r.fatalError
***REMOVED***

func (r *Lexer) AddError(e error) ***REMOVED***
	if r.fatalError == nil ***REMOVED***
		r.fatalError = e
	***REMOVED***
***REMOVED***

func (r *Lexer) AddNonFatalError(e error) ***REMOVED***
	r.addNonfatalError(&LexerError***REMOVED***
		Offset: r.start,
		Data:   string(r.Data[r.start:r.pos]),
		Reason: e.Error(),
	***REMOVED***)
***REMOVED***

func (r *Lexer) addNonfatalError(err *LexerError) ***REMOVED***
	if r.UseMultipleErrors ***REMOVED***
		// We don't want to add errors with the same offset.
		if len(r.multipleErrors) != 0 && r.multipleErrors[len(r.multipleErrors)-1].Offset == err.Offset ***REMOVED***
			return
		***REMOVED***
		r.multipleErrors = append(r.multipleErrors, err)
		return
	***REMOVED***
	r.fatalError = err
***REMOVED***

func (r *Lexer) GetNonFatalErrors() []*LexerError ***REMOVED***
	return r.multipleErrors
***REMOVED***

// JsonNumber fetches and json.Number from 'encoding/json' package.
// Both int, float or string, contains them are valid values
func (r *Lexer) JsonNumber() json.Number ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***
	if !r.Ok() ***REMOVED***
		r.errInvalidToken("json.Number")
		return json.Number("")
	***REMOVED***

	switch r.token.kind ***REMOVED***
	case tokenString:
		return json.Number(r.String())
	case tokenNumber:
		return json.Number(r.Raw())
	case tokenNull:
		r.Null()
		return json.Number("")
	default:
		r.errSyntax()
		return json.Number("")
	***REMOVED***
***REMOVED***

// Interface fetches an interface***REMOVED******REMOVED*** analogous to the 'encoding/json' package.
func (r *Lexer) Interface() interface***REMOVED******REMOVED*** ***REMOVED***
	if r.token.kind == tokenUndef && r.Ok() ***REMOVED***
		r.FetchToken()
	***REMOVED***

	if !r.Ok() ***REMOVED***
		return nil
	***REMOVED***
	switch r.token.kind ***REMOVED***
	case tokenString:
		return r.String()
	case tokenNumber:
		return r.Float64()
	case tokenBool:
		return r.Bool()
	case tokenNull:
		r.Null()
		return nil
	***REMOVED***

	if r.token.delimValue == '***REMOVED***' ***REMOVED***
		r.consume()

		ret := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for !r.IsDelim('***REMOVED***') ***REMOVED***
			key := r.String()
			r.WantColon()
			ret[key] = r.Interface()
			r.WantComma()
		***REMOVED***
		r.Delim('***REMOVED***')

		if r.Ok() ***REMOVED***
			return ret
		***REMOVED*** else ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED*** else if r.token.delimValue == '[' ***REMOVED***
		r.consume()

		ret := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
		for !r.IsDelim(']') ***REMOVED***
			ret = append(ret, r.Interface())
			r.WantComma()
		***REMOVED***
		r.Delim(']')

		if r.Ok() ***REMOVED***
			return ret
		***REMOVED*** else ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	r.errSyntax()
	return nil
***REMOVED***

// WantComma requires a comma to be present before fetching next token.
func (r *Lexer) WantComma() ***REMOVED***
	r.wantSep = ','
	r.firstElement = false
***REMOVED***

// WantColon requires a colon to be present before fetching next token.
func (r *Lexer) WantColon() ***REMOVED***
	r.wantSep = ':'
	r.firstElement = false
***REMOVED***
