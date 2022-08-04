package cmd

import (
	"testing"
	"time"

	"github.com/mstoykov/envconfig"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/lib/types"
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
	t.Parallel()
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
			t.Parallel()
			for _, test := range data.Tests ***REMOVED***
				t.Run(`"`+test.Name+`"`, func(t *testing.T) ***REMOVED***
					t.Parallel()
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

func TestConfigEnv(t *testing.T) ***REMOVED***
	t.Parallel()
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
			"":         func(c Config) ***REMOVED*** assert.Equal(t, []string***REMOVED******REMOVED***, c.Out) ***REMOVED***,
			"influxdb": func(c Config) ***REMOVED*** assert.Equal(t, []string***REMOVED***"influxdb"***REMOVED***, c.Out) ***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for field, data := range testdata ***REMOVED***
		field, data := field, data
		t.Run(field.Name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			for value, fn := range data ***REMOVED***
				value, fn := value, fn
				t.Run(`"`+value+`"`, func(t *testing.T) ***REMOVED***
					t.Parallel()
					var config Config
					assert.NoError(t, envconfig.Process("", &config, func(key string) (string, bool) ***REMOVED***
						if key == field.Key ***REMOVED***
							return value, true
						***REMOVED***
						return "", false
					***REMOVED***))
					fn(config)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestConfigApply(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("Linger", func(t *testing.T) ***REMOVED***
		t.Parallel()
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Linger: null.BoolFrom(true)***REMOVED***)
		assert.Equal(t, null.BoolFrom(true), conf.Linger)
	***REMOVED***)
	t.Run("NoUsageReport", func(t *testing.T) ***REMOVED***
		t.Parallel()
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***NoUsageReport: null.BoolFrom(true)***REMOVED***)
		assert.Equal(t, null.BoolFrom(true), conf.NoUsageReport)
	***REMOVED***)
	t.Run("Out", func(t *testing.T) ***REMOVED***
		t.Parallel()
		conf := Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Out: []string***REMOVED***"influxdb"***REMOVED******REMOVED***)
		assert.Equal(t, []string***REMOVED***"influxdb"***REMOVED***, conf.Out)

		conf = Config***REMOVED******REMOVED***.Apply(Config***REMOVED***Out: []string***REMOVED***"influxdb", "json"***REMOVED******REMOVED***)
		assert.Equal(t, []string***REMOVED***"influxdb", "json"***REMOVED***, conf.Out)
	***REMOVED***)
***REMOVED***

func TestDeriveAndValidateConfig(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		name   string
		conf   Config
		isExec bool
		err    string
	***REMOVED******REMOVED***
		***REMOVED***"defaultOK", Config***REMOVED******REMOVED***, true, ""***REMOVED***,
		***REMOVED***
			"defaultErr",
			Config***REMOVED******REMOVED***,
			false,
			"executor default: function 'default' not found in exports",
		***REMOVED***,
		***REMOVED***
			"nonDefaultOK", Config***REMOVED***Options: lib.Options***REMOVED***Scenarios: lib.ScenarioConfigs***REMOVED***
				"per_vu_iters": executor.PerVUIterationsConfig***REMOVED***
					BaseConfig: executor.BaseConfig***REMOVED***
						Name: "per_vu_iters", Type: "per-vu-iterations", Exec: null.StringFrom("nonDefault"),
					***REMOVED***,
					VUs:         null.IntFrom(1),
					Iterations:  null.IntFrom(1),
					MaxDuration: types.NullDurationFrom(time.Second),
				***REMOVED***,
			***REMOVED******REMOVED******REMOVED***, true, "",
		***REMOVED***,
		***REMOVED***
			"nonDefaultErr",
			Config***REMOVED***Options: lib.Options***REMOVED***Scenarios: lib.ScenarioConfigs***REMOVED***
				"per_vu_iters": executor.PerVUIterationsConfig***REMOVED***
					BaseConfig: executor.BaseConfig***REMOVED***
						Name: "per_vu_iters", Type: "per-vu-iterations", Exec: null.StringFrom("nonDefaultErr"),
					***REMOVED***,
					VUs:         null.IntFrom(1),
					Iterations:  null.IntFrom(1),
					MaxDuration: types.NullDurationFrom(time.Second),
				***REMOVED***,
			***REMOVED******REMOVED******REMOVED***,
			false,
			"executor per_vu_iters: function 'nonDefaultErr' not found in exports",
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			_, err := deriveAndValidateConfig(tc.conf,
				func(_ string) bool ***REMOVED*** return tc.isExec ***REMOVED***, nil)
			if tc.err != "" ***REMOVED***
				var ecerr errext.HasExitCode
				assert.ErrorAs(t, err, &ecerr)
				assert.Equal(t, exitcodes.InvalidConfig, ecerr.ExitCode())
				assert.Contains(t, err.Error(), tc.err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
