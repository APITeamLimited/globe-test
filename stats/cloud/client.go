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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	// RequestTimeout is the default cloud request timeout
	RequestTimeout = 20 * time.Second
	// RetryInterval is the default cloud request retry interval
	RetryInterval = 500 * time.Millisecond
	// MaxRetries specifies max retry attempts
	MaxRetries = 3

	k6IdempotencyKeyHeader = "k6-Idempotency-Key"
)

// Client handles communication with Load Impact cloud API.
type Client struct ***REMOVED***
	client  *http.Client
	token   string
	baseURL string
	version string

	retries       int
	retryInterval time.Duration
***REMOVED***

func NewClient(token, host, version string) *Client ***REMOVED***
	c := &Client***REMOVED***
		client:        &http.Client***REMOVED***Timeout: RequestTimeout***REMOVED***,
		token:         token,
		baseURL:       fmt.Sprintf("%s/v1", host),
		version:       version,
		retries:       MaxRetries,
		retryInterval: RetryInterval,
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

	req, err := http.NewRequest(method, url, buf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if shouldAddIdempotencyKey(req) ***REMOVED***
		req.Header.Set(k6IdempotencyKeyHeader, randomStrHex())
	***REMOVED***

	return req, nil
***REMOVED***

func (c *Client) Do(req *http.Request, v interface***REMOVED******REMOVED***) error ***REMOVED***
	var originalBody []byte
	var err error

	if req.Body != nil ***REMOVED***
		originalBody, err = ioutil.ReadAll(req.Body)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if cerr := req.Body.Close(); cerr != nil ***REMOVED***
			err = cerr
		***REMOVED***
	***REMOVED***

	for i := 1; i <= c.retries; i++ ***REMOVED***
		if len(originalBody) > 0 ***REMOVED***
			req.Body = ioutil.NopCloser(bytes.NewBuffer(originalBody))
		***REMOVED***

		retry, err := c.do(req, v, i)

		if retry ***REMOVED***
			time.Sleep(c.retryInterval)
			continue
		***REMOVED***

		return err
	***REMOVED***

	return err
***REMOVED***

func (c *Client) do(req *http.Request, v interface***REMOVED******REMOVED***, attempt int) (retry bool, err error) ***REMOVED***
	if req.Header.Get("Content-Type") == "" ***REMOVED***
		req.Header.Set("Content-Type", "application/json")
	***REMOVED***
	if c.token != "" ***REMOVED***
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.token))
	***REMOVED***
	req.Header.Set("User-Agent", "k6cloud/"+c.version)
	resp, err := c.client.Do(req)

	defer func() ***REMOVED***
		if resp != nil ***REMOVED***
			if cerr := resp.Body.Close(); cerr != nil && err == nil ***REMOVED***
				err = cerr
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if shouldRetry(resp, err, attempt, c.retries) ***REMOVED***
		return true, err
	***REMOVED***

	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if err = checkResponse(resp); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if v != nil ***REMOVED***
		if err = json.NewDecoder(resp.Body).Decode(v); err == io.EOF ***REMOVED***
			err = nil // Ignore EOF from empty body
		***REMOVED***
	***REMOVED***

	return false, err
***REMOVED***

func checkResponse(r *http.Response) error ***REMOVED***
	if r == nil ***REMOVED***
		return ErrUnknown
	***REMOVED***

	if c := r.StatusCode; c >= 200 && c <= 299 ***REMOVED***
		return nil
	***REMOVED***

	data, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var payload struct ***REMOVED***
		Error ErrorResponse `json:"error"`
	***REMOVED***
	if err := json.Unmarshal(data, &payload); err != nil ***REMOVED***
		if r.StatusCode == http.StatusUnauthorized ***REMOVED***
			return ErrNotAuthenticated
		***REMOVED***
		if r.StatusCode == http.StatusForbidden ***REMOVED***
			return ErrNotAuthorized
		***REMOVED***
		return errors.Errorf(
			"Unexpected HTTP error from %s: %d %s",
			r.Request.URL,
			r.StatusCode,
			http.StatusText(r.StatusCode),
		)
	***REMOVED***
	payload.Error.Response = r
	return payload.Error
***REMOVED***

func shouldRetry(resp *http.Response, err error, attempt, maxAttempts int) bool ***REMOVED***
	if attempt >= maxAttempts ***REMOVED***
		return false
	***REMOVED***

	if resp == nil || err != nil ***REMOVED***
		return true
	***REMOVED***

	if resp.StatusCode >= 500 || resp.StatusCode == 429 ***REMOVED***
		return true
	***REMOVED***

	return false
***REMOVED***

func shouldAddIdempotencyKey(req *http.Request) bool ***REMOVED***
	switch req.Method ***REMOVED***
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return false
	default:
		return req.Header.Get(k6IdempotencyKeyHeader) == ""
	***REMOVED***
***REMOVED***

// randomStrHex returns a hex string which can be used
// for session token id or idempotency key.
//nolint:gosec
func randomStrHex() string ***REMOVED***
	// 16 hex characters
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
***REMOVED***

func init() ***REMOVED***
	rand.Seed(time.Now().UTC().UnixNano())
***REMOVED***
