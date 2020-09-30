// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"sync"
)

var (
	nodeMu sync.Mutex
	ifname string  // name of interface being used
	nodeID [6]byte // hardware for version 1 UUIDs
	zeroID [6]byte // nodeID with only 0's
)

// NodeInterface returns the name of the interface from which the NodeID was
// derived.  The interface "user" is returned if the NodeID was set by
// SetNodeID.
func NodeInterface() string ***REMOVED***
	defer nodeMu.Unlock()
	nodeMu.Lock()
	return ifname
***REMOVED***

// SetNodeInterface selects the hardware address to be used for Version 1 UUIDs.
// If name is "" then the first usable interface found will be used or a random
// Node ID will be generated.  If a named interface cannot be found then false
// is returned.
//
// SetNodeInterface never fails when name is "".
func SetNodeInterface(name string) bool ***REMOVED***
	defer nodeMu.Unlock()
	nodeMu.Lock()
	return setNodeInterface(name)
***REMOVED***

func setNodeInterface(name string) bool ***REMOVED***
	iname, addr := getHardwareInterface(name) // null implementation for js
	if iname != "" && addr != nil ***REMOVED***
		ifname = iname
		copy(nodeID[:], addr)
		return true
	***REMOVED***

	// We found no interfaces with a valid hardware address.  If name
	// does not specify a specific interface generate a random Node ID
	// (section 4.1.6)
	if name == "" ***REMOVED***
		ifname = "random"
		randomBits(nodeID[:])
		return true
	***REMOVED***
	return false
***REMOVED***

// NodeID returns a slice of a copy of the current Node ID, setting the Node ID
// if not already set.
func NodeID() []byte ***REMOVED***
	defer nodeMu.Unlock()
	nodeMu.Lock()
	if nodeID == zeroID ***REMOVED***
		setNodeInterface("")
	***REMOVED***
	nid := nodeID
	return nid[:]
***REMOVED***

// SetNodeID sets the Node ID to be used for Version 1 UUIDs.  The first 6 bytes
// of id are used.  If id is less than 6 bytes then false is returned and the
// Node ID is not set.
func SetNodeID(id []byte) bool ***REMOVED***
	if len(id) < 6 ***REMOVED***
		return false
	***REMOVED***
	defer nodeMu.Unlock()
	nodeMu.Lock()
	copy(nodeID[:], id)
	ifname = "user"
	return true
***REMOVED***

// NodeID returns the 6 byte node id encoded in uuid.  It returns nil if uuid is
// not valid.  The NodeID is only well defined for version 1 and 2 UUIDs.
func (uuid UUID) NodeID() []byte ***REMOVED***
	var node [6]byte
	copy(node[:], uuid[10:])
	return node[:]
***REMOVED***
