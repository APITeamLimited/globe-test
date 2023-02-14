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
}

type FormattedSamples struct {
	Samples    []workerMetrics.Sample `json:"samples" protobuf:"bytes,1"`
	FlushCount int                    `json:"flushCount" protobuf:"varint,2"`
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
		o.flushIncrementMutex.Unlock()
	}()

	sampleContainers := o.GetBufferedSamples()

	samples := make([]workerMetrics.Sample, 0)
	count := 0

	for _, sc := range sampleContainers {
		samples = append(samples, sc.GetSamples()...)
		count += len(samples)
	}

	marshalledSamples, err := json.Marshal(FormattedSamples{
		Samples:    samples,
		FlushCount: o.flushCount,
	})

	if err != nil {
		libWorker.HandleError(o.gs, err)
		return
	}
	//fmt.Print("\n\n\n\n")

	if o.flushCount == 15 {
		fmt.Println((string(marshalledSamples)))
	}

	//printSize(marshalledSamples)
	libWorker.DispatchMessage(o.gs, string(marshalledSamples), "METRICS")
}

func printSize(butes []byte) {
	// Check if b, kb or mb

	if len(butes) > 1000000 {
		fmt.Printf("%f mb", float64(len(butes))/1000000)
		return

	} else if len(butes) > 1000 {
		fmt.Printf("%f kb", float64(len(butes))/1000)
		return
	}

	fmt.Printf("%d bytes", len(butes))
}
