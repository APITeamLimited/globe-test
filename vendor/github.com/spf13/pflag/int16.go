package pflag

import "strconv"

// -- int16 Value
type int16Value int16

func newInt16Value(val int16, p *int16) *int16Value ***REMOVED***
	*p = val
	return (*int16Value)(p)
***REMOVED***

func (i *int16Value) Set(s string) error ***REMOVED***
	v, err := strconv.ParseInt(s, 0, 16)
	*i = int16Value(v)
	return err
***REMOVED***

func (i *int16Value) Type() string ***REMOVED***
	return "int16"
***REMOVED***

func (i *int16Value) String() string ***REMOVED*** return strconv.FormatInt(int64(*i), 10) ***REMOVED***

func int16Conv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v, err := strconv.ParseInt(sval, 0, 16)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return int16(v), nil
***REMOVED***

// GetInt16 returns the int16 value of a flag with the given name
func (f *FlagSet) GetInt16(name string) (int16, error) ***REMOVED***
	val, err := f.getFlagType(name, "int16", int16Conv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(int16), nil
***REMOVED***

// Int16Var defines an int16 flag with specified name, default value, and usage string.
// The argument p points to an int16 variable in which to store the value of the flag.
func (f *FlagSet) Int16Var(p *int16, name string, value int16, usage string) ***REMOVED***
	f.VarP(newInt16Value(value, p), name, "", usage)
***REMOVED***

// Int16VarP is like Int16Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int16VarP(p *int16, name, shorthand string, value int16, usage string) ***REMOVED***
	f.VarP(newInt16Value(value, p), name, shorthand, usage)
***REMOVED***

// Int16Var defines an int16 flag with specified name, default value, and usage string.
// The argument p points to an int16 variable in which to store the value of the flag.
func Int16Var(p *int16, name string, value int16, usage string) ***REMOVED***
	CommandLine.VarP(newInt16Value(value, p), name, "", usage)
***REMOVED***

// Int16VarP is like Int16Var, but accepts a shorthand letter that can be used after a single dash.
func Int16VarP(p *int16, name, shorthand string, value int16, usage string) ***REMOVED***
	CommandLine.VarP(newInt16Value(value, p), name, shorthand, usage)
***REMOVED***

// Int16 defines an int16 flag with specified name, default value, and usage string.
// The return value is the address of an int16 variable that stores the value of the flag.
func (f *FlagSet) Int16(name string, value int16, usage string) *int16 ***REMOVED***
	p := new(int16)
	f.Int16VarP(p, name, "", value, usage)
	return p
***REMOVED***

// Int16P is like Int16, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int16P(name, shorthand string, value int16, usage string) *int16 ***REMOVED***
	p := new(int16)
	f.Int16VarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Int16 defines an int16 flag with specified name, default value, and usage string.
// The return value is the address of an int16 variable that stores the value of the flag.
func Int16(name string, value int16, usage string) *int16 ***REMOVED***
	return CommandLine.Int16P(name, "", value, usage)
***REMOVED***

// Int16P is like Int16, but accepts a shorthand letter that can be used after a single dash.
func Int16P(name, shorthand string, value int16, usage string) *int16 ***REMOVED***
	return CommandLine.Int16P(name, shorthand, value, usage)
***REMOVED***
