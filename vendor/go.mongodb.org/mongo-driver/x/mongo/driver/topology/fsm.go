// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"bytes"
	"fmt"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
)

var (
	// SupportedWireVersions is the range of wire versions supported by the driver.
	SupportedWireVersions = description.NewVersionRange(2, 17)
)

const (
	// MinSupportedMongoDBVersion is the version string for the lowest MongoDB version supported by the driver.
	MinSupportedMongoDBVersion = "2.6"
)

type fsm struct ***REMOVED***
	description.Topology
	maxElectionID    primitive.ObjectID
	maxSetVersion    uint32
	compatible       atomic.Value
	compatibilityErr error
***REMOVED***

func newFSM() *fsm ***REMOVED***
	f := fsm***REMOVED******REMOVED***
	f.compatible.Store(true)
	return &f
***REMOVED***

// apply takes a new server description and modifies the FSM's topology description based on it. It returns the
// updated topology description as well as a server description. The returned server description is either the same
// one that was passed in, or a new one in the case that it had to be changed.
//
// apply should operation on immutable descriptions so we don't have to lock for the entire time we're applying the
// server description.
func (f *fsm) apply(s description.Server) (description.Topology, description.Server) ***REMOVED***
	newServers := make([]description.Server, len(f.Servers))
	copy(newServers, f.Servers)

	oldMinutes := f.SessionTimeoutMinutes
	f.Topology = description.Topology***REMOVED***
		Kind:    f.Kind,
		Servers: newServers,
		SetName: f.SetName,
	***REMOVED***

	// For data bearing servers, set SessionTimeoutMinutes to the lowest among them
	if oldMinutes == 0 ***REMOVED***
		// If timeout currently 0, check all servers to see if any still don't have a timeout
		// If they all have timeout, pick the lowest.
		timeout := s.SessionTimeoutMinutes
		for _, server := range f.Servers ***REMOVED***
			if server.DataBearing() && server.SessionTimeoutMinutes < timeout ***REMOVED***
				timeout = server.SessionTimeoutMinutes
			***REMOVED***
		***REMOVED***
		f.SessionTimeoutMinutes = timeout
	***REMOVED*** else ***REMOVED***
		if s.DataBearing() && oldMinutes > s.SessionTimeoutMinutes ***REMOVED***
			f.SessionTimeoutMinutes = s.SessionTimeoutMinutes
		***REMOVED*** else ***REMOVED***
			f.SessionTimeoutMinutes = oldMinutes
		***REMOVED***
	***REMOVED***

	if _, ok := f.findServer(s.Addr); !ok ***REMOVED***
		return f.Topology, s
	***REMOVED***

	updatedDesc := s
	switch f.Kind ***REMOVED***
	case description.Unknown:
		updatedDesc = f.applyToUnknown(s)
	case description.Sharded:
		updatedDesc = f.applyToSharded(s)
	case description.ReplicaSetNoPrimary:
		updatedDesc = f.applyToReplicaSetNoPrimary(s)
	case description.ReplicaSetWithPrimary:
		updatedDesc = f.applyToReplicaSetWithPrimary(s)
	case description.Single:
		updatedDesc = f.applyToSingle(s)
	***REMOVED***

	for _, server := range f.Servers ***REMOVED***
		if server.WireVersion != nil ***REMOVED***
			if server.WireVersion.Max < SupportedWireVersions.Min ***REMOVED***
				f.compatible.Store(false)
				f.compatibilityErr = fmt.Errorf(
					"server at %s reports wire version %d, but this version of the Go driver requires "+
						"at least %d (MongoDB %s)",
					server.Addr.String(),
					server.WireVersion.Max,
					SupportedWireVersions.Min,
					MinSupportedMongoDBVersion,
				)
				f.Topology.CompatibilityErr = f.compatibilityErr
				return f.Topology, s
			***REMOVED***

			if server.WireVersion.Min > SupportedWireVersions.Max ***REMOVED***
				f.compatible.Store(false)
				f.compatibilityErr = fmt.Errorf(
					"server at %s requires wire version %d, but this version of the Go driver only supports up to %d",
					server.Addr.String(),
					server.WireVersion.Min,
					SupportedWireVersions.Max,
				)
				f.Topology.CompatibilityErr = f.compatibilityErr
				return f.Topology, s
			***REMOVED***
		***REMOVED***
	***REMOVED***

	f.compatible.Store(true)
	f.compatibilityErr = nil
	return f.Topology, updatedDesc
***REMOVED***

func (f *fsm) applyToReplicaSetNoPrimary(s description.Server) description.Server ***REMOVED***
	switch s.Kind ***REMOVED***
	case description.Standalone, description.Mongos:
		f.removeServerByAddr(s.Addr)
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.updateRSWithoutPrimary(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
	***REMOVED***

	return s
***REMOVED***

func (f *fsm) applyToReplicaSetWithPrimary(s description.Server) description.Server ***REMOVED***
	switch s.Kind ***REMOVED***
	case description.Standalone, description.Mongos:
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.updateRSWithPrimaryFromMember(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
		f.checkIfHasPrimary()
	***REMOVED***

	return s
***REMOVED***

func (f *fsm) applyToSharded(s description.Server) description.Server ***REMOVED***
	switch s.Kind ***REMOVED***
	case description.Mongos, description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		f.removeServerByAddr(s.Addr)
	***REMOVED***

	return s
***REMOVED***

func (f *fsm) applyToSingle(s description.Server) description.Server ***REMOVED***
	switch s.Kind ***REMOVED***
	case description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.Mongos:
		if f.SetName != "" ***REMOVED***
			f.removeServerByAddr(s.Addr)
			return s
		***REMOVED***

		f.replaceServer(s)
	case description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		// A replica set name can be provided when creating a direct connection. In this case, if the set name returned
		// by the hello response doesn't match up with the one provided during configuration, the server description
		// is replaced with a default Unknown description.
		//
		// We create a new server description rather than doing s.Kind = description.Unknown because the other fields,
		// such as RTT, need to be cleared for Unknown descriptions as well.
		if f.SetName != "" && f.SetName != s.SetName ***REMOVED***
			s = description.Server***REMOVED***
				Addr: s.Addr,
				Kind: description.Unknown,
			***REMOVED***
		***REMOVED***

		f.replaceServer(s)
	***REMOVED***

	return s
***REMOVED***

func (f *fsm) applyToUnknown(s description.Server) description.Server ***REMOVED***
	switch s.Kind ***REMOVED***
	case description.Mongos:
		f.setKind(description.Sharded)
		f.replaceServer(s)
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.setKind(description.ReplicaSetNoPrimary)
		f.updateRSWithoutPrimary(s)
	case description.Standalone:
		f.updateUnknownWithStandalone(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
	***REMOVED***

	return s
***REMOVED***

func (f *fsm) checkIfHasPrimary() ***REMOVED***
	if _, ok := f.findPrimary(); ok ***REMOVED***
		f.setKind(description.ReplicaSetWithPrimary)
	***REMOVED*** else ***REMOVED***
		f.setKind(description.ReplicaSetNoPrimary)
	***REMOVED***
***REMOVED***

func (f *fsm) updateRSFromPrimary(s description.Server) ***REMOVED***
	if f.SetName == "" ***REMOVED***
		f.SetName = s.SetName
	***REMOVED*** else if f.SetName != s.SetName ***REMOVED***
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	***REMOVED***

	if s.SetVersion != 0 && !s.ElectionID.IsZero() ***REMOVED***
		if f.maxSetVersion > s.SetVersion || bytes.Compare(f.maxElectionID[:], s.ElectionID[:]) == 1 ***REMOVED***
			f.replaceServer(description.Server***REMOVED***
				Addr:      s.Addr,
				LastError: fmt.Errorf("was a primary, but its set version or election id is stale"),
			***REMOVED***)
			f.checkIfHasPrimary()
			return
		***REMOVED***

		f.maxElectionID = s.ElectionID
	***REMOVED***

	if s.SetVersion > f.maxSetVersion ***REMOVED***
		f.maxSetVersion = s.SetVersion
	***REMOVED***

	if j, ok := f.findPrimary(); ok ***REMOVED***
		f.setServer(j, description.Server***REMOVED***
			Addr:      f.Servers[j].Addr,
			LastError: fmt.Errorf("was a primary, but a new primary was discovered"),
		***REMOVED***)
	***REMOVED***

	f.replaceServer(s)

	for j := len(f.Servers) - 1; j >= 0; j-- ***REMOVED***
		found := false
		for _, member := range s.Members ***REMOVED***
			if member == f.Servers[j].Addr ***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			f.removeServer(j)
		***REMOVED***
	***REMOVED***

	for _, member := range s.Members ***REMOVED***
		if _, ok := f.findServer(member); !ok ***REMOVED***
			f.addServer(member)
		***REMOVED***
	***REMOVED***

	f.checkIfHasPrimary()
***REMOVED***

func (f *fsm) updateRSWithPrimaryFromMember(s description.Server) ***REMOVED***
	if f.SetName != s.SetName ***REMOVED***
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	***REMOVED***

	if s.Addr != s.CanonicalAddr ***REMOVED***
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	***REMOVED***

	f.replaceServer(s)

	if _, ok := f.findPrimary(); !ok ***REMOVED***
		f.setKind(description.ReplicaSetNoPrimary)
	***REMOVED***
***REMOVED***

func (f *fsm) updateRSWithoutPrimary(s description.Server) ***REMOVED***
	if f.SetName == "" ***REMOVED***
		f.SetName = s.SetName
	***REMOVED*** else if f.SetName != s.SetName ***REMOVED***
		f.removeServerByAddr(s.Addr)
		return
	***REMOVED***

	for _, member := range s.Members ***REMOVED***
		if _, ok := f.findServer(member); !ok ***REMOVED***
			f.addServer(member)
		***REMOVED***
	***REMOVED***

	if s.Addr != s.CanonicalAddr ***REMOVED***
		f.removeServerByAddr(s.Addr)
		return
	***REMOVED***

	f.replaceServer(s)
***REMOVED***

func (f *fsm) updateUnknownWithStandalone(s description.Server) ***REMOVED***
	if len(f.Servers) > 1 ***REMOVED***
		f.removeServerByAddr(s.Addr)
		return
	***REMOVED***

	f.setKind(description.Single)
	f.replaceServer(s)
***REMOVED***

func (f *fsm) addServer(addr address.Address) ***REMOVED***
	f.Servers = append(f.Servers, description.Server***REMOVED***
		Addr: addr.Canonicalize(),
	***REMOVED***)
***REMOVED***

func (f *fsm) findPrimary() (int, bool) ***REMOVED***
	for i, s := range f.Servers ***REMOVED***
		if s.Kind == description.RSPrimary ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***

	return 0, false
***REMOVED***

func (f *fsm) findServer(addr address.Address) (int, bool) ***REMOVED***
	canon := addr.Canonicalize()
	for i, s := range f.Servers ***REMOVED***
		if canon == s.Addr ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***

	return 0, false
***REMOVED***

func (f *fsm) removeServer(i int) ***REMOVED***
	f.Servers = append(f.Servers[:i], f.Servers[i+1:]...)
***REMOVED***

func (f *fsm) removeServerByAddr(addr address.Address) ***REMOVED***
	if i, ok := f.findServer(addr); ok ***REMOVED***
		f.removeServer(i)
	***REMOVED***
***REMOVED***

func (f *fsm) replaceServer(s description.Server) ***REMOVED***
	if i, ok := f.findServer(s.Addr); ok ***REMOVED***
		f.setServer(i, s)
	***REMOVED***
***REMOVED***

func (f *fsm) setServer(i int, s description.Server) ***REMOVED***
	f.Servers[i] = s
***REMOVED***

func (f *fsm) setKind(k description.TopologyKind) ***REMOVED***
	f.Kind = k
***REMOVED***
