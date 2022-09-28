// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"fmt"
	"math"
	"sort"
)

// RFC 7540, Section 5.3.5: the default weight is 16.
const priorityDefaultWeight = 15 // 16 = 15 + 1

// PriorityWriteSchedulerConfig configures a priorityWriteScheduler.
type PriorityWriteSchedulerConfig struct ***REMOVED***
	// MaxClosedNodesInTree controls the maximum number of closed streams to
	// retain in the priority tree. Setting this to zero saves a small amount
	// of memory at the cost of performance.
	//
	// See RFC 7540, Section 5.3.4:
	//   "It is possible for a stream to become closed while prioritization
	//   information ... is in transit. ... This potentially creates suboptimal
	//   prioritization, since the stream could be given a priority that is
	//   different from what is intended. To avoid these problems, an endpoint
	//   SHOULD retain stream prioritization state for a period after streams
	//   become closed. The longer state is retained, the lower the chance that
	//   streams are assigned incorrect or default priority values."
	MaxClosedNodesInTree int

	// MaxIdleNodesInTree controls the maximum number of idle streams to
	// retain in the priority tree. Setting this to zero saves a small amount
	// of memory at the cost of performance.
	//
	// See RFC 7540, Section 5.3.4:
	//   Similarly, streams that are in the "idle" state can be assigned
	//   priority or become a parent of other streams. This allows for the
	//   creation of a grouping node in the dependency tree, which enables
	//   more flexible expressions of priority. Idle streams begin with a
	//   default priority (Section 5.3.5).
	MaxIdleNodesInTree int

	// ThrottleOutOfOrderWrites enables write throttling to help ensure that
	// data is delivered in priority order. This works around a race where
	// stream B depends on stream A and both streams are about to call Write
	// to queue DATA frames. If B wins the race, a naive scheduler would eagerly
	// write as much data from B as possible, but this is suboptimal because A
	// is a higher-priority stream. With throttling enabled, we write a small
	// amount of data from B to minimize the amount of bandwidth that B can
	// steal from A.
	ThrottleOutOfOrderWrites bool
***REMOVED***

// NewPriorityWriteScheduler constructs a WriteScheduler that schedules
// frames by following HTTP/2 priorities as described in RFC 7540 Section 5.3.
// If cfg is nil, default options are used.
func NewPriorityWriteScheduler(cfg *PriorityWriteSchedulerConfig) WriteScheduler ***REMOVED***
	if cfg == nil ***REMOVED***
		// For justification of these defaults, see:
		// https://docs.google.com/document/d/1oLhNg1skaWD4_DtaoCxdSRN5erEXrH-KnLrMwEpOtFY
		cfg = &PriorityWriteSchedulerConfig***REMOVED***
			MaxClosedNodesInTree:     10,
			MaxIdleNodesInTree:       10,
			ThrottleOutOfOrderWrites: false,
		***REMOVED***
	***REMOVED***

	ws := &priorityWriteScheduler***REMOVED***
		nodes:                make(map[uint32]*priorityNode),
		maxClosedNodesInTree: cfg.MaxClosedNodesInTree,
		maxIdleNodesInTree:   cfg.MaxIdleNodesInTree,
		enableWriteThrottle:  cfg.ThrottleOutOfOrderWrites,
	***REMOVED***
	ws.nodes[0] = &ws.root
	if cfg.ThrottleOutOfOrderWrites ***REMOVED***
		ws.writeThrottleLimit = 1024
	***REMOVED*** else ***REMOVED***
		ws.writeThrottleLimit = math.MaxInt32
	***REMOVED***
	return ws
***REMOVED***

type priorityNodeState int

const (
	priorityNodeOpen priorityNodeState = iota
	priorityNodeClosed
	priorityNodeIdle
)

// priorityNode is a node in an HTTP/2 priority tree.
// Each node is associated with a single stream ID.
// See RFC 7540, Section 5.3.
type priorityNode struct ***REMOVED***
	q            writeQueue        // queue of pending frames to write
	id           uint32            // id of the stream, or 0 for the root of the tree
	weight       uint8             // the actual weight is weight+1, so the value is in [1,256]
	state        priorityNodeState // open | closed | idle
	bytes        int64             // number of bytes written by this node, or 0 if closed
	subtreeBytes int64             // sum(node.bytes) of all nodes in this subtree

	// These links form the priority tree.
	parent     *priorityNode
	kids       *priorityNode // start of the kids list
	prev, next *priorityNode // doubly-linked list of siblings
***REMOVED***

func (n *priorityNode) setParent(parent *priorityNode) ***REMOVED***
	if n == parent ***REMOVED***
		panic("setParent to self")
	***REMOVED***
	if n.parent == parent ***REMOVED***
		return
	***REMOVED***
	// Unlink from current parent.
	if parent := n.parent; parent != nil ***REMOVED***
		if n.prev == nil ***REMOVED***
			parent.kids = n.next
		***REMOVED*** else ***REMOVED***
			n.prev.next = n.next
		***REMOVED***
		if n.next != nil ***REMOVED***
			n.next.prev = n.prev
		***REMOVED***
	***REMOVED***
	// Link to new parent.
	// If parent=nil, remove n from the tree.
	// Always insert at the head of parent.kids (this is assumed by walkReadyInOrder).
	n.parent = parent
	if parent == nil ***REMOVED***
		n.next = nil
		n.prev = nil
	***REMOVED*** else ***REMOVED***
		n.next = parent.kids
		n.prev = nil
		if n.next != nil ***REMOVED***
			n.next.prev = n
		***REMOVED***
		parent.kids = n
	***REMOVED***
***REMOVED***

func (n *priorityNode) addBytes(b int64) ***REMOVED***
	n.bytes += b
	for ; n != nil; n = n.parent ***REMOVED***
		n.subtreeBytes += b
	***REMOVED***
***REMOVED***

// walkReadyInOrder iterates over the tree in priority order, calling f for each node
// with a non-empty write queue. When f returns true, this function returns true and the
// walk halts. tmp is used as scratch space for sorting.
//
// f(n, openParent) takes two arguments: the node to visit, n, and a bool that is true
// if any ancestor p of n is still open (ignoring the root node).
func (n *priorityNode) walkReadyInOrder(openParent bool, tmp *[]*priorityNode, f func(*priorityNode, bool) bool) bool ***REMOVED***
	if !n.q.empty() && f(n, openParent) ***REMOVED***
		return true
	***REMOVED***
	if n.kids == nil ***REMOVED***
		return false
	***REMOVED***

	// Don't consider the root "open" when updating openParent since
	// we can't send data frames on the root stream (only control frames).
	if n.id != 0 ***REMOVED***
		openParent = openParent || (n.state == priorityNodeOpen)
	***REMOVED***

	// Common case: only one kid or all kids have the same weight.
	// Some clients don't use weights; other clients (like web browsers)
	// use mostly-linear priority trees.
	w := n.kids.weight
	needSort := false
	for k := n.kids.next; k != nil; k = k.next ***REMOVED***
		if k.weight != w ***REMOVED***
			needSort = true
			break
		***REMOVED***
	***REMOVED***
	if !needSort ***REMOVED***
		for k := n.kids; k != nil; k = k.next ***REMOVED***
			if k.walkReadyInOrder(openParent, tmp, f) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***

	// Uncommon case: sort the child nodes. We remove the kids from the parent,
	// then re-insert after sorting so we can reuse tmp for future sort calls.
	*tmp = (*tmp)[:0]
	for n.kids != nil ***REMOVED***
		*tmp = append(*tmp, n.kids)
		n.kids.setParent(nil)
	***REMOVED***
	sort.Sort(sortPriorityNodeSiblings(*tmp))
	for i := len(*tmp) - 1; i >= 0; i-- ***REMOVED***
		(*tmp)[i].setParent(n) // setParent inserts at the head of n.kids
	***REMOVED***
	for k := n.kids; k != nil; k = k.next ***REMOVED***
		if k.walkReadyInOrder(openParent, tmp, f) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

type sortPriorityNodeSiblings []*priorityNode

func (z sortPriorityNodeSiblings) Len() int      ***REMOVED*** return len(z) ***REMOVED***
func (z sortPriorityNodeSiblings) Swap(i, k int) ***REMOVED*** z[i], z[k] = z[k], z[i] ***REMOVED***
func (z sortPriorityNodeSiblings) Less(i, k int) bool ***REMOVED***
	// Prefer the subtree that has sent fewer bytes relative to its weight.
	// See sections 5.3.2 and 5.3.4.
	wi, bi := float64(z[i].weight+1), float64(z[i].subtreeBytes)
	wk, bk := float64(z[k].weight+1), float64(z[k].subtreeBytes)
	if bi == 0 && bk == 0 ***REMOVED***
		return wi >= wk
	***REMOVED***
	if bk == 0 ***REMOVED***
		return false
	***REMOVED***
	return bi/bk <= wi/wk
***REMOVED***

type priorityWriteScheduler struct ***REMOVED***
	// root is the root of the priority tree, where root.id = 0.
	// The root queues control frames that are not associated with any stream.
	root priorityNode

	// nodes maps stream ids to priority tree nodes.
	nodes map[uint32]*priorityNode

	// maxID is the maximum stream id in nodes.
	maxID uint32

	// lists of nodes that have been closed or are idle, but are kept in
	// the tree for improved prioritization. When the lengths exceed either
	// maxClosedNodesInTree or maxIdleNodesInTree, old nodes are discarded.
	closedNodes, idleNodes []*priorityNode

	// From the config.
	maxClosedNodesInTree int
	maxIdleNodesInTree   int
	writeThrottleLimit   int32
	enableWriteThrottle  bool

	// tmp is scratch space for priorityNode.walkReadyInOrder to reduce allocations.
	tmp []*priorityNode

	// pool of empty queues for reuse.
	queuePool writeQueuePool
***REMOVED***

func (ws *priorityWriteScheduler) OpenStream(streamID uint32, options OpenStreamOptions) ***REMOVED***
	// The stream may be currently idle but cannot be opened or closed.
	if curr := ws.nodes[streamID]; curr != nil ***REMOVED***
		if curr.state != priorityNodeIdle ***REMOVED***
			panic(fmt.Sprintf("stream %d already opened", streamID))
		***REMOVED***
		curr.state = priorityNodeOpen
		return
	***REMOVED***

	// RFC 7540, Section 5.3.5:
	//  "All streams are initially assigned a non-exclusive dependency on stream 0x0.
	//  Pushed streams initially depend on their associated stream. In both cases,
	//  streams are assigned a default weight of 16."
	parent := ws.nodes[options.PusherID]
	if parent == nil ***REMOVED***
		parent = &ws.root
	***REMOVED***
	n := &priorityNode***REMOVED***
		q:      *ws.queuePool.get(),
		id:     streamID,
		weight: priorityDefaultWeight,
		state:  priorityNodeOpen,
	***REMOVED***
	n.setParent(parent)
	ws.nodes[streamID] = n
	if streamID > ws.maxID ***REMOVED***
		ws.maxID = streamID
	***REMOVED***
***REMOVED***

func (ws *priorityWriteScheduler) CloseStream(streamID uint32) ***REMOVED***
	if streamID == 0 ***REMOVED***
		panic("violation of WriteScheduler interface: cannot close stream 0")
	***REMOVED***
	if ws.nodes[streamID] == nil ***REMOVED***
		panic(fmt.Sprintf("violation of WriteScheduler interface: unknown stream %d", streamID))
	***REMOVED***
	if ws.nodes[streamID].state != priorityNodeOpen ***REMOVED***
		panic(fmt.Sprintf("violation of WriteScheduler interface: stream %d already closed", streamID))
	***REMOVED***

	n := ws.nodes[streamID]
	n.state = priorityNodeClosed
	n.addBytes(-n.bytes)

	q := n.q
	ws.queuePool.put(&q)
	n.q.s = nil
	if ws.maxClosedNodesInTree > 0 ***REMOVED***
		ws.addClosedOrIdleNode(&ws.closedNodes, ws.maxClosedNodesInTree, n)
	***REMOVED*** else ***REMOVED***
		ws.removeNode(n)
	***REMOVED***
***REMOVED***

func (ws *priorityWriteScheduler) AdjustStream(streamID uint32, priority PriorityParam) ***REMOVED***
	if streamID == 0 ***REMOVED***
		panic("adjustPriority on root")
	***REMOVED***

	// If streamID does not exist, there are two cases:
	// - A closed stream that has been removed (this will have ID <= maxID)
	// - An idle stream that is being used for "grouping" (this will have ID > maxID)
	n := ws.nodes[streamID]
	if n == nil ***REMOVED***
		if streamID <= ws.maxID || ws.maxIdleNodesInTree == 0 ***REMOVED***
			return
		***REMOVED***
		ws.maxID = streamID
		n = &priorityNode***REMOVED***
			q:      *ws.queuePool.get(),
			id:     streamID,
			weight: priorityDefaultWeight,
			state:  priorityNodeIdle,
		***REMOVED***
		n.setParent(&ws.root)
		ws.nodes[streamID] = n
		ws.addClosedOrIdleNode(&ws.idleNodes, ws.maxIdleNodesInTree, n)
	***REMOVED***

	// Section 5.3.1: A dependency on a stream that is not currently in the tree
	// results in that stream being given a default priority (Section 5.3.5).
	parent := ws.nodes[priority.StreamDep]
	if parent == nil ***REMOVED***
		n.setParent(&ws.root)
		n.weight = priorityDefaultWeight
		return
	***REMOVED***

	// Ignore if the client tries to make a node its own parent.
	if n == parent ***REMOVED***
		return
	***REMOVED***

	// Section 5.3.3:
	//   "If a stream is made dependent on one of its own dependencies, the
	//   formerly dependent stream is first moved to be dependent on the
	//   reprioritized stream's previous parent. The moved dependency retains
	//   its weight."
	//
	// That is: if parent depends on n, move parent to depend on n.parent.
	for x := parent.parent; x != nil; x = x.parent ***REMOVED***
		if x == n ***REMOVED***
			parent.setParent(n.parent)
			break
		***REMOVED***
	***REMOVED***

	// Section 5.3.3: The exclusive flag causes the stream to become the sole
	// dependency of its parent stream, causing other dependencies to become
	// dependent on the exclusive stream.
	if priority.Exclusive ***REMOVED***
		k := parent.kids
		for k != nil ***REMOVED***
			next := k.next
			if k != n ***REMOVED***
				k.setParent(n)
			***REMOVED***
			k = next
		***REMOVED***
	***REMOVED***

	n.setParent(parent)
	n.weight = priority.Weight
***REMOVED***

func (ws *priorityWriteScheduler) Push(wr FrameWriteRequest) ***REMOVED***
	var n *priorityNode
	if wr.isControl() ***REMOVED***
		n = &ws.root
	***REMOVED*** else ***REMOVED***
		id := wr.StreamID()
		n = ws.nodes[id]
		if n == nil ***REMOVED***
			// id is an idle or closed stream. wr should not be a HEADERS or
			// DATA frame. In other case, we push wr onto the root, rather
			// than creating a new priorityNode.
			if wr.DataSize() > 0 ***REMOVED***
				panic("add DATA on non-open stream")
			***REMOVED***
			n = &ws.root
		***REMOVED***
	***REMOVED***
	n.q.push(wr)
***REMOVED***

func (ws *priorityWriteScheduler) Pop() (wr FrameWriteRequest, ok bool) ***REMOVED***
	ws.root.walkReadyInOrder(false, &ws.tmp, func(n *priorityNode, openParent bool) bool ***REMOVED***
		limit := int32(math.MaxInt32)
		if openParent ***REMOVED***
			limit = ws.writeThrottleLimit
		***REMOVED***
		wr, ok = n.q.consume(limit)
		if !ok ***REMOVED***
			return false
		***REMOVED***
		n.addBytes(int64(wr.DataSize()))
		// If B depends on A and B continuously has data available but A
		// does not, gradually increase the throttling limit to allow B to
		// steal more and more bandwidth from A.
		if openParent ***REMOVED***
			ws.writeThrottleLimit += 1024
			if ws.writeThrottleLimit < 0 ***REMOVED***
				ws.writeThrottleLimit = math.MaxInt32
			***REMOVED***
		***REMOVED*** else if ws.enableWriteThrottle ***REMOVED***
			ws.writeThrottleLimit = 1024
		***REMOVED***
		return true
	***REMOVED***)
	return wr, ok
***REMOVED***

func (ws *priorityWriteScheduler) addClosedOrIdleNode(list *[]*priorityNode, maxSize int, n *priorityNode) ***REMOVED***
	if maxSize == 0 ***REMOVED***
		return
	***REMOVED***
	if len(*list) == maxSize ***REMOVED***
		// Remove the oldest node, then shift left.
		ws.removeNode((*list)[0])
		x := (*list)[1:]
		copy(*list, x)
		*list = (*list)[:len(x)]
	***REMOVED***
	*list = append(*list, n)
***REMOVED***

func (ws *priorityWriteScheduler) removeNode(n *priorityNode) ***REMOVED***
	for k := n.kids; k != nil; k = k.next ***REMOVED***
		k.setParent(n.parent)
	***REMOVED***
	n.setParent(nil)
	delete(ws.nodes, n.id)
***REMOVED***
