package sarama

import (
	"fmt"
)

// TestReporter has methods matching go's testing.T to avoid importing
// `testing` in the main part of the library.
type TestReporter interface ***REMOVED***
	Error(...interface***REMOVED******REMOVED***)
	Errorf(string, ...interface***REMOVED******REMOVED***)
	Fatal(...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// MockResponse is a response builder interface it defines one method that
// allows generating a response based on a request body. MockResponses are used
// to program behavior of MockBroker in tests.
type MockResponse interface ***REMOVED***
	For(reqBody versionedDecoder) (res encoder)
***REMOVED***

// MockWrapper is a mock response builder that returns a particular concrete
// response regardless of the actual request passed to the `For` method.
type MockWrapper struct ***REMOVED***
	res encoder
***REMOVED***

func (mw *MockWrapper) For(reqBody versionedDecoder) (res encoder) ***REMOVED***
	return mw.res
***REMOVED***

func NewMockWrapper(res encoder) *MockWrapper ***REMOVED***
	return &MockWrapper***REMOVED***res: res***REMOVED***
***REMOVED***

// MockSequence is a mock response builder that is created from a sequence of
// concrete responses. Every time when a `MockBroker` calls its `For` method
// the next response from the sequence is returned. When the end of the
// sequence is reached the last element from the sequence is returned.
type MockSequence struct ***REMOVED***
	responses []MockResponse
***REMOVED***

func NewMockSequence(responses ...interface***REMOVED******REMOVED***) *MockSequence ***REMOVED***
	ms := &MockSequence***REMOVED******REMOVED***
	ms.responses = make([]MockResponse, len(responses))
	for i, res := range responses ***REMOVED***
		switch res := res.(type) ***REMOVED***
		case MockResponse:
			ms.responses[i] = res
		case encoder:
			ms.responses[i] = NewMockWrapper(res)
		default:
			panic(fmt.Sprintf("Unexpected response type: %T", res))
		***REMOVED***
	***REMOVED***
	return ms
***REMOVED***

func (mc *MockSequence) For(reqBody versionedDecoder) (res encoder) ***REMOVED***
	res = mc.responses[0].For(reqBody)
	if len(mc.responses) > 1 ***REMOVED***
		mc.responses = mc.responses[1:]
	***REMOVED***
	return res
***REMOVED***

// MockMetadataResponse is a `MetadataResponse` builder.
type MockMetadataResponse struct ***REMOVED***
	leaders map[string]map[int32]int32
	brokers map[string]int32
	t       TestReporter
***REMOVED***

func NewMockMetadataResponse(t TestReporter) *MockMetadataResponse ***REMOVED***
	return &MockMetadataResponse***REMOVED***
		leaders: make(map[string]map[int32]int32),
		brokers: make(map[string]int32),
		t:       t,
	***REMOVED***
***REMOVED***

func (mmr *MockMetadataResponse) SetLeader(topic string, partition, brokerID int32) *MockMetadataResponse ***REMOVED***
	partitions := mmr.leaders[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]int32)
		mmr.leaders[topic] = partitions
	***REMOVED***
	partitions[partition] = brokerID
	return mmr
***REMOVED***

func (mmr *MockMetadataResponse) SetBroker(addr string, brokerID int32) *MockMetadataResponse ***REMOVED***
	mmr.brokers[addr] = brokerID
	return mmr
***REMOVED***

func (mmr *MockMetadataResponse) For(reqBody versionedDecoder) encoder ***REMOVED***
	metadataRequest := reqBody.(*MetadataRequest)
	metadataResponse := &MetadataResponse***REMOVED******REMOVED***
	for addr, brokerID := range mmr.brokers ***REMOVED***
		metadataResponse.AddBroker(addr, brokerID)
	***REMOVED***
	if len(metadataRequest.Topics) == 0 ***REMOVED***
		for topic, partitions := range mmr.leaders ***REMOVED***
			for partition, brokerID := range partitions ***REMOVED***
				metadataResponse.AddTopicPartition(topic, partition, brokerID, nil, nil, ErrNoError)
			***REMOVED***
		***REMOVED***
		return metadataResponse
	***REMOVED***
	for _, topic := range metadataRequest.Topics ***REMOVED***
		for partition, brokerID := range mmr.leaders[topic] ***REMOVED***
			metadataResponse.AddTopicPartition(topic, partition, brokerID, nil, nil, ErrNoError)
		***REMOVED***
	***REMOVED***
	return metadataResponse
***REMOVED***

// MockOffsetResponse is an `OffsetResponse` builder.
type MockOffsetResponse struct ***REMOVED***
	offsets map[string]map[int32]map[int64]int64
	t       TestReporter
	version int16
***REMOVED***

func NewMockOffsetResponse(t TestReporter) *MockOffsetResponse ***REMOVED***
	return &MockOffsetResponse***REMOVED***
		offsets: make(map[string]map[int32]map[int64]int64),
		t:       t,
	***REMOVED***
***REMOVED***

func (mor *MockOffsetResponse) SetVersion(version int16) *MockOffsetResponse ***REMOVED***
	mor.version = version
	return mor
***REMOVED***

func (mor *MockOffsetResponse) SetOffset(topic string, partition int32, time, offset int64) *MockOffsetResponse ***REMOVED***
	partitions := mor.offsets[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]map[int64]int64)
		mor.offsets[topic] = partitions
	***REMOVED***
	times := partitions[partition]
	if times == nil ***REMOVED***
		times = make(map[int64]int64)
		partitions[partition] = times
	***REMOVED***
	times[time] = offset
	return mor
***REMOVED***

func (mor *MockOffsetResponse) For(reqBody versionedDecoder) encoder ***REMOVED***
	offsetRequest := reqBody.(*OffsetRequest)
	offsetResponse := &OffsetResponse***REMOVED***Version: mor.version***REMOVED***
	for topic, partitions := range offsetRequest.blocks ***REMOVED***
		for partition, block := range partitions ***REMOVED***
			offset := mor.getOffset(topic, partition, block.time)
			offsetResponse.AddTopicPartition(topic, partition, offset)
		***REMOVED***
	***REMOVED***
	return offsetResponse
***REMOVED***

func (mor *MockOffsetResponse) getOffset(topic string, partition int32, time int64) int64 ***REMOVED***
	partitions := mor.offsets[topic]
	if partitions == nil ***REMOVED***
		mor.t.Errorf("missing topic: %s", topic)
	***REMOVED***
	times := partitions[partition]
	if times == nil ***REMOVED***
		mor.t.Errorf("missing partition: %d", partition)
	***REMOVED***
	offset, ok := times[time]
	if !ok ***REMOVED***
		mor.t.Errorf("missing time: %d", time)
	***REMOVED***
	return offset
***REMOVED***

// MockFetchResponse is a `FetchResponse` builder.
type MockFetchResponse struct ***REMOVED***
	messages       map[string]map[int32]map[int64]Encoder
	highWaterMarks map[string]map[int32]int64
	t              TestReporter
	batchSize      int
	version        int16
***REMOVED***

func NewMockFetchResponse(t TestReporter, batchSize int) *MockFetchResponse ***REMOVED***
	return &MockFetchResponse***REMOVED***
		messages:       make(map[string]map[int32]map[int64]Encoder),
		highWaterMarks: make(map[string]map[int32]int64),
		t:              t,
		batchSize:      batchSize,
	***REMOVED***
***REMOVED***

func (mfr *MockFetchResponse) SetVersion(version int16) *MockFetchResponse ***REMOVED***
	mfr.version = version
	return mfr
***REMOVED***

func (mfr *MockFetchResponse) SetMessage(topic string, partition int32, offset int64, msg Encoder) *MockFetchResponse ***REMOVED***
	partitions := mfr.messages[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]map[int64]Encoder)
		mfr.messages[topic] = partitions
	***REMOVED***
	messages := partitions[partition]
	if messages == nil ***REMOVED***
		messages = make(map[int64]Encoder)
		partitions[partition] = messages
	***REMOVED***
	messages[offset] = msg
	return mfr
***REMOVED***

func (mfr *MockFetchResponse) SetHighWaterMark(topic string, partition int32, offset int64) *MockFetchResponse ***REMOVED***
	partitions := mfr.highWaterMarks[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]int64)
		mfr.highWaterMarks[topic] = partitions
	***REMOVED***
	partitions[partition] = offset
	return mfr
***REMOVED***

func (mfr *MockFetchResponse) For(reqBody versionedDecoder) encoder ***REMOVED***
	fetchRequest := reqBody.(*FetchRequest)
	res := &FetchResponse***REMOVED***
		Version: mfr.version,
	***REMOVED***
	for topic, partitions := range fetchRequest.blocks ***REMOVED***
		for partition, block := range partitions ***REMOVED***
			initialOffset := block.fetchOffset
			offset := initialOffset
			maxOffset := initialOffset + int64(mfr.getMessageCount(topic, partition))
			for i := 0; i < mfr.batchSize && offset < maxOffset; ***REMOVED***
				msg := mfr.getMessage(topic, partition, offset)
				if msg != nil ***REMOVED***
					res.AddMessage(topic, partition, nil, msg, offset)
					i++
				***REMOVED***
				offset++
			***REMOVED***
			fb := res.GetBlock(topic, partition)
			if fb == nil ***REMOVED***
				res.AddError(topic, partition, ErrNoError)
				fb = res.GetBlock(topic, partition)
			***REMOVED***
			fb.HighWaterMarkOffset = mfr.getHighWaterMark(topic, partition)
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

func (mfr *MockFetchResponse) getMessage(topic string, partition int32, offset int64) Encoder ***REMOVED***
	partitions := mfr.messages[topic]
	if partitions == nil ***REMOVED***
		return nil
	***REMOVED***
	messages := partitions[partition]
	if messages == nil ***REMOVED***
		return nil
	***REMOVED***
	return messages[offset]
***REMOVED***

func (mfr *MockFetchResponse) getMessageCount(topic string, partition int32) int ***REMOVED***
	partitions := mfr.messages[topic]
	if partitions == nil ***REMOVED***
		return 0
	***REMOVED***
	messages := partitions[partition]
	if messages == nil ***REMOVED***
		return 0
	***REMOVED***
	return len(messages)
***REMOVED***

func (mfr *MockFetchResponse) getHighWaterMark(topic string, partition int32) int64 ***REMOVED***
	partitions := mfr.highWaterMarks[topic]
	if partitions == nil ***REMOVED***
		return 0
	***REMOVED***
	return partitions[partition]
***REMOVED***

// MockConsumerMetadataResponse is a `ConsumerMetadataResponse` builder.
type MockConsumerMetadataResponse struct ***REMOVED***
	coordinators map[string]interface***REMOVED******REMOVED***
	t            TestReporter
***REMOVED***

func NewMockConsumerMetadataResponse(t TestReporter) *MockConsumerMetadataResponse ***REMOVED***
	return &MockConsumerMetadataResponse***REMOVED***
		coordinators: make(map[string]interface***REMOVED******REMOVED***),
		t:            t,
	***REMOVED***
***REMOVED***

func (mr *MockConsumerMetadataResponse) SetCoordinator(group string, broker *MockBroker) *MockConsumerMetadataResponse ***REMOVED***
	mr.coordinators[group] = broker
	return mr
***REMOVED***

func (mr *MockConsumerMetadataResponse) SetError(group string, kerror KError) *MockConsumerMetadataResponse ***REMOVED***
	mr.coordinators[group] = kerror
	return mr
***REMOVED***

func (mr *MockConsumerMetadataResponse) For(reqBody versionedDecoder) encoder ***REMOVED***
	req := reqBody.(*ConsumerMetadataRequest)
	group := req.ConsumerGroup
	res := &ConsumerMetadataResponse***REMOVED******REMOVED***
	v := mr.coordinators[group]
	switch v := v.(type) ***REMOVED***
	case *MockBroker:
		res.Coordinator = &Broker***REMOVED***id: v.BrokerID(), addr: v.Addr()***REMOVED***
	case KError:
		res.Err = v
	***REMOVED***
	return res
***REMOVED***

// MockOffsetCommitResponse is a `OffsetCommitResponse` builder.
type MockOffsetCommitResponse struct ***REMOVED***
	errors map[string]map[string]map[int32]KError
	t      TestReporter
***REMOVED***

func NewMockOffsetCommitResponse(t TestReporter) *MockOffsetCommitResponse ***REMOVED***
	return &MockOffsetCommitResponse***REMOVED***t: t***REMOVED***
***REMOVED***

func (mr *MockOffsetCommitResponse) SetError(group, topic string, partition int32, kerror KError) *MockOffsetCommitResponse ***REMOVED***
	if mr.errors == nil ***REMOVED***
		mr.errors = make(map[string]map[string]map[int32]KError)
	***REMOVED***
	topics := mr.errors[group]
	if topics == nil ***REMOVED***
		topics = make(map[string]map[int32]KError)
		mr.errors[group] = topics
	***REMOVED***
	partitions := topics[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]KError)
		topics[topic] = partitions
	***REMOVED***
	partitions[partition] = kerror
	return mr
***REMOVED***

func (mr *MockOffsetCommitResponse) For(reqBody versionedDecoder) encoder ***REMOVED***
	req := reqBody.(*OffsetCommitRequest)
	group := req.ConsumerGroup
	res := &OffsetCommitResponse***REMOVED******REMOVED***
	for topic, partitions := range req.blocks ***REMOVED***
		for partition := range partitions ***REMOVED***
			res.AddError(topic, partition, mr.getError(group, topic, partition))
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

func (mr *MockOffsetCommitResponse) getError(group, topic string, partition int32) KError ***REMOVED***
	topics := mr.errors[group]
	if topics == nil ***REMOVED***
		return ErrNoError
	***REMOVED***
	partitions := topics[topic]
	if partitions == nil ***REMOVED***
		return ErrNoError
	***REMOVED***
	kerror, ok := partitions[partition]
	if !ok ***REMOVED***
		return ErrNoError
	***REMOVED***
	return kerror
***REMOVED***

// MockProduceResponse is a `ProduceResponse` builder.
type MockProduceResponse struct ***REMOVED***
	version int16
	errors  map[string]map[int32]KError
	t       TestReporter
***REMOVED***

func NewMockProduceResponse(t TestReporter) *MockProduceResponse ***REMOVED***
	return &MockProduceResponse***REMOVED***t: t***REMOVED***
***REMOVED***

func (mr *MockProduceResponse) SetVersion(version int16) *MockProduceResponse ***REMOVED***
	mr.version = version
	return mr
***REMOVED***

func (mr *MockProduceResponse) SetError(topic string, partition int32, kerror KError) *MockProduceResponse ***REMOVED***
	if mr.errors == nil ***REMOVED***
		mr.errors = make(map[string]map[int32]KError)
	***REMOVED***
	partitions := mr.errors[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]KError)
		mr.errors[topic] = partitions
	***REMOVED***
	partitions[partition] = kerror
	return mr
***REMOVED***

func (mr *MockProduceResponse) For(reqBody versionedDecoder) encoder ***REMOVED***
	req := reqBody.(*ProduceRequest)
	res := &ProduceResponse***REMOVED***
		Version: mr.version,
	***REMOVED***
	for topic, partitions := range req.records ***REMOVED***
		for partition := range partitions ***REMOVED***
			res.AddTopicPartition(topic, partition, mr.getError(topic, partition))
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

func (mr *MockProduceResponse) getError(topic string, partition int32) KError ***REMOVED***
	partitions := mr.errors[topic]
	if partitions == nil ***REMOVED***
		return ErrNoError
	***REMOVED***
	kerror, ok := partitions[partition]
	if !ok ***REMOVED***
		return ErrNoError
	***REMOVED***
	return kerror
***REMOVED***

// MockOffsetFetchResponse is a `OffsetFetchResponse` builder.
type MockOffsetFetchResponse struct ***REMOVED***
	offsets map[string]map[string]map[int32]*OffsetFetchResponseBlock
	t       TestReporter
***REMOVED***

func NewMockOffsetFetchResponse(t TestReporter) *MockOffsetFetchResponse ***REMOVED***
	return &MockOffsetFetchResponse***REMOVED***t: t***REMOVED***
***REMOVED***

func (mr *MockOffsetFetchResponse) SetOffset(group, topic string, partition int32, offset int64, metadata string, kerror KError) *MockOffsetFetchResponse ***REMOVED***
	if mr.offsets == nil ***REMOVED***
		mr.offsets = make(map[string]map[string]map[int32]*OffsetFetchResponseBlock)
	***REMOVED***
	topics := mr.offsets[group]
	if topics == nil ***REMOVED***
		topics = make(map[string]map[int32]*OffsetFetchResponseBlock)
		mr.offsets[group] = topics
	***REMOVED***
	partitions := topics[topic]
	if partitions == nil ***REMOVED***
		partitions = make(map[int32]*OffsetFetchResponseBlock)
		topics[topic] = partitions
	***REMOVED***
	partitions[partition] = &OffsetFetchResponseBlock***REMOVED***offset, metadata, kerror***REMOVED***
	return mr
***REMOVED***

func (mr *MockOffsetFetchResponse) For(reqBody versionedDecoder) encoder ***REMOVED***
	req := reqBody.(*OffsetFetchRequest)
	group := req.ConsumerGroup
	res := &OffsetFetchResponse***REMOVED******REMOVED***
	for topic, partitions := range mr.offsets[group] ***REMOVED***
		for partition, block := range partitions ***REMOVED***
			res.AddBlock(topic, partition, block)
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***
