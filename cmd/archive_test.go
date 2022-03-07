package cmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/errext/exitcodes"
)

func TestArchiveThresholds(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		name         string
		noThresholds bool
		testFilename string

		wantErr bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:         "archive should fail with exit status 104 on a malformed threshold expression",
			noThresholds: false,
			testFilename: "testdata/thresholds/malformed_expression.js",
			wantErr:      true,
		***REMOVED***,
		***REMOVED***
			name:         "archive should on a malformed threshold expression but --no-thresholds flag set",
			noThresholds: true,
			testFilename: "testdata/thresholds/malformed_expression.js",
			wantErr:      false,
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			testScript, err := ioutil.ReadFile(testCase.testFilename)
			require.NoError(t, err)

			testState := newGlobalTestState(t)
			require.NoError(t, afero.WriteFile(testState.fs, filepath.Join(testState.cwd, testCase.testFilename), testScript, 0o644))
			testState.args = []string***REMOVED***"k6", "archive", testCase.testFilename***REMOVED***
			if testCase.noThresholds ***REMOVED***
				testState.args = append(testState.args, "--no-thresholds")
			***REMOVED***

			if testCase.wantErr ***REMOVED***
				testState.expectedExitCode = int(exitcodes.InvalidConfig)
			***REMOVED***
			newRootCommand(testState.globalState).execute()
		***REMOVED***)
	***REMOVED***
***REMOVED***
