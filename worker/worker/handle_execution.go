package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/core"
	"github.com/APITeamLimited/globe-test/worker/core/local"
	"github.com/APITeamLimited/globe-test/worker/errext"
	"github.com/APITeamLimited/globe-test/worker/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
)

/*
This is the main function that is called when the worker is started.
It is responsible for running a job and reporting on its status
*/
func handleExecution(ctx context.Context,
	client *redis.Client, job libOrch.ChildJob, workerId string) ***REMOVED***
	fmt.Printf("\033[1;32mGot child job %s\033[0m\n", job.ChildJobId)
	gs := newGlobalState(ctx, client, job, workerId)

	libWorker.UpdateStatus(gs, "LOADING")

	workerInfo := loadWorkerInfo(ctx, client, job, workerId, gs)

	test, err := loadAndConfigureTest(gs, job, workerInfo)
	if err != nil ***REMOVED***
		go libWorker.HandleStringError(gs, fmt.Sprintf("failed to load test: %s", err))
		return
	***REMOVED***

	// Write the full options back to the Runner.
	testRunState, err := test.buildTestRunState(test.derivedConfig.Options)
	if err != nil ***REMOVED***
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error building testRunState %s", err.Error()))
		return
	***REMOVED***

	// Don't know if these can be removed easily without unexpected side effects

	// We prepare a bunch of contexts:
	//  - The runCtx is cancelled as soon as the Engine's run() lambda finishes,
	//    and can trigger things like the usage report and end of test summary.
	//    Crucially, metrics processing by the Engine will still work after this
	//    context is cancelled!
	//  - The globalCtx is cancelled only after we're completely done with the
	//    test execution, so that the Engine  can start winding down its metrics
	//    processing.
	globalCtx, globalCancel := context.WithCancel(ctx)
	defer globalCancel()
	runCtx, runCancel := context.WithCancel(globalCtx)
	defer runCancel()

	execScheduler, err := local.NewExecutionScheduler(testRunState)
	if err != nil ***REMOVED***
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error initializing the execution scheduler: %s", err.Error()))
		return
	***REMOVED***

	// Create all outputs.
	outputs, err := createOutputs(gs)
	if err != nil ***REMOVED***
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error creating outputs %s", err.Error()))
		return
	***REMOVED***

	// Create the engine.
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil ***REMOVED***
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error creating engine %s", err.Error()))
		return
	***REMOVED***

	// Wait for the job to be started on redis
	// TODO: implement as a blocking redis call

	startChannel := workerInfo.Client.Subscribe(ctx, fmt.Sprintf("%s:go", job.ChildJobId)).Channel()

	libWorker.UpdateStatus(gs, "READY")

	// Wait for start message on the channel
	<-startChannel

	// We do this here so we can get any output URLs below.
	err = engine.OutputManager.StartOutputs()
	if err != nil ***REMOVED***
		libWorker.HandleStringError(gs, fmt.Sprintf("Error starting outputs %s", err.Error()))
		return
	***REMOVED***
	defer engine.OutputManager.StopOutputs()

	// Trap Interrupts, SIGINTs and SIGTERMs.
	gracefulStop := func(sig os.Signal) ***REMOVED***
		libWorker.DispatchMessage(gs, fmt.Sprintf("Stopping worker in response to signal %s", sig), "DEBUG")
	***REMOVED***
	onHardStop := func(sig os.Signal) ***REMOVED***
		libWorker.DispatchMessage(gs, fmt.Sprintf("Hard stop in response to signal %s", sig), "DEBUG")
		globalCancel() // not that it matters, given the following command...
	***REMOVED***
	stopSignalHandling := handleTestAbortSignals(gs, gracefulStop, onHardStop)
	defer stopSignalHandling()

	// Initialize the engine
	libWorker.DispatchMessage(gs, "Initializing VU(s)...", "DEBUG")
	engineRun, engineWait, err := engine.Init(globalCtx, runCtx, workerInfo)
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		libWorker.HandleError(gs, errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		return
	***REMOVED***

	// Start the test run
	libWorker.UpdateStatus(gs, "RUNNING")
	var interrupt error
	err = engineRun()
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		if errext.IsInterruptError(err) ***REMOVED***
			interrupt = err
		***REMOVED***
		if interrupt == nil ***REMOVED***
			fmt.Println(errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		***REMOVED***
	***REMOVED***
	runCancel()
	libWorker.DispatchMessage(gs, "Engine run terminated cleanly", "DEBUG")

	executionState := execScheduler.GetState()
	// Warn if no iterations could be completed.
	if executionState.GetFullIterationCount() == 0 ***REMOVED***
		libWorker.DispatchMessage(gs, "No script iterations finished, consider making the test duration longer", "DEBUG")
	***REMOVED***

	engine.MetricsEngine.MetricsLock.Lock() // TODO: refactor so this is not needed
	marshalledMetrics, err := test.initRunner.RetrieveMetricsJSON(globalCtx, &libWorker.Summary***REMOVED***
		Metrics:         engine.MetricsEngine.ObservedMetrics,
		RootGroup:       execScheduler.GetRunner().GetDefaultGroup(),
		TestRunDuration: executionState.GetCurrentTestRunDuration(),
	***REMOVED***)
	engine.MetricsEngine.MetricsLock.Unlock()

	// Retrive collection and environment variables
	if workerInfo.Collection != nil ***REMOVED***
		collectionVariables, err := json.Marshal(workerInfo.Collection.Variables)

		if err != nil ***REMOVED***
			libWorker.HandleStringError(gs, fmt.Sprintf("Error marshalling collection variables %s", err.Error()))
		***REMOVED*** else ***REMOVED***
			libWorker.DispatchMessage(gs, string(collectionVariables), "COLLECTION_VARIABLES")
		***REMOVED***
	***REMOVED***

	if workerInfo.Environment != nil ***REMOVED***
		environmentVariables, err := json.Marshal(workerInfo.Environment.Variables)

		if err != nil ***REMOVED***
			libWorker.HandleStringError(gs, fmt.Sprintf("Error marshalling environment variables %s", err.Error()))
		***REMOVED*** else ***REMOVED***
			libWorker.DispatchMessage(gs, string(environmentVariables), "ENVIRONMENT_VARIABLES")
		***REMOVED***
	***REMOVED***

	if err == nil ***REMOVED***
		libWorker.DispatchMessage(gs, string(marshalledMetrics), "SUMMARY_METRICS")
	***REMOVED*** else ***REMOVED***
		libWorker.HandleError(gs, err)
	***REMOVED***

	libWorker.UpdateStatus(gs, "SUCCESS")

	globalCancel() // signal the Engine that it should wind down
	libWorker.DispatchMessage(gs, "Waiting for the Engine to finish...", "DEBUG")
	engineWait()
	libWorker.DispatchMessage(gs, "Everything has finished, exiting worker", "DEBUG")
	if interrupt != nil ***REMOVED***
		return
	***REMOVED***
	if engine.IsTainted() ***REMOVED***
		libWorker.HandleError(gs, errext.WithExitCodeIfNone(errors.New("some thresholds have failed"), exitcodes.ThresholdsHaveFailed))
		return
	***REMOVED***
***REMOVED***

func (lct *workerLoadedAndConfiguredTest) buildTestRunState(
	configToReinject libWorker.Options,
) (*libWorker.TestRunState, error) ***REMOVED***
	// This might be the full derived or just the consolidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO: init atlas root worker, etc.

	return &libWorker.TestRunState***REMOVED***
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.derivedConfig.Options, // we will always run with the derived options
	***REMOVED***, nil
***REMOVED***
