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
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyslogFromConfigLine(t *testing.T) ***REMOVED***
	t.Parallel()
	tests := [...]struct ***REMOVED***
		line string
		err  bool
		res  lokiHook
	***REMOVED******REMOVED***
		***REMOVED***
			line: "loki", // default settings
			res: lokiHook***REMOVED***
				ctx:           context.Background(),
				addr:          "http://127.0.0.1:3100/loki/api/v1/push",
				limit:         100,
				pushPeriod:    time.Second * 1,
				levels:        logrus.AllLevels,
				msgMaxSize:    1024 * 1024,
				droppedLabels: map[string]string***REMOVED***"level": "warning"***REMOVED***,
				droppedMsg:    "k6 dropped %d log messages because they were above the limit of %d messages / %s",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			line: "loki=somewhere:1233,label.something=else,label.foo=bar,limit=32,level=info,allowedLabels=[something],pushPeriod=5m32s,msgMaxSize=1231",
			res: lokiHook***REMOVED***
				ctx:           context.Background(),
				addr:          "somewhere:1233",
				limit:         32,
				pushPeriod:    time.Minute*5 + time.Second*32,
				levels:        logrus.AllLevels[:5],
				labels:        [][2]string***REMOVED******REMOVED***"something", "else"***REMOVED***, ***REMOVED***"foo", "bar"***REMOVED******REMOVED***,
				msgMaxSize:    1231,
				allowedLabels: []string***REMOVED***"something"***REMOVED***,
				droppedLabels: map[string]string***REMOVED***"something": "else"***REMOVED***,
				droppedMsg:    "k6 dropped %d log messages because they were above the limit of %d messages / %s foo=bar level=warning",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			line: "lokino",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line: "loki=something,limit=word",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line: "loki=something,level=notlevel",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line: "loki=something,unknownoption",
			err:  true,
		***REMOVED***,
		***REMOVED***
			line: "loki=something,label=somethng",
			err:  true,
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		test := test
		t.Run(test.line, func(t *testing.T) ***REMOVED***
			// no parallel because this is way too fast and parallel will only slow it down

			res, err := LokiFromConfigLine(context.Background(), nil, test.line)

			if test.err ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			test.res.client = res.(*lokiHook).client
			test.res.ch = res.(*lokiHook).ch
			require.Equal(t, &test.res, res)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestParseArray(t *testing.T) ***REMOVED***
	cases := [...]struct ***REMOVED***
		key, value string
		i          int
		args       []string
		result     []string
		resultI    int
		err        bool
	***REMOVED******REMOVED***
		***REMOVED***
			key:     "test",
			value:   "[some",
			i:       2,
			args:    []string***REMOVED***"else=asa", "e=s", "test=[some", "else]"***REMOVED***,
			result:  []string***REMOVED***"some", "else"***REMOVED***,
			resultI: 3,
		***REMOVED***,
		***REMOVED***
			key:   "test",
			value: "[some",
			i:     2,
			args:  []string***REMOVED***"else=asa", "e=s", "test=[some", "else"***REMOVED***,
			err:   true,
		***REMOVED***,
		***REMOVED***
			key:   "test",
			value: "[some",
			i:     2,
			args:  []string***REMOVED***"else=asa", "e=s", "test=[some", "", "s]"***REMOVED***,
			err:   true,
		***REMOVED***,
		***REMOVED***
			key:   "test",
			value: "some",
			i:     2,
			args:  []string***REMOVED***"else=asa", "e=s", "test=some", "else]"***REMOVED***,
			err:   true,
		***REMOVED***,
		***REMOVED***
			key:     "test",
			value:   "",
			i:       2,
			args:    []string***REMOVED***"else=asa", "e=s", "test=", "sdasa"***REMOVED***,
			result:  []string***REMOVED******REMOVED***,
			resultI: 2,
		***REMOVED***,
		***REMOVED***
			key:     "test",
			value:   "[some]",
			i:       2,
			args:    []string***REMOVED***"else=asa", "e=s", "test=[some]", "else=sa"***REMOVED***,
			result:  []string***REMOVED***"some"***REMOVED***,
			resultI: 2,
		***REMOVED***,
	***REMOVED***

	for i, c := range cases ***REMOVED***
		c := c
		t.Run(fmt.Sprint(i), func(t *testing.T) ***REMOVED***
			result, i, err := parseArray(c.key, c.value, c.i, c.args)
			assert.Equal(t, c.result, result)
			assert.Equal(t, c.resultI, i)
			if c.err ***REMOVED***
				assert.Error(t, err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLogEntryMarshal(t *testing.T) ***REMOVED***
	entry := logEntry***REMOVED***
		t:   9223372036854775807, // the actual max
		msg: "something",
	***REMOVED***
	expected := []byte(`["9223372036854775807","something"]`)
	s, err := json.Marshal(entry)
	require.NoError(t, err)

	require.Equal(t, expected, s)
***REMOVED***
