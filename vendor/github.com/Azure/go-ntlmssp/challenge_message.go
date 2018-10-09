package ntlmssp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type challengeMessageFields struct ***REMOVED***
	messageHeader
	TargetName      varField
	NegotiateFlags  negotiateFlags
	ServerChallenge [8]byte
	_               [8]byte
	TargetInfo      varField
***REMOVED***

func (m challengeMessageFields) IsValid() bool ***REMOVED***
	return m.messageHeader.IsValid() && m.MessageType == 2
***REMOVED***

type challengeMessage struct ***REMOVED***
	challengeMessageFields
	TargetName    string
	TargetInfo    map[avID][]byte
	TargetInfoRaw []byte
***REMOVED***

func (m *challengeMessage) UnmarshalBinary(data []byte) error ***REMOVED***
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, &m.challengeMessageFields)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !m.challengeMessageFields.IsValid() ***REMOVED***
		return fmt.Errorf("Message is not a valid challenge message: %+v", m.challengeMessageFields.messageHeader)
	***REMOVED***

	if m.challengeMessageFields.TargetName.Len > 0 ***REMOVED***
		m.TargetName, err = m.challengeMessageFields.TargetName.ReadStringFrom(data, m.NegotiateFlags.Has(negotiateFlagNTLMSSPNEGOTIATEUNICODE))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if m.challengeMessageFields.TargetInfo.Len > 0 ***REMOVED***
		d, err := m.challengeMessageFields.TargetInfo.ReadFrom(data)
		m.TargetInfoRaw = d
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		m.TargetInfo = make(map[avID][]byte)
		r := bytes.NewReader(d)
		for ***REMOVED***
			var id avID
			var l uint16
			err = binary.Read(r, binary.LittleEndian, &id)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if id == avIDMsvAvEOL ***REMOVED***
				break
			***REMOVED***

			err = binary.Read(r, binary.LittleEndian, &l)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			value := make([]byte, l)
			n, err := r.Read(value)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if n != int(l) ***REMOVED***
				return fmt.Errorf("Expected to read %d bytes, got only %d", l, n)
			***REMOVED***
			m.TargetInfo[id] = value
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
