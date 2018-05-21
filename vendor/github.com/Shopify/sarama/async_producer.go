package sarama

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/eapache/queue"
)

// AsyncProducer publishes Kafka messages using a non-blocking API. It routes messages
// to the correct broker for the provided topic-partition, refreshing metadata as appropriate,
// and parses responses for errors. You must read from the Errors() channel or the
// producer will deadlock. You must call Close() or AsyncClose() on a producer to avoid
// leaks: it will not be garbage-collected automatically when it passes out of
// scope.
type AsyncProducer interface ***REMOVED***

	// AsyncClose triggers a shutdown of the producer. The shutdown has completed
	// when both the Errors and Successes channels have been closed. When calling
	// AsyncClose, you *must* continue to read from those channels in order to
	// drain the results of any messages in flight.
	AsyncClose()

	// Close shuts down the producer and waits for any buffered messages to be
	// flushed. You must call this function before a producer object passes out of
	// scope, as it may otherwise leak memory. You must call this before calling
	// Close on the underlying client.
	Close() error

	// Input is the input channel for the user to write messages to that they
	// wish to send.
	Input() chan<- *ProducerMessage

	// Successes is the success output channel back to the user when Return.Successes is
	// enabled. If Return.Successes is true, you MUST read from this channel or the
	// Producer will deadlock. It is suggested that you send and read messages
	// together in a single select statement.
	Successes() <-chan *ProducerMessage

	// Errors is the error output channel back to the user. You MUST read from this
	// channel or the Producer will deadlock when the channel is full. Alternatively,
	// you can set Producer.Return.Errors in your config to false, which prevents
	// errors to be returned.
	Errors() <-chan *ProducerError
***REMOVED***

type asyncProducer struct ***REMOVED***
	client    Client
	conf      *Config
	ownClient bool

	errors                    chan *ProducerError
	input, successes, retries chan *ProducerMessage
	inFlight                  sync.WaitGroup

	brokers    map[*Broker]chan<- *ProducerMessage
	brokerRefs map[chan<- *ProducerMessage]int
	brokerLock sync.Mutex
***REMOVED***

// NewAsyncProducer creates a new AsyncProducer using the given broker addresses and configuration.
func NewAsyncProducer(addrs []string, conf *Config) (AsyncProducer, error) ***REMOVED***
	client, err := NewClient(addrs, conf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p, err := NewAsyncProducerFromClient(client)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	p.(*asyncProducer).ownClient = true
	return p, nil
***REMOVED***

// NewAsyncProducerFromClient creates a new Producer using the given client. It is still
// necessary to call Close() on the underlying client when shutting down this producer.
func NewAsyncProducerFromClient(client Client) (AsyncProducer, error) ***REMOVED***
	// Check that we are not dealing with a closed Client before processing any other arguments
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	p := &asyncProducer***REMOVED***
		client:     client,
		conf:       client.Config(),
		errors:     make(chan *ProducerError),
		input:      make(chan *ProducerMessage),
		successes:  make(chan *ProducerMessage),
		retries:    make(chan *ProducerMessage),
		brokers:    make(map[*Broker]chan<- *ProducerMessage),
		brokerRefs: make(map[chan<- *ProducerMessage]int),
	***REMOVED***

	// launch our singleton dispatchers
	go withRecover(p.dispatcher)
	go withRecover(p.retryHandler)

	return p, nil
***REMOVED***

type flagSet int8

const (
	syn      flagSet = 1 << iota // first message from partitionProducer to brokerProducer
	fin                          // final message from partitionProducer to brokerProducer and back
	shutdown                     // start the shutdown process
)

// ProducerMessage is the collection of elements passed to the Producer in order to send a message.
type ProducerMessage struct ***REMOVED***
	Topic string // The Kafka topic for this message.
	// The partitioning key for this message. Pre-existing Encoders include
	// StringEncoder and ByteEncoder.
	Key Encoder
	// The actual message to store in Kafka. Pre-existing Encoders include
	// StringEncoder and ByteEncoder.
	Value Encoder

	// The headers are key-value pairs that are transparently passed
	// by Kafka between producers and consumers.
	Headers []RecordHeader

	// This field is used to hold arbitrary data you wish to include so it
	// will be available when receiving on the Successes and Errors channels.
	// Sarama completely ignores this field and is only to be used for
	// pass-through data.
	Metadata interface***REMOVED******REMOVED***

	// Below this point are filled in by the producer as the message is processed

	// Offset is the offset of the message stored on the broker. This is only
	// guaranteed to be defined if the message was successfully delivered and
	// RequiredAcks is not NoResponse.
	Offset int64
	// Partition is the partition that the message was sent to. This is only
	// guaranteed to be defined if the message was successfully delivered.
	Partition int32
	// Timestamp is the timestamp assigned to the message by the broker. This
	// is only guaranteed to be defined if the message was successfully
	// delivered, RequiredAcks is not NoResponse, and the Kafka broker is at
	// least version 0.10.0.
	Timestamp time.Time

	retries int
	flags   flagSet
***REMOVED***

const producerMessageOverhead = 26 // the metadata overhead of CRC, flags, etc.

func (m *ProducerMessage) byteSize(version int) int ***REMOVED***
	var size int
	if version >= 2 ***REMOVED***
		size = maximumRecordOverhead
		for _, h := range m.Headers ***REMOVED***
			size += len(h.Key) + len(h.Value) + 2*binary.MaxVarintLen32
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		size = producerMessageOverhead
	***REMOVED***
	if m.Key != nil ***REMOVED***
		size += m.Key.Length()
	***REMOVED***
	if m.Value != nil ***REMOVED***
		size += m.Value.Length()
	***REMOVED***
	return size
***REMOVED***

func (m *ProducerMessage) clear() ***REMOVED***
	m.flags = 0
	m.retries = 0
***REMOVED***

// ProducerError is the type of error generated when the producer fails to deliver a message.
// It contains the original ProducerMessage as well as the actual error value.
type ProducerError struct ***REMOVED***
	Msg *ProducerMessage
	Err error
***REMOVED***

func (pe ProducerError) Error() string ***REMOVED***
	return fmt.Sprintf("kafka: Failed to produce message to topic %s: %s", pe.Msg.Topic, pe.Err)
***REMOVED***

// ProducerErrors is a type that wraps a batch of "ProducerError"s and implements the Error interface.
// It can be returned from the Producer's Close method to avoid the need to manually drain the Errors channel
// when closing a producer.
type ProducerErrors []*ProducerError

func (pe ProducerErrors) Error() string ***REMOVED***
	return fmt.Sprintf("kafka: Failed to deliver %d messages.", len(pe))
***REMOVED***

func (p *asyncProducer) Errors() <-chan *ProducerError ***REMOVED***
	return p.errors
***REMOVED***

func (p *asyncProducer) Successes() <-chan *ProducerMessage ***REMOVED***
	return p.successes
***REMOVED***

func (p *asyncProducer) Input() chan<- *ProducerMessage ***REMOVED***
	return p.input
***REMOVED***

func (p *asyncProducer) Close() error ***REMOVED***
	p.AsyncClose()

	if p.conf.Producer.Return.Successes ***REMOVED***
		go withRecover(func() ***REMOVED***
			for range p.successes ***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	var errors ProducerErrors
	if p.conf.Producer.Return.Errors ***REMOVED***
		for event := range p.errors ***REMOVED***
			errors = append(errors, event)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		<-p.errors
	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		return errors
	***REMOVED***
	return nil
***REMOVED***

func (p *asyncProducer) AsyncClose() ***REMOVED***
	go withRecover(p.shutdown)
***REMOVED***

// singleton
// dispatches messages by topic
func (p *asyncProducer) dispatcher() ***REMOVED***
	handlers := make(map[string]chan<- *ProducerMessage)
	shuttingDown := false

	for msg := range p.input ***REMOVED***
		if msg == nil ***REMOVED***
			Logger.Println("Something tried to send a nil message, it was ignored.")
			continue
		***REMOVED***

		if msg.flags&shutdown != 0 ***REMOVED***
			shuttingDown = true
			p.inFlight.Done()
			continue
		***REMOVED*** else if msg.retries == 0 ***REMOVED***
			if shuttingDown ***REMOVED***
				// we can't just call returnError here because that decrements the wait group,
				// which hasn't been incremented yet for this message, and shouldn't be
				pErr := &ProducerError***REMOVED***Msg: msg, Err: ErrShuttingDown***REMOVED***
				if p.conf.Producer.Return.Errors ***REMOVED***
					p.errors <- pErr
				***REMOVED*** else ***REMOVED***
					Logger.Println(pErr)
				***REMOVED***
				continue
			***REMOVED***
			p.inFlight.Add(1)
		***REMOVED***

		version := 1
		if p.conf.Version.IsAtLeast(V0_11_0_0) ***REMOVED***
			version = 2
		***REMOVED***
		if msg.byteSize(version) > p.conf.Producer.MaxMessageBytes ***REMOVED***
			p.returnError(msg, ErrMessageSizeTooLarge)
			continue
		***REMOVED***

		handler := handlers[msg.Topic]
		if handler == nil ***REMOVED***
			handler = p.newTopicProducer(msg.Topic)
			handlers[msg.Topic] = handler
		***REMOVED***

		handler <- msg
	***REMOVED***

	for _, handler := range handlers ***REMOVED***
		close(handler)
	***REMOVED***
***REMOVED***

// one per topic
// partitions messages, then dispatches them by partition
type topicProducer struct ***REMOVED***
	parent *asyncProducer
	topic  string
	input  <-chan *ProducerMessage

	breaker     *breaker.Breaker
	handlers    map[int32]chan<- *ProducerMessage
	partitioner Partitioner
***REMOVED***

func (p *asyncProducer) newTopicProducer(topic string) chan<- *ProducerMessage ***REMOVED***
	input := make(chan *ProducerMessage, p.conf.ChannelBufferSize)
	tp := &topicProducer***REMOVED***
		parent:      p,
		topic:       topic,
		input:       input,
		breaker:     breaker.New(3, 1, 10*time.Second),
		handlers:    make(map[int32]chan<- *ProducerMessage),
		partitioner: p.conf.Producer.Partitioner(topic),
	***REMOVED***
	go withRecover(tp.dispatch)
	return input
***REMOVED***

func (tp *topicProducer) dispatch() ***REMOVED***
	for msg := range tp.input ***REMOVED***
		if msg.retries == 0 ***REMOVED***
			if err := tp.partitionMessage(msg); err != nil ***REMOVED***
				tp.parent.returnError(msg, err)
				continue
			***REMOVED***
		***REMOVED***

		handler := tp.handlers[msg.Partition]
		if handler == nil ***REMOVED***
			handler = tp.parent.newPartitionProducer(msg.Topic, msg.Partition)
			tp.handlers[msg.Partition] = handler
		***REMOVED***

		handler <- msg
	***REMOVED***

	for _, handler := range tp.handlers ***REMOVED***
		close(handler)
	***REMOVED***
***REMOVED***

func (tp *topicProducer) partitionMessage(msg *ProducerMessage) error ***REMOVED***
	var partitions []int32

	err := tp.breaker.Run(func() (err error) ***REMOVED***
		if tp.partitioner.RequiresConsistency() ***REMOVED***
			partitions, err = tp.parent.client.Partitions(msg.Topic)
		***REMOVED*** else ***REMOVED***
			partitions, err = tp.parent.client.WritablePartitions(msg.Topic)
		***REMOVED***
		return
	***REMOVED***)

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	numPartitions := int32(len(partitions))

	if numPartitions == 0 ***REMOVED***
		return ErrLeaderNotAvailable
	***REMOVED***

	choice, err := tp.partitioner.Partition(msg, numPartitions)

	if err != nil ***REMOVED***
		return err
	***REMOVED*** else if choice < 0 || choice >= numPartitions ***REMOVED***
		return ErrInvalidPartition
	***REMOVED***

	msg.Partition = partitions[choice]

	return nil
***REMOVED***

// one per partition per topic
// dispatches messages to the appropriate broker
// also responsible for maintaining message order during retries
type partitionProducer struct ***REMOVED***
	parent    *asyncProducer
	topic     string
	partition int32
	input     <-chan *ProducerMessage

	leader  *Broker
	breaker *breaker.Breaker
	output  chan<- *ProducerMessage

	// highWatermark tracks the "current" retry level, which is the only one where we actually let messages through,
	// all other messages get buffered in retryState[msg.retries].buf to preserve ordering
	// retryState[msg.retries].expectChaser simply tracks whether we've seen a fin message for a given level (and
	// therefore whether our buffer is complete and safe to flush)
	highWatermark int
	retryState    []partitionRetryState
***REMOVED***

type partitionRetryState struct ***REMOVED***
	buf          []*ProducerMessage
	expectChaser bool
***REMOVED***

func (p *asyncProducer) newPartitionProducer(topic string, partition int32) chan<- *ProducerMessage ***REMOVED***
	input := make(chan *ProducerMessage, p.conf.ChannelBufferSize)
	pp := &partitionProducer***REMOVED***
		parent:    p,
		topic:     topic,
		partition: partition,
		input:     input,

		breaker:    breaker.New(3, 1, 10*time.Second),
		retryState: make([]partitionRetryState, p.conf.Producer.Retry.Max+1),
	***REMOVED***
	go withRecover(pp.dispatch)
	return input
***REMOVED***

func (pp *partitionProducer) dispatch() ***REMOVED***
	// try to prefetch the leader; if this doesn't work, we'll do a proper call to `updateLeader`
	// on the first message
	pp.leader, _ = pp.parent.client.Leader(pp.topic, pp.partition)
	if pp.leader != nil ***REMOVED***
		pp.output = pp.parent.getBrokerProducer(pp.leader)
		pp.parent.inFlight.Add(1) // we're generating a syn message; track it so we don't shut down while it's still inflight
		pp.output <- &ProducerMessage***REMOVED***Topic: pp.topic, Partition: pp.partition, flags: syn***REMOVED***
	***REMOVED***

	for msg := range pp.input ***REMOVED***
		if msg.retries > pp.highWatermark ***REMOVED***
			// a new, higher, retry level; handle it and then back off
			pp.newHighWatermark(msg.retries)
			time.Sleep(pp.parent.conf.Producer.Retry.Backoff)
		***REMOVED*** else if pp.highWatermark > 0 ***REMOVED***
			// we are retrying something (else highWatermark would be 0) but this message is not a *new* retry level
			if msg.retries < pp.highWatermark ***REMOVED***
				// in fact this message is not even the current retry level, so buffer it for now (unless it's a just a fin)
				if msg.flags&fin == fin ***REMOVED***
					pp.retryState[msg.retries].expectChaser = false
					pp.parent.inFlight.Done() // this fin is now handled and will be garbage collected
				***REMOVED*** else ***REMOVED***
					pp.retryState[msg.retries].buf = append(pp.retryState[msg.retries].buf, msg)
				***REMOVED***
				continue
			***REMOVED*** else if msg.flags&fin == fin ***REMOVED***
				// this message is of the current retry level (msg.retries == highWatermark) and the fin flag is set,
				// meaning this retry level is done and we can go down (at least) one level and flush that
				pp.retryState[pp.highWatermark].expectChaser = false
				pp.flushRetryBuffers()
				pp.parent.inFlight.Done() // this fin is now handled and will be garbage collected
				continue
			***REMOVED***
		***REMOVED***

		// if we made it this far then the current msg contains real data, and can be sent to the next goroutine
		// without breaking any of our ordering guarantees

		if pp.output == nil ***REMOVED***
			if err := pp.updateLeader(); err != nil ***REMOVED***
				pp.parent.returnError(msg, err)
				time.Sleep(pp.parent.conf.Producer.Retry.Backoff)
				continue
			***REMOVED***
			Logger.Printf("producer/leader/%s/%d selected broker %d\n", pp.topic, pp.partition, pp.leader.ID())
		***REMOVED***

		pp.output <- msg
	***REMOVED***

	if pp.output != nil ***REMOVED***
		pp.parent.unrefBrokerProducer(pp.leader, pp.output)
	***REMOVED***
***REMOVED***

func (pp *partitionProducer) newHighWatermark(hwm int) ***REMOVED***
	Logger.Printf("producer/leader/%s/%d state change to [retrying-%d]\n", pp.topic, pp.partition, hwm)
	pp.highWatermark = hwm

	// send off a fin so that we know when everything "in between" has made it
	// back to us and we can safely flush the backlog (otherwise we risk re-ordering messages)
	pp.retryState[pp.highWatermark].expectChaser = true
	pp.parent.inFlight.Add(1) // we're generating a fin message; track it so we don't shut down while it's still inflight
	pp.output <- &ProducerMessage***REMOVED***Topic: pp.topic, Partition: pp.partition, flags: fin, retries: pp.highWatermark - 1***REMOVED***

	// a new HWM means that our current broker selection is out of date
	Logger.Printf("producer/leader/%s/%d abandoning broker %d\n", pp.topic, pp.partition, pp.leader.ID())
	pp.parent.unrefBrokerProducer(pp.leader, pp.output)
	pp.output = nil
***REMOVED***

func (pp *partitionProducer) flushRetryBuffers() ***REMOVED***
	Logger.Printf("producer/leader/%s/%d state change to [flushing-%d]\n", pp.topic, pp.partition, pp.highWatermark)
	for ***REMOVED***
		pp.highWatermark--

		if pp.output == nil ***REMOVED***
			if err := pp.updateLeader(); err != nil ***REMOVED***
				pp.parent.returnErrors(pp.retryState[pp.highWatermark].buf, err)
				goto flushDone
			***REMOVED***
			Logger.Printf("producer/leader/%s/%d selected broker %d\n", pp.topic, pp.partition, pp.leader.ID())
		***REMOVED***

		for _, msg := range pp.retryState[pp.highWatermark].buf ***REMOVED***
			pp.output <- msg
		***REMOVED***

	flushDone:
		pp.retryState[pp.highWatermark].buf = nil
		if pp.retryState[pp.highWatermark].expectChaser ***REMOVED***
			Logger.Printf("producer/leader/%s/%d state change to [retrying-%d]\n", pp.topic, pp.partition, pp.highWatermark)
			break
		***REMOVED*** else if pp.highWatermark == 0 ***REMOVED***
			Logger.Printf("producer/leader/%s/%d state change to [normal]\n", pp.topic, pp.partition)
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (pp *partitionProducer) updateLeader() error ***REMOVED***
	return pp.breaker.Run(func() (err error) ***REMOVED***
		if err = pp.parent.client.RefreshMetadata(pp.topic); err != nil ***REMOVED***
			return err
		***REMOVED***

		if pp.leader, err = pp.parent.client.Leader(pp.topic, pp.partition); err != nil ***REMOVED***
			return err
		***REMOVED***

		pp.output = pp.parent.getBrokerProducer(pp.leader)
		pp.parent.inFlight.Add(1) // we're generating a syn message; track it so we don't shut down while it's still inflight
		pp.output <- &ProducerMessage***REMOVED***Topic: pp.topic, Partition: pp.partition, flags: syn***REMOVED***

		return nil
	***REMOVED***)
***REMOVED***

// one per broker; also constructs an associated flusher
func (p *asyncProducer) newBrokerProducer(broker *Broker) chan<- *ProducerMessage ***REMOVED***
	var (
		input     = make(chan *ProducerMessage)
		bridge    = make(chan *produceSet)
		responses = make(chan *brokerProducerResponse)
	)

	bp := &brokerProducer***REMOVED***
		parent:         p,
		broker:         broker,
		input:          input,
		output:         bridge,
		responses:      responses,
		buffer:         newProduceSet(p),
		currentRetries: make(map[string]map[int32]error),
	***REMOVED***
	go withRecover(bp.run)

	// minimal bridge to make the network response `select`able
	go withRecover(func() ***REMOVED***
		for set := range bridge ***REMOVED***
			request := set.buildRequest()

			response, err := broker.Produce(request)

			responses <- &brokerProducerResponse***REMOVED***
				set: set,
				err: err,
				res: response,
			***REMOVED***
		***REMOVED***
		close(responses)
	***REMOVED***)

	return input
***REMOVED***

type brokerProducerResponse struct ***REMOVED***
	set *produceSet
	err error
	res *ProduceResponse
***REMOVED***

// groups messages together into appropriately-sized batches for sending to the broker
// handles state related to retries etc
type brokerProducer struct ***REMOVED***
	parent *asyncProducer
	broker *Broker

	input     <-chan *ProducerMessage
	output    chan<- *produceSet
	responses <-chan *brokerProducerResponse

	buffer     *produceSet
	timer      <-chan time.Time
	timerFired bool

	closing        error
	currentRetries map[string]map[int32]error
***REMOVED***

func (bp *brokerProducer) run() ***REMOVED***
	var output chan<- *produceSet
	Logger.Printf("producer/broker/%d starting up\n", bp.broker.ID())

	for ***REMOVED***
		select ***REMOVED***
		case msg := <-bp.input:
			if msg == nil ***REMOVED***
				bp.shutdown()
				return
			***REMOVED***

			if msg.flags&syn == syn ***REMOVED***
				Logger.Printf("producer/broker/%d state change to [open] on %s/%d\n",
					bp.broker.ID(), msg.Topic, msg.Partition)
				if bp.currentRetries[msg.Topic] == nil ***REMOVED***
					bp.currentRetries[msg.Topic] = make(map[int32]error)
				***REMOVED***
				bp.currentRetries[msg.Topic][msg.Partition] = nil
				bp.parent.inFlight.Done()
				continue
			***REMOVED***

			if reason := bp.needsRetry(msg); reason != nil ***REMOVED***
				bp.parent.retryMessage(msg, reason)

				if bp.closing == nil && msg.flags&fin == fin ***REMOVED***
					// we were retrying this partition but we can start processing again
					delete(bp.currentRetries[msg.Topic], msg.Partition)
					Logger.Printf("producer/broker/%d state change to [closed] on %s/%d\n",
						bp.broker.ID(), msg.Topic, msg.Partition)
				***REMOVED***

				continue
			***REMOVED***

			if bp.buffer.wouldOverflow(msg) ***REMOVED***
				if err := bp.waitForSpace(msg); err != nil ***REMOVED***
					bp.parent.retryMessage(msg, err)
					continue
				***REMOVED***
			***REMOVED***

			if err := bp.buffer.add(msg); err != nil ***REMOVED***
				bp.parent.returnError(msg, err)
				continue
			***REMOVED***

			if bp.parent.conf.Producer.Flush.Frequency > 0 && bp.timer == nil ***REMOVED***
				bp.timer = time.After(bp.parent.conf.Producer.Flush.Frequency)
			***REMOVED***
		case <-bp.timer:
			bp.timerFired = true
		case output <- bp.buffer:
			bp.rollOver()
		case response := <-bp.responses:
			bp.handleResponse(response)
		***REMOVED***

		if bp.timerFired || bp.buffer.readyToFlush() ***REMOVED***
			output = bp.output
		***REMOVED*** else ***REMOVED***
			output = nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (bp *brokerProducer) shutdown() ***REMOVED***
	for !bp.buffer.empty() ***REMOVED***
		select ***REMOVED***
		case response := <-bp.responses:
			bp.handleResponse(response)
		case bp.output <- bp.buffer:
			bp.rollOver()
		***REMOVED***
	***REMOVED***
	close(bp.output)
	for response := range bp.responses ***REMOVED***
		bp.handleResponse(response)
	***REMOVED***

	Logger.Printf("producer/broker/%d shut down\n", bp.broker.ID())
***REMOVED***

func (bp *brokerProducer) needsRetry(msg *ProducerMessage) error ***REMOVED***
	if bp.closing != nil ***REMOVED***
		return bp.closing
	***REMOVED***

	return bp.currentRetries[msg.Topic][msg.Partition]
***REMOVED***

func (bp *brokerProducer) waitForSpace(msg *ProducerMessage) error ***REMOVED***
	Logger.Printf("producer/broker/%d maximum request accumulated, waiting for space\n", bp.broker.ID())

	for ***REMOVED***
		select ***REMOVED***
		case response := <-bp.responses:
			bp.handleResponse(response)
			// handling a response can change our state, so re-check some things
			if reason := bp.needsRetry(msg); reason != nil ***REMOVED***
				return reason
			***REMOVED*** else if !bp.buffer.wouldOverflow(msg) ***REMOVED***
				return nil
			***REMOVED***
		case bp.output <- bp.buffer:
			bp.rollOver()
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (bp *brokerProducer) rollOver() ***REMOVED***
	bp.timer = nil
	bp.timerFired = false
	bp.buffer = newProduceSet(bp.parent)
***REMOVED***

func (bp *brokerProducer) handleResponse(response *brokerProducerResponse) ***REMOVED***
	if response.err != nil ***REMOVED***
		bp.handleError(response.set, response.err)
	***REMOVED*** else ***REMOVED***
		bp.handleSuccess(response.set, response.res)
	***REMOVED***

	if bp.buffer.empty() ***REMOVED***
		bp.rollOver() // this can happen if the response invalidated our buffer
	***REMOVED***
***REMOVED***

func (bp *brokerProducer) handleSuccess(sent *produceSet, response *ProduceResponse) ***REMOVED***
	// we iterate through the blocks in the request set, not the response, so that we notice
	// if the response is missing a block completely
	sent.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) ***REMOVED***
		if response == nil ***REMOVED***
			// this only happens when RequiredAcks is NoResponse, so we have to assume success
			bp.parent.returnSuccesses(msgs)
			return
		***REMOVED***

		block := response.GetBlock(topic, partition)
		if block == nil ***REMOVED***
			bp.parent.returnErrors(msgs, ErrIncompleteResponse)
			return
		***REMOVED***

		switch block.Err ***REMOVED***
		// Success
		case ErrNoError:
			if bp.parent.conf.Version.IsAtLeast(V0_10_0_0) && !block.Timestamp.IsZero() ***REMOVED***
				for _, msg := range msgs ***REMOVED***
					msg.Timestamp = block.Timestamp
				***REMOVED***
			***REMOVED***
			for i, msg := range msgs ***REMOVED***
				msg.Offset = block.Offset + int64(i)
			***REMOVED***
			bp.parent.returnSuccesses(msgs)
		// Retriable errors
		case ErrInvalidMessage, ErrUnknownTopicOrPartition, ErrLeaderNotAvailable, ErrNotLeaderForPartition,
			ErrRequestTimedOut, ErrNotEnoughReplicas, ErrNotEnoughReplicasAfterAppend:
			Logger.Printf("producer/broker/%d state change to [retrying] on %s/%d because %v\n",
				bp.broker.ID(), topic, partition, block.Err)
			bp.currentRetries[topic][partition] = block.Err
			bp.parent.retryMessages(msgs, block.Err)
			bp.parent.retryMessages(bp.buffer.dropPartition(topic, partition), block.Err)
		// Other non-retriable errors
		default:
			bp.parent.returnErrors(msgs, block.Err)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (bp *brokerProducer) handleError(sent *produceSet, err error) ***REMOVED***
	switch err.(type) ***REMOVED***
	case PacketEncodingError:
		sent.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) ***REMOVED***
			bp.parent.returnErrors(msgs, err)
		***REMOVED***)
	default:
		Logger.Printf("producer/broker/%d state change to [closing] because %s\n", bp.broker.ID(), err)
		bp.parent.abandonBrokerConnection(bp.broker)
		_ = bp.broker.Close()
		bp.closing = err
		sent.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) ***REMOVED***
			bp.parent.retryMessages(msgs, err)
		***REMOVED***)
		bp.buffer.eachPartition(func(topic string, partition int32, msgs []*ProducerMessage) ***REMOVED***
			bp.parent.retryMessages(msgs, err)
		***REMOVED***)
		bp.rollOver()
	***REMOVED***
***REMOVED***

// singleton
// effectively a "bridge" between the flushers and the dispatcher in order to avoid deadlock
// based on https://godoc.org/github.com/eapache/channels#InfiniteChannel
func (p *asyncProducer) retryHandler() ***REMOVED***
	var msg *ProducerMessage
	buf := queue.New()

	for ***REMOVED***
		if buf.Length() == 0 ***REMOVED***
			msg = <-p.retries
		***REMOVED*** else ***REMOVED***
			select ***REMOVED***
			case msg = <-p.retries:
			case p.input <- buf.Peek().(*ProducerMessage):
				buf.Remove()
				continue
			***REMOVED***
		***REMOVED***

		if msg == nil ***REMOVED***
			return
		***REMOVED***

		buf.Add(msg)
	***REMOVED***
***REMOVED***

// utility functions

func (p *asyncProducer) shutdown() ***REMOVED***
	Logger.Println("Producer shutting down.")
	p.inFlight.Add(1)
	p.input <- &ProducerMessage***REMOVED***flags: shutdown***REMOVED***

	p.inFlight.Wait()

	if p.ownClient ***REMOVED***
		err := p.client.Close()
		if err != nil ***REMOVED***
			Logger.Println("producer/shutdown failed to close the embedded client:", err)
		***REMOVED***
	***REMOVED***

	close(p.input)
	close(p.retries)
	close(p.errors)
	close(p.successes)
***REMOVED***

func (p *asyncProducer) returnError(msg *ProducerMessage, err error) ***REMOVED***
	msg.clear()
	pErr := &ProducerError***REMOVED***Msg: msg, Err: err***REMOVED***
	if p.conf.Producer.Return.Errors ***REMOVED***
		p.errors <- pErr
	***REMOVED*** else ***REMOVED***
		Logger.Println(pErr)
	***REMOVED***
	p.inFlight.Done()
***REMOVED***

func (p *asyncProducer) returnErrors(batch []*ProducerMessage, err error) ***REMOVED***
	for _, msg := range batch ***REMOVED***
		p.returnError(msg, err)
	***REMOVED***
***REMOVED***

func (p *asyncProducer) returnSuccesses(batch []*ProducerMessage) ***REMOVED***
	for _, msg := range batch ***REMOVED***
		if p.conf.Producer.Return.Successes ***REMOVED***
			msg.clear()
			p.successes <- msg
		***REMOVED***
		p.inFlight.Done()
	***REMOVED***
***REMOVED***

func (p *asyncProducer) retryMessage(msg *ProducerMessage, err error) ***REMOVED***
	if msg.retries >= p.conf.Producer.Retry.Max ***REMOVED***
		p.returnError(msg, err)
	***REMOVED*** else ***REMOVED***
		msg.retries++
		p.retries <- msg
	***REMOVED***
***REMOVED***

func (p *asyncProducer) retryMessages(batch []*ProducerMessage, err error) ***REMOVED***
	for _, msg := range batch ***REMOVED***
		p.retryMessage(msg, err)
	***REMOVED***
***REMOVED***

func (p *asyncProducer) getBrokerProducer(broker *Broker) chan<- *ProducerMessage ***REMOVED***
	p.brokerLock.Lock()
	defer p.brokerLock.Unlock()

	bp := p.brokers[broker]

	if bp == nil ***REMOVED***
		bp = p.newBrokerProducer(broker)
		p.brokers[broker] = bp
		p.brokerRefs[bp] = 0
	***REMOVED***

	p.brokerRefs[bp]++

	return bp
***REMOVED***

func (p *asyncProducer) unrefBrokerProducer(broker *Broker, bp chan<- *ProducerMessage) ***REMOVED***
	p.brokerLock.Lock()
	defer p.brokerLock.Unlock()

	p.brokerRefs[bp]--
	if p.brokerRefs[bp] == 0 ***REMOVED***
		close(bp)
		delete(p.brokerRefs, bp)

		if p.brokers[broker] == bp ***REMOVED***
			delete(p.brokers, broker)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *asyncProducer) abandonBrokerConnection(broker *Broker) ***REMOVED***
	p.brokerLock.Lock()
	defer p.brokerLock.Unlock()

	delete(p.brokers, broker)
***REMOVED***
