package sarama

import "github.com/rcrowley/go-metrics"

// RequiredAcks is used in Produce Requests to tell the broker how many replica acknowledgements
// it must see before responding. Any of the constants defined here are valid. On broker versions
// prior to 0.8.2.0 any other positive int16 is also valid (the broker will wait for that many
// acknowledgements) but in 0.8.2.0 and later this will raise an exception (it has been replaced
// by setting the `min.isr` value in the brokers configuration).
type RequiredAcks int16

const (
	// NoResponse doesn't send any response, the TCP ACK is all you get.
	NoResponse RequiredAcks = 0
	// WaitForLocal waits for only the local commit to succeed before responding.
	WaitForLocal RequiredAcks = 1
	// WaitForAll waits for all in-sync replicas to commit before responding.
	// The minimum number of in-sync replicas is configured on the broker via
	// the `min.insync.replicas` configuration key.
	WaitForAll RequiredAcks = -1
)

type ProduceRequest struct ***REMOVED***
	TransactionalID *string
	RequiredAcks    RequiredAcks
	Timeout         int32
	Version         int16 // v1 requires Kafka 0.9, v2 requires Kafka 0.10, v3 requires Kafka 0.11
	records         map[string]map[int32]Records
***REMOVED***

func updateMsgSetMetrics(msgSet *MessageSet, compressionRatioMetric metrics.Histogram,
	topicCompressionRatioMetric metrics.Histogram) int64 ***REMOVED***
	var topicRecordCount int64
	for _, messageBlock := range msgSet.Messages ***REMOVED***
		// Is this a fake "message" wrapping real messages?
		if messageBlock.Msg.Set != nil ***REMOVED***
			topicRecordCount += int64(len(messageBlock.Msg.Set.Messages))
		***REMOVED*** else ***REMOVED***
			// A single uncompressed message
			topicRecordCount++
		***REMOVED***
		// Better be safe than sorry when computing the compression ratio
		if messageBlock.Msg.compressedSize != 0 ***REMOVED***
			compressionRatio := float64(len(messageBlock.Msg.Value)) /
				float64(messageBlock.Msg.compressedSize)
			// Histogram do not support decimal values, let's multiple it by 100 for better precision
			intCompressionRatio := int64(100 * compressionRatio)
			compressionRatioMetric.Update(intCompressionRatio)
			topicCompressionRatioMetric.Update(intCompressionRatio)
		***REMOVED***
	***REMOVED***
	return topicRecordCount
***REMOVED***

func updateBatchMetrics(recordBatch *RecordBatch, compressionRatioMetric metrics.Histogram,
	topicCompressionRatioMetric metrics.Histogram) int64 ***REMOVED***
	if recordBatch.compressedRecords != nil ***REMOVED***
		compressionRatio := int64(float64(recordBatch.recordsLen) / float64(len(recordBatch.compressedRecords)) * 100)
		compressionRatioMetric.Update(compressionRatio)
		topicCompressionRatioMetric.Update(compressionRatio)
	***REMOVED***

	return int64(len(recordBatch.Records))
***REMOVED***

func (r *ProduceRequest) encode(pe packetEncoder) error ***REMOVED***
	if r.Version >= 3 ***REMOVED***
		if err := pe.putNullableString(r.TransactionalID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	pe.putInt16(int16(r.RequiredAcks))
	pe.putInt32(r.Timeout)
	metricRegistry := pe.metricRegistry()
	var batchSizeMetric metrics.Histogram
	var compressionRatioMetric metrics.Histogram
	if metricRegistry != nil ***REMOVED***
		batchSizeMetric = getOrRegisterHistogram("batch-size", metricRegistry)
		compressionRatioMetric = getOrRegisterHistogram("compression-ratio", metricRegistry)
	***REMOVED***
	totalRecordCount := int64(0)

	err := pe.putArrayLength(len(r.records))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, partitions := range r.records ***REMOVED***
		err = pe.putString(topic)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = pe.putArrayLength(len(partitions))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		topicRecordCount := int64(0)
		var topicCompressionRatioMetric metrics.Histogram
		if metricRegistry != nil ***REMOVED***
			topicCompressionRatioMetric = getOrRegisterTopicHistogram("compression-ratio", topic, metricRegistry)
		***REMOVED***
		for id, records := range partitions ***REMOVED***
			startOffset := pe.offset()
			pe.putInt32(id)
			pe.push(&lengthField***REMOVED******REMOVED***)
			err = records.encode(pe)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			err = pe.pop()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if metricRegistry != nil ***REMOVED***
				if r.Version >= 3 ***REMOVED***
					topicRecordCount += updateBatchMetrics(records.recordBatch, compressionRatioMetric, topicCompressionRatioMetric)
				***REMOVED*** else ***REMOVED***
					topicRecordCount += updateMsgSetMetrics(records.msgSet, compressionRatioMetric, topicCompressionRatioMetric)
				***REMOVED***
				batchSize := int64(pe.offset() - startOffset)
				batchSizeMetric.Update(batchSize)
				getOrRegisterTopicHistogram("batch-size", topic, metricRegistry).Update(batchSize)
			***REMOVED***
		***REMOVED***
		if topicRecordCount > 0 ***REMOVED***
			getOrRegisterTopicMeter("record-send-rate", topic, metricRegistry).Mark(topicRecordCount)
			getOrRegisterTopicHistogram("records-per-request", topic, metricRegistry).Update(topicRecordCount)
			totalRecordCount += topicRecordCount
		***REMOVED***
	***REMOVED***
	if totalRecordCount > 0 ***REMOVED***
		metrics.GetOrRegisterMeter("record-send-rate", metricRegistry).Mark(totalRecordCount)
		getOrRegisterHistogram("records-per-request", metricRegistry).Update(totalRecordCount)
	***REMOVED***

	return nil
***REMOVED***

func (r *ProduceRequest) decode(pd packetDecoder, version int16) error ***REMOVED***
	r.Version = version

	if version >= 3 ***REMOVED***
		id, err := pd.getNullableString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.TransactionalID = id
	***REMOVED***
	requiredAcks, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.RequiredAcks = RequiredAcks(requiredAcks)
	if r.Timeout, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	topicCount, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if topicCount == 0 ***REMOVED***
		return nil
	***REMOVED***

	r.records = make(map[string]map[int32]Records)
	for i := 0; i < topicCount; i++ ***REMOVED***
		topic, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		partitionCount, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.records[topic] = make(map[int32]Records)

		for j := 0; j < partitionCount; j++ ***REMOVED***
			partition, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			size, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			recordsDecoder, err := pd.getSubset(int(size))
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			var records Records
			if err := records.decode(recordsDecoder); err != nil ***REMOVED***
				return err
			***REMOVED***
			r.records[topic][partition] = records
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *ProduceRequest) key() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ProduceRequest) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *ProduceRequest) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_9_0_0
	case 2:
		return V0_10_0_0
	case 3:
		return V0_11_0_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

func (r *ProduceRequest) ensureRecords(topic string, partition int32) ***REMOVED***
	if r.records == nil ***REMOVED***
		r.records = make(map[string]map[int32]Records)
	***REMOVED***

	if r.records[topic] == nil ***REMOVED***
		r.records[topic] = make(map[int32]Records)
	***REMOVED***
***REMOVED***

func (r *ProduceRequest) AddMessage(topic string, partition int32, msg *Message) ***REMOVED***
	r.ensureRecords(topic, partition)
	set := r.records[topic][partition].msgSet

	if set == nil ***REMOVED***
		set = new(MessageSet)
		r.records[topic][partition] = newLegacyRecords(set)
	***REMOVED***

	set.addMessage(msg)
***REMOVED***

func (r *ProduceRequest) AddSet(topic string, partition int32, set *MessageSet) ***REMOVED***
	r.ensureRecords(topic, partition)
	r.records[topic][partition] = newLegacyRecords(set)
***REMOVED***

func (r *ProduceRequest) AddBatch(topic string, partition int32, batch *RecordBatch) ***REMOVED***
	r.ensureRecords(topic, partition)
	r.records[topic][partition] = newDefaultRecords(batch)
***REMOVED***
