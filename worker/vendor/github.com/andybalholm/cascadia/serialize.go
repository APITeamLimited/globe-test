package cascadia

import (
	"fmt"
	"strconv"
	"strings"
)

// implements the reverse operation Sel -> string

var specialCharReplacer *strings.Replacer

func init() ***REMOVED***
	var pairs []string
	for _, s := range ",!\"#$%&'()*+ -./:;<=>?@[\\]^`***REMOVED***|***REMOVED***~" ***REMOVED***
		pairs = append(pairs, string(s), "\\"+string(s))
	***REMOVED***
	specialCharReplacer = strings.NewReplacer(pairs...)
***REMOVED***

// espace special CSS char
func escape(s string) string ***REMOVED*** return specialCharReplacer.Replace(s) ***REMOVED***

func (c tagSelector) String() string ***REMOVED***
	return c.tag
***REMOVED***

func (c idSelector) String() string ***REMOVED***
	return "#" + escape(c.id)
***REMOVED***

func (c classSelector) String() string ***REMOVED***
	return "." + escape(c.class)
***REMOVED***

func (c attrSelector) String() string ***REMOVED***
	val := c.val
	if c.operation == "#=" ***REMOVED***
		val = c.regexp.String()
	***REMOVED*** else if c.operation != "" ***REMOVED***
		val = fmt.Sprintf(`"%s"`, val)
	***REMOVED***

	ignoreCase := ""

	if c.insensitive ***REMOVED***
		ignoreCase = " i"
	***REMOVED***

	return fmt.Sprintf(`[%s%s%s%s]`, c.key, c.operation, val, ignoreCase)
***REMOVED***

func (c relativePseudoClassSelector) String() string ***REMOVED***
	return fmt.Sprintf(":%s(%s)", c.name, c.match.String())
***REMOVED***

func (c containsPseudoClassSelector) String() string ***REMOVED***
	s := "contains"
	if c.own ***REMOVED***
		s += "Own"
	***REMOVED***
	return fmt.Sprintf(`:%s("%s")`, s, c.value)
***REMOVED***

func (c regexpPseudoClassSelector) String() string ***REMOVED***
	s := "matches"
	if c.own ***REMOVED***
		s += "Own"
	***REMOVED***
	return fmt.Sprintf(":%s(%s)", s, c.regexp.String())
***REMOVED***

func (c nthPseudoClassSelector) String() string ***REMOVED***
	if c.a == 0 && c.b == 1 ***REMOVED*** // special cases
		s := ":first-"
		if c.last ***REMOVED***
			s = ":last-"
		***REMOVED***
		if c.ofType ***REMOVED***
			s += "of-type"
		***REMOVED*** else ***REMOVED***
			s += "child"
		***REMOVED***
		return s
	***REMOVED***
	var name string
	switch [2]bool***REMOVED***c.last, c.ofType***REMOVED*** ***REMOVED***
	case [2]bool***REMOVED***true, true***REMOVED***:
		name = "nth-last-of-type"
	case [2]bool***REMOVED***true, false***REMOVED***:
		name = "nth-last-child"
	case [2]bool***REMOVED***false, true***REMOVED***:
		name = "nth-of-type"
	case [2]bool***REMOVED***false, false***REMOVED***:
		name = "nth-child"
	***REMOVED***
	s := fmt.Sprintf("+%d", c.b)
	if c.b < 0 ***REMOVED*** // avoid +-8 invalid syntax
		s = strconv.Itoa(c.b)
	***REMOVED***
	return fmt.Sprintf(":%s(%dn%s)", name, c.a, s)
***REMOVED***

func (c onlyChildPseudoClassSelector) String() string ***REMOVED***
	if c.ofType ***REMOVED***
		return ":only-of-type"
	***REMOVED***
	return ":only-child"
***REMOVED***

func (c inputPseudoClassSelector) String() string ***REMOVED***
	return ":input"
***REMOVED***

func (c emptyElementPseudoClassSelector) String() string ***REMOVED***
	return ":empty"
***REMOVED***

func (c rootPseudoClassSelector) String() string ***REMOVED***
	return ":root"
***REMOVED***

func (c linkPseudoClassSelector) String() string ***REMOVED***
	return ":link"
***REMOVED***

func (c langPseudoClassSelector) String() string ***REMOVED***
	return fmt.Sprintf(":lang(%s)", c.lang)
***REMOVED***

func (c neverMatchSelector) String() string ***REMOVED***
	return c.value
***REMOVED***

func (c enabledPseudoClassSelector) String() string ***REMOVED***
	return ":enabled"
***REMOVED***

func (c disabledPseudoClassSelector) String() string ***REMOVED***
	return ":disabled"
***REMOVED***

func (c checkedPseudoClassSelector) String() string ***REMOVED***
	return ":checked"
***REMOVED***

func (c compoundSelector) String() string ***REMOVED***
	if len(c.selectors) == 0 && c.pseudoElement == "" ***REMOVED***
		return "*"
	***REMOVED***
	chunks := make([]string, len(c.selectors))
	for i, sel := range c.selectors ***REMOVED***
		chunks[i] = sel.String()
	***REMOVED***
	s := strings.Join(chunks, "")
	if c.pseudoElement != "" ***REMOVED***
		s += "::" + c.pseudoElement
	***REMOVED***
	return s
***REMOVED***

func (c combinedSelector) String() string ***REMOVED***
	start := c.first.String()
	if c.second != nil ***REMOVED***
		start += fmt.Sprintf(" %s %s", string(c.combinator), c.second.String())
	***REMOVED***
	return start
***REMOVED***

func (c SelectorGroup) String() string ***REMOVED***
	ck := make([]string, len(c))
	for i, s := range c ***REMOVED***
		ck[i] = s.String()
	***REMOVED***
	return strings.Join(ck, ", ")
***REMOVED***
