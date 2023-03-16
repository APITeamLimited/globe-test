package globetest

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/aggregator"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/APITeamLimited/globe-test/worker/output"
	"google.golang.org/protobuf/proto"
)

const (
	uniqueOutputLimit = 10000
	flushPeriod       = 6000 * time.Millisecond
)

type Output struct {
	output.SampleBuffer

	periodicFlusher *output.PeriodicFlusher

	gs       libWorker.BaseGlobalState
	location string

	seenMetrics map[string]struct{}

	flushCount int
	flushMutex sync.Mutex

	consoleMessageCount int
	addedSinceLastFlush bool

	consoleMessages []*aggregator.ConsoleMessage
	consoleMutex    sync.Mutex
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

		flushMutex: sync.Mutex{},

		consoleMessages: make([]*aggregator.ConsoleMessage, 0),
		consoleMutex:    sync.Mutex{},
	}, nil

}

func (o *Output) Start() error {
	pf, err := output.NewPeriodicFlusher(flushPeriod, func() {
		o.flushMetrics()
		o.flushConsoleMessages()
	})
	if err != nil {
		return err
	}
	o.periodicFlusher = pf

	go o.listenOnLoggerChannel()

	return nil
}

func (o *Output) Stop() error {
	o.periodicFlusher.Stop()
	return nil
}

func (o *Output) flushMetrics() {
	o.flushMutex.Lock()

	defer func() {
		o.flushCount++
		o.flushMutex.Unlock()
	}()

	flushedMetrics := make(map[string]metrics.Metric)

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

		if _, ok := flushedMetrics[metricName]; !ok {
			metric := *sample.Metric
			metric.Name = metricName
			flushedMetrics[metricName] = metric
		}

		// Add sample to sink
		flushedMetrics[metricName].Sink.Add(sample)

		outputCount++
	}

	interval := aggregator.Interval{
		Period: int32(o.flushCount),
		Sinks:  make(map[string]*aggregator.Sink, len(flushedMetrics)),
	}

	for metricName, metric := range flushedMetrics {
		sinkType, err := getSinkType(metric.Sink)
		if err != nil {
			libWorker.HandleError(o.gs, err)
			return
		}

		interval.Sinks[metricName] = &aggregator.Sink{
			Type:   sinkType,
			Labels: metric.Sink.Format(),
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

func getSinkType(sink metrics.Sink) (aggregator.SinkType, error) {
	switch sink.(type) {
	case *metrics.CounterSink:
		return aggregator.SinkType_Counter, nil
	case *metrics.GaugeSink:
		return aggregator.SinkType_Gauge, nil
	case *metrics.TrendSink:
		return aggregator.SinkType_Trend, nil
	case *metrics.RateSink:
		return aggregator.SinkType_Rate, nil
	default:
		// Incorrect sink type being returned here but faster than using a
		// custom struct to return nil and an error
		return aggregator.SinkType_Counter, errors.New("unknown sink type")
	}
}
