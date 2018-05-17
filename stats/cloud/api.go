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
	"compress/gzip"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/loadimpact/k6/lib"
	"github.com/pkg/errors"
)

type ThresholdResult map[string]map[string]bool

type TestRun struct ***REMOVED***
	Name       string              `json:"name"`
	ProjectID  int64               `json:"project_id,omitempty"`
	VUsMax     int64               `json:"vus"`
	Thresholds map[string][]string `json:"thresholds"`
	// Duration of test in seconds. -1 for unknown length, 0 for continuous running.
	Duration int64 `json:"duration"`
***REMOVED***

type CreateTestRunResponse struct ***REMOVED***
	ReferenceID    string  `json:"reference_id"`
	ConfigOverride *Config `json:"config"`
***REMOVED***

type TestProgressResponse struct ***REMOVED***
	RunStatusText string  `json:"run_status_text"`
	RunStatus     int     `json:"run_status"`
	ResultStatus  int     `json:"result_status"`
	Progress      float64 `json:"progress"`
***REMOVED***

type LoginResponse struct ***REMOVED***
	Token string `json:"token"`
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

func (c *Client) PushMetric(referenceID string, noCompress bool, samples []*Sample) error ***REMOVED***
	url := fmt.Sprintf("%s/metrics/%s", c.baseURL, referenceID)

	if noCompress ***REMOVED***
		req, err := c.NewRequest("POST", url, samples)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return c.Do(req, nil)
	***REMOVED***

	var buf bytes.Buffer
	if samples != nil ***REMOVED***
		b, err := json.Marshal(&samples)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		g := gzip.NewWriter(&buf)
		if _, err = g.Write(b); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err = g.Close(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.Header.Set("Content-Encoding", "gzip")
	return c.Do(req, nil)
***REMOVED***

func (c *Client) StartCloudTestRun(name string, projectID int64, arc *lib.Archive) (string, error) ***REMOVED***
	requestUrl := fmt.Sprintf("%s/archive-upload", c.baseURL)

	var buf bytes.Buffer
	mp := multipart.NewWriter(&buf)

	if err := mp.WriteField("name", name); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if projectID != 0 ***REMOVED***
		if err := mp.WriteField("project_id", strconv.FormatInt(projectID, 10)); err != nil ***REMOVED***
			return "", err
		***REMOVED***
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

	req, err := http.NewRequest("POST", requestUrl, &buf)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	req.Header.Set("Content-Type", mp.FormDataContentType())

	ctrr := CreateTestRunResponse***REMOVED******REMOVED***
	if err := c.Do(req, &ctrr); err != nil ***REMOVED***
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
		ResultStatus int             `json:"result_status"`
		Thresholds   ThresholdResult `json:"thresholds"`
	***REMOVED******REMOVED***
		status,
		thresholds,
	***REMOVED***

	req, err := c.NewRequest("POST", url, data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return c.Do(req, nil)
***REMOVED***

func (c *Client) GetTestProgress(referenceID string) (*TestProgressResponse, error) ***REMOVED***
	url := fmt.Sprintf("%s/test-progress/%s", c.baseURL, referenceID)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ctrr := TestProgressResponse***REMOVED******REMOVED***
	err = c.Do(req, &ctrr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &ctrr, nil
***REMOVED***

func (c *Client) StopCloudTestRun(referenceID string) error ***REMOVED***
	url := fmt.Sprintf("%s/tests/%s/stop", c.baseURL, referenceID)

	req, err := c.NewRequest("POST", url, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return c.Do(req, nil)
***REMOVED***

func (c *Client) ValidateOptions(options lib.Options) error ***REMOVED***
	url := fmt.Sprintf("%s/validate-options", c.baseURL)

	data := struct ***REMOVED***
		Options lib.Options `json:"options"`
	***REMOVED******REMOVED***
		options,
	***REMOVED***

	req, err := c.NewRequest("POST", url, data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return c.Do(req, nil)
***REMOVED***

func (c *Client) Login(email string, password string) (*LoginResponse, error) ***REMOVED***
	url := fmt.Sprintf("%s/login", c.baseURL)

	data := struct ***REMOVED***
		Email    string `json:"email"`
		Password string `json:"password"`
	***REMOVED******REMOVED***
		email,
		password,
	***REMOVED***

	req, err := c.NewRequest("POST", url, data)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	lr := LoginResponse***REMOVED******REMOVED***
	err = c.Do(req, &lr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &lr, nil
***REMOVED***
