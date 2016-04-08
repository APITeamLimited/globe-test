package ping

import (
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/comm"
	"github.com/loadimpact/speedboat/util/testutil"
	"testing"
)

func TestParseNothing(t *testing.T) ***REMOVED***
	err := testutil.WithAppContext("", Command, func(c *cli.Context) ***REMOVED***
		topic, local, err := Parse(c)
		if err != nil ***REMOVED***
			t.Error("Error:", err)
		***REMOVED***
		if topic != comm.MasterTopic ***REMOVED***
			t.Error("Default topic not master", topic)
		***REMOVED***
		if local ***REMOVED***
			t.Error("Default allows local")
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestParseWorker(t *testing.T) ***REMOVED***
	err := testutil.WithAppContext("--worker", Command, func(c *cli.Context) ***REMOVED***
		topic, _, err := Parse(c)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		if topic != comm.WorkerTopic ***REMOVED***
			t.Fail()
		***REMOVED***
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***

func TestParseLocal(t *testing.T) ***REMOVED***
	err := testutil.WithAppContext("--local", Command, func(c *cli.Context) ***REMOVED***
		_, local, err := Parse(c)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		if !local ***REMOVED***
			t.Fail()
		***REMOVED***
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
***REMOVED***
