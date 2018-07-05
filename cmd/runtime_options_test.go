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
	"os"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EnvVarTest struct ***REMOVED***
	name      string
	useSysEnv bool // Whether to include the system env vars by default (run) or not (cloud/archive/inspect)
	systemEnv map[string]string
	cliOpts   []string
	expErr    bool
	expEnv    map[string]string
***REMOVED***

var envVarTestCases = []EnvVarTest***REMOVED***
	***REMOVED***
		"empty env",
		true,
		map[string]string***REMOVED******REMOVED***,
		[]string***REMOVED******REMOVED***,
		false,
		map[string]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		"disabled sys env by default",
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED******REMOVED***,
		false,
		map[string]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		"disabled sys env by cli 1",
		true,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED***"--include-system-env-vars=false"***REMOVED***,
		false,
		map[string]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		"disabled sys env by cli 2",
		true,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED***"--include-system-env-vars=0"***REMOVED***,
		false,
		map[string]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		"enabled sys env by default",
		true,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED******REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"enabled sys env by cli 1",
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED***"--include-system-env-vars"***REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"enabled sys env by cli 2",
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED***"--include-system-env-vars=true"***REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"run only system env",
		true,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED******REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"mixed system and cli env",
		true,
		map[string]string***REMOVED***"test1": "val1", "test2": ""***REMOVED***,
		[]string***REMOVED***"--env", "test3=val3", "-e", "test4", "-e", "test5="***REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "val1", "test2": "", "test3": "val3", "test4": "", "test5": ""***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"mixed system and cli env 2",
		false,
		map[string]string***REMOVED***"test1": "val1", "test2": ""***REMOVED***,
		[]string***REMOVED***"--env", "test3=val3", "-e", "test4", "-e", "test5=", "--include-system-env-vars=1"***REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "val1", "test2": "", "test3": "val3", "test4": "", "test5": ""***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"disabled system env with cli params",
		false,
		map[string]string***REMOVED***"test1": "val1"***REMOVED***,
		[]string***REMOVED***"-e", "test2=overwriten", "-e", "test2=val2"***REMOVED***,
		false,
		map[string]string***REMOVED***"test2": "val2"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"overwriting system env with cli param",
		true,
		map[string]string***REMOVED***"test1": "val1sys"***REMOVED***,
		[]string***REMOVED***"--env", "test1=val1cli"***REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "val1cli"***REMOVED***,
	***REMOVED***,
	***REMOVED***
		"error invalid cli var name 1",
		true,
		map[string]string***REMOVED******REMOVED***,
		[]string***REMOVED***"--env", "test a=error"***REMOVED***,
		true,
		map[string]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		"error invalid cli var name 2",
		true,
		map[string]string***REMOVED******REMOVED***,
		[]string***REMOVED***"--env", "1var=error"***REMOVED***,
		true,
		map[string]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		"error invalid cli var name 3",
		true,
		map[string]string***REMOVED******REMOVED***,
		[]string***REMOVED***"--env", "уникод=unicode-disabled"***REMOVED***,
		true,
		map[string]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		"valid env vars with spaces",
		true,
		map[string]string***REMOVED***"test1": "value 1"***REMOVED***,
		[]string***REMOVED***"--env", "test2=value 2"***REMOVED***,
		false,
		map[string]string***REMOVED***"test1": "value 1", "test2": "value 2"***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestEnvVars(t *testing.T) ***REMOVED***
	for _, tc := range envVarTestCases ***REMOVED***
		t.Run(fmt.Sprintf("EnvVar test '%s'", tc.name), func(t *testing.T) ***REMOVED***
			os.Clearenv()
			for key, val := range tc.systemEnv ***REMOVED***
				require.NoError(t, os.Setenv(key, val))
			***REMOVED***
			flags := runtimeOptionFlagSet(tc.useSysEnv)
			require.NoError(t, flags.Parse(tc.cliOpts))

			rtOpts, err := getRuntimeOptions(flags)
			if tc.expErr ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			require.EqualValues(t, tc.expEnv, rtOpts.Env)

			// Clear the env again so real system values don't accidentally pollute the end-to-end test
			os.Clearenv()

			jsCode := "export default function() ***REMOVED***\n"
			for key, val := range tc.expEnv ***REMOVED***
				jsCode += fmt.Sprintf(
					"if (__ENV.%s !== '%s') ***REMOVED*** throw new Error('Invalid %s: ' + __ENV.%s); ***REMOVED***\n",
					key, val, key, key,
				)
			***REMOVED***
			jsCode += "***REMOVED***"

			runner, err := newRunner(
				&lib.SourceData***REMOVED***
					Data:     []byte(jsCode),
					Filename: "/script.js",
				***REMOVED***,
				typeJS,
				afero.NewOsFs(),
				rtOpts,
			)
			require.NoError(t, err)

			archive := runner.MakeArchive()
			archiveBuf := &bytes.Buffer***REMOVED******REMOVED***
			assert.NoError(t, archive.Write(archiveBuf))

			getRunnerErr := func(rtOpts lib.RuntimeOptions) (lib.Runner, error) ***REMOVED***
				r, err := newRunner(
					&lib.SourceData***REMOVED***
						Data:     []byte(archiveBuf.Bytes()),
						Filename: "/script.tar",
					***REMOVED***,
					typeArchive,
					afero.NewOsFs(),
					rtOpts,
				)
				return r, err
			***REMOVED***

			_, err = getRunnerErr(lib.RuntimeOptions***REMOVED******REMOVED***)
			require.NoError(t, err)
			for key, val := range tc.expEnv ***REMOVED***
				r, err := getRunnerErr(lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***key: "almost " + val***REMOVED******REMOVED***)
				assert.NoError(t, err)
				assert.Equal(t, r.MakeArchive().Env[key], "almost "+val)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
