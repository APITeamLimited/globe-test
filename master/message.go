package master

import (
	"encoding/json"
)

type Message struct ***REMOVED***
	Type string `json:"type"`
	Body string `json:"body"`
***REMOVED***

func DecodeMessage(data []byte) (msg Message, err error) ***REMOVED***
	err = json.Unmarshal(data, msg)
	return msg, err
***REMOVED***

func (msg *Message) Encode() ([]byte, error) ***REMOVED***
	return json.Marshal(msg)
***REMOVED***
