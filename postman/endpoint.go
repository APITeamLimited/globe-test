package postman

import (
	"bytes"
	"errors"
	"github.com/robertkrimen/otto"
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

	Tests      []*otto.Script
	PreRequest []*otto.Script

	URLString string
	BodyMap   map[string]string
	HeaderMap map[string]string
***REMOVED***

func MakeEndpoints(c Collection, vm *otto.Otto) ([]Endpoint, error) ***REMOVED***
	eps := make([]Endpoint, 0)
	for _, item := range c.Item ***REMOVED***
		if err := makeEndpointsFrom(item, vm, &eps); err != nil ***REMOVED***
			return eps, err
		***REMOVED***
	***REMOVED***

	return eps, nil
***REMOVED***

func makeEndpointsFrom(i Item, vm *otto.Otto, eps *[]Endpoint) error ***REMOVED***
	if i.Request.URL != "" ***REMOVED***
		ep, err := MakeEndpoint(i, vm)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*eps = append(*eps, ep)
	***REMOVED***

	for _, item := range i.Item ***REMOVED***
		if err := makeEndpointsFrom(item, vm, eps); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func MakeEndpoint(i Item, vm *otto.Otto) (Endpoint, error) ***REMOVED***
	if i.Request.URL == "" ***REMOVED***
		return Endpoint***REMOVED******REMOVED***, ErrItemHasNoRequest
	***REMOVED***

	endpoint := Endpoint***REMOVED***
		Method:    i.Request.Method,
		URLString: i.Request.URL,
	***REMOVED***

	u, err := url.Parse(i.Request.URL)
	if err != nil ***REMOVED***
		return endpoint, err
	***REMOVED***
	endpoint.URL = u

	endpoint.Header = make(http.Header)
	endpoint.HeaderMap = make(map[string]string)
	for _, item := range i.Request.Header ***REMOVED***
		endpoint.Header[item.Key] = append(endpoint.Header[item.Key], item.Value)
		endpoint.HeaderMap[item.Key] = item.Value
	***REMOVED***

	switch i.Request.Body.Mode ***REMOVED***
	case "raw":
		endpoint.Body = []byte(i.Request.Body.Raw)
		endpoint.BodyMap = make(map[string]string)
	case "urlencoded":
		values := make(url.Values)
		endpoint.BodyMap = make(map[string]string)
		for _, field := range i.Request.Body.URLEncoded ***REMOVED***
			if !field.Enabled ***REMOVED***
				continue
			***REMOVED***
			values[field.Key] = append(values[field.Key], field.Value)
			endpoint.BodyMap[field.Key] = field.Value
		***REMOVED***
		endpoint.Body = []byte(values.Encode())
	case "formdata":
		endpoint.Body = make([]byte, 0)
		endpoint.BodyMap = make(map[string]string)
		w := multipart.NewWriter(bytes.NewBuffer(endpoint.Body))
		for _, field := range i.Request.Body.FormData ***REMOVED***
			if !field.Enabled ***REMOVED***
				continue
			***REMOVED***

			if err := w.WriteField(field.Key, field.Value); err != nil ***REMOVED***
				return endpoint, err
			***REMOVED***
			endpoint.BodyMap[field.Key] = field.Value
		***REMOVED***
	***REMOVED***

	if vm != nil ***REMOVED***
		for _, event := range i.Event ***REMOVED***
			if event.Disabled ***REMOVED***
				continue
			***REMOVED***

			script, err := vm.Compile("event", string(event.Script.Exec))
			if err != nil ***REMOVED***
				return endpoint, err
			***REMOVED***

			switch event.Listen ***REMOVED***
			case "test":
				endpoint.Tests = append(endpoint.Tests, script)
			case "prerequest":
				endpoint.PreRequest = append(endpoint.PreRequest, script)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return endpoint, nil
***REMOVED***

func (ep Endpoint) Request() http.Request ***REMOVED***
	return http.Request***REMOVED***
		Method:        ep.Method,
		URL:           ep.URL,
		Header:        ep.Header,
		Body:          ioutil.NopCloser(bytes.NewBuffer(ep.Body)),
		ContentLength: int64(len(ep.Body)),
	***REMOVED***
***REMOVED***
