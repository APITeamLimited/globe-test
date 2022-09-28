package stringprep

import (
	"golang.org/x/text/unicode/norm"
)

// Profile represents a stringprep profile.
type Profile struct ***REMOVED***
	Mappings  []Mapping
	Normalize bool
	Prohibits []Set
	CheckBiDi bool
***REMOVED***

var errProhibited = "prohibited character"

// Prepare transforms an input string to an output string following
// the rules defined in the profile as defined by RFC-3454.
func (p Profile) Prepare(s string) (string, error) ***REMOVED***
	// Optimistically, assume output will be same length as input
	temp := make([]rune, 0, len(s))

	// Apply maps
	for _, r := range s ***REMOVED***
		rs, ok := p.applyMaps(r)
		if ok ***REMOVED***
			temp = append(temp, rs...)
		***REMOVED*** else ***REMOVED***
			temp = append(temp, r)
		***REMOVED***
	***REMOVED***

	// Normalize
	var out string
	if p.Normalize ***REMOVED***
		out = norm.NFKC.String(string(temp))
	***REMOVED*** else ***REMOVED***
		out = string(temp)
	***REMOVED***

	// Check prohibited
	for _, r := range out ***REMOVED***
		if p.runeIsProhibited(r) ***REMOVED***
			return "", Error***REMOVED***Msg: errProhibited, Rune: r***REMOVED***
		***REMOVED***
	***REMOVED***

	// Check BiDi allowed
	if p.CheckBiDi ***REMOVED***
		if err := passesBiDiRules(out); err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***

	return out, nil
***REMOVED***

func (p Profile) applyMaps(r rune) ([]rune, bool) ***REMOVED***
	for _, m := range p.Mappings ***REMOVED***
		rs, ok := m.Map(r)
		if ok ***REMOVED***
			return rs, true
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func (p Profile) runeIsProhibited(r rune) bool ***REMOVED***
	for _, s := range p.Prohibits ***REMOVED***
		if s.Contains(r) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
