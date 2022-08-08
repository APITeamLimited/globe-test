package redis

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v9/internal"
	"github.com/go-redis/redis/v9/internal/pool"
	"github.com/go-redis/redis/v9/internal/proto"
)

// PubSub implements Pub/Sub commands as described in
// http://redis.io/topics/pubsub. Message receiving is NOT safe
// for concurrent use by multiple goroutines.
//
// PubSub automatically reconnects to Redis Server and resubscribes
// to the channels in case of network errors.
type PubSub struct ***REMOVED***
	opt *Options

	newConn   func(ctx context.Context, channels []string) (*pool.Conn, error)
	closeConn func(*pool.Conn) error

	mu       sync.Mutex
	cn       *pool.Conn
	channels map[string]struct***REMOVED******REMOVED***
	patterns map[string]struct***REMOVED******REMOVED***

	closed bool
	exit   chan struct***REMOVED******REMOVED***

	cmd *Cmd

	chOnce sync.Once
	msgCh  *channel
	allCh  *channel
***REMOVED***

func (c *PubSub) init() ***REMOVED***
	c.exit = make(chan struct***REMOVED******REMOVED***)
***REMOVED***

func (c *PubSub) String() string ***REMOVED***
	channels := mapKeys(c.channels)
	channels = append(channels, mapKeys(c.patterns)...)
	return fmt.Sprintf("PubSub(%s)", strings.Join(channels, ", "))
***REMOVED***

func (c *PubSub) connWithLock(ctx context.Context) (*pool.Conn, error) ***REMOVED***
	c.mu.Lock()
	cn, err := c.conn(ctx, nil)
	c.mu.Unlock()
	return cn, err
***REMOVED***

func (c *PubSub) conn(ctx context.Context, newChannels []string) (*pool.Conn, error) ***REMOVED***
	if c.closed ***REMOVED***
		return nil, pool.ErrClosed
	***REMOVED***
	if c.cn != nil ***REMOVED***
		return c.cn, nil
	***REMOVED***

	channels := mapKeys(c.channels)
	channels = append(channels, newChannels...)

	cn, err := c.newConn(ctx, channels)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := c.resubscribe(ctx, cn); err != nil ***REMOVED***
		_ = c.closeConn(cn)
		return nil, err
	***REMOVED***

	c.cn = cn
	return cn, nil
***REMOVED***

func (c *PubSub) writeCmd(ctx context.Context, cn *pool.Conn, cmd Cmder) error ***REMOVED***
	return cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error ***REMOVED***
		return writeCmd(wr, cmd)
	***REMOVED***)
***REMOVED***

func (c *PubSub) resubscribe(ctx context.Context, cn *pool.Conn) error ***REMOVED***
	var firstErr error

	if len(c.channels) > 0 ***REMOVED***
		firstErr = c._subscribe(ctx, cn, "subscribe", mapKeys(c.channels))
	***REMOVED***

	if len(c.patterns) > 0 ***REMOVED***
		err := c._subscribe(ctx, cn, "psubscribe", mapKeys(c.patterns))
		if err != nil && firstErr == nil ***REMOVED***
			firstErr = err
		***REMOVED***
	***REMOVED***

	return firstErr
***REMOVED***

func mapKeys(m map[string]struct***REMOVED******REMOVED***) []string ***REMOVED***
	s := make([]string, len(m))
	i := 0
	for k := range m ***REMOVED***
		s[i] = k
		i++
	***REMOVED***
	return s
***REMOVED***

func (c *PubSub) _subscribe(
	ctx context.Context, cn *pool.Conn, redisCmd string, channels []string,
) error ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 0, 1+len(channels))
	args = append(args, redisCmd)
	for _, channel := range channels ***REMOVED***
		args = append(args, channel)
	***REMOVED***
	cmd := NewSliceCmd(ctx, args...)
	return c.writeCmd(ctx, cn, cmd)
***REMOVED***

func (c *PubSub) releaseConnWithLock(
	ctx context.Context,
	cn *pool.Conn,
	err error,
	allowTimeout bool,
) ***REMOVED***
	c.mu.Lock()
	c.releaseConn(ctx, cn, err, allowTimeout)
	c.mu.Unlock()
***REMOVED***

func (c *PubSub) releaseConn(ctx context.Context, cn *pool.Conn, err error, allowTimeout bool) ***REMOVED***
	if c.cn != cn ***REMOVED***
		return
	***REMOVED***
	if isBadConn(err, allowTimeout, c.opt.Addr) ***REMOVED***
		c.reconnect(ctx, err)
	***REMOVED***
***REMOVED***

func (c *PubSub) reconnect(ctx context.Context, reason error) ***REMOVED***
	_ = c.closeTheCn(reason)
	_, _ = c.conn(ctx, nil)
***REMOVED***

func (c *PubSub) closeTheCn(reason error) error ***REMOVED***
	if c.cn == nil ***REMOVED***
		return nil
	***REMOVED***
	if !c.closed ***REMOVED***
		//internal.Logger.Printf(c.getContext(), "redis: discarding bad PubSub connection: %s", reason)
	***REMOVED***
	err := c.closeConn(c.cn)
	c.cn = nil
	return err
***REMOVED***

func (c *PubSub) Close() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed ***REMOVED***
		return pool.ErrClosed
	***REMOVED***
	c.closed = true
	close(c.exit)

	return c.closeTheCn(pool.ErrClosed)
***REMOVED***

// Subscribe the client to the specified channels. It returns
// empty subscription if there are no channels.
func (c *PubSub) Subscribe(ctx context.Context, channels ...string) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.subscribe(ctx, "subscribe", channels...)
	if c.channels == nil ***REMOVED***
		c.channels = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	for _, s := range channels ***REMOVED***
		c.channels[s] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return err
***REMOVED***

// PSubscribe the client to the given patterns. It returns
// empty subscription if there are no patterns.
func (c *PubSub) PSubscribe(ctx context.Context, patterns ...string) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.subscribe(ctx, "psubscribe", patterns...)
	if c.patterns == nil ***REMOVED***
		c.patterns = make(map[string]struct***REMOVED******REMOVED***)
	***REMOVED***
	for _, s := range patterns ***REMOVED***
		c.patterns[s] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return err
***REMOVED***

// Unsubscribe the client from the given channels, or from all of
// them if none is given.
func (c *PubSub) Unsubscribe(ctx context.Context, channels ...string) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range channels ***REMOVED***
		delete(c.channels, channel)
	***REMOVED***
	err := c.subscribe(ctx, "unsubscribe", channels...)
	return err
***REMOVED***

// PUnsubscribe the client from the given patterns, or from all of
// them if none is given.
func (c *PubSub) PUnsubscribe(ctx context.Context, patterns ...string) error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, pattern := range patterns ***REMOVED***
		delete(c.patterns, pattern)
	***REMOVED***
	err := c.subscribe(ctx, "punsubscribe", patterns...)
	return err
***REMOVED***

func (c *PubSub) subscribe(ctx context.Context, redisCmd string, channels ...string) error ***REMOVED***
	cn, err := c.conn(ctx, channels)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = c._subscribe(ctx, cn, redisCmd, channels)
	c.releaseConn(ctx, cn, err, false)
	return err
***REMOVED***

func (c *PubSub) Ping(ctx context.Context, payload ...string) error ***REMOVED***
	args := []interface***REMOVED******REMOVED******REMOVED***"ping"***REMOVED***
	if len(payload) == 1 ***REMOVED***
		args = append(args, payload[0])
	***REMOVED***
	cmd := NewCmd(ctx, args...)

	c.mu.Lock()
	defer c.mu.Unlock()

	cn, err := c.conn(ctx, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = c.writeCmd(ctx, cn, cmd)
	c.releaseConn(ctx, cn, err, false)
	return err
***REMOVED***

// Subscription received after a successful subscription to channel.
type Subscription struct ***REMOVED***
	// Can be "subscribe", "unsubscribe", "psubscribe" or "punsubscribe".
	Kind string
	// Channel name we have subscribed to.
	Channel string
	// Number of channels we are currently subscribed to.
	Count int
***REMOVED***

func (m *Subscription) String() string ***REMOVED***
	return fmt.Sprintf("%s: %s", m.Kind, m.Channel)
***REMOVED***

// Message received as result of a PUBLISH command issued by another client.
type Message struct ***REMOVED***
	Channel      string
	Pattern      string
	Payload      string
	PayloadSlice []string
***REMOVED***

func (m *Message) String() string ***REMOVED***
	return fmt.Sprintf("Message<%s: %s>", m.Channel, m.Payload)
***REMOVED***

// Pong received as result of a PING command issued by another client.
type Pong struct ***REMOVED***
	Payload string
***REMOVED***

func (p *Pong) String() string ***REMOVED***
	if p.Payload != "" ***REMOVED***
		return fmt.Sprintf("Pong<%s>", p.Payload)
	***REMOVED***
	return "Pong"
***REMOVED***

func (c *PubSub) newMessage(reply interface***REMOVED******REMOVED***) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch reply := reply.(type) ***REMOVED***
	case string:
		return &Pong***REMOVED***
			Payload: reply,
		***REMOVED***, nil
	case []interface***REMOVED******REMOVED***:
		switch kind := reply[0].(string); kind ***REMOVED***
		case "subscribe", "unsubscribe", "psubscribe", "punsubscribe":
			// Can be nil in case of "unsubscribe".
			channel, _ := reply[1].(string)
			return &Subscription***REMOVED***
				Kind:    kind,
				Channel: channel,
				Count:   int(reply[2].(int64)),
			***REMOVED***, nil
		case "message":
			switch payload := reply[2].(type) ***REMOVED***
			case string:
				return &Message***REMOVED***
					Channel: reply[1].(string),
					Payload: payload,
				***REMOVED***, nil
			case []interface***REMOVED******REMOVED***:
				ss := make([]string, len(payload))
				for i, s := range payload ***REMOVED***
					ss[i] = s.(string)
				***REMOVED***
				return &Message***REMOVED***
					Channel:      reply[1].(string),
					PayloadSlice: ss,
				***REMOVED***, nil
			default:
				return nil, fmt.Errorf("redis: unsupported pubsub message payload: %T", payload)
			***REMOVED***
		case "pmessage":
			return &Message***REMOVED***
				Pattern: reply[1].(string),
				Channel: reply[2].(string),
				Payload: reply[3].(string),
			***REMOVED***, nil
		case "pong":
			return &Pong***REMOVED***
				Payload: reply[1].(string),
			***REMOVED***, nil
		default:
			return nil, fmt.Errorf("redis: unsupported pubsub message: %q", kind)
		***REMOVED***
	default:
		return nil, fmt.Errorf("redis: unsupported pubsub message: %#v", reply)
	***REMOVED***
***REMOVED***

// ReceiveTimeout acts like Receive but returns an error if message
// is not received in time. This is low-level API and in most cases
// Channel should be used instead.
func (c *PubSub) ReceiveTimeout(ctx context.Context, timeout time.Duration) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if c.cmd == nil ***REMOVED***
		c.cmd = NewCmd(ctx)
	***REMOVED***

	// Don't hold the lock to allow subscriptions and pings.

	cn, err := c.connWithLock(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = cn.WithReader(ctx, timeout, func(rd *proto.Reader) error ***REMOVED***
		return c.cmd.readReply(rd)
	***REMOVED***)

	c.releaseConnWithLock(ctx, cn, err, timeout > 0)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return c.newMessage(c.cmd.Val())
***REMOVED***

// Receive returns a message as a Subscription, Message, Pong or error.
// See PubSub example for details. This is low-level API and in most cases
// Channel should be used instead.
func (c *PubSub) Receive(ctx context.Context) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return c.ReceiveTimeout(ctx, 0)
***REMOVED***

// ReceiveMessage returns a Message or error ignoring Subscription and Pong
// messages. This is low-level API and in most cases Channel should be used
// instead.
func (c *PubSub) ReceiveMessage(ctx context.Context) (*Message, error) ***REMOVED***
	for ***REMOVED***
		msg, err := c.Receive(ctx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch msg := msg.(type) ***REMOVED***
		case *Subscription:
			// Ignore.
		case *Pong:
			// Ignore.
		case *Message:
			return msg, nil
		default:
			err := fmt.Errorf("redis: unknown message: %T", msg)
			return nil, err
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *PubSub) getContext() context.Context ***REMOVED***
	if c.cmd != nil ***REMOVED***
		return c.cmd.ctx
	***REMOVED***
	return context.Background()
***REMOVED***

//------------------------------------------------------------------------------

// Channel returns a Go channel for concurrently receiving messages.
// The channel is closed together with the PubSub. If the Go channel
// is blocked full for 30 seconds the message is dropped.
// Receive* APIs can not be used after channel is created.
//
// go-redis periodically sends ping messages to test connection health
// and re-subscribes if ping can not not received for 30 seconds.
func (c *PubSub) Channel(opts ...ChannelOption) <-chan *Message ***REMOVED***
	c.chOnce.Do(func() ***REMOVED***
		c.msgCh = newChannel(c, opts...)
		c.msgCh.initMsgChan()
	***REMOVED***)
	if c.msgCh == nil ***REMOVED***
		err := fmt.Errorf("redis: Channel can't be called after ChannelWithSubscriptions")
		panic(err)
	***REMOVED***
	return c.msgCh.msgCh
***REMOVED***

// ChannelSize is like Channel, but creates a Go channel
// with specified buffer size.
//
// Deprecated: use Channel(WithChannelSize(size)), remove in v9.
func (c *PubSub) ChannelSize(size int) <-chan *Message ***REMOVED***
	return c.Channel(WithChannelSize(size))
***REMOVED***

// ChannelWithSubscriptions is like Channel, but message type can be either
// *Subscription or *Message. Subscription messages can be used to detect
// reconnections.
//
// ChannelWithSubscriptions can not be used together with Channel or ChannelSize.
func (c *PubSub) ChannelWithSubscriptions(opts ...ChannelOption) <-chan interface***REMOVED******REMOVED*** ***REMOVED***
	c.chOnce.Do(func() ***REMOVED***
		c.allCh = newChannel(c, opts...)
		c.allCh.initAllChan()
	***REMOVED***)
	if c.allCh == nil ***REMOVED***
		err := fmt.Errorf("redis: ChannelWithSubscriptions can't be called after Channel")
		panic(err)
	***REMOVED***
	return c.allCh.allCh
***REMOVED***

type ChannelOption func(c *channel)

// WithChannelSize specifies the Go chan size that is used to buffer incoming messages.
//
// The default is 100 messages.
func WithChannelSize(size int) ChannelOption ***REMOVED***
	return func(c *channel) ***REMOVED***
		c.chanSize = size
	***REMOVED***
***REMOVED***

// WithChannelHealthCheckInterval specifies the health check interval.
// PubSub will ping Redis Server if it does not receive any messages within the interval.
// To disable health check, use zero interval.
//
// The default is 3 seconds.
func WithChannelHealthCheckInterval(d time.Duration) ChannelOption ***REMOVED***
	return func(c *channel) ***REMOVED***
		c.checkInterval = d
	***REMOVED***
***REMOVED***

// WithChannelSendTimeout specifies the channel send timeout after which
// the message is dropped.
//
// The default is 60 seconds.
func WithChannelSendTimeout(d time.Duration) ChannelOption ***REMOVED***
	return func(c *channel) ***REMOVED***
		c.chanSendTimeout = d
	***REMOVED***
***REMOVED***

type channel struct ***REMOVED***
	pubSub *PubSub

	msgCh chan *Message
	allCh chan interface***REMOVED******REMOVED***
	ping  chan struct***REMOVED******REMOVED***

	chanSize        int
	chanSendTimeout time.Duration
	checkInterval   time.Duration
***REMOVED***

func newChannel(pubSub *PubSub, opts ...ChannelOption) *channel ***REMOVED***
	c := &channel***REMOVED***
		pubSub: pubSub,

		chanSize:        100,
		chanSendTimeout: time.Minute,
		checkInterval:   3 * time.Second,
	***REMOVED***
	for _, opt := range opts ***REMOVED***
		opt(c)
	***REMOVED***
	if c.checkInterval > 0 ***REMOVED***
		c.initHealthCheck()
	***REMOVED***
	return c
***REMOVED***

func (c *channel) initHealthCheck() ***REMOVED***
	ctx := context.TODO()
	c.ping = make(chan struct***REMOVED******REMOVED***, 1)

	go func() ***REMOVED***
		timer := time.NewTimer(time.Minute)
		timer.Stop()

		for ***REMOVED***
			timer.Reset(c.checkInterval)
			select ***REMOVED***
			case <-c.ping:
				if !timer.Stop() ***REMOVED***
					<-timer.C
				***REMOVED***
			case <-timer.C:
				if pingErr := c.pubSub.Ping(ctx); pingErr != nil ***REMOVED***
					c.pubSub.mu.Lock()
					c.pubSub.reconnect(ctx, pingErr)
					c.pubSub.mu.Unlock()
				***REMOVED***
			case <-c.pubSub.exit:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

// initMsgChan must be in sync with initAllChan.
func (c *channel) initMsgChan() ***REMOVED***
	ctx := context.TODO()
	c.msgCh = make(chan *Message, c.chanSize)

	go func() ***REMOVED***
		timer := time.NewTimer(time.Minute)
		timer.Stop()

		var errCount int
		for ***REMOVED***
			msg, err := c.pubSub.Receive(ctx)
			if err != nil ***REMOVED***
				if err == pool.ErrClosed ***REMOVED***
					close(c.msgCh)
					return
				***REMOVED***
				if errCount > 0 ***REMOVED***
					time.Sleep(100 * time.Millisecond)
				***REMOVED***
				errCount++
				continue
			***REMOVED***

			errCount = 0

			// Any message is as good as a ping.
			select ***REMOVED***
			case c.ping <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			default:
			***REMOVED***

			switch msg := msg.(type) ***REMOVED***
			case *Subscription:
				// Ignore.
			case *Pong:
				// Ignore.
			case *Message:
				timer.Reset(c.chanSendTimeout)
				select ***REMOVED***
				case c.msgCh <- msg:
					if !timer.Stop() ***REMOVED***
						<-timer.C
					***REMOVED***
				case <-timer.C:
					internal.Logger.Printf(
						ctx, "redis: %s channel is full for %s (message is dropped)",
						c, c.chanSendTimeout)
				***REMOVED***
			default:
				internal.Logger.Printf(ctx, "redis: unknown message type: %T", msg)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

// initAllChan must be in sync with initMsgChan.
func (c *channel) initAllChan() ***REMOVED***
	ctx := context.TODO()
	c.allCh = make(chan interface***REMOVED******REMOVED***, c.chanSize)

	go func() ***REMOVED***
		timer := time.NewTimer(time.Minute)
		timer.Stop()

		var errCount int
		for ***REMOVED***
			msg, err := c.pubSub.Receive(ctx)
			if err != nil ***REMOVED***
				if err == pool.ErrClosed ***REMOVED***
					close(c.allCh)
					return
				***REMOVED***
				if errCount > 0 ***REMOVED***
					time.Sleep(100 * time.Millisecond)
				***REMOVED***
				errCount++
				continue
			***REMOVED***

			errCount = 0

			// Any message is as good as a ping.
			select ***REMOVED***
			case c.ping <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
			default:
			***REMOVED***

			switch msg := msg.(type) ***REMOVED***
			case *Pong:
				// Ignore.
			case *Subscription, *Message:
				timer.Reset(c.chanSendTimeout)
				select ***REMOVED***
				case c.allCh <- msg:
					if !timer.Stop() ***REMOVED***
						<-timer.C
					***REMOVED***
				case <-timer.C:
					internal.Logger.Printf(
						ctx, "redis: %s channel is full for %s (message is dropped)",
						c, c.chanSendTimeout)
				***REMOVED***
			default:
				internal.Logger.Printf(ctx, "redis: unknown message type: %T", msg)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***
