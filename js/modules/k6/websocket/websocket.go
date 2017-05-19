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

package websocket

import (
	"bytes"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"
	"github.com/loadimpact/k6/js/common"
)

type Websocket struct***REMOVED******REMOVED***

type Socket struct ***REMOVED***
	ctx           context.Context
	conn          *websocket.Conn
	eventHandlers map[string]goja.Callable
	scheduled     chan goja.Callable
	done          chan struct***REMOVED******REMOVED***
***REMOVED***

const writeWait = 10 * time.Second

var (
	newline = []byte***REMOVED***'\n'***REMOVED***
	space   = []byte***REMOVED***' '***REMOVED***
)

func (*Websocket) Connect(ctx context.Context, url string, setupFn goja.Callable) ***REMOVED***
	rt := common.GetRuntime(ctx)
	conn, _, connErr := websocket.DefaultDialer.Dial(url, nil)

	socket := Socket***REMOVED***
		ctx:           ctx,
		conn:          conn,
		eventHandlers: make(map[string]goja.Callable),
		scheduled:     make(chan goja.Callable),
		done:          make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// Run the user-provided set up function
	setupFn(goja.Undefined(), rt.ToValue(&socket))

	if connErr != nil ***REMOVED***
		// Pass the error to the user script before exiting immediately
		socket.handleEvent("error", rt.ToValue(connErr))
		return
	***REMOVED***
	defer conn.Close()

	conn.SetPongHandler(func(string) error ***REMOVED*** socket.handleEvent("pong"); return nil ***REMOVED***)

	// The connection is now open, emit the event
	socket.handleEvent("open")

	readDataChan := make(chan []byte)
	readErrChan := make(chan error)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wraps a couple of channels around conn.ReadMessage
	go readPump(conn, readDataChan, readErrChan)

	// This is the main control loop. All JS code (including error handlers)
	// should only be executed by this thread to avoid race conditions
	for ***REMOVED***
		select ***REMOVED***
		case readData := <-readDataChan:
			socket.handleEvent("message", rt.ToValue(string(readData)))
		case readErr := <-readErrChan:
			socket.handleEvent("error", rt.ToValue(readErr))

		case scheduledFn := <-socket.scheduled:
			scheduledFn(goja.Undefined())

		case <-interrupt:
			socket.handleEvent("close", rt.ToValue("Interrupt"))
			socket.closeConnection(websocket.CloseGoingAway)

		case <-socket.done:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *Socket) On(event string, handler goja.Callable) ***REMOVED***
	s.eventHandlers[event] = handler
***REMOVED***

func (s *Socket) handleEvent(event string, args ...goja.Value) ***REMOVED***
	if handler, ok := s.eventHandlers[event]; ok ***REMOVED***
		handler(goja.Undefined(), args...)
	***REMOVED***
***REMOVED***

func (s *Socket) Send(message string) ***REMOVED***
	// NOTE: No binary message support for the time being since goja doesn't
	// support typed arrays.
	rt := common.GetRuntime(s.ctx)

	if err := s.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil ***REMOVED***
		s.handleEvent("error", rt.ToValue(err))
	***REMOVED***
***REMOVED***

func (s *Socket) Ping() ***REMOVED***
	rt := common.GetRuntime(s.ctx)
	deadline := time.Now().Add(writeWait)

	err := s.conn.WriteControl(websocket.PingMessage, []byte***REMOVED******REMOVED***, deadline)
	if err != nil ***REMOVED***
		s.handleEvent("error", rt.ToValue(err))
	***REMOVED***
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
	if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) ***REMOVED***
		code = int(args[0].Export().(int64))
	***REMOVED***

	s.closeConnection(code)
***REMOVED***

func (s *Socket) closeConnection(code int) error ***REMOVED***
	// Attempts to close the websocket gracefully

	rt := common.GetRuntime(s.ctx)

	err := s.conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(code, ""),
		time.Now().Add(writeWait),
	)
	if err != nil ***REMOVED***
		s.handleEvent("error", rt.ToValue(err))

		// Try to close the connection anyway
		s.conn.Close()
		return err
	***REMOVED***

	select ***REMOVED***
	case <-s.done:
	case <-time.After(time.Second):
	***REMOVED***

	s.conn.Close()
	close(s.done)

	return nil
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

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		readChan <- message
	***REMOVED***
***REMOVED***
