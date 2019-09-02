package js

import (
	"context"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
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
		Throw:        null.BoolFrom(true),
		MaxRedirects: null.IntFrom(10),
		Hosts:        tb.Dialer.Hosts,
	***REMOVED***)

	var ch = make(chan stats.SampleContainer, 100)
	go func() ***REMOVED*** // read the channel so it doesn't block
		for ***REMOVED***
			<-ch
		***REMOVED***
	***REMOVED***()
	vu, err := r.NewVU(ch)
	if !assert.NoError(b, err) ***REMOVED***
		return
	***REMOVED***
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		err = vu.RunOnce(context.Background())
		assert.NoError(b, err)
	***REMOVED***
***REMOVED***
