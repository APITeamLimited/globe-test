package sarama

import (
	"sync"
	"time"
)

// Offset Manager

// OffsetManager uses Kafka to store and fetch consumed partition offsets.
type OffsetManager interface ***REMOVED***
	// ManagePartition creates a PartitionOffsetManager on the given topic/partition.
	// It will return an error if this OffsetManager is already managing the given
	// topic/partition.
	ManagePartition(topic string, partition int32) (PartitionOffsetManager, error)

	// Close stops the OffsetManager from managing offsets. It is required to call
	// this function before an OffsetManager object passes out of scope, as it
	// will otherwise leak memory. You must call this after all the
	// PartitionOffsetManagers are closed.
	Close() error
***REMOVED***

type offsetManager struct ***REMOVED***
	client Client
	conf   *Config
	group  string

	lock sync.Mutex
	poms map[string]map[int32]*partitionOffsetManager
	boms map[*Broker]*brokerOffsetManager
***REMOVED***

// NewOffsetManagerFromClient creates a new OffsetManager from the given client.
// It is still necessary to call Close() on the underlying client when finished with the partition manager.
func NewOffsetManagerFromClient(group string, client Client) (OffsetManager, error) ***REMOVED***
	// Check that we are not dealing with a closed Client before processing any other arguments
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	om := &offsetManager***REMOVED***
		client: client,
		conf:   client.Config(),
		group:  group,
		poms:   make(map[string]map[int32]*partitionOffsetManager),
		boms:   make(map[*Broker]*brokerOffsetManager),
	***REMOVED***

	return om, nil
***REMOVED***

func (om *offsetManager) ManagePartition(topic string, partition int32) (PartitionOffsetManager, error) ***REMOVED***
	pom, err := om.newPartitionOffsetManager(topic, partition)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	om.lock.Lock()
	defer om.lock.Unlock()

	topicManagers := om.poms[topic]
	if topicManagers == nil ***REMOVED***
		topicManagers = make(map[int32]*partitionOffsetManager)
		om.poms[topic] = topicManagers
	***REMOVED***

	if topicManagers[partition] != nil ***REMOVED***
		return nil, ConfigurationError("That topic/partition is already being managed")
	***REMOVED***

	topicManagers[partition] = pom
	return pom, nil
***REMOVED***

func (om *offsetManager) Close() error ***REMOVED***
	return nil
***REMOVED***

func (om *offsetManager) refBrokerOffsetManager(broker *Broker) *brokerOffsetManager ***REMOVED***
	om.lock.Lock()
	defer om.lock.Unlock()

	bom := om.boms[broker]
	if bom == nil ***REMOVED***
		bom = om.newBrokerOffsetManager(broker)
		om.boms[broker] = bom
	***REMOVED***

	bom.refs++

	return bom
***REMOVED***

func (om *offsetManager) unrefBrokerOffsetManager(bom *brokerOffsetManager) ***REMOVED***
	om.lock.Lock()
	defer om.lock.Unlock()

	bom.refs--

	if bom.refs == 0 ***REMOVED***
		close(bom.updateSubscriptions)
		if om.boms[bom.broker] == bom ***REMOVED***
			delete(om.boms, bom.broker)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (om *offsetManager) abandonBroker(bom *brokerOffsetManager) ***REMOVED***
	om.lock.Lock()
	defer om.lock.Unlock()

	delete(om.boms, bom.broker)
***REMOVED***

func (om *offsetManager) abandonPartitionOffsetManager(pom *partitionOffsetManager) ***REMOVED***
	om.lock.Lock()
	defer om.lock.Unlock()

	delete(om.poms[pom.topic], pom.partition)
	if len(om.poms[pom.topic]) == 0 ***REMOVED***
		delete(om.poms, pom.topic)
	***REMOVED***
***REMOVED***

// Partition Offset Manager

// PartitionOffsetManager uses Kafka to store and fetch consumed partition offsets. You MUST call Close()
// on a partition offset manager to avoid leaks, it will not be garbage-collected automatically when it passes
// out of scope.
type PartitionOffsetManager interface ***REMOVED***
	// NextOffset returns the next offset that should be consumed for the managed
	// partition, accompanied by metadata which can be used to reconstruct the state
	// of the partition consumer when it resumes. NextOffset() will return
	// `config.Consumer.Offsets.Initial` and an empty metadata string if no offset
	// was committed for this partition yet.
	NextOffset() (int64, string)

	// MarkOffset marks the provided offset, alongside a metadata string
	// that represents the state of the partition consumer at that point in time. The
	// metadata string can be used by another consumer to restore that state, so it
	// can resume consumption.
	//
	// To follow upstream conventions, you are expected to mark the offset of the
	// next message to read, not the last message read. Thus, when calling `MarkOffset`
	// you should typically add one to the offset of the last consumed message.
	//
	// Note: calling MarkOffset does not necessarily commit the offset to the backend
	// store immediately for efficiency reasons, and it may never be committed if
	// your application crashes. This means that you may end up processing the same
	// message twice, and your processing should ideally be idempotent.
	MarkOffset(offset int64, metadata string)

	// ResetOffset resets to the provided offset, alongside a metadata string that
	// represents the state of the partition consumer at that point in time. Reset
	// acts as a counterpart to MarkOffset, the difference being that it allows to
	// reset an offset to an earlier or smaller value, where MarkOffset only
	// allows incrementing the offset. cf MarkOffset for more details.
	ResetOffset(offset int64, metadata string)

	// Errors returns a read channel of errors that occur during offset management, if
	// enabled. By default, errors are logged and not returned over this channel. If
	// you want to implement any custom error handling, set your config's
	// Consumer.Return.Errors setting to true, and read from this channel.
	Errors() <-chan *ConsumerError

	// AsyncClose initiates a shutdown of the PartitionOffsetManager. This method will
	// return immediately, after which you should wait until the 'errors' channel has
	// been drained and closed. It is required to call this function, or Close before
	// a consumer object passes out of scope, as it will otherwise leak memory. You
	// must call this before calling Close on the underlying client.
	AsyncClose()

	// Close stops the PartitionOffsetManager from managing offsets. It is required to
	// call this function (or AsyncClose) before a PartitionOffsetManager object
	// passes out of scope, as it will otherwise leak memory. You must call this
	// before calling Close on the underlying client.
	Close() error
***REMOVED***

type partitionOffsetManager struct ***REMOVED***
	parent    *offsetManager
	topic     string
	partition int32

	lock     sync.Mutex
	offset   int64
	metadata string
	dirty    bool
	clean    sync.Cond
	broker   *brokerOffsetManager

	errors    chan *ConsumerError
	rebalance chan none
	dying     chan none
***REMOVED***

func (om *offsetManager) newPartitionOffsetManager(topic string, partition int32) (*partitionOffsetManager, error) ***REMOVED***
	pom := &partitionOffsetManager***REMOVED***
		parent:    om,
		topic:     topic,
		partition: partition,
		errors:    make(chan *ConsumerError, om.conf.ChannelBufferSize),
		rebalance: make(chan none, 1),
		dying:     make(chan none),
	***REMOVED***
	pom.clean.L = &pom.lock

	if err := pom.selectBroker(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := pom.fetchInitialOffset(om.conf.Metadata.Retry.Max); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pom.broker.updateSubscriptions <- pom

	go withRecover(pom.mainLoop)

	return pom, nil
***REMOVED***

func (pom *partitionOffsetManager) mainLoop() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-pom.rebalance:
			if err := pom.selectBroker(); err != nil ***REMOVED***
				pom.handleError(err)
				pom.rebalance <- none***REMOVED******REMOVED***
			***REMOVED*** else ***REMOVED***
				pom.broker.updateSubscriptions <- pom
			***REMOVED***
		case <-pom.dying:
			if pom.broker != nil ***REMOVED***
				select ***REMOVED***
				case <-pom.rebalance:
				case pom.broker.updateSubscriptions <- pom:
				***REMOVED***
				pom.parent.unrefBrokerOffsetManager(pom.broker)
			***REMOVED***
			pom.parent.abandonPartitionOffsetManager(pom)
			close(pom.errors)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (pom *partitionOffsetManager) selectBroker() error ***REMOVED***
	if pom.broker != nil ***REMOVED***
		pom.parent.unrefBrokerOffsetManager(pom.broker)
		pom.broker = nil
	***REMOVED***

	var broker *Broker
	var err error

	if err = pom.parent.client.RefreshCoordinator(pom.parent.group); err != nil ***REMOVED***
		return err
	***REMOVED***

	if broker, err = pom.parent.client.Coordinator(pom.parent.group); err != nil ***REMOVED***
		return err
	***REMOVED***

	pom.broker = pom.parent.refBrokerOffsetManager(broker)
	return nil
***REMOVED***

func (pom *partitionOffsetManager) fetchInitialOffset(retries int) error ***REMOVED***
	request := new(OffsetFetchRequest)
	request.Version = 1
	request.ConsumerGroup = pom.parent.group
	request.AddPartition(pom.topic, pom.partition)

	response, err := pom.broker.broker.FetchOffset(request)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	block := response.GetBlock(pom.topic, pom.partition)
	if block == nil ***REMOVED***
		return ErrIncompleteResponse
	***REMOVED***

	switch block.Err ***REMOVED***
	case ErrNoError:
		pom.offset = block.Offset
		pom.metadata = block.Metadata
		return nil
	case ErrNotCoordinatorForConsumer:
		if retries <= 0 ***REMOVED***
			return block.Err
		***REMOVED***
		if err := pom.selectBroker(); err != nil ***REMOVED***
			return err
		***REMOVED***
		return pom.fetchInitialOffset(retries - 1)
	case ErrOffsetsLoadInProgress:
		if retries <= 0 ***REMOVED***
			return block.Err
		***REMOVED***
		time.Sleep(pom.parent.conf.Metadata.Retry.Backoff)
		return pom.fetchInitialOffset(retries - 1)
	default:
		return block.Err
	***REMOVED***
***REMOVED***

func (pom *partitionOffsetManager) handleError(err error) ***REMOVED***
	cErr := &ConsumerError***REMOVED***
		Topic:     pom.topic,
		Partition: pom.partition,
		Err:       err,
	***REMOVED***

	if pom.parent.conf.Consumer.Return.Errors ***REMOVED***
		pom.errors <- cErr
	***REMOVED*** else ***REMOVED***
		Logger.Println(cErr)
	***REMOVED***
***REMOVED***

func (pom *partitionOffsetManager) Errors() <-chan *ConsumerError ***REMOVED***
	return pom.errors
***REMOVED***

func (pom *partitionOffsetManager) MarkOffset(offset int64, metadata string) ***REMOVED***
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if offset > pom.offset ***REMOVED***
		pom.offset = offset
		pom.metadata = metadata
		pom.dirty = true
	***REMOVED***
***REMOVED***

func (pom *partitionOffsetManager) ResetOffset(offset int64, metadata string) ***REMOVED***
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if offset <= pom.offset ***REMOVED***
		pom.offset = offset
		pom.metadata = metadata
		pom.dirty = true
	***REMOVED***
***REMOVED***

func (pom *partitionOffsetManager) updateCommitted(offset int64, metadata string) ***REMOVED***
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if pom.offset == offset && pom.metadata == metadata ***REMOVED***
		pom.dirty = false
		pom.clean.Signal()
	***REMOVED***
***REMOVED***

func (pom *partitionOffsetManager) NextOffset() (int64, string) ***REMOVED***
	pom.lock.Lock()
	defer pom.lock.Unlock()

	if pom.offset >= 0 ***REMOVED***
		return pom.offset, pom.metadata
	***REMOVED***

	return pom.parent.conf.Consumer.Offsets.Initial, ""
***REMOVED***

func (pom *partitionOffsetManager) AsyncClose() ***REMOVED***
	go func() ***REMOVED***
		pom.lock.Lock()
		defer pom.lock.Unlock()

		for pom.dirty ***REMOVED***
			pom.clean.Wait()
		***REMOVED***

		close(pom.dying)
	***REMOVED***()
***REMOVED***

func (pom *partitionOffsetManager) Close() error ***REMOVED***
	pom.AsyncClose()

	var errors ConsumerErrors
	for err := range pom.errors ***REMOVED***
		errors = append(errors, err)
	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		return errors
	***REMOVED***
	return nil
***REMOVED***

// Broker Offset Manager

type brokerOffsetManager struct ***REMOVED***
	parent              *offsetManager
	broker              *Broker
	timer               *time.Ticker
	updateSubscriptions chan *partitionOffsetManager
	subscriptions       map[*partitionOffsetManager]none
	refs                int
***REMOVED***

func (om *offsetManager) newBrokerOffsetManager(broker *Broker) *brokerOffsetManager ***REMOVED***
	bom := &brokerOffsetManager***REMOVED***
		parent:              om,
		broker:              broker,
		timer:               time.NewTicker(om.conf.Consumer.Offsets.CommitInterval),
		updateSubscriptions: make(chan *partitionOffsetManager),
		subscriptions:       make(map[*partitionOffsetManager]none),
	***REMOVED***

	go withRecover(bom.mainLoop)

	return bom
***REMOVED***

func (bom *brokerOffsetManager) mainLoop() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-bom.timer.C:
			if len(bom.subscriptions) > 0 ***REMOVED***
				bom.flushToBroker()
			***REMOVED***
		case s, ok := <-bom.updateSubscriptions:
			if !ok ***REMOVED***
				bom.timer.Stop()
				return
			***REMOVED***
			if _, ok := bom.subscriptions[s]; ok ***REMOVED***
				delete(bom.subscriptions, s)
			***REMOVED*** else ***REMOVED***
				bom.subscriptions[s] = none***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (bom *brokerOffsetManager) flushToBroker() ***REMOVED***
	request := bom.constructRequest()
	if request == nil ***REMOVED***
		return
	***REMOVED***

	response, err := bom.broker.CommitOffset(request)

	if err != nil ***REMOVED***
		bom.abort(err)
		return
	***REMOVED***

	for s := range bom.subscriptions ***REMOVED***
		if request.blocks[s.topic] == nil || request.blocks[s.topic][s.partition] == nil ***REMOVED***
			continue
		***REMOVED***

		var err KError
		var ok bool

		if response.Errors[s.topic] == nil ***REMOVED***
			s.handleError(ErrIncompleteResponse)
			delete(bom.subscriptions, s)
			s.rebalance <- none***REMOVED******REMOVED***
			continue
		***REMOVED***
		if err, ok = response.Errors[s.topic][s.partition]; !ok ***REMOVED***
			s.handleError(ErrIncompleteResponse)
			delete(bom.subscriptions, s)
			s.rebalance <- none***REMOVED******REMOVED***
			continue
		***REMOVED***

		switch err ***REMOVED***
		case ErrNoError:
			block := request.blocks[s.topic][s.partition]
			s.updateCommitted(block.offset, block.metadata)
		case ErrNotLeaderForPartition, ErrLeaderNotAvailable,
			ErrConsumerCoordinatorNotAvailable, ErrNotCoordinatorForConsumer:
			// not a critical error, we just need to redispatch
			delete(bom.subscriptions, s)
			s.rebalance <- none***REMOVED******REMOVED***
		case ErrOffsetMetadataTooLarge, ErrInvalidCommitOffsetSize:
			// nothing we can do about this, just tell the user and carry on
			s.handleError(err)
		case ErrOffsetsLoadInProgress:
			// nothing wrong but we didn't commit, we'll get it next time round
			break
		case ErrUnknownTopicOrPartition:
			// let the user know *and* try redispatching - if topic-auto-create is
			// enabled, redispatching should trigger a metadata request and create the
			// topic; if not then re-dispatching won't help, but we've let the user
			// know and it shouldn't hurt either (see https://github.com/Shopify/sarama/issues/706)
			fallthrough
		default:
			// dunno, tell the user and try redispatching
			s.handleError(err)
			delete(bom.subscriptions, s)
			s.rebalance <- none***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (bom *brokerOffsetManager) constructRequest() *OffsetCommitRequest ***REMOVED***
	var r *OffsetCommitRequest
	var perPartitionTimestamp int64
	if bom.parent.conf.Consumer.Offsets.Retention == 0 ***REMOVED***
		perPartitionTimestamp = ReceiveTime
		r = &OffsetCommitRequest***REMOVED***
			Version:                 1,
			ConsumerGroup:           bom.parent.group,
			ConsumerGroupGeneration: GroupGenerationUndefined,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r = &OffsetCommitRequest***REMOVED***
			Version:                 2,
			RetentionTime:           int64(bom.parent.conf.Consumer.Offsets.Retention / time.Millisecond),
			ConsumerGroup:           bom.parent.group,
			ConsumerGroupGeneration: GroupGenerationUndefined,
		***REMOVED***

	***REMOVED***

	for s := range bom.subscriptions ***REMOVED***
		s.lock.Lock()
		if s.dirty ***REMOVED***
			r.AddBlock(s.topic, s.partition, s.offset, perPartitionTimestamp, s.metadata)
		***REMOVED***
		s.lock.Unlock()
	***REMOVED***

	if len(r.blocks) > 0 ***REMOVED***
		return r
	***REMOVED***

	return nil
***REMOVED***

func (bom *brokerOffsetManager) abort(err error) ***REMOVED***
	_ = bom.broker.Close() // we don't care about the error this might return, we already have one
	bom.parent.abandonBroker(bom)

	for pom := range bom.subscriptions ***REMOVED***
		pom.handleError(err)
		pom.rebalance <- none***REMOVED******REMOVED***
	***REMOVED***

	for s := range bom.updateSubscriptions ***REMOVED***
		if _, ok := bom.subscriptions[s]; !ok ***REMOVED***
			s.handleError(err)
			s.rebalance <- none***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	bom.subscriptions = make(map[*partitionOffsetManager]none)
***REMOVED***
