package redis

import (
	"context"
	"sync"
)

type pipelineExecer func(context.Context, []Cmder) error

// Pipeliner is an mechanism to realise Redis Pipeline technique.
//
// Pipelining is a technique to extremely speed up processing by packing
// operations to batches, send them at once to Redis and read a replies in a
// singe step.
// See https://redis.io/topics/pipelining
//
// Pay attention, that Pipeline is not a transaction, so you can get unexpected
// results in case of big pipelines and small read/write timeouts.
// Redis client has retransmission logic in case of timeouts, pipeline
// can be retransmitted and commands can be executed more then once.
// To avoid this: it is good idea to use reasonable bigger read/write timeouts
// depends of your batch size and/or use TxPipeline.
type Pipeliner interface ***REMOVED***
	StatefulCmdable
	Len() int
	Do(ctx context.Context, args ...interface***REMOVED******REMOVED***) *Cmd
	Process(ctx context.Context, cmd Cmder) error
	Discard()
	Exec(ctx context.Context) ([]Cmder, error)
***REMOVED***

var _ Pipeliner = (*Pipeline)(nil)

// Pipeline implements pipelining as described in
// http://redis.io/topics/pipelining. It's safe for concurrent use
// by multiple goroutines.
type Pipeline struct ***REMOVED***
	cmdable
	statefulCmdable

	exec pipelineExecer

	mu   sync.Mutex
	cmds []Cmder
***REMOVED***

func (c *Pipeline) init() ***REMOVED***
	c.cmdable = c.Process
	c.statefulCmdable = c.Process
***REMOVED***

// Len returns the number of queued commands.
func (c *Pipeline) Len() int ***REMOVED***
	c.mu.Lock()
	ln := len(c.cmds)
	c.mu.Unlock()
	return ln
***REMOVED***

// Do queues the custom command for later execution.
func (c *Pipeline) Do(ctx context.Context, args ...interface***REMOVED******REMOVED***) *Cmd ***REMOVED***
	cmd := NewCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Process queues the cmd for later execution.
func (c *Pipeline) Process(ctx context.Context, cmd Cmder) error ***REMOVED***
	c.mu.Lock()
	c.cmds = append(c.cmds, cmd)
	c.mu.Unlock()
	return nil
***REMOVED***

// Discard resets the pipeline and discards queued commands.
func (c *Pipeline) Discard() ***REMOVED***
	c.mu.Lock()
	c.cmds = c.cmds[:0]
	c.mu.Unlock()
***REMOVED***

// Exec executes all previously queued commands using one
// client-server roundtrip.
//
// Exec always returns list of commands and error of the first failed
// command if any.
func (c *Pipeline) Exec(ctx context.Context) ([]Cmder, error) ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.cmds) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	cmds := c.cmds
	c.cmds = nil

	return cmds, c.exec(ctx, cmds)
***REMOVED***

func (c *Pipeline) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	if err := fn(c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return c.Exec(ctx)
***REMOVED***

func (c *Pipeline) Pipeline() Pipeliner ***REMOVED***
	return c
***REMOVED***

func (c *Pipeline) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.Pipelined(ctx, fn)
***REMOVED***

func (c *Pipeline) TxPipeline() Pipeliner ***REMOVED***
	return c
***REMOVED***
