// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Node represents a server session in a linked list
type Node struct ***REMOVED***
	*Server
	next *Node
	prev *Node
***REMOVED***

// topologyDescription is used to track a subset of the fields present in a description.Topology instance that are
// relevant for determining session expiration.
type topologyDescription struct ***REMOVED***
	kind           description.TopologyKind
	timeoutMinutes uint32
***REMOVED***

// Pool is a pool of server sessions that can be reused.
type Pool struct ***REMOVED***
	descChan       <-chan description.Topology
	head           *Node
	tail           *Node
	latestTopology topologyDescription
	mutex          sync.Mutex // mutex to protect list and sessionTimeout

	checkedOut int // number of sessions checked out of pool
***REMOVED***

func (p *Pool) createServerSession() (*Server, error) ***REMOVED***
	s, err := newServerSession()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p.checkedOut++
	return s, nil
***REMOVED***

// NewPool creates a new server session pool
func NewPool(descChan <-chan description.Topology) *Pool ***REMOVED***
	p := &Pool***REMOVED***
		descChan: descChan,
	***REMOVED***

	return p
***REMOVED***

// assumes caller has mutex to protect the pool
func (p *Pool) updateTimeout() ***REMOVED***
	select ***REMOVED***
	case newDesc := <-p.descChan:
		p.latestTopology = topologyDescription***REMOVED***
			kind:           newDesc.Kind,
			timeoutMinutes: newDesc.SessionTimeoutMinutes,
		***REMOVED***
	default:
		// no new description waiting
	***REMOVED***
***REMOVED***

// GetSession retrieves an unexpired session from the pool.
func (p *Pool) GetSession() (*Server, error) ***REMOVED***
	p.mutex.Lock() // prevent changing the linked list while seeing if sessions have expired
	defer p.mutex.Unlock()

	// empty pool
	if p.head == nil && p.tail == nil ***REMOVED***
		return p.createServerSession()
	***REMOVED***

	p.updateTimeout()
	for p.head != nil ***REMOVED***
		// pull session from head of queue and return if it is valid for at least 1 more minute
		if p.head.expired(p.latestTopology) ***REMOVED***
			p.head = p.head.next
			continue
		***REMOVED***

		// found unexpired session
		session := p.head.Server
		if p.head.next != nil ***REMOVED***
			p.head.next.prev = nil
		***REMOVED***
		if p.tail == p.head ***REMOVED***
			p.tail = nil
			p.head = nil
		***REMOVED*** else ***REMOVED***
			p.head = p.head.next
		***REMOVED***

		p.checkedOut++
		return session, nil
	***REMOVED***

	// no valid session found
	p.tail = nil // empty list
	return p.createServerSession()
***REMOVED***

// ReturnSession returns a session to the pool if it has not expired.
func (p *Pool) ReturnSession(ss *Server) ***REMOVED***
	if ss == nil ***REMOVED***
		return
	***REMOVED***

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.checkedOut--
	p.updateTimeout()
	// check sessions at end of queue for expired
	// stop checking after hitting the first valid session
	for p.tail != nil && p.tail.expired(p.latestTopology) ***REMOVED***
		if p.tail.prev != nil ***REMOVED***
			p.tail.prev.next = nil
		***REMOVED***
		p.tail = p.tail.prev
	***REMOVED***

	// session expired
	if ss.expired(p.latestTopology) ***REMOVED***
		return
	***REMOVED***

	// session is dirty
	if ss.Dirty ***REMOVED***
		return
	***REMOVED***

	newNode := &Node***REMOVED***
		Server: ss,
		next:   nil,
		prev:   nil,
	***REMOVED***

	// empty list
	if p.tail == nil ***REMOVED***
		p.head = newNode
		p.tail = newNode
		return
	***REMOVED***

	// at least 1 valid session in list
	newNode.next = p.head
	p.head.prev = newNode
	p.head = newNode
***REMOVED***

// IDSlice returns a slice of session IDs for each session in the pool
func (p *Pool) IDSlice() []bsoncore.Document ***REMOVED***
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var ids []bsoncore.Document
	for node := p.head; node != nil; node = node.next ***REMOVED***
		ids = append(ids, node.SessionID)
	***REMOVED***

	return ids
***REMOVED***

// String implements the Stringer interface
func (p *Pool) String() string ***REMOVED***
	p.mutex.Lock()
	defer p.mutex.Unlock()

	s := ""
	for head := p.head; head != nil; head = head.next ***REMOVED***
		s += head.SessionID.String() + "\n"
	***REMOVED***

	return s
***REMOVED***

// CheckedOut returns number of sessions checked out from pool.
func (p *Pool) CheckedOut() int ***REMOVED***
	return p.checkedOut
***REMOVED***
