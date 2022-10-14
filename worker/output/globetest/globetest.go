package globetest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

const flushPeriod = 200 * time.Millisecond

type Output struct {
	output.SampleBuffer

	periodicFlusher *output.PeriodicFlusher

	workerInfo  *libWorker.WorkerInfo
	seenMetrics map[string]struct{}
	thresholds  map[string]workerMetrics.Thresholds
}

func New(workerInfo *libWorker.WorkerInfo) (output.Output, error) {
	return &Output{
		workerInfo:  workerInfo,
		seenMetrics: make(map[string]struct{}),
	}, nil
}

func (o *Output) Description() string {
	return fmt.Sprintf("GlobeTest output for job %s", o.workerInfo.JobId)
}

func (o *Output) Start() error {
	pf, err := output.NewPeriodicFlusher(flushPeriod, o.flushMetrics)
	if err != nil {
		return err
	}
	o.periodicFlusher = pf

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
	samples := o.GetBufferedSamples()
	var count int

	formattedSamples := make([]SampleEnvelope, 0)

	for _, sc := range samples {
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples {
			sample := sample

			wrapped := wrapSample(sample)

			formattedSamples = append(formattedSamples, wrapped)
		}
	}

	if len(formattedSamples) > 0 {
		marshalled, err := json.Marshal(formattedSamples)
		if err != nil {
			libWorker.HandleError(o.workerInfo.Ctx, o.workerInfo.Client, o.workerInfo.JobId, o.workerInfo.WorkerId, err)
			return
		}

		libWorker.DispatchMessage(o.workerInfo.Ctx, o.workerInfo.Client, o.workerInfo.JobId, o.workerInfo.WorkerId, string(marshalled), "METRICS")
	}
}
