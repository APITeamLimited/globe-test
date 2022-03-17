/*
 *
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
 *
 */

package resolver

type addressMapEntry struct ***REMOVED***
	addr  Address
	value interface***REMOVED******REMOVED***
***REMOVED***

// AddressMap is a map of addresses to arbitrary values taking into account
// Attributes.  BalancerAttributes are ignored, as are Metadata and Type.
// Multiple accesses may not be performed concurrently.  Must be created via
// NewAddressMap; do not construct directly.
type AddressMap struct ***REMOVED***
	m map[string]addressMapEntryList
***REMOVED***

type addressMapEntryList []*addressMapEntry

// NewAddressMap creates a new AddressMap.
func NewAddressMap() *AddressMap ***REMOVED***
	return &AddressMap***REMOVED***m: make(map[string]addressMapEntryList)***REMOVED***
***REMOVED***

// find returns the index of addr in the addressMapEntry slice, or -1 if not
// present.
func (l addressMapEntryList) find(addr Address) int ***REMOVED***
	if len(l) == 0 ***REMOVED***
		return -1
	***REMOVED***
	for i, entry := range l ***REMOVED***
		if entry.addr.ServerName == addr.ServerName &&
			entry.addr.Attributes.Equal(addr.Attributes) ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// Get returns the value for the address in the map, if present.
func (a *AddressMap) Get(addr Address) (value interface***REMOVED******REMOVED***, ok bool) ***REMOVED***
	entryList := a.m[addr.Addr]
	if entry := entryList.find(addr); entry != -1 ***REMOVED***
		return entryList[entry].value, true
	***REMOVED***
	return nil, false
***REMOVED***

// Set updates or adds the value to the address in the map.
func (a *AddressMap) Set(addr Address, value interface***REMOVED******REMOVED***) ***REMOVED***
	entryList := a.m[addr.Addr]
	if entry := entryList.find(addr); entry != -1 ***REMOVED***
		a.m[addr.Addr][entry].value = value
		return
	***REMOVED***
	a.m[addr.Addr] = append(a.m[addr.Addr], &addressMapEntry***REMOVED***addr: addr, value: value***REMOVED***)
***REMOVED***

// Delete removes addr from the map.
func (a *AddressMap) Delete(addr Address) ***REMOVED***
	entryList := a.m[addr.Addr]
	entry := entryList.find(addr)
	if entry == -1 ***REMOVED***
		return
	***REMOVED***
	if len(entryList) == 1 ***REMOVED***
		entryList = nil
	***REMOVED*** else ***REMOVED***
		copy(entryList[entry:], entryList[entry+1:])
		entryList = entryList[:len(entryList)-1]
	***REMOVED***
	a.m[addr.Addr] = entryList
***REMOVED***

// Len returns the number of entries in the map.
func (a *AddressMap) Len() int ***REMOVED***
	ret := 0
	for _, entryList := range a.m ***REMOVED***
		ret += len(entryList)
	***REMOVED***
	return ret
***REMOVED***

// Keys returns a slice of all current map keys.
func (a *AddressMap) Keys() []Address ***REMOVED***
	ret := make([]Address, 0, a.Len())
	for _, entryList := range a.m ***REMOVED***
		for _, entry := range entryList ***REMOVED***
			ret = append(ret, entry.addr)
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***
