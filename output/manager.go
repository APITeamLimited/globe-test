package output

import (
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/stats"
)

// Manager can be used to manage multiple outputs at the same time.
type Manager struct ***REMOVED***
	outputs []Output
	logger  logrus.FieldLogger

	testStopCallback func(error)
***REMOVED***

// NewManager returns a new manager for the given outputs.
func NewManager(outputs []Output, logger logrus.FieldLogger, testStopCallback func(error)) *Manager ***REMOVED***
	return &Manager***REMOVED***
		outputs:          outputs,
		logger:           logger.WithField("component", "output-manager"),
		testStopCallback: testStopCallback,
	***REMOVED***
***REMOVED***

// StartOutputs spins up all configured outputs. If some output fails to start,
// it stops the already started ones. This may take some time, since some
// outputs make initial network requests to set up whatever remote services are
// going to listen to them.
func (om *Manager) StartOutputs() error ***REMOVED***
	om.logger.Debugf("Starting %d outputs...", len(om.outputs))
	for i, out := range om.outputs ***REMOVED***
		if stopOut, ok := out.(WithTestRunStop); ok ***REMOVED***
			stopOut.SetTestRunStopCallback(om.testStopCallback)
		***REMOVED***

		if err := out.Start(); err != nil ***REMOVED***
			om.stopOutputs(i)
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// StopOutputs stops all configured outputs.
func (om *Manager) StopOutputs() ***REMOVED***
	om.stopOutputs(len(om.outputs))
***REMOVED***

func (om *Manager) stopOutputs(upToID int) ***REMOVED***
	om.logger.Debugf("Stopping %d outputs...", upToID)
	for i := 0; i < upToID; i++ ***REMOVED***
		if err := om.outputs[i].Stop(); err != nil ***REMOVED***
			om.logger.WithError(err).Errorf("Stopping output %d failed", i)
		***REMOVED***
	***REMOVED***
***REMOVED***

// SetRunStatus checks which outputs implement the WithRunStatusUpdates
// interface and sets the provided RunStatus to them.
func (om *Manager) SetRunStatus(status lib.RunStatus) ***REMOVED***
	for _, out := range om.outputs ***REMOVED***
		if statUpdOut, ok := out.(WithRunStatusUpdates); ok ***REMOVED***
			statUpdOut.SetRunStatus(status)
		***REMOVED***
	***REMOVED***
***REMOVED***

// AddMetricSamples is a temporary method to make the Manager usable in the
// current Engine. It needs to be replaced with the full metric pump.
//
// TODO: refactor
func (om *Manager) AddMetricSamples(sampleContainers []stats.SampleContainer) ***REMOVED***
	if len(sampleContainers) == 0 ***REMOVED***
		return
	***REMOVED***

	for _, out := range om.outputs ***REMOVED***
		out.AddMetricSamples(sampleContainers)
	***REMOVED***
***REMOVED***
