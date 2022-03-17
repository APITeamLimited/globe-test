package cascadia

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Matcher is the interface for basic selector functionality.
// Match returns whether a selector matches n.
type Matcher interface ***REMOVED***
	Match(n *html.Node) bool
***REMOVED***

// Sel is the interface for all the functionality provided by selectors.
type Sel interface ***REMOVED***
	Matcher
	Specificity() Specificity

	// Returns a CSS input compiling to this selector.
	String() string

	// Returns a pseudo-element, or an empty string.
	PseudoElement() string
***REMOVED***

// Parse parses a selector. Use `ParseWithPseudoElement`
// if you need support for pseudo-elements.
func Parse(sel string) (Sel, error) ***REMOVED***
	p := &parser***REMOVED***s: sel***REMOVED***
	compiled, err := p.parseSelector()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if p.i < len(sel) ***REMOVED***
		return nil, fmt.Errorf("parsing %q: %d bytes left over", sel, len(sel)-p.i)
	***REMOVED***

	return compiled, nil
***REMOVED***

// ParseWithPseudoElement parses a single selector,
// with support for pseudo-element.
func ParseWithPseudoElement(sel string) (Sel, error) ***REMOVED***
	p := &parser***REMOVED***s: sel, acceptPseudoElements: true***REMOVED***
	compiled, err := p.parseSelector()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if p.i < len(sel) ***REMOVED***
		return nil, fmt.Errorf("parsing %q: %d bytes left over", sel, len(sel)-p.i)
	***REMOVED***

	return compiled, nil
***REMOVED***

// ParseGroup parses a selector, or a group of selectors separated by commas.
// Use `ParseGroupWithPseudoElements`
// if you need support for pseudo-elements.
func ParseGroup(sel string) (SelectorGroup, error) ***REMOVED***
	p := &parser***REMOVED***s: sel***REMOVED***
	compiled, err := p.parseSelectorGroup()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if p.i < len(sel) ***REMOVED***
		return nil, fmt.Errorf("parsing %q: %d bytes left over", sel, len(sel)-p.i)
	***REMOVED***

	return compiled, nil
***REMOVED***

// ParseGroupWithPseudoElements parses a selector, or a group of selectors separated by commas.
// It supports pseudo-elements.
func ParseGroupWithPseudoElements(sel string) (SelectorGroup, error) ***REMOVED***
	p := &parser***REMOVED***s: sel, acceptPseudoElements: true***REMOVED***
	compiled, err := p.parseSelectorGroup()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if p.i < len(sel) ***REMOVED***
		return nil, fmt.Errorf("parsing %q: %d bytes left over", sel, len(sel)-p.i)
	***REMOVED***

	return compiled, nil
***REMOVED***

// A Selector is a function which tells whether a node matches or not.
//
// This type is maintained for compatibility; I recommend using the newer and
// more idiomatic interfaces Sel and Matcher.
type Selector func(*html.Node) bool

// Compile parses a selector and returns, if successful, a Selector object
// that can be used to match against html.Node objects.
func Compile(sel string) (Selector, error) ***REMOVED***
	compiled, err := ParseGroup(sel)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return Selector(compiled.Match), nil
***REMOVED***

// MustCompile is like Compile, but panics instead of returning an error.
func MustCompile(sel string) Selector ***REMOVED***
	compiled, err := Compile(sel)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return compiled
***REMOVED***

// MatchAll returns a slice of the nodes that match the selector,
// from n and its children.
func (s Selector) MatchAll(n *html.Node) []*html.Node ***REMOVED***
	return s.matchAllInto(n, nil)
***REMOVED***

func (s Selector) matchAllInto(n *html.Node, storage []*html.Node) []*html.Node ***REMOVED***
	if s(n) ***REMOVED***
		storage = append(storage, n)
	***REMOVED***

	for child := n.FirstChild; child != nil; child = child.NextSibling ***REMOVED***
		storage = s.matchAllInto(child, storage)
	***REMOVED***

	return storage
***REMOVED***

func queryInto(n *html.Node, m Matcher, storage []*html.Node) []*html.Node ***REMOVED***
	for child := n.FirstChild; child != nil; child = child.NextSibling ***REMOVED***
		if m.Match(child) ***REMOVED***
			storage = append(storage, child)
		***REMOVED***
		storage = queryInto(child, m, storage)
	***REMOVED***

	return storage
***REMOVED***

// QueryAll returns a slice of all the nodes that match m, from the descendants
// of n.
func QueryAll(n *html.Node, m Matcher) []*html.Node ***REMOVED***
	return queryInto(n, m, nil)
***REMOVED***

// Match returns true if the node matches the selector.
func (s Selector) Match(n *html.Node) bool ***REMOVED***
	return s(n)
***REMOVED***

// MatchFirst returns the first node that matches s, from n and its children.
func (s Selector) MatchFirst(n *html.Node) *html.Node ***REMOVED***
	if s.Match(n) ***REMOVED***
		return n
	***REMOVED***

	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		m := s.MatchFirst(c)
		if m != nil ***REMOVED***
			return m
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Query returns the first node that matches m, from the descendants of n.
// If none matches, it returns nil.
func Query(n *html.Node, m Matcher) *html.Node ***REMOVED***
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if m.Match(c) ***REMOVED***
			return c
		***REMOVED***
		if matched := Query(c, m); matched != nil ***REMOVED***
			return matched
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Filter returns the nodes in nodes that match the selector.
func (s Selector) Filter(nodes []*html.Node) (result []*html.Node) ***REMOVED***
	for _, n := range nodes ***REMOVED***
		if s(n) ***REMOVED***
			result = append(result, n)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// Filter returns the nodes that match m.
func Filter(nodes []*html.Node, m Matcher) (result []*html.Node) ***REMOVED***
	for _, n := range nodes ***REMOVED***
		if m.Match(n) ***REMOVED***
			result = append(result, n)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

type tagSelector struct ***REMOVED***
	tag string
***REMOVED***

// Matches elements with a given tag name.
func (t tagSelector) Match(n *html.Node) bool ***REMOVED***
	return n.Type == html.ElementNode && n.Data == t.tag
***REMOVED***

func (c tagSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 0, 1***REMOVED***
***REMOVED***

func (c tagSelector) PseudoElement() string ***REMOVED***
	return ""
***REMOVED***

type classSelector struct ***REMOVED***
	class string
***REMOVED***

// Matches elements by class attribute.
func (t classSelector) Match(n *html.Node) bool ***REMOVED***
	return matchAttribute(n, "class", func(s string) bool ***REMOVED***
		return matchInclude(t.class, s, false)
	***REMOVED***)
***REMOVED***

func (c classSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

func (c classSelector) PseudoElement() string ***REMOVED***
	return ""
***REMOVED***

type idSelector struct ***REMOVED***
	id string
***REMOVED***

// Matches elements by id attribute.
func (t idSelector) Match(n *html.Node) bool ***REMOVED***
	return matchAttribute(n, "id", func(s string) bool ***REMOVED***
		return s == t.id
	***REMOVED***)
***REMOVED***

func (c idSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***1, 0, 0***REMOVED***
***REMOVED***

func (c idSelector) PseudoElement() string ***REMOVED***
	return ""
***REMOVED***

type attrSelector struct ***REMOVED***
	key, val, operation string
	regexp              *regexp.Regexp
	insensitive         bool
***REMOVED***

// Matches elements by attribute value.
func (t attrSelector) Match(n *html.Node) bool ***REMOVED***
	switch t.operation ***REMOVED***
	case "":
		return matchAttribute(n, t.key, func(string) bool ***REMOVED*** return true ***REMOVED***)
	case "=":
		return matchAttribute(n, t.key, func(s string) bool ***REMOVED*** return matchInsensitiveValue(s, t.val, t.insensitive) ***REMOVED***)
	case "!=":
		return attributeNotEqualMatch(t.key, t.val, n, t.insensitive)
	case "~=":
		// matches elements where the attribute named key is a whitespace-separated list that includes val.
		return matchAttribute(n, t.key, func(s string) bool ***REMOVED*** return matchInclude(t.val, s, t.insensitive) ***REMOVED***)
	case "|=":
		return attributeDashMatch(t.key, t.val, n, t.insensitive)
	case "^=":
		return attributePrefixMatch(t.key, t.val, n, t.insensitive)
	case "$=":
		return attributeSuffixMatch(t.key, t.val, n, t.insensitive)
	case "*=":
		return attributeSubstringMatch(t.key, t.val, n, t.insensitive)
	case "#=":
		return attributeRegexMatch(t.key, t.regexp, n)
	default:
		panic(fmt.Sprintf("unsuported operation : %s", t.operation))
	***REMOVED***
***REMOVED***

// matches elements where we ignore (or not) the case of the attribute value
// the user attribute is the value set by the user to match elements
// the real attribute is the attribute value found in the code parsed
func matchInsensitiveValue(userAttr string, realAttr string, ignoreCase bool) bool ***REMOVED***
	if ignoreCase ***REMOVED***
		return strings.EqualFold(userAttr, realAttr)
	***REMOVED***
	return userAttr == realAttr

***REMOVED***

// matches elements where the attribute named key satisifes the function f.
func matchAttribute(n *html.Node, key string, f func(string) bool) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	for _, a := range n.Attr ***REMOVED***
		if a.Key == key && f(a.Val) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// attributeNotEqualMatch matches elements where
// the attribute named key does not have the value val.
func attributeNotEqualMatch(key, val string, n *html.Node, ignoreCase bool) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	for _, a := range n.Attr ***REMOVED***
		if a.Key == key && matchInsensitiveValue(a.Val, val, ignoreCase) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// returns true if s is a whitespace-separated list that includes val.
func matchInclude(val string, s string, ignoreCase bool) bool ***REMOVED***
	for s != "" ***REMOVED***
		i := strings.IndexAny(s, " \t\r\n\f")
		if i == -1 ***REMOVED***
			return matchInsensitiveValue(s, val, ignoreCase)
		***REMOVED***
		if matchInsensitiveValue(s[:i], val, ignoreCase) ***REMOVED***
			return true
		***REMOVED***
		s = s[i+1:]
	***REMOVED***
	return false
***REMOVED***

//  matches elements where the attribute named key equals val or starts with val plus a hyphen.
func attributeDashMatch(key, val string, n *html.Node, ignoreCase bool) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if matchInsensitiveValue(s, val, ignoreCase) ***REMOVED***
				return true
			***REMOVED***
			if len(s) <= len(val) ***REMOVED***
				return false
			***REMOVED***
			if matchInsensitiveValue(s[:len(val)], val, ignoreCase) && s[len(val)] == '-' ***REMOVED***
				return true
			***REMOVED***
			return false
		***REMOVED***)
***REMOVED***

// attributePrefixMatch returns a Selector that matches elements where
// the attribute named key starts with val.
func attributePrefixMatch(key, val string, n *html.Node, ignoreCase bool) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			if ignoreCase ***REMOVED***
				return strings.HasPrefix(strings.ToLower(s), strings.ToLower(val))
			***REMOVED***
			return strings.HasPrefix(s, val)
		***REMOVED***)
***REMOVED***

// attributeSuffixMatch matches elements where
// the attribute named key ends with val.
func attributeSuffixMatch(key, val string, n *html.Node, ignoreCase bool) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			if ignoreCase ***REMOVED***
				return strings.HasSuffix(strings.ToLower(s), strings.ToLower(val))
			***REMOVED***
			return strings.HasSuffix(s, val)
		***REMOVED***)
***REMOVED***

// attributeSubstringMatch matches nodes where
// the attribute named key contains val.
func attributeSubstringMatch(key, val string, n *html.Node, ignoreCase bool) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			if ignoreCase ***REMOVED***
				return strings.Contains(strings.ToLower(s), strings.ToLower(val))
			***REMOVED***
			return strings.Contains(s, val)
		***REMOVED***)
***REMOVED***

// attributeRegexMatch  matches nodes where
// the attribute named key matches the regular expression rx
func attributeRegexMatch(key string, rx *regexp.Regexp, n *html.Node) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			return rx.MatchString(s)
		***REMOVED***)
***REMOVED***

func (c attrSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

func (c attrSelector) PseudoElement() string ***REMOVED***
	return ""
***REMOVED***

// see pseudo_classes.go for pseudo classes selectors

// on a static context, some selectors can't match anything
type neverMatchSelector struct ***REMOVED***
	value string
***REMOVED***

func (s neverMatchSelector) Match(n *html.Node) bool ***REMOVED***
	return false
***REMOVED***

func (s neverMatchSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 0, 0***REMOVED***
***REMOVED***

func (c neverMatchSelector) PseudoElement() string ***REMOVED***
	return ""
***REMOVED***

type compoundSelector struct ***REMOVED***
	selectors     []Sel
	pseudoElement string
***REMOVED***

// Matches elements if each sub-selectors matches.
func (t compoundSelector) Match(n *html.Node) bool ***REMOVED***
	if len(t.selectors) == 0 ***REMOVED***
		return n.Type == html.ElementNode
	***REMOVED***

	for _, sel := range t.selectors ***REMOVED***
		if !sel.Match(n) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (s compoundSelector) Specificity() Specificity ***REMOVED***
	var out Specificity
	for _, sel := range s.selectors ***REMOVED***
		out = out.Add(sel.Specificity())
	***REMOVED***
	if s.pseudoElement != "" ***REMOVED***
		// https://drafts.csswg.org/selectors-3/#specificity
		out = out.Add(Specificity***REMOVED***0, 0, 1***REMOVED***)
	***REMOVED***
	return out
***REMOVED***

func (c compoundSelector) PseudoElement() string ***REMOVED***
	return c.pseudoElement
***REMOVED***

type combinedSelector struct ***REMOVED***
	first      Sel
	combinator byte
	second     Sel
***REMOVED***

func (t combinedSelector) Match(n *html.Node) bool ***REMOVED***
	if t.first == nil ***REMOVED***
		return false // maybe we should panic
	***REMOVED***
	switch t.combinator ***REMOVED***
	case 0:
		return t.first.Match(n)
	case ' ':
		return descendantMatch(t.first, t.second, n)
	case '>':
		return childMatch(t.first, t.second, n)
	case '+':
		return siblingMatch(t.first, t.second, true, n)
	case '~':
		return siblingMatch(t.first, t.second, false, n)
	default:
		panic("unknown combinator")
	***REMOVED***
***REMOVED***

// matches an element if it matches d and has an ancestor that matches a.
func descendantMatch(a, d Matcher, n *html.Node) bool ***REMOVED***
	if !d.Match(n) ***REMOVED***
		return false
	***REMOVED***

	for p := n.Parent; p != nil; p = p.Parent ***REMOVED***
		if a.Match(p) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// matches an element if it matches d and its parent matches a.
func childMatch(a, d Matcher, n *html.Node) bool ***REMOVED***
	return d.Match(n) && n.Parent != nil && a.Match(n.Parent)
***REMOVED***

// matches an element if it matches s2 and is preceded by an element that matches s1.
// If adjacent is true, the sibling must be immediately before the element.
func siblingMatch(s1, s2 Matcher, adjacent bool, n *html.Node) bool ***REMOVED***
	if !s2.Match(n) ***REMOVED***
		return false
	***REMOVED***

	if adjacent ***REMOVED***
		for n = n.PrevSibling; n != nil; n = n.PrevSibling ***REMOVED***
			if n.Type == html.TextNode || n.Type == html.CommentNode ***REMOVED***
				continue
			***REMOVED***
			return s1.Match(n)
		***REMOVED***
		return false
	***REMOVED***

	// Walk backwards looking for element that matches s1
	for c := n.PrevSibling; c != nil; c = c.PrevSibling ***REMOVED***
		if s1.Match(c) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (s combinedSelector) Specificity() Specificity ***REMOVED***
	spec := s.first.Specificity()
	if s.second != nil ***REMOVED***
		spec = spec.Add(s.second.Specificity())
	***REMOVED***
	return spec
***REMOVED***

// on combinedSelector, a pseudo-element only makes sens on the last
// selector, although others increase specificity.
func (c combinedSelector) PseudoElement() string ***REMOVED***
	if c.second == nil ***REMOVED***
		return ""
	***REMOVED***
	return c.second.PseudoElement()
***REMOVED***

// A SelectorGroup is a list of selectors, which matches if any of the
// individual selectors matches.
type SelectorGroup []Sel

// Match returns true if the node matches one of the single selectors.
func (s SelectorGroup) Match(n *html.Node) bool ***REMOVED***
	for _, sel := range s ***REMOVED***
		if sel.Match(n) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
