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
	"fmt"

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

func (codec) Marshal(v interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	vv, ok := v.(proto.Message)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	***REMOVED***
	return proto.Marshal(vv)
***REMOVED***

func (codec) Unmarshal(data []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	vv, ok := v.(proto.Message)
	if !ok ***REMOVED***
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	***REMOVED***
	return proto.Unmarshal(data, vv)
***REMOVED***

func (codec) Name() string ***REMOVED***
	return Name
***REMOVED***
