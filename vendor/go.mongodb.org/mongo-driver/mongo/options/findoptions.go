// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"
)

// FindOptions represents options that can be used to configure a Find operation.
type FindOptions struct ***REMOVED***
	// AllowDiskUse specifies whether the server can write temporary data to disk while executing the Find operation.
	// This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.2 will report an error if this option
	// is specified. For server versions < 3.2, the driver will return a client-side error if this option is specified.
	// The default value is false.
	AllowDiskUse *bool

	// AllowPartial results specifies whether the Find operation on a sharded cluster can return partial results if some
	// shards are down rather than returning an error. The default value is false.
	AllowPartialResults *bool

	// BatchSize is the maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// Collation specifies a collation to use for string comparisons during the operation. This option is only valid for
	// MongoDB versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string that will be included in server logs, profiling logs, and currentOp queries to help trace the operation.
	// The default is nil, which means that no comment will be included in the logs.
	Comment *string

	// CursorType specifies the type of cursor that should be created for the operation. The default is NonTailable, which
	// means that the cursor will be closed by the server when the last batch of documents is retrieved.
	CursorType *CursorType

	// Hint is the index to use for the Find operation. This should either be the index name as a string or the index
	// specification as a document. The driver will return an error if the hint parameter is a multi-key map. The default
	// value is nil, which means that no hint will be sent.
	Hint interface***REMOVED******REMOVED***

	// Limit is the maximum number of documents to return. The default value is 0, which means that all documents matching the
	// filter will be returned. A negative limit specifies that the resulting documents should be returned in a single
	// batch. The default value is 0.
	Limit *int64

	// Max is a document specifying the exclusive upper bound for a specific index. The default value is nil, which means that
	// there is no maximum value.
	Max interface***REMOVED******REMOVED***

	// MaxAwaitTime is the maximum amount of time that the server should wait for new documents to satisfy a tailable cursor
	// query. This option is only valid for tailable await cursors (see the CursorType option for more information) and
	// MongoDB versions >= 3.2. For other cursor types or previous server versions, this option is ignored.
	MaxAwaitTime *time.Duration

	// MaxTime is the maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used in its
	// place to control the amount of time that a single operation can run before returning an error. MaxTime is ignored if
	// Timeout is set on the client.
	MaxTime *time.Duration

	// Min is a document specifying the inclusive lower bound for a specific index. The default value is 0, which means that
	// there is no minimum value.
	Min interface***REMOVED******REMOVED***

	// NoCursorTimeout specifies whether the cursor created by the operation will not timeout after a period of inactivity.
	// The default value is false.
	NoCursorTimeout *bool

	// OplogReplay is for internal replication use only and should not be set.
	//
	// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is
	// set.
	OplogReplay *bool

	// Project is a document describing which fields will be included in the documents returned by the Find operation. The
	// default value is nil, which means all fields will be included.
	Projection interface***REMOVED******REMOVED***

	// ReturnKey specifies whether the documents returned by the Find operation will only contain fields corresponding to the
	// index used. The default value is false.
	ReturnKey *bool

	// ShowRecordID specifies whether a $recordId field with a record identifier will be included in the documents returned by
	// the Find operation. The default value is false.
	ShowRecordID *bool

	// Skip is the number of documents to skip before adding documents to the result. The default value is 0.
	Skip *int64

	// Snapshot specifies whether the cursor will not return a document more than once because of an intervening write operation.
	// The default value is false.
	//
	// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
	Snapshot *bool

	// Sort is a document specifying the order in which documents should be returned.  The driver will return an error if the
	// sort parameter is a multi-key map.
	Sort interface***REMOVED******REMOVED***

	// Let specifies parameters for the find expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***
***REMOVED***

// Find creates a new FindOptions instance.
func Find() *FindOptions ***REMOVED***
	return &FindOptions***REMOVED******REMOVED***
***REMOVED***

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (f *FindOptions) SetAllowDiskUse(b bool) *FindOptions ***REMOVED***
	f.AllowDiskUse = &b
	return f
***REMOVED***

// SetAllowPartialResults sets the value for the AllowPartialResults field.
func (f *FindOptions) SetAllowPartialResults(b bool) *FindOptions ***REMOVED***
	f.AllowPartialResults = &b
	return f
***REMOVED***

// SetBatchSize sets the value for the BatchSize field.
func (f *FindOptions) SetBatchSize(i int32) *FindOptions ***REMOVED***
	f.BatchSize = &i
	return f
***REMOVED***

// SetCollation sets the value for the Collation field.
func (f *FindOptions) SetCollation(collation *Collation) *FindOptions ***REMOVED***
	f.Collation = collation
	return f
***REMOVED***

// SetComment sets the value for the Comment field.
func (f *FindOptions) SetComment(comment string) *FindOptions ***REMOVED***
	f.Comment = &comment
	return f
***REMOVED***

// SetCursorType sets the value for the CursorType field.
func (f *FindOptions) SetCursorType(ct CursorType) *FindOptions ***REMOVED***
	f.CursorType = &ct
	return f
***REMOVED***

// SetHint sets the value for the Hint field.
func (f *FindOptions) SetHint(hint interface***REMOVED******REMOVED***) *FindOptions ***REMOVED***
	f.Hint = hint
	return f
***REMOVED***

// SetLet sets the value for the Let field.
func (f *FindOptions) SetLet(let interface***REMOVED******REMOVED***) *FindOptions ***REMOVED***
	f.Let = let
	return f
***REMOVED***

// SetLimit sets the value for the Limit field.
func (f *FindOptions) SetLimit(i int64) *FindOptions ***REMOVED***
	f.Limit = &i
	return f
***REMOVED***

// SetMax sets the value for the Max field.
func (f *FindOptions) SetMax(max interface***REMOVED******REMOVED***) *FindOptions ***REMOVED***
	f.Max = max
	return f
***REMOVED***

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (f *FindOptions) SetMaxAwaitTime(d time.Duration) *FindOptions ***REMOVED***
	f.MaxAwaitTime = &d
	return f
***REMOVED***

// SetMaxTime specifies the max time to allow the query to run.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used used in its place to control the amount of time that a single operation
// can run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOptions) SetMaxTime(d time.Duration) *FindOptions ***REMOVED***
	f.MaxTime = &d
	return f
***REMOVED***

// SetMin sets the value for the Min field.
func (f *FindOptions) SetMin(min interface***REMOVED******REMOVED***) *FindOptions ***REMOVED***
	f.Min = min
	return f
***REMOVED***

// SetNoCursorTimeout sets the value for the NoCursorTimeout field.
func (f *FindOptions) SetNoCursorTimeout(b bool) *FindOptions ***REMOVED***
	f.NoCursorTimeout = &b
	return f
***REMOVED***

// SetOplogReplay sets the value for the OplogReplay field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is set.
func (f *FindOptions) SetOplogReplay(b bool) *FindOptions ***REMOVED***
	f.OplogReplay = &b
	return f
***REMOVED***

// SetProjection sets the value for the Projection field.
func (f *FindOptions) SetProjection(projection interface***REMOVED******REMOVED***) *FindOptions ***REMOVED***
	f.Projection = projection
	return f
***REMOVED***

// SetReturnKey sets the value for the ReturnKey field.
func (f *FindOptions) SetReturnKey(b bool) *FindOptions ***REMOVED***
	f.ReturnKey = &b
	return f
***REMOVED***

// SetShowRecordID sets the value for the ShowRecordID field.
func (f *FindOptions) SetShowRecordID(b bool) *FindOptions ***REMOVED***
	f.ShowRecordID = &b
	return f
***REMOVED***

// SetSkip sets the value for the Skip field.
func (f *FindOptions) SetSkip(i int64) *FindOptions ***REMOVED***
	f.Skip = &i
	return f
***REMOVED***

// SetSnapshot sets the value for the Snapshot field.
//
// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
func (f *FindOptions) SetSnapshot(b bool) *FindOptions ***REMOVED***
	f.Snapshot = &b
	return f
***REMOVED***

// SetSort sets the value for the Sort field.
func (f *FindOptions) SetSort(sort interface***REMOVED******REMOVED***) *FindOptions ***REMOVED***
	f.Sort = sort
	return f
***REMOVED***

// MergeFindOptions combines the given FindOptions instances into a single FindOptions in a last-one-wins fashion.
func MergeFindOptions(opts ...*FindOptions) *FindOptions ***REMOVED***
	fo := Find()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.AllowDiskUse != nil ***REMOVED***
			fo.AllowDiskUse = opt.AllowDiskUse
		***REMOVED***
		if opt.AllowPartialResults != nil ***REMOVED***
			fo.AllowPartialResults = opt.AllowPartialResults
		***REMOVED***
		if opt.BatchSize != nil ***REMOVED***
			fo.BatchSize = opt.BatchSize
		***REMOVED***
		if opt.Collation != nil ***REMOVED***
			fo.Collation = opt.Collation
		***REMOVED***
		if opt.Comment != nil ***REMOVED***
			fo.Comment = opt.Comment
		***REMOVED***
		if opt.CursorType != nil ***REMOVED***
			fo.CursorType = opt.CursorType
		***REMOVED***
		if opt.Hint != nil ***REMOVED***
			fo.Hint = opt.Hint
		***REMOVED***
		if opt.Let != nil ***REMOVED***
			fo.Let = opt.Let
		***REMOVED***
		if opt.Limit != nil ***REMOVED***
			fo.Limit = opt.Limit
		***REMOVED***
		if opt.Max != nil ***REMOVED***
			fo.Max = opt.Max
		***REMOVED***
		if opt.MaxAwaitTime != nil ***REMOVED***
			fo.MaxAwaitTime = opt.MaxAwaitTime
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			fo.MaxTime = opt.MaxTime
		***REMOVED***
		if opt.Min != nil ***REMOVED***
			fo.Min = opt.Min
		***REMOVED***
		if opt.NoCursorTimeout != nil ***REMOVED***
			fo.NoCursorTimeout = opt.NoCursorTimeout
		***REMOVED***
		if opt.OplogReplay != nil ***REMOVED***
			fo.OplogReplay = opt.OplogReplay
		***REMOVED***
		if opt.Projection != nil ***REMOVED***
			fo.Projection = opt.Projection
		***REMOVED***
		if opt.ReturnKey != nil ***REMOVED***
			fo.ReturnKey = opt.ReturnKey
		***REMOVED***
		if opt.ShowRecordID != nil ***REMOVED***
			fo.ShowRecordID = opt.ShowRecordID
		***REMOVED***
		if opt.Skip != nil ***REMOVED***
			fo.Skip = opt.Skip
		***REMOVED***
		if opt.Snapshot != nil ***REMOVED***
			fo.Snapshot = opt.Snapshot
		***REMOVED***
		if opt.Sort != nil ***REMOVED***
			fo.Sort = opt.Sort
		***REMOVED***
	***REMOVED***

	return fo
***REMOVED***

// FindOneOptions represents options that can be used to configure a FindOne operation.
type FindOneOptions struct ***REMOVED***
	// If true, an operation on a sharded cluster can return partial results if some shards are down rather than
	// returning an error. The default value is false.
	AllowPartialResults *bool

	// The maximum number of documents to be included in each batch returned by the server.
	//
	// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
	BatchSize *int32

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string that will be included in server logs, profiling logs, and currentOp queries to help trace the operation.
	// The default is nil, which means that no comment will be included in the logs.
	Comment *string

	// Specifies the type of cursor that should be created for the operation. The default is NonTailable, which means
	// that the cursor will be closed by the server when the last batch of documents is retrieved.
	//
	// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
	CursorType *CursorType

	// The index to use for the aggregation. This should either be the index name as a string or the index specification
	// as a document. The driver will return an error if the hint parameter is a multi-key map. The default value is nil,
	// which means that no hint will be sent.
	Hint interface***REMOVED******REMOVED***

	// A document specifying the exclusive upper bound for a specific index. The default value is nil, which means that
	// there is no maximum value.
	Max interface***REMOVED******REMOVED***

	// The maximum amount of time that the server should wait for new documents to satisfy a tailable cursor query.
	// This option is only valid for tailable await cursors (see the CursorType option for more information) and
	// MongoDB versions >= 3.2. For other cursor types or previous server versions, this option is ignored.
	//
	// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
	MaxAwaitTime *time.Duration

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration

	// A document specifying the inclusive lower bound for a specific index. The default value is 0, which means that
	// there is no minimum value.
	Min interface***REMOVED******REMOVED***

	// If true, the cursor created by the operation will not timeout after a period of inactivity. The default value
	// is false.
	//
	// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
	NoCursorTimeout *bool

	// This option is for internal replication use only and should not be set.
	//
	// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is
	// set.
	OplogReplay *bool

	// A document describing which fields will be included in the document returned by the operation. The default value
	// is nil, which means all fields will be included.
	Projection interface***REMOVED******REMOVED***

	// If true, the document returned by the operation will only contain fields corresponding to the index used. The
	// default value is false.
	ReturnKey *bool

	// If true, a $recordId field with a record identifier will be included in the document returned by the operation.
	// The default value is false.
	ShowRecordID *bool

	// The number of documents to skip before selecting the document to be returned. The default value is 0.
	Skip *int64

	// If true, the cursor will not return a document more than once because of an intervening write operation. The
	// default value is false.
	//
	// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
	Snapshot *bool

	// A document specifying the sort order to apply to the query. The first document in the sorted order will be
	// returned. The driver will return an error if the sort parameter is a multi-key map.
	Sort interface***REMOVED******REMOVED***
***REMOVED***

// FindOne creates a new FindOneOptions instance.
func FindOne() *FindOneOptions ***REMOVED***
	return &FindOneOptions***REMOVED******REMOVED***
***REMOVED***

// SetAllowPartialResults sets the value for the AllowPartialResults field.
func (f *FindOneOptions) SetAllowPartialResults(b bool) *FindOneOptions ***REMOVED***
	f.AllowPartialResults = &b
	return f
***REMOVED***

// SetBatchSize sets the value for the BatchSize field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetBatchSize(i int32) *FindOneOptions ***REMOVED***
	f.BatchSize = &i
	return f
***REMOVED***

// SetCollation sets the value for the Collation field.
func (f *FindOneOptions) SetCollation(collation *Collation) *FindOneOptions ***REMOVED***
	f.Collation = collation
	return f
***REMOVED***

// SetComment sets the value for the Comment field.
func (f *FindOneOptions) SetComment(comment string) *FindOneOptions ***REMOVED***
	f.Comment = &comment
	return f
***REMOVED***

// SetCursorType sets the value for the CursorType field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetCursorType(ct CursorType) *FindOneOptions ***REMOVED***
	f.CursorType = &ct
	return f
***REMOVED***

// SetHint sets the value for the Hint field.
func (f *FindOneOptions) SetHint(hint interface***REMOVED******REMOVED***) *FindOneOptions ***REMOVED***
	f.Hint = hint
	return f
***REMOVED***

// SetMax sets the value for the Max field.
func (f *FindOneOptions) SetMax(max interface***REMOVED******REMOVED***) *FindOneOptions ***REMOVED***
	f.Max = max
	return f
***REMOVED***

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetMaxAwaitTime(d time.Duration) *FindOneOptions ***REMOVED***
	f.MaxAwaitTime = &d
	return f
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneOptions) SetMaxTime(d time.Duration) *FindOneOptions ***REMOVED***
	f.MaxTime = &d
	return f
***REMOVED***

// SetMin sets the value for the Min field.
func (f *FindOneOptions) SetMin(min interface***REMOVED******REMOVED***) *FindOneOptions ***REMOVED***
	f.Min = min
	return f
***REMOVED***

// SetNoCursorTimeout sets the value for the NoCursorTimeout field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
func (f *FindOneOptions) SetNoCursorTimeout(b bool) *FindOneOptions ***REMOVED***
	f.NoCursorTimeout = &b
	return f
***REMOVED***

// SetOplogReplay sets the value for the OplogReplay field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is
// set.
func (f *FindOneOptions) SetOplogReplay(b bool) *FindOneOptions ***REMOVED***
	f.OplogReplay = &b
	return f
***REMOVED***

// SetProjection sets the value for the Projection field.
func (f *FindOneOptions) SetProjection(projection interface***REMOVED******REMOVED***) *FindOneOptions ***REMOVED***
	f.Projection = projection
	return f
***REMOVED***

// SetReturnKey sets the value for the ReturnKey field.
func (f *FindOneOptions) SetReturnKey(b bool) *FindOneOptions ***REMOVED***
	f.ReturnKey = &b
	return f
***REMOVED***

// SetShowRecordID sets the value for the ShowRecordID field.
func (f *FindOneOptions) SetShowRecordID(b bool) *FindOneOptions ***REMOVED***
	f.ShowRecordID = &b
	return f
***REMOVED***

// SetSkip sets the value for the Skip field.
func (f *FindOneOptions) SetSkip(i int64) *FindOneOptions ***REMOVED***
	f.Skip = &i
	return f
***REMOVED***

// SetSnapshot sets the value for the Snapshot field.
//
// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
func (f *FindOneOptions) SetSnapshot(b bool) *FindOneOptions ***REMOVED***
	f.Snapshot = &b
	return f
***REMOVED***

// SetSort sets the value for the Sort field.
func (f *FindOneOptions) SetSort(sort interface***REMOVED******REMOVED***) *FindOneOptions ***REMOVED***
	f.Sort = sort
	return f
***REMOVED***

// MergeFindOneOptions combines the given FindOneOptions instances into a single FindOneOptions in a last-one-wins
// fashion.
func MergeFindOneOptions(opts ...*FindOneOptions) *FindOneOptions ***REMOVED***
	fo := FindOne()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.AllowPartialResults != nil ***REMOVED***
			fo.AllowPartialResults = opt.AllowPartialResults
		***REMOVED***
		if opt.BatchSize != nil ***REMOVED***
			fo.BatchSize = opt.BatchSize
		***REMOVED***
		if opt.Collation != nil ***REMOVED***
			fo.Collation = opt.Collation
		***REMOVED***
		if opt.Comment != nil ***REMOVED***
			fo.Comment = opt.Comment
		***REMOVED***
		if opt.CursorType != nil ***REMOVED***
			fo.CursorType = opt.CursorType
		***REMOVED***
		if opt.Hint != nil ***REMOVED***
			fo.Hint = opt.Hint
		***REMOVED***
		if opt.Max != nil ***REMOVED***
			fo.Max = opt.Max
		***REMOVED***
		if opt.MaxAwaitTime != nil ***REMOVED***
			fo.MaxAwaitTime = opt.MaxAwaitTime
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			fo.MaxTime = opt.MaxTime
		***REMOVED***
		if opt.Min != nil ***REMOVED***
			fo.Min = opt.Min
		***REMOVED***
		if opt.NoCursorTimeout != nil ***REMOVED***
			fo.NoCursorTimeout = opt.NoCursorTimeout
		***REMOVED***
		if opt.OplogReplay != nil ***REMOVED***
			fo.OplogReplay = opt.OplogReplay
		***REMOVED***
		if opt.Projection != nil ***REMOVED***
			fo.Projection = opt.Projection
		***REMOVED***
		if opt.ReturnKey != nil ***REMOVED***
			fo.ReturnKey = opt.ReturnKey
		***REMOVED***
		if opt.ShowRecordID != nil ***REMOVED***
			fo.ShowRecordID = opt.ShowRecordID
		***REMOVED***
		if opt.Skip != nil ***REMOVED***
			fo.Skip = opt.Skip
		***REMOVED***
		if opt.Snapshot != nil ***REMOVED***
			fo.Snapshot = opt.Snapshot
		***REMOVED***
		if opt.Sort != nil ***REMOVED***
			fo.Sort = opt.Sort
		***REMOVED***
	***REMOVED***

	return fo
***REMOVED***

// FindOneAndReplaceOptions represents options that can be used to configure a FindOneAndReplace instance.
type FindOneAndReplaceOptions struct ***REMOVED***
	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration

	// A document describing which fields will be included in the document returned by the operation. The default value
	// is nil, which means all fields will be included.
	Projection interface***REMOVED******REMOVED***

	// Specifies whether the original or replaced document should be returned by the operation. The default value is
	// Before, which means the original document will be returned from before the replacement is performed.
	ReturnDocument *ReturnDocument

	// A document specifying which document should be replaced if the filter used by the operation matches multiple
	// documents in the collection. If set, the first document in the sorted order will be replaced. The driver will
	// return an error if the sort parameter is a multi-key map. The default value is nil.
	Sort interface***REMOVED******REMOVED***

	// If true, a new document will be inserted if the filter does not match any documents in the collection. The
	// default value is false.
	Upsert *bool

	// The index to use for the operation. This should either be the index name as a string or the index specification
	// as a document. This option is only valid for MongoDB versions >= 4.4. MongoDB version 4.2 will report an error if
	// this option is specified. For server versions < 4.2, the driver will return an error if this option is specified.
	// The driver will return an error if this option is used with during an unacknowledged write operation. The driver
	// will return an error if the hint parameter is a multi-key map. The default value is nil, which means that no hint
	// will be sent.
	Hint interface***REMOVED******REMOVED***

	// Specifies parameters for the find one and replace expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***
***REMOVED***

// FindOneAndReplace creates a new FindOneAndReplaceOptions instance.
func FindOneAndReplace() *FindOneAndReplaceOptions ***REMOVED***
	return &FindOneAndReplaceOptions***REMOVED******REMOVED***
***REMOVED***

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (f *FindOneAndReplaceOptions) SetBypassDocumentValidation(b bool) *FindOneAndReplaceOptions ***REMOVED***
	f.BypassDocumentValidation = &b
	return f
***REMOVED***

// SetCollation sets the value for the Collation field.
func (f *FindOneAndReplaceOptions) SetCollation(collation *Collation) *FindOneAndReplaceOptions ***REMOVED***
	f.Collation = collation
	return f
***REMOVED***

// SetComment sets the value for the Comment field.
func (f *FindOneAndReplaceOptions) SetComment(comment interface***REMOVED******REMOVED***) *FindOneAndReplaceOptions ***REMOVED***
	f.Comment = comment
	return f
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndReplaceOptions) SetMaxTime(d time.Duration) *FindOneAndReplaceOptions ***REMOVED***
	f.MaxTime = &d
	return f
***REMOVED***

// SetProjection sets the value for the Projection field.
func (f *FindOneAndReplaceOptions) SetProjection(projection interface***REMOVED******REMOVED***) *FindOneAndReplaceOptions ***REMOVED***
	f.Projection = projection
	return f
***REMOVED***

// SetReturnDocument sets the value for the ReturnDocument field.
func (f *FindOneAndReplaceOptions) SetReturnDocument(rd ReturnDocument) *FindOneAndReplaceOptions ***REMOVED***
	f.ReturnDocument = &rd
	return f
***REMOVED***

// SetSort sets the value for the Sort field.
func (f *FindOneAndReplaceOptions) SetSort(sort interface***REMOVED******REMOVED***) *FindOneAndReplaceOptions ***REMOVED***
	f.Sort = sort
	return f
***REMOVED***

// SetUpsert sets the value for the Upsert field.
func (f *FindOneAndReplaceOptions) SetUpsert(b bool) *FindOneAndReplaceOptions ***REMOVED***
	f.Upsert = &b
	return f
***REMOVED***

// SetHint sets the value for the Hint field.
func (f *FindOneAndReplaceOptions) SetHint(hint interface***REMOVED******REMOVED***) *FindOneAndReplaceOptions ***REMOVED***
	f.Hint = hint
	return f
***REMOVED***

// SetLet sets the value for the Let field.
func (f *FindOneAndReplaceOptions) SetLet(let interface***REMOVED******REMOVED***) *FindOneAndReplaceOptions ***REMOVED***
	f.Let = let
	return f
***REMOVED***

// MergeFindOneAndReplaceOptions combines the given FindOneAndReplaceOptions instances into a single
// FindOneAndReplaceOptions in a last-one-wins fashion.
func MergeFindOneAndReplaceOptions(opts ...*FindOneAndReplaceOptions) *FindOneAndReplaceOptions ***REMOVED***
	fo := FindOneAndReplace()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.BypassDocumentValidation != nil ***REMOVED***
			fo.BypassDocumentValidation = opt.BypassDocumentValidation
		***REMOVED***
		if opt.Collation != nil ***REMOVED***
			fo.Collation = opt.Collation
		***REMOVED***
		if opt.Comment != nil ***REMOVED***
			fo.Comment = opt.Comment
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			fo.MaxTime = opt.MaxTime
		***REMOVED***
		if opt.Projection != nil ***REMOVED***
			fo.Projection = opt.Projection
		***REMOVED***
		if opt.ReturnDocument != nil ***REMOVED***
			fo.ReturnDocument = opt.ReturnDocument
		***REMOVED***
		if opt.Sort != nil ***REMOVED***
			fo.Sort = opt.Sort
		***REMOVED***
		if opt.Upsert != nil ***REMOVED***
			fo.Upsert = opt.Upsert
		***REMOVED***
		if opt.Hint != nil ***REMOVED***
			fo.Hint = opt.Hint
		***REMOVED***
		if opt.Let != nil ***REMOVED***
			fo.Let = opt.Let
		***REMOVED***
	***REMOVED***

	return fo
***REMOVED***

// FindOneAndUpdateOptions represents options that can be used to configure a FindOneAndUpdate options.
type FindOneAndUpdateOptions struct ***REMOVED***
	// A set of filters specifying to which array elements an update should apply. This option is only valid for MongoDB
	// versions >= 3.6. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the update will apply to all array elements.
	ArrayFilters *ArrayFilters

	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime is
	// ignored if Timeout is set on the client.
	MaxTime *time.Duration

	// A document describing which fields will be included in the document returned by the operation. The default value
	// is nil, which means all fields will be included.
	Projection interface***REMOVED******REMOVED***

	// Specifies whether the original or replaced document should be returned by the operation. The default value is
	// Before, which means the original document will be returned before the replacement is performed.
	ReturnDocument *ReturnDocument

	// A document specifying which document should be updated if the filter used by the operation matches multiple
	// documents in the collection. If set, the first document in the sorted order will be updated. The driver will
	// return an error if the sort parameter is a multi-key map. The default value is nil.
	Sort interface***REMOVED******REMOVED***

	// If true, a new document will be inserted if the filter does not match any documents in the collection. The
	// default value is false.
	Upsert *bool

	// The index to use for the operation. This should either be the index name as a string or the index specification
	// as a document. This option is only valid for MongoDB versions >= 4.4. MongoDB version 4.2 will report an error if
	// this option is specified. For server versions < 4.2, the driver will return an error if this option is specified.
	// The driver will return an error if this option is used with during an unacknowledged write operation. The driver
	// will return an error if the hint parameter is a multi-key map. The default value is nil, which means that no hint
	// will be sent.
	Hint interface***REMOVED******REMOVED***

	// Specifies parameters for the find one and update expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***
***REMOVED***

// FindOneAndUpdate creates a new FindOneAndUpdateOptions instance.
func FindOneAndUpdate() *FindOneAndUpdateOptions ***REMOVED***
	return &FindOneAndUpdateOptions***REMOVED******REMOVED***
***REMOVED***

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (f *FindOneAndUpdateOptions) SetBypassDocumentValidation(b bool) *FindOneAndUpdateOptions ***REMOVED***
	f.BypassDocumentValidation = &b
	return f
***REMOVED***

// SetArrayFilters sets the value for the ArrayFilters field.
func (f *FindOneAndUpdateOptions) SetArrayFilters(filters ArrayFilters) *FindOneAndUpdateOptions ***REMOVED***
	f.ArrayFilters = &filters
	return f
***REMOVED***

// SetCollation sets the value for the Collation field.
func (f *FindOneAndUpdateOptions) SetCollation(collation *Collation) *FindOneAndUpdateOptions ***REMOVED***
	f.Collation = collation
	return f
***REMOVED***

// SetComment sets the value for the Comment field.
func (f *FindOneAndUpdateOptions) SetComment(comment interface***REMOVED******REMOVED***) *FindOneAndUpdateOptions ***REMOVED***
	f.Comment = comment
	return f
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndUpdateOptions) SetMaxTime(d time.Duration) *FindOneAndUpdateOptions ***REMOVED***
	f.MaxTime = &d
	return f
***REMOVED***

// SetProjection sets the value for the Projection field.
func (f *FindOneAndUpdateOptions) SetProjection(projection interface***REMOVED******REMOVED***) *FindOneAndUpdateOptions ***REMOVED***
	f.Projection = projection
	return f
***REMOVED***

// SetReturnDocument sets the value for the ReturnDocument field.
func (f *FindOneAndUpdateOptions) SetReturnDocument(rd ReturnDocument) *FindOneAndUpdateOptions ***REMOVED***
	f.ReturnDocument = &rd
	return f
***REMOVED***

// SetSort sets the value for the Sort field.
func (f *FindOneAndUpdateOptions) SetSort(sort interface***REMOVED******REMOVED***) *FindOneAndUpdateOptions ***REMOVED***
	f.Sort = sort
	return f
***REMOVED***

// SetUpsert sets the value for the Upsert field.
func (f *FindOneAndUpdateOptions) SetUpsert(b bool) *FindOneAndUpdateOptions ***REMOVED***
	f.Upsert = &b
	return f
***REMOVED***

// SetHint sets the value for the Hint field.
func (f *FindOneAndUpdateOptions) SetHint(hint interface***REMOVED******REMOVED***) *FindOneAndUpdateOptions ***REMOVED***
	f.Hint = hint
	return f
***REMOVED***

// SetLet sets the value for the Let field.
func (f *FindOneAndUpdateOptions) SetLet(let interface***REMOVED******REMOVED***) *FindOneAndUpdateOptions ***REMOVED***
	f.Let = let
	return f
***REMOVED***

// MergeFindOneAndUpdateOptions combines the given FindOneAndUpdateOptions instances into a single
// FindOneAndUpdateOptions in a last-one-wins fashion.
func MergeFindOneAndUpdateOptions(opts ...*FindOneAndUpdateOptions) *FindOneAndUpdateOptions ***REMOVED***
	fo := FindOneAndUpdate()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.ArrayFilters != nil ***REMOVED***
			fo.ArrayFilters = opt.ArrayFilters
		***REMOVED***
		if opt.BypassDocumentValidation != nil ***REMOVED***
			fo.BypassDocumentValidation = opt.BypassDocumentValidation
		***REMOVED***
		if opt.Collation != nil ***REMOVED***
			fo.Collation = opt.Collation
		***REMOVED***
		if opt.Comment != nil ***REMOVED***
			fo.Comment = opt.Comment
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			fo.MaxTime = opt.MaxTime
		***REMOVED***
		if opt.Projection != nil ***REMOVED***
			fo.Projection = opt.Projection
		***REMOVED***
		if opt.ReturnDocument != nil ***REMOVED***
			fo.ReturnDocument = opt.ReturnDocument
		***REMOVED***
		if opt.Sort != nil ***REMOVED***
			fo.Sort = opt.Sort
		***REMOVED***
		if opt.Upsert != nil ***REMOVED***
			fo.Upsert = opt.Upsert
		***REMOVED***
		if opt.Hint != nil ***REMOVED***
			fo.Hint = opt.Hint
		***REMOVED***
		if opt.Let != nil ***REMOVED***
			fo.Let = opt.Let
		***REMOVED***
	***REMOVED***

	return fo
***REMOVED***

// FindOneAndDeleteOptions represents options that can be used to configure a FindOneAndDelete operation.
type FindOneAndDeleteOptions struct ***REMOVED***
	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration

	// A document describing which fields will be included in the document returned by the operation. The default value
	// is nil, which means all fields will be included.
	Projection interface***REMOVED******REMOVED***

	// A document specifying which document should be replaced if the filter used by the operation matches multiple
	// documents in the collection. If set, the first document in the sorted order will be selected for replacement.
	// The driver will return an error if the sort parameter is a multi-key map. The default value is nil.
	Sort interface***REMOVED******REMOVED***

	// The index to use for the operation. This should either be the index name as a string or the index specification
	// as a document. This option is only valid for MongoDB versions >= 4.4. MongoDB version 4.2 will report an error if
	// this option is specified. For server versions < 4.2, the driver will return an error if this option is specified.
	// The driver will return an error if this option is used with during an unacknowledged write operation. The driver
	// will return an error if the hint parameter is a multi-key map. The default value is nil, which means that no hint
	// will be sent.
	Hint interface***REMOVED******REMOVED***

	// Specifies parameters for the find one and delete expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***
***REMOVED***

// FindOneAndDelete creates a new FindOneAndDeleteOptions instance.
func FindOneAndDelete() *FindOneAndDeleteOptions ***REMOVED***
	return &FindOneAndDeleteOptions***REMOVED******REMOVED***
***REMOVED***

// SetCollation sets the value for the Collation field.
func (f *FindOneAndDeleteOptions) SetCollation(collation *Collation) *FindOneAndDeleteOptions ***REMOVED***
	f.Collation = collation
	return f
***REMOVED***

// SetComment sets the value for the Comment field.
func (f *FindOneAndDeleteOptions) SetComment(comment interface***REMOVED******REMOVED***) *FindOneAndDeleteOptions ***REMOVED***
	f.Comment = comment
	return f
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndDeleteOptions) SetMaxTime(d time.Duration) *FindOneAndDeleteOptions ***REMOVED***
	f.MaxTime = &d
	return f
***REMOVED***

// SetProjection sets the value for the Projection field.
func (f *FindOneAndDeleteOptions) SetProjection(projection interface***REMOVED******REMOVED***) *FindOneAndDeleteOptions ***REMOVED***
	f.Projection = projection
	return f
***REMOVED***

// SetSort sets the value for the Sort field.
func (f *FindOneAndDeleteOptions) SetSort(sort interface***REMOVED******REMOVED***) *FindOneAndDeleteOptions ***REMOVED***
	f.Sort = sort
	return f
***REMOVED***

// SetHint sets the value for the Hint field.
func (f *FindOneAndDeleteOptions) SetHint(hint interface***REMOVED******REMOVED***) *FindOneAndDeleteOptions ***REMOVED***
	f.Hint = hint
	return f
***REMOVED***

// SetLet sets the value for the Let field.
func (f *FindOneAndDeleteOptions) SetLet(let interface***REMOVED******REMOVED***) *FindOneAndDeleteOptions ***REMOVED***
	f.Let = let
	return f
***REMOVED***

// MergeFindOneAndDeleteOptions combines the given FindOneAndDeleteOptions instances into a single
// FindOneAndDeleteOptions in a last-one-wins fashion.
func MergeFindOneAndDeleteOptions(opts ...*FindOneAndDeleteOptions) *FindOneAndDeleteOptions ***REMOVED***
	fo := FindOneAndDelete()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.Collation != nil ***REMOVED***
			fo.Collation = opt.Collation
		***REMOVED***
		if opt.Comment != nil ***REMOVED***
			fo.Comment = opt.Comment
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			fo.MaxTime = opt.MaxTime
		***REMOVED***
		if opt.Projection != nil ***REMOVED***
			fo.Projection = opt.Projection
		***REMOVED***
		if opt.Sort != nil ***REMOVED***
			fo.Sort = opt.Sort
		***REMOVED***
		if opt.Hint != nil ***REMOVED***
			fo.Hint = opt.Hint
		***REMOVED***
		if opt.Let != nil ***REMOVED***
			fo.Let = opt.Let
		***REMOVED***
	***REMOVED***

	return fo
***REMOVED***
