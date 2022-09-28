// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

// IndexOptionsBuilder specifies options for a new index.
//
// Deprecated: Use the IndexOptions type in the mongo/options package instead.
type IndexOptionsBuilder struct ***REMOVED***
	document bson.D
***REMOVED***

// NewIndexOptionsBuilder creates a new IndexOptionsBuilder.
//
// Deprecated: Use the Index function in mongo/options instead.
func NewIndexOptionsBuilder() *IndexOptionsBuilder ***REMOVED***
	return &IndexOptionsBuilder***REMOVED******REMOVED***
***REMOVED***

// Background specifies a value for the background option.
//
// Deprecated: Use the IndexOptions.SetBackground function in mongo/options instead.
func (iob *IndexOptionsBuilder) Background(background bool) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"background", background***REMOVED***)
	return iob
***REMOVED***

// ExpireAfterSeconds specifies a value for the expireAfterSeconds option.
//
// Deprecated: Use the IndexOptions.SetExpireAfterSeconds function in mongo/options instead.
func (iob *IndexOptionsBuilder) ExpireAfterSeconds(expireAfterSeconds int32) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"expireAfterSeconds", expireAfterSeconds***REMOVED***)
	return iob
***REMOVED***

// Name specifies a value for the name option.
//
// Deprecated: Use the IndexOptions.SetName function in mongo/options instead.
func (iob *IndexOptionsBuilder) Name(name string) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"name", name***REMOVED***)
	return iob
***REMOVED***

// Sparse specifies a value for the sparse option.
//
// Deprecated: Use the IndexOptions.SetSparse function in mongo/options instead.
func (iob *IndexOptionsBuilder) Sparse(sparse bool) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"sparse", sparse***REMOVED***)
	return iob
***REMOVED***

// StorageEngine specifies a value for the storageEngine option.
//
// Deprecated: Use the IndexOptions.SetStorageEngine function in mongo/options instead.
func (iob *IndexOptionsBuilder) StorageEngine(storageEngine interface***REMOVED******REMOVED***) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"storageEngine", storageEngine***REMOVED***)
	return iob
***REMOVED***

// Unique specifies a value for the unique option.
//
// Deprecated: Use the IndexOptions.SetUnique function in mongo/options instead.
func (iob *IndexOptionsBuilder) Unique(unique bool) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"unique", unique***REMOVED***)
	return iob
***REMOVED***

// Version specifies a value for the version option.
//
// Deprecated: Use the IndexOptions.SetVersion function in mongo/options instead.
func (iob *IndexOptionsBuilder) Version(version int32) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"v", version***REMOVED***)
	return iob
***REMOVED***

// DefaultLanguage specifies a value for the default_language option.
//
// Deprecated: Use the IndexOptions.SetDefaultLanguage function in mongo/options instead.
func (iob *IndexOptionsBuilder) DefaultLanguage(defaultLanguage string) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"default_language", defaultLanguage***REMOVED***)
	return iob
***REMOVED***

// LanguageOverride specifies a value for the language_override option.
//
// Deprecated: Use the IndexOptions.SetLanguageOverride function in mongo/options instead.
func (iob *IndexOptionsBuilder) LanguageOverride(languageOverride string) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"language_override", languageOverride***REMOVED***)
	return iob
***REMOVED***

// TextVersion specifies a value for the textIndexVersion option.
//
// Deprecated: Use the IndexOptions.SetTextVersion function in mongo/options instead.
func (iob *IndexOptionsBuilder) TextVersion(textVersion int32) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"textIndexVersion", textVersion***REMOVED***)
	return iob
***REMOVED***

// Weights specifies a value for the weights option.
//
// Deprecated: Use the IndexOptions.SetWeights function in mongo/options instead.
func (iob *IndexOptionsBuilder) Weights(weights interface***REMOVED******REMOVED***) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"weights", weights***REMOVED***)
	return iob
***REMOVED***

// SphereVersion specifies a value for the 2dsphereIndexVersion option.
//
// Deprecated: Use the IndexOptions.SetSphereVersion function in mongo/options instead.
func (iob *IndexOptionsBuilder) SphereVersion(sphereVersion int32) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"2dsphereIndexVersion", sphereVersion***REMOVED***)
	return iob
***REMOVED***

// Bits specifies a value for the bits option.
//
// Deprecated: Use the IndexOptions.SetBits function in mongo/options instead.
func (iob *IndexOptionsBuilder) Bits(bits int32) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"bits", bits***REMOVED***)
	return iob
***REMOVED***

// Max specifies a value for the max option.
//
// Deprecated: Use the IndexOptions.SetMax function in mongo/options instead.
func (iob *IndexOptionsBuilder) Max(max float64) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"max", max***REMOVED***)
	return iob
***REMOVED***

// Min specifies a value for the min option.
//
// Deprecated: Use the IndexOptions.SetMin function in mongo/options instead.
func (iob *IndexOptionsBuilder) Min(min float64) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"min", min***REMOVED***)
	return iob
***REMOVED***

// BucketSize specifies a value for the bucketSize option.
//
// Deprecated: Use the IndexOptions.SetBucketSize function in mongo/options instead.
func (iob *IndexOptionsBuilder) BucketSize(bucketSize int32) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"bucketSize", bucketSize***REMOVED***)
	return iob
***REMOVED***

// PartialFilterExpression specifies a value for the partialFilterExpression option.
//
// Deprecated: Use the IndexOptions.SetPartialFilterExpression function in mongo/options instead.
func (iob *IndexOptionsBuilder) PartialFilterExpression(partialFilterExpression interface***REMOVED******REMOVED***) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"partialFilterExpression", partialFilterExpression***REMOVED***)
	return iob
***REMOVED***

// Collation specifies a value for the collation option.
//
// Deprecated: Use the IndexOptions.SetCollation function in mongo/options instead.
func (iob *IndexOptionsBuilder) Collation(collation interface***REMOVED******REMOVED***) *IndexOptionsBuilder ***REMOVED***
	iob.document = append(iob.document, bson.E***REMOVED***"collation", collation***REMOVED***)
	return iob
***REMOVED***

// Build finishes constructing an the builder.
//
// Deprecated: Use the IndexOptions type in the mongo/options package instead.
func (iob *IndexOptionsBuilder) Build() bson.D ***REMOVED***
	return iob.document
***REMOVED***
