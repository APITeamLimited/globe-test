package js

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/runner"
	"github.com/valyala/fasthttp"
	"gopkg.in/olebedev/go-duktape.v2"
	"net/url"
	"strings"
	"time"
)

type apiFunc func(r *Runner, c *duktape.Context, ch chan<- runner.Result) int

func apiHTTPDo(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	method := argString(c, 0)
	if method == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing method in http call")***REMOVED***
		return 0
	***REMOVED***

	u := argString(c, 1)
	if u == "" ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Missing URL in http call")***REMOVED***
		return 0
	***REMOVED***

	body := ""
	switch c.GetType(2) ***REMOVED***
	case duktape.TypeNone, duktape.TypeNull, duktape.TypeUndefined:
	case duktape.TypeString, duktape.TypeNumber, duktape.TypeBoolean:
		body = c.ToString(2)
	case duktape.TypeObject:
		body = c.JsonEncode(2)
	default:
		ch <- runner.Result***REMOVED***Error: errors.New("Unknown type for request body")***REMOVED***
		return 0
	***REMOVED***

	args := struct ***REMOVED***
		Quiet   bool              `json:"quiet"`
		Headers map[string]string `json:"headers"`
	***REMOVED******REMOVED******REMOVED***
	if err := argJSON(c, 3, &args); err != nil ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Invalid arguments to http call")***REMOVED***
		return 0
	***REMOVED***

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	req.Header.SetMethod(method)

	if method == "GET" ***REMOVED***
		if body != "" && body[0] == '***REMOVED***' ***REMOVED***
			rawItems := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
			if err := json.Unmarshal([]byte(body), &rawItems); err != nil ***REMOVED***
				ch <- runner.Result***REMOVED***Error: err***REMOVED***
				return 0
			***REMOVED***
			parts := []string***REMOVED******REMOVED***
			for key, value := range rawItems ***REMOVED***
				value := url.QueryEscape(fmt.Sprint(value))
				parts = append(parts, fmt.Sprintf("%s=%s", key, value))
			***REMOVED***
			req.SetRequestURI(u + "?" + strings.Join(parts, "&"))
		***REMOVED*** else ***REMOVED***
			req.SetRequestURI(u)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		req.SetRequestURI(u)
		req.SetBodyString(body)
	***REMOVED***

	for key, value := range args.Headers ***REMOVED***
		req.Header.Set(key, value)
	***REMOVED***

	startTime := time.Now()
	err := r.Client.Do(req, res)
	duration := time.Since(startTime)

	if !args.Quiet ***REMOVED***
		ch <- runner.Result***REMOVED***Error: err, Time: duration***REMOVED***
	***REMOVED***

	index := c.PushObject()
	***REMOVED***
		c.PushNumber(float64(res.StatusCode()))
		c.PutPropString(-2, "status")

		c.PushString(string(res.Body()))
		c.PutPropString(-2, "body")

		c.PushObject()
		res.Header.VisitAll(func(key, value []byte) ***REMOVED***
			c.PushString(string(value))
			c.PutPropString(-2, string(key))
		***REMOVED***)
		c.PutPropString(-2, "headers")
	***REMOVED***

	c.PushGlobalObject()
	c.GetPropString(-1, "HTTPResponse")
	c.SetPrototype(index)
	c.Pop()

	return 1
***REMOVED***

func apiHTTPSetMaxConnectionsPerHost(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	num := int(argNumber(c, 0))
	if num < 1 ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Max connections per host must be at least 1")***REMOVED***
		return 0
	***REMOVED***
	r.Client.MaxConnsPerHost = num
	return 0
***REMOVED***

func apiLogType(r *Runner, c *duktape.Context, ch chan<- runner.Result) int ***REMOVED***
	kind := argString(c, 0)
	text := argString(c, 1)
	extra := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	if err := argJSON(c, 2, &extra); err != nil ***REMOVED***
		ch <- runner.Result***REMOVED***Error: errors.New("Log context is not an object")***REMOVED***
		return 0
	***REMOVED***

	l := log.WithFields(log.Fields(extra))
	switch kind ***REMOVED***
	case "debug":
		l.Debug(text)
	case "info":
		l.Info(text)
	case "warn":
		l.Warn(text)
	case "error":
		l.Error(text)
	***REMOVED***

	return 0
***REMOVED***
