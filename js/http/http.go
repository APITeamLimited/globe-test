package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loadimpact/speedboat/sampler"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	neturl "net/url"
	"time"
)

type ContextKey int

const (
	clientKey = ContextKey(iota)
)

var ErrNoClient = errors.New("No client in context")

var mDuration *sampler.Metric
var mErrors *sampler.Metric

func init() ***REMOVED***
	mDuration = sampler.Stats("request.duration")
	mErrors = sampler.Counter("request.error")
***REMOVED***

type Args struct ***REMOVED***
	Quiet   bool              `json:"quiet"`
	Headers map[string]string `json:"headers"`
***REMOVED***

type Response struct ***REMOVED***
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
***REMOVED***

func WithDefaultClient(ctx context.Context) context.Context ***REMOVED***
	return WithClient(ctx, &fasthttp.Client***REMOVED******REMOVED***)
***REMOVED***

func WithClient(ctx context.Context, c *fasthttp.Client) context.Context ***REMOVED***
	return context.WithValue(ctx, clientKey, c)
***REMOVED***

func GetClient(ctx context.Context) *fasthttp.Client ***REMOVED***
	return ctx.Value(clientKey).(*fasthttp.Client)
***REMOVED***

func Request(ctx context.Context, method, url, body string, args Args) (Response, error) ***REMOVED***
	client := GetClient(ctx)
	if client == nil ***REMOVED***
		return Response***REMOVED******REMOVED***, ErrNoClient
	***REMOVED***

	if method == "GET" && body != "" ***REMOVED***
		u, err := neturl.Parse(url)
		if err != nil ***REMOVED***
			return Response***REMOVED******REMOVED***, err
		***REMOVED***

		var params map[string]interface***REMOVED******REMOVED***
		if err = json.Unmarshal([]byte(body), &params); err != nil ***REMOVED***
			return Response***REMOVED******REMOVED***, err
		***REMOVED***

		q := u.Query()
		for key, val := range params ***REMOVED***
			q.Set(key, fmt.Sprint(val))
		***REMOVED***
		u.RawQuery = q.Encode()
		url = u.String()
	***REMOVED***

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	if method != "GET" ***REMOVED***
		req.SetBodyString(body)
	***REMOVED***

	for key, value := range args.Headers ***REMOVED***
		req.Header.Set(key, value)
	***REMOVED***

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	startTime := time.Now()
	err := client.Do(req, res)
	duration := time.Since(startTime)

	if !args.Quiet ***REMOVED***
		mDuration.WithFields(sampler.Fields***REMOVED***
			"url":    url,
			"method": method,
			"status": res.StatusCode(),
		***REMOVED***).Duration(duration)
	***REMOVED***

	if err != nil ***REMOVED***
		if !args.Quiet ***REMOVED***
			mErrors.WithFields(sampler.Fields***REMOVED***
				"url":    url,
				"method": method,
				"error":  err,
			***REMOVED***).Int(1)
		***REMOVED***
		return Response***REMOVED******REMOVED***, err
	***REMOVED***

	resHeaders := make(map[string]string)
	res.Header.VisitAll(func(key, value []byte) ***REMOVED***
		resHeaders[string(key)] = string(value)
	***REMOVED***)

	return Response***REMOVED***
		Status:  res.StatusCode(),
		Body:    string(res.Body()),
		Headers: resHeaders,
	***REMOVED***, nil
***REMOVED***
