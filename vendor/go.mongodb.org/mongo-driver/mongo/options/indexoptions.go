// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"
)

// CreateIndexesOptions represents options that can be used to configure IndexView.CreateOne and IndexView.CreateMany
// operations.
type CreateIndexesOptions struct ***REMOVED***
	// The number of data-bearing members of a replica set, including the primary, that must complete the index builds
	// successfully before the primary marks the indexes as ready. This should either be a string or int32 value. The
	// semantics of the values are as follows:
	//
	// 1. String: specifies a tag. All members with that tag must complete the build.
	// 2. int: the number of members that must complete the build.
	// 3. "majority": A special value to indicate that more than half the nodes must complete the build.
	// 4. "votingMembers": A special value to indicate that all voting data-bearing nodes must complete.
	//
	// This option is only available on MongoDB versions >= 4.4. A client-side error will be returned if the option
	// is specified for MongoDB versions <= 4.2. The default value is nil, meaning that the server-side default will be
	// used. See dochub.mongodb.org/core/index-commit-quorum for more information.
	CommitQuorum interface***REMOVED******REMOVED***

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration
***REMOVED***

// CreateIndexes creates a new CreateIndexesOptions instance.
func CreateIndexes() *CreateIndexesOptions ***REMOVED***
	return &CreateIndexesOptions***REMOVED******REMOVED***
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (c *CreateIndexesOptions) SetMaxTime(d time.Duration) *CreateIndexesOptions ***REMOVED***
	c.MaxTime = &d
	return c
***REMOVED***

// SetCommitQuorumInt sets the value for the CommitQuorum field as an int32.
func (c *CreateIndexesOptions) SetCommitQuorumInt(quorum int32) *CreateIndexesOptions ***REMOVED***
	c.CommitQuorum = quorum
	return c
***REMOVED***

// SetCommitQuorumString sets the value for the CommitQuorum field as a string.
func (c *CreateIndexesOptions) SetCommitQuorumString(quorum string) *CreateIndexesOptions ***REMOVED***
	c.CommitQuorum = quorum
	return c
***REMOVED***

// SetCommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
func (c *CreateIndexesOptions) SetCommitQuorumMajority() *CreateIndexesOptions ***REMOVED***
	c.CommitQuorum = "majority"
	return c
***REMOVED***

// SetCommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
func (c *CreateIndexesOptions) SetCommitQuorumVotingMembers() *CreateIndexesOptions ***REMOVED***
	c.CommitQuorum = "votingMembers"
	return c
***REMOVED***

// MergeCreateIndexesOptions combines the given CreateIndexesOptions into a single CreateIndexesOptions in a last one
// wins fashion.
func MergeCreateIndexesOptions(opts ...*CreateIndexesOptions) *CreateIndexesOptions ***REMOVED***
	c := CreateIndexes()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			c.MaxTime = opt.MaxTime
		***REMOVED***
		if opt.CommitQuorum != nil ***REMOVED***
			c.CommitQuorum = opt.CommitQuorum
		***REMOVED***
	***REMOVED***

	return c
***REMOVED***

// DropIndexesOptions represents options that can be used to configure IndexView.DropOne and IndexView.DropAll
// operations.
type DropIndexesOptions struct ***REMOVED***
	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration
***REMOVED***

// DropIndexes creates a new DropIndexesOptions instance.
func DropIndexes() *DropIndexesOptions ***REMOVED***
	return &DropIndexesOptions***REMOVED******REMOVED***
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (d *DropIndexesOptions) SetMaxTime(duration time.Duration) *DropIndexesOptions ***REMOVED***
	d.MaxTime = &duration
	return d
***REMOVED***

// MergeDropIndexesOptions combines the given DropIndexesOptions into a single DropIndexesOptions in a last-one-wins
// fashion.
func MergeDropIndexesOptions(opts ...*DropIndexesOptions) *DropIndexesOptions ***REMOVED***
	c := DropIndexes()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			c.MaxTime = opt.MaxTime
		***REMOVED***
	***REMOVED***

	return c
***REMOVED***

// ListIndexesOptions represents options that can be used to configure an IndexView.List operation.
type ListIndexesOptions struct ***REMOVED***
	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	//
	// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout option may be used
	// in its place to control the amount of time that a single operation can run before returning an error. MaxTime
	// is ignored if Timeout is set on the client.
	MaxTime *time.Duration
***REMOVED***

// ListIndexes creates a new ListIndexesOptions instance.
func ListIndexes() *ListIndexesOptions ***REMOVED***
	return &ListIndexesOptions***REMOVED******REMOVED***
***REMOVED***

// SetBatchSize sets the value for the BatchSize field.
func (l *ListIndexesOptions) SetBatchSize(i int32) *ListIndexesOptions ***REMOVED***
	l.BatchSize = &i
	return l
***REMOVED***

// SetMaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (l *ListIndexesOptions) SetMaxTime(d time.Duration) *ListIndexesOptions ***REMOVED***
	l.MaxTime = &d
	return l
***REMOVED***

// MergeListIndexesOptions combines the given ListIndexesOptions instances into a single *ListIndexesOptions in a
// last-one-wins fashion.
func MergeListIndexesOptions(opts ...*ListIndexesOptions) *ListIndexesOptions ***REMOVED***
	c := ListIndexes()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.BatchSize != nil ***REMOVED***
			c.BatchSize = opt.BatchSize
		***REMOVED***
		if opt.MaxTime != nil ***REMOVED***
			c.MaxTime = opt.MaxTime
		***REMOVED***
	***REMOVED***

	return c
***REMOVED***

// IndexOptions represents options that can be used to configure a new index created through the IndexView.CreateOne
// or IndexView.CreateMany operations.
type IndexOptions struct ***REMOVED***
	// If true, the index will be built in the background on the server and will not block other tasks. The default
	// value is false.
	//
	// Deprecated: This option has been deprecated in MongoDB version 4.2.
	Background *bool

	// The length of time, in seconds, for documents to remain in the collection. The default value is 0, which means
	// that documents will remain in the collection until they're explicitly deleted or the collection is dropped.
	ExpireAfterSeconds *int32

	// The name of the index. The default value is "[field1]_[direction1]_[field2]_[direction2]...". For example, an
	// index with the specification ***REMOVED***name: 1, age: -1***REMOVED*** will be named "name_1_age_-1".
	Name *string

	// If true, the index will only reference documents that contain the fields specified in the index. The default is
	// false.
	Sparse *bool

	// Specifies the storage engine to use for the index. The value must be a document in the form
	// ***REMOVED***<storage engine name>: <options>***REMOVED***. The default value is nil, which means that the default storage engine
	// will be used. This option is only applicable for MongoDB versions >= 3.0 and is ignored for previous server
	// versions.
	StorageEngine interface***REMOVED******REMOVED***

	// If true, the collection will not accept insertion or update of documents where the index key value matches an
	// existing value in the index. The default is false.
	Unique *bool

	// The index version number, either 0 or 1.
	Version *int32

	// The language that determines the list of stop words and the rules for the stemmer and tokenizer. This option
	// is only applicable for text indexes and is ignored for other index types. The default value is "english".
	DefaultLanguage *string

	// The name of the field in the collection's documents that contains the override language for the document. This
	// option is only applicable for text indexes and is ignored for other index types. The default value is the value
	// of the DefaultLanguage option.
	LanguageOverride *string

	// The index version number for a text index. See https://www.mongodb.com/docs/manual/core/index-text/#text-versions for
	// information about different version numbers.
	TextVersion *int32

	// A document that contains field and weight pairs. The weight is an integer ranging from 1 to 99,999, inclusive,
	// indicating the significance of the field relative to the other indexed fields in terms of the score. This option
	// is only applicable for text indexes and is ignored for other index types. The default value is nil, which means
	// that every field will have a weight of 1.
	Weights interface***REMOVED******REMOVED***

	// The index version number for a 2D sphere index. See https://www.mongodb.com/docs/manual/core/2dsphere/#dsphere-v2 for
	// information about different version numbers.
	SphereVersion *int32

	// The precision of the stored geohash value of the location data. This option only applies to 2D indexes and is
	// ignored for other index types. The value must be between 1 and 32, inclusive. The default value is 26.
	Bits *int32

	// The upper inclusive boundary for longitude and latitude values. This option is only applicable to 2D indexes and
	// is ignored for other index types. The default value is 180.0.
	Max *float64

	// The lower inclusive boundary for longitude and latitude values. This option is only applicable to 2D indexes and
	// is ignored for other index types. The default value is -180.0.
	Min *float64

	// The number of units within which to group location values. Location values that are within BucketSize units of
	// each other will be grouped in the same bucket. This option is only applicable to geoHaystack indexes and is
	// ignored for other index types. The value must be greater than 0.
	BucketSize *int32

	// A document that defines which collection documents the index should reference. This option is only valid for
	// MongoDB versions >= 3.2 and is ignored for previous server versions.
	PartialFilterExpression interface***REMOVED******REMOVED***

	// The collation to use for string comparisons for the index. This option is only valid for MongoDB versions >= 3.4.
	// For previous server versions, the driver will return an error if this option is used.
	Collation *Collation

	// A document that defines the wildcard projection for the index.
	WildcardProjection interface***REMOVED******REMOVED***

	// If true, the index will exist on the target collection but will not be used by the query planner when executing
	// operations. This option is only valid for MongoDB versions >= 4.4. The default value is false.
	Hidden *bool
***REMOVED***

// Index creates a new IndexOptions instance.
func Index() *IndexOptions ***REMOVED***
	return &IndexOptions***REMOVED******REMOVED***
***REMOVED***

// SetBackground sets value for the Background field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.2.
func (i *IndexOptions) SetBackground(background bool) *IndexOptions ***REMOVED***
	i.Background = &background
	return i
***REMOVED***

// SetExpireAfterSeconds sets value for the ExpireAfterSeconds field.
func (i *IndexOptions) SetExpireAfterSeconds(seconds int32) *IndexOptions ***REMOVED***
	i.ExpireAfterSeconds = &seconds
	return i
***REMOVED***

// SetName sets the value for the Name field.
func (i *IndexOptions) SetName(name string) *IndexOptions ***REMOVED***
	i.Name = &name
	return i
***REMOVED***

// SetSparse sets the value of the Sparse field.
func (i *IndexOptions) SetSparse(sparse bool) *IndexOptions ***REMOVED***
	i.Sparse = &sparse
	return i
***REMOVED***

// SetStorageEngine sets the value for the StorageEngine field.
func (i *IndexOptions) SetStorageEngine(engine interface***REMOVED******REMOVED***) *IndexOptions ***REMOVED***
	i.StorageEngine = engine
	return i
***REMOVED***

// SetUnique sets the value for the Unique field.
func (i *IndexOptions) SetUnique(unique bool) *IndexOptions ***REMOVED***
	i.Unique = &unique
	return i
***REMOVED***

// SetVersion sets the value for the Version field.
func (i *IndexOptions) SetVersion(version int32) *IndexOptions ***REMOVED***
	i.Version = &version
	return i
***REMOVED***

// SetDefaultLanguage sets the value for the DefaultLanguage field.
func (i *IndexOptions) SetDefaultLanguage(language string) *IndexOptions ***REMOVED***
	i.DefaultLanguage = &language
	return i
***REMOVED***

// SetLanguageOverride sets the value of the LanguageOverride field.
func (i *IndexOptions) SetLanguageOverride(override string) *IndexOptions ***REMOVED***
	i.LanguageOverride = &override
	return i
***REMOVED***

// SetTextVersion sets the value for the TextVersion field.
func (i *IndexOptions) SetTextVersion(version int32) *IndexOptions ***REMOVED***
	i.TextVersion = &version
	return i
***REMOVED***

// SetWeights sets the value for the Weights field.
func (i *IndexOptions) SetWeights(weights interface***REMOVED******REMOVED***) *IndexOptions ***REMOVED***
	i.Weights = weights
	return i
***REMOVED***

// SetSphereVersion sets the value for the SphereVersion field.
func (i *IndexOptions) SetSphereVersion(version int32) *IndexOptions ***REMOVED***
	i.SphereVersion = &version
	return i
***REMOVED***

// SetBits sets the value for the Bits field.
func (i *IndexOptions) SetBits(bits int32) *IndexOptions ***REMOVED***
	i.Bits = &bits
	return i
***REMOVED***

// SetMax sets the value for the Max field.
func (i *IndexOptions) SetMax(max float64) *IndexOptions ***REMOVED***
	i.Max = &max
	return i
***REMOVED***

// SetMin sets the value for the Min field.
func (i *IndexOptions) SetMin(min float64) *IndexOptions ***REMOVED***
	i.Min = &min
	return i
***REMOVED***

// SetBucketSize sets the value for the BucketSize field
func (i *IndexOptions) SetBucketSize(bucketSize int32) *IndexOptions ***REMOVED***
	i.BucketSize = &bucketSize
	return i
***REMOVED***

// SetPartialFilterExpression sets the value for the PartialFilterExpression field.
func (i *IndexOptions) SetPartialFilterExpression(expression interface***REMOVED******REMOVED***) *IndexOptions ***REMOVED***
	i.PartialFilterExpression = expression
	return i
***REMOVED***

// SetCollation sets the value for the Collation field.
func (i *IndexOptions) SetCollation(collation *Collation) *IndexOptions ***REMOVED***
	i.Collation = collation
	return i
***REMOVED***

// SetWildcardProjection sets the value for the WildcardProjection field.
func (i *IndexOptions) SetWildcardProjection(wildcardProjection interface***REMOVED******REMOVED***) *IndexOptions ***REMOVED***
	i.WildcardProjection = wildcardProjection
	return i
***REMOVED***

// SetHidden sets the value for the Hidden field.
func (i *IndexOptions) SetHidden(hidden bool) *IndexOptions ***REMOVED***
	i.Hidden = &hidden
	return i
***REMOVED***

// MergeIndexOptions combines the given IndexOptions into a single IndexOptions in a last-one-wins fashion.
func MergeIndexOptions(opts ...*IndexOptions) *IndexOptions ***REMOVED***
	i := Index()

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.Background != nil ***REMOVED***
			i.Background = opt.Background
		***REMOVED***
		if opt.ExpireAfterSeconds != nil ***REMOVED***
			i.ExpireAfterSeconds = opt.ExpireAfterSeconds
		***REMOVED***
		if opt.Name != nil ***REMOVED***
			i.Name = opt.Name
		***REMOVED***
		if opt.Sparse != nil ***REMOVED***
			i.Sparse = opt.Sparse
		***REMOVED***
		if opt.StorageEngine != nil ***REMOVED***
			i.StorageEngine = opt.StorageEngine
		***REMOVED***
		if opt.Unique != nil ***REMOVED***
			i.Unique = opt.Unique
		***REMOVED***
		if opt.Version != nil ***REMOVED***
			i.Version = opt.Version
		***REMOVED***
		if opt.DefaultLanguage != nil ***REMOVED***
			i.DefaultLanguage = opt.DefaultLanguage
		***REMOVED***
		if opt.LanguageOverride != nil ***REMOVED***
			i.LanguageOverride = opt.LanguageOverride
		***REMOVED***
		if opt.TextVersion != nil ***REMOVED***
			i.TextVersion = opt.TextVersion
		***REMOVED***
		if opt.Weights != nil ***REMOVED***
			i.Weights = opt.Weights
		***REMOVED***
		if opt.SphereVersion != nil ***REMOVED***
			i.SphereVersion = opt.SphereVersion
		***REMOVED***
		if opt.Bits != nil ***REMOVED***
			i.Bits = opt.Bits
		***REMOVED***
		if opt.Max != nil ***REMOVED***
			i.Max = opt.Max
		***REMOVED***
		if opt.Min != nil ***REMOVED***
			i.Min = opt.Min
		***REMOVED***
		if opt.BucketSize != nil ***REMOVED***
			i.BucketSize = opt.BucketSize
		***REMOVED***
		if opt.PartialFilterExpression != nil ***REMOVED***
			i.PartialFilterExpression = opt.PartialFilterExpression
		***REMOVED***
		if opt.Collation != nil ***REMOVED***
			i.Collation = opt.Collation
		***REMOVED***
		if opt.WildcardProjection != nil ***REMOVED***
			i.WildcardProjection = opt.WildcardProjection
		***REMOVED***
		if opt.Hidden != nil ***REMOVED***
			i.Hidden = opt.Hidden
		***REMOVED***
	***REMOVED***

	return i
***REMOVED***
