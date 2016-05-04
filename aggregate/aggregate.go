package aggregate

import (
	"github.com/loadimpact/speedboat/runner"
)

func Aggregate(stats *Stats, in <-chan runner.Result) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		defer stats.End()
		for res := range in ***REMOVED***
			if res.Abort ***REMOVED***
				continue
			***REMOVED***
			stats.Ingest(&res)
			ch <- res
		***REMOVED***
	***REMOVED***()

	return ch
***REMOVED***
