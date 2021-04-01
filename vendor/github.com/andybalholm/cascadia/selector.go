package cascadia

import (
	"bytes"
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
// It is currently the same as Matcher, but other methods may be added in the
// future.
type Sel interface ***REMOVED***
	Matcher
	Specificity() Specificity
***REMOVED***

// Parse parses a selector.
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

// ParseGroup parses a selector, or a group of selectors separated by commas.
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

type classSelector struct ***REMOVED***
	class string
***REMOVED***

// Matches elements by class attribute.
func (t classSelector) Match(n *html.Node) bool ***REMOVED***
	return matchAttribute(n, "class", func(s string) bool ***REMOVED***
		return matchInclude(t.class, s)
	***REMOVED***)
***REMOVED***

func (c classSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
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

type attrSelector struct ***REMOVED***
	key, val, operation string
	regexp              *regexp.Regexp
***REMOVED***

// Matches elements by attribute value.
func (t attrSelector) Match(n *html.Node) bool ***REMOVED***
	switch t.operation ***REMOVED***
	case "":
		return matchAttribute(n, t.key, func(string) bool ***REMOVED*** return true ***REMOVED***)
	case "=":
		return matchAttribute(n, t.key, func(s string) bool ***REMOVED*** return s == t.val ***REMOVED***)
	case "!=":
		return attributeNotEqualMatch(t.key, t.val, n)
	case "~=":
		// matches elements where the attribute named key is a whitespace-separated list that includes val.
		return matchAttribute(n, t.key, func(s string) bool ***REMOVED*** return matchInclude(t.val, s) ***REMOVED***)
	case "|=":
		return attributeDashMatch(t.key, t.val, n)
	case "^=":
		return attributePrefixMatch(t.key, t.val, n)
	case "$=":
		return attributeSuffixMatch(t.key, t.val, n)
	case "*=":
		return attributeSubstringMatch(t.key, t.val, n)
	case "#=":
		return attributeRegexMatch(t.key, t.regexp, n)
	default:
		panic(fmt.Sprintf("unsuported operation : %s", t.operation))
	***REMOVED***
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
func attributeNotEqualMatch(key, val string, n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	for _, a := range n.Attr ***REMOVED***
		if a.Key == key && a.Val == val ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// returns true if s is a whitespace-separated list that includes val.
func matchInclude(val, s string) bool ***REMOVED***
	for s != "" ***REMOVED***
		i := strings.IndexAny(s, " \t\r\n\f")
		if i == -1 ***REMOVED***
			return s == val
		***REMOVED***
		if s[:i] == val ***REMOVED***
			return true
		***REMOVED***
		s = s[i+1:]
	***REMOVED***
	return false
***REMOVED***

//  matches elements where the attribute named key equals val or starts with val plus a hyphen.
func attributeDashMatch(key, val string, n *html.Node) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if s == val ***REMOVED***
				return true
			***REMOVED***
			if len(s) <= len(val) ***REMOVED***
				return false
			***REMOVED***
			if s[:len(val)] == val && s[len(val)] == '-' ***REMOVED***
				return true
			***REMOVED***
			return false
		***REMOVED***)
***REMOVED***

// attributePrefixMatch returns a Selector that matches elements where
// the attribute named key starts with val.
func attributePrefixMatch(key, val string, n *html.Node) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			return strings.HasPrefix(s, val)
		***REMOVED***)
***REMOVED***

// attributeSuffixMatch matches elements where
// the attribute named key ends with val.
func attributeSuffixMatch(key, val string, n *html.Node) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			return strings.HasSuffix(s, val)
		***REMOVED***)
***REMOVED***

// attributeSubstringMatch matches nodes where
// the attribute named key contains val.
func attributeSubstringMatch(key, val string, n *html.Node) bool ***REMOVED***
	return matchAttribute(n, key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
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

// ---------------- Pseudo class selectors ----------------
// we use severals concrete types of pseudo-class selectors

type relativePseudoClassSelector struct ***REMOVED***
	name  string // one of "not", "has", "haschild"
	match SelectorGroup
***REMOVED***

func (s relativePseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	switch s.name ***REMOVED***
	case "not":
		// matches elements that do not match a.
		return !s.match.Match(n)
	case "has":
		//  matches elements with any descendant that matches a.
		return hasDescendantMatch(n, s.match)
	case "haschild":
		// matches elements with a child that matches a.
		return hasChildMatch(n, s.match)
	default:
		panic(fmt.Sprintf("unsupported relative pseudo class selector : %s", s.name))
	***REMOVED***
***REMOVED***

// hasChildMatch returns whether n has any child that matches a.
func hasChildMatch(n *html.Node, a Matcher) bool ***REMOVED***
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if a.Match(c) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// hasDescendantMatch performs a depth-first search of n's descendants,
// testing whether any of them match a. It returns true as soon as a match is
// found, or false if no match is found.
func hasDescendantMatch(n *html.Node, a Matcher) bool ***REMOVED***
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if a.Match(c) || (c.Type == html.ElementNode && hasDescendantMatch(c, a)) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Specificity returns the specificity of the most specific selectors
// in the pseudo-class arguments.
// See https://www.w3.org/TR/selectors/#specificity-rules
func (s relativePseudoClassSelector) Specificity() Specificity ***REMOVED***
	var max Specificity
	for _, sel := range s.match ***REMOVED***
		newSpe := sel.Specificity()
		if max.Less(newSpe) ***REMOVED***
			max = newSpe
		***REMOVED***
	***REMOVED***
	return max
***REMOVED***

type containsPseudoClassSelector struct ***REMOVED***
	own   bool
	value string
***REMOVED***

func (s containsPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	var text string
	if s.own ***REMOVED***
		// matches nodes that directly contain the given text
		text = strings.ToLower(nodeOwnText(n))
	***REMOVED*** else ***REMOVED***
		// matches nodes that contain the given text.
		text = strings.ToLower(nodeText(n))
	***REMOVED***
	return strings.Contains(text, s.value)
***REMOVED***

func (s containsPseudoClassSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

type regexpPseudoClassSelector struct ***REMOVED***
	own    bool
	regexp *regexp.Regexp
***REMOVED***

func (s regexpPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	var text string
	if s.own ***REMOVED***
		// matches nodes whose text directly matches the specified regular expression
		text = nodeOwnText(n)
	***REMOVED*** else ***REMOVED***
		// matches nodes whose text matches the specified regular expression
		text = nodeText(n)
	***REMOVED***
	return s.regexp.MatchString(text)
***REMOVED***

// writeNodeText writes the text contained in n and its descendants to b.
func writeNodeText(n *html.Node, b *bytes.Buffer) ***REMOVED***
	switch n.Type ***REMOVED***
	case html.TextNode:
		b.WriteString(n.Data)
	case html.ElementNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			writeNodeText(c, b)
		***REMOVED***
	***REMOVED***
***REMOVED***

// nodeText returns the text contained in n and its descendants.
func nodeText(n *html.Node) string ***REMOVED***
	var b bytes.Buffer
	writeNodeText(n, &b)
	return b.String()
***REMOVED***

// nodeOwnText returns the contents of the text nodes that are direct
// children of n.
func nodeOwnText(n *html.Node) string ***REMOVED***
	var b bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if c.Type == html.TextNode ***REMOVED***
			b.WriteString(c.Data)
		***REMOVED***
	***REMOVED***
	return b.String()
***REMOVED***

func (s regexpPseudoClassSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

type nthPseudoClassSelector struct ***REMOVED***
	a, b         int
	last, ofType bool
***REMOVED***

func (s nthPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if s.a == 0 ***REMOVED***
		if s.last ***REMOVED***
			return simpleNthLastChildMatch(s.b, s.ofType, n)
		***REMOVED*** else ***REMOVED***
			return simpleNthChildMatch(s.b, s.ofType, n)
		***REMOVED***
	***REMOVED***
	return nthChildMatch(s.a, s.b, s.last, s.ofType, n)
***REMOVED***

// nthChildMatch implements :nth-child(an+b).
// If last is true, implements :nth-last-child instead.
// If ofType is true, implements :nth-of-type instead.
func nthChildMatch(a, b int, last, ofType bool, n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***

	parent := n.Parent
	if parent == nil ***REMOVED***
		return false
	***REMOVED***

	if parent.Type == html.DocumentNode ***REMOVED***
		return false
	***REMOVED***

	i := -1
	count := 0
	for c := parent.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if (c.Type != html.ElementNode) || (ofType && c.Data != n.Data) ***REMOVED***
			continue
		***REMOVED***
		count++
		if c == n ***REMOVED***
			i = count
			if !last ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if i == -1 ***REMOVED***
		// This shouldn't happen, since n should always be one of its parent's children.
		return false
	***REMOVED***

	if last ***REMOVED***
		i = count - i + 1
	***REMOVED***

	i -= b
	if a == 0 ***REMOVED***
		return i == 0
	***REMOVED***

	return i%a == 0 && i/a >= 0
***REMOVED***

// simpleNthChildMatch implements :nth-child(b).
// If ofType is true, implements :nth-of-type instead.
func simpleNthChildMatch(b int, ofType bool, n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***

	parent := n.Parent
	if parent == nil ***REMOVED***
		return false
	***REMOVED***

	if parent.Type == html.DocumentNode ***REMOVED***
		return false
	***REMOVED***

	count := 0
	for c := parent.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if c.Type != html.ElementNode || (ofType && c.Data != n.Data) ***REMOVED***
			continue
		***REMOVED***
		count++
		if c == n ***REMOVED***
			return count == b
		***REMOVED***
		if count >= b ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// simpleNthLastChildMatch implements :nth-last-child(b).
// If ofType is true, implements :nth-last-of-type instead.
func simpleNthLastChildMatch(b int, ofType bool, n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***

	parent := n.Parent
	if parent == nil ***REMOVED***
		return false
	***REMOVED***

	if parent.Type == html.DocumentNode ***REMOVED***
		return false
	***REMOVED***

	count := 0
	for c := parent.LastChild; c != nil; c = c.PrevSibling ***REMOVED***
		if c.Type != html.ElementNode || (ofType && c.Data != n.Data) ***REMOVED***
			continue
		***REMOVED***
		count++
		if c == n ***REMOVED***
			return count == b
		***REMOVED***
		if count >= b ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Specificity for nth-child pseudo-class.
// Does not support a list of selectors
func (s nthPseudoClassSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

type onlyChildPseudoClassSelector struct ***REMOVED***
	ofType bool
***REMOVED***

// Match implements :only-child.
// If `ofType` is true, it implements :only-of-type instead.
func (s onlyChildPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***

	parent := n.Parent
	if parent == nil ***REMOVED***
		return false
	***REMOVED***

	if parent.Type == html.DocumentNode ***REMOVED***
		return false
	***REMOVED***

	count := 0
	for c := parent.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if (c.Type != html.ElementNode) || (s.ofType && c.Data != n.Data) ***REMOVED***
			continue
		***REMOVED***
		count++
		if count > 1 ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return count == 1
***REMOVED***

func (s onlyChildPseudoClassSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

type inputPseudoClassSelector struct***REMOVED******REMOVED***

// Matches input, select, textarea and button elements.
func (s inputPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	return n.Type == html.ElementNode && (n.Data == "input" || n.Data == "select" || n.Data == "textarea" || n.Data == "button")
***REMOVED***

func (s inputPseudoClassSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

type emptyElementPseudoClassSelector struct***REMOVED******REMOVED***

// Matches empty elements.
func (s emptyElementPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***

	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		switch c.Type ***REMOVED***
		case html.ElementNode, html.TextNode:
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func (s emptyElementPseudoClassSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

type rootPseudoClassSelector struct***REMOVED******REMOVED***

// Match implements :root
func (s rootPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	if n.Parent == nil ***REMOVED***
		return false
	***REMOVED***
	return n.Parent.Type == html.DocumentNode
***REMOVED***

func (s rootPseudoClassSelector) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

type compoundSelector struct ***REMOVED***
	selectors []Sel
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
	return out
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
