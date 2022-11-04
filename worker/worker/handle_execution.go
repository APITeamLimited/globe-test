package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
func handleExecution(csdasdtx context.Context,
	client *redis.Client, job libOrch.ChildJob, workerId string) bool {
	ctx := context.Background()

	gs := newGlobalState(ctx, client, job, workerId)

	libWorker.UpdateStatus(gs, "LOADING")

	workerInfo := loadWorkerInfo(ctx, client, job, workerId, gs)

	test, err := loadAndConfigureTest(gs, job, workerInfo)
	if err != nil {
		go libWorker.HandleStringError(gs, fmt.Sprintf("failed to load test: %s", err))
		return false
	}

	// Write the full options back to the Runner.
	testRunState, err := test.buildTestRunState(test.derivedConfig.Options)
	if err != nil {
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error building testRunState %s", err.Error()))
		return false
	}

	startChannel := workerInfo.Client.Subscribe(ctx, fmt.Sprintf("%s:go", job.ChildJobId)).Channel()

	libWorker.UpdateStatus(gs, "READY")

	// Wait for start message on the channel
	<-startChannel

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

	// Create handler for test aborts
	childJobUpdatesChannel := client.Subscribe(ctx, fmt.Sprintf("childJobUserUpdates:%s", job.ChildJobId)).Channel()

	go func() {
		for msg := range childJobUpdatesChannel {
			var abortMessage = jobUserUpdate{}
			if err := json.Unmarshal([]byte(msg.Payload), &abortMessage); err != nil {
				libWorker.HandleStringError(gs, fmt.Sprintf("Error unmarshalling abort message: %s", err.Error()))
				continue
			}

			if abortMessage.UpdateType == "CANCEL" {
				fmt.Println("Aborting child job due to a request from the orchestrator")
				runCancel()
				return
			}
		}
	}()

	execScheduler, err := local.NewExecutionScheduler(testRunState)
	if err != nil {
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error initializing the execution scheduler: %s", err.Error()))
		return false
	}

	// Create all outputs.
	outputs, err := createOutputs(gs)
	if err != nil {
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error creating outputs %s", err.Error()))
		return false
	}

	// Create the engine.
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil {
		go libWorker.HandleStringError(gs, fmt.Sprintf("Error creating engine %s", err.Error()))
		return false
	}

	// Wait for the job to be started on redis
	// TODO: implement as a blocking redis call

	// We do this here so we can get any output URLs below.
	err = engine.OutputManager.StartOutputs()
	if err != nil {
		libWorker.HandleStringError(gs, fmt.Sprintf("Error starting outputs %s", err.Error()))
		return false
	}
	defer engine.OutputManager.StopOutputs()

	// Initialize the engine
	engineRun, _, err := engine.Init(globalCtx, runCtx, workerInfo)
	if err != nil {
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		libWorker.HandleError(gs, errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		return false
	}

	// Start the test run
	libWorker.UpdateStatus(gs, "RUNNING")
	var interrupt error
	err = engineRun()
	if err != nil {
		err = common.UnwrapGojaInterruptedError(err)
		if errext.IsInterruptError(err) {
			interrupt = err
		}
		if interrupt == nil {
			fmt.Println(errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		}
	}

	runCancel()

	executionState := execScheduler.GetState()

	engine.MetricsEngine.MetricsLock.Lock() // TODO: refactor so this is not needed
	marshalledMetrics, err := test.initRunner.RetrieveMetricsJSON(globalCtx, &libWorker.Summary{
		Metrics:         engine.MetricsEngine.ObservedMetrics,
		RootGroup:       execScheduler.GetRunner().GetDefaultGroup(),
		TestRunDuration: executionState.GetCurrentTestRunDuration(),
	})
	engine.MetricsEngine.MetricsLock.Unlock()

	// Retrive collection and environment variables
	if workerInfo.Collection != nil {
		collectionVariables, err := json.Marshal(workerInfo.Collection.Variables)

		if err != nil {
			libWorker.HandleStringError(gs, fmt.Sprintf("Error marshalling collection variables %s", err.Error()))
		} else {
			libWorker.DispatchMessage(gs, string(collectionVariables), "COLLECTION_VARIABLES")
		}
	}

	if workerInfo.Environment != nil {
		environmentVariables, err := json.Marshal(workerInfo.Environment.Variables)

		if err != nil {
			libWorker.HandleStringError(gs, fmt.Sprintf("Error marshalling environment variables %s", err.Error()))
		} else {
			libWorker.DispatchMessage(gs, string(environmentVariables), "ENVIRONMENT_VARIABLES")
		}
	}

	if err == nil {
		libWorker.DispatchMessage(gs, string(marshalledMetrics), "SUMMARY_METRICS")
	} else {
		libWorker.HandleError(gs, err)
	}

	libWorker.UpdateStatus(gs, "SUCCESS")

	globalCancel() // signal the Engine that it should wind down
	if interrupt != nil {
		return false
	}
	if engine.IsTainted() {
		libWorker.HandleError(gs, errext.WithExitCodeIfNone(errors.New("some thresholds have failed"), exitcodes.ThresholdsHaveFailed))
		return false
	}

	return true
}

func (lct *workerLoadedAndConfiguredTest) buildTestRunState(
	configToReinject libWorker.Options,
) (*libWorker.TestRunState, error) {
	// This might be the full derived or just the consolidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil {
		return nil, err
	}

	// TODO: init atlas root worker, etc.

	return &libWorker.TestRunState{
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.derivedConfig.Options, // we will always run with the derived options
	}, nil
}
