package executor

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/lib/testutils"
	"github.com/APITeamLimited/k6-worker/lib/testutils/minirunner"
	"github.com/APITeamLimited/k6-worker/metrics"
)

func simpleRunner(vuFn func(context.Context, *lib.State) error) lib.Runner ***REMOVED***
	return &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, state *lib.State, _ chan<- metrics.SampleContainer) error ***REMOVED***
			return vuFn(ctx, state)
		***REMOVED***,
	***REMOVED***
***REMOVED***

func getTestRunState(tb testing.TB, options lib.Options, runner lib.Runner) *lib.TestRunState ***REMOVED***
	reg := metrics.NewRegistry()
	piState := &lib.TestPreInitState***REMOVED***
		Logger:         testutils.NewLogger(tb),
		RuntimeOptions: lib.RuntimeOptions***REMOVED******REMOVED***,
		Registry:       reg,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(reg),
	***REMOVED***

	require.NoError(tb, runner.SetOptions(options))

	return &lib.TestRunState***REMOVED***
		TestPreInitState: piState,
		Options:          options,
		Runner:           runner,
	***REMOVED***
***REMOVED***

func setupExecutor(t testing.TB, config lib.ExecutorConfig, es *lib.ExecutionState) (
	context.Context, context.CancelFunc, lib.Executor, *testutils.SimpleLogrusHook,
) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	engineOut := make(chan metrics.SampleContainer, 100) // TODO: return this for more complicated tests?

	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)
	logEntry := logrus.NewEntry(testLog)

	initVUFunc := func(_ context.Context, logger *logrus.Entry, workerInfo *lib.WorkerInfo) (lib.InitializedVU, error) ***REMOVED***
		idl, idg := es.GetUniqueVUIdentifiers()
		return es.Test.Runner.NewVU(idl, idg, engineOut, lib.GetTestWorkerInfo())
	***REMOVED***
	es.SetInitVUFunc(initVUFunc)

	maxPlannedVUs := lib.GetMaxPlannedVUs(config.GetExecutionRequirements(es.ExecutionTuple))
	initializeVUs(ctx, t, logEntry, es, maxPlannedVUs, initVUFunc)

	executor, err := config.NewExecutor(es, logEntry)
	require.NoError(t, err)

	err = executor.Init(ctx)
	require.NoError(t, err)
	return ctx, cancel, executor, logHook
***REMOVED***

func initializeVUs(
	ctx context.Context, t testing.TB, logEntry *logrus.Entry, es *lib.ExecutionState, number uint64, initVU lib.InitVUFunc,
) ***REMOVED***
	// This is not how the local ExecutionScheduler initializes VUs, but should do the same job
	for i := uint64(0); i < number; i++ ***REMOVED***
		// Not calling es.InitializeNewVU() here to avoid a double increment of initializedVUs,
		// which is done in es.AddInitializedVU().
		vu, err := initVU(ctx, logEntry, lib.GetTestWorkerInfo())
		require.NoError(t, err)
		es.AddInitializedVU(vu)
	***REMOVED***
***REMOVED***

type executorTest struct ***REMOVED***
	options lib.Options
	state   *lib.ExecutionState

	ctx      context.Context //nolint
	cancel   context.CancelFunc
	executor lib.Executor
	logHook  *testutils.SimpleLogrusHook
***REMOVED***

func setupExecutorTest(
	t testing.TB, segmentStr, sequenceStr string, extraOptions lib.Options,
	runner lib.Runner, config lib.ExecutorConfig,
) *executorTest ***REMOVED***
	var err error
	var segment *lib.ExecutionSegment
	if segmentStr != "" ***REMOVED***
		segment, err = lib.NewExecutionSegmentFromString(segmentStr)
		require.NoError(t, err)
	***REMOVED***

	var sequence lib.ExecutionSegmentSequence
	if sequenceStr != "" ***REMOVED***
		sequence, err = lib.NewExecutionSegmentSequenceFromString(sequenceStr)
		require.NoError(t, err)
	***REMOVED***

	et, err := lib.NewExecutionTuple(segment, &sequence)
	require.NoError(t, err)

	options := lib.Options***REMOVED***
		ExecutionSegment:         segment,
		ExecutionSegmentSequence: &sequence,
	***REMOVED***.Apply(runner.GetOptions()).Apply(extraOptions)

	testRunState := getTestRunState(t, options, runner)

	execReqs := config.GetExecutionRequirements(et)
	es := lib.NewExecutionState(testRunState, et, lib.GetMaxPlannedVUs(execReqs), lib.GetMaxPossibleVUs(execReqs))
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
