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
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/robertkrimen/otto"
)

const (
	DefaultMaxRedirect = 10
)

var ErrDefaultExport = errors.New("you must export a 'default' function")

const entrypoint = "__$$entrypoint$$__"

type Runner struct ***REMOVED***
	Runtime      *Runtime
	DefaultGroup *lib.Group
	Groups       []*lib.Group
	Checks       []*lib.Check
	Options      lib.Options
	SrcData      *lib.SourceData

	HTTPTransport *http.Transport

	groupIDCounter int64
	groupsMutex    sync.Mutex
	checkIDCounter int64
	checksMutex    sync.Mutex
***REMOVED***

func NewRunner(runtime *Runtime, src *lib.SourceData, exports otto.Value) (*Runner, error) ***REMOVED***
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
	if err := runtime.VM.Set(entrypoint, callable); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	r := &Runner***REMOVED***
		Runtime: runtime,
		Options: runtime.Options,
		SrcData: src,
		HTTPTransport: &http.Transport***REMOVED***
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer***REMOVED***
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			***REMOVED***).DialContext,
			MaxIdleConns:        math.MaxInt32,
			MaxIdleConnsPerHost: math.MaxInt32,
		***REMOVED***,
	***REMOVED***
	r.DefaultGroup = lib.NewGroup("", nil, nil)
	r.Groups = []*lib.Group***REMOVED***r.DefaultGroup***REMOVED***

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

	if err := u.vm.Set("__jsapi__", JSAPI***REMOVED***u***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return u, nil
***REMOVED***

func (r *Runner) GetGroups() []*lib.Group ***REMOVED***
	return r.Groups
***REMOVED***

func (r *Runner) GetChecks() []*lib.Check ***REMOVED***
	return r.Checks
***REMOVED***

func (r *Runner) GetOptions() lib.Options ***REMOVED***
	return r.Options
***REMOVED***

func (r *Runner) ApplyOptions(opts lib.Options) ***REMOVED***
	r.Options = r.Options.Apply(opts)
***REMOVED***

func (r *Runner) GetSourceData() *lib.SourceData ***REMOVED***
	return r.SrcData
***REMOVED***

type VU struct ***REMOVED***
	ID       int64
	IDString string
	Samples  []stats.Sample
	Taint    bool

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

	u.started = time.Now()
	u.ctx = ctx
	_, err := u.callable.Call(otto.UndefinedValue())
	u.ctx = nil

	if u.Taint ***REMOVED***
		u.Taint = false
		if err == nil ***REMOVED***
			err = lib.ErrVUWantsTaint
		***REMOVED***
	***REMOVED***

	samples := u.Samples
	u.Samples = nil
	return samples, err
***REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.IDString = strconv.FormatInt(u.ID, 10)
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
