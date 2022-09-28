// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WriteModel is an interface implemented by models that can be used in a BulkWrite operation. Each WriteModel
// represents a write.
//
// This interface is implemented by InsertOneModel, DeleteOneModel, DeleteManyModel, ReplaceOneModel, UpdateOneModel,
// and UpdateManyModel. Custom implementations of this interface must not be used.
type WriteModel interface ***REMOVED***
	writeModel()
***REMOVED***

// InsertOneModel is used to insert a single document in a BulkWrite operation.
type InsertOneModel struct ***REMOVED***
	Document interface***REMOVED******REMOVED***
***REMOVED***

// NewInsertOneModel creates a new InsertOneModel.
func NewInsertOneModel() *InsertOneModel ***REMOVED***
	return &InsertOneModel***REMOVED******REMOVED***
***REMOVED***

// SetDocument specifies the document to be inserted. The document cannot be nil. If it does not have an _id field when
// transformed into BSON, one will be added automatically to the marshalled document. The original document will not be
// modified.
func (iom *InsertOneModel) SetDocument(doc interface***REMOVED******REMOVED***) *InsertOneModel ***REMOVED***
	iom.Document = doc
	return iom
***REMOVED***

func (*InsertOneModel) writeModel() ***REMOVED******REMOVED***

// DeleteOneModel is used to delete at most one document in a BulkWriteOperation.
type DeleteOneModel struct ***REMOVED***
	Filter    interface***REMOVED******REMOVED***
	Collation *options.Collation
	Hint      interface***REMOVED******REMOVED***
***REMOVED***

// NewDeleteOneModel creates a new DeleteOneModel.
func NewDeleteOneModel() *DeleteOneModel ***REMOVED***
	return &DeleteOneModel***REMOVED******REMOVED***
***REMOVED***

// SetFilter specifies a filter to use to select the document to delete. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (dom *DeleteOneModel) SetFilter(filter interface***REMOVED******REMOVED***) *DeleteOneModel ***REMOVED***
	dom.Filter = filter
	return dom
***REMOVED***

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (dom *DeleteOneModel) SetCollation(collation *options.Collation) *DeleteOneModel ***REMOVED***
	dom.Collation = collation
	return dom
***REMOVED***

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (dom *DeleteOneModel) SetHint(hint interface***REMOVED******REMOVED***) *DeleteOneModel ***REMOVED***
	dom.Hint = hint
	return dom
***REMOVED***

func (*DeleteOneModel) writeModel() ***REMOVED******REMOVED***

// DeleteManyModel is used to delete multiple documents in a BulkWrite operation.
type DeleteManyModel struct ***REMOVED***
	Filter    interface***REMOVED******REMOVED***
	Collation *options.Collation
	Hint      interface***REMOVED******REMOVED***
***REMOVED***

// NewDeleteManyModel creates a new DeleteManyModel.
func NewDeleteManyModel() *DeleteManyModel ***REMOVED***
	return &DeleteManyModel***REMOVED******REMOVED***
***REMOVED***

// SetFilter specifies a filter to use to select documents to delete. The filter must be a document containing query
// operators. It cannot be nil.
func (dmm *DeleteManyModel) SetFilter(filter interface***REMOVED******REMOVED***) *DeleteManyModel ***REMOVED***
	dmm.Filter = filter
	return dmm
***REMOVED***

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (dmm *DeleteManyModel) SetCollation(collation *options.Collation) *DeleteManyModel ***REMOVED***
	dmm.Collation = collation
	return dmm
***REMOVED***

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (dmm *DeleteManyModel) SetHint(hint interface***REMOVED******REMOVED***) *DeleteManyModel ***REMOVED***
	dmm.Hint = hint
	return dmm
***REMOVED***

func (*DeleteManyModel) writeModel() ***REMOVED******REMOVED***

// ReplaceOneModel is used to replace at most one document in a BulkWrite operation.
type ReplaceOneModel struct ***REMOVED***
	Collation   *options.Collation
	Upsert      *bool
	Filter      interface***REMOVED******REMOVED***
	Replacement interface***REMOVED******REMOVED***
	Hint        interface***REMOVED******REMOVED***
***REMOVED***

// NewReplaceOneModel creates a new ReplaceOneModel.
func NewReplaceOneModel() *ReplaceOneModel ***REMOVED***
	return &ReplaceOneModel***REMOVED******REMOVED***
***REMOVED***

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (rom *ReplaceOneModel) SetHint(hint interface***REMOVED******REMOVED***) *ReplaceOneModel ***REMOVED***
	rom.Hint = hint
	return rom
***REMOVED***

// SetFilter specifies a filter to use to select the document to replace. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (rom *ReplaceOneModel) SetFilter(filter interface***REMOVED******REMOVED***) *ReplaceOneModel ***REMOVED***
	rom.Filter = filter
	return rom
***REMOVED***

// SetReplacement specifies a document that will be used to replace the selected document. It cannot be nil and cannot
// contain any update operators (https://www.mongodb.com/docs/manual/reference/operator/update/).
func (rom *ReplaceOneModel) SetReplacement(rep interface***REMOVED******REMOVED***) *ReplaceOneModel ***REMOVED***
	rom.Replacement = rep
	return rom
***REMOVED***

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (rom *ReplaceOneModel) SetCollation(collation *options.Collation) *ReplaceOneModel ***REMOVED***
	rom.Collation = collation
	return rom
***REMOVED***

// SetUpsert specifies whether or not the replacement document should be inserted if no document matching the filter is
// found. If an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (rom *ReplaceOneModel) SetUpsert(upsert bool) *ReplaceOneModel ***REMOVED***
	rom.Upsert = &upsert
	return rom
***REMOVED***

func (*ReplaceOneModel) writeModel() ***REMOVED******REMOVED***

// UpdateOneModel is used to update at most one document in a BulkWrite operation.
type UpdateOneModel struct ***REMOVED***
	Collation    *options.Collation
	Upsert       *bool
	Filter       interface***REMOVED******REMOVED***
	Update       interface***REMOVED******REMOVED***
	ArrayFilters *options.ArrayFilters
	Hint         interface***REMOVED******REMOVED***
***REMOVED***

// NewUpdateOneModel creates a new UpdateOneModel.
func NewUpdateOneModel() *UpdateOneModel ***REMOVED***
	return &UpdateOneModel***REMOVED******REMOVED***
***REMOVED***

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (uom *UpdateOneModel) SetHint(hint interface***REMOVED******REMOVED***) *UpdateOneModel ***REMOVED***
	uom.Hint = hint
	return uom
***REMOVED***

// SetFilter specifies a filter to use to select the document to update. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (uom *UpdateOneModel) SetFilter(filter interface***REMOVED******REMOVED***) *UpdateOneModel ***REMOVED***
	uom.Filter = filter
	return uom
***REMOVED***

// SetUpdate specifies the modifications to be made to the selected document. The value must be a document containing
// update operators (https://www.mongodb.com/docs/manual/reference/operator/update/). It cannot be nil or empty.
func (uom *UpdateOneModel) SetUpdate(update interface***REMOVED******REMOVED***) *UpdateOneModel ***REMOVED***
	uom.Update = update
	return uom
***REMOVED***

// SetArrayFilters specifies a set of filters to determine which elements should be modified when updating an array
// field.
func (uom *UpdateOneModel) SetArrayFilters(filters options.ArrayFilters) *UpdateOneModel ***REMOVED***
	uom.ArrayFilters = &filters
	return uom
***REMOVED***

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (uom *UpdateOneModel) SetCollation(collation *options.Collation) *UpdateOneModel ***REMOVED***
	uom.Collation = collation
	return uom
***REMOVED***

// SetUpsert specifies whether or not a new document should be inserted if no document matching the filter is found. If
// an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (uom *UpdateOneModel) SetUpsert(upsert bool) *UpdateOneModel ***REMOVED***
	uom.Upsert = &upsert
	return uom
***REMOVED***

func (*UpdateOneModel) writeModel() ***REMOVED******REMOVED***

// UpdateManyModel is used to update multiple documents in a BulkWrite operation.
type UpdateManyModel struct ***REMOVED***
	Collation    *options.Collation
	Upsert       *bool
	Filter       interface***REMOVED******REMOVED***
	Update       interface***REMOVED******REMOVED***
	ArrayFilters *options.ArrayFilters
	Hint         interface***REMOVED******REMOVED***
***REMOVED***

// NewUpdateManyModel creates a new UpdateManyModel.
func NewUpdateManyModel() *UpdateManyModel ***REMOVED***
	return &UpdateManyModel***REMOVED******REMOVED***
***REMOVED***

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (umm *UpdateManyModel) SetHint(hint interface***REMOVED******REMOVED***) *UpdateManyModel ***REMOVED***
	umm.Hint = hint
	return umm
***REMOVED***

// SetFilter specifies a filter to use to select documents to update. The filter must be a document containing query
// operators. It cannot be nil.
func (umm *UpdateManyModel) SetFilter(filter interface***REMOVED******REMOVED***) *UpdateManyModel ***REMOVED***
	umm.Filter = filter
	return umm
***REMOVED***

// SetUpdate specifies the modifications to be made to the selected documents. The value must be a document containing
// update operators (https://www.mongodb.com/docs/manual/reference/operator/update/). It cannot be nil or empty.
func (umm *UpdateManyModel) SetUpdate(update interface***REMOVED******REMOVED***) *UpdateManyModel ***REMOVED***
	umm.Update = update
	return umm
***REMOVED***

// SetArrayFilters specifies a set of filters to determine which elements should be modified when updating an array
// field.
func (umm *UpdateManyModel) SetArrayFilters(filters options.ArrayFilters) *UpdateManyModel ***REMOVED***
	umm.ArrayFilters = &filters
	return umm
***REMOVED***

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (umm *UpdateManyModel) SetCollation(collation *options.Collation) *UpdateManyModel ***REMOVED***
	umm.Collation = collation
	return umm
***REMOVED***

// SetUpsert specifies whether or not a new document should be inserted if no document matching the filter is found. If
// an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (umm *UpdateManyModel) SetUpsert(upsert bool) *UpdateManyModel ***REMOVED***
	umm.Upsert = &upsert
	return umm
***REMOVED***

func (*UpdateManyModel) writeModel() ***REMOVED******REMOVED***
