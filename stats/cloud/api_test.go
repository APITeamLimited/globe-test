package cloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

	if err != nil ***REMOVED***
		t.Fatal(err, nil)
	***REMOVED***
	if resp.ReferenceID != "1" ***REMOVED***
		t.Fatal(resp.ReferenceID, "1")
	***REMOVED***
***REMOVED***

func TestPublishMetric(t *testing.T) ***REMOVED***
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintf(w, "")
	***REMOVED***))
	defer server.Close()

	client := NewClient("token", server.URL, "1.0")

	samples := []*sample***REMOVED***
		***REMOVED***
			Type:   "Point",
			Metric: "metric",
			Data: sampleData***REMOVED***
				Type:  1,
				Time:  time.Now(),
				Value: 1.2,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := client.PushMetric("1", samples)

	if err != nil ***REMOVED***
		t.Fatal(err, nil)
	***REMOVED***
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

	if err != nil ***REMOVED***
		t.Fatal(err, nil)
	***REMOVED***
***REMOVED***
