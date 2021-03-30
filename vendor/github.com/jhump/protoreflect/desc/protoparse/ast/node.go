package ast

// Node is the interface implemented by all nodes in the AST. It
// provides information about the span of this AST node in terms
// of location in the source file. It also provides information
// about all prior comments (attached as leading comments) and
// optional subsequent comments (attached as trailing comments).
type Node interface ***REMOVED***
	Start() *SourcePos
	End() *SourcePos
	LeadingComments() []Comment
	TrailingComments() []Comment
***REMOVED***

// TerminalNode represents a leaf in the AST. These represent
// the tokens/lexemes in the protobuf language. Comments and
// whitespace are accumulated by the lexer and associated with
// the following lexed token.
type TerminalNode interface ***REMOVED***
	Node
	// PopLeadingComment removes the first leading comment from this
	// token and returns it. If the node has no leading comments then
	// this method will panic.
	PopLeadingComment() Comment
	// PushTrailingComment appends the given comment to the token's
	// trailing comments.
	PushTrailingComment(Comment)
	// LeadingWhitespace returns any whitespace between the prior comment
	// (last leading comment), if any, or prior lexed token and this token.
	LeadingWhitespace() string
	// RawText returns the raw text of the token as read from the source.
	RawText() string
***REMOVED***

var _ TerminalNode = (*StringLiteralNode)(nil)
var _ TerminalNode = (*UintLiteralNode)(nil)
var _ TerminalNode = (*FloatLiteralNode)(nil)
var _ TerminalNode = (*IdentNode)(nil)
var _ TerminalNode = (*BoolLiteralNode)(nil)
var _ TerminalNode = (*SpecialFloatLiteralNode)(nil)
var _ TerminalNode = (*KeywordNode)(nil)
var _ TerminalNode = (*RuneNode)(nil)

// TokenInfo represents state accumulated by the lexer to associated with a
// token (aka terminal node).
type TokenInfo struct ***REMOVED***
	// The location of the token in the source file.
	PosRange
	// The raw text of the token.
	RawText string
	// Any comments encountered preceding this token.
	LeadingComments []Comment
	// Any leading whitespace immediately preceding this token.
	LeadingWhitespace string
	// Any trailing comments following this token. This is usually
	// empty as tokens are created by the lexer immediately and
	// trailing comments are accounted for afterwards, added using
	// the node's PushTrailingComment method.
	TrailingComments []Comment
***REMOVED***

func (t *TokenInfo) asTerminalNode() terminalNode ***REMOVED***
	return terminalNode***REMOVED***
		posRange:          t.PosRange,
		leadingComments:   t.LeadingComments,
		leadingWhitespace: t.LeadingWhitespace,
		trailingComments:  t.TrailingComments,
		raw:               t.RawText,
	***REMOVED***
***REMOVED***

// CompositeNode represents any non-terminal node in the tree. These
// are interior or root nodes and have child nodes.
type CompositeNode interface ***REMOVED***
	Node
	// All AST nodes that are immediate children of this one.
	Children() []Node
***REMOVED***

// terminalNode contains book-keeping shared by all TerminalNode
// implementations. It is embedded in all such node types in this
// package. It provides the implementation of the TerminalNode
// interface.
type terminalNode struct ***REMOVED***
	posRange          PosRange
	leadingComments   []Comment
	leadingWhitespace string
	trailingComments  []Comment
	raw               string
***REMOVED***

func (n *terminalNode) Start() *SourcePos ***REMOVED***
	return &n.posRange.Start
***REMOVED***

func (n *terminalNode) End() *SourcePos ***REMOVED***
	return &n.posRange.End
***REMOVED***

func (n *terminalNode) LeadingComments() []Comment ***REMOVED***
	return n.leadingComments
***REMOVED***

func (n *terminalNode) TrailingComments() []Comment ***REMOVED***
	return n.trailingComments
***REMOVED***

func (n *terminalNode) PopLeadingComment() Comment ***REMOVED***
	c := n.leadingComments[0]
	n.leadingComments = n.leadingComments[1:]
	return c
***REMOVED***

func (n *terminalNode) PushTrailingComment(c Comment) ***REMOVED***
	n.trailingComments = append(n.trailingComments, c)
***REMOVED***

func (n *terminalNode) LeadingWhitespace() string ***REMOVED***
	return n.leadingWhitespace
***REMOVED***

func (n *terminalNode) RawText() string ***REMOVED***
	return n.raw
***REMOVED***

// compositeNode contains book-keeping shared by all CompositeNode
// implementations. It is embedded in all such node types in this
// package. It provides the implementation of the CompositeNode
// interface.
type compositeNode struct ***REMOVED***
	children []Node
***REMOVED***

func (n *compositeNode) Children() []Node ***REMOVED***
	return n.children
***REMOVED***

func (n *compositeNode) Start() *SourcePos ***REMOVED***
	return n.children[0].Start()
***REMOVED***

func (n *compositeNode) End() *SourcePos ***REMOVED***
	return n.children[len(n.children)-1].End()
***REMOVED***

func (n *compositeNode) LeadingComments() []Comment ***REMOVED***
	return n.children[0].LeadingComments()
***REMOVED***

func (n *compositeNode) TrailingComments() []Comment ***REMOVED***
	return n.children[len(n.children)-1].TrailingComments()
***REMOVED***

// RuneNode represents a single rune in protobuf source. Runes
// are typically collected into tokens, but some runes stand on
// their own, such as punctuation/symbols like commas, semicolons,
// equals signs, open and close symbols (braces, brackets, angles,
// and parentheses), and periods/dots.
type RuneNode struct ***REMOVED***
	terminalNode
	Rune rune
***REMOVED***

// NewRuneNode creates a new *RuneNode with the given properties.
func NewRuneNode(r rune, info TokenInfo) *RuneNode ***REMOVED***
	return &RuneNode***REMOVED***
		terminalNode: info.asTerminalNode(),
		Rune:         r,
	***REMOVED***
***REMOVED***

// EmptyDeclNode represents an empty declaration in protobuf source.
// These amount to extra semicolons, with no actual content preceding
// the semicolon.
type EmptyDeclNode struct ***REMOVED***
	compositeNode
	Semicolon *RuneNode
***REMOVED***

// NewEmptyDeclNode creates a new *EmptyDeclNode. The one argument must
// be non-nil.
func NewEmptyDeclNode(semicolon *RuneNode) *EmptyDeclNode ***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	return &EmptyDeclNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: []Node***REMOVED***semicolon***REMOVED***,
		***REMOVED***,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

func (e *EmptyDeclNode) fileElement()    ***REMOVED******REMOVED***
func (e *EmptyDeclNode) msgElement()     ***REMOVED******REMOVED***
func (e *EmptyDeclNode) extendElement()  ***REMOVED******REMOVED***
func (e *EmptyDeclNode) oneOfElement()   ***REMOVED******REMOVED***
func (e *EmptyDeclNode) enumElement()    ***REMOVED******REMOVED***
func (e *EmptyDeclNode) serviceElement() ***REMOVED******REMOVED***
func (e *EmptyDeclNode) methodElement()  ***REMOVED******REMOVED***
