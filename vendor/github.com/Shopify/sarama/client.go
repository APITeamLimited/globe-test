package sarama

import (
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Client is a generic Kafka client. It manages connections to one or more Kafka brokers.
// You MUST call Close() on a client to avoid leaks, it will not be garbage-collected
// automatically when it passes out of scope. It is safe to share a client amongst many
// users, however Kafka will process requests from a single client strictly in serial,
// so it is generally more efficient to use the default one client per producer/consumer.
type Client interface ***REMOVED***
	// Config returns the Config struct of the client. This struct should not be
	// altered after it has been created.
	Config() *Config

	// Brokers returns the current set of active brokers as retrieved from cluster metadata.
	Brokers() []*Broker

	// Topics returns the set of available topics as retrieved from cluster metadata.
	Topics() ([]string, error)

	// Partitions returns the sorted list of all partition IDs for the given topic.
	Partitions(topic string) ([]int32, error)

	// WritablePartitions returns the sorted list of all writable partition IDs for
	// the given topic, where "writable" means "having a valid leader accepting
	// writes".
	WritablePartitions(topic string) ([]int32, error)

	// Leader returns the broker object that is the leader of the current
	// topic/partition, as determined by querying the cluster metadata.
	Leader(topic string, partitionID int32) (*Broker, error)

	// Replicas returns the set of all replica IDs for the given partition.
	Replicas(topic string, partitionID int32) ([]int32, error)

	// InSyncReplicas returns the set of all in-sync replica IDs for the given
	// partition. In-sync replicas are replicas which are fully caught up with
	// the partition leader.
	InSyncReplicas(topic string, partitionID int32) ([]int32, error)

	// RefreshMetadata takes a list of topics and queries the cluster to refresh the
	// available metadata for those topics. If no topics are provided, it will refresh
	// metadata for all topics.
	RefreshMetadata(topics ...string) error

	// GetOffset queries the cluster to get the most recent available offset at the
	// given time (in milliseconds) on the topic/partition combination.
	// Time should be OffsetOldest for the earliest available offset,
	// OffsetNewest for the offset of the message that will be produced next, or a time.
	GetOffset(topic string, partitionID int32, time int64) (int64, error)

	// Coordinator returns the coordinating broker for a consumer group. It will
	// return a locally cached value if it's available. You can call
	// RefreshCoordinator to update the cached value. This function only works on
	// Kafka 0.8.2 and higher.
	Coordinator(consumerGroup string) (*Broker, error)

	// RefreshCoordinator retrieves the coordinator for a consumer group and stores it
	// in local cache. This function only works on Kafka 0.8.2 and higher.
	RefreshCoordinator(consumerGroup string) error

	// Close shuts down all broker connections managed by this client. It is required
	// to call this function before a client object passes out of scope, as it will
	// otherwise leak memory. You must close any Producers or Consumers using a client
	// before you close the client.
	Close() error

	// Closed returns true if the client has already had Close called on it
	Closed() bool
***REMOVED***

const (
	// OffsetNewest stands for the log head offset, i.e. the offset that will be
	// assigned to the next message that will be produced to the partition. You
	// can send this to a client's GetOffset method to get this offset, or when
	// calling ConsumePartition to start consuming new messages.
	OffsetNewest int64 = -1
	// OffsetOldest stands for the oldest offset available on the broker for a
	// partition. You can send this to a client's GetOffset method to get this
	// offset, or when calling ConsumePartition to start consuming from the
	// oldest offset that is still available on the broker.
	OffsetOldest int64 = -2
)

type client struct ***REMOVED***
	conf           *Config
	closer, closed chan none // for shutting down background metadata updater

	// the broker addresses given to us through the constructor are not guaranteed to be returned in
	// the cluster metadata (I *think* it only returns brokers who are currently leading partitions?)
	// so we store them separately
	seedBrokers []*Broker
	deadSeeds   []*Broker

	brokers      map[int32]*Broker                       // maps broker ids to brokers
	metadata     map[string]map[int32]*PartitionMetadata // maps topics to partition ids to metadata
	coordinators map[string]int32                        // Maps consumer group names to coordinating broker IDs

	// If the number of partitions is large, we can get some churn calling cachedPartitions,
	// so the result is cached.  It is important to update this value whenever metadata is changed
	cachedPartitionsResults map[string][maxPartitionIndex][]int32

	lock sync.RWMutex // protects access to the maps that hold cluster state.
***REMOVED***

// NewClient creates a new Client. It connects to one of the given broker addresses
// and uses that broker to automatically fetch metadata on the rest of the kafka cluster. If metadata cannot
// be retrieved from any of the given broker addresses, the client is not created.
func NewClient(addrs []string, conf *Config) (Client, error) ***REMOVED***
	Logger.Println("Initializing new client")

	if conf == nil ***REMOVED***
		conf = NewConfig()
	***REMOVED***

	if err := conf.Validate(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(addrs) < 1 ***REMOVED***
		return nil, ConfigurationError("You must provide at least one broker address")
	***REMOVED***

	client := &client***REMOVED***
		conf:                    conf,
		closer:                  make(chan none),
		closed:                  make(chan none),
		brokers:                 make(map[int32]*Broker),
		metadata:                make(map[string]map[int32]*PartitionMetadata),
		cachedPartitionsResults: make(map[string][maxPartitionIndex][]int32),
		coordinators:            make(map[string]int32),
	***REMOVED***

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, index := range random.Perm(len(addrs)) ***REMOVED***
		client.seedBrokers = append(client.seedBrokers, NewBroker(addrs[index]))
	***REMOVED***

	if conf.Metadata.Full ***REMOVED***
		// do an initial fetch of all cluster metadata by specifying an empty list of topics
		err := client.RefreshMetadata()
		switch err ***REMOVED***
		case nil:
			break
		case ErrLeaderNotAvailable, ErrReplicaNotAvailable, ErrTopicAuthorizationFailed, ErrClusterAuthorizationFailed:
			// indicates that maybe part of the cluster is down, but is not fatal to creating the client
			Logger.Println(err)
		default:
			close(client.closed) // we haven't started the background updater yet, so we have to do this manually
			_ = client.Close()
			return nil, err
		***REMOVED***
	***REMOVED***
	go withRecover(client.backgroundMetadataUpdater)

	Logger.Println("Successfully initialized new client")

	return client, nil
***REMOVED***

func (client *client) Config() *Config ***REMOVED***
	return client.conf
***REMOVED***

func (client *client) Brokers() []*Broker ***REMOVED***
	client.lock.RLock()
	defer client.lock.RUnlock()
	brokers := make([]*Broker, 0)
	for _, broker := range client.brokers ***REMOVED***
		brokers = append(brokers, broker)
	***REMOVED***
	return brokers
***REMOVED***

func (client *client) Close() error ***REMOVED***
	if client.Closed() ***REMOVED***
		// Chances are this is being called from a defer() and the error will go unobserved
		// so we go ahead and log the event in this case.
		Logger.Printf("Close() called on already closed client")
		return ErrClosedClient
	***REMOVED***

	// shutdown and wait for the background thread before we take the lock, to avoid races
	close(client.closer)
	<-client.closed

	client.lock.Lock()
	defer client.lock.Unlock()
	Logger.Println("Closing Client")

	for _, broker := range client.brokers ***REMOVED***
		safeAsyncClose(broker)
	***REMOVED***

	for _, broker := range client.seedBrokers ***REMOVED***
		safeAsyncClose(broker)
	***REMOVED***

	client.brokers = nil
	client.metadata = nil

	return nil
***REMOVED***

func (client *client) Closed() bool ***REMOVED***
	return client.brokers == nil
***REMOVED***

func (client *client) Topics() ([]string, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	client.lock.RLock()
	defer client.lock.RUnlock()

	ret := make([]string, 0, len(client.metadata))
	for topic := range client.metadata ***REMOVED***
		ret = append(ret, topic)
	***REMOVED***

	return ret, nil
***REMOVED***

func (client *client) Partitions(topic string) ([]int32, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	partitions := client.cachedPartitions(topic, allPartitions)

	if len(partitions) == 0 ***REMOVED***
		err := client.RefreshMetadata(topic)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		partitions = client.cachedPartitions(topic, allPartitions)
	***REMOVED***

	if partitions == nil ***REMOVED***
		return nil, ErrUnknownTopicOrPartition
	***REMOVED***

	return partitions, nil
***REMOVED***

func (client *client) WritablePartitions(topic string) ([]int32, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	partitions := client.cachedPartitions(topic, writablePartitions)

	// len==0 catches when it's nil (no such topic) and the odd case when every single
	// partition is undergoing leader election simultaneously. Callers have to be able to handle
	// this function returning an empty slice (which is a valid return value) but catching it
	// here the first time (note we *don't* catch it below where we return ErrUnknownTopicOrPartition) triggers
	// a metadata refresh as a nicety so callers can just try again and don't have to manually
	// trigger a refresh (otherwise they'd just keep getting a stale cached copy).
	if len(partitions) == 0 ***REMOVED***
		err := client.RefreshMetadata(topic)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		partitions = client.cachedPartitions(topic, writablePartitions)
	***REMOVED***

	if partitions == nil ***REMOVED***
		return nil, ErrUnknownTopicOrPartition
	***REMOVED***

	return partitions, nil
***REMOVED***

func (client *client) Replicas(topic string, partitionID int32) ([]int32, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	metadata := client.cachedMetadata(topic, partitionID)

	if metadata == nil ***REMOVED***
		err := client.RefreshMetadata(topic)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		metadata = client.cachedMetadata(topic, partitionID)
	***REMOVED***

	if metadata == nil ***REMOVED***
		return nil, ErrUnknownTopicOrPartition
	***REMOVED***

	if metadata.Err == ErrReplicaNotAvailable ***REMOVED***
		return dupInt32Slice(metadata.Replicas), metadata.Err
	***REMOVED***
	return dupInt32Slice(metadata.Replicas), nil
***REMOVED***

func (client *client) InSyncReplicas(topic string, partitionID int32) ([]int32, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	metadata := client.cachedMetadata(topic, partitionID)

	if metadata == nil ***REMOVED***
		err := client.RefreshMetadata(topic)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		metadata = client.cachedMetadata(topic, partitionID)
	***REMOVED***

	if metadata == nil ***REMOVED***
		return nil, ErrUnknownTopicOrPartition
	***REMOVED***

	if metadata.Err == ErrReplicaNotAvailable ***REMOVED***
		return dupInt32Slice(metadata.Isr), metadata.Err
	***REMOVED***
	return dupInt32Slice(metadata.Isr), nil
***REMOVED***

func (client *client) Leader(topic string, partitionID int32) (*Broker, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	leader, err := client.cachedLeader(topic, partitionID)

	if leader == nil ***REMOVED***
		err = client.RefreshMetadata(topic)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		leader, err = client.cachedLeader(topic, partitionID)
	***REMOVED***

	return leader, err
***REMOVED***

func (client *client) RefreshMetadata(topics ...string) error ***REMOVED***
	if client.Closed() ***REMOVED***
		return ErrClosedClient
	***REMOVED***

	// Prior to 0.8.2, Kafka will throw exceptions on an empty topic and not return a proper
	// error. This handles the case by returning an error instead of sending it
	// off to Kafka. See: https://github.com/Shopify/sarama/pull/38#issuecomment-26362310
	for _, topic := range topics ***REMOVED***
		if len(topic) == 0 ***REMOVED***
			return ErrInvalidTopic // this is the error that 0.8.2 and later correctly return
		***REMOVED***
	***REMOVED***

	return client.tryRefreshMetadata(topics, client.conf.Metadata.Retry.Max)
***REMOVED***

func (client *client) GetOffset(topic string, partitionID int32, time int64) (int64, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return -1, ErrClosedClient
	***REMOVED***

	offset, err := client.getOffset(topic, partitionID, time)

	if err != nil ***REMOVED***
		if err := client.RefreshMetadata(topic); err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		return client.getOffset(topic, partitionID, time)
	***REMOVED***

	return offset, err
***REMOVED***

func (client *client) Coordinator(consumerGroup string) (*Broker, error) ***REMOVED***
	if client.Closed() ***REMOVED***
		return nil, ErrClosedClient
	***REMOVED***

	coordinator := client.cachedCoordinator(consumerGroup)

	if coordinator == nil ***REMOVED***
		if err := client.RefreshCoordinator(consumerGroup); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		coordinator = client.cachedCoordinator(consumerGroup)
	***REMOVED***

	if coordinator == nil ***REMOVED***
		return nil, ErrConsumerCoordinatorNotAvailable
	***REMOVED***

	_ = coordinator.Open(client.conf)
	return coordinator, nil
***REMOVED***

func (client *client) RefreshCoordinator(consumerGroup string) error ***REMOVED***
	if client.Closed() ***REMOVED***
		return ErrClosedClient
	***REMOVED***

	response, err := client.getConsumerMetadata(consumerGroup, client.conf.Metadata.Retry.Max)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	client.lock.Lock()
	defer client.lock.Unlock()
	client.registerBroker(response.Coordinator)
	client.coordinators[consumerGroup] = response.Coordinator.ID()
	return nil
***REMOVED***

// private broker management helpers

// registerBroker makes sure a broker received by a Metadata or Coordinator request is registered
// in the brokers map. It returns the broker that is registered, which may be the provided broker,
// or a previously registered Broker instance. You must hold the write lock before calling this function.
func (client *client) registerBroker(broker *Broker) ***REMOVED***
	if client.brokers[broker.ID()] == nil ***REMOVED***
		client.brokers[broker.ID()] = broker
		Logger.Printf("client/brokers registered new broker #%d at %s", broker.ID(), broker.Addr())
	***REMOVED*** else if broker.Addr() != client.brokers[broker.ID()].Addr() ***REMOVED***
		safeAsyncClose(client.brokers[broker.ID()])
		client.brokers[broker.ID()] = broker
		Logger.Printf("client/brokers replaced registered broker #%d with %s", broker.ID(), broker.Addr())
	***REMOVED***
***REMOVED***

// deregisterBroker removes a broker from the seedsBroker list, and if it's
// not the seedbroker, removes it from brokers map completely.
func (client *client) deregisterBroker(broker *Broker) ***REMOVED***
	client.lock.Lock()
	defer client.lock.Unlock()

	if len(client.seedBrokers) > 0 && broker == client.seedBrokers[0] ***REMOVED***
		client.deadSeeds = append(client.deadSeeds, broker)
		client.seedBrokers = client.seedBrokers[1:]
	***REMOVED*** else ***REMOVED***
		// we do this so that our loop in `tryRefreshMetadata` doesn't go on forever,
		// but we really shouldn't have to; once that loop is made better this case can be
		// removed, and the function generally can be renamed from `deregisterBroker` to
		// `nextSeedBroker` or something
		Logger.Printf("client/brokers deregistered broker #%d at %s", broker.ID(), broker.Addr())
		delete(client.brokers, broker.ID())
	***REMOVED***
***REMOVED***

func (client *client) resurrectDeadBrokers() ***REMOVED***
	client.lock.Lock()
	defer client.lock.Unlock()

	Logger.Printf("client/brokers resurrecting %d dead seed brokers", len(client.deadSeeds))
	client.seedBrokers = append(client.seedBrokers, client.deadSeeds...)
	client.deadSeeds = nil
***REMOVED***

func (client *client) any() *Broker ***REMOVED***
	client.lock.RLock()
	defer client.lock.RUnlock()

	if len(client.seedBrokers) > 0 ***REMOVED***
		_ = client.seedBrokers[0].Open(client.conf)
		return client.seedBrokers[0]
	***REMOVED***

	// not guaranteed to be random *or* deterministic
	for _, broker := range client.brokers ***REMOVED***
		_ = broker.Open(client.conf)
		return broker
	***REMOVED***

	return nil
***REMOVED***

// private caching/lazy metadata helpers

type partitionType int

const (
	allPartitions partitionType = iota
	writablePartitions
	// If you add any more types, update the partition cache in update()

	// Ensure this is the last partition type value
	maxPartitionIndex
)

func (client *client) cachedMetadata(topic string, partitionID int32) *PartitionMetadata ***REMOVED***
	client.lock.RLock()
	defer client.lock.RUnlock()

	partitions := client.metadata[topic]
	if partitions != nil ***REMOVED***
		return partitions[partitionID]
	***REMOVED***

	return nil
***REMOVED***

func (client *client) cachedPartitions(topic string, partitionSet partitionType) []int32 ***REMOVED***
	client.lock.RLock()
	defer client.lock.RUnlock()

	partitions, exists := client.cachedPartitionsResults[topic]

	if !exists ***REMOVED***
		return nil
	***REMOVED***
	return partitions[partitionSet]
***REMOVED***

func (client *client) setPartitionCache(topic string, partitionSet partitionType) []int32 ***REMOVED***
	partitions := client.metadata[topic]

	if partitions == nil ***REMOVED***
		return nil
	***REMOVED***

	ret := make([]int32, 0, len(partitions))
	for _, partition := range partitions ***REMOVED***
		if partitionSet == writablePartitions && partition.Err == ErrLeaderNotAvailable ***REMOVED***
			continue
		***REMOVED***
		ret = append(ret, partition.ID)
	***REMOVED***

	sort.Sort(int32Slice(ret))
	return ret
***REMOVED***

func (client *client) cachedLeader(topic string, partitionID int32) (*Broker, error) ***REMOVED***
	client.lock.RLock()
	defer client.lock.RUnlock()

	partitions := client.metadata[topic]
	if partitions != nil ***REMOVED***
		metadata, ok := partitions[partitionID]
		if ok ***REMOVED***
			if metadata.Err == ErrLeaderNotAvailable ***REMOVED***
				return nil, ErrLeaderNotAvailable
			***REMOVED***
			b := client.brokers[metadata.Leader]
			if b == nil ***REMOVED***
				return nil, ErrLeaderNotAvailable
			***REMOVED***
			_ = b.Open(client.conf)
			return b, nil
		***REMOVED***
	***REMOVED***

	return nil, ErrUnknownTopicOrPartition
***REMOVED***

func (client *client) getOffset(topic string, partitionID int32, time int64) (int64, error) ***REMOVED***
	broker, err := client.Leader(topic, partitionID)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	request := &OffsetRequest***REMOVED******REMOVED***
	if client.conf.Version.IsAtLeast(V0_10_1_0) ***REMOVED***
		request.Version = 1
	***REMOVED***
	request.AddBlock(topic, partitionID, time, 1)

	response, err := broker.GetAvailableOffsets(request)
	if err != nil ***REMOVED***
		_ = broker.Close()
		return -1, err
	***REMOVED***

	block := response.GetBlock(topic, partitionID)
	if block == nil ***REMOVED***
		_ = broker.Close()
		return -1, ErrIncompleteResponse
	***REMOVED***
	if block.Err != ErrNoError ***REMOVED***
		return -1, block.Err
	***REMOVED***
	if len(block.Offsets) != 1 ***REMOVED***
		return -1, ErrOffsetOutOfRange
	***REMOVED***

	return block.Offsets[0], nil
***REMOVED***

// core metadata update logic

func (client *client) backgroundMetadataUpdater() ***REMOVED***
	defer close(client.closed)

	if client.conf.Metadata.RefreshFrequency == time.Duration(0) ***REMOVED***
		return
	***REMOVED***

	ticker := time.NewTicker(client.conf.Metadata.RefreshFrequency)
	defer ticker.Stop()

	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			topics := []string***REMOVED******REMOVED***
			if !client.conf.Metadata.Full ***REMOVED***
				if specificTopics, err := client.Topics(); err != nil ***REMOVED***
					Logger.Println("Client background metadata topic load:", err)
					break
				***REMOVED*** else if len(specificTopics) == 0 ***REMOVED***
					Logger.Println("Client background metadata update: no specific topics to update")
					break
				***REMOVED*** else ***REMOVED***
					topics = specificTopics
				***REMOVED***
			***REMOVED***

			if err := client.RefreshMetadata(topics...); err != nil ***REMOVED***
				Logger.Println("Client background metadata update:", err)
			***REMOVED***
		case <-client.closer:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (client *client) tryRefreshMetadata(topics []string, attemptsRemaining int) error ***REMOVED***
	retry := func(err error) error ***REMOVED***
		if attemptsRemaining > 0 ***REMOVED***
			Logger.Printf("client/metadata retrying after %dms... (%d attempts remaining)\n", client.conf.Metadata.Retry.Backoff/time.Millisecond, attemptsRemaining)
			time.Sleep(client.conf.Metadata.Retry.Backoff)
			return client.tryRefreshMetadata(topics, attemptsRemaining-1)
		***REMOVED***
		return err
	***REMOVED***

	for broker := client.any(); broker != nil; broker = client.any() ***REMOVED***
		if len(topics) > 0 ***REMOVED***
			Logger.Printf("client/metadata fetching metadata for %v from broker %s\n", topics, broker.addr)
		***REMOVED*** else ***REMOVED***
			Logger.Printf("client/metadata fetching metadata for all topics from broker %s\n", broker.addr)
		***REMOVED***
		response, err := broker.GetMetadata(&MetadataRequest***REMOVED***Topics: topics***REMOVED***)

		switch err.(type) ***REMOVED***
		case nil:
			// valid response, use it
			shouldRetry, err := client.updateMetadata(response)
			if shouldRetry ***REMOVED***
				Logger.Println("client/metadata found some partitions to be leaderless")
				return retry(err) // note: err can be nil
			***REMOVED***
			return err

		case PacketEncodingError:
			// didn't even send, return the error
			return err
		default:
			// some other error, remove that broker and try again
			Logger.Println("client/metadata got error from broker while fetching metadata:", err)
			_ = broker.Close()
			client.deregisterBroker(broker)
		***REMOVED***
	***REMOVED***

	Logger.Println("client/metadata no available broker to send metadata request to")
	client.resurrectDeadBrokers()
	return retry(ErrOutOfBrokers)
***REMOVED***

// if no fatal error, returns a list of topics that need retrying due to ErrLeaderNotAvailable
func (client *client) updateMetadata(data *MetadataResponse) (retry bool, err error) ***REMOVED***
	client.lock.Lock()
	defer client.lock.Unlock()

	// For all the brokers we received:
	// - if it is a new ID, save it
	// - if it is an existing ID, but the address we have is stale, discard the old one and save it
	// - otherwise ignore it, replacing our existing one would just bounce the connection
	for _, broker := range data.Brokers ***REMOVED***
		client.registerBroker(broker)
	***REMOVED***

	for _, topic := range data.Topics ***REMOVED***
		delete(client.metadata, topic.Name)
		delete(client.cachedPartitionsResults, topic.Name)

		switch topic.Err ***REMOVED***
		case ErrNoError:
			break
		case ErrInvalidTopic, ErrTopicAuthorizationFailed: // don't retry, don't store partial results
			err = topic.Err
			continue
		case ErrUnknownTopicOrPartition: // retry, do not store partial partition results
			err = topic.Err
			retry = true
			continue
		case ErrLeaderNotAvailable: // retry, but store partial partition results
			retry = true
			break
		default: // don't retry, don't store partial results
			Logger.Printf("Unexpected topic-level metadata error: %s", topic.Err)
			err = topic.Err
			continue
		***REMOVED***

		client.metadata[topic.Name] = make(map[int32]*PartitionMetadata, len(topic.Partitions))
		for _, partition := range topic.Partitions ***REMOVED***
			client.metadata[topic.Name][partition.ID] = partition
			if partition.Err == ErrLeaderNotAvailable ***REMOVED***
				retry = true
			***REMOVED***
		***REMOVED***

		var partitionCache [maxPartitionIndex][]int32
		partitionCache[allPartitions] = client.setPartitionCache(topic.Name, allPartitions)
		partitionCache[writablePartitions] = client.setPartitionCache(topic.Name, writablePartitions)
		client.cachedPartitionsResults[topic.Name] = partitionCache
	***REMOVED***

	return
***REMOVED***

func (client *client) cachedCoordinator(consumerGroup string) *Broker ***REMOVED***
	client.lock.RLock()
	defer client.lock.RUnlock()
	if coordinatorID, ok := client.coordinators[consumerGroup]; ok ***REMOVED***
		return client.brokers[coordinatorID]
	***REMOVED***
	return nil
***REMOVED***

func (client *client) getConsumerMetadata(consumerGroup string, attemptsRemaining int) (*ConsumerMetadataResponse, error) ***REMOVED***
	retry := func(err error) (*ConsumerMetadataResponse, error) ***REMOVED***
		if attemptsRemaining > 0 ***REMOVED***
			Logger.Printf("client/coordinator retrying after %dms... (%d attempts remaining)\n", client.conf.Metadata.Retry.Backoff/time.Millisecond, attemptsRemaining)
			time.Sleep(client.conf.Metadata.Retry.Backoff)
			return client.getConsumerMetadata(consumerGroup, attemptsRemaining-1)
		***REMOVED***
		return nil, err
	***REMOVED***

	for broker := client.any(); broker != nil; broker = client.any() ***REMOVED***
		Logger.Printf("client/coordinator requesting coordinator for consumergroup %s from %s\n", consumerGroup, broker.Addr())

		request := new(ConsumerMetadataRequest)
		request.ConsumerGroup = consumerGroup

		response, err := broker.GetConsumerMetadata(request)

		if err != nil ***REMOVED***
			Logger.Printf("client/coordinator request to broker %s failed: %s\n", broker.Addr(), err)

			switch err.(type) ***REMOVED***
			case PacketEncodingError:
				return nil, err
			default:
				_ = broker.Close()
				client.deregisterBroker(broker)
				continue
			***REMOVED***
		***REMOVED***

		switch response.Err ***REMOVED***
		case ErrNoError:
			Logger.Printf("client/coordinator coordinator for consumergroup %s is #%d (%s)\n", consumerGroup, response.Coordinator.ID(), response.Coordinator.Addr())
			return response, nil

		case ErrConsumerCoordinatorNotAvailable:
			Logger.Printf("client/coordinator coordinator for consumer group %s is not available\n", consumerGroup)

			// This is very ugly, but this scenario will only happen once per cluster.
			// The __consumer_offsets topic only has to be created one time.
			// The number of partitions not configurable, but partition 0 should always exist.
			if _, err := client.Leader("__consumer_offsets", 0); err != nil ***REMOVED***
				Logger.Printf("client/coordinator the __consumer_offsets topic is not initialized completely yet. Waiting 2 seconds...\n")
				time.Sleep(2 * time.Second)
			***REMOVED***

			return retry(ErrConsumerCoordinatorNotAvailable)
		default:
			return nil, response.Err
		***REMOVED***
	***REMOVED***

	Logger.Println("client/coordinator no available broker to send consumer metadata request to")
	client.resurrectDeadBrokers()
	return retry(ErrOutOfBrokers)
***REMOVED***
