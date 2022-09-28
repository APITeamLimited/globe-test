// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/tag"
)

// ServerSelector is an interface implemented by types that can perform server selection given a topology description
// and list of candidate servers. The selector should filter the provided candidates list and return a subset that
// matches some criteria.
type ServerSelector interface ***REMOVED***
	SelectServer(Topology, []Server) ([]Server, error)
***REMOVED***

// ServerSelectorFunc is a function that can be used as a ServerSelector.
type ServerSelectorFunc func(Topology, []Server) ([]Server, error)

// SelectServer implements the ServerSelector interface.
func (ssf ServerSelectorFunc) SelectServer(t Topology, s []Server) ([]Server, error) ***REMOVED***
	return ssf(t, s)
***REMOVED***

type compositeSelector struct ***REMOVED***
	selectors []ServerSelector
***REMOVED***

// CompositeSelector combines multiple selectors into a single selector by applying them in order to the candidates
// list.
//
// For example, if the initial candidates list is [s0, s1, s2, s3] and two selectors are provided where the first
// matches s0 and s1 and the second matches s1 and s2, the following would occur during server selection:
//
// 1. firstSelector([s0, s1, s2, s3]) -> [s0, s1]
// 2. secondSelector([s0, s1]) -> [s1]
//
// The final list of candidates returned by the composite selector would be [s1].
func CompositeSelector(selectors []ServerSelector) ServerSelector ***REMOVED***
	return &compositeSelector***REMOVED***selectors: selectors***REMOVED***
***REMOVED***

func (cs *compositeSelector) SelectServer(t Topology, candidates []Server) ([]Server, error) ***REMOVED***
	var err error
	for _, sel := range cs.selectors ***REMOVED***
		candidates, err = sel.SelectServer(t, candidates)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return candidates, nil
***REMOVED***

type latencySelector struct ***REMOVED***
	latency time.Duration
***REMOVED***

// LatencySelector creates a ServerSelector which selects servers based on their average RTT values.
func LatencySelector(latency time.Duration) ServerSelector ***REMOVED***
	return &latencySelector***REMOVED***latency: latency***REMOVED***
***REMOVED***

func (ls *latencySelector) SelectServer(t Topology, candidates []Server) ([]Server, error) ***REMOVED***
	if ls.latency < 0 ***REMOVED***
		return candidates, nil
	***REMOVED***
	if t.Kind == LoadBalanced ***REMOVED***
		// In LoadBalanced mode, there should only be one server in the topology and it must be selected.
		return candidates, nil
	***REMOVED***

	switch len(candidates) ***REMOVED***
	case 0, 1:
		return candidates, nil
	default:
		min := time.Duration(math.MaxInt64)
		for _, candidate := range candidates ***REMOVED***
			if candidate.AverageRTTSet ***REMOVED***
				if candidate.AverageRTT < min ***REMOVED***
					min = candidate.AverageRTT
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if min == math.MaxInt64 ***REMOVED***
			return candidates, nil
		***REMOVED***

		max := min + ls.latency

		var result []Server
		for _, candidate := range candidates ***REMOVED***
			if candidate.AverageRTTSet ***REMOVED***
				if candidate.AverageRTT <= max ***REMOVED***
					result = append(result, candidate)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return result, nil
	***REMOVED***
***REMOVED***

// WriteSelector selects all the writable servers.
func WriteSelector() ServerSelector ***REMOVED***
	return ServerSelectorFunc(func(t Topology, candidates []Server) ([]Server, error) ***REMOVED***
		switch t.Kind ***REMOVED***
		case Single, LoadBalanced:
			return candidates, nil
		default:
			result := []Server***REMOVED******REMOVED***
			for _, candidate := range candidates ***REMOVED***
				switch candidate.Kind ***REMOVED***
				case Mongos, RSPrimary, Standalone:
					result = append(result, candidate)
				***REMOVED***
			***REMOVED***
			return result, nil
		***REMOVED***
	***REMOVED***)
***REMOVED***

// ReadPrefSelector selects servers based on the provided read preference.
func ReadPrefSelector(rp *readpref.ReadPref) ServerSelector ***REMOVED***
	return readPrefSelector(rp, false)
***REMOVED***

// OutputAggregateSelector selects servers based on the provided read preference given that the underlying operation is
// aggregate with an output stage.
func OutputAggregateSelector(rp *readpref.ReadPref) ServerSelector ***REMOVED***
	return readPrefSelector(rp, true)
***REMOVED***

func readPrefSelector(rp *readpref.ReadPref, isOutputAggregate bool) ServerSelector ***REMOVED***
	return ServerSelectorFunc(func(t Topology, candidates []Server) ([]Server, error) ***REMOVED***
		if t.Kind == LoadBalanced ***REMOVED***
			// In LoadBalanced mode, there should only be one server in the topology and it must be selected. We check
			// this before checking MaxStaleness support because there's no monitoring in this mode, so the candidate
			// server wouldn't have a wire version set, which would result in an error.
			return candidates, nil
		***REMOVED***

		if _, set := rp.MaxStaleness(); set ***REMOVED***
			for _, s := range candidates ***REMOVED***
				if s.Kind != Unknown ***REMOVED***
					if err := maxStalenessSupported(s.WireVersion); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		switch t.Kind ***REMOVED***
		case Single:
			return candidates, nil
		case ReplicaSetNoPrimary, ReplicaSetWithPrimary:
			return selectForReplicaSet(rp, isOutputAggregate, t, candidates)
		case Sharded:
			return selectByKind(candidates, Mongos), nil
		***REMOVED***

		return nil, nil
	***REMOVED***)
***REMOVED***

// maxStalenessSupported returns an error if the given server version does not support max staleness.
func maxStalenessSupported(wireVersion *VersionRange) error ***REMOVED***
	if wireVersion != nil && wireVersion.Max < 5 ***REMOVED***
		return fmt.Errorf("max staleness is only supported for servers 3.4 or newer")
	***REMOVED***

	return nil
***REMOVED***

func selectForReplicaSet(rp *readpref.ReadPref, isOutputAggregate bool, t Topology, candidates []Server) ([]Server, error) ***REMOVED***
	if err := verifyMaxStaleness(rp, t); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// If underlying operation is an aggregate with an output stage, only apply read preference
	// if all candidates are 5.0+. Otherwise, operate under primary read preference.
	if isOutputAggregate ***REMOVED***
		for _, s := range candidates ***REMOVED***
			if s.WireVersion.Max < 13 ***REMOVED***
				return selectByKind(candidates, RSPrimary), nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	switch rp.Mode() ***REMOVED***
	case readpref.PrimaryMode:
		return selectByKind(candidates, RSPrimary), nil
	case readpref.PrimaryPreferredMode:
		selected := selectByKind(candidates, RSPrimary)

		if len(selected) == 0 ***REMOVED***
			selected = selectSecondaries(rp, candidates)
			return selectByTagSet(selected, rp.TagSets()), nil
		***REMOVED***

		return selected, nil
	case readpref.SecondaryPreferredMode:
		selected := selectSecondaries(rp, candidates)
		selected = selectByTagSet(selected, rp.TagSets())
		if len(selected) > 0 ***REMOVED***
			return selected, nil
		***REMOVED***
		return selectByKind(candidates, RSPrimary), nil
	case readpref.SecondaryMode:
		selected := selectSecondaries(rp, candidates)
		return selectByTagSet(selected, rp.TagSets()), nil
	case readpref.NearestMode:
		selected := selectByKind(candidates, RSPrimary)
		selected = append(selected, selectSecondaries(rp, candidates)...)
		return selectByTagSet(selected, rp.TagSets()), nil
	***REMOVED***

	return nil, fmt.Errorf("unsupported mode: %d", rp.Mode())
***REMOVED***

func selectSecondaries(rp *readpref.ReadPref, candidates []Server) []Server ***REMOVED***
	secondaries := selectByKind(candidates, RSSecondary)
	if len(secondaries) == 0 ***REMOVED***
		return secondaries
	***REMOVED***
	if maxStaleness, set := rp.MaxStaleness(); set ***REMOVED***
		primaries := selectByKind(candidates, RSPrimary)
		if len(primaries) == 0 ***REMOVED***
			baseTime := secondaries[0].LastWriteTime
			for i := 1; i < len(secondaries); i++ ***REMOVED***
				if secondaries[i].LastWriteTime.After(baseTime) ***REMOVED***
					baseTime = secondaries[i].LastWriteTime
				***REMOVED***
			***REMOVED***

			var selected []Server
			for _, secondary := range secondaries ***REMOVED***
				estimatedStaleness := baseTime.Sub(secondary.LastWriteTime) + secondary.HeartbeatInterval
				if estimatedStaleness <= maxStaleness ***REMOVED***
					selected = append(selected, secondary)
				***REMOVED***
			***REMOVED***

			return selected
		***REMOVED***

		primary := primaries[0]

		var selected []Server
		for _, secondary := range secondaries ***REMOVED***
			estimatedStaleness := secondary.LastUpdateTime.Sub(secondary.LastWriteTime) - primary.LastUpdateTime.Sub(primary.LastWriteTime) + secondary.HeartbeatInterval
			if estimatedStaleness <= maxStaleness ***REMOVED***
				selected = append(selected, secondary)
			***REMOVED***
		***REMOVED***
		return selected
	***REMOVED***

	return secondaries
***REMOVED***

func selectByTagSet(candidates []Server, tagSets []tag.Set) []Server ***REMOVED***
	if len(tagSets) == 0 ***REMOVED***
		return candidates
	***REMOVED***

	for _, ts := range tagSets ***REMOVED***
		// If this tag set is empty, we can take a fast path because the empty list is a subset of all tag sets, so
		// all candidate servers will be selected.
		if len(ts) == 0 ***REMOVED***
			return candidates
		***REMOVED***

		var results []Server
		for _, s := range candidates ***REMOVED***
			// ts is non-empty, so only servers with a non-empty set of tags need to be checked.
			if len(s.Tags) > 0 && s.Tags.ContainsAll(ts) ***REMOVED***
				results = append(results, s)
			***REMOVED***
		***REMOVED***

		if len(results) > 0 ***REMOVED***
			return results
		***REMOVED***
	***REMOVED***

	return []Server***REMOVED******REMOVED***
***REMOVED***

func selectByKind(candidates []Server, kind ServerKind) []Server ***REMOVED***
	var result []Server
	for _, s := range candidates ***REMOVED***
		if s.Kind == kind ***REMOVED***
			result = append(result, s)
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***

func verifyMaxStaleness(rp *readpref.ReadPref, t Topology) error ***REMOVED***
	maxStaleness, set := rp.MaxStaleness()
	if !set ***REMOVED***
		return nil
	***REMOVED***

	if maxStaleness < 90*time.Second ***REMOVED***
		return fmt.Errorf("max staleness (%s) must be greater than or equal to 90s", maxStaleness)
	***REMOVED***

	if len(t.Servers) < 1 ***REMOVED***
		// Maybe we should return an error here instead?
		return nil
	***REMOVED***

	// we'll assume all candidates have the same heartbeat interval.
	s := t.Servers[0]
	idleWritePeriod := 10 * time.Second

	if maxStaleness < s.HeartbeatInterval+idleWritePeriod ***REMOVED***
		return fmt.Errorf(
			"max staleness (%s) must be greater than or equal to the heartbeat interval (%s) plus idle write period (%s)",
			maxStaleness, s.HeartbeatInterval, idleWritePeriod,
		)
	***REMOVED***

	return nil
***REMOVED***
