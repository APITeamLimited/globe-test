package cloudapi

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

//go:generate easyjson -pkg -no_std_marshalers -gen_build_flags -mod=mod .

//easyjson:json
type msg struct ***REMOVED***
	Streams        []msgStreams        `json:"streams"`
	DroppedEntries []msgDroppedEntries `json:"dropped_entries"`
***REMOVED***

//easyjson:json
type msgStreams struct ***REMOVED***
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"` // this can be optimized
***REMOVED***

//easyjson:json
type msgDroppedEntries struct ***REMOVED***
	Labels    map[string]string `json:"labels"`
	Timestamp string            `json:"timestamp"`
***REMOVED***

// Log writes the Streams and Dropped Entries to the passed logger.
// It returns the most recent timestamp seen overall messages.
func (m *msg) Log(logger logrus.FieldLogger) int64 ***REMOVED***
	var level string
	var ts int64

	for _, stream := range m.Streams ***REMOVED***
		fields := labelsToLogrusFields(stream.Stream)
		var ok bool
		if level, ok = stream.Stream["level"]; ok ***REMOVED***
			delete(fields, "level")
		***REMOVED***

		for _, value := range stream.Values ***REMOVED***
			nsec, _ := strconv.ParseInt(value[0], 10, 64)
			e := logger.WithFields(fields).WithTime(time.Unix(0, nsec))
			lvl, err := logrus.ParseLevel(level)
			if err != nil ***REMOVED***
				e.Info(value[1])
				e.Warn("last message had unknown level " + level)
			***REMOVED*** else ***REMOVED***
				e.Log(lvl, value[1])
			***REMOVED***

			// find the latest seen message
			if nsec > ts ***REMOVED***
				ts = nsec
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, dropped := range m.DroppedEntries ***REMOVED***
		nsec, _ := strconv.ParseInt(dropped.Timestamp, 10, 64)
		logger.WithFields(labelsToLogrusFields(dropped.Labels)).WithTime(time.Unix(0, nsec)).Warn("dropped")

		if nsec > ts ***REMOVED***
			ts = nsec
		***REMOVED***
	***REMOVED***

	return ts
***REMOVED***

func labelsToLogrusFields(labels map[string]string) logrus.Fields ***REMOVED***
	fields := make(logrus.Fields, len(labels))

	for key, val := range labels ***REMOVED***
		fields[key] = val
	***REMOVED***

	return fields
***REMOVED***

func (c *Config) logtailConn(ctx context.Context, referenceID string, since time.Time) (*websocket.Conn, error) ***REMOVED***
	u, err := url.Parse(c.LogsTailURL.String)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("couldn't parse cloud logs host %w", err)
	***REMOVED***

	u.RawQuery = fmt.Sprintf(`query=***REMOVED***test_run_id="%s"***REMOVED***&start=%d`, referenceID, since.UnixNano())

	headers := make(http.Header)
	headers.Add("Sec-WebSocket-Protocol", "token="+c.Token.String)

	var conn *websocket.Conn
	err = retry(sleeperFunc(time.Sleep), 3, 5*time.Second, 2*time.Minute, func() (err error) ***REMOVED***
		// We don't need to close the http body or use it for anything until we want to actually log
		// what the server returned as body when it errors out
		conn, _, err = websocket.DefaultDialer.DialContext(ctx, u.String(), headers) //nolint:bodyclose
		return err
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return conn, nil
***REMOVED***

// StreamLogsToLogger streams the logs for the configured test to the provided logger until ctx is
// Done or an error occurs.
func (c *Config) StreamLogsToLogger(
	ctx context.Context, logger logrus.FieldLogger, referenceID string, tailFrom time.Duration,
) error ***REMOVED***
	var mconn sync.Mutex

	conn, err := c.logtailConn(ctx, referenceID, time.Now().Add(-tailFrom))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	go func() ***REMOVED***
		<-ctx.Done()

		mconn.Lock()
		defer mconn.Unlock()

		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "closing"),
			time.Now().Add(time.Second))

		_ = conn.Close()
	***REMOVED***()

	msgBuffer := make(chan []byte, 10)
	defer close(msgBuffer)

	var mostRecent int64
	go func() ***REMOVED***
		for message := range msgBuffer ***REMOVED***
			var m msg
			err := easyjson.Unmarshal(message, &m)
			if err != nil ***REMOVED***
				logger.WithError(err).Errorf("couldn't unmarshal a message from the cloud: %s", string(message))

				continue
			***REMOVED***
			ts := m.Log(logger)
			atomic.StoreInt64(&mostRecent, ts)
		***REMOVED***
	***REMOVED***()

	for ***REMOVED***
		_, message, err := conn.ReadMessage()
		select ***REMOVED*** // check if we should stop before continuing
		case <-ctx.Done():
			return nil
		default:
		***REMOVED***

		if err != nil ***REMOVED***
			logger.WithError(err).Warn("error reading a log message from the cloud, trying to establish a fresh connection with the logs service...") //nolint:lll

			var since time.Time
			if ts := atomic.LoadInt64(&mostRecent); ts > 0 ***REMOVED***
				// add 1ns for avoid possible repetition
				since = time.Unix(0, ts).Add(time.Nanosecond)
			***REMOVED*** else ***REMOVED***
				since = time.Now()
			***REMOVED***

			// TODO: avoid the "logical" race condition
			// The case explained:
			// * The msgBuffer consumer is slow
			// * ReadMessage is fast and adds at least one more message in the buffer
			// * An error is got in the meantime and the re-dialing procedure is tried
			// * Then the latest timestamp used will not be the real latest received
			// * because it is still waiting to be processed.
			// In the case the connection will be restored then the first message will be a duplicate.
			newconn, errd := c.logtailConn(ctx, referenceID, since)
			if errd != nil ***REMOVED***
				// return the main error
				return err
			***REMOVED***

			mconn.Lock()
			conn = newconn
			mconn.Unlock()
			continue
		***REMOVED***

		select ***REMOVED***
		case <-ctx.Done():
			return nil
		case msgBuffer <- message:
		***REMOVED***
	***REMOVED***
***REMOVED***

// sleeper represents an abstraction for waiting an amount of time.
type sleeper interface ***REMOVED***
	Sleep(d time.Duration)
***REMOVED***

// sleeperFunc uses the underhood function for implementing the wait operation.
type sleeperFunc func(time.Duration)

func (sfn sleeperFunc) Sleep(d time.Duration) ***REMOVED***
	sfn(d)
***REMOVED***

// retry retries to execute a provided function until it isn't successful
// or the maximum number of attempts is hit. It waits the specified interval
// between the latest iteration and the next retry.
// Interval is used as the base to compute an exponential backoff,
// if the computed interval overtakes the max interval then max will be used.
func retry(s sleeper, attempts uint, interval, max time.Duration, do func() error) (err error) ***REMOVED***
	baseInterval := math.Abs(interval.Truncate(time.Second).Seconds())
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec

	for i := 0; i < int(attempts); i++ ***REMOVED***
		if i > 0 ***REMOVED***
			// wait = (interval ^ i) + random milliseconds
			wait := time.Duration(math.Pow(baseInterval, float64(i))) * time.Second
			wait += time.Duration(r.Int63n(1000)) * time.Millisecond

			if wait > max ***REMOVED***
				wait = max
			***REMOVED***
			s.Sleep(wait)
		***REMOVED***
		err = do()
		if err == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
