package pflag

import "strconv"

// optional interface to indicate boolean flags that can be
// supplied without "=value" text
type boolFlag interface ***REMOVED***
	Value
	IsBoolFlag() bool
***REMOVED***

// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue ***REMOVED***
	*p = val
	return (*boolValue)(p)
***REMOVED***

func (b *boolValue) Set(s string) error ***REMOVED***
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
***REMOVED***

func (b *boolValue) Type() string ***REMOVED***
	return "bool"
***REMOVED***

func (b *boolValue) String() string ***REMOVED*** return strconv.FormatBool(bool(*b)) ***REMOVED***

func (b *boolValue) IsBoolFlag() bool ***REMOVED*** return true ***REMOVED***

func boolConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return strconv.ParseBool(sval)
***REMOVED***

// GetBool return the bool value of a flag with the given name
func (f *FlagSet) GetBool(name string) (bool, error) ***REMOVED***
	val, err := f.getFlagType(name, "bool", boolConv)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return val.(bool), nil
***REMOVED***

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) ***REMOVED***
	f.BoolVarP(p, name, "", value, usage)
***REMOVED***

// BoolVarP is like BoolVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolVarP(p *bool, name, shorthand string, value bool, usage string) ***REMOVED***
	flag := f.VarPF(newBoolValue(value, p), name, shorthand, usage)
	flag.NoOptDefVal = "true"
***REMOVED***

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func BoolVar(p *bool, name string, value bool, usage string) ***REMOVED***
	BoolVarP(p, name, "", value, usage)
***REMOVED***

// BoolVarP is like BoolVar, but accepts a shorthand letter that can be used after a single dash.
func BoolVarP(p *bool, name, shorthand string, value bool, usage string) ***REMOVED***
	flag := CommandLine.VarPF(newBoolValue(value, p), name, shorthand, usage)
	flag.NoOptDefVal = "true"
***REMOVED***

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func (f *FlagSet) Bool(name string, value bool, usage string) *bool ***REMOVED***
	return f.BoolP(name, "", value, usage)
***REMOVED***

// BoolP is like Bool, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolP(name, shorthand string, value bool, usage string) *bool ***REMOVED***
	p := new(bool)
	f.BoolVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func Bool(name string, value bool, usage string) *bool ***REMOVED***
	return BoolP(name, "", value, usage)
***REMOVED***

// BoolP is like Bool, but accepts a shorthand letter that can be used after a single dash.
func BoolP(name, shorthand string, value bool, usage string) *bool ***REMOVED***
	b := CommandLine.BoolP(name, shorthand, value, usage)
	return b
***REMOVED***
