// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// ListCollectionsOptions represents options that can be used to configure a ListCollections operation.
type ListCollectionsOptions struct ***REMOVED***
	// If true, each collection document will only contain a field for the collection name. The default value is false.
	NameOnly *bool

	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// If true, and NameOnly is true, limits the documents returned to only contain collections the user is authorized to use. The default value
	// is false. This option is only valid for MongoDB server versions >= 4.0. Server versions < 4.0 ignore this option.
	AuthorizedCollections *bool
***REMOVED***

// ListCollections creates a new ListCollectionsOptions instance.
func ListCollections() *ListCollectionsOptions ***REMOVED***
	return &ListCollectionsOptions***REMOVED******REMOVED***
***REMOVED***

// SetNameOnly sets the value for the NameOnly field.
func (lc *ListCollectionsOptions) SetNameOnly(b bool) *ListCollectionsOptions ***REMOVED***
	lc.NameOnly = &b
	return lc
***REMOVED***

// SetBatchSize sets the value for the BatchSize field.
func (lc *ListCollectionsOptions) SetBatchSize(size int32) *ListCollectionsOptions ***REMOVED***
	lc.BatchSize = &size
	return lc
***REMOVED***

// SetAuthorizedCollections sets the value for the AuthorizedCollections field. This option is only valid for MongoDB server versions >= 4.0. Server
// versions < 4.0 ignore this option.
func (lc *ListCollectionsOptions) SetAuthorizedCollections(b bool) *ListCollectionsOptions ***REMOVED***
	lc.AuthorizedCollections = &b
	return lc
***REMOVED***

// MergeListCollectionsOptions combines the given ListCollectionsOptions instances into a single *ListCollectionsOptions
// in a last-one-wins fashion.
func MergeListCollectionsOptions(opts ...*ListCollectionsOptions) *ListCollectionsOptions ***REMOVED***
	lc := ListCollections()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.NameOnly != nil ***REMOVED***
			lc.NameOnly = opt.NameOnly
		***REMOVED***
		if opt.BatchSize != nil ***REMOVED***
			lc.BatchSize = opt.BatchSize
		***REMOVED***
		if opt.AuthorizedCollections != nil ***REMOVED***
			lc.AuthorizedCollections = opt.AuthorizedCollections
		***REMOVED***
	***REMOVED***

	return lc
***REMOVED***
