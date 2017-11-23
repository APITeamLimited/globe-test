package pflag

import "strconv"

// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue ***REMOVED***
	*p = val
	return (*intValue)(p)
***REMOVED***

func (i *intValue) Set(s string) error ***REMOVED***
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
***REMOVED***

func (i *intValue) Type() string ***REMOVED***
	return "int"
***REMOVED***

func (i *intValue) String() string ***REMOVED*** return strconv.Itoa(int(*i)) ***REMOVED***

func intConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return strconv.Atoi(sval)
***REMOVED***

// GetInt return the int value of a flag with the given name
func (f *FlagSet) GetInt(name string) (int, error) ***REMOVED***
	val, err := f.getFlagType(name, "int", intConv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(int), nil
***REMOVED***

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) ***REMOVED***
	f.VarP(newIntValue(value, p), name, "", usage)
***REMOVED***

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntVarP(p *int, name, shorthand string, value int, usage string) ***REMOVED***
	f.VarP(newIntValue(value, p), name, shorthand, usage)
***REMOVED***

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string) ***REMOVED***
	CommandLine.VarP(newIntValue(value, p), name, "", usage)
***REMOVED***

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func IntVarP(p *int, name, shorthand string, value int, usage string) ***REMOVED***
	CommandLine.VarP(newIntValue(value, p), name, shorthand, usage)
***REMOVED***

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (f *FlagSet) Int(name string, value int, usage string) *int ***REMOVED***
	p := new(int)
	f.IntVarP(p, name, "", value, usage)
	return p
***REMOVED***

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntP(name, shorthand string, value int, usage string) *int ***REMOVED***
	p := new(int)
	f.IntVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string) *int ***REMOVED***
	return CommandLine.IntP(name, "", value, usage)
***REMOVED***

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func IntP(name, shorthand string, value int, usage string) *int ***REMOVED***
	return CommandLine.IntP(name, shorthand, value, usage)
***REMOVED***
