package pflag

import "strconv"

// -- float32 Value
type float32Value float32

func newFloat32Value(val float32, p *float32) *float32Value ***REMOVED***
	*p = val
	return (*float32Value)(p)
***REMOVED***

func (f *float32Value) Set(s string) error ***REMOVED***
	v, err := strconv.ParseFloat(s, 32)
	*f = float32Value(v)
	return err
***REMOVED***

func (f *float32Value) Type() string ***REMOVED***
	return "float32"
***REMOVED***

func (f *float32Value) String() string ***REMOVED*** return strconv.FormatFloat(float64(*f), 'g', -1, 32) ***REMOVED***

func float32Conv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v, err := strconv.ParseFloat(sval, 32)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return float32(v), nil
***REMOVED***

// GetFloat32 return the float32 value of a flag with the given name
func (f *FlagSet) GetFloat32(name string) (float32, error) ***REMOVED***
	val, err := f.getFlagType(name, "float32", float32Conv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(float32), nil
***REMOVED***

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func (f *FlagSet) Float32Var(p *float32, name string, value float32, usage string) ***REMOVED***
	f.VarP(newFloat32Value(value, p), name, "", usage)
***REMOVED***

// Float32VarP is like Float32Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32VarP(p *float32, name, shorthand string, value float32, usage string) ***REMOVED***
	f.VarP(newFloat32Value(value, p), name, shorthand, usage)
***REMOVED***

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func Float32Var(p *float32, name string, value float32, usage string) ***REMOVED***
	CommandLine.VarP(newFloat32Value(value, p), name, "", usage)
***REMOVED***

// Float32VarP is like Float32Var, but accepts a shorthand letter that can be used after a single dash.
func Float32VarP(p *float32, name, shorthand string, value float32, usage string) ***REMOVED***
	CommandLine.VarP(newFloat32Value(value, p), name, shorthand, usage)
***REMOVED***

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func (f *FlagSet) Float32(name string, value float32, usage string) *float32 ***REMOVED***
	p := new(float32)
	f.Float32VarP(p, name, "", value, usage)
	return p
***REMOVED***

// Float32P is like Float32, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32P(name, shorthand string, value float32, usage string) *float32 ***REMOVED***
	p := new(float32)
	f.Float32VarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func Float32(name string, value float32, usage string) *float32 ***REMOVED***
	return CommandLine.Float32P(name, "", value, usage)
***REMOVED***

// Float32P is like Float32, but accepts a shorthand letter that can be used after a single dash.
func Float32P(name, shorthand string, value float32, usage string) *float32 ***REMOVED***
	return CommandLine.Float32P(name, shorthand, value, usage)
***REMOVED***
