package sarama

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

type crcPolynomial int8

const (
	crcIEEE crcPolynomial = iota
	crcCastagnoli
)

var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)

// crc32Field implements the pushEncoder and pushDecoder interfaces for calculating CRC32s.
type crc32Field struct ***REMOVED***
	startOffset int
	polynomial  crcPolynomial
***REMOVED***

func (c *crc32Field) saveOffset(in int) ***REMOVED***
	c.startOffset = in
***REMOVED***

func (c *crc32Field) reserveLength() int ***REMOVED***
	return 4
***REMOVED***

func newCRC32Field(polynomial crcPolynomial) *crc32Field ***REMOVED***
	return &crc32Field***REMOVED***polynomial: polynomial***REMOVED***
***REMOVED***

func (c *crc32Field) run(curOffset int, buf []byte) error ***REMOVED***
	crc, err := c.crc(curOffset, buf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	binary.BigEndian.PutUint32(buf[c.startOffset:], crc)
	return nil
***REMOVED***

func (c *crc32Field) check(curOffset int, buf []byte) error ***REMOVED***
	crc, err := c.crc(curOffset, buf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	expected := binary.BigEndian.Uint32(buf[c.startOffset:])
	if crc != expected ***REMOVED***
		return PacketDecodingError***REMOVED***fmt.Sprintf("CRC didn't match expected %#x got %#x", expected, crc)***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
func (c *crc32Field) crc(curOffset int, buf []byte) (uint32, error) ***REMOVED***
	var tab *crc32.Table
	switch c.polynomial ***REMOVED***
	case crcIEEE:
		tab = crc32.IEEETable
	case crcCastagnoli:
		tab = castagnoliTable
	default:
		return 0, PacketDecodingError***REMOVED***"invalid CRC type"***REMOVED***
	***REMOVED***
	return crc32.Checksum(buf[c.startOffset+4:curOffset], tab), nil
***REMOVED***
