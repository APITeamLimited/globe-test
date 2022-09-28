/*
 *
 * Copyright 2020 gRPC authors.
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

// Package metadata contains functions to set and get metadata from addresses.
//
// This package is experimental.
package metadata

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)

type mdKeyType string

const mdKey = mdKeyType("grpc.internal.address.metadata")

type mdValue metadata.MD

func (m mdValue) Equal(o interface***REMOVED******REMOVED***) bool ***REMOVED***
	om, ok := o.(mdValue)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	if len(m) != len(om) ***REMOVED***
		return false
	***REMOVED***
	for k, v := range m ***REMOVED***
		ov := om[k]
		if len(ov) != len(v) ***REMOVED***
			return false
		***REMOVED***
		for i, ve := range v ***REMOVED***
			if ov[i] != ve ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Get returns the metadata of addr.
func Get(addr resolver.Address) metadata.MD ***REMOVED***
	attrs := addr.Attributes
	if attrs == nil ***REMOVED***
		return nil
	***REMOVED***
	md, _ := attrs.Value(mdKey).(mdValue)
	return metadata.MD(md)
***REMOVED***

// Set sets (overrides) the metadata in addr.
//
// When a SubConn is created with this address, the RPCs sent on it will all
// have this metadata.
func Set(addr resolver.Address, md metadata.MD) resolver.Address ***REMOVED***
	addr.Attributes = addr.Attributes.WithValue(mdKey, mdValue(md))
	return addr
***REMOVED***

// Validate returns an error if the input md contains invalid keys or values.
//
// If the header is not a pseudo-header, the following items are checked:
// - header names must contain one or more characters from this set [0-9 a-z _ - .].
// - if the header-name ends with a "-bin" suffix, no validation of the header value is performed.
// - otherwise, the header value must contain one or more characters from the set [%x20-%x7E].
func Validate(md metadata.MD) error ***REMOVED***
	for k, vals := range md ***REMOVED***
		// pseudo-header will be ignored
		if k[0] == ':' ***REMOVED***
			continue
		***REMOVED***
		// check key, for i that saving a conversion if not using for range
		for i := 0; i < len(k); i++ ***REMOVED***
			r := k[i]
			if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '.' && r != '-' && r != '_' ***REMOVED***
				return fmt.Errorf("header key %q contains illegal characters not in [0-9a-z-_.]", k)
			***REMOVED***
		***REMOVED***
		if strings.HasSuffix(k, "-bin") ***REMOVED***
			continue
		***REMOVED***
		// check value
		for _, val := range vals ***REMOVED***
			if hasNotPrintable(val) ***REMOVED***
				return fmt.Errorf("header key %q contains value with non-printable ASCII characters", k)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// hasNotPrintable return true if msg contains any characters which are not in %x20-%x7E
func hasNotPrintable(msg string) bool ***REMOVED***
	// for i that saving a conversion if not using for range
	for i := 0; i < len(msg); i++ ***REMOVED***
		if msg[i] < 0x20 || msg[i] > 0x7E ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
