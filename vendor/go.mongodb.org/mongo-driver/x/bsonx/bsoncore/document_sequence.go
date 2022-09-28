// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncore

import (
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// DocumentSequenceStyle is used to represent how a document sequence is laid out in a slice of
// bytes.
type DocumentSequenceStyle uint32

// These constants are the valid styles for a DocumentSequence.
const (
	_ DocumentSequenceStyle = iota
	SequenceStyle
	ArrayStyle
)

// DocumentSequence represents a sequence of documents. The Style field indicates how the documents
// are laid out inside of the Data field.
type DocumentSequence struct ***REMOVED***
	Style DocumentSequenceStyle
	Data  []byte
	Pos   int
***REMOVED***

// ErrCorruptedDocument is returned when a full document couldn't be read from the sequence.
var ErrCorruptedDocument = errors.New("invalid DocumentSequence: corrupted document")

// ErrNonDocument is returned when a DocumentSequence contains a non-document BSON value.
var ErrNonDocument = errors.New("invalid DocumentSequence: a non-document value was found in sequence")

// ErrInvalidDocumentSequenceStyle is returned when an unknown DocumentSequenceStyle is set on a
// DocumentSequence.
var ErrInvalidDocumentSequenceStyle = errors.New("invalid DocumentSequenceStyle")

// DocumentCount returns the number of documents in the sequence.
func (ds *DocumentSequence) DocumentCount() int ***REMOVED***
	if ds == nil ***REMOVED***
		return 0
	***REMOVED***
	switch ds.Style ***REMOVED***
	case SequenceStyle:
		var count int
		var ok bool
		rem := ds.Data
		for len(rem) > 0 ***REMOVED***
			_, rem, ok = ReadDocument(rem)
			if !ok ***REMOVED***
				return 0
			***REMOVED***
			count++
		***REMOVED***
		return count
	case ArrayStyle:
		_, rem, ok := ReadLength(ds.Data)
		if !ok ***REMOVED***
			return 0
		***REMOVED***

		var count int
		for len(rem) > 1 ***REMOVED***
			_, rem, ok = ReadElement(rem)
			if !ok ***REMOVED***
				return 0
			***REMOVED***
			count++
		***REMOVED***
		return count
	default:
		return 0
	***REMOVED***
***REMOVED***

// Empty returns true if the sequence is empty. It always returns true for unknown sequence styles.
func (ds *DocumentSequence) Empty() bool ***REMOVED***
	if ds == nil ***REMOVED***
		return true
	***REMOVED***

	switch ds.Style ***REMOVED***
	case SequenceStyle:
		return len(ds.Data) == 0
	case ArrayStyle:
		return len(ds.Data) <= 5
	default:
		return true
	***REMOVED***
***REMOVED***

//ResetIterator resets the iteration point for the Next method to the beginning of the document
//sequence.
func (ds *DocumentSequence) ResetIterator() ***REMOVED***
	if ds == nil ***REMOVED***
		return
	***REMOVED***
	ds.Pos = 0
***REMOVED***

// Documents returns a slice of the documents. If nil either the Data field is also nil or could not
// be properly read.
func (ds *DocumentSequence) Documents() ([]Document, error) ***REMOVED***
	if ds == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	switch ds.Style ***REMOVED***
	case SequenceStyle:
		rem := ds.Data
		var docs []Document
		var doc Document
		var ok bool
		for ***REMOVED***
			doc, rem, ok = ReadDocument(rem)
			if !ok ***REMOVED***
				if len(rem) == 0 ***REMOVED***
					break
				***REMOVED***
				return nil, ErrCorruptedDocument
			***REMOVED***
			docs = append(docs, doc)
		***REMOVED***
		return docs, nil
	case ArrayStyle:
		if len(ds.Data) == 0 ***REMOVED***
			return nil, nil
		***REMOVED***
		vals, err := Document(ds.Data).Values()
		if err != nil ***REMOVED***
			return nil, ErrCorruptedDocument
		***REMOVED***
		docs := make([]Document, 0, len(vals))
		for _, v := range vals ***REMOVED***
			if v.Type != bsontype.EmbeddedDocument ***REMOVED***
				return nil, ErrNonDocument
			***REMOVED***
			docs = append(docs, v.Data)
		***REMOVED***
		return docs, nil
	default:
		return nil, ErrInvalidDocumentSequenceStyle
	***REMOVED***
***REMOVED***

// Next retrieves the next document from this sequence and returns it. This method will return
// io.EOF when it has reached the end of the sequence.
func (ds *DocumentSequence) Next() (Document, error) ***REMOVED***
	if ds == nil || ds.Pos >= len(ds.Data) ***REMOVED***
		return nil, io.EOF
	***REMOVED***
	switch ds.Style ***REMOVED***
	case SequenceStyle:
		doc, _, ok := ReadDocument(ds.Data[ds.Pos:])
		if !ok ***REMOVED***
			return nil, ErrCorruptedDocument
		***REMOVED***
		ds.Pos += len(doc)
		return doc, nil
	case ArrayStyle:
		if ds.Pos < 4 ***REMOVED***
			if len(ds.Data) < 4 ***REMOVED***
				return nil, ErrCorruptedDocument
			***REMOVED***
			ds.Pos = 4 // Skip the length of the document
		***REMOVED***
		if len(ds.Data[ds.Pos:]) == 1 && ds.Data[ds.Pos] == 0x00 ***REMOVED***
			return nil, io.EOF // At the end of the document
		***REMOVED***
		elem, _, ok := ReadElement(ds.Data[ds.Pos:])
		if !ok ***REMOVED***
			return nil, ErrCorruptedDocument
		***REMOVED***
		ds.Pos += len(elem)
		val := elem.Value()
		if val.Type != bsontype.EmbeddedDocument ***REMOVED***
			return nil, ErrNonDocument
		***REMOVED***
		return val.Data, nil
	default:
		return nil, ErrInvalidDocumentSequenceStyle
	***REMOVED***
***REMOVED***
