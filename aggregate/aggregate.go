package aggregate

import (
	"github.com/loadimpact/speedboat/runner"
)

func Aggregate(stats *Stats, in <-chan runner.Result) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		for res := range in ***REMOVED***
			stats.Ingest(&res)
			ch <- res
		***REMOVED***

		stats.End()
	***REMOVED***()

	return ch
***REMOVED***
