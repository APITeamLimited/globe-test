package js

import (
	"github.com/loadimpact/speedboat"
	"gopkg.in/olebedev/go-duktape.v2"
	"testing"
	"time"
)

func testJSContext() (*Runner, *duktape.Context) ***REMOVED***
	r := New("")
	c, err := r.newJSContext(speedboat.Test***REMOVED***
		URL:    "http://example.com",
		Script: "script.js",
		Stages: []speedboat.TestStage***REMOVED***
			speedboat.TestStage***REMOVED***StartVUs: 10, EndVUs: 100, Duration: 10 * time.Second***REMOVED***,
		***REMOVED***,
	***REMOVED***, 1)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return r, c
***REMOVED***

func TestHTTPSetMaxConnectionsPerHost(t *testing.T) ***REMOVED***
	src := `require('http').setMaxConnectionsPerHost(200);`
	r, c := testJSContext()
	if err := c.PevalString(src); err != nil ***REMOVED***
		t.Fatalf("Couldn't run script: %s", err)
	***REMOVED***
	if r.Client.MaxConnsPerHost != 200 ***REMOVED***
		t.Fatalf("Incorrect number of max connections: %d", r.Client.MaxConnsPerHost)
	***REMOVED***
***REMOVED***

func TestHTTPSetMaxConnectionsPerHostNegative(t *testing.T) ***REMOVED***
	src := `require('http').setMaxConnectionsPerHost(-1);`
	r, c := testJSContext()
	before := r.Client.MaxConnsPerHost
	if err := c.PevalString(src); err != nil ***REMOVED***
		t.Fatalf("Couldn't run script: %s", err)
	***REMOVED***
	if r.Client.MaxConnsPerHost != before ***REMOVED***
		t.Errorf("Max connections changed! %d", r.Client.MaxConnsPerHost)
	***REMOVED***
***REMOVED***
