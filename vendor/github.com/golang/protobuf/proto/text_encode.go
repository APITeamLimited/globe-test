// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const wrapTextMarshalV2 = false

// TextMarshaler is a configurable text format marshaler.
type TextMarshaler struct ***REMOVED***
	Compact   bool // use compact text format (one line)
	ExpandAny bool // expand google.protobuf.Any messages of known types
***REMOVED***

// Marshal writes the proto text format of m to w.
func (tm *TextMarshaler) Marshal(w io.Writer, m Message) error ***REMOVED***
	b, err := tm.marshal(m)
	if len(b) > 0 ***REMOVED***
		if _, err := w.Write(b); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

// Text returns a proto text formatted string of m.
func (tm *TextMarshaler) Text(m Message) string ***REMOVED***
	b, _ := tm.marshal(m)
	return string(b)
***REMOVED***

func (tm *TextMarshaler) marshal(m Message) ([]byte, error) ***REMOVED***
	mr := MessageReflect(m)
	if mr == nil || !mr.IsValid() ***REMOVED***
		return []byte("<nil>"), nil
	***REMOVED***

	if wrapTextMarshalV2 ***REMOVED***
		if m, ok := m.(encoding.TextMarshaler); ok ***REMOVED***
			return m.MarshalText()
		***REMOVED***

		opts := prototext.MarshalOptions***REMOVED***
			AllowPartial: true,
			EmitUnknown:  true,
		***REMOVED***
		if !tm.Compact ***REMOVED***
			opts.Indent = "  "
		***REMOVED***
		if !tm.ExpandAny ***REMOVED***
			opts.Resolver = (*protoregistry.Types)(nil)
		***REMOVED***
		return opts.Marshal(mr.Interface())
	***REMOVED*** else ***REMOVED***
		w := &textWriter***REMOVED***
			compact:   tm.Compact,
			expandAny: tm.ExpandAny,
			complete:  true,
		***REMOVED***

		if m, ok := m.(encoding.TextMarshaler); ok ***REMOVED***
			b, err := m.MarshalText()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			w.Write(b)
			return w.buf, nil
		***REMOVED***

		err := w.writeMessage(mr)
		return w.buf, err
	***REMOVED***
***REMOVED***

var (
	defaultTextMarshaler = TextMarshaler***REMOVED******REMOVED***
	compactTextMarshaler = TextMarshaler***REMOVED***Compact: true***REMOVED***
)

// MarshalText writes the proto text format of m to w.
func MarshalText(w io.Writer, m Message) error ***REMOVED*** return defaultTextMarshaler.Marshal(w, m) ***REMOVED***

// MarshalTextString returns a proto text formatted string of m.
func MarshalTextString(m Message) string ***REMOVED*** return defaultTextMarshaler.Text(m) ***REMOVED***

// CompactText writes the compact proto text format of m to w.
func CompactText(w io.Writer, m Message) error ***REMOVED*** return compactTextMarshaler.Marshal(w, m) ***REMOVED***

// CompactTextString returns a compact proto text formatted string of m.
func CompactTextString(m Message) string ***REMOVED*** return compactTextMarshaler.Text(m) ***REMOVED***

var (
	newline         = []byte("\n")
	endBraceNewline = []byte("***REMOVED***\n")
	posInf          = []byte("inf")
	negInf          = []byte("-inf")
	nan             = []byte("nan")
)

// textWriter is an io.Writer that tracks its indentation level.
type textWriter struct ***REMOVED***
	compact   bool // same as TextMarshaler.Compact
	expandAny bool // same as TextMarshaler.ExpandAny
	complete  bool // whether the current position is a complete line
	indent    int  // indentation level; never negative
	buf       []byte
***REMOVED***

func (w *textWriter) Write(p []byte) (n int, _ error) ***REMOVED***
	newlines := bytes.Count(p, newline)
	if newlines == 0 ***REMOVED***
		if !w.compact && w.complete ***REMOVED***
			w.writeIndent()
		***REMOVED***
		w.buf = append(w.buf, p...)
		w.complete = false
		return len(p), nil
	***REMOVED***

	frags := bytes.SplitN(p, newline, newlines+1)
	if w.compact ***REMOVED***
		for i, frag := range frags ***REMOVED***
			if i > 0 ***REMOVED***
				w.buf = append(w.buf, ' ')
				n++
			***REMOVED***
			w.buf = append(w.buf, frag...)
			n += len(frag)
		***REMOVED***
		return n, nil
	***REMOVED***

	for i, frag := range frags ***REMOVED***
		if w.complete ***REMOVED***
			w.writeIndent()
		***REMOVED***
		w.buf = append(w.buf, frag...)
		n += len(frag)
		if i+1 < len(frags) ***REMOVED***
			w.buf = append(w.buf, '\n')
			n++
		***REMOVED***
	***REMOVED***
	w.complete = len(frags[len(frags)-1]) == 0
	return n, nil
***REMOVED***

func (w *textWriter) WriteByte(c byte) error ***REMOVED***
	if w.compact && c == '\n' ***REMOVED***
		c = ' '
	***REMOVED***
	if !w.compact && w.complete ***REMOVED***
		w.writeIndent()
	***REMOVED***
	w.buf = append(w.buf, c)
	w.complete = c == '\n'
	return nil
***REMOVED***

func (w *textWriter) writeName(fd protoreflect.FieldDescriptor) ***REMOVED***
	if !w.compact && w.complete ***REMOVED***
		w.writeIndent()
	***REMOVED***
	w.complete = false

	if fd.Kind() != protoreflect.GroupKind ***REMOVED***
		w.buf = append(w.buf, fd.Name()...)
		w.WriteByte(':')
	***REMOVED*** else ***REMOVED***
		// Use message type name for group field name.
		w.buf = append(w.buf, fd.Message().Name()...)
	***REMOVED***

	if !w.compact ***REMOVED***
		w.WriteByte(' ')
	***REMOVED***
***REMOVED***

func requiresQuotes(u string) bool ***REMOVED***
	// When type URL contains any characters except [0-9A-Za-z./\-]*, it must be quoted.
	for _, ch := range u ***REMOVED***
		switch ***REMOVED***
		case ch == '.' || ch == '/' || ch == '_':
			continue
		case '0' <= ch && ch <= '9':
			continue
		case 'A' <= ch && ch <= 'Z':
			continue
		case 'a' <= ch && ch <= 'z':
			continue
		default:
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// writeProto3Any writes an expanded google.protobuf.Any message.
//
// It returns (false, nil) if sv value can't be unmarshaled (e.g. because
// required messages are not linked in).
//
// It returns (true, error) when sv was written in expanded format or an error
// was encountered.
func (w *textWriter) writeProto3Any(m protoreflect.Message) (bool, error) ***REMOVED***
	md := m.Descriptor()
	fdURL := md.Fields().ByName("type_url")
	fdVal := md.Fields().ByName("value")

	url := m.Get(fdURL).String()
	mt, err := protoregistry.GlobalTypes.FindMessageByURL(url)
	if err != nil ***REMOVED***
		return false, nil
	***REMOVED***

	b := m.Get(fdVal).Bytes()
	m2 := mt.New()
	if err := proto.Unmarshal(b, m2.Interface()); err != nil ***REMOVED***
		return false, nil
	***REMOVED***
	w.Write([]byte("["))
	if requiresQuotes(url) ***REMOVED***
		w.writeQuotedString(url)
	***REMOVED*** else ***REMOVED***
		w.Write([]byte(url))
	***REMOVED***
	if w.compact ***REMOVED***
		w.Write([]byte("]:<"))
	***REMOVED*** else ***REMOVED***
		w.Write([]byte("]: <\n"))
		w.indent++
	***REMOVED***
	if err := w.writeMessage(m2); err != nil ***REMOVED***
		return true, err
	***REMOVED***
	if w.compact ***REMOVED***
		w.Write([]byte("> "))
	***REMOVED*** else ***REMOVED***
		w.indent--
		w.Write([]byte(">\n"))
	***REMOVED***
	return true, nil
***REMOVED***

func (w *textWriter) writeMessage(m protoreflect.Message) error ***REMOVED***
	md := m.Descriptor()
	if w.expandAny && md.FullName() == "google.protobuf.Any" ***REMOVED***
		if canExpand, err := w.writeProto3Any(m); canExpand ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	fds := md.Fields()
	for i := 0; i < fds.Len(); ***REMOVED***
		fd := fds.Get(i)
		if od := fd.ContainingOneof(); od != nil ***REMOVED***
			fd = m.WhichOneof(od)
			i += od.Fields().Len()
		***REMOVED*** else ***REMOVED***
			i++
		***REMOVED***
		if fd == nil || !m.Has(fd) ***REMOVED***
			continue
		***REMOVED***

		switch ***REMOVED***
		case fd.IsList():
			lv := m.Get(fd).List()
			for j := 0; j < lv.Len(); j++ ***REMOVED***
				w.writeName(fd)
				v := lv.Get(j)
				if err := w.writeSingularValue(v, fd); err != nil ***REMOVED***
					return err
				***REMOVED***
				w.WriteByte('\n')
			***REMOVED***
		case fd.IsMap():
			kfd := fd.MapKey()
			vfd := fd.MapValue()
			mv := m.Get(fd).Map()

			type entry struct***REMOVED*** key, val protoreflect.Value ***REMOVED***
			var entries []entry
			mv.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool ***REMOVED***
				entries = append(entries, entry***REMOVED***k.Value(), v***REMOVED***)
				return true
			***REMOVED***)
			sort.Slice(entries, func(i, j int) bool ***REMOVED***
				switch kfd.Kind() ***REMOVED***
				case protoreflect.BoolKind:
					return !entries[i].key.Bool() && entries[j].key.Bool()
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
					return entries[i].key.Int() < entries[j].key.Int()
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					return entries[i].key.Uint() < entries[j].key.Uint()
				case protoreflect.StringKind:
					return entries[i].key.String() < entries[j].key.String()
				default:
					panic("invalid kind")
				***REMOVED***
			***REMOVED***)
			for _, entry := range entries ***REMOVED***
				w.writeName(fd)
				w.WriteByte('<')
				if !w.compact ***REMOVED***
					w.WriteByte('\n')
				***REMOVED***
				w.indent++
				w.writeName(kfd)
				if err := w.writeSingularValue(entry.key, kfd); err != nil ***REMOVED***
					return err
				***REMOVED***
				w.WriteByte('\n')
				w.writeName(vfd)
				if err := w.writeSingularValue(entry.val, vfd); err != nil ***REMOVED***
					return err
				***REMOVED***
				w.WriteByte('\n')
				w.indent--
				w.WriteByte('>')
				w.WriteByte('\n')
			***REMOVED***
		default:
			w.writeName(fd)
			if err := w.writeSingularValue(m.Get(fd), fd); err != nil ***REMOVED***
				return err
			***REMOVED***
			w.WriteByte('\n')
		***REMOVED***
	***REMOVED***

	if b := m.GetUnknown(); len(b) > 0 ***REMOVED***
		w.writeUnknownFields(b)
	***REMOVED***
	return w.writeExtensions(m)
***REMOVED***

func (w *textWriter) writeSingularValue(v protoreflect.Value, fd protoreflect.FieldDescriptor) error ***REMOVED***
	switch fd.Kind() ***REMOVED***
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		switch vf := v.Float(); ***REMOVED***
		case math.IsInf(vf, +1):
			w.Write(posInf)
		case math.IsInf(vf, -1):
			w.Write(negInf)
		case math.IsNaN(vf):
			w.Write(nan)
		default:
			fmt.Fprint(w, v.Interface())
		***REMOVED***
	case protoreflect.StringKind:
		// NOTE: This does not validate UTF-8 for historical reasons.
		w.writeQuotedString(string(v.String()))
	case protoreflect.BytesKind:
		w.writeQuotedString(string(v.Bytes()))
	case protoreflect.MessageKind, protoreflect.GroupKind:
		var bra, ket byte = '<', '>'
		if fd.Kind() == protoreflect.GroupKind ***REMOVED***
			bra, ket = '***REMOVED***', '***REMOVED***'
		***REMOVED***
		w.WriteByte(bra)
		if !w.compact ***REMOVED***
			w.WriteByte('\n')
		***REMOVED***
		w.indent++
		m := v.Message()
		if m2, ok := m.Interface().(encoding.TextMarshaler); ok ***REMOVED***
			b, err := m2.MarshalText()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			w.Write(b)
		***REMOVED*** else ***REMOVED***
			w.writeMessage(m)
		***REMOVED***
		w.indent--
		w.WriteByte(ket)
	case protoreflect.EnumKind:
		if ev := fd.Enum().Values().ByNumber(v.Enum()); ev != nil ***REMOVED***
			fmt.Fprint(w, ev.Name())
		***REMOVED*** else ***REMOVED***
			fmt.Fprint(w, v.Enum())
		***REMOVED***
	default:
		fmt.Fprint(w, v.Interface())
	***REMOVED***
	return nil
***REMOVED***

// writeQuotedString writes a quoted string in the protocol buffer text format.
func (w *textWriter) writeQuotedString(s string) ***REMOVED***
	w.WriteByte('"')
	for i := 0; i < len(s); i++ ***REMOVED***
		switch c := s[i]; c ***REMOVED***
		case '\n':
			w.buf = append(w.buf, `\n`...)
		case '\r':
			w.buf = append(w.buf, `\r`...)
		case '\t':
			w.buf = append(w.buf, `\t`...)
		case '"':
			w.buf = append(w.buf, `\"`...)
		case '\\':
			w.buf = append(w.buf, `\\`...)
		default:
			if isPrint := c >= 0x20 && c < 0x7f; isPrint ***REMOVED***
				w.buf = append(w.buf, c)
			***REMOVED*** else ***REMOVED***
				w.buf = append(w.buf, fmt.Sprintf(`\%03o`, c)...)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	w.WriteByte('"')
***REMOVED***

func (w *textWriter) writeUnknownFields(b []byte) ***REMOVED***
	if !w.compact ***REMOVED***
		fmt.Fprintf(w, "/* %d unknown bytes */\n", len(b))
	***REMOVED***

	for len(b) > 0 ***REMOVED***
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 ***REMOVED***
			return
		***REMOVED***
		b = b[n:]

		if wtyp == protowire.EndGroupType ***REMOVED***
			w.indent--
			w.Write(endBraceNewline)
			continue
		***REMOVED***
		fmt.Fprint(w, num)
		if wtyp != protowire.StartGroupType ***REMOVED***
			w.WriteByte(':')
		***REMOVED***
		if !w.compact || wtyp == protowire.StartGroupType ***REMOVED***
			w.WriteByte(' ')
		***REMOVED***
		switch wtyp ***REMOVED***
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(b)
			if n < 0 ***REMOVED***
				return
			***REMOVED***
			b = b[n:]
			fmt.Fprint(w, v)
		case protowire.Fixed32Type:
			v, n := protowire.ConsumeFixed32(b)
			if n < 0 ***REMOVED***
				return
			***REMOVED***
			b = b[n:]
			fmt.Fprint(w, v)
		case protowire.Fixed64Type:
			v, n := protowire.ConsumeFixed64(b)
			if n < 0 ***REMOVED***
				return
			***REMOVED***
			b = b[n:]
			fmt.Fprint(w, v)
		case protowire.BytesType:
			v, n := protowire.ConsumeBytes(b)
			if n < 0 ***REMOVED***
				return
			***REMOVED***
			b = b[n:]
			fmt.Fprintf(w, "%q", v)
		case protowire.StartGroupType:
			w.WriteByte('***REMOVED***')
			w.indent++
		default:
			fmt.Fprintf(w, "/* unknown wire type %d */", wtyp)
		***REMOVED***
		w.WriteByte('\n')
	***REMOVED***
***REMOVED***

// writeExtensions writes all the extensions in m.
func (w *textWriter) writeExtensions(m protoreflect.Message) error ***REMOVED***
	md := m.Descriptor()
	if md.ExtensionRanges().Len() == 0 ***REMOVED***
		return nil
	***REMOVED***

	type ext struct ***REMOVED***
		desc protoreflect.FieldDescriptor
		val  protoreflect.Value
	***REMOVED***
	var exts []ext
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		if fd.IsExtension() ***REMOVED***
			exts = append(exts, ext***REMOVED***fd, v***REMOVED***)
		***REMOVED***
		return true
	***REMOVED***)
	sort.Slice(exts, func(i, j int) bool ***REMOVED***
		return exts[i].desc.Number() < exts[j].desc.Number()
	***REMOVED***)

	for _, ext := range exts ***REMOVED***
		// For message set, use the name of the message as the extension name.
		name := string(ext.desc.FullName())
		if isMessageSet(ext.desc.ContainingMessage()) ***REMOVED***
			name = strings.TrimSuffix(name, ".message_set_extension")
		***REMOVED***

		if !ext.desc.IsList() ***REMOVED***
			if err := w.writeSingularExtension(name, ext.val, ext.desc); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			lv := ext.val.List()
			for i := 0; i < lv.Len(); i++ ***REMOVED***
				if err := w.writeSingularExtension(name, lv.Get(i), ext.desc); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (w *textWriter) writeSingularExtension(name string, v protoreflect.Value, fd protoreflect.FieldDescriptor) error ***REMOVED***
	fmt.Fprintf(w, "[%s]:", name)
	if !w.compact ***REMOVED***
		w.WriteByte(' ')
	***REMOVED***
	if err := w.writeSingularValue(v, fd); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteByte('\n')
	return nil
***REMOVED***

func (w *textWriter) writeIndent() ***REMOVED***
	if !w.complete ***REMOVED***
		return
	***REMOVED***
	for i := 0; i < w.indent*2; i++ ***REMOVED***
		w.buf = append(w.buf, ' ')
	***REMOVED***
	w.complete = false
***REMOVED***
