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

	ArrayLiteral struct ***REMOVED***
		LeftBracket  file.Idx
		RightBracket file.Idx
		Value        []Expression
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

	FunctionLiteral struct ***REMOVED***
		Function      file.Idx
		Name          *Identifier
		ParameterList *ParameterList
		Body          Statement
		Source        string

		DeclarationList []Declaration
	***REMOVED***

	Identifier struct ***REMOVED***
		Name unistring.String
		Idx  file.Idx
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

	ParameterList struct ***REMOVED***
		Opening file.Idx
		List    []*Identifier
		Closing file.Idx
	***REMOVED***

	Property struct ***REMOVED***
		Key   Expression
		Kind  string
		Value Expression
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

	ThisExpression struct ***REMOVED***
		Idx file.Idx
	***REMOVED***

	UnaryExpression struct ***REMOVED***
		Operator token.Token
		Idx      file.Idx // If a prefix operation
		Operand  Expression
		Postfix  bool
	***REMOVED***

	VariableExpression struct ***REMOVED***
		Name        unistring.String
		Idx         file.Idx
		Initializer Expression
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
func (*FunctionLiteral) _expressionNode()       ***REMOVED******REMOVED***
func (*Identifier) _expressionNode()            ***REMOVED******REMOVED***
func (*NewExpression) _expressionNode()         ***REMOVED******REMOVED***
func (*NullLiteral) _expressionNode()           ***REMOVED******REMOVED***
func (*NumberLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*ObjectLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*RegExpLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*SequenceExpression) _expressionNode()    ***REMOVED******REMOVED***
func (*StringLiteral) _expressionNode()         ***REMOVED******REMOVED***
func (*ThisExpression) _expressionNode()        ***REMOVED******REMOVED***
func (*UnaryExpression) _expressionNode()       ***REMOVED******REMOVED***
func (*VariableExpression) _expressionNode()    ***REMOVED******REMOVED***
func (*MetaProperty) _expressionNode()          ***REMOVED******REMOVED***

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
		Parameter *Identifier
		Body      Statement
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
		Into   Expression
		Source Expression
		Body   Statement
	***REMOVED***

	ForOfStatement struct ***REMOVED***
		For    file.Idx
		Into   Expression
		Source Expression
		Body   Statement
	***REMOVED***

	ForStatement struct ***REMOVED***
		For         file.Idx
		Initializer Expression
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
		Body    Statement
		Catch   *CatchStatement
		Finally Statement
	***REMOVED***

	VariableStatement struct ***REMOVED***
		Var  file.Idx
		List []Expression
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

// =========== //
// Declaration //
// =========== //

type (
	// All declaration nodes implement the Declaration interface.
	Declaration interface ***REMOVED***
		_declarationNode()
	***REMOVED***

	FunctionDeclaration struct ***REMOVED***
		Function *FunctionLiteral
	***REMOVED***

	VariableDeclaration struct ***REMOVED***
		Var  file.Idx
		List []*VariableExpression
	***REMOVED***
)

// _declarationNode

func (*FunctionDeclaration) _declarationNode() ***REMOVED******REMOVED***
func (*VariableDeclaration) _declarationNode() ***REMOVED******REMOVED***

// ==== //
// Node //
// ==== //

type Program struct ***REMOVED***
	Body []Statement

	DeclarationList []Declaration

	File *file.File
***REMOVED***

// ==== //
// Idx0 //
// ==== //

func (self *ArrayLiteral) Idx0() file.Idx          ***REMOVED*** return self.LeftBracket ***REMOVED***
func (self *AssignExpression) Idx0() file.Idx      ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *BadExpression) Idx0() file.Idx         ***REMOVED*** return self.From ***REMOVED***
func (self *BinaryExpression) Idx0() file.Idx      ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *BooleanLiteral) Idx0() file.Idx        ***REMOVED*** return self.Idx ***REMOVED***
func (self *BracketExpression) Idx0() file.Idx     ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *CallExpression) Idx0() file.Idx        ***REMOVED*** return self.Callee.Idx0() ***REMOVED***
func (self *ConditionalExpression) Idx0() file.Idx ***REMOVED*** return self.Test.Idx0() ***REMOVED***
func (self *DotExpression) Idx0() file.Idx         ***REMOVED*** return self.Left.Idx0() ***REMOVED***
func (self *FunctionLiteral) Idx0() file.Idx       ***REMOVED*** return self.Function ***REMOVED***
func (self *Identifier) Idx0() file.Idx            ***REMOVED*** return self.Idx ***REMOVED***
func (self *NewExpression) Idx0() file.Idx         ***REMOVED*** return self.New ***REMOVED***
func (self *NullLiteral) Idx0() file.Idx           ***REMOVED*** return self.Idx ***REMOVED***
func (self *NumberLiteral) Idx0() file.Idx         ***REMOVED*** return self.Idx ***REMOVED***
func (self *ObjectLiteral) Idx0() file.Idx         ***REMOVED*** return self.LeftBrace ***REMOVED***
func (self *RegExpLiteral) Idx0() file.Idx         ***REMOVED*** return self.Idx ***REMOVED***
func (self *SequenceExpression) Idx0() file.Idx    ***REMOVED*** return self.Sequence[0].Idx0() ***REMOVED***
func (self *StringLiteral) Idx0() file.Idx         ***REMOVED*** return self.Idx ***REMOVED***
func (self *ThisExpression) Idx0() file.Idx        ***REMOVED*** return self.Idx ***REMOVED***
func (self *UnaryExpression) Idx0() file.Idx       ***REMOVED*** return self.Idx ***REMOVED***
func (self *VariableExpression) Idx0() file.Idx    ***REMOVED*** return self.Idx ***REMOVED***
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

// ==== //
// Idx1 //
// ==== //

func (self *ArrayLiteral) Idx1() file.Idx          ***REMOVED*** return self.RightBracket ***REMOVED***
func (self *AssignExpression) Idx1() file.Idx      ***REMOVED*** return self.Right.Idx1() ***REMOVED***
func (self *BadExpression) Idx1() file.Idx         ***REMOVED*** return self.To ***REMOVED***
func (self *BinaryExpression) Idx1() file.Idx      ***REMOVED*** return self.Right.Idx1() ***REMOVED***
func (self *BooleanLiteral) Idx1() file.Idx        ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *BracketExpression) Idx1() file.Idx     ***REMOVED*** return self.RightBracket + 1 ***REMOVED***
func (self *CallExpression) Idx1() file.Idx        ***REMOVED*** return self.RightParenthesis + 1 ***REMOVED***
func (self *ConditionalExpression) Idx1() file.Idx ***REMOVED*** return self.Test.Idx1() ***REMOVED***
func (self *DotExpression) Idx1() file.Idx         ***REMOVED*** return self.Identifier.Idx1() ***REMOVED***
func (self *FunctionLiteral) Idx1() file.Idx       ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *Identifier) Idx1() file.Idx            ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Name)) ***REMOVED***
func (self *NewExpression) Idx1() file.Idx         ***REMOVED*** return self.RightParenthesis + 1 ***REMOVED***
func (self *NullLiteral) Idx1() file.Idx           ***REMOVED*** return file.Idx(int(self.Idx) + 4) ***REMOVED*** // "null"
func (self *NumberLiteral) Idx1() file.Idx         ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *ObjectLiteral) Idx1() file.Idx         ***REMOVED*** return self.RightBrace ***REMOVED***
func (self *RegExpLiteral) Idx1() file.Idx         ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *SequenceExpression) Idx1() file.Idx    ***REMOVED*** return self.Sequence[0].Idx1() ***REMOVED***
func (self *StringLiteral) Idx1() file.Idx         ***REMOVED*** return file.Idx(int(self.Idx) + len(self.Literal)) ***REMOVED***
func (self *ThisExpression) Idx1() file.Idx        ***REMOVED*** return self.Idx ***REMOVED***
func (self *UnaryExpression) Idx1() file.Idx ***REMOVED***
	if self.Postfix ***REMOVED***
		return self.Operand.Idx1() + 2 // ++ --
	***REMOVED***
	return self.Operand.Idx1()
***REMOVED***
func (self *VariableExpression) Idx1() file.Idx ***REMOVED***
	if self.Initializer == nil ***REMOVED***
		return file.Idx(int(self.Idx) + len(self.Name) + 1)
	***REMOVED***
	return self.Initializer.Idx1()
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
func (self *ReturnStatement) Idx1() file.Idx   ***REMOVED*** return self.Return ***REMOVED***
func (self *SwitchStatement) Idx1() file.Idx   ***REMOVED*** return self.Body[len(self.Body)-1].Idx1() ***REMOVED***
func (self *ThrowStatement) Idx1() file.Idx    ***REMOVED*** return self.Throw ***REMOVED***
func (self *TryStatement) Idx1() file.Idx      ***REMOVED*** return self.Try ***REMOVED***
func (self *VariableStatement) Idx1() file.Idx ***REMOVED*** return self.List[len(self.List)-1].Idx1() ***REMOVED***
func (self *WhileStatement) Idx1() file.Idx    ***REMOVED*** return self.Body.Idx1() ***REMOVED***
func (self *WithStatement) Idx1() file.Idx     ***REMOVED*** return self.Body.Idx1() ***REMOVED***
