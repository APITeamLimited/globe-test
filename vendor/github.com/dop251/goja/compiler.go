package goja

import (
	"fmt"
	"github.com/dop251/goja/token"
	"sort"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/unistring"
)

type blockType int

const (
	blockLoop blockType = iota
	blockLoopEnum
	blockTry
	blockLabel
	blockSwitch
	blockWith
	blockScope
	blockIterScope
	blockOptChain
)

const (
	maskConst     = 1 << 31
	maskVar       = 1 << 30
	maskDeletable = 1 << 29
	maskStrict    = maskDeletable

	maskTyp = maskConst | maskVar | maskDeletable
)

type varType byte

const (
	varTypeVar varType = iota
	varTypeLet
	varTypeStrictConst
	varTypeConst
)

const thisBindingName = " this" // must not be a valid identifier

type CompilerError struct ***REMOVED***
	Message string
	File    *file.File
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
	src      *file.File
	srcMap   []srcMapItem
***REMOVED***

type compiler struct ***REMOVED***
	p     *Program
	scope *scope
	block *block

	classScope *classScope

	enumGetExpr compiledEnumGetExpr

	evalVM *vm // VM used to evaluate constant expressions
	ctxVM  *vm // VM in which an eval() code is compiled
***REMOVED***

type binding struct ***REMOVED***
	scope        *scope
	name         unistring.String
	accessPoints map[*scope]*[]int
	isConst      bool
	isStrict     bool
	isArg        bool
	isVar        bool
	inStash      bool
***REMOVED***

func (b *binding) getAccessPointsForScope(s *scope) *[]int ***REMOVED***
	m := b.accessPoints[s]
	if m == nil ***REMOVED***
		a := make([]int, 0, 1)
		m = &a
		if b.accessPoints == nil ***REMOVED***
			b.accessPoints = make(map[*scope]*[]int)
		***REMOVED***
		b.accessPoints[s] = m
	***REMOVED***
	return m
***REMOVED***

func (b *binding) markAccessPointAt(pos int) ***REMOVED***
	scope := b.scope.c.scope
	m := b.getAccessPointsForScope(scope)
	*m = append(*m, pos-scope.base)
***REMOVED***

func (b *binding) markAccessPointAtScope(scope *scope, pos int) ***REMOVED***
	m := b.getAccessPointsForScope(scope)
	*m = append(*m, pos-scope.base)
***REMOVED***

func (b *binding) markAccessPoint() ***REMOVED***
	scope := b.scope.c.scope
	m := b.getAccessPointsForScope(scope)
	*m = append(*m, len(scope.prg.code)-scope.base)
***REMOVED***

func (b *binding) emitGet() ***REMOVED***
	b.markAccessPoint()
	if b.isVar && !b.isArg ***REMOVED***
		b.scope.c.emit(loadStack(0))
	***REMOVED*** else ***REMOVED***
		b.scope.c.emit(loadStackLex(0))
	***REMOVED***
***REMOVED***

func (b *binding) emitGetAt(pos int) ***REMOVED***
	b.markAccessPointAt(pos)
	if b.isVar && !b.isArg ***REMOVED***
		b.scope.c.p.code[pos] = loadStack(0)
	***REMOVED*** else ***REMOVED***
		b.scope.c.p.code[pos] = loadStackLex(0)
	***REMOVED***
***REMOVED***

func (b *binding) emitGetP() ***REMOVED***
	if b.isVar && !b.isArg ***REMOVED***
		// no-op
	***REMOVED*** else ***REMOVED***
		// make sure TDZ is checked
		b.markAccessPoint()
		b.scope.c.emit(loadStackLex(0), pop)
	***REMOVED***
***REMOVED***

func (b *binding) emitSet() ***REMOVED***
	if b.isConst ***REMOVED***
		if b.isStrict || b.scope.c.scope.strict ***REMOVED***
			b.scope.c.emit(throwAssignToConst)
		***REMOVED***
		return
	***REMOVED***
	b.markAccessPoint()
	if b.isVar && !b.isArg ***REMOVED***
		b.scope.c.emit(storeStack(0))
	***REMOVED*** else ***REMOVED***
		b.scope.c.emit(storeStackLex(0))
	***REMOVED***
***REMOVED***

func (b *binding) emitSetP() ***REMOVED***
	if b.isConst ***REMOVED***
		if b.isStrict || b.scope.c.scope.strict ***REMOVED***
			b.scope.c.emit(throwAssignToConst)
		***REMOVED***
		return
	***REMOVED***
	b.markAccessPoint()
	if b.isVar && !b.isArg ***REMOVED***
		b.scope.c.emit(storeStackP(0))
	***REMOVED*** else ***REMOVED***
		b.scope.c.emit(storeStackLexP(0))
	***REMOVED***
***REMOVED***

func (b *binding) emitInitP() ***REMOVED***
	if !b.isVar && b.scope.outer == nil ***REMOVED***
		b.scope.c.emit(initGlobalP(b.name))
	***REMOVED*** else ***REMOVED***
		b.markAccessPoint()
		b.scope.c.emit(initStackP(0))
	***REMOVED***
***REMOVED***

func (b *binding) emitInit() ***REMOVED***
	if !b.isVar && b.scope.outer == nil ***REMOVED***
		b.scope.c.emit(initGlobal(b.name))
	***REMOVED*** else ***REMOVED***
		b.markAccessPoint()
		b.scope.c.emit(initStack(0))
	***REMOVED***
***REMOVED***

func (b *binding) emitInitAt(pos int) ***REMOVED***
	if !b.isVar && b.scope.outer == nil ***REMOVED***
		b.scope.c.p.code[pos] = initGlobal(b.name)
	***REMOVED*** else ***REMOVED***
		b.markAccessPointAt(pos)
		b.scope.c.p.code[pos] = initStack(0)
	***REMOVED***
***REMOVED***

func (b *binding) emitInitAtScope(scope *scope, pos int) ***REMOVED***
	if !b.isVar && scope.outer == nil ***REMOVED***
		scope.c.p.code[pos] = initGlobal(b.name)
	***REMOVED*** else ***REMOVED***
		b.markAccessPointAtScope(scope, pos)
		scope.c.p.code[pos] = initStack(0)
	***REMOVED***
***REMOVED***

func (b *binding) emitInitPAtScope(scope *scope, pos int) ***REMOVED***
	if !b.isVar && scope.outer == nil ***REMOVED***
		scope.c.p.code[pos] = initGlobalP(b.name)
	***REMOVED*** else ***REMOVED***
		b.markAccessPointAtScope(scope, pos)
		scope.c.p.code[pos] = initStackP(0)
	***REMOVED***
***REMOVED***

func (b *binding) emitGetVar(callee bool) ***REMOVED***
	b.markAccessPoint()
	if b.isVar && !b.isArg ***REMOVED***
		b.scope.c.emit(&loadMixed***REMOVED***name: b.name, callee: callee***REMOVED***)
	***REMOVED*** else ***REMOVED***
		b.scope.c.emit(&loadMixedLex***REMOVED***name: b.name, callee: callee***REMOVED***)
	***REMOVED***
***REMOVED***

func (b *binding) emitResolveVar(strict bool) ***REMOVED***
	b.markAccessPoint()
	if b.isVar && !b.isArg ***REMOVED***
		b.scope.c.emit(&resolveMixed***REMOVED***name: b.name, strict: strict, typ: varTypeVar***REMOVED***)
	***REMOVED*** else ***REMOVED***
		var typ varType
		if b.isConst ***REMOVED***
			if b.isStrict ***REMOVED***
				typ = varTypeStrictConst
			***REMOVED*** else ***REMOVED***
				typ = varTypeConst
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			typ = varTypeLet
		***REMOVED***
		b.scope.c.emit(&resolveMixed***REMOVED***name: b.name, strict: strict, typ: typ***REMOVED***)
	***REMOVED***
***REMOVED***

func (b *binding) moveToStash() ***REMOVED***
	if b.isArg && !b.scope.argsInStash ***REMOVED***
		b.scope.moveArgsToStash()
	***REMOVED*** else ***REMOVED***
		b.inStash = true
		b.scope.needStash = true
	***REMOVED***
***REMOVED***

func (b *binding) useCount() (count int) ***REMOVED***
	for _, a := range b.accessPoints ***REMOVED***
		count += len(*a)
	***REMOVED***
	return
***REMOVED***

type scope struct ***REMOVED***
	c          *compiler
	prg        *Program
	outer      *scope
	nested     []*scope
	boundNames map[unistring.String]*binding
	bindings   []*binding
	base       int
	numArgs    int

	// function type. If not funcNone, this is a function or a top-level lexical environment
	funcType funcType

	// in strict mode
	strict bool
	// eval top-level scope
	eval bool
	// at least one inner scope has direct eval() which can lookup names dynamically (by name)
	dynLookup bool
	// at least one binding has been marked for placement in stash
	needStash bool

	// is a variable environment, i.e. the target for dynamically created var bindings
	variable bool
	// a function scope that has at least one direct eval() and non-strict, so the variables can be added dynamically
	dynamic bool
	// arguments have been marked for placement in stash (functions only)
	argsInStash bool
	// need 'arguments' object (functions only)
	argsNeeded bool
***REMOVED***

type block struct ***REMOVED***
	typ        blockType
	label      unistring.String
	cont       int
	breaks     []int
	conts      []int
	outer      *block
	breaking   *block // set when the 'finally' block is an empty break statement sequence
	needResult bool
***REMOVED***

func (c *compiler) leaveScopeBlock(enter *enterBlock) ***REMOVED***
	c.updateEnterBlock(enter)
	leave := &leaveBlock***REMOVED***
		stackSize: enter.stackSize,
		popStash:  enter.stashSize > 0,
	***REMOVED***
	c.emit(leave)
	for _, pc := range c.block.breaks ***REMOVED***
		c.p.code[pc] = leave
	***REMOVED***
	c.block.breaks = nil
	c.leaveBlock()
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
		c:      c,
		prg:    c.p,
		outer:  c.scope,
		strict: strict,
	***REMOVED***
***REMOVED***

func (c *compiler) newBlockScope() ***REMOVED***
	c.newScope()
	if outer := c.scope.outer; outer != nil ***REMOVED***
		outer.nested = append(outer.nested, c.scope)
	***REMOVED***
	c.scope.base = len(c.p.code)
***REMOVED***

func (c *compiler) popScope() ***REMOVED***
	c.scope = c.scope.outer
***REMOVED***

func newCompiler() *compiler ***REMOVED***
	c := &compiler***REMOVED***
		p: &Program***REMOVED******REMOVED***,
	***REMOVED***

	c.enumGetExpr.init(c, file.Idx(0))

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
	dumpInitFields := func(initFields *Program) ***REMOVED***
		i := indent + ">"
		logger("%s ---- init_fields:", i)
		initFields._dumpCode(i, logger)
		logger("%s ----", i)
	***REMOVED***
	for pc, ins := range p.code ***REMOVED***
		logger("%s %d: %T(%v)", indent, pc, ins, ins)
		var prg *Program
		switch f := ins.(type) ***REMOVED***
		case *newFunc:
			prg = f.prg
		case *newArrowFunc:
			prg = f.prg
		case *newMethod:
			prg = f.prg
		case *newDerivedClass:
			if f.initFields != nil ***REMOVED***
				dumpInitFields(f.initFields)
			***REMOVED***
			prg = f.ctor
		case *newClass:
			if f.initFields != nil ***REMOVED***
				dumpInitFields(f.initFields)
			***REMOVED***
		case *newStaticFieldInit:
			if f.initFields != nil ***REMOVED***
				dumpInitFields(f.initFields)
			***REMOVED***
		***REMOVED***
		if prg != nil ***REMOVED***
			prg._dumpCode(indent+">", logger)
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

func (p *Program) addSrcMap(srcPos int) ***REMOVED***
	if len(p.srcMap) > 0 && p.srcMap[len(p.srcMap)-1].srcPos == srcPos ***REMOVED***
		return
	***REMOVED***
	p.srcMap = append(p.srcMap, srcMapItem***REMOVED***pc: len(p.code), srcPos: srcPos***REMOVED***)
***REMOVED***

func (s *scope) lookupName(name unistring.String) (binding *binding, noDynamics bool) ***REMOVED***
	noDynamics = true
	toStash := false
	for curScope := s; ; curScope = curScope.outer ***REMOVED***
		if curScope.outer != nil ***REMOVED***
			if b, exists := curScope.boundNames[name]; exists ***REMOVED***
				if toStash && !b.inStash ***REMOVED***
					b.moveToStash()
				***REMOVED***
				binding = b
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			noDynamics = false
			return
		***REMOVED***
		if curScope.dynamic ***REMOVED***
			noDynamics = false
		***REMOVED***
		if name == "arguments" && curScope.funcType != funcNone && curScope.funcType != funcArrow ***REMOVED***
			if curScope.funcType == funcClsInit ***REMOVED***
				s.c.throwSyntaxError(0, "'arguments' is not allowed in class field initializer or static initialization block")
			***REMOVED***
			curScope.argsNeeded = true
			binding, _ = curScope.bindName(name)
			return
		***REMOVED***
		if curScope.isFunction() ***REMOVED***
			toStash = true
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *scope) lookupThis() (*binding, bool) ***REMOVED***
	toStash := false
	for curScope := s; curScope != nil; curScope = curScope.outer ***REMOVED***
		if curScope.outer == nil ***REMOVED***
			if curScope.eval ***REMOVED***
				return nil, true
			***REMOVED***
		***REMOVED***
		if b, exists := curScope.boundNames[thisBindingName]; exists ***REMOVED***
			if toStash && !b.inStash ***REMOVED***
				b.moveToStash()
			***REMOVED***
			return b, false
		***REMOVED***
		if curScope.isFunction() ***REMOVED***
			toStash = true
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func (s *scope) ensureBoundNamesCreated() ***REMOVED***
	if s.boundNames == nil ***REMOVED***
		s.boundNames = make(map[unistring.String]*binding)
	***REMOVED***
***REMOVED***

func (s *scope) addBinding(offset int) *binding ***REMOVED***
	if len(s.bindings) >= (1<<24)-1 ***REMOVED***
		s.c.throwSyntaxError(offset, "Too many variables")
	***REMOVED***
	b := &binding***REMOVED***
		scope: s,
	***REMOVED***
	s.bindings = append(s.bindings, b)
	return b
***REMOVED***

func (s *scope) bindNameLexical(name unistring.String, unique bool, offset int) (*binding, bool) ***REMOVED***
	if b := s.boundNames[name]; b != nil ***REMOVED***
		if unique ***REMOVED***
			s.c.throwSyntaxError(offset, "Identifier '%s' has already been declared", name)
		***REMOVED***
		return b, false
	***REMOVED***
	b := s.addBinding(offset)
	b.name = name
	s.ensureBoundNamesCreated()
	s.boundNames[name] = b
	return b, true
***REMOVED***

func (s *scope) createThisBinding() *binding ***REMOVED***
	thisBinding, _ := s.bindNameLexical(thisBindingName, false, 0)
	thisBinding.isVar = true // don't check on load
	return thisBinding
***REMOVED***

func (s *scope) bindName(name unistring.String) (*binding, bool) ***REMOVED***
	if !s.isFunction() && !s.variable && s.outer != nil ***REMOVED***
		return s.outer.bindName(name)
	***REMOVED***
	b, created := s.bindNameLexical(name, false, 0)
	if created ***REMOVED***
		b.isVar = true
	***REMOVED***
	return b, created
***REMOVED***

func (s *scope) bindNameShadow(name unistring.String) (*binding, bool) ***REMOVED***
	if !s.isFunction() && s.outer != nil ***REMOVED***
		return s.outer.bindNameShadow(name)
	***REMOVED***

	_, exists := s.boundNames[name]
	b := &binding***REMOVED***
		scope: s,
		name:  name,
	***REMOVED***
	s.bindings = append(s.bindings, b)
	s.ensureBoundNamesCreated()
	s.boundNames[name] = b
	return b, !exists
***REMOVED***

func (s *scope) nearestFunction() *scope ***REMOVED***
	for sc := s; sc != nil; sc = sc.outer ***REMOVED***
		if sc.isFunction() ***REMOVED***
			return sc
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (s *scope) nearestThis() *scope ***REMOVED***
	for sc := s; sc != nil; sc = sc.outer ***REMOVED***
		if sc.eval || sc.isFunction() && sc.funcType != funcArrow ***REMOVED***
			return sc
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (s *scope) finaliseVarAlloc(stackOffset int) (stashSize, stackSize int) ***REMOVED***
	argsInStash := false
	if f := s.nearestFunction(); f != nil ***REMOVED***
		argsInStash = f.argsInStash
	***REMOVED***
	stackIdx, stashIdx := 0, 0
	allInStash := s.isDynamic()
	var derivedCtor bool
	if fs := s.nearestThis(); fs != nil && fs.funcType == funcDerivedCtor ***REMOVED***
		derivedCtor = true
	***REMOVED***
	for i, b := range s.bindings ***REMOVED***
		var this bool
		if b.name == thisBindingName ***REMOVED***
			this = true
		***REMOVED***
		if allInStash || b.inStash ***REMOVED***
			for scope, aps := range b.accessPoints ***REMOVED***
				var level uint32
				for sc := scope; sc != nil && sc != s; sc = sc.outer ***REMOVED***
					if sc.needStash || sc.isDynamic() ***REMOVED***
						level++
					***REMOVED***
				***REMOVED***
				if level > 255 ***REMOVED***
					s.c.throwSyntaxError(0, "Maximum nesting level (256) exceeded")
				***REMOVED***
				idx := (level << 24) | uint32(stashIdx)
				base := scope.base
				code := scope.prg.code
				if this ***REMOVED***
					if derivedCtor ***REMOVED***
						for _, pc := range *aps ***REMOVED***
							ap := &code[base+pc]
							switch (*ap).(type) ***REMOVED***
							case loadStack:
								*ap = loadThisStash(idx)
							case initStack:
								*ap = initStash(idx)
							case resolveThisStack:
								*ap = resolveThisStash(idx)
							case _ret:
								*ap = cret(idx)
							default:
								s.c.assert(false, s.c.p.sourceOffset(pc), "Unsupported instruction for 'this'")
							***REMOVED***
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						for _, pc := range *aps ***REMOVED***
							ap := &code[base+pc]
							switch (*ap).(type) ***REMOVED***
							case loadStack:
								*ap = loadStash(idx)
							case initStack:
								*ap = initStash(idx)
							default:
								s.c.assert(false, s.c.p.sourceOffset(pc), "Unsupported instruction for 'this'")
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					for _, pc := range *aps ***REMOVED***
						ap := &code[base+pc]
						switch i := (*ap).(type) ***REMOVED***
						case loadStack:
							*ap = loadStash(idx)
						case storeStack:
							*ap = storeStash(idx)
						case storeStackP:
							*ap = storeStashP(idx)
						case loadStackLex:
							*ap = loadStashLex(idx)
						case storeStackLex:
							*ap = storeStashLex(idx)
						case storeStackLexP:
							*ap = storeStashLexP(idx)
						case initStackP:
							*ap = initStashP(idx)
						case initStack:
							*ap = initStash(idx)
						case *loadMixed:
							i.idx = idx
						case *loadMixedLex:
							i.idx = idx
						case *resolveMixed:
							i.idx = idx
						default:
							s.c.assert(false, s.c.p.sourceOffset(pc), "Unsupported instruction for binding: %T", i)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
			stashIdx++
		***REMOVED*** else ***REMOVED***
			var idx int
			if !this ***REMOVED***
				if i < s.numArgs ***REMOVED***
					idx = -(i + 1)
				***REMOVED*** else ***REMOVED***
					stackIdx++
					idx = stackIdx + stackOffset
				***REMOVED***
			***REMOVED***
			for scope, aps := range b.accessPoints ***REMOVED***
				var level int
				for sc := scope; sc != nil && sc != s; sc = sc.outer ***REMOVED***
					if sc.needStash || sc.isDynamic() ***REMOVED***
						level++
					***REMOVED***
				***REMOVED***
				if level > 255 ***REMOVED***
					s.c.throwSyntaxError(0, "Maximum nesting level (256) exceeded")
				***REMOVED***
				code := scope.prg.code
				base := scope.base
				if this ***REMOVED***
					if derivedCtor ***REMOVED***
						for _, pc := range *aps ***REMOVED***
							ap := &code[base+pc]
							switch (*ap).(type) ***REMOVED***
							case loadStack:
								*ap = loadThisStack***REMOVED******REMOVED***
							case initStack:
								// no-op
							case resolveThisStack:
								// no-op
							case _ret:
								// no-op, already in the right place
							default:
								s.c.assert(false, s.c.p.sourceOffset(pc), "Unsupported instruction for 'this'")
							***REMOVED***
						***REMOVED***
					***REMOVED*** /*else ***REMOVED***
						no-op
					***REMOVED****/
				***REMOVED*** else if argsInStash ***REMOVED***
					for _, pc := range *aps ***REMOVED***
						ap := &code[base+pc]
						switch i := (*ap).(type) ***REMOVED***
						case loadStack:
							*ap = loadStack1(idx)
						case storeStack:
							*ap = storeStack1(idx)
						case storeStackP:
							*ap = storeStack1P(idx)
						case loadStackLex:
							*ap = loadStack1Lex(idx)
						case storeStackLex:
							*ap = storeStack1Lex(idx)
						case storeStackLexP:
							*ap = storeStack1LexP(idx)
						case initStackP:
							*ap = initStack1P(idx)
						case initStack:
							*ap = initStack1(idx)
						case *loadMixed:
							*ap = &loadMixedStack1***REMOVED***name: i.name, idx: idx, level: uint8(level), callee: i.callee***REMOVED***
						case *loadMixedLex:
							*ap = &loadMixedStack1Lex***REMOVED***name: i.name, idx: idx, level: uint8(level), callee: i.callee***REMOVED***
						case *resolveMixed:
							*ap = &resolveMixedStack1***REMOVED***typ: i.typ, name: i.name, idx: idx, level: uint8(level), strict: i.strict***REMOVED***
						default:
							s.c.assert(false, s.c.p.sourceOffset(pc), "Unsupported instruction for binding: %T", i)
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					for _, pc := range *aps ***REMOVED***
						ap := &code[base+pc]
						switch i := (*ap).(type) ***REMOVED***
						case loadStack:
							*ap = loadStack(idx)
						case storeStack:
							*ap = storeStack(idx)
						case storeStackP:
							*ap = storeStackP(idx)
						case loadStackLex:
							*ap = loadStackLex(idx)
						case storeStackLex:
							*ap = storeStackLex(idx)
						case storeStackLexP:
							*ap = storeStackLexP(idx)
						case initStack:
							*ap = initStack(idx)
						case initStackP:
							*ap = initStackP(idx)
						case *loadMixed:
							*ap = &loadMixedStack***REMOVED***name: i.name, idx: idx, level: uint8(level), callee: i.callee***REMOVED***
						case *loadMixedLex:
							*ap = &loadMixedStackLex***REMOVED***name: i.name, idx: idx, level: uint8(level), callee: i.callee***REMOVED***
						case *resolveMixed:
							*ap = &resolveMixedStack***REMOVED***typ: i.typ, name: i.name, idx: idx, level: uint8(level), strict: i.strict***REMOVED***
						default:
							s.c.assert(false, s.c.p.sourceOffset(pc), "Unsupported instruction for binding: %T", i)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, nested := range s.nested ***REMOVED***
		nested.finaliseVarAlloc(stackIdx + stackOffset)
	***REMOVED***
	return stashIdx, stackIdx
***REMOVED***

func (s *scope) moveArgsToStash() ***REMOVED***
	for _, b := range s.bindings ***REMOVED***
		if !b.isArg ***REMOVED***
			break
		***REMOVED***
		b.inStash = true
	***REMOVED***
	s.argsInStash = true
	s.needStash = true
***REMOVED***

func (s *scope) trimCode(delta int) ***REMOVED***
	s.c.p.code = s.c.p.code[delta:]
	srcMap := s.c.p.srcMap
	for i := range srcMap ***REMOVED***
		srcMap[i].pc -= delta
	***REMOVED***
	s.adjustBase(-delta)
***REMOVED***

func (s *scope) adjustBase(delta int) ***REMOVED***
	s.base += delta
	for _, nested := range s.nested ***REMOVED***
		nested.adjustBase(delta)
	***REMOVED***
***REMOVED***

func (s *scope) makeNamesMap() map[unistring.String]uint32 ***REMOVED***
	l := len(s.bindings)
	if l == 0 ***REMOVED***
		return nil
	***REMOVED***
	names := make(map[unistring.String]uint32, l)
	for i, b := range s.bindings ***REMOVED***
		idx := uint32(i)
		if b.isConst ***REMOVED***
			idx |= maskConst
			if b.isStrict ***REMOVED***
				idx |= maskStrict
			***REMOVED***
		***REMOVED***
		if b.isVar ***REMOVED***
			idx |= maskVar
		***REMOVED***
		names[b.name] = idx
	***REMOVED***
	return names
***REMOVED***

func (s *scope) isDynamic() bool ***REMOVED***
	return s.dynLookup || s.dynamic
***REMOVED***

func (s *scope) isFunction() bool ***REMOVED***
	return s.funcType != funcNone && !s.eval
***REMOVED***

func (s *scope) deleteBinding(b *binding) ***REMOVED***
	idx := 0
	for i, bb := range s.bindings ***REMOVED***
		if bb == b ***REMOVED***
			idx = i
			goto found
		***REMOVED***
	***REMOVED***
	return
found:
	delete(s.boundNames, b.name)
	copy(s.bindings[idx:], s.bindings[idx+1:])
	l := len(s.bindings) - 1
	s.bindings[l] = nil
	s.bindings = s.bindings[:l]
***REMOVED***

func (c *compiler) compile(in *ast.Program, strict, inGlobal bool, evalVm *vm) ***REMOVED***
	c.ctxVM = evalVm

	eval := evalVm != nil
	c.p.src = in.File
	c.newScope()
	scope := c.scope
	scope.dynamic = true
	scope.eval = eval
	if !strict && len(in.Body) > 0 ***REMOVED***
		strict = c.isStrict(in.Body) != nil
	***REMOVED***
	scope.strict = strict
	ownVarScope := eval && strict
	ownLexScope := !inGlobal || eval
	if ownVarScope ***REMOVED***
		c.newBlockScope()
		scope = c.scope
		scope.variable = true
	***REMOVED***
	if eval && !inGlobal ***REMOVED***
		for s := evalVm.stash; s != nil; s = s.outer ***REMOVED***
			if ft := s.funcType; ft != funcNone && ft != funcArrow ***REMOVED***
				scope.funcType = ft
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	funcs := c.extractFunctions(in.Body)
	c.createFunctionBindings(funcs)
	numFuncs := len(scope.bindings)
	if inGlobal && !ownVarScope ***REMOVED***
		if numFuncs == len(funcs) ***REMOVED***
			c.compileFunctionsGlobalAllUnique(funcs)
		***REMOVED*** else ***REMOVED***
			c.compileFunctionsGlobal(funcs)
		***REMOVED***
	***REMOVED***
	c.compileDeclList(in.DeclarationList, false)
	numVars := len(scope.bindings) - numFuncs
	vars := make([]unistring.String, len(scope.bindings))
	for i, b := range scope.bindings ***REMOVED***
		vars[i] = b.name
	***REMOVED***
	if len(vars) > 0 && !ownVarScope && ownLexScope ***REMOVED***
		if inGlobal ***REMOVED***
			c.emit(&bindGlobal***REMOVED***
				vars:      vars[numFuncs:],
				funcs:     vars[:numFuncs],
				deletable: eval,
			***REMOVED***)
		***REMOVED*** else ***REMOVED***
			c.emit(&bindVars***REMOVED***names: vars, deletable: eval***REMOVED***)
		***REMOVED***
	***REMOVED***
	var enter *enterBlock
	if c.compileLexicalDeclarations(in.Body, ownVarScope || !ownLexScope) ***REMOVED***
		if ownLexScope ***REMOVED***
			c.block = &block***REMOVED***
				outer:      c.block,
				typ:        blockScope,
				needResult: true,
			***REMOVED***
			enter = &enterBlock***REMOVED******REMOVED***
			c.emit(enter)
		***REMOVED***
	***REMOVED***
	if len(scope.bindings) > 0 && !ownLexScope ***REMOVED***
		var lets, consts []unistring.String
		for _, b := range c.scope.bindings[numFuncs+numVars:] ***REMOVED***
			if b.isConst ***REMOVED***
				consts = append(consts, b.name)
			***REMOVED*** else ***REMOVED***
				lets = append(lets, b.name)
			***REMOVED***
		***REMOVED***
		c.emit(&bindGlobal***REMOVED***
			vars:   vars[numFuncs:],
			funcs:  vars[:numFuncs],
			lets:   lets,
			consts: consts,
		***REMOVED***)
	***REMOVED***
	if !inGlobal || ownVarScope ***REMOVED***
		c.compileFunctions(funcs)
	***REMOVED***
	c.compileStatements(in.Body, true)
	if enter != nil ***REMOVED***
		c.leaveScopeBlock(enter)
		c.popScope()
	***REMOVED***

	c.p.code = append(c.p.code, halt)

	scope.finaliseVarAlloc(0)
***REMOVED***

func (c *compiler) compileDeclList(v []*ast.VariableDeclaration, inFunc bool) ***REMOVED***
	for _, value := range v ***REMOVED***
		c.createVarBindings(value, inFunc)
	***REMOVED***
***REMOVED***

func (c *compiler) extractLabelled(st ast.Statement) ast.Statement ***REMOVED***
	if st, ok := st.(*ast.LabelledStatement); ok ***REMOVED***
		return c.extractLabelled(st.Statement)
	***REMOVED***
	return st
***REMOVED***

func (c *compiler) extractFunctions(list []ast.Statement) (funcs []*ast.FunctionDeclaration) ***REMOVED***
	for _, st := range list ***REMOVED***
		var decl *ast.FunctionDeclaration
		switch st := c.extractLabelled(st).(type) ***REMOVED***
		case *ast.FunctionDeclaration:
			decl = st
		case *ast.LabelledStatement:
			if st1, ok := st.Statement.(*ast.FunctionDeclaration); ok ***REMOVED***
				decl = st1
			***REMOVED*** else ***REMOVED***
				continue
			***REMOVED***
		default:
			continue
		***REMOVED***
		funcs = append(funcs, decl)
	***REMOVED***
	return
***REMOVED***

func (c *compiler) createFunctionBindings(funcs []*ast.FunctionDeclaration) ***REMOVED***
	s := c.scope
	if s.outer != nil ***REMOVED***
		unique := !s.isFunction() && !s.variable && s.strict
		for _, decl := range funcs ***REMOVED***
			s.bindNameLexical(decl.Function.Name.Name, unique, int(decl.Function.Name.Idx1())-1)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, decl := range funcs ***REMOVED***
			s.bindName(decl.Function.Name.Name)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunctions(list []*ast.FunctionDeclaration) ***REMOVED***
	for _, decl := range list ***REMOVED***
		c.compileFunction(decl)
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunctionsGlobalAllUnique(list []*ast.FunctionDeclaration) ***REMOVED***
	for _, decl := range list ***REMOVED***
		c.compileFunctionLiteral(decl.Function, false).emitGetter(true)
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunctionsGlobal(list []*ast.FunctionDeclaration) ***REMOVED***
	m := make(map[unistring.String]int, len(list))
	for i := len(list) - 1; i >= 0; i-- ***REMOVED***
		name := list[i].Function.Name.Name
		if _, exists := m[name]; !exists ***REMOVED***
			m[name] = i
		***REMOVED***
	***REMOVED***
	idx := 0
	for i, decl := range list ***REMOVED***
		name := decl.Function.Name.Name
		if m[name] == i ***REMOVED***
			c.compileFunctionLiteral(decl.Function, false).emitGetter(true)
			c.scope.bindings[idx] = c.scope.boundNames[name]
			idx++
		***REMOVED*** else ***REMOVED***
			leave := c.enterDummyMode()
			c.compileFunctionLiteral(decl.Function, false).emitGetter(false)
			leave()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) createVarIdBinding(name unistring.String, offset int, inFunc bool) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierLName(name, offset)
		c.checkIdentifierName(name, offset)
	***REMOVED***
	if !inFunc || name != "arguments" ***REMOVED***
		c.scope.bindName(name)
	***REMOVED***
***REMOVED***

func (c *compiler) createBindings(target ast.Expression, createIdBinding func(name unistring.String, offset int)) ***REMOVED***
	switch target := target.(type) ***REMOVED***
	case *ast.Identifier:
		createIdBinding(target.Name, int(target.Idx)-1)
	case *ast.ObjectPattern:
		for _, prop := range target.Properties ***REMOVED***
			switch prop := prop.(type) ***REMOVED***
			case *ast.PropertyShort:
				createIdBinding(prop.Name.Name, int(prop.Name.Idx)-1)
			case *ast.PropertyKeyed:
				c.createBindings(prop.Value, createIdBinding)
			default:
				c.throwSyntaxError(int(target.Idx0()-1), "unsupported property type in ObjectPattern: %T", prop)
			***REMOVED***
		***REMOVED***
		if target.Rest != nil ***REMOVED***
			c.createBindings(target.Rest, createIdBinding)
		***REMOVED***
	case *ast.ArrayPattern:
		for _, elt := range target.Elements ***REMOVED***
			if elt != nil ***REMOVED***
				c.createBindings(elt, createIdBinding)
			***REMOVED***
		***REMOVED***
		if target.Rest != nil ***REMOVED***
			c.createBindings(target.Rest, createIdBinding)
		***REMOVED***
	case *ast.AssignExpression:
		c.createBindings(target.Left, createIdBinding)
	default:
		c.throwSyntaxError(int(target.Idx0()-1), "unsupported binding target: %T", target)
	***REMOVED***
***REMOVED***

func (c *compiler) createVarBinding(target ast.Expression, inFunc bool) ***REMOVED***
	c.createBindings(target, func(name unistring.String, offset int) ***REMOVED***
		c.createVarIdBinding(name, offset, inFunc)
	***REMOVED***)
***REMOVED***

func (c *compiler) createVarBindings(v *ast.VariableDeclaration, inFunc bool) ***REMOVED***
	for _, item := range v.List ***REMOVED***
		c.createVarBinding(item.Target, inFunc)
	***REMOVED***
***REMOVED***

func (c *compiler) createLexicalIdBinding(name unistring.String, isConst bool, offset int) *binding ***REMOVED***
	if name == "let" ***REMOVED***
		c.throwSyntaxError(offset, "let is disallowed as a lexically bound name")
	***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierLName(name, offset)
		c.checkIdentifierName(name, offset)
	***REMOVED***
	b, _ := c.scope.bindNameLexical(name, true, offset)
	if isConst ***REMOVED***
		b.isConst, b.isStrict = true, true
	***REMOVED***
	return b
***REMOVED***

func (c *compiler) createLexicalIdBindingFuncBody(name unistring.String, isConst bool, offset int, calleeBinding *binding) *binding ***REMOVED***
	if name == "let" ***REMOVED***
		c.throwSyntaxError(offset, "let is disallowed as a lexically bound name")
	***REMOVED***
	if c.scope.strict ***REMOVED***
		c.checkIdentifierLName(name, offset)
		c.checkIdentifierName(name, offset)
	***REMOVED***
	paramScope := c.scope.outer
	parentBinding := paramScope.boundNames[name]
	if parentBinding != nil ***REMOVED***
		if parentBinding != calleeBinding && (name != "arguments" || !paramScope.argsNeeded) ***REMOVED***
			c.throwSyntaxError(offset, "Identifier '%s' has already been declared", name)
		***REMOVED***
	***REMOVED***
	b, _ := c.scope.bindNameLexical(name, true, offset)
	if isConst ***REMOVED***
		b.isConst, b.isStrict = true, true
	***REMOVED***
	return b
***REMOVED***

func (c *compiler) createLexicalBinding(target ast.Expression, isConst bool) ***REMOVED***
	c.createBindings(target, func(name unistring.String, offset int) ***REMOVED***
		c.createLexicalIdBinding(name, isConst, offset)
	***REMOVED***)
***REMOVED***

func (c *compiler) createLexicalBindings(lex *ast.LexicalDeclaration) ***REMOVED***
	for _, d := range lex.List ***REMOVED***
		c.createLexicalBinding(d.Target, lex.Token == token.CONST)
	***REMOVED***
***REMOVED***

func (c *compiler) compileLexicalDeclarations(list []ast.Statement, scopeDeclared bool) bool ***REMOVED***
	for _, st := range list ***REMOVED***
		if lex, ok := st.(*ast.LexicalDeclaration); ok ***REMOVED***
			if !scopeDeclared ***REMOVED***
				c.newBlockScope()
				scopeDeclared = true
			***REMOVED***
			c.createLexicalBindings(lex)
		***REMOVED*** else if cls, ok := st.(*ast.ClassDeclaration); ok ***REMOVED***
			if !scopeDeclared ***REMOVED***
				c.newBlockScope()
				scopeDeclared = true
			***REMOVED***
			c.createLexicalIdBinding(cls.Class.Name.Name, false, int(cls.Class.Name.Idx)-1)
		***REMOVED***
	***REMOVED***
	return scopeDeclared
***REMOVED***

func (c *compiler) compileLexicalDeclarationsFuncBody(list []ast.Statement, calleeBinding *binding) ***REMOVED***
	for _, st := range list ***REMOVED***
		if lex, ok := st.(*ast.LexicalDeclaration); ok ***REMOVED***
			isConst := lex.Token == token.CONST
			for _, d := range lex.List ***REMOVED***
				c.createBindings(d.Target, func(name unistring.String, offset int) ***REMOVED***
					c.createLexicalIdBindingFuncBody(name, isConst, offset, calleeBinding)
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *compiler) compileFunction(v *ast.FunctionDeclaration) ***REMOVED***
	name := v.Function.Name.Name
	b := c.scope.boundNames[name]
	if b == nil || b.isVar ***REMOVED***
		e := &compiledIdentifierExpr***REMOVED***
			name: v.Function.Name.Name,
		***REMOVED***
		e.init(c, v.Function.Idx0())
		e.emitSetter(c.compileFunctionLiteral(v.Function, false), false)
	***REMOVED*** else ***REMOVED***
		c.compileFunctionLiteral(v.Function, false).emitGetter(true)
		b.emitInitP()
	***REMOVED***
***REMOVED***

func (c *compiler) compileStandaloneFunctionDecl(v *ast.FunctionDeclaration) ***REMOVED***
	if c.scope.strict ***REMOVED***
		c.throwSyntaxError(int(v.Idx0())-1, "In strict mode code, functions can only be declared at top level or inside a block.")
	***REMOVED***
	c.throwSyntaxError(int(v.Idx0())-1, "In non-strict mode code, functions can only be declared at top level, inside a block, or as the body of an if statement.")
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

func (c *compiler) isStrict(list []ast.Statement) *ast.StringLiteral ***REMOVED***
	for _, st := range list ***REMOVED***
		if st, ok := st.(*ast.ExpressionStatement); ok ***REMOVED***
			if e, ok := st.Expression.(*ast.StringLiteral); ok ***REMOVED***
				if e.Literal == `"use strict"` || e.Literal == `'use strict'` ***REMOVED***
					return e
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *compiler) isStrictStatement(s ast.Statement) *ast.StringLiteral ***REMOVED***
	if s, ok := s.(*ast.BlockStatement); ok ***REMOVED***
		return c.isStrict(s.List)
	***REMOVED***
	return nil
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

// Enter a 'dummy' compilation mode. Any code produced after this method is called will be discarded after
// leaveFunc is called with no additional side effects. This is useful for compiling code inside a
// constant falsy condition 'if' branch or a loop (i.e 'if (false) ***REMOVED*** ... ***REMOVED*** or while (false) ***REMOVED*** ... ***REMOVED***).
// Such code should not be included in the final compilation result as it's never called, but it must
// still produce compilation errors if there are any.
// TODO: make sure variable lookups do not de-optimise parent scopes
func (c *compiler) enterDummyMode() (leaveFunc func()) ***REMOVED***
	savedBlock, savedProgram := c.block, c.p
	if savedBlock != nil ***REMOVED***
		c.block = &block***REMOVED***
			typ:      savedBlock.typ,
			label:    savedBlock.label,
			outer:    savedBlock.outer,
			breaking: savedBlock.breaking,
		***REMOVED***
	***REMOVED***
	c.p = &Program***REMOVED******REMOVED***
	c.newScope()
	return func() ***REMOVED***
		c.block, c.p = savedBlock, savedProgram
		c.popScope()
	***REMOVED***
***REMOVED***

func (c *compiler) compileStatementDummy(statement ast.Statement) ***REMOVED***
	leave := c.enterDummyMode()
	c.compileStatement(statement, false)
	leave()
***REMOVED***

func (c *compiler) assert(cond bool, offset int, msg string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if !cond ***REMOVED***
		c.throwSyntaxError(offset, "Compiler bug: "+msg, args...)
	***REMOVED***
***REMOVED***

func privateIdString(desc unistring.String) unistring.String ***REMOVED***
	return asciiString("#").concat(stringValueFromRaw(desc)).string()
***REMOVED***

type privateName struct ***REMOVED***
	idx                  int
	isStatic             bool
	isMethod             bool
	hasGetter, hasSetter bool
***REMOVED***

type resolvedPrivateName struct ***REMOVED***
	name     unistring.String
	idx      uint32
	level    uint8
	isStatic bool
	isMethod bool
***REMOVED***

func (r *resolvedPrivateName) string() unistring.String ***REMOVED***
	return privateIdString(r.name)
***REMOVED***

type privateEnvRegistry struct ***REMOVED***
	fields, methods []unistring.String
***REMOVED***

type classScope struct ***REMOVED***
	c            *compiler
	privateNames map[unistring.String]*privateName

	instanceEnv, staticEnv privateEnvRegistry

	outer *classScope
***REMOVED***

func (r *privateEnvRegistry) createPrivateMethodId(name unistring.String) int ***REMOVED***
	r.methods = append(r.methods, name)
	return len(r.methods) - 1
***REMOVED***

func (r *privateEnvRegistry) createPrivateFieldId(name unistring.String) int ***REMOVED***
	r.fields = append(r.fields, name)
	return len(r.fields) - 1
***REMOVED***

func (s *classScope) declarePrivateId(name unistring.String, kind ast.PropertyKind, isStatic bool, offset int) ***REMOVED***
	pn := s.privateNames[name]
	if pn != nil ***REMOVED***
		if pn.isStatic == isStatic ***REMOVED***
			switch kind ***REMOVED***
			case ast.PropertyKindGet:
				if pn.hasSetter && !pn.hasGetter ***REMOVED***
					pn.hasGetter = true
					return
				***REMOVED***
			case ast.PropertyKindSet:
				if pn.hasGetter && !pn.hasSetter ***REMOVED***
					pn.hasSetter = true
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
		s.c.throwSyntaxError(offset, "Identifier '#%s' has already been declared", name)
		panic("unreachable")
	***REMOVED***
	var env *privateEnvRegistry
	if isStatic ***REMOVED***
		env = &s.staticEnv
	***REMOVED*** else ***REMOVED***
		env = &s.instanceEnv
	***REMOVED***

	pn = &privateName***REMOVED***
		isStatic:  isStatic,
		hasGetter: kind == ast.PropertyKindGet,
		hasSetter: kind == ast.PropertyKindSet,
	***REMOVED***
	if kind != ast.PropertyKindValue ***REMOVED***
		pn.idx = env.createPrivateMethodId(name)
		pn.isMethod = true
	***REMOVED*** else ***REMOVED***
		pn.idx = env.createPrivateFieldId(name)
	***REMOVED***

	if s.privateNames == nil ***REMOVED***
		s.privateNames = make(map[unistring.String]*privateName)
	***REMOVED***
	s.privateNames[name] = pn
***REMOVED***

func (s *classScope) getDeclaredPrivateId(name unistring.String) *privateName ***REMOVED***
	if n := s.privateNames[name]; n != nil ***REMOVED***
		return n
	***REMOVED***
	s.c.assert(false, 0, "getDeclaredPrivateId() for undeclared id")
	panic("unreachable")
***REMOVED***

func (c *compiler) resolvePrivateName(name unistring.String, offset int) (*resolvedPrivateName, *privateId) ***REMOVED***
	level := 0
	for s := c.classScope; s != nil; s = s.outer ***REMOVED***
		if len(s.privateNames) > 0 ***REMOVED***
			if pn := s.privateNames[name]; pn != nil ***REMOVED***
				return &resolvedPrivateName***REMOVED***
					name:     name,
					idx:      uint32(pn.idx),
					level:    uint8(level),
					isStatic: pn.isStatic,
					isMethod: pn.isMethod,
				***REMOVED***, nil
			***REMOVED***
			level++
		***REMOVED***
	***REMOVED***
	if c.ctxVM != nil ***REMOVED***
		for s := c.ctxVM.privEnv; s != nil; s = s.outer ***REMOVED***
			if id := s.names[name]; id != nil ***REMOVED***
				return nil, id
			***REMOVED***
		***REMOVED***
	***REMOVED***
	c.throwSyntaxError(offset, "Private field '#%s' must be declared in an enclosing class", name)
	panic("unreachable")
***REMOVED***
