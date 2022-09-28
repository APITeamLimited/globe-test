// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncore

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// ValidationError is an error type returned when attempting to validate a document or array.
type ValidationError string

func (ve ValidationError) Error() string ***REMOVED*** return string(ve) ***REMOVED***

// NewDocumentLengthError creates and returns an error for when the length of a document exceeds the
// bytes available.
func NewDocumentLengthError(length, rem int) error ***REMOVED***
	return lengthError("document", length, rem)
***REMOVED***

func lengthError(bufferType string, length, rem int) error ***REMOVED***
	return ValidationError(fmt.Sprintf("%v length exceeds available bytes. length=%d remainingBytes=%d",
		bufferType, length, rem))
***REMOVED***

// InsufficientBytesError indicates that there were not enough bytes to read the next component.
type InsufficientBytesError struct ***REMOVED***
	Source    []byte
	Remaining []byte
***REMOVED***

// NewInsufficientBytesError creates a new InsufficientBytesError with the given Document and
// remaining bytes.
func NewInsufficientBytesError(src, rem []byte) InsufficientBytesError ***REMOVED***
	return InsufficientBytesError***REMOVED***Source: src, Remaining: rem***REMOVED***
***REMOVED***

// Error implements the error interface.
func (ibe InsufficientBytesError) Error() string ***REMOVED***
	return "too few bytes to read next component"
***REMOVED***

// Equal checks that err2 also is an ErrTooSmall.
func (ibe InsufficientBytesError) Equal(err2 error) bool ***REMOVED***
	switch err2.(type) ***REMOVED***
	case InsufficientBytesError:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// InvalidDepthTraversalError is returned when attempting a recursive Lookup when one component of
// the path is neither an embedded document nor an array.
type InvalidDepthTraversalError struct ***REMOVED***
	Key  string
	Type bsontype.Type
***REMOVED***

func (idte InvalidDepthTraversalError) Error() string ***REMOVED***
	return fmt.Sprintf(
		"attempt to traverse into %s, but it's type is %s, not %s nor %s",
		idte.Key, idte.Type, bsontype.EmbeddedDocument, bsontype.Array,
	)
***REMOVED***

// ErrMissingNull is returned when a document or array's last byte is not null.
const ErrMissingNull ValidationError = "document or array end is missing null byte"

// ErrInvalidLength indicates that a length in a binary representation of a BSON document or array
// is invalid.
const ErrInvalidLength ValidationError = "document or array length is invalid"

// ErrNilReader indicates that an operation was attempted on a nil io.Reader.
var ErrNilReader = errors.New("nil reader")

// ErrEmptyKey indicates that no key was provided to a Lookup method.
var ErrEmptyKey = errors.New("empty key provided")

// ErrElementNotFound indicates that an Element matching a certain condition does not exist.
var ErrElementNotFound = errors.New("element not found")

// ErrOutOfBounds indicates that an index provided to access something was invalid.
var ErrOutOfBounds = errors.New("out of bounds")

// Document is a raw bytes representation of a BSON document.
type Document []byte

// NewDocumentFromReader reads a document from r. This function will only validate the length is
// correct and that the document ends with a null byte.
func NewDocumentFromReader(r io.Reader) (Document, error) ***REMOVED***
	return newBufferFromReader(r)
***REMOVED***

func newBufferFromReader(r io.Reader) ([]byte, error) ***REMOVED***
	if r == nil ***REMOVED***
		return nil, ErrNilReader
	***REMOVED***

	var lengthBytes [4]byte

	// ReadFull guarantees that we will have read at least len(lengthBytes) if err == nil
	_, err := io.ReadFull(r, lengthBytes[:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	length, _, _ := readi32(lengthBytes[:]) // ignore ok since we always have enough bytes to read a length
	if length < 0 ***REMOVED***
		return nil, ErrInvalidLength
	***REMOVED***
	buffer := make([]byte, length)

	copy(buffer, lengthBytes[:])

	_, err = io.ReadFull(r, buffer[4:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if buffer[length-1] != 0x00 ***REMOVED***
		return nil, ErrMissingNull
	***REMOVED***

	return buffer, nil
***REMOVED***

// Lookup searches the document, potentially recursively, for the given key. If there are multiple
// keys provided, this method will recurse down, as long as the top and intermediate nodes are
// either documents or arrays. If an error occurs or if the value doesn't exist, an empty Value is
// returned.
func (d Document) Lookup(key ...string) Value ***REMOVED***
	val, _ := d.LookupErr(key...)
	return val
***REMOVED***

// LookupErr is the same as Lookup, except it returns an error in addition to an empty Value.
func (d Document) LookupErr(key ...string) (Value, error) ***REMOVED***
	if len(key) < 1 ***REMOVED***
		return Value***REMOVED******REMOVED***, ErrEmptyKey
	***REMOVED***
	length, rem, ok := ReadLength(d)
	if !ok ***REMOVED***
		return Value***REMOVED******REMOVED***, NewInsufficientBytesError(d, rem)
	***REMOVED***

	length -= 4

	var elem Element
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return Value***REMOVED******REMOVED***, NewInsufficientBytesError(d, rem)
		***REMOVED***
		// We use `KeyBytes` rather than `Key` to avoid a needless string alloc.
		if string(elem.KeyBytes()) != key[0] ***REMOVED***
			continue
		***REMOVED***
		if len(key) > 1 ***REMOVED***
			tt := bsontype.Type(elem[0])
			switch tt ***REMOVED***
			case bsontype.EmbeddedDocument:
				val, err := elem.Value().Document().LookupErr(key[1:]...)
				if err != nil ***REMOVED***
					return Value***REMOVED******REMOVED***, err
				***REMOVED***
				return val, nil
			case bsontype.Array:
				// Convert to Document to continue Lookup recursion.
				val, err := Document(elem.Value().Array()).LookupErr(key[1:]...)
				if err != nil ***REMOVED***
					return Value***REMOVED******REMOVED***, err
				***REMOVED***
				return val, nil
			default:
				return Value***REMOVED******REMOVED***, InvalidDepthTraversalError***REMOVED***Key: elem.Key(), Type: tt***REMOVED***
			***REMOVED***
		***REMOVED***
		return elem.ValueErr()
	***REMOVED***
	return Value***REMOVED******REMOVED***, ErrElementNotFound
***REMOVED***

// Index searches for and retrieves the element at the given index. This method will panic if
// the document is invalid or if the index is out of bounds.
func (d Document) Index(index uint) Element ***REMOVED***
	elem, err := d.IndexErr(index)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return elem
***REMOVED***

// IndexErr searches for and retrieves the element at the given index.
func (d Document) IndexErr(index uint) (Element, error) ***REMOVED***
	return indexErr(d, index)
***REMOVED***

func indexErr(b []byte, index uint) (Element, error) ***REMOVED***
	length, rem, ok := ReadLength(b)
	if !ok ***REMOVED***
		return nil, NewInsufficientBytesError(b, rem)
	***REMOVED***

	length -= 4

	var current uint
	var elem Element
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return nil, NewInsufficientBytesError(b, rem)
		***REMOVED***
		if current != index ***REMOVED***
			current++
			continue
		***REMOVED***
		return elem, nil
	***REMOVED***
	return nil, ErrOutOfBounds
***REMOVED***

// DebugString outputs a human readable version of Document. It will attempt to stringify the
// valid components of the document even if the entire document is not valid.
func (d Document) DebugString() string ***REMOVED***
	if len(d) < 5 ***REMOVED***
		return "<malformed>"
	***REMOVED***
	var buf bytes.Buffer
	buf.WriteString("Document")
	length, rem, _ := ReadLength(d) // We know we have enough bytes to read the length
	buf.WriteByte('(')
	buf.WriteString(strconv.Itoa(int(length)))
	length -= 4
	buf.WriteString(")***REMOVED***")
	var elem Element
	var ok bool
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			buf.WriteString(fmt.Sprintf("<malformed (%d)>", length))
			break
		***REMOVED***
		fmt.Fprintf(&buf, "%s ", elem.DebugString())
	***REMOVED***
	buf.WriteByte('***REMOVED***')

	return buf.String()
***REMOVED***

// String outputs an ExtendedJSON version of Document. If the document is not valid, this method
// returns an empty string.
func (d Document) String() string ***REMOVED***
	if len(d) < 5 ***REMOVED***
		return ""
	***REMOVED***
	var buf bytes.Buffer
	buf.WriteByte('***REMOVED***')

	length, rem, _ := ReadLength(d) // We know we have enough bytes to read the length

	length -= 4

	var elem Element
	var ok bool
	first := true
	for length > 1 ***REMOVED***
		if !first ***REMOVED***
			buf.WriteByte(',')
		***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		fmt.Fprintf(&buf, "%s", elem.String())
		first = false
	***REMOVED***
	buf.WriteByte('***REMOVED***')

	return buf.String()
***REMOVED***

// Elements returns this document as a slice of elements. The returned slice will contain valid
// elements. If the document is not valid, the elements up to the invalid point will be returned
// along with an error.
func (d Document) Elements() ([]Element, error) ***REMOVED***
	length, rem, ok := ReadLength(d)
	if !ok ***REMOVED***
		return nil, NewInsufficientBytesError(d, rem)
	***REMOVED***

	length -= 4

	var elem Element
	var elems []Element
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return elems, NewInsufficientBytesError(d, rem)
		***REMOVED***
		if err := elem.Validate(); err != nil ***REMOVED***
			return elems, err
		***REMOVED***
		elems = append(elems, elem)
	***REMOVED***
	return elems, nil
***REMOVED***

// Values returns this document as a slice of values. The returned slice will contain valid values.
// If the document is not valid, the values up to the invalid point will be returned along with an
// error.
func (d Document) Values() ([]Value, error) ***REMOVED***
	return values(d)
***REMOVED***

func values(b []byte) ([]Value, error) ***REMOVED***
	length, rem, ok := ReadLength(b)
	if !ok ***REMOVED***
		return nil, NewInsufficientBytesError(b, rem)
	***REMOVED***

	length -= 4

	var elem Element
	var vals []Value
	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return vals, NewInsufficientBytesError(b, rem)
		***REMOVED***
		if err := elem.Value().Validate(); err != nil ***REMOVED***
			return vals, err
		***REMOVED***
		vals = append(vals, elem.Value())
	***REMOVED***
	return vals, nil
***REMOVED***

// Validate validates the document and ensures the elements contained within are valid.
func (d Document) Validate() error ***REMOVED***
	length, rem, ok := ReadLength(d)
	if !ok ***REMOVED***
		return NewInsufficientBytesError(d, rem)
	***REMOVED***
	if int(length) > len(d) ***REMOVED***
		return NewDocumentLengthError(int(length), len(d))
	***REMOVED***
	if d[length-1] != 0x00 ***REMOVED***
		return ErrMissingNull
	***REMOVED***

	length -= 4
	var elem Element

	for length > 1 ***REMOVED***
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok ***REMOVED***
			return NewInsufficientBytesError(d, rem)
		***REMOVED***
		err := elem.Validate()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if len(rem) < 1 || rem[0] != 0x00 ***REMOVED***
		return ErrMissingNull
	***REMOVED***
	return nil
***REMOVED***
