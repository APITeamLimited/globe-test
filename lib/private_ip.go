package lib

import (
	"fmt"
	"net"
)

var privateIPBlocks []*net.IPNet

func init() ***REMOVED***
	for _, cidr := range []string***REMOVED***
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	***REMOVED*** ***REMOVED***
		_, block, err := net.ParseCIDR(cidr)
		if err != nil ***REMOVED***
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		***REMOVED***
		privateIPBlocks = append(privateIPBlocks, block)
	***REMOVED***
***REMOVED***

func IsPrivateIP(ip net.IP) bool ***REMOVED***
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ***REMOVED***
		return true
	***REMOVED***

	for _, block := range privateIPBlocks ***REMOVED***
		if block.Contains(ip) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func IsPrivateIPString(ip string) bool ***REMOVED***
	parsedIp := net.ParseIP(ip)

	if parsedIp == nil ***REMOVED***
		return false
	***REMOVED***

	return IsPrivateIP(parsedIp)
***REMOVED***
