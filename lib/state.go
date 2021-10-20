/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/stats"
)

// DialContexter is an interface that can dial with a context
type DialContexter interface ***REMOVED***
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
***REMOVED***

// State provides the volatile state for a VU.
type State struct ***REMOVED***
	// Global options.
	Options Options

	// Logger. Avoid using the global logger.
	// TODO change to logrus.FieldLogger when there is time to fix all the tests
	Logger *logrus.Logger

	// Current group; all emitted metrics are tagged with this.
	Group *Group

	// Networking equipment.
	Transport http.RoundTripper
	Dialer    DialContexter
	CookieJar *cookiejar.Jar
	TLSConfig *tls.Config

	// Rate limits.
	RPSLimit *rate.Limiter

	// Sample channel, possibly buffered
	Samples chan<- stats.SampleContainer

	// Buffer pool; use instead of allocating fresh buffers when possible.
	// TODO: maybe use https://golang.org/pkg/sync/#Pool ?
	BPool *bpool.BufferPool

	VUID, VUIDGlobal uint64
	Iteration        int64
	Tags             *TagMap
	// These will be assigned on VU activation.
	// Returns the iteration number of this VU in the current scenario.
	GetScenarioVUIter func() uint64
	// Returns the iteration number across all VUs in the current scenario
	// unique to this single k6 instance.
	// TODO: Maybe this doesn't belong here but in ScenarioState?
	GetScenarioLocalVUIter func() uint64
	// Returns the iteration number across all VUs in the current scenario
	// unique globally across k6 instances (taking into account execution
	// segments).
	GetScenarioGlobalVUIter func() uint64

	BuiltinMetrics *metrics.BuiltinMetrics
***REMOVED***

// CloneTags makes a copy of the tags map and returns it.
func (s *State) CloneTags() map[string]string ***REMOVED***
	return s.Tags.Clone()
***REMOVED***

// TagMap is a safe-concurrent Tags lookup.
type TagMap struct ***REMOVED***
	m     map[string]string
	mutex sync.RWMutex
***REMOVED***

// NewTagMap creates a TagMap,
// if a not-nil map is passed then it will be used as the internal map
// otherwise a new one will be created.
func NewTagMap(m map[string]string) *TagMap ***REMOVED***
	if m == nil ***REMOVED***
		m = make(map[string]string)
	***REMOVED***
	return &TagMap***REMOVED***
		m:     m,
		mutex: sync.RWMutex***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// Set sets a Tag.
func (tg *TagMap) Set(k, v string) ***REMOVED***
	tg.mutex.Lock()
	defer tg.mutex.Unlock()
	tg.m[k] = v
***REMOVED***

// Get returns the Tag value and true
// if the provided key has been found.
func (tg *TagMap) Get(k string) (string, bool) ***REMOVED***
	tg.mutex.RLock()
	defer tg.mutex.RUnlock()
	v, ok := tg.m[k]
	return v, ok
***REMOVED***

// Len returns the number of the set keys.
func (tg *TagMap) Len() int ***REMOVED***
	tg.mutex.RLock()
	defer tg.mutex.RUnlock()
	return len(tg.m)
***REMOVED***

// Delete deletes a map's item based on the provided key.
func (tg *TagMap) Delete(k string) ***REMOVED***
	tg.mutex.Lock()
	defer tg.mutex.Unlock()
	delete(tg.m, k)
***REMOVED***

// Clone returns a map with the entire set of items.
func (tg *TagMap) Clone() map[string]string ***REMOVED***
	tg.mutex.RLock()
	defer tg.mutex.RUnlock()

	tags := make(map[string]string, len(tg.m))
	for k, v := range tg.m ***REMOVED***
		tags[k] = v
	***REMOVED***
	return tags
***REMOVED***
