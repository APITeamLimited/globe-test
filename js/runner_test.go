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

package js

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go/build"
	"io/ioutil"
	stdlog "log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/modules/k6"
	k6http "github.com/loadimpact/k6/js/modules/k6/http"
	k6metrics "github.com/loadimpact/k6/js/modules/k6/metrics"
	"github.com/loadimpact/k6/js/modules/k6/ws"
	"github.com/loadimpact/k6/lib"
	_ "github.com/loadimpact/k6/lib/executor" // TODO: figure out something better
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/dummy"
)

func TestRunnerNew(t *testing.T) ***REMOVED***
	t.Run("Valid", func(t *testing.T) ***REMOVED***
		r, err := getSimpleRunner("/script.js", `
			var counter = 0;
			exports.default = function() ***REMOVED*** counter++; ***REMOVED***
		`)
		assert.NoError(t, err)

		t.Run("NewVU", func(t *testing.T) ***REMOVED***
			initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
			assert.NoError(t, err)
			vuc, ok := initVU.(*VU)
			assert.True(t, ok)
			assert.Equal(t, int64(0), vuc.Runtime.Get("counter").Export())

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			t.Run("RunOnce", func(t *testing.T) ***REMOVED***
				err = vu.RunOnce()
				assert.NoError(t, err)
				assert.Equal(t, int64(1), vuc.Runtime.Get("counter").Export())
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := getSimpleRunner("/script.js", `blarg`)
		assert.EqualError(t, err, "ReferenceError: blarg is not defined at file:///script.js:1:1(0)")
	***REMOVED***)
***REMOVED***

func TestRunnerGetDefaultGroup(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `exports.default = function() ***REMOVED******REMOVED***;`)
	if assert.NoError(t, err) ***REMOVED***
		assert.NotNil(t, r1.GetDefaultGroup())
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if assert.NoError(t, err) ***REMOVED***
		assert.NotNil(t, r2.GetDefaultGroup())
	***REMOVED***
***REMOVED***

func TestRunnerOptions(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `exports.default = function() ***REMOVED******REMOVED***;`)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(false, false), r.Bundle.Options.Paused)
			r.SetOptions(lib.Options***REMOVED***Paused: null.BoolFrom(true)***REMOVED***)
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(true, true), r.Bundle.Options.Paused)
			r.SetOptions(lib.Options***REMOVED***Paused: null.BoolFrom(false)***REMOVED***)
			assert.Equal(t, r.Bundle.Options, r.GetOptions())
			assert.Equal(t, null.NewBool(false, true), r.Bundle.Options.Paused)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestOptionsSettingToScript(t *testing.T) ***REMOVED***
	t.Parallel()

	optionVariants := []string***REMOVED***
		"",
		"var options = null;",
		"var options = undefined;",
		"var options = ***REMOVED******REMOVED***;",
		"var options = ***REMOVED***teardownTimeout: '1s'***REMOVED***;",
	***REMOVED***

	for i, variant := range optionVariants ***REMOVED***
		variant := variant
		t.Run(fmt.Sprintf("Variant#%d", i), func(t *testing.T) ***REMOVED***
			t.Parallel()
			data := variant + `
					exports.default = function() ***REMOVED***
						if (!options) ***REMOVED***
							throw new Error("Expected options to be defined!");
						***REMOVED***
						if (options.teardownTimeout != __ENV.expectedTeardownTimeout) ***REMOVED***
							throw new Error("expected teardownTimeout to be " + __ENV.expectedTeardownTimeout + " but it was " + options.teardownTimeout);
						***REMOVED***
					***REMOVED***;`
			r, err := getSimpleRunner("/script.js", data,
				lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"expectedTeardownTimeout": "4s"***REMOVED******REMOVED***)
			require.NoError(t, err)

			newOptions := lib.Options***REMOVED***TeardownTimeout: types.NullDurationFrom(4 * time.Second)***REMOVED***
			r.SetOptions(newOptions)
			require.Equal(t, newOptions, r.GetOptions())

			samples := make(chan stats.SampleContainer, 100)
			initVU, err := r.NewVU(1, samples)
			if assert.NoError(t, err) ***REMOVED***
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
				err := vu.RunOnce()
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestOptionsPropagationToScript(t *testing.T) ***REMOVED***
	t.Parallel()
	data := `
			var options = ***REMOVED*** setupTimeout: "1s", myOption: "test" ***REMOVED***;
			exports.options = options;
			exports.default = function() ***REMOVED***
				if (options.external) ***REMOVED***
					throw new Error("Unexpected property external!");
				***REMOVED***
				if (options.myOption != "test") ***REMOVED***
					throw new Error("expected myOption to remain unchanged but it was '" + options.myOption + "'");
				***REMOVED***
				if (options.setupTimeout != __ENV.expectedSetupTimeout) ***REMOVED***
					throw new Error("expected setupTimeout to be " + __ENV.expectedSetupTimeout + " but it was " + options.setupTimeout);
				***REMOVED***
			***REMOVED***;`

	expScriptOptions := lib.Options***REMOVED***SetupTimeout: types.NullDurationFrom(1 * time.Second)***REMOVED***
	r1, err := getSimpleRunner("/script.js", data,
		lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"expectedSetupTimeout": "1s"***REMOVED******REMOVED***)
	require.NoError(t, err)
	require.Equal(t, expScriptOptions, r1.GetOptions())

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"expectedSetupTimeout": "3s"***REMOVED******REMOVED***)

	require.NoError(t, err)
	require.Equal(t, expScriptOptions, r2.GetOptions())

	newOptions := lib.Options***REMOVED***SetupTimeout: types.NullDurationFrom(3 * time.Second)***REMOVED***
	require.NoError(t, r2.SetOptions(newOptions))
	require.Equal(t, newOptions, r2.GetOptions())

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			samples := make(chan stats.SampleContainer, 100)

			initVU, err := r.NewVU(1, samples)
			if assert.NoError(t, err) ***REMOVED***
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
				err := vu.RunOnce()
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMetricName(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	script := tb.Replacer.Replace(`
		var Counter = require("k6/metrics").Counter;

		var myCounter = new Counter("not ok name @");

		exports.default = function(data) ***REMOVED***
			myCounter.add(1);
		***REMOVED***
	`)

	_, err := getSimpleRunner("/script.js", script)
	require.Error(t, err)
***REMOVED***

func TestSetupDataIsolation(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	script := tb.Replacer.Replace(`
		var Counter = require("k6/metrics").Counter;

		exports.options = ***REMOVED***
			scenarios: ***REMOVED***
				shared_iters: ***REMOVED***
					executor: "shared-iterations",
					vus: 5,
					iterations: 500,
				***REMOVED***,
			***REMOVED***,
			teardownTimeout: "5s",
			setupTimeout: "5s",
		***REMOVED***;
		var myCounter = new Counter("mycounter");

		exports.setup = function() ***REMOVED***
			return ***REMOVED*** v: 0 ***REMOVED***;
		***REMOVED***

		exports.default = function(data) ***REMOVED***
			if (data.v !== __ITER) ***REMOVED***
				throw new Error("default: wrong data for iter " + __ITER + ": " + JSON.stringify(data));
			***REMOVED***
			data.v += 1;
			myCounter.add(1);
		***REMOVED***

		exports.teardown = function(data) ***REMOVED***
			if (data.v !== 0) ***REMOVED***
				throw new Error("teardown: wrong data: " + data.v);
			***REMOVED***
			myCounter.add(1);
		***REMOVED***
	`)

	runner, err := getSimpleRunner("/script.js", script)
	require.NoError(t, err)

	options := runner.GetOptions()
	require.Empty(t, options.Validate())

	execScheduler, err := local.NewExecutionScheduler(runner, logrus.StandardLogger())
	require.NoError(t, err)
	engine, err := core.NewEngine(execScheduler, options, logrus.StandardLogger())
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	run, wait, err := engine.Init(ctx, ctx)
	require.NoError(t, err)

	collector := &dummy.Collector***REMOVED******REMOVED***
	engine.Collectors = []lib.Collector***REMOVED***collector***REMOVED***

	errC := make(chan error)
	go func() ***REMOVED*** errC <- run() ***REMOVED***()

	select ***REMOVED***
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("Test timed out")
	case err := <-errC:
		cancel()
		require.NoError(t, err)
		wait()
		require.False(t, engine.IsTainted())
	***REMOVED***
	var count int
	for _, s := range collector.Samples ***REMOVED***
		if s.Metric.Name == "mycounter" ***REMOVED***
			count += int(s.Value)
		***REMOVED***
	***REMOVED***
	require.Equal(t, 501, count, "mycounter should be the number of iterations + 1 for the teardown")
***REMOVED***

func testSetupDataHelper(t *testing.T, data string) ***REMOVED***
	t.Helper()
	expScriptOptions := lib.Options***REMOVED***
		SetupTimeout:    types.NullDurationFrom(1 * time.Second),
		TeardownTimeout: types.NullDurationFrom(1 * time.Second),
	***REMOVED***
	r1, err := getSimpleRunner("/script.js", data) // TODO fix this
	require.NoError(t, err)
	require.Equal(t, expScriptOptions, r1.GetOptions())

	testdata := map[string]*Runner***REMOVED***"Source": r1***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			samples := make(chan stats.SampleContainer, 100)

			if !assert.NoError(t, r.Setup(context.Background(), samples)) ***REMOVED***
				return
			***REMOVED***
			initVU, err := r.NewVU(1, samples)
			if assert.NoError(t, err) ***REMOVED***
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
				err := vu.RunOnce()
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSetupDataReturnValue(t *testing.T) ***REMOVED***
	testSetupDataHelper(t, `
	exports.options = ***REMOVED*** setupTimeout: "1s", teardownTimeout: "1s" ***REMOVED***;
	exports.setup = function() ***REMOVED***
		return 42;
	***REMOVED***
	exports.default = function(data) ***REMOVED***
		if (data != 42) ***REMOVED***
			throw new Error("default: wrong data: " + JSON.stringify(data))
		***REMOVED***
	***REMOVED***;

	exports.teardown = function(data) ***REMOVED***
		if (data != 42) ***REMOVED***
			throw new Error("teardown: wrong data: " + JSON.stringify(data))
		***REMOVED***
	***REMOVED***;`)
***REMOVED***

func TestSetupDataNoSetup(t *testing.T) ***REMOVED***
	testSetupDataHelper(t, `
	exports.options = ***REMOVED*** setupTimeout: "1s", teardownTimeout: "1s" ***REMOVED***;
	exports.default = function(data) ***REMOVED***
		if (data !== undefined) ***REMOVED***
			throw new Error("default: wrong data: " + JSON.stringify(data))
		***REMOVED***
	***REMOVED***;

	exports.teardown = function(data) ***REMOVED***
		if (data !== undefined) ***REMOVED***
			console.log(data);
			throw new Error("teardown: wrong data: " + JSON.stringify(data))
		***REMOVED***
	***REMOVED***;`)
***REMOVED***

func TestConsoleInInitContext(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `
			console.log("1");
			exports.default = function(data) ***REMOVED***
			***REMOVED***;
		`)
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			samples := make(chan stats.SampleContainer, 100)
			initVU, err := r.NewVU(1, samples)
			if assert.NoError(t, err) ***REMOVED***
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
				err := vu.RunOnce()
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSetupDataNoReturn(t *testing.T) ***REMOVED***
	testSetupDataHelper(t, `
	exports.options = ***REMOVED*** setupTimeout: "1s", teardownTimeout: "1s" ***REMOVED***;
	exports.setup = function() ***REMOVED*** ***REMOVED***
	exports.default = function(data) ***REMOVED***
		if (data !== undefined) ***REMOVED***
			throw new Error("default: wrong data: " + JSON.stringify(data))
		***REMOVED***
	***REMOVED***;

	exports.teardown = function(data) ***REMOVED***
		if (data !== undefined) ***REMOVED***
			throw new Error("teardown: wrong data: " + JSON.stringify(data))
		***REMOVED***
	***REMOVED***;`)
***REMOVED***

func TestRunnerIntegrationImports(t *testing.T) ***REMOVED***
	t.Run("Modules", func(t *testing.T) ***REMOVED***
		modules := []string***REMOVED***
			"k6",
			"k6/http",
			"k6/metrics",
			"k6/html",
		***REMOVED***
		rtOpts := lib.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom("extended")***REMOVED***
		for _, mod := range modules ***REMOVED***
			mod := mod
			t.Run(mod, func(t *testing.T) ***REMOVED***
				t.Run("Source", func(t *testing.T) ***REMOVED***
					_, err := getSimpleRunner("/script.js", fmt.Sprintf(`import "%s"; exports.default = function() ***REMOVED******REMOVED***`, mod), rtOpts)
					assert.NoError(t, err)
				***REMOVED***)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)

	t.Run("Files", func(t *testing.T) ***REMOVED***
		fs := afero.NewMemMapFs()
		require.NoError(t, fs.MkdirAll("/path/to", 0755))
		require.NoError(t, afero.WriteFile(fs, "/path/to/lib.js", []byte(`exports.default = "hi!";`), 0644))

		testdata := map[string]struct***REMOVED*** filename, path string ***REMOVED******REMOVED***
			"Absolute":       ***REMOVED***"/path/script.js", "/path/to/lib.js"***REMOVED***,
			"Relative":       ***REMOVED***"/path/script.js", "./to/lib.js"***REMOVED***,
			"Adjacent":       ***REMOVED***"/path/to/script.js", "./lib.js"***REMOVED***,
			"STDIN-Absolute": ***REMOVED***"-", "/path/to/lib.js"***REMOVED***,
			"STDIN-Relative": ***REMOVED***"-", "./path/to/lib.js"***REMOVED***,
		***REMOVED***
		for name, data := range testdata ***REMOVED***
			name, data := name, data
			t.Run(name, func(t *testing.T) ***REMOVED***
				r1, err := getSimpleRunner(data.filename, fmt.Sprintf(`
					var hi = require("%s").default;
					exports.default = function() ***REMOVED***
						if (hi != "hi!") ***REMOVED*** throw new Error("incorrect value"); ***REMOVED***
					***REMOVED***`, data.path), fs)
				require.NoError(t, err)

				r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
				require.NoError(t, err)

				testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
				for name, r := range testdata ***REMOVED***
					r := r
					t.Run(name, func(t *testing.T) ***REMOVED***
						initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
						require.NoError(t, err)
						ctx, cancel := context.WithCancel(context.Background())
						defer cancel()
						vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
						err = vu.RunOnce()
						require.NoError(t, err)
					***REMOVED***)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestVURunContext(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `
		exports.options = ***REMOVED*** vus: 10 ***REMOVED***;
		exports.default = function() ***REMOVED*** fn(); ***REMOVED***
		`)
	require.NoError(t, err)
	r1.SetOptions(r1.GetOptions().Apply(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***))

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.newVU(1, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			fnCalled := false
			vu.Runtime.Set("fn", func() ***REMOVED***
				fnCalled = true

				assert.Equal(t, vu.Runtime, common.GetRuntime(*vu.Context), "incorrect runtime in context")

				state := lib.GetState(*vu.Context)
				if assert.NotNil(t, state) ***REMOVED***
					assert.Equal(t, null.IntFrom(10), state.Options.VUs)
					assert.Equal(t, null.BoolFrom(true), state.Options.Throw)
					assert.NotNil(t, state.Logger)
					assert.Equal(t, r.GetDefaultGroup(), state.Group)
					assert.Equal(t, vu.Transport, state.Transport)
				***REMOVED***
			***REMOVED***)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			assert.NoError(t, err)
			assert.True(t, fnCalled, "fn() not called")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVURunInterrupt(t *testing.T) ***REMOVED***
	// TODO: figure out why interrupt sometimes fails... data race in goja?
	if isWindows ***REMOVED***
		t.Skip()
	***REMOVED***

	r1, err := getSimpleRunner("/script.js", `
		exports.default = function() ***REMOVED*** while(true) ***REMOVED******REMOVED*** ***REMOVED***
		`)
	require.NoError(t, err)
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***))

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)
	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		name, r := name, r
		t.Run(name, func(t *testing.T) ***REMOVED***
			samples := make(chan stats.SampleContainer, 100)
			defer close(samples)
			go func() ***REMOVED***
				for range samples ***REMOVED***
				***REMOVED***
			***REMOVED***()

			vu, err := r.newVU(1, samples)
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "context cancelled")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVURunInterruptDoesntPanic(t *testing.T) ***REMOVED***
	// TODO: figure out why interrupt sometimes fails... data race in goja?
	if isWindows ***REMOVED***
		t.Skip()
	***REMOVED***

	r1, err := getSimpleRunner("/script.js", `
		exports.default = function() ***REMOVED*** while(true) ***REMOVED******REMOVED*** ***REMOVED***
		`)
	require.NoError(t, err)
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***))

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)
	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		name, r := name, r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			samples := make(chan stats.SampleContainer, 100)
			go func() ***REMOVED***
				for range samples ***REMOVED***
				***REMOVED***
			***REMOVED***()
			var wg sync.WaitGroup

			initVU, err := r.newVU(1, samples)
			require.NoError(t, err)
			for i := 0; i < 1000; i++ ***REMOVED***
				wg.Add(1)
				newCtx, newCancel := context.WithCancel(ctx)
				vu := initVU.Activate(&lib.VUActivationParams***REMOVED***
					RunContext:         newCtx,
					DeactivateCallback: func(_ lib.InitializedVU) ***REMOVED*** wg.Done() ***REMOVED***,
				***REMOVED***)
				ch := make(chan struct***REMOVED******REMOVED***)
				go func() ***REMOVED***
					close(ch)
					vuErr := vu.RunOnce()
					assert.Error(t, vuErr)
					assert.Contains(t, vuErr.Error(), "context cancelled")
				***REMOVED***()
				<-ch
				time.Sleep(time.Millisecond * 1) // NOTE: increase this in case of problems ;)
				newCancel()
				wg.Wait()
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationGroups(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `
		var group = require("k6").group;
		exports.default = function() ***REMOVED***
			fnOuter();
			group("my group", function() ***REMOVED***
				fnInner();
				group("nested group", function() ***REMOVED***
					fnNested();
				***REMOVED***)
			***REMOVED***);
		***REMOVED***
		`)
	require.NoError(t, err)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			vu, err := r.newVU(1, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			fnOuterCalled := false
			fnInnerCalled := false
			fnNestedCalled := false
			vu.Runtime.Set("fnOuter", func() ***REMOVED***
				fnOuterCalled = true
				assert.Equal(t, r.GetDefaultGroup(), lib.GetState(*vu.Context).Group)
			***REMOVED***)
			vu.Runtime.Set("fnInner", func() ***REMOVED***
				fnInnerCalled = true
				g := lib.GetState(*vu.Context).Group
				assert.Equal(t, "my group", g.Name)
				assert.Equal(t, r.GetDefaultGroup(), g.Parent)
			***REMOVED***)
			vu.Runtime.Set("fnNested", func() ***REMOVED***
				fnNestedCalled = true
				g := lib.GetState(*vu.Context).Group
				assert.Equal(t, "nested group", g.Name)
				assert.Equal(t, "my group", g.Parent.Name)
				assert.Equal(t, r.GetDefaultGroup(), g.Parent.Parent)
			***REMOVED***)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			assert.NoError(t, err)
			assert.True(t, fnOuterCalled, "fnOuter() not called")
			assert.True(t, fnInnerCalled, "fnInner() not called")
			assert.True(t, fnNestedCalled, "fnNested() not called")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationMetrics(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `
		var group = require("k6").group;
		var Trend = require("k6/metrics").Trend;
		var myMetric = new Trend("my_metric");
		exports.default = function() ***REMOVED*** myMetric.add(5); ***REMOVED***
		`)
	require.NoError(t, err)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			samples := make(chan stats.SampleContainer, 100)
			vu, err := r.newVU(1, samples)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			assert.NoError(t, err)
			sampleCount := 0
			for i, sampleC := range stats.GetBufferedSamples(samples) ***REMOVED***
				for j, s := range sampleC.GetSamples() ***REMOVED***
					sampleCount++
					switch i + j ***REMOVED***
					case 0:
						assert.Equal(t, 5.0, s.Value)
						assert.Equal(t, "my_metric", s.Metric.Name)
						assert.Equal(t, stats.Trend, s.Metric.Type)
					case 1:
						assert.Equal(t, 0.0, s.Value)
						assert.Equal(t, metrics.DataSent, s.Metric, "`data_sent` sample is before `data_received` and `iteration_duration`")
					case 2:
						assert.Equal(t, 0.0, s.Value)
						assert.Equal(t, metrics.DataReceived, s.Metric, "`data_received` sample is after `data_received`")
					case 3:
						assert.Equal(t, metrics.IterationDuration, s.Metric, "`iteration-duration` sample is after `data_received`")
					case 4:
						assert.Equal(t, metrics.Iterations, s.Metric, "`iterations` sample is after `iteration_duration`")
						assert.Equal(t, float64(1), s.Value)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			assert.Equal(t, sampleCount, 5)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationInsecureRequests(t *testing.T) ***REMOVED***
	testdata := map[string]struct ***REMOVED***
		opts   lib.Options
		errMsg string
	***REMOVED******REMOVED***
		"Null": ***REMOVED***
			lib.Options***REMOVED******REMOVED***,
			"x509: certificate has expired or is not yet valid",
		***REMOVED***,
		"False": ***REMOVED***
			lib.Options***REMOVED***InsecureSkipTLSVerify: null.BoolFrom(false)***REMOVED***,
			"x509: certificate has expired or is not yet valid",
		***REMOVED***,
		"True": ***REMOVED***
			lib.Options***REMOVED***InsecureSkipTLSVerify: null.BoolFrom(true)***REMOVED***,
			"",
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		data := data
		t.Run(name, func(t *testing.T) ***REMOVED***
			r1, err := getSimpleRunner("/script.js", `
					var http = require("k6/http");;
					exports.default = function() ***REMOVED*** http.get("https://expired.badssl.com/"); ***REMOVED***
				`)
			require.NoError(t, err)
			require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***.Apply(data.opts)))

			r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
			require.NoError(t, err)
			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				r := r
				t.Run(name, func(t *testing.T) ***REMOVED***
					r.Logger, _ = logtest.NewNullLogger()

					initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					err = vu.RunOnce()
					if data.errMsg != "" ***REMOVED***
						require.Error(t, err)
						assert.Contains(t, err.Error(), data.errMsg)
					***REMOVED*** else ***REMOVED***
						assert.NoError(t, err)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationBlacklistOption(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `
					var http = require("k6/http");;
					exports.default = function() ***REMOVED*** http.get("http://10.1.2.3/"); ***REMOVED***
				`)
	require.NoError(t, err)

	cidr, err := lib.ParseCIDR("10.0.0.0/8")

	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***
		Throw:        null.BoolFrom(true),
		BlacklistIPs: []*lib.IPNet***REMOVED***cidr***REMOVED***,
	***REMOVED***))

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "IP (10.1.2.3) is in a blacklisted range (10.0.0.0/8)")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationBlacklistScript(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `
					var http = require("k6/http");;

					exports.options = ***REMOVED***
						throw: true,
						blacklistIPs: ["10.0.0.0/8"],
					***REMOVED***;

					exports.default = function() ***REMOVED*** http.get("http://10.1.2.3/"); ***REMOVED***
				`)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***

	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "IP (10.1.2.3) is in a blacklisted range (10.0.0.0/8)")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationHosts(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	r1, err := getSimpleRunner("/script.js",
		tb.Replacer.Replace(`
					var k6 = require("k6");
					var check = k6.check;
					var fail = k6.fail;
					var http = require("k6/http");;
					exports.default = function() ***REMOVED***
						var res = http.get("http://test.loadimpact.com:HTTPBIN_PORT/");
						check(res, ***REMOVED***
							"is correct IP": function(r) ***REMOVED*** return r.remote_ip === "127.0.0.1" ***REMOVED***
						***REMOVED***) || fail("failed to override dns");
					***REMOVED***
				`))
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	r1.SetOptions(lib.Options***REMOVED***
		Throw: null.BoolFrom(true),
		Hosts: map[string]net.IP***REMOVED***
			"test.loadimpact.com": net.ParseIP("127.0.0.1"),
		***REMOVED***,
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationTLSConfig(t *testing.T) ***REMOVED***
	unsupportedVersionErrorMsg := "remote error: tls: handshake failure"
	for _, tag := range build.Default.ReleaseTags ***REMOVED***
		if tag == "go1.12" ***REMOVED***
			unsupportedVersionErrorMsg = "tls: no supported versions satisfy MinVersion and MaxVersion"
			break
		***REMOVED***
	***REMOVED***
	testdata := map[string]struct ***REMOVED***
		opts   lib.Options
		errMsg string
	***REMOVED******REMOVED***
		"NullCipherSuites": ***REMOVED***
			lib.Options***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"SupportedCipherSuite": ***REMOVED***
			lib.Options***REMOVED***TLSCipherSuites: &lib.TLSCipherSuites***REMOVED***tls.TLS_RSA_WITH_AES_128_GCM_SHA256***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"UnsupportedCipherSuite": ***REMOVED***
			lib.Options***REMOVED***TLSCipherSuites: &lib.TLSCipherSuites***REMOVED***tls.TLS_RSA_WITH_RC4_128_SHA***REMOVED******REMOVED***,
			"remote error: tls: handshake failure",
		***REMOVED***,
		"NullVersion": ***REMOVED***
			lib.Options***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"SupportedVersion": ***REMOVED***
			lib.Options***REMOVED***TLSVersion: &lib.TLSVersions***REMOVED***Min: tls.VersionTLS12, Max: tls.VersionTLS12***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		"UnsupportedVersion": ***REMOVED***
			lib.Options***REMOVED***TLSVersion: &lib.TLSVersions***REMOVED***Min: tls.VersionSSL30, Max: tls.VersionSSL30***REMOVED******REMOVED***,
			unsupportedVersionErrorMsg,
		***REMOVED***,
	***REMOVED***
	for name, data := range testdata ***REMOVED***
		data := data
		t.Run(name, func(t *testing.T) ***REMOVED***
			r1, err := getSimpleRunner("/script.js", `
					var http = require("k6/http");;
					exports.default = function() ***REMOVED*** http.get("https://sha256.badssl.com/"); ***REMOVED***
				`)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***.Apply(data.opts)))

			r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				r := r
				t.Run(name, func(t *testing.T) ***REMOVED***
					r.Logger, _ = logtest.NewNullLogger()

					initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
					if !assert.NoError(t, err) ***REMOVED***
						return
					***REMOVED***
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					err = vu.RunOnce()
					if data.errMsg != "" ***REMOVED***
						require.Error(t, err)
						assert.Contains(t, err.Error(), data.errMsg)
					***REMOVED*** else ***REMOVED***
						assert.NoError(t, err)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationOpenFunctionError(t *testing.T) ***REMOVED***
	r, err := getSimpleRunner("/script.js", `
			exports.default = function() ***REMOVED*** open("/tmp/foo") ***REMOVED***
		`)
	assert.NoError(t, err)

	initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	err = vu.RunOnce()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only available in the init stage")
***REMOVED***

func TestVUIntegrationOpenFunctionErrorWhenSneaky(t *testing.T) ***REMOVED***
	r, err := getSimpleRunner("/script.js", `
			var sneaky = open;
			exports.default = function() ***REMOVED*** sneaky("/tmp/foo") ***REMOVED***
		`)
	assert.NoError(t, err)

	initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	err = vu.RunOnce()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only available in the init stage")
***REMOVED***

func TestVUIntegrationCookiesReset(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	r1, err := getSimpleRunner("/script.js", tb.Replacer.Replace(`
			var http = require("k6/http");;
			exports.default = function() ***REMOVED***
				var url = "HTTPBIN_URL";
				var preRes = http.get(url + "/cookies");
				if (preRes.status != 200) ***REMOVED*** throw new Error("wrong status (pre): " + preRes.status); ***REMOVED***
				if (preRes.json().k1 || preRes.json().k2) ***REMOVED***
					throw new Error("cookies persisted: " + preRes.body);
				***REMOVED***

				var res = http.get(url + "/cookies/set?k2=v2&k1=v1");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
				if (res.json().k1 != "v1" || res.json().k2 != "v2") ***REMOVED***
					throw new Error("wrong cookies: " + res.body);
				***REMOVED***
			***REMOVED***
		`))
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***
		Throw:        null.BoolFrom(true),
		MaxRedirects: null.IntFrom(10),
		Hosts:        tb.Dialer.Hosts,
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			for i := 0; i < 2; i++ ***REMOVED***
				err = vu.RunOnce()
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationCookiesNoReset(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	r1, err := getSimpleRunner("/script.js", tb.Replacer.Replace(`
			var http = require("k6/http");;
			exports.default = function() ***REMOVED***
				var url = "HTTPBIN_URL";
				if (__ITER == 0) ***REMOVED***
					var res = http.get(url + "/cookies/set?k2=v2&k1=v1");
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
					if (res.json().k1 != "v1" || res.json().k2 != "v2") ***REMOVED***
						throw new Error("wrong cookies: " + res.body);
					***REMOVED***
				***REMOVED***

				if (__ITER == 1) ***REMOVED***
					var res = http.get(url + "/cookies");
					if (res.status != 200) ***REMOVED*** throw new Error("wrong status (pre): " + res.status); ***REMOVED***
					if (res.json().k1 != "v1" || res.json().k2 != "v2") ***REMOVED***
						throw new Error("wrong cookies: " + res.body);
					***REMOVED***
				***REMOVED***
			***REMOVED***
		`))
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***
		Throw:          null.BoolFrom(true),
		MaxRedirects:   null.IntFrom(10),
		Hosts:          tb.Dialer.Hosts,
		NoCookiesReset: null.BoolFrom(true),
	***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			assert.NoError(t, err)

			err = vu.RunOnce()
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationVUID(t *testing.T) ***REMOVED***
	r1, err := getSimpleRunner("/script.js", `
			exports.default = function() ***REMOVED***
				if (__VU != 1234) ***REMOVED*** throw new Error("wrong __VU: " + __VU); ***REMOVED***
			***REMOVED***`,
	)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)

	r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			initVU, err := r.NewVU(1234, make(chan stats.SampleContainer, 100))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationClientCerts(t *testing.T) ***REMOVED***
	clientCAPool := x509.NewCertPool()
	assert.True(t, clientCAPool.AppendCertsFromPEM(
		[]byte("-----BEGIN CERTIFICATE-----\n"+
			"MIIBYzCCAQqgAwIBAgIUMYw1pqZ1XhXdFG0S2ITXhfHBsWgwCgYIKoZIzj0EAwIw\n"+
			"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE1MTYxODAwWhcNMjIwODE0MTYxODAw\n"+
			"WjAQMQ4wDAYDVQQDEwVNeSBDQTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABFWO\n"+
			"fg4dgL8cdvjoSWDQFLBJxlbQFlZfOSyUR277a4g91BD07KWX+9ny+Q8WuUODog06\n"+
			"xH1g8fc6zuaejllfzM6jQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTAD\n"+
			"AQH/MB0GA1UdDgQWBBTeoSFylGCmyqj1X4sWez1r6hkhjDAKBggqhkjOPQQDAgNH\n"+
			"ADBEAiAfuKi6u/BVXenCkgnU2sfXsYjel6rACuXEcx01yaaWuQIgXAtjrDisdlf4\n"+
			"0ZdoIoYjNhDAXUtnyRBt+V6+rIklv/8=\n"+
			"-----END CERTIFICATE-----"),
	))
	serverCert, err := tls.X509KeyPair(
		[]byte("-----BEGIN CERTIFICATE-----\n"+
			"MIIBxjCCAW2gAwIBAgIUICcYHG1bI28NZm676wHlMPxL+CEwCgYIKoZIzj0EAwIw\n"+
			"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE3MTQwNjAwWhcNMTgwODE3MTQwNjAw\n"+
			"WjAZMRcwFQYDVQQDEw4xMjcuMC4wLjE6Njk2OTBZMBMGByqGSM49AgEGCCqGSM49\n"+
			"AwEHA0IABCdD1IqowucJ5oUjGYCZZnXvgi7EMD4jD1osbOkzOFFnHSLRvdm6fcJu\n"+
			"vPUcl4g8zUs466sC0AVUNpk21XbA/QajgZswgZgwDgYDVR0PAQH/BAQDAgWgMB0G\n"+
			"A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1Ud\n"+
			"DgQWBBTeAc8HY3sgGIV+fu/lY0OKr2Ho0jAfBgNVHSMEGDAWgBTeoSFylGCmyqj1\n"+
			"X4sWez1r6hkhjDAZBgNVHREEEjAQgg4xMjcuMC4wLjE6Njk2OTAKBggqhkjOPQQD\n"+
			"AgNHADBEAiAt3gC5FGQfSJXQ5DloXAOeJDFnKIL7d6xhftgPS5O08QIgRuAyysB8\n"+
			"5JXHvvze5DMN/clHYptos9idVFc+weUZAUQ=\n"+
			"-----END CERTIFICATE-----\n"+
			"-----BEGIN CERTIFICATE-----\n"+
			"MIIBYzCCAQqgAwIBAgIUMYw1pqZ1XhXdFG0S2ITXhfHBsWgwCgYIKoZIzj0EAwIw\n"+
			"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE1MTYxODAwWhcNMjIwODE0MTYxODAw\n"+
			"WjAQMQ4wDAYDVQQDEwVNeSBDQTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABFWO\n"+
			"fg4dgL8cdvjoSWDQFLBJxlbQFlZfOSyUR277a4g91BD07KWX+9ny+Q8WuUODog06\n"+
			"xH1g8fc6zuaejllfzM6jQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTAD\n"+
			"AQH/MB0GA1UdDgQWBBTeoSFylGCmyqj1X4sWez1r6hkhjDAKBggqhkjOPQQDAgNH\n"+
			"ADBEAiAfuKi6u/BVXenCkgnU2sfXsYjel6rACuXEcx01yaaWuQIgXAtjrDisdlf4\n"+
			"0ZdoIoYjNhDAXUtnyRBt+V6+rIklv/8=\n"+
			"-----END CERTIFICATE-----"),
		[]byte("-----BEGIN EC PRIVATE KEY-----\n"+
			"MHcCAQEEIKYptA4VtQ8UOKL+d1wkhl+51aPpvO+ppY62nLF9Z1w5oAoGCCqGSM49\n"+
			"AwEHoUQDQgAEJ0PUiqjC5wnmhSMZgJlmde+CLsQwPiMPWixs6TM4UWcdItG92bp9\n"+
			"wm689RyXiDzNSzjrqwLQBVQ2mTbVdsD9Bg==\n"+
			"-----END EC PRIVATE KEY-----"),
	)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config***REMOVED***
		Certificates: []tls.Certificate***REMOVED***serverCert***REMOVED***,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAPool,
	***REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	defer func() ***REMOVED*** _ = listener.Close() ***REMOVED***()
	srv := &http.Server***REMOVED***
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
			_, _ = fmt.Fprintf(w, "ok")
		***REMOVED***),
		ErrorLog: stdlog.New(ioutil.Discard, "", 0),
	***REMOVED***
	go func() ***REMOVED*** _ = srv.Serve(listener) ***REMOVED***()

	r1, err := getSimpleRunner("/script.js", fmt.Sprintf(`
			var http = require("k6/http");;
			exports.default = function() ***REMOVED*** http.get("https://%s")***REMOVED***
		`, listener.Addr().String()))
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***
		Throw:                 null.BoolFrom(true),
		InsecureSkipTLSVerify: null.BoolFrom(true),
	***REMOVED***))

	t.Run("Unauthenticated", func(t *testing.T) ***REMOVED***
		r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***

		runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
		for name, r := range runners ***REMOVED***
			r := r
			t.Run(name, func(t *testing.T) ***REMOVED***
				r.Logger, _ = logtest.NewNullLogger()
				initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
				if assert.NoError(t, err) ***REMOVED***
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					err := vu.RunOnce()
					require.Error(t, err)
					assert.Contains(t, err.Error(), "remote error: tls: bad certificate")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)

	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***
		TLSAuth: []*lib.TLSAuth***REMOVED***
			***REMOVED***
				TLSAuthFields: lib.TLSAuthFields***REMOVED***
					Domains: []string***REMOVED***"127.0.0.1"***REMOVED***,
					Cert: "-----BEGIN CERTIFICATE-----\n" +
						"MIIBoTCCAUigAwIBAgIUd6XedDxP+rGo+kq0APqHElGZzs4wCgYIKoZIzj0EAwIw\n" +
						"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE3MTUwNjAwWhcNMTgwODE3MTUwNjAw\n" +
						"WjARMQ8wDQYDVQQDEwZjbGllbnQwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATL\n" +
						"mi/a1RVvk05FyrYmartbo/9cW+53DrQLW1twurII2q5ZfimdMX05A32uB3Ycoy/J\n" +
						"x+w7Ifyd/YRw0zEc3NHQo38wfTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYI\n" +
						"KwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFN2SR/TD\n" +
						"yNW5DQWxZSkoXHQWsLY+MB8GA1UdIwQYMBaAFN6hIXKUYKbKqPVfixZ7PWvqGSGM\n" +
						"MAoGCCqGSM49BAMCA0cAMEQCICtETmyOmupmg4w3tw59VYJyOBqRTxg6SK+rOQmq\n" +
						"kE1VAiAUvsflDfmWBZ8EMPu46OhX6RX6MbvJ9NNvRco2G5ek1w==\n" +
						"-----END CERTIFICATE-----",
					Key: "-----BEGIN EC PRIVATE KEY-----\n" +
						"MHcCAQEEIOrnhT05alCeQEX66HgnSHah/m5LazjJHLDawYRnhUtZoAoGCCqGSM49\n" +
						"AwEHoUQDQgAEy5ov2tUVb5NORcq2Jmq7W6P/XFvudw60C1tbcLqyCNquWX4pnTF9\n" +
						"OQN9rgd2HKMvycfsOyH8nf2EcNMxHNzR0A==\n" +
						"-----END EC PRIVATE KEY-----",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***))

	t.Run("Authenticated", func(t *testing.T) ***REMOVED***
		r2, err := NewFromArchive(r1.MakeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***

		runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
		for name, r := range runners ***REMOVED***
			r := r
			t.Run(name, func(t *testing.T) ***REMOVED***
				initVU, err := r.NewVU(1, make(chan stats.SampleContainer, 100))
				if assert.NoError(t, err) ***REMOVED***
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					err := vu.RunOnce()
					assert.NoError(t, err)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestHTTPRequestInInitContext(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	_, err := getSimpleRunner("/script.js", tb.Replacer.Replace(`
					var k6 = require("k6");
					var check = k6.check;
					var fail = k6.fail;
					var http = require("k6/http");;
					var res = http.get("HTTPBIN_URL/");
					exports.default = function() ***REMOVED***
						console.log(test);
					***REMOVED***
				`))
	if assert.Error(t, err) ***REMOVED***
		assert.Equal(
			t,
			"GoError: "+k6http.ErrHTTPForbiddenInInitContext.Error(),
			err.Error())
	***REMOVED***
***REMOVED***

func TestInitContextForbidden(t *testing.T) ***REMOVED***
	table := [...][3]string***REMOVED***
		***REMOVED***
			"http.request",
			`var http = require("k6/http");;
			 var res = http.get("HTTPBIN_URL");
			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			k6http.ErrHTTPForbiddenInInitContext.Error(),
		***REMOVED***,
		***REMOVED***
			"http.batch",
			`var http = require("k6/http");;
			 var res = http.batch("HTTPBIN_URL/something", "HTTPBIN_URL/else");
			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			k6http.ErrBatchForbiddenInInitContext.Error(),
		***REMOVED***,
		***REMOVED***
			"http.cookieJar",
			`var http = require("k6/http");;
			 var jar = http.cookieJar();
			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			k6http.ErrJarForbiddenInInitContext.Error(),
		***REMOVED***,
		***REMOVED***
			"check",
			`var check = require("k6").check;
			 check("test", ***REMOVED***'is test': function(test) ***REMOVED*** return test == "test"***REMOVED******REMOVED***)
			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			k6.ErrCheckInInitContext.Error(),
		***REMOVED***,
		***REMOVED***
			"group",
			`var group = require("k6").group;
			 group("group1", function () ***REMOVED*** console.log("group1");***REMOVED***)
			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			k6.ErrGroupInInitContext.Error(),
		***REMOVED***,
		***REMOVED***
			"ws",
			`var ws = require("k6/ws");
			 var url = "ws://echo.websocket.org";
			 var params = ***REMOVED*** "tags": ***REMOVED*** "my_tag": "hello" ***REMOVED*** ***REMOVED***;
			 var response = ws.connect(url, params, function (socket) ***REMOVED***
			   socket.on('open', function open() ***REMOVED***
					console.log('connected');
			   ***REMOVED***)
		   ***REMOVED***);

			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			ws.ErrWSInInitContext.Error(),
		***REMOVED***,
		***REMOVED***
			"metric",
			`var Counter = require("k6/metrics").Counter;
			 var counter = Counter("myCounter");
			 counter.add(1);
			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			k6metrics.ErrMetricsAddInInitContext.Error(),
		***REMOVED***,
	***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	for _, test := range table ***REMOVED***
		test := test
		t.Run(test[0], func(t *testing.T) ***REMOVED***
			_, err := getSimpleRunner("/script.js", tb.Replacer.Replace(test[1]))
			if assert.Error(t, err) ***REMOVED***
				assert.Equal(
					t,
					"GoError: "+test[2],
					err.Error())
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestArchiveRunningIntegrity(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	fs := afero.NewMemMapFs()
	data := tb.Replacer.Replace(`
			var fput = open("/home/somebody/test.json");
			exports.options = ***REMOVED*** setupTimeout: "10s", teardownTimeout: "10s" ***REMOVED***;
			exports.setup = function () ***REMOVED***
				return JSON.parse(fput);
			***REMOVED***
			exports.default = function(data) ***REMOVED***
				if (data != 42) ***REMOVED***
					throw new Error("incorrect answer " + data);
				***REMOVED***
			***REMOVED***
		`)
	require.NoError(t, afero.WriteFile(fs, "/home/somebody/test.json", []byte(`42`), os.ModePerm))
	require.NoError(t, afero.WriteFile(fs, "/script.js", []byte(data), os.ModePerm))
	r1, err := getSimpleRunner("/script.js", data, fs)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, r1.MakeArchive().Write(buf))

	arc, err := lib.ReadArchive(buf)
	require.NoError(t, err)
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := make(chan stats.SampleContainer, 100)
			err = r.Setup(context.Background(), ch)
			require.NoError(t, err)
			initVU, err := r.NewVU(1, ch)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestArchiveNotPanicking(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/non/existent", []byte(`42`), os.ModePerm))
	r1, err := getSimpleRunner("/script.js", tb.Replacer.Replace(`
			var fput = open("/non/existent");
			exports.default = function(data) ***REMOVED******REMOVED***
		`), fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	arc.Filesystems = map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs()***REMOVED***
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	// we do want this to error here as this is where we find out that a given file is not in the
	// archive
	require.Error(t, err)
	require.Nil(t, r2)
***REMOVED***

func TestStuffNotPanicking(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	r, err := getSimpleRunner("/script.js", tb.Replacer.Replace(`
			var http = require("k6/http");
			var ws = require("k6/ws");
			var group = require("k6").group;
			var parseHTML = require("k6/html").parseHTML;

			exports.options = ***REMOVED*** iterations: 1, vus: 1, vusMax: 1 ***REMOVED***;

			exports.default = function() ***REMOVED***
				var doc = parseHTML(http.get("HTTPBIN_URL/html").body);

				var testCases = [
					function() ***REMOVED*** return group()***REMOVED***,
					function() ***REMOVED*** return group("test")***REMOVED***,
					function() ***REMOVED*** return group("test", "wat")***REMOVED***,
					function() ***REMOVED*** return doc.find('p').each()***REMOVED***,
					function() ***REMOVED*** return doc.find('p').each("wat")***REMOVED***,
					function() ***REMOVED*** return doc.find('p').map()***REMOVED***,
					function() ***REMOVED*** return doc.find('p').map("wat")***REMOVED***,
					function() ***REMOVED*** return ws.connect("WSBIN_URL/ws-echo")***REMOVED***,
				];

				testCases.forEach(function(fn, idx) ***REMOVED***
					var hasException;
					try ***REMOVED***
						fn();
						hasException = false;
					***REMOVED*** catch (e) ***REMOVED***
						hasException = true;
					***REMOVED***

					if (hasException === false) ***REMOVED***
						throw new Error("Expected test case #" + idx + " to return an error");
					***REMOVED*** else if (hasException === undefined) ***REMOVED***
						throw new Error("Something strange happened with test case #" + idx);
					***REMOVED***
				***REMOVED***);
			***REMOVED***
		`))
	require.NoError(t, err)

	ch := make(chan stats.SampleContainer, 1000)
	initVU, err := r.NewVU(1, ch)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	errC := make(chan error)
	go func() ***REMOVED*** errC <- vu.RunOnce() ***REMOVED***()

	select ***REMOVED***
	case <-time.After(15 * time.Second):
		cancel()
		t.Fatal("Test timed out")
	case err := <-errC:
		cancel()
		require.NoError(t, err)
	***REMOVED***
***REMOVED***

func TestSystemTags(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	// Handle paths with custom logic
	tb.Mux.HandleFunc("/wrong-redirect", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Add("Location", "%")
		w.WriteHeader(http.StatusTemporaryRedirect)
	***REMOVED***)

	r, err := getSimpleRunner("/script.js", tb.Replacer.Replace(`
		var http = require("k6/http");

		exports.http_get = function() ***REMOVED***
			http.get("HTTPBIN_IP_URL");
		***REMOVED***;
		exports.https_get = function() ***REMOVED***
			http.get("HTTPSBIN_IP_URL");
		***REMOVED***;
		exports.bad_url_get = function() ***REMOVED***
			http.get("http://127.0.0.1:1");
		***REMOVED***;
		exports.noop = function() ***REMOVED******REMOVED***;
	`), lib.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom("base")***REMOVED***)
	require.NoError(t, err)

	httpURL, err := url.Parse(tb.ServerHTTP.URL)
	require.NoError(t, err)

	testedSystemTags := []struct***REMOVED*** tag, exec, expVal string ***REMOVED******REMOVED***
		***REMOVED***"proto", "http_get", "HTTP/1.1"***REMOVED***,
		***REMOVED***"status", "http_get", "200"***REMOVED***,
		***REMOVED***"method", "http_get", "GET"***REMOVED***,
		***REMOVED***"url", "http_get", tb.ServerHTTP.URL***REMOVED***,
		***REMOVED***"url", "https_get", tb.ServerHTTPS.URL***REMOVED***,
		***REMOVED***"ip", "http_get", httpURL.Hostname()***REMOVED***,
		***REMOVED***"name", "http_get", tb.ServerHTTP.URL***REMOVED***,
		***REMOVED***"group", "http_get", ""***REMOVED***,
		***REMOVED***"vu", "http_get", "8"***REMOVED***,
		***REMOVED***"vu", "noop", "9"***REMOVED***,
		***REMOVED***"iter", "http_get", "0"***REMOVED***,
		***REMOVED***"iter", "noop", "0"***REMOVED***,
		***REMOVED***"tls_version", "https_get", "tls1.3"***REMOVED***,
		***REMOVED***"ocsp_status", "https_get", "unknown"***REMOVED***,
		***REMOVED***"error", "bad_url_get", `dial: connection refused`***REMOVED***,
		***REMOVED***"error_code", "bad_url_get", "1212"***REMOVED***,
		//TODO: add more tests
	***REMOVED***

	samples := make(chan stats.SampleContainer, 100)
	for num, tc := range testedSystemTags ***REMOVED***
		num, tc := num, tc
		t.Run(fmt.Sprintf("TC %d with only %s", num, tc.tag), func(t *testing.T) ***REMOVED***
			require.NoError(t, r.SetOptions(r.GetOptions().Apply(lib.Options***REMOVED***
				Throw:                 null.BoolFrom(false),
				TLSVersion:            &lib.TLSVersions***REMOVED***Max: lib.TLSVersion13***REMOVED***,
				SystemTags:            stats.ToSystemTagSet([]string***REMOVED***tc.tag***REMOVED***),
				InsecureSkipTLSVerify: null.BoolFrom(true),
			***REMOVED***)))

			vu, err := r.NewVU(int64(num), samples)
			require.NoError(t, err)
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***
				RunContext: context.Background(),
				Exec:       tc.exec,
			***REMOVED***)
			require.NoError(t, activeVU.RunOnce())

			bufSamples := stats.GetBufferedSamples(samples)
			assert.NotEmpty(t, bufSamples)
			for _, sample := range bufSamples[0].GetSamples() ***REMOVED***
				assert.NotEmpty(t, sample.Tags)
				for emittedTag, emittedVal := range sample.Tags.CloneTags() ***REMOVED***
					assert.Equal(t, tc.tag, emittedTag)
					assert.Equal(t, tc.expVal, emittedVal)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
