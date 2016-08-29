package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loadimpact/speedboat/lib"
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

func New(addr string) (*Client, error) ***REMOVED***
	if addr == "" ***REMOVED***
		return nil, errNoAddress
	***REMOVED***

	return &Client***REMOVED***
		BaseURL: url.URL***REMOVED***Scheme: "http", Host: addr***REMOVED***,
		Client:  http.Client***REMOVED******REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (c *Client) call(method string, relative url.URL, out interface***REMOVED******REMOVED***) error ***REMOVED***
	req := http.Request***REMOVED***
		Method: method,
		URL:    c.BaseURL.ResolveReference(&relative),
	***REMOVED***
	res, err := c.Client.Do(&req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer res.Body.Close()

	if res.StatusCode >= 400 ***REMOVED***
		var envelope struct ***REMOVED***
			Error string `json:"error"`
		***REMOVED***
		body, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, &envelope); err != nil ***REMOVED***
			return err
		***REMOVED***
		if envelope.Error == "" ***REMOVED***
			envelope.Error = res.Status
		***REMOVED***
		return errors.New(envelope.Error)
	***REMOVED***

	if out == nil ***REMOVED***
		return nil
	***REMOVED***

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, out); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (c *Client) Ping() error ***REMOVED***
	if err := c.call("GET", url.URL***REMOVED***Path: "/ping"***REMOVED***, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Status returns the status of the currently running test.
func (c *Client) Status() (lib.Status, error) ***REMOVED***
	var status lib.Status
	if err := c.call("GET", url.URL***REMOVED***Path: "/v1/status"***REMOVED***, &status); err != nil ***REMOVED***
		return status, err
	***REMOVED***
	return status, nil
***REMOVED***

// Scales the currently running test.
func (c *Client) Scale(vus int64) error ***REMOVED***
	u := url.URL***REMOVED***Path: "/v1/scale", RawQuery: fmt.Sprintf("vus=%d", vus)***REMOVED***
	if err := c.call("POST", u, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Aborts the currently running test.
func (c *Client) Abort() error ***REMOVED***
	if err := c.call("POST", url.URL***REMOVED***Path: "/v1/abort"***REMOVED***, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
