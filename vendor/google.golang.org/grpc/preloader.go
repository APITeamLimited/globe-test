/*
 *
 * Copyright 2019 gRPC authors.
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

package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PreparedMsg is responsible for creating a Marshalled and Compressed object.
//
// This API is EXPERIMENTAL.
type PreparedMsg struct ***REMOVED***
	// Struct for preparing msg before sending them
	encodedData []byte
	hdr         []byte
	payload     []byte
***REMOVED***

// Encode marshalls and compresses the message using the codec and compressor for the stream.
func (p *PreparedMsg) Encode(s Stream, msg interface***REMOVED******REMOVED***) error ***REMOVED***
	ctx := s.Context()
	rpcInfo, ok := rpcInfoFromContext(ctx)
	if !ok ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: unable to get rpcInfo")
	***REMOVED***

	// check if the context has the relevant information to prepareMsg
	if rpcInfo.preloaderInfo == nil ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: rpcInfo.preloaderInfo is nil")
	***REMOVED***
	if rpcInfo.preloaderInfo.codec == nil ***REMOVED***
		return status.Errorf(codes.Internal, "grpc: rpcInfo.preloaderInfo.codec is nil")
	***REMOVED***

	// prepare the msg
	data, err := encode(rpcInfo.preloaderInfo.codec, msg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.encodedData = data
	compData, err := compress(data, rpcInfo.preloaderInfo.cp, rpcInfo.preloaderInfo.comp)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.hdr, p.payload = msgHeader(data, compData)
	return nil
***REMOVED***
