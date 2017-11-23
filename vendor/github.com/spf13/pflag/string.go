package pflag

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue ***REMOVED***
	*p = val
	return (*stringValue)(p)
***REMOVED***

func (s *stringValue) Set(val string) error ***REMOVED***
	*s = stringValue(val)
	return nil
***REMOVED***
func (s *stringValue) Type() string ***REMOVED***
	return "string"
***REMOVED***

func (s *stringValue) String() string ***REMOVED*** return string(*s) ***REMOVED***

func stringConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	return sval, nil
***REMOVED***

// GetString return the string value of a flag with the given name
func (f *FlagSet) GetString(name string) (string, error) ***REMOVED***
	val, err := f.getFlagType(name, "string", stringConv)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return val.(string), nil
***REMOVED***

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (f *FlagSet) StringVar(p *string, name string, value string, usage string) ***REMOVED***
	f.VarP(newStringValue(value, p), name, "", usage)
***REMOVED***

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringVarP(p *string, name, shorthand string, value string, usage string) ***REMOVED***
	f.VarP(newStringValue(value, p), name, shorthand, usage)
***REMOVED***

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string) ***REMOVED***
	CommandLine.VarP(newStringValue(value, p), name, "", usage)
***REMOVED***

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
func StringVarP(p *string, name, shorthand string, value string, usage string) ***REMOVED***
	CommandLine.VarP(newStringValue(value, p), name, shorthand, usage)
***REMOVED***

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (f *FlagSet) String(name string, value string, usage string) *string ***REMOVED***
	p := new(string)
	f.StringVarP(p, name, "", value, usage)
	return p
***REMOVED***

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringP(name, shorthand string, value string, usage string) *string ***REMOVED***
	p := new(string)
	f.StringVarP(p, name, shorthand, value, usage)
	return p
***REMOVED***

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string) *string ***REMOVED***
	return CommandLine.StringP(name, "", value, usage)
***REMOVED***

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func StringP(name, shorthand string, value string, usage string) *string ***REMOVED***
	return CommandLine.StringP(name, shorthand, value, usage)
***REMOVED***
