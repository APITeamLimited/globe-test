package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cespare/xxhash/v2"
	rendezvous "github.com/dgryski/go-rendezvous" //nolint

	"github.com/go-redis/redis/v9/internal"
	"github.com/go-redis/redis/v9/internal/hashtag"
	"github.com/go-redis/redis/v9/internal/pool"
	"github.com/go-redis/redis/v9/internal/rand"
)

var errRingShardsDown = errors.New("redis: all ring shards are down")

//------------------------------------------------------------------------------

type ConsistentHash interface ***REMOVED***
	Get(string) string
***REMOVED***

type rendezvousWrapper struct ***REMOVED***
	*rendezvous.Rendezvous
***REMOVED***

func (w rendezvousWrapper) Get(key string) string ***REMOVED***
	return w.Lookup(key)
***REMOVED***

func newRendezvous(shards []string) ConsistentHash ***REMOVED***
	return rendezvousWrapper***REMOVED***rendezvous.New(shards, xxhash.Sum64String)***REMOVED***
***REMOVED***

//------------------------------------------------------------------------------

// RingOptions are used to configure a ring client and should be
// passed to NewRing.
type RingOptions struct ***REMOVED***
	// Map of name => host:port addresses of ring shards.
	Addrs map[string]string

	// NewClient creates a shard client with provided name and options.
	NewClient func(name string, opt *Options) *Client

	// Frequency of PING commands sent to check shards availability.
	// Shard is considered down after 3 subsequent failed checks.
	HeartbeatFrequency time.Duration

	// NewConsistentHash returns a consistent hash that is used
	// to distribute keys across the shards.
	//
	// See https://medium.com/@dgryski/consistent-hashing-algorithmic-tradeoffs-ef6b8e2fcae8
	// for consistent hashing algorithmic tradeoffs.
	NewConsistentHash func(shards []string) ConsistentHash

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

	// PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolFIFO bool

	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration

	TLSConfig *tls.Config
	Limiter   Limiter
***REMOVED***

func (opt *RingOptions) init() ***REMOVED***
	if opt.NewClient == nil ***REMOVED***
		opt.NewClient = func(name string, opt *Options) *Client ***REMOVED***
			return NewClient(opt)
		***REMOVED***
	***REMOVED***

	if opt.HeartbeatFrequency == 0 ***REMOVED***
		opt.HeartbeatFrequency = 500 * time.Millisecond
	***REMOVED***

	if opt.NewConsistentHash == nil ***REMOVED***
		opt.NewConsistentHash = newRendezvous
	***REMOVED***

	if opt.MaxRetries == -1 ***REMOVED***
		opt.MaxRetries = 0
	***REMOVED*** else if opt.MaxRetries == 0 ***REMOVED***
		opt.MaxRetries = 3
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
***REMOVED***

func (opt *RingOptions) clientOptions() *Options ***REMOVED***
	return &Options***REMOVED***
		Dialer:    opt.Dialer,
		OnConnect: opt.OnConnect,

		Username: opt.Username,
		Password: opt.Password,
		DB:       opt.DB,

		MaxRetries: -1,

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
		Limiter:   opt.Limiter,
	***REMOVED***
***REMOVED***

//------------------------------------------------------------------------------

type ringShard struct ***REMOVED***
	Client *Client
	down   int32
***REMOVED***

func newRingShard(opt *RingOptions, name, addr string) *ringShard ***REMOVED***
	clopt := opt.clientOptions()
	clopt.Addr = addr

	return &ringShard***REMOVED***
		Client: opt.NewClient(name, clopt),
	***REMOVED***
***REMOVED***

func (shard *ringShard) String() string ***REMOVED***
	var state string
	if shard.IsUp() ***REMOVED***
		state = "up"
	***REMOVED*** else ***REMOVED***
		state = "down"
	***REMOVED***
	return fmt.Sprintf("%s is %s", shard.Client, state)
***REMOVED***

func (shard *ringShard) IsDown() bool ***REMOVED***
	const threshold = 3
	return atomic.LoadInt32(&shard.down) >= threshold
***REMOVED***

func (shard *ringShard) IsUp() bool ***REMOVED***
	return !shard.IsDown()
***REMOVED***

// Vote votes to set shard state and returns true if state was changed.
func (shard *ringShard) Vote(up bool) bool ***REMOVED***
	if up ***REMOVED***
		changed := shard.IsDown()
		atomic.StoreInt32(&shard.down, 0)
		return changed
	***REMOVED***

	if shard.IsDown() ***REMOVED***
		return false
	***REMOVED***

	atomic.AddInt32(&shard.down, 1)
	return shard.IsDown()
***REMOVED***

//------------------------------------------------------------------------------

type ringShards struct ***REMOVED***
	opt *RingOptions

	mu       sync.RWMutex
	hash     ConsistentHash
	shards   map[string]*ringShard // read only
	list     []*ringShard          // read only
	numShard int
	closed   bool
***REMOVED***

func newRingShards(opt *RingOptions) *ringShards ***REMOVED***
	shards := make(map[string]*ringShard, len(opt.Addrs))
	list := make([]*ringShard, 0, len(shards))

	for name, addr := range opt.Addrs ***REMOVED***
		shard := newRingShard(opt, name, addr)
		shards[name] = shard

		list = append(list, shard)
	***REMOVED***

	c := &ringShards***REMOVED***
		opt: opt,

		shards: shards,
		list:   list,
	***REMOVED***
	c.rebalance()

	return c
***REMOVED***

func (c *ringShards) List() []*ringShard ***REMOVED***
	var list []*ringShard

	c.mu.RLock()
	if !c.closed ***REMOVED***
		list = c.list
	***REMOVED***
	c.mu.RUnlock()

	return list
***REMOVED***

func (c *ringShards) Hash(key string) string ***REMOVED***
	key = hashtag.Key(key)

	var hash string

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.numShard > 0 ***REMOVED***
		hash = c.hash.Get(key)
	***REMOVED***

	return hash
***REMOVED***

func (c *ringShards) GetByKey(key string) (*ringShard, error) ***REMOVED***
	key = hashtag.Key(key)

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed ***REMOVED***
		return nil, pool.ErrClosed
	***REMOVED***

	if c.numShard == 0 ***REMOVED***
		return nil, errRingShardsDown
	***REMOVED***

	hash := c.hash.Get(key)
	if hash == "" ***REMOVED***
		return nil, errRingShardsDown
	***REMOVED***

	return c.shards[hash], nil
***REMOVED***

func (c *ringShards) GetByName(shardName string) (*ringShard, error) ***REMOVED***
	if shardName == "" ***REMOVED***
		return c.Random()
	***REMOVED***

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.shards[shardName], nil
***REMOVED***

func (c *ringShards) Random() (*ringShard, error) ***REMOVED***
	return c.GetByKey(strconv.Itoa(rand.Int()))
***REMOVED***

// Heartbeat monitors state of each shard in the ring.
func (c *ringShards) Heartbeat(ctx context.Context, frequency time.Duration) ***REMOVED***
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			var rebalance bool

			for _, shard := range c.List() ***REMOVED***
				err := shard.Client.Ping(ctx).Err()
				isUp := err == nil || err == pool.ErrPoolTimeout
				if shard.Vote(isUp) ***REMOVED***
					internal.Logger.Printf(ctx, "ring shard state changed: %s", shard)
					rebalance = true
				***REMOVED***
			***REMOVED***

			if rebalance ***REMOVED***
				c.rebalance()
			***REMOVED***
		case <-ctx.Done():
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// rebalance removes dead shards from the Ring.
func (c *ringShards) rebalance() ***REMOVED***
	c.mu.RLock()
	shards := c.shards
	c.mu.RUnlock()

	liveShards := make([]string, 0, len(shards))

	for name, shard := range shards ***REMOVED***
		if shard.IsUp() ***REMOVED***
			liveShards = append(liveShards, name)
		***REMOVED***
	***REMOVED***

	hash := c.opt.NewConsistentHash(liveShards)

	c.mu.Lock()
	c.hash = hash
	c.numShard = len(liveShards)
	c.mu.Unlock()
***REMOVED***

func (c *ringShards) Len() int ***REMOVED***
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.numShard
***REMOVED***

func (c *ringShards) Close() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed ***REMOVED***
		return nil
	***REMOVED***
	c.closed = true

	var firstErr error
	for _, shard := range c.shards ***REMOVED***
		if err := shard.Client.Close(); err != nil && firstErr == nil ***REMOVED***
			firstErr = err
		***REMOVED***
	***REMOVED***

	c.hash = nil
	c.shards = nil
	c.numShard = 0
	c.list = nil

	return firstErr
***REMOVED***

//------------------------------------------------------------------------------

type ring struct ***REMOVED***
	opt               *RingOptions
	shards            *ringShards
	cmdsInfoCache     *cmdsInfoCache //nolint:structcheck
	heartbeatCancelFn context.CancelFunc
***REMOVED***

// Ring is a Redis client that uses consistent hashing to distribute
// keys across multiple Redis servers (shards). It's safe for
// concurrent use by multiple goroutines.
//
// Ring monitors the state of each shard and removes dead shards from
// the ring. When a shard comes online it is added back to the ring. This
// gives you maximum availability and partition tolerance, but no
// consistency between different shards or even clients. Each client
// uses shards that are available to the client and does not do any
// coordination when shard state is changed.
//
// Ring should be used when you need multiple Redis servers for caching
// and can tolerate losing data when one of the servers dies.
// Otherwise you should use Redis Cluster.
type Ring struct ***REMOVED***
	*ring
	cmdable
	hooks
***REMOVED***

func NewRing(opt *RingOptions) *Ring ***REMOVED***
	opt.init()

	hbCtx, hbCancel := context.WithCancel(context.Background())

	ring := Ring***REMOVED***
		ring: &ring***REMOVED***
			opt:               opt,
			shards:            newRingShards(opt),
			heartbeatCancelFn: hbCancel,
		***REMOVED***,
	***REMOVED***

	ring.cmdsInfoCache = newCmdsInfoCache(ring.cmdsInfo)
	ring.cmdable = ring.Process

	go ring.shards.Heartbeat(hbCtx, opt.HeartbeatFrequency)

	return &ring
***REMOVED***

// Do creates a Cmd from the args and processes the cmd.
func (c *Ring) Do(ctx context.Context, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	cmd := NewCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

func (c *Ring) Process(ctx context.Context, cmd Cmder) error ***REMOVED***
	return c.hooks.process(ctx, cmd, c.process)
***REMOVED***

// Options returns read-only Options that were used to create the client.
func (c *Ring) Options() *RingOptions ***REMOVED***
	return c.opt
***REMOVED***

func (c *Ring) retryBackoff(attempt int) time.Duration ***REMOVED***
	return internal.RetryBackoff(attempt, c.opt.MinRetryBackoff, c.opt.MaxRetryBackoff)
***REMOVED***

// PoolStats returns accumulated connection pool stats.
func (c *Ring) PoolStats() *PoolStats ***REMOVED***
	shards := c.shards.List()
	var acc PoolStats
	for _, shard := range shards ***REMOVED***
		s := shard.Client.connPool.Stats()
		acc.Hits += s.Hits
		acc.Misses += s.Misses
		acc.Timeouts += s.Timeouts
		acc.TotalConns += s.TotalConns
		acc.IdleConns += s.IdleConns
	***REMOVED***
	return &acc
***REMOVED***

// Len returns the current number of shards in the ring.
func (c *Ring) Len() int ***REMOVED***
	return c.shards.Len()
***REMOVED***

// Subscribe subscribes the client to the specified channels.
func (c *Ring) Subscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	if len(channels) == 0 ***REMOVED***
		panic("at least one channel is required")
	***REMOVED***

	shard, err := c.shards.GetByKey(channels[0])
	if err != nil ***REMOVED***
		// TODO: return PubSub with sticky error
		panic(err)
	***REMOVED***
	return shard.Client.Subscribe(ctx, channels...)
***REMOVED***

// PSubscribe subscribes the client to the given patterns.
func (c *Ring) PSubscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	if len(channels) == 0 ***REMOVED***
		panic("at least one channel is required")
	***REMOVED***

	shard, err := c.shards.GetByKey(channels[0])
	if err != nil ***REMOVED***
		// TODO: return PubSub with sticky error
		panic(err)
	***REMOVED***
	return shard.Client.PSubscribe(ctx, channels...)
***REMOVED***

// ForEachShard concurrently calls the fn on each live shard in the ring.
// It returns the first error if any.
func (c *Ring) ForEachShard(
	ctx context.Context,
	fn func(ctx context.Context, client *Client) error,
) error ***REMOVED***
	shards := c.shards.List()
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	for _, shard := range shards ***REMOVED***
		if shard.IsDown() ***REMOVED***
			continue
		***REMOVED***

		wg.Add(1)
		go func(shard *ringShard) ***REMOVED***
			defer wg.Done()
			err := fn(ctx, shard.Client)
			if err != nil ***REMOVED***
				select ***REMOVED***
				case errCh <- err:
				default:
				***REMOVED***
			***REMOVED***
		***REMOVED***(shard)
	***REMOVED***
	wg.Wait()

	select ***REMOVED***
	case err := <-errCh:
		return err
	default:
		return nil
	***REMOVED***
***REMOVED***

func (c *Ring) cmdsInfo(ctx context.Context) (map[string]*CommandInfo, error) ***REMOVED***
	shards := c.shards.List()
	var firstErr error
	for _, shard := range shards ***REMOVED***
		cmdsInfo, err := shard.Client.Command(ctx).Result()
		if err == nil ***REMOVED***
			return cmdsInfo, nil
		***REMOVED***
		if firstErr == nil ***REMOVED***
			firstErr = err
		***REMOVED***
	***REMOVED***
	if firstErr == nil ***REMOVED***
		return nil, errRingShardsDown
	***REMOVED***
	return nil, firstErr
***REMOVED***

func (c *Ring) cmdInfo(ctx context.Context, name string) *CommandInfo ***REMOVED***
	cmdsInfo, err := c.cmdsInfoCache.Get(ctx)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	info := cmdsInfo[name]
	if info == nil ***REMOVED***
		internal.Logger.Printf(ctx, "info for cmd=%s not found", name)
	***REMOVED***
	return info
***REMOVED***

func (c *Ring) cmdShard(ctx context.Context, cmd Cmder) (*ringShard, error) ***REMOVED***
	cmdInfo := c.cmdInfo(ctx, cmd.Name())
	pos := cmdFirstKeyPos(cmd, cmdInfo)
	if pos == 0 ***REMOVED***
		return c.shards.Random()
	***REMOVED***
	firstKey := cmd.stringArg(pos)
	return c.shards.GetByKey(firstKey)
***REMOVED***

func (c *Ring) process(ctx context.Context, cmd Cmder) error ***REMOVED***
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ ***REMOVED***
		if attempt > 0 ***REMOVED***
			if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		shard, err := c.cmdShard(ctx, cmd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		lastErr = shard.Client.Process(ctx, cmd)
		if lastErr == nil || !shouldRetry(lastErr, cmd.readTimeout() == nil) ***REMOVED***
			return lastErr
		***REMOVED***
	***REMOVED***
	return lastErr
***REMOVED***

func (c *Ring) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.Pipeline().Pipelined(ctx, fn)
***REMOVED***

func (c *Ring) Pipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

func (c *Ring) processPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processPipeline(ctx, cmds, func(ctx context.Context, cmds []Cmder) error ***REMOVED***
		return c.generalProcessPipeline(ctx, cmds, false)
	***REMOVED***)
***REMOVED***

func (c *Ring) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.TxPipeline().Pipelined(ctx, fn)
***REMOVED***

func (c *Ring) TxPipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processTxPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

func (c *Ring) processTxPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processPipeline(ctx, cmds, func(ctx context.Context, cmds []Cmder) error ***REMOVED***
		return c.generalProcessPipeline(ctx, cmds, true)
	***REMOVED***)
***REMOVED***

func (c *Ring) generalProcessPipeline(
	ctx context.Context, cmds []Cmder, tx bool,
) error ***REMOVED***
	cmdsMap := make(map[string][]Cmder)
	for _, cmd := range cmds ***REMOVED***
		cmdInfo := c.cmdInfo(ctx, cmd.Name())
		hash := cmd.stringArg(cmdFirstKeyPos(cmd, cmdInfo))
		if hash != "" ***REMOVED***
			hash = c.shards.Hash(hash)
		***REMOVED***
		cmdsMap[hash] = append(cmdsMap[hash], cmd)
	***REMOVED***

	var wg sync.WaitGroup
	for hash, cmds := range cmdsMap ***REMOVED***
		wg.Add(1)
		go func(hash string, cmds []Cmder) ***REMOVED***
			defer wg.Done()

			_ = c.processShardPipeline(ctx, hash, cmds, tx)
		***REMOVED***(hash, cmds)
	***REMOVED***

	wg.Wait()
	return cmdsFirstErr(cmds)
***REMOVED***

func (c *Ring) processShardPipeline(
	ctx context.Context, hash string, cmds []Cmder, tx bool,
) error ***REMOVED***
	// TODO: retry?
	shard, err := c.shards.GetByName(hash)
	if err != nil ***REMOVED***
		setCmdsErr(cmds, err)
		return err
	***REMOVED***

	if tx ***REMOVED***
		return shard.Client.processTxPipeline(ctx, cmds)
	***REMOVED***
	return shard.Client.processPipeline(ctx, cmds)
***REMOVED***

func (c *Ring) Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error ***REMOVED***
	if len(keys) == 0 ***REMOVED***
		return fmt.Errorf("redis: Watch requires at least one key")
	***REMOVED***

	var shards []*ringShard
	for _, key := range keys ***REMOVED***
		if key != "" ***REMOVED***
			shard, err := c.shards.GetByKey(hashtag.Key(key))
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			shards = append(shards, shard)
		***REMOVED***
	***REMOVED***

	if len(shards) == 0 ***REMOVED***
		return fmt.Errorf("redis: Watch requires at least one shard")
	***REMOVED***

	if len(shards) > 1 ***REMOVED***
		for _, shard := range shards[1:] ***REMOVED***
			if shard.Client != shards[0].Client ***REMOVED***
				err := fmt.Errorf("redis: Watch requires all keys to be in the same shard")
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return shards[0].Client.Watch(ctx, fn, keys...)
***REMOVED***

// Close closes the ring client, releasing any open resources.
//
// It is rare to Close a Ring, as the Ring is meant to be long-lived
// and shared between many goroutines.
func (c *Ring) Close() error ***REMOVED***
	c.heartbeatCancelFn()

	return c.shards.Close()
***REMOVED***
