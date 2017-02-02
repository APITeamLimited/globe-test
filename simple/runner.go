/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package simple

import (
	"context"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
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

type Runner struct ***REMOVED***
	URL       *url.URL
	Transport *http.Transport
	Options   lib.Options

	defaultGroup *lib.Group
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
		defaultGroup: &lib.Group***REMOVED******REMOVED***,
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

func (r *Runner) GetDefaultGroup() *lib.Group ***REMOVED***
	return &lib.Group***REMOVED******REMOVED***
***REMOVED***

func (r Runner) GetOptions() lib.Options ***REMOVED***
	return r.Options
***REMOVED***

func (r *Runner) ApplyOptions(opts lib.Options) ***REMOVED***
	r.Options = r.Options.Apply(opts)
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
	tags := map[string]string***REMOVED***
		"vu":     u.IDString,
		"status": "0",
		"method": "GET",
		"url":    u.URLString,
	***REMOVED***

	resp, err := u.Client.Do(u.Request.WithContext(httptrace.WithClientTrace(ctx, u.cTrace)))
	if err != nil ***REMOVED***
		return u.tracer.Done().Samples(tags), err
	***REMOVED***
	tags["status"] = strconv.Itoa(resp.StatusCode)

	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil ***REMOVED***
		return u.tracer.Done().Samples(tags), err
	***REMOVED***
	_ = resp.Body.Close()

	return u.tracer.Done().Samples(tags), nil
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.IDString = strconv.FormatInt(id, 10)
	return nil
***REMOVED***
