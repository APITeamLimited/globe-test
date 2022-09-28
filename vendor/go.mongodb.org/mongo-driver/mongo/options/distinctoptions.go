// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "time"

// DistinctOptions represents options that can be used to configure a Distinct operation.
type DistinctOptions struct ***REMOVED***
	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation. The default value is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be
	// used in its place to control the amount of time that a single operation can run before returning an error.
	// MaxTime is ignored if Timeout is set on the client.
	MaxTime *time.Duration
***REMOVED***

// Distinct creates a new DistinctOptions instance.
func Distinct() *DistinctOptions ***REMOVED***
	return &DistinctOptions***REMOVED******REMOVED***
***REMOVED***

// SetCollation sets the value for the Collation field.
func (do *DistinctOptions) SetCollation(c *Collation) *DistinctOptions ***REMOVED***
	do.Collation = c
	return do
***REMOVED***

// SetComment sets the value for the Comment field.
func (do *DistinctOptions) SetComment(comment interface***REMOVED******REMOVED***) *DistinctOptions ***REMOVED***
	do.Comment = comment
	return do
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (do *DistinctOptions) SetMaxTime(d time.Duration) *DistinctOptions ***REMOVED***
	do.MaxTime = &d
	return do
***REMOVED***

// MergeDistinctOptions combines the given DistinctOptions instances into a single DistinctOptions in a last-one-wins
// fashion.
func MergeDistinctOptions(opts ...*DistinctOptions) *DistinctOptions ***REMOVED***
	distinctOpts := Distinct()
	for _, do := range opts ***REMOVED***
		if do == nil ***REMOVED***
			continue
		***REMOVED***
		if do.Collation != nil ***REMOVED***
			distinctOpts.Collation = do.Collation
		***REMOVED***
		if do.Comment != nil ***REMOVED***
			distinctOpts.Comment = do.Comment
		***REMOVED***
		if do.MaxTime != nil ***REMOVED***
			distinctOpts.MaxTime = do.MaxTime
		***REMOVED***
	***REMOVED***

	return distinctOpts
***REMOVED***
