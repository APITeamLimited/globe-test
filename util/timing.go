package util

import (
	"runtime"
	"time"
)

func Time(fn func()) time.Duration ***REMOVED***
	m := runtime.MemStats***REMOVED******REMOVED***
	runtime.ReadMemStats(&m)

	numGC1 := m.NumGC

	startTime := time.Now()
	fn()
	duration := time.Since(startTime)

	runtime.ReadMemStats(&m)
	numGC2 := m.NumGC

	gcTotal := uint64(0)
	for i := numGC1; i < numGC2; i++ ***REMOVED***
		gcTotal += m.PauseNs[(i+255)%256]
	***REMOVED***

	return duration - time.Duration(gcTotal)
***REMOVED***
