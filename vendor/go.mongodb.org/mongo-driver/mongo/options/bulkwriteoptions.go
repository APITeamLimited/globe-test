// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DefaultOrdered is the default value for the Ordered option in BulkWriteOptions.
var DefaultOrdered = true

// BulkWriteOptions represents options that can be used to configure a BulkWrite operation.
type BulkWriteOptions struct ***REMOVED***
	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***

	// If true, no writes will be executed after one fails. The default value is true.
	Ordered *bool

	// Specifies parameters for all update and delete commands in the BulkWrite. This option is only valid for MongoDB
	// versions >= 5.0. Older servers will report an error for using this option. This must be a document mapping
	// parameter names to values. Values must be constant or closed expressions that do not reference document fields.
	// Parameters can then be accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***
***REMOVED***

// BulkWrite creates a new *BulkWriteOptions instance.
func BulkWrite() *BulkWriteOptions ***REMOVED***
	return &BulkWriteOptions***REMOVED***
		Ordered: &DefaultOrdered,
	***REMOVED***
***REMOVED***

// SetComment sets the value for the Comment field.
func (b *BulkWriteOptions) SetComment(comment interface***REMOVED******REMOVED***) *BulkWriteOptions ***REMOVED***
	b.Comment = comment
	return b
***REMOVED***

// SetOrdered sets the value for the Ordered field.
func (b *BulkWriteOptions) SetOrdered(ordered bool) *BulkWriteOptions ***REMOVED***
	b.Ordered = &ordered
	return b
***REMOVED***

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (b *BulkWriteOptions) SetBypassDocumentValidation(bypass bool) *BulkWriteOptions ***REMOVED***
	b.BypassDocumentValidation = &bypass
	return b
***REMOVED***

// SetLet sets the value for the Let field. Let specifies parameters for all update and delete commands in the BulkWrite.
// This option is only valid for MongoDB versions >= 5.0. Older servers will report an error for using this option.
// This must be a document mapping parameter names to values. Values must be constant or closed expressions that do not
// reference document fields. Parameters can then be accessed as variables in an aggregate expression context (e.g. "$$var").
func (b *BulkWriteOptions) SetLet(let interface***REMOVED******REMOVED***) *BulkWriteOptions ***REMOVED***
	b.Let = &let
	return b
***REMOVED***

// MergeBulkWriteOptions combines the given BulkWriteOptions instances into a single BulkWriteOptions in a last-one-wins
// fashion.
func MergeBulkWriteOptions(opts ...*BulkWriteOptions) *BulkWriteOptions ***REMOVED***
	b := BulkWrite()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.Comment != nil ***REMOVED***
			b.Comment = opt.Comment
		***REMOVED***
		if opt.Ordered != nil ***REMOVED***
			b.Ordered = opt.Ordered
		***REMOVED***
		if opt.BypassDocumentValidation != nil ***REMOVED***
			b.BypassDocumentValidation = opt.BypassDocumentValidation
		***REMOVED***
		if opt.Let != nil ***REMOVED***
			b.Let = opt.Let
		***REMOVED***
	***REMOVED***

	return b
***REMOVED***
