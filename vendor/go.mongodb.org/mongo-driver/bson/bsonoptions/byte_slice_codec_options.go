// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonoptions

// ByteSliceCodecOptions represents all possible options for byte slice encoding and decoding.
type ByteSliceCodecOptions struct ***REMOVED***
	EncodeNilAsEmpty *bool // Specifies if a nil byte slice should encode as an empty binary instead of null. Defaults to false.
***REMOVED***

// ByteSliceCodec creates a new *ByteSliceCodecOptions
func ByteSliceCodec() *ByteSliceCodecOptions ***REMOVED***
	return &ByteSliceCodecOptions***REMOVED******REMOVED***
***REMOVED***

// SetEncodeNilAsEmpty specifies  if a nil byte slice should encode as an empty binary instead of null. Defaults to false.
func (bs *ByteSliceCodecOptions) SetEncodeNilAsEmpty(b bool) *ByteSliceCodecOptions ***REMOVED***
	bs.EncodeNilAsEmpty = &b
	return bs
***REMOVED***

// MergeByteSliceCodecOptions combines the given *ByteSliceCodecOptions into a single *ByteSliceCodecOptions in a last one wins fashion.
func MergeByteSliceCodecOptions(opts ...*ByteSliceCodecOptions) *ByteSliceCodecOptions ***REMOVED***
	bs := ByteSliceCodec()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.EncodeNilAsEmpty != nil ***REMOVED***
			bs.EncodeNilAsEmpty = opt.EncodeNilAsEmpty
		***REMOVED***
	***REMOVED***

	return bs
***REMOVED***
