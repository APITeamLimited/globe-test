// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Topology contains information about a MongoDB cluster.
type Topology struct ***REMOVED***
	Servers               []Server
	SetName               string
	Kind                  TopologyKind
	SessionTimeoutMinutes uint32
	CompatibilityErr      error
***REMOVED***

// String implements the Stringer interface.
func (t Topology) String() string ***REMOVED***
	var serversStr string
	for _, s := range t.Servers ***REMOVED***
		serversStr += "***REMOVED*** " + s.String() + " ***REMOVED***, "
	***REMOVED***
	return fmt.Sprintf("Type: %s, Servers: [%s]", t.Kind, serversStr)
***REMOVED***

// Equal compares two topology descriptions and returns true if they are equal.
func (t Topology) Equal(other Topology) bool ***REMOVED***
	if t.Kind != other.Kind ***REMOVED***
		return false
	***REMOVED***

	topoServers := make(map[string]Server)
	for _, s := range t.Servers ***REMOVED***
		topoServers[s.Addr.String()] = s
	***REMOVED***

	otherServers := make(map[string]Server)
	for _, s := range other.Servers ***REMOVED***
		otherServers[s.Addr.String()] = s
	***REMOVED***

	if len(topoServers) != len(otherServers) ***REMOVED***
		return false
	***REMOVED***

	for _, server := range topoServers ***REMOVED***
		otherServer := otherServers[server.Addr.String()]

		if !server.Equal(otherServer) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

// HasReadableServer returns true if the topology contains a server suitable for reading.
//
// If the Topology's kind is Single or Sharded, the mode parameter is ignored and the function contains true if any of
// the servers in the Topology are of a known type.
//
// For replica sets, the function returns true if the cluster contains a server that matches the provided read
// preference mode.
func (t Topology) HasReadableServer(mode readpref.Mode) bool ***REMOVED***
	switch t.Kind ***REMOVED***
	case Single, Sharded:
		return hasAvailableServer(t.Servers, 0)
	case ReplicaSetWithPrimary:
		return hasAvailableServer(t.Servers, mode)
	case ReplicaSetNoPrimary, ReplicaSet:
		if mode == readpref.PrimaryMode ***REMOVED***
			return false
		***REMOVED***
		// invalid read preference
		if !mode.IsValid() ***REMOVED***
			return false
		***REMOVED***

		return hasAvailableServer(t.Servers, mode)
	***REMOVED***
	return false
***REMOVED***

// HasWritableServer returns true if a topology has a server available for writing.
//
// If the Topology's kind is Single or Sharded, this function returns true if any of the servers in the Topology are of
// a known type.
//
// For replica sets, the function returns true if the replica set contains a primary.
func (t Topology) HasWritableServer() bool ***REMOVED***
	return t.HasReadableServer(readpref.PrimaryMode)
***REMOVED***

// hasAvailableServer returns true if any servers are available based on the read preference.
func hasAvailableServer(servers []Server, mode readpref.Mode) bool ***REMOVED***
	switch mode ***REMOVED***
	case readpref.PrimaryMode:
		for _, s := range servers ***REMOVED***
			if s.Kind == RSPrimary ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	case readpref.PrimaryPreferredMode, readpref.SecondaryPreferredMode, readpref.NearestMode:
		for _, s := range servers ***REMOVED***
			if s.Kind == RSPrimary || s.Kind == RSSecondary ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	case readpref.SecondaryMode:
		for _, s := range servers ***REMOVED***
			if s.Kind == RSSecondary ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***

	// read preference is not specified
	for _, s := range servers ***REMOVED***
		switch s.Kind ***REMOVED***
		case Standalone,
			RSMember,
			RSPrimary,
			RSSecondary,
			RSArbiter,
			RSGhost,
			Mongos:
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***
