// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"context"
	"crypto"
	"sync"
	"time"
)

// renewJitter is the maximum deviation from Manager.RenewBefore.
const renewJitter = time.Hour

// domainRenewal tracks the state used by the periodic timers
// renewing a single domain's cert.
type domainRenewal struct ***REMOVED***
	m      *Manager
	domain string
	key    crypto.Signer

	timerMu sync.Mutex
	timer   *time.Timer
***REMOVED***

// start starts a cert renewal timer at the time
// defined by the certificate expiration time exp.
//
// If the timer is already started, calling start is a noop.
func (dr *domainRenewal) start(exp time.Time) ***REMOVED***
	dr.timerMu.Lock()
	defer dr.timerMu.Unlock()
	if dr.timer != nil ***REMOVED***
		return
	***REMOVED***
	dr.timer = time.AfterFunc(dr.next(exp), dr.renew)
***REMOVED***

// stop stops the cert renewal timer.
// If the timer is already stopped, calling stop is a noop.
func (dr *domainRenewal) stop() ***REMOVED***
	dr.timerMu.Lock()
	defer dr.timerMu.Unlock()
	if dr.timer == nil ***REMOVED***
		return
	***REMOVED***
	dr.timer.Stop()
	dr.timer = nil
***REMOVED***

// renew is called periodically by a timer.
// The first renew call is kicked off by dr.start.
func (dr *domainRenewal) renew() ***REMOVED***
	dr.timerMu.Lock()
	defer dr.timerMu.Unlock()
	if dr.timer == nil ***REMOVED***
		return
	***REMOVED***

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	// TODO: rotate dr.key at some point?
	next, err := dr.do(ctx)
	if err != nil ***REMOVED***
		next = renewJitter / 2
		next += time.Duration(pseudoRand.int63n(int64(next)))
	***REMOVED***
	dr.timer = time.AfterFunc(next, dr.renew)
	testDidRenewLoop(next, err)
***REMOVED***

// do is similar to Manager.createCert but it doesn't lock a Manager.state item.
// Instead, it requests a new certificate independently and, upon success,
// replaces dr.m.state item with a new one and updates cache for the given domain.
//
// It may return immediately if the expiration date of the currently cached cert
// is far enough in the future.
//
// The returned value is a time interval after which the renewal should occur again.
func (dr *domainRenewal) do(ctx context.Context) (time.Duration, error) ***REMOVED***
	// a race is likely unavoidable in a distributed environment
	// but we try nonetheless
	if tlscert, err := dr.m.cacheGet(ctx, dr.domain); err == nil ***REMOVED***
		next := dr.next(tlscert.Leaf.NotAfter)
		if next > dr.m.renewBefore()+renewJitter ***REMOVED***
			return next, nil
		***REMOVED***
	***REMOVED***

	der, leaf, err := dr.m.authorizedCert(ctx, dr.key, dr.domain)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	state := &certState***REMOVED***
		key:  dr.key,
		cert: der,
		leaf: leaf,
	***REMOVED***
	tlscert, err := state.tlscert()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	dr.m.cachePut(ctx, dr.domain, tlscert)
	dr.m.stateMu.Lock()
	defer dr.m.stateMu.Unlock()
	// m.state is guaranteed to be non-nil at this point
	dr.m.state[dr.domain] = state
	return dr.next(leaf.NotAfter), nil
***REMOVED***

func (dr *domainRenewal) next(expiry time.Time) time.Duration ***REMOVED***
	d := expiry.Sub(timeNow()) - dr.m.renewBefore()
	// add a bit of randomness to renew deadline
	n := pseudoRand.int63n(int64(renewJitter))
	d -= time.Duration(n)
	if d < 0 ***REMOVED***
		return 0
	***REMOVED***
	return d
***REMOVED***

var testDidRenewLoop = func(next time.Duration, err error) ***REMOVED******REMOVED***
