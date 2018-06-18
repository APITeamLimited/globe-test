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
	"testing"
	"time"

	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/dummy"
	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

type testErrorWithString string

func (e testErrorWithString) Error() string  ***REMOVED*** return string(e) ***REMOVED***
func (e testErrorWithString) String() string ***REMOVED*** return string(e) ***REMOVED***

// Apply a null logger to the engine and return the hook.
func applyNullLogger(e *Engine) *logtest.Hook ***REMOVED***
	logger, hook := logtest.NewNullLogger()
	e.SetLogger(logger)
	return hook
***REMOVED***

// Wrapper around newEngine that applies a null logger.
func newTestEngine(ex lib.Executor, opts lib.Options) (*Engine, error, *logtest.Hook) ***REMOVED***
	e, err := NewEngine(ex, opts)
	if err != nil ***REMOVED***
		return e, err, nil
	***REMOVED***
	hook := applyNullLogger(e)
	return e, nil, hook
***REMOVED***

func L(r lib.Runner) lib.Executor ***REMOVED***
	return local.New(r)
***REMOVED***

func LF(fn func(ctx context.Context) ([]stats.SampleContainer, error)) lib.Executor ***REMOVED***
	return L(&lib.MiniRunner***REMOVED***Fn: fn***REMOVED***)
***REMOVED***

func TestNewEngine(t *testing.T) ***REMOVED***
	_, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
***REMOVED***

func TestNewEngineOptions(t *testing.T) ***REMOVED***
	t.Run("Duration", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Duration: types.NullDurationFrom(10 * time.Second),
		***REMOVED***)
		assert.NoError(t, err)
		assert.Nil(t, e.Executor.GetStages())
		assert.Equal(t, types.NullDurationFrom(10*time.Second), e.Executor.GetEndTime())

		t.Run("Infinite", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Duration: types.NullDuration***REMOVED******REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.Nil(t, e.Executor.GetStages())
			assert.Equal(t, types.NullDuration***REMOVED******REMOVED***, e.Executor.GetEndTime())
		***REMOVED***)
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Stages: []lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Executor.GetStages(), 1) ***REMOVED***
			assert.Equal(t, e.Executor.GetStages()[0], lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Stages/Duration", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Duration: types.NullDurationFrom(60 * time.Second),
			Stages: []lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		if assert.Len(t, e.Executor.GetStages(), 1) ***REMOVED***
			assert.Equal(t, e.Executor.GetStages()[0], lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***)
		***REMOVED***
		assert.Equal(t, types.NullDurationFrom(60*time.Second), e.Executor.GetEndTime())
	***REMOVED***)
	t.Run("Iterations", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Iterations: null.IntFrom(100)***REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, null.IntFrom(100), e.Executor.GetEndIterations())
	***REMOVED***)
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(0), e.Executor.GetVUsMax())
			assert.Equal(t, int64(0), e.Executor.GetVUs())
		***REMOVED***)
		t.Run("set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.Executor.GetVUsMax())
			assert.Equal(t, int64(0), e.Executor.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		t.Run("no max", func(t *testing.T) ***REMOVED***
			_, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUs: null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "can't raise vu count (to 10) above vu cap (0)")
		***REMOVED***)
		t.Run("negative max", func(t *testing.T) ***REMOVED***
			_, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(-1),
			***REMOVED***)
			assert.EqualError(t, err, "vu cap can't be negative")
		***REMOVED***)
		t.Run("max too low", func(t *testing.T) ***REMOVED***
			_, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(1),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "can't raise vu count (to 10) above vu cap (1)")
		***REMOVED***)
		t.Run("max higher", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(1),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.Executor.GetVUsMax())
			assert.Equal(t, int64(1), e.Executor.GetVUs())
		***REMOVED***)
		t.Run("max just right", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.Executor.GetVUsMax())
			assert.Equal(t, int64(10), e.Executor.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("Paused", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.Executor.IsPaused())
		***REMOVED***)
		t.Run("false", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				Paused: null.BoolFrom(false),
			***REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.Executor.IsPaused())
		***REMOVED***)
		t.Run("true", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				Paused: null.BoolFrom(true),
			***REMOVED***)
			assert.NoError(t, err)
			assert.True(t, e.Executor.IsPaused())
		***REMOVED***)
	***REMOVED***)
	t.Run("thresholds", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Thresholds: map[string]stats.Thresholds***REMOVED***
				"my_metric": ***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)
		assert.Contains(t, e.thresholds, "my_metric")

		t.Run("submetrics", func(t *testing.T) ***REMOVED***
			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
				Thresholds: map[string]stats.Thresholds***REMOVED***
					"my_metric***REMOVED***tag:value***REMOVED***": ***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED***)
			assert.NoError(t, err)
			assert.Contains(t, e.thresholds, "my_metric***REMOVED***tag:value***REMOVED***")
			assert.Contains(t, e.submetrics, "my_metric")
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestEngineRun(t *testing.T) ***REMOVED***
	log.SetLevel(log.DebugLevel)
	t.Run("exits with context", func(t *testing.T) ***REMOVED***
		duration := 100 * time.Millisecond
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()
		startTime := time.Now()
		assert.NoError(t, e.Run(ctx))
		assert.WithinDuration(t, startTime.Add(duration), time.Now(), 100*time.Millisecond)
	***REMOVED***)
	t.Run("exits with executor", func(t *testing.T) ***REMOVED***
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			VUs:        null.IntFrom(10),
			VUsMax:     null.IntFrom(10),
			Iterations: null.IntFrom(100),
		***REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.Run(context.Background()))
		assert.Equal(t, int64(100), e.Executor.GetIterations())
	***REMOVED***)

	// Make sure samples are discarded after context close (using "cutoff" timestamp in local.go)
	t.Run("collects samples", func(t *testing.T) ***REMOVED***
		testMetric := stats.New("test_metric", stats.Trend)

		signalChan := make(chan interface***REMOVED******REMOVED***)
		var e *Engine
		e, err, _ := newTestEngine(LF(func(ctx context.Context) (samples []stats.SampleContainer, err error) ***REMOVED***
			samples = append(samples, stats.Sample***REMOVED***Metric: testMetric, Time: time.Now(), Value: 1***REMOVED***)
			close(signalChan)
			<-ctx.Done()

			// HACK(robin): Add a sleep here to temporarily workaround two problems with this test:
			// 1. The sample times are compared against the `cutoff` in core/local/local.go and sometimes the
			//    second sample (below) gets a `Time` smaller than `cutoff` because the lines below get executed
			//    before the `<-ctx.Done()` select in local.go:Run() on multi-core systems where
			//    goroutines can run in parallel.
			// 2. Sometimes the `case samples := <-vuOut` gets selected before the `<-ctx.Done()` in
			//    core/local/local.go:Run() causing all samples from this mocked "RunOnce()" function to be accepted.
			time.Sleep(time.Millisecond * 10)
			samples = append(samples, stats.Sample***REMOVED***Metric: testMetric, Time: time.Now(), Value: 2***REMOVED***)
			return samples, err
		***REMOVED***), lib.Options***REMOVED***
			VUs:        null.IntFrom(1),
			VUsMax:     null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***

		c := &dummy.Collector***REMOVED******REMOVED***
		e.Collectors = []lib.Collector***REMOVED***c***REMOVED***

		ctx, cancel := context.WithCancel(context.Background())
		errC := make(chan error)
		go func() ***REMOVED*** errC <- e.Run(ctx) ***REMOVED***()
		<-signalChan
		cancel()
		assert.NoError(t, <-errC)

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
	e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	assert.NoError(t, e.Run(ctx))
***REMOVED***

func TestEngineCollector(t *testing.T) ***REMOVED***
	testMetric := stats.New("test_metric", stats.Trend)

	e, err, _ := newTestEngine(LF(func(ctx context.Context) ([]stats.SampleContainer, error) ***REMOVED***
		return []stats.SampleContainer***REMOVED***stats.Sample***REMOVED***Metric: testMetric***REMOVED******REMOVED***, nil
	***REMOVED***), lib.Options***REMOVED***VUs: null.IntFrom(1), VUsMax: null.IntFrom(1), Iterations: null.IntFrom(1)***REMOVED***)
	assert.NoError(t, err)

	c := &dummy.Collector***REMOVED******REMOVED***
	e.Collectors = []lib.Collector***REMOVED***c***REMOVED***

	assert.NoError(t, e.Run(context.Background()))

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
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		e.processSamples(
			stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED***,
		)

		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric"].Sink)
	***REMOVED***)
	t.Run("submetric", func(t *testing.T) ***REMOVED***
		ths, err := stats.NewThresholds([]string***REMOVED***`1+1==2`***REMOVED***)
		assert.NoError(t, err)

		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***
			Thresholds: map[string]stats.Thresholds***REMOVED***
				"my_metric***REMOVED***a:1***REMOVED***": ths,
			***REMOVED***,
		***REMOVED***)
		assert.NoError(t, err)

		sms := e.submetrics["my_metric"]
		assert.Len(t, sms, 1)
		assert.Equal(t, "my_metric***REMOVED***a:1***REMOVED***", sms[0].Name)
		assert.EqualValues(t, map[string]string***REMOVED***"a": "1"***REMOVED***, sms[0].Tags.CloneTags())

		e.processSamples(
			stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED***)***REMOVED***,
		)

		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric"].Sink)
		assert.IsType(t, &stats.GaugeSink***REMOVED******REMOVED***, e.Metrics["my_metric***REMOVED***a:1***REMOVED***"].Sink)
	***REMOVED***)
***REMOVED***

func TestEngine_runThresholds(t *testing.T) ***REMOVED***
	metric := stats.New("my_metric", stats.Gauge)
	thresholds := make(map[string]stats.Thresholds, 1)

	ths, err := stats.NewThresholds([]string***REMOVED***"1+1==3"***REMOVED***)
	assert.NoError(t, err)

	t.Run("aborted", func(t *testing.T) ***REMOVED***
		ths.Thresholds[0].AbortOnFail = true
		thresholds[metric.Name] = ths
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Thresholds: thresholds***REMOVED***)
		assert.NoError(t, err)

		e.processSamples(
			stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED***,
		)

		ctx, cancel := context.WithCancel(context.Background())
		aborted := false

		cancelFunc := func() ***REMOVED***
			cancel()
			aborted = true
		***REMOVED***

		e.runThresholds(ctx, cancelFunc)

		assert.True(t, aborted)
	***REMOVED***)

	t.Run("canceled", func(t *testing.T) ***REMOVED***
		ths.Abort = false
		thresholds[metric.Name] = ths
		e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Thresholds: thresholds***REMOVED***)
		assert.NoError(t, err)

		e.processSamples(
			stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED***,
		)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		done := make(chan struct***REMOVED******REMOVED***)
		go func() ***REMOVED***
			defer close(done)
			e.runThresholds(ctx, cancel)
		***REMOVED***()

		select ***REMOVED***
		case <-done:
			return
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Test should have completed within a second")
		***REMOVED***
	***REMOVED***)
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
		t.Run(name, func(t *testing.T) ***REMOVED***
			thresholds := make(map[string]stats.Thresholds, len(data.ths))
			for m, srcs := range data.ths ***REMOVED***
				ths, err := stats.NewThresholds(srcs)
				assert.NoError(t, err)
				ths.Thresholds[0].AbortOnFail = data.abort
				thresholds[m] = ths
			***REMOVED***

			e, err, _ := newTestEngine(nil, lib.Options***REMOVED***Thresholds: thresholds***REMOVED***)
			assert.NoError(t, err)

			e.processSamples(
				stats.Sample***REMOVED***Metric: metric, Value: 1.25, Tags: stats.IntoSampleTags(&map[string]string***REMOVED***"a": "1"***REMOVED***)***REMOVED***,
			)

			abortCalled := false

			abortFunc := func() ***REMOVED***
				abortCalled = true
			***REMOVED***

			e.processThresholds(abortFunc)

			assert.Equal(t, data.pass, !e.IsTainted())
			if data.abort ***REMOVED***
				assert.True(t, abortCalled)
			***REMOVED***
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
func TestSentReceivedMetrics(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := testutils.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	tr := tb.Replacer.Replace

	const expectedHeaderMaxLength = 500

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
		***REMOVED***tr(`import ws from "k6/ws";
			let data = "0123456789".repeat(100);
			export default function() ***REMOVED***
				ws.connect("ws://HTTPBIN_IP:HTTPBIN_PORT/ws-echo", null, function (socket) ***REMOVED***
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
		***REMOVED***1, 1***REMOVED***, ***REMOVED***1, 2***REMOVED***, ***REMOVED***2, 1***REMOVED***, ***REMOVED***5, 2***REMOVED***, ***REMOVED***25, 2***REMOVED***, ***REMOVED***50, 5***REMOVED***,
	***REMOVED***

	runTest := func(t *testing.T, ts testScript, tc testCase, noConnReuse bool) (float64, float64) ***REMOVED***
		r, err := js.New(
			&lib.SourceData***REMOVED***Filename: "/script.js", Data: []byte(ts.Code)***REMOVED***,
			afero.NewMemMapFs(),
			lib.RuntimeOptions***REMOVED******REMOVED***,
		)
		require.NoError(t, err)

		options := lib.Options***REMOVED***
			Iterations: null.IntFrom(tc.Iterations),
			VUs:        null.IntFrom(tc.VUs),
			VUsMax:     null.IntFrom(tc.VUs),
			Hosts:      tb.Dialer.Hosts,
			InsecureSkipTLSVerify: null.BoolFrom(true),
			NoVUConnectionReuse:   null.BoolFrom(noConnReuse),
		***REMOVED***

		r.SetOptions(options)
		engine, err := NewEngine(local.New(r), options)
		require.NoError(t, err)

		collector := &dummy.Collector***REMOVED******REMOVED***
		engine.Collectors = []lib.Collector***REMOVED***collector***REMOVED***

		ctx, cancel := context.WithCancel(context.Background())
		errC := make(chan error)
		go func() ***REMOVED*** errC <- engine.Run(ctx) ***REMOVED***()

		select ***REMOVED***
		case <-time.After(10 * time.Second):
			cancel()
			t.Fatal("Test timed out")
		case err := <-errC:
			cancel()
			require.NoError(t, err)
		***REMOVED***

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
				t.Errorf("noReuseSent=%f is greater than reuseSent=%f", noReuseSent, reuseSent)
			***REMOVED***
			if noReuseReceived < reuseReceived ***REMOVED***
				t.Errorf("noReuseReceived=%f is greater than reuseReceived=%f", noReuseReceived, reuseReceived)
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
	tb := testutils.NewHTTPMultiBin(t)
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
				var response = ws.connect("wss://HTTPSBIN_IP:HTTPSBIN_PORT/ws-echo", params, function (socket) ***REMOVED***
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
		&lib.SourceData***REMOVED***Filename: "/script.js", Data: script***REMOVED***,
		afero.NewMemMapFs(),
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)
	require.NoError(t, err)

	options := lib.Options***REMOVED***
		Iterations:            null.IntFrom(3),
		VUs:                   null.IntFrom(2),
		VUsMax:                null.IntFrom(2),
		Hosts:                 tb.Dialer.Hosts,
		RunTags:               runTags,
		SystemTags:            lib.GetTagSet(lib.DefaultSystemTagList...),
		InsecureSkipTLSVerify: null.BoolFrom(true),
	***REMOVED***

	r.SetOptions(options)
	engine, err := NewEngine(local.New(r), options)
	require.NoError(t, err)

	collector := &dummy.Collector***REMOVED******REMOVED***
	engine.Collectors = []lib.Collector***REMOVED***collector***REMOVED***

	ctx, cancel := context.WithCancel(context.Background())
	errC := make(chan error)
	go func() ***REMOVED*** errC <- engine.Run(ctx) ***REMOVED***()

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("Test timed out")
	case err := <-errC:
		cancel()
		require.NoError(t, err)
	***REMOVED***

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
