package goquery

import (
	"bytes"

	"golang.org/x/net/html"
)

// used to determine if a set (map[*html.Node]bool) should be used
// instead of iterating over a slice. The set uses more memory and
// is slower than slice iteration for small N.
const minNodesForSet = 1000

var nodeNames = []string***REMOVED***
	html.ErrorNode:    "#error",
	html.TextNode:     "#text",
	html.DocumentNode: "#document",
	html.CommentNode:  "#comment",
***REMOVED***

// NodeName returns the node name of the first element in the selection.
// It tries to behave in a similar way as the DOM's nodeName property
// (https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeName).
//
// Go's net/html package defines the following node types, listed with
// the corresponding returned value from this function:
//
//     ErrorNode : #error
//     TextNode : #text
//     DocumentNode : #document
//     ElementNode : the element's tag name
//     CommentNode : #comment
//     DoctypeNode : the name of the document type
//
func NodeName(s *Selection) string ***REMOVED***
	if s.Length() == 0 ***REMOVED***
		return ""
	***REMOVED***
	return nodeName(s.Get(0))
***REMOVED***

// nodeName returns the node name of the given html node.
// See NodeName for additional details on behaviour.
func nodeName(node *html.Node) string ***REMOVED***
	if node == nil ***REMOVED***
		return ""
	***REMOVED***

	switch node.Type ***REMOVED***
	case html.ElementNode, html.DoctypeNode:
		return node.Data
	default:
		if node.Type >= 0 && int(node.Type) < len(nodeNames) ***REMOVED***
			return nodeNames[node.Type]
		***REMOVED***
		return ""
	***REMOVED***
***REMOVED***

// OuterHtml returns the outer HTML rendering of the first item in
// the selection - that is, the HTML including the first element's
// tag and attributes.
//
// Unlike InnerHtml, this is a function and not a method on the Selection,
// because this is not a jQuery method (in javascript-land, this is
// a property provided by the DOM).
func OuterHtml(s *Selection) (string, error) ***REMOVED***
	var buf bytes.Buffer

	if s.Length() == 0 ***REMOVED***
		return "", nil
	***REMOVED***
	n := s.Get(0)
	if err := html.Render(&buf, n); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return buf.String(), nil
***REMOVED***

// Loop through all container nodes to search for the target node.
func sliceContains(container []*html.Node, contained *html.Node) bool ***REMOVED***
	for _, n := range container ***REMOVED***
		if nodeContains(n, contained) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// Checks if the contained node is within the container node.
func nodeContains(container *html.Node, contained *html.Node) bool ***REMOVED***
	// Check if the parent of the contained node is the container node, traversing
	// upward until the top is reached, or the container is found.
	for contained = contained.Parent; contained != nil; contained = contained.Parent ***REMOVED***
		if container == contained ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Checks if the target node is in the slice of nodes.
func isInSlice(slice []*html.Node, node *html.Node) bool ***REMOVED***
	return indexInSlice(slice, node) > -1
***REMOVED***

// Returns the index of the target node in the slice, or -1.
func indexInSlice(slice []*html.Node, node *html.Node) int ***REMOVED***
	if node != nil ***REMOVED***
		for i, n := range slice ***REMOVED***
			if n == node ***REMOVED***
				return i
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// Appends the new nodes to the target slice, making sure no duplicate is added.
// There is no check to the original state of the target slice, so it may still
// contain duplicates. The target slice is returned because append() may create
// a new underlying array. If targetSet is nil, a local set is created with the
// target if len(target) + len(nodes) is greater than minNodesForSet.
func appendWithoutDuplicates(target []*html.Node, nodes []*html.Node, targetSet map[*html.Node]bool) []*html.Node ***REMOVED***
	// if there are not that many nodes, don't use the map, faster to just use nested loops
	// (unless a non-nil targetSet is passed, in which case the caller knows better).
	if targetSet == nil && len(target)+len(nodes) < minNodesForSet ***REMOVED***
		for _, n := range nodes ***REMOVED***
			if !isInSlice(target, n) ***REMOVED***
				target = append(target, n)
			***REMOVED***
		***REMOVED***
		return target
	***REMOVED***

	// if a targetSet is passed, then assume it is reliable, otherwise create one
	// and initialize it with the current target contents.
	if targetSet == nil ***REMOVED***
		targetSet = make(map[*html.Node]bool, len(target))
		for _, n := range target ***REMOVED***
			targetSet[n] = true
		***REMOVED***
	***REMOVED***
	for _, n := range nodes ***REMOVED***
		if !targetSet[n] ***REMOVED***
			target = append(target, n)
			targetSet[n] = true
		***REMOVED***
	***REMOVED***

	return target
***REMOVED***

// Loop through a selection, returning only those nodes that pass the predicate
// function.
func grep(sel *Selection, predicate func(i int, s *Selection) bool) (result []*html.Node) ***REMOVED***
	for i, n := range sel.Nodes ***REMOVED***
		if predicate(i, newSingleSelection(n, sel.document)) ***REMOVED***
			result = append(result, n)
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// Creates a new Selection object based on the specified nodes, and keeps the
// source Selection object on the stack (linked list).
func pushStack(fromSel *Selection, nodes []*html.Node) *Selection ***REMOVED***
	result := &Selection***REMOVED***nodes, fromSel.document, fromSel***REMOVED***
	return result
***REMOVED***
