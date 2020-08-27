/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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
