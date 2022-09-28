// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// InsertOneOptions represents options that can be used to configure an InsertOne operation.
type InsertOneOptions struct ***REMOVED***
	// If true, writes executed as part of the operation will opt out of document-level validation on the server. This
	// option is valid for MongoDB versions >= 3.2 and is ignored for previous server versions. The default value is
	// false. See https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about document
	// validation.
	BypassDocumentValidation *bool

	// A string or document that will be included in server logs, profiling logs, and currentOp queries to help trace
	// the operation.  The default value is nil, which means that no comment will be included in the logs.
	Comment interface***REMOVED******REMOVED***
***REMOVED***

// InsertOne creates a new InsertOneOptions instance.
func InsertOne() *InsertOneOptions ***REMOVED***
	return &InsertOneOptions***REMOVED******REMOVED***
***REMOVED***

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (ioo *InsertOneOptions) SetBypassDocumentValidation(b bool) *InsertOneOptions ***REMOVED***
	ioo.BypassDocumentValidation = &b
	return ioo
***REMOVED***

// SetComment sets the value for the Comment field.
func (ioo *InsertOneOptions) SetComment(comment interface***REMOVED******REMOVED***) *InsertOneOptions ***REMOVED***
	ioo.Comment = comment
	return ioo
***REMOVED***

// MergeInsertOneOptions combines the given InsertOneOptions instances into a single InsertOneOptions in a last-one-wins
// fashion.
func MergeInsertOneOptions(opts ...*InsertOneOptions) *InsertOneOptions ***REMOVED***
	ioOpts := InsertOne()
	for _, ioo := range opts ***REMOVED***
		if ioo == nil ***REMOVED***
			continue
		***REMOVED***
		if ioo.BypassDocumentValidation != nil ***REMOVED***
			ioOpts.BypassDocumentValidation = ioo.BypassDocumentValidation
		***REMOVED***
		if ioo.Comment != nil ***REMOVED***
			ioOpts.Comment = ioo.Comment
		***REMOVED***
	***REMOVED***

	return ioOpts
***REMOVED***

// InsertManyOptions represents options that can be used to configure an InsertMany operation.
type InsertManyOptions struct ***REMOVED***
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
***REMOVED***

// InsertMany creates a new InsertManyOptions instance.
func InsertMany() *InsertManyOptions ***REMOVED***
	return &InsertManyOptions***REMOVED***
		Ordered: &DefaultOrdered,
	***REMOVED***
***REMOVED***

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (imo *InsertManyOptions) SetBypassDocumentValidation(b bool) *InsertManyOptions ***REMOVED***
	imo.BypassDocumentValidation = &b
	return imo
***REMOVED***

// SetComment sets the value for the Comment field.
func (imo *InsertManyOptions) SetComment(comment interface***REMOVED******REMOVED***) *InsertManyOptions ***REMOVED***
	imo.Comment = comment
	return imo
***REMOVED***

// SetOrdered sets the value for the Ordered field.
func (imo *InsertManyOptions) SetOrdered(b bool) *InsertManyOptions ***REMOVED***
	imo.Ordered = &b
	return imo
***REMOVED***

// MergeInsertManyOptions combines the given InsertManyOptions instances into a single InsertManyOptions in a last one
// wins fashion.
func MergeInsertManyOptions(opts ...*InsertManyOptions) *InsertManyOptions ***REMOVED***
	imOpts := InsertMany()
	for _, imo := range opts ***REMOVED***
		if imo == nil ***REMOVED***
			continue
		***REMOVED***
		if imo.BypassDocumentValidation != nil ***REMOVED***
			imOpts.BypassDocumentValidation = imo.BypassDocumentValidation
		***REMOVED***
		if imo.Comment != nil ***REMOVED***
			imOpts.Comment = imo.Comment
		***REMOVED***
		if imo.Ordered != nil ***REMOVED***
			imOpts.Ordered = imo.Ordered
		***REMOVED***
	***REMOVED***

	return imOpts
***REMOVED***
