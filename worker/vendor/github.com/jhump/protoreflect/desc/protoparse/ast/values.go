package ast

import (
	"fmt"
	"math"
	"strings"
)

// ValueNode is an AST node that represents a literal value.
//
// It also includes references (e.g. IdentifierValueNode), which can be
// used as values in some contexts, such as describing the default value
// for a field, which can refer to an enum value.
//
// This also allows NoSourceNode to be used in place of a real value node
// for some usages.
type ValueNode interface ***REMOVED***
	Node
	// Value returns a Go representation of the value. For scalars, this
	// will be a string, int64, uint64, float64, or bool. This could also
	// be an Identifier (e.g. IdentValueNodes). It can also be a composite
	// literal:
	//   * For array literals, the type returned will be []ValueNode
	//   * For message literals, the type returned will be []*MessageFieldNode
	Value() interface***REMOVED******REMOVED***
***REMOVED***

var _ ValueNode = (*IdentNode)(nil)
var _ ValueNode = (*CompoundIdentNode)(nil)
var _ ValueNode = (*StringLiteralNode)(nil)
var _ ValueNode = (*CompoundStringLiteralNode)(nil)
var _ ValueNode = (*UintLiteralNode)(nil)
var _ ValueNode = (*PositiveUintLiteralNode)(nil)
var _ ValueNode = (*NegativeIntLiteralNode)(nil)
var _ ValueNode = (*FloatLiteralNode)(nil)
var _ ValueNode = (*SpecialFloatLiteralNode)(nil)
var _ ValueNode = (*SignedFloatLiteralNode)(nil)
var _ ValueNode = (*BoolLiteralNode)(nil)
var _ ValueNode = (*ArrayLiteralNode)(nil)
var _ ValueNode = (*MessageLiteralNode)(nil)
var _ ValueNode = NoSourceNode***REMOVED******REMOVED***

// StringValueNode is an AST node that represents a string literal.
// Such a node can be a single literal (*StringLiteralNode) or a
// concatenation of multiple literals (*CompoundStringLiteralNode).
type StringValueNode interface ***REMOVED***
	ValueNode
	AsString() string
***REMOVED***

var _ StringValueNode = (*StringLiteralNode)(nil)
var _ StringValueNode = (*CompoundStringLiteralNode)(nil)

// StringLiteralNode represents a simple string literal. Example:
//
//  "proto2"
type StringLiteralNode struct ***REMOVED***
	terminalNode
	// Val is the actual string value that the literal indicates.
	Val string
***REMOVED***

// NewStringLiteralNode creates a new *StringLiteralNode with the given val.
func NewStringLiteralNode(val string, info TokenInfo) *StringLiteralNode ***REMOVED***
	return &StringLiteralNode***REMOVED***
		terminalNode: info.asTerminalNode(),
		Val:          val,
	***REMOVED***
***REMOVED***

func (n *StringLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.AsString()
***REMOVED***

func (n *StringLiteralNode) AsString() string ***REMOVED***
	return n.Val
***REMOVED***

// CompoundStringLiteralNode represents a compound string literal, which is
// the concatenaton of adjacent string literals. Example:
//
//  "this "  "is"   " all one "   "string"
type CompoundStringLiteralNode struct ***REMOVED***
	compositeNode
	Val string
***REMOVED***

// NewCompoundLiteralStringNode creates a new *CompoundStringLiteralNode that
// consists of the given string components. The components argument may not be
// empty.
func NewCompoundLiteralStringNode(components ...*StringLiteralNode) *CompoundStringLiteralNode ***REMOVED***
	if len(components) == 0 ***REMOVED***
		panic("must have at least one component")
	***REMOVED***
	children := make([]Node, len(components))
	var b strings.Builder
	for i, comp := range components ***REMOVED***
		children[i] = comp
		b.WriteString(comp.Val)
	***REMOVED***
	return &CompoundStringLiteralNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Val: b.String(),
	***REMOVED***
***REMOVED***

func (n *CompoundStringLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.AsString()
***REMOVED***

func (n *CompoundStringLiteralNode) AsString() string ***REMOVED***
	return n.Val
***REMOVED***

// IntValueNode is an AST node that represents an integer literal. If
// an integer literal is too large for an int64 (or uint64 for
// positive literals), it is represented instead by a FloatValueNode.
type IntValueNode interface ***REMOVED***
	ValueNode
	AsInt64() (int64, bool)
	AsUint64() (uint64, bool)
***REMOVED***

// AsInt32 range checks the given int value and returns its value is
// in the range or 0, false if it is outside the range.
func AsInt32(n IntValueNode, min, max int32) (int32, bool) ***REMOVED***
	i, ok := n.AsInt64()
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	if i < int64(min) || i > int64(max) ***REMOVED***
		return 0, false
	***REMOVED***
	return int32(i), true
***REMOVED***

var _ IntValueNode = (*UintLiteralNode)(nil)
var _ IntValueNode = (*PositiveUintLiteralNode)(nil)
var _ IntValueNode = (*NegativeIntLiteralNode)(nil)

// UintLiteralNode represents a simple integer literal with no sign character.
type UintLiteralNode struct ***REMOVED***
	terminalNode
	// Val is the numeric value indicated by the literal
	Val uint64
***REMOVED***

// NewUintLiteralNode creates a new *UintLiteralNode with the given val.
func NewUintLiteralNode(val uint64, info TokenInfo) *UintLiteralNode ***REMOVED***
	return &UintLiteralNode***REMOVED***
		terminalNode: info.asTerminalNode(),
		Val:          val,
	***REMOVED***
***REMOVED***

func (n *UintLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.Val
***REMOVED***

func (n *UintLiteralNode) AsInt64() (int64, bool) ***REMOVED***
	if n.Val > math.MaxInt64 ***REMOVED***
		return 0, false
	***REMOVED***
	return int64(n.Val), true
***REMOVED***

func (n *UintLiteralNode) AsUint64() (uint64, bool) ***REMOVED***
	return n.Val, true
***REMOVED***

func (n *UintLiteralNode) AsFloat() float64 ***REMOVED***
	return float64(n.Val)
***REMOVED***

// PositiveUintLiteralNode represents an integer literal with a positive (+) sign.
type PositiveUintLiteralNode struct ***REMOVED***
	compositeNode
	Plus *RuneNode
	Uint *UintLiteralNode
	Val  uint64
***REMOVED***

// NewPositiveUintLiteralNode creates a new *PositiveUintLiteralNode. Both
// arguments must be non-nil.
func NewPositiveUintLiteralNode(sign *RuneNode, i *UintLiteralNode) *PositiveUintLiteralNode ***REMOVED***
	if sign == nil ***REMOVED***
		panic("sign is nil")
	***REMOVED***
	if i == nil ***REMOVED***
		panic("i is nil")
	***REMOVED***
	children := []Node***REMOVED***sign, i***REMOVED***
	return &PositiveUintLiteralNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Plus: sign,
		Uint: i,
		Val:  i.Val,
	***REMOVED***
***REMOVED***

func (n *PositiveUintLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.Val
***REMOVED***

func (n *PositiveUintLiteralNode) AsInt64() (int64, bool) ***REMOVED***
	if n.Val > math.MaxInt64 ***REMOVED***
		return 0, false
	***REMOVED***
	return int64(n.Val), true
***REMOVED***

func (n *PositiveUintLiteralNode) AsUint64() (uint64, bool) ***REMOVED***
	return n.Val, true
***REMOVED***

// NegativeIntLiteralNode represents an integer literal with a negative (-) sign.
type NegativeIntLiteralNode struct ***REMOVED***
	compositeNode
	Minus *RuneNode
	Uint  *UintLiteralNode
	Val   int64
***REMOVED***

// NewNegativeIntLiteralNode creates a new *NegativeIntLiteralNode. Both
// arguments must be non-nil.
func NewNegativeIntLiteralNode(sign *RuneNode, i *UintLiteralNode) *NegativeIntLiteralNode ***REMOVED***
	if sign == nil ***REMOVED***
		panic("sign is nil")
	***REMOVED***
	if i == nil ***REMOVED***
		panic("i is nil")
	***REMOVED***
	children := []Node***REMOVED***sign, i***REMOVED***
	return &NegativeIntLiteralNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Minus: sign,
		Uint:  i,
		Val:   -int64(i.Val),
	***REMOVED***
***REMOVED***

func (n *NegativeIntLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.Val
***REMOVED***

func (n *NegativeIntLiteralNode) AsInt64() (int64, bool) ***REMOVED***
	return n.Val, true
***REMOVED***

func (n *NegativeIntLiteralNode) AsUint64() (uint64, bool) ***REMOVED***
	if n.Val < 0 ***REMOVED***
		return 0, false
	***REMOVED***
	return uint64(n.Val), true
***REMOVED***

// FloatValueNode is an AST node that represents a numeric literal with
// a floating point, in scientific notation, or too large to fit in an
// int64 or uint64.
type FloatValueNode interface ***REMOVED***
	ValueNode
	AsFloat() float64
***REMOVED***

var _ FloatValueNode = (*FloatLiteralNode)(nil)
var _ FloatValueNode = (*SpecialFloatLiteralNode)(nil)
var _ FloatValueNode = (*UintLiteralNode)(nil)

// FloatLiteralNode represents a floating point numeric literal.
type FloatLiteralNode struct ***REMOVED***
	terminalNode
	// Val is the numeric value indicated by the literal
	Val float64
***REMOVED***

// NewFloatLiteralNode creates a new *FloatLiteralNode with the given val.
func NewFloatLiteralNode(val float64, info TokenInfo) *FloatLiteralNode ***REMOVED***
	return &FloatLiteralNode***REMOVED***
		terminalNode: info.asTerminalNode(),
		Val:          val,
	***REMOVED***
***REMOVED***

func (n *FloatLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.AsFloat()
***REMOVED***

func (n *FloatLiteralNode) AsFloat() float64 ***REMOVED***
	return n.Val
***REMOVED***

// SpecialFloatLiteralNode represents a special floating point numeric literal
// for "inf" and "nan" values.
type SpecialFloatLiteralNode struct ***REMOVED***
	*KeywordNode
	Val float64
***REMOVED***

// NewSpecialFloatLiteralNode returns a new *SpecialFloatLiteralNode for the
// given keyword, which must be "inf" or "nan".
func NewSpecialFloatLiteralNode(name *KeywordNode) *SpecialFloatLiteralNode ***REMOVED***
	var f float64
	if name.Val == "inf" ***REMOVED***
		f = math.Inf(1)
	***REMOVED*** else ***REMOVED***
		f = math.NaN()
	***REMOVED***
	return &SpecialFloatLiteralNode***REMOVED***
		KeywordNode: name,
		Val:         f,
	***REMOVED***
***REMOVED***

func (n *SpecialFloatLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.AsFloat()
***REMOVED***

func (n *SpecialFloatLiteralNode) AsFloat() float64 ***REMOVED***
	return n.Val
***REMOVED***

// SignedFloatLiteralNode represents a signed floating point number.
type SignedFloatLiteralNode struct ***REMOVED***
	compositeNode
	Sign  *RuneNode
	Float FloatValueNode
	Val   float64
***REMOVED***

// NewSignedFloatLiteralNode creates a new *SignedFloatLiteralNode. Both
// arguments must be non-nil.
func NewSignedFloatLiteralNode(sign *RuneNode, f FloatValueNode) *SignedFloatLiteralNode ***REMOVED***
	if sign == nil ***REMOVED***
		panic("sign is nil")
	***REMOVED***
	if f == nil ***REMOVED***
		panic("f is nil")
	***REMOVED***
	children := []Node***REMOVED***sign, f***REMOVED***
	val := f.AsFloat()
	if sign.Rune == '-' ***REMOVED***
		val = -val
	***REMOVED***
	return &SignedFloatLiteralNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Sign:  sign,
		Float: f,
		Val:   val,
	***REMOVED***
***REMOVED***

func (n *SignedFloatLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.Val
***REMOVED***

func (n *SignedFloatLiteralNode) AsFloat() float64 ***REMOVED***
	return n.Val
***REMOVED***

// BoolLiteralNode represents a boolean literal.
//
// Deprecated: The AST uses IdentNode for boolean literals, where the
// identifier value is "true" or "false". This is required because an
// identifier "true" is not necessarily a boolean value as it could also
// be an enum value named "true" (ditto for "false").
type BoolLiteralNode struct ***REMOVED***
	*KeywordNode
	Val bool
***REMOVED***

// NewBoolLiteralNode returns a new *BoolLiteralNode for the given keyword,
// which must be "true" or "false".
func NewBoolLiteralNode(name *KeywordNode) *BoolLiteralNode ***REMOVED***
	return &BoolLiteralNode***REMOVED***
		KeywordNode: name,
		Val:         name.Val == "true",
	***REMOVED***
***REMOVED***

func (n *BoolLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.Val
***REMOVED***

// ArrayLiteralNode represents an array literal, which is only allowed inside of
// a MessageLiteralNode, to indicate values for a repeated field. Example:
//
//  ["foo", "bar", "baz"]
type ArrayLiteralNode struct ***REMOVED***
	compositeNode
	OpenBracket *RuneNode
	Elements    []ValueNode
	// Commas represent the separating ',' characters between elements. The
	// length of this slice must be exactly len(Elements)-1, with each item
	// in Elements having a corresponding item in this slice *except the last*
	// (since a trailing comma is not allowed).
	Commas       []*RuneNode
	CloseBracket *RuneNode
***REMOVED***

// NewArrayLiteralNode creates a new *ArrayLiteralNode. The openBracket and
// closeBracket args must be non-nil and represent the "[" and "]" runes that
// surround the array values. The given commas arg must have a length that is
// one less than the length of the vals arg. However, vals may be empty, in
// which case commas must also be empty.
func NewArrayLiteralNode(openBracket *RuneNode, vals []ValueNode, commas []*RuneNode, closeBracket *RuneNode) *ArrayLiteralNode ***REMOVED***
	if openBracket == nil ***REMOVED***
		panic("openBracket is nil")
	***REMOVED***
	if closeBracket == nil ***REMOVED***
		panic("closeBracket is nil")
	***REMOVED***
	if len(vals) == 0 && len(commas) != 0 ***REMOVED***
		panic("vals is empty but commas is not")
	***REMOVED***
	if len(vals) > 0 && len(commas) != len(vals)-1 ***REMOVED***
		panic(fmt.Sprintf("%d vals requires %d commas, not %d", len(vals), len(vals)-1, len(commas)))
	***REMOVED***
	children := make([]Node, 0, len(vals)*2+1)
	children = append(children, openBracket)
	for i, val := range vals ***REMOVED***
		if i > 0 ***REMOVED***
			if commas[i-1] == nil ***REMOVED***
				panic(fmt.Sprintf("commas[%d] is nil", i-1))
			***REMOVED***
			children = append(children, commas[i-1])
		***REMOVED***
		if val == nil ***REMOVED***
			panic(fmt.Sprintf("vals[%d] is nil", i))
		***REMOVED***
		children = append(children, val)
	***REMOVED***
	children = append(children, closeBracket)

	return &ArrayLiteralNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		OpenBracket:  openBracket,
		Elements:     vals,
		Commas:       commas,
		CloseBracket: closeBracket,
	***REMOVED***
***REMOVED***

func (n *ArrayLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.Elements
***REMOVED***

// MessageLiteralNode represents a message literal, which is compatible with the
// protobuf text format and can be used for custom options with message types.
// Example:
//
//   ***REMOVED*** foo:1 foo:2 foo:3 bar:<name:"abc" id:123> ***REMOVED***
type MessageLiteralNode struct ***REMOVED***
	compositeNode
	Open     *RuneNode // should be '***REMOVED***' or '<'
	Elements []*MessageFieldNode
	// Separator characters between elements, which can be either ','
	// or ';' if present. This slice must be exactly len(Elements) in
	// length, with each item in Elements having one corresponding item
	// in Seps. Separators in message literals are optional, so a given
	// item in this slice may be nil to indicate absence of a separator.
	Seps  []*RuneNode
	Close *RuneNode // should be '***REMOVED***' or '>', depending on Open
***REMOVED***

// NewMessageLiteralNode creates a new *MessageLiteralNode. The openSym and
// closeSym runes must not be nil and should be "***REMOVED***" and "***REMOVED***" or "<" and ">".
//
// Unlike separators (dots and commas) used for other AST nodes that represent
// a list of elements, the seps arg must be the SAME length as vals, and it may
// contain nil values to indicate absence of a separator (in fact, it could be
// all nils).
func NewMessageLiteralNode(openSym *RuneNode, vals []*MessageFieldNode, seps []*RuneNode, closeSym *RuneNode) *MessageLiteralNode ***REMOVED***
	if openSym == nil ***REMOVED***
		panic("openSym is nil")
	***REMOVED***
	if closeSym == nil ***REMOVED***
		panic("closeSym is nil")
	***REMOVED***
	if len(seps) != len(vals) ***REMOVED***
		panic(fmt.Sprintf("%d vals requires %d commas, not %d", len(vals), len(vals), len(seps)))
	***REMOVED***
	numChildren := len(vals) + 2
	for _, sep := range seps ***REMOVED***
		if sep != nil ***REMOVED***
			numChildren++
		***REMOVED***
	***REMOVED***
	children := make([]Node, 0, numChildren)
	children = append(children, openSym)
	for i, val := range vals ***REMOVED***
		if val == nil ***REMOVED***
			panic(fmt.Sprintf("vals[%d] is nil", i))
		***REMOVED***
		children = append(children, val)
		if seps[i] != nil ***REMOVED***
			children = append(children, seps[i])
		***REMOVED***
	***REMOVED***
	children = append(children, closeSym)

	return &MessageLiteralNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Open:     openSym,
		Elements: vals,
		Seps:     seps,
		Close:    closeSym,
	***REMOVED***
***REMOVED***

func (n *MessageLiteralNode) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.Elements
***REMOVED***

// MessageFieldNode represents a single field (name and value) inside of a
// message literal. Example:
//
//   foo:"bar"
type MessageFieldNode struct ***REMOVED***
	compositeNode
	Name *FieldReferenceNode
	// Sep represents the ':' separator between the name and value. If
	// the value is a message literal (and thus starts with '<' or '***REMOVED***')
	// or an array literal (starting with '[') then the separator is
	// optional, and thus may be nil.
	Sep *RuneNode
	Val ValueNode
***REMOVED***

// NewMessageFieldNode creates a new *MessageFieldNode. All args except sep
// must be non-nil.
func NewMessageFieldNode(name *FieldReferenceNode, sep *RuneNode, val ValueNode) *MessageFieldNode ***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if val == nil ***REMOVED***
		panic("val is nil")
	***REMOVED***
	numChildren := 2
	if sep != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	children = append(children, name)
	if sep != nil ***REMOVED***
		children = append(children, sep)
	***REMOVED***
	children = append(children, val)

	return &MessageFieldNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Name: name,
		Sep:  sep,
		Val:  val,
	***REMOVED***
***REMOVED***
