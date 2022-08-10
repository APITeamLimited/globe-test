package csv

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"gopkg.in/guregu/null.v3"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/types"
)

func TestNewConfig(t *testing.T) ***REMOVED***
	config := NewConfig()
	assert.Equal(t, "file.csv", config.FileName.String)
	assert.Equal(t, "1s", config.SaveInterval.String())
	assert.Equal(t, "unix", config.TimeFormat.String)
***REMOVED***

func TestApply(t *testing.T) ***REMOVED***
	configs := []Config***REMOVED***
		***REMOVED***
			FileName:     null.StringFrom(""),
			SaveInterval: types.NullDurationFrom(2 * time.Second),
			TimeFormat:   null.StringFrom("unix"),
		***REMOVED***,
		***REMOVED***
			FileName:     null.StringFrom("newPath"),
			SaveInterval: types.NewNullDuration(time.Duration(1), false),
			TimeFormat:   null.StringFrom("rfc3339"),
		***REMOVED***,
	***REMOVED***
	expected := []struct ***REMOVED***
		FileName     string
		SaveInterval string
		TimeFormat   string
	***REMOVED******REMOVED***
		***REMOVED***
			FileName:     "",
			SaveInterval: "2s",
			TimeFormat:   "unix",
		***REMOVED***,
		***REMOVED***
			FileName:     "newPath",
			SaveInterval: "1s",
			TimeFormat:   "rfc3339",
		***REMOVED***,
	***REMOVED***

	for i := range configs ***REMOVED***
		config := configs[i]
		expected := expected[i]
		t.Run(expected.FileName+"_"+expected.SaveInterval, func(t *testing.T) ***REMOVED***
			baseConfig := NewConfig()
			baseConfig = baseConfig.Apply(config)

			assert.Equal(t, expected.FileName, baseConfig.FileName.String)
			assert.Equal(t, expected.SaveInterval, baseConfig.SaveInterval.String())
			assert.Equal(t, expected.TimeFormat, baseConfig.TimeFormat.String)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestParseArg(t *testing.T) ***REMOVED***
	cases := map[string]struct ***REMOVED***
		config             Config
		expectedLogEntries []string
		expectedErr        bool
	***REMOVED******REMOVED***
		"test_file.csv": ***REMOVED***
			config: Config***REMOVED***
				FileName:     null.StringFrom("test_file.csv"),
				SaveInterval: types.NewNullDuration(1*time.Second, false),
				TimeFormat:   null.NewString("unix", false),
			***REMOVED***,
		***REMOVED***,
		"save_interval=5s": ***REMOVED***
			config: Config***REMOVED***
				FileName:     null.NewString("file.csv", false),
				SaveInterval: types.NullDurationFrom(5 * time.Second),
				TimeFormat:   null.NewString("unix", false),
			***REMOVED***,
			expectedLogEntries: []string***REMOVED***
				"CSV output argument 'save_interval' is deprecated, please use 'saveInterval' instead.",
			***REMOVED***,
		***REMOVED***,
		"saveInterval=5s": ***REMOVED***
			config: Config***REMOVED***
				FileName:     null.NewString("file.csv", false),
				SaveInterval: types.NullDurationFrom(5 * time.Second),
				TimeFormat:   null.NewString("unix", false),
			***REMOVED***,
		***REMOVED***,
		"file_name=test.csv,save_interval=5s": ***REMOVED***
			config: Config***REMOVED***
				FileName:     null.StringFrom("test.csv"),
				SaveInterval: types.NullDurationFrom(5 * time.Second),
				TimeFormat:   null.NewString("unix", false),
			***REMOVED***,
			expectedLogEntries: []string***REMOVED***
				"CSV output argument 'file_name' is deprecated, please use 'fileName' instead.",
				"CSV output argument 'save_interval' is deprecated, please use 'saveInterval' instead.",
			***REMOVED***,
		***REMOVED***,
		"fileName=test.csv,save_interval=5s": ***REMOVED***
			config: Config***REMOVED***
				FileName:     null.StringFrom("test.csv"),
				SaveInterval: types.NullDurationFrom(5 * time.Second),
				TimeFormat:   null.NewString("unix", false),
			***REMOVED***,
			expectedLogEntries: []string***REMOVED***
				"CSV output argument 'save_interval' is deprecated, please use 'saveInterval' instead.",
			***REMOVED***,
		***REMOVED***,
		"filename=test.csv,save_interval=5s": ***REMOVED***
			expectedErr: true,
		***REMOVED***,
		"fileName=test.csv,timeFormat=rfc3339": ***REMOVED***
			config: Config***REMOVED***
				FileName:     null.StringFrom("test.csv"),
				SaveInterval: types.NewNullDuration(1*time.Second, false),
				TimeFormat:   null.StringFrom("rfc3339"),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for arg, testCase := range cases ***REMOVED***
		arg := arg
		testCase := testCase

		testLogger, hook := test.NewNullLogger()
		testLogger.SetOutput(testutils.NewTestOutput(t))

		t.Run(arg, func(t *testing.T) ***REMOVED***
			config, err := ParseArg(arg, testLogger)

			if testCase.expectedErr ***REMOVED***
				assert.Error(t, err)
				return
			***REMOVED***

			require.NoError(t, err)
			assert.Equal(t, testCase.config, config)

			var entries []string
			for _, v := range hook.AllEntries() ***REMOVED***
				assert.Equal(t, v.Level, logrus.WarnLevel)
				entries = append(entries, v.Message)
			***REMOVED***
			assert.ElementsMatch(t, entries, testCase.expectedLogEntries)
		***REMOVED***)
	***REMOVED***
***REMOVED***
