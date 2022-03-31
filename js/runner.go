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
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/oxtoacart/bpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/net/http2"
	"golang.org/x/time/rate"

	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/eventloop"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/consts"
	"go.k6.io/k6/lib/netext"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
)

// Ensure Runner implements the lib.Runner interface
var _ lib.Runner = &Runner***REMOVED******REMOVED***

// TODO: https://github.com/grafana/k6/issues/2186
// An advanced TLS support should cover the rid of the warning
//
// nolint:gochecknoglobals
var nameToCertWarning sync.Once

type Runner struct ***REMOVED***
	Bundle         *Bundle
	Logger         *logrus.Logger
	defaultGroup   *lib.Group
	builtinMetrics *metrics.BuiltinMetrics
	registry       *metrics.Registry

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
	rs *lib.RuntimeState, src *loader.SourceData, filesystems map[string]afero.Fs,
) (*Runner, error) ***REMOVED***
	bundle, err := NewBundle(rs.Logger, src, filesystems, rs.RuntimeOptions, rs.Registry)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewFromBundle(rs, bundle)
***REMOVED***

// NewFromArchive returns a new Runner from the source in the provided archive
func NewFromArchive(rs *lib.RuntimeState, arc *lib.Archive) (*Runner, error) ***REMOVED***
	bundle, err := NewBundleFromArchive(rs.Logger, arc, rs.RuntimeOptions, rs.Registry)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewFromBundle(rs, bundle)
***REMOVED***

// NewFromBundle returns a new Runner from the provided Bundle
func NewFromBundle(rs *lib.RuntimeState, b *Bundle) (*Runner, error) ***REMOVED***
	defaultGroup, err := lib.NewGroup("", nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defDNS := types.DefaultDNSConfig()
	r := &Runner***REMOVED***
		Bundle:       b,
		Logger:       rs.Logger,
		defaultGroup: defaultGroup,
		BaseDialer: net.Dialer***REMOVED***
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		***REMOVED***,
		console: newConsole(rs.Logger),
		Resolver: netext.NewResolver(
			net.LookupIP, 0, defDNS.Select.DNSSelect, defDNS.Policy.DNSPolicy),
		ActualResolver: net.LookupIP,
		builtinMetrics: rs.BuiltinMetrics,
		registry:       rs.Registry,
	***REMOVED***

	err = r.SetOptions(r.Bundle.Options)

	return r, err
***REMOVED***

func (r *Runner) MakeArchive() *lib.Archive ***REMOVED***
	return r.Bundle.makeArchive()
***REMOVED***

// NewVU returns a new initialized VU.
func (r *Runner) NewVU(idLocal, idGlobal uint64, samplesOut chan<- metrics.SampleContainer) (lib.InitializedVU, error) ***REMOVED***
	vu, err := r.newVU(idLocal, idGlobal, samplesOut)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return lib.InitializedVU(vu), nil
***REMOVED***

// nolint:funlen
func (r *Runner) newVU(idLocal, idGlobal uint64, samplesOut chan<- metrics.SampleContainer) (*VU, error) ***REMOVED***
	// Instantiate a new bundle, make a VU out of it.
	moduleVUImpl := newModuleVUImpl()
	bi, err := r.Bundle.Instantiate(r.Logger, idLocal, moduleVUImpl)
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
		cert, errC := auth.Certificate()
		if errC != nil ***REMOVED***
			return nil, errC
		***REMOVED***
		certs[i] = *cert
		for _, name := range auth.Domains ***REMOVED***
			nameToCert[name] = cert
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
		if idLocal > 0 ***REMOVED***
			ipIndex = idLocal - 1
		***REMOVED***
		dialer.Dialer.LocalAddr = &net.TCPAddr***REMOVED***IP: r.Bundle.Options.LocalIPs.Pool.GetIP(ipIndex)***REMOVED***
	***REMOVED***

	tlsConfig := &tls.Config***REMOVED***
		InsecureSkipVerify: r.Bundle.Options.InsecureSkipTLSVerify.Bool, //nolint:gosec
		CipherSuites:       cipherSuites,
		MinVersion:         uint16(tlsVersions.Min),
		MaxVersion:         uint16(tlsVersions.Max),
		Certificates:       certs,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	***REMOVED***
	// Follow NameToCertificate in https://pkg.go.dev/crypto/tls@go1.17.6#Config, leave this field nil
	// when it is empty
	if len(nameToCert) > 0 ***REMOVED***
		nameToCertWarning.Do(func() ***REMOVED***
			r.Logger.Warn("tlsAuth.domains option could be removed in the next releases, it's recommended to leave it empty " +
				"and let k6 automatically detect from the provided certificate. It follows the Go's NameToCertificate " +
				"deprecation - https://pkg.go.dev/crypto/tls@go1.17#Config.")
		***REMOVED***)
		// nolint:staticcheck // ignore SA1019 we can deprecate it but we have to continue to support the previous code.
		tlsConfig.NameToCertificate = nameToCert
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

	if forceHTTP1() ***REMOVED***
		transport.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper) // send over h1 protocol
	***REMOVED*** else ***REMOVED***
		_ = http2.ConfigureTransport(transport) // send over h2 protocol
	***REMOVED***

	cookieJar, err := cookiejar.New(nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vu := &VU***REMOVED***
		ID:             idLocal,
		IDGlobal:       idGlobal,
		iteration:      int64(-1),
		BundleInstance: *bi,
		Runner:         r,
		Transport:      transport,
		Dialer:         dialer,
		CookieJar:      cookieJar,
		TLSConfig:      tlsConfig,
		Console:        r.console,
		BPool:          bpool.NewBufferPool(100),
		Samples:        samplesOut,
		scenarioIter:   make(map[string]uint64),
		moduleVUImpl:   moduleVUImpl,
	***REMOVED***

	vu.state = &lib.State***REMOVED***
		Logger:         vu.Runner.Logger,
		Options:        vu.Runner.Bundle.Options,
		Transport:      vu.Transport,
		Dialer:         vu.Dialer,
		TLSConfig:      vu.TLSConfig,
		CookieJar:      cookieJar,
		RPSLimit:       vu.Runner.RPSLimit,
		BPool:          vu.BPool,
		VUID:           vu.ID,
		VUIDGlobal:     vu.IDGlobal,
		Samples:        vu.Samples,
		Tags:           lib.NewTagMap(vu.Runner.Bundle.Options.RunTags.CloneTags()),
		Group:          r.defaultGroup,
		BuiltinMetrics: r.builtinMetrics,
	***REMOVED***
	vu.moduleVUImpl.state = vu.state
	_ = vu.Runtime.Set("console", vu.Console)

	// This is here mostly so if someone tries they get a nice message
	// instead of "Value is not an object: undefined  ..."
	_ = vu.Runtime.GlobalObject().Set("open",
		func() ***REMOVED***
			common.Throw(vu.Runtime, errors.New(openCantBeUsedOutsideInitContextMsg))
		***REMOVED***)

	return vu, nil
***REMOVED***

// forceHTTP1 checks if force http1 env variable has been set in order to force requests to be sent over h1
// TODO: This feature is temporary until #936 is resolved
func forceHTTP1() bool ***REMOVED***
	godebug := os.Getenv("GODEBUG")
	if godebug == "" ***REMOVED***
		return false
	***REMOVED***
	variables := strings.SplitAfter(godebug, ",")

	for _, v := range variables ***REMOVED***
		if strings.Trim(v, ",") == "http2client=0" ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Setup runs the setup function if there is one and sets the setupData to the returned value
func (r *Runner) Setup(ctx context.Context, out chan<- metrics.SampleContainer) error ***REMOVED***
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

// Teardown runs the teardown function if there is one.
func (r *Runner) Teardown(ctx context.Context, out chan<- metrics.SampleContainer) error ***REMOVED***
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
	summaryDataForJS := summarizeMetricsToObject(summary, r.Bundle.Options, r.setupData)

	out := make(chan metrics.SampleContainer, 100)
	defer close(out)

	go func() ***REMOVED*** // discard all metrics
		for range out ***REMOVED***
		***REMOVED***
	***REMOVED***()

	vu, err := r.newVU(0, 0, out)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	handleSummaryFn := goja.Undefined()
	if exported := vu.Runtime.Get("exports").ToObject(vu.Runtime); exported != nil ***REMOVED***
		fn := exported.Get(consts.HandleSummaryFn)
		if _, ok := goja.AssertFunction(fn); ok ***REMOVED***
			handleSummaryFn = fn
		***REMOVED*** else if fn != nil ***REMOVED***
			return nil, fmt.Errorf("exported identifier %s must be a function", consts.HandleSummaryFn)
		***REMOVED***
	***REMOVED***

	ctx, cancel := context.WithTimeout(ctx, r.getTimeoutFor(consts.HandleSummaryFn))
	defer cancel()
	go func() ***REMOVED***
		<-ctx.Done()
		vu.Runtime.Interrupt(context.Canceled)
	***REMOVED***()
	*vu.Context = ctx

	wrapper := strings.Replace(summaryWrapperLambdaCode, "/*JSLIB_SUMMARY_CODE*/", jslibSummaryCode, 1)
	handleSummaryWrapperRaw, err := vu.Runtime.RunString(wrapper)
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
	***REMOVED***
	rawResult, _, _, err := vu.runFn(ctx, false, handleSummaryWrapper, nil, wrapperArgs...)

	// TODO: refactor the whole JS runner to avoid copy-pasting these complicated bits...
	// deadline is reached so we have timeouted but this might've not been registered correctly
	if deadline, ok := ctx.Deadline(); ok && time.Now().After(deadline) ***REMOVED***
		// we could have an error that is not context.Canceled in which case we should return it instead
		if err, ok := err.(*goja.InterruptedError); ok && rawResult != nil && err.Value() != context.Canceled ***REMOVED***
			// TODO: silence this error?
			return nil, err
		***REMOVED***
		// otherwise we have timeouted
		return nil, newTimeoutError(consts.HandleSummaryFn, r.getTimeoutFor(consts.HandleSummaryFn))
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
func (r *Runner) runPart(
	ctx context.Context,
	out chan<- metrics.SampleContainer,
	name string,
	arg interface***REMOVED******REMOVED***,
) (goja.Value, error) ***REMOVED***
	vu, err := r.newVU(0, 0, out)
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

	if r.Bundle.Options.SystemTags.Has(metrics.TagGroup) ***REMOVED***
		vu.state.Tags.Set("group", group.Path)
	***REMOVED***
	vu.state.Group = group

	v, _, _, err := vu.runFn(ctx, false, fn, nil, vu.Runtime.ToValue(arg))

	// deadline is reached so we have timeouted but this might've not been registered correctly
	if deadline, ok := ctx.Deadline(); ok && time.Now().After(deadline) ***REMOVED***
		// we could have an error that is not context.Canceled in which case we should return it instead
		if err, ok := err.(*goja.InterruptedError); ok && v != nil && err.Value() != context.Canceled ***REMOVED***
			// TODO: silence this error?
			return v, err
		***REMOVED***
		// otherwise we have timeouted
		return v, newTimeoutError(name, r.getTimeoutFor(name))
	***REMOVED***
	return v, err
***REMOVED***

// getTimeoutFor returns the timeout duration for given special script function.
func (r *Runner) getTimeoutFor(stage string) time.Duration ***REMOVED***
	d := time.Duration(0)
	switch stage ***REMOVED***
	case consts.SetupFn:
		return r.Bundle.Options.SetupTimeout.TimeDuration()
	case consts.TeardownFn:
		return r.Bundle.Options.TeardownTimeout.TimeDuration()
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
	ID        uint64 // local to the current instance
	IDGlobal  uint64 // global across all instances
	iteration int64

	Console *console
	BPool   *bpool.BufferPool

	Samples chan<- metrics.SampleContainer

	setupData goja.Value

	state *lib.State
	// count of iterations executed by this VU in each scenario
	scenarioIter map[string]uint64

	moduleVUImpl *moduleVUImpl
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

	scenarioName              string
	getNextIterationCounters  func() (uint64, uint64)
	scIterLocal, scIterGlobal uint64
***REMOVED***

// GetID returns the unique VU ID.
func (u *VU) GetID() uint64 ***REMOVED***
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
	u.state.Tags = lib.NewTagMap(opts.RunTags.CloneTags())
	for k, v := range params.Tags ***REMOVED***
		u.state.Tags.Set(k, v)
	***REMOVED***
	if opts.SystemTags.Has(metrics.TagVU) ***REMOVED***
		u.state.Tags.Set("vu", strconv.FormatUint(u.ID, 10))
	***REMOVED***
	if opts.SystemTags.Has(metrics.TagIter) ***REMOVED***
		u.state.Tags.Set("iter", strconv.FormatInt(u.iteration, 10))
	***REMOVED***
	if opts.SystemTags.Has(metrics.TagGroup) ***REMOVED***
		u.state.Tags.Set("group", u.state.Group.Path)
	***REMOVED***
	if opts.SystemTags.Has(metrics.TagScenario) ***REMOVED***
		u.state.Tags.Set("scenario", params.Scenario)
	***REMOVED***

	ctx := params.RunContext
	*u.Context = ctx

	u.state.GetScenarioVUIter = func() uint64 ***REMOVED***
		return u.scenarioIter[params.Scenario]
	***REMOVED***

	avu := &ActiveVU***REMOVED***
		VU:                       u,
		VUActivationParams:       params,
		busy:                     make(chan struct***REMOVED******REMOVED***, 1),
		scenarioName:             params.Scenario,
		scIterLocal:              ^uint64(0),
		scIterGlobal:             ^uint64(0),
		getNextIterationCounters: params.GetNextIterationCounters,
	***REMOVED***

	u.state.GetScenarioLocalVUIter = func() uint64 ***REMOVED***
		return avu.scIterLocal
	***REMOVED***
	u.state.GetScenarioGlobalVUIter = func() uint64 ***REMOVED***
		return avu.scIterGlobal
	***REMOVED***

	go func() ***REMOVED***
		// Wait for the run context to be over
		<-ctx.Done()
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

	u.incrIteration()
	if err := u.Runtime.Set("__ITER", u.iteration); err != nil ***REMOVED***
		panic(fmt.Errorf("error setting __ITER in goja runtime: %w", err))
	***REMOVED***

	ctx, cancel := context.WithCancel(u.RunContext)
	defer cancel()
	u.moduleVUImpl.ctx = ctx
	// Call the exported function.
	_, isFullIteration, totalTime, err := u.runFn(ctx, true, fn, cancel, u.setupData)
	if err != nil ***REMOVED***
		var x *goja.InterruptedError
		if errors.As(err, &x) ***REMOVED***
			if v, ok := x.Value().(*common.InterruptError); ok ***REMOVED***
				v.Reason = x.Error()
				err = v
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// If MinIterationDuration is specified and the iteration wasn't canceled
	// and was less than it, sleep for the remainder
	if isFullIteration && u.Runner.Bundle.Options.MinIterationDuration.Valid ***REMOVED***
		durationDiff := u.Runner.Bundle.Options.MinIterationDuration.TimeDuration() - totalTime
		if durationDiff > 0 ***REMOVED***
			select ***REMOVED***
			case <-time.After(durationDiff):
			case <-u.RunContext.Done():
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return err
***REMOVED***

// if isDefault is true, cancel also needs to be provided and it should cancel the provided context
// TODO remove the need for the above through refactoring of this function and its callees
func (u *VU) runFn(
	ctx context.Context, isDefault bool, fn goja.Callable, cancel func(), args ...goja.Value,
) (v goja.Value, isFullIteration bool, t time.Duration, err error) ***REMOVED***
	if !u.Runner.Bundle.Options.NoCookiesReset.ValueOrZero() ***REMOVED***
		u.state.CookieJar, err = cookiejar.New(nil)
		if err != nil ***REMOVED***
			return goja.Undefined(), false, time.Duration(0), err
		***REMOVED***
	***REMOVED***

	opts := &u.Runner.Bundle.Options
	if opts.SystemTags.Has(metrics.TagIter) ***REMOVED***
		u.state.Tags.Set("iter", strconv.FormatInt(u.state.Iteration, 10))
	***REMOVED***

	startTime := time.Now()

	if u.moduleVUImpl.eventLoop == nil ***REMOVED***
		u.moduleVUImpl.eventLoop = eventloop.New(u.moduleVUImpl)
	***REMOVED***
	err = common.RunWithPanicCatching(u.state.Logger, u.Runtime, func() error ***REMOVED***
		return u.moduleVUImpl.eventLoop.Start(func() (err error) ***REMOVED***
			v, err = fn(goja.Undefined(), args...) // Actually run the JS script
			return err
		***REMOVED***)
	***REMOVED***)

	select ***REMOVED***
	case <-ctx.Done():
		isFullIteration = false
	default:
		isFullIteration = true
	***REMOVED***

	if cancel != nil ***REMOVED***
		cancel()
		u.moduleVUImpl.eventLoop.WaitOnRegistered()
	***REMOVED***
	endTime := time.Now()
	var exception *goja.Exception
	if errors.As(err, &exception) ***REMOVED***
		err = &scriptException***REMOVED***inner: exception***REMOVED***
	***REMOVED***

	if u.Runner.Bundle.Options.NoVUConnectionReuse.Bool ***REMOVED***
		u.Transport.CloseIdleConnections()
	***REMOVED***

	sampleTags := metrics.NewSampleTags(u.state.CloneTags())
	u.state.Samples <- u.Dialer.GetTrail(
		startTime, endTime, isFullIteration, isDefault, sampleTags, u.Runner.builtinMetrics)

	return v, isFullIteration, endTime.Sub(startTime), err
***REMOVED***

func (u *ActiveVU) incrIteration() ***REMOVED***
	u.iteration++
	u.state.Iteration = u.iteration

	if _, ok := u.scenarioIter[u.scenarioName]; ok ***REMOVED***
		u.scenarioIter[u.scenarioName]++
	***REMOVED*** else ***REMOVED***
		u.scenarioIter[u.scenarioName] = 0
	***REMOVED***
	// TODO remove this
	if u.getNextIterationCounters != nil ***REMOVED***
		u.scIterLocal, u.scIterGlobal = u.getNextIterationCounters()
	***REMOVED***
***REMOVED***

type scriptException struct ***REMOVED***
	inner *goja.Exception
***REMOVED***

var (
	_ errext.Exception   = &scriptException***REMOVED******REMOVED***
	_ errext.HasExitCode = &scriptException***REMOVED******REMOVED***
	_ errext.HasHint     = &scriptException***REMOVED******REMOVED***
)

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

func (s *scriptException) Hint() string ***REMOVED***
	return "script exception"
***REMOVED***

func (s *scriptException) ExitCode() errext.ExitCode ***REMOVED***
	return exitcodes.ScriptException
***REMOVED***
