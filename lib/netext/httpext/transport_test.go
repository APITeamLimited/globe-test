/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package httpext

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/stats"
	"github.com/sirupsen/logrus"
)

func BenchmarkMeasureAndEmitMetrics(b *testing.B) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	samples := make(chan stats.SampleContainer, 10)
	defer close(samples)
	go func() ***REMOVED***
		for range samples ***REMOVED***
		***REMOVED***
	***REMOVED***()
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	state := &lib.State***REMOVED***
		Options: lib.Options***REMOVED***
			RunTags:    &stats.SampleTags***REMOVED******REMOVED***,
			SystemTags: &stats.DefaultSystemTagSet,
		***REMOVED***,
		Samples: samples,
		Logger:  logger,
	***REMOVED***
	t := transport***REMOVED***
		state: state,
		ctx:   ctx,
	***REMOVED***

	b.ResetTimer()
	unfRequest := &unfinishedRequest***REMOVED***
		tracer: &Tracer***REMOVED******REMOVED***,
		response: &http.Response***REMOVED***
			StatusCode: 200,
		***REMOVED***,
		request: &http.Request***REMOVED***
			URL: &url.URL***REMOVED***
				Host:   "example.com",
				Scheme: "https",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	b.Run("no responseCallback", func(b *testing.B) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			t.measureAndEmitMetrics(unfRequest)
		***REMOVED***
	***REMOVED***)

	t.responseCallback = func(n int) bool ***REMOVED*** return true ***REMOVED***

	b.Run("responseCallback", func(b *testing.B) ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			t.measureAndEmitMetrics(unfRequest)
		***REMOVED***
	***REMOVED***)
***REMOVED***
