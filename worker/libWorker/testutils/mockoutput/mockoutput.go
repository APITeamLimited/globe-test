package mockoutput

import (
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
)

// New exists so that the usage from tests avoids repetition, i.e. is
// mockoutput.New() instead of &mockoutput.MockOutput{}
func New() *MockOutput {
	return &MockOutput{}
}

// MockOutput can be used in tests to mock an actual output.
type MockOutput struct {
	SampleContainers []workerMetrics.SampleContainer
	Samples          []workerMetrics.Sample
	RunStatus        libWorker.RunStatus

	DescFn  func() string
	StartFn func() error
	StopFn  func() error
}

var _ output.WithRunStatusUpdates = &MockOutput{}

// AddMetricSamples just saves the results in memory.
func (mo *MockOutput) AddMetricSamples(scs []workerMetrics.SampleContainer) {
	mo.SampleContainers = append(mo.SampleContainers, scs...)
	for _, sc := range scs {
		mo.Samples = append(mo.Samples, sc.GetSamples()...)
	}
}

// SetRunStatus updates the RunStatus property.
func (mo *MockOutput) SetRunStatus(latestStatus libWorker.RunStatus) {
	mo.RunStatus = latestStatus
}

// Description calls the supplied DescFn callback, if available.
func (mo *MockOutput) Description() string {
	if mo.DescFn != nil {
		return mo.DescFn()
	}
	return "mock"
}

// Start calls the supplied StartFn callback, if available.
func (mo *MockOutput) Start() error {
	if mo.StartFn != nil {
		return mo.StartFn()
	}
	return nil
}

// Stop calls the supplied StopFn callback, if available.
func (mo *MockOutput) Stop() error {
	if mo.StopFn != nil {
		return mo.StopFn()
	}
	return nil
}
