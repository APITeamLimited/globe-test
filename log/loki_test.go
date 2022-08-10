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

			res, err := LokiFromConfigLine(context.Background(), nil, test.line, make(chan struct***REMOVED******REMOVED***))

			if test.err ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			test.res.client = res.(*lokiHook).client
			test.res.ch = res.(*lokiHook).ch
			test.res.lokiStopped = res.(*lokiHook).lokiStopped
			require.Equal(t, &test.res, res)
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

func TestFilterLabels(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		allowedLabels  []string
		labels         map[string]string
		expectedLabels map[string]string
		msg            string
		result         string
	***REMOVED******REMOVED***
		***REMOVED***
			allowedLabels:  []string***REMOVED***"a", "b"***REMOVED***,
			labels:         map[string]string***REMOVED***"a": "1", "b": "2", "d": "3", "c": "4", "e": "5"***REMOVED***,
			expectedLabels: map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED***,
			msg:            "some msg",
			result:         "some msg c=4 d=3 e=5",
		***REMOVED***,
		***REMOVED***
			allowedLabels:  []string***REMOVED***"a", "b"***REMOVED***,
			labels:         map[string]string***REMOVED***"d": "3", "c": "4", "e": "5"***REMOVED***,
			expectedLabels: map[string]string***REMOVED******REMOVED***,
			msg:            "some msg",
			result:         "some msg c=4 d=3 e=5",
		***REMOVED***,
		***REMOVED***
			allowedLabels:  []string***REMOVED***"a", "b"***REMOVED***,
			labels:         map[string]string***REMOVED***"a": "1", "d": "3", "c": "4", "e": "5"***REMOVED***,
			expectedLabels: map[string]string***REMOVED***"a": "1"***REMOVED***,
			msg:            "some msg",
			result:         "some msg c=4 d=3 e=5",
		***REMOVED***,
	***REMOVED***

	for i, c := range cases ***REMOVED***
		c := c
		t.Run(fmt.Sprint(i), func(t *testing.T) ***REMOVED***
			h := &lokiHook***REMOVED******REMOVED***
			h.allowedLabels = c.allowedLabels
			result := h.filterLabels(c.labels, c.msg)
			assert.Equal(t, c.result, result)
			assert.Equal(t, c.expectedLabels, c.labels)
		***REMOVED***)
	***REMOVED***
***REMOVED***
