package pflag

import "strconv"

// -- count Value
type countValue int

func newCountValue(val int, p *int) *countValue ***REMOVED***
	*p = val
	return (*countValue)(p)
***REMOVED***

func (i *countValue) Set(s string) error ***REMOVED***
	// "+1" means that no specific value was passed, so increment
	if s == "+1" ***REMOVED***
		*i = countValue(*i + 1)
		return nil
	***REMOVED***
	v, err := strconv.ParseInt(s, 0, 0)
	*i = countValue(v)
	return err
***REMOVED***

func (i *countValue) Type() string ***REMOVED***
	return "count"
***REMOVED***

func (i *countValue) String() string ***REMOVED*** return strconv.Itoa(int(*i)) ***REMOVED***

func countConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	i, err := strconv.Atoi(sval)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return i, nil
***REMOVED***

// GetCount return the int value of a flag with the given name
func (f *FlagSet) GetCount(name string) (int, error) ***REMOVED***
	val, err := f.getFlagType(name, "count", countConv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return val.(int), nil
***REMOVED***

// CountVar defines a count flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
// A count flag will add 1 to its value every time it is found on the command line
func (f *FlagSet) CountVar(p *int, name string, usage string) ***REMOVED***
	f.CountVarP(p, name, "", usage)
***REMOVED***

// CountVarP is like CountVar only take a shorthand for the flag name.
func (f *FlagSet) CountVarP(p *int, name, shorthand string, usage string) ***REMOVED***
	flag := f.VarPF(newCountValue(0, p), name, shorthand, usage)
	flag.NoOptDefVal = "+1"
***REMOVED***

// CountVar like CountVar only the flag is placed on the CommandLine instead of a given flag set
func CountVar(p *int, name string, usage string) ***REMOVED***
	CommandLine.CountVar(p, name, usage)
***REMOVED***

// CountVarP is like CountVar only take a shorthand for the flag name.
func CountVarP(p *int, name, shorthand string, usage string) ***REMOVED***
	CommandLine.CountVarP(p, name, shorthand, usage)
***REMOVED***

// Count defines a count flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// A count flag will add 1 to its value every time it is found on the command line
func (f *FlagSet) Count(name string, usage string) *int ***REMOVED***
	p := new(int)
	f.CountVarP(p, name, "", usage)
	return p
***REMOVED***

// CountP is like Count only takes a shorthand for the flag name.
func (f *FlagSet) CountP(name, shorthand string, usage string) *int ***REMOVED***
	p := new(int)
	f.CountVarP(p, name, shorthand, usage)
	return p
***REMOVED***

// Count defines a count flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// A count flag will add 1 to its value evey time it is found on the command line
func Count(name string, usage string) *int ***REMOVED***
	return CommandLine.CountP(name, "", usage)
***REMOVED***

// CountP is like Count only takes a shorthand for the flag name.
func CountP(name, shorthand string, usage string) *int ***REMOVED***
	return CommandLine.CountP(name, shorthand, usage)
***REMOVED***
