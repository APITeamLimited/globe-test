// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package timeseries implements a time series structure for stats collection.
package timeseries // import "golang.org/x/net/internal/timeseries"

import (
	"fmt"
	"log"
	"time"
)

const (
	timeSeriesNumBuckets       = 64
	minuteHourSeriesNumBuckets = 60
)

var timeSeriesResolutions = []time.Duration***REMOVED***
	1 * time.Second,
	10 * time.Second,
	1 * time.Minute,
	10 * time.Minute,
	1 * time.Hour,
	6 * time.Hour,
	24 * time.Hour,          // 1 day
	7 * 24 * time.Hour,      // 1 week
	4 * 7 * 24 * time.Hour,  // 4 weeks
	16 * 7 * 24 * time.Hour, // 16 weeks
***REMOVED***

var minuteHourSeriesResolutions = []time.Duration***REMOVED***
	1 * time.Second,
	1 * time.Minute,
***REMOVED***

// An Observable is a kind of data that can be aggregated in a time series.
type Observable interface ***REMOVED***
	Multiply(ratio float64)    // Multiplies the data in self by a given ratio
	Add(other Observable)      // Adds the data from a different observation to self
	Clear()                    // Clears the observation so it can be reused.
	CopyFrom(other Observable) // Copies the contents of a given observation to self
***REMOVED***

// Float attaches the methods of Observable to a float64.
type Float float64

// NewFloat returns a Float.
func NewFloat() Observable ***REMOVED***
	f := Float(0)
	return &f
***REMOVED***

// String returns the float as a string.
func (f *Float) String() string ***REMOVED*** return fmt.Sprintf("%g", f.Value()) ***REMOVED***

// Value returns the float's value.
func (f *Float) Value() float64 ***REMOVED*** return float64(*f) ***REMOVED***

func (f *Float) Multiply(ratio float64) ***REMOVED*** *f *= Float(ratio) ***REMOVED***

func (f *Float) Add(other Observable) ***REMOVED***
	o := other.(*Float)
	*f += *o
***REMOVED***

func (f *Float) Clear() ***REMOVED*** *f = 0 ***REMOVED***

func (f *Float) CopyFrom(other Observable) ***REMOVED***
	o := other.(*Float)
	*f = *o
***REMOVED***

// A Clock tells the current time.
type Clock interface ***REMOVED***
	Time() time.Time
***REMOVED***

type defaultClock int

var defaultClockInstance defaultClock

func (defaultClock) Time() time.Time ***REMOVED*** return time.Now() ***REMOVED***

// Information kept per level. Each level consists of a circular list of
// observations. The start of the level may be derived from end and the
// len(buckets) * sizeInMillis.
type tsLevel struct ***REMOVED***
	oldest   int               // index to oldest bucketed Observable
	newest   int               // index to newest bucketed Observable
	end      time.Time         // end timestamp for this level
	size     time.Duration     // duration of the bucketed Observable
	buckets  []Observable      // collections of observations
	provider func() Observable // used for creating new Observable
***REMOVED***

func (l *tsLevel) Clear() ***REMOVED***
	l.oldest = 0
	l.newest = len(l.buckets) - 1
	l.end = time.Time***REMOVED******REMOVED***
	for i := range l.buckets ***REMOVED***
		if l.buckets[i] != nil ***REMOVED***
			l.buckets[i].Clear()
			l.buckets[i] = nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (l *tsLevel) InitLevel(size time.Duration, numBuckets int, f func() Observable) ***REMOVED***
	l.size = size
	l.provider = f
	l.buckets = make([]Observable, numBuckets)
***REMOVED***

// Keeps a sequence of levels. Each level is responsible for storing data at
// a given resolution. For example, the first level stores data at a one
// minute resolution while the second level stores data at a one hour
// resolution.

// Each level is represented by a sequence of buckets. Each bucket spans an
// interval equal to the resolution of the level. New observations are added
// to the last bucket.
type timeSeries struct ***REMOVED***
	provider    func() Observable // make more Observable
	numBuckets  int               // number of buckets in each level
	levels      []*tsLevel        // levels of bucketed Observable
	lastAdd     time.Time         // time of last Observable tracked
	total       Observable        // convenient aggregation of all Observable
	clock       Clock             // Clock for getting current time
	pending     Observable        // observations not yet bucketed
	pendingTime time.Time         // what time are we keeping in pending
	dirty       bool              // if there are pending observations
***REMOVED***

// init initializes a level according to the supplied criteria.
func (ts *timeSeries) init(resolutions []time.Duration, f func() Observable, numBuckets int, clock Clock) ***REMOVED***
	ts.provider = f
	ts.numBuckets = numBuckets
	ts.clock = clock
	ts.levels = make([]*tsLevel, len(resolutions))

	for i := range resolutions ***REMOVED***
		if i > 0 && resolutions[i-1] >= resolutions[i] ***REMOVED***
			log.Print("timeseries: resolutions must be monotonically increasing")
			break
		***REMOVED***
		newLevel := new(tsLevel)
		newLevel.InitLevel(resolutions[i], ts.numBuckets, ts.provider)
		ts.levels[i] = newLevel
	***REMOVED***

	ts.Clear()
***REMOVED***

// Clear removes all observations from the time series.
func (ts *timeSeries) Clear() ***REMOVED***
	ts.lastAdd = time.Time***REMOVED******REMOVED***
	ts.total = ts.resetObservation(ts.total)
	ts.pending = ts.resetObservation(ts.pending)
	ts.pendingTime = time.Time***REMOVED******REMOVED***
	ts.dirty = false

	for i := range ts.levels ***REMOVED***
		ts.levels[i].Clear()
	***REMOVED***
***REMOVED***

// Add records an observation at the current time.
func (ts *timeSeries) Add(observation Observable) ***REMOVED***
	ts.AddWithTime(observation, ts.clock.Time())
***REMOVED***

// AddWithTime records an observation at the specified time.
func (ts *timeSeries) AddWithTime(observation Observable, t time.Time) ***REMOVED***

	smallBucketDuration := ts.levels[0].size

	if t.After(ts.lastAdd) ***REMOVED***
		ts.lastAdd = t
	***REMOVED***

	if t.After(ts.pendingTime) ***REMOVED***
		ts.advance(t)
		ts.mergePendingUpdates()
		ts.pendingTime = ts.levels[0].end
		ts.pending.CopyFrom(observation)
		ts.dirty = true
	***REMOVED*** else if t.After(ts.pendingTime.Add(-1 * smallBucketDuration)) ***REMOVED***
		// The observation is close enough to go into the pending bucket.
		// This compensates for clock skewing and small scheduling delays
		// by letting the update stay in the fast path.
		ts.pending.Add(observation)
		ts.dirty = true
	***REMOVED*** else ***REMOVED***
		ts.mergeValue(observation, t)
	***REMOVED***
***REMOVED***

// mergeValue inserts the observation at the specified time in the past into all levels.
func (ts *timeSeries) mergeValue(observation Observable, t time.Time) ***REMOVED***
	for _, level := range ts.levels ***REMOVED***
		index := (ts.numBuckets - 1) - int(level.end.Sub(t)/level.size)
		if 0 <= index && index < ts.numBuckets ***REMOVED***
			bucketNumber := (level.oldest + index) % ts.numBuckets
			if level.buckets[bucketNumber] == nil ***REMOVED***
				level.buckets[bucketNumber] = level.provider()
			***REMOVED***
			level.buckets[bucketNumber].Add(observation)
		***REMOVED***
	***REMOVED***
	ts.total.Add(observation)
***REMOVED***

// mergePendingUpdates applies the pending updates into all levels.
func (ts *timeSeries) mergePendingUpdates() ***REMOVED***
	if ts.dirty ***REMOVED***
		ts.mergeValue(ts.pending, ts.pendingTime)
		ts.pending = ts.resetObservation(ts.pending)
		ts.dirty = false
	***REMOVED***
***REMOVED***

// advance cycles the buckets at each level until the latest bucket in
// each level can hold the time specified.
func (ts *timeSeries) advance(t time.Time) ***REMOVED***
	if !t.After(ts.levels[0].end) ***REMOVED***
		return
	***REMOVED***
	for i := 0; i < len(ts.levels); i++ ***REMOVED***
		level := ts.levels[i]
		if !level.end.Before(t) ***REMOVED***
			break
		***REMOVED***

		// If the time is sufficiently far, just clear the level and advance
		// directly.
		if !t.Before(level.end.Add(level.size * time.Duration(ts.numBuckets))) ***REMOVED***
			for _, b := range level.buckets ***REMOVED***
				ts.resetObservation(b)
			***REMOVED***
			level.end = time.Unix(0, (t.UnixNano()/level.size.Nanoseconds())*level.size.Nanoseconds())
		***REMOVED***

		for t.After(level.end) ***REMOVED***
			level.end = level.end.Add(level.size)
			level.newest = level.oldest
			level.oldest = (level.oldest + 1) % ts.numBuckets
			ts.resetObservation(level.buckets[level.newest])
		***REMOVED***

		t = level.end
	***REMOVED***
***REMOVED***

// Latest returns the sum of the num latest buckets from the level.
func (ts *timeSeries) Latest(level, num int) Observable ***REMOVED***
	now := ts.clock.Time()
	if ts.levels[0].end.Before(now) ***REMOVED***
		ts.advance(now)
	***REMOVED***

	ts.mergePendingUpdates()

	result := ts.provider()
	l := ts.levels[level]
	index := l.newest

	for i := 0; i < num; i++ ***REMOVED***
		if l.buckets[index] != nil ***REMOVED***
			result.Add(l.buckets[index])
		***REMOVED***
		if index == 0 ***REMOVED***
			index = ts.numBuckets
		***REMOVED***
		index--
	***REMOVED***

	return result
***REMOVED***

// LatestBuckets returns a copy of the num latest buckets from level.
func (ts *timeSeries) LatestBuckets(level, num int) []Observable ***REMOVED***
	if level < 0 || level > len(ts.levels) ***REMOVED***
		log.Print("timeseries: bad level argument: ", level)
		return nil
	***REMOVED***
	if num < 0 || num >= ts.numBuckets ***REMOVED***
		log.Print("timeseries: bad num argument: ", num)
		return nil
	***REMOVED***

	results := make([]Observable, num)
	now := ts.clock.Time()
	if ts.levels[0].end.Before(now) ***REMOVED***
		ts.advance(now)
	***REMOVED***

	ts.mergePendingUpdates()

	l := ts.levels[level]
	index := l.newest

	for i := 0; i < num; i++ ***REMOVED***
		result := ts.provider()
		results[i] = result
		if l.buckets[index] != nil ***REMOVED***
			result.CopyFrom(l.buckets[index])
		***REMOVED***

		if index == 0 ***REMOVED***
			index = ts.numBuckets
		***REMOVED***
		index -= 1
	***REMOVED***
	return results
***REMOVED***

// ScaleBy updates observations by scaling by factor.
func (ts *timeSeries) ScaleBy(factor float64) ***REMOVED***
	for _, l := range ts.levels ***REMOVED***
		for i := 0; i < ts.numBuckets; i++ ***REMOVED***
			l.buckets[i].Multiply(factor)
		***REMOVED***
	***REMOVED***

	ts.total.Multiply(factor)
	ts.pending.Multiply(factor)
***REMOVED***

// Range returns the sum of observations added over the specified time range.
// If start or finish times don't fall on bucket boundaries of the same
// level, then return values are approximate answers.
func (ts *timeSeries) Range(start, finish time.Time) Observable ***REMOVED***
	return ts.ComputeRange(start, finish, 1)[0]
***REMOVED***

// Recent returns the sum of observations from the last delta.
func (ts *timeSeries) Recent(delta time.Duration) Observable ***REMOVED***
	now := ts.clock.Time()
	return ts.Range(now.Add(-delta), now)
***REMOVED***

// Total returns the total of all observations.
func (ts *timeSeries) Total() Observable ***REMOVED***
	ts.mergePendingUpdates()
	return ts.total
***REMOVED***

// ComputeRange computes a specified number of values into a slice using
// the observations recorded over the specified time period. The return
// values are approximate if the start or finish times don't fall on the
// bucket boundaries at the same level or if the number of buckets spanning
// the range is not an integral multiple of num.
func (ts *timeSeries) ComputeRange(start, finish time.Time, num int) []Observable ***REMOVED***
	if start.After(finish) ***REMOVED***
		log.Printf("timeseries: start > finish, %v>%v", start, finish)
		return nil
	***REMOVED***

	if num < 0 ***REMOVED***
		log.Printf("timeseries: num < 0, %v", num)
		return nil
	***REMOVED***

	results := make([]Observable, num)

	for _, l := range ts.levels ***REMOVED***
		if !start.Before(l.end.Add(-l.size * time.Duration(ts.numBuckets))) ***REMOVED***
			ts.extract(l, start, finish, num, results)
			return results
		***REMOVED***
	***REMOVED***

	// Failed to find a level that covers the desired range. So just
	// extract from the last level, even if it doesn't cover the entire
	// desired range.
	ts.extract(ts.levels[len(ts.levels)-1], start, finish, num, results)

	return results
***REMOVED***

// RecentList returns the specified number of values in slice over the most
// recent time period of the specified range.
func (ts *timeSeries) RecentList(delta time.Duration, num int) []Observable ***REMOVED***
	if delta < 0 ***REMOVED***
		return nil
	***REMOVED***
	now := ts.clock.Time()
	return ts.ComputeRange(now.Add(-delta), now, num)
***REMOVED***

// extract returns a slice of specified number of observations from a given
// level over a given range.
func (ts *timeSeries) extract(l *tsLevel, start, finish time.Time, num int, results []Observable) ***REMOVED***
	ts.mergePendingUpdates()

	srcInterval := l.size
	dstInterval := finish.Sub(start) / time.Duration(num)
	dstStart := start
	srcStart := l.end.Add(-srcInterval * time.Duration(ts.numBuckets))

	srcIndex := 0

	// Where should scanning start?
	if dstStart.After(srcStart) ***REMOVED***
		advance := dstStart.Sub(srcStart) / srcInterval
		srcIndex += int(advance)
		srcStart = srcStart.Add(advance * srcInterval)
	***REMOVED***

	// The i'th value is computed as show below.
	// interval = (finish/start)/num
	// i'th value = sum of observation in range
	//   [ start + i       * interval,
	//     start + (i + 1) * interval )
	for i := 0; i < num; i++ ***REMOVED***
		results[i] = ts.resetObservation(results[i])
		dstEnd := dstStart.Add(dstInterval)
		for srcIndex < ts.numBuckets && srcStart.Before(dstEnd) ***REMOVED***
			srcEnd := srcStart.Add(srcInterval)
			if srcEnd.After(ts.lastAdd) ***REMOVED***
				srcEnd = ts.lastAdd
			***REMOVED***

			if !srcEnd.Before(dstStart) ***REMOVED***
				srcValue := l.buckets[(srcIndex+l.oldest)%ts.numBuckets]
				if !srcStart.Before(dstStart) && !srcEnd.After(dstEnd) ***REMOVED***
					// dst completely contains src.
					if srcValue != nil ***REMOVED***
						results[i].Add(srcValue)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					// dst partially overlaps src.
					overlapStart := maxTime(srcStart, dstStart)
					overlapEnd := minTime(srcEnd, dstEnd)
					base := srcEnd.Sub(srcStart)
					fraction := overlapEnd.Sub(overlapStart).Seconds() / base.Seconds()

					used := ts.provider()
					if srcValue != nil ***REMOVED***
						used.CopyFrom(srcValue)
					***REMOVED***
					used.Multiply(fraction)
					results[i].Add(used)
				***REMOVED***

				if srcEnd.After(dstEnd) ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			srcIndex++
			srcStart = srcStart.Add(srcInterval)
		***REMOVED***
		dstStart = dstStart.Add(dstInterval)
	***REMOVED***
***REMOVED***

// resetObservation clears the content so the struct may be reused.
func (ts *timeSeries) resetObservation(observation Observable) Observable ***REMOVED***
	if observation == nil ***REMOVED***
		observation = ts.provider()
	***REMOVED*** else ***REMOVED***
		observation.Clear()
	***REMOVED***
	return observation
***REMOVED***

// TimeSeries tracks data at granularities from 1 second to 16 weeks.
type TimeSeries struct ***REMOVED***
	timeSeries
***REMOVED***

// NewTimeSeries creates a new TimeSeries using the function provided for creating new Observable.
func NewTimeSeries(f func() Observable) *TimeSeries ***REMOVED***
	return NewTimeSeriesWithClock(f, defaultClockInstance)
***REMOVED***

// NewTimeSeriesWithClock creates a new TimeSeries using the function provided for creating new Observable and the clock for
// assigning timestamps.
func NewTimeSeriesWithClock(f func() Observable, clock Clock) *TimeSeries ***REMOVED***
	ts := new(TimeSeries)
	ts.timeSeries.init(timeSeriesResolutions, f, timeSeriesNumBuckets, clock)
	return ts
***REMOVED***

// MinuteHourSeries tracks data at granularities of 1 minute and 1 hour.
type MinuteHourSeries struct ***REMOVED***
	timeSeries
***REMOVED***

// NewMinuteHourSeries creates a new MinuteHourSeries using the function provided for creating new Observable.
func NewMinuteHourSeries(f func() Observable) *MinuteHourSeries ***REMOVED***
	return NewMinuteHourSeriesWithClock(f, defaultClockInstance)
***REMOVED***

// NewMinuteHourSeriesWithClock creates a new MinuteHourSeries using the function provided for creating new Observable and the clock for
// assigning timestamps.
func NewMinuteHourSeriesWithClock(f func() Observable, clock Clock) *MinuteHourSeries ***REMOVED***
	ts := new(MinuteHourSeries)
	ts.timeSeries.init(minuteHourSeriesResolutions, f,
		minuteHourSeriesNumBuckets, clock)
	return ts
***REMOVED***

func (ts *MinuteHourSeries) Minute() Observable ***REMOVED***
	return ts.timeSeries.Latest(0, 60)
***REMOVED***

func (ts *MinuteHourSeries) Hour() Observable ***REMOVED***
	return ts.timeSeries.Latest(1, 60)
***REMOVED***

func minTime(a, b time.Time) time.Time ***REMOVED***
	if a.Before(b) ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func maxTime(a, b time.Time) time.Time ***REMOVED***
	if a.After(b) ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***
