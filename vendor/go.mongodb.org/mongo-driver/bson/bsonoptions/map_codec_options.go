// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonoptions

// MapCodecOptions represents all possible options for map encoding and decoding.
type MapCodecOptions struct ***REMOVED***
	DecodeZerosMap   *bool // Specifies if the map should be zeroed before decoding into it. Defaults to false.
	EncodeNilAsEmpty *bool // Specifies if a nil map should encode as an empty document instead of null. Defaults to false.
	// Specifies how keys should be handled. If false, the behavior matches encoding/json, where the encoding key type must
	// either be a string, an integer type, or implement bsoncodec.KeyMarshaler and the decoding key type must either be a
	// string, an integer type, or implement bsoncodec.KeyUnmarshaler. If true, keys are encoded with fmt.Sprint() and the
	// encoding key type must be a string, an integer type, or a float. If true, the use of Stringer will override
	// TextMarshaler/TextUnmarshaler. Defaults to false.
	EncodeKeysWithStringer *bool
***REMOVED***

// MapCodec creates a new *MapCodecOptions
func MapCodec() *MapCodecOptions ***REMOVED***
	return &MapCodecOptions***REMOVED******REMOVED***
***REMOVED***

// SetDecodeZerosMap specifies if the map should be zeroed before decoding into it. Defaults to false.
func (t *MapCodecOptions) SetDecodeZerosMap(b bool) *MapCodecOptions ***REMOVED***
	t.DecodeZerosMap = &b
	return t
***REMOVED***

// SetEncodeNilAsEmpty specifies if a nil map should encode as an empty document instead of null. Defaults to false.
func (t *MapCodecOptions) SetEncodeNilAsEmpty(b bool) *MapCodecOptions ***REMOVED***
	t.EncodeNilAsEmpty = &b
	return t
***REMOVED***

// SetEncodeKeysWithStringer specifies how keys should be handled. If false, the behavior matches encoding/json, where the
// encoding key type must either be a string, an integer type, or implement bsoncodec.KeyMarshaler and the decoding key
// type must either be a string, an integer type, or implement bsoncodec.KeyUnmarshaler. If true, keys are encoded with
// fmt.Sprint() and the encoding key type must be a string, an integer type, or a float. If true, the use of Stringer
// will override TextMarshaler/TextUnmarshaler. Defaults to false.
func (t *MapCodecOptions) SetEncodeKeysWithStringer(b bool) *MapCodecOptions ***REMOVED***
	t.EncodeKeysWithStringer = &b
	return t
***REMOVED***

// MergeMapCodecOptions combines the given *MapCodecOptions into a single *MapCodecOptions in a last one wins fashion.
func MergeMapCodecOptions(opts ...*MapCodecOptions) *MapCodecOptions ***REMOVED***
	s := MapCodec()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.DecodeZerosMap != nil ***REMOVED***
			s.DecodeZerosMap = opt.DecodeZerosMap
		***REMOVED***
		if opt.EncodeNilAsEmpty != nil ***REMOVED***
			s.EncodeNilAsEmpty = opt.EncodeNilAsEmpty
		***REMOVED***
		if opt.EncodeKeysWithStringer != nil ***REMOVED***
			s.EncodeKeysWithStringer = opt.EncodeKeysWithStringer
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***
