// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonoptions

// EmptyInterfaceCodecOptions represents all possible options for interface***REMOVED******REMOVED*** encoding and decoding.
type EmptyInterfaceCodecOptions struct ***REMOVED***
	DecodeBinaryAsSlice *bool // Specifies if Old and Generic type binarys should default to []slice instead of primitive.Binary. Defaults to false.
***REMOVED***

// EmptyInterfaceCodec creates a new *EmptyInterfaceCodecOptions
func EmptyInterfaceCodec() *EmptyInterfaceCodecOptions ***REMOVED***
	return &EmptyInterfaceCodecOptions***REMOVED******REMOVED***
***REMOVED***

// SetDecodeBinaryAsSlice specifies if Old and Generic type binarys should default to []slice instead of primitive.Binary. Defaults to false.
func (e *EmptyInterfaceCodecOptions) SetDecodeBinaryAsSlice(b bool) *EmptyInterfaceCodecOptions ***REMOVED***
	e.DecodeBinaryAsSlice = &b
	return e
***REMOVED***

// MergeEmptyInterfaceCodecOptions combines the given *EmptyInterfaceCodecOptions into a single *EmptyInterfaceCodecOptions in a last one wins fashion.
func MergeEmptyInterfaceCodecOptions(opts ...*EmptyInterfaceCodecOptions) *EmptyInterfaceCodecOptions ***REMOVED***
	e := EmptyInterfaceCodec()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.DecodeBinaryAsSlice != nil ***REMOVED***
			e.DecodeBinaryAsSlice = opt.DecodeBinaryAsSlice
		***REMOVED***
	***REMOVED***

	return e
***REMOVED***
