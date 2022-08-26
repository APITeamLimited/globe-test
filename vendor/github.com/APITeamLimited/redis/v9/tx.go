package redis

import (
	"context"

	"github.com/APITeamLimited/redis/v9/internal/pool"
	"github.com/APITeamLimited/redis/v9/internal/proto"
)

// TxFailedErr transaction redis failed.
const TxFailedErr = proto.RedisError("redis: transaction failed")

// Tx implements Redis transactions as described in
// http://redis.io/topics/transactions. It's NOT safe for concurrent use
// by multiple goroutines, because Exec resets list of watched keys.
//
// If you don't need WATCH, use Pipeline instead.
type Tx struct ***REMOVED***
	baseClient
	cmdable
	statefulCmdable
	hooks
***REMOVED***

func (c *Client) newTx() *Tx ***REMOVED***
	tx := Tx***REMOVED***
		baseClient: baseClient***REMOVED***
			opt:      c.opt,
			connPool: pool.NewStickyConnPool(c.connPool),
		***REMOVED***,
		hooks: c.hooks.clone(),
	***REMOVED***
	tx.init()
	return &tx
***REMOVED***

func (c *Tx) init() ***REMOVED***
	c.cmdable = c.Process
	c.statefulCmdable = c.Process
***REMOVED***

func (c *Tx) Process(ctx context.Context, cmd Cmder) error ***REMOVED***
	return c.hooks.process(ctx, cmd, c.baseClient.process)
***REMOVED***

// Watch prepares a transaction and marks the keys to be watched
// for conditional execution if there are any keys.
//
// The transaction is automatically closed when fn exits.
func (c *Client) Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error ***REMOVED***
	tx := c.newTx()
	defer tx.Close(ctx)
	if len(keys) > 0 ***REMOVED***
		if err := tx.Watch(ctx, keys...).Err(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return fn(tx)
***REMOVED***

// Close closes the transaction, releasing any open resources.
func (c *Tx) Close(ctx context.Context) error ***REMOVED***
	_ = c.Unwatch(ctx).Err()
	return c.baseClient.Close()
***REMOVED***

// Watch marks the keys to be watched for conditional execution
// of a transaction.
func (c *Tx) Watch(ctx context.Context, keys ...string) *StatusCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "watch"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewStatusCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Unwatch flushes all the previously watched keys for a transaction.
func (c *Tx) Unwatch(ctx context.Context, keys ...string) *StatusCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 1+len(keys))
	args[0] = "unwatch"
	for i, key := range keys ***REMOVED***
		args[1+i] = key
	***REMOVED***
	cmd := NewStatusCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
***REMOVED***

// Pipeline creates a pipeline. Usually it is more convenient to use Pipelined.
func (c *Tx) Pipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: func(ctx context.Context, cmds []Cmder) error ***REMOVED***
			return c.hooks.processPipeline(ctx, cmds, c.baseClient.processPipeline)
		***REMOVED***,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***

// Pipelined executes commands queued in the fn outside of the transaction.
// Use TxPipelined if you need transactional behavior.
func (c *Tx) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.Pipeline().Pipelined(ctx, fn)
***REMOVED***

// TxPipelined executes commands queued in the fn in the transaction.
//
// When using WATCH, EXEC will execute commands only if the watched keys
// were not modified, allowing for a check-and-set mechanism.
//
// Exec always returns list of commands. If transaction fails
// TxFailedErr is returned. Otherwise Exec returns an error of the first
// failed command or nil.
func (c *Tx) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) ***REMOVED***
	return c.TxPipeline().Pipelined(ctx, fn)
***REMOVED***

// TxPipeline creates a pipeline. Usually it is more convenient to use TxPipelined.
func (c *Tx) TxPipeline() Pipeliner ***REMOVED***
	pipe := Pipeline***REMOVED***
		exec: func(ctx context.Context, cmds []Cmder) error ***REMOVED***
			return c.hooks.processTxPipeline(ctx, cmds, c.baseClient.processTxPipeline)
		***REMOVED***,
	***REMOVED***
	pipe.init()
	return &pipe
***REMOVED***
