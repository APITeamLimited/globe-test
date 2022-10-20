package orchMetrics

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/output/globetest"
)

func calculateMean(wrappedEnvelopes []*wrappedEnvelope, envelope *wrappedEnvelope, currentTime time.Time, vuCount int) []*wrappedEnvelope {
	// Zero weighting has no effect on the rate
	if vuCount == 0 {
		return wrappedEnvelopes
	}

	// Check if wrapped envelope exists in the slice
	metricName := envelope.Metric.Name
	accEnvelope := &wrappedEnvelope{}

	for _, wrappedEnvelopeCurrent := range wrappedEnvelopes {
		if wrappedEnvelopeCurrent.Metric.Name == metricName {
			accEnvelope = wrappedEnvelopeCurrent
			break
		}
	}

	// If wrapped envelope does not exist, create a new one
	if accEnvelope.Metric == nil {
		accEnvelope := wrappedEnvelope{
			SampleEnvelope: globetest.SampleEnvelope{
				Type: "Point",
				Data: globetest.SampleData{
					Time:  currentTime,
					Value: 0,
				},
				Metric: envelope.Metric,
			},
			workerId: envelope.workerId,
			location: envelope.location,
		}

		wrappedEnvelopes = append(wrappedEnvelopes, &accEnvelope)
	}

	// Rate tracks the percentage of values that are non-zero
	// Perform weighted average

	weightingOld := accEnvelope.weighting
	weightingEnvelope := envelope.weighting

	rateOld := accEnvelope.Data.Value
	rateEnvelope := envelope.Data.Value

	weightingNew := weightingOld + weightingEnvelope

	// Cannot divide by zero
	if weightingNew == 0 {
		return wrappedEnvelopes
	}

	rateNew := (rateOld*float64(weightingOld) + rateEnvelope*float64(weightingEnvelope)) / float64(weightingNew)

	accEnvelope.weighting = weightingNew
	accEnvelope.Data.Value = rateNew

	return wrappedEnvelopes
}
