package sampler

import (
	"sync"
	"time"
)

const (
	DefaultIntent = iota
	TimeIntent
)

const (
	StatsType = iota
	CounterType
)

type Fields map[string]interface***REMOVED******REMOVED***

type Entry struct ***REMOVED***
	Metric *Metric
	Time   time.Time
	Fields map[string]interface***REMOVED******REMOVED***
	Value  int64
***REMOVED***

func (e *Entry) WithField(key string, value interface***REMOVED******REMOVED***) *Entry ***REMOVED***
	e.Fields[key] = value
	return e
***REMOVED***

func (e *Entry) WithFields(fields Fields) *Entry ***REMOVED***
	for key, value := range fields ***REMOVED***
		e.Fields[key] = value
	***REMOVED***
	return e
***REMOVED***

func (e *Entry) Int(v int) ***REMOVED***
	e.Value = int64(v)
	e.Metric.Write(e)
***REMOVED***

func (e *Entry) Int64(v int64) ***REMOVED***
	e.Value = v
	e.Metric.Write(e)
***REMOVED***

func (e *Entry) Duration(d time.Duration) ***REMOVED***
	e.Value = d.Nanoseconds()
	e.Metric.Intent = TimeIntent
	e.Metric.Write(e)
***REMOVED***

type Metric struct ***REMOVED***
	Name    string
	Sampler *Sampler

	Type   int
	Intent int

	values     []int64
	valueMutex sync.Mutex
***REMOVED***

func (m *Metric) Entry() *Entry ***REMOVED***
	return &Entry***REMOVED***Metric: m, Fields: make(map[string]interface***REMOVED******REMOVED***)***REMOVED***
***REMOVED***

func (m *Metric) WithField(key string, value interface***REMOVED******REMOVED***) *Entry ***REMOVED***
	return m.Entry().WithField(key, value)
***REMOVED***

func (m *Metric) WithFields(fields Fields) *Entry ***REMOVED***
	return m.Entry().WithFields(fields)
***REMOVED***

func (m *Metric) Int(v int) ***REMOVED***
	m.Entry().Int(v)
***REMOVED***

func (m *Metric) Int64(v int64) ***REMOVED***
	m.Entry().Int64(v)
***REMOVED***

func (m *Metric) Duration(d time.Duration) ***REMOVED***
	m.Entry().Duration(d)
***REMOVED***

func (m *Metric) Write(e *Entry) ***REMOVED***
	m.valueMutex.Lock()
	defer m.valueMutex.Unlock()

	m.values = append(m.values, e.Value)
	m.Sampler.Write(m, e)
***REMOVED***

func (m *Metric) Min() int64 ***REMOVED***
	var min int64
	for _, v := range m.values ***REMOVED***
		if min == 0 || v < min ***REMOVED***
			min = v
		***REMOVED***
	***REMOVED***
	return min
***REMOVED***

func (m *Metric) Max() int64 ***REMOVED***
	var max int64
	for _, v := range m.values ***REMOVED***
		if v > max ***REMOVED***
			max = v
		***REMOVED***
	***REMOVED***
	return max
***REMOVED***

func (m *Metric) Avg() int64 ***REMOVED***
	if len(m.values) == 0 ***REMOVED***
		return 0
	***REMOVED***

	var sum int64
	for _, v := range m.values ***REMOVED***
		sum += v
	***REMOVED***
	return sum / int64(len(m.values))
***REMOVED***

func (m *Metric) Med() int64 ***REMOVED***
	idx := len(m.values) / 2
	if idx >= len(m.values) ***REMOVED***
		idx = len(m.values) - 1
	***REMOVED***
	return m.values[idx]
***REMOVED***

func (m *Metric) Last() int64 ***REMOVED***
	return m.values[len(m.values)-1]
***REMOVED***

type Sampler struct ***REMOVED***
	Metrics map[string]*Metric
	Outputs []Output
	OnError func(error)

	MetricMutex sync.Mutex
***REMOVED***

func New() *Sampler ***REMOVED***
	return &Sampler***REMOVED***Metrics: make(map[string]*Metric)***REMOVED***
***REMOVED***

func (s *Sampler) Get(name string) *Metric ***REMOVED***
	s.MetricMutex.Lock()
	defer s.MetricMutex.Unlock()

	metric, ok := s.Metrics[name]
	if !ok ***REMOVED***
		metric = &Metric***REMOVED***Name: name, Sampler: s***REMOVED***
		s.Metrics[name] = metric
	***REMOVED***
	return metric
***REMOVED***

func (s *Sampler) GetAs(name string, t int) *Metric ***REMOVED***
	m := s.Get(name)
	m.Type = t
	return m
***REMOVED***

func (s *Sampler) Stats(name string) *Metric ***REMOVED***
	return s.GetAs(name, StatsType)
***REMOVED***

func (s *Sampler) Counter(name string) *Metric ***REMOVED***
	return s.GetAs(name, CounterType)
***REMOVED***

func (s *Sampler) Write(m *Metric, e *Entry) ***REMOVED***
	for _, out := range s.Outputs ***REMOVED***
		if err := out.Write(m, e); err != nil ***REMOVED***
			if s.OnError != nil ***REMOVED***
				s.OnError(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *Sampler) Commit() error ***REMOVED***
	for _, out := range s.Outputs ***REMOVED***
		if err := out.Commit(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type Output interface ***REMOVED***
	Write(m *Metric, e *Entry) error
	Commit() error
***REMOVED***
