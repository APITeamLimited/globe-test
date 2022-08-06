package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v9/internal"
	"github.com/go-redis/redis/v9/internal/pool"
	"github.com/go-redis/redis/v9/internal/proto"
)

// Nil reply returned by Redis when key does not exist.
const Nil = proto.Nil

// SetLogger set custom log
func SetLogger(logger internal.Logging) ***REMOVED***
	internal.Logger = logger
***REMOVED***

//------------------------------------------------------------------------------

type Hook interface ***REMOVED***
	BeforeProcess(ctx context.Context, cmd Cmder) (context.Context, error)
	AfterProcess(ctx context.Context, cmd Cmder) error

	BeforeProcessPipeline(ctx context.Context, cmds []Cmder) (context.Context, error)
	AfterProcessPipeline(ctx context.Context, cmds []Cmder) error
***REMOVED***

type hooks struct ***REMOVED***
	hooks []Hook
***REMOVED***

func (hs *hooks) lock() ***REMOVED***
	hs.hooks = hs.hooks[:len(hs.hooks):len(hs.hooks)]
***REMOVED***

func (hs hooks) clone() hooks ***REMOVED***
	clone := hs
	clone.lock()
	return clone
***REMOVED***

func (hs *hooks) AddHook(hook Hook) ***REMOVED***
	hs.hooks = append(hs.hooks, hook)
***REMOVED***

func (hs hooks) process(
	ctx context.Context, cmd Cmder, fn func(context.Context, Cmder) error,
) error ***REMOVED***
	if len(hs.hooks) == 0 ***REMOVED***
		err := fn(ctx, cmd)
		cmd.SetErr(err)
		return err
	***REMOVED***

	var hookIndex int
	var retErr error

	for ; hookIndex < len(hs.hooks) && retErr == nil; hookIndex++ ***REMOVED***
		ctx, retErr = hs.hooks[hookIndex].BeforeProcess(ctx, cmd)
		if retErr != nil ***REMOVED***
			cmd.SetErr(retErr)
		***REMOVED***
	***REMOVED***

	if retErr == nil ***REMOVED***
		retErr = fn(ctx, cmd)
		cmd.SetErr(retErr)
	***REMOVED***

	for hookIndex--; hookIndex >= 0; hookIndex-- ***REMOVED***
		if err := hs.hooks[hookIndex].AfterProcess(ctx, cmd); err != nil ***REMOVED***
			retErr = err
			cmd.SetErr(retErr)
		***REMOVED***
	***REMOVED***

	return retErr
***REMOVED***

func (hs hooks) processPipeline(
	ctx context.Context, cmds []Cmder, fn func(context.Context, []Cmder) error,
) error ***REMOVED***
	if len(hs.hooks) == 0 ***REMOVED***
		err := fn(ctx, cmds)
		return err
	***REMOVED***

	var hookIndex int
	var retErr error

	for ; hookIndex < len(hs.hooks) && retErr == nil; hookIndex++ ***REMOVED***
		ctx, retErr = hs.hooks[hookIndex].BeforeProcessPipeline(ctx, cmds)
		if retErr != nil ***REMOVED***
			setCmdsErr(cmds, retErr)
		***REMOVED***
	***REMOVED***

	if retErr == nil ***REMOVED***
		retErr = fn(ctx, cmds)
	***REMOVED***

	for hookIndex--; hookIndex >= 0; hookIndex-- ***REMOVED***
		if err := hs.hooks[hookIndex].AfterProcessPipeline(ctx, cmds); err != nil ***REMOVED***
			retErr = err
			setCmdsErr(cmds, retErr)
		***REMOVED***
	***REMOVED***

	return retErr
***REMOVED***

func (hs hooks) processTxPipeline(
	ctx context.Context, cmds []Cmder, fn func(context.Context, []Cmder) error,
) error ***REMOVED***
	cmds = wrapMultiExec(ctx, cmds)
	return hs.processPipeline(ctx, cmds, fn)
***REMOVED***

//------------------------------------------------------------------------------

type baseClient struct ***REMOVED***
	opt      *Options
	connPool pool.Pooler

	onClose func() error // hook called when client is closed
***REMOVED***

func newBaseClient(opt *Options, connPool pool.Pooler) *baseClient ***REMOVED***
	return &baseClient***REMOVED***
		opt:      opt,
		connPool: connPool,
	***REMOVED***
***REMOVED***

func (c *baseClient) clone() *baseClient ***REMOVED***
	clone := *c
	return &clone
***REMOVED***

func (c *baseClient) withTimeout(timeout time.Duration) *baseClient ***REMOVED***
	opt := c.opt.clone()
	opt.ReadTimeout = timeout
	opt.WriteTimeout = timeout

	clone := c.clone()
	clone.opt = opt

	return clone
***REMOVED***

func (c *baseClient) String() string ***REMOVED***
	return fmt.Sprintf("Redis<%s db:%d>", c.getAddr(), c.opt.DB)
***REMOVED***

func (c *baseClient) newConn(ctx context.Context) (*pool.Conn, error) ***REMOVED***
	cn, err := c.connPool.NewConn(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = c.initConn(ctx, cn)
	if err != nil ***REMOVED***
		_ = c.connPool.CloseConn(cn)
		return nil, err
	***REMOVED***

	return cn, nil
***REMOVED***

func (c *baseClient) getConn(ctx context.Context) (*pool.Conn, error) ***REMOVED***
	if c.opt.Limiter != nil ***REMOVED***
		err := c.opt.Limiter.Allow()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	cn, err := c._getConn(ctx)
	if err != nil ***REMOVED***
		if c.opt.Limiter != nil ***REMOVED***
			c.opt.Limiter.ReportResult(err)
		***REMOVED***
		return nil, err
	***REMOVED***

	return cn, nil
***REMOVED***

func (c *baseClient) _getConn(ctx context.Context) (*pool.Conn, error) ***REMOVED***
	cn, err := c.connPool.Get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if cn.Inited ***REMOVED***
		return cn, nil
	***REMOVED***

	if err := c.initConn(ctx, cn); err != nil ***REMOVED***
		c.connPool.Remove(ctx, cn, err)
		if err := errors.Unwrap(err); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, err
	***REMOVED***

	return cn, nil
***REMOVED***

func (c *baseClient) initConn(ctx context.Context, cn *pool.Conn) error ***REMOVED***
	if cn.Inited ***REMOVED***
		return nil
	***REMOVED***
	cn.Inited = true

	username, password := c.opt.Username, c.opt.Password
	if c.opt.CredentialsProvider != nil ***REMOVED***
		username, password = c.opt.CredentialsProvider()
	***REMOVED***

	connPool := pool.NewSingleConnPool(c.connPool, cn)
	conn := newConn(c.opt, connPool)

	var auth bool

	// For redis-server < 6.0 that does not support the Hello command,
	// we continue to provide services with RESP2.
	if err := conn.Hello(ctx, 3, username, password, "").Err(); err == nil ***REMOVED***
		auth = true
	***REMOVED*** else if !strings.HasPrefix(err.Error(), "ERR unknown command") ***REMOVED***
		return err
	***REMOVED***

	_, err := conn.Pipelined(ctx, func(pipe Pipeliner) error ***REMOVED***
		if !auth && password != "" ***REMOVED***
			if username != "" ***REMOVED***
				pipe.AuthACL(ctx, username, password)
			***REMOVED*** else ***REMOVED***
				pipe.Auth(ctx, password)
			***REMOVED***
		***REMOVED***

		if c.opt.DB > 0 ***REMOVED***
			pipe.Select(ctx, c.opt.DB)
		***REMOVED***

		if c.opt.readOnly ***REMOVED***
			pipe.ReadOnly(ctx)
		***REMOVED***

		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if c.opt.OnConnect != nil ***REMOVED***
		return c.opt.OnConnect(ctx, conn)
	***REMOVED***
	return nil
***REMOVED***

func (c *baseClient) releaseConn(ctx context.Context, cn *pool.Conn, err error) ***REMOVED***
	if c.opt.Limiter != nil ***REMOVED***
		c.opt.Limiter.ReportResult(err)
	***REMOVED***

	if isBadConn(err, false, c.opt.Addr) ***REMOVED***
		c.connPool.Remove(ctx, cn, err)
	***REMOVED*** else ***REMOVED***
		c.connPool.Put(ctx, cn)
	***REMOVED***
***REMOVED***

func (c *baseClient) withConn(
	ctx context.Context, fn func(context.Context, *pool.Conn) error,
) error ***REMOVED***
	cn, err := c.getConn(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	defer func() ***REMOVED***
		c.releaseConn(ctx, cn, err)
	***REMOVED***()

	done := ctx.Done() //nolint:ifshort

	if done == nil ***REMOVED***
		err = fn(ctx, cn)
		return err
	***REMOVED***

	errc := make(chan error, 1)
	go func() ***REMOVED*** errc <- fn(ctx, cn) ***REMOVED***()

	select ***REMOVED***
	case <-done:
		_ = cn.Close()
		// Wait for the goroutine to finish and send something.
		<-errc

		err = ctx.Err()
		return err
	case err = <-errc:
		return err
	***REMOVED***
***REMOVED***

func (c *baseClient) process(ctx context.Context, cmd Cmder) error ***REMOVED***
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ ***REMOVED***
		attempt := attempt

		retry, err := c._process(ctx, cmd, attempt)
		if err == nil || !retry ***REMOVED***
			return err
		***REMOVED***

		lastErr = err
	***REMOVED***
	return lastErr
***REMOVED***

func (c *baseClient) _process(ctx context.Context, cmd Cmder, attempt int) (bool, error) ***REMOVED***
	if attempt > 0 ***REMOVED***
		if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil ***REMOVED***
			return false, err
		***REMOVED***
	***REMOVED***

	retryTimeout := uint32(1)
	err := c.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error ***REMOVED***
		err := cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error ***REMOVED***
			return writeCmd(wr, cmd)
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = cn.WithReader(ctx, c.cmdTimeout(cmd), cmd.readReply)
		if err != nil ***REMOVED***
			if cmd.readTimeout() == nil ***REMOVED***
				atomic.StoreUint32(&retryTimeout, 1)
			***REMOVED***
			return err
		***REMOVED***

		return nil
	***REMOVED***)
	if err == nil ***REMOVED***
		return false, nil
	***REMOVED***

	retry := shouldRetry(err, atomic.LoadUint32(&retryTimeout) == 1)
	return retry, err
***REMOVED***

func (c *baseClient) retryBackoff(attempt int) time.Duration ***REMOVED***
	return internal.RetryBackoff(attempt, c.opt.MinRetryBackoff, c.opt.MaxRetryBackoff)
***REMOVED***

func (c *baseClient) cmdTimeout(cmd Cmder) time.Duration ***REMOVED***
	if timeout := cmd.readTimeout(); timeout != nil ***REMOVED***
		t := *timeout
		if t == 0 ***REMOVED***
			return 0
		***REMOVED***
		return t + 10*time.Second
	***REMOVED***
	return c.opt.ReadTimeout
***REMOVED***

// Close closes the client, releasing any open resources.
//
// It is rare to Close a Client, as the Client is meant to be
// long-lived and shared between many goroutines.
func (c *baseClient) Close() error ***REMOVED***
	var firstErr error
	if c.onClose != nil ***REMOVED***
		if err := c.onClose(); err != nil ***REMOVED***
			firstErr = err
		***REMOVED***
	***REMOVED***
	if err := c.connPool.Close(); err != nil && firstErr == nil ***REMOVED***
		firstErr = err
	***REMOVED***
	return firstErr
***REMOVED***

func (c *baseClient) getAddr() string ***REMOVED***
	return c.opt.Addr
***REMOVED***

func (c *baseClient) processPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.generalProcessPipeline(ctx, cmds, c.pipelineProcessCmds)
***REMOVED***

func (c *baseClient) processTxPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.generalProcessPipeline(ctx, cmds, c.txPipelineProcessCmds)
***REMOVED***

type pipelineProcessor func(context.Context, *pool.Conn, []Cmder) (bool, error)

func (c *baseClient) generalProcessPipeline(
	ctx context.Context, cmds []Cmder, p pipelineProcessor,
) error ***REMOVED***
	err := c._generalProcessPipeline(ctx, cmds, p)
	if err != nil ***REMOVED***
		setCmdsErr(cmds, err)
		return err
	***REMOVED***
	return cmdsFirstErr(cmds)
***REMOVED***

func (c *baseClient) _generalProcessPipeline(
	ctx context.Context, cmds []Cmder, p pipelineProcessor,
) error ***REMOVED***
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ ***REMOVED***
		if attempt > 0 ***REMOVED***
			if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		var canRetry bool
		lastErr = c.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error ***REMOVED***
			var err error
			canRetry, err = p(ctx, cn, cmds)
			return err
		***REMOVED***)
		if lastErr == nil || !canRetry || !shouldRetry(lastErr, true) ***REMOVED***
			return lastErr
		***REMOVED***
	***REMOVED***
	return lastErr
***REMOVED***

func (c *baseClient) pipelineProcessCmds(
	ctx context.Context, cn *pool.Conn, cmds []Cmder,
) (bool, error) ***REMOVED***
	err := cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error ***REMOVED***
		return writeCmds(wr, cmds)
	***REMOVED***)
	if err != nil ***REMOVED***
		return true, err
	***REMOVED***

	err = cn.WithReader(ctx, c.opt.ReadTimeout, func(rd *proto.Reader) error ***REMOVED***
		return pipelineReadCmds(rd, cmds)
	***REMOVED***)
	return true, err
***REMOVED***

func pipelineReadCmds(rd *proto.Reader, cmds []Cmder) error ***REMOVED***
	for _, cmd := range cmds ***REMOVED***
		err := cmd.readReply(rd)
		cmd.SetErr(err)
		if err != nil && !isRedisError(err) ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *baseClient) txPipelineProcessCmds(
	ctx context.Context, cn *pool.Conn, cmds []Cmder,
) (bool, error) ***REMOVED***
	err := cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error ***REMOVED***
		return writeCmds(wr, cmds)
	***REMOVED***)
	if err != nil ***REMOVED***
		return true, err
	***REMOVED***

	err = cn.WithReader(ctx, c.opt.ReadTimeout, func(rd *proto.Reader) error ***REMOVED***
		statusCmd := cmds[0].(*StatusCmd)
		// Trim multi and exec.
		cmds = cmds[1 : len(cmds)-1]

		err := txPipelineReadQueued(rd, statusCmd, cmds)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		return pipelineReadCmds(rd, cmds)
	***REMOVED***)
	return false, err
***REMOVED***

func wrapMultiExec(ctx context.Context, cmds []Cmder) []Cmder ***REMOVED***
	if len(cmds) == 0 ***REMOVED***
		panic("not reached")
	***REMOVED***
	cmdCopy := make([]Cmder, len(cmds)+2)
	cmdCopy[0] = NewStatusCmd(ctx, "multi")
	copy(cmdCopy[1:], cmds)
	cmdCopy[len(cmdCopy)-1] = NewSliceCmd(ctx, "exec")
	return cmdCopy
***REMOVED***

func txPipelineReadQueued(rd *proto.Reader, statusCmd *StatusCmd, cmds []Cmder) error ***REMOVED***
	// Parse +OK.
	if err := statusCmd.readReply(rd); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Parse +QUEUED.
	for range cmds ***REMOVED***
		if err := statusCmd.readReply(rd); err != nil && !isRedisError(err) ***REMOVED***
			return err
		***REMOVED***
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

//------------------------------------------------------------------------------

// Client is a Redis client representing a pool of zero or more underlying connections.
// It's safe for concurrent use by multiple goroutines.
//
// Client creates and frees connections automatically; it also maintains a free pool
// of idle connections. You can control the pool size with Config.PoolSize option.
type Client struct ***REMOVED***
	*baseClient
	cmdable
	hooks
***REMOVED***

// NewClient returns a client to the Redis Server specified by Options.
func NewClient(opt *Options) *Client ***REMOVED***
	opt.init()

	c := Client***REMOVED***
		baseClient: newBaseClient(opt, newConnPool(opt)),
	***REMOVED***
	c.cmdable = c.Process

	return &c
***REMOVED***

func (c *Client) clone() *Client ***REMOVED***
	clone := *c
	clone.cmdable = clone.Process
	clone.hooks.lock()
	return &clone
***REMOVED***

func (c *Client) WithTimeout(timeout time.Duration) *Client ***REMOVED***
	clone := c.clone()
	clone.baseClient = c.baseClient.withTimeout(timeout)
	return clone
***REMOVED***

func (c *Client) Conn() *Conn ***REMOVED***
	return newConn(c.opt, pool.NewStickyConnPool(c.connPool))
***REMOVED***

// Do creates a Cmd from the args and processes the cmd.
func (c *Client) Do(ctx context.Context, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	cmd := NewCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

func (c *Client) Process(ctx context.Context, cmd Cmder) error ***REMOVED***
	return c.hooks.process(ctx, cmd, c.baseClient.process)
***REMOVED***

func (c *Client) processPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processPipeline(ctx, cmds, c.baseClient.processPipeline)
***REMOVED***

func (c *Client) processTxPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processTxPipeline(ctx, cmds, c.baseClient.processTxPipeline)
***REMOVED***

// Options returns read-only Options that were used to create the client.
func (c *Client) Options() *Options ***REMOVED***
	return c.opt
***REMOVED***

type PoolStats pool.Stats

// PoolStats returns connection pool stats.
func (c *Client) PoolStats() *PoolStats ***REMOVED***
	stats := c.connPool.Stats()
	return (*PoolStats)(stats)
***REMOVED***

func (c *Client) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.Pipeline().Pipelined(ctx, fn)
***REMOVED***

func (c *Client) Pipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

func (c *Client) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.TxPipeline().Pipelined(ctx, fn)
***REMOVED***

// TxPipeline acts like Pipeline, but wraps queued commands with MULTI/EXEC.
func (c *Client) TxPipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processTxPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

func (c *Client) pubSub() *PubSub ***REMOVED***
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

// Subscribe subscribes the client to the specified channels.
// Channels can be omitted to create empty subscription.
// Note that this method does not wait on a response from Redis, so the
// subscription may not be active immediately. To force the connection to wait,
// you may call the Receive() method on the returned *PubSub like so:
//
//    sub := client.Subscribe(queryResp)
//    iface, err := sub.Receive()
//    if err != nil ***REMOVED***
//        // handle error
//    ***REMOVED***
//
//    // Should be *Subscription, but others are possible if other actions have been
//    // taken on sub since it was created.
//    switch iface.(type) ***REMOVED***
//    case *Subscription:
//        // subscribe succeeded
//    case *Message:
//        // received first message
//    case *Pong:
//        // pong received
//    default:
//        // handle error
//    ***REMOVED***
//
//    ch := sub.Channel()
func (c *Client) Subscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	pubsub := c.pubSub()
	if len(channels) > 0 ***REMOVED***
		_ = pubsub.Subscribe(ctx, channels...)
	***REMOVED***
	return pubsub
***REMOVED***

// PSubscribe subscribes the client to the given patterns.
// Patterns can be omitted to create empty subscription.
func (c *Client) PSubscribe(ctx context.Context, channels ...string) *PubSub ***REMOVED***
	pubsub := c.pubSub()
	if len(channels) > 0 ***REMOVED***
		_ = pubsub.PSubscribe(ctx, channels...)
	***REMOVED***
	return pubsub
***REMOVED***

//------------------------------------------------------------------------------

type conn struct ***REMOVED***
	baseClient
	cmdable
	statefulCmdable
	hooks // TODO: inherit hooks
***REMOVED***

// Conn represents a single Redis connection rather than a pool of connections.
// Prefer running commands from Client unless there is a specific need
// for a continuous single Redis connection.
type Conn struct ***REMOVED***
	*conn
***REMOVED***

func newConn(opt *Options, connPool pool.Pooler) *Conn ***REMOVED***
	c := Conn***REMOVED***
		conn: &conn***REMOVED***
			baseClient: baseClient***REMOVED***
				opt:      opt,
				connPool: connPool,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	c.cmdable = c.Process
	c.statefulCmdable = c.Process
	return &c
***REMOVED***

func (c *Conn) Process(ctx context.Context, cmd Cmder) error ***REMOVED***
	return c.hooks.process(ctx, cmd, c.baseClient.process)
***REMOVED***

func (c *Conn) processPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processPipeline(ctx, cmds, c.baseClient.processPipeline)
***REMOVED***

func (c *Conn) processTxPipeline(ctx context.Context, cmds []Cmder) error ***REMOVED***
	return c.hooks.processTxPipeline(ctx, cmds, c.baseClient.processTxPipeline)
***REMOVED***

func (c *Conn) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.Pipeline().Pipelined(ctx, fn)
***REMOVED***

func (c *Conn) Pipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

func (c *Conn) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.TxPipeline().Pipelined(ctx, fn)
***REMOVED***

// TxPipeline acts like Pipeline, but wraps queued commands with MULTI/EXEC.
func (c *Conn) TxPipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: c.processTxPipeline,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***
