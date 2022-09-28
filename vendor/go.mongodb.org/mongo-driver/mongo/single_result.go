// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

// ErrNoDocuments is returned by SingleResult methods when the operation that created the SingleResult did not return
// any documents.
var ErrNoDocuments = errors.New("mongo: no documents in result")

// SingleResult represents a single document returned from an operation. If the operation resulted in an error, all
// SingleResult methods will return that error. If the operation did not return any documents, all SingleResult methods
// will return ErrNoDocuments.
type SingleResult struct ***REMOVED***
	err error
	cur *Cursor
	rdr bson.Raw
	reg *bsoncodec.Registry
***REMOVED***

// NewSingleResultFromDocument creates a SingleResult with the provided error, registry, and an underlying Cursor pre-loaded with
// the provided document, error and registry. If no registry is provided, bson.DefaultRegistry will be used. If an error distinct
// from the one provided occurs during creation of the SingleResult, that error will be stored on the returned SingleResult.
//
// The document parameter must be a non-nil document.
func NewSingleResultFromDocument(document interface***REMOVED******REMOVED***, err error, registry *bsoncodec.Registry) *SingleResult ***REMOVED***
	if document == nil ***REMOVED***
		return &SingleResult***REMOVED***err: ErrNilDocument***REMOVED***
	***REMOVED***
	if registry == nil ***REMOVED***
		registry = bson.DefaultRegistry
	***REMOVED***

	cur, createErr := NewCursorFromDocuments([]interface***REMOVED******REMOVED******REMOVED***document***REMOVED***, err, registry)
	if createErr != nil ***REMOVED***
		return &SingleResult***REMOVED***err: createErr***REMOVED***
	***REMOVED***

	return &SingleResult***REMOVED***
		cur: cur,
		err: err,
		reg: registry,
	***REMOVED***
***REMOVED***

// Decode will unmarshal the document represented by this SingleResult into v. If there was an error from the operation
// that created this SingleResult, that error will be returned. If the operation returned no documents, Decode will
// return ErrNoDocuments.
//
// If the operation was successful and returned a document, Decode will return any errors from the unmarshalling process
// without any modification. If v is nil or is a typed nil, an error will be returned.
func (sr *SingleResult) Decode(v interface***REMOVED******REMOVED***) error ***REMOVED***
	if sr.err != nil ***REMOVED***
		return sr.err
	***REMOVED***
	if sr.reg == nil ***REMOVED***
		return bson.ErrNilRegistry
	***REMOVED***

	if sr.err = sr.setRdrContents(); sr.err != nil ***REMOVED***
		return sr.err
	***REMOVED***
	return bson.UnmarshalWithRegistry(sr.reg, sr.rdr, v)
***REMOVED***

// DecodeBytes will return the document represented by this SingleResult as a bson.Raw. If there was an error from the
// operation that created this SingleResult, both the result and that error will be returned. If the operation returned
// no documents, this will return (nil, ErrNoDocuments).
func (sr *SingleResult) DecodeBytes() (bson.Raw, error) ***REMOVED***
	if sr.err != nil ***REMOVED***
		return sr.rdr, sr.err
	***REMOVED***

	if sr.err = sr.setRdrContents(); sr.err != nil ***REMOVED***
		return nil, sr.err
	***REMOVED***
	return sr.rdr, nil
***REMOVED***

// setRdrContents will set the contents of rdr by iterating the underlying cursor if necessary.
func (sr *SingleResult) setRdrContents() error ***REMOVED***
	switch ***REMOVED***
	case sr.err != nil:
		return sr.err
	case sr.rdr != nil:
		return nil
	case sr.cur != nil:
		defer sr.cur.Close(context.TODO())

		if !sr.cur.Next(context.TODO()) ***REMOVED***
			if err := sr.cur.Err(); err != nil ***REMOVED***
				return err
			***REMOVED***

			return ErrNoDocuments
		***REMOVED***
		sr.rdr = sr.cur.Current
		return nil
	***REMOVED***

	return ErrNoDocuments
***REMOVED***

// Err returns the error from the operation that created this SingleResult. If the operation was successful but did not
// return any documents, Err will return ErrNoDocuments. If the operation was successful and returned a document, Err
// will return nil.
func (sr *SingleResult) Err() error ***REMOVED***
	sr.err = sr.setRdrContents()

	return sr.err
***REMOVED***
