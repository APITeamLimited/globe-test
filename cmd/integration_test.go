package cmd

import (
	"bytes"
	"path/filepath"
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

// TODO: add a hell of a lot more integration tests, including some that spin up
// a test HTTP server and actually check if k6 hits it

// TODO: also add a test that starts multiple k6 "instances", for example:
//  - one with `k6 run --paused` and another with `k6 resume`
//  - one with `k6 run` and another with `k6 stats` or `k6 status`
