package cascadia

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// This file implements the pseudo classes selectors,
// which share the implementation of PseudoElement() and Specificity()

type abstractPseudoClass struct***REMOVED******REMOVED***

func (s abstractPseudoClass) Specificity() Specificity ***REMOVED***
	return Specificity***REMOVED***0, 1, 0***REMOVED***
***REMOVED***

func (c abstractPseudoClass) PseudoElement() string ***REMOVED***
	return ""
***REMOVED***

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

func (c relativePseudoClassSelector) PseudoElement() string ***REMOVED***
	return ""
***REMOVED***

type containsPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
	value string
	own   bool
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

type regexpPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
	regexp *regexp.Regexp
	own    bool
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

type nthPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
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

type onlyChildPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
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

type inputPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
***REMOVED***

// Matches input, select, textarea and button elements.
func (s inputPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	return n.Type == html.ElementNode && (n.Data == "input" || n.Data == "select" || n.Data == "textarea" || n.Data == "button")
***REMOVED***

type emptyElementPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
***REMOVED***

// Matches empty elements.
func (s emptyElementPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***

	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		switch c.Type ***REMOVED***
		case html.ElementNode:
			return false
		case html.TextNode:
			if strings.TrimSpace(nodeText(c)) == "" ***REMOVED***
				continue
			***REMOVED*** else ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

type rootPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
***REMOVED***

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

func hasAttr(n *html.Node, attr string) bool ***REMOVED***
	return matchAttribute(n, attr, func(string) bool ***REMOVED*** return true ***REMOVED***)
***REMOVED***

type linkPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
***REMOVED***

// Match implements :link
func (s linkPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	return (n.DataAtom == atom.A || n.DataAtom == atom.Area || n.DataAtom == atom.Link) && hasAttr(n, "href")
***REMOVED***

type langPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
	lang string
***REMOVED***

func (s langPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	own := matchAttribute(n, "lang", func(val string) bool ***REMOVED***
		return val == s.lang || strings.HasPrefix(val, s.lang+"-")
	***REMOVED***)
	if n.Parent == nil ***REMOVED***
		return own
	***REMOVED***
	return own || s.Match(n.Parent)
***REMOVED***

type enabledPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
***REMOVED***

func (s enabledPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	switch n.DataAtom ***REMOVED***
	case atom.A, atom.Area, atom.Link:
		return hasAttr(n, "href")
	case atom.Optgroup, atom.Menuitem, atom.Fieldset:
		return !hasAttr(n, "disabled")
	case atom.Button, atom.Input, atom.Select, atom.Textarea, atom.Option:
		return !hasAttr(n, "disabled") && !inDisabledFieldset(n)
	***REMOVED***
	return false
***REMOVED***

type disabledPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
***REMOVED***

func (s disabledPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	switch n.DataAtom ***REMOVED***
	case atom.Optgroup, atom.Menuitem, atom.Fieldset:
		return hasAttr(n, "disabled")
	case atom.Button, atom.Input, atom.Select, atom.Textarea, atom.Option:
		return hasAttr(n, "disabled") || inDisabledFieldset(n)
	***REMOVED***
	return false
***REMOVED***

func hasLegendInPreviousSiblings(n *html.Node) bool ***REMOVED***
	for s := n.PrevSibling; s != nil; s = s.PrevSibling ***REMOVED***
		if s.DataAtom == atom.Legend ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func inDisabledFieldset(n *html.Node) bool ***REMOVED***
	if n.Parent == nil ***REMOVED***
		return false
	***REMOVED***
	if n.Parent.DataAtom == atom.Fieldset && hasAttr(n.Parent, "disabled") &&
		(n.DataAtom != atom.Legend || hasLegendInPreviousSiblings(n)) ***REMOVED***
		return true
	***REMOVED***
	return inDisabledFieldset(n.Parent)
***REMOVED***

type checkedPseudoClassSelector struct ***REMOVED***
	abstractPseudoClass
***REMOVED***

func (s checkedPseudoClassSelector) Match(n *html.Node) bool ***REMOVED***
	if n.Type != html.ElementNode ***REMOVED***
		return false
	***REMOVED***
	switch n.DataAtom ***REMOVED***
	case atom.Input, atom.Menuitem:
		return hasAttr(n, "checked") && matchAttribute(n, "type", func(val string) bool ***REMOVED***
			t := toLowerASCII(val)
			return t == "checkbox" || t == "radio"
		***REMOVED***)
	case atom.Option:
		return hasAttr(n, "selected")
	***REMOVED***
	return false
***REMOVED***
