package protoparse

import "fmt"

// This file defines all of the nodes in the proto AST.

// SourcePos identifies a location in a proto source file.
type SourcePos struct ***REMOVED***
	Filename  string
	Line, Col int
	Offset    int
***REMOVED***

func (pos SourcePos) String() string ***REMOVED***
	if pos.Line <= 0 || pos.Col <= 0 ***REMOVED***
		return pos.Filename
	***REMOVED***
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Col)
***REMOVED***

func unknownPos(filename string) *SourcePos ***REMOVED***
	return &SourcePos***REMOVED***Filename: filename***REMOVED***
***REMOVED***

// node is the interface implemented by all nodes in the AST
type node interface ***REMOVED***
	start() *SourcePos
	end() *SourcePos
	leadingComments() []comment
	trailingComments() []comment
***REMOVED***

type terminalNode interface ***REMOVED***
	node
	popLeadingComment() comment
	pushTrailingComment(comment)
***REMOVED***

var _ terminalNode = (*basicNode)(nil)
var _ terminalNode = (*stringLiteralNode)(nil)
var _ terminalNode = (*intLiteralNode)(nil)
var _ terminalNode = (*floatLiteralNode)(nil)
var _ terminalNode = (*identNode)(nil)

type fileDecl interface ***REMOVED***
	node
	getSyntax() node
***REMOVED***

var _ fileDecl = (*fileNode)(nil)
var _ fileDecl = (*noSourceNode)(nil)

type optionDecl interface ***REMOVED***
	node
	getName() node
	getValue() valueNode
***REMOVED***

var _ optionDecl = (*optionNode)(nil)
var _ optionDecl = (*noSourceNode)(nil)

type fieldDecl interface ***REMOVED***
	node
	fieldLabel() node
	fieldName() node
	fieldType() node
	fieldTag() node
	fieldExtendee() node
	getGroupKeyword() node
***REMOVED***

var _ fieldDecl = (*fieldNode)(nil)
var _ fieldDecl = (*groupNode)(nil)
var _ fieldDecl = (*mapFieldNode)(nil)
var _ fieldDecl = (*syntheticMapField)(nil)
var _ fieldDecl = (*noSourceNode)(nil)

type rangeDecl interface ***REMOVED***
	node
	rangeStart() node
	rangeEnd() node
***REMOVED***

var _ rangeDecl = (*rangeNode)(nil)
var _ rangeDecl = (*noSourceNode)(nil)

type enumValueDecl interface ***REMOVED***
	node
	getName() node
	getNumber() node
***REMOVED***

var _ enumValueDecl = (*enumValueNode)(nil)
var _ enumValueDecl = (*noSourceNode)(nil)

type msgDecl interface ***REMOVED***
	node
	messageName() node
***REMOVED***

var _ msgDecl = (*messageNode)(nil)
var _ msgDecl = (*groupNode)(nil)
var _ msgDecl = (*mapFieldNode)(nil)
var _ msgDecl = (*noSourceNode)(nil)

type methodDecl interface ***REMOVED***
	node
	getInputType() node
	getOutputType() node
***REMOVED***

var _ methodDecl = (*methodNode)(nil)
var _ methodDecl = (*noSourceNode)(nil)

type posRange struct ***REMOVED***
	start, end SourcePos
***REMOVED***

type basicNode struct ***REMOVED***
	posRange
	leading  []comment
	trailing []comment
***REMOVED***

func (n *basicNode) start() *SourcePos ***REMOVED***
	return &n.posRange.start
***REMOVED***

func (n *basicNode) end() *SourcePos ***REMOVED***
	return &n.posRange.end
***REMOVED***

func (n *basicNode) leadingComments() []comment ***REMOVED***
	return n.leading
***REMOVED***

func (n *basicNode) trailingComments() []comment ***REMOVED***
	return n.trailing
***REMOVED***

func (n *basicNode) popLeadingComment() comment ***REMOVED***
	c := n.leading[0]
	n.leading = n.leading[1:]
	return c
***REMOVED***

func (n *basicNode) pushTrailingComment(c comment) ***REMOVED***
	n.trailing = append(n.trailing, c)
***REMOVED***

type comment struct ***REMOVED***
	posRange
	text string
***REMOVED***

type basicCompositeNode struct ***REMOVED***
	first node
	last  node
***REMOVED***

func (n *basicCompositeNode) start() *SourcePos ***REMOVED***
	return n.first.start()
***REMOVED***

func (n *basicCompositeNode) end() *SourcePos ***REMOVED***
	return n.last.end()
***REMOVED***

func (n *basicCompositeNode) leadingComments() []comment ***REMOVED***
	return n.first.leadingComments()
***REMOVED***

func (n *basicCompositeNode) trailingComments() []comment ***REMOVED***
	return n.last.trailingComments()
***REMOVED***

func (n *basicCompositeNode) setRange(first, last node) ***REMOVED***
	n.first = first
	n.last = last
***REMOVED***

type fileNode struct ***REMOVED***
	basicCompositeNode
	syntax *syntaxNode
	decls  []*fileElement

	// This field is populated after parsing, to make it easier to find
	// source locations by import name for constructing link errors.
	imports []*importNode
***REMOVED***

func (n *fileNode) getSyntax() node ***REMOVED***
	return n.syntax
***REMOVED***

type fileElement struct ***REMOVED***
	// a discriminated union: only one field will be set
	imp     *importNode
	pkg     *packageNode
	option  *optionNode
	message *messageNode
	enum    *enumNode
	extend  *extendNode
	service *serviceNode
	empty   *basicNode
***REMOVED***

func (n *fileElement) start() *SourcePos ***REMOVED***
	return n.get().start()
***REMOVED***

func (n *fileElement) end() *SourcePos ***REMOVED***
	return n.get().end()
***REMOVED***

func (n *fileElement) leadingComments() []comment ***REMOVED***
	return n.get().leadingComments()
***REMOVED***

func (n *fileElement) trailingComments() []comment ***REMOVED***
	return n.get().trailingComments()
***REMOVED***

func (n *fileElement) get() node ***REMOVED***
	switch ***REMOVED***
	case n.imp != nil:
		return n.imp
	case n.pkg != nil:
		return n.pkg
	case n.option != nil:
		return n.option
	case n.message != nil:
		return n.message
	case n.enum != nil:
		return n.enum
	case n.extend != nil:
		return n.extend
	case n.service != nil:
		return n.service
	default:
		return n.empty
	***REMOVED***
***REMOVED***

type syntaxNode struct ***REMOVED***
	basicCompositeNode
	syntax *compoundStringNode
***REMOVED***

type importNode struct ***REMOVED***
	basicCompositeNode
	name   *compoundStringNode
	public bool
	weak   bool
***REMOVED***

type packageNode struct ***REMOVED***
	basicCompositeNode
	name *compoundIdentNode
***REMOVED***

type identifier string

type identNode struct ***REMOVED***
	basicNode
	val string
***REMOVED***

func (n *identNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return identifier(n.val)
***REMOVED***

type compoundIdentNode struct ***REMOVED***
	basicCompositeNode
	val string
***REMOVED***

func (n *compoundIdentNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return identifier(n.val)
***REMOVED***

type compactOptionsNode struct ***REMOVED***
	basicCompositeNode
	decls []*optionNode
***REMOVED***

func (n *compactOptionsNode) Elements() []*optionNode ***REMOVED***
	if n == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.decls
***REMOVED***

type optionNode struct ***REMOVED***
	basicCompositeNode
	name *optionNameNode
	val  valueNode
***REMOVED***

func (n *optionNode) getName() node ***REMOVED***
	return n.name
***REMOVED***

func (n *optionNode) getValue() valueNode ***REMOVED***
	return n.val
***REMOVED***

type optionNameNode struct ***REMOVED***
	basicCompositeNode
	parts []*optionNamePartNode
***REMOVED***

type optionNamePartNode struct ***REMOVED***
	basicCompositeNode
	text        *compoundIdentNode
	offset      int
	length      int
	isExtension bool
	st, en      *SourcePos
***REMOVED***

func (n *optionNamePartNode) start() *SourcePos ***REMOVED***
	if n.isExtension ***REMOVED***
		return n.basicCompositeNode.start()
	***REMOVED***
	return n.st
***REMOVED***

func (n *optionNamePartNode) end() *SourcePos ***REMOVED***
	if n.isExtension ***REMOVED***
		return n.basicCompositeNode.end()
	***REMOVED***
	return n.en
***REMOVED***

func (n *optionNamePartNode) setRange(first, last node) ***REMOVED***
	n.basicCompositeNode.setRange(first, last)
	if !n.isExtension ***REMOVED***
		st := *first.start()
		st.Col += n.offset
		n.st = &st
		en := st
		en.Col += n.length
		n.en = &en
	***REMOVED***
***REMOVED***

type valueNode interface ***REMOVED***
	node
	value() interface***REMOVED******REMOVED***
***REMOVED***

var _ valueNode = (*identNode)(nil)
var _ valueNode = (*compoundIdentNode)(nil)
var _ valueNode = (*stringLiteralNode)(nil)
var _ valueNode = (*compoundStringNode)(nil)
var _ valueNode = (*intLiteralNode)(nil)
var _ valueNode = (*compoundIntNode)(nil)
var _ valueNode = (*compoundUintNode)(nil)
var _ valueNode = (*floatLiteralNode)(nil)
var _ valueNode = (*compoundFloatNode)(nil)
var _ valueNode = (*boolLiteralNode)(nil)
var _ valueNode = (*sliceLiteralNode)(nil)
var _ valueNode = (*aggregateLiteralNode)(nil)
var _ valueNode = (*noSourceNode)(nil)

type stringLiteralNode struct ***REMOVED***
	basicNode
	val string
***REMOVED***

func (n *stringLiteralNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

type compoundStringNode struct ***REMOVED***
	basicCompositeNode
	val string
***REMOVED***

func (n *compoundStringNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

type intLiteral interface ***REMOVED***
	asInt32(min, max int32) (int32, bool)
	value() interface***REMOVED******REMOVED***
***REMOVED***

type intLiteralNode struct ***REMOVED***
	basicNode
	val uint64
***REMOVED***

var _ intLiteral = (*intLiteralNode)(nil)

func (n *intLiteralNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

func (n *intLiteralNode) asInt32(min, max int32) (int32, bool) ***REMOVED***
	if (min >= 0 && n.val < uint64(min)) || n.val > uint64(max) ***REMOVED***
		return 0, false
	***REMOVED***
	return int32(n.val), true
***REMOVED***

type compoundUintNode struct ***REMOVED***
	basicCompositeNode
	val uint64
***REMOVED***

var _ intLiteral = (*compoundUintNode)(nil)

func (n *compoundUintNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

func (n *compoundUintNode) asInt32(min, max int32) (int32, bool) ***REMOVED***
	if (min >= 0 && n.val < uint64(min)) || n.val > uint64(max) ***REMOVED***
		return 0, false
	***REMOVED***
	return int32(n.val), true
***REMOVED***

type compoundIntNode struct ***REMOVED***
	basicCompositeNode
	val int64
***REMOVED***

var _ intLiteral = (*compoundIntNode)(nil)

func (n *compoundIntNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

func (n *compoundIntNode) asInt32(min, max int32) (int32, bool) ***REMOVED***
	if n.val < int64(min) || n.val > int64(max) ***REMOVED***
		return 0, false
	***REMOVED***
	return int32(n.val), true
***REMOVED***

type floatLiteralNode struct ***REMOVED***
	basicNode
	val float64
***REMOVED***

func (n *floatLiteralNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

type compoundFloatNode struct ***REMOVED***
	basicCompositeNode
	val float64
***REMOVED***

func (n *compoundFloatNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

type boolLiteralNode struct ***REMOVED***
	*identNode
	val bool
***REMOVED***

func (n *boolLiteralNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.val
***REMOVED***

type sliceLiteralNode struct ***REMOVED***
	basicCompositeNode
	elements []valueNode
***REMOVED***

func (n *sliceLiteralNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.elements
***REMOVED***

type aggregateLiteralNode struct ***REMOVED***
	basicCompositeNode
	elements []*aggregateEntryNode
***REMOVED***

func (n *aggregateLiteralNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.elements
***REMOVED***

type aggregateEntryNode struct ***REMOVED***
	basicCompositeNode
	name *aggregateNameNode
	val  valueNode
***REMOVED***

type aggregateNameNode struct ***REMOVED***
	basicCompositeNode
	name        *compoundIdentNode
	isExtension bool
***REMOVED***

func (a *aggregateNameNode) value() string ***REMOVED***
	if a.isExtension ***REMOVED***
		return "[" + a.name.val + "]"
	***REMOVED*** else ***REMOVED***
		return a.name.val
	***REMOVED***
***REMOVED***

type fieldNode struct ***REMOVED***
	basicCompositeNode
	label   fieldLabel
	fldType *compoundIdentNode
	name    *identNode
	tag     *intLiteralNode
	options *compactOptionsNode

	// This field is populated after parsing, to allow lookup of extendee source
	// locations when field extendees cannot be linked. (Otherwise, this is just
	// stored as a string in the field descriptors defined inside the extend
	// block).
	extendee *extendNode
***REMOVED***

func (n *fieldNode) fieldLabel() node ***REMOVED***
	// proto3 fields and fields inside one-ofs will not have a label and we need
	// this check in order to return a nil node -- otherwise we'd return a
	// non-nil node that has a nil pointer value in it :/
	if n.label.identNode == nil ***REMOVED***
		return nil
	***REMOVED***
	return n.label.identNode
***REMOVED***

func (n *fieldNode) fieldName() node ***REMOVED***
	return n.name
***REMOVED***

func (n *fieldNode) fieldType() node ***REMOVED***
	return n.fldType
***REMOVED***

func (n *fieldNode) fieldTag() node ***REMOVED***
	return n.tag
***REMOVED***

func (n *fieldNode) fieldExtendee() node ***REMOVED***
	if n.extendee != nil ***REMOVED***
		return n.extendee.extendee
	***REMOVED***
	return nil
***REMOVED***

func (n *fieldNode) getGroupKeyword() node ***REMOVED***
	return nil
***REMOVED***

type fieldLabel struct ***REMOVED***
	*identNode
	repeated bool
	required bool
***REMOVED***

type groupNode struct ***REMOVED***
	basicCompositeNode
	groupKeyword *identNode
	label        fieldLabel
	name         *identNode
	tag          *intLiteralNode
	decls        []*messageElement
	options      *compactOptionsNode

	// This field is populated after parsing, to allow lookup of extendee source
	// locations when field extendees cannot be linked. (Otherwise, this is just
	// stored as a string in the field descriptors defined inside the extend
	// block).
	extendee *extendNode
***REMOVED***

func (n *groupNode) fieldLabel() node ***REMOVED***
	if n.label.identNode == nil ***REMOVED***
		// return nil interface to indicate absence, not a typed nil
		return nil
	***REMOVED***
	return n.label.identNode
***REMOVED***

func (n *groupNode) fieldName() node ***REMOVED***
	return n.name
***REMOVED***

func (n *groupNode) fieldType() node ***REMOVED***
	return n.groupKeyword
***REMOVED***

func (n *groupNode) fieldTag() node ***REMOVED***
	return n.tag
***REMOVED***

func (n *groupNode) fieldExtendee() node ***REMOVED***
	if n.extendee != nil ***REMOVED***
		return n.extendee.extendee
	***REMOVED***
	return nil
***REMOVED***

func (n *groupNode) getGroupKeyword() node ***REMOVED***
	return n.groupKeyword
***REMOVED***

func (n *groupNode) messageName() node ***REMOVED***
	return n.name
***REMOVED***

type oneOfNode struct ***REMOVED***
	basicCompositeNode
	name  *identNode
	decls []*oneOfElement
***REMOVED***

type oneOfElement struct ***REMOVED***
	// a discriminated union: only one field will be set
	option *optionNode
	field  *fieldNode
	group  *groupNode
	empty  *basicNode
***REMOVED***

func (n *oneOfElement) start() *SourcePos ***REMOVED***
	return n.get().start()
***REMOVED***

func (n *oneOfElement) end() *SourcePos ***REMOVED***
	return n.get().end()
***REMOVED***

func (n *oneOfElement) leadingComments() []comment ***REMOVED***
	return n.get().leadingComments()
***REMOVED***

func (n *oneOfElement) trailingComments() []comment ***REMOVED***
	return n.get().trailingComments()
***REMOVED***

func (n *oneOfElement) get() node ***REMOVED***
	switch ***REMOVED***
	case n.option != nil:
		return n.option
	case n.field != nil:
		return n.field
	default:
		return n.empty
	***REMOVED***
***REMOVED***

type mapTypeNode struct ***REMOVED***
	basicCompositeNode
	mapKeyword *identNode
	keyType    *identNode
	valueType  *compoundIdentNode
***REMOVED***

type mapFieldNode struct ***REMOVED***
	basicCompositeNode
	mapType *mapTypeNode
	name    *identNode
	tag     *intLiteralNode
	options *compactOptionsNode
***REMOVED***

func (n *mapFieldNode) fieldLabel() node ***REMOVED***
	return nil
***REMOVED***

func (n *mapFieldNode) fieldName() node ***REMOVED***
	return n.name
***REMOVED***

func (n *mapFieldNode) fieldType() node ***REMOVED***
	return n.mapType
***REMOVED***

func (n *mapFieldNode) fieldTag() node ***REMOVED***
	return n.tag
***REMOVED***

func (n *mapFieldNode) fieldExtendee() node ***REMOVED***
	return nil
***REMOVED***

func (n *mapFieldNode) getGroupKeyword() node ***REMOVED***
	return nil
***REMOVED***

func (n *mapFieldNode) messageName() node ***REMOVED***
	return n.name
***REMOVED***

func (n *mapFieldNode) keyField() *syntheticMapField ***REMOVED***
	k := n.mapType.keyType
	t := &compoundIdentNode***REMOVED***val: k.val***REMOVED***
	t.setRange(k, k)
	return newSyntheticMapField(t, 1)
***REMOVED***

func (n *mapFieldNode) valueField() *syntheticMapField ***REMOVED***
	return newSyntheticMapField(n.mapType.valueType, 2)
***REMOVED***

func newSyntheticMapField(ident *compoundIdentNode, tagNum uint64) *syntheticMapField ***REMOVED***
	tag := &intLiteralNode***REMOVED***
		basicNode: basicNode***REMOVED***
			posRange: posRange***REMOVED***start: *ident.start(), end: *ident.end()***REMOVED***,
		***REMOVED***,
		val: tagNum,
	***REMOVED***
	return &syntheticMapField***REMOVED***ident: ident, tag: tag***REMOVED***
***REMOVED***

type syntheticMapField struct ***REMOVED***
	ident *compoundIdentNode
	tag   *intLiteralNode
***REMOVED***

func (n *syntheticMapField) start() *SourcePos ***REMOVED***
	return n.ident.start()
***REMOVED***

func (n *syntheticMapField) end() *SourcePos ***REMOVED***
	return n.ident.end()
***REMOVED***

func (n *syntheticMapField) leadingComments() []comment ***REMOVED***
	return nil
***REMOVED***

func (n *syntheticMapField) trailingComments() []comment ***REMOVED***
	return nil
***REMOVED***

func (n *syntheticMapField) fieldLabel() node ***REMOVED***
	return n.ident
***REMOVED***

func (n *syntheticMapField) fieldName() node ***REMOVED***
	return n.ident
***REMOVED***

func (n *syntheticMapField) fieldType() node ***REMOVED***
	return n.ident
***REMOVED***

func (n *syntheticMapField) fieldTag() node ***REMOVED***
	return n.tag
***REMOVED***

func (n *syntheticMapField) fieldExtendee() node ***REMOVED***
	return nil
***REMOVED***

func (n *syntheticMapField) getGroupKeyword() node ***REMOVED***
	return nil
***REMOVED***

type extensionRangeNode struct ***REMOVED***
	basicCompositeNode
	ranges  []*rangeNode
	options *compactOptionsNode
***REMOVED***

type rangeNode struct ***REMOVED***
	basicCompositeNode
	startNode, endNode node
	endMax             bool
***REMOVED***

func (n *rangeNode) rangeStart() node ***REMOVED***
	return n.startNode
***REMOVED***

func (n *rangeNode) rangeEnd() node ***REMOVED***
	if n.endNode == nil ***REMOVED***
		return n.startNode
	***REMOVED***
	return n.endNode
***REMOVED***

func (n *rangeNode) startValue() interface***REMOVED******REMOVED*** ***REMOVED***
	return n.startNode.(intLiteral).value()
***REMOVED***

func (n *rangeNode) startValueAsInt32(min, max int32) (int32, bool) ***REMOVED***
	return n.startNode.(intLiteral).asInt32(min, max)
***REMOVED***

func (n *rangeNode) endValue() interface***REMOVED******REMOVED*** ***REMOVED***
	l, ok := n.endNode.(intLiteral)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return l.value()
***REMOVED***

func (n *rangeNode) endValueAsInt32(min, max int32) (int32, bool) ***REMOVED***
	if n.endMax ***REMOVED***
		return max, true
	***REMOVED***
	if n.endNode == nil ***REMOVED***
		return n.startValueAsInt32(min, max)
	***REMOVED***
	return n.endNode.(intLiteral).asInt32(min, max)
***REMOVED***

type reservedNode struct ***REMOVED***
	basicCompositeNode
	ranges []*rangeNode
	names  []*compoundStringNode
***REMOVED***

type enumNode struct ***REMOVED***
	basicCompositeNode
	name  *identNode
	decls []*enumElement
***REMOVED***

type enumElement struct ***REMOVED***
	// a discriminated union: only one field will be set
	option   *optionNode
	value    *enumValueNode
	reserved *reservedNode
	empty    *basicNode
***REMOVED***

func (n *enumElement) start() *SourcePos ***REMOVED***
	return n.get().start()
***REMOVED***

func (n *enumElement) end() *SourcePos ***REMOVED***
	return n.get().end()
***REMOVED***

func (n *enumElement) leadingComments() []comment ***REMOVED***
	return n.get().leadingComments()
***REMOVED***

func (n *enumElement) trailingComments() []comment ***REMOVED***
	return n.get().trailingComments()
***REMOVED***

func (n *enumElement) get() node ***REMOVED***
	switch ***REMOVED***
	case n.option != nil:
		return n.option
	case n.value != nil:
		return n.value
	default:
		return n.empty
	***REMOVED***
***REMOVED***

type enumValueNode struct ***REMOVED***
	basicCompositeNode
	name    *identNode
	options *compactOptionsNode
	number  *compoundIntNode
***REMOVED***

func (n *enumValueNode) getName() node ***REMOVED***
	return n.name
***REMOVED***

func (n *enumValueNode) getNumber() node ***REMOVED***
	return n.number
***REMOVED***

type messageNode struct ***REMOVED***
	basicCompositeNode
	name  *identNode
	decls []*messageElement
***REMOVED***

func (n *messageNode) messageName() node ***REMOVED***
	return n.name
***REMOVED***

type messageElement struct ***REMOVED***
	// a discriminated union: only one field will be set
	option         *optionNode
	field          *fieldNode
	mapField       *mapFieldNode
	oneOf          *oneOfNode
	group          *groupNode
	nested         *messageNode
	enum           *enumNode
	extend         *extendNode
	extensionRange *extensionRangeNode
	reserved       *reservedNode
	empty          *basicNode
***REMOVED***

func (n *messageElement) start() *SourcePos ***REMOVED***
	return n.get().start()
***REMOVED***

func (n *messageElement) end() *SourcePos ***REMOVED***
	return n.get().end()
***REMOVED***

func (n *messageElement) leadingComments() []comment ***REMOVED***
	return n.get().leadingComments()
***REMOVED***

func (n *messageElement) trailingComments() []comment ***REMOVED***
	return n.get().trailingComments()
***REMOVED***

func (n *messageElement) get() node ***REMOVED***
	switch ***REMOVED***
	case n.option != nil:
		return n.option
	case n.field != nil:
		return n.field
	case n.mapField != nil:
		return n.mapField
	case n.oneOf != nil:
		return n.oneOf
	case n.group != nil:
		return n.group
	case n.nested != nil:
		return n.nested
	case n.enum != nil:
		return n.enum
	case n.extend != nil:
		return n.extend
	case n.extensionRange != nil:
		return n.extensionRange
	case n.reserved != nil:
		return n.reserved
	default:
		return n.empty
	***REMOVED***
***REMOVED***

type extendNode struct ***REMOVED***
	basicCompositeNode
	extendee *compoundIdentNode
	decls    []*extendElement
***REMOVED***

type extendElement struct ***REMOVED***
	// a discriminated union: only one field will be set
	field *fieldNode
	group *groupNode
	empty *basicNode
***REMOVED***

func (n *extendElement) start() *SourcePos ***REMOVED***
	return n.get().start()
***REMOVED***

func (n *extendElement) end() *SourcePos ***REMOVED***
	return n.get().end()
***REMOVED***

func (n *extendElement) leadingComments() []comment ***REMOVED***
	return n.get().leadingComments()
***REMOVED***

func (n *extendElement) trailingComments() []comment ***REMOVED***
	return n.get().trailingComments()
***REMOVED***

func (n *extendElement) get() node ***REMOVED***
	switch ***REMOVED***
	case n.field != nil:
		return n.field
	case n.group != nil:
		return n.group
	default:
		return n.empty
	***REMOVED***
***REMOVED***

type serviceNode struct ***REMOVED***
	basicCompositeNode
	name  *identNode
	decls []*serviceElement
***REMOVED***

type serviceElement struct ***REMOVED***
	// a discriminated union: only one field will be set
	option *optionNode
	rpc    *methodNode
	empty  *basicNode
***REMOVED***

func (n *serviceElement) start() *SourcePos ***REMOVED***
	return n.get().start()
***REMOVED***

func (n *serviceElement) end() *SourcePos ***REMOVED***
	return n.get().end()
***REMOVED***

func (n *serviceElement) leadingComments() []comment ***REMOVED***
	return n.get().leadingComments()
***REMOVED***

func (n *serviceElement) trailingComments() []comment ***REMOVED***
	return n.get().trailingComments()
***REMOVED***

func (n *serviceElement) get() node ***REMOVED***
	switch ***REMOVED***
	case n.option != nil:
		return n.option
	case n.rpc != nil:
		return n.rpc
	default:
		return n.empty
	***REMOVED***
***REMOVED***

type methodNode struct ***REMOVED***
	basicCompositeNode
	name    *identNode
	input   *rpcTypeNode
	output  *rpcTypeNode
	options []*optionNode
***REMOVED***

func (n *methodNode) getInputType() node ***REMOVED***
	return n.input.msgType
***REMOVED***

func (n *methodNode) getOutputType() node ***REMOVED***
	return n.output.msgType
***REMOVED***

type rpcTypeNode struct ***REMOVED***
	basicCompositeNode
	msgType       *compoundIdentNode
	streamKeyword node
***REMOVED***

type noSourceNode struct ***REMOVED***
	pos *SourcePos
***REMOVED***

func (n noSourceNode) start() *SourcePos ***REMOVED***
	return n.pos
***REMOVED***

func (n noSourceNode) end() *SourcePos ***REMOVED***
	return n.pos
***REMOVED***

func (n noSourceNode) leadingComments() []comment ***REMOVED***
	return nil
***REMOVED***

func (n noSourceNode) trailingComments() []comment ***REMOVED***
	return nil
***REMOVED***

func (n noSourceNode) getSyntax() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) getName() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) getValue() valueNode ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) fieldLabel() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) fieldName() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) fieldType() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) fieldTag() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) fieldExtendee() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) getGroupKeyword() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) rangeStart() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) rangeEnd() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) getNumber() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) messageName() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) getInputType() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) getOutputType() node ***REMOVED***
	return n
***REMOVED***

func (n noSourceNode) value() interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***
