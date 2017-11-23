package pflag

import (
	"fmt"
	"net"
	"strings"
)

// IPNet adapts net.IPNet for use as a flag.
type ipNetValue net.IPNet

func (ipnet ipNetValue) String() string ***REMOVED***
	n := net.IPNet(ipnet)
	return n.String()
***REMOVED***

func (ipnet *ipNetValue) Set(value string) error ***REMOVED***
	_, n, err := net.ParseCIDR(strings.TrimSpace(value))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*ipnet = ipNetValue(*n)
	return nil
***REMOVED***

func (*ipNetValue) Type() string ***REMOVED***
	return "ipNet"
***REMOVED***

func newIPNetValue(val net.IPNet, p *net.IPNet) *ipNetValue ***REMOVED***
	*p = val
	return (*ipNetValue)(p)
***REMOVED***

func ipNetConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	_, n, err := net.ParseCIDR(strings.TrimSpace(sval))
	if err == nil ***REMOVED***
		return *n, nil
	***REMOVED***
	return nil, fmt.Errorf("invalid string being converted to IPNet: %s", sval)
***REMOVED***

// GetIPNet return the net.IPNet value of a flag with the given name
func (f *FlagSet) GetIPNet(name string) (net.IPNet, error) ***REMOVED***
	val, err := f.getFlagType(name, "ipNet", ipNetConv)
	if err != nil ***REMOVED***
		return net.IPNet***REMOVED******REMOVED***, err
	***REMOVED***
	return val.(net.IPNet), nil
***REMOVED***

// IPNetVar defines an net.IPNet flag with specified name, default value, and usage string.
// The argument p points to an net.IPNet variable in which to store the value of the flag.
func (f *FlagSet) IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string) ***REMOVED***
	f.VarP(newIPNetValue(value, p), name, "", usage)
***REMOVED***

// IPNetVarP is like IPNetVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPNetVarP(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) ***REMOVED***
	f.VarP(newIPNetValue(value, p), name, shorthand, usage)
***REMOVED***

// IPNetVar defines an net.IPNet flag with specified name, default value, and usage string.
// The argument p points to an net.IPNet variable in which to store the value of the flag.
func IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string) ***REMOVED***
	CommandLine.VarP(newIPNetValue(value, p), name, "", usage)
***REMOVED***

// IPNetVarP is like IPNetVar, but accepts a shorthand letter that can be used after a single dash.
func IPNetVarP(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) ***REMOVED***
	CommandLine.VarP(newIPNetValue(value, p), name, shorthand, usage)
***REMOVED***

// IPNet defines an net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of an net.IPNet variable that stores the value of the flag.
func (f *FlagSet) IPNet(name string, value net.IPNet, usage string) *net.IPNet ***REMOVED***
	p := new(net.IPNet)
	f.IPNetVarP(p, name, "", value, usage)
	return p
***REMOVED***

// IPNetP is like IPNet, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet ***REMOVED***
	p := new(net.IPNet)
	f.IPNetVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// IPNet defines an net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of an net.IPNet variable that stores the value of the flag.
func IPNet(name string, value net.IPNet, usage string) *net.IPNet ***REMOVED***
	return CommandLine.IPNetP(name, "", value, usage)
***REMOVED***

// IPNetP is like IPNet, but accepts a shorthand letter that can be used after a single dash.
func IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet ***REMOVED***
	return CommandLine.IPNetP(name, shorthand, value, usage)
***REMOVED***
