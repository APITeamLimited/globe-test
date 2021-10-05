/*
 *
 * Copyright 2018 gRPC authors.
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

// Package channelz defines APIs for enabling channelz service, entry
// registration/deletion, and accessing channelz data. It also defines channelz
// metric struct formats.
//
// All APIs in this package are experimental.
package channelz

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/grpclog"
)

const (
	defaultMaxTraceEntry int32 = 30
)

var (
	db    dbWrapper
	idGen idGenerator
	// EntryPerPage defines the number of channelz entries to be shown on a web page.
	EntryPerPage  = int64(50)
	curState      int32
	maxTraceEntry = defaultMaxTraceEntry
)

// TurnOn turns on channelz data collection.
func TurnOn() ***REMOVED***
	if !IsOn() ***REMOVED***
		NewChannelzStorage()
		atomic.StoreInt32(&curState, 1)
	***REMOVED***
***REMOVED***

// IsOn returns whether channelz data collection is on.
func IsOn() bool ***REMOVED***
	return atomic.CompareAndSwapInt32(&curState, 1, 1)
***REMOVED***

// SetMaxTraceEntry sets maximum number of trace entry per entity (i.e. channel/subchannel).
// Setting it to 0 will disable channel tracing.
func SetMaxTraceEntry(i int32) ***REMOVED***
	atomic.StoreInt32(&maxTraceEntry, i)
***REMOVED***

// ResetMaxTraceEntryToDefault resets the maximum number of trace entry per entity to default.
func ResetMaxTraceEntryToDefault() ***REMOVED***
	atomic.StoreInt32(&maxTraceEntry, defaultMaxTraceEntry)
***REMOVED***

func getMaxTraceEntry() int ***REMOVED***
	i := atomic.LoadInt32(&maxTraceEntry)
	return int(i)
***REMOVED***

// dbWarpper wraps around a reference to internal channelz data storage, and
// provide synchronized functionality to set and get the reference.
type dbWrapper struct ***REMOVED***
	mu sync.RWMutex
	DB *channelMap
***REMOVED***

func (d *dbWrapper) set(db *channelMap) ***REMOVED***
	d.mu.Lock()
	d.DB = db
	d.mu.Unlock()
***REMOVED***

func (d *dbWrapper) get() *channelMap ***REMOVED***
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.DB
***REMOVED***

// NewChannelzStorage initializes channelz data storage and id generator.
//
// This function returns a cleanup function to wait for all channelz state to be reset by the
// grpc goroutines when those entities get closed. By using this cleanup function, we make sure tests
// don't mess up each other, i.e. lingering goroutine from previous test doing entity removal happen
// to remove some entity just register by the new test, since the id space is the same.
//
// Note: This function is exported for testing purpose only. User should not call
// it in most cases.
func NewChannelzStorage() (cleanup func() error) ***REMOVED***
	db.set(&channelMap***REMOVED***
		topLevelChannels: make(map[int64]struct***REMOVED******REMOVED***),
		channels:         make(map[int64]*channel),
		listenSockets:    make(map[int64]*listenSocket),
		normalSockets:    make(map[int64]*normalSocket),
		servers:          make(map[int64]*server),
		subChannels:      make(map[int64]*subChannel),
	***REMOVED***)
	idGen.reset()
	return func() error ***REMOVED***
		var err error
		cm := db.get()
		if cm == nil ***REMOVED***
			return nil
		***REMOVED***
		for i := 0; i < 1000; i++ ***REMOVED***
			cm.mu.Lock()
			if len(cm.topLevelChannels) == 0 && len(cm.servers) == 0 && len(cm.channels) == 0 && len(cm.subChannels) == 0 && len(cm.listenSockets) == 0 && len(cm.normalSockets) == 0 ***REMOVED***
				cm.mu.Unlock()
				// all things stored in the channelz map have been cleared.
				return nil
			***REMOVED***
			cm.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
		***REMOVED***

		cm.mu.Lock()
		err = fmt.Errorf("after 10s the channelz map has not been cleaned up yet, topchannels: %d, servers: %d, channels: %d, subchannels: %d, listen sockets: %d, normal sockets: %d", len(cm.topLevelChannels), len(cm.servers), len(cm.channels), len(cm.subChannels), len(cm.listenSockets), len(cm.normalSockets))
		cm.mu.Unlock()
		return err
	***REMOVED***
***REMOVED***

// GetTopChannels returns a slice of top channel's ChannelMetric, along with a
// boolean indicating whether there's more top channels to be queried for.
//
// The arg id specifies that only top channel with id at or above it will be included
// in the result. The returned slice is up to a length of the arg maxResults or
// EntryPerPage if maxResults is zero, and is sorted in ascending id order.
func GetTopChannels(id int64, maxResults int64) ([]*ChannelMetric, bool) ***REMOVED***
	return db.get().GetTopChannels(id, maxResults)
***REMOVED***

// GetServers returns a slice of server's ServerMetric, along with a
// boolean indicating whether there's more servers to be queried for.
//
// The arg id specifies that only server with id at or above it will be included
// in the result. The returned slice is up to a length of the arg maxResults or
// EntryPerPage if maxResults is zero, and is sorted in ascending id order.
func GetServers(id int64, maxResults int64) ([]*ServerMetric, bool) ***REMOVED***
	return db.get().GetServers(id, maxResults)
***REMOVED***

// GetServerSockets returns a slice of server's (identified by id) normal socket's
// SocketMetric, along with a boolean indicating whether there's more sockets to
// be queried for.
//
// The arg startID specifies that only sockets with id at or above it will be
// included in the result. The returned slice is up to a length of the arg maxResults
// or EntryPerPage if maxResults is zero, and is sorted in ascending id order.
func GetServerSockets(id int64, startID int64, maxResults int64) ([]*SocketMetric, bool) ***REMOVED***
	return db.get().GetServerSockets(id, startID, maxResults)
***REMOVED***

// GetChannel returns the ChannelMetric for the channel (identified by id).
func GetChannel(id int64) *ChannelMetric ***REMOVED***
	return db.get().GetChannel(id)
***REMOVED***

// GetSubChannel returns the SubChannelMetric for the subchannel (identified by id).
func GetSubChannel(id int64) *SubChannelMetric ***REMOVED***
	return db.get().GetSubChannel(id)
***REMOVED***

// GetSocket returns the SocketInternalMetric for the socket (identified by id).
func GetSocket(id int64) *SocketMetric ***REMOVED***
	return db.get().GetSocket(id)
***REMOVED***

// GetServer returns the ServerMetric for the server (identified by id).
func GetServer(id int64) *ServerMetric ***REMOVED***
	return db.get().GetServer(id)
***REMOVED***

// RegisterChannel registers the given channel c in channelz database with ref
// as its reference name, and add it to the child list of its parent (identified
// by pid). pid = 0 means no parent. It returns the unique channelz tracking id
// assigned to this channel.
func RegisterChannel(c Channel, pid int64, ref string) int64 ***REMOVED***
	id := idGen.genID()
	cn := &channel***REMOVED***
		refName:     ref,
		c:           c,
		subChans:    make(map[int64]string),
		nestedChans: make(map[int64]string),
		id:          id,
		pid:         pid,
		trace:       &channelTrace***REMOVED***createdTime: time.Now(), events: make([]*TraceEvent, 0, getMaxTraceEntry())***REMOVED***,
	***REMOVED***
	if pid == 0 ***REMOVED***
		db.get().addChannel(id, cn, true, pid, ref)
	***REMOVED*** else ***REMOVED***
		db.get().addChannel(id, cn, false, pid, ref)
	***REMOVED***
	return id
***REMOVED***

// RegisterSubChannel registers the given channel c in channelz database with ref
// as its reference name, and add it to the child list of its parent (identified
// by pid). It returns the unique channelz tracking id assigned to this subchannel.
func RegisterSubChannel(c Channel, pid int64, ref string) int64 ***REMOVED***
	if pid == 0 ***REMOVED***
		logger.Error("a SubChannel's parent id cannot be 0")
		return 0
	***REMOVED***
	id := idGen.genID()
	sc := &subChannel***REMOVED***
		refName: ref,
		c:       c,
		sockets: make(map[int64]string),
		id:      id,
		pid:     pid,
		trace:   &channelTrace***REMOVED***createdTime: time.Now(), events: make([]*TraceEvent, 0, getMaxTraceEntry())***REMOVED***,
	***REMOVED***
	db.get().addSubChannel(id, sc, pid, ref)
	return id
***REMOVED***

// RegisterServer registers the given server s in channelz database. It returns
// the unique channelz tracking id assigned to this server.
func RegisterServer(s Server, ref string) int64 ***REMOVED***
	id := idGen.genID()
	svr := &server***REMOVED***
		refName:       ref,
		s:             s,
		sockets:       make(map[int64]string),
		listenSockets: make(map[int64]string),
		id:            id,
	***REMOVED***
	db.get().addServer(id, svr)
	return id
***REMOVED***

// RegisterListenSocket registers the given listen socket s in channelz database
// with ref as its reference name, and add it to the child list of its parent
// (identified by pid). It returns the unique channelz tracking id assigned to
// this listen socket.
func RegisterListenSocket(s Socket, pid int64, ref string) int64 ***REMOVED***
	if pid == 0 ***REMOVED***
		logger.Error("a ListenSocket's parent id cannot be 0")
		return 0
	***REMOVED***
	id := idGen.genID()
	ls := &listenSocket***REMOVED***refName: ref, s: s, id: id, pid: pid***REMOVED***
	db.get().addListenSocket(id, ls, pid, ref)
	return id
***REMOVED***

// RegisterNormalSocket registers the given normal socket s in channelz database
// with ref as its reference name, and add it to the child list of its parent
// (identified by pid). It returns the unique channelz tracking id assigned to
// this normal socket.
func RegisterNormalSocket(s Socket, pid int64, ref string) int64 ***REMOVED***
	if pid == 0 ***REMOVED***
		logger.Error("a NormalSocket's parent id cannot be 0")
		return 0
	***REMOVED***
	id := idGen.genID()
	ns := &normalSocket***REMOVED***refName: ref, s: s, id: id, pid: pid***REMOVED***
	db.get().addNormalSocket(id, ns, pid, ref)
	return id
***REMOVED***

// RemoveEntry removes an entry with unique channelz trakcing id to be id from
// channelz database.
func RemoveEntry(id int64) ***REMOVED***
	db.get().removeEntry(id)
***REMOVED***

// TraceEventDesc is what the caller of AddTraceEvent should provide to describe the event to be added
// to the channel trace.
// The Parent field is optional. It is used for event that will be recorded in the entity's parent
// trace also.
type TraceEventDesc struct ***REMOVED***
	Desc     string
	Severity Severity
	Parent   *TraceEventDesc
***REMOVED***

// AddTraceEvent adds trace related to the entity with specified id, using the provided TraceEventDesc.
func AddTraceEvent(l grpclog.DepthLoggerV2, id int64, depth int, desc *TraceEventDesc) ***REMOVED***
	for d := desc; d != nil; d = d.Parent ***REMOVED***
		switch d.Severity ***REMOVED***
		case CtUnknown, CtInfo:
			l.InfoDepth(depth+1, d.Desc)
		case CtWarning:
			l.WarningDepth(depth+1, d.Desc)
		case CtError:
			l.ErrorDepth(depth+1, d.Desc)
		***REMOVED***
	***REMOVED***
	if getMaxTraceEntry() == 0 ***REMOVED***
		return
	***REMOVED***
	db.get().traceEvent(id, desc)
***REMOVED***

// channelMap is the storage data structure for channelz.
// Methods of channelMap can be divided in two two categories with respect to locking.
// 1. Methods acquire the global lock.
// 2. Methods that can only be called when global lock is held.
// A second type of method need always to be called inside a first type of method.
type channelMap struct ***REMOVED***
	mu               sync.RWMutex
	topLevelChannels map[int64]struct***REMOVED******REMOVED***
	servers          map[int64]*server
	channels         map[int64]*channel
	subChannels      map[int64]*subChannel
	listenSockets    map[int64]*listenSocket
	normalSockets    map[int64]*normalSocket
***REMOVED***

func (c *channelMap) addServer(id int64, s *server) ***REMOVED***
	c.mu.Lock()
	s.cm = c
	c.servers[id] = s
	c.mu.Unlock()
***REMOVED***

func (c *channelMap) addChannel(id int64, cn *channel, isTopChannel bool, pid int64, ref string) ***REMOVED***
	c.mu.Lock()
	cn.cm = c
	cn.trace.cm = c
	c.channels[id] = cn
	if isTopChannel ***REMOVED***
		c.topLevelChannels[id] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		c.findEntry(pid).addChild(id, cn)
	***REMOVED***
	c.mu.Unlock()
***REMOVED***

func (c *channelMap) addSubChannel(id int64, sc *subChannel, pid int64, ref string) ***REMOVED***
	c.mu.Lock()
	sc.cm = c
	sc.trace.cm = c
	c.subChannels[id] = sc
	c.findEntry(pid).addChild(id, sc)
	c.mu.Unlock()
***REMOVED***

func (c *channelMap) addListenSocket(id int64, ls *listenSocket, pid int64, ref string) ***REMOVED***
	c.mu.Lock()
	ls.cm = c
	c.listenSockets[id] = ls
	c.findEntry(pid).addChild(id, ls)
	c.mu.Unlock()
***REMOVED***

func (c *channelMap) addNormalSocket(id int64, ns *normalSocket, pid int64, ref string) ***REMOVED***
	c.mu.Lock()
	ns.cm = c
	c.normalSockets[id] = ns
	c.findEntry(pid).addChild(id, ns)
	c.mu.Unlock()
***REMOVED***

// removeEntry triggers the removal of an entry, which may not indeed delete the entry, if it has to
// wait on the deletion of its children and until no other entity's channel trace references it.
// It may lead to a chain of entry deletion. For example, deleting the last socket of a gracefully
// shutting down server will lead to the server being also deleted.
func (c *channelMap) removeEntry(id int64) ***REMOVED***
	c.mu.Lock()
	c.findEntry(id).triggerDelete()
	c.mu.Unlock()
***REMOVED***

// c.mu must be held by the caller
func (c *channelMap) decrTraceRefCount(id int64) ***REMOVED***
	e := c.findEntry(id)
	if v, ok := e.(tracedChannel); ok ***REMOVED***
		v.decrTraceRefCount()
		e.deleteSelfIfReady()
	***REMOVED***
***REMOVED***

// c.mu must be held by the caller.
func (c *channelMap) findEntry(id int64) entry ***REMOVED***
	var v entry
	var ok bool
	if v, ok = c.channels[id]; ok ***REMOVED***
		return v
	***REMOVED***
	if v, ok = c.subChannels[id]; ok ***REMOVED***
		return v
	***REMOVED***
	if v, ok = c.servers[id]; ok ***REMOVED***
		return v
	***REMOVED***
	if v, ok = c.listenSockets[id]; ok ***REMOVED***
		return v
	***REMOVED***
	if v, ok = c.normalSockets[id]; ok ***REMOVED***
		return v
	***REMOVED***
	return &dummyEntry***REMOVED***idNotFound: id***REMOVED***
***REMOVED***

// c.mu must be held by the caller
// deleteEntry simply deletes an entry from the channelMap. Before calling this
// method, caller must check this entry is ready to be deleted, i.e removeEntry()
// has been called on it, and no children still exist.
// Conditionals are ordered by the expected frequency of deletion of each entity
// type, in order to optimize performance.
func (c *channelMap) deleteEntry(id int64) ***REMOVED***
	var ok bool
	if _, ok = c.normalSockets[id]; ok ***REMOVED***
		delete(c.normalSockets, id)
		return
	***REMOVED***
	if _, ok = c.subChannels[id]; ok ***REMOVED***
		delete(c.subChannels, id)
		return
	***REMOVED***
	if _, ok = c.channels[id]; ok ***REMOVED***
		delete(c.channels, id)
		delete(c.topLevelChannels, id)
		return
	***REMOVED***
	if _, ok = c.listenSockets[id]; ok ***REMOVED***
		delete(c.listenSockets, id)
		return
	***REMOVED***
	if _, ok = c.servers[id]; ok ***REMOVED***
		delete(c.servers, id)
		return
	***REMOVED***
***REMOVED***

func (c *channelMap) traceEvent(id int64, desc *TraceEventDesc) ***REMOVED***
	c.mu.Lock()
	child := c.findEntry(id)
	childTC, ok := child.(tracedChannel)
	if !ok ***REMOVED***
		c.mu.Unlock()
		return
	***REMOVED***
	childTC.getChannelTrace().append(&TraceEvent***REMOVED***Desc: desc.Desc, Severity: desc.Severity, Timestamp: time.Now()***REMOVED***)
	if desc.Parent != nil ***REMOVED***
		parent := c.findEntry(child.getParentID())
		var chanType RefChannelType
		switch child.(type) ***REMOVED***
		case *channel:
			chanType = RefChannel
		case *subChannel:
			chanType = RefSubChannel
		***REMOVED***
		if parentTC, ok := parent.(tracedChannel); ok ***REMOVED***
			parentTC.getChannelTrace().append(&TraceEvent***REMOVED***
				Desc:      desc.Parent.Desc,
				Severity:  desc.Parent.Severity,
				Timestamp: time.Now(),
				RefID:     id,
				RefName:   childTC.getRefName(),
				RefType:   chanType,
			***REMOVED***)
			childTC.incrTraceRefCount()
		***REMOVED***
	***REMOVED***
	c.mu.Unlock()
***REMOVED***

type int64Slice []int64

func (s int64Slice) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s int64Slice) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s int64Slice) Less(i, j int) bool ***REMOVED*** return s[i] < s[j] ***REMOVED***

func copyMap(m map[int64]string) map[int64]string ***REMOVED***
	n := make(map[int64]string)
	for k, v := range m ***REMOVED***
		n[k] = v
	***REMOVED***
	return n
***REMOVED***

func min(a, b int64) int64 ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func (c *channelMap) GetTopChannels(id int64, maxResults int64) ([]*ChannelMetric, bool) ***REMOVED***
	if maxResults <= 0 ***REMOVED***
		maxResults = EntryPerPage
	***REMOVED***
	c.mu.RLock()
	l := int64(len(c.topLevelChannels))
	ids := make([]int64, 0, l)
	cns := make([]*channel, 0, min(l, maxResults))

	for k := range c.topLevelChannels ***REMOVED***
		ids = append(ids, k)
	***REMOVED***
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool ***REMOVED*** return ids[i] >= id ***REMOVED***)
	count := int64(0)
	var end bool
	var t []*ChannelMetric
	for i, v := range ids[idx:] ***REMOVED***
		if count == maxResults ***REMOVED***
			break
		***REMOVED***
		if cn, ok := c.channels[v]; ok ***REMOVED***
			cns = append(cns, cn)
			t = append(t, &ChannelMetric***REMOVED***
				NestedChans: copyMap(cn.nestedChans),
				SubChans:    copyMap(cn.subChans),
			***REMOVED***)
			count++
		***REMOVED***
		if i == len(ids[idx:])-1 ***REMOVED***
			end = true
			break
		***REMOVED***
	***REMOVED***
	c.mu.RUnlock()
	if count == 0 ***REMOVED***
		end = true
	***REMOVED***

	for i, cn := range cns ***REMOVED***
		t[i].ChannelData = cn.c.ChannelzMetric()
		t[i].ID = cn.id
		t[i].RefName = cn.refName
		t[i].Trace = cn.trace.dumpData()
	***REMOVED***
	return t, end
***REMOVED***

func (c *channelMap) GetServers(id, maxResults int64) ([]*ServerMetric, bool) ***REMOVED***
	if maxResults <= 0 ***REMOVED***
		maxResults = EntryPerPage
	***REMOVED***
	c.mu.RLock()
	l := int64(len(c.servers))
	ids := make([]int64, 0, l)
	ss := make([]*server, 0, min(l, maxResults))
	for k := range c.servers ***REMOVED***
		ids = append(ids, k)
	***REMOVED***
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool ***REMOVED*** return ids[i] >= id ***REMOVED***)
	count := int64(0)
	var end bool
	var s []*ServerMetric
	for i, v := range ids[idx:] ***REMOVED***
		if count == maxResults ***REMOVED***
			break
		***REMOVED***
		if svr, ok := c.servers[v]; ok ***REMOVED***
			ss = append(ss, svr)
			s = append(s, &ServerMetric***REMOVED***
				ListenSockets: copyMap(svr.listenSockets),
			***REMOVED***)
			count++
		***REMOVED***
		if i == len(ids[idx:])-1 ***REMOVED***
			end = true
			break
		***REMOVED***
	***REMOVED***
	c.mu.RUnlock()
	if count == 0 ***REMOVED***
		end = true
	***REMOVED***

	for i, svr := range ss ***REMOVED***
		s[i].ServerData = svr.s.ChannelzMetric()
		s[i].ID = svr.id
		s[i].RefName = svr.refName
	***REMOVED***
	return s, end
***REMOVED***

func (c *channelMap) GetServerSockets(id int64, startID int64, maxResults int64) ([]*SocketMetric, bool) ***REMOVED***
	if maxResults <= 0 ***REMOVED***
		maxResults = EntryPerPage
	***REMOVED***
	var svr *server
	var ok bool
	c.mu.RLock()
	if svr, ok = c.servers[id]; !ok ***REMOVED***
		// server with id doesn't exist.
		c.mu.RUnlock()
		return nil, true
	***REMOVED***
	svrskts := svr.sockets
	l := int64(len(svrskts))
	ids := make([]int64, 0, l)
	sks := make([]*normalSocket, 0, min(l, maxResults))
	for k := range svrskts ***REMOVED***
		ids = append(ids, k)
	***REMOVED***
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool ***REMOVED*** return ids[i] >= startID ***REMOVED***)
	count := int64(0)
	var end bool
	for i, v := range ids[idx:] ***REMOVED***
		if count == maxResults ***REMOVED***
			break
		***REMOVED***
		if ns, ok := c.normalSockets[v]; ok ***REMOVED***
			sks = append(sks, ns)
			count++
		***REMOVED***
		if i == len(ids[idx:])-1 ***REMOVED***
			end = true
			break
		***REMOVED***
	***REMOVED***
	c.mu.RUnlock()
	if count == 0 ***REMOVED***
		end = true
	***REMOVED***
	s := make([]*SocketMetric, 0, len(sks))
	for _, ns := range sks ***REMOVED***
		sm := &SocketMetric***REMOVED******REMOVED***
		sm.SocketData = ns.s.ChannelzMetric()
		sm.ID = ns.id
		sm.RefName = ns.refName
		s = append(s, sm)
	***REMOVED***
	return s, end
***REMOVED***

func (c *channelMap) GetChannel(id int64) *ChannelMetric ***REMOVED***
	cm := &ChannelMetric***REMOVED******REMOVED***
	var cn *channel
	var ok bool
	c.mu.RLock()
	if cn, ok = c.channels[id]; !ok ***REMOVED***
		// channel with id doesn't exist.
		c.mu.RUnlock()
		return nil
	***REMOVED***
	cm.NestedChans = copyMap(cn.nestedChans)
	cm.SubChans = copyMap(cn.subChans)
	// cn.c can be set to &dummyChannel***REMOVED******REMOVED*** when deleteSelfFromMap is called. Save a copy of cn.c when
	// holding the lock to prevent potential data race.
	chanCopy := cn.c
	c.mu.RUnlock()
	cm.ChannelData = chanCopy.ChannelzMetric()
	cm.ID = cn.id
	cm.RefName = cn.refName
	cm.Trace = cn.trace.dumpData()
	return cm
***REMOVED***

func (c *channelMap) GetSubChannel(id int64) *SubChannelMetric ***REMOVED***
	cm := &SubChannelMetric***REMOVED******REMOVED***
	var sc *subChannel
	var ok bool
	c.mu.RLock()
	if sc, ok = c.subChannels[id]; !ok ***REMOVED***
		// subchannel with id doesn't exist.
		c.mu.RUnlock()
		return nil
	***REMOVED***
	cm.Sockets = copyMap(sc.sockets)
	// sc.c can be set to &dummyChannel***REMOVED******REMOVED*** when deleteSelfFromMap is called. Save a copy of sc.c when
	// holding the lock to prevent potential data race.
	chanCopy := sc.c
	c.mu.RUnlock()
	cm.ChannelData = chanCopy.ChannelzMetric()
	cm.ID = sc.id
	cm.RefName = sc.refName
	cm.Trace = sc.trace.dumpData()
	return cm
***REMOVED***

func (c *channelMap) GetSocket(id int64) *SocketMetric ***REMOVED***
	sm := &SocketMetric***REMOVED******REMOVED***
	c.mu.RLock()
	if ls, ok := c.listenSockets[id]; ok ***REMOVED***
		c.mu.RUnlock()
		sm.SocketData = ls.s.ChannelzMetric()
		sm.ID = ls.id
		sm.RefName = ls.refName
		return sm
	***REMOVED***
	if ns, ok := c.normalSockets[id]; ok ***REMOVED***
		c.mu.RUnlock()
		sm.SocketData = ns.s.ChannelzMetric()
		sm.ID = ns.id
		sm.RefName = ns.refName
		return sm
	***REMOVED***
	c.mu.RUnlock()
	return nil
***REMOVED***

func (c *channelMap) GetServer(id int64) *ServerMetric ***REMOVED***
	sm := &ServerMetric***REMOVED******REMOVED***
	var svr *server
	var ok bool
	c.mu.RLock()
	if svr, ok = c.servers[id]; !ok ***REMOVED***
		c.mu.RUnlock()
		return nil
	***REMOVED***
	sm.ListenSockets = copyMap(svr.listenSockets)
	c.mu.RUnlock()
	sm.ID = svr.id
	sm.RefName = svr.refName
	sm.ServerData = svr.s.ChannelzMetric()
	return sm
***REMOVED***

type idGenerator struct ***REMOVED***
	id int64
***REMOVED***

func (i *idGenerator) reset() ***REMOVED***
	atomic.StoreInt64(&i.id, 0)
***REMOVED***

func (i *idGenerator) genID() int64 ***REMOVED***
	return atomic.AddInt64(&i.id, 1)
***REMOVED***
