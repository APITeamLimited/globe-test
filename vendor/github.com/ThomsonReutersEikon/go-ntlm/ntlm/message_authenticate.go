//Copyright 2013 Thomson Reuters Global Resources. BSD License please see License file for more information

package ntlm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type AuthenticateMessage struct ***REMOVED***
	// sig - 8 bytes
	Signature []byte
	// message type - 4 bytes
	MessageType uint32

	// The LmChallenge Response can be v1 or v2
	LmChallengeResponse *PayloadStruct // 8 bytes
	LmV1Response        *LmV1Response
	LmV2Response        *LmV2Response

	// The NtChallengeResponse can be v1 or v2
	NtChallengeResponseFields *PayloadStruct // 8 bytes
	NtlmV1Response            *NtlmV1Response
	NtlmV2Response            *NtlmV2Response

	DomainName  *PayloadStruct // 8 bytes
	UserName    *PayloadStruct // 8 bytes
	Workstation *PayloadStruct // 8 bytes

	// If the NTLMSSP_NEGOTIATE_KEY_EXCH flag is set in the neogitate flags then this will point to the offset in the payload
	// with the key, otherwise it will have Len = 0. According to Davenport these bytes are optional (see Type3 message).
	// The MS-NLMP docs do not mention this.
	EncryptedRandomSessionKey *PayloadStruct // 8 bytes

	/// MS-NLMP 2.2.1.3 - In connectionless mode, a NEGOTIATE structure that contains a set of bit flags (section 2.2.2.5) and represents the
	// conclusion of negotiation—the choices the client has made from the options the server offered in the CHALLENGE_MESSAGE.
	// In connection-oriented mode, a NEGOTIATE structure that contains the set of bit flags (section 2.2.2.5) negotiated in
	// the previous
	NegotiateFlags uint32 // 4 bytes

	// Version (8 bytes): A VERSION structure (section 2.2.2.10) that is present only when the NTLMSSP_NEGOTIATE_VERSION
	// flag is set in the NegotiateFlags field. This structure is used for debugging purposes only. In normal protocol
	// messages, it is ignored and does not affect the NTLM message processing.<9>
	Version *VersionStruct

	// The message integrity for the NTLM NEGOTIATE_MESSAGE, CHALLENGE_MESSAGE, and AUTHENTICATE_MESSAGE.<10>
	Mic []byte // 16 bytes

	// payload - variable
	Payload []byte
***REMOVED***

func ParseAuthenticateMessage(body []byte, ntlmVersion int) (*AuthenticateMessage, error) ***REMOVED***
	am := new(AuthenticateMessage)

	am.Signature = body[0:8]
	if !bytes.Equal(am.Signature, []byte("NTLMSSP\x00")) ***REMOVED***
		return nil, errors.New("Invalid NTLM message signature")
	***REMOVED***

	am.MessageType = binary.LittleEndian.Uint32(body[8:12])
	if am.MessageType != 3 ***REMOVED***
		return nil, errors.New("Invalid NTLM message type should be 0x00000003 for authenticate message")
	***REMOVED***

	var err error

	am.LmChallengeResponse, err = ReadBytePayload(12, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if ntlmVersion == 2 ***REMOVED***
		am.LmV2Response = ReadLmV2Response(am.LmChallengeResponse.Payload)
	***REMOVED*** else ***REMOVED***
		am.LmV1Response = ReadLmV1Response(am.LmChallengeResponse.Payload)
	***REMOVED***

	am.NtChallengeResponseFields, err = ReadBytePayload(20, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Check to see if this is a v1 or v2 response
	if ntlmVersion == 2 ***REMOVED***
		am.NtlmV2Response, err = ReadNtlmV2Response(am.NtChallengeResponseFields.Payload)
	***REMOVED*** else ***REMOVED***
		am.NtlmV1Response, err = ReadNtlmV1Response(am.NtChallengeResponseFields.Payload)
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	am.DomainName, err = ReadStringPayload(28, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	am.UserName, err = ReadStringPayload(36, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	am.Workstation, err = ReadStringPayload(44, body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	lowestOffset := am.getLowestPayloadOffset()
	offset := 52

	// If the lowest payload offset is 52 then:
	// The Session Key, flags, and OS Version structure are omitted. The data (payload) block in this case starts after the Workstation Name
	// security buffer header, at offset 52. This form is seen in older Win9x-based systems. This is from the davenport notes about Type 3
	// messages and this information does not seem to be present in the MS-NLMP document
	if lowestOffset > 52 ***REMOVED***
		am.EncryptedRandomSessionKey, err = ReadBytePayload(offset, body)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		offset = offset + 8

		am.NegotiateFlags = binary.LittleEndian.Uint32(body[offset : offset+4])
		offset = offset + 4

		// Version (8 bytes): A VERSION structure (section 2.2.2.10) that is present only when the NTLMSSP_NEGOTIATE_VERSION flag is set in the NegotiateFlags field. This structure is used for debugging purposes only. In normal protocol messages, it is ignored and does not affect the NTLM message processing.<9>
		if NTLMSSP_NEGOTIATE_VERSION.IsSet(am.NegotiateFlags) ***REMOVED***
			am.Version, err = ReadVersionStruct(body[offset : offset+8])
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			offset = offset + 8
		***REMOVED***

		// The MS-NLMP has this to say about the MIC
		//   "An AUTHENTICATE_MESSAGE indicates the presence of a MIC field if the TargetInfo field has an AV_PAIR structure whose two fields are:
		//   AvId == MsvAvFlags Value bit 0x2 == 1"
		// However there is no TargetInfo structure in the Authenticate Message! There is one in the Challenge Message though. So I'm using
		// a hack to check to see if there is a MIC. I look to see if there is room for the MIC before the payload starts. If so I assume
		// there is a MIC and read it out.
		var lowestOffset = am.getLowestPayloadOffset()
		if lowestOffset > offset ***REMOVED***
			// MIC - 16 bytes
			am.Mic = body[offset : offset+16]
			offset = offset + 16
		***REMOVED***
	***REMOVED***

	am.Payload = body[offset:]

	return am, nil
***REMOVED***

func (a *AuthenticateMessage) ClientChallenge() (response []byte) ***REMOVED***
	if a.NtlmV2Response != nil ***REMOVED***
		response = a.NtlmV2Response.NtlmV2ClientChallenge.ChallengeFromClient
	***REMOVED*** else if a.NtlmV1Response != nil && NTLMSSP_NEGOTIATE_EXTENDED_SESSIONSECURITY.IsSet(a.NegotiateFlags) ***REMOVED***
		response = a.LmV1Response.Response[0:8]
	***REMOVED***

	return response
***REMOVED***

func (a *AuthenticateMessage) getLowestPayloadOffset() int ***REMOVED***
	payloadStructs := [...]*PayloadStruct***REMOVED***a.LmChallengeResponse, a.NtChallengeResponseFields, a.DomainName, a.UserName, a.Workstation, a.EncryptedRandomSessionKey***REMOVED***

	// Find the lowest offset value
	lowest := 9999
	for i := range payloadStructs ***REMOVED***
		p := payloadStructs[i]
		if p != nil && p.Offset > 0 && int(p.Offset) < lowest ***REMOVED***
			lowest = int(p.Offset)
		***REMOVED***
	***REMOVED***

	return lowest
***REMOVED***

func (a *AuthenticateMessage) Bytes() []byte ***REMOVED***
	payloadLen := int(a.LmChallengeResponse.Len + a.NtChallengeResponseFields.Len + a.DomainName.Len + a.UserName.Len + a.Workstation.Len + a.EncryptedRandomSessionKey.Len)
	messageLen := 8 + 4 + 6*8 + 4 + 8 + 16
	payloadOffset := uint32(messageLen)

	messageBytes := make([]byte, 0, messageLen+payloadLen)
	buffer := bytes.NewBuffer(messageBytes)

	buffer.Write(a.Signature)

	binary.Write(buffer, binary.LittleEndian, a.MessageType)

	a.LmChallengeResponse.Offset = payloadOffset
	payloadOffset += uint32(a.LmChallengeResponse.Len)
	buffer.Write(a.LmChallengeResponse.Bytes())

	a.NtChallengeResponseFields.Offset = payloadOffset
	payloadOffset += uint32(a.NtChallengeResponseFields.Len)
	buffer.Write(a.NtChallengeResponseFields.Bytes())

	a.DomainName.Offset = payloadOffset
	payloadOffset += uint32(a.DomainName.Len)
	buffer.Write(a.DomainName.Bytes())

	a.UserName.Offset = payloadOffset
	payloadOffset += uint32(a.UserName.Len)
	buffer.Write(a.UserName.Bytes())

	a.Workstation.Offset = payloadOffset
	payloadOffset += uint32(a.Workstation.Len)
	buffer.Write(a.Workstation.Bytes())

	a.EncryptedRandomSessionKey.Offset = payloadOffset
	payloadOffset += uint32(a.EncryptedRandomSessionKey.Len)
	buffer.Write(a.EncryptedRandomSessionKey.Bytes())

	buffer.Write(uint32ToBytes(a.NegotiateFlags))

	if a.Version != nil ***REMOVED***
		buffer.Write(a.Version.Bytes())
	***REMOVED*** else ***REMOVED***
		buffer.Write(make([]byte, 8))
	***REMOVED***

	if a.Mic != nil ***REMOVED***
		buffer.Write(a.Mic)
	***REMOVED*** else ***REMOVED***
		buffer.Write(make([]byte, 16))
	***REMOVED***

	// Write out the payloads
	buffer.Write(a.LmChallengeResponse.Payload)
	buffer.Write(a.NtChallengeResponseFields.Payload)
	buffer.Write(a.DomainName.Payload)
	buffer.Write(a.UserName.Payload)
	buffer.Write(a.Workstation.Payload)
	buffer.Write(a.EncryptedRandomSessionKey.Payload)

	return buffer.Bytes()
***REMOVED***

func (a *AuthenticateMessage) String() string ***REMOVED***
	var buffer bytes.Buffer

	buffer.WriteString("Authenticate NTLM Message\n")
	buffer.WriteString(fmt.Sprintf("Payload Offset: %d Length: %d\n", a.getLowestPayloadOffset(), len(a.Payload)))

	if a.LmV2Response != nil ***REMOVED***
		buffer.WriteString(a.LmV2Response.String())
		buffer.WriteString("\n")
	***REMOVED***

	if a.LmV1Response != nil ***REMOVED***
		buffer.WriteString(a.LmV1Response.String())
		buffer.WriteString("\n")
	***REMOVED***

	if a.NtlmV2Response != nil ***REMOVED***
		buffer.WriteString(a.NtlmV2Response.String())
		buffer.WriteString("\n")
	***REMOVED***

	if a.NtlmV1Response != nil ***REMOVED***
		buffer.WriteString(fmt.Sprintf("NtlmResponse Length: %d\n", a.NtChallengeResponseFields.Len))
		buffer.WriteString(a.NtlmV1Response.String())
		buffer.WriteString("\n")
	***REMOVED***

	buffer.WriteString(fmt.Sprintf("UserName: %s\n", a.UserName.String()))
	buffer.WriteString(fmt.Sprintf("DomainName: %s\n", a.DomainName.String()))
	buffer.WriteString(fmt.Sprintf("Workstation: %s\n", a.Workstation.String()))

	if a.EncryptedRandomSessionKey != nil ***REMOVED***
		buffer.WriteString(fmt.Sprintf("EncryptedRandomSessionKey: %s\n", a.EncryptedRandomSessionKey.String()))
	***REMOVED***

	if a.Version != nil ***REMOVED***
		buffer.WriteString(fmt.Sprintf("Version: %s\n", a.Version.String()))
	***REMOVED***

	if a.Mic != nil ***REMOVED***
		buffer.WriteString(fmt.Sprintf("MIC: %s\n", hex.EncodeToString(a.Mic)))
	***REMOVED***

	buffer.WriteString(fmt.Sprintf("Flags %d\n", a.NegotiateFlags))
	buffer.WriteString(FlagsToString(a.NegotiateFlags))

	return buffer.String()
***REMOVED***
