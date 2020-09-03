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
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

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

func (m *msg) Log(logger logrus.FieldLogger) ***REMOVED***
	var level string

	for _, stream := range m.Streams ***REMOVED***
		fields := labelsToLogrusFields(stream.Stream)
		var ok bool
		if level, ok = stream.Stream["level"]; ok ***REMOVED***
			delete(fields, "level")
		***REMOVED***

		for _, value := range stream.Values ***REMOVED***
			nsec, _ := strconv.Atoi(value[0])
			e := logger.WithFields(fields).WithTime(time.Unix(0, int64(nsec)))
			lvl, err := logrus.ParseLevel(level)
			if err != nil ***REMOVED***
				e.Info(value[1])
				e.Warn("last message had unknown level " + level)
			***REMOVED*** else ***REMOVED***
				e.Log(lvl, value[1])
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, dropped := range m.DroppedEntries ***REMOVED***
		nsec, _ := strconv.Atoi(dropped.Timestamp)
		logger.WithFields(labelsToLogrusFields(dropped.Labels)).WithTime(time.Unix(0, int64(nsec))).Warn("dropped")
	***REMOVED***
***REMOVED***

func labelsToLogrusFields(labels map[string]string) logrus.Fields ***REMOVED***
	fields := make(logrus.Fields, len(labels))

	for key, val := range labels ***REMOVED***
		fields[key] = val
	***REMOVED***

	return fields
***REMOVED***

func (c *Config) getRequest(referenceID string, start time.Duration) (*url.URL, error) ***REMOVED***
	u, err := url.Parse(c.LogsHost.String)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("couldn't parse cloud logs host %w", err)
	***REMOVED***

	u.RawQuery = fmt.Sprintf(`query=***REMOVED***test_run_id="%s"***REMOVED***&start=%d`,
		referenceID,
		time.Now().Add(-start).UnixNano(),
	)

	return u, nil
***REMOVED***

// StreamLogsToLogger streams the logs for the configured test to the provided logger until ctx is
// Done or an error occurs.
func (c *Config) StreamLogsToLogger(
	ctx context.Context, logger logrus.FieldLogger, referenceID string, start time.Duration,
) error ***REMOVED***
	u, err := c.getRequest(referenceID, start)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	headers := make(http.Header)
	headers.Add("Sec-WebSocket-Protocol", "token="+c.Token.String)

	// We don't need to close the http body or use it for anything until we want to actually log
	// what the server returned as body when it errors out
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), headers) //nolint:bodyclose
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	go func() ***REMOVED***
		<-ctx.Done()

		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "closing"),
			time.Now().Add(time.Second))

		_ = conn.Close()
	***REMOVED***()

	msgBuffer := make(chan []byte, 10)

	defer close(msgBuffer)

	go func() ***REMOVED***
		for message := range msgBuffer ***REMOVED***
			var m msg
			err := easyjson.Unmarshal(message, &m)
			if err != nil ***REMOVED***
				logger.WithError(err).Errorf("couldn't unmarshal a message from the cloud: %s", string(message))

				continue
			***REMOVED***

			m.Log(logger)
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
			logger.WithError(err).Warn("error reading a message from the cloud")

			return err
		***REMOVED***

		select ***REMOVED***
		case <-ctx.Done():
			return nil
		case msgBuffer <- message:
		***REMOVED***
	***REMOVED***
***REMOVED***
