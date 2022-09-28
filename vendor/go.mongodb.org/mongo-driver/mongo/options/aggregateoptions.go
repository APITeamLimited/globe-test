// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// AggregateOptions represents options that can be used to configure an Aggregate operation.
type AggregateOptions struct ***REMOVED***
	// If true, the operation can write to temporary files in the _tmp subdirectory of the database directory path on
	// the server. The default value is false.
	AllowDiskUse *bool

	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration

	// The maximum amount of time that the server should wait for new documents to satisfy a tailable cursor query.
	// This option is only valid for MongoDB versions >= 3.2 and is ignored for previous server versions.
	MaxAwaitTime *time.Duration

	// A string that will be included in server logs, profiling logs, and currentOp queries to help trace the operation.
	// The default is nil, which means that no comment will be included in the logs.
	Comment *string

	// The index to use for the aggregation. This should either be the index name as a string or the index specification
	// as a document. The hint does not apply to $lookup and $graphLookup aggregation stages. The driver will return an
	// error if the hint parameter is a multi-key map. The default value is nil, which means that no hint will be sent.
	Hint interface***REMOVED******REMOVED***

	// Specifies parameters for the aggregate expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***

	// Custom options to be added to aggregate expression. Key-value pairs of the BSON map should correlate with desired
	// option names and values. Values must be Marshalable. Custom options may conflict with non-custom options, and custom
	// options bypass client-side validation. Prefer using non-custom options where possible.
	Custom bson.M
***REMOVED***

// Aggregate creates a new AggregateOptions instance.
func Aggregate() *AggregateOptions ***REMOVED***
	return &AggregateOptions***REMOVED******REMOVED***
***REMOVED***

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (ao *AggregateOptions) SetAllowDiskUse(b bool) *AggregateOptions ***REMOVED***
	ao.AllowDiskUse = &b
	return ao
***REMOVED***

// SetBatchSize sets the value for the BatchSize field.
func (ao *AggregateOptions) SetBatchSize(i int32) *AggregateOptions ***REMOVED***
	ao.BatchSize = &i
	return ao
***REMOVED***

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (ao *AggregateOptions) SetBypassDocumentValidation(b bool) *AggregateOptions ***REMOVED***
	ao.BypassDocumentValidation = &b
	return ao
***REMOVED***

// SetCollation sets the value for the Collation field.
func (ao *AggregateOptions) SetCollation(c *Collation) *AggregateOptions ***REMOVED***
	ao.Collation = c
	return ao
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (ao *AggregateOptions) SetMaxTime(d time.Duration) *AggregateOptions ***REMOVED***
	ao.MaxTime = &d
	return ao
***REMOVED***

// SetMaxAwaitTime sets the value for the MaxAwaitTime field.
func (ao *AggregateOptions) SetMaxAwaitTime(d time.Duration) *AggregateOptions ***REMOVED***
	ao.MaxAwaitTime = &d
	return ao
***REMOVED***

// SetComment sets the value for the Comment field.
func (ao *AggregateOptions) SetComment(s string) *AggregateOptions ***REMOVED***
	ao.Comment = &s
	return ao
***REMOVED***

// SetHint sets the value for the Hint field.
func (ao *AggregateOptions) SetHint(h interface***REMOVED******REMOVED***) *AggregateOptions ***REMOVED***
	ao.Hint = h
	return ao
***REMOVED***

// SetLet sets the value for the Let field.
func (ao *AggregateOptions) SetLet(let interface***REMOVED******REMOVED***) *AggregateOptions ***REMOVED***
	ao.Let = let
	return ao
***REMOVED***

// SetCustom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (ao *AggregateOptions) SetCustom(c bson.M) *AggregateOptions ***REMOVED***
	ao.Custom = c
	return ao
***REMOVED***

// MergeAggregateOptions combines the given AggregateOptions instances into a single AggregateOptions in a last-one-wins
// fashion.
func MergeAggregateOptions(opts ...*AggregateOptions) *AggregateOptions ***REMOVED***
	aggOpts := Aggregate()
	for _, ao := range opts ***REMOVED***
		if ao == nil ***REMOVED***
			continue
		***REMOVED***
		if ao.AllowDiskUse != nil ***REMOVED***
			aggOpts.AllowDiskUse = ao.AllowDiskUse
		***REMOVED***
		if ao.BatchSize != nil ***REMOVED***
			aggOpts.BatchSize = ao.BatchSize
		***REMOVED***
		if ao.BypassDocumentValidation != nil ***REMOVED***
			aggOpts.BypassDocumentValidation = ao.BypassDocumentValidation
		***REMOVED***
		if ao.Collation != nil ***REMOVED***
			aggOpts.Collation = ao.Collation
		***REMOVED***
		if ao.MaxTime != nil ***REMOVED***
			aggOpts.MaxTime = ao.MaxTime
		***REMOVED***
		if ao.MaxAwaitTime != nil ***REMOVED***
			aggOpts.MaxAwaitTime = ao.MaxAwaitTime
		***REMOVED***
		if ao.Comment != nil ***REMOVED***
			aggOpts.Comment = ao.Comment
		***REMOVED***
		if ao.Hint != nil ***REMOVED***
			aggOpts.Hint = ao.Hint
		***REMOVED***
		if ao.Let != nil ***REMOVED***
			aggOpts.Let = ao.Let
		***REMOVED***
		if ao.Custom != nil ***REMOVED***
			aggOpts.Custom = ao.Custom
		***REMOVED***
	***REMOVED***

	return aggOpts
***REMOVED***
