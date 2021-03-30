package protoparse

import "github.com/jhump/protoreflect/desc/protoparse/ast"

// SourcePos is the same as ast.SourcePos. This alias exists for
// backwards compatibility (SourcePos used to be defined in this package.)
type SourcePos = ast.SourcePos

// the types below are accumulator types: linked lists that are
// constructed during parsing and then converted to slices of AST nodes
// once the whole list has been parsed

type compactOptionList struct ***REMOVED***
	option *ast.OptionNode
	comma  *ast.RuneNode
	next   *compactOptionList
***REMOVED***

func (list *compactOptionList) toNodes() ([]*ast.OptionNode, []*ast.RuneNode) ***REMOVED***
	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	opts := make([]*ast.OptionNode, l)
	commas := make([]*ast.RuneNode, l-1)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		opts[i] = cur.option
		if cur.comma != nil ***REMOVED***
			commas[i] = cur.comma
		***REMOVED***
	***REMOVED***
	return opts, commas
***REMOVED***

type stringList struct ***REMOVED***
	str  *ast.StringLiteralNode
	next *stringList
***REMOVED***

func (list *stringList) toStringValueNode() ast.StringValueNode ***REMOVED***
	if list.next == nil ***REMOVED***
		// single name
		return list.str
	***REMOVED***

	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	strs := make([]*ast.StringLiteralNode, l)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		strs[i] = cur.str
	***REMOVED***
	return ast.NewCompoundLiteralStringNode(strs...)
***REMOVED***

type nameList struct ***REMOVED***
	name  ast.StringValueNode
	comma *ast.RuneNode
	next  *nameList
***REMOVED***

func (list *nameList) toNodes() ([]ast.StringValueNode, []*ast.RuneNode) ***REMOVED***
	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	names := make([]ast.StringValueNode, l)
	commas := make([]*ast.RuneNode, l-1)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		names[i] = cur.name
		if cur.comma != nil ***REMOVED***
			commas[i] = cur.comma
		***REMOVED***
	***REMOVED***
	return names, commas
***REMOVED***

type rangeList struct ***REMOVED***
	rng   *ast.RangeNode
	comma *ast.RuneNode
	next  *rangeList
***REMOVED***

func (list *rangeList) toNodes() ([]*ast.RangeNode, []*ast.RuneNode) ***REMOVED***
	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	ranges := make([]*ast.RangeNode, l)
	commas := make([]*ast.RuneNode, l-1)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		ranges[i] = cur.rng
		if cur.comma != nil ***REMOVED***
			commas[i] = cur.comma
		***REMOVED***
	***REMOVED***
	return ranges, commas
***REMOVED***

type valueList struct ***REMOVED***
	val   ast.ValueNode
	comma *ast.RuneNode
	next  *valueList
***REMOVED***

func (list *valueList) toNodes() ([]ast.ValueNode, []*ast.RuneNode) ***REMOVED***
	if list == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	vals := make([]ast.ValueNode, l)
	commas := make([]*ast.RuneNode, l-1)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		vals[i] = cur.val
		if cur.comma != nil ***REMOVED***
			commas[i] = cur.comma
		***REMOVED***
	***REMOVED***
	return vals, commas
***REMOVED***

type fieldRefList struct ***REMOVED***
	ref  *ast.FieldReferenceNode
	dot  *ast.RuneNode
	next *fieldRefList
***REMOVED***

func (list *fieldRefList) toNodes() ([]*ast.FieldReferenceNode, []*ast.RuneNode) ***REMOVED***
	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	refs := make([]*ast.FieldReferenceNode, l)
	dots := make([]*ast.RuneNode, l-1)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		refs[i] = cur.ref
		if cur.dot != nil ***REMOVED***
			dots[i] = cur.dot
		***REMOVED***
	***REMOVED***

	return refs, dots
***REMOVED***

type identList struct ***REMOVED***
	ident *ast.IdentNode
	dot   *ast.RuneNode
	next  *identList
***REMOVED***

func (list *identList) toIdentValueNode(leadingDot *ast.RuneNode) ast.IdentValueNode ***REMOVED***
	if list.next == nil && leadingDot == nil ***REMOVED***
		// single name
		return list.ident
	***REMOVED***

	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	idents := make([]*ast.IdentNode, l)
	dots := make([]*ast.RuneNode, l-1)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		idents[i] = cur.ident
		if cur.dot != nil ***REMOVED***
			dots[i] = cur.dot
		***REMOVED***
	***REMOVED***

	return ast.NewCompoundIdentNode(leadingDot, idents, dots)
***REMOVED***

type messageFieldEntry struct ***REMOVED***
	field     *ast.MessageFieldNode
	delimiter *ast.RuneNode
***REMOVED***

type messageFieldList struct ***REMOVED***
	field *messageFieldEntry
	next  *messageFieldList
***REMOVED***

func (list *messageFieldList) toNodes() ([]*ast.MessageFieldNode, []*ast.RuneNode) ***REMOVED***
	if list == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	l := 0
	for cur := list; cur != nil; cur = cur.next ***REMOVED***
		l++
	***REMOVED***
	fields := make([]*ast.MessageFieldNode, l)
	delimiters := make([]*ast.RuneNode, l)
	for cur, i := list, 0; cur != nil; cur, i = cur.next, i+1 ***REMOVED***
		fields[i] = cur.field.field
		if cur.field.delimiter != nil ***REMOVED***
			delimiters[i] = cur.field.delimiter
		***REMOVED***
	***REMOVED***

	return fields, delimiters
***REMOVED***
