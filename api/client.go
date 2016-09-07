package api

import (
	"encoding/json"
	"errors"
	"github.com/google/jsonapi"
	"github.com/loadimpact/speedboat/lib"
	"io"
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

func (c *Client) call(method string, relative url.URL) (io.ReadCloser, error) ***REMOVED***
	req := http.Request***REMOVED***
		Method: method,
		URL:    c.BaseURL.ResolveReference(&relative),
	***REMOVED***
	res, err := c.Client.Do(&req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if res.StatusCode >= 400 ***REMOVED***
		body, _ := ioutil.ReadAll(res.Body)

		var envelope ErrorResponse
		if err := json.Unmarshal(body, &envelope); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if len(envelope.Errors) == 0 ***REMOVED***
			return nil, errors.New("Unknown error")
		***REMOVED***
		return nil, errors.New(envelope.Errors[0].Title)
	***REMOVED***

	return res.Body, nil
***REMOVED***

func (c *Client) callSingle(method string, relative url.URL, out interface***REMOVED******REMOVED***) error ***REMOVED***
	body, err := c.call(method, relative)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer body.Close()

	return jsonapi.UnmarshalPayload(body, out)
***REMOVED***

func (c *Client) Ping() error ***REMOVED***
	body, err := c.call("GET", url.URL***REMOVED***Path: "/ping"***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	body.Close()
	return nil
***REMOVED***

// Status returns the status of the currently running test.
func (c *Client) Status() (lib.Status, error) ***REMOVED***
	var status lib.Status
	if err := c.callSingle("GET", url.URL***REMOVED***Path: "/v1/status"***REMOVED***, &status); err != nil ***REMOVED***
		return status, err
	***REMOVED***
	return status, nil
***REMOVED***

// Scales the currently running test.
func (c *Client) Scale(vus int64) error ***REMOVED***
	// u := url.URL***REMOVED***Path: "/v1/scale", RawQuery: fmt.Sprintf("vus=%d", vus)***REMOVED***
	// if err := c.call("POST", u, nil); err != nil ***REMOVED***
	// 	return err
	// ***REMOVED***
	return nil
***REMOVED***

// Aborts the currently running test.
func (c *Client) Abort() error ***REMOVED***
	// if err := c.call("POST", url.URL***REMOVED***Path: "/v1/abort"***REMOVED***, nil); err != nil ***REMOVED***
	// 	return err
	// ***REMOVED***
	return nil
***REMOVED***
