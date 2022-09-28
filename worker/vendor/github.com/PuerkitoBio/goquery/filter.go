package goquery

import "golang.org/x/net/html"

// Filter reduces the set of matched elements to those that match the selector string.
// It returns a new Selection object for this subset of matching elements.
func (s *Selection) Filter(selector string) *Selection ***REMOVED***
	return s.FilterMatcher(compileMatcher(selector))
***REMOVED***

// FilterMatcher reduces the set of matched elements to those that match
// the given matcher. It returns a new Selection object for this subset
// of matching elements.
func (s *Selection) FilterMatcher(m Matcher) *Selection ***REMOVED***
	return pushStack(s, winnow(s, m, true))
***REMOVED***

// Not removes elements from the Selection that match the selector string.
// It returns a new Selection object with the matching elements removed.
func (s *Selection) Not(selector string) *Selection ***REMOVED***
	return s.NotMatcher(compileMatcher(selector))
***REMOVED***

// NotMatcher removes elements from the Selection that match the given matcher.
// It returns a new Selection object with the matching elements removed.
func (s *Selection) NotMatcher(m Matcher) *Selection ***REMOVED***
	return pushStack(s, winnow(s, m, false))
***REMOVED***

// FilterFunction reduces the set of matched elements to those that pass the function's test.
// It returns a new Selection object for this subset of elements.
func (s *Selection) FilterFunction(f func(int, *Selection) bool) *Selection ***REMOVED***
	return pushStack(s, winnowFunction(s, f, true))
***REMOVED***

// NotFunction removes elements from the Selection that pass the function's test.
// It returns a new Selection object with the matching elements removed.
func (s *Selection) NotFunction(f func(int, *Selection) bool) *Selection ***REMOVED***
	return pushStack(s, winnowFunction(s, f, false))
***REMOVED***

// FilterNodes reduces the set of matched elements to those that match the specified nodes.
// It returns a new Selection object for this subset of elements.
func (s *Selection) FilterNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return pushStack(s, winnowNodes(s, nodes, true))
***REMOVED***

// NotNodes removes elements from the Selection that match the specified nodes.
// It returns a new Selection object with the matching elements removed.
func (s *Selection) NotNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return pushStack(s, winnowNodes(s, nodes, false))
***REMOVED***

// FilterSelection reduces the set of matched elements to those that match a
// node in the specified Selection object.
// It returns a new Selection object for this subset of elements.
func (s *Selection) FilterSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return pushStack(s, winnowNodes(s, nil, true))
	***REMOVED***
	return pushStack(s, winnowNodes(s, sel.Nodes, true))
***REMOVED***

// NotSelection removes elements from the Selection that match a node in the specified
// Selection object. It returns a new Selection object with the matching elements removed.
func (s *Selection) NotSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return pushStack(s, winnowNodes(s, nil, false))
	***REMOVED***
	return pushStack(s, winnowNodes(s, sel.Nodes, false))
***REMOVED***

// Intersection is an alias for FilterSelection.
func (s *Selection) Intersection(sel *Selection) *Selection ***REMOVED***
	return s.FilterSelection(sel)
***REMOVED***

// Has reduces the set of matched elements to those that have a descendant
// that matches the selector.
// It returns a new Selection object with the matching elements.
func (s *Selection) Has(selector string) *Selection ***REMOVED***
	return s.HasSelection(s.document.Find(selector))
***REMOVED***

// HasMatcher reduces the set of matched elements to those that have a descendant
// that matches the matcher.
// It returns a new Selection object with the matching elements.
func (s *Selection) HasMatcher(m Matcher) *Selection ***REMOVED***
	return s.HasSelection(s.document.FindMatcher(m))
***REMOVED***

// HasNodes reduces the set of matched elements to those that have a
// descendant that matches one of the nodes.
// It returns a new Selection object with the matching elements.
func (s *Selection) HasNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return s.FilterFunction(func(_ int, sel *Selection) bool ***REMOVED***
		// Add all nodes that contain one of the specified nodes
		for _, n := range nodes ***REMOVED***
			if sel.Contains(n) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***)
***REMOVED***

// HasSelection reduces the set of matched elements to those that have a
// descendant that matches one of the nodes of the specified Selection object.
// It returns a new Selection object with the matching elements.
func (s *Selection) HasSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.HasNodes()
	***REMOVED***
	return s.HasNodes(sel.Nodes...)
***REMOVED***

// End ends the most recent filtering operation in the current chain and
// returns the set of matched elements to its previous state.
func (s *Selection) End() *Selection ***REMOVED***
	if s.prevSel != nil ***REMOVED***
		return s.prevSel
	***REMOVED***
	return newEmptySelection(s.document)
***REMOVED***

// Filter based on the matcher, and the indicator to keep (Filter) or
// to get rid of (Not) the matching elements.
func winnow(sel *Selection, m Matcher, keep bool) []*html.Node ***REMOVED***
	// Optimize if keep is requested
	if keep ***REMOVED***
		return m.Filter(sel.Nodes)
	***REMOVED***
	// Use grep
	return grep(sel, func(i int, s *Selection) bool ***REMOVED***
		return !m.Match(s.Get(0))
	***REMOVED***)
***REMOVED***

// Filter based on an array of nodes, and the indicator to keep (Filter) or
// to get rid of (Not) the matching elements.
func winnowNodes(sel *Selection, nodes []*html.Node, keep bool) []*html.Node ***REMOVED***
	if len(nodes)+len(sel.Nodes) < minNodesForSet ***REMOVED***
		return grep(sel, func(i int, s *Selection) bool ***REMOVED***
			return isInSlice(nodes, s.Get(0)) == keep
		***REMOVED***)
	***REMOVED***

	set := make(map[*html.Node]bool)
	for _, n := range nodes ***REMOVED***
		set[n] = true
	***REMOVED***
	return grep(sel, func(i int, s *Selection) bool ***REMOVED***
		return set[s.Get(0)] == keep
	***REMOVED***)
***REMOVED***

// Filter based on a function test, and the indicator to keep (Filter) or
// to get rid of (Not) the matching elements.
func winnowFunction(sel *Selection, f func(int, *Selection) bool, keep bool) []*html.Node ***REMOVED***
	return grep(sel, func(i int, s *Selection) bool ***REMOVED***
		return f(i, s) == keep
	***REMOVED***)
***REMOVED***
