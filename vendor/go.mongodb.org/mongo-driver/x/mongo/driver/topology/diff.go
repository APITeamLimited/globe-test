// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import "go.mongodb.org/mongo-driver/mongo/description"

// hostlistDiff is the difference between a topology and a host list.
type hostlistDiff struct ***REMOVED***
	Added   []string
	Removed []string
***REMOVED***

// diffHostList compares the topology description and host list and returns the difference.
func diffHostList(t description.Topology, hostlist []string) hostlistDiff ***REMOVED***
	var diff hostlistDiff

	oldServers := make(map[string]bool)
	for _, s := range t.Servers ***REMOVED***
		oldServers[s.Addr.String()] = true
	***REMOVED***

	for _, addr := range hostlist ***REMOVED***
		if oldServers[addr] ***REMOVED***
			delete(oldServers, addr)
		***REMOVED*** else ***REMOVED***
			diff.Added = append(diff.Added, addr)
		***REMOVED***
	***REMOVED***

	for addr := range oldServers ***REMOVED***
		diff.Removed = append(diff.Removed, addr)
	***REMOVED***

	return diff
***REMOVED***

// topologyDiff is the difference between two different topology descriptions.
type topologyDiff struct ***REMOVED***
	Added   []description.Server
	Removed []description.Server
***REMOVED***

// diffTopology compares the two topology descriptions and returns the difference.
func diffTopology(old, new description.Topology) topologyDiff ***REMOVED***
	var diff topologyDiff

	oldServers := make(map[string]bool)
	for _, s := range old.Servers ***REMOVED***
		oldServers[s.Addr.String()] = true
	***REMOVED***

	for _, s := range new.Servers ***REMOVED***
		addr := s.Addr.String()
		if oldServers[addr] ***REMOVED***
			delete(oldServers, addr)
		***REMOVED*** else ***REMOVED***
			diff.Added = append(diff.Added, s)
		***REMOVED***
	***REMOVED***

	for _, s := range old.Servers ***REMOVED***
		addr := s.Addr.String()
		if oldServers[addr] ***REMOVED***
			diff.Removed = append(diff.Removed, s)
		***REMOVED***
	***REMOVED***

	return diff
***REMOVED***
