// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonoptions

var defaultOverwriteDuplicatedInlinedFields = true

// StructCodecOptions represents all possible options for struct encoding and decoding.
type StructCodecOptions struct ***REMOVED***
	DecodeZeroStruct                 *bool // Specifies if structs should be zeroed before decoding into them. Defaults to false.
	DecodeDeepZeroInline             *bool // Specifies if structs should be recursively zeroed when a inline value is decoded. Defaults to false.
	EncodeOmitDefaultStruct          *bool // Specifies if default structs should be considered empty by omitempty. Defaults to false.
	AllowUnexportedFields            *bool // Specifies if unexported fields should be marshaled/unmarshaled. Defaults to false.
	OverwriteDuplicatedInlinedFields *bool // Specifies if fields in inlined structs can be overwritten by higher level struct fields with the same key. Defaults to true.
***REMOVED***

// StructCodec creates a new *StructCodecOptions
func StructCodec() *StructCodecOptions ***REMOVED***
	return &StructCodecOptions***REMOVED******REMOVED***
***REMOVED***

// SetDecodeZeroStruct specifies if structs should be zeroed before decoding into them. Defaults to false.
func (t *StructCodecOptions) SetDecodeZeroStruct(b bool) *StructCodecOptions ***REMOVED***
	t.DecodeZeroStruct = &b
	return t
***REMOVED***

// SetDecodeDeepZeroInline specifies if structs should be zeroed before decoding into them. Defaults to false.
func (t *StructCodecOptions) SetDecodeDeepZeroInline(b bool) *StructCodecOptions ***REMOVED***
	t.DecodeDeepZeroInline = &b
	return t
***REMOVED***

// SetEncodeOmitDefaultStruct specifies if default structs should be considered empty by omitempty. A default struct has all
// its values set to their default value. Defaults to false.
func (t *StructCodecOptions) SetEncodeOmitDefaultStruct(b bool) *StructCodecOptions ***REMOVED***
	t.EncodeOmitDefaultStruct = &b
	return t
***REMOVED***

// SetOverwriteDuplicatedInlinedFields specifies if inlined struct fields can be overwritten by higher level struct fields with the
// same bson key. When true and decoding, values will be written to the outermost struct with a matching key, and when
// encoding, keys will have the value of the top-most matching field. When false, decoding and encoding will error if
// there are duplicate keys after the struct is inlined. Defaults to true.
func (t *StructCodecOptions) SetOverwriteDuplicatedInlinedFields(b bool) *StructCodecOptions ***REMOVED***
	t.OverwriteDuplicatedInlinedFields = &b
	return t
***REMOVED***

// SetAllowUnexportedFields specifies if unexported fields should be marshaled/unmarshaled. Defaults to false.
func (t *StructCodecOptions) SetAllowUnexportedFields(b bool) *StructCodecOptions ***REMOVED***
	t.AllowUnexportedFields = &b
	return t
***REMOVED***

// MergeStructCodecOptions combines the given *StructCodecOptions into a single *StructCodecOptions in a last one wins fashion.
func MergeStructCodecOptions(opts ...*StructCodecOptions) *StructCodecOptions ***REMOVED***
	s := &StructCodecOptions***REMOVED***
		OverwriteDuplicatedInlinedFields: &defaultOverwriteDuplicatedInlinedFields,
	***REMOVED***
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***

		if opt.DecodeZeroStruct != nil ***REMOVED***
			s.DecodeZeroStruct = opt.DecodeZeroStruct
		***REMOVED***
		if opt.DecodeDeepZeroInline != nil ***REMOVED***
			s.DecodeDeepZeroInline = opt.DecodeDeepZeroInline
		***REMOVED***
		if opt.EncodeOmitDefaultStruct != nil ***REMOVED***
			s.EncodeOmitDefaultStruct = opt.EncodeOmitDefaultStruct
		***REMOVED***
		if opt.OverwriteDuplicatedInlinedFields != nil ***REMOVED***
			s.OverwriteDuplicatedInlinedFields = opt.OverwriteDuplicatedInlinedFields
		***REMOVED***
		if opt.AllowUnexportedFields != nil ***REMOVED***
			s.AllowUnexportedFields = opt.AllowUnexportedFields
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***
