// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// ClientOptions represents all possible options for creating a client session.
type ClientOptions struct ***REMOVED***
	CausalConsistency     *bool
	DefaultReadConcern    *readconcern.ReadConcern
	DefaultWriteConcern   *writeconcern.WriteConcern
	DefaultReadPreference *readpref.ReadPref
	DefaultMaxCommitTime  *time.Duration
	Snapshot              *bool
***REMOVED***

// TransactionOptions represents all possible options for starting a transaction in a session.
type TransactionOptions struct ***REMOVED***
	ReadConcern    *readconcern.ReadConcern
	WriteConcern   *writeconcern.WriteConcern
	ReadPreference *readpref.ReadPref
	MaxCommitTime  *time.Duration
***REMOVED***

func mergeClientOptions(opts ...*ClientOptions) *ClientOptions ***REMOVED***
	c := &ClientOptions***REMOVED******REMOVED***
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.CausalConsistency != nil ***REMOVED***
			c.CausalConsistency = opt.CausalConsistency
		***REMOVED***
		if opt.DefaultReadConcern != nil ***REMOVED***
			c.DefaultReadConcern = opt.DefaultReadConcern
		***REMOVED***
		if opt.DefaultReadPreference != nil ***REMOVED***
			c.DefaultReadPreference = opt.DefaultReadPreference
		***REMOVED***
		if opt.DefaultWriteConcern != nil ***REMOVED***
			c.DefaultWriteConcern = opt.DefaultWriteConcern
		***REMOVED***
		if opt.DefaultMaxCommitTime != nil ***REMOVED***
			c.DefaultMaxCommitTime = opt.DefaultMaxCommitTime
		***REMOVED***
		if opt.Snapshot != nil ***REMOVED***
			c.Snapshot = opt.Snapshot
		***REMOVED***
	***REMOVED***

	return c
***REMOVED***
