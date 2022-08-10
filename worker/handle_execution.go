package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis/v9"
	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/ui/pb"
)

/*
This is the main function that is called when the worker is started.
It is responsible for running a job and reporting on its status
*/
func handleExecution(ctx context.Context,
	client *redis.Client, job map[string]string, workerId string) {
	// Check if redis message is a uuid

	fmt.Println("\033[1;32mGot job", job["id"], "\033[0m")

	go updateStatus(ctx, client, job["id"], workerId, "LOADING")

	globalState := newGlobalState(ctx, client, job["id"], workerId)

	test, err := loadAndConfigureTest(globalState, job)

	if err != nil {
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("failed to loadAndConfigureTest: %s", err.Error()))
		return
	}

	go dispatchMessage(ctx, client, job["id"], workerId, fmt.Sprintf("Loaded test %s", test.workerLoadedTest.sourceRootPath), "DEBUG")

	// Write the full consolidated *and derived* options back to the Runner.
	conf := test.derivedConfig
	testRunState, err := test.buildTestRunState(conf.Options)
	if err != nil {
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error building testRunState %s", err.Error()))
		return
	}

	// We prepare a bunch of contexts:
	//  - The runCtx is cancelled as soon as the Engine's run() lambda finishes,
	//    and can trigger things like the usage report and end of test summary.
	//    Crucially, metrics processing by the Engine will still work after this
	//    context is cancelled!
	//  - The lingerCtx is cancelled by Ctrl+C, and is used to wait for that
	//    event when k6 was ran with the --linger option.
	//  - The globalCtx is cancelled only after we're completely done with the
	//    test execution and any --linger has been cleared, so that the Engine
	//    can start winding down its metrics processing.
	globalCtx, globalCancel := context.WithCancel(globalState.ctx)
	defer globalCancel()
	lingerCtx, lingerCancel := context.WithCancel(globalCtx)
	defer lingerCancel()
	runCtx, runCancel := context.WithCancel(lingerCtx)
	defer runCancel()

	logger := testRunState.Logger
	logger.Debug("Initializing the execution scheduler...")
	execScheduler, err := local.NewExecutionScheduler(testRunState)
	if err != nil {
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error initializing the execution scheduler: %s", err.Error()))
		return
	}

	// This is manually triggered after the Engine's Run() has completed,
	// and things like a single Ctrl+C don't affect it. We use it to make
	// sure that the progressbars finish updating with the latest execution
	// state one last time, after the test run has finished.
	_, progressCancel := context.WithCancel(globalCtx)
	defer progressCancel()
	progressBarWG := &sync.WaitGroup{}
	progressBarWG.Add(1)
	go func() {
		pbs := []*pb.ProgressBar{execScheduler.GetInitProgressBar()}
		for _, s := range execScheduler.GetExecutors() {
			pbs = append(pbs, s.GetProgress())
		}
		progressBarWG.Done()
	}()

	// Create all outputs.
	executionPlan := execScheduler.GetExecutionPlan()
	outputs, err := createOutputs(globalState, test, executionPlan)
	if err != nil {
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error creating outputs %s", err.Error()))
		return
	}

	// TODO: create a MetricsEngine here and add its ingester to the list of
	// outputs (unless both NoThresholds and NoSummary were enabled)

	// TODO: remove this completely
	// Create the engine.
	go dispatchMessage(ctx, client, job["id"], workerId, "Initializing the Engine...", "DEBUG")
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil {
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error creating engine %s", err.Error()))
		return
	}

	// We do this here so we can get any output URLs below.
	err = engine.OutputManager.StartOutputs()
	if err != nil {
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error starting outputs %s", err.Error()))
		return
	}
	defer engine.OutputManager.StopOutputs()

	// Trap Interrupts, SIGINTs and SIGTERMs.
	gracefulStop := func(sig os.Signal) {
		logger.WithField("sig", sig).Debug("Stopping k6 in response to signal...")
		lingerCancel() // stop the test run, metric processing is cancelled below
	}
	onHardStop := func(sig os.Signal) {
		logger.WithField("sig", sig).Error("Aborting k6 in response to signal")
		globalCancel() // not that it matters, given the following command...
	}
	stopSignalHandling := handleTestAbortSignals(globalState, gracefulStop, onHardStop)
	defer stopSignalHandling()

	// Initialize the engine
	go dispatchMessage(ctx, client, job["id"], workerId, "Initializing VU(s)...", "DEBUG")
	engineRun, engineWait, err := engine.Init(globalCtx, runCtx)
	if err != nil {
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		go handleError(ctx, client, job["id"], workerId, errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		return
	}

	// Start the test run
	go updateStatus(ctx, client, job["id"], workerId, "RUNNING")
	var interrupt error
	err = engineRun()
	if err != nil {
		err = common.UnwrapGojaInterruptedError(err)
		if errext.IsInterruptError(err) {
			// Don't return here since we need to work with --linger,
			// show the end-of-test summary and exit cleanly.
			interrupt = err
		}
		if !conf.Linger.Bool && interrupt == nil {
			fmt.Println(errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		}
	}
	runCancel()
	go dispatchMessage(ctx, client, job["id"], workerId, "Engine run terminated cleanly", "DEBUG")

	progressCancel()

	executionState := execScheduler.GetState()
	// Warn if no iterations could be completed.
	if executionState.GetFullIterationCount() == 0 {
		go dispatchMessage(ctx, client, job["id"], workerId, "No script iterations finished, consider making the test duration longer", "DEBUG")
	}

	// Handle the end-of-test summary.
	if !testRunState.RuntimeOptions.NoSummary.Bool {
		engine.MetricsEngine.MetricsLock.Lock() // TODO: refactor so this is not needed
		summaryResult := &lib.Summary{
			Metrics:         engine.MetricsEngine.ObservedMetrics,
			RootGroup:       execScheduler.GetRunner().GetDefaultGroup(),
			TestRunDuration: executionState.GetCurrentTestRunDuration(),
			NoColor:         globalState.flags.noColor,
		}
		engine.MetricsEngine.MetricsLock.Unlock()
		summaryResultMarshalled, err := json.Marshal(summaryResult)
		if err == nil {
			dispatchMessage(ctx, client, job["id"], workerId, string(summaryResultMarshalled), "RESULTS")
		} else {
			handleError(ctx, client, job["id"], workerId, err)
		}
	}

	updateStatus(ctx, client, job["id"], workerId, "SUCCESS")

	globalCancel() // signal the Engine that it should wind down
	logger.Debug("Waiting for engine processes to finish...")
	engineWait()
	logger.Debug("Everything has finished, exiting k6!")
	if test.keyLogger != nil {
		if err := test.keyLogger.Close(); err != nil {
			logger.WithError(err).Warn("Error while closing the SSLKEYLOGFILE")
		}
	}
	if interrupt != nil {
		return
	}
	if engine.IsTainted() {
		fmt.Println(errext.WithExitCodeIfNone(errors.New("some thresholds have failed"), exitcodes.ThresholdsHaveFailed))
		return
	}
}
