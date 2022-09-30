package minirunner

import (
	"context"
	"io"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

// Ensure mock implementations conform to the interfaces.
var (
	_ libWorker.Runner        = &MiniRunner{}
	_ libWorker.InitializedVU = &VU{}
	_ libWorker.ActiveVU      = &ActiveVU{}
)

// MiniRunner partially implements the libWorker.Runner interface, but instead of
// using a real JS runtime, it allows us to directly specify the options and
// functions with Go code.
type MiniRunner struct {
	Fn              func(ctx context.Context, state *libWorker.State, out chan<- workerMetrics.SampleContainer) error
	SetupFn         func(ctx context.Context, out chan<- workerMetrics.SampleContainer) ([]byte, error)
	TeardownFn      func(ctx context.Context, out chan<- workerMetrics.SampleContainer) error
	HandleSummaryFn func(context.Context, *libWorker.Summary) (map[string]io.Reader, error)

	SetupData []byte

	Group   *libWorker.Group
	Options libWorker.Options
}

// MakeArchive isn't implemented, it always returns nil and is just here to
// satisfy the libWorker.Runner interface.
func (r MiniRunner) MakeArchive() *libWorker.Archive {
	return nil
}

// NewVU returns a new VU with an incremental ID.
func (r *MiniRunner) NewVU(idLocal, idGlobal uint64, out chan<- workerMetrics.SampleContainer, workerInfo *libWorker.WorkerInfo) (libWorker.InitializedVU, error) {
	state := &libWorker.State{VUID: idLocal, VUIDGlobal: idGlobal, Iteration: int64(-1)}
	return &VU{
		R:            r,
		Out:          out,
		ID:           idLocal,
		IDGlobal:     idGlobal,
		state:        state,
		scenarioIter: make(map[string]uint64),
	}, nil
}

// Setup calls the supplied mock setup() function, if present.
func (r *MiniRunner) Setup(ctx context.Context, out chan<- workerMetrics.SampleContainer) (err error) {
	if fn := r.SetupFn; fn != nil {
		r.SetupData, err = fn(ctx, out)
	}
	return
}

// GetSetupData returns json representation of the setup data if setup() is
// specified and was ran, nil otherwise.
func (r MiniRunner) GetSetupData() []byte {
	return r.SetupData
}

// SetSetupData saves the externally supplied setup data as JSON in the runner.
func (r *MiniRunner) SetSetupData(data []byte) {
	r.SetupData = data
}

// Teardown calls the supplied mock teardown() function, if present.
func (r MiniRunner) Teardown(ctx context.Context, out chan<- workerMetrics.SampleContainer) error {
	if fn := r.TeardownFn; fn != nil {
		return fn(ctx, out)
	}
	return nil
}

// GetDefaultGroup returns the default group.
func (r MiniRunner) GetDefaultGroup() *libWorker.Group {
	if r.Group == nil {
		r.Group = &libWorker.Group{}
	}
	return r.Group
}

// IsExecutable satisfies libWorker.Runner, but is mocked for MiniRunner since
// it doesn't deal with JS.
func (r MiniRunner) IsExecutable(name string) bool {
	return true
}

// GetOptions returns the supplied options struct.
func (r MiniRunner) GetOptions() libWorker.Options {
	return r.Options
}

// SetOptions allows you to override the runner options.
func (r *MiniRunner) SetOptions(opts libWorker.Options) error {
	r.Options = opts
	return nil
}

// HandleSummary calls the specified summary callback, if supplied.
func (r *MiniRunner) HandleSummary(ctx context.Context, s *libWorker.Summary) (map[string]io.Reader, error) {
	if r.HandleSummaryFn != nil {
		return r.HandleSummaryFn(ctx, s)
	}
	return nil, nil
}

// VU is a mock VU, spawned by a MiniRunner.
type VU struct {
	R            *MiniRunner
	Out          chan<- workerMetrics.SampleContainer
	ID, IDGlobal uint64
	Iteration    int64
	state        *libWorker.State
	// count of iterations executed by this VU in each scenario
	scenarioIter map[string]uint64
}

// ActiveVU holds a VU and its activation parameters
type ActiveVU struct {
	*VU
	*libWorker.VUActivationParams
	busy chan struct{}

	scenarioName              string
	getNextIterations         func() (uint64, uint64)
	scIterLocal, scIterGlobal uint64
}

// GetID returns the unique VU ID.
func (vu *VU) GetID() uint64 {
	return vu.ID
}

// State returns the VU's State.
func (vu *VU) State() *libWorker.State {
	return vu.state
}

// Activate the VU so it will be able to run code.
func (vu *VU) Activate(params *libWorker.VUActivationParams) libWorker.ActiveVU {
	ctx := params.RunContext

	vu.state.GetScenarioVUIter = func() uint64 {
		return vu.scenarioIter[params.Scenario]
	}

	avu := &ActiveVU{
		VU:                 vu,
		VUActivationParams: params,
		busy:               make(chan struct{}, 1),
		scenarioName:       params.Scenario,
		scIterLocal:        ^uint64(0),
		scIterGlobal:       ^uint64(0),
		getNextIterations:  params.GetNextIterationCounters,
	}

	vu.state.GetScenarioLocalVUIter = func() uint64 {
		return avu.scIterLocal
	}
	vu.state.GetScenarioGlobalVUIter = func() uint64 {
		return avu.scIterGlobal
	}

	go func() {
		<-ctx.Done()

		// Wait for the VU to stop running, if it was, and prevent it from
		// running again for this activation
		avu.busy <- struct{}{}

		if params.DeactivateCallback != nil {
			params.DeactivateCallback(vu)
		}
	}()

	return avu
}

func (vu *ActiveVU) incrIteration() {
	vu.Iteration++
	vu.state.Iteration = vu.Iteration

	if _, ok := vu.scenarioIter[vu.scenarioName]; ok {
		vu.scenarioIter[vu.scenarioName]++
	} else {
		vu.scenarioIter[vu.scenarioName] = 0
	}
	vu.scIterLocal, vu.scIterGlobal = vu.getNextIterations()
}

// RunOnce runs the mock default function once, incrementing its iteration.
func (vu *ActiveVU) RunOnce() error {
	if vu.R.Fn == nil {
		return nil
	}

	select {
	case <-vu.RunContext.Done():
		return vu.RunContext.Err() // we are done, return
	case vu.busy <- struct{}{}:
		// nothing else can run now, and the VU cannot be deactivated
	}
	defer func() {
		<-vu.busy // unlock deactivation again
	}()

	vu.incrIteration()
	return vu.R.Fn(vu.RunContext, vu.State(), vu.Out)
}
