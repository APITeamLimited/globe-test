//Copyright 2013 Thomson Reuters Global Resources. BSD License please see License file for more information

package ntlm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

const (
	UnicodeStringPayload = iota
	OemStringPayload
	BytesPayload
)

type PayloadStruct struct ***REMOVED***
	Type    int
	Len     uint16
	MaxLen  uint16
	Offset  uint32
	Payload []byte
***REMOVED***

func (p *PayloadStruct) Bytes() []byte ***REMOVED***
	dest := make([]byte, 0, 8)
	buffer := bytes.NewBuffer(dest)

	binary.Write(buffer, binary.LittleEndian, p.Len)
	binary.Write(buffer, binary.LittleEndian, p.MaxLen)
	binary.Write(buffer, binary.LittleEndian, p.Offset)

	return buffer.Bytes()
***REMOVED***

func (p *PayloadStruct) String() string ***REMOVED***
	var returnString string

	switch p.Type ***REMOVED***
	case UnicodeStringPayload:
		returnString = utf16ToString(p.Payload)
	case OemStringPayload:
		returnString = string(p.Payload)
	case BytesPayload:
		returnString = hex.EncodeToString(p.Payload)
	default:
		returnString = "unknown type"
	***REMOVED***
	return returnString
***REMOVED***

func CreateBytePayload(bytes []byte) (*PayloadStruct, error) ***REMOVED***
	p := new(PayloadStruct)
	p.Type = BytesPayload
	p.Len = uint16(len(bytes))
	p.MaxLen = uint16(len(bytes))
	p.Payload = bytes // TODO: Copy these bytes instead of keeping a reference
	return p, nil
***REMOVED***

func CreateStringPayload(value string) (*PayloadStruct, error) ***REMOVED***
	// Create UTF16 unicode bytes from string
	bytes := utf16FromString(value)
	p := new(PayloadStruct)
	p.Type = UnicodeStringPayload
	p.Len = uint16(len(bytes))
	p.MaxLen = uint16(len(bytes))
	p.Payload = bytes // TODO: Copy these bytes instead of keeping a reference
	return p, nil
***REMOVED***

func ReadStringPayload(startByte int, bytes []byte) (*PayloadStruct, error) ***REMOVED***
	return ReadPayloadStruct(startByte, bytes, UnicodeStringPayload)
***REMOVED***

func ReadBytePayload(startByte int, bytes []byte) (*PayloadStruct, error) ***REMOVED***
	return ReadPayloadStruct(startByte, bytes, BytesPayload)
***REMOVED***

func ReadPayloadStruct(startByte int, bytes []byte, PayloadType int) (*PayloadStruct, error) ***REMOVED***
	p := new(PayloadStruct)

	p.Type = PayloadType
	p.Len = binary.LittleEndian.Uint16(bytes[startByte : startByte+2])
	p.MaxLen = binary.LittleEndian.Uint16(bytes[startByte+2 : startByte+4])
	p.Offset = binary.LittleEndian.Uint32(bytes[startByte+4 : startByte+8])

	if p.Len > 0 ***REMOVED***
		endOffset := p.Offset + uint32(p.Len)
		p.Payload = bytes[p.Offset:endOffset]
	***REMOVED***

	return p, nil
***REMOVED***
