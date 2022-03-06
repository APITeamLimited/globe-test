package cmd

import (
	"bytes"
	"context"
	"os/signal"
	"runtime"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/lib/testutils"
)

type globalTestState struct ***REMOVED***
	*globalState
	cancel func()

	stdOut, stdErr *bytes.Buffer
	loggerHook     *testutils.SimpleLogrusHook

	cwd string
***REMOVED***

func newGlobalTestState(t *testing.T) *globalTestState ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	fs := &afero.MemMapFs***REMOVED******REMOVED***
	cwd := "/test/"
	if runtime.GOOS == "windows" ***REMOVED***
		cwd = "c:\\test\\"
	***REMOVED***
	require.NoError(t, fs.MkdirAll(cwd, 0o755))

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.Out = testutils.NewTestOutput(t)
	hook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: logrus.AllLevels***REMOVED***
	logger.AddHook(hook)

	ts := &globalTestState***REMOVED***
		cwd:        cwd,
		cancel:     cancel,
		loggerHook: hook,
		stdOut:     new(bytes.Buffer),
		stdErr:     new(bytes.Buffer),
	***REMOVED***

	outMutex := &sync.Mutex***REMOVED******REMOVED***
	defaultFlags := getDefaultFlags(".config")
	ts.globalState = &globalState***REMOVED***
		ctx:            ctx,
		fs:             fs,
		getwd:          func() (string, error) ***REMOVED*** return ts.cwd, nil ***REMOVED***,
		args:           []string***REMOVED******REMOVED***,
		envVars:        map[string]string***REMOVED******REMOVED***,
		defaultFlags:   defaultFlags,
		flags:          defaultFlags,
		outMutex:       outMutex,
		stdOut:         &consoleWriter***REMOVED***nil, ts.stdOut, false, outMutex, nil***REMOVED***,
		stdErr:         &consoleWriter***REMOVED***nil, ts.stdErr, false, outMutex, nil***REMOVED***,
		stdIn:          new(bytes.Buffer),
		signalNotify:   signal.Notify,
		signalStop:     signal.Stop,
		logger:         logger,
		fallbackLogger: testutils.NewLogger(t).WithField("fallback", true),
	***REMOVED***
	return ts
***REMOVED***

func TestDeprecatedOptionWarning(t *testing.T) ***REMOVED***
	t.Parallel()

	ts := newGlobalTestState(t)
	ts.args = []string***REMOVED***"k6", "--logformat", "json", "run", "-"***REMOVED***
	ts.stdIn = bytes.NewBuffer([]byte(`
		console.log('foo');
		export default function() ***REMOVED*** console.log('bar'); ***REMOVED***;
	`))

	root := newRootCommand(ts.globalState)

	require.NoError(t, root.cmd.Execute())

	logMsgs := ts.loggerHook.Drain()
	assert.True(t, testutils.LogContains(logMsgs, logrus.InfoLevel, "foo"))
	assert.True(t, testutils.LogContains(logMsgs, logrus.InfoLevel, "bar"))
	assert.Contains(t, ts.stdErr.String(), `"level":"info","msg":"foo","source":"console"`)
	assert.Contains(t, ts.stdErr.String(), `"level":"info","msg":"bar","source":"console"`)

	// TODO: after we get rid of cobra, actually emit this message to stderr
	// and, ideally, through the log, not just print it...
	assert.False(t, testutils.LogContains(logMsgs, logrus.InfoLevel, "logformat"))
	assert.Contains(t, ts.stdOut.String(), `--logformat has been deprecated`)
***REMOVED***
