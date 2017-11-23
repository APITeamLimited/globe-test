package pflag

import "strconv"

// -- int32 Value
type int32Value int32

func newInt32Value(val int32, p *int32) *int32Value ***REMOVED***
	*p = val
	return (*int32Value)(p)
***REMOVED***

func (i *int32Value) Set(s string) error ***REMOVED***
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
***REMOVED***

func (i *int32Value) Type() string ***REMOVED***
	return "int32"
***REMOVED***

func (i *int32Value) String() string ***REMOVED*** return strconv.FormatInt(int64(*i), 10) ***REMOVED***

func int32Conv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v, err := strconv.ParseInt(sval, 0, 32)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return int32(v), nil
***REMOVED***

// GetInt32 return the int32 value of a flag with the given name
func (f *FlagSet) GetInt32(name string) (int32, error) ***REMOVED***
	val, err := f.getFlagType(name, "int32", int32Conv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(int32), nil
***REMOVED***

// Int32Var defines an int32 flag with specified name, default value, and usage string.
// The argument p points to an int32 variable in which to store the value of the flag.
func (f *FlagSet) Int32Var(p *int32, name string, value int32, usage string) ***REMOVED***
	f.VarP(newInt32Value(value, p), name, "", usage)
***REMOVED***

// Int32VarP is like Int32Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int32VarP(p *int32, name, shorthand string, value int32, usage string) ***REMOVED***
	f.VarP(newInt32Value(value, p), name, shorthand, usage)
***REMOVED***

// Int32Var defines an int32 flag with specified name, default value, and usage string.
// The argument p points to an int32 variable in which to store the value of the flag.
func Int32Var(p *int32, name string, value int32, usage string) ***REMOVED***
	CommandLine.VarP(newInt32Value(value, p), name, "", usage)
***REMOVED***

// Int32VarP is like Int32Var, but accepts a shorthand letter that can be used after a single dash.
func Int32VarP(p *int32, name, shorthand string, value int32, usage string) ***REMOVED***
	CommandLine.VarP(newInt32Value(value, p), name, shorthand, usage)
***REMOVED***

// Int32 defines an int32 flag with specified name, default value, and usage string.
// The return value is the address of an int32 variable that stores the value of the flag.
func (f *FlagSet) Int32(name string, value int32, usage string) *int32 ***REMOVED***
	p := new(int32)
	f.Int32VarP(p, name, "", value, usage)
	return p
***REMOVED***

// Int32P is like Int32, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int32P(name, shorthand string, value int32, usage string) *int32 ***REMOVED***
	p := new(int32)
	f.Int32VarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Int32 defines an int32 flag with specified name, default value, and usage string.
// The return value is the address of an int32 variable that stores the value of the flag.
func Int32(name string, value int32, usage string) *int32 ***REMOVED***
	return CommandLine.Int32P(name, "", value, usage)
***REMOVED***

// Int32P is like Int32, but accepts a shorthand letter that can be used after a single dash.
func Int32P(name, shorthand string, value int32, usage string) *int32 ***REMOVED***
	return CommandLine.Int32P(name, shorthand, value, usage)
***REMOVED***
