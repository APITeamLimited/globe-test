package ast

import "fmt"

// ServiceNode represents a service declaration. Example:
//
//  service Foo ***REMOVED***
//    rpc Bar (Baz) returns (Bob);
//    rpc Frobnitz (stream Parts) returns (Gyzmeaux);
//  ***REMOVED***
type ServiceNode struct ***REMOVED***
	compositeNode
	Keyword    *KeywordNode
	Name       *IdentNode
	OpenBrace  *RuneNode
	Decls      []ServiceElement
	CloseBrace *RuneNode
***REMOVED***

func (*ServiceNode) fileElement() ***REMOVED******REMOVED***

// NewServiceNode creates a new *ServiceNode. All arguments must be non-nil.
//  - keyword: The token corresponding to the "service" keyword.
//  - name: The token corresponding to the service's name.
//  - openBrace: The token corresponding to the "***REMOVED***" rune that starts the body.
//  - decls: All declarations inside the service body.
//  - closeBrace: The token corresponding to the "***REMOVED***" rune that ends the body.
func NewServiceNode(keyword *KeywordNode, name *IdentNode, openBrace *RuneNode, decls []ServiceElement, closeBrace *RuneNode) *ServiceNode ***REMOVED***
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
		switch decl := decl.(type) ***REMOVED***
		case *OptionNode, *RPCNode, *EmptyDeclNode:
		default:
			panic(fmt.Sprintf("invalid ServiceElement type: %T", decl))
		***REMOVED***
	***REMOVED***

	return &ServiceNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:    keyword,
		Name:       name,
		OpenBrace:  openBrace,
		Decls:      decls,
		CloseBrace: closeBrace,
	***REMOVED***
***REMOVED***

// ServiceElement is an interface implemented by all AST nodes that can
// appear in the body of a service declaration.
type ServiceElement interface ***REMOVED***
	Node
	serviceElement()
***REMOVED***

var _ ServiceElement = (*OptionNode)(nil)
var _ ServiceElement = (*RPCNode)(nil)
var _ ServiceElement = (*EmptyDeclNode)(nil)

// RPCDeclNode is a placeholder interface for AST nodes that represent RPC
// declarations. This allows NoSourceNode to be used in place of *RPCNode
// for some usages.
type RPCDeclNode interface ***REMOVED***
	Node
	GetInputType() Node
	GetOutputType() Node
***REMOVED***

var _ RPCDeclNode = (*RPCNode)(nil)
var _ RPCDeclNode = NoSourceNode***REMOVED******REMOVED***

// RPCNode represents an RPC declaration. Example:
//
//  rpc Foo (Bar) returns (Baz);
type RPCNode struct ***REMOVED***
	compositeNode
	Keyword    *KeywordNode
	Name       *IdentNode
	Input      *RPCTypeNode
	Returns    *KeywordNode
	Output     *RPCTypeNode
	Semicolon  *RuneNode
	OpenBrace  *RuneNode
	Decls      []RPCElement
	CloseBrace *RuneNode
***REMOVED***

func (n *RPCNode) serviceElement() ***REMOVED******REMOVED***

// NewRPCNode creates a new *RPCNode with no body. All arguments must be non-nil.
//  - keyword: The token corresponding to the "rpc" keyword.
//  - name: The token corresponding to the RPC's name.
//  - input: The token corresponding to the RPC input message type.
//  - returns: The token corresponding to the "returns" keyword that precedes the output type.
//  - output: The token corresponding to the RPC output message type.
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewRPCNode(keyword *KeywordNode, name *IdentNode, input *RPCTypeNode, returns *KeywordNode, output *RPCTypeNode, semicolon *RuneNode) *RPCNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if input == nil ***REMOVED***
		panic("input is nil")
	***REMOVED***
	if returns == nil ***REMOVED***
		panic("returns is nil")
	***REMOVED***
	if output == nil ***REMOVED***
		panic("output is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	children := []Node***REMOVED***keyword, name, input, returns, output, semicolon***REMOVED***
	return &RPCNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Name:      name,
		Input:     input,
		Returns:   returns,
		Output:    output,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

// NewRPCNodeWithBody creates a new *RPCNode that includes a body (and possibly
// options). All arguments must be non-nil.
//  - keyword: The token corresponding to the "rpc" keyword.
//  - name: The token corresponding to the RPC's name.
//  - input: The token corresponding to the RPC input message type.
//  - returns: The token corresponding to the "returns" keyword that precedes the output type.
//  - output: The token corresponding to the RPC output message type.
//  - openBrace: The token corresponding to the "***REMOVED***" rune that starts the body.
//  - decls: All declarations inside the RPC body.
//  - closeBrace: The token corresponding to the "***REMOVED***" rune that ends the body.
func NewRPCNodeWithBody(keyword *KeywordNode, name *IdentNode, input *RPCTypeNode, returns *KeywordNode, output *RPCTypeNode, openBrace *RuneNode, decls []RPCElement, closeBrace *RuneNode) *RPCNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if input == nil ***REMOVED***
		panic("input is nil")
	***REMOVED***
	if returns == nil ***REMOVED***
		panic("returns is nil")
	***REMOVED***
	if output == nil ***REMOVED***
		panic("output is nil")
	***REMOVED***
	if openBrace == nil ***REMOVED***
		panic("openBrace is nil")
	***REMOVED***
	if closeBrace == nil ***REMOVED***
		panic("closeBrace is nil")
	***REMOVED***
	children := make([]Node, 0, 7+len(decls))
	children = append(children, keyword, name, input, returns, output, openBrace)
	for _, decl := range decls ***REMOVED***
		children = append(children, decl)
	***REMOVED***
	children = append(children, closeBrace)

	for _, decl := range decls ***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *OptionNode, *EmptyDeclNode:
		default:
			panic(fmt.Sprintf("invalid RPCElement type: %T", decl))
		***REMOVED***
	***REMOVED***

	return &RPCNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:    keyword,
		Name:       name,
		Input:      input,
		Returns:    returns,
		Output:     output,
		OpenBrace:  openBrace,
		Decls:      decls,
		CloseBrace: closeBrace,
	***REMOVED***
***REMOVED***

func (n *RPCNode) GetInputType() Node ***REMOVED***
	return n.Input.MessageType
***REMOVED***

func (n *RPCNode) GetOutputType() Node ***REMOVED***
	return n.Output.MessageType
***REMOVED***

// RPCElement is an interface implemented by all AST nodes that can
// appear in the body of an rpc declaration (aka method).
type RPCElement interface ***REMOVED***
	Node
	methodElement()
***REMOVED***

var _ RPCElement = (*OptionNode)(nil)
var _ RPCElement = (*EmptyDeclNode)(nil)

// RPCTypeNode represents the declaration of a request or response type for an
// RPC. Example:
//
//  (stream foo.Bar)
type RPCTypeNode struct ***REMOVED***
	compositeNode
	OpenParen   *RuneNode
	Stream      *KeywordNode
	MessageType IdentValueNode
	CloseParen  *RuneNode
***REMOVED***

// NewRPCTypeNode creates a new *RPCTypeNode. All arguments must be non-nil
// except stream, which may be nil.
//  - openParen: The token corresponding to the "(" rune that starts the declaration.
//  - stream: The token corresponding to the "stream" keyword or nil if not present.
//  - msgType: The token corresponding to the message type's name.
//  - closeParen: The token corresponding to the ")" rune that ends the declaration.
func NewRPCTypeNode(openParen *RuneNode, stream *KeywordNode, msgType IdentValueNode, closeParen *RuneNode) *RPCTypeNode ***REMOVED***
	if openParen == nil ***REMOVED***
		panic("openParen is nil")
	***REMOVED***
	if msgType == nil ***REMOVED***
		panic("msgType is nil")
	***REMOVED***
	if closeParen == nil ***REMOVED***
		panic("closeParen is nil")
	***REMOVED***
	var children []Node
	if stream != nil ***REMOVED***
		children = []Node***REMOVED***openParen, stream, msgType, closeParen***REMOVED***
	***REMOVED*** else ***REMOVED***
		children = []Node***REMOVED***openParen, msgType, closeParen***REMOVED***
	***REMOVED***

	return &RPCTypeNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		OpenParen:   openParen,
		Stream:      stream,
		MessageType: msgType,
		CloseParen:  closeParen,
	***REMOVED***
***REMOVED***
