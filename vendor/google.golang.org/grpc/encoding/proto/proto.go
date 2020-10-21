/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"math"
	"sync"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

func init() ***REMOVED***
	encoding.RegisterCodec(codec***REMOVED******REMOVED***)
***REMOVED***

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type codec struct***REMOVED******REMOVED***

type cachedProtoBuffer struct ***REMOVED***
	lastMarshaledSize uint32
	proto.Buffer
***REMOVED***

func capToMaxInt32(val int) uint32 ***REMOVED***
	if val > math.MaxInt32 ***REMOVED***
		return uint32(math.MaxInt32)
	***REMOVED***
	return uint32(val)
***REMOVED***

func marshal(v interface***REMOVED******REMOVED***, cb *cachedProtoBuffer) ([]byte, error) ***REMOVED***
	protoMsg := v.(proto.Message)
	newSlice := make([]byte, 0, cb.lastMarshaledSize)

	cb.SetBuf(newSlice)
	cb.Reset()
	if err := cb.Marshal(protoMsg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	out := cb.Bytes()
	cb.lastMarshaledSize = capToMaxInt32(len(out))
	return out, nil
***REMOVED***

func (codec) Marshal(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if pm, ok := v.(proto.Marshaler); ok ***REMOVED***
		// object can marshal itself, no need for buffer
		return pm.Marshal()
	***REMOVED***

	cb := protoBufferPool.Get().(*cachedProtoBuffer)
	out, err := marshal(v, cb)

	// put back buffer and lose the ref to the slice
	cb.SetBuf(nil)
	protoBufferPool.Put(cb)
	return out, err
***REMOVED***

func (codec) Unmarshal(data []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	protoMsg := v.(proto.Message)
	protoMsg.Reset()

	if pu, ok := protoMsg.(proto.Unmarshaler); ok ***REMOVED***
		// object can unmarshal itself, no need for buffer
		return pu.Unmarshal(data)
	***REMOVED***

	cb := protoBufferPool.Get().(*cachedProtoBuffer)
	cb.SetBuf(data)
	err := cb.Unmarshal(protoMsg)
	cb.SetBuf(nil)
	protoBufferPool.Put(cb)
	return err
***REMOVED***

func (codec) Name() string ***REMOVED***
	return Name
***REMOVED***

var protoBufferPool = &sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return &cachedProtoBuffer***REMOVED***
			Buffer:            proto.Buffer***REMOVED******REMOVED***,
			lastMarshaledSize: 16,
		***REMOVED***
	***REMOVED***,
***REMOVED***
