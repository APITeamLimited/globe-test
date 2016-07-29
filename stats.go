package main

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/loadimpact/speedboat/stats/accumulate"
	"github.com/loadimpact/speedboat/stats/writer"
	"io"
)

type Summarizer struct ***REMOVED***
	Accumulator *accumulate.Backend
	Formatter   writer.Formatter
***REMOVED***

func (s *Summarizer) Codify() map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	data := make(map[string]interface***REMOVED******REMOVED***)

	for stat, dimensions := range s.Accumulator.Data ***REMOVED***
		statData := make(map[string]interface***REMOVED******REMOVED***)

		switch stat.Type ***REMOVED***
		case stats.CounterType:
			for dname, dim := range dimensions ***REMOVED***
				val := stats.ApplyIntent(dim.Sum(), stat.Intent)
				if len(dimensions) == 1 ***REMOVED***
					data[stat.Name] = val
				***REMOVED*** else ***REMOVED***
					statData[dname] = val
				***REMOVED***
			***REMOVED***
		case stats.GaugeType:
			for dname, dim := range dimensions ***REMOVED***
				if dim.Last == 0 ***REMOVED***
					continue
				***REMOVED***

				val := stats.ApplyIntent(dim.Last, stat.Intent)
				if len(dimensions) == 1 ***REMOVED***
					data[stat.Name] = val
				***REMOVED*** else ***REMOVED***
					statData[dname] = val
				***REMOVED***
			***REMOVED***
		case stats.HistogramType:
			count := 0
			for dname, dim := range dimensions ***REMOVED***
				l := len(dim.Values)
				if l > count ***REMOVED***
					count = l
				***REMOVED***

				statData[dname] = map[string]interface***REMOVED******REMOVED******REMOVED***
					"min": stats.ApplyIntent(dim.Min(), stat.Intent),
					"max": stats.ApplyIntent(dim.Max(), stat.Intent),
					"avg": stats.ApplyIntent(dim.Avg(), stat.Intent),
					"med": stats.ApplyIntent(dim.Med(), stat.Intent),
					"p90": stats.ApplyIntent(dim.Pct(0.90), stat.Intent),
					"p95": stats.ApplyIntent(dim.Pct(0.95), stat.Intent),
					"p99": stats.ApplyIntent(dim.Pct(0.99), stat.Intent),
				***REMOVED***
			***REMOVED***

			statData["count"] = count
		***REMOVED***

		if len(statData) > 0 ***REMOVED***
			data[stat.Name] = statData
		***REMOVED***
	***REMOVED***

	return data
***REMOVED***

func (s *Summarizer) Format() ([]byte, error) ***REMOVED***
	return s.Formatter.Format(s.Codify())
***REMOVED***

func (s *Summarizer) Print(w io.Writer) error ***REMOVED***
	data, err := s.Format()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err := w.Write(data); err != nil ***REMOVED***
		return err
	***REMOVED***
	if data[len(data)-1] != '\n' ***REMOVED***
		if _, err := w.Write([]byte***REMOVED***'\n'***REMOVED***); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
