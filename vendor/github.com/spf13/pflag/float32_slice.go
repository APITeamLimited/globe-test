package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- float32Slice Value
type float32SliceValue struct ***REMOVED***
	value   *[]float32
	changed bool
***REMOVED***

func newFloat32SliceValue(val []float32, p *[]float32) *float32SliceValue ***REMOVED***
	isv := new(float32SliceValue)
	isv.value = p
	*isv.value = val
	return isv
***REMOVED***

func (s *float32SliceValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]float32, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		var temp64 float64
		temp64, err = strconv.ParseFloat(d, 32)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		out[i] = float32(temp64)

	***REMOVED***
	if !s.changed ***REMOVED***
		*s.value = out
	***REMOVED*** else ***REMOVED***
		*s.value = append(*s.value, out...)
	***REMOVED***
	s.changed = true
	return nil
***REMOVED***

func (s *float32SliceValue) Type() string ***REMOVED***
	return "float32Slice"
***REMOVED***

func (s *float32SliceValue) String() string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = fmt.Sprintf("%f", d)
	***REMOVED***
	return "[" + strings.Join(out, ",") + "]"
***REMOVED***

func (s *float32SliceValue) fromString(val string) (float32, error) ***REMOVED***
	t64, err := strconv.ParseFloat(val, 32)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return float32(t64), nil
***REMOVED***

func (s *float32SliceValue) toString(val float32) string ***REMOVED***
	return fmt.Sprintf("%f", val)
***REMOVED***

func (s *float32SliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *float32SliceValue) Replace(val []string) error ***REMOVED***
	out := make([]float32, len(val))
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

func (s *float32SliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func float32SliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []float32***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]float32, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		var temp64 float64
		temp64, err = strconv.ParseFloat(d, 32)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out[i] = float32(temp64)

	***REMOVED***
	return out, nil
***REMOVED***

// GetFloat32Slice return the []float32 value of a flag with the given name
func (f *FlagSet) GetFloat32Slice(name string) ([]float32, error) ***REMOVED***
	val, err := f.getFlagType(name, "float32Slice", float32SliceConv)
	if err != nil ***REMOVED***
		return []float32***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]float32), nil
***REMOVED***

// Float32SliceVar defines a float32Slice flag with specified name, default value, and usage string.
// The argument p points to a []float32 variable in which to store the value of the flag.
func (f *FlagSet) Float32SliceVar(p *[]float32, name string, value []float32, usage string) ***REMOVED***
	f.VarP(newFloat32SliceValue(value, p), name, "", usage)
***REMOVED***

// Float32SliceVarP is like Float32SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32SliceVarP(p *[]float32, name, shorthand string, value []float32, usage string) ***REMOVED***
	f.VarP(newFloat32SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Float32SliceVar defines a float32[] flag with specified name, default value, and usage string.
// The argument p points to a float32[] variable in which to store the value of the flag.
func Float32SliceVar(p *[]float32, name string, value []float32, usage string) ***REMOVED***
	CommandLine.VarP(newFloat32SliceValue(value, p), name, "", usage)
***REMOVED***

// Float32SliceVarP is like Float32SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Float32SliceVarP(p *[]float32, name, shorthand string, value []float32, usage string) ***REMOVED***
	CommandLine.VarP(newFloat32SliceValue(value, p), name, shorthand, usage)
***REMOVED***

// Float32Slice defines a []float32 flag with specified name, default value, and usage string.
// The return value is the address of a []float32 variable that stores the value of the flag.
func (f *FlagSet) Float32Slice(name string, value []float32, usage string) *[]float32 ***REMOVED***
	p := []float32***REMOVED******REMOVED***
	f.Float32SliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// Float32SliceP is like Float32Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32SliceP(name, shorthand string, value []float32, usage string) *[]float32 ***REMOVED***
	p := []float32***REMOVED******REMOVED***
	f.Float32SliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// Float32Slice defines a []float32 flag with specified name, default value, and usage string.
// The return value is the address of a []float32 variable that stores the value of the flag.
func Float32Slice(name string, value []float32, usage string) *[]float32 ***REMOVED***
	return CommandLine.Float32SliceP(name, "", value, usage)
***REMOVED***

// Float32SliceP is like Float32Slice, but accepts a shorthand letter that can be used after a single dash.
func Float32SliceP(name, shorthand string, value []float32, usage string) *[]float32 ***REMOVED***
	return CommandLine.Float32SliceP(name, shorthand, value, usage)
***REMOVED***
