// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"bytes"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MDoc is an unordered, type safe, concise BSON document representation. This type should not be
// used if you require ordering of values or duplicate keys.
type MDoc map[string]Val

// ReadMDoc will create a Doc using the provided slice of bytes. If the
// slice of bytes is not a valid BSON document, this method will return an error.
func ReadMDoc(b []byte) (MDoc, error) ***REMOVED***
	doc := make(MDoc)
	err := doc.UnmarshalBSON(b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return doc, nil
***REMOVED***

// Copy makes a shallow copy of this document.
func (d MDoc) Copy() MDoc ***REMOVED***
	d2 := make(MDoc, len(d))
	for k, v := range d ***REMOVED***
		d2[k] = v
	***REMOVED***
	return d2
***REMOVED***

// Lookup searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Value if they key does not exist. To know if they key actually
// exists, use LookupErr.
func (d MDoc) Lookup(key ...string) Val ***REMOVED***
	val, _ := d.LookupErr(key...)
	return val
***REMOVED***

// LookupErr searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
func (d MDoc) LookupErr(key ...string) (Val, error) ***REMOVED***
	elem, err := d.LookupElementErr(key...)
	return elem.Value, err
***REMOVED***

// LookupElement searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Element if they key does not exist. To know if they key actually
// exists, use LookupElementErr.
func (d MDoc) LookupElement(key ...string) Elem ***REMOVED***
	elem, _ := d.LookupElementErr(key...)
	return elem
***REMOVED***

// LookupElementErr searches the document and potentially subdocuments for the
// provided key. Each key provided to this method represents a layer of depth.
func (d MDoc) LookupElementErr(key ...string) (Elem, error) ***REMOVED***
	// KeyNotFound operates by being created where the error happens and then the depth is
	// incremented by 1 as each function unwinds. Whenever this function returns, it also assigns
	// the Key slice to the key slice it has. This ensures that the proper depth is identified and
	// the proper keys.
	if len(key) == 0 ***REMOVED***
		return Elem***REMOVED******REMOVED***, KeyNotFound***REMOVED***Key: key***REMOVED***
	***REMOVED***

	var elem Elem
	var err error
	val, ok := d[key[0]]
	if !ok ***REMOVED***
		return Elem***REMOVED******REMOVED***, KeyNotFound***REMOVED***Key: key***REMOVED***
	***REMOVED***

	if len(key) == 1 ***REMOVED***
		return Elem***REMOVED***Key: key[0], Value: val***REMOVED***, nil
	***REMOVED***

	switch val.Type() ***REMOVED***
	case bsontype.EmbeddedDocument:
		switch tt := val.primitive.(type) ***REMOVED***
		case Doc:
			elem, err = tt.LookupElementErr(key[1:]...)
		case MDoc:
			elem, err = tt.LookupElementErr(key[1:]...)
		***REMOVED***
	default:
		return Elem***REMOVED******REMOVED***, KeyNotFound***REMOVED***Type: val.Type()***REMOVED***
	***REMOVED***
	switch tt := err.(type) ***REMOVED***
	case KeyNotFound:
		tt.Depth++
		tt.Key = key
		return Elem***REMOVED******REMOVED***, tt
	case nil:
		return elem, nil
	default:
		return Elem***REMOVED******REMOVED***, err // We can't actually hit this.
	***REMOVED***
***REMOVED***

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
//
// This method will never return an error.
func (d MDoc) MarshalBSONValue() (bsontype.Type, []byte, error) ***REMOVED***
	if d == nil ***REMOVED***
		// TODO: Should we do this?
		return bsontype.Null, nil, nil
	***REMOVED***
	data, _ := d.MarshalBSON()
	return bsontype.EmbeddedDocument, data, nil
***REMOVED***

// MarshalBSON implements the Marshaler interface.
//
// This method will never return an error.
func (d MDoc) MarshalBSON() ([]byte, error) ***REMOVED*** return d.AppendMarshalBSON(nil) ***REMOVED***

// AppendMarshalBSON marshals Doc to BSON bytes, appending to dst.
//
// This method will never return an error.
func (d MDoc) AppendMarshalBSON(dst []byte) ([]byte, error) ***REMOVED***
	idx, dst := bsoncore.ReserveLength(dst)
	for k, v := range d ***REMOVED***
		t, data, _ := v.MarshalBSONValue() // Value.MarshalBSONValue never returns an error.
		dst = append(dst, byte(t))
		dst = append(dst, k...)
		dst = append(dst, 0x00)
		dst = append(dst, data...)
	***REMOVED***
	dst = append(dst, 0x00)
	dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst, nil
***REMOVED***

// UnmarshalBSON implements the Unmarshaler interface.
func (d *MDoc) UnmarshalBSON(b []byte) error ***REMOVED***
	if d == nil ***REMOVED***
		return ErrNilDocument
	***REMOVED***

	if err := bsoncore.Document(b).Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***

	elems, err := bsoncore.Document(b).Elements()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var val Val
	for _, elem := range elems ***REMOVED***
		rawv := elem.Value()
		err = val.UnmarshalBSONValue(rawv.Type, rawv.Data)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		(*d)[elem.Key()] = val
	***REMOVED***
	return nil
***REMOVED***

// Equal compares this document to another, returning true if they are equal.
func (d MDoc) Equal(id IDoc) bool ***REMOVED***
	switch tt := id.(type) ***REMOVED***
	case MDoc:
		d2 := tt
		if len(d) != len(d2) ***REMOVED***
			return false
		***REMOVED***
		for key, value := range d ***REMOVED***
			value2, ok := d2[key]
			if !ok ***REMOVED***
				return false
			***REMOVED***
			if !value.Equal(value2) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	case Doc:
		unique := make(map[string]struct***REMOVED******REMOVED***)
		for _, elem := range tt ***REMOVED***
			unique[elem.Key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			val, ok := d[elem.Key]
			if !ok ***REMOVED***
				return false
			***REMOVED***
			if !val.Equal(elem.Value) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if len(unique) != len(d) ***REMOVED***
			return false
		***REMOVED***
	case nil:
		return d == nil
	default:
		return false
	***REMOVED***

	return true
***REMOVED***

// String implements the fmt.Stringer interface.
func (d MDoc) String() string ***REMOVED***
	var buf bytes.Buffer
	buf.Write([]byte("bson.Document***REMOVED***"))
	first := true
	for key, value := range d ***REMOVED***
		if !first ***REMOVED***
			buf.Write([]byte(", "))
		***REMOVED***
		fmt.Fprintf(&buf, "%v", Elem***REMOVED***Key: key, Value: value***REMOVED***)
		first = false
	***REMOVED***
	buf.WriteByte('***REMOVED***')

	return buf.String()
***REMOVED***

func (MDoc) idoc() ***REMOVED******REMOVED***
