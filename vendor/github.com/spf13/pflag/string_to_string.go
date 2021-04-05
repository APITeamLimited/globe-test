package pflag

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

// -- stringToString Value
type stringToStringValue struct ***REMOVED***
	value   *map[string]string
	changed bool
***REMOVED***

func newStringToStringValue(val map[string]string, p *map[string]string) *stringToStringValue ***REMOVED***
	ssv := new(stringToStringValue)
	ssv.value = p
	*ssv.value = val
	return ssv
***REMOVED***

// Format: a=1,b=2
func (s *stringToStringValue) Set(val string) error ***REMOVED***
	var ss []string
	n := strings.Count(val, "=")
	switch n ***REMOVED***
	case 0:
		return fmt.Errorf("%s must be formatted as key=value", val)
	case 1:
		ss = append(ss, strings.Trim(val, `"`))
	default:
		r := csv.NewReader(strings.NewReader(val))
		var err error
		ss, err = r.Read()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	out := make(map[string]string, len(ss))
	for _, pair := range ss ***REMOVED***
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 ***REMOVED***
			return fmt.Errorf("%s must be formatted as key=value", pair)
		***REMOVED***
		out[kv[0]] = kv[1]
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

func (s *stringToStringValue) Type() string ***REMOVED***
	return "stringToString"
***REMOVED***

func (s *stringToStringValue) String() string ***REMOVED***
	records := make([]string, 0, len(*s.value)>>1)
	for k, v := range *s.value ***REMOVED***
		records = append(records, k+"="+v)
	***REMOVED***

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(records); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	w.Flush()
	return "[" + strings.TrimSpace(buf.String()) + "]"
***REMOVED***

func stringToStringConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// An empty string would cause an empty map
	if len(val) == 0 ***REMOVED***
		return map[string]string***REMOVED******REMOVED***, nil
	***REMOVED***
	r := csv.NewReader(strings.NewReader(val))
	ss, err := r.Read()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	out := make(map[string]string, len(ss))
	for _, pair := range ss ***REMOVED***
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 ***REMOVED***
			return nil, fmt.Errorf("%s must be formatted as key=value", pair)
		***REMOVED***
		out[kv[0]] = kv[1]
	***REMOVED***
	return out, nil
***REMOVED***

// GetStringToString return the map[string]string value of a flag with the given name
func (f *FlagSet) GetStringToString(name string) (map[string]string, error) ***REMOVED***
	val, err := f.getFlagType(name, "stringToString", stringToStringConv)
	if err != nil ***REMOVED***
		return map[string]string***REMOVED******REMOVED***, err
	***REMOVED***
	return val.(map[string]string), nil
***REMOVED***

// StringToStringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a map[string]string variable in which to store the values of the multiple flags.
// The value of each argument will not try to be separated by comma
func (f *FlagSet) StringToStringVar(p *map[string]string, name string, value map[string]string, usage string) ***REMOVED***
	f.VarP(newStringToStringValue(value, p), name, "", usage)
***REMOVED***

// StringToStringVarP is like StringToStringVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringToStringVarP(p *map[string]string, name, shorthand string, value map[string]string, usage string) ***REMOVED***
	f.VarP(newStringToStringValue(value, p), name, shorthand, usage)
***REMOVED***

// StringToStringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a map[string]string variable in which to store the value of the flag.
// The value of each argument will not try to be separated by comma
func StringToStringVar(p *map[string]string, name string, value map[string]string, usage string) ***REMOVED***
	CommandLine.VarP(newStringToStringValue(value, p), name, "", usage)
***REMOVED***

// StringToStringVarP is like StringToStringVar, but accepts a shorthand letter that can be used after a single dash.
func StringToStringVarP(p *map[string]string, name, shorthand string, value map[string]string, usage string) ***REMOVED***
	CommandLine.VarP(newStringToStringValue(value, p), name, shorthand, usage)
***REMOVED***

// StringToString defines a string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]string variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma
func (f *FlagSet) StringToString(name string, value map[string]string, usage string) *map[string]string ***REMOVED***
	p := map[string]string***REMOVED******REMOVED***
	f.StringToStringVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// StringToStringP is like StringToString, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringToStringP(name, shorthand string, value map[string]string, usage string) *map[string]string ***REMOVED***
	p := map[string]string***REMOVED******REMOVED***
	f.StringToStringVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// StringToString defines a string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]string variable that stores the value of the flag.
// The value of each argument will not try to be separated by comma
func StringToString(name string, value map[string]string, usage string) *map[string]string ***REMOVED***
	return CommandLine.StringToStringP(name, "", value, usage)
***REMOVED***

// StringToStringP is like StringToString, but accepts a shorthand letter that can be used after a single dash.
func StringToStringP(name, shorthand string, value map[string]string, usage string) *map[string]string ***REMOVED***
	return CommandLine.StringToStringP(name, shorthand, value, usage)
***REMOVED***
