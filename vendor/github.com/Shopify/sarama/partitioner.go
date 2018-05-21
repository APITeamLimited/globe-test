package sarama

import (
	"hash"
	"hash/fnv"
	"math/rand"
	"time"
)

// Partitioner is anything that, given a Kafka message and a number of partitions indexed [0...numPartitions-1],
// decides to which partition to send the message. RandomPartitioner, RoundRobinPartitioner and HashPartitioner are provided
// as simple default implementations.
type Partitioner interface ***REMOVED***
	// Partition takes a message and partition count and chooses a partition
	Partition(message *ProducerMessage, numPartitions int32) (int32, error)

	// RequiresConsistency indicates to the user of the partitioner whether the
	// mapping of key->partition is consistent or not. Specifically, if a
	// partitioner requires consistency then it must be allowed to choose from all
	// partitions (even ones known to be unavailable), and its choice must be
	// respected by the caller. The obvious example is the HashPartitioner.
	RequiresConsistency() bool
***REMOVED***

// PartitionerConstructor is the type for a function capable of constructing new Partitioners.
type PartitionerConstructor func(topic string) Partitioner

type manualPartitioner struct***REMOVED******REMOVED***

// NewManualPartitioner returns a Partitioner which uses the partition manually set in the provided
// ProducerMessage's Partition field as the partition to produce to.
func NewManualPartitioner(topic string) Partitioner ***REMOVED***
	return new(manualPartitioner)
***REMOVED***

func (p *manualPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) ***REMOVED***
	return message.Partition, nil
***REMOVED***

func (p *manualPartitioner) RequiresConsistency() bool ***REMOVED***
	return true
***REMOVED***

type randomPartitioner struct ***REMOVED***
	generator *rand.Rand
***REMOVED***

// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
func NewRandomPartitioner(topic string) Partitioner ***REMOVED***
	p := new(randomPartitioner)
	p.generator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return p
***REMOVED***

func (p *randomPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) ***REMOVED***
	return int32(p.generator.Intn(int(numPartitions))), nil
***REMOVED***

func (p *randomPartitioner) RequiresConsistency() bool ***REMOVED***
	return false
***REMOVED***

type roundRobinPartitioner struct ***REMOVED***
	partition int32
***REMOVED***

// NewRoundRobinPartitioner returns a Partitioner which walks through the available partitions one at a time.
func NewRoundRobinPartitioner(topic string) Partitioner ***REMOVED***
	return &roundRobinPartitioner***REMOVED******REMOVED***
***REMOVED***

func (p *roundRobinPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) ***REMOVED***
	if p.partition >= numPartitions ***REMOVED***
		p.partition = 0
	***REMOVED***
	ret := p.partition
	p.partition++
	return ret, nil
***REMOVED***

func (p *roundRobinPartitioner) RequiresConsistency() bool ***REMOVED***
	return false
***REMOVED***

type hashPartitioner struct ***REMOVED***
	random Partitioner
	hasher hash.Hash32
***REMOVED***

// NewCustomHashPartitioner is a wrapper around NewHashPartitioner, allowing the use of custom hasher.
// The argument is a function providing the instance, implementing the hash.Hash32 interface. This is to ensure that
// each partition dispatcher gets its own hasher, to avoid concurrency issues by sharing an instance.
func NewCustomHashPartitioner(hasher func() hash.Hash32) PartitionerConstructor ***REMOVED***
	return func(topic string) Partitioner ***REMOVED***
		p := new(hashPartitioner)
		p.random = NewRandomPartitioner(topic)
		p.hasher = hasher()
		return p
	***REMOVED***
***REMOVED***

// NewHashPartitioner returns a Partitioner which behaves as follows. If the message's key is nil then a
// random partition is chosen. Otherwise the FNV-1a hash of the encoded bytes of the message key is used,
// modulus the number of partitions. This ensures that messages with the same key always end up on the
// same partition.
func NewHashPartitioner(topic string) Partitioner ***REMOVED***
	p := new(hashPartitioner)
	p.random = NewRandomPartitioner(topic)
	p.hasher = fnv.New32a()
	return p
***REMOVED***

func (p *hashPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) ***REMOVED***
	if message.Key == nil ***REMOVED***
		return p.random.Partition(message, numPartitions)
	***REMOVED***
	bytes, err := message.Key.Encode()
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	p.hasher.Reset()
	_, err = p.hasher.Write(bytes)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	partition := int32(p.hasher.Sum32()) % numPartitions
	if partition < 0 ***REMOVED***
		partition = -partition
	***REMOVED***
	return partition, nil
***REMOVED***

func (p *hashPartitioner) RequiresConsistency() bool ***REMOVED***
	return true
***REMOVED***
