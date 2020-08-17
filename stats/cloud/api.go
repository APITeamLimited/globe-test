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
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/mailru/easyjson"
	"github.com/pkg/errors"

	"github.com/loadimpact/k6/lib"
	"github.com/sirupsen/logrus"
)

type ResultStatus int

const (
	ResultStatusPassed ResultStatus = 0
	ResultStatusFailed ResultStatus = 1
)
const expectedGzipRatio = 6 // based on test it is around 6.8, but we don't need to be that accurate

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
	RunStatusText string        `json:"run_status_text"`
	RunStatus     lib.RunStatus `json:"run_status"`
	ResultStatus  ResultStatus  `json:"result_status"`
	Progress      float64       `json:"progress"`
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

// PushMetric pushes the provided metric samples for the given referenceID
func (c *Client) PushMetric(referenceID string, noCompress bool, s []*Sample) error ***REMOVED***
	start := time.Now()
	url := fmt.Sprintf("%s/metrics/%s", c.baseURL, referenceID)

	jsonStart := time.Now()
	b, err := easyjson.Marshal(samples(s))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	jsonTime := time.Since(jsonStart)

	// TODO: change the context, maybe to one with a timeout
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	req.Header.Set("X-Payload-Sample-Count", strconv.Itoa(len(s)))
	var additionalFields logrus.Fields

	if !noCompress ***REMOVED***
		buf := c.pushBufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer c.pushBufferPool.Put(buf)
		unzippedSize := len(b)
		buf.Grow(unzippedSize / expectedGzipRatio)
		gzipStart := time.Now()
		***REMOVED***
			g, _ := gzip.NewWriterLevel(buf, gzip.BestSpeed)
			if _, err = g.Write(b); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err = g.Close(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		gzipTime := time.Since(gzipStart)

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("X-Payload-Byte-Count", strconv.Itoa(unzippedSize))

		additionalFields = logrus.Fields***REMOVED***
			"unzipped_size":  unzippedSize,
			"gzip_t":         gzipTime,
			"content_length": buf.Len(),
		***REMOVED***

		b = buf.Bytes()
	***REMOVED***

	req.Header.Set("Content-Length", strconv.Itoa(len(b)))
	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	req.GetBody = func() (io.ReadCloser, error) ***REMOVED***
		return ioutil.NopCloser(bytes.NewReader(b)), nil
	***REMOVED***

	err = c.Do(req, nil)

	c.logger.WithFields(logrus.Fields***REMOVED***
		"t":         time.Since(start),
		"json_t":    jsonTime,
		"part_size": len(s),
	***REMOVED***).WithFields(additionalFields).Debug("Pushed part to cloud")

	return err
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

func (c *Client) TestFinished(referenceID string, thresholds ThresholdResult, tained bool, runStatus lib.RunStatus) error ***REMOVED***
	url := fmt.Sprintf("%s/tests/%s", c.baseURL, referenceID)

	resultStatus := ResultStatusPassed
	if tained ***REMOVED***
		resultStatus = ResultStatusFailed
	***REMOVED***

	data := struct ***REMOVED***
		ResultStatus ResultStatus    `json:"result_status"`
		RunStatus    lib.RunStatus   `json:"run_status"`
		Thresholds   ThresholdResult `json:"thresholds"`
	***REMOVED******REMOVED***
		resultStatus,
		runStatus,
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
