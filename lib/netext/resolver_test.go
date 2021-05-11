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

package netext

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib/testutils/mockresolver"
	"go.k6.io/k6/lib/types"
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
