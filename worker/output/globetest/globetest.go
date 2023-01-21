package globetest

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

const flushPeriod = 1000 * time.Millisecond

type Output struct {
	output.SampleBuffer

	periodicFlusher *output.PeriodicFlusher

	gs          libWorker.BaseGlobalState
	seenMetrics map[string]struct{}
	thresholds  map[string]workerMetrics.Thresholds

	flushCount          int
	flushIncrementMutex sync.Mutex

	// In case 2 sets of metrics are sent in the same flush, we need to know if
	// we've already received the vu count so it can be added to the next flush
	// instead
	// 	vuMetrics1 *SampleEnvelope
	// 	vuMetrics2 *SampleEnvelope
}

type WrappedFormattedSamples struct {
	SampleEnvelopes []SampleEnvelope `json:"samples"`
	FlushCount      int              `json:"flush_count"`
}

func New(gs libWorker.BaseGlobalState) (output.Output, error) {
	return &Output{
		gs:          gs,
		seenMetrics: make(map[string]struct{}),

		flushCount:          0,
		flushIncrementMutex: sync.Mutex{},
	}, nil
}

func (o *Output) Description() string {
	return fmt.Sprintf("GlobeTest output for job %s", o.gs.JobId())
}

func (o *Output) Start() error {
	pf, err := output.NewPeriodicFlusher(flushPeriod, o.flushMetrics)
	if err != nil {
		return err
	}
	o.periodicFlusher = pf

	// Flush once immediately

	return nil
}

func (o *Output) Stop() error {
	o.periodicFlusher.Stop()
	return nil
}

// SetThresholds receives the thresholds before the output is Start()-ed.
func (o *Output) SetThresholds(thresholds map[string]workerMetrics.Thresholds) {
	if len(thresholds) == 0 {
		return
	}
	o.thresholds = make(map[string]workerMetrics.Thresholds, len(thresholds))
	for name, t := range thresholds {
		o.thresholds[name] = t
	}
}

func (o *Output) flushMetrics() {
	defer func() {
		o.flushIncrementMutex.Lock()
		o.flushCount++

		// if o.vuMetrics2 != nil {
		// 	o.vuMetrics1 = o.vuMetrics2
		// 	o.vuMetrics2 = nil
		// } else {
		// 	o.vuMetrics1 = nil
		// }

		o.flushIncrementMutex.Unlock()
	}()

	samples := o.GetBufferedSamples()
	var count int

	formattedSamples := make([]SampleEnvelope, 0)

	for _, sc := range samples {
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples {
			wrapped := wrapSample(sample)

			formattedSamples = append(formattedSamples, wrapped)
		}
	}

	// Get vu count from formattedSamples
	/*var foundVus *SampleEnvelope
	for _, sample := range formattedSamples {
		if sample.Metric.Name == "vus" {
			foundVus = &sample
		}
	}

	// If formattedSamples contains vus and vuMetrics is nil, set vuMetrics to

	if o.vuMetrics1 == nil && foundVus != nil {
		o.vuMetrics1 = foundVus
	} else if o.vuMetrics1 != nil && foundVus != nil {
		o.vuMetrics2 = foundVus

		// Remove the vus from formattedSamples
		for i, sample := range formattedSamples {
			if sample.Metric.Name == "vus" {
				formattedSamples = append(formattedSamples[:i], formattedSamples[i+1:]...)
				break
			}
		}
	} else if foundVus == nil && o.vuMetrics1 != nil {
		// Add the vuMetrics to formattedSamples
		formattedSamples = append(formattedSamples, *o.vuMetrics1)
	}*/

	marshalledWrappedSamples, err := json.Marshal(WrappedFormattedSamples{
		SampleEnvelopes: aggregateSampleEnvelopes(formattedSamples),
		FlushCount:      o.flushCount,
	})

	libWorker.DispatchMessage(o.gs, string(marshalledWrappedSamples), "METRICS")
}
