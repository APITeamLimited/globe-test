package orchMetrics

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/output/globetest"
	"gopkg.in/guregu/null.v3"
)

func calculateTotal(wrappedEnvelopes []*wrappedEnvelope, envelope *wrappedEnvelope, currentTime time.Time) []*wrappedEnvelope ***REMOVED***
	// Check if wrapped envelope exists in the slice
	metricName := envelope.Metric.Name
	accEnvelope := &wrappedEnvelope***REMOVED******REMOVED***

	for _, wrappedEnvelopeCurrent := range wrappedEnvelopes ***REMOVED***
		if wrappedEnvelopeCurrent.Metric.Name == metricName ***REMOVED***
			accEnvelope = wrappedEnvelopeCurrent
			break
		***REMOVED***
	***REMOVED***

	// If wrapped envelope does not exist, create a new one
	if accEnvelope.Metric == nil ***REMOVED***
		accEnvelope := wrappedEnvelope***REMOVED***
			SampleEnvelope: globetest.SampleEnvelope***REMOVED***
				Type: "Point",
				Data: globetest.SampleData***REMOVED***
					Time:  currentTime,
					Value: 0,
				***REMOVED***,
				Metric: envelope.Metric,
			***REMOVED***,
			workerId: envelope.workerId,
			location: envelope.location,
		***REMOVED***

		wrappedEnvelopes = append(wrappedEnvelopes, &accEnvelope)
	***REMOVED***

	// Counter is a cumulative metric, so we need to add the value to the existing value
	accEnvelope.Data.Value += envelope.Data.Value

	if envelope.Metric.Tainted.Bool ***REMOVED***
		accEnvelope.Metric.Tainted = null.BoolFrom(true)
	***REMOVED***

	return wrappedEnvelopes
***REMOVED***
