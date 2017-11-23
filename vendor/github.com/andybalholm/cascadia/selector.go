package cascadia

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// the Selector type, and functions for creating them

// A Selector is a function which tells whether a node matches or not.
type Selector func(*html.Node) bool

// hasChildMatch returns whether n has any child that matches a.
func hasChildMatch(n *html.Node, a Selector) bool ***REMOVED***
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if a(c) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// hasDescendantMatch performs a depth-first search of n's descendants,
// testing whether any of them match a. It returns true as soon as a match is
// found, or false if no match is found.
func hasDescendantMatch(n *html.Node, a Selector) bool ***REMOVED***
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		if a(c) || (c.Type == html.ElementNode && hasDescendantMatch(c, a)) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Compile parses a selector and returns, if successful, a Selector object
// that can be used to match against html.Node objects.
func Compile(sel string) (Selector, error) ***REMOVED***
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

// Filter returns the nodes in nodes that match the selector.
func (s Selector) Filter(nodes []*html.Node) (result []*html.Node) ***REMOVED***
	for _, n := range nodes ***REMOVED***
		if s(n) ***REMOVED***
			result = append(result, n)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// typeSelector returns a Selector that matches elements with a given tag name.
func typeSelector(tag string) Selector ***REMOVED***
	tag = toLowerASCII(tag)
	return func(n *html.Node) bool ***REMOVED***
		return n.Type == html.ElementNode && n.Data == tag
	***REMOVED***
***REMOVED***

// toLowerASCII returns s with all ASCII capital letters lowercased.
func toLowerASCII(s string) string ***REMOVED***
	var b []byte
	for i := 0; i < len(s); i++ ***REMOVED***
		if c := s[i]; 'A' <= c && c <= 'Z' ***REMOVED***
			if b == nil ***REMOVED***
				b = make([]byte, len(s))
				copy(b, s)
			***REMOVED***
			b[i] = s[i] + ('a' - 'A')
		***REMOVED***
	***REMOVED***

	if b == nil ***REMOVED***
		return s
	***REMOVED***

	return string(b)
***REMOVED***

// attributeSelector returns a Selector that matches elements
// where the attribute named key satisifes the function f.
func attributeSelector(key string, f func(string) bool) Selector ***REMOVED***
	key = toLowerASCII(key)
	return func(n *html.Node) bool ***REMOVED***
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
***REMOVED***

// attributeExistsSelector returns a Selector that matches elements that have
// an attribute named key.
func attributeExistsSelector(key string) Selector ***REMOVED***
	return attributeSelector(key, func(string) bool ***REMOVED*** return true ***REMOVED***)
***REMOVED***

// attributeEqualsSelector returns a Selector that matches elements where
// the attribute named key has the value val.
func attributeEqualsSelector(key, val string) Selector ***REMOVED***
	return attributeSelector(key,
		func(s string) bool ***REMOVED***
			return s == val
		***REMOVED***)
***REMOVED***

// attributeNotEqualSelector returns a Selector that matches elements where
// the attribute named key does not have the value val.
func attributeNotEqualSelector(key, val string) Selector ***REMOVED***
	key = toLowerASCII(key)
	return func(n *html.Node) bool ***REMOVED***
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
***REMOVED***

// attributeIncludesSelector returns a Selector that matches elements where
// the attribute named key is a whitespace-separated list that includes val.
func attributeIncludesSelector(key, val string) Selector ***REMOVED***
	return attributeSelector(key,
		func(s string) bool ***REMOVED***
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
		***REMOVED***)
***REMOVED***

// attributeDashmatchSelector returns a Selector that matches elements where
// the attribute named key equals val or starts with val plus a hyphen.
func attributeDashmatchSelector(key, val string) Selector ***REMOVED***
	return attributeSelector(key,
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

// attributePrefixSelector returns a Selector that matches elements where
// the attribute named key starts with val.
func attributePrefixSelector(key, val string) Selector ***REMOVED***
	return attributeSelector(key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			return strings.HasPrefix(s, val)
		***REMOVED***)
***REMOVED***

// attributeSuffixSelector returns a Selector that matches elements where
// the attribute named key ends with val.
func attributeSuffixSelector(key, val string) Selector ***REMOVED***
	return attributeSelector(key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			return strings.HasSuffix(s, val)
		***REMOVED***)
***REMOVED***

// attributeSubstringSelector returns a Selector that matches nodes where
// the attribute named key contains val.
func attributeSubstringSelector(key, val string) Selector ***REMOVED***
	return attributeSelector(key,
		func(s string) bool ***REMOVED***
			if strings.TrimSpace(s) == "" ***REMOVED***
				return false
			***REMOVED***
			return strings.Contains(s, val)
		***REMOVED***)
***REMOVED***

// attributeRegexSelector returns a Selector that matches nodes where
// the attribute named key matches the regular expression rx
func attributeRegexSelector(key string, rx *regexp.Regexp) Selector ***REMOVED***
	return attributeSelector(key,
		func(s string) bool ***REMOVED***
			return rx.MatchString(s)
		***REMOVED***)
***REMOVED***

// intersectionSelector returns a selector that matches nodes that match
// both a and b.
func intersectionSelector(a, b Selector) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		return a(n) && b(n)
	***REMOVED***
***REMOVED***

// unionSelector returns a selector that matches elements that match
// either a or b.
func unionSelector(a, b Selector) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		return a(n) || b(n)
	***REMOVED***
***REMOVED***

// negatedSelector returns a selector that matches elements that do not match a.
func negatedSelector(a Selector) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		if n.Type != html.ElementNode ***REMOVED***
			return false
		***REMOVED***
		return !a(n)
	***REMOVED***
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

// textSubstrSelector returns a selector that matches nodes that
// contain the given text.
func textSubstrSelector(val string) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		text := strings.ToLower(nodeText(n))
		return strings.Contains(text, val)
	***REMOVED***
***REMOVED***

// ownTextSubstrSelector returns a selector that matches nodes that
// directly contain the given text
func ownTextSubstrSelector(val string) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		text := strings.ToLower(nodeOwnText(n))
		return strings.Contains(text, val)
	***REMOVED***
***REMOVED***

// textRegexSelector returns a selector that matches nodes whose text matches
// the specified regular expression
func textRegexSelector(rx *regexp.Regexp) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		return rx.MatchString(nodeText(n))
	***REMOVED***
***REMOVED***

// ownTextRegexSelector returns a selector that matches nodes whose text
// directly matches the specified regular expression
func ownTextRegexSelector(rx *regexp.Regexp) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		return rx.MatchString(nodeOwnText(n))
	***REMOVED***
***REMOVED***

// hasChildSelector returns a selector that matches elements
// with a child that matches a.
func hasChildSelector(a Selector) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		if n.Type != html.ElementNode ***REMOVED***
			return false
		***REMOVED***
		return hasChildMatch(n, a)
	***REMOVED***
***REMOVED***

// hasDescendantSelector returns a selector that matches elements
// with any descendant that matches a.
func hasDescendantSelector(a Selector) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		if n.Type != html.ElementNode ***REMOVED***
			return false
		***REMOVED***
		return hasDescendantMatch(n, a)
	***REMOVED***
***REMOVED***

// nthChildSelector returns a selector that implements :nth-child(an+b).
// If last is true, implements :nth-last-child instead.
// If ofType is true, implements :nth-of-type instead.
func nthChildSelector(a, b int, last, ofType bool) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
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
***REMOVED***

// simpleNthChildSelector returns a selector that implements :nth-child(b).
// If ofType is true, implements :nth-of-type instead.
func simpleNthChildSelector(b int, ofType bool) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
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
***REMOVED***

// simpleNthLastChildSelector returns a selector that implements
// :nth-last-child(b). If ofType is true, implements :nth-last-of-type
// instead.
func simpleNthLastChildSelector(b int, ofType bool) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
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
***REMOVED***

// onlyChildSelector returns a selector that implements :only-child.
// If ofType is true, it implements :only-of-type instead.
func onlyChildSelector(ofType bool) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
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
			if (c.Type != html.ElementNode) || (ofType && c.Data != n.Data) ***REMOVED***
				continue
			***REMOVED***
			count++
			if count > 1 ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		return count == 1
	***REMOVED***
***REMOVED***

// inputSelector is a Selector that matches input, select, textarea and button elements.
func inputSelector(n *html.Node) bool ***REMOVED***
	return n.Type == html.ElementNode && (n.Data == "input" || n.Data == "select" || n.Data == "textarea" || n.Data == "button")
***REMOVED***

// emptyElementSelector is a Selector that matches empty elements.
func emptyElementSelector(n *html.Node) bool ***REMOVED***
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

// descendantSelector returns a Selector that matches an element if
// it matches d and has an ancestor that matches a.
func descendantSelector(a, d Selector) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		if !d(n) ***REMOVED***
			return false
		***REMOVED***

		for p := n.Parent; p != nil; p = p.Parent ***REMOVED***
			if a(p) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***

		return false
	***REMOVED***
***REMOVED***

// childSelector returns a Selector that matches an element if
// it matches d and its parent matches a.
func childSelector(a, d Selector) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		return d(n) && n.Parent != nil && a(n.Parent)
	***REMOVED***
***REMOVED***

// siblingSelector returns a Selector that matches an element
// if it matches s2 and in is preceded by an element that matches s1.
// If adjacent is true, the sibling must be immediately before the element.
func siblingSelector(s1, s2 Selector, adjacent bool) Selector ***REMOVED***
	return func(n *html.Node) bool ***REMOVED***
		if !s2(n) ***REMOVED***
			return false
		***REMOVED***

		if adjacent ***REMOVED***
			for n = n.PrevSibling; n != nil; n = n.PrevSibling ***REMOVED***
				if n.Type == html.TextNode || n.Type == html.CommentNode ***REMOVED***
					continue
				***REMOVED***
				return s1(n)
			***REMOVED***
			return false
		***REMOVED***

		// Walk backwards looking for element that matches s1
		for c := n.PrevSibling; c != nil; c = c.PrevSibling ***REMOVED***
			if s1(c) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***

		return false
	***REMOVED***
***REMOVED***

// rootSelector implements :root
func rootSelector(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	if n.Parent == nil ***REMOVED***
		return false
	***REMOVED***
	return n.Parent.Type == html.DocumentNode
***REMOVED***
