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
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// supportedFlags is a list of all the character flags supported by fmt package.
const supportedFlags = "0-+# "

// formatState implements the fmt.Formatter interface and contains information
// about the state of a formatting operation.  The NewFormatter function can
// be used to get a new Formatter which can be used directly as arguments
// in standard fmt package printing calls.
type formatState struct ***REMOVED***
	value          interface***REMOVED******REMOVED***
	fs             fmt.State
	depth          int
	pointers       map[uintptr]int
	ignoreNextType bool
	cs             *ConfigState
***REMOVED***

// buildDefaultFormat recreates the original format string without precision
// and width information to pass in to fmt.Sprintf in the case of an
// unrecognized type.  Unless new types are added to the language, this
// function won't ever be called.
func (f *formatState) buildDefaultFormat() (format string) ***REMOVED***
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range supportedFlags ***REMOVED***
		if f.fs.Flag(int(flag)) ***REMOVED***
			buf.WriteRune(flag)
		***REMOVED***
	***REMOVED***

	buf.WriteRune('v')

	format = buf.String()
	return format
***REMOVED***

// constructOrigFormat recreates the original format string including precision
// and width information to pass along to the standard fmt package.  This allows
// automatic deferral of all format strings this package doesn't support.
func (f *formatState) constructOrigFormat(verb rune) (format string) ***REMOVED***
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range supportedFlags ***REMOVED***
		if f.fs.Flag(int(flag)) ***REMOVED***
			buf.WriteRune(flag)
		***REMOVED***
	***REMOVED***

	if width, ok := f.fs.Width(); ok ***REMOVED***
		buf.WriteString(strconv.Itoa(width))
	***REMOVED***

	if precision, ok := f.fs.Precision(); ok ***REMOVED***
		buf.Write(precisionBytes)
		buf.WriteString(strconv.Itoa(precision))
	***REMOVED***

	buf.WriteRune(verb)

	format = buf.String()
	return format
***REMOVED***

// unpackValue returns values inside of non-nil interfaces when possible and
// ensures that types for values which have been unpacked from an interface
// are displayed when the show types flag is also set.
// This is useful for data types like structs, arrays, slices, and maps which
// can contain varying types packed inside an interface.
func (f *formatState) unpackValue(v reflect.Value) reflect.Value ***REMOVED***
	if v.Kind() == reflect.Interface ***REMOVED***
		f.ignoreNextType = false
		if !v.IsNil() ***REMOVED***
			v = v.Elem()
		***REMOVED***
	***REMOVED***
	return v
***REMOVED***

// formatPtr handles formatting of pointers by indirecting them as necessary.
func (f *formatState) formatPtr(v reflect.Value) ***REMOVED***
	// Display nil if top level pointer is nil.
	showTypes := f.fs.Flag('#')
	if v.IsNil() && (!showTypes || f.ignoreNextType) ***REMOVED***
		f.fs.Write(nilAngleBytes)
		return
	***REMOVED***

	// Remove pointers at or below the current depth from map used to detect
	// circular refs.
	for k, depth := range f.pointers ***REMOVED***
		if depth >= f.depth ***REMOVED***
			delete(f.pointers, k)
		***REMOVED***
	***REMOVED***

	// Keep list of all dereferenced pointers to possibly show later.
	pointerChain := make([]uintptr, 0)

	// Figure out how many levels of indirection there are by derferencing
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
		if pd, ok := f.pointers[addr]; ok && pd < f.depth ***REMOVED***
			cycleFound = true
			indirects--
			break
		***REMOVED***
		f.pointers[addr] = f.depth

		ve = ve.Elem()
		if ve.Kind() == reflect.Interface ***REMOVED***
			if ve.IsNil() ***REMOVED***
				nilFound = true
				break
			***REMOVED***
			ve = ve.Elem()
		***REMOVED***
	***REMOVED***

	// Display type or indirection level depending on flags.
	if showTypes && !f.ignoreNextType ***REMOVED***
		f.fs.Write(openParenBytes)
		f.fs.Write(bytes.Repeat(asteriskBytes, indirects))
		f.fs.Write([]byte(ve.Type().String()))
		f.fs.Write(closeParenBytes)
	***REMOVED*** else ***REMOVED***
		if nilFound || cycleFound ***REMOVED***
			indirects += strings.Count(ve.Type().String(), "*")
		***REMOVED***
		f.fs.Write(openAngleBytes)
		f.fs.Write([]byte(strings.Repeat("*", indirects)))
		f.fs.Write(closeAngleBytes)
	***REMOVED***

	// Display pointer information depending on flags.
	if f.fs.Flag('+') && (len(pointerChain) > 0) ***REMOVED***
		f.fs.Write(openParenBytes)
		for i, addr := range pointerChain ***REMOVED***
			if i > 0 ***REMOVED***
				f.fs.Write(pointerChainBytes)
			***REMOVED***
			printHexPtr(f.fs, addr)
		***REMOVED***
		f.fs.Write(closeParenBytes)
	***REMOVED***

	// Display dereferenced value.
	switch ***REMOVED***
	case nilFound:
		f.fs.Write(nilAngleBytes)

	case cycleFound:
		f.fs.Write(circularShortBytes)

	default:
		f.ignoreNextType = true
		f.format(ve)
	***REMOVED***
***REMOVED***

// format is the main workhorse for providing the Formatter interface.  It
// uses the passed reflect value to figure out what kind of object we are
// dealing with and formats it appropriately.  It is a recursive function,
// however circular data structures are detected and handled properly.
func (f *formatState) format(v reflect.Value) ***REMOVED***
	// Handle invalid reflect values immediately.
	kind := v.Kind()
	if kind == reflect.Invalid ***REMOVED***
		f.fs.Write(invalidAngleBytes)
		return
	***REMOVED***

	// Handle pointers specially.
	if kind == reflect.Ptr ***REMOVED***
		f.formatPtr(v)
		return
	***REMOVED***

	// Print type information unless already handled elsewhere.
	if !f.ignoreNextType && f.fs.Flag('#') ***REMOVED***
		f.fs.Write(openParenBytes)
		f.fs.Write([]byte(v.Type().String()))
		f.fs.Write(closeParenBytes)
	***REMOVED***
	f.ignoreNextType = false

	// Call Stringer/error interfaces if they exist and the handle methods
	// flag is enabled.
	if !f.cs.DisableMethods ***REMOVED***
		if (kind != reflect.Invalid) && (kind != reflect.Interface) ***REMOVED***
			if handled := handleMethods(f.cs, f.fs, v); handled ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	switch kind ***REMOVED***
	case reflect.Invalid:
		// Do nothing.  We should never get here since invalid has already
		// been handled above.

	case reflect.Bool:
		printBool(f.fs, v.Bool())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		printInt(f.fs, v.Int(), 10)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		printUint(f.fs, v.Uint(), 10)

	case reflect.Float32:
		printFloat(f.fs, v.Float(), 32)

	case reflect.Float64:
		printFloat(f.fs, v.Float(), 64)

	case reflect.Complex64:
		printComplex(f.fs, v.Complex(), 32)

	case reflect.Complex128:
		printComplex(f.fs, v.Complex(), 64)

	case reflect.Slice:
		if v.IsNil() ***REMOVED***
			f.fs.Write(nilAngleBytes)
			break
		***REMOVED***
		fallthrough

	case reflect.Array:
		f.fs.Write(openBracketBytes)
		f.depth++
		if (f.cs.MaxDepth != 0) && (f.depth > f.cs.MaxDepth) ***REMOVED***
			f.fs.Write(maxShortBytes)
		***REMOVED*** else ***REMOVED***
			numEntries := v.Len()
			for i := 0; i < numEntries; i++ ***REMOVED***
				if i > 0 ***REMOVED***
					f.fs.Write(spaceBytes)
				***REMOVED***
				f.ignoreNextType = true
				f.format(f.unpackValue(v.Index(i)))
			***REMOVED***
		***REMOVED***
		f.depth--
		f.fs.Write(closeBracketBytes)

	case reflect.String:
		f.fs.Write([]byte(v.String()))

	case reflect.Interface:
		// The only time we should get here is for nil interfaces due to
		// unpackValue calls.
		if v.IsNil() ***REMOVED***
			f.fs.Write(nilAngleBytes)
		***REMOVED***

	case reflect.Ptr:
		// Do nothing.  We should never get here since pointers have already
		// been handled above.

	case reflect.Map:
		// nil maps should be indicated as different than empty maps
		if v.IsNil() ***REMOVED***
			f.fs.Write(nilAngleBytes)
			break
		***REMOVED***

		f.fs.Write(openMapBytes)
		f.depth++
		if (f.cs.MaxDepth != 0) && (f.depth > f.cs.MaxDepth) ***REMOVED***
			f.fs.Write(maxShortBytes)
		***REMOVED*** else ***REMOVED***
			keys := v.MapKeys()
			if f.cs.SortKeys ***REMOVED***
				sortValues(keys, f.cs)
			***REMOVED***
			for i, key := range keys ***REMOVED***
				if i > 0 ***REMOVED***
					f.fs.Write(spaceBytes)
				***REMOVED***
				f.ignoreNextType = true
				f.format(f.unpackValue(key))
				f.fs.Write(colonBytes)
				f.ignoreNextType = true
				f.format(f.unpackValue(v.MapIndex(key)))
			***REMOVED***
		***REMOVED***
		f.depth--
		f.fs.Write(closeMapBytes)

	case reflect.Struct:
		numFields := v.NumField()
		f.fs.Write(openBraceBytes)
		f.depth++
		if (f.cs.MaxDepth != 0) && (f.depth > f.cs.MaxDepth) ***REMOVED***
			f.fs.Write(maxShortBytes)
		***REMOVED*** else ***REMOVED***
			vt := v.Type()
			for i := 0; i < numFields; i++ ***REMOVED***
				if i > 0 ***REMOVED***
					f.fs.Write(spaceBytes)
				***REMOVED***
				vtf := vt.Field(i)
				if f.fs.Flag('+') || f.fs.Flag('#') ***REMOVED***
					f.fs.Write([]byte(vtf.Name))
					f.fs.Write(colonBytes)
				***REMOVED***
				f.format(f.unpackValue(v.Field(i)))
			***REMOVED***
		***REMOVED***
		f.depth--
		f.fs.Write(closeBraceBytes)

	case reflect.Uintptr:
		printHexPtr(f.fs, uintptr(v.Uint()))

	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		printHexPtr(f.fs, v.Pointer())

	// There were not any other types at the time this code was written, but
	// fall back to letting the default fmt package handle it if any get added.
	default:
		format := f.buildDefaultFormat()
		if v.CanInterface() ***REMOVED***
			fmt.Fprintf(f.fs, format, v.Interface())
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(f.fs, format, v.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

// Format satisfies the fmt.Formatter interface. See NewFormatter for usage
// details.
func (f *formatState) Format(fs fmt.State, verb rune) ***REMOVED***
	f.fs = fs

	// Use standard formatting for verbs that are not v.
	if verb != 'v' ***REMOVED***
		format := f.constructOrigFormat(verb)
		fmt.Fprintf(fs, format, f.value)
		return
	***REMOVED***

	if f.value == nil ***REMOVED***
		if fs.Flag('#') ***REMOVED***
			fs.Write(interfaceBytes)
		***REMOVED***
		fs.Write(nilAngleBytes)
		return
	***REMOVED***

	f.format(reflect.ValueOf(f.value))
***REMOVED***

// newFormatter is a helper function to consolidate the logic from the various
// public methods which take varying config states.
func newFormatter(cs *ConfigState, v interface***REMOVED******REMOVED***) fmt.Formatter ***REMOVED***
	fs := &formatState***REMOVED***value: v, cs: cs***REMOVED***
	fs.pointers = make(map[uintptr]int)
	return fs
***REMOVED***

/*
NewFormatter returns a custom formatter that satisfies the fmt.Formatter
interface.  As a result, it integrates cleanly with standard fmt package
printing functions.  The formatter is useful for inline printing of smaller data
types similar to the standard %v format specifier.

The custom formatter only responds to the %v (most compact), %+v (adds pointer
addresses), %#v (adds types), or %#+v (adds types and pointer addresses) verb
combinations.  Any other verbs such as %x and %q will be sent to the the
standard fmt package for formatting.  In addition, the custom formatter ignores
the width and precision arguments (however they will still work on the format
specifiers not handled by the custom formatter).

Typically this function shouldn't be called directly.  It is much easier to make
use of the custom formatter by calling one of the convenience functions such as
Printf, Println, or Fprintf.
*/
func NewFormatter(v interface***REMOVED******REMOVED***) fmt.Formatter ***REMOVED***
	return newFormatter(&Config, v)
***REMOVED***
