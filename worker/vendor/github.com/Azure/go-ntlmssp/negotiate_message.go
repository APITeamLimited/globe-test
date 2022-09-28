package ntlmssp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

const expMsgBodyLen = 40

type negotiateMessageFields struct ***REMOVED***
	messageHeader
	NegotiateFlags negotiateFlags

	Domain      varField
	Workstation varField

	Version
***REMOVED***

var defaultFlags = negotiateFlagNTLMSSPNEGOTIATETARGETINFO |
	negotiateFlagNTLMSSPNEGOTIATE56 |
	negotiateFlagNTLMSSPNEGOTIATE128 |
	negotiateFlagNTLMSSPNEGOTIATEUNICODE |
	negotiateFlagNTLMSSPNEGOTIATEEXTENDEDSESSIONSECURITY

//NewNegotiateMessage creates a new NEGOTIATE message with the
//flags that this package supports.
func NewNegotiateMessage(domainName, workstationName string) ([]byte, error) ***REMOVED***
	payloadOffset := expMsgBodyLen
	flags := defaultFlags

	if domainName != "" ***REMOVED***
		flags |= negotiateFlagNTLMSSPNEGOTIATEOEMDOMAINSUPPLIED
	***REMOVED***

	if workstationName != "" ***REMOVED***
		flags |= negotiateFlagNTLMSSPNEGOTIATEOEMWORKSTATIONSUPPLIED
	***REMOVED***

	msg := negotiateMessageFields***REMOVED***
		messageHeader:  newMessageHeader(1),
		NegotiateFlags: flags,
		Domain:         newVarField(&payloadOffset, len(domainName)),
		Workstation:    newVarField(&payloadOffset, len(workstationName)),
		Version:        DefaultVersion(),
	***REMOVED***

	b := bytes.Buffer***REMOVED******REMOVED***
	if err := binary.Write(&b, binary.LittleEndian, &msg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if b.Len() != expMsgBodyLen ***REMOVED***
		return nil, errors.New("incorrect body length")
	***REMOVED***

	payload := strings.ToUpper(domainName + workstationName)
	if _, err := b.WriteString(payload); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return b.Bytes(), nil
***REMOVED***
