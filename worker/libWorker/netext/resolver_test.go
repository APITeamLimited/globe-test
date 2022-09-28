package netext

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/APITeamLimited/globe-test/worker/libWorker/testutils/mockresolver"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
)

func TestResolver(t *testing.T) ***REMOVED***
	t.Parallel()

	host := "myhost"
	mr := mockresolver.New(map[string][]net.IP***REMOVED***
		host: ***REMOVED***
			net.ParseIP("127.0.0.10"),
			net.ParseIP("127.0.0.11"),
			net.ParseIP("127.0.0.12"),
			net.ParseIP("2001:db8::10"),
			net.ParseIP("2001:db8::11"),
			net.ParseIP("2001:db8::12"),
		***REMOVED***,
	***REMOVED***, nil)

	t.Run("LookupIP", func(t *testing.T) ***REMOVED***
		t.Parallel()
		testCases := []struct ***REMOVED***
			ttl   time.Duration
			sel   types.DNSSelect
			pol   types.DNSPolicy
			expIP []net.IP
		***REMOVED******REMOVED***
			***REMOVED***
				0, types.DNSfirst, types.DNSpreferIPv4,
				[]net.IP***REMOVED***net.ParseIP("127.0.0.10")***REMOVED***,
			***REMOVED***,
			***REMOVED***
				time.Second, types.DNSfirst, types.DNSpreferIPv4,
				[]net.IP***REMOVED***net.ParseIP("127.0.0.10")***REMOVED***,
			***REMOVED***,
			***REMOVED***0, types.DNSroundRobin, types.DNSonlyIPv6, []net.IP***REMOVED***
				net.ParseIP("2001:db8::10"),
				net.ParseIP("2001:db8::11"),
				net.ParseIP("2001:db8::12"),
				net.ParseIP("2001:db8::10"),
			***REMOVED******REMOVED***,
			***REMOVED***
				0, types.DNSfirst, types.DNSpreferIPv6,
				[]net.IP***REMOVED***net.ParseIP("2001:db8::10")***REMOVED***,
			***REMOVED***,
			***REMOVED***0, types.DNSroundRobin, types.DNSpreferIPv4, []net.IP***REMOVED***
				net.ParseIP("127.0.0.10"),
				net.ParseIP("127.0.0.11"),
				net.ParseIP("127.0.0.12"),
				net.ParseIP("127.0.0.10"),
			***REMOVED******REMOVED***,
		***REMOVED***

		for _, tc := range testCases ***REMOVED***
			tc := tc
			t.Run(fmt.Sprintf("%s_%s_%s", tc.ttl, tc.sel, tc.pol), func(t *testing.T) ***REMOVED***
				t.Parallel()
				r := NewResolver(mr.LookupIPAll, tc.ttl, tc.sel, tc.pol)
				ip, err := r.LookupIP(host)
				require.NoError(t, err)
				assert.Equal(t, tc.expIP[0], ip)

				if tc.ttl > 0 ***REMOVED***
					require.IsType(t, &cacheResolver***REMOVED******REMOVED***, r)
					cr := r.(*cacheResolver)
					assert.Len(t, cr.cache, 1)
					assert.Equal(t, tc.ttl, cr.ttl)
					firstLookup := cr.cache[host].lastLookup
					time.Sleep(cr.ttl + 100*time.Millisecond)
					_, err = r.LookupIP(host)
					require.NoError(t, err)
					assert.True(t, cr.cache[host].lastLookup.After(firstLookup))
				***REMOVED***

				if tc.sel == types.DNSroundRobin ***REMOVED***
					ips := []net.IP***REMOVED***ip***REMOVED***
					for i := 0; i < 3; i++ ***REMOVED***
						ip, err = r.LookupIP(host)
						require.NoError(t, err)
						ips = append(ips, ip)
					***REMOVED***
					assert.Equal(t, tc.expIP, ips)
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***
