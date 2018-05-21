package sarama

import "encoding/binary"

// LengthField implements the PushEncoder and PushDecoder interfaces for calculating 4-byte lengths.
type lengthField struct ***REMOVED***
	startOffset int
***REMOVED***

func (l *lengthField) saveOffset(in int) ***REMOVED***
	l.startOffset = in
***REMOVED***

func (l *lengthField) reserveLength() int ***REMOVED***
	return 4
***REMOVED***

func (l *lengthField) run(curOffset int, buf []byte) error ***REMOVED***
	binary.BigEndian.PutUint32(buf[l.startOffset:], uint32(curOffset-l.startOffset-4))
	return nil
***REMOVED***

func (l *lengthField) check(curOffset int, buf []byte) error ***REMOVED***
	if uint32(curOffset-l.startOffset-4) != binary.BigEndian.Uint32(buf[l.startOffset:]) ***REMOVED***
		return PacketDecodingError***REMOVED***"length field invalid"***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

type varintLengthField struct ***REMOVED***
	startOffset int
	length      int64
***REMOVED***

func (l *varintLengthField) decode(pd packetDecoder) error ***REMOVED***
	var err error
	l.length, err = pd.getVarint()
	return err
***REMOVED***

func (l *varintLengthField) saveOffset(in int) ***REMOVED***
	l.startOffset = in
***REMOVED***

func (l *varintLengthField) adjustLength(currOffset int) int ***REMOVED***
	oldFieldSize := l.reserveLength()
	l.length = int64(currOffset - l.startOffset - oldFieldSize)

	return l.reserveLength() - oldFieldSize
***REMOVED***

func (l *varintLengthField) reserveLength() int ***REMOVED***
	var tmp [binary.MaxVarintLen64]byte
	return binary.PutVarint(tmp[:], l.length)
***REMOVED***

func (l *varintLengthField) run(curOffset int, buf []byte) error ***REMOVED***
	binary.PutVarint(buf[l.startOffset:], l.length)
	return nil
***REMOVED***

func (l *varintLengthField) check(curOffset int, buf []byte) error ***REMOVED***
	if int64(curOffset-l.startOffset-l.reserveLength()) != l.length ***REMOVED***
		return PacketDecodingError***REMOVED***"length field invalid"***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
