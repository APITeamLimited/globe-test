package worker

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/APITeamLimited/redis/v9"
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
	client *redis.Client, job map[string]string, workerId string) {
	fmt.Println("\033[1;32mGot job", job["id"], "\033[0m")
	go lib.UpdateStatus(ctx, client, job["id"], workerId, "LOADING")

	globalState := newGlobalState(ctx, client, job["id"], workerId)

	workerInfo, err := loadWorkerInfo(ctx, client, job, workerId)
	if err != nil {
		go lib.HandleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("failed to load test: %s", err))
		return
	}

	test, err := loadAndConfigureTest(globalState, job, workerInfo)
	if err != nil {
		go lib.HandleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("failed to load test: %s", err))
		return
	}

	go lib.DispatchMessage(ctx, client, job["id"], workerId, fmt.Sprintf("Loaded test %s", test.workerLoadedTest.sourceRootPath), "DEBUG")

	// Write the full consolidated *and derived* options back to the Runner.
	testRunState, err := test.buildTestRunState(test.derivedConfig.Options)
	if err != nil {
		go lib.HandleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error building testRunState %s", err.Error()))
		return
	}

	// Don't know if these can be removed easily without unexpected side effects

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
	globalCtx, globalCancel := context.WithCancel(ctx)
	defer globalCancel()
	lingerCtx, lingerCancel := context.WithCancel(globalCtx)
	defer lingerCancel()
	runCtx, runCancel := context.WithCancel(lingerCtx)
	defer runCancel()

	execScheduler, err := local.NewExecutionScheduler(testRunState)
	if err != nil {
		go lib.HandleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error initializing the execution scheduler: %s", err.Error()))
		return
	}

	// Create all outputs.
	executionPlan := execScheduler.GetExecutionPlan()
	outputs, err := createOutputs(globalState, test, executionPlan, workerInfo)
	if err != nil {
		go lib.HandleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error creating outputs %s", err.Error()))
		return
	}

	// TODO: create a MetricsEngine here and add its ingester to the list of
	// outputs (unless both NoThresholds and NoSummary were enabled)

	// TODO: remove this completely
	// Create the engine.
	go lib.DispatchMessage(ctx, client, job["id"], workerId, "Initializing the Engine...", "DEBUG")
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil {
		go lib.HandleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error creating engine %s", err.Error()))
		return
	}

	// We do this here so we can get any output URLs below.
	err = engine.OutputManager.StartOutputs()
	if err != nil {
		go lib.HandleStringError(ctx, client, job["id"], workerId, fmt.Sprintf("Error starting outputs %s", err.Error()))
		return
	}
	defer engine.OutputManager.StopOutputs()

	// Trap Interrupts, SIGINTs and SIGTERMs.
	gracefulStop := func(sig os.Signal) {
		go lib.DispatchMessage(ctx, client, job["id"], workerId, fmt.Sprintf("Stopping worker in response to signal %s", sig), "DEBUG")
		lingerCancel() // stop the test run, metric processing is cancelled below
	}
	onHardStop := func(sig os.Signal) {
		go lib.DispatchMessage(ctx, client, job["id"], workerId, fmt.Sprintf("Hard stop in response to signal %s", sig), "DEBUG")
		globalCancel() // not that it matters, given the following command...
	}
	stopSignalHandling := handleTestAbortSignals(globalState, gracefulStop, onHardStop)
	defer stopSignalHandling()

	// Initialize the engine
	go lib.DispatchMessage(ctx, client, job["id"], workerId, "Initializing VU(s)...", "DEBUG")
	engineRun, engineWait, err := engine.Init(globalCtx, runCtx, workerInfo)
	if err != nil {
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		go lib.HandleError(ctx, client, job["id"], workerId, errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		return
	}

	// Start the test run
	go lib.UpdateStatus(ctx, client, job["id"], workerId, "RUNNING")
	var interrupt error
	err = engineRun()
	if err != nil {
		err = common.UnwrapGojaInterruptedError(err)
		if errext.IsInterruptError(err) {
			// Don't return here since we need to work with --linger,
			// show the end-of-test summary and exit cleanly.
			interrupt = err
		}
		if !test.derivedConfig.Linger.Bool && interrupt == nil {
			fmt.Println(errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		}
	}
	runCancel()
	go lib.DispatchMessage(ctx, client, job["id"], workerId, "Engine run terminated cleanly", "DEBUG")

	executionState := execScheduler.GetState()
	// Warn if no iterations could be completed.
	if executionState.GetFullIterationCount() == 0 {
		go lib.DispatchMessage(ctx, client, job["id"], workerId, "No script iterations finished, consider making the test duration longer", "DEBUG")
	}

	// Handle the end-of-test summary.
	if !testRunState.RuntimeOptions.NoSummary.Bool {
		engine.MetricsEngine.MetricsLock.Lock() // TODO: refactor so this is not needed
		marshalledMetrics, err := test.initRunner.RetrieveMetricsJSON(globalCtx, &lib.Summary{
			Metrics:         engine.MetricsEngine.ObservedMetrics,
			RootGroup:       execScheduler.GetRunner().GetDefaultGroup(),
			TestRunDuration: executionState.GetCurrentTestRunDuration(),
		})
		engine.MetricsEngine.MetricsLock.Unlock()

		if err == nil {
			go lib.DispatchMessage(ctx, client, job["id"], workerId, string(marshalledMetrics), "SUMMARY_METRICS")
		} else {
			go lib.HandleError(ctx, client, job["id"], workerId, err)
		}
	}

	lib.UpdateStatus(ctx, client, job["id"], workerId, "SUCCESS")

	globalCancel() // signal the Engine that it should wind down
	go lib.DispatchMessage(ctx, client, job["id"], workerId, "Waiting for the Engine to finish...", "DEBUG")
	engineWait()
	go lib.DispatchMessage(ctx, client, job["id"], workerId, "Everything has finished, exiting worker", "DEBUG")
	if interrupt != nil {
		return
	}
	if engine.IsTainted() {
		go lib.HandleError(ctx, client, job["id"], workerId, errext.WithExitCodeIfNone(errors.New("some thresholds have failed"), exitcodes.ThresholdsHaveFailed))
		return
	}
}
