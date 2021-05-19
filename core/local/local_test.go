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

package local

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/js"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/netext"
	"go.k6.io/k6/lib/netext/httpext"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/lib/testutils/minirunner"
	"go.k6.io/k6/lib/testutils/mockresolver"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/stats"
)

func newTestExecutionScheduler(
	t *testing.T, runner lib.Runner, logger *logrus.Logger, opts lib.Options,
) (ctx context.Context, cancel func(), execScheduler *ExecutionScheduler, samples chan stats.SampleContainer) ***REMOVED***
	if runner == nil ***REMOVED***
		runner = &minirunner.MiniRunner***REMOVED******REMOVED***
	***REMOVED***
	ctx, cancel = context.WithCancel(context.Background())
	newOpts, err := executor.DeriveScenariosFromShortcuts(lib.Options***REMOVED***
		MetricSamplesBufferSize: null.NewInt(200, false),
	***REMOVED***.Apply(runner.GetOptions()).Apply(opts))
	require.NoError(t, err)
	require.Empty(t, newOpts.Validate())

	require.NoError(t, runner.SetOptions(newOpts))

	if logger == nil ***REMOVED***
		logger = logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
	***REMOVED***

	execScheduler, err = NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	samples = make(chan stats.SampleContainer, newOpts.MetricSamplesBufferSize.Int64)
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

	return ctx, cancel, execScheduler, samples
***REMOVED***

func TestExecutionSchedulerRun(t *testing.T) ***REMOVED***
	t.Parallel()
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, nil, nil, lib.Options***REMOVED******REMOVED***)
	defer cancel()

	err := make(chan error, 1)
	go func() ***REMOVED*** err <- execScheduler.Run(ctx, ctx, samples) ***REMOVED***()
	assert.NoError(t, <-err)
***REMOVED***

func TestExecutionSchedulerRunNonDefault(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		name, script, expErr string
	***REMOVED******REMOVED***
		***REMOVED***"defaultOK", `export default function () ***REMOVED******REMOVED***`, ""***REMOVED***,
		***REMOVED***"nonDefaultOK", `
	export let options = ***REMOVED***
		scenarios: ***REMOVED***
			per_vu_iters: ***REMOVED***
				executor: "per-vu-iterations",
				vus: 1,
				iterations: 1,
				exec: "nonDefault",
			***REMOVED***,
		***REMOVED***
	***REMOVED***
	export function nonDefault() ***REMOVED******REMOVED***`, ""***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			logger := logrus.New()
			logger.SetOutput(testutils.NewTestOutput(t))
			runner, err := js.New(logger, &loader.SourceData***REMOVED***
				URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: []byte(tc.script),
			***REMOVED***,
				nil, lib.RuntimeOptions***REMOVED******REMOVED***)
			require.NoError(t, err)

			execScheduler, err := NewExecutionScheduler(runner, logger)
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			done := make(chan struct***REMOVED******REMOVED***)
			samples := make(chan stats.SampleContainer)
			go func() ***REMOVED***
				err := execScheduler.Init(ctx, samples)
				if tc.expErr != "" ***REMOVED***
					assert.EqualError(t, err, tc.expErr)
				***REMOVED*** else ***REMOVED***
					assert.NoError(t, err)
					assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
				***REMOVED***
				close(done)
			***REMOVED***()
			for ***REMOVED***
				select ***REMOVED***
				case <-samples:
				case <-done:
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerRunEnv(t *testing.T) ***REMOVED***
	t.Parallel()

	scriptTemplate := `
	import ***REMOVED*** Counter ***REMOVED*** from "k6/metrics";

	let errors = new Counter("errors");

	export let options = ***REMOVED***
		scenarios: ***REMOVED***
			executor: ***REMOVED***
				executor: "%[1]s",
				gracefulStop: "0.5s",
				%[2]s
			***REMOVED***
		***REMOVED***
	***REMOVED***

	export default function () ***REMOVED***
		if (__ENV.TESTVAR !== "%[3]s") ***REMOVED***
		    console.error('Wrong env var value. Expected: %[3]s, actual: ', __ENV.TESTVAR);
			errors.add(1);
		***REMOVED***
	***REMOVED***`

	executorConfigs := map[string]string***REMOVED***
		"constant-arrival-rate": `
			rate: 1,
			timeUnit: "0.5s",
			duration: "0.5s",
			preAllocatedVUs: 1,
			maxVUs: 2,`,
		"constant-vus": `
			vus: 1,
			duration: "0.5s",`,
		"externally-controlled": `
			vus: 1,
			duration: "0.5s",`,
		"per-vu-iterations": `
			vus: 1,
			iterations: 1,`,
		"shared-iterations": `
			vus: 1,
			iterations: 1,`,
		"ramping-arrival-rate": `
			startRate: 1,
			timeUnit: "0.5s",
			preAllocatedVUs: 1,
			maxVUs: 2,
			stages: [ ***REMOVED*** target: 1, duration: "0.5s" ***REMOVED*** ],`,
		"ramping-vus": `
			startVUs: 1,
			stages: [ ***REMOVED*** target: 1, duration: "0.5s" ***REMOVED*** ],`,
	***REMOVED***

	testCases := []struct***REMOVED*** name, script string ***REMOVED******REMOVED******REMOVED***

	// Generate tests using global env and with env override
	for ename, econf := range executorConfigs ***REMOVED***
		testCases = append(testCases, struct***REMOVED*** name, script string ***REMOVED******REMOVED***
			"global/" + ename, fmt.Sprintf(scriptTemplate, ename, econf, "global"),
		***REMOVED***)
		configWithEnvOverride := econf + "env: ***REMOVED*** TESTVAR: 'overridden' ***REMOVED***"
		testCases = append(testCases, struct***REMOVED*** name, script string ***REMOVED******REMOVED***
			"override/" + ename, fmt.Sprintf(scriptTemplate, ename, configWithEnvOverride, "overridden"),
		***REMOVED***)
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			logger := logrus.New()
			logger.SetOutput(testutils.NewTestOutput(t))
			runner, err := js.New(logger, &loader.SourceData***REMOVED***
				URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
				Data: []byte(tc.script),
			***REMOVED***,
				nil, lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"TESTVAR": "global"***REMOVED******REMOVED***)
			require.NoError(t, err)

			execScheduler, err := NewExecutionScheduler(runner, logger)
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			done := make(chan struct***REMOVED******REMOVED***)
			samples := make(chan stats.SampleContainer)
			go func() ***REMOVED***
				assert.NoError(t, execScheduler.Init(ctx, samples))
				assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
				close(done)
			***REMOVED***()
			for ***REMOVED***
				select ***REMOVED***
				case sample := <-samples:
					if s, ok := sample.(stats.Sample); ok && s.Metric.Name == "errors" ***REMOVED***
						assert.FailNow(t, "received error sample from test")
					***REMOVED***
				case <-done:
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerSystemTags(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace

	script := sr(`
	import http from "k6/http";

	export let options = ***REMOVED***
		scenarios: ***REMOVED***
			per_vu_test: ***REMOVED***
				executor: "per-vu-iterations",
				gracefulStop: "0s",
				vus: 1,
				iterations: 1,
			***REMOVED***,
			shared_test: ***REMOVED***
				executor: "shared-iterations",
				gracefulStop: "0s",
				vus: 1,
				iterations: 1,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	export default function () ***REMOVED***
		http.get("HTTPBIN_IP_URL/");
	***REMOVED***`)

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	runner, err := js.New(logger, &loader.SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
		Data: []byte(script),
	***REMOVED***,
		nil, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	require.NoError(t, runner.SetOptions(runner.GetOptions().Apply(lib.Options***REMOVED***
		SystemTags: &stats.DefaultSystemTagSet,
	***REMOVED***)))

	execScheduler, err := NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	samples := make(chan stats.SampleContainer)
	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer close(done)
		require.NoError(t, execScheduler.Init(ctx, samples))
		require.NoError(t, execScheduler.Run(ctx, ctx, samples))
	***REMOVED***()

	expCommonTrailTags := stats.IntoSampleTags(&map[string]string***REMOVED***
		"group":             "",
		"method":            "GET",
		"name":              sr("HTTPBIN_IP_URL/"),
		"url":               sr("HTTPBIN_IP_URL/"),
		"proto":             "HTTP/1.1",
		"status":            "200",
		"expected_response": "true",
	***REMOVED***)
	expTrailPVUTagsRaw := expCommonTrailTags.CloneTags()
	expTrailPVUTagsRaw["scenario"] = "per_vu_test"
	expTrailPVUTags := stats.IntoSampleTags(&expTrailPVUTagsRaw)
	expTrailSITagsRaw := expCommonTrailTags.CloneTags()
	expTrailSITagsRaw["scenario"] = "shared_test"
	expTrailSITags := stats.IntoSampleTags(&expTrailSITagsRaw)
	expNetTrailPVUTags := stats.IntoSampleTags(&map[string]string***REMOVED***
		"group":    "",
		"scenario": "per_vu_test",
	***REMOVED***)
	expNetTrailSITags := stats.IntoSampleTags(&map[string]string***REMOVED***
		"group":    "",
		"scenario": "shared_test",
	***REMOVED***)

	var gotCorrectTags int
	for ***REMOVED***
		select ***REMOVED***
		case sample := <-samples:
			switch s := sample.(type) ***REMOVED***
			case *httpext.Trail:
				if s.Tags.IsEqual(expTrailPVUTags) || s.Tags.IsEqual(expTrailSITags) ***REMOVED***
					gotCorrectTags++
				***REMOVED***
			case *netext.NetTrail:
				if s.Tags.IsEqual(expNetTrailPVUTags) || s.Tags.IsEqual(expNetTrailSITags) ***REMOVED***
					gotCorrectTags++
				***REMOVED***
			***REMOVED***
		case <-done:
			require.Equal(t, 4, gotCorrectTags, "received wrong amount of samples with expected tags")
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerRunCustomTags(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace

	scriptTemplate := sr(`
	import http from "k6/http";

	export let options = ***REMOVED***
		scenarios: ***REMOVED***
			executor: ***REMOVED***
				executor: "%s",
				gracefulStop: "0.5s",
				%s
			***REMOVED***
		***REMOVED***
	***REMOVED***

	export default function () ***REMOVED***
		http.get("HTTPBIN_IP_URL/");
	***REMOVED***`)

	executorConfigs := map[string]string***REMOVED***
		"constant-arrival-rate": `
			rate: 1,
			timeUnit: "0.5s",
			duration: "0.5s",
			preAllocatedVUs: 1,
			maxVUs: 2,`,
		"constant-vus": `
			vus: 1,
			duration: "0.5s",`,
		"externally-controlled": `
			vus: 1,
			duration: "0.5s",`,
		"per-vu-iterations": `
			vus: 1,
			iterations: 1,`,
		"shared-iterations": `
			vus: 1,
			iterations: 1,`,
		"ramping-arrival-rate": `
			startRate: 5,
			timeUnit: "0.5s",
			preAllocatedVUs: 1,
			maxVUs: 2,
			stages: [ ***REMOVED*** target: 10, duration: "1s" ***REMOVED*** ],`,
		"ramping-vus": `
			startVUs: 1,
			stages: [ ***REMOVED*** target: 1, duration: "0.5s" ***REMOVED*** ],`,
	***REMOVED***

	testCases := []struct***REMOVED*** name, script string ***REMOVED******REMOVED******REMOVED***

	// Generate tests using custom tags
	for ename, econf := range executorConfigs ***REMOVED***
		configWithCustomTag := econf + "tags: ***REMOVED*** customTag: 'value' ***REMOVED***"
		testCases = append(testCases, struct***REMOVED*** name, script string ***REMOVED******REMOVED***
			ename, fmt.Sprintf(scriptTemplate, ename, configWithCustomTag),
		***REMOVED***)
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			logger := logrus.New()
			logger.SetOutput(testutils.NewTestOutput(t))

			runner, err := js.New(logger, &loader.SourceData***REMOVED***
				URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
				Data: []byte(tc.script),
			***REMOVED***,
				nil, lib.RuntimeOptions***REMOVED******REMOVED***)
			require.NoError(t, err)

			execScheduler, err := NewExecutionScheduler(runner, logger)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			done := make(chan struct***REMOVED******REMOVED***)
			samples := make(chan stats.SampleContainer)
			go func() ***REMOVED***
				defer close(done)
				require.NoError(t, execScheduler.Init(ctx, samples))
				require.NoError(t, execScheduler.Run(ctx, ctx, samples))
			***REMOVED***()
			var gotTrailTag, gotNetTrailTag bool
			for ***REMOVED***
				select ***REMOVED***
				case sample := <-samples:
					if trail, ok := sample.(*httpext.Trail); ok && !gotTrailTag ***REMOVED***
						tags := trail.Tags.CloneTags()
						if v, ok := tags["customTag"]; ok && v == "value" ***REMOVED***
							gotTrailTag = true
						***REMOVED***
					***REMOVED***
					if netTrail, ok := sample.(*netext.NetTrail); ok && !gotNetTrailTag ***REMOVED***
						tags := netTrail.Tags.CloneTags()
						if v, ok := tags["customTag"]; ok && v == "value" ***REMOVED***
							gotNetTrailTag = true
						***REMOVED***
					***REMOVED***
				case <-done:
					if !gotTrailTag || !gotNetTrailTag ***REMOVED***
						assert.FailNow(t, "a sample with expected tag wasn't received")
					***REMOVED***
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// Ensure that custom executor settings are unique per executor and
// that there's no "crossover"/"pollution" between executors.
// Also test that custom tags are properly set on checks and groups metrics.
func TestExecutionSchedulerRunCustomConfigNoCrossover(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	script := tb.Replacer.Replace(`
	import http from "k6/http";
	import ws from 'k6/ws';
	import ***REMOVED*** Counter ***REMOVED*** from 'k6/metrics';
	import ***REMOVED*** check, group ***REMOVED*** from 'k6';

	let errors = new Counter('errors');

	export let options = ***REMOVED***
		// Required for WS tests
		hosts: ***REMOVED*** 'httpbin.local': '127.0.0.1' ***REMOVED***,
		scenarios: ***REMOVED***
			scenario1: ***REMOVED***
				executor: 'per-vu-iterations',
				vus: 1,
				iterations: 1,
				gracefulStop: '0s',
				maxDuration: '1s',
				exec: 's1func',
				env: ***REMOVED*** TESTVAR1: 'scenario1' ***REMOVED***,
				tags: ***REMOVED*** testtag1: 'scenario1' ***REMOVED***,
			***REMOVED***,
			scenario2: ***REMOVED***
				executor: 'shared-iterations',
				vus: 1,
				iterations: 1,
				gracefulStop: '1s',
				startTime: '0.5s',
				maxDuration: '2s',
				exec: 's2func',
				env: ***REMOVED*** TESTVAR2: 'scenario2' ***REMOVED***,
				tags: ***REMOVED*** testtag2: 'scenario2' ***REMOVED***,
			***REMOVED***,
			scenario3: ***REMOVED***
				executor: 'per-vu-iterations',
				vus: 1,
				iterations: 1,
				gracefulStop: '1s',
				exec: 's3funcWS',
				env: ***REMOVED*** TESTVAR3: 'scenario3' ***REMOVED***,
				tags: ***REMOVED*** testtag3: 'scenario3' ***REMOVED***,
			***REMOVED***,
		***REMOVED***
	***REMOVED***

	function checkVar(name, expected) ***REMOVED***
		if (__ENV[name] !== expected) ***REMOVED***
		    console.error('Wrong ' + name + " env var value. Expected: '"
						+ expected + "', actual: '" + __ENV[name] + "'");
			errors.add(1);
		***REMOVED***
	***REMOVED***

	export function s1func() ***REMOVED***
		checkVar('TESTVAR1', 'scenario1');
		checkVar('TESTVAR2', undefined);
		checkVar('TESTVAR3', undefined);
		checkVar('TESTGLOBALVAR', 'global');

		// Intentionally try to pollute the env
		__ENV.TESTVAR2 = 'overridden';

		http.get('HTTPBIN_IP_URL/', ***REMOVED*** tags: ***REMOVED*** reqtag: 'scenario1' ***REMOVED******REMOVED***);
	***REMOVED***

	export function s2func() ***REMOVED***
		checkVar('TESTVAR1', undefined);
		checkVar('TESTVAR2', 'scenario2');
		checkVar('TESTVAR3', undefined);
		checkVar('TESTGLOBALVAR', 'global');

		http.get('HTTPBIN_IP_URL/', ***REMOVED*** tags: ***REMOVED*** reqtag: 'scenario2' ***REMOVED******REMOVED***);
	***REMOVED***

	export function s3funcWS() ***REMOVED***
		checkVar('TESTVAR1', undefined);
		checkVar('TESTVAR2', undefined);
		checkVar('TESTVAR3', 'scenario3');
		checkVar('TESTGLOBALVAR', 'global');

		const customTags = ***REMOVED*** wstag: 'scenario3' ***REMOVED***;
		group('wsgroup', function() ***REMOVED***
			const response = ws.connect('WSBIN_URL/ws-echo', ***REMOVED*** tags: customTags ***REMOVED***,
				function (socket) ***REMOVED***
					socket.on('open', function() ***REMOVED***
						socket.send('hello');
					***REMOVED***);
					socket.on('message', function(msg) ***REMOVED***
						if (msg != 'hello') ***REMOVED***
						    console.error("Expected to receive 'hello' but got '" + msg + "' instead!");
							errors.add(1);
						***REMOVED***
						socket.close()
					***REMOVED***);
					socket.on('error', function (e) ***REMOVED***
						console.log('ws error: ' + e.error());
						errors.add(1);
					***REMOVED***);
				***REMOVED***
			);
			check(response, ***REMOVED*** 'status is 101': (r) => r && r.status === 101 ***REMOVED***, customTags);
		***REMOVED***);
	***REMOVED***
`)
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	runner, err := js.New(logger, &loader.SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
		Data: []byte(script),
	***REMOVED***,
		nil, lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"TESTGLOBALVAR": "global"***REMOVED******REMOVED***)
	require.NoError(t, err)

	execScheduler, err := NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	samples := make(chan stats.SampleContainer)
	go func() ***REMOVED***
		assert.NoError(t, execScheduler.Init(ctx, samples))
		assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
		close(samples)
	***REMOVED***()

	expectedTrailTags := []map[string]string***REMOVED***
		***REMOVED***"testtag1": "scenario1", "reqtag": "scenario1"***REMOVED***,
		***REMOVED***"testtag2": "scenario2", "reqtag": "scenario2"***REMOVED***,
	***REMOVED***
	expectedNetTrailTags := []map[string]string***REMOVED***
		***REMOVED***"testtag1": "scenario1"***REMOVED***,
		***REMOVED***"testtag2": "scenario2"***REMOVED***,
	***REMOVED***
	expectedConnSampleTags := map[string]string***REMOVED***
		"testtag3": "scenario3", "wstag": "scenario3",
	***REMOVED***
	expectedPlainSampleTags := []map[string]string***REMOVED***
		***REMOVED***"testtag3": "scenario3"***REMOVED***,
		***REMOVED***"testtag3": "scenario3", "wstag": "scenario3"***REMOVED***,
	***REMOVED***
	var gotSampleTags int
	for sample := range samples ***REMOVED***
		switch s := sample.(type) ***REMOVED***
		case stats.Sample:
			if s.Metric.Name == "errors" ***REMOVED***
				assert.FailNow(t, "received error sample from test")
			***REMOVED***
			if s.Metric.Name == "checks" || s.Metric.Name == "group_duration" ***REMOVED***
				tags := s.Tags.CloneTags()
				for _, expTags := range expectedPlainSampleTags ***REMOVED***
					if reflect.DeepEqual(expTags, tags) ***REMOVED***
						gotSampleTags++
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case *httpext.Trail:
			tags := s.Tags.CloneTags()
			for _, expTags := range expectedTrailTags ***REMOVED***
				if reflect.DeepEqual(expTags, tags) ***REMOVED***
					gotSampleTags++
				***REMOVED***
			***REMOVED***
		case *netext.NetTrail:
			tags := s.Tags.CloneTags()
			for _, expTags := range expectedNetTrailTags ***REMOVED***
				if reflect.DeepEqual(expTags, tags) ***REMOVED***
					gotSampleTags++
				***REMOVED***
			***REMOVED***
		case stats.ConnectedSamples:
			for _, sm := range s.Samples ***REMOVED***
				tags := sm.Tags.CloneTags()
				if reflect.DeepEqual(expectedConnSampleTags, tags) ***REMOVED***
					gotSampleTags++
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	require.Equal(t, 8, gotSampleTags, "received wrong amount of samples with expected tags")
***REMOVED***

func TestExecutionSchedulerSetupTeardownRun(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Normal", func(t *testing.T) ***REMOVED***
		t.Parallel()
		setupC := make(chan struct***REMOVED******REMOVED***)
		teardownC := make(chan struct***REMOVED******REMOVED***)
		runner := &minirunner.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				close(setupC)
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				close(teardownC)
				return nil
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED******REMOVED***)

		err := make(chan error, 1)
		go func() ***REMOVED*** err <- execScheduler.Run(ctx, ctx, samples) ***REMOVED***()
		defer cancel()
		<-setupC
		<-teardownC
		assert.NoError(t, <-err)
	***REMOVED***)
	t.Run("Setup Error", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, errors.New("setup error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED******REMOVED***)
		defer cancel()
		assert.EqualError(t, execScheduler.Run(ctx, ctx, samples), "setup error")
	***REMOVED***)
	t.Run("Don't Run Setup", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, errors.New("setup error")
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
			NoSetup:    null.BoolFrom(true),
			VUs:        null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		defer cancel()
		assert.EqualError(t, execScheduler.Run(ctx, ctx, samples), "teardown error")
	***REMOVED***)

	t.Run("Teardown Error", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
			VUs:        null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		defer cancel()

		assert.EqualError(t, execScheduler.Run(ctx, ctx, samples), "teardown error")
	***REMOVED***)
	t.Run("Don't Run Teardown", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED***
			SetupFn: func(ctx context.Context, out chan<- stats.SampleContainer) ([]byte, error) ***REMOVED***
				return nil, nil
			***REMOVED***,
			TeardownFn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
				return errors.New("teardown error")
			***REMOVED***,
		***REMOVED***
		ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
			NoTeardown: null.BoolFrom(true),
			VUs:        null.IntFrom(1),
			Iterations: null.IntFrom(1),
		***REMOVED***)
		defer cancel()
		assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
	***REMOVED***)
***REMOVED***

func TestExecutionSchedulerStages(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		Duration time.Duration
		Stages   []lib.Stage
	***REMOVED******REMOVED***
		"one": ***REMOVED***
			1 * time.Second,
			[]lib.Stage***REMOVED******REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(1)***REMOVED******REMOVED***,
		***REMOVED***,
		"two": ***REMOVED***
			2 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(1)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(2)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"four": ***REMOVED***
			4 * time.Second,
			[]lib.Stage***REMOVED***
				***REMOVED***Duration: types.NullDurationFrom(1 * time.Second), Target: null.IntFrom(5)***REMOVED***,
				***REMOVED***Duration: types.NullDurationFrom(3 * time.Second), Target: null.IntFrom(10)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		data := data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			runner := &minirunner.MiniRunner***REMOVED***
				Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
					time.Sleep(100 * time.Millisecond)
					return nil
				***REMOVED***,
			***REMOVED***
			ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
				VUs:    null.IntFrom(1),
				Stages: data.Stages,
			***REMOVED***)
			defer cancel()
			assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
			assert.True(t, execScheduler.GetState().GetCurrentTestRunDuration() >= data.Duration)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerEndTime(t *testing.T) ***REMOVED***
	t.Parallel()
	runner := &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			time.Sleep(100 * time.Millisecond)
			return nil
		***REMOVED***,
	***REMOVED***
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED***
		VUs:      null.IntFrom(10),
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***)
	defer cancel()

	endTime, isFinal := lib.GetEndOffset(execScheduler.GetExecutionPlan())
	assert.Equal(t, 31*time.Second, endTime) // because of the default 30s gracefulStop
	assert.True(t, isFinal)

	startTime := time.Now()
	assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
	runTime := time.Since(startTime)
	assert.True(t, runTime > 1*time.Second, "test did not take 1s")
	assert.True(t, runTime < 10*time.Second, "took more than 10 seconds")
***REMOVED***

func TestExecutionSchedulerRuntimeErrors(t *testing.T) ***REMOVED***
	t.Parallel()
	runner := &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			time.Sleep(10 * time.Millisecond)
			return errors.New("hi")
		***REMOVED***,
		Options: lib.Options***REMOVED***
			VUs:      null.IntFrom(10),
			Duration: types.NullDurationFrom(1 * time.Second),
		***REMOVED***,
	***REMOVED***
	logger, hook := logtest.NewNullLogger()
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, lib.Options***REMOVED******REMOVED***)
	defer cancel()

	endTime, isFinal := lib.GetEndOffset(execScheduler.GetExecutionPlan())
	assert.Equal(t, 31*time.Second, endTime) // because of the default 30s gracefulStop
	assert.True(t, isFinal)

	startTime := time.Now()
	assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
	runTime := time.Since(startTime)
	assert.True(t, runTime > 1*time.Second, "test did not take 1s")
	assert.True(t, runTime < 10*time.Second, "took more than 10 seconds")

	assert.NotEmpty(t, hook.Entries)
	for _, e := range hook.Entries ***REMOVED***
		assert.Equal(t, "hi", e.Message)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerEndErrors(t *testing.T) ***REMOVED***
	t.Parallel()

	exec := executor.NewConstantVUsConfig("we_need_hard_stop")
	exec.VUs = null.IntFrom(10)
	exec.Duration = types.NullDurationFrom(1 * time.Second)
	exec.GracefulStop = types.NullDurationFrom(0 * time.Second)

	runner := &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			<-ctx.Done()
			return errors.New("hi")
		***REMOVED***,
		Options: lib.Options***REMOVED***
			Scenarios: lib.ScenarioConfigs***REMOVED***exec.GetName(): exec***REMOVED***,
		***REMOVED***,
	***REMOVED***
	logger, hook := logtest.NewNullLogger()
	ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, lib.Options***REMOVED******REMOVED***)
	defer cancel()

	endTime, isFinal := lib.GetEndOffset(execScheduler.GetExecutionPlan())
	assert.Equal(t, 1*time.Second, endTime) // because of the 0s gracefulStop
	assert.True(t, isFinal)

	startTime := time.Now()
	assert.NoError(t, execScheduler.Run(ctx, ctx, samples))
	runTime := time.Since(startTime)
	assert.True(t, runTime > 1*time.Second, "test did not take 1s")
	assert.True(t, runTime < 10*time.Second, "took more than 10 seconds")

	assert.Empty(t, hook.Entries)
***REMOVED***

func TestExecutionSchedulerEndIterations(t *testing.T) ***REMOVED***
	t.Parallel()
	metric := &stats.Metric***REMOVED***Name: "test_metric"***REMOVED***

	options, err := executor.DeriveScenariosFromShortcuts(lib.Options***REMOVED***
		VUs:        null.IntFrom(1),
		Iterations: null.IntFrom(100),
	***REMOVED***)
	require.NoError(t, err)
	require.Empty(t, options.Validate())

	var i int64
	runner := &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
			default:
				atomic.AddInt64(&i, 1)
			***REMOVED***
			out <- stats.Sample***REMOVED***Metric: metric, Value: 1.0***REMOVED***
			return nil
		***REMOVED***,
		Options: options,
	***REMOVED***

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	execScheduler, err := NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	samples := make(chan stats.SampleContainer, 300)
	require.NoError(t, execScheduler.Init(ctx, samples))
	require.NoError(t, execScheduler.Run(ctx, ctx, samples))

	assert.Equal(t, uint64(100), execScheduler.GetState().GetFullIterationCount())
	assert.Equal(t, uint64(0), execScheduler.GetState().GetPartialIterationCount())
	assert.Equal(t, int64(100), i)
	require.Equal(t, 100, len(samples)) // TODO: change to 200 https://github.com/k6io/k6/issues/1250
	for i := 0; i < 100; i++ ***REMOVED***
		mySample, ok := <-samples
		require.True(t, ok)
		assert.Equal(t, stats.Sample***REMOVED***Metric: metric, Value: 1.0***REMOVED***, mySample)
	***REMOVED***
***REMOVED***

func TestExecutionSchedulerIsRunning(t *testing.T) ***REMOVED***
	t.Parallel()
	runner := &minirunner.MiniRunner***REMOVED***
		Fn: func(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
			<-ctx.Done()
			return nil
		***REMOVED***,
	***REMOVED***
	ctx, cancel, execScheduler, _ := newTestExecutionScheduler(t, runner, nil, lib.Options***REMOVED******REMOVED***)
	state := execScheduler.GetState()

	err := make(chan error)
	go func() ***REMOVED*** err <- execScheduler.Run(ctx, ctx, nil) ***REMOVED***()
	for !state.HasStarted() ***REMOVED***
		time.Sleep(10 * time.Microsecond)
	***REMOVED***
	cancel()
	for !state.HasEnded() ***REMOVED***
		time.Sleep(10 * time.Microsecond)
	***REMOVED***
	assert.NoError(t, <-err)
***REMOVED***

// TestDNSResolver checks the DNS resolution behavior at the ExecutionScheduler level.
func TestDNSResolver(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	sr := tb.Replacer.Replace
	script := sr(`
		import http from "k6/http";
		import ***REMOVED*** sleep ***REMOVED*** from "k6";

		export let options = ***REMOVED***
			vus: 1,
			iterations: 8,
			noConnectionReuse: true,
		***REMOVED***

		export default function () ***REMOVED***
			const res = http.get("http://myhost:HTTPBIN_PORT/", ***REMOVED*** timeout: 50 ***REMOVED***);
			sleep(0.7);  // somewhat uneven multiple of 0.5 to minimize races with asserts
		***REMOVED***`)

	t.Run("cache", func(t *testing.T) ***REMOVED***
		t.Parallel()
		testCases := map[string]struct ***REMOVED***
			opts          lib.Options
			expLogEntries int
		***REMOVED******REMOVED***
			"default": ***REMOVED*** // IPs are cached for 5m
				lib.Options***REMOVED***DNS: types.DefaultDNSConfig()***REMOVED***, 0,
			***REMOVED***,
			"0": ***REMOVED*** // cache is disabled, every request does a DNS lookup
				lib.Options***REMOVED***DNS: types.DNSConfig***REMOVED***
					TTL:    null.StringFrom("0"),
					Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSfirst, Valid: true***REMOVED***,
					Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv4, Valid: false***REMOVED***,
				***REMOVED******REMOVED***, 5,
			***REMOVED***,
			"1000": ***REMOVED*** // cache IPs for 1s, check that unitless values are interpreted as ms
				lib.Options***REMOVED***DNS: types.DNSConfig***REMOVED***
					TTL:    null.StringFrom("1000"),
					Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSfirst, Valid: true***REMOVED***,
					Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv4, Valid: false***REMOVED***,
				***REMOVED******REMOVED***, 4,
			***REMOVED***,
			"3s": ***REMOVED***
				lib.Options***REMOVED***DNS: types.DNSConfig***REMOVED***
					TTL:    null.StringFrom("3s"),
					Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSfirst, Valid: true***REMOVED***,
					Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv4, Valid: false***REMOVED***,
				***REMOVED******REMOVED***, 3,
			***REMOVED***,
		***REMOVED***

		expErr := sr(`dial tcp 127.0.0.254:HTTPBIN_PORT: connect: connection refused`)
		if runtime.GOOS == "windows" ***REMOVED***
			expErr = "request timeout"
		***REMOVED***
		for name, tc := range testCases ***REMOVED***
			tc := tc
			t.Run(name, func(t *testing.T) ***REMOVED***
				t.Parallel()
				logger := logrus.New()
				logger.SetOutput(ioutil.Discard)
				logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
				logger.AddHook(&logHook)

				runner, err := js.New(logger, &loader.SourceData***REMOVED***
					URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: []byte(script),
				***REMOVED***, nil, lib.RuntimeOptions***REMOVED******REMOVED***)
				require.NoError(t, err)

				mr := mockresolver.New(nil, net.LookupIP)
				runner.ActualResolver = mr.LookupIPAll

				ctx, cancel, execScheduler, samples := newTestExecutionScheduler(t, runner, logger, tc.opts)
				defer cancel()

				mr.Set("myhost", sr("HTTPBIN_IP"))
				time.AfterFunc(1700*time.Millisecond, func() ***REMOVED***
					mr.Set("myhost", "127.0.0.254")
				***REMOVED***)
				defer mr.Unset("myhost")

				errCh := make(chan error, 1)
				go func() ***REMOVED*** errCh <- execScheduler.Run(ctx, ctx, samples) ***REMOVED***()

				select ***REMOVED***
				case err := <-errCh:
					require.NoError(t, err)
					entries := logHook.Drain()
					require.Len(t, entries, tc.expLogEntries)
					for _, entry := range entries ***REMOVED***
						require.IsType(t, &url.Error***REMOVED******REMOVED***, entry.Data["error"])
						assert.EqualError(t, entry.Data["error"].(*url.Error).Err, expErr)
					***REMOVED***
				case <-time.After(10 * time.Second):
					t.Fatal("timed out")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestRealTimeAndSetupTeardownMetrics(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
	import ***REMOVED*** Counter ***REMOVED*** from "k6/metrics";
	import ***REMOVED*** sleep ***REMOVED*** from "k6";

	var counter = new Counter("test_counter");

	export function setup() ***REMOVED***
		console.log("setup(), sleeping for 1 second");
		counter.add(1, ***REMOVED*** place: "setupBeforeSleep" ***REMOVED***);
		sleep(1);
		console.log("setup sleep is done");
		counter.add(2, ***REMOVED*** place: "setupAfterSleep" ***REMOVED***);
		return ***REMOVED*** "some": ["data"], "v": 1 ***REMOVED***;
	***REMOVED***

	export function teardown(data) ***REMOVED***
		console.log("teardown(" + JSON.stringify(data) + "), sleeping for 1 second");
		counter.add(3, ***REMOVED*** place: "teardownBeforeSleep" ***REMOVED***);
		sleep(1);
		if (!data || data.v != 1) ***REMOVED***
			throw new Error("incorrect data: " + JSON.stringify(data));
		***REMOVED***
		console.log("teardown sleep is done");
		counter.add(4, ***REMOVED*** place: "teardownAfterSleep" ***REMOVED***);
	***REMOVED***

	export default function (data) ***REMOVED***
		console.log("default(" + JSON.stringify(data) + ") with ENV=" + JSON.stringify(__ENV) + " for in ITER " + __ITER + " and VU " + __VU);
		counter.add(5, ***REMOVED*** place: "defaultBeforeSleep" ***REMOVED***);
		if (!data || data.v != 1) ***REMOVED***
			throw new Error("incorrect data: " + JSON.stringify(data));
		***REMOVED***
		sleep(1);
		console.log("default() for in ITER " + __ITER + " and VU " + __VU + " done!");
		counter.add(6, ***REMOVED*** place: "defaultAfterSleep" ***REMOVED***);
	***REMOVED***`)

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	runner, err := js.New(logger, &loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: script***REMOVED***, nil, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	options, err := executor.DeriveScenariosFromShortcuts(runner.GetOptions().Apply(lib.Options***REMOVED***
		Iterations:      null.IntFrom(2),
		VUs:             null.IntFrom(1),
		SystemTags:      &stats.DefaultSystemTagSet,
		SetupTimeout:    types.NullDurationFrom(4 * time.Second),
		TeardownTimeout: types.NullDurationFrom(4 * time.Second),
	***REMOVED***))
	require.NoError(t, err)
	require.NoError(t, runner.SetOptions(options))

	execScheduler, err := NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct***REMOVED******REMOVED***)
	sampleContainers := make(chan stats.SampleContainer)
	go func() ***REMOVED***
		require.NoError(t, execScheduler.Init(ctx, sampleContainers))
		assert.NoError(t, execScheduler.Run(ctx, ctx, sampleContainers))
		close(done)
	***REMOVED***()

	expectIn := func(from, to time.Duration, expected stats.SampleContainer) ***REMOVED***
		start := time.Now()
		from *= time.Millisecond
		to *= time.Millisecond
		for ***REMOVED***
			select ***REMOVED***
			case sampleContainer := <-sampleContainers:
				now := time.Now()
				elapsed := now.Sub(start)
				if elapsed < from ***REMOVED***
					t.Errorf("Received sample earlier (%s) than expected (%s)", elapsed, from)
					return
				***REMOVED***
				assert.IsType(t, expected, sampleContainer)
				expSamples := expected.GetSamples()
				gotSamples := sampleContainer.GetSamples()
				if assert.Len(t, gotSamples, len(expSamples)) ***REMOVED***
					for i, s := range gotSamples ***REMOVED***
						expS := expSamples[i]
						if s.Metric != metrics.IterationDuration ***REMOVED***
							assert.Equal(t, expS.Value, s.Value)
						***REMOVED***
						assert.Equal(t, expS.Metric.Name, s.Metric.Name)
						assert.Equal(t, expS.Tags.CloneTags(), s.Tags.CloneTags())
						assert.InDelta(t, 0, now.Sub(s.Time), float64(50*time.Millisecond))
					***REMOVED***
				***REMOVED***
				return
			case <-time.After(to):
				t.Errorf("Did not receive sample in the maximum allotted time (%s)", to)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	getTags := func(args ...string) *stats.SampleTags ***REMOVED***
		tags := map[string]string***REMOVED******REMOVED***
		for i := 0; i < len(args)-1; i += 2 ***REMOVED***
			tags[args[i]] = args[i+1]
		***REMOVED***
		return stats.IntoSampleTags(&tags)
	***REMOVED***
	testCounter := stats.New("test_counter", stats.Counter)
	getSample := func(expValue float64, expMetric *stats.Metric, expTags ...string) stats.SampleContainer ***REMOVED***
		return stats.Sample***REMOVED***
			Metric: expMetric,
			Time:   time.Now(),
			Tags:   getTags(expTags...),
			Value:  expValue,
		***REMOVED***
	***REMOVED***
	getDummyTrail := func(group string, emitIterations bool, addExpTags ...string) stats.SampleContainer ***REMOVED***
		expTags := []string***REMOVED***"group", group***REMOVED***
		expTags = append(expTags, addExpTags...)
		return netext.NewDialer(
			net.Dialer***REMOVED******REMOVED***,
			netext.NewResolver(net.LookupIP, 0, types.DNSfirst, types.DNSpreferIPv4),
		).GetTrail(time.Now(), time.Now(),
			true, emitIterations, getTags(expTags...))
	***REMOVED***

	// Initially give a long time (5s) for the execScheduler to start
	expectIn(0, 5000, getSample(1, testCounter, "group", "::setup", "place", "setupBeforeSleep"))
	expectIn(900, 1100, getSample(2, testCounter, "group", "::setup", "place", "setupAfterSleep"))
	expectIn(0, 100, getDummyTrail("::setup", false))

	expectIn(0, 100, getSample(5, testCounter, "group", "", "place", "defaultBeforeSleep", "scenario", "default"))
	expectIn(900, 1100, getSample(6, testCounter, "group", "", "place", "defaultAfterSleep", "scenario", "default"))
	expectIn(0, 100, getDummyTrail("", true, "scenario", "default"))

	expectIn(0, 100, getSample(5, testCounter, "group", "", "place", "defaultBeforeSleep", "scenario", "default"))
	expectIn(900, 1100, getSample(6, testCounter, "group", "", "place", "defaultAfterSleep", "scenario", "default"))
	expectIn(0, 100, getDummyTrail("", true, "scenario", "default"))

	expectIn(0, 1000, getSample(3, testCounter, "group", "::teardown", "place", "teardownBeforeSleep"))
	expectIn(900, 1100, getSample(4, testCounter, "group", "::teardown", "place", "teardownAfterSleep"))
	expectIn(0, 100, getDummyTrail("::teardown", false))

	for ***REMOVED***
		select ***REMOVED***
		case s := <-sampleContainers:
			t.Fatalf("Did not expect anything in the sample channel bug got %#v", s)
		case <-time.After(3 * time.Second):
			t.Fatalf("Local execScheduler took way to long to finish")
		case <-done:
			return // Exit normally
		***REMOVED***
	***REMOVED***
***REMOVED***

// Just a lib.PausableExecutor implementation that can return an error
type pausableExecutor struct ***REMOVED***
	lib.Executor
	err error
***REMOVED***

func (p pausableExecutor) SetPaused(bool) error ***REMOVED***
	return p.err
***REMOVED***

func TestSetPaused(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("second pause is an error", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		sched, err := NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: nil***REMOVED******REMOVED***

		require.NoError(t, sched.SetPaused(true))
		err = sched.SetPaused(true)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution is already paused")
	***REMOVED***)

	t.Run("unpause at the start is an error", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		sched, err := NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: nil***REMOVED******REMOVED***
		err = sched.SetPaused(false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution wasn't paused")
	***REMOVED***)

	t.Run("second unpause is an error", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		sched, err := NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: nil***REMOVED******REMOVED***
		require.NoError(t, sched.SetPaused(true))
		require.NoError(t, sched.SetPaused(false))
		err = sched.SetPaused(false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution wasn't paused")
	***REMOVED***)

	t.Run("an error on pausing is propagated", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED******REMOVED***
		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		sched, err := NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		expectedErr := errors.New("testing pausable executor error")
		sched.executors = []lib.Executor***REMOVED***pausableExecutor***REMOVED***err: expectedErr***REMOVED******REMOVED***
		err = sched.SetPaused(true)
		require.Error(t, err)
		require.Equal(t, err, expectedErr)
	***REMOVED***)

	t.Run("can't pause unpausable executor", func(t *testing.T) ***REMOVED***
		t.Parallel()
		runner := &minirunner.MiniRunner***REMOVED******REMOVED***
		options, err := executor.DeriveScenariosFromShortcuts(lib.Options***REMOVED***
			Iterations: null.IntFrom(2),
			VUs:        null.IntFrom(1),
		***REMOVED***.Apply(runner.GetOptions()))
		require.NoError(t, err)
		require.NoError(t, runner.SetOptions(options))

		logger := logrus.New()
		logger.SetOutput(testutils.NewTestOutput(t))
		sched, err := NewExecutionScheduler(runner, logger)
		require.NoError(t, err)
		err = sched.SetPaused(true)
		require.Error(t, err)
		require.Contains(t, err.Error(), "doesn't support pause and resume operations after its start")
	***REMOVED***)
***REMOVED***

func TestNewExecutionSchedulerHasWork(t *testing.T) ***REMOVED***
	t.Parallel()
	script := []byte(`
		import http from 'k6/http';

		export let options = ***REMOVED***
			executionSegment: "3/4:1",
			executionSegmentSequence: "0,1/4,2/4,3/4,1",
			scenarios: ***REMOVED***
				shared_iters1: ***REMOVED***
					executor: "shared-iterations",
					vus: 3,
					iterations: 3,
				***REMOVED***,
				shared_iters2: ***REMOVED***
					executor: "shared-iterations",
					vus: 4,
					iterations: 4,
				***REMOVED***,
				constant_arr_rate: ***REMOVED***
					executor: "constant-arrival-rate",
					rate: 3,
					timeUnit: "1s",
					duration: "20s",
					preAllocatedVUs: 4,
					maxVUs: 4,
				***REMOVED***,
		    ***REMOVED***,
		***REMOVED***;

		export default function() ***REMOVED***
			const response = http.get("http://test.loadimpact.com");
		***REMOVED***;
`)

	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))

	runner, err := js.New(
		logger,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
			Data: script,
		***REMOVED***,
		nil,
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)
	require.NoError(t, err)

	execScheduler, err := NewExecutionScheduler(runner, logger)
	require.NoError(t, err)

	assert.Len(t, execScheduler.executors, 2)
	assert.Len(t, execScheduler.executorConfigs, 3)
***REMOVED***
