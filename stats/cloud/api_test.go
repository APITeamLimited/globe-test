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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func init() ***REMOVED***
	_ = os.Setenv("K6CLOUD_HOST", "")
	_ = os.Setenv("K6CLOUD_TOKEN", "")
***REMOVED***

func TestTimestampMarshaling(t *testing.T) ***REMOVED***
	oldTimeFormat, err := time.Parse(
		time.RFC3339,
		//1521806137415652223 as a unix nanosecond timestamp
		"2018-03-23T13:55:37.415652223+02:00",
	)
	require.NoError(t, err)

	testCases := []struct ***REMOVED***
		t   time.Time
		exp string
	***REMOVED******REMOVED***
		***REMOVED***oldTimeFormat, `"1521806137415652"`***REMOVED***,
		***REMOVED***time.Unix(1521806137, 415652223), `"1521806137415652"`***REMOVED***,
		***REMOVED***time.Unix(1521806137, 0), `"1521806137000000"`***REMOVED***,
		***REMOVED***time.Unix(0, 0), `"0"`***REMOVED***,
		***REMOVED***time.Unix(0, 1), `"0"`***REMOVED***,
		***REMOVED***time.Unix(0, 1000), `"1"`***REMOVED***,
		***REMOVED***time.Unix(1, 0), `"1000000"`***REMOVED***,
	***REMOVED***

	for i, tc := range testCases ***REMOVED***
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) ***REMOVED***
			res, err := json.Marshal(Timestamp(tc.t))
			require.NoError(t, err)
			assert.Equal(t, string(res), tc.exp)

			var rev Timestamp
			require.NoError(t, json.Unmarshal(res, &rev))
			diff := tc.t.Sub(time.Time(rev))
			if diff < -time.Microsecond || diff > time.Microsecond ***REMOVED***
				t.Errorf(
					"Expected the difference to be under a microsecond, but is %s (%d and %d)",
					diff,
					tc.t.UnixNano(),
					time.Time(rev).UnixNano(),
				)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestCreateTestRun(t *testing.T) ***REMOVED***
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintf(w, `***REMOVED***"reference_id": "1"***REMOVED***`)
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	tr := &TestRun***REMOVED***
		Name: "test",
	***REMOVED***
	resp, err := client.CreateTestRun(tr)

	assert.Nil(t, err)
	assert.Equal(t, resp.ReferenceID, "1")
***REMOVED***

func TestPublishMetric(t *testing.T) ***REMOVED***
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintf(w, "")
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	samples := []*Sample***REMOVED***
		***REMOVED***
			Type:   "Point",
			Metric: "metric",
			Data: SampleDataSingle***REMOVED***
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
		fmt.Fprintf(w, "")
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	thresholds := map[string]map[string]bool***REMOVED***
		"threshold": ***REMOVED***
			"max < 10": true,
		***REMOVED***,
	***REMOVED***
	err := client.TestFinished("1", thresholds, true)

	assert.Nil(t, err)
***REMOVED***

func TestAuthorizedError(t *testing.T) ***REMOVED***
	called := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		called++
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `***REMOVED***"error": ***REMOVED***"code": 5, "message": "Not allowed"***REMOVED******REMOVED***`)
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	resp, err := client.CreateTestRun(&TestRun***REMOVED***Name: "test"***REMOVED***)

	assert.Equal(t, 1, called)
	assert.Nil(t, resp)
	assert.EqualError(t, err, ErrNotAuthorized.Error())
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
			fmt.Fprintf(w, `***REMOVED***"reference_id": "1"***REMOVED***`)
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
