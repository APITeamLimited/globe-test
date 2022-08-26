package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/redis/v9/internal"
	"github.com/APITeamLimited/redis/v9/internal/hashtag"
	"github.com/APITeamLimited/redis/v9/internal/pool"
	"github.com/APITeamLimited/redis/v9/internal/proto"
	"github.com/APITeamLimited/redis/v9/internal/rand"
)

var errClusterNoNodes = fmt.Errorf("redis: cluster has no nodes")

// ClusterOptions are used to configure a cluster client and should be
// passed to NewClusterClient.
type ClusterOptions struct ***REMOVED***
	// A seed list of host:port addresses of cluster nodes.
	Addrs []string

	// NewClient creates a cluster node client with provided name and options.
	NewClient func(opt *Options) *Client

	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 3 retries.
	MaxRedirects int

	// Enables read-only commands on slave nodes.
	ReadOnly bool
	// Allows routing read-only commands to the closest master or slave node.
	// It automatically enables ReadOnly.
	RouteByLatency bool
	// Allows routing read-only commands to the random master or slave node.
	// It automatically enables ReadOnly.
	RouteRandomly bool

	// Optional function that returns cluster slots information.
	// It is useful to manually create cluster of standalone Redis servers
	// and load-balance read/write operations between master and slaves.
	// It can use service like ZooKeeper to maintain configuration information
	// and Cluster.ReloadState to manually trigger state reloading.
	ClusterSlots func(context.Context) ([]ClusterSlot, error)

	// Following options are copied from Options struct.

	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	OnConnect func(ctx context.Context, cn *Conn) error

	Username string
	Password string

	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolFIFO bool

	// PoolSize applies per cluster node and not for the whole cluster.
	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration

	TLSConfig *tls.Config
***REMOVED***

func (opt *ClusterOptions) init() ***REMOVED***
	if opt.MaxRedirects == -1 ***REMOVED***
		opt.MaxRedirects = 0
	***REMOVED*** else if opt.MaxRedirects == 0 ***REMOVED***
		opt.MaxRedirects = 3
	***REMOVED***

	if opt.RouteByLatency || opt.RouteRandomly ***REMOVED***
		opt.ReadOnly = true
	***REMOVED***

	if opt.PoolSize == 0 ***REMOVED***
		opt.PoolSize = 5 * runtime.GOMAXPROCS(0)
	***REMOVED***

	switch opt.ReadTimeout ***REMOVED***
	case -1:
		opt.ReadTimeout = 0
	case 0:
		opt.ReadTimeout = 3 * time.Second
	***REMOVED***
	switch opt.WriteTimeout ***REMOVED***
	case -1:
		opt.WriteTimeout = 0
	case 0:
		opt.WriteTimeout = opt.ReadTimeout
	***REMOVED***

	if opt.MaxRetries == 0 ***REMOVED***
		opt.MaxRetries = -1
	***REMOVED***
	switch opt.MinRetryBackoff ***REMOVED***
	case -1:
		opt.MinRetryBackoff = 0
	case 0:
		opt.MinRetryBackoff = 8 * time.Millisecond
	***REMOVED***
	switch opt.MaxRetryBackoff ***REMOVED***
	case -1:
		opt.MaxRetryBackoff = 0
	case 0:
		opt.MaxRetryBackoff = 512 * time.Millisecond
	***REMOVED***

	if opt.NewClient == nil ***REMOVED***
		opt.NewClient = NewClient
	***REMOVED***
***REMOVED***

func (opt *ClusterOptions) clientOptions() *Options ***REMOVED***
	return &Options***REMOVED***
		Dialer:    opt.Dialer,
		OnConnect: opt.OnConnect,

		Username: opt.Username,
		Password: opt.Password,

		MaxRetries:      opt.MaxRetries,
		MinRetryBackoff: opt.MinRetryBackoff,
		MaxRetryBackoff: opt.MaxRetryBackoff,

		DialTimeout:  opt.DialTimeout,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,

		PoolFIFO:        opt.PoolFIFO,
		PoolSize:        opt.PoolSize,
		PoolTimeout:     opt.PoolTimeout,
		MinIdleConns:    opt.MinIdleConns,
		MaxIdleConns:    opt.MaxIdleConns,
		ConnMaxIdleTime: opt.ConnMaxIdleTime,
		ConnMaxLifetime: opt.ConnMaxLifetime,

		TLSConfig: opt.TLSConfig,
		// If ClusterSlots is populated, then we probably have an artificial
		// cluster whose nodes are not in clustering mode (otherwise there isn't
		// much use for ClusterSlots config).  This means we cannot execute the
		// READONLY command against that node -- setting readOnly to false in such
		// situations in the options below will prevent that from happening.
		readOnly: opt.ReadOnly && opt.ClusterSlots == nil,
	***REMOVED***
***REMOVED***

//------------------------------------------------------------------------------

type clusterNode struct ***REMOVED***
	Client *Client

	latency    uint32 // atomic
	generation uint32 // atomic
	failing    uint32 // atomic
***REMOVED***

func newClusterNode(clOpt *ClusterOptions, addr string) *clusterNode ***REMOVED***
	opt := clOpt.clientOptions()
	opt.Addr = addr
	node := clusterNode***REMOVED***
		Client: clOpt.NewClient(opt),
	***REMOVED***

	node.latency = math.MaxUint32
	if clOpt.RouteByLatency ***REMOVED***
		go node.updateLatency()
	***REMOVED***

	return &node
***REMOVED***

func (n *clusterNode) String() string ***REMOVED***
	return n.Client.String()
***REMOVED***

func (n *clusterNode) Close() error ***REMOVED***
	return n.Client.Close()
***REMOVED***

func (n *clusterNode) updateLatency() ***REMOVED***
	const numProbe = 10
	var dur uint64

	successes := 0
	for i := 0; i < numProbe; i++ ***REMOVED***
		time.Sleep(time.Duration(10+rand.Intn(10)) * time.Millisecond)

		start := time.Now()
		err := n.Client.Ping(context.TODO()).Err()
		if err == nil ***REMOVED***
			dur += uint64(time.Since(start) / time.Microsecond)
			successes++
		***REMOVED***
	***REMOVED***

	var latency float64
	if successes == 0 ***REMOVED***
		// If none of the pings worked, set latency to some arbitrarily high value so this node gets
		// least priority.
		latency = float64((1 * time.Minute) / time.Microsecond)
	***REMOVED*** else ***REMOVED***
		latency = float64(dur) / float64(successes)
	***REMOVED***
	atomic.StoreUint32(&n.latency, uint32(latency+0.5))
***REMOVED***

func (n *clusterNode) Latency() time.Duration ***REMOVED***
	latency := atomic.LoadUint32(&n.latency)
	return time.Duration(latency) * time.Microsecond
***REMOVED***

func (n *clusterNode) MarkAsFailing() ***REMOVED***
	atomic.StoreUint32(&n.failing, uint32(time.Now().Unix()))
***REMOVED***

func (n *clusterNode) Failing() bool ***REMOVED***
	const timeout = 15 // 15 seconds

	failing := atomic.LoadUint32(&n.failing)
	if failing == 0 ***REMOVED***
		return false
	***REMOVED***
	if time.Now().Unix()-int64(failing) < timeout ***REMOVED***
		return true
	***REMOVED***
	atomic.StoreUint32(&n.failing, 0)
	return false
***REMOVED***

func (n *clusterNode) Generation() uint32 ***REMOVED***
	return atomic.LoadUint32(&n.generation)
***REMOVED***

func (n *clusterNode) SetGeneration(gen uint32) ***REMOVED***
	for ***REMOVED***
		v := atomic.LoadUint32(&n.generation)
		if gen < v || atomic.CompareAndSwapUint32(&n.generation, v, gen) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

//------------------------------------------------------------------------------

type clusterNodes struct ***REMOVED***
	opt *ClusterOptions

	mu          sync.RWMutex
	addrs       []string
	nodes       map[string]*clusterNode
	activeAddrs []string
	closed      bool

	_generation uint32 // atomic
***REMOVED***

func newClusterNodes(opt *ClusterOptions) *clusterNodes ***REMOVED***
	return &clusterNodes***REMOVED***
		opt: opt,

		addrs: opt.Addrs,
		nodes: make(map[string]*clusterNode),
	***REMOVED***
***REMOVED***

func (c *clusterNodes) Close() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed ***REMOVED***
		return nil
	***REMOVED***
	c.closed = true

	var firstErr error
	for _, node := range c.nodes ***REMOVED***
		if err := node.Client.Close(); err != nil && firstErr == nil ***REMOVED***
			firstErr = err
		***REMOVED***
	***REMOVED***

	c.nodes = nil
	c.activeAddrs = nil

	return firstErr
***REMOVED***

func (c *clusterNodes) Addrs() ([]string, error) ***REMOVED***
	var addrs []string

	c.mu.RLock()
	closed := c.closed //nolint:ifshort
	if !closed ***REMOVED***
		if len(c.activeAddrs) > 0 ***REMOVED***
			addrs = c.activeAddrs
		***REMOVED*** else ***REMOVED***
			addrs = c.addrs
		***REMOVED***
	***REMOVED***
	c.mu.RUnlock()

	if closed ***REMOVED***
		return nil, pool.ErrClosed
	***REMOVED***
	if len(addrs) == 0 ***REMOVED***
		return nil, errClusterNoNodes
	***REMOVED***
	return addrs, nil
***REMOVED***

func (c *clusterNodes) NextGeneration() uint32 ***REMOVED***
	return atomic.AddUint32(&c._generation, 1)
***REMOVED***

// GC removes unused nodes.
func (c *clusterNodes) GC(generation uint32) ***REMOVED***
	//nolint:prealloc
	var collected []*clusterNode

	c.mu.Lock()

	c.activeAddrs = c.activeAddrs[:0]
	for addr, node := range c.nodes ***REMOVED***
		if node.Generation() >= generation ***REMOVED***
			c.activeAddrs = append(c.activeAddrs, addr)
			if c.opt.RouteByLatency ***REMOVED***
				go node.updateLatency()
			***REMOVED***
			continue
		***REMOVED***

		delete(c.nodes, addr)
		collected = append(collected, node)
	***REMOVED***

	c.mu.Unlock()

	for _, node := range collected ***REMOVED***
		_ = node.Client.Close()
	***REMOVED***
***REMOVED***

func (c *clusterNodes) GetOrCreate(addr string) (*clusterNode, error) ***REMOVED***
	node, err := c.get(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if node != nil ***REMOVED***
		return node, nil
	***REMOVED***

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed ***REMOVED***
		return nil, pool.ErrClosed
	***REMOVED***

	node, ok := c.nodes[addr]
	if ok ***REMOVED***
		return node, nil
	***REMOVED***

	node = newClusterNode(c.opt, addr)

	c.addrs = appendIfNotExists(c.addrs, addr)
	c.nodes[addr] = node

	return node, nil
***REMOVED***

func (c *clusterNodes) get(addr string) (*clusterNode, error) ***REMOVED***
	var node *clusterNode
	var err error
	c.mu.RLock()
	if c.closed ***REMOVED***
		err = pool.ErrClosed
	***REMOVED*** else ***REMOVED***
		node = c.nodes[addr]
	***REMOVED***
	c.mu.RUnlock()
	return node, err
***REMOVED***

func (c *clusterNodes) All() ([]*clusterNode, error) ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed ***REMOVED***
		return nil, pool.ErrClosed
	***REMOVED***

	cp := make([]*clusterNode, 0, len(c.nodes))
	for _, node := range c.nodes ***REMOVED***
		cp = append(cp, node)
	***REMOVED***
	return cp, nil
***REMOVED***

func (c *clusterNodes) Random() (*clusterNode, error) ***REMOVED***
	addrs, err := c.Addrs()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	n := rand.Intn(len(addrs))
	return c.GetOrCreate(addrs[n])
***REMOVED***

//------------------------------------------------------------------------------

type clusterSlot struct ***REMOVED***
	start, end int
	nodes      []*clusterNode
***REMOVED***

type clusterSlotSlice []*clusterSlot

func (p clusterSlotSlice) Len() int ***REMOVED***
	return len(p)
***REMOVED***

func (p clusterSlotSlice) Less(i, j int) bool ***REMOVED***
	return p[i].start < p[j].start
***REMOVED***

func (p clusterSlotSlice) Swap(i, j int) ***REMOVED***
	p[i], p[j] = p[j], p[i]
***REMOVED***

type clusterState struct ***REMOVED***
	nodes   *clusterNodes
	Masters []*clusterNode
	Slaves  []*clusterNode

	slots []*clusterSlot

	generation uint32
	createdAt  time.Time
***REMOVED***

func newClusterState(
	nodes *clusterNodes, slots []ClusterSlot, origin string,
) (*clusterState, error) ***REMOVED***
	c := clusterState***REMOVED***
		nodes: nodes,

		slots: make([]*clusterSlot, 0, len(slots)),

		generation: nodes.NextGeneration(),
		createdAt:  time.Now(),
	***REMOVED***

	originHost, _, _ := net.SplitHostPort(origin)
	isLoopbackOrigin := isLoopback(originHost)

	for _, slot := range slots ***REMOVED***
		var nodes []*clusterNode
		for i, slotNode := range slot.Nodes ***REMOVED***
			addr := slotNode.Addr
			if !isLoopbackOrigin ***REMOVED***
				addr = replaceLoopbackHost(addr, originHost)
			***REMOVED***

			node, err := c.nodes.GetOrCreate(addr)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			node.SetGeneration(c.generation)
			nodes = append(nodes, node)

			if i == 0 ***REMOVED***
				c.Masters = appendUniqueNode(c.Masters, node)
			***REMOVED*** else ***REMOVED***
				c.Slaves = appendUniqueNode(c.Slaves, node)
			***REMOVED***
		***REMOVED***

		c.slots = append(c.slots, &clusterSlot***REMOVED***
			start: slot.Start,
			end:   slot.End,
			nodes: nodes,
		***REMOVED***)
	***REMOVED***

	sort.Sort(clusterSlotSlice(c.slots))

	time.AfterFunc(time.Minute, func() ***REMOVED***
		nodes.GC(c.generation)
	***REMOVED***)

	return &c, nil
***REMOVED***

func replaceLoopbackHost(nodeAddr, originHost string) string ***REMOVED***
	nodeHost, nodePort, err := net.SplitHostPort(nodeAddr)
	if err != nil ***REMOVED***
		return nodeAddr
	***REMOVED***

	nodeIP := net.ParseIP(nodeHost)
	if nodeIP == nil ***REMOVED***
		return nodeAddr
	***REMOVED***

	if !nodeIP.IsLoopback() ***REMOVED***
		return nodeAddr
	***REMOVED***

	// Use origin host which is not loopback and node port.
	return net.JoinHostPort(originHost, nodePort)
***REMOVED***

func isLoopback(host string) bool ***REMOVED***
	ip := net.ParseIP(host)
	if ip == nil ***REMOVED***
		return true
	***REMOVED***
	return ip.IsLoopback()
***REMOVED***

func (c *clusterState) slotMasterNode(slot int) (*clusterNode, error) ***REMOVED***
	nodes := c.slotNodes(slot)
	if len(nodes) > 0 ***REMOVED***
		return nodes[0], nil
	***REMOVED***
	return c.nodes.Random()
***REMOVED***

func (c *clusterState) slotSlaveNode(slot int) (*clusterNode, error) ***REMOVED***
	nodes := c.slotNodes(slot)
	switch len(nodes) ***REMOVED***
	case 0:
		return c.nodes.Random()
	case 1:
		return nodes[0], nil
	case 2:
		if slave := nodes[1]; !slave.Failing() ***REMOVED***
			return slave, nil
		***REMOVED***
		return nodes[0], nil
	default:
		var slave *clusterNode
		for i := 0; i < 10; i++ ***REMOVED***
			n := rand.Intn(len(nodes)-1) + 1
			slave = nodes[n]
			if !slave.Failing() ***REMOVED***
				return slave, nil
			***REMOVED***
		***REMOVED***

		// All slaves are loading - use master.
		return nodes[0], nil
	***REMOVED***
***REMOVED***

func (c *clusterState) slotClosestNode(slot int) (*clusterNode, error) ***REMOVED***
	nodes := c.slotNodes(slot)
	if len(nodes) == 0 ***REMOVED***
		return c.nodes.Random()
	***REMOVED***

	var node *clusterNode
	for _, n := range nodes ***REMOVED***
		if n.Failing() ***REMOVED***
			continue
		***REMOVED***
		if node == nil || n.Latency() < node.Latency() ***REMOVED***
			node = n
		***REMOVED***
	***REMOVED***
	if node != nil ***REMOVED***
		return node, nil
	***REMOVED***

	// If all nodes are failing - return random node
	return c.nodes.Random()
***REMOVED***

func (c *clusterState) slotRandomNode(slot int) (*clusterNode, error) ***REMOVED***
	nodes := c.slotNodes(slot)
	if len(nodes) == 0 ***REMOVED***
		return c.nodes.Random()
	***REMOVED***
	if len(nodes) == 1 ***REMOVED***
		return nodes[0], nil
	***REMOVED***
	randomNodes := rand.Perm(len(nodes))
	for _, idx := range randomNodes ***REMOVED***
		if node := nodes[idx]; !node.Failing() ***REMOVED***
			return node, nil
		***REMOVED***
	***REMOVED***
	return nodes[randomNodes[0]], nil
***REMOVED***

func (c *clusterState) slotNodes(slot int) []*clusterNode ***REMOVED***
	i := sort.Search(len(c.slots), func(i int) bool ***REMOVED***
		return c.slots[i].end >= slot
	***REMOVED***)
	if i >= len(c.slots) ***REMOVED***
		return nil
	***REMOVED***
	x := c.slots[i]
	if slot >= x.start && slot <= x.end ***REMOVED***
		return x.nodes
	***REMOVED***
	return nil
***REMOVED***

//------------------------------------------------------------------------------

type clusterStateHolder struct ***REMOVED***
	load func(ctx context.Context) (*clusterState, error)

	state     atomic.Value
	reloading uint32 // atomic
***REMOVED***

func newClusterStateHolder(fn func(ctx context.Context) (*clusterState, error)) *clusterStateHolder ***REMOVED***
	return &clusterStateHolder***REMOVED***
		load: fn,
	***REMOVED***
***REMOVED***

func (c *clusterStateHolder) Reload(ctx context.Context) (*clusterState, error) ***REMOVED***
	state, err := c.load(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.state.Store(state)
	return state, nil
***REMOVED***

func (c *clusterStateHolder) LazyReload() ***REMOVED***
	if !atomic.CompareAndSwapUint32(&c.reloading, 0, 1) ***REMOVED***
		return
	***REMOVED***
	go func() ***REMOVED***
		defer atomic.StoreUint32(&c.reloading, 0)

		_, err := c.Reload(context.Background())
		if err != nil ***REMOVED***
			return
		***REMOVED***
		time.Sleep(200 * time.Millisecond)
	***REMOVED***()
***REMOVED***

func (c *clusterStateHolder) Get(ctx context.Context) (*clusterState, error) ***REMOVED***
	v := c.state.Load()
	if v == nil ***REMOVED***
		return c.Reload(ctx)
	***REMOVED***

	state := v.(*clusterState)
	if time.Since(state.createdAt) > 10*time.Second ***REMOVED***
		c.LazyReload()
	***REMOVED***
	return state, nil
***REMOVED***

func (c *clusterStateHolder) ReloadOrGet(ctx context.Context) (*clusterState, error) ***REMOVED***
	state, err := c.Reload(ctx)
	if err == nil ***REMOVED***
		return state, nil
	***REMOVED***
	return c.Get(ctx)
***REMOVED***

//------------------------------------------------------------------------------

type clusterClient struct ***REMOVED***
	opt           *ClusterOptions
	nodes         *clusterNodes
	state         *clusterStateHolder //nolint:structcheck
	cmdsInfoCache *cmdsInfoCache      //nolint:structcheck
***REMOVED***

// ClusterClient is a Redis Cluster client representing a pool of zero
// or more underlying connections. It's safe for concurrent use by
// multiple goroutines.
type ClusterClient struct ***REMOVED***
	*clusterClient
	cmdable
	hooks
***REMOVED***

// NewClusterClient returns a Redis Cluster client as described in
// http://redis.io/topics/cluster-spec.
func NewClusterClient(opt *ClusterOptions) *ClusterClient ***REMOVED***
	opt.init()

	c := &ClusterClient***REMOVED***
		clusterClient: &clusterClient***REMOVED***
			opt:   opt,
			nodes: newClusterNodes(opt),
		***REMOVED***,
	***REMOVED***
	c.state = newClusterStateHolder(c.loadState)
	c.cmdsInfoCache = newCmdsInfoCache(c.cmdsInfo)
	c.cmdable = c.Process

	return c
***REMOVED***

// Options returns read-only Options that were used to create the client.
func (c *ClusterClient) Options() *ClusterOptions ***REMOVED***
	return c.opt
***REMOVED***

// ReloadState reloads cluster state. If available it calls ClusterSlots func
// to get cluster slots information.
func (c *ClusterClient) ReloadState(ctx context.Context) ***REMOVED***
	c.state.LazyReload()
***REMOVED***

// Close closes the cluster client, releasing any open resources.
//
// It is rare to Close a ClusterClient, as the ClusterClient is meant
// to be long-lived and shared between many goroutines.
func (c *ClusterClient) Close() error ***REMOVED***
	return c.nodes.Close()
***REMOVED***

// Do creates a Cmd from the args and processes the cmd.
func (c *ClusterClient) Do(ctx context.Context, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	cmd := NewCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

func (c *ClusterClient) Process(ctx context.Context, cmd Cmder) error ***REMOVED***
	return c.hooks.process(ctx, cmd, c.process)
***REMOVED***

func (c *ClusterClient) process(ctx context.Context, cmd Cmder) error ***REMOVED***
	cmdInfo := c.cmdInfo(ctx, cmd.Name())
	slot := c.cmdSlot(ctx, cmd)

	var node *clusterNode
	var ask bool
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRedirects; attempt++ ***REMOVED***
		if attempt > 0 ***REMOVED***
			if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if node == nil ***REMOVED***
			var err error
			node, err = c.cmdNode(ctx, cmdInfo, slot)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if ask ***REMOVED***
			pipe := node.Client.Pipeline()
			_ = pipe.Process(ctx, NewCmd(ctx, "asking"))
			_ = pipe.Process(ctx, cmd)
			_, lastErr = pipe.Exec(ctx)
			ask = false
		***REMOVED*** else ***REMOVED***
			lastErr = node.Client.Process(ctx, cmd)
		***REMOVED***

		// If there is no error - we are done.
		if lastErr == nil ***REMOVED***
			return nil
		***REMOVED***
		if isReadOnly := isReadOnlyError(lastErr); isReadOnly || lastErr == pool.ErrClosed ***REMOVED***
			if isReadOnly ***REMOVED***
				c.state.LazyReload()
			***REMOVED***
			node = nil
			continue
		***REMOVED***

		// If slave is loading - pick another node.
		if c.opt.ReadOnly && isLoadingError(lastErr) ***REMOVED***
			node.MarkAsFailing()
			node = nil
			continue
		***REMOVED***

		var moved bool
		var addr string
		moved, ask, addr = isMovedError(lastErr)
		if moved || ask ***REMOVED***
			c.state.LazyReload()

			var err error
			node, err = c.nodes.GetOrCreate(addr)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		if shouldRetry(lastErr, cmd.readTimeout() == nil) ***REMOVED***
			// First retry the same node.
			if attempt == 0 ***REMOVED***
				continue
			***REMOVED***

			// Second try another node.
			node.MarkAsFailing()
			node = nil
			continue
		***REMOVED***

		return lastErr
	***REMOVED***
	return lastErr
***REMOVED***

// ForEachMaster concurrently calls the fn on each master node in the cluster.
// It returns the first error if any.
func (c *ClusterClient) ForEachMaster(
	ctx context.Context,
	fn func(ctx context.Context, client *Client) error,
) error ***REMOVED***
	state, err := c.state.ReloadOrGet(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, master := range state.Masters ***REMOVED***
		wg.Add(1)
		go func(node *clusterNode) ***REMOVED***
			defer wg.Done()
			err := fn(ctx, node.Client)
			if err != nil ***REMOVED***
				select ***REMOVED***
				case errCh <- err:
				default:
				***REMOVED***
			***REMOVED***
		***REMOVED***(master)
	***REMOVED***

	wg.Wait()

	select ***REMOVED***
	case err := <-errCh:
		return err
	default:
		return nil
	***REMOVED***
***REMOVED***

// ForEachSlave concurrently calls the fn on each slave node in the cluster.
// It returns the first error if any.
func (c *ClusterClient) ForEachSlave(
	ctx context.Context,
	fn func(ctx context.Context, client *Client) error,
) error ***REMOVED***
	state, err := c.state.ReloadOrGet(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, slave := range state.Slaves ***REMOVED***
		wg.Add(1)
		go func(node *clusterNode) ***REMOVED***
			defer wg.Done()
			err := fn(ctx, node.Client)
			if err != nil ***REMOVED***
				select ***REMOVED***
				case errCh <- err:
				default:
				***REMOVED***
			***REMOVED***
		***REMOVED***(slave)
	***REMOVED***

	wg.Wait()

	select ***REMOVED***
	case err := <-errCh:
		return err
	default:
		return nil
	***REMOVED***
***REMOVED***

// ForEachShard concurrently calls the fn on each known node in the cluster.
// It returns the first error if any.
func (c *ClusterClient) ForEachShard(
	ctx context.Context,
	fn func(ctx context.Context, client *Client) error,
) error ***REMOVED***
	state, err := c.state.ReloadOrGet(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	worker := func(node *clusterNode) ***REMOVED***
		defer wg.Done()
		err := fn(ctx, node.Client)
		if err != nil ***REMOVED***
			select ***REMOVED***
			case errCh <- err:
			default:
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, node := range state.Masters ***REMOVED***
		wg.Add(1)
		go worker(node)
	***REMOVED***
	for _, node := range state.Slaves ***REMOVED***
		wg.Add(1)
		go worker(node)
	***REMOVED***

	wg.Wait()

	select ***REMOVED***
	case err := <-errCh:
		return err
	default:
		return nil
	***REMOVED***
***REMOVED***

// PoolStats returns accumulated connection pool stats.
func (c *ClusterClient) PoolStats() *PoolStats ***REMOVED***
	var acc PoolStats

	state, _ := c.state.Get(context.TODO())
	if state == nil ***REMOVED***
		return &acc
	***REMOVED***

	for _, node := range state.Masters ***REMOVED***
		s := node.Client.connPool.Stats()
		acc.Hits += s.Hits
		acc.Misses += s.Misses
		acc.Timeouts += s.Timeouts

		acc.TotalConns += s.TotalConns
		acc.IdleConns += s.IdleConns
		acc.StaleConns += s.StaleConns
	***REMOVED***

	for _, node := range state.Slaves ***REMOVED***
		s := node.Client.connPool.Stats()
		acc.Hits += s.Hits
		acc.Misses += s.Misses
		acc.Timeouts += s.Timeouts

		acc.TotalConns += s.TotalConns
		acc.IdleConns += s.IdleConns
		acc.StaleConns += s.StaleConns
	***REMOVED***

	return &acc
***REMOVED***

func (c *ClusterClient) loadState(ctx context.Context) (*clusterState, error) ***REMOVED***
	if c.opt.ClusterSlots != nil ***REMOVED***
		slots, err := c.opt.ClusterSlots(ctx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return newClusterState(c.nodes, slots, "")
	***REMOVED***

	addrs, err := c.nodes.Addrs()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var firstErr error

	for _, idx := range rand.Perm(len(addrs)) ***REMOVED***
		addr := addrs[idx]

		node, err := c.nodes.GetOrCreate(addr)
		if err != nil ***REMOVED***
			if firstErr == nil ***REMOVED***
				firstErr = err
			***REMOVED***
			continue
		***REMOVED***

		slots, err := node.Client.ClusterSlots(ctx).Result()
		if err != nil ***REMOVED***
			if firstErr == nil ***REMOVED***
				firstErr = err
			***REMOVED***
			continue
		***REMOVED***

		return newClusterState(c.nodes, slots, node.Client.opt.Addr)
	***REMOVED***

	/*
	 * No node is connectable. It's possible that all nodes' IP has changed.
	 * Clear activeAddrs to let client be able to re-connect using the initial
	 * setting of the addresses (e.g. [redis-cluster-0:6379, redis-cluster-1:6379]),
	 * which might have chance to resolve domain name and get updated IP address.
	 */
	c.nodes.mu.Lock()
	c.nodes.activeAddrs = nil
	c.nodes.mu.Unlock()

	return nil, firstErr
***REMOVED***

func (c *ClusterClient) Pipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

func (c *ClusterClient) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.Pipeline().Pipelined(ctx, fn)
***REMOVED***

func (c *ClusterClient) processPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processPipeline(ctx, cmds, c._processPipeline)
***REMOVED***

func (c *ClusterClient) _processPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	cmdsMap := newCmdsMap()
	err := c.mapCmdsByNode(ctx, cmdsMap, cmds)
	if err != nil ***REMOVED***
		setCmdsErr(cmds, err)
		return err
	***REMOVED***

	for attempt := 0; attempt <= c.opt.MaxRedirects; attempt++ ***REMOVED***
		if attempt > 0 ***REMOVED***
			if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil ***REMOVED***
				setCmdsErr(cmds, err)
				return err
			***REMOVED***
		***REMOVED***

		failedCmds := newCmdsMap()
		var wg sync.WaitGroup

		for node, cmds := range cmdsMap.m ***REMOVED***
			wg.Add(1)
			go func(node *clusterNode, cmds []Cmder) ***REMOVED***
				defer wg.Done()

				err := c._processPipelineNode(ctx, node, cmds, failedCmds)
				if err == nil ***REMOVED***
					return
				***REMOVED***
				if attempt < c.opt.MaxRedirects ***REMOVED***
					if err := c.mapCmdsByNode(ctx, failedCmds, cmds); err != nil ***REMOVED***
						setCmdsErr(cmds, err)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					setCmdsErr(cmds, err)
				***REMOVED***
			***REMOVED***(node, cmds)
		***REMOVED***

		wg.Wait()
		if len(failedCmds.m) == 0 ***REMOVED***
			break
		***REMOVED***
		cmdsMap = failedCmds
	***REMOVED***

	return cmdsFirstErr(cmds)
***REMOVED***

func (c *ClusterClient) mapCmdsByNode(ctx context.Context, cmdsMap *cmdsMap, cmds []Cmder) error ***REMOVED***
	state, err := c.state.Get(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if c.opt.ReadOnly && c.cmdsAreReadOnly(ctx, cmds) ***REMOVED***
		for _, cmd := range cmds ***REMOVED***
			slot := c.cmdSlot(ctx, cmd)
			node, err := c.slotReadOnlyNode(state, slot)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			cmdsMap.Add(node, cmd)
		***REMOVED***
		return nil
	***REMOVED***

	for _, cmd := range cmds ***REMOVED***
		slot := c.cmdSlot(ctx, cmd)
		node, err := state.slotMasterNode(slot)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		cmdsMap.Add(node, cmd)
	***REMOVED***
	return nil
***REMOVED***

func (c *ClusterClient) cmdsAreReadOnly(ctx context.Context, cmds []Cmder) bool ***REMOVED***
	for _, cmd := range cmds ***REMOVED***
		cmdInfo := c.cmdInfo(ctx, cmd.Name())
		if cmdInfo == nil || !cmdInfo.ReadOnly ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (c *ClusterClient) _processPipelineNode(
	ctx context.Context, node *clusterNode, cmds []Cmder, failedCmds *cmdsMap,
) error ***REMOVED***
	return node.Client.hooks.processPipeline(ctx, cmds, func(ctx context.Context, cmds []Cmder) error ***REMOVED***
		return node.Client.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error ***REMOVED***
			err := cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error ***REMOVED***
				return writeCmds(wr, cmds)
			***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			return cn.WithReader(ctx, c.opt.ReadTimeout, func(rd *proto.Reader) error ***REMOVED***
				return c.pipelineReadCmds(ctx, node, rd, cmds, failedCmds)
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func (c *ClusterClient) pipelineReadCmds(
	ctx context.Context,
	node *clusterNode,
	rd *proto.Reader,
	cmds []Cmder,
	failedCmds *cmdsMap,
) error ***REMOVED***
	for _, cmd := range cmds ***REMOVED***
		err := cmd.readReply(rd)
		cmd.SetErr(err)

		if err == nil ***REMOVED***
			continue
		***REMOVED***

		if c.checkMovedErr(ctx, cmd, err, failedCmds) ***REMOVED***
			continue
		***REMOVED***

		if c.opt.ReadOnly && (isLoadingError(err) || !isRedisError(err)) ***REMOVED***
			node.MarkAsFailing()
			return err
		***REMOVED***
		if isRedisError(err) ***REMOVED***
			continue
		***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (c *ClusterClient) checkMovedErr(
	ctx context.Context, cmd Cmder, err error, failedCmds *cmdsMap,
) bool ***REMOVED***
	moved, ask, addr := isMovedError(err)
	if !moved && !ask ***REMOVED***
		return false
	***REMOVED***

	node, err := c.nodes.GetOrCreate(addr)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	if moved ***REMOVED***
		c.state.LazyReload()
		failedCmds.Add(node, cmd)
		return true
	***REMOVED***

	if ask ***REMOVED***
		failedCmds.Add(node, NewCmd(ctx, "asking"), cmd)
		return true
	***REMOVED***

	panic("not reached")
***REMOVED***

// TxPipeline acts like Pipeline, but wraps queued commands with MULTI/EXEC.
func (c *ClusterClient) TxPipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processTxPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

func (c *ClusterClient) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.TxPipeline().Pipelined(ctx, fn)
***REMOVED***

func (c *ClusterClient) processTxPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processTxPipeline(ctx, cmds, c._processTxPipeline)
***REMOVED***

func (c *ClusterClient) _processTxPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	// Trim multi .. exec.
	cmds = cmds[1 : len(cmds)-1]

	state, err := c.state.Get(ctx)
	if err != nil ***REMOVED***
		setCmdsErr(cmds, err)
		return err
	***REMOVED***

	cmdsMap := c.mapCmdsBySlot(ctx, cmds)
	for slot, cmds := range cmdsMap ***REMOVED***
		node, err := state.slotMasterNode(slot)
		if err != nil ***REMOVED***
			setCmdsErr(cmds, err)
			continue
		***REMOVED***

		cmdsMap := map[*clusterNode][]Cmder***REMOVED***node: cmds***REMOVED***
		for attempt := 0; attempt <= c.opt.MaxRedirects; attempt++ ***REMOVED***
			if attempt > 0 ***REMOVED***
				if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil ***REMOVED***
					setCmdsErr(cmds, err)
					return err
				***REMOVED***
			***REMOVED***

			failedCmds := newCmdsMap()
			var wg sync.WaitGroup

			for node, cmds := range cmdsMap ***REMOVED***
				wg.Add(1)
				go func(node *clusterNode, cmds []Cmder) ***REMOVED***
					defer wg.Done()

					err := c._processTxPipelineNode(ctx, node, cmds, failedCmds)
					if err == nil ***REMOVED***
						return
					***REMOVED***

					if attempt < c.opt.MaxRedirects ***REMOVED***
						if err := c.mapCmdsByNode(ctx, failedCmds, cmds); err != nil ***REMOVED***
							setCmdsErr(cmds, err)
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						setCmdsErr(cmds, err)
					***REMOVED***
				***REMOVED***(node, cmds)
			***REMOVED***

			wg.Wait()
			if len(failedCmds.m) == 0 ***REMOVED***
				break
			***REMOVED***
			cmdsMap = failedCmds.m
		***REMOVED***
	***REMOVED***

	return cmdsFirstErr(cmds)
***REMOVED***

func (c *ClusterClient) mapCmdsBySlot(ctx context.Context, cmds []Cmder) map[int][]Cmder ***REMOVED***
	cmdsMap := make(map[int][]Cmder)
	for _, cmd := range cmds ***REMOVED***
		slot := c.cmdSlot(ctx, cmd)
		cmdsMap[slot] = append(cmdsMap[slot], cmd)
	***REMOVED***
	return cmdsMap
***REMOVED***

func (c *ClusterClient) _processTxPipelineNode(
	ctx context.Context, node *clusterNode, cmds []Cmder, failedCmds *cmdsMap,
) error ***REMOVED***
	return node.Client.hooks.processTxPipeline(ctx, cmds, func(ctx context.Context, cmds []Cmder) error ***REMOVED***
		return node.Client.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error ***REMOVED***
			err := cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error ***REMOVED***
				return writeCmds(wr, cmds)
			***REMOVED***)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			return cn.WithReader(ctx, c.opt.ReadTimeout, func(rd *proto.Reader) error ***REMOVED***
				statusCmd := cmds[0].(*StatusCmd)
				// Trim multi and exec.
				cmds = cmds[1 : len(cmds)-1]

				err := c.txPipelineReadQueued(ctx, rd, statusCmd, cmds, failedCmds)
				if err != nil ***REMOVED***
					moved, ask, addr := isMovedError(err)
					if moved || ask ***REMOVED***
						return c.cmdsMoved(ctx, cmds, moved, ask, addr, failedCmds)
					***REMOVED***
					return err
				***REMOVED***

				return pipelineReadCmds(rd, cmds)
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func (c *ClusterClient) txPipelineReadQueued(
	ctx context.Context,
	rd *proto.Reader,
	statusCmd *StatusCmd,
	cmds []Cmder,
	failedCmds *cmdsMap,
) error ***REMOVED***
	// Parse queued replies.
	if err := statusCmd.readReply(rd); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, cmd := range cmds ***REMOVED***
		err := statusCmd.readReply(rd)
		if err == nil || c.checkMovedErr(ctx, cmd, err, failedCmds) || isRedisError(err) ***REMOVED***
			continue
		***REMOVED***
		return err
	***REMOVED***

	// Parse number of replies.
	line, err := rd.ReadLine()
	if err != nil ***REMOVED***
		if err == Nil ***REMOVED***
			err = TxFailedErr
		***REMOVED***
		return err
	***REMOVED***

	if line[0] != proto.RespArray ***REMOVED***
		return fmt.Errorf("redis: expected '*', but got line %q", line)
	***REMOVED***

	return nil
***REMOVED***

func (c *ClusterClient) cmdsMoved(
	ctx context.Context, cmds []Cmder,
	moved, ask bool,
	addr string,
	failedCmds *cmdsMap,
) error ***REMOVED***
	node, err := c.nodes.GetOrCreate(addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if moved ***REMOVED***
		c.state.LazyReload()
		for _, cmd := range cmds ***REMOVED***
			failedCmds.Add(node, cmd)
		***REMOVED***
		return nil
	***REMOVED***

	if ask ***REMOVED***
		for _, cmd := range cmds ***REMOVED***
			failedCmds.Add(node, NewCmd(ctx, "asking"), cmd)
		***REMOVED***
		return nil
	***REMOVED***

	return nil
***REMOVED***

func (c *ClusterClient) Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error ***REMOVED***
	if len(keys) == 0 ***REMOVED***
		return fmt.Errorf("redis: Watch requires at least one key")
	***REMOVED***

	slot := hashtag.Slot(keys[0])
	for _, key := range keys[1:] ***REMOVED***
		if hashtag.Slot(key) != slot ***REMOVED***
			err := fmt.Errorf("redis: Watch requires all keys to be in the same slot")
			return err
		***REMOVED***
	***REMOVED***

	node, err := c.slotMasterNode(ctx, slot)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for attempt := 0; attempt <= c.opt.MaxRedirects; attempt++ ***REMOVED***
		if attempt > 0 ***REMOVED***
			if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		err = node.Client.Watch(ctx, fn, keys...)
		if err == nil ***REMOVED***
			break
		***REMOVED***

		moved, ask, addr := isMovedError(err)
		if moved || ask ***REMOVED***
			node, err = c.nodes.GetOrCreate(addr)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		if isReadOnly := isReadOnlyError(err); isReadOnly || err == pool.ErrClosed ***REMOVED***
			if isReadOnly ***REMOVED***
				c.state.LazyReload()
			***REMOVED***
			node, err = c.slotMasterNode(ctx, slot)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		if shouldRetry(err, true) ***REMOVED***
			continue
		***REMOVED***

		return err
	***REMOVED***

	return err
***REMOVED***

func (c *ClusterClient) pubSub() *PubSub ***REMOVED***
	var node *clusterNode
	pubsub := &PubSub***REMOVED***
		opt: c.opt.clientOptions(),

		newConn: func(ctx context.Context, channels []string) (*pool.Conn, error) ***REMOVED***
			if node != nil ***REMOVED***
				panic("node != nil")
			***REMOVED***

			var err error
			if len(channels) > 0 ***REMOVED***
				slot := hashtag.Slot(channels[0])
				node, err = c.slotMasterNode(ctx, slot)
			***REMOVED*** else ***REMOVED***
				node, err = c.nodes.Random()
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			cn, err := node.Client.newConn(context.TODO())
			if err != nil ***REMOVED***
				node = nil

				return nil, err
			***REMOVED***

			return cn, nil
		***REMOVED***,
		closeConn: func(cn *pool.Conn) error ***REMOVED***
			err := node.Client.connPool.CloseConn(cn)
			node = nil
			return err
		***REMOVED***,
	***REMOVED***
	pubsub.init()

	return pubsub
***REMOVED***

// Subscribe subscribes the client to the specified channels.
// Channels can be omitted to create empty subscription.
func (c *ClusterClient) Subscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	pubsub := c.pubSub()
	if len(channels) > 0 ***REMOVED***
		_ = pubsub.Subscribe(ctx, channels...)
	***REMOVED***
	return pubsub
***REMOVED***

// PSubscribe subscribes the client to the given patterns.
// Patterns can be omitted to create empty subscription.
func (c *ClusterClient) PSubscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	pubsub := c.pubSub()
	if len(channels) > 0 ***REMOVED***
		_ = pubsub.PSubscribe(ctx, channels...)
	***REMOVED***
	return pubsub
***REMOVED***

func (c *ClusterClient) retryBackoff(attempt int) time.Duration ***REMOVED***
	return internal.RetryBackoff(attempt, c.opt.MinRetryBackoff, c.opt.MaxRetryBackoff)
***REMOVED***

func (c *ClusterClient) cmdsInfo(ctx context.Context) (map[string]*CommandInfo, error) ***REMOVED***
	// Try 3 random nodes.
	const nodeLimit = 3

	addrs, err := c.nodes.Addrs()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var firstErr error

	perm := rand.Perm(len(addrs))
	if len(perm) > nodeLimit ***REMOVED***
		perm = perm[:nodeLimit]
	***REMOVED***

	for _, idx := range perm ***REMOVED***
		addr := addrs[idx]

		node, err := c.nodes.GetOrCreate(addr)
		if err != nil ***REMOVED***
			if firstErr == nil ***REMOVED***
				firstErr = err
			***REMOVED***
			continue
		***REMOVED***

		info, err := node.Client.Command(ctx).Result()
		if err == nil ***REMOVED***
			return info, nil
		***REMOVED***
		if firstErr == nil ***REMOVED***
			firstErr = err
		***REMOVED***
	***REMOVED***

	if firstErr == nil ***REMOVED***
		panic("not reached")
	***REMOVED***
	return nil, firstErr
***REMOVED***

func (c *ClusterClient) cmdInfo(ctx context.Context, name string) *CommandInfo ***REMOVED***
	cmdsInfo, err := c.cmdsInfoCache.Get(ctx)
	if err != nil ***REMOVED***
		internal.Logger.Printf(context.TODO(), "getting command info: %s", err)
		return nil
	***REMOVED***

	info := cmdsInfo[name]
	if info == nil ***REMOVED***
		internal.Logger.Printf(context.TODO(), "info for cmd=%s not found", name)
	***REMOVED***
	return info
***REMOVED***

func (c *ClusterClient) cmdSlot(ctx context.Context, cmd Cmder) int ***REMOVED***
	args := cmd.Args()
	if args[0] == "cluster" && args[1] == "getkeysinslot" ***REMOVED***
		return args[2].(int)
	***REMOVED***

	cmdInfo := c.cmdInfo(ctx, cmd.Name())
	return cmdSlot(cmd, cmdFirstKeyPos(cmd, cmdInfo))
***REMOVED***

func cmdSlot(cmd Cmder, pos int) int ***REMOVED***
	if pos == 0 ***REMOVED***
		return hashtag.RandomSlot()
	***REMOVED***
	firstKey := cmd.stringArg(pos)
	return hashtag.Slot(firstKey)
***REMOVED***

func (c *ClusterClient) cmdNode(
	ctx context.Context,
	cmdInfo *CommandInfo,
	slot int,
) (*clusterNode, error) ***REMOVED***
	state, err := c.state.Get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if c.opt.ReadOnly && cmdInfo != nil && cmdInfo.ReadOnly ***REMOVED***
		return c.slotReadOnlyNode(state, slot)
	***REMOVED***
	return state.slotMasterNode(slot)
***REMOVED***

func (c *clusterClient) slotReadOnlyNode(state *clusterState, slot int) (*clusterNode, error) ***REMOVED***
	if c.opt.RouteByLatency ***REMOVED***
		return state.slotClosestNode(slot)
	***REMOVED***
	if c.opt.RouteRandomly ***REMOVED***
		return state.slotRandomNode(slot)
	***REMOVED***
	return state.slotSlaveNode(slot)
***REMOVED***

func (c *ClusterClient) slotMasterNode(ctx context.Context, slot int) (*clusterNode, error) ***REMOVED***
	state, err := c.state.Get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return state.slotMasterNode(slot)
***REMOVED***

// SlaveForKey gets a client for a replica node to run any command on it.
// This is especially useful if we want to run a particular lua script which has
// only read only commands on the replica.
// This is because other redis commands generally have a flag that points that
// they are read only and automatically run on the replica nodes
// if ClusterOptions.ReadOnly flag is set to true.
func (c *ClusterClient) SlaveForKey(ctx context.Context, key string) (*Client, error) ***REMOVED***
	state, err := c.state.Get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	slot := hashtag.Slot(key)
	node, err := c.slotReadOnlyNode(state, slot)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return node.Client, err
***REMOVED***

// MasterForKey return a client to the master node for a particular key.
func (c *ClusterClient) MasterForKey(ctx context.Context, key string) (*Client, error) ***REMOVED***
	slot := hashtag.Slot(key)
	node, err := c.slotMasterNode(ctx, slot)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return node.Client, err
***REMOVED***

func appendUniqueNode(nodes []*clusterNode, node *clusterNode) []*clusterNode ***REMOVED***
	for _, n := range nodes ***REMOVED***
		if n == node ***REMOVED***
			return nodes
		***REMOVED***
	***REMOVED***
	return append(nodes, node)
***REMOVED***

func appendIfNotExists(ss []string, es ...string) []string ***REMOVED***
loop:
	for _, e := range es ***REMOVED***
		for _, s := range ss ***REMOVED***
			if s == e ***REMOVED***
				continue loop
			***REMOVED***
		***REMOVED***
		ss = append(ss, e)
	***REMOVED***
	return ss
***REMOVED***

//------------------------------------------------------------------------------

type cmdsMap struct ***REMOVED***
	mu sync.Mutex
	m  map[*clusterNode][]Cmder
***REMOVED***

func newCmdsMap() *cmdsMap ***REMOVED***
	return &cmdsMap***REMOVED***
		m: make(map[*clusterNode][]Cmder),
	***REMOVED***
***REMOVED***

func (m *cmdsMap) Add(node *clusterNode, cmds ...Cmder) ***REMOVED***
	m.mu.Lock()
	m.m[node] = append(m.m[node], cmds...)
	m.mu.Unlock()
***REMOVED***
