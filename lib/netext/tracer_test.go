/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package netext

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
)

func TestTracer(t *testing.T) ***REMOVED***
	client := &http.Client***REMOVED***
		Transport: &http.Transport***REMOVED***DialContext: NewDialer(net.Dialer***REMOVED******REMOVED***).DialContext***REMOVED***,
	***REMOVED***

	for _, isReuse := range []bool***REMOVED***false, true***REMOVED*** ***REMOVED***
		name := "First"
		if isReuse ***REMOVED***
			name = "Reuse"
		***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			tracer := &Tracer***REMOVED******REMOVED***
			req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			res, err := client.Do(req.WithContext(WithTracer(context.Background(), tracer)))
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			_, err = io.Copy(ioutil.Discard, res.Body)
			assert.NoError(t, err)
			assert.NoError(t, res.Body.Close())
			samples := tracer.Done().Samples(map[string]string***REMOVED***"tag": "value"***REMOVED***)

			assert.Len(t, samples, 8)
			seenMetrics := map[*stats.Metric]bool***REMOVED******REMOVED***
			for _, s := range samples ***REMOVED***
				assert.NotContains(t, seenMetrics, s.Metric)
				seenMetrics[s.Metric] = true

				assert.False(t, s.Time.IsZero())
				assert.Equal(t, map[string]string***REMOVED***"tag": "value"***REMOVED***, s.Tags)

				switch s.Metric ***REMOVED***
				case metrics.HTTPReqs:
					assert.Equal(t, 1.0, s.Value)
				case metrics.HTTPReqConnecting:
					if isReuse ***REMOVED***
						assert.Equal(t, 0.0, s.Value)
						break
					***REMOVED***
					fallthrough
				case metrics.HTTPReqDuration, metrics.HTTPReqBlocked, metrics.HTTPReqSending, metrics.HTTPReqWaiting, metrics.HTTPReqReceiving:
					assert.True(t, s.Value > 0.0, "%s is <= 0", s.Metric.Name)
				case metrics.HTTPReqTLSHandshaking:
					if !isReuse ***REMOVED***
						assert.True(t, s.Value > 0.0, "%s is <= 0", s.Metric.Name)
						continue
					***REMOVED***
					assert.True(t, s.Value == 0.0, "%s is <> 0", s.Metric.Name)
				default:
					t.Errorf("unexpected metric: %s", s.Metric.Name)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
