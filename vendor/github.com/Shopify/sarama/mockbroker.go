package sarama

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	expectationTimeout = 500 * time.Millisecond
)

type requestHandlerFunc func(req *request) (res encoder)

// RequestNotifierFunc is invoked when a mock broker processes a request successfully
// and will provides the number of bytes read and written.
type RequestNotifierFunc func(bytesRead, bytesWritten int)

// MockBroker is a mock Kafka broker that is used in unit tests. It is exposed
// to facilitate testing of higher level or specialized consumers and producers
// built on top of Sarama. Note that it does not 'mimic' the Kafka API protocol,
// but rather provides a facility to do that. It takes care of the TCP
// transport, request unmarshaling, response marshaling, and makes it the test
// writer responsibility to program correct according to the Kafka API protocol
// MockBroker behaviour.
//
// MockBroker is implemented as a TCP server listening on a kernel-selected
// localhost port that can accept many connections. It reads Kafka requests
// from that connection and returns responses programmed by the SetHandlerByMap
// function. If a MockBroker receives a request that it has no programmed
// response for, then it returns nothing and the request times out.
//
// A set of MockRequest builders to define mappings used by MockBroker is
// provided by Sarama. But users can develop MockRequests of their own and use
// them along with or instead of the standard ones.
//
// When running tests with MockBroker it is strongly recommended to specify
// a timeout to `go test` so that if the broker hangs waiting for a response,
// the test panics.
//
// It is not necessary to prefix message length or correlation ID to your
// response bytes, the server does that automatically as a convenience.
type MockBroker struct ***REMOVED***
	brokerID     int32
	port         int32
	closing      chan none
	stopper      chan none
	expectations chan encoder
	listener     net.Listener
	t            TestReporter
	latency      time.Duration
	handler      requestHandlerFunc
	notifier     RequestNotifierFunc
	history      []RequestResponse
	lock         sync.Mutex
***REMOVED***

// RequestResponse represents a Request/Response pair processed by MockBroker.
type RequestResponse struct ***REMOVED***
	Request  protocolBody
	Response encoder
***REMOVED***

// SetLatency makes broker pause for the specified period every time before
// replying.
func (b *MockBroker) SetLatency(latency time.Duration) ***REMOVED***
	b.latency = latency
***REMOVED***

// SetHandlerByMap defines mapping of Request types to MockResponses. When a
// request is received by the broker, it looks up the request type in the map
// and uses the found MockResponse instance to generate an appropriate reply.
// If the request type is not found in the map then nothing is sent.
func (b *MockBroker) SetHandlerByMap(handlerMap map[string]MockResponse) ***REMOVED***
	b.setHandler(func(req *request) (res encoder) ***REMOVED***
		reqTypeName := reflect.TypeOf(req.body).Elem().Name()
		mockResponse := handlerMap[reqTypeName]
		if mockResponse == nil ***REMOVED***
			return nil
		***REMOVED***
		return mockResponse.For(req.body)
	***REMOVED***)
***REMOVED***

// SetNotifier set a function that will get invoked whenever a request has been
// processed successfully and will provide the number of bytes read and written
func (b *MockBroker) SetNotifier(notifier RequestNotifierFunc) ***REMOVED***
	b.lock.Lock()
	b.notifier = notifier
	b.lock.Unlock()
***REMOVED***

// BrokerID returns broker ID assigned to the broker.
func (b *MockBroker) BrokerID() int32 ***REMOVED***
	return b.brokerID
***REMOVED***

// History returns a slice of RequestResponse pairs in the order they were
// processed by the broker. Note that in case of multiple connections to the
// broker the order expected by a test can be different from the order recorded
// in the history, unless some synchronization is implemented in the test.
func (b *MockBroker) History() []RequestResponse ***REMOVED***
	b.lock.Lock()
	history := make([]RequestResponse, len(b.history))
	copy(history, b.history)
	b.lock.Unlock()
	return history
***REMOVED***

// Port returns the TCP port number the broker is listening for requests on.
func (b *MockBroker) Port() int32 ***REMOVED***
	return b.port
***REMOVED***

// Addr returns the broker connection string in the form "<address>:<port>".
func (b *MockBroker) Addr() string ***REMOVED***
	return b.listener.Addr().String()
***REMOVED***

// Close terminates the broker blocking until it stops internal goroutines and
// releases all resources.
func (b *MockBroker) Close() ***REMOVED***
	close(b.expectations)
	if len(b.expectations) > 0 ***REMOVED***
		buf := bytes.NewBufferString(fmt.Sprintf("mockbroker/%d: not all expectations were satisfied! Still waiting on:\n", b.BrokerID()))
		for e := range b.expectations ***REMOVED***
			_, _ = buf.WriteString(spew.Sdump(e))
		***REMOVED***
		b.t.Error(buf.String())
	***REMOVED***
	close(b.closing)
	<-b.stopper
***REMOVED***

// setHandler sets the specified function as the request handler. Whenever
// a mock broker reads a request from the wire it passes the request to the
// function and sends back whatever the handler function returns.
func (b *MockBroker) setHandler(handler requestHandlerFunc) ***REMOVED***
	b.lock.Lock()
	b.handler = handler
	b.lock.Unlock()
***REMOVED***

func (b *MockBroker) serverLoop() ***REMOVED***
	defer close(b.stopper)
	var err error
	var conn net.Conn

	go func() ***REMOVED***
		<-b.closing
		err := b.listener.Close()
		if err != nil ***REMOVED***
			b.t.Error(err)
		***REMOVED***
	***REMOVED***()

	wg := &sync.WaitGroup***REMOVED******REMOVED***
	i := 0
	for conn, err = b.listener.Accept(); err == nil; conn, err = b.listener.Accept() ***REMOVED***
		wg.Add(1)
		go b.handleRequests(conn, i, wg)
		i++
	***REMOVED***
	wg.Wait()
	Logger.Printf("*** mockbroker/%d: listener closed, err=%v", b.BrokerID(), err)
***REMOVED***

func (b *MockBroker) handleRequests(conn net.Conn, idx int, wg *sync.WaitGroup) ***REMOVED***
	defer wg.Done()
	defer func() ***REMOVED***
		_ = conn.Close()
	***REMOVED***()
	Logger.Printf("*** mockbroker/%d/%d: connection opened", b.BrokerID(), idx)
	var err error

	abort := make(chan none)
	defer close(abort)
	go func() ***REMOVED***
		select ***REMOVED***
		case <-b.closing:
			_ = conn.Close()
		case <-abort:
		***REMOVED***
	***REMOVED***()

	resHeader := make([]byte, 8)
	for ***REMOVED***
		req, bytesRead, err := decodeRequest(conn)
		if err != nil ***REMOVED***
			Logger.Printf("*** mockbroker/%d/%d: invalid request: err=%+v, %+v", b.brokerID, idx, err, spew.Sdump(req))
			b.serverError(err)
			break
		***REMOVED***

		if b.latency > 0 ***REMOVED***
			time.Sleep(b.latency)
		***REMOVED***

		b.lock.Lock()
		res := b.handler(req)
		b.history = append(b.history, RequestResponse***REMOVED***req.body, res***REMOVED***)
		b.lock.Unlock()

		if res == nil ***REMOVED***
			Logger.Printf("*** mockbroker/%d/%d: ignored %v", b.brokerID, idx, spew.Sdump(req))
			continue
		***REMOVED***
		Logger.Printf("*** mockbroker/%d/%d: served %v -> %v", b.brokerID, idx, req, res)

		encodedRes, err := encode(res, nil)
		if err != nil ***REMOVED***
			b.serverError(err)
			break
		***REMOVED***
		if len(encodedRes) == 0 ***REMOVED***
			b.lock.Lock()
			if b.notifier != nil ***REMOVED***
				b.notifier(bytesRead, 0)
			***REMOVED***
			b.lock.Unlock()
			continue
		***REMOVED***

		binary.BigEndian.PutUint32(resHeader, uint32(len(encodedRes)+4))
		binary.BigEndian.PutUint32(resHeader[4:], uint32(req.correlationID))
		if _, err = conn.Write(resHeader); err != nil ***REMOVED***
			b.serverError(err)
			break
		***REMOVED***
		if _, err = conn.Write(encodedRes); err != nil ***REMOVED***
			b.serverError(err)
			break
		***REMOVED***

		b.lock.Lock()
		if b.notifier != nil ***REMOVED***
			b.notifier(bytesRead, len(resHeader)+len(encodedRes))
		***REMOVED***
		b.lock.Unlock()
	***REMOVED***
	Logger.Printf("*** mockbroker/%d/%d: connection closed, err=%v", b.BrokerID(), idx, err)
***REMOVED***

func (b *MockBroker) defaultRequestHandler(req *request) (res encoder) ***REMOVED***
	select ***REMOVED***
	case res, ok := <-b.expectations:
		if !ok ***REMOVED***
			return nil
		***REMOVED***
		return res
	case <-time.After(expectationTimeout):
		return nil
	***REMOVED***
***REMOVED***

func (b *MockBroker) serverError(err error) ***REMOVED***
	isConnectionClosedError := false
	if _, ok := err.(*net.OpError); ok ***REMOVED***
		isConnectionClosedError = true
	***REMOVED*** else if err == io.EOF ***REMOVED***
		isConnectionClosedError = true
	***REMOVED*** else if err.Error() == "use of closed network connection" ***REMOVED***
		isConnectionClosedError = true
	***REMOVED***

	if isConnectionClosedError ***REMOVED***
		return
	***REMOVED***

	b.t.Errorf(err.Error())
***REMOVED***

// NewMockBroker launches a fake Kafka broker. It takes a TestReporter as provided by the
// test framework and a channel of responses to use.  If an error occurs it is
// simply logged to the TestReporter and the broker exits.
func NewMockBroker(t TestReporter, brokerID int32) *MockBroker ***REMOVED***
	return NewMockBrokerAddr(t, brokerID, "localhost:0")
***REMOVED***

// NewMockBrokerAddr behaves like newMockBroker but listens on the address you give
// it rather than just some ephemeral port.
func NewMockBrokerAddr(t TestReporter, brokerID int32, addr string) *MockBroker ***REMOVED***
	listener, err := net.Listen("tcp", addr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	return NewMockBrokerListener(t, brokerID, listener)
***REMOVED***

// NewMockBrokerListener behaves like newMockBrokerAddr but accepts connections on the listener specified.
func NewMockBrokerListener(t TestReporter, brokerID int32, listener net.Listener) *MockBroker ***REMOVED***
	var err error

	broker := &MockBroker***REMOVED***
		closing:      make(chan none),
		stopper:      make(chan none),
		t:            t,
		brokerID:     brokerID,
		expectations: make(chan encoder, 512),
		listener:     listener,
	***REMOVED***
	broker.handler = broker.defaultRequestHandler

	Logger.Printf("*** mockbroker/%d listening on %s\n", brokerID, broker.listener.Addr().String())
	_, portStr, err := net.SplitHostPort(broker.listener.Addr().String())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tmp, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	broker.port = int32(tmp)

	go broker.serverLoop()

	return broker
***REMOVED***

func (b *MockBroker) Returns(e encoder) ***REMOVED***
	b.expectations <- e
***REMOVED***
