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
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
)

func BenchmarkAggregateHTTP(b *testing.B) ***REMOVED***
	script := &loader.SourceData***REMOVED***
		Data: []byte(""),
		URL:  &url.URL***REMOVED***Path: "/script.js"***REMOVED***,
	***REMOVED***

	options := lib.Options***REMOVED***
		Duration: types.NullDurationFrom(1 * time.Second),
	***REMOVED***

	config := NewConfig().Apply(Config***REMOVED***
		NoCompress:              null.BoolFrom(true),
		AggregationCalcInterval: types.NullDurationFrom(time.Millisecond * 200),
		AggregationPeriod:       types.NullDurationFrom(time.Millisecond * 200),
	***REMOVED***)
	collector, err := New(config, script, options, []lib.ExecutionStep***REMOVED******REMOVED***, "1.0")
	require.NoError(b, err)
	now := time.Now()
	collector.referenceID = "something"
	var containersCount = time.Duration(500000)

	for _, tagCount := range []int***REMOVED***1, 5, 35, 315, 3645***REMOVED*** ***REMOVED***
		tagCount := tagCount
		b.Run(fmt.Sprintf("tags:%d", tagCount), func(b *testing.B) ***REMOVED***
			b.ResetTimer()
			for s := 0; s < b.N; s++ ***REMOVED***
				b.StopTimer()
				var container = make([]stats.SampleContainer, containersCount)
				for i := time.Duration(1); i <= containersCount; i++ ***REMOVED***
					var status = "200"
					if int(i)%tagCount%7 == 6 ***REMOVED***
						status = "404"
					***REMOVED*** else if int(i)%tagCount%7 == 5 ***REMOVED***
						status = "500"
					***REMOVED***

					container[i-1] = &httpext.Trail***REMOVED***
						Blocked:        i % 200 * 100 * time.Millisecond,
						Connecting:     i % 200 * 200 * time.Millisecond,
						TLSHandshaking: i % 200 * 300 * time.Millisecond,
						Sending:        i % 200 * 400 * time.Millisecond,
						Waiting:        500 * time.Millisecond,
						Receiving:      600 * time.Millisecond,
						EndTime:        now.Add(i % 100 * 100),
						ConnDuration:   500 * time.Millisecond,
						Duration:       i % 150 * 1500 * time.Millisecond,
						Tags: stats.IntoSampleTags(&map[string]string***REMOVED***
							"test": "mest", "a": "b",
							"custom": fmt.Sprintf("group%d", int(i)%tagCount%9),
							"group":  fmt.Sprintf("group%d", int(i)%tagCount%5),
							"status": status,
							"url":    fmt.Sprintf("something%d", int(i)%tagCount%11),
							"name":   fmt.Sprintf("else%d", int(i)%tagCount%11),
						***REMOVED***),
					***REMOVED***
				***REMOVED***
				collector.Collect(container)
				b.StartTimer()
				collector.aggregateHTTPTrails(time.Millisecond * 200)
				collector.bufferSamples = nil
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
