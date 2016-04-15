package report

import (
	"fmt"
	"github.com/loadimpact/speedboat/runner"
	"io"
	"time"
)

type CSVReporter struct***REMOVED******REMOVED***

func (CSVReporter) Begin(w io.Writer) ***REMOVED******REMOVED***

func (CSVReporter) Report(w io.Writer, res *runner.Result) ***REMOVED***
	// TODO: Timestamp events themselves!
	t := time.Now()
	errString := ""
	if res.Error != nil ***REMOVED***
		errString = res.Error.Error()
	***REMOVED***
	fmt.Fprintf(w, "%d;%d;%s;%s\n", t.Unix(), res.Time.Nanoseconds(), res.Text, errString)
***REMOVED***

func (CSVReporter) End(w io.Writer) ***REMOVED******REMOVED***
