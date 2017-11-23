package pflag

import "strconv"

// -- uint32 value
type uint32Value uint32

func newUint32Value(val uint32, p *uint32) *uint32Value ***REMOVED***
	*p = val
	return (*uint32Value)(p)
***REMOVED***

func (i *uint32Value) Set(s string) error ***REMOVED***
	v, err := strconv.ParseUint(s, 0, 32)
	*i = uint32Value(v)
	return err
***REMOVED***

func (i *uint32Value) Type() string ***REMOVED***
	return "uint32"
***REMOVED***

func (i *uint32Value) String() string ***REMOVED*** return strconv.FormatUint(uint64(*i), 10) ***REMOVED***

func uint32Conv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v, err := strconv.ParseUint(sval, 0, 32)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return uint32(v), nil
***REMOVED***

// GetUint32 return the uint32 value of a flag with the given name
func (f *FlagSet) GetUint32(name string) (uint32, error) ***REMOVED***
	val, err := f.getFlagType(name, "uint32", uint32Conv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(uint32), nil
***REMOVED***

// Uint32Var defines a uint32 flag with specified name, default value, and usage string.
// The argument p points to a uint32 variable in which to store the value of the flag.
func (f *FlagSet) Uint32Var(p *uint32, name string, value uint32, usage string) ***REMOVED***
	f.VarP(newUint32Value(value, p), name, "", usage)
***REMOVED***

// Uint32VarP is like Uint32Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Uint32VarP(p *uint32, name, shorthand string, value uint32, usage string) ***REMOVED***
	f.VarP(newUint32Value(value, p), name, shorthand, usage)
***REMOVED***

// Uint32Var defines a uint32 flag with specified name, default value, and usage string.
// The argument p points to a uint32  variable in which to store the value of the flag.
func Uint32Var(p *uint32, name string, value uint32, usage string) ***REMOVED***
	CommandLine.VarP(newUint32Value(value, p), name, "", usage)
***REMOVED***

// Uint32VarP is like Uint32Var, but accepts a shorthand letter that can be used after a single dash.
func Uint32VarP(p *uint32, name, shorthand string, value uint32, usage string) ***REMOVED***
	CommandLine.VarP(newUint32Value(value, p), name, shorthand, usage)
***REMOVED***

// Uint32 defines a uint32 flag with specified name, default value, and usage string.
// The return value is the address of a uint32  variable that stores the value of the flag.
func (f *FlagSet) Uint32(name string, value uint32, usage string) *uint32 ***REMOVED***
	p := new(uint32)
	f.Uint32VarP(p, name, "", value, usage)
	return p
***REMOVED***

// Uint32P is like Uint32, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Uint32P(name, shorthand string, value uint32, usage string) *uint32 ***REMOVED***
	p := new(uint32)
	f.Uint32VarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Uint32 defines a uint32 flag with specified name, default value, and usage string.
// The return value is the address of a uint32  variable that stores the value of the flag.
func Uint32(name string, value uint32, usage string) *uint32 ***REMOVED***
	return CommandLine.Uint32P(name, "", value, usage)
***REMOVED***

// Uint32P is like Uint32, but accepts a shorthand letter that can be used after a single dash.
func Uint32P(name, shorthand string, value uint32, usage string) *uint32 ***REMOVED***
	return CommandLine.Uint32P(name, shorthand, value, usage)
***REMOVED***
