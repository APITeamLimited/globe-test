package pflag

import (
	"fmt"
	"net"
	"strings"
)

// -- net.IP value
type ipValue net.IP

func newIPValue(val net.IP, p *net.IP) *ipValue ***REMOVED***
	*p = val
	return (*ipValue)(p)
***REMOVED***

func (i *ipValue) String() string ***REMOVED*** return net.IP(*i).String() ***REMOVED***
func (i *ipValue) Set(s string) error ***REMOVED***
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip == nil ***REMOVED***
		return fmt.Errorf("failed to parse IP: %q", s)
	***REMOVED***
	*i = ipValue(ip)
	return nil
***REMOVED***

func (i *ipValue) Type() string ***REMOVED***
	return "ip"
***REMOVED***

func ipConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	ip := net.ParseIP(sval)
	if ip != nil ***REMOVED***
		return ip, nil
	***REMOVED***
	return nil, fmt.Errorf("invalid string being converted to IP address: %s", sval)
***REMOVED***

// GetIP return the net.IP value of a flag with the given name
func (f *FlagSet) GetIP(name string) (net.IP, error) ***REMOVED***
	val, err := f.getFlagType(name, "ip", ipConv)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return val.(net.IP), nil
***REMOVED***

// IPVar defines an net.IP flag with specified name, default value, and usage string.
// The argument p points to an net.IP variable in which to store the value of the flag.
func (f *FlagSet) IPVar(p *net.IP, name string, value net.IP, usage string) ***REMOVED***
	f.VarP(newIPValue(value, p), name, "", usage)
***REMOVED***

// IPVarP is like IPVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPVarP(p *net.IP, name, shorthand string, value net.IP, usage string) ***REMOVED***
	f.VarP(newIPValue(value, p), name, shorthand, usage)
***REMOVED***

// IPVar defines an net.IP flag with specified name, default value, and usage string.
// The argument p points to an net.IP variable in which to store the value of the flag.
func IPVar(p *net.IP, name string, value net.IP, usage string) ***REMOVED***
	CommandLine.VarP(newIPValue(value, p), name, "", usage)
***REMOVED***

// IPVarP is like IPVar, but accepts a shorthand letter that can be used after a single dash.
func IPVarP(p *net.IP, name, shorthand string, value net.IP, usage string) ***REMOVED***
	CommandLine.VarP(newIPValue(value, p), name, shorthand, usage)
***REMOVED***

// IP defines an net.IP flag with specified name, default value, and usage string.
// The return value is the address of an net.IP variable that stores the value of the flag.
func (f *FlagSet) IP(name string, value net.IP, usage string) *net.IP ***REMOVED***
	p := new(net.IP)
	f.IPVarP(p, name, "", value, usage)
	return p
***REMOVED***

// IPP is like IP, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPP(name, shorthand string, value net.IP, usage string) *net.IP ***REMOVED***
	p := new(net.IP)
	f.IPVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// IP defines an net.IP flag with specified name, default value, and usage string.
// The return value is the address of an net.IP variable that stores the value of the flag.
func IP(name string, value net.IP, usage string) *net.IP ***REMOVED***
	return CommandLine.IPP(name, "", value, usage)
***REMOVED***

// IPP is like IP, but accepts a shorthand letter that can be used after a single dash.
func IPP(name, shorthand string, value net.IP, usage string) *net.IP ***REMOVED***
	return CommandLine.IPP(name, shorthand, value, usage)
***REMOVED***
