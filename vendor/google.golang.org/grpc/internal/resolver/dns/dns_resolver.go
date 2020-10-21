/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package dns implements a dns resolver to be installed as the default resolver
// in grpc.
package dns

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	grpclbstate "google.golang.org/grpc/balancer/grpclb/state"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/envconfig"
	"google.golang.org/grpc/internal/grpcrand"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

// EnableSRVLookups controls whether the DNS resolver attempts to fetch gRPCLB
// addresses from SRV records.  Must not be changed after init time.
var EnableSRVLookups = false

var logger = grpclog.Component("dns")

func init() ***REMOVED***
	resolver.Register(NewBuilder())
***REMOVED***

const (
	defaultPort       = "443"
	defaultDNSSvrPort = "53"
	golang            = "GO"
	// txtPrefix is the prefix string to be prepended to the host name for txt record lookup.
	txtPrefix = "_grpc_config."
	// In DNS, service config is encoded in a TXT record via the mechanism
	// described in RFC-1464 using the attribute name grpc_config.
	txtAttribute = "grpc_config="
)

var (
	errMissingAddr = errors.New("dns resolver: missing address")

	// Addresses ending with a colon that is supposed to be the separator
	// between host and port is not allowed.  E.g. "::" is a valid address as
	// it is an IPv6 address (host only) and "[::]:" is invalid as it ends with
	// a colon as the host and port separator
	errEndsWithColon = errors.New("dns resolver: missing port after port-separator colon")
)

var (
	defaultResolver netResolver = net.DefaultResolver
	// To prevent excessive re-resolution, we enforce a rate limit on DNS
	// resolution requests.
	minDNSResRate = 30 * time.Second
)

var customAuthorityDialler = func(authority string) func(ctx context.Context, network, address string) (net.Conn, error) ***REMOVED***
	return func(ctx context.Context, network, address string) (net.Conn, error) ***REMOVED***
		var dialer net.Dialer
		return dialer.DialContext(ctx, network, authority)
	***REMOVED***
***REMOVED***

var customAuthorityResolver = func(authority string) (netResolver, error) ***REMOVED***
	host, port, err := parseTarget(authority, defaultDNSSvrPort)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	authorityWithPort := net.JoinHostPort(host, port)

	return &net.Resolver***REMOVED***
		PreferGo: true,
		Dial:     customAuthorityDialler(authorityWithPort),
	***REMOVED***, nil
***REMOVED***

// NewBuilder creates a dnsBuilder which is used to factory DNS resolvers.
func NewBuilder() resolver.Builder ***REMOVED***
	return &dnsBuilder***REMOVED******REMOVED***
***REMOVED***

type dnsBuilder struct***REMOVED******REMOVED***

// Build creates and starts a DNS resolver that watches the name resolution of the target.
func (b *dnsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) ***REMOVED***
	host, port, err := parseTarget(target.Endpoint, defaultPort)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// IP address.
	if ipAddr, ok := formatIP(host); ok ***REMOVED***
		addr := []resolver.Address***REMOVED******REMOVED***Addr: ipAddr + ":" + port***REMOVED******REMOVED***
		cc.UpdateState(resolver.State***REMOVED***Addresses: addr***REMOVED***)
		return deadResolver***REMOVED******REMOVED***, nil
	***REMOVED***

	// DNS address (non-IP).
	ctx, cancel := context.WithCancel(context.Background())
	d := &dnsResolver***REMOVED***
		host:                 host,
		port:                 port,
		ctx:                  ctx,
		cancel:               cancel,
		cc:                   cc,
		rn:                   make(chan struct***REMOVED******REMOVED***, 1),
		disableServiceConfig: opts.DisableServiceConfig,
	***REMOVED***

	if target.Authority == "" ***REMOVED***
		d.resolver = defaultResolver
	***REMOVED*** else ***REMOVED***
		d.resolver, err = customAuthorityResolver(target.Authority)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	d.wg.Add(1)
	go d.watcher()
	d.ResolveNow(resolver.ResolveNowOptions***REMOVED******REMOVED***)
	return d, nil
***REMOVED***

// Scheme returns the naming scheme of this resolver builder, which is "dns".
func (b *dnsBuilder) Scheme() string ***REMOVED***
	return "dns"
***REMOVED***

type netResolver interface ***REMOVED***
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupTXT(ctx context.Context, name string) (txts []string, err error)
***REMOVED***

// deadResolver is a resolver that does nothing.
type deadResolver struct***REMOVED******REMOVED***

func (deadResolver) ResolveNow(resolver.ResolveNowOptions) ***REMOVED******REMOVED***

func (deadResolver) Close() ***REMOVED******REMOVED***

// dnsResolver watches for the name resolution update for a non-IP target.
type dnsResolver struct ***REMOVED***
	host     string
	port     string
	resolver netResolver
	ctx      context.Context
	cancel   context.CancelFunc
	cc       resolver.ClientConn
	// rn channel is used by ResolveNow() to force an immediate resolution of the target.
	rn chan struct***REMOVED******REMOVED***
	// wg is used to enforce Close() to return after the watcher() goroutine has finished.
	// Otherwise, data race will be possible. [Race Example] in dns_resolver_test we
	// replace the real lookup functions with mocked ones to facilitate testing.
	// If Close() doesn't wait for watcher() goroutine finishes, race detector sometimes
	// will warns lookup (READ the lookup function pointers) inside watcher() goroutine
	// has data race with replaceNetFunc (WRITE the lookup function pointers).
	wg                   sync.WaitGroup
	disableServiceConfig bool
***REMOVED***

// ResolveNow invoke an immediate resolution of the target that this dnsResolver watches.
func (d *dnsResolver) ResolveNow(resolver.ResolveNowOptions) ***REMOVED***
	select ***REMOVED***
	case d.rn <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	default:
	***REMOVED***
***REMOVED***

// Close closes the dnsResolver.
func (d *dnsResolver) Close() ***REMOVED***
	d.cancel()
	d.wg.Wait()
***REMOVED***

func (d *dnsResolver) watcher() ***REMOVED***
	defer d.wg.Done()
	for ***REMOVED***
		select ***REMOVED***
		case <-d.ctx.Done():
			return
		case <-d.rn:
		***REMOVED***

		state, err := d.lookup()
		if err != nil ***REMOVED***
			d.cc.ReportError(err)
		***REMOVED*** else ***REMOVED***
			d.cc.UpdateState(*state)
		***REMOVED***

		// Sleep to prevent excessive re-resolutions. Incoming resolution requests
		// will be queued in d.rn.
		t := time.NewTimer(minDNSResRate)
		select ***REMOVED***
		case <-t.C:
		case <-d.ctx.Done():
			t.Stop()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *dnsResolver) lookupSRV() ([]resolver.Address, error) ***REMOVED***
	if !EnableSRVLookups ***REMOVED***
		return nil, nil
	***REMOVED***
	var newAddrs []resolver.Address
	_, srvs, err := d.resolver.LookupSRV(d.ctx, "grpclb", "tcp", d.host)
	if err != nil ***REMOVED***
		err = handleDNSError(err, "SRV") // may become nil
		return nil, err
	***REMOVED***
	for _, s := range srvs ***REMOVED***
		lbAddrs, err := d.resolver.LookupHost(d.ctx, s.Target)
		if err != nil ***REMOVED***
			err = handleDNSError(err, "A") // may become nil
			if err == nil ***REMOVED***
				// If there are other SRV records, look them up and ignore this
				// one that does not exist.
				continue
			***REMOVED***
			return nil, err
		***REMOVED***
		for _, a := range lbAddrs ***REMOVED***
			ip, ok := formatIP(a)
			if !ok ***REMOVED***
				return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
			***REMOVED***
			addr := ip + ":" + strconv.Itoa(int(s.Port))
			newAddrs = append(newAddrs, resolver.Address***REMOVED***Addr: addr, ServerName: s.Target***REMOVED***)
		***REMOVED***
	***REMOVED***
	return newAddrs, nil
***REMOVED***

var filterError = func(err error) error ***REMOVED***
	if dnsErr, ok := err.(*net.DNSError); ok && !dnsErr.IsTimeout && !dnsErr.IsTemporary ***REMOVED***
		// Timeouts and temporary errors should be communicated to gRPC to
		// attempt another DNS query (with backoff).  Other errors should be
		// suppressed (they may represent the absence of a TXT record).
		return nil
	***REMOVED***
	return err
***REMOVED***

func handleDNSError(err error, lookupType string) error ***REMOVED***
	err = filterError(err)
	if err != nil ***REMOVED***
		err = fmt.Errorf("dns: %v record lookup error: %v", lookupType, err)
		logger.Info(err)
	***REMOVED***
	return err
***REMOVED***

func (d *dnsResolver) lookupTXT() *serviceconfig.ParseResult ***REMOVED***
	ss, err := d.resolver.LookupTXT(d.ctx, txtPrefix+d.host)
	if err != nil ***REMOVED***
		if envconfig.TXTErrIgnore ***REMOVED***
			return nil
		***REMOVED***
		if err = handleDNSError(err, "TXT"); err != nil ***REMOVED***
			return &serviceconfig.ParseResult***REMOVED***Err: err***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	var res string
	for _, s := range ss ***REMOVED***
		res += s
	***REMOVED***

	// TXT record must have "grpc_config=" attribute in order to be used as service config.
	if !strings.HasPrefix(res, txtAttribute) ***REMOVED***
		logger.Warningf("dns: TXT record %v missing %v attribute", res, txtAttribute)
		// This is not an error; it is the equivalent of not having a service config.
		return nil
	***REMOVED***
	sc := canaryingSC(strings.TrimPrefix(res, txtAttribute))
	return d.cc.ParseServiceConfig(sc)
***REMOVED***

func (d *dnsResolver) lookupHost() ([]resolver.Address, error) ***REMOVED***
	var newAddrs []resolver.Address
	addrs, err := d.resolver.LookupHost(d.ctx, d.host)
	if err != nil ***REMOVED***
		err = handleDNSError(err, "A")
		return nil, err
	***REMOVED***
	for _, a := range addrs ***REMOVED***
		ip, ok := formatIP(a)
		if !ok ***REMOVED***
			return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
		***REMOVED***
		addr := ip + ":" + d.port
		newAddrs = append(newAddrs, resolver.Address***REMOVED***Addr: addr***REMOVED***)
	***REMOVED***
	return newAddrs, nil
***REMOVED***

func (d *dnsResolver) lookup() (*resolver.State, error) ***REMOVED***
	srv, srvErr := d.lookupSRV()
	addrs, hostErr := d.lookupHost()
	if hostErr != nil && (srvErr != nil || len(srv) == 0) ***REMOVED***
		return nil, hostErr
	***REMOVED***

	state := resolver.State***REMOVED***Addresses: addrs***REMOVED***
	if len(srv) > 0 ***REMOVED***
		state = grpclbstate.Set(state, &grpclbstate.State***REMOVED***BalancerAddresses: srv***REMOVED***)
	***REMOVED***
	if !d.disableServiceConfig ***REMOVED***
		state.ServiceConfig = d.lookupTXT()
	***REMOVED***
	return &state, nil
***REMOVED***

// formatIP returns ok = false if addr is not a valid textual representation of an IP address.
// If addr is an IPv4 address, return the addr and ok = true.
// If addr is an IPv6 address, return the addr enclosed in square brackets and ok = true.
func formatIP(addr string) (addrIP string, ok bool) ***REMOVED***
	ip := net.ParseIP(addr)
	if ip == nil ***REMOVED***
		return "", false
	***REMOVED***
	if ip.To4() != nil ***REMOVED***
		return addr, true
	***REMOVED***
	return "[" + addr + "]", true
***REMOVED***

// parseTarget takes the user input target string and default port, returns formatted host and port info.
// If target doesn't specify a port, set the port to be the defaultPort.
// If target is in IPv6 format and host-name is enclosed in square brackets, brackets
// are stripped when setting the host.
// examples:
// target: "www.google.com" defaultPort: "443" returns host: "www.google.com", port: "443"
// target: "ipv4-host:80" defaultPort: "443" returns host: "ipv4-host", port: "80"
// target: "[ipv6-host]" defaultPort: "443" returns host: "ipv6-host", port: "443"
// target: ":80" defaultPort: "443" returns host: "localhost", port: "80"
func parseTarget(target, defaultPort string) (host, port string, err error) ***REMOVED***
	if target == "" ***REMOVED***
		return "", "", errMissingAddr
	***REMOVED***
	if ip := net.ParseIP(target); ip != nil ***REMOVED***
		// target is an IPv4 or IPv6(without brackets) address
		return target, defaultPort, nil
	***REMOVED***
	if host, port, err = net.SplitHostPort(target); err == nil ***REMOVED***
		if port == "" ***REMOVED***
			// If the port field is empty (target ends with colon), e.g. "[::1]:", this is an error.
			return "", "", errEndsWithColon
		***REMOVED***
		// target has port, i.e ipv4-host:port, [ipv6-host]:port, host-name:port
		if host == "" ***REMOVED***
			// Keep consistent with net.Dial(): If the host is empty, as in ":80", the local system is assumed.
			host = "localhost"
		***REMOVED***
		return host, port, nil
	***REMOVED***
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil ***REMOVED***
		// target doesn't have port
		return host, port, nil
	***REMOVED***
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
***REMOVED***

type rawChoice struct ***REMOVED***
	ClientLanguage *[]string        `json:"clientLanguage,omitempty"`
	Percentage     *int             `json:"percentage,omitempty"`
	ClientHostName *[]string        `json:"clientHostName,omitempty"`
	ServiceConfig  *json.RawMessage `json:"serviceConfig,omitempty"`
***REMOVED***

func containsString(a *[]string, b string) bool ***REMOVED***
	if a == nil ***REMOVED***
		return true
	***REMOVED***
	for _, c := range *a ***REMOVED***
		if c == b ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func chosenByPercentage(a *int) bool ***REMOVED***
	if a == nil ***REMOVED***
		return true
	***REMOVED***
	return grpcrand.Intn(100)+1 <= *a
***REMOVED***

func canaryingSC(js string) string ***REMOVED***
	if js == "" ***REMOVED***
		return ""
	***REMOVED***
	var rcs []rawChoice
	err := json.Unmarshal([]byte(js), &rcs)
	if err != nil ***REMOVED***
		logger.Warningf("dns: error parsing service config json: %v", err)
		return ""
	***REMOVED***
	cliHostname, err := os.Hostname()
	if err != nil ***REMOVED***
		logger.Warningf("dns: error getting client hostname: %v", err)
		return ""
	***REMOVED***
	var sc string
	for _, c := range rcs ***REMOVED***
		if !containsString(c.ClientLanguage, golang) ||
			!chosenByPercentage(c.Percentage) ||
			!containsString(c.ClientHostName, cliHostname) ||
			c.ServiceConfig == nil ***REMOVED***
			continue
		***REMOVED***
		sc = string(*c.ServiceConfig)
		break
	***REMOVED***
	return sc
***REMOVED***
