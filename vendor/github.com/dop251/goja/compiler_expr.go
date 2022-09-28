package goja

import (
	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
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

type compiledTemplateLiteral struct ***REMOVED***
	baseCompiledExpr
	tag         compiledExpr
	elements    []*ast.TemplateElement
	expressions []compiledExpr
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

type funcType uint8

const (
	funcNone funcType = iota
	funcRegular
	funcArrow
	funcMethod
	funcClsInit
	funcCtor
	funcDerivedCtor
)

type compiledFunctionLiteral struct ***REMOVED***
	baseCompiledExpr
	name            *ast.Identifier
	parameterList   *ast.ParameterList
	body            []ast.Statement
	source          string
	declarationList []*ast.VariableDeclaration
	lhsName         unistring.String
	strict          *ast.StringLiteral
	homeObjOffset   uint32
	typ             funcType
	isExpr          bool
***REMOVED***

type compiledBracketExpr struct ***REMOVED***
	baseCompiledExpr
	left, member compiledExpr
***REMOVED***

type compiledThisExpr struct ***REMOVED***
	baseCompiledExpr
***REMOVED***

type compiledSuperExpr struct ***REMOVED***
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

type compiledCoalesce struct ***REMOVED***
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

type compiledOptionalChain struct ***REMOVED***
	baseCompiledExpr
	expr compiledExpr
***REMOVED***

type compiledOptional struct ***REMOVED***
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
	case *ast.TemplateLiteral:
		return c.compileTemplateLiteral(v)
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
	case *ast.ClassLiteral:
		return c.compileClassLiteral(v, true)
	case *ast.DotExpression:
		return c.compileDotExpression(v)
	case *ast.PrivateDotExpression:
		return c.compilePrivateDotExpression(v)
	case *ast.BracketExpression:
		return c.compileBracketExpression(v)
	case *ast.ThisExpression:
		r := &compiledThisExpr***REMOVED******REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.SuperExpression:
		c.throwSyntaxError(int(v.Idx0())-1, "'super' keyword unexpected here")
		panic("unreachable")
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
	case *ast.OptionalChain:
		r := &compiledOptionalChain***REMOVED***
			expr: c.compileExpression(v.Expression),
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	case *ast.Optional:
		r := &compiledOptional***REMOVED***
			expr: c.compileExpression(v.Expression),
		***REMOVED***
		r.init(c, v.Idx0())
		return r
	default:
		c.assert(false, int(v.Idx0())-1, "Unknown expression type: %T", v)
		panic("unreachable")
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
	e.c.assert(false, e.offset, "Cannot emit reference for this type of expression")
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
	if e.offset >= 0 ***REMOVED***
		e.c.p.addSrcMap(e.offset)
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
		e.c.assert(b != nil, e.offset, "No dynamics and not found")
		if putOnStack ***REMOVED***
			b.emitGet()
		***REMOVED*** else ***REMOVED***
			b.emitGetP()
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
		e.c.assert(b != nil, e.offset, "No dynamics and not found")
		b.emitGet()
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
		e.c.assert(b != nil, e.offset, "No dynamics and not found")
		e.c.emit(loadUndef)
		b.emitGet()
	***REMOVED*** else ***REMOVED***
		if b != nil ***REMOVED***
			b.emitGetVar(true)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadDynamicCallee(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitVarSetter1(putOnStack bool, emitRight func(isRef bool)) ***REMOVED***
	e.addSrcMap()
	c := e.c

	if b, noDynamics := c.scope.lookupName(e.name); noDynamics ***REMOVED***
		if c.scope.strict ***REMOVED***
			c.checkIdentifierLName(e.name, e.offset)
		***REMOVED***
		emitRight(false)
		if b != nil ***REMOVED***
			if putOnStack ***REMOVED***
				b.emitSet()
			***REMOVED*** else ***REMOVED***
				b.emitSetP()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if c.scope.strict ***REMOVED***
				c.emit(setGlobalStrict(e.name))
			***REMOVED*** else ***REMOVED***
				c.emit(setGlobal(e.name))
			***REMOVED***
			if !putOnStack ***REMOVED***
				c.emit(pop)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.emitVarRef(e.name, e.offset, b)
		emitRight(true)
		if putOnStack ***REMOVED***
			c.emit(putValue)
		***REMOVED*** else ***REMOVED***
			c.emit(putValueP)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledIdentifierExpr) emitVarSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.emitVarSetter1(putOnStack, func(bool) ***REMOVED***
		e.c.emitNamedOrConst(valueExpr, e.name)
	***REMOVED***)
***REMOVED***

func (c *compiler) emitVarRef(name unistring.String, offset int, b *binding) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierLName(name, offset)
	***REMOVED***

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
	b, _ := e.c.scope.lookupName(e.name)
	e.c.emitVarRef(e.name, e.offset, b)
***REMOVED***

func (e *compiledIdentifierExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.emitVarSetter(valueExpr, putOnStack)
***REMOVED***

func (e *compiledIdentifierExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.emitVarSetter1(true, func(isRef bool) ***REMOVED***
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
		e.emitVarSetter1(false, func(isRef bool) ***REMOVED***
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

type compiledSuperDotExpr struct ***REMOVED***
	baseCompiledExpr
	name unistring.String
***REMOVED***

func (e *compiledSuperDotExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.c.emitLoadThis()
	e.c.emit(loadSuper)
	e.addSrcMap()
	e.c.emit(getPropRecv(e.name))
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledSuperDotExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.c.emitLoadThis()
	e.c.emit(loadSuper)
	valueExpr.emitGetter(true)
	e.addSrcMap()
	if putOnStack ***REMOVED***
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setPropRecvStrict(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropRecv(e.name))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setPropRecvStrictP(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropRecvP(e.name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledSuperDotExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emitLoadThis()
		e.c.emit(loadSuper, dupLast(2), getPropRecv(e.name))
		body()
		e.addSrcMap()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setPropRecvStrictP(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropRecvP(e.name))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.c.emitLoadThis()
			e.c.emit(loadSuper, dupLast(2), getPropRecv(e.name))
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			e.addSrcMap()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setPropRecvStrict(e.name))
			***REMOVED*** else ***REMOVED***
				e.c.emit(setPropRecv(e.name))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.c.emitLoadThis()
			e.c.emit(loadSuper, dupLast(2), getPropRecv(e.name))
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(3))
			body()
			e.addSrcMap()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setPropRecvStrictP(e.name))
			***REMOVED*** else ***REMOVED***
				e.c.emit(setPropRecvP(e.name))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledSuperDotExpr) emitRef() ***REMOVED***
	e.c.emitLoadThis()
	e.c.emit(loadSuper)
	if e.c.scope.strict ***REMOVED***
		e.c.emit(getPropRefRecvStrict(e.name))
	***REMOVED*** else ***REMOVED***
		e.c.emit(getPropRefRecv(e.name))
	***REMOVED***
***REMOVED***

func (e *compiledSuperDotExpr) deleteExpr() compiledExpr ***REMOVED***
	return e.c.superDeleteError(e.offset)
***REMOVED***

type compiledDotExpr struct ***REMOVED***
	baseCompiledExpr
	left compiledExpr
	name unistring.String
***REMOVED***

type compiledPrivateDotExpr struct ***REMOVED***
	baseCompiledExpr
	left compiledExpr
	name unistring.String
***REMOVED***

func (c *compiler) checkSuperBase(idx file.Idx) ***REMOVED***
	if s := c.scope.nearestThis(); s != nil ***REMOVED***
		switch s.funcType ***REMOVED***
		case funcMethod, funcClsInit, funcCtor, funcDerivedCtor:
			return
		***REMOVED***
	***REMOVED***
	c.throwSyntaxError(int(idx)-1, "'super' keyword unexpected here")
	panic("unreachable")
***REMOVED***

func (c *compiler) compileDotExpression(v *ast.DotExpression) compiledExpr ***REMOVED***
	if sup, ok := v.Left.(*ast.SuperExpression); ok ***REMOVED***
		c.checkSuperBase(sup.Idx)
		r := &compiledSuperDotExpr***REMOVED***
			name: v.Identifier.Name,
		***REMOVED***
		r.init(c, v.Identifier.Idx)
		return r
	***REMOVED***

	r := &compiledDotExpr***REMOVED***
		left: c.compileExpression(v.Left),
		name: v.Identifier.Name,
	***REMOVED***
	r.init(c, v.Identifier.Idx)
	return r
***REMOVED***

func (c *compiler) compilePrivateDotExpression(v *ast.PrivateDotExpression) compiledExpr ***REMOVED***
	r := &compiledPrivateDotExpr***REMOVED***
		left: c.compileExpression(v.Left),
		name: v.Identifier.Name,
	***REMOVED***
	r.init(c, v.Identifier.Idx)
	return r
***REMOVED***

func (e *compiledPrivateDotExpr) _emitGetter(rn *resolvedPrivateName, id *privateId) ***REMOVED***
	if rn != nil ***REMOVED***
		e.c.emit((*getPrivatePropRes)(rn))
	***REMOVED*** else ***REMOVED***
		e.c.emit((*getPrivatePropId)(id))
	***REMOVED***
***REMOVED***

func (e *compiledPrivateDotExpr) _emitSetter(rn *resolvedPrivateName, id *privateId) ***REMOVED***
	if rn != nil ***REMOVED***
		e.c.emit((*setPrivatePropRes)(rn))
	***REMOVED*** else ***REMOVED***
		e.c.emit((*setPrivatePropId)(id))
	***REMOVED***
***REMOVED***

func (e *compiledPrivateDotExpr) _emitSetterP(rn *resolvedPrivateName, id *privateId) ***REMOVED***
	if rn != nil ***REMOVED***
		e.c.emit((*setPrivatePropResP)(rn))
	***REMOVED*** else ***REMOVED***
		e.c.emit((*setPrivatePropIdP)(id))
	***REMOVED***
***REMOVED***

func (e *compiledPrivateDotExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.left.emitGetter(true)
	e.addSrcMap()
	rn, id := e.c.resolvePrivateName(e.name, e.offset)
	e._emitGetter(rn, id)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledPrivateDotExpr) emitSetter(v compiledExpr, putOnStack bool) ***REMOVED***
	rn, id := e.c.resolvePrivateName(e.name, e.offset)
	e.left.emitGetter(true)
	v.emitGetter(true)
	e.addSrcMap()
	if putOnStack ***REMOVED***
		e._emitSetter(rn, id)
	***REMOVED*** else ***REMOVED***
		e._emitSetterP(rn, id)
	***REMOVED***
***REMOVED***

func (e *compiledPrivateDotExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	rn, id := e.c.resolvePrivateName(e.name, e.offset)
	if !putOnStack ***REMOVED***
		e.left.emitGetter(true)
		e.c.emit(dup)
		e._emitGetter(rn, id)
		body()
		e.addSrcMap()
		e._emitSetterP(rn, id)
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.left.emitGetter(true)
			e.c.emit(dup)
			e._emitGetter(rn, id)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			e.addSrcMap()
			e._emitSetter(rn, id)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.left.emitGetter(true)
			e.c.emit(dup)
			e._emitGetter(rn, id)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(2))
			body()
			e.addSrcMap()
			e._emitSetterP(rn, id)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledPrivateDotExpr) deleteExpr() compiledExpr ***REMOVED***
	e.c.throwSyntaxError(e.offset, "Private fields can not be deleted")
	panic("unreachable")
***REMOVED***

func (e *compiledPrivateDotExpr) emitRef() ***REMOVED***
	e.left.emitGetter(true)
	rn, id := e.c.resolvePrivateName(e.name, e.offset)
	if rn != nil ***REMOVED***
		e.c.emit((*getPrivateRefRes)(rn))
	***REMOVED*** else ***REMOVED***
		e.c.emit((*getPrivateRefId)(id))
	***REMOVED***
***REMOVED***

type compiledSuperBracketExpr struct ***REMOVED***
	baseCompiledExpr
	member compiledExpr
***REMOVED***

func (e *compiledSuperBracketExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.c.emitLoadThis()
	e.member.emitGetter(true)
	e.c.emit(loadSuper)
	e.addSrcMap()
	e.c.emit(getElemRecv)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledSuperBracketExpr) emitSetter(valueExpr compiledExpr, putOnStack bool) ***REMOVED***
	e.c.emitLoadThis()
	e.member.emitGetter(true)
	e.c.emit(loadSuper)
	valueExpr.emitGetter(true)
	e.addSrcMap()
	if putOnStack ***REMOVED***
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setElemRecvStrict)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElemRecv)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setElemRecvStrictP)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElemRecvP)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledSuperBracketExpr) emitUnary(prepare, body func(), postfix, putOnStack bool) ***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emitLoadThis()
		e.member.emitGetter(true)
		e.c.emit(loadSuper, dupLast(3), getElemRecv)
		body()
		e.addSrcMap()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setElemRecvStrictP)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElemRecvP)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.c.emitLoadThis()
			e.member.emitGetter(true)
			e.c.emit(loadSuper, dupLast(3), getElemRecv)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			e.addSrcMap()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setElemRecvStrict)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setElemRecv)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.c.emitLoadThis()
			e.member.emitGetter(true)
			e.c.emit(loadSuper, dupLast(3), getElemRecv)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(4))
			body()
			e.addSrcMap()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setElemRecvStrictP)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setElemRecvP)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledSuperBracketExpr) emitRef() ***REMOVED***
	e.c.emitLoadThis()
	e.member.emitGetter(true)
	e.c.emit(loadSuper)
	if e.c.scope.strict ***REMOVED***
		e.c.emit(getElemRefRecvStrict)
	***REMOVED*** else ***REMOVED***
		e.c.emit(getElemRefRecv)
	***REMOVED***
***REMOVED***

func (c *compiler) superDeleteError(offset int) compiledExpr ***REMOVED***
	return c.compileEmitterExpr(func() ***REMOVED***
		c.emit(throwConst***REMOVED***referenceError("Unsupported reference to 'super'")***REMOVED***)
	***REMOVED***, file.Idx(offset+1))
***REMOVED***

func (e *compiledSuperBracketExpr) deleteExpr() compiledExpr ***REMOVED***
	return e.c.superDeleteError(e.offset)
***REMOVED***

func (c *compiler) checkConstantString(expr compiledExpr) (unistring.String, bool) ***REMOVED***
	if expr.constant() ***REMOVED***
		if val, ex := c.evalConst(expr); ex == nil ***REMOVED***
			if s, ok := val.(valueString); ok ***REMOVED***
				return s.string(), true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return "", false
***REMOVED***

func (c *compiler) compileBracketExpression(v *ast.BracketExpression) compiledExpr ***REMOVED***
	if sup, ok := v.Left.(*ast.SuperExpression); ok ***REMOVED***
		c.checkSuperBase(sup.Idx)
		member := c.compileExpression(v.Member)
		if name, ok := c.checkConstantString(member); ok ***REMOVED***
			r := &compiledSuperDotExpr***REMOVED***
				name: name,
			***REMOVED***
			r.init(c, v.LeftBracket)
			return r
		***REMOVED***

		r := &compiledSuperBracketExpr***REMOVED***
			member: member,
		***REMOVED***
		r.init(c, v.LeftBracket)
		return r
	***REMOVED***

	left := c.compileExpression(v.Left)
	member := c.compileExpression(v.Member)
	if name, ok := c.checkConstantString(member); ok ***REMOVED***
		r := &compiledDotExpr***REMOVED***
			left: left,
			name: name,
		***REMOVED***
		r.init(c, v.LeftBracket)
		return r
	***REMOVED***

	r := &compiledBracketExpr***REMOVED***
		left:   left,
		member: member,
	***REMOVED***
	r.init(c, v.LeftBracket)
	return r
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
		e.addSrcMap()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setPropStrictP(e.name))
		***REMOVED*** else ***REMOVED***
			e.c.emit(setPropP(e.name))
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
			e.addSrcMap()
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
			e.addSrcMap()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setPropStrictP(e.name))
			***REMOVED*** else ***REMOVED***
				e.c.emit(setPropP(e.name))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledDotExpr) deleteExpr() compiledExpr ***REMOVED***
	r := &deletePropExpr***REMOVED***
		left: e.left,
		name: e.name,
	***REMOVED***
	r.init(e.c, file.Idx(e.offset)+1)
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
	e.addSrcMap()
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
		e.c.emit(dupLast(2), getElem)
		body()
		e.addSrcMap()
		if e.c.scope.strict ***REMOVED***
			e.c.emit(setElemStrict, pop)
		***REMOVED*** else ***REMOVED***
			e.c.emit(setElem, pop)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !postfix ***REMOVED***
			e.left.emitGetter(true)
			e.member.emitGetter(true)
			e.c.emit(dupLast(2), getElem)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			body()
			e.addSrcMap()
			if e.c.scope.strict ***REMOVED***
				e.c.emit(setElemStrict)
			***REMOVED*** else ***REMOVED***
				e.c.emit(setElem)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef)
			e.left.emitGetter(true)
			e.member.emitGetter(true)
			e.c.emit(dupLast(2), getElem)
			if prepare != nil ***REMOVED***
				prepare()
			***REMOVED***
			e.c.emit(rdupN(3))
			body()
			e.addSrcMap()
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
	r.init(e.c, file.Idx(e.offset)+1)
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
	switch e.operator ***REMOVED***
	case token.ASSIGN:
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
	case token.EXPONENT:
		e.left.emitUnary(nil, func() ***REMOVED***
			e.right.emitGetter(true)
			e.c.emit(exp)
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
		e.c.assert(false, e.offset, "Unknown assign operator: %s", e.operator.String())
		panic("unreachable")
	***REMOVED***
***REMOVED***

func (e *compiledLiteral) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.c.emit(loadVal(e.c.p.defineLiteralValue(e.val)))
	***REMOVED***
***REMOVED***

func (e *compiledLiteral) constant() bool ***REMOVED***
	return true
***REMOVED***

func (e *compiledTemplateLiteral) emitGetter(putOnStack bool) ***REMOVED***
	if e.tag == nil ***REMOVED***
		if len(e.elements) == 0 ***REMOVED***
			e.c.emit(loadVal(e.c.p.defineLiteralValue(stringEmpty)))
		***REMOVED*** else ***REMOVED***
			tail := e.elements[len(e.elements)-1].Parsed
			if len(e.elements) == 1 ***REMOVED***
				e.c.emit(loadVal(e.c.p.defineLiteralValue(stringValueFromRaw(tail))))
			***REMOVED*** else ***REMOVED***
				stringCount := 0
				if head := e.elements[0].Parsed; head != "" ***REMOVED***
					e.c.emit(loadVal(e.c.p.defineLiteralValue(stringValueFromRaw(head))))
					stringCount++
				***REMOVED***
				e.expressions[0].emitGetter(true)
				e.c.emit(_toString***REMOVED******REMOVED***)
				stringCount++
				for i := 1; i < len(e.elements)-1; i++ ***REMOVED***
					if elt := e.elements[i].Parsed; elt != "" ***REMOVED***
						e.c.emit(loadVal(e.c.p.defineLiteralValue(stringValueFromRaw(elt))))
						stringCount++
					***REMOVED***
					e.expressions[i].emitGetter(true)
					e.c.emit(_toString***REMOVED******REMOVED***)
					stringCount++
				***REMOVED***
				if tail != "" ***REMOVED***
					e.c.emit(loadVal(e.c.p.defineLiteralValue(stringValueFromRaw(tail))))
					stringCount++
				***REMOVED***
				e.c.emit(concatStrings(stringCount))
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		cooked := make([]Value, len(e.elements))
		raw := make([]Value, len(e.elements))
		for i, elt := range e.elements ***REMOVED***
			raw[i] = &valueProperty***REMOVED***
				enumerable: true,
				value:      newStringValue(elt.Literal),
			***REMOVED***
			var cookedVal Value
			if elt.Valid ***REMOVED***
				cookedVal = stringValueFromRaw(elt.Parsed)
			***REMOVED*** else ***REMOVED***
				cookedVal = _undefined
			***REMOVED***
			cooked[i] = &valueProperty***REMOVED***
				enumerable: true,
				value:      cookedVal,
			***REMOVED***
		***REMOVED***
		e.c.emitCallee(e.tag)
		e.c.emit(&getTaggedTmplObject***REMOVED***
			raw:    raw,
			cooked: cooked,
		***REMOVED***)
		for _, expr := range e.expressions ***REMOVED***
			expr.emitGetter(true)
		***REMOVED***
		e.c.emit(call(len(e.expressions) + 1))
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
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

func (c *compiler) newCode(length, minCap int) (buf []instruction) ***REMOVED***
	if c.codeScratchpad != nil ***REMOVED***
		buf = c.codeScratchpad
		c.codeScratchpad = nil
	***REMOVED***
	if cap(buf) < minCap ***REMOVED***
		buf = make([]instruction, length, minCap)
	***REMOVED*** else ***REMOVED***
		buf = buf[:length]
	***REMOVED***
	return
***REMOVED***

func (e *compiledFunctionLiteral) compile() (prg *Program, name unistring.String, length int, strict bool) ***REMOVED***
	e.c.assert(e.typ != funcNone, e.offset, "compiledFunctionLiteral.typ is not set")

	savedPrg := e.c.p
	preambleLen := 8 // enter, boxThis, loadStack(0), initThis, createArgs, set, loadCallee, init
	e.c.p = &Program***REMOVED***
		src:  e.c.p.src,
		code: e.c.newCode(preambleLen, 16),
	***REMOVED***
	e.c.newScope()
	s := e.c.scope
	s.funcType = e.typ

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

		if firstDupIdx >= 0 && (hasPatterns || hasInits || s.strict || e.typ == funcArrow || e.typ == funcMethod) ***REMOVED***
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

	var thisBinding *binding
	if e.typ != funcArrow ***REMOVED***
		thisBinding = s.createThisBinding()
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
					e.c.emitPatternLexicalAssign(target, init)
				***REMOVED***, false)
			***REMOVED*** else if item.Initializer != nil ***REMOVED***
				markGet := len(e.c.p.code)
				e.c.emit(nil)
				mark := len(e.c.p.code)
				e.c.emit(nil)
				e.c.emitExpr(e.c.compileExpression(item.Initializer), true)
				if firstForwardRef == -1 && (s.isDynamic() || s.bindings[i].useCount() > 0) ***REMOVED***
					firstForwardRef = i
				***REMOVED***
				if firstForwardRef == -1 ***REMOVED***
					s.bindings[i].emitGetAt(markGet)
				***REMOVED*** else ***REMOVED***
					e.c.p.code[markGet] = loadStackLex(-i - 1)
				***REMOVED***
				s.bindings[i].emitInitP()
				e.c.p.code[mark] = jdefP(len(e.c.p.code) - mark)
			***REMOVED*** else ***REMOVED***
				if firstForwardRef == -1 && s.bindings[i].useCount() > 0 ***REMOVED***
					firstForwardRef = i
				***REMOVED***
				if firstForwardRef != -1 ***REMOVED***
					e.c.emit(loadStackLex(-i - 1))
					s.bindings[i].emitInitP()
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
					e.c.emitPatternLexicalAssign(target, init)
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
			calleeBinding.emitInitP()
		***REMOVED***
	***REMOVED***

	e.c.compileFunctions(funcs)
	e.c.compileStatements(body, false)

	var last ast.Statement
	if l := len(body); l > 0 ***REMOVED***
		last = body[l-1]
	***REMOVED***
	if _, ok := last.(*ast.ReturnStatement); !ok ***REMOVED***
		if e.typ == funcDerivedCtor ***REMOVED***
			e.c.emit(loadUndef)
			thisBinding.markAccessPoint()
			e.c.emit(ret)
		***REMOVED*** else ***REMOVED***
			e.c.emit(loadUndef, ret)
		***REMOVED***
	***REMOVED***

	delta := 0
	code := e.c.p.code

	if s.isDynamic() && !s.argsInStash ***REMOVED***
		s.moveArgsToStash()
	***REMOVED***

	if s.argsNeeded || s.isDynamic() && e.typ != funcArrow && e.typ != funcClsInit ***REMOVED***
		if e.typ == funcClsInit ***REMOVED***
			e.c.throwSyntaxError(e.offset, "'arguments' is not allowed in class field initializer or static initialization block")
		***REMOVED***
		b, created := s.bindNameLexical("arguments", false, 0)
		if created || b.isVar ***REMOVED***
			if !s.argsInStash ***REMOVED***
				s.moveArgsToStash()
			***REMOVED***
			if s.strict ***REMOVED***
				b.isConst = true
			***REMOVED*** else ***REMOVED***
				b.isVar = e.c.scope.isFunction()
			***REMOVED***
			pos := preambleLen - 2
			delta += 2
			if s.strict || hasPatterns || hasInits ***REMOVED***
				code[pos] = createArgsUnmapped(paramsCount)
			***REMOVED*** else ***REMOVED***
				code[pos] = createArgsMapped(paramsCount)
			***REMOVED***
			pos++
			b.emitInitPAtScope(s, pos)
		***REMOVED***
	***REMOVED***

	if calleeBinding != nil ***REMOVED***
		if !s.isDynamic() && calleeBinding.useCount() == 0 ***REMOVED***
			s.deleteBinding(calleeBinding)
			calleeBinding = nil
		***REMOVED*** else ***REMOVED***
			delta++
			calleeBinding.emitInitPAtScope(s, preambleLen-delta)
			delta++
			code[preambleLen-delta] = loadCallee
		***REMOVED***
	***REMOVED***

	if thisBinding != nil ***REMOVED***
		if !s.isDynamic() && thisBinding.useCount() == 0 ***REMOVED***
			s.deleteBinding(thisBinding)
			thisBinding = nil
		***REMOVED*** else ***REMOVED***
			if thisBinding.inStash || s.isDynamic() ***REMOVED***
				delta++
				thisBinding.emitInitAtScope(s, preambleLen-delta)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	stashSize, stackSize := s.finaliseVarAlloc(0)

	if thisBinding != nil && thisBinding.inStash && (!s.argsInStash || stackSize > 0) ***REMOVED***
		delta++
		code[preambleLen-delta] = loadStack(0)
	***REMOVED*** // otherwise, 'this' will be at stack[sp-1], no need to load

	if !s.strict && thisBinding != nil ***REMOVED***
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
				funcType:    e.typ,
			***REMOVED***
			if s.isDynamic() ***REMOVED***
				enter1.names = s.makeNamesMap()
			***REMOVED***
			enter = &enter1
			if enterFunc2Mark != -1 ***REMOVED***
				ef2 := &enterFuncBody***REMOVED***
					extensible: e.c.scope.dynamic,
					funcType:   e.typ,
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
				funcType:   e.typ,
			***REMOVED***
			if s.isDynamic() ***REMOVED***
				enter1.names = s.makeNamesMap()
			***REMOVED***
			enter = &enter1
			if enterFunc2Mark != -1 ***REMOVED***
				ef2 := &enterFuncBody***REMOVED***
					adjustStack: true,
					extensible:  e.c.scope.dynamic,
					funcType:    e.typ,
				***REMOVED***
				e.c.updateEnterBlock(&ef2.enterBlock)
				e.c.p.code[enterFunc2Mark] = ef2
			***REMOVED***
		***REMOVED***
		if emitArgsRestMark != -1 && s.argsInStash ***REMOVED***
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
				funcType:   e.typ,
			***REMOVED***
			e.c.updateEnterBlock(&ef2.enterBlock)
			e.c.p.code[enterFunc2Mark] = ef2
		***REMOVED***
	***REMOVED***
	code[delta] = enter
	s.trimCode(delta)

	strict = s.strict
	prg = e.c.p
	// e.c.p.dumpCode()
	if enterFunc2Mark != -1 ***REMOVED***
		e.c.popScope()
	***REMOVED***
	e.c.popScope()
	e.c.p = savedPrg

	return
***REMOVED***

func (e *compiledFunctionLiteral) emitGetter(putOnStack bool) ***REMOVED***
	p, name, length, strict := e.compile()
	switch e.typ ***REMOVED***
	case funcArrow:
		e.c.emit(&newArrowFunc***REMOVED***newFunc: newFunc***REMOVED***prg: p, length: length, name: name, source: e.source, strict: strict***REMOVED******REMOVED***)
	case funcMethod, funcClsInit:
		e.c.emit(&newMethod***REMOVED***newFunc: newFunc***REMOVED***prg: p, length: length, name: name, source: e.source, strict: strict***REMOVED***, homeObjOffset: e.homeObjOffset***REMOVED***)
	case funcRegular:
		e.c.emit(&newFunc***REMOVED***prg: p, length: length, name: name, source: e.source, strict: strict***REMOVED***)
	default:
		e.c.throwSyntaxError(e.offset, "Unsupported func type: %v", e.typ)
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
		typ:             funcRegular,
		strict:          strictBody,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

type compiledClassLiteral struct ***REMOVED***
	baseCompiledExpr
	name       *ast.Identifier
	superClass compiledExpr
	body       []ast.ClassElement
	lhsName    unistring.String
	source     string
	isExpr     bool
***REMOVED***

func (c *compiler) processKey(expr ast.Expression) (val unistring.String, computed bool) ***REMOVED***
	keyExpr := c.compileExpression(expr)
	if keyExpr.constant() ***REMOVED***
		v, ex := c.evalConst(keyExpr)
		if ex == nil ***REMOVED***
			return v.string(), false
		***REMOVED***
	***REMOVED***
	keyExpr.emitGetter(true)
	computed = true
	return
***REMOVED***

func (e *compiledClassLiteral) processClassKey(expr ast.Expression) (privateName *privateName, key unistring.String, computed bool) ***REMOVED***
	if p, ok := expr.(*ast.PrivateIdentifier); ok ***REMOVED***
		privateName = e.c.classScope.getDeclaredPrivateId(p.Name)
		key = privateIdString(p.Name)
		return
	***REMOVED***
	key, computed = e.c.processKey(expr)
	return
***REMOVED***

type clsElement struct ***REMOVED***
	key         unistring.String
	privateName *privateName
	initializer compiledExpr
	body        *compiledFunctionLiteral
	computed    bool
***REMOVED***

func (e *compiledClassLiteral) emitGetter(putOnStack bool) ***REMOVED***
	e.c.newBlockScope()
	s := e.c.scope
	s.strict = true

	enter := &enterBlock***REMOVED******REMOVED***
	mark0 := len(e.c.p.code)
	e.c.emit(enter)
	e.c.block = &block***REMOVED***
		typ:   blockScope,
		outer: e.c.block,
	***REMOVED***
	var clsBinding *binding
	var clsName unistring.String
	if name := e.name; name != nil ***REMOVED***
		clsName = name.Name
		clsBinding = e.c.createLexicalIdBinding(clsName, true, int(name.Idx)-1)
	***REMOVED*** else ***REMOVED***
		clsName = e.lhsName
	***REMOVED***

	var ctorMethod *ast.MethodDefinition
	ctorMethodIdx := -1
	staticsCount := 0
	instanceFieldsCount := 0
	hasStaticPrivateMethods := false
	cs := &classScope***REMOVED***
		c:     e.c,
		outer: e.c.classScope,
	***REMOVED***

	for idx, elt := range e.body ***REMOVED***
		switch elt := elt.(type) ***REMOVED***
		case *ast.ClassStaticBlock:
			if len(elt.Block.List) > 0 ***REMOVED***
				staticsCount++
			***REMOVED***
		case *ast.FieldDefinition:
			if id, ok := elt.Key.(*ast.PrivateIdentifier); ok ***REMOVED***
				cs.declarePrivateId(id.Name, ast.PropertyKindValue, elt.Static, int(elt.Idx)-1)
			***REMOVED***
			if elt.Static ***REMOVED***
				staticsCount++
			***REMOVED*** else ***REMOVED***
				instanceFieldsCount++
			***REMOVED***
		case *ast.MethodDefinition:
			if !elt.Static ***REMOVED***
				if id, ok := elt.Key.(*ast.StringLiteral); ok ***REMOVED***
					if !elt.Computed && id.Value == "constructor" ***REMOVED***
						if ctorMethod != nil ***REMOVED***
							e.c.throwSyntaxError(int(id.Idx)-1, "A class may only have one constructor")
						***REMOVED***
						ctorMethod = elt
						ctorMethodIdx = idx
						continue
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if id, ok := elt.Key.(*ast.PrivateIdentifier); ok ***REMOVED***
				cs.declarePrivateId(id.Name, elt.Kind, elt.Static, int(elt.Idx)-1)
				if elt.Static ***REMOVED***
					hasStaticPrivateMethods = true
				***REMOVED***
			***REMOVED***
		default:
			e.c.assert(false, int(elt.Idx0())-1, "Unsupported static element: %T", elt)
		***REMOVED***
	***REMOVED***

	var staticInit *newStaticFieldInit
	if staticsCount > 0 || hasStaticPrivateMethods ***REMOVED***
		staticInit = &newStaticFieldInit***REMOVED******REMOVED***
		e.c.emit(staticInit)
	***REMOVED***

	var derived bool
	var newClassIns *newClass
	if superClass := e.superClass; superClass != nil ***REMOVED***
		derived = true
		superClass.emitGetter(true)
		ndc := &newDerivedClass***REMOVED***
			newClass: newClass***REMOVED***
				name:   clsName,
				source: e.source,
			***REMOVED***,
		***REMOVED***
		e.addSrcMap()
		e.c.emit(ndc)
		newClassIns = &ndc.newClass
	***REMOVED*** else ***REMOVED***
		newClassIns = &newClass***REMOVED***
			name:   clsName,
			source: e.source,
		***REMOVED***
		e.addSrcMap()
		e.c.emit(newClassIns)
	***REMOVED***

	e.c.classScope = cs

	if ctorMethod != nil ***REMOVED***
		newClassIns.ctor, newClassIns.length = e.c.compileCtor(ctorMethod.Body, derived)
	***REMOVED***

	curIsPrototype := false

	instanceFields := make([]clsElement, 0, instanceFieldsCount)
	staticElements := make([]clsElement, 0, staticsCount)

	// stack at this point:
	//
	// staticFieldInit (if staticsCount > 0 || hasStaticPrivateMethods)
	// prototype
	// class function
	// <- sp

	for idx, elt := range e.body ***REMOVED***
		if idx == ctorMethodIdx ***REMOVED***
			continue
		***REMOVED***
		switch elt := elt.(type) ***REMOVED***
		case *ast.ClassStaticBlock:
			if len(elt.Block.List) > 0 ***REMOVED***
				f := e.c.compileFunctionLiteral(&ast.FunctionLiteral***REMOVED***
					Function:        elt.Idx0(),
					ParameterList:   &ast.ParameterList***REMOVED******REMOVED***,
					Body:            elt.Block,
					Source:          elt.Source,
					DeclarationList: elt.DeclarationList,
				***REMOVED***, true)
				f.typ = funcClsInit
				//f.lhsName = "<static_initializer>"
				f.homeObjOffset = 1
				staticElements = append(staticElements, clsElement***REMOVED***
					body: f,
				***REMOVED***)
			***REMOVED***
		case *ast.FieldDefinition:
			privateName, key, computed := e.processClassKey(elt.Key)
			var el clsElement
			if elt.Initializer != nil ***REMOVED***
				el.initializer = e.c.compileExpression(elt.Initializer)
			***REMOVED***
			el.computed = computed
			if computed ***REMOVED***
				if elt.Static ***REMOVED***
					if curIsPrototype ***REMOVED***
						e.c.emit(defineComputedKey(5))
					***REMOVED*** else ***REMOVED***
						e.c.emit(defineComputedKey(4))
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if curIsPrototype ***REMOVED***
						e.c.emit(defineComputedKey(3))
					***REMOVED*** else ***REMOVED***
						e.c.emit(defineComputedKey(2))
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				el.privateName = privateName
				el.key = key
			***REMOVED***
			if elt.Static ***REMOVED***
				staticElements = append(staticElements, el)
			***REMOVED*** else ***REMOVED***
				instanceFields = append(instanceFields, el)
			***REMOVED***
		case *ast.MethodDefinition:
			if elt.Static ***REMOVED***
				if curIsPrototype ***REMOVED***
					e.c.emit(pop)
					curIsPrototype = false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if !curIsPrototype ***REMOVED***
					e.c.emit(dupN(1))
					curIsPrototype = true
				***REMOVED***
			***REMOVED***
			privateName, key, computed := e.processClassKey(elt.Key)
			lit := e.c.compileFunctionLiteral(elt.Body, true)
			lit.typ = funcMethod
			if computed ***REMOVED***
				e.c.emit(_toPropertyKey***REMOVED******REMOVED***)
				lit.homeObjOffset = 2
			***REMOVED*** else ***REMOVED***
				lit.homeObjOffset = 1
				lit.lhsName = key
			***REMOVED***
			lit.emitGetter(true)
			if privateName != nil ***REMOVED***
				var offset int
				if elt.Static ***REMOVED***
					if curIsPrototype ***REMOVED***
						/*
							staticInit
							proto
							cls
							proto
							method
							<- sp
						*/
						offset = 5
					***REMOVED*** else ***REMOVED***
						/*
							staticInit
							proto
							cls
							method
							<- sp
						*/
						offset = 4
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if curIsPrototype ***REMOVED***
						offset = 3
					***REMOVED*** else ***REMOVED***
						offset = 2
					***REMOVED***
				***REMOVED***
				switch elt.Kind ***REMOVED***
				case ast.PropertyKindGet:
					e.c.emit(&definePrivateGetter***REMOVED***
						definePrivateMethod: definePrivateMethod***REMOVED***
							idx:          privateName.idx,
							targetOffset: offset,
						***REMOVED***,
					***REMOVED***)
				case ast.PropertyKindSet:
					e.c.emit(&definePrivateSetter***REMOVED***
						definePrivateMethod: definePrivateMethod***REMOVED***
							idx:          privateName.idx,
							targetOffset: offset,
						***REMOVED***,
					***REMOVED***)
				default:
					e.c.emit(&definePrivateMethod***REMOVED***
						idx:          privateName.idx,
						targetOffset: offset,
					***REMOVED***)
				***REMOVED***
			***REMOVED*** else if computed ***REMOVED***
				switch elt.Kind ***REMOVED***
				case ast.PropertyKindGet:
					e.c.emit(&defineGetter***REMOVED******REMOVED***)
				case ast.PropertyKindSet:
					e.c.emit(&defineSetter***REMOVED******REMOVED***)
				default:
					e.c.emit(&defineMethod***REMOVED******REMOVED***)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				switch elt.Kind ***REMOVED***
				case ast.PropertyKindGet:
					e.c.emit(&defineGetterKeyed***REMOVED***key: key***REMOVED***)
				case ast.PropertyKindSet:
					e.c.emit(&defineSetterKeyed***REMOVED***key: key***REMOVED***)
				default:
					e.c.emit(&defineMethodKeyed***REMOVED***key: key***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if curIsPrototype ***REMOVED***
		e.c.emit(pop)
	***REMOVED***

	if len(instanceFields) > 0 ***REMOVED***
		newClassIns.initFields = e.compileFieldsAndStaticBlocks(instanceFields, "<instance_members_initializer>")
	***REMOVED***
	if staticInit != nil ***REMOVED***
		if len(staticElements) > 0 ***REMOVED***
			staticInit.initFields = e.compileFieldsAndStaticBlocks(staticElements, "<static_initializer>")
		***REMOVED***
	***REMOVED***

	env := e.c.classScope.instanceEnv
	if s.dynLookup ***REMOVED***
		newClassIns.privateMethods, newClassIns.privateFields = env.methods, env.fields
	***REMOVED***
	newClassIns.numPrivateMethods = uint32(len(env.methods))
	newClassIns.numPrivateFields = uint32(len(env.fields))
	newClassIns.hasPrivateEnv = len(e.c.classScope.privateNames) > 0

	if (clsBinding != nil && clsBinding.useCount() > 0) || s.dynLookup ***REMOVED***
		if clsBinding != nil ***REMOVED***
			// Because this block may be in the middle of an expression, it's initial stack position
			// cannot be known, and therefore it may not have any stack variables.
			// Note, because clsBinding would be accessed through a function, it should already be in stash,
			// this is just to make sure.
			clsBinding.moveToStash()
			clsBinding.emitInit()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if clsBinding != nil ***REMOVED***
			s.deleteBinding(clsBinding)
			clsBinding = nil
		***REMOVED***
		e.c.p.code[mark0] = jump(1)
	***REMOVED***

	if staticsCount > 0 || hasStaticPrivateMethods ***REMOVED***
		ise := &initStaticElements***REMOVED******REMOVED***
		e.c.emit(ise)
		env := e.c.classScope.staticEnv
		staticInit.numPrivateFields = uint32(len(env.fields))
		staticInit.numPrivateMethods = uint32(len(env.methods))
		if s.dynLookup ***REMOVED***
			// These cannot be set on staticInit, because it is executed before ClassHeritage, and therefore
			// the VM's PrivateEnvironment is still not set.
			ise.privateFields = env.fields
			ise.privateMethods = env.methods
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		e.c.emit(endVariadic) // re-using as semantics match
	***REMOVED***

	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***

	if clsBinding != nil || s.dynLookup ***REMOVED***
		e.c.leaveScopeBlock(enter)
		e.c.assert(enter.stackSize == 0, e.offset, "enter.StackSize != 0 in compiledClassLiteral")
	***REMOVED*** else ***REMOVED***
		e.c.block = e.c.block.outer
	***REMOVED***
	if len(e.c.classScope.privateNames) > 0 ***REMOVED***
		e.c.emit(popPrivateEnv***REMOVED******REMOVED***)
	***REMOVED***
	e.c.classScope = e.c.classScope.outer
	e.c.popScope()
***REMOVED***

func (e *compiledClassLiteral) compileFieldsAndStaticBlocks(elements []clsElement, funcName unistring.String) *Program ***REMOVED***
	savedPrg := e.c.p
	savedBlock := e.c.block
	defer func() ***REMOVED***
		e.c.p = savedPrg
		e.c.block = savedBlock
	***REMOVED***()

	e.c.block = &block***REMOVED***
		typ: blockScope,
	***REMOVED***

	e.c.p = &Program***REMOVED***
		src:      savedPrg.src,
		funcName: funcName,
		code:     e.c.newCode(2, 16),
	***REMOVED***

	e.c.newScope()
	s := e.c.scope
	s.funcType = funcClsInit
	thisBinding := s.createThisBinding()

	valIdx := 0
	for _, elt := range elements ***REMOVED***
		if elt.body != nil ***REMOVED***
			e.c.emit(dup) // this
			elt.body.emitGetter(true)
			elt.body.addSrcMap()
			e.c.emit(call(0), pop)
		***REMOVED*** else ***REMOVED***
			if elt.computed ***REMOVED***
				e.c.emit(loadComputedKey(valIdx))
				valIdx++
			***REMOVED***
			if init := elt.initializer; init != nil ***REMOVED***
				if !elt.computed ***REMOVED***
					e.c.emitNamedOrConst(init, elt.key)
				***REMOVED*** else ***REMOVED***
					e.c.emitExpr(init, true)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				e.c.emit(loadUndef)
			***REMOVED***
			if elt.privateName != nil ***REMOVED***
				e.c.emit(&definePrivateProp***REMOVED***
					idx: elt.privateName.idx,
				***REMOVED***)
			***REMOVED*** else if elt.computed ***REMOVED***
				e.c.emit(defineProp***REMOVED******REMOVED***)
			***REMOVED*** else ***REMOVED***
				e.c.emit(definePropKeyed(elt.key))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	e.c.emit(halt)
	if s.isDynamic() || thisBinding.useCount() > 0 ***REMOVED***
		if s.isDynamic() || thisBinding.inStash ***REMOVED***
			thisBinding.emitInitAt(1)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		s.deleteBinding(thisBinding)
	***REMOVED***
	stashSize, stackSize := s.finaliseVarAlloc(0)
	e.c.assert(stackSize == 0, e.offset, "stackSize != 0 in initFields")
	if stashSize > 0 ***REMOVED***
		e.c.assert(stashSize == 1, e.offset, "stashSize != 1 in initFields")
		enter := &enterFunc***REMOVED***
			stashSize: 1,
			funcType:  funcClsInit,
		***REMOVED***
		if s.dynLookup ***REMOVED***
			enter.names = s.makeNamesMap()
		***REMOVED***
		e.c.p.code[0] = enter
		s.trimCode(0)
	***REMOVED*** else ***REMOVED***
		s.trimCode(2)
	***REMOVED***
	res := e.c.p
	e.c.popScope()
	return res
***REMOVED***

func (c *compiler) compileClassLiteral(v *ast.ClassLiteral, isExpr bool) *compiledClassLiteral ***REMOVED***
	if v.Name != nil ***REMOVED***
		c.checkIdentifierLName(v.Name.Name, int(v.Name.Idx)-1)
	***REMOVED***
	r := &compiledClassLiteral***REMOVED***
		name:       v.Name,
		superClass: c.compileExpression(v.SuperClass),
		body:       v.Body,
		source:     v.Source,
		isExpr:     isExpr,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) compileCtor(ctor *ast.FunctionLiteral, derived bool) (p *Program, length int) ***REMOVED***
	f := c.compileFunctionLiteral(ctor, true)
	if derived ***REMOVED***
		f.typ = funcDerivedCtor
	***REMOVED*** else ***REMOVED***
		f.typ = funcCtor
	***REMOVED***
	p, _, length, _ = f.compile()
	return
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
		typ:             funcArrow,
		strict:          strictBody,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

func (c *compiler) emitLoadThis() ***REMOVED***
	b, eval := c.scope.lookupThis()
	if b != nil ***REMOVED***
		b.emitGet()
	***REMOVED*** else ***REMOVED***
		if eval ***REMOVED***
			c.emit(getThisDynamic***REMOVED******REMOVED***)
		***REMOVED*** else ***REMOVED***
			c.emit(loadGlobalObject)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *compiledThisExpr) emitGetter(putOnStack bool) ***REMOVED***
	e.addSrcMap()
	e.c.emitLoadThis()
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledSuperExpr) emitGetter(putOnStack bool) ***REMOVED***
	if putOnStack ***REMOVED***
		e.c.emit(loadSuper)
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
	if s := e.c.scope.nearestThis(); s == nil || s.funcType == funcNone ***REMOVED***
		e.c.throwSyntaxError(e.offset, "new.target expression is not allowed here")
	***REMOVED***
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
	c.assert(false, 0, "unknown exception type thrown while evaluating constant expression: %s", v.String())
	panic("unreachable")
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
		e.addSrcMap()
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
		e.c.assert(false, e.offset, "Unknown unary operator: %s", e.operator.String())
		panic("unreachable")
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
	e.c.emitExpr(e.right, true)
	e.c.p.code[j] = jeq1(len(e.c.p.code) - j)
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledCoalesce) constant() bool ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if v != _null && v != _undefined ***REMOVED***
				return true
			***REMOVED***
			return e.right.constant()
		***REMOVED*** else ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (e *compiledCoalesce) emitGetter(putOnStack bool) ***REMOVED***
	if e.left.constant() ***REMOVED***
		if v, ex := e.c.evalConst(e.left); ex == nil ***REMOVED***
			if v == _undefined || v == _null ***REMOVED***
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
	e.c.emitExpr(e.right, true)
	e.c.p.code[j] = jcoalesc(len(e.c.p.code) - j)
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
	case token.EXPONENT:
		e.c.emit(exp)
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
		e.c.assert(false, e.offset, "Unknown operator: %s", e.operator.String())
		panic("unreachable")
	***REMOVED***

	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compileBinaryExpression(v *ast.BinaryExpression) compiledExpr ***REMOVED***

	switch v.Operator ***REMOVED***
	case token.LOGICAL_OR:
		return c.compileLogicalOr(v.Left, v.Right, v.Idx0())
	case token.COALESCE:
		return c.compileCoalesce(v.Left, v.Right, v.Idx0())
	case token.LOGICAL_AND:
		return c.compileLogicalAnd(v.Left, v.Right, v.Idx0())
	***REMOVED***

	if id, ok := v.Left.(*ast.PrivateIdentifier); ok ***REMOVED***
		return c.compilePrivateIn(id, v.Right, id.Idx)
	***REMOVED***

	r := &compiledBinaryExpr***REMOVED***
		left:     c.compileExpression(v.Left),
		right:    c.compileExpression(v.Right),
		operator: v.Operator,
	***REMOVED***
	r.init(c, v.Idx0())
	return r
***REMOVED***

type compiledPrivateIn struct ***REMOVED***
	baseCompiledExpr
	id    unistring.String
	right compiledExpr
***REMOVED***

func (e *compiledPrivateIn) emitGetter(putOnStack bool) ***REMOVED***
	e.right.emitGetter(true)
	rn, id := e.c.resolvePrivateName(e.id, e.offset)
	if rn != nil ***REMOVED***
		e.c.emit((*privateInRes)(rn))
	***REMOVED*** else ***REMOVED***
		e.c.emit((*privateInId)(id))
	***REMOVED***
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (c *compiler) compilePrivateIn(id *ast.PrivateIdentifier, right ast.Expression, idx file.Idx) compiledExpr ***REMOVED***
	r := &compiledPrivateIn***REMOVED***
		id:    id.Name,
		right: c.compileExpression(right),
	***REMOVED***
	r.init(c, idx)
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

func (c *compiler) compileCoalesce(left, right ast.Expression, idx file.Idx) compiledExpr ***REMOVED***
	r := &compiledCoalesce***REMOVED***
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
	hasProto := false
	for _, prop := range e.expr.Value ***REMOVED***
		switch prop := prop.(type) ***REMOVED***
		case *ast.PropertyKeyed:
			key, computed := e.c.processKey(prop.Key)
			valueExpr := e.c.compileExpression(prop.Value)
			var ne namedEmitter
			if fn, ok := valueExpr.(*compiledFunctionLiteral); ok ***REMOVED***
				if fn.name == nil ***REMOVED***
					ne = fn
				***REMOVED***
				switch prop.Kind ***REMOVED***
				case ast.PropertyKindMethod, ast.PropertyKindGet, ast.PropertyKindSet:
					fn.typ = funcMethod
					if computed ***REMOVED***
						fn.homeObjOffset = 2
					***REMOVED*** else ***REMOVED***
						fn.homeObjOffset = 1
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if v, ok := valueExpr.(namedEmitter); ok ***REMOVED***
				ne = v
			***REMOVED***
			if computed ***REMOVED***
				e.c.emit(_toPropertyKey***REMOVED******REMOVED***)
				e.c.emitExpr(valueExpr, true)
				switch prop.Kind ***REMOVED***
				case ast.PropertyKindValue:
					if ne != nil ***REMOVED***
						e.c.emit(setElem1Named)
					***REMOVED*** else ***REMOVED***
						e.c.emit(setElem1)
					***REMOVED***
				case ast.PropertyKindMethod:
					e.c.emit(&defineMethod***REMOVED***enumerable: true***REMOVED***)
				case ast.PropertyKindGet:
					e.c.emit(&defineGetter***REMOVED***enumerable: true***REMOVED***)
				case ast.PropertyKindSet:
					e.c.emit(&defineSetter***REMOVED***enumerable: true***REMOVED***)
				default:
					e.c.assert(false, e.offset, "unknown property kind: %s", prop.Kind)
					panic("unreachable")
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				isProto := key == __proto__ && !prop.Computed
				if isProto ***REMOVED***
					if hasProto ***REMOVED***
						e.c.throwSyntaxError(int(prop.Idx0())-1, "Duplicate __proto__ fields are not allowed in object literals")
					***REMOVED*** else ***REMOVED***
						hasProto = true
					***REMOVED***
				***REMOVED***
				if ne != nil && !isProto ***REMOVED***
					ne.emitNamed(key)
				***REMOVED*** else ***REMOVED***
					e.c.emitExpr(valueExpr, true)
				***REMOVED***
				switch prop.Kind ***REMOVED***
				case ast.PropertyKindValue:
					if isProto ***REMOVED***
						e.c.emit(setProto)
					***REMOVED*** else ***REMOVED***
						e.c.emit(putProp(key))
					***REMOVED***
				case ast.PropertyKindMethod:
					e.c.emit(&defineMethodKeyed***REMOVED***key: key, enumerable: true***REMOVED***)
				case ast.PropertyKindGet:
					e.c.emit(&defineGetterKeyed***REMOVED***key: key, enumerable: true***REMOVED***)
				case ast.PropertyKindSet:
					e.c.emit(&defineSetterKeyed***REMOVED***key: key, enumerable: true***REMOVED***)
				default:
					e.c.assert(false, e.offset, "unknown property kind: %s", prop.Kind)
					panic("unreachable")
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
			e.c.emit(putProp(key))
		case *ast.SpreadElement:
			e.c.compileExpression(prop.Expression).emitGetter(true)
			e.c.emit(copySpread)
		default:
			e.c.assert(false, e.offset, "unknown Property type: %T", prop)
			panic("unreachable")
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
				e.c.emitExpr(e.c.compileExpression(v), true)
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

func (c *compiler) emitCallee(callee compiledExpr) (calleeName unistring.String) ***REMOVED***
	switch callee := callee.(type) ***REMOVED***
	case *compiledDotExpr:
		callee.left.emitGetter(true)
		c.emit(getPropCallee(callee.name))
	case *compiledPrivateDotExpr:
		callee.left.emitGetter(true)
		rn, id := c.resolvePrivateName(callee.name, callee.offset)
		if rn != nil ***REMOVED***
			c.emit((*getPrivatePropResCallee)(rn))
		***REMOVED*** else ***REMOVED***
			c.emit((*getPrivatePropIdCallee)(id))
		***REMOVED***
	case *compiledSuperDotExpr:
		c.emitLoadThis()
		c.emit(loadSuper)
		c.emit(getPropRecvCallee(callee.name))
	case *compiledBracketExpr:
		callee.left.emitGetter(true)
		callee.member.emitGetter(true)
		c.emit(getElemCallee)
	case *compiledSuperBracketExpr:
		c.emitLoadThis()
		c.emit(loadSuper)
		callee.member.emitGetter(true)
		c.emit(getElemRecvCallee)
	case *compiledIdentifierExpr:
		calleeName = callee.name
		callee.emitGetterAndCallee()
	case *compiledOptionalChain:
		c.startOptChain()
		c.emitCallee(callee.expr)
		c.endOptChain()
	case *compiledOptional:
		c.emitCallee(callee.expr)
		c.block.conts = append(c.block.conts, len(c.p.code))
		c.emit(nil)
	case *compiledSuperExpr:
		// no-op
	default:
		c.emit(loadUndef)
		callee.emitGetter(true)
	***REMOVED***
	return
***REMOVED***

func (e *compiledCallExpr) emitGetter(putOnStack bool) ***REMOVED***
	if e.isVariadic ***REMOVED***
		e.c.emit(startVariadic)
	***REMOVED***
	calleeName := e.c.emitCallee(e.callee)

	for _, expr := range e.args ***REMOVED***
		expr.emitGetter(true)
	***REMOVED***

	e.addSrcMap()
	if _, ok := e.callee.(*compiledSuperExpr); ok ***REMOVED***
		b, eval := e.c.scope.lookupThis()
		e.c.assert(eval || b != nil, e.offset, "super call, but no 'this' binding")
		if eval ***REMOVED***
			e.c.emit(resolveThisDynamic***REMOVED******REMOVED***)
		***REMOVED*** else ***REMOVED***
			b.markAccessPoint()
			e.c.emit(resolveThisStack***REMOVED******REMOVED***)
		***REMOVED***
		if e.isVariadic ***REMOVED***
			e.c.emit(superCallVariadic)
		***REMOVED*** else ***REMOVED***
			e.c.emit(superCall(len(e.args)))
		***REMOVED***
	***REMOVED*** else if calleeName == "eval" ***REMOVED***
		foundVar := false
		for sc := e.c.scope; sc != nil; sc = sc.outer ***REMOVED***
			if !foundVar && (sc.variable || sc.isFunction()) ***REMOVED***
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

func (c *compiler) compileCallee(v ast.Expression) compiledExpr ***REMOVED***
	if sup, ok := v.(*ast.SuperExpression); ok ***REMOVED***
		if s := c.scope.nearestThis(); s != nil && s.funcType == funcDerivedCtor ***REMOVED***
			e := &compiledSuperExpr***REMOVED******REMOVED***
			e.init(c, sup.Idx)
			return e
		***REMOVED***
		c.throwSyntaxError(int(v.Idx0())-1, "'super' keyword unexpected here")
		panic("unreachable")
	***REMOVED***
	return c.compileExpression(v)
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
		callee:     c.compileCallee(v.Callee),
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
	if c.scope.strict && len(v.Literal) > 1 && v.Literal[0] == '0' && v.Literal[1] <= '7' && v.Literal[1] >= '0' ***REMOVED***
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
		c.assert(false, int(v.Idx)-1, "Unsupported number literal type: %T", v.Value)
		panic("unreachable")
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

func (c *compiler) compileTemplateLiteral(v *ast.TemplateLiteral) compiledExpr ***REMOVED***
	r := &compiledTemplateLiteral***REMOVED******REMOVED***
	if v.Tag != nil ***REMOVED***
		r.tag = c.compileExpression(v.Tag)
	***REMOVED***
	ce := make([]compiledExpr, len(v.Expressions))
	for i, expr := range v.Expressions ***REMOVED***
		ce[i] = c.compileExpression(expr)
	***REMOVED***
	r.expressions = ce
	r.elements = v.Elements
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

func (c *compiler) emitExpr(expr compiledExpr, putOnStack bool) ***REMOVED***
	if expr.constant() ***REMOVED***
		c.emitConst(expr, putOnStack)
	***REMOVED*** else ***REMOVED***
		expr.emitGetter(putOnStack)
	***REMOVED***
***REMOVED***

type namedEmitter interface ***REMOVED***
	emitNamed(name unistring.String)
***REMOVED***

func (c *compiler) emitNamed(expr compiledExpr, name unistring.String) ***REMOVED***
	if en, ok := expr.(namedEmitter); ok ***REMOVED***
		en.emitNamed(name)
	***REMOVED*** else ***REMOVED***
		expr.emitGetter(true)
	***REMOVED***
***REMOVED***

func (c *compiler) emitNamedOrConst(expr compiledExpr, name unistring.String) ***REMOVED***
	if expr.constant() ***REMOVED***
		c.emitConst(expr, true)
	***REMOVED*** else ***REMOVED***
		c.emitNamed(expr, name)
	***REMOVED***
***REMOVED***

func (e *compiledFunctionLiteral) emitNamed(name unistring.String) ***REMOVED***
	e.lhsName = name
	e.emitGetter(true)
***REMOVED***

func (e *compiledClassLiteral) emitNamed(name unistring.String) ***REMOVED***
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
		c.assert(false, int(pattern.Idx0())-1, "unsupported Pattern: %T", pattern)
		panic("unreachable")
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
	c.emit(iterate)
	for _, elt := range pattern.Elements ***REMOVED***
		switch elt := elt.(type) ***REMOVED***
		case nil:
			c.emit(iterGetNextOrUndef***REMOVED******REMOVED***, pop)
		case *ast.AssignExpression:
			c.emitAssign(elt.Left, c.compilePatternInitExpr(func() ***REMOVED***
				c.emit(iterGetNextOrUndef***REMOVED******REMOVED***)
			***REMOVED***, elt.Right, elt.Idx0()), emitAssign)
		default:
			c.emitAssign(elt, c.compileEmitterExpr(func() ***REMOVED***
				c.emit(iterGetNextOrUndef***REMOVED******REMOVED***)
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
		e.c.emitExpr(e.def, true)
		e.c.p.code[mark] = jdef(len(e.c.p.code) - mark)
	***REMOVED***
***REMOVED***

func (e *compiledPatternInitExpr) emitNamed(name unistring.String) ***REMOVED***
	e.emitSrc()
	if e.def != nil ***REMOVED***
		mark := len(e.c.p.code)
		e.c.emit(nil)
		e.c.emitNamedOrConst(e.def, name)
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

func (e *compiledSpreadCallArgument) emitGetter(putOnStack bool) ***REMOVED***
	e.expr.emitGetter(putOnStack)
	if putOnStack ***REMOVED***
		e.c.emit(pushSpread)
	***REMOVED***
***REMOVED***

func (c *compiler) startOptChain() ***REMOVED***
	c.block = &block***REMOVED***
		typ:   blockOptChain,
		outer: c.block,
	***REMOVED***
***REMOVED***

func (c *compiler) endOptChain() ***REMOVED***
	lbl := len(c.p.code)
	for _, item := range c.block.breaks ***REMOVED***
		c.p.code[item] = jopt(lbl - item)
	***REMOVED***
	for _, item := range c.block.conts ***REMOVED***
		c.p.code[item] = joptc(lbl - item)
	***REMOVED***
	c.block = c.block.outer
***REMOVED***

func (e *compiledOptionalChain) emitGetter(putOnStack bool) ***REMOVED***
	e.c.startOptChain()
	e.expr.emitGetter(true)
	e.c.endOptChain()
	if !putOnStack ***REMOVED***
		e.c.emit(pop)
	***REMOVED***
***REMOVED***

func (e *compiledOptional) emitGetter(putOnStack bool) ***REMOVED***
	e.expr.emitGetter(putOnStack)
	if putOnStack ***REMOVED***
		e.c.block.breaks = append(e.c.block.breaks, len(e.c.p.code))
		e.c.emit(nil)
	***REMOVED***
***REMOVED***
