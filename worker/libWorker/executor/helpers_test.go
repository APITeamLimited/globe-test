package executor

import "github.com/APITeamLimited/globe-test/worker/workerMetrics"

func sumMetricValues(samples chan workerMetrics.SampleContainer, metricName string) (sum float64) ***REMOVED*** //nolint:unparam
	for _, sc := range workerMetrics.GetBufferedSamples(samples) ***REMOVED***
		samples := sc.GetSamples()
		for _, s := range samples ***REMOVED***
			if s.Metric.Name == metricName ***REMOVED***
				sum += s.Value
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return sum
***REMOVED***
