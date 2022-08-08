package node

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/afero"
	"go.k6.io/k6/core"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/ui/pb"
)

/*
This is the main function that is called when the node is started.
It is responsible for running a job and reporting on its status
*/
func handleExecution(ctx context.Context,
	client *redis.Client, job map[string]string, nodeId string) ***REMOVED***
	// Check if redis message is a uuid

	fmt.Println("\033[1;32mGot job", job["id"], "\033[0m")

	go updateStatus(ctx, client, job["id"], nodeId, "LOADING")

	globalState := newGlobalState(ctx, client, job["id"], nodeId)

	test, err := loadAndConfigureTest(globalState, job)

	if err != nil ***REMOVED***
		go handleError(ctx, client, job["id"], nodeId, fmt.Errorf("failed to loadAndConfigureTest: %s", err.Error()))
		return
	***REMOVED***

	fmt.Println("Loaded test", test.nodeLoadedTest.sourceRootPath)

	// Write the full consolidated *and derived* options back to the Runner.
	conf := test.derivedConfig
	testRunState, err := test.buildTestRunState(conf.Options)
	if err != nil ***REMOVED***
		go handleError(ctx, client, job["id"], nodeId, fmt.Errorf("Error building testRunState", err))
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
		go handleError(ctx, client, job["id"], nodeId, fmt.Errorf("Error initializing the execution scheduler:", err))
		return
	***REMOVED***

	// This is manually triggered after the Engine's Run() has completed,
	// and things like a single Ctrl+C don't affect it. We use it to make
	// sure that the progressbars finish updating with the latest execution
	// state one last time, after the test run has finished.
	_, progressCancel := context.WithCancel(globalCtx)
	defer progressCancel()
	initBar := execScheduler.GetInitProgressBar()
	progressBarWG := &sync.WaitGroup***REMOVED******REMOVED***
	progressBarWG.Add(1)
	go func() ***REMOVED***
		pbs := []*pb.ProgressBar***REMOVED***execScheduler.GetInitProgressBar()***REMOVED***
		for _, s := range execScheduler.GetExecutors() ***REMOVED***
			pbs = append(pbs, s.GetProgress())
		***REMOVED***
		progressBarWG.Done()
	***REMOVED***()

	// Create all outputs.
	executionPlan := execScheduler.GetExecutionPlan()
	outputs, err := createOutputs(globalState, test, executionPlan)
	if err != nil ***REMOVED***
		go handleError(ctx, client, job["id"], nodeId, fmt.Errorf("Error creating outputs", err))
		return
	***REMOVED***

	// TODO: create a MetricsEngine here and add its ingester to the list of
	// outputs (unless both NoThresholds and NoSummary were enabled)

	// TODO: remove this completely
	// Create the engine.
	initBar.Modify(pb.WithConstProgress(0, "Init engine"))
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil ***REMOVED***
		go handleError(ctx, client, job["id"], nodeId, fmt.Errorf("Error creating engine", err))
		return
	***REMOVED***

	// We do this here so we can get any output URLs below.
	initBar.Modify(pb.WithConstProgress(0, "Starting outputs"))
	// TODO: directly create the MutputManager here, not in the Engine
	err = engine.OutputManager.StartOutputs()
	if err != nil ***REMOVED***
		go handleError(ctx, client, job["id"], nodeId, fmt.Errorf("Error starting outputs", err))
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
	go dispatchMessage(ctx, client, job["id"], nodeId, "Init VUs...")
	engineRun, engineWait, err := engine.Init(globalCtx, runCtx)
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		go handleError(ctx, client, job["id"], nodeId, errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		return
	***REMOVED***

	// Start the test run
	go updateStatus(ctx, client, job["id"], nodeId, "RUNNING")
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
	go dispatchMessage(ctx, client, job["id"], nodeId, "Engine run terminated cleanly")

	progressCancel()

	executionState := execScheduler.GetState()
	// Warn if no iterations could be completed.
	if executionState.GetFullIterationCount() == 0 ***REMOVED***
		go dispatchMessage(ctx, client, job["id"], nodeId, "No script iterations finished, consider making the test duration longer")
	***REMOVED***

	// Handle the end-of-test summary.
	if !testRunState.RuntimeOptions.NoSummary.Bool ***REMOVED***
		engine.MetricsEngine.MetricsLock.Lock() // TODO: refactor so this is not needed
		summaryResult, err := test.initRunner.HandleSummary(globalCtx, &lib.Summary***REMOVED***
			Metrics:         engine.MetricsEngine.ObservedMetrics,
			RootGroup:       execScheduler.GetRunner().GetDefaultGroup(),
			TestRunDuration: executionState.GetCurrentTestRunDuration(),
			NoColor:         globalState.flags.noColor,
		***REMOVED***)
		engine.MetricsEngine.MetricsLock.Unlock()
		if err == nil ***REMOVED***
			fmt.Println("TODO: handle summary result")
			err = handleSummaryResult(globalState.fs, globalState.stdOut, globalState.stdErr, summaryResult)
		***REMOVED***
		if err != nil ***REMOVED***
			handleError(ctx, client, job["id"], nodeId, err)
		***REMOVED***
	***REMOVED***

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

func handleSummaryResult(fs afero.Fs, stdOut, stdErr io.Writer, result map[string]io.Reader) error ***REMOVED***
	var errs []error

	getWriter := func(path string) (io.Writer, error) ***REMOVED***
		switch path ***REMOVED***
		case "stdout":
			return stdOut, nil
		case "stderr":
			return stdErr, nil
		default:
			return fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		***REMOVED***
	***REMOVED***

	for path, value := range result ***REMOVED***
		if writer, err := getWriter(path); err != nil ***REMOVED***
			errs = append(errs, fmt.Errorf("could not open '%s': %w", path, err))
		***REMOVED*** else if n, err := io.Copy(writer, value); err != nil ***REMOVED***
			errs = append(errs, fmt.Errorf("error saving summary to '%s' after %d bytes: %w", path, n, err))
		***REMOVED***
	***REMOVED***

	return consolidateErrorMessage(errs, "Could not save some summary information:")
***REMOVED***
