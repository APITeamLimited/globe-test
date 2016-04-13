package util

import (
	"runtime"
	"time"
)

func Time(fn func()) time.Duration ***REMOVED***
	m := runtime.MemStats***REMOVED******REMOVED***
	runtime.ReadMemStats(&m)

	gcTotal1 := m.PauseTotalNs

	startTime := time.Now()
	fn()
	duration := time.Since(startTime)

	runtime.ReadMemStats(&m)
	gcTotal2 := m.PauseTotalNs

	gcTotal := time.Duration(gcTotal2 - gcTotal1)
	return duration - gcTotal
***REMOVED***
