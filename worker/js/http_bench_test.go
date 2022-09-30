package js

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils/httpmultibin"
)

func BenchmarkHTTPRequests(b *testing.B) ***REMOVED***
	b.StopTimer()
	tb := httpmultibin.NewHTTPMultiBin(b)

	r, err := getSimpleRunner(b, "/script.js", tb.Replacer.Replace(`
			import http from "k6/http";
			export default function() ***REMOVED***
				let url = "HTTPBIN_URL";
				let res = http.get(url + "/cookies/set?k2=v2&k1=v1");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			***REMOVED***
		`), libWorker.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom("extended")***REMOVED***)
	require.NoError(b, err)
	err = r.SetOptions(libWorker.Options***REMOVED***
		Throw:          null.BoolFrom(true),
		MaxRedirects:   null.IntFrom(10),
		Hosts:          tb.Dialer.Hosts,
		NoCookiesReset: null.BoolFrom(true),
	***REMOVED***)
	require.NoError(b, err)

	ch := make(chan workerMetrics.SampleContainer, 100)
	defer close(ch)
	go func() ***REMOVED*** // read the channel so it doesn't block
		for range ch ***REMOVED***
		***REMOVED***
	***REMOVED***()
	initVU, err := r.NewVU(1, 1, ch, libWorker.GetTestWorkerInfo())
	require.NoError(b, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&libWorker.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		assert.NoError(b, vu.RunOnce())
	***REMOVED***
***REMOVED***

func BenchmarkHTTPRequestsBase(b *testing.B) ***REMOVED***
	b.StopTimer()
	tb := httpmultibin.NewHTTPMultiBin(b)

	r, err := getSimpleRunner(b, "/script.js", tb.Replacer.Replace(`
			var http = require("k6/http");
			exports.default = function() ***REMOVED***
				var url = "HTTPBIN_URL";
				var res = http.get(url + "/cookies/set?k2=v2&k1=v1");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status) ***REMOVED***
			***REMOVED***
		`))
	require.NoError(b, err)
	err = r.SetOptions(libWorker.Options***REMOVED***
		Throw:          null.BoolFrom(true),
		MaxRedirects:   null.IntFrom(10),
		Hosts:          tb.Dialer.Hosts,
		NoCookiesReset: null.BoolFrom(true),
	***REMOVED***)
	require.NoError(b, err)

	ch := make(chan workerMetrics.SampleContainer, 100)
	defer close(ch)
	go func() ***REMOVED*** // read the channel so it doesn't block
		for range ch ***REMOVED***
		***REMOVED***
	***REMOVED***()
	initVU, err := r.NewVU(1, 1, ch, libWorker.GetTestWorkerInfo())
	require.NoError(b, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vu := initVU.Activate(&libWorker.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
	b.StartTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		assert.NoError(b, vu.RunOnce())
	***REMOVED***
***REMOVED***
