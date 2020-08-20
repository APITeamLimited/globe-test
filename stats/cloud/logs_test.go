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

package cloud

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib/testutils"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgParsing(t *testing.T) ***REMOVED***
	m := `***REMOVED***
  "streams": [
    ***REMOVED***
      "stream": ***REMOVED***
      	"key1": "value1",
	   	"key2": "value2"
      ***REMOVED***,
      "values": [
        [
      	"1598282752000000000",
		"something to log"
        ]
      ]
    ***REMOVED***
  ],
  "dropped_entries": [
    ***REMOVED***
      "labels": ***REMOVED***
      	"key3": "value1",
	   	"key4": "value2"
      ***REMOVED***,
      "timestamp": "1598282752000000000"
    ***REMOVED***
  ]
***REMOVED***
`
	expectMsg := msg***REMOVED***
		Streams: []msgStreams***REMOVED***
			***REMOVED***
				Stream: map[string]string***REMOVED***"key1": "value1", "key2": "value2"***REMOVED***,
				Values: [][2]string***REMOVED******REMOVED***"1598282752000000000", "something to log"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		DroppedEntries: []msgDroppedEntries***REMOVED***
			***REMOVED***
				Labels:    map[string]string***REMOVED***"key3": "value1", "key4": "value2"***REMOVED***,
				Timestamp: "1598282752000000000",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	var message msg
	require.NoError(t, easyjson.Unmarshal([]byte(m), &message))
	require.Equal(t, expectMsg, message)
***REMOVED***

func TestMSGLog(t *testing.T) ***REMOVED***
	expectMsg := msg***REMOVED***
		Streams: []msgStreams***REMOVED***
			***REMOVED***
				Stream: map[string]string***REMOVED***"key1": "value1", "key2": "value2"***REMOVED***,
				Values: [][2]string***REMOVED******REMOVED***"1598282752000000000", "something to log"***REMOVED******REMOVED***,
			***REMOVED***,
			***REMOVED***
				Stream: map[string]string***REMOVED***"key1": "value1", "key2": "value2", "level": "warn"***REMOVED***,
				Values: [][2]string***REMOVED******REMOVED***"1598282752000000000", "something else log"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
		DroppedEntries: []msgDroppedEntries***REMOVED***
			***REMOVED***
				Labels:    map[string]string***REMOVED***"key3": "value1", "key4": "value2", "level": "panic"***REMOVED***,
				Timestamp: "1598282752000000000",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	logger := logrus.New()
	logger.Out = ioutil.Discard
	hook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: logrus.AllLevels***REMOVED***
	logger.AddHook(hook)
	expectMsg.Log(logger)
	logLines := hook.Drain()
	assert.Equal(t, 4, len(logLines))
	expectTime := time.Unix(0, 1598282752000000000)
	for i, entry := range logLines ***REMOVED***
		var expectedMsg string
		switch i ***REMOVED***
		case 0:
			expectedMsg = "something to log"
		case 1:
			expectedMsg = "last message had unknown level "
		case 2:
			expectedMsg = "something else log"
		case 3:
			expectedMsg = "dropped"
		***REMOVED***
		require.Equal(t, expectedMsg, entry.Message)
		require.Equal(t, expectTime, entry.Time)
	***REMOVED***
***REMOVED***
