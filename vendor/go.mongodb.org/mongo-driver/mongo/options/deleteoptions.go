// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DeleteOptions represents options that can be used to configure DeleteOne and DeleteMany operations.
type DeleteOptions struct ***REMOVED***
	// Specifies a collation to use for string comparisons during the operation. This option is only valid for MongoDB
	// versions >= 3.4. For previous server versions, the driver will return an error if this option is used. The
	// default value is nil, which means the default collation of the collection will be used.
	Collation *Collation

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***

	// The index to use for the operation. This should either be the index name as a string or the index specification
	// as a document. This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.4 will return an error
	// if this option is specified. For server versions < 3.4, the driver will return a client-side error if this option
	// is specified. The driver will return an error if this option is specified during an unacknowledged write
	// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil,
	// which means that no hint will be sent.
	Hint interface***REMOVED******REMOVED***

	// Specifies parameters for the delete expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***
***REMOVED***

// Delete creates a new DeleteOptions instance.
func Delete() *DeleteOptions ***REMOVED***
	return &DeleteOptions***REMOVED******REMOVED***
***REMOVED***

// SetCollation sets the value for the Collation field.
func (do *DeleteOptions) SetCollation(c *Collation) *DeleteOptions ***REMOVED***
	do.Collation = c
	return do
***REMOVED***

// SetComment sets the value for the Comment field.
func (do *DeleteOptions) SetComment(comment interface***REMOVED******REMOVED***) *DeleteOptions ***REMOVED***
	do.Comment = comment
	return do
***REMOVED***

// SetHint sets the value for the Hint field.
func (do *DeleteOptions) SetHint(hint interface***REMOVED******REMOVED***) *DeleteOptions ***REMOVED***
	do.Hint = hint
	return do
***REMOVED***

// SetLet sets the value for the Let field.
func (do *DeleteOptions) SetLet(let interface***REMOVED******REMOVED***) *DeleteOptions ***REMOVED***
	do.Let = let
	return do
***REMOVED***

// MergeDeleteOptions combines the given DeleteOptions instances into a single DeleteOptions in a last-one-wins fashion.
func MergeDeleteOptions(opts ...*DeleteOptions) *DeleteOptions ***REMOVED***
	dOpts := Delete()
	for _, do := range opts ***REMOVED***
		if do == nil ***REMOVED***
			continue
		***REMOVED***
		if do.Collation != nil ***REMOVED***
			dOpts.Collation = do.Collation
		***REMOVED***
		if do.Comment != nil ***REMOVED***
			dOpts.Comment = do.Comment
		***REMOVED***
		if do.Hint != nil ***REMOVED***
			dOpts.Hint = do.Hint
		***REMOVED***
		if do.Let != nil ***REMOVED***
			dOpts.Let = do.Let
		***REMOVED***
	***REMOVED***

	return dOpts
***REMOVED***
