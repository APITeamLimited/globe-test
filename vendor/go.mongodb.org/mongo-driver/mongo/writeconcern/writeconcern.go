// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package writeconcern defines write concerns for MongoDB operations.
package writeconcern // import "go.mongodb.org/mongo-driver/mongo/writeconcern"

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrInconsistent indicates that an inconsistent write concern was specified.
var ErrInconsistent = errors.New("a write concern cannot have both w=0 and j=true")

// ErrEmptyWriteConcern indicates that a write concern has no fields set.
var ErrEmptyWriteConcern = errors.New("a write concern must have at least one field set")

// ErrNegativeW indicates that a negative integer `w` field was specified.
var ErrNegativeW = errors.New("write concern `w` field cannot be a negative number")

// ErrNegativeWTimeout indicates that a negative WTimeout was specified.
var ErrNegativeWTimeout = errors.New("write concern `wtimeout` field cannot be negative")

// WriteConcern describes the level of acknowledgement requested from MongoDB for write operations
// to a standalone mongod or to replica sets or to sharded clusters.
type WriteConcern struct ***REMOVED***
	w interface***REMOVED******REMOVED***
	j bool

	// NOTE(benjirewis): wTimeout will be deprecated in a future release. The more general Timeout
	// option may be used in its place to control the amount of time that a single operation can run
	// before returning an error. Using wTimeout and setting Timeout on the client will result in
	// undefined behavior.
	wTimeout time.Duration
***REMOVED***

// Option is an option to provide when creating a WriteConcern.
type Option func(concern *WriteConcern)

// New constructs a new WriteConcern.
func New(options ...Option) *WriteConcern ***REMOVED***
	concern := &WriteConcern***REMOVED******REMOVED***

	for _, option := range options ***REMOVED***
		option(concern)
	***REMOVED***

	return concern
***REMOVED***

// W requests acknowledgement that write operations propagate to the specified number of mongod
// instances.
func W(w int) Option ***REMOVED***
	return func(concern *WriteConcern) ***REMOVED***
		concern.w = w
	***REMOVED***
***REMOVED***

// WMajority requests acknowledgement that write operations propagate to the majority of mongod
// instances.
func WMajority() Option ***REMOVED***
	return func(concern *WriteConcern) ***REMOVED***
		concern.w = "majority"
	***REMOVED***
***REMOVED***

// WTagSet requests acknowledgement that write operations propagate to the specified mongod
// instance.
func WTagSet(tag string) Option ***REMOVED***
	return func(concern *WriteConcern) ***REMOVED***
		concern.w = tag
	***REMOVED***
***REMOVED***

// J requests acknowledgement from MongoDB that write operations are written to
// the journal.
func J(j bool) Option ***REMOVED***
	return func(concern *WriteConcern) ***REMOVED***
		concern.j = j
	***REMOVED***
***REMOVED***

// WTimeout specifies specifies a time limit for the write concern.
//
// NOTE(benjirewis): wTimeout will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can run
// before returning an error. Using wTimeout and setting Timeout on the client will result in
// undefined behavior.
func WTimeout(d time.Duration) Option ***REMOVED***
	return func(concern *WriteConcern) ***REMOVED***
		concern.wTimeout = d
	***REMOVED***
***REMOVED***

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (wc *WriteConcern) MarshalBSONValue() (bsontype.Type, []byte, error) ***REMOVED***
	if !wc.IsValid() ***REMOVED***
		return bsontype.Type(0), nil, ErrInconsistent
	***REMOVED***

	var elems []byte

	if wc.w != nil ***REMOVED***
		switch t := wc.w.(type) ***REMOVED***
		case int:
			if t < 0 ***REMOVED***
				return bsontype.Type(0), nil, ErrNegativeW
			***REMOVED***

			elems = bsoncore.AppendInt32Element(elems, "w", int32(t))
		case string:
			elems = bsoncore.AppendStringElement(elems, "w", t)
		***REMOVED***
	***REMOVED***

	if wc.j ***REMOVED***
		elems = bsoncore.AppendBooleanElement(elems, "j", wc.j)
	***REMOVED***

	if wc.wTimeout < 0 ***REMOVED***
		return bsontype.Type(0), nil, ErrNegativeWTimeout
	***REMOVED***

	if wc.wTimeout != 0 ***REMOVED***
		elems = bsoncore.AppendInt64Element(elems, "wtimeout", int64(wc.wTimeout/time.Millisecond))
	***REMOVED***

	if len(elems) == 0 ***REMOVED***
		return bsontype.Type(0), nil, ErrEmptyWriteConcern
	***REMOVED***
	return bsontype.EmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
***REMOVED***

// AcknowledgedValue returns true if a BSON RawValue for a write concern represents an acknowledged write concern.
// The element's value must be a document representing a write concern.
func AcknowledgedValue(rawv bson.RawValue) bool ***REMOVED***
	doc, ok := bsoncore.Value***REMOVED***Type: rawv.Type, Data: rawv.Value***REMOVED***.DocumentOK()
	if !ok ***REMOVED***
		return false
	***REMOVED***

	val, err := doc.LookupErr("w")
	if err != nil ***REMOVED***
		// key w not found --> acknowledged
		return true
	***REMOVED***

	i32, ok := val.Int32OK()
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return i32 != 0
***REMOVED***

// Acknowledged indicates whether or not a write with the given write concern will be acknowledged.
func (wc *WriteConcern) Acknowledged() bool ***REMOVED***
	if wc == nil || wc.j ***REMOVED***
		return true
	***REMOVED***

	switch v := wc.w.(type) ***REMOVED***
	case int:
		if v == 0 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// IsValid checks whether the write concern is invalid.
func (wc *WriteConcern) IsValid() bool ***REMOVED***
	if !wc.j ***REMOVED***
		return true
	***REMOVED***

	switch v := wc.w.(type) ***REMOVED***
	case int:
		if v == 0 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// GetW returns the write concern w level.
func (wc *WriteConcern) GetW() interface***REMOVED******REMOVED*** ***REMOVED***
	return wc.w
***REMOVED***

// GetJ returns the write concern journaling level.
func (wc *WriteConcern) GetJ() bool ***REMOVED***
	return wc.j
***REMOVED***

// GetWTimeout returns the write concern timeout.
func (wc *WriteConcern) GetWTimeout() time.Duration ***REMOVED***
	return wc.wTimeout
***REMOVED***

// WithOptions returns a copy of this WriteConcern with the options set.
func (wc *WriteConcern) WithOptions(options ...Option) *WriteConcern ***REMOVED***
	if wc == nil ***REMOVED***
		return New(options...)
	***REMOVED***
	newWC := &WriteConcern***REMOVED******REMOVED***
	*newWC = *wc

	for _, option := range options ***REMOVED***
		option(newWC)
	***REMOVED***

	return newWC
***REMOVED***

// AckWrite returns true if a write concern represents an acknowledged write
func AckWrite(wc *WriteConcern) bool ***REMOVED***
	return wc == nil || wc.Acknowledged()
***REMOVED***
