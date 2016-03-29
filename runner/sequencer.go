package runner

import (
	"time"
)

type Sequencer struct ***REMOVED***
	Metrics []Metric
***REMOVED***

type Stat struct ***REMOVED***
	Min, Max, Avg, Med float64
***REMOVED***

type Stats struct ***REMOVED***
	Duration Stat
***REMOVED***

func NewSequencer() Sequencer ***REMOVED***
	return Sequencer***REMOVED******REMOVED***
***REMOVED***

func (s *Sequencer) Add(m Metric) ***REMOVED***
	s.Metrics = append(s.Metrics, m)
***REMOVED***

func (s *Sequencer) Count() int ***REMOVED***
	return len(s.Metrics)
***REMOVED***

func (s *Sequencer) StatDuration() (st Stat) ***REMOVED***
	count := s.Count()
	if count == 0 ***REMOVED***
		return st
	***REMOVED***

	total := time.Duration(0)
	min := time.Duration(0)
	max := time.Duration(0)
	for i := 0; i < count; i++ ***REMOVED***
		m := s.Metrics[i]
		total += m.Duration
		if m.Duration < min || min == time.Duration(0) ***REMOVED***
			min = m.Duration
		***REMOVED***
		if m.Duration > max ***REMOVED***
			max = m.Duration
		***REMOVED***
	***REMOVED***

	avg := time.Duration(total.Nanoseconds() / int64(count))
	med := s.Metrics[len(s.Metrics)/2].Duration

	return Stat***REMOVED***
		Min: min.Seconds(),
		Max: max.Seconds(),
		Avg: avg.Seconds(),
		Med: med.Seconds(),
	***REMOVED***
***REMOVED***

func (s *Sequencer) Stats() (st Stats) ***REMOVED***
	st.Duration = s.StatDuration()
	return st
***REMOVED***
