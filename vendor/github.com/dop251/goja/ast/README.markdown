# ast
--
    import "github.com/dop251/goja/ast"

Package ast declares types representing a JavaScript AST.


### Warning

The parser and AST interfaces are still works-in-progress (particularly where
node types are concerned) and may change in the future.

## Usage

#### type ArrayLiteral

```go
type ArrayLiteral struct ***REMOVED***
	LeftBracket  file.Idx
	RightBracket file.Idx
	Value        []Expression
***REMOVED***
```


#### func (*ArrayLiteral) Idx0

```go
func (self *ArrayLiteral) Idx0() file.Idx
```

#### func (*ArrayLiteral) Idx1

```go
func (self *ArrayLiteral) Idx1() file.Idx
```

#### type AssignExpression

```go
type AssignExpression struct ***REMOVED***
	Operator token.Token
	Left     Expression
	Right    Expression
***REMOVED***
```


#### func (*AssignExpression) Idx0

```go
func (self *AssignExpression) Idx0() file.Idx
```

#### func (*AssignExpression) Idx1

```go
func (self *AssignExpression) Idx1() file.Idx
```

#### type BadExpression

```go
type BadExpression struct ***REMOVED***
	From file.Idx
	To   file.Idx
***REMOVED***
```


#### func (*BadExpression) Idx0

```go
func (self *BadExpression) Idx0() file.Idx
```

#### func (*BadExpression) Idx1

```go
func (self *BadExpression) Idx1() file.Idx
```

#### type BadStatement

```go
type BadStatement struct ***REMOVED***
	From file.Idx
	To   file.Idx
***REMOVED***
```


#### func (*BadStatement) Idx0

```go
func (self *BadStatement) Idx0() file.Idx
```

#### func (*BadStatement) Idx1

```go
func (self *BadStatement) Idx1() file.Idx
```

#### type BinaryExpression

```go
type BinaryExpression struct ***REMOVED***
	Operator   token.Token
	Left       Expression
	Right      Expression
	Comparison bool
***REMOVED***
```


#### func (*BinaryExpression) Idx0

```go
func (self *BinaryExpression) Idx0() file.Idx
```

#### func (*BinaryExpression) Idx1

```go
func (self *BinaryExpression) Idx1() file.Idx
```

#### type BlockStatement

```go
type BlockStatement struct ***REMOVED***
	LeftBrace  file.Idx
	List       []Statement
	RightBrace file.Idx
***REMOVED***
```


#### func (*BlockStatement) Idx0

```go
func (self *BlockStatement) Idx0() file.Idx
```

#### func (*BlockStatement) Idx1

```go
func (self *BlockStatement) Idx1() file.Idx
```

#### type BooleanLiteral

```go
type BooleanLiteral struct ***REMOVED***
	Idx     file.Idx
	Literal string
	Value   bool
***REMOVED***
```


#### func (*BooleanLiteral) Idx0

```go
func (self *BooleanLiteral) Idx0() file.Idx
```

#### func (*BooleanLiteral) Idx1

```go
func (self *BooleanLiteral) Idx1() file.Idx
```

#### type BracketExpression

```go
type BracketExpression struct ***REMOVED***
	Left         Expression
	Member       Expression
	LeftBracket  file.Idx
	RightBracket file.Idx
***REMOVED***
```


#### func (*BracketExpression) Idx0

```go
func (self *BracketExpression) Idx0() file.Idx
```

#### func (*BracketExpression) Idx1

```go
func (self *BracketExpression) Idx1() file.Idx
```

#### type BranchStatement

```go
type BranchStatement struct ***REMOVED***
	Idx   file.Idx
	Token token.Token
	Label *Identifier
***REMOVED***
```


#### func (*BranchStatement) Idx0

```go
func (self *BranchStatement) Idx0() file.Idx
```

#### func (*BranchStatement) Idx1

```go
func (self *BranchStatement) Idx1() file.Idx
```

#### type CallExpression

```go
type CallExpression struct ***REMOVED***
	Callee           Expression
	LeftParenthesis  file.Idx
	ArgumentList     []Expression
	RightParenthesis file.Idx
***REMOVED***
```


#### func (*CallExpression) Idx0

```go
func (self *CallExpression) Idx0() file.Idx
```

#### func (*CallExpression) Idx1

```go
func (self *CallExpression) Idx1() file.Idx
```

#### type CaseStatement

```go
type CaseStatement struct ***REMOVED***
	Case       file.Idx
	Test       Expression
	Consequent []Statement
***REMOVED***
```


#### func (*CaseStatement) Idx0

```go
func (self *CaseStatement) Idx0() file.Idx
```

#### func (*CaseStatement) Idx1

```go
func (self *CaseStatement) Idx1() file.Idx
```

#### type CatchStatement

```go
type CatchStatement struct ***REMOVED***
	Catch     file.Idx
	Parameter *Identifier
	Body      Statement
***REMOVED***
```


#### func (*CatchStatement) Idx0

```go
func (self *CatchStatement) Idx0() file.Idx
```

#### func (*CatchStatement) Idx1

```go
func (self *CatchStatement) Idx1() file.Idx
```

#### type ConditionalExpression

```go
type ConditionalExpression struct ***REMOVED***
	Test       Expression
	Consequent Expression
	Alternate  Expression
***REMOVED***
```


#### func (*ConditionalExpression) Idx0

```go
func (self *ConditionalExpression) Idx0() file.Idx
```

#### func (*ConditionalExpression) Idx1

```go
func (self *ConditionalExpression) Idx1() file.Idx
```

#### type DebuggerStatement

```go
type DebuggerStatement struct ***REMOVED***
	Debugger file.Idx
***REMOVED***
```


#### func (*DebuggerStatement) Idx0

```go
func (self *DebuggerStatement) Idx0() file.Idx
```

#### func (*DebuggerStatement) Idx1

```go
func (self *DebuggerStatement) Idx1() file.Idx
```

#### type Declaration

```go
type Declaration interface ***REMOVED***
	// contains filtered or unexported methods
***REMOVED***
```

All declaration nodes implement the Declaration interface.

#### type DoWhileStatement

```go
type DoWhileStatement struct ***REMOVED***
	Do   file.Idx
	Test Expression
	Body Statement
***REMOVED***
```


#### func (*DoWhileStatement) Idx0

```go
func (self *DoWhileStatement) Idx0() file.Idx
```

#### func (*DoWhileStatement) Idx1

```go
func (self *DoWhileStatement) Idx1() file.Idx
```

#### type DotExpression

```go
type DotExpression struct ***REMOVED***
	Left       Expression
	Identifier Identifier
***REMOVED***
```


#### func (*DotExpression) Idx0

```go
func (self *DotExpression) Idx0() file.Idx
```

#### func (*DotExpression) Idx1

```go
func (self *DotExpression) Idx1() file.Idx
```

#### type EmptyStatement

```go
type EmptyStatement struct ***REMOVED***
	Semicolon file.Idx
***REMOVED***
```


#### func (*EmptyStatement) Idx0

```go
func (self *EmptyStatement) Idx0() file.Idx
```

#### func (*EmptyStatement) Idx1

```go
func (self *EmptyStatement) Idx1() file.Idx
```

#### type Expression

```go
type Expression interface ***REMOVED***
	Node
	// contains filtered or unexported methods
***REMOVED***
```

All expression nodes implement the Expression interface.

#### type ExpressionStatement

```go
type ExpressionStatement struct ***REMOVED***
	Expression Expression
***REMOVED***
```


#### func (*ExpressionStatement) Idx0

```go
func (self *ExpressionStatement) Idx0() file.Idx
```

#### func (*ExpressionStatement) Idx1

```go
func (self *ExpressionStatement) Idx1() file.Idx
```

#### type ForInStatement

```go
type ForInStatement struct ***REMOVED***
	For    file.Idx
	Into   Expression
	Source Expression
	Body   Statement
***REMOVED***
```


#### func (*ForInStatement) Idx0

```go
func (self *ForInStatement) Idx0() file.Idx
```

#### func (*ForInStatement) Idx1

```go
func (self *ForInStatement) Idx1() file.Idx
```

#### type ForStatement

```go
type ForStatement struct ***REMOVED***
	For         file.Idx
	Initializer Expression
	Update      Expression
	Test        Expression
	Body        Statement
***REMOVED***
```


#### func (*ForStatement) Idx0

```go
func (self *ForStatement) Idx0() file.Idx
```

#### func (*ForStatement) Idx1

```go
func (self *ForStatement) Idx1() file.Idx
```

#### type FunctionDeclaration

```go
type FunctionDeclaration struct ***REMOVED***
	Function *FunctionLiteral
***REMOVED***
```


#### type FunctionLiteral

```go
type FunctionLiteral struct ***REMOVED***
	Function      file.Idx
	Name          *Identifier
	ParameterList *ParameterList
	Body          Statement
	Source        string

	DeclarationList []Declaration
***REMOVED***
```


#### func (*FunctionLiteral) Idx0

```go
func (self *FunctionLiteral) Idx0() file.Idx
```

#### func (*FunctionLiteral) Idx1

```go
func (self *FunctionLiteral) Idx1() file.Idx
```

#### type Identifier

```go
type Identifier struct ***REMOVED***
	Name string
	Idx  file.Idx
***REMOVED***
```


#### func (*Identifier) Idx0

```go
func (self *Identifier) Idx0() file.Idx
```

#### func (*Identifier) Idx1

```go
func (self *Identifier) Idx1() file.Idx
```

#### type IfStatement

```go
type IfStatement struct ***REMOVED***
	If         file.Idx
	Test       Expression
	Consequent Statement
	Alternate  Statement
***REMOVED***
```


#### func (*IfStatement) Idx0

```go
func (self *IfStatement) Idx0() file.Idx
```

#### func (*IfStatement) Idx1

```go
func (self *IfStatement) Idx1() file.Idx
```

#### type LabelledStatement

```go
type LabelledStatement struct ***REMOVED***
	Label     *Identifier
	Colon     file.Idx
	Statement Statement
***REMOVED***
```


#### func (*LabelledStatement) Idx0

```go
func (self *LabelledStatement) Idx0() file.Idx
```

#### func (*LabelledStatement) Idx1

```go
func (self *LabelledStatement) Idx1() file.Idx
```

#### type NewExpression

```go
type NewExpression struct ***REMOVED***
	New              file.Idx
	Callee           Expression
	LeftParenthesis  file.Idx
	ArgumentList     []Expression
	RightParenthesis file.Idx
***REMOVED***
```


#### func (*NewExpression) Idx0

```go
func (self *NewExpression) Idx0() file.Idx
```

#### func (*NewExpression) Idx1

```go
func (self *NewExpression) Idx1() file.Idx
```

#### type Node

```go
type Node interface ***REMOVED***
	Idx0() file.Idx // The index of the first character belonging to the node
	Idx1() file.Idx // The index of the first character immediately after the node
***REMOVED***
```

All nodes implement the Node interface.

#### type NullLiteral

```go
type NullLiteral struct ***REMOVED***
	Idx     file.Idx
	Literal string
***REMOVED***
```


#### func (*NullLiteral) Idx0

```go
func (self *NullLiteral) Idx0() file.Idx
```

#### func (*NullLiteral) Idx1

```go
func (self *NullLiteral) Idx1() file.Idx
```

#### type NumberLiteral

```go
type NumberLiteral struct ***REMOVED***
	Idx     file.Idx
	Literal string
	Value   interface***REMOVED******REMOVED***
***REMOVED***
```


#### func (*NumberLiteral) Idx0

```go
func (self *NumberLiteral) Idx0() file.Idx
```

#### func (*NumberLiteral) Idx1

```go
func (self *NumberLiteral) Idx1() file.Idx
```

#### type ObjectLiteral

```go
type ObjectLiteral struct ***REMOVED***
	LeftBrace  file.Idx
	RightBrace file.Idx
	Value      []Property
***REMOVED***
```


#### func (*ObjectLiteral) Idx0

```go
func (self *ObjectLiteral) Idx0() file.Idx
```

#### func (*ObjectLiteral) Idx1

```go
func (self *ObjectLiteral) Idx1() file.Idx
```

#### type ParameterList

```go
type ParameterList struct ***REMOVED***
	Opening file.Idx
	List    []*Identifier
	Closing file.Idx
***REMOVED***
```


#### type Program

```go
type Program struct ***REMOVED***
	Body []Statement

	DeclarationList []Declaration

	File *file.File
***REMOVED***
```


#### func (*Program) Idx0

```go
func (self *Program) Idx0() file.Idx
```

#### func (*Program) Idx1

```go
func (self *Program) Idx1() file.Idx
```

#### type Property

```go
type Property struct ***REMOVED***
	Key   string
	Kind  string
	Value Expression
***REMOVED***
```


#### type RegExpLiteral

```go
type RegExpLiteral struct ***REMOVED***
	Idx     file.Idx
	Literal string
	Pattern string
	Flags   string
	Value   string
***REMOVED***
```


#### func (*RegExpLiteral) Idx0

```go
func (self *RegExpLiteral) Idx0() file.Idx
```

#### func (*RegExpLiteral) Idx1

```go
func (self *RegExpLiteral) Idx1() file.Idx
```

#### type ReturnStatement

```go
type ReturnStatement struct ***REMOVED***
	Return   file.Idx
	Argument Expression
***REMOVED***
```


#### func (*ReturnStatement) Idx0

```go
func (self *ReturnStatement) Idx0() file.Idx
```

#### func (*ReturnStatement) Idx1

```go
func (self *ReturnStatement) Idx1() file.Idx
```

#### type SequenceExpression

```go
type SequenceExpression struct ***REMOVED***
	Sequence []Expression
***REMOVED***
```


#### func (*SequenceExpression) Idx0

```go
func (self *SequenceExpression) Idx0() file.Idx
```

#### func (*SequenceExpression) Idx1

```go
func (self *SequenceExpression) Idx1() file.Idx
```

#### type Statement

```go
type Statement interface ***REMOVED***
	Node
	// contains filtered or unexported methods
***REMOVED***
```

All statement nodes implement the Statement interface.

#### type StringLiteral

```go
type StringLiteral struct ***REMOVED***
	Idx     file.Idx
	Literal string
	Value   string
***REMOVED***
```


#### func (*StringLiteral) Idx0

```go
func (self *StringLiteral) Idx0() file.Idx
```

#### func (*StringLiteral) Idx1

```go
func (self *StringLiteral) Idx1() file.Idx
```

#### type SwitchStatement

```go
type SwitchStatement struct ***REMOVED***
	Switch       file.Idx
	Discriminant Expression
	Default      int
	Body         []*CaseStatement
***REMOVED***
```


#### func (*SwitchStatement) Idx0

```go
func (self *SwitchStatement) Idx0() file.Idx
```

#### func (*SwitchStatement) Idx1

```go
func (self *SwitchStatement) Idx1() file.Idx
```

#### type ThisExpression

```go
type ThisExpression struct ***REMOVED***
	Idx file.Idx
***REMOVED***
```


#### func (*ThisExpression) Idx0

```go
func (self *ThisExpression) Idx0() file.Idx
```

#### func (*ThisExpression) Idx1

```go
func (self *ThisExpression) Idx1() file.Idx
```

#### type ThrowStatement

```go
type ThrowStatement struct ***REMOVED***
	Throw    file.Idx
	Argument Expression
***REMOVED***
```


#### func (*ThrowStatement) Idx0

```go
func (self *ThrowStatement) Idx0() file.Idx
```

#### func (*ThrowStatement) Idx1

```go
func (self *ThrowStatement) Idx1() file.Idx
```

#### type TryStatement

```go
type TryStatement struct ***REMOVED***
	Try     file.Idx
	Body    Statement
	Catch   *CatchStatement
	Finally Statement
***REMOVED***
```


#### func (*TryStatement) Idx0

```go
func (self *TryStatement) Idx0() file.Idx
```

#### func (*TryStatement) Idx1

```go
func (self *TryStatement) Idx1() file.Idx
```

#### type UnaryExpression

```go
type UnaryExpression struct ***REMOVED***
	Operator token.Token
	Idx      file.Idx // If a prefix operation
	Operand  Expression
	Postfix  bool
***REMOVED***
```


#### func (*UnaryExpression) Idx0

```go
func (self *UnaryExpression) Idx0() file.Idx
```

#### func (*UnaryExpression) Idx1

```go
func (self *UnaryExpression) Idx1() file.Idx
```

#### type VariableDeclaration

```go
type VariableDeclaration struct ***REMOVED***
	Var  file.Idx
	List []*VariableExpression
***REMOVED***
```


#### type VariableExpression

```go
type VariableExpression struct ***REMOVED***
	Name        string
	Idx         file.Idx
	Initializer Expression
***REMOVED***
```


#### func (*VariableExpression) Idx0

```go
func (self *VariableExpression) Idx0() file.Idx
```

#### func (*VariableExpression) Idx1

```go
func (self *VariableExpression) Idx1() file.Idx
```

#### type VariableStatement

```go
type VariableStatement struct ***REMOVED***
	Var  file.Idx
	List []Expression
***REMOVED***
```


#### func (*VariableStatement) Idx0

```go
func (self *VariableStatement) Idx0() file.Idx
```

#### func (*VariableStatement) Idx1

```go
func (self *VariableStatement) Idx1() file.Idx
```

#### type WhileStatement

```go
type WhileStatement struct ***REMOVED***
	While file.Idx
	Test  Expression
	Body  Statement
***REMOVED***
```


#### func (*WhileStatement) Idx0

```go
func (self *WhileStatement) Idx0() file.Idx
```

#### func (*WhileStatement) Idx1

```go
func (self *WhileStatement) Idx1() file.Idx
```

#### type WithStatement

```go
type WithStatement struct ***REMOVED***
	With   file.Idx
	Object Expression
	Body   Statement
***REMOVED***
```


#### func (*WithStatement) Idx0

```go
func (self *WithStatement) Idx0() file.Idx
```

#### func (*WithStatement) Idx1

```go
func (self *WithStatement) Idx1() file.Idx
```

--
**godocdown** http://github.com/robertkrimen/godocdown
