package goja

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/unistring"
)

const (
	blockLoop = iota
	blockLoopEnum
	blockTry
	blockBranch
	blockSwitch
	blockWith
)

type CompilerError struct ***REMOVED***
	Message string
	File    *SrcFile
	Offset  int
***REMOVED***

type CompilerSyntaxError struct ***REMOVED***
	CompilerError
***REMOVED***

type CompilerReferenceError struct ***REMOVED***
	CompilerError
***REMOVED***

type srcMapItem struct ***REMOVED***
	pc     int
	srcPos int
***REMOVED***

type Program struct ***REMOVED***
	code   []instruction
	values []Value

	funcName unistring.String
	src      *SrcFile
	srcMap   []srcMapItem
***REMOVED***

type compiler struct ***REMOVED***
	p          *Program
	scope      *scope
	block      *block
	blockStart int

	enumGetExpr compiledEnumGetExpr

	evalVM *vm
***REMOVED***

type scope struct ***REMOVED***
	names      map[unistring.String]uint32
	outer      *scope
	strict     bool
	eval       bool
	lexical    bool
	dynamic    bool
	accessed   bool
	argsNeeded bool
	thisNeeded bool

	namesMap    map[unistring.String]unistring.String
	lastFreeTmp int
***REMOVED***

type block struct ***REMOVED***
	typ        int
	label      unistring.String
	needResult bool
	cont       int
	breaks     []int
	conts      []int
	outer      *block
***REMOVED***

func (c *compiler) leaveBlock() ***REMOVED***
	lbl := len(c.p.code)
	for _, item := range c.block.breaks ***REMOVED***
		c.p.code[item] = jump(lbl - item)
	***REMOVED***
	if t := c.block.typ; t == blockLoop || t == blockLoopEnum ***REMOVED***
		for _, item := range c.block.conts ***REMOVED***
			c.p.code[item] = jump(c.block.cont - item)
		***REMOVED***
	***REMOVED***
	c.block = c.block.outer
***REMOVED***

func (e *CompilerSyntaxError) Error() string ***REMOVED***
	if e.File != nil ***REMOVED***
		return fmt.Sprintf("SyntaxError: %s at %s", e.Message, e.File.Position(e.Offset))
	***REMOVED***
	return fmt.Sprintf("SyntaxError: %s", e.Message)
***REMOVED***

func (e *CompilerReferenceError) Error() string ***REMOVED***
	return fmt.Sprintf("ReferenceError: %s", e.Message)
***REMOVED***

func (c *compiler) newScope() ***REMOVED***
	strict := false
	if c.scope != nil ***REMOVED***
		strict = c.scope.strict
	***REMOVED***
	c.scope = &scope***REMOVED***
		outer:    c.scope,
		names:    make(map[unistring.String]uint32),
		strict:   strict,
		namesMap: make(map[unistring.String]unistring.String),
	***REMOVED***
***REMOVED***

func (c *compiler) popScope() ***REMOVED***
	c.scope = c.scope.outer
***REMOVED***

func newCompiler() *compiler ***REMOVED***
	c := &compiler***REMOVED***
		p: &Program***REMOVED******REMOVED***,
	***REMOVED***

	c.enumGetExpr.init(c, file.Idx(0))

	c.newScope()
	c.scope.dynamic = true
	return c
***REMOVED***

func (p *Program) defineLiteralValue(val Value) uint32 ***REMOVED***
	for idx, v := range p.values ***REMOVED***
		if v.SameAs(val) ***REMOVED***
			return uint32(idx)
		***REMOVED***
	***REMOVED***
	idx := uint32(len(p.values))
	p.values = append(p.values, val)
	return idx
***REMOVED***

func (p *Program) dumpCode(logger func(format string, args ...interface***REMOVED******REMOVED***)) ***REMOVED***
	p._dumpCode("", logger)
***REMOVED***

func (p *Program) _dumpCode(indent string, logger func(format string, args ...interface***REMOVED******REMOVED***)) ***REMOVED***
	logger("values: %+v", p.values)
	for pc, ins := range p.code ***REMOVED***
		logger("%s %d: %T(%v)", indent, pc, ins, ins)
		if f, ok := ins.(*newFunc); ok ***REMOVED***
			f.prg._dumpCode(indent+">", logger)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *Program) sourceOffset(pc int) int ***REMOVED***
	i := sort.Search(len(p.srcMap), func(idx int) bool ***REMOVED***
		return p.srcMap[idx].pc > pc
	***REMOVED***) - 1
	if i >= 0 ***REMOVED***
		return p.srcMap[i].srcPos
	***REMOVED***

	return 0
***REMOVED***

func (s *scope) isFunction() bool ***REMOVED***
	if !s.lexical ***REMOVED***
		return s.outer != nil
	***REMOVED***
	return s.outer.isFunction()
***REMOVED***

func (s *scope) lookupName(name unistring.String) (idx uint32, found, noDynamics bool) ***REMOVED***
	var level uint32 = 0
	noDynamics = true
	for curScope := s; curScope != nil; curScope = curScope.outer ***REMOVED***
		if curScope != s ***REMOVED***
			curScope.accessed = true
		***REMOVED***
		if curScope.dynamic ***REMOVED***
			noDynamics = false
		***REMOVED*** else ***REMOVED***
			var mapped unistring.String
			if m, exists := curScope.namesMap[name]; exists ***REMOVED***
				mapped = m
			***REMOVED*** else ***REMOVED***
				mapped = name
			***REMOVED***
			if i, exists := curScope.names[mapped]; exists ***REMOVED***
				idx = i | (level << 24)
				found = true
				return
			***REMOVED***
		***REMOVED***
		if name == "arguments" && !s.lexical && s.isFunction() ***REMOVED***
			s.argsNeeded = true
			s.accessed = true
			idx, _ = s.bindName(name)
			found = true
			return
		***REMOVED***
		level++
	***REMOVED***
	return
***REMOVED***

func (s *scope) bindName(name unistring.String) (uint32, bool) ***REMOVED***
	if s.lexical ***REMOVED***
		return s.outer.bindName(name)
	***REMOVED***

	if idx, exists := s.names[name]; exists ***REMOVED***
		return idx, false
	***REMOVED***
	idx := uint32(len(s.names))
	s.names[name] = idx
	return idx, true
***REMOVED***

func (s *scope) bindNameShadow(name unistring.String) (uint32, bool) ***REMOVED***
	if s.lexical ***REMOVED***
		return s.outer.bindName(name)
	***REMOVED***

	unique := true

	if idx, exists := s.names[name]; exists ***REMOVED***
		unique = false
		// shadow the var
		delete(s.names, name)
		n := unistring.String(strconv.Itoa(int(idx)))
		s.names[n] = idx
	***REMOVED***
	idx := uint32(len(s.names))
	s.names[name] = idx
	return idx, unique
***REMOVED***

func (c *compiler) markBlockStart() ***REMOVED***
	c.blockStart = len(c.p.code)
***REMOVED***

func (c *compiler) compile(in *ast.Program) ***REMOVED***
	c.p.src = NewSrcFile(in.File.Name(), in.File.Source(), in.SourceMap)

	if len(in.Body) > 0 ***REMOVED***
		if !c.scope.strict ***REMOVED***
			c.scope.strict = c.isStrict(in.Body)
		***REMOVED***
	***REMOVED***

	c.compileDeclList(in.DeclarationList, false)
	c.compileFunctions(in.DeclarationList)

	c.markBlockStart()
	c.compileStatements(in.Body, true)

	c.p.code = append(c.p.code, halt)
	code := c.p.code
	c.p.code = make([]instruction, 0, len(code)+len(c.scope.names)+2)
	if c.scope.eval ***REMOVED***
		if !c.scope.strict ***REMOVED***
			c.emit(jne(2), newStash)
		***REMOVED*** else ***REMOVED***
			c.emit(pop, newStash)
		***REMOVED***
	***REMOVED***
	l := len(c.p.code)
	c.p.code = c.p.code[:l+len(c.scope.names)]
	for name, nameIdx := range c.scope.names ***REMOVED***
		c.p.code[l+int(nameIdx)] = bindName(name)
	***REMOVED***

	c.p.code = append(c.p.code, code...)
	for i := range c.p.srcMap ***REMOVED***
		c.p.srcMap[i].pc += len(c.scope.names)
	***REMOVED***

***REMOVED***

func (c *compiler) compileDeclList(v []ast.Declaration, inFunc bool) ***REMOVED***
	for _, value := range v ***REMOVED***
		switch value := value.(type) ***REMOVED***
		case *ast.FunctionDeclaration:
			c.compileFunctionDecl(value)
		case *ast.VariableDeclaration:
			c.compileVarDecl(value, inFunc)
		default:
			panic(fmt.Errorf("Unsupported declaration: %T", value))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunctions(v []ast.Declaration) ***REMOVED***
	for _, value := range v ***REMOVED***
		if value, ok := value.(*ast.FunctionDeclaration); ok ***REMOVED***
			c.compileFunction(value)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileVarDecl(v *ast.VariableDeclaration, inFunc bool) ***REMOVED***
	for _, item := range v.List ***REMOVED***
		if c.scope.strict ***REMOVED***
			c.checkIdentifierLName(item.Name, int(item.Idx)-1)
			c.checkIdentifierName(item.Name, int(item.Idx)-1)
		***REMOVED***
		if !inFunc || item.Name != "arguments" ***REMOVED***
			idx, ok := c.scope.bindName(item.Name)
			_ = idx
			//log.Printf("Define var: %s: %x", item.Name, idx)
			if !ok ***REMOVED***
				// TODO: error
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) addDecls() []instruction ***REMOVED***
	code := make([]instruction, len(c.scope.names))
	for name, nameIdx := range c.scope.names ***REMOVED***
		code[nameIdx] = bindName(name)
	***REMOVED***
	return code
***REMOVED***

func (c *compiler) convertInstrToStashless(instr uint32, args int) (newIdx int, convert bool) ***REMOVED***
	level := instr >> 24
	idx := instr & 0x00FFFFFF
	if level > 0 ***REMOVED***
		level--
		newIdx = int((level << 24) | idx)
	***REMOVED*** else ***REMOVED***
		iidx := int(idx)
		if iidx < args ***REMOVED***
			newIdx = -iidx - 1
		***REMOVED*** else ***REMOVED***
			newIdx = iidx - args + 1
		***REMOVED***
		convert = true
	***REMOVED***
	return
***REMOVED***

func (c *compiler) convertFunctionToStashless(code []instruction, args int) ***REMOVED***
	code[0] = enterFuncStashless***REMOVED***stackSize: uint32(len(c.scope.names) - args), args: uint32(args)***REMOVED***
	for pc := 1; pc < len(code); pc++ ***REMOVED***
		instr := code[pc]
		if instr == ret ***REMOVED***
			code[pc] = retStashless
		***REMOVED***
		switch instr := instr.(type) ***REMOVED***
		case getLocal:
			if newIdx, convert := c.convertInstrToStashless(uint32(instr), args); convert ***REMOVED***
				code[pc] = loadStack(newIdx)
			***REMOVED*** else ***REMOVED***
				code[pc] = getLocal(newIdx)
			***REMOVED***
		case setLocal:
			if newIdx, convert := c.convertInstrToStashless(uint32(instr), args); convert ***REMOVED***
				code[pc] = storeStack(newIdx)
			***REMOVED*** else ***REMOVED***
				code[pc] = setLocal(newIdx)
			***REMOVED***
		case setLocalP:
			if newIdx, convert := c.convertInstrToStashless(uint32(instr), args); convert ***REMOVED***
				code[pc] = storeStackP(newIdx)
			***REMOVED*** else ***REMOVED***
				code[pc] = setLocalP(newIdx)
			***REMOVED***
		case getVar:
			level := instr.idx >> 24
			idx := instr.idx & 0x00FFFFFF
			level--
			instr.idx = level<<24 | idx
			code[pc] = instr
		case setVar:
			level := instr.idx >> 24
			idx := instr.idx & 0x00FFFFFF
			level--
			instr.idx = level<<24 | idx
			code[pc] = instr
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunctionDecl(v *ast.FunctionDeclaration) ***REMOVED***
	idx, ok := c.scope.bindName(v.Function.Name.Name)
	if !ok ***REMOVED***
		// TODO: error
	***REMOVED***
	_ = idx
	// log.Printf("Define function: %s: %x", v.Function.Name.Name, idx)
***REMOVED***

func (c *compiler) compileFunction(v *ast.FunctionDeclaration) ***REMOVED***
	e := &compiledIdentifierExpr***REMOVED***
		name: v.Function.Name.Name,
	***REMOVED***
	e.init(c, v.Function.Idx0())
	e.emitSetter(c.compileFunctionLiteral(v.Function, false))
	c.emit(pop)
***REMOVED***

func (c *compiler) emit(instructions ...instruction) ***REMOVED***
	c.p.code = append(c.p.code, instructions...)
***REMOVED***

func (c *compiler) throwSyntaxError(offset int, format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	panic(&CompilerSyntaxError***REMOVED***
		CompilerError: CompilerError***REMOVED***
			File:    c.p.src,
			Offset:  offset,
			Message: fmt.Sprintf(format, args...),
		***REMOVED***,
	***REMOVED***)
***REMOVED***

func (c *compiler) isStrict(list []ast.Statement) bool ***REMOVED***
	for _, st := range list ***REMOVED***
		if st, ok := st.(*ast.ExpressionStatement); ok ***REMOVED***
			if e, ok := st.Expression.(*ast.StringLiteral); ok ***REMOVED***
				if e.Literal == `"use strict"` || e.Literal == `'use strict'` ***REMOVED***
					return true
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (c *compiler) isStrictStatement(s ast.Statement) bool ***REMOVED***
	if s, ok := s.(*ast.BlockStatement); ok ***REMOVED***
		return c.isStrict(s.List)
	***REMOVED***
	return false
***REMOVED***

func (c *compiler) checkIdentifierName(name unistring.String, offset int) ***REMOVED***
	switch name ***REMOVED***
	case "implements", "interface", "let", "package", "private", "protected", "public", "static", "yield":
		c.throwSyntaxError(offset, "Unexpected strict mode reserved word")
	***REMOVED***
***REMOVED***

func (c *compiler) checkIdentifierLName(name unistring.String, offset int) ***REMOVED***
	switch name ***REMOVED***
	case "eval", "arguments":
		c.throwSyntaxError(offset, "Assignment to eval or arguments is not allowed in strict mode")
	***REMOVED***
***REMOVED***
