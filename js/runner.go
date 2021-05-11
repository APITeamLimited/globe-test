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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/dop251/goja"
	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/net/http2"
	"golang.org/x/time/rate"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/lib/netext"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/stats"
)

// Ensure Runner implements the lib.Runner interface
var _ lib.Runner = &Runner***REMOVED******REMOVED***

type Runner struct ***REMOVED***
	Bundle       *Bundle
	Logger       *logrus.Logger
	defaultGroup *lib.Group

	BaseDialer net.Dialer
	Resolver   netext.Resolver
	// TODO: Remove ActualResolver, it's a hack to simplify mocking in tests.
	ActualResolver netext.MultiResolver
	RPSLimit       *rate.Limiter

	console   *console
	setupData []byte
***REMOVED***

// New returns a new Runner for the provide source
func New(
	logger *logrus.Logger, src *loader.SourceData, filesystems map[string]afero.Fs, rtOpts lib.RuntimeOptions,
) (*Runner, error) ***REMOVED***
	bundle, err := NewBundle(logger, src, filesystems, rtOpts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return newFromBundle(logger, bundle)
***REMOVED***

// NewFromArchive returns a new Runner from the source in the provided archive
func NewFromArchive(logger *logrus.Logger, arc *lib.Archive, rtOpts lib.RuntimeOptions) (*Runner, error) ***REMOVED***
	bundle, err := NewBundleFromArchive(logger, arc, rtOpts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return newFromBundle(logger, bundle)
***REMOVED***

func newFromBundle(logger *logrus.Logger, b *Bundle) (*Runner, error) ***REMOVED***
	defaultGroup, err := lib.NewGroup("", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defDNS := types.DefaultDNSConfig()
	r := &Runner***REMOVED***
		Bundle:       b,
		Logger:       logger,
		defaultGroup: defaultGroup,
		BaseDialer: net.Dialer***REMOVED***
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		***REMOVED***,
		console: newConsole(logger),
		Resolver: netext.NewResolver(
			net.LookupIP, 0, defDNS.Select.DNSSelect, defDNS.Policy.DNSPolicy),
		ActualResolver: net.LookupIP,
	***REMOVED***

	err = r.SetOptions(r.Bundle.Options)

	return r, err
***REMOVED***

func (r *Runner) MakeArchive() *lib.Archive ***REMOVED***
	return r.Bundle.makeArchive()
***REMOVED***

// NewVU returns a new initialized VU.
func (r *Runner) NewVU(id int64, samplesOut chan<- stats.SampleContainer) (lib.InitializedVU, error) ***REMOVED***
	vu, err := r.newVU(id, samplesOut)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return lib.InitializedVU(vu), nil
***REMOVED***

// nolint:funlen
func (r *Runner) newVU(id int64, samplesOut chan<- stats.SampleContainer) (*VU, error) ***REMOVED***
	// Instantiate a new bundle, make a VU out of it.
	bi, err := r.Bundle.Instantiate(r.Logger, id)
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
		Dialer:           r.BaseDialer,
		Resolver:         r.Resolver,
		Blacklist:        r.Bundle.Options.BlacklistIPs,
		BlockedHostnames: r.Bundle.Options.BlockedHostnames.Trie,
		Hosts:            r.Bundle.Options.Hosts,
	***REMOVED***
	if r.Bundle.Options.LocalIPs.Valid ***REMOVED***
		var ipIndex uint64
		if id > 0 ***REMOVED***
			ipIndex = uint64(id - 1)
		***REMOVED***
		dialer.Dialer.LocalAddr = &net.TCPAddr***REMOVED***IP: r.Bundle.Options.LocalIPs.Pool.GetIP(ipIndex)***REMOVED***
	***REMOVED***

	tlsConfig := &tls.Config***REMOVED***
		InsecureSkipVerify: r.Bundle.Options.InsecureSkipTLSVerify.Bool, //nolint:gosec
		CipherSuites:       cipherSuites,
		MinVersion:         uint16(tlsVersions.Min),
		MaxVersion:         uint16(tlsVersions.Max),
		Certificates:       certs,
		NameToCertificate:  nameToCert,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	***REMOVED***
	transport := &http.Transport***REMOVED***
		Proxy:               http.ProxyFromEnvironment,
		TLSClientConfig:     tlsConfig,
		DialContext:         dialer.DialContext,
		DisableCompression:  true,
		DisableKeepAlives:   r.Bundle.Options.NoConnectionReuse.Bool,
		MaxIdleConns:        int(r.Bundle.Options.Batch.Int64),
		MaxIdleConnsPerHost: int(r.Bundle.Options.BatchPerHost.Int64),
	***REMOVED***
	_ = http2.ConfigureTransport(transport)

	cookieJar, err := cookiejar.New(nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vu := &VU***REMOVED***
		ID:             id,
		BundleInstance: *bi,
		Runner:         r,
		Transport:      transport,
		Dialer:         dialer,
		CookieJar:      cookieJar,
		TLSConfig:      tlsConfig,
		Console:        r.console,
		BPool:          bpool.NewBufferPool(100),
		Samples:        samplesOut,
	***REMOVED***

	vu.state = &lib.State***REMOVED***
		Logger:    vu.Runner.Logger,
		Options:   vu.Runner.Bundle.Options,
		Transport: vu.Transport,
		Dialer:    vu.Dialer,
		TLSConfig: vu.TLSConfig,
		CookieJar: cookieJar,
		RPSLimit:  vu.Runner.RPSLimit,
		BPool:     vu.BPool,
		Vu:        vu.ID,
		Samples:   vu.Samples,
		Iteration: vu.Iteration,
		Tags:      vu.Runner.Bundle.Options.RunTags.CloneTags(),
		Group:     r.defaultGroup,
	***REMOVED***
	vu.Runtime.Set("console", common.Bind(vu.Runtime, vu.Console, vu.Context))

	// This is here mostly so if someone tries they get a nice message
	// instead of "Value is not an object: undefined  ..."
	common.BindToGlobal(vu.Runtime, map[string]interface***REMOVED******REMOVED******REMOVED***
		"open": func() ***REMOVED***
			common.Throw(vu.Runtime, errors.New(openCantBeUsedOutsideInitContextMsg))
		***REMOVED***,
	***REMOVED***)

	return vu, nil
***REMOVED***

// Setup runs the setup function if there is one and sets the setupData to the returned value
func (r *Runner) Setup(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
	setupCtx, setupCancel := context.WithTimeout(ctx, r.getTimeoutFor(consts.SetupFn))
	defer setupCancel()

	v, err := r.runPart(setupCtx, out, consts.SetupFn, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// r.setupData = nil is special it means undefined from this moment forward
	if goja.IsUndefined(v) ***REMOVED***
		r.setupData = nil
		return nil
	***REMOVED***

	r.setupData, err = json.Marshal(v.Export())
	if err != nil ***REMOVED***
		return fmt.Errorf("error marshaling setup() data to JSON: %w", err)
	***REMOVED***
	var tmp interface***REMOVED******REMOVED***
	return json.Unmarshal(r.setupData, &tmp)
***REMOVED***

// GetSetupData returns the setup data as json if Setup() was specified and executed, nil otherwise
func (r *Runner) GetSetupData() []byte ***REMOVED***
	return r.setupData
***REMOVED***

// SetSetupData saves the externally supplied setup data as json in the runner, so it can be used in VUs
func (r *Runner) SetSetupData(data []byte) ***REMOVED***
	r.setupData = data
***REMOVED***

func (r *Runner) Teardown(ctx context.Context, out chan<- stats.SampleContainer) error ***REMOVED***
	teardownCtx, teardownCancel := context.WithTimeout(ctx, r.getTimeoutFor(consts.TeardownFn))
	defer teardownCancel()

	var data interface***REMOVED******REMOVED***
	if r.setupData != nil ***REMOVED***
		if err := json.Unmarshal(r.setupData, &data); err != nil ***REMOVED***
			return fmt.Errorf("error unmarshaling setup data for teardown() from JSON: %w", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		data = goja.Undefined()
	***REMOVED***
	_, err := r.runPart(teardownCtx, out, consts.TeardownFn, data)
	return err
***REMOVED***

func (r *Runner) GetDefaultGroup() *lib.Group ***REMOVED***
	return r.defaultGroup
***REMOVED***

func (r *Runner) GetOptions() lib.Options ***REMOVED***
	return r.Bundle.Options
***REMOVED***

// IsExecutable returns whether the given name is an exported and
// executable function in the script.
func (r *Runner) IsExecutable(name string) bool ***REMOVED***
	_, exists := r.Bundle.exports[name]
	return exists
***REMOVED***

// HandleSummary calls the specified summary callback, if supplied.
func (r *Runner) HandleSummary(ctx context.Context, summary *lib.Summary) (map[string]io.Reader, error) ***REMOVED***
	summaryDataForJS := summarizeMetricsToObject(summary, r.Bundle.Options)

	out := make(chan stats.SampleContainer, 100)
	defer close(out)

	go func() ***REMOVED*** // discard all metrics
		for range out ***REMOVED***
		***REMOVED***
	***REMOVED***()

	vu, err := r.newVU(0, out)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	handleSummaryFn := goja.Undefined()
	if exported := vu.Runtime.Get("exports").ToObject(vu.Runtime); exported != nil ***REMOVED***
		fn := exported.Get(consts.HandleSummaryFn)
		if _, ok := goja.AssertFunction(fn); ok ***REMOVED***
			handleSummaryFn = fn
		***REMOVED*** else if fn != nil ***REMOVED***
			return nil, fmt.Errorf("exported identfier %s must be a function", consts.HandleSummaryFn)
		***REMOVED***
	***REMOVED***
	ctx = common.WithRuntime(ctx, vu.Runtime)
	ctx = lib.WithState(ctx, vu.state)
	ctx, cancel := context.WithTimeout(ctx, r.getTimeoutFor(consts.HandleSummaryFn))
	defer cancel()
	go func() ***REMOVED***
		<-ctx.Done()
		vu.Runtime.Interrupt(context.Canceled)
	***REMOVED***()
	*vu.Context = ctx

	handleSummaryWrapperRaw, err := vu.Runtime.RunString(summaryWrapperLambdaCode)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unexpected error while getting the summary wrapper: %w", err)
	***REMOVED***
	handleSummaryWrapper, ok := goja.AssertFunction(handleSummaryWrapperRaw)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("unexpected error did not get a callable summary wrapper")
	***REMOVED***

	wrapperArgs := []goja.Value***REMOVED***
		handleSummaryFn,
		vu.Runtime.ToValue(r.Bundle.RuntimeOptions.SummaryExport.String),
		vu.Runtime.ToValue(summaryDataForJS),
		vu.Runtime.ToValue(getOldTextSummaryFunc(summary, r.Bundle.Options)), // TODO: remove
	***REMOVED***
	rawResult, _, _, err := vu.runFn(ctx, false, handleSummaryWrapper, wrapperArgs...)

	// TODO: refactor the whole JS runner to avoid copy-pasting these complicated bits...
	// deadline is reached so we have timeouted but this might've not been registered correctly
	if deadline, ok := ctx.Deadline(); ok && time.Now().After(deadline) ***REMOVED***
		// we could have an error that is not context.Canceled in which case we should return it instead
		if err, ok := err.(*goja.InterruptedError); ok && rawResult != nil && err.Value() != context.Canceled ***REMOVED***
			// TODO: silence this error?
			return nil, err
		***REMOVED***
		// otherwise we have timeouted
		return nil, lib.NewTimeoutError(consts.HandleSummaryFn, r.getTimeoutFor(consts.HandleSummaryFn))
	***REMOVED***

	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unexpected error while generating the summary: %w", err)
	***REMOVED***
	return getSummaryResult(rawResult)
***REMOVED***

func (r *Runner) SetOptions(opts lib.Options) error ***REMOVED***
	r.Bundle.Options = opts
	r.RPSLimit = nil
	if rps := opts.RPS; rps.Valid ***REMOVED***
		r.RPSLimit = rate.NewLimiter(rate.Limit(rps.Int64), 1)
	***REMOVED***

	// TODO: validate that all exec values are either nil or valid exported methods (or HTTP requests in the future)

	if opts.ConsoleOutput.Valid ***REMOVED***
		c, err := newFileConsole(opts.ConsoleOutput.String, r.Logger.Formatter)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		r.console = c
	***REMOVED***

	// FIXME: Resolver probably shouldn't be reset here...
	// It's done because the js.Runner is created before the full
	// configuration has been processed, at which point we don't have
	// access to the DNSConfig, and need to wait for this SetOptions
	// call that happens after all config has been assembled.
	// We could make DNSConfig part of RuntimeOptions, but that seems
	// conceptually wrong since the JS runtime doesn't care about it
	// (it needs the actual resolver, not the config), and it would
	// require an additional field on Bundle to pass the config through,
	// which is arguably worse than this.
	if err := r.setResolver(opts.DNS); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (r *Runner) setResolver(dns types.DNSConfig) error ***REMOVED***
	ttl, err := parseTTL(dns.TTL.String)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	dnsSel := dns.Select
	if !dnsSel.Valid ***REMOVED***
		dnsSel = types.DefaultDNSConfig().Select
	***REMOVED***
	dnsPol := dns.Policy
	if !dnsPol.Valid ***REMOVED***
		dnsPol = types.DefaultDNSConfig().Policy
	***REMOVED***
	r.Resolver = netext.NewResolver(
		r.ActualResolver, ttl, dnsSel.DNSSelect, dnsPol.DNSPolicy)

	return nil
***REMOVED***

func parseTTL(ttlS string) (time.Duration, error) ***REMOVED***
	ttl := time.Duration(0)
	switch ttlS ***REMOVED***
	case "inf":
		// cache "infinitely"
		ttl = time.Hour * 24 * 365
	case "0":
		// disable cache
	case "":
		ttlS = types.DefaultDNSConfig().TTL.String
		fallthrough
	default:
		var err error
		ttl, err = types.ParseExtendedDuration(ttlS)
		if ttl < 0 || err != nil ***REMOVED***
			return ttl, fmt.Errorf("invalid DNS TTL: %s", ttlS)
		***REMOVED***
	***REMOVED***
	return ttl, nil
***REMOVED***

// Runs an exported function in its own temporary VU, optionally with an argument. Execution is
// interrupted if the context expires. No error is returned if the part does not exist.
func (r *Runner) runPart(ctx context.Context, out chan<- stats.SampleContainer, name string, arg interface***REMOVED******REMOVED***) (goja.Value, error) ***REMOVED***
	vu, err := r.newVU(0, out)
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

	ctx = common.WithRuntime(ctx, vu.Runtime)
	ctx = lib.WithState(ctx, vu.state)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() ***REMOVED***
		<-ctx.Done()
		vu.Runtime.Interrupt(context.Canceled)
	***REMOVED***()
	*vu.Context = ctx

	group, err := r.GetDefaultGroup().Group(name)
	if err != nil ***REMOVED***
		return goja.Undefined(), err
	***REMOVED***

	if r.Bundle.Options.SystemTags.Has(stats.TagGroup) ***REMOVED***
		vu.state.Tags["group"] = group.Path
	***REMOVED***
	vu.state.Group = group

	v, _, _, err := vu.runFn(ctx, false, fn, vu.Runtime.ToValue(arg))

	// deadline is reached so we have timeouted but this might've not been registered correctly
	if deadline, ok := ctx.Deadline(); ok && time.Now().After(deadline) ***REMOVED***
		// we could have an error that is not context.Canceled in which case we should return it instead
		if err, ok := err.(*goja.InterruptedError); ok && v != nil && err.Value() != context.Canceled ***REMOVED***
			// TODO: silence this error?
			return v, err
		***REMOVED***
		// otherwise we have timeouted
		return v, lib.NewTimeoutError(name, r.getTimeoutFor(name))
	***REMOVED***
	return v, err
***REMOVED***

// getTimeoutFor returns the timeout duration for given special script function.
func (r *Runner) getTimeoutFor(stage string) time.Duration ***REMOVED***
	d := time.Duration(0)
	switch stage ***REMOVED***
	case consts.SetupFn:
		return time.Duration(r.Bundle.Options.SetupTimeout.Duration)
	case consts.TeardownFn:
		return time.Duration(r.Bundle.Options.TeardownTimeout.Duration)
	case consts.HandleSummaryFn:
		return 2 * time.Minute // TODO: make configurable
	***REMOVED***
	return d
***REMOVED***

type VU struct ***REMOVED***
	BundleInstance

	Runner    *Runner
	Transport *http.Transport
	Dialer    *netext.Dialer
	CookieJar *cookiejar.Jar
	TLSConfig *tls.Config
	ID        int64
	Iteration int64

	Console *console
	BPool   *bpool.BufferPool

	Samples chan<- stats.SampleContainer

	setupData goja.Value

	state *lib.State
***REMOVED***

// Verify that interfaces are implemented
var (
	_ lib.ActiveVU      = &ActiveVU***REMOVED******REMOVED***
	_ lib.InitializedVU = &VU***REMOVED******REMOVED***
)

// ActiveVU holds a VU and its activation parameters
type ActiveVU struct ***REMOVED***
	*VU
	*lib.VUActivationParams
	busy chan struct***REMOVED******REMOVED***
***REMOVED***

// GetID returns the unique VU ID.
func (u *VU) GetID() int64 ***REMOVED***
	return u.ID
***REMOVED***

// Activate the VU so it will be able to run code.
func (u *VU) Activate(params *lib.VUActivationParams) lib.ActiveVU ***REMOVED***
	u.Runtime.ClearInterrupt()

	if params.Exec == "" ***REMOVED***
		params.Exec = consts.DefaultFn
	***REMOVED***

	// Override the preset global env with any custom env vars
	env := make(map[string]string, len(u.env)+len(params.Env))
	for key, value := range u.env ***REMOVED***
		env[key] = value
	***REMOVED***
	for key, value := range params.Env ***REMOVED***
		env[key] = value
	***REMOVED***
	u.Runtime.Set("__ENV", env)

	opts := u.Runner.Bundle.Options
	// TODO: maybe we can cache the original tags only clone them and add (if any) new tags on top ?
	u.state.Tags = opts.RunTags.CloneTags()
	for k, v := range params.Tags ***REMOVED***
		u.state.Tags[k] = v
	***REMOVED***
	if opts.SystemTags.Has(stats.TagVU) ***REMOVED***
		u.state.Tags["vu"] = strconv.FormatInt(u.ID, 10)
	***REMOVED***
	if opts.SystemTags.Has(stats.TagIter) ***REMOVED***
		u.state.Tags["iter"] = strconv.FormatInt(u.Iteration, 10)
	***REMOVED***
	if opts.SystemTags.Has(stats.TagGroup) ***REMOVED***
		u.state.Tags["group"] = u.state.Group.Path
	***REMOVED***
	if opts.SystemTags.Has(stats.TagScenario) ***REMOVED***
		u.state.Tags["scenario"] = params.Scenario
	***REMOVED***

	params.RunContext = common.WithRuntime(params.RunContext, u.Runtime)
	params.RunContext = lib.WithState(params.RunContext, u.state)
	*u.Context = params.RunContext

	avu := &ActiveVU***REMOVED***
		VU:                 u,
		VUActivationParams: params,
		busy:               make(chan struct***REMOVED******REMOVED***, 1),
	***REMOVED***

	go func() ***REMOVED***
		// Wait for the run context to be over
		<-params.RunContext.Done()
		// Interrupt the JS runtime
		u.Runtime.Interrupt(context.Canceled)
		// Wait for the VU to stop running, if it was, and prevent it from
		// running again for this activation
		avu.busy <- struct***REMOVED******REMOVED******REMOVED******REMOVED***

		if params.DeactivateCallback != nil ***REMOVED***
			params.DeactivateCallback(u)
		***REMOVED***
	***REMOVED***()

	return avu
***REMOVED***

// RunOnce runs the configured Exec function once.
func (u *ActiveVU) RunOnce() error ***REMOVED***
	select ***REMOVED***
	case <-u.RunContext.Done():
		return u.RunContext.Err() // we are done, return
	case u.busy <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
		// nothing else can run now, and the VU cannot be deactivated
	***REMOVED***
	defer func() ***REMOVED***
		<-u.busy // unlock deactivation again
	***REMOVED***()

	// Unmarshall the setupData only the first time for each VU so that VUs are isolated but we
	// still don't use too much CPU in the middle test
	if u.setupData == nil ***REMOVED***
		if u.Runner.setupData != nil ***REMOVED***
			var data interface***REMOVED******REMOVED***
			if err := json.Unmarshal(u.Runner.setupData, &data); err != nil ***REMOVED***
				return fmt.Errorf("error unmarshaling setup data for the iteration from JSON: %w", err)
			***REMOVED***
			u.setupData = u.Runtime.ToValue(data)
		***REMOVED*** else ***REMOVED***
			u.setupData = goja.Undefined()
		***REMOVED***
	***REMOVED***

	fn, ok := u.exports[u.Exec]
	if !ok ***REMOVED***
		// Shouldn't happen; this is validated in cmd.validateScenarioConfig()
		panic(fmt.Sprintf("function '%s' not found in exports", u.Exec))
	***REMOVED***

	// Call the exported function.
	_, isFullIteration, totalTime, err := u.runFn(u.RunContext, true, fn, u.setupData)

	// If MinIterationDuration is specified and the iteration wasn't canceled
	// and was less than it, sleep for the remainder
	if isFullIteration && u.Runner.Bundle.Options.MinIterationDuration.Valid ***REMOVED***
		durationDiff := time.Duration(u.Runner.Bundle.Options.MinIterationDuration.Duration) - totalTime
		if durationDiff > 0 ***REMOVED***
			time.Sleep(durationDiff)
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

func (u *VU) runFn(
	ctx context.Context, isDefault bool, fn goja.Callable, args ...goja.Value,
) (v goja.Value, isFullIteration bool, t time.Duration, err error) ***REMOVED***
	if !u.Runner.Bundle.Options.NoCookiesReset.ValueOrZero() ***REMOVED***
		u.state.CookieJar, err = cookiejar.New(nil)
		if err != nil ***REMOVED***
			return goja.Undefined(), false, time.Duration(0), err
		***REMOVED***
	***REMOVED***

	opts := &u.Runner.Bundle.Options
	if opts.SystemTags.Has(stats.TagIter) ***REMOVED***
		u.state.Tags["iter"] = strconv.FormatInt(u.Iteration, 10)
	***REMOVED***

	// TODO: this seems like the wrong place for the iteration incrementation
	// also this means that teardown and setup have __ITER defined
	// maybe move it to RunOnce ?
	u.Runtime.Set("__ITER", u.Iteration)
	u.Iteration++

	defer func() ***REMOVED***
		if r := recover(); r != nil ***REMOVED***
			gojaStack := u.Runtime.CaptureCallStack(20, nil)
			err = fmt.Errorf("a panic occurred in VU code but was caught: %s", r)
			// TODO figure out how to use PanicLevel without panicing .. this might require changing
			// the logger we use see
			// https://github.com/sirupsen/logrus/issues/1028
			// https://github.com/sirupsen/logrus/issues/993
			b := new(bytes.Buffer)
			for _, s := range gojaStack ***REMOVED***
				s.Write(b)
			***REMOVED***
			u.state.Logger.Log(logrus.ErrorLevel, "panic: ", r, "\n", string(debug.Stack()), "\nGoja stack:\n", b.String())
		***REMOVED***
	***REMOVED***()

	startTime := time.Now()
	v, err = fn(goja.Undefined(), args...) // Actually run the JS script
	endTime := time.Now()
	var exception *goja.Exception
	if errors.As(err, &exception) ***REMOVED***
		err = &scriptException***REMOVED***inner: exception***REMOVED***
	***REMOVED***

	select ***REMOVED***
	case <-ctx.Done():
		isFullIteration = false
	default:
		isFullIteration = true
	***REMOVED***

	if u.Runner.Bundle.Options.NoVUConnectionReuse.Bool ***REMOVED***
		u.Transport.CloseIdleConnections()
	***REMOVED***

	u.state.Samples <- u.Dialer.GetTrail(startTime, endTime, isFullIteration, isDefault, stats.NewSampleTags(u.state.Tags))

	return v, isFullIteration, endTime.Sub(startTime), err
***REMOVED***

type scriptException struct ***REMOVED***
	inner *goja.Exception
***REMOVED***

var _ types.ScriptException = &scriptException***REMOVED******REMOVED***

func (s *scriptException) Error() string ***REMOVED***
	// this calls String instead of error so that by default if it's printed to print the stacktrace
	return s.inner.String()
***REMOVED***

func (s *scriptException) StackTrace() string ***REMOVED***
	return s.inner.String()
***REMOVED***

func (s *scriptException) Unwrap() error ***REMOVED***
	return s.inner
***REMOVED***
