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

type Output struct ***REMOVED***
	output.SampleBuffer

	periodicFlusher *output.PeriodicFlusher

	workerInfo  *libWorker.WorkerInfo
	seenMetrics map[string]struct***REMOVED******REMOVED***
	thresholds  map[string]workerMetrics.Thresholds
***REMOVED***

func New(workerInfo *libWorker.WorkerInfo) (output.Output, error) ***REMOVED***
	return &Output***REMOVED***
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
func (o *Output) SetThresholds(thresholds map[string]workerMetrics.Thresholds) ***REMOVED***
	if len(thresholds) == 0 ***REMOVED***
		return
	***REMOVED***
	o.thresholds = make(map[string]workerMetrics.Thresholds, len(thresholds))
	for name, t := range thresholds ***REMOVED***
		o.thresholds[name] = t
	***REMOVED***
***REMOVED***

func (o *Output) flushMetrics() ***REMOVED***
	samples := o.GetBufferedSamples()
	var count int

	formattedSamples := make([]SampleEnvelope, 0)

	for _, sc := range samples ***REMOVED***
		samples := sc.GetSamples()
		count += len(samples)
		for _, sample := range samples ***REMOVED***
			sample := sample

			wrapped := wrapSample(sample)

			formattedSamples = append(formattedSamples, wrapped)
		***REMOVED***
	***REMOVED***

	if len(formattedSamples) > 0 ***REMOVED***
		marshalled, err := json.Marshal(formattedSamples)
		if err != nil ***REMOVED***
			libWorker.HandleError(o.workerInfo.Ctx, o.workerInfo.Client, o.workerInfo.JobId, o.workerInfo.WorkerId, err)
			return
		***REMOVED***

		libWorker.DispatchMessage(o.workerInfo.Ctx, o.workerInfo.Client, o.workerInfo.JobId, o.workerInfo.WorkerId, string(marshalled), "METRICS")
	***REMOVED***
***REMOVED***
