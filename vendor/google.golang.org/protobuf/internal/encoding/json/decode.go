// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/errors"
)

// call specifies which Decoder method was invoked.
type call uint8

const (
	readCall call = iota
	peekCall
)

const unexpectedFmt = "unexpected token %s"

// ErrUnexpectedEOF means that EOF was encountered in the middle of the input.
var ErrUnexpectedEOF = errors.New("%v", io.ErrUnexpectedEOF)

// Decoder is a token-based JSON decoder.
type Decoder struct ***REMOVED***
	// lastCall is last method called, either readCall or peekCall.
	// Initial value is readCall.
	lastCall call

	// lastToken contains the last read token.
	lastToken Token

	// lastErr contains the last read error.
	lastErr error

	// openStack is a stack containing ObjectOpen and ArrayOpen values. The
	// top of stack represents the object or the array the current value is
	// directly located in.
	openStack []Kind

	// orig is used in reporting line and column.
	orig []byte
	// in contains the unconsumed input.
	in []byte
***REMOVED***

// NewDecoder returns a Decoder to read the given []byte.
func NewDecoder(b []byte) *Decoder ***REMOVED***
	return &Decoder***REMOVED***orig: b, in: b***REMOVED***
***REMOVED***

// Peek looks ahead and returns the next token kind without advancing a read.
func (d *Decoder) Peek() (Token, error) ***REMOVED***
	defer func() ***REMOVED*** d.lastCall = peekCall ***REMOVED***()
	if d.lastCall == readCall ***REMOVED***
		d.lastToken, d.lastErr = d.Read()
	***REMOVED***
	return d.lastToken, d.lastErr
***REMOVED***

// Read returns the next JSON token.
// It will return an error if there is no valid token.
func (d *Decoder) Read() (Token, error) ***REMOVED***
	const scalar = Null | Bool | Number | String

	defer func() ***REMOVED*** d.lastCall = readCall ***REMOVED***()
	if d.lastCall == peekCall ***REMOVED***
		return d.lastToken, d.lastErr
	***REMOVED***

	tok, err := d.parseNext()
	if err != nil ***REMOVED***
		return Token***REMOVED******REMOVED***, err
	***REMOVED***

	switch tok.kind ***REMOVED***
	case EOF:
		if len(d.openStack) != 0 ||
			d.lastToken.kind&scalar|ObjectClose|ArrayClose == 0 ***REMOVED***
			return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
		***REMOVED***

	case Null:
		if !d.isValueNext() ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(tok.pos, unexpectedFmt, tok.RawString())
		***REMOVED***

	case Bool, Number:
		if !d.isValueNext() ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(tok.pos, unexpectedFmt, tok.RawString())
		***REMOVED***

	case String:
		if d.isValueNext() ***REMOVED***
			break
		***REMOVED***
		// This string token should only be for a field name.
		if d.lastToken.kind&(ObjectOpen|comma) == 0 ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(tok.pos, unexpectedFmt, tok.RawString())
		***REMOVED***
		if len(d.in) == 0 ***REMOVED***
			return Token***REMOVED******REMOVED***, ErrUnexpectedEOF
		***REMOVED***
		if c := d.in[0]; c != ':' ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(d.currPos(), `unexpected character %s, missing ":" after field name`, string(c))
		***REMOVED***
		tok.kind = Name
		d.consume(1)

	case ObjectOpen, ArrayOpen:
		if !d.isValueNext() ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(tok.pos, unexpectedFmt, tok.RawString())
		***REMOVED***
		d.openStack = append(d.openStack, tok.kind)

	case ObjectClose:
		if len(d.openStack) == 0 ||
			d.lastToken.kind == comma ||
			d.openStack[len(d.openStack)-1] != ObjectOpen ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(tok.pos, unexpectedFmt, tok.RawString())
		***REMOVED***
		d.openStack = d.openStack[:len(d.openStack)-1]

	case ArrayClose:
		if len(d.openStack) == 0 ||
			d.lastToken.kind == comma ||
			d.openStack[len(d.openStack)-1] != ArrayOpen ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(tok.pos, unexpectedFmt, tok.RawString())
		***REMOVED***
		d.openStack = d.openStack[:len(d.openStack)-1]

	case comma:
		if len(d.openStack) == 0 ||
			d.lastToken.kind&(scalar|ObjectClose|ArrayClose) == 0 ***REMOVED***
			return Token***REMOVED******REMOVED***, d.newSyntaxError(tok.pos, unexpectedFmt, tok.RawString())
		***REMOVED***
	***REMOVED***

	// Update d.lastToken only after validating token to be in the right sequence.
	d.lastToken = tok

	if d.lastToken.kind == comma ***REMOVED***
		return d.Read()
	***REMOVED***
	return tok, nil
***REMOVED***

// Any sequence that looks like a non-delimiter (for error reporting).
var errRegexp = regexp.MustCompile(`^([-+._a-zA-Z0-9]***REMOVED***1,32***REMOVED***|.)`)

// parseNext parses for the next JSON token. It returns a Token object for
// different types, except for Name. It does not handle whether the next token
// is in a valid sequence or not.
func (d *Decoder) parseNext() (Token, error) ***REMOVED***
	// Trim leading spaces.
	d.consume(0)

	in := d.in
	if len(in) == 0 ***REMOVED***
		return d.consumeToken(EOF, 0), nil
	***REMOVED***

	switch in[0] ***REMOVED***
	case 'n':
		if n := matchWithDelim("null", in); n != 0 ***REMOVED***
			return d.consumeToken(Null, n), nil
		***REMOVED***

	case 't':
		if n := matchWithDelim("true", in); n != 0 ***REMOVED***
			return d.consumeBoolToken(true, n), nil
		***REMOVED***

	case 'f':
		if n := matchWithDelim("false", in); n != 0 ***REMOVED***
			return d.consumeBoolToken(false, n), nil
		***REMOVED***

	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if n, ok := parseNumber(in); ok ***REMOVED***
			return d.consumeToken(Number, n), nil
		***REMOVED***

	case '"':
		s, n, err := d.parseString(in)
		if err != nil ***REMOVED***
			return Token***REMOVED******REMOVED***, err
		***REMOVED***
		return d.consumeStringToken(s, n), nil

	case '***REMOVED***':
		return d.consumeToken(ObjectOpen, 1), nil

	case '***REMOVED***':
		return d.consumeToken(ObjectClose, 1), nil

	case '[':
		return d.consumeToken(ArrayOpen, 1), nil

	case ']':
		return d.consumeToken(ArrayClose, 1), nil

	case ',':
		return d.consumeToken(comma, 1), nil
	***REMOVED***
	return Token***REMOVED******REMOVED***, d.newSyntaxError(d.currPos(), "invalid value %s", errRegexp.Find(in))
***REMOVED***

// newSyntaxError returns an error with line and column information useful for
// syntax errors.
func (d *Decoder) newSyntaxError(pos int, f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	e := errors.New(f, x...)
	line, column := d.Position(pos)
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

// currPos returns the current index position of d.in from d.orig.
func (d *Decoder) currPos() int ***REMOVED***
	return len(d.orig) - len(d.in)
***REMOVED***

// matchWithDelim matches s with the input b and verifies that the match
// terminates with a delimiter of some form (e.g., r"[^-+_.a-zA-Z0-9]").
// As a special case, EOF is considered a delimiter. It returns the length of s
// if there is a match, else 0.
func matchWithDelim(s string, b []byte) int ***REMOVED***
	if !bytes.HasPrefix(b, []byte(s)) ***REMOVED***
		return 0
	***REMOVED***

	n := len(s)
	if n < len(b) && isNotDelim(b[n]) ***REMOVED***
		return 0
	***REMOVED***
	return n
***REMOVED***

// isNotDelim returns true if given byte is a not delimiter character.
func isNotDelim(c byte) bool ***REMOVED***
	return (c == '-' || c == '+' || c == '.' || c == '_' ||
		('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9'))
***REMOVED***

// consume consumes n bytes of input and any subsequent whitespace.
func (d *Decoder) consume(n int) ***REMOVED***
	d.in = d.in[n:]
	for len(d.in) > 0 ***REMOVED***
		switch d.in[0] ***REMOVED***
		case ' ', '\n', '\r', '\t':
			d.in = d.in[1:]
		default:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// isValueNext returns true if next type should be a JSON value: Null,
// Number, String or Bool.
func (d *Decoder) isValueNext() bool ***REMOVED***
	if len(d.openStack) == 0 ***REMOVED***
		return d.lastToken.kind == 0
	***REMOVED***

	start := d.openStack[len(d.openStack)-1]
	switch start ***REMOVED***
	case ObjectOpen:
		return d.lastToken.kind&Name != 0
	case ArrayOpen:
		return d.lastToken.kind&(ArrayOpen|comma) != 0
	***REMOVED***
	panic(fmt.Sprintf(
		"unreachable logic in Decoder.isValueNext, lastToken.kind: %v, openStack: %v",
		d.lastToken.kind, start))
***REMOVED***

// consumeToken constructs a Token for given Kind with raw value derived from
// current d.in and given size, and consumes the given size-lenght of it.
func (d *Decoder) consumeToken(kind Kind, size int) Token ***REMOVED***
	tok := Token***REMOVED***
		kind: kind,
		raw:  d.in[:size],
		pos:  len(d.orig) - len(d.in),
	***REMOVED***
	d.consume(size)
	return tok
***REMOVED***

// consumeBoolToken constructs a Token for a Bool kind with raw value derived from
// current d.in and given size.
func (d *Decoder) consumeBoolToken(b bool, size int) Token ***REMOVED***
	tok := Token***REMOVED***
		kind: Bool,
		raw:  d.in[:size],
		pos:  len(d.orig) - len(d.in),
		boo:  b,
	***REMOVED***
	d.consume(size)
	return tok
***REMOVED***

// consumeStringToken constructs a Token for a String kind with raw value derived
// from current d.in and given size.
func (d *Decoder) consumeStringToken(s string, size int) Token ***REMOVED***
	tok := Token***REMOVED***
		kind: String,
		raw:  d.in[:size],
		pos:  len(d.orig) - len(d.in),
		str:  s,
	***REMOVED***
	d.consume(size)
	return tok
***REMOVED***

// Clone returns a copy of the Decoder for use in reading ahead the next JSON
// object, array or other values without affecting current Decoder.
func (d *Decoder) Clone() *Decoder ***REMOVED***
	ret := *d
	ret.openStack = append([]Kind(nil), ret.openStack...)
	return &ret
***REMOVED***
