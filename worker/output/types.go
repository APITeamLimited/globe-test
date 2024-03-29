// Package output contains the interfaces that k6 outputs (and output
// extensions) have to implement, as well as some helpers to make their
// implementation and management easier.
package output

import (
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/metrics"
)

// TODO: make v2 with buffered channels?

// An Output abstracts the process of funneling samples to an external storage
// backend, such as a file or something like an InfluxDB instance.
//
// N.B: All outputs should have non-blocking AddMetricSamples() methods and
// should spawn their own goroutine to flush metrics asynchronously.
type Output interface {
	// Start is called before the Engine tries to use the output and should be
	// used for any long initialization tasks, as well as for starting a
	// goroutine to asynchronously flush metrics to the output.
	Start() error

	// A method to receive the latest metric samples from the Engine. This
	// method is never called concurrently, so do not do anything blocking here
	// that might take a long time. Preferably, just use the SampleBuffer or
	// something like it to buffer metrics until they are flushed.
	AddMetricSamples(samples []metrics.SampleContainer)

	// Flush all remaining metrics and finalize the test run.
	Stop() error
}

// WithThresholds is an output that requires the Engine to give it the
// thresholds before it can be started.
type WithThresholds interface {
	Output
	SetThresholds(map[string]metrics.Thresholds)
}

// WithTestRunStop is an output that can stop the Engine mid-test, interrupting
// the whole test run execution if some internal condition occurs, completely
// independently from the thresholds. It requires a callback function which
// expects an error and triggers the Engine to stop.
type WithTestRunStop interface {
	Output
	SetTestRunStopCallback(func(error))
}

// WithRunStatusUpdates means the output can receive test run status updates.
type WithRunStatusUpdates interface {
	Output
	SetRunStatus(latestStatus libWorker.RunStatus)
}

// WithBuiltinMetrics means the output can receive the builtin metrics.
type WithBuiltinMetrics interface {
	Output
	SetBuiltinMetrics(builtinMetrics *metrics.BuiltinMetrics)
}
