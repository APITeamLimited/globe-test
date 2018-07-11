package pflag

import (
	"bytes"
	"encoding/csv"
	"strings"
)

// -- stringSlice Value
type stringSliceValue struct ***REMOVED***
	value   *[]string
	changed bool
***REMOVED***

func newStringSliceValue(val []string, p *[]string) *stringSliceValue ***REMOVED***
	ssv := new(stringSliceValue)
	ssv.value = p
	*ssv.value = val
	return ssv
***REMOVED***

func readAsCSV(val string) ([]string, error) ***REMOVED***
	if val == "" ***REMOVED***
		return []string***REMOVED******REMOVED***, nil
	***REMOVED***
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
***REMOVED***

func writeAsCSV(vals []string) (string, error) ***REMOVED***
	b := &bytes.Buffer***REMOVED******REMOVED***
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
***REMOVED***

func (s *stringSliceValue) Set(val string) error ***REMOVED***
	v, err := readAsCSV(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !s.changed ***REMOVED***
		*s.value = v
	***REMOVED*** else ***REMOVED***
		*s.value = append(*s.value, v...)
	***REMOVED***
	s.changed = true
	return nil
***REMOVED***

func (s *stringSliceValue) Type() string ***REMOVED***
	return "stringSlice"
***REMOVED***

func (s *stringSliceValue) String() string ***REMOVED***
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
***REMOVED***

func stringSliceConv(sval string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	sval = sval[1 : len(sval)-1]
	// An empty string would cause a slice with one (empty) string
	if len(sval) == 0 ***REMOVED***
		return []string***REMOVED******REMOVED***, nil
	***REMOVED***
	return readAsCSV(sval)
***REMOVED***

// GetStringSlice return the []string value of a flag with the given name
func (f *FlagSet) GetStringSlice(name string) ([]string, error) ***REMOVED***
	val, err := f.getFlagType(name, "stringSlice", stringSliceConv)
	if err != nil ***REMOVED***
		return []string***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]string), nil
***REMOVED***

// StringSliceVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" -ss="v3"
// will result in
//   []string***REMOVED***"v1", "v2", "v3"***REMOVED***
func (f *FlagSet) StringSliceVar(p *[]string, name string, value []string, usage string) ***REMOVED***
	f.VarP(newStringSliceValue(value, p), name, "", usage)
***REMOVED***

// StringSliceVarP is like StringSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringSliceVarP(p *[]string, name, shorthand string, value []string, usage string) ***REMOVED***
	f.VarP(newStringSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// StringSliceVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" -ss="v3"
// will result in
//   []string***REMOVED***"v1", "v2", "v3"***REMOVED***
func StringSliceVar(p *[]string, name string, value []string, usage string) ***REMOVED***
	CommandLine.VarP(newStringSliceValue(value, p), name, "", usage)
***REMOVED***

// StringSliceVarP is like StringSliceVar, but accepts a shorthand letter that can be used after a single dash.
func StringSliceVarP(p *[]string, name, shorthand string, value []string, usage string) ***REMOVED***
	CommandLine.VarP(newStringSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// StringSlice defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" -ss="v3"
// will result in
//   []string***REMOVED***"v1", "v2", "v3"***REMOVED***
func (f *FlagSet) StringSlice(name string, value []string, usage string) *[]string ***REMOVED***
	p := []string***REMOVED******REMOVED***
	f.StringSliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// StringSliceP is like StringSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringSliceP(name, shorthand string, value []string, usage string) *[]string ***REMOVED***
	p := []string***REMOVED******REMOVED***
	f.StringSliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// StringSlice defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" -ss="v3"
// will result in
//   []string***REMOVED***"v1", "v2", "v3"***REMOVED***
func StringSlice(name string, value []string, usage string) *[]string ***REMOVED***
	return CommandLine.StringSliceP(name, "", value, usage)
***REMOVED***

// StringSliceP is like StringSlice, but accepts a shorthand letter that can be used after a single dash.
func StringSliceP(name, shorthand string, value []string, usage string) *[]string ***REMOVED***
	return CommandLine.StringSliceP(name, shorthand, value, usage)
***REMOVED***
