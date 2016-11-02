package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
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
	res.Body.Close()

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
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Status returns the status of the currently running test.
func (c *Client) Status() (lib.Status, error) ***REMOVED***
	var status lib.Status
	if err := c.call("GET", "/v1/status", nil, &status); err != nil ***REMOVED***
		return status, err
	***REMOVED***
	return status, nil
***REMOVED***

// Updates the status of the currently running test.
func (c *Client) UpdateStatus(status lib.Status) (lib.Status, error) ***REMOVED***
	data, err := jsonapi.Marshal(status)
	if err != nil ***REMOVED***
		return status, err
	***REMOVED***
	if err := c.call("PATCH", "/v1/status", data, &status); err != nil ***REMOVED***
		return status, err
	***REMOVED***
	return status, nil
***REMOVED***

// Returns a snapshot of metrics for the currently running test.
func (c *Client) Metrics() ([]stats.Metric, error) ***REMOVED***
	var metrics []stats.Metric
	if err := c.call("GET", "/v1/metrics", nil, &metrics); err != nil ***REMOVED***
		return metrics, err
	***REMOVED***
	return metrics, nil
***REMOVED***

func (c *Client) Metric(name string) (stats.Metric, error) ***REMOVED***
	var metric stats.Metric
	if err := c.call("GET", "/v1/metrics/"+name, nil, &metric); err != nil ***REMOVED***
		return metric, err
	***REMOVED***
	return metric, nil
***REMOVED***
