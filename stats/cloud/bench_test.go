/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
)

// script to clean the logs: `perl -p -e  "s/time=\".*\n//g"`
// TODO: find what sed magic needs to be used to make it work and use it in order to be able to do
// inplace
// TODO: Add a more versatile test not only with metrics that will be aggregated all the time and
// not only httpext.Trail
func BenchmarkCloud(b *testing.B) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(b)
	var maxMetricSamplesPerPackage = 20
	tb.Mux.HandleFunc("/v1/tests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		_, err := fmt.Fprintf(w, `***REMOVED***
			"reference_id": "12",
			"config": ***REMOVED***
				"metricPushInterval": "200ms",
				"aggregationPeriod": "100ms",
				"maxMetricSamplesPerPackage": %d,
				"aggregationCalcInterval": "100ms",
				"aggregationWaitPeriod": "100ms"
			***REMOVED***
		***REMOVED***`, maxMetricSamplesPerPackage)
		require.NoError(b, err)
	***REMOVED***))
	defer tb.Cleanup()

	script := &loader.SourceData***REMOVED***
		Data: []byte(""),
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***

	options := lib.Options***REMOVED***
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	config := NewConfig().Apply(Config***REMOVED***
		Host:       null.StringFrom(tb.ServerHTTP.URL),
		NoCompress: null.BoolFrom(true),
	***REMOVED***)
	collector, err := New(config, script, options, []lib.ExecutionStep***REMOVED******REMOVED***, "1.0")
	require.NoError(b, err)
	now := time.Now()
	tags := stats.IntoSampleTags(&map[string]string***REMOVED***"test": "mest", "a": "b", "url": "something", "name": "else"***REMOVED***)
	var gotTheLimit = false
	var m sync.Mutex

	tb.Mux.HandleFunc(fmt.Sprintf("/v1/metrics/%s", collector.referenceID),
		func(_ http.ResponseWriter, r *http.Request) ***REMOVED***
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(b, err)
			receivedSamples := []Sample***REMOVED******REMOVED***
			assert.NoError(b, json.Unmarshal(body, &receivedSamples))
			assert.True(b, len(receivedSamples) <= maxMetricSamplesPerPackage)
			if len(receivedSamples) == maxMetricSamplesPerPackage ***REMOVED***
				m.Lock()
				gotTheLimit = true
				m.Unlock()
			***REMOVED***
		***REMOVED***)

	require.NoError(b, collector.Init())
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	go func() ***REMOVED***
		collector.Run(ctx)
		wg.Done()
	***REMOVED***()

	for s := 0; s < b.N; s++ ***REMOVED***
		for j := time.Duration(1); j <= 200; j++ ***REMOVED***
			var container = make([]stats.SampleContainer, 0, 500)
			for i := time.Duration(1); i <= 50; i++ ***REMOVED***
				container = append(container, &httpext.Trail***REMOVED***
					Blocked:        i % 200 * 100 * time.Millisecond,
					Connecting:     i % 200 * 200 * time.Millisecond,
					TLSHandshaking: i % 200 * 300 * time.Millisecond,
					Sending:        i * i * 400 * time.Millisecond,
					Waiting:        500 * time.Millisecond,
					Receiving:      600 * time.Millisecond,

					EndTime:      now.Add(i * 100),
					ConnDuration: 500 * time.Millisecond,
					Duration:     j * i * 1500 * time.Millisecond,
					Tags:         stats.NewSampleTags(tags.CloneTags()),
				***REMOVED***)
			***REMOVED***
			collector.Collect(container)
		***REMOVED***
	***REMOVED***

	cancel()
	wg.Wait()
	require.True(b, gotTheLimit)
***REMOVED***
