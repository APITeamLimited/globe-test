package goquery

import (
	"strings"

	"golang.org/x/net/html"
)

// After applies the selector from the root document and inserts the matched elements
// after the elements in the set of matched elements.
//
// If one of the matched elements in the selection is not currently in the
// document, it's impossible to insert nodes after it, so it will be ignored.
//
// This follows the same rules as Selection.Append.
func (s *Selection) After(selector string) *Selection ***REMOVED***
	return s.AfterMatcher(compileMatcher(selector))
***REMOVED***

// AfterMatcher applies the matcher from the root document and inserts the matched elements
// after the elements in the set of matched elements.
//
// If one of the matched elements in the selection is not currently in the
// document, it's impossible to insert nodes after it, so it will be ignored.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AfterMatcher(m Matcher) *Selection ***REMOVED***
	return s.AfterNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// AfterSelection inserts the elements in the selection after each element in the set of matched
// elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AfterSelection(sel *Selection) *Selection ***REMOVED***
	return s.AfterNodes(sel.Nodes...)
***REMOVED***

// AfterHtml parses the html and inserts it after the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AfterHtml(htmlStr string) *Selection ***REMOVED***
	return s.eachNodeHtml(htmlStr, true, func(node *html.Node, nodes []*html.Node) ***REMOVED***
		nextSibling := node.NextSibling
		for _, n := range nodes ***REMOVED***
			if node.Parent != nil ***REMOVED***
				node.Parent.InsertBefore(n, nextSibling)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// AfterNodes inserts the nodes after each element in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AfterNodes(ns ...*html.Node) *Selection ***REMOVED***
	return s.manipulateNodes(ns, true, func(sn *html.Node, n *html.Node) ***REMOVED***
		if sn.Parent != nil ***REMOVED***
			sn.Parent.InsertBefore(n, sn.NextSibling)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Append appends the elements specified by the selector to the end of each element
// in the set of matched elements, following those rules:
//
// 1) The selector is applied to the root document.
//
// 2) Elements that are part of the document will be moved to the new location.
//
// 3) If there are multiple locations to append to, cloned nodes will be
// appended to all target locations except the last one, which will be moved
// as noted in (2).
func (s *Selection) Append(selector string) *Selection ***REMOVED***
	return s.AppendMatcher(compileMatcher(selector))
***REMOVED***

// AppendMatcher appends the elements specified by the matcher to the end of each element
// in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AppendMatcher(m Matcher) *Selection ***REMOVED***
	return s.AppendNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// AppendSelection appends the elements in the selection to the end of each element
// in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AppendSelection(sel *Selection) *Selection ***REMOVED***
	return s.AppendNodes(sel.Nodes...)
***REMOVED***

// AppendHtml parses the html and appends it to the set of matched elements.
func (s *Selection) AppendHtml(htmlStr string) *Selection ***REMOVED***
	return s.eachNodeHtml(htmlStr, false, func(node *html.Node, nodes []*html.Node) ***REMOVED***
		for _, n := range nodes ***REMOVED***
			node.AppendChild(n)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// AppendNodes appends the specified nodes to each node in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) AppendNodes(ns ...*html.Node) *Selection ***REMOVED***
	return s.manipulateNodes(ns, false, func(sn *html.Node, n *html.Node) ***REMOVED***
		sn.AppendChild(n)
	***REMOVED***)
***REMOVED***

// Before inserts the matched elements before each element in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) Before(selector string) *Selection ***REMOVED***
	return s.BeforeMatcher(compileMatcher(selector))
***REMOVED***

// BeforeMatcher inserts the matched elements before each element in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) BeforeMatcher(m Matcher) *Selection ***REMOVED***
	return s.BeforeNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// BeforeSelection inserts the elements in the selection before each element in the set of matched
// elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) BeforeSelection(sel *Selection) *Selection ***REMOVED***
	return s.BeforeNodes(sel.Nodes...)
***REMOVED***

// BeforeHtml parses the html and inserts it before the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) BeforeHtml(htmlStr string) *Selection ***REMOVED***
	return s.eachNodeHtml(htmlStr, true, func(node *html.Node, nodes []*html.Node) ***REMOVED***
		for _, n := range nodes ***REMOVED***
			if node.Parent != nil ***REMOVED***
				node.Parent.InsertBefore(n, node)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// BeforeNodes inserts the nodes before each element in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) BeforeNodes(ns ...*html.Node) *Selection ***REMOVED***
	return s.manipulateNodes(ns, false, func(sn *html.Node, n *html.Node) ***REMOVED***
		if sn.Parent != nil ***REMOVED***
			sn.Parent.InsertBefore(n, sn)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Clone creates a deep copy of the set of matched nodes. The new nodes will not be
// attached to the document.
func (s *Selection) Clone() *Selection ***REMOVED***
	ns := newEmptySelection(s.document)
	ns.Nodes = cloneNodes(s.Nodes)
	return ns
***REMOVED***

// Empty removes all children nodes from the set of matched elements.
// It returns the children nodes in a new Selection.
func (s *Selection) Empty() *Selection ***REMOVED***
	var nodes []*html.Node

	for _, n := range s.Nodes ***REMOVED***
		for c := n.FirstChild; c != nil; c = n.FirstChild ***REMOVED***
			n.RemoveChild(c)
			nodes = append(nodes, c)
		***REMOVED***
	***REMOVED***

	return pushStack(s, nodes)
***REMOVED***

// Prepend prepends the elements specified by the selector to each element in
// the set of matched elements, following the same rules as Append.
func (s *Selection) Prepend(selector string) *Selection ***REMOVED***
	return s.PrependMatcher(compileMatcher(selector))
***REMOVED***

// PrependMatcher prepends the elements specified by the matcher to each
// element in the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) PrependMatcher(m Matcher) *Selection ***REMOVED***
	return s.PrependNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// PrependSelection prepends the elements in the selection to each element in
// the set of matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) PrependSelection(sel *Selection) *Selection ***REMOVED***
	return s.PrependNodes(sel.Nodes...)
***REMOVED***

// PrependHtml parses the html and prepends it to the set of matched elements.
func (s *Selection) PrependHtml(htmlStr string) *Selection ***REMOVED***
	return s.eachNodeHtml(htmlStr, false, func(node *html.Node, nodes []*html.Node) ***REMOVED***
		firstChild := node.FirstChild
		for _, n := range nodes ***REMOVED***
			node.InsertBefore(n, firstChild)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// PrependNodes prepends the specified nodes to each node in the set of
// matched elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) PrependNodes(ns ...*html.Node) *Selection ***REMOVED***
	return s.manipulateNodes(ns, true, func(sn *html.Node, n *html.Node) ***REMOVED***
		// sn.FirstChild may be nil, in which case this functions like
		// sn.AppendChild()
		sn.InsertBefore(n, sn.FirstChild)
	***REMOVED***)
***REMOVED***

// Remove removes the set of matched elements from the document.
// It returns the same selection, now consisting of nodes not in the document.
func (s *Selection) Remove() *Selection ***REMOVED***
	for _, n := range s.Nodes ***REMOVED***
		if n.Parent != nil ***REMOVED***
			n.Parent.RemoveChild(n)
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

// RemoveFiltered removes from the current set of matched elements those that
// match the selector filter. It returns the Selection of removed nodes.
//
// For example if the selection s contains "<h1>", "<h2>" and "<h3>"
// and s.RemoveFiltered("h2") is called, only the "<h2>" node is removed
// (and returned), while "<h1>" and "<h3>" are kept in the document.
func (s *Selection) RemoveFiltered(selector string) *Selection ***REMOVED***
	return s.RemoveMatcher(compileMatcher(selector))
***REMOVED***

// RemoveMatcher removes from the current set of matched elements those that
// match the Matcher filter. It returns the Selection of removed nodes.
// See RemoveFiltered for additional information.
func (s *Selection) RemoveMatcher(m Matcher) *Selection ***REMOVED***
	return s.FilterMatcher(m).Remove()
***REMOVED***

// ReplaceWith replaces each element in the set of matched elements with the
// nodes matched by the given selector.
// It returns the removed elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) ReplaceWith(selector string) *Selection ***REMOVED***
	return s.ReplaceWithMatcher(compileMatcher(selector))
***REMOVED***

// ReplaceWithMatcher replaces each element in the set of matched elements with
// the nodes matched by the given Matcher.
// It returns the removed elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) ReplaceWithMatcher(m Matcher) *Selection ***REMOVED***
	return s.ReplaceWithNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// ReplaceWithSelection replaces each element in the set of matched elements with
// the nodes from the given Selection.
// It returns the removed elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) ReplaceWithSelection(sel *Selection) *Selection ***REMOVED***
	return s.ReplaceWithNodes(sel.Nodes...)
***REMOVED***

// ReplaceWithHtml replaces each element in the set of matched elements with
// the parsed HTML.
// It returns the removed elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) ReplaceWithHtml(htmlStr string) *Selection ***REMOVED***
	s.eachNodeHtml(htmlStr, true, func(node *html.Node, nodes []*html.Node) ***REMOVED***
		nextSibling := node.NextSibling
		for _, n := range nodes ***REMOVED***
			if node.Parent != nil ***REMOVED***
				node.Parent.InsertBefore(n, nextSibling)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return s.Remove()
***REMOVED***

// ReplaceWithNodes replaces each element in the set of matched elements with
// the given nodes.
// It returns the removed elements.
//
// This follows the same rules as Selection.Append.
func (s *Selection) ReplaceWithNodes(ns ...*html.Node) *Selection ***REMOVED***
	s.AfterNodes(ns...)
	return s.Remove()
***REMOVED***

// SetHtml sets the html content of each element in the selection to
// specified html string.
func (s *Selection) SetHtml(htmlStr string) *Selection ***REMOVED***
	for _, context := range s.Nodes ***REMOVED***
		for c := context.FirstChild; c != nil; c = context.FirstChild ***REMOVED***
			context.RemoveChild(c)
		***REMOVED***
	***REMOVED***
	return s.eachNodeHtml(htmlStr, false, func(node *html.Node, nodes []*html.Node) ***REMOVED***
		for _, n := range nodes ***REMOVED***
			node.AppendChild(n)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// SetText sets the content of each element in the selection to specified content.
// The provided text string is escaped.
func (s *Selection) SetText(text string) *Selection ***REMOVED***
	return s.SetHtml(html.EscapeString(text))
***REMOVED***

// Unwrap removes the parents of the set of matched elements, leaving the matched
// elements (and their siblings, if any) in their place.
// It returns the original selection.
func (s *Selection) Unwrap() *Selection ***REMOVED***
	s.Parent().Each(func(i int, ss *Selection) ***REMOVED***
		// For some reason, jquery allows unwrap to remove the <head> element, so
		// allowing it here too. Same for <html>. Why it allows those elements to
		// be unwrapped while not allowing body is a mystery to me.
		if ss.Nodes[0].Data != "body" ***REMOVED***
			ss.ReplaceWithSelection(ss.Contents())
		***REMOVED***
	***REMOVED***)

	return s
***REMOVED***

// Wrap wraps each element in the set of matched elements inside the first
// element matched by the given selector. The matched child is cloned before
// being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) Wrap(selector string) *Selection ***REMOVED***
	return s.WrapMatcher(compileMatcher(selector))
***REMOVED***

// WrapMatcher wraps each element in the set of matched elements inside the
// first element matched by the given matcher. The matched child is cloned
// before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapMatcher(m Matcher) *Selection ***REMOVED***
	return s.wrapNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// WrapSelection wraps each element in the set of matched elements inside the
// first element in the given Selection. The element is cloned before being
// inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapSelection(sel *Selection) *Selection ***REMOVED***
	return s.wrapNodes(sel.Nodes...)
***REMOVED***

// WrapHtml wraps each element in the set of matched elements inside the inner-
// most child of the given HTML.
//
// It returns the original set of elements.
func (s *Selection) WrapHtml(htmlStr string) *Selection ***REMOVED***
	nodesMap := make(map[string][]*html.Node)
	for _, context := range s.Nodes ***REMOVED***
		var parent *html.Node
		if context.Parent != nil ***REMOVED***
			parent = context.Parent
		***REMOVED*** else ***REMOVED***
			parent = &html.Node***REMOVED***Type: html.ElementNode***REMOVED***
		***REMOVED***
		nodes, found := nodesMap[nodeName(parent)]
		if !found ***REMOVED***
			nodes = parseHtmlWithContext(htmlStr, parent)
			nodesMap[nodeName(parent)] = nodes
		***REMOVED***
		newSingleSelection(context, s.document).wrapAllNodes(cloneNodes(nodes)...)
	***REMOVED***
	return s
***REMOVED***

// WrapNode wraps each element in the set of matched elements inside the inner-
// most child of the given node. The given node is copied before being inserted
// into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapNode(n *html.Node) *Selection ***REMOVED***
	return s.wrapNodes(n)
***REMOVED***

func (s *Selection) wrapNodes(ns ...*html.Node) *Selection ***REMOVED***
	s.Each(func(i int, ss *Selection) ***REMOVED***
		ss.wrapAllNodes(ns...)
	***REMOVED***)

	return s
***REMOVED***

// WrapAll wraps a single HTML structure, matched by the given selector, around
// all elements in the set of matched elements. The matched child is cloned
// before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapAll(selector string) *Selection ***REMOVED***
	return s.WrapAllMatcher(compileMatcher(selector))
***REMOVED***

// WrapAllMatcher wraps a single HTML structure, matched by the given Matcher,
// around all elements in the set of matched elements. The matched child is
// cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapAllMatcher(m Matcher) *Selection ***REMOVED***
	return s.wrapAllNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// WrapAllSelection wraps a single HTML structure, the first node of the given
// Selection, around all elements in the set of matched elements. The matched
// child is cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapAllSelection(sel *Selection) *Selection ***REMOVED***
	return s.wrapAllNodes(sel.Nodes...)
***REMOVED***

// WrapAllHtml wraps the given HTML structure around all elements in the set of
// matched elements. The matched child is cloned before being inserted into the
// document.
//
// It returns the original set of elements.
func (s *Selection) WrapAllHtml(htmlStr string) *Selection ***REMOVED***
	var context *html.Node
	var nodes []*html.Node
	if len(s.Nodes) > 0 ***REMOVED***
		context = s.Nodes[0]
		if context.Parent != nil ***REMOVED***
			nodes = parseHtmlWithContext(htmlStr, context)
		***REMOVED*** else ***REMOVED***
			nodes = parseHtml(htmlStr)
		***REMOVED***
	***REMOVED***
	return s.wrapAllNodes(nodes...)
***REMOVED***

func (s *Selection) wrapAllNodes(ns ...*html.Node) *Selection ***REMOVED***
	if len(ns) > 0 ***REMOVED***
		return s.WrapAllNode(ns[0])
	***REMOVED***
	return s
***REMOVED***

// WrapAllNode wraps the given node around the first element in the Selection,
// making all other nodes in the Selection children of the given node. The node
// is cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapAllNode(n *html.Node) *Selection ***REMOVED***
	if s.Size() == 0 ***REMOVED***
		return s
	***REMOVED***

	wrap := cloneNode(n)

	first := s.Nodes[0]
	if first.Parent != nil ***REMOVED***
		first.Parent.InsertBefore(wrap, first)
		first.Parent.RemoveChild(first)
	***REMOVED***

	for c := getFirstChildEl(wrap); c != nil; c = getFirstChildEl(wrap) ***REMOVED***
		wrap = c
	***REMOVED***

	newSingleSelection(wrap, s.document).AppendSelection(s)

	return s
***REMOVED***

// WrapInner wraps an HTML structure, matched by the given selector, around the
// content of element in the set of matched elements. The matched child is
// cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapInner(selector string) *Selection ***REMOVED***
	return s.WrapInnerMatcher(compileMatcher(selector))
***REMOVED***

// WrapInnerMatcher wraps an HTML structure, matched by the given selector,
// around the content of element in the set of matched elements. The matched
// child is cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapInnerMatcher(m Matcher) *Selection ***REMOVED***
	return s.wrapInnerNodes(m.MatchAll(s.document.rootNode)...)
***REMOVED***

// WrapInnerSelection wraps an HTML structure, matched by the given selector,
// around the content of element in the set of matched elements. The matched
// child is cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapInnerSelection(sel *Selection) *Selection ***REMOVED***
	return s.wrapInnerNodes(sel.Nodes...)
***REMOVED***

// WrapInnerHtml wraps an HTML structure, matched by the given selector, around
// the content of element in the set of matched elements. The matched child is
// cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapInnerHtml(htmlStr string) *Selection ***REMOVED***
	nodesMap := make(map[string][]*html.Node)
	for _, context := range s.Nodes ***REMOVED***
		nodes, found := nodesMap[nodeName(context)]
		if !found ***REMOVED***
			nodes = parseHtmlWithContext(htmlStr, context)
			nodesMap[nodeName(context)] = nodes
		***REMOVED***
		newSingleSelection(context, s.document).wrapInnerNodes(cloneNodes(nodes)...)
	***REMOVED***
	return s
***REMOVED***

// WrapInnerNode wraps an HTML structure, matched by the given selector, around
// the content of element in the set of matched elements. The matched child is
// cloned before being inserted into the document.
//
// It returns the original set of elements.
func (s *Selection) WrapInnerNode(n *html.Node) *Selection ***REMOVED***
	return s.wrapInnerNodes(n)
***REMOVED***

func (s *Selection) wrapInnerNodes(ns ...*html.Node) *Selection ***REMOVED***
	if len(ns) == 0 ***REMOVED***
		return s
	***REMOVED***

	s.Each(func(i int, s *Selection) ***REMOVED***
		contents := s.Contents()

		if contents.Size() > 0 ***REMOVED***
			contents.wrapAllNodes(ns...)
		***REMOVED*** else ***REMOVED***
			s.AppendNodes(cloneNode(ns[0]))
		***REMOVED***
	***REMOVED***)

	return s
***REMOVED***

func parseHtml(h string) []*html.Node ***REMOVED***
	// Errors are only returned when the io.Reader returns any error besides
	// EOF, but strings.Reader never will
	nodes, err := html.ParseFragment(strings.NewReader(h), &html.Node***REMOVED***Type: html.ElementNode***REMOVED***)
	if err != nil ***REMOVED***
		panic("goquery: failed to parse HTML: " + err.Error())
	***REMOVED***
	return nodes
***REMOVED***

func parseHtmlWithContext(h string, context *html.Node) []*html.Node ***REMOVED***
	// Errors are only returned when the io.Reader returns any error besides
	// EOF, but strings.Reader never will
	nodes, err := html.ParseFragment(strings.NewReader(h), context)
	if err != nil ***REMOVED***
		panic("goquery: failed to parse HTML: " + err.Error())
	***REMOVED***
	return nodes
***REMOVED***

// Get the first child that is an ElementNode
func getFirstChildEl(n *html.Node) *html.Node ***REMOVED***
	c := n.FirstChild
	for c != nil && c.Type != html.ElementNode ***REMOVED***
		c = c.NextSibling
	***REMOVED***
	return c
***REMOVED***

// Deep copy a slice of nodes.
func cloneNodes(ns []*html.Node) []*html.Node ***REMOVED***
	cns := make([]*html.Node, 0, len(ns))

	for _, n := range ns ***REMOVED***
		cns = append(cns, cloneNode(n))
	***REMOVED***

	return cns
***REMOVED***

// Deep copy a node. The new node has clones of all the original node's
// children but none of its parents or siblings.
func cloneNode(n *html.Node) *html.Node ***REMOVED***
	nn := &html.Node***REMOVED***
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	***REMOVED***

	copy(nn.Attr, n.Attr)
	for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
		nn.AppendChild(cloneNode(c))
	***REMOVED***

	return nn
***REMOVED***

func (s *Selection) manipulateNodes(ns []*html.Node, reverse bool,
	f func(sn *html.Node, n *html.Node)) *Selection ***REMOVED***

	lasti := s.Size() - 1

	// net.Html doesn't provide document fragments for insertion, so to get
	// things in the correct order with After() and Prepend(), the callback
	// needs to be called on the reverse of the nodes.
	if reverse ***REMOVED***
		for i, j := 0, len(ns)-1; i < j; i, j = i+1, j-1 ***REMOVED***
			ns[i], ns[j] = ns[j], ns[i]
		***REMOVED***
	***REMOVED***

	for i, sn := range s.Nodes ***REMOVED***
		for _, n := range ns ***REMOVED***
			if i != lasti ***REMOVED***
				f(sn, cloneNode(n))
			***REMOVED*** else ***REMOVED***
				if n.Parent != nil ***REMOVED***
					n.Parent.RemoveChild(n)
				***REMOVED***
				f(sn, n)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

// eachNodeHtml parses the given html string and inserts the resulting nodes in the dom with the mergeFn.
// The parsed nodes are inserted for each element of the selection.
// isParent can be used to indicate that the elements of the selection should be treated as the parent for the parsed html.
// A cache is used to avoid parsing the html multiple times should the elements of the selection result in the same context.
func (s *Selection) eachNodeHtml(htmlStr string, isParent bool, mergeFn func(n *html.Node, nodes []*html.Node)) *Selection ***REMOVED***
	// cache to avoid parsing the html for the same context multiple times
	nodeCache := make(map[string][]*html.Node)
	var context *html.Node
	for _, n := range s.Nodes ***REMOVED***
		if isParent ***REMOVED***
			context = n.Parent
		***REMOVED*** else ***REMOVED***
			if n.Type != html.ElementNode ***REMOVED***
				continue
			***REMOVED***
			context = n
		***REMOVED***
		if context != nil ***REMOVED***
			nodes, found := nodeCache[nodeName(context)]
			if !found ***REMOVED***
				nodes = parseHtmlWithContext(htmlStr, context)
				nodeCache[nodeName(context)] = nodes
			***REMOVED***
			mergeFn(n, cloneNodes(nodes))
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***
