package dnscache
// Package dnscache caches DNS lookups

import (
  "net"
  "sync"
  "time"
)

type Resolver struct ***REMOVED***
  lock sync.RWMutex
  cache map[string][]net.IP
***REMOVED***

func New(refreshRate time.Duration) *Resolver ***REMOVED***
  resolver := &Resolver ***REMOVED***
    cache: make(map[string][]net.IP, 64),
  ***REMOVED***
  if refreshRate > 0 ***REMOVED***
    go resolver.autoRefresh(refreshRate)
  ***REMOVED***
  return resolver
***REMOVED***

func (r *Resolver) Fetch(address string) ([]net.IP, error) ***REMOVED***
  r.lock.RLock()
  ips, exists := r.cache[address]
  r.lock.RUnlock()
  if exists ***REMOVED*** return ips, nil ***REMOVED***

  return r.Lookup(address)
***REMOVED***

func (r *Resolver) FetchOne(address string) (net.IP, error) ***REMOVED***
  ips, err := r.Fetch(address)
  if err != nil || len(ips) == 0 ***REMOVED*** return nil, err***REMOVED***
  return ips[0], nil
***REMOVED***

func (r *Resolver) FetchOneString(address string) (string, error) ***REMOVED***
  ip, err := r.FetchOne(address)
  if err != nil || ip == nil ***REMOVED*** return "", err ***REMOVED***
  return ip.String(), nil
***REMOVED***

func (r *Resolver) Refresh() ***REMOVED***
  i := 0
  r.lock.RLock()
  addresses := make([]string, len(r.cache))
  for key, _ := range r.cache ***REMOVED***
    addresses[i] = key
    i++
  ***REMOVED***
  r.lock.RUnlock()

  for _, address := range addresses ***REMOVED***
    r.Lookup(address)
    time.Sleep(time.Second * 2)
  ***REMOVED***
***REMOVED***

func (r *Resolver) Lookup(address string) ([]net.IP, error) ***REMOVED***
  ips, err := net.LookupIP(address)
  if err != nil ***REMOVED*** return nil, err ***REMOVED***

  r.lock.Lock()
  r.cache[address] = ips
  r.lock.Unlock()
  return ips, nil
***REMOVED***

func (r *Resolver) autoRefresh(rate time.Duration) ***REMOVED***
  for ***REMOVED***
    time.Sleep(rate)
    r.Refresh()
  ***REMOVED***
***REMOVED***
