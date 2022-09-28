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
	// The underlying map is keyed by an Address with fields that we don't care
	// about being set to their zero values. The only fields that we care about
	// are `Addr`, `ServerName` and `Attributes`. Since we need to be able to
	// distinguish between addresses with same `Addr` and `ServerName`, but
	// different `Attributes`, we cannot store the `Attributes` in the map key.
	//
	// The comparison operation for structs work as follows:
	//  Struct values are comparable if all their fields are comparable. Two
	//  struct values are equal if their corresponding non-blank fields are equal.
	//
	// The value type of the map contains a slice of addresses which match the key
	// in their `Addr` and `ServerName` fields and contain the corresponding value
	// associated with them.
	m map[Address]addressMapEntryList
***REMOVED***

func toMapKey(addr *Address) Address ***REMOVED***
	return Address***REMOVED***Addr: addr.Addr, ServerName: addr.ServerName***REMOVED***
***REMOVED***

type addressMapEntryList []*addressMapEntry

// NewAddressMap creates a new AddressMap.
func NewAddressMap() *AddressMap ***REMOVED***
	return &AddressMap***REMOVED***m: make(map[Address]addressMapEntryList)***REMOVED***
***REMOVED***

// find returns the index of addr in the addressMapEntry slice, or -1 if not
// present.
func (l addressMapEntryList) find(addr Address) int ***REMOVED***
	for i, entry := range l ***REMOVED***
		// Attributes are the only thing to match on here, since `Addr` and
		// `ServerName` are already equal.
		if entry.addr.Attributes.Equal(addr.Attributes) ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// Get returns the value for the address in the map, if present.
func (a *AddressMap) Get(addr Address) (value interface***REMOVED******REMOVED***, ok bool) ***REMOVED***
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	if entry := entryList.find(addr); entry != -1 ***REMOVED***
		return entryList[entry].value, true
	***REMOVED***
	return nil, false
***REMOVED***

// Set updates or adds the value to the address in the map.
func (a *AddressMap) Set(addr Address, value interface***REMOVED******REMOVED***) ***REMOVED***
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	if entry := entryList.find(addr); entry != -1 ***REMOVED***
		entryList[entry].value = value
		return
	***REMOVED***
	a.m[addrKey] = append(entryList, &addressMapEntry***REMOVED***addr: addr, value: value***REMOVED***)
***REMOVED***

// Delete removes addr from the map.
func (a *AddressMap) Delete(addr Address) ***REMOVED***
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
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
	a.m[addrKey] = entryList
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

// Values returns a slice of all current map values.
func (a *AddressMap) Values() []interface***REMOVED******REMOVED*** ***REMOVED***
	ret := make([]interface***REMOVED******REMOVED***, 0, a.Len())
	for _, entryList := range a.m ***REMOVED***
		for _, entry := range entryList ***REMOVED***
			ret = append(ret, entry.value)
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***
