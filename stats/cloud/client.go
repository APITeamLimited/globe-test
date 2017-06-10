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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	RequestTimeout = 10 * time.Second
)

// Client handles communication with Load Impact cloud API.
type Client struct ***REMOVED***
	client  *http.Client
	token   string
	baseURL string
	version string
***REMOVED***

func NewClient(token, host, version string) *Client ***REMOVED***
	client := &http.Client***REMOVED***
		Timeout: RequestTimeout,
	***REMOVED***

	hostEnv := os.Getenv("K6CLOUD_HOST")
	if hostEnv != "" ***REMOVED***
		host = hostEnv
	***REMOVED***
	if host == "" ***REMOVED***
		host = "https://ingest.loadimpact.com"
	***REMOVED***

	baseURL := fmt.Sprintf("%s/v1", host)

	c := &Client***REMOVED***
		client:  client,
		token:   token,
		baseURL: baseURL,
		version: version,
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
	req.Header.Set("User-Agent", "k6cloud/"+c.version)

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

	if err = checkResponse(resp); err != nil ***REMOVED***
		return err
	***REMOVED***

	if v != nil ***REMOVED***
		if err = json.NewDecoder(resp.Body).Decode(v); err == io.EOF ***REMOVED***
			err = nil // Ignore EOF from empty body
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

func checkResponse(r *http.Response) error ***REMOVED***
	if c := r.StatusCode; c >= 200 && c <= 299 ***REMOVED***
		return nil
	***REMOVED***

	if r.StatusCode == 401 ***REMOVED***
		return ErrNotAuthenticated
	***REMOVED*** else if r.StatusCode == 403 ***REMOVED***
		return ErrNotAuthorized
	***REMOVED***

	// Struct of errors set back from API
	errorStruct := &struct ***REMOVED***
		ErrorData struct ***REMOVED***
			Message string `json:"message"`
			Code    int    `json:"code"`
		***REMOVED*** `json:"error"`
	***REMOVED******REMOVED******REMOVED***

	if err := json.NewDecoder(r.Body).Decode(errorStruct); err != nil ***REMOVED***
		return errors.Wrap(err, "Non-standard API error response")
	***REMOVED***

	errorResponse := &ErrorResponse***REMOVED***
		Response: r,
		Message:  errorStruct.ErrorData.Message,
		Code:     errorStruct.ErrorData.Code,
	***REMOVED***

	return errorResponse
***REMOVED***
