package metrics

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// DuplicateMetric is the error returned by Registry.Register when a metric
// already exists.  If you mean to Register that metric you must first
// Unregister the existing metric.
type DuplicateMetric string

func (err DuplicateMetric) Error() string ***REMOVED***
	return fmt.Sprintf("duplicate metric: %s", string(err))
***REMOVED***

// A Registry holds references to a set of metrics by name and can iterate
// over them, calling callback functions provided by the user.
//
// This is an interface so as to encourage other structs to implement
// the Registry API as appropriate.
type Registry interface ***REMOVED***

	// Call the given function for each registered metric.
	Each(func(string, interface***REMOVED******REMOVED***))

	// Get the metric by the given name or nil if none is registered.
	Get(string) interface***REMOVED******REMOVED***

	// GetAll metrics in the Registry.
	GetAll() map[string]map[string]interface***REMOVED******REMOVED***

	// Gets an existing metric or registers the given one.
	// The interface can be the metric to register if not found in registry,
	// or a function returning the metric for lazy instantiation.
	GetOrRegister(string, interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED***

	// Register the given metric under the given name.
	Register(string, interface***REMOVED******REMOVED***) error

	// Run all registered healthchecks.
	RunHealthchecks()

	// Unregister the metric with the given name.
	Unregister(string)

	// Unregister all metrics.  (Mostly for testing.)
	UnregisterAll()
***REMOVED***

// The standard implementation of a Registry is a mutex-protected map
// of names to metrics.
type StandardRegistry struct ***REMOVED***
	metrics map[string]interface***REMOVED******REMOVED***
	mutex   sync.RWMutex
***REMOVED***

// Create a new registry.
func NewRegistry() Registry ***REMOVED***
	return &StandardRegistry***REMOVED***metrics: make(map[string]interface***REMOVED******REMOVED***)***REMOVED***
***REMOVED***

// Call the given function for each registered metric.
func (r *StandardRegistry) Each(f func(string, interface***REMOVED******REMOVED***)) ***REMOVED***
	for name, i := range r.registered() ***REMOVED***
		f(name, i)
	***REMOVED***
***REMOVED***

// Get the metric by the given name or nil if none is registered.
func (r *StandardRegistry) Get(name string) interface***REMOVED******REMOVED*** ***REMOVED***
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.metrics[name]
***REMOVED***

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (r *StandardRegistry) GetOrRegister(name string, i interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	// access the read lock first which should be re-entrant
	r.mutex.RLock()
	metric, ok := r.metrics[name]
	r.mutex.RUnlock()
	if ok ***REMOVED***
		return metric
	***REMOVED***

	// only take the write lock if we'll be modifying the metrics map
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if metric, ok := r.metrics[name]; ok ***REMOVED***
		return metric
	***REMOVED***
	if v := reflect.ValueOf(i); v.Kind() == reflect.Func ***REMOVED***
		i = v.Call(nil)[0].Interface()
	***REMOVED***
	r.register(name, i)
	return i
***REMOVED***

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func (r *StandardRegistry) Register(name string, i interface***REMOVED******REMOVED***) error ***REMOVED***
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.register(name, i)
***REMOVED***

// Run all registered healthchecks.
func (r *StandardRegistry) RunHealthchecks() ***REMOVED***
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, i := range r.metrics ***REMOVED***
		if h, ok := i.(Healthcheck); ok ***REMOVED***
			h.Check()
		***REMOVED***
	***REMOVED***
***REMOVED***

// GetAll metrics in the Registry
func (r *StandardRegistry) GetAll() map[string]map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	data := make(map[string]map[string]interface***REMOVED******REMOVED***)
	r.Each(func(name string, i interface***REMOVED******REMOVED***) ***REMOVED***
		values := make(map[string]interface***REMOVED******REMOVED***)
		switch metric := i.(type) ***REMOVED***
		case Counter:
			values["count"] = metric.Count()
		case Gauge:
			values["value"] = metric.Value()
		case GaugeFloat64:
			values["value"] = metric.Value()
		case Healthcheck:
			values["error"] = nil
			metric.Check()
			if err := metric.Error(); nil != err ***REMOVED***
				values["error"] = metric.Error().Error()
			***REMOVED***
		case Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64***REMOVED***0.5, 0.75, 0.95, 0.99, 0.999***REMOVED***)
			values["count"] = h.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
		case Meter:
			m := metric.Snapshot()
			values["count"] = m.Count()
			values["1m.rate"] = m.Rate1()
			values["5m.rate"] = m.Rate5()
			values["15m.rate"] = m.Rate15()
			values["mean.rate"] = m.RateMean()
		case Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64***REMOVED***0.5, 0.75, 0.95, 0.99, 0.999***REMOVED***)
			values["count"] = t.Count()
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["stddev"] = t.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
			values["1m.rate"] = t.Rate1()
			values["5m.rate"] = t.Rate5()
			values["15m.rate"] = t.Rate15()
			values["mean.rate"] = t.RateMean()
		***REMOVED***
		data[name] = values
	***REMOVED***)
	return data
***REMOVED***

// Unregister the metric with the given name.
func (r *StandardRegistry) Unregister(name string) ***REMOVED***
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.stop(name)
	delete(r.metrics, name)
***REMOVED***

// Unregister all metrics.  (Mostly for testing.)
func (r *StandardRegistry) UnregisterAll() ***REMOVED***
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for name, _ := range r.metrics ***REMOVED***
		r.stop(name)
		delete(r.metrics, name)
	***REMOVED***
***REMOVED***

func (r *StandardRegistry) register(name string, i interface***REMOVED******REMOVED***) error ***REMOVED***
	if _, ok := r.metrics[name]; ok ***REMOVED***
		return DuplicateMetric(name)
	***REMOVED***
	switch i.(type) ***REMOVED***
	case Counter, Gauge, GaugeFloat64, Healthcheck, Histogram, Meter, Timer:
		r.metrics[name] = i
	***REMOVED***
	return nil
***REMOVED***

func (r *StandardRegistry) registered() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	r.mutex.Lock()
	defer r.mutex.Unlock()
	metrics := make(map[string]interface***REMOVED******REMOVED***, len(r.metrics))
	for name, i := range r.metrics ***REMOVED***
		metrics[name] = i
	***REMOVED***
	return metrics
***REMOVED***

func (r *StandardRegistry) stop(name string) ***REMOVED***
	if i, ok := r.metrics[name]; ok ***REMOVED***
		if s, ok := i.(Stoppable); ok ***REMOVED***
			s.Stop()
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stoppable defines the metrics which has to be stopped.
type Stoppable interface ***REMOVED***
	Stop()
***REMOVED***

type PrefixedRegistry struct ***REMOVED***
	underlying Registry
	prefix     string
***REMOVED***

func NewPrefixedRegistry(prefix string) Registry ***REMOVED***
	return &PrefixedRegistry***REMOVED***
		underlying: NewRegistry(),
		prefix:     prefix,
	***REMOVED***
***REMOVED***

func NewPrefixedChildRegistry(parent Registry, prefix string) Registry ***REMOVED***
	return &PrefixedRegistry***REMOVED***
		underlying: parent,
		prefix:     prefix,
	***REMOVED***
***REMOVED***

// Call the given function for each registered metric.
func (r *PrefixedRegistry) Each(fn func(string, interface***REMOVED******REMOVED***)) ***REMOVED***
	wrappedFn := func(prefix string) func(string, interface***REMOVED******REMOVED***) ***REMOVED***
		return func(name string, iface interface***REMOVED******REMOVED***) ***REMOVED***
			if strings.HasPrefix(name, prefix) ***REMOVED***
				fn(name, iface)
			***REMOVED*** else ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	baseRegistry, prefix := findPrefix(r, "")
	baseRegistry.Each(wrappedFn(prefix))
***REMOVED***

func findPrefix(registry Registry, prefix string) (Registry, string) ***REMOVED***
	switch r := registry.(type) ***REMOVED***
	case *PrefixedRegistry:
		return findPrefix(r.underlying, r.prefix+prefix)
	case *StandardRegistry:
		return r, prefix
	***REMOVED***
	return nil, ""
***REMOVED***

// Get the metric by the given name or nil if none is registered.
func (r *PrefixedRegistry) Get(name string) interface***REMOVED******REMOVED*** ***REMOVED***
	realName := r.prefix + name
	return r.underlying.Get(realName)
***REMOVED***

// Gets an existing metric or registers the given one.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (r *PrefixedRegistry) GetOrRegister(name string, metric interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	realName := r.prefix + name
	return r.underlying.GetOrRegister(realName, metric)
***REMOVED***

// Register the given metric under the given name. The name will be prefixed.
func (r *PrefixedRegistry) Register(name string, metric interface***REMOVED******REMOVED***) error ***REMOVED***
	realName := r.prefix + name
	return r.underlying.Register(realName, metric)
***REMOVED***

// Run all registered healthchecks.
func (r *PrefixedRegistry) RunHealthchecks() ***REMOVED***
	r.underlying.RunHealthchecks()
***REMOVED***

// GetAll metrics in the Registry
func (r *PrefixedRegistry) GetAll() map[string]map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return r.underlying.GetAll()
***REMOVED***

// Unregister the metric with the given name. The name will be prefixed.
func (r *PrefixedRegistry) Unregister(name string) ***REMOVED***
	realName := r.prefix + name
	r.underlying.Unregister(realName)
***REMOVED***

// Unregister all metrics.  (Mostly for testing.)
func (r *PrefixedRegistry) UnregisterAll() ***REMOVED***
	r.underlying.UnregisterAll()
***REMOVED***

var DefaultRegistry Registry = NewRegistry()

// Call the given function for each registered metric.
func Each(f func(string, interface***REMOVED******REMOVED***)) ***REMOVED***
	DefaultRegistry.Each(f)
***REMOVED***

// Get the metric by the given name or nil if none is registered.
func Get(name string) interface***REMOVED******REMOVED*** ***REMOVED***
	return DefaultRegistry.Get(name)
***REMOVED***

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
func GetOrRegister(name string, i interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return DefaultRegistry.GetOrRegister(name, i)
***REMOVED***

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func Register(name string, i interface***REMOVED******REMOVED***) error ***REMOVED***
	return DefaultRegistry.Register(name, i)
***REMOVED***

// Register the given metric under the given name.  Panics if a metric by the
// given name is already registered.
func MustRegister(name string, i interface***REMOVED******REMOVED***) ***REMOVED***
	if err := Register(name, i); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// Run all registered healthchecks.
func RunHealthchecks() ***REMOVED***
	DefaultRegistry.RunHealthchecks()
***REMOVED***

// Unregister the metric with the given name.
func Unregister(name string) ***REMOVED***
	DefaultRegistry.Unregister(name)
***REMOVED***
