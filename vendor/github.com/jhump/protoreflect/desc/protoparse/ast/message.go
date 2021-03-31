package ast

import "fmt"

// MessageDeclNode is a node in the AST that defines a message type. This
// includes normal message fields as well as implicit messages:
//  - *MessageNode
//  - *GroupNode (the group is a field and inline message type)
//  - *MapFieldNode (map fields implicitly define a MapEntry message type)
// This also allows NoSourceNode to be used in place of one of the above
// for some usages.
type MessageDeclNode interface ***REMOVED***
	Node
	MessageName() Node
***REMOVED***

var _ MessageDeclNode = (*MessageNode)(nil)
var _ MessageDeclNode = (*GroupNode)(nil)
var _ MessageDeclNode = (*MapFieldNode)(nil)
var _ MessageDeclNode = NoSourceNode***REMOVED******REMOVED***

// MessageNode represents a message declaration. Example:
//
//  message Foo ***REMOVED***
//    string name = 1;
//    repeated string labels = 2;
//    bytes extra = 3;
//  ***REMOVED***
type MessageNode struct ***REMOVED***
	compositeNode
	Keyword *KeywordNode
	Name    *IdentNode
	MessageBody
***REMOVED***

func (*MessageNode) fileElement() ***REMOVED******REMOVED***
func (*MessageNode) msgElement()  ***REMOVED******REMOVED***

// NewMessageNode creates a new *MessageNode. All arguments must be non-nil.
//  - keyword: The token corresponding to the "message" keyword.
//  - name: The token corresponding to the field's name.
//  - openBrace: The token corresponding to the "***REMOVED***" rune that starts the body.
//  - decls: All declarations inside the message body.
//  - closeBrace: The token corresponding to the "***REMOVED***" rune that ends the body.
func NewMessageNode(keyword *KeywordNode, name *IdentNode, openBrace *RuneNode, decls []MessageElement, closeBrace *RuneNode) *MessageNode ***REMOVED***
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

	ret := &MessageNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword: keyword,
		Name:    name,
	***REMOVED***
	populateMessageBody(&ret.MessageBody, openBrace, decls, closeBrace)
	return ret
***REMOVED***

func (n *MessageNode) MessageName() Node ***REMOVED***
	return n.Name
***REMOVED***

// MessageBody represents the body of a message. It is used by both
// MessageNodes and GroupNodes.
type MessageBody struct ***REMOVED***
	OpenBrace  *RuneNode
	Decls      []MessageElement
	CloseBrace *RuneNode
***REMOVED***

func populateMessageBody(m *MessageBody, openBrace *RuneNode, decls []MessageElement, closeBrace *RuneNode) ***REMOVED***
	m.OpenBrace = openBrace
	m.Decls = decls
	for _, decl := range decls ***REMOVED***
		switch decl.(type) ***REMOVED***
		case *OptionNode, *FieldNode, *MapFieldNode, *GroupNode, *OneOfNode,
			*MessageNode, *EnumNode, *ExtendNode, *ExtensionRangeNode,
			*ReservedNode, *EmptyDeclNode:
		default:
			panic(fmt.Sprintf("invalid MessageElement type: %T", decl))
		***REMOVED***
	***REMOVED***
	m.CloseBrace = closeBrace
***REMOVED***

// MessageElement is an interface implemented by all AST nodes that can
// appear in a message body.
type MessageElement interface ***REMOVED***
	Node
	msgElement()
***REMOVED***

var _ MessageElement = (*OptionNode)(nil)
var _ MessageElement = (*FieldNode)(nil)
var _ MessageElement = (*MapFieldNode)(nil)
var _ MessageElement = (*OneOfNode)(nil)
var _ MessageElement = (*GroupNode)(nil)
var _ MessageElement = (*MessageNode)(nil)
var _ MessageElement = (*EnumNode)(nil)
var _ MessageElement = (*ExtendNode)(nil)
var _ MessageElement = (*ExtensionRangeNode)(nil)
var _ MessageElement = (*ReservedNode)(nil)
var _ MessageElement = (*EmptyDeclNode)(nil)

// ExtendNode represents a declaration of extension fields. Example:
//
//  extend google.protobuf.FieldOptions ***REMOVED***
//    bool redacted = 33333;
//  ***REMOVED***
type ExtendNode struct ***REMOVED***
	compositeNode
	Keyword    *KeywordNode
	Extendee   IdentValueNode
	OpenBrace  *RuneNode
	Decls      []ExtendElement
	CloseBrace *RuneNode
***REMOVED***

func (*ExtendNode) fileElement() ***REMOVED******REMOVED***
func (*ExtendNode) msgElement()  ***REMOVED******REMOVED***

// NewExtendNode creates a new *ExtendNode. All arguments must be non-nil.
//  - keyword: The token corresponding to the "extend" keyword.
//  - extendee: The token corresponding to the name of the extended message.
//  - openBrace: The token corresponding to the "***REMOVED***" rune that starts the body.
//  - decls: All declarations inside the message body.
//  - closeBrace: The token corresponding to the "***REMOVED***" rune that ends the body.
func NewExtendNode(keyword *KeywordNode, extendee IdentValueNode, openBrace *RuneNode, decls []ExtendElement, closeBrace *RuneNode) *ExtendNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if extendee == nil ***REMOVED***
		panic("extendee is nil")
	***REMOVED***
	if openBrace == nil ***REMOVED***
		panic("openBrace is nil")
	***REMOVED***
	if closeBrace == nil ***REMOVED***
		panic("closeBrace is nil")
	***REMOVED***
	children := make([]Node, 0, 4+len(decls))
	children = append(children, keyword, extendee, openBrace)
	for _, decl := range decls ***REMOVED***
		children = append(children, decl)
	***REMOVED***
	children = append(children, closeBrace)

	ret := &ExtendNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:    keyword,
		Extendee:   extendee,
		OpenBrace:  openBrace,
		Decls:      decls,
		CloseBrace: closeBrace,
	***REMOVED***
	for _, decl := range decls ***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *FieldNode:
			decl.Extendee = ret
		case *GroupNode:
			decl.Extendee = ret
		case *EmptyDeclNode:
		default:
			panic(fmt.Sprintf("invalid ExtendElement type: %T", decl))
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

// ExtendElement is an interface implemented by all AST nodes that can
// appear in the body of an extends declaration.
type ExtendElement interface ***REMOVED***
	Node
	extendElement()
***REMOVED***

var _ ExtendElement = (*FieldNode)(nil)
var _ ExtendElement = (*GroupNode)(nil)
var _ ExtendElement = (*EmptyDeclNode)(nil)
