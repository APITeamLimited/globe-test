// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.7

package context

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// An emptyCtx is never canceled, has no values, and has no deadline. It is not
// struct***REMOVED******REMOVED***, since vars of this type must have distinct addresses.
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) ***REMOVED***
	return
***REMOVED***

func (*emptyCtx) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

func (*emptyCtx) Err() error ***REMOVED***
	return nil
***REMOVED***

func (*emptyCtx) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

func (e *emptyCtx) String() string ***REMOVED***
	switch e ***REMOVED***
	case background:
		return "context.Background"
	case todo:
		return "context.TODO"
	***REMOVED***
	return "unknown empty Context"
***REMOVED***

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

// Canceled is the error returned by Context.Err when the context is canceled.
var Canceled = errors.New("context canceled")

// DeadlineExceeded is the error returned by Context.Err when the context's
// deadline passes.
var DeadlineExceeded = errors.New("context deadline exceeded")

// WithCancel returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called
// or when the parent context's Done channel is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) ***REMOVED***
	c := newCancelCtx(parent)
	propagateCancel(parent, c)
	return c, func() ***REMOVED*** c.cancel(true, Canceled) ***REMOVED***
***REMOVED***

// newCancelCtx returns an initialized cancelCtx.
func newCancelCtx(parent Context) *cancelCtx ***REMOVED***
	return &cancelCtx***REMOVED***
		Context: parent,
		done:    make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent Context, child canceler) ***REMOVED***
	if parent.Done() == nil ***REMOVED***
		return // parent is never canceled
	***REMOVED***
	if p, ok := parentCancelCtx(parent); ok ***REMOVED***
		p.mu.Lock()
		if p.err != nil ***REMOVED***
			// parent has already been canceled
			child.cancel(false, p.err)
		***REMOVED*** else ***REMOVED***
			if p.children == nil ***REMOVED***
				p.children = make(map[canceler]bool)
			***REMOVED***
			p.children[child] = true
		***REMOVED***
		p.mu.Unlock()
	***REMOVED*** else ***REMOVED***
		go func() ***REMOVED***
			select ***REMOVED***
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			***REMOVED***
		***REMOVED***()
	***REMOVED***
***REMOVED***

// parentCancelCtx follows a chain of parent references until it finds a
// *cancelCtx. This function understands how each of the concrete types in this
// package represents its parent.
func parentCancelCtx(parent Context) (*cancelCtx, bool) ***REMOVED***
	for ***REMOVED***
		switch c := parent.(type) ***REMOVED***
		case *cancelCtx:
			return c, true
		case *timerCtx:
			return c.cancelCtx, true
		case *valueCtx:
			parent = c.Context
		default:
			return nil, false
		***REMOVED***
	***REMOVED***
***REMOVED***

// removeChild removes a context from its parent.
func removeChild(parent Context, child canceler) ***REMOVED***
	p, ok := parentCancelCtx(parent)
	if !ok ***REMOVED***
		return
	***REMOVED***
	p.mu.Lock()
	if p.children != nil ***REMOVED***
		delete(p.children, child)
	***REMOVED***
	p.mu.Unlock()
***REMOVED***

// A canceler is a context type that can be canceled directly. The
// implementations are *cancelCtx and *timerCtx.
type canceler interface ***REMOVED***
	cancel(removeFromParent bool, err error)
	Done() <-chan struct***REMOVED******REMOVED***
***REMOVED***

// A cancelCtx can be canceled. When canceled, it also cancels any children
// that implement canceler.
type cancelCtx struct ***REMOVED***
	Context

	done chan struct***REMOVED******REMOVED*** // closed by the first cancel call.

	mu       sync.Mutex
	children map[canceler]bool // set to nil by the first cancel call
	err      error             // set to non-nil by the first cancel call
***REMOVED***

func (c *cancelCtx) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return c.done
***REMOVED***

func (c *cancelCtx) Err() error ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
***REMOVED***

func (c *cancelCtx) String() string ***REMOVED***
	return fmt.Sprintf("%v.WithCancel", c.Context)
***REMOVED***

// cancel closes c.done, cancels each of c's children, and, if
// removeFromParent is true, removes c from its parent's children.
func (c *cancelCtx) cancel(removeFromParent bool, err error) ***REMOVED***
	if err == nil ***REMOVED***
		panic("context: internal error: missing cancel error")
	***REMOVED***
	c.mu.Lock()
	if c.err != nil ***REMOVED***
		c.mu.Unlock()
		return // already canceled
	***REMOVED***
	c.err = err
	close(c.done)
	for child := range c.children ***REMOVED***
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err)
	***REMOVED***
	c.children = nil
	c.mu.Unlock()

	if removeFromParent ***REMOVED***
		removeChild(c.Context, c)
	***REMOVED***
***REMOVED***

// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d. If the parent's deadline is already earlier than d,
// WithDeadline(parent, d) is semantically equivalent to parent. The returned
// context's Done channel is closed when the deadline expires, when the returned
// cancel function is called, or when the parent context's Done channel is
// closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) ***REMOVED***
	if cur, ok := parent.Deadline(); ok && cur.Before(deadline) ***REMOVED***
		// The current deadline is already sooner than the new one.
		return WithCancel(parent)
	***REMOVED***
	c := &timerCtx***REMOVED***
		cancelCtx: newCancelCtx(parent),
		deadline:  deadline,
	***REMOVED***
	propagateCancel(parent, c)
	d := deadline.Sub(time.Now())
	if d <= 0 ***REMOVED***
		c.cancel(true, DeadlineExceeded) // deadline has already passed
		return c, func() ***REMOVED*** c.cancel(true, Canceled) ***REMOVED***
	***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil ***REMOVED***
		c.timer = time.AfterFunc(d, func() ***REMOVED***
			c.cancel(true, DeadlineExceeded)
		***REMOVED***)
	***REMOVED***
	return c, func() ***REMOVED*** c.cancel(true, Canceled) ***REMOVED***
***REMOVED***

// A timerCtx carries a timer and a deadline. It embeds a cancelCtx to
// implement Done and Err. It implements cancel by stopping its timer then
// delegating to cancelCtx.cancel.
type timerCtx struct ***REMOVED***
	*cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
***REMOVED***

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) ***REMOVED***
	return c.deadline, true
***REMOVED***

func (c *timerCtx) String() string ***REMOVED***
	return fmt.Sprintf("%v.WithDeadline(%s [%s])", c.cancelCtx.Context, c.deadline, c.deadline.Sub(time.Now()))
***REMOVED***

func (c *timerCtx) cancel(removeFromParent bool, err error) ***REMOVED***
	c.cancelCtx.cancel(false, err)
	if removeFromParent ***REMOVED***
		// Remove this timerCtx from its parent cancelCtx's children.
		removeChild(c.cancelCtx.Context, c)
	***REMOVED***
	c.mu.Lock()
	if c.timer != nil ***REMOVED***
		c.timer.Stop()
		c.timer = nil
	***REMOVED***
	c.mu.Unlock()
***REMOVED***

// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete:
//
// 	func slowOperationWithTimeout(ctx context.Context) (Result, error) ***REMOVED***
// 		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
// 		defer cancel()  // releases resources if slowOperation completes before timeout elapses
// 		return slowOperation(ctx)
// 	***REMOVED***
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) ***REMOVED***
	return WithDeadline(parent, time.Now().Add(timeout))
***REMOVED***

// WithValue returns a copy of parent in which the value associated with key is
// val.
//
// Use context Values only for request-scoped data that transits processes and
// APIs, not for passing optional parameters to functions.
func WithValue(parent Context, key interface***REMOVED******REMOVED***, val interface***REMOVED******REMOVED***) Context ***REMOVED***
	return &valueCtx***REMOVED***parent, key, val***REMOVED***
***REMOVED***

// A valueCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valueCtx struct ***REMOVED***
	Context
	key, val interface***REMOVED******REMOVED***
***REMOVED***

func (c *valueCtx) String() string ***REMOVED***
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", c.Context, c.key, c.val)
***REMOVED***

func (c *valueCtx) Value(key interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if c.key == key ***REMOVED***
		return c.val
	***REMOVED***
	return c.Context.Value(key)
***REMOVED***
