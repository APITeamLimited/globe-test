package goquery

import "golang.org/x/net/html"

// Add adds the selector string's matching nodes to those in the current
// selection and returns a new Selection object.
// The selector string is run in the context of the document of the current
// Selection object.
func (s *Selection) Add(selector string) *Selection ***REMOVED***
	return s.AddNodes(findWithMatcher([]*html.Node***REMOVED***s.document.rootNode***REMOVED***, compileMatcher(selector))...)
***REMOVED***

// AddMatcher adds the matcher's matching nodes to those in the current
// selection and returns a new Selection object.
// The matcher is run in the context of the document of the current
// Selection object.
func (s *Selection) AddMatcher(m Matcher) *Selection ***REMOVED***
	return s.AddNodes(findWithMatcher([]*html.Node***REMOVED***s.document.rootNode***REMOVED***, m)...)
***REMOVED***

// AddSelection adds the specified Selection object's nodes to those in the
// current selection and returns a new Selection object.
func (s *Selection) AddSelection(sel *Selection) *Selection ***REMOVED***
	if sel == nil ***REMOVED***
		return s.AddNodes()
	***REMOVED***
	return s.AddNodes(sel.Nodes...)
***REMOVED***

// Union is an alias for AddSelection.
func (s *Selection) Union(sel *Selection) *Selection ***REMOVED***
	return s.AddSelection(sel)
***REMOVED***

// AddNodes adds the specified nodes to those in the
// current selection and returns a new Selection object.
func (s *Selection) AddNodes(nodes ...*html.Node) *Selection ***REMOVED***
	return pushStack(s, appendWithoutDuplicates(s.Nodes, nodes, nil))
***REMOVED***

// AndSelf adds the previous set of elements on the stack to the current set.
// It returns a new Selection object containing the current Selection combined
// with the previous one.
func (s *Selection) AndSelf() *Selection ***REMOVED***
	return s.AddSelection(s.prevSel)
***REMOVED***
