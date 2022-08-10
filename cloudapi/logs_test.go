package cloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/lib/testutils/httpmultibin"
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

func TestRetry(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("Success", func(t *testing.T) ***REMOVED***
		t.Parallel()

		tests := []struct ***REMOVED***
			name     string
			attempts int
			expWaits []time.Duration // pow(abs(interval), attempt index)
		***REMOVED******REMOVED***
			***REMOVED***
				name:     "NoRetry",
				attempts: 1,
			***REMOVED***,
			***REMOVED***
				name:     "TwoAttempts",
				attempts: 2,
				expWaits: []time.Duration***REMOVED***5 * time.Second***REMOVED***,
			***REMOVED***,
			***REMOVED***
				name:     "MaximumExceeded",
				attempts: 4,
				expWaits: []time.Duration***REMOVED***5 * time.Second, 25 * time.Second, 2 * time.Minute***REMOVED***,
			***REMOVED***,
			***REMOVED***
				name:     "AttemptsLimit",
				attempts: 5,
				expWaits: []time.Duration***REMOVED***5 * time.Second, 25 * time.Second, 2 * time.Minute, 2 * time.Minute***REMOVED***,
			***REMOVED***,
		***REMOVED***

		for _, tt := range tests ***REMOVED***
			t.Run(tt.name, func(t *testing.T) ***REMOVED***
				var sleepRequests []time.Duration
				// sleepCollector tracks the request duration value for sleep requests.
				sleepCollector := sleeperFunc(func(d time.Duration) ***REMOVED***
					sleepRequests = append(sleepRequests, d)
				***REMOVED***)

				var iterations int
				err := retry(sleepCollector, 5, 5*time.Second, 2*time.Minute, func() error ***REMOVED***
					iterations++
					if iterations < tt.attempts ***REMOVED***
						return fmt.Errorf("unexpected error")
					***REMOVED***
					return nil
				***REMOVED***)
				require.NoError(t, err)
				require.Equal(t, tt.attempts, iterations)
				require.Equal(t, len(tt.expWaits), len(sleepRequests))

				// the added random milliseconds makes difficult to know the exact value
				// so it asserts that expwait <= actual <= expwait + 1s
				for i, expwait := range tt.expWaits ***REMOVED***
					assert.GreaterOrEqual(t, sleepRequests[i], expwait)
					assert.LessOrEqual(t, sleepRequests[i], expwait+(1*time.Second))
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Fail", func(t *testing.T) ***REMOVED***
		t.Parallel()

		mock := sleeperFunc(func(time.Duration) ***REMOVED*** /* noop - nowait */ ***REMOVED***)
		err := retry(mock, 5, 5*time.Second, 30*time.Second, func() error ***REMOVED***
			return fmt.Errorf("unexpected error")
		***REMOVED***)

		assert.Error(t, err, "unexpected error")
	***REMOVED***)
***REMOVED***

func TestStreamLogsToLogger(t *testing.T) ***REMOVED***
	t.Parallel()

	// It registers an handler for the logtail endpoint
	// It upgrades as websocket the HTTP handler and invokes the provided callback.
	logtailHandleFunc := func(tb *httpmultibin.HTTPMultiBin, fn func(*websocket.Conn, *http.Request)) ***REMOVED***
		upgrader := websocket.Upgrader***REMOVED***
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		***REMOVED***
		tb.Mux.HandleFunc("/api/v1/tail", func(w http.ResponseWriter, req *http.Request) ***REMOVED***
			conn, err := upgrader.Upgrade(w, req, nil)
			require.NoError(t, err)

			fn(conn, req)
			_ = conn.Close()
		***REMOVED***)
	***REMOVED***

	// a basic config with the logtail endpoint set
	configFromHTTPMultiBin := func(tb *httpmultibin.HTTPMultiBin) Config ***REMOVED***
		wsurl := strings.TrimPrefix(tb.ServerHTTP.URL, "http://")
		return Config***REMOVED***
			LogsTailURL: null.NewString(fmt.Sprintf("ws://%s/api/v1/tail", wsurl), false),
		***REMOVED***
	***REMOVED***

	// get all messages from the mocked logger
	logLines := func(hook *testutils.SimpleLogrusHook) (lines []string) ***REMOVED***
		for _, e := range hook.Drain() ***REMOVED***
			lines = append(lines, e.Message)
		***REMOVED***
		return
	***REMOVED***

	generateLogline := func(key string, ts uint64, msg string) string ***REMOVED***
		return fmt.Sprintf(`***REMOVED***"streams":[***REMOVED***"stream":***REMOVED***"key":%q,"level":"warn"***REMOVED***,"values":[["%d",%q]]***REMOVED***],"dropped_entities":[]***REMOVED***`, key, ts, msg)
	***REMOVED***

	t.Run("Success", func(t *testing.T) ***REMOVED***
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tb := httpmultibin.NewHTTPMultiBin(t)
		logtailHandleFunc(tb, func(conn *websocket.Conn, _ *http.Request) ***REMOVED***
			rawmsg := json.RawMessage(generateLogline("stream1", 1598282752000000000, "logline1"))
			err := conn.WriteJSON(rawmsg)
			require.NoError(t, err)

			rawmsg = json.RawMessage(generateLogline("stream2", 1598282752000000001, "logline2"))
			err = conn.WriteJSON(rawmsg)
			require.NoError(t, err)

			// wait the flush on the network
			time.Sleep(5 * time.Millisecond)
			cancel()
		***REMOVED***)

		logger := logrus.New()
		logger.Out = ioutil.Discard
		hook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: logrus.AllLevels***REMOVED***
		logger.AddHook(hook)

		c := configFromHTTPMultiBin(tb)
		err := c.StreamLogsToLogger(ctx, logger, "ref_id", 0)
		require.NoError(t, err)

		assert.Equal(t, []string***REMOVED***"logline1", "logline2"***REMOVED***, logLines(hook))
	***REMOVED***)

	t.Run("RestoreConnFromLatestMessage", func(t *testing.T) ***REMOVED***
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		startFilter := func(u url.URL) (start time.Time, err error) ***REMOVED***
			rawstart, err := strconv.ParseInt(u.Query().Get("start"), 10, 64)
			if err != nil ***REMOVED***
				return start, err
			***REMOVED***

			start = time.Unix(0, rawstart)
			return
		***REMOVED***

		var requestsCount uint64

		tb := httpmultibin.NewHTTPMultiBin(t)
		logtailHandleFunc(tb, func(conn *websocket.Conn, req *http.Request) ***REMOVED***
			requests := atomic.AddUint64(&requestsCount, 1)

			start, err := startFilter(*req.URL)
			require.NoError(t, err)

			if requests <= 1 ***REMOVED***
				t0 := time.Date(2021, time.July, 27, 0, 0, 0, 0, time.UTC).UnixNano()
				t1 := time.Date(2021, time.July, 27, 1, 0, 0, 0, time.UTC).UnixNano()
				t2 := time.Date(2021, time.July, 27, 2, 0, 0, 0, time.UTC).UnixNano()

				// send a correct logline so we will able to assert
				// that the connection is restored from t2 as expected
				rawmsg := json.RawMessage(fmt.Sprintf(`***REMOVED***"streams":[***REMOVED***"stream":***REMOVED***"key":"stream1","level":"warn"***REMOVED***,"values":[["%d","newest logline"],["%d","second logline"],["%d","oldest logline"]]***REMOVED***],"dropped_entities":[]***REMOVED***`, t2, t1, t0))
				err = conn.WriteJSON(rawmsg)
				require.NoError(t, err)

				// wait the flush of the message on the network
				time.Sleep(20 * time.Millisecond)

				// it generates a failure closing the connection
				// in a rude way
				err = conn.Close()
				require.NoError(t, err)
				return
			***REMOVED***

			// assert that the client created the request with `start`
			// populated from the most recent seen value (t2+1ns)
			require.Equal(t, time.Unix(0, 1627351200000000001), start)

			// send a correct logline so we will able to assert
			// that the connection is restored as expected
			err = conn.WriteJSON(json.RawMessage(generateLogline("stream3", 1627358400000000000, "logline-after-restored-conn")))
			require.NoError(t, err)

			// wait the flush of the message on the network
			time.Sleep(20 * time.Millisecond)
			cancel()
		***REMOVED***)

		logger := logrus.New()
		logger.Out = ioutil.Discard
		hook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: logrus.AllLevels***REMOVED***
		logger.AddHook(hook)

		c := configFromHTTPMultiBin(tb)
		err := c.StreamLogsToLogger(ctx, logger, "ref_id", 0)
		require.NoError(t, err)

		assert.Equal(t,
			[]string***REMOVED***
				"newest logline",
				"second logline",
				"oldest logline",
				"error reading a log message from the cloud, trying to establish a fresh connection with the logs service...",
				"logline-after-restored-conn",
			***REMOVED***, logLines(hook))
	***REMOVED***)

	t.Run("RestoreConnFromTimeNow", func(t *testing.T) ***REMOVED***
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		startFilter := func(u url.URL) (start time.Time, err error) ***REMOVED***
			rawstart, err := strconv.ParseInt(u.Query().Get("start"), 10, 64)
			if err != nil ***REMOVED***
				return start, err
			***REMOVED***

			start = time.Unix(0, rawstart)
			return
		***REMOVED***

		var requestsCount uint64
		t0 := time.Now()

		tb := httpmultibin.NewHTTPMultiBin(t)
		logtailHandleFunc(tb, func(conn *websocket.Conn, req *http.Request) ***REMOVED***
			requests := atomic.AddUint64(&requestsCount, 1)

			start, err := startFilter(*req.URL)
			require.NoError(t, err)

			if requests <= 1 ***REMOVED***
				// if it's the first attempt then
				// it generates a failure closing the connection
				// in a rude way
				err = conn.Close()
				require.NoError(t, err)
				return
			***REMOVED***

			// it asserts that the second attempt
			// has a `start` after the test run
			require.True(t, start.After(t0))

			// send a correct logline so we will able to assert
			// that the connection is restored as expected
			err = conn.WriteJSON(json.RawMessage(`***REMOVED***"streams":[***REMOVED***"stream":***REMOVED***"key":"stream1","level":"warn"***REMOVED***,"values":[["1598282752000000000","logline-after-restored-conn"]]***REMOVED***],"dropped_entities":[]***REMOVED***`))
			require.NoError(t, err)

			// wait the flush of the message on the network
			time.Sleep(20 * time.Millisecond)
			cancel()
		***REMOVED***)

		logger := logrus.New()
		logger.Out = ioutil.Discard
		hook := &testutils.SimpleLogrusHook***REMOVED***HookedLevels: logrus.AllLevels***REMOVED***
		logger.AddHook(hook)

		c := configFromHTTPMultiBin(tb)
		err := c.StreamLogsToLogger(ctx, logger, "ref_id", 0)
		require.NoError(t, err)

		assert.Equal(t,
			[]string***REMOVED***
				"error reading a log message from the cloud, trying to establish a fresh connection with the logs service...",
				"logline-after-restored-conn",
			***REMOVED***, logLines(hook))
	***REMOVED***)
***REMOVED***
