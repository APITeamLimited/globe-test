package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/stats"
)

// Client handles communication with Load Impact cloud API.
type Client struct ***REMOVED***
	client  *http.Client
	token   string
	baseURL string
***REMOVED***

type ErrorResponse struct ***REMOVED***
	Message string `json:"message"`
***REMOVED***

func NewClient(token string) *Client ***REMOVED***

	var client = &http.Client***REMOVED***
		Timeout: 30 * time.Second,
	***REMOVED***

	host := os.Getenv("K6CLOUD_HOST")
	if host == "" ***REMOVED***
		host = "http://localhost:5000"
	***REMOVED***

	baseURL := fmt.Sprintf("%s/v1", host)

	c := &Client***REMOVED***
		client:  client,
		token:   token,
		baseURL: baseURL,
	***REMOVED***
	return c
***REMOVED***

func (c *Client) NewRequest(method, url string, data interface***REMOVED******REMOVED***) (*http.Request, error) ***REMOVED***
	var buf io.Reader

	if data != nil ***REMOVED***
		b, err := json.Marshal(&data)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		buf = bytes.NewBuffer(b)
	***REMOVED***

	return http.NewRequest(method, url, buf)
***REMOVED***

func (c *Client) Do(req *http.Request, v interface***REMOVED******REMOVED***) error ***REMOVED***

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.token))

	resp, err := c.client.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		err := resp.Body.Close()
		if err != nil ***REMOVED***
			log.Errorln(err)
		***REMOVED***
	***REMOVED***()

	if v != nil ***REMOVED***
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF ***REMOVED***
			err = nil // Ignore EOF from empty body
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

type TestRun struct ***REMOVED***
	Name       string              `json:"name"`
	ProjectID  int                 `json:"project_id,omitempty"`
	Thresholds map[string][]string `json:"thresholds"`
	// Duration of test in seconds. -1 for unknown length, 0 for continuous running.
	Duration int64 `json:"duration"`
***REMOVED***

type CreateTestRunResponse struct ***REMOVED***
	ReferenceID string `json:"reference_id"`
***REMOVED***

func (c *Client) CreateTestRun(testRun *TestRun) *CreateTestRunResponse ***REMOVED***
	url := fmt.Sprintf("%s/tests", c.baseURL)
	req, err := c.NewRequest("POST", url, testRun)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	var ctrr = CreateTestRunResponse***REMOVED******REMOVED***
	err = c.Do(req, &ctrr)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***

	return &ctrr
***REMOVED***

func (c *Client) PushMetric(referenceID string, samples []*Sample) ***REMOVED***
	url := fmt.Sprintf("%s/metrics/%s", c.baseURL, referenceID)

	req, err := c.NewRequest("POST", url, samples)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = c.Do(req, nil)
	if err != nil ***REMOVED***
		return
	***REMOVED***
***REMOVED***

func (c *Client) TestFinished(referenceID string) ***REMOVED***
	url := fmt.Sprintf("%s/tests/%s", c.baseURL, referenceID)

	data := struct ***REMOVED***
		Status int `json:"status"`
	***REMOVED******REMOVED***
		1,
	***REMOVED***

	req, err := c.NewRequest("POST", url, data)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = c.Do(req, nil)
	if err != nil ***REMOVED***
		return
	***REMOVED***
***REMOVED***

type Sample struct ***REMOVED***
	Type   string     `json:"type"`
	Metric string     `json:"metric"`
	Data   SampleData `json:"data"`
***REMOVED***

type SampleData struct ***REMOVED***
	Type  stats.MetricType  `json:"type"`
	Time  time.Time         `json:"time"`
	Value float64           `json:"value"`
	Tags  map[string]string `json:"tags,omitempty"`
***REMOVED***
