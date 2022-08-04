package outputtest

import (
	"strconv"

	"github.com/spf13/afero"
	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

func init() ***REMOVED***
	output.RegisterExtension("outputtest", func(params output.Params) (output.Output, error) ***REMOVED***
		return &Output***REMOVED***params: params***REMOVED***, nil
	***REMOVED***)
***REMOVED***

// Output is meant to test xk6 and the output extension sub-system of k6.
type Output struct ***REMOVED***
	params     output.Params
	calcResult float64
	outputFile afero.File
***REMOVED***

// Description returns a human-readable description of the output.
func (o *Output) Description() string ***REMOVED***
	return "test output extension"
***REMOVED***

// Start opens the specified output file.
func (o *Output) Start() error ***REMOVED***
	out, err := o.params.FS.Create(o.params.ConfigArgument)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.outputFile = out

	return nil
***REMOVED***

// AddMetricSamples just plucks out the metric we're interested in.
func (o *Output) AddMetricSamples(sampleContainers []metrics.SampleContainer) ***REMOVED***
	for _, sc := range sampleContainers ***REMOVED***
		for _, sample := range sc.GetSamples() ***REMOVED***
			if sample.Metric.Name == "foos" ***REMOVED***
				o.calcResult += sample.Value
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Stop saves the dummy results and closes the file.
func (o *Output) Stop() error ***REMOVED***
	_, err := o.outputFile.Write([]byte(strconv.FormatFloat(o.calcResult, 'f', 0, 64)))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return o.outputFile.Close()
***REMOVED***
