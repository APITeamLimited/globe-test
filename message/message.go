package message

import (
	"bytes"
	"encoding/json"
)

const ClientTopic string = "client"
const MasterTopic string = "master"
const WorkerTopic string = "worker"

type Message struct ***REMOVED***
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Body  string `json:"body"`
***REMOVED***

func NewToMaster(t string, b string) Message ***REMOVED***
	return Message***REMOVED***
		Topic: MasterTopic,
		Type:  t,
		Body:  b,
	***REMOVED***
***REMOVED***

func NewToClient(t string, b string) Message ***REMOVED***
	return Message***REMOVED***
		Topic: ClientTopic,
		Type:  t,
		Body:  b,
	***REMOVED***
***REMOVED***

func NewToWorker(t string, b string) Message ***REMOVED***
	return Message***REMOVED***
		Topic: ClientTopic,
		Type:  t,
		Body:  b,
	***REMOVED***
***REMOVED***

func Decode(data []byte, msg interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	nullIndex := bytes.IndexByte(data, 0)
	return json.Unmarshal(data[nullIndex+1:], msg)
***REMOVED***

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
