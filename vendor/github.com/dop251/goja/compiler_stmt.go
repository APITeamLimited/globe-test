package goja

import (
	"fmt"
	"strconv"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
)

func (c *compiler) compileStatement(v ast.Statement, needResult bool) ***REMOVED***
	// log.Printf("compileStatement(): %T", v)

	switch v := v.(type) ***REMOVED***
	case *ast.BlockStatement:
		c.compileBlockStatement(v, needResult)
	case *ast.ExpressionStatement:
		c.compileExpressionStatement(v, needResult)
	case *ast.VariableStatement:
		c.compileVariableStatement(v, needResult)
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
		c.compileBranchStatement(v, needResult)
	case *ast.TryStatement:
		c.compileTryStatement(v)
		if needResult ***REMOVED***
			c.emit(loadUndef)
		***REMOVED***
	case *ast.ThrowStatement:
		c.compileThrowStatement(v)
	case *ast.SwitchStatement:
		c.compileSwitchStatement(v, needResult)
	case *ast.LabelledStatement:
		c.compileLabeledStatement(v, needResult)
	case *ast.EmptyStatement:
		c.compileEmptyStatement(needResult)
	case *ast.WithStatement:
		c.compileWithStatement(v, needResult)
	case *ast.DebuggerStatement:
	default:
		panic(fmt.Errorf("Unknown statement type: %T", v))
	***REMOVED***
***REMOVED***

func (c *compiler) compileLabeledStatement(v *ast.LabelledStatement, needResult bool) ***REMOVED***
	label := v.Label.Name
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
		c.compileGenericLabeledStatement(v.Statement, needResult, label)
	***REMOVED***
***REMOVED***

func (c *compiler) compileTryStatement(v *ast.TryStatement) ***REMOVED***
	if c.scope.strict && v.Catch != nil ***REMOVED***
		switch v.Catch.Parameter.Name ***REMOVED***
		case "arguments", "eval":
			c.throwSyntaxError(int(v.Catch.Parameter.Idx)-1, "Catch variable may not be eval or arguments in strict mode")
		***REMOVED***
	***REMOVED***
	c.block = &block***REMOVED***
		typ:   blockTry,
		outer: c.block,
	***REMOVED***
	lbl := len(c.p.code)
	c.emit(nil)
	c.compileStatement(v.Body, false)
	c.emit(halt)
	lbl2 := len(c.p.code)
	c.emit(nil)
	var catchOffset int
	dynamicCatch := true
	if v.Catch != nil ***REMOVED***
		dyn := nearestNonLexical(c.scope).dynamic
		accessed := c.scope.accessed
		c.newScope()
		c.scope.bindName(v.Catch.Parameter.Name)
		c.scope.lexical = true
		start := len(c.p.code)
		c.emit(nil)
		catchOffset = len(c.p.code) - lbl
		c.emit(enterCatch(v.Catch.Parameter.Name))
		c.compileStatement(v.Catch.Body, false)
		dyn1 := c.scope.dynamic
		accessed1 := c.scope.accessed
		c.popScope()
		if !dyn && !dyn1 && !accessed1 ***REMOVED***
			c.scope.accessed = accessed
			dynamicCatch = false
			code := c.p.code[start+1:]
			m := make(map[uint32]uint32)
			remap := func(instr uint32) uint32 ***REMOVED***
				level := instr >> 24
				idx := instr & 0x00FFFFFF
				if level > 0 ***REMOVED***
					level--
					return (level << 24) | idx
				***REMOVED*** else ***REMOVED***
					// remap
					newIdx, exists := m[idx]
					if !exists ***REMOVED***
						exname := unistring.String(" __tmp" + strconv.Itoa(c.scope.lastFreeTmp))
						c.scope.lastFreeTmp++
						newIdx, _ = c.scope.bindName(exname)
						m[idx] = newIdx
					***REMOVED***
					return newIdx
				***REMOVED***
			***REMOVED***
			for pc, instr := range code ***REMOVED***
				switch instr := instr.(type) ***REMOVED***
				case getLocal:
					code[pc] = getLocal(remap(uint32(instr)))
				case setLocal:
					code[pc] = setLocal(remap(uint32(instr)))
				case setLocalP:
					code[pc] = setLocalP(remap(uint32(instr)))
				***REMOVED***
			***REMOVED***
			c.p.code[start+1] = pop
			if catchVarIdx, exists := m[0]; exists ***REMOVED***
				c.p.code[start] = setLocal(catchVarIdx)
				catchOffset--
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			c.scope.accessed = true
		***REMOVED***

		/*
			if true/*sc.dynamic/ ***REMOVED***
				dynamicCatch = true
				c.scope.accessed = true
				c.newScope()
				c.scope.bindName(v.Catch.Parameter.Name)
				c.scope.lexical = true
				c.emit(enterCatch(v.Catch.Parameter.Name))
				c.compileStatement(v.Catch.Body, false)
				c.popScope()
			***REMOVED*** else ***REMOVED***
				exname := " __tmp" + strconv.Itoa(c.scope.lastFreeTmp)
				c.scope.lastFreeTmp++
				catchVarIdx, _ := c.scope.bindName(exname)
				c.emit(setLocal(catchVarIdx), pop)
				saved, wasSaved := c.scope.namesMap[v.Catch.Parameter.Name]
				c.scope.namesMap[v.Catch.Parameter.Name] = exname
				c.compileStatement(v.Catch.Body, false)
				if wasSaved ***REMOVED***
					c.scope.namesMap[v.Catch.Parameter.Name] = saved
				***REMOVED*** else ***REMOVED***
					delete(c.scope.namesMap, v.Catch.Parameter.Name)
				***REMOVED***
				c.scope.lastFreeTmp--
			***REMOVED****/
		c.emit(halt)
	***REMOVED***
	var finallyOffset int
	if v.Finally != nil ***REMOVED***
		lbl1 := len(c.p.code)
		c.emit(nil)
		finallyOffset = len(c.p.code) - lbl
		c.compileStatement(v.Finally, false)
		c.emit(halt, retFinally)
		c.p.code[lbl1] = jump(len(c.p.code) - lbl1)
	***REMOVED***
	c.p.code[lbl] = try***REMOVED***catchOffset: int32(catchOffset), finallyOffset: int32(finallyOffset), dynamic: dynamicCatch***REMOVED***
	c.p.code[lbl2] = jump(len(c.p.code) - lbl2)
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileThrowStatement(v *ast.ThrowStatement) ***REMOVED***
	//c.p.srcMap = append(c.p.srcMap, srcMapItem***REMOVED***pc: len(c.p.code), srcPos: int(v.Throw) - 1***REMOVED***)
	c.compileExpression(v.Argument).emitGetter(true)
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

	if needResult ***REMOVED***
		c.emit(jump(2))
	***REMOVED***
	start := len(c.p.code)
	if needResult ***REMOVED***
		c.emit(pop)
	***REMOVED***
	c.markBlockStart()
	c.compileStatement(v.Body, needResult)
	c.block.cont = len(c.p.code)
	c.emitExpr(c.compileExpression(v.Test), true)
	c.emit(jeq(start - len(c.p.code)))
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileForStatement(v *ast.ForStatement, needResult bool) ***REMOVED***
	c.compileLabeledForStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileLabeledForStatement(v *ast.ForStatement, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockLoop,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***

	if v.Initializer != nil ***REMOVED***
		c.compileExpression(v.Initializer).emitGetter(false)
	***REMOVED***
	if needResult ***REMOVED***
		c.emit(loadUndef) // initial result
	***REMOVED***
	start := len(c.p.code)
	c.markBlockStart()
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
					// TODO: Properly implement dummy compilation (no garbage in block, scope, etc..)
					/*
						p := c.p
						c.p = &program***REMOVED******REMOVED***
						c.compileStatement(v.Body, false)
						if v.Update != nil ***REMOVED***
							c.compileExpression(v.Update).emitGetter(false)
						***REMOVED***
						c.p = p*/
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
		c.emit(pop) // remove last result
	***REMOVED***
	c.markBlockStart()
	c.compileStatement(v.Body, needResult)
	c.block.cont = len(c.p.code)
	if v.Update != nil ***REMOVED***
		c.compileExpression(v.Update).emitGetter(false)
	***REMOVED***
	c.emit(jump(start - len(c.p.code)))
	if v.Test != nil ***REMOVED***
		if !testConst ***REMOVED***
			c.p.code[j] = jne(len(c.p.code) - j)
		***REMOVED***
	***REMOVED***
end:
	c.leaveBlock()
	c.markBlockStart()
***REMOVED***

func (c *compiler) compileForInStatement(v *ast.ForInStatement, needResult bool) ***REMOVED***
	c.compileLabeledForInStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileLabeledForInStatement(v *ast.ForInStatement, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockLoopEnum,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***

	c.compileExpression(v.Source).emitGetter(true)
	c.emit(enumerate)
	if needResult ***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
	start := len(c.p.code)
	c.markBlockStart()
	c.block.cont = start
	c.emit(nil)
	c.compileExpression(v.Into).emitSetter(&c.enumGetExpr)
	c.emit(pop)
	if needResult ***REMOVED***
		c.emit(pop) // remove last result
	***REMOVED***
	c.markBlockStart()
	c.compileStatement(v.Body, needResult)
	c.emit(jump(start - len(c.p.code)))
	c.p.code[start] = enumNext(len(c.p.code) - start)
	c.leaveBlock()
	c.markBlockStart()
	c.emit(enumPop)
***REMOVED***

func (c *compiler) compileForOfStatement(v *ast.ForOfStatement, needResult bool) ***REMOVED***
	c.compileLabeledForOfStatement(v, needResult, "")
***REMOVED***

func (c *compiler) compileLabeledForOfStatement(v *ast.ForOfStatement, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockLoopEnum,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***

	c.compileExpression(v.Source).emitGetter(true)
	c.emit(iterate)
	if needResult ***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
	start := len(c.p.code)
	c.markBlockStart()
	c.block.cont = start

	c.emit(nil)
	c.compileExpression(v.Into).emitSetter(&c.enumGetExpr)
	c.emit(pop)
	if needResult ***REMOVED***
		c.emit(pop) // remove last result
	***REMOVED***
	c.markBlockStart()
	c.compileStatement(v.Body, needResult)
	c.emit(jump(start - len(c.p.code)))
	c.p.code[start] = iterNext(len(c.p.code) - start)
	c.leaveBlock()
	c.markBlockStart()
	c.emit(enumPop)
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
		c.emit(loadUndef)
	***REMOVED***
	start := len(c.p.code)
	c.markBlockStart()
	c.block.cont = start
	expr := c.compileExpression(v.Test)
	testTrue := false
	var j int
	if expr.constant() ***REMOVED***
		if t, ex := c.evalConst(expr); ex == nil ***REMOVED***
			if t.ToBoolean() ***REMOVED***
				testTrue = true
			***REMOVED*** else ***REMOVED***
				p := c.p
				c.p = &Program***REMOVED******REMOVED***
				c.compileStatement(v.Body, false)
				c.p = p
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
		c.emit(pop)
	***REMOVED***
	c.markBlockStart()
	c.compileStatement(v.Body, needResult)
	c.emit(jump(start - len(c.p.code)))
	if !testTrue ***REMOVED***
		c.p.code[j] = jne(len(c.p.code) - j)
	***REMOVED***
end:
	c.leaveBlock()
	c.markBlockStart()
***REMOVED***

func (c *compiler) compileEmptyStatement(needResult bool) ***REMOVED***
	if needResult ***REMOVED***
		if len(c.p.code) == c.blockStart ***REMOVED***
			// first statement in block, use undefined as result
			c.emit(loadUndef)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileBranchStatement(v *ast.BranchStatement, needResult bool) ***REMOVED***
	switch v.Token ***REMOVED***
	case token.BREAK:
		c.compileBreak(v.Label, v.Idx)
	case token.CONTINUE:
		c.compileContinue(v.Label, v.Idx)
	default:
		panic(fmt.Errorf("Unknown branch statement token: %s", v.Token.String()))
	***REMOVED***
***REMOVED***

func (c *compiler) findBranchBlock(st *ast.BranchStatement) *block ***REMOVED***
	switch st.Token ***REMOVED***
	case token.BREAK:
		return c.findBreakBlock(st.Label)
	case token.CONTINUE:
		return c.findContinueBlock(st.Label)
	***REMOVED***
	return nil
***REMOVED***

func (c *compiler) findContinueBlock(label *ast.Identifier) (block *block) ***REMOVED***
	if label != nil ***REMOVED***
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			if (b.typ == blockLoop || b.typ == blockLoopEnum) && b.label == label.Name ***REMOVED***
				block = b
				break
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// find the nearest loop
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			if b.typ == blockLoop || b.typ == blockLoopEnum ***REMOVED***
				block = b
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (c *compiler) findBreakBlock(label *ast.Identifier) (block *block) ***REMOVED***
	if label != nil ***REMOVED***
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			if b.label == label.Name ***REMOVED***
				block = b
				break
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// find the nearest loop or switch
	L:
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			switch b.typ ***REMOVED***
			case blockLoop, blockLoopEnum, blockSwitch:
				block = b
				break L
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (c *compiler) compileBreak(label *ast.Identifier, idx file.Idx) ***REMOVED***
	var block *block
	if label != nil ***REMOVED***
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			switch b.typ ***REMOVED***
			case blockTry:
				c.emit(halt)
			case blockWith:
				c.emit(leaveWith)
			***REMOVED***
			if b.label == label.Name ***REMOVED***
				block = b
				break
			***REMOVED***
		***REMOVED***
		if block == nil ***REMOVED***
			c.throwSyntaxError(int(idx)-1, "Undefined label '%s'", label.Name)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// find the nearest loop or switch
	L:
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			switch b.typ ***REMOVED***
			case blockTry:
				c.emit(halt)
			case blockWith:
				c.emit(leaveWith)
			case blockLoop, blockLoopEnum, blockSwitch:
				block = b
				break L
			***REMOVED***
		***REMOVED***
		if block == nil ***REMOVED***
			c.throwSyntaxError(int(idx)-1, "Could not find block")
			return
		***REMOVED***
	***REMOVED***

	if len(c.p.code) == c.blockStart && block.needResult ***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
	block.breaks = append(block.breaks, len(c.p.code))
	c.emit(nil)
***REMOVED***

func (c *compiler) compileContinue(label *ast.Identifier, idx file.Idx) ***REMOVED***
	var block *block
	if label != nil ***REMOVED***
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			if b.typ == blockTry ***REMOVED***
				c.emit(halt)
			***REMOVED*** else if (b.typ == blockLoop || b.typ == blockLoopEnum) && b.label == label.Name ***REMOVED***
				block = b
				break
			***REMOVED***
		***REMOVED***
		if block == nil ***REMOVED***
			c.throwSyntaxError(int(idx)-1, "Undefined label '%s'", label.Name)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// find the nearest loop
		for b := c.block; b != nil; b = b.outer ***REMOVED***
			if b.typ == blockTry ***REMOVED***
				c.emit(halt)
			***REMOVED*** else if b.typ == blockLoop || b.typ == blockLoopEnum ***REMOVED***
				block = b
				break
			***REMOVED***
		***REMOVED***
		if block == nil ***REMOVED***
			c.throwSyntaxError(int(idx)-1, "Could not find block")
			return
		***REMOVED***
	***REMOVED***

	if len(c.p.code) == c.blockStart && block.needResult ***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
	block.conts = append(block.conts, len(c.p.code))
	c.emit(nil)
***REMOVED***

func (c *compiler) compileIfStatement(v *ast.IfStatement, needResult bool) ***REMOVED***
	test := c.compileExpression(v.Test)
	if test.constant() ***REMOVED***
		r, ex := c.evalConst(test)
		if ex != nil ***REMOVED***
			test.addSrcMap()
			c.emitThrow(ex.val)
			return
		***REMOVED***
		if r.ToBoolean() ***REMOVED***
			c.markBlockStart()
			c.compileStatement(v.Consequent, needResult)
			if v.Alternate != nil ***REMOVED***
				p := c.p
				c.p = &Program***REMOVED******REMOVED***
				c.markBlockStart()
				c.compileStatement(v.Alternate, false)
				c.p = p
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// TODO: Properly implement dummy compilation (no garbage in block, scope, etc..)
			p := c.p
			c.p = &Program***REMOVED******REMOVED***
			c.compileStatement(v.Consequent, false)
			c.p = p
			if v.Alternate != nil ***REMOVED***
				c.compileStatement(v.Alternate, needResult)
			***REMOVED*** else ***REMOVED***
				if needResult ***REMOVED***
					c.emit(loadUndef)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***
	test.emitGetter(true)
	jmp := len(c.p.code)
	c.emit(nil)
	c.markBlockStart()
	c.compileStatement(v.Consequent, needResult)
	if v.Alternate != nil ***REMOVED***
		jmp1 := len(c.p.code)
		c.emit(nil)
		c.p.code[jmp] = jne(len(c.p.code) - jmp)
		c.markBlockStart()
		c.compileStatement(v.Alternate, needResult)
		c.p.code[jmp1] = jump(len(c.p.code) - jmp1)
		c.markBlockStart()
	***REMOVED*** else ***REMOVED***
		if needResult ***REMOVED***
			c.emit(jump(2))
			c.p.code[jmp] = jne(len(c.p.code) - jmp)
			c.emit(loadUndef)
			c.markBlockStart()
		***REMOVED*** else ***REMOVED***
			c.p.code[jmp] = jne(len(c.p.code) - jmp)
			c.markBlockStart()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileReturnStatement(v *ast.ReturnStatement) ***REMOVED***
	if v.Argument != nil ***REMOVED***
		c.compileExpression(v.Argument).emitGetter(true)
		//c.emit(checkResolve)
	***REMOVED*** else ***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
	for b := c.block; b != nil; b = b.outer ***REMOVED***
		switch b.typ ***REMOVED***
		case blockTry:
			c.emit(halt)
		case blockLoopEnum:
			c.emit(enumPop)
		***REMOVED***
	***REMOVED***
	c.emit(ret)
***REMOVED***

func (c *compiler) compileVariableStatement(v *ast.VariableStatement, needResult bool) ***REMOVED***
	for _, expr := range v.List ***REMOVED***
		c.compileExpression(expr).emitGetter(false)
	***REMOVED***
	if needResult ***REMOVED***
		c.emit(loadUndef)
	***REMOVED***
***REMOVED***

func (c *compiler) getFirstNonEmptyStatement(st ast.Statement) ast.Statement ***REMOVED***
	switch st := st.(type) ***REMOVED***
	case *ast.BlockStatement:
		return c.getFirstNonEmptyStatementList(st.List)
	case *ast.LabelledStatement:
		return c.getFirstNonEmptyStatement(st.Statement)
	***REMOVED***
	return st
***REMOVED***

func (c *compiler) getFirstNonEmptyStatementList(list []ast.Statement) ast.Statement ***REMOVED***
	for _, st := range list ***REMOVED***
		switch st := st.(type) ***REMOVED***
		case *ast.EmptyStatement:
			continue
		case *ast.BlockStatement:
			return c.getFirstNonEmptyStatementList(st.List)
		case *ast.LabelledStatement:
			return c.getFirstNonEmptyStatement(st.Statement)
		***REMOVED***
		return st
	***REMOVED***
	return nil
***REMOVED***

func (c *compiler) compileStatements(list []ast.Statement, needResult bool) ***REMOVED***
	if len(list) > 0 ***REMOVED***
		cur := list[0]
		for idx := 0; idx < len(list); ***REMOVED***
			var next ast.Statement
			// find next non-empty statement
			for idx++; idx < len(list); idx++ ***REMOVED***
				if _, empty := list[idx].(*ast.EmptyStatement); !empty ***REMOVED***
					next = list[idx]
					break
				***REMOVED***
			***REMOVED***

			if next != nil ***REMOVED***
				bs := c.getFirstNonEmptyStatement(next)
				if bs, ok := bs.(*ast.BranchStatement); ok ***REMOVED***
					block := c.findBranchBlock(bs)
					if block != nil ***REMOVED***
						c.compileStatement(cur, block.needResult)
						cur = next
						continue
					***REMOVED***
				***REMOVED***
				c.compileStatement(cur, false)
				cur = next
			***REMOVED*** else ***REMOVED***
				c.compileStatement(cur, needResult)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if needResult ***REMOVED***
			c.emit(loadUndef)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileGenericLabeledStatement(v ast.Statement, needResult bool, label unistring.String) ***REMOVED***
	c.block = &block***REMOVED***
		typ:        blockBranch,
		outer:      c.block,
		label:      label,
		needResult: needResult,
	***REMOVED***
	c.compileStatement(v, needResult)
	c.leaveBlock()
***REMOVED***

func (c *compiler) compileBlockStatement(v *ast.BlockStatement, needResult bool) ***REMOVED***
	c.compileStatements(v.List, needResult)
***REMOVED***

func (c *compiler) compileExpressionStatement(v *ast.ExpressionStatement, needResult bool) ***REMOVED***
	expr := c.compileExpression(v.Expression)
	if expr.constant() ***REMOVED***
		c.emitConst(expr, needResult)
	***REMOVED*** else ***REMOVED***
		expr.emitGetter(needResult)
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
	c.newScope()
	c.scope.dynamic = true
	c.scope.lexical = true
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

	jumps := make([]int, len(v.Body))

	for i, s := range v.Body ***REMOVED***
		if s.Test != nil ***REMOVED***
			c.emit(dup)
			c.compileExpression(s.Test).emitGetter(true)
			c.emit(op_strict_eq)
			c.emit(jne(3), pop)
			jumps[i] = len(c.p.code)
			c.emit(nil)
		***REMOVED***
	***REMOVED***

	c.emit(pop)
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
			c.markBlockStart()
		***REMOVED***
		nr := false
		c.markBlockStart()
		if needResult ***REMOVED***
			if i < len(v.Body)-1 ***REMOVED***
				st := c.getFirstNonEmptyStatementList(v.Body[i+1].Consequent)
				if st, ok := st.(*ast.BranchStatement); ok && st.Token == token.BREAK ***REMOVED***
					if c.findBreakBlock(st.Label) != nil ***REMOVED***
						stmts := append(s.Consequent, st)
						c.compileStatements(stmts, false)
						continue
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				nr = true
			***REMOVED***
		***REMOVED***
		c.compileStatements(s.Consequent, nr)
	***REMOVED***
	if jumpNoMatch != -1 ***REMOVED***
		if needResult ***REMOVED***
			c.emit(jump(2))
		***REMOVED***
		c.p.code[jumpNoMatch] = jump(len(c.p.code) - jumpNoMatch)
		if needResult ***REMOVED***
			c.emit(loadUndef)
		***REMOVED***
	***REMOVED***
	c.leaveBlock()
	c.markBlockStart()
***REMOVED***
