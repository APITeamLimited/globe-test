/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/sirupsen/logrus"

	"github.com/loadimpact/k6/api/v1"
)

// Client is a simple HTTP client for the REST API.
type Client struct ***REMOVED***
	BaseURL    *url.URL
	httpClient *http.Client
	logger     *logrus.Entry
***REMOVED***

// Option function are helpers that enable the flexible configuration of the
// REST API client.
type Option func(*Client)

// New returns a newly configured REST API Client.
func New(base string, options ...Option) (*Client, error) ***REMOVED***
	baseURL, err := url.Parse("http://" + base)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c := &Client***REMOVED***
		BaseURL:    baseURL,
		httpClient: http.DefaultClient,
	***REMOVED***

	for _, option := range options ***REMOVED***
		option(c)
	***REMOVED***

	return c, nil
***REMOVED***

// WithHTTPClient configures the supplied HTTP client to be used when making
// REST API requests.
func WithHTTPClient(httpClient *http.Client) Option ***REMOVED***
	return Option(func(c *Client) ***REMOVED***
		c.httpClient = httpClient
	***REMOVED***)
***REMOVED***

// WithLogger sets the specifield logger to the client.
func WithLogger(logger *logrus.Entry) Option ***REMOVED***
	return Option(func(c *Client) ***REMOVED***
		c.logger = logger
	***REMOVED***)
***REMOVED***

// Call executes the desired REST API request.
func (c *Client) Call(ctx context.Context, method string, rel *url.URL, body, out interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	if c.logger != nil ***REMOVED***
		c.logger.Debugf("[REST API] Making a %s request to '%s'", method, rel.String())
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				c.logger.WithError(err).Error("[REST API] Error")
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	var bodyReader io.ReadCloser
	if body != nil ***REMOVED***
		var bodyData []byte
		switch val := body.(type) ***REMOVED***
		case []byte:
			bodyData = val
		case string:
			bodyData = []byte(val)
		default:
			bodyData, err = jsonapi.Marshal(body)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		bodyReader = ioutil.NopCloser(bytes.NewBuffer(bodyData))
	***REMOVED***

	req := &http.Request***REMOVED***
		Method: method,
		URL:    c.BaseURL.ResolveReference(rel),
		Body:   bodyReader,
	***REMOVED***
	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED*** _ = res.Body.Close() ***REMOVED***()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if res.StatusCode >= 400 ***REMOVED***
		var errs v1.ErrorResponse
		if err := json.Unmarshal(data, &errs); err != nil ***REMOVED***
			return err
		***REMOVED***
		return errs.Errors[0]
	***REMOVED***

	if out != nil ***REMOVED***
		return jsonapi.Unmarshal(data, out)
	***REMOVED***
	return nil
***REMOVED***
