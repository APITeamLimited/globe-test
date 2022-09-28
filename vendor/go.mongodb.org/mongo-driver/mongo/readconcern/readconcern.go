// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package readconcern defines read concerns for MongoDB operations.
package readconcern // import "go.mongodb.org/mongo-driver/mongo/readconcern"

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ReadConcern for replica sets and replica set shards determines which data to return from a query.
type ReadConcern struct ***REMOVED***
	level string
***REMOVED***

// Option is an option to provide when creating a ReadConcern.
type Option func(concern *ReadConcern)

// Level creates an option that sets the level of a ReadConcern.
func Level(level string) Option ***REMOVED***
	return func(concern *ReadConcern) ***REMOVED***
		concern.level = level
	***REMOVED***
***REMOVED***

// Local specifies that the query should return the instance’s most recent data.
func Local() *ReadConcern ***REMOVED***
	return New(Level("local"))
***REMOVED***

// Majority specifies that the query should return the instance’s most recent data acknowledged as
// having been written to a majority of members in the replica set.
func Majority() *ReadConcern ***REMOVED***
	return New(Level("majority"))
***REMOVED***

// Linearizable specifies that the query should return data that reflects all successful writes
// issued with a write concern of "majority" and acknowledged prior to the start of the read operation.
func Linearizable() *ReadConcern ***REMOVED***
	return New(Level("linearizable"))
***REMOVED***

// Available specifies that the query should return data from the instance with no guarantee
// that the data has been written to a majority of the replica set members (i.e. may be rolled back).
func Available() *ReadConcern ***REMOVED***
	return New(Level("available"))
***REMOVED***

// Snapshot is only available for operations within multi-document transactions.
func Snapshot() *ReadConcern ***REMOVED***
	return New(Level("snapshot"))
***REMOVED***

// New constructs a new read concern from the given string.
func New(options ...Option) *ReadConcern ***REMOVED***
	concern := &ReadConcern***REMOVED******REMOVED***

	for _, option := range options ***REMOVED***
		option(concern)
	***REMOVED***

	return concern
***REMOVED***

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (rc *ReadConcern) MarshalBSONValue() (bsontype.Type, []byte, error) ***REMOVED***
	var elems []byte

	if len(rc.level) > 0 ***REMOVED***
		elems = bsoncore.AppendStringElement(elems, "level", rc.level)
	***REMOVED***

	return bsontype.EmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
***REMOVED***

// GetLevel returns the read concern level.
func (rc *ReadConcern) GetLevel() string ***REMOVED***
	return rc.level
***REMOVED***
