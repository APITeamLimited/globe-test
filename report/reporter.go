package report

import (
	"github.com/loadimpact/speedboat/runner"
	"io"
)

type Reporter interface ***REMOVED***
	Begin(w io.Writer)
	Report(w io.Writer, res *runner.Result)
	End(w io.Writer)
***REMOVED***

func Report(r Reporter, w io.Writer, in <-chan runner.Result) <-chan runner.Result ***REMOVED***
	ch := make(chan runner.Result)

	go func() ***REMOVED***
		defer close(ch)

		r.Begin(w)
		for res := range in ***REMOVED***
			r.Report(w, &res)
			ch <- res
		***REMOVED***
		r.End(w)
	***REMOVED***()

	return ch
***REMOVED***
