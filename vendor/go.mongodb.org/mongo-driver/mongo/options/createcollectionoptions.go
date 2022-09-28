// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DefaultIndexOptions represents the default options for a collection to apply on new indexes. This type can be used
// when creating a new collection through the CreateCollectionOptions.SetDefaultIndexOptions method.
type DefaultIndexOptions struct ***REMOVED***
	// Specifies the storage engine to use for the index. The value must be a document in the form
	// ***REMOVED***<storage engine name>: <options>***REMOVED***. The default value is nil, which means that the default storage engine
	// will be used.
	StorageEngine interface***REMOVED******REMOVED***
***REMOVED***

// DefaultIndex creates a new DefaultIndexOptions instance.
func DefaultIndex() *DefaultIndexOptions ***REMOVED***
	return &DefaultIndexOptions***REMOVED******REMOVED***
***REMOVED***

// SetStorageEngine sets the value for the StorageEngine field.
func (d *DefaultIndexOptions) SetStorageEngine(storageEngine interface***REMOVED******REMOVED***) *DefaultIndexOptions ***REMOVED***
	d.StorageEngine = storageEngine
	return d
***REMOVED***

// TimeSeriesOptions specifies options on a time-series collection.
type TimeSeriesOptions struct ***REMOVED***
	// Name of the top-level field to be used for time. Inserted documents must have this field,
	// and the field must be of the BSON UTC datetime type (0x9).
	TimeField string

	// Optional name of the top-level field describing the series. This field is used to group
	// related data and may be of any BSON type, except for array. This name may not be the same
	// as the TimeField or _id.
	MetaField *string

	// Optional string specifying granularity of time-series data. Allowed granularity options are
	// "seconds", "minutes" and "hours".
	Granularity *string
***REMOVED***

// TimeSeries creates a new TimeSeriesOptions instance.
func TimeSeries() *TimeSeriesOptions ***REMOVED***
	return &TimeSeriesOptions***REMOVED******REMOVED***
***REMOVED***

// SetTimeField sets the value for the TimeField.
func (tso *TimeSeriesOptions) SetTimeField(timeField string) *TimeSeriesOptions ***REMOVED***
	tso.TimeField = timeField
	return tso
***REMOVED***

// SetMetaField sets the value for the MetaField.
func (tso *TimeSeriesOptions) SetMetaField(metaField string) *TimeSeriesOptions ***REMOVED***
	tso.MetaField = &metaField
	return tso
***REMOVED***

// SetGranularity sets the value for Granularity.
func (tso *TimeSeriesOptions) SetGranularity(granularity string) *TimeSeriesOptions ***REMOVED***
	tso.Granularity = &granularity
	return tso
***REMOVED***

// CreateCollectionOptions represents options that can be used to configure a CreateCollection operation.
type CreateCollectionOptions struct ***REMOVED***
	// Specifies if the collection is capped (see https://www.mongodb.com/docs/manual/core/capped-collections/). If true,
	// the SizeInBytes option must also be specified. The default value is false.
	Capped *bool

	// Specifies the default collation for the new collection. This option is only valid for MongoDB versions >= 3.4.
	// For previous server versions, the driver will return an error if this option is used. The default value is nil.
	Collation *Collation

	// Specifies how change streams opened against the collection can return pre- and post-images of updated
	// documents. The value must be a document in the form ***REMOVED***<option name>: <options>***REMOVED***. This option is only valid for
	// MongoDB versions >= 6.0. The default value is nil, which means that change streams opened against the collection
	// will not return pre- and post-images of updated documents in any way.
	ChangeStreamPreAndPostImages interface***REMOVED******REMOVED***

	// Specifies a default configuration for indexes on the collection. This option is only valid for MongoDB versions
	// >= 3.4. The default value is nil, meaning indexes will be configured using server defaults.
	DefaultIndexOptions *DefaultIndexOptions

	// Specifies the maximum number of documents allowed in a capped collection. The limit specified by the SizeInBytes
	// option takes precedence over this option. If a capped collection reaches its size limit, old documents will be
	// removed, regardless of the number of documents in the collection. The default value is 0, meaning the maximum
	// number of documents is unbounded.
	MaxDocuments *int64

	// Specifies the maximum size in bytes for a capped collection. The default value is 0.
	SizeInBytes *int64

	// Specifies the storage engine to use for the index. The value must be a document in the form
	// ***REMOVED***<storage engine name>: <options>***REMOVED***. The default value is nil, which means that the default storage engine
	// will be used.
	StorageEngine interface***REMOVED******REMOVED***

	// Specifies what should happen if a document being inserted does not pass validation. Valid values are "error" and
	// "warn". See https://www.mongodb.com/docs/manual/core/schema-validation/#accept-or-reject-invalid-documents for more
	// information. This option is only valid for MongoDB versions >= 3.2. The default value is "error".
	ValidationAction *string

	// Specifies how strictly the server applies validation rules to existing documents in the collection during update
	// operations. Valid values are "off", "strict", and "moderate". See
	// https://www.mongodb.com/docs/manual/core/schema-validation/#existing-documents for more information. This option is
	// only valid for MongoDB versions >= 3.2. The default value is "strict".
	ValidationLevel *string

	// A document specifying validation rules for the collection. See
	// https://www.mongodb.com/docs/manual/core/schema-validation/ for more information about schema validation. This option
	// is only valid for MongoDB versions >= 3.2. The default value is nil, meaning no validator will be used for the
	// collection.
	Validator interface***REMOVED******REMOVED***

	// Value indicating after how many seconds old time-series data should be deleted. See
	// https://www.mongodb.com/docs/manual/reference/command/create/ for supported options, and
	// https://www.mongodb.com/docs/manual/core/timeseries-collections/ for more information on time-series
	// collections.
	//
	// This option is only valid for MongoDB versions >= 5.0
	ExpireAfterSeconds *int64

	// Options for specifying a time-series collection. See
	// https://www.mongodb.com/docs/manual/reference/command/create/ for supported options, and
	// https://www.mongodb.com/docs/manual/core/timeseries-collections/ for more information on time-series
	// collections.
	//
	// This option is only valid for MongoDB versions >= 5.0
	TimeSeriesOptions *TimeSeriesOptions

	// EncryptedFields configures encrypted fields.
	//
	// This option is only valid for MongoDB versions >= 6.0
	EncryptedFields interface***REMOVED******REMOVED***

	// ClusteredIndex is used to create a collection with a clustered index.
	//
	// This option is only valid for MongoDB versions >= 5.3
	ClusteredIndex interface***REMOVED******REMOVED***
***REMOVED***

// CreateCollection creates a new CreateCollectionOptions instance.
func CreateCollection() *CreateCollectionOptions ***REMOVED***
	return &CreateCollectionOptions***REMOVED******REMOVED***
***REMOVED***

// SetCapped sets the value for the Capped field.
func (c *CreateCollectionOptions) SetCapped(capped bool) *CreateCollectionOptions ***REMOVED***
	c.Capped = &capped
	return c
***REMOVED***

// SetCollation sets the value for the Collation field.
func (c *CreateCollectionOptions) SetCollation(collation *Collation) *CreateCollectionOptions ***REMOVED***
	c.Collation = collation
	return c
***REMOVED***

// SetChangeStreamPreAndPostImages sets the value for the ChangeStreamPreAndPostImages field.
func (c *CreateCollectionOptions) SetChangeStreamPreAndPostImages(csppi interface***REMOVED******REMOVED***) *CreateCollectionOptions ***REMOVED***
	c.ChangeStreamPreAndPostImages = &csppi
	return c
***REMOVED***

// SetDefaultIndexOptions sets the value for the DefaultIndexOptions field.
func (c *CreateCollectionOptions) SetDefaultIndexOptions(opts *DefaultIndexOptions) *CreateCollectionOptions ***REMOVED***
	c.DefaultIndexOptions = opts
	return c
***REMOVED***

// SetMaxDocuments sets the value for the MaxDocuments field.
func (c *CreateCollectionOptions) SetMaxDocuments(max int64) *CreateCollectionOptions ***REMOVED***
	c.MaxDocuments = &max
	return c
***REMOVED***

// SetSizeInBytes sets the value for the SizeInBytes field.
func (c *CreateCollectionOptions) SetSizeInBytes(size int64) *CreateCollectionOptions ***REMOVED***
	c.SizeInBytes = &size
	return c
***REMOVED***

// SetStorageEngine sets the value for the StorageEngine field.
func (c *CreateCollectionOptions) SetStorageEngine(storageEngine interface***REMOVED******REMOVED***) *CreateCollectionOptions ***REMOVED***
	c.StorageEngine = &storageEngine
	return c
***REMOVED***

// SetValidationAction sets the value for the ValidationAction field.
func (c *CreateCollectionOptions) SetValidationAction(action string) *CreateCollectionOptions ***REMOVED***
	c.ValidationAction = &action
	return c
***REMOVED***

// SetValidationLevel sets the value for the ValidationLevel field.
func (c *CreateCollectionOptions) SetValidationLevel(level string) *CreateCollectionOptions ***REMOVED***
	c.ValidationLevel = &level
	return c
***REMOVED***

// SetValidator sets the value for the Validator field.
func (c *CreateCollectionOptions) SetValidator(validator interface***REMOVED******REMOVED***) *CreateCollectionOptions ***REMOVED***
	c.Validator = validator
	return c
***REMOVED***

// SetExpireAfterSeconds sets the value for the ExpireAfterSeconds field.
func (c *CreateCollectionOptions) SetExpireAfterSeconds(eas int64) *CreateCollectionOptions ***REMOVED***
	c.ExpireAfterSeconds = &eas
	return c
***REMOVED***

// SetTimeSeriesOptions sets the options for time-series collections.
func (c *CreateCollectionOptions) SetTimeSeriesOptions(timeSeriesOpts *TimeSeriesOptions) *CreateCollectionOptions ***REMOVED***
	c.TimeSeriesOptions = timeSeriesOpts
	return c
***REMOVED***

// SetEncryptedFields sets the encrypted fields for encrypted collections.
func (c *CreateCollectionOptions) SetEncryptedFields(encryptedFields interface***REMOVED******REMOVED***) *CreateCollectionOptions ***REMOVED***
	c.EncryptedFields = encryptedFields
	return c
***REMOVED***

// SetClusteredIndex sets the value for the ClusteredIndex field.
func (c *CreateCollectionOptions) SetClusteredIndex(clusteredIndex interface***REMOVED******REMOVED***) *CreateCollectionOptions ***REMOVED***
	c.ClusteredIndex = clusteredIndex
	return c
***REMOVED***

// MergeCreateCollectionOptions combines the given CreateCollectionOptions instances into a single
// CreateCollectionOptions in a last-one-wins fashion.
func MergeCreateCollectionOptions(opts ...*CreateCollectionOptions) *CreateCollectionOptions ***REMOVED***
	cc := CreateCollection()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***

		if opt.Capped != nil ***REMOVED***
			cc.Capped = opt.Capped
		***REMOVED***
		if opt.Collation != nil ***REMOVED***
			cc.Collation = opt.Collation
		***REMOVED***
		if opt.ChangeStreamPreAndPostImages != nil ***REMOVED***
			cc.ChangeStreamPreAndPostImages = opt.ChangeStreamPreAndPostImages
		***REMOVED***
		if opt.DefaultIndexOptions != nil ***REMOVED***
			cc.DefaultIndexOptions = opt.DefaultIndexOptions
		***REMOVED***
		if opt.MaxDocuments != nil ***REMOVED***
			cc.MaxDocuments = opt.MaxDocuments
		***REMOVED***
		if opt.SizeInBytes != nil ***REMOVED***
			cc.SizeInBytes = opt.SizeInBytes
		***REMOVED***
		if opt.StorageEngine != nil ***REMOVED***
			cc.StorageEngine = opt.StorageEngine
		***REMOVED***
		if opt.ValidationAction != nil ***REMOVED***
			cc.ValidationAction = opt.ValidationAction
		***REMOVED***
		if opt.ValidationLevel != nil ***REMOVED***
			cc.ValidationLevel = opt.ValidationLevel
		***REMOVED***
		if opt.Validator != nil ***REMOVED***
			cc.Validator = opt.Validator
		***REMOVED***
		if opt.ExpireAfterSeconds != nil ***REMOVED***
			cc.ExpireAfterSeconds = opt.ExpireAfterSeconds
		***REMOVED***
		if opt.TimeSeriesOptions != nil ***REMOVED***
			cc.TimeSeriesOptions = opt.TimeSeriesOptions
		***REMOVED***
		if opt.EncryptedFields != nil ***REMOVED***
			cc.EncryptedFields = opt.EncryptedFields
		***REMOVED***
		if opt.ClusteredIndex != nil ***REMOVED***
			cc.ClusteredIndex = opt.ClusteredIndex
		***REMOVED***
	***REMOVED***

	return cc
***REMOVED***

// CreateViewOptions represents options that can be used to configure a CreateView operation.
type CreateViewOptions struct ***REMOVED***
	// Specifies the default collation for the new collection. This option is only valid for MongoDB versions >= 3.4.
	// For previous server versions, the driver will return an error if this option is used. The default value is nil.
	Collation *Collation
***REMOVED***

// CreateView creates an new CreateViewOptions instance.
func CreateView() *CreateViewOptions ***REMOVED***
	return &CreateViewOptions***REMOVED******REMOVED***
***REMOVED***

// SetCollation sets the value for the Collation field.
func (c *CreateViewOptions) SetCollation(collation *Collation) *CreateViewOptions ***REMOVED***
	c.Collation = collation
	return c
***REMOVED***

// MergeCreateViewOptions combines the given CreateViewOptions instances into a single CreateViewOptions in a
// last-one-wins fashion.
func MergeCreateViewOptions(opts ...*CreateViewOptions) *CreateViewOptions ***REMOVED***
	cv := CreateView()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***

		if opt.Collation != nil ***REMOVED***
			cv.Collation = opt.Collation
		***REMOVED***
	***REMOVED***

	return cv
***REMOVED***
