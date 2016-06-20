package js

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"time"
)

type HTTPParams struct ***REMOVED***
	Quiet   bool
	Headers map[string]string
***REMOVED***

type HTTPResponse struct ***REMOVED***
	Status  int
	Headers map[string]string
	Body    string
***REMOVED***

func (res HTTPResponse) ToValue(vm *otto.Otto) (otto.Value, error) ***REMOVED***
	return vm.ToValue(map[string]interface***REMOVED******REMOVED******REMOVED***
		"status":  res.Status,
		"headers": res.Headers,
		"body":    res.Body,
	***REMOVED***)
***REMOVED***

func (u *VU) HTTPRequest(method, url, body string, params HTTPParams) (HTTPResponse, error) ***REMOVED***
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	if body != "" ***REMOVED***
		req.SetBodyString(body)
	***REMOVED***

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	startTime := time.Now()
	err := u.Client.Do(req, resp)
	duration := time.Since(startTime)

	u.Runner.mDuration.WithFields(sampler.Fields***REMOVED***
		"url":    u.Runner.Test.URL,
		"method": "GET",
		"status": resp.StatusCode(),
	***REMOVED***).Duration(duration)

	if err != nil ***REMOVED***
		u.Runner.mErrors.WithField("url", u.Runner.Test.URL).Int(1)
		return HTTPResponse***REMOVED******REMOVED***, err
	***REMOVED***

	headers := make(map[string]string)
	resp.Header.VisitAll(func(key []byte, value []byte) ***REMOVED***
		headers[string(key)] = string(value)
	***REMOVED***)

	return HTTPResponse***REMOVED***
		Status:  resp.StatusCode(),
		Headers: headers,
		Body:    string(resp.Body()),
	***REMOVED***, nil
***REMOVED***

func (u *VU) Sleep(t float64) ***REMOVED***
	time.Sleep(time.Duration(t * float64(time.Second)))
***REMOVED***

func (u *VU) Log(level, msg string, fields map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	e := u.Runner.logger.WithFields(log.Fields(fields))

	switch level ***REMOVED***
	case "debug":
		e.Debug(msg)
	case "info":
		e.Info(msg)
	case "warn":
		e.Warn(msg)
	case "error":
		e.Error(msg)
	***REMOVED***
***REMOVED***
