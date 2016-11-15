package lib

import (
	// "math"
	"time"
)

// StageAt returns the stage at the specified offset (in nanoseconds) and the time remaining of
// said stage. If the interval is past the end of the test, an empty stage and 0 is returned.
func StageAt(stages []Stage, offset time.Duration) (s Stage, stageLeft time.Duration, ok bool) ***REMOVED***
	var counter time.Duration
	for _, stage := range stages ***REMOVED***
		counter += time.Duration(stage.Duration.Int64)
		if counter >= offset ***REMOVED***
			return stage, counter - offset, true
		***REMOVED***
	***REMOVED***
	return Stage***REMOVED******REMOVED***, 0, false
***REMOVED***

// Ease eases a value x towards y over time, so that: f=f(t) : f(tx)=x, f(ty)=y.
func Ease(t, tx, ty, x, y int64) int64 ***REMOVED***
	return x*(ty-t)/(ty-tx) + y*(t-tx)/(ty-tx)
***REMOVED***
