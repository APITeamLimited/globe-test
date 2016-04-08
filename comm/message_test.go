package comm

import (
	"testing"
)

func TestEncodeDecode(t *testing.T) ***REMOVED***
	msg1 := Message***REMOVED***
		Topic:   "test",
		Type:    "test",
		Payload: []byte("test"),
	***REMOVED***
	data, err := msg1.Encode()
	msg2 := &Message***REMOVED******REMOVED***
	err = Decode(data, msg2)
	if err != nil ***REMOVED***
		t.Errorf("Couldn't decode: %s", err)
	***REMOVED***

	if msg2.Topic != msg1.Topic ***REMOVED***
		t.Errorf("Topic mismatch: %s != %s", msg2.Topic, msg1.Topic)
	***REMOVED***
	if msg2.Type != msg1.Type ***REMOVED***
		t.Errorf("Type mismatch: %s != %s", msg2.Type, msg1.Type)
	***REMOVED***
	if string(msg2.Payload) != string(msg1.Payload) ***REMOVED***
		t.Errorf("Payload mismatch: \"%s\" != \"%s\"", string(msg2.Payload), string(msg1.Payload))
	***REMOVED***
***REMOVED***
