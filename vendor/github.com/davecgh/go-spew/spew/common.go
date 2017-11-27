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
	"io"
	"reflect"
	"sort"
	"strconv"
)

// Some constants in the form of bytes to avoid string overhead.  This mirrors
// the technique used in the fmt package.
var (
	panicBytes            = []byte("(PANIC=")
	plusBytes             = []byte("+")
	iBytes                = []byte("i")
	trueBytes             = []byte("true")
	falseBytes            = []byte("false")
	interfaceBytes        = []byte("(interface ***REMOVED******REMOVED***)")
	commaNewlineBytes     = []byte(",\n")
	newlineBytes          = []byte("\n")
	openBraceBytes        = []byte("***REMOVED***")
	openBraceNewlineBytes = []byte("***REMOVED***\n")
	closeBraceBytes       = []byte("***REMOVED***")
	asteriskBytes         = []byte("*")
	colonBytes            = []byte(":")
	colonSpaceBytes       = []byte(": ")
	openParenBytes        = []byte("(")
	closeParenBytes       = []byte(")")
	spaceBytes            = []byte(" ")
	pointerChainBytes     = []byte("->")
	nilAngleBytes         = []byte("<nil>")
	maxNewlineBytes       = []byte("<max depth reached>\n")
	maxShortBytes         = []byte("<max>")
	circularBytes         = []byte("<already shown>")
	circularShortBytes    = []byte("<shown>")
	invalidAngleBytes     = []byte("<invalid>")
	openBracketBytes      = []byte("[")
	closeBracketBytes     = []byte("]")
	percentBytes          = []byte("%")
	precisionBytes        = []byte(".")
	openAngleBytes        = []byte("<")
	closeAngleBytes       = []byte(">")
	openMapBytes          = []byte("map[")
	closeMapBytes         = []byte("]")
	lenEqualsBytes        = []byte("len=")
	capEqualsBytes        = []byte("cap=")
)

// hexDigits is used to map a decimal value to a hex digit.
var hexDigits = "0123456789abcdef"

// catchPanic handles any panics that might occur during the handleMethods
// calls.
func catchPanic(w io.Writer, v reflect.Value) ***REMOVED***
	if err := recover(); err != nil ***REMOVED***
		w.Write(panicBytes)
		fmt.Fprintf(w, "%v", err)
		w.Write(closeParenBytes)
	***REMOVED***
***REMOVED***

// handleMethods attempts to call the Error and String methods on the underlying
// type the passed reflect.Value represents and outputes the result to Writer w.
//
// It handles panics in any called methods by catching and displaying the error
// as the formatted value.
func handleMethods(cs *ConfigState, w io.Writer, v reflect.Value) (handled bool) ***REMOVED***
	// We need an interface to check if the type implements the error or
	// Stringer interface.  However, the reflect package won't give us an
	// interface on certain things like unexported struct fields in order
	// to enforce visibility rules.  We use unsafe, when it's available,
	// to bypass these restrictions since this package does not mutate the
	// values.
	if !v.CanInterface() ***REMOVED***
		if UnsafeDisabled ***REMOVED***
			return false
		***REMOVED***

		v = unsafeReflectValue(v)
	***REMOVED***

	// Choose whether or not to do error and Stringer interface lookups against
	// the base type or a pointer to the base type depending on settings.
	// Technically calling one of these methods with a pointer receiver can
	// mutate the value, however, types which choose to satisify an error or
	// Stringer interface with a pointer receiver should not be mutating their
	// state inside these interface methods.
	if !cs.DisablePointerMethods && !UnsafeDisabled && !v.CanAddr() ***REMOVED***
		v = unsafeReflectValue(v)
	***REMOVED***
	if v.CanAddr() ***REMOVED***
		v = v.Addr()
	***REMOVED***

	// Is it an error or Stringer?
	switch iface := v.Interface().(type) ***REMOVED***
	case error:
		defer catchPanic(w, v)
		if cs.ContinueOnMethod ***REMOVED***
			w.Write(openParenBytes)
			w.Write([]byte(iface.Error()))
			w.Write(closeParenBytes)
			w.Write(spaceBytes)
			return false
		***REMOVED***

		w.Write([]byte(iface.Error()))
		return true

	case fmt.Stringer:
		defer catchPanic(w, v)
		if cs.ContinueOnMethod ***REMOVED***
			w.Write(openParenBytes)
			w.Write([]byte(iface.String()))
			w.Write(closeParenBytes)
			w.Write(spaceBytes)
			return false
		***REMOVED***
		w.Write([]byte(iface.String()))
		return true
	***REMOVED***
	return false
***REMOVED***

// printBool outputs a boolean value as true or false to Writer w.
func printBool(w io.Writer, val bool) ***REMOVED***
	if val ***REMOVED***
		w.Write(trueBytes)
	***REMOVED*** else ***REMOVED***
		w.Write(falseBytes)
	***REMOVED***
***REMOVED***

// printInt outputs a signed integer value to Writer w.
func printInt(w io.Writer, val int64, base int) ***REMOVED***
	w.Write([]byte(strconv.FormatInt(val, base)))
***REMOVED***

// printUint outputs an unsigned integer value to Writer w.
func printUint(w io.Writer, val uint64, base int) ***REMOVED***
	w.Write([]byte(strconv.FormatUint(val, base)))
***REMOVED***

// printFloat outputs a floating point value using the specified precision,
// which is expected to be 32 or 64bit, to Writer w.
func printFloat(w io.Writer, val float64, precision int) ***REMOVED***
	w.Write([]byte(strconv.FormatFloat(val, 'g', -1, precision)))
***REMOVED***

// printComplex outputs a complex value using the specified float precision
// for the real and imaginary parts to Writer w.
func printComplex(w io.Writer, c complex128, floatPrecision int) ***REMOVED***
	r := real(c)
	w.Write(openParenBytes)
	w.Write([]byte(strconv.FormatFloat(r, 'g', -1, floatPrecision)))
	i := imag(c)
	if i >= 0 ***REMOVED***
		w.Write(plusBytes)
	***REMOVED***
	w.Write([]byte(strconv.FormatFloat(i, 'g', -1, floatPrecision)))
	w.Write(iBytes)
	w.Write(closeParenBytes)
***REMOVED***

// printHexPtr outputs a uintptr formatted as hexadecimal with a leading '0x'
// prefix to Writer w.
func printHexPtr(w io.Writer, p uintptr) ***REMOVED***
	// Null pointer.
	num := uint64(p)
	if num == 0 ***REMOVED***
		w.Write(nilAngleBytes)
		return
	***REMOVED***

	// Max uint64 is 16 bytes in hex + 2 bytes for '0x' prefix
	buf := make([]byte, 18)

	// It's simpler to construct the hex string right to left.
	base := uint64(16)
	i := len(buf) - 1
	for num >= base ***REMOVED***
		buf[i] = hexDigits[num%base]
		num /= base
		i--
	***REMOVED***
	buf[i] = hexDigits[num]

	// Add '0x' prefix.
	i--
	buf[i] = 'x'
	i--
	buf[i] = '0'

	// Strip unused leading bytes.
	buf = buf[i:]
	w.Write(buf)
***REMOVED***

// valuesSorter implements sort.Interface to allow a slice of reflect.Value
// elements to be sorted.
type valuesSorter struct ***REMOVED***
	values  []reflect.Value
	strings []string // either nil or same len and values
	cs      *ConfigState
***REMOVED***

// newValuesSorter initializes a valuesSorter instance, which holds a set of
// surrogate keys on which the data should be sorted.  It uses flags in
// ConfigState to decide if and how to populate those surrogate keys.
func newValuesSorter(values []reflect.Value, cs *ConfigState) sort.Interface ***REMOVED***
	vs := &valuesSorter***REMOVED***values: values, cs: cs***REMOVED***
	if canSortSimply(vs.values[0].Kind()) ***REMOVED***
		return vs
	***REMOVED***
	if !cs.DisableMethods ***REMOVED***
		vs.strings = make([]string, len(values))
		for i := range vs.values ***REMOVED***
			b := bytes.Buffer***REMOVED******REMOVED***
			if !handleMethods(cs, &b, vs.values[i]) ***REMOVED***
				vs.strings = nil
				break
			***REMOVED***
			vs.strings[i] = b.String()
		***REMOVED***
	***REMOVED***
	if vs.strings == nil && cs.SpewKeys ***REMOVED***
		vs.strings = make([]string, len(values))
		for i := range vs.values ***REMOVED***
			vs.strings[i] = Sprintf("%#v", vs.values[i].Interface())
		***REMOVED***
	***REMOVED***
	return vs
***REMOVED***

// canSortSimply tests whether a reflect.Kind is a primitive that can be sorted
// directly, or whether it should be considered for sorting by surrogate keys
// (if the ConfigState allows it).
func canSortSimply(kind reflect.Kind) bool ***REMOVED***
	// This switch parallels valueSortLess, except for the default case.
	switch kind ***REMOVED***
	case reflect.Bool:
		return true
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return true
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	case reflect.Uintptr:
		return true
	case reflect.Array:
		return true
	***REMOVED***
	return false
***REMOVED***

// Len returns the number of values in the slice.  It is part of the
// sort.Interface implementation.
func (s *valuesSorter) Len() int ***REMOVED***
	return len(s.values)
***REMOVED***

// Swap swaps the values at the passed indices.  It is part of the
// sort.Interface implementation.
func (s *valuesSorter) Swap(i, j int) ***REMOVED***
	s.values[i], s.values[j] = s.values[j], s.values[i]
	if s.strings != nil ***REMOVED***
		s.strings[i], s.strings[j] = s.strings[j], s.strings[i]
	***REMOVED***
***REMOVED***

// valueSortLess returns whether the first value should sort before the second
// value.  It is used by valueSorter.Less as part of the sort.Interface
// implementation.
func valueSortLess(a, b reflect.Value) bool ***REMOVED***
	switch a.Kind() ***REMOVED***
	case reflect.Bool:
		return !a.Bool() && b.Bool()
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return a.Int() < b.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return a.Uint() < b.Uint()
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float()
	case reflect.String:
		return a.String() < b.String()
	case reflect.Uintptr:
		return a.Uint() < b.Uint()
	case reflect.Array:
		// Compare the contents of both arrays.
		l := a.Len()
		for i := 0; i < l; i++ ***REMOVED***
			av := a.Index(i)
			bv := b.Index(i)
			if av.Interface() == bv.Interface() ***REMOVED***
				continue
			***REMOVED***
			return valueSortLess(av, bv)
		***REMOVED***
	***REMOVED***
	return a.String() < b.String()
***REMOVED***

// Less returns whether the value at index i should sort before the
// value at index j.  It is part of the sort.Interface implementation.
func (s *valuesSorter) Less(i, j int) bool ***REMOVED***
	if s.strings == nil ***REMOVED***
		return valueSortLess(s.values[i], s.values[j])
	***REMOVED***
	return s.strings[i] < s.strings[j]
***REMOVED***

// sortValues is a sort function that handles both native types and any type that
// can be converted to error or Stringer.  Other inputs are sorted according to
// their Value.String() value to ensure display stability.
func sortValues(values []reflect.Value, cs *ConfigState) ***REMOVED***
	if len(values) == 0 ***REMOVED***
		return
	***REMOVED***
	sort.Sort(newValuesSorter(values, cs))
***REMOVED***
