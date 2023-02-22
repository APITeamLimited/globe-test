package globetest

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/aggregator"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/APITeamLimited/globe-test/worker/output"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	uniqueOutputLimit = 10000
	flushPeriod       = 1000 * time.Millisecond
)

type Output struct {
	output.SampleBuffer

	periodicFlusher *output.PeriodicFlusher

	gs       libWorker.BaseGlobalState
	location string

	seenMetrics map[string]struct{}
	thresholds  map[string]metrics.Thresholds

	flushCount          int
	flushIncrementMutex sync.Mutex

	uniqueConsoleMessages []string
	consoleMessages       []*aggregator.ConsoleMessage
	consoleMutex          sync.Mutex

	startTime time.Time
}

type FormattedSamples struct {
	Samples    []metrics.Sample `json:"samples" protobuf:"bytes,1"`
	FlushCount int              `json:"flushCount" protobuf:"varint,2"`
}

func New(gs libWorker.BaseGlobalState, location string) (output.Output, error) {
	return &Output{
		gs:       gs,
		location: location,

		seenMetrics: make(map[string]struct{}),

		flushCount:          0,
		flushIncrementMutex: sync.Mutex{},

		uniqueConsoleMessages: make([]string, 0),
		consoleMessages:       make([]*aggregator.ConsoleMessage, 0),
		consoleMutex:          sync.Mutex{},
	}, nil

}

func (o *Output) Start() error {
	o.startTime = time.Now()

	pf, err := output.NewPeriodicFlusher(flushPeriod, func() {
		o.flushMetrics()
		o.flushConsoleMessages()
	})
	if err != nil {
		return err
	}
	o.periodicFlusher = pf

	go func() {
		for log := range o.gs.GetLoggerChannel() {
			level := log["level"].(string)

			var msg string

			if log["error"] != nil {
				msg = log["error"].(string)
			} else {
				msg = log["msg"].(string)
			}

			consoleMessage := &aggregator.ConsoleMessage{
				Message: msg,
				Level:   level,
				Count: map[string]int32{
					o.location: 1,
				},
				FirstOccurred: timestamppb.New(time.Now()),
				LastOccurred:  timestamppb.New(time.Now()),
			}

			o.addConsoleMessage(consoleMessage)
		}
	}()

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
	defer func() {
		o.flushIncrementMutex.Lock()
		o.flushCount++
		o.flushIncrementMutex.Unlock()
	}()

	flushMetrics := make(map[string]metrics.Metric)

	var metricName string

	outputCount := 0

	samples := make([]metrics.Sample, 0)
	for _, sampleContainer := range o.GetBufferedSamples() {
		samples = append(samples, sampleContainer.GetSamples()...)
	}

	for _, sample := range samples {
		if outputCount > uniqueOutputLimit {
			libWorker.DispatchMessage(o.gs, "MAX_OUTPUTS_REACHED", "MESSAGE")
			break
		}

		for tag, value := range sample.Tags.CloneTags() {
			metricName = fmt.Sprintf("%s::%s::%s", sample.Metric.Name, tag, value)
		}

		if _, ok := flushMetrics[metricName]; !ok {
			metric := *sample.Metric
			metric.Name = metricName
			flushMetrics[metricName] = metric
		}

		// Add sample to sink
		flushMetrics[metricName].Sink.Add(sample)

		outputCount++
	}

	timeSinceStart := time.Since(o.startTime)

	interval := aggregator.Interval{
		Period: int32(o.flushCount),
		Sinks:  make(map[string]*aggregator.Sink, len(flushMetrics)),
	}

	for metricName, metric := range flushMetrics {
		interval.Sinks[metricName] = &aggregator.Sink{
			Labels: metric.Sink.Format(timeSinceStart),
		}
	}

	streamedData := &aggregator.StreamedData{
		DataPoints: []*aggregator.DataPoint{
			{
				Data: &aggregator.DataPoint_Interval{
					Interval: &interval,
				},
			},
		},
	}

	// Encode interval as protobuf
	encodedBytes, err := proto.Marshal(streamedData)
	if err != nil {
		libWorker.HandleError(o.gs, err)
		return
	}

	libWorker.DispatchMessage(o.gs, base64.StdEncoding.EncodeToString(encodedBytes), "INTERVAL")
}
