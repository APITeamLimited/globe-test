package simple

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
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
		Runner: r,
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
	Runner *Runner
	ID     int64

	Request *http.Request
	Client  *http.Client

	tracer *lib.Tracer
	cTrace *httptrace.ClientTrace
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	resp, err := u.Client.Do(u.Request.WithContext(httptrace.WithClientTrace(ctx, u.cTrace)))
	if err != nil ***REMOVED***
		u.tracer.Done()
		return err
	***REMOVED***

	_, _ = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	trail := u.tracer.Done()

	log.WithFields(log.Fields***REMOVED***
		"duration":   trail.Duration,
		"blocked":    trail.Blocked,
		"looking_up": trail.LookingUp,
		"connecting": trail.Connecting,
		"sending":    trail.Sending,
		"waiting":    trail.Waiting,
		"receiving":  trail.Receiving,
		"reused":     trail.ConnReused,
		"addr":       trail.ConnRemoteAddr,
	***REMOVED***).Info("Request")

	return nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	return nil
***REMOVED***
