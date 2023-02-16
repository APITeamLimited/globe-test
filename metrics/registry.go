package metrics

import (
	"encoding/json"
	"fmt"
	"regexp"
	"runtime/metrics"
	"sync"
)

// Registry is what can create metrics
type Registry struct {
	metrics map[string]*Metric
	l       sync.RWMutex
}

// NewRegistry returns a new registry
func NewRegistry() *Registry {
	return &Registry{
		metrics: make(map[string]*Metric),
	}
}

const nameRegexString = "^[\\p{L}\\p{N}\\._ !\\?/&#\\(\\)<>%-]{1,128}$"

var compileNameRegex = regexp.MustCompile(nameRegexString)

func checkName(name string) bool {
	return compileNameRegex.Match([]byte(name))
}

// Adds remote worker samples
func (r *Registry) AddSamples(rawMessage string) error {
	samples := []metrics.Sample{}

	err := json.Unmarshal([]byte(rawMessage), &samples)
	if err != nil {
		return err
	}

	r.l.Lock()
	for _, sample := range samples {

		metrics := sample.GetSamples()

		// Assign the metric if it's not already registered
		if _, ok := r.metrics[sample.Metric.Name]; !ok {
			m := sample.Metric
			r.metrics[sample.Metric.Name] = m
		}
	}
	r.l.Unlock()
}

// NewMetric returns new metric registered to this registry
// TODO have multiple versions returning specific metric types when we have such things
func (r *Registry) NewMetric(name string, typ MetricType, t ...ValueType) (*Metric, error) {
	r.l.Lock()
	defer r.l.Unlock()

	if !checkName(name) {
		return nil, fmt.Errorf("invalid metric name: '%s'", name) //nolint:golint,stylecheck
	}
	oldMetric, ok := r.metrics[name]

	if !ok {
		m := InstantiateMetric(name, typ, t...)
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
func (r *Registry) MustNewMetric(name string, typ MetricType, t ...ValueType) *Metric {
	m, err := r.NewMetric(name, typ, t...)
	if err != nil {
		panic(err)
	}
	return m
}

// Get returns the Metric with the given name. If that metric doesn't exist,
// Get() will return a nil value.
func (r *Registry) Get(name string) *Metric {
	return r.metrics[name]
}
