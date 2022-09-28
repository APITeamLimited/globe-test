package ast

import "fmt"

// ExtensionRangeNode represents an extension range declaration in an extendable
// message. Example:
//
//  extensions 100 to max;
type ExtensionRangeNode struct ***REMOVED***
	compositeNode
	Keyword *KeywordNode
	Ranges  []*RangeNode
	// Commas represent the separating ',' characters between ranges. The
	// length of this slice must be exactly len(Ranges)-1, each item in Ranges
	// having a corresponding item in this slice *except the last* (since a
	// trailing comma is not allowed).
	Commas    []*RuneNode
	Options   *CompactOptionsNode
	Semicolon *RuneNode
***REMOVED***

func (e *ExtensionRangeNode) msgElement() ***REMOVED******REMOVED***

// NewExtensionRangeNode creates a new *ExtensionRangeNode. All args must be
// non-nil except opts, which may be nil.
//  - keyword: The token corresponding to the "extends" keyword.
//  - ranges: One or more range expressions.
//  - commas: Tokens that represent the "," runes that delimit the range expressions.
//    The length of commas must be one less than the length of ranges.
//  - opts: The node corresponding to options that apply to each of the ranges.
//  - semicolon The token corresponding to the ";" rune that ends the declaration.
func NewExtensionRangeNode(keyword *KeywordNode, ranges []*RangeNode, commas []*RuneNode, opts *CompactOptionsNode, semicolon *RuneNode) *ExtensionRangeNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	if len(ranges) == 0 ***REMOVED***
		panic("must have at least one range")
	***REMOVED***
	if len(commas) != len(ranges)-1 ***REMOVED***
		panic(fmt.Sprintf("%d ranges requires %d commas, not %d", len(ranges), len(ranges)-1, len(commas)))
	***REMOVED***
	numChildren := len(ranges)*2 + 1
	if opts != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	children = append(children, keyword)
	for i, rng := range ranges ***REMOVED***
		if i > 0 ***REMOVED***
			if commas[i-1] == nil ***REMOVED***
				panic(fmt.Sprintf("commas[%d] is nil", i-1))
			***REMOVED***
			children = append(children, commas[i-1])
		***REMOVED***
		if rng == nil ***REMOVED***
			panic(fmt.Sprintf("ranges[%d] is nil", i))
		***REMOVED***
		children = append(children, rng)
	***REMOVED***
	if opts != nil ***REMOVED***
		children = append(children, opts)
	***REMOVED***
	children = append(children, semicolon)
	return &ExtensionRangeNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Ranges:    ranges,
		Commas:    commas,
		Options:   opts,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

// RangeDeclNode is a placeholder interface for AST nodes that represent
// numeric values. This allows NoSourceNode to be used in place of *RangeNode
// for some usages.
type RangeDeclNode interface ***REMOVED***
	Node
	RangeStart() Node
	RangeEnd() Node
***REMOVED***

var _ RangeDeclNode = (*RangeNode)(nil)
var _ RangeDeclNode = NoSourceNode***REMOVED******REMOVED***

// RangeNode represents a range expression, used in both extension ranges and
// reserved ranges. Example:
//
//  1000 to max
type RangeNode struct ***REMOVED***
	compositeNode
	StartVal IntValueNode
	// if To is non-nil, then exactly one of EndVal or Max must also be non-nil
	To *KeywordNode
	// EndVal and Max are mutually exclusive
	EndVal IntValueNode
	Max    *KeywordNode
***REMOVED***

// NewRangeNode creates a new *RangeNode. The start argument must be non-nil.
// The to argument represents the "to" keyword. If present (i.e. if it is non-nil),
// then so must be exactly one of end or max. If max is non-nil, it indicates a
// "100 to max" style range. But if end is non-nil, the end of the range is a
// literal, such as "100 to 200".
func NewRangeNode(start IntValueNode, to *KeywordNode, end IntValueNode, max *KeywordNode) *RangeNode ***REMOVED***
	if start == nil ***REMOVED***
		panic("start is nil")
	***REMOVED***
	numChildren := 1
	if to != nil ***REMOVED***
		if end == nil && max == nil ***REMOVED***
			panic("to is not nil, but end and max both are")
		***REMOVED***
		if end != nil && max != nil ***REMOVED***
			panic("end and max cannot be both non-nil")
		***REMOVED***
		numChildren = 3
	***REMOVED*** else ***REMOVED***
		if end != nil ***REMOVED***
			panic("to is nil, but end is not")
		***REMOVED***
		if max != nil ***REMOVED***
			panic("to is nil, but max is not")
		***REMOVED***
	***REMOVED***
	children := make([]Node, 0, numChildren)
	children = append(children, start)
	if to != nil ***REMOVED***
		children = append(children, to)
		if end != nil ***REMOVED***
			children = append(children, end)
		***REMOVED*** else ***REMOVED***
			children = append(children, max)
		***REMOVED***
	***REMOVED***
	return &RangeNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		StartVal: start,
		To:       to,
		EndVal:   end,
		Max:      max,
	***REMOVED***
***REMOVED***

func (n *RangeNode) RangeStart() Node ***REMOVED***
	return n.StartVal
***REMOVED***

func (n *RangeNode) RangeEnd() Node ***REMOVED***
	if n.Max != nil ***REMOVED***
		return n.Max
	***REMOVED***
	if n.EndVal != nil ***REMOVED***
		return n.EndVal
	***REMOVED***
	return n.StartVal
***REMOVED***

func (n *RangeNode) StartValue() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.StartVal.Value()
***REMOVED***

func (n *RangeNode) StartValueAsInt32(min, max int32) (int32, bool) ***REMOVED***
	return AsInt32(n.StartVal, min, max)
***REMOVED***

func (n *RangeNode) EndValue() interface***REMOVED******REMOVED*** ***REMOVED***
	if n.EndVal == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.EndVal.Value()
***REMOVED***

func (n *RangeNode) EndValueAsInt32(min, max int32) (int32, bool) ***REMOVED***
	if n.Max != nil ***REMOVED***
		return max, true
	***REMOVED***
	if n.EndVal == nil ***REMOVED***
		return n.StartValueAsInt32(min, max)
	***REMOVED***
	return AsInt32(n.EndVal, min, max)
***REMOVED***

// ReservedNode represents reserved declaration, which can be used to reserve
// either names or numbers. Examples:
//
//   reserved 1, 10-12, 15;
//   reserved "foo", "bar", "baz";
type ReservedNode struct ***REMOVED***
	compositeNode
	Keyword *KeywordNode
	// If non-empty, this node represents reserved ranges and Names will be empty.
	Ranges []*RangeNode
	// If non-empty, this node represents reserved names and Ranges will be empty.
	Names []StringValueNode
	// Commas represent the separating ',' characters between options. The
	// length of this slice must be exactly len(Ranges)-1 or len(Names)-1, depending
	// on whether this node represents reserved ranges or reserved names. Each item
	// in Ranges or Names has a corresponding item in this slice *except the last*
	// (since a trailing comma is not allowed).
	Commas    []*RuneNode
	Semicolon *RuneNode
***REMOVED***

func (*ReservedNode) msgElement()  ***REMOVED******REMOVED***
func (*ReservedNode) enumElement() ***REMOVED******REMOVED***

// NewReservedRangesNode creates a new *ReservedNode that represents reserved
// numeric ranges. All args must be non-nil.
//  - keyword: The token corresponding to the "reserved" keyword.
//  - ranges: One or more range expressions.
//  - commas: Tokens that represent the "," runes that delimit the range expressions.
//    The length of commas must be one less than the length of ranges.
//  - semicolon The token corresponding to the ";" rune that ends the declaration.
func NewReservedRangesNode(keyword *KeywordNode, ranges []*RangeNode, commas []*RuneNode, semicolon *RuneNode) *ReservedNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	if len(ranges) == 0 ***REMOVED***
		panic("must have at least one range")
	***REMOVED***
	if len(commas) != len(ranges)-1 ***REMOVED***
		panic(fmt.Sprintf("%d ranges requires %d commas, not %d", len(ranges), len(ranges)-1, len(commas)))
	***REMOVED***
	children := make([]Node, 0, len(ranges)*2+1)
	children = append(children, keyword)
	for i, rng := range ranges ***REMOVED***
		if i > 0 ***REMOVED***
			if commas[i-1] == nil ***REMOVED***
				panic(fmt.Sprintf("commas[%d] is nil", i-1))
			***REMOVED***
			children = append(children, commas[i-1])
		***REMOVED***
		if rng == nil ***REMOVED***
			panic(fmt.Sprintf("ranges[%d] is nil", i))
		***REMOVED***
		children = append(children, rng)
	***REMOVED***
	children = append(children, semicolon)
	return &ReservedNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Ranges:    ranges,
		Commas:    commas,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

// NewReservedNamesNode creates a new *ReservedNode that represents reserved
// names. All args must be non-nil.
//  - keyword: The token corresponding to the "reserved" keyword.
//  - names: One or more names.
//  - commas: Tokens that represent the "," runes that delimit the names.
//    The length of commas must be one less than the length of names.
//  - semicolon The token corresponding to the ";" rune that ends the declaration.
func NewReservedNamesNode(keyword *KeywordNode, names []StringValueNode, commas []*RuneNode, semicolon *RuneNode) *ReservedNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	if len(names) == 0 ***REMOVED***
		panic("must have at least one name")
	***REMOVED***
	if len(commas) != len(names)-1 ***REMOVED***
		panic(fmt.Sprintf("%d names requires %d commas, not %d", len(names), len(names)-1, len(commas)))
	***REMOVED***
	children := make([]Node, 0, len(names)*2+1)
	children = append(children, keyword)
	for i, name := range names ***REMOVED***
		if i > 0 ***REMOVED***
			if commas[i-1] == nil ***REMOVED***
				panic(fmt.Sprintf("commas[%d] is nil", i-1))
			***REMOVED***
			children = append(children, commas[i-1])
		***REMOVED***
		if name == nil ***REMOVED***
			panic(fmt.Sprintf("names[%d] is nil", i))
		***REMOVED***
		children = append(children, name)
	***REMOVED***
	children = append(children, semicolon)
	return &ReservedNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Names:     names,
		Commas:    commas,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***
