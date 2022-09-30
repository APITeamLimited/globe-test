package workerMetrics

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

// B formats a boolean value for emission.
func B(b bool) float64 ***REMOVED***
	if b ***REMOVED***
		return 1
	***REMOVED***
	return 0
***REMOVED***
