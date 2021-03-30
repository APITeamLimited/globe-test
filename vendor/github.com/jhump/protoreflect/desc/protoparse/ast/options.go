package ast

import "fmt"

// OptionDeclNode is a placeholder interface for AST nodes that represent
// options. This allows NoSourceNode to be used in place of *OptionNode
// for some usages.
type OptionDeclNode interface ***REMOVED***
	Node
	GetName() Node
	GetValue() ValueNode
***REMOVED***

var _ OptionDeclNode = (*OptionNode)(nil)
var _ OptionDeclNode = NoSourceNode***REMOVED******REMOVED***

// OptionNode represents the declaration of a single option for an element.
// It is used both for normal option declarations (start with "option" keyword
// and end with semicolon) and for compact options found in fields, enum values,
// and extension ranges. Example:
//
//  option (custom.option) = "foo";
type OptionNode struct ***REMOVED***
	compositeNode
	Keyword   *KeywordNode // absent for compact options
	Name      *OptionNameNode
	Equals    *RuneNode
	Val       ValueNode
	Semicolon *RuneNode // absent for compact options
***REMOVED***

func (e *OptionNode) fileElement()    ***REMOVED******REMOVED***
func (e *OptionNode) msgElement()     ***REMOVED******REMOVED***
func (e *OptionNode) oneOfElement()   ***REMOVED******REMOVED***
func (e *OptionNode) enumElement()    ***REMOVED******REMOVED***
func (e *OptionNode) serviceElement() ***REMOVED******REMOVED***
func (e *OptionNode) methodElement()  ***REMOVED******REMOVED***

// NewOptionNode creates a new *OptionNode for a full option declaration (as
// used in files, messages, oneofs, enums, services, and methods). All arguments
// must be non-nil. (Also see NewCompactOptionNode.)
//  - keyword: The token corresponding to the "option" keyword.
//  - name: The token corresponding to the name of the option.
//  - equals: The token corresponding to the "=" rune after the name.
//  - val: The token corresponding to the option value.
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewOptionNode(keyword *KeywordNode, name *OptionNameNode, equals *RuneNode, val ValueNode, semicolon *RuneNode) *OptionNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if equals == nil ***REMOVED***
		panic("equals is nil")
	***REMOVED***
	if val == nil ***REMOVED***
		panic("val is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	children := []Node***REMOVED***keyword, name, equals, val, semicolon***REMOVED***
	return &OptionNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:   keyword,
		Name:      name,
		Equals:    equals,
		Val:       val,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

// NewCompactOptionNode creates a new *OptionNode for a full compact declaration
// (as used in fields, enum values, and extension ranges). All arguments must be
// non-nil.
//  - name: The token corresponding to the name of the option.
//  - equals: The token corresponding to the "=" rune after the name.
//  - val: The token corresponding to the option value.
func NewCompactOptionNode(name *OptionNameNode, equals *RuneNode, val ValueNode) *OptionNode ***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if equals == nil ***REMOVED***
		panic("equals is nil")
	***REMOVED***
	if val == nil ***REMOVED***
		panic("val is nil")
	***REMOVED***
	children := []Node***REMOVED***name, equals, val***REMOVED***
	return &OptionNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Name:   name,
		Equals: equals,
		Val:    val,
	***REMOVED***
***REMOVED***

func (n *OptionNode) GetName() Node ***REMOVED***
	return n.Name
***REMOVED***

func (n *OptionNode) GetValue() ValueNode ***REMOVED***
	return n.Val
***REMOVED***

// OptionNameNode represents an option name or even a traversal through message
// types to name a nested option field. Example:
//
//   (foo.bar).baz.(bob)
type OptionNameNode struct ***REMOVED***
	compositeNode
	Parts []*FieldReferenceNode
	// Dots represent the separating '.' characters between name parts. The
	// length of this slice must be exactly len(Parts)-1, each item in Parts
	// having a corresponding item in this slice *except the last* (since a
	// trailing dot is not allowed).
	//
	// These do *not* include dots that are inside of an extension name. For
	// example: (foo.bar).baz.(bob) has three parts:
	//    1. (foo.bar)  - an extension name
	//    2. baz        - a regular field in foo.bar
	//    3. (bob)      - an extension field in baz
	// Note that the dot in foo.bar will thus not be present in Dots but is
	// instead in Parts[0].
	Dots []*RuneNode
***REMOVED***

// NewOptionNameNode creates a new *OptionNameNode. The dots arg must have a
// length that is one less than the length of parts. The parts arg must not be
// empty.
func NewOptionNameNode(parts []*FieldReferenceNode, dots []*RuneNode) *OptionNameNode ***REMOVED***
	if len(parts) == 0 ***REMOVED***
		panic("must have at least one part")
	***REMOVED***
	if len(dots) != len(parts)-1 ***REMOVED***
		panic(fmt.Sprintf("%d parts requires %d dots, not %d", len(parts), len(parts)-1, len(dots)))
	***REMOVED***
	children := make([]Node, 0, len(parts)*2-1)
	for i, part := range parts ***REMOVED***
		if part == nil ***REMOVED***
			panic(fmt.Sprintf("parts[%d] is nil", i))
		***REMOVED***
		if i > 0 ***REMOVED***
			if dots[i-1] == nil ***REMOVED***
				panic(fmt.Sprintf("dots[%d] is nil", i-1))
			***REMOVED***
			children = append(children, dots[i-1])
		***REMOVED***
		children = append(children, part)
	***REMOVED***
	return &OptionNameNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Parts: parts,
		Dots:  dots,
	***REMOVED***
***REMOVED***

// FieldReferenceNode is a reference to a field name. It can indicate a regular
// field (simple unqualified name) or an extension field (possibly-qualified
// name that is enclosed either in brackets or parentheses).
//
// This is used in options to indicate the names of custom options (which are
// actually extensions), in which case the name is enclosed in parentheses "("
// and ")". It is also used in message literals to set extension fields, in
// which case the name is enclosed in square brackets "[" and "]".
//
// Example:
//   (foo.bar)
type FieldReferenceNode struct ***REMOVED***
	compositeNode
	Open  *RuneNode // only present for extension names
	Name  IdentValueNode
	Close *RuneNode // only present for extension names
***REMOVED***

// NewFieldReferenceNode creates a new *FieldReferenceNode for a regular field.
// The name arg must not be nil.
func NewFieldReferenceNode(name *IdentNode) *FieldReferenceNode ***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	children := []Node***REMOVED***name***REMOVED***
	return &FieldReferenceNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Name: name,
	***REMOVED***
***REMOVED***

// NewExtensionFieldReferenceNode creates a new *FieldReferenceNode for an
// extension field. All args must be non-nil. The openSym and closeSym runes
// should be "(" and ")" or "[" and "]".
func NewExtensionFieldReferenceNode(openSym *RuneNode, name IdentValueNode, closeSym *RuneNode) *FieldReferenceNode ***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if openSym == nil ***REMOVED***
		panic("openSym is nil")
	***REMOVED***
	if closeSym == nil ***REMOVED***
		panic("closeSym is nil")
	***REMOVED***
	children := []Node***REMOVED***openSym, name, closeSym***REMOVED***
	return &FieldReferenceNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Open:  openSym,
		Name:  name,
		Close: closeSym,
	***REMOVED***
***REMOVED***

// IsExtension reports if this is an extension name or not (e.g. enclosed in
// punctuation, such as parentheses or brackets).
func (a *FieldReferenceNode) IsExtension() bool ***REMOVED***
	return a.Open != nil
***REMOVED***

func (a *FieldReferenceNode) Value() string ***REMOVED***
	if a.Open != nil ***REMOVED***
		return string(a.Open.Rune) + string(a.Name.AsIdentifier()) + string(a.Close.Rune)
	***REMOVED*** else ***REMOVED***
		return string(a.Name.AsIdentifier())
	***REMOVED***
***REMOVED***

// CompactOptionsNode represents a compact options declaration, as used with
// fields, enum values, and extension ranges. Example:
//
//  [deprecated = true, json_name = "foo_bar"]
type CompactOptionsNode struct ***REMOVED***
	compositeNode
	OpenBracket *RuneNode
	Options     []*OptionNode
	// Commas represent the separating ',' characters between options. The
	// length of this slice must be exactly len(Options)-1, with each item
	// in Options having a corresponding item in this slice *except the last*
	// (since a trailing comma is not allowed).
	Commas       []*RuneNode
	CloseBracket *RuneNode
***REMOVED***

// NewCompactOptionsNode creates a *CompactOptionsNode. All args must be
// non-nil. The commas arg must have a length that is one less than the
// length of opts. The opts arg must not be empty.
func NewCompactOptionsNode(openBracket *RuneNode, opts []*OptionNode, commas []*RuneNode, closeBracket *RuneNode) *CompactOptionsNode ***REMOVED***
	if openBracket == nil ***REMOVED***
		panic("openBracket is nil")
	***REMOVED***
	if closeBracket == nil ***REMOVED***
		panic("closeBracket is nil")
	***REMOVED***
	if len(opts) == 0 ***REMOVED***
		panic("must have at least one part")
	***REMOVED***
	if len(commas) != len(opts)-1 ***REMOVED***
		panic(fmt.Sprintf("%d opts requires %d commas, not %d", len(opts), len(opts)-1, len(commas)))
	***REMOVED***
	children := make([]Node, 0, len(opts)*2+1)
	children = append(children, openBracket)
	for i, opt := range opts ***REMOVED***
		if i > 0 ***REMOVED***
			if commas[i-1] == nil ***REMOVED***
				panic(fmt.Sprintf("commas[%d] is nil", i-1))
			***REMOVED***
			children = append(children, commas[i-1])
		***REMOVED***
		if opt == nil ***REMOVED***
			panic(fmt.Sprintf("opts[%d] is nil", i))
		***REMOVED***
		children = append(children, opt)
	***REMOVED***
	children = append(children, closeBracket)

	return &CompactOptionsNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		OpenBracket:  openBracket,
		Options:      opts,
		Commas:       commas,
		CloseBracket: closeBracket,
	***REMOVED***
***REMOVED***

func (e *CompactOptionsNode) GetElements() []*OptionNode ***REMOVED***
	if e == nil ***REMOVED***
		return nil
	***REMOVED***
	return e.Options
***REMOVED***
