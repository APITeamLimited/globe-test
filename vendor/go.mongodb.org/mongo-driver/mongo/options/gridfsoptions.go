// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// DefaultName is the default name for a GridFS bucket.
var DefaultName = "fs"

// DefaultChunkSize is the default size of each file chunk in bytes (255 KiB).
var DefaultChunkSize int32 = 255 * 1024

// DefaultRevision is the default revision number for a download by name operation.
var DefaultRevision int32 = -1

// BucketOptions represents options that can be used to configure GridFS bucket.
type BucketOptions struct ***REMOVED***
	// The name of the bucket. The default value is "fs".
	Name *string

	// The number of bytes in each chunk in the bucket. The default value is 255 KiB.
	ChunkSizeBytes *int32

	// The write concern for the bucket. The default value is the write concern of the database from which the bucket
	// is created.
	WriteConcern *writeconcern.WriteConcern

	// The read concern for the bucket. The default value is the read concern of the database from which the bucket
	// is created.
	ReadConcern *readconcern.ReadConcern

	// The read preference for the bucket. The default value is the read preference of the database from which the
	// bucket is created.
	ReadPreference *readpref.ReadPref
***REMOVED***

// GridFSBucket creates a new BucketOptions instance.
func GridFSBucket() *BucketOptions ***REMOVED***
	return &BucketOptions***REMOVED***
		Name:           &DefaultName,
		ChunkSizeBytes: &DefaultChunkSize,
	***REMOVED***
***REMOVED***

// SetName sets the value for the Name field.
func (b *BucketOptions) SetName(name string) *BucketOptions ***REMOVED***
	b.Name = &name
	return b
***REMOVED***

// SetChunkSizeBytes sets the value for the ChunkSize field.
func (b *BucketOptions) SetChunkSizeBytes(i int32) *BucketOptions ***REMOVED***
	b.ChunkSizeBytes = &i
	return b
***REMOVED***

// SetWriteConcern sets the value for the WriteConcern field.
func (b *BucketOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *BucketOptions ***REMOVED***
	b.WriteConcern = wc
	return b
***REMOVED***

// SetReadConcern sets the value for the ReadConcern field.
func (b *BucketOptions) SetReadConcern(rc *readconcern.ReadConcern) *BucketOptions ***REMOVED***
	b.ReadConcern = rc
	return b
***REMOVED***

// SetReadPreference sets the value for the ReadPreference field.
func (b *BucketOptions) SetReadPreference(rp *readpref.ReadPref) *BucketOptions ***REMOVED***
	b.ReadPreference = rp
	return b
***REMOVED***

// MergeBucketOptions combines the given BucketOptions instances into a single BucketOptions in a last-one-wins fashion.
func MergeBucketOptions(opts ...*BucketOptions) *BucketOptions ***REMOVED***
	b := GridFSBucket()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.Name != nil ***REMOVED***
			b.Name = opt.Name
		***REMOVED***
		if opt.ChunkSizeBytes != nil ***REMOVED***
			b.ChunkSizeBytes = opt.ChunkSizeBytes
		***REMOVED***
		if opt.WriteConcern != nil ***REMOVED***
			b.WriteConcern = opt.WriteConcern
		***REMOVED***
		if opt.ReadConcern != nil ***REMOVED***
			b.ReadConcern = opt.ReadConcern
		***REMOVED***
		if opt.ReadPreference != nil ***REMOVED***
			b.ReadPreference = opt.ReadPreference
		***REMOVED***
	***REMOVED***

	return b
***REMOVED***

// UploadOptions represents options that can be used to configure a GridFS upload operation.
type UploadOptions struct ***REMOVED***
	// The number of bytes in each chunk in the bucket. The default value is DefaultChunkSize (255 KiB).
	ChunkSizeBytes *int32

	// Additional application data that will be stored in the "metadata" field of the document in the files collection.
	// The default value is nil, which means that the document in the files collection will not contain a "metadata"
	// field.
	Metadata interface***REMOVED******REMOVED***

	// The BSON registry to use for converting filters to BSON documents. The default value is bson.DefaultRegistry.
	Registry *bsoncodec.Registry
***REMOVED***

// GridFSUpload creates a new UploadOptions instance.
func GridFSUpload() *UploadOptions ***REMOVED***
	return &UploadOptions***REMOVED***Registry: bson.DefaultRegistry***REMOVED***
***REMOVED***

// SetChunkSizeBytes sets the value for the ChunkSize field.
func (u *UploadOptions) SetChunkSizeBytes(i int32) *UploadOptions ***REMOVED***
	u.ChunkSizeBytes = &i
	return u
***REMOVED***

// SetMetadata sets the value for the Metadata field.
func (u *UploadOptions) SetMetadata(doc interface***REMOVED******REMOVED***) *UploadOptions ***REMOVED***
	u.Metadata = doc
	return u
***REMOVED***

// MergeUploadOptions combines the given UploadOptions instances into a single UploadOptions in a last-one-wins fashion.
func MergeUploadOptions(opts ...*UploadOptions) *UploadOptions ***REMOVED***
	u := GridFSUpload()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.ChunkSizeBytes != nil ***REMOVED***
			u.ChunkSizeBytes = opt.ChunkSizeBytes
		***REMOVED***
		if opt.Metadata != nil ***REMOVED***
			u.Metadata = opt.Metadata
		***REMOVED***
		if opt.Registry != nil ***REMOVED***
			u.Registry = opt.Registry
		***REMOVED***
	***REMOVED***

	return u
***REMOVED***

// NameOptions represents options that can be used to configure a GridFS DownloadByName operation.
type NameOptions struct ***REMOVED***
	// Specifies the revision of the file to retrieve. Revision numbers are defined as follows:
	//
	// * 0 = the original stored file
	// * 1 = the first revision
	// * 2 = the second revision
	// * etc..
	// * -2 = the second most recent revision
	// * -1 = the most recent revision.
	//
	// The default value is -1
	Revision *int32
***REMOVED***

// GridFSName creates a new NameOptions instance.
func GridFSName() *NameOptions ***REMOVED***
	return &NameOptions***REMOVED******REMOVED***
***REMOVED***

// SetRevision sets the value for the Revision field.
func (n *NameOptions) SetRevision(r int32) *NameOptions ***REMOVED***
	n.Revision = &r
	return n
***REMOVED***

// MergeNameOptions combines the given NameOptions instances into a single *NameOptions in a last-one-wins fashion.
func MergeNameOptions(opts ...*NameOptions) *NameOptions ***REMOVED***
	n := GridFSName()
	n.Revision = &DefaultRevision

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.Revision != nil ***REMOVED***
			n.Revision = opt.Revision
		***REMOVED***
	***REMOVED***

	return n
***REMOVED***

// GridFSFindOptions represents options that can be used to configure a GridFS Find operation.
type GridFSFindOptions struct ***REMOVED***
	// If true, the server can write temporary data to disk while executing the find operation. The default value
	// is false. This option is only valid for MongoDB versions >= 4.4. For previous server versions, the server will
	// return an error if this option is used.
	AllowDiskUse *bool

	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// The maximum number of documents to return. The default value is 0, which means that all documents matching the
	// filter will be returned. A negative limit specifies that the resulting documents should be returned in a single
	// batch. The default value is 0.
	Limit *int32

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration

	// If true, the cursor created by the operation will not timeout after a period of inactivity. The default value
	// is false.
	NoCursorTimeout *bool

	// The number of documents to skip before adding documents to the result. The default value is 0.
	Skip *int32

	// A document specifying the order in which documents should be returned.  The driver will return an error if the
	// sort parameter is a multi-key map.
	Sort interface***REMOVED******REMOVED***
***REMOVED***

// GridFSFind creates a new GridFSFindOptions instance.
func GridFSFind() *GridFSFindOptions ***REMOVED***
	return &GridFSFindOptions***REMOVED******REMOVED***
***REMOVED***

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (f *GridFSFindOptions) SetAllowDiskUse(b bool) *GridFSFindOptions ***REMOVED***
	f.AllowDiskUse = &b
	return f
***REMOVED***

// SetBatchSize sets the value for the BatchSize field.
func (f *GridFSFindOptions) SetBatchSize(i int32) *GridFSFindOptions ***REMOVED***
	f.BatchSize = &i
	return f
***REMOVED***

// SetLimit sets the value for the Limit field.
func (f *GridFSFindOptions) SetLimit(i int32) *GridFSFindOptions ***REMOVED***
	f.Limit = &i
	return f
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *GridFSFindOptions) SetMaxTime(d time.Duration) *GridFSFindOptions ***REMOVED***
	f.MaxTime = &d
	return f
***REMOVED***

// SetNoCursorTimeout sets the value for the NoCursorTimeout field.
func (f *GridFSFindOptions) SetNoCursorTimeout(b bool) *GridFSFindOptions ***REMOVED***
	f.NoCursorTimeout = &b
	return f
***REMOVED***

// SetSkip sets the value for the Skip field.
func (f *GridFSFindOptions) SetSkip(i int32) *GridFSFindOptions ***REMOVED***
	f.Skip = &i
	return f
***REMOVED***

// SetSort sets the value for the Sort field.
func (f *GridFSFindOptions) SetSort(sort interface***REMOVED******REMOVED***) *GridFSFindOptions ***REMOVED***
	f.Sort = sort
	return f
***REMOVED***

// MergeGridFSFindOptions combines the given GridFSFindOptions instances into a single GridFSFindOptions in a
// last-one-wins fashion.
func MergeGridFSFindOptions(opts ...*GridFSFindOptions) *GridFSFindOptions ***REMOVED***
	fo := GridFSFind()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.AllowDiskUse != nil ***REMOVED***
			fo.AllowDiskUse = opt.AllowDiskUse
		***REMOVED***
		if opt.BatchSize != nil ***REMOVED***
			fo.BatchSize = opt.BatchSize
		***REMOVED***
		if opt.Limit != nil ***REMOVED***
			fo.Limit = opt.Limit
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			fo.MaxTime = opt.MaxTime
		***REMOVED***
		if opt.NoCursorTimeout != nil ***REMOVED***
			fo.NoCursorTimeout = opt.NoCursorTimeout
		***REMOVED***
		if opt.Skip != nil ***REMOVED***
			fo.Skip = opt.Skip
		***REMOVED***
		if opt.Sort != nil ***REMOVED***
			fo.Sort = opt.Sort
		***REMOVED***
	***REMOVED***

	return fo
***REMOVED***
