package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis/v9"
	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
)

/*
This is the main function that is called when the worker is started.
It is responsible for running a job and reporting on its status
*/
func handleExecution(ctx context.Context,
	client *redis.Client, job map[string]string, workerId string) ***REMOVED***
	// Check if redis message is a uuid

	fmt.Println("\033[1;32mGot job", job["id"], "\033[0m")

	go updateStatus(ctx, client, job["id"], workerId, "LOADING")

	globalState := newGlobalState(ctx, client, job["id"], workerId)

	test, err := loadAndConfigureTest(globalState, job)

	if err != nil ***REMOVED***
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("failed to loadAndConfigureTest: %s", err.Error()))
		return
	***REMOVED***

	go dispatchMessage(ctx, client, job["id"], workerId, fmt.Sprintf("Loaded test %s", test.workerLoadedTest.sourceRootPath), "DEBUG")

	// Write the full consolidated *and derived* options back to the Runner.
	conf := test.derivedConfig
	testRunState, err := test.buildTestRunState(conf.Options)
	if err != nil ***REMOVED***
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error building testRunState %s", err.Error()))
		return
	***REMOVED***

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
	if err != nil ***REMOVED***
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error initializing the execution scheduler: %s", err.Error()))
		return
	***REMOVED***

	// Create all outputs.
	executionPlan := execScheduler.GetExecutionPlan()
	outputs, err := createOutputs(globalState, test, executionPlan)
	if err != nil ***REMOVED***
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error creating outputs %s", err.Error()))
		return
	***REMOVED***

	// TODO: create a MetricsEngine here and add its ingester to the list of
	// outputs (unless both NoThresholds and NoSummary were enabled)

	// TODO: remove this completely
	// Create the engine.
	go dispatchMessage(ctx, client, job["id"], workerId, "Initializing the Engine...", "DEBUG")
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil ***REMOVED***
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error creating engine %s", err.Error()))
		return
	***REMOVED***

	// We do this here so we can get any output URLs below.
	err = engine.OutputManager.StartOutputs()
	if err != nil ***REMOVED***
		go handleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error starting outputs %s", err.Error()))
		return
	***REMOVED***
	defer engine.OutputManager.StopOutputs()

	// Trap Interrupts, SIGINTs and SIGTERMs.
	gracefulStop := func(sig os.Signal) ***REMOVED***
		logger.WithField("sig", sig).Debug("Stopping k6 in response to signal...")
		lingerCancel() // stop the test run, metric processing is cancelled below
	***REMOVED***
	onHardStop := func(sig os.Signal) ***REMOVED***
		logger.WithField("sig", sig).Error("Aborting k6 in response to signal")
		globalCancel() // not that it matters, given the following command...
	***REMOVED***
	stopSignalHandling := handleTestAbortSignals(globalState, gracefulStop, onHardStop)
	defer stopSignalHandling()

	// Initialize the engine
	go dispatchMessage(ctx, client, job["id"], workerId, "Initializing VU(s)...", "DEBUG")
	engineRun, engineWait, err := engine.Init(globalCtx, runCtx)
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		go handleError(ctx, client, job["id"], workerId, errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		return
	***REMOVED***

	// Start the test run
	go updateStatus(ctx, client, job["id"], workerId, "RUNNING")
	var interrupt error
	err = engineRun()
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		if errext.IsInterruptError(err) ***REMOVED***
			// Don't return here since we need to work with --linger,
			// show the end-of-test summary and exit cleanly.
			interrupt = err
		***REMOVED***
		if !conf.Linger.Bool && interrupt == nil ***REMOVED***
			fmt.Println(errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		***REMOVED***
	***REMOVED***
	runCancel()
	go dispatchMessage(ctx, client, job["id"], workerId, "Engine run terminated cleanly", "DEBUG")

	executionState := execScheduler.GetState()
	// Warn if no iterations could be completed.
	if executionState.GetFullIterationCount() == 0 ***REMOVED***
		go dispatchMessage(ctx, client, job["id"], workerId, "No script iterations finished, consider making the test duration longer", "DEBUG")
	***REMOVED***

	// Handle the end-of-test summary.
	if !testRunState.RuntimeOptions.NoSummary.Bool ***REMOVED***
		engine.MetricsEngine.MetricsLock.Lock() // TODO: refactor so this is not needed
		summaryResult := &lib.Summary***REMOVED***
			Metrics:         engine.MetricsEngine.ObservedMetrics,
			RootGroup:       execScheduler.GetRunner().GetDefaultGroup(),
			TestRunDuration: executionState.GetCurrentTestRunDuration(),
			NoColor:         globalState.flags.noColor,
		***REMOVED***
		engine.MetricsEngine.MetricsLock.Unlock()
		summaryResultMarshalled, err := json.Marshal(summaryResult)
		if err == nil ***REMOVED***
			dispatchMessage(ctx, client, job["id"], workerId, string(summaryResultMarshalled), "RESULTS")
		***REMOVED*** else ***REMOVED***
			handleError(ctx, client, job["id"], workerId, err)
		***REMOVED***
	***REMOVED***

	updateStatus(ctx, client, job["id"], workerId, "SUCCESS")

	globalCancel() // signal the Engine that it should wind down
	logger.Debug("Waiting for engine processes to finish...")
	engineWait()
	logger.Debug("Everything has finished, exiting k6!")
	if test.keyLogger != nil ***REMOVED***
		if err := test.keyLogger.Close(); err != nil ***REMOVED***
			logger.WithError(err).Warn("Error while closing the SSLKEYLOGFILE")
		***REMOVED***
	***REMOVED***
	if interrupt != nil ***REMOVED***
		return
	***REMOVED***
	if engine.IsTainted() ***REMOVED***
		fmt.Println(errext.WithExitCodeIfNone(errors.New("some thresholds have failed"), exitcodes.ThresholdsHaveFailed))
		return
	***REMOVED***
***REMOVED***
