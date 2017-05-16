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
	"crypto/tls"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
)

type Runner struct ***REMOVED***
	URL       *url.URL
	Transport *http.Transport
	Options   lib.Options

	defaultGroup *lib.Group
***REMOVED***

func New(u *url.URL) (*Runner, error) ***REMOVED***
	return &Runner***REMOVED***
		URL: u,
		Transport: &http.Transport***REMOVED***
			Proxy: http.ProxyFromEnvironment,
			DialContext: netext.NewDialer(net.Dialer***REMOVED***
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			***REMOVED***).DialContext,
			TLSClientConfig:     &tls.Config***REMOVED******REMOVED***,
			MaxIdleConns:        math.MaxInt32,
			MaxIdleConnsPerHost: math.MaxInt32,
		***REMOVED***,
		defaultGroup: &lib.Group***REMOVED******REMOVED***,
	***REMOVED***, nil
***REMOVED***

func (r *Runner) MakeArchive() *lib.Archive ***REMOVED***
	return &lib.Archive***REMOVED***
		Type:     "url",
		Filename: r.URL.String(),
	***REMOVED***
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
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
		tracer: &netext.Tracer***REMOVED******REMOVED***,
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
	r.Transport.TLSClientConfig.InsecureSkipVerify = opts.InsecureSkipTLSVerify.Bool
***REMOVED***

type VU struct ***REMOVED***
	Runner   *Runner
	ID       int64
	IDString string

	URLString string
	Request   *http.Request
	Client    *http.Client

	tracer *netext.Tracer
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	tags := map[string]string***REMOVED***
		"vu":     u.IDString,
		"status": "0",
		"method": "GET",
		"url":    u.URLString,
	***REMOVED***

	resp, err := u.Client.Do(u.Request.WithContext(netext.WithTracer(ctx, u.tracer)))
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
