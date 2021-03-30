package ast

import "fmt"

// FieldDeclNode is a node in the AST that defines a field. This includes
// normal message fields as well as extensions. There are multiple types
// of AST nodes that declare fields:
//  - *FieldNode
//  - *GroupNode
//  - *MapFieldNode
// This also allows NoSourceNode to be used in place of one of the above
// for some usages.
type FieldDeclNode interface ***REMOVED***
	Node
	FieldLabel() Node
	FieldName() Node
	FieldType() Node
	FieldTag() Node
	FieldExtendee() Node
	GetGroupKeyword() Node
	GetOptions() *CompactOptionsNode
***REMOVED***

var _ FieldDeclNode = (*FieldNode)(nil)
var _ FieldDeclNode = (*GroupNode)(nil)
var _ FieldDeclNode = (*MapFieldNode)(nil)
var _ FieldDeclNode = (*SyntheticMapField)(nil)
var _ FieldDeclNode = NoSourceNode***REMOVED******REMOVED***

// FieldNode represents a normal field declaration (not groups or maps). It
// can represent extension fields as well as non-extension fields (both inside
// of messages and inside of one-ofs). Example:
//
//  optional string foo = 1;
type FieldNode struct ***REMOVED***
	compositeNode
	Label     FieldLabel
	FldType   IdentValueNode
	Name      *IdentNode
	Equals    *RuneNode
	Tag       *UintLiteralNode
	Options   *CompactOptionsNode
	Semicolon *RuneNode

	// This is an up-link to the containing *ExtendNode for fields
	// that are defined inside of "extend" blocks.
	Extendee *ExtendNode
***REMOVED***

func (*FieldNode) msgElement()    ***REMOVED******REMOVED***
func (*FieldNode) oneOfElement()  ***REMOVED******REMOVED***
func (*FieldNode) extendElement() ***REMOVED******REMOVED***

// NewFieldNode creates a new *FieldNode. The label and options arguments may be
// nil but the others must be non-nil.
//  - label: The token corresponding to the label keyword if present ("optional",
//    "required", or "repeated").
//  - fieldType: The token corresponding to the field's type.
//  - name: The token corresponding to the field's name.
//  - equals: The token corresponding to the '=' rune after the name.
//  - tag: The token corresponding to the field's tag number.
//  - opts: Optional set of field options.
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewFieldNode(label *KeywordNode, fieldType IdentValueNode, name *IdentNode, equals *RuneNode, tag *UintLiteralNode, opts *CompactOptionsNode, semicolon *RuneNode) *FieldNode ***REMOVED***
	if fieldType == nil ***REMOVED***
		panic("fieldType is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if equals == nil ***REMOVED***
		panic("equals is nil")
	***REMOVED***
	if tag == nil ***REMOVED***
		panic("tag is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	numChildren := 5
	if label != nil ***REMOVED***
		numChildren++
	***REMOVED***
	if opts != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	if label != nil ***REMOVED***
		children = append(children, label)
	***REMOVED***
	children = append(children, fieldType, name, equals, tag)
	if opts != nil ***REMOVED***
		children = append(children, opts)
	***REMOVED***
	children = append(children, semicolon)

	return &FieldNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Label:     newFieldLabel(label),
		FldType:   fieldType,
		Name:      name,
		Equals:    equals,
		Tag:       tag,
		Options:   opts,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

func (n *FieldNode) FieldLabel() Node ***REMOVED***
	// proto3 fields and fields inside one-ofs will not have a label and we need
	// this check in order to return a nil node -- otherwise we'd return a
	// non-nil node that has a nil pointer value in it :/
	if n.Label.KeywordNode == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.Label.KeywordNode
***REMOVED***

func (n *FieldNode) FieldName() Node ***REMOVED***
	return n.Name
***REMOVED***

func (n *FieldNode) FieldType() Node ***REMOVED***
	return n.FldType
***REMOVED***

func (n *FieldNode) FieldTag() Node ***REMOVED***
	return n.Tag
***REMOVED***

func (n *FieldNode) FieldExtendee() Node ***REMOVED***
	if n.Extendee != nil ***REMOVED***
		return n.Extendee.Extendee
	***REMOVED***
	return nil
***REMOVED***

func (n *FieldNode) GetGroupKeyword() Node ***REMOVED***
	return nil
***REMOVED***

func (n *FieldNode) GetOptions() *CompactOptionsNode ***REMOVED***
	return n.Options
***REMOVED***

// FieldLabel represents the label of a field, which indicates its cardinality
// (i.e. whether it is optional, required, or repeated).
type FieldLabel struct ***REMOVED***
	*KeywordNode
	Repeated bool
	Required bool
***REMOVED***

func newFieldLabel(lbl *KeywordNode) FieldLabel ***REMOVED***
	repeated, required := false, false
	if lbl != nil ***REMOVED***
		repeated = lbl.Val == "repeated"
		required = lbl.Val == "required"
	***REMOVED***
	return FieldLabel***REMOVED***
		KeywordNode: lbl,
		Repeated:    repeated,
		Required:    required,
	***REMOVED***
***REMOVED***

// IsPresent returns true if a label keyword was present in the declaration
// and false if it was absent.
func (f *FieldLabel) IsPresent() bool ***REMOVED***
	return f.KeywordNode != nil
***REMOVED***

// GroupNode represents a group declaration, which doubles as a field and inline
// message declaration. It can represent extension fields as well as
// non-extension fields (both inside of messages and inside of one-ofs).
// Example:
//
//  optional group Key = 4 ***REMOVED***
//    optional uint64 id = 1;
//    optional string name = 2;
//  ***REMOVED***
type GroupNode struct ***REMOVED***
	compositeNode
	Label   FieldLabel
	Keyword *KeywordNode
	Name    *IdentNode
	Equals  *RuneNode
	Tag     *UintLiteralNode
	Options *CompactOptionsNode
	MessageBody

	// This is an up-link to the containing *ExtendNode for groups
	// that are defined inside of "extend" blocks.
	Extendee *ExtendNode
***REMOVED***

func (*GroupNode) msgElement()    ***REMOVED******REMOVED***
func (*GroupNode) oneOfElement()  ***REMOVED******REMOVED***
func (*GroupNode) extendElement() ***REMOVED******REMOVED***

// NewGroupNode creates a new *GroupNode. The label and options arguments may be
// nil but the others must be non-nil.
//  - label: The token corresponding to the label keyword if present ("optional",
//    "required", or "repeated").
//  - keyword: The token corresponding to the "group" keyword.
//  - name: The token corresponding to the field's name.
//  - equals: The token corresponding to the '=' rune after the name.
//  - tag: The token corresponding to the field's tag number.
//  - opts: Optional set of field options.
//  - openBrace: The token corresponding to the "***REMOVED***" rune that starts the body.
//  - decls: All declarations inside the group body.
//  - closeBrace: The token corresponding to the "***REMOVED***" rune that ends the body.
func NewGroupNode(label *KeywordNode, keyword *KeywordNode, name *IdentNode, equals *RuneNode, tag *UintLiteralNode, opts *CompactOptionsNode, openBrace *RuneNode, decls []MessageElement, closeBrace *RuneNode) *GroupNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("fieldType is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if equals == nil ***REMOVED***
		panic("equals is nil")
	***REMOVED***
	if tag == nil ***REMOVED***
		panic("tag is nil")
	***REMOVED***
	if openBrace == nil ***REMOVED***
		panic("openBrace is nil")
	***REMOVED***
	if closeBrace == nil ***REMOVED***
		panic("closeBrace is nil")
	***REMOVED***
	numChildren := 6 + len(decls)
	if label != nil ***REMOVED***
		numChildren++
	***REMOVED***
	if opts != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	if label != nil ***REMOVED***
		children = append(children, label)
	***REMOVED***
	children = append(children, keyword, name, equals, tag)
	if opts != nil ***REMOVED***
		children = append(children, opts)
	***REMOVED***
	children = append(children, openBrace)
	for _, decl := range decls ***REMOVED***
		children = append(children, decl)
	***REMOVED***
	children = append(children, closeBrace)

	ret := &GroupNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Label:   newFieldLabel(label),
		Keyword: keyword,
		Name:    name,
		Equals:  equals,
		Tag:     tag,
		Options: opts,
	***REMOVED***
	populateMessageBody(&ret.MessageBody, openBrace, decls, closeBrace)
	return ret
***REMOVED***

func (n *GroupNode) FieldLabel() Node ***REMOVED***
	if n.Label.KeywordNode == nil ***REMOVED***
		// return nil interface to indicate absence, not a typed nil
		return nil
	***REMOVED***
	return n.Label.KeywordNode
***REMOVED***

func (n *GroupNode) FieldName() Node ***REMOVED***
	return n.Name
***REMOVED***

func (n *GroupNode) FieldType() Node ***REMOVED***
	return n.Keyword
***REMOVED***

func (n *GroupNode) FieldTag() Node ***REMOVED***
	return n.Tag
***REMOVED***

func (n *GroupNode) FieldExtendee() Node ***REMOVED***
	if n.Extendee != nil ***REMOVED***
		return n.Extendee.Extendee
	***REMOVED***
	return nil
***REMOVED***

func (n *GroupNode) GetGroupKeyword() Node ***REMOVED***
	return n.Keyword
***REMOVED***

func (n *GroupNode) GetOptions() *CompactOptionsNode ***REMOVED***
	return n.Options
***REMOVED***

func (n *GroupNode) MessageName() Node ***REMOVED***
	return n.Name
***REMOVED***

// OneOfNode represents a one-of declaration. Example:
//
//  oneof query ***REMOVED***
//    string by_name = 2;
//    Type by_type = 3;
//    Address by_address = 4;
//    Labels by_label = 5;
//  ***REMOVED***
type OneOfNode struct ***REMOVED***
	compositeNode
	Keyword    *KeywordNode
	Name       *IdentNode
	OpenBrace  *RuneNode
	Decls      []OneOfElement
	CloseBrace *RuneNode
***REMOVED***

func (*OneOfNode) msgElement() ***REMOVED******REMOVED***

// NewOneOfNode creates a new *OneOfNode. All arguments must be non-nil. While
// it is technically allowed for decls to be nil or empty, the resulting node
// will not be a valid oneof, which must have at least one field.
//  - keyword: The token corresponding to the "oneof" keyword.
//  - name: The token corresponding to the oneof's name.
//  - openBrace: The token corresponding to the "***REMOVED***" rune that starts the body.
//  - decls: All declarations inside the oneof body.
//  - closeBrace: The token corresponding to the "***REMOVED***" rune that ends the body.
func NewOneOfNode(keyword *KeywordNode, name *IdentNode, openBrace *RuneNode, decls []OneOfElement, closeBrace *RuneNode) *OneOfNode ***REMOVED***
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
		case *OptionNode, *FieldNode, *GroupNode, *EmptyDeclNode:
		default:
			panic(fmt.Sprintf("invalid OneOfElement type: %T", decl))
		***REMOVED***
	***REMOVED***

	return &OneOfNode***REMOVED***
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

// OneOfElement is an interface implemented by all AST nodes that can
// appear in the body of a oneof declaration.
type OneOfElement interface ***REMOVED***
	Node
	oneOfElement()
***REMOVED***

var _ OneOfElement = (*OptionNode)(nil)
var _ OneOfElement = (*FieldNode)(nil)
var _ OneOfElement = (*GroupNode)(nil)
var _ OneOfElement = (*EmptyDeclNode)(nil)

// MapTypeNode represents the type declaration for a map field. It defines
// both the key and value types for the map. Example:
//
//  map<string, Values>
type MapTypeNode struct ***REMOVED***
	compositeNode
	Keyword    *KeywordNode
	OpenAngle  *RuneNode
	KeyType    *IdentNode
	Comma      *RuneNode
	ValueType  IdentValueNode
	CloseAngle *RuneNode
***REMOVED***

// NewMapTypeNode creates a new *MapTypeNode. All arguments must be non-nil.
//  - keyword: The token corresponding to the "map" keyword.
//  - openAngle: The token corresponding to the "<" rune after the keyword.
//  - keyType: The token corresponding to the key type for the map.
//  - comma: The token corresponding to the "," rune between key and value types.
//  - valType: The token corresponding to the value type for the map.
//  - closeAngle: The token corresponding to the ">" rune that ends the declaration.
func NewMapTypeNode(keyword *KeywordNode, openAngle *RuneNode, keyType *IdentNode, comma *RuneNode, valType IdentValueNode, closeAngle *RuneNode) *MapTypeNode ***REMOVED***
	if keyword == nil ***REMOVED***
		panic("keyword is nil")
	***REMOVED***
	if openAngle == nil ***REMOVED***
		panic("openAngle is nil")
	***REMOVED***
	if keyType == nil ***REMOVED***
		panic("keyType is nil")
	***REMOVED***
	if comma == nil ***REMOVED***
		panic("comma is nil")
	***REMOVED***
	if valType == nil ***REMOVED***
		panic("valType is nil")
	***REMOVED***
	if closeAngle == nil ***REMOVED***
		panic("closeAngle is nil")
	***REMOVED***
	children := []Node***REMOVED***keyword, openAngle, keyType, comma, valType, closeAngle***REMOVED***
	return &MapTypeNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		Keyword:    keyword,
		OpenAngle:  openAngle,
		KeyType:    keyType,
		Comma:      comma,
		ValueType:  valType,
		CloseAngle: closeAngle,
	***REMOVED***
***REMOVED***

// MapFieldNode represents a map field declaration. Example:
//
//  map<string,string> replacements = 3 [deprecated = true];
type MapFieldNode struct ***REMOVED***
	compositeNode
	MapType   *MapTypeNode
	Name      *IdentNode
	Equals    *RuneNode
	Tag       *UintLiteralNode
	Options   *CompactOptionsNode
	Semicolon *RuneNode
***REMOVED***

func (*MapFieldNode) msgElement() ***REMOVED******REMOVED***

// NewMapFieldNode creates a new *MapFieldNode. All arguments must be non-nil
// except opts, which may be nil.
//  - mapType: The token corresponding to the map type.
//  - name: The token corresponding to the field's name.
//  - equals: The token corresponding to the '=' rune after the name.
//  - tag: The token corresponding to the field's tag number.
//  - opts: Optional set of field options.
//  - semicolon: The token corresponding to the ";" rune that ends the declaration.
func NewMapFieldNode(mapType *MapTypeNode, name *IdentNode, equals *RuneNode, tag *UintLiteralNode, opts *CompactOptionsNode, semicolon *RuneNode) *MapFieldNode ***REMOVED***
	if mapType == nil ***REMOVED***
		panic("mapType is nil")
	***REMOVED***
	if name == nil ***REMOVED***
		panic("name is nil")
	***REMOVED***
	if equals == nil ***REMOVED***
		panic("equals is nil")
	***REMOVED***
	if tag == nil ***REMOVED***
		panic("tag is nil")
	***REMOVED***
	if semicolon == nil ***REMOVED***
		panic("semicolon is nil")
	***REMOVED***
	numChildren := 5
	if opts != nil ***REMOVED***
		numChildren++
	***REMOVED***
	children := make([]Node, 0, numChildren)
	children = append(children, mapType, name, equals, tag)
	if opts != nil ***REMOVED***
		children = append(children, opts)
	***REMOVED***
	children = append(children, semicolon)

	return &MapFieldNode***REMOVED***
		compositeNode: compositeNode***REMOVED***
			children: children,
		***REMOVED***,
		MapType:   mapType,
		Name:      name,
		Equals:    equals,
		Tag:       tag,
		Options:   opts,
		Semicolon: semicolon,
	***REMOVED***
***REMOVED***

func (n *MapFieldNode) FieldLabel() Node ***REMOVED***
	return nil
***REMOVED***

func (n *MapFieldNode) FieldName() Node ***REMOVED***
	return n.Name
***REMOVED***

func (n *MapFieldNode) FieldType() Node ***REMOVED***
	return n.MapType
***REMOVED***

func (n *MapFieldNode) FieldTag() Node ***REMOVED***
	return n.Tag
***REMOVED***

func (n *MapFieldNode) FieldExtendee() Node ***REMOVED***
	return nil
***REMOVED***

func (n *MapFieldNode) GetGroupKeyword() Node ***REMOVED***
	return nil
***REMOVED***

func (n *MapFieldNode) GetOptions() *CompactOptionsNode ***REMOVED***
	return n.Options
***REMOVED***

func (n *MapFieldNode) MessageName() Node ***REMOVED***
	return n.Name
***REMOVED***

func (n *MapFieldNode) KeyField() *SyntheticMapField ***REMOVED***
	return NewSyntheticMapField(n.MapType.KeyType, 1)
***REMOVED***

func (n *MapFieldNode) ValueField() *SyntheticMapField ***REMOVED***
	return NewSyntheticMapField(n.MapType.ValueType, 2)
***REMOVED***

// SyntheticMapField is not an actual node in the AST but a synthetic node
// that implements FieldDeclNode. These are used to represent the implicit
// field declarations of the "key" and "value" fields in a map entry.
type SyntheticMapField struct ***REMOVED***
	Ident IdentValueNode
	Tag   *UintLiteralNode
***REMOVED***

// NewSyntheticMapField creates a new *SyntheticMapField for the given
// identifier (either a key or value type in a map declaration) and tag
// number (1 for key, 2 for value).
func NewSyntheticMapField(ident IdentValueNode, tagNum uint64) *SyntheticMapField ***REMOVED***
	tag := &UintLiteralNode***REMOVED***
		terminalNode: terminalNode***REMOVED***
			posRange: PosRange***REMOVED***Start: *ident.Start(), End: *ident.End()***REMOVED***,
		***REMOVED***,
		Val: tagNum,
	***REMOVED***
	return &SyntheticMapField***REMOVED***Ident: ident, Tag: tag***REMOVED***
***REMOVED***

func (n *SyntheticMapField) Start() *SourcePos ***REMOVED***
	return n.Ident.Start()
***REMOVED***

func (n *SyntheticMapField) End() *SourcePos ***REMOVED***
	return n.Ident.End()
***REMOVED***

func (n *SyntheticMapField) LeadingComments() []Comment ***REMOVED***
	return nil
***REMOVED***

func (n *SyntheticMapField) TrailingComments() []Comment ***REMOVED***
	return nil
***REMOVED***

func (n *SyntheticMapField) FieldLabel() Node ***REMOVED***
	return n.Ident
***REMOVED***

func (n *SyntheticMapField) FieldName() Node ***REMOVED***
	return n.Ident
***REMOVED***

func (n *SyntheticMapField) FieldType() Node ***REMOVED***
	return n.Ident
***REMOVED***

func (n *SyntheticMapField) FieldTag() Node ***REMOVED***
	return n.Tag
***REMOVED***

func (n *SyntheticMapField) FieldExtendee() Node ***REMOVED***
	return nil
***REMOVED***

func (n *SyntheticMapField) GetGroupKeyword() Node ***REMOVED***
	return nil
***REMOVED***

func (n *SyntheticMapField) GetOptions() *CompactOptionsNode ***REMOVED***
	return nil
***REMOVED***
