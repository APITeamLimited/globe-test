package sarama

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rcrowley/go-metrics"
)

// Broker represents a single Kafka broker connection. All operations on this object are entirely concurrency-safe.
type Broker struct ***REMOVED***
	id   int32
	addr string

	conf          *Config
	correlationID int32
	conn          net.Conn
	connErr       error
	lock          sync.Mutex
	opened        int32

	responses chan responsePromise
	done      chan bool

	incomingByteRate       metrics.Meter
	requestRate            metrics.Meter
	requestSize            metrics.Histogram
	requestLatency         metrics.Histogram
	outgoingByteRate       metrics.Meter
	responseRate           metrics.Meter
	responseSize           metrics.Histogram
	brokerIncomingByteRate metrics.Meter
	brokerRequestRate      metrics.Meter
	brokerRequestSize      metrics.Histogram
	brokerRequestLatency   metrics.Histogram
	brokerOutgoingByteRate metrics.Meter
	brokerResponseRate     metrics.Meter
	brokerResponseSize     metrics.Histogram
***REMOVED***

type responsePromise struct ***REMOVED***
	requestTime   time.Time
	correlationID int32
	packets       chan []byte
	errors        chan error
***REMOVED***

// NewBroker creates and returns a Broker targeting the given host:port address.
// This does not attempt to actually connect, you have to call Open() for that.
func NewBroker(addr string) *Broker ***REMOVED***
	return &Broker***REMOVED***id: -1, addr: addr***REMOVED***
***REMOVED***

// Open tries to connect to the Broker if it is not already connected or connecting, but does not block
// waiting for the connection to complete. This means that any subsequent operations on the broker will
// block waiting for the connection to succeed or fail. To get the effect of a fully synchronous Open call,
// follow it by a call to Connected(). The only errors Open will return directly are ConfigurationError or
// AlreadyConnected. If conf is nil, the result of NewConfig() is used.
func (b *Broker) Open(conf *Config) error ***REMOVED***
	if !atomic.CompareAndSwapInt32(&b.opened, 0, 1) ***REMOVED***
		return ErrAlreadyConnected
	***REMOVED***

	if conf == nil ***REMOVED***
		conf = NewConfig()
	***REMOVED***

	err := conf.Validate()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	b.lock.Lock()

	go withRecover(func() ***REMOVED***
		defer b.lock.Unlock()

		dialer := net.Dialer***REMOVED***
			Timeout:   conf.Net.DialTimeout,
			KeepAlive: conf.Net.KeepAlive,
		***REMOVED***

		if conf.Net.TLS.Enable ***REMOVED***
			b.conn, b.connErr = tls.DialWithDialer(&dialer, "tcp", b.addr, conf.Net.TLS.Config)
		***REMOVED*** else ***REMOVED***
			b.conn, b.connErr = dialer.Dial("tcp", b.addr)
		***REMOVED***
		if b.connErr != nil ***REMOVED***
			Logger.Printf("Failed to connect to broker %s: %s\n", b.addr, b.connErr)
			b.conn = nil
			atomic.StoreInt32(&b.opened, 0)
			return
		***REMOVED***
		b.conn = newBufConn(b.conn)

		b.conf = conf

		// Create or reuse the global metrics shared between brokers
		b.incomingByteRate = metrics.GetOrRegisterMeter("incoming-byte-rate", conf.MetricRegistry)
		b.requestRate = metrics.GetOrRegisterMeter("request-rate", conf.MetricRegistry)
		b.requestSize = getOrRegisterHistogram("request-size", conf.MetricRegistry)
		b.requestLatency = getOrRegisterHistogram("request-latency-in-ms", conf.MetricRegistry)
		b.outgoingByteRate = metrics.GetOrRegisterMeter("outgoing-byte-rate", conf.MetricRegistry)
		b.responseRate = metrics.GetOrRegisterMeter("response-rate", conf.MetricRegistry)
		b.responseSize = getOrRegisterHistogram("response-size", conf.MetricRegistry)
		// Do not gather metrics for seeded broker (only used during bootstrap) because they share
		// the same id (-1) and are already exposed through the global metrics above
		if b.id >= 0 ***REMOVED***
			b.brokerIncomingByteRate = getOrRegisterBrokerMeter("incoming-byte-rate", b, conf.MetricRegistry)
			b.brokerRequestRate = getOrRegisterBrokerMeter("request-rate", b, conf.MetricRegistry)
			b.brokerRequestSize = getOrRegisterBrokerHistogram("request-size", b, conf.MetricRegistry)
			b.brokerRequestLatency = getOrRegisterBrokerHistogram("request-latency-in-ms", b, conf.MetricRegistry)
			b.brokerOutgoingByteRate = getOrRegisterBrokerMeter("outgoing-byte-rate", b, conf.MetricRegistry)
			b.brokerResponseRate = getOrRegisterBrokerMeter("response-rate", b, conf.MetricRegistry)
			b.brokerResponseSize = getOrRegisterBrokerHistogram("response-size", b, conf.MetricRegistry)
		***REMOVED***

		if conf.Net.SASL.Enable ***REMOVED***
			b.connErr = b.sendAndReceiveSASLPlainAuth()
			if b.connErr != nil ***REMOVED***
				err = b.conn.Close()
				if err == nil ***REMOVED***
					Logger.Printf("Closed connection to broker %s\n", b.addr)
				***REMOVED*** else ***REMOVED***
					Logger.Printf("Error while closing connection to broker %s: %s\n", b.addr, err)
				***REMOVED***
				b.conn = nil
				atomic.StoreInt32(&b.opened, 0)
				return
			***REMOVED***
		***REMOVED***

		b.done = make(chan bool)
		b.responses = make(chan responsePromise, b.conf.Net.MaxOpenRequests-1)

		if b.id >= 0 ***REMOVED***
			Logger.Printf("Connected to broker at %s (registered as #%d)\n", b.addr, b.id)
		***REMOVED*** else ***REMOVED***
			Logger.Printf("Connected to broker at %s (unregistered)\n", b.addr)
		***REMOVED***
		go withRecover(b.responseReceiver)
	***REMOVED***)

	return nil
***REMOVED***

// Connected returns true if the broker is connected and false otherwise. If the broker is not
// connected but it had tried to connect, the error from that connection attempt is also returned.
func (b *Broker) Connected() (bool, error) ***REMOVED***
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.conn != nil, b.connErr
***REMOVED***

func (b *Broker) Close() error ***REMOVED***
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.conn == nil ***REMOVED***
		return ErrNotConnected
	***REMOVED***

	close(b.responses)
	<-b.done

	err := b.conn.Close()

	b.conn = nil
	b.connErr = nil
	b.done = nil
	b.responses = nil

	if b.id >= 0 ***REMOVED***
		b.conf.MetricRegistry.Unregister(getMetricNameForBroker("incoming-byte-rate", b))
		b.conf.MetricRegistry.Unregister(getMetricNameForBroker("request-rate", b))
		b.conf.MetricRegistry.Unregister(getMetricNameForBroker("outgoing-byte-rate", b))
		b.conf.MetricRegistry.Unregister(getMetricNameForBroker("response-rate", b))
	***REMOVED***

	if err == nil ***REMOVED***
		Logger.Printf("Closed connection to broker %s\n", b.addr)
	***REMOVED*** else ***REMOVED***
		Logger.Printf("Error while closing connection to broker %s: %s\n", b.addr, err)
	***REMOVED***

	atomic.StoreInt32(&b.opened, 0)

	return err
***REMOVED***

// ID returns the broker ID retrieved from Kafka's metadata, or -1 if that is not known.
func (b *Broker) ID() int32 ***REMOVED***
	return b.id
***REMOVED***

// Addr returns the broker address as either retrieved from Kafka's metadata or passed to NewBroker.
func (b *Broker) Addr() string ***REMOVED***
	return b.addr
***REMOVED***

func (b *Broker) GetMetadata(request *MetadataRequest) (*MetadataResponse, error) ***REMOVED***
	response := new(MetadataResponse)

	err := b.sendAndReceive(request, response)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) GetConsumerMetadata(request *ConsumerMetadataRequest) (*ConsumerMetadataResponse, error) ***REMOVED***
	response := new(ConsumerMetadataResponse)

	err := b.sendAndReceive(request, response)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) GetAvailableOffsets(request *OffsetRequest) (*OffsetResponse, error) ***REMOVED***
	response := new(OffsetResponse)

	err := b.sendAndReceive(request, response)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) Produce(request *ProduceRequest) (*ProduceResponse, error) ***REMOVED***
	var response *ProduceResponse
	var err error

	if request.RequiredAcks == NoResponse ***REMOVED***
		err = b.sendAndReceive(request, nil)
	***REMOVED*** else ***REMOVED***
		response = new(ProduceResponse)
		err = b.sendAndReceive(request, response)
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) Fetch(request *FetchRequest) (*FetchResponse, error) ***REMOVED***
	response := new(FetchResponse)

	err := b.sendAndReceive(request, response)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) CommitOffset(request *OffsetCommitRequest) (*OffsetCommitResponse, error) ***REMOVED***
	response := new(OffsetCommitResponse)

	err := b.sendAndReceive(request, response)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) FetchOffset(request *OffsetFetchRequest) (*OffsetFetchResponse, error) ***REMOVED***
	response := new(OffsetFetchResponse)

	err := b.sendAndReceive(request, response)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) JoinGroup(request *JoinGroupRequest) (*JoinGroupResponse, error) ***REMOVED***
	response := new(JoinGroupResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) SyncGroup(request *SyncGroupRequest) (*SyncGroupResponse, error) ***REMOVED***
	response := new(SyncGroupResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) LeaveGroup(request *LeaveGroupRequest) (*LeaveGroupResponse, error) ***REMOVED***
	response := new(LeaveGroupResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) Heartbeat(request *HeartbeatRequest) (*HeartbeatResponse, error) ***REMOVED***
	response := new(HeartbeatResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) ListGroups(request *ListGroupsRequest) (*ListGroupsResponse, error) ***REMOVED***
	response := new(ListGroupsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) DescribeGroups(request *DescribeGroupsRequest) (*DescribeGroupsResponse, error) ***REMOVED***
	response := new(DescribeGroupsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) ApiVersions(request *ApiVersionsRequest) (*ApiVersionsResponse, error) ***REMOVED***
	response := new(ApiVersionsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) CreateTopics(request *CreateTopicsRequest) (*CreateTopicsResponse, error) ***REMOVED***
	response := new(CreateTopicsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) DeleteTopics(request *DeleteTopicsRequest) (*DeleteTopicsResponse, error) ***REMOVED***
	response := new(DeleteTopicsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) DescribeAcls(request *DescribeAclsRequest) (*DescribeAclsResponse, error) ***REMOVED***
	response := new(DescribeAclsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) CreateAcls(request *CreateAclsRequest) (*CreateAclsResponse, error) ***REMOVED***
	response := new(CreateAclsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) DeleteAcls(request *DeleteAclsRequest) (*DeleteAclsResponse, error) ***REMOVED***
	response := new(DeleteAclsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) InitProducerID(request *InitProducerIDRequest) (*InitProducerIDResponse, error) ***REMOVED***
	response := new(InitProducerIDResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) AddPartitionsToTxn(request *AddPartitionsToTxnRequest) (*AddPartitionsToTxnResponse, error) ***REMOVED***
	response := new(AddPartitionsToTxnResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) AddOffsetsToTxn(request *AddOffsetsToTxnRequest) (*AddOffsetsToTxnResponse, error) ***REMOVED***
	response := new(AddOffsetsToTxnResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) EndTxn(request *EndTxnRequest) (*EndTxnResponse, error) ***REMOVED***
	response := new(EndTxnResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) TxnOffsetCommit(request *TxnOffsetCommitRequest) (*TxnOffsetCommitResponse, error) ***REMOVED***
	response := new(TxnOffsetCommitResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) DescribeConfigs(request *DescribeConfigsRequest) (*DescribeConfigsResponse, error) ***REMOVED***
	response := new(DescribeConfigsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***

func (b *Broker) AlterConfigs(request *AlterConfigsRequest) (*AlterConfigsResponse, error) ***REMOVED***
	response := new(AlterConfigsResponse)

	err := b.sendAndReceive(request, response)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return response, nil
***REMOVED***
func (b *Broker) send(rb protocolBody, promiseResponse bool) (*responsePromise, error) ***REMOVED***
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.conn == nil ***REMOVED***
		if b.connErr != nil ***REMOVED***
			return nil, b.connErr
		***REMOVED***
		return nil, ErrNotConnected
	***REMOVED***

	if !b.conf.Version.IsAtLeast(rb.requiredVersion()) ***REMOVED***
		return nil, ErrUnsupportedVersion
	***REMOVED***

	req := &request***REMOVED***correlationID: b.correlationID, clientID: b.conf.ClientID, body: rb***REMOVED***
	buf, err := encode(req, b.conf.MetricRegistry)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = b.conn.SetWriteDeadline(time.Now().Add(b.conf.Net.WriteTimeout))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	requestTime := time.Now()
	bytes, err := b.conn.Write(buf)
	b.updateOutgoingCommunicationMetrics(bytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	b.correlationID++

	if !promiseResponse ***REMOVED***
		// Record request latency without the response
		b.updateRequestLatencyMetrics(time.Since(requestTime))
		return nil, nil
	***REMOVED***

	promise := responsePromise***REMOVED***requestTime, req.correlationID, make(chan []byte), make(chan error)***REMOVED***
	b.responses <- promise

	return &promise, nil
***REMOVED***

func (b *Broker) sendAndReceive(req protocolBody, res versionedDecoder) error ***REMOVED***
	promise, err := b.send(req, res != nil)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if promise == nil ***REMOVED***
		return nil
	***REMOVED***

	select ***REMOVED***
	case buf := <-promise.packets:
		return versionedDecode(buf, res, req.version())
	case err = <-promise.errors:
		return err
	***REMOVED***
***REMOVED***

func (b *Broker) decode(pd packetDecoder) (err error) ***REMOVED***
	b.id, err = pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	host, err := pd.getString()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	port, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	b.addr = net.JoinHostPort(host, fmt.Sprint(port))
	if _, _, err := net.SplitHostPort(b.addr); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (b *Broker) encode(pe packetEncoder) (err error) ***REMOVED***

	host, portstr, err := net.SplitHostPort(b.addr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	port, err := strconv.Atoi(portstr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt32(b.id)

	err = pe.putString(host)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pe.putInt32(int32(port))

	return nil
***REMOVED***

func (b *Broker) responseReceiver() ***REMOVED***
	var dead error
	header := make([]byte, 8)
	for response := range b.responses ***REMOVED***
		if dead != nil ***REMOVED***
			response.errors <- dead
			continue
		***REMOVED***

		err := b.conn.SetReadDeadline(time.Now().Add(b.conf.Net.ReadTimeout))
		if err != nil ***REMOVED***
			dead = err
			response.errors <- err
			continue
		***REMOVED***

		bytesReadHeader, err := io.ReadFull(b.conn, header)
		requestLatency := time.Since(response.requestTime)
		if err != nil ***REMOVED***
			b.updateIncomingCommunicationMetrics(bytesReadHeader, requestLatency)
			dead = err
			response.errors <- err
			continue
		***REMOVED***

		decodedHeader := responseHeader***REMOVED******REMOVED***
		err = decode(header, &decodedHeader)
		if err != nil ***REMOVED***
			b.updateIncomingCommunicationMetrics(bytesReadHeader, requestLatency)
			dead = err
			response.errors <- err
			continue
		***REMOVED***
		if decodedHeader.correlationID != response.correlationID ***REMOVED***
			b.updateIncomingCommunicationMetrics(bytesReadHeader, requestLatency)
			// TODO if decoded ID < cur ID, discard until we catch up
			// TODO if decoded ID > cur ID, save it so when cur ID catches up we have a response
			dead = PacketDecodingError***REMOVED***fmt.Sprintf("correlation ID didn't match, wanted %d, got %d", response.correlationID, decodedHeader.correlationID)***REMOVED***
			response.errors <- dead
			continue
		***REMOVED***

		buf := make([]byte, decodedHeader.length-4)
		bytesReadBody, err := io.ReadFull(b.conn, buf)
		b.updateIncomingCommunicationMetrics(bytesReadHeader+bytesReadBody, requestLatency)
		if err != nil ***REMOVED***
			dead = err
			response.errors <- err
			continue
		***REMOVED***

		response.packets <- buf
	***REMOVED***
	close(b.done)
***REMOVED***

func (b *Broker) sendAndReceiveSASLPlainHandshake() error ***REMOVED***
	rb := &SaslHandshakeRequest***REMOVED***"PLAIN"***REMOVED***
	req := &request***REMOVED***correlationID: b.correlationID, clientID: b.conf.ClientID, body: rb***REMOVED***
	buf, err := encode(req, b.conf.MetricRegistry)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = b.conn.SetWriteDeadline(time.Now().Add(b.conf.Net.WriteTimeout))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	requestTime := time.Now()
	bytes, err := b.conn.Write(buf)
	b.updateOutgoingCommunicationMetrics(bytes)
	if err != nil ***REMOVED***
		Logger.Printf("Failed to send SASL handshake %s: %s\n", b.addr, err.Error())
		return err
	***REMOVED***
	b.correlationID++
	//wait for the response
	header := make([]byte, 8) // response header
	_, err = io.ReadFull(b.conn, header)
	if err != nil ***REMOVED***
		Logger.Printf("Failed to read SASL handshake header : %s\n", err.Error())
		return err
	***REMOVED***
	length := binary.BigEndian.Uint32(header[:4])
	payload := make([]byte, length-4)
	n, err := io.ReadFull(b.conn, payload)
	if err != nil ***REMOVED***
		Logger.Printf("Failed to read SASL handshake payload : %s\n", err.Error())
		return err
	***REMOVED***
	b.updateIncomingCommunicationMetrics(n+8, time.Since(requestTime))
	res := &SaslHandshakeResponse***REMOVED******REMOVED***
	err = versionedDecode(payload, res, 0)
	if err != nil ***REMOVED***
		Logger.Printf("Failed to parse SASL handshake : %s\n", err.Error())
		return err
	***REMOVED***
	if res.Err != ErrNoError ***REMOVED***
		Logger.Printf("Invalid SASL Mechanism : %s\n", res.Err.Error())
		return res.Err
	***REMOVED***
	Logger.Print("Successful SASL handshake")
	return nil
***REMOVED***

// Kafka 0.10.0 plans to support SASL Plain and Kerberos as per PR #812 (KIP-43)/(JIRA KAFKA-3149)
// Some hosted kafka services such as IBM Message Hub already offer SASL/PLAIN auth with Kafka 0.9
//
// In SASL Plain, Kafka expects the auth header to be in the following format
// Message format (from https://tools.ietf.org/html/rfc4616):
//
//   message   = [authzid] UTF8NUL authcid UTF8NUL passwd
//   authcid   = 1*SAFE ; MUST accept up to 255 octets
//   authzid   = 1*SAFE ; MUST accept up to 255 octets
//   passwd    = 1*SAFE ; MUST accept up to 255 octets
//   UTF8NUL   = %x00 ; UTF-8 encoded NUL character
//
//   SAFE      = UTF1 / UTF2 / UTF3 / UTF4
//                  ;; any UTF-8 encoded Unicode character except NUL
//
// When credentials are valid, Kafka returns a 4 byte array of null characters.
// When credentials are invalid, Kafka closes the connection. This does not seem to be the ideal way
// of responding to bad credentials but thats how its being done today.
func (b *Broker) sendAndReceiveSASLPlainAuth() error ***REMOVED***
	if b.conf.Net.SASL.Handshake ***REMOVED***
		handshakeErr := b.sendAndReceiveSASLPlainHandshake()
		if handshakeErr != nil ***REMOVED***
			Logger.Printf("Error while performing SASL handshake %s\n", b.addr)
			return handshakeErr
		***REMOVED***
	***REMOVED***
	length := 1 + len(b.conf.Net.SASL.User) + 1 + len(b.conf.Net.SASL.Password)
	authBytes := make([]byte, length+4) //4 byte length header + auth data
	binary.BigEndian.PutUint32(authBytes, uint32(length))
	copy(authBytes[4:], []byte("\x00"+b.conf.Net.SASL.User+"\x00"+b.conf.Net.SASL.Password))

	err := b.conn.SetWriteDeadline(time.Now().Add(b.conf.Net.WriteTimeout))
	if err != nil ***REMOVED***
		Logger.Printf("Failed to set write deadline when doing SASL auth with broker %s: %s\n", b.addr, err.Error())
		return err
	***REMOVED***

	requestTime := time.Now()
	bytesWritten, err := b.conn.Write(authBytes)
	b.updateOutgoingCommunicationMetrics(bytesWritten)
	if err != nil ***REMOVED***
		Logger.Printf("Failed to write SASL auth header to broker %s: %s\n", b.addr, err.Error())
		return err
	***REMOVED***

	header := make([]byte, 4)
	n, err := io.ReadFull(b.conn, header)
	b.updateIncomingCommunicationMetrics(n, time.Since(requestTime))
	// If the credentials are valid, we would get a 4 byte response filled with null characters.
	// Otherwise, the broker closes the connection and we get an EOF
	if err != nil ***REMOVED***
		Logger.Printf("Failed to read response while authenticating with SASL to broker %s: %s\n", b.addr, err.Error())
		return err
	***REMOVED***

	Logger.Printf("SASL authentication successful with broker %s:%v - %v\n", b.addr, n, header)
	return nil
***REMOVED***

func (b *Broker) updateIncomingCommunicationMetrics(bytes int, requestLatency time.Duration) ***REMOVED***
	b.updateRequestLatencyMetrics(requestLatency)
	b.responseRate.Mark(1)
	if b.brokerResponseRate != nil ***REMOVED***
		b.brokerResponseRate.Mark(1)
	***REMOVED***
	responseSize := int64(bytes)
	b.incomingByteRate.Mark(responseSize)
	if b.brokerIncomingByteRate != nil ***REMOVED***
		b.brokerIncomingByteRate.Mark(responseSize)
	***REMOVED***
	b.responseSize.Update(responseSize)
	if b.brokerResponseSize != nil ***REMOVED***
		b.brokerResponseSize.Update(responseSize)
	***REMOVED***
***REMOVED***

func (b *Broker) updateRequestLatencyMetrics(requestLatency time.Duration) ***REMOVED***
	requestLatencyInMs := int64(requestLatency / time.Millisecond)
	b.requestLatency.Update(requestLatencyInMs)
	if b.brokerRequestLatency != nil ***REMOVED***
		b.brokerRequestLatency.Update(requestLatencyInMs)
	***REMOVED***
***REMOVED***

func (b *Broker) updateOutgoingCommunicationMetrics(bytes int) ***REMOVED***
	b.requestRate.Mark(1)
	if b.brokerRequestRate != nil ***REMOVED***
		b.brokerRequestRate.Mark(1)
	***REMOVED***
	requestSize := int64(bytes)
	b.outgoingByteRate.Mark(requestSize)
	if b.brokerOutgoingByteRate != nil ***REMOVED***
		b.brokerOutgoingByteRate.Mark(requestSize)
	***REMOVED***
	b.requestSize.Update(requestSize)
	if b.brokerRequestSize != nil ***REMOVED***
		b.brokerRequestSize.Update(requestSize)
	***REMOVED***
***REMOVED***
