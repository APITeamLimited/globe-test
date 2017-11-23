// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package proto

// Functions for parsing the Text protocol buffer format.
// TODO: message sets.

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Error string emitted when deserializing Any and fields are already set
const anyRepeatedlyUnpacked = "Any message unpacked multiple times, or %q already set"

type ParseError struct ***REMOVED***
	Message string
	Line    int // 1-based line number
	Offset  int // 0-based byte offset from start of input
***REMOVED***

func (p *ParseError) Error() string ***REMOVED***
	if p.Line == 1 ***REMOVED***
		// show offset only for first line
		return fmt.Sprintf("line 1.%d: %v", p.Offset, p.Message)
	***REMOVED***
	return fmt.Sprintf("line %d: %v", p.Line, p.Message)
***REMOVED***

type token struct ***REMOVED***
	value    string
	err      *ParseError
	line     int    // line number
	offset   int    // byte number from start of input, not start of line
	unquoted string // the unquoted version of value, if it was a quoted string
***REMOVED***

func (t *token) String() string ***REMOVED***
	if t.err == nil ***REMOVED***
		return fmt.Sprintf("%q (line=%d, offset=%d)", t.value, t.line, t.offset)
	***REMOVED***
	return fmt.Sprintf("parse error: %v", t.err)
***REMOVED***

type textParser struct ***REMOVED***
	s            string // remaining input
	done         bool   // whether the parsing is finished (success or error)
	backed       bool   // whether back() was called
	offset, line int
	cur          token
***REMOVED***

func newTextParser(s string) *textParser ***REMOVED***
	p := new(textParser)
	p.s = s
	p.line = 1
	p.cur.line = 1
	return p
***REMOVED***

func (p *textParser) errorf(format string, a ...interface***REMOVED******REMOVED***) *ParseError ***REMOVED***
	pe := &ParseError***REMOVED***fmt.Sprintf(format, a...), p.cur.line, p.cur.offset***REMOVED***
	p.cur.err = pe
	p.done = true
	return pe
***REMOVED***

// Numbers and identifiers are matched by [-+._A-Za-z0-9]
func isIdentOrNumberChar(c byte) bool ***REMOVED***
	switch ***REMOVED***
	case 'A' <= c && c <= 'Z', 'a' <= c && c <= 'z':
		return true
	case '0' <= c && c <= '9':
		return true
	***REMOVED***
	switch c ***REMOVED***
	case '-', '+', '.', '_':
		return true
	***REMOVED***
	return false
***REMOVED***

func isWhitespace(c byte) bool ***REMOVED***
	switch c ***REMOVED***
	case ' ', '\t', '\n', '\r':
		return true
	***REMOVED***
	return false
***REMOVED***

func isQuote(c byte) bool ***REMOVED***
	switch c ***REMOVED***
	case '"', '\'':
		return true
	***REMOVED***
	return false
***REMOVED***

func (p *textParser) skipWhitespace() ***REMOVED***
	i := 0
	for i < len(p.s) && (isWhitespace(p.s[i]) || p.s[i] == '#') ***REMOVED***
		if p.s[i] == '#' ***REMOVED***
			// comment; skip to end of line or input
			for i < len(p.s) && p.s[i] != '\n' ***REMOVED***
				i++
			***REMOVED***
			if i == len(p.s) ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		if p.s[i] == '\n' ***REMOVED***
			p.line++
		***REMOVED***
		i++
	***REMOVED***
	p.offset += i
	p.s = p.s[i:len(p.s)]
	if len(p.s) == 0 ***REMOVED***
		p.done = true
	***REMOVED***
***REMOVED***

func (p *textParser) advance() ***REMOVED***
	// Skip whitespace
	p.skipWhitespace()
	if p.done ***REMOVED***
		return
	***REMOVED***

	// Start of non-whitespace
	p.cur.err = nil
	p.cur.offset, p.cur.line = p.offset, p.line
	p.cur.unquoted = ""
	switch p.s[0] ***REMOVED***
	case '<', '>', '***REMOVED***', '***REMOVED***', ':', '[', ']', ';', ',', '/':
		// Single symbol
		p.cur.value, p.s = p.s[0:1], p.s[1:len(p.s)]
	case '"', '\'':
		// Quoted string
		i := 1
		for i < len(p.s) && p.s[i] != p.s[0] && p.s[i] != '\n' ***REMOVED***
			if p.s[i] == '\\' && i+1 < len(p.s) ***REMOVED***
				// skip escaped char
				i++
			***REMOVED***
			i++
		***REMOVED***
		if i >= len(p.s) || p.s[i] != p.s[0] ***REMOVED***
			p.errorf("unmatched quote")
			return
		***REMOVED***
		unq, err := unquoteC(p.s[1:i], rune(p.s[0]))
		if err != nil ***REMOVED***
			p.errorf("invalid quoted string %s: %v", p.s[0:i+1], err)
			return
		***REMOVED***
		p.cur.value, p.s = p.s[0:i+1], p.s[i+1:len(p.s)]
		p.cur.unquoted = unq
	default:
		i := 0
		for i < len(p.s) && isIdentOrNumberChar(p.s[i]) ***REMOVED***
			i++
		***REMOVED***
		if i == 0 ***REMOVED***
			p.errorf("unexpected byte %#x", p.s[0])
			return
		***REMOVED***
		p.cur.value, p.s = p.s[0:i], p.s[i:len(p.s)]
	***REMOVED***
	p.offset += len(p.cur.value)
***REMOVED***

var (
	errBadUTF8 = errors.New("proto: bad UTF-8")
	errBadHex  = errors.New("proto: bad hexadecimal")
)

func unquoteC(s string, quote rune) (string, error) ***REMOVED***
	// This is based on C++'s tokenizer.cc.
	// Despite its name, this is *not* parsing C syntax.
	// For instance, "\0" is an invalid quoted string.

	// Avoid allocation in trivial cases.
	simple := true
	for _, r := range s ***REMOVED***
		if r == '\\' || r == quote ***REMOVED***
			simple = false
			break
		***REMOVED***
	***REMOVED***
	if simple ***REMOVED***
		return s, nil
	***REMOVED***

	buf := make([]byte, 0, 3*len(s)/2)
	for len(s) > 0 ***REMOVED***
		r, n := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && n == 1 ***REMOVED***
			return "", errBadUTF8
		***REMOVED***
		s = s[n:]
		if r != '\\' ***REMOVED***
			if r < utf8.RuneSelf ***REMOVED***
				buf = append(buf, byte(r))
			***REMOVED*** else ***REMOVED***
				buf = append(buf, string(r)...)
			***REMOVED***
			continue
		***REMOVED***

		ch, tail, err := unescape(s)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		buf = append(buf, ch...)
		s = tail
	***REMOVED***
	return string(buf), nil
***REMOVED***

func unescape(s string) (ch string, tail string, err error) ***REMOVED***
	r, n := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && n == 1 ***REMOVED***
		return "", "", errBadUTF8
	***REMOVED***
	s = s[n:]
	switch r ***REMOVED***
	case 'a':
		return "\a", s, nil
	case 'b':
		return "\b", s, nil
	case 'f':
		return "\f", s, nil
	case 'n':
		return "\n", s, nil
	case 'r':
		return "\r", s, nil
	case 't':
		return "\t", s, nil
	case 'v':
		return "\v", s, nil
	case '?':
		return "?", s, nil // trigraph workaround
	case '\'', '"', '\\':
		return string(r), s, nil
	case '0', '1', '2', '3', '4', '5', '6', '7', 'x', 'X':
		if len(s) < 2 ***REMOVED***
			return "", "", fmt.Errorf(`\%c requires 2 following digits`, r)
		***REMOVED***
		base := 8
		ss := s[:2]
		s = s[2:]
		if r == 'x' || r == 'X' ***REMOVED***
			base = 16
		***REMOVED*** else ***REMOVED***
			ss = string(r) + ss
		***REMOVED***
		i, err := strconv.ParseUint(ss, base, 8)
		if err != nil ***REMOVED***
			return "", "", err
		***REMOVED***
		return string([]byte***REMOVED***byte(i)***REMOVED***), s, nil
	case 'u', 'U':
		n := 4
		if r == 'U' ***REMOVED***
			n = 8
		***REMOVED***
		if len(s) < n ***REMOVED***
			return "", "", fmt.Errorf(`\%c requires %d digits`, r, n)
		***REMOVED***

		bs := make([]byte, n/2)
		for i := 0; i < n; i += 2 ***REMOVED***
			a, ok1 := unhex(s[i])
			b, ok2 := unhex(s[i+1])
			if !ok1 || !ok2 ***REMOVED***
				return "", "", errBadHex
			***REMOVED***
			bs[i/2] = a<<4 | b
		***REMOVED***
		s = s[n:]
		return string(bs), s, nil
	***REMOVED***
	return "", "", fmt.Errorf(`unknown escape \%c`, r)
***REMOVED***

// Adapted from src/pkg/strconv/quote.go.
func unhex(b byte) (v byte, ok bool) ***REMOVED***
	switch ***REMOVED***
	case '0' <= b && b <= '9':
		return b - '0', true
	case 'a' <= b && b <= 'f':
		return b - 'a' + 10, true
	case 'A' <= b && b <= 'F':
		return b - 'A' + 10, true
	***REMOVED***
	return 0, false
***REMOVED***

// Back off the parser by one token. Can only be done between calls to next().
// It makes the next advance() a no-op.
func (p *textParser) back() ***REMOVED*** p.backed = true ***REMOVED***

// Advances the parser and returns the new current token.
func (p *textParser) next() *token ***REMOVED***
	if p.backed || p.done ***REMOVED***
		p.backed = false
		return &p.cur
	***REMOVED***
	p.advance()
	if p.done ***REMOVED***
		p.cur.value = ""
	***REMOVED*** else if len(p.cur.value) > 0 && isQuote(p.cur.value[0]) ***REMOVED***
		// Look for multiple quoted strings separated by whitespace,
		// and concatenate them.
		cat := p.cur
		for ***REMOVED***
			p.skipWhitespace()
			if p.done || !isQuote(p.s[0]) ***REMOVED***
				break
			***REMOVED***
			p.advance()
			if p.cur.err != nil ***REMOVED***
				return &p.cur
			***REMOVED***
			cat.value += " " + p.cur.value
			cat.unquoted += p.cur.unquoted
		***REMOVED***
		p.done = false // parser may have seen EOF, but we want to return cat
		p.cur = cat
	***REMOVED***
	return &p.cur
***REMOVED***

func (p *textParser) consumeToken(s string) error ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return tok.err
	***REMOVED***
	if tok.value != s ***REMOVED***
		p.back()
		return p.errorf("expected %q, found %q", s, tok.value)
	***REMOVED***
	return nil
***REMOVED***

// Return a RequiredNotSetError indicating which required field was not set.
func (p *textParser) missingRequiredFieldError(sv reflect.Value) *RequiredNotSetError ***REMOVED***
	st := sv.Type()
	sprops := GetProperties(st)
	for i := 0; i < st.NumField(); i++ ***REMOVED***
		if !isNil(sv.Field(i)) ***REMOVED***
			continue
		***REMOVED***

		props := sprops.Prop[i]
		if props.Required ***REMOVED***
			return &RequiredNotSetError***REMOVED***fmt.Sprintf("%v.%v", st, props.OrigName)***REMOVED***
		***REMOVED***
	***REMOVED***
	return &RequiredNotSetError***REMOVED***fmt.Sprintf("%v.<unknown field name>", st)***REMOVED*** // should not happen
***REMOVED***

// Returns the index in the struct for the named field, as well as the parsed tag properties.
func structFieldByName(sprops *StructProperties, name string) (int, *Properties, bool) ***REMOVED***
	i, ok := sprops.decoderOrigNames[name]
	if ok ***REMOVED***
		return i, sprops.Prop[i], true
	***REMOVED***
	return -1, nil, false
***REMOVED***

// Consume a ':' from the input stream (if the next token is a colon),
// returning an error if a colon is needed but not present.
func (p *textParser) checkForColon(props *Properties, typ reflect.Type) *ParseError ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return tok.err
	***REMOVED***
	if tok.value != ":" ***REMOVED***
		// Colon is optional when the field is a group or message.
		needColon := true
		switch props.Wire ***REMOVED***
		case "group":
			needColon = false
		case "bytes":
			// A "bytes" field is either a message, a string, or a repeated field;
			// those three become *T, *string and []T respectively, so we can check for
			// this field being a pointer to a non-string.
			if typ.Kind() == reflect.Ptr ***REMOVED***
				// *T or *string
				if typ.Elem().Kind() == reflect.String ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if typ.Kind() == reflect.Slice ***REMOVED***
				// []T or []*T
				if typ.Elem().Kind() != reflect.Ptr ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if typ.Kind() == reflect.String ***REMOVED***
				// The proto3 exception is for a string field,
				// which requires a colon.
				break
			***REMOVED***
			needColon = false
		***REMOVED***
		if needColon ***REMOVED***
			return p.errorf("expected ':', found %q", tok.value)
		***REMOVED***
		p.back()
	***REMOVED***
	return nil
***REMOVED***

func (p *textParser) readStruct(sv reflect.Value, terminator string) error ***REMOVED***
	st := sv.Type()
	sprops := GetProperties(st)
	reqCount := sprops.reqCount
	var reqFieldErr error
	fieldSet := make(map[string]bool)
	// A struct is a sequence of "name: value", terminated by one of
	// '>' or '***REMOVED***', or the end of the input.  A name may also be
	// "[extension]" or "[type/url]".
	//
	// The whole struct can also be an expanded Any message, like:
	// [type/url] < ... struct contents ... >
	for ***REMOVED***
		tok := p.next()
		if tok.err != nil ***REMOVED***
			return tok.err
		***REMOVED***
		if tok.value == terminator ***REMOVED***
			break
		***REMOVED***
		if tok.value == "[" ***REMOVED***
			// Looks like an extension or an Any.
			//
			// TODO: Check whether we need to handle
			// namespace rooted names (e.g. ".something.Foo").
			extName, err := p.consumeExtName()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if s := strings.LastIndex(extName, "/"); s >= 0 ***REMOVED***
				// If it contains a slash, it's an Any type URL.
				messageName := extName[s+1:]
				mt := MessageType(messageName)
				if mt == nil ***REMOVED***
					return p.errorf("unrecognized message %q in google.protobuf.Any", messageName)
				***REMOVED***
				tok = p.next()
				if tok.err != nil ***REMOVED***
					return tok.err
				***REMOVED***
				// consume an optional colon
				if tok.value == ":" ***REMOVED***
					tok = p.next()
					if tok.err != nil ***REMOVED***
						return tok.err
					***REMOVED***
				***REMOVED***
				var terminator string
				switch tok.value ***REMOVED***
				case "<":
					terminator = ">"
				case "***REMOVED***":
					terminator = "***REMOVED***"
				default:
					return p.errorf("expected '***REMOVED***' or '<', found %q", tok.value)
				***REMOVED***
				v := reflect.New(mt.Elem())
				if pe := p.readStruct(v.Elem(), terminator); pe != nil ***REMOVED***
					return pe
				***REMOVED***
				b, err := Marshal(v.Interface().(Message))
				if err != nil ***REMOVED***
					return p.errorf("failed to marshal message of type %q: %v", messageName, err)
				***REMOVED***
				if fieldSet["type_url"] ***REMOVED***
					return p.errorf(anyRepeatedlyUnpacked, "type_url")
				***REMOVED***
				if fieldSet["value"] ***REMOVED***
					return p.errorf(anyRepeatedlyUnpacked, "value")
				***REMOVED***
				sv.FieldByName("TypeUrl").SetString(extName)
				sv.FieldByName("Value").SetBytes(b)
				fieldSet["type_url"] = true
				fieldSet["value"] = true
				continue
			***REMOVED***

			var desc *ExtensionDesc
			// This could be faster, but it's functional.
			// TODO: Do something smarter than a linear scan.
			for _, d := range RegisteredExtensions(reflect.New(st).Interface().(Message)) ***REMOVED***
				if d.Name == extName ***REMOVED***
					desc = d
					break
				***REMOVED***
			***REMOVED***
			if desc == nil ***REMOVED***
				return p.errorf("unrecognized extension %q", extName)
			***REMOVED***

			props := &Properties***REMOVED******REMOVED***
			props.Parse(desc.Tag)

			typ := reflect.TypeOf(desc.ExtensionType)
			if err := p.checkForColon(props, typ); err != nil ***REMOVED***
				return err
			***REMOVED***

			rep := desc.repeated()

			// Read the extension structure, and set it in
			// the value we're constructing.
			var ext reflect.Value
			if !rep ***REMOVED***
				ext = reflect.New(typ).Elem()
			***REMOVED*** else ***REMOVED***
				ext = reflect.New(typ.Elem()).Elem()
			***REMOVED***
			if err := p.readAny(ext, props); err != nil ***REMOVED***
				if _, ok := err.(*RequiredNotSetError); !ok ***REMOVED***
					return err
				***REMOVED***
				reqFieldErr = err
			***REMOVED***
			ep := sv.Addr().Interface().(Message)
			if !rep ***REMOVED***
				SetExtension(ep, desc, ext.Interface())
			***REMOVED*** else ***REMOVED***
				old, err := GetExtension(ep, desc)
				var sl reflect.Value
				if err == nil ***REMOVED***
					sl = reflect.ValueOf(old) // existing slice
				***REMOVED*** else ***REMOVED***
					sl = reflect.MakeSlice(typ, 0, 1)
				***REMOVED***
				sl = reflect.Append(sl, ext)
				SetExtension(ep, desc, sl.Interface())
			***REMOVED***
			if err := p.consumeOptionalSeparator(); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		// This is a normal, non-extension field.
		name := tok.value
		var dst reflect.Value
		fi, props, ok := structFieldByName(sprops, name)
		if ok ***REMOVED***
			dst = sv.Field(fi)
		***REMOVED*** else if oop, ok := sprops.OneofTypes[name]; ok ***REMOVED***
			// It is a oneof.
			props = oop.Prop
			nv := reflect.New(oop.Type.Elem())
			dst = nv.Elem().Field(0)
			field := sv.Field(oop.Field)
			if !field.IsNil() ***REMOVED***
				return p.errorf("field '%s' would overwrite already parsed oneof '%s'", name, sv.Type().Field(oop.Field).Name)
			***REMOVED***
			field.Set(nv)
		***REMOVED***
		if !dst.IsValid() ***REMOVED***
			return p.errorf("unknown field name %q in %v", name, st)
		***REMOVED***

		if dst.Kind() == reflect.Map ***REMOVED***
			// Consume any colon.
			if err := p.checkForColon(props, dst.Type()); err != nil ***REMOVED***
				return err
			***REMOVED***

			// Construct the map if it doesn't already exist.
			if dst.IsNil() ***REMOVED***
				dst.Set(reflect.MakeMap(dst.Type()))
			***REMOVED***
			key := reflect.New(dst.Type().Key()).Elem()
			val := reflect.New(dst.Type().Elem()).Elem()

			// The map entry should be this sequence of tokens:
			//	< key : KEY value : VALUE >
			// However, implementations may omit key or value, and technically
			// we should support them in any order.  See b/28924776 for a time
			// this went wrong.

			tok := p.next()
			var terminator string
			switch tok.value ***REMOVED***
			case "<":
				terminator = ">"
			case "***REMOVED***":
				terminator = "***REMOVED***"
			default:
				return p.errorf("expected '***REMOVED***' or '<', found %q", tok.value)
			***REMOVED***
			for ***REMOVED***
				tok := p.next()
				if tok.err != nil ***REMOVED***
					return tok.err
				***REMOVED***
				if tok.value == terminator ***REMOVED***
					break
				***REMOVED***
				switch tok.value ***REMOVED***
				case "key":
					if err := p.consumeToken(":"); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := p.readAny(key, props.mkeyprop); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := p.consumeOptionalSeparator(); err != nil ***REMOVED***
						return err
					***REMOVED***
				case "value":
					if err := p.checkForColon(props.mvalprop, dst.Type().Elem()); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := p.readAny(val, props.mvalprop); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := p.consumeOptionalSeparator(); err != nil ***REMOVED***
						return err
					***REMOVED***
				default:
					p.back()
					return p.errorf(`expected "key", "value", or %q, found %q`, terminator, tok.value)
				***REMOVED***
			***REMOVED***

			dst.SetMapIndex(key, val)
			continue
		***REMOVED***

		// Check that it's not already set if it's not a repeated field.
		if !props.Repeated && fieldSet[name] ***REMOVED***
			return p.errorf("non-repeated field %q was repeated", name)
		***REMOVED***

		if err := p.checkForColon(props, dst.Type()); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Parse into the field.
		fieldSet[name] = true
		if err := p.readAny(dst, props); err != nil ***REMOVED***
			if _, ok := err.(*RequiredNotSetError); !ok ***REMOVED***
				return err
			***REMOVED***
			reqFieldErr = err
		***REMOVED***
		if props.Required ***REMOVED***
			reqCount--
		***REMOVED***

		if err := p.consumeOptionalSeparator(); err != nil ***REMOVED***
			return err
		***REMOVED***

	***REMOVED***

	if reqCount > 0 ***REMOVED***
		return p.missingRequiredFieldError(sv)
	***REMOVED***
	return reqFieldErr
***REMOVED***

// consumeExtName consumes extension name or expanded Any type URL and the
// following ']'. It returns the name or URL consumed.
func (p *textParser) consumeExtName() (string, error) ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return "", tok.err
	***REMOVED***

	// If extension name or type url is quoted, it's a single token.
	if len(tok.value) > 2 && isQuote(tok.value[0]) && tok.value[len(tok.value)-1] == tok.value[0] ***REMOVED***
		name, err := unquoteC(tok.value[1:len(tok.value)-1], rune(tok.value[0]))
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		return name, p.consumeToken("]")
	***REMOVED***

	// Consume everything up to "]"
	var parts []string
	for tok.value != "]" ***REMOVED***
		parts = append(parts, tok.value)
		tok = p.next()
		if tok.err != nil ***REMOVED***
			return "", p.errorf("unrecognized type_url or extension name: %s", tok.err)
		***REMOVED***
	***REMOVED***
	return strings.Join(parts, ""), nil
***REMOVED***

// consumeOptionalSeparator consumes an optional semicolon or comma.
// It is used in readStruct to provide backward compatibility.
func (p *textParser) consumeOptionalSeparator() error ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return tok.err
	***REMOVED***
	if tok.value != ";" && tok.value != "," ***REMOVED***
		p.back()
	***REMOVED***
	return nil
***REMOVED***

func (p *textParser) readAny(v reflect.Value, props *Properties) error ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return tok.err
	***REMOVED***
	if tok.value == "" ***REMOVED***
		return p.errorf("unexpected EOF")
	***REMOVED***

	switch fv := v; fv.Kind() ***REMOVED***
	case reflect.Slice:
		at := v.Type()
		if at.Elem().Kind() == reflect.Uint8 ***REMOVED***
			// Special case for []byte
			if tok.value[0] != '"' && tok.value[0] != '\'' ***REMOVED***
				// Deliberately written out here, as the error after
				// this switch statement would write "invalid []byte: ...",
				// which is not as user-friendly.
				return p.errorf("invalid string: %v", tok.value)
			***REMOVED***
			bytes := []byte(tok.unquoted)
			fv.Set(reflect.ValueOf(bytes))
			return nil
		***REMOVED***
		// Repeated field.
		if tok.value == "[" ***REMOVED***
			// Repeated field with list notation, like [1,2,3].
			for ***REMOVED***
				fv.Set(reflect.Append(fv, reflect.New(at.Elem()).Elem()))
				err := p.readAny(fv.Index(fv.Len()-1), props)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				tok := p.next()
				if tok.err != nil ***REMOVED***
					return tok.err
				***REMOVED***
				if tok.value == "]" ***REMOVED***
					break
				***REMOVED***
				if tok.value != "," ***REMOVED***
					return p.errorf("Expected ']' or ',' found %q", tok.value)
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***
		// One value of the repeated field.
		p.back()
		fv.Set(reflect.Append(fv, reflect.New(at.Elem()).Elem()))
		return p.readAny(fv.Index(fv.Len()-1), props)
	case reflect.Bool:
		// true/1/t/True or false/f/0/False.
		switch tok.value ***REMOVED***
		case "true", "1", "t", "True":
			fv.SetBool(true)
			return nil
		case "false", "0", "f", "False":
			fv.SetBool(false)
			return nil
		***REMOVED***
	case reflect.Float32, reflect.Float64:
		v := tok.value
		// Ignore 'f' for compatibility with output generated by C++, but don't
		// remove 'f' when the value is "-inf" or "inf".
		if strings.HasSuffix(v, "f") && tok.value != "-inf" && tok.value != "inf" ***REMOVED***
			v = v[:len(v)-1]
		***REMOVED***
		if f, err := strconv.ParseFloat(v, fv.Type().Bits()); err == nil ***REMOVED***
			fv.SetFloat(f)
			return nil
		***REMOVED***
	case reflect.Int32:
		if x, err := strconv.ParseInt(tok.value, 0, 32); err == nil ***REMOVED***
			fv.SetInt(x)
			return nil
		***REMOVED***

		if len(props.Enum) == 0 ***REMOVED***
			break
		***REMOVED***
		m, ok := enumValueMaps[props.Enum]
		if !ok ***REMOVED***
			break
		***REMOVED***
		x, ok := m[tok.value]
		if !ok ***REMOVED***
			break
		***REMOVED***
		fv.SetInt(int64(x))
		return nil
	case reflect.Int64:
		if x, err := strconv.ParseInt(tok.value, 0, 64); err == nil ***REMOVED***
			fv.SetInt(x)
			return nil
		***REMOVED***

	case reflect.Ptr:
		// A basic field (indirected through pointer), or a repeated message/group
		p.back()
		fv.Set(reflect.New(fv.Type().Elem()))
		return p.readAny(fv.Elem(), props)
	case reflect.String:
		if tok.value[0] == '"' || tok.value[0] == '\'' ***REMOVED***
			fv.SetString(tok.unquoted)
			return nil
		***REMOVED***
	case reflect.Struct:
		var terminator string
		switch tok.value ***REMOVED***
		case "***REMOVED***":
			terminator = "***REMOVED***"
		case "<":
			terminator = ">"
		default:
			return p.errorf("expected '***REMOVED***' or '<', found %q", tok.value)
		***REMOVED***
		// TODO: Handle nested messages which implement encoding.TextUnmarshaler.
		return p.readStruct(fv, terminator)
	case reflect.Uint32:
		if x, err := strconv.ParseUint(tok.value, 0, 32); err == nil ***REMOVED***
			fv.SetUint(x)
			return nil
		***REMOVED***
	case reflect.Uint64:
		if x, err := strconv.ParseUint(tok.value, 0, 64); err == nil ***REMOVED***
			fv.SetUint(x)
			return nil
		***REMOVED***
	***REMOVED***
	return p.errorf("invalid %v: %v", v.Type(), tok.value)
***REMOVED***

// UnmarshalText reads a protocol buffer in Text format. UnmarshalText resets pb
// before starting to unmarshal, so any existing data in pb is always removed.
// If a required field is not set and no other error occurs,
// UnmarshalText returns *RequiredNotSetError.
func UnmarshalText(s string, pb Message) error ***REMOVED***
	if um, ok := pb.(encoding.TextUnmarshaler); ok ***REMOVED***
		err := um.UnmarshalText([]byte(s))
		return err
	***REMOVED***
	pb.Reset()
	v := reflect.ValueOf(pb)
	if pe := newTextParser(s).readStruct(v.Elem(), ""); pe != nil ***REMOVED***
		return pe
	***REMOVED***
	return nil
***REMOVED***
