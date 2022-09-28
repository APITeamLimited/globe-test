// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonoptions

// SliceCodecOptions represents all possible options for slice encoding and decoding.
type SliceCodecOptions struct ***REMOVED***
	EncodeNilAsEmpty *bool // Specifies if a nil slice should encode as an empty array instead of null. Defaults to false.
***REMOVED***

// SliceCodec creates a new *SliceCodecOptions
func SliceCodec() *SliceCodecOptions ***REMOVED***
	return &SliceCodecOptions***REMOVED******REMOVED***
***REMOVED***

// SetEncodeNilAsEmpty specifies  if a nil slice should encode as an empty array instead of null. Defaults to false.
func (s *SliceCodecOptions) SetEncodeNilAsEmpty(b bool) *SliceCodecOptions ***REMOVED***
	s.EncodeNilAsEmpty = &b
	return s
***REMOVED***

// MergeSliceCodecOptions combines the given *SliceCodecOptions into a single *SliceCodecOptions in a last one wins fashion.
func MergeSliceCodecOptions(opts ...*SliceCodecOptions) *SliceCodecOptions ***REMOVED***
	s := SliceCodec()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.EncodeNilAsEmpty != nil ***REMOVED***
			s.EncodeNilAsEmpty = opt.EncodeNilAsEmpty
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***
