package syntax

import (
	"bytes"
	"errors"
)

type ReplacerData struct ***REMOVED***
	Rep     string
	Strings []string
	Rules   []int
***REMOVED***

const (
	replaceSpecials     = 4
	replaceLeftPortion  = -1
	replaceRightPortion = -2
	replaceLastGroup    = -3
	replaceWholeString  = -4
)

//ErrReplacementError is a general error during parsing the replacement text
var ErrReplacementError = errors.New("Replacement pattern error.")

// NewReplacerData will populate a reusable replacer data struct based on the given replacement string
// and the capture group data from a regexp
func NewReplacerData(rep string, caps map[int]int, capsize int, capnames map[string]int, op RegexOptions) (*ReplacerData, error) ***REMOVED***
	p := parser***REMOVED***
		options:  op,
		caps:     caps,
		capsize:  capsize,
		capnames: capnames,
	***REMOVED***
	p.setPattern(rep)
	concat, err := p.scanReplacement()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if concat.t != ntConcatenate ***REMOVED***
		panic(ErrReplacementError)
	***REMOVED***

	sb := &bytes.Buffer***REMOVED******REMOVED***
	var (
		strings []string
		rules   []int
	)

	for _, child := range concat.children ***REMOVED***
		switch child.t ***REMOVED***
		case ntMulti:
			child.writeStrToBuf(sb)

		case ntOne:
			sb.WriteRune(child.ch)

		case ntRef:
			if sb.Len() > 0 ***REMOVED***
				rules = append(rules, len(strings))
				strings = append(strings, sb.String())
				sb.Reset()
			***REMOVED***
			slot := child.m

			if len(caps) > 0 && slot >= 0 ***REMOVED***
				slot = caps[slot]
			***REMOVED***

			rules = append(rules, -replaceSpecials-1-slot)

		default:
			panic(ErrReplacementError)
		***REMOVED***
	***REMOVED***

	if sb.Len() > 0 ***REMOVED***
		rules = append(rules, len(strings))
		strings = append(strings, sb.String())
	***REMOVED***

	return &ReplacerData***REMOVED***
		Rep:     rep,
		Strings: strings,
		Rules:   rules,
	***REMOVED***, nil
***REMOVED***
