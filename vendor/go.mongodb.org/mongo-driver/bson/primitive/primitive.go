// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package primitive contains types similar to Go primitives for BSON types that do not have direct
// Go primitive representations.
package primitive // import "go.mongodb.org/mongo-driver/bson/primitive"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Binary represents a BSON binary value.
type Binary struct ***REMOVED***
	Subtype byte
	Data    []byte
***REMOVED***

// Equal compares bp to bp2 and returns true if they are equal.
func (bp Binary) Equal(bp2 Binary) bool ***REMOVED***
	if bp.Subtype != bp2.Subtype ***REMOVED***
		return false
	***REMOVED***
	return bytes.Equal(bp.Data, bp2.Data)
***REMOVED***

// IsZero returns if bp is the empty Binary.
func (bp Binary) IsZero() bool ***REMOVED***
	return bp.Subtype == 0 && len(bp.Data) == 0
***REMOVED***

// Undefined represents the BSON undefined value type.
type Undefined struct***REMOVED******REMOVED***

// DateTime represents the BSON datetime value.
type DateTime int64

var _ json.Marshaler = DateTime(0)
var _ json.Unmarshaler = (*DateTime)(nil)

// MarshalJSON marshal to time type.
func (d DateTime) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(d.Time())
***REMOVED***

// UnmarshalJSON creates a primitive.DateTime from a JSON string.
func (d *DateTime) UnmarshalJSON(data []byte) error ***REMOVED***
	// Ignore "null" to keep parity with the time.Time type and the standard library. Decoding "null" into a non-pointer
	// DateTime field will leave the field unchanged. For pointer values, the encoding/json will set the pointer to nil
	// and will not defer to the UnmarshalJSON hook.
	if string(data) == "null" ***REMOVED***
		return nil
	***REMOVED***

	var tempTime time.Time
	if err := json.Unmarshal(data, &tempTime); err != nil ***REMOVED***
		return err
	***REMOVED***

	*d = NewDateTimeFromTime(tempTime)
	return nil
***REMOVED***

// Time returns the date as a time type.
func (d DateTime) Time() time.Time ***REMOVED***
	return time.Unix(int64(d)/1000, int64(d)%1000*1000000)
***REMOVED***

// NewDateTimeFromTime creates a new DateTime from a Time.
func NewDateTimeFromTime(t time.Time) DateTime ***REMOVED***
	return DateTime(t.Unix()*1e3 + int64(t.Nanosecond())/1e6)
***REMOVED***

// Null represents the BSON null value.
type Null struct***REMOVED******REMOVED***

// Regex represents a BSON regex value.
type Regex struct ***REMOVED***
	Pattern string
	Options string
***REMOVED***

func (rp Regex) String() string ***REMOVED***
	return fmt.Sprintf(`***REMOVED***"pattern": "%s", "options": "%s"***REMOVED***`, rp.Pattern, rp.Options)
***REMOVED***

// Equal compares rp to rp2 and returns true if they are equal.
func (rp Regex) Equal(rp2 Regex) bool ***REMOVED***
	return rp.Pattern == rp2.Pattern && rp.Options == rp2.Options
***REMOVED***

// IsZero returns if rp is the empty Regex.
func (rp Regex) IsZero() bool ***REMOVED***
	return rp.Pattern == "" && rp.Options == ""
***REMOVED***

// DBPointer represents a BSON dbpointer value.
type DBPointer struct ***REMOVED***
	DB      string
	Pointer ObjectID
***REMOVED***

func (d DBPointer) String() string ***REMOVED***
	return fmt.Sprintf(`***REMOVED***"db": "%s", "pointer": "%s"***REMOVED***`, d.DB, d.Pointer)
***REMOVED***

// Equal compares d to d2 and returns true if they are equal.
func (d DBPointer) Equal(d2 DBPointer) bool ***REMOVED***
	return d == d2
***REMOVED***

// IsZero returns if d is the empty DBPointer.
func (d DBPointer) IsZero() bool ***REMOVED***
	return d.DB == "" && d.Pointer.IsZero()
***REMOVED***

// JavaScript represents a BSON JavaScript code value.
type JavaScript string

// Symbol represents a BSON symbol value.
type Symbol string

// CodeWithScope represents a BSON JavaScript code with scope value.
type CodeWithScope struct ***REMOVED***
	Code  JavaScript
	Scope interface***REMOVED******REMOVED***
***REMOVED***

func (cws CodeWithScope) String() string ***REMOVED***
	return fmt.Sprintf(`***REMOVED***"code": "%s", "scope": %v***REMOVED***`, cws.Code, cws.Scope)
***REMOVED***

// Timestamp represents a BSON timestamp value.
type Timestamp struct ***REMOVED***
	T uint32
	I uint32
***REMOVED***

// Equal compares tp to tp2 and returns true if they are equal.
func (tp Timestamp) Equal(tp2 Timestamp) bool ***REMOVED***
	return tp.T == tp2.T && tp.I == tp2.I
***REMOVED***

// IsZero returns if tp is the zero Timestamp.
func (tp Timestamp) IsZero() bool ***REMOVED***
	return tp.T == 0 && tp.I == 0
***REMOVED***

// CompareTimestamp returns an integer comparing two Timestamps, where T is compared first, followed by I.
// Returns 0 if tp = tp2, 1 if tp > tp2, -1 if tp < tp2.
func CompareTimestamp(tp, tp2 Timestamp) int ***REMOVED***
	if tp.Equal(tp2) ***REMOVED***
		return 0
	***REMOVED***

	if tp.T > tp2.T ***REMOVED***
		return 1
	***REMOVED***
	if tp.T < tp2.T ***REMOVED***
		return -1
	***REMOVED***
	// Compare I values because T values are equal
	if tp.I > tp2.I ***REMOVED***
		return 1
	***REMOVED***
	return -1
***REMOVED***

// MinKey represents the BSON minkey value.
type MinKey struct***REMOVED******REMOVED***

// MaxKey represents the BSON maxkey value.
type MaxKey struct***REMOVED******REMOVED***

// D is an ordered representation of a BSON document. This type should be used when the order of the elements matters,
// such as MongoDB command documents. If the order of the elements does not matter, an M should be used instead.
//
// Example usage:
//
// 		bson.D***REMOVED******REMOVED***"foo", "bar"***REMOVED***, ***REMOVED***"hello", "world"***REMOVED***, ***REMOVED***"pi", 3.14159***REMOVED******REMOVED***
type D []E

// Map creates a map from the elements of the D.
func (d D) Map() M ***REMOVED***
	m := make(M, len(d))
	for _, e := range d ***REMOVED***
		m[e.Key] = e.Value
	***REMOVED***
	return m
***REMOVED***

// E represents a BSON element for a D. It is usually used inside a D.
type E struct ***REMOVED***
	Key   string
	Value interface***REMOVED******REMOVED***
***REMOVED***

// M is an unordered representation of a BSON document. This type should be used when the order of the elements does not
// matter. This type is handled as a regular map[string]interface***REMOVED******REMOVED*** when encoding and decoding. Elements will be
// serialized in an undefined, random order. If the order of the elements matters, a D should be used instead.
//
// Example usage:
//
// 		bson.M***REMOVED***"foo": "bar", "hello": "world", "pi": 3.14159***REMOVED***
type M map[string]interface***REMOVED******REMOVED***

// An A is an ordered representation of a BSON array.
//
// Example usage:
//
// 		bson.A***REMOVED***"bar", "world", 3.14159, bson.D***REMOVED******REMOVED***"qux", 12345***REMOVED******REMOVED******REMOVED***
type A []interface***REMOVED******REMOVED***
