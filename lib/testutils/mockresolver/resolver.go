package mockresolver

import (
	"fmt"
	"net"
	"sync"
)

// MockResolver implements netext.Resolver, and allows changing the host
// mapping at runtime.
type MockResolver struct ***REMOVED***
	m        sync.RWMutex
	hosts    map[string][]net.IP
	fallback func(host string) ([]net.IP, error)
***REMOVED***

// New returns a new MockResolver.
func New(hosts map[string][]net.IP, fallback func(host string) ([]net.IP, error)) *MockResolver ***REMOVED***
	if hosts == nil ***REMOVED***
		hosts = make(map[string][]net.IP)
	***REMOVED***
	return &MockResolver***REMOVED***hosts: hosts, fallback: fallback***REMOVED***
***REMOVED***

// LookupIP returns the first IP mapped for host.
func (r *MockResolver) LookupIP(host string) (net.IP, error) ***REMOVED***
	if ips, err := r.LookupIPAll(host); err != nil ***REMOVED***
		return nil, err
	***REMOVED*** else if len(ips) > 0 ***REMOVED***
		return ips[0], nil
	***REMOVED***
	return nil, nil
***REMOVED***

// LookupIPAll returns all IPs mapped for host. It mimics the net.LookupIP
// signature so that it can be used to mock netext.LookupIP in tests.
func (r *MockResolver) LookupIPAll(host string) ([]net.IP, error) ***REMOVED***
	r.m.RLock()
	defer r.m.RUnlock()
	if ips, ok := r.hosts[host]; ok ***REMOVED***
		return ips, nil
	***REMOVED***
	if r.fallback != nil ***REMOVED***
		return r.fallback(host)
	***REMOVED***
	return nil, fmt.Errorf("lookup %s: no such host", host)
***REMOVED***

// Set the host to resolve to ip.
func (r *MockResolver) Set(host, ip string) ***REMOVED***
	r.m.Lock()
	defer r.m.Unlock()
	r.hosts[host] = []net.IP***REMOVED***net.ParseIP(ip)***REMOVED***
***REMOVED***

// Unset removes the host.
func (r *MockResolver) Unset(host string) ***REMOVED***
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.hosts, host)
***REMOVED***
