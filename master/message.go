package master

import (
	"encoding/json"
)

type Message struct ***REMOVED***
	Type string `json:"type"`
	Body string `json:"body"`
***REMOVED***

func DecodeMessage(data []byte, msg interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	return json.Unmarshal(data, msg)
***REMOVED***

func (msg *Message) Encode() ([]byte, error) ***REMOVED***
	return json.Marshal(msg)
***REMOVED***
