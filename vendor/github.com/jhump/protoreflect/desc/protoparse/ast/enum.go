package ast

import "fmt"

// EnumNode represents an enum declaration. Example:
//
//  enum Foo ***REMOVED*** BAR = 0; BAZ = 1 ***REMOVED***
type EnumNode struct ***REMOVED***
	compositeNode
	Keyword    *KeywordNode
	Name       *IdentNode
	OpenBrace  *RuneNode
	Decls      []EnumElement
	CloseBrace *RuneNode
***REMOVED***

func (*EnumNode) fileElement() ***REMOVED******REMOVED***
func (*EnumNode) msgElement()  ***REMOVED******REMOVED***

// NewEnumNode creates a new *EnumNode. All arguments must be non-nil. While
// it is technically allowed for decls to be nil or empty, the resulting node
// will not be a valid enum, which must have at least one value.
//  - keyword: The token corresponding to the "enum" keyword.
//  - name: The token corresponding to the enum's name.
//  - openBrace: The token corresponding to the "***REMOVED***" rune that starts the body.
//  - decls: All declarations inside the enum body.
//  - closeBrace: The token corresponding to the "***REMOVED***" rune that ends the body.
func NewEnumNode(keyword *KeywordNode, name *IdentNode, openBrace *RuneNode, decls []EnumElement, closeBrace *RuneNode) *EnumNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if openBrace == nil ***REMOVED***
		panic("openBrace is nil")
	***REMOVED***
	if closeBrace == nil ***REMOVED***
		panic("closeBrace is nil")
	***REMOVED***
	children := make([]Node, 0, 4+len(decls))
	children = append(children, keyword, name, openBrace)
	for _, decl := range decls ***REMOVED***
		children = append(children, decl)
	***REMOVED***
	children = append(children, closeBrace)

	for _, decl := range decls ***REMOVED***
		switch decl.(type) ***REMOVED***
		case *OptionNode, *EnumValueNode, *ReservedNode, *EmptyDeclNode:
		default:
			panic(fmt.Sprintf("invalid EnumElement type: %T", decl))
		***REMOVED***
	***REMOVED***

	return &EnumNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:    keyword,
		Name:       name,
		OpenBrace:  openBrace,
		CloseBrace: closeBrace,
		Decls:      decls,
	***REMOVED***
***REMOVED***

// EnumElement is an interface implemented by all AST nodes that can
// appear in the body of an enum declaration.
type EnumElement interface ***REMOVED***
	Node
	enumElement()
***REMOVED***

var _ EnumElement = (*OptionNode)(nil)
var _ EnumElement = (*EnumValueNode)(nil)
var _ EnumElement = (*ReservedNode)(nil)
var _ EnumElement = (*EmptyDeclNode)(nil)

// EnumValueDeclNode is a placeholder interface for AST nodes that represent
// enum values. This allows NoSourceNode to be used in place of *EnumValueNode
// for some usages.
type EnumValueDeclNode interface ***REMOVED***
	Node
	GetName() Node
	GetNumber() Node
***REMOVED***

var _ EnumValueDeclNode = (*EnumValueNode)(nil)
var _ EnumValueDeclNode = NoSourceNode***REMOVED******REMOVED***

// EnumNode represents an enum declaration. Example:
//
//  UNSET = 0 [deprecated = true];
type EnumValueNode struct ***REMOVED***
	compositeNode
	Name      *IdentNode
	Equals    *RuneNode
	Number    IntValueNode
	Options   *CompactOptionsNode
	Semicolon *RuneNode
***REMOVED***

func (*EnumValueNode) enumElement() ***REMOVED******REMOVED***

// NewEnumValueNode creates a new *EnumValueNode. All arguments must be non-nil
// except opts which is only non-nil if the declaration included options.
//  - name: The token corresponding to the enum value's name.
//  - equals: The token corresponding to the '=' rune after the name.
//  - number: The token corresponding to the enum value's number.
//  - opts: Optional set of enum value options.
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewEnumValueNode(name *IdentNode, equals *RuneNode, number IntValueNode, opts *CompactOptionsNode, semicolon *RuneNode) *EnumValueNode ***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if equals == nil ***REMOVED***
		panic("equals is nil")
	***REMOVED***
	if number == nil ***REMOVED***
		panic("number is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	numChildren := 4
	if opts != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	children = append(children, name, equals, number)
	if opts != nil ***REMOVED***
		children = append(children, opts)
	***REMOVED***
	children = append(children, semicolon)
	return &EnumValueNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Name:      name,
		Equals:    equals,
		Number:    number,
		Options:   opts,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

func (e *EnumValueNode) GetName() Node ***REMOVED***
	return e.Name
***REMOVED***

func (e *EnumValueNode) GetNumber() Node ***REMOVED***
	return e.Number
***REMOVED***
