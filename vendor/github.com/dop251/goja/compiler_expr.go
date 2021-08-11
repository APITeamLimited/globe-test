package goja

import (
	"fmt"
	"regexp"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
)

var (
	octalRegexp = regexp.MustCompile(`^0[0-7]`)
)

type compiledExpr interface ***REMOVED***
	emitGetter(putOnStack bool)
	emitSetter(valueExpr compiledExpr, putOnStack bool)
	emitRef()
	emitUnary(prepare, body func(), postfix, putOnStack bool)
	deleteExpr() compiledExpr
	constant() bool
	addSrcMap()
***REMOVED***

type compiledExprOrRef interface ***REMOVED***
	compiledExpr
	emitGetterOrRef()
***REMOVED***

type compiledCallExpr struct ***REMOVED***
	baseCompiledExpr
	args   []compiledExpr
	callee compiledExpr

	isVariadic bool
***REMOVED***

type compiledNewExpr struct ***REMOVED***
	compiledCallExpr
***REMOVED***

type compiledObjectLiteral struct ***REMOVED***
	baseCompiledExpr
	expr *ast.ObjectLiteral
***REMOVED***

type compiledArrayLiteral struct ***REMOVED***
	baseCompiledExpr
	expr *ast.ArrayLiteral
***REMOVED***

type compiledRegexpLiteral struct ***REMOVED***
	baseCompiledExpr
	expr *ast.RegExpLiteral
***REMOVED***

type compiledLiteral struct ***REMOVED***
	baseCompiledExpr
	val Value
***REMOVED***

type compiledAssignExpr struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
	operator    token.Token
***REMOVED***

type compiledObjectAssignmentPattern struct ***REMOVED***
	baseCompiledExpr
	expr *ast.ObjectPattern
***REMOVED***

type compiledArrayAssignmentPattern struct ***REMOVED***
	baseCompiledExpr
	expr *ast.ArrayPattern
***REMOVED***

type deleteGlobalExpr struct ***REMOVED***
	baseCompiledExpr
	name unistring.String
***REMOVED***

type deleteVarExpr struct ***REMOVED***
	baseCompiledExpr
	name unistring.String
***REMOVED***

type deletePropExpr struct ***REMOVED***
	baseCompiledExpr
	left compiledExpr
	name unistring.String
***REMOVED***

type deleteElemExpr struct ***REMOVED***
	baseCompiledExpr
	left, member compiledExpr
***REMOVED***

type constantExpr struct ***REMOVED***
	baseCompiledExpr
	val Value
***REMOVED***

type baseCompiledExpr struct ***REMOVED***
	c      *compiler
	offset int
***REMOVED***

type compiledIdentifierExpr struct ***REMOVED***
	baseCompiledExpr
	name unistring.String
***REMOVED***

type compiledFunctionLiteral struct ***REMOVED***
	baseCompiledExpr
	name            *ast.Identifier
	parameterList   *ast.ParameterList
	body            []ast.Statement
	source          string
	declarationList []*ast.VariableDeclaration
	lhsName         unistring.String
	strict          *ast.StringLiteral
	isExpr          bool
	isArrow         bool
***REMOVED***

type compiledBracketExpr struct ***REMOVED***
	baseCompiledExpr
	left, member compiledExpr
***REMOVED***

type compiledThisExpr struct ***REMOVED***
	baseCompiledExpr
***REMOVED***

type compiledNewTarget struct ***REMOVED***
	baseCompiledExpr
***REMOVED***

type compiledSequenceExpr struct ***REMOVED***
	baseCompiledExpr
	sequence []compiledExpr
***REMOVED***

type compiledUnaryExpr struct ***REMOVED***
	baseCompiledExpr
	operand  compiledExpr
	operator token.Token
	postfix  bool
***REMOVED***

type compiledConditionalExpr struct ***REMOVED***
	baseCompiledExpr
	test, consequent, alternate compiledExpr
***REMOVED***

type compiledLogicalOr struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
***REMOVED***

type compiledLogicalAnd struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
***REMOVED***

type compiledBinaryExpr struct ***REMOVED***
	baseCompiledExpr
	left, right compiledExpr
	operator    token.Token
***REMOVED***

type compiledEnumGetExpr struct ***REMOVED***
	baseCompiledExpr
***REMOVED***

type defaultDeleteExpr struct ***REMOVED***
	baseCompiledExpr
	expr compiledExpr
***REMOVED***

type compiledSpreadCallArgument struct ***REMOVED***
	baseCompiledExpr
	expr compiledExpr
***REMOVED***

func (e *defaultDeleteExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.expr.emitGetter(false)
	if putOnStack ***REMOVED***
		e.c.emit(loadVal(e.c.p.defineLiteralValue(valueTrue)))
	***REMOVED***
***REMOVED***

func (c *compiler) compileExpression(v ast.Expression) compiledExpr ***REMOVED***
	// log.Printf("compileExpression: %T", v)
	switch v := v.(type) ***REMOVED***
	case nil:
		return nil
	case *ast.AssignExpression:
		return c.compileAssignExpression(v)
	case *ast.NumberLiteral:
		return c.compileNumberLiteral(v)
	case *ast.StringLiteral:
		return c.compileStringLiteral(v)
	case *ast.BooleanLiteral:
		return c.compileBooleanLiteral(v)
	case *ast.NullLiteral:
		r := &compiledLiteral***REMOVED***
			val: _null,
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.Identifier:
		return c.compileIdentifierExpression(v)
	case *ast.CallExpression:
		return c.compileCallExpression(v)
	case *ast.ObjectLiteral:
		return c.compileObjectLiteral(v)
	case *ast.ArrayLiteral:
		return c.compileArrayLiteral(v)
	case *ast.RegExpLiteral:
		return c.compileRegexpLiteral(v)
	case *ast.BinaryExpression:
		return c.compileBinaryExpression(v)
	case *ast.UnaryExpression:
		return c.compileUnaryExpression(v)
	case *ast.ConditionalExpression:
		return c.compileConditionalExpression(v)
	case *ast.FunctionLiteral:
		return c.compileFunctionLiteral(v, true)
	case *ast.ArrowFunctionLiteral:
		return c.compileArrowFunctionLiteral(v)
	case *ast.DotExpression:
		r := &compiledDotExpr***REMOVED***
			left: c.compileExpression(v.Left),
			name: v.Identifier.Name,
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.BracketExpression:
		r := &compiledBracketExpr***REMOVED***
			left:   c.compileExpression(v.Left),
			member: c.compileExpression(v.Member),
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.ThisExpression:
		r := &compiledThisExpr***REMOVED******REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.SequenceExpression:
		return c.compileSequenceExpression(v)
	case *ast.NewExpression:
		return c.compileNewExpression(v)
	case *ast.MetaProperty:
		return c.compileMetaProperty(v)
	case *ast.ObjectPattern:
		return c.compileObjectAssignmentPattern(v)
	case *ast.ArrayPattern:
		return c.compileArrayAssignmentPattern(v)
	default:
		panic(fmt.Errorf("Unknown expression type: %T", v))
	***REMOVED***
***REMOVED***

func (e *baseCompiledExpr) constant() bool ***REMOVED***
	return false
***REMOVED***

func (e *baseCompiledExpr) init(c *compiler, idx file.Idx) ***REMOVED***
	e.c = c
	e.offset = int(idx) - 1
***REMOVED***

func (e *baseCompiledExpr) emitSetter(compiledExpr, bool) ***REMOVED***
	e.c.throwSyntaxError(e.offset, "Not a valid left-value expression")
***REMOVED***

func (e *baseCompiledExpr) emitRef() ***REMOVED***
	e.c.throwSyntaxError(e.offset, "Cannot emit reference for this type of expression")
***REMOVED***

func (e *baseCompiledExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &constantExpr***REMOVED***
		val: valueTrue,
	***REMOVED***
	r.init(e.c, file.Idx(e.offset+1))
	return r
***REMOVED***

func (e *baseCompiledExpr) emitUnary(func(), func(), bool, bool) ***REMOVED***
	e.c.throwSyntaxError(e.offset, "Not a valid left-value expression")
***REMOVED***

func (e *baseCompiledExpr) addSrcMap() ***REMOVED***
	if e.offset > 0 ***REMOVED***
		e.c.p.srcMap = append(e.c.p.srcMap, srcMapItem***REMOVED***pc: len(e.c.p.code), srcPos: e.offset***REMOVED***)
	***REMOVED***
***REMOVED***

func (e *constantExpr) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		e.c.emit(loadVal(e.c.p.defineLiteralValue(e.val)))
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b != nil ***REMOVED***
			if putOnStack ***REMOVED***
				b.emitGet()
			***REMOVED*** else ***REMOVED***
				b.emitGetP()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			panic("No dynamics and not found")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGetVar(false)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadDynamic(e.name))
		***REMOVED***
		if !putOnStack ***REMOVED***
			e.c.emit(pop)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitGetterOrRef() ***REMOVED***
	e.addSrcMap()
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGet()
		***REMOVED*** else ***REMOVED***
			panic("No dynamics and not found")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGetVar(false)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadDynamicRef(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitGetterAndCallee() ***REMOVED***
	e.addSrcMap()
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b != nil ***REMOVED***
			e.c.emit(loadUndef)
			b.emitGet()
		***REMOVED*** else ***REMOVED***
			panic("No dynamics and not found")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGetVar(true)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadDynamicCallee(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) emitVarSetter1(name unistring.String, offset int, putOnStack bool, emitRight func(isRef bool)) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierLName(name, offset)
	***REMOVED***

	if b, noDynamics := c.scope.lookupName(name); noDynamics ***REMOVED***
		emitRight(false)
		if b != nil ***REMOVED***
			if putOnStack ***REMOVED***
				b.emitSet()
			***REMOVED*** else ***REMOVED***
				b.emitSetP()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if c.scope.strict ***REMOVED***
				c.emit(setGlobalStrict(name))
			***REMOVED*** else ***REMOVED***
				c.emit(setGlobal(name))
			***REMOVED***
			if !putOnStack ***REMOVED***
				c.emit(pop)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitResolveVar(c.scope.strict)
		***REMOVED*** else ***REMOVED***
			if c.scope.strict ***REMOVED***
				c.emit(resolveVar1Strict(name))
			***REMOVED*** else ***REMOVED***
				c.emit(resolveVar1(name))
			***REMOVED***
		***REMOVED***
		emitRight(true)
		if putOnStack ***REMOVED***
			c.emit(putValue)
		***REMOVED*** else ***REMOVED***
			c.emit(putValueP)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) emitVarSetter(name unistring.String, offset int, valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	c.emitVarSetter1(name, offset, putOnStack, func(bool) ***REMOVED***
		c.emitExpr(valueExpr, true)
	***REMOVED***)
***REMOVED***

func (c *compiler) emitVarRef(name unistring.String, offset int) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierLName(name, offset)
	***REMOVED***

	b, _ := c.scope.lookupName(name)
	if b != nil ***REMOVED***
		b.emitResolveVar(c.scope.strict)
	***REMOVED*** else ***REMOVED***
		if c.scope.strict ***REMOVED***
			c.emit(resolveVar1Strict(name))
		***REMOVED*** else ***REMOVED***
			c.emit(resolveVar1(name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitRef() ***REMOVED***
	e.c.emitVarRef(e.name, e.offset)
***REMOVED***

func (e *compiledIdentifierExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.c.emitVarSetter(e.name, e.offset, valueExpr, putOnStack)
***REMOVED***

func (e *compiledIdentifierExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.c.emitVarSetter1(e.name, e.offset, true, func(isRef bool) ***REMOVED***
			e.c.emit(loadUndef)
			if isRef ***REMOVED***
				e.c.emit(getValue)
			***REMOVED*** else ***REMOVED***
				e.emitGetter(true)
			***REMOVED***
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			if !postfix ***REMOVED***
				body()
			***REMOVED***
			e.c.emit(rdupN(1))
			if postfix ***REMOVED***
				body()
			***REMOVED***
		***REMOVED***)
		e.c.emit(pop)
	***REMOVED*** else ***REMOVED***
		e.c.emitVarSetter1(e.name, e.offset, false, func(isRef bool) ***REMOVED***
			if isRef ***REMOVED***
				e.c.emit(getValue)
			***REMOVED*** else ***REMOVED***
				e.emitGetter(true)
			***REMOVED***
			body()
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) deleteExpr() compiledExpr ***REMOVED***
	if e.c.scope.strict ***REMOVED***
		e.c.throwSyntaxError(e.offset, "Delete of an unqualified identifier in strict mode")
		panic("Unreachable")
	***REMOVED***
	if b, noDynamics := e.c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if b == nil ***REMOVED***
			r := &deleteGlobalExpr***REMOVED***
				name: e.name,
			***REMOVED***
			r.init(e.c, file.Idx(0))
			return r
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b == nil ***REMOVED***
			r := &deleteVarExpr***REMOVED***
				name: e.name,
			***REMOVED***
			r.init(e.c, file.Idx(e.offset+1))
			return r
		***REMOVED***
	***REMOVED***
	r := &compiledLiteral***REMOVED***
		val: valueFalse,
	***REMOVED***
	r.init(e.c, file.Idx(e.offset+1))
	return r
***REMOVED***

type compiledDotExpr struct ***REMOVED***
	baseCompiledExpr
	left compiledExpr
	name unistring.String
***REMOVED***

func (e *compiledDotExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.addSrcMap()
	e.c.emit(getProp(e.name))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) emitRef() ***REMOVED***
	e.left.emitGetter(true)
	if e.c.scope.strict ***REMOVED***
		e.c.emit(getPropRefStrict(e.name))
	***REMOVED*** else ***REMOVED***
		e.c.emit(getPropRef(e.name))
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	valueExpr.emitGetter(true)
	if e.c.scope.strict ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setPropStrict(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropStrictP(e.name))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setProp(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropP(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if !putOnStack ***REMOVED***
		e.left.emitGetter(true)
		e.c.emit(dup)
		e.c.emit(getProp(e.name))
		body()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setPropStrict(e.name), pop)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setProp(e.name), pop)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.left.emitGetter(true)
			e.c.emit(dup)
			e.c.emit(getProp(e.name))
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setPropStrict(e.name))
			***REMOVED*** else ***REMOVED***
				e.c.emit(setProp(e.name))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.left.emitGetter(true)
			e.c.emit(dup)
			e.c.emit(getProp(e.name))
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(2))
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setPropStrict(e.name))
			***REMOVED*** else ***REMOVED***
				e.c.emit(setProp(e.name))
			***REMOVED***
			e.c.emit(pop)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &deletePropExpr***REMOVED***
		left: e.left,
		name: e.name,
	***REMOVED***
	r.init(e.c, file.Idx(0))
	return r
***REMOVED***

func (e *compiledBracketExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.member.emitGetter(true)
	e.addSrcMap()
	e.c.emit(getElem)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledBracketExpr) emitRef() ***REMOVED***
	e.left.emitGetter(true)
	e.member.emitGetter(true)
	if e.c.scope.strict ***REMOVED***
		e.c.emit(getElemRefStrict)
	***REMOVED*** else ***REMOVED***
		e.c.emit(getElemRef)
	***REMOVED***
***REMOVED***

func (e *compiledBracketExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.member.emitGetter(true)
	valueExpr.emitGetter(true)
	if e.c.scope.strict ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setElemStrict)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElemStrictP)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if putOnStack ***REMOVED***
			e.c.emit(setElem)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElemP)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledBracketExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if !putOnStack ***REMOVED***
		e.left.emitGetter(true)
		e.member.emitGetter(true)
		e.c.emit(dupN(1), dupN(1))
		e.c.emit(getElem)
		body()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setElemStrict, pop)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElem, pop)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.left.emitGetter(true)
			e.member.emitGetter(true)
			e.c.emit(dupN(1), dupN(1))
			e.c.emit(getElem)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setElemStrict)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setElem)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.left.emitGetter(true)
			e.member.emitGetter(true)
			e.c.emit(dupN(1), dupN(1))
			e.c.emit(getElem)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(3))
			body()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setElemStrict, pop)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setElem, pop)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledBracketExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &deleteElemExpr***REMOVED***
		left:   e.left,
		member: e.member,
	***REMOVED***
	r.init(e.c, file.Idx(0))
	return r
***REMOVED***

func (e *deleteElemExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.member.emitGetter(true)
	e.addSrcMap()
	if e.c.scope.strict ***REMOVED***
		e.c.emit(deleteElemStrict)
	***REMOVED*** else ***REMOVED***
		e.c.emit(deleteElem)
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *deletePropExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.addSrcMap()
	if e.c.scope.strict ***REMOVED***
		e.c.emit(deletePropStrict(e.name))
	***REMOVED*** else ***REMOVED***
		e.c.emit(deleteProp(e.name))
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *deleteVarExpr) emitGetter(putOnStack bool) ***REMOVED***
	/*if e.c.scope.strict ***REMOVED***
		e.c.throwSyntaxError(e.offset, "Delete of an unqualified identifier in strict mode")
		return
	***REMOVED****/
	e.c.emit(deleteVar(e.name))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *deleteGlobalExpr) emitGetter(putOnStack bool) ***REMOVED***
	/*if e.c.scope.strict ***REMOVED***
		e.c.throwSyntaxError(e.offset, "Delete of an unqualified identifier in strict mode")
		return
	***REMOVED****/

	e.c.emit(deleteGlobal(e.name))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledAssignExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	switch e.operator ***REMOVED***
	case token.ASSIGN:
		if fn, ok := e.right.(*compiledFunctionLiteral); ok ***REMOVED***
			if fn.name == nil ***REMOVED***
				if id, ok := e.left.(*compiledIdentifierExpr); ok ***REMOVED***
					fn.lhsName = id.name
				***REMOVED***
			***REMOVED***
		***REMOVED***
		e.left.emitSetter(e.right, putOnStack)
	case token.PLUS:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(add)
		***REMOVED***, false, putOnStack)
	case token.MINUS:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(sub)
		***REMOVED***, false, putOnStack)
	case token.MULTIPLY:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(mul)
		***REMOVED***, false, putOnStack)
	case token.SLASH:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(div)
		***REMOVED***, false, putOnStack)
	case token.REMAINDER:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(mod)
		***REMOVED***, false, putOnStack)
	case token.OR:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(or)
		***REMOVED***, false, putOnStack)
	case token.AND:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(and)
		***REMOVED***, false, putOnStack)
	case token.EXCLUSIVE_OR:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(xor)
		***REMOVED***, false, putOnStack)
	case token.SHIFT_LEFT:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(sal)
		***REMOVED***, false, putOnStack)
	case token.SHIFT_RIGHT:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(sar)
		***REMOVED***, false, putOnStack)
	case token.UNSIGNED_SHIFT_RIGHT:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(shr)
		***REMOVED***, false, putOnStack)
	default:
		panic(fmt.Errorf("Unknown assign operator: %s", e.operator.String()))
	***REMOVED***
***REMOVED***

func (e *compiledLiteral) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		e.c.emit(loadVal(e.c.p.defineLiteralValue(e.val)))
	***REMOVED***
***REMOVED***

func (e *compiledLiteral) constant() bool ***REMOVED***
	return true
***REMOVED***

func (c *compiler) compileParameterBindingIdentifier(name unistring.String, offset int) (*binding, bool) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierName(name, offset)
		c.checkIdentifierLName(name, offset)
	***REMOVED***
	return c.scope.bindNameShadow(name)
***REMOVED***

func (c *compiler) compileParameterPatternIdBinding(name unistring.String, offset int) ***REMOVED***
	if _, unique := c.compileParameterBindingIdentifier(name, offset); !unique ***REMOVED***
		c.throwSyntaxError(offset, "Duplicate parameter name not allowed in this context")
	***REMOVED***
***REMOVED***

func (c *compiler) compileParameterPatternBinding(item ast.Expression) ***REMOVED***
	c.createBindings(item, c.compileParameterPatternIdBinding)
***REMOVED***

func (e *compiledFunctionLiteral) emitGetter(putOnStack bool) ***REMOVED***
	savedPrg := e.c.p
	e.c.p = &Program***REMOVED***
		src: e.c.p.src,
	***REMOVED***
	e.c.newScope()
	s := e.c.scope
	s.function = true
	s.arrow = e.isArrow

	var name unistring.String
	if e.name != nil ***REMOVED***
		name = e.name.Name
	***REMOVED*** else ***REMOVED***
		name = e.lhsName
	***REMOVED***

	if name != "" ***REMOVED***
		e.c.p.funcName = name
	***REMOVED***
	savedBlock := e.c.block
	defer func() ***REMOVED***
		e.c.block = savedBlock
	***REMOVED***()

	e.c.block = &block***REMOVED***
		typ: blockScope,
	***REMOVED***

	if !s.strict ***REMOVED***
		s.strict = e.strict != nil
	***REMOVED***

	hasPatterns := false
	hasInits := false
	firstDupIdx := -1
	length := 0

	if e.parameterList.Rest != nil ***REMOVED***
		hasPatterns = true // strictly speaking not, but we need to activate all the checks
	***REMOVED***

	// First, make sure that the first bindings correspond to the formal parameters
	for _, item := range e.parameterList.List ***REMOVED***
		switch tgt := item.Target.(type) ***REMOVED***
		case *ast.Identifier:
			offset := int(tgt.Idx) - 1
			b, unique := e.c.compileParameterBindingIdentifier(tgt.Name, offset)
			if !unique ***REMOVED***
				firstDupIdx = offset
			***REMOVED***
			b.isArg = true
		case ast.Pattern:
			b := s.addBinding(int(item.Idx0()) - 1)
			b.isArg = true
			hasPatterns = true
		default:
			e.c.throwSyntaxError(int(item.Idx0())-1, "Unsupported BindingElement type: %T", item)
			return
		***REMOVED***
		if item.Initializer != nil ***REMOVED***
			hasInits = true
		***REMOVED***

		if firstDupIdx >= 0 && (hasPatterns || hasInits || s.strict || e.isArrow) ***REMOVED***
			e.c.throwSyntaxError(firstDupIdx, "Duplicate parameter name not allowed in this context")
			return
		***REMOVED***

		if (hasPatterns || hasInits) && e.strict != nil ***REMOVED***
			e.c.throwSyntaxError(int(e.strict.Idx)-1, "Illegal 'use strict' directive in function with non-simple parameter list")
			return
		***REMOVED***

		if !hasInits ***REMOVED***
			length++
		***REMOVED***
	***REMOVED***

	// create pattern bindings
	if hasPatterns ***REMOVED***
		for _, item := range e.parameterList.List ***REMOVED***
			switch tgt := item.Target.(type) ***REMOVED***
			case *ast.Identifier:
				// we already created those in the previous loop, skipping
			default:
				e.c.compileParameterPatternBinding(tgt)
			***REMOVED***
		***REMOVED***
		if rest := e.parameterList.Rest; rest != nil ***REMOVED***
			e.c.compileParameterPatternBinding(rest)
		***REMOVED***
	***REMOVED***

	paramsCount := len(e.parameterList.List)

	s.numArgs = paramsCount
	body := e.body
	funcs := e.c.extractFunctions(body)
	var calleeBinding *binding
	preambleLen := 4 // enter, boxThis, createArgs, set
	e.c.p.code = make([]instruction, preambleLen, 8)

	emitArgsRestMark := -1
	firstForwardRef := -1
	enterFunc2Mark := -1

	if hasPatterns || hasInits ***REMOVED***
		if e.isExpr && e.name != nil ***REMOVED***
			if b, created := s.bindNameLexical(e.name.Name, false, 0); created ***REMOVED***
				b.isConst = true
				calleeBinding = b
			***REMOVED***
		***REMOVED***
		if calleeBinding != nil ***REMOVED***
			e.c.emit(loadCallee)
			calleeBinding.emitInit()
		***REMOVED***
		for i, item := range e.parameterList.List ***REMOVED***
			if pattern, ok := item.Target.(ast.Pattern); ok ***REMOVED***
				i := i
				e.c.compilePatternInitExpr(func() ***REMOVED***
					if firstForwardRef == -1 ***REMOVED***
						s.bindings[i].emitGet()
					***REMOVED*** else ***REMOVED***
						e.c.emit(loadStackLex(-i - 1))
					***REMOVED***
				***REMOVED***, item.Initializer, item.Target.Idx0()).emitGetter(true)
				e.c.emitPattern(pattern, func(target, init compiledExpr) ***REMOVED***
					e.c.emitPatternLexicalAssign(target, init, false)
				***REMOVED***, false)
			***REMOVED*** else if item.Initializer != nil ***REMOVED***
				markGet := len(e.c.p.code)
				e.c.emit(nil)
				mark := len(e.c.p.code)
				e.c.emit(nil)
				e.c.compileExpression(item.Initializer).emitGetter(true)
				if firstForwardRef == -1 && (s.isDynamic() || s.bindings[i].useCount() > 0) ***REMOVED***
					firstForwardRef = i
				***REMOVED***
				if firstForwardRef == -1 ***REMOVED***
					s.bindings[i].emitGetAt(markGet)
				***REMOVED*** else ***REMOVED***
					e.c.p.code[markGet] = loadStackLex(-i - 1)
				***REMOVED***
				s.bindings[i].emitInit()
				e.c.p.code[mark] = jdefP(len(e.c.p.code) - mark)
			***REMOVED*** else ***REMOVED***
				if firstForwardRef == -1 && s.bindings[i].useCount() > 0 ***REMOVED***
					firstForwardRef = i
				***REMOVED***
				if firstForwardRef != -1 ***REMOVED***
					e.c.emit(loadStackLex(-i - 1))
					s.bindings[i].emitInit()
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if rest := e.parameterList.Rest; rest != nil ***REMOVED***
			e.c.emitAssign(rest, e.c.compileEmitterExpr(
				func() ***REMOVED***
					emitArgsRestMark = len(e.c.p.code)
					e.c.emit(createArgsRestStack(paramsCount))
				***REMOVED***, rest.Idx0()),
				func(target, init compiledExpr) ***REMOVED***
					e.c.emitPatternLexicalAssign(target, init, false)
				***REMOVED***)
		***REMOVED***
		if firstForwardRef != -1 ***REMOVED***
			for _, b := range s.bindings ***REMOVED***
				b.inStash = true
			***REMOVED***
			s.argsInStash = true
			s.needStash = true
		***REMOVED***

		e.c.newBlockScope()
		varScope := e.c.scope
		varScope.variable = true
		enterFunc2Mark = len(e.c.p.code)
		e.c.emit(nil)
		e.c.compileDeclList(e.declarationList, false)
		e.c.createFunctionBindings(funcs)
		e.c.compileLexicalDeclarationsFuncBody(body, calleeBinding)
		for _, b := range varScope.bindings ***REMOVED***
			if b.isVar ***REMOVED***
				if parentBinding := s.boundNames[b.name]; parentBinding != nil && parentBinding != calleeBinding ***REMOVED***
					parentBinding.emitGet()
					b.emitSetP()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// To avoid triggering variable conflict when binding from non-strict direct eval().
		// Parameters are supposed to be in a parent scope, hence no conflict.
		for _, b := range s.bindings[:paramsCount] ***REMOVED***
			b.isVar = true
		***REMOVED***
		e.c.compileDeclList(e.declarationList, true)
		e.c.createFunctionBindings(funcs)
		e.c.compileLexicalDeclarations(body, true)
		if e.isExpr && e.name != nil ***REMOVED***
			if b, created := s.bindNameLexical(e.name.Name, false, 0); created ***REMOVED***
				b.isConst = true
				calleeBinding = b
			***REMOVED***
		***REMOVED***
		if calleeBinding != nil ***REMOVED***
			e.c.emit(loadCallee)
			calleeBinding.emitInit()
		***REMOVED***
	***REMOVED***

	e.c.compileFunctions(funcs)
	e.c.compileStatements(body, false)

	var last ast.Statement
	if l := len(body); l > 0 ***REMOVED***
		last = body[l-1]
	***REMOVED***
	if _, ok := last.(*ast.ReturnStatement); !ok ***REMOVED***
		e.c.emit(loadUndef, ret)
	***REMOVED***

	delta := 0
	code := e.c.p.code

	if calleeBinding != nil && !s.isDynamic() && calleeBinding.useCount() == 1 ***REMOVED***
		s.deleteBinding(calleeBinding)
		preambleLen += 2
	***REMOVED***

	if !s.argsInStash && (s.argsNeeded || s.isDynamic()) ***REMOVED***
		s.moveArgsToStash()
	***REMOVED***

	if s.argsNeeded ***REMOVED***
		b, created := s.bindNameLexical("arguments", false, 0)
		if !created && !b.isVar ***REMOVED***
			s.argsNeeded = false
		***REMOVED*** else ***REMOVED***
			if s.strict ***REMOVED***
				b.isConst = true
			***REMOVED*** else ***REMOVED***
				b.isVar = e.c.scope.function
			***REMOVED***
			pos := preambleLen - 2
			delta += 2
			if s.strict || hasPatterns || hasInits ***REMOVED***
				code[pos] = createArgsUnmapped(paramsCount)
			***REMOVED*** else ***REMOVED***
				code[pos] = createArgsMapped(paramsCount)
			***REMOVED***
			pos++
			b.markAccessPointAtScope(s, pos)
			code[pos] = storeStashP(0)
		***REMOVED***
	***REMOVED***

	stashSize, stackSize := s.finaliseVarAlloc(0)

	if !s.strict && s.thisNeeded ***REMOVED***
		delta++
		code[preambleLen-delta] = boxThis
	***REMOVED***
	delta++
	delta = preambleLen - delta
	var enter instruction
	if stashSize > 0 || s.argsInStash ***REMOVED***
		if firstForwardRef == -1 ***REMOVED***
			enter1 := enterFunc***REMOVED***
				numArgs:     uint32(paramsCount),
				argsToStash: s.argsInStash,
				stashSize:   uint32(stashSize),
				stackSize:   uint32(stackSize),
				extensible:  s.dynamic,
			***REMOVED***
			if s.isDynamic() ***REMOVED***
				enter1.names = s.makeNamesMap()
			***REMOVED***
			enter = &enter1
			if enterFunc2Mark != -1 ***REMOVED***
				ef2 := &enterFuncBody***REMOVED***
					extensible: e.c.scope.dynamic,
				***REMOVED***
				e.c.updateEnterBlock(&ef2.enterBlock)
				e.c.p.code[enterFunc2Mark] = ef2
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			enter1 := enterFunc1***REMOVED***
				stashSize:  uint32(stashSize),
				numArgs:    uint32(paramsCount),
				argsToCopy: uint32(firstForwardRef),
				extensible: s.dynamic,
			***REMOVED***
			if s.isDynamic() ***REMOVED***
				enter1.names = s.makeNamesMap()
			***REMOVED***
			enter = &enter1
			if enterFunc2Mark != -1 ***REMOVED***
				ef2 := &enterFuncBody***REMOVED***
					adjustStack: true,
					extensible:  e.c.scope.dynamic,
				***REMOVED***
				e.c.updateEnterBlock(&ef2.enterBlock)
				e.c.p.code[enterFunc2Mark] = ef2
			***REMOVED***
		***REMOVED***
		if emitArgsRestMark != -1 ***REMOVED***
			e.c.p.code[emitArgsRestMark] = createArgsRestStash
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		enter = &enterFuncStashless***REMOVED***
			stackSize: uint32(stackSize),
			args:      uint32(paramsCount),
		***REMOVED***
		if enterFunc2Mark != -1 ***REMOVED***
			ef2 := &enterFuncBody***REMOVED***
				extensible: e.c.scope.dynamic,
			***REMOVED***
			e.c.updateEnterBlock(&ef2.enterBlock)
			e.c.p.code[enterFunc2Mark] = ef2
		***REMOVED***
	***REMOVED***
	code[delta] = enter
	if delta != 0 ***REMOVED***
		e.c.p.code = code[delta:]
		for i := range e.c.p.srcMap ***REMOVED***
			e.c.p.srcMap[i].pc -= delta
		***REMOVED***
		s.adjustBase(-delta)
	***REMOVED***

	strict := s.strict
	p := e.c.p
	// e.c.p.dumpCode()
	if enterFunc2Mark != -1 ***REMOVED***
		e.c.popScope()
	***REMOVED***
	e.c.popScope()
	e.c.p = savedPrg
	if e.isArrow ***REMOVED***
		e.c.emit(&newArrowFunc***REMOVED***newFunc: newFunc***REMOVED***prg: p, length: uint32(length), name: name, source: e.source, strict: strict***REMOVED******REMOVED***)
	***REMOVED*** else ***REMOVED***
		e.c.emit(&newFunc***REMOVED***prg: p, length: uint32(length), name: name, source: e.source, strict: strict***REMOVED***)
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunctionLiteral(v *ast.FunctionLiteral, isExpr bool) *compiledFunctionLiteral ***REMOVED***
	strictBody := c.isStrictStatement(v.Body)
	if v.Name != nil && (c.scope.strict || strictBody != nil) ***REMOVED***
		c.checkIdentifierLName(v.Name.Name, int(v.Name.Idx)-1)
	***REMOVED***
	r := &compiledFunctionLiteral***REMOVED***
		name:            v.Name,
		parameterList:   v.ParameterList,
		body:            v.Body.List,
		source:          v.Source,
		declarationList: v.DeclarationList,
		isExpr:          isExpr,
		strict:          strictBody,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileArrowFunctionLiteral(v *ast.ArrowFunctionLiteral) *compiledFunctionLiteral ***REMOVED***
	var strictBody *ast.StringLiteral
	var body []ast.Statement
	switch b := v.Body.(type) ***REMOVED***
	case *ast.BlockStatement:
		strictBody = c.isStrictStatement(b)
		body = b.List
	case *ast.ExpressionBody:
		body = []ast.Statement***REMOVED***
			&ast.ReturnStatement***REMOVED***
				Argument: b.Expression,
			***REMOVED***,
		***REMOVED***
	default:
		c.throwSyntaxError(int(b.Idx0())-1, "Unsupported ConciseBody type: %T", b)
	***REMOVED***
	r := &compiledFunctionLiteral***REMOVED***
		parameterList:   v.ParameterList,
		body:            body,
		source:          v.Source,
		declarationList: v.DeclarationList,
		isExpr:          true,
		isArrow:         true,
		strict:          strictBody,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledThisExpr) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		scope := e.c.scope
		for ; scope != nil && !scope.function && !scope.eval; scope = scope.outer ***REMOVED***
		***REMOVED***

		if scope != nil ***REMOVED***
			if !scope.arrow ***REMOVED***
				scope.thisNeeded = true
			***REMOVED***
			e.c.emit(loadStack(0))
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadGlobalObject)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledNewExpr) emitGetter(putOnStack bool) ***REMOVED***
	if e.isVariadic ***REMOVED***
		e.c.emit(startVariadic)
	***REMOVED***
	e.callee.emitGetter(true)
	for _, expr := range e.args ***REMOVED***
		expr.emitGetter(true)
	***REMOVED***
	e.addSrcMap()
	if e.isVariadic ***REMOVED***
		e.c.emit(newVariadic, endVariadic)
	***REMOVED*** else ***REMOVED***
		e.c.emit(_new(len(e.args)))
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileCallArgs(list []ast.Expression) (args []compiledExpr, isVariadic bool) ***REMOVED***
	args = make([]compiledExpr, len(list))
	for i, argExpr := range list ***REMOVED***
		if spread, ok := argExpr.(*ast.SpreadElement); ok ***REMOVED***
			args[i] = c.compileSpreadCallArgument(spread)
			isVariadic = true
		***REMOVED*** else ***REMOVED***
			args[i] = c.compileExpression(argExpr)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (c *compiler) compileNewExpression(v *ast.NewExpression) compiledExpr ***REMOVED***
	args, isVariadic := c.compileCallArgs(v.ArgumentList)
	r := &compiledNewExpr***REMOVED***
		compiledCallExpr: compiledCallExpr***REMOVED***
			callee:     c.compileExpression(v.Callee),
			args:       args,
			isVariadic: isVariadic,
		***REMOVED***,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledNewTarget) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.addSrcMap()
		e.c.emit(loadNewTarget)
	***REMOVED***
***REMOVED***

func (c *compiler) compileMetaProperty(v *ast.MetaProperty) compiledExpr ***REMOVED***
	if v.Meta.Name == "new" || v.Property.Name != "target" ***REMOVED***
		r := &compiledNewTarget***REMOVED******REMOVED***
		r.init(c, v.Idx0())
		return r
	***REMOVED***
	c.throwSyntaxError(int(v.Idx)-1, "Unsupported meta property: %s.%s", v.Meta.Name, v.Property.Name)
	return nil
***REMOVED***

func (e *compiledSequenceExpr) emitGetter(putOnStack bool) ***REMOVED***
	if len(e.sequence) > 0 ***REMOVED***
		for i := 0; i < len(e.sequence)-1; i++ ***REMOVED***
			e.sequence[i].emitGetter(false)
		***REMOVED***
		e.sequence[len(e.sequence)-1].emitGetter(putOnStack)
	***REMOVED***
***REMOVED***

func (c *compiler) compileSequenceExpression(v *ast.SequenceExpression) compiledExpr ***REMOVED***
	s := make([]compiledExpr, len(v.Sequence))
	for i, expr := range v.Sequence ***REMOVED***
		s[i] = c.compileExpression(expr)
	***REMOVED***
	r := &compiledSequenceExpr***REMOVED***
		sequence: s,
	***REMOVED***
	var idx file.Idx
	if len(v.Sequence) > 0 ***REMOVED***
		idx = v.Idx0()
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (c *compiler) emitThrow(v Value) ***REMOVED***
	if o, ok := v.(*Object); ok ***REMOVED***
		t := nilSafe(o.self.getStr("name", nil)).toString().String()
		switch t ***REMOVED***
		case "TypeError":
			c.emit(loadDynamic(t))
			msg := o.self.getStr("message", nil)
			if msg != nil ***REMOVED***
				c.emit(loadVal(c.p.defineLiteralValue(msg)))
				c.emit(_new(1))
			***REMOVED*** else ***REMOVED***
				c.emit(_new(0))
			***REMOVED***
			c.emit(throw)
			return
		***REMOVED***
	***REMOVED***
	panic(fmt.Errorf("unknown exception type thrown while evaliating constant expression: %s", v.String()))
***REMOVED***

func (c *compiler) emitConst(expr compiledExpr, putOnStack bool) ***REMOVED***
	v, ex := c.evalConst(expr)
	if ex == nil ***REMOVED***
		if putOnStack ***REMOVED***
			c.emit(loadVal(c.p.defineLiteralValue(v)))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.emitThrow(ex.val)
	***REMOVED***
***REMOVED***

func (c *compiler) emitExpr(expr compiledExpr, putOnStack bool) ***REMOVED***
	if expr.constant() ***REMOVED***
		c.emitConst(expr, putOnStack)
	***REMOVED*** else ***REMOVED***
		expr.emitGetter(putOnStack)
	***REMOVED***
***REMOVED***

func (c *compiler) evalConst(expr compiledExpr) (Value, *Exception) ***REMOVED***
	if expr, ok := expr.(*compiledLiteral); ok ***REMOVED***
		return expr.val, nil
	***REMOVED***
	if c.evalVM == nil ***REMOVED***
		c.evalVM = New().vm
	***REMOVED***
	var savedPrg *Program
	createdPrg := false
	if c.evalVM.prg == nil ***REMOVED***
		c.evalVM.prg = &Program***REMOVED******REMOVED***
		savedPrg = c.p
		c.p = c.evalVM.prg
		createdPrg = true
	***REMOVED***
	savedPc := len(c.p.code)
	expr.emitGetter(true)
	c.emit(halt)
	c.evalVM.pc = savedPc
	ex := c.evalVM.runTry()
	if createdPrg ***REMOVED***
		c.evalVM.prg = nil
		c.evalVM.pc = 0
		c.p = savedPrg
	***REMOVED*** else ***REMOVED***
		c.evalVM.prg.code = c.evalVM.prg.code[:savedPc]
		c.p.code = c.evalVM.prg.code
	***REMOVED***
	if ex == nil ***REMOVED***
		return c.evalVM.pop(), nil
	***REMOVED***
	return nil, ex
***REMOVED***

func (e *compiledUnaryExpr) constant() bool ***REMOVED***
	return e.operand.constant()
***REMOVED***

func (e *compiledUnaryExpr) emitGetter(putOnStack bool) ***REMOVED***
	var prepare, body func()

	toNumber := func() ***REMOVED***
		e.c.emit(toNumber)
	***REMOVED***

	switch e.operator ***REMOVED***
	case token.NOT:
		e.operand.emitGetter(true)
		e.c.emit(not)
		goto end
	case token.BITWISE_NOT:
		e.operand.emitGetter(true)
		e.c.emit(bnot)
		goto end
	case token.TYPEOF:
		if o, ok := e.operand.(compiledExprOrRef); ok ***REMOVED***
			o.emitGetterOrRef()
		***REMOVED*** else ***REMOVED***
			e.operand.emitGetter(true)
		***REMOVED***
		e.c.emit(typeof)
		goto end
	case token.DELETE:
		e.operand.deleteExpr().emitGetter(putOnStack)
		return
	case token.MINUS:
		e.c.emitExpr(e.operand, true)
		e.c.emit(neg)
		goto end
	case token.PLUS:
		e.c.emitExpr(e.operand, true)
		e.c.emit(plus)
		goto end
	case token.INCREMENT:
		prepare = toNumber
		body = func() ***REMOVED***
			e.c.emit(inc)
		***REMOVED***
	case token.DECREMENT:
		prepare = toNumber
		body = func() ***REMOVED***
			e.c.emit(dec)
		***REMOVED***
	case token.VOID:
		e.c.emitExpr(e.operand, false)
		if putOnStack ***REMOVED***
			e.c.emit(loadUndef)
		***REMOVED***
		return
	default:
		panic(fmt.Errorf("Unknown unary operator: %s", e.operator.String()))
	***REMOVED***

	e.operand.emitUnary(prepare, body, e.postfix, putOnStack)
	return

end:
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileUnaryExpression(v *ast.UnaryExpression) compiledExpr ***REMOVED***
	r := &compiledUnaryExpr***REMOVED***
		operand:  c.compileExpression(v.Operand),
		operator: v.Operator,
		postfix:  v.Postfix,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledConditionalExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.test.emitGetter(true)
	j := len(e.c.p.code)
	e.c.emit(nil)
	e.consequent.emitGetter(putOnStack)
	j1 := len(e.c.p.code)
	e.c.emit(nil)
	e.c.p.code[j] = jne(len(e.c.p.code) - j)
	e.alternate.emitGetter(putOnStack)
	e.c.p.code[j1] = jump(len(e.c.p.code) - j1)
***REMOVED***

func (c *compiler) compileConditionalExpression(v *ast.ConditionalExpression) compiledExpr ***REMOVED***
	r := &compiledConditionalExpr***REMOVED***
		test:       c.compileExpression(v.Test),
		consequent: c.compileExpression(v.Consequent),
		alternate:  c.compileExpression(v.Alternate),
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledLogicalOr) constant() bool ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if v.ToBoolean() ***REMOVED***
				return true
			***REMOVED***
			return e.right.constant()
		***REMOVED*** else ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (e *compiledLogicalOr) emitGetter(putOnStack bool) ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if !v.ToBoolean() ***REMOVED***
				e.c.emitExpr(e.right, putOnStack)
			***REMOVED*** else ***REMOVED***
				if putOnStack ***REMOVED***
					e.c.emit(loadVal(e.c.p.defineLiteralValue(v)))
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emitThrow(ex.val)
		***REMOVED***
		return
	***REMOVED***
	e.c.emitExpr(e.left, true)
	j := len(e.c.p.code)
	e.addSrcMap()
	e.c.emit(nil)
	e.c.emit(pop)
	e.c.emitExpr(e.right, true)
	e.c.p.code[j] = jeq1(len(e.c.p.code) - j)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledLogicalAnd) constant() bool ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if !v.ToBoolean() ***REMOVED***
				return true
			***REMOVED*** else ***REMOVED***
				return e.right.constant()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (e *compiledLogicalAnd) emitGetter(putOnStack bool) ***REMOVED***
	var j int
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if !v.ToBoolean() ***REMOVED***
				e.c.emit(loadVal(e.c.p.defineLiteralValue(v)))
			***REMOVED*** else ***REMOVED***
				e.c.emitExpr(e.right, putOnStack)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emitThrow(ex.val)
		***REMOVED***
		return
	***REMOVED***
	e.left.emitGetter(true)
	j = len(e.c.p.code)
	e.addSrcMap()
	e.c.emit(nil)
	e.c.emit(pop)
	e.c.emitExpr(e.right, true)
	e.c.p.code[j] = jneq1(len(e.c.p.code) - j)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledBinaryExpr) constant() bool ***REMOVED***
	return e.left.constant() && e.right.constant()
***REMOVED***

func (e *compiledBinaryExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.c.emitExpr(e.left, true)
	e.c.emitExpr(e.right, true)
	e.addSrcMap()

	switch e.operator ***REMOVED***
	case token.LESS:
		e.c.emit(op_lt)
	case token.GREATER:
		e.c.emit(op_gt)
	case token.LESS_OR_EQUAL:
		e.c.emit(op_lte)
	case token.GREATER_OR_EQUAL:
		e.c.emit(op_gte)
	case token.EQUAL:
		e.c.emit(op_eq)
	case token.NOT_EQUAL:
		e.c.emit(op_neq)
	case token.STRICT_EQUAL:
		e.c.emit(op_strict_eq)
	case token.STRICT_NOT_EQUAL:
		e.c.emit(op_strict_neq)
	case token.PLUS:
		e.c.emit(add)
	case token.MINUS:
		e.c.emit(sub)
	case token.MULTIPLY:
		e.c.emit(mul)
	case token.SLASH:
		e.c.emit(div)
	case token.REMAINDER:
		e.c.emit(mod)
	case token.AND:
		e.c.emit(and)
	case token.OR:
		e.c.emit(or)
	case token.EXCLUSIVE_OR:
		e.c.emit(xor)
	case token.INSTANCEOF:
		e.c.emit(op_instanceof)
	case token.IN:
		e.c.emit(op_in)
	case token.SHIFT_LEFT:
		e.c.emit(sal)
	case token.SHIFT_RIGHT:
		e.c.emit(sar)
	case token.UNSIGNED_SHIFT_RIGHT:
		e.c.emit(shr)
	default:
		panic(fmt.Errorf("Unknown operator: %s", e.operator.String()))
	***REMOVED***

	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileBinaryExpression(v *ast.BinaryExpression) compiledExpr ***REMOVED***

	switch v.Operator ***REMOVED***
	case token.LOGICAL_OR:
		return c.compileLogicalOr(v.Left, v.Right, v.Idx0())
	case token.LOGICAL_AND:
		return c.compileLogicalAnd(v.Left, v.Right, v.Idx0())
	***REMOVED***

	r := &compiledBinaryExpr***REMOVED***
		left:     c.compileExpression(v.Left),
		right:    c.compileExpression(v.Right),
		operator: v.Operator,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileLogicalOr(left, right ast.Expression, idx file.Idx) compiledExpr ***REMOVED***
	r := &compiledLogicalOr***REMOVED***
		left:  c.compileExpression(left),
		right: c.compileExpression(right),
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (c *compiler) compileLogicalAnd(left, right ast.Expression, idx file.Idx) compiledExpr ***REMOVED***
	r := &compiledLogicalAnd***REMOVED***
		left:  c.compileExpression(left),
		right: c.compileExpression(right),
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (e *compiledObjectLiteral) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	e.c.emit(newObject)
	for _, prop := range e.expr.Value ***REMOVED***
		switch prop := prop.(type) ***REMOVED***
		case *ast.PropertyKeyed:
			keyExpr := e.c.compileExpression(prop.Key)
			computed := false
			var key unistring.String
			switch keyExpr := keyExpr.(type) ***REMOVED***
			case *compiledLiteral:
				key = keyExpr.val.string()
			default:
				keyExpr.emitGetter(true)
				computed = true
			***REMOVED***
			valueExpr := e.c.compileExpression(prop.Value)
			var anonFn *compiledFunctionLiteral
			if fn, ok := valueExpr.(*compiledFunctionLiteral); ok ***REMOVED***
				if fn.name == nil ***REMOVED***
					anonFn = fn
					fn.lhsName = key
				***REMOVED***
			***REMOVED***
			if computed ***REMOVED***
				e.c.emit(_toPropertyKey***REMOVED******REMOVED***)
				valueExpr.emitGetter(true)
				switch prop.Kind ***REMOVED***
				case ast.PropertyKindValue, ast.PropertyKindMethod:
					if anonFn != nil ***REMOVED***
						e.c.emit(setElem1Named)
					***REMOVED*** else ***REMOVED***
						e.c.emit(setElem1)
					***REMOVED***
				case ast.PropertyKindGet:
					e.c.emit(setPropGetter1)
				case ast.PropertyKindSet:
					e.c.emit(setPropSetter1)
				default:
					panic(fmt.Errorf("unknown property kind: %s", prop.Kind))
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if anonFn != nil ***REMOVED***
					anonFn.lhsName = key
				***REMOVED***
				valueExpr.emitGetter(true)
				switch prop.Kind ***REMOVED***
				case ast.PropertyKindValue:
					if key == __proto__ ***REMOVED***
						e.c.emit(setProto)
					***REMOVED*** else ***REMOVED***
						e.c.emit(setProp1(key))
					***REMOVED***
				case ast.PropertyKindMethod:
					e.c.emit(setProp1(key))
				case ast.PropertyKindGet:
					e.c.emit(setPropGetter(key))
				case ast.PropertyKindSet:
					e.c.emit(setPropSetter(key))
				default:
					panic(fmt.Errorf("unknown property kind: %s", prop.Kind))
				***REMOVED***
			***REMOVED***
		case *ast.PropertyShort:
			key := prop.Name.Name
			if prop.Initializer != nil ***REMOVED***
				e.c.throwSyntaxError(int(prop.Initializer.Idx0())-1, "Invalid shorthand property initializer")
			***REMOVED***
			if e.c.scope.strict && key == "let" ***REMOVED***
				e.c.throwSyntaxError(e.offset, "'let' cannot be used as a shorthand property in strict mode")
			***REMOVED***
			e.c.compileIdentifierExpression(&prop.Name).emitGetter(true)
			e.c.emit(setProp1(key))
		case *ast.SpreadElement:
			e.c.compileExpression(prop.Expression).emitGetter(true)
			e.c.emit(copySpread)
		default:
			panic(fmt.Errorf("unknown Property type: %T", prop))
		***REMOVED***
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileObjectLiteral(v *ast.ObjectLiteral) compiledExpr ***REMOVED***
	r := &compiledObjectLiteral***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledArrayLiteral) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	hasSpread := false
	mark := len(e.c.p.code)
	e.c.emit(nil)
	for _, v := range e.expr.Value ***REMOVED***
		if spread, ok := v.(*ast.SpreadElement); ok ***REMOVED***
			hasSpread = true
			e.c.compileExpression(spread.Expression).emitGetter(true)
			e.c.emit(pushArraySpread)
		***REMOVED*** else ***REMOVED***
			if v != nil ***REMOVED***
				e.c.compileExpression(v).emitGetter(true)
			***REMOVED*** else ***REMOVED***
				e.c.emit(loadNil)
			***REMOVED***
			e.c.emit(pushArrayItem)
		***REMOVED***
	***REMOVED***
	var objCount uint32
	if !hasSpread ***REMOVED***
		objCount = uint32(len(e.expr.Value))
	***REMOVED***
	e.c.p.code[mark] = newArray(objCount)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileArrayLiteral(v *ast.ArrayLiteral) compiledExpr ***REMOVED***
	r := &compiledArrayLiteral***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledRegexpLiteral) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		pattern, err := compileRegexp(e.expr.Pattern, e.expr.Flags)
		if err != nil ***REMOVED***
			e.c.throwSyntaxError(e.offset, err.Error())
		***REMOVED***

		e.c.emit(&newRegexp***REMOVED***pattern: pattern, src: newStringValue(e.expr.Pattern)***REMOVED***)
	***REMOVED***
***REMOVED***

func (c *compiler) compileRegexpLiteral(v *ast.RegExpLiteral) compiledExpr ***REMOVED***
	r := &compiledRegexpLiteral***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledCallExpr) emitGetter(putOnStack bool) ***REMOVED***
	var calleeName unistring.String
	if e.isVariadic ***REMOVED***
		e.c.emit(startVariadic)
	***REMOVED***
	switch callee := e.callee.(type) ***REMOVED***
	case *compiledDotExpr:
		callee.left.emitGetter(true)
		e.c.emit(dup)
		e.c.emit(getPropCallee(callee.name))
	case *compiledBracketExpr:
		callee.left.emitGetter(true)
		e.c.emit(dup)
		callee.member.emitGetter(true)
		e.c.emit(getElemCallee)
	case *compiledIdentifierExpr:
		calleeName = callee.name
		callee.emitGetterAndCallee()
	default:
		e.c.emit(loadUndef)
		callee.emitGetter(true)
	***REMOVED***

	for _, expr := range e.args ***REMOVED***
		expr.emitGetter(true)
	***REMOVED***

	e.addSrcMap()
	if calleeName == "eval" ***REMOVED***
		foundFunc, foundVar := false, false
		for sc := e.c.scope; sc != nil; sc = sc.outer ***REMOVED***
			if !foundFunc && sc.function && !sc.arrow ***REMOVED***
				foundFunc = true
				sc.thisNeeded, sc.argsNeeded = true, true
			***REMOVED***
			if !foundVar && (sc.variable || sc.function) ***REMOVED***
				foundVar = true
				if !sc.strict ***REMOVED***
					sc.dynamic = true
				***REMOVED***
			***REMOVED***
			sc.dynLookup = true
		***REMOVED***

		if e.c.scope.strict ***REMOVED***
			if e.isVariadic ***REMOVED***
				e.c.emit(callEvalVariadicStrict)
			***REMOVED*** else ***REMOVED***
				e.c.emit(callEvalStrict(len(e.args)))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if e.isVariadic ***REMOVED***
				e.c.emit(callEvalVariadic)
			***REMOVED*** else ***REMOVED***
				e.c.emit(callEval(len(e.args)))
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if e.isVariadic ***REMOVED***
			e.c.emit(callVariadic)
		***REMOVED*** else ***REMOVED***
			e.c.emit(call(len(e.args)))
		***REMOVED***
	***REMOVED***
	if e.isVariadic ***REMOVED***
		e.c.emit(endVariadic)
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledCallExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &defaultDeleteExpr***REMOVED***
		expr: e,
	***REMOVED***
	r.init(e.c, file.Idx(e.offset+1))
	return r
***REMOVED***

func (c *compiler) compileSpreadCallArgument(spread *ast.SpreadElement) compiledExpr ***REMOVED***
	r := &compiledSpreadCallArgument***REMOVED***
		expr: c.compileExpression(spread.Expression),
	***REMOVED***
	r.init(c, spread.Idx0())
	return r
***REMOVED***

func (c *compiler) compileCallExpression(v *ast.CallExpression) compiledExpr ***REMOVED***

	args := make([]compiledExpr, len(v.ArgumentList))
	isVariadic := false
	for i, argExpr := range v.ArgumentList ***REMOVED***
		if spread, ok := argExpr.(*ast.SpreadElement); ok ***REMOVED***
			args[i] = c.compileSpreadCallArgument(spread)
			isVariadic = true
		***REMOVED*** else ***REMOVED***
			args[i] = c.compileExpression(argExpr)
		***REMOVED***
	***REMOVED***

	r := &compiledCallExpr***REMOVED***
		args:       args,
		callee:     c.compileExpression(v.Callee),
		isVariadic: isVariadic,
	***REMOVED***
	r.init(c, v.LeftParenthesis)
	return r
***REMOVED***

func (c *compiler) compileIdentifierExpression(v *ast.Identifier) compiledExpr ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierName(v.Name, int(v.Idx)-1)
	***REMOVED***

	r := &compiledIdentifierExpr***REMOVED***
		name: v.Name,
	***REMOVED***
	r.offset = int(v.Idx) - 1
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileNumberLiteral(v *ast.NumberLiteral) compiledExpr ***REMOVED***
	if c.scope.strict && octalRegexp.MatchString(v.Literal) ***REMOVED***
		c.throwSyntaxError(int(v.Idx)-1, "Octal literals are not allowed in strict mode")
		panic("Unreachable")
	***REMOVED***
	var val Value
	switch num := v.Value.(type) ***REMOVED***
	case int64:
		val = intToValue(num)
	case float64:
		val = floatToValue(num)
	default:
		panic(fmt.Errorf("Unsupported number literal type: %T", v.Value))
	***REMOVED***
	r := &compiledLiteral***REMOVED***
		val: val,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileStringLiteral(v *ast.StringLiteral) compiledExpr ***REMOVED***
	r := &compiledLiteral***REMOVED***
		val: stringValueFromRaw(v.Value),
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileBooleanLiteral(v *ast.BooleanLiteral) compiledExpr ***REMOVED***
	var val Value
	if v.Value ***REMOVED***
		val = valueTrue
	***REMOVED*** else ***REMOVED***
		val = valueFalse
	***REMOVED***

	r := &compiledLiteral***REMOVED***
		val: val,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileAssignExpression(v *ast.AssignExpression) compiledExpr ***REMOVED***
	// log.Printf("compileAssignExpression(): %+v", v)

	r := &compiledAssignExpr***REMOVED***
		left:     c.compileExpression(v.Left),
		right:    c.compileExpression(v.Right),
		operator: v.Operator,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledEnumGetExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.c.emit(enumGet)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileObjectAssignmentPattern(v *ast.ObjectPattern) compiledExpr ***REMOVED***
	r := &compiledObjectAssignmentPattern***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledObjectAssignmentPattern) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.c.emit(loadUndef)
	***REMOVED***
***REMOVED***

func (c *compiler) compileArrayAssignmentPattern(v *ast.ArrayPattern) compiledExpr ***REMOVED***
	r := &compiledArrayAssignmentPattern***REMOVED***
		expr: v,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (e *compiledArrayAssignmentPattern) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.c.emit(loadUndef)
	***REMOVED***
***REMOVED***

func (c *compiler) emitNamed(expr compiledExpr, name unistring.String) ***REMOVED***
	if en, ok := expr.(interface ***REMOVED***
		emitNamed(name unistring.String)
	***REMOVED***); ok ***REMOVED***
		en.emitNamed(name)
	***REMOVED*** else ***REMOVED***
		expr.emitGetter(true)
	***REMOVED***
***REMOVED***

func (e *compiledFunctionLiteral) emitNamed(name unistring.String) ***REMOVED***
	e.lhsName = name
	e.emitGetter(true)
***REMOVED***

func (c *compiler) emitPattern(pattern ast.Pattern, emitter func(target, init compiledExpr), putOnStack bool) ***REMOVED***
	switch pattern := pattern.(type) ***REMOVED***
	case *ast.ObjectPattern:
		c.emitObjectPattern(pattern, emitter, putOnStack)
	case *ast.ArrayPattern:
		c.emitArrayPattern(pattern, emitter, putOnStack)
	default:
		panic(fmt.Errorf("unsupported Pattern: %T", pattern))
	***REMOVED***
***REMOVED***

func (c *compiler) emitAssign(target ast.Expression, init compiledExpr, emitAssignSimple func(target, init compiledExpr)) ***REMOVED***
	pattern, isPattern := target.(ast.Pattern)
	if isPattern ***REMOVED***
		init.emitGetter(true)
		c.emitPattern(pattern, emitAssignSimple, false)
	***REMOVED*** else ***REMOVED***
		emitAssignSimple(c.compileExpression(target), init)
	***REMOVED***
***REMOVED***

func (c *compiler) emitObjectPattern(pattern *ast.ObjectPattern, emitAssign func(target, init compiledExpr), putOnStack bool) ***REMOVED***
	if pattern.Rest != nil ***REMOVED***
		c.emit(createDestructSrc)
	***REMOVED*** else ***REMOVED***
		c.emit(checkObjectCoercible)
	***REMOVED***
	for _, prop := range pattern.Properties ***REMOVED***
		switch prop := prop.(type) ***REMOVED***
		case *ast.PropertyShort:
			c.emit(dup)
			emitAssign(c.compileIdentifierExpression(&prop.Name), c.compilePatternInitExpr(func() ***REMOVED***
				c.emit(getProp(prop.Name.Name))
			***REMOVED***, prop.Initializer, prop.Idx0()))
		case *ast.PropertyKeyed:
			c.emit(dup)
			c.compileExpression(prop.Key).emitGetter(true)
			c.emit(_toPropertyKey***REMOVED******REMOVED***)
			var target ast.Expression
			var initializer ast.Expression
			if e, ok := prop.Value.(*ast.AssignExpression); ok ***REMOVED***
				target = e.Left
				initializer = e.Right
			***REMOVED*** else ***REMOVED***
				target = prop.Value
			***REMOVED***
			c.emitAssign(target, c.compilePatternInitExpr(func() ***REMOVED***
				c.emit(getKey)
			***REMOVED***, initializer, prop.Idx0()), emitAssign)
		default:
			c.throwSyntaxError(int(prop.Idx0()-1), "Unsupported AssignmentProperty type: %T", prop)
		***REMOVED***
	***REMOVED***
	if pattern.Rest != nil ***REMOVED***
		emitAssign(c.compileExpression(pattern.Rest), c.compileEmitterExpr(func() ***REMOVED***
			c.emit(copyRest)
		***REMOVED***, pattern.Rest.Idx0()))
		c.emit(pop)
	***REMOVED***
	if !putOnStack ***REMOVED***
		c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) emitArrayPattern(pattern *ast.ArrayPattern, emitAssign func(target, init compiledExpr), putOnStack bool) ***REMOVED***
	var marks []int
	c.emit(iterate)
	for _, elt := range pattern.Elements ***REMOVED***
		switch elt := elt.(type) ***REMOVED***
		case nil:
			marks = append(marks, len(c.p.code))
			c.emit(nil)
		case *ast.AssignExpression:
			c.emitAssign(elt.Left, c.compilePatternInitExpr(func() ***REMOVED***
				marks = append(marks, len(c.p.code))
				c.emit(nil, enumGet)
			***REMOVED***, elt.Right, elt.Idx0()), emitAssign)
		default:
			c.emitAssign(elt, c.compileEmitterExpr(func() ***REMOVED***
				marks = append(marks, len(c.p.code))
				c.emit(nil, enumGet)
			***REMOVED***, elt.Idx0()), emitAssign)
		***REMOVED***
	***REMOVED***
	if pattern.Rest != nil ***REMOVED***
		c.emitAssign(pattern.Rest, c.compileEmitterExpr(func() ***REMOVED***
			c.emit(newArrayFromIter)
		***REMOVED***, pattern.Rest.Idx0()), emitAssign)
	***REMOVED*** else ***REMOVED***
		c.emit(enumPopClose)
	***REMOVED***
	mark1 := len(c.p.code)
	c.emit(nil)

	for i, elt := range pattern.Elements ***REMOVED***
		switch elt := elt.(type) ***REMOVED***
		case nil:
			c.p.code[marks[i]] = iterNext(len(c.p.code) - marks[i])
		case *ast.Identifier:
			emitAssign(c.compileIdentifierExpression(elt), c.compileEmitterExpr(func() ***REMOVED***
				c.p.code[marks[i]] = iterNext(len(c.p.code) - marks[i])
				c.emit(loadUndef)
			***REMOVED***, elt.Idx0()))
		case *ast.AssignExpression:
			c.emitAssign(elt.Left, c.compileNamedEmitterExpr(func(name unistring.String) ***REMOVED***
				c.p.code[marks[i]] = iterNext(len(c.p.code) - marks[i])
				c.emitNamed(c.compileExpression(elt.Right), name)
			***REMOVED***, elt.Idx0()), emitAssign)
		default:
			c.emitAssign(elt, c.compileEmitterExpr(
				func() ***REMOVED***
					c.p.code[marks[i]] = iterNext(len(c.p.code) - marks[i])
					c.emit(loadUndef)
				***REMOVED***, elt.Idx0()), emitAssign)
		***REMOVED***
	***REMOVED***
	c.emit(enumPop)
	if pattern.Rest != nil ***REMOVED***
		c.emitAssign(pattern.Rest, c.compileExpression(
			&ast.ArrayLiteral***REMOVED***
				LeftBracket:  pattern.Rest.Idx0(),
				RightBracket: pattern.Rest.Idx0(),
			***REMOVED***), emitAssign)
	***REMOVED***
	c.p.code[mark1] = jump(len(c.p.code) - mark1)

	if !putOnStack ***REMOVED***
		c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledObjectAssignmentPattern) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	valueExpr.emitGetter(true)
	e.c.emitObjectPattern(e.expr, e.c.emitPatternAssign, putOnStack)
***REMOVED***

func (e *compiledArrayAssignmentPattern) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	valueExpr.emitGetter(true)
	e.c.emitArrayPattern(e.expr, e.c.emitPatternAssign, putOnStack)
***REMOVED***

type compiledPatternInitExpr struct ***REMOVED***
	baseCompiledExpr
	emitSrc func()
	def     compiledExpr
***REMOVED***

func (e *compiledPatternInitExpr) emitGetter(putOnStack bool) ***REMOVED***
	if !putOnStack ***REMOVED***
		return
	***REMOVED***
	e.emitSrc()
	if e.def != nil ***REMOVED***
		mark := len(e.c.p.code)
		e.c.emit(nil)
		e.def.emitGetter(true)
		e.c.p.code[mark] = jdef(len(e.c.p.code) - mark)
	***REMOVED***
***REMOVED***

func (e *compiledPatternInitExpr) emitNamed(name unistring.String) ***REMOVED***
	e.emitSrc()
	if e.def != nil ***REMOVED***
		mark := len(e.c.p.code)
		e.c.emit(nil)
		e.c.emitNamed(e.def, name)
		e.c.p.code[mark] = jdef(len(e.c.p.code) - mark)
	***REMOVED***
***REMOVED***

func (c *compiler) compilePatternInitExpr(emitSrc func(), def ast.Expression, idx file.Idx) compiledExpr ***REMOVED***
	r := &compiledPatternInitExpr***REMOVED***
		emitSrc: emitSrc,
		def:     c.compileExpression(def),
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

type compiledEmitterExpr struct ***REMOVED***
	baseCompiledExpr
	emitter      func()
	namedEmitter func(name unistring.String)
***REMOVED***

func (e *compiledEmitterExpr) emitGetter(putOnStack bool) ***REMOVED***
	if e.emitter != nil ***REMOVED***
		e.emitter()
	***REMOVED*** else ***REMOVED***
		e.namedEmitter("")
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledEmitterExpr) emitNamed(name unistring.String) ***REMOVED***
	if e.namedEmitter != nil ***REMOVED***
		e.namedEmitter(name)
	***REMOVED*** else ***REMOVED***
		e.emitter()
	***REMOVED***
***REMOVED***

func (c *compiler) compileEmitterExpr(emitter func(), idx file.Idx) *compiledEmitterExpr ***REMOVED***
	r := &compiledEmitterExpr***REMOVED***
		emitter: emitter,
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (c *compiler) compileNamedEmitterExpr(namedEmitter func(unistring.String), idx file.Idx) *compiledEmitterExpr ***REMOVED***
	r := &compiledEmitterExpr***REMOVED***
		namedEmitter: namedEmitter,
	***REMOVED***
	r.init(c, idx)
	return r
***REMOVED***

func (e *compiledSpreadCallArgument) emitGetter(putOnStack bool) ***REMOVED***
	e.expr.emitGetter(putOnStack)
	if putOnStack ***REMOVED***
		e.c.emit(pushSpread)
	***REMOVED***
***REMOVED***
