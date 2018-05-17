package sarama

import (
	"time"
)

type AbortedTransaction struct ***REMOVED***
	ProducerID  int64
	FirstOffset int64
***REMOVED***

func (t *AbortedTransaction) decode(pd packetDecoder) (err error) ***REMOVED***
	if t.ProducerID, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if t.FirstOffset, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (t *AbortedTransaction) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.putInt64(t.ProducerID)
	pe.putInt64(t.FirstOffset)

	return nil
***REMOVED***

type FetchResponseBlock struct ***REMOVED***
	Err                 KError
	HighWaterMarkOffset int64
	LastStableOffset    int64
	AbortedTransactions []*AbortedTransaction
	Records             *Records // deprecated: use FetchResponseBlock.Records
	RecordsSet          []*Records
	Partial             bool
***REMOVED***

func (b *FetchResponseBlock) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	tmp, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.Err = KError(tmp)

	b.HighWaterMarkOffset, err = pd.getInt64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if version >= 4 ***REMOVED***
		b.LastStableOffset, err = pd.getInt64()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		numTransact, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if numTransact >= 0 ***REMOVED***
			b.AbortedTransactions = make([]*AbortedTransaction, numTransact)
		***REMOVED***

		for i := 0; i < numTransact; i++ ***REMOVED***
			transact := new(AbortedTransaction)
			if err = transact.decode(pd); err != nil ***REMOVED***
				return err
			***REMOVED***
			b.AbortedTransactions[i] = transact
		***REMOVED***
	***REMOVED***

	recordsSize, err := pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	recordsDecoder, err := pd.getSubset(int(recordsSize))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	b.RecordsSet = []*Records***REMOVED******REMOVED***

	for recordsDecoder.remaining() > 0 ***REMOVED***
		records := &Records***REMOVED******REMOVED***
		if err := records.decode(recordsDecoder); err != nil ***REMOVED***
			// If we have at least one decoded records, this is not an error
			if err == ErrInsufficientData ***REMOVED***
				if len(b.RecordsSet) == 0 ***REMOVED***
					b.Partial = true
				***REMOVED***
				break
			***REMOVED***
			return err
		***REMOVED***

		partial, err := records.isPartial()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// If we have at least one full records, we skip incomplete ones
		if partial && len(b.RecordsSet) > 0 ***REMOVED***
			break
		***REMOVED***

		b.RecordsSet = append(b.RecordsSet, records)

		if b.Records == nil ***REMOVED***
			b.Records = records
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (b *FetchResponseBlock) numRecords() (int, error) ***REMOVED***
	sum := 0

	for _, records := range b.RecordsSet ***REMOVED***
		count, err := records.numRecords()
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***

		sum += count
	***REMOVED***

	return sum, nil
***REMOVED***

func (b *FetchResponseBlock) isPartial() (bool, error) ***REMOVED***
	if b.Partial ***REMOVED***
		return true, nil
	***REMOVED***

	if len(b.RecordsSet) == 1 ***REMOVED***
		return b.RecordsSet[0].isPartial()
	***REMOVED***

	return false, nil
***REMOVED***

func (b *FetchResponseBlock) encode(pe packetEncoder, version int16) (err error) ***REMOVED***
	pe.putInt16(int16(b.Err))

	pe.putInt64(b.HighWaterMarkOffset)

	if version >= 4 ***REMOVED***
		pe.putInt64(b.LastStableOffset)

		if err = pe.putArrayLength(len(b.AbortedTransactions)); err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, transact := range b.AbortedTransactions ***REMOVED***
			if err = transact.encode(pe); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	pe.push(&lengthField***REMOVED******REMOVED***)
	for _, records := range b.RecordsSet ***REMOVED***
		err = records.encode(pe)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return pe.pop()
***REMOVED***

type FetchResponse struct ***REMOVED***
	Blocks       map[string]map[int32]*FetchResponseBlock
	ThrottleTime time.Duration
	Version      int16 // v1 requires 0.9+, v2 requires 0.10+
***REMOVED***

func (r *FetchResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	r.Version = version

	if r.Version >= 1 ***REMOVED***
		throttle, err := pd.getInt32()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r.ThrottleTime = time.Duration(throttle) * time.Millisecond
	***REMOVED***

	numTopics, err := pd.getArrayLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.Blocks = make(map[string]map[int32]*FetchResponseBlock, numTopics)
	for i := 0; i < numTopics; i++ ***REMOVED***
		name, err := pd.getString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		numBlocks, err := pd.getArrayLength()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.Blocks[name] = make(map[int32]*FetchResponseBlock, numBlocks)

		for j := 0; j < numBlocks; j++ ***REMOVED***
			id, err := pd.getInt32()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			block := new(FetchResponseBlock)
			err = block.decode(pd, version)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			r.Blocks[name][id] = block
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (r *FetchResponse) encode(pe packetEncoder) (err error) ***REMOVED***
	if r.Version >= 1 ***REMOVED***
		pe.putInt32(int32(r.ThrottleTime / time.Millisecond))
	***REMOVED***

	err = pe.putArrayLength(len(r.Blocks))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for topic, partitions := range r.Blocks ***REMOVED***
		err = pe.putString(topic)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = pe.putArrayLength(len(partitions))
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		for id, block := range partitions ***REMOVED***
			pe.putInt32(id)
			err = block.encode(pe, r.Version)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

	***REMOVED***
	return nil
***REMOVED***

func (r *FetchResponse) key() int16 ***REMOVED***
	return 1
***REMOVED***

func (r *FetchResponse) version() int16 ***REMOVED***
	return r.Version
***REMOVED***

func (r *FetchResponse) requiredVersion() KafkaVersion ***REMOVED***
	switch r.Version ***REMOVED***
	case 1:
		return V0_9_0_0
	case 2:
		return V0_10_0_0
	case 3:
		return V0_10_1_0
	case 4:
		return V0_11_0_0
	default:
		return minVersion
	***REMOVED***
***REMOVED***

func (r *FetchResponse) GetBlock(topic string, partition int32) *FetchResponseBlock ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		return nil
	***REMOVED***

	if r.Blocks[topic] == nil ***REMOVED***
		return nil
	***REMOVED***

	return r.Blocks[topic][partition]
***REMOVED***

func (r *FetchResponse) AddError(topic string, partition int32, err KError) ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		r.Blocks = make(map[string]map[int32]*FetchResponseBlock)
	***REMOVED***
	partitions, ok := r.Blocks[topic]
	if !ok ***REMOVED***
		partitions = make(map[int32]*FetchResponseBlock)
		r.Blocks[topic] = partitions
	***REMOVED***
	frb, ok := partitions[partition]
	if !ok ***REMOVED***
		frb = new(FetchResponseBlock)
		partitions[partition] = frb
	***REMOVED***
	frb.Err = err
***REMOVED***

func (r *FetchResponse) getOrCreateBlock(topic string, partition int32) *FetchResponseBlock ***REMOVED***
	if r.Blocks == nil ***REMOVED***
		r.Blocks = make(map[string]map[int32]*FetchResponseBlock)
	***REMOVED***
	partitions, ok := r.Blocks[topic]
	if !ok ***REMOVED***
		partitions = make(map[int32]*FetchResponseBlock)
		r.Blocks[topic] = partitions
	***REMOVED***
	frb, ok := partitions[partition]
	if !ok ***REMOVED***
		frb = new(FetchResponseBlock)
		partitions[partition] = frb
	***REMOVED***

	return frb
***REMOVED***

func encodeKV(key, value Encoder) ([]byte, []byte) ***REMOVED***
	var kb []byte
	var vb []byte
	if key != nil ***REMOVED***
		kb, _ = key.Encode()
	***REMOVED***
	if value != nil ***REMOVED***
		vb, _ = value.Encode()
	***REMOVED***

	return kb, vb
***REMOVED***

func (r *FetchResponse) AddMessage(topic string, partition int32, key, value Encoder, offset int64) ***REMOVED***
	frb := r.getOrCreateBlock(topic, partition)
	kb, vb := encodeKV(key, value)
	msg := &Message***REMOVED***Key: kb, Value: vb***REMOVED***
	msgBlock := &MessageBlock***REMOVED***Msg: msg, Offset: offset***REMOVED***
	if len(frb.RecordsSet) == 0 ***REMOVED***
		records := newLegacyRecords(&MessageSet***REMOVED******REMOVED***)
		frb.RecordsSet = []*Records***REMOVED***&records***REMOVED***
	***REMOVED***
	set := frb.RecordsSet[0].msgSet
	set.Messages = append(set.Messages, msgBlock)
***REMOVED***

func (r *FetchResponse) AddRecord(topic string, partition int32, key, value Encoder, offset int64) ***REMOVED***
	frb := r.getOrCreateBlock(topic, partition)
	kb, vb := encodeKV(key, value)
	rec := &Record***REMOVED***Key: kb, Value: vb, OffsetDelta: offset***REMOVED***
	if len(frb.RecordsSet) == 0 ***REMOVED***
		records := newDefaultRecords(&RecordBatch***REMOVED***Version: 2***REMOVED***)
		frb.RecordsSet = []*Records***REMOVED***&records***REMOVED***
	***REMOVED***
	batch := frb.RecordsSet[0].recordBatch
	batch.addRecord(rec)
***REMOVED***

func (r *FetchResponse) SetLastOffsetDelta(topic string, partition int32, offset int32) ***REMOVED***
	frb := r.getOrCreateBlock(topic, partition)
	if len(frb.RecordsSet) == 0 ***REMOVED***
		records := newDefaultRecords(&RecordBatch***REMOVED***Version: 2***REMOVED***)
		frb.RecordsSet = []*Records***REMOVED***&records***REMOVED***
	***REMOVED***
	batch := frb.RecordsSet[0].recordBatch
	batch.LastOffsetDelta = offset
***REMOVED***

func (r *FetchResponse) SetLastStableOffset(topic string, partition int32, offset int64) ***REMOVED***
	frb := r.getOrCreateBlock(topic, partition)
	frb.LastStableOffset = offset
***REMOVED***
