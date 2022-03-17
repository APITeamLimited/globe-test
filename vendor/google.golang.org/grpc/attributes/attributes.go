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
// Experimental
//
// Notice: This package is EXPERIMENTAL and may be changed or removed in a
// later release.
package attributes

// Attributes is an immutable struct for storing and retrieving generic
// key/value pairs.  Keys must be hashable, and users should define their own
// types for keys.  Values should not be modified after they are added to an
// Attributes or if they were received from one.  If values implement 'Equal(o
// interface***REMOVED******REMOVED***) bool', it will be called by (*Attributes).Equal to determine
// whether two values with the same key should be considered equal.
type Attributes struct ***REMOVED***
	m map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
***REMOVED***

// New returns a new Attributes containing the key/value pair.
func New(key, value interface***REMOVED******REMOVED***) *Attributes ***REMOVED***
	return &Attributes***REMOVED***m: map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED******REMOVED***key: value***REMOVED******REMOVED***
***REMOVED***

// WithValue returns a new Attributes containing the previous keys and values
// and the new key/value pair.  If the same key appears multiple times, the
// last value overwrites all previous values for that key.  To remove an
// existing key, use a nil value.  value should not be modified later.
func (a *Attributes) WithValue(key, value interface***REMOVED******REMOVED***) *Attributes ***REMOVED***
	if a == nil ***REMOVED***
		return New(key, value)
	***REMOVED***
	n := &Attributes***REMOVED***m: make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, len(a.m)+1)***REMOVED***
	for k, v := range a.m ***REMOVED***
		n.m[k] = v
	***REMOVED***
	n.m[key] = value
	return n
***REMOVED***

// Value returns the value associated with these attributes for key, or nil if
// no value is associated with key.  The returned value should not be modified.
func (a *Attributes) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if a == nil ***REMOVED***
		return nil
	***REMOVED***
	return a.m[key]
***REMOVED***

// Equal returns whether a and o are equivalent.  If 'Equal(o interface***REMOVED******REMOVED***)
// bool' is implemented for a value in the attributes, it is called to
// determine if the value matches the one stored in the other attributes.  If
// Equal is not implemented, standard equality is used to determine if the two
// values are equal. Note that some types (e.g. maps) aren't comparable by
// default, so they must be wrapped in a struct, or in an alias type, with Equal
// defined.
func (a *Attributes) Equal(o *Attributes) bool ***REMOVED***
	if a == nil && o == nil ***REMOVED***
		return true
	***REMOVED***
	if a == nil || o == nil ***REMOVED***
		return false
	***REMOVED***
	if len(a.m) != len(o.m) ***REMOVED***
		return false
	***REMOVED***
	for k, v := range a.m ***REMOVED***
		ov, ok := o.m[k]
		if !ok ***REMOVED***
			// o missing element of a
			return false
		***REMOVED***
		if eq, ok := v.(interface***REMOVED*** Equal(o interface***REMOVED******REMOVED***) bool ***REMOVED***); ok ***REMOVED***
			if !eq.Equal(ov) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else if v != ov ***REMOVED***
			// Fallback to a standard equality check if Value is unimplemented.
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
