// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package readpref

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/tag"
)

// ErrInvalidTagSet indicates that an invalid set of tags was specified.
var ErrInvalidTagSet = errors.New("an even number of tags must be specified")

// Option configures a read preference
type Option func(*ReadPref) error

// WithMaxStaleness sets the maximum staleness a
// server is allowed.
func WithMaxStaleness(ms time.Duration) Option ***REMOVED***
	return func(rp *ReadPref) error ***REMOVED***
		rp.maxStaleness = ms
		rp.maxStalenessSet = true
		return nil
	***REMOVED***
***REMOVED***

// WithTags sets a single tag set used to match
// a server. The last call to WithTags or WithTagSets
// overrides all previous calls to either method.
func WithTags(tags ...string) Option ***REMOVED***
	return func(rp *ReadPref) error ***REMOVED***
		length := len(tags)
		if length < 2 || length%2 != 0 ***REMOVED***
			return ErrInvalidTagSet
		***REMOVED***

		tagset := make(tag.Set, 0, length/2)

		for i := 1; i < length; i += 2 ***REMOVED***
			tagset = append(tagset, tag.Tag***REMOVED***Name: tags[i-1], Value: tags[i]***REMOVED***)
		***REMOVED***

		return WithTagSets(tagset)(rp)
	***REMOVED***
***REMOVED***

// WithTagSets sets the tag sets used to match
// a server. The last call to WithTags or WithTagSets
// overrides all previous calls to either method.
func WithTagSets(tagSets ...tag.Set) Option ***REMOVED***
	return func(rp *ReadPref) error ***REMOVED***
		rp.tagSets = tagSets
		return nil
	***REMOVED***
***REMOVED***

// WithHedgeEnabled specifies whether or not hedged reads should be enabled in the server. This feature requires MongoDB
// server version 4.4 or higher. For more information about hedged reads, see
// https://www.mongodb.com/docs/manual/core/sharded-cluster-query-router/#mongos-hedged-reads. If not specified, the default
// is to not send a value to the server, which will result in the server defaults being used.
func WithHedgeEnabled(hedgeEnabled bool) Option ***REMOVED***
	return func(rp *ReadPref) error ***REMOVED***
		rp.hedgeEnabled = &hedgeEnabled
		return nil
	***REMOVED***
***REMOVED***
