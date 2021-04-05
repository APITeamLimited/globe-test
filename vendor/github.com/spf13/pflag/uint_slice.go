package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- uintSlice Value
type uintSliceValue struct ***REMOVED***
	value   *[]uint
	changed bool
***REMOVED***

func newUintSliceValue(val []uint, p *[]uint) *uintSliceValue ***REMOVED***
	uisv := new(uintSliceValue)
	uisv.value = p
	*uisv.value = val
	return uisv
***REMOVED***

func (s *uintSliceValue) Set(val string) error ***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]uint, len(ss))
	for i, d := range ss ***REMOVED***
		u, err := strconv.ParseUint(d, 10, 0)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		out[i] = uint(u)
	***REMOVED***
	if !s.changed ***REMOVED***
		*s.value = out
	***REMOVED*** else ***REMOVED***
		*s.value = append(*s.value, out...)
	***REMOVED***
	s.changed = true
	return nil
***REMOVED***

func (s *uintSliceValue) Type() string ***REMOVED***
	return "uintSlice"
***REMOVED***

func (s *uintSliceValue) String() string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = fmt.Sprintf("%d", d)
	***REMOVED***
	return "[" + strings.Join(out, ",") + "]"
***REMOVED***

func (s *uintSliceValue) fromString(val string) (uint, error) ***REMOVED***
	t, err := strconv.ParseUint(val, 10, 0)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return uint(t), nil
***REMOVED***

func (s *uintSliceValue) toString(val uint) string ***REMOVED***
	return fmt.Sprintf("%d", val)
***REMOVED***

func (s *uintSliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *uintSliceValue) Replace(val []string) error ***REMOVED***
	out := make([]uint, len(val))
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

func (s *uintSliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func uintSliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []uint***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]uint, len(ss))
	for i, d := range ss ***REMOVED***
		u, err := strconv.ParseUint(d, 10, 0)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out[i] = uint(u)
	***REMOVED***
	return out, nil
***REMOVED***

// GetUintSlice returns the []uint value of a flag with the given name.
func (f *FlagSet) GetUintSlice(name string) ([]uint, error) ***REMOVED***
	val, err := f.getFlagType(name, "uintSlice", uintSliceConv)
	if err != nil ***REMOVED***
		return []uint***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]uint), nil
***REMOVED***

// UintSliceVar defines a uintSlice flag with specified name, default value, and usage string.
// The argument p points to a []uint variable in which to store the value of the flag.
func (f *FlagSet) UintSliceVar(p *[]uint, name string, value []uint, usage string) ***REMOVED***
	f.VarP(newUintSliceValue(value, p), name, "", usage)
***REMOVED***

// UintSliceVarP is like UintSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UintSliceVarP(p *[]uint, name, shorthand string, value []uint, usage string) ***REMOVED***
	f.VarP(newUintSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// UintSliceVar defines a uint[] flag with specified name, default value, and usage string.
// The argument p points to a uint[] variable in which to store the value of the flag.
func UintSliceVar(p *[]uint, name string, value []uint, usage string) ***REMOVED***
	CommandLine.VarP(newUintSliceValue(value, p), name, "", usage)
***REMOVED***

// UintSliceVarP is like the UintSliceVar, but accepts a shorthand letter that can be used after a single dash.
func UintSliceVarP(p *[]uint, name, shorthand string, value []uint, usage string) ***REMOVED***
	CommandLine.VarP(newUintSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// UintSlice defines a []uint flag with specified name, default value, and usage string.
// The return value is the address of a []uint variable that stores the value of the flag.
func (f *FlagSet) UintSlice(name string, value []uint, usage string) *[]uint ***REMOVED***
	p := []uint***REMOVED******REMOVED***
	f.UintSliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// UintSliceP is like UintSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UintSliceP(name, shorthand string, value []uint, usage string) *[]uint ***REMOVED***
	p := []uint***REMOVED******REMOVED***
	f.UintSliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// UintSlice defines a []uint flag with specified name, default value, and usage string.
// The return value is the address of a []uint variable that stores the value of the flag.
func UintSlice(name string, value []uint, usage string) *[]uint ***REMOVED***
	return CommandLine.UintSliceP(name, "", value, usage)
***REMOVED***

// UintSliceP is like UintSlice, but accepts a shorthand letter that can be used after a single dash.
func UintSliceP(name, shorthand string, value []uint, usage string) *[]uint ***REMOVED***
	return CommandLine.UintSliceP(name, shorthand, value, usage)
***REMOVED***
