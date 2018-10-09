package ntlmssp

import (
	"bytes"
)

var signature = [8]byte***REMOVED***'N', 'T', 'L', 'M', 'S', 'S', 'P', 0***REMOVED***

type messageHeader struct ***REMOVED***
	Signature   [8]byte
	MessageType uint32
***REMOVED***

func (h messageHeader) IsValid() bool ***REMOVED***
	return bytes.Equal(h.Signature[:], signature[:]) &&
		h.MessageType > 0 && h.MessageType < 4
***REMOVED***

func newMessageHeader(messageType uint32) messageHeader ***REMOVED***
	return messageHeader***REMOVED***signature, messageType***REMOVED***
***REMOVED***
