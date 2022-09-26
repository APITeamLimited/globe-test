package netext

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/lib/testutils/mockresolver"
	"github.com/APITeamLimited/k6-worker/lib/types"
)

func TestDialerAddr(t *testing.T) ***REMOVED***
	t.Parallel()
	dialer := NewDialer(net.Dialer***REMOVED******REMOVED***, newResolver())
	dialer.Hosts = map[string]*lib.HostAddress***REMOVED***
		"example.com":                ***REMOVED***IP: net.ParseIP("3.4.5.6")***REMOVED***,
		"example.com:443":            ***REMOVED***IP: net.ParseIP("3.4.5.6"), Port: 8443***REMOVED***,
		"example.com:8080":           ***REMOVED***IP: net.ParseIP("3.4.5.6"), Port: 9090***REMOVED***,
		"example-deny-host.com":      ***REMOVED***IP: net.ParseIP("8.9.10.11")***REMOVED***,
		"example-ipv6.com":           ***REMOVED***IP: net.ParseIP("2001:db8::68")***REMOVED***,
		"example-ipv6.com:443":       ***REMOVED***IP: net.ParseIP("2001:db8::68"), Port: 8443***REMOVED***,
		"example-ipv6-deny-host.com": ***REMOVED***IP: net.ParseIP("::1")***REMOVED***,
	***REMOVED***

	ipNet, err := lib.ParseCIDR("8.9.10.0/24")
	require.NoError(t, err)

	ipV6Net, err := lib.ParseCIDR("::1/24")
	require.NoError(t, err)

	dialer.Blacklist = []*lib.IPNet***REMOVED***ipNet, ipV6Net***REMOVED***

	testCases := []struct ***REMOVED***
		address, expAddress, expErr string
	***REMOVED******REMOVED***
		// IPv4
		***REMOVED***"example-resolver.com:80", "1.2.3.4:80", ""***REMOVED***,
		***REMOVED***"example.com:80", "3.4.5.6:80", ""***REMOVED***,
		***REMOVED***"example.com:443", "3.4.5.6:8443", ""***REMOVED***,
		***REMOVED***"example.com:8080", "3.4.5.6:9090", ""***REMOVED***,
		***REMOVED***"1.2.3.4:80", "1.2.3.4:80", ""***REMOVED***,
		***REMOVED***"1.2.3.4", "", "address 1.2.3.4: missing port in address"***REMOVED***,
		***REMOVED***"example-deny-resolver.com:80", "", "IP (8.9.10.11) is in a blacklisted range (8.9.10.0/24)"***REMOVED***,
		***REMOVED***"example-deny-host.com:80", "", "IP (8.9.10.11) is in a blacklisted range (8.9.10.0/24)"***REMOVED***,
		***REMOVED***"no-such-host.com:80", "", "lookup no-such-host.com: no such host"***REMOVED***,

		// IPv6
		***REMOVED***"example-ipv6.com:443", "[2001:db8::68]:8443", ""***REMOVED***,
		***REMOVED***"[2001:db8:aaaa:1::100]:443", "[2001:db8:aaaa:1::100]:443", ""***REMOVED***,
		***REMOVED***"[::1.2.3.4]", "", "address [::1.2.3.4]: missing port in address"***REMOVED***,
		***REMOVED***"example-ipv6-deny-resolver.com:80", "", "IP (::1) is in a blacklisted range (::/24)"***REMOVED***,
		***REMOVED***"example-ipv6-deny-host.com:80", "", "IP (::1) is in a blacklisted range (::/24)"***REMOVED***,
		***REMOVED***"example-ipv6-deny-host.com:80", "", "IP (::1) is in a blacklisted range (::/24)"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc

		t.Run(tc.address, func(t *testing.T) ***REMOVED***
			t.Parallel()
			addr, err := dialer.getDialAddr(tc.address)

			if tc.expErr != "" ***REMOVED***
				require.EqualError(t, err, tc.expErr)
			***REMOVED*** else ***REMOVED***
				require.NoError(t, err)
				require.Equal(t, tc.expAddress, addr)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestDialerAddrBlockHostnamesStar(t *testing.T) ***REMOVED***
	t.Parallel()
	dialer := NewDialer(net.Dialer***REMOVED******REMOVED***, newResolver())
	dialer.Hosts = map[string]*lib.HostAddress***REMOVED***
		"example.com": ***REMOVED***IP: net.ParseIP("3.4.5.6")***REMOVED***,
	***REMOVED***

	blocked, err := types.NewHostnameTrie([]string***REMOVED***"*"***REMOVED***)
	require.NoError(t, err)
	dialer.BlockedHostnames = blocked
	testCases := []struct ***REMOVED***
		address, expAddress, expErr string
	***REMOVED******REMOVED***
		// IPv4
		***REMOVED***"example.com:80", "", "hostname (example.com) is in a blocked pattern (*)"***REMOVED***,
		***REMOVED***"example.com:443", "", "hostname (example.com) is in a blocked pattern (*)"***REMOVED***,
		***REMOVED***"not.com:30", "", "hostname (not.com) is in a blocked pattern (*)"***REMOVED***,
		***REMOVED***"1.2.3.4:80", "1.2.3.4:80", ""***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc

		t.Run(tc.address, func(t *testing.T) ***REMOVED***
			t.Parallel()
			addr, err := dialer.getDialAddr(tc.address)

			if tc.expErr != "" ***REMOVED***
				require.EqualError(t, err, tc.expErr)
			***REMOVED*** else ***REMOVED***
				require.NoError(t, err)
				require.Equal(t, tc.expAddress, addr)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func newResolver() *mockresolver.MockResolver ***REMOVED***
	return mockresolver.New(
		map[string][]net.IP***REMOVED***
			"example-resolver.com":           ***REMOVED***net.ParseIP("1.2.3.4")***REMOVED***,
			"example-deny-resolver.com":      ***REMOVED***net.ParseIP("8.9.10.11")***REMOVED***,
			"example-ipv6-deny-resolver.com": ***REMOVED***net.ParseIP("::1")***REMOVED***,
		***REMOVED***, nil,
	)
***REMOVED***
