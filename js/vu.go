package js

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"time"
)

var (
	mRequests = stats.Stat***REMOVED***Name: "requests", Type: stats.HistogramType, Intent: stats.TimeIntent***REMOVED***
	mErrors   = stats.Stat***REMOVED***Name: "errors", Type: stats.CounterType***REMOVED***

	ErrTooManyRedirects = errors.New("too many redirects")
)

type HTTPParams struct ***REMOVED***
	Follow  bool
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

func (u *VU) HTTPRequest(method, url, body string, params HTTPParams, redirects int) (HTTPResponse, error) ***REMOVED***
	parsedURL, err := neturl.Parse(url)
	if err != nil ***REMOVED***
		return HTTPResponse***REMOVED******REMOVED***, err
	***REMOVED***

	req := http.Request***REMOVED***
		Method: method,
		URL:    parsedURL,
		Header: make(http.Header),
	***REMOVED***

	if method == "GET" || method == "HEAD" ***REMOVED***
		req.URL.RawQuery = body
	***REMOVED*** else ***REMOVED***
		// NOT IMPLEMENTED! I'm just testing stuff out.
		// req.SetBodyString(body)
	***REMOVED***

	for key, value := range params.Headers ***REMOVED***
		req.Header[key] = []string***REMOVED***value***REMOVED***
	***REMOVED***

	startTime := time.Now()
	resp, err := u.Client.Do(&req)
	duration := time.Since(startTime)

	var status int
	var respBody []byte
	if err == nil ***REMOVED***
		status = resp.StatusCode
		respBody, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	***REMOVED***

	tags := stats.Tags***REMOVED***
		"url":    url,
		"method": method,
		"status": status,
		"proto":  resp.Proto,
	***REMOVED***

	if !params.Quiet ***REMOVED***
		u.Collector.Add(stats.Sample***REMOVED***
			Stat:   &mRequests,
			Tags:   tags,
			Values: stats.Values***REMOVED***"duration": float64(duration)***REMOVED***,
		***REMOVED***)
	***REMOVED***

	if err != nil ***REMOVED***
		if !params.Quiet ***REMOVED***
			u.Collector.Add(stats.Sample***REMOVED***
				Stat:   &mErrors,
				Tags:   tags,
				Values: stats.Value(1),
			***REMOVED***)
		***REMOVED***
		return HTTPResponse***REMOVED******REMOVED***, err
	***REMOVED***

	// switch resp.StatusCode ***REMOVED***
	// case 301, 302, 303, 307, 308:
	// 	if !params.Follow ***REMOVED***
	// 		break
	// 	***REMOVED***
	// 	if redirects >= u.FollowDepth ***REMOVED***
	// 		return HTTPResponse***REMOVED******REMOVED***, ErrTooManyRedirects
	// 	***REMOVED***

	// 	redirectURL := url
	// 	resp.Header.VisitAll(func(key, value []byte) ***REMOVED***
	// 		if string(key) != "Location" ***REMOVED***
	// 			return
	// 		***REMOVED***

	// 		redirectURL = resolveRedirect(url, string(value))
	// 	***REMOVED***)

	// 	redirectMethod := method
	// 	redirectBody := body
	// 	if status == 301 || status == 302 || status == 303 ***REMOVED***
	// 		redirectMethod = "GET"
	// 		redirectBody = ""
	// 	***REMOVED***
	// 	return u.HTTPRequest(redirectMethod, redirectURL, redirectBody, params, redirects+1)
	// ***REMOVED***

	headers := make(map[string]string)
	for key, vals := range resp.Header ***REMOVED***
		headers[key] = vals[0]
	***REMOVED***

	return HTTPResponse***REMOVED***
		Status:  resp.StatusCode,
		Headers: headers,
		Body:    string(respBody),
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
