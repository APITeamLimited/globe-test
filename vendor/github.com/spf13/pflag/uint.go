package pflag

import "strconv"

// -- uint Value
type uintValue uint

func newUintValue(val uint, p *uint) *uintValue ***REMOVED***
	*p = val
	return (*uintValue)(p)
***REMOVED***

func (i *uintValue) Set(s string) error ***REMOVED***
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
***REMOVED***

func (i *uintValue) Type() string ***REMOVED***
	return "uint"
***REMOVED***

func (i *uintValue) String() string ***REMOVED*** return strconv.FormatUint(uint64(*i), 10) ***REMOVED***

func uintConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v, err := strconv.ParseUint(sval, 0, 0)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return uint(v), nil
***REMOVED***

// GetUint return the uint value of a flag with the given name
func (f *FlagSet) GetUint(name string) (uint, error) ***REMOVED***
	val, err := f.getFlagType(name, "uint", uintConv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(uint), nil
***REMOVED***

// UintVar defines a uint flag with specified name, default value, and usage string.
// The argument p points to a uint variable in which to store the value of the flag.
func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) ***REMOVED***
	f.VarP(newUintValue(value, p), name, "", usage)
***REMOVED***

// UintVarP is like UintVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UintVarP(p *uint, name, shorthand string, value uint, usage string) ***REMOVED***
	f.VarP(newUintValue(value, p), name, shorthand, usage)
***REMOVED***

// UintVar defines a uint flag with specified name, default value, and usage string.
// The argument p points to a uint  variable in which to store the value of the flag.
func UintVar(p *uint, name string, value uint, usage string) ***REMOVED***
	CommandLine.VarP(newUintValue(value, p), name, "", usage)
***REMOVED***

// UintVarP is like UintVar, but accepts a shorthand letter that can be used after a single dash.
func UintVarP(p *uint, name, shorthand string, value uint, usage string) ***REMOVED***
	CommandLine.VarP(newUintValue(value, p), name, shorthand, usage)
***REMOVED***

// Uint defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint  variable that stores the value of the flag.
func (f *FlagSet) Uint(name string, value uint, usage string) *uint ***REMOVED***
	p := new(uint)
	f.UintVarP(p, name, "", value, usage)
	return p
***REMOVED***

// UintP is like Uint, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UintP(name, shorthand string, value uint, usage string) *uint ***REMOVED***
	p := new(uint)
	f.UintVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Uint defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint  variable that stores the value of the flag.
func Uint(name string, value uint, usage string) *uint ***REMOVED***
	return CommandLine.UintP(name, "", value, usage)
***REMOVED***

// UintP is like Uint, but accepts a shorthand letter that can be used after a single dash.
func UintP(name, shorthand string, value uint, usage string) *uint ***REMOVED***
	return CommandLine.UintP(name, shorthand, value, usage)
***REMOVED***
