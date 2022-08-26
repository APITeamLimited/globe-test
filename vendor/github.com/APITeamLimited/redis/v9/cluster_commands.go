package redis

import (
	"context"
	"sync"
	"sync/atomic"
)

func (c *ClusterClient) DBSize(ctx context.Context) *IntCmd ***REMOVED***
	cmd := NewIntCmd(ctx, "dbsize")
	_ = c.hooks.process(ctx, cmd, func(ctx context.Context, _ Cmder) error ***REMOVED***
		var size int64
		err := c.ForEachMaster(ctx, func(ctx context.Context, master *Client) error ***REMOVED***
			n, err := master.DBSize(ctx).Result()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			atomic.AddInt64(&size, n)
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			cmd.SetErr(err)
		***REMOVED*** else ***REMOVED***
			cmd.val = size
		***REMOVED***
		return nil
	***REMOVED***)
	return cmd
***REMOVED***

func (c *ClusterClient) ScriptLoad(ctx context.Context, script string) *StringCmd ***REMOVED***
	cmd := NewStringCmd(ctx, "script", "load", script)
	_ = c.hooks.process(ctx, cmd, func(ctx context.Context, _ Cmder) error ***REMOVED***
		mu := &sync.Mutex***REMOVED******REMOVED***
		err := c.ForEachShard(ctx, func(ctx context.Context, shard *Client) error ***REMOVED***
			val, err := shard.ScriptLoad(ctx, script).Result()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			mu.Lock()
			if cmd.Val() == "" ***REMOVED***
				cmd.val = val
			***REMOVED***
			mu.Unlock()

			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			cmd.SetErr(err)
		***REMOVED***
		return nil
	***REMOVED***)
	return cmd
***REMOVED***

func (c *ClusterClient) ScriptFlush(ctx context.Context) *StatusCmd ***REMOVED***
	cmd := NewStatusCmd(ctx, "script", "flush")
	_ = c.hooks.process(ctx, cmd, func(ctx context.Context, _ Cmder) error ***REMOVED***
		err := c.ForEachShard(ctx, func(ctx context.Context, shard *Client) error ***REMOVED***
			return shard.ScriptFlush(ctx).Err()
		***REMOVED***)
		if err != nil ***REMOVED***
			cmd.SetErr(err)
		***REMOVED***
		return nil
	***REMOVED***)
	return cmd
***REMOVED***

func (c *ClusterClient) ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd ***REMOVED***
	args := make([]interface***REMOVED******REMOVED***, 2+len(hashes))
	args[0] = "script"
	args[1] = "exists"
	for i, hash := range hashes ***REMOVED***
		args[2+i] = hash
	***REMOVED***
	cmd := NewBoolSliceCmd(ctx, args...)

	result := make([]bool, len(hashes))
	for i := range result ***REMOVED***
		result[i] = true
	***REMOVED***

	_ = c.hooks.process(ctx, cmd, func(ctx context.Context, _ Cmder) error ***REMOVED***
		mu := &sync.Mutex***REMOVED******REMOVED***
		err := c.ForEachShard(ctx, func(ctx context.Context, shard *Client) error ***REMOVED***
			val, err := shard.ScriptExists(ctx, hashes...).Result()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			mu.Lock()
			for i, v := range val ***REMOVED***
				result[i] = result[i] && v
			***REMOVED***
			mu.Unlock()

			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			cmd.SetErr(err)
		***REMOVED*** else ***REMOVED***
			cmd.val = result
		***REMOVED***
		return nil
	***REMOVED***)
	return cmd
***REMOVED***
