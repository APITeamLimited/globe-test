// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Collation allows users to specify language-specific rules for string comparison, such as
// rules for lettercase and accent marks.
type Collation struct ***REMOVED***
	Locale          string `bson:",omitempty"` // The locale
	CaseLevel       bool   `bson:",omitempty"` // The case level
	CaseFirst       string `bson:",omitempty"` // The case ordering
	Strength        int    `bson:",omitempty"` // The number of comparison levels to use
	NumericOrdering bool   `bson:",omitempty"` // Whether to order numbers based on numerical order and not collation order
	Alternate       string `bson:",omitempty"` // Whether spaces and punctuation are considered base characters
	MaxVariable     string `bson:",omitempty"` // Which characters are affected by alternate: "shifted"
	Normalization   bool   `bson:",omitempty"` // Causes text to be normalized into Unicode NFD
	Backwards       bool   `bson:",omitempty"` // Causes secondary differences to be considered in reverse order, as it is done in the French language
***REMOVED***

// ToDocument converts the Collation to a bson.Raw.
func (co *Collation) ToDocument() bson.Raw ***REMOVED***
	idx, doc := bsoncore.AppendDocumentStart(nil)
	if co.Locale != "" ***REMOVED***
		doc = bsoncore.AppendStringElement(doc, "locale", co.Locale)
	***REMOVED***
	if co.CaseLevel ***REMOVED***
		doc = bsoncore.AppendBooleanElement(doc, "caseLevel", true)
	***REMOVED***
	if co.CaseFirst != "" ***REMOVED***
		doc = bsoncore.AppendStringElement(doc, "caseFirst", co.CaseFirst)
	***REMOVED***
	if co.Strength != 0 ***REMOVED***
		doc = bsoncore.AppendInt32Element(doc, "strength", int32(co.Strength))
	***REMOVED***
	if co.NumericOrdering ***REMOVED***
		doc = bsoncore.AppendBooleanElement(doc, "numericOrdering", true)
	***REMOVED***
	if co.Alternate != "" ***REMOVED***
		doc = bsoncore.AppendStringElement(doc, "alternate", co.Alternate)
	***REMOVED***
	if co.MaxVariable != "" ***REMOVED***
		doc = bsoncore.AppendStringElement(doc, "maxVariable", co.MaxVariable)
	***REMOVED***
	if co.Normalization ***REMOVED***
		doc = bsoncore.AppendBooleanElement(doc, "normalization", true)
	***REMOVED***
	if co.Backwards ***REMOVED***
		doc = bsoncore.AppendBooleanElement(doc, "backwards", true)
	***REMOVED***
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc
***REMOVED***

// CursorType specifies whether a cursor should close when the last data is retrieved. See
// NonTailable, Tailable, and TailableAwait.
type CursorType int8

const (
	// NonTailable specifies that a cursor should close after retrieving the last data.
	NonTailable CursorType = iota
	// Tailable specifies that a cursor should not close when the last data is retrieved and can be resumed later.
	Tailable
	// TailableAwait specifies that a cursor should not close when the last data is retrieved and
	// that it should block for a certain amount of time for new data before returning no data.
	TailableAwait
)

// ReturnDocument specifies whether a findAndUpdate operation should return the document as it was
// before the update or as it is after the update.
type ReturnDocument int8

const (
	// Before specifies that findAndUpdate should return the document as it was before the update.
	Before ReturnDocument = iota
	// After specifies that findAndUpdate should return the document as it is after the update.
	After
)

// FullDocument specifies how a change stream should return the modified document.
type FullDocument string

const (
	// Default does not include a document copy.
	Default FullDocument = "default"
	// Off is the same as sending no value for fullDocumentBeforeChange.
	Off FullDocument = "off"
	// Required is the same as WhenAvailable but raises a server-side error if the post-image is not available.
	Required FullDocument = "required"
	// UpdateLookup includes a delta describing the changes to the document and a copy of the entire document that
	// was changed.
	UpdateLookup FullDocument = "updateLookup"
	// WhenAvailable includes a post-image of the the modified document for replace and update change events
	// if the post-image for this event is available.
	WhenAvailable FullDocument = "whenAvailable"
)

// ArrayFilters is used to hold filters for the array filters CRUD option. If a registry is nil, bson.DefaultRegistry
// will be used when converting the filter interfaces to BSON.
type ArrayFilters struct ***REMOVED***
	Registry *bsoncodec.Registry // The registry to use for converting filters. Defaults to bson.DefaultRegistry.
	Filters  []interface***REMOVED******REMOVED***       // The filters to apply
***REMOVED***

// ToArray builds a []bson.Raw from the provided ArrayFilters.
func (af *ArrayFilters) ToArray() ([]bson.Raw, error) ***REMOVED***
	registry := af.Registry
	if registry == nil ***REMOVED***
		registry = bson.DefaultRegistry
	***REMOVED***
	filters := make([]bson.Raw, 0, len(af.Filters))
	for _, f := range af.Filters ***REMOVED***
		filter, err := bson.MarshalWithRegistry(registry, f)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		filters = append(filters, filter)
	***REMOVED***
	return filters, nil
***REMOVED***

// ToArrayDocument builds a BSON array for the array filters CRUD option. If the registry for af is nil,
// bson.DefaultRegistry will be used when converting the filter interfaces to BSON.
func (af *ArrayFilters) ToArrayDocument() (bson.Raw, error) ***REMOVED***
	registry := af.Registry
	if registry == nil ***REMOVED***
		registry = bson.DefaultRegistry
	***REMOVED***

	idx, arr := bsoncore.AppendArrayStart(nil)
	for i, f := range af.Filters ***REMOVED***
		filter, err := bson.MarshalWithRegistry(registry, f)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		arr = bsoncore.AppendDocumentElement(arr, strconv.Itoa(i), filter)
	***REMOVED***
	arr, _ = bsoncore.AppendArrayEnd(arr, idx)
	return arr, nil
***REMOVED***

// MarshalError is returned when attempting to transform a value into a document
// results in an error.
type MarshalError struct ***REMOVED***
	Value interface***REMOVED******REMOVED***
	Err   error
***REMOVED***

// Error implements the error interface.
func (me MarshalError) Error() string ***REMOVED***
	return fmt.Sprintf("cannot transform type %s to a bson.Raw", reflect.TypeOf(me.Value))
***REMOVED***
