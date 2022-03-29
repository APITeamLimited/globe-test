package cmd

// TODO: convert this into the integration tests, once https://github.com/grafana/k6/issues/2459 will be done

import (
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib/testutils"
)

// alarmist is a mock module that do a panic
type alarmist struct ***REMOVED***
	vu modules.VU
***REMOVED***

var _ modules.Module = &alarmist***REMOVED******REMOVED***

func (a *alarmist) NewModuleInstance(vu modules.VU) modules.Instance ***REMOVED***
	return &alarmist***REMOVED***
		vu: vu,
	***REMOVED***
***REMOVED***

func (a *alarmist) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***
		Named: map[string]interface***REMOVED******REMOVED******REMOVED***
			"panic": a.panic,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (a *alarmist) panic(s string) ***REMOVED***
	panic(s)
***REMOVED***

func init() ***REMOVED***
	modules.Register("k6/x/alarmist", new(alarmist))
***REMOVED***

func TestRunScriptPanicsErrorsAndAbort(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		caseName, testScript, expectedLogMessage string
	***REMOVED******REMOVED***
		***REMOVED***
			caseName: "panic in the VU context",
			testScript: `
			import ***REMOVED*** panic ***REMOVED*** from 'k6/x/alarmist';

			export default function() ***REMOVED***
				panic('hey');
				console.log('lorem ipsum');
			***REMOVED***
			`,
			expectedLogMessage: "a panic occurred during JS execution: hey",
		***REMOVED***,
		***REMOVED***
			caseName: "panic in the init context",
			testScript: `
			import ***REMOVED*** panic ***REMOVED*** from 'k6/x/alarmist';

			panic('hey');
			export default function() ***REMOVED***
				console.log('lorem ipsum');
			***REMOVED***
			`,
			expectedLogMessage: "a panic occurred during JS execution: hey",
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		name := tc.caseName

		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			testFilename := "script.js"
			testState := newGlobalTestState(t)
			require.NoError(t, afero.WriteFile(testState.fs, filepath.Join(testState.cwd, testFilename), []byte(tc.testScript), 0o644))
			testState.args = []string***REMOVED***"k6", "run", testFilename***REMOVED***

			testState.expectedExitCode = int(exitcodes.ScriptAborted)
			newRootCommand(testState.globalState).execute()

			logs := testState.loggerHook.Drain()

			assert.True(t, testutils.LogContains(logs, logrus.ErrorLevel, tc.expectedLogMessage))
			assert.False(t, testutils.LogContains(logs, logrus.InfoLevel, "lorem ipsum"))
		***REMOVED***)
	***REMOVED***
***REMOVED***
