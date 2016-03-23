package message

import (
	"bytes"
	"encoding/json"
)

const ClientTopic string = "client" // Subscription topic for clients
const MasterTopic string = "master" // Subscription topic for masters
const WorkerTopic string = "worker" // Subscription topic for workers

// A directed message.
type Message struct ***REMOVED***
	Topic  string `json:"-"`
	Type   string `json:"type"`
	Fields Fields `json:"fields"`
***REMOVED***

// A set of fields in a message.
type Fields map[string]interface***REMOVED******REMOVED***

// Creates a message directed to the master server.
func NewToMaster(t string, f Fields) Message ***REMOVED***
	return Message***REMOVED***
		Topic:  MasterTopic,
		Type:   t,
		Fields: f,
	***REMOVED***
***REMOVED***

// Creates a message directed to clients.
func NewToClient(t string, f Fields) Message ***REMOVED***
	return Message***REMOVED***
		Topic:  ClientTopic,
		Type:   t,
		Fields: f,
	***REMOVED***
***REMOVED***

// Creates a message directed to workers.
func NewToWorker(t string, f Fields) Message ***REMOVED***
	return Message***REMOVED***
		Topic:  ClientTopic,
		Type:   t,
		Fields: f,
	***REMOVED***
***REMOVED***

// Decodes a message from wire format.
func Decode(data []byte, msg *Message) (err error) ***REMOVED***
	nullIndex := bytes.IndexByte(data, 0)
	msg.Topic = string(data[:nullIndex])
	return json.Unmarshal(data[nullIndex+1:], msg)
***REMOVED***

// Encodes a message to wire format.
func (msg *Message) Encode() ([]byte, error) ***REMOVED***
	jdata, err := json.Marshal(msg)
	if err != nil ***REMOVED***
		return jdata, err
	***REMOVED***

	buf := bytes.NewBufferString(msg.Topic)
	buf.WriteByte(0)
	buf.Write(jdata)
	return buf.Bytes(), nil
***REMOVED***
