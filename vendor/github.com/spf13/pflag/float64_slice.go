package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- float64Slice Value
type float64SliceValue struct ***REMOVED***
	value   *[]float64
	changed bool
***REMOVED***

func newFloat64SliceValue(val []float64, p *[]float64) *float64SliceValue ***REMOVED***
	isv := new(float64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
***REMOVED***

func (s *float64SliceValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]float64, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = strconv.ParseFloat(d, 64)
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

func (s *float64SliceValue) Type() string ***REMOVED***
	return "float64Slice"
***REMOVED***

func (s *float64SliceValue) String() string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = fmt.Sprintf("%f", d)
	***REMOVED***
	return "[" + strings.Join(out, ",") + "]"
***REMOVED***

func (s *float64SliceValue) fromString(val string) (float64, error) ***REMOVED***
	return strconv.ParseFloat(val, 64)
***REMOVED***

func (s *float64SliceValue) toString(val float64) string ***REMOVED***
	return fmt.Sprintf("%f", val)
***REMOVED***

func (s *float64SliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *float64SliceValue) Replace(val []string) error ***REMOVED***
	out := make([]float64, len(val))
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

func (s *float64SliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func float64SliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []float64***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]float64, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = strconv.ParseFloat(d, 64)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

	***REMOVED***
	return out, nil
***REMOVED***

// GetFloat64Slice return the []float64 value of a flag with the given name
func (f *FlagSet) GetFloat64Slice(name string) ([]float64, error) ***REMOVED***
	val, err := f.getFlagType(name, "float64Slice", float64SliceConv)
	if err != nil ***REMOVED***
		return []float64***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]float64), nil
***REMOVED***

// Float64SliceVar defines a float64Slice flag with specified name, default value, and usage string.
// The argument p points to a []float64 variable in which to store the value of the flag.
func (f *FlagSet) Float64SliceVar(p *[]float64, name string, value []float64, usage string) ***REMOVED***
	f.VarP(newFloat64SliceValue(value, p), name, "", usage)
***REMOVED***

// Float64SliceVarP is like Float64SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float64SliceVarP(p *[]float64, name, shorthand string, value []float64, usage string) ***REMOVED***
	f.VarP(newFloat64SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Float64SliceVar defines a float64[] flag with specified name, default value, and usage string.
// The argument p points to a float64[] variable in which to store the value of the flag.
func Float64SliceVar(p *[]float64, name string, value []float64, usage string) ***REMOVED***
	CommandLine.VarP(newFloat64SliceValue(value, p), name, "", usage)
***REMOVED***

// Float64SliceVarP is like Float64SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Float64SliceVarP(p *[]float64, name, shorthand string, value []float64, usage string) ***REMOVED***
	CommandLine.VarP(newFloat64SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Float64Slice defines a []float64 flag with specified name, default value, and usage string.
// The return value is the address of a []float64 variable that stores the value of the flag.
func (f *FlagSet) Float64Slice(name string, value []float64, usage string) *[]float64 ***REMOVED***
	p := []float64***REMOVED******REMOVED***
	f.Float64SliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// Float64SliceP is like Float64Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float64SliceP(name, shorthand string, value []float64, usage string) *[]float64 ***REMOVED***
	p := []float64***REMOVED******REMOVED***
	f.Float64SliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// Float64Slice defines a []float64 flag with specified name, default value, and usage string.
// The return value is the address of a []float64 variable that stores the value of the flag.
func Float64Slice(name string, value []float64, usage string) *[]float64 ***REMOVED***
	return CommandLine.Float64SliceP(name, "", value, usage)
***REMOVED***

// Float64SliceP is like Float64Slice, but accepts a shorthand letter that can be used after a single dash.
func Float64SliceP(name, shorthand string, value []float64, usage string) *[]float64 ***REMOVED***
	return CommandLine.Float64SliceP(name, shorthand, value, usage)
***REMOVED***
