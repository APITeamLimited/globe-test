package stats

import (
	"time"
)

const timeUnit = time.Millisecond

// D formats a duration for emission.
// The reverse of D() is ToD().
func D(d time.Duration) float64 ***REMOVED***
	return float64(d) / float64(timeUnit)
***REMOVED***

// ToD converts an emitted duration to a time.Duration.
// The reverse of ToD() is D().
func ToD(d float64) time.Duration ***REMOVED***
	return time.Duration(d * float64(timeUnit))
***REMOVED***
