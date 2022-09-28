package stats

import (
	"sort"
	"time"
)

// float64ToInt rounds a float64 to an int
func float64ToInt(input float64) (output int) ***REMOVED***
	r, _ := Round(input, 0)
	return int(r)
***REMOVED***

// unixnano returns nanoseconds from UTC epoch
func unixnano() int64 ***REMOVED***
	return time.Now().UTC().UnixNano()
***REMOVED***

// copyslice copies a slice of float64s
func copyslice(input Float64Data) Float64Data ***REMOVED***
	s := make(Float64Data, input.Len())
	copy(s, input)
	return s
***REMOVED***

// sortedCopy returns a sorted copy of float64s
func sortedCopy(input Float64Data) (copy Float64Data) ***REMOVED***
	copy = copyslice(input)
	sort.Float64s(copy)
	return
***REMOVED***

// sortedCopyDif returns a sorted copy of float64s
// only if the original data isn't sorted.
// Only use this if returned slice won't be manipulated!
func sortedCopyDif(input Float64Data) (copy Float64Data) ***REMOVED***
	if sort.Float64sAreSorted(input) ***REMOVED***
		return input
	***REMOVED***
	copy = copyslice(input)
	sort.Float64s(copy)
	return
***REMOVED***
