package pflag

import (
	"fmt"
	"strings"
	"time"
)

// -- durationSlice Value
type durationSliceValue struct ***REMOVED***
	value   *[]time.Duration
	changed bool
***REMOVED***

func newDurationSliceValue(val []time.Duration, p *[]time.Duration) *durationSliceValue ***REMOVED***
	dsv := new(durationSliceValue)
	dsv.value = p
	*dsv.value = val
	return dsv
***REMOVED***

func (s *durationSliceValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]time.Duration, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = time.ParseDuration(d)
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

func (s *durationSliceValue) Type() string ***REMOVED***
	return "durationSlice"
***REMOVED***

func (s *durationSliceValue) String() string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = fmt.Sprintf("%s", d)
	***REMOVED***
	return "[" + strings.Join(out, ",") + "]"
***REMOVED***

func (s *durationSliceValue) fromString(val string) (time.Duration, error) ***REMOVED***
	return time.ParseDuration(val)
***REMOVED***

func (s *durationSliceValue) toString(val time.Duration) string ***REMOVED***
	return fmt.Sprintf("%s", val)
***REMOVED***

func (s *durationSliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *durationSliceValue) Replace(val []string) error ***REMOVED***
	out := make([]time.Duration, len(val))
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

func (s *durationSliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func durationSliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []time.Duration***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]time.Duration, len(ss))
	for i, d := range ss ***REMOVED***
		var err error
		out[i], err = time.ParseDuration(d)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

	***REMOVED***
	return out, nil
***REMOVED***

// GetDurationSlice returns the []time.Duration value of a flag with the given name
func (f *FlagSet) GetDurationSlice(name string) ([]time.Duration, error) ***REMOVED***
	val, err := f.getFlagType(name, "durationSlice", durationSliceConv)
	if err != nil ***REMOVED***
		return []time.Duration***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]time.Duration), nil
***REMOVED***

// DurationSliceVar defines a durationSlice flag with specified name, default value, and usage string.
// The argument p points to a []time.Duration variable in which to store the value of the flag.
func (f *FlagSet) DurationSliceVar(p *[]time.Duration, name string, value []time.Duration, usage string) ***REMOVED***
	f.VarP(newDurationSliceValue(value, p), name, "", usage)
***REMOVED***

// DurationSliceVarP is like DurationSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) DurationSliceVarP(p *[]time.Duration, name, shorthand string, value []time.Duration, usage string) ***REMOVED***
	f.VarP(newDurationSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// DurationSliceVar defines a duration[] flag with specified name, default value, and usage string.
// The argument p points to a duration[] variable in which to store the value of the flag.
func DurationSliceVar(p *[]time.Duration, name string, value []time.Duration, usage string) ***REMOVED***
	CommandLine.VarP(newDurationSliceValue(value, p), name, "", usage)
***REMOVED***

// DurationSliceVarP is like DurationSliceVar, but accepts a shorthand letter that can be used after a single dash.
func DurationSliceVarP(p *[]time.Duration, name, shorthand string, value []time.Duration, usage string) ***REMOVED***
	CommandLine.VarP(newDurationSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// DurationSlice defines a []time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a []time.Duration variable that stores the value of the flag.
func (f *FlagSet) DurationSlice(name string, value []time.Duration, usage string) *[]time.Duration ***REMOVED***
	p := []time.Duration***REMOVED******REMOVED***
	f.DurationSliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// DurationSliceP is like DurationSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) DurationSliceP(name, shorthand string, value []time.Duration, usage string) *[]time.Duration ***REMOVED***
	p := []time.Duration***REMOVED******REMOVED***
	f.DurationSliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// DurationSlice defines a []time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a []time.Duration variable that stores the value of the flag.
func DurationSlice(name string, value []time.Duration, usage string) *[]time.Duration ***REMOVED***
	return CommandLine.DurationSliceP(name, "", value, usage)
***REMOVED***

// DurationSliceP is like DurationSlice, but accepts a shorthand letter that can be used after a single dash.
func DurationSliceP(name, shorthand string, value []time.Duration, usage string) *[]time.Duration ***REMOVED***
	return CommandLine.DurationSliceP(name, shorthand, value, usage)
***REMOVED***
