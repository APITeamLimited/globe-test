// +build go1.5

package metrics

import "runtime"

func gcCPUFraction(memStats *runtime.MemStats) float64 ***REMOVED***
	return memStats.GCCPUFraction
***REMOVED***
