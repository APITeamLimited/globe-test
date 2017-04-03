package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	neturl "net/url"
	"strconv"
	"strings"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js2/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
)

type HTTPResponseTimings struct ***REMOVED***
	Duration, Blocked, LookingUp, Connecting, Sending, Waiting, Receiving float64
***REMOVED***

type HTTPResponse struct ***REMOVED***
	URL     string
	Status  int
	Headers map[string]string
	Body    string
	Timings HTTPResponseTimings

	runtime    *goja.Runtime
	cachedJSON goja.Value
***REMOVED***

func (res *HTTPResponse) Json() goja.Value ***REMOVED***
	if res.cachedJSON == nil ***REMOVED***
		var v interface***REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(res.Body), &v); err != nil ***REMOVED***
			common.Throw(res.runtime, err)
		***REMOVED***
		res.cachedJSON = res.runtime.ToValue(v)
	***REMOVED***
	return res.cachedJSON
***REMOVED***

type HTTP struct***REMOVED******REMOVED***

func (*HTTP) Request(ctx context.Context, method, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	var bodyReader io.Reader
	var contentType string
	if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) ***REMOVED***
		var data map[string]goja.Value
		if rt.ExportTo(args[0], &data) == nil ***REMOVED***
			bodyQuery := make(neturl.Values, len(data))
			for k, v := range data ***REMOVED***
				bodyQuery.Set(k, v.String())
			***REMOVED***
			bodyReader = bytes.NewBufferString(bodyQuery.Encode())
			contentType = "application/x-www-form-urlencoded"
		***REMOVED*** else ***REMOVED***
			bodyReader = bytes.NewBufferString(args[0].String())
		***REMOVED***
	***REMOVED***

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if contentType != "" ***REMOVED***
		req.Header.Set("Content-Type", contentType)
	***REMOVED***

	tags := map[string]string***REMOVED***
		"status": "0",
		"method": method,
		"url":    url,
		"group":  state.Group.Path,
	***REMOVED***

	if len(args) > 1 ***REMOVED***
		paramsV := args[1]
		if !goja.IsUndefined(paramsV) && !goja.IsNull(paramsV) ***REMOVED***
			params := paramsV.ToObject(rt)
			for _, k := range params.Keys() ***REMOVED***
				switch k ***REMOVED***
				case "headers":
					headersV := params.Get(k)
					if goja.IsUndefined(headersV) || goja.IsNull(headersV) ***REMOVED***
						continue
					***REMOVED***
					headers := headersV.ToObject(rt)
					if headers == nil ***REMOVED***
						continue
					***REMOVED***
					for _, key := range headers.Keys() ***REMOVED***
						req.Header.Set(key, headers.Get(key).String())
					***REMOVED***
				case "tags":
					tagsV := params.Get(k)
					if goja.IsUndefined(tagsV) || goja.IsNull(tagsV) ***REMOVED***
						continue
					***REMOVED***
					tagObj := tagsV.ToObject(rt)
					if tagObj == nil ***REMOVED***
						continue
					***REMOVED***
					for _, key := range tagObj.Keys() ***REMOVED***
						tags[key] = tagObj.Get(key).String()
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	client := http.Client***REMOVED******REMOVED***
	tracer := lib.Tracer***REMOVED******REMOVED***
	res, err := client.Do(req.WithContext(httptrace.WithClientTrace(ctx, tracer.Trace())))
	if err != nil ***REMOVED***
		state.Samples = append(state.Samples, tracer.Done().Samples(tags)...)
		return nil, err
	***REMOVED***

	body, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		state.Samples = append(state.Samples, tracer.Done().Samples(tags)...)
		return nil, err
	***REMOVED***
	_ = res.Body.Close()
	trail := tracer.Done()

	tags["status"] = strconv.Itoa(res.StatusCode)
	state.Samples = append(state.Samples, trail.Samples(tags)...)

	headers := make(map[string]string, len(res.Header))
	for k, vs := range res.Header ***REMOVED***
		headers[k] = strings.Join(vs, ", ")
	***REMOVED***
	return &HTTPResponse***REMOVED***
		URL:     res.Request.URL.String(),
		Status:  res.StatusCode,
		Headers: headers,
		Body:    string(body),
		Timings: HTTPResponseTimings***REMOVED***
			Duration:   stats.D(trail.Duration),
			Blocked:    stats.D(trail.Blocked),
			LookingUp:  stats.D(trail.LookingUp),
			Connecting: stats.D(trail.Connecting),
			Sending:    stats.D(trail.Sending),
			Waiting:    stats.D(trail.Waiting),
			Receiving:  stats.D(trail.Receiving),
		***REMOVED***,

		runtime: rt,
	***REMOVED***, nil
***REMOVED***

func (http *HTTP) Get(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, "GET", url, args...)
***REMOVED***

func (http *HTTP) Head(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	// The body argument is always undefined for GETs and HEADs.
	args = append([]goja.Value***REMOVED***goja.Undefined()***REMOVED***, args...)
	return http.Request(ctx, "HEAD", url, args...)
***REMOVED***

func (http *HTTP) Post(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "POST", url, args...)
***REMOVED***

func (http *HTTP) Put(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "PUT", url, args...)
***REMOVED***

func (http *HTTP) Patch(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "PATCH", url, args...)
***REMOVED***

func (http *HTTP) Del(ctx context.Context, url string, args ...goja.Value) (*HTTPResponse, error) ***REMOVED***
	return http.Request(ctx, "DELETE", url, args...)
***REMOVED***
