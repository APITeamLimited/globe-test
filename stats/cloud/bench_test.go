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

	for _, tagCount := range []int***REMOVED***1, 5, 10, 100, 1000***REMOVED*** ***REMOVED***
		tagCount := tagCount
		b.Run(fmt.Sprintf("tags:%d", tagCount), func(b *testing.B) ***REMOVED***
			tags := make([]*stats.SampleTags, tagCount)
			for i := range tags ***REMOVED***
				tags[i] = stats.IntoSampleTags(&map[string]string***REMOVED***
					"test": "mest", "a": "b",
					"url":  fmt.Sprintf("something%d", i),
					"name": fmt.Sprintf("else%d", i),
				***REMOVED***)
			***REMOVED***
			b.ResetTimer()
			for s := 0; s < b.N; s++ ***REMOVED***
				for j := time.Duration(1); j <= 200; j++ ***REMOVED***
					var container = make([]stats.SampleContainer, 0, 500)
					for i := time.Duration(1); i <= 500; i++ ***REMOVED***
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
							Tags:         stats.NewSampleTags(tags[int(i+j)%len(tags)].CloneTags()),
						***REMOVED***)
					***REMOVED***
					collector.Collect(container)
				***REMOVED***
				collector.aggregateHTTPTrails(time.Millisecond * 200)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
