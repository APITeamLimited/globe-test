package ws

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"

	"github.com/APITeamLimited/globe-test/worker/js/common"
	"github.com/APITeamLimited/globe-test/worker/js/modules"
	httpModule "github.com/APITeamLimited/globe-test/worker/js/modules/k6/http"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct***REMOVED******REMOVED***

	// WS represents a module instance of the WebSocket module.
	WS struct ***REMOVED***
		vu  modules.VU
		obj *goja.Object
	***REMOVED***
)

var (
	_ modules.Module   = &RootModule***REMOVED******REMOVED***
	_ modules.Instance = &WS***REMOVED******REMOVED***
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule ***REMOVED***
	return &RootModule***REMOVED******REMOVED***
***REMOVED***

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(m modules.VU) modules.Instance ***REMOVED***
	rt := m.Runtime()
	mi := &WS***REMOVED***
		vu: m,
	***REMOVED***
	obj := rt.NewObject()
	if err := obj.Set("connect", mi.Connect); err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	mi.obj = obj
	return mi
***REMOVED***

// ErrWSInInitContext is returned when websockets are using in the init context
var ErrWSInInitContext = common.NewInitContextError("using websockets in the init context is not supported")

type Socket struct ***REMOVED***
	rt            *goja.Runtime
	ctx           context.Context
	conn          *websocket.Conn
	eventHandlers map[string][]goja.Callable
	scheduled     chan goja.Callable
	done          chan struct***REMOVED******REMOVED***
	shutdownOnce  sync.Once

	pingSendTimestamps map[string]time.Time
	pingSendCounter    int

	sampleTags     *workerMetrics.SampleTags
	samplesOutput  chan<- workerMetrics.SampleContainer
	builtinMetrics *workerMetrics.BuiltinMetrics
***REMOVED***

type WSHTTPResponse struct ***REMOVED***
	URL     string            `json:"url"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Error   string            `json:"error"`
***REMOVED***

type message struct ***REMOVED***
	mtype int // message type consts as defined in gorilla/websocket/conn.go
	data  []byte
***REMOVED***

const writeWait = 10 * time.Second

// Exports returns the exports of the ws module.
func (mi *WS) Exports() modules.Exports ***REMOVED***
	return modules.Exports***REMOVED***Default: mi.obj***REMOVED***
***REMOVED***

// Connect establishes a WebSocket connection based on the parameters provided.
// TODO: refactor to reduce the method complexity
//nolint:funlen,gocognit,gocyclo
func (mi *WS) Connect(url string, args ...goja.Value) (*WSHTTPResponse, error) ***REMOVED***
	ctx := mi.vu.Context()
	rt := mi.vu.Runtime()
	state := mi.vu.State()
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

	header := make(http.Header)
	header.Set("User-Agent", state.Options.UserAgent.String)

	enableCompression := false

	tags := state.CloneTags()
	jar := state.CookieJar

	// Parse the optional second argument (params)
	if !goja.IsUndefined(paramsV) && !goja.IsNull(paramsV) ***REMOVED***
		params := paramsV.ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "headers":
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
			case "jar":
				jarV := params.Get(k)
				if goja.IsUndefined(jarV) || goja.IsNull(jarV) ***REMOVED***
					continue
				***REMOVED***
				if v, ok := jarV.Export().(*httpModule.CookieJar); ok ***REMOVED***
					jar = v.Jar
				***REMOVED***
			case "compression":
				// deflate compression algorithm is supported - as defined in RFC7692
				// compression here relies on the implementation in gorilla/websocket package, usage is
				// experimental and may result in decreased performance. package supports
				// only "no context takeover" scenario

				algoString := strings.TrimSpace(params.Get(k).ToString().String())
				if algoString == "" ***REMOVED***
					continue
				***REMOVED***

				if algoString != "deflate" ***REMOVED***
					return nil, fmt.Errorf("unsupported compression algorithm '%s', supported algorithm is 'deflate'", algoString)
				***REMOVED***

				enableCompression = true
			***REMOVED***
		***REMOVED***

	***REMOVED***

	if state.Options.SystemTags.Has(workerMetrics.TagURL) ***REMOVED***
		tags["url"] = url
	***REMOVED***

	// Overriding the NextProtos to avoid talking http2
	var tlsConfig *tls.Config
	if state.TLSConfig != nil ***REMOVED***
		tlsConfig = state.TLSConfig.Clone()
		tlsConfig.NextProtos = []string***REMOVED***"http/1.1"***REMOVED***
	***REMOVED***

	wsd := websocket.Dialer***REMOVED***
		HandshakeTimeout: time.Second * 60, // TODO configurable
		// Pass a custom net.DialContext function to websocket.Dialer that will substitute
		// the underlying net.Conn with our own tracked netext.Conn
		NetDialContext:    state.Dialer.DialContext,
		Proxy:             http.ProxyFromEnvironment,
		TLSClientConfig:   tlsConfig,
		EnableCompression: enableCompression,
		Jar:               jar,
	***REMOVED***
	if jar == nil ***REMOVED*** // this is needed because of how interfaces work and that wsd.Jar is http.Cookiejar
		wsd.Jar = nil
	***REMOVED***

	start := time.Now()
	conn, httpResponse, connErr := wsd.DialContext(ctx, url, header)
	connectionEnd := time.Now()
	connectionDuration := workerMetrics.D(connectionEnd.Sub(start))

	if state.Options.SystemTags.Has(workerMetrics.TagIP) && conn.RemoteAddr() != nil ***REMOVED***
		if ip, _, err := net.SplitHostPort(conn.RemoteAddr().String()); err == nil ***REMOVED***
			tags["ip"] = ip
		***REMOVED***
	***REMOVED***

	if httpResponse != nil ***REMOVED***
		if state.Options.SystemTags.Has(workerMetrics.TagStatus) ***REMOVED***
			tags["status"] = strconv.Itoa(httpResponse.StatusCode)
		***REMOVED***

		if state.Options.SystemTags.Has(workerMetrics.TagSubproto) ***REMOVED***
			tags["subproto"] = httpResponse.Header.Get("Sec-WebSocket-Protocol")
		***REMOVED***
	***REMOVED***

	socket := Socket***REMOVED***
		ctx:                ctx,
		rt:                 rt,
		conn:               conn,
		eventHandlers:      make(map[string][]goja.Callable),
		pingSendTimestamps: make(map[string]time.Time),
		scheduled:          make(chan goja.Callable),
		done:               make(chan struct***REMOVED******REMOVED***),
		samplesOutput:      state.Samples,
		sampleTags:         workerMetrics.IntoSampleTags(&tags),
		builtinMetrics:     state.BuiltinMetrics,
	***REMOVED***

	workerMetrics.PushIfNotDone(ctx, state.Samples, workerMetrics.ConnectedSamples***REMOVED***
		Samples: []workerMetrics.Sample***REMOVED***
			***REMOVED***Metric: state.BuiltinMetrics.WSSessions, Time: start, Tags: socket.sampleTags, Value: 1***REMOVED***,
			***REMOVED***Metric: state.BuiltinMetrics.WSConnecting, Time: start, Tags: socket.sampleTags, Value: connectionDuration***REMOVED***,
		***REMOVED***,
		Tags: socket.sampleTags,
		Time: start,
	***REMOVED***)

	if connErr != nil ***REMOVED***
		// Pass the error to the user script before exiting immediately
		socket.handleEvent("error", rt.ToValue(connErr))
		if state.Options.Throw.Bool ***REMOVED***
			return nil, connErr
		***REMOVED***
		state.Logger.WithError(connErr).Warn("Attempt to establish a WebSocket connection failed")
		return &WSHTTPResponse***REMOVED***
			Error: connErr.Error(),
		***REMOVED***, nil
	***REMOVED***

	// Run the user-provided set up function
	if _, err := setupFn(goja.Undefined(), rt.ToValue(&socket)); err != nil ***REMOVED***
		_ = socket.closeConnection(websocket.CloseGoingAway)
		return nil, err
	***REMOVED***
	wsResponse, wsRespErr := wrapHTTPResponse(httpResponse)
	if wsRespErr != nil ***REMOVED***
		return nil, wsRespErr
	***REMOVED***
	wsResponse.URL = url

	defer func() ***REMOVED*** _ = conn.Close() ***REMOVED***()

	// The connection is now open, emit the event
	socket.handleEvent("open")

	// Make the default close handler a noop to avoid duplicate closes,
	// since we use custom closing logic to call user's event
	// handlers and for cleanup. See closeConnection.
	// closeConnection is not set directly as a handler here to
	// avoid race conditions when calling the Goja runtime.
	conn.SetCloseHandler(func(code int, text string) error ***REMOVED*** return nil ***REMOVED***)

	// Pass ping/pong events through the main control loop
	pingChan := make(chan string)
	pongChan := make(chan string)
	conn.SetPingHandler(func(msg string) error ***REMOVED*** pingChan <- msg; return nil ***REMOVED***)
	conn.SetPongHandler(func(pingID string) error ***REMOVED*** pongChan <- pingID; return nil ***REMOVED***)

	readDataChan := make(chan *message)
	readCloseChan := make(chan int)
	readErrChan := make(chan error)

	// Wraps a couple of channels around conn.ReadMessage
	go socket.readPump(readDataChan, readErrChan, readCloseChan)

	// we do it here as below we can panic, which translates to an exception in js code
	defer func() ***REMOVED***
		socket.Close() // just in case
		end := time.Now()
		sessionDuration := workerMetrics.D(end.Sub(start))

		workerMetrics.PushIfNotDone(ctx, state.Samples, workerMetrics.Sample***REMOVED***
			Metric: socket.builtinMetrics..WSSessionDuration,
			Tags:   socket.sampleTags,
			Time:   start,
			Value:  sessionDuration,
		***REMOVED***)
	***REMOVED***()

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

		case msg := <-readDataChan:
			workerMetrics.PushIfNotDone(ctx, socket.samplesOutput, workerMetrics.Sample***REMOVED***
				Metric: socket.builtinMetrics..WSMessagesReceived,
				Time:   time.Now(),
				Tags:   socket.sampleTags,
				Value:  1,
			***REMOVED***)

			if msg.mtype == websocket.BinaryMessage ***REMOVED***
				ab := rt.NewArrayBuffer(msg.data)
				socket.handleEvent("binaryMessage", rt.ToValue(&ab))
			***REMOVED*** else ***REMOVED***
				socket.handleEvent("message", rt.ToValue(string(msg.data)))
			***REMOVED***

		case readErr := <-readErrChan:
			socket.handleEvent("error", rt.ToValue(readErr))

		case code := <-readCloseChan:
			_ = socket.closeConnection(code)

		case scheduledFn := <-socket.scheduled:
			if _, err := scheduledFn(goja.Undefined()); err != nil ***REMOVED***
				_ = socket.closeConnection(websocket.CloseGoingAway)
				return nil, err
			***REMOVED***

		case <-ctx.Done():
			// VU is shutting down during an interrupt
			// socket events will not be forwarded to the VU
			_ = socket.closeConnection(websocket.CloseGoingAway)

		case <-socket.done:
			// This is the final exit point normally triggered by closeConnection
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
				common.Throw(s.rt, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Send writes the given string message to the connection.
func (s *Socket) Send(message string) ***REMOVED***
	if err := s.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil ***REMOVED***
		s.handleEvent("error", s.rt.ToValue(err))
	***REMOVED***

	workerMetrics.PushIfNotDone(s.ctx, s.samplesOutput, workerMetrics.Sample***REMOVED***
		Metric: s.builtinMetrics..WSMessagesSent,
		Time:   time.Now(),
		Tags:   s.sampleTags,
		Value:  1,
	***REMOVED***)
***REMOVED***

// SendBinary writes the given ArrayBuffer message to the connection.
func (s *Socket) SendBinary(message goja.Value) ***REMOVED***
	if message == nil ***REMOVED***
		common.Throw(s.rt, errors.New("missing argument, expected ArrayBuffer"))
	***REMOVED***

	msg := message.Export()
	if ab, ok := msg.(goja.ArrayBuffer); ok ***REMOVED***
		if err := s.conn.WriteMessage(websocket.BinaryMessage, ab.Bytes()); err != nil ***REMOVED***
			s.handleEvent("error", s.rt.ToValue(err))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var jsType string
		switch ***REMOVED***
		case goja.IsNull(message), goja.IsUndefined(message):
			jsType = message.String()
		default:
			jsType = message.ToObject(s.rt).ClassName()
		***REMOVED***
		common.Throw(s.rt, fmt.Errorf("expected ArrayBuffer as argument, received: %s", jsType))
	***REMOVED***

	workerMetrics.PushIfNotDone(s.ctx, s.samplesOutput, workerMetrics.Sample***REMOVED***
		Metric: s.builtinMetrics..WSMessagesSent,
		Time:   time.Now(),
		Tags:   s.sampleTags,
		Value:  1,
	***REMOVED***)
***REMOVED***

func (s *Socket) Ping() ***REMOVED***
	deadline := time.Now().Add(writeWait)
	pingID := strconv.Itoa(s.pingSendCounter)
	data := []byte(pingID)

	err := s.conn.WriteControl(websocket.PingMessage, data, deadline)
	if err != nil ***REMOVED***
		s.handleEvent("error", s.rt.ToValue(err))
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

	workerMetrics.PushIfNotDone(s.ctx, s.samplesOutput, workerMetrics.Sample***REMOVED***
		Metric: s.builtinMetrics..WSPing,
		Time:   pongTimestamp,
		Tags:   s.sampleTags,
		Value:  workerMetrics.D(pongTimestamp.Sub(pingTimestamp)),
	***REMOVED***)
***REMOVED***

// SetTimeout executes the provided function inside the socket's event loop after at least the provided
// timeout, which is in ms, has elapsed
func (s *Socket) SetTimeout(fn goja.Callable, timeoutMs float64) error ***REMOVED***
	// Starts a goroutine, blocks once on the timeout and pushes the callable
	// back to the main loop through the scheduled channel.
	//
	// Intentionally not using the generic GetDurationValue() helper, since this
	// API is meant to use ms, similar to the original SetTimeout() JS API.
	d := time.Duration(timeoutMs * float64(time.Millisecond))
	if d <= 0 ***REMOVED***
		return fmt.Errorf("setTimeout requires a >0 timeout parameter, received %.2f", timeoutMs)
	***REMOVED***
	go func() ***REMOVED***
		select ***REMOVED***
		case <-time.After(d):
			select ***REMOVED***
			case s.scheduled <- fn:
			case <-s.done:
				return
			***REMOVED***

		case <-s.done:
			return
		***REMOVED***
	***REMOVED***()

	return nil
***REMOVED***

// SetInterval executes the provided function inside the socket's event loop each interval time, which is
// in ms
func (s *Socket) SetInterval(fn goja.Callable, intervalMs float64) error ***REMOVED***
	// Starts a goroutine, blocks forever on the ticker and pushes the callable
	// back to the main loop through the scheduled channel.
	//
	// Intentionally not using the generic GetDurationValue() helper, since this
	// API is meant to use ms, similar to the original SetInterval() JS API.
	d := time.Duration(intervalMs * float64(time.Millisecond))
	if d <= 0 ***REMOVED***
		return fmt.Errorf("setInterval requires a >0 timeout parameter, received %.2f", intervalMs)
	***REMOVED***
	go func() ***REMOVED***
		ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
		defer ticker.Stop()

		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				select ***REMOVED***
				case s.scheduled <- fn:
				case <-s.done:
					return
				***REMOVED***

			case <-s.done:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return nil
***REMOVED***

func (s *Socket) Close(args ...goja.Value) ***REMOVED***
	code := websocket.CloseGoingAway
	if len(args) > 0 ***REMOVED***
		code = int(args[0].ToInteger())
	***REMOVED***

	_ = s.closeConnection(code)
***REMOVED***

// closeConnection cleanly closes the WebSocket connection.
// Returns an error if sending the close control frame fails.
func (s *Socket) closeConnection(code int) error ***REMOVED***
	var err error

	s.shutdownOnce.Do(func() ***REMOVED***
		// this is because handleEvent can panic ... on purpose so we just make sure we
		// close the connection and the channel
		defer func() ***REMOVED***
			_ = s.conn.Close()

			// Stop the main control loop
			close(s.done)
		***REMOVED***()
		err = s.conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(code, ""),
			time.Now().Add(writeWait),
		)
		if err != nil ***REMOVED***
			// Call the user-defined error handler
			s.handleEvent("error", s.rt.ToValue(err))
		***REMOVED***

		// Call the user-defined close handler
		s.handleEvent("close", s.rt.ToValue(code))
	***REMOVED***)

	return err
***REMOVED***

// Wraps conn.ReadMessage in a channel
func (s *Socket) readPump(readChan chan *message, errorChan chan error, closeChan chan int) ***REMOVED***
	for ***REMOVED***
		messageType, data, err := s.conn.ReadMessage()
		if err != nil ***REMOVED***
			if websocket.IsUnexpectedCloseError(
				err, websocket.CloseNormalClosure, websocket.CloseGoingAway) ***REMOVED***
				// Report an unexpected closure
				select ***REMOVED***
				case errorChan <- err:
				case <-s.done:
					return
				***REMOVED***
			***REMOVED***
			code := websocket.CloseGoingAway
			if e, ok := err.(*websocket.CloseError); ok ***REMOVED***
				code = e.Code
			***REMOVED***
			select ***REMOVED***
			case closeChan <- code:
			case <-s.done:
			***REMOVED***
			return
		***REMOVED***

		select ***REMOVED***
		case readChan <- &message***REMOVED***messageType, data***REMOVED***:
		case <-s.done:
			return
		***REMOVED***
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
