package js

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/stats"
)

func BenchmarkHTTPRequests(b *testing.B) ***REMOVED***
	b.StopTimer()
	tb := httpmultibin.NewHTTPMultiBin(b)
	defer tb.Cleanup()

	r, err := getSimpleRunner("/script.js", tb.Replacer.Replace(`
			import http from "k6/http";
			export default function() ***REMOVED***
				let url = "HTTPBIN_URL";
				let res = http.get(url + "/cookies/set?k2=v2&k1=v1");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			***REMOVED***
		`))
	if !assert.NoError(b, err) ***REMOVED***
		return
	***REMOVED***
	r.SetOptions(lib.Options***REMOVED***
		Throw:          null.BoolFrom(true),
		MaxRedirects:   null.IntFrom(10),
		Hosts:          tb.Dialer.Hosts,
		NoCookiesReset: null.BoolFrom(true),
	***REMOVED***)

	var ch = make(chan stats.SampleContainer, 100)
	go func() ***REMOVED*** // read the channel so it doesn't block
		for ***REMOVED***
			<-ch
		***REMOVED***
	***REMOVED***()
	initVU, err := r.NewVU(1, ch)
	if !assert.NoError(b, err) ***REMOVED***
		return
	***REMOVED***
	vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: context.Background()***REMOVED***)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		err = vu.RunOnce()
		assert.NoError(b, err)
	***REMOVED***
***REMOVED***
