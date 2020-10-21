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

// Package attributes defines a generic key/value store used in various gRPC
// components.
//
// All APIs in this package are EXPERIMENTAL.
package attributes

import "fmt"

// Attributes is an immutable struct for storing and retrieving generic
// key/value pairs.  Keys must be hashable, and users should define their own
// types for keys.
type Attributes struct ***REMOVED***
	m map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
***REMOVED***

// New returns a new Attributes containing all key/value pairs in kvs.  If the
// same key appears multiple times, the last value overwrites all previous
// values for that key.  Panics if len(kvs) is not even.
func New(kvs ...interface***REMOVED******REMOVED***) *Attributes ***REMOVED***
	if len(kvs)%2 != 0 ***REMOVED***
		panic(fmt.Sprintf("attributes.New called with unexpected input: len(kvs) = %v", len(kvs)))
	***REMOVED***
	a := &Attributes***REMOVED***m: make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, len(kvs)/2)***REMOVED***
	for i := 0; i < len(kvs)/2; i++ ***REMOVED***
		a.m[kvs[i*2]] = kvs[i*2+1]
	***REMOVED***
	return a
***REMOVED***

// WithValues returns a new Attributes containing all key/value pairs in a and
// kvs.  Panics if len(kvs) is not even.  If the same key appears multiple
// times, the last value overwrites all previous values for that key.  To
// remove an existing key, use a nil value.
func (a *Attributes) WithValues(kvs ...interface***REMOVED******REMOVED***) *Attributes ***REMOVED***
	if a == nil ***REMOVED***
		return New(kvs...)
	***REMOVED***
	if len(kvs)%2 != 0 ***REMOVED***
		panic(fmt.Sprintf("attributes.New called with unexpected input: len(kvs) = %v", len(kvs)))
	***REMOVED***
	n := &Attributes***REMOVED***m: make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, len(a.m)+len(kvs)/2)***REMOVED***
	for k, v := range a.m ***REMOVED***
		n.m[k] = v
	***REMOVED***
	for i := 0; i < len(kvs)/2; i++ ***REMOVED***
		n.m[kvs[i*2]] = kvs[i*2+1]
	***REMOVED***
	return n
***REMOVED***

// Value returns the value associated with these attributes for key, or nil if
// no value is associated with key.
func (a *Attributes) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if a == nil ***REMOVED***
		return nil
	***REMOVED***
	return a.m[key]
***REMOVED***
