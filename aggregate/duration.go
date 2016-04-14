package aggregate

import (
	"time"
)

type DurationStat struct ***REMOVED***
	Min, Max, Avg, Med time.Duration

	// TODO: Implement a rolling average/median algorithm instead.
	Values []time.Duration
***REMOVED***

func (s *DurationStat) Ingest(d time.Duration) ***REMOVED***
	if d < s.Min || s.Min == time.Duration(0) ***REMOVED***
		s.Min = d
	***REMOVED***
	if d > s.Max ***REMOVED***
		s.Max = d
	***REMOVED***
	s.Values = append(s.Values, d)
***REMOVED***

func (s *DurationStat) End() ***REMOVED***
	sum := time.Duration(0)
	for _, d := range s.Values ***REMOVED***
		sum += d
	***REMOVED***
	s.Avg = sum / time.Duration(len(s.Values))
	s.Med = s.Values[len(s.Values)/2]
***REMOVED***
