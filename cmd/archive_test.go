package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/lib/testutils"
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

			cmd := getArchiveCmd(testutils.NewLogger(t), newCommandFlags())
			filename, err := filepath.Abs(testCase.testFilename)
			require.NoError(t, err)
			args := []string***REMOVED***filename***REMOVED***
			if testCase.noThresholds ***REMOVED***
				args = append(args, "--no-thresholds")
			***REMOVED***
			cmd.SetArgs(args)
			wantExitCode := exitcodes.InvalidConfig

			var gotErrExt errext.HasExitCode
			gotErr := cmd.Execute()

			assert.Equal(t,
				testCase.wantErr,
				gotErr != nil,
				"archive command error = %v, wantErr %v", gotErr, testCase.wantErr,
			)

			if testCase.wantErr ***REMOVED***
				require.ErrorAs(t, gotErr, &gotErrExt)
				assert.Equalf(t, wantExitCode, gotErrExt.ExitCode(),
					"status code must be %d", wantExitCode,
				)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
