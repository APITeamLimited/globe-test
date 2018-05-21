package sarama

import "fmt"

const (
	unknownRecords = iota
	legacyRecords
	defaultRecords

	magicOffset = 16
	magicLength = 1
)

// Records implements a union type containing either a RecordBatch or a legacy MessageSet.
type Records struct ***REMOVED***
	recordsType int
	msgSet      *MessageSet
	recordBatch *RecordBatch
***REMOVED***

func newLegacyRecords(msgSet *MessageSet) Records ***REMOVED***
	return Records***REMOVED***recordsType: legacyRecords, msgSet: msgSet***REMOVED***
***REMOVED***

func newDefaultRecords(batch *RecordBatch) Records ***REMOVED***
	return Records***REMOVED***recordsType: defaultRecords, recordBatch: batch***REMOVED***
***REMOVED***

// setTypeFromFields sets type of Records depending on which of msgSet or recordBatch is not nil.
// The first return value indicates whether both fields are nil (and the type is not set).
// If both fields are not nil, it returns an error.
func (r *Records) setTypeFromFields() (bool, error) ***REMOVED***
	if r.msgSet == nil && r.recordBatch == nil ***REMOVED***
		return true, nil
	***REMOVED***
	if r.msgSet != nil && r.recordBatch != nil ***REMOVED***
		return false, fmt.Errorf("both msgSet and recordBatch are set, but record type is unknown")
	***REMOVED***
	r.recordsType = defaultRecords
	if r.msgSet != nil ***REMOVED***
		r.recordsType = legacyRecords
	***REMOVED***
	return false, nil
***REMOVED***

func (r *Records) encode(pe packetEncoder) error ***REMOVED***
	if r.recordsType == unknownRecords ***REMOVED***
		if empty, err := r.setTypeFromFields(); err != nil || empty ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	switch r.recordsType ***REMOVED***
	case legacyRecords:
		if r.msgSet == nil ***REMOVED***
			return nil
		***REMOVED***
		return r.msgSet.encode(pe)
	case defaultRecords:
		if r.recordBatch == nil ***REMOVED***
			return nil
		***REMOVED***
		return r.recordBatch.encode(pe)
	***REMOVED***

	return fmt.Errorf("unknown records type: %v", r.recordsType)
***REMOVED***

func (r *Records) setTypeFromMagic(pd packetDecoder) error ***REMOVED***
	magic, err := magicValue(pd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	r.recordsType = defaultRecords
	if magic < 2 ***REMOVED***
		r.recordsType = legacyRecords
	***REMOVED***

	return nil
***REMOVED***

func (r *Records) decode(pd packetDecoder) error ***REMOVED***
	if r.recordsType == unknownRecords ***REMOVED***
		if err := r.setTypeFromMagic(pd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	switch r.recordsType ***REMOVED***
	case legacyRecords:
		r.msgSet = &MessageSet***REMOVED******REMOVED***
		return r.msgSet.decode(pd)
	case defaultRecords:
		r.recordBatch = &RecordBatch***REMOVED******REMOVED***
		return r.recordBatch.decode(pd)
	***REMOVED***
	return fmt.Errorf("unknown records type: %v", r.recordsType)
***REMOVED***

func (r *Records) numRecords() (int, error) ***REMOVED***
	if r.recordsType == unknownRecords ***REMOVED***
		if empty, err := r.setTypeFromFields(); err != nil || empty ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***

	switch r.recordsType ***REMOVED***
	case legacyRecords:
		if r.msgSet == nil ***REMOVED***
			return 0, nil
		***REMOVED***
		return len(r.msgSet.Messages), nil
	case defaultRecords:
		if r.recordBatch == nil ***REMOVED***
			return 0, nil
		***REMOVED***
		return len(r.recordBatch.Records), nil
	***REMOVED***
	return 0, fmt.Errorf("unknown records type: %v", r.recordsType)
***REMOVED***

func (r *Records) isPartial() (bool, error) ***REMOVED***
	if r.recordsType == unknownRecords ***REMOVED***
		if empty, err := r.setTypeFromFields(); err != nil || empty ***REMOVED***
			return false, err
		***REMOVED***
	***REMOVED***

	switch r.recordsType ***REMOVED***
	case unknownRecords:
		return false, nil
	case legacyRecords:
		if r.msgSet == nil ***REMOVED***
			return false, nil
		***REMOVED***
		return r.msgSet.PartialTrailingMessage, nil
	case defaultRecords:
		if r.recordBatch == nil ***REMOVED***
			return false, nil
		***REMOVED***
		return r.recordBatch.PartialTrailingRecord, nil
	***REMOVED***
	return false, fmt.Errorf("unknown records type: %v", r.recordsType)
***REMOVED***

func (r *Records) isControl() (bool, error) ***REMOVED***
	if r.recordsType == unknownRecords ***REMOVED***
		if empty, err := r.setTypeFromFields(); err != nil || empty ***REMOVED***
			return false, err
		***REMOVED***
	***REMOVED***

	switch r.recordsType ***REMOVED***
	case legacyRecords:
		return false, nil
	case defaultRecords:
		if r.recordBatch == nil ***REMOVED***
			return false, nil
		***REMOVED***
		return r.recordBatch.Control, nil
	***REMOVED***
	return false, fmt.Errorf("unknown records type: %v", r.recordsType)
***REMOVED***

func magicValue(pd packetDecoder) (int8, error) ***REMOVED***
	dec, err := pd.peek(magicOffset, magicLength)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return dec.getInt8()
***REMOVED***
