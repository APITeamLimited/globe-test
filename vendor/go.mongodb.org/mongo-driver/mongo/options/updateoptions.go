// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// UpdateOptions represents options that can be used to configure UpdateOne and UpdateMany operations.
type UpdateOptions struct ***REMOVED***
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

	// The index to use for the operation. This should either be the index name as a string or the index specification
	// as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will return an error
	// if this option is specified. For server versions < 3.4, the driver will return a client-side error if this option
	// is specified. The driver will return an error if this option is specified during an unacknowledged write
	// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil,
	// which means that no hint will be sent.
	Hint interface***REMOVED******REMOVED***

	// If true, a new document will be inserted if the filter does not match any documents in the collection. The
	// default value is false.
	Upsert *bool

	// Specifies parameters for the update expression. This option is only valid for MongoDB versions >= 5.0. Older
	// servers will report an error for using this option. This must be a document mapping parameter names to values.
	// Values must be constant or closed expressions that do not reference document fields. Parameters can then be
	// accessed as variables in an aggregate expression context (e.g. "$$var").
	Let interface***REMOVED******REMOVED***
***REMOVED***

// Update creates a new UpdateOptions instance.
func Update() *UpdateOptions ***REMOVED***
	return &UpdateOptions***REMOVED******REMOVED***
***REMOVED***

// SetArrayFilters sets the value for the ArrayFilters field.
func (uo *UpdateOptions) SetArrayFilters(af ArrayFilters) *UpdateOptions ***REMOVED***
	uo.ArrayFilters = &af
	return uo
***REMOVED***

// SetBypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (uo *UpdateOptions) SetBypassDocumentValidation(b bool) *UpdateOptions ***REMOVED***
	uo.BypassDocumentValidation = &b
	return uo
***REMOVED***

// SetCollation sets the value for the Collation field.
func (uo *UpdateOptions) SetCollation(c *Collation) *UpdateOptions ***REMOVED***
	uo.Collation = c
	return uo
***REMOVED***

// SetComment sets the value for the Comment field.
func (uo *UpdateOptions) SetComment(comment interface***REMOVED******REMOVED***) *UpdateOptions ***REMOVED***
	uo.Comment = comment
	return uo
***REMOVED***

// SetHint sets the value for the Hint field.
func (uo *UpdateOptions) SetHint(h interface***REMOVED******REMOVED***) *UpdateOptions ***REMOVED***
	uo.Hint = h
	return uo
***REMOVED***

// SetUpsert sets the value for the Upsert field.
func (uo *UpdateOptions) SetUpsert(b bool) *UpdateOptions ***REMOVED***
	uo.Upsert = &b
	return uo
***REMOVED***

// SetLet sets the value for the Let field.
func (uo *UpdateOptions) SetLet(l interface***REMOVED******REMOVED***) *UpdateOptions ***REMOVED***
	uo.Let = l
	return uo
***REMOVED***

// MergeUpdateOptions combines the given UpdateOptions instances into a single UpdateOptions in a last-one-wins fashion.
func MergeUpdateOptions(opts ...*UpdateOptions) *UpdateOptions ***REMOVED***
	uOpts := Update()
	for _, uo := range opts ***REMOVED***
		if uo == nil ***REMOVED***
			continue
		***REMOVED***
		if uo.ArrayFilters != nil ***REMOVED***
			uOpts.ArrayFilters = uo.ArrayFilters
		***REMOVED***
		if uo.BypassDocumentValidation != nil ***REMOVED***
			uOpts.BypassDocumentValidation = uo.BypassDocumentValidation
		***REMOVED***
		if uo.Collation != nil ***REMOVED***
			uOpts.Collation = uo.Collation
		***REMOVED***
		if uo.Comment != nil ***REMOVED***
			uOpts.Comment = uo.Comment
		***REMOVED***
		if uo.Hint != nil ***REMOVED***
			uOpts.Hint = uo.Hint
		***REMOVED***
		if uo.Upsert != nil ***REMOVED***
			uOpts.Upsert = uo.Upsert
		***REMOVED***
		if uo.Let != nil ***REMOVED***
			uOpts.Let = uo.Let
		***REMOVED***
	***REMOVED***

	return uOpts
***REMOVED***
