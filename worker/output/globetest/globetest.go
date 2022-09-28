package globetest

import (
	"fmt"
	"strings"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/APITeamLimited/globe-test/worker/output"
	jwriter "github.com/mailru/easyjson/jwriter"
)

const flushPeriod = 200 * time.Millisecond

type Output struct {
	output.SampleBuffer

	periodicFlusher *output.PeriodicFlusher

	workerInfo  *libWorker.WorkerInfo
	seenMetrics map[string]struct{}
	thresholds  map[string]metrics.Thresholds
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
func (o *Output) SetThresholds(thresholds map[string]metrics.Thresholds) {
	if len(thresholds) == 0 {
		return
	}
	o.thresholds = make(map[string]metrics.Thresholds, len(thresholds))
	for name, t := range thresholds {
		o.thresholds[name] = t
	}
}

func (o *Output) flushMetrics() {
	samples := o.GetBufferedSamples()
	var count int
	jw := new(jwriter.Writer)
	for _, sc := range samples {
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples {
			sample := sample
			o.handleMetric(sample.Metric, jw)
			wrapSample(sample).MarshalEasyJSON(jw)
			jw.RawByte('\n')
		}
	}

	buffer := jw.Buffer.BuildBytes()

	formatted := strings.ReplaceAll(fmt.Sprintf("[%s]", strings.ReplaceAll(string(buffer), "\n", ",")), ",]", "]")

	if count > 0 {
		libWorker.DispatchMessage(o.workerInfo.Ctx, o.workerInfo.Client, o.workerInfo.JobId, o.workerInfo.WorkerId, formatted, "METRICS")
	}
}

func (o *Output) handleMetric(m *metrics.Metric, jw *jwriter.Writer) {
	if _, ok := o.seenMetrics[m.Name]; ok {
		return
	}
	o.seenMetrics[m.Name] = struct{}{}

	wrapped := metricEnvelope{
		Type:   "Metric",
		Metric: m.Name,
	}
	wrapped.Data.Name = m.Name
	wrapped.Data.Type = m.Type
	wrapped.Data.Contains = m.Contains
	wrapped.Data.Submetrics = m.Submetrics

	if ts, ok := o.thresholds[m.Name]; ok {
		wrapped.Data.Thresholds = ts
	}

	wrapped.MarshalEasyJSON(jw)
	jw.RawByte('\n')
}
