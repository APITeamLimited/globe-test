package sarama

import "sync"

// SyncProducer publishes Kafka messages, blocking until they have been acknowledged. It routes messages to the correct
// broker, refreshing metadata as appropriate, and parses responses for errors. You must call Close() on a producer
// to avoid leaks, it may not be garbage-collected automatically when it passes out of scope.
//
// The SyncProducer comes with two caveats: it will generally be less efficient than the AsyncProducer, and the actual
// durability guarantee provided when a message is acknowledged depend on the configured value of `Producer.RequiredAcks`.
// There are configurations where a message acknowledged by the SyncProducer can still sometimes be lost.
//
// For implementation reasons, the SyncProducer requires `Producer.Return.Errors` and `Producer.Return.Successes` to
// be set to true in its configuration.
type SyncProducer interface ***REMOVED***

	// SendMessage produces a given message, and returns only when it either has
	// succeeded or failed to produce. It will return the partition and the offset
	// of the produced message, or an error if the message failed to produce.
	SendMessage(msg *ProducerMessage) (partition int32, offset int64, err error)

	// SendMessages produces a given set of messages, and returns only when all
	// messages in the set have either succeeded or failed. Note that messages
	// can succeed and fail individually; if some succeed and some fail,
	// SendMessages will return an error.
	SendMessages(msgs []*ProducerMessage) error

	// Close shuts down the producer and waits for any buffered messages to be
	// flushed. You must call this function before a producer object passes out of
	// scope, as it may otherwise leak memory. You must call this before calling
	// Close on the underlying client.
	Close() error
***REMOVED***

type syncProducer struct ***REMOVED***
	producer *asyncProducer
	wg       sync.WaitGroup
***REMOVED***

// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
func NewSyncProducer(addrs []string, config *Config) (SyncProducer, error) ***REMOVED***
	if config == nil ***REMOVED***
		config = NewConfig()
		config.Producer.Return.Successes = true
	***REMOVED***

	if err := verifyProducerConfig(config); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p, err := NewAsyncProducer(addrs, config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newSyncProducerFromAsyncProducer(p.(*asyncProducer)), nil
***REMOVED***

// NewSyncProducerFromClient creates a new SyncProducer using the given client. It is still
// necessary to call Close() on the underlying client when shutting down this producer.
func NewSyncProducerFromClient(client Client) (SyncProducer, error) ***REMOVED***
	if err := verifyProducerConfig(client.Config()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p, err := NewAsyncProducerFromClient(client)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newSyncProducerFromAsyncProducer(p.(*asyncProducer)), nil
***REMOVED***

func newSyncProducerFromAsyncProducer(p *asyncProducer) *syncProducer ***REMOVED***
	sp := &syncProducer***REMOVED***producer: p***REMOVED***

	sp.wg.Add(2)
	go withRecover(sp.handleSuccesses)
	go withRecover(sp.handleErrors)

	return sp
***REMOVED***

func verifyProducerConfig(config *Config) error ***REMOVED***
	if !config.Producer.Return.Errors ***REMOVED***
		return ConfigurationError("Producer.Return.Errors must be true to be used in a SyncProducer")
	***REMOVED***
	if !config.Producer.Return.Successes ***REMOVED***
		return ConfigurationError("Producer.Return.Successes must be true to be used in a SyncProducer")
	***REMOVED***
	return nil
***REMOVED***

func (sp *syncProducer) SendMessage(msg *ProducerMessage) (partition int32, offset int64, err error) ***REMOVED***
	oldMetadata := msg.Metadata
	defer func() ***REMOVED***
		msg.Metadata = oldMetadata
	***REMOVED***()

	expectation := make(chan *ProducerError, 1)
	msg.Metadata = expectation
	sp.producer.Input() <- msg

	if err := <-expectation; err != nil ***REMOVED***
		return -1, -1, err.Err
	***REMOVED***

	return msg.Partition, msg.Offset, nil
***REMOVED***

func (sp *syncProducer) SendMessages(msgs []*ProducerMessage) error ***REMOVED***
	savedMetadata := make([]interface***REMOVED******REMOVED***, len(msgs))
	for i := range msgs ***REMOVED***
		savedMetadata[i] = msgs[i].Metadata
	***REMOVED***
	defer func() ***REMOVED***
		for i := range msgs ***REMOVED***
			msgs[i].Metadata = savedMetadata[i]
		***REMOVED***
	***REMOVED***()

	expectations := make(chan chan *ProducerError, len(msgs))
	go func() ***REMOVED***
		for _, msg := range msgs ***REMOVED***
			expectation := make(chan *ProducerError, 1)
			msg.Metadata = expectation
			sp.producer.Input() <- msg
			expectations <- expectation
		***REMOVED***
		close(expectations)
	***REMOVED***()

	var errors ProducerErrors
	for expectation := range expectations ***REMOVED***
		if err := <-expectation; err != nil ***REMOVED***
			errors = append(errors, err)
		***REMOVED***
	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		return errors
	***REMOVED***
	return nil
***REMOVED***

func (sp *syncProducer) handleSuccesses() ***REMOVED***
	defer sp.wg.Done()
	for msg := range sp.producer.Successes() ***REMOVED***
		expectation := msg.Metadata.(chan *ProducerError)
		expectation <- nil
	***REMOVED***
***REMOVED***

func (sp *syncProducer) handleErrors() ***REMOVED***
	defer sp.wg.Done()
	for err := range sp.producer.Errors() ***REMOVED***
		expectation := err.Msg.Metadata.(chan *ProducerError)
		expectation <- err
	***REMOVED***
***REMOVED***

func (sp *syncProducer) Close() error ***REMOVED***
	sp.producer.AsyncClose()
	sp.wg.Wait()
	return nil
***REMOVED***
