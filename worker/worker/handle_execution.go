package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
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
func handleExecution(csdasdtx context.Context, client *redis.Client, job libOrch.ChildJob,
	workerId string, creditsClient *redis.Client, standalone bool) bool ***REMOVED***
	ctx := context.Background()

	gs := newGlobalState(ctx, client, job, workerId)

	libWorker.UpdateStatus(gs, "LOADING")

	workerInfo := loadWorkerInfo(ctx, client, job, workerId, gs, creditsClient, standalone)

	test, err := loadAndConfigureTest(gs, job, workerInfo)
	if err != nil ***REMOVED***
		libWorker.HandleStringError(gs, fmt.Sprintf("failed to load test: %s", err))
		return false
	***REMOVED***

	// Write the full options back to the Runner.
	testRunState, err := test.buildTestRunState(test.derivedConfig.Options)
	if err != nil ***REMOVED***
		libWorker.HandleStringError(gs, fmt.Sprintf("Error building testRunState %s", err.Error()))
		return false
	***REMOVED***

	startSubscription := workerInfo.Client.Subscribe(ctx, fmt.Sprintf("%s:go", job.ChildJobId))

	libWorker.UpdateStatus(gs, "READY")

	// Wait for start message on the channel
	<-startSubscription.Channel()
	startSubscription.Close()

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

	childUpdatesKey := fmt.Sprintf("childjobUserUpdates:%s", job.ChildJobId)

	childUpdatesSubscription := client.Subscribe(ctx, childUpdatesKey)
	childJobUpdatesChannel := childUpdatesSubscription.Channel()

	// Create handler for test aborts
	defer workerInfo.CreditsManager.StopCreditsCapturing()
	defer childUpdatesSubscription.Close()

	go func() ***REMOVED***
		for msg := range childJobUpdatesChannel ***REMOVED***
			var abortMessage = JobUserUpdate***REMOVED******REMOVED***
			if err := json.Unmarshal([]byte(msg.Payload), &abortMessage); err != nil ***REMOVED***
				libWorker.HandleStringError(gs, fmt.Sprintf("Error unmarshalling abort message: %s", err.Error()))
				continue
			***REMOVED***

			if abortMessage.UpdateType == "CANCEL" ***REMOVED***
				fmt.Println("Aborting child job due to a request from the orchestrator")
				runCancel()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	execScheduler, err := local.NewExecutionScheduler(testRunState)
	if err != nil ***REMOVED***
		libWorker.HandleStringError(gs, fmt.Sprintf("Error initializing the execution scheduler: %s", err.Error()))
		return false
	***REMOVED***

	// Create all outputs.
	outputs, err := createOutputs(gs)
	if err != nil ***REMOVED***
		libWorker.HandleStringError(gs, fmt.Sprintf("Error creating outputs %s", err.Error()))
		return false
	***REMOVED***

	// Create the engine.
	engine, err := core.NewEngine(testRunState, execScheduler, outputs)
	if err != nil ***REMOVED***
		libWorker.HandleStringError(gs, fmt.Sprintf("Error creating engine %s", err.Error()))
		return false
	***REMOVED***

	// Wait for the job to be started on redis
	// TODO: implement as a blocking redis call

	// We do this here so we can get any output URLs below.
	err = engine.OutputManager.StartOutputs()
	if err != nil ***REMOVED***
		libWorker.HandleStringError(gs, fmt.Sprintf("Error starting outputs %s", err.Error()))
		return false
	***REMOVED***
	defer engine.OutputManager.StopOutputs()

	// Initialize the engine
	engineRun, _, err := engine.Init(globalCtx, runCtx, workerInfo)
	if err != nil ***REMOVED***
		err = common.UnwrapGojaInterruptedError(err)
		// Add a generic engine exit code if we don't have a more specific one
		libWorker.HandleError(gs, errext.WithExitCodeIfNone(err, exitcodes.GenericEngine))
		return false
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

	executionState := execScheduler.GetState()

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
	if interrupt != nil ***REMOVED***
		return false
	***REMOVED***
	if engine.IsTainted() ***REMOVED***
		libWorker.HandleError(gs, errext.WithExitCodeIfNone(errors.New("some thresholds have failed"), exitcodes.ThresholdsHaveFailed))
		return false
	***REMOVED***

	return true
***REMOVED***

func (lct *workerLoadedAndConfiguredTest) buildTestRunState(
	configToReinject libWorker.Options,
) (*libWorker.TestRunState, error) ***REMOVED***
	// This might be the full derived or just the consolidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &libWorker.TestRunState***REMOVED***
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.derivedConfig.Options, // we will always run with the derived options
	***REMOVED***, nil
***REMOVED***

func loadWorkerInfo(ctx context.Context,
	client *redis.Client, job libOrch.ChildJob, workerId string, gs libWorker.BaseGlobalState,
	creditsClient *redis.Client, standalone bool) *libWorker.WorkerInfo ***REMOVED***
	workerInfo := &libWorker.WorkerInfo***REMOVED***
		Client:          client,
		JobId:           job.Id,
		ChildJobId:      job.ChildJobId,
		OrchestratorId:  job.AssignedOrchestrator,
		WorkerId:        workerId,
		Ctx:             ctx,
		WorkerOptions:   job.Options,
		Gs:              &gs,
		VerifiedDomains: job.VerifiedDomains,
		SubFraction:     job.SubFraction,
		CreditsManager:  lib.CreateCreditsManager(ctx, job.Scope.Variant, job.Scope.VariantTargetId, creditsClient),
		Standalone:      standalone,
	***REMOVED***

	if job.CollectionContext != nil && job.CollectionContext.Name != "" ***REMOVED***
		collectionVariables := make(map[string]string)

		for _, variable := range job.CollectionContext.Variables ***REMOVED***
			collectionVariables[variable.Key] = variable.Value
		***REMOVED***

		workerInfo.Collection = &libWorker.Collection***REMOVED***
			Variables: collectionVariables,
			Name:      job.CollectionContext.Name,
		***REMOVED***
	***REMOVED***

	if job.EnvironmentContext != nil && job.EnvironmentContext.Name != "" ***REMOVED***
		environmentVariables := make(map[string]string)

		for _, variable := range job.EnvironmentContext.Variables ***REMOVED***
			environmentVariables[variable.Key] = variable.Value
		***REMOVED***

		workerInfo.Environment = &libWorker.Environment***REMOVED***
			Variables: environmentVariables,
			Name:      job.EnvironmentContext.Name,
		***REMOVED***
	***REMOVED***

	workerInfo.FinalRequest = job.FinalRequest
	workerInfo.UnderlyingRequest = job.UnderlyingRequest

	return workerInfo
***REMOVED***
