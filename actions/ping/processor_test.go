package ping

import (
	"github.com/loadimpact/speedboat/comm"
	"testing"
	"time"
)

func TestProcessPing(t *testing.T) ***REMOVED***
	p := PingProcessor***REMOVED******REMOVED***
	now := time.Now()
	res := <-p.ProcessPing(PingMessage***REMOVED***Time: now***REMOVED***)

	if res.Topic != comm.ClientTopic ***REMOVED***
		t.Error("Message not to client:", res.Topic)
	***REMOVED***
	if res.Type != "ping.pong" ***REMOVED***
		t.Error("Wrong message type:", res.Type)
	***REMOVED***

	data := PingMessage***REMOVED******REMOVED***
	if err := res.Take(&data); err != nil ***REMOVED***
		t.Fatal("Couldn't decode pong:", err)
	***REMOVED***

	if data.Time.Unix() != now.Unix() ***REMOVED***
		t.Error("Wrong timestamp:", data.Time)
	***REMOVED***
***REMOVED***
