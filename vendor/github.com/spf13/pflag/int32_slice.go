package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- int32Slice Value
type int32SliceValue struct ***REMOVED***
	value   *[]int32
	changed bool
***REMOVED***

func newInt32SliceValue(val []int32, p *[]int32) *int32SliceValue ***REMOVED***
	isv := new(int32SliceValue)
	isv.value = p
	*isv.value = val
	return isv
***REMOVED***

func (s *int32SliceValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]int32, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		var temp64 int64
		temp64, err = strconv.ParseInt(d, 0, 32)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		out[i] = int32(temp64)

	***REMOVED***
	if !s.changed ***REMOVED***
		*s.value = out
	***REMOVED*** else ***REMOVED***
		*s.value = append(*s.value, out...)
	***REMOVED***
	s.changed = true
	return nil
***REMOVED***

func (s *int32SliceValue) Type() string ***REMOVED***
	return "int32Slice"
***REMOVED***

func (s *int32SliceValue) String() string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = fmt.Sprintf("%d", d)
	***REMOVED***
	return "[" + strings.Join(out, ",") + "]"
***REMOVED***

func (s *int32SliceValue) fromString(val string) (int32, error) ***REMOVED***
	t64, err := strconv.ParseInt(val, 0, 32)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return int32(t64), nil
***REMOVED***

func (s *int32SliceValue) toString(val int32) string ***REMOVED***
	return fmt.Sprintf("%d", val)
***REMOVED***

func (s *int32SliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *int32SliceValue) Replace(val []string) error ***REMOVED***
	out := make([]int32, len(val))
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

func (s *int32SliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func int32SliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []int32***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]int32, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		var temp64 int64
		temp64, err = strconv.ParseInt(d, 0, 32)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out[i] = int32(temp64)

	***REMOVED***
	return out, nil
***REMOVED***

// GetInt32Slice return the []int32 value of a flag with the given name
func (f *FlagSet) GetInt32Slice(name string) ([]int32, error) ***REMOVED***
	val, err := f.getFlagType(name, "int32Slice", int32SliceConv)
	if err != nil ***REMOVED***
		return []int32***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]int32), nil
***REMOVED***

// Int32SliceVar defines a int32Slice flag with specified name, default value, and usage string.
// The argument p points to a []int32 variable in which to store the value of the flag.
func (f *FlagSet) Int32SliceVar(p *[]int32, name string, value []int32, usage string) ***REMOVED***
	f.VarP(newInt32SliceValue(value, p), name, "", usage)
***REMOVED***

// Int32SliceVarP is like Int32SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int32SliceVarP(p *[]int32, name, shorthand string, value []int32, usage string) ***REMOVED***
	f.VarP(newInt32SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Int32SliceVar defines a int32[] flag with specified name, default value, and usage string.
// The argument p points to a int32[] variable in which to store the value of the flag.
func Int32SliceVar(p *[]int32, name string, value []int32, usage string) ***REMOVED***
	CommandLine.VarP(newInt32SliceValue(value, p), name, "", usage)
***REMOVED***

// Int32SliceVarP is like Int32SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Int32SliceVarP(p *[]int32, name, shorthand string, value []int32, usage string) ***REMOVED***
	CommandLine.VarP(newInt32SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Int32Slice defines a []int32 flag with specified name, default value, and usage string.
// The return value is the address of a []int32 variable that stores the value of the flag.
func (f *FlagSet) Int32Slice(name string, value []int32, usage string) *[]int32 ***REMOVED***
	p := []int32***REMOVED******REMOVED***
	f.Int32SliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// Int32SliceP is like Int32Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int32SliceP(name, shorthand string, value []int32, usage string) *[]int32 ***REMOVED***
	p := []int32***REMOVED******REMOVED***
	f.Int32SliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// Int32Slice defines a []int32 flag with specified name, default value, and usage string.
// The return value is the address of a []int32 variable that stores the value of the flag.
func Int32Slice(name string, value []int32, usage string) *[]int32 ***REMOVED***
	return CommandLine.Int32SliceP(name, "", value, usage)
***REMOVED***

// Int32SliceP is like Int32Slice, but accepts a shorthand letter that can be used after a single dash.
func Int32SliceP(name, shorthand string, value []int32, usage string) *[]int32 ***REMOVED***
	return CommandLine.Int32SliceP(name, shorthand, value, usage)
***REMOVED***
