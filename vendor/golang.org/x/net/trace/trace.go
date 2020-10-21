// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package trace implements tracing of requests and long-lived objects.
It exports HTTP interfaces on /debug/requests and /debug/events.

A trace.Trace provides tracing for short-lived objects, usually requests.
A request handler might be implemented like this:

	func fooHandler(w http.ResponseWriter, req *http.Request) ***REMOVED***
		tr := trace.New("mypkg.Foo", req.URL.Path)
		defer tr.Finish()
		...
		tr.LazyPrintf("some event %q happened", str)
		...
		if err := somethingImportant(); err != nil ***REMOVED***
			tr.LazyPrintf("somethingImportant failed: %v", err)
			tr.SetError()
		***REMOVED***
	***REMOVED***

The /debug/requests HTTP endpoint organizes the traces by family,
errors, and duration.  It also provides histogram of request duration
for each family.

A trace.EventLog provides tracing for long-lived objects, such as RPC
connections.

	// A Fetcher fetches URL paths for a single domain.
	type Fetcher struct ***REMOVED***
		domain string
		events trace.EventLog
	***REMOVED***

	func NewFetcher(domain string) *Fetcher ***REMOVED***
		return &Fetcher***REMOVED***
			domain,
			trace.NewEventLog("mypkg.Fetcher", domain),
		***REMOVED***
	***REMOVED***

	func (f *Fetcher) Fetch(path string) (string, error) ***REMOVED***
		resp, err := http.Get("http://" + f.domain + "/" + path)
		if err != nil ***REMOVED***
			f.events.Errorf("Get(%q) = %v", path, err)
			return "", err
		***REMOVED***
		f.events.Printf("Get(%q) = %s", path, resp.Status)
		...
	***REMOVED***

	func (f *Fetcher) Close() error ***REMOVED***
		f.events.Finish()
		return nil
	***REMOVED***

The /debug/events HTTP endpoint organizes the event logs by family and
by time since the last error.  The expanded view displays recent log
entries and the log's call stack.
*/
package trace // import "golang.org/x/net/trace"

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/internal/timeseries"
)

// DebugUseAfterFinish controls whether to debug uses of Trace values after finishing.
// FOR DEBUGGING ONLY. This will slow down the program.
var DebugUseAfterFinish = false

// AuthRequest determines whether a specific request is permitted to load the
// /debug/requests or /debug/events pages.
//
// It returns two bools; the first indicates whether the page may be viewed at all,
// and the second indicates whether sensitive events will be shown.
//
// AuthRequest may be replaced by a program to customize its authorization requirements.
//
// The default AuthRequest function returns (true, true) if and only if the request
// comes from localhost/127.0.0.1/[::1].
var AuthRequest = func(req *http.Request) (any, sensitive bool) ***REMOVED***
	// RemoteAddr is commonly in the form "IP" or "IP:port".
	// If it is in the form "IP:port", split off the port.
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil ***REMOVED***
		host = req.RemoteAddr
	***REMOVED***
	switch host ***REMOVED***
	case "localhost", "127.0.0.1", "::1":
		return true, true
	default:
		return false, false
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	_, pat := http.DefaultServeMux.Handler(&http.Request***REMOVED***URL: &url.URL***REMOVED***Path: "/debug/requests"***REMOVED******REMOVED***)
	if pat != "" ***REMOVED***
		panic("/debug/requests is already registered. You may have two independent copies of " +
			"golang.org/x/net/trace in your binary, trying to maintain separate state. This may " +
			"involve a vendored copy of golang.org/x/net/trace.")
	***REMOVED***

	// TODO(jbd): Serve Traces from /debug/traces in the future?
	// There is no requirement for a request to be present to have traces.
	http.HandleFunc("/debug/requests", Traces)
	http.HandleFunc("/debug/events", Events)
***REMOVED***

// NewContext returns a copy of the parent context
// and associates it with a Trace.
func NewContext(ctx context.Context, tr Trace) context.Context ***REMOVED***
	return context.WithValue(ctx, contextKey, tr)
***REMOVED***

// FromContext returns the Trace bound to the context, if any.
func FromContext(ctx context.Context) (tr Trace, ok bool) ***REMOVED***
	tr, ok = ctx.Value(contextKey).(Trace)
	return
***REMOVED***

// Traces responds with traces from the program.
// The package initialization registers it in http.DefaultServeMux
// at /debug/requests.
//
// It performs authorization by running AuthRequest.
func Traces(w http.ResponseWriter, req *http.Request) ***REMOVED***
	any, sensitive := AuthRequest(req)
	if !any ***REMOVED***
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	***REMOVED***
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	Render(w, req, sensitive)
***REMOVED***

// Events responds with a page of events collected by EventLogs.
// The package initialization registers it in http.DefaultServeMux
// at /debug/events.
//
// It performs authorization by running AuthRequest.
func Events(w http.ResponseWriter, req *http.Request) ***REMOVED***
	any, sensitive := AuthRequest(req)
	if !any ***REMOVED***
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	***REMOVED***
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	RenderEvents(w, req, sensitive)
***REMOVED***

// Render renders the HTML page typically served at /debug/requests.
// It does not do any auth checking. The request may be nil.
//
// Most users will use the Traces handler.
func Render(w io.Writer, req *http.Request, sensitive bool) ***REMOVED***
	data := &struct ***REMOVED***
		Families         []string
		ActiveTraceCount map[string]int
		CompletedTraces  map[string]*family

		// Set when a bucket has been selected.
		Traces        traceList
		Family        string
		Bucket        int
		Expanded      bool
		Traced        bool
		Active        bool
		ShowSensitive bool // whether to show sensitive events

		Histogram       template.HTML
		HistogramWindow string // e.g. "last minute", "last hour", "all time"

		// If non-zero, the set of traces is a partial set,
		// and this is the total number.
		Total int
	***REMOVED******REMOVED***
		CompletedTraces: completedTraces,
	***REMOVED***

	data.ShowSensitive = sensitive
	if req != nil ***REMOVED***
		// Allow show_sensitive=0 to force hiding of sensitive data for testing.
		// This only goes one way; you can't use show_sensitive=1 to see things.
		if req.FormValue("show_sensitive") == "0" ***REMOVED***
			data.ShowSensitive = false
		***REMOVED***

		if exp, err := strconv.ParseBool(req.FormValue("exp")); err == nil ***REMOVED***
			data.Expanded = exp
		***REMOVED***
		if exp, err := strconv.ParseBool(req.FormValue("rtraced")); err == nil ***REMOVED***
			data.Traced = exp
		***REMOVED***
	***REMOVED***

	completedMu.RLock()
	data.Families = make([]string, 0, len(completedTraces))
	for fam := range completedTraces ***REMOVED***
		data.Families = append(data.Families, fam)
	***REMOVED***
	completedMu.RUnlock()
	sort.Strings(data.Families)

	// We are careful here to minimize the time spent locking activeMu,
	// since that lock is required every time an RPC starts and finishes.
	data.ActiveTraceCount = make(map[string]int, len(data.Families))
	activeMu.RLock()
	for fam, s := range activeTraces ***REMOVED***
		data.ActiveTraceCount[fam] = s.Len()
	***REMOVED***
	activeMu.RUnlock()

	var ok bool
	data.Family, data.Bucket, ok = parseArgs(req)
	switch ***REMOVED***
	case !ok:
		// No-op
	case data.Bucket == -1:
		data.Active = true
		n := data.ActiveTraceCount[data.Family]
		data.Traces = getActiveTraces(data.Family)
		if len(data.Traces) < n ***REMOVED***
			data.Total = n
		***REMOVED***
	case data.Bucket < bucketsPerFamily:
		if b := lookupBucket(data.Family, data.Bucket); b != nil ***REMOVED***
			data.Traces = b.Copy(data.Traced)
		***REMOVED***
	default:
		if f := getFamily(data.Family, false); f != nil ***REMOVED***
			var obs timeseries.Observable
			f.LatencyMu.RLock()
			switch o := data.Bucket - bucketsPerFamily; o ***REMOVED***
			case 0:
				obs = f.Latency.Minute()
				data.HistogramWindow = "last minute"
			case 1:
				obs = f.Latency.Hour()
				data.HistogramWindow = "last hour"
			case 2:
				obs = f.Latency.Total()
				data.HistogramWindow = "all time"
			***REMOVED***
			f.LatencyMu.RUnlock()
			if obs != nil ***REMOVED***
				data.Histogram = obs.(*histogram).html()
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if data.Traces != nil ***REMOVED***
		defer data.Traces.Free()
		sort.Sort(data.Traces)
	***REMOVED***

	completedMu.RLock()
	defer completedMu.RUnlock()
	if err := pageTmpl().ExecuteTemplate(w, "Page", data); err != nil ***REMOVED***
		log.Printf("net/trace: Failed executing template: %v", err)
	***REMOVED***
***REMOVED***

func parseArgs(req *http.Request) (fam string, b int, ok bool) ***REMOVED***
	if req == nil ***REMOVED***
		return "", 0, false
	***REMOVED***
	fam, bStr := req.FormValue("fam"), req.FormValue("b")
	if fam == "" || bStr == "" ***REMOVED***
		return "", 0, false
	***REMOVED***
	b, err := strconv.Atoi(bStr)
	if err != nil || b < -1 ***REMOVED***
		return "", 0, false
	***REMOVED***

	return fam, b, true
***REMOVED***

func lookupBucket(fam string, b int) *traceBucket ***REMOVED***
	f := getFamily(fam, false)
	if f == nil || b < 0 || b >= len(f.Buckets) ***REMOVED***
		return nil
	***REMOVED***
	return f.Buckets[b]
***REMOVED***

type contextKeyT string

var contextKey = contextKeyT("golang.org/x/net/trace.Trace")

// Trace represents an active request.
type Trace interface ***REMOVED***
	// LazyLog adds x to the event log. It will be evaluated each time the
	// /debug/requests page is rendered. Any memory referenced by x will be
	// pinned until the trace is finished and later discarded.
	LazyLog(x fmt.Stringer, sensitive bool)

	// LazyPrintf evaluates its arguments with fmt.Sprintf each time the
	// /debug/requests page is rendered. Any memory referenced by a will be
	// pinned until the trace is finished and later discarded.
	LazyPrintf(format string, a ...interface***REMOVED******REMOVED***)

	// SetError declares that this trace resulted in an error.
	SetError()

	// SetRecycler sets a recycler for the trace.
	// f will be called for each event passed to LazyLog at a time when
	// it is no longer required, whether while the trace is still active
	// and the event is discarded, or when a completed trace is discarded.
	SetRecycler(f func(interface***REMOVED******REMOVED***))

	// SetTraceInfo sets the trace info for the trace.
	// This is currently unused.
	SetTraceInfo(traceID, spanID uint64)

	// SetMaxEvents sets the maximum number of events that will be stored
	// in the trace. This has no effect if any events have already been
	// added to the trace.
	SetMaxEvents(m int)

	// Finish declares that this trace is complete.
	// The trace should not be used after calling this method.
	Finish()
***REMOVED***

type lazySprintf struct ***REMOVED***
	format string
	a      []interface***REMOVED******REMOVED***
***REMOVED***

func (l *lazySprintf) String() string ***REMOVED***
	return fmt.Sprintf(l.format, l.a...)
***REMOVED***

// New returns a new Trace with the specified family and title.
func New(family, title string) Trace ***REMOVED***
	tr := newTrace()
	tr.ref()
	tr.Family, tr.Title = family, title
	tr.Start = time.Now()
	tr.maxEvents = maxEventsPerTrace
	tr.events = tr.eventsBuf[:0]

	activeMu.RLock()
	s := activeTraces[tr.Family]
	activeMu.RUnlock()
	if s == nil ***REMOVED***
		activeMu.Lock()
		s = activeTraces[tr.Family] // check again
		if s == nil ***REMOVED***
			s = new(traceSet)
			activeTraces[tr.Family] = s
		***REMOVED***
		activeMu.Unlock()
	***REMOVED***
	s.Add(tr)

	// Trigger allocation of the completed trace structure for this family.
	// This will cause the family to be present in the request page during
	// the first trace of this family. We don't care about the return value,
	// nor is there any need for this to run inline, so we execute it in its
	// own goroutine, but only if the family isn't allocated yet.
	completedMu.RLock()
	if _, ok := completedTraces[tr.Family]; !ok ***REMOVED***
		go allocFamily(tr.Family)
	***REMOVED***
	completedMu.RUnlock()

	return tr
***REMOVED***

func (tr *trace) Finish() ***REMOVED***
	elapsed := time.Now().Sub(tr.Start)
	tr.mu.Lock()
	tr.Elapsed = elapsed
	tr.mu.Unlock()

	if DebugUseAfterFinish ***REMOVED***
		buf := make([]byte, 4<<10) // 4 KB should be enough
		n := runtime.Stack(buf, false)
		tr.finishStack = buf[:n]
	***REMOVED***

	activeMu.RLock()
	m := activeTraces[tr.Family]
	activeMu.RUnlock()
	m.Remove(tr)

	f := getFamily(tr.Family, true)
	tr.mu.RLock() // protects tr fields in Cond.match calls
	for _, b := range f.Buckets ***REMOVED***
		if b.Cond.match(tr) ***REMOVED***
			b.Add(tr)
		***REMOVED***
	***REMOVED***
	tr.mu.RUnlock()

	// Add a sample of elapsed time as microseconds to the family's timeseries
	h := new(histogram)
	h.addMeasurement(elapsed.Nanoseconds() / 1e3)
	f.LatencyMu.Lock()
	f.Latency.Add(h)
	f.LatencyMu.Unlock()

	tr.unref() // matches ref in New
***REMOVED***

const (
	bucketsPerFamily    = 9
	tracesPerBucket     = 10
	maxActiveTraces     = 20 // Maximum number of active traces to show.
	maxEventsPerTrace   = 10
	numHistogramBuckets = 38
)

var (
	// The active traces.
	activeMu     sync.RWMutex
	activeTraces = make(map[string]*traceSet) // family -> traces

	// Families of completed traces.
	completedMu     sync.RWMutex
	completedTraces = make(map[string]*family) // family -> traces
)

type traceSet struct ***REMOVED***
	mu sync.RWMutex
	m  map[*trace]bool

	// We could avoid the entire map scan in FirstN by having a slice of all the traces
	// ordered by start time, and an index into that from the trace struct, with a periodic
	// repack of the slice after enough traces finish; we could also use a skip list or similar.
	// However, that would shift some of the expense from /debug/requests time to RPC time,
	// which is probably the wrong trade-off.
***REMOVED***

func (ts *traceSet) Len() int ***REMOVED***
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return len(ts.m)
***REMOVED***

func (ts *traceSet) Add(tr *trace) ***REMOVED***
	ts.mu.Lock()
	if ts.m == nil ***REMOVED***
		ts.m = make(map[*trace]bool)
	***REMOVED***
	ts.m[tr] = true
	ts.mu.Unlock()
***REMOVED***

func (ts *traceSet) Remove(tr *trace) ***REMOVED***
	ts.mu.Lock()
	delete(ts.m, tr)
	ts.mu.Unlock()
***REMOVED***

// FirstN returns the first n traces ordered by time.
func (ts *traceSet) FirstN(n int) traceList ***REMOVED***
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if n > len(ts.m) ***REMOVED***
		n = len(ts.m)
	***REMOVED***
	trl := make(traceList, 0, n)

	// Fast path for when no selectivity is needed.
	if n == len(ts.m) ***REMOVED***
		for tr := range ts.m ***REMOVED***
			tr.ref()
			trl = append(trl, tr)
		***REMOVED***
		sort.Sort(trl)
		return trl
	***REMOVED***

	// Pick the oldest n traces.
	// This is inefficient. See the comment in the traceSet struct.
	for tr := range ts.m ***REMOVED***
		// Put the first n traces into trl in the order they occur.
		// When we have n, sort trl, and thereafter maintain its order.
		if len(trl) < n ***REMOVED***
			tr.ref()
			trl = append(trl, tr)
			if len(trl) == n ***REMOVED***
				// This is guaranteed to happen exactly once during this loop.
				sort.Sort(trl)
			***REMOVED***
			continue
		***REMOVED***
		if tr.Start.After(trl[n-1].Start) ***REMOVED***
			continue
		***REMOVED***

		// Find where to insert this one.
		tr.ref()
		i := sort.Search(n, func(i int) bool ***REMOVED*** return trl[i].Start.After(tr.Start) ***REMOVED***)
		trl[n-1].unref()
		copy(trl[i+1:], trl[i:])
		trl[i] = tr
	***REMOVED***

	return trl
***REMOVED***

func getActiveTraces(fam string) traceList ***REMOVED***
	activeMu.RLock()
	s := activeTraces[fam]
	activeMu.RUnlock()
	if s == nil ***REMOVED***
		return nil
	***REMOVED***
	return s.FirstN(maxActiveTraces)
***REMOVED***

func getFamily(fam string, allocNew bool) *family ***REMOVED***
	completedMu.RLock()
	f := completedTraces[fam]
	completedMu.RUnlock()
	if f == nil && allocNew ***REMOVED***
		f = allocFamily(fam)
	***REMOVED***
	return f
***REMOVED***

func allocFamily(fam string) *family ***REMOVED***
	completedMu.Lock()
	defer completedMu.Unlock()
	f := completedTraces[fam]
	if f == nil ***REMOVED***
		f = newFamily()
		completedTraces[fam] = f
	***REMOVED***
	return f
***REMOVED***

// family represents a set of trace buckets and associated latency information.
type family struct ***REMOVED***
	// traces may occur in multiple buckets.
	Buckets [bucketsPerFamily]*traceBucket

	// latency time series
	LatencyMu sync.RWMutex
	Latency   *timeseries.MinuteHourSeries
***REMOVED***

func newFamily() *family ***REMOVED***
	return &family***REMOVED***
		Buckets: [bucketsPerFamily]*traceBucket***REMOVED***
			***REMOVED***Cond: minCond(0)***REMOVED***,
			***REMOVED***Cond: minCond(50 * time.Millisecond)***REMOVED***,
			***REMOVED***Cond: minCond(100 * time.Millisecond)***REMOVED***,
			***REMOVED***Cond: minCond(200 * time.Millisecond)***REMOVED***,
			***REMOVED***Cond: minCond(500 * time.Millisecond)***REMOVED***,
			***REMOVED***Cond: minCond(1 * time.Second)***REMOVED***,
			***REMOVED***Cond: minCond(10 * time.Second)***REMOVED***,
			***REMOVED***Cond: minCond(100 * time.Second)***REMOVED***,
			***REMOVED***Cond: errorCond***REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		Latency: timeseries.NewMinuteHourSeries(func() timeseries.Observable ***REMOVED*** return new(histogram) ***REMOVED***),
	***REMOVED***
***REMOVED***

// traceBucket represents a size-capped bucket of historic traces,
// along with a condition for a trace to belong to the bucket.
type traceBucket struct ***REMOVED***
	Cond cond

	// Ring buffer implementation of a fixed-size FIFO queue.
	mu     sync.RWMutex
	buf    [tracesPerBucket]*trace
	start  int // < tracesPerBucket
	length int // <= tracesPerBucket
***REMOVED***

func (b *traceBucket) Add(tr *trace) ***REMOVED***
	b.mu.Lock()
	defer b.mu.Unlock()

	i := b.start + b.length
	if i >= tracesPerBucket ***REMOVED***
		i -= tracesPerBucket
	***REMOVED***
	if b.length == tracesPerBucket ***REMOVED***
		// "Remove" an element from the bucket.
		b.buf[i].unref()
		b.start++
		if b.start == tracesPerBucket ***REMOVED***
			b.start = 0
		***REMOVED***
	***REMOVED***
	b.buf[i] = tr
	if b.length < tracesPerBucket ***REMOVED***
		b.length++
	***REMOVED***
	tr.ref()
***REMOVED***

// Copy returns a copy of the traces in the bucket.
// If tracedOnly is true, only the traces with trace information will be returned.
// The logs will be ref'd before returning; the caller should call
// the Free method when it is done with them.
// TODO(dsymonds): keep track of traced requests in separate buckets.
func (b *traceBucket) Copy(tracedOnly bool) traceList ***REMOVED***
	b.mu.RLock()
	defer b.mu.RUnlock()

	trl := make(traceList, 0, b.length)
	for i, x := 0, b.start; i < b.length; i++ ***REMOVED***
		tr := b.buf[x]
		if !tracedOnly || tr.spanID != 0 ***REMOVED***
			tr.ref()
			trl = append(trl, tr)
		***REMOVED***
		x++
		if x == b.length ***REMOVED***
			x = 0
		***REMOVED***
	***REMOVED***
	return trl
***REMOVED***

func (b *traceBucket) Empty() bool ***REMOVED***
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.length == 0
***REMOVED***

// cond represents a condition on a trace.
type cond interface ***REMOVED***
	match(t *trace) bool
	String() string
***REMOVED***

type minCond time.Duration

func (m minCond) match(t *trace) bool ***REMOVED*** return t.Elapsed >= time.Duration(m) ***REMOVED***
func (m minCond) String() string      ***REMOVED*** return fmt.Sprintf("≥%gs", time.Duration(m).Seconds()) ***REMOVED***

type errorCond struct***REMOVED******REMOVED***

func (e errorCond) match(t *trace) bool ***REMOVED*** return t.IsError ***REMOVED***
func (e errorCond) String() string      ***REMOVED*** return "errors" ***REMOVED***

type traceList []*trace

// Free calls unref on each element of the list.
func (trl traceList) Free() ***REMOVED***
	for _, t := range trl ***REMOVED***
		t.unref()
	***REMOVED***
***REMOVED***

// traceList may be sorted in reverse chronological order.
func (trl traceList) Len() int           ***REMOVED*** return len(trl) ***REMOVED***
func (trl traceList) Less(i, j int) bool ***REMOVED*** return trl[i].Start.After(trl[j].Start) ***REMOVED***
func (trl traceList) Swap(i, j int)      ***REMOVED*** trl[i], trl[j] = trl[j], trl[i] ***REMOVED***

// An event is a timestamped log entry in a trace.
type event struct ***REMOVED***
	When       time.Time
	Elapsed    time.Duration // since previous event in trace
	NewDay     bool          // whether this event is on a different day to the previous event
	Recyclable bool          // whether this event was passed via LazyLog
	Sensitive  bool          // whether this event contains sensitive information
	What       interface***REMOVED******REMOVED***   // string or fmt.Stringer
***REMOVED***

// WhenString returns a string representation of the elapsed time of the event.
// It will include the date if midnight was crossed.
func (e event) WhenString() string ***REMOVED***
	if e.NewDay ***REMOVED***
		return e.When.Format("2006/01/02 15:04:05.000000")
	***REMOVED***
	return e.When.Format("15:04:05.000000")
***REMOVED***

// discarded represents a number of discarded events.
// It is stored as *discarded to make it easier to update in-place.
type discarded int

func (d *discarded) String() string ***REMOVED***
	return fmt.Sprintf("(%d events discarded)", int(*d))
***REMOVED***

// trace represents an active or complete request,
// either sent or received by this program.
type trace struct ***REMOVED***
	// Family is the top-level grouping of traces to which this belongs.
	Family string

	// Title is the title of this trace.
	Title string

	// Start time of the this trace.
	Start time.Time

	mu        sync.RWMutex
	events    []event // Append-only sequence of events (modulo discards).
	maxEvents int
	recycler  func(interface***REMOVED******REMOVED***)
	IsError   bool          // Whether this trace resulted in an error.
	Elapsed   time.Duration // Elapsed time for this trace, zero while active.
	traceID   uint64        // Trace information if non-zero.
	spanID    uint64

	refs int32     // how many buckets this is in
	disc discarded // scratch space to avoid allocation

	finishStack []byte // where finish was called, if DebugUseAfterFinish is set

	eventsBuf [4]event // preallocated buffer in case we only log a few events
***REMOVED***

func (tr *trace) reset() ***REMOVED***
	// Clear all but the mutex. Mutexes may not be copied, even when unlocked.
	tr.Family = ""
	tr.Title = ""
	tr.Start = time.Time***REMOVED******REMOVED***

	tr.mu.Lock()
	tr.Elapsed = 0
	tr.traceID = 0
	tr.spanID = 0
	tr.IsError = false
	tr.maxEvents = 0
	tr.events = nil
	tr.recycler = nil
	tr.mu.Unlock()

	tr.refs = 0
	tr.disc = 0
	tr.finishStack = nil
	for i := range tr.eventsBuf ***REMOVED***
		tr.eventsBuf[i] = event***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// delta returns the elapsed time since the last event or the trace start,
// and whether it spans midnight.
// L >= tr.mu
func (tr *trace) delta(t time.Time) (time.Duration, bool) ***REMOVED***
	if len(tr.events) == 0 ***REMOVED***
		return t.Sub(tr.Start), false
	***REMOVED***
	prev := tr.events[len(tr.events)-1].When
	return t.Sub(prev), prev.Day() != t.Day()
***REMOVED***

func (tr *trace) addEvent(x interface***REMOVED******REMOVED***, recyclable, sensitive bool) ***REMOVED***
	if DebugUseAfterFinish && tr.finishStack != nil ***REMOVED***
		buf := make([]byte, 4<<10) // 4 KB should be enough
		n := runtime.Stack(buf, false)
		log.Printf("net/trace: trace used after finish:\nFinished at:\n%s\nUsed at:\n%s", tr.finishStack, buf[:n])
	***REMOVED***

	/*
		NOTE TO DEBUGGERS

		If you are here because your program panicked in this code,
		it is almost definitely the fault of code using this package,
		and very unlikely to be the fault of this code.

		The most likely scenario is that some code elsewhere is using
		a trace.Trace after its Finish method is called.
		You can temporarily set the DebugUseAfterFinish var
		to help discover where that is; do not leave that var set,
		since it makes this package much less efficient.
	*/

	e := event***REMOVED***When: time.Now(), What: x, Recyclable: recyclable, Sensitive: sensitive***REMOVED***
	tr.mu.Lock()
	e.Elapsed, e.NewDay = tr.delta(e.When)
	if len(tr.events) < tr.maxEvents ***REMOVED***
		tr.events = append(tr.events, e)
	***REMOVED*** else ***REMOVED***
		// Discard the middle events.
		di := int((tr.maxEvents - 1) / 2)
		if d, ok := tr.events[di].What.(*discarded); ok ***REMOVED***
			(*d)++
		***REMOVED*** else ***REMOVED***
			// disc starts at two to count for the event it is replacing,
			// plus the next one that we are about to drop.
			tr.disc = 2
			if tr.recycler != nil && tr.events[di].Recyclable ***REMOVED***
				go tr.recycler(tr.events[di].What)
			***REMOVED***
			tr.events[di].What = &tr.disc
		***REMOVED***
		// The timestamp of the discarded meta-event should be
		// the time of the last event it is representing.
		tr.events[di].When = tr.events[di+1].When

		if tr.recycler != nil && tr.events[di+1].Recyclable ***REMOVED***
			go tr.recycler(tr.events[di+1].What)
		***REMOVED***
		copy(tr.events[di+1:], tr.events[di+2:])
		tr.events[tr.maxEvents-1] = e
	***REMOVED***
	tr.mu.Unlock()
***REMOVED***

func (tr *trace) LazyLog(x fmt.Stringer, sensitive bool) ***REMOVED***
	tr.addEvent(x, true, sensitive)
***REMOVED***

func (tr *trace) LazyPrintf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	tr.addEvent(&lazySprintf***REMOVED***format, a***REMOVED***, false, false)
***REMOVED***

func (tr *trace) SetError() ***REMOVED***
	tr.mu.Lock()
	tr.IsError = true
	tr.mu.Unlock()
***REMOVED***

func (tr *trace) SetRecycler(f func(interface***REMOVED******REMOVED***)) ***REMOVED***
	tr.mu.Lock()
	tr.recycler = f
	tr.mu.Unlock()
***REMOVED***

func (tr *trace) SetTraceInfo(traceID, spanID uint64) ***REMOVED***
	tr.mu.Lock()
	tr.traceID, tr.spanID = traceID, spanID
	tr.mu.Unlock()
***REMOVED***

func (tr *trace) SetMaxEvents(m int) ***REMOVED***
	tr.mu.Lock()
	// Always keep at least three events: first, discarded count, last.
	if len(tr.events) == 0 && m > 3 ***REMOVED***
		tr.maxEvents = m
	***REMOVED***
	tr.mu.Unlock()
***REMOVED***

func (tr *trace) ref() ***REMOVED***
	atomic.AddInt32(&tr.refs, 1)
***REMOVED***

func (tr *trace) unref() ***REMOVED***
	if atomic.AddInt32(&tr.refs, -1) == 0 ***REMOVED***
		tr.mu.RLock()
		if tr.recycler != nil ***REMOVED***
			// freeTrace clears tr, so we hold tr.recycler and tr.events here.
			go func(f func(interface***REMOVED******REMOVED***), es []event) ***REMOVED***
				for _, e := range es ***REMOVED***
					if e.Recyclable ***REMOVED***
						f(e.What)
					***REMOVED***
				***REMOVED***
			***REMOVED***(tr.recycler, tr.events)
		***REMOVED***
		tr.mu.RUnlock()

		freeTrace(tr)
	***REMOVED***
***REMOVED***

func (tr *trace) When() string ***REMOVED***
	return tr.Start.Format("2006/01/02 15:04:05.000000")
***REMOVED***

func (tr *trace) ElapsedTime() string ***REMOVED***
	tr.mu.RLock()
	t := tr.Elapsed
	tr.mu.RUnlock()

	if t == 0 ***REMOVED***
		// Active trace.
		t = time.Since(tr.Start)
	***REMOVED***
	return fmt.Sprintf("%.6f", t.Seconds())
***REMOVED***

func (tr *trace) Events() []event ***REMOVED***
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.events
***REMOVED***

var traceFreeList = make(chan *trace, 1000) // TODO(dsymonds): Use sync.Pool?

// newTrace returns a trace ready to use.
func newTrace() *trace ***REMOVED***
	select ***REMOVED***
	case tr := <-traceFreeList:
		return tr
	default:
		return new(trace)
	***REMOVED***
***REMOVED***

// freeTrace adds tr to traceFreeList if there's room.
// This is non-blocking.
func freeTrace(tr *trace) ***REMOVED***
	if DebugUseAfterFinish ***REMOVED***
		return // never reuse
	***REMOVED***
	tr.reset()
	select ***REMOVED***
	case traceFreeList <- tr:
	default:
	***REMOVED***
***REMOVED***

func elapsed(d time.Duration) string ***REMOVED***
	b := []byte(fmt.Sprintf("%.6f", d.Seconds()))

	// For subsecond durations, blank all zeros before decimal point,
	// and all zeros between the decimal point and the first non-zero digit.
	if d < time.Second ***REMOVED***
		dot := bytes.IndexByte(b, '.')
		for i := 0; i < dot; i++ ***REMOVED***
			b[i] = ' '
		***REMOVED***
		for i := dot + 1; i < len(b); i++ ***REMOVED***
			if b[i] == '0' ***REMOVED***
				b[i] = ' '
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return string(b)
***REMOVED***

var pageTmplCache *template.Template
var pageTmplOnce sync.Once

func pageTmpl() *template.Template ***REMOVED***
	pageTmplOnce.Do(func() ***REMOVED***
		pageTmplCache = template.Must(template.New("Page").Funcs(template.FuncMap***REMOVED***
			"elapsed": elapsed,
			"add":     func(a, b int) int ***REMOVED*** return a + b ***REMOVED***,
		***REMOVED***).Parse(pageHTML))
	***REMOVED***)
	return pageTmplCache
***REMOVED***

const pageHTML = `
***REMOVED******REMOVED***template "Prolog" .***REMOVED******REMOVED***
***REMOVED******REMOVED***template "StatusTable" .***REMOVED******REMOVED***
***REMOVED******REMOVED***template "Epilog" .***REMOVED******REMOVED***

***REMOVED******REMOVED***define "Prolog"***REMOVED******REMOVED***
<html>
	<head>
	<title>/debug/requests</title>
	<style type="text/css">
		body ***REMOVED***
			font-family: sans-serif;
		***REMOVED***
		table#tr-status td.family ***REMOVED***
			padding-right: 2em;
		***REMOVED***
		table#tr-status td.active ***REMOVED***
			padding-right: 1em;
		***REMOVED***
		table#tr-status td.latency-first ***REMOVED***
			padding-left: 1em;
		***REMOVED***
		table#tr-status td.empty ***REMOVED***
			color: #aaa;
		***REMOVED***
		table#reqs ***REMOVED***
			margin-top: 1em;
		***REMOVED***
		table#reqs tr.first ***REMOVED***
			***REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***font-weight: bold;***REMOVED******REMOVED***end***REMOVED******REMOVED***
		***REMOVED***
		table#reqs td ***REMOVED***
			font-family: monospace;
		***REMOVED***
		table#reqs td.when ***REMOVED***
			text-align: right;
			white-space: nowrap;
		***REMOVED***
		table#reqs td.elapsed ***REMOVED***
			padding: 0 0.5em;
			text-align: right;
			white-space: pre;
			width: 10em;
		***REMOVED***
		address ***REMOVED***
			font-size: smaller;
			margin-top: 5em;
		***REMOVED***
	</style>
	</head>
	<body>

<h1>/debug/requests</h1>
***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***/* end of Prolog */***REMOVED******REMOVED***

***REMOVED******REMOVED***define "StatusTable"***REMOVED******REMOVED***
<table id="tr-status">
	***REMOVED******REMOVED***range $fam := .Families***REMOVED******REMOVED***
	<tr>
		<td class="family">***REMOVED******REMOVED***$fam***REMOVED******REMOVED***</td>

		***REMOVED******REMOVED***$n := index $.ActiveTraceCount $fam***REMOVED******REMOVED***
		<td class="active ***REMOVED******REMOVED***if not $n***REMOVED******REMOVED***empty***REMOVED******REMOVED***end***REMOVED******REMOVED***">
			***REMOVED******REMOVED***if $n***REMOVED******REMOVED***<a href="?fam=***REMOVED******REMOVED***$fam***REMOVED******REMOVED***&b=-1***REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***&exp=1***REMOVED******REMOVED***end***REMOVED******REMOVED***">***REMOVED******REMOVED***end***REMOVED******REMOVED***
			[***REMOVED******REMOVED***$n***REMOVED******REMOVED*** active]
			***REMOVED******REMOVED***if $n***REMOVED******REMOVED***</a>***REMOVED******REMOVED***end***REMOVED******REMOVED***
		</td>

		***REMOVED******REMOVED***$f := index $.CompletedTraces $fam***REMOVED******REMOVED***
		***REMOVED******REMOVED***range $i, $b := $f.Buckets***REMOVED******REMOVED***
		***REMOVED******REMOVED***$empty := $b.Empty***REMOVED******REMOVED***
		<td ***REMOVED******REMOVED***if $empty***REMOVED******REMOVED***class="empty"***REMOVED******REMOVED***end***REMOVED******REMOVED***>
		***REMOVED******REMOVED***if not $empty***REMOVED******REMOVED***<a href="?fam=***REMOVED******REMOVED***$fam***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$i***REMOVED******REMOVED******REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***&exp=1***REMOVED******REMOVED***end***REMOVED******REMOVED***">***REMOVED******REMOVED***end***REMOVED******REMOVED***
		[***REMOVED******REMOVED***.Cond***REMOVED******REMOVED***]
		***REMOVED******REMOVED***if not $empty***REMOVED******REMOVED***</a>***REMOVED******REMOVED***end***REMOVED******REMOVED***
		</td>
		***REMOVED******REMOVED***end***REMOVED******REMOVED***

		***REMOVED******REMOVED***$nb := len $f.Buckets***REMOVED******REMOVED***
		<td class="latency-first">
		<a href="?fam=***REMOVED******REMOVED***$fam***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$nb***REMOVED******REMOVED***">[minute]</a>
		</td>
		<td>
		<a href="?fam=***REMOVED******REMOVED***$fam***REMOVED******REMOVED***&b=***REMOVED******REMOVED***add $nb 1***REMOVED******REMOVED***">[hour]</a>
		</td>
		<td>
		<a href="?fam=***REMOVED******REMOVED***$fam***REMOVED******REMOVED***&b=***REMOVED******REMOVED***add $nb 2***REMOVED******REMOVED***">[total]</a>
		</td>

	</tr>
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
</table>
***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***/* end of StatusTable */***REMOVED******REMOVED***

***REMOVED******REMOVED***define "Epilog"***REMOVED******REMOVED***
***REMOVED******REMOVED***if $.Traces***REMOVED******REMOVED***
<hr />
<h3>Family: ***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***</h3>

***REMOVED******REMOVED***if or $.Expanded $.Traced***REMOVED******REMOVED***
  <a href="?fam=***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$.Bucket***REMOVED******REMOVED***">[Normal/Summary]</a>
***REMOVED******REMOVED***else***REMOVED******REMOVED***
  [Normal/Summary]
***REMOVED******REMOVED***end***REMOVED******REMOVED***

***REMOVED******REMOVED***if or (not $.Expanded) $.Traced***REMOVED******REMOVED***
  <a href="?fam=***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$.Bucket***REMOVED******REMOVED***&exp=1">[Normal/Expanded]</a>
***REMOVED******REMOVED***else***REMOVED******REMOVED***
  [Normal/Expanded]
***REMOVED******REMOVED***end***REMOVED******REMOVED***

***REMOVED******REMOVED***if not $.Active***REMOVED******REMOVED***
	***REMOVED******REMOVED***if or $.Expanded (not $.Traced)***REMOVED******REMOVED***
	<a href="?fam=***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$.Bucket***REMOVED******REMOVED***&rtraced=1">[Traced/Summary]</a>
	***REMOVED******REMOVED***else***REMOVED******REMOVED***
	[Traced/Summary]
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED******REMOVED***if or (not $.Expanded) (not $.Traced)***REMOVED******REMOVED***
	<a href="?fam=***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$.Bucket***REMOVED******REMOVED***&exp=1&rtraced=1">[Traced/Expanded]</a>
        ***REMOVED******REMOVED***else***REMOVED******REMOVED***
	[Traced/Expanded]
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED***

***REMOVED******REMOVED***if $.Total***REMOVED******REMOVED***
<p><em>Showing <b>***REMOVED******REMOVED***len $.Traces***REMOVED******REMOVED***</b> of <b>***REMOVED******REMOVED***$.Total***REMOVED******REMOVED***</b> traces.</em></p>
***REMOVED******REMOVED***end***REMOVED******REMOVED***

<table id="reqs">
	<caption>
		***REMOVED******REMOVED***if $.Active***REMOVED******REMOVED***Active***REMOVED******REMOVED***else***REMOVED******REMOVED***Completed***REMOVED******REMOVED***end***REMOVED******REMOVED*** Requests
	</caption>
	<tr><th>When</th><th>Elapsed&nbsp;(s)</th></tr>
	***REMOVED******REMOVED***range $tr := $.Traces***REMOVED******REMOVED***
	<tr class="first">
		<td class="when">***REMOVED******REMOVED***$tr.When***REMOVED******REMOVED***</td>
		<td class="elapsed">***REMOVED******REMOVED***$tr.ElapsedTime***REMOVED******REMOVED***</td>
		<td>***REMOVED******REMOVED***$tr.Title***REMOVED******REMOVED***</td>
		***REMOVED******REMOVED***/* TODO: include traceID/spanID */***REMOVED******REMOVED***
	</tr>
	***REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***
	***REMOVED******REMOVED***range $tr.Events***REMOVED******REMOVED***
	<tr>
		<td class="when">***REMOVED******REMOVED***.WhenString***REMOVED******REMOVED***</td>
		<td class="elapsed">***REMOVED******REMOVED***elapsed .Elapsed***REMOVED******REMOVED***</td>
		<td>***REMOVED******REMOVED***if or $.ShowSensitive (not .Sensitive)***REMOVED******REMOVED***... ***REMOVED******REMOVED***.What***REMOVED******REMOVED******REMOVED******REMOVED***else***REMOVED******REMOVED***<em>[redacted]</em>***REMOVED******REMOVED***end***REMOVED******REMOVED***</td>
	</tr>
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
</table>
***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***/* if $.Traces */***REMOVED******REMOVED***

***REMOVED******REMOVED***if $.Histogram***REMOVED******REMOVED***
<h4>Latency (&micro;s) of ***REMOVED******REMOVED***$.Family***REMOVED******REMOVED*** over ***REMOVED******REMOVED***$.HistogramWindow***REMOVED******REMOVED***</h4>
***REMOVED******REMOVED***$.Histogram***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***/* if $.Histogram */***REMOVED******REMOVED***

	</body>
</html>
***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***/* end of Epilog */***REMOVED******REMOVED***
`
