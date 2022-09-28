package ast

import (
	"fmt"
	"strings"
)

// Identifier is a possibly-qualified name. This is used to distinguish
// ValueNode values that are references/identifiers vs. those that are
// string literals.
type Identifier string

// IdentValueNode is an AST node that represents an identifier.
type IdentValueNode interface ***REMOVED***
	ValueNode
	AsIdentifier() Identifier
***REMOVED***

var _ IdentValueNode = (*IdentNode)(nil)
var _ IdentValueNode = (*CompoundIdentNode)(nil)

// IdentNode represents a simple, unqualified identifier. These are used to name
// elements declared in a protobuf file or to refer to elements. Example:
//
//  foobar
type IdentNode struct ***REMOVED***
	terminalNode
	Val string
***REMOVED***

// NewIdentNode creates a new *IdentNode. The given val is the identifier text.
func NewIdentNode(val string, info TokenInfo) *IdentNode ***REMOVED***
	return &IdentNode***REMOVED***
		terminalNode: info.asTerminalNode(),
		Val:          val,
	***REMOVED***
***REMOVED***

func (n *IdentNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.AsIdentifier()
***REMOVED***

func (n *IdentNode) AsIdentifier() Identifier ***REMOVED***
	return Identifier(n.Val)
***REMOVED***

// ToKeyword is used to convert identifiers to keywords. Since keywords are not
// reserved in the protobuf language, they are initially lexed as identifiers
// and then converted to keywords based on context.
func (n *IdentNode) ToKeyword() *KeywordNode ***REMOVED***
	return (*KeywordNode)(n)
***REMOVED***

// CompoundIdentNode represents a qualified identifier. A qualified identifier
// has at least one dot and possibly multiple identifier names (all separated by
// dots). If the identifier has a leading dot, then it is a *fully* qualified
// identifier. Example:
//
//  .com.foobar.Baz
type CompoundIdentNode struct ***REMOVED***
	compositeNode
	// Optional leading dot, indicating that the identifier is fully qualified.
	LeadingDot *RuneNode
	Components []*IdentNode
	// Dots[0] is the dot after Components[0]. The length of Dots is always
	// one less than the length of Components.
	Dots []*RuneNode
	// The text value of the identifier, with all components and dots
	// concatenated.
	Val string
***REMOVED***

// NewCompoundIdentNode creates a *CompoundIdentNode. The leadingDot may be nil.
// The dots arg must have a length that is one less than the length of
// components. The components arg must not be empty.
func NewCompoundIdentNode(leadingDot *RuneNode, components []*IdentNode, dots []*RuneNode) *CompoundIdentNode ***REMOVED***
	if len(components) == 0 ***REMOVED***
		panic("must have at least one component")
	***REMOVED***
	if len(dots) != len(components)-1 ***REMOVED***
		panic(fmt.Sprintf("%d components requires %d dots, not %d", len(components), len(components)-1, len(dots)))
	***REMOVED***
	numChildren := len(components)*2 - 1
	if leadingDot != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	var b strings.Builder
	if leadingDot != nil ***REMOVED***
		children = append(children, leadingDot)
		b.WriteRune(leadingDot.Rune)
	***REMOVED***
	for i, comp := range components ***REMOVED***
		if i > 0 ***REMOVED***
			dot := dots[i-1]
			children = append(children, dot)
			b.WriteRune(dot.Rune)
		***REMOVED***
		children = append(children, comp)
		b.WriteString(comp.Val)
	***REMOVED***
	return &CompoundIdentNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		LeadingDot: leadingDot,
		Components: components,
		Dots:       dots,
		Val:        b.String(),
	***REMOVED***
***REMOVED***

func (n *CompoundIdentNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.AsIdentifier()
***REMOVED***

func (n *CompoundIdentNode) AsIdentifier() Identifier ***REMOVED***
	return Identifier(n.Val)
***REMOVED***

// KeywordNode is an AST node that represents a keyword. Keywords are
// like identifiers, but they have special meaning in particular contexts.
// Example:
//
//  message
type KeywordNode IdentNode

// NewKeywordNode creates a new *KeywordNode. The given val is the keyword.
func NewKeywordNode(val string, info TokenInfo) *KeywordNode ***REMOVED***
	return &KeywordNode***REMOVED***
		terminalNode: info.asTerminalNode(),
		Val:          val,
	***REMOVED***
***REMOVED***
