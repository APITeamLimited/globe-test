// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// CollectionOptions represents options that can be used to configure a Collection.
type CollectionOptions struct ***REMOVED***
	// ReadConcern is the read concern to use for operations executed on the Collection. The default value is nil, which means that
	// the read concern of the Database used to configure the Collection will be used.
	ReadConcern *readconcern.ReadConcern

	// WriteConcern is the write concern to use for operations executed on the Collection. The default value is nil, which means that
	// the write concern of the Database used to configure the Collection will be used.
	WriteConcern *writeconcern.WriteConcern

	// ReadPreference is the read preference to use for operations executed on the Collection. The default value is nil, which means that
	// the read preference of the Database used to configure the Collection will be used.
	ReadPreference *readpref.ReadPref

	// Registry is the BSON registry to marshal and unmarshal documents for operations executed on the Collection. The default value
	// is nil, which means that the registry of the Database used to configure the Collection will be used.
	Registry *bsoncodec.Registry
***REMOVED***

// Collection creates a new CollectionOptions instance.
func Collection() *CollectionOptions ***REMOVED***
	return &CollectionOptions***REMOVED******REMOVED***
***REMOVED***

// SetReadConcern sets the value for the ReadConcern field.
func (c *CollectionOptions) SetReadConcern(rc *readconcern.ReadConcern) *CollectionOptions ***REMOVED***
	c.ReadConcern = rc
	return c
***REMOVED***

// SetWriteConcern sets the value for the WriteConcern field.
func (c *CollectionOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *CollectionOptions ***REMOVED***
	c.WriteConcern = wc
	return c
***REMOVED***

// SetReadPreference sets the value for the ReadPreference field.
func (c *CollectionOptions) SetReadPreference(rp *readpref.ReadPref) *CollectionOptions ***REMOVED***
	c.ReadPreference = rp
	return c
***REMOVED***

// SetRegistry sets the value for the Registry field.
func (c *CollectionOptions) SetRegistry(r *bsoncodec.Registry) *CollectionOptions ***REMOVED***
	c.Registry = r
	return c
***REMOVED***

// MergeCollectionOptions combines the given CollectionOptions instances into a single *CollectionOptions in a
// last-one-wins fashion.
func MergeCollectionOptions(opts ...*CollectionOptions) *CollectionOptions ***REMOVED***
	c := Collection()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.ReadConcern != nil ***REMOVED***
			c.ReadConcern = opt.ReadConcern
		***REMOVED***
		if opt.WriteConcern != nil ***REMOVED***
			c.WriteConcern = opt.WriteConcern
		***REMOVED***
		if opt.ReadPreference != nil ***REMOVED***
			c.ReadPreference = opt.ReadPreference
		***REMOVED***
		if opt.Registry != nil ***REMOVED***
			c.Registry = opt.Registry
		***REMOVED***
	***REMOVED***

	return c
***REMOVED***
