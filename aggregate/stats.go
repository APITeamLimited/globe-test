package aggregate

import (
	"github.com/loadimpact/speedboat/runner"
)

type Stats struct ***REMOVED***
	Results int64
	Time    DurationStat
***REMOVED***

func (s *Stats) Ingest(res *runner.Result) ***REMOVED***
	s.Results++
	s.Time.Ingest(res.Time)
***REMOVED***

func (s *Stats) End() ***REMOVED***
	s.Time.End()
***REMOVED***
