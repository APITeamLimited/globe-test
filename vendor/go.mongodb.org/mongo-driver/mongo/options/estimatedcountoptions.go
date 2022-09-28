// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// EstimatedDocumentCountOptions represents options that can be used to configure an EstimatedDocumentCount operation.
type EstimatedDocumentCountOptions struct ***REMOVED***
	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration
***REMOVED***

// EstimatedDocumentCount creates a new EstimatedDocumentCountOptions instance.
func EstimatedDocumentCount() *EstimatedDocumentCountOptions ***REMOVED***
	return &EstimatedDocumentCountOptions***REMOVED******REMOVED***
***REMOVED***

// SetComment sets the value for the Comment field.
func (eco *EstimatedDocumentCountOptions) SetComment(comment interface***REMOVED******REMOVED***) *EstimatedDocumentCountOptions ***REMOVED***
	eco.Comment = comment
	return eco
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option
// may be used in its place to control the amount of time that a single operation can run before
// returning an error. MaxTime is ignored if Timeout is set on the client.
func (eco *EstimatedDocumentCountOptions) SetMaxTime(d time.Duration) *EstimatedDocumentCountOptions ***REMOVED***
	eco.MaxTime = &d
	return eco
***REMOVED***

// MergeEstimatedDocumentCountOptions combines the given EstimatedDocumentCountOptions instances into a single
// EstimatedDocumentCountOptions in a last-one-wins fashion.
func MergeEstimatedDocumentCountOptions(opts ...*EstimatedDocumentCountOptions) *EstimatedDocumentCountOptions ***REMOVED***
	e := EstimatedDocumentCount()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.Comment != nil ***REMOVED***
			e.Comment = opt.Comment
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			e.MaxTime = opt.MaxTime
		***REMOVED***
	***REMOVED***

	return e
***REMOVED***
