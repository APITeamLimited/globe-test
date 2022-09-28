// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncore

import (
	"bytes"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// MalformedElementError represents a class of errors that RawElement methods return.
type MalformedElementError string

func (mee MalformedElementError) Error() string ***REMOVED*** return string(mee) ***REMOVED***

// ErrElementMissingKey is returned when a RawElement is missing a key.
const ErrElementMissingKey MalformedElementError = "element is missing key"

// ErrElementMissingType is returned when a RawElement is missing a type.
const ErrElementMissingType MalformedElementError = "element is missing type"

// Element is a raw bytes representation of a BSON element.
type Element []byte

// Key returns the key for this element. If the element is not valid, this method returns an empty
// string. If knowing if the element is valid is important, use KeyErr.
func (e Element) Key() string ***REMOVED***
	key, _ := e.KeyErr()
	return key
***REMOVED***

// KeyBytes returns the key for this element as a []byte. If the element is not valid, this method
// returns an empty string. If knowing if the element is valid is important, use KeyErr. This method
// will not include the null byte at the end of the key in the slice of bytes.
func (e Element) KeyBytes() []byte ***REMOVED***
	key, _ := e.KeyBytesErr()
	return key
***REMOVED***

// KeyErr returns the key for this element, returning an error if the element is not valid.
func (e Element) KeyErr() (string, error) ***REMOVED***
	key, err := e.KeyBytesErr()
	return string(key), err
***REMOVED***

// KeyBytesErr returns the key for this element as a []byte, returning an error if the element is
// not valid.
func (e Element) KeyBytesErr() ([]byte, error) ***REMOVED***
	if len(e) <= 0 ***REMOVED***
		return nil, ErrElementMissingType
	***REMOVED***
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 ***REMOVED***
		return nil, ErrElementMissingKey
	***REMOVED***
	return e[1 : idx+1], nil
***REMOVED***

// Validate ensures the element is a valid BSON element.
func (e Element) Validate() error ***REMOVED***
	if len(e) < 1 ***REMOVED***
		return ErrElementMissingType
	***REMOVED***
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 ***REMOVED***
		return ErrElementMissingKey
	***REMOVED***
	return Value***REMOVED***Type: bsontype.Type(e[0]), Data: e[idx+2:]***REMOVED***.Validate()
***REMOVED***

// CompareKey will compare this element's key to key. This method makes it easy to compare keys
// without needing to allocate a string. The key may be null terminated. If a valid key cannot be
// read this method will return false.
func (e Element) CompareKey(key []byte) bool ***REMOVED***
	if len(e) < 2 ***REMOVED***
		return false
	***REMOVED***
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 ***REMOVED***
		return false
	***REMOVED***
	if index := bytes.IndexByte(key, 0x00); index > -1 ***REMOVED***
		key = key[:index]
	***REMOVED***
	return bytes.Equal(e[1:idx+1], key)
***REMOVED***

// Value returns the value of this element. If the element is not valid, this method returns an
// empty Value. If knowing if the element is valid is important, use ValueErr.
func (e Element) Value() Value ***REMOVED***
	val, _ := e.ValueErr()
	return val
***REMOVED***

// ValueErr returns the value for this element, returning an error if the element is not valid.
func (e Element) ValueErr() (Value, error) ***REMOVED***
	if len(e) <= 0 ***REMOVED***
		return Value***REMOVED******REMOVED***, ErrElementMissingType
	***REMOVED***
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 ***REMOVED***
		return Value***REMOVED******REMOVED***, ErrElementMissingKey
	***REMOVED***

	val, rem, exists := ReadValue(e[idx+2:], bsontype.Type(e[0]))
	if !exists ***REMOVED***
		return Value***REMOVED******REMOVED***, NewInsufficientBytesError(e, rem)
	***REMOVED***
	return val, nil
***REMOVED***

// String implements the fmt.String interface. The output will be in extended JSON format.
func (e Element) String() string ***REMOVED***
	if len(e) <= 0 ***REMOVED***
		return ""
	***REMOVED***
	t := bsontype.Type(e[0])
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 ***REMOVED***
		return ""
	***REMOVED***
	key, valBytes := []byte(e[1:idx+1]), []byte(e[idx+2:])
	val, _, valid := ReadValue(valBytes, t)
	if !valid ***REMOVED***
		return ""
	***REMOVED***
	return fmt.Sprintf(`"%s": %v`, key, val)
***REMOVED***

// DebugString outputs a human readable version of RawElement. It will attempt to stringify the
// valid components of the element even if the entire element is not valid.
func (e Element) DebugString() string ***REMOVED***
	if len(e) <= 0 ***REMOVED***
		return "<malformed>"
	***REMOVED***
	t := bsontype.Type(e[0])
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 ***REMOVED***
		return fmt.Sprintf(`bson.Element***REMOVED***[%s]<malformed>***REMOVED***`, t)
	***REMOVED***
	key, valBytes := []byte(e[1:idx+1]), []byte(e[idx+2:])
	val, _, valid := ReadValue(valBytes, t)
	if !valid ***REMOVED***
		return fmt.Sprintf(`bson.Element***REMOVED***[%s]"%s": <malformed>***REMOVED***`, t, key)
	***REMOVED***
	return fmt.Sprintf(`bson.Element***REMOVED***[%s]"%s": %v***REMOVED***`, t, key, val)
***REMOVED***
