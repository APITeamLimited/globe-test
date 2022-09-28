// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"
)

const maxEventsPerLog = 100

type bucket struct ***REMOVED***
	MaxErrAge time.Duration
	String    string
***REMOVED***

var buckets = []bucket***REMOVED***
	***REMOVED***0, "total"***REMOVED***,
	***REMOVED***10 * time.Second, "errs<10s"***REMOVED***,
	***REMOVED***1 * time.Minute, "errs<1m"***REMOVED***,
	***REMOVED***10 * time.Minute, "errs<10m"***REMOVED***,
	***REMOVED***1 * time.Hour, "errs<1h"***REMOVED***,
	***REMOVED***10 * time.Hour, "errs<10h"***REMOVED***,
	***REMOVED***24000 * time.Hour, "errors"***REMOVED***,
***REMOVED***

// RenderEvents renders the HTML page typically served at /debug/events.
// It does not do any auth checking. The request may be nil.
//
// Most users will use the Events handler.
func RenderEvents(w http.ResponseWriter, req *http.Request, sensitive bool) ***REMOVED***
	now := time.Now()
	data := &struct ***REMOVED***
		Families []string // family names
		Buckets  []bucket
		Counts   [][]int // eventLog count per family/bucket

		// Set when a bucket has been selected.
		Family    string
		Bucket    int
		EventLogs eventLogs
		Expanded  bool
	***REMOVED******REMOVED***
		Buckets: buckets,
	***REMOVED***

	data.Families = make([]string, 0, len(families))
	famMu.RLock()
	for name := range families ***REMOVED***
		data.Families = append(data.Families, name)
	***REMOVED***
	famMu.RUnlock()
	sort.Strings(data.Families)

	// Count the number of eventLogs in each family for each error age.
	data.Counts = make([][]int, len(data.Families))
	for i, name := range data.Families ***REMOVED***
		// TODO(sameer): move this loop under the family lock.
		f := getEventFamily(name)
		data.Counts[i] = make([]int, len(data.Buckets))
		for j, b := range data.Buckets ***REMOVED***
			data.Counts[i][j] = f.Count(now, b.MaxErrAge)
		***REMOVED***
	***REMOVED***

	if req != nil ***REMOVED***
		var ok bool
		data.Family, data.Bucket, ok = parseEventsArgs(req)
		if !ok ***REMOVED***
			// No-op
		***REMOVED*** else ***REMOVED***
			data.EventLogs = getEventFamily(data.Family).Copy(now, buckets[data.Bucket].MaxErrAge)
		***REMOVED***
		if data.EventLogs != nil ***REMOVED***
			defer data.EventLogs.Free()
			sort.Sort(data.EventLogs)
		***REMOVED***
		if exp, err := strconv.ParseBool(req.FormValue("exp")); err == nil ***REMOVED***
			data.Expanded = exp
		***REMOVED***
	***REMOVED***

	famMu.RLock()
	defer famMu.RUnlock()
	if err := eventsTmpl().Execute(w, data); err != nil ***REMOVED***
		log.Printf("net/trace: Failed executing template: %v", err)
	***REMOVED***
***REMOVED***

func parseEventsArgs(req *http.Request) (fam string, b int, ok bool) ***REMOVED***
	fam, bStr := req.FormValue("fam"), req.FormValue("b")
	if fam == "" || bStr == "" ***REMOVED***
		return "", 0, false
	***REMOVED***
	b, err := strconv.Atoi(bStr)
	if err != nil || b < 0 || b >= len(buckets) ***REMOVED***
		return "", 0, false
	***REMOVED***
	return fam, b, true
***REMOVED***

// An EventLog provides a log of events associated with a specific object.
type EventLog interface ***REMOVED***
	// Printf formats its arguments with fmt.Sprintf and adds the
	// result to the event log.
	Printf(format string, a ...interface***REMOVED******REMOVED***)

	// Errorf is like Printf, but it marks this event as an error.
	Errorf(format string, a ...interface***REMOVED******REMOVED***)

	// Finish declares that this event log is complete.
	// The event log should not be used after calling this method.
	Finish()
***REMOVED***

// NewEventLog returns a new EventLog with the specified family name
// and title.
func NewEventLog(family, title string) EventLog ***REMOVED***
	el := newEventLog()
	el.ref()
	el.Family, el.Title = family, title
	el.Start = time.Now()
	el.events = make([]logEntry, 0, maxEventsPerLog)
	el.stack = make([]uintptr, 32)
	n := runtime.Callers(2, el.stack)
	el.stack = el.stack[:n]

	getEventFamily(family).add(el)
	return el
***REMOVED***

func (el *eventLog) Finish() ***REMOVED***
	getEventFamily(el.Family).remove(el)
	el.unref() // matches ref in New
***REMOVED***

var (
	famMu    sync.RWMutex
	families = make(map[string]*eventFamily) // family name => family
)

func getEventFamily(fam string) *eventFamily ***REMOVED***
	famMu.Lock()
	defer famMu.Unlock()
	f := families[fam]
	if f == nil ***REMOVED***
		f = &eventFamily***REMOVED******REMOVED***
		families[fam] = f
	***REMOVED***
	return f
***REMOVED***

type eventFamily struct ***REMOVED***
	mu        sync.RWMutex
	eventLogs eventLogs
***REMOVED***

func (f *eventFamily) add(el *eventLog) ***REMOVED***
	f.mu.Lock()
	f.eventLogs = append(f.eventLogs, el)
	f.mu.Unlock()
***REMOVED***

func (f *eventFamily) remove(el *eventLog) ***REMOVED***
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, el0 := range f.eventLogs ***REMOVED***
		if el == el0 ***REMOVED***
			copy(f.eventLogs[i:], f.eventLogs[i+1:])
			f.eventLogs = f.eventLogs[:len(f.eventLogs)-1]
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *eventFamily) Count(now time.Time, maxErrAge time.Duration) (n int) ***REMOVED***
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, el := range f.eventLogs ***REMOVED***
		if el.hasRecentError(now, maxErrAge) ***REMOVED***
			n++
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (f *eventFamily) Copy(now time.Time, maxErrAge time.Duration) (els eventLogs) ***REMOVED***
	f.mu.RLock()
	defer f.mu.RUnlock()
	els = make(eventLogs, 0, len(f.eventLogs))
	for _, el := range f.eventLogs ***REMOVED***
		if el.hasRecentError(now, maxErrAge) ***REMOVED***
			el.ref()
			els = append(els, el)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

type eventLogs []*eventLog

// Free calls unref on each element of the list.
func (els eventLogs) Free() ***REMOVED***
	for _, el := range els ***REMOVED***
		el.unref()
	***REMOVED***
***REMOVED***

// eventLogs may be sorted in reverse chronological order.
func (els eventLogs) Len() int           ***REMOVED*** return len(els) ***REMOVED***
func (els eventLogs) Less(i, j int) bool ***REMOVED*** return els[i].Start.After(els[j].Start) ***REMOVED***
func (els eventLogs) Swap(i, j int)      ***REMOVED*** els[i], els[j] = els[j], els[i] ***REMOVED***

// A logEntry is a timestamped log entry in an event log.
type logEntry struct ***REMOVED***
	When    time.Time
	Elapsed time.Duration // since previous event in log
	NewDay  bool          // whether this event is on a different day to the previous event
	What    string
	IsErr   bool
***REMOVED***

// WhenString returns a string representation of the elapsed time of the event.
// It will include the date if midnight was crossed.
func (e logEntry) WhenString() string ***REMOVED***
	if e.NewDay ***REMOVED***
		return e.When.Format("2006/01/02 15:04:05.000000")
	***REMOVED***
	return e.When.Format("15:04:05.000000")
***REMOVED***

// An eventLog represents an active event log.
type eventLog struct ***REMOVED***
	// Family is the top-level grouping of event logs to which this belongs.
	Family string

	// Title is the title of this event log.
	Title string

	// Timing information.
	Start time.Time

	// Call stack where this event log was created.
	stack []uintptr

	// Append-only sequence of events.
	//
	// TODO(sameer): change this to a ring buffer to avoid the array copy
	// when we hit maxEventsPerLog.
	mu            sync.RWMutex
	events        []logEntry
	LastErrorTime time.Time
	discarded     int

	refs int32 // how many buckets this is in
***REMOVED***

func (el *eventLog) reset() ***REMOVED***
	// Clear all but the mutex. Mutexes may not be copied, even when unlocked.
	el.Family = ""
	el.Title = ""
	el.Start = time.Time***REMOVED******REMOVED***
	el.stack = nil
	el.events = nil
	el.LastErrorTime = time.Time***REMOVED******REMOVED***
	el.discarded = 0
	el.refs = 0
***REMOVED***

func (el *eventLog) hasRecentError(now time.Time, maxErrAge time.Duration) bool ***REMOVED***
	if maxErrAge == 0 ***REMOVED***
		return true
	***REMOVED***
	el.mu.RLock()
	defer el.mu.RUnlock()
	return now.Sub(el.LastErrorTime) < maxErrAge
***REMOVED***

// delta returns the elapsed time since the last event or the log start,
// and whether it spans midnight.
// L >= el.mu
func (el *eventLog) delta(t time.Time) (time.Duration, bool) ***REMOVED***
	if len(el.events) == 0 ***REMOVED***
		return t.Sub(el.Start), false
	***REMOVED***
	prev := el.events[len(el.events)-1].When
	return t.Sub(prev), prev.Day() != t.Day()

***REMOVED***

func (el *eventLog) Printf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	el.printf(false, format, a...)
***REMOVED***

func (el *eventLog) Errorf(format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	el.printf(true, format, a...)
***REMOVED***

func (el *eventLog) printf(isErr bool, format string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
	e := logEntry***REMOVED***When: time.Now(), IsErr: isErr, What: fmt.Sprintf(format, a...)***REMOVED***
	el.mu.Lock()
	e.Elapsed, e.NewDay = el.delta(e.When)
	if len(el.events) < maxEventsPerLog ***REMOVED***
		el.events = append(el.events, e)
	***REMOVED*** else ***REMOVED***
		// Discard the oldest event.
		if el.discarded == 0 ***REMOVED***
			// el.discarded starts at two to count for the event it
			// is replacing, plus the next one that we are about to
			// drop.
			el.discarded = 2
		***REMOVED*** else ***REMOVED***
			el.discarded++
		***REMOVED***
		// TODO(sameer): if this causes allocations on a critical path,
		// change eventLog.What to be a fmt.Stringer, as in trace.go.
		el.events[0].What = fmt.Sprintf("(%d events discarded)", el.discarded)
		// The timestamp of the discarded meta-event should be
		// the time of the last event it is representing.
		el.events[0].When = el.events[1].When
		copy(el.events[1:], el.events[2:])
		el.events[maxEventsPerLog-1] = e
	***REMOVED***
	if e.IsErr ***REMOVED***
		el.LastErrorTime = e.When
	***REMOVED***
	el.mu.Unlock()
***REMOVED***

func (el *eventLog) ref() ***REMOVED***
	atomic.AddInt32(&el.refs, 1)
***REMOVED***

func (el *eventLog) unref() ***REMOVED***
	if atomic.AddInt32(&el.refs, -1) == 0 ***REMOVED***
		freeEventLog(el)
	***REMOVED***
***REMOVED***

func (el *eventLog) When() string ***REMOVED***
	return el.Start.Format("2006/01/02 15:04:05.000000")
***REMOVED***

func (el *eventLog) ElapsedTime() string ***REMOVED***
	elapsed := time.Since(el.Start)
	return fmt.Sprintf("%.6f", elapsed.Seconds())
***REMOVED***

func (el *eventLog) Stack() string ***REMOVED***
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 1, 8, 1, '\t', 0)
	printStackRecord(tw, el.stack)
	tw.Flush()
	return buf.String()
***REMOVED***

// printStackRecord prints the function + source line information
// for a single stack trace.
// Adapted from runtime/pprof/pprof.go.
func printStackRecord(w io.Writer, stk []uintptr) ***REMOVED***
	for _, pc := range stk ***REMOVED***
		f := runtime.FuncForPC(pc)
		if f == nil ***REMOVED***
			continue
		***REMOVED***
		file, line := f.FileLine(pc)
		name := f.Name()
		// Hide runtime.goexit and any runtime functions at the beginning.
		if strings.HasPrefix(name, "runtime.") ***REMOVED***
			continue
		***REMOVED***
		fmt.Fprintf(w, "#   %s\t%s:%d\n", name, file, line)
	***REMOVED***
***REMOVED***

func (el *eventLog) Events() []logEntry ***REMOVED***
	el.mu.RLock()
	defer el.mu.RUnlock()
	return el.events
***REMOVED***

// freeEventLogs is a freelist of *eventLog
var freeEventLogs = make(chan *eventLog, 1000)

// newEventLog returns a event log ready to use.
func newEventLog() *eventLog ***REMOVED***
	select ***REMOVED***
	case el := <-freeEventLogs:
		return el
	default:
		return new(eventLog)
	***REMOVED***
***REMOVED***

// freeEventLog adds el to freeEventLogs if there's room.
// This is non-blocking.
func freeEventLog(el *eventLog) ***REMOVED***
	el.reset()
	select ***REMOVED***
	case freeEventLogs <- el:
	default:
	***REMOVED***
***REMOVED***

var eventsTmplCache *template.Template
var eventsTmplOnce sync.Once

func eventsTmpl() *template.Template ***REMOVED***
	eventsTmplOnce.Do(func() ***REMOVED***
		eventsTmplCache = template.Must(template.New("events").Funcs(template.FuncMap***REMOVED***
			"elapsed":   elapsed,
			"trimSpace": strings.TrimSpace,
		***REMOVED***).Parse(eventsHTML))
	***REMOVED***)
	return eventsTmplCache
***REMOVED***

const eventsHTML = `
<html>
	<head>
		<title>events</title>
	</head>
	<style type="text/css">
		body ***REMOVED***
			font-family: sans-serif;
		***REMOVED***
		table#req-status td.family ***REMOVED***
			padding-right: 2em;
		***REMOVED***
		table#req-status td.active ***REMOVED***
			padding-right: 1em;
		***REMOVED***
		table#req-status td.empty ***REMOVED***
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
	<body>

<h1>/debug/events</h1>

<table id="req-status">
	***REMOVED******REMOVED***range $i, $fam := .Families***REMOVED******REMOVED***
	<tr>
		<td class="family">***REMOVED******REMOVED***$fam***REMOVED******REMOVED***</td>

	        ***REMOVED******REMOVED***range $j, $bucket := $.Buckets***REMOVED******REMOVED***
	        ***REMOVED******REMOVED***$n := index $.Counts $i $j***REMOVED******REMOVED***
		<td class="***REMOVED******REMOVED***if not $bucket.MaxErrAge***REMOVED******REMOVED***active***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***if not $n***REMOVED******REMOVED***empty***REMOVED******REMOVED***end***REMOVED******REMOVED***">
	                ***REMOVED******REMOVED***if $n***REMOVED******REMOVED***<a href="?fam=***REMOVED******REMOVED***$fam***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$j***REMOVED******REMOVED******REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***&exp=1***REMOVED******REMOVED***end***REMOVED******REMOVED***">***REMOVED******REMOVED***end***REMOVED******REMOVED***
		        [***REMOVED******REMOVED***$n***REMOVED******REMOVED*** ***REMOVED******REMOVED***$bucket.String***REMOVED******REMOVED***]
			***REMOVED******REMOVED***if $n***REMOVED******REMOVED***</a>***REMOVED******REMOVED***end***REMOVED******REMOVED***
		</td>
                ***REMOVED******REMOVED***end***REMOVED******REMOVED***

	</tr>***REMOVED******REMOVED***end***REMOVED******REMOVED***
</table>

***REMOVED******REMOVED***if $.EventLogs***REMOVED******REMOVED***
<hr />
<h3>Family: ***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***</h3>

***REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***<a href="?fam=***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$.Bucket***REMOVED******REMOVED***">***REMOVED******REMOVED***end***REMOVED******REMOVED***
[Summary]***REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***</a>***REMOVED******REMOVED***end***REMOVED******REMOVED***

***REMOVED******REMOVED***if not $.Expanded***REMOVED******REMOVED***<a href="?fam=***REMOVED******REMOVED***$.Family***REMOVED******REMOVED***&b=***REMOVED******REMOVED***$.Bucket***REMOVED******REMOVED***&exp=1">***REMOVED******REMOVED***end***REMOVED******REMOVED***
[Expanded]***REMOVED******REMOVED***if not $.Expanded***REMOVED******REMOVED***</a>***REMOVED******REMOVED***end***REMOVED******REMOVED***

<table id="reqs">
	<tr><th>When</th><th>Elapsed</th></tr>
	***REMOVED******REMOVED***range $el := $.EventLogs***REMOVED******REMOVED***
	<tr class="first">
		<td class="when">***REMOVED******REMOVED***$el.When***REMOVED******REMOVED***</td>
		<td class="elapsed">***REMOVED******REMOVED***$el.ElapsedTime***REMOVED******REMOVED***</td>
		<td>***REMOVED******REMOVED***$el.Title***REMOVED******REMOVED***
	</tr>
	***REMOVED******REMOVED***if $.Expanded***REMOVED******REMOVED***
	<tr>
		<td class="when"></td>
		<td class="elapsed"></td>
		<td><pre>***REMOVED******REMOVED***$el.Stack|trimSpace***REMOVED******REMOVED***</pre></td>
	</tr>
	***REMOVED******REMOVED***range $el.Events***REMOVED******REMOVED***
	<tr>
		<td class="when">***REMOVED******REMOVED***.WhenString***REMOVED******REMOVED***</td>
		<td class="elapsed">***REMOVED******REMOVED***elapsed .Elapsed***REMOVED******REMOVED***</td>
		<td>.***REMOVED******REMOVED***if .IsErr***REMOVED******REMOVED***E***REMOVED******REMOVED***else***REMOVED******REMOVED***.***REMOVED******REMOVED***end***REMOVED******REMOVED***. ***REMOVED******REMOVED***.What***REMOVED******REMOVED***</td>
	</tr>
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED******REMOVED***end***REMOVED******REMOVED***
</table>
***REMOVED******REMOVED***end***REMOVED******REMOVED***
	</body>
</html>
`
