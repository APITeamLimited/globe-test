package cloud

import (
	"fmt"
	"time"

	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
)

type sample struct ***REMOVED***
	Type   string     `json:"type"`
	Metric string     `json:"metric"`
	Data   sampleData `json:"data"`
***REMOVED***

type sampleData struct ***REMOVED***
	Type  stats.MetricType  `json:"type"`
	Time  time.Time         `json:"time"`
	Value float64           `json:"value"`
	Tags  map[string]string `json:"tags,omitempty"`
***REMOVED***

type ThresholdResult map[string]map[string]bool

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

func (c *Client) CreateTestRun(testRun *TestRun) (*CreateTestRunResponse, error) ***REMOVED***
	url := fmt.Sprintf("%s/tests", c.baseURL)
	req, err := c.NewRequest("POST", url, testRun)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var ctrr = CreateTestRunResponse***REMOVED******REMOVED***
	err = c.Do(req, &ctrr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if ctrr.ReferenceID == "" ***REMOVED***
		return nil, errors.Errorf("Failed to get a reference ID")
	***REMOVED***

	return &ctrr, nil
***REMOVED***

func (c *Client) PushMetric(referenceID string, samples []*sample) error ***REMOVED***
	url := fmt.Sprintf("%s/metrics/%s", c.baseURL, referenceID)

	req, err := c.NewRequest("POST", url, samples)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = c.Do(req, nil)
	return err
***REMOVED***

func (c *Client) TestFinished(referenceID string, thresholds ThresholdResult, tained bool) error ***REMOVED***
	url := fmt.Sprintf("%s/tests/%s", c.baseURL, referenceID)

	status := 1

	if tained ***REMOVED***
		status = 2
	***REMOVED***

	data := struct ***REMOVED***
		Status     int             `json:"status"`
		Thresholds ThresholdResult `json:"thresholds"`
	***REMOVED******REMOVED***
		status,
		thresholds,
	***REMOVED***

	req, err := c.NewRequest("POST", url, data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = c.Do(req, nil)
	return err
***REMOVED***
