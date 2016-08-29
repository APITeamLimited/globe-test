package simple

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/lib"
	"net"
	"net/http"
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
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	return &VU***REMOVED***
		Runner: r,
		Request: http.Request***REMOVED***
			Method: "GET",
			URL:    r.URL,
		***REMOVED***,
		Client: http.Client***REMOVED***
			Transport: r.Transport,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

type VU struct ***REMOVED***
	Runner *Runner
	ID     int64

	Request http.Request
	Client  http.Client
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) error ***REMOVED***
	log.WithField("id", u.ID).Info("Running")
	return nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	return nil
***REMOVED***
