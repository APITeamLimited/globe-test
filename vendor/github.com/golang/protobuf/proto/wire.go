// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	protoV2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

// Size returns the size in bytes of the wire-format encoding of m.
func Size(m Message) int ***REMOVED***
	if m == nil ***REMOVED***
		return 0
	***REMOVED***
	mi := MessageV2(m)
	return protoV2.Size(mi)
***REMOVED***

// Marshal returns the wire-format encoding of m.
func Marshal(m Message) ([]byte, error) ***REMOVED***
	b, err := marshalAppend(nil, m, false)
	if b == nil ***REMOVED***
		b = zeroBytes
	***REMOVED***
	return b, err
***REMOVED***

var zeroBytes = make([]byte, 0, 0)

func marshalAppend(buf []byte, m Message, deterministic bool) ([]byte, error) ***REMOVED***
	if m == nil ***REMOVED***
		return nil, ErrNil
	***REMOVED***
	mi := MessageV2(m)
	nbuf, err := protoV2.MarshalOptions***REMOVED***
		Deterministic: deterministic,
		AllowPartial:  true,
	***REMOVED***.MarshalAppend(buf, mi)
	if err != nil ***REMOVED***
		return buf, err
	***REMOVED***
	if len(buf) == len(nbuf) ***REMOVED***
		if !mi.ProtoReflect().IsValid() ***REMOVED***
			return buf, ErrNil
		***REMOVED***
	***REMOVED***
	return nbuf, checkRequiredNotSet(mi)
***REMOVED***

// Unmarshal parses a wire-format message in b and places the decoded results in m.
//
// Unmarshal resets m before starting to unmarshal, so any existing data in m is always
// removed. Use UnmarshalMerge to preserve and append to existing data.
func Unmarshal(b []byte, m Message) error ***REMOVED***
	m.Reset()
	return UnmarshalMerge(b, m)
***REMOVED***

// UnmarshalMerge parses a wire-format message in b and places the decoded results in m.
func UnmarshalMerge(b []byte, m Message) error ***REMOVED***
	mi := MessageV2(m)
	out, err := protoV2.UnmarshalOptions***REMOVED***
		AllowPartial: true,
		Merge:        true,
	***REMOVED***.UnmarshalState(protoiface.UnmarshalInput***REMOVED***
		Buf:     b,
		Message: mi.ProtoReflect(),
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if out.Flags&protoiface.UnmarshalInitialized > 0 ***REMOVED***
		return nil
	***REMOVED***
	return checkRequiredNotSet(mi)
***REMOVED***
