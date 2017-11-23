package pflag

import "strconv"

// -- int64 Value
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value ***REMOVED***
	*p = val
	return (*int64Value)(p)
***REMOVED***

func (i *int64Value) Set(s string) error ***REMOVED***
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
***REMOVED***

func (i *int64Value) Type() string ***REMOVED***
	return "int64"
***REMOVED***

func (i *int64Value) String() string ***REMOVED*** return strconv.FormatInt(int64(*i), 10) ***REMOVED***

func int64Conv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return strconv.ParseInt(sval, 0, 64)
***REMOVED***

// GetInt64 return the int64 value of a flag with the given name
func (f *FlagSet) GetInt64(name string) (int64, error) ***REMOVED***
	val, err := f.getFlagType(name, "int64", int64Conv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(int64), nil
***REMOVED***

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) ***REMOVED***
	f.VarP(newInt64Value(value, p), name, "", usage)
***REMOVED***

// Int64VarP is like Int64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64VarP(p *int64, name, shorthand string, value int64, usage string) ***REMOVED***
	f.VarP(newInt64Value(value, p), name, shorthand, usage)
***REMOVED***

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func Int64Var(p *int64, name string, value int64, usage string) ***REMOVED***
	CommandLine.VarP(newInt64Value(value, p), name, "", usage)
***REMOVED***

// Int64VarP is like Int64Var, but accepts a shorthand letter that can be used after a single dash.
func Int64VarP(p *int64, name, shorthand string, value int64, usage string) ***REMOVED***
	CommandLine.VarP(newInt64Value(value, p), name, shorthand, usage)
***REMOVED***

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func (f *FlagSet) Int64(name string, value int64, usage string) *int64 ***REMOVED***
	p := new(int64)
	f.Int64VarP(p, name, "", value, usage)
	return p
***REMOVED***

// Int64P is like Int64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64P(name, shorthand string, value int64, usage string) *int64 ***REMOVED***
	p := new(int64)
	f.Int64VarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func Int64(name string, value int64, usage string) *int64 ***REMOVED***
	return CommandLine.Int64P(name, "", value, usage)
***REMOVED***

// Int64P is like Int64, but accepts a shorthand letter that can be used after a single dash.
func Int64P(name, shorthand string, value int64, usage string) *int64 ***REMOVED***
	return CommandLine.Int64P(name, shorthand, value, usage)
***REMOVED***
