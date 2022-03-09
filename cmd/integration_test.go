package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/lib/testutils"
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

// TODO: add a hell of a lot more integration tests, including some that spin up
// a test HTTP server and actually check if k6 hits it

// TODO: also add a test that starts multiple k6 "instances", for example:
//  - one with `k6 run --paused` and another with `k6 resume`
//  - one with `k6 run` and another with `k6 stats` or `k6 status`
