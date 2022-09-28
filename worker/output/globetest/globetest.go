package globetest

import (
	"fmt"
	"strings"
	"time"

	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/metrics"
	"github.com/APITeamLimited/k6-worker/output"
	jwriter "github.com/mailru/easyjson/jwriter"
)

const flushPeriod = 200 * time.Millisecond

type Output struct ***REMOVED***
	output.SampleBuffer

	periodicFlusher *output.PeriodicFlusher

	params      output.Params
	workerInfo  *lib.WorkerInfo
	seenMetrics map[string]struct***REMOVED******REMOVED***
	thresholds  map[string]metrics.Thresholds
***REMOVED***

func New(params output.Params, workerInfo *lib.WorkerInfo) (output.Output, error) ***REMOVED***
	return &Output***REMOVED***
		params:      params,
		workerInfo:  workerInfo,
		seenMetrics: make(map[string]struct***REMOVED******REMOVED***),
	***REMOVED***, nil
***REMOVED***

func (o *Output) Description() string ***REMOVED***
	return fmt.Sprintf("GlobeTest output for job %s", o.workerInfo.JobId)
***REMOVED***

func (o *Output) Start() error ***REMOVED***
	pf, err := output.NewPeriodicFlusher(flushPeriod, o.flushMetrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.periodicFlusher = pf

	return nil
***REMOVED***

func (o *Output) Stop() error ***REMOVED***
	o.periodicFlusher.Stop()
	return nil
***REMOVED***

// SetThresholds receives the thresholds before the output is Start()-ed.
func (o *Output) SetThresholds(thresholds map[string]metrics.Thresholds) ***REMOVED***
	if len(thresholds) == 0 ***REMOVED***
		return
	***REMOVED***
	o.thresholds = make(map[string]metrics.Thresholds, len(thresholds))
	for name, t := range thresholds ***REMOVED***
		o.thresholds[name] = t
	***REMOVED***
***REMOVED***

func (o *Output) flushMetrics() ***REMOVED***
	samples := o.GetBufferedSamples()
	var count int
	jw := new(jwriter.Writer)
	for _, sc := range samples ***REMOVED***
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples ***REMOVED***
			sample := sample
			o.handleMetric(sample.Metric, jw)
			wrapSample(sample).MarshalEasyJSON(jw)
			jw.RawByte('\n')
		***REMOVED***
	***REMOVED***

	buffer := jw.Buffer.BuildBytes()

	formatted := strings.ReplaceAll(fmt.Sprintf("[%s]", strings.ReplaceAll(string(buffer), "\n", ",")), ",]", "]")

	if count > 0 ***REMOVED***
		lib.DispatchMessage(o.workerInfo.Ctx, o.workerInfo.Client, o.workerInfo.JobId, o.workerInfo.WorkerId, formatted, "METRICS")
	***REMOVED***
***REMOVED***

func (o *Output) handleMetric(m *metrics.Metric, jw *jwriter.Writer) ***REMOVED***
	if _, ok := o.seenMetrics[m.Name]; ok ***REMOVED***
		return
	***REMOVED***
	o.seenMetrics[m.Name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	wrapped := metricEnvelope***REMOVED***
		Type:   "Metric",
		Metric: m.Name,
	***REMOVED***
	wrapped.Data.Name = m.Name
	wrapped.Data.Type = m.Type
	wrapped.Data.Contains = m.Contains
	wrapped.Data.Submetrics = m.Submetrics

	if ts, ok := o.thresholds[m.Name]; ok ***REMOVED***
		wrapped.Data.Thresholds = ts
	***REMOVED***

	wrapped.MarshalEasyJSON(jw)
	jw.RawByte('\n')
***REMOVED***
