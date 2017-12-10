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
        "bytes"
	"fmt"
        "mime/multipart"
        "net/http"
	"time"

        "github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/pkg/errors"
)

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

	ctrr := CreateTestRunResponse***REMOVED******REMOVED***
	err = c.Do(req, &ctrr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if ctrr.ReferenceID == "" ***REMOVED***
		return nil, errors.Errorf("Failed to get a reference ID")
	***REMOVED***

	return &ctrr, nil
***REMOVED***

func (c *Client) PushMetric(referenceID string, samples []*Sample) error ***REMOVED***
        url := fmt.Sprintf("%s/metrics/%s", c.baseURL, referenceID)

        req, err := c.NewRequest("POST", url, samples)
        if err != nil ***REMOVED***
                return err
        ***REMOVED***

        err = c.Do(req, nil)
        return err
***REMOVED***

func (c *Client) StartCloudTestRun(name string, arc *lib.Archive) (string, error) ***REMOVED***
        requestUrl := fmt.Sprintf("%s/archive-upload", c.baseURL)

        buf := bytes.NewBuffer(nil)
        mp := multipart.NewWriter(buf)

        if err := mp.WriteField("name", name); err != nil ***REMOVED***
                return "", err
        ***REMOVED***

        fw, err := mp.CreateFormFile("file", "archive.tar")
        if err != nil ***REMOVED***
                return "", err
        ***REMOVED***

        if err := arc.Write(fw); err != nil ***REMOVED***
                return "", err
        ***REMOVED***

        if err := mp.Close(); err != nil ***REMOVED***
                return "", err
        ***REMOVED***

        req, err := http.NewRequest("POST", requestUrl, buf)
        if err != nil ***REMOVED***
                return "", err
        ***REMOVED***

        req.Header.Set("Content-Type", mp.FormDataContentType())

        ctrr := CreateTestRunResponse***REMOVED******REMOVED***
        err = c.Do(req, &ctrr)
	if err != nil ***REMOVED***
                return "", err
	***REMOVED***

        return ctrr.ReferenceID, nil
***REMOVED***

func (c *Client) TestFinished(referenceID string, thresholds ThresholdResult, tained bool) error ***REMOVED***
	url := fmt.Sprintf("%s/tests/%s", c.baseURL, referenceID)

	status := 0

	if tained ***REMOVED***
		status = 1
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

func (c *Client) ValidateConfig(arc *lib.Archive) error ***REMOVED***
        url := fmt.Sprintf("%s/validate-config", c.baseURL)

        data := struct ***REMOVED***
                Config *lib.Archive `json:"config"`
        ***REMOVED******REMOVED***
                arc,
        ***REMOVED***

        req, err := c.NewRequest("POST", url, data)
        if err != nil ***REMOVED***
                return err
        ***REMOVED***

        err = c.Do(req, nil)
        return err

***REMOVED***
