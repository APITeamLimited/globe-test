//Copyright 2013 Thomson Reuters Global Resources. BSD License please see License file for more information

package ntlm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type AvPairType uint16

// MS-NLMP - 2.2.2.1 AV_PAIR
const (
	// Indicates that this is the last AV_PAIR in the list. AvLen MUST be 0. This type of information MUST be present in the AV pair list.
	MsvAvEOL AvPairType = iota
	// The server's NetBIOS computer name. The name MUST be in Unicode, and is not null-terminated. This type of information MUST be present in the AV_pair list.
	MsvAvNbComputerName
	// The server's NetBIOS domain name. The name MUST be in Unicode, and is not null-terminated. This type of information MUST be present in the AV_pair list.
	MsvAvNbDomainName
	// The fully qualified domain name (FQDN (1)) of the computer. The name MUST be in Unicode, and is not null-terminated.
	MsvAvDnsComputerName
	// The FQDN (2) of the domain. The name MUST be in Unicode, and is not null-terminate.
	MsvAvDnsDomainName
	// The FQDN (2) of the forest. The name MUST be in Unicode, and is not null-terminated.<11>
	MsvAvDnsTreeName
	// A 32-bit value indicating server or client configuration.
	// 0x00000001: indicates to the client that the account authentication is constrained.
	// 0x00000002: indicates that the client is providing message integrity in the MIC field (section 2.2.1.3) in the AUTHENTICATE_MESSAGE.<12>
	// 0x00000004: indicates that the client is providing a target SPN generated from an untrusted source.<13>
	MsvAvFlags
	// A FILETIME structure ([MS-DTYP] section 2.3.1) in little-endian byte order that contains the server local time.<14>
	MsvAvTimestamp
	//A Restriction_Encoding (section 2.2.2.2) structure. The Value field contains a structure representing the integrity level of the security principal, as well as a MachineID created at computer startup to identify the calling machine.<15>
	MsAvRestrictions
	// The SPN of the target server. The name MUST be in Unicode and is not null-terminated.<16>
	MsvAvTargetName
	// annel bindings hash. The Value field contains an MD5 hash ([RFC4121] section 4.1.1.2) of a gss_channel_bindings_struct ([RFC2744] section 3.11).
	// An all-zero value of the hash is used to indicate absence of channel bindings.<17>
	MsvChannelBindings
)

// Helper struct that contains a list of AvPairs with helper methods for running through them
type AvPairs struct ***REMOVED***
	List []AvPair
***REMOVED***

func (p *AvPairs) AddAvPair(avId AvPairType, bytes []byte) ***REMOVED***
	a := &AvPair***REMOVED***AvId: avId, AvLen: uint16(len(bytes)), Value: bytes***REMOVED***
	p.List = append(p.List, *a)
***REMOVED***

func ReadAvPairs(data []byte) *AvPairs ***REMOVED***
	pairs := new(AvPairs)

	// Get the number of AvPairs and allocate enough AvPair structures to hold them
	offset := 0
	for i := 0; len(data) > 0 && i < 11; i++ ***REMOVED***
		pair := ReadAvPair(data, offset)
		offset = offset + 4 + int(pair.AvLen)
		pairs.List = append(pairs.List, *pair)
		if pair.AvId == MsvAvEOL ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return pairs
***REMOVED***

func (p *AvPairs) Bytes() (result []byte) ***REMOVED***
	totalLength := 0
	for i := range p.List ***REMOVED***
		a := p.List[i]
		totalLength = totalLength + int(a.AvLen) + 4
	***REMOVED***

	result = make([]byte, 0, totalLength)
	for i := range p.List ***REMOVED***
		a := p.List[i]
		result = append(result, a.Bytes()...)
	***REMOVED***

	return result
***REMOVED***

func (p *AvPairs) String() string ***REMOVED***
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Av Pairs (Total %d pairs)\n", len(p.List)))

	for i := range p.List ***REMOVED***
		buffer.WriteString(p.List[i].String())
		buffer.WriteString("\n")
	***REMOVED***

	return buffer.String()
***REMOVED***

func (p *AvPairs) Find(avType AvPairType) (result *AvPair) ***REMOVED***
	for i := range p.List ***REMOVED***
		pair := p.List[i]
		if avType == pair.AvId ***REMOVED***
			result = &pair
			break
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (p *AvPairs) ByteValue(avType AvPairType) (result []byte) ***REMOVED***
	pair := p.Find(avType)
	if pair != nil ***REMOVED***
		result = pair.Value
	***REMOVED***
	return
***REMOVED***

func (p *AvPairs) StringValue(avType AvPairType) (result string) ***REMOVED***
	pair := p.Find(avType)
	if pair != nil ***REMOVED***
		result = pair.UnicodeStringValue()
	***REMOVED***
	return
***REMOVED***

// AvPair as described by MS-NLMP
type AvPair struct ***REMOVED***
	AvId  AvPairType
	AvLen uint16
	Value []byte
***REMOVED***

func ReadAvPair(data []byte, offset int) *AvPair ***REMOVED***
	pair := new(AvPair)
	pair.AvId = AvPairType(binary.LittleEndian.Uint16(data[offset : offset+2]))
	pair.AvLen = binary.LittleEndian.Uint16(data[offset+2 : offset+4])
	pair.Value = data[offset+4 : offset+4+int(pair.AvLen)]
	return pair
***REMOVED***

func (a *AvPair) UnicodeStringValue() string ***REMOVED***
	return utf16ToString(a.Value)
***REMOVED***

func (a *AvPair) Bytes() (result []byte) ***REMOVED***
	result = make([]byte, 4, a.AvLen+4)
	result[0] = byte(a.AvId)
	result[1] = byte(a.AvId >> 8)
	result[2] = byte(a.AvLen)
	result[3] = byte(a.AvLen >> 8)
	result = append(result, a.Value...)
	return
***REMOVED***

func (a *AvPair) String() string ***REMOVED***
	var outString string

	switch a.AvId ***REMOVED***
	case MsvAvEOL:
		outString = "MsvAvEOL"
	case MsvAvNbComputerName:
		outString = "MsAvNbComputerName: " + a.UnicodeStringValue()
	case MsvAvNbDomainName:
		outString = "MsvAvNbDomainName: " + a.UnicodeStringValue()
	case MsvAvDnsComputerName:
		outString = "MsvAvDnsComputerName: " + a.UnicodeStringValue()
	case MsvAvDnsDomainName:
		outString = "MsvAvDnsDomainName: " + a.UnicodeStringValue()
	case MsvAvDnsTreeName:
		outString = "MsvAvDnsTreeName: " + a.UnicodeStringValue()
	case MsvAvFlags:
		outString = "MsvAvFlags: " + hex.EncodeToString(a.Value)
	case MsvAvTimestamp:
		outString = "MsvAvTimestamp: " + hex.EncodeToString(a.Value)
	case MsAvRestrictions:
		outString = "MsAvRestrictions: " + hex.EncodeToString(a.Value)
	case MsvAvTargetName:
		outString = "MsvAvTargetName: " + a.UnicodeStringValue()
	case MsvChannelBindings:
		outString = "MsvChannelBindings: " + hex.EncodeToString(a.Value)
	default:
		outString = fmt.Sprintf("unknown pair type: '%d'", a.AvId)
	***REMOVED***

	return outString
***REMOVED***
