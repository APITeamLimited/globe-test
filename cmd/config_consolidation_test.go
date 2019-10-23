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
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/scheduler"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
)

// A helper funcion for setting arbitrary environment variables and
// restoring the old ones at the end, usually by deferring the returned callback
//TODO: remove these hacks when we improve the configuration... we shouldn't
// have to mess with the global environment at all...
func setEnv(t *testing.T, newEnv []string) (restoreEnv func()) ***REMOVED***
	actuallSetEnv := func(env []string, abortOnSetErr bool) ***REMOVED***
		os.Clearenv()
		for _, e := range env ***REMOVED***
			val := ""
			pair := strings.SplitN(e, "=", 2)
			if len(pair) > 1 ***REMOVED***
				val = pair[1]
			***REMOVED***
			err := os.Setenv(pair[0], val)
			if abortOnSetErr ***REMOVED***
				require.NoError(t, err)
			***REMOVED*** else if err != nil ***REMOVED***
				t.Logf(
					"Received a non-aborting but unexpected error '%s' when setting env.var '%s' to '%s'",
					err, pair[0], val,
				)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	oldEnv := os.Environ()
	actuallSetEnv(newEnv, true)

	return func() ***REMOVED***
		actuallSetEnv(oldEnv, false)
	***REMOVED***
***REMOVED***

func verifyOneIterPerOneVU(t *testing.T, c Config) ***REMOVED***
	// No config anywhere should result in a 1 VU with a 1 uninterruptible iteration config
	sched := c.Execution[lib.DefaultSchedulerName]
	require.NotEmpty(t, sched)
	require.IsType(t, scheduler.PerVUIteationsConfig***REMOVED******REMOVED***, sched)
	perVuIters, ok := sched.(scheduler.PerVUIteationsConfig)
	require.True(t, ok)
	assert.Equal(t, null.NewInt(1, false), perVuIters.Iterations)
	assert.Equal(t, null.NewInt(1, false), perVuIters.VUs)
***REMOVED***

func verifySharedIters(vus, iters null.Int) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		sched := c.Execution[lib.DefaultSchedulerName]
		require.NotEmpty(t, sched)
		require.IsType(t, scheduler.SharedIteationsConfig***REMOVED******REMOVED***, sched)
		sharedIterConfig, ok := sched.(scheduler.SharedIteationsConfig)
		require.True(t, ok)
		assert.Equal(t, vus, sharedIterConfig.VUs)
		assert.Equal(t, iters, sharedIterConfig.Iterations)
		assert.Equal(t, vus, c.VUs)
		assert.Equal(t, iters, c.Iterations)
	***REMOVED***
***REMOVED***

func verifyConstLoopingVUs(vus null.Int, duration time.Duration) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		sched := c.Execution[lib.DefaultSchedulerName]
		require.NotEmpty(t, sched)
		require.IsType(t, scheduler.ConstantLoopingVUsConfig***REMOVED******REMOVED***, sched)
		clvc, ok := sched.(scheduler.ConstantLoopingVUsConfig)
		require.True(t, ok)
		assert.Equal(t, vus, clvc.VUs)
		assert.Equal(t, types.NullDurationFrom(duration), clvc.Duration)
		assert.Equal(t, vus, c.VUs)
		assert.Equal(t, types.NullDurationFrom(duration), c.Duration)
	***REMOVED***
***REMOVED***

func verifyVarLoopingVUs(startVus null.Int, stages []scheduler.Stage) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		sched := c.Execution[lib.DefaultSchedulerName]
		require.NotEmpty(t, sched)
		require.IsType(t, scheduler.VariableLoopingVUsConfig***REMOVED******REMOVED***, sched)
		clvc, ok := sched.(scheduler.VariableLoopingVUsConfig)
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
func buildStages(durationsAndVUs ...int64) []scheduler.Stage ***REMOVED***
	l := len(durationsAndVUs)
	if l%2 != 0 ***REMOVED***
		panic("wrong len")
	***REMOVED***
	result := make([]scheduler.Stage, 0, l/2)
	for i := 0; i < l; i += 2 ***REMOVED***
		result = append(result, scheduler.Stage***REMOVED***
			Duration: types.NullDurationFrom(time.Duration(durationsAndVUs[i]) * time.Second),
			Target:   null.IntFrom(durationsAndVUs[i+1]),
		***REMOVED***)
	***REMOVED***
	return result
***REMOVED***

func mostFlagSets() []flagSetInit ***REMOVED***
	//TODO: make this unnecessary... currently these are the only commands in which
	// getConsolidatedConfig() is used, but they also have differences in their CLI flags :/
	// sigh... compromises...
	result := []flagSetInit***REMOVED******REMOVED***
	for i, fsi := range []flagSetInit***REMOVED***runCmdFlagSet, archiveCmdFlagSet, cloudCmdFlagSet***REMOVED*** ***REMOVED***
		i, fsi := i, fsi // go...
		result = append(result, func() *pflag.FlagSet ***REMOVED***
			flags := pflag.NewFlagSet(fmt.Sprintf("superContrivedFlags_%d", i), pflag.ContinueOnError)
			flags.AddFlagSet(rootCmdPersistentFlagSet())
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
		must(afero.WriteFile(fs, f.filepath, []byte(f.contents), 0644)) // modes don't matter in the afero.MemMapFs
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

	//TODO: remove this when the configuration is more reproducible and sane...
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
	//TODO: remove after fixing the config, obviously a dirty hack
	exitOnRunning = false
	configFilePath = ""
	runType = ""
	runNoSetup = false
	runNoTeardown = false
***REMOVED***

// Something that makes the test also be a valid io.Writer, useful for passing it
// as an output for logs and CLI flag help messages...
type testOutput struct***REMOVED*** *testing.T ***REMOVED***

func (to testOutput) Write(p []byte) (n int, err error) ***REMOVED***
	to.Logf("%s", p)
	return len(p), nil
***REMOVED***

var _ io.Writer = testOutput***REMOVED******REMOVED***

// exp contains the different events or errors we expect our test case to trigger.
// for space and clarity, we use the fact that by default, all of the struct values are false
type exp struct ***REMOVED***
	cliParseError      bool
	cliReadError       bool
	consolidationError bool
	derivationError    bool
	validationErrors   bool
	logWarning         bool //TODO: remove in the next version?
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
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--iterations", "blah"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--execution", ""***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--stage", "10:20s"***REMOVED******REMOVED***, exp***REMOVED***cliReadError: true***REMOVED***, nil***REMOVED***,
		// Check if CLI shortcuts generate correct execution values
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "1", "--iterations", "5"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(1), I(5))***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "2", "-i", "6"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(2), I(6))***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-d", "123s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(null.NewInt(1, false), 123*time.Second)***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-d", "30s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(3), 30*time.Second)***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "4", "--duration", "60s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(4), 1*time.Minute)***REMOVED***,
		***REMOVED***
			opts***REMOVED***cli: []string***REMOVED***"--stage", "20s:10", "-s", "3m:5"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***,
			verifyVarLoopingVUs(null.NewInt(1, false), buildStages(20, 10, 180, 5)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***cli: []string***REMOVED***"-s", "1m6s:5", "--vus", "10"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***,
			verifyVarLoopingVUs(null.NewInt(10, true), buildStages(66, 5)),
		***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "1", "-i", "6", "-d", "10s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			verifySharedIters(I(1), I(6))(t, c)
			sharedIterConfig := c.Execution[lib.DefaultSchedulerName].(scheduler.SharedIteationsConfig)
			assert.Equal(t, time.Duration(sharedIterConfig.MaxDuration.Duration), 10*time.Second)
		***REMOVED******REMOVED***,
		// This should get a validation error since VUs are more than the shared iterations
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "10", "-i", "6"***REMOVED******REMOVED***, exp***REMOVED***validationErrors: true***REMOVED***, verifySharedIters(I(10), I(6))***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-s", "10s:5", "-s", "10s:"***REMOVED******REMOVED***, exp***REMOVED***validationErrors: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***fs: defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s"***REMOVED***], "vus": 10***REMOVED***`)***REMOVED***, exp***REMOVED***validationErrors: true***REMOVED***, nil***REMOVED***,
		// These should emit a warning
		//TODO: in next version, those should be an error
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "2", "-d", "10s", "-s", "10s:20"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-i", "5", "-s", "10s:20"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-d", "0"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil***REMOVED***,
		***REMOVED***
			opts***REMOVED***runner: &lib.Options***REMOVED***
				VUs:      null.IntFrom(5),
				Duration: types.NullDurationFrom(44 * time.Second),
				Stages: []lib.Stage***REMOVED***
					***REMOVED***Duration: types.NullDurationFrom(3 * time.Second), Target: I(20)***REMOVED***,
				***REMOVED***,
			***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil,
		***REMOVED***,
		***REMOVED***opts***REMOVED***fs: defaultConfig(`***REMOVED***"execution": ***REMOVED******REMOVED******REMOVED***`)***REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, verifyOneIterPerOneVU***REMOVED***,
		// Test if environment variable shortcuts are working as expected
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_VUS=5", "K6_ITERATIONS=15"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(5), I(15))***REMOVED***,
		***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_VUS=10", "K6_DURATION=20s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(10), 20*time.Second)***REMOVED***,
		***REMOVED***
			opts***REMOVED***env: []string***REMOVED***"K6_STAGES=2m30s:11,1h1m:100"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***,
			verifyVarLoopingVUs(null.NewInt(1, false), buildStages(150, 11, 3660, 100)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***env: []string***REMOVED***"K6_STAGES=100s:100,0m30s:0", "K6_VUS=0"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***,
			verifyVarLoopingVUs(null.NewInt(0, true), buildStages(100, 100, 30, 0)),
		***REMOVED***,
		// Test if JSON configs work as expected
		***REMOVED***opts***REMOVED***fs: defaultConfig(`***REMOVED***"iterations": 77, "vus": 7***REMOVED***`)***REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(I(7), I(77))***REMOVED***,
		***REMOVED***opts***REMOVED***fs: defaultConfig(`wrong-json`)***REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,
		***REMOVED***opts***REMOVED***fs: getFS(nil), cli: []string***REMOVED***"--config", "/my/config.file"***REMOVED******REMOVED***, exp***REMOVED***consolidationError: true***REMOVED***, nil***REMOVED***,

		// Test combinations between options and levels
		***REMOVED***
			opts***REMOVED***
				fs:  getFS([]file***REMOVED******REMOVED***"/my/config.file", `***REMOVED***"vus": 8, "duration": "2m"***REMOVED***`***REMOVED******REMOVED***),
				cli: []string***REMOVED***"--config", "/my/config.file"***REMOVED***,
			***REMOVED***, exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(8), 120*time.Second),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:  defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s", "target": 20***REMOVED***], "vus": 10***REMOVED***`),
				env: []string***REMOVED***"K6_DURATION=15s"***REMOVED***,
				cli: []string***REMOVED***"--stage", ""***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(10), 15*time.Second),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(5), Duration: types.NullDurationFrom(50 * time.Second)***REMOVED***,
				cli:    []string***REMOVED***"--stage", "5s:5"***REMOVED***,
			***REMOVED***,
			//TODO: this shouldn't be a warning in the next version, but the result will be different
			exp***REMOVED***logWarning: true***REMOVED***, verifyConstLoopingVUs(I(5), 50*time.Second),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:     defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s", "target": 10***REMOVED***]***REMOVED***`),
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(5)***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyVarLoopingVUs(null.NewInt(5, true), buildStages(20, 10)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:     defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "20s", "target": 10***REMOVED***]***REMOVED***`),
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(5)***REMOVED***,
				env:    []string***REMOVED***"K6_VUS=15", "K6_ITERATIONS=15"***REMOVED***,
			***REMOVED***,
			exp***REMOVED***logWarning: true***REMOVED***, //TODO: this won't be a warning in the next version, but the result will be different
			verifySharedIters(I(15), I(15)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs:     defaultConfig(`***REMOVED***"stages": [***REMOVED***"duration": "11s", "target": 11***REMOVED***]***REMOVED***`),
				runner: &lib.Options***REMOVED***VUs: null.IntFrom(22)***REMOVED***,
				env:    []string***REMOVED***"K6_VUS=33"***REMOVED***,
				cli:    []string***REMOVED***"--stage", "44s:44", "-s", "55s:55"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***,
			verifyVarLoopingVUs(null.NewInt(33, true), buildStages(44, 44, 55, 55)),
		***REMOVED***,

		//TODO: test the future full overwriting of the duration/iterations/stages/execution options
		***REMOVED***
			opts***REMOVED***
				fs: defaultConfig(`***REMOVED***
					"execution": ***REMOVED*** "someKey": ***REMOVED***
						"type": "constant-looping-vus", "vus": 10, "duration": "60s", "interruptible": false,
						"iterationTimeout": "10s", "startTime": "70s", "env": ***REMOVED***"test": "mest"***REMOVED***, "exec": "someFunc"
					***REMOVED******REMOVED******REMOVED***`),
				env: []string***REMOVED***"K6_ITERATIONS=25"***REMOVED***,
				cli: []string***REMOVED***"--vus", "12"***REMOVED***,
			***REMOVED***,
			exp***REMOVED******REMOVED***, verifySharedIters(I(12), I(25)),
		***REMOVED***,
		***REMOVED***
			opts***REMOVED***
				fs: defaultConfig(`
					***REMOVED***
						"execution": ***REMOVED***
							"default": ***REMOVED***
								"type": "constant-looping-vus",
								"vus": 10,
								"duration": "60s"
							***REMOVED***
						***REMOVED***,
						"vus": 10,
						"duration": "60s"
					***REMOVED***`,
				),
			***REMOVED***,
			exp***REMOVED******REMOVED***, verifyConstLoopingVUs(I(10), 60*time.Second),
		***REMOVED***,
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
			opts***REMOVED***runner: &lib.Options***REMOVED***
				SystemTags: stats.NewSystemTagSet(stats.TagSubproto, stats.TagURL)***REMOVED***,
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
		***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--summary-trend-stats", `""`***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, func(t *testing.T, c Config) ***REMOVED***
			assert.Equal(t, []string***REMOVED******REMOVED***, c.Options.SummaryTrendStats)
		***REMOVED******REMOVED***,
		***REMOVED***
			opts***REMOVED***runner: &lib.Options***REMOVED***SummaryTrendStats: []string***REMOVED***"avg", "p(90)", "count"***REMOVED******REMOVED******REMOVED***,
			exp***REMOVED******REMOVED***,
			func(t *testing.T, c Config) ***REMOVED***
				assert.Equal(t, []string***REMOVED***"avg", "p(90)", "count"***REMOVED***, c.Options.SummaryTrendStats)
			***REMOVED***,
		***REMOVED***,
		//TODO: test for differences between flagsets
		//TODO: more tests in general, especially ones not related to execution parameters...
	***REMOVED***
***REMOVED***

func runTestCase(
	t *testing.T,
	testCase configConsolidationTestCase,
	newFlagSet flagSetInit,
	logHook *testutils.SimpleLogrusHook,
) ***REMOVED***
	t.Logf("Test with opts=%#v and exp=%#v\n", testCase.options, testCase.expected)
	logrus.SetOutput(testOutput***REMOVED***t***REMOVED***)
	logHook.Drain()

	restoreEnv := setEnv(t, testCase.options.env)
	defer restoreEnv()

	flagSet := newFlagSet()
	defer resetStickyGlobalVars()
	flagSet.SetOutput(testOutput***REMOVED***t***REMOVED***)
	//flagSet.PrintDefaults()

	cliErr := flagSet.Parse(testCase.options.cli)
	if testCase.expected.cliParseError ***REMOVED***
		require.Error(t, cliErr)
		return
	***REMOVED***
	require.NoError(t, cliErr)

	//TODO: remove these hacks when we improve the configuration...
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

	var runner lib.Runner
	if testCase.options.runner != nil ***REMOVED***
		runner = &lib.MiniRunner***REMOVED***Options: *testCase.options.runner***REMOVED***
	***REMOVED***
	if testCase.options.fs == nil ***REMOVED***
		t.Logf("Creating an empty FS for this test")
		testCase.options.fs = afero.NewMemMapFs() // create an empty FS if it wasn't supplied
	***REMOVED***

	consolidatedConfig, err := getConsolidatedConfig(testCase.options.fs, cliConf, runner)
	if testCase.expected.consolidationError ***REMOVED***
		require.Error(t, err)
		return
	***REMOVED***
	require.NoError(t, err)

	derivedConfig, err := deriveExecutionConfig(consolidatedConfig)
	if testCase.expected.derivationError ***REMOVED***
		require.Error(t, err)
		return
	***REMOVED***
	require.NoError(t, err)

	warnings := logHook.Drain()
	if testCase.expected.logWarning ***REMOVED***
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

func TestConfigConsolidation(t *testing.T) ***REMOVED***
	// This test and its subtests shouldn't be ran in parallel, since they unfortunately have
	// to mess with shared global objects (env vars, variables, the log, ... santa?)
	logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel***REMOVED******REMOVED***
	logrus.AddHook(&logHook)
	logrus.SetOutput(ioutil.Discard)
	defer logrus.SetOutput(os.Stderr)

	for tcNum, testCase := range getConfigConsolidationTestCases() ***REMOVED***
		flagSetInits := testCase.options.cliFlagSetInits
		if flagSetInits == nil ***REMOVED*** // handle the most common case
			flagSetInits = mostFlagSets()
		***REMOVED***
		for fsNum, flagSet := range flagSetInits ***REMOVED***
			// I want to paralelize this, but I cannot... due to global variables and other
			// questionable architectural choices... :|
			testCase, flagSet := testCase, flagSet
			t.Run(
				fmt.Sprintf("TestCase#%d_FlagSet#%d", tcNum, fsNum),
				func(t *testing.T) ***REMOVED*** runTestCase(t, testCase, flagSet, &logHook) ***REMOVED***,
			)
		***REMOVED***
	***REMOVED***
***REMOVED***
