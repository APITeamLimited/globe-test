package orchMetrics

/*func calculateTotal(wrappedEnvelopes []*wrappedEnvelope, envelope *wrappedEnvelope, currentTime time.Time) []*wrappedEnvelope {
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

	// Counter is a cumulative metric, so we need to add the value to the existing value
	accEnvelope.Data.Value += envelope.Data.Value

	if envelope.Metric.Tainted.Bool {
		accEnvelope.Metric.Tainted = null.BoolFrom(true)
	}

	return wrappedEnvelopes
}
*/
