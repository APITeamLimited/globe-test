package goquery

import "golang.org/x/net/html"

type siblingType int

// Sibling type, used internally when iterating over children at the same
// level (siblings) to specify which nodes are requested.
const (
	siblingPrevUntil siblingType = iota - 3
	siblingPrevAll
	siblingPrev
	siblingAll
	siblingNext
	siblingNextAll
	siblingNextUntil
	siblingAllIncludingNonElements
)

// Find gets the descendants of each element in the current set of matched
// elements, filtered by a selector. It returns a new Selection object
// containing these matched elements.
func (s *Selection) Find(selector string) *Selection ***REMOVED***
	return pushStack(s, findWithMatcher(s.Nodes, compileMatcher(selector)))
***REMOVED***

// FindMatcher gets the descendants of each element in the current set of matched
// elements, filtered by the matcher. It returns a new Selection object
// containing these matched elements.
func (s *Selection) FindMatcher(m Matcher) *Selection ***REMOVED***
	return pushStack(s, findWithMatcher(s.Nodes, m))
***REMOVED***

// FindSelection gets the descendants of each element in the current
// Selection, filtered by a Selection. It returns a new Selection object
// containing these matched elements.
func (s *Selection) FindSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return pushStack(s, nil)
	***REMOVED***
	return s.FindNodes(sel.Nodes...)
***REMOVED***

// FindNodes gets the descendants of each element in the current
// Selection, filtered by some nodes. It returns a new Selection object
// containing these matched elements.
func (s *Selection) FindNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return pushStack(s, mapNodes(nodes, func(i int, n *html.Node) []*html.Node ***REMOVED***
		if sliceContains(s.Nodes, n) ***REMOVED***
			return []*html.Node***REMOVED***n***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***))
***REMOVED***

// Contents gets the children of each element in the Selection,
// including text and comment nodes. It returns a new Selection object
// containing these elements.
func (s *Selection) Contents() *Selection ***REMOVED***
	return pushStack(s, getChildrenNodes(s.Nodes, siblingAllIncludingNonElements))
***REMOVED***

// ContentsFiltered gets the children of each element in the Selection,
// filtered by the specified selector. It returns a new Selection
// object containing these elements. Since selectors only act on Element nodes,
// this function is an alias to ChildrenFiltered unless the selector is empty,
// in which case it is an alias to Contents.
func (s *Selection) ContentsFiltered(selector string) *Selection ***REMOVED***
	if selector != "" ***REMOVED***
		return s.ChildrenFiltered(selector)
	***REMOVED***
	return s.Contents()
***REMOVED***

// ContentsMatcher gets the children of each element in the Selection,
// filtered by the specified matcher. It returns a new Selection
// object containing these elements. Since matchers only act on Element nodes,
// this function is an alias to ChildrenMatcher.
func (s *Selection) ContentsMatcher(m Matcher) *Selection ***REMOVED***
	return s.ChildrenMatcher(m)
***REMOVED***

// Children gets the child elements of each element in the Selection.
// It returns a new Selection object containing these elements.
func (s *Selection) Children() *Selection ***REMOVED***
	return pushStack(s, getChildrenNodes(s.Nodes, siblingAll))
***REMOVED***

// ChildrenFiltered gets the child elements of each element in the Selection,
// filtered by the specified selector. It returns a new
// Selection object containing these elements.
func (s *Selection) ChildrenFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getChildrenNodes(s.Nodes, siblingAll), compileMatcher(selector))
***REMOVED***

// ChildrenMatcher gets the child elements of each element in the Selection,
// filtered by the specified matcher. It returns a new
// Selection object containing these elements.
func (s *Selection) ChildrenMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getChildrenNodes(s.Nodes, siblingAll), m)
***REMOVED***

// Parent gets the parent of each element in the Selection. It returns a
// new Selection object containing the matched elements.
func (s *Selection) Parent() *Selection ***REMOVED***
	return pushStack(s, getParentNodes(s.Nodes))
***REMOVED***

// ParentFiltered gets the parent of each element in the Selection filtered by a
// selector. It returns a new Selection object containing the matched elements.
func (s *Selection) ParentFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getParentNodes(s.Nodes), compileMatcher(selector))
***REMOVED***

// ParentMatcher gets the parent of each element in the Selection filtered by a
// matcher. It returns a new Selection object containing the matched elements.
func (s *Selection) ParentMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getParentNodes(s.Nodes), m)
***REMOVED***

// Closest gets the first element that matches the selector by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (s *Selection) Closest(selector string) *Selection ***REMOVED***
	cs := compileMatcher(selector)
	return s.ClosestMatcher(cs)
***REMOVED***

// ClosestMatcher gets the first element that matches the matcher by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (s *Selection) ClosestMatcher(m Matcher) *Selection ***REMOVED***
	return pushStack(s, mapNodes(s.Nodes, func(i int, n *html.Node) []*html.Node ***REMOVED***
		// For each node in the selection, test the node itself, then each parent
		// until a match is found.
		for ; n != nil; n = n.Parent ***REMOVED***
			if m.Match(n) ***REMOVED***
				return []*html.Node***REMOVED***n***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***))
***REMOVED***

// ClosestNodes gets the first element that matches one of the nodes by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (s *Selection) ClosestNodes(nodes ...*html.Node) *Selection ***REMOVED***
	set := make(map[*html.Node]bool)
	for _, n := range nodes ***REMOVED***
		set[n] = true
	***REMOVED***
	return pushStack(s, mapNodes(s.Nodes, func(i int, n *html.Node) []*html.Node ***REMOVED***
		// For each node in the selection, test the node itself, then each parent
		// until a match is found.
		for ; n != nil; n = n.Parent ***REMOVED***
			if set[n] ***REMOVED***
				return []*html.Node***REMOVED***n***REMOVED***
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***))
***REMOVED***

// ClosestSelection gets the first element that matches one of the nodes in the
// Selection by testing the element itself and traversing up through its ancestors
// in the DOM tree.
func (s *Selection) ClosestSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return pushStack(s, nil)
	***REMOVED***
	return s.ClosestNodes(sel.Nodes...)
***REMOVED***

// Parents gets the ancestors of each element in the current Selection. It
// returns a new Selection object with the matched elements.
func (s *Selection) Parents() *Selection ***REMOVED***
	return pushStack(s, getParentsNodes(s.Nodes, nil, nil))
***REMOVED***

// ParentsFiltered gets the ancestors of each element in the current
// Selection. It returns a new Selection object with the matched elements.
func (s *Selection) ParentsFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nil), compileMatcher(selector))
***REMOVED***

// ParentsMatcher gets the ancestors of each element in the current
// Selection. It returns a new Selection object with the matched elements.
func (s *Selection) ParentsMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nil), m)
***REMOVED***

// ParentsUntil gets the ancestors of each element in the Selection, up to but
// not including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (s *Selection) ParentsUntil(selector string) *Selection ***REMOVED***
	return pushStack(s, getParentsNodes(s.Nodes, compileMatcher(selector), nil))
***REMOVED***

// ParentsUntilMatcher gets the ancestors of each element in the Selection, up to but
// not including the element matched by the matcher. It returns a new Selection
// object containing the matched elements.
func (s *Selection) ParentsUntilMatcher(m Matcher) *Selection ***REMOVED***
	return pushStack(s, getParentsNodes(s.Nodes, m, nil))
***REMOVED***

// ParentsUntilSelection gets the ancestors of each element in the Selection,
// up to but not including the elements in the specified Selection. It returns a
// new Selection object containing the matched elements.
func (s *Selection) ParentsUntilSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.Parents()
	***REMOVED***
	return s.ParentsUntilNodes(sel.Nodes...)
***REMOVED***

// ParentsUntilNodes gets the ancestors of each element in the Selection,
// up to but not including the specified nodes. It returns a
// new Selection object containing the matched elements.
func (s *Selection) ParentsUntilNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return pushStack(s, getParentsNodes(s.Nodes, nil, nodes))
***REMOVED***

// ParentsFilteredUntil is like ParentsUntil, with the option to filter the
// results based on a selector string. It returns a new Selection
// object containing the matched elements.
func (s *Selection) ParentsFilteredUntil(filterSelector, untilSelector string) *Selection ***REMOVED***
	return filterAndPush(s, getParentsNodes(s.Nodes, compileMatcher(untilSelector), nil), compileMatcher(filterSelector))
***REMOVED***

// ParentsFilteredUntilMatcher is like ParentsUntilMatcher, with the option to filter the
// results based on a matcher. It returns a new Selection object containing the matched elements.
func (s *Selection) ParentsFilteredUntilMatcher(filter, until Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getParentsNodes(s.Nodes, until, nil), filter)
***REMOVED***

// ParentsFilteredUntilSelection is like ParentsUntilSelection, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsFilteredUntilSelection(filterSelector string, sel *Selection) *Selection ***REMOVED***
	return s.ParentsMatcherUntilSelection(compileMatcher(filterSelector), sel)
***REMOVED***

// ParentsMatcherUntilSelection is like ParentsUntilSelection, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsMatcherUntilSelection(filter Matcher, sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.ParentsMatcher(filter)
	***REMOVED***
	return s.ParentsMatcherUntilNodes(filter, sel.Nodes...)
***REMOVED***

// ParentsFilteredUntilNodes is like ParentsUntilNodes, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection ***REMOVED***
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nodes), compileMatcher(filterSelector))
***REMOVED***

// ParentsMatcherUntilNodes is like ParentsUntilNodes, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsMatcherUntilNodes(filter Matcher, nodes ...*html.Node) *Selection ***REMOVED***
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nodes), filter)
***REMOVED***

// Siblings gets the siblings of each element in the Selection. It returns
// a new Selection object containing the matched elements.
func (s *Selection) Siblings() *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingAll, nil, nil))
***REMOVED***

// SiblingsFiltered gets the siblings of each element in the Selection
// filtered by a selector. It returns a new Selection object containing the
// matched elements.
func (s *Selection) SiblingsFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingAll, nil, nil), compileMatcher(selector))
***REMOVED***

// SiblingsMatcher gets the siblings of each element in the Selection
// filtered by a matcher. It returns a new Selection object containing the
// matched elements.
func (s *Selection) SiblingsMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingAll, nil, nil), m)
***REMOVED***

// Next gets the immediately following sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) Next() *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNext, nil, nil))
***REMOVED***

// NextFiltered gets the immediately following sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNext, nil, nil), compileMatcher(selector))
***REMOVED***

// NextMatcher gets the immediately following sibling of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNext, nil, nil), m)
***REMOVED***

// NextAll gets all the following siblings of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) NextAll() *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextAll, nil, nil))
***REMOVED***

// NextAllFiltered gets all the following siblings of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextAllFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextAll, nil, nil), compileMatcher(selector))
***REMOVED***

// NextAllMatcher gets all the following siblings of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextAllMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextAll, nil, nil), m)
***REMOVED***

// Prev gets the immediately preceding sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) Prev() *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrev, nil, nil))
***REMOVED***

// PrevFiltered gets the immediately preceding sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrev, nil, nil), compileMatcher(selector))
***REMOVED***

// PrevMatcher gets the immediately preceding sibling of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrev, nil, nil), m)
***REMOVED***

// PrevAll gets all the preceding siblings of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) PrevAll() *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevAll, nil, nil))
***REMOVED***

// PrevAllFiltered gets all the preceding siblings of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevAllFiltered(selector string) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevAll, nil, nil), compileMatcher(selector))
***REMOVED***

// PrevAllMatcher gets all the preceding siblings of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevAllMatcher(m Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevAll, nil, nil), m)
***REMOVED***

// NextUntil gets all following siblings of each element up to but not
// including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntil(selector string) *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		compileMatcher(selector), nil))
***REMOVED***

// NextUntilMatcher gets all following siblings of each element up to but not
// including the element matched by the matcher. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntilMatcher(m Matcher) *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		m, nil))
***REMOVED***

// NextUntilSelection gets all following siblings of each element up to but not
// including the element matched by the Selection. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntilSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.NextAll()
	***REMOVED***
	return s.NextUntilNodes(sel.Nodes...)
***REMOVED***

// NextUntilNodes gets all following siblings of each element up to but not
// including the element matched by the nodes. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntilNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		nil, nodes))
***REMOVED***

// PrevUntil gets all preceding siblings of each element up to but not
// including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntil(selector string) *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		compileMatcher(selector), nil))
***REMOVED***

// PrevUntilMatcher gets all preceding siblings of each element up to but not
// including the element matched by the matcher. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntilMatcher(m Matcher) *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		m, nil))
***REMOVED***

// PrevUntilSelection gets all preceding siblings of each element up to but not
// including the element matched by the Selection. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntilSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.PrevAll()
	***REMOVED***
	return s.PrevUntilNodes(sel.Nodes...)
***REMOVED***

// PrevUntilNodes gets all preceding siblings of each element up to but not
// including the element matched by the nodes. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntilNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		nil, nodes))
***REMOVED***

// NextFilteredUntil is like NextUntil, with the option to filter
// the results based on a selector string.
// It returns a new Selection object containing the matched elements.
func (s *Selection) NextFilteredUntil(filterSelector, untilSelector string) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		compileMatcher(untilSelector), nil), compileMatcher(filterSelector))
***REMOVED***

// NextFilteredUntilMatcher is like NextUntilMatcher, with the option to filter
// the results based on a matcher.
// It returns a new Selection object containing the matched elements.
func (s *Selection) NextFilteredUntilMatcher(filter, until Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		until, nil), filter)
***REMOVED***

// NextFilteredUntilSelection is like NextUntilSelection, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextFilteredUntilSelection(filterSelector string, sel *Selection) *Selection ***REMOVED***
	return s.NextMatcherUntilSelection(compileMatcher(filterSelector), sel)
***REMOVED***

// NextMatcherUntilSelection is like NextUntilSelection, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextMatcherUntilSelection(filter Matcher, sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.NextMatcher(filter)
	***REMOVED***
	return s.NextMatcherUntilNodes(filter, sel.Nodes...)
***REMOVED***

// NextFilteredUntilNodes is like NextUntilNodes, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		nil, nodes), compileMatcher(filterSelector))
***REMOVED***

// NextMatcherUntilNodes is like NextUntilNodes, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextMatcherUntilNodes(filter Matcher, nodes ...*html.Node) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		nil, nodes), filter)
***REMOVED***

// PrevFilteredUntil is like PrevUntil, with the option to filter
// the results based on a selector string.
// It returns a new Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntil(filterSelector, untilSelector string) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		compileMatcher(untilSelector), nil), compileMatcher(filterSelector))
***REMOVED***

// PrevFilteredUntilMatcher is like PrevUntilMatcher, with the option to filter
// the results based on a matcher.
// It returns a new Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntilMatcher(filter, until Matcher) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		until, nil), filter)
***REMOVED***

// PrevFilteredUntilSelection is like PrevUntilSelection, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntilSelection(filterSelector string, sel *Selection) *Selection ***REMOVED***
	return s.PrevMatcherUntilSelection(compileMatcher(filterSelector), sel)
***REMOVED***

// PrevMatcherUntilSelection is like PrevUntilSelection, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevMatcherUntilSelection(filter Matcher, sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.PrevMatcher(filter)
	***REMOVED***
	return s.PrevMatcherUntilNodes(filter, sel.Nodes...)
***REMOVED***

// PrevFilteredUntilNodes is like PrevUntilNodes, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		nil, nodes), compileMatcher(filterSelector))
***REMOVED***

// PrevMatcherUntilNodes is like PrevUntilNodes, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevMatcherUntilNodes(filter Matcher, nodes ...*html.Node) *Selection ***REMOVED***
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		nil, nodes), filter)
***REMOVED***

// Filter and push filters the nodes based on a matcher, and pushes the results
// on the stack, with the srcSel as previous selection.
func filterAndPush(srcSel *Selection, nodes []*html.Node, m Matcher) *Selection ***REMOVED***
	// Create a temporary Selection with the specified nodes to filter using winnow
	sel := &Selection***REMOVED***nodes, srcSel.document, nil***REMOVED***
	// Filter based on matcher and push on stack
	return pushStack(srcSel, winnow(sel, m, true))
***REMOVED***

// Internal implementation of Find that return raw nodes.
func findWithMatcher(nodes []*html.Node, m Matcher) []*html.Node ***REMOVED***
	// Map nodes to find the matches within the children of each node
	return mapNodes(nodes, func(i int, n *html.Node) (result []*html.Node) ***REMOVED***
		// Go down one level, becausejQuery's Find selects only within descendants
		for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			if c.Type == html.ElementNode ***REMOVED***
				result = append(result, m.MatchAll(c)...)
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***)
***REMOVED***

// Internal implementation to get all parent nodes, stopping at the specified
// node (or nil if no stop).
func getParentsNodes(nodes []*html.Node, stopm Matcher, stopNodes []*html.Node) []*html.Node ***REMOVED***
	return mapNodes(nodes, func(i int, n *html.Node) (result []*html.Node) ***REMOVED***
		for p := n.Parent; p != nil; p = p.Parent ***REMOVED***
			sel := newSingleSelection(p, nil)
			if stopm != nil ***REMOVED***
				if sel.IsMatcher(stopm) ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if len(stopNodes) > 0 ***REMOVED***
				if sel.IsNodes(stopNodes...) ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			if p.Type == html.ElementNode ***REMOVED***
				result = append(result, p)
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***)
***REMOVED***

// Internal implementation of sibling nodes that return a raw slice of matches.
func getSiblingNodes(nodes []*html.Node, st siblingType, untilm Matcher, untilNodes []*html.Node) []*html.Node ***REMOVED***
	var f func(*html.Node) bool

	// If the requested siblings are ...Until, create the test function to
	// determine if the until condition is reached (returns true if it is)
	if st == siblingNextUntil || st == siblingPrevUntil ***REMOVED***
		f = func(n *html.Node) bool ***REMOVED***
			if untilm != nil ***REMOVED***
				// Matcher-based condition
				sel := newSingleSelection(n, nil)
				return sel.IsMatcher(untilm)
			***REMOVED*** else if len(untilNodes) > 0 ***REMOVED***
				// Nodes-based condition
				sel := newSingleSelection(n, nil)
				return sel.IsNodes(untilNodes...)
			***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return mapNodes(nodes, func(i int, n *html.Node) []*html.Node ***REMOVED***
		return getChildrenWithSiblingType(n.Parent, st, n, f)
	***REMOVED***)
***REMOVED***

// Gets the children nodes of each node in the specified slice of nodes,
// based on the sibling type request.
func getChildrenNodes(nodes []*html.Node, st siblingType) []*html.Node ***REMOVED***
	return mapNodes(nodes, func(i int, n *html.Node) []*html.Node ***REMOVED***
		return getChildrenWithSiblingType(n, st, nil, nil)
	***REMOVED***)
***REMOVED***

// Gets the children of the specified parent, based on the requested sibling
// type, skipping a specified node if required.
func getChildrenWithSiblingType(parent *html.Node, st siblingType, skipNode *html.Node,
	untilFunc func(*html.Node) bool) (result []*html.Node) ***REMOVED***

	// Create the iterator function
	var iter = func(cur *html.Node) (ret *html.Node) ***REMOVED***
		// Based on the sibling type requested, iterate the right way
		for ***REMOVED***
			switch st ***REMOVED***
			case siblingAll, siblingAllIncludingNonElements:
				if cur == nil ***REMOVED***
					// First iteration, start with first child of parent
					// Skip node if required
					if ret = parent.FirstChild; ret == skipNode && skipNode != nil ***REMOVED***
						ret = skipNode.NextSibling
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					// Skip node if required
					if ret = cur.NextSibling; ret == skipNode && skipNode != nil ***REMOVED***
						ret = skipNode.NextSibling
					***REMOVED***
				***REMOVED***
			case siblingPrev, siblingPrevAll, siblingPrevUntil:
				if cur == nil ***REMOVED***
					// Start with previous sibling of the skip node
					ret = skipNode.PrevSibling
				***REMOVED*** else ***REMOVED***
					ret = cur.PrevSibling
				***REMOVED***
			case siblingNext, siblingNextAll, siblingNextUntil:
				if cur == nil ***REMOVED***
					// Start with next sibling of the skip node
					ret = skipNode.NextSibling
				***REMOVED*** else ***REMOVED***
					ret = cur.NextSibling
				***REMOVED***
			default:
				panic("Invalid sibling type.")
			***REMOVED***
			if ret == nil || ret.Type == html.ElementNode || st == siblingAllIncludingNonElements ***REMOVED***
				return
			***REMOVED***
			// Not a valid node, try again from this one
			cur = ret
		***REMOVED***
	***REMOVED***

	for c := iter(nil); c != nil; c = iter(c) ***REMOVED***
		// If this is an ...Until case, test before append (returns true
		// if the until condition is reached)
		if st == siblingNextUntil || st == siblingPrevUntil ***REMOVED***
			if untilFunc(c) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		result = append(result, c)
		if st == siblingNext || st == siblingPrev ***REMOVED***
			// Only one node was requested (immediate next or previous), so exit
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// Internal implementation of parent nodes that return a raw slice of Nodes.
func getParentNodes(nodes []*html.Node) []*html.Node ***REMOVED***
	return mapNodes(nodes, func(i int, n *html.Node) []*html.Node ***REMOVED***
		if n.Parent != nil && n.Parent.Type == html.ElementNode ***REMOVED***
			return []*html.Node***REMOVED***n.Parent***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***)
***REMOVED***

// Internal map function used by many traversing methods. Takes the source nodes
// to iterate on and the mapping function that returns an array of nodes.
// Returns an array of nodes mapped by calling the callback function once for
// each node in the source nodes.
func mapNodes(nodes []*html.Node, f func(int, *html.Node) []*html.Node) (result []*html.Node) ***REMOVED***
	set := make(map[*html.Node]bool)
	for i, n := range nodes ***REMOVED***
		if vals := f(i, n); len(vals) > 0 ***REMOVED***
			result = appendWithoutDuplicates(result, vals, set)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***
