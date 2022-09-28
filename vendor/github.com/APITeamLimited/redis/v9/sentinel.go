package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/APITeamLimited/redis/v9/internal"
	"github.com/APITeamLimited/redis/v9/internal/pool"
	"github.com/APITeamLimited/redis/v9/internal/rand"
)

//------------------------------------------------------------------------------

// FailoverOptions are used to configure a failover client and should
// be passed to NewFailoverClient.
type FailoverOptions struct ***REMOVED***
	// The master name.
	MasterName string
	// A seed list of host:port addresses of sentinel nodes.
	SentinelAddrs []string

	// If specified with SentinelPassword, enables ACL-based authentication (via
	// AUTH <user> <pass>).
	SentinelUsername string
	// Sentinel password from "requirepass <password>" (if enabled) in Sentinel
	// configuration, or, if SentinelUsername is also supplied, used for ACL-based
	// authentication.
	SentinelPassword string

	// Allows routing read-only commands to the closest master or replica node.
	// This option only works with NewFailoverClusterClient.
	RouteByLatency bool
	// Allows routing read-only commands to the random master or replica node.
	// This option only works with NewFailoverClusterClient.
	RouteRandomly bool

	// Route all commands to replica read-only nodes.
	ReplicaOnly bool

	// Use replicas disconnected with master when cannot get connected replicas
	// Now, this option only works in RandomReplicaAddr function.
	UseDisconnectedReplicas bool

	// Following options are copied from Options struct.

	Dialer    func(ctx context.Context, network, addr string) (net.Conn, error)
	OnConnect func(ctx context.Context, cn *Conn) error

	Username string
	Password string
	DB       int

	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PoolFIFO bool

	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration

	TLSConfig *tls.Config
***REMOVED***

func (opt *FailoverOptions) clientOptions() *Options ***REMOVED***
	return &Options***REMOVED***
		Addr: "FailoverClient",

		Dialer:    opt.Dialer,
		OnConnect: opt.OnConnect,

		DB:       opt.DB,
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
	***REMOVED***
***REMOVED***

func (opt *FailoverOptions) sentinelOptions(addr string) *Options ***REMOVED***
	return &Options***REMOVED***
		Addr: addr,

		Dialer:    opt.Dialer,
		OnConnect: opt.OnConnect,

		DB:       0,
		Username: opt.SentinelUsername,
		Password: opt.SentinelPassword,

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
	***REMOVED***
***REMOVED***

func (opt *FailoverOptions) clusterOptions() *ClusterOptions ***REMOVED***
	return &ClusterOptions***REMOVED***
		Dialer:    opt.Dialer,
		OnConnect: opt.OnConnect,

		Username: opt.Username,
		Password: opt.Password,

		MaxRedirects: opt.MaxRetries,

		RouteByLatency: opt.RouteByLatency,
		RouteRandomly:  opt.RouteRandomly,

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
	***REMOVED***
***REMOVED***

// NewFailoverClient returns a Redis client that uses Redis Sentinel
// for automatic failover. It's safe for concurrent use by multiple
// goroutines.
func NewFailoverClient(failoverOpt *FailoverOptions) *Client ***REMOVED***
	if failoverOpt.RouteByLatency ***REMOVED***
		panic("to route commands by latency, use NewFailoverClusterClient")
	***REMOVED***
	if failoverOpt.RouteRandomly ***REMOVED***
		panic("to route commands randomly, use NewFailoverClusterClient")
	***REMOVED***

	sentinelAddrs := make([]string, len(failoverOpt.SentinelAddrs))
	copy(sentinelAddrs, failoverOpt.SentinelAddrs)

	rand.Shuffle(len(sentinelAddrs), func(i, j int) ***REMOVED***
		sentinelAddrs[i], sentinelAddrs[j] = sentinelAddrs[j], sentinelAddrs[i]
	***REMOVED***)

	failover := &sentinelFailover***REMOVED***
		opt:           failoverOpt,
		sentinelAddrs: sentinelAddrs,
	***REMOVED***

	opt := failoverOpt.clientOptions()
	opt.Dialer = masterReplicaDialer(failover)
	opt.init()

	connPool := newConnPool(opt)

	failover.mu.Lock()
	failover.onFailover = func(ctx context.Context, addr string) ***REMOVED***
		_ = connPool.Filter(func(cn *pool.Conn) bool ***REMOVED***
			return cn.RemoteAddr().String() != addr
		***REMOVED***)
	***REMOVED***
	failover.mu.Unlock()

	c := Client***REMOVED***
		baseClient: newBaseClient(opt, connPool),
	***REMOVED***
	c.cmdable = c.Process
	c.onClose = failover.Close

	return &c
***REMOVED***

func masterReplicaDialer(
	failover *sentinelFailover,
) func(ctx context.Context, network, addr string) (net.Conn, error) ***REMOVED***
	return func(ctx context.Context, network, _ string) (net.Conn, error) ***REMOVED***
		var addr string
		var err error

		if failover.opt.ReplicaOnly ***REMOVED***
			addr, err = failover.RandomReplicaAddr(ctx)
		***REMOVED*** else ***REMOVED***
			addr, err = failover.MasterAddr(ctx)
			if err == nil ***REMOVED***
				failover.trySwitchMaster(ctx, addr)
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if failover.opt.Dialer != nil ***REMOVED***
			return failover.opt.Dialer(ctx, network, addr)
		***REMOVED***

		netDialer := &net.Dialer***REMOVED***
			Timeout:   failover.opt.DialTimeout,
			KeepAlive: 5 * time.Minute,
		***REMOVED***
		if failover.opt.TLSConfig == nil ***REMOVED***
			return netDialer.DialContext(ctx, network, addr)
		***REMOVED***
		return tls.DialWithDialer(netDialer, network, addr, failover.opt.TLSConfig)
	***REMOVED***
***REMOVED***

//------------------------------------------------------------------------------

// SentinelClient is a client for a Redis Sentinel.
type SentinelClient struct ***REMOVED***
	*baseClient
	hooks
***REMOVED***

func NewSentinelClient(opt *Options) *SentinelClient ***REMOVED***
	opt.init()
	c := &SentinelClient***REMOVED***
		baseClient: &baseClient***REMOVED***
			opt:      opt,
			connPool: newConnPool(opt),
		***REMOVED***,
	***REMOVED***
	return c
***REMOVED***

func (c *SentinelClient) Process(ctx context.Context, cmd Cmder) error ***REMOVED***
	return c.hooks.process(ctx, cmd, c.baseClient.process)
***REMOVED***

func (c *SentinelClient) pubSub() *PubSub ***REMOVED***
	pubsub := &PubSub***REMOVED***
		opt: c.opt,

		newConn: func(ctx context.Context, channels []string) (*pool.Conn, error) ***REMOVED***
			return c.newConn(ctx)
		***REMOVED***,
		closeConn: c.connPool.CloseConn,
	***REMOVED***
	pubsub.init()
	return pubsub
***REMOVED***

// Ping is used to test if a connection is still alive, or to
// measure latency.
func (c *SentinelClient) Ping(ctx context.Context) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "ping")
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Subscribe subscribes the client to the specified channels.
// Channels can be omitted to create empty subscription.
func (c *SentinelClient) Subscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	pubsub := c.pubSub()
	if len(channels) > 0 ***REMOVED***
		_ = pubsub.Subscribe(ctx, channels...)
	***REMOVED***
	return pubsub
***REMOVED***

// PSubscribe subscribes the client to the given patterns.
// Patterns can be omitted to create empty subscription.
func (c *SentinelClient) PSubscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	pubsub := c.pubSub()
	if len(channels) > 0 ***REMOVED***
		_ = pubsub.PSubscribe(ctx, channels...)
	***REMOVED***
	return pubsub
***REMOVED***

func (c *SentinelClient) GetMasterAddrByName(ctx context.Context, name string) *StringSliceCmd ***REMOVED***
	cmd := NewStringSliceCmd(ctx, "sentinel", "get-master-addr-by-name", name)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

func (c *SentinelClient) Sentinels(ctx context.Context, name string) *MapStringStringSliceCmd ***REMOVED***
	cmd := NewMapStringStringSliceCmd(ctx, "sentinel", "sentinels", name)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Failover forces a failover as if the master was not reachable, and without
// asking for agreement to other Sentinels.
func (c *SentinelClient) Failover(ctx context.Context, name string) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "sentinel", "failover", name)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Reset resets all the masters with matching name. The pattern argument is a
// glob-style pattern. The reset process clears any previous state in a master
// (including a failover in progress), and removes every replica and sentinel
// already discovered and associated with the master.
func (c *SentinelClient) Reset(ctx context.Context, pattern string) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "sentinel", "reset", pattern)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// FlushConfig forces Sentinel to rewrite its configuration on disk, including
// the current Sentinel state.
func (c *SentinelClient) FlushConfig(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "sentinel", "flushconfig")
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Master shows the state and info of the specified master.
func (c *SentinelClient) Master(ctx context.Context, name string) *MapStringStringCmd ***REMOVED***
	cmd := NewMapStringStringCmd(ctx, "sentinel", "master", name)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Masters shows a list of monitored masters and their state.
func (c *SentinelClient) Masters(ctx context.Context) *SliceCmd ***REMOVED***
	cmd := NewSliceCmd(ctx, "sentinel", "masters")
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Replicas shows a list of replicas for the specified master and their state.
func (c *SentinelClient) Replicas(ctx context.Context, name string) *MapStringStringSliceCmd ***REMOVED***
	cmd := NewMapStringStringSliceCmd(ctx, "sentinel", "replicas", name)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// CkQuorum checks if the current Sentinel configuration is able to reach the
// quorum needed to failover a master, and the majority needed to authorize the
// failover. This command should be used in monitoring systems to check if a
// Sentinel deployment is ok.
func (c *SentinelClient) CkQuorum(ctx context.Context, name string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "sentinel", "ckquorum", name)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Monitor tells the Sentinel to start monitoring a new master with the specified
// name, ip, port, and quorum.
func (c *SentinelClient) Monitor(ctx context.Context, name, ip, port, quorum string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "sentinel", "monitor", name, ip, port, quorum)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Set is used in order to change configuration parameters of a specific master.
func (c *SentinelClient) Set(ctx context.Context, name, option, value string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "sentinel", "set", name, option, value)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Remove is used in order to remove the specified master: the master will no
// longer be monitored, and will totally be removed from the internal state of
// the Sentinel.
func (c *SentinelClient) Remove(ctx context.Context, name string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "sentinel", "remove", name)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

//------------------------------------------------------------------------------

type sentinelFailover struct ***REMOVED***
	opt *FailoverOptions

	sentinelAddrs []string

	onFailover func(ctx context.Context, addr string)
	onUpdate   func(ctx context.Context)

	mu          sync.RWMutex
	_masterAddr string
	sentinel    *SentinelClient
	pubsub      *PubSub
***REMOVED***

func (c *sentinelFailover) Close() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sentinel != nil ***REMOVED***
		return c.closeSentinel()
	***REMOVED***
	return nil
***REMOVED***

func (c *sentinelFailover) closeSentinel() error ***REMOVED***
	firstErr := c.pubsub.Close()
	c.pubsub = nil

	err := c.sentinel.Close()
	if err != nil && firstErr == nil ***REMOVED***
		firstErr = err
	***REMOVED***
	c.sentinel = nil

	return firstErr
***REMOVED***

func (c *sentinelFailover) RandomReplicaAddr(ctx context.Context) (string, error) ***REMOVED***
	if c.opt == nil ***REMOVED***
		return "", errors.New("opt is nil")
	***REMOVED***

	addresses, err := c.replicaAddrs(ctx, false)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if len(addresses) == 0 && c.opt.UseDisconnectedReplicas ***REMOVED***
		addresses, err = c.replicaAddrs(ctx, true)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***

	if len(addresses) == 0 ***REMOVED***
		return c.MasterAddr(ctx)
	***REMOVED***
	return addresses[rand.Intn(len(addresses))], nil
***REMOVED***

func (c *sentinelFailover) MasterAddr(ctx context.Context) (string, error) ***REMOVED***
	c.mu.RLock()
	sentinel := c.sentinel
	c.mu.RUnlock()

	if sentinel != nil ***REMOVED***
		addr := c.getMasterAddr(ctx, sentinel)
		if addr != "" ***REMOVED***
			return addr, nil
		***REMOVED***
	***REMOVED***

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sentinel != nil ***REMOVED***
		addr := c.getMasterAddr(ctx, c.sentinel)
		if addr != "" ***REMOVED***
			return addr, nil
		***REMOVED***
		_ = c.closeSentinel()
	***REMOVED***

	for i, sentinelAddr := range c.sentinelAddrs ***REMOVED***
		sentinel := NewSentinelClient(c.opt.sentinelOptions(sentinelAddr))

		masterAddr, err := sentinel.GetMasterAddrByName(ctx, c.opt.MasterName).Result()
		if err != nil ***REMOVED***
			internal.Logger.Printf(ctx, "sentinel: GetMasterAddrByName master=%q failed: %s",
				c.opt.MasterName, err)
			_ = sentinel.Close()
			continue
		***REMOVED***

		// Push working sentinel to the top.
		c.sentinelAddrs[0], c.sentinelAddrs[i] = c.sentinelAddrs[i], c.sentinelAddrs[0]
		c.setSentinel(ctx, sentinel)

		addr := net.JoinHostPort(masterAddr[0], masterAddr[1])
		return addr, nil
	***REMOVED***

	return "", errors.New("redis: all sentinels specified in configuration are unreachable")
***REMOVED***

func (c *sentinelFailover) replicaAddrs(ctx context.Context, useDisconnected bool) ([]string, error) ***REMOVED***
	c.mu.RLock()
	sentinel := c.sentinel
	c.mu.RUnlock()

	if sentinel != nil ***REMOVED***
		addrs := c.getReplicaAddrs(ctx, sentinel)
		if len(addrs) > 0 ***REMOVED***
			return addrs, nil
		***REMOVED***
	***REMOVED***

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sentinel != nil ***REMOVED***
		addrs := c.getReplicaAddrs(ctx, c.sentinel)
		if len(addrs) > 0 ***REMOVED***
			return addrs, nil
		***REMOVED***
		_ = c.closeSentinel()
	***REMOVED***

	var sentinelReachable bool

	for i, sentinelAddr := range c.sentinelAddrs ***REMOVED***
		sentinel := NewSentinelClient(c.opt.sentinelOptions(sentinelAddr))

		replicas, err := sentinel.Replicas(ctx, c.opt.MasterName).Result()
		if err != nil ***REMOVED***
			internal.Logger.Printf(ctx, "sentinel: Replicas master=%q failed: %s",
				c.opt.MasterName, err)
			_ = sentinel.Close()
			continue
		***REMOVED***
		sentinelReachable = true
		addrs := parseReplicaAddrs(replicas, useDisconnected)
		if len(addrs) == 0 ***REMOVED***
			continue
		***REMOVED***
		// Push working sentinel to the top.
		c.sentinelAddrs[0], c.sentinelAddrs[i] = c.sentinelAddrs[i], c.sentinelAddrs[0]
		c.setSentinel(ctx, sentinel)

		return addrs, nil
	***REMOVED***

	if sentinelReachable ***REMOVED***
		return []string***REMOVED******REMOVED***, nil
	***REMOVED***
	return []string***REMOVED******REMOVED***, errors.New("redis: all sentinels specified in configuration are unreachable")
***REMOVED***

func (c *sentinelFailover) getMasterAddr(ctx context.Context, sentinel *SentinelClient) string ***REMOVED***
	addr, err := sentinel.GetMasterAddrByName(ctx, c.opt.MasterName).Result()
	if err != nil ***REMOVED***
		internal.Logger.Printf(ctx, "sentinel: GetMasterAddrByName name=%q failed: %s",
			c.opt.MasterName, err)
		return ""
	***REMOVED***
	return net.JoinHostPort(addr[0], addr[1])
***REMOVED***

func (c *sentinelFailover) getReplicaAddrs(ctx context.Context, sentinel *SentinelClient) []string ***REMOVED***
	addrs, err := sentinel.Replicas(ctx, c.opt.MasterName).Result()
	if err != nil ***REMOVED***
		internal.Logger.Printf(ctx, "sentinel: Replicas name=%q failed: %s",
			c.opt.MasterName, err)
		return nil
	***REMOVED***
	return parseReplicaAddrs(addrs, false)
***REMOVED***

func parseReplicaAddrs(addrs []map[string]string, keepDisconnected bool) []string ***REMOVED***
	nodes := make([]string, 0, len(addrs))
	for _, node := range addrs ***REMOVED***
		isDown := false
		if flags, ok := node["flags"]; ok ***REMOVED***
			for _, flag := range strings.Split(flags, ",") ***REMOVED***
				switch flag ***REMOVED***
				case "s_down", "o_down":
					isDown = true
				case "disconnected":
					if !keepDisconnected ***REMOVED***
						isDown = true
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if !isDown && node["ip"] != "" && node["port"] != "" ***REMOVED***
			nodes = append(nodes, net.JoinHostPort(node["ip"], node["port"]))
		***REMOVED***
	***REMOVED***

	return nodes
***REMOVED***

func (c *sentinelFailover) trySwitchMaster(ctx context.Context, addr string) ***REMOVED***
	c.mu.RLock()
	currentAddr := c._masterAddr //nolint:ifshort
	c.mu.RUnlock()

	if addr == currentAddr ***REMOVED***
		return
	***REMOVED***

	c.mu.Lock()
	defer c.mu.Unlock()

	if addr == c._masterAddr ***REMOVED***
		return
	***REMOVED***
	c._masterAddr = addr

	internal.Logger.Printf(ctx, "sentinel: new master=%q addr=%q",
		c.opt.MasterName, addr)
	if c.onFailover != nil ***REMOVED***
		c.onFailover(ctx, addr)
	***REMOVED***
***REMOVED***

func (c *sentinelFailover) setSentinel(ctx context.Context, sentinel *SentinelClient) ***REMOVED***
	if c.sentinel != nil ***REMOVED***
		panic("not reached")
	***REMOVED***
	c.sentinel = sentinel
	c.discoverSentinels(ctx)

	c.pubsub = sentinel.Subscribe(ctx, "+switch-master", "+replica-reconf-done")
	go c.listen(c.pubsub)
***REMOVED***

func (c *sentinelFailover) discoverSentinels(ctx context.Context) ***REMOVED***
	sentinels, err := c.sentinel.Sentinels(ctx, c.opt.MasterName).Result()
	if err != nil ***REMOVED***
		internal.Logger.Printf(ctx, "sentinel: Sentinels master=%q failed: %s", c.opt.MasterName, err)
		return
	***REMOVED***
	for _, sentinel := range sentinels ***REMOVED***
		ip, ok := sentinel["ip"]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		port, ok := sentinel["port"]
		if !ok ***REMOVED***
			continue
		***REMOVED***
		if ip != "" && port != "" ***REMOVED***
			sentinelAddr := net.JoinHostPort(ip, port)
			if !contains(c.sentinelAddrs, sentinelAddr) ***REMOVED***
				internal.Logger.Printf(ctx, "sentinel: discovered new sentinel=%q for master=%q",
					sentinelAddr, c.opt.MasterName)
				c.sentinelAddrs = append(c.sentinelAddrs, sentinelAddr)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *sentinelFailover) listen(pubsub *PubSub) ***REMOVED***
	ctx := context.TODO()

	if c.onUpdate != nil ***REMOVED***
		c.onUpdate(ctx)
	***REMOVED***

	ch := pubsub.Channel()
	for msg := range ch ***REMOVED***
		if msg.Channel == "+switch-master" ***REMOVED***
			parts := strings.Split(msg.Payload, " ")
			if parts[0] != c.opt.MasterName ***REMOVED***
				internal.Logger.Printf(pubsub.getContext(), "sentinel: ignore addr for master=%q", parts[0])
				continue
			***REMOVED***
			addr := net.JoinHostPort(parts[3], parts[4])
			c.trySwitchMaster(pubsub.getContext(), addr)
		***REMOVED***

		if c.onUpdate != nil ***REMOVED***
			c.onUpdate(ctx)
		***REMOVED***
	***REMOVED***
***REMOVED***

func contains(slice []string, str string) bool ***REMOVED***
	for _, s := range slice ***REMOVED***
		if s == str ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

//------------------------------------------------------------------------------

// NewFailoverClusterClient returns a client that supports routing read-only commands
// to a replica node.
func NewFailoverClusterClient(failoverOpt *FailoverOptions) *ClusterClient ***REMOVED***
	sentinelAddrs := make([]string, len(failoverOpt.SentinelAddrs))
	copy(sentinelAddrs, failoverOpt.SentinelAddrs)

	failover := &sentinelFailover***REMOVED***
		opt:           failoverOpt,
		sentinelAddrs: sentinelAddrs,
	***REMOVED***

	opt := failoverOpt.clusterOptions()
	opt.ClusterSlots = func(ctx context.Context) ([]ClusterSlot, error) ***REMOVED***
		masterAddr, err := failover.MasterAddr(ctx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		nodes := []ClusterNode***REMOVED******REMOVED***
			Addr: masterAddr,
		***REMOVED******REMOVED***

		replicaAddrs, err := failover.replicaAddrs(ctx, false)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		for _, replicaAddr := range replicaAddrs ***REMOVED***
			nodes = append(nodes, ClusterNode***REMOVED***
				Addr: replicaAddr,
			***REMOVED***)
		***REMOVED***

		slots := []ClusterSlot***REMOVED***
			***REMOVED***
				Start: 0,
				End:   16383,
				Nodes: nodes,
			***REMOVED***,
		***REMOVED***
		return slots, nil
	***REMOVED***

	c := NewClusterClient(opt)

	failover.mu.Lock()
	failover.onUpdate = func(ctx context.Context) ***REMOVED***
		c.ReloadState(ctx)
	***REMOVED***
	failover.mu.Unlock()

	return c
***REMOVED***
