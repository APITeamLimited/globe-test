// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"sync"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pool generation state constants.
const (
	generationDisconnected int64 = iota
	generationConnected
)

// generationStats represents the version of a pool. It tracks the generation number as well as the number of
// connections that have been created in the generation.
type generationStats struct ***REMOVED***
	generation uint64
	numConns   uint64
***REMOVED***

// poolGenerationMap tracks the version for each service ID present in a pool. For deployments that are not behind a
// load balancer, there is only one service ID: primitive.NilObjectID. For load-balanced deployments, each server behind
// the load balancer will have a unique service ID.
type poolGenerationMap struct ***REMOVED***
	// state must be accessed using the atomic package and should be at the beginning of the struct.
	// - atomic bug: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	// - suggested layout: https://go101.org/article/memory-layout.html
	state         int64
	generationMap map[primitive.ObjectID]*generationStats

	sync.Mutex
***REMOVED***

func newPoolGenerationMap() *poolGenerationMap ***REMOVED***
	pgm := &poolGenerationMap***REMOVED***
		generationMap: make(map[primitive.ObjectID]*generationStats),
	***REMOVED***
	pgm.generationMap[primitive.NilObjectID] = &generationStats***REMOVED******REMOVED***
	return pgm
***REMOVED***

func (p *poolGenerationMap) connect() ***REMOVED***
	atomic.StoreInt64(&p.state, generationConnected)
***REMOVED***

func (p *poolGenerationMap) disconnect() ***REMOVED***
	atomic.StoreInt64(&p.state, generationDisconnected)
***REMOVED***

// addConnection increments the connection count for the generation associated with the given service ID and returns the
// generation number for the connection.
func (p *poolGenerationMap) addConnection(serviceIDPtr *primitive.ObjectID) uint64 ***REMOVED***
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	stats, ok := p.generationMap[serviceID]
	if ok ***REMOVED***
		// If the serviceID is already being tracked, we only need to increment the connection count.
		stats.numConns++
		return stats.generation
	***REMOVED***

	// If the serviceID is untracked, create a new entry with a starting generation number of 0.
	stats = &generationStats***REMOVED***
		numConns: 1,
	***REMOVED***
	p.generationMap[serviceID] = stats
	return 0
***REMOVED***

func (p *poolGenerationMap) removeConnection(serviceIDPtr *primitive.ObjectID) ***REMOVED***
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	stats, ok := p.generationMap[serviceID]
	if !ok ***REMOVED***
		return
	***REMOVED***

	// If the serviceID is being tracked, decrement the connection count and delete this serviceID to prevent the map
	// from growing unboundedly. This case would happen if a server behind a load-balancer was permanently removed
	// and its connections were pruned after a network error or idle timeout.
	stats.numConns--
	if stats.numConns == 0 ***REMOVED***
		delete(p.generationMap, serviceID)
	***REMOVED***
***REMOVED***

func (p *poolGenerationMap) clear(serviceIDPtr *primitive.ObjectID) ***REMOVED***
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	if stats, ok := p.generationMap[serviceID]; ok ***REMOVED***
		stats.generation++
	***REMOVED***
***REMOVED***

func (p *poolGenerationMap) stale(serviceIDPtr *primitive.ObjectID, knownGeneration uint64) bool ***REMOVED***
	// If the map has been disconnected, all connections should be considered stale to ensure that they're closed.
	if atomic.LoadInt64(&p.state) == generationDisconnected ***REMOVED***
		return true
	***REMOVED***

	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	if stats, ok := p.generationMap[serviceID]; ok ***REMOVED***
		return knownGeneration < stats.generation
	***REMOVED***
	return false
***REMOVED***

func (p *poolGenerationMap) getGeneration(serviceIDPtr *primitive.ObjectID) uint64 ***REMOVED***
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	if stats, ok := p.generationMap[serviceID]; ok ***REMOVED***
		return stats.generation
	***REMOVED***
	return 0
***REMOVED***

func (p *poolGenerationMap) getNumConns(serviceIDPtr *primitive.ObjectID) uint64 ***REMOVED***
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	if stats, ok := p.generationMap[serviceID]; ok ***REMOVED***
		return stats.numConns
	***REMOVED***
	return 0
***REMOVED***

func getServiceID(oid *primitive.ObjectID) primitive.ObjectID ***REMOVED***
	if oid == nil ***REMOVED***
		return primitive.NilObjectID
	***REMOVED***
	return *oid
***REMOVED***
