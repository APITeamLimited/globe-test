// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"bytes"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrNilDocument indicates that an operation was attempted on a nil *bson.Document.
var ErrNilDocument = errors.New("document is nil")

// KeyNotFound is an error type returned from the Lookup methods on Document. This type contains
// information about which key was not found and if it was actually not found or if a component of
// the key except the last was not a document nor array.
type KeyNotFound struct ***REMOVED***
	Key   []string      // The keys that were searched for.
	Depth uint          // Which key either was not found or was an incorrect type.
	Type  bsontype.Type // The type of the key that was found but was an incorrect type.
***REMOVED***

func (knf KeyNotFound) Error() string ***REMOVED***
	depth := knf.Depth
	if depth >= uint(len(knf.Key)) ***REMOVED***
		depth = uint(len(knf.Key)) - 1
	***REMOVED***

	if len(knf.Key) == 0 ***REMOVED***
		return "no keys were provided for lookup"
	***REMOVED***

	if knf.Type != bsontype.Type(0) ***REMOVED***
		return fmt.Sprintf(`key "%s" was found but was not valid to traverse BSON type %s`, knf.Key[depth], knf.Type)
	***REMOVED***

	return fmt.Sprintf(`key "%s" was not found`, knf.Key[depth])
***REMOVED***

// Doc is a type safe, concise BSON document representation.
type Doc []Elem

// ReadDoc will create a Document using the provided slice of bytes. If the
// slice of bytes is not a valid BSON document, this method will return an error.
func ReadDoc(b []byte) (Doc, error) ***REMOVED***
	doc := make(Doc, 0)
	err := doc.UnmarshalBSON(b)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return doc, nil
***REMOVED***

// Copy makes a shallow copy of this document.
func (d Doc) Copy() Doc ***REMOVED***
	d2 := make(Doc, len(d))
	copy(d2, d)
	return d2
***REMOVED***

// Append adds an element to the end of the document, creating it from the key and value provided.
func (d Doc) Append(key string, val Val) Doc ***REMOVED***
	return append(d, Elem***REMOVED***Key: key, Value: val***REMOVED***)
***REMOVED***

// Prepend adds an element to the beginning of the document, creating it from the key and value provided.
func (d Doc) Prepend(key string, val Val) Doc ***REMOVED***
	// TODO: should we just modify d itself instead of doing an alloc here?
	return append(Doc***REMOVED******REMOVED***Key: key, Value: val***REMOVED******REMOVED***, d...)
***REMOVED***

// Set replaces an element of a document. If an element with a matching key is
// found, the element will be replaced with the one provided. If the document
// does not have an element with that key, the element is appended to the
// document instead.
func (d Doc) Set(key string, val Val) Doc ***REMOVED***
	idx := d.IndexOf(key)
	if idx == -1 ***REMOVED***
		return append(d, Elem***REMOVED***Key: key, Value: val***REMOVED***)
	***REMOVED***
	d[idx] = Elem***REMOVED***Key: key, Value: val***REMOVED***
	return d
***REMOVED***

// IndexOf returns the index of the first element with a key of key, or -1 if no element with a key
// was found.
func (d Doc) IndexOf(key string) int ***REMOVED***
	for i, e := range d ***REMOVED***
		if e.Key == key ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// Delete removes the element with key if it exists and returns the updated Doc.
func (d Doc) Delete(key string) Doc ***REMOVED***
	idx := d.IndexOf(key)
	if idx == -1 ***REMOVED***
		return d
	***REMOVED***
	return append(d[:idx], d[idx+1:]...)
***REMOVED***

// Lookup searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Value if they key does not exist. To know if they key actually
// exists, use LookupErr.
func (d Doc) Lookup(key ...string) Val ***REMOVED***
	val, _ := d.LookupErr(key...)
	return val
***REMOVED***

// LookupErr searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
func (d Doc) LookupErr(key ...string) (Val, error) ***REMOVED***
	elem, err := d.LookupElementErr(key...)
	return elem.Value, err
***REMOVED***

// LookupElement searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Element if they key does not exist. To know if they key actually
// exists, use LookupElementErr.
func (d Doc) LookupElement(key ...string) Elem ***REMOVED***
	elem, _ := d.LookupElementErr(key...)
	return elem
***REMOVED***

// LookupElementErr searches the document and potentially subdocuments for the
// provided key. Each key provided to this method represents a layer of depth.
func (d Doc) LookupElementErr(key ...string) (Elem, error) ***REMOVED***
	// KeyNotFound operates by being created where the error happens and then the depth is
	// incremented by 1 as each function unwinds. Whenever this function returns, it also assigns
	// the Key slice to the key slice it has. This ensures that the proper depth is identified and
	// the proper keys.
	if len(key) == 0 ***REMOVED***
		return Elem***REMOVED******REMOVED***, KeyNotFound***REMOVED***Key: key***REMOVED***
	***REMOVED***

	var elem Elem
	var err error
	idx := d.IndexOf(key[0])
	if idx == -1 ***REMOVED***
		return Elem***REMOVED******REMOVED***, KeyNotFound***REMOVED***Key: key***REMOVED***
	***REMOVED***

	elem = d[idx]
	if len(key) == 1 ***REMOVED***
		return elem, nil
	***REMOVED***

	switch elem.Value.Type() ***REMOVED***
	case bsontype.EmbeddedDocument:
		switch tt := elem.Value.primitive.(type) ***REMOVED***
		case Doc:
			elem, err = tt.LookupElementErr(key[1:]...)
		case MDoc:
			elem, err = tt.LookupElementErr(key[1:]...)
		***REMOVED***
	default:
		return Elem***REMOVED******REMOVED***, KeyNotFound***REMOVED***Type: elem.Value.Type()***REMOVED***
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
func (d Doc) MarshalBSONValue() (bsontype.Type, []byte, error) ***REMOVED***
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
func (d Doc) MarshalBSON() ([]byte, error) ***REMOVED*** return d.AppendMarshalBSON(nil) ***REMOVED***

// AppendMarshalBSON marshals Doc to BSON bytes, appending to dst.
//
// This method will never return an error.
func (d Doc) AppendMarshalBSON(dst []byte) ([]byte, error) ***REMOVED***
	idx, dst := bsoncore.ReserveLength(dst)
	for _, elem := range d ***REMOVED***
		t, data, _ := elem.Value.MarshalBSONValue() // Value.MarshalBSONValue never returns an error.
		dst = append(dst, byte(t))
		dst = append(dst, elem.Key...)
		dst = append(dst, 0x00)
		dst = append(dst, data...)
	***REMOVED***
	dst = append(dst, 0x00)
	dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst, nil
***REMOVED***

// UnmarshalBSON implements the Unmarshaler interface.
func (d *Doc) UnmarshalBSON(b []byte) error ***REMOVED***
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
		*d = d.Append(elem.Key(), val)
	***REMOVED***
	return nil
***REMOVED***

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (d *Doc) UnmarshalBSONValue(t bsontype.Type, data []byte) error ***REMOVED***
	if t != bsontype.EmbeddedDocument ***REMOVED***
		return fmt.Errorf("cannot unmarshal %s into a bsonx.Doc", t)
	***REMOVED***
	return d.UnmarshalBSON(data)
***REMOVED***

// Equal compares this document to another, returning true if they are equal.
func (d Doc) Equal(id IDoc) bool ***REMOVED***
	switch tt := id.(type) ***REMOVED***
	case Doc:
		d2 := tt
		if len(d) != len(d2) ***REMOVED***
			return false
		***REMOVED***
		for idx := range d ***REMOVED***
			if !d[idx].Equal(d2[idx]) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	case MDoc:
		unique := make(map[string]struct***REMOVED******REMOVED***)
		for _, elem := range d ***REMOVED***
			unique[elem.Key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			val, ok := tt[elem.Key]
			if !ok ***REMOVED***
				return false
			***REMOVED***
			if !val.Equal(elem.Value) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if len(unique) != len(tt) ***REMOVED***
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
func (d Doc) String() string ***REMOVED***
	var buf bytes.Buffer
	buf.Write([]byte("bson.Document***REMOVED***"))
	for idx, elem := range d ***REMOVED***
		if idx > 0 ***REMOVED***
			buf.Write([]byte(", "))
		***REMOVED***
		fmt.Fprintf(&buf, "%v", elem)
	***REMOVED***
	buf.WriteByte('***REMOVED***')

	return buf.String()
***REMOVED***

func (Doc) idoc() ***REMOVED******REMOVED***
