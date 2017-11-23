package pflag

import (
	"fmt"
	"net"
	"strconv"
)

// -- net.IPMask value
type ipMaskValue net.IPMask

func newIPMaskValue(val net.IPMask, p *net.IPMask) *ipMaskValue ***REMOVED***
	*p = val
	return (*ipMaskValue)(p)
***REMOVED***

func (i *ipMaskValue) String() string ***REMOVED*** return net.IPMask(*i).String() ***REMOVED***
func (i *ipMaskValue) Set(s string) error ***REMOVED***
	ip := ParseIPv4Mask(s)
	if ip == nil ***REMOVED***
		return fmt.Errorf("failed to parse IP mask: %q", s)
	***REMOVED***
	*i = ipMaskValue(ip)
	return nil
***REMOVED***

func (i *ipMaskValue) Type() string ***REMOVED***
	return "ipMask"
***REMOVED***

// ParseIPv4Mask written in IP form (e.g. 255.255.255.0).
// This function should really belong to the net package.
func ParseIPv4Mask(s string) net.IPMask ***REMOVED***
	mask := net.ParseIP(s)
	if mask == nil ***REMOVED***
		if len(s) != 8 ***REMOVED***
			return nil
		***REMOVED***
		// net.IPMask.String() actually outputs things like ffffff00
		// so write a horrible parser for that as well  :-(
		m := []int***REMOVED******REMOVED***
		for i := 0; i < 4; i++ ***REMOVED***
			b := "0x" + s[2*i:2*i+2]
			d, err := strconv.ParseInt(b, 0, 0)
			if err != nil ***REMOVED***
				return nil
			***REMOVED***
			m = append(m, int(d))
		***REMOVED***
		s := fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
		mask = net.ParseIP(s)
		if mask == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return net.IPv4Mask(mask[12], mask[13], mask[14], mask[15])
***REMOVED***

func parseIPv4Mask(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	mask := ParseIPv4Mask(sval)
	if mask == nil ***REMOVED***
		return nil, fmt.Errorf("unable to parse %s as net.IPMask", sval)
	***REMOVED***
	return mask, nil
***REMOVED***

// GetIPv4Mask return the net.IPv4Mask value of a flag with the given name
func (f *FlagSet) GetIPv4Mask(name string) (net.IPMask, error) ***REMOVED***
	val, err := f.getFlagType(name, "ipMask", parseIPv4Mask)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return val.(net.IPMask), nil
***REMOVED***

// IPMaskVar defines an net.IPMask flag with specified name, default value, and usage string.
// The argument p points to an net.IPMask variable in which to store the value of the flag.
func (f *FlagSet) IPMaskVar(p *net.IPMask, name string, value net.IPMask, usage string) ***REMOVED***
	f.VarP(newIPMaskValue(value, p), name, "", usage)
***REMOVED***

// IPMaskVarP is like IPMaskVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPMaskVarP(p *net.IPMask, name, shorthand string, value net.IPMask, usage string) ***REMOVED***
	f.VarP(newIPMaskValue(value, p), name, shorthand, usage)
***REMOVED***

// IPMaskVar defines an net.IPMask flag with specified name, default value, and usage string.
// The argument p points to an net.IPMask variable in which to store the value of the flag.
func IPMaskVar(p *net.IPMask, name string, value net.IPMask, usage string) ***REMOVED***
	CommandLine.VarP(newIPMaskValue(value, p), name, "", usage)
***REMOVED***

// IPMaskVarP is like IPMaskVar, but accepts a shorthand letter that can be used after a single dash.
func IPMaskVarP(p *net.IPMask, name, shorthand string, value net.IPMask, usage string) ***REMOVED***
	CommandLine.VarP(newIPMaskValue(value, p), name, shorthand, usage)
***REMOVED***

// IPMask defines an net.IPMask flag with specified name, default value, and usage string.
// The return value is the address of an net.IPMask variable that stores the value of the flag.
func (f *FlagSet) IPMask(name string, value net.IPMask, usage string) *net.IPMask ***REMOVED***
	p := new(net.IPMask)
	f.IPMaskVarP(p, name, "", value, usage)
	return p
***REMOVED***

// IPMaskP is like IPMask, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPMaskP(name, shorthand string, value net.IPMask, usage string) *net.IPMask ***REMOVED***
	p := new(net.IPMask)
	f.IPMaskVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// IPMask defines an net.IPMask flag with specified name, default value, and usage string.
// The return value is the address of an net.IPMask variable that stores the value of the flag.
func IPMask(name string, value net.IPMask, usage string) *net.IPMask ***REMOVED***
	return CommandLine.IPMaskP(name, "", value, usage)
***REMOVED***

// IPMaskP is like IP, but accepts a shorthand letter that can be used after a single dash.
func IPMaskP(name, shorthand string, value net.IPMask, usage string) *net.IPMask ***REMOVED***
	return CommandLine.IPMaskP(name, shorthand, value, usage)
***REMOVED***
