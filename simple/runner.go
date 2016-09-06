package simple

import (
	"context"
	"github.com/loadimpact/speedboat/lib"
	"github.com/loadimpact/speedboat/stats"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
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

type Runner struct ***REMOVED***
	URL       *url.URL
	Transport *http.Transport
***REMOVED***

func New(rawurl string) (*Runner, error) ***REMOVED***
	u, err := url.Parse(rawurl)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Runner***REMOVED***
		URL: u,
		Transport: &http.Transport***REMOVED***
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer***REMOVED***
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			***REMOVED***).DialContext,
			MaxIdleConns:        math.MaxInt32,
			MaxIdleConnsPerHost: math.MaxInt32,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	tracer := &lib.Tracer***REMOVED******REMOVED***

	return &VU***REMOVED***
		Runner:    r,
		URLString: r.URL.String(),
		Request: &http.Request***REMOVED***
			Method: "GET",
			URL:    r.URL,
		***REMOVED***,
		Client: &http.Client***REMOVED***
			Transport: r.Transport,
		***REMOVED***,
		tracer: tracer,
		cTrace: tracer.Trace(),
	***REMOVED***, nil
***REMOVED***

type VU struct ***REMOVED***
	Runner   *Runner
	ID       int64
	IDString string

	URLString string
	Request   *http.Request
	Client    *http.Client

	tracer *lib.Tracer
	cTrace *httptrace.ClientTrace
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	resp, err := u.Client.Do(u.Request.WithContext(httptrace.WithClientTrace(ctx, u.cTrace)))
	if err != nil ***REMOVED***
		u.tracer.Done()
		return nil, err
	***REMOVED***

	_, _ = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	trail := u.tracer.Done()

	tags := map[string]string***REMOVED***
		"vu":     u.IDString,
		"method": "GET",
		"url":    u.URLString,
	***REMOVED***

	t := time.Now()
	return []stats.Sample***REMOVED***
		stats.Sample***REMOVED***Metric: MetricReqs, Time: t, Tags: tags, Value: 1***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqDuration, Time: t, Tags: tags, Value: float64(trail.Duration)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqBlocked, Time: t, Tags: tags, Value: float64(trail.Blocked)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqLookingUp, Time: t, Tags: tags, Value: float64(trail.LookingUp)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqConnecting, Time: t, Tags: tags, Value: float64(trail.Connecting)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqSending, Time: t, Tags: tags, Value: float64(trail.Sending)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqWaiting, Time: t, Tags: tags, Value: float64(trail.Waiting)***REMOVED***,
		stats.Sample***REMOVED***Metric: MetricReqReceiving, Time: t, Tags: tags, Value: float64(trail.Receiving)***REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.IDString = strconv.FormatInt(id, 10)
	return nil
***REMOVED***
