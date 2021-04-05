package pflag

import (
	"io"
	"strconv"
	"strings"
)

// -- boolSlice Value
type boolSliceValue struct ***REMOVED***
	value   *[]bool
	changed bool
***REMOVED***

func newBoolSliceValue(val []bool, p *[]bool) *boolSliceValue ***REMOVED***
	bsv := new(boolSliceValue)
	bsv.value = p
	*bsv.value = val
	return bsv
***REMOVED***

// Set converts, and assigns, the comma-separated boolean argument string representation as the []bool value of this flag.
// If Set is called on a flag that already has a []bool assigned, the newly converted values will be appended.
func (s *boolSliceValue) Set(val string) error ***REMOVED***

	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	boolStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF ***REMOVED***
		return err
	***REMOVED***

	// parse boolean values into slice
	out := make([]bool, 0, len(boolStrSlice))
	for _, boolStr := range boolStrSlice ***REMOVED***
		b, err := strconv.ParseBool(strings.TrimSpace(boolStr))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		out = append(out, b)
	***REMOVED***

	if !s.changed ***REMOVED***
		*s.value = out
	***REMOVED*** else ***REMOVED***
		*s.value = append(*s.value, out...)
	***REMOVED***

	s.changed = true

	return nil
***REMOVED***

// Type returns a string that uniquely represents this flag's type.
func (s *boolSliceValue) Type() string ***REMOVED***
	return "boolSlice"
***REMOVED***

// String defines a "native" format for this boolean slice flag value.
func (s *boolSliceValue) String() string ***REMOVED***

	boolStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value ***REMOVED***
		boolStrSlice[i] = strconv.FormatBool(b)
	***REMOVED***

	out, _ := writeAsCSV(boolStrSlice)

	return "[" + out + "]"
***REMOVED***

func (s *boolSliceValue) fromString(val string) (bool, error) ***REMOVED***
	return strconv.ParseBool(val)
***REMOVED***

func (s *boolSliceValue) toString(val bool) string ***REMOVED***
	return strconv.FormatBool(val)
***REMOVED***

func (s *boolSliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *boolSliceValue) Replace(val []string) error ***REMOVED***
	out := make([]bool, len(val))
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

func (s *boolSliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func boolSliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []bool***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]bool, len(ss))
	for i, t := range ss ***REMOVED***
		var err error
		out[i], err = strconv.ParseBool(t)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***

// GetBoolSlice returns the []bool value of a flag with the given name.
func (f *FlagSet) GetBoolSlice(name string) ([]bool, error) ***REMOVED***
	val, err := f.getFlagType(name, "boolSlice", boolSliceConv)
	if err != nil ***REMOVED***
		return []bool***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]bool), nil
***REMOVED***

// BoolSliceVar defines a boolSlice flag with specified name, default value, and usage string.
// The argument p points to a []bool variable in which to store the value of the flag.
func (f *FlagSet) BoolSliceVar(p *[]bool, name string, value []bool, usage string) ***REMOVED***
	f.VarP(newBoolSliceValue(value, p), name, "", usage)
***REMOVED***

// BoolSliceVarP is like BoolSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolSliceVarP(p *[]bool, name, shorthand string, value []bool, usage string) ***REMOVED***
	f.VarP(newBoolSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// BoolSliceVar defines a []bool flag with specified name, default value, and usage string.
// The argument p points to a []bool variable in which to store the value of the flag.
func BoolSliceVar(p *[]bool, name string, value []bool, usage string) ***REMOVED***
	CommandLine.VarP(newBoolSliceValue(value, p), name, "", usage)
***REMOVED***

// BoolSliceVarP is like BoolSliceVar, but accepts a shorthand letter that can be used after a single dash.
func BoolSliceVarP(p *[]bool, name, shorthand string, value []bool, usage string) ***REMOVED***
	CommandLine.VarP(newBoolSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func (f *FlagSet) BoolSlice(name string, value []bool, usage string) *[]bool ***REMOVED***
	p := []bool***REMOVED******REMOVED***
	f.BoolSliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// BoolSliceP is like BoolSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool ***REMOVED***
	p := []bool***REMOVED******REMOVED***
	f.BoolSliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func BoolSlice(name string, value []bool, usage string) *[]bool ***REMOVED***
	return CommandLine.BoolSliceP(name, "", value, usage)
***REMOVED***

// BoolSliceP is like BoolSlice, but accepts a shorthand letter that can be used after a single dash.
func BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool ***REMOVED***
	return CommandLine.BoolSliceP(name, shorthand, value, usage)
***REMOVED***
