package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- intSlice Value
type intSliceValue struct ***REMOVED***
	value   *[]int
	changed bool
***REMOVED***

func newIntSliceValue(val []int, p *[]int) *intSliceValue ***REMOVED***
	isv := new(intSliceValue)
	isv.value = p
	*isv.value = val
	return isv
***REMOVED***

func (s *intSliceValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]int, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = strconv.Atoi(d)
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

func (s *intSliceValue) Type() string ***REMOVED***
	return "intSlice"
***REMOVED***

func (s *intSliceValue) String() string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = fmt.Sprintf("%d", d)
	***REMOVED***
	return "[" + strings.Join(out, ",") + "]"
***REMOVED***

func (s *intSliceValue) Append(val string) error ***REMOVED***
	i, err := strconv.Atoi(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *intSliceValue) Replace(val []string) error ***REMOVED***
	out := make([]int, len(val))
	for i, d := range val ***REMOVED***
		var err error
		out[i], err = strconv.Atoi(d)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	*s.value = out
	return nil
***REMOVED***

func (s *intSliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = strconv.Itoa(d)
	***REMOVED***
	return out
***REMOVED***

func intSliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []int***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]int, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = strconv.Atoi(d)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

	***REMOVED***
	return out, nil
***REMOVED***

// GetIntSlice return the []int value of a flag with the given name
func (f *FlagSet) GetIntSlice(name string) ([]int, error) ***REMOVED***
	val, err := f.getFlagType(name, "intSlice", intSliceConv)
	if err != nil ***REMOVED***
		return []int***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]int), nil
***REMOVED***

// IntSliceVar defines a intSlice flag with specified name, default value, and usage string.
// The argument p points to a []int variable in which to store the value of the flag.
func (f *FlagSet) IntSliceVar(p *[]int, name string, value []int, usage string) ***REMOVED***
	f.VarP(newIntSliceValue(value, p), name, "", usage)
***REMOVED***

// IntSliceVarP is like IntSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntSliceVarP(p *[]int, name, shorthand string, value []int, usage string) ***REMOVED***
	f.VarP(newIntSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// IntSliceVar defines a int[] flag with specified name, default value, and usage string.
// The argument p points to a int[] variable in which to store the value of the flag.
func IntSliceVar(p *[]int, name string, value []int, usage string) ***REMOVED***
	CommandLine.VarP(newIntSliceValue(value, p), name, "", usage)
***REMOVED***

// IntSliceVarP is like IntSliceVar, but accepts a shorthand letter that can be used after a single dash.
func IntSliceVarP(p *[]int, name, shorthand string, value []int, usage string) ***REMOVED***
	CommandLine.VarP(newIntSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// IntSlice defines a []int flag with specified name, default value, and usage string.
// The return value is the address of a []int variable that stores the value of the flag.
func (f *FlagSet) IntSlice(name string, value []int, usage string) *[]int ***REMOVED***
	p := []int***REMOVED******REMOVED***
	f.IntSliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// IntSliceP is like IntSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntSliceP(name, shorthand string, value []int, usage string) *[]int ***REMOVED***
	p := []int***REMOVED******REMOVED***
	f.IntSliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// IntSlice defines a []int flag with specified name, default value, and usage string.
// The return value is the address of a []int variable that stores the value of the flag.
func IntSlice(name string, value []int, usage string) *[]int ***REMOVED***
	return CommandLine.IntSliceP(name, "", value, usage)
***REMOVED***

// IntSliceP is like IntSlice, but accepts a shorthand letter that can be used after a single dash.
func IntSliceP(name, shorthand string, value []int, usage string) *[]int ***REMOVED***
	return CommandLine.IntSliceP(name, shorthand, value, usage)
***REMOVED***
