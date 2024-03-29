package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/js/common"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/core"
	"github.com/APITeamLimited/globe-test/worker/core/local"
	"github.com/APITeamLimited/globe-test/worker/errext"
	"github.com/APITeamLimited/globe-test/worker/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/redis/v9"
	"github.com/gorilla/websocket"
)

/*
This is the main function that is called when the worker is started.
It is responsible for running a job and reporting on its status
*/
func handleExecution(ctx context.Context, conn *websocket.Conn, job *libOrch.ChildJob,
	workerId string, creditsClient *redis.Client, standalone bool, connReadMutex, connWriteMutex *sync.Mutex) bool {
	gs := newGlobalState(ctx, conn, job, workerId, job.FuncModeInfo, connReadMutex, connWriteMutex)
	eventChannels := getEventChannels(gs)
	defer close(eventChannels.goMessageChannel)
	defer close(eventChannels.childUpdatesChannel)

	// Prestart-abort callback is for when not yet running but the test is aborted
	preAbortChannel := make(chan bool)
	gs.SetRunAbortFunc(func() {
		preAbortChannel <- true
	})

	libWorker.UpdateStatus(gs, "LOADING")
	workerInfo := loadWorkerInfo(ctx, conn, job, workerId, gs, creditsClient, standalone)

	test, err := loadAndConfigureTest(gs, job, workerInfo)
	if err != nil {
		libWorker.HandleStringError(gs, fmt.Sprintf("failed to load test: %s", err))
		return false
	}

	// Write the full options back to the Runner.
	testRunState, err := test.buildTestRunState(test.derivedConfig.Options)
	if err != nil {
		libWorker.HandleStringError(gs, fmt.Sprintf("Error building testRunState %s", err.Error()))
		return false
	}

	startChannel := testStartChannel(gs, eventChannels, preAbortChannel)
	libWorker.UpdateStatus(gs, "READY")

	// Wait for the start signal from the orchestrator
	startTime := <-startChannel
	if startTime == nil {
		return false
	}

	// Wait till start time
	if startTime.After(time.Now()) {
		time.Sleep(time.Until(*startTime))
	}

	// Test starts here

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

	gs.SetRunAbortFunc(func() {
		runCancel()

		time.Sleep(1 * time.Second)

		globalCancel()
	})

	// Regularly deduct credits
	workerInfo.CreditsManager.StartMonitoringCredits(func() {
		libWorker.HandleStringError(gs, "Test stopped due to lack of credits")
		runCancel()
	})
	defer func() {
		// Run as goroutine to reduce response time
		go workerInfo.CreditsManager.BillFinalCredits()
	}()

	execScheduler, err := local.NewExecutionScheduler(testRunState)
	if err != nil {
		libWorker.HandleStringError(gs, fmt.Sprintf("Error initializing the execution scheduler: %s", err.Error()))
		return false
	}

	// Create all outputs.
	outputs, err := createOutputs(gs, job.Location)
	if err != nil {
		libWorker.HandleStringError(gs, fmt.Sprintf("Error creating outputs %s", err.Error()))
		return false
	}

	// Create the engine.
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil {
		libWorker.HandleStringError(gs, fmt.Sprintf("Error creating engine %s", err.Error()))
		return false
	}

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

	go func() {
		for msg := range eventChannels.childUpdatesChannel {
			var updateMessage = JobUserUpdate{}
			if err := json.Unmarshal([]byte(msg), &updateMessage); err != nil {
				libWorker.HandleStringError(gs, fmt.Sprintf("Error unmarshalling abort message: %s", err.Error()))
				continue
			}

			if updateMessage.UpdateType == "CANCEL" {
				fmt.Println("Aborting child job due to a request from the orchestrator")

				runCancel()
			}
		}
	}()

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

	libWorker.UpdateStatus(gs, "SUCCESS")

	globalCancel() // signal the Engine that it should wind down

	return interrupt == nil
}

func (lct *workerLoadedAndConfiguredTest) buildTestRunState(
	configToReinject libWorker.Options,
) (*libWorker.TestRunState, error) {
	// This might be the full derived or just the consolidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil {
		return nil, err
	}

	return &libWorker.TestRunState{
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.derivedConfig.Options, // we will always run with the derived options
	}, nil
}

func loadWorkerInfo(ctx context.Context,
	conn *websocket.Conn, job *libOrch.ChildJob, workerId string, gs libWorker.BaseGlobalState,
	creditsClient *redis.Client, standalone bool) *libWorker.WorkerInfo {
	workerInfo := &libWorker.WorkerInfo{
		Conn:            conn,
		JobId:           job.Id,
		ChildJobId:      job.ChildJobId,
		OrchestratorId:  job.AssignedOrchestrator,
		WorkerId:        workerId,
		Ctx:             ctx,
		WorkerOptions:   job.ChildOptions,
		Gs:              &gs,
		VerifiedDomains: job.VerifiedDomains,
		SubFraction:     job.SubFraction,
		Standalone:      standalone,
	}

	workerInfo.DomainLimiter = libWorker.CreateDomainLimiter(standalone, job.VerifiedDomains, workerInfo)

	if gs.FuncModeInfo() != nil && creditsClient != nil {
		workerInfo.CreditsManager = lib.CreateCreditsManager(ctx, job.Scope.Variant, job.Scope.VariantTargetId, creditsClient, *gs.FuncModeInfo())
	}

	if job.CollectionContext != nil && job.CollectionContext.Name != "" {
		collectionVariables := make(map[string]string)

		for _, variable := range job.CollectionContext.Variables {
			collectionVariables[variable.Key] = variable.Value
		}

		workerInfo.Collection = &libWorker.Collection{
			Variables: collectionVariables,
			Name:      job.CollectionContext.Name,
		}
	}

	if job.EnvironmentContext != nil && job.EnvironmentContext.Name != "" {
		environmentVariables := make(map[string]string)

		for _, variable := range job.EnvironmentContext.Variables {
			environmentVariables[variable.Key] = variable.Value
		}

		workerInfo.Environment = &libWorker.Environment{
			Variables: environmentVariables,
			Name:      job.EnvironmentContext.Name,
		}
	}

	workerInfo.TestData = *job.TestData

	return workerInfo
}
