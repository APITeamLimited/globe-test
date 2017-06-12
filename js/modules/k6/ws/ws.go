/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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

package ws

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
)

type WS struct***REMOVED******REMOVED***

type Socket struct ***REMOVED***
	ctx           context.Context
	conn          *websocket.Conn
	eventHandlers map[string][]goja.Callable
	scheduled     chan goja.Callable
	done          chan struct***REMOVED******REMOVED***
	shutdownOnce  sync.Once

	msgsSent, msgsReceived float64

	pingSendTimestamps map[string]time.Time
	pingSendCounter    int
	pingDeltas         []time.Duration
***REMOVED***

const writeWait = 10 * time.Second

var (
	newline = []byte***REMOVED***'\n'***REMOVED***
	space   = []byte***REMOVED***' '***REMOVED***
)

func (*WS) Connect(ctx context.Context, url string, args ...goja.Value) *http.Response ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	var callableV goja.Value
	var paramsV goja.Value

	// The params argument is optional
	switch len(args) ***REMOVED***
	case 2:
		paramsV = args[0]
		callableV = args[1]
	case 1:
		paramsV = goja.Undefined()
		callableV = args[0]
	default:
		common.Throw(rt, errors.New("Invalid number of arguments to ws.connect"))
		return nil
	***REMOVED***

	// Get the callable (required)
	setupFn, isFunc := goja.AssertFunction(callableV)
	if !isFunc ***REMOVED***
		common.Throw(rt, errors.New("Last argument to ws.connect must be a function"))
		return nil
	***REMOVED***

	// Leave header to nil by default so we can pass it directly to the Dialer
	var header http.Header

	tags := map[string]string***REMOVED***
		"url":         url,
		"group":       state.Group.Path,
		"status":      "0",
		"subprotocol": "",
	***REMOVED***

	// Parse the optional second argument (params)
	if !goja.IsUndefined(paramsV) && !goja.IsNull(paramsV) ***REMOVED***
		params := paramsV.ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "headers":
				header = http.Header***REMOVED******REMOVED***
				headersV := params.Get(k)
				if goja.IsUndefined(headersV) || goja.IsNull(headersV) ***REMOVED***
					continue
				***REMOVED***
				headersObj := headersV.ToObject(rt)
				if headersObj == nil ***REMOVED***
					continue
				***REMOVED***
				for _, key := range headersObj.Keys() ***REMOVED***
					header.Set(key, headersObj.Get(key).String())
				***REMOVED***
			case "tags":
				tagsV := params.Get(k)
				if goja.IsUndefined(tagsV) || goja.IsNull(tagsV) ***REMOVED***
					continue
				***REMOVED***
				tagObj := tagsV.ToObject(rt)
				if tagObj == nil ***REMOVED***
					continue
				***REMOVED***
				for _, key := range tagObj.Keys() ***REMOVED***
					tags[key] = tagObj.Get(key).String()
				***REMOVED***
			***REMOVED***
		***REMOVED***

	***REMOVED***

	var bytesRead, bytesWritten int64
	// Pass a custom net.Dial function to websocket.Dialer that will substitute
	// the underlying net.Conn with our own tracked netext.Conn
	netDial := func(network, address string) (net.Conn, error) ***REMOVED***
		conn, err := state.Dialer.DialContext(ctx, network, address)
		if err != nil ***REMOVED***
			return conn, err
		***REMOVED***

		return netext.TrackConn(conn, &bytesRead, &bytesWritten), nil
	***REMOVED***

	wsd := websocket.Dialer***REMOVED***
		NetDial: netDial,
		Proxy:   http.ProxyFromEnvironment,
	***REMOVED***

	start := time.Now()
	conn, response, connErr := wsd.Dial(url, header)
	connectionEnd := time.Now()
	connectionDuration := stats.D(connectionEnd.Sub(start))

	socket := Socket***REMOVED***
		ctx:                ctx,
		conn:               conn,
		eventHandlers:      make(map[string][]goja.Callable),
		pingSendTimestamps: make(map[string]time.Time),
		scheduled:          make(chan goja.Callable),
		done:               make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// Run the user-provided set up function
	setupFn(goja.Undefined(), rt.ToValue(&socket))

	if connErr != nil ***REMOVED***
		// Pass the error to the user script before exiting immediately
		socket.handleEvent("error", rt.ToValue(connErr))
		return response
	***REMOVED***
	defer conn.Close()

	tags["status"] = strconv.Itoa(response.StatusCode)
	tags["subprotocol"] = response.Header.Get("Sec-WebSocket-Protocol")

	// The connection is now open, emit the event
	socket.handleEvent("open")

	// Pass ping/pong events through the main control loop
	pingChan := make(chan string)
	pongChan := make(chan string)
	conn.SetPingHandler(func(msg string) error ***REMOVED*** pingChan <- msg; return nil ***REMOVED***)
	conn.SetPongHandler(func(pingID string) error ***REMOVED*** pongChan <- pingID; return nil ***REMOVED***)

	readDataChan := make(chan []byte)
	readErrChan := make(chan error)

	// Wraps a couple of channels around conn.ReadMessage
	go readPump(conn, readDataChan, readErrChan)

	// This is the main control loop. All JS code (including error handlers)
	// should only be executed by this thread to avoid race conditions
	for ***REMOVED***
		select ***REMOVED***
		case <-pingChan:
			// Handle pings received from the server
			socket.handleEvent("ping")

		case pingID := <-pongChan:
			// Handle pong responses to our pings
			socket.trackPong(pingID)
			socket.handleEvent("pong")

		case readData := <-readDataChan:
			socket.msgsReceived++
			socket.handleEvent("message", rt.ToValue(string(readData)))

		case readErr := <-readErrChan:
			socket.handleEvent("error", rt.ToValue(readErr))

		case scheduledFn := <-socket.scheduled:
			scheduledFn(goja.Undefined())

		case <-ctx.Done():
			// This means that K6 is shutting down (e.g., during an interrupt)
			socket.handleEvent("close", rt.ToValue("Interrupt"))
			socket.closeConnection(websocket.CloseGoingAway)

		case <-socket.done:
			// This is the final exit point normally triggered by closeConnection
			end := time.Now()
			sessionDuration := stats.D(end.Sub(start))

			samples := []stats.Sample***REMOVED***
				***REMOVED***Metric: metrics.WSSessions, Time: end, Tags: tags, Value: 1***REMOVED***,
				***REMOVED***Metric: metrics.WSMessagesSent, Time: end, Tags: tags, Value: socket.msgsReceived***REMOVED***,
				***REMOVED***Metric: metrics.WSMessagesReceived, Time: end, Tags: tags, Value: socket.msgsSent***REMOVED***,
				***REMOVED***Metric: metrics.WSConnecting, Time: end, Tags: tags, Value: connectionDuration***REMOVED***,
				***REMOVED***Metric: metrics.WSSessionDuration, Time: end, Tags: tags, Value: sessionDuration***REMOVED***,
				***REMOVED***Metric: metrics.DataReceived, Time: end, Tags: tags, Value: float64(bytesRead)***REMOVED***,
				***REMOVED***Metric: metrics.DataSent, Time: end, Tags: tags, Value: float64(bytesWritten)***REMOVED***,
			***REMOVED***

			for _, pingDelta := range socket.pingDeltas ***REMOVED***
				samples = append(samples, stats.Sample***REMOVED***
					Metric: metrics.WSPing,
					Time:   end,
					Tags:   tags,
					Value:  stats.D(pingDelta),
				***REMOVED***)
			***REMOVED***

			state.Samples = append(state.Samples, samples...)

			return response
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *Socket) On(event string, handler goja.Value) ***REMOVED***
	if handler, ok := goja.AssertFunction(handler); ok ***REMOVED***
		s.eventHandlers[event] = append(s.eventHandlers[event], handler)
	***REMOVED***
***REMOVED***

func (s *Socket) handleEvent(event string, args ...goja.Value) ***REMOVED***
	if handlers, ok := s.eventHandlers[event]; ok ***REMOVED***
		for _, handler := range handlers ***REMOVED***
			handler(goja.Undefined(), args...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *Socket) Send(message string) ***REMOVED***
	// NOTE: No binary message support for the time being since goja doesn't
	// support typed arrays.
	rt := common.GetRuntime(s.ctx)

	writeData := []byte(message)
	if err := s.conn.WriteMessage(websocket.TextMessage, writeData); err != nil ***REMOVED***
		s.handleEvent("error", rt.ToValue(err))
	***REMOVED***

	s.msgsSent++
***REMOVED***

func (s *Socket) Ping() ***REMOVED***
	rt := common.GetRuntime(s.ctx)
	deadline := time.Now().Add(writeWait)
	pingID := strconv.Itoa(s.pingSendCounter)
	data := []byte(pingID)

	err := s.conn.WriteControl(websocket.PingMessage, data, deadline)
	if err != nil ***REMOVED***
		s.handleEvent("error", rt.ToValue(err))
		return
	***REMOVED***

	s.pingSendTimestamps[pingID] = time.Now()
	s.pingSendCounter++
***REMOVED***

func (s *Socket) trackPong(pingID string) ***REMOVED***
	pongTimestamp := time.Now()

	if _, ok := s.pingSendTimestamps[pingID]; !ok ***REMOVED***
		// We received a pong for a ping we didn't send; ignore
		// (this shouldn't happen with a compliant server)
		return
	***REMOVED***
	pingTimestamp := s.pingSendTimestamps[pingID]

	s.pingDeltas = append(s.pingDeltas, pongTimestamp.Sub(pingTimestamp))
***REMOVED***

func (s *Socket) SetTimeout(fn goja.Callable, timeoutMs int) ***REMOVED***
	// Starts a goroutine, blocks once on the timeout and pushes the callable
	// back to the main loop through the scheduled channel
	go func() ***REMOVED***
		select ***REMOVED***
		case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
			s.scheduled <- fn

		case <-s.done:
			return
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (s *Socket) SetInterval(fn goja.Callable, intervalMs int) ***REMOVED***
	// Starts a goroutine, blocks forever on the ticker and pushes the callable
	// back to the main loop through the scheduled channel
	go func() ***REMOVED***
		ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
		defer ticker.Stop()

		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				s.scheduled <- fn

			case <-s.done:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (s *Socket) Close(args ...goja.Value) ***REMOVED***
	code := websocket.CloseGoingAway
	if len(args) > 0 ***REMOVED***
		code = int(args[0].ToInteger())
	***REMOVED***

	s.closeConnection(code)
***REMOVED***

// Attempts to close the websocket gracefully
func (s *Socket) closeConnection(code int) error ***REMOVED***
	var err error

	s.shutdownOnce.Do(func() ***REMOVED***
		rt := common.GetRuntime(s.ctx)

		writeErr := s.conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(code, ""),
			time.Now().Add(writeWait),
		)
		if writeErr != nil ***REMOVED***
			// Just call the handler, we'll try to close the connection anyway
			s.handleEvent("error", rt.ToValue(err))
			err = writeErr
		***REMOVED***
		s.conn.Close()

		// Stops the main control loop
		close(s.done)
	***REMOVED***)

	return err
***REMOVED***

// Wraps conn.ReadMessage in a channel
func readPump(conn *websocket.Conn, readChan chan []byte, errorChan chan error) ***REMOVED***
	defer conn.Close()

	for ***REMOVED***
		_, message, err := conn.ReadMessage()
		if err != nil ***REMOVED***
			// Only emit the error if we didn't close the socket ourselves
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) ***REMOVED***
				errorChan <- err
			***REMOVED***

			return
		***REMOVED***

		readChan <- message
	***REMOVED***
***REMOVED***
