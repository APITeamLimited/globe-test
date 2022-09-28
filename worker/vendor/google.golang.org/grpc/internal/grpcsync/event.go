/*
 *
 * Copyright 2018 gRPC authors.
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

// Package grpcsync implements additional synchronization primitives built upon
// the sync package.
package grpcsync

import (
	"sync"
	"sync/atomic"
)

// Event represents a one-time event that may occur in the future.
type Event struct ***REMOVED***
	fired int32
	c     chan struct***REMOVED******REMOVED***
	o     sync.Once
***REMOVED***

// Fire causes e to complete.  It is safe to call multiple times, and
// concurrently.  It returns true iff this call to Fire caused the signaling
// channel returned by Done to close.
func (e *Event) Fire() bool ***REMOVED***
	ret := false
	e.o.Do(func() ***REMOVED***
		atomic.StoreInt32(&e.fired, 1)
		close(e.c)
		ret = true
	***REMOVED***)
	return ret
***REMOVED***

// Done returns a channel that will be closed when Fire is called.
func (e *Event) Done() <-chan struct***REMOVED******REMOVED*** ***REMOVED***
	return e.c
***REMOVED***

// HasFired returns true if Fire has been called.
func (e *Event) HasFired() bool ***REMOVED***
	return atomic.LoadInt32(&e.fired) == 1
***REMOVED***

// NewEvent returns a new, ready-to-use Event.
func NewEvent() *Event ***REMOVED***
	return &Event***REMOVED***c: make(chan struct***REMOVED******REMOVED***)***REMOVED***
***REMOVED***
