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
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib/scheduler"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/types"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"

	"github.com/kelseyhightower/envconfig"
	"github.com/loadimpact/k6/lib"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

type testCmdData struct ***REMOVED***
	Name  string
	Tests []testCmdTest
***REMOVED***

type testCmdTest struct ***REMOVED***
	Args     []string
	Expected []string
	Name     string
***REMOVED***

func TestConfigCmd(t *testing.T) ***REMOVED***

	testdata := []testCmdData***REMOVED***
		***REMOVED***
			Name: "Out",

			Tests: []testCmdTest***REMOVED***
				***REMOVED***
					Name:     "NoArgs",
					Args:     []string***REMOVED***""***REMOVED***,
					Expected: []string***REMOVED******REMOVED***,
				***REMOVED***,
				***REMOVED***
					Name:     "SingleArg",
					Args:     []string***REMOVED***"--out", "influxdb=http://localhost:8086/k6"***REMOVED***,
					Expected: []string***REMOVED***"influxdb=http://localhost:8086/k6"***REMOVED***,
				***REMOVED***,
				***REMOVED***
					Name:     "MultiArg",
					Args:     []string***REMOVED***"--out", "influxdb=http://localhost:8086/k6", "--out", "json=test.json"***REMOVED***,
					Expected: []string***REMOVED***"influxdb=http://localhost:8086/k6", "json=test.json"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, data := range testdata ***REMOVED***
		t.Run(data.Name, func(t *testing.T) ***REMOVED***
			for _, test := range data.Tests ***REMOVED***
				t.Run(`"`+test.Name+`"`, func(t *testing.T) ***REMOVED***
					fs := configFlagSet()
					fs.AddFlagSet(optionFlagSet())
					assert.NoError(t, fs.Parse(test.Args))

					config, err := getConfig(fs)
					assert.NoError(t, err)
					assert.Equal(t, test.Expected, config.Out)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

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
	assert.Equal(t, null.NewBool(false, false), perVuIters.Interruptible)
	assert.Equal(t, null.NewInt(1, false), perVuIters.Iterations)
	assert.Equal(t, null.NewInt(1, false), perVuIters.VUs)
	//TODO: verify shortcut options as well?
***REMOVED***

func verifySharedIters(vus, iters int64) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		sched := c.Execution[lib.DefaultSchedulerName]
		require.NotEmpty(t, sched)
		require.IsType(t, scheduler.SharedIteationsConfig***REMOVED******REMOVED***, sched)
		sharedIterConfig, ok := sched.(scheduler.SharedIteationsConfig)
		require.True(t, ok)
		assert.Equal(t, null.NewInt(vus, true), sharedIterConfig.VUs)
		assert.Equal(t, null.NewInt(iters, true), sharedIterConfig.Iterations)
		//TODO: verify shortcut options as well?
	***REMOVED***
***REMOVED***

func verifyConstantLoopingVUs(vus int64, duration time.Duration) func(t *testing.T, c Config) ***REMOVED***
	return func(t *testing.T, c Config) ***REMOVED***
		sched := c.Execution[lib.DefaultSchedulerName]
		require.NotEmpty(t, sched)
		require.IsType(t, scheduler.ConstantLoopingVUsConfig***REMOVED******REMOVED***, sched)
		clvc, ok := sched.(scheduler.ConstantLoopingVUsConfig)
		require.True(t, ok)
		assert.Equal(t, null.NewBool(true, false), clvc.Interruptible)
		assert.Equal(t, null.NewInt(vus, true), clvc.VUs)
		assert.Equal(t, types.NullDurationFrom(duration), clvc.Duration)
		//TODO: verify shortcut options as well?
	***REMOVED***
***REMOVED***

func mostFlagSets() []*pflag.FlagSet ***REMOVED***
	//TODO: make this unneccesary... currently these are the only commands in which
	// getConsolidatedConfig() is used, but they also have differences in their CLI flags :/
	// sigh... compromises...
	return []*pflag.FlagSet***REMOVED***runCmdFlagSet(), archiveCmdFlagSet(), cloudCmdFlagSet()***REMOVED***
***REMOVED***

type opts struct ***REMOVED***
	cliFlagSets []*pflag.FlagSet
	cli         []string
	env         []string
	runner      *lib.Options
	//TODO: test the JSON config as well... after most of https://github.com/loadimpact/k6/issues/883#issuecomment-468646291 is fixed
***REMOVED***

// exp contains the different events or errors we expect our test case to trigger.
// for space and clarity, we use the fact that by default, all of the struct values are false
type exp struct ***REMOVED***
	cliParseError      bool
	cliReadError       bool
	consolidationError bool
	validationErrors   bool
	logWarning         bool //TODO: remove in the next version?
***REMOVED***

// A hell of a complicated test case, that still doesn't test things fully...
type configConsolidationTestCase struct ***REMOVED***
	options         opts
	expected        exp
	customValidator func(t *testing.T, c Config)
***REMOVED***

var configConsolidationTestCases = []configConsolidationTestCase***REMOVED***
	// Check that no options will result in 1 VU 1 iter value for execution
	***REMOVED***opts***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyOneIterPerOneVU***REMOVED***,
	// Verify some CLI errors
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--blah", "blah"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--duration", "blah"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--iterations", "blah"***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--execution", ""***REMOVED******REMOVED***, exp***REMOVED***cliParseError: true***REMOVED***, nil***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--stage", "10:20s"***REMOVED******REMOVED***, exp***REMOVED***cliReadError: true***REMOVED***, nil***REMOVED***,
	// Check if CLI shortcuts generate correct execution values
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "1", "--iterations", "5"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(1, 5)***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "2", "-i", "6"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(2, 6)***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-d", "30s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstantLoopingVUs(3, 30*time.Second)***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "4", "--duration", "60s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstantLoopingVUs(4, 1*time.Minute)***REMOVED***,
	//TODO: verify stages
	// This should get a validation error since VUs are more than the shared iterations
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"--vus", "10", "-i", "6"***REMOVED******REMOVED***, exp***REMOVED***validationErrors: true***REMOVED***, verifySharedIters(10, 6)***REMOVED***,
	// These should emit a warning
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "1", "-i", "6", "-d", "10s"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "2", "-d", "10s", "-s", "10s:20"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-i", "5", "-s", "10s:20"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil***REMOVED***,
	***REMOVED***opts***REMOVED***cli: []string***REMOVED***"-u", "3", "-d", "0"***REMOVED******REMOVED***, exp***REMOVED***logWarning: true***REMOVED***, nil***REMOVED***,
	// Test if environment variable shortcuts are working as expected
	***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_VUS=5", "K6_ITERATIONS=15"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifySharedIters(5, 15)***REMOVED***,
	***REMOVED***opts***REMOVED***env: []string***REMOVED***"K6_VUS=10", "K6_DURATION=20s"***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyConstantLoopingVUs(10, 20*time.Second)***REMOVED***,

	//TODO: test combinations between options and levels
	//TODO: test the future full overwriting of the duration/iterations/stages/execution options

	// Just in case, verify that no options will result in the same 1 vu 1 iter config
	***REMOVED***opts***REMOVED******REMOVED***, exp***REMOVED******REMOVED***, verifyOneIterPerOneVU***REMOVED***,
	//TODO: test for differences between flagsets
	//TODO: more tests in general...
***REMOVED***

func TestConfigConsolidation(t *testing.T) ***REMOVED***
	// This test and its subtests shouldn't be ran in parallel, since they unfortunately have
	// to mess with shared global objects (env vars, variables, the log, ... santa?)
	logHook := testutils.SimpleLogrusHook***REMOVED***HookedLevels: []log.Level***REMOVED***log.WarnLevel***REMOVED******REMOVED***
	log.AddHook(&logHook)
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)

	runTestCase := func(t *testing.T, testCase configConsolidationTestCase, flagSet *pflag.FlagSet) ***REMOVED***
		t.Logf("Test with opts=%#v and exp=%#v\n", testCase.options, testCase.expected)
		logHook.Drain()

		restoreEnv := setEnv(t, testCase.options.env)
		defer restoreEnv()

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
		fs := afero.NewMemMapFs() //TODO: test JSON configs as well!
		result, err := getConsolidatedConfig(fs, cliConf, runner)
		if testCase.expected.consolidationError ***REMOVED***
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

		validationErrors := result.Validate()
		if testCase.expected.validationErrors ***REMOVED***
			assert.NotEmpty(t, validationErrors)
		***REMOVED*** else ***REMOVED***
			assert.Empty(t, validationErrors)
		***REMOVED***

		if testCase.customValidator != nil ***REMOVED***
			testCase.customValidator(t, result)
		***REMOVED***
	***REMOVED***

	for tcNum, testCase := range configConsolidationTestCases ***REMOVED***
		flagSets := testCase.options.cliFlagSets
		if flagSets == nil ***REMOVED*** // handle the most common case
			flagSets = mostFlagSets()
		***REMOVED***
		for fsNum, flagSet := range flagSets ***REMOVED***
			// I want to paralelize this, but I cannot... due to global variables and other
			// questionable architectural choices... :|
			t.Run(
				fmt.Sprintf("TestCase#%d_FlagSet#%d", tcNum, fsNum),
				func(t *testing.T) ***REMOVED*** runTestCase(t, testCase, flagSet) ***REMOVED***,
			)
		***REMOVED***
	***REMOVED***
***REMOVED***
func TestConfigEnv(t *testing.T) ***REMOVED***
	testdata := map[struct***REMOVED*** Name, Key string ***REMOVED***]map[string]func(Config)***REMOVED***
		***REMOVED***"Linger", "K6_LINGER"***REMOVED***: ***REMOVED***
			"":      func(c Config) ***REMOVED*** assert.Equal(t, null.Bool***REMOVED******REMOVED***, c.Linger) ***REMOVED***,
			"true":  func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(true), c.Linger) ***REMOVED***,
			"false": func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(false), c.Linger) ***REMOVED***,
		***REMOVED***,
		***REMOVED***"NoUsageReport", "K6_NO_USAGE_REPORT"***REMOVED***: ***REMOVED***
			"":      func(c Config) ***REMOVED*** assert.Equal(t, null.Bool***REMOVED******REMOVED***, c.NoUsageReport) ***REMOVED***,
			"true":  func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(true), c.NoUsageReport) ***REMOVED***,
			"false": func(c Config) ***REMOVED*** assert.Equal(t, null.BoolFrom(false), c.NoUsageReport) ***REMOVED***,
		***REMOVED***,
		***REMOVED***"Out", "K6_OUT"***REMOVED***: ***REMOVED***
			"":         func(c Config) ***REMOVED*** assert.Equal(t, []string***REMOVED***""***REMOVED***, c.Out) ***REMOVED***,
			"influxdb": func(c Config) ***REMOVED*** assert.Equal(t, []string***REMOVED***"influxdb"***REMOVED***, c.Out) ***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for field, data := range testdata ***REMOVED***
		os.Clearenv()
		t.Run(field.Name, func(t *testing.T) ***REMOVED***
			for value, fn := range data ***REMOVED***
				t.Run(`"`+value+`"`, func(t *testing.T) ***REMOVED***
					assert.NoError(t, os.Setenv(field.Key, value))
					var config Config
					assert.NoError(t, envconfig.Process("k6", &config))
					fn(config)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConfigApply(t *testing.T) ***REMOVED***
	t.Run("Linger", func(t *testing.T) ***REMOVED***
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Linger: null.BoolFrom(true)***REMOVED***)
		assert.Equal(t, null.BoolFrom(true), conf.Linger)
	***REMOVED***)
	t.Run("NoUsageReport", func(t *testing.T) ***REMOVED***
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***NoUsageReport: null.BoolFrom(true)***REMOVED***)
		assert.Equal(t, null.BoolFrom(true), conf.NoUsageReport)
	***REMOVED***)
	t.Run("Out", func(t *testing.T) ***REMOVED***
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Out: []string***REMOVED***"influxdb"***REMOVED******REMOVED***)
		assert.Equal(t, []string***REMOVED***"influxdb"***REMOVED***, conf.Out)

		conf = Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Out: []string***REMOVED***"influxdb", "json"***REMOVED******REMOVED***)
		assert.Equal(t, []string***REMOVED***"influxdb", "json"***REMOVED***, conf.Out)
	***REMOVED***)
***REMOVED***
