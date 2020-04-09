/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package core

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/executor"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/lib/testutils/minirunner"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/dummy"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

const isWindows = runtime.GOOS == "windows"

// Wrapper around NewEngine that applies a logger and manages the options.
func newTestEngine( //nolint:golint
	t *testing.T, runCtx context.Context, runner lib.Runner, opts lib.Options,
) (engine *Engine, run func() error, wait func()) ***REMOVED***
	if runner == nil ***REMOVED***
		runner = &minirunner.MiniRunner***REMOVED******REMOVED***
	***REMOVED***
	globalCtx, globalCancel := context.WithCancel(context.Background())
	var runCancel func()
	if runCtx == nil ***REMOVED***
		runCtx, runCancel = context.WithCancel(globalCtx)
	***REMOVED***

	newOpts, err := executor.DeriveExecutionFromShortcuts(lib.Options***REMOVED***
		MetricSamplesBufferSize: null.NewInt(200, false),
	***REMOVED***.Apply(runner.GetOptions()).Apply(opts))
	require.NoError(t, err)
	require.Empty(t, newOpts.Validate())

	require.NoError(t, runner.SetOptions(newOpts))

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	execScheduler, err := local.NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	engine, err = NewEngine(execScheduler, opts, logger)
	require.NoError(t, err)

	run, waitFn, err := engine.Init(globalCtx, runCtx)
	require.NoError(t, err)

	return engine, run, func() ***REMOVED***
		if runCancel != nil ***REMOVED***
			runCancel()
		***REMOVED***
		globalCancel()
		waitFn()
	***REMOVED***
***REMOVED***

func TestNewEngine(t *testing.T) ***REMOVED***
	newTestEngine(t, nil, nil, lib.Options***REMOVED******REMOVED***)
***REMOVED***

func TestEngineRun(t *testing.T) ***REMOVED***
	logrus.SetLevel(logrus.DebugLevel)
	t.Run("exits with context", func(t *testing.T) ***REMOVED***
		done := make(chan struct***REMOVED******REMOVED***)
		runner := &minirunner.MiniRunner***REMOVED***Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			<-ctx.Done()
			close(done)
			return nil
		***REMOVED******REMOVED***

		duration := 100 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		_, run, wait := newTestEngine(t, ctx, runner, lib.Options***REMOVED******REMOVED***)
		defer wait()

		startTime := time.Now()
		assert.NoError(t, run())
		assert.WithinDuration(t, startTime.Add(duration), time.Now(), 100*time.Millisecond)
		<-done
	***REMOVED***)
	t.Run("exits with executor", func(t *testing.T) ***REMOVED***
		e, run, wait := newTestEngine(t, nil, nil, lib.Options***REMOVED***
			VUs:        null.IntFrom(10),
			Iterations: null.IntFrom(100),
		***REMOVED***)
		defer wait()
		assert.NoError(t, run())
		assert.Equal(t, uint64(100), e.ExecutionScheduler.GetState().GetFullIterationCount())
	***REMOVED***)
	// Make sure samples are discarded after context close (using "cutoff" timestamp in local.go)
	t.Run("collects samples", func(t *testing.T) ***REMOVED***
		testMetric := stats.New("test_metric", stats.Trend)

		signalChan := make(chan interface***REMOVED******REMOVED***)

		runner := &minirunner.MiniRunner***REMOVED***Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			stats.PushIfNotDone(ctx, out, stats.Sample***REMOVED***Metric: testMetric, Time: time.Now(), Value: 1***REMOVED***)
			close(signalChan)
			<-ctx.Done()
			stats.PushIfNotDone(ctx, out, stats.Sample***REMOVED***Metric: testMetric, Time: time.Now(), Value: 1***REMOVED***)
			return nil
		***REMOVED******REMOVED***

		ctx, cancel := context.WithCancel(context.Background())
		e, run, wait := newTestEngine(t, ctx, runner, lib.Options***REMOVED***
			VUs:        null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)

		c := &dummy.Collector***REMOVED******REMOVED***
		e.Collectors = []lib.Collector***REMOVED***c***REMOVED***

		errC := make(chan error)
		go func() ***REMOVED*** errC <- run() ***REMOVED***()
		<-signalChan
		cancel()
		assert.NoError(t, <-errC)
		wait()

		found := 0
		for _, s := range c.Samples ***REMOVED***
			if s.Metric != testMetric ***REMOVED***
				continue
			***REMOVED***
			found++
			assert.Equal(t, 1.0, s.Value, "wrong value")
		***REMOVED***
		assert.Equal(t, 1, found, "wrong number of samples")
	***REMOVED***)
***REMOVED***

func TestEngineAtTime(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, run, wait := newTestEngine(t, ctx, nil, lib.Options***REMOVED***
		VUs:      null.IntFrom(2),
		Duration: types.NullDurationFrom(20 * time.Second),
	***REMOVED***)
	defer wait()

	assert.NoError(t, run())
***REMOVED***

func TestEngineStopped(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	e, run, wait := newTestEngine(t, ctx, nil, lib.Options***REMOVED***
		VUs:      null.IntFrom(1),
		Duration: types.NullDurationFrom(20 * time.Second),
	***REMOVED***)
	defer wait()

	assert.NoError(t, run())
	assert.Equal(t, false, e.IsStopped(), "engine should be running")
	e.Stop()
	assert.Equal(t, true, e.IsStopped(), "engine should be stopped")
***REMOVED***

func TestEngineCollector(t *testing.T) ***REMOVED***
	testMetric := stats.New("test_metric", stats.Trend)

	runner := &minirunner.MiniRunner***REMOVED***Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
		out <- stats.Sample***REMOVED***Metric: testMetric***REMOVED***
		return nil
	***REMOVED******REMOVED***

	e, run, wait := newTestEngine(t, nil, runner, lib.Options***REMOVED***VUs: null.IntFrom(1), Iterations: null.IntFrom(1)***REMOVED***)

	c := &dummy.Collector***REMOVED******REMOVED***
	e.Collectors = []lib.Collector***REMOVED***c***REMOVED***

	assert.NoError(t, run())
	wait()

	cSamples := []stats.Sample***REMOVED******REMOVED***
	for _, sample := range c.Samples ***REMOVED***
		if sample.Metric == testMetric ***REMOVED***
			cSamples = append(cSamples, sample)
		***REMOVED***
	***REMOVED***
	metric := e.Metrics["test_metric"]
	if assert.NotNil(t, metric) ***REMOVED***
		sink := metric.Sink.(*stats.TrendSink)
		if assert.NotNil(t, sink) ***REMOVED***
			numCollectorSamples := len(cSamples)
			numEngineSamples := len(sink.Values)
			assert.Equal(t, numEngineSamples, numCollectorSamples)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEngine_processSamples(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)

	t.Run("metric", func(t *testing.T) ***REMOVED***
		e, _, wait := newTestEngine(t, nil, nil, lib.Options***REMOVED******REMOVED***)
		defer wait()

		e.processSamples(
			[]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED******REMOVED***,
		)

		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric"].Sink)
	***REMOVED***)
	t.Run("submetric", func(t *testing.T) ***REMOVED***
		ths, err := stats.NewThresholds([]string***REMOVED***`1+1==2`***REMOVED***)
		assert.NoError(t, err)

		e, _, wait := newTestEngine(t, nil, nil, lib.Options***REMOVED***
			Thresholds: map[string]stats.Thresholds***REMOVED***
				"my_metric***REMOVED***a:1***REMOVED***": ths,
			***REMOVED***,
		***REMOVED***)
		defer wait()

		sms := e.submetrics["my_metric"]
		assert.Len(t, sms, 1)
		assert.Equal(t, "my_metric***REMOVED***a:1***REMOVED***", sms[0].Name)
		assert.EqualValues(t, map[string]string***REMOVED***"a": "1"***REMOVED***, sms[0].Tags.CloneTags())

		e.processSamples(
			[]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED***)***REMOVED******REMOVED***,
		)

		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric"].Sink)
		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric***REMOVED***a:1***REMOVED***"].Sink)
	***REMOVED***)
***REMOVED***

func TestEngineThresholdsWillAbort(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)

	ths, err := stats.NewThresholds([]string***REMOVED***"1+1==3"***REMOVED***)
	assert.NoError(t, err)
	ths.Thresholds[0].AbortOnFail = true

	thresholds := map[string]stats.Thresholds***REMOVED***metric.Name: ths***REMOVED***

	e, _, wait := newTestEngine(t, nil, nil, lib.Options***REMOVED***Thresholds: thresholds***REMOVED***)
	defer wait()

	e.processSamples(
		[]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED******REMOVED***,
	)
	assert.True(t, e.processThresholds())
***REMOVED***

func TestEngineAbortedByThresholds(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)

	ths, err := stats.NewThresholds([]string***REMOVED***"1+1==3"***REMOVED***)
	assert.NoError(t, err)
	ths.Thresholds[0].AbortOnFail = true

	thresholds := map[string]stats.Thresholds***REMOVED***metric.Name: ths***REMOVED***

	done := make(chan struct***REMOVED******REMOVED***)
	runner := &minirunner.MiniRunner***REMOVED***Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
		out <- stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED***
		<-ctx.Done()
		close(done)
		return nil
	***REMOVED******REMOVED***

	_, run, wait := newTestEngine(t, nil, runner, lib.Options***REMOVED***Thresholds: thresholds***REMOVED***)
	defer wait()

	go func() ***REMOVED***
		assert.NoError(t, run())
	***REMOVED***()

	select ***REMOVED***
	case <-done:
		return
	case <-time.After(10 * time.Second):
		assert.Fail(t, "Test should have completed within 10 seconds")
	***REMOVED***
***REMOVED***

func TestEngine_processThresholds(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)

	testdata := map[string]struct ***REMOVED***
		pass  bool
		ths   map[string][]string
		abort bool
	***REMOVED******REMOVED***
		"passing":  ***REMOVED***true, map[string][]string***REMOVED***"my_metric": ***REMOVED***"1+1==2"***REMOVED******REMOVED***, false***REMOVED***,
		"failing":  ***REMOVED***false, map[string][]string***REMOVED***"my_metric": ***REMOVED***"1+1==3"***REMOVED******REMOVED***, false***REMOVED***,
		"aborting": ***REMOVED***false, map[string][]string***REMOVED***"my_metric": ***REMOVED***"1+1==3"***REMOVED******REMOVED***, true***REMOVED***,

		"submetric,match,passing":   ***REMOVED***true, map[string][]string***REMOVED***"my_metric***REMOVED***a:1***REMOVED***": ***REMOVED***"1+1==2"***REMOVED******REMOVED***, false***REMOVED***,
		"submetric,match,failing":   ***REMOVED***false, map[string][]string***REMOVED***"my_metric***REMOVED***a:1***REMOVED***": ***REMOVED***"1+1==3"***REMOVED******REMOVED***, false***REMOVED***,
		"submetric,nomatch,passing": ***REMOVED***true, map[string][]string***REMOVED***"my_metric***REMOVED***a:2***REMOVED***": ***REMOVED***"1+1==2"***REMOVED******REMOVED***, false***REMOVED***,
		"submetric,nomatch,failing": ***REMOVED***true, map[string][]string***REMOVED***"my_metric***REMOVED***a:2***REMOVED***": ***REMOVED***"1+1==3"***REMOVED******REMOVED***, false***REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		name, data := name, data
		t.Run(name, func(t *testing.T) ***REMOVED***
			thresholds := make(map[string]stats.Thresholds, len(data.ths))
			for m, srcs := range data.ths ***REMOVED***
				ths, err := stats.NewThresholds(srcs)
				assert.NoError(t, err)
				ths.Thresholds[0].AbortOnFail = data.abort
				thresholds[m] = ths
			***REMOVED***

			e, _, wait := newTestEngine(t, nil, nil, lib.Options***REMOVED***Thresholds: thresholds***REMOVED***)
			defer wait()

			e.processSamples(
				[]stats.SampleContainer***REMOVED***stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED******REMOVED***,
			)

			assert.Equal(t, data.abort, e.processThresholds())
			assert.Equal(t, data.pass, !e.IsTainted())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func getMetricSum(collector *dummy.Collector, name string) (result float64) ***REMOVED***
	for _, sc := range collector.SampleContainers ***REMOVED***
		for _, s := range sc.GetSamples() ***REMOVED***
			if s.Metric.Name == name ***REMOVED***
				result += s.Value
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func getMetricCount(collector *dummy.Collector, name string) (result uint) ***REMOVED***
	for _, sc := range collector.SampleContainers ***REMOVED***
		for _, s := range sc.GetSamples() ***REMOVED***
			if s.Metric.Name == name ***REMOVED***
				result++
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

const expectedHeaderMaxLength = 500

// FIXME: This test is too brittle, consider simplifying.
func TestSentReceivedMetrics(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	tr := tb.Replacer.Replace

	type testScript struct ***REMOVED***
		Code                 string
		NumRequests          int64
		ExpectedDataSent     int64
		ExpectedDataReceived int64
	***REMOVED***
	testScripts := []testScript***REMOVED***
		***REMOVED***tr(`import http from "k6/http";
			export default function() ***REMOVED***
				http.get("HTTPBIN_URL/bytes/15000");
			***REMOVED***`), 1, 0, 15000***REMOVED***,
		***REMOVED***tr(`import http from "k6/http";
			export default function() ***REMOVED***
				http.get("HTTPBIN_URL/bytes/5000");
				http.get("HTTPSBIN_URL/bytes/5000");
				http.batch(["HTTPBIN_URL/bytes/10000", "HTTPBIN_URL/bytes/20000", "HTTPSBIN_URL/bytes/10000"]);
			***REMOVED***`), 5, 0, 50000***REMOVED***,
		***REMOVED***tr(`import http from "k6/http";
			let data = "0123456789".repeat(100);
			export default function() ***REMOVED***
				http.post("HTTPBIN_URL/ip", ***REMOVED***
					file: http.file(data, "test.txt")
				***REMOVED***);
			***REMOVED***`), 1, 1000, 100***REMOVED***,
		// NOTE(imiric): This needs to keep testing against /ws-echo-invalid because
		// this test is highly sensitive to metric data, and slightly differing
		// WS server implementations might introduce flakiness.
		// See https://github.com/loadimpact/k6/pull/1149
		***REMOVED***tr(`import ws from "k6/ws";
			let data = "0123456789".repeat(100);
			export default function() ***REMOVED***
				ws.connect("WSBIN_URL/ws-echo-invalid", null, function (socket) ***REMOVED***
					socket.on('open', function open() ***REMOVED***
						socket.send(data);
					***REMOVED***);
					socket.on('message', function (message) ***REMOVED***
						socket.close();
					***REMOVED***);
				***REMOVED***);
			***REMOVED***`), 2, 1000, 1000***REMOVED***,
	***REMOVED***

	type testCase struct***REMOVED*** Iterations, VUs int64 ***REMOVED***
	testCases := []testCase***REMOVED***
		***REMOVED***1, 1***REMOVED***, ***REMOVED***2, 2***REMOVED***, ***REMOVED***2, 1***REMOVED***, ***REMOVED***5, 2***REMOVED***, ***REMOVED***25, 2***REMOVED***, ***REMOVED***50, 5***REMOVED***,
	***REMOVED***

	runTest := func(t *testing.T, ts testScript, tc testCase, noConnReuse bool) (float64, float64) ***REMOVED***
		r, err := js.New(
			&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: []byte(ts.Code)***REMOVED***,
			nil,
			lib.RuntimeOptions***REMOVED******REMOVED***,
		)
		require.NoError(t, err)

		engine, run, wait := newTestEngine(t, nil, r, lib.Options***REMOVED***
			Iterations:            null.IntFrom(tc.Iterations),
			VUs:                   null.IntFrom(tc.VUs),
			Hosts:                 tb.Dialer.Hosts,
			InsecureSkipTLSVerify: null.BoolFrom(true),
			NoVUConnectionReuse:   null.BoolFrom(noConnReuse),
			Batch:                 null.IntFrom(20),
		***REMOVED***)

		collector := &dummy.Collector***REMOVED******REMOVED***
		engine.Collectors = []lib.Collector***REMOVED***collector***REMOVED***

		errC := make(chan error)
		go func() ***REMOVED*** errC <- run() ***REMOVED***()

		select ***REMOVED***
		case <-time.After(10 * time.Second):
			t.Fatal("Test timed out")
		case err := <-errC:
			require.NoError(t, err)
		***REMOVED***
		wait()

		checkData := func(name string, expected int64) float64 ***REMOVED***
			data := getMetricSum(collector, name)
			expectedDataMin := float64(expected * tc.Iterations)
			expectedDataMax := float64((expected + ts.NumRequests*expectedHeaderMaxLength) * tc.Iterations)

			if data < expectedDataMin || data > expectedDataMax ***REMOVED***
				t.Errorf(
					"The %s sum should be in the interval [%f, %f] but was %f",
					name, expectedDataMin, expectedDataMax, data,
				)
			***REMOVED***
			return data
		***REMOVED***

		return checkData(metrics.DataSent.Name, ts.ExpectedDataSent),
			checkData(metrics.DataReceived.Name, ts.ExpectedDataReceived)
	***REMOVED***

	getTestCase := func(t *testing.T, ts testScript, tc testCase) func(t *testing.T) ***REMOVED***
		return func(t *testing.T) ***REMOVED***
			t.Parallel()
			noReuseSent, noReuseReceived := runTest(t, ts, tc, true)
			reuseSent, reuseReceived := runTest(t, ts, tc, false)

			if noReuseSent < reuseSent ***REMOVED***
				t.Errorf("reuseSent=%f is greater than noReuseSent=%f", reuseSent, noReuseSent)
			***REMOVED***
			if noReuseReceived < reuseReceived ***REMOVED***
				t.Errorf("reuseReceived=%f is greater than noReuseReceived=%f", reuseReceived, noReuseReceived)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// This Run will not return until the parallel subtests complete.
	t.Run("group", func(t *testing.T) ***REMOVED***
		for tsNum, ts := range testScripts ***REMOVED***
			for tcNum, tc := range testCases ***REMOVED***
				t.Run(
					fmt.Sprintf("SentReceivedMetrics_script[%d]_case[%d](%d,%d)", tsNum, tcNum, tc.Iterations, tc.VUs),
					getTestCase(t, ts, tc),
				)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestRunTags(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	runTagsMap := map[string]string***REMOVED***"foo": "bar", "test": "mest", "over": "written"***REMOVED***
	runTags := stats.NewSampleTags(runTagsMap)

	script := []byte(tb.Replacer.Replace(`
		import http from "k6/http";
		import ws from "k6/ws";
		import ***REMOVED*** Counter ***REMOVED*** from "k6/metrics";
		import ***REMOVED*** group, check, fail ***REMOVED*** from "k6";

		let customTags =  ***REMOVED*** "over": "the rainbow" ***REMOVED***;
		let params = ***REMOVED*** "tags": customTags***REMOVED***;
		let statusCheck = ***REMOVED*** "status is 200": (r) => r.status === 200 ***REMOVED***

		let myCounter = new Counter("mycounter");

		export default function() ***REMOVED***

			group("http", function() ***REMOVED***
				check(http.get("HTTPSBIN_URL", params), statusCheck, customTags);
				check(http.get("HTTPBIN_URL/status/418", params), statusCheck, customTags);
			***REMOVED***)

			group("websockets", function() ***REMOVED***
				var response = ws.connect("WSBIN_URL/ws-echo", params, function (socket) ***REMOVED***
					socket.on('open', function open() ***REMOVED***
						console.log('ws open and say hello');
						socket.send("hello");
					***REMOVED***);

					socket.on('message', function (message) ***REMOVED***
						console.log('ws got message ' + message);
						if (message != "hello") ***REMOVED***
							fail("Expected to receive 'hello' but got '" + message + "' instead !");
						***REMOVED***
						console.log('ws closing socket...');
						socket.close();
					***REMOVED***);

					socket.on('close', function () ***REMOVED***
						console.log('ws close');
					***REMOVED***);

					socket.on('error', function (e) ***REMOVED***
						console.log('ws error: ' + e.error());
					***REMOVED***);
				***REMOVED***);
				console.log('connect returned');
				check(response, ***REMOVED*** "status is 101": (r) => r && r.status === 101 ***REMOVED***, customTags);
			***REMOVED***)

			myCounter.add(1, customTags);
		***REMOVED***
	`))

	r, err := js.New(
		&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: script***REMOVED***,
		nil,
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)
	require.NoError(t, err)

	engine, run, wait := newTestEngine(t, nil, r, lib.Options***REMOVED***
		Iterations:            null.IntFrom(3),
		VUs:                   null.IntFrom(2),
		Hosts:                 tb.Dialer.Hosts,
		RunTags:               runTags,
		SystemTags:            &stats.DefaultSystemTagSet,
		InsecureSkipTLSVerify: null.BoolFrom(true),
	***REMOVED***)

	collector := &dummy.Collector***REMOVED******REMOVED***
	engine.Collectors = []lib.Collector***REMOVED***collector***REMOVED***

	errC := make(chan error)
	go func() ***REMOVED*** errC <- run() ***REMOVED***()

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out")
	case err := <-errC:
		require.NoError(t, err)
	***REMOVED***
	wait()

	systemMetrics := []*stats.Metric***REMOVED***
		metrics.VUs, metrics.VUsMax, metrics.Iterations, metrics.IterationDuration,
		metrics.GroupDuration, metrics.DataSent, metrics.DataReceived,
	***REMOVED***

	getExpectedOverVal := func(metricName string) string ***REMOVED***
		for _, sysMetric := range systemMetrics ***REMOVED***
			if sysMetric.Name == metricName ***REMOVED***
				return runTagsMap["over"]
			***REMOVED***
		***REMOVED***
		return "the rainbow"
	***REMOVED***

	for _, s := range collector.Samples ***REMOVED***
		for key, expVal := range runTagsMap ***REMOVED***
			val, ok := s.Tags.Get(key)

			if key == "over" ***REMOVED***
				expVal = getExpectedOverVal(s.Metric.Name)
			***REMOVED***

			assert.True(t, ok)
			assert.Equalf(t, expVal, val, "Wrong tag value in sample for metric %#v", s.Metric)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSetupTeardownThresholds(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	script := []byte(tb.Replacer.Replace(`
		import http from "k6/http";
		import ***REMOVED*** check ***REMOVED*** from "k6";
		import ***REMOVED*** Counter ***REMOVED*** from "k6/metrics";

		let statusCheck = ***REMOVED*** "status is 200": (r) => r.status === 200 ***REMOVED***
		let myCounter = new Counter("setup_teardown");

		export let options = ***REMOVED***
			iterations: 5,
			thresholds: ***REMOVED***
				"setup_teardown": ["count == 2"],
				"iterations": ["count == 5"],
				"http_reqs": ["count == 7"],
			***REMOVED***,
		***REMOVED***;

		export function setup() ***REMOVED***
			check(http.get("HTTPBIN_IP_URL"), statusCheck) && myCounter.add(1);
		***REMOVED***;

		export default function () ***REMOVED***
			check(http.get("HTTPBIN_IP_URL"), statusCheck);
		***REMOVED***;

		export function teardown() ***REMOVED***
			check(http.get("HTTPBIN_IP_URL"), statusCheck) && myCounter.add(1);
		***REMOVED***;
	`))

	runner, err := js.New(
		&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: script***REMOVED***,
		nil,
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)
	require.NoError(t, err)

	engine, run, wait := newTestEngine(t, nil, runner, lib.Options***REMOVED***
		SystemTags:      &stats.DefaultSystemTagSet,
		SetupTimeout:    types.NullDurationFrom(3 * time.Second),
		TeardownTimeout: types.NullDurationFrom(3 * time.Second),
		VUs:             null.IntFrom(3),
	***REMOVED***)
	defer wait()

	errC := make(chan error)
	go func() ***REMOVED*** errC <- run() ***REMOVED***()

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out")
	case err := <-errC:
		require.NoError(t, err)
		require.False(t, engine.IsTainted())
	***REMOVED***
***REMOVED***

func TestEmittedMetricsWhenScalingDown(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	script := []byte(tb.Replacer.Replace(`
		import http from "k6/http";
		import ***REMOVED*** sleep ***REMOVED*** from "k6";

		export let options = ***REMOVED***
			systemTags: ["iter", "vu", "url"],
			execution: ***REMOVED***
				we_need_hard_stop_and_ramp_down: ***REMOVED***
					type: "variable-looping-vus",
					// Start with 2 VUs for 4 seconds and then quickly scale down to 1 for the next 4s and then quit
					startVUs: 2,
					stages: [
						***REMOVED*** duration: "4s", target: 2 ***REMOVED***,
						***REMOVED*** duration: "0s", target: 1 ***REMOVED***,
						***REMOVED*** duration: "4s", target: 1 ***REMOVED***,
					],
					gracefulStop: "0s",
					gracefulRampDown: "0s",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***;

		export default function () ***REMOVED***
			console.log("VU " + __VU + " starting iteration #" + __ITER);
			http.get("HTTPBIN_IP_URL/bytes/15000");
			sleep(3.1);
			http.get("HTTPBIN_IP_URL/bytes/15000");
			console.log("VU " + __VU + " ending iteration #" + __ITER);
		***REMOVED***;
	`))

	runner, err := js.New(
		&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: script***REMOVED***,
		nil,
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)
	require.NoError(t, err)

	engine, run, wait := newTestEngine(t, nil, runner, lib.Options***REMOVED******REMOVED***)

	collector := &dummy.Collector***REMOVED******REMOVED***
	engine.Collectors = []lib.Collector***REMOVED***collector***REMOVED***

	errC := make(chan error)
	go func() ***REMOVED*** errC <- run() ***REMOVED***()

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out")
	case err := <-errC:
		require.NoError(t, err)
		wait()
		require.False(t, engine.IsTainted())
	***REMOVED***

	// The 3.1 sleep in the default function would cause the first VU to complete 2 full iterations
	// and stat executing its third one, while the second VU will only fully complete 1 iteration
	// and will be canceled in the middle of its second one.
	assert.Equal(t, 3.0, getMetricSum(collector, metrics.Iterations.Name))

	// That means that we expect to see 8 HTTP requests in total, 3*2=6 from the complete iterations
	// and one each from the two iterations that would be canceled in the middle of their execution
	assert.Equal(t, 8.0, getMetricSum(collector, metrics.HTTPReqs.Name))

	// And we expect to see the data_received for all 8 of those requests. Previously, the data for
	// the 8th request (the 3rd one in the first VU before the test ends) was cut off by the engine
	// because it was emitted after the test officially ended. But that was mostly an unintended
	// consequence of the fact that those metrics were emitted only after an iteration ended when
	// it was interrupted.
	dataReceivedExpectedMin := 15000.0 * 8
	dataReceivedExpectedMax := (15000.0 + expectedHeaderMaxLength) * 8
	dataReceivedActual := getMetricSum(collector, metrics.DataReceived.Name)
	if dataReceivedActual < dataReceivedExpectedMin || dataReceivedActual > dataReceivedExpectedMax ***REMOVED***
		t.Errorf(
			"The data_received sum should be in the interval [%f, %f] but was %f",
			dataReceivedExpectedMin, dataReceivedExpectedMax, dataReceivedActual,
		)
	***REMOVED***

	// Also, the interrupted iterations shouldn't affect the average iteration_duration in any way, only
	// complete iterations should be taken into account
	durationCount := float64(getMetricCount(collector, metrics.IterationDuration.Name))
	assert.Equal(t, 3.0, durationCount)
	durationSum := getMetricSum(collector, metrics.IterationDuration.Name)
	assert.InDelta(t, 3.35, durationSum/(1000*durationCount), 0.25)
***REMOVED***

func TestMetricsEmission(t *testing.T) ***REMOVED***
	if !isWindows ***REMOVED***
		t.Parallel()
	***REMOVED***

	testCases := []struct ***REMOVED***
		method             string
		minIterDuration    string
		defaultBody        string
		expCount, expIters float64
	***REMOVED******REMOVED***
		// Since emission of Iterations happens before the minIterationDuration
		// sleep is done, we expect to receive metrics for all executions of
		// the `default` function, despite of the lower overall duration setting.
		***REMOVED***"minIterationDuration", `"300ms"`, "testCounter.add(1);", 16.0, 16.0***REMOVED***,
		// With the manual sleep method and no minIterationDuration, the last
		// `default` execution will be cutoff by the duration setting, so only
		// 3 sets of metrics are expected.
		***REMOVED***"sleepBeforeCounterAdd", "null", "sleep(0.3); testCounter.add(1); ", 12.0, 12.0***REMOVED***,
		// The counter should be sent, but the last iteration will be incomplete
		***REMOVED***"sleepAfterCounterAdd", "null", "testCounter.add(1); sleep(0.3); ", 16.0, 12.0***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.method, func(t *testing.T) ***REMOVED***
			if !isWindows ***REMOVED***
				t.Parallel()
			***REMOVED***
			runner, err := js.New(
				&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: []byte(fmt.Sprintf(`
				import ***REMOVED*** sleep ***REMOVED*** from "k6";
				import ***REMOVED*** Counter ***REMOVED*** from "k6/metrics";

				let testCounter = new Counter("testcounter");

				export let options = ***REMOVED***
					execution: ***REMOVED***
						we_need_hard_stop: ***REMOVED***
							type: "constant-looping-vus",
							vus: 4,
							duration: "1s",
							gracefulStop: "0s",
						***REMOVED***,
					***REMOVED***,
					minIterationDuration: %s,
				***REMOVED***;

				export default function() ***REMOVED***
					%s
				***REMOVED***
				`, tc.minIterDuration, tc.defaultBody))***REMOVED***,
				nil,
				lib.RuntimeOptions***REMOVED******REMOVED***,
			)
			require.NoError(t, err)

			engine, run, wait := newTestEngine(t, nil, runner, runner.GetOptions())

			collector := &dummy.Collector***REMOVED******REMOVED***
			engine.Collectors = []lib.Collector***REMOVED***collector***REMOVED***

			errC := make(chan error)
			go func() ***REMOVED*** errC <- run() ***REMOVED***()

			select ***REMOVED***
			case <-time.After(10 * time.Second):
				t.Fatal("Test timed out")
			case err := <-errC:
				require.NoError(t, err)
				wait()
				require.False(t, engine.IsTainted())
			***REMOVED***

			assert.Equal(t, tc.expIters, getMetricSum(collector, metrics.Iterations.Name))
			assert.Equal(t, tc.expCount, getMetricSum(collector, "testcounter"))
		***REMOVED***)
	***REMOVED***
***REMOVED***

//nolint: funlen
func TestMinIterationDurationInSetupTeardownStage(t *testing.T) ***REMOVED***
	t.Parallel()
	setupScript := `
		import ***REMOVED*** sleep ***REMOVED*** from "k6";

		export function setup() ***REMOVED***
			sleep(1);
		***REMOVED***

		export let options = ***REMOVED***
			minIterationDuration: "2s",
			execution: ***REMOVED***
				we_need_hard_stop: ***REMOVED***
					type: "constant-looping-vus",
					vus: 2,
					duration: "1.9s",
					gracefulStop: "0s",
				***REMOVED***,
			***REMOVED***,
			setupTimeout: "2s",
		***REMOVED***;

		export default function () ***REMOVED***
		***REMOVED***;`
	teardownScript := `
		import ***REMOVED*** sleep ***REMOVED*** from "k6";

		export let options = ***REMOVED***
			minIterationDuration: "2s",
			execution: ***REMOVED***
				we_need_hard_stop: ***REMOVED***
					type: "constant-looping-vus",
					vus: 2,
					duration: "1.9s",
					gracefulStop: "0s",
				***REMOVED***,
			***REMOVED***,
			teardownTimeout: "2s",
		***REMOVED***;

		export default function () ***REMOVED***
		***REMOVED***;

		export function teardown() ***REMOVED***
			sleep(1);
		***REMOVED***
`
	tests := []struct ***REMOVED***
		name, script string
	***REMOVED******REMOVED***
		***REMOVED***"Test setup", setupScript***REMOVED***,
		***REMOVED***"Test teardown", teardownScript***REMOVED***,
	***REMOVED***
	for _, tc := range tests ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			runner, err := js.New(
				&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: []byte(tc.script)***REMOVED***,
				nil,
				lib.RuntimeOptions***REMOVED******REMOVED***,
			)
			require.NoError(t, err)

			engine, run, wait := newTestEngine(t, nil, runner, runner.GetOptions())

			errC := make(chan error)
			go func() ***REMOVED*** errC <- run() ***REMOVED***()
			select ***REMOVED***
			case <-time.After(10 * time.Second):
				t.Fatal("Test timed out")
			case err := <-errC:
				require.NoError(t, err)
				wait()
				require.False(t, engine.IsTainted())
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
