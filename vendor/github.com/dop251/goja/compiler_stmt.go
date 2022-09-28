package goja

import (
	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
)

func (c *compiler) compileStatement(v ast.Statement, needResult bool) ***REMOVED***

	switch v := v.(type) ***REMOVED***
	case *ast.BlockStatement:
		c.compileBlockStatement(v, needResult)
	case *ast.ExpressionStatement:
		c.compileExpressionStatement(v, needResult)
	case *ast.VariableStatement:
		c.compileVariableStatement(v)
	case *ast.LexicalDeclaration:
		c.compileLexicalDeclaration(v)
	case *ast.ReturnStatement:
		c.compileReturnStatement(v)
	case *ast.IfStatement:
		c.compileIfStatement(v, needResult)
	case *ast.DoWhileStatement:
		c.compileDoWhileStatement(v, needResult)
	case *ast.ForStatement:
		c.compileForStatement(v, needResult)
	case *ast.ForInStatement:
		c.compileForInStatement(v, needResult)
	case *ast.ForOfStatement:
		c.compileForOfStatement(v, needResult)
	case *ast.WhileStatement:
		c.compileWhileStatement(v, needResult)
	case *ast.BranchStatement:
		c.compileBranchStatement(v)
	case *ast.TryStatement:
		c.compileTryStatement(v, needResult)
	case *ast.ThrowStatement:
		c.compileThrowStatement(v)
	case *ast.SwitchStatement:
		c.compileSwitchStatement(v, needResult)
	case *ast.LabelledStatement:
		c.compileLabeledStatement(v, needResult)
	case *ast.EmptyStatement:
		c.compileEmptyStatement(needResult)
	case *ast.FunctionDeclaration:
		c.compileStandaloneFunctionDecl(v)
		// note functions inside blocks are hoisted to the top of the block and are compiled using compileFunctions()
	case *ast.ClassDeclaration:
		c.compileClassDeclaration(v)
	case *ast.WithStatement:
		c.compileWithStatement(v, needResult)
	case *ast.DebuggerStatement:
	default:
		c.assert(false, int(v.Idx0())-1, "Unknown statement type: %T", v)
		panic("unreachable")
	***REMOVED***
***REMOVED***

func (c *compiler) compileLabeledStatement(v *ast.LabelledStatement, needResult bool) ***REMOVED***
	label := v.Label.Name
	if c.scope.strict ***REMOVED***
		c.checkIdentifierName(label, int(v.Label.Idx)-1)
	***REMOVED***
	for b := c.block; b != nil; b = b.outer ***REMOVED***
		if b.label == label ***REMOVED***
			c.throwSyntaxError(int(v.Label.Idx-1), "Label '%s' has already been declared", label)
		***REMOVED***
	***REMOVED***
	switch s := v.Statement.(type) ***REMOVED***
	case *ast.ForInStatement:
		c.compileLabeledForInStatement(s, needResult, label)
	case *ast.ForOfStatement:
		c.compileLabeledForOfStatement(s, needResult, label)
	case *ast.ForStatement:
		c.compileLabeledForStatement(s, needResult, label)
	case *ast.WhileStatement:
		c.compileLabeledWhileStatement(s, needResult, label)
	case *ast.DoWhileStatement:
		c.compileLabeledDoWhileStatement(s, needResult, label)
	default:
		c.compileGenericLabeledStatement(s, needResult, label)
	***REMOVED***
***REMOVED***

func (c *compiler) updateEnterBlock(enter *enterBlock) ***REMOVED***
	scope := c.scope
	stashSize, stackSize := 0, 0
	if scope.dynLookup ***REMOVED***
		stashSize = len(scope.bindings)
		enter.names = scope.makeNamesMap()
	***REMOVED*** else ***REMOVED***
		for _, b := range scope.bindings ***REMOVED***
			if b.inStash ***REMOVED***
				stashSize++
			***REMOVED*** else ***REMOVED***
				stackSize++
			***REMOVED***
		***REMOVED***
	***REMOVED***
	enter.stashSize, enter.stackSize = uint32(stashSize), uint32(stackSize)
***REMOVED***

func (c *compiler) compileTryStatement(v *ast.TryStatement, needResult bool) ***REMOVED***
	c.block = &block***REMOVED***
		typ:   blockTry,
		outer: c.block,
	***REMOVED***
	var lp int
	var bodyNeedResult bool
	var finallyBreaking *block
	if v.Finally != nil ***REMOVED***
		lp, finallyBreaking = c.scanStatements(v.Finally.List)
	***REMOVED***
	if finallyBreaking != nil ***REMOVED***
		c.block.breaking = finallyBreaking
		if lp == -1 ***REMOVED***
			bodyNeedResult = finallyBreaking.needResult
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		bodyNeedResult = needResult
	***REMOVED***
	lbl := len(c.p.code)
	c.emit(nil)
	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
	c.compileBlockStatement(v.Body, bodyNeedResult)
	c.emit(halt)
	lbl2 := len(c.p.code)
	c.emit(nil)
	var catchOffset int
	if v.Catch != nil ***REMOVED***
		catchOffset = len(c.p.code) - lbl
		if v.Catch.Parameter != nil ***REMOVED***
			c.block = &block***REMOVED***
				typ:   blockScope,
				outer: c.block,
			***REMOVED***
			c.newBlockScope()
			list := v.Catch.Body.List
			funcs := c.extractFunctions(list)
			if _, ok := v.Catch.Parameter.(ast.Pattern); ok ***REMOVED***
				// add anonymous binding for the catch parameter, note it must be first
				c.scope.addBinding(int(v.Catch.Idx0()) - 1)
			***REMOVED***
			c.createBindings(v.Catch.Parameter, func(name unistring.String, offset int) ***REMOVED***
				if c.scope.strict ***REMOVED***
					switch name ***REMOVED***
					case "arguments", "eval":
						c.throwSyntaxError(offset, "Catch variable may not be eval or arguments in strict mode")
					***REMOVED***
				***REMOVED***
				c.scope.bindNameLexical(name, true, offset)
			***REMOVED***)
			enter := &enterBlock***REMOVED******REMOVED***
			c.emit(enter)
			if pattern, ok := v.Catch.Parameter.(ast.Pattern); ok ***REMOVED***
				c.scope.bindings[0].emitGet()
				c.emitPattern(pattern, func(target, init compiledExpr) ***REMOVED***
					c.emitPatternLexicalAssign(target, init)
				***REMOVED***, false)
			***REMOVED***
			for _, decl := range funcs ***REMOVED***
				c.scope.bindNameLexical(decl.Function.Name.Name, true, int(decl.Function.Name.Idx1())-1)
			***REMOVED***
			c.compileLexicalDeclarations(list, true)
			c.compileFunctions(funcs)
			c.compileStatements(list, bodyNeedResult)
			c.leaveScopeBlock(enter)
			if c.scope.dynLookup || c.scope.bindings[0].inStash ***REMOVED***
				c.p.code[lbl+catchOffset] = &enterCatchBlock***REMOVED***
					names:     enter.names,
					stashSize: enter.stashSize,
					stackSize: enter.stackSize,
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				enter.stackSize--
			***REMOVED***
			c.popScope()
		***REMOVED*** else ***REMOVED***
			c.emit(pop)
			c.compileBlockStatement(v.Catch.Body, bodyNeedResult)
		***REMOVED***
		c.emit(halt)
	***REMOVED***
	var finallyOffset int
	if v.Finally != nil ***REMOVED***
		lbl1 := len(c.p.code)
		c.emit(nil)
		finallyOffset = len(c.p.code) - lbl
		if bodyNeedResult && finallyBreaking != nil && lp == -1 ***REMOVED***
			c.emit(clearResult)
		***REMOVED***
		c.compileBlockStatement(v.Finally, false)
		c.emit(halt, retFinally)

		c.p.code[lbl1] = jump(len(c.p.code) - lbl1)
	***REMOVED***
	c.p.code[lbl] = try***REMOVED***catchOffset: int32(catchOffset), finallyOffset: int32(finallyOffset)***REMOVED***
	c.p.code[lbl2] = jump(len(c.p.code) - lbl2)
	c.leaveBlock()
***REMOVED***

func (c *compiler) addSrcMap(node ast.Node) ***REMOVED***
	c.p.addSrcMap(int(node.Idx0()) - 1)
***REMOVED***

func (c *compiler) compileThrowStatement(v *ast.ThrowStatement) ***REMOVED***
	c.compileExpression(v.Argument).emitGetter(true)
	c.addSrcMap(v)
	c.emit(throw)
***REMOVED***

func (c *compiler) compileDoWhileStatement(v *ast.DoWhileStatement, needResult bool) ***REMOVED***
	c.compileLabeledDoWhileStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileLabeledDoWhileStatement(v *ast.DoWhileStatement, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockLoop,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***

	start := len(c.p.code)
	c.compileStatement(v.Body, needResult)
	c.block.cont = len(c.p.code)
	c.emitExpr(c.compileExpression(v.Test), true)
	c.emit(jeq(start - len(c.p.code)))
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileForStatement(v *ast.ForStatement, needResult bool) ***REMOVED***
	c.compileLabeledForStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileForHeadLexDecl(decl *ast.LexicalDeclaration, needResult bool) *enterBlock ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockIterScope,
		outer:      c.block,
		needResult: needResult,
	***REMOVED***

	c.newBlockScope()
	enterIterBlock := &enterBlock***REMOVED******REMOVED***
	c.emit(enterIterBlock)
	c.createLexicalBindings(decl)
	c.compileLexicalDeclaration(decl)
	return enterIterBlock
***REMOVED***

func (c *compiler) compileLabeledForStatement(v *ast.ForStatement, needResult bool, label unistring.String) ***REMOVED***
	loopBlock := &block***REMOVED***
		typ:        blockLoop,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***
	c.block = loopBlock

	var enterIterBlock *enterBlock
	switch init := v.Initializer.(type) ***REMOVED***
	case nil:
		// no-op
	case *ast.ForLoopInitializerLexicalDecl:
		enterIterBlock = c.compileForHeadLexDecl(&init.LexicalDeclaration, needResult)
	case *ast.ForLoopInitializerVarDeclList:
		for _, expr := range init.List ***REMOVED***
			c.compileVarBinding(expr)
		***REMOVED***
	case *ast.ForLoopInitializerExpression:
		c.compileExpression(init.Expression).emitGetter(false)
	default:
		c.assert(false, int(v.For)-1, "Unsupported for loop initializer: %T", init)
		panic("unreachable")
	***REMOVED***

	if needResult ***REMOVED***
		c.emit(clearResult) // initial result
	***REMOVED***

	if enterIterBlock != nil ***REMOVED***
		c.emit(jump(1))
	***REMOVED***

	start := len(c.p.code)
	var j int
	testConst := false
	if v.Test != nil ***REMOVED***
		expr := c.compileExpression(v.Test)
		if expr.constant() ***REMOVED***
			r, ex := c.evalConst(expr)
			if ex == nil ***REMOVED***
				if r.ToBoolean() ***REMOVED***
					testConst = true
				***REMOVED*** else ***REMOVED***
					leave := c.enterDummyMode()
					c.compileStatement(v.Body, false)
					if v.Update != nil ***REMOVED***
						c.compileExpression(v.Update).emitGetter(false)
					***REMOVED***
					leave()
					goto end
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				expr.addSrcMap()
				c.emitThrow(ex.val)
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			expr.emitGetter(true)
			j = len(c.p.code)
			c.emit(nil)
		***REMOVED***
	***REMOVED***
	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
	c.compileStatement(v.Body, needResult)
	loopBlock.cont = len(c.p.code)
	if enterIterBlock != nil ***REMOVED***
		c.emit(jump(1))
	***REMOVED***
	if v.Update != nil ***REMOVED***
		c.compileExpression(v.Update).emitGetter(false)
	***REMOVED***
	if enterIterBlock != nil ***REMOVED***
		if c.scope.needStash || c.scope.isDynamic() ***REMOVED***
			c.p.code[start-1] = copyStash***REMOVED******REMOVED***
			c.p.code[loopBlock.cont] = copyStash***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			if l := len(c.p.code); l > loopBlock.cont ***REMOVED***
				loopBlock.cont++
			***REMOVED*** else ***REMOVED***
				c.p.code = c.p.code[:l-1]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	c.emit(jump(start - len(c.p.code)))
	if v.Test != nil ***REMOVED***
		if !testConst ***REMOVED***
			c.p.code[j] = jne(len(c.p.code) - j)
		***REMOVED***
	***REMOVED***
end:
	if enterIterBlock != nil ***REMOVED***
		c.leaveScopeBlock(enterIterBlock)
		c.popScope()
	***REMOVED***
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileForInStatement(v *ast.ForInStatement, needResult bool) ***REMOVED***
	c.compileLabeledForInStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileForInto(into ast.ForInto, needResult bool) (enter *enterBlock) ***REMOVED***
	switch into := into.(type) ***REMOVED***
	case *ast.ForIntoExpression:
		c.compileExpression(into.Expression).emitSetter(&c.enumGetExpr, false)
	case *ast.ForIntoVar:
		if c.scope.strict && into.Binding.Initializer != nil ***REMOVED***
			c.throwSyntaxError(int(into.Binding.Initializer.Idx0())-1, "for-in loop variable declaration may not have an initializer.")
		***REMOVED***
		switch target := into.Binding.Target.(type) ***REMOVED***
		case *ast.Identifier:
			c.compileIdentifierExpression(target).emitSetter(&c.enumGetExpr, false)
		case ast.Pattern:
			c.emit(enumGet)
			c.emitPattern(target, c.emitPatternVarAssign, false)
		default:
			c.throwSyntaxError(int(target.Idx0()-1), "unsupported for-in var target: %T", target)
		***REMOVED***
	case *ast.ForDeclaration:

		c.block = &block***REMOVED***
			typ:        blockIterScope,
			outer:      c.block,
			needResult: needResult,
		***REMOVED***

		c.newBlockScope()
		enter = &enterBlock***REMOVED******REMOVED***
		c.emit(enter)
		switch target := into.Target.(type) ***REMOVED***
		case *ast.Identifier:
			b := c.createLexicalIdBinding(target.Name, into.IsConst, int(into.Idx)-1)
			c.emit(enumGet)
			b.emitInitP()
		case ast.Pattern:
			c.createLexicalBinding(target, into.IsConst)
			c.emit(enumGet)
			c.emitPattern(target, func(target, init compiledExpr) ***REMOVED***
				c.emitPatternLexicalAssign(target, init)
			***REMOVED***, false)
		default:
			c.assert(false, int(into.Idx)-1, "Unsupported ForBinding: %T", into.Target)
		***REMOVED***
	default:
		c.assert(false, int(into.Idx0())-1, "Unsupported for-into: %T", into)
		panic("unreachable")
	***REMOVED***

	return
***REMOVED***

func (c *compiler) compileLabeledForInOfStatement(into ast.ForInto, source ast.Expression, body ast.Statement, iter, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockLoopEnum,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***
	enterPos := -1
	if forDecl, ok := into.(*ast.ForDeclaration); ok ***REMOVED***
		c.block = &block***REMOVED***
			typ:        blockScope,
			outer:      c.block,
			needResult: false,
		***REMOVED***
		c.newBlockScope()
		enterPos = len(c.p.code)
		c.emit(jump(1))
		c.createLexicalBinding(forDecl.Target, forDecl.IsConst)
	***REMOVED***
	c.compileExpression(source).emitGetter(true)
	if enterPos != -1 ***REMOVED***
		s := c.scope
		used := len(c.block.breaks) > 0 || s.isDynamic()
		if !used ***REMOVED***
			for _, b := range s.bindings ***REMOVED***
				if b.useCount() > 0 ***REMOVED***
					used = true
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if used ***REMOVED***
			// We need the stack untouched because it contains the source.
			// This is not the most optimal way, but it's an edge case, hopefully quite rare.
			for _, b := range s.bindings ***REMOVED***
				b.moveToStash()
			***REMOVED***
			enter := &enterBlock***REMOVED******REMOVED***
			c.p.code[enterPos] = enter
			c.leaveScopeBlock(enter)
		***REMOVED*** else ***REMOVED***
			c.block = c.block.outer
		***REMOVED***
		c.popScope()
	***REMOVED***
	if iter ***REMOVED***
		c.emit(iterateP)
	***REMOVED*** else ***REMOVED***
		c.emit(enumerate)
	***REMOVED***
	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
	start := len(c.p.code)
	c.block.cont = start
	c.emit(nil)
	enterIterBlock := c.compileForInto(into, needResult)
	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
	c.compileStatement(body, needResult)
	if enterIterBlock != nil ***REMOVED***
		c.leaveScopeBlock(enterIterBlock)
		c.popScope()
	***REMOVED***
	c.emit(jump(start - len(c.p.code)))
	if iter ***REMOVED***
		c.p.code[start] = iterNext(len(c.p.code) - start)
	***REMOVED*** else ***REMOVED***
		c.p.code[start] = enumNext(len(c.p.code) - start)
	***REMOVED***
	c.emit(enumPop, jump(2))
	c.leaveBlock()
	c.emit(enumPopClose)
***REMOVED***

func (c *compiler) compileLabeledForInStatement(v *ast.ForInStatement, needResult bool, label unistring.String) ***REMOVED***
	c.compileLabeledForInOfStatement(v.Into, v.Source, v.Body, false, needResult, label)
***REMOVED***

func (c *compiler) compileForOfStatement(v *ast.ForOfStatement, needResult bool) ***REMOVED***
	c.compileLabeledForOfStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileLabeledForOfStatement(v *ast.ForOfStatement, needResult bool, label unistring.String) ***REMOVED***
	c.compileLabeledForInOfStatement(v.Into, v.Source, v.Body, true, needResult, label)
***REMOVED***

func (c *compiler) compileWhileStatement(v *ast.WhileStatement, needResult bool) ***REMOVED***
	c.compileLabeledWhileStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileLabeledWhileStatement(v *ast.WhileStatement, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockLoop,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***

	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
	start := len(c.p.code)
	c.block.cont = start
	expr := c.compileExpression(v.Test)
	testTrue := false
	var j int
	if expr.constant() ***REMOVED***
		if t, ex := c.evalConst(expr); ex == nil ***REMOVED***
			if t.ToBoolean() ***REMOVED***
				testTrue = true
			***REMOVED*** else ***REMOVED***
				c.compileStatementDummy(v.Body)
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			c.emitThrow(ex.val)
			goto end
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		expr.emitGetter(true)
		j = len(c.p.code)
		c.emit(nil)
	***REMOVED***
	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
	c.compileStatement(v.Body, needResult)
	c.emit(jump(start - len(c.p.code)))
	if !testTrue ***REMOVED***
		c.p.code[j] = jne(len(c.p.code) - j)
	***REMOVED***
end:
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileEmptyStatement(needResult bool) ***REMOVED***
	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
***REMOVED***

func (c *compiler) compileBranchStatement(v *ast.BranchStatement) ***REMOVED***
	switch v.Token ***REMOVED***
	case token.BREAK:
		c.compileBreak(v.Label, v.Idx)
	case token.CONTINUE:
		c.compileContinue(v.Label, v.Idx)
	default:
		c.assert(false, int(v.Idx0())-1, "Unknown branch statement token: %s", v.Token.String())
		panic("unreachable")
	***REMOVED***
***REMOVED***

func (c *compiler) findBranchBlock(st *ast.BranchStatement) *block ***REMOVED***
	switch st.Token ***REMOVED***
	case token.BREAK:
		return c.findBreakBlock(st.Label, true)
	case token.CONTINUE:
		return c.findBreakBlock(st.Label, false)
	***REMOVED***
	return nil
***REMOVED***

func (c *compiler) findBreakBlock(label *ast.Identifier, isBreak bool) (res *block) ***REMOVED***
	if label != nil ***REMOVED***
		var found *block
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			if res == nil ***REMOVED***
				if bb := b.breaking; bb != nil ***REMOVED***
					res = bb
					if isBreak ***REMOVED***
						return
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if b.label == label.Name ***REMOVED***
				found = b
				break
			***REMOVED***
		***REMOVED***
		if !isBreak && found != nil && found.typ != blockLoop && found.typ != blockLoopEnum ***REMOVED***
			c.throwSyntaxError(int(label.Idx)-1, "Illegal continue statement: '%s' does not denote an iteration statement", label.Name)
		***REMOVED***
		if res == nil ***REMOVED***
			res = found
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// find the nearest loop or switch (if break)
	L:
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			if bb := b.breaking; bb != nil ***REMOVED***
				return bb
			***REMOVED***
			switch b.typ ***REMOVED***
			case blockLoop, blockLoopEnum:
				res = b
				break L
			case blockSwitch:
				if isBreak ***REMOVED***
					res = b
					break L
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (c *compiler) emitBlockExitCode(label *ast.Identifier, idx file.Idx, isBreak bool) *block ***REMOVED***
	block := c.findBreakBlock(label, isBreak)
	if block == nil ***REMOVED***
		c.throwSyntaxError(int(idx)-1, "Could not find block")
		panic("unreachable")
	***REMOVED***
L:
	for b := c.block; b != block; b = b.outer ***REMOVED***
		switch b.typ ***REMOVED***
		case blockIterScope:
			if !isBreak && b.outer == block ***REMOVED***
				break L
			***REMOVED***
			fallthrough
		case blockScope:
			b.breaks = append(b.breaks, len(c.p.code))
			c.emit(nil)
		case blockTry:
			c.emit(halt)
		case blockWith:
			c.emit(leaveWith)
		case blockLoopEnum:
			c.emit(enumPopClose)
		***REMOVED***
	***REMOVED***
	return block
***REMOVED***

func (c *compiler) compileBreak(label *ast.Identifier, idx file.Idx) ***REMOVED***
	block := c.emitBlockExitCode(label, idx, true)
	block.breaks = append(block.breaks, len(c.p.code))
	c.emit(nil)
***REMOVED***

func (c *compiler) compileContinue(label *ast.Identifier, idx file.Idx) ***REMOVED***
	block := c.emitBlockExitCode(label, idx, false)
	block.conts = append(block.conts, len(c.p.code))
	c.emit(nil)
***REMOVED***

func (c *compiler) compileIfBody(s ast.Statement, needResult bool) ***REMOVED***
	if !c.scope.strict ***REMOVED***
		if s, ok := s.(*ast.FunctionDeclaration); ok ***REMOVED***
			c.compileFunction(s)
			if needResult ***REMOVED***
				c.emit(clearResult)
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	c.compileStatement(s, needResult)
***REMOVED***

func (c *compiler) compileIfBodyDummy(s ast.Statement) ***REMOVED***
	leave := c.enterDummyMode()
	defer leave()
	c.compileIfBody(s, false)
***REMOVED***

func (c *compiler) compileIfStatement(v *ast.IfStatement, needResult bool) ***REMOVED***
	test := c.compileExpression(v.Test)
	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***
	if test.constant() ***REMOVED***
		r, ex := c.evalConst(test)
		if ex != nil ***REMOVED***
			test.addSrcMap()
			c.emitThrow(ex.val)
			return
		***REMOVED***
		if r.ToBoolean() ***REMOVED***
			c.compileIfBody(v.Consequent, needResult)
			if v.Alternate != nil ***REMOVED***
				c.compileIfBodyDummy(v.Alternate)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			c.compileIfBodyDummy(v.Consequent)
			if v.Alternate != nil ***REMOVED***
				c.compileIfBody(v.Alternate, needResult)
			***REMOVED*** else ***REMOVED***
				if needResult ***REMOVED***
					c.emit(clearResult)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***
	test.emitGetter(true)
	jmp := len(c.p.code)
	c.emit(nil)
	c.compileIfBody(v.Consequent, needResult)
	if v.Alternate != nil ***REMOVED***
		jmp1 := len(c.p.code)
		c.emit(nil)
		c.p.code[jmp] = jne(len(c.p.code) - jmp)
		c.compileIfBody(v.Alternate, needResult)
		c.p.code[jmp1] = jump(len(c.p.code) - jmp1)
	***REMOVED*** else ***REMOVED***
		if needResult ***REMOVED***
			c.emit(jump(2))
			c.p.code[jmp] = jne(len(c.p.code) - jmp)
			c.emit(clearResult)
		***REMOVED*** else ***REMOVED***
			c.p.code[jmp] = jne(len(c.p.code) - jmp)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileReturnStatement(v *ast.ReturnStatement) ***REMOVED***
	if s := c.scope.nearestFunction(); s != nil && s.funcType == funcClsInit ***REMOVED***
		c.throwSyntaxError(int(v.Return)-1, "Illegal return statement")
	***REMOVED***
	if v.Argument != nil ***REMOVED***
		c.emitExpr(c.compileExpression(v.Argument), true)
	***REMOVED*** else ***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
	for b := c.block; b != nil; b = b.outer ***REMOVED***
		switch b.typ ***REMOVED***
		case blockTry:
			c.emit(halt)
		case blockLoopEnum:
			c.emit(enumPopClose)
		***REMOVED***
	***REMOVED***
	if s := c.scope.nearestFunction(); s != nil && s.funcType == funcDerivedCtor ***REMOVED***
		b := s.boundNames[thisBindingName]
		c.assert(b != nil, int(v.Return)-1, "Derived constructor, but no 'this' binding")
		b.markAccessPoint()
	***REMOVED***
	c.emit(ret)
***REMOVED***

func (c *compiler) checkVarConflict(name unistring.String, offset int) ***REMOVED***
	for sc := c.scope; sc != nil; sc = sc.outer ***REMOVED***
		if b, exists := sc.boundNames[name]; exists && !b.isVar && !(b.isArg && sc != c.scope) ***REMOVED***
			c.throwSyntaxError(offset, "Identifier '%s' has already been declared", name)
		***REMOVED***
		if sc.isFunction() ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) emitVarAssign(name unistring.String, offset int, init compiledExpr) ***REMOVED***
	c.checkVarConflict(name, offset)
	if init != nil ***REMOVED***
		b, noDyn := c.scope.lookupName(name)
		if noDyn ***REMOVED***
			c.emitNamedOrConst(init, name)
			c.p.addSrcMap(offset)
			b.emitInitP()
		***REMOVED*** else ***REMOVED***
			c.emitVarRef(name, offset, b)
			c.emitNamedOrConst(init, name)
			c.p.addSrcMap(offset)
			c.emit(initValueP)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileVarBinding(expr *ast.Binding) ***REMOVED***
	switch target := expr.Target.(type) ***REMOVED***
	case *ast.Identifier:
		c.emitVarAssign(target.Name, int(target.Idx)-1, c.compileExpression(expr.Initializer))
	case ast.Pattern:
		c.compileExpression(expr.Initializer).emitGetter(true)
		c.emitPattern(target, c.emitPatternVarAssign, false)
	default:
		c.throwSyntaxError(int(target.Idx0()-1), "unsupported variable binding target: %T", target)
	***REMOVED***
***REMOVED***

func (c *compiler) emitLexicalAssign(name unistring.String, offset int, init compiledExpr) ***REMOVED***
	b := c.scope.boundNames[name]
	c.assert(b != nil, offset, "Lexical declaration for an unbound name")
	if init != nil ***REMOVED***
		c.emitNamedOrConst(init, name)
		c.p.addSrcMap(offset)
	***REMOVED*** else ***REMOVED***
		if b.isConst ***REMOVED***
			c.throwSyntaxError(offset, "Missing initializer in const declaration")
		***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
	b.emitInitP()
***REMOVED***

func (c *compiler) emitPatternVarAssign(target, init compiledExpr) ***REMOVED***
	id := target.(*compiledIdentifierExpr)
	c.emitVarAssign(id.name, id.offset, init)
***REMOVED***

func (c *compiler) emitPatternLexicalAssign(target, init compiledExpr) ***REMOVED***
	id := target.(*compiledIdentifierExpr)
	c.emitLexicalAssign(id.name, id.offset, init)
***REMOVED***

func (c *compiler) emitPatternAssign(target, init compiledExpr) ***REMOVED***
	if id, ok := target.(*compiledIdentifierExpr); ok ***REMOVED***
		b, noDyn := c.scope.lookupName(id.name)
		if noDyn ***REMOVED***
			c.emitNamedOrConst(init, id.name)
			b.emitSetP()
		***REMOVED*** else ***REMOVED***
			c.emitVarRef(id.name, id.offset, b)
			c.emitNamedOrConst(init, id.name)
			c.emit(putValueP)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		target.emitRef()
		c.emitExpr(init, true)
		c.emit(putValueP)
	***REMOVED***
***REMOVED***

func (c *compiler) compileLexicalBinding(expr *ast.Binding) ***REMOVED***
	switch target := expr.Target.(type) ***REMOVED***
	case *ast.Identifier:
		c.emitLexicalAssign(target.Name, int(target.Idx)-1, c.compileExpression(expr.Initializer))
	case ast.Pattern:
		c.compileExpression(expr.Initializer).emitGetter(true)
		c.emitPattern(target, func(target, init compiledExpr) ***REMOVED***
			c.emitPatternLexicalAssign(target, init)
		***REMOVED***, false)
	default:
		c.throwSyntaxError(int(target.Idx0()-1), "unsupported lexical binding target: %T", target)
	***REMOVED***
***REMOVED***

func (c *compiler) compileVariableStatement(v *ast.VariableStatement) ***REMOVED***
	for _, expr := range v.List ***REMOVED***
		c.compileVarBinding(expr)
	***REMOVED***
***REMOVED***

func (c *compiler) compileLexicalDeclaration(v *ast.LexicalDeclaration) ***REMOVED***
	for _, e := range v.List ***REMOVED***
		c.compileLexicalBinding(e)
	***REMOVED***
***REMOVED***

func (c *compiler) isEmptyResult(st ast.Statement) bool ***REMOVED***
	switch st := st.(type) ***REMOVED***
	case *ast.EmptyStatement, *ast.VariableStatement, *ast.LexicalDeclaration, *ast.FunctionDeclaration,
		*ast.ClassDeclaration, *ast.BranchStatement, *ast.DebuggerStatement:
		return true
	case *ast.LabelledStatement:
		return c.isEmptyResult(st.Statement)
	case *ast.BlockStatement:
		for _, s := range st.List ***REMOVED***
			if _, ok := s.(*ast.BranchStatement); ok ***REMOVED***
				return true
			***REMOVED***
			if !c.isEmptyResult(s) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (c *compiler) scanStatements(list []ast.Statement) (lastProducingIdx int, breakingBlock *block) ***REMOVED***
	lastProducingIdx = -1
	for i, st := range list ***REMOVED***
		if bs, ok := st.(*ast.BranchStatement); ok ***REMOVED***
			if blk := c.findBranchBlock(bs); blk != nil ***REMOVED***
				breakingBlock = blk
			***REMOVED***
			break
		***REMOVED***
		if !c.isEmptyResult(st) ***REMOVED***
			lastProducingIdx = i
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (c *compiler) compileStatementsNeedResult(list []ast.Statement, lastProducingIdx int) ***REMOVED***
	if lastProducingIdx >= 0 ***REMOVED***
		for _, st := range list[:lastProducingIdx] ***REMOVED***
			if _, ok := st.(*ast.FunctionDeclaration); ok ***REMOVED***
				continue
			***REMOVED***
			c.compileStatement(st, false)
		***REMOVED***
		c.compileStatement(list[lastProducingIdx], true)
	***REMOVED***
	var leave func()
	defer func() ***REMOVED***
		if leave != nil ***REMOVED***
			leave()
		***REMOVED***
	***REMOVED***()
	for _, st := range list[lastProducingIdx+1:] ***REMOVED***
		if _, ok := st.(*ast.FunctionDeclaration); ok ***REMOVED***
			continue
		***REMOVED***
		c.compileStatement(st, false)
		if leave == nil ***REMOVED***
			if _, ok := st.(*ast.BranchStatement); ok ***REMOVED***
				leave = c.enterDummyMode()
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileStatements(list []ast.Statement, needResult bool) ***REMOVED***
	lastProducingIdx, blk := c.scanStatements(list)
	if blk != nil ***REMOVED***
		needResult = blk.needResult
	***REMOVED***
	if needResult ***REMOVED***
		c.compileStatementsNeedResult(list, lastProducingIdx)
		return
	***REMOVED***
	for _, st := range list ***REMOVED***
		if _, ok := st.(*ast.FunctionDeclaration); ok ***REMOVED***
			continue
		***REMOVED***
		c.compileStatement(st, false)
	***REMOVED***
***REMOVED***

func (c *compiler) compileGenericLabeledStatement(v ast.Statement, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockLabel,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***
	c.compileStatement(v, needResult)
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileBlockStatement(v *ast.BlockStatement, needResult bool) ***REMOVED***
	var scopeDeclared bool
	funcs := c.extractFunctions(v.List)
	if len(funcs) > 0 ***REMOVED***
		c.newBlockScope()
		scopeDeclared = true
	***REMOVED***
	c.createFunctionBindings(funcs)
	scopeDeclared = c.compileLexicalDeclarations(v.List, scopeDeclared)

	var enter *enterBlock
	if scopeDeclared ***REMOVED***
		c.block = &block***REMOVED***
			outer:      c.block,
			typ:        blockScope,
			needResult: needResult,
		***REMOVED***
		enter = &enterBlock***REMOVED******REMOVED***
		c.emit(enter)
	***REMOVED***
	c.compileFunctions(funcs)
	c.compileStatements(v.List, needResult)
	if scopeDeclared ***REMOVED***
		c.leaveScopeBlock(enter)
		c.popScope()
	***REMOVED***
***REMOVED***

func (c *compiler) compileExpressionStatement(v *ast.ExpressionStatement, needResult bool) ***REMOVED***
	c.emitExpr(c.compileExpression(v.Expression), needResult)
	if needResult ***REMOVED***
		c.emit(saveResult)
	***REMOVED***
***REMOVED***

func (c *compiler) compileWithStatement(v *ast.WithStatement, needResult bool) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.throwSyntaxError(int(v.With)-1, "Strict mode code may not include a with statement")
		return
	***REMOVED***
	c.compileExpression(v.Object).emitGetter(true)
	c.emit(enterWith)
	c.block = &block***REMOVED***
		outer:      c.block,
		typ:        blockWith,
		needResult: needResult,
	***REMOVED***
	c.newBlockScope()
	c.scope.dynamic = true
	c.compileStatement(v.Body, needResult)
	c.emit(leaveWith)
	c.leaveBlock()
	c.popScope()
***REMOVED***

func (c *compiler) compileSwitchStatement(v *ast.SwitchStatement, needResult bool) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockSwitch,
		outer:      c.block,
		needResult: needResult,
	***REMOVED***

	c.compileExpression(v.Discriminant).emitGetter(true)

	var funcs []*ast.FunctionDeclaration
	for _, s := range v.Body ***REMOVED***
		f := c.extractFunctions(s.Consequent)
		funcs = append(funcs, f...)
	***REMOVED***
	var scopeDeclared bool
	if len(funcs) > 0 ***REMOVED***
		c.newBlockScope()
		scopeDeclared = true
		c.createFunctionBindings(funcs)
	***REMOVED***

	for _, s := range v.Body ***REMOVED***
		scopeDeclared = c.compileLexicalDeclarations(s.Consequent, scopeDeclared)
	***REMOVED***

	var enter *enterBlock
	var db *binding
	if scopeDeclared ***REMOVED***
		c.block = &block***REMOVED***
			typ:        blockScope,
			outer:      c.block,
			needResult: needResult,
		***REMOVED***
		enter = &enterBlock***REMOVED******REMOVED***
		c.emit(enter)
		// create anonymous variable for the discriminant
		bindings := c.scope.bindings
		var bb []*binding
		if cap(bindings) == len(bindings) ***REMOVED***
			bb = make([]*binding, len(bindings)+1)
		***REMOVED*** else ***REMOVED***
			bb = bindings[:len(bindings)+1]
		***REMOVED***
		copy(bb[1:], bindings)
		db = &binding***REMOVED***
			scope:    c.scope,
			isConst:  true,
			isStrict: true,
		***REMOVED***
		bb[0] = db
		c.scope.bindings = bb
	***REMOVED***

	c.compileFunctions(funcs)

	if needResult ***REMOVED***
		c.emit(clearResult)
	***REMOVED***

	jumps := make([]int, len(v.Body))

	for i, s := range v.Body ***REMOVED***
		if s.Test != nil ***REMOVED***
			if db != nil ***REMOVED***
				db.emitGet()
			***REMOVED*** else ***REMOVED***
				c.emit(dup)
			***REMOVED***
			c.compileExpression(s.Test).emitGetter(true)
			c.emit(op_strict_eq)
			if db != nil ***REMOVED***
				c.emit(jne(2))
			***REMOVED*** else ***REMOVED***
				c.emit(jne(3), pop)
			***REMOVED***
			jumps[i] = len(c.p.code)
			c.emit(nil)
		***REMOVED***
	***REMOVED***

	if db == nil ***REMOVED***
		c.emit(pop)
	***REMOVED***
	jumpNoMatch := -1
	if v.Default != -1 ***REMOVED***
		if v.Default != 0 ***REMOVED***
			jumps[v.Default] = len(c.p.code)
			c.emit(nil)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		jumpNoMatch = len(c.p.code)
		c.emit(nil)
	***REMOVED***

	for i, s := range v.Body ***REMOVED***
		if s.Test != nil || i != 0 ***REMOVED***
			c.p.code[jumps[i]] = jump(len(c.p.code) - jumps[i])
		***REMOVED***
		c.compileStatements(s.Consequent, needResult)
	***REMOVED***

	if jumpNoMatch != -1 ***REMOVED***
		c.p.code[jumpNoMatch] = jump(len(c.p.code) - jumpNoMatch)
	***REMOVED***
	if enter != nil ***REMOVED***
		c.leaveScopeBlock(enter)
		enter.stackSize--
		c.popScope()
	***REMOVED***
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileClassDeclaration(v *ast.ClassDeclaration) ***REMOVED***
	c.emitLexicalAssign(v.Class.Name.Name, int(v.Class.Class)-1, c.compileClassLiteral(v.Class, false))
***REMOVED***
