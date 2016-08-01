package js

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/stats"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

var (
	mRequests = stats.Stat***REMOVED***Name: "requests", Type: stats.HistogramType, Intent: stats.TimeIntent***REMOVED***
	mErrors   = stats.Stat***REMOVED***Name: "errors", Type: stats.CounterType***REMOVED***

	ErrTooManyRedirects = errors.New("too many redirects")

	errInternalHandleRedirect = errors.New("[internal] handle redirect")
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

type stringReadCloser struct ***REMOVED***
	*strings.Reader
***REMOVED***

func (stringReadCloser) Close() error ***REMOVED*** return nil ***REMOVED***

func (u *VU) HTTPRequest(method, url, body string, params HTTPParams, redirects int) (HTTPResponse, error) ***REMOVED***
	log.WithFields(log.Fields***REMOVED***"method": method, "url": url, "body": body, "params": params***REMOVED***).Debug("Request")

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
		req.Body = stringReadCloser***REMOVED***strings.NewReader(body)***REMOVED***
		req.ContentLength = int64(len(body))
	***REMOVED***

	for key, value := range params.Headers ***REMOVED***
		req.Header[key] = []string***REMOVED***value***REMOVED***
	***REMOVED***

	startTime := time.Now()
	resp, err := u.Client.Do(&req)
	duration := time.Since(startTime)

	tags := stats.Tags***REMOVED***
		"url":    url,
		"method": method,
		"status": 0,
	***REMOVED***

	var respBody []byte
	if resp != nil ***REMOVED***
		tags["status"] = resp.StatusCode
		tags["proto"] = resp.Proto
		respBody, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	***REMOVED***

	if !params.Quiet ***REMOVED***
		u.Collector.Add(stats.Sample***REMOVED***
			Stat:   &mRequests,
			Tags:   tags,
			Values: stats.Values***REMOVED***"duration": float64(duration)***REMOVED***,
		***REMOVED***)
	***REMOVED***

	switch e := err.(type) ***REMOVED***
	case nil:
		// Do nothing
	case *neturl.Error:
		if e.Err != errInternalHandleRedirect ***REMOVED***
			if !params.Quiet ***REMOVED***
				u.Collector.Add(stats.Sample***REMOVED***Stat: &mErrors, Tags: tags, Values: stats.Value(1)***REMOVED***)
			***REMOVED***
			return HTTPResponse***REMOVED******REMOVED***, err
		***REMOVED***

		if !params.Follow ***REMOVED***
			break
		***REMOVED***

		if redirects >= u.FollowDepth ***REMOVED***
			return HTTPResponse***REMOVED******REMOVED***, ErrTooManyRedirects
		***REMOVED***

		redirectURL := resolveRedirect(url, resp.Header.Get("Location"))
		redirectMethod := method
		redirectBody := ""
		if resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 303 ***REMOVED***
			redirectMethod = "GET"
			if redirectMethod != method ***REMOVED***
				redirectBody = ""
			***REMOVED***
		***REMOVED***

		return u.HTTPRequest(redirectMethod, redirectURL, redirectBody, params, redirects+1)
	default:
		if !params.Quiet ***REMOVED***
			u.Collector.Add(stats.Sample***REMOVED***Stat: &mErrors, Tags: tags, Values: stats.Value(1)***REMOVED***)
		***REMOVED***
		return HTTPResponse***REMOVED******REMOVED***, err
	***REMOVED***

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
