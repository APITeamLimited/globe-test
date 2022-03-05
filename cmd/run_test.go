/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib/fsext"
)

type mockWriter struct ***REMOVED***
	err      error
	errAfter int
***REMOVED***

func (fw mockWriter) Write(p []byte) (n int, err error) ***REMOVED***
	if fw.err != nil ***REMOVED***
		return fw.errAfter, fw.err
	***REMOVED***
	return len(p), nil
***REMOVED***

var _ io.Writer = mockWriter***REMOVED******REMOVED***

func getFiles(t *testing.T, fs afero.Fs) map[string]*bytes.Buffer ***REMOVED***
	result := map[string]*bytes.Buffer***REMOVED******REMOVED***
	walkFn := func(filePath string, _ os.FileInfo, err error) error ***REMOVED***
		if filePath == "/" || filePath == "\\" ***REMOVED***
			return nil
		***REMOVED***
		require.NoError(t, err)
		contents, err := afero.ReadFile(fs, filePath)
		require.NoError(t, err)
		result[filePath] = bytes.NewBuffer(contents)
		return nil
	***REMOVED***

	err := fsext.Walk(fs, afero.FilePathSeparator, filepath.WalkFunc(walkFn))
	require.NoError(t, err)

	return result
***REMOVED***

func assertEqual(t *testing.T, exp string, actual io.Reader) ***REMOVED***
	act, err := ioutil.ReadAll(actual)
	require.NoError(t, err)
	assert.Equal(t, []byte(exp), act)
***REMOVED***

func initVars() (
	content map[string]io.Reader, stdout *bytes.Buffer, stderr *bytes.Buffer, fs afero.Fs,
) ***REMOVED***
	return map[string]io.Reader***REMOVED******REMOVED***, bytes.NewBuffer([]byte***REMOVED******REMOVED***), bytes.NewBuffer([]byte***REMOVED******REMOVED***), afero.NewMemMapFs()
***REMOVED***

func TestHandleSummaryResultSimple(t *testing.T) ***REMOVED***
	t.Parallel()
	content, stdout, stderr, fs := initVars()

	// Test noop
	assert.NoError(t, handleSummaryResult(fs, stdout, stderr, content))
	require.Empty(t, getFiles(t, fs))
	require.Empty(t, stdout.Bytes())
	require.Empty(t, stderr.Bytes())

	// Test stdout only
	content["stdout"] = bytes.NewBufferString("some stdout summary")
	assert.NoError(t, handleSummaryResult(fs, stdout, stderr, content))
	require.Empty(t, getFiles(t, fs))
	assertEqual(t, "some stdout summary", stdout)
	require.Empty(t, stderr.Bytes())
***REMOVED***

func TestHandleSummaryResultError(t *testing.T) ***REMOVED***
	t.Parallel()
	content, _, stderr, fs := initVars()

	expErr := errors.New("test error")
	stdout := mockWriter***REMOVED***err: expErr, errAfter: 10***REMOVED***

	filePath1 := "/path/file1"
	filePath2 := "/path/file2"
	if runtime.GOOS == "windows" ***REMOVED***
		filePath1 = "\\path\\file1"
		filePath2 = "\\path\\file2"
	***REMOVED***

	content["stdout"] = bytes.NewBufferString("some stdout summary")
	content["stderr"] = bytes.NewBufferString("some stderr summary")
	content[filePath1] = bytes.NewBufferString("file summary 1")
	content[filePath2] = bytes.NewBufferString("file summary 2")
	err := handleSummaryResult(fs, stdout, stderr, content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expErr.Error())
	files := getFiles(t, fs)
	assertEqual(t, "file summary 1", files[filePath1])
	assertEqual(t, "file summary 2", files[filePath2])
***REMOVED***

func TestRunScriptErrorsAndAbort(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct ***REMOVED***
		testFilename, name   string
		expErr, expLogOutput string
		expExitCode          errext.ExitCode
		extraArgs            []string
	***REMOVED******REMOVED***
		***REMOVED***
			testFilename: "abort.js",
			expErr:       common.AbortTest,
			expExitCode:  exitcodes.ScriptAborted,
		***REMOVED***,
		***REMOVED***
			testFilename: "abort_initerr.js",
			expErr:       common.AbortTest,
			expExitCode:  exitcodes.ScriptAborted,
		***REMOVED***,
		***REMOVED***
			testFilename: "abort_initvu.js",
			expErr:       common.AbortTest,
			expExitCode:  exitcodes.ScriptAborted,
		***REMOVED***,
		***REMOVED***
			testFilename: "abort_teardown.js",
			expErr:       common.AbortTest,
			expExitCode:  exitcodes.ScriptAborted,
			expLogOutput: "Calling teardown function after test.abort()",
		***REMOVED***,
		***REMOVED***
			testFilename: "initerr.js",
			expErr:       "ReferenceError: someUndefinedVar is not defined",
			expExitCode:  exitcodes.ScriptException,
		***REMOVED***,
		***REMOVED***
			testFilename: "thresholds/malformed_expression.js",
			name:         "run should fail with exit status 104 on a malformed threshold expression",
			expErr:       "malformed threshold expression",
			expExitCode:  exitcodes.InvalidConfig,
		***REMOVED***,
		***REMOVED***
			testFilename: "thresholds/malformed_expression.js",
			name:         "run should on a malformed threshold expression but --no-thresholds flag set",
			extraArgs:    []string***REMOVED***"--no-thresholds"***REMOVED***,
			// we don't expect an error
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		name := tc.testFilename
		if tc.name != "" ***REMOVED***
			name = fmt.Sprintf("%s (%s)", tc.testFilename, tc.name)
		***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			testScript, err := ioutil.ReadFile(path.Join("testdata", tc.testFilename))
			require.NoError(t, err)

			testState := newGlobalTestState(t)
			require.NoError(t, afero.WriteFile(testState.fs, filepath.Join(testState.cwd, tc.testFilename), testScript, 0o644))
			testState.args = append([]string***REMOVED***"k6", "run", tc.testFilename***REMOVED***, tc.extraArgs...)

			err = newRootCommand(testState.globalState).cmd.Execute()

			if tc.expErr != "" ***REMOVED***
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expErr)
			***REMOVED*** else ***REMOVED***
				require.NoError(t, err)
			***REMOVED***

			if tc.expExitCode != 0 ***REMOVED***
				var e errext.HasExitCode
				require.ErrorAs(t, err, &e)
				assert.Equalf(t, tc.expExitCode, e.ExitCode(), "Status code must be %d", tc.expExitCode)
			***REMOVED***

			if tc.expLogOutput != "" ***REMOVED***
				var gotMsg bool
				for _, entry := range testState.loggerHook.Drain() ***REMOVED***
					if strings.Contains(entry.Message, tc.expLogOutput) ***REMOVED***
						gotMsg = true
						break
					***REMOVED***
				***REMOVED***
				assert.True(t, gotMsg)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
