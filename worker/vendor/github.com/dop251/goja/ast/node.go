/*
Package ast declares types representing a JavaScript AST.

Warning

The parser and AST interfaces are still works-in-progress (particularly where
node types are concerned) and may change in the future.

*/
package ast

import (
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
)

type PropertyKind string

const (
	PropertyKindValue  PropertyKind = "value"
	PropertyKindGet    PropertyKind = "get"
	PropertyKindSet    PropertyKind = "set"
	PropertyKindMethod PropertyKind = "method"
)

// All nodes implement the Node interface.
type Node interface ***REMOVED***
	Idx0() file.Idx // The index of the first character belonging to the node
	Idx1() file.Idx // The index of the first character immediately after the node
***REMOVED***

// ========== //
// Expression //
// ========== //

type (
	// All expression nodes implement the Expression interface.
	Expression interface ***REMOVED***
		Node
		_expressionNode()
	***REMOVED***

	BindingTarget interface ***REMOVED***
		Expression
		_bindingTarget()
	***REMOVED***

	Binding struct ***REMOVED***
		Target      BindingTarget
		Initializer Expression
	***REMOVED***

	Pattern interface ***REMOVED***
		BindingTarget
		_pattern()
	***REMOVED***

	ArrayLiteral struct ***REMOVED***
		LeftBracket  file.Idx
		RightBracket file.Idx
		Value        []Expression
	***REMOVED***

	ArrayPattern struct ***REMOVED***
		LeftBracket  file.Idx
		RightBracket file.Idx
		Elements     []Expression
		Rest         Expression
	***REMOVED***

	AssignExpression struct ***REMOVED***
		Operator token.Token
		Left     Expression
		Right    Expression
	***REMOVED***

	BadExpression struct ***REMOVED***
		From file.Idx
		To   file.Idx
	***REMOVED***

	BinaryExpression struct ***REMOVED***
		Operator   token.Token
		Left       Expression
		Right      Expression
		Comparison bool
	***REMOVED***

	BooleanLiteral struct ***REMOVED***
		Idx     file.Idx
		Literal string
		Value   bool
	***REMOVED***

	BracketExpression struct ***REMOVED***
		Left         Expression
		Member       Expression
		LeftBracket  file.Idx
		RightBracket file.Idx
	***REMOVED***

	CallExpression struct ***REMOVED***
		Callee           Expression
		LeftParenthesis  file.Idx
		ArgumentList     []Expression
		RightParenthesis file.Idx
	***REMOVED***

	ConditionalExpression struct ***REMOVED***
		Test       Expression
		Consequent Expression
		Alternate  Expression
	***REMOVED***

	DotExpression struct ***REMOVED***
		Left       Expression
		Identifier Identifier
	***REMOVED***

	PrivateDotExpression struct ***REMOVED***
		Left       Expression
		Identifier PrivateIdentifier
	***REMOVED***

	OptionalChain struct ***REMOVED***
		Expression
	***REMOVED***

	Optional struct ***REMOVED***
		Expression
	***REMOVED***

	FunctionLiteral struct ***REMOVED***
		Function      file.Idx
		Name          *Identifier
		ParameterList *ParameterList
		Body          *BlockStatement
		Source        string

		DeclarationList []*VariableDeclaration
	***REMOVED***

	ClassLiteral struct ***REMOVED***
		Class      file.Idx
		RightBrace file.Idx
		Name       *Identifier
		SuperClass Expression
		Body       []ClassElement
		Source     string
	***REMOVED***

	ConciseBody interface ***REMOVED***
		Node
		_conciseBody()
	***REMOVED***

	ExpressionBody struct ***REMOVED***
		Expression Expression
	***REMOVED***

	ArrowFunctionLiteral struct ***REMOVED***
		Start           file.Idx
		ParameterList   *ParameterList
		Body            ConciseBody
		Source          string
		DeclarationList []*VariableDeclaration
	***REMOVED***

	Identifier struct ***REMOVED***
		Name unistring.String
		Idx  file.Idx
	***REMOVED***

	PrivateIdentifier struct ***REMOVED***
		Identifier
	***REMOVED***

	NewExpression struct ***REMOVED***
		New              file.Idx
		Callee           Expression
		LeftParenthesis  file.Idx
		ArgumentList     []Expression
		RightParenthesis file.Idx
	***REMOVED***

	NullLiteral struct ***REMOVED***
		Idx     file.Idx
		Literal string
	***REMOVED***

	NumberLiteral struct ***REMOVED***
		Idx     file.Idx
		Literal string
		Value   interface***REMOVED******REMOVED***
	***REMOVED***

	ObjectLiteral struct ***REMOVED***
		LeftBrace  file.Idx
		RightBrace file.Idx
		Value      []Property
	***REMOVED***

	ObjectPattern struct ***REMOVED***
		LeftBrace  file.Idx
		RightBrace file.Idx
		Properties []Property
		Rest       Expression
	***REMOVED***

	ParameterList struct ***REMOVED***
		Opening file.Idx
		List    []*Binding
		Rest    Expression
		Closing file.Idx
	***REMOVED***

	Property interface ***REMOVED***
		Expression
		_property()
	***REMOVED***

	PropertyShort struct ***REMOVED***
		Name        Identifier
		Initializer Expression
	***REMOVED***

	PropertyKeyed struct ***REMOVED***
		Key      Expression
		Kind     PropertyKind
		Value    Expression
		Computed bool
	***REMOVED***

	SpreadElement struct ***REMOVED***
		Expression
	***REMOVED***

	RegExpLiteral struct ***REMOVED***
		Idx     file.Idx
		Literal string
		Pattern string
		Flags   string
	***REMOVED***

	SequenceExpression struct ***REMOVED***
		Sequence []Expression
	***REMOVED***

	StringLiteral struct ***REMOVED***
		Idx     file.Idx
		Literal string
		Value   unistring.String
	***REMOVED***

	TemplateElement struct ***REMOVED***
		Idx     file.Idx
		Literal string
		Parsed  unistring.String
		Valid   bool
	***REMOVED***

	TemplateLiteral struct ***REMOVED***
		OpenQuote   file.Idx
		CloseQuote  file.Idx
		Tag         Expression
		Elements    []*TemplateElement
		Expressions []Expression
	***REMOVED***

	ThisExpression struct ***REMOVED***
		Idx file.Idx
	***REMOVED***

	SuperExpression struct ***REMOVED***
		Idx file.Idx
	***REMOVED***

	UnaryExpression struct ***REMOVED***
		Operator token.Token
		Idx      file.Idx // If a prefix operation
		Operand  Expression
		Postfix  bool
	***REMOVED***

	MetaProperty struct ***REMOVED***
		Meta, Property *Identifier
		Idx            file.Idx
	***REMOVED***
)

// _expressionNode

func (*ArrayLiteral) _expressionNode()          ***REMOVED******REMOVED***
func (*AssignExpression) _expressionNode()      ***REMOVED******REMOVED***
func (*BadExpression) _expressionNode()         ***REMOVED******REMOVED***
func (*BinaryExpression) _expressionNode()      ***REMOVED******REMOVED***
func (*BooleanLiteral) _expressionNode()        ***REMOVED******REMOVED***
func (*BracketExpression) _expressionNode()     ***REMOVED******REMOVED***
func (*CallExpression) _expressionNode()        ***REMOVED******REMOVED***
func (*ConditionalExpression) _expressionNode() ***REMOVED******REMOVED***
func (*DotExpression) _expressionNode()         ***REMOVED******REMOVED***
func (*PrivateDotExpression) _expressionNode()  ***REMOVED******REMOVED***
func (*FunctionLiteral) _expressionNode()       ***REMOVED******REMOVED***
func (*ClassLiteral) _expressionNode()          ***REMOVED******REMOVED***
func (*ArrowFunctionLiteral) _expressionNode()  ***REMOVED******REMOVED***
func (*Identifier) _expressionNode()            ***REMOVED******REMOVED***
func (*NewExpression) _expressionNode()         ***REMOVED******REMOVED***
func (*NullLiteral) _expressionNode()           ***REMOVED******REMOVED***
func (*NumberLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*ObjectLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*RegExpLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*SequenceExpression) _expressionNode()    ***REMOVED******REMOVED***
func (*StringLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*TemplateLiteral) _expressionNode()       ***REMOVED******REMOVED***
func (*ThisExpression) _expressionNode()        ***REMOVED******REMOVED***
func (*SuperExpression) _expressionNode()       ***REMOVED******REMOVED***
func (*UnaryExpression) _expressionNode()       ***REMOVED******REMOVED***
func (*MetaProperty) _expressionNode()          ***REMOVED******REMOVED***
func (*ObjectPattern) _expressionNode()         ***REMOVED******REMOVED***
func (*ArrayPattern) _expressionNode()          ***REMOVED******REMOVED***
func (*Binding) _expressionNode()               ***REMOVED******REMOVED***

func (*PropertyShort) _expressionNode() ***REMOVED******REMOVED***
func (*PropertyKeyed) _expressionNode() ***REMOVED******REMOVED***

// ========= //
// Statement //
// ========= //

type (
	// All statement nodes implement the Statement interface.
	Statement interface ***REMOVED***
		Node
		_statementNode()
	***REMOVED***

	BadStatement struct ***REMOVED***
		From file.Idx
		To   file.Idx
	***REMOVED***

	BlockStatement struct ***REMOVED***
		LeftBrace  file.Idx
		List       []Statement
		RightBrace file.Idx
	***REMOVED***

	BranchStatement struct ***REMOVED***
		Idx   file.Idx
		Token token.Token
		Label *Identifier
	***REMOVED***

	CaseStatement struct ***REMOVED***
		Case       file.Idx
		Test       Expression
		Consequent []Statement
	***REMOVED***

	CatchStatement struct ***REMOVED***
		Catch     file.Idx
		Parameter BindingTarget
		Body      *BlockStatement
	***REMOVED***

	DebuggerStatement struct ***REMOVED***
		Debugger file.Idx
	***REMOVED***

	DoWhileStatement struct ***REMOVED***
		Do   file.Idx
		Test Expression
		Body Statement
	***REMOVED***

	EmptyStatement struct ***REMOVED***
		Semicolon file.Idx
	***REMOVED***

	ExpressionStatement struct ***REMOVED***
		Expression Expression
	***REMOVED***

	ForInStatement struct ***REMOVED***
		For    file.Idx
		Into   ForInto
		Source Expression
		Body   Statement
	***REMOVED***

	ForOfStatement struct ***REMOVED***
		For    file.Idx
		Into   ForInto
		Source Expression
		Body   Statement
	***REMOVED***

	ForStatement struct ***REMOVED***
		For         file.Idx
		Initializer ForLoopInitializer
		Update      Expression
		Test        Expression
		Body        Statement
	***REMOVED***

	IfStatement struct ***REMOVED***
		If         file.Idx
		Test       Expression
		Consequent Statement
		Alternate  Statement
	***REMOVED***

	LabelledStatement struct ***REMOVED***
		Label     *Identifier
		Colon     file.Idx
		Statement Statement
	***REMOVED***

	ReturnStatement struct ***REMOVED***
		Return   file.Idx
		Argument Expression
	***REMOVED***

	SwitchStatement struct ***REMOVED***
		Switch       file.Idx
		Discriminant Expression
		Default      int
		Body         []*CaseStatement
	***REMOVED***

	ThrowStatement struct ***REMOVED***
		Throw    file.Idx
		Argument Expression
	***REMOVED***

	TryStatement struct ***REMOVED***
		Try     file.Idx
		Body    *BlockStatement
		Catch   *CatchStatement
		Finally *BlockStatement
	***REMOVED***

	VariableStatement struct ***REMOVED***
		Var  file.Idx
		List []*Binding
	***REMOVED***

	LexicalDeclaration struct ***REMOVED***
		Idx   file.Idx
		Token token.Token
		List  []*Binding
	***REMOVED***

	WhileStatement struct ***REMOVED***
		While file.Idx
		Test  Expression
		Body  Statement
	***REMOVED***

	WithStatement struct ***REMOVED***
		With   file.Idx
		Object Expression
		Body   Statement
	***REMOVED***

	FunctionDeclaration struct ***REMOVED***
		Function *FunctionLiteral
	***REMOVED***

	ClassDeclaration struct ***REMOVED***
		Class *ClassLiteral
	***REMOVED***
)

// _statementNode

func (*BadStatement) _statementNode()        ***REMOVED******REMOVED***
func (*BlockStatement) _statementNode()      ***REMOVED******REMOVED***
func (*BranchStatement) _statementNode()     ***REMOVED******REMOVED***
func (*CaseStatement) _statementNode()       ***REMOVED******REMOVED***
func (*CatchStatement) _statementNode()      ***REMOVED******REMOVED***
func (*DebuggerStatement) _statementNode()   ***REMOVED******REMOVED***
func (*DoWhileStatement) _statementNode()    ***REMOVED******REMOVED***
func (*EmptyStatement) _statementNode()      ***REMOVED******REMOVED***
func (*ExpressionStatement) _statementNode() ***REMOVED******REMOVED***
func (*ForInStatement) _statementNode()      ***REMOVED******REMOVED***
func (*ForOfStatement) _statementNode()      ***REMOVED******REMOVED***
func (*ForStatement) _statementNode()        ***REMOVED******REMOVED***
func (*IfStatement) _statementNode()         ***REMOVED******REMOVED***
func (*LabelledStatement) _statementNode()   ***REMOVED******REMOVED***
func (*ReturnStatement) _statementNode()     ***REMOVED******REMOVED***
func (*SwitchStatement) _statementNode()     ***REMOVED******REMOVED***
func (*ThrowStatement) _statementNode()      ***REMOVED******REMOVED***
func (*TryStatement) _statementNode()        ***REMOVED******REMOVED***
func (*VariableStatement) _statementNode()   ***REMOVED******REMOVED***
func (*WhileStatement) _statementNode()      ***REMOVED******REMOVED***
func (*WithStatement) _statementNode()       ***REMOVED******REMOVED***
func (*LexicalDeclaration) _statementNode()  ***REMOVED******REMOVED***
func (*FunctionDeclaration) _statementNode() ***REMOVED******REMOVED***
func (*ClassDeclaration) _statementNode()    ***REMOVED******REMOVED***

// =========== //
// Declaration //
// =========== //

type (
	VariableDeclaration struct ***REMOVED***
		Var  file.Idx
		List []*Binding
	***REMOVED***

	ClassElement interface ***REMOVED***
		Node
		_classElement()
	***REMOVED***

	FieldDefinition struct ***REMOVED***
		Idx         file.Idx
		Key         Expression
		Initializer Expression
		Computed    bool
		Static      bool
	***REMOVED***

	MethodDefinition struct ***REMOVED***
		Idx      file.Idx
		Key      Expression
		Kind     PropertyKind // "method", "get" or "set"
		Body     *FunctionLiteral
		Computed bool
		Static   bool
	***REMOVED***

	ClassStaticBlock struct ***REMOVED***
		Static          file.Idx
		Block           *BlockStatement
		Source          string
		DeclarationList []*VariableDeclaration
	***REMOVED***
)

type (
	ForLoopInitializer interface ***REMOVED***
		_forLoopInitializer()
	***REMOVED***

	ForLoopInitializerExpression struct ***REMOVED***
		Expression Expression
	***REMOVED***

	ForLoopInitializerVarDeclList struct ***REMOVED***
		Var  file.Idx
		List []*Binding
	***REMOVED***

	ForLoopInitializerLexicalDecl struct ***REMOVED***
		LexicalDeclaration LexicalDeclaration
	***REMOVED***

	ForInto interface ***REMOVED***
		Node
		_forInto()
	***REMOVED***

	ForIntoVar struct ***REMOVED***
		Binding *Binding
	***REMOVED***

	ForDeclaration struct ***REMOVED***
		Idx     file.Idx
		IsConst bool
		Target  BindingTarget
	***REMOVED***

	ForIntoExpression struct ***REMOVED***
		Expression Expression
	***REMOVED***
)

func (*ForLoopInitializerExpression) _forLoopInitializer()  ***REMOVED******REMOVED***
func (*ForLoopInitializerVarDeclList) _forLoopInitializer() ***REMOVED******REMOVED***
func (*ForLoopInitializerLexicalDecl) _forLoopInitializer() ***REMOVED******REMOVED***

func (*ForIntoVar) _forInto()        ***REMOVED******REMOVED***
func (*ForDeclaration) _forInto()    ***REMOVED******REMOVED***
func (*ForIntoExpression) _forInto() ***REMOVED******REMOVED***

func (*ArrayPattern) _pattern()       ***REMOVED******REMOVED***
func (*ArrayPattern) _bindingTarget() ***REMOVED******REMOVED***

func (*ObjectPattern) _pattern()       ***REMOVED******REMOVED***
func (*ObjectPattern) _bindingTarget() ***REMOVED******REMOVED***

func (*BadExpression) _bindingTarget() ***REMOVED******REMOVED***

func (*PropertyShort) _property() ***REMOVED******REMOVED***
func (*PropertyKeyed) _property() ***REMOVED******REMOVED***
func (*SpreadElement) _property() ***REMOVED******REMOVED***

func (*Identifier) _bindingTarget() ***REMOVED******REMOVED***

func (*BlockStatement) _conciseBody() ***REMOVED******REMOVED***
func (*ExpressionBody) _conciseBody() ***REMOVED******REMOVED***

func (*FieldDefinition) _classElement()  ***REMOVED******REMOVED***
func (*MethodDefinition) _classElement() ***REMOVED******REMOVED***
func (*ClassStaticBlock) _classElement() ***REMOVED******REMOVED***

// ==== //
// Node //
// ==== //

type Program struct ***REMOVED***
	Body []Statement

	DeclarationList []*VariableDeclaration

	File *file.File
***REMOVED***

// ==== //
// Idx0 //
// ==== //

func (self *ArrayLiteral) Idx0() file.Idx          ***REMOVED*** return self.LeftBracket ***REMOVED***
func (self *ArrayPattern) Idx0() file.Idx          ***REMOVED*** return self.LeftBracket ***REMOVED***
func (self *ObjectPattern) Idx0() file.Idx         ***REMOVED*** return self.LeftBrace ***REMOVED***
func (self *AssignExpression) Idx0() file.Idx      ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *BadExpression) Idx0() file.Idx         ***REMOVED*** return self.From ***REMOVED***
func (self *BinaryExpression) Idx0() file.Idx      ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *BooleanLiteral) Idx0() file.Idx        ***REMOVED*** return self.Idx ***REMOVED***
func (self *BracketExpression) Idx0() file.Idx     ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *CallExpression) Idx0() file.Idx        ***REMOVED*** return self.Callee.Idx0() ***REMOVED***
func (self *ConditionalExpression) Idx0() file.Idx ***REMOVED*** return self.Test.Idx0() ***REMOVED***
func (self *DotExpression) Idx0() file.Idx         ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *PrivateDotExpression) Idx0() file.Idx  ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *FunctionLiteral) Idx0() file.Idx       ***REMOVED*** return self.Function ***REMOVED***
func (self *ClassLiteral) Idx0() file.Idx          ***REMOVED*** return self.Class ***REMOVED***
func (self *ArrowFunctionLiteral) Idx0() file.Idx  ***REMOVED*** return self.Start ***REMOVED***
func (self *Identifier) Idx0() file.Idx            ***REMOVED*** return self.Idx ***REMOVED***
func (self *NewExpression) Idx0() file.Idx         ***REMOVED*** return self.New ***REMOVED***
func (self *NullLiteral) Idx0() file.Idx           ***REMOVED*** return self.Idx ***REMOVED***
func (self *NumberLiteral) Idx0() file.Idx         ***REMOVED*** return self.Idx ***REMOVED***
func (self *ObjectLiteral) Idx0() file.Idx         ***REMOVED*** return self.LeftBrace ***REMOVED***
func (self *RegExpLiteral) Idx0() file.Idx         ***REMOVED*** return self.Idx ***REMOVED***
func (self *SequenceExpression) Idx0() file.Idx    ***REMOVED*** return self.Sequence[0].Idx0() ***REMOVED***
func (self *StringLiteral) Idx0() file.Idx         ***REMOVED*** return self.Idx ***REMOVED***
func (self *TemplateLiteral) Idx0() file.Idx       ***REMOVED*** return self.OpenQuote ***REMOVED***
func (self *ThisExpression) Idx0() file.Idx        ***REMOVED*** return self.Idx ***REMOVED***
func (self *SuperExpression) Idx0() file.Idx       ***REMOVED*** return self.Idx ***REMOVED***
func (self *UnaryExpression) Idx0() file.Idx       ***REMOVED*** return self.Idx ***REMOVED***
func (self *MetaProperty) Idx0() file.Idx          ***REMOVED*** return self.Idx ***REMOVED***

func (self *BadStatement) Idx0() file.Idx        ***REMOVED*** return self.From ***REMOVED***
func (self *BlockStatement) Idx0() file.Idx      ***REMOVED*** return self.LeftBrace ***REMOVED***
func (self *BranchStatement) Idx0() file.Idx     ***REMOVED*** return self.Idx ***REMOVED***
func (self *CaseStatement) Idx0() file.Idx       ***REMOVED*** return self.Case ***REMOVED***
func (self *CatchStatement) Idx0() file.Idx      ***REMOVED*** return self.Catch ***REMOVED***
func (self *DebuggerStatement) Idx0() file.Idx   ***REMOVED*** return self.Debugger ***REMOVED***
func (self *DoWhileStatement) Idx0() file.Idx    ***REMOVED*** return self.Do ***REMOVED***
func (self *EmptyStatement) Idx0() file.Idx      ***REMOVED*** return self.Semicolon ***REMOVED***
func (self *ExpressionStatement) Idx0() file.Idx ***REMOVED*** return self.Expression.Idx0() ***REMOVED***
func (self *ForInStatement) Idx0() file.Idx      ***REMOVED*** return self.For ***REMOVED***
func (self *ForOfStatement) Idx0() file.Idx      ***REMOVED*** return self.For ***REMOVED***
func (self *ForStatement) Idx0() file.Idx        ***REMOVED*** return self.For ***REMOVED***
func (self *IfStatement) Idx0() file.Idx         ***REMOVED*** return self.If ***REMOVED***
func (self *LabelledStatement) Idx0() file.Idx   ***REMOVED*** return self.Label.Idx0() ***REMOVED***
func (self *Program) Idx0() file.Idx             ***REMOVED*** return self.Body[0].Idx0() ***REMOVED***
func (self *ReturnStatement) Idx0() file.Idx     ***REMOVED*** return self.Return ***REMOVED***
func (self *SwitchStatement) Idx0() file.Idx     ***REMOVED*** return self.Switch ***REMOVED***
func (self *ThrowStatement) Idx0() file.Idx      ***REMOVED*** return self.Throw ***REMOVED***
func (self *TryStatement) Idx0() file.Idx        ***REMOVED*** return self.Try ***REMOVED***
func (self *VariableStatement) Idx0() file.Idx   ***REMOVED*** return self.Var ***REMOVED***
func (self *WhileStatement) Idx0() file.Idx      ***REMOVED*** return self.While ***REMOVED***
func (self *WithStatement) Idx0() file.Idx       ***REMOVED*** return self.With ***REMOVED***
func (self *LexicalDeclaration) Idx0() file.Idx  ***REMOVED*** return self.Idx ***REMOVED***
func (self *FunctionDeclaration) Idx0() file.Idx ***REMOVED*** return self.Function.Idx0() ***REMOVED***
func (self *ClassDeclaration) Idx0() file.Idx    ***REMOVED*** return self.Class.Idx0() ***REMOVED***
func (self *Binding) Idx0() file.Idx             ***REMOVED*** return self.Target.Idx0() ***REMOVED***

func (self *ForLoopInitializerVarDeclList) Idx0() file.Idx ***REMOVED*** return self.List[0].Idx0() ***REMOVED***
func (self *PropertyShort) Idx0() file.Idx                 ***REMOVED*** return self.Name.Idx ***REMOVED***
func (self *PropertyKeyed) Idx0() file.Idx                 ***REMOVED*** return self.Key.Idx0() ***REMOVED***
func (self *ExpressionBody) Idx0() file.Idx                ***REMOVED*** return self.Expression.Idx0() ***REMOVED***

func (self *FieldDefinition) Idx0() file.Idx  ***REMOVED*** return self.Idx ***REMOVED***
func (self *MethodDefinition) Idx0() file.Idx ***REMOVED*** return self.Idx ***REMOVED***
func (self *ClassStaticBlock) Idx0() file.Idx ***REMOVED*** return self.Static ***REMOVED***

func (self *ForDeclaration) Idx0() file.Idx    ***REMOVED*** return self.Idx ***REMOVED***
func (self *ForIntoVar) Idx0() file.Idx        ***REMOVED*** return self.Binding.Idx0() ***REMOVED***
func (self *ForIntoExpression) Idx0() file.Idx ***REMOVED*** return self.Expression.Idx0() ***REMOVED***

// ==== //
// Idx1 //
// ==== //

func (self *ArrayLiteral) Idx1() file.Idx          ***REMOVED*** return self.RightBracket + 1 ***REMOVED***
func (self *ArrayPattern) Idx1() file.Idx          ***REMOVED*** return self.RightBracket + 1 ***REMOVED***
func (self *AssignExpression) Idx1() file.Idx      ***REMOVED*** return self.Right.Idx1() ***REMOVED***
func (self *BadExpression) Idx1() file.Idx         ***REMOVED*** return self.To ***REMOVED***
func (self *BinaryExpression) Idx1() file.Idx      ***REMOVED*** return self.Right.Idx1() ***REMOVED***
func (self *BooleanLiteral) Idx1() file.Idx        ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *BracketExpression) Idx1() file.Idx     ***REMOVED*** return self.RightBracket + 1 ***REMOVED***
func (self *CallExpression) Idx1() file.Idx        ***REMOVED*** return self.RightParenthesis + 1 ***REMOVED***
func (self *ConditionalExpression) Idx1() file.Idx ***REMOVED*** return self.Test.Idx1() ***REMOVED***
func (self *DotExpression) Idx1() file.Idx         ***REMOVED*** return self.Identifier.Idx1() ***REMOVED***
func (self *PrivateDotExpression) Idx1() file.Idx  ***REMOVED*** return self.Identifier.Idx1() ***REMOVED***
func (self *FunctionLiteral) Idx1() file.Idx       ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *ClassLiteral) Idx1() file.Idx          ***REMOVED*** return self.RightBrace + 1 ***REMOVED***
func (self *ArrowFunctionLiteral) Idx1() file.Idx  ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *Identifier) Idx1() file.Idx            ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Name)) ***REMOVED***
func (self *NewExpression) Idx1() file.Idx ***REMOVED***
	if self.ArgumentList != nil ***REMOVED***
		return self.RightParenthesis + 1
	***REMOVED*** else ***REMOVED***
		return self.Callee.Idx1()
	***REMOVED***
***REMOVED***
func (self *NullLiteral) Idx1() file.Idx        ***REMOVED*** return file.Idx(int(self.Idx) + 4) ***REMOVED*** // "null"
func (self *NumberLiteral) Idx1() file.Idx      ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *ObjectLiteral) Idx1() file.Idx      ***REMOVED*** return self.RightBrace + 1 ***REMOVED***
func (self *ObjectPattern) Idx1() file.Idx      ***REMOVED*** return self.RightBrace + 1 ***REMOVED***
func (self *RegExpLiteral) Idx1() file.Idx      ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *SequenceExpression) Idx1() file.Idx ***REMOVED*** return self.Sequence[len(self.Sequence)-1].Idx1() ***REMOVED***
func (self *StringLiteral) Idx1() file.Idx      ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *TemplateLiteral) Idx1() file.Idx    ***REMOVED*** return self.CloseQuote + 1 ***REMOVED***
func (self *ThisExpression) Idx1() file.Idx     ***REMOVED*** return self.Idx + 4 ***REMOVED***
func (self *SuperExpression) Idx1() file.Idx    ***REMOVED*** return self.Idx + 5 ***REMOVED***
func (self *UnaryExpression) Idx1() file.Idx ***REMOVED***
	if self.Postfix ***REMOVED***
		return self.Operand.Idx1() + 2 // ++ --
	***REMOVED***
	return self.Operand.Idx1()
***REMOVED***
func (self *MetaProperty) Idx1() file.Idx ***REMOVED***
	return self.Property.Idx1()
***REMOVED***

func (self *BadStatement) Idx1() file.Idx        ***REMOVED*** return self.To ***REMOVED***
func (self *BlockStatement) Idx1() file.Idx      ***REMOVED*** return self.RightBrace + 1 ***REMOVED***
func (self *BranchStatement) Idx1() file.Idx     ***REMOVED*** return self.Idx ***REMOVED***
func (self *CaseStatement) Idx1() file.Idx       ***REMOVED*** return self.Consequent[len(self.Consequent)-1].Idx1() ***REMOVED***
func (self *CatchStatement) Idx1() file.Idx      ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *DebuggerStatement) Idx1() file.Idx   ***REMOVED*** return self.Debugger + 8 ***REMOVED***
func (self *DoWhileStatement) Idx1() file.Idx    ***REMOVED*** return self.Test.Idx1() ***REMOVED***
func (self *EmptyStatement) Idx1() file.Idx      ***REMOVED*** return self.Semicolon + 1 ***REMOVED***
func (self *ExpressionStatement) Idx1() file.Idx ***REMOVED*** return self.Expression.Idx1() ***REMOVED***
func (self *ForInStatement) Idx1() file.Idx      ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *ForOfStatement) Idx1() file.Idx      ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *ForStatement) Idx1() file.Idx        ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *IfStatement) Idx1() file.Idx ***REMOVED***
	if self.Alternate != nil ***REMOVED***
		return self.Alternate.Idx1()
	***REMOVED***
	return self.Consequent.Idx1()
***REMOVED***
func (self *LabelledStatement) Idx1() file.Idx ***REMOVED*** return self.Colon + 1 ***REMOVED***
func (self *Program) Idx1() file.Idx           ***REMOVED*** return self.Body[len(self.Body)-1].Idx1() ***REMOVED***
func (self *ReturnStatement) Idx1() file.Idx   ***REMOVED*** return self.Return + 6 ***REMOVED***
func (self *SwitchStatement) Idx1() file.Idx   ***REMOVED*** return self.Body[len(self.Body)-1].Idx1() ***REMOVED***
func (self *ThrowStatement) Idx1() file.Idx    ***REMOVED*** return self.Argument.Idx1() ***REMOVED***
func (self *TryStatement) Idx1() file.Idx ***REMOVED***
	if self.Finally != nil ***REMOVED***
		return self.Finally.Idx1()
	***REMOVED***
	if self.Catch != nil ***REMOVED***
		return self.Catch.Idx1()
	***REMOVED***
	return self.Body.Idx1()
***REMOVED***
func (self *VariableStatement) Idx1() file.Idx   ***REMOVED*** return self.List[len(self.List)-1].Idx1() ***REMOVED***
func (self *WhileStatement) Idx1() file.Idx      ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *WithStatement) Idx1() file.Idx       ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *LexicalDeclaration) Idx1() file.Idx  ***REMOVED*** return self.List[len(self.List)-1].Idx1() ***REMOVED***
func (self *FunctionDeclaration) Idx1() file.Idx ***REMOVED*** return self.Function.Idx1() ***REMOVED***
func (self *ClassDeclaration) Idx1() file.Idx    ***REMOVED*** return self.Class.Idx1() ***REMOVED***
func (self *Binding) Idx1() file.Idx ***REMOVED***
	if self.Initializer != nil ***REMOVED***
		return self.Initializer.Idx1()
	***REMOVED***
	return self.Target.Idx1()
***REMOVED***

func (self *ForLoopInitializerVarDeclList) Idx1() file.Idx ***REMOVED*** return self.List[len(self.List)-1].Idx1() ***REMOVED***

func (self *PropertyShort) Idx1() file.Idx ***REMOVED***
	if self.Initializer != nil ***REMOVED***
		return self.Initializer.Idx1()
	***REMOVED***
	return self.Name.Idx1()
***REMOVED***

func (self *PropertyKeyed) Idx1() file.Idx ***REMOVED*** return self.Value.Idx1() ***REMOVED***

func (self *ExpressionBody) Idx1() file.Idx ***REMOVED*** return self.Expression.Idx1() ***REMOVED***

func (self *FieldDefinition) Idx1() file.Idx ***REMOVED***
	if self.Initializer != nil ***REMOVED***
		return self.Initializer.Idx1()
	***REMOVED***
	return self.Key.Idx1()
***REMOVED***

func (self *MethodDefinition) Idx1() file.Idx ***REMOVED***
	return self.Body.Idx1()
***REMOVED***

func (self *ClassStaticBlock) Idx1() file.Idx ***REMOVED***
	return self.Block.Idx1()
***REMOVED***

func (self *ForDeclaration) Idx1() file.Idx    ***REMOVED*** return self.Target.Idx1() ***REMOVED***
func (self *ForIntoVar) Idx1() file.Idx        ***REMOVED*** return self.Binding.Idx1() ***REMOVED***
func (self *ForIntoExpression) Idx1() file.Idx ***REMOVED*** return self.Expression.Idx1() ***REMOVED***
