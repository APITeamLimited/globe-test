// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncore

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// NewArrayLengthError creates and returns an error for when the length of an array exceeds the
// bytes available.
func NewArrayLengthError(length, rem int) error ***REMOVED***
	return lengthError("array", length, rem)
***REMOVED***

// Array is a raw bytes representation of a BSON array.
type Array []byte

// NewArrayFromReader reads an array from r. This function will only validate the length is
// correct and that the array ends with a null byte.
func NewArrayFromReader(r io.Reader) (Array, error) ***REMOVED***
	return newBufferFromReader(r)
***REMOVED***

// Index searches for and retrieves the value at the given index. This method will panic if
// the array is invalid or if the index is out of bounds.
func (a Array) Index(index uint) Value ***REMOVED***
	value, err := a.IndexErr(index)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return value
***REMOVED***

// IndexErr searches for and retrieves the value at the given index.
func (a Array) IndexErr(index uint) (Value, error) ***REMOVED***
	elem, err := indexErr(a, index)
	if err != nil ***REMOVED***
		return Value***REMOVED******REMOVED***, err
	***REMOVED***
	return elem.Value(), err
***REMOVED***

// DebugString outputs a human readable version of Array. It will attempt to stringify the
// valid components of the array even if the entire array is not valid.
func (a Array) DebugString() string ***REMOVED***
	if len(a) < 5 ***REMOVED***
		return "<malformed>"
	***REMOVED***
	var buf bytes.Buffer
	buf.WriteString("Array")
	length, rem, _ := ReadLength(a) // We know we have enough bytes to read the length
	buf.WriteByte('(')
	buf.WriteString(strconv.Itoa(int(length)))
	length -= 4
	buf.WriteString(")[")
	var elem Element
	var ok bool
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			buf.WriteString(fmt.Sprintf("<malformed (%d)>", length))
			break
		***REMOVED***
		fmt.Fprintf(&buf, "%s", elem.Value().DebugString())
		if length != 1 ***REMOVED***
			buf.WriteByte(',')
		***REMOVED***
	***REMOVED***
	buf.WriteByte(']')

	return buf.String()
***REMOVED***

// String outputs an ExtendedJSON version of Array. If the Array is not valid, this method
// returns an empty string.
func (a Array) String() string ***REMOVED***
	if len(a) < 5 ***REMOVED***
		return ""
	***REMOVED***
	var buf bytes.Buffer
	buf.WriteByte('[')

	length, rem, _ := ReadLength(a) // We know we have enough bytes to read the length

	length -= 4

	var elem Element
	var ok bool
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		fmt.Fprintf(&buf, "%s", elem.Value().String())
		if length > 1 ***REMOVED***
			buf.WriteByte(',')
		***REMOVED***
	***REMOVED***
	if length != 1 ***REMOVED*** // Missing final null byte or inaccurate length
		return ""
	***REMOVED***

	buf.WriteByte(']')
	return buf.String()
***REMOVED***

// Values returns this array as a slice of values. The returned slice will contain valid values.
// If the array is not valid, the values up to the invalid point will be returned along with an
// error.
func (a Array) Values() ([]Value, error) ***REMOVED***
	return values(a)
***REMOVED***

// Validate validates the array and ensures the elements contained within are valid.
func (a Array) Validate() error ***REMOVED***
	length, rem, ok := ReadLength(a)
	if !ok ***REMOVED***
		return NewInsufficientBytesError(a, rem)
	***REMOVED***
	if int(length) > len(a) ***REMOVED***
		return NewArrayLengthError(int(length), len(a))
	***REMOVED***
	if a[length-1] != 0x00 ***REMOVED***
		return ErrMissingNull
	***REMOVED***

	length -= 4
	var elem Element

	var keyNum int64
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return NewInsufficientBytesError(a, rem)
		***REMOVED***

		// validate element
		err := elem.Validate()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// validate keys increase numerically
		if fmt.Sprint(keyNum) != elem.Key() ***REMOVED***
			return fmt.Errorf("array key %q is out of order or invalid", elem.Key())
		***REMOVED***
		keyNum++
	***REMOVED***

	if len(rem) < 1 || rem[0] != 0x00 ***REMOVED***
		return ErrMissingNull
	***REMOVED***
	return nil
***REMOVED***
