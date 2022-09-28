// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on gopkg.in/mgo.v2/bson by Gustavo Niemeyer
// See THIRD-PARTY-NOTICES for original license terms.

package primitive

import (
	"crypto/rand"
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// ErrInvalidHex indicates that a hex string cannot be converted to an ObjectID.
var ErrInvalidHex = errors.New("the provided hex string is not a valid ObjectID")

// ObjectID is the BSON ObjectID type.
type ObjectID [12]byte

// NilObjectID is the zero value for ObjectID.
var NilObjectID ObjectID

var objectIDCounter = readRandomUint32()
var processUnique = processUniqueBytes()

var _ encoding.TextMarshaler = ObjectID***REMOVED******REMOVED***
var _ encoding.TextUnmarshaler = &ObjectID***REMOVED******REMOVED***

// NewObjectID generates a new ObjectID.
func NewObjectID() ObjectID ***REMOVED***
	return NewObjectIDFromTimestamp(time.Now())
***REMOVED***

// NewObjectIDFromTimestamp generates a new ObjectID based on the given time.
func NewObjectIDFromTimestamp(timestamp time.Time) ObjectID ***REMOVED***
	var b [12]byte

	binary.BigEndian.PutUint32(b[0:4], uint32(timestamp.Unix()))
	copy(b[4:9], processUnique[:])
	putUint24(b[9:12], atomic.AddUint32(&objectIDCounter, 1))

	return b
***REMOVED***

// Timestamp extracts the time part of the ObjectId.
func (id ObjectID) Timestamp() time.Time ***REMOVED***
	unixSecs := binary.BigEndian.Uint32(id[0:4])
	return time.Unix(int64(unixSecs), 0).UTC()
***REMOVED***

// Hex returns the hex encoding of the ObjectID as a string.
func (id ObjectID) Hex() string ***REMOVED***
	return hex.EncodeToString(id[:])
***REMOVED***

func (id ObjectID) String() string ***REMOVED***
	return fmt.Sprintf("ObjectID(%q)", id.Hex())
***REMOVED***

// IsZero returns true if id is the empty ObjectID.
func (id ObjectID) IsZero() bool ***REMOVED***
	return id == NilObjectID
***REMOVED***

// ObjectIDFromHex creates a new ObjectID from a hex string. It returns an error if the hex string is not a
// valid ObjectID.
func ObjectIDFromHex(s string) (ObjectID, error) ***REMOVED***
	if len(s) != 24 ***REMOVED***
		return NilObjectID, ErrInvalidHex
	***REMOVED***

	b, err := hex.DecodeString(s)
	if err != nil ***REMOVED***
		return NilObjectID, err
	***REMOVED***

	var oid [12]byte
	copy(oid[:], b)

	return oid, nil
***REMOVED***

// IsValidObjectID returns true if the provided hex string represents a valid ObjectID and false if not.
func IsValidObjectID(s string) bool ***REMOVED***
	_, err := ObjectIDFromHex(s)
	return err == nil
***REMOVED***

// MarshalText returns the ObjectID as UTF-8-encoded text. Implementing this allows us to use ObjectID
// as a map key when marshalling JSON. See https://pkg.go.dev/encoding#TextMarshaler
func (id ObjectID) MarshalText() ([]byte, error) ***REMOVED***
	return []byte(id.Hex()), nil
***REMOVED***

// UnmarshalText populates the byte slice with the ObjectID. Implementing this allows us to use ObjectID
// as a map key when unmarshalling JSON. See https://pkg.go.dev/encoding#TextUnmarshaler
func (id *ObjectID) UnmarshalText(b []byte) error ***REMOVED***
	oid, err := ObjectIDFromHex(string(b))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*id = oid
	return nil
***REMOVED***

// MarshalJSON returns the ObjectID as a string
func (id ObjectID) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(id.Hex())
***REMOVED***

// UnmarshalJSON populates the byte slice with the ObjectID. If the byte slice is 24 bytes long, it
// will be populated with the hex representation of the ObjectID. If the byte slice is twelve bytes
// long, it will be populated with the BSON representation of the ObjectID. This method also accepts empty strings and
// decodes them as NilObjectID. For any other inputs, an error will be returned.
func (id *ObjectID) UnmarshalJSON(b []byte) error ***REMOVED***
	// Ignore "null" to keep parity with the standard library. Decoding a JSON null into a non-pointer ObjectID field
	// will leave the field unchanged. For pointer values, encoding/json will set the pointer to nil and will not
	// enter the UnmarshalJSON hook.
	if string(b) == "null" ***REMOVED***
		return nil
	***REMOVED***

	var err error
	switch len(b) ***REMOVED***
	case 12:
		copy(id[:], b)
	default:
		// Extended JSON
		var res interface***REMOVED******REMOVED***
		err := json.Unmarshal(b, &res)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		str, ok := res.(string)
		if !ok ***REMOVED***
			m, ok := res.(map[string]interface***REMOVED******REMOVED***)
			if !ok ***REMOVED***
				return errors.New("not an extended JSON ObjectID")
			***REMOVED***
			oid, ok := m["$oid"]
			if !ok ***REMOVED***
				return errors.New("not an extended JSON ObjectID")
			***REMOVED***
			str, ok = oid.(string)
			if !ok ***REMOVED***
				return errors.New("not an extended JSON ObjectID")
			***REMOVED***
		***REMOVED***

		// An empty string is not a valid ObjectID, but we treat it as a special value that decodes as NilObjectID.
		if len(str) == 0 ***REMOVED***
			copy(id[:], NilObjectID[:])
			return nil
		***REMOVED***

		if len(str) != 24 ***REMOVED***
			return fmt.Errorf("cannot unmarshal into an ObjectID, the length must be 24 but it is %d", len(str))
		***REMOVED***

		_, err = hex.Decode(id[:], []byte(str))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

func processUniqueBytes() [5]byte ***REMOVED***
	var b [5]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil ***REMOVED***
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %v", err))
	***REMOVED***

	return b
***REMOVED***

func readRandomUint32() uint32 ***REMOVED***
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil ***REMOVED***
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %v", err))
	***REMOVED***

	return (uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
***REMOVED***

func putUint24(b []byte, v uint32) ***REMOVED***
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
***REMOVED***
