package regexp2

import (
	"bytes"
	"errors"

	"github.com/dlclark/regexp2/syntax"
)

const (
	replaceSpecials     = 4
	replaceLeftPortion  = -1
	replaceRightPortion = -2
	replaceLastGroup    = -3
	replaceWholeString  = -4
)

// MatchEvaluator is a function that takes a match and returns a replacement string to be used
type MatchEvaluator func(Match) string

// Three very similar algorithms appear below: replace (pattern),
// replace (evaluator), and split.

// Replace Replaces all occurrences of the regex in the string with the
// replacement pattern.
//
// Note that the special case of no matches is handled on its own:
// with no matches, the input string is returned unchanged.
// The right-to-left case is split out because StringBuilder
// doesn't handle right-to-left string building directly very well.
func replace(regex *Regexp, data *syntax.ReplacerData, evaluator MatchEvaluator, input string, startAt, count int) (string, error) ***REMOVED***
	if count < -1 ***REMOVED***
		return "", errors.New("Count too small")
	***REMOVED***
	if count == 0 ***REMOVED***
		return "", nil
	***REMOVED***

	m, err := regex.FindStringMatchStartingAt(input, startAt)

	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if m == nil ***REMOVED***
		return input, nil
	***REMOVED***

	buf := &bytes.Buffer***REMOVED******REMOVED***
	text := m.text

	if !regex.RightToLeft() ***REMOVED***
		prevat := 0
		for m != nil ***REMOVED***
			if m.Index != prevat ***REMOVED***
				buf.WriteString(string(text[prevat:m.Index]))
			***REMOVED***
			prevat = m.Index + m.Length
			if evaluator == nil ***REMOVED***
				replacementImpl(data, buf, m)
			***REMOVED*** else ***REMOVED***
				buf.WriteString(evaluator(*m))
			***REMOVED***

			count--
			if count == 0 ***REMOVED***
				break
			***REMOVED***
			m, err = regex.FindNextMatch(m)
			if err != nil ***REMOVED***
				return "", nil
			***REMOVED***
		***REMOVED***

		if prevat < len(text) ***REMOVED***
			buf.WriteString(string(text[prevat:]))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		prevat := len(text)
		var al []string

		for m != nil ***REMOVED***
			if m.Index+m.Length != prevat ***REMOVED***
				al = append(al, string(text[m.Index+m.Length:prevat]))
			***REMOVED***
			prevat = m.Index
			if evaluator == nil ***REMOVED***
				replacementImplRTL(data, &al, m)
			***REMOVED*** else ***REMOVED***
				al = append(al, evaluator(*m))
			***REMOVED***

			count--
			if count == 0 ***REMOVED***
				break
			***REMOVED***
			m, err = regex.FindNextMatch(m)
			if err != nil ***REMOVED***
				return "", nil
			***REMOVED***
		***REMOVED***

		if prevat > 0 ***REMOVED***
			buf.WriteString(string(text[:prevat]))
		***REMOVED***

		for i := len(al) - 1; i >= 0; i-- ***REMOVED***
			buf.WriteString(al[i])
		***REMOVED***
	***REMOVED***

	return buf.String(), nil
***REMOVED***

// Given a Match, emits into the StringBuilder the evaluated
// substitution pattern.
func replacementImpl(data *syntax.ReplacerData, buf *bytes.Buffer, m *Match) ***REMOVED***
	for _, r := range data.Rules ***REMOVED***

		if r >= 0 ***REMOVED*** // string lookup
			buf.WriteString(data.Strings[r])
		***REMOVED*** else if r < -replaceSpecials ***REMOVED*** // group lookup
			m.groupValueAppendToBuf(-replaceSpecials-1-r, buf)
		***REMOVED*** else ***REMOVED***
			switch -replaceSpecials - 1 - r ***REMOVED*** // special insertion patterns
			case replaceLeftPortion:
				for i := 0; i < m.Index; i++ ***REMOVED***
					buf.WriteRune(m.text[i])
				***REMOVED***
			case replaceRightPortion:
				for i := m.Index + m.Length; i < len(m.text); i++ ***REMOVED***
					buf.WriteRune(m.text[i])
				***REMOVED***
			case replaceLastGroup:
				m.groupValueAppendToBuf(m.GroupCount()-1, buf)
			case replaceWholeString:
				for i := 0; i < len(m.text); i++ ***REMOVED***
					buf.WriteRune(m.text[i])
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func replacementImplRTL(data *syntax.ReplacerData, al *[]string, m *Match) ***REMOVED***
	l := *al
	buf := &bytes.Buffer***REMOVED******REMOVED***

	for _, r := range data.Rules ***REMOVED***
		buf.Reset()
		if r >= 0 ***REMOVED*** // string lookup
			l = append(l, data.Strings[r])
		***REMOVED*** else if r < -replaceSpecials ***REMOVED*** // group lookup
			m.groupValueAppendToBuf(-replaceSpecials-1-r, buf)
			l = append(l, buf.String())
		***REMOVED*** else ***REMOVED***
			switch -replaceSpecials - 1 - r ***REMOVED*** // special insertion patterns
			case replaceLeftPortion:
				for i := 0; i < m.Index; i++ ***REMOVED***
					buf.WriteRune(m.text[i])
				***REMOVED***
			case replaceRightPortion:
				for i := m.Index + m.Length; i < len(m.text); i++ ***REMOVED***
					buf.WriteRune(m.text[i])
				***REMOVED***
			case replaceLastGroup:
				m.groupValueAppendToBuf(m.GroupCount()-1, buf)
			case replaceWholeString:
				for i := 0; i < len(m.text); i++ ***REMOVED***
					buf.WriteRune(m.text[i])
				***REMOVED***
			***REMOVED***
			l = append(l, buf.String())
		***REMOVED***
	***REMOVED***

	*al = l
***REMOVED***
