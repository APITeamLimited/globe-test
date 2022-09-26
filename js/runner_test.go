package js

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"go/build"
	"io/ioutil"
	stdlog "log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/test/grpc_testing"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/k6-worker/core"
	"github.com/APITeamLimited/k6-worker/core/local"
	"github.com/APITeamLimited/k6-worker/errext"
	"github.com/APITeamLimited/k6-worker/js/modules/k6"
	k6http "github.com/APITeamLimited/k6-worker/js/modules/k6/http"
	k6metrics "github.com/APITeamLimited/k6-worker/js/modules/k6/metrics"
	"github.com/APITeamLimited/k6-worker/js/modules/k6/ws"
	"github.com/APITeamLimited/k6-worker/lib"
	_ "github.com/APITeamLimited/k6-worker/lib/executor" // TODO: figure out something better
	"github.com/APITeamLimited/k6-worker/lib/fsext"
	"github.com/APITeamLimited/k6-worker/lib/testutils"
	"github.com/APITeamLimited/k6-worker/lib/testutils/httpmultibin"
	"github.com/APITeamLimited/k6-worker/lib/testutils/mockoutput"
	"github.com/APITeamLimited/k6-worker/lib/types"
	"github.com/APITeamLimited/k6-worker/loader"
	"github.com/APITeamLimited/k6-worker/metrics"
	"github.com/APITeamLimited/k6-worker/output"
)

func TestRunnerNew(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Valid", func(t *testing.T) ***REMOVED***
		t.Parallel()
		r, err := getSimpleRunner(t, "/script.js", `
			exports.counter = 0;
			exports.default = function() ***REMOVED*** exports.counter++; ***REMOVED***
		`)
		require.NoError(t, err)

		t.Run("NewVU", func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)
			vuc, ok := initVU.(*VU)
			require.True(t, ok)
			assert.Equal(t, int64(0), vuc.pgm.exports.Get("counter").Export())

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			t.Run("RunOnce", func(t *testing.T) ***REMOVED***
				err = vu.RunOnce()
				require.NoError(t, err)
				assert.Equal(t, int64(1), vuc.pgm.exports.Get("counter").Export())
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)

	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		t.Parallel()
		_, err := getSimpleRunner(t, "/script.js", `blarg`)
		assert.EqualError(t, err, "ReferenceError: blarg is not defined\n\tat file:///script.js:2:1(1)\n\tat native\n")
	***REMOVED***)
***REMOVED***

func TestRunnerGetDefaultGroup(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `exports.default = function() ***REMOVED******REMOVED***;`)
	require.NoError(t, err)
	assert.NotNil(t, r1.GetDefaultGroup())

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)
	assert.NotNil(t, r2.GetDefaultGroup())
***REMOVED***

func TestRunnerOptions(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `exports.default = function() ***REMOVED******REMOVED***;`)
	require.NoError(t, err)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		name, r := name, r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
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

func TestRunnerRPSLimit(t *testing.T) ***REMOVED***
	t.Parallel()

	var nilLimiter *rate.Limiter

	variants := []struct ***REMOVED***
		name    string
		options lib.Options
		limiter *rate.Limiter
	***REMOVED******REMOVED***
		***REMOVED***
			name:    "RPS not defined",
			options: lib.Options***REMOVED******REMOVED***,
			limiter: nilLimiter,
		***REMOVED***,
		***REMOVED***
			name:    "RPS set to non-zero int",
			options: lib.Options***REMOVED***RPS: null.IntFrom(9)***REMOVED***,
			limiter: rate.NewLimiter(rate.Limit(9), 1),
		***REMOVED***,
		***REMOVED***
			name:    "RPS set to zero",
			options: lib.Options***REMOVED***RPS: null.IntFrom(0)***REMOVED***,
			limiter: nilLimiter,
		***REMOVED***,
		***REMOVED***
			name:    "RPS set to below zero value",
			options: lib.Options***REMOVED***RPS: null.IntFrom(-1)***REMOVED***,
			limiter: nilLimiter,
		***REMOVED***,
	***REMOVED***

	for _, variant := range variants ***REMOVED***
		variant := variant

		t.Run(variant.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			r, err := getSimpleRunner(t, "/script.js", `exports.default = function() ***REMOVED******REMOVED***;`)
			require.NoError(t, err)
			err = r.SetOptions(variant.options)
			require.NoError(t, err)
			assert.Equal(t, variant.limiter, r.RPSLimit)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestOptionsSettingToScript(t *testing.T) ***REMOVED***
	t.Parallel()

	optionVariants := []string***REMOVED***
		"export var options = ***REMOVED******REMOVED***;",
		"export var options = ***REMOVED***teardownTimeout: '1s'***REMOVED***;",
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
			r, err := getSimpleRunner(t, "/script.js", data,
				lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"expectedTeardownTimeout": "4s"***REMOVED******REMOVED***)
			require.NoError(t, err)

			newOptions := lib.Options***REMOVED***TeardownTimeout: types.NullDurationFrom(4 * time.Second)***REMOVED***
			r.SetOptions(newOptions)
			require.Equal(t, newOptions, r.GetOptions())

			samples := make(chan metrics.SampleContainer, 100)
			initVU, err := r.NewVU(1, 1, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			require.NoError(t, vu.RunOnce())
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
	r1, err := getSimpleRunner(t, "/script.js", data,
		lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"expectedSetupTimeout": "1s"***REMOVED******REMOVED***)
	require.NoError(t, err)
	require.Equal(t, expScriptOptions, r1.GetOptions())

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
			RuntimeOptions: lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***"expectedSetupTimeout": "3s"***REMOVED******REMOVED***,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())

	require.NoError(t, err)
	require.Equal(t, expScriptOptions, r2.GetOptions())

	newOptions := lib.Options***REMOVED***SetupTimeout: types.NullDurationFrom(3 * time.Second)***REMOVED***
	require.NoError(t, r2.SetOptions(newOptions))
	require.Equal(t, newOptions, r2.GetOptions())

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			samples := make(chan metrics.SampleContainer, 100)

			initVU, err := r.NewVU(1, 1, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			require.NoError(t, vu.RunOnce())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMetricName(t *testing.T) ***REMOVED***
	t.Parallel()

	script := `
		var Counter = require("k6/metrics").Counter;

		var myCounter = new Counter("not ok name @");

		exports.default = function(data) ***REMOVED***
			myCounter.add(1);
		***REMOVED***
	`

	_, err := getSimpleRunner(t, "/script.js", script)
	require.Error(t, err)
***REMOVED***

func TestSetupDataIsolation(t *testing.T) ***REMOVED***
	t.Parallel()

	script := `
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
	`

	runner, err := getSimpleRunner(t, "/script.js", script)
	require.NoError(t, err)

	options := runner.GetOptions()
	require.Empty(t, options.Validate())

	testRunState := &lib.TestRunState***REMOVED***
		TestPreInitState: runner.preInitState,
		Options:          options,
		Runner:           runner,
	***REMOVED***

	execScheduler, err := local.NewExecutionScheduler(testRunState)
	require.NoError(t, err)

	mockOutput := mockoutput.New()
	engine, err := core.NewEngine(testRunState, execScheduler, []output.Output***REMOVED***mockOutput***REMOVED***)
	require.NoError(t, err)
	require.NoError(t, engine.OutputManager.StartOutputs())
	defer engine.OutputManager.StopOutputs()

	ctx, cancel := context.WithCancel(context.Background())
	run, wait, err := engine.Init(ctx, ctx, lib.GetTestWorkerInfo())
	require.NoError(t, err)

	require.Empty(t, runner.defaultGroup.Groups)

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
	require.Contains(t, runner.defaultGroup.Groups, "setup")
	require.Contains(t, runner.defaultGroup.Groups, "teardown")
	var count int
	for _, s := range mockOutput.Samples ***REMOVED***
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
	r1, err := getSimpleRunner(t, "/script.js", data) // TODO fix this
	require.NoError(t, err)
	require.Equal(t, expScriptOptions, r1.GetOptions())

	testdata := map[string]*Runner***REMOVED***"Source": r1***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			samples := make(chan metrics.SampleContainer, 100)

			require.NoError(t, r.Setup(ctx, samples))
			initVU, err := r.NewVU(1, 1, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			require.NoError(t, vu.RunOnce())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSetupDataReturnValue(t *testing.T) ***REMOVED***
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
			console.log("1");
			exports.default = function(data) ***REMOVED***
			***REMOVED***;
		`)
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			samples := make(chan metrics.SampleContainer, 100)
			initVU, err := r.NewVU(1, 1, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			require.NoError(t, vu.RunOnce())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSetupDataNoReturn(t *testing.T) ***REMOVED***
	t.Parallel()
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
	t.Parallel()
	t.Run("Modules", func(t *testing.T) ***REMOVED***
		t.Parallel()
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
					_, err := getSimpleRunner(t, "/script.js", fmt.Sprintf(`import "%s"; exports.default = function() ***REMOVED******REMOVED***`, mod), rtOpts)
					require.NoError(t, err)
				***REMOVED***)
			***REMOVED***)
		***REMOVED***
	***REMOVED***)

	t.Run("Files", func(t *testing.T) ***REMOVED***
		t.Parallel()

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
				t.Parallel()
				fs := afero.NewMemMapFs()
				require.NoError(t, fs.MkdirAll("/path/to", 0o755))
				require.NoError(t, afero.WriteFile(fs, "/path/to/lib.js", []byte(`exports.default = "hi!";`), 0o644))
				r1, err := getSimpleRunner(t, data.filename, fmt.Sprintf(`
					var hi = require("%s").default;
					exports.default = function() ***REMOVED***
						if (hi != "hi!") ***REMOVED*** throw new Error("incorrect value"); ***REMOVED***
					***REMOVED***`, data.path), fs)
				require.NoError(t, err)

				registry := metrics.NewRegistry()
				builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
				r2, err := NewFromArchive(
					&lib.TestPreInitState***REMOVED***
						Logger:         testutils.NewLogger(t),
						BuiltinMetrics: builtinMetrics,
						Registry:       registry,
					***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
				require.NoError(t, err)

				testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
				for name, r := range testdata ***REMOVED***
					r := r
					t.Run(name, func(t *testing.T) ***REMOVED***
						initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
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
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
		exports.options = ***REMOVED*** vus: 10 ***REMOVED***;
		exports.default = function() ***REMOVED*** fn(); ***REMOVED***
	`)
	require.NoError(t, err)
	r1.SetOptions(r1.GetOptions().Apply(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***))

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			vu, err := r.newVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			fnCalled := false
			vu.Runtime.Set("fn", func() ***REMOVED***
				fnCalled = true

				require.NotNil(t, vu.moduleVUImpl.Runtime())
				require.Nil(t, vu.moduleVUImpl.InitEnv())

				state := vu.moduleVUImpl.State()
				require.NotNil(t, state)
				assert.Equal(t, null.IntFrom(10), state.Options.VUs)
				assert.Equal(t, null.BoolFrom(true), state.Options.Throw)
				assert.NotNil(t, state.Logger)
				assert.Equal(t, r.GetDefaultGroup(), state.Group)
				assert.Equal(t, vu.Transport, state.Transport)
			***REMOVED***)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			require.NoError(t, err)
			assert.True(t, fnCalled, "fn() not called")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVURunInterrupt(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
		exports.default = function() ***REMOVED*** while(true) ***REMOVED******REMOVED*** ***REMOVED***
		`)
	require.NoError(t, err)
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***))

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)
	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		name, r := name, r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			samples := make(chan metrics.SampleContainer, 100)
			defer close(samples)
			go func() ***REMOVED***
				for range samples ***REMOVED***
				***REMOVED***
			***REMOVED***()

			vu, err := r.newVU(1, 1, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "context canceled")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVURunInterruptDoesntPanic(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
		exports.default = function() ***REMOVED*** while(true) ***REMOVED******REMOVED*** ***REMOVED***
		`)
	require.NoError(t, err)
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***))

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)
	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			samples := make(chan metrics.SampleContainer, 100)
			defer close(samples)
			go func() ***REMOVED***
				for range samples ***REMOVED***
				***REMOVED***
			***REMOVED***()
			var wg sync.WaitGroup

			initVU, err := r.newVU(1, 1, samples, lib.GetTestWorkerInfo())
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
					require.Error(t, vuErr)
					assert.Contains(t, vuErr.Error(), "context canceled")
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
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
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

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			vu, err := r.newVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			fnOuterCalled := false
			fnInnerCalled := false
			fnNestedCalled := false
			vu.Runtime.Set("fnOuter", func() ***REMOVED***
				fnOuterCalled = true
				assert.Equal(t, r.GetDefaultGroup(), vu.state.Group)
			***REMOVED***)
			vu.Runtime.Set("fnInner", func() ***REMOVED***
				fnInnerCalled = true
				g := vu.state.Group
				assert.Equal(t, "my group", g.Name)
				assert.Equal(t, r.GetDefaultGroup(), g.Parent)
			***REMOVED***)
			vu.Runtime.Set("fnNested", func() ***REMOVED***
				fnNestedCalled = true
				g := vu.state.Group
				assert.Equal(t, "nested group", g.Name)
				assert.Equal(t, "my group", g.Parent.Name)
				assert.Equal(t, r.GetDefaultGroup(), g.Parent.Parent)
			***REMOVED***)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			require.NoError(t, err)
			assert.True(t, fnOuterCalled, "fnOuter() not called")
			assert.True(t, fnInnerCalled, "fnInner() not called")
			assert.True(t, fnNestedCalled, "fnNested() not called")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationMetrics(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
		var group = require("k6").group;
		var Trend = require("k6/metrics").Trend;
		var myMetric = new Trend("my_metric");
		exports.default = function() ***REMOVED*** myMetric.add(5); ***REMOVED***
		`)
	require.NoError(t, err)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	testdata := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range testdata ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			samples := make(chan metrics.SampleContainer, 100)
			defer close(samples)
			vu, err := r.newVU(1, 1, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = activeVU.RunOnce()
			require.NoError(t, err)
			sampleCount := 0
			for i, sampleC := range metrics.GetBufferedSamples(samples) ***REMOVED***
				for j, s := range sampleC.GetSamples() ***REMOVED***
					sampleCount++
					switch i + j ***REMOVED***
					case 0:
						assert.Equal(t, 5.0, s.Value)
						assert.Equal(t, "my_metric", s.Metric.Name)
						assert.Equal(t, metrics.Trend, s.Metric.Type)
					case 1:
						assert.Equal(t, 0.0, s.Value)
						assert.Equal(t, builtinMetrics.DataSent, s.Metric, "`data_sent` sample is before `data_received` and `iteration_duration`")
					case 2:
						assert.Equal(t, 0.0, s.Value)
						assert.Equal(t, builtinMetrics.DataReceived, s.Metric, "`data_received` sample is after `data_received`")
					case 3:
						assert.Equal(t, builtinMetrics.IterationDuration, s.Metric, "`iteration-duration` sample is after `data_received`")
					case 4:
						assert.Equal(t, builtinMetrics.Iterations, s.Metric, "`iterations` sample is after `iteration_duration`")
						assert.Equal(t, float64(1), s.Value)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			assert.Equal(t, sampleCount, 5)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func GenerateTLSCertificate(t *testing.T, host string, notBefore time.Time, validFor time.Duration) ([]byte, []byte) ***REMOVED***
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature
	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
	// the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.
	keyUsage |= x509.KeyUsageKeyEncipherment

	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	require.NoError(t, err)

	template := x509.Certificate***REMOVED***
		SerialNumber: serialNumber,
		Subject: pkix.Name***REMOVED***
			Organization: []string***REMOVED***"Acme Co"***REMOVED***,
		***REMOVED***,
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage***REMOVED***x509.ExtKeyUsageServerAuth***REMOVED***,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
	***REMOVED***

	hosts := strings.Split(host, ",")
	for _, h := range hosts ***REMOVED***
		if ip := net.ParseIP(h); ip != nil ***REMOVED***
			template.IPAddresses = append(template.IPAddresses, ip)
		***REMOVED*** else ***REMOVED***
			template.DNSNames = append(template.DNSNames, h)
		***REMOVED***
	***REMOVED***

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)

	certPem := pem.EncodeToMemory(&pem.Block***REMOVED***Type: "CERTIFICATE", Bytes: derBytes***REMOVED***)
	require.NoError(t, err)

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	require.NoError(t, err)
	keyPem := pem.EncodeToMemory(&pem.Block***REMOVED***Type: "PRIVATE KEY", Bytes: privBytes***REMOVED***)
	require.NoError(t, err)
	return certPem, keyPem
***REMOVED***

func GetTestServerWithCertificate(t *testing.T, certPem, key []byte) *httptest.Server ***REMOVED***
	server := &http.Server***REMOVED***
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
			w.WriteHeader(200)
		***REMOVED***),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Second,
	***REMOVED***
	s := &httptest.Server***REMOVED******REMOVED***
	s.Config = server

	s.TLS = new(tls.Config)
	if s.TLS.NextProtos == nil ***REMOVED***
		nextProtos := []string***REMOVED***"http/1.1"***REMOVED***
		if s.EnableHTTP2 ***REMOVED***
			nextProtos = []string***REMOVED***"h2"***REMOVED***
		***REMOVED***
		s.TLS.NextProtos = nextProtos
	***REMOVED***
	cert, err := tls.X509KeyPair(certPem, key)
	require.NoError(t, err)
	s.TLS.Certificates = append(s.TLS.Certificates, cert)
	for _, suite := range tls.CipherSuites() ***REMOVED***
		if !strings.Contains(suite.Name, "256") ***REMOVED***
			continue
		***REMOVED***
		s.TLS.CipherSuites = append(s.TLS.CipherSuites, suite.ID)
	***REMOVED***
	certpool := x509.NewCertPool()
	certificate, err := x509.ParseCertificate(cert.Certificate[0])
	require.NoError(t, err)
	certpool.AddCert(certificate)
	client := &http.Client***REMOVED***Transport: &http.Transport***REMOVED******REMOVED******REMOVED***
	client.Transport = &http.Transport***REMOVED***
		TLSClientConfig: &tls.Config***REMOVED*** //nolint:gosec
			RootCAs: certpool,
		***REMOVED***,
		ForceAttemptHTTP2: s.EnableHTTP2,
	***REMOVED***
	s.Listener, err = net.Listen("tcp", "")
	require.NoError(t, err)
	s.Listener = tls.NewListener(s.Listener, s.TLS)
	s.URL = "https://" + s.Listener.Addr().String()
	return s
***REMOVED***

func TestVUIntegrationInsecureRequests(t *testing.T) ***REMOVED***
	t.Parallel()
	certPem, keyPem := GenerateTLSCertificate(t, "mybadssl.localhost", time.Now(), 0)
	s := GetTestServerWithCertificate(t, certPem, keyPem)
	go func() ***REMOVED***
		_ = s.Config.Serve(s.Listener)
	***REMOVED***()
	t.Cleanup(func() ***REMOVED***
		require.NoError(t, s.Config.Close())
	***REMOVED***)
	host, port, err := net.SplitHostPort(s.Listener.Addr().String())
	require.NoError(t, err)
	ip := net.ParseIP(host)
	mybadsslHostname, err := lib.NewHostAddress(ip, port)
	require.NoError(t, err)
	cert, err := x509.ParseCertificate(s.TLS.Certificates[0].Certificate[0])
	require.NoError(t, err)

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
			t.Parallel()
			r1, err := getSimpleRunner(t, "/script.js", `
			  var http = require("k6/http");;
        exports.default = function() ***REMOVED*** http.get("https://mybadssl.localhost/"); ***REMOVED***
				`)
			require.NoError(t, err)
			require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***.Apply(data.opts)))

			r1.Bundle.Options.Hosts = map[string]*lib.HostAddress***REMOVED***
				"mybadssl.localhost": mybadsslHostname,
			***REMOVED***
			registry := metrics.NewRegistry()
			builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
			r2, err := NewFromArchive(
				&lib.TestPreInitState***REMOVED***
					Logger:         testutils.NewLogger(t),
					BuiltinMetrics: builtinMetrics,
					Registry:       registry,
				***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
			require.NoError(t, err)
			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				r := r
				t.Run(name, func(t *testing.T) ***REMOVED***
					t.Parallel()
					r.preInitState.Logger, _ = logtest.NewNullLogger()

					initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
					require.NoError(t, err)
					initVU.(*VU).TLSConfig.RootCAs = x509.NewCertPool() //nolint:forcetypeassert
					initVU.(*VU).TLSConfig.RootCAs.AddCert(cert)        //nolint:forcetypeassert

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					err = vu.RunOnce()
					if data.errMsg != "" ***REMOVED***
						require.Error(t, err)
						assert.Contains(t, err.Error(), data.errMsg)
					***REMOVED*** else ***REMOVED***
						require.NoError(t, err)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationBlacklistOption(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
					var http = require("k6/http");;
					exports.default = function() ***REMOVED*** http.get("http://10.1.2.3/"); ***REMOVED***
				`)
	require.NoError(t, err)

	cidr, err := lib.ParseCIDR("10.0.0.0/8")

	require.NoError(t, err)
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***
		Throw:        null.BoolFrom(true),
		BlacklistIPs: []*lib.IPNet***REMOVED***cidr***REMOVED***,
	***REMOVED***))

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)
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
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
					var http = require("k6/http");;

					exports.options = ***REMOVED***
						throw: true,
						blacklistIPs: ["10.0.0.0/8"],
					***REMOVED***;

					exports.default = function() ***REMOVED*** http.get("http://10.1.2.3/"); ***REMOVED***
				`)
	require.NoError(t, err)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***

	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "IP (10.1.2.3) is in a blacklisted range (10.0.0.0/8)")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationBlockHostnamesOption(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
					var http = require("k6/http");
					exports.default = function() ***REMOVED*** http.get("https://k6.io/"); ***REMOVED***
				`)
	require.NoError(t, err)

	hostnames, err := types.NewNullHostnameTrie([]string***REMOVED***"*.io"***REMOVED***)
	require.NoError(t, err)
	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***
		Throw:            null.BoolFrom(true),
		BlockedHostnames: hostnames,
	***REMOVED***))

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***

	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVu, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "hostname (k6.io) is in a blocked pattern (*.io)")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationBlockHostnamesScript(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
					var http = require("k6/http");

					exports.options = ***REMOVED***
						throw: true,
						blockHostnames: ["*.io"],
					***REMOVED***;

					exports.default = function() ***REMOVED*** http.get("https://k6.io/"); ***REMOVED***
				`)
	require.NoError(t, err)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***

	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVu, err := r.NewVU(0, 0, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "hostname (k6.io) is in a blocked pattern (*.io)")
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationHosts(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	r1, err := getSimpleRunner(t, "/script.js",
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
	require.NoError(t, err)

	r1.SetOptions(lib.Options***REMOVED***
		Throw: null.BoolFrom(true),
		Hosts: map[string]*lib.HostAddress***REMOVED***
			"test.loadimpact.com": ***REMOVED***IP: net.ParseIP("127.0.0.1")***REMOVED***,
		***REMOVED***,
	***REMOVED***)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationTLSConfig(t *testing.T) ***REMOVED***
	t.Parallel()
	certPem, keyPem := GenerateTLSCertificate(t, "sha256-badssl.localhost", time.Now(), time.Hour)
	s := GetTestServerWithCertificate(t, certPem, keyPem)
	go func() ***REMOVED***
		_ = s.Config.Serve(s.Listener)
	***REMOVED***()
	t.Cleanup(func() ***REMOVED***
		require.NoError(t, s.Config.Close())
	***REMOVED***)
	host, port, err := net.SplitHostPort(s.Listener.Addr().String())
	require.NoError(t, err)
	ip := net.ParseIP(host)
	mybadsslHostname, err := lib.NewHostAddress(ip, port)
	require.NoError(t, err)
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
			lib.Options***REMOVED***
				TLSCipherSuites: &lib.TLSCipherSuites***REMOVED***tls.TLS_RSA_WITH_RC4_128_SHA***REMOVED***,
				TLSVersion:      &lib.TLSVersions***REMOVED***Max: tls.VersionTLS12***REMOVED***,
			***REMOVED***,
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
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	cert, err := x509.ParseCertificate(s.TLS.Certificates[0].Certificate[0])
	require.NoError(t, err)
	for name, data := range testdata ***REMOVED***
		data := data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			r1, err := getSimpleRunner(t, "/script.js", `
					var http = require("k6/http");;
					exports.default = function() ***REMOVED*** http.get("https://sha256-badssl.localhost/"); ***REMOVED***
				`)
			require.NoError(t, err)
			require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***.Apply(data.opts)))

			r1.Bundle.Options.Hosts = map[string]*lib.HostAddress***REMOVED***
				"sha256-badssl.localhost": mybadsslHostname,
			***REMOVED***
			r2, err := NewFromArchive(
				&lib.TestPreInitState***REMOVED***
					Logger:         testutils.NewLogger(t),
					BuiltinMetrics: builtinMetrics,
					Registry:       registry,
				***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				r := r
				t.Run(name, func(t *testing.T) ***REMOVED***
					t.Parallel()
					r.preInitState.Logger, _ = logtest.NewNullLogger()

					initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
					require.NoError(t, err)
					initVU.(*VU).TLSConfig.RootCAs = x509.NewCertPool() //nolint:forcetypeassert
					initVU.(*VU).TLSConfig.RootCAs.AddCert(cert)        //nolint:forcetypeassert
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					err = vu.RunOnce()
					if data.errMsg != "" ***REMOVED***
						require.Error(t, err, "for message %q", data.errMsg)
						assert.Contains(t, err.Error(), data.errMsg)
					***REMOVED*** else ***REMOVED***
						require.NoError(t, err)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationOpenFunctionError(t *testing.T) ***REMOVED***
	t.Parallel()
	r, err := getSimpleRunner(t, "/script.js", `
			exports.default = function() ***REMOVED*** open("/tmp/foo") ***REMOVED***
		`)
	require.NoError(t, err)

	initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	err = vu.RunOnce()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only available in the init stage")
***REMOVED***

func TestVUIntegrationOpenFunctionErrorWhenSneaky(t *testing.T) ***REMOVED***
	t.Parallel()
	r, err := getSimpleRunner(t, "/script.js", `
			var sneaky = open;
			exports.default = function() ***REMOVED*** sneaky("/tmp/foo") ***REMOVED***
		`)
	require.NoError(t, err)

	initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	err = vu.RunOnce()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only available in the init stage")
***REMOVED***

func TestVUDoesOpenUnderV0Condition(t *testing.T) ***REMOVED***
	t.Parallel()

	baseFS := afero.NewMemMapFs()
	data := `
			if (__VU == 0) ***REMOVED***
				let data = open("/home/somebody/test.json");
			***REMOVED***
			exports.default = function() ***REMOVED***
				console.log("hey")
			***REMOVED***
		`
	require.NoError(t, afero.WriteFile(baseFS, "/home/somebody/test.json", []byte(`42`), os.ModePerm))
	require.NoError(t, afero.WriteFile(baseFS, "/script.js", []byte(data), os.ModePerm))

	fs := fsext.NewCacheOnReadFs(baseFS, afero.NewMemMapFs(), 0)

	r, err := getSimpleRunner(t, "/script.js", data, fs)
	require.NoError(t, err)

	_, err = r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
	require.NoError(t, err)
***REMOVED***

func TestVUDoesNotOpenUnderConditions(t *testing.T) ***REMOVED***
	t.Parallel()

	baseFS := afero.NewMemMapFs()
	data := `
			if (__VU > 0) ***REMOVED***
				let data = open("/home/somebody/test.json");
			***REMOVED***
			exports.default = function() ***REMOVED***
				console.log("hey")
			***REMOVED***
		`
	require.NoError(t, afero.WriteFile(baseFS, "/home/somebody/test.json", []byte(`42`), os.ModePerm))
	require.NoError(t, afero.WriteFile(baseFS, "/script.js", []byte(data), os.ModePerm))

	fs := fsext.NewCacheOnReadFs(baseFS, afero.NewMemMapFs(), 0)

	r, err := getSimpleRunner(t, "/script.js", data, fs)
	require.NoError(t, err)

	_, err = r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "open() can't be used with files that weren't previously opened during initialization (__VU==0)")
***REMOVED***

func TestVUDoesNonExistingPathnUnderConditions(t *testing.T) ***REMOVED***
	t.Parallel()

	baseFS := afero.NewMemMapFs()
	data := `
			if (__VU == 1) ***REMOVED***
				let data = open("/home/nobody");
			***REMOVED***
			exports.default = function() ***REMOVED***
				console.log("hey")
			***REMOVED***
		`
	require.NoError(t, afero.WriteFile(baseFS, "/script.js", []byte(data), os.ModePerm))

	fs := fsext.NewCacheOnReadFs(baseFS, afero.NewMemMapFs(), 0)

	r, err := getSimpleRunner(t, "/script.js", data, fs)
	require.NoError(t, err)

	_, err = r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "open() can't be used with files that weren't previously opened during initialization (__VU==0)")
***REMOVED***

func TestVUIntegrationCookiesReset(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	r1, err := getSimpleRunner(t, "/script.js", tb.Replacer.Replace(`
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
	require.NoError(t, err)
	r1.SetOptions(lib.Options***REMOVED***
		Throw:        null.BoolFrom(true),
		MaxRedirects: null.IntFrom(10),
		Hosts:        tb.Dialer.Hosts,
	***REMOVED***)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			for i := 0; i < 2; i++ ***REMOVED***
				require.NoError(t, vu.RunOnce())
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationCookiesNoReset(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	r1, err := getSimpleRunner(t, "/script.js", tb.Replacer.Replace(`
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
	require.NoError(t, err)
	r1.SetOptions(lib.Options***REMOVED***
		Throw:          null.BoolFrom(true),
		MaxRedirects:   null.IntFrom(10),
		Hosts:          tb.Dialer.Hosts,
		NoCookiesReset: null.BoolFrom(true),
	***REMOVED***)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)

			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestVUIntegrationVUID(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
			exports.default = function() ***REMOVED***
				if (__VU != 1234) ***REMOVED*** throw new Error("wrong __VU: " + __VU); ***REMOVED***
			***REMOVED***`,
	)
	require.NoError(t, err)
	r1.SetOptions(lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1234, 1234, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

/*
CA key:
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIDEm8bxihqYfAsWP39o5DpkAksPBw+3rlDHNX+d69oYGoAoGCCqGSM49
AwEHoUQDQgAEeeuCFQsdraFJr8JaKbAKfjYpZ2U+p3r/OzcmAsjFO8EckmV9uFZs
Gq3JurKi9Z3dDKQcwinHQ1malicbwWhamQ==
-----END EC PRIVATE KEY-----
*/
func TestVUIntegrationClientCerts(t *testing.T) ***REMOVED***
	t.Parallel()
	clientCAPool := x509.NewCertPool()
	assert.True(t, clientCAPool.AppendCertsFromPEM(
		[]byte("-----BEGIN CERTIFICATE-----\n"+
			"MIIBWzCCAQGgAwIBAgIJAIQMBgLi+DV6MAoGCCqGSM49BAMCMBAxDjAMBgNVBAMM\n"+
			"BU15IENBMCAXDTIyMDEyMTEyMjkzNloYDzMwMjEwNTI0MTIyOTM2WjAQMQ4wDAYD\n"+
			"VQQDDAVNeSBDQTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHnrghULHa2hSa/C\n"+
			"WimwCn42KWdlPqd6/zs3JgLIxTvBHJJlfbhWbBqtybqyovWd3QykHMIpx0NZmpYn\n"+
			"G8FoWpmjQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1Ud\n"+
			"DgQWBBSkukBA8lgFvvBJAYKsoSUR+PX71jAKBggqhkjOPQQDAgNIADBFAiEAiFF7\n"+
			"Y54CMNRSBSVMgd4mQgrzJInRH88KpLsQ7VeOAaQCIEa0vaLln9zxIDZQKocml4Db\n"+
			"AEJr8tDzMKIds6sRTBT4\n"+
			"-----END CERTIFICATE-----"),
	))
	serverCert, err := tls.X509KeyPair(
		[]byte("-----BEGIN CERTIFICATE-----\n"+
			"MIIBcTCCARigAwIBAgIJAIP0njRt16gbMAoGCCqGSM49BAMCMBAxDjAMBgNVBAMM\n"+
			"BU15IENBMCAXDTIyMDEyMTE1MTA0OVoYDzMwMjEwNTI0MTUxMDQ5WjAZMRcwFQYD\n"+
			"VQQDDA4xMjcuMC4wLjE6Njk2OTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABH8Y\n"+
			"exy5LI9r+RNwVpf/5ZX86EigMYHp9YOyiUMmfUfvDig+BGhlwjm7Lh2941Gz4amO\n"+
			"lpN2YAkcd0wnNLHkVOmjUDBOMA4GA1UdDwEB/wQEAwIBBjAMBgNVHRMBAf8EAjAA\n"+
			"MB0GA1UdDgQWBBQ9cIYUwwzfzBXPyRGB5tNpAgHWujAPBgNVHREECDAGhwR/AAAB\n"+
			"MAoGCCqGSM49BAMCA0cAMEQCIDjRZlg+jKgI9K99HOM2wS9+URr6R1/FYLZYBtMc\n"+
			"pq3hAiB9NQxNqV459fgN0BpbiLrEvJjquRFoUr9BWsG+hHrHtQ==\n"+
			"-----END CERTIFICATE-----\n"+
			"-----BEGIN CERTIFICATE-----\n"+
			"MIIBWzCCAQGgAwIBAgIJAIQMBgLi+DV6MAoGCCqGSM49BAMCMBAxDjAMBgNVBAMM\n"+
			"BU15IENBMCAXDTIyMDEyMTEyMjkzNloYDzMwMjEwNTI0MTIyOTM2WjAQMQ4wDAYD\n"+
			"VQQDDAVNeSBDQTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHnrghULHa2hSa/C\n"+
			"WimwCn42KWdlPqd6/zs3JgLIxTvBHJJlfbhWbBqtybqyovWd3QykHMIpx0NZmpYn\n"+
			"G8FoWpmjQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1Ud\n"+
			"DgQWBBSkukBA8lgFvvBJAYKsoSUR+PX71jAKBggqhkjOPQQDAgNIADBFAiEAiFF7\n"+
			"Y54CMNRSBSVMgd4mQgrzJInRH88KpLsQ7VeOAaQCIEa0vaLln9zxIDZQKocml4Db\n"+
			"AEJr8tDzMKIds6sRTBT4\n"+
			"-----END CERTIFICATE-----"),
		[]byte("-----BEGIN EC PRIVATE KEY-----\n"+
			"MHcCAQEEIHNpjs0P9/ejoUYF5Agzf9clHR4PwBsVfZ+JgslfuBg1oAoGCCqGSM49\n"+
			"AwEHoUQDQgAEfxh7HLksj2v5E3BWl//llfzoSKAxgen1g7KJQyZ9R+8OKD4EaGXC\n"+
			"ObsuHb3jUbPhqY6Wk3ZgCRx3TCc0seRU6Q==\n"+
			"-----END EC PRIVATE KEY-----"),
	)
	require.NoError(t, err)

	testdata := map[string]struct ***REMOVED***
		withClientCert     bool
		withDomains        bool
		insecureSkipVerify bool
		errMsg             string
	***REMOVED******REMOVED***
		"WithoutCert":      ***REMOVED***false, false, true, "remote error: tls: bad certificate"***REMOVED***,
		"WithCert":         ***REMOVED***true, true, true, ""***REMOVED***,
		"VerifyServerCert": ***REMOVED***true, false, false, "certificate signed by unknown authority"***REMOVED***,
		"WithoutDomains":   ***REMOVED***true, false, true, ""***REMOVED***,
	***REMOVED***

	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config***REMOVED***
		Certificates: []tls.Certificate***REMOVED***serverCert***REMOVED***,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAPool,
	***REMOVED***)
	require.NoError(t, err)
	srv := &http.Server***REMOVED***
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) ***REMOVED***
			_, _ = fmt.Fprintf(w, "ok")
		***REMOVED***),
		ErrorLog: stdlog.New(ioutil.Discard, "", 0),
	***REMOVED***
	go func() ***REMOVED*** _ = srv.Serve(listener) ***REMOVED***()
	t.Cleanup(func() ***REMOVED*** _ = listener.Close() ***REMOVED***)
	for name, data := range testdata ***REMOVED***
		data := data

		registry := metrics.NewRegistry()
		builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			r1, err := getSimpleRunner(t, "/script.js", fmt.Sprintf(`
			var http = require("k6/http");
			var k6 = require("k6");
			var check = k6.check;
			exports.default = function() ***REMOVED***
				const res = http.get("https://%s")
				check(res, ***REMOVED***
					'is status 200': (r) => r.status === 200,
					'verify resp': (r) => r.body.includes('ok'),
				***REMOVED***)
			***REMOVED***`, listener.Addr().String()))
			require.NoError(t, err)

			opt := lib.Options***REMOVED***Throw: null.BoolFrom(true)***REMOVED***
			if data.insecureSkipVerify ***REMOVED***
				opt.InsecureSkipTLSVerify = null.BoolFrom(true)
			***REMOVED***
			if data.withClientCert ***REMOVED***
				opt.TLSAuth = []*lib.TLSAuth***REMOVED***
					***REMOVED***
						TLSAuthFields: lib.TLSAuthFields***REMOVED***
							Cert: "-----BEGIN CERTIFICATE-----\n" +
								"MIIBVzCB/6ADAgECAgkAg/SeNG3XqB0wCgYIKoZIzj0EAwIwEDEOMAwGA1UEAwwF\n" +
								"TXkgQ0EwIBcNMjIwMTIxMTUxMjM0WhgPMzAyMTA1MjQxNTEyMzRaMBExDzANBgNV\n" +
								"BAMMBmNsaWVudDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABKM7OJQMYG4KLtDA\n" +
								"gZ8zOg2PimHMmQnjD2HtI4cSwIUJJnvHWLowbFe9fk6XeP9b3dK1ImUI++/EZdVr\n" +
								"ABAcngejPzA9MA4GA1UdDwEB/wQEAwIBBjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQW\n" +
								"BBSttJe1mcPEnBOZ6wvKPG4zL0m1CzAKBggqhkjOPQQDAgNHADBEAiBPSLgKA/r9\n" +
								"u/FW6W+oy6Odm1kdNMGCI472iTn545GwJgIgb3UQPOUTOj0IN4JLJYfmYyXviqsy\n" +
								"zk9eWNHFXDA9U6U=\n" +
								"-----END CERTIFICATE-----",
							Key: "-----BEGIN EC PRIVATE KEY-----\n" +
								"MHcCAQEEINDaMGkOT3thu1A0LfLJr3Jd011/aEG6OArmEQaujwgpoAoGCCqGSM49\n" +
								"AwEHoUQDQgAEozs4lAxgbgou0MCBnzM6DY+KYcyZCeMPYe0jhxLAhQkme8dYujBs\n" +
								"V71+Tpd4/1vd0rUiZQj778Rl1WsAEByeBw==\n" +
								"-----END EC PRIVATE KEY-----",
						***REMOVED***,
					***REMOVED***,
				***REMOVED***
				if data.withDomains ***REMOVED***
					opt.TLSAuth[0].TLSAuthFields.Domains = []string***REMOVED***"127.0.0.1"***REMOVED***
				***REMOVED***
				_, _ = opt.TLSAuth[0].Certificate()
			***REMOVED***
			require.NoError(t, r1.SetOptions(opt))
			r2, err := NewFromArchive(
				&lib.TestPreInitState***REMOVED***
					Logger:         testutils.NewLogger(t),
					BuiltinMetrics: builtinMetrics,
					Registry:       registry,
				***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				r := r
				t.Run(name, func(t *testing.T) ***REMOVED***
					t.Parallel()
					r.preInitState.Logger, _ = logtest.NewNullLogger()
					initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
					require.NoError(t, err)
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					err = vu.RunOnce()
					if len(data.errMsg) > 0 ***REMOVED***
						require.Error(t, err)
						assert.Contains(t, err.Error(), data.errMsg)
					***REMOVED*** else ***REMOVED***
						require.NoError(t, err)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestHTTPRequestInInitContext(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	_, err := getSimpleRunner(t, "/script.js", tb.Replacer.Replace(`
					var k6 = require("k6");
					var check = k6.check;
					var fail = k6.fail;
					var http = require("k6/http");;
					var res = http.get("HTTPBIN_URL/");
					exports.default = function() ***REMOVED***
						console.log(test);
					***REMOVED***
				`))
	require.Error(t, err)
	assert.Contains(
		t,
		err.Error(),
		k6http.ErrHTTPForbiddenInInitContext.Error())
***REMOVED***

func TestInitContextForbidden(t *testing.T) ***REMOVED***
	t.Parallel()
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
			"abortTest",
			`var test = require("k6/execution").test;
			 test.abort();
			 exports.default = function() ***REMOVED*** console.log("p"); ***REMOVED***`,
			errext.AbortTest,
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

	for _, test := range table ***REMOVED***
		test := test
		t.Run(test[0], func(t *testing.T) ***REMOVED***
			t.Parallel()
			_, err := getSimpleRunner(t, "/script.js", tb.Replacer.Replace(test[1]))
			require.Error(t, err)
			assert.Contains(
				t,
				err.Error(),
				test[2])
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestArchiveRunningIntegrity(t *testing.T) ***REMOVED***
	t.Parallel()

	fs := afero.NewMemMapFs()
	data := `
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
		`
	require.NoError(t, afero.WriteFile(fs, "/home/somebody/test.json", []byte(`42`), os.ModePerm))
	require.NoError(t, afero.WriteFile(fs, "/script.js", []byte(data), os.ModePerm))
	r1, err := getSimpleRunner(t, "/script.js", data, fs)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, r1.MakeArchive().Write(buf))

	arc, err := lib.ReadArchive(buf)
	require.NoError(t, err)
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, arc, lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			var err error
			ch := make(chan metrics.SampleContainer, 100)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err = r.Setup(ctx, ch)
			cancel()
			require.NoError(t, err)
			initVU, err := r.NewVU(1, 1, ch, lib.GetTestWorkerInfo())
			require.NoError(t, err)
			ctx, cancel = context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestArchiveNotPanicking(t *testing.T) ***REMOVED***
	t.Parallel()
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/non/existent", []byte(`42`), os.ModePerm))
	r1, err := getSimpleRunner(t, "/script.js", `
			var fput = open("/non/existent");
			exports.default = function(data) ***REMOVED******REMOVED***
		`, fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	arc.Filesystems = map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs()***REMOVED***
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, arc, lib.GetTestWorkerInfo())
	// we do want this to error here as this is where we find out that a given file is not in the
	// archive
	require.Error(t, err)
	require.Nil(t, r2)
***REMOVED***

func TestStuffNotPanicking(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	r, err := getSimpleRunner(t, "/script.js", tb.Replacer.Replace(`
			var http = require("k6/http");
			var ws = require("k6/ws");
			var group = require("k6").group;
			var parseHTML = require("k6/html").parseHTML;

			exports.options = ***REMOVED*** iterations: 1, vus: 1 ***REMOVED***;

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

	ch := make(chan metrics.SampleContainer, 1000)
	initVU, err := r.NewVU(1, 1, ch, lib.GetTestWorkerInfo())
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

func TestPanicOnSimpleHTML(t *testing.T) ***REMOVED***
	t.Parallel()

	r, err := getSimpleRunner(t, "/script.js", `
			var parseHTML = require("k6/html").parseHTML;

			exports.options = ***REMOVED*** iterations: 1, vus: 1 ***REMOVED***;

			exports.default = function() ***REMOVED***
				var doc = parseHTML("<html>");
				var o = doc.find(".something").slice(0, 4).toArray()
			***REMOVED***;
		`)
	require.NoError(t, err)

	ch := make(chan metrics.SampleContainer, 1000)
	initVU, err := r.NewVU(1, 1, ch, lib.GetTestWorkerInfo())
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

	// Handle paths with custom logic
	tb.Mux.HandleFunc("/wrong-redirect", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Add("Location", "%")
		w.WriteHeader(http.StatusTemporaryRedirect)
	***REMOVED***)

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
		***REMOVED***"scenario", "http_get", "default"***REMOVED***,
		// TODO: add more tests
	***REMOVED***

	for num, tc := range testedSystemTags ***REMOVED***
		num, tc := num, tc
		t.Run(fmt.Sprintf("TC %d with only %s", num, tc.tag), func(t *testing.T) ***REMOVED***
			t.Parallel()
			samples := make(chan metrics.SampleContainer, 100)
			r, err := getSimpleRunner(t, "/script.js", tb.Replacer.Replace(`
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
			require.NoError(t, r.SetOptions(r.GetOptions().Apply(lib.Options***REMOVED***
				Throw:                 null.BoolFrom(false),
				TLSVersion:            &lib.TLSVersions***REMOVED***Max: tls.VersionTLS13***REMOVED***,
				SystemTags:            metrics.ToSystemTagSet([]string***REMOVED***tc.tag***REMOVED***),
				InsecureSkipTLSVerify: null.BoolFrom(true),
			***REMOVED***)))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			vu, err := r.NewVU(uint64(num), 0, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)
			activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***
				RunContext: ctx,
				Exec:       tc.exec,
				Scenario:   "default",
			***REMOVED***)
			require.NoError(t, activeVU.RunOnce())

			bufSamples := metrics.GetBufferedSamples(samples)
			require.NotEmpty(t, bufSamples)
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

func TestVUPanic(t *testing.T) ***REMOVED***
	t.Parallel()
	r1, err := getSimpleRunner(t, "/script.js", `
			var group = require("k6").group;
			exports.default = function() ***REMOVED***
				group("panic here", function() ***REMOVED***
					if (__ITER == 0) ***REMOVED***
						panic("here we panic");
					***REMOVED***
					console.log("here we don't");
				***REMOVED***)
			***REMOVED***`,
	)
	require.NoError(t, err)

	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	r2, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         testutils.NewLogger(t),
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
		***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			initVU, err := r.NewVU(1, 1234, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			logger := logrus.New()
			logger.SetLevel(logrus.InfoLevel)
			logger.Out = ioutil.Discard
			hook := testutils.SimpleLogrusHook***REMOVED***
				HookedLevels: []logrus.Level***REMOVED***logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel***REMOVED***,
			***REMOVED***
			logger.AddHook(&hook)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			vu.(*ActiveVU).Runtime.Set("panic", func(str string) ***REMOVED*** panic(str) ***REMOVED***)
			vu.(*ActiveVU).state.Logger = logger

			vu.(*ActiveVU).Console.logger = logger.WithField("source", "console")
			err = vu.RunOnce()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "a panic occurred during JS execution: here we panic")
			entries := hook.Drain()
			require.Len(t, entries, 1)
			assert.Equal(t, logrus.ErrorLevel, entries[0].Level)
			require.True(t, strings.HasPrefix(entries[0].Message, "panic: here we panic"))
			// broken since goja@f3cfc97811c0b4d8337902c3e42fb2371ba1d524 see
			// https://github.com/dop251/goja/issues/179#issuecomment-783572020
			// require.True(t, strings.HasSuffix(entries[0].Message, "Goja stack:\nfile:///script.js:3:4(12)"))

			err = vu.RunOnce()
			require.NoError(t, err)

			entries = hook.Drain()
			require.Len(t, entries, 1)
			assert.Equal(t, logrus.InfoLevel, entries[0].Level)
			require.Contains(t, entries[0].Message, "here we don't")
		***REMOVED***)
	***REMOVED***
***REMOVED***

type multiFileTestCase struct ***REMOVED***
	fses       map[string]afero.Fs
	rtOpts     lib.RuntimeOptions
	cwd        string
	script     string
	expInitErr bool
	expVUErr   bool
	samples    chan metrics.SampleContainer
***REMOVED***

func runMultiFileTestCase(t *testing.T, tc multiFileTestCase, tb *httpmultibin.HTTPMultiBin) ***REMOVED***
	t.Helper()
	logger := testutils.NewLogger(t)
	registry := metrics.NewRegistry()
	builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
	runner, err := New(
		&lib.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
			RuntimeOptions: tc.rtOpts,
		***REMOVED***,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: tc.cwd + "/script.js", Scheme: "file"***REMOVED***,
			Data: []byte(tc.script),
		***REMOVED***,
		tc.fses, lib.GetTestWorkerInfo(),
	)
	if tc.expInitErr ***REMOVED***
		require.Error(t, err)
		return
	***REMOVED***
	require.NoError(t, err)

	options := runner.GetOptions()
	require.Empty(t, options.Validate())

	vu, err := runner.NewVU(1, 1, tc.samples, lib.GetTestWorkerInfo())
	require.NoError(t, err)

	jsVU, ok := vu.(*VU)
	require.True(t, ok)
	jsVU.state.Dialer = tb.Dialer
	jsVU.state.TLSConfig = tb.TLSClientConfig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	activeVU := vu.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)

	err = activeVU.RunOnce()
	if tc.expVUErr ***REMOVED***
		require.Error(t, err)
	***REMOVED*** else ***REMOVED***
		require.NoError(t, err)
	***REMOVED***

	arc := runner.MakeArchive()
	runnerFromArc, err := NewFromArchive(
		&lib.TestPreInitState***REMOVED***
			Logger:         logger,
			BuiltinMetrics: builtinMetrics,
			Registry:       registry,
			RuntimeOptions: tc.rtOpts,
		***REMOVED***, arc, lib.GetTestWorkerInfo())
	require.NoError(t, err)
	vuFromArc, err := runnerFromArc.NewVU(2, 2, tc.samples, lib.GetTestWorkerInfo())
	require.NoError(t, err)
	jsVUFromArc, ok := vuFromArc.(*VU)
	require.True(t, ok)
	jsVUFromArc.state.Dialer = tb.Dialer
	jsVUFromArc.state.TLSConfig = tb.TLSClientConfig
	activeVUFromArc := jsVUFromArc.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	err = activeVUFromArc.RunOnce()
	if tc.expVUErr ***REMOVED***
		require.Error(t, err)
		return
	***REMOVED***
	require.NoError(t, err)
***REMOVED***

func TestComplicatedFileImportsForGRPC(t *testing.T) ***REMOVED***
	t.Parallel()
	tb := httpmultibin.NewHTTPMultiBin(t)

	tb.GRPCStub.UnaryCallFunc = func(ctx context.Context, sreq *grpc_testing.SimpleRequest) (
		*grpc_testing.SimpleResponse, error,
	) ***REMOVED***
		return &grpc_testing.SimpleResponse***REMOVED***
			Username: "foo",
		***REMOVED***, nil
	***REMOVED***

	fs := afero.NewMemMapFs()
	protoFile, err := ioutil.ReadFile("../vendor/google.golang.org/grpc/test/grpc_testing/test.proto")
	require.NoError(t, err)
	require.NoError(t, afero.WriteFile(fs, "/path/to/service.proto", protoFile, 0o644))
	require.NoError(t, afero.WriteFile(fs, "/path/to/same-dir.proto", []byte(
		`syntax = "proto3";package whatever;import "service.proto";`,
	), 0o644))
	require.NoError(t, afero.WriteFile(fs, "/path/subdir.proto", []byte(
		`syntax = "proto3";package whatever;import "to/service.proto";`,
	), 0o644))
	require.NoError(t, afero.WriteFile(fs, "/path/to/abs.proto", []byte(
		`syntax = "proto3";package whatever;import "/path/to/service.proto";`,
	), 0o644))

	grpcTestCase := func(expInitErr, expVUErr bool, cwd, loadCode string) multiFileTestCase ***REMOVED***
		script := tb.Replacer.Replace(fmt.Sprintf(`
			var grpc = require('k6/net/grpc');
			var client = new grpc.Client();

			%s // load statements

			exports.default = function() ***REMOVED***
				client.connect('GRPCBIN_ADDR', ***REMOVED***timeout: '3s'***REMOVED***);
				try ***REMOVED***
					var resp = client.invoke('grpc.testing.TestService/UnaryCall', ***REMOVED******REMOVED***)
					if (!resp.message || resp.error || resp.message.username !== 'foo') ***REMOVED***
						throw new Error('unexpected response message: ' + JSON.stringify(resp.message))
					***REMOVED***
				***REMOVED*** finally ***REMOVED***
					client.close();
				***REMOVED***
			***REMOVED***
		`, loadCode))

		return multiFileTestCase***REMOVED***
			fses:    map[string]afero.Fs***REMOVED***"file": fs, "https": afero.NewMemMapFs()***REMOVED***,
			rtOpts:  lib.RuntimeOptions***REMOVED***CompatibilityMode: null.NewString("base", true)***REMOVED***,
			samples: make(chan metrics.SampleContainer, 100),
			cwd:     cwd, expInitErr: expInitErr, expVUErr: expVUErr, script: script,
		***REMOVED***
	***REMOVED***

	testCases := []multiFileTestCase***REMOVED***
		grpcTestCase(false, true, "/", `/* no grpc loads */`), // exp VU error with no proto files loaded

		// Init errors when the protobuf file can't be loaded
		grpcTestCase(true, false, "/", `client.load(null, 'service.proto');`),
		grpcTestCase(true, false, "/", `client.load(null, '/wrong/path/to/service.proto');`),
		grpcTestCase(true, false, "/", `client.load(['/', '/path/'], 'service.proto');`),

		// Direct imports of service.proto
		grpcTestCase(false, false, "/", `client.load(null, '/path/to/service.proto');`), // full path should be fine
		grpcTestCase(false, false, "/path/to/", `client.load([], 'service.proto');`),    // file name from same folder
		grpcTestCase(false, false, "/", `client.load(['./path//to/'], 'service.proto');`),
		grpcTestCase(false, false, "/path/", `client.load(['./to/'], 'service.proto');`),

		grpcTestCase(false, false, "/whatever", `client.load(['/path/to/'], 'service.proto');`),  // with import paths
		grpcTestCase(false, false, "/path", `client.load(['/', '/path/to/'], 'service.proto');`), // with import paths
		grpcTestCase(false, false, "/whatever", `client.load(['../path/to/'], 'service.proto');`),

		// Import another file that imports "service.proto" directly
		grpcTestCase(true, false, "/", `client.load([], '/path/to/same-dir.proto');`),
		grpcTestCase(true, false, "/path/", `client.load([], 'to/same-dir.proto');`),
		grpcTestCase(true, false, "/", `client.load(['/path/'], 'to/same-dir.proto');`),
		grpcTestCase(false, false, "/path/to/", `client.load([], 'same-dir.proto');`),
		grpcTestCase(false, false, "/", `client.load(['/path/to/'], 'same-dir.proto');`),
		grpcTestCase(false, false, "/whatever", `client.load(['/other', '/path/to/'], 'same-dir.proto');`),
		grpcTestCase(false, false, "/", `client.load(['./path//to/'], 'same-dir.proto');`),
		grpcTestCase(false, false, "/path/", `client.load(['./to/'], 'same-dir.proto');`),
		grpcTestCase(false, false, "/whatever", `client.load(['../path/to/'], 'same-dir.proto');`),

		// Import another file that imports "to/service.proto" directly
		grpcTestCase(true, false, "/", `client.load([], '/path/to/subdir.proto');`),
		grpcTestCase(false, false, "/path/", `client.load([], 'subdir.proto');`),
		grpcTestCase(false, false, "/", `client.load(['/path/'], 'subdir.proto');`),
		grpcTestCase(false, false, "/", `client.load(['./path/'], 'subdir.proto');`),
		grpcTestCase(false, false, "/whatever", `client.load(['/other', '/path/'], 'subdir.proto');`),
		grpcTestCase(false, false, "/whatever", `client.load(['../other', '../path/'], 'subdir.proto');`),

		// Import another file that imports "/path/to/service.proto" directly
		grpcTestCase(true, false, "/", `client.load(['/path'], '/path/to/abs.proto');`),
		grpcTestCase(false, false, "/", `client.load([], '/path/to/abs.proto');`),
		grpcTestCase(false, false, "/whatever", `client.load(['/'], '/path/to/abs.proto');`),
	***REMOVED***

	for i, tc := range testCases ***REMOVED***
		i, tc := i, tc
		t.Run(fmt.Sprintf("TestCase_%d", i), func(t *testing.T) ***REMOVED***
			t.Parallel()
			t.Logf(
				"CWD: %s, expInitErr: %t, expVUErr: %t, script injected with: `%s`",
				tc.cwd, tc.expInitErr, tc.expVUErr, tc.script,
			)
			runMultiFileTestCase(t, tc, tb)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMinIterationDurationIsCancellable(t *testing.T) ***REMOVED***
	t.Parallel()

	r, err := getSimpleRunner(t, "/script.js", `
			exports.options = ***REMOVED*** iterations: 1, vus: 1, minIterationDuration: '1m' ***REMOVED***;

			exports.default = function() ***REMOVED*** /* do nothing */ ***REMOVED***;
		`)
	require.NoError(t, err)

	ch := make(chan metrics.SampleContainer, 1000)
	initVU, err := r.NewVU(1, 1, ch, lib.GetTestWorkerInfo())
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	errC := make(chan error)
	go func() ***REMOVED*** errC <- vu.RunOnce() ***REMOVED***()

	time.Sleep(200 * time.Millisecond) // give it some time to actually start

	cancel() // simulate the end of gracefulStop or a Ctrl+C event

	select ***REMOVED***
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out or minIterationDuration prevailed")
	case err := <-errC:
		require.NoError(t, err)
	***REMOVED***
***REMOVED***

//nolint:paralleltest
func TestForceHTTP1Feature(t *testing.T) ***REMOVED***
	cases := map[string]struct ***REMOVED***
		godebug               string
		expectedForceH1Result bool
		protocol              string
	***REMOVED******REMOVED***
		"Force H1 Enabled. Checking for H1": ***REMOVED***
			godebug:               "http2client=0,gctrace=1",
			expectedForceH1Result: true,
			protocol:              "HTTP/1.1",
		***REMOVED***,
		"Force H1 Disabled. Checking for H2": ***REMOVED***
			godebug:               "test=0",
			expectedForceH1Result: false,
			protocol:              "HTTP/2.0",
		***REMOVED***,
	***REMOVED***

	for name, tc := range cases ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			err := os.Setenv("GODEBUG", tc.godebug)
			require.NoError(t, err)
			defer func() ***REMOVED***
				err = os.Unsetenv("GODEBUG")
				require.NoError(t, err)
			***REMOVED***()
			assert.Equal(t, tc.expectedForceH1Result, forceHTTP1())

			tb := httpmultibin.NewHTTPMultiBin(t)

			data := fmt.Sprintf(`var k6 = require("k6");
			var check = k6.check;
			var fail = k6.fail;
			var http = require("k6/http");;
			exports.default = function() ***REMOVED***
				var res = http.get("HTTP2BIN_URL");
				if (
					!check(res, ***REMOVED***
					'checking to see if status was 200': (res) => res.status === 200,
					'checking to see protocol': (res) => res.proto === '%s'
					***REMOVED***)
				) ***REMOVED***
					fail('test failed')
				***REMOVED***
			***REMOVED***`, tc.protocol)

			r1, err := getSimpleRunner(t, "/script.js", tb.Replacer.Replace(data))
			require.NoError(t, err)

			err = r1.SetOptions(lib.Options***REMOVED***
				Hosts: tb.Dialer.Hosts,
				// We disable TLS verify so that we don't get a TLS handshake error since
				// the certificates on the endpoint are not certified by a certificate authority
				InsecureSkipTLSVerify: null.BoolFrom(true),
			***REMOVED***)

			require.NoError(t, err)

			registry := metrics.NewRegistry()
			builtinMetrics := metrics.RegisterBuiltinMetrics(registry)
			r2, err := NewFromArchive(
				&lib.TestPreInitState***REMOVED***
					Logger:         testutils.NewLogger(t),
					BuiltinMetrics: builtinMetrics,
					Registry:       registry,
				***REMOVED***, r1.MakeArchive(), lib.GetTestWorkerInfo())
			require.NoError(t, err)

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				r := r
				t.Run(name, func(t *testing.T) ***REMOVED***
					initVU, err := r.NewVU(1, 1, make(chan metrics.SampleContainer, 100), lib.GetTestWorkerInfo())
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
***REMOVED***

func TestExecutionInfo(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		name, script, expErr string
	***REMOVED******REMOVED***
		***REMOVED***name: "vu_ok", script: `
		var exec = require('k6/execution');

		exports.default = function() ***REMOVED***
			if (exec.vu.idInInstance !== 1) throw new Error('unexpected VU ID: '+exec.vu.idInInstance);
			if (exec.vu.idInTest !== 10) throw new Error('unexpected global VU ID: '+exec.vu.idInTest);
			if (exec.vu.iterationInInstance !== 0) throw new Error('unexpected VU iteration: '+exec.vu.iterationInInstance);
			if (exec.vu.iterationInScenario !== 0) throw new Error('unexpected scenario iteration: '+exec.vu.iterationInScenario);
		***REMOVED***`***REMOVED***,
		***REMOVED***name: "vu_err", script: `
		var exec = require('k6/execution');
		exec.vu;
		`, expErr: "getting VU information in the init context is not supported"***REMOVED***,
		***REMOVED***name: "scenario_ok", script: `
		var exec = require('k6/execution');
		var sleep = require('k6').sleep;

		exports.default = function() ***REMOVED***
			var si = exec.scenario;
			sleep(0.1);
			if (si.name !== 'default') throw new Error('unexpected scenario name: '+si.name);
			if (si.executor !== 'test-exec') throw new Error('unexpected executor: '+si.executor);
			if (si.startTime > new Date().getTime()) throw new Error('unexpected startTime: '+si.startTime);
			if (si.progress !== 0.1) throw new Error('unexpected progress: '+si.progress);
			if (si.iterationInInstance !== 3) throw new Error('unexpected scenario local iteration: '+si.iterationInInstance);
			if (si.iterationInTest !== 4) throw new Error('unexpected scenario local iteration: '+si.iterationInTest);
		***REMOVED***`***REMOVED***,
		***REMOVED***name: "scenario_err", script: `
		var exec = require('k6/execution');
		exec.scenario;
		`, expErr: "getting scenario information outside of the VU context is not supported"***REMOVED***,
		***REMOVED***name: "test_ok", script: `
		var exec = require('k6/execution');

		exports.default = function() ***REMOVED***
			var ti = exec.instance;
			if (ti.currentTestRunDuration !== 0) throw new Error('unexpected test duration: '+ti.currentTestRunDuration);
			if (ti.vusActive !== 1) throw new Error('unexpected vusActive: '+ti.vusActive);
			if (ti.vusInitialized !== 0) throw new Error('unexpected vusInitialized: '+ti.vusInitialized);
			if (ti.iterationsCompleted !== 0) throw new Error('unexpected iterationsCompleted: '+ti.iterationsCompleted);
			if (ti.iterationsInterrupted !== 0) throw new Error('unexpected iterationsInterrupted: '+ti.iterationsInterrupted);
		***REMOVED***`***REMOVED***,
		***REMOVED***name: "test_err", script: `
		var exec = require('k6/execution');
		exec.instance;
		`, expErr: "getting instance information in the init context is not supported"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			r, err := getSimpleRunner(t, "/script.js", tc.script)
			if tc.expErr != "" ***REMOVED***
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expErr)
				return
			***REMOVED***
			require.NoError(t, err)

			r.Bundle.Options.SystemTags = metrics.NewSystemTagSet(metrics.DefaultSystemTagSet)
			samples := make(chan metrics.SampleContainer, 100)
			initVU, err := r.NewVU(1, 10, samples, lib.GetTestWorkerInfo())
			require.NoError(t, err)

			testRunState := &lib.TestRunState***REMOVED***
				TestPreInitState: r.preInitState,
				Options:          r.GetOptions(),
				Runner:           r,
			***REMOVED***

			execScheduler, err := local.NewExecutionScheduler(testRunState)
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctx = lib.WithExecutionState(ctx, execScheduler.GetState())
			ctx = lib.WithScenarioState(ctx, &lib.ScenarioState***REMOVED***
				Name:      "default",
				Executor:  "test-exec",
				StartTime: time.Now(),
				ProgressFn: func() (float64, []string) ***REMOVED***
					return 0.1, nil
				***REMOVED***,
			***REMOVED***)
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***
				RunContext:               ctx,
				Exec:                     "default",
				GetNextIterationCounters: func() (uint64, uint64) ***REMOVED*** return 3, 4 ***REMOVED***,
			***REMOVED***)

			execState := execScheduler.GetState()
			execState.ModCurrentlyActiveVUsCount(+1)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***
