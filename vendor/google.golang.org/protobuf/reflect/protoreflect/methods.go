// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protoreflect

import (
	"google.golang.org/protobuf/internal/pragma"
)

// The following types are used by the fast-path Message.ProtoMethods method.
//
// To avoid polluting the public protoreflect API with types used only by
// low-level implementations, the canonical definitions of these types are
// in the runtime/protoiface package. The definitions here and in protoiface
// must be kept in sync.
type (
	methods = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Flags            supportFlags
		Size             func(sizeInput) sizeOutput
		Marshal          func(marshalInput) (marshalOutput, error)
		Unmarshal        func(unmarshalInput) (unmarshalOutput, error)
		Merge            func(mergeInput) mergeOutput
		CheckInitialized func(checkInitializedInput) (checkInitializedOutput, error)
	***REMOVED***
	supportFlags = uint64
	sizeInput    = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Message Message
		Flags   uint8
	***REMOVED***
	sizeOutput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Size int
	***REMOVED***
	marshalInput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Message Message
		Buf     []byte
		Flags   uint8
	***REMOVED***
	marshalOutput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Buf []byte
	***REMOVED***
	unmarshalInput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Message  Message
		Buf      []byte
		Flags    uint8
		Resolver interface ***REMOVED***
			FindExtensionByName(field FullName) (ExtensionType, error)
			FindExtensionByNumber(message FullName, field FieldNumber) (ExtensionType, error)
		***REMOVED***
	***REMOVED***
	unmarshalOutput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Flags uint8
	***REMOVED***
	mergeInput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Source      Message
		Destination Message
	***REMOVED***
	mergeOutput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Flags uint8
	***REMOVED***
	checkInitializedInput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
		Message Message
	***REMOVED***
	checkInitializedOutput = struct ***REMOVED***
		pragma.NoUnkeyedLiterals
	***REMOVED***
)
