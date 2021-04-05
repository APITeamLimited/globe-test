package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- int64Slice Value
type int64SliceValue struct ***REMOVED***
	value   *[]int64
	changed bool
***REMOVED***

func newInt64SliceValue(val []int64, p *[]int64) *int64SliceValue ***REMOVED***
	isv := new(int64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
***REMOVED***

func (s *int64SliceValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]int64, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = strconv.ParseInt(d, 0, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

	***REMOVED***
	if !s.changed ***REMOVED***
		*s.value = out
	***REMOVED*** else ***REMOVED***
		*s.value = append(*s.value, out...)
	***REMOVED***
	s.changed = true
	return nil
***REMOVED***

func (s *int64SliceValue) Type() string ***REMOVED***
	return "int64Slice"
***REMOVED***

func (s *int64SliceValue) String() string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = fmt.Sprintf("%d", d)
	***REMOVED***
	return "[" + strings.Join(out, ",") + "]"
***REMOVED***

func (s *int64SliceValue) fromString(val string) (int64, error) ***REMOVED***
	return strconv.ParseInt(val, 0, 64)
***REMOVED***

func (s *int64SliceValue) toString(val int64) string ***REMOVED***
	return fmt.Sprintf("%d", val)
***REMOVED***

func (s *int64SliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *int64SliceValue) Replace(val []string) error ***REMOVED***
	out := make([]int64, len(val))
	for i, d := range val ***REMOVED***
		var err error
		out[i], err = s.fromString(d)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	*s.value = out
	return nil
***REMOVED***

func (s *int64SliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func int64SliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []int64***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]int64, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = strconv.ParseInt(d, 0, 64)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

	***REMOVED***
	return out, nil
***REMOVED***

// GetInt64Slice return the []int64 value of a flag with the given name
func (f *FlagSet) GetInt64Slice(name string) ([]int64, error) ***REMOVED***
	val, err := f.getFlagType(name, "int64Slice", int64SliceConv)
	if err != nil ***REMOVED***
		return []int64***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]int64), nil
***REMOVED***

// Int64SliceVar defines a int64Slice flag with specified name, default value, and usage string.
// The argument p points to a []int64 variable in which to store the value of the flag.
func (f *FlagSet) Int64SliceVar(p *[]int64, name string, value []int64, usage string) ***REMOVED***
	f.VarP(newInt64SliceValue(value, p), name, "", usage)
***REMOVED***

// Int64SliceVarP is like Int64SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64SliceVarP(p *[]int64, name, shorthand string, value []int64, usage string) ***REMOVED***
	f.VarP(newInt64SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Int64SliceVar defines a int64[] flag with specified name, default value, and usage string.
// The argument p points to a int64[] variable in which to store the value of the flag.
func Int64SliceVar(p *[]int64, name string, value []int64, usage string) ***REMOVED***
	CommandLine.VarP(newInt64SliceValue(value, p), name, "", usage)
***REMOVED***

// Int64SliceVarP is like Int64SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Int64SliceVarP(p *[]int64, name, shorthand string, value []int64, usage string) ***REMOVED***
	CommandLine.VarP(newInt64SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func (f *FlagSet) Int64Slice(name string, value []int64, usage string) *[]int64 ***REMOVED***
	p := []int64***REMOVED******REMOVED***
	f.Int64SliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// Int64SliceP is like Int64Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64SliceP(name, shorthand string, value []int64, usage string) *[]int64 ***REMOVED***
	p := []int64***REMOVED******REMOVED***
	f.Int64SliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func Int64Slice(name string, value []int64, usage string) *[]int64 ***REMOVED***
	return CommandLine.Int64SliceP(name, "", value, usage)
***REMOVED***

// Int64SliceP is like Int64Slice, but accepts a shorthand letter that can be used after a single dash.
func Int64SliceP(name, shorthand string, value []int64, usage string) *[]int64 ***REMOVED***
	return CommandLine.Int64SliceP(name, shorthand, value, usage)
***REMOVED***
