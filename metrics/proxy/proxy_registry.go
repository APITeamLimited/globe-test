// Moves metrics every second to the orchestrator
package proxy_registry

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/metrics"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"gopkg.in/guregu/null.v3"
)

const proxyFlushInterval = 1 * time.Second

// ProxyRegistry is what can create metrics
type ProxyRegistry struct {
	gs          libWorker.BaseGlobalState
	metrics     map[string]*metrics.Metric
	samples     []metrics.Sample
	l           sync.RWMutex
	samplesChan chan metrics.SampleContainer
	ticker      *time.Ticker
}

// NewProxyRegistry returns a new registry
func NewProxyRegistry(metricSamplesBufferSize null.Int, gs libWorker.BaseGlobalState) *ProxyRegistry {
	return &ProxyRegistry{
		gs:          gs,
		metrics:     make(map[string]*metrics.Metric),
		samplesChan: make(chan metrics.SampleContainer, metricSamplesBufferSize.Int64),
	}
}

func (r *ProxyRegistry) Start() {
	r.ticker = time.NewTicker(proxyFlushInterval)

	go func() {
		for range r.ticker.C {
			r.flush()
		}
	}()

	go func() {
		// On samples from the engine, register them to the metrics
		for sample := range r.samplesChan {
			r.l.Lock()
			metrics := sample.GetSamples()

			r.samples = append(r.samples, metrics...)

			for _, s := range metrics {
				// Assign the metric if it's not already registered
				if _, ok := r.metrics[s.Metric.Name]; !ok {
					m := s.Metric
					r.metrics[s.Metric.Name] = m
				}
			}

			r.l.Unlock()

		}
	}()
}

func (r *ProxyRegistry) Stop() {
	r.ticker.Stop()
}

const nameRegexString = "^[\\p{L}\\p{N}\\._ !\\?/&#\\(\\)<>%-]{1,128}$"

var compileNameRegex = regexp.MustCompile(nameRegexString)

func checkName(name string) bool {
	return compileNameRegex.Match([]byte(name))
}

func (r *ProxyRegistry) GetSamplesChan() chan metrics.SampleContainer {
	return r.samplesChan
}

func (r *ProxyRegistry) flush() {
	r.l.RLock()
	defer r.l.RUnlock()

	marshalledMetrics, err := json.Marshal(r.samples)
	if err != nil {
		panic(err)
	}

	r.samples = make([]metrics.Sample, 0)

	libWorker.DispatchMessage(r.gs, string(marshalledMetrics), "METRICS")

}

// NewMetric returns new metric registered to this registry
// TODO have multiple versions returning specific metric types when we have such things
func (r *ProxyRegistry) NewMetric(name string, typ metrics.MetricType, t ...metrics.ValueType) (*metrics.Metric, error) {
	r.l.Lock()
	defer r.l.Unlock()

	if !checkName(name) {
		return nil, fmt.Errorf("invalid metric name: '%s'", name)
	}
	oldMetric, ok := r.metrics[name]

	if !ok {
		m := metrics.InstantiateMetric(name, typ, t...)
		r.metrics[name] = m
		return m, nil
	}
	if oldMetric.Type != typ {
		return nil, fmt.Errorf("metric '%s' already exists but with type %s, instead of %s", name, oldMetric.Type, typ)
	}
	if len(t) > 0 {
		if t[0] != oldMetric.Contains {
			return nil, fmt.Errorf("metric '%s' already exists but with a value type %s, instead of %s",
				name, oldMetric.Contains, t[0])
		}
	}
	return oldMetric, nil
}

// MustNewMetric is like NewMetric, but will panic if there is an error
func (r *ProxyRegistry) MustNewMetric(name string, typ metrics.MetricType, t ...metrics.ValueType) *metrics.Metric {
	m, err := r.NewMetric(name, typ, t...)
	if err != nil {
		panic(err)
	}
	return m
}

// THis may be required for some built in functionality on workers
func RegisterBuiltinMetrics(registry *ProxyRegistry) *metrics.BuiltinMetrics {
	return &metrics.BuiltinMetrics{
		VUs:               registry.MustNewMetric(metrics.VUsName, metrics.Gauge),
		VUsMax:            registry.MustNewMetric(metrics.VUsMaxName, metrics.Gauge),
		Iterations:        registry.MustNewMetric(metrics.IterationsName, metrics.Counter),
		IterationDuration: registry.MustNewMetric(metrics.IterationDurationName, metrics.Trend, metrics.Time),
		DroppedIterations: registry.MustNewMetric(metrics.DroppedIterationsName, metrics.Counter),

		Checks:        registry.MustNewMetric(metrics.ChecksName, metrics.Rate),
		GroupDuration: registry.MustNewMetric(metrics.GroupDurationName, metrics.Trend, metrics.Time),

		HTTPReqs:              registry.MustNewMetric(metrics.HTTPReqsName, metrics.Counter),
		HTTPReqFailed:         registry.MustNewMetric(metrics.HTTPReqFailedName, metrics.Rate),
		HTTPReqDuration:       registry.MustNewMetric(metrics.HTTPReqDurationName, metrics.Trend, metrics.Time),
		HTTPReqBlocked:        registry.MustNewMetric(metrics.HTTPReqBlockedName, metrics.Trend, metrics.Time),
		HTTPReqConnecting:     registry.MustNewMetric(metrics.HTTPReqConnectingName, metrics.Trend, metrics.Time),
		HTTPReqTLSHandshaking: registry.MustNewMetric(metrics.HTTPReqTLSHandshakingName, metrics.Trend, metrics.Time),
		HTTPReqSending:        registry.MustNewMetric(metrics.HTTPReqSendingName, metrics.Trend, metrics.Time),
		HTTPReqWaiting:        registry.MustNewMetric(metrics.HTTPReqWaitingName, metrics.Trend, metrics.Time),
		HTTPReqReceiving:      registry.MustNewMetric(metrics.HTTPReqReceivingName, metrics.Trend, metrics.Time),

		WSSessions:         registry.MustNewMetric(metrics.WSSessionsName, metrics.Counter),
		WSMessagesSent:     registry.MustNewMetric(metrics.WSMessagesSentName, metrics.Counter),
		WSMessagesReceived: registry.MustNewMetric(metrics.WSMessagesReceivedName, metrics.Counter),
		WSPing:             registry.MustNewMetric(metrics.WSPingName, metrics.Trend, metrics.Time),
		WSSessionDuration:  registry.MustNewMetric(metrics.WSSessionDurationName, metrics.Trend, metrics.Time),
		WSConnecting:       registry.MustNewMetric(metrics.WSConnectingName, metrics.Trend, metrics.Time),

		GRPCReqDuration: registry.MustNewMetric(metrics.GRPCReqDurationName, metrics.Trend, metrics.Time),

		DataSent:     registry.MustNewMetric(metrics.DataSentName, metrics.Counter, metrics.Data),
		DataReceived: registry.MustNewMetric(metrics.DataReceivedName, metrics.Counter, metrics.Data),
	}
}
