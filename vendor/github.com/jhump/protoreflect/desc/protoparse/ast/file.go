package ast

import "fmt"

// FileDeclNode is a placeholder interface for AST nodes that represent files.
// This allows NoSourceNode to be used in place of *FileNode for some usages.
type FileDeclNode interface ***REMOVED***
	Node
	GetSyntax() Node
***REMOVED***

var _ FileDeclNode = (*FileNode)(nil)
var _ FileDeclNode = NoSourceNode***REMOVED******REMOVED***

// FileNode is the root of the AST hierarchy. It represents an entire
// protobuf source file.
type FileNode struct ***REMOVED***
	compositeNode
	Syntax *SyntaxNode // nil if file has no syntax declaration
	Decls  []FileElement

	// Any comments that follow the last token in the file.
	FinalComments []Comment
	// Any whitespace at the end of the file (after the last token or
	// last comment in the file).
	FinalWhitespace string
***REMOVED***

// NewFileElement creates a new *FileNode. The syntax parameter is optional. If it
// is absent, it means the file had no syntax declaration.
//
// This function panics if the concrete type of any element of decls is not
// from this package.
func NewFileNode(syntax *SyntaxNode, decls []FileElement) *FileNode ***REMOVED***
	numChildren := len(decls)
	if syntax != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	if syntax != nil ***REMOVED***
		children = append(children, syntax)
	***REMOVED***
	for _, decl := range decls ***REMOVED***
		children = append(children, decl)
	***REMOVED***

	for _, decl := range decls ***REMOVED***
		switch decl := decl.(type) ***REMOVED***
		case *PackageNode, *ImportNode, *OptionNode, *MessageNode,
			*EnumNode, *ExtendNode, *ServiceNode, *EmptyDeclNode:
		default:
			panic(fmt.Sprintf("invalid FileElement type: %T", decl))
		***REMOVED***
	***REMOVED***

	return &FileNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Syntax: syntax,
		Decls:  decls,
	***REMOVED***
***REMOVED***

func NewEmptyFileNode(filename string) *FileNode ***REMOVED***
	return &FileNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: []Node***REMOVED***NewNoSourceNode(filename)***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (f *FileNode) GetSyntax() Node ***REMOVED***
	return f.Syntax
***REMOVED***

// FileElement is an interface implemented by all AST nodes that are
// allowed as top-level declarations in the file.
type FileElement interface ***REMOVED***
	Node
	fileElement()
***REMOVED***

var _ FileElement = (*ImportNode)(nil)
var _ FileElement = (*PackageNode)(nil)
var _ FileElement = (*OptionNode)(nil)
var _ FileElement = (*MessageNode)(nil)
var _ FileElement = (*EnumNode)(nil)
var _ FileElement = (*ExtendNode)(nil)
var _ FileElement = (*ServiceNode)(nil)
var _ FileElement = (*EmptyDeclNode)(nil)

// SyntaxNode represents a syntax declaration, which if present must be
// the first non-comment content. Example:
//
//  syntax = "proto2";
//
// Files that don't have a syntax node are assumed to use proto2 syntax.
type SyntaxNode struct ***REMOVED***
	compositeNode
	Keyword   *KeywordNode
	Equals    *RuneNode
	Syntax    StringValueNode
	Semicolon *RuneNode
***REMOVED***

// NewSyntaxNode creates a new *SyntaxNode. All four arguments must be non-nil:
//  - keyword: The token corresponding to the "syntax" keyword.
//  - equals: The token corresponding to the "=" rune.
//  - syntax: The actual syntax value, e.g. "proto2" or "proto3".
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewSyntaxNode(keyword *KeywordNode, equals *RuneNode, syntax StringValueNode, semicolon *RuneNode) *SyntaxNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if equals == nil ***REMOVED***
		panic("equals is nil")
	***REMOVED***
	if syntax == nil ***REMOVED***
		panic("syntax is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	children := []Node***REMOVED***keyword, equals, syntax, semicolon***REMOVED***
	return &SyntaxNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Equals:    equals,
		Syntax:    syntax,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

// ImportNode represents an import statement. Example:
//
//  import "google/protobuf/empty.proto";
type ImportNode struct ***REMOVED***
	compositeNode
	Keyword *KeywordNode
	// Optional; if present indicates this is a public import
	Public *KeywordNode
	// Optional; if present indicates this is a weak import
	Weak      *KeywordNode
	Name      StringValueNode
	Semicolon *RuneNode
***REMOVED***

// NewImportNode creates a new *ImportNode. The public and weak arguments are optional
// and only one or the other (or neither) may be specified, not both. When public is
// non-nil, it indicates the "public" keyword in the import statement and means this is
// a public import. When weak is non-nil, it indicates the "weak" keyword in the import
// statement and means this is a weak import. When both are nil, this is a normal import.
// The other arguments must be non-nil:
//  - keyword: The token corresponding to the "import" keyword.
//  - public: The token corresponding to the optional "public" keyword.
//  - weak: The token corresponding to the optional "weak" keyword.
//  - name: The actual imported file name.
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewImportNode(keyword *KeywordNode, public *KeywordNode, weak *KeywordNode, name StringValueNode, semicolon *RuneNode) *ImportNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	numChildren := 3
	if public != nil || weak != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	children = append(children, keyword)
	if public != nil ***REMOVED***
		children = append(children, public)
	***REMOVED*** else if weak != nil ***REMOVED***
		children = append(children, weak)
	***REMOVED***
	children = append(children, name, semicolon)

	return &ImportNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Public:    public,
		Weak:      weak,
		Name:      name,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

func (*ImportNode) fileElement() ***REMOVED******REMOVED***

// PackageNode represents a package declaration. Example:
//
//  package foobar.com;
type PackageNode struct ***REMOVED***
	compositeNode
	Keyword   *KeywordNode
	Name      IdentValueNode
	Semicolon *RuneNode
***REMOVED***

func (*PackageNode) fileElement() ***REMOVED******REMOVED***

// NewPackageNode creates a new *PackageNode. All three arguments must be non-nil:
//  - keyword: The token corresponding to the "package" keyword.
//  - name: The package name declared for the file.
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewPackageNode(keyword *KeywordNode, name IdentValueNode, semicolon *RuneNode) *PackageNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	children := []Node***REMOVED***keyword, name, semicolon***REMOVED***
	return &PackageNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Name:      name,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***
