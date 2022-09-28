// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package readpref defines read preferences for MongoDB queries.
package readpref // import "go.mongodb.org/mongo-driver/mongo/readpref"

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/tag"
)

var (
	errInvalidReadPreference = errors.New("can not specify tags, max staleness, or hedge with mode primary")
)

var primary = ReadPref***REMOVED***mode: PrimaryMode***REMOVED***

// Primary constructs a read preference with a PrimaryMode.
func Primary() *ReadPref ***REMOVED***
	return &primary
***REMOVED***

// PrimaryPreferred constructs a read preference with a PrimaryPreferredMode.
func PrimaryPreferred(opts ...Option) *ReadPref ***REMOVED***
	// New only returns an error with a mode of Primary
	rp, _ := New(PrimaryPreferredMode, opts...)
	return rp
***REMOVED***

// SecondaryPreferred constructs a read preference with a SecondaryPreferredMode.
func SecondaryPreferred(opts ...Option) *ReadPref ***REMOVED***
	// New only returns an error with a mode of Primary
	rp, _ := New(SecondaryPreferredMode, opts...)
	return rp
***REMOVED***

// Secondary constructs a read preference with a SecondaryMode.
func Secondary(opts ...Option) *ReadPref ***REMOVED***
	// New only returns an error with a mode of Primary
	rp, _ := New(SecondaryMode, opts...)
	return rp
***REMOVED***

// Nearest constructs a read preference with a NearestMode.
func Nearest(opts ...Option) *ReadPref ***REMOVED***
	// New only returns an error with a mode of Primary
	rp, _ := New(NearestMode, opts...)
	return rp
***REMOVED***

// New creates a new ReadPref.
func New(mode Mode, opts ...Option) (*ReadPref, error) ***REMOVED***
	rp := &ReadPref***REMOVED***
		mode: mode,
	***REMOVED***

	if mode == PrimaryMode && len(opts) != 0 ***REMOVED***
		return nil, errInvalidReadPreference
	***REMOVED***

	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		err := opt(rp)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return rp, nil
***REMOVED***

// ReadPref determines which servers are considered suitable for read operations.
type ReadPref struct ***REMOVED***
	maxStaleness    time.Duration
	maxStalenessSet bool
	mode            Mode
	tagSets         []tag.Set
	hedgeEnabled    *bool
***REMOVED***

// MaxStaleness is the maximum amount of time to allow
// a server to be considered eligible for selection. The
// second return value indicates if this value has been set.
func (r *ReadPref) MaxStaleness() (time.Duration, bool) ***REMOVED***
	return r.maxStaleness, r.maxStalenessSet
***REMOVED***

// Mode indicates the mode of the read preference.
func (r *ReadPref) Mode() Mode ***REMOVED***
	return r.mode
***REMOVED***

// TagSets are multiple tag sets indicating
// which servers should be considered.
func (r *ReadPref) TagSets() []tag.Set ***REMOVED***
	return r.tagSets
***REMOVED***

// HedgeEnabled returns whether or not hedged reads are enabled for this read preference. If this option was not
// specified during read preference construction, nil is returned.
func (r *ReadPref) HedgeEnabled() *bool ***REMOVED***
	return r.hedgeEnabled
***REMOVED***

// String returns a human-readable description of the read preference.
func (r *ReadPref) String() string ***REMOVED***
	var b bytes.Buffer
	b.WriteString(r.mode.String())
	delim := "("
	if r.maxStalenessSet ***REMOVED***
		fmt.Fprintf(&b, "%smaxStaleness=%v", delim, r.maxStaleness)
		delim = " "
	***REMOVED***
	for _, tagSet := range r.tagSets ***REMOVED***
		fmt.Fprintf(&b, "%stagSet=%s", delim, tagSet.String())
		delim = " "
	***REMOVED***
	if r.hedgeEnabled != nil ***REMOVED***
		fmt.Fprintf(&b, "%shedgeEnabled=%v", delim, *r.hedgeEnabled)
		delim = " "
	***REMOVED***
	if delim != "(" ***REMOVED***
		b.WriteString(")")
	***REMOVED***
	return b.String()
***REMOVED***
