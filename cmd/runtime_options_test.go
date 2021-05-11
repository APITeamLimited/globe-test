/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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
	"fmt"
	"net/url"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/loader"
)

type runtimeOptionsTestCase struct ***REMOVED***
	useSysEnv bool // Whether to include the system env vars by default (run) or not (cloud/archive/inspect)
	expErr    bool
	cliFlags  []string
	systemEnv map[string]string
	expRTOpts lib.RuntimeOptions
***REMOVED***

//nolint:gochecknoglobals
var (
	defaultCompatMode  = null.NewString("extended", false)
	baseCompatMode     = null.NewString("base", true)
	extendedCompatMode = null.NewString("extended", true)
)

var runtimeOptionsTestCases = map[string]runtimeOptionsTestCase***REMOVED*** //nolint:gochecknoglobals
	"empty env": ***REMOVED***
		useSysEnv: true,
		// everything else is empty
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  nil,
		***REMOVED***,
	***REMOVED***,
	"disabled sys env by default": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"disabled sys env by default with ext compat mode": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"test1": "val1", "K6_COMPATIBILITY_MODE": "extended"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, false),
			CompatibilityMode:    extendedCompatMode,
			Env:                  map[string]string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"disabled sys env by cli 1": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"test1": "val1", "K6_COMPATIBILITY_MODE": "base"***REMOVED***,
		cliFlags:  []string***REMOVED***"--include-system-env-vars=false"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, true),
			CompatibilityMode:    baseCompatMode,
			Env:                  map[string]string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"disabled sys env by cli 2": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"K6_INCLUDE_SYSTEM_ENV_VARS": "true", "K6_COMPATIBILITY_MODE": "extended"***REMOVED***,
		cliFlags:  []string***REMOVED***"--include-system-env-vars=0", "--compatibility-mode=base"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, true),
			CompatibilityMode:    baseCompatMode,
			Env:                  map[string]string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"disabled sys env by env": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"K6_INCLUDE_SYSTEM_ENV_VARS": "false", "K6_COMPATIBILITY_MODE": "extended"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, true),
			CompatibilityMode:    extendedCompatMode,
			Env:                  map[string]string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"enabled sys env by env": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"K6_INCLUDE_SYSTEM_ENV_VARS": "true", "K6_COMPATIBILITY_MODE": "extended"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, true),
			CompatibilityMode:    extendedCompatMode,
			Env:                  map[string]string***REMOVED***"K6_INCLUDE_SYSTEM_ENV_VARS": "true", "K6_COMPATIBILITY_MODE": "extended"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"enabled sys env by default": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		cliFlags:  []string***REMOVED******REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"enabled sys env by cli 1": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		cliFlags:  []string***REMOVED***"--include-system-env-vars"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, true),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"enabled sys env by cli 2": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		cliFlags:  []string***REMOVED***"--include-system-env-vars=true"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, true),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"run only system env": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		cliFlags:  []string***REMOVED******REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"mixed system and cli env": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"test1": "val1", "test2": ""***REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "test3=val3", "-e", "test4", "-e", "test5="***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "val1", "test2": "", "test3": "val3", "test4": "", "test5": ""***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"mixed system and cli env 2": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"test1": "val1", "test2": ""***REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "test3=val3", "-e", "test4", "-e", "test5=", "--include-system-env-vars=1"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, true),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "val1", "test2": "", "test3": "val3", "test4": "", "test5": ""***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"disabled system env with cli params": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		cliFlags:  []string***REMOVED***"-e", "test2=overwriten", "-e", "test2=val2"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test2": "val2"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"overwriting system env with cli param": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"test1": "val1sys"***REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "test1=val1cli"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "val1cli"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"error wrong compat mode env var value": ***REMOVED***
		systemEnv: map[string]string***REMOVED***"K6_COMPATIBILITY_MODE": "asdf"***REMOVED***,
		expErr:    true,
	***REMOVED***,
	"error wrong compat mode env var value even with CLI flag": ***REMOVED***
		systemEnv: map[string]string***REMOVED***"K6_COMPATIBILITY_MODE": "asdf"***REMOVED***,
		cliFlags:  []string***REMOVED***"--compatibility-mode", "true"***REMOVED***,
		expErr:    true,
	***REMOVED***,
	"error wrong compat mode cli flag value": ***REMOVED***
		cliFlags: []string***REMOVED***"--compatibility-mode", "whatever"***REMOVED***,
		expErr:   true,
	***REMOVED***,
	"error invalid cli var name 1": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED******REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "test a=error"***REMOVED***,
		expErr:    true,
	***REMOVED***,
	"error invalid cli var name 2": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED******REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "1var=error"***REMOVED***,
		expErr:    true,
	***REMOVED***,
	"error invalid cli var name 3": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED******REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "уникод=unicode-disabled"***REMOVED***,
		expErr:    true,
	***REMOVED***,
	"valid env vars with spaces": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"test1": "value 1"***REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "test2=value 2"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "value 1", "test2": "value 2"***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"valid env vars with special chars": ***REMOVED***
		useSysEnv: true,
		systemEnv: map[string]string***REMOVED***"test1": "value 1"***REMOVED***,
		cliFlags:  []string***REMOVED***"--env", "test2=value,2", "-e", `test3= ,  ,,, value, ,, 2!'@#,"`***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(true, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED***"test1": "value 1", "test2": "value,2", "test3": ` ,  ,,, value, ,, 2!'@#,"`***REMOVED***,
		***REMOVED***,
	***REMOVED***,
	"summary and thresholds from env": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"K6_NO_THRESHOLDS": "false", "K6_NO_SUMMARY": "0", "K6_SUMMARY_EXPORT": "foo"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED******REMOVED***,
			NoThresholds:         null.NewBool(false, true),
			NoSummary:            null.NewBool(false, true),
			SummaryExport:        null.NewString("foo", true),
		***REMOVED***,
	***REMOVED***,
	"summary and thresholds from env overwritten by CLI": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"K6_NO_THRESHOLDS": "FALSE", "K6_NO_SUMMARY": "0", "K6_SUMMARY_EXPORT": "foo"***REMOVED***,
		cliFlags:  []string***REMOVED***"--no-thresholds", "true", "--no-summary", "true", "--summary-export", "bar"***REMOVED***,
		expRTOpts: lib.RuntimeOptions***REMOVED***
			IncludeSystemEnvVars: null.NewBool(false, false),
			CompatibilityMode:    defaultCompatMode,
			Env:                  map[string]string***REMOVED******REMOVED***,
			NoThresholds:         null.NewBool(true, true),
			NoSummary:            null.NewBool(true, true),
			SummaryExport:        null.NewString("bar", true),
		***REMOVED***,
	***REMOVED***,
	"env var error detected even when CLI flags overwrite 1": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"K6_NO_THRESHOLDS": "boo"***REMOVED***,
		cliFlags:  []string***REMOVED***"--no-thresholds", "true"***REMOVED***,
		expErr:    true,
	***REMOVED***,
	"env var error detected even when CLI flags overwrite 2": ***REMOVED***
		useSysEnv: false,
		systemEnv: map[string]string***REMOVED***"K6_NO_SUMMARY": "hoo"***REMOVED***,
		cliFlags:  []string***REMOVED***"--no-summary", "true"***REMOVED***,
		expErr:    true,
	***REMOVED***,
***REMOVED***

func testRuntimeOptionsCase(t *testing.T, tc runtimeOptionsTestCase) ***REMOVED***
	flags := runtimeOptionFlagSet(tc.useSysEnv)
	require.NoError(t, flags.Parse(tc.cliFlags))

	rtOpts, err := getRuntimeOptions(flags, tc.systemEnv)
	if tc.expErr ***REMOVED***
		require.Error(t, err)
		return
	***REMOVED***
	require.NoError(t, err)
	require.Equal(t, tc.expRTOpts, rtOpts)

	compatMode, err := lib.ValidateCompatibilityMode(rtOpts.CompatibilityMode.String)
	require.NoError(t, err)

	jsCode := new(bytes.Buffer)
	if compatMode == lib.CompatibilityModeExtended ***REMOVED***
		fmt.Fprint(jsCode, "export default function() ***REMOVED***")
	***REMOVED*** else ***REMOVED***
		fmt.Fprint(jsCode, "module.exports.default = function() ***REMOVED***")
	***REMOVED***

	for key, val := range tc.expRTOpts.Env ***REMOVED***
		fmt.Fprintf(jsCode,
			"if (__ENV.%s !== `%s`) ***REMOVED*** throw new Error('Invalid %s: ' + __ENV.%s); ***REMOVED***",
			key, val, key, key,
		)
	***REMOVED***
	fmt.Fprint(jsCode, "***REMOVED***")

	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/script.js", jsCode.Bytes(), 0o644))
	runner, err := newRunner(
		testutils.NewLogger(t),
		&loader.SourceData***REMOVED***Data: jsCode.Bytes(), URL: &url.URL***REMOVED***Path: "/script.js", Scheme: "file"***REMOVED******REMOVED***,
		typeJS,
		map[string]afero.Fs***REMOVED***"file": fs***REMOVED***,
		rtOpts,
	)
	require.NoError(t, err)

	archive := runner.MakeArchive()
	archiveBuf := &bytes.Buffer***REMOVED******REMOVED***
	require.NoError(t, archive.Write(archiveBuf))

	getRunnerErr := func(rtOpts lib.RuntimeOptions) (lib.Runner, error) ***REMOVED***
		return newRunner(
			testutils.NewLogger(t),
			&loader.SourceData***REMOVED***
				Data: archiveBuf.Bytes(),
				URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
			***REMOVED***,
			typeArchive,
			nil,
			rtOpts,
		)
	***REMOVED***

	_, err = getRunnerErr(lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)
	for key, val := range tc.expRTOpts.Env ***REMOVED***
		r, err := getRunnerErr(lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***key: "almost " + val***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, r.MakeArchive().Env[key], "almost "+val)
	***REMOVED***
***REMOVED***

func TestRuntimeOptions(t *testing.T) ***REMOVED***
	for name, tc := range runtimeOptionsTestCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("RuntimeOptions test '%s'", name), func(t *testing.T) ***REMOVED***
			t.Parallel()
			testRuntimeOptionsCase(t, tc)
		***REMOVED***)
	***REMOVED***
***REMOVED***
