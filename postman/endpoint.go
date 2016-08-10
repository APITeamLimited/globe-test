package postman

import (
	"bytes"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

var (
	ErrItemHasNoRequest = errors.New("can't make an endpoint out of an item with no request")
)

type Endpoint struct ***REMOVED***
	Method string
	URL    *url.URL
	Header http.Header
	Body   []byte
***REMOVED***

func MakeEndpoint(i Item) (Endpoint, error) ***REMOVED***
	if i.Request.URL == "" ***REMOVED***
		return Endpoint***REMOVED******REMOVED***, ErrItemHasNoRequest
	***REMOVED***

	u, err := url.Parse(i.Request.URL)
	if err != nil ***REMOVED***
		return Endpoint***REMOVED******REMOVED***, err
	***REMOVED***

	header := make(http.Header)
	for _, item := range i.Request.Header ***REMOVED***
		header[item.Key] = append(header[item.Key], item.Value)
	***REMOVED***

	var body []byte
	switch i.Request.Body.Mode ***REMOVED***
	case "raw":
		body = []byte(i.Request.Body.Raw)
	case "urlencoded":
		values := make(url.Values)
		for _, field := range i.Request.Body.URLEncoded ***REMOVED***
			if !field.Enabled ***REMOVED***
				continue
			***REMOVED***
			values[field.Key] = append(values[field.Key], field.Value)
		***REMOVED***
		body = []byte(values.Encode())
	case "formdata":
		body = make([]byte, 0)
		w := multipart.NewWriter(bytes.NewBuffer(body))
		for _, field := range i.Request.Body.FormData ***REMOVED***
			if !field.Enabled ***REMOVED***
				continue
			***REMOVED***

			if err := w.WriteField(field.Key, field.Value); err != nil ***REMOVED***
				return Endpoint***REMOVED******REMOVED***, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return Endpoint***REMOVED***i.Request.Method, u, header, body***REMOVED***, nil
***REMOVED***

func (ep Endpoint) Request() http.Request ***REMOVED***
	return http.Request***REMOVED***
		Method: ep.Method,
		URL:    ep.URL,
		Header: ep.Header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(ep.Body)),
	***REMOVED***
***REMOVED***
