package js

import (
	"github.com/loadimpact/speedboat/loadtest"
	"github.com/loadimpact/speedboat/runner"
	"gopkg.in/olebedev/go-duktape.v2"
	"testing"
	"time"
)

func testJSContext() (*Runner, *duktape.Context, chan runner.Result) ***REMOVED***
	ch := make(chan runner.Result, 100)
	r := New()
	c, err := r.newJSContext(loadtest.LoadTest***REMOVED***
		URL:    "http://example.com",
		Script: "script.js",
		Source: "~ not actually valid JS ~",
		Stages: []loadtest.Stage***REMOVED***
			loadtest.Stage***REMOVED***VUs: loadtest.VUSpec***REMOVED***Start: 10, End: 100***REMOVED***, Duration: 10 * time.Second***REMOVED***,
		***REMOVED***,
	***REMOVED***, 1, ch)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return r, c, ch
***REMOVED***

func TestHTTPSetMaxConnectionsPerHost(t *testing.T) ***REMOVED***
	src := `require('http').setMaxConnectionsPerHost(200);`
	r, c, _ := testJSContext()
	if err := c.PevalString(src); err != nil ***REMOVED***
		t.Fatalf("Couldn't run script: %s", err)
	***REMOVED***
	if r.Client.MaxConnsPerHost != 200 ***REMOVED***
		t.Fatalf("Incorrect number of max connections: %d", r.Client.MaxConnsPerHost)
	***REMOVED***
***REMOVED***

func TestHTTPSetMaxConnectionsPerHostNegative(t *testing.T) ***REMOVED***
	src := `require('http').setMaxConnectionsPerHost(-1);`
	r, c, ch := testJSContext()
	before := r.Client.MaxConnsPerHost
	if err := c.PevalString(src); err != nil ***REMOVED***
		t.Fatalf("Couldn't run script: %s", err)
	***REMOVED***
	select ***REMOVED***
	case res := <-ch:
		if res.Error == nil ***REMOVED***
			t.Error("No error reported!")
		***REMOVED***
		if r.Client.MaxConnsPerHost != before ***REMOVED***
			t.Errorf("Max connections changed! %d", r.Client.MaxConnsPerHost)
		***REMOVED***
	default:
		t.Error("No results")
	***REMOVED***
***REMOVED***
