package executor

import "go.k6.io/k6/metrics"

func sumMetricValues(samples chan metrics.SampleContainer, metricName string) (sum float64) ***REMOVED*** //nolint:unparam
	for _, sc := range metrics.GetBufferedSamples(samples) ***REMOVED***
		samples := sc.GetSamples()
		for _, s := range samples ***REMOVED***
			if s.Metric.Name == metricName ***REMOVED***
				sum += s.Value
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return sum
***REMOVED***
