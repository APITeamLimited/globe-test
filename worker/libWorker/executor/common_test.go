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

func simpleRunner(vuFn func(context.Context, *libWorker.State) error) libWorker.Runner {
	return &minirunner.MiniRunner{
		Fn: func(ctx context.Context, state *libWorker.State, _ chan<- workerMetrics.SampleContainer) error {
			return vuFn(ctx, state)
		},
	}
}

func getTestRunState(tb testing.TB, options libWorker.Options, runner libWorker.Runner) *libWorker.TestRunState {
	reg := workerMetrics.NewRegistry()
	piState := &libWorker.TestPreInitState{
		Logger:         testutils.NewLogger(tb),
		RuntimeOptions: libWorker.RuntimeOptions{},
		Registry:       reg,
		BuiltinMetrics: workerMetrics.RegisterBuiltinMetrics(reg),
	}

	require.NoError(tb, runner.SetOptions(options))

	return &libWorker.TestRunState{
		TestPreInitState: piState,
		Options:          options,
		Runner:           runner,
	}
}

func setupExecutor(t testing.TB, config libWorker.ExecutorConfig, es *libWorker.ExecutionState) (
	context.Context, context.CancelFunc, libWorker.Executor, *testutils.SimpleLogrusHook,
) {
	ctx, cancel := context.WithCancel(context.Background())
	engineOut := make(chan workerMetrics.SampleContainer, 100) // TODO: return this for more complicated tests?

	logHook := &testutils.SimpleLogrusHook{HookedLevels: []logrus.Level{logrus.WarnLevel}}
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)
	logEntry := logrus.NewEntry(testLog)

	initVUFunc := func(_ context.Context, logger *logrus.Entry, workerInfo *libWorker.WorkerInfo) (libWorker.InitializedVU, error) {
		idl, idg := es.GetUniqueVUIdentifiers()
		return es.Test.Runner.NewVU(idl, idg, engineOut, libWorker.GetTestWorkerInfo())
	}
	es.SetInitVUFunc(initVUFunc)

	maxPlannedVUs := libWorker.GetMaxPlannedVUs(config.GetExecutionRequirements(es.ExecutionTuple))
	initializeVUs(ctx, t, logEntry, es, maxPlannedVUs, initVUFunc)

	executor, err := config.NewExecutor(es, logEntry)
	require.NoError(t, err)

	err = executor.Init(ctx)
	require.NoError(t, err)
	return ctx, cancel, executor, logHook
}

func initializeVUs(
	ctx context.Context, t testing.TB, logEntry *logrus.Entry, es *libWorker.ExecutionState, number uint64, initVU libWorker.InitVUFunc,
) {
	// This is not how the local ExecutionScheduler initializes VUs, but should do the same job
	for i := uint64(0); i < number; i++ {
		// Not calling es.InitializeNewVU() here to avoid a double increment of initializedVUs,
		// which is done in es.AddInitializedVU().
		vu, err := initVU(ctx, logEntry, libWorker.GetTestWorkerInfo())
		require.NoError(t, err)
		es.AddInitializedVU(vu)
	}
}

type executorTest struct {
	options libWorker.Options
	state   *libWorker.ExecutionState

	ctx      context.Context //nolint
	cancel   context.CancelFunc
	executor libWorker.Executor
	logHook  *testutils.SimpleLogrusHook
}

func setupExecutorTest(
	t testing.TB, segmentStr, sequenceStr string, extraOptions libWorker.Options,
	runner libWorker.Runner, config libWorker.ExecutorConfig,
) *executorTest {
	var err error
	var segment *libWorker.ExecutionSegment
	if segmentStr != "" {
		segment, err = libWorker.NewExecutionSegmentFromString(segmentStr)
		require.NoError(t, err)
	}

	var sequence libWorker.ExecutionSegmentSequence
	if sequenceStr != "" {
		sequence, err = libWorker.NewExecutionSegmentSequenceFromString(sequenceStr)
		require.NoError(t, err)
	}

	et, err := libWorker.NewExecutionTuple(segment, &sequence)
	require.NoError(t, err)

	options := libWorker.Options{
		ExecutionSegment:         segment,
		ExecutionSegmentSequence: &sequence,
	}.Apply(runner.GetOptions()).Apply(extraOptions)

	testRunState := getTestRunState(t, options, runner)

	execReqs := config.GetExecutionRequirements(et)
	es := libWorker.NewExecutionState(testRunState, et, libWorker.GetMaxPlannedVUs(execReqs), libWorker.GetMaxPossibleVUs(execReqs))
	ctx, cancel, executor, logHook := setupExecutor(t, config, es)

	return &executorTest{
		options:  options,
		state:    es,
		ctx:      ctx,
		cancel:   cancel,
		executor: executor,
		logHook:  logHook,
	}
}
