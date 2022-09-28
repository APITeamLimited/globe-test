// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx // import "go.mongodb.org/mongo-driver/x/bsonx"

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrNilArray indicates that an operation was attempted on a nil *Array.
var ErrNilArray = errors.New("array is nil")

// Arr represents an array in BSON.
type Arr []Val

// String implements the fmt.Stringer interface.
func (a Arr) String() string ***REMOVED***
	var buf bytes.Buffer
	buf.Write([]byte("bson.Array["))
	for idx, val := range a ***REMOVED***
		if idx > 0 ***REMOVED***
			buf.Write([]byte(", "))
		***REMOVED***
		fmt.Fprintf(&buf, "%s", val)
	***REMOVED***
	buf.WriteByte(']')

	return buf.String()
***REMOVED***

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
func (a Arr) MarshalBSONValue() (bsontype.Type, []byte, error) ***REMOVED***
	if a == nil ***REMOVED***
		// TODO: Should we do this?
		return bsontype.Null, nil, nil
	***REMOVED***

	idx, dst := bsoncore.ReserveLength(nil)
	for idx, value := range a ***REMOVED***
		t, data, _ := value.MarshalBSONValue() // marshalBSONValue never returns an error.
		dst = append(dst, byte(t))
		dst = append(dst, strconv.Itoa(idx)...)
		dst = append(dst, 0x00)
		dst = append(dst, data...)
	***REMOVED***
	dst = append(dst, 0x00)
	dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	return bsontype.Array, dst, nil
***REMOVED***

// UnmarshalBSONValue implements the bsoncodec.ValueUnmarshaler interface.
func (a *Arr) UnmarshalBSONValue(t bsontype.Type, data []byte) error ***REMOVED***
	if a == nil ***REMOVED***
		return ErrNilArray
	***REMOVED***
	*a = (*a)[:0]

	elements, err := bsoncore.Document(data).Elements()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, elem := range elements ***REMOVED***
		var val Val
		rawval := elem.Value()
		err = val.UnmarshalBSONValue(rawval.Type, rawval.Data)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*a = append(*a, val)
	***REMOVED***
	return nil
***REMOVED***

// Equal compares this document to another, returning true if they are equal.
func (a Arr) Equal(a2 Arr) bool ***REMOVED***
	if len(a) != len(a2) ***REMOVED***
		return false
	***REMOVED***
	for idx := range a ***REMOVED***
		if !a[idx].Equal(a2[idx]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (Arr) idoc() ***REMOVED******REMOVED***
