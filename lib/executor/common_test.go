package executor

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/testutils/minirunner"
	"github.com/loadimpact/k6/stats"
)

func simpleRunner(vuFn func(context.Context) error) lib.Runner ***REMOVED***
	return &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, _ chan<- stats.SampleContainer) error ***REMOVED***
			return vuFn(ctx)
		***REMOVED***,
	***REMOVED***
***REMOVED***

func setupExecutor(t *testing.T, config lib.ExecutorConfig, es *lib.ExecutionState, runner lib.Runner) (
	context.Context, context.CancelFunc, lib.Executor, *testutils.SimpleLogrusHook,
) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	engineOut := make(chan stats.SampleContainer, 100) // TODO: return this for more complicated tests?

	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	testLog := logrus.New()
	testLog.AddHook(logHook)
	testLog.SetOutput(ioutil.Discard)
	logEntry := logrus.NewEntry(testLog)

	es.SetInitVUFunc(func(_ context.Context, logger *logrus.Entry) (lib.VU, error) ***REMOVED***
		return runner.NewVU(engineOut)
	***REMOVED***)

	segment := es.Options.ExecutionSegment
	maxVUs := lib.GetMaxPossibleVUs(config.GetExecutionRequirements(segment))
	initializeVUs(ctx, t, logEntry, es, maxVUs)

	executor, err := config.NewExecutor(es, logEntry)
	require.NoError(t, err)

	err = executor.Init(ctx)
	require.NoError(t, err)
	return ctx, cancel, executor, logHook
***REMOVED***

func initializeVUs(
	ctx context.Context, t testing.TB, logEntry *logrus.Entry, es *lib.ExecutionState, number uint64,
) ***REMOVED***
	// This is not how the local ExecutionScheduler initializes VUs, but should do the same job
	for i := uint64(0); i < number; i++ ***REMOVED***
		vu, err := es.InitializeNewVU(ctx, logEntry)
		require.NoError(t, err)
		es.AddInitializedVU(vu)
	***REMOVED***
***REMOVED***
