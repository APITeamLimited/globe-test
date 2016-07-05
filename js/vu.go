package js

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"time"
)

var (
	mRequests = stats.Stat***REMOVED***Name: "requests", Type: stats.HistogramType, Intent: stats.TimeIntent***REMOVED***
	mErrors   = stats.Stat***REMOVED***Name: "errors", Type: stats.CounterType***REMOVED***
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
	obj, err := Make(vm, "HTTPResponse")
	if err != nil ***REMOVED***
		return otto.UndefinedValue(), err
	***REMOVED***

	obj.Set("status", res.Status)
	obj.Set("headers", res.Headers)
	obj.Set("body", res.Body)

	return vm.ToValue(obj)
***REMOVED***

func (u *VU) HTTPRequest(method, url, body string, params HTTPParams) (HTTPResponse, error) ***REMOVED***
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)

	if method == "GET" || method == "HEAD" ***REMOVED***
		req.SetRequestURI(putBodyInURL(url, body))
	***REMOVED*** else ***REMOVED***
		req.SetRequestURI(url)
		req.SetBodyString(body)
	***REMOVED***

	for key, value := range params.Headers ***REMOVED***
		req.Header.Set(key, value)
	***REMOVED***

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	startTime := time.Now()
	err := u.Client.Do(req, resp)
	duration := time.Since(startTime)

	if !params.Quiet ***REMOVED***
		u.Collector.Add(stats.Point***REMOVED***
			Stat: &mRequests,
			Tags: stats.Tags***REMOVED***
				"url":    url,
				"method": method,
				"status": resp.StatusCode(),
			***REMOVED***,
			Values: stats.Values***REMOVED***"duration": float64(duration)***REMOVED***,
		***REMOVED***)
	***REMOVED***

	if err != nil ***REMOVED***
		if !params.Quiet ***REMOVED***
			u.Collector.Add(stats.Point***REMOVED***
				Stat: &mErrors,
				Tags: stats.Tags***REMOVED***
					"url":    url,
					"method": method,
					"status": resp.StatusCode(),
				***REMOVED***,
				Values: stats.Value(1),
			***REMOVED***)
		***REMOVED***
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
