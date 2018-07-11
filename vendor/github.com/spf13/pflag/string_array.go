package pflag

// -- stringArray Value
type stringArrayValue struct ***REMOVED***
	value   *[]string
	changed bool
***REMOVED***

func newStringArrayValue(val []string, p *[]string) *stringArrayValue ***REMOVED***
	ssv := new(stringArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
***REMOVED***

func (s *stringArrayValue) Set(val string) error ***REMOVED***
	if !s.changed ***REMOVED***
		*s.value = []string***REMOVED***val***REMOVED***
		s.changed = true
	***REMOVED*** else ***REMOVED***
		*s.value = append(*s.value, val)
	***REMOVED***
	return nil
***REMOVED***

func (s *stringArrayValue) Type() string ***REMOVED***
	return "stringArray"
***REMOVED***

func (s *stringArrayValue) String() string ***REMOVED***
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
***REMOVED***

func stringArrayConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	sval = sval[1 : len(sval)-1]
	// An empty string would cause a array with one (empty) string
	if len(sval) == 0 ***REMOVED***
		return []string***REMOVED******REMOVED***, nil
	***REMOVED***
	return readAsCSV(sval)
***REMOVED***

// GetStringArray return the []string value of a flag with the given name
func (f *FlagSet) GetStringArray(name string) ([]string, error) ***REMOVED***
	val, err := f.getFlagType(name, "stringArray", stringArrayConv)
	if err != nil ***REMOVED***
		return []string***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]string), nil
***REMOVED***

// StringArrayVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the values of the multiple flags.
// The value of each argument will not try to be separated by comma. Use a StringSlice for that.
func (f *FlagSet) StringArrayVar(p *[]string, name string, value []string, usage string) ***REMOVED***
	f.VarP(newStringArrayValue(value, p), name, "", usage)
***REMOVED***

// StringArrayVarP is like StringArrayVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringArrayVarP(p *[]string, name, shorthand string, value []string, usage string) ***REMOVED***
	f.VarP(newStringArrayValue(value, p), name, shorthand, usage)
***REMOVED***

// StringArrayVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
// The value of each argument will not try to be separated by comma. Use a StringSlice for that.
func StringArrayVar(p *[]string, name string, value []string, usage string) ***REMOVED***
	CommandLine.VarP(newStringArrayValue(value, p), name, "", usage)
***REMOVED***

// StringArrayVarP is like StringArrayVar, but accepts a shorthand letter that can be used after a single dash.
func StringArrayVarP(p *[]string, name, shorthand string, value []string, usage string) ***REMOVED***
	CommandLine.VarP(newStringArrayValue(value, p), name, shorthand, usage)
***REMOVED***

// StringArray defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma. Use a StringSlice for that.
func (f *FlagSet) StringArray(name string, value []string, usage string) *[]string ***REMOVED***
	p := []string***REMOVED******REMOVED***
	f.StringArrayVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// StringArrayP is like StringArray, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringArrayP(name, shorthand string, value []string, usage string) *[]string ***REMOVED***
	p := []string***REMOVED******REMOVED***
	f.StringArrayVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// StringArray defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma. Use a StringSlice for that.
func StringArray(name string, value []string, usage string) *[]string ***REMOVED***
	return CommandLine.StringArrayP(name, "", value, usage)
***REMOVED***

// StringArrayP is like StringArray, but accepts a shorthand letter that can be used after a single dash.
func StringArrayP(name, shorthand string, value []string, usage string) *[]string ***REMOVED***
	return CommandLine.StringArrayP(name, shorthand, value, usage)
***REMOVED***
