package execution

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/APITeamLimited/k6-worker/errext"
	"github.com/APITeamLimited/k6-worker/js/common"
	"github.com/APITeamLimited/k6-worker/js/modulestest"
	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/lib/executor"
	"github.com/APITeamLimited/k6-worker/lib/testutils"
	"github.com/APITeamLimited/k6-worker/lib/types"
	"github.com/APITeamLimited/k6-worker/metrics"
	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func setupTagsExecEnv(t *testing.T) *modulestest.Runtime ***REMOVED***
	testRuntime := modulestest.NewRuntime(t)
	m, ok := New().NewModuleInstance(testRuntime.VU).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, testRuntime.VU.Runtime().Set("exec", m.Exports().Default))

	return testRuntime
***REMOVED***

func TestVUTags(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("Get", func(t *testing.T) ***REMOVED***
		t.Parallel()

		tenv := setupTagsExecEnv(t)
		tenv.MoveToVUContext(&lib.State***REMOVED***Tags: lib.NewTagMap(map[string]string***REMOVED***"vu": "42"***REMOVED***)***REMOVED***)
		tag, err := tenv.VU.Runtime().RunString(`exec.vu.tags["vu"]`)
		require.NoError(t, err)
		assert.Equal(t, "42", tag.String())

		// not found
		tag, err = tenv.VU.Runtime().RunString(`exec.vu.tags["not-existing-tag"]`)
		require.NoError(t, err)
		assert.Equal(t, "undefined", tag.String())
	***REMOVED***)

	t.Run("JSONEncoding", func(t *testing.T) ***REMOVED***
		t.Parallel()

		tenv := setupTagsExecEnv(t)
		tenv.MoveToVUContext(&lib.State***REMOVED***
			Options: lib.Options***REMOVED***
				SystemTags: metrics.NewSystemTagSet(metrics.TagVU),
			***REMOVED***,
			Tags: lib.NewTagMap(map[string]string***REMOVED***
				"vu": "42",
			***REMOVED***),
		***REMOVED***)
		state := tenv.VU.State()
		state.Tags.Set("custom-tag", "mytag1")

		encoded, err := tenv.VU.Runtime().RunString(`JSON.stringify(exec.vu.tags)`)
		require.NoError(t, err)
		assert.JSONEq(t, `***REMOVED***"vu":"42","custom-tag":"mytag1"***REMOVED***`, encoded.String())
	***REMOVED***)

	t.Run("Set", func(t *testing.T) ***REMOVED***
		t.Parallel()

		t.Run("SuccessAccetedTypes", func(t *testing.T) ***REMOVED***
			t.Parallel()

			// bool and numbers are implicitly converted into string

			tests := map[string]struct ***REMOVED***
				v   interface***REMOVED******REMOVED***
				exp string
			***REMOVED******REMOVED***
				"string": ***REMOVED***v: `"tag1"`, exp: "tag1"***REMOVED***,
				"bool":   ***REMOVED***v: true, exp: "true"***REMOVED***,
				"int":    ***REMOVED***v: 101, exp: "101"***REMOVED***,
				"float":  ***REMOVED***v: 3.14, exp: "3.14"***REMOVED***,
			***REMOVED***

			tenv := setupTagsExecEnv(t)
			tenv.MoveToVUContext(&lib.State***REMOVED***Tags: lib.NewTagMap(map[string]string***REMOVED***"vu": "42"***REMOVED***)***REMOVED***)

			for _, tc := range tests ***REMOVED***
				_, err := tenv.VU.Runtime().RunString(fmt.Sprintf(`exec.vu.tags["mytag"] = %v`, tc.v))
				require.NoError(t, err)

				val, err := tenv.VU.Runtime().RunString(`exec.vu.tags["mytag"]`)
				require.NoError(t, err)

				assert.Equal(t, tc.exp, val.String())
			***REMOVED***
		***REMOVED***)

		t.Run("SuccessOverwriteSystemTag", func(t *testing.T) ***REMOVED***
			t.Parallel()

			tenv := setupTagsExecEnv(t)
			tenv.MoveToVUContext(&lib.State***REMOVED***Tags: lib.NewTagMap(map[string]string***REMOVED***"vu": "42"***REMOVED***)***REMOVED***)

			_, err := tenv.VU.Runtime().RunString(`exec.vu.tags["vu"] = "vu101"`)
			require.NoError(t, err)
			val, err := tenv.VU.Runtime().RunString(`exec.vu.tags["vu"]`)
			require.NoError(t, err)
			assert.Equal(t, "vu101", val.String())
		***REMOVED***)

		t.Run("DiscardWrongTypeAndRaisingError", func(t *testing.T) ***REMOVED***
			t.Parallel()

			tenv := setupTagsExecEnv(t)
			tenv.MoveToVUContext(&lib.State***REMOVED***Tags: lib.NewTagMap(map[string]string***REMOVED***"vu": "42"***REMOVED***)***REMOVED***)

			state := tenv.VU.State()
			state.Options.Throw = null.BoolFrom(true)
			require.NotNil(t, state)

			cases := []string***REMOVED***
				`[1, 3, 5]`,             // array
				`***REMOVED***f1: "value1", f2: 4***REMOVED***`, // object
			***REMOVED***

			for _, val := range cases ***REMOVED***
				_, err := tenv.VU.Runtime().RunString(`exec.vu.tags["custom-tag"] = ` + val)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "TypeError:")
				assert.Contains(t, err.Error(), "only String, Boolean and Number")
			***REMOVED***
		***REMOVED***)

		t.Run("DiscardWrongTypeOnlyWarning", func(t *testing.T) ***REMOVED***
			t.Parallel()
			logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
			testLog := logrus.New()
			testLog.AddHook(logHook)
			testLog.SetOutput(ioutil.Discard)

			tenv := setupTagsExecEnv(t)
			tenv.MoveToVUContext(&lib.State***REMOVED***
				Options: lib.Options***REMOVED***
					SystemTags: metrics.NewSystemTagSet(metrics.TagVU),
				***REMOVED***,
				Tags: lib.NewTagMap(map[string]string***REMOVED***
					"vu": "42",
				***REMOVED***),
				Logger: testLog,
			***REMOVED***)
			_, err := tenv.VU.Runtime().RunString(`exec.vu.tags["custom-tag"] = [1, 3, 5]`)
			require.NoError(t, err)

			entries := logHook.Drain()
			require.Len(t, entries, 1)
			assert.Contains(t, entries[0].Message, "discarded")
		***REMOVED***)

		t.Run("DiscardNullOrUndefined", func(t *testing.T) ***REMOVED***
			t.Parallel()

			logHook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
			testLog := logrus.New()
			testLog.AddHook(logHook)
			testLog.SetOutput(ioutil.Discard)

			cases := []string***REMOVED***"null", "undefined"***REMOVED***
			tenv := setupTagsExecEnv(t)
			tenv.MoveToVUContext(&lib.State***REMOVED***
				Options: lib.Options***REMOVED***
					SystemTags: metrics.NewSystemTagSet(metrics.TagVU),
				***REMOVED***,
				Tags:   lib.NewTagMap(map[string]string***REMOVED***"vu": "42"***REMOVED***),
				Logger: testLog,
			***REMOVED***)
			for _, val := range cases ***REMOVED***
				_, err := tenv.VU.Runtime().RunString(`exec.vu.tags["custom-tag"] = ` + val)
				require.NoError(t, err)

				entries := logHook.Drain()
				require.Len(t, entries, 1)
				assert.Contains(t, entries[0].Message, "discarded")
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestAbortTest(t *testing.T) ***REMOVED*** //nolint:tparallel
	t.Parallel()

	var (
		rt    = goja.New()
		state = &lib.State***REMOVED******REMOVED***
		ctx   = context.Background()
	)

	m, ok := New().NewModuleInstance(
		&modulestest.VU***REMOVED***
			RuntimeField: rt,
			InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
			CtxField:     ctx,
			StateField:   state,
		***REMOVED***,
	).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("exec", m.Exports().Default))

	prove := func(t *testing.T, script, reason string) ***REMOVED***
		_, err := rt.RunString(script)
		require.NotNil(t, err)
		var x *goja.InterruptedError
		assert.ErrorAs(t, err, &x)
		v, ok := x.Value().(*errext.InterruptError)
		require.True(t, ok)
		require.Equal(t, v.Reason, reason)
	***REMOVED***

	t.Run("default reason", func(t *testing.T) ***REMOVED*** //nolint:paralleltest
		prove(t, "exec.test.abort()", errext.AbortTest)
	***REMOVED***)
	t.Run("custom reason", func(t *testing.T) ***REMOVED*** //nolint:paralleltest
		prove(t, `exec.test.abort("mayday")`, fmt.Sprintf("%s: mayday", errext.AbortTest))
	***REMOVED***)
***REMOVED***

func TestOptionsTestFull(t *testing.T) ***REMOVED***
	t.Parallel()

	expected := `***REMOVED***"paused":true,"scenarios":***REMOVED***"const-vus":***REMOVED***"executor":"constant-vus","startTime":"10s","gracefulStop":"30s","env":***REMOVED***"FOO":"bar"***REMOVED***,"exec":"default","tags":***REMOVED***"tagkey":"tagvalue"***REMOVED***,"vus":50,"duration":"10m0s"***REMOVED******REMOVED***,"executionSegment":"0:1/4","executionSegmentSequence":"0,1/4,1/2,1","noSetup":true,"setupTimeout":"1m0s","noTeardown":true,"teardownTimeout":"5m0s","rps":100,"dns":***REMOVED***"ttl":"1m","select":"roundRobin","policy":"any"***REMOVED***,"maxRedirects":3,"userAgent":"k6-user-agent","batch":15,"batchPerHost":5,"httpDebug":"full","insecureSkipTLSVerify":true,"tlsCipherSuites":["TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"],"tlsVersion":***REMOVED***"min":"tls1.2","max":"tls1.3"***REMOVED***,"tlsAuth":[***REMOVED***"domains":["example.com"],"cert":"mycert.pem","key":"mycert-key.pem","password":"mypwd"***REMOVED***],"throw":true,"thresholds":***REMOVED***"http_req_duration":[***REMOVED***"threshold":"rate>0.01","abortOnFail":true,"delayAbortEval":"10s"***REMOVED***]***REMOVED***,"blacklistIPs":["192.0.2.0/24"],"blockHostnames":["test.k6.io","*.example.com"],"hosts":***REMOVED***"test.k6.io":"1.2.3.4:8443"***REMOVED***,"noConnectionReuse":true,"noVUConnectionReuse":true,"minIterationDuration":"10s","ext":***REMOVED***"ext-one":***REMOVED***"rawkey":"rawvalue"***REMOVED******REMOVED***,"summaryTrendStats":["avg","min","max"],"summaryTimeUnit":"ms","systemTags":["iter","vu"],"tags":null,"metricSamplesBufferSize":8,"noCookiesReset":true,"discardResponseBodies":true,"consoleOutput":"loadtest.log","tags":***REMOVED***"runtag-key":"runtag-value"***REMOVED***,"localIPs":"192.168.20.12-192.168.20.15,192.168.10.0/27"***REMOVED***`

	var (
		rt    = goja.New()
		state = &lib.State***REMOVED***
			Options: lib.Options***REMOVED***
				Paused: null.BoolFrom(true),
				Scenarios: map[string]lib.ExecutorConfig***REMOVED***
					"const-vus": executor.ConstantVUsConfig***REMOVED***
						BaseConfig: executor.BaseConfig***REMOVED***
							Name:         "const-vus",
							Type:         "constant-vus",
							StartTime:    types.NullDurationFrom(10 * time.Second),
							GracefulStop: types.NullDurationFrom(30 * time.Second),
							Env: map[string]string***REMOVED***
								"FOO": "bar",
							***REMOVED***,
							Exec: null.StringFrom("default"),
							Tags: map[string]string***REMOVED***
								"tagkey": "tagvalue",
							***REMOVED***,
						***REMOVED***,
						VUs:      null.IntFrom(50),
						Duration: types.NullDurationFrom(10 * time.Minute),
					***REMOVED***,
				***REMOVED***,
				ExecutionSegment: func() *lib.ExecutionSegment ***REMOVED***
					seg, err := lib.NewExecutionSegmentFromString("0:1/4")
					require.NoError(t, err)
					return seg
				***REMOVED***(),
				ExecutionSegmentSequence: func() *lib.ExecutionSegmentSequence ***REMOVED***
					seq, err := lib.NewExecutionSegmentSequenceFromString("0,1/4,1/2,1")
					require.NoError(t, err)
					return &seq
				***REMOVED***(),
				NoSetup:               null.BoolFrom(true),
				NoTeardown:            null.BoolFrom(true),
				NoConnectionReuse:     null.BoolFrom(true),
				NoVUConnectionReuse:   null.BoolFrom(true),
				InsecureSkipTLSVerify: null.BoolFrom(true),
				Throw:                 null.BoolFrom(true),
				NoCookiesReset:        null.BoolFrom(true),
				DiscardResponseBodies: null.BoolFrom(true),
				RPS:                   null.IntFrom(100),
				MaxRedirects:          null.IntFrom(3),
				UserAgent:             null.StringFrom("k6-user-agent"),
				Batch:                 null.IntFrom(15),
				BatchPerHost:          null.IntFrom(5),
				SetupTimeout:          types.NullDurationFrom(1 * time.Minute),
				TeardownTimeout:       types.NullDurationFrom(5 * time.Minute),
				MinIterationDuration:  types.NullDurationFrom(10 * time.Second),
				HTTPDebug:             null.StringFrom("full"),
				DNS: types.DNSConfig***REMOVED***
					TTL:    null.StringFrom("1m"),
					Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSroundRobin, Valid: true***REMOVED***,
					Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSany, Valid: true***REMOVED***,
					Valid:  true,
				***REMOVED***,
				TLSVersion: &lib.TLSVersions***REMOVED***
					Min: tls.VersionTLS12,
					Max: tls.VersionTLS13,
				***REMOVED***,
				TLSAuth: []*lib.TLSAuth***REMOVED***
					***REMOVED***
						TLSAuthFields: lib.TLSAuthFields***REMOVED***
							Cert:     "mycert.pem",
							Key:      "mycert-key.pem",
							Password: null.StringFrom("mypwd"),
							Domains:  []string***REMOVED***"example.com"***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				TLSCipherSuites: &lib.TLSCipherSuites***REMOVED***
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				***REMOVED***,
				BlacklistIPs: []*lib.IPNet***REMOVED***
					***REMOVED***
						IPNet: func() net.IPNet ***REMOVED***
							_, ipv4net, err := net.ParseCIDR("192.0.2.1/24")
							require.NoError(t, err)
							return *ipv4net
						***REMOVED***(),
					***REMOVED***,
				***REMOVED***,
				Thresholds: map[string]metrics.Thresholds***REMOVED***
					"http_req_duration": ***REMOVED***
						Thresholds: []*metrics.Threshold***REMOVED***
							***REMOVED***
								Source:           "rate>0.01",
								LastFailed:       true,
								AbortOnFail:      true,
								AbortGracePeriod: types.NullDurationFrom(10 * time.Second),
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				BlockedHostnames: func() types.NullHostnameTrie ***REMOVED***
					bh, err := types.NewNullHostnameTrie([]string***REMOVED***"test.k6.io", "*.example.com"***REMOVED***)
					require.NoError(t, err)
					return bh
				***REMOVED***(),
				Hosts: map[string]*lib.HostAddress***REMOVED***
					"test.k6.io": ***REMOVED***
						IP:   []byte***REMOVED***0x01, 0x02, 0x03, 0x04***REMOVED***,
						Port: 8443,
					***REMOVED***,
				***REMOVED***,
				External: map[string]json.RawMessage***REMOVED***
					"ext-one": json.RawMessage(`***REMOVED***"rawkey":"rawvalue"***REMOVED***`),
				***REMOVED***,
				SummaryTrendStats: []string***REMOVED***"avg", "min", "max"***REMOVED***,
				SummaryTimeUnit:   null.StringFrom("ms"),
				SystemTags: func() *metrics.SystemTagSet ***REMOVED***
					sysm := metrics.TagIter | metrics.TagVU
					return &sysm
				***REMOVED***(),
				RunTags:                 map[string]string***REMOVED***"runtag-key": "runtag-value"***REMOVED***,
				MetricSamplesBufferSize: null.IntFrom(8),
				ConsoleOutput:           null.StringFrom("loadtest.log"),
				LocalIPs: func() types.NullIPPool ***REMOVED***
					npool := types.NullIPPool***REMOVED******REMOVED***
					err := npool.UnmarshalText([]byte("192.168.20.12-192.168.20.15,192.168.10.0/27"))
					require.NoError(t, err)
					return npool
				***REMOVED***(),

				// The following fields are not expected to be
				// in the final test.options object
				VUs:        null.IntFrom(50),
				Iterations: null.IntFrom(100),
				Duration:   types.NullDurationFrom(10 * time.Second),
				Stages: []lib.Stage***REMOVED***
					***REMOVED***
						Duration: types.NullDurationFrom(2 * time.Second),
						Target:   null.IntFrom(2),
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		ctx = context.Background()
	)

	m, ok := New().NewModuleInstance(
		&modulestest.VU***REMOVED***
			RuntimeField: rt,
			CtxField:     ctx,
			StateField:   state,
		***REMOVED***,
	).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("exec", m.Exports().Default))

	opts, err := rt.RunString(`JSON.stringify(exec.test.options)`)
	require.NoError(t, err)
	require.NotNil(t, opts)
	assert.JSONEq(t, expected, opts.String())
***REMOVED***

func TestOptionsTestSetPropertyDenied(t *testing.T) ***REMOVED***
	t.Parallel()

	rt := goja.New()
	m, ok := New().NewModuleInstance(
		&modulestest.VU***REMOVED***
			RuntimeField: rt,
			CtxField:     context.Background(),
			StateField: &lib.State***REMOVED***
				Options: lib.Options***REMOVED***
					Paused: null.BoolFrom(true),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("exec", m.Exports().Default))

	_, err := rt.RunString(`exec.test.options.paused = false`)
	require.NoError(t, err)
	paused, err := rt.RunString(`exec.test.options.paused`)
	require.NoError(t, err)
	assert.Equal(t, true, rt.ToValue(paused).ToBoolean())
***REMOVED***

func TestScenarioNoAvailableInInitContext(t *testing.T) ***REMOVED***
	t.Parallel()

	rt := goja.New()
	m, ok := New().NewModuleInstance(
		&modulestest.VU***REMOVED***
			RuntimeField: rt,
			CtxField:     context.Background(),
			StateField: &lib.State***REMOVED***
				Options: lib.Options***REMOVED***
					Paused: null.BoolFrom(true),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	).(*ModuleInstance)
	require.True(t, ok)
	require.NoError(t, rt.Set("exec", m.Exports().Default))

	scenarioExportedProps := []string***REMOVED***"name", "executor", "startTime", "progress", "iterationInInstance", "iterationInTest"***REMOVED***

	for _, code := range scenarioExportedProps ***REMOVED***
		prop := fmt.Sprintf("exec.scenario.%s", code)
		_, err := rt.RunString(prop)
		require.Error(t, err)
		require.ErrorContains(t, err, "getting scenario information outside of the VU context is not supported")
	***REMOVED***
***REMOVED***
