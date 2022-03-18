/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpc

import (
	"context"
	"io"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/status"
)

// pickerWrapper is a wrapper of balancer.Picker. It blocks on certain pick
// actions and unblock when there's a picker update.
type pickerWrapper struct ***REMOVED***
	mu         sync.Mutex
	done       bool
	blockingCh chan struct***REMOVED******REMOVED***
	picker     balancer.Picker
***REMOVED***

func newPickerWrapper() *pickerWrapper ***REMOVED***
	return &pickerWrapper***REMOVED***blockingCh: make(chan struct***REMOVED******REMOVED***)***REMOVED***
***REMOVED***

// updatePicker is called by UpdateBalancerState. It unblocks all blocked pick.
func (pw *pickerWrapper) updatePicker(p balancer.Picker) ***REMOVED***
	pw.mu.Lock()
	if pw.done ***REMOVED***
		pw.mu.Unlock()
		return
	***REMOVED***
	pw.picker = p
	// pw.blockingCh should never be nil.
	close(pw.blockingCh)
	pw.blockingCh = make(chan struct***REMOVED******REMOVED***)
	pw.mu.Unlock()
***REMOVED***

func doneChannelzWrapper(acw *acBalancerWrapper, done func(balancer.DoneInfo)) func(balancer.DoneInfo) ***REMOVED***
	acw.mu.Lock()
	ac := acw.ac
	acw.mu.Unlock()
	ac.incrCallsStarted()
	return func(b balancer.DoneInfo) ***REMOVED***
		if b.Err != nil && b.Err != io.EOF ***REMOVED***
			ac.incrCallsFailed()
		***REMOVED*** else ***REMOVED***
			ac.incrCallsSucceeded()
		***REMOVED***
		if done != nil ***REMOVED***
			done(b)
		***REMOVED***
	***REMOVED***
***REMOVED***

// pick returns the transport that will be used for the RPC.
// It may block in the following cases:
// - there's no picker
// - the current picker returns ErrNoSubConnAvailable
// - the current picker returns other errors and failfast is false.
// - the subConn returned by the current picker is not READY
// When one of these situations happens, pick blocks until the picker gets updated.
func (pw *pickerWrapper) pick(ctx context.Context, failfast bool, info balancer.PickInfo) (transport.ClientTransport, func(balancer.DoneInfo), error) ***REMOVED***
	var ch chan struct***REMOVED******REMOVED***

	var lastPickErr error
	for ***REMOVED***
		pw.mu.Lock()
		if pw.done ***REMOVED***
			pw.mu.Unlock()
			return nil, nil, ErrClientConnClosing
		***REMOVED***

		if pw.picker == nil ***REMOVED***
			ch = pw.blockingCh
		***REMOVED***
		if ch == pw.blockingCh ***REMOVED***
			// This could happen when either:
			// - pw.picker is nil (the previous if condition), or
			// - has called pick on the current picker.
			pw.mu.Unlock()
			select ***REMOVED***
			case <-ctx.Done():
				var errStr string
				if lastPickErr != nil ***REMOVED***
					errStr = "latest balancer error: " + lastPickErr.Error()
				***REMOVED*** else ***REMOVED***
					errStr = ctx.Err().Error()
				***REMOVED***
				switch ctx.Err() ***REMOVED***
				case context.DeadlineExceeded:
					return nil, nil, status.Error(codes.DeadlineExceeded, errStr)
				case context.Canceled:
					return nil, nil, status.Error(codes.Canceled, errStr)
				***REMOVED***
			case <-ch:
			***REMOVED***
			continue
		***REMOVED***

		ch = pw.blockingCh
		p := pw.picker
		pw.mu.Unlock()

		pickResult, err := p.Pick(info)

		if err != nil ***REMOVED***
			if err == balancer.ErrNoSubConnAvailable ***REMOVED***
				continue
			***REMOVED***
			if _, ok := status.FromError(err); ok ***REMOVED***
				// Status error: end the RPC unconditionally with this status.
				return nil, nil, err
			***REMOVED***
			// For all other errors, wait for ready RPCs should block and other
			// RPCs should fail with unavailable.
			if !failfast ***REMOVED***
				lastPickErr = err
				continue
			***REMOVED***
			return nil, nil, status.Error(codes.Unavailable, err.Error())
		***REMOVED***

		acw, ok := pickResult.SubConn.(*acBalancerWrapper)
		if !ok ***REMOVED***
			logger.Errorf("subconn returned from pick is type %T, not *acBalancerWrapper", pickResult.SubConn)
			continue
		***REMOVED***
		if t := acw.getAddrConn().getReadyTransport(); t != nil ***REMOVED***
			if channelz.IsOn() ***REMOVED***
				return t, doneChannelzWrapper(acw, pickResult.Done), nil
			***REMOVED***
			return t, pickResult.Done, nil
		***REMOVED***
		if pickResult.Done != nil ***REMOVED***
			// Calling done with nil error, no bytes sent and no bytes received.
			// DoneInfo with default value works.
			pickResult.Done(balancer.DoneInfo***REMOVED******REMOVED***)
		***REMOVED***
		logger.Infof("blockingPicker: the picked transport is not ready, loop back to repick")
		// If ok == false, ac.state is not READY.
		// A valid picker always returns READY subConn. This means the state of ac
		// just changed, and picker will be updated shortly.
		// continue back to the beginning of the for loop to repick.
	***REMOVED***
***REMOVED***

func (pw *pickerWrapper) close() ***REMOVED***
	pw.mu.Lock()
	defer pw.mu.Unlock()
	if pw.done ***REMOVED***
		return
	***REMOVED***
	pw.done = true
	close(pw.blockingCh)
***REMOVED***
