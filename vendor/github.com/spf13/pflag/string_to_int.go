package pflag

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// -- stringToInt Value
type stringToIntValue struct ***REMOVED***
	value   *map[string]int
	changed bool
***REMOVED***

func newStringToIntValue(val map[string]int, p *map[string]int) *stringToIntValue ***REMOVED***
	ssv := new(stringToIntValue)
	ssv.value = p
	*ssv.value = val
	return ssv
***REMOVED***

// Format: a=1,b=2
func (s *stringToIntValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make(map[string]int, len(ss))
	for _, pair := range ss ***REMOVED***
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 ***REMOVED***
			return fmt.Errorf("%s must be formatted as key=value", pair)
		***REMOVED***
		var err error
		out[kv[0]], err = strconv.Atoi(kv[1])
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if !s.changed ***REMOVED***
		*s.value = out
	***REMOVED*** else ***REMOVED***
		for k, v := range out ***REMOVED***
			(*s.value)[k] = v
		***REMOVED***
	***REMOVED***
	s.changed = true
	return nil
***REMOVED***

func (s *stringToIntValue) Type() string ***REMOVED***
	return "stringToInt"
***REMOVED***

func (s *stringToIntValue) String() string ***REMOVED***
	var buf bytes.Buffer
	i := 0
	for k, v := range *s.value ***REMOVED***
		if i > 0 ***REMOVED***
			buf.WriteRune(',')
		***REMOVED***
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.Itoa(v))
		i++
	***REMOVED***
	return "[" + buf.String() + "]"
***REMOVED***

func stringToIntConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// An empty string would cause an empty map
	if len(val) == 0 ***REMOVED***
		return map[string]int***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make(map[string]int, len(ss))
	for _, pair := range ss ***REMOVED***
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 ***REMOVED***
			return nil, fmt.Errorf("%s must be formatted as key=value", pair)
		***REMOVED***
		var err error
		out[kv[0]], err = strconv.Atoi(kv[1])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***

// GetStringToInt return the map[string]int value of a flag with the given name
func (f *FlagSet) GetStringToInt(name string) (map[string]int, error) ***REMOVED***
	val, err := f.getFlagType(name, "stringToInt", stringToIntConv)
	if err != nil ***REMOVED***
		return map[string]int***REMOVED******REMOVED***, err
	***REMOVED***
	return val.(map[string]int), nil
***REMOVED***

// StringToIntVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a map[string]int variable in which to store the values of the multiple flags.
// The value of each argument will not try to be separated by comma
func (f *FlagSet) StringToIntVar(p *map[string]int, name string, value map[string]int, usage string) ***REMOVED***
	f.VarP(newStringToIntValue(value, p), name, "", usage)
***REMOVED***

// StringToIntVarP is like StringToIntVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringToIntVarP(p *map[string]int, name, shorthand string, value map[string]int, usage string) ***REMOVED***
	f.VarP(newStringToIntValue(value, p), name, shorthand, usage)
***REMOVED***

// StringToIntVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a map[string]int variable in which to store the value of the flag.
// The value of each argument will not try to be separated by comma
func StringToIntVar(p *map[string]int, name string, value map[string]int, usage string) ***REMOVED***
	CommandLine.VarP(newStringToIntValue(value, p), name, "", usage)
***REMOVED***

// StringToIntVarP is like StringToIntVar, but accepts a shorthand letter that can be used after a single dash.
func StringToIntVarP(p *map[string]int, name, shorthand string, value map[string]int, usage string) ***REMOVED***
	CommandLine.VarP(newStringToIntValue(value, p), name, shorthand, usage)
***REMOVED***

// StringToInt defines a string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma
func (f *FlagSet) StringToInt(name string, value map[string]int, usage string) *map[string]int ***REMOVED***
	p := map[string]int***REMOVED******REMOVED***
	f.StringToIntVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// StringToIntP is like StringToInt, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringToIntP(name, shorthand string, value map[string]int, usage string) *map[string]int ***REMOVED***
	p := map[string]int***REMOVED******REMOVED***
	f.StringToIntVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// StringToInt defines a string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma
func StringToInt(name string, value map[string]int, usage string) *map[string]int ***REMOVED***
	return CommandLine.StringToIntP(name, "", value, usage)
***REMOVED***

// StringToIntP is like StringToInt, but accepts a shorthand letter that can be used after a single dash.
func StringToIntP(name, shorthand string, value map[string]int, usage string) *map[string]int ***REMOVED***
	return CommandLine.StringToIntP(name, shorthand, value, usage)
***REMOVED***
