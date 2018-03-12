// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// This file contains utilities for generating code.

// TODO: other write methods like:
// - slices, maps, types, etc.

// CodeWriter is a utility for writing structured code. It computes the content
// hash and size of written content. It ensures there are newlines between
// written code blocks.
type CodeWriter struct ***REMOVED***
	buf  bytes.Buffer
	Size int
	Hash hash.Hash32 // content hash
	gob  *gob.Encoder
	// For comments we skip the usual one-line separator if they are followed by
	// a code block.
	skipSep bool
***REMOVED***

func (w *CodeWriter) Write(p []byte) (n int, err error) ***REMOVED***
	return w.buf.Write(p)
***REMOVED***

// NewCodeWriter returns a new CodeWriter.
func NewCodeWriter() *CodeWriter ***REMOVED***
	h := fnv.New32()
	return &CodeWriter***REMOVED***Hash: h, gob: gob.NewEncoder(h)***REMOVED***
***REMOVED***

// WriteGoFile appends the buffer with the total size of all created structures
// and writes it as a Go file to the the given file with the given package name.
func (w *CodeWriter) WriteGoFile(filename, pkg string) ***REMOVED***
	f, err := os.Create(filename)
	if err != nil ***REMOVED***
		log.Fatalf("Could not create file %s: %v", filename, err)
	***REMOVED***
	defer f.Close()
	if _, err = w.WriteGo(f, pkg); err != nil ***REMOVED***
		log.Fatalf("Error writing file %s: %v", filename, err)
	***REMOVED***
***REMOVED***

// WriteGo appends the buffer with the total size of all created structures and
// writes it as a Go file to the the given writer with the given package name.
func (w *CodeWriter) WriteGo(out io.Writer, pkg string) (n int, err error) ***REMOVED***
	sz := w.Size
	w.WriteComment("Total table size %d bytes (%dKiB); checksum: %X\n", sz, sz/1024, w.Hash.Sum32())
	defer w.buf.Reset()
	return WriteGo(out, pkg, w.buf.Bytes())
***REMOVED***

func (w *CodeWriter) printf(f string, x ...interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Fprintf(w, f, x...)
***REMOVED***

func (w *CodeWriter) insertSep() ***REMOVED***
	if w.skipSep ***REMOVED***
		w.skipSep = false
		return
	***REMOVED***
	// Use at least two newlines to ensure a blank space between the previous
	// block. WriteGoFile will remove extraneous newlines.
	w.printf("\n\n")
***REMOVED***

// WriteComment writes a comment block. All line starts are prefixed with "//".
// Initial empty lines are gobbled. The indentation for the first line is
// stripped from consecutive lines.
func (w *CodeWriter) WriteComment(comment string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	s := fmt.Sprintf(comment, args...)
	s = strings.Trim(s, "\n")

	// Use at least two newlines to ensure a blank space between the previous
	// block. WriteGoFile will remove extraneous newlines.
	w.printf("\n\n// ")
	w.skipSep = true

	// strip first indent level.
	sep := "\n"
	for ; len(s) > 0 && (s[0] == '\t' || s[0] == ' '); s = s[1:] ***REMOVED***
		sep += s[:1]
	***REMOVED***

	strings.NewReplacer(sep, "\n// ", "\n", "\n// ").WriteString(w, s)

	w.printf("\n")
***REMOVED***

func (w *CodeWriter) writeSizeInfo(size int) ***REMOVED***
	w.printf("// Size: %d bytes\n", size)
***REMOVED***

// WriteConst writes a constant of the given name and value.
func (w *CodeWriter) WriteConst(name string, x interface***REMOVED******REMOVED***) ***REMOVED***
	w.insertSep()
	v := reflect.ValueOf(x)

	switch v.Type().Kind() ***REMOVED***
	case reflect.String:
		w.printf("const %s %s = ", name, typeName(x))
		w.WriteString(v.String())
		w.printf("\n")
	default:
		w.printf("const %s = %#v\n", name, x)
	***REMOVED***
***REMOVED***

// WriteVar writes a variable of the given name and value.
func (w *CodeWriter) WriteVar(name string, x interface***REMOVED******REMOVED***) ***REMOVED***
	w.insertSep()
	v := reflect.ValueOf(x)
	oldSize := w.Size
	sz := int(v.Type().Size())
	w.Size += sz

	switch v.Type().Kind() ***REMOVED***
	case reflect.String:
		w.printf("var %s %s = ", name, typeName(x))
		w.WriteString(v.String())
	case reflect.Struct:
		w.gob.Encode(x)
		fallthrough
	case reflect.Slice, reflect.Array:
		w.printf("var %s = ", name)
		w.writeValue(v)
		w.writeSizeInfo(w.Size - oldSize)
	default:
		w.printf("var %s %s = ", name, typeName(x))
		w.gob.Encode(x)
		w.writeValue(v)
		w.writeSizeInfo(w.Size - oldSize)
	***REMOVED***
	w.printf("\n")
***REMOVED***

func (w *CodeWriter) writeValue(v reflect.Value) ***REMOVED***
	x := v.Interface()
	switch v.Kind() ***REMOVED***
	case reflect.String:
		w.WriteString(v.String())
	case reflect.Array:
		// Don't double count: callers of WriteArray count on the size being
		// added, so we need to discount it here.
		w.Size -= int(v.Type().Size())
		w.writeSlice(x, true)
	case reflect.Slice:
		w.writeSlice(x, false)
	case reflect.Struct:
		w.printf("%s***REMOVED***\n", typeName(v.Interface()))
		t := v.Type()
		for i := 0; i < v.NumField(); i++ ***REMOVED***
			w.printf("%s: ", t.Field(i).Name)
			w.writeValue(v.Field(i))
			w.printf(",\n")
		***REMOVED***
		w.printf("***REMOVED***")
	default:
		w.printf("%#v", x)
	***REMOVED***
***REMOVED***

// WriteString writes a string literal.
func (w *CodeWriter) WriteString(s string) ***REMOVED***
	s = strings.Replace(s, `\`, `\\`, -1)
	io.WriteString(w.Hash, s) // content hash
	w.Size += len(s)

	const maxInline = 40
	if len(s) <= maxInline ***REMOVED***
		w.printf("%q", s)
		return
	***REMOVED***

	// We will render the string as a multi-line string.
	const maxWidth = 80 - 4 - len(`"`) - len(`" +`)

	// When starting on its own line, go fmt indents line 2+ an extra level.
	n, max := maxWidth, maxWidth-4

	// As per https://golang.org/issue/18078, the compiler has trouble
	// compiling the concatenation of many strings, s0 + s1 + s2 + ... + sN,
	// for large N. We insert redundant, explicit parentheses to work around
	// that, lowering the N at any given step: (s0 + s1 + ... + s63) + (s64 +
	// ... + s127) + etc + (etc + ... + sN).
	explicitParens, extraComment := len(s) > 128*1024, ""
	if explicitParens ***REMOVED***
		w.printf(`(`)
		extraComment = "; the redundant, explicit parens are for https://golang.org/issue/18078"
	***REMOVED***

	// Print "" +\n, if a string does not start on its own line.
	b := w.buf.Bytes()
	if p := len(bytes.TrimRight(b, " \t")); p > 0 && b[p-1] != '\n' ***REMOVED***
		w.printf("\"\" + // Size: %d bytes%s\n", len(s), extraComment)
		n, max = maxWidth, maxWidth
	***REMOVED***

	w.printf(`"`)

	for sz, p, nLines := 0, 0, 0; p < len(s); ***REMOVED***
		var r rune
		r, sz = utf8.DecodeRuneInString(s[p:])
		out := s[p : p+sz]
		chars := 1
		if !unicode.IsPrint(r) || r == utf8.RuneError || r == '"' ***REMOVED***
			switch sz ***REMOVED***
			case 1:
				out = fmt.Sprintf("\\x%02x", s[p])
			case 2, 3:
				out = fmt.Sprintf("\\u%04x", r)
			case 4:
				out = fmt.Sprintf("\\U%08x", r)
			***REMOVED***
			chars = len(out)
		***REMOVED***
		if n -= chars; n < 0 ***REMOVED***
			nLines++
			if explicitParens && nLines&63 == 63 ***REMOVED***
				w.printf("\") + (\"")
			***REMOVED***
			w.printf("\" +\n\"")
			n = max - len(out)
		***REMOVED***
		w.printf("%s", out)
		p += sz
	***REMOVED***
	w.printf(`"`)
	if explicitParens ***REMOVED***
		w.printf(`)`)
	***REMOVED***
***REMOVED***

// WriteSlice writes a slice value.
func (w *CodeWriter) WriteSlice(x interface***REMOVED******REMOVED***) ***REMOVED***
	w.writeSlice(x, false)
***REMOVED***

// WriteArray writes an array value.
func (w *CodeWriter) WriteArray(x interface***REMOVED******REMOVED***) ***REMOVED***
	w.writeSlice(x, true)
***REMOVED***

func (w *CodeWriter) writeSlice(x interface***REMOVED******REMOVED***, isArray bool) ***REMOVED***
	v := reflect.ValueOf(x)
	w.gob.Encode(v.Len())
	w.Size += v.Len() * int(v.Type().Elem().Size())
	name := typeName(x)
	if isArray ***REMOVED***
		name = fmt.Sprintf("[%d]%s", v.Len(), name[strings.Index(name, "]")+1:])
	***REMOVED***
	if isArray ***REMOVED***
		w.printf("%s***REMOVED***\n", name)
	***REMOVED*** else ***REMOVED***
		w.printf("%s***REMOVED*** // %d elements\n", name, v.Len())
	***REMOVED***

	switch kind := v.Type().Elem().Kind(); kind ***REMOVED***
	case reflect.String:
		for _, s := range x.([]string) ***REMOVED***
			w.WriteString(s)
			w.printf(",\n")
		***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// nLine and nBlock are the number of elements per line and block.
		nLine, nBlock, format := 8, 64, "%d,"
		switch kind ***REMOVED***
		case reflect.Uint8:
			format = "%#02x,"
		case reflect.Uint16:
			format = "%#04x,"
		case reflect.Uint32:
			nLine, nBlock, format = 4, 32, "%#08x,"
		case reflect.Uint, reflect.Uint64:
			nLine, nBlock, format = 4, 32, "%#016x,"
		case reflect.Int8:
			nLine = 16
		***REMOVED***
		n := nLine
		for i := 0; i < v.Len(); i++ ***REMOVED***
			if i%nBlock == 0 && v.Len() > nBlock ***REMOVED***
				w.printf("// Entry %X - %X\n", i, i+nBlock-1)
			***REMOVED***
			x := v.Index(i).Interface()
			w.gob.Encode(x)
			w.printf(format, x)
			if n--; n == 0 ***REMOVED***
				n = nLine
				w.printf("\n")
			***REMOVED***
		***REMOVED***
		w.printf("\n")
	case reflect.Struct:
		zero := reflect.Zero(v.Type().Elem()).Interface()
		for i := 0; i < v.Len(); i++ ***REMOVED***
			x := v.Index(i).Interface()
			w.gob.EncodeValue(v)
			if !reflect.DeepEqual(zero, x) ***REMOVED***
				line := fmt.Sprintf("%#v,\n", x)
				line = line[strings.IndexByte(line, '***REMOVED***'):]
				w.printf("%d: ", i)
				w.printf(line)
			***REMOVED***
		***REMOVED***
	case reflect.Array:
		for i := 0; i < v.Len(); i++ ***REMOVED***
			w.printf("%d: %#v,\n", i, v.Index(i).Interface())
		***REMOVED***
	default:
		panic("gen: slice elem type not supported")
	***REMOVED***
	w.printf("***REMOVED***")
***REMOVED***

// WriteType writes a definition of the type of the given value and returns the
// type name.
func (w *CodeWriter) WriteType(x interface***REMOVED******REMOVED***) string ***REMOVED***
	t := reflect.TypeOf(x)
	w.printf("type %s struct ***REMOVED***\n", t.Name())
	for i := 0; i < t.NumField(); i++ ***REMOVED***
		w.printf("\t%s %s\n", t.Field(i).Name, t.Field(i).Type)
	***REMOVED***
	w.printf("***REMOVED***\n")
	return t.Name()
***REMOVED***

// typeName returns the name of the go type of x.
func typeName(x interface***REMOVED******REMOVED***) string ***REMOVED***
	t := reflect.ValueOf(x).Type()
	return strings.Replace(fmt.Sprint(t), "main.", "", 1)
***REMOVED***
