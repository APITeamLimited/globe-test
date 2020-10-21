// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	protoV2 "google.golang.org/protobuf/proto"
)

var (
	// Deprecated: No longer returned.
	ErrNil = errors.New("proto: Marshal called with nil")

	// Deprecated: No longer returned.
	ErrTooLarge = errors.New("proto: message encodes to over 2 GB")

	// Deprecated: No longer returned.
	ErrInternalBadWireType = errors.New("proto: internal error: bad wiretype for oneof")
)

// Deprecated: Do not use.
type Stats struct***REMOVED*** Emalloc, Dmalloc, Encode, Decode, Chit, Cmiss, Size uint64 ***REMOVED***

// Deprecated: Do not use.
func GetStats() Stats ***REMOVED*** return Stats***REMOVED******REMOVED*** ***REMOVED***

// Deprecated: Do not use.
func MarshalMessageSet(interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return nil, errors.New("proto: not implemented")
***REMOVED***

// Deprecated: Do not use.
func UnmarshalMessageSet([]byte, interface***REMOVED******REMOVED***) error ***REMOVED***
	return errors.New("proto: not implemented")
***REMOVED***

// Deprecated: Do not use.
func MarshalMessageSetJSON(interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return nil, errors.New("proto: not implemented")
***REMOVED***

// Deprecated: Do not use.
func UnmarshalMessageSetJSON([]byte, interface***REMOVED******REMOVED***) error ***REMOVED***
	return errors.New("proto: not implemented")
***REMOVED***

// Deprecated: Do not use.
func RegisterMessageSetType(Message, int32, string) ***REMOVED******REMOVED***

// Deprecated: Do not use.
func EnumName(m map[int32]string, v int32) string ***REMOVED***
	s, ok := m[v]
	if ok ***REMOVED***
		return s
	***REMOVED***
	return strconv.Itoa(int(v))
***REMOVED***

// Deprecated: Do not use.
func UnmarshalJSONEnum(m map[string]int32, data []byte, enumName string) (int32, error) ***REMOVED***
	if data[0] == '"' ***REMOVED***
		// New style: enums are strings.
		var repr string
		if err := json.Unmarshal(data, &repr); err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		val, ok := m[repr]
		if !ok ***REMOVED***
			return 0, fmt.Errorf("unrecognized enum %s value %q", enumName, repr)
		***REMOVED***
		return val, nil
	***REMOVED***
	// Old style: enums are ints.
	var val int32
	if err := json.Unmarshal(data, &val); err != nil ***REMOVED***
		return 0, fmt.Errorf("cannot unmarshal %#q into enum %s", data, enumName)
	***REMOVED***
	return val, nil
***REMOVED***

// Deprecated: Do not use; this type existed for intenal-use only.
type InternalMessageInfo struct***REMOVED******REMOVED***

// Deprecated: Do not use; this method existed for intenal-use only.
func (*InternalMessageInfo) DiscardUnknown(m Message) ***REMOVED***
	DiscardUnknown(m)
***REMOVED***

// Deprecated: Do not use; this method existed for intenal-use only.
func (*InternalMessageInfo) Marshal(b []byte, m Message, deterministic bool) ([]byte, error) ***REMOVED***
	return protoV2.MarshalOptions***REMOVED***Deterministic: deterministic***REMOVED***.MarshalAppend(b, MessageV2(m))
***REMOVED***

// Deprecated: Do not use; this method existed for intenal-use only.
func (*InternalMessageInfo) Merge(dst, src Message) ***REMOVED***
	protoV2.Merge(MessageV2(dst), MessageV2(src))
***REMOVED***

// Deprecated: Do not use; this method existed for intenal-use only.
func (*InternalMessageInfo) Size(m Message) int ***REMOVED***
	return protoV2.Size(MessageV2(m))
***REMOVED***

// Deprecated: Do not use; this method existed for intenal-use only.
func (*InternalMessageInfo) Unmarshal(m Message, b []byte) error ***REMOVED***
	return protoV2.UnmarshalOptions***REMOVED***Merge: true***REMOVED***.Unmarshal(b, MessageV2(m))
***REMOVED***
