package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/httpmultibin"
)

const (
	noopDefaultFunc   = `export default function() ***REMOVED******REMOVED***;`
	fooLogDefaultFunc = `export default function() ***REMOVED*** console.log('foo'); ***REMOVED***;`
	noopHandleSummary = `
		export function handleSummary(data) ***REMOVED***
			return ***REMOVED******REMOVED***; // silence the end of test summary
		***REMOVED***;
	`
)

func TestVersion(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)
	ts.args = []string***REMOVED***"k6", "version"***REMOVED***
	newRootCommand(ts.globalState).execute()

	stdOut := ts.stdOut.String()
	assert.Contains(t, stdOut, "k6 v"+consts.Version)
	assert.Contains(t, stdOut, runtime.Version())
	assert.Contains(t, stdOut, runtime.GOOS)
	assert.Contains(t, stdOut, runtime.GOARCH)
	assert.NotContains(t, stdOut[:len(stdOut)-1], "\n")

	assert.Empty(t, ts.stdErr.Bytes())
	assert.Empty(t, ts.loggerHook.Drain())
***REMOVED***

func TestSimpleTestStdin(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)
	ts.args = []string***REMOVED***"k6", "run", "-"***REMOVED***
	ts.stdIn = bytes.NewBufferString(noopDefaultFunc)
	newRootCommand(ts.globalState).execute()

	stdOut := ts.stdOut.String()
	assert.Contains(t, stdOut, "default: 1 iterations for each of 1 VUs")
	assert.Contains(t, stdOut, "1 complete and 0 interrupted iterations")
	assert.Empty(t, ts.stdErr.Bytes())
	assert.Empty(t, ts.loggerHook.Drain())
***REMOVED***

func TestStdoutAndStderrAreEmptyWithQuietAndHandleSummary(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)
	ts.args = []string***REMOVED***"k6", "--quiet", "run", "-"***REMOVED***
	ts.stdIn = bytes.NewBufferString(noopDefaultFunc + noopHandleSummary)
	newRootCommand(ts.globalState).execute()

	assert.Empty(t, ts.stdErr.Bytes())
	assert.Empty(t, ts.stdOut.Bytes())
	assert.Empty(t, ts.loggerHook.Drain())
***REMOVED***

func TestStdoutAndStderrAreEmptyWithQuietAndLogsForwarded(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)

	// TODO: add a test with relative path
	logFilePath := filepath.Join(ts.cwd, "test.log")

	ts.args = []string***REMOVED***
		"k6", "--quiet", "--log-output", "file=" + logFilePath,
		"--log-format", "raw", "run", "--no-summary", "-",
	***REMOVED***
	ts.stdIn = bytes.NewBufferString(fooLogDefaultFunc)
	newRootCommand(ts.globalState).execute()

	// The test state hook still catches this message
	assert.True(t, testutils.LogContains(ts.loggerHook.Drain(), logrus.InfoLevel, `foo`))

	// But it's not shown on stderr or stdout
	assert.Empty(t, ts.stdErr.Bytes())
	assert.Empty(t, ts.stdOut.Bytes())

	// Instead it should be in the log file
	logContents, err := afero.ReadFile(ts.fs, logFilePath)
	require.NoError(t, err)
	assert.Equal(t, "foo\n", string(logContents))
***REMOVED***

func TestRelativeLogPathWithSetupAndTeardown(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)

	ts.args = []string***REMOVED***"k6", "--log-output", "file=test.log", "--log-format", "raw", "run", "-i", "2", "-"***REMOVED***
	ts.stdIn = bytes.NewBufferString(fooLogDefaultFunc + `
		export function setup() ***REMOVED*** console.log('bar'); ***REMOVED***;
		export function teardown() ***REMOVED*** console.log('baz'); ***REMOVED***;
	`)
	newRootCommand(ts.globalState).execute()

	// The test state hook still catches these messages
	logEntries := ts.loggerHook.Drain()
	assert.True(t, testutils.LogContains(logEntries, logrus.InfoLevel, `foo`))
	assert.True(t, testutils.LogContains(logEntries, logrus.InfoLevel, `bar`))
	assert.True(t, testutils.LogContains(logEntries, logrus.InfoLevel, `baz`))

	// And check that the log file also contains everything
	logContents, err := afero.ReadFile(ts.fs, filepath.Join(ts.cwd, "test.log"))
	require.NoError(t, err)
	assert.Equal(t, "bar\nfoo\nfoo\nbaz\n", string(logContents))
***REMOVED***

func TestWrongCliFlagIterations(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)
	ts.args = []string***REMOVED***"k6", "run", "--iterations", "foo", "-"***REMOVED***
	ts.stdIn = bytes.NewBufferString(noopDefaultFunc)
	// TODO: check for exitcodes.InvalidConfig after https://github.com/loadimpact/k6/issues/883 is done...
	ts.expectedExitCode = -1
	newRootCommand(ts.globalState).execute()
	assert.True(t, testutils.LogContains(ts.loggerHook.Drain(), logrus.ErrorLevel, `invalid argument "foo"`))
***REMOVED***

func TestWrongEnvVarIterations(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)
	ts.args = []string***REMOVED***"k6", "run", "--vus", "2", "-"***REMOVED***
	ts.envVars = map[string]string***REMOVED***"K6_ITERATIONS": "4"***REMOVED***
	ts.stdIn = bytes.NewBufferString(noopDefaultFunc)

	newRootCommand(ts.globalState).execute()

	stdOut := ts.stdOut.String()
	t.Logf(stdOut)
	assert.Contains(t, stdOut, "4 iterations shared among 2 VUs")
	assert.Contains(t, stdOut, "4 complete and 0 interrupted iterations")
	assert.Empty(t, ts.stdErr.Bytes())
	assert.Empty(t, ts.loggerHook.Drain())
***REMOVED***

func TestMetricsAndThresholds(t *testing.T) ***REMOVED***
	t.Parallel()
	script := `
		import ***REMOVED*** Counter ***REMOVED*** from 'k6/metrics';

		var setupCounter = new Counter('setup_counter');
		var teardownCounter = new Counter('teardown_counter');
		var defaultCounter = new Counter('default_counter');
		let unusedCounter = new Counter('unused_counter');

		export const options = ***REMOVED***
			scenarios: ***REMOVED***
				sc1: ***REMOVED***
					executor: 'per-vu-iterations',
					vus: 1,
					iterations: 1,
				***REMOVED***,
				sc2: ***REMOVED***
					executor: 'shared-iterations',
					vus: 1,
					iterations: 1,
				***REMOVED***,
			***REMOVED***,
			thresholds: ***REMOVED***
				'setup_counter': ['count == 1'],
				'teardown_counter': ['count == 1'],
				'default_counter': ['count == 2'],
				'default_counter***REMOVED***scenario:sc1***REMOVED***': ['count == 1'],
				'default_counter***REMOVED***scenario:sc2***REMOVED***': ['count == 1'],
				'iterations': ['count == 2'],
				'iterations***REMOVED***scenario:sc1***REMOVED***': ['count == 1'],
				'iterations***REMOVED***scenario:sc2***REMOVED***': ['count == 1'],
				'default_counter***REMOVED***nonexistent:tag***REMOVED***': ['count == 0'],
				'unused_counter': ['count == 0'],
				'http_req_duration***REMOVED***status:200***REMOVED***': [' max == 0'], // no HTTP requests
			***REMOVED***,
		***REMOVED***;

		export function setup() ***REMOVED***
			console.log('setup() start');
			setupCounter.add(1);
			console.log('setup() end');
			return ***REMOVED*** foo: 'bar' ***REMOVED***
		***REMOVED***

		export default function (data) ***REMOVED***
			console.log('default(' + JSON.stringify(data) + ')');
			defaultCounter.add(1);
		***REMOVED***

		export function teardown(data) ***REMOVED***
			console.log('teardown(' + JSON.stringify(data) + ')');
			teardownCounter.add(1);
		***REMOVED***

		export function handleSummary(data) ***REMOVED***
			console.log('handleSummary()');
			return ***REMOVED*** stdout: JSON.stringify(data, null, 4) ***REMOVED***
		***REMOVED***
	`
	ts := newGlobalTestState(t)
	require.NoError(t, afero.WriteFile(ts.fs, filepath.Join(ts.cwd, "test.js"), []byte(script), 0o644))
	ts.args = []string***REMOVED***"k6", "run", "--quiet", "--log-format=raw", "test.js"***REMOVED***

	newRootCommand(ts.globalState).execute()

	expLogLines := []string***REMOVED***
		`setup() start`, `setup() end`, `default(***REMOVED***"foo":"bar"***REMOVED***)`,
		`default(***REMOVED***"foo":"bar"***REMOVED***)`, `teardown(***REMOVED***"foo":"bar"***REMOVED***)`, `handleSummary()`,
	***REMOVED***

	logHookEntries := ts.loggerHook.Drain()
	require.Len(t, logHookEntries, len(expLogLines))
	for i, expLogLine := range expLogLines ***REMOVED***
		assert.Equal(t, expLogLine, logHookEntries[i].Message)
	***REMOVED***

	assert.Equal(t, strings.Join(expLogLines, "\n")+"\n", ts.stdErr.String())

	var summary map[string]interface***REMOVED******REMOVED***
	require.NoError(t, json.Unmarshal(ts.stdOut.Bytes(), &summary))

	metrics, ok := summary["metrics"].(map[string]interface***REMOVED******REMOVED***)
	require.True(t, ok)

	teardownCounter, ok := metrics["teardown_counter"].(map[string]interface***REMOVED******REMOVED***)
	require.True(t, ok)

	teardownThresholds, ok := teardownCounter["thresholds"].(map[string]interface***REMOVED******REMOVED***)
	require.True(t, ok)

	expected := map[string]interface***REMOVED******REMOVED******REMOVED***"count == 1": map[string]interface***REMOVED******REMOVED******REMOVED***"ok": true***REMOVED******REMOVED***
	require.Equal(t, expected, teardownThresholds)
***REMOVED***

func TestSSLKEYLOGFILE(t *testing.T) ***REMOVED***
	t.Parallel()

	// TODO don't use insecureSkipTLSVerify when/if tlsConfig is given to the runner from outside
	tb := httpmultibin.NewHTTPMultiBin(t)
	ts := newGlobalTestState(t)
	ts.args = []string***REMOVED***"k6", "run", "-"***REMOVED***
	ts.envVars = map[string]string***REMOVED***"SSLKEYLOGFILE": "./ssl.log"***REMOVED***
	ts.stdIn = bytes.NewReader([]byte(tb.Replacer.Replace(`
    import http from "k6/http"
    export const options = ***REMOVED***
      hosts: ***REMOVED***
        "HTTPSBIN_DOMAIN": "HTTPSBIN_IP",
      ***REMOVED***,
      insecureSkipTLSVerify: true,
    ***REMOVED***

    export default () => ***REMOVED***
      http.get("HTTPSBIN_URL/get");
    ***REMOVED***
  `)))

	newRootCommand(ts.globalState).execute()

	assert.True(t,
		testutils.LogContains(ts.loggerHook.Drain(), logrus.WarnLevel, "SSLKEYLOGFILE was specified"))
	sslloglines, err := afero.ReadFile(ts.fs, filepath.Join(ts.cwd, "ssl.log"))
	require.NoError(t, err)
	// TODO maybe have multiple depending on the ciphers used as that seems to change it
	require.Regexp(t, "^CLIENT_[A-Z_]+ [0-9a-f]+ [0-9a-f]+\n", string(sslloglines))
***REMOVED***

// TODO: add a hell of a lot more integration tests, including some that spin up
// a test HTTP server and actually check if k6 hits it

// TODO: also add a test that starts multiple k6 "instances", for example:
//  - one with `k6 run --paused` and another with `k6 resume`
//  - one with `k6 run` and another with `k6 stats` or `k6 status`

func TestExecutionTestOptionsDefaultValues(t *testing.T) ***REMOVED***
	t.Parallel()
	script := `
		import exec from 'k6/execution';

		export default function () ***REMOVED***
			console.log(exec.test.options)
		***REMOVED***
	`

	ts := newGlobalTestState(t)
	require.NoError(t, afero.WriteFile(ts.fs, filepath.Join(ts.cwd, "test.js"), []byte(script), 0o644))
	ts.args = []string***REMOVED***"k6", "run", "--iterations", "1", "test.js"***REMOVED***

	newRootCommand(ts.globalState).execute()

	loglines := ts.loggerHook.Drain()
	require.Len(t, loglines, 1)

	expected := `***REMOVED***"paused":null,"executionSegment":null,"executionSegmentSequence":null,"noSetup":null,"setupTimeout":null,"noTeardown":null,"teardownTimeout":null,"rps":null,"dns":***REMOVED***"ttl":null,"select":null,"policy":null***REMOVED***,"maxRedirects":null,"userAgent":null,"batch":null,"batchPerHost":null,"httpDebug":null,"insecureSkipTLSVerify":null,"tlsCipherSuites":null,"tlsVersion":null,"tlsAuth":null,"throw":null,"thresholds":null,"blacklistIPs":null,"blockHostnames":null,"hosts":null,"noConnectionReuse":null,"noVUConnectionReuse":null,"minIterationDuration":null,"ext":null,"summaryTrendStats":["avg", "min", "med", "max", "p(90)", "p(95)"],"summaryTimeUnit":null,"systemTags":["check","error","error_code","expected_response","group","method","name","proto","scenario","service","status","subproto","tls_version","url"],"tags":null,"metricSamplesBufferSize":null,"noCookiesReset":null,"discardResponseBodies":null,"consoleOutput":null,"scenarios":***REMOVED***"default":***REMOVED***"vus":null,"iterations":1,"executor":"shared-iterations","maxDuration":null,"startTime":null,"env":null,"tags":null,"gracefulStop":null,"exec":null***REMOVED******REMOVED***,"localIPs":null***REMOVED***`
	assert.JSONEq(t, expected, loglines[0].Message)
***REMOVED***
