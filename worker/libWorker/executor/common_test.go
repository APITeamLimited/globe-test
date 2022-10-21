package executor

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils/minirunner"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

func simpleRunner(vuFn func(context.Context, *libWorker.State) error) libWorker.Runner ***REMOVED***
	return &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, state *libWorker.State, _ chan<- workerMetrics.SampleContainer) error ***REMOVED***
			return vuFn(ctx, state)
		***REMOVED***,
	***REMOVED***
***REMOVED***

func getTestRunState(tb testing.TB, options libWorker.Options, runner libWorker.Runner) *libWorker.TestRunState ***REMOVED***
	reg := workerMetrics.NewRegistry()
	piState := &libWorker.TestPreInitState***REMOVED***
		Logger:         testutils.NewLogger(tb),
		RuntimeOptions: libWorker.RuntimeOptions***REMOVED******REMOVED***,
		Registry:       reg,
		BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(reg),
	***REMOVED***

	require.NoError(tb, runner.SetOptions(options))

	return &libWorker.TestRunState***REMOVED***
		TestPreInitState: piState,
		Options:          options,
		Runner:           runner,
	***REMOVED***
***REMOVED***

func setupExecutor(t testing.TB, config libWorker.ExecutorConfig, es *libWorker.ExecutionState) (
	context.Context, context.CancelFunc, libWorker.Executor, *testutils.SimpleLogrusHook,
) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	engineOut := make(chan workerMetrics.SampleContainer, 100) // TODO: return this for more complicated tests?

	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)
	logEntry := logrus.NewEntry(testLog)

	initVUFunc := func(_ context.Context, logger *logrus.Entry, workerInfo *libWorker.WorkerInfo) (libWorker.InitializedVU, error) ***REMOVED***
		idl, idg := es.GetUniqueVUIdentifiers()
		return es.Test.Runner.NewVU(idl, idg, engineOut, libWorker.GetTestWorkerInfo())
	***REMOVED***
	es.SetInitVUFunc(initVUFunc)

	maxPlannedVUs := libWorker.GetMaxPlannedVUs(config.GetExecutionRequirements(es.ExecutionTuple))
	initializeVUs(ctx, t, logEntry, es, maxPlannedVUs, initVUFunc)

	executor, err := config.NewExecutor(es, logEntry)
	require.NoError(t, err)

	err = executor.Init(ctx)
	require.NoError(t, err)
	return ctx, cancel, executor, logHook
***REMOVED***

func initializeVUs(
	ctx context.Context, t testing.TB, logEntry *logrus.Entry, es *libWorker.ExecutionState, number uint64, initVU libWorker.InitVUFunc,
) ***REMOVED***
	// This is not how the local ExecutionScheduler initializes VUs, but should do the same job
	for i := uint64(0); i < number; i++ ***REMOVED***
		// Not calling es.InitializeNewVU() here to avoid a double increment of initializedVUs,
		// which is done in es.AddInitializedVU().
		vu, err := initVU(ctx, logEntry, libWorker.GetTestWorkerInfo())
		require.NoError(t, err)
		es.AddInitializedVU(vu)
	***REMOVED***
***REMOVED***

type executorTest struct ***REMOVED***
	options libWorker.Options
	state   *libWorker.ExecutionState

	ctx      context.Context //nolint
	cancel   context.CancelFunc
	executor libWorker.Executor
	logHook  *testutils.SimpleLogrusHook
***REMOVED***

func setupExecutorTest(
	t testing.TB, segmentStr, sequenceStr string, extraOptions libWorker.Options,
	runner libWorker.Runner, config libWorker.ExecutorConfig,
) *executorTest ***REMOVED***
	var err error
	var segment *libWorker.ExecutionSegment
	if segmentStr != "" ***REMOVED***
		segment, err = libWorker.NewExecutionSegmentFromString(segmentStr)
		require.NoError(t, err)
	***REMOVED***

	var sequence libWorker.ExecutionSegmentSequence
	if sequenceStr != "" ***REMOVED***
		sequence, err = libWorker.NewExecutionSegmentSequenceFromString(sequenceStr)
		require.NoError(t, err)
	***REMOVED***

	et, err := libWorker.NewExecutionTuple(segment, &sequence)
	require.NoError(t, err)

	options := libWorker.Options***REMOVED***
		ExecutionSegment:         segment,
		ExecutionSegmentSequence: &sequence,
	***REMOVED***.Apply(runner.GetOptions()).Apply(extraOptions)

	testRunState := getTestRunState(t, options, runner)

	execReqs := config.GetExecutionRequirements(et)
	es := libWorker.NewExecutionState(testRunState, et, libWorker.GetMaxPlannedVUs(execReqs), libWorker.GetMaxPossibleVUs(execReqs))
	ctx, cancel, executor, logHook := setupExecutor(t, config, es)

	return &executorTest***REMOVED***
		options:  options,
		state:    es,
		ctx:      ctx,
		cancel:   cancel,
		executor: executor,
		logHook:  logHook,
	***REMOVED***
***REMOVED***
