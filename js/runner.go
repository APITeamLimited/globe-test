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
	"encoding/json"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	"github.com/oxtoacart/bpool"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/viki-org/dnscache"
	"golang.org/x/net/http2"
	"golang.org/x/time/rate"
)

var errInterrupt = errors.New("context cancelled")

// Ensure Runner implements the lib.Runner interface
var _ lib.Runner = &Runner***REMOVED******REMOVED***

type Runner struct ***REMOVED***
	Bundle       *Bundle
	Logger       *log.Logger
	defaultGroup *lib.Group

	BaseDialer net.Dialer
	Resolver   *dnscache.Resolver
	RPSLimit   *rate.Limiter

	setupData interface***REMOVED******REMOVED***
***REMOVED***

func New(src *lib.SourceData, fs afero.Fs, rtOpts lib.RuntimeOptions) (*Runner, error) ***REMOVED***
	bundle, err := NewBundle(src, fs, rtOpts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewFromBundle(bundle)
***REMOVED***

func NewFromArchive(arc *lib.Archive, rtOpts lib.RuntimeOptions) (*Runner, error) ***REMOVED***
	bundle, err := NewBundleFromArchive(arc, rtOpts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewFromBundle(bundle)
***REMOVED***

func NewFromBundle(b *Bundle) (*Runner, error) ***REMOVED***
	defaultGroup, err := lib.NewGroup("", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	r := &Runner***REMOVED***
		Bundle:       b,
		Logger:       log.StandardLogger(),
		defaultGroup: defaultGroup,
		BaseDialer: net.Dialer***REMOVED***
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		***REMOVED***,
		Resolver: dnscache.New(0),
	***REMOVED***
	r.SetOptions(r.Bundle.Options)
	return r, nil
***REMOVED***

func (r *Runner) MakeArchive() *lib.Archive ***REMOVED***
	return r.Bundle.MakeArchive()
***REMOVED***

func (r *Runner) NewVU() (lib.VU, error) ***REMOVED***
	vu, err := r.newVU()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return lib.VU(vu), nil
***REMOVED***

func (r *Runner) newVU() (*VU, error) ***REMOVED***
	// Instantiate a new bundle, make a VU out of it.
	bi, err := r.Bundle.Instantiate()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var cipherSuites []uint16
	if r.Bundle.Options.TLSCipherSuites != nil ***REMOVED***
		cipherSuites = *r.Bundle.Options.TLSCipherSuites
	***REMOVED***

	var tlsVersions lib.TLSVersions
	if r.Bundle.Options.TLSVersion != nil ***REMOVED***
		tlsVersions = *r.Bundle.Options.TLSVersion
	***REMOVED***

	tlsAuth := r.Bundle.Options.TLSAuth
	certs := make([]tls.Certificate, len(tlsAuth))
	nameToCert := make(map[string]*tls.Certificate)
	for i, auth := range tlsAuth ***REMOVED***
		for _, name := range auth.Domains ***REMOVED***
			cert, err := auth.Certificate()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			certs[i] = *cert
			nameToCert[name] = &certs[i]
		***REMOVED***
	***REMOVED***

	dialer := &netext.Dialer***REMOVED***
		Dialer:    r.BaseDialer,
		Resolver:  r.Resolver,
		Blacklist: r.Bundle.Options.BlacklistIPs,
		Hosts:     r.Bundle.Options.Hosts,
	***REMOVED***
	tlsConfig := &tls.Config***REMOVED***
		InsecureSkipVerify: r.Bundle.Options.InsecureSkipTLSVerify.Bool,
		CipherSuites:       cipherSuites,
		MinVersion:         uint16(tlsVersions.Min),
		MaxVersion:         uint16(tlsVersions.Max),
		Certificates:       certs,
		NameToCertificate:  nameToCert,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	***REMOVED***
	transport := &http.Transport***REMOVED***
		Proxy:              http.ProxyFromEnvironment,
		TLSClientConfig:    tlsConfig,
		DialContext:        dialer.DialContext,
		DisableCompression: true,
		DisableKeepAlives:  r.Bundle.Options.NoConnectionReuse.Bool,
	***REMOVED***
	_ = http2.ConfigureTransport(transport)

	vu := &VU***REMOVED***
		BundleInstance: *bi,
		Runner:         r,
		HTTPTransport:  netext.NewHTTPTransport(transport),
		Dialer:         dialer,
		TLSConfig:      tlsConfig,
		Console:        NewConsole(),
		BPool:          bpool.NewBufferPool(100),
	***REMOVED***
	vu.Runtime.Set("console", common.Bind(vu.Runtime, vu.Console, vu.Context))
	common.BindToGlobal(vu.Runtime, map[string]interface***REMOVED******REMOVED******REMOVED***
		"open": func() ***REMOVED***
			common.Throw(vu.Runtime, errors.New("\"open\" function is only available to the init code (aka global scope), see https://docs.k6.io/docs/test-life-cycle for more information"))
		***REMOVED***,
	***REMOVED***)

	// Give the VU an initial sense of identity.
	if err := vu.Reconfigure(0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return vu, nil
***REMOVED***

func (r *Runner) Setup(ctx context.Context) error ***REMOVED***
	setupCtx, setupCancel := context.WithTimeout(
		ctx,
		time.Duration(r.Bundle.Options.SetupTimeout.Duration),
	)
	defer setupCancel()

	v, err := r.runPart(setupCtx, "setup", nil)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "setup")
	***REMOVED***
	data, err := json.Marshal(v.Export())
	if err != nil ***REMOVED***
		return errors.Wrap(err, "setup")
	***REMOVED***
	return json.Unmarshal(data, &r.setupData)
***REMOVED***

// GetSetupData returns the setup data if Setup() was specified and executed, nil otherwise
func (r *Runner) GetSetupData() interface***REMOVED******REMOVED*** ***REMOVED***
	return r.setupData
***REMOVED***

// SetSetupData saves the externally supplied setup data in the runner, so it can be used in VUs
func (r *Runner) SetSetupData(data interface***REMOVED******REMOVED***) ***REMOVED***
	r.setupData = data
***REMOVED***

func (r *Runner) Teardown(ctx context.Context) error ***REMOVED***
	teardownCtx, teardownCancel := context.WithTimeout(
		ctx,
		time.Duration(r.Bundle.Options.TeardownTimeout.Duration),
	)
	defer teardownCancel()

	_, err := r.runPart(teardownCtx, "teardown", r.setupData)
	return err
***REMOVED***

func (r *Runner) GetDefaultGroup() *lib.Group ***REMOVED***
	return r.defaultGroup
***REMOVED***

func (r *Runner) GetOptions() lib.Options ***REMOVED***
	return r.Bundle.Options
***REMOVED***

func (r *Runner) SetOptions(opts lib.Options) ***REMOVED***
	r.Bundle.Options = opts

	r.RPSLimit = nil
	if rps := opts.RPS; rps.Valid ***REMOVED***
		r.RPSLimit = rate.NewLimiter(rate.Limit(rps.Int64), 1)
	***REMOVED***
***REMOVED***

// Runs an exported function in its own temporary VU, optionally with an argument. Execution is
// interrupted if the context expires. No error is returned if the part does not exist.
func (r *Runner) runPart(ctx context.Context, name string, arg interface***REMOVED******REMOVED***) (goja.Value, error) ***REMOVED***
	vu, err := r.newVU()
	if err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***
	exp := vu.Runtime.Get("exports").ToObject(vu.Runtime)
	if exp == nil ***REMOVED***
		return goja.Undefined(), nil
	***REMOVED***
	fn, ok := goja.AssertFunction(exp.Get(name))
	if !ok ***REMOVED***
		return goja.Undefined(), nil
	***REMOVED***

	ctx, cancel := context.WithCancel(ctx)
	go func() ***REMOVED***
		<-ctx.Done()
		vu.Runtime.Interrupt(errInterrupt)
	***REMOVED***()
	v, _, err := vu.runFn(ctx, fn, vu.Runtime.ToValue(arg))
	cancel()
	return v, err
***REMOVED***

type VU struct ***REMOVED***
	BundleInstance

	Runner        *Runner
	HTTPTransport *netext.HTTPTransport
	Dialer        *netext.Dialer
	TLSConfig     *tls.Config
	ID            int64
	Iteration     int64

	Console *Console
	BPool   *bpool.BufferPool

	setupData goja.Value

	// A VU will track the last context it was called with for cancellation.
	// Note that interruptTrackedCtx is the context that is currently being tracked, while
	// interruptCancel cancels an unrelated context that terminates the tracking goroutine
	// without triggering an interrupt (for if the context changes).
	// There are cleaner ways of handling the interruption problem, but this is a hot path that
	// needs to be called thousands of times per second, which rules out anything that spawns a
	// goroutine per call.
	interruptTrackedCtx context.Context
	interruptCancel     context.CancelFunc
***REMOVED***

// Verify that VU implements lib.VU
var _ lib.VU = &VU***REMOVED******REMOVED***

func (u *VU) Reconfigure(id int64) error ***REMOVED***
	u.ID = id
	u.Iteration = 0
	u.Runtime.Set("__VU", u.ID)
	return nil
***REMOVED***

func (u *VU) RunOnce(ctx context.Context) ([]stats.SampleContainer, error) ***REMOVED***
	// Track the context and interrupt JS execution if it's cancelled.
	if u.interruptTrackedCtx != ctx ***REMOVED***
		interCtx, interCancel := context.WithCancel(context.Background())
		if u.interruptCancel != nil ***REMOVED***
			u.interruptCancel()
		***REMOVED***
		u.interruptCancel = interCancel
		u.interruptTrackedCtx = ctx
		go func() ***REMOVED***
			select ***REMOVED***
			case <-interCtx.Done():
			case <-ctx.Done():
				u.Runtime.Interrupt(errInterrupt)
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// Lazily JS-ify setupData on first run. This is lightweight enough that we can get away with
	// it, and alleviates a problem where setupData wouldn't get populated properly if NewVU() was
	// called before Setup(), which is hard to avoid with how the Executor works w/o complicating
	// the local executor further by deferring SetVUsMax() calls to within the Run() function.
	if u.setupData == nil && u.Runner.setupData != nil ***REMOVED***
		u.setupData = u.Runtime.ToValue(u.Runner.setupData)
	***REMOVED***

	// Call the default function.
	_, state, err := u.runFn(ctx, u.Default, u.setupData)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return state.Samples, nil
***REMOVED***

func (u *VU) runFn(ctx context.Context, fn goja.Callable, args ...goja.Value) (goja.Value, *common.State, error) ***REMOVED***
	cookieJar, err := cookiejar.New(nil)
	if err != nil ***REMOVED***
		return goja.Undefined(), nil, err
	***REMOVED***

	state := &common.State***REMOVED***
		Logger:        u.Runner.Logger,
		Options:       u.Runner.Bundle.Options,
		Group:         u.Runner.defaultGroup,
		HTTPTransport: u.HTTPTransport,
		Dialer:        u.Dialer,
		TLSConfig:     u.TLSConfig,
		CookieJar:     cookieJar,
		RPSLimit:      u.Runner.RPSLimit,
		BPool:         u.BPool,
		Vu:            u.ID,
		Iteration:     u.Iteration,
	***REMOVED***

	newctx := common.WithRuntime(ctx, u.Runtime)
	newctx = common.WithState(newctx, state)
	*u.Context = newctx

	u.Runtime.Set("__ITER", u.Iteration)
	iter := u.Iteration
	u.Iteration++

	startTime := time.Now()
	v, err := fn(goja.Undefined(), args...) // Actually run the JS script
	endTime := time.Now()

	tags := state.Options.RunTags.CloneTags()
	if state.Options.SystemTags["vu"] ***REMOVED***
		tags["vu"] = strconv.FormatInt(u.ID, 10)
	***REMOVED***
	if state.Options.SystemTags["iter"] ***REMOVED***
		tags["iter"] = strconv.FormatInt(iter, 10)
	***REMOVED***
	sampleTags := stats.IntoSampleTags(&tags)

	if u.Runner.Bundle.Options.NoVUConnectionReuse.Bool ***REMOVED***
		u.HTTPTransport.CloseIdleConnections()
	***REMOVED***

	state.Samples = append(state.Samples, u.Dialer.GetTrail(startTime, endTime, sampleTags))

	return v, state, err
***REMOVED***
