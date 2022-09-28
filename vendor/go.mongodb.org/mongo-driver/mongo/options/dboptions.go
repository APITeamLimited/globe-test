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

// DatabaseOptions represents options that can be used to configure a Database.
type DatabaseOptions struct ***REMOVED***
	// ReadConcern is the read concern to use for operations executed on the Database. The default value is nil, which means that
	// the read concern of the Client used to configure the Database will be used.
	ReadConcern *readconcern.ReadConcern

	// WriteConcern is the write concern to use for operations executed on the Database. The default value is nil, which means that the
	// write concern of the Client used to configure the Database will be used.
	WriteConcern *writeconcern.WriteConcern

	// ReadPreference is the read preference to use for operations executed on the Database. The default value is nil, which means that
	// the read preference of the Client used to configure the Database will be used.
	ReadPreference *readpref.ReadPref

	// Registry is the BSON registry to marshal and unmarshal documents for operations executed on the Database. The default value
	// is nil, which means that the registry of the Client used to configure the Database will be used.
	Registry *bsoncodec.Registry
***REMOVED***

// Database creates a new DatabaseOptions instance.
func Database() *DatabaseOptions ***REMOVED***
	return &DatabaseOptions***REMOVED******REMOVED***
***REMOVED***

// SetReadConcern sets the value for the ReadConcern field.
func (d *DatabaseOptions) SetReadConcern(rc *readconcern.ReadConcern) *DatabaseOptions ***REMOVED***
	d.ReadConcern = rc
	return d
***REMOVED***

// SetWriteConcern sets the value for the WriteConcern field.
func (d *DatabaseOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *DatabaseOptions ***REMOVED***
	d.WriteConcern = wc
	return d
***REMOVED***

// SetReadPreference sets the value for the ReadPreference field.
func (d *DatabaseOptions) SetReadPreference(rp *readpref.ReadPref) *DatabaseOptions ***REMOVED***
	d.ReadPreference = rp
	return d
***REMOVED***

// SetRegistry sets the value for the Registry field.
func (d *DatabaseOptions) SetRegistry(r *bsoncodec.Registry) *DatabaseOptions ***REMOVED***
	d.Registry = r
	return d
***REMOVED***

// MergeDatabaseOptions combines the given DatabaseOptions instances into a single DatabaseOptions in a last-one-wins
// fashion.
func MergeDatabaseOptions(opts ...*DatabaseOptions) *DatabaseOptions ***REMOVED***
	d := Database()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.ReadConcern != nil ***REMOVED***
			d.ReadConcern = opt.ReadConcern
		***REMOVED***
		if opt.WriteConcern != nil ***REMOVED***
			d.WriteConcern = opt.WriteConcern
		***REMOVED***
		if opt.ReadPreference != nil ***REMOVED***
			d.ReadPreference = opt.ReadPreference
		***REMOVED***
		if opt.Registry != nil ***REMOVED***
			d.Registry = opt.Registry
		***REMOVED***
	***REMOVED***

	return d
***REMOVED***
