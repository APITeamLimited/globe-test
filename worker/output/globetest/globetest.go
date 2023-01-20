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

	flushCount      int
	flushCountMutex sync.Mutex
}

type WrappedFormattedSamples struct {
	SampleEnvelopes []SampleEnvelope `json:"samples"`
	FlushCount      int              `json:"flush_count"`
}

func New(gs libWorker.BaseGlobalState) (output.Output, error) {
	return &Output{
		gs:          gs,
		seenMetrics: make(map[string]struct{}),

		flushCount:      0,
		flushCountMutex: sync.Mutex{},
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
		o.flushCountMutex.Lock()
		o.flushCount++
		o.flushCountMutex.Unlock()

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

	marshalledWrappedSamples, err := json.Marshal(WrappedFormattedSamples{
		SampleEnvelopes: formattedSamples,
		FlushCount:      o.flushCount,
	})

	if err != nil {
		libWorker.HandleError(o.gs, err)
		return
	}

	libWorker.DispatchMessage(o.gs, string(marshalledWrappedSamples), "METRICS")
}
