package pflag

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// -- stringToInt64 Value
type stringToInt64Value struct ***REMOVED***
	value   *map[string]int64
	changed bool
***REMOVED***

func newStringToInt64Value(val map[string]int64, p *map[string]int64) *stringToInt64Value ***REMOVED***
	ssv := new(stringToInt64Value)
	ssv.value = p
	*ssv.value = val
	return ssv
***REMOVED***

// Format: a=1,b=2
func (s *stringToInt64Value) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make(map[string]int64, len(ss))
	for _, pair := range ss ***REMOVED***
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 ***REMOVED***
			return fmt.Errorf("%s must be formatted as key=value", pair)
		***REMOVED***
		var err error
		out[kv[0]], err = strconv.ParseInt(kv[1], 10, 64)
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

func (s *stringToInt64Value) Type() string ***REMOVED***
	return "stringToInt64"
***REMOVED***

func (s *stringToInt64Value) String() string ***REMOVED***
	var buf bytes.Buffer
	i := 0
	for k, v := range *s.value ***REMOVED***
		if i > 0 ***REMOVED***
			buf.WriteRune(',')
		***REMOVED***
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.FormatInt(v, 10))
		i++
	***REMOVED***
	return "[" + buf.String() + "]"
***REMOVED***

func stringToInt64Conv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// An empty string would cause an empty map
	if len(val) == 0 ***REMOVED***
		return map[string]int64***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make(map[string]int64, len(ss))
	for _, pair := range ss ***REMOVED***
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 ***REMOVED***
			return nil, fmt.Errorf("%s must be formatted as key=value", pair)
		***REMOVED***
		var err error
		out[kv[0]], err = strconv.ParseInt(kv[1], 10, 64)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***

// GetStringToInt64 return the map[string]int64 value of a flag with the given name
func (f *FlagSet) GetStringToInt64(name string) (map[string]int64, error) ***REMOVED***
	val, err := f.getFlagType(name, "stringToInt64", stringToInt64Conv)
	if err != nil ***REMOVED***
		return map[string]int64***REMOVED******REMOVED***, err
	***REMOVED***
	return val.(map[string]int64), nil
***REMOVED***

// StringToInt64Var defines a string flag with specified name, default value, and usage string.
// The argument p point64s to a map[string]int64 variable in which to store the values of the multiple flags.
// The value of each argument will not try to be separated by comma
func (f *FlagSet) StringToInt64Var(p *map[string]int64, name string, value map[string]int64, usage string) ***REMOVED***
	f.VarP(newStringToInt64Value(value, p), name, "", usage)
***REMOVED***

// StringToInt64VarP is like StringToInt64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringToInt64VarP(p *map[string]int64, name, shorthand string, value map[string]int64, usage string) ***REMOVED***
	f.VarP(newStringToInt64Value(value, p), name, shorthand, usage)
***REMOVED***

// StringToInt64Var defines a string flag with specified name, default value, and usage string.
// The argument p point64s to a map[string]int64 variable in which to store the value of the flag.
// The value of each argument will not try to be separated by comma
func StringToInt64Var(p *map[string]int64, name string, value map[string]int64, usage string) ***REMOVED***
	CommandLine.VarP(newStringToInt64Value(value, p), name, "", usage)
***REMOVED***

// StringToInt64VarP is like StringToInt64Var, but accepts a shorthand letter that can be used after a single dash.
func StringToInt64VarP(p *map[string]int64, name, shorthand string, value map[string]int64, usage string) ***REMOVED***
	CommandLine.VarP(newStringToInt64Value(value, p), name, shorthand, usage)
***REMOVED***

// StringToInt64 defines a string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int64 variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma
func (f *FlagSet) StringToInt64(name string, value map[string]int64, usage string) *map[string]int64 ***REMOVED***
	p := map[string]int64***REMOVED******REMOVED***
	f.StringToInt64VarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// StringToInt64P is like StringToInt64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringToInt64P(name, shorthand string, value map[string]int64, usage string) *map[string]int64 ***REMOVED***
	p := map[string]int64***REMOVED******REMOVED***
	f.StringToInt64VarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// StringToInt64 defines a string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int64 variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma
func StringToInt64(name string, value map[string]int64, usage string) *map[string]int64 ***REMOVED***
	return CommandLine.StringToInt64P(name, "", value, usage)
***REMOVED***

// StringToInt64P is like StringToInt64, but accepts a shorthand letter that can be used after a single dash.
func StringToInt64P(name, shorthand string, value map[string]int64, usage string) *map[string]int64 ***REMOVED***
	return CommandLine.StringToInt64P(name, shorthand, value, usage)
***REMOVED***
