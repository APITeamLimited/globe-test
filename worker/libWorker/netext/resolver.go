package netext

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

// MultiResolver returns all IP addresses for the given host.
type MultiResolver func(host string) ([]net.IP, error)

// Resolver is an interface that returns DNS information about a given host.
type Resolver interface ***REMOVED***
	LookupIP(host string) (net.IP, error)
***REMOVED***

type resolver struct ***REMOVED***
	resolve     MultiResolver
	selectIndex types.DNSSelect
	policy      types.DNSPolicy
	rrm         *sync.Mutex
	rand        *rand.Rand
	roundRobin  map[string]uint8
***REMOVED***

type cacheRecord struct ***REMOVED***
	ips        []net.IP
	lastLookup time.Time
***REMOVED***

type cacheResolver struct ***REMOVED***
	resolver
	ttl   time.Duration
	cm    *sync.Mutex
	cache map[string]cacheRecord
***REMOVED***

// NewResolver returns a new DNS resolver. If ttl is not 0, responses
// will be cached per host for the specified period. The IP returned from
// LookupIP() will be selected based on the given sel and pol values.
func NewResolver(
	actRes MultiResolver, ttl time.Duration, sel types.DNSSelect, pol types.DNSPolicy,
) Resolver ***REMOVED***
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	res := resolver***REMOVED***
		resolve:     actRes,
		selectIndex: sel,
		policy:      pol,
		rrm:         &sync.Mutex***REMOVED******REMOVED***,
		rand:        r,
		roundRobin:  make(map[string]uint8),
	***REMOVED***
	if ttl == 0 ***REMOVED***
		return &res
	***REMOVED***
	return &cacheResolver***REMOVED***
		resolver: res,
		ttl:      ttl,
		cm:       &sync.Mutex***REMOVED******REMOVED***,
		cache:    make(map[string]cacheRecord),
	***REMOVED***
***REMOVED***

// LookupIP returns a single IP resolved for host, selected according to the
// configured select and policy options.
func (r *resolver) LookupIP(host string) (net.IP, error) ***REMOVED***
	ips, err := r.resolve(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ips = r.applyPolicy(ips)
	return r.selectOne(host, ips), nil
***REMOVED***

// LookupIP returns a single IP resolved for host, selected according to the
// configured select and policy options. Results are cached per host and will be
// refreshed if the last lookup time exceeds the configured TTL (not the TTL
// returned in the DNS record).
func (r *cacheResolver) LookupIP(host string) (net.IP, error) ***REMOVED***
	r.cm.Lock()

	var ips []net.IP
	// TODO: Invalidate? When?
	if cr, ok := r.cache[host]; ok && time.Now().Before(cr.lastLookup.Add(r.ttl)) ***REMOVED***
		ips = cr.ips
	***REMOVED*** else ***REMOVED***
		r.cm.Unlock() // The lookup could take some time, so unlock momentarily.
		var err error
		ips, err = r.resolve(host)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ips = r.applyPolicy(ips)
		r.cm.Lock()
		r.cache[host] = cacheRecord***REMOVED***ips: ips, lastLookup: time.Now()***REMOVED***
	***REMOVED***

	r.cm.Unlock()

	return r.selectOne(host, ips), nil
***REMOVED***

func (r *resolver) selectOne(host string, ips []net.IP) net.IP ***REMOVED***
	if len(ips) == 0 ***REMOVED***
		return nil
	***REMOVED***

	var ip net.IP
	switch r.selectIndex ***REMOVED***
	case types.DNSfirst:
		return ips[0]
	case types.DNSroundRobin:
		r.rrm.Lock()
		// NOTE: This index approach is not stable and might result in returning
		// repeated or skipped IPs if the records change during a test run.
		ip = ips[int(r.roundRobin[host])%len(ips)]
		r.roundRobin[host]++
		r.rrm.Unlock()
	case types.DNSrandom:
		r.rrm.Lock()
		ip = ips[r.rand.Intn(len(ips))]
		r.rrm.Unlock()
	***REMOVED***

	return ip
***REMOVED***

func (r *resolver) applyPolicy(ips []net.IP) (retIPs []net.IP) ***REMOVED***
	if r.policy == types.DNSany ***REMOVED***
		return ips
	***REMOVED***
	ip4, ip6 := groupByVersion(ips)
	switch r.policy ***REMOVED***
	case types.DNSpreferIPv4:
		retIPs = ip4
		if len(retIPs) == 0 ***REMOVED***
			retIPs = ip6
		***REMOVED***
	case types.DNSpreferIPv6:
		retIPs = ip6
		if len(retIPs) == 0 ***REMOVED***
			retIPs = ip4
		***REMOVED***
	case types.DNSonlyIPv4:
		retIPs = ip4
	case types.DNSonlyIPv6:
		retIPs = ip6
	// Already checked above, but added to satisfy 'exhaustive' linter.
	case types.DNSany:
		retIPs = ips
	***REMOVED***

	return
***REMOVED***

func groupByVersion(ips []net.IP) (ip4 []net.IP, ip6 []net.IP) ***REMOVED***
	for _, ip := range ips ***REMOVED***
		if ip.To4() != nil ***REMOVED***
			ip4 = append(ip4, ip)
		***REMOVED*** else ***REMOVED***
			ip6 = append(ip6, ip)
		***REMOVED***
	***REMOVED***

	return
***REMOVED***
