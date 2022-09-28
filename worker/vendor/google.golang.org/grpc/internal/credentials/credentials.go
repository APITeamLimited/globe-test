/*
 * Copyright 2021 gRPC authors.
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
 */

package credentials

import (
	"context"
)

// requestInfoKey is a struct to be used as the key to store RequestInfo in a
// context.
type requestInfoKey struct***REMOVED******REMOVED***

// NewRequestInfoContext creates a context with ri.
func NewRequestInfoContext(ctx context.Context, ri interface***REMOVED******REMOVED***) context.Context ***REMOVED***
	return context.WithValue(ctx, requestInfoKey***REMOVED******REMOVED***, ri)
***REMOVED***

// RequestInfoFromContext extracts the RequestInfo from ctx.
func RequestInfoFromContext(ctx context.Context) interface***REMOVED******REMOVED*** ***REMOVED***
	return ctx.Value(requestInfoKey***REMOVED******REMOVED***)
***REMOVED***

// clientHandshakeInfoKey is a struct used as the key to store
// ClientHandshakeInfo in a context.
type clientHandshakeInfoKey struct***REMOVED******REMOVED***

// ClientHandshakeInfoFromContext extracts the ClientHandshakeInfo from ctx.
func ClientHandshakeInfoFromContext(ctx context.Context) interface***REMOVED******REMOVED*** ***REMOVED***
	return ctx.Value(clientHandshakeInfoKey***REMOVED******REMOVED***)
***REMOVED***

// NewClientHandshakeInfoContext creates a context with chi.
func NewClientHandshakeInfoContext(ctx context.Context, chi interface***REMOVED******REMOVED***) context.Context ***REMOVED***
	return context.WithValue(ctx, clientHandshakeInfoKey***REMOVED******REMOVED***, chi)
***REMOVED***
