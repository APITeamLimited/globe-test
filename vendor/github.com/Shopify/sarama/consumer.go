package sarama

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ConsumerMessage encapsulates a Kafka message returned by the consumer.
type ConsumerMessage struct ***REMOVED***
	Key, Value     []byte
	Topic          string
	Partition      int32
	Offset         int64
	Timestamp      time.Time       // only set if kafka is version 0.10+, inner message timestamp
	BlockTimestamp time.Time       // only set if kafka is version 0.10+, outer (compressed) block timestamp
	Headers        []*RecordHeader // only set if kafka is version 0.11+
***REMOVED***

// ConsumerError is what is provided to the user when an error occurs.
// It wraps an error and includes the topic and partition.
type ConsumerError struct ***REMOVED***
	Topic     string
	Partition int32
	Err       error
***REMOVED***

func (ce ConsumerError) Error() string ***REMOVED***
	return fmt.Sprintf("kafka: error while consuming %s/%d: %s", ce.Topic, ce.Partition, ce.Err)
***REMOVED***

// ConsumerErrors is a type that wraps a batch of errors and implements the Error interface.
// It can be returned from the PartitionConsumer's Close methods to avoid the need to manually drain errors
// when stopping.
type ConsumerErrors []*ConsumerError

func (ce ConsumerErrors) Error() string ***REMOVED***
	return fmt.Sprintf("kafka: %d errors while consuming", len(ce))
***REMOVED***

// Consumer manages PartitionConsumers which process Kafka messages from brokers. You MUST call Close()
// on a consumer to avoid leaks, it will not be garbage-collected automatically when it passes out of
// scope.
//
// Sarama's Consumer type does not currently support automatic consumer-group rebalancing and offset tracking.
// For Zookeeper-based tracking (Kafka 0.8.2 and earlier), the https://github.com/wvanbergen/kafka library
// builds on Sarama to add this support. For Kafka-based tracking (Kafka 0.9 and later), the
// https://github.com/bsm/sarama-cluster library builds on Sarama to add this support.
type Consumer interface ***REMOVED***

	// Topics returns the set of available topics as retrieved from the cluster
	// metadata. This method is the same as Client.Topics(), and is provided for
	// convenience.
	Topics() ([]string, error)

	// Partitions returns the sorted list of all partition IDs for the given topic.
	// This method is the same as Client.Partitions(), and is provided for convenience.
	Partitions(topic string) ([]int32, error)

	// ConsumePartition creates a PartitionConsumer on the given topic/partition with
	// the given offset. It will return an error if this Consumer is already consuming
	// on the given topic/partition. Offset can be a literal offset, or OffsetNewest
	// or OffsetOldest
	ConsumePartition(topic string, partition int32, offset int64) (PartitionConsumer, error)

	// HighWaterMarks returns the current high water marks for each topic and partition.
	// Consistency between partitions is not guaranteed since high water marks are updated separately.
	HighWaterMarks() map[string]map[int32]int64

	// Close shuts down the consumer. It must be called after all child
	// PartitionConsumers have already been closed.
	Close() error
***REMOVED***

type consumer struct ***REMOVED***
	client    Client
	conf      *Config
	ownClient bool

	lock            sync.Mutex
	children        map[string]map[int32]*partitionConsumer
	brokerConsumers map[*Broker]*brokerConsumer
***REMOVED***

// NewConsumer creates a new consumer using the given broker addresses and configuration.
func NewConsumer(addrs []string, config *Config) (Consumer, error) ***REMOVED***
	client, err := NewClient(addrs, config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c, err := NewConsumerFromClient(client)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.(*consumer).ownClient = true
	return c, nil
***REMOVED***

// NewConsumerFromClient creates a new consumer using the given client. It is still
// necessary to call Close() on the underlying client when shutting down this consumer.
func NewConsumerFromClient(client Client) (Consumer, error) ***REMOVED***
	// Check that we are not dealing with a closed Client before processing any other arguments
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	c := &consumer***REMOVED***
		client:          client,
		conf:            client.Config(),
		children:        make(map[string]map[int32]*partitionConsumer),
		brokerConsumers: make(map[*Broker]*brokerConsumer),
	***REMOVED***

	return c, nil
***REMOVED***

func (c *consumer) Close() error ***REMOVED***
	if c.ownClient ***REMOVED***
		return c.client.Close()
	***REMOVED***
	return nil
***REMOVED***

func (c *consumer) Topics() ([]string, error) ***REMOVED***
	return c.client.Topics()
***REMOVED***

func (c *consumer) Partitions(topic string) ([]int32, error) ***REMOVED***
	return c.client.Partitions(topic)
***REMOVED***

func (c *consumer) ConsumePartition(topic string, partition int32, offset int64) (PartitionConsumer, error) ***REMOVED***
	child := &partitionConsumer***REMOVED***
		consumer:  c,
		conf:      c.conf,
		topic:     topic,
		partition: partition,
		messages:  make(chan *ConsumerMessage, c.conf.ChannelBufferSize),
		errors:    make(chan *ConsumerError, c.conf.ChannelBufferSize),
		feeder:    make(chan *FetchResponse, 1),
		trigger:   make(chan none, 1),
		dying:     make(chan none),
		fetchSize: c.conf.Consumer.Fetch.Default,
	***REMOVED***

	if err := child.chooseStartingOffset(offset); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var leader *Broker
	var err error
	if leader, err = c.client.Leader(child.topic, child.partition); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := c.addChild(child); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	go withRecover(child.dispatcher)
	go withRecover(child.responseFeeder)

	child.broker = c.refBrokerConsumer(leader)
	child.broker.input <- child

	return child, nil
***REMOVED***

func (c *consumer) HighWaterMarks() map[string]map[int32]int64 ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	hwms := make(map[string]map[int32]int64)
	for topic, p := range c.children ***REMOVED***
		hwm := make(map[int32]int64, len(p))
		for partition, pc := range p ***REMOVED***
			hwm[partition] = pc.HighWaterMarkOffset()
		***REMOVED***
		hwms[topic] = hwm
	***REMOVED***

	return hwms
***REMOVED***

func (c *consumer) addChild(child *partitionConsumer) error ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	topicChildren := c.children[child.topic]
	if topicChildren == nil ***REMOVED***
		topicChildren = make(map[int32]*partitionConsumer)
		c.children[child.topic] = topicChildren
	***REMOVED***

	if topicChildren[child.partition] != nil ***REMOVED***
		return ConfigurationError("That topic/partition is already being consumed")
	***REMOVED***

	topicChildren[child.partition] = child
	return nil
***REMOVED***

func (c *consumer) removeChild(child *partitionConsumer) ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.children[child.topic], child.partition)
***REMOVED***

func (c *consumer) refBrokerConsumer(broker *Broker) *brokerConsumer ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	bc := c.brokerConsumers[broker]
	if bc == nil ***REMOVED***
		bc = c.newBrokerConsumer(broker)
		c.brokerConsumers[broker] = bc
	***REMOVED***

	bc.refs++

	return bc
***REMOVED***

func (c *consumer) unrefBrokerConsumer(brokerWorker *brokerConsumer) ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	brokerWorker.refs--

	if brokerWorker.refs == 0 ***REMOVED***
		close(brokerWorker.input)
		if c.brokerConsumers[brokerWorker.broker] == brokerWorker ***REMOVED***
			delete(c.brokerConsumers, brokerWorker.broker)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *consumer) abandonBrokerConsumer(brokerWorker *brokerConsumer) ***REMOVED***
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.brokerConsumers, brokerWorker.broker)
***REMOVED***

// PartitionConsumer

// PartitionConsumer processes Kafka messages from a given topic and partition. You MUST call one of Close() or
// AsyncClose() on a PartitionConsumer to avoid leaks; it will not be garbage-collected automatically when it passes out
// of scope.
//
// The simplest way of using a PartitionConsumer is to loop over its Messages channel using a for/range
// loop. The PartitionConsumer will only stop itself in one case: when the offset being consumed is reported
// as out of range by the brokers. In this case you should decide what you want to do (try a different offset,
// notify a human, etc) and handle it appropriately. For all other error cases, it will just keep retrying.
// By default, it logs these errors to sarama.Logger; if you want to be notified directly of all errors, set
// your config's Consumer.Return.Errors to true and read from the Errors channel, using a select statement
// or a separate goroutine. Check out the Consumer examples to see implementations of these different approaches.
//
// To terminate such a for/range loop while the loop is executing, call AsyncClose. This will kick off the process of
// consumer tear-down & return imediately. Continue to loop, servicing the Messages channel until the teardown process
// AsyncClose initiated closes it (thus terminating the for/range loop). If you've already ceased reading Messages, call
// Close; this will signal the PartitionConsumer's goroutines to begin shutting down (just like AsyncClose), but will
// also drain the Messages channel, harvest all errors & return them once cleanup has completed.
type PartitionConsumer interface ***REMOVED***

	// AsyncClose initiates a shutdown of the PartitionConsumer. This method will return immediately, after which you
	// should continue to service the 'Messages' and 'Errors' channels until they are empty. It is required to call this
	// function, or Close before a consumer object passes out of scope, as it will otherwise leak memory. You must call
	// this before calling Close on the underlying client.
	AsyncClose()

	// Close stops the PartitionConsumer from fetching messages. It will initiate a shutdown just like AsyncClose, drain
	// the Messages channel, harvest any errors & return them to the caller. Note that if you are continuing to service
	// the Messages channel when this function is called, you will be competing with Close for messages; consider
	// calling AsyncClose, instead. It is required to call this function (or AsyncClose) before a consumer object passes
	// out of scope, as it will otherwise leak memory. You must call this before calling Close on the underlying client.
	Close() error

	// Messages returns the read channel for the messages that are returned by
	// the broker.
	Messages() <-chan *ConsumerMessage

	// Errors returns a read channel of errors that occurred during consuming, if
	// enabled. By default, errors are logged and not returned over this channel.
	// If you want to implement any custom error handling, set your config's
	// Consumer.Return.Errors setting to true, and read from this channel.
	Errors() <-chan *ConsumerError

	// HighWaterMarkOffset returns the high water mark offset of the partition,
	// i.e. the offset that will be used for the next message that will be produced.
	// You can use this to determine how far behind the processing is.
	HighWaterMarkOffset() int64
***REMOVED***

type partitionConsumer struct ***REMOVED***
	highWaterMarkOffset int64 // must be at the top of the struct because https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	consumer            *consumer
	conf                *Config
	topic               string
	partition           int32

	broker   *brokerConsumer
	messages chan *ConsumerMessage
	errors   chan *ConsumerError
	feeder   chan *FetchResponse

	trigger, dying chan none
	responseResult error

	fetchSize int32
	offset    int64
***REMOVED***

var errTimedOut = errors.New("timed out feeding messages to the user") // not user-facing

func (child *partitionConsumer) sendError(err error) ***REMOVED***
	cErr := &ConsumerError***REMOVED***
		Topic:     child.topic,
		Partition: child.partition,
		Err:       err,
	***REMOVED***

	if child.conf.Consumer.Return.Errors ***REMOVED***
		child.errors <- cErr
	***REMOVED*** else ***REMOVED***
		Logger.Println(cErr)
	***REMOVED***
***REMOVED***

func (child *partitionConsumer) dispatcher() ***REMOVED***
	for range child.trigger ***REMOVED***
		select ***REMOVED***
		case <-child.dying:
			close(child.trigger)
		case <-time.After(child.conf.Consumer.Retry.Backoff):
			if child.broker != nil ***REMOVED***
				child.consumer.unrefBrokerConsumer(child.broker)
				child.broker = nil
			***REMOVED***

			Logger.Printf("consumer/%s/%d finding new broker\n", child.topic, child.partition)
			if err := child.dispatch(); err != nil ***REMOVED***
				child.sendError(err)
				child.trigger <- none***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if child.broker != nil ***REMOVED***
		child.consumer.unrefBrokerConsumer(child.broker)
	***REMOVED***
	child.consumer.removeChild(child)
	close(child.feeder)
***REMOVED***

func (child *partitionConsumer) dispatch() error ***REMOVED***
	if err := child.consumer.client.RefreshMetadata(child.topic); err != nil ***REMOVED***
		return err
	***REMOVED***

	var leader *Broker
	var err error
	if leader, err = child.consumer.client.Leader(child.topic, child.partition); err != nil ***REMOVED***
		return err
	***REMOVED***

	child.broker = child.consumer.refBrokerConsumer(leader)

	child.broker.input <- child

	return nil
***REMOVED***

func (child *partitionConsumer) chooseStartingOffset(offset int64) error ***REMOVED***
	newestOffset, err := child.consumer.client.GetOffset(child.topic, child.partition, OffsetNewest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	oldestOffset, err := child.consumer.client.GetOffset(child.topic, child.partition, OffsetOldest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	switch ***REMOVED***
	case offset == OffsetNewest:
		child.offset = newestOffset
	case offset == OffsetOldest:
		child.offset = oldestOffset
	case offset >= oldestOffset && offset <= newestOffset:
		child.offset = offset
	default:
		return ErrOffsetOutOfRange
	***REMOVED***

	return nil
***REMOVED***

func (child *partitionConsumer) Messages() <-chan *ConsumerMessage ***REMOVED***
	return child.messages
***REMOVED***

func (child *partitionConsumer) Errors() <-chan *ConsumerError ***REMOVED***
	return child.errors
***REMOVED***

func (child *partitionConsumer) AsyncClose() ***REMOVED***
	// this triggers whatever broker owns this child to abandon it and close its trigger channel, which causes
	// the dispatcher to exit its loop, which removes it from the consumer then closes its 'messages' and
	// 'errors' channel (alternatively, if the child is already at the dispatcher for some reason, that will
	// also just close itself)
	close(child.dying)
***REMOVED***

func (child *partitionConsumer) Close() error ***REMOVED***
	child.AsyncClose()

	go withRecover(func() ***REMOVED***
		for range child.messages ***REMOVED***
			// drain
		***REMOVED***
	***REMOVED***)

	var errors ConsumerErrors
	for err := range child.errors ***REMOVED***
		errors = append(errors, err)
	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		return errors
	***REMOVED***
	return nil
***REMOVED***

func (child *partitionConsumer) HighWaterMarkOffset() int64 ***REMOVED***
	return atomic.LoadInt64(&child.highWaterMarkOffset)
***REMOVED***

func (child *partitionConsumer) responseFeeder() ***REMOVED***
	var msgs []*ConsumerMessage
	expiryTicker := time.NewTicker(child.conf.Consumer.MaxProcessingTime)
	firstAttempt := true

feederLoop:
	for response := range child.feeder ***REMOVED***
		msgs, child.responseResult = child.parseResponse(response)

		for i, msg := range msgs ***REMOVED***
		messageSelect:
			select ***REMOVED***
			case child.messages <- msg:
				firstAttempt = true
			case <-expiryTicker.C:
				if !firstAttempt ***REMOVED***
					child.responseResult = errTimedOut
					child.broker.acks.Done()
					for _, msg = range msgs[i:] ***REMOVED***
						child.messages <- msg
					***REMOVED***
					child.broker.input <- child
					expiryTicker.Stop()
					continue feederLoop
				***REMOVED*** else ***REMOVED***
					// current message has not been sent, return to select
					// statement
					firstAttempt = false
					goto messageSelect
				***REMOVED***
			***REMOVED***
		***REMOVED***

		child.broker.acks.Done()
	***REMOVED***

	expiryTicker.Stop()
	close(child.messages)
	close(child.errors)
***REMOVED***

func (child *partitionConsumer) parseMessages(msgSet *MessageSet) ([]*ConsumerMessage, error) ***REMOVED***
	var messages []*ConsumerMessage
	var incomplete bool
	prelude := true

	for _, msgBlock := range msgSet.Messages ***REMOVED***
		for _, msg := range msgBlock.Messages() ***REMOVED***
			offset := msg.Offset
			if msg.Msg.Version >= 1 ***REMOVED***
				baseOffset := msgBlock.Offset - msgBlock.Messages()[len(msgBlock.Messages())-1].Offset
				offset += baseOffset
			***REMOVED***
			if prelude && offset < child.offset ***REMOVED***
				continue
			***REMOVED***
			prelude = false

			if offset >= child.offset ***REMOVED***
				messages = append(messages, &ConsumerMessage***REMOVED***
					Topic:          child.topic,
					Partition:      child.partition,
					Key:            msg.Msg.Key,
					Value:          msg.Msg.Value,
					Offset:         offset,
					Timestamp:      msg.Msg.Timestamp,
					BlockTimestamp: msgBlock.Msg.Timestamp,
				***REMOVED***)
				child.offset = offset + 1
			***REMOVED*** else ***REMOVED***
				incomplete = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if incomplete || len(messages) == 0 ***REMOVED***
		return nil, ErrIncompleteResponse
	***REMOVED***
	return messages, nil
***REMOVED***

func (child *partitionConsumer) parseRecords(batch *RecordBatch) ([]*ConsumerMessage, error) ***REMOVED***
	var messages []*ConsumerMessage
	var incomplete bool
	prelude := true
	originalOffset := child.offset

	for _, rec := range batch.Records ***REMOVED***
		offset := batch.FirstOffset + rec.OffsetDelta
		if prelude && offset < child.offset ***REMOVED***
			continue
		***REMOVED***
		prelude = false

		if offset >= child.offset ***REMOVED***
			messages = append(messages, &ConsumerMessage***REMOVED***
				Topic:     child.topic,
				Partition: child.partition,
				Key:       rec.Key,
				Value:     rec.Value,
				Offset:    offset,
				Timestamp: batch.FirstTimestamp.Add(rec.TimestampDelta),
				Headers:   rec.Headers,
			***REMOVED***)
			child.offset = offset + 1
		***REMOVED*** else ***REMOVED***
			incomplete = true
		***REMOVED***
	***REMOVED***

	if incomplete ***REMOVED***
		return nil, ErrIncompleteResponse
	***REMOVED***

	child.offset = batch.FirstOffset + int64(batch.LastOffsetDelta) + 1
	if child.offset <= originalOffset ***REMOVED***
		return nil, ErrConsumerOffsetNotAdvanced
	***REMOVED***

	return messages, nil
***REMOVED***

func (child *partitionConsumer) parseResponse(response *FetchResponse) ([]*ConsumerMessage, error) ***REMOVED***
	block := response.GetBlock(child.topic, child.partition)
	if block == nil ***REMOVED***
		return nil, ErrIncompleteResponse
	***REMOVED***

	if block.Err != ErrNoError ***REMOVED***
		return nil, block.Err
	***REMOVED***

	nRecs, err := block.numRecords()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if nRecs == 0 ***REMOVED***
		partialTrailingMessage, err := block.isPartial()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// We got no messages. If we got a trailing one then we need to ask for more data.
		// Otherwise we just poll again and wait for one to be produced...
		if partialTrailingMessage ***REMOVED***
			if child.conf.Consumer.Fetch.Max > 0 && child.fetchSize == child.conf.Consumer.Fetch.Max ***REMOVED***
				// we can't ask for more data, we've hit the configured limit
				child.sendError(ErrMessageTooLarge)
				child.offset++ // skip this one so we can keep processing future messages
			***REMOVED*** else ***REMOVED***
				child.fetchSize *= 2
				if child.conf.Consumer.Fetch.Max > 0 && child.fetchSize > child.conf.Consumer.Fetch.Max ***REMOVED***
					child.fetchSize = child.conf.Consumer.Fetch.Max
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return nil, nil
	***REMOVED***

	// we got messages, reset our fetch size in case it was increased for a previous request
	child.fetchSize = child.conf.Consumer.Fetch.Default
	atomic.StoreInt64(&child.highWaterMarkOffset, block.HighWaterMarkOffset)

	messages := []*ConsumerMessage***REMOVED******REMOVED***
	for _, records := range block.RecordsSet ***REMOVED***
		if control, err := records.isControl(); err != nil || control ***REMOVED***
			continue
		***REMOVED***

		switch records.recordsType ***REMOVED***
		case legacyRecords:
			messageSetMessages, err := child.parseMessages(records.msgSet)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			messages = append(messages, messageSetMessages...)
		case defaultRecords:
			recordBatchMessages, err := child.parseRecords(records.recordBatch)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			messages = append(messages, recordBatchMessages...)
		default:
			return nil, fmt.Errorf("unknown records type: %v", records.recordsType)
		***REMOVED***
	***REMOVED***

	return messages, nil
***REMOVED***

// brokerConsumer

type brokerConsumer struct ***REMOVED***
	consumer         *consumer
	broker           *Broker
	input            chan *partitionConsumer
	newSubscriptions chan []*partitionConsumer
	wait             chan none
	subscriptions    map[*partitionConsumer]none
	acks             sync.WaitGroup
	refs             int
***REMOVED***

func (c *consumer) newBrokerConsumer(broker *Broker) *brokerConsumer ***REMOVED***
	bc := &brokerConsumer***REMOVED***
		consumer:         c,
		broker:           broker,
		input:            make(chan *partitionConsumer),
		newSubscriptions: make(chan []*partitionConsumer),
		wait:             make(chan none),
		subscriptions:    make(map[*partitionConsumer]none),
		refs:             0,
	***REMOVED***

	go withRecover(bc.subscriptionManager)
	go withRecover(bc.subscriptionConsumer)

	return bc
***REMOVED***

func (bc *brokerConsumer) subscriptionManager() ***REMOVED***
	var buffer []*partitionConsumer

	// The subscriptionManager constantly accepts new subscriptions on `input` (even when the main subscriptionConsumer
	// goroutine is in the middle of a network request) and batches it up. The main worker goroutine picks
	// up a batch of new subscriptions between every network request by reading from `newSubscriptions`, so we give
	// it nil if no new subscriptions are available. We also write to `wait` only when new subscriptions is available,
	// so the main goroutine can block waiting for work if it has none.
	for ***REMOVED***
		if len(buffer) > 0 ***REMOVED***
			select ***REMOVED***
			case event, ok := <-bc.input:
				if !ok ***REMOVED***
					goto done
				***REMOVED***
				buffer = append(buffer, event)
			case bc.newSubscriptions <- buffer:
				buffer = nil
			case bc.wait <- none***REMOVED******REMOVED***:
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			select ***REMOVED***
			case event, ok := <-bc.input:
				if !ok ***REMOVED***
					goto done
				***REMOVED***
				buffer = append(buffer, event)
			case bc.newSubscriptions <- nil:
			***REMOVED***
		***REMOVED***
	***REMOVED***

done:
	close(bc.wait)
	if len(buffer) > 0 ***REMOVED***
		bc.newSubscriptions <- buffer
	***REMOVED***
	close(bc.newSubscriptions)
***REMOVED***

func (bc *brokerConsumer) subscriptionConsumer() ***REMOVED***
	<-bc.wait // wait for our first piece of work

	// the subscriptionConsumer ensures we will get nil right away if no new subscriptions is available
	for newSubscriptions := range bc.newSubscriptions ***REMOVED***
		bc.updateSubscriptions(newSubscriptions)

		if len(bc.subscriptions) == 0 ***REMOVED***
			// We're about to be shut down or we're about to receive more subscriptions.
			// Either way, the signal just hasn't propagated to our goroutine yet.
			<-bc.wait
			continue
		***REMOVED***

		response, err := bc.fetchNewMessages()

		if err != nil ***REMOVED***
			Logger.Printf("consumer/broker/%d disconnecting due to error processing FetchRequest: %s\n", bc.broker.ID(), err)
			bc.abort(err)
			return
		***REMOVED***

		bc.acks.Add(len(bc.subscriptions))
		for child := range bc.subscriptions ***REMOVED***
			child.feeder <- response
		***REMOVED***
		bc.acks.Wait()
		bc.handleResponses()
	***REMOVED***
***REMOVED***

func (bc *brokerConsumer) updateSubscriptions(newSubscriptions []*partitionConsumer) ***REMOVED***
	for _, child := range newSubscriptions ***REMOVED***
		bc.subscriptions[child] = none***REMOVED******REMOVED***
		Logger.Printf("consumer/broker/%d added subscription to %s/%d\n", bc.broker.ID(), child.topic, child.partition)
	***REMOVED***

	for child := range bc.subscriptions ***REMOVED***
		select ***REMOVED***
		case <-child.dying:
			Logger.Printf("consumer/broker/%d closed dead subscription to %s/%d\n", bc.broker.ID(), child.topic, child.partition)
			close(child.trigger)
			delete(bc.subscriptions, child)
		default:
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (bc *brokerConsumer) handleResponses() ***REMOVED***
	// handles the response codes left for us by our subscriptions, and abandons ones that have been closed
	for child := range bc.subscriptions ***REMOVED***
		result := child.responseResult
		child.responseResult = nil

		switch result ***REMOVED***
		case nil:
			break
		case errTimedOut:
			Logger.Printf("consumer/broker/%d abandoned subscription to %s/%d because consuming was taking too long\n",
				bc.broker.ID(), child.topic, child.partition)
			delete(bc.subscriptions, child)
		case ErrOffsetOutOfRange:
			// there's no point in retrying this it will just fail the same way again
			// shut it down and force the user to choose what to do
			child.sendError(result)
			Logger.Printf("consumer/%s/%d shutting down because %s\n", child.topic, child.partition, result)
			close(child.trigger)
			delete(bc.subscriptions, child)
		case ErrUnknownTopicOrPartition, ErrNotLeaderForPartition, ErrLeaderNotAvailable, ErrReplicaNotAvailable:
			// not an error, but does need redispatching
			Logger.Printf("consumer/broker/%d abandoned subscription to %s/%d because %s\n",
				bc.broker.ID(), child.topic, child.partition, result)
			child.trigger <- none***REMOVED******REMOVED***
			delete(bc.subscriptions, child)
		default:
			// dunno, tell the user and try redispatching
			child.sendError(result)
			Logger.Printf("consumer/broker/%d abandoned subscription to %s/%d because %s\n",
				bc.broker.ID(), child.topic, child.partition, result)
			child.trigger <- none***REMOVED******REMOVED***
			delete(bc.subscriptions, child)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (bc *brokerConsumer) abort(err error) ***REMOVED***
	bc.consumer.abandonBrokerConsumer(bc)
	_ = bc.broker.Close() // we don't care about the error this might return, we already have one

	for child := range bc.subscriptions ***REMOVED***
		child.sendError(err)
		child.trigger <- none***REMOVED******REMOVED***
	***REMOVED***

	for newSubscriptions := range bc.newSubscriptions ***REMOVED***
		if len(newSubscriptions) == 0 ***REMOVED***
			<-bc.wait
			continue
		***REMOVED***
		for _, child := range newSubscriptions ***REMOVED***
			child.sendError(err)
			child.trigger <- none***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (bc *brokerConsumer) fetchNewMessages() (*FetchResponse, error) ***REMOVED***
	request := &FetchRequest***REMOVED***
		MinBytes:    bc.consumer.conf.Consumer.Fetch.Min,
		MaxWaitTime: int32(bc.consumer.conf.Consumer.MaxWaitTime / time.Millisecond),
	***REMOVED***
	if bc.consumer.conf.Version.IsAtLeast(V0_10_0_0) ***REMOVED***
		request.Version = 2
	***REMOVED***
	if bc.consumer.conf.Version.IsAtLeast(V0_10_1_0) ***REMOVED***
		request.Version = 3
		request.MaxBytes = MaxResponseSize
	***REMOVED***
	if bc.consumer.conf.Version.IsAtLeast(V0_11_0_0) ***REMOVED***
		request.Version = 4
		request.Isolation = ReadUncommitted // We don't support yet transactions.
	***REMOVED***

	for child := range bc.subscriptions ***REMOVED***
		request.AddBlock(child.topic, child.partition, child.offset, child.fetchSize)
	***REMOVED***

	return bc.broker.Fetch(request)
***REMOVED***
