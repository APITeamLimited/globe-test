// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errgroup provides synchronization, error propagation, and Context
// cancelation for groups of goroutines working on subtasks of a common task.
package errgroup

import (
	"context"
	"sync"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel on error.
type Group struct ***REMOVED***
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
***REMOVED***

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	return &Group***REMOVED***cancel: cancel***REMOVED***, ctx
***REMOVED***

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error ***REMOVED***
	g.wg.Wait()
	if g.cancel != nil ***REMOVED***
		g.cancel()
	***REMOVED***
	return g.err
***REMOVED***

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(f func() error) ***REMOVED***
	g.wg.Add(1)

	go func() ***REMOVED***
		defer g.wg.Done()

		if err := f(); err != nil ***REMOVED***
			g.errOnce.Do(func() ***REMOVED***
				g.err = err
				if g.cancel != nil ***REMOVED***
					g.cancel()
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***()
***REMOVED***
