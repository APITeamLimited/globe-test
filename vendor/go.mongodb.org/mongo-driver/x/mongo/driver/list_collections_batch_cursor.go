// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"errors"
	"io"
	"strings"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ListCollectionsBatchCursor is a special batch cursor returned from ListCollections that properly
// handles current and legacy ListCollections operations.
type ListCollectionsBatchCursor struct ***REMOVED***
	legacy       bool // server version < 3.0
	bc           *BatchCursor
	currentBatch *bsoncore.DocumentSequence
	err          error
***REMOVED***

// NewListCollectionsBatchCursor creates a new non-legacy ListCollectionsCursor.
func NewListCollectionsBatchCursor(bc *BatchCursor) (*ListCollectionsBatchCursor, error) ***REMOVED***
	if bc == nil ***REMOVED***
		return nil, errors.New("batch cursor must not be nil")
	***REMOVED***
	return &ListCollectionsBatchCursor***REMOVED***bc: bc, currentBatch: new(bsoncore.DocumentSequence)***REMOVED***, nil
***REMOVED***

// NewLegacyListCollectionsBatchCursor creates a new legacy ListCollectionsCursor.
func NewLegacyListCollectionsBatchCursor(bc *BatchCursor) (*ListCollectionsBatchCursor, error) ***REMOVED***
	if bc == nil ***REMOVED***
		return nil, errors.New("batch cursor must not be nil")
	***REMOVED***
	return &ListCollectionsBatchCursor***REMOVED***legacy: true, bc: bc, currentBatch: new(bsoncore.DocumentSequence)***REMOVED***, nil
***REMOVED***

// ID returns the cursor ID for this batch cursor.
func (lcbc *ListCollectionsBatchCursor) ID() int64 ***REMOVED***
	return lcbc.bc.ID()
***REMOVED***

// Next indicates if there is another batch available. Returning false does not necessarily indicate
// that the cursor is closed. This method will return false when an empty batch is returned.
//
// If Next returns true, there is a valid batch of documents available. If Next returns false, there
// is not a valid batch of documents available.
func (lcbc *ListCollectionsBatchCursor) Next(ctx context.Context) bool ***REMOVED***
	if !lcbc.bc.Next(ctx) ***REMOVED***
		return false
	***REMOVED***

	if !lcbc.legacy ***REMOVED***
		lcbc.currentBatch.Style = lcbc.bc.currentBatch.Style
		lcbc.currentBatch.Data = lcbc.bc.currentBatch.Data
		lcbc.currentBatch.ResetIterator()
		return true
	***REMOVED***

	lcbc.currentBatch.Style = bsoncore.SequenceStyle
	lcbc.currentBatch.Data = lcbc.currentBatch.Data[:0]

	var doc bsoncore.Document
	for ***REMOVED***
		doc, lcbc.err = lcbc.bc.currentBatch.Next()
		if lcbc.err != nil ***REMOVED***
			if lcbc.err == io.EOF ***REMOVED***
				lcbc.err = nil
				break
			***REMOVED***
			return false
		***REMOVED***
		doc, lcbc.err = lcbc.projectNameElement(doc)
		if lcbc.err != nil ***REMOVED***
			return false
		***REMOVED***
		lcbc.currentBatch.Data = append(lcbc.currentBatch.Data, doc...)
	***REMOVED***

	return true
***REMOVED***

// Batch will return a DocumentSequence for the current batch of documents. The returned
// DocumentSequence is only valid until the next call to Next or Close.
func (lcbc *ListCollectionsBatchCursor) Batch() *bsoncore.DocumentSequence ***REMOVED*** return lcbc.currentBatch ***REMOVED***

// Server returns a pointer to the cursor's server.
func (lcbc *ListCollectionsBatchCursor) Server() Server ***REMOVED*** return lcbc.bc.server ***REMOVED***

// Err returns the latest error encountered.
func (lcbc *ListCollectionsBatchCursor) Err() error ***REMOVED***
	if lcbc.err != nil ***REMOVED***
		return lcbc.err
	***REMOVED***
	return lcbc.bc.Err()
***REMOVED***

// Close closes this batch cursor.
func (lcbc *ListCollectionsBatchCursor) Close(ctx context.Context) error ***REMOVED*** return lcbc.bc.Close(ctx) ***REMOVED***

// project out the database name for a legacy server
func (*ListCollectionsBatchCursor) projectNameElement(rawDoc bsoncore.Document) (bsoncore.Document, error) ***REMOVED***
	elems, err := rawDoc.Elements()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var filteredElems []byte
	for _, elem := range elems ***REMOVED***
		key := elem.Key()
		if key != "name" ***REMOVED***
			filteredElems = append(filteredElems, elem...)
			continue
		***REMOVED***

		name := elem.Value().StringValue()
		collName := name[strings.Index(name, ".")+1:]
		filteredElems = bsoncore.AppendStringElement(filteredElems, "name", collName)
	***REMOVED***

	var filteredDoc []byte
	filteredDoc = bsoncore.BuildDocument(filteredDoc, filteredElems)
	return filteredDoc, nil
***REMOVED***
