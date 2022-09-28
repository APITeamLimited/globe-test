// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build zos && s390x
// +build zos,s390x

package unix

import (
	"sync"
)

// This file simulates epoll on z/OS using poll.

// Analogous to epoll_event on Linux.
// TODO(neeilan): Pad is because the Linux kernel expects a 96-bit struct. We never pass this to the kernel; remove?
type EpollEvent struct ***REMOVED***
	Events uint32
	Fd     int32
	Pad    int32
***REMOVED***

const (
	EPOLLERR      = 0x8
	EPOLLHUP      = 0x10
	EPOLLIN       = 0x1
	EPOLLMSG      = 0x400
	EPOLLOUT      = 0x4
	EPOLLPRI      = 0x2
	EPOLLRDBAND   = 0x80
	EPOLLRDNORM   = 0x40
	EPOLLWRBAND   = 0x200
	EPOLLWRNORM   = 0x100
	EPOLL_CTL_ADD = 0x1
	EPOLL_CTL_DEL = 0x2
	EPOLL_CTL_MOD = 0x3
	// The following constants are part of the epoll API, but represent
	// currently unsupported functionality on z/OS.
	// EPOLL_CLOEXEC  = 0x80000
	// EPOLLET        = 0x80000000
	// EPOLLONESHOT   = 0x40000000
	// EPOLLRDHUP     = 0x2000     // Typically used with edge-triggered notis
	// EPOLLEXCLUSIVE = 0x10000000 // Exclusive wake-up mode
	// EPOLLWAKEUP    = 0x20000000 // Relies on Linux's BLOCK_SUSPEND capability
)

// TODO(neeilan): We can eliminate these epToPoll / pToEpoll calls by using identical mask values for POLL/EPOLL
// constants where possible The lower 16 bits of epoll events (uint32) can fit any system poll event (int16).

// epToPollEvt converts epoll event field to poll equivalent.
// In epoll, Events is a 32-bit field, while poll uses 16 bits.
func epToPollEvt(events uint32) int16 ***REMOVED***
	var ep2p = map[uint32]int16***REMOVED***
		EPOLLIN:  POLLIN,
		EPOLLOUT: POLLOUT,
		EPOLLHUP: POLLHUP,
		EPOLLPRI: POLLPRI,
		EPOLLERR: POLLERR,
	***REMOVED***

	var pollEvts int16 = 0
	for epEvt, pEvt := range ep2p ***REMOVED***
		if (events & epEvt) != 0 ***REMOVED***
			pollEvts |= pEvt
		***REMOVED***
	***REMOVED***

	return pollEvts
***REMOVED***

// pToEpollEvt converts 16 bit poll event bitfields to 32-bit epoll event fields.
func pToEpollEvt(revents int16) uint32 ***REMOVED***
	var p2ep = map[int16]uint32***REMOVED***
		POLLIN:  EPOLLIN,
		POLLOUT: EPOLLOUT,
		POLLHUP: EPOLLHUP,
		POLLPRI: EPOLLPRI,
		POLLERR: EPOLLERR,
	***REMOVED***

	var epollEvts uint32 = 0
	for pEvt, epEvt := range p2ep ***REMOVED***
		if (revents & pEvt) != 0 ***REMOVED***
			epollEvts |= epEvt
		***REMOVED***
	***REMOVED***

	return epollEvts
***REMOVED***

// Per-process epoll implementation.
type epollImpl struct ***REMOVED***
	mu       sync.Mutex
	epfd2ep  map[int]*eventPoll
	nextEpfd int
***REMOVED***

// eventPoll holds a set of file descriptors being watched by the process. A process can have multiple epoll instances.
// On Linux, this is an in-kernel data structure accessed through a fd.
type eventPoll struct ***REMOVED***
	mu  sync.Mutex
	fds map[int]*EpollEvent
***REMOVED***

// epoll impl for this process.
var impl epollImpl = epollImpl***REMOVED***
	epfd2ep:  make(map[int]*eventPoll),
	nextEpfd: 0,
***REMOVED***

func (e *epollImpl) epollcreate(size int) (epfd int, err error) ***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()
	epfd = e.nextEpfd
	e.nextEpfd++

	e.epfd2ep[epfd] = &eventPoll***REMOVED***
		fds: make(map[int]*EpollEvent),
	***REMOVED***
	return epfd, nil
***REMOVED***

func (e *epollImpl) epollcreate1(flag int) (fd int, err error) ***REMOVED***
	return e.epollcreate(4)
***REMOVED***

func (e *epollImpl) epollctl(epfd int, op int, fd int, event *EpollEvent) (err error) ***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()

	ep, ok := e.epfd2ep[epfd]
	if !ok ***REMOVED***

		return EBADF
	***REMOVED***

	switch op ***REMOVED***
	case EPOLL_CTL_ADD:
		// TODO(neeilan): When we make epfds and fds disjoint, detect epoll
		// loops here (instances watching each other) and return ELOOP.
		if _, ok := ep.fds[fd]; ok ***REMOVED***
			return EEXIST
		***REMOVED***
		ep.fds[fd] = event
	case EPOLL_CTL_MOD:
		if _, ok := ep.fds[fd]; !ok ***REMOVED***
			return ENOENT
		***REMOVED***
		ep.fds[fd] = event
	case EPOLL_CTL_DEL:
		if _, ok := ep.fds[fd]; !ok ***REMOVED***
			return ENOENT
		***REMOVED***
		delete(ep.fds, fd)

	***REMOVED***
	return nil
***REMOVED***

// Must be called while holding ep.mu
func (ep *eventPoll) getFds() []int ***REMOVED***
	fds := make([]int, len(ep.fds))
	for fd := range ep.fds ***REMOVED***
		fds = append(fds, fd)
	***REMOVED***
	return fds
***REMOVED***

func (e *epollImpl) epollwait(epfd int, events []EpollEvent, msec int) (n int, err error) ***REMOVED***
	e.mu.Lock() // in [rare] case of concurrent epollcreate + epollwait
	ep, ok := e.epfd2ep[epfd]

	if !ok ***REMOVED***
		e.mu.Unlock()
		return 0, EBADF
	***REMOVED***

	pollfds := make([]PollFd, 4)
	for fd, epollevt := range ep.fds ***REMOVED***
		pollfds = append(pollfds, PollFd***REMOVED***Fd: int32(fd), Events: epToPollEvt(epollevt.Events)***REMOVED***)
	***REMOVED***
	e.mu.Unlock()

	n, err = Poll(pollfds, msec)
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***

	i := 0
	for _, pFd := range pollfds ***REMOVED***
		if pFd.Revents != 0 ***REMOVED***
			events[i] = EpollEvent***REMOVED***Fd: pFd.Fd, Events: pToEpollEvt(pFd.Revents)***REMOVED***
			i++
		***REMOVED***

		if i == n ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return n, nil
***REMOVED***

func EpollCreate(size int) (fd int, err error) ***REMOVED***
	return impl.epollcreate(size)
***REMOVED***

func EpollCreate1(flag int) (fd int, err error) ***REMOVED***
	return impl.epollcreate1(flag)
***REMOVED***

func EpollCtl(epfd int, op int, fd int, event *EpollEvent) (err error) ***REMOVED***
	return impl.epollctl(epfd, op, fd, event)
***REMOVED***

// Because EpollWait mutates events, the caller is expected to coordinate
// concurrent access if calling with the same epfd from multiple goroutines.
func EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error) ***REMOVED***
	return impl.epollwait(epfd, events, msec)
***REMOVED***
