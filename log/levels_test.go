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
