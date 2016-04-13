package aggregate

import (
	"time"
)

type DurationStat struct ***REMOVED***
	Min, Max, Avg, Med time.Duration

	// TODO: Implement a rolling average/median algorithm instead.
	values []time.Duration
***REMOVED***

func (s *DurationStat) Ingest(d time.Duration) ***REMOVED***
	if d < s.Min || s.Min == time.Duration(0) ***REMOVED***
		s.Min = d
	***REMOVED***
	if d > s.Max ***REMOVED***
		s.Max = d
	***REMOVED***
	s.values = append(s.values, d)
***REMOVED***

func (s *DurationStat) End() ***REMOVED***
	sum := time.Duration(0)
	for _, d := range s.values ***REMOVED***
		sum += d
	***REMOVED***
	s.Avg = sum / time.Duration(len(s.values))
	s.Med = s.values[len(s.values)/2]
***REMOVED***
