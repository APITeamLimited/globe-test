package message

import (
	"testing"
)

func TestEncodeDecode(t *testing.T) ***REMOVED***
	msg1 := Message***REMOVED***
		Topic: "test",
		Type:  "test",
		Body:  "Abc123",
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
	if msg2.Body != msg2.Body ***REMOVED***
		t.Errorf("Body mismatch: \"%s\" != \"%s\"", msg2.Body, msg1.Body)
	***REMOVED***
***REMOVED***
