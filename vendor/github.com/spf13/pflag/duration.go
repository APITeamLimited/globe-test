package pflag

import (
	"time"
)

// -- time.Duration Value
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue ***REMOVED***
	*p = val
	return (*durationValue)(p)
***REMOVED***

func (d *durationValue) Set(s string) error ***REMOVED***
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
***REMOVED***

func (d *durationValue) Type() string ***REMOVED***
	return "duration"
***REMOVED***

func (d *durationValue) String() string ***REMOVED*** return (*time.Duration)(d).String() ***REMOVED***

func durationConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return time.ParseDuration(sval)
***REMOVED***

// GetDuration return the duration value of a flag with the given name
func (f *FlagSet) GetDuration(name string) (time.Duration, error) ***REMOVED***
	val, err := f.getFlagType(name, "duration", durationConv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(time.Duration), nil
***REMOVED***

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) ***REMOVED***
	f.VarP(newDurationValue(value, p), name, "", usage)
***REMOVED***

// DurationVarP is like DurationVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) DurationVarP(p *time.Duration, name, shorthand string, value time.Duration, usage string) ***REMOVED***
	f.VarP(newDurationValue(value, p), name, shorthand, usage)
***REMOVED***

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func DurationVar(p *time.Duration, name string, value time.Duration, usage string) ***REMOVED***
	CommandLine.VarP(newDurationValue(value, p), name, "", usage)
***REMOVED***

// DurationVarP is like DurationVar, but accepts a shorthand letter that can be used after a single dash.
func DurationVarP(p *time.Duration, name, shorthand string, value time.Duration, usage string) ***REMOVED***
	CommandLine.VarP(newDurationValue(value, p), name, shorthand, usage)
***REMOVED***

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration ***REMOVED***
	p := new(time.Duration)
	f.DurationVarP(p, name, "", value, usage)
	return p
***REMOVED***

// DurationP is like Duration, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) DurationP(name, shorthand string, value time.Duration, usage string) *time.Duration ***REMOVED***
	p := new(time.Duration)
	f.DurationVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func Duration(name string, value time.Duration, usage string) *time.Duration ***REMOVED***
	return CommandLine.DurationP(name, "", value, usage)
***REMOVED***

// DurationP is like Duration, but accepts a shorthand letter that can be used after a single dash.
func DurationP(name, shorthand string, value time.Duration, usage string) *time.Duration ***REMOVED***
	return CommandLine.DurationP(name, shorthand, value, usage)
***REMOVED***
