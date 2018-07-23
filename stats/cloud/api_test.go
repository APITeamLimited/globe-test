/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
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
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/loadimpact/k6/lib/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() ***REMOVED***
	_ = os.Setenv("K6CLOUD_HOST", "")
	_ = os.Setenv("K6CLOUD_TOKEN", "")
***REMOVED***

func fprintf(t *testing.T, w io.Writer, format string, a ...interface***REMOVED******REMOVED***) int ***REMOVED***
	n, err := fmt.Fprintf(w, format, a...)
	require.NoError(t, err)
	return n
***REMOVED***

func TestCreateTestRun(t *testing.T) ***REMOVED***
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fprintf(t, w, `***REMOVED***"reference_id": "1", "config": ***REMOVED***"aggregationPeriod": "2s"***REMOVED******REMOVED***`)
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	tr := &TestRun***REMOVED***
		Name: "test",
	***REMOVED***
	resp, err := client.CreateTestRun(tr)

	assert.Nil(t, err)
	assert.Equal(t, resp.ReferenceID, "1")
	assert.NotNil(t, resp.ConfigOverride)
	assert.True(t, resp.ConfigOverride.AggregationPeriod.Valid)
	assert.Equal(t, types.Duration(2*time.Second), resp.ConfigOverride.AggregationPeriod.Duration)
	assert.False(t, resp.ConfigOverride.AggregationMinSamples.Valid)
***REMOVED***

func TestPublishMetric(t *testing.T) ***REMOVED***
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fprintf(t, w, "")
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	samples := []*Sample***REMOVED***
		***REMOVED***
			Type:   "Point",
			Metric: "metric",
			Data: &SampleDataSingle***REMOVED***
				Type:  1,
				Time:  Timestamp(time.Now()),
				Value: 1.2,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := client.PushMetric("1", false, samples)

	assert.Nil(t, err)
***REMOVED***

func TestFinished(t *testing.T) ***REMOVED***
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fprintf(t, w, "")
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	thresholds := map[string]map[string]bool***REMOVED***
		"threshold": ***REMOVED***
			"max < 10": true,
		***REMOVED***,
	***REMOVED***
	err := client.TestFinished("1", thresholds, true, 0)

	assert.Nil(t, err)
***REMOVED***

func TestAuthorizedError(t *testing.T) ***REMOVED***
	called := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		called++
		w.WriteHeader(http.StatusForbidden)
		fprintf(t, w, `***REMOVED***"error": ***REMOVED***"code": 5, "message": "Not allowed"***REMOVED******REMOVED***`)
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	resp, err := client.CreateTestRun(&TestRun***REMOVED***Name: "test"***REMOVED***)

	assert.Equal(t, 1, called)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "403 Not allowed [err code 5]")
***REMOVED***

func TestRetry(t *testing.T) ***REMOVED***
	called := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		called++
		w.WriteHeader(500)
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")
	client.retryInterval = 1 * time.Millisecond
	resp, err := client.CreateTestRun(&TestRun***REMOVED***Name: "test"***REMOVED***)

	assert.Equal(t, 3, called)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
***REMOVED***

func TestRetrySuccessOnSecond(t *testing.T) ***REMOVED***
	called := 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		called++
		if called == 2 ***REMOVED***
			fprintf(t, w, `***REMOVED***"reference_id": "1"***REMOVED***`)
			return
		***REMOVED***
		w.WriteHeader(500)
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")
	client.retryInterval = 1 * time.Millisecond
	resp, err := client.CreateTestRun(&TestRun***REMOVED***Name: "test"***REMOVED***)

	assert.Equal(t, 2, called)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
***REMOVED***
