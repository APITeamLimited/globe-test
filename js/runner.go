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

package js

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
	"github.com/robertkrimen/otto"
	"math"
	"net"
	"net/http"
	"strconv"
	"time"
)

var ErrDefaultExport = errors.New("you must export a 'default' function")

const entrypoint = "__$$entrypoint$$__"

type Runner struct ***REMOVED***
	Runtime      *Runtime
	DefaultGroup *lib.Group
	Options      lib.Options

	HTTPTransport *http.Transport
***REMOVED***

func NewRunner(rt *Runtime, exports otto.Value) (*Runner, error) ***REMOVED***
	expObj := exports.Object()
	if expObj == nil ***REMOVED***
		return nil, ErrDefaultExport
	***REMOVED***

	// Values "remember" which VM they belong to, so to get a callable that works across VM copies,
	// we have to stick it in the global scope, then retrieve it again from the new instance.
	callable, err := expObj.Get("default")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if !callable.IsFunction() ***REMOVED***
		return nil, ErrDefaultExport
	***REMOVED***
	if err := rt.VM.Set(entrypoint, callable); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defaultGroup, err := lib.NewGroup("", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	r := &Runner***REMOVED***
		Runtime:      rt,
		DefaultGroup: defaultGroup,
		Options:      rt.Options,
		HTTPTransport: &http.Transport***REMOVED***
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer***REMOVED***
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			***REMOVED***).DialContext,
			TLSClientConfig:     &tls.Config***REMOVED******REMOVED***,
			MaxIdleConns:        math.MaxInt32,
			MaxIdleConnsPerHost: math.MaxInt32,
		***REMOVED***,
	***REMOVED***

	return r, nil
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	u := &VU***REMOVED***
		runner: r,
		vm:     r.Runtime.VM.Copy(),
		group:  r.DefaultGroup,
	***REMOVED***

	u.CookieJar = lib.NewCookieJar()
	u.HTTPClient = &http.Client***REMOVED***
		Transport:     r.HTTPTransport,
		CheckRedirect: u.checkRedirect,
		Jar:           u.CookieJar,
	***REMOVED***

	callable, err := u.vm.Get(entrypoint)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	u.callable = callable

	if err := u.vm.Set("__jsapi__", &JSAPI***REMOVED***u***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return u, nil
***REMOVED***

func (r *Runner) GetDefaultGroup() *lib.Group ***REMOVED***
	return r.DefaultGroup
***REMOVED***

func (r *Runner) GetOptions() lib.Options ***REMOVED***
	return r.Options
***REMOVED***

func (r *Runner) ApplyOptions(opts lib.Options) ***REMOVED***
	r.Options = r.Options.Apply(opts)
	r.HTTPTransport.TLSClientConfig.InsecureSkipVerify = opts.InsecureSkipTLSVerify.Bool
***REMOVED***

type VU struct ***REMOVED***
	ID        int64
	IDString  string
	Iteration int64
	Samples   []stats.Sample

	runner   *Runner
	vm       *otto.Otto
	callable otto.Value

	HTTPClient *http.Client
	CookieJar  *lib.CookieJar

	started time.Time
	ctx     context.Context
	group   *lib.Group
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	u.CookieJar.Clear()

	if err := u.vm.Set("__ITER", u.Iteration); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	u.started = time.Now()
	u.Samples = []stats.Sample***REMOVED***
		***REMOVED***Time: u.started, Metric: metrics.Iterations, Value: 1.0***REMOVED***,
	***REMOVED***
	u.ctx = ctx
	_, err := u.callable.Call(otto.UndefinedValue())
	u.ctx = nil

	u.Iteration++

	samples := u.Samples
	u.Samples = nil
	return samples, err
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.IDString = strconv.FormatInt(u.ID, 10)
	u.Iteration = 0

	if err := u.vm.Set("__VU", u.ID); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := u.vm.Set("__ITER", u.Iteration); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (u *VU) checkRedirect(req *http.Request, via []*http.Request) error ***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"from": via[len(via)-1].URL.String(),
		"to":   req.URL.String(),
	***REMOVED***).Debug("-> Redirect")
	if int64(len(via)) >= u.runner.Options.MaxRedirects.Int64 ***REMOVED***
		return errors.New(fmt.Sprintf("stopped after %d redirects", u.runner.Options.MaxRedirects.Int64))
	***REMOVED***
	return nil
***REMOVED***
