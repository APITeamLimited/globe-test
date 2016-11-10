package js

import (
	// "github.com/robertkrimen/otto"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

var (
	MetricReqs          = stats.New("http_reqs", stats.Counter)
	MetricReqDuration   = stats.New("http_req_duration", stats.Trend, stats.Time)
	MetricReqBlocked    = stats.New("http_req_blocked", stats.Trend, stats.Time)
	MetricReqLookingUp  = stats.New("http_req_looking_up", stats.Trend, stats.Time)
	MetricReqConnecting = stats.New("http_req_connecting", stats.Trend, stats.Time)
	MetricReqSending    = stats.New("http_req_sending", stats.Trend, stats.Time)
	MetricReqWaiting    = stats.New("http_req_waiting", stats.Trend, stats.Time)
	MetricReqReceiving  = stats.New("http_req_receiving", stats.Trend, stats.Time)
)

type HTTPResponse struct ***REMOVED***
	Status int
***REMOVED***

func (a JSAPI) HTTPRequest(method, url, body string, params map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"method": method,
		"url":    url,
		"body":   body,
	***REMOVED***).Debug("JS: Request")
	bodyReader := io.Reader(nil)
	if body != "" ***REMOVED***
		bodyReader = strings.NewReader(body)
	***REMOVED***
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***

	if h, ok := params["headers"]; ok ***REMOVED***
		headers, ok := h.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			panic(a.vu.vm.MakeTypeError("headers must be an object"))
		***REMOVED***
		for key, v := range headers ***REMOVED***
			value, ok := v.(string)
			if !ok ***REMOVED***
				panic(a.vu.vm.MakeTypeError("header values must be strings"))
			***REMOVED***
			req.Header.Set(key, value)
		***REMOVED***
	***REMOVED***

	tracer := lib.Tracer***REMOVED******REMOVED***
	res, err := a.vu.HTTPClient.Do(req.WithContext(httptrace.WithClientTrace(a.vu.ctx, tracer.Trace())))
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		throw(a.vu.vm, err)
	***REMOVED***
	res.Body.Close()

	trail := tracer.Done()
	t := time.Now()
	tags := map[string]string***REMOVED***
		"vu":     a.vu.IDString,
		"method": method,
		"url":    url,
	***REMOVED***
	a.vu.Samples = append(a.vu.Samples,
		stats.Sample***REMOVED***Metric: MetricReqs, Time: t, Tags: tags, Value: 1***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqDuration, Time: t, Tags: tags, Value: float64(trail.Duration)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqBlocked, Time: t, Tags: tags, Value: float64(trail.Blocked)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqLookingUp, Time: t, Tags: tags, Value: float64(trail.LookingUp)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqConnecting, Time: t, Tags: tags, Value: float64(trail.Connecting)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqSending, Time: t, Tags: tags, Value: float64(trail.Sending)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqWaiting, Time: t, Tags: tags, Value: float64(trail.Waiting)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqReceiving, Time: t, Tags: tags, Value: float64(trail.Receiving)***REMOVED***,
	)

	headers := make(map[string]string)
	for k, v := range res.Header ***REMOVED***
		headers[k] = strings.Join(v, ", ")
	***REMOVED***
	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"status":  res.StatusCode,
		"body":    string(resBody),
		"headers": headers,
	***REMOVED***
***REMOVED***

func (a JSAPI) HTTPSetMaxRedirects(n int) ***REMOVED***
	a.vu.MaxRedirects = n
***REMOVED***
