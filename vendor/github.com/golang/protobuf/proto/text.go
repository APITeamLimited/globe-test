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

// Functions for writing the text protocol buffer format.

import (
	"bufio"
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
	"sort"
	"strings"
)

var (
	newline         = []byte("\n")
	spaces          = []byte("                                        ")
	gtNewline       = []byte(">\n")
	endBraceNewline = []byte("***REMOVED***\n")
	backslashN      = []byte***REMOVED***'\\', 'n'***REMOVED***
	backslashR      = []byte***REMOVED***'\\', 'r'***REMOVED***
	backslashT      = []byte***REMOVED***'\\', 't'***REMOVED***
	backslashDQ     = []byte***REMOVED***'\\', '"'***REMOVED***
	backslashBS     = []byte***REMOVED***'\\', '\\'***REMOVED***
	posInf          = []byte("inf")
	negInf          = []byte("-inf")
	nan             = []byte("nan")
)

type writer interface ***REMOVED***
	io.Writer
	WriteByte(byte) error
***REMOVED***

// textWriter is an io.Writer that tracks its indentation level.
type textWriter struct ***REMOVED***
	ind      int
	complete bool // if the current position is a complete line
	compact  bool // whether to write out as a one-liner
	w        writer
***REMOVED***

func (w *textWriter) WriteString(s string) (n int, err error) ***REMOVED***
	if !strings.Contains(s, "\n") ***REMOVED***
		if !w.compact && w.complete ***REMOVED***
			w.writeIndent()
		***REMOVED***
		w.complete = false
		return io.WriteString(w.w, s)
	***REMOVED***
	// WriteString is typically called without newlines, so this
	// codepath and its copy are rare.  We copy to avoid
	// duplicating all of Write's logic here.
	return w.Write([]byte(s))
***REMOVED***

func (w *textWriter) Write(p []byte) (n int, err error) ***REMOVED***
	newlines := bytes.Count(p, newline)
	if newlines == 0 ***REMOVED***
		if !w.compact && w.complete ***REMOVED***
			w.writeIndent()
		***REMOVED***
		n, err = w.w.Write(p)
		w.complete = false
		return n, err
	***REMOVED***

	frags := bytes.SplitN(p, newline, newlines+1)
	if w.compact ***REMOVED***
		for i, frag := range frags ***REMOVED***
			if i > 0 ***REMOVED***
				if err := w.w.WriteByte(' '); err != nil ***REMOVED***
					return n, err
				***REMOVED***
				n++
			***REMOVED***
			nn, err := w.w.Write(frag)
			n += nn
			if err != nil ***REMOVED***
				return n, err
			***REMOVED***
		***REMOVED***
		return n, nil
	***REMOVED***

	for i, frag := range frags ***REMOVED***
		if w.complete ***REMOVED***
			w.writeIndent()
		***REMOVED***
		nn, err := w.w.Write(frag)
		n += nn
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		if i+1 < len(frags) ***REMOVED***
			if err := w.w.WriteByte('\n'); err != nil ***REMOVED***
				return n, err
			***REMOVED***
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
	err := w.w.WriteByte(c)
	w.complete = c == '\n'
	return err
***REMOVED***

func (w *textWriter) indent() ***REMOVED*** w.ind++ ***REMOVED***

func (w *textWriter) unindent() ***REMOVED***
	if w.ind == 0 ***REMOVED***
		log.Print("proto: textWriter unindented too far")
		return
	***REMOVED***
	w.ind--
***REMOVED***

func writeName(w *textWriter, props *Properties) error ***REMOVED***
	if _, err := w.WriteString(props.OrigName); err != nil ***REMOVED***
		return err
	***REMOVED***
	if props.Wire != "group" ***REMOVED***
		return w.WriteByte(':')
	***REMOVED***
	return nil
***REMOVED***

// raw is the interface satisfied by RawMessage.
type raw interface ***REMOVED***
	Bytes() []byte
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

// isAny reports whether sv is a google.protobuf.Any message
func isAny(sv reflect.Value) bool ***REMOVED***
	type wkt interface ***REMOVED***
		XXX_WellKnownType() string
	***REMOVED***
	t, ok := sv.Addr().Interface().(wkt)
	return ok && t.XXX_WellKnownType() == "Any"
***REMOVED***

// writeProto3Any writes an expanded google.protobuf.Any message.
//
// It returns (false, nil) if sv value can't be unmarshaled (e.g. because
// required messages are not linked in).
//
// It returns (true, error) when sv was written in expanded format or an error
// was encountered.
func (tm *TextMarshaler) writeProto3Any(w *textWriter, sv reflect.Value) (bool, error) ***REMOVED***
	turl := sv.FieldByName("TypeUrl")
	val := sv.FieldByName("Value")
	if !turl.IsValid() || !val.IsValid() ***REMOVED***
		return true, errors.New("proto: invalid google.protobuf.Any message")
	***REMOVED***

	b, ok := val.Interface().([]byte)
	if !ok ***REMOVED***
		return true, errors.New("proto: invalid google.protobuf.Any message")
	***REMOVED***

	parts := strings.Split(turl.String(), "/")
	mt := MessageType(parts[len(parts)-1])
	if mt == nil ***REMOVED***
		return false, nil
	***REMOVED***
	m := reflect.New(mt.Elem())
	if err := Unmarshal(b, m.Interface().(Message)); err != nil ***REMOVED***
		return false, nil
	***REMOVED***
	w.Write([]byte("["))
	u := turl.String()
	if requiresQuotes(u) ***REMOVED***
		writeString(w, u)
	***REMOVED*** else ***REMOVED***
		w.Write([]byte(u))
	***REMOVED***
	if w.compact ***REMOVED***
		w.Write([]byte("]:<"))
	***REMOVED*** else ***REMOVED***
		w.Write([]byte("]: <\n"))
		w.ind++
	***REMOVED***
	if err := tm.writeStruct(w, m.Elem()); err != nil ***REMOVED***
		return true, err
	***REMOVED***
	if w.compact ***REMOVED***
		w.Write([]byte("> "))
	***REMOVED*** else ***REMOVED***
		w.ind--
		w.Write([]byte(">\n"))
	***REMOVED***
	return true, nil
***REMOVED***

func (tm *TextMarshaler) writeStruct(w *textWriter, sv reflect.Value) error ***REMOVED***
	if tm.ExpandAny && isAny(sv) ***REMOVED***
		if canExpand, err := tm.writeProto3Any(w, sv); canExpand ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	st := sv.Type()
	sprops := GetProperties(st)
	for i := 0; i < sv.NumField(); i++ ***REMOVED***
		fv := sv.Field(i)
		props := sprops.Prop[i]
		name := st.Field(i).Name

		if strings.HasPrefix(name, "XXX_") ***REMOVED***
			// There are two XXX_ fields:
			//   XXX_unrecognized []byte
			//   XXX_extensions   map[int32]proto.Extension
			// The first is handled here;
			// the second is handled at the bottom of this function.
			if name == "XXX_unrecognized" && !fv.IsNil() ***REMOVED***
				if err := writeUnknownStruct(w, fv.Interface().([]byte)); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		if fv.Kind() == reflect.Ptr && fv.IsNil() ***REMOVED***
			// Field not filled in. This could be an optional field or
			// a required field that wasn't filled in. Either way, there
			// isn't anything we can show for it.
			continue
		***REMOVED***
		if fv.Kind() == reflect.Slice && fv.IsNil() ***REMOVED***
			// Repeated field that is empty, or a bytes field that is unused.
			continue
		***REMOVED***

		if props.Repeated && fv.Kind() == reflect.Slice ***REMOVED***
			// Repeated field.
			for j := 0; j < fv.Len(); j++ ***REMOVED***
				if err := writeName(w, props); err != nil ***REMOVED***
					return err
				***REMOVED***
				if !w.compact ***REMOVED***
					if err := w.WriteByte(' '); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				v := fv.Index(j)
				if v.Kind() == reflect.Ptr && v.IsNil() ***REMOVED***
					// A nil message in a repeated field is not valid,
					// but we can handle that more gracefully than panicking.
					if _, err := w.Write([]byte("<nil>\n")); err != nil ***REMOVED***
						return err
					***REMOVED***
					continue
				***REMOVED***
				if err := tm.writeAny(w, v, props); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := w.WriteByte('\n'); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		if fv.Kind() == reflect.Map ***REMOVED***
			// Map fields are rendered as a repeated struct with key/value fields.
			keys := fv.MapKeys()
			sort.Sort(mapKeys(keys))
			for _, key := range keys ***REMOVED***
				val := fv.MapIndex(key)
				if err := writeName(w, props); err != nil ***REMOVED***
					return err
				***REMOVED***
				if !w.compact ***REMOVED***
					if err := w.WriteByte(' '); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				// open struct
				if err := w.WriteByte('<'); err != nil ***REMOVED***
					return err
				***REMOVED***
				if !w.compact ***REMOVED***
					if err := w.WriteByte('\n'); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				w.indent()
				// key
				if _, err := w.WriteString("key:"); err != nil ***REMOVED***
					return err
				***REMOVED***
				if !w.compact ***REMOVED***
					if err := w.WriteByte(' '); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				if err := tm.writeAny(w, key, props.mkeyprop); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := w.WriteByte('\n'); err != nil ***REMOVED***
					return err
				***REMOVED***
				// nil values aren't legal, but we can avoid panicking because of them.
				if val.Kind() != reflect.Ptr || !val.IsNil() ***REMOVED***
					// value
					if _, err := w.WriteString("value:"); err != nil ***REMOVED***
						return err
					***REMOVED***
					if !w.compact ***REMOVED***
						if err := w.WriteByte(' '); err != nil ***REMOVED***
							return err
						***REMOVED***
					***REMOVED***
					if err := tm.writeAny(w, val, props.mvalprop); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := w.WriteByte('\n'); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				// close struct
				w.unindent()
				if err := w.WriteByte('>'); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := w.WriteByte('\n'); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		if props.proto3 && fv.Kind() == reflect.Slice && fv.Len() == 0 ***REMOVED***
			// empty bytes field
			continue
		***REMOVED***
		if fv.Kind() != reflect.Ptr && fv.Kind() != reflect.Slice ***REMOVED***
			// proto3 non-repeated scalar field; skip if zero value
			if isProto3Zero(fv) ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		if fv.Kind() == reflect.Interface ***REMOVED***
			// Check if it is a oneof.
			if st.Field(i).Tag.Get("protobuf_oneof") != "" ***REMOVED***
				// fv is nil, or holds a pointer to generated struct.
				// That generated struct has exactly one field,
				// which has a protobuf struct tag.
				if fv.IsNil() ***REMOVED***
					continue
				***REMOVED***
				inner := fv.Elem().Elem() // interface -> *T -> T
				tag := inner.Type().Field(0).Tag.Get("protobuf")
				props = new(Properties) // Overwrite the outer props var, but not its pointee.
				props.Parse(tag)
				// Write the value in the oneof, not the oneof itself.
				fv = inner.Field(0)

				// Special case to cope with malformed messages gracefully:
				// If the value in the oneof is a nil pointer, don't panic
				// in writeAny.
				if fv.Kind() == reflect.Ptr && fv.IsNil() ***REMOVED***
					// Use errors.New so writeAny won't render quotes.
					msg := errors.New("/* nil */")
					fv = reflect.ValueOf(&msg).Elem()
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err := writeName(w, props); err != nil ***REMOVED***
			return err
		***REMOVED***
		if !w.compact ***REMOVED***
			if err := w.WriteByte(' '); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if b, ok := fv.Interface().(raw); ok ***REMOVED***
			if err := writeRaw(w, b.Bytes()); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		// Enums have a String method, so writeAny will work fine.
		if err := tm.writeAny(w, fv, props); err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := w.WriteByte('\n'); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Extensions (the XXX_extensions field).
	pv := sv.Addr()
	if _, ok := extendable(pv.Interface()); ok ***REMOVED***
		if err := tm.writeExtensions(w, pv); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// writeRaw writes an uninterpreted raw message.
func writeRaw(w *textWriter, b []byte) error ***REMOVED***
	if err := w.WriteByte('<'); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !w.compact ***REMOVED***
		if err := w.WriteByte('\n'); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	w.indent()
	if err := writeUnknownStruct(w, b); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.unindent()
	if err := w.WriteByte('>'); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// writeAny writes an arbitrary field.
func (tm *TextMarshaler) writeAny(w *textWriter, v reflect.Value, props *Properties) error ***REMOVED***
	v = reflect.Indirect(v)

	// Floats have special cases.
	if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 ***REMOVED***
		x := v.Float()
		var b []byte
		switch ***REMOVED***
		case math.IsInf(x, 1):
			b = posInf
		case math.IsInf(x, -1):
			b = negInf
		case math.IsNaN(x):
			b = nan
		***REMOVED***
		if b != nil ***REMOVED***
			_, err := w.Write(b)
			return err
		***REMOVED***
		// Other values are handled below.
	***REMOVED***

	// We don't attempt to serialise every possible value type; only those
	// that can occur in protocol buffers.
	switch v.Kind() ***REMOVED***
	case reflect.Slice:
		// Should only be a []byte; repeated fields are handled in writeStruct.
		if err := writeString(w, string(v.Bytes())); err != nil ***REMOVED***
			return err
		***REMOVED***
	case reflect.String:
		if err := writeString(w, v.String()); err != nil ***REMOVED***
			return err
		***REMOVED***
	case reflect.Struct:
		// Required/optional group/message.
		var bra, ket byte = '<', '>'
		if props != nil && props.Wire == "group" ***REMOVED***
			bra, ket = '***REMOVED***', '***REMOVED***'
		***REMOVED***
		if err := w.WriteByte(bra); err != nil ***REMOVED***
			return err
		***REMOVED***
		if !w.compact ***REMOVED***
			if err := w.WriteByte('\n'); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		w.indent()
		if etm, ok := v.Interface().(encoding.TextMarshaler); ok ***REMOVED***
			text, err := etm.MarshalText()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if _, err = w.Write(text); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else if err := tm.writeStruct(w, v); err != nil ***REMOVED***
			return err
		***REMOVED***
		w.unindent()
		if err := w.WriteByte(ket); err != nil ***REMOVED***
			return err
		***REMOVED***
	default:
		_, err := fmt.Fprint(w, v.Interface())
		return err
	***REMOVED***
	return nil
***REMOVED***

// equivalent to C's isprint.
func isprint(c byte) bool ***REMOVED***
	return c >= 0x20 && c < 0x7f
***REMOVED***

// writeString writes a string in the protocol buffer text format.
// It is similar to strconv.Quote except we don't use Go escape sequences,
// we treat the string as a byte sequence, and we use octal escapes.
// These differences are to maintain interoperability with the other
// languages' implementations of the text format.
func writeString(w *textWriter, s string) error ***REMOVED***
	// use WriteByte here to get any needed indent
	if err := w.WriteByte('"'); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Loop over the bytes, not the runes.
	for i := 0; i < len(s); i++ ***REMOVED***
		var err error
		// Divergence from C++: we don't escape apostrophes.
		// There's no need to escape them, and the C++ parser
		// copes with a naked apostrophe.
		switch c := s[i]; c ***REMOVED***
		case '\n':
			_, err = w.w.Write(backslashN)
		case '\r':
			_, err = w.w.Write(backslashR)
		case '\t':
			_, err = w.w.Write(backslashT)
		case '"':
			_, err = w.w.Write(backslashDQ)
		case '\\':
			_, err = w.w.Write(backslashBS)
		default:
			if isprint(c) ***REMOVED***
				err = w.w.WriteByte(c)
			***REMOVED*** else ***REMOVED***
				_, err = fmt.Fprintf(w.w, "\\%03o", c)
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return w.WriteByte('"')
***REMOVED***

func writeUnknownStruct(w *textWriter, data []byte) (err error) ***REMOVED***
	if !w.compact ***REMOVED***
		if _, err := fmt.Fprintf(w, "/* %d unknown bytes */\n", len(data)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	b := NewBuffer(data)
	for b.index < len(b.buf) ***REMOVED***
		x, err := b.DecodeVarint()
		if err != nil ***REMOVED***
			_, err := fmt.Fprintf(w, "/* %v */\n", err)
			return err
		***REMOVED***
		wire, tag := x&7, x>>3
		if wire == WireEndGroup ***REMOVED***
			w.unindent()
			if _, err := w.Write(endBraceNewline); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***
		if _, err := fmt.Fprint(w, tag); err != nil ***REMOVED***
			return err
		***REMOVED***
		if wire != WireStartGroup ***REMOVED***
			if err := w.WriteByte(':'); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if !w.compact || wire == WireStartGroup ***REMOVED***
			if err := w.WriteByte(' '); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		switch wire ***REMOVED***
		case WireBytes:
			buf, e := b.DecodeRawBytes(false)
			if e == nil ***REMOVED***
				_, err = fmt.Fprintf(w, "%q", buf)
			***REMOVED*** else ***REMOVED***
				_, err = fmt.Fprintf(w, "/* %v */", e)
			***REMOVED***
		case WireFixed32:
			x, err = b.DecodeFixed32()
			err = writeUnknownInt(w, x, err)
		case WireFixed64:
			x, err = b.DecodeFixed64()
			err = writeUnknownInt(w, x, err)
		case WireStartGroup:
			err = w.WriteByte('***REMOVED***')
			w.indent()
		case WireVarint:
			x, err = b.DecodeVarint()
			err = writeUnknownInt(w, x, err)
		default:
			_, err = fmt.Fprintf(w, "/* unknown wire type %d */", wire)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err = w.WriteByte('\n'); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func writeUnknownInt(w *textWriter, x uint64, err error) error ***REMOVED***
	if err == nil ***REMOVED***
		_, err = fmt.Fprint(w, x)
	***REMOVED*** else ***REMOVED***
		_, err = fmt.Fprintf(w, "/* %v */", err)
	***REMOVED***
	return err
***REMOVED***

type int32Slice []int32

func (s int32Slice) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s int32Slice) Less(i, j int) bool ***REMOVED*** return s[i] < s[j] ***REMOVED***
func (s int32Slice) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***

// writeExtensions writes all the extensions in pv.
// pv is assumed to be a pointer to a protocol message struct that is extendable.
func (tm *TextMarshaler) writeExtensions(w *textWriter, pv reflect.Value) error ***REMOVED***
	emap := extensionMaps[pv.Type().Elem()]
	ep, _ := extendable(pv.Interface())

	// Order the extensions by ID.
	// This isn't strictly necessary, but it will give us
	// canonical output, which will also make testing easier.
	m, mu := ep.extensionsRead()
	if m == nil ***REMOVED***
		return nil
	***REMOVED***
	mu.Lock()
	ids := make([]int32, 0, len(m))
	for id := range m ***REMOVED***
		ids = append(ids, id)
	***REMOVED***
	sort.Sort(int32Slice(ids))
	mu.Unlock()

	for _, extNum := range ids ***REMOVED***
		ext := m[extNum]
		var desc *ExtensionDesc
		if emap != nil ***REMOVED***
			desc = emap[extNum]
		***REMOVED***
		if desc == nil ***REMOVED***
			// Unknown extension.
			if err := writeUnknownStruct(w, ext.enc); err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		pb, err := GetExtension(ep, desc)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed getting extension: %v", err)
		***REMOVED***

		// Repeated extensions will appear as a slice.
		if !desc.repeated() ***REMOVED***
			if err := tm.writeExtension(w, desc.Name, pb); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			v := reflect.ValueOf(pb)
			for i := 0; i < v.Len(); i++ ***REMOVED***
				if err := tm.writeExtension(w, desc.Name, v.Index(i).Interface()); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (tm *TextMarshaler) writeExtension(w *textWriter, name string, pb interface***REMOVED******REMOVED***) error ***REMOVED***
	if _, err := fmt.Fprintf(w, "[%s]:", name); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !w.compact ***REMOVED***
		if err := w.WriteByte(' '); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := tm.writeAny(w, reflect.ValueOf(pb), nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := w.WriteByte('\n'); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (w *textWriter) writeIndent() ***REMOVED***
	if !w.complete ***REMOVED***
		return
	***REMOVED***
	remain := w.ind * 2
	for remain > 0 ***REMOVED***
		n := remain
		if n > len(spaces) ***REMOVED***
			n = len(spaces)
		***REMOVED***
		w.w.Write(spaces[:n])
		remain -= n
	***REMOVED***
	w.complete = false
***REMOVED***

// TextMarshaler is a configurable text format marshaler.
type TextMarshaler struct ***REMOVED***
	Compact   bool // use compact text format (one line).
	ExpandAny bool // expand google.protobuf.Any messages of known types
***REMOVED***

// Marshal writes a given protocol buffer in text format.
// The only errors returned are from w.
func (tm *TextMarshaler) Marshal(w io.Writer, pb Message) error ***REMOVED***
	val := reflect.ValueOf(pb)
	if pb == nil || val.IsNil() ***REMOVED***
		w.Write([]byte("<nil>"))
		return nil
	***REMOVED***
	var bw *bufio.Writer
	ww, ok := w.(writer)
	if !ok ***REMOVED***
		bw = bufio.NewWriter(w)
		ww = bw
	***REMOVED***
	aw := &textWriter***REMOVED***
		w:        ww,
		complete: true,
		compact:  tm.Compact,
	***REMOVED***

	if etm, ok := pb.(encoding.TextMarshaler); ok ***REMOVED***
		text, err := etm.MarshalText()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err = aw.Write(text); err != nil ***REMOVED***
			return err
		***REMOVED***
		if bw != nil ***REMOVED***
			return bw.Flush()
		***REMOVED***
		return nil
	***REMOVED***
	// Dereference the received pointer so we don't have outer < and >.
	v := reflect.Indirect(val)
	if err := tm.writeStruct(aw, v); err != nil ***REMOVED***
		return err
	***REMOVED***
	if bw != nil ***REMOVED***
		return bw.Flush()
	***REMOVED***
	return nil
***REMOVED***

// Text is the same as Marshal, but returns the string directly.
func (tm *TextMarshaler) Text(pb Message) string ***REMOVED***
	var buf bytes.Buffer
	tm.Marshal(&buf, pb)
	return buf.String()
***REMOVED***

var (
	defaultTextMarshaler = TextMarshaler***REMOVED******REMOVED***
	compactTextMarshaler = TextMarshaler***REMOVED***Compact: true***REMOVED***
)

// TODO: consider removing some of the Marshal functions below.

// MarshalText writes a given protocol buffer in text format.
// The only errors returned are from w.
func MarshalText(w io.Writer, pb Message) error ***REMOVED*** return defaultTextMarshaler.Marshal(w, pb) ***REMOVED***

// MarshalTextString is the same as MarshalText, but returns the string directly.
func MarshalTextString(pb Message) string ***REMOVED*** return defaultTextMarshaler.Text(pb) ***REMOVED***

// CompactText writes a given protocol buffer in compact text format (one line).
func CompactText(w io.Writer, pb Message) error ***REMOVED*** return compactTextMarshaler.Marshal(w, pb) ***REMOVED***

// CompactTextString is the same as CompactText, but returns the string directly.
func CompactTextString(pb Message) string ***REMOVED*** return compactTextMarshaler.Text(pb) ***REMOVED***
