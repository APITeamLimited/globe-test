// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonoptions

var defaultDecodeOIDAsHex = true

// StringCodecOptions represents all possible options for string encoding and decoding.
type StringCodecOptions struct ***REMOVED***
	DecodeObjectIDAsHex *bool // Specifies if we should decode ObjectID as the hex value. Defaults to true.
***REMOVED***

// StringCodec creates a new *StringCodecOptions
func StringCodec() *StringCodecOptions ***REMOVED***
	return &StringCodecOptions***REMOVED******REMOVED***
***REMOVED***

// SetDecodeObjectIDAsHex specifies if object IDs should be decoded as their hex representation. If false, a string made
// from the raw object ID bytes will be used. Defaults to true.
func (t *StringCodecOptions) SetDecodeObjectIDAsHex(b bool) *StringCodecOptions ***REMOVED***
	t.DecodeObjectIDAsHex = &b
	return t
***REMOVED***

// MergeStringCodecOptions combines the given *StringCodecOptions into a single *StringCodecOptions in a last one wins fashion.
func MergeStringCodecOptions(opts ...*StringCodecOptions) *StringCodecOptions ***REMOVED***
	s := &StringCodecOptions***REMOVED***&defaultDecodeOIDAsHex***REMOVED***
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.DecodeObjectIDAsHex != nil ***REMOVED***
			s.DecodeObjectIDAsHex = opt.DecodeObjectIDAsHex
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***
