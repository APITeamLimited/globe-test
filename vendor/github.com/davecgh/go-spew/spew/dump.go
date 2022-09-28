/*
 * Copyright (c) 2013-2016 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package spew

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	// uint8Type is a reflect.Type representing a uint8.  It is used to
	// convert cgo types to uint8 slices for hexdumping.
	uint8Type = reflect.TypeOf(uint8(0))

	// cCharRE is a regular expression that matches a cgo char.
	// It is used to detect character arrays to hexdump them.
	cCharRE = regexp.MustCompile(`^.*\._Ctype_char$`)

	// cUnsignedCharRE is a regular expression that matches a cgo unsigned
	// char.  It is used to detect unsigned character arrays to hexdump
	// them.
	cUnsignedCharRE = regexp.MustCompile(`^.*\._Ctype_unsignedchar$`)

	// cUint8tCharRE is a regular expression that matches a cgo uint8_t.
	// It is used to detect uint8_t arrays to hexdump them.
	cUint8tCharRE = regexp.MustCompile(`^.*\._Ctype_uint8_t$`)
)

// dumpState contains information about the state of a dump operation.
type dumpState struct ***REMOVED***
	w                io.Writer
	depth            int
	pointers         map[uintptr]int
	ignoreNextType   bool
	ignoreNextIndent bool
	cs               *ConfigState
***REMOVED***

// indent performs indentation according to the depth level and cs.Indent
// option.
func (d *dumpState) indent() ***REMOVED***
	if d.ignoreNextIndent ***REMOVED***
		d.ignoreNextIndent = false
		return
	***REMOVED***
	d.w.Write(bytes.Repeat([]byte(d.cs.Indent), d.depth))
***REMOVED***

// unpackValue returns values inside of non-nil interfaces when possible.
// This is useful for data types like structs, arrays, slices, and maps which
// can contain varying types packed inside an interface.
func (d *dumpState) unpackValue(v reflect.Value) reflect.Value ***REMOVED***
	if v.Kind() == reflect.Interface && !v.IsNil() ***REMOVED***
		v = v.Elem()
	***REMOVED***
	return v
***REMOVED***

// dumpPtr handles formatting of pointers by indirecting them as necessary.
func (d *dumpState) dumpPtr(v reflect.Value) ***REMOVED***
	// Remove pointers at or below the current depth from map used to detect
	// circular refs.
	for k, depth := range d.pointers ***REMOVED***
		if depth >= d.depth ***REMOVED***
			delete(d.pointers, k)
		***REMOVED***
	***REMOVED***

	// Keep list of all dereferenced pointers to show later.
	pointerChain := make([]uintptr, 0)

	// Figure out how many levels of indirection there are by dereferencing
	// pointers and unpacking interfaces down the chain while detecting circular
	// references.
	nilFound := false
	cycleFound := false
	indirects := 0
	ve := v
	for ve.Kind() == reflect.Ptr ***REMOVED***
		if ve.IsNil() ***REMOVED***
			nilFound = true
			break
		***REMOVED***
		indirects++
		addr := ve.Pointer()
		pointerChain = append(pointerChain, addr)
		if pd, ok := d.pointers[addr]; ok && pd < d.depth ***REMOVED***
			cycleFound = true
			indirects--
			break
		***REMOVED***
		d.pointers[addr] = d.depth

		ve = ve.Elem()
		if ve.Kind() == reflect.Interface ***REMOVED***
			if ve.IsNil() ***REMOVED***
				nilFound = true
				break
			***REMOVED***
			ve = ve.Elem()
		***REMOVED***
	***REMOVED***

	// Display type information.
	d.w.Write(openParenBytes)
	d.w.Write(bytes.Repeat(asteriskBytes, indirects))
	d.w.Write([]byte(ve.Type().String()))
	d.w.Write(closeParenBytes)

	// Display pointer information.
	if !d.cs.DisablePointerAddresses && len(pointerChain) > 0 ***REMOVED***
		d.w.Write(openParenBytes)
		for i, addr := range pointerChain ***REMOVED***
			if i > 0 ***REMOVED***
				d.w.Write(pointerChainBytes)
			***REMOVED***
			printHexPtr(d.w, addr)
		***REMOVED***
		d.w.Write(closeParenBytes)
	***REMOVED***

	// Display dereferenced value.
	d.w.Write(openParenBytes)
	switch ***REMOVED***
	case nilFound:
		d.w.Write(nilAngleBytes)

	case cycleFound:
		d.w.Write(circularBytes)

	default:
		d.ignoreNextType = true
		d.dump(ve)
	***REMOVED***
	d.w.Write(closeParenBytes)
***REMOVED***

// dumpSlice handles formatting of arrays and slices.  Byte (uint8 under
// reflection) arrays and slices are dumped in hexdump -C fashion.
func (d *dumpState) dumpSlice(v reflect.Value) ***REMOVED***
	// Determine whether this type should be hex dumped or not.  Also,
	// for types which should be hexdumped, try to use the underlying data
	// first, then fall back to trying to convert them to a uint8 slice.
	var buf []uint8
	doConvert := false
	doHexDump := false
	numEntries := v.Len()
	if numEntries > 0 ***REMOVED***
		vt := v.Index(0).Type()
		vts := vt.String()
		switch ***REMOVED***
		// C types that need to be converted.
		case cCharRE.MatchString(vts):
			fallthrough
		case cUnsignedCharRE.MatchString(vts):
			fallthrough
		case cUint8tCharRE.MatchString(vts):
			doConvert = true

		// Try to use existing uint8 slices and fall back to converting
		// and copying if that fails.
		case vt.Kind() == reflect.Uint8:
			// We need an addressable interface to convert the type
			// to a byte slice.  However, the reflect package won't
			// give us an interface on certain things like
			// unexported struct fields in order to enforce
			// visibility rules.  We use unsafe, when available, to
			// bypass these restrictions since this package does not
			// mutate the values.
			vs := v
			if !vs.CanInterface() || !vs.CanAddr() ***REMOVED***
				vs = unsafeReflectValue(vs)
			***REMOVED***
			if !UnsafeDisabled ***REMOVED***
				vs = vs.Slice(0, numEntries)

				// Use the existing uint8 slice if it can be
				// type asserted.
				iface := vs.Interface()
				if slice, ok := iface.([]uint8); ok ***REMOVED***
					buf = slice
					doHexDump = true
					break
				***REMOVED***
			***REMOVED***

			// The underlying data needs to be converted if it can't
			// be type asserted to a uint8 slice.
			doConvert = true
		***REMOVED***

		// Copy and convert the underlying type if needed.
		if doConvert && vt.ConvertibleTo(uint8Type) ***REMOVED***
			// Convert and copy each element into a uint8 byte
			// slice.
			buf = make([]uint8, numEntries)
			for i := 0; i < numEntries; i++ ***REMOVED***
				vv := v.Index(i)
				buf[i] = uint8(vv.Convert(uint8Type).Uint())
			***REMOVED***
			doHexDump = true
		***REMOVED***
	***REMOVED***

	// Hexdump the entire slice as needed.
	if doHexDump ***REMOVED***
		indent := strings.Repeat(d.cs.Indent, d.depth)
		str := indent + hex.Dump(buf)
		str = strings.Replace(str, "\n", "\n"+indent, -1)
		str = strings.TrimRight(str, d.cs.Indent)
		d.w.Write([]byte(str))
		return
	***REMOVED***

	// Recursively call dump for each item.
	for i := 0; i < numEntries; i++ ***REMOVED***
		d.dump(d.unpackValue(v.Index(i)))
		if i < (numEntries - 1) ***REMOVED***
			d.w.Write(commaNewlineBytes)
		***REMOVED*** else ***REMOVED***
			d.w.Write(newlineBytes)
		***REMOVED***
	***REMOVED***
***REMOVED***

// dump is the main workhorse for dumping a value.  It uses the passed reflect
// value to figure out what kind of object we are dealing with and formats it
// appropriately.  It is a recursive function, however circular data structures
// are detected and handled properly.
func (d *dumpState) dump(v reflect.Value) ***REMOVED***
	// Handle invalid reflect values immediately.
	kind := v.Kind()
	if kind == reflect.Invalid ***REMOVED***
		d.w.Write(invalidAngleBytes)
		return
	***REMOVED***

	// Handle pointers specially.
	if kind == reflect.Ptr ***REMOVED***
		d.indent()
		d.dumpPtr(v)
		return
	***REMOVED***

	// Print type information unless already handled elsewhere.
	if !d.ignoreNextType ***REMOVED***
		d.indent()
		d.w.Write(openParenBytes)
		d.w.Write([]byte(v.Type().String()))
		d.w.Write(closeParenBytes)
		d.w.Write(spaceBytes)
	***REMOVED***
	d.ignoreNextType = false

	// Display length and capacity if the built-in len and cap functions
	// work with the value's kind and the len/cap itself is non-zero.
	valueLen, valueCap := 0, 0
	switch v.Kind() ***REMOVED***
	case reflect.Array, reflect.Slice, reflect.Chan:
		valueLen, valueCap = v.Len(), v.Cap()
	case reflect.Map, reflect.String:
		valueLen = v.Len()
	***REMOVED***
	if valueLen != 0 || !d.cs.DisableCapacities && valueCap != 0 ***REMOVED***
		d.w.Write(openParenBytes)
		if valueLen != 0 ***REMOVED***
			d.w.Write(lenEqualsBytes)
			printInt(d.w, int64(valueLen), 10)
		***REMOVED***
		if !d.cs.DisableCapacities && valueCap != 0 ***REMOVED***
			if valueLen != 0 ***REMOVED***
				d.w.Write(spaceBytes)
			***REMOVED***
			d.w.Write(capEqualsBytes)
			printInt(d.w, int64(valueCap), 10)
		***REMOVED***
		d.w.Write(closeParenBytes)
		d.w.Write(spaceBytes)
	***REMOVED***

	// Call Stringer/error interfaces if they exist and the handle methods flag
	// is enabled
	if !d.cs.DisableMethods ***REMOVED***
		if (kind != reflect.Invalid) && (kind != reflect.Interface) ***REMOVED***
			if handled := handleMethods(d.cs, d.w, v); handled ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	switch kind ***REMOVED***
	case reflect.Invalid:
		// Do nothing.  We should never get here since invalid has already
		// been handled above.

	case reflect.Bool:
		printBool(d.w, v.Bool())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		printInt(d.w, v.Int(), 10)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		printUint(d.w, v.Uint(), 10)

	case reflect.Float32:
		printFloat(d.w, v.Float(), 32)

	case reflect.Float64:
		printFloat(d.w, v.Float(), 64)

	case reflect.Complex64:
		printComplex(d.w, v.Complex(), 32)

	case reflect.Complex128:
		printComplex(d.w, v.Complex(), 64)

	case reflect.Slice:
		if v.IsNil() ***REMOVED***
			d.w.Write(nilAngleBytes)
			break
		***REMOVED***
		fallthrough

	case reflect.Array:
		d.w.Write(openBraceNewlineBytes)
		d.depth++
		if (d.cs.MaxDepth != 0) && (d.depth > d.cs.MaxDepth) ***REMOVED***
			d.indent()
			d.w.Write(maxNewlineBytes)
		***REMOVED*** else ***REMOVED***
			d.dumpSlice(v)
		***REMOVED***
		d.depth--
		d.indent()
		d.w.Write(closeBraceBytes)

	case reflect.String:
		d.w.Write([]byte(strconv.Quote(v.String())))

	case reflect.Interface:
		// The only time we should get here is for nil interfaces due to
		// unpackValue calls.
		if v.IsNil() ***REMOVED***
			d.w.Write(nilAngleBytes)
		***REMOVED***

	case reflect.Ptr:
		// Do nothing.  We should never get here since pointers have already
		// been handled above.

	case reflect.Map:
		// nil maps should be indicated as different than empty maps
		if v.IsNil() ***REMOVED***
			d.w.Write(nilAngleBytes)
			break
		***REMOVED***

		d.w.Write(openBraceNewlineBytes)
		d.depth++
		if (d.cs.MaxDepth != 0) && (d.depth > d.cs.MaxDepth) ***REMOVED***
			d.indent()
			d.w.Write(maxNewlineBytes)
		***REMOVED*** else ***REMOVED***
			numEntries := v.Len()
			keys := v.MapKeys()
			if d.cs.SortKeys ***REMOVED***
				sortValues(keys, d.cs)
			***REMOVED***
			for i, key := range keys ***REMOVED***
				d.dump(d.unpackValue(key))
				d.w.Write(colonSpaceBytes)
				d.ignoreNextIndent = true
				d.dump(d.unpackValue(v.MapIndex(key)))
				if i < (numEntries - 1) ***REMOVED***
					d.w.Write(commaNewlineBytes)
				***REMOVED*** else ***REMOVED***
					d.w.Write(newlineBytes)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		d.depth--
		d.indent()
		d.w.Write(closeBraceBytes)

	case reflect.Struct:
		d.w.Write(openBraceNewlineBytes)
		d.depth++
		if (d.cs.MaxDepth != 0) && (d.depth > d.cs.MaxDepth) ***REMOVED***
			d.indent()
			d.w.Write(maxNewlineBytes)
		***REMOVED*** else ***REMOVED***
			vt := v.Type()
			numFields := v.NumField()
			for i := 0; i < numFields; i++ ***REMOVED***
				d.indent()
				vtf := vt.Field(i)
				d.w.Write([]byte(vtf.Name))
				d.w.Write(colonSpaceBytes)
				d.ignoreNextIndent = true
				d.dump(d.unpackValue(v.Field(i)))
				if i < (numFields - 1) ***REMOVED***
					d.w.Write(commaNewlineBytes)
				***REMOVED*** else ***REMOVED***
					d.w.Write(newlineBytes)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		d.depth--
		d.indent()
		d.w.Write(closeBraceBytes)

	case reflect.Uintptr:
		printHexPtr(d.w, uintptr(v.Uint()))

	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		printHexPtr(d.w, v.Pointer())

	// There were not any other types at the time this code was written, but
	// fall back to letting the default fmt package handle it in case any new
	// types are added.
	default:
		if v.CanInterface() ***REMOVED***
			fmt.Fprintf(d.w, "%v", v.Interface())
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(d.w, "%v", v.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

// fdump is a helper function to consolidate the logic from the various public
// methods which take varying writers and config states.
func fdump(cs *ConfigState, w io.Writer, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	for _, arg := range a ***REMOVED***
		if arg == nil ***REMOVED***
			w.Write(interfaceBytes)
			w.Write(spaceBytes)
			w.Write(nilAngleBytes)
			w.Write(newlineBytes)
			continue
		***REMOVED***

		d := dumpState***REMOVED***w: w, cs: cs***REMOVED***
		d.pointers = make(map[uintptr]int)
		d.dump(reflect.ValueOf(arg))
		d.w.Write(newlineBytes)
	***REMOVED***
***REMOVED***

// Fdump formats and displays the passed arguments to io.Writer w.  It formats
// exactly the same as Dump.
func Fdump(w io.Writer, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	fdump(&Config, w, a...)
***REMOVED***

// Sdump returns a string with the passed arguments formatted exactly the same
// as Dump.
func Sdump(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	var buf bytes.Buffer
	fdump(&Config, &buf, a...)
	return buf.String()
***REMOVED***

/*
Dump displays the passed parameters to standard out with newlines, customizable
indentation, and additional debug information such as complete types and all
pointer addresses used to indirect to the final value.  It provides the
following features over the built-in printing facilities provided by the fmt
package:

	* Pointers are dereferenced and followed
	* Circular data structures are detected and handled properly
	* Custom Stringer/error interfaces are optionally invoked, including
	  on unexported types
	* Custom types which only implement the Stringer/error interfaces via
	  a pointer receiver are optionally invoked when passing non-pointer
	  variables
	* Byte arrays and slices are dumped like the hexdump -C command which
	  includes offsets, byte values in hex, and ASCII output

The configuration options are controlled by an exported package global,
spew.Config.  See ConfigState for options documentation.

See Fdump if you would prefer dumping to an arbitrary io.Writer or Sdump to
get the formatted result as a string.
*/
func Dump(a ...interface***REMOVED******REMOVED***) ***REMOVED***
	fdump(&Config, os.Stdout, a...)
***REMOVED***
