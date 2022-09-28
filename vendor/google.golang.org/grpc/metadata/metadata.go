/*
 *
 * Copyright 2014 gRPC authors.
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

// Package metadata define the structure of the metadata supported by gRPC library.
// Please refer to https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
// for more information about custom-metadata.
package metadata // import "google.golang.org/grpc/metadata"

import (
	"context"
	"fmt"
	"strings"
)

// DecodeKeyValue returns k, v, nil.
//
// Deprecated: use k and v directly instead.
func DecodeKeyValue(k, v string) (string, string, error) ***REMOVED***
	return k, v, nil
***REMOVED***

// MD is a mapping from metadata keys to values. Users should use the following
// two convenience functions New and Pairs to generate MD.
type MD map[string][]string

// New creates an MD from a given key-value map.
//
// Only the following ASCII characters are allowed in keys:
//  - digits: 0-9
//  - uppercase letters: A-Z (normalized to lower)
//  - lowercase letters: a-z
//  - special characters: -_.
// Uppercase letters are automatically converted to lowercase.
//
// Keys beginning with "grpc-" are reserved for grpc-internal use only and may
// result in errors if set in metadata.
func New(m map[string]string) MD ***REMOVED***
	md := MD***REMOVED******REMOVED***
	for k, val := range m ***REMOVED***
		key := strings.ToLower(k)
		md[key] = append(md[key], val)
	***REMOVED***
	return md
***REMOVED***

// Pairs returns an MD formed by the mapping of key, value ...
// Pairs panics if len(kv) is odd.
//
// Only the following ASCII characters are allowed in keys:
//  - digits: 0-9
//  - uppercase letters: A-Z (normalized to lower)
//  - lowercase letters: a-z
//  - special characters: -_.
// Uppercase letters are automatically converted to lowercase.
//
// Keys beginning with "grpc-" are reserved for grpc-internal use only and may
// result in errors if set in metadata.
func Pairs(kv ...string) MD ***REMOVED***
	if len(kv)%2 == 1 ***REMOVED***
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	***REMOVED***
	md := MD***REMOVED******REMOVED***
	for i := 0; i < len(kv); i += 2 ***REMOVED***
		key := strings.ToLower(kv[i])
		md[key] = append(md[key], kv[i+1])
	***REMOVED***
	return md
***REMOVED***

// Len returns the number of items in md.
func (md MD) Len() int ***REMOVED***
	return len(md)
***REMOVED***

// Copy returns a copy of md.
func (md MD) Copy() MD ***REMOVED***
	return Join(md)
***REMOVED***

// Get obtains the values for a given key.
//
// k is converted to lowercase before searching in md.
func (md MD) Get(k string) []string ***REMOVED***
	k = strings.ToLower(k)
	return md[k]
***REMOVED***

// Set sets the value of a given key with a slice of values.
//
// k is converted to lowercase before storing in md.
func (md MD) Set(k string, vals ...string) ***REMOVED***
	if len(vals) == 0 ***REMOVED***
		return
	***REMOVED***
	k = strings.ToLower(k)
	md[k] = vals
***REMOVED***

// Append adds the values to key k, not overwriting what was already stored at
// that key.
//
// k is converted to lowercase before storing in md.
func (md MD) Append(k string, vals ...string) ***REMOVED***
	if len(vals) == 0 ***REMOVED***
		return
	***REMOVED***
	k = strings.ToLower(k)
	md[k] = append(md[k], vals...)
***REMOVED***

// Delete removes the values for a given key k which is converted to lowercase
// before removing it from md.
func (md MD) Delete(k string) ***REMOVED***
	k = strings.ToLower(k)
	delete(md, k)
***REMOVED***

// Join joins any number of mds into a single MD.
//
// The order of values for each key is determined by the order in which the mds
// containing those values are presented to Join.
func Join(mds ...MD) MD ***REMOVED***
	out := MD***REMOVED******REMOVED***
	for _, md := range mds ***REMOVED***
		for k, v := range md ***REMOVED***
			out[k] = append(out[k], v...)
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***

type mdIncomingKey struct***REMOVED******REMOVED***
type mdOutgoingKey struct***REMOVED******REMOVED***

// NewIncomingContext creates a new context with incoming md attached.
func NewIncomingContext(ctx context.Context, md MD) context.Context ***REMOVED***
	return context.WithValue(ctx, mdIncomingKey***REMOVED******REMOVED***, md)
***REMOVED***

// NewOutgoingContext creates a new context with outgoing md attached. If used
// in conjunction with AppendToOutgoingContext, NewOutgoingContext will
// overwrite any previously-appended metadata.
func NewOutgoingContext(ctx context.Context, md MD) context.Context ***REMOVED***
	return context.WithValue(ctx, mdOutgoingKey***REMOVED******REMOVED***, rawMD***REMOVED***md: md***REMOVED***)
***REMOVED***

// AppendToOutgoingContext returns a new context with the provided kv merged
// with any existing metadata in the context. Please refer to the documentation
// of Pairs for a description of kv.
func AppendToOutgoingContext(ctx context.Context, kv ...string) context.Context ***REMOVED***
	if len(kv)%2 == 1 ***REMOVED***
		panic(fmt.Sprintf("metadata: AppendToOutgoingContext got an odd number of input pairs for metadata: %d", len(kv)))
	***REMOVED***
	md, _ := ctx.Value(mdOutgoingKey***REMOVED******REMOVED***).(rawMD)
	added := make([][]string, len(md.added)+1)
	copy(added, md.added)
	added[len(added)-1] = make([]string, len(kv))
	copy(added[len(added)-1], kv)
	return context.WithValue(ctx, mdOutgoingKey***REMOVED******REMOVED***, rawMD***REMOVED***md: md.md, added: added***REMOVED***)
***REMOVED***

// FromIncomingContext returns the incoming metadata in ctx if it exists.
//
// All keys in the returned MD are lowercase.
func FromIncomingContext(ctx context.Context) (MD, bool) ***REMOVED***
	md, ok := ctx.Value(mdIncomingKey***REMOVED******REMOVED***).(MD)
	if !ok ***REMOVED***
		return nil, false
	***REMOVED***
	out := MD***REMOVED******REMOVED***
	for k, v := range md ***REMOVED***
		// We need to manually convert all keys to lower case, because MD is a
		// map, and there's no guarantee that the MD attached to the context is
		// created using our helper functions.
		key := strings.ToLower(k)
		s := make([]string, len(v))
		copy(s, v)
		out[key] = s
	***REMOVED***
	return out, true
***REMOVED***

// FromOutgoingContextRaw returns the un-merged, intermediary contents of rawMD.
//
// Remember to perform strings.ToLower on the keys, for both the returned MD (MD
// is a map, there's no guarantee it's created using our helper functions) and
// the extra kv pairs (AppendToOutgoingContext doesn't turn them into
// lowercase).
//
// This is intended for gRPC-internal use ONLY. Users should use
// FromOutgoingContext instead.
func FromOutgoingContextRaw(ctx context.Context) (MD, [][]string, bool) ***REMOVED***
	raw, ok := ctx.Value(mdOutgoingKey***REMOVED******REMOVED***).(rawMD)
	if !ok ***REMOVED***
		return nil, nil, false
	***REMOVED***

	return raw.md, raw.added, true
***REMOVED***

// FromOutgoingContext returns the outgoing metadata in ctx if it exists.
//
// All keys in the returned MD are lowercase.
func FromOutgoingContext(ctx context.Context) (MD, bool) ***REMOVED***
	raw, ok := ctx.Value(mdOutgoingKey***REMOVED******REMOVED***).(rawMD)
	if !ok ***REMOVED***
		return nil, false
	***REMOVED***

	out := MD***REMOVED******REMOVED***
	for k, v := range raw.md ***REMOVED***
		// We need to manually convert all keys to lower case, because MD is a
		// map, and there's no guarantee that the MD attached to the context is
		// created using our helper functions.
		key := strings.ToLower(k)
		s := make([]string, len(v))
		copy(s, v)
		out[key] = s
	***REMOVED***
	for _, added := range raw.added ***REMOVED***
		if len(added)%2 == 1 ***REMOVED***
			panic(fmt.Sprintf("metadata: FromOutgoingContext got an odd number of input pairs for metadata: %d", len(added)))
		***REMOVED***

		for i := 0; i < len(added); i += 2 ***REMOVED***
			key := strings.ToLower(added[i])
			out[key] = append(out[key], added[i+1])
		***REMOVED***
	***REMOVED***
	return out, ok
***REMOVED***

type rawMD struct ***REMOVED***
	md    MD
	added [][]string
***REMOVED***
