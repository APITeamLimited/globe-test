package libWorker

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// DialContexter is an interface that can dial with a context
type DialContexter interface ***REMOVED***
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
***REMOVED***

// State provides the volatile state for a VU.
type State struct ***REMOVED***
	// Global options and built-in workerMetrics.
	//
	// TODO: remove them from here, the built-in metrics and the script options
	// are not part of a VU's unique "state", they are global and the same for
	// all VUs. Figure out how to thread them some other way, e.g. through the
	// TestPreInitState. The Samples channel might also benefit from that...
	Options        Options
	BuiltinMetrics *workerMetrics.BuiltinMetrics

	// Logger. Avoid using the global logger.
	// TODO: change to logrus.FieldLogger when there is time to fix all the tests
	Logger *logrus.Logger

	// Current group; all emitted metrics are tagged with this.
	Group *Group

	// Networking equipment.
	Dialer DialContexter

	// TODO: move a lot of the things below to the k6/http ModuleInstance, see
	// https://github.com/grafana/k6/issues/2293.
	Transport http.RoundTripper
	CookieJar *cookiejar.Jar
	TLSConfig *tls.Config

	// Rate limits.
	RPSLimit *rate.Limiter

	// Sample channel, possibly buffered
	Samples chan<- workerMetrics.SampleContainer

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
