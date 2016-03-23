package master

import (
	"testing"
)

func TestEncodeDecode(t *testing.T) ***REMOVED***
	msg1 := Message***REMOVED***
		Type: "test",
		Body: "Abc123",
	***REMOVED***
	data, err := msg1.Encode()
	if err != nil ***REMOVED***
		t.Errorf("Couldn't encode message: %s", err)
	***REMOVED***
	if string(data) != `***REMOVED***"type":"test","body":"Abc123"***REMOVED***` ***REMOVED***
		t.Errorf("Incorrect JSON representation: %s", string(data))
	***REMOVED***

	msg2 := &Message***REMOVED******REMOVED***
	err = DecodeMessage(data, msg2)
	if err != nil ***REMOVED***
		t.Errorf("Couldn't decode: %s", err)
	***REMOVED***

	if msg2.Type != msg1.Type ***REMOVED***
		t.Errorf("Type mismatch: %s != %s", msg2.Type, msg1.Type)
	***REMOVED***
	if msg2.Body != msg2.Body ***REMOVED***
		t.Errorf("Body mismatch: \"%s\" != \"%s\"", msg2.Body, msg1.Body)
	***REMOVED***
***REMOVED***
