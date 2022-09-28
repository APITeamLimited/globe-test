package ast

// VisitFunc is used to examine a node in the AST when walking the tree.
// It returns true or false as to whether or not the descendants of the
// given node should be visited. If it returns true, the node's children
// will be visisted; if false, they will not. When returning true, it
// can also return a new VisitFunc to use for the children. If it returns
// (true, nil), then the current function will be re-used when visiting
// the children.
//
// See also the Visitor type.
type VisitFunc func(Node) (bool, VisitFunc)

// Walk conducts a walk of the AST rooted at the given root using the
// given function. It performs a "pre-order traversal", visiting a
// given AST node before it visits that node's descendants.
func Walk(root Node, v VisitFunc) ***REMOVED***
	ok, next := v(root)
	if !ok ***REMOVED***
		return
	***REMOVED***
	if next != nil ***REMOVED***
		v = next
	***REMOVED***
	if comp, ok := root.(CompositeNode); ok ***REMOVED***
		for _, child := range comp.Children() ***REMOVED***
			Walk(child, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Visitor provides a technique for walking the AST that allows for
// dynamic dispatch, where a particular function is invoked based on
// the runtime type of the argument.
//
// It consists of a number of functions, each of which matches a
// concrete Node type. It also includes functions for sub-interfaces
// of Node and the Node interface itself, to be used as broader
// "catch all" functions.
//
// To use a visitor, provide a function for the node types of
// interest and pass visitor.Visit as the function to a Walk operation.
// When a node is traversed, the corresponding function field of
// the visitor is invoked, if not nil. If the function for a node's
// concrete type is nil/absent but the function for an interface it
// implements is present, that interface visit function will be used
// instead. If no matching function is present, the traversal will
// continue. If a matching function is present, it will be invoked
// and its response determines how the traversal proceeds.
//
// Every visit function returns (bool, *Visitor). If the bool returned
// is false, the visited node's descendants are skipped. Otherwise,
// traversal will continue into the node's children. If the returned
// visitor is nil, the current visitor will continue to be used. But
// if a non-nil visitor is returned, it will be used to visit the
// node's children.
type Visitor struct ***REMOVED***
	// VisitFileNode is invoked when visiting a *FileNode in the AST.
	VisitFileNode func(*FileNode) (bool, *Visitor)
	// VisitSyntaxNode is invoked when visiting a *SyntaxNode in the AST.
	VisitSyntaxNode func(*SyntaxNode) (bool, *Visitor)
	// VisitPackageNode is invoked when visiting a *PackageNode in the AST.
	VisitPackageNode func(*PackageNode) (bool, *Visitor)
	// VisitImportNode is invoked when visiting an *ImportNode in the AST.
	VisitImportNode func(*ImportNode) (bool, *Visitor)
	// VisitOptionNode is invoked when visiting an *OptionNode in the AST.
	VisitOptionNode func(*OptionNode) (bool, *Visitor)
	// VisitOptionNameNode is invoked when visiting an *OptionNameNode in the AST.
	VisitOptionNameNode func(*OptionNameNode) (bool, *Visitor)
	// VisitFieldReferenceNode is invoked when visiting a *FieldReferenceNode in the AST.
	VisitFieldReferenceNode func(*FieldReferenceNode) (bool, *Visitor)
	// VisitCompactOptionsNode is invoked when visiting a *CompactOptionsNode in the AST.
	VisitCompactOptionsNode func(*CompactOptionsNode) (bool, *Visitor)
	// VisitMessageNode is invoked when visiting a *MessageNode in the AST.
	VisitMessageNode func(*MessageNode) (bool, *Visitor)
	// VisitExtendNode is invoked when visiting an *ExtendNode in the AST.
	VisitExtendNode func(*ExtendNode) (bool, *Visitor)
	// VisitExtensionRangeNode is invoked when visiting an *ExtensionRangeNode in the AST.
	VisitExtensionRangeNode func(*ExtensionRangeNode) (bool, *Visitor)
	// VisitReservedNode is invoked when visiting a *ReservedNode in the AST.
	VisitReservedNode func(*ReservedNode) (bool, *Visitor)
	// VisitRangeNode is invoked when visiting a *RangeNode in the AST.
	VisitRangeNode func(*RangeNode) (bool, *Visitor)
	// VisitFieldNode is invoked when visiting a *FieldNode in the AST.
	VisitFieldNode func(*FieldNode) (bool, *Visitor)
	// VisitGroupNode is invoked when visiting a *GroupNode in the AST.
	VisitGroupNode func(*GroupNode) (bool, *Visitor)
	// VisitMapFieldNode is invoked when visiting a *MapFieldNode in the AST.
	VisitMapFieldNode func(*MapFieldNode) (bool, *Visitor)
	// VisitMapTypeNode is invoked when visiting a *MapTypeNode in the AST.
	VisitMapTypeNode func(*MapTypeNode) (bool, *Visitor)
	// VisitOneOfNode is invoked when visiting a *OneOfNode in the AST.
	VisitOneOfNode func(*OneOfNode) (bool, *Visitor)
	// VisitEnumNode is invoked when visiting an *EnumNode in the AST.
	VisitEnumNode func(*EnumNode) (bool, *Visitor)
	// VisitEnumValueNode is invoked when visiting an *EnumValueNode in the AST.
	VisitEnumValueNode func(*EnumValueNode) (bool, *Visitor)
	// VisitServiceNode is invoked when visiting a *ServiceNode in the AST.
	VisitServiceNode func(*ServiceNode) (bool, *Visitor)
	// VisitRPCNode is invoked when visiting an *RPCNode in the AST.
	VisitRPCNode func(*RPCNode) (bool, *Visitor)
	// VisitRPCTypeNode is invoked when visiting an *RPCTypeNode in the AST.
	VisitRPCTypeNode func(*RPCTypeNode) (bool, *Visitor)
	// VisitIdentNode is invoked when visiting an *IdentNode in the AST.
	VisitIdentNode func(*IdentNode) (bool, *Visitor)
	// VisitCompoundIdentNode is invoked when visiting a *CompoundIdentNode in the AST.
	VisitCompoundIdentNode func(*CompoundIdentNode) (bool, *Visitor)
	// VisitStringLiteralNode is invoked when visiting a *StringLiteralNode in the AST.
	VisitStringLiteralNode func(*StringLiteralNode) (bool, *Visitor)
	// VisitCompoundStringLiteralNode is invoked when visiting a *CompoundStringLiteralNode in the AST.
	VisitCompoundStringLiteralNode func(*CompoundStringLiteralNode) (bool, *Visitor)
	// VisitUintLiteralNode is invoked when visiting a *UintLiteralNode in the AST.
	VisitUintLiteralNode func(*UintLiteralNode) (bool, *Visitor)
	// VisitPositiveUintLiteralNode is invoked when visiting a *PositiveUintLiteralNode in the AST.
	VisitPositiveUintLiteralNode func(*PositiveUintLiteralNode) (bool, *Visitor)
	// VisitNegativeIntLiteralNode is invoked when visiting a *NegativeIntLiteralNode in the AST.
	VisitNegativeIntLiteralNode func(*NegativeIntLiteralNode) (bool, *Visitor)
	// VisitFloatLiteralNode is invoked when visiting a *FloatLiteralNode in the AST.
	VisitFloatLiteralNode func(*FloatLiteralNode) (bool, *Visitor)
	// VisitSpecialFloatLiteralNode is invoked when visiting a *SpecialFloatLiteralNode in the AST.
	VisitSpecialFloatLiteralNode func(*SpecialFloatLiteralNode) (bool, *Visitor)
	// VisitSignedFloatLiteralNode is invoked when visiting a *SignedFloatLiteralNode in the AST.
	VisitSignedFloatLiteralNode func(*SignedFloatLiteralNode) (bool, *Visitor)
	// VisitBoolLiteralNode is invoked when visiting a *BoolLiteralNode in the AST.
	VisitBoolLiteralNode func(*BoolLiteralNode) (bool, *Visitor)
	// VisitArrayLiteralNode is invoked when visiting an *ArrayLiteralNode in the AST.
	VisitArrayLiteralNode func(*ArrayLiteralNode) (bool, *Visitor)
	// VisitMessageLiteralNode is invoked when visiting a *MessageLiteralNode in the AST.
	VisitMessageLiteralNode func(*MessageLiteralNode) (bool, *Visitor)
	// VisitMessageFieldNode is invoked when visiting a *MessageFieldNode in the AST.
	VisitMessageFieldNode func(*MessageFieldNode) (bool, *Visitor)
	// VisitKeywordNode is invoked when visiting a *KeywordNode in the AST.
	VisitKeywordNode func(*KeywordNode) (bool, *Visitor)
	// VisitRuneNode is invoked when visiting a *RuneNode in the AST.
	VisitRuneNode func(*RuneNode) (bool, *Visitor)
	// VisitEmptyDeclNode is invoked when visiting a *EmptyDeclNode in the AST.
	VisitEmptyDeclNode func(*EmptyDeclNode) (bool, *Visitor)

	// VisitFieldDeclNode is invoked when visiting a FieldDeclNode in the AST.
	// This function is used when no concrete type function is provided. If
	// both this and VisitMessageDeclNode are provided, and a node implements
	// both (such as *GroupNode and *MapFieldNode), this function will be
	// invoked and not the other.
	VisitFieldDeclNode func(FieldDeclNode) (bool, *Visitor)
	// VisitMessageDeclNode is invoked when visiting a MessageDeclNode in the AST.
	// This function is used when no concrete type function is provided.
	VisitMessageDeclNode func(MessageDeclNode) (bool, *Visitor)

	// VisitIdentValueNode is invoked when visiting an IdentValueNode in the AST.
	// This function is used when no concrete type function is provided.
	VisitIdentValueNode func(IdentValueNode) (bool, *Visitor)
	// VisitStringValueNode is invoked when visiting a StringValueNode in the AST.
	// This function is used when no concrete type function is provided.
	VisitStringValueNode func(StringValueNode) (bool, *Visitor)
	// VisitIntValueNode is invoked when visiting an IntValueNode in the AST.
	// This function is used when no concrete type function is provided. If
	// both this and VisitFloatValueNode are provided, and a node implements
	// both (such as *UintLiteralNode), this function will be invoked and
	// not the other.
	VisitIntValueNode func(IntValueNode) (bool, *Visitor)
	// VisitFloatValueNode is invoked when visiting a FloatValueNode in the AST.
	// This function is used when no concrete type function is provided.
	VisitFloatValueNode func(FloatValueNode) (bool, *Visitor)
	// VisitValueNode is invoked when visiting a ValueNode in the AST. This
	// function is used when no concrete type function is provided and no
	// more specific ValueNode function is provided that matches the node.
	VisitValueNode func(ValueNode) (bool, *Visitor)

	// VisitTerminalNode is invoked when visiting a TerminalNode in the AST.
	// This function is used when no concrete type function is provided
	// no more specific interface type function is provided.
	VisitTerminalNode func(TerminalNode) (bool, *Visitor)
	// VisitCompositeNode is invoked when visiting a CompositeNode in the AST.
	// This function is used when no concrete type function is provided
	// no more specific interface type function is provided.
	VisitCompositeNode func(CompositeNode) (bool, *Visitor)
	// VisitNode is invoked when visiting a Node in the AST. This
	// function is only used when no other more specific function is
	// provided.
	VisitNode func(Node) (bool, *Visitor)
***REMOVED***

// Visit provides the Visitor's implementation of VisitFunc, to be
// used with Walk operations.
func (v *Visitor) Visit(n Node) (bool, VisitFunc) ***REMOVED***
	var ok, matched bool
	var next *Visitor
	switch n := n.(type) ***REMOVED***
	case *FileNode:
		if v.VisitFileNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitFileNode(n)
		***REMOVED***
	case *SyntaxNode:
		if v.VisitSyntaxNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitSyntaxNode(n)
		***REMOVED***
	case *PackageNode:
		if v.VisitPackageNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitPackageNode(n)
		***REMOVED***
	case *ImportNode:
		if v.VisitImportNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitImportNode(n)
		***REMOVED***
	case *OptionNode:
		if v.VisitOptionNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitOptionNode(n)
		***REMOVED***
	case *OptionNameNode:
		if v.VisitOptionNameNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitOptionNameNode(n)
		***REMOVED***
	case *FieldReferenceNode:
		if v.VisitFieldReferenceNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitFieldReferenceNode(n)
		***REMOVED***
	case *CompactOptionsNode:
		if v.VisitCompactOptionsNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitCompactOptionsNode(n)
		***REMOVED***
	case *MessageNode:
		if v.VisitMessageNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitMessageNode(n)
		***REMOVED***
	case *ExtendNode:
		if v.VisitExtendNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitExtendNode(n)
		***REMOVED***
	case *ExtensionRangeNode:
		if v.VisitExtensionRangeNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitExtensionRangeNode(n)
		***REMOVED***
	case *ReservedNode:
		if v.VisitReservedNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitReservedNode(n)
		***REMOVED***
	case *RangeNode:
		if v.VisitRangeNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitRangeNode(n)
		***REMOVED***
	case *FieldNode:
		if v.VisitFieldNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitFieldNode(n)
		***REMOVED***
	case *GroupNode:
		if v.VisitGroupNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitGroupNode(n)
		***REMOVED***
	case *MapFieldNode:
		if v.VisitMapFieldNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitMapFieldNode(n)
		***REMOVED***
	case *MapTypeNode:
		if v.VisitMapTypeNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitMapTypeNode(n)
		***REMOVED***
	case *OneOfNode:
		if v.VisitOneOfNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitOneOfNode(n)
		***REMOVED***
	case *EnumNode:
		if v.VisitEnumNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitEnumNode(n)
		***REMOVED***
	case *EnumValueNode:
		if v.VisitEnumValueNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitEnumValueNode(n)
		***REMOVED***
	case *ServiceNode:
		if v.VisitServiceNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitServiceNode(n)
		***REMOVED***
	case *RPCNode:
		if v.VisitRPCNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitRPCNode(n)
		***REMOVED***
	case *RPCTypeNode:
		if v.VisitRPCTypeNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitRPCTypeNode(n)
		***REMOVED***
	case *IdentNode:
		if v.VisitIdentNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitIdentNode(n)
		***REMOVED***
	case *CompoundIdentNode:
		if v.VisitCompoundIdentNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitCompoundIdentNode(n)
		***REMOVED***
	case *StringLiteralNode:
		if v.VisitStringLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitStringLiteralNode(n)
		***REMOVED***
	case *CompoundStringLiteralNode:
		if v.VisitCompoundStringLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitCompoundStringLiteralNode(n)
		***REMOVED***
	case *UintLiteralNode:
		if v.VisitUintLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitUintLiteralNode(n)
		***REMOVED***
	case *PositiveUintLiteralNode:
		if v.VisitPositiveUintLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitPositiveUintLiteralNode(n)
		***REMOVED***
	case *NegativeIntLiteralNode:
		if v.VisitNegativeIntLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitNegativeIntLiteralNode(n)
		***REMOVED***
	case *FloatLiteralNode:
		if v.VisitFloatLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitFloatLiteralNode(n)
		***REMOVED***
	case *SpecialFloatLiteralNode:
		if v.VisitSpecialFloatLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitSpecialFloatLiteralNode(n)
		***REMOVED***
	case *SignedFloatLiteralNode:
		if v.VisitSignedFloatLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitSignedFloatLiteralNode(n)
		***REMOVED***
	case *BoolLiteralNode:
		if v.VisitBoolLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitBoolLiteralNode(n)
		***REMOVED***
	case *ArrayLiteralNode:
		if v.VisitArrayLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitArrayLiteralNode(n)
		***REMOVED***
	case *MessageLiteralNode:
		if v.VisitMessageLiteralNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitMessageLiteralNode(n)
		***REMOVED***
	case *MessageFieldNode:
		if v.VisitMessageFieldNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitMessageFieldNode(n)
		***REMOVED***
	case *KeywordNode:
		if v.VisitKeywordNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitKeywordNode(n)
		***REMOVED***
	case *RuneNode:
		if v.VisitRuneNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitRuneNode(n)
		***REMOVED***
	case *EmptyDeclNode:
		if v.VisitEmptyDeclNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitEmptyDeclNode(n)
		***REMOVED***
	***REMOVED***

	if !matched ***REMOVED***
		// Visitor provided no concrete type visit function, so
		// check interface types. We do this in several passes
		// to provide "priority" for matched interfaces for nodes
		// that actually implement more than one interface.
		//
		// For example, StringLiteralNode implements both
		// StringValueNode and ValueNode. Both cases could match
		// so the first case is what would match. So if we want
		// to test against either, they need to be in different
		// switch statements.
		switch n := n.(type) ***REMOVED***
		case FieldDeclNode:
			if v.VisitFieldDeclNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitFieldDeclNode(n)
			***REMOVED***
		case IdentValueNode:
			if v.VisitIdentValueNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitIdentValueNode(n)
			***REMOVED***
		case StringValueNode:
			if v.VisitStringValueNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitStringValueNode(n)
			***REMOVED***
		case IntValueNode:
			if v.VisitIntValueNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitIntValueNode(n)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !matched ***REMOVED***
		// These two are excluded from the above switch so that
		// if visitor provides both VisitIntValueNode and
		// VisitFloatValueNode, we'll prefer VisitIntValueNode
		// for *UintLiteralNode (which implements both). Similarly,
		// that way we prefer VisitFieldDeclNode over
		// VisitMessageDeclNode when visiting a *GroupNode.
		switch n := n.(type) ***REMOVED***
		case FloatValueNode:
			if v.VisitFloatValueNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitFloatValueNode(n)
			***REMOVED***
		case MessageDeclNode:
			if v.VisitMessageDeclNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitMessageDeclNode(n)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !matched ***REMOVED***
		switch n := n.(type) ***REMOVED***
		case ValueNode:
			if v.VisitValueNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitValueNode(n)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !matched ***REMOVED***
		switch n := n.(type) ***REMOVED***
		case TerminalNode:
			if v.VisitTerminalNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitTerminalNode(n)
			***REMOVED***
		case CompositeNode:
			if v.VisitCompositeNode != nil ***REMOVED***
				matched = true
				ok, next = v.VisitCompositeNode(n)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !matched ***REMOVED***
		// finally, fallback to most generic visit function
		if v.VisitNode != nil ***REMOVED***
			matched = true
			ok, next = v.VisitNode(n)
		***REMOVED***
	***REMOVED***

	if !matched ***REMOVED***
		// keep descending with the current visitor
		return true, nil
	***REMOVED***

	if !ok ***REMOVED***
		return false, nil
	***REMOVED***
	if next != nil ***REMOVED***
		return true, next.Visit
	***REMOVED***
	return true, v.Visit
***REMOVED***
