package tests

import (
	"context"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/cmd/integration_tests/testmodules/events"
	"go.k6.io/k6/core/local"
	"go.k6.io/k6/js"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
	"gopkg.in/guregu/null.v3"
)

func eventLoopTest(t *testing.T, script []byte, testHandle func(context.Context, lib.Runner, error, *testutils.SimpleLogrusHook)) ***REMOVED***
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel***REMOVED******REMOVED***
	logger.AddHook(logHook)

	registry := metrics.NewRegistry()
	piState := &lib.TestPreInitState***REMOVED***
		Logger:         logger,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
	***REMOVED***

	script = []byte("import ***REMOVED***setTimeout***REMOVED*** from 'k6/x/events';\n" + string(script))
	runner, err := js.New(piState, &loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: script***REMOVED***, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newOpts, err := executor.DeriveScenariosFromShortcuts(lib.Options***REMOVED***
		MetricSamplesBufferSize: null.NewInt(200, false),
		TeardownTimeout:         types.NullDurationFrom(time.Second),
		SetupTimeout:            types.NullDurationFrom(time.Second),
	***REMOVED***.Apply(runner.GetOptions()), nil)
	require.NoError(t, err)
	require.Empty(t, newOpts.Validate())
	require.NoError(t, runner.SetOptions(newOpts))

	testState := &lib.TestRunState***REMOVED***
		TestPreInitState: piState,
		Options:          newOpts,
		Runner:           runner,
	***REMOVED***

	execScheduler, err := local.NewExecutionScheduler(testState)
	require.NoError(t, err)

	samples := make(chan metrics.SampleContainer, newOpts.MetricSamplesBufferSize.Int64)
	go func() ***REMOVED***
		for ***REMOVED***
			select ***REMOVED***
			case <-samples:
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	require.NoError(t, execScheduler.Init(ctx, samples))

	errCh := make(chan error, 1)
	go func() ***REMOVED*** errCh <- execScheduler.Run(ctx, ctx, samples) ***REMOVED***()

	select ***REMOVED***
	case err := <-errCh:
		testHandle(ctx, runner, err, logHook)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out")
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	modules.Register("k6/x/events", events.New())
***REMOVED***

func TestEventLoop(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
		setTimeout(()=> ***REMOVED***console.log("initcontext setTimeout")***REMOVED***, 200)
		console.log("initcontext");
		export default function() ***REMOVED***
			setTimeout(()=> ***REMOVED***console.log("default setTimeout")***REMOVED***, 200)
			console.log("default");
		***REMOVED***;
		export function setup() ***REMOVED***
			setTimeout(()=> ***REMOVED***console.log("setup setTimeout")***REMOVED***, 200)
			console.log("setup");
		***REMOVED***;
		export function teardown() ***REMOVED***
			setTimeout(()=> ***REMOVED***console.log("teardown setTimeout")***REMOVED***, 200)
			console.log("teardown");
		***REMOVED***;
		export function handleSummary() ***REMOVED***
			setTimeout(()=> ***REMOVED***console.log("handleSummary setTimeout")***REMOVED***, 200)
			console.log("handleSummary");
		***REMOVED***;
`)
	eventLoopTest(t, script, func(ctx context.Context, runner lib.Runner, err error, logHook *testutils.SimpleLogrusHook) ***REMOVED***
		require.NoError(t, err)
		_, err = runner.HandleSummary(ctx, &lib.Summary***REMOVED***RootGroup: &lib.Group***REMOVED******REMOVED******REMOVED***)
		require.NoError(t, err)
		entries := logHook.Drain()
		msgs := make([]string, len(entries))
		for i, entry := range entries ***REMOVED***
			msgs[i] = entry.Message
		***REMOVED***
		require.Equal(t, []string***REMOVED***
			"initcontext", // first initialization
			"initcontext setTimeout",
			"initcontext", // for vu
			"initcontext setTimeout",
			"initcontext", // for setup
			"initcontext setTimeout",
			"setup", // setup
			"setup setTimeout",
			"default", // one iteration
			"default setTimeout",
			"initcontext", // for teardown
			"initcontext setTimeout",
			"teardown", // teardown
			"teardown setTimeout",
			"initcontext", // for handleSummary
			"initcontext setTimeout",
			"handleSummary", // handleSummary
			"handleSummary setTimeout",
		***REMOVED***, msgs)
	***REMOVED***)
***REMOVED***

func TestEventLoopCrossScenario(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
import exec from "k6/execution"
export const options = ***REMOVED***
        scenarios: ***REMOVED***
                "first":***REMOVED***
                        executor: "shared-iterations",
                        maxDuration: "1s",
                        iterations: 1,
                        vus: 1,
                        gracefulStop:"1s",
                ***REMOVED***,
                "second": ***REMOVED***
                        executor: "shared-iterations",
                        maxDuration: "1s",
                        iterations: 1,
                        vus: 1,
                        startTime: "3s",
                ***REMOVED***
        ***REMOVED***
***REMOVED***
export default function() ***REMOVED***
	let i = exec.scenario.name
	setTimeout(()=> ***REMOVED***console.log(i)***REMOVED***, 3000)
***REMOVED***
`)

	eventLoopTest(t, script, func(_ context.Context, _ lib.Runner, err error, logHook *testutils.SimpleLogrusHook) ***REMOVED***
		require.NoError(t, err)
		entries := logHook.Drain()
		msgs := make([]string, len(entries))
		for i, entry := range entries ***REMOVED***
			msgs[i] = entry.Message
		***REMOVED***
		require.Equal(t, []string***REMOVED***
			"setTimeout 1 was stopped because the VU iteration was interrupted",
			"second",
		***REMOVED***, msgs)
	***REMOVED***)
***REMOVED***

func TestEventLoopDoesntCrossIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
import ***REMOVED*** sleep ***REMOVED*** from "k6"
export const options = ***REMOVED***
  iterations: 2,
  vus: 1,
***REMOVED***

export default function() ***REMOVED***
  let i = __ITER;
	setTimeout(()=> ***REMOVED*** console.log(i) ***REMOVED***, 1000)
  if (__ITER == 0) ***REMOVED***
    throw "just error"
  ***REMOVED*** else ***REMOVED***
    sleep(1)
  ***REMOVED***
***REMOVED***
`)

	eventLoopTest(t, script, func(_ context.Context, _ lib.Runner, err error, logHook *testutils.SimpleLogrusHook) ***REMOVED***
		require.NoError(t, err)
		entries := logHook.Drain()
		msgs := make([]string, len(entries))
		for i, entry := range entries ***REMOVED***
			msgs[i] = entry.Message
		***REMOVED***
		require.Equal(t, []string***REMOVED***
			"setTimeout 1 was stopped because the VU iteration was interrupted",
			"just error\n\tat /script.js:13:4(15)\n\tat native\n", "1",
		***REMOVED***, msgs)
	***REMOVED***)
***REMOVED***
