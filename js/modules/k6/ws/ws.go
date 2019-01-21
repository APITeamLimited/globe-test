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
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
)

// ErrWSInInitContext is returned when websockets are using in the init context
var ErrWSInInitContext = common.NewInitContextError("using websockets in the init context is not supported")

type WS struct***REMOVED******REMOVED***

type pingDelta struct ***REMOVED***
	ping time.Time
	pong time.Time
***REMOVED***

type Socket struct ***REMOVED***
	ctx           context.Context
	conn          *websocket.Conn
	eventHandlers map[string][]goja.Callable
	scheduled     chan goja.Callable
	done          chan struct***REMOVED******REMOVED***
	shutdownOnce  sync.Once

	msgSentTimestamps     []time.Time
	msgReceivedTimestamps []time.Time

	pingSendTimestamps map[string]time.Time
	pingSendCounter    int
	pingTimestamps     []pingDelta
***REMOVED***

type WSHTTPResponse struct ***REMOVED***
	URL     string            `json:"url"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Error   string            `json:"error"`
***REMOVED***

const writeWait = 10 * time.Second

func New() *WS ***REMOVED***
	return &WS***REMOVED******REMOVED***
***REMOVED***

func (*WS) Connect(ctx context.Context, url string, args ...goja.Value) (*WSHTTPResponse, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)
	if state == nil ***REMOVED***
		return nil, ErrWSInInitContext
	***REMOVED***

	// The params argument is optional
	var callableV, paramsV goja.Value
	switch len(args) ***REMOVED***
	case 2:
		paramsV = args[0]
		callableV = args[1]
	case 1:
		paramsV = goja.Undefined()
		callableV = args[0]
	default:
		return nil, errors.New("invalid number of arguments to ws.connect")
	***REMOVED***

	// Get the callable (required)
	setupFn, isFunc := goja.AssertFunction(callableV)
	if !isFunc ***REMOVED***
		return nil, errors.New("last argument to ws.connect must be a function")
	***REMOVED***

	// Leave header to nil by default so we can pass it directly to the Dialer
	var header http.Header

	tags := state.Options.RunTags.CloneTags()

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

	if state.Options.SystemTags["url"] ***REMOVED***
		tags["url"] = url
	***REMOVED***
	if state.Options.SystemTags["group"] ***REMOVED***
		tags["group"] = state.Group.Path
	***REMOVED***

	// Pass a custom net.Dial function to websocket.Dialer that will substitute
	// the underlying net.Conn with our own tracked netext.Conn
	netDial := func(network, address string) (net.Conn, error) ***REMOVED***
		return state.Dialer.DialContext(ctx, network, address)
	***REMOVED***

	// Overriding the NextProtos to avoid talking http2
	var tlsConfig *tls.Config
	if state.TLSConfig != nil ***REMOVED***
		tlsConfig = state.TLSConfig.Clone()
		tlsConfig.NextProtos = []string***REMOVED***"http/1.1"***REMOVED***
	***REMOVED***

	wsd := websocket.Dialer***REMOVED***
		NetDial:         netDial,
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	***REMOVED***

	start := time.Now()
	conn, httpResponse, connErr := wsd.Dial(url, header)
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

	if state.Options.SystemTags["ip"] && conn.RemoteAddr() != nil ***REMOVED***
		if ip, _, err := net.SplitHostPort(conn.RemoteAddr().String()); err == nil ***REMOVED***
			tags["ip"] = ip
		***REMOVED***
	***REMOVED***

	// Run the user-provided set up function
	if _, err := setupFn(goja.Undefined(), rt.ToValue(&socket)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if connErr != nil ***REMOVED***
		// Pass the error to the user script before exiting immediately
		socket.handleEvent("error", rt.ToValue(connErr))

		return nil, connErr
	***REMOVED***

	wsResponse, wsRespErr := wrapHTTPResponse(httpResponse)
	if wsRespErr != nil ***REMOVED***
		return nil, wsRespErr
	***REMOVED***
	wsResponse.URL = url

	defer func() ***REMOVED*** _ = conn.Close() ***REMOVED***()

	if state.Options.SystemTags["status"] ***REMOVED***
		tags["status"] = strconv.Itoa(httpResponse.StatusCode)
	***REMOVED***
	if state.Options.SystemTags["subproto"] ***REMOVED***
		tags["subproto"] = httpResponse.Header.Get("Sec-WebSocket-Protocol")
	***REMOVED***

	// The connection is now open, emit the event
	socket.handleEvent("open")

	// Pass ping/pong events through the main control loop
	pingChan := make(chan string)
	pongChan := make(chan string)
	conn.SetPingHandler(func(msg string) error ***REMOVED*** pingChan <- msg; return nil ***REMOVED***)
	conn.SetPongHandler(func(pingID string) error ***REMOVED*** pongChan <- pingID; return nil ***REMOVED***)

	readDataChan := make(chan []byte)
	readCloseChan := make(chan int)
	readErrChan := make(chan error)

	// Wraps a couple of channels around conn.ReadMessage
	go readPump(conn, readDataChan, readErrChan, readCloseChan)

	// This is the main control loop. All JS code (including error handlers)
	// should only be executed by this thread to avoid race conditions
	for ***REMOVED***
		select ***REMOVED***
		case pingData := <-pingChan:
			// Handle pings received from the server
			// - trigger the `ping` event
			// - reply with pong (needed when `SetPingHandler` is overwritten)
			err := socket.conn.WriteControl(websocket.PongMessage, []byte(pingData), time.Now().Add(writeWait))
			if err != nil ***REMOVED***
				socket.handleEvent("error", rt.ToValue(err))
			***REMOVED***
			socket.handleEvent("ping")

		case pingID := <-pongChan:
			// Handle pong responses to our pings
			socket.trackPong(pingID)
			socket.handleEvent("pong")

		case readData := <-readDataChan:
			socket.msgReceivedTimestamps = append(socket.msgReceivedTimestamps, time.Now())
			socket.handleEvent("message", rt.ToValue(string(readData)))

		case readErr := <-readErrChan:
			socket.handleEvent("error", rt.ToValue(readErr))

		case readClose := <-readCloseChan:
			// handle server close
			socket.handleEvent("close", rt.ToValue(readClose))

		case scheduledFn := <-socket.scheduled:
			if _, err := scheduledFn(goja.Undefined()); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

		case <-ctx.Done():
			// VU is shutting down during an interrupt
			// socket events will not be forwarded to the VU
			_ = socket.closeConnection(websocket.CloseGoingAway)

		case <-socket.done:
			// This is the final exit point normally triggered by closeConnection
			end := time.Now()
			sessionDuration := stats.D(end.Sub(start))

			sampleTags := stats.IntoSampleTags(&tags)

			stats.PushIfNotCancelled(ctx, state.Samples, stats.ConnectedSamples***REMOVED***
				Samples: []stats.Sample***REMOVED***
					***REMOVED***Metric: metrics.WSSessions, Time: start, Tags: sampleTags, Value: 1***REMOVED***,
					***REMOVED***Metric: metrics.WSConnecting, Time: start, Tags: sampleTags, Value: connectionDuration***REMOVED***,
					***REMOVED***Metric: metrics.WSSessionDuration, Time: start, Tags: sampleTags, Value: sessionDuration***REMOVED***,
				***REMOVED***,
				Tags: sampleTags,
				Time: start,
			***REMOVED***)

			for _, msgSentTimestamp := range socket.msgSentTimestamps ***REMOVED***
				stats.PushIfNotCancelled(ctx, state.Samples, stats.Sample***REMOVED***
					Metric: metrics.WSMessagesSent,
					Time:   msgSentTimestamp,
					Tags:   sampleTags,
					Value:  1,
				***REMOVED***)
			***REMOVED***

			for _, msgReceivedTimestamp := range socket.msgReceivedTimestamps ***REMOVED***
				stats.PushIfNotCancelled(ctx, state.Samples, stats.Sample***REMOVED***
					Metric: metrics.WSMessagesReceived,
					Time:   msgReceivedTimestamp,
					Tags:   sampleTags,
					Value:  1,
				***REMOVED***)
			***REMOVED***

			for _, pingDelta := range socket.pingTimestamps ***REMOVED***
				stats.PushIfNotCancelled(ctx, state.Samples, stats.Sample***REMOVED***
					Metric: metrics.WSPing,
					Time:   pingDelta.pong,
					Tags:   sampleTags,
					Value:  stats.D(pingDelta.pong.Sub(pingDelta.ping)),
				***REMOVED***)
			***REMOVED***

			return wsResponse, nil
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
			if _, err := handler(goja.Undefined(), args...); err != nil ***REMOVED***
				common.Throw(common.GetRuntime(s.ctx), err)
			***REMOVED***
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

	s.msgSentTimestamps = append(s.msgSentTimestamps, time.Now())
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

	s.pingTimestamps = append(s.pingTimestamps, pingDelta***REMOVED***pingTimestamp, pongTimestamp***REMOVED***)
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

	_ = s.closeConnection(code)
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

		// trigger `close` event when the client closes the connection
		s.handleEvent("close", rt.ToValue(code))
		_ = s.conn.Close()

		// Stops the main control loop
		close(s.done)
	***REMOVED***)

	return err
***REMOVED***

// Wraps conn.ReadMessage in a channel
func readPump(conn *websocket.Conn, readChan chan []byte, errorChan chan error, closeChan chan int) ***REMOVED***
	defer func() ***REMOVED*** _ = conn.Close() ***REMOVED***()

	for ***REMOVED***
		_, message, err := conn.ReadMessage()
		if err != nil ***REMOVED***

			if websocket.IsCloseError(err, websocket.CloseNormalClosure) ***REMOVED***
				closeChan <- err.(*websocket.CloseError).Code
			***REMOVED*** else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) ***REMOVED***
				// Emit the error if it is not CloseNormalClosure
				// and the error is not  originated from closing the socket ourselves with `CloseGoingAway`
				errorChan <- err
			***REMOVED***

			//CloseGoingAway errors are ignored
			return
		***REMOVED***

		readChan <- message
	***REMOVED***
***REMOVED***

// Wrap the raw HTTPResponse we received to a WSHTTPResponse we can pass to the user
func wrapHTTPResponse(httpResponse *http.Response) (*WSHTTPResponse, error) ***REMOVED***
	wsResponse := WSHTTPResponse***REMOVED***
		Status: httpResponse.StatusCode,
	***REMOVED***

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = httpResponse.Body.Close()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	wsResponse.Body = string(body)

	wsResponse.Headers = make(map[string]string, len(httpResponse.Header))
	for k, vs := range httpResponse.Header ***REMOVED***
		wsResponse.Headers[k] = strings.Join(vs, ", ")
	***REMOVED***

	return &wsResponse, nil
***REMOVED***
