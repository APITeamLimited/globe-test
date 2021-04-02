package pflag

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// -- ipSlice Value
type ipSliceValue struct ***REMOVED***
	value   *[]net.IP
	changed bool
***REMOVED***

func newIPSliceValue(val []net.IP, p *[]net.IP) *ipSliceValue ***REMOVED***
	ipsv := new(ipSliceValue)
	ipsv.value = p
	*ipsv.value = val
	return ipsv
***REMOVED***

// Set converts, and assigns, the comma-separated IP argument string representation as the []net.IP value of this flag.
// If Set is called on a flag that already has a []net.IP assigned, the newly converted values will be appended.
func (s *ipSliceValue) Set(val string) error ***REMOVED***

	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	ipStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF ***REMOVED***
		return err
	***REMOVED***

	// parse ip values into slice
	out := make([]net.IP, 0, len(ipStrSlice))
	for _, ipStr := range ipStrSlice ***REMOVED***
		ip := net.ParseIP(strings.TrimSpace(ipStr))
		if ip == nil ***REMOVED***
			return fmt.Errorf("invalid string being converted to IP address: %s", ipStr)
		***REMOVED***
		out = append(out, ip)
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
func (s *ipSliceValue) Type() string ***REMOVED***
	return "ipSlice"
***REMOVED***

// String defines a "native" format for this net.IP slice flag value.
func (s *ipSliceValue) String() string ***REMOVED***

	ipStrSlice := make([]string, len(*s.value))
	for i, ip := range *s.value ***REMOVED***
		ipStrSlice[i] = ip.String()
	***REMOVED***

	out, _ := writeAsCSV(ipStrSlice)

	return "[" + out + "]"
***REMOVED***

func (s *ipSliceValue) fromString(val string) (net.IP, error) ***REMOVED***
	return net.ParseIP(strings.TrimSpace(val)), nil
***REMOVED***

func (s *ipSliceValue) toString(val net.IP) string ***REMOVED***
	return val.String()
***REMOVED***

func (s *ipSliceValue) Append(val string) error ***REMOVED***
	i, err := s.fromString(val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*s.value = append(*s.value, i)
	return nil
***REMOVED***

func (s *ipSliceValue) Replace(val []string) error ***REMOVED***
	out := make([]net.IP, len(val))
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

func (s *ipSliceValue) GetSlice() []string ***REMOVED***
	out := make([]string, len(*s.value))
	for i, d := range *s.value ***REMOVED***
		out[i] = s.toString(d)
	***REMOVED***
	return out
***REMOVED***

func ipSliceConv(val string) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 ***REMOVED***
		return []net.IP***REMOVED******REMOVED***, nil
	***REMOVED***
	ss := strings.Split(val, ",")
	out := make([]net.IP, len(ss))
	for i, sval := range ss ***REMOVED***
		ip := net.ParseIP(strings.TrimSpace(sval))
		if ip == nil ***REMOVED***
			return nil, fmt.Errorf("invalid string being converted to IP address: %s", sval)
		***REMOVED***
		out[i] = ip
	***REMOVED***
	return out, nil
***REMOVED***

// GetIPSlice returns the []net.IP value of a flag with the given name
func (f *FlagSet) GetIPSlice(name string) ([]net.IP, error) ***REMOVED***
	val, err := f.getFlagType(name, "ipSlice", ipSliceConv)
	if err != nil ***REMOVED***
		return []net.IP***REMOVED******REMOVED***, err
	***REMOVED***
	return val.([]net.IP), nil
***REMOVED***

// IPSliceVar defines a ipSlice flag with specified name, default value, and usage string.
// The argument p points to a []net.IP variable in which to store the value of the flag.
func (f *FlagSet) IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string) ***REMOVED***
	f.VarP(newIPSliceValue(value, p), name, "", usage)
***REMOVED***

// IPSliceVarP is like IPSliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPSliceVarP(p *[]net.IP, name, shorthand string, value []net.IP, usage string) ***REMOVED***
	f.VarP(newIPSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// IPSliceVar defines a []net.IP flag with specified name, default value, and usage string.
// The argument p points to a []net.IP variable in which to store the value of the flag.
func IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string) ***REMOVED***
	CommandLine.VarP(newIPSliceValue(value, p), name, "", usage)
***REMOVED***

// IPSliceVarP is like IPSliceVar, but accepts a shorthand letter that can be used after a single dash.
func IPSliceVarP(p *[]net.IP, name, shorthand string, value []net.IP, usage string) ***REMOVED***
	CommandLine.VarP(newIPSliceValue(value, p), name, shorthand, usage)
***REMOVED***

// IPSlice defines a []net.IP flag with specified name, default value, and usage string.
// The return value is the address of a []net.IP variable that stores the value of that flag.
func (f *FlagSet) IPSlice(name string, value []net.IP, usage string) *[]net.IP ***REMOVED***
	p := []net.IP***REMOVED******REMOVED***
	f.IPSliceVarP(&p, name, "", value, usage)
	return &p
***REMOVED***

// IPSliceP is like IPSlice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPSliceP(name, shorthand string, value []net.IP, usage string) *[]net.IP ***REMOVED***
	p := []net.IP***REMOVED******REMOVED***
	f.IPSliceVarP(&p, name, shorthand, value, usage)
	return &p
***REMOVED***

// IPSlice defines a []net.IP flag with specified name, default value, and usage string.
// The return value is the address of a []net.IP variable that stores the value of the flag.
func IPSlice(name string, value []net.IP, usage string) *[]net.IP ***REMOVED***
	return CommandLine.IPSliceP(name, "", value, usage)
***REMOVED***

// IPSliceP is like IPSlice, but accepts a shorthand letter that can be used after a single dash.
func IPSliceP(name, shorthand string, value []net.IP, usage string) *[]net.IP ***REMOVED***
	return CommandLine.IPSliceP(name, shorthand, value, usage)
***REMOVED***
