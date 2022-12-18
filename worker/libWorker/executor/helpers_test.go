package executor

import "github.com/APITeamLimited/globe-test/worker/workerMetrics"

func sumMetricValues(samples chan workerMetrics.SampleContainer, metricName string) (sum float64) { //nolint:unparam
	for _, sc := range workerMetrics.GetBufferedSamples(samples) {
		samples := sc.GetSamples()
		for _, s := range samples {
			if s.Metric.Name == metricName {
				sum += s.Value
			}
		}
	}
	return sum
}
