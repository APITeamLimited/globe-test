// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/prototext"
	protoV2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const wrapTextUnmarshalV2 = false

// ParseError is returned by UnmarshalText.
type ParseError struct ***REMOVED***
	Message string

	// Deprecated: Do not use.
	Line, Offset int
***REMOVED***

func (e *ParseError) Error() string ***REMOVED***
	if wrapTextUnmarshalV2 ***REMOVED***
		return e.Message
	***REMOVED***
	if e.Line == 1 ***REMOVED***
		return fmt.Sprintf("line 1.%d: %v", e.Offset, e.Message)
	***REMOVED***
	return fmt.Sprintf("line %d: %v", e.Line, e.Message)
***REMOVED***

// UnmarshalText parses a proto text formatted string into m.
func UnmarshalText(s string, m Message) error ***REMOVED***
	if u, ok := m.(encoding.TextUnmarshaler); ok ***REMOVED***
		return u.UnmarshalText([]byte(s))
	***REMOVED***

	m.Reset()
	mi := MessageV2(m)

	if wrapTextUnmarshalV2 ***REMOVED***
		err := prototext.UnmarshalOptions***REMOVED***
			AllowPartial: true,
		***REMOVED***.Unmarshal([]byte(s), mi)
		if err != nil ***REMOVED***
			return &ParseError***REMOVED***Message: err.Error()***REMOVED***
		***REMOVED***
		return checkRequiredNotSet(mi)
	***REMOVED*** else ***REMOVED***
		if err := newTextParser(s).unmarshalMessage(mi.ProtoReflect(), ""); err != nil ***REMOVED***
			return err
		***REMOVED***
		return checkRequiredNotSet(mi)
	***REMOVED***
***REMOVED***

type textParser struct ***REMOVED***
	s            string // remaining input
	done         bool   // whether the parsing is finished (success or error)
	backed       bool   // whether back() was called
	offset, line int
	cur          token
***REMOVED***

type token struct ***REMOVED***
	value    string
	err      *ParseError
	line     int    // line number
	offset   int    // byte number from start of input, not start of line
	unquoted string // the unquoted version of value, if it was a quoted string
***REMOVED***

func newTextParser(s string) *textParser ***REMOVED***
	p := new(textParser)
	p.s = s
	p.line = 1
	p.cur.line = 1
	return p
***REMOVED***

func (p *textParser) unmarshalMessage(m protoreflect.Message, terminator string) (err error) ***REMOVED***
	md := m.Descriptor()
	fds := md.Fields()

	// A struct is a sequence of "name: value", terminated by one of
	// '>' or '***REMOVED***', or the end of the input.  A name may also be
	// "[extension]" or "[type/url]".
	//
	// The whole struct can also be an expanded Any message, like:
	// [type/url] < ... struct contents ... >
	seen := make(map[protoreflect.FieldNumber]bool)
	for ***REMOVED***
		tok := p.next()
		if tok.err != nil ***REMOVED***
			return tok.err
		***REMOVED***
		if tok.value == terminator ***REMOVED***
			break
		***REMOVED***
		if tok.value == "[" ***REMOVED***
			if err := p.unmarshalExtensionOrAny(m, seen); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		// This is a normal, non-extension field.
		name := protoreflect.Name(tok.value)
		fd := fds.ByName(name)
		switch ***REMOVED***
		case fd == nil:
			gd := fds.ByName(protoreflect.Name(strings.ToLower(string(name))))
			if gd != nil && gd.Kind() == protoreflect.GroupKind && gd.Message().Name() == name ***REMOVED***
				fd = gd
			***REMOVED***
		case fd.Kind() == protoreflect.GroupKind && fd.Message().Name() != name:
			fd = nil
		case fd.IsWeak() && fd.Message().IsPlaceholder():
			fd = nil
		***REMOVED***
		if fd == nil ***REMOVED***
			typeName := string(md.FullName())
			if m, ok := m.Interface().(Message); ok ***REMOVED***
				t := reflect.TypeOf(m)
				if t.Kind() == reflect.Ptr ***REMOVED***
					typeName = t.Elem().String()
				***REMOVED***
			***REMOVED***
			return p.errorf("unknown field name %q in %v", name, typeName)
		***REMOVED***
		if od := fd.ContainingOneof(); od != nil && m.WhichOneof(od) != nil ***REMOVED***
			return p.errorf("field '%s' would overwrite already parsed oneof '%s'", name, od.Name())
		***REMOVED***
		if fd.Cardinality() != protoreflect.Repeated && seen[fd.Number()] ***REMOVED***
			return p.errorf("non-repeated field %q was repeated", fd.Name())
		***REMOVED***
		seen[fd.Number()] = true

		// Consume any colon.
		if err := p.checkForColon(fd); err != nil ***REMOVED***
			return err
		***REMOVED***

		// Parse into the field.
		v := m.Get(fd)
		if !m.Has(fd) && (fd.IsList() || fd.IsMap() || fd.Message() != nil) ***REMOVED***
			v = m.Mutable(fd)
		***REMOVED***
		if v, err = p.unmarshalValue(v, fd); err != nil ***REMOVED***
			return err
		***REMOVED***
		m.Set(fd, v)

		if err := p.consumeOptionalSeparator(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (p *textParser) unmarshalExtensionOrAny(m protoreflect.Message, seen map[protoreflect.FieldNumber]bool) error ***REMOVED***
	name, err := p.consumeExtensionOrAnyName()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If it contains a slash, it's an Any type URL.
	if slashIdx := strings.LastIndex(name, "/"); slashIdx >= 0 ***REMOVED***
		tok := p.next()
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

		mt, err := protoregistry.GlobalTypes.FindMessageByURL(name)
		if err != nil ***REMOVED***
			return p.errorf("unrecognized message %q in google.protobuf.Any", name[slashIdx+len("/"):])
		***REMOVED***
		m2 := mt.New()
		if err := p.unmarshalMessage(m2, terminator); err != nil ***REMOVED***
			return err
		***REMOVED***
		b, err := protoV2.Marshal(m2.Interface())
		if err != nil ***REMOVED***
			return p.errorf("failed to marshal message of type %q: %v", name[slashIdx+len("/"):], err)
		***REMOVED***

		urlFD := m.Descriptor().Fields().ByName("type_url")
		valFD := m.Descriptor().Fields().ByName("value")
		if seen[urlFD.Number()] ***REMOVED***
			return p.errorf("Any message unpacked multiple times, or %q already set", urlFD.Name())
		***REMOVED***
		if seen[valFD.Number()] ***REMOVED***
			return p.errorf("Any message unpacked multiple times, or %q already set", valFD.Name())
		***REMOVED***
		m.Set(urlFD, protoreflect.ValueOfString(name))
		m.Set(valFD, protoreflect.ValueOfBytes(b))
		seen[urlFD.Number()] = true
		seen[valFD.Number()] = true
		return nil
	***REMOVED***

	xname := protoreflect.FullName(name)
	xt, _ := protoregistry.GlobalTypes.FindExtensionByName(xname)
	if xt == nil && isMessageSet(m.Descriptor()) ***REMOVED***
		xt, _ = protoregistry.GlobalTypes.FindExtensionByName(xname.Append("message_set_extension"))
	***REMOVED***
	if xt == nil ***REMOVED***
		return p.errorf("unrecognized extension %q", name)
	***REMOVED***
	fd := xt.TypeDescriptor()
	if fd.ContainingMessage().FullName() != m.Descriptor().FullName() ***REMOVED***
		return p.errorf("extension field %q does not extend message %q", name, m.Descriptor().FullName())
	***REMOVED***

	if err := p.checkForColon(fd); err != nil ***REMOVED***
		return err
	***REMOVED***

	v := m.Get(fd)
	if !m.Has(fd) && (fd.IsList() || fd.IsMap() || fd.Message() != nil) ***REMOVED***
		v = m.Mutable(fd)
	***REMOVED***
	v, err = p.unmarshalValue(v, fd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	m.Set(fd, v)
	return p.consumeOptionalSeparator()
***REMOVED***

func (p *textParser) unmarshalValue(v protoreflect.Value, fd protoreflect.FieldDescriptor) (protoreflect.Value, error) ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return v, tok.err
	***REMOVED***
	if tok.value == "" ***REMOVED***
		return v, p.errorf("unexpected EOF")
	***REMOVED***

	switch ***REMOVED***
	case fd.IsList():
		lv := v.List()
		var err error
		if tok.value == "[" ***REMOVED***
			// Repeated field with list notation, like [1,2,3].
			for ***REMOVED***
				vv := lv.NewElement()
				vv, err = p.unmarshalSingularValue(vv, fd)
				if err != nil ***REMOVED***
					return v, err
				***REMOVED***
				lv.Append(vv)

				tok := p.next()
				if tok.err != nil ***REMOVED***
					return v, tok.err
				***REMOVED***
				if tok.value == "]" ***REMOVED***
					break
				***REMOVED***
				if tok.value != "," ***REMOVED***
					return v, p.errorf("Expected ']' or ',' found %q", tok.value)
				***REMOVED***
			***REMOVED***
			return v, nil
		***REMOVED***

		// One value of the repeated field.
		p.back()
		vv := lv.NewElement()
		vv, err = p.unmarshalSingularValue(vv, fd)
		if err != nil ***REMOVED***
			return v, err
		***REMOVED***
		lv.Append(vv)
		return v, nil
	case fd.IsMap():
		// The map entry should be this sequence of tokens:
		//	< key : KEY value : VALUE >
		// However, implementations may omit key or value, and technically
		// we should support them in any order.
		var terminator string
		switch tok.value ***REMOVED***
		case "<":
			terminator = ">"
		case "***REMOVED***":
			terminator = "***REMOVED***"
		default:
			return v, p.errorf("expected '***REMOVED***' or '<', found %q", tok.value)
		***REMOVED***

		keyFD := fd.MapKey()
		valFD := fd.MapValue()

		mv := v.Map()
		kv := keyFD.Default()
		vv := mv.NewValue()
		for ***REMOVED***
			tok := p.next()
			if tok.err != nil ***REMOVED***
				return v, tok.err
			***REMOVED***
			if tok.value == terminator ***REMOVED***
				break
			***REMOVED***
			var err error
			switch tok.value ***REMOVED***
			case "key":
				if err := p.consumeToken(":"); err != nil ***REMOVED***
					return v, err
				***REMOVED***
				if kv, err = p.unmarshalSingularValue(kv, keyFD); err != nil ***REMOVED***
					return v, err
				***REMOVED***
				if err := p.consumeOptionalSeparator(); err != nil ***REMOVED***
					return v, err
				***REMOVED***
			case "value":
				if err := p.checkForColon(valFD); err != nil ***REMOVED***
					return v, err
				***REMOVED***
				if vv, err = p.unmarshalSingularValue(vv, valFD); err != nil ***REMOVED***
					return v, err
				***REMOVED***
				if err := p.consumeOptionalSeparator(); err != nil ***REMOVED***
					return v, err
				***REMOVED***
			default:
				p.back()
				return v, p.errorf(`expected "key", "value", or %q, found %q`, terminator, tok.value)
			***REMOVED***
		***REMOVED***
		mv.Set(kv.MapKey(), vv)
		return v, nil
	default:
		p.back()
		return p.unmarshalSingularValue(v, fd)
	***REMOVED***
***REMOVED***

func (p *textParser) unmarshalSingularValue(v protoreflect.Value, fd protoreflect.FieldDescriptor) (protoreflect.Value, error) ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return v, tok.err
	***REMOVED***
	if tok.value == "" ***REMOVED***
		return v, p.errorf("unexpected EOF")
	***REMOVED***

	switch fd.Kind() ***REMOVED***
	case protoreflect.BoolKind:
		switch tok.value ***REMOVED***
		case "true", "1", "t", "True":
			return protoreflect.ValueOfBool(true), nil
		case "false", "0", "f", "False":
			return protoreflect.ValueOfBool(false), nil
		***REMOVED***
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if x, err := strconv.ParseInt(tok.value, 0, 32); err == nil ***REMOVED***
			return protoreflect.ValueOfInt32(int32(x)), nil
		***REMOVED***

		// The C++ parser accepts large positive hex numbers that uses
		// two's complement arithmetic to represent negative numbers.
		// This feature is here for backwards compatibility with C++.
		if strings.HasPrefix(tok.value, "0x") ***REMOVED***
			if x, err := strconv.ParseUint(tok.value, 0, 32); err == nil ***REMOVED***
				return protoreflect.ValueOfInt32(int32(-(int64(^x) + 1))), nil
			***REMOVED***
		***REMOVED***
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if x, err := strconv.ParseInt(tok.value, 0, 64); err == nil ***REMOVED***
			return protoreflect.ValueOfInt64(int64(x)), nil
		***REMOVED***

		// The C++ parser accepts large positive hex numbers that uses
		// two's complement arithmetic to represent negative numbers.
		// This feature is here for backwards compatibility with C++.
		if strings.HasPrefix(tok.value, "0x") ***REMOVED***
			if x, err := strconv.ParseUint(tok.value, 0, 64); err == nil ***REMOVED***
				return protoreflect.ValueOfInt64(int64(-(int64(^x) + 1))), nil
			***REMOVED***
		***REMOVED***
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if x, err := strconv.ParseUint(tok.value, 0, 32); err == nil ***REMOVED***
			return protoreflect.ValueOfUint32(uint32(x)), nil
		***REMOVED***
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if x, err := strconv.ParseUint(tok.value, 0, 64); err == nil ***REMOVED***
			return protoreflect.ValueOfUint64(uint64(x)), nil
		***REMOVED***
	case protoreflect.FloatKind:
		// Ignore 'f' for compatibility with output generated by C++,
		// but don't remove 'f' when the value is "-inf" or "inf".
		v := tok.value
		if strings.HasSuffix(v, "f") && v != "-inf" && v != "inf" ***REMOVED***
			v = v[:len(v)-len("f")]
		***REMOVED***
		if x, err := strconv.ParseFloat(v, 32); err == nil ***REMOVED***
			return protoreflect.ValueOfFloat32(float32(x)), nil
		***REMOVED***
	case protoreflect.DoubleKind:
		// Ignore 'f' for compatibility with output generated by C++,
		// but don't remove 'f' when the value is "-inf" or "inf".
		v := tok.value
		if strings.HasSuffix(v, "f") && v != "-inf" && v != "inf" ***REMOVED***
			v = v[:len(v)-len("f")]
		***REMOVED***
		if x, err := strconv.ParseFloat(v, 64); err == nil ***REMOVED***
			return protoreflect.ValueOfFloat64(float64(x)), nil
		***REMOVED***
	case protoreflect.StringKind:
		if isQuote(tok.value[0]) ***REMOVED***
			return protoreflect.ValueOfString(tok.unquoted), nil
		***REMOVED***
	case protoreflect.BytesKind:
		if isQuote(tok.value[0]) ***REMOVED***
			return protoreflect.ValueOfBytes([]byte(tok.unquoted)), nil
		***REMOVED***
	case protoreflect.EnumKind:
		if x, err := strconv.ParseInt(tok.value, 0, 32); err == nil ***REMOVED***
			return protoreflect.ValueOfEnum(protoreflect.EnumNumber(x)), nil
		***REMOVED***
		vd := fd.Enum().Values().ByName(protoreflect.Name(tok.value))
		if vd != nil ***REMOVED***
			return protoreflect.ValueOfEnum(vd.Number()), nil
		***REMOVED***
	case protoreflect.MessageKind, protoreflect.GroupKind:
		var terminator string
		switch tok.value ***REMOVED***
		case "***REMOVED***":
			terminator = "***REMOVED***"
		case "<":
			terminator = ">"
		default:
			return v, p.errorf("expected '***REMOVED***' or '<', found %q", tok.value)
		***REMOVED***
		err := p.unmarshalMessage(v.Message(), terminator)
		return v, err
	default:
		panic(fmt.Sprintf("invalid kind %v", fd.Kind()))
	***REMOVED***
	return v, p.errorf("invalid %v: %v", fd.Kind(), tok.value)
***REMOVED***

// Consume a ':' from the input stream (if the next token is a colon),
// returning an error if a colon is needed but not present.
func (p *textParser) checkForColon(fd protoreflect.FieldDescriptor) *ParseError ***REMOVED***
	tok := p.next()
	if tok.err != nil ***REMOVED***
		return tok.err
	***REMOVED***
	if tok.value != ":" ***REMOVED***
		if fd.Message() == nil ***REMOVED***
			return p.errorf("expected ':', found %q", tok.value)
		***REMOVED***
		p.back()
	***REMOVED***
	return nil
***REMOVED***

// consumeExtensionOrAnyName consumes an extension name or an Any type URL and
// the following ']'. It returns the name or URL consumed.
func (p *textParser) consumeExtensionOrAnyName() (string, error) ***REMOVED***
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
		if p.done && tok.value != "]" ***REMOVED***
			return "", p.errorf("unclosed type_url or extension name")
		***REMOVED***
	***REMOVED***
	return strings.Join(parts, ""), nil
***REMOVED***

// consumeOptionalSeparator consumes an optional semicolon or comma.
// It is used in unmarshalMessage to provide backward compatibility.
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

func (p *textParser) errorf(format string, a ...interface***REMOVED******REMOVED***) *ParseError ***REMOVED***
	pe := &ParseError***REMOVED***fmt.Sprintf(format, a...), p.cur.line, p.cur.offset***REMOVED***
	p.cur.err = pe
	p.done = true
	return pe
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

var errBadUTF8 = errors.New("proto: bad UTF-8")

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
	case '0', '1', '2', '3', '4', '5', '6', '7':
		if len(s) < 2 ***REMOVED***
			return "", "", fmt.Errorf(`\%c requires 2 following digits`, r)
		***REMOVED***
		ss := string(r) + s[:2]
		s = s[2:]
		i, err := strconv.ParseUint(ss, 8, 8)
		if err != nil ***REMOVED***
			return "", "", fmt.Errorf(`\%s contains non-octal digits`, ss)
		***REMOVED***
		return string([]byte***REMOVED***byte(i)***REMOVED***), s, nil
	case 'x', 'X', 'u', 'U':
		var n int
		switch r ***REMOVED***
		case 'x', 'X':
			n = 2
		case 'u':
			n = 4
		case 'U':
			n = 8
		***REMOVED***
		if len(s) < n ***REMOVED***
			return "", "", fmt.Errorf(`\%c requires %d following digits`, r, n)
		***REMOVED***
		ss := s[:n]
		s = s[n:]
		i, err := strconv.ParseUint(ss, 16, 64)
		if err != nil ***REMOVED***
			return "", "", fmt.Errorf(`\%c%s contains non-hexadecimal digits`, r, ss)
		***REMOVED***
		if r == 'x' || r == 'X' ***REMOVED***
			return string([]byte***REMOVED***byte(i)***REMOVED***), s, nil
		***REMOVED***
		if i > utf8.MaxRune ***REMOVED***
			return "", "", fmt.Errorf(`\%c%s is not a valid Unicode code point`, r, ss)
		***REMOVED***
		return string(rune(i)), s, nil
	***REMOVED***
	return "", "", fmt.Errorf(`unknown escape \%c`, r)
***REMOVED***

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
