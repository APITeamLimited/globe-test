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

package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	errNoAddress = errors.New("no address given")
)

type Client struct ***REMOVED***
	BaseURL url.URL
	Client  http.Client
***REMOVED***

func NewClient(addr string) (*Client, error) ***REMOVED***
	if addr == "" ***REMOVED***
		return nil, errNoAddress
	***REMOVED***

	return &Client***REMOVED***
		BaseURL: url.URL***REMOVED***Scheme: "http", Host: addr***REMOVED***,
		Client:  http.Client***REMOVED******REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (c *Client) request(method, path string, body []byte) ([]byte, error) ***REMOVED***
	relative := url.URL***REMOVED***Path: path***REMOVED***
	req := http.Request***REMOVED***
		Method: method,
		URL:    c.BaseURL.ResolveReference(&relative),
	***REMOVED***
	if body != nil ***REMOVED***
		req.ContentLength = int64(len(body))
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	***REMOVED***

	res, err := c.Client.Do(&req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	data, _ := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()

	if res.StatusCode >= 400 ***REMOVED***
		var envelope api2go.HTTPError
		if err := json.Unmarshal(data, &envelope); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if len(envelope.Errors) == 0 ***REMOVED***
			return nil, errors.New("Unknown error")
		***REMOVED***
		return nil, errors.New(envelope.Errors[0].Title)
	***REMOVED***

	return data, nil
***REMOVED***

func (c *Client) call(method, path string, body []byte, out interface***REMOVED******REMOVED***) error ***REMOVED***
	body, err := c.request(method, path, body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return jsonapi.Unmarshal(body, out)
***REMOVED***

func (c *Client) Ping() error ***REMOVED***
	_, err := c.request("GET", "/ping", nil)
	return err
***REMOVED***

// Status returns the status of the currently running test.
func (c *Client) Status() (lib.Status, error) ***REMOVED***
	var status lib.Status
	err := c.call("GET", "/v1/status", nil, &status)
	return status, err
***REMOVED***

// Updates the status of the currently running test.
func (c *Client) UpdateStatus(status lib.Status) (lib.Status, error) ***REMOVED***
	data, err := jsonapi.Marshal(status)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	err = c.call("PATCH", "/v1/status", data, &status)
	return status, err
***REMOVED***

// Returns a snapshot of metrics for the currently running test.
func (c *Client) Metrics() ([]stats.Metric, error) ***REMOVED***
	var metrics []stats.Metric
	err := c.call("GET", "/v1/metrics", nil, &metrics)
	return metrics, err
***REMOVED***

func (c *Client) Metric(name string) (stats.Metric, error) ***REMOVED***
	var metric stats.Metric
	err := c.call("GET", "/v1/metrics/"+name, nil, &metric)
	return metric, err
***REMOVED***
