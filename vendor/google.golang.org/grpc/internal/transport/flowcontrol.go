/*
 *
 * Copyright 2014 gRPC authors.
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

package transport

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

// writeQuota is a soft limit on the amount of data a stream can
// schedule before some of it is written out.
type writeQuota struct ***REMOVED***
	quota int32
	// get waits on read from when quota goes less than or equal to zero.
	// replenish writes on it when quota goes positive again.
	ch chan struct***REMOVED******REMOVED***
	// done is triggered in error case.
	done <-chan struct***REMOVED******REMOVED***
	// replenish is called by loopyWriter to give quota back to.
	// It is implemented as a field so that it can be updated
	// by tests.
	replenish func(n int)
***REMOVED***

func newWriteQuota(sz int32, done <-chan struct***REMOVED******REMOVED***) *writeQuota ***REMOVED***
	w := &writeQuota***REMOVED***
		quota: sz,
		ch:    make(chan struct***REMOVED******REMOVED***, 1),
		done:  done,
	***REMOVED***
	w.replenish = w.realReplenish
	return w
***REMOVED***

func (w *writeQuota) get(sz int32) error ***REMOVED***
	for ***REMOVED***
		if atomic.LoadInt32(&w.quota) > 0 ***REMOVED***
			atomic.AddInt32(&w.quota, -sz)
			return nil
		***REMOVED***
		select ***REMOVED***
		case <-w.ch:
			continue
		case <-w.done:
			return errStreamDone
		***REMOVED***
	***REMOVED***
***REMOVED***

func (w *writeQuota) realReplenish(n int) ***REMOVED***
	sz := int32(n)
	a := atomic.AddInt32(&w.quota, sz)
	b := a - sz
	if b <= 0 && a > 0 ***REMOVED***
		select ***REMOVED***
		case w.ch <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***

type trInFlow struct ***REMOVED***
	limit               uint32
	unacked             uint32
	effectiveWindowSize uint32
***REMOVED***

func (f *trInFlow) newLimit(n uint32) uint32 ***REMOVED***
	d := n - f.limit
	f.limit = n
	f.updateEffectiveWindowSize()
	return d
***REMOVED***

func (f *trInFlow) onData(n uint32) uint32 ***REMOVED***
	f.unacked += n
	if f.unacked >= f.limit/4 ***REMOVED***
		w := f.unacked
		f.unacked = 0
		f.updateEffectiveWindowSize()
		return w
	***REMOVED***
	f.updateEffectiveWindowSize()
	return 0
***REMOVED***

func (f *trInFlow) reset() uint32 ***REMOVED***
	w := f.unacked
	f.unacked = 0
	f.updateEffectiveWindowSize()
	return w
***REMOVED***

func (f *trInFlow) updateEffectiveWindowSize() ***REMOVED***
	atomic.StoreUint32(&f.effectiveWindowSize, f.limit-f.unacked)
***REMOVED***

func (f *trInFlow) getSize() uint32 ***REMOVED***
	return atomic.LoadUint32(&f.effectiveWindowSize)
***REMOVED***

// TODO(mmukhi): Simplify this code.
// inFlow deals with inbound flow control
type inFlow struct ***REMOVED***
	mu sync.Mutex
	// The inbound flow control limit for pending data.
	limit uint32
	// pendingData is the overall data which have been received but not been
	// consumed by applications.
	pendingData uint32
	// The amount of data the application has consumed but grpc has not sent
	// window update for them. Used to reduce window update frequency.
	pendingUpdate uint32
	// delta is the extra window update given by receiver when an application
	// is reading data bigger in size than the inFlow limit.
	delta uint32
***REMOVED***

// newLimit updates the inflow window to a new value n.
// It assumes that n is always greater than the old limit.
func (f *inFlow) newLimit(n uint32) uint32 ***REMOVED***
	f.mu.Lock()
	d := n - f.limit
	f.limit = n
	f.mu.Unlock()
	return d
***REMOVED***

func (f *inFlow) maybeAdjust(n uint32) uint32 ***REMOVED***
	if n > uint32(math.MaxInt32) ***REMOVED***
		n = uint32(math.MaxInt32)
	***REMOVED***
	f.mu.Lock()
	defer f.mu.Unlock()
	// estSenderQuota is the receiver's view of the maximum number of bytes the sender
	// can send without a window update.
	estSenderQuota := int32(f.limit - (f.pendingData + f.pendingUpdate))
	// estUntransmittedData is the maximum number of bytes the sends might not have put
	// on the wire yet. A value of 0 or less means that we have already received all or
	// more bytes than the application is requesting to read.
	estUntransmittedData := int32(n - f.pendingData) // Casting into int32 since it could be negative.
	// This implies that unless we send a window update, the sender won't be able to send all the bytes
	// for this message. Therefore we must send an update over the limit since there's an active read
	// request from the application.
	if estUntransmittedData > estSenderQuota ***REMOVED***
		// Sender's window shouldn't go more than 2^31 - 1 as specified in the HTTP spec.
		if f.limit+n > maxWindowSize ***REMOVED***
			f.delta = maxWindowSize - f.limit
		***REMOVED*** else ***REMOVED***
			// Send a window update for the whole message and not just the difference between
			// estUntransmittedData and estSenderQuota. This will be helpful in case the message
			// is padded; We will fallback on the current available window(at least a 1/4th of the limit).
			f.delta = n
		***REMOVED***
		return f.delta
	***REMOVED***
	return 0
***REMOVED***

// onData is invoked when some data frame is received. It updates pendingData.
func (f *inFlow) onData(n uint32) error ***REMOVED***
	f.mu.Lock()
	f.pendingData += n
	if f.pendingData+f.pendingUpdate > f.limit+f.delta ***REMOVED***
		limit := f.limit
		rcvd := f.pendingData + f.pendingUpdate
		f.mu.Unlock()
		return fmt.Errorf("received %d-bytes data exceeding the limit %d bytes", rcvd, limit)
	***REMOVED***
	f.mu.Unlock()
	return nil
***REMOVED***

// onRead is invoked when the application reads the data. It returns the window size
// to be sent to the peer.
func (f *inFlow) onRead(n uint32) uint32 ***REMOVED***
	f.mu.Lock()
	if f.pendingData == 0 ***REMOVED***
		f.mu.Unlock()
		return 0
	***REMOVED***
	f.pendingData -= n
	if n > f.delta ***REMOVED***
		n -= f.delta
		f.delta = 0
	***REMOVED*** else ***REMOVED***
		f.delta -= n
		n = 0
	***REMOVED***
	f.pendingUpdate += n
	if f.pendingUpdate >= f.limit/4 ***REMOVED***
		wu := f.pendingUpdate
		f.pendingUpdate = 0
		f.mu.Unlock()
		return wu
	***REMOVED***
	f.mu.Unlock()
	return 0
***REMOVED***
