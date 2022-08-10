package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_getLevels(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := [...]struct ***REMOVED***
		level  string
		err    bool
		levels []logrus.Level
	***REMOVED******REMOVED***
		***REMOVED***
			level: "info",
			err:   false,
			levels: []logrus.Level***REMOVED***
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
				logrus.InfoLevel,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			level: "error",
			err:   false,
			levels: []logrus.Level***REMOVED***
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			level:  "tea",
			err:    true,
			levels: nil,
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		test := test
		t.Run(test.level, func(t *testing.T) ***REMOVED***
			t.Parallel()

			levels, err := parseLevels(test.level)

			if test.err ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***

			require.NoError(t, err)
			require.Equal(t, test.levels, levels)
		***REMOVED***)
	***REMOVED***
***REMOVED***
