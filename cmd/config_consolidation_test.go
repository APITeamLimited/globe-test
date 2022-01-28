/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/minirunner"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
)

func verifyOneIterPerOneVU(t *testing.T, c Config) ***REMOVED***
	// No config anywhere should result in a 1 VU with a 1 iteration config
	exec := c.Scenarios[lib.DefaultScenarioName]
	require.NotEmpty(t, exec)
	require.IsType(t, executor.PerVUIterationsConfig***REMOVED******REMOVED***, exec)
	perVuIters, ok := exec.(executor.PerVUIterationsConfig)
	require.True(t, ok)
	assert.Equal(t, null.NewInt(1, false), perVuIters.Iterations)
	assert.Equal(t, null.NewInt(1, false), perVuIters.VUs)
***REMOVED***

func verifySharedIters(vus, iters null.Int) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		exec := c.Scenarios[lib.DefaultScenarioName]
		require.NotEmpty(t, exec)
		require.IsType(t, executor.SharedIterationsConfig***REMOVED******REMOVED***, exec)
		sharedIterConfig, ok := exec.(executor.SharedIterationsConfig)
		require.True(t, ok)
		assert.Equal(t, vus, sharedIterConfig.VUs)
		assert.Equal(t, iters, sharedIterConfig.Iterations)
		assert.Equal(t, vus, c.VUs)
		assert.Equal(t, iters, c.Iterations)
	***REMOVED***
***REMOVED***

func verifyConstLoopingVUs(vus null.Int, duration time.Duration) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		exec := c.Scenarios[lib.DefaultScenarioName]
		require.NotEmpty(t, exec)
		require.IsType(t, executor.ConstantVUsConfig***REMOVED******REMOVED***, exec)
		clvc, ok := exec.(executor.ConstantVUsConfig)
		require.True(t, ok)
		assert.Equal(t, vus, clvc.VUs)
		assert.Equal(t, types.NullDurationFrom(duration), clvc.Duration)
		assert.Equal(t, vus, c.VUs)
		assert.Equal(t, types.NullDurationFrom(duration), c.Duration)
	***REMOVED***
***REMOVED***

func verifyExternallyExecuted(scenarioName string, vus null.Int, duration time.Duration) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		exec := c.Scenarios[scenarioName]
		require.NotEmpty(t, exec)
		require.IsType(t, executor.ExternallyControlledConfig***REMOVED******REMOVED***, exec)
		ecc, ok := exec.(executor.ExternallyControlledConfig)
		require.True(t, ok)
		assert.Equal(t, vus, ecc.VUs)
		assert.Equal(t, types.NullDurationFrom(duration), ecc.Duration)
		assert.Equal(t, vus, ecc.MaxVUs) // MaxVUs defaults to VUs unless specified
	***REMOVED***
***REMOVED***

func verifyRampingVUs(startVus null.Int, stages []executor.Stage) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		exec := c.Scenarios[lib.DefaultScenarioName]
		require.NotEmpty(t, exec)
		require.IsType(t, executor.RampingVUsConfig***REMOVED******REMOVED***, exec)
		clvc, ok := exec.(executor.RampingVUsConfig)
		require.True(t, ok)
		assert.Equal(t, startVus, clvc.StartVUs)
		assert.Equal(t, startVus, c.VUs)
		assert.Equal(t, stages, clvc.Stages)
		assert.Len(t, c.Stages, len(stages))
		for i, s := range stages ***REMOVED***
			assert.Equal(t, s.Duration, c.Stages[i].Duration)
			assert.Equal(t, s.Target, c.Stages[i].Target)
		***REMOVED***
	***REMOVED***
***REMOVED***

// A helper function that accepts (duration in second, VUs) pairs and returns
// a valid slice of stage structs
func buildStages(durationsAndVUs ...int64) []executor.Stage ***REMOVED***
	l := len(durationsAndVUs)
	if l%2 != 0 ***REMOVED***
		panic("wrong len")
	***REMOVED***
	result := make([]executor.Stage, 0, l/2)
	for i := 0; i < l; i += 2 ***REMOVED***
		result = append(result, executor.Stage***REMOVED***
			Duration: types.NullDurationFrom(time.Duration(durationsAndVUs[i]) * time.Second),
			Target:   null.IntFrom(durationsAndVUs[i+1]),
		***REMOVED***)
	***REMOVED***
	return result
***REMOVED***

func mostFlagSets() []flagSetInit ***REMOVED***
	// TODO: make this unnecessary... currently these are the only commands in which
	// getConsolidatedConfig() is used, but they also have differences in their CLI flags :/
	// sigh... compromises...
	result := []flagSetInit***REMOVED******REMOVED***
	for i, fsi := range []flagSetInit***REMOVED***runCmdFlagSet, archiveCmdFlagSet, cloudCmdFlagSet***REMOVED*** ***REMOVED***
		i, fsi := i, fsi // go...
		result = append(result, func() *pflag.FlagSet ***REMOVED***
			flags := pflag.NewFlagSet(fmt.Sprintf("superContrivedFlags_%d", i), pflag.ContinueOnError)
			flags.AddFlagSet(new(rootCommand).rootCmdPersistentFlagSet())
			flags.AddFlagSet(fsi())
			return flags
		***REMOVED***)
	***REMOVED***
	return result
***REMOVED***

type file struct ***REMOVED***
	filepath, contents string
***REMOVED***

func getFS(files []file) afero.Fs ***REMOVED***
	fs := afero.NewMemMapFs()
	for _, f := range files ***REMOVED***
		must(afero.WriteFile(fs, f.filepath, []byte(f.contents), 0o644)) // modes don't matter in the afero.MemMapFs
	***REMOVED***
	return fs
***REMOVED***

func defaultConfig(jsonConfig string) afero.Fs ***REMOVED***
	return getFS([]file***REMOVED******REMOVED***defaultConfigFilePath, jsonConfig***REMOVED******REMOVED***)
***REMOVED***

type flagSetInit func() *pflag.FlagSet

type opts struct ***REMOVED***
	cli    []string
	env    []string
	runner *lib.Options
	fs     afero.Fs

	// TODO: remove this when the configuration is more reproducible and sane...
	// We use a func, because initializing a FlagSet that points to variables
	// actually will change those variables to their default values :| In our
	// case, this happens only some of the time, for global variables that
	// are configurable only via CLI flags, but not environment variables.
	//
	// For the rest, their default value is their current value, since that
	// has been set from the environment variable. That has a bunch of other
	// issues on its own, and the func() doesn't help at all, and we need to
	// use the resetStickyGlobalVars() hack on top of that...
	cliFlagSetInits []flagSetInit
***REMOVED***

func resetStickyGlobalVars() ***REMOVED***
	// TODO: remove after fixing the config, obviously a dirty hack
	exitOnRunning = false
	configFilePath = ""
	runType = ""
***REMOVED***

// exp contains the different events or errors we expect our test case to trigger.
// for space and clarity, we use the fact that by default, all of the struct values are false
type exp struct ***REMOVED***
	cliParseError      bool
	cliReadError       bool
	consolidationError bool // Note: consolidationError includes validation errors from envconfig.Process()
	derivationError    bool
	validationErrors   bool
	logWarning         bool
***REMOVED***

// A hell of a complicated test case, that still doesn't test things fully...
type configConsolidationTestCase struct ***REMOVED***
	options         opts
	expected        exp
	customValidator func(t *testing.T, c Config)
***REMOVED***

func getConfigConsolidationTestCases() []configConsolidationTestCase ***REMOVED***
	I := null.IntFrom // shortcut for "Valid" (i.e. user-specified) ints
	// This is a function, because some of these test cases actually need for the init() functions
	// to be executed, since they depend on defaultConfigFilePath
	return []configConsolidationTestCase***REMOVED***
		// Check that no options will result in 1 VU 1 iter value for execution
		***REMOVED***opts***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyOneIterPerOneVU***REMOVED***,
		// Verify some CLI errors
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--blah", "blah"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--duration", "blah"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--duration", "1000"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***, // intentionally unsupported
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--iterations", "blah"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--execution", ""***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--stage", "10:20s"***REMOVED******REMOVED***, exp***REMOVED***cliReadError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--stage", "1000:20"***REMOVED******REMOVED***, exp***REMOVED***cliReadError: true***REMOVED***, nil***REMOVED***, // intentionally unsupported
		// Check if CLI shortcuts generate correct execution values
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "1", "--iterations", "5"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(1), I(5))***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "2", "-i", "6"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(2), I(6))***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-d", "123s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(null.NewInt(1, false), 123*time.Second)***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-d", "30s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(3), 30*time.Second)***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "4", "--duration", "60s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(4), 1*time.Minute)***REMOVED***,
		***REMOVED***
			opts***REMOVED***cli: []string***REMOVED***"--stage", "20s:10", "-s", "3m:5"***REMOVED******REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyRampingVUs(null.NewInt(1, false), buildStages(20, 10, 180, 5)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***cli: []string***REMOVED***"-s", "1m6s:5", "--vus", "10"***REMOVED******REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyRampingVUs(null.NewInt(10, true), buildStages(66, 5)),
		***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "1", "-i", "6", "-d", "10s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			verifySharedIters(I(1), I(6))(t, c)
			sharedIterConfig, ok := c.Scenarios[lib.DefaultScenarioName].(executor.SharedIterationsConfig)
			require.True(t, ok)
			assert.Equal(t, sharedIterConfig.MaxDuration.TimeDuration(), 10*time.Second)
		***REMOVED******REMOVED***,
		// This should get a validation error since VUs are more than the shared iterations
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "10", "-i", "6"***REMOVED******REMOVED***, exp***REMOVED***validationErrors: true***REMOVED***, verifySharedIters(I(10), I(6))***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-s", "10s:5", "-s", "10s:"***REMOVED******REMOVED***, exp***REMOVED***validationErrors: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***fs: defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s"***REMOVED***], "vus": 10***REMOVED***`)***REMOVED***, exp***REMOVED***validationErrors: true***REMOVED***, nil***REMOVED***,
		// These should emit a derivation error
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "2", "-d", "10s", "-s", "10s:20"***REMOVED******REMOVED***, exp***REMOVED***derivationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-i", "5", "-s", "10s:20"***REMOVED******REMOVED***, exp***REMOVED***derivationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-d", "0"***REMOVED******REMOVED***, exp***REMOVED***derivationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***
			opts***REMOVED***runner: &lib.Options***REMOVED***
				VUs:      null.IntFrom(5),
				Duration: types.NullDurationFrom(44 * time.Second),
				Stages: []lib.Stage***REMOVED***
					***REMOVED***Duration: types.NullDurationFrom(3 * time.Second), Target: I(20)***REMOVED***,
				***REMOVED***,
			***REMOVED******REMOVED***, exp***REMOVED***derivationError: true***REMOVED***, nil,
		***REMOVED***,
		***REMOVED***opts***REMOVED***fs: defaultConfig(`***REMOVED***"scenarios": ***REMOVED******REMOVED******REMOVED***`)***REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, verifyOneIterPerOneVU***REMOVED***,
		// Test if environment variable shortcuts are working as expected
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_VUS=5", "K6_ITERATIONS=15"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(5), I(15))***REMOVED***,
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_VUS=10", "K6_DURATION=20s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(10), 20*time.Second)***REMOVED***,
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_VUS=10", "K6_DURATION=10000"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(10), 10*time.Second)***REMOVED***,
		***REMOVED***
			opts***REMOVED***env: []string***REMOVED***"K6_STAGES=2m30s:11,1h1m:100"***REMOVED******REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyRampingVUs(null.NewInt(1, false), buildStages(150, 11, 3660, 100)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***env: []string***REMOVED***"K6_STAGES=100s:100,0m30s:0", "K6_VUS=0"***REMOVED******REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyRampingVUs(null.NewInt(0, true), buildStages(100, 100, 30, 0)),
		***REMOVED***,
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_STAGES=1000:100"***REMOVED******REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***, // intentionally unsupported
		// Test if JSON configs work as expected
		***REMOVED***opts***REMOVED***fs: defaultConfig(`***REMOVED***"iterations": 77, "vus": 7***REMOVED***`)***REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(7), I(77))***REMOVED***,
		***REMOVED***opts***REMOVED***fs: defaultConfig(`wrong-json`)***REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***fs: getFS(nil), cli: []string***REMOVED***"--config", "/my/config.file"***REMOVED******REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,

		// Test combinations between options and levels
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "1"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyOneIterPerOneVU***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "10"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, verifyOneIterPerOneVU***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:  getFS([]file***REMOVED******REMOVED***"/my/config.file", `***REMOVED***"vus": 8, "duration": "2m"***REMOVED***`***REMOVED******REMOVED***),
				cli: []string***REMOVED***"--config", "/my/config.file"***REMOVED***,
			***REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(8), 120*time.Second),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:  getFS([]file***REMOVED******REMOVED***"/my/config.file", `***REMOVED***"duration": 20000***REMOVED***`***REMOVED******REMOVED***),
				cli: []string***REMOVED***"--config", "/my/config.file"***REMOVED***,
			***REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(null.NewInt(1, false), 20*time.Second),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:  defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s", "target": 20***REMOVED***], "vus": 10***REMOVED***`),
				env: []string***REMOVED***"K6_DURATION=15s"***REMOVED***,
				cli: []string***REMOVED***"--stage", ""***REMOVED***,
			***REMOVED***,
			exp***REMOVED***logWarning: true***REMOVED***,
			verifyOneIterPerOneVU,
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(5), Duration: types.NullDurationFrom(50 * time.Second)***REMOVED***,
				cli:    []string***REMOVED***"--stage", "5s:5"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyRampingVUs(I(5), buildStages(5, 5)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:     defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s", "target": 10***REMOVED***]***REMOVED***`),
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(5)***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyRampingVUs(I(5), buildStages(20, 10)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:     defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s", "target": 10***REMOVED***]***REMOVED***`),
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(5)***REMOVED***,
				env:    []string***REMOVED***"K6_VUS=15", "K6_ITERATIONS=17"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifySharedIters(I(15), I(17)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:     defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "11s", "target": 11***REMOVED***]***REMOVED***`),
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(22)***REMOVED***,
				env:    []string***REMOVED***"K6_VUS=33"***REMOVED***,
				cli:    []string***REMOVED***"--stage", "44s:44", "-s", "55s:55"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyRampingVUs(null.NewInt(33, true), buildStages(44, 44, 55, 55)),
		***REMOVED***,

		// TODO: test the future full overwriting of the duration/iterations/stages/execution options
		***REMOVED***
			opts***REMOVED***
				fs: defaultConfig(`***REMOVED***
					"scenarios": ***REMOVED*** "someKey": ***REMOVED***
						"executor": "constant-vus", "vus": 10, "duration": "60s", "gracefulStop": "10s",
						"startTime": "70s", "env": ***REMOVED***"test": "mest"***REMOVED***, "exec": "someFunc"
					***REMOVED******REMOVED******REMOVED***`),
				env: []string***REMOVED***"K6_ITERATIONS=25"***REMOVED***,
				cli: []string***REMOVED***"--vus", "12"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifySharedIters(I(12), I(25)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs: defaultConfig(`***REMOVED***"scenarios": ***REMOVED*** "foo": ***REMOVED***
					"executor": "constant-vus", "vus": 2, "duration": "1d",
					"gracefulStop": "10000", "startTime": 1000.5
				***REMOVED******REMOVED******REMOVED***`),
			***REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
				exec := c.Scenarios["foo"]
				require.NotEmpty(t, exec)
				require.IsType(t, executor.ConstantVUsConfig***REMOVED******REMOVED***, exec)
				clvc, ok := exec.(executor.ConstantVUsConfig)
				require.True(t, ok)
				assert.Equal(t, null.IntFrom(2), clvc.VUs)
				assert.Equal(t, types.NullDurationFrom(24*time.Hour), clvc.Duration)
				assert.Equal(t, types.NullDurationFrom(time.Second+500*time.Microsecond), clvc.StartTime)
				assert.Equal(t, types.NullDurationFrom(10*time.Second), clvc.GracefulStop)
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs: defaultConfig(`***REMOVED***"scenarios": ***REMOVED*** "def": ***REMOVED***
					"executor": "externally-controlled", "vus": 15, "duration": "2h"
				***REMOVED******REMOVED******REMOVED***`),
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyExternallyExecuted("def", I(15), 2*time.Hour),
		***REMOVED***,
		// TODO: test execution-segment

		// Just in case, verify that no options will result in the same 1 vu 1 iter config
		***REMOVED***opts***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyOneIterPerOneVU***REMOVED***,

		// Test system tags
		***REMOVED***opts***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, &stats.DefaultSystemTagSet, c.Options.SystemTags)
		***REMOVED******REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--system-tags", `""`***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, stats.SystemTagSet(0), *c.Options.SystemTags)
		***REMOVED******REMOVED***,
		***REMOVED***
			opts***REMOVED***
				runner: &lib.Options***REMOVED***
					SystemTags: stats.NewSystemTagSet(stats.TagSubproto, stats.TagURL),
				***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			func(t *testing.T, c Config) ***REMOVED***
				assert.Equal(
					t,
					*stats.NewSystemTagSet(stats.TagSubproto, stats.TagURL),
					*c.Options.SystemTags,
				)
			***REMOVED***,
		***REMOVED***,
		// Test summary trend stats
		***REMOVED***opts***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, lib.DefaultSummaryTrendStats, c.Options.SummaryTrendStats)
		***REMOVED******REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--summary-trend-stats", ""***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, []string***REMOVED******REMOVED***, c.Options.SummaryTrendStats)
		***REMOVED******REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--summary-trend-stats", "coun"***REMOVED******REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--summary-trend-stats", "med,avg,p("***REMOVED******REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--summary-trend-stats", "med,avg,p(-1)"***REMOVED******REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--summary-trend-stats", "med,avg,p(101)"***REMOVED******REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--summary-trend-stats", "med,avg,p(99.999)"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, []string***REMOVED***"med", "avg", "p(99.999)"***REMOVED***, c.Options.SummaryTrendStats)
		***REMOVED******REMOVED***,
		***REMOVED***
			opts***REMOVED***runner: &lib.Options***REMOVED***SummaryTrendStats: []string***REMOVED***"avg", "p(90)", "count"***REMOVED******REMOVED******REMOVED***,
			exp***REMOVED******REMOVED***,
			func(t *testing.T, c Config) ***REMOVED***
				assert.Equal(t, []string***REMOVED***"avg", "p(90)", "count"***REMOVED***, c.Options.SummaryTrendStats)
			***REMOVED***,
		***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED******REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, types.DNSConfig***REMOVED***
				TTL:    null.NewString("5m", false),
				Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSrandom, Valid: false***REMOVED***,
				Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv4, Valid: false***REMOVED***,
			***REMOVED***, c.Options.DNS)
		***REMOVED******REMOVED***,
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_DNS=ttl=5,select=roundRobin"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, types.DNSConfig***REMOVED***
				TTL:    null.StringFrom("5"),
				Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSroundRobin, Valid: true***REMOVED***,
				Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv4, Valid: false***REMOVED***,
			***REMOVED***, c.Options.DNS)
		***REMOVED******REMOVED***,
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_DNS=ttl=inf,select=random,policy=preferIPv6"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, types.DNSConfig***REMOVED***
				TTL:    null.StringFrom("inf"),
				Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSrandom, Valid: true***REMOVED***,
				Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv6, Valid: true***REMOVED***,
			***REMOVED***, c.Options.DNS)
		***REMOVED******REMOVED***,
		// This is functionally invalid, but will error out in validation done in js.parseTTL().
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--dns", "ttl=-1"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, types.DNSConfig***REMOVED***
				TTL:    null.StringFrom("-1"),
				Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSrandom, Valid: false***REMOVED***,
				Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv4, Valid: false***REMOVED***,
			***REMOVED***, c.Options.DNS)
		***REMOVED******REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--dns", "ttl=0,blah=nope"***REMOVED******REMOVED***, exp***REMOVED***cliReadError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--dns", "ttl=0"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, types.DNSConfig***REMOVED***
				TTL:    null.StringFrom("0"),
				Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSrandom, Valid: false***REMOVED***,
				Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSpreferIPv4, Valid: false***REMOVED***,
			***REMOVED***, c.Options.DNS)
		***REMOVED******REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--dns", "ttl=5s,select="***REMOVED******REMOVED***, exp***REMOVED***cliReadError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***
			opts***REMOVED***fs: defaultConfig(`***REMOVED***"dns": ***REMOVED***"ttl": "0", "select": "roundRobin", "policy": "onlyIPv4"***REMOVED******REMOVED***`)***REMOVED***,
			exp***REMOVED******REMOVED***,
			func(t *testing.T, c Config) ***REMOVED***
				assert.Equal(t, types.DNSConfig***REMOVED***
					TTL:    null.StringFrom("0"),
					Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSroundRobin, Valid: true***REMOVED***,
					Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSonlyIPv4, Valid: true***REMOVED***,
				***REMOVED***, c.Options.DNS)
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:  defaultConfig(`***REMOVED***"dns": ***REMOVED***"ttl": "0"***REMOVED******REMOVED***`),
				env: []string***REMOVED***"K6_DNS=ttl=30,policy=any"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			func(t *testing.T, c Config) ***REMOVED***
				assert.Equal(t, types.DNSConfig***REMOVED***
					TTL:    null.StringFrom("30"),
					Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSrandom, Valid: false***REMOVED***,
					Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSany, Valid: true***REMOVED***,
				***REMOVED***, c.Options.DNS)
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			// CLI overrides all, falling back to env
			opts***REMOVED***
				fs:  defaultConfig(`***REMOVED***"dns": ***REMOVED***"ttl": "60", "select": "first"***REMOVED******REMOVED***`),
				env: []string***REMOVED***"K6_DNS=ttl=30,select=random,policy=any"***REMOVED***,
				cli: []string***REMOVED***"--dns", "ttl=5"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			func(t *testing.T, c Config) ***REMOVED***
				assert.Equal(t, types.DNSConfig***REMOVED***
					TTL:    null.StringFrom("5"),
					Select: types.NullDNSSelect***REMOVED***DNSSelect: types.DNSrandom, Valid: true***REMOVED***,
					Policy: types.NullDNSPolicy***REMOVED***DNSPolicy: types.DNSany, Valid: true***REMOVED***,
				***REMOVED***, c.Options.DNS)
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***env: []string***REMOVED***"K6_NO_SETUP=true", "K6_NO_TEARDOWN=false"***REMOVED******REMOVED***,
			exp***REMOVED******REMOVED***,
			func(t *testing.T, c Config) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), c.Options.NoSetup)
				assert.Equal(t, null.BoolFrom(false), c.Options.NoTeardown)
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***env: []string***REMOVED***"K6_NO_SETUP=false", "K6_NO_TEARDOWN=bool"***REMOVED******REMOVED***,
			exp***REMOVED***
				consolidationError: true,
			***REMOVED***,
			nil,
		***REMOVED***,
		// TODO: test for differences between flagsets
		// TODO: more tests in general, especially ones not related to execution parameters...
	***REMOVED***
***REMOVED***

func runTestCase(
	t *testing.T,
	testCase configConsolidationTestCase,
	newFlagSet flagSetInit,
) ***REMOVED***
	t.Helper()
	t.Logf("Test with opts=%#v and exp=%#v\n", testCase.options, testCase.expected)
	output := testutils.NewTestOutput(t)
	logHook := &testutils.SimpleLogrusHook***REMOVED***
		HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED***,
	***REMOVED***

	logHook.Drain()
	logger := logrus.New()
	logger.AddHook(logHook)
	logger.SetOutput(output)

	flagSet := newFlagSet()
	defer resetStickyGlobalVars()
	flagSet.SetOutput(output)
	// flagSet.PrintDefaults()

	cliErr := flagSet.Parse(testCase.options.cli)
	if testCase.expected.cliParseError ***REMOVED***
		require.Error(t, cliErr)
		return
	***REMOVED***
	require.NoError(t, cliErr)

	// TODO: remove these hacks when we improve the configuration...
	var cliConf Config
	if flagSet.Lookup("out") != nil ***REMOVED***
		cliConf, cliErr = getConfig(flagSet)
	***REMOVED*** else ***REMOVED***
		opts, errOpts := getOptions(flagSet)
		cliConf, cliErr = Config***REMOVED***Options: opts***REMOVED***, errOpts
	***REMOVED***
	if testCase.expected.cliReadError ***REMOVED***
		require.Error(t, cliErr)
		return
	***REMOVED***
	require.NoError(t, cliErr)

	var runnerOpts lib.Options
	if testCase.options.runner != nil ***REMOVED***
		runnerOpts = minirunner.MiniRunner***REMOVED***Options: *testCase.options.runner***REMOVED***.GetOptions()
	***REMOVED***
	// without runner creation, values in runnerOpts will simply be invalid

	if testCase.options.fs == nil ***REMOVED***
		t.Logf("Creating an empty FS for this test")
		testCase.options.fs = afero.NewMemMapFs() // create an empty FS if it wasn't supplied
	***REMOVED***

	consolidatedConfig, err := getConsolidatedConfig(testCase.options.fs, cliConf, runnerOpts,
		// TODO: just make testcase.options.env in map[string]string
		buildEnvMap(testCase.options.env))
	if testCase.expected.consolidationError ***REMOVED***
		require.Error(t, err)
		return
	***REMOVED***
	require.NoError(t, err)

	derivedConfig := consolidatedConfig
	derivedConfig.Options, err = executor.DeriveScenariosFromShortcuts(consolidatedConfig.Options, logger)
	if testCase.expected.derivationError ***REMOVED***
		require.Error(t, err)
		return
	***REMOVED***
	require.NoError(t, err)

	if warnings := logHook.Drain(); testCase.expected.logWarning ***REMOVED***
		assert.NotEmpty(t, warnings)
	***REMOVED*** else ***REMOVED***
		assert.Empty(t, warnings)
	***REMOVED***

	validationErrors := derivedConfig.Validate()
	if testCase.expected.validationErrors ***REMOVED***
		assert.NotEmpty(t, validationErrors)
	***REMOVED*** else ***REMOVED***
		assert.Empty(t, validationErrors)
	***REMOVED***

	if testCase.customValidator != nil ***REMOVED***
		testCase.customValidator(t, derivedConfig)
	***REMOVED***
***REMOVED***

//nolint:paralleltest // see comments in test
func TestConfigConsolidation(t *testing.T) ***REMOVED***
	// This test and its subtests shouldn't be ran in parallel, since they unfortunately have
	// to mess with shared global objects (variables, ... santa?)

	for tcNum, testCase := range getConfigConsolidationTestCases() ***REMOVED***
		tcNum, testCase := tcNum, testCase
		flagSetInits := testCase.options.cliFlagSetInits
		if flagSetInits == nil ***REMOVED*** // handle the most common case
			flagSetInits = mostFlagSets()
		***REMOVED***
		for fsNum, flagSet := range flagSetInits ***REMOVED***
			// I want to paralelize this, but I cannot... due to global variables and other
			// questionable architectural choices... :|
			fsNum, flagSet := fsNum, flagSet
			t.Run(
				fmt.Sprintf("TestCase#%d_FlagSet#%d", tcNum, fsNum),
				func(t *testing.T) ***REMOVED*** runTestCase(t, testCase, flagSet) ***REMOVED***,
			)
		***REMOVED***
	***REMOVED***
***REMOVED***
