// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"golang.org/x/net/html/atom"
)

// A NodeType is the type of a Node.
type NodeType uint32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
	// RawNode nodes are not returned by the parser, but can be part of the
	// Node tree passed to func Render to insert raw HTML (without escaping).
	// If so, this package makes no guarantee that the rendered HTML is secure
	// (from e.g. Cross Site Scripting attacks) or well-formed.
	RawNode
	scopeMarkerNode
)

// Section 12.2.4.3 says "The markers are inserted when entering applet,
// object, marquee, template, td, th, and caption elements, and are used
// to prevent formatting from "leaking" into applet, object, marquee,
// template, td, th, and caption elements".
var scopeMarker = Node***REMOVED***Type: scopeMarkerNode***REMOVED***

// A Node consists of a NodeType and some Data (tag name for element nodes,
// content for text) and are part of a tree of Nodes. Element nodes may also
// have a Namespace and contain a slice of Attributes. Data is unescaped, so
// that it looks like "a<b" rather than "a&lt;b". For element nodes, DataAtom
// is the atom for Data, or zero if Data is not a known tag name.
//
// An empty Namespace implies a "http://www.w3.org/1999/xhtml" namespace.
// Similarly, "math" is short for "http://www.w3.org/1998/Math/MathML", and
// "svg" is short for "http://www.w3.org/2000/svg".
type Node struct ***REMOVED***
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	Type      NodeType
	DataAtom  atom.Atom
	Data      string
	Namespace string
	Attr      []Attribute
***REMOVED***

// InsertBefore inserts newChild as a child of n, immediately before oldChild
// in the sequence of n's children. oldChild may be nil, in which case newChild
// is appended to the end of n's children.
//
// It will panic if newChild already has a parent or siblings.
func (n *Node) InsertBefore(newChild, oldChild *Node) ***REMOVED***
	if newChild.Parent != nil || newChild.PrevSibling != nil || newChild.NextSibling != nil ***REMOVED***
		panic("html: InsertBefore called for an attached child Node")
	***REMOVED***
	var prev, next *Node
	if oldChild != nil ***REMOVED***
		prev, next = oldChild.PrevSibling, oldChild
	***REMOVED*** else ***REMOVED***
		prev = n.LastChild
	***REMOVED***
	if prev != nil ***REMOVED***
		prev.NextSibling = newChild
	***REMOVED*** else ***REMOVED***
		n.FirstChild = newChild
	***REMOVED***
	if next != nil ***REMOVED***
		next.PrevSibling = newChild
	***REMOVED*** else ***REMOVED***
		n.LastChild = newChild
	***REMOVED***
	newChild.Parent = n
	newChild.PrevSibling = prev
	newChild.NextSibling = next
***REMOVED***

// AppendChild adds a node c as a child of n.
//
// It will panic if c already has a parent or siblings.
func (n *Node) AppendChild(c *Node) ***REMOVED***
	if c.Parent != nil || c.PrevSibling != nil || c.NextSibling != nil ***REMOVED***
		panic("html: AppendChild called for an attached child Node")
	***REMOVED***
	last := n.LastChild
	if last != nil ***REMOVED***
		last.NextSibling = c
	***REMOVED*** else ***REMOVED***
		n.FirstChild = c
	***REMOVED***
	n.LastChild = c
	c.Parent = n
	c.PrevSibling = last
***REMOVED***

// RemoveChild removes a node c that is a child of n. Afterwards, c will have
// no parent and no siblings.
//
// It will panic if c's parent is not n.
func (n *Node) RemoveChild(c *Node) ***REMOVED***
	if c.Parent != n ***REMOVED***
		panic("html: RemoveChild called for a non-child Node")
	***REMOVED***
	if n.FirstChild == c ***REMOVED***
		n.FirstChild = c.NextSibling
	***REMOVED***
	if c.NextSibling != nil ***REMOVED***
		c.NextSibling.PrevSibling = c.PrevSibling
	***REMOVED***
	if n.LastChild == c ***REMOVED***
		n.LastChild = c.PrevSibling
	***REMOVED***
	if c.PrevSibling != nil ***REMOVED***
		c.PrevSibling.NextSibling = c.NextSibling
	***REMOVED***
	c.Parent = nil
	c.PrevSibling = nil
	c.NextSibling = nil
***REMOVED***

// reparentChildren reparents all of src's child nodes to dst.
func reparentChildren(dst, src *Node) ***REMOVED***
	for ***REMOVED***
		child := src.FirstChild
		if child == nil ***REMOVED***
			break
		***REMOVED***
		src.RemoveChild(child)
		dst.AppendChild(child)
	***REMOVED***
***REMOVED***

// clone returns a new node with the same type, data and attributes.
// The clone has no parent, no siblings and no children.
func (n *Node) clone() *Node ***REMOVED***
	m := &Node***REMOVED***
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]Attribute, len(n.Attr)),
	***REMOVED***
	copy(m.Attr, n.Attr)
	return m
***REMOVED***

// nodeStack is a stack of nodes.
type nodeStack []*Node

// pop pops the stack. It will panic if s is empty.
func (s *nodeStack) pop() *Node ***REMOVED***
	i := len(*s)
	n := (*s)[i-1]
	*s = (*s)[:i-1]
	return n
***REMOVED***

// top returns the most recently pushed node, or nil if s is empty.
func (s *nodeStack) top() *Node ***REMOVED***
	if i := len(*s); i > 0 ***REMOVED***
		return (*s)[i-1]
	***REMOVED***
	return nil
***REMOVED***

// index returns the index of the top-most occurrence of n in the stack, or -1
// if n is not present.
func (s *nodeStack) index(n *Node) int ***REMOVED***
	for i := len(*s) - 1; i >= 0; i-- ***REMOVED***
		if (*s)[i] == n ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// contains returns whether a is within s.
func (s *nodeStack) contains(a atom.Atom) bool ***REMOVED***
	for _, n := range *s ***REMOVED***
		if n.DataAtom == a && n.Namespace == "" ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// insert inserts a node at the given index.
func (s *nodeStack) insert(i int, n *Node) ***REMOVED***
	(*s) = append(*s, nil)
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = n
***REMOVED***

// remove removes a node from the stack. It is a no-op if n is not present.
func (s *nodeStack) remove(n *Node) ***REMOVED***
	i := s.index(n)
	if i == -1 ***REMOVED***
		return
	***REMOVED***
	copy((*s)[i:], (*s)[i+1:])
	j := len(*s) - 1
	(*s)[j] = nil
	*s = (*s)[:j]
***REMOVED***

type insertionModeStack []insertionMode

func (s *insertionModeStack) pop() (im insertionMode) ***REMOVED***
	i := len(*s)
	im = (*s)[i-1]
	*s = (*s)[:i-1]
	return im
***REMOVED***

func (s *insertionModeStack) top() insertionMode ***REMOVED***
	if i := len(*s); i > 0 ***REMOVED***
		return (*s)[i-1]
	***REMOVED***
	return nil
***REMOVED***
