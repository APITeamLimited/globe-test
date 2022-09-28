package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/APITeamLimited/globe-test/worker/js"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils"
	"github.com/APITeamLimited/globe-test/worker/loader"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutionInfoVUSharing(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
		import exec from 'k6/execution';
		import ***REMOVED*** sleep ***REMOVED*** from 'k6';

		// The cvus scenario should reuse the two VUs created for the carr scenario.
		export let options = ***REMOVED***
			scenarios: ***REMOVED***
				carr: ***REMOVED***
					executor: 'constant-arrival-rate',
					exec: 'carr',
					rate: 9,
					timeUnit: '0.95s',
					duration: '1s',
					preAllocatedVUs: 2,
					maxVUs: 10,
					gracefulStop: '100ms',
				***REMOVED***,
			    cvus: ***REMOVED***
					executor: 'constant-vus',
					exec: 'cvus',
					vus: 2,
					duration: '1s',
					startTime: '2s',
					gracefulStop: '0s',
			    ***REMOVED***,
		    ***REMOVED***,
		***REMOVED***;

		export function cvus() ***REMOVED***
			const info = Object.assign(***REMOVED***scenario: 'cvus'***REMOVED***, exec.vu);
			console.log(JSON.stringify(info));
			sleep(0.2);
		***REMOVED***;

		export function carr() ***REMOVED***
			const info = Object.assign(***REMOVED***scenario: 'carr'***REMOVED***, exec.vu);
			console.log(JSON.stringify(info));
		***REMOVED***;
`)

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel***REMOVED******REMOVED***
	logger.AddHook(&logHook)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	runner, err := js.New(
		&libWorker.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
			Data: script,
		***REMOVED***,
		nil, libWorker.GetTestWorkerInfo(),
	)
	require.NoError(t, err)

	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, libWorker.Options***REMOVED******REMOVED***)
	defer cancel()

	type vuStat struct ***REMOVED***
		iteration uint64
		scIter    map[string]uint64
	***REMOVED***
	vuStats := map[uint64]*vuStat***REMOVED******REMOVED***

	type logEntry struct ***REMOVED***
		IDInInstance        uint64
		Scenario            string
		IterationInInstance uint64
		IterationInScenario uint64
	***REMOVED***

	errCh := make(chan error, 1)
	go func() ***REMOVED*** errCh <- execScheduler.Run(ctx, ctx, samples, libWorker.GetTestWorkerInfo()) ***REMOVED***()

	select ***REMOVED***
	case err := <-errCh:
		require.NoError(t, err)
		entries := logHook.Drain()
		assert.InDelta(t, 20, len(entries), 2)
		le := &logEntry***REMOVED******REMOVED***
		for _, entry := range entries ***REMOVED***
			err = json.Unmarshal([]byte(entry.Message), le)
			require.NoError(t, err)
			assert.Contains(t, []uint64***REMOVED***1, 2***REMOVED***, le.IDInInstance)
			if _, ok := vuStats[le.IDInInstance]; !ok ***REMOVED***
				vuStats[le.IDInInstance] = &vuStat***REMOVED***0, make(map[string]uint64)***REMOVED***
			***REMOVED***
			if le.IterationInInstance > vuStats[le.IDInInstance].iteration ***REMOVED***
				vuStats[le.IDInInstance].iteration = le.IterationInInstance
			***REMOVED***
			if le.IterationInScenario > vuStats[le.IDInInstance].scIter[le.Scenario] ***REMOVED***
				vuStats[le.IDInInstance].scIter[le.Scenario] = le.IterationInScenario
			***REMOVED***
		***REMOVED***
		require.Len(t, vuStats, 2)
		// Both VUs should complete 10 iterations each globally, but 5
		// iterations each per scenario (iterations are 0-based)
		for _, v := range vuStats ***REMOVED***
			assert.Equal(t, uint64(9), v.iteration)
			assert.Equal(t, uint64(4), v.scIter["cvus"])
			assert.Equal(t, uint64(4), v.scIter["carr"])
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timed out")
	***REMOVED***
***REMOVED***

func TestExecutionInfoScenarioIter(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
		import exec from 'k6/execution';

		// The pvu scenario should reuse the two VUs created for the carr scenario.
		export let options = ***REMOVED***
			scenarios: ***REMOVED***
				carr: ***REMOVED***
					executor: 'constant-arrival-rate',
					exec: 'carr',
					rate: 9,
					timeUnit: '0.95s',
					duration: '1s',
					preAllocatedVUs: 2,
					maxVUs: 10,
					gracefulStop: '100ms',
				***REMOVED***,
				pvu: ***REMOVED***
					executor: 'per-vu-iterations',
					exec: 'pvu',
					vus: 2,
					iterations: 5,
					startTime: '2s',
					gracefulStop: '100ms',
				***REMOVED***,
			***REMOVED***,
		***REMOVED***;

		export function pvu() ***REMOVED***
			const info = Object.assign(***REMOVED***VUID: __VU***REMOVED***, exec.scenario);
			console.log(JSON.stringify(info));
		***REMOVED***

		export function carr() ***REMOVED***
			const info = Object.assign(***REMOVED***VUID: __VU***REMOVED***, exec.scenario);
			console.log(JSON.stringify(info));
		***REMOVED***;
`)

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel***REMOVED******REMOVED***
	logger.AddHook(&logHook)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	runner, err := js.New(
		&libWorker.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
			Data: script,
		***REMOVED***,
		nil, libWorker.GetTestWorkerInfo(),
	)
	require.NoError(t, err)

	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, libWorker.Options***REMOVED******REMOVED***)
	defer cancel()

	errCh := make(chan error, 1)
	go func() ***REMOVED*** errCh <- execScheduler.Run(ctx, ctx, samples, libWorker.GetTestWorkerInfo()) ***REMOVED***()

	scStats := map[string]uint64***REMOVED******REMOVED***

	type logEntry struct ***REMOVED***
		Name                      string
		IterationInInstance, VUID uint64
	***REMOVED***

	select ***REMOVED***
	case err := <-errCh:
		require.NoError(t, err)
		entries := logHook.Drain()
		require.Len(t, entries, 20)
		le := &logEntry***REMOVED******REMOVED***
		for _, entry := range entries ***REMOVED***
			err = json.Unmarshal([]byte(entry.Message), le)
			require.NoError(t, err)
			assert.Contains(t, []uint64***REMOVED***1, 2***REMOVED***, le.VUID)
			if le.IterationInInstance > scStats[le.Name] ***REMOVED***
				scStats[le.Name] = le.IterationInInstance
			***REMOVED***
		***REMOVED***
		require.Len(t, scStats, 2)
		// The global per scenario iteration count should be 9 (iterations
		// start at 0), despite VUs being shared or more than 1 being used.
		for _, v := range scStats ***REMOVED***
			assert.Equal(t, uint64(9), v)
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("timed out")
	***REMOVED***
***REMOVED***

// Ensure that scenario iterations returned from k6/execution are
// stable during the execution of an iteration.
func TestSharedIterationsStable(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
		import ***REMOVED*** sleep ***REMOVED*** from 'k6';
		import exec from 'k6/execution';

		export let options = ***REMOVED***
			scenarios: ***REMOVED***
				test: ***REMOVED***
					executor: 'shared-iterations',
					vus: 50,
					iterations: 50,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***;
		export default function () ***REMOVED***
			sleep(1);
			console.log(JSON.stringify(Object.assign(***REMOVED***VUID: __VU***REMOVED***, exec.scenario)));
		***REMOVED***
`)

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel***REMOVED******REMOVED***
	logger.AddHook(&logHook)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	runner, err := js.New(
		&libWorker.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
			Data: script,
		***REMOVED***,
		nil, libWorker.GetTestWorkerInfo(),
	)
	require.NoError(t, err)

	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, libWorker.Options***REMOVED******REMOVED***)
	defer cancel()

	errCh := make(chan error, 1)
	go func() ***REMOVED*** errCh <- execScheduler.Run(ctx, ctx, samples, libWorker.GetTestWorkerInfo()) ***REMOVED***()

	expIters := [50]int64***REMOVED******REMOVED***
	for i := 0; i < 50; i++ ***REMOVED***
		expIters[i] = int64(i)
	***REMOVED***
	gotLocalIters, gotGlobalIters := []int64***REMOVED******REMOVED***, []int64***REMOVED******REMOVED***

	type logEntry struct***REMOVED*** IterationInInstance, IterationInTest int64 ***REMOVED***

	select ***REMOVED***
	case err := <-errCh:
		require.NoError(t, err)
		entries := logHook.Drain()
		require.Len(t, entries, 50)
		le := &logEntry***REMOVED******REMOVED***
		for _, entry := range entries ***REMOVED***
			err = json.Unmarshal([]byte(entry.Message), le)
			require.NoError(t, err)
			require.Equal(t, le.IterationInInstance, le.IterationInTest)
			gotLocalIters = append(gotLocalIters, le.IterationInInstance)
			gotGlobalIters = append(gotGlobalIters, le.IterationInTest)
		***REMOVED***

		assert.ElementsMatch(t, expIters, gotLocalIters)
		assert.ElementsMatch(t, expIters, gotGlobalIters)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out")
	***REMOVED***
***REMOVED***

func TestExecutionInfoAll(t *testing.T) ***REMOVED***
	t.Parallel()

	scriptTemplate := `
	import ***REMOVED*** sleep ***REMOVED*** from 'k6';
	import exec from "k6/execution";

	export let options = ***REMOVED***
		scenarios: ***REMOVED***
			executor: ***REMOVED***
				executor: "%[1]s",
				%[2]s
			***REMOVED***
		***REMOVED***
	***REMOVED***

	export default function () ***REMOVED***
		sleep(0.2);
		console.log(JSON.stringify(exec));
	***REMOVED***`

	executorConfigs := map[string]string***REMOVED***
		"constant-arrival-rate": `
			rate: 1,
			timeUnit: "1s",
			duration: "1s",
			preAllocatedVUs: 1,
			maxVUs: 2,
			gracefulStop: "0s",`,
		"constant-vus": `
			vus: 1,
			duration: "1s",
			gracefulStop: "0s",`,
		"externally-controlled": `
			vus: 1,
			duration: "1s",`,
		"per-vu-iterations": `
			vus: 1,
			iterations: 1,
			gracefulStop: "0s",`,
		"shared-iterations": `
			vus: 1,
			iterations: 1,
			gracefulStop: "0s",`,
		"ramping-arrival-rate": `
			startRate: 1,
			timeUnit: "0.5s",
			preAllocatedVUs: 1,
			maxVUs: 2,
			stages: [ ***REMOVED*** target: 1, duration: "1s" ***REMOVED*** ],
			gracefulStop: "0s",`,
		"ramping-vus": `
			startVUs: 1,
			stages: [ ***REMOVED*** target: 1, duration: "1s" ***REMOVED*** ],
			gracefulStop: "0s",`,
	***REMOVED***

	testCases := []struct***REMOVED*** name, script string ***REMOVED******REMOVED******REMOVED***

	for ename, econf := range executorConfigs ***REMOVED***
		testCases = append(testCases, struct***REMOVED*** name, script string ***REMOVED******REMOVED***
			ename, fmt.Sprintf(scriptTemplate, ename, econf),
		***REMOVED***)
	***REMOVED***

	// We're only checking a small subset of all properties, to ensure
	// there were no errors with accessing any of the top-level ones.
	// Most of the others are time-based, and would be difficult/flaky to check.
	type logEntry struct ***REMOVED***
		Scenario struct***REMOVED*** Executor string ***REMOVED***
		Instance struct***REMOVED*** VUsActive int ***REMOVED***
		VU       struct***REMOVED*** IDInTest int ***REMOVED***
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			logger := logrus.New()
			logger.SetOutput(ioutil.Discard)
			logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel***REMOVED******REMOVED***
			logger.AddHook(&logHook)

			registry := metrics.NewRegistry()
			builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
			runner, err := js.New(
				&libWorker.TestPreInitState***REMOVED***
					Logger:         logger,
					BuiltinMetrics: builtinMetrics,
					Registry:       registry,
				***REMOVED***,
				&loader.SourceData***REMOVED***
					URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
					Data: []byte(tc.script),
				***REMOVED***, nil, libWorker.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, libWorker.Options***REMOVED******REMOVED***)
			defer cancel()

			errCh := make(chan error, 1)
			go func() ***REMOVED*** errCh <- execScheduler.Run(ctx, ctx, samples, libWorker.GetTestWorkerInfo()) ***REMOVED***()

			select ***REMOVED***
			case err := <-errCh:
				require.NoError(t, err)
				entries := logHook.Drain()
				require.GreaterOrEqual(t, len(entries), 1)

				le := &logEntry***REMOVED******REMOVED***
				err = json.Unmarshal([]byte(entries[0].Message), le)
				require.NoError(t, err)

				assert.Equal(t, tc.name, le.Scenario.Executor)
				assert.Equal(t, 1, le.Instance.VUsActive)
				assert.Equal(t, 1, le.VU.IDInTest)
			case <-time.After(5 * time.Second):
				t.Fatal("timed out")
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
