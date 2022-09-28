// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package topology contains types that handles the discovery, monitoring, and selection
// of servers. This package is designed to expose enough inner workings of service discovery
// and monitoring to allow low level applications to have fine grained control, while hiding
// most of the detailed implementation of the algorithms.
package topology // import "go.mongodb.org/mongo-driver/x/mongo/driver/topology"

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/randutil"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/dns"
)

// Topology state constants.
const (
	topologyDisconnected int64 = iota
	topologyDisconnecting
	topologyConnected
	topologyConnecting
)

// ErrSubscribeAfterClosed is returned when a user attempts to subscribe to a
// closed Server or Topology.
var ErrSubscribeAfterClosed = errors.New("cannot subscribe after closeConnection")

// ErrTopologyClosed is returned when a user attempts to call a method on a
// closed Topology.
var ErrTopologyClosed = errors.New("topology is closed")

// ErrTopologyConnected is returned whena  user attempts to Connect to an
// already connected Topology.
var ErrTopologyConnected = errors.New("topology is connected or connecting")

// ErrServerSelectionTimeout is returned from server selection when the server
// selection process took longer than allowed by the timeout.
var ErrServerSelectionTimeout = errors.New("server selection timeout")

// MonitorMode represents the way in which a server is monitored.
type MonitorMode uint8

// random is a package-global pseudo-random number generator.
var random = randutil.NewLockedRand()

// These constants are the available monitoring modes.
const (
	AutomaticMode MonitorMode = iota
	SingleMode
)

// Topology represents a MongoDB deployment.
type Topology struct ***REMOVED***
	state int64

	cfg *config

	desc atomic.Value // holds a description.Topology

	dnsResolver *dns.Resolver

	done chan struct***REMOVED******REMOVED***

	pollingRequired   bool
	pollingDone       chan struct***REMOVED******REMOVED***
	pollingwg         sync.WaitGroup
	rescanSRVInterval time.Duration
	pollHeartbeatTime atomic.Value // holds a bool

	updateCallback updateTopologyCallback
	fsm            *fsm

	// This should really be encapsulated into it's own type. This will likely
	// require a redesign so we can share a minimum of data between the
	// subscribers and the topology.
	subscribers         map[uint64]chan description.Topology
	currentSubscriberID uint64
	subscriptionsClosed bool
	subLock             sync.Mutex

	// We should redesign how we Connect and handle individal servers. This is
	// too difficult to maintain and it's rather easy to accidentally access
	// the servers without acquiring the lock or checking if the servers are
	// closed. This lock should also be an RWMutex.
	serversLock   sync.Mutex
	serversClosed bool
	servers       map[address.Address]*Server

	id primitive.ObjectID
***REMOVED***

var _ driver.Deployment = &Topology***REMOVED******REMOVED***
var _ driver.Subscriber = &Topology***REMOVED******REMOVED***

type serverSelectionState struct ***REMOVED***
	selector    description.ServerSelector
	timeoutChan <-chan time.Time
***REMOVED***

func newServerSelectionState(selector description.ServerSelector, timeoutChan <-chan time.Time) serverSelectionState ***REMOVED***
	return serverSelectionState***REMOVED***
		selector:    selector,
		timeoutChan: timeoutChan,
	***REMOVED***
***REMOVED***

// New creates a new topology.
func New(opts ...Option) (*Topology, error) ***REMOVED***
	cfg, err := newConfig(opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	t := &Topology***REMOVED***
		cfg:               cfg,
		done:              make(chan struct***REMOVED******REMOVED***),
		pollingDone:       make(chan struct***REMOVED******REMOVED***),
		rescanSRVInterval: 60 * time.Second,
		fsm:               newFSM(),
		subscribers:       make(map[uint64]chan description.Topology),
		servers:           make(map[address.Address]*Server),
		dnsResolver:       dns.DefaultResolver,
		id:                primitive.NewObjectID(),
	***REMOVED***
	t.desc.Store(description.Topology***REMOVED******REMOVED***)
	t.updateCallback = func(desc description.Server) description.Server ***REMOVED***
		return t.apply(context.TODO(), desc)
	***REMOVED***

	if t.cfg.uri != "" ***REMOVED***
		t.pollingRequired = strings.HasPrefix(t.cfg.uri, "mongodb+srv://") && !t.cfg.loadBalanced
	***REMOVED***

	t.publishTopologyOpeningEvent()

	return t, nil
***REMOVED***

// Connect initializes a Topology and starts the monitoring process. This function
// must be called to properly monitor the topology.
func (t *Topology) Connect() error ***REMOVED***
	if !atomic.CompareAndSwapInt64(&t.state, topologyDisconnected, topologyConnecting) ***REMOVED***
		return ErrTopologyConnected
	***REMOVED***

	t.desc.Store(description.Topology***REMOVED******REMOVED***)
	var err error
	t.serversLock.Lock()

	// A replica set name sets the initial topology type to ReplicaSetNoPrimary unless a direct connection is also
	// specified, in which case the initial type is Single.
	if t.cfg.replicaSetName != "" ***REMOVED***
		t.fsm.SetName = t.cfg.replicaSetName
		t.fsm.Kind = description.ReplicaSetNoPrimary
	***REMOVED***

	// A direct connection unconditionally sets the topology type to Single.
	if t.cfg.mode == SingleMode ***REMOVED***
		t.fsm.Kind = description.Single
	***REMOVED***

	for _, a := range t.cfg.seedList ***REMOVED***
		addr := address.Address(a).Canonicalize()
		t.fsm.Servers = append(t.fsm.Servers, description.NewDefaultServer(addr))
	***REMOVED***

	switch ***REMOVED***
	case t.cfg.loadBalanced:
		// In LoadBalanced mode, we mock a series of events: TopologyDescriptionChanged from Unknown to LoadBalanced,
		// ServerDescriptionChanged from Unknown to LoadBalancer, and then TopologyDescriptionChanged to reflect the
		// previous ServerDescriptionChanged event. We publish all of these events here because we don't start server
		// monitoring routines in this mode, so we have to mock state changes.

		// Transition from Unknown with no servers to LoadBalanced with a single Unknown server.
		t.fsm.Kind = description.LoadBalanced
		t.publishTopologyDescriptionChangedEvent(description.Topology***REMOVED******REMOVED***, t.fsm.Topology)

		addr := address.Address(t.cfg.seedList[0]).Canonicalize()
		if err := t.addServer(addr); err != nil ***REMOVED***
			t.serversLock.Unlock()
			return err
		***REMOVED***

		// Transition the server from Unknown to LoadBalancer.
		newServerDesc := t.servers[addr].Description()
		t.publishServerDescriptionChangedEvent(t.fsm.Servers[0], newServerDesc)

		// Transition from LoadBalanced with an Unknown server to LoadBalanced with a LoadBalancer.
		oldDesc := t.fsm.Topology
		t.fsm.Servers = []description.Server***REMOVED***newServerDesc***REMOVED***
		t.desc.Store(t.fsm.Topology)
		t.publishTopologyDescriptionChangedEvent(oldDesc, t.fsm.Topology)
	default:
		// In non-LB mode, we only publish an initial TopologyDescriptionChanged event from Unknown with no servers to
		// the current state (e.g. Unknown with one or more servers if we're discovering or Single with one server if
		// we're connecting directly). Other events are published when state changes occur due to responses in the
		// server monitoring goroutines.

		newDesc := description.Topology***REMOVED***
			Kind:                  t.fsm.Kind,
			Servers:               t.fsm.Servers,
			SessionTimeoutMinutes: t.fsm.SessionTimeoutMinutes,
		***REMOVED***
		t.desc.Store(newDesc)
		t.publishTopologyDescriptionChangedEvent(description.Topology***REMOVED******REMOVED***, t.fsm.Topology)
		for _, a := range t.cfg.seedList ***REMOVED***
			addr := address.Address(a).Canonicalize()
			err = t.addServer(addr)
			if err != nil ***REMOVED***
				t.serversLock.Unlock()
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	t.serversLock.Unlock()
	if t.pollingRequired ***REMOVED***
		go t.pollSRVRecords()
		t.pollingwg.Add(1)
	***REMOVED***

	t.subscriptionsClosed = false // explicitly set in case topology was disconnected and then reconnected

	atomic.StoreInt64(&t.state, topologyConnected)
	return nil
***REMOVED***

// Disconnect closes the topology. It stops the monitoring thread and
// closes all open subscriptions.
func (t *Topology) Disconnect(ctx context.Context) error ***REMOVED***
	if !atomic.CompareAndSwapInt64(&t.state, topologyConnected, topologyDisconnecting) ***REMOVED***
		return ErrTopologyClosed
	***REMOVED***

	servers := make(map[address.Address]*Server)
	t.serversLock.Lock()
	t.serversClosed = true
	for addr, server := range t.servers ***REMOVED***
		servers[addr] = server
	***REMOVED***
	t.serversLock.Unlock()

	for _, server := range servers ***REMOVED***
		_ = server.Disconnect(ctx)
		t.publishServerClosedEvent(server.address)
	***REMOVED***

	t.subLock.Lock()
	for id, ch := range t.subscribers ***REMOVED***
		close(ch)
		delete(t.subscribers, id)
	***REMOVED***
	t.subscriptionsClosed = true
	t.subLock.Unlock()

	if t.pollingRequired ***REMOVED***
		t.pollingDone <- struct***REMOVED******REMOVED******REMOVED******REMOVED***
		t.pollingwg.Wait()
	***REMOVED***

	t.desc.Store(description.Topology***REMOVED******REMOVED***)

	atomic.StoreInt64(&t.state, topologyDisconnected)
	t.publishTopologyClosedEvent()
	return nil
***REMOVED***

// Description returns a description of the topology.
func (t *Topology) Description() description.Topology ***REMOVED***
	td, ok := t.desc.Load().(description.Topology)
	if !ok ***REMOVED***
		td = description.Topology***REMOVED******REMOVED***
	***REMOVED***
	return td
***REMOVED***

// Kind returns the topology kind of this Topology.
func (t *Topology) Kind() description.TopologyKind ***REMOVED*** return t.Description().Kind ***REMOVED***

// Subscribe returns a Subscription on which all updated description.Topologys
// will be sent. The channel of the subscription will have a buffer size of one,
// and will be pre-populated with the current description.Topology.
// Subscribe implements the driver.Subscriber interface.
func (t *Topology) Subscribe() (*driver.Subscription, error) ***REMOVED***
	if atomic.LoadInt64(&t.state) != topologyConnected ***REMOVED***
		return nil, errors.New("cannot subscribe to Topology that is not connected")
	***REMOVED***
	ch := make(chan description.Topology, 1)
	td, ok := t.desc.Load().(description.Topology)
	if !ok ***REMOVED***
		td = description.Topology***REMOVED******REMOVED***
	***REMOVED***
	ch <- td

	t.subLock.Lock()
	defer t.subLock.Unlock()
	if t.subscriptionsClosed ***REMOVED***
		return nil, ErrSubscribeAfterClosed
	***REMOVED***
	id := t.currentSubscriberID
	t.subscribers[id] = ch
	t.currentSubscriberID++

	return &driver.Subscription***REMOVED***
		Updates: ch,
		ID:      id,
	***REMOVED***, nil
***REMOVED***

// Unsubscribe unsubscribes the given subscription from the topology and closes the subscription channel.
// Unsubscribe implements the driver.Subscriber interface.
func (t *Topology) Unsubscribe(sub *driver.Subscription) error ***REMOVED***
	t.subLock.Lock()
	defer t.subLock.Unlock()

	if t.subscriptionsClosed ***REMOVED***
		return nil
	***REMOVED***

	ch, ok := t.subscribers[sub.ID]
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	close(ch)
	delete(t.subscribers, sub.ID)
	return nil
***REMOVED***

// RequestImmediateCheck will send heartbeats to all the servers in the
// topology right away, instead of waiting for the heartbeat timeout.
func (t *Topology) RequestImmediateCheck() ***REMOVED***
	if atomic.LoadInt64(&t.state) != topologyConnected ***REMOVED***
		return
	***REMOVED***
	t.serversLock.Lock()
	for _, server := range t.servers ***REMOVED***
		server.RequestImmediateCheck()
	***REMOVED***
	t.serversLock.Unlock()
***REMOVED***

// SelectServer selects a server with given a selector. SelectServer complies with the
// server selection spec, and will time out after serverSelectionTimeout or when the
// parent context is done.
func (t *Topology) SelectServer(ctx context.Context, ss description.ServerSelector) (driver.Server, error) ***REMOVED***
	if atomic.LoadInt64(&t.state) != topologyConnected ***REMOVED***
		return nil, ErrTopologyClosed
	***REMOVED***
	var ssTimeoutCh <-chan time.Time

	if t.cfg.serverSelectionTimeout > 0 ***REMOVED***
		ssTimeout := time.NewTimer(t.cfg.serverSelectionTimeout)
		ssTimeoutCh = ssTimeout.C
		defer ssTimeout.Stop()
	***REMOVED***

	var doneOnce bool
	var sub *driver.Subscription
	selectionState := newServerSelectionState(ss, ssTimeoutCh)
	for ***REMOVED***
		var suitable []description.Server
		var selectErr error

		if !doneOnce ***REMOVED***
			// for the first pass, select a server from the current description.
			// this improves selection speed for up-to-date topology descriptions.
			suitable, selectErr = t.selectServerFromDescription(t.Description(), selectionState)
			doneOnce = true
		***REMOVED*** else ***REMOVED***
			// if the first pass didn't select a server, the previous description did not contain a suitable server, so
			// we subscribe to the topology and attempt to obtain a server from that subscription
			if sub == nil ***REMOVED***
				var err error
				sub, err = t.Subscribe()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				defer t.Unsubscribe(sub)
			***REMOVED***

			suitable, selectErr = t.selectServerFromSubscription(ctx, sub.Updates, selectionState)
		***REMOVED***
		if selectErr != nil ***REMOVED***
			return nil, selectErr
		***REMOVED***

		if len(suitable) == 0 ***REMOVED***
			// try again if there are no servers available
			continue
		***REMOVED***

		// If there's only one suitable server description, try to find the associated server and
		// return it. This is an optimization primarily for standalone and load-balanced deployments.
		if len(suitable) == 1 ***REMOVED***
			server, err := t.FindServer(suitable[0])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if server == nil ***REMOVED***
				continue
			***REMOVED***
			return server, nil
		***REMOVED***

		// Randomly select 2 suitable server descriptions and find servers for them. We select two
		// so we can pick the one with the one with fewer in-progress operations below.
		desc1, desc2 := pick2(suitable)
		server1, err := t.FindServer(desc1)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		server2, err := t.FindServer(desc2)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// If we don't have an actual server for one or both of the provided descriptions, either
		// return the one server we have, or try again if they're both nil. This could happen for a
		// number of reasons, including that the server has since stopped being a part of this
		// topology.
		if server1 == nil || server2 == nil ***REMOVED***
			if server1 == nil && server2 == nil ***REMOVED***
				continue
			***REMOVED***
			if server1 != nil ***REMOVED***
				return server1, nil
			***REMOVED***
			return server2, nil
		***REMOVED***

		// Of the two randomly selected suitable servers, pick the one with fewer in-use connections.
		// We use in-use connections as an analog for in-progress operations because they are almost
		// always the same value for a given server.
		if server1.OperationCount() < server2.OperationCount() ***REMOVED***
			return server1, nil
		***REMOVED***
		return server2, nil
	***REMOVED***
***REMOVED***

// pick2 returns 2 random server descriptions from the input slice of server descriptions,
// guaranteeing that the same element from the slice is not picked twice. The order of server
// descriptions in the input slice may be modified. If fewer than 2 server descriptions are
// provided, pick2 will panic.
func pick2(ds []description.Server) (description.Server, description.Server) ***REMOVED***
	// Select a random index from the input slice and keep the server description from that index.
	idx := random.Intn(len(ds))
	s1 := ds[idx]

	// Swap the selected index to the end and reslice to remove it so we don't pick the same server
	// description twice.
	ds[idx], ds[len(ds)-1] = ds[len(ds)-1], ds[idx]
	ds = ds[:len(ds)-1]

	// Select another random index from the input slice and return both selected server descriptions.
	return s1, ds[random.Intn(len(ds))]
***REMOVED***

// FindServer will attempt to find a server that fits the given server description.
// This method will return nil, nil if a matching server could not be found.
func (t *Topology) FindServer(selected description.Server) (*SelectedServer, error) ***REMOVED***
	if atomic.LoadInt64(&t.state) != topologyConnected ***REMOVED***
		return nil, ErrTopologyClosed
	***REMOVED***
	t.serversLock.Lock()
	defer t.serversLock.Unlock()
	server, ok := t.servers[selected.Addr]
	if !ok ***REMOVED***
		return nil, nil
	***REMOVED***

	desc := t.Description()
	return &SelectedServer***REMOVED***
		Server: server,
		Kind:   desc.Kind,
	***REMOVED***, nil
***REMOVED***

// selectServerFromSubscription loops until a topology description is available for server selection. It returns
// when the given context expires, server selection timeout is reached, or a description containing a selectable
// server is available.
func (t *Topology) selectServerFromSubscription(ctx context.Context, subscriptionCh <-chan description.Topology,
	selectionState serverSelectionState) ([]description.Server, error) ***REMOVED***

	current := t.Description()
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil, ServerSelectionError***REMOVED***Wrapped: ctx.Err(), Desc: current***REMOVED***
		case <-selectionState.timeoutChan:
			return nil, ServerSelectionError***REMOVED***Wrapped: ErrServerSelectionTimeout, Desc: current***REMOVED***
		case current = <-subscriptionCh:
		***REMOVED***

		suitable, err := t.selectServerFromDescription(current, selectionState)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if len(suitable) > 0 ***REMOVED***
			return suitable, nil
		***REMOVED***
		t.RequestImmediateCheck()
	***REMOVED***
***REMOVED***

// selectServerFromDescription process the given topology description and returns a slice of suitable servers.
func (t *Topology) selectServerFromDescription(desc description.Topology,
	selectionState serverSelectionState) ([]description.Server, error) ***REMOVED***

	// Unlike selectServerFromSubscription, this code path does not check ctx.Done or selectionState.timeoutChan because
	// selecting a server from a description is not a blocking operation.

	if desc.CompatibilityErr != nil ***REMOVED***
		return nil, desc.CompatibilityErr
	***REMOVED***

	// If the topology kind is LoadBalanced, the LB is the only server and it is always considered selectable. The
	// selectors exported by the driver should already return the LB as a candidate, so this but this check ensures that
	// the LB is always selectable even if a user of the low-level driver provides a custom selector.
	if desc.Kind == description.LoadBalanced ***REMOVED***
		return desc.Servers, nil
	***REMOVED***

	var allowed []description.Server
	for _, s := range desc.Servers ***REMOVED***
		if s.Kind != description.Unknown ***REMOVED***
			allowed = append(allowed, s)
		***REMOVED***
	***REMOVED***

	suitable, err := selectionState.selector.SelectServer(desc, allowed)
	if err != nil ***REMOVED***
		return nil, ServerSelectionError***REMOVED***Wrapped: err, Desc: desc***REMOVED***
	***REMOVED***
	return suitable, nil
***REMOVED***

func (t *Topology) pollSRVRecords() ***REMOVED***
	defer t.pollingwg.Done()

	serverConfig := newServerConfig(t.cfg.serverOpts...)
	heartbeatInterval := serverConfig.heartbeatInterval

	pollTicker := time.NewTicker(t.rescanSRVInterval)
	defer pollTicker.Stop()
	t.pollHeartbeatTime.Store(false)
	var doneOnce bool
	defer func() ***REMOVED***
		//  ¯\_(ツ)_/¯
		if r := recover(); r != nil && !doneOnce ***REMOVED***
			<-t.pollingDone
		***REMOVED***
	***REMOVED***()

	// remove the scheme
	uri := t.cfg.uri[14:]
	hosts := uri
	if idx := strings.IndexAny(uri, "/?@"); idx != -1 ***REMOVED***
		hosts = uri[:idx]
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case <-pollTicker.C:
		case <-t.pollingDone:
			doneOnce = true
			return
		***REMOVED***
		topoKind := t.Description().Kind
		if !(topoKind == description.Unknown || topoKind == description.Sharded) ***REMOVED***
			break
		***REMOVED***

		parsedHosts, err := t.dnsResolver.ParseHosts(hosts, t.cfg.srvServiceName, false)
		// DNS problem or no verified hosts returned
		if err != nil || len(parsedHosts) == 0 ***REMOVED***
			if !t.pollHeartbeatTime.Load().(bool) ***REMOVED***
				pollTicker.Stop()
				pollTicker = time.NewTicker(heartbeatInterval)
				t.pollHeartbeatTime.Store(true)
			***REMOVED***
			continue
		***REMOVED***
		if t.pollHeartbeatTime.Load().(bool) ***REMOVED***
			pollTicker.Stop()
			pollTicker = time.NewTicker(t.rescanSRVInterval)
			t.pollHeartbeatTime.Store(false)
		***REMOVED***

		cont := t.processSRVResults(parsedHosts)
		if !cont ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	<-t.pollingDone
	doneOnce = true
***REMOVED***

func (t *Topology) processSRVResults(parsedHosts []string) bool ***REMOVED***
	t.serversLock.Lock()
	defer t.serversLock.Unlock()

	if t.serversClosed ***REMOVED***
		return false
	***REMOVED***
	prev := t.fsm.Topology
	diff := diffHostList(t.fsm.Topology, parsedHosts)

	if len(diff.Added) == 0 && len(diff.Removed) == 0 ***REMOVED***
		return true
	***REMOVED***

	for _, r := range diff.Removed ***REMOVED***
		addr := address.Address(r).Canonicalize()
		s, ok := t.servers[addr]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		go func() ***REMOVED***
			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()
			_ = s.Disconnect(cancelCtx)
		***REMOVED***()
		delete(t.servers, addr)
		t.fsm.removeServerByAddr(addr)
		t.publishServerClosedEvent(s.address)
	***REMOVED***

	// Now that we've removed all the hosts that disappeared from the SRV record, we need to add any
	// new hosts added to the SRV record. If adding all of the new hosts would increase the number
	// of servers past srvMaxHosts, shuffle the list of added hosts.
	if t.cfg.srvMaxHosts > 0 && len(t.servers)+len(diff.Added) > t.cfg.srvMaxHosts ***REMOVED***
		random.Shuffle(len(diff.Added), func(i, j int) ***REMOVED***
			diff.Added[i], diff.Added[j] = diff.Added[j], diff.Added[i]
		***REMOVED***)
	***REMOVED***
	// Add all added hosts until the number of servers reaches srvMaxHosts.
	for _, a := range diff.Added ***REMOVED***
		if t.cfg.srvMaxHosts > 0 && len(t.servers) >= t.cfg.srvMaxHosts ***REMOVED***
			break
		***REMOVED***
		addr := address.Address(a).Canonicalize()
		_ = t.addServer(addr)
		t.fsm.addServer(addr)
	***REMOVED***

	//store new description
	newDesc := description.Topology***REMOVED***
		Kind:                  t.fsm.Kind,
		Servers:               t.fsm.Servers,
		SessionTimeoutMinutes: t.fsm.SessionTimeoutMinutes,
	***REMOVED***
	t.desc.Store(newDesc)

	if !prev.Equal(newDesc) ***REMOVED***
		t.publishTopologyDescriptionChangedEvent(prev, newDesc)
	***REMOVED***

	t.subLock.Lock()
	for _, ch := range t.subscribers ***REMOVED***
		// We drain the description if there's one in the channel
		select ***REMOVED***
		case <-ch:
		default:
		***REMOVED***
		ch <- newDesc
	***REMOVED***
	t.subLock.Unlock()

	return true
***REMOVED***

// apply updates the Topology and its underlying FSM based on the provided server description and returns the server
// description that should be stored.
func (t *Topology) apply(ctx context.Context, desc description.Server) description.Server ***REMOVED***
	t.serversLock.Lock()
	defer t.serversLock.Unlock()

	ind, ok := t.fsm.findServer(desc.Addr)
	if t.serversClosed || !ok ***REMOVED***
		return desc
	***REMOVED***

	prev := t.fsm.Topology
	oldDesc := t.fsm.Servers[ind]
	if oldDesc.TopologyVersion.CompareToIncoming(desc.TopologyVersion) > 0 ***REMOVED***
		return oldDesc
	***REMOVED***

	var current description.Topology
	current, desc = t.fsm.apply(desc)

	if !oldDesc.Equal(desc) ***REMOVED***
		t.publishServerDescriptionChangedEvent(oldDesc, desc)
	***REMOVED***

	diff := diffTopology(prev, current)

	for _, removed := range diff.Removed ***REMOVED***
		if s, ok := t.servers[removed.Addr]; ok ***REMOVED***
			go func() ***REMOVED***
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				_ = s.Disconnect(cancelCtx)
			***REMOVED***()
			delete(t.servers, removed.Addr)
			t.publishServerClosedEvent(s.address)
		***REMOVED***
	***REMOVED***

	for _, added := range diff.Added ***REMOVED***
		_ = t.addServer(added.Addr)
	***REMOVED***

	t.desc.Store(current)
	if !prev.Equal(current) ***REMOVED***
		t.publishTopologyDescriptionChangedEvent(prev, current)
	***REMOVED***

	t.subLock.Lock()
	for _, ch := range t.subscribers ***REMOVED***
		// We drain the description if there's one in the channel
		select ***REMOVED***
		case <-ch:
		default:
		***REMOVED***
		ch <- current
	***REMOVED***
	t.subLock.Unlock()

	return desc
***REMOVED***

func (t *Topology) addServer(addr address.Address) error ***REMOVED***
	if _, ok := t.servers[addr]; ok ***REMOVED***
		return nil
	***REMOVED***

	svr, err := ConnectServer(addr, t.updateCallback, t.id, t.cfg.serverOpts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	t.servers[addr] = svr

	return nil
***REMOVED***

// String implements the Stringer interface
func (t *Topology) String() string ***REMOVED***
	desc := t.Description()

	serversStr := ""
	t.serversLock.Lock()
	defer t.serversLock.Unlock()
	for _, s := range t.servers ***REMOVED***
		serversStr += "***REMOVED*** " + s.String() + " ***REMOVED***, "
	***REMOVED***
	return fmt.Sprintf("Type: %s, Servers: [%s]", desc.Kind, serversStr)
***REMOVED***

// publishes a ServerDescriptionChangedEvent to indicate the server description has changed
func (t *Topology) publishServerDescriptionChangedEvent(prev description.Server, current description.Server) ***REMOVED***
	serverDescriptionChanged := &event.ServerDescriptionChangedEvent***REMOVED***
		Address:             current.Addr,
		TopologyID:          t.id,
		PreviousDescription: prev,
		NewDescription:      current,
	***REMOVED***

	if t.cfg.serverMonitor != nil && t.cfg.serverMonitor.ServerDescriptionChanged != nil ***REMOVED***
		t.cfg.serverMonitor.ServerDescriptionChanged(serverDescriptionChanged)
	***REMOVED***
***REMOVED***

// publishes a ServerClosedEvent to indicate the server has closed
func (t *Topology) publishServerClosedEvent(addr address.Address) ***REMOVED***
	serverClosed := &event.ServerClosedEvent***REMOVED***
		Address:    addr,
		TopologyID: t.id,
	***REMOVED***

	if t.cfg.serverMonitor != nil && t.cfg.serverMonitor.ServerClosed != nil ***REMOVED***
		t.cfg.serverMonitor.ServerClosed(serverClosed)
	***REMOVED***
***REMOVED***

// publishes a TopologyDescriptionChangedEvent to indicate the topology description has changed
func (t *Topology) publishTopologyDescriptionChangedEvent(prev description.Topology, current description.Topology) ***REMOVED***
	topologyDescriptionChanged := &event.TopologyDescriptionChangedEvent***REMOVED***
		TopologyID:          t.id,
		PreviousDescription: prev,
		NewDescription:      current,
	***REMOVED***

	if t.cfg.serverMonitor != nil && t.cfg.serverMonitor.TopologyDescriptionChanged != nil ***REMOVED***
		t.cfg.serverMonitor.TopologyDescriptionChanged(topologyDescriptionChanged)
	***REMOVED***
***REMOVED***

// publishes a TopologyOpeningEvent to indicate the topology is being initialized
func (t *Topology) publishTopologyOpeningEvent() ***REMOVED***
	topologyOpening := &event.TopologyOpeningEvent***REMOVED***
		TopologyID: t.id,
	***REMOVED***

	if t.cfg.serverMonitor != nil && t.cfg.serverMonitor.TopologyOpening != nil ***REMOVED***
		t.cfg.serverMonitor.TopologyOpening(topologyOpening)
	***REMOVED***
***REMOVED***

// publishes a TopologyClosedEvent to indicate the topology has been closed
func (t *Topology) publishTopologyClosedEvent() ***REMOVED***
	topologyClosed := &event.TopologyClosedEvent***REMOVED***
		TopologyID: t.id,
	***REMOVED***

	if t.cfg.serverMonitor != nil && t.cfg.serverMonitor.TopologyClosed != nil ***REMOVED***
		t.cfg.serverMonitor.TopologyClosed(topologyClosed)
	***REMOVED***
***REMOVED***
