// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/errors"
)

// Decoder is a token-based textproto decoder.
type Decoder struct ***REMOVED***
	// lastCall is last method called, either readCall or peekCall.
	// Initial value is readCall.
	lastCall call

	// lastToken contains the last read token.
	lastToken Token

	// lastErr contains the last read error.
	lastErr error

	// openStack is a stack containing the byte characters for MessageOpen and
	// ListOpen kinds. The top of stack represents the message or the list that
	// the current token is nested in. An empty stack means the current token is
	// at the top level message. The characters '***REMOVED***' and '<' both represent the
	// MessageOpen kind.
	openStack []byte

	// orig is used in reporting line and column.
	orig []byte
	// in contains the unconsumed input.
	in []byte
***REMOVED***

// NewDecoder returns a Decoder to read the given []byte.
func NewDecoder(b []byte) *Decoder ***REMOVED***
	return &Decoder***REMOVED***orig: b, in: b***REMOVED***
***REMOVED***

// ErrUnexpectedEOF means that EOF was encountered in the middle of the input.
var ErrUnexpectedEOF = errors.New("%v", io.ErrUnexpectedEOF)

// call specifies which Decoder method was invoked.
type call uint8

const (
	readCall call = iota
	peekCall
)

// Peek looks ahead and returns the next token and error without advancing a read.
func (d *Decoder) Peek() (Token, error) ***REMOVED***
	defer func() ***REMOVED*** d.lastCall = peekCall ***REMOVED***()
	if d.lastCall == readCall ***REMOVED***
		d.lastToken, d.lastErr = d.Read()
	***REMOVED***
	return d.lastToken, d.lastErr
***REMOVED***

// Read returns the next token.
// It will return an error if there is no valid token.
func (d *Decoder) Read() (Token, error) ***REMOVED***
	defer func() ***REMOVED*** d.lastCall = readCall ***REMOVED***()
	if d.lastCall == peekCall ***REMOVED***
		return d.lastToken, d.lastErr
	***REMOVED***

	tok, err := d.parseNext(d.lastToken.kind)
	if err != nil ***REMOVED***
		return Token***REMOVED******REMOVED***, err
	***REMOVED***

	switch tok.kind ***REMOVED***
	case comma, semicolon:
		tok, err = d.parseNext(tok.kind)
		if err != nil ***REMOVED***
			return Token***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***
	d.lastToken = tok
	return tok, nil
***REMOVED***

const (
	mismatchedFmt = "mismatched close character %q"
	unexpectedFmt = "unexpected character %q"
)

// parseNext parses the next Token based on given last kind.
func (d *Decoder) parseNext(lastKind Kind) (Token, error) ***REMOVED***
	// Trim leading spaces.
	d.consume(0)
	isEOF := false
	if len(d.in) == 0 ***REMOVED***
		isEOF = true
	***REMOVED***

	switch lastKind ***REMOVED***
	case EOF:
		return d.consumeToken(EOF, 0, 0), nil

	case bof:
		// Start of top level message. Next token can be EOF or Name.
		if isEOF ***REMOVED***
			return d.consumeToken(EOF, 0, 0), nil
		***REMOVED***
		return d.parseFieldName()

	case Name:
		// Next token can be MessageOpen, ListOpen or Scalar.
		if isEOF ***REMOVED***
			return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
		***REMOVED***
		switch ch := d.in[0]; ch ***REMOVED***
		case '***REMOVED***', '<':
			d.pushOpenStack(ch)
			return d.consumeToken(MessageOpen, 1, 0), nil
		case '[':
			d.pushOpenStack(ch)
			return d.consumeToken(ListOpen, 1, 0), nil
		default:
			return d.parseScalar()
		***REMOVED***

	case Scalar:
		openKind, closeCh := d.currentOpenKind()
		switch openKind ***REMOVED***
		case bof:
			// Top level message.
			// 	Next token can be EOF, comma, semicolon or Name.
			if isEOF ***REMOVED***
				return d.consumeToken(EOF, 0, 0), nil
			***REMOVED***
			switch d.in[0] ***REMOVED***
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			***REMOVED***

		case MessageOpen:
			// Next token can be MessageClose, comma, semicolon or Name.
			if isEOF ***REMOVED***
				return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token***REMOVED******REMOVED***, d.newSyntaxError(mismatchedFmt, ch)
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			***REMOVED***

		case ListOpen:
			// Next token can be ListClose or comma.
			if isEOF ***REMOVED***
				return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case ']':
				d.popOpenStack()
				return d.consumeToken(ListClose, 1, 0), nil
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			default:
				return Token***REMOVED******REMOVED***, d.newSyntaxError(unexpectedFmt, ch)
			***REMOVED***
		***REMOVED***

	case MessageOpen:
		// Next token can be MessageClose or Name.
		if isEOF ***REMOVED***
			return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
		***REMOVED***
		_, closeCh := d.currentOpenKind()
		switch ch := d.in[0]; ch ***REMOVED***
		case closeCh:
			d.popOpenStack()
			return d.consumeToken(MessageClose, 1, 0), nil
		case otherCloseChar[closeCh]:
			return Token***REMOVED******REMOVED***, d.newSyntaxError(mismatchedFmt, ch)
		default:
			return d.parseFieldName()
		***REMOVED***

	case MessageClose:
		openKind, closeCh := d.currentOpenKind()
		switch openKind ***REMOVED***
		case bof:
			// Top level message.
			// Next token can be EOF, comma, semicolon or Name.
			if isEOF ***REMOVED***
				return d.consumeToken(EOF, 0, 0), nil
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			***REMOVED***

		case MessageOpen:
			// Next token can be MessageClose, comma, semicolon or Name.
			if isEOF ***REMOVED***
				return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token***REMOVED******REMOVED***, d.newSyntaxError(mismatchedFmt, ch)
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			***REMOVED***

		case ListOpen:
			// Next token can be ListClose or comma
			if isEOF ***REMOVED***
				return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(ListClose, 1, 0), nil
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			default:
				return Token***REMOVED******REMOVED***, d.newSyntaxError(unexpectedFmt, ch)
			***REMOVED***
		***REMOVED***

	case ListOpen:
		// Next token can be ListClose, MessageStart or Scalar.
		if isEOF ***REMOVED***
			return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
		***REMOVED***
		switch ch := d.in[0]; ch ***REMOVED***
		case ']':
			d.popOpenStack()
			return d.consumeToken(ListClose, 1, 0), nil
		case '***REMOVED***', '<':
			d.pushOpenStack(ch)
			return d.consumeToken(MessageOpen, 1, 0), nil
		default:
			return d.parseScalar()
		***REMOVED***

	case ListClose:
		openKind, closeCh := d.currentOpenKind()
		switch openKind ***REMOVED***
		case bof:
			// Top level message.
			// Next token can be EOF, comma, semicolon or Name.
			if isEOF ***REMOVED***
				return d.consumeToken(EOF, 0, 0), nil
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			***REMOVED***

		case MessageOpen:
			// Next token can be MessageClose, comma, semicolon or Name.
			if isEOF ***REMOVED***
				return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token***REMOVED******REMOVED***, d.newSyntaxError(mismatchedFmt, ch)
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			***REMOVED***

		default:
			// It is not possible to have this case. Let it panic below.
		***REMOVED***

	case comma, semicolon:
		openKind, closeCh := d.currentOpenKind()
		switch openKind ***REMOVED***
		case bof:
			// Top level message. Next token can be EOF or Name.
			if isEOF ***REMOVED***
				return d.consumeToken(EOF, 0, 0), nil
			***REMOVED***
			return d.parseFieldName()

		case MessageOpen:
			// Next token can be MessageClose or Name.
			if isEOF ***REMOVED***
				return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token***REMOVED******REMOVED***, d.newSyntaxError(mismatchedFmt, ch)
			default:
				return d.parseFieldName()
			***REMOVED***

		case ListOpen:
			if lastKind == semicolon ***REMOVED***
				// It is not be possible to have this case as logic here
				// should not have produced a semicolon Token when inside a
				// list. Let it panic below.
				break
			***REMOVED***
			// Next token can be MessageOpen or Scalar.
			if isEOF ***REMOVED***
				return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
			***REMOVED***
			switch ch := d.in[0]; ch ***REMOVED***
			case '***REMOVED***', '<':
				d.pushOpenStack(ch)
				return d.consumeToken(MessageOpen, 1, 0), nil
			default:
				return d.parseScalar()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	line, column := d.Position(len(d.orig) - len(d.in))
	panic(fmt.Sprintf("Decoder.parseNext: bug at handling line %d:%d with lastKind=%v", line, column, lastKind))
***REMOVED***

var otherCloseChar = map[byte]byte***REMOVED***
	'***REMOVED***': '>',
	'>': '***REMOVED***',
***REMOVED***

// currentOpenKind indicates whether current position is inside a message, list
// or top-level message by returning MessageOpen, ListOpen or bof respectively.
// If the returned kind is either a MessageOpen or ListOpen, it also returns the
// corresponding closing character.
func (d *Decoder) currentOpenKind() (Kind, byte) ***REMOVED***
	if len(d.openStack) == 0 ***REMOVED***
		return bof, 0
	***REMOVED***
	openCh := d.openStack[len(d.openStack)-1]
	switch openCh ***REMOVED***
	case '***REMOVED***':
		return MessageOpen, '***REMOVED***'
	case '<':
		return MessageOpen, '>'
	case '[':
		return ListOpen, ']'
	***REMOVED***
	panic(fmt.Sprintf("Decoder: openStack contains invalid byte %c", openCh))
***REMOVED***

func (d *Decoder) pushOpenStack(ch byte) ***REMOVED***
	d.openStack = append(d.openStack, ch)
***REMOVED***

func (d *Decoder) popOpenStack() ***REMOVED***
	d.openStack = d.openStack[:len(d.openStack)-1]
***REMOVED***

// parseFieldName parses field name and separator.
func (d *Decoder) parseFieldName() (tok Token, err error) ***REMOVED***
	defer func() ***REMOVED***
		if err == nil && d.tryConsumeChar(':') ***REMOVED***
			tok.attrs |= hasSeparator
		***REMOVED***
	***REMOVED***()

	// Extension or Any type URL.
	if d.in[0] == '[' ***REMOVED***
		return d.parseTypeName()
	***REMOVED***

	// Identifier.
	if size := parseIdent(d.in, false); size > 0 ***REMOVED***
		return d.consumeToken(Name, size, uint8(IdentName)), nil
	***REMOVED***

	// Field number. Identify if input is a valid number that is not negative
	// and is decimal integer within 32-bit range.
	if num := parseNumber(d.in); num.size > 0 ***REMOVED***
		if !num.neg && num.kind == numDec ***REMOVED***
			if _, err := strconv.ParseInt(string(d.in[:num.size]), 10, 32); err == nil ***REMOVED***
				return d.consumeToken(Name, num.size, uint8(FieldNumber)), nil
			***REMOVED***
		***REMOVED***
		return Token***REMOVED******REMOVED***, d.newSyntaxError("invalid field number: %s", d.in[:num.size])
	***REMOVED***

	return Token***REMOVED******REMOVED***, d.newSyntaxError("invalid field name: %s", errId(d.in))
***REMOVED***

// parseTypeName parses Any type URL or extension field name. The name is
// enclosed in [ and ] characters. The C++ parser does not handle many legal URL
// strings. This implementation is more liberal and allows for the pattern
// ^[-_a-zA-Z0-9]+([./][-_a-zA-Z0-9]+)*`). Whitespaces and comments are allowed
// in between [ ], '.', '/' and the sub names.
func (d *Decoder) parseTypeName() (Token, error) ***REMOVED***
	startPos := len(d.orig) - len(d.in)
	// Use alias s to advance first in order to use d.in for error handling.
	// Caller already checks for [ as first character.
	s := consume(d.in[1:], 0)
	if len(s) == 0 ***REMOVED***
		return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
	***REMOVED***

	var name []byte
	for len(s) > 0 && isTypeNameChar(s[0]) ***REMOVED***
		name = append(name, s[0])
		s = s[1:]
	***REMOVED***
	s = consume(s, 0)

	var closed bool
	for len(s) > 0 && !closed ***REMOVED***
		switch ***REMOVED***
		case s[0] == ']':
			s = s[1:]
			closed = true

		case s[0] == '/', s[0] == '.':
			if len(name) > 0 && (name[len(name)-1] == '/' || name[len(name)-1] == '.') ***REMOVED***
				return Token***REMOVED******REMOVED***, d.newSyntaxError("invalid type URL/extension field name: %s",
					d.orig[startPos:len(d.orig)-len(s)+1])
			***REMOVED***
			name = append(name, s[0])
			s = s[1:]
			s = consume(s, 0)
			for len(s) > 0 && isTypeNameChar(s[0]) ***REMOVED***
				name = append(name, s[0])
				s = s[1:]
			***REMOVED***
			s = consume(s, 0)

		default:
			return Token***REMOVED******REMOVED***, d.newSyntaxError(
				"invalid type URL/extension field name: %s", d.orig[startPos:len(d.orig)-len(s)+1])
		***REMOVED***
	***REMOVED***

	if !closed ***REMOVED***
		return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
	***REMOVED***

	// First character cannot be '.'. Last character cannot be '.' or '/'.
	size := len(name)
	if size == 0 || name[0] == '.' || name[size-1] == '.' || name[size-1] == '/' ***REMOVED***
		return Token***REMOVED******REMOVED***, d.newSyntaxError("invalid type URL/extension field name: %s",
			d.orig[startPos:len(d.orig)-len(s)])
	***REMOVED***

	d.in = s
	endPos := len(d.orig) - len(d.in)
	d.consume(0)

	return Token***REMOVED***
		kind:  Name,
		attrs: uint8(TypeName),
		pos:   startPos,
		raw:   d.orig[startPos:endPos],
		str:   string(name),
	***REMOVED***, nil
***REMOVED***

func isTypeNameChar(b byte) bool ***REMOVED***
	return (b == '-' || b == '_' ||
		('0' <= b && b <= '9') ||
		('a' <= b && b <= 'z') ||
		('A' <= b && b <= 'Z'))
***REMOVED***

func isWhiteSpace(b byte) bool ***REMOVED***
	switch b ***REMOVED***
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// parseIdent parses an unquoted proto identifier and returns size.
// If allowNeg is true, it allows '-' to be the first character in the
// identifier. This is used when parsing literal values like -infinity, etc.
// Regular expression matches an identifier: `^[_a-zA-Z][_a-zA-Z0-9]*`
func parseIdent(input []byte, allowNeg bool) int ***REMOVED***
	var size int

	s := input
	if len(s) == 0 ***REMOVED***
		return 0
	***REMOVED***

	if allowNeg && s[0] == '-' ***REMOVED***
		s = s[1:]
		size++
		if len(s) == 0 ***REMOVED***
			return 0
		***REMOVED***
	***REMOVED***

	switch ***REMOVED***
	case s[0] == '_',
		'a' <= s[0] && s[0] <= 'z',
		'A' <= s[0] && s[0] <= 'Z':
		s = s[1:]
		size++
	default:
		return 0
	***REMOVED***

	for len(s) > 0 && (s[0] == '_' ||
		'a' <= s[0] && s[0] <= 'z' ||
		'A' <= s[0] && s[0] <= 'Z' ||
		'0' <= s[0] && s[0] <= '9') ***REMOVED***
		s = s[1:]
		size++
	***REMOVED***

	if len(s) > 0 && !isDelim(s[0]) ***REMOVED***
		return 0
	***REMOVED***

	return size
***REMOVED***

// parseScalar parses for a string, literal or number value.
func (d *Decoder) parseScalar() (Token, error) ***REMOVED***
	if d.in[0] == '"' || d.in[0] == '\'' ***REMOVED***
		return d.parseStringValue()
	***REMOVED***

	if tok, ok := d.parseLiteralValue(); ok ***REMOVED***
		return tok, nil
	***REMOVED***

	if tok, ok := d.parseNumberValue(); ok ***REMOVED***
		return tok, nil
	***REMOVED***

	return Token***REMOVED******REMOVED***, d.newSyntaxError("invalid scalar value: %s", errId(d.in))
***REMOVED***

// parseLiteralValue parses a literal value. A literal value is used for
// bools, special floats and enums. This function simply identifies that the
// field value is a literal.
func (d *Decoder) parseLiteralValue() (Token, bool) ***REMOVED***
	size := parseIdent(d.in, true)
	if size == 0 ***REMOVED***
		return Token***REMOVED******REMOVED***, false
	***REMOVED***
	return d.consumeToken(Scalar, size, literalValue), true
***REMOVED***

// consumeToken constructs a Token for given Kind from d.in and consumes given
// size-length from it.
func (d *Decoder) consumeToken(kind Kind, size int, attrs uint8) Token ***REMOVED***
	// Important to compute raw and pos before consuming.
	tok := Token***REMOVED***
		kind:  kind,
		attrs: attrs,
		pos:   len(d.orig) - len(d.in),
		raw:   d.in[:size],
	***REMOVED***
	d.consume(size)
	return tok
***REMOVED***

// newSyntaxError returns a syntax error with line and column information for
// current position.
func (d *Decoder) newSyntaxError(f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	e := errors.New(f, x...)
	line, column := d.Position(len(d.orig) - len(d.in))
	return errors.New("syntax error (line %d:%d): %v", line, column, e)
***REMOVED***

// Position returns line and column number of given index of the original input.
// It will panic if index is out of range.
func (d *Decoder) Position(idx int) (line int, column int) ***REMOVED***
	b := d.orig[:idx]
	line = bytes.Count(b, []byte("\n")) + 1
	if i := bytes.LastIndexByte(b, '\n'); i >= 0 ***REMOVED***
		b = b[i+1:]
	***REMOVED***
	column = utf8.RuneCount(b) + 1 // ignore multi-rune characters
	return line, column
***REMOVED***

func (d *Decoder) tryConsumeChar(c byte) bool ***REMOVED***
	if len(d.in) > 0 && d.in[0] == c ***REMOVED***
		d.consume(1)
		return true
	***REMOVED***
	return false
***REMOVED***

// consume consumes n bytes of input and any subsequent whitespace or comments.
func (d *Decoder) consume(n int) ***REMOVED***
	d.in = consume(d.in, n)
	return
***REMOVED***

// consume consumes n bytes of input and any subsequent whitespace or comments.
func consume(b []byte, n int) []byte ***REMOVED***
	b = b[n:]
	for len(b) > 0 ***REMOVED***
		switch b[0] ***REMOVED***
		case ' ', '\n', '\r', '\t':
			b = b[1:]
		case '#':
			if i := bytes.IndexByte(b, '\n'); i >= 0 ***REMOVED***
				b = b[i+len("\n"):]
			***REMOVED*** else ***REMOVED***
				b = nil
			***REMOVED***
		default:
			return b
		***REMOVED***
	***REMOVED***
	return b
***REMOVED***

// errId extracts a byte sequence that looks like an invalid ID
// (for the purposes of error reporting).
func errId(seq []byte) []byte ***REMOVED***
	const maxLen = 32
	for i := 0; i < len(seq); ***REMOVED***
		if i > maxLen ***REMOVED***
			return append(seq[:i:i], "…"...)
		***REMOVED***
		r, size := utf8.DecodeRune(seq[i:])
		if r > utf8.RuneSelf || (r != '/' && isDelim(byte(r))) ***REMOVED***
			if i == 0 ***REMOVED***
				// Either the first byte is invalid UTF-8 or a
				// delimiter, or the first rune is non-ASCII.
				// Return it as-is.
				i = size
			***REMOVED***
			return seq[:i:i]
		***REMOVED***
		i += size
	***REMOVED***
	// No delimiter found.
	return seq
***REMOVED***

// isDelim returns true if given byte is a delimiter character.
func isDelim(c byte) bool ***REMOVED***
	return !(c == '-' || c == '+' || c == '.' || c == '_' ||
		('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9'))
***REMOVED***
