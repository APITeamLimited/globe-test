package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

const (
	TIMEOUT = 10 * time.Second
)

// Client handles communication with Load Impact cloud API.
type Client struct ***REMOVED***
	client  *http.Client
	token   string
	baseURL string
	version string
***REMOVED***

func NewClient(token, host, version string) *Client ***REMOVED***

	var client = &http.Client***REMOVED***
		Timeout: TIMEOUT,
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

	err = checkResponse(resp)

	if v != nil ***REMOVED***
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF ***REMOVED***
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
		return AuthenticateError
	***REMOVED*** else if r.StatusCode == 403 ***REMOVED***
		return AuthorizeError
	***REMOVED***

	// Struct of errors set back from API
	errorStruct := &struct ***REMOVED***
		ErrorData struct ***REMOVED***
			Message string `json:"message"`
			Code    int    `json:"code"`
		***REMOVED*** `json:"error"`
	***REMOVED******REMOVED******REMOVED***

	err := json.NewDecoder(r.Body).Decode(errorStruct)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "Non-standard API error response")
	***REMOVED***

	errorResponse := &ErrorResponse***REMOVED***
		Response: r,
		Message:  errorStruct.ErrorData.Message,
		Code:     errorStruct.ErrorData.Code,
	***REMOVED***

	return errorResponse
***REMOVED***
