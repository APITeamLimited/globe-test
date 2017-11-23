package humanize

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Seconds-based time units
const (
	Day      = 24 * time.Hour
	Week     = 7 * Day
	Month    = 30 * Day
	Year     = 12 * Month
	LongTime = 37 * Year
)

// Time formats a time into a relative string.
//
// Time(someT) -> "3 weeks ago"
func Time(then time.Time) string ***REMOVED***
	return RelTime(then, time.Now(), "ago", "from now")
***REMOVED***

// A RelTimeMagnitude struct contains a relative time point at which
// the relative format of time will switch to a new format string.  A
// slice of these in ascending order by their "D" field is passed to
// CustomRelTime to format durations.
//
// The Format field is a string that may contain a "%s" which will be
// replaced with the appropriate signed label (e.g. "ago" or "from
// now") and a "%d" that will be replaced by the quantity.
//
// The DivBy field is the amount of time the time difference must be
// divided by in order to display correctly.
//
// e.g. if D is 2*time.Minute and you want to display "%d minutes %s"
// DivBy should be time.Minute so whatever the duration is will be
// expressed in minutes.
type RelTimeMagnitude struct ***REMOVED***
	D      time.Duration
	Format string
	DivBy  time.Duration
***REMOVED***

var defaultMagnitudes = []RelTimeMagnitude***REMOVED***
	***REMOVED***time.Second, "now", time.Second***REMOVED***,
	***REMOVED***2 * time.Second, "1 second %s", 1***REMOVED***,
	***REMOVED***time.Minute, "%d seconds %s", time.Second***REMOVED***,
	***REMOVED***2 * time.Minute, "1 minute %s", 1***REMOVED***,
	***REMOVED***time.Hour, "%d minutes %s", time.Minute***REMOVED***,
	***REMOVED***2 * time.Hour, "1 hour %s", 1***REMOVED***,
	***REMOVED***Day, "%d hours %s", time.Hour***REMOVED***,
	***REMOVED***2 * Day, "1 day %s", 1***REMOVED***,
	***REMOVED***Week, "%d days %s", Day***REMOVED***,
	***REMOVED***2 * Week, "1 week %s", 1***REMOVED***,
	***REMOVED***Month, "%d weeks %s", Week***REMOVED***,
	***REMOVED***2 * Month, "1 month %s", 1***REMOVED***,
	***REMOVED***Year, "%d months %s", Month***REMOVED***,
	***REMOVED***18 * Month, "1 year %s", 1***REMOVED***,
	***REMOVED***2 * Year, "2 years %s", 1***REMOVED***,
	***REMOVED***LongTime, "%d years %s", Year***REMOVED***,
	***REMOVED***math.MaxInt64, "a long while %s", 1***REMOVED***,
***REMOVED***

// RelTime formats a time into a relative string.
//
// It takes two times and two labels.  In addition to the generic time
// delta string (e.g. 5 minutes), the labels are used applied so that
// the label corresponding to the smaller time is applied.
//
// RelTime(timeInPast, timeInFuture, "earlier", "later") -> "3 weeks earlier"
func RelTime(a, b time.Time, albl, blbl string) string ***REMOVED***
	return CustomRelTime(a, b, albl, blbl, defaultMagnitudes)
***REMOVED***

// CustomRelTime formats a time into a relative string.
//
// It takes two times two labels and a table of relative time formats.
// In addition to the generic time delta string (e.g. 5 minutes), the
// labels are used applied so that the label corresponding to the
// smaller time is applied.
func CustomRelTime(a, b time.Time, albl, blbl string, magnitudes []RelTimeMagnitude) string ***REMOVED***
	lbl := albl
	diff := b.Sub(a)

	if a.After(b) ***REMOVED***
		lbl = blbl
		diff = a.Sub(b)
	***REMOVED***

	n := sort.Search(len(magnitudes), func(i int) bool ***REMOVED***
		return magnitudes[i].D > diff
	***REMOVED***)

	if n >= len(magnitudes) ***REMOVED***
		n = len(magnitudes) - 1
	***REMOVED***
	mag := magnitudes[n]
	args := []interface***REMOVED******REMOVED******REMOVED******REMOVED***
	escaped := false
	for _, ch := range mag.Format ***REMOVED***
		if escaped ***REMOVED***
			switch ch ***REMOVED***
			case 's':
				args = append(args, lbl)
			case 'd':
				args = append(args, diff/mag.DivBy)
			***REMOVED***
			escaped = false
		***REMOVED*** else ***REMOVED***
			escaped = ch == '%'
		***REMOVED***
	***REMOVED***
	return fmt.Sprintf(mag.Format, args...)
***REMOVED***
