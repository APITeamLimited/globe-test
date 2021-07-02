package goja

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/dop251/goja/unistring"
)

const (
	maxInt = 1 << 53
)

type valueStack []Value

type stash struct ***REMOVED***
	values    []Value
	extraArgs []Value
	names     map[unistring.String]uint32
	obj       *Object

	outer *stash

	function bool
***REMOVED***

type context struct ***REMOVED***
	prg       *Program
	funcName  unistring.String
	stash     *stash
	newTarget Value
	result    Value
	pc, sb    int
	args      int
***REMOVED***

type iterStackItem struct ***REMOVED***
	val  Value
	f    iterNextFunc
	iter *Object
***REMOVED***

type ref interface ***REMOVED***
	get() Value
	set(Value)
	refname() unistring.String
***REMOVED***

type stashRef struct ***REMOVED***
	n   unistring.String
	v   *[]Value
	idx int
***REMOVED***

func (r *stashRef) get() Value ***REMOVED***
	return nilSafe((*r.v)[r.idx])
***REMOVED***

func (r *stashRef) set(v Value) ***REMOVED***
	(*r.v)[r.idx] = v
***REMOVED***

func (r *stashRef) refname() unistring.String ***REMOVED***
	return r.n
***REMOVED***

type stashRefLex struct ***REMOVED***
	stashRef
***REMOVED***

func (r *stashRefLex) get() Value ***REMOVED***
	v := (*r.v)[r.idx]
	if v == nil ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	return v
***REMOVED***

func (r *stashRefLex) set(v Value) ***REMOVED***
	p := &(*r.v)[r.idx]
	if *p == nil ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	*p = v
***REMOVED***

type stashRefConst struct ***REMOVED***
	stashRefLex
***REMOVED***

func (r *stashRefConst) set(v Value) ***REMOVED***
	panic(errAssignToConst)
***REMOVED***

type objRef struct ***REMOVED***
	base   objectImpl
	name   unistring.String
	strict bool
***REMOVED***

func (r *objRef) get() Value ***REMOVED***
	return r.base.getStr(r.name, nil)
***REMOVED***

func (r *objRef) set(v Value) ***REMOVED***
	r.base.setOwnStr(r.name, v, r.strict)
***REMOVED***

func (r *objRef) refname() unistring.String ***REMOVED***
	return r.name
***REMOVED***

type unresolvedRef struct ***REMOVED***
	runtime *Runtime
	name    unistring.String
***REMOVED***

func (r *unresolvedRef) get() Value ***REMOVED***
	r.runtime.throwReferenceError(r.name)
	panic("Unreachable")
***REMOVED***

func (r *unresolvedRef) set(Value) ***REMOVED***
	r.get()
***REMOVED***

func (r *unresolvedRef) refname() unistring.String ***REMOVED***
	return r.name
***REMOVED***

type vm struct ***REMOVED***
	r            *Runtime
	prg          *Program
	funcName     unistring.String
	pc           int
	stack        valueStack
	sp, sb, args int

	stash     *stash
	callStack []context
	iterStack []iterStackItem
	refStack  []ref
	newTarget Value
	result    Value

	maxCallStackSize int

	stashAllocs int
	halt        bool

	interrupted   uint32
	interruptVal  interface***REMOVED******REMOVED***
	interruptLock sync.Mutex
***REMOVED***

type instruction interface ***REMOVED***
	exec(*vm)
***REMOVED***

func intToValue(i int64) Value ***REMOVED***
	if i >= -maxInt && i <= maxInt ***REMOVED***
		if i >= -128 && i <= 127 ***REMOVED***
			return intCache[i+128]
		***REMOVED***
		return valueInt(i)
	***REMOVED***
	return valueFloat(i)
***REMOVED***

func floatToInt(f float64) (result int64, ok bool) ***REMOVED***
	if (f != 0 || !math.Signbit(f)) && !math.IsInf(f, 0) && f == math.Trunc(f) && f >= -maxInt && f <= maxInt ***REMOVED***
		return int64(f), true
	***REMOVED***
	return 0, false
***REMOVED***

func floatToValue(f float64) (result Value) ***REMOVED***
	if i, ok := floatToInt(f); ok ***REMOVED***
		return intToValue(i)
	***REMOVED***
	switch ***REMOVED***
	case f == 0:
		return _negativeZero
	case math.IsNaN(f):
		return _NaN
	case math.IsInf(f, 1):
		return _positiveInf
	case math.IsInf(f, -1):
		return _negativeInf
	***REMOVED***
	return valueFloat(f)
***REMOVED***

func assertInt64(v Value) (int64, bool) ***REMOVED***
	num := v.ToNumber()
	if i, ok := num.(valueInt); ok ***REMOVED***
		return int64(i), true
	***REMOVED***
	if f, ok := num.(valueFloat); ok ***REMOVED***
		if i, ok := floatToInt(float64(f)); ok ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

func toIntIgnoreNegZero(v Value) (int64, bool) ***REMOVED***
	num := v.ToNumber()
	if i, ok := num.(valueInt); ok ***REMOVED***
		return int64(i), true
	***REMOVED***
	if f, ok := num.(valueFloat); ok ***REMOVED***
		if v == _negativeZero ***REMOVED***
			return 0, true
		***REMOVED***
		if i, ok := floatToInt(float64(f)); ok ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

func (s *valueStack) expand(idx int) ***REMOVED***
	if idx < len(*s) ***REMOVED***
		return
	***REMOVED***
	idx++
	if idx < cap(*s) ***REMOVED***
		*s = (*s)[:idx]
	***REMOVED*** else ***REMOVED***
		var newCap int
		if idx < 1024 ***REMOVED***
			newCap = idx * 2
		***REMOVED*** else ***REMOVED***
			newCap = (idx + 1025) &^ 1023
		***REMOVED***
		n := make([]Value, idx, newCap)
		copy(n, *s)
		*s = n
	***REMOVED***
***REMOVED***

func stashObjHas(obj *Object, name unistring.String) bool ***REMOVED***
	if obj.self.hasPropertyStr(name) ***REMOVED***
		if unscopables, ok := obj.self.getSym(SymUnscopables, nil).(*Object); ok ***REMOVED***
			if b := unscopables.self.getStr(name, nil); b != nil ***REMOVED***
				return !b.ToBoolean()
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (s *stash) initByIdx(idx uint32, v Value) ***REMOVED***
	if s.obj != nil ***REMOVED***
		panic("Attempt to init by idx into an object scope")
	***REMOVED***
	s.values[idx] = v
***REMOVED***

func (s *stash) initByName(name unistring.String, v Value) ***REMOVED***
	if idx, exists := s.names[name]; exists ***REMOVED***
		s.values[idx&^maskTyp] = v
	***REMOVED*** else ***REMOVED***
		panic(referenceError(fmt.Sprintf("%s is not defined", name)))
	***REMOVED***
***REMOVED***

func (s *stash) getByIdx(idx uint32) Value ***REMOVED***
	return s.values[idx]
***REMOVED***

func (s *stash) getByName(name unistring.String) (v Value, exists bool) ***REMOVED***
	if s.obj != nil ***REMOVED***
		if stashObjHas(s.obj, name) ***REMOVED***
			return nilSafe(s.obj.self.getStr(name, nil)), true
		***REMOVED***
		return nil, false
	***REMOVED***
	if idx, exists := s.names[name]; exists ***REMOVED***
		v := s.values[idx&^maskTyp]
		if v == nil ***REMOVED***
			if idx&maskVar == 0 ***REMOVED***
				panic(errAccessBeforeInit)
			***REMOVED*** else ***REMOVED***
				v = _undefined
			***REMOVED***
		***REMOVED***
		return v, true
	***REMOVED***
	return nil, false
***REMOVED***

func (s *stash) getRefByName(name unistring.String, strict bool) ref ***REMOVED***
	if obj := s.obj; obj != nil ***REMOVED***
		if stashObjHas(obj, name) ***REMOVED***
			return &objRef***REMOVED***
				base:   obj.self,
				name:   name,
				strict: strict,
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if idx, exists := s.names[name]; exists ***REMOVED***
			if idx&maskVar == 0 ***REMOVED***
				if idx&maskConst == 0 ***REMOVED***
					return &stashRefLex***REMOVED***
						stashRef: stashRef***REMOVED***
							n:   name,
							v:   &s.values,
							idx: int(idx &^ maskTyp),
						***REMOVED***,
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					return &stashRefConst***REMOVED***
						stashRefLex: stashRefLex***REMOVED***
							stashRef: stashRef***REMOVED***
								n:   name,
								v:   &s.values,
								idx: int(idx &^ maskTyp),
							***REMOVED***,
						***REMOVED***,
					***REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return &stashRef***REMOVED***
					n:   name,
					v:   &s.values,
					idx: int(idx &^ maskTyp),
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (s *stash) createBinding(name unistring.String, deletable bool) ***REMOVED***
	if s.names == nil ***REMOVED***
		s.names = make(map[unistring.String]uint32)
	***REMOVED***
	if _, exists := s.names[name]; !exists ***REMOVED***
		idx := uint32(len(s.names)) | maskVar
		if deletable ***REMOVED***
			idx |= maskDeletable
		***REMOVED***
		s.names[name] = idx
		s.values = append(s.values, _undefined)
	***REMOVED***
***REMOVED***

func (s *stash) createLexBinding(name unistring.String, isConst bool) ***REMOVED***
	if s.names == nil ***REMOVED***
		s.names = make(map[unistring.String]uint32)
	***REMOVED***
	if _, exists := s.names[name]; !exists ***REMOVED***
		idx := uint32(len(s.names))
		if isConst ***REMOVED***
			idx |= maskConst
		***REMOVED***
		s.names[name] = idx
		s.values = append(s.values, nil)
	***REMOVED***
***REMOVED***

func (s *stash) deleteBinding(name unistring.String) ***REMOVED***
	delete(s.names, name)
***REMOVED***

func (vm *vm) newStash() ***REMOVED***
	vm.stash = &stash***REMOVED***
		outer: vm.stash,
	***REMOVED***
	vm.stashAllocs++
***REMOVED***

func (vm *vm) init() ***REMOVED***
	vm.sb = -1
	vm.stash = &vm.r.global.stash
	vm.maxCallStackSize = math.MaxInt32
***REMOVED***

func (vm *vm) run() ***REMOVED***
	vm.halt = false
	interrupted := false
	ticks := 0
	for !vm.halt ***REMOVED***
		if interrupted = atomic.LoadUint32(&vm.interrupted) != 0; interrupted ***REMOVED***
			break
		***REMOVED***
		vm.prg.code[vm.pc].exec(vm)
		ticks++
		if ticks > 10000 ***REMOVED***
			runtime.Gosched()
			ticks = 0
		***REMOVED***
	***REMOVED***

	if interrupted ***REMOVED***
		vm.interruptLock.Lock()
		v := &InterruptedError***REMOVED***
			iface: vm.interruptVal,
		***REMOVED***
		atomic.StoreUint32(&vm.interrupted, 0)
		vm.interruptVal = nil
		vm.interruptLock.Unlock()
		panic(&uncatchableException***REMOVED***
			stack: &v.stack,
			err:   v,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (vm *vm) Interrupt(v interface***REMOVED******REMOVED***) ***REMOVED***
	vm.interruptLock.Lock()
	vm.interruptVal = v
	atomic.StoreUint32(&vm.interrupted, 1)
	vm.interruptLock.Unlock()
***REMOVED***

func (vm *vm) ClearInterrupt() ***REMOVED***
	atomic.StoreUint32(&vm.interrupted, 0)
***REMOVED***

func (vm *vm) captureStack(stack []StackFrame, ctxOffset int) []StackFrame ***REMOVED***
	// Unroll the context stack
	if vm.pc != -1 ***REMOVED***
		stack = append(stack, StackFrame***REMOVED***prg: vm.prg, pc: vm.pc, funcName: vm.funcName***REMOVED***)
	***REMOVED***
	for i := len(vm.callStack) - 1; i > ctxOffset-1; i-- ***REMOVED***
		if vm.callStack[i].pc != -1 ***REMOVED***
			stack = append(stack, StackFrame***REMOVED***prg: vm.callStack[i].prg, pc: vm.callStack[i].pc - 1, funcName: vm.callStack[i].funcName***REMOVED***)
		***REMOVED***
	***REMOVED***
	return stack
***REMOVED***

func (vm *vm) try(f func()) (ex *Exception) ***REMOVED***
	var ctx context
	vm.saveCtx(&ctx)

	ctxOffset := len(vm.callStack)
	sp := vm.sp
	iterLen := len(vm.iterStack)
	refLen := len(vm.refStack)

	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			defer func() ***REMOVED***
				vm.callStack = vm.callStack[:ctxOffset]
				vm.restoreCtx(&ctx)
				vm.sp = sp

				// Restore other stacks
				iterTail := vm.iterStack[iterLen:]
				for i := range iterTail ***REMOVED***
					if iter := iterTail[i].iter; iter != nil ***REMOVED***
						vm.try(func() ***REMOVED***
							returnIter(iter)
						***REMOVED***)
					***REMOVED***
					iterTail[i] = iterStackItem***REMOVED******REMOVED***
				***REMOVED***
				vm.iterStack = vm.iterStack[:iterLen]
				refTail := vm.refStack[refLen:]
				for i := range refTail ***REMOVED***
					refTail[i] = nil
				***REMOVED***
				vm.refStack = vm.refStack[:refLen]
			***REMOVED***()
			switch x1 := x.(type) ***REMOVED***
			case Value:
				ex = &Exception***REMOVED***
					val: x1,
				***REMOVED***
			case *Exception:
				ex = x1
			case *uncatchableException:
				*x1.stack = vm.captureStack(*x1.stack, ctxOffset)
				panic(x1)
			case typeError:
				ex = &Exception***REMOVED***
					val: vm.r.NewTypeError(string(x1)),
				***REMOVED***
			case referenceError:
				ex = &Exception***REMOVED***
					val: vm.r.newError(vm.r.global.ReferenceError, string(x1)),
				***REMOVED***
			case rangeError:
				ex = &Exception***REMOVED***
					val: vm.r.newError(vm.r.global.RangeError, string(x1)),
				***REMOVED***
			default:
				/*
					if vm.prg != nil ***REMOVED***
						vm.prg.dumpCode(log.Printf)
					***REMOVED***
					log.Print("Stack: ", string(debug.Stack()))
					panic(fmt.Errorf("Panic at %d: %v", vm.pc, x))
				*/
				panic(x)
			***REMOVED***
			ex.stack = vm.captureStack(ex.stack, ctxOffset)
		***REMOVED***
	***REMOVED***()

	f()
	return
***REMOVED***

func (vm *vm) runTry() (ex *Exception) ***REMOVED***
	return vm.try(vm.run)
***REMOVED***

func (vm *vm) push(v Value) ***REMOVED***
	vm.stack.expand(vm.sp)
	vm.stack[vm.sp] = v
	vm.sp++
***REMOVED***

func (vm *vm) pop() Value ***REMOVED***
	vm.sp--
	return vm.stack[vm.sp]
***REMOVED***

func (vm *vm) peek() Value ***REMOVED***
	return vm.stack[vm.sp-1]
***REMOVED***

func (vm *vm) saveCtx(ctx *context) ***REMOVED***
	ctx.prg, ctx.stash, ctx.newTarget, ctx.result, ctx.pc, ctx.sb, ctx.args =
		vm.prg, vm.stash, vm.newTarget, vm.result, vm.pc, vm.sb, vm.args
	if vm.funcName != "" ***REMOVED***
		ctx.funcName = vm.funcName
	***REMOVED*** else if ctx.prg != nil && ctx.prg.funcName != "" ***REMOVED***
		ctx.funcName = ctx.prg.funcName
	***REMOVED***
***REMOVED***

func (vm *vm) pushCtx() ***REMOVED***
	if len(vm.callStack) > vm.maxCallStackSize ***REMOVED***
		ex := &StackOverflowError***REMOVED******REMOVED***
		panic(&uncatchableException***REMOVED***
			stack: &ex.stack,
			err:   ex,
		***REMOVED***)
	***REMOVED***
	vm.callStack = append(vm.callStack, context***REMOVED******REMOVED***)
	ctx := &vm.callStack[len(vm.callStack)-1]
	vm.saveCtx(ctx)
***REMOVED***

func (vm *vm) restoreCtx(ctx *context) ***REMOVED***
	vm.prg, vm.funcName, vm.stash, vm.newTarget, vm.result, vm.pc, vm.sb, vm.args =
		ctx.prg, ctx.funcName, ctx.stash, ctx.newTarget, ctx.result, ctx.pc, ctx.sb, ctx.args
***REMOVED***

func (vm *vm) popCtx() ***REMOVED***
	l := len(vm.callStack) - 1
	ctx := &vm.callStack[l]
	vm.restoreCtx(ctx)

	ctx.prg = nil
	ctx.stash = nil
	ctx.result = nil
	ctx.newTarget = nil

	vm.callStack = vm.callStack[:l]
***REMOVED***

func (vm *vm) toCallee(v Value) *Object ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		return obj
	***REMOVED***
	switch unresolved := v.(type) ***REMOVED***
	case valueUnresolved:
		unresolved.throw()
		panic("Unreachable")
	case memberUnresolved:
		panic(vm.r.NewTypeError("Object has no member '%s'", unresolved.ref))
	***REMOVED***
	panic(vm.r.NewTypeError("Value is not an object: %s", v.toString()))
***REMOVED***

type loadVal uint32

func (l loadVal) exec(vm *vm) ***REMOVED***
	vm.push(vm.prg.values[l])
	vm.pc++
***REMOVED***

type _loadUndef struct***REMOVED******REMOVED***

var loadUndef _loadUndef

func (_loadUndef) exec(vm *vm) ***REMOVED***
	vm.push(_undefined)
	vm.pc++
***REMOVED***

type _loadNil struct***REMOVED******REMOVED***

var loadNil _loadNil

func (_loadNil) exec(vm *vm) ***REMOVED***
	vm.push(nil)
	vm.pc++
***REMOVED***

type _saveResult struct***REMOVED******REMOVED***

var saveResult _saveResult

func (_saveResult) exec(vm *vm) ***REMOVED***
	vm.sp--
	vm.result = vm.stack[vm.sp]
	vm.pc++
***REMOVED***

type _clearResult struct***REMOVED******REMOVED***

var clearResult _clearResult

func (_clearResult) exec(vm *vm) ***REMOVED***
	vm.result = _undefined
	vm.pc++
***REMOVED***

type _loadGlobalObject struct***REMOVED******REMOVED***

var loadGlobalObject _loadGlobalObject

func (_loadGlobalObject) exec(vm *vm) ***REMOVED***
	vm.push(vm.r.globalObject)
	vm.pc++
***REMOVED***

type loadStack int

func (l loadStack) exec(vm *vm) ***REMOVED***
	// l > 0 -- var<l-1>
	// l == 0 -- this

	if l > 0 ***REMOVED***
		vm.push(nilSafe(vm.stack[vm.sb+vm.args+int(l)]))
	***REMOVED*** else ***REMOVED***
		vm.push(vm.stack[vm.sb])
	***REMOVED***
	vm.pc++
***REMOVED***

type loadStack1 int

func (l loadStack1) exec(vm *vm) ***REMOVED***
	// args are in stash
	// l > 0 -- var<l-1>
	// l == 0 -- this

	if l > 0 ***REMOVED***
		vm.push(nilSafe(vm.stack[vm.sb+int(l)]))
	***REMOVED*** else ***REMOVED***
		vm.push(vm.stack[vm.sb])
	***REMOVED***
	vm.pc++
***REMOVED***

type loadStackLex int

func (l loadStackLex) exec(vm *vm) ***REMOVED***
	// l < 0 -- arg<-l-1>
	// l > 0 -- var<l-1>
	var p *Value
	if l < 0 ***REMOVED***
		arg := int(-l)
		if arg > vm.args ***REMOVED***
			vm.push(_undefined)
			vm.pc++
			return
		***REMOVED*** else ***REMOVED***
			p = &vm.stack[vm.sb+arg]
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		p = &vm.stack[vm.sb+vm.args+int(l)]
	***REMOVED***
	if *p == nil ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	vm.push(*p)
	vm.pc++
***REMOVED***

type loadStack1Lex int

func (l loadStack1Lex) exec(vm *vm) ***REMOVED***
	p := &vm.stack[vm.sb+int(l)]
	if *p == nil ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	vm.push(*p)
	vm.pc++
***REMOVED***

type _loadCallee struct***REMOVED******REMOVED***

var loadCallee _loadCallee

func (_loadCallee) exec(vm *vm) ***REMOVED***
	vm.push(vm.stack[vm.sb-1])
	vm.pc++
***REMOVED***

func (vm *vm) storeStack(s int) ***REMOVED***
	// l > 0 -- var<l-1>

	if s > 0 ***REMOVED***
		vm.stack[vm.sb+vm.args+s] = vm.stack[vm.sp-1]
	***REMOVED*** else ***REMOVED***
		panic("Illegal stack var index")
	***REMOVED***
	vm.pc++
***REMOVED***

func (vm *vm) storeStack1(s int) ***REMOVED***
	// args are in stash
	// l > 0 -- var<l-1>

	if s > 0 ***REMOVED***
		vm.stack[vm.sb+s] = vm.stack[vm.sp-1]
	***REMOVED*** else ***REMOVED***
		panic("Illegal stack var index")
	***REMOVED***
	vm.pc++
***REMOVED***

func (vm *vm) storeStackLex(s int) ***REMOVED***
	// l < 0 -- arg<-l-1>
	// l > 0 -- var<l-1>
	var p *Value
	if s < 0 ***REMOVED***
		p = &vm.stack[vm.sb-s]
	***REMOVED*** else ***REMOVED***
		p = &vm.stack[vm.sb+vm.args+s]
	***REMOVED***

	if *p != nil ***REMOVED***
		*p = vm.stack[vm.sp-1]
	***REMOVED*** else ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	vm.pc++
***REMOVED***

func (vm *vm) storeStack1Lex(s int) ***REMOVED***
	// args are in stash
	// s > 0 -- var<l-1>
	if s <= 0 ***REMOVED***
		panic("Illegal stack var index")
	***REMOVED***
	p := &vm.stack[vm.sb+s]
	if *p != nil ***REMOVED***
		*p = vm.stack[vm.sp-1]
	***REMOVED*** else ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	vm.pc++
***REMOVED***

func (vm *vm) initStack(s int) ***REMOVED***
	if s <= 0 ***REMOVED***
		panic("Illegal stack var index")
	***REMOVED***
	vm.stack[vm.sb+vm.args+s] = vm.stack[vm.sp-1]
	vm.pc++
***REMOVED***

func (vm *vm) initStack1(s int) ***REMOVED***
	if s <= 0 ***REMOVED***
		panic("Illegal stack var index")
	***REMOVED***
	vm.stack[vm.sb+s] = vm.stack[vm.sp-1]
	vm.pc++
***REMOVED***

type storeStack int

func (s storeStack) exec(vm *vm) ***REMOVED***
	vm.storeStack(int(s))
***REMOVED***

type storeStack1 int

func (s storeStack1) exec(vm *vm) ***REMOVED***
	vm.storeStack1(int(s))
***REMOVED***

type storeStackLex int

func (s storeStackLex) exec(vm *vm) ***REMOVED***
	vm.storeStackLex(int(s))
***REMOVED***

type storeStack1Lex int

func (s storeStack1Lex) exec(vm *vm) ***REMOVED***
	vm.storeStack1Lex(int(s))
***REMOVED***

type initStack int

func (s initStack) exec(vm *vm) ***REMOVED***
	vm.initStack(int(s))
	vm.sp--
***REMOVED***

type initStack1 int

func (s initStack1) exec(vm *vm) ***REMOVED***
	vm.initStack1(int(s))
	vm.sp--
***REMOVED***

type storeStackP int

func (s storeStackP) exec(vm *vm) ***REMOVED***
	vm.storeStack(int(s))
	vm.sp--
***REMOVED***

type storeStack1P int

func (s storeStack1P) exec(vm *vm) ***REMOVED***
	vm.storeStack1(int(s))
	vm.sp--
***REMOVED***

type storeStackLexP int

func (s storeStackLexP) exec(vm *vm) ***REMOVED***
	vm.storeStackLex(int(s))
	vm.sp--
***REMOVED***

type storeStack1LexP int

func (s storeStack1LexP) exec(vm *vm) ***REMOVED***
	vm.storeStack1Lex(int(s))
	vm.sp--
***REMOVED***

type _toNumber struct***REMOVED******REMOVED***

var toNumber _toNumber

func (_toNumber) exec(vm *vm) ***REMOVED***
	vm.stack[vm.sp-1] = vm.stack[vm.sp-1].ToNumber()
	vm.pc++
***REMOVED***

type _add struct***REMOVED******REMOVED***

var add _add

func (_add) exec(vm *vm) ***REMOVED***
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]

	if o, ok := left.(*Object); ok ***REMOVED***
		left = o.toPrimitive()
	***REMOVED***

	if o, ok := right.(*Object); ok ***REMOVED***
		right = o.toPrimitive()
	***REMOVED***

	var ret Value

	leftString, isLeftString := left.(valueString)
	rightString, isRightString := right.(valueString)

	if isLeftString || isRightString ***REMOVED***
		if !isLeftString ***REMOVED***
			leftString = left.toString()
		***REMOVED***
		if !isRightString ***REMOVED***
			rightString = right.toString()
		***REMOVED***
		ret = leftString.concat(rightString)
	***REMOVED*** else ***REMOVED***
		if leftInt, ok := left.(valueInt); ok ***REMOVED***
			if rightInt, ok := right.(valueInt); ok ***REMOVED***
				ret = intToValue(int64(leftInt) + int64(rightInt))
			***REMOVED*** else ***REMOVED***
				ret = floatToValue(float64(leftInt) + right.ToFloat())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			ret = floatToValue(left.ToFloat() + right.ToFloat())
		***REMOVED***
	***REMOVED***

	vm.stack[vm.sp-2] = ret
	vm.sp--
	vm.pc++
***REMOVED***

type _sub struct***REMOVED******REMOVED***

var sub _sub

func (_sub) exec(vm *vm) ***REMOVED***
	right := vm.stack[vm.sp-1]
	left := vm.stack[vm.sp-2]

	var result Value

	if left, ok := left.(valueInt); ok ***REMOVED***
		if right, ok := right.(valueInt); ok ***REMOVED***
			result = intToValue(int64(left) - int64(right))
			goto end
		***REMOVED***
	***REMOVED***

	result = floatToValue(left.ToFloat() - right.ToFloat())
end:
	vm.sp--
	vm.stack[vm.sp-1] = result
	vm.pc++
***REMOVED***

type _mul struct***REMOVED******REMOVED***

var mul _mul

func (_mul) exec(vm *vm) ***REMOVED***
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var result Value

	if left, ok := assertInt64(left); ok ***REMOVED***
		if right, ok := assertInt64(right); ok ***REMOVED***
			if left == 0 && right == -1 || left == -1 && right == 0 ***REMOVED***
				result = _negativeZero
				goto end
			***REMOVED***
			res := left * right
			// check for overflow
			if left == 0 || right == 0 || res/left == right ***REMOVED***
				result = intToValue(res)
				goto end
			***REMOVED***

		***REMOVED***
	***REMOVED***

	result = floatToValue(left.ToFloat() * right.ToFloat())

end:
	vm.sp--
	vm.stack[vm.sp-1] = result
	vm.pc++
***REMOVED***

type _div struct***REMOVED******REMOVED***

var div _div

func (_div) exec(vm *vm) ***REMOVED***
	left := vm.stack[vm.sp-2].ToFloat()
	right := vm.stack[vm.sp-1].ToFloat()

	var result Value

	if math.IsNaN(left) || math.IsNaN(right) ***REMOVED***
		result = _NaN
		goto end
	***REMOVED***
	if math.IsInf(left, 0) && math.IsInf(right, 0) ***REMOVED***
		result = _NaN
		goto end
	***REMOVED***
	if left == 0 && right == 0 ***REMOVED***
		result = _NaN
		goto end
	***REMOVED***

	if math.IsInf(left, 0) ***REMOVED***
		if math.Signbit(left) == math.Signbit(right) ***REMOVED***
			result = _positiveInf
			goto end
		***REMOVED*** else ***REMOVED***
			result = _negativeInf
			goto end
		***REMOVED***
	***REMOVED***
	if math.IsInf(right, 0) ***REMOVED***
		if math.Signbit(left) == math.Signbit(right) ***REMOVED***
			result = _positiveZero
			goto end
		***REMOVED*** else ***REMOVED***
			result = _negativeZero
			goto end
		***REMOVED***
	***REMOVED***
	if right == 0 ***REMOVED***
		if math.Signbit(left) == math.Signbit(right) ***REMOVED***
			result = _positiveInf
			goto end
		***REMOVED*** else ***REMOVED***
			result = _negativeInf
			goto end
		***REMOVED***
	***REMOVED***

	result = floatToValue(left / right)

end:
	vm.sp--
	vm.stack[vm.sp-1] = result
	vm.pc++
***REMOVED***

type _mod struct***REMOVED******REMOVED***

var mod _mod

func (_mod) exec(vm *vm) ***REMOVED***
	left := vm.stack[vm.sp-2]
	right := vm.stack[vm.sp-1]

	var result Value

	if leftInt, ok := assertInt64(left); ok ***REMOVED***
		if rightInt, ok := assertInt64(right); ok ***REMOVED***
			if rightInt == 0 ***REMOVED***
				result = _NaN
				goto end
			***REMOVED***
			r := leftInt % rightInt
			if r == 0 && leftInt < 0 ***REMOVED***
				result = _negativeZero
			***REMOVED*** else ***REMOVED***
				result = intToValue(leftInt % rightInt)
			***REMOVED***
			goto end
		***REMOVED***
	***REMOVED***

	result = floatToValue(math.Mod(left.ToFloat(), right.ToFloat()))
end:
	vm.sp--
	vm.stack[vm.sp-1] = result
	vm.pc++
***REMOVED***

type _neg struct***REMOVED******REMOVED***

var neg _neg

func (_neg) exec(vm *vm) ***REMOVED***
	operand := vm.stack[vm.sp-1]

	var result Value

	if i, ok := assertInt64(operand); ok ***REMOVED***
		if i == 0 ***REMOVED***
			result = _negativeZero
		***REMOVED*** else ***REMOVED***
			result = valueInt(-i)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		f := operand.ToFloat()
		if !math.IsNaN(f) ***REMOVED***
			f = -f
		***REMOVED***
		result = valueFloat(f)
	***REMOVED***

	vm.stack[vm.sp-1] = result
	vm.pc++
***REMOVED***

type _plus struct***REMOVED******REMOVED***

var plus _plus

func (_plus) exec(vm *vm) ***REMOVED***
	vm.stack[vm.sp-1] = vm.stack[vm.sp-1].ToNumber()
	vm.pc++
***REMOVED***

type _inc struct***REMOVED******REMOVED***

var inc _inc

func (_inc) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-1]

	if i, ok := assertInt64(v); ok ***REMOVED***
		v = intToValue(i + 1)
		goto end
	***REMOVED***

	v = valueFloat(v.ToFloat() + 1)

end:
	vm.stack[vm.sp-1] = v
	vm.pc++
***REMOVED***

type _dec struct***REMOVED******REMOVED***

var dec _dec

func (_dec) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-1]

	if i, ok := assertInt64(v); ok ***REMOVED***
		v = intToValue(i - 1)
		goto end
	***REMOVED***

	v = valueFloat(v.ToFloat() - 1)

end:
	vm.stack[vm.sp-1] = v
	vm.pc++
***REMOVED***

type _and struct***REMOVED******REMOVED***

var and _and

func (_and) exec(vm *vm) ***REMOVED***
	left := toInt32(vm.stack[vm.sp-2])
	right := toInt32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left & right))
	vm.sp--
	vm.pc++
***REMOVED***

type _or struct***REMOVED******REMOVED***

var or _or

func (_or) exec(vm *vm) ***REMOVED***
	left := toInt32(vm.stack[vm.sp-2])
	right := toInt32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left | right))
	vm.sp--
	vm.pc++
***REMOVED***

type _xor struct***REMOVED******REMOVED***

var xor _xor

func (_xor) exec(vm *vm) ***REMOVED***
	left := toInt32(vm.stack[vm.sp-2])
	right := toInt32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left ^ right))
	vm.sp--
	vm.pc++
***REMOVED***

type _bnot struct***REMOVED******REMOVED***

var bnot _bnot

func (_bnot) exec(vm *vm) ***REMOVED***
	op := toInt32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-1] = intToValue(int64(^op))
	vm.pc++
***REMOVED***

type _sal struct***REMOVED******REMOVED***

var sal _sal

func (_sal) exec(vm *vm) ***REMOVED***
	left := toInt32(vm.stack[vm.sp-2])
	right := toUint32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left << (right & 0x1F)))
	vm.sp--
	vm.pc++
***REMOVED***

type _sar struct***REMOVED******REMOVED***

var sar _sar

func (_sar) exec(vm *vm) ***REMOVED***
	left := toInt32(vm.stack[vm.sp-2])
	right := toUint32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left >> (right & 0x1F)))
	vm.sp--
	vm.pc++
***REMOVED***

type _shr struct***REMOVED******REMOVED***

var shr _shr

func (_shr) exec(vm *vm) ***REMOVED***
	left := toUint32(vm.stack[vm.sp-2])
	right := toUint32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left >> (right & 0x1F)))
	vm.sp--
	vm.pc++
***REMOVED***

type _halt struct***REMOVED******REMOVED***

var halt _halt

func (_halt) exec(vm *vm) ***REMOVED***
	vm.halt = true
	vm.pc++
***REMOVED***

type jump int32

func (j jump) exec(vm *vm) ***REMOVED***
	vm.pc += int(j)
***REMOVED***

type _setElem struct***REMOVED******REMOVED***

var setElem _setElem

func (_setElem) exec(vm *vm) ***REMOVED***
	obj := vm.stack[vm.sp-3].ToObject(vm.r)
	propName := toPropertyKey(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	obj.setOwn(propName, val, false)

	vm.sp -= 2
	vm.stack[vm.sp-1] = val
	vm.pc++
***REMOVED***

type _setElemP struct***REMOVED******REMOVED***

var setElemP _setElemP

func (_setElemP) exec(vm *vm) ***REMOVED***
	obj := vm.stack[vm.sp-3].ToObject(vm.r)
	propName := toPropertyKey(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	obj.setOwn(propName, val, false)

	vm.sp -= 3
	vm.pc++
***REMOVED***

type _setElemStrict struct***REMOVED******REMOVED***

var setElemStrict _setElemStrict

func (_setElemStrict) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-3])
	propName := toPropertyKey(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	obj.setOwn(propName, val, true)

	vm.sp -= 2
	vm.stack[vm.sp-1] = val
	vm.pc++
***REMOVED***

type _setElemStrictP struct***REMOVED******REMOVED***

var setElemStrictP _setElemStrictP

func (_setElemStrictP) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-3])
	propName := toPropertyKey(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	obj.setOwn(propName, val, true)

	vm.sp -= 3
	vm.pc++
***REMOVED***

type _deleteElem struct***REMOVED******REMOVED***

var deleteElem _deleteElem

func (_deleteElem) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	propName := toPropertyKey(vm.stack[vm.sp-1])
	if obj.delete(propName, false) ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type _deleteElemStrict struct***REMOVED******REMOVED***

var deleteElemStrict _deleteElemStrict

func (_deleteElemStrict) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	propName := toPropertyKey(vm.stack[vm.sp-1])
	obj.delete(propName, true)
	vm.stack[vm.sp-2] = valueTrue
	vm.sp--
	vm.pc++
***REMOVED***

type deleteProp unistring.String

func (d deleteProp) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-1])
	if obj.self.deleteStr(unistring.String(d), false) ***REMOVED***
		vm.stack[vm.sp-1] = valueTrue
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-1] = valueFalse
	***REMOVED***
	vm.pc++
***REMOVED***

type deletePropStrict unistring.String

func (d deletePropStrict) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-1])
	obj.self.deleteStr(unistring.String(d), true)
	vm.stack[vm.sp-1] = valueTrue
	vm.pc++
***REMOVED***

type setProp unistring.String

func (p setProp) exec(vm *vm) ***REMOVED***
	val := vm.stack[vm.sp-1]
	vm.stack[vm.sp-2].ToObject(vm.r).self.setOwnStr(unistring.String(p), val, false)
	vm.stack[vm.sp-2] = val
	vm.sp--
	vm.pc++
***REMOVED***

type setPropP unistring.String

func (p setPropP) exec(vm *vm) ***REMOVED***
	val := vm.stack[vm.sp-1]
	vm.stack[vm.sp-2].ToObject(vm.r).self.setOwnStr(unistring.String(p), val, false)
	vm.sp -= 2
	vm.pc++
***REMOVED***

type setPropStrict unistring.String

func (p setPropStrict) exec(vm *vm) ***REMOVED***
	obj := vm.stack[vm.sp-2]
	val := vm.stack[vm.sp-1]

	obj1 := vm.r.toObject(obj)
	obj1.self.setOwnStr(unistring.String(p), val, true)
	vm.stack[vm.sp-2] = val
	vm.sp--
	vm.pc++
***REMOVED***

type setPropStrictP unistring.String

func (p setPropStrictP) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	obj.self.setOwnStr(unistring.String(p), val, true)
	vm.sp -= 2
	vm.pc++
***REMOVED***

type setProp1 unistring.String

func (p setProp1) exec(vm *vm) ***REMOVED***
	vm.r.toObject(vm.stack[vm.sp-2]).self._putProp(unistring.String(p), vm.stack[vm.sp-1], true, true, true)

	vm.sp--
	vm.pc++
***REMOVED***

type _setProto struct***REMOVED******REMOVED***

var setProto _setProto

func (_setProto) exec(vm *vm) ***REMOVED***
	vm.r.toObject(vm.stack[vm.sp-2]).self.setProto(vm.r.toProto(vm.stack[vm.sp-1]), true)

	vm.sp--
	vm.pc++
***REMOVED***

type setPropGetter unistring.String

func (s setPropGetter) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	descr := PropertyDescriptor***REMOVED***
		Getter:       val,
		Configurable: FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
	***REMOVED***

	obj.self.defineOwnPropertyStr(unistring.String(s), descr, false)

	vm.sp--
	vm.pc++
***REMOVED***

type setPropSetter unistring.String

func (s setPropSetter) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	descr := PropertyDescriptor***REMOVED***
		Setter:       val,
		Configurable: FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
	***REMOVED***

	obj.self.defineOwnPropertyStr(unistring.String(s), descr, false)

	vm.sp--
	vm.pc++
***REMOVED***

type getProp unistring.String

func (g getProp) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-1]
	obj := v.baseObject(vm.r)
	if obj == nil ***REMOVED***
		panic(vm.r.NewTypeError("Cannot read property '%s' of undefined", g))
	***REMOVED***
	vm.stack[vm.sp-1] = nilSafe(obj.self.getStr(unistring.String(g), v))

	vm.pc++
***REMOVED***

type getPropCallee unistring.String

func (g getPropCallee) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-1]
	obj := v.baseObject(vm.r)
	n := unistring.String(g)
	if obj == nil ***REMOVED***
		panic(vm.r.NewTypeError("Cannot read property '%s' of undefined or null", n))
	***REMOVED***
	prop := obj.self.getStr(n, v)
	if prop == nil ***REMOVED***
		prop = memberUnresolved***REMOVED***valueUnresolved***REMOVED***r: vm.r, ref: n***REMOVED******REMOVED***
	***REMOVED***
	vm.stack[vm.sp-1] = prop

	vm.pc++
***REMOVED***

type _getElem struct***REMOVED******REMOVED***

var getElem _getElem

func (_getElem) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-2]
	obj := v.baseObject(vm.r)
	propName := toPropertyKey(vm.stack[vm.sp-1])
	if obj == nil ***REMOVED***
		panic(vm.r.NewTypeError("Cannot read property '%s' of undefined", propName.String()))
	***REMOVED***

	vm.stack[vm.sp-2] = nilSafe(obj.get(propName, v))

	vm.sp--
	vm.pc++
***REMOVED***

type _getElemCallee struct***REMOVED******REMOVED***

var getElemCallee _getElemCallee

func (_getElemCallee) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-2]
	obj := v.baseObject(vm.r)
	propName := toPropertyKey(vm.stack[vm.sp-1])
	if obj == nil ***REMOVED***
		panic(vm.r.NewTypeError("Cannot read property '%s' of undefined", propName.String()))
	***REMOVED***

	prop := obj.get(propName, v)
	if prop == nil ***REMOVED***
		prop = memberUnresolved***REMOVED***valueUnresolved***REMOVED***r: vm.r, ref: propName.string()***REMOVED******REMOVED***
	***REMOVED***
	vm.stack[vm.sp-2] = prop

	vm.sp--
	vm.pc++
***REMOVED***

type _dup struct***REMOVED******REMOVED***

var dup _dup

func (_dup) exec(vm *vm) ***REMOVED***
	vm.push(vm.stack[vm.sp-1])
	vm.pc++
***REMOVED***

type dupN uint32

func (d dupN) exec(vm *vm) ***REMOVED***
	vm.push(vm.stack[vm.sp-1-int(d)])
	vm.pc++
***REMOVED***

type rdupN uint32

func (d rdupN) exec(vm *vm) ***REMOVED***
	vm.stack[vm.sp-1-int(d)] = vm.stack[vm.sp-1]
	vm.pc++
***REMOVED***

type _newObject struct***REMOVED******REMOVED***

var newObject _newObject

func (_newObject) exec(vm *vm) ***REMOVED***
	vm.push(vm.r.NewObject())
	vm.pc++
***REMOVED***

type newArray uint32

func (l newArray) exec(vm *vm) ***REMOVED***
	values := make([]Value, l)
	if l > 0 ***REMOVED***
		copy(values, vm.stack[vm.sp-int(l):vm.sp])
	***REMOVED***
	obj := vm.r.newArrayValues(values)
	if l > 0 ***REMOVED***
		vm.sp -= int(l) - 1
		vm.stack[vm.sp-1] = obj
	***REMOVED*** else ***REMOVED***
		vm.push(obj)
	***REMOVED***
	vm.pc++
***REMOVED***

type newArraySparse struct ***REMOVED***
	l, objCount int
***REMOVED***

func (n *newArraySparse) exec(vm *vm) ***REMOVED***
	values := make([]Value, n.l)
	copy(values, vm.stack[vm.sp-int(n.l):vm.sp])
	arr := vm.r.newArrayObject()
	setArrayValues(arr, values)
	arr.objCount = n.objCount
	vm.sp -= int(n.l) - 1
	vm.stack[vm.sp-1] = arr.val
	vm.pc++
***REMOVED***

type newRegexp struct ***REMOVED***
	pattern *regexpPattern
	src     valueString
***REMOVED***

func (n *newRegexp) exec(vm *vm) ***REMOVED***
	vm.push(vm.r.newRegExpp(n.pattern.clone(), n.src, vm.r.global.RegExpPrototype).val)
	vm.pc++
***REMOVED***

func (vm *vm) setLocalLex(s int) ***REMOVED***
	v := vm.stack[vm.sp-1]
	level := s >> 24
	idx := uint32(s & 0x00FFFFFF)
	stash := vm.stash
	for i := 0; i < level; i++ ***REMOVED***
		stash = stash.outer
	***REMOVED***
	p := &stash.values[idx]
	if *p == nil ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	*p = v
	vm.pc++
***REMOVED***

func (vm *vm) initLocal(s int) ***REMOVED***
	v := vm.stack[vm.sp-1]
	level := s >> 24
	idx := uint32(s & 0x00FFFFFF)
	stash := vm.stash
	for i := 0; i < level; i++ ***REMOVED***
		stash = stash.outer
	***REMOVED***
	stash.initByIdx(idx, v)
	vm.pc++
***REMOVED***

type storeStash uint32

func (s storeStash) exec(vm *vm) ***REMOVED***
	vm.initLocal(int(s))
***REMOVED***

type storeStashP uint32

func (s storeStashP) exec(vm *vm) ***REMOVED***
	vm.initLocal(int(s))
	vm.sp--
***REMOVED***

type storeStashLex uint32

func (s storeStashLex) exec(vm *vm) ***REMOVED***
	vm.setLocalLex(int(s))
***REMOVED***

type storeStashLexP uint32

func (s storeStashLexP) exec(vm *vm) ***REMOVED***
	vm.setLocalLex(int(s))
	vm.sp--
***REMOVED***

type initStash uint32

func (s initStash) exec(vm *vm) ***REMOVED***
	vm.initLocal(int(s))
	vm.sp--
***REMOVED***

type initGlobal unistring.String

func (s initGlobal) exec(vm *vm) ***REMOVED***
	vm.sp--
	vm.r.global.stash.initByName(unistring.String(s), vm.stack[vm.sp])
	vm.pc++
***REMOVED***

type resolveVar1 unistring.String

func (s resolveVar1) exec(vm *vm) ***REMOVED***
	name := unistring.String(s)
	var ref ref
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		ref = stash.getRefByName(name, false)
		if ref != nil ***REMOVED***
			goto end
		***REMOVED***
	***REMOVED***

	ref = &objRef***REMOVED***
		base: vm.r.globalObject.self,
		name: name,
	***REMOVED***

end:
	vm.refStack = append(vm.refStack, ref)
	vm.pc++
***REMOVED***

type deleteVar unistring.String

func (d deleteVar) exec(vm *vm) ***REMOVED***
	name := unistring.String(d)
	ret := true
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if stash.obj != nil ***REMOVED***
			if stashObjHas(stash.obj, name) ***REMOVED***
				ret = stash.obj.self.deleteStr(name, false)
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if idx, exists := stash.names[name]; exists ***REMOVED***
				if idx&(maskVar|maskDeletable) == maskVar|maskDeletable ***REMOVED***
					stash.deleteBinding(name)
				***REMOVED*** else ***REMOVED***
					ret = false
				***REMOVED***
				goto end
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if vm.r.globalObject.self.hasPropertyStr(name) ***REMOVED***
		ret = vm.r.globalObject.self.deleteStr(name, false)
	***REMOVED***

end:
	if ret ***REMOVED***
		vm.push(valueTrue)
	***REMOVED*** else ***REMOVED***
		vm.push(valueFalse)
	***REMOVED***
	vm.pc++
***REMOVED***

type deleteGlobal unistring.String

func (d deleteGlobal) exec(vm *vm) ***REMOVED***
	name := unistring.String(d)
	var ret bool
	if vm.r.globalObject.self.hasPropertyStr(name) ***REMOVED***
		ret = vm.r.globalObject.self.deleteStr(name, false)
		if ret ***REMOVED***
			delete(vm.r.global.varNames, name)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ret = true
	***REMOVED***
	if ret ***REMOVED***
		vm.push(valueTrue)
	***REMOVED*** else ***REMOVED***
		vm.push(valueFalse)
	***REMOVED***
	vm.pc++
***REMOVED***

type resolveVar1Strict unistring.String

func (s resolveVar1Strict) exec(vm *vm) ***REMOVED***
	name := unistring.String(s)
	var ref ref
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		ref = stash.getRefByName(name, true)
		if ref != nil ***REMOVED***
			goto end
		***REMOVED***
	***REMOVED***

	if vm.r.globalObject.self.hasPropertyStr(name) ***REMOVED***
		ref = &objRef***REMOVED***
			base:   vm.r.globalObject.self,
			name:   name,
			strict: true,
		***REMOVED***
		goto end
	***REMOVED***

	ref = &unresolvedRef***REMOVED***
		runtime: vm.r,
		name:    name,
	***REMOVED***

end:
	vm.refStack = append(vm.refStack, ref)
	vm.pc++
***REMOVED***

type setGlobal unistring.String

func (s setGlobal) exec(vm *vm) ***REMOVED***
	vm.r.setGlobal(unistring.String(s), vm.peek(), false)
	vm.pc++
***REMOVED***

type setGlobalStrict unistring.String

func (s setGlobalStrict) exec(vm *vm) ***REMOVED***
	vm.r.setGlobal(unistring.String(s), vm.peek(), true)
	vm.pc++
***REMOVED***

// Load a var from stash
type loadStash uint32

func (g loadStash) exec(vm *vm) ***REMOVED***
	level := int(g >> 24)
	idx := uint32(g & 0x00FFFFFF)
	stash := vm.stash
	for i := 0; i < level; i++ ***REMOVED***
		stash = stash.outer
	***REMOVED***

	vm.push(nilSafe(stash.getByIdx(idx)))
	vm.pc++
***REMOVED***

// Load a lexical binding from stash
type loadStashLex uint32

func (g loadStashLex) exec(vm *vm) ***REMOVED***
	level := int(g >> 24)
	idx := uint32(g & 0x00FFFFFF)
	stash := vm.stash
	for i := 0; i < level; i++ ***REMOVED***
		stash = stash.outer
	***REMOVED***

	v := stash.getByIdx(idx)
	if v == nil ***REMOVED***
		panic(errAccessBeforeInit)
	***REMOVED***
	vm.push(v)
	vm.pc++
***REMOVED***

// scan dynamic stashes up to the given level (encoded as 8 most significant bits of idx), if not found
// return the indexed var binding value from stash
type loadMixed struct ***REMOVED***
	name   unistring.String
	idx    uint32
	callee bool
***REMOVED***

func (g *loadMixed) exec(vm *vm) ***REMOVED***
	level := int(g.idx >> 24)
	idx := g.idx & 0x00FFFFFF
	stash := vm.stash
	name := g.name
	for i := 0; i < level; i++ ***REMOVED***
		if v, found := stash.getByName(name); found ***REMOVED***
			if g.callee ***REMOVED***
				if stash.obj != nil ***REMOVED***
					vm.push(stash.obj)
				***REMOVED*** else ***REMOVED***
					vm.push(_undefined)
				***REMOVED***
			***REMOVED***
			vm.push(v)
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***
	if g.callee ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	if stash != nil ***REMOVED***
		vm.push(nilSafe(stash.getByIdx(idx)))
	***REMOVED***
end:
	vm.pc++
***REMOVED***

// scan dynamic stashes up to the given level (encoded as 8 most significant bits of idx), if not found
// return the indexed lexical binding value from stash
type loadMixedLex loadMixed

func (g *loadMixedLex) exec(vm *vm) ***REMOVED***
	level := int(g.idx >> 24)
	idx := g.idx & 0x00FFFFFF
	stash := vm.stash
	name := g.name
	for i := 0; i < level; i++ ***REMOVED***
		if v, found := stash.getByName(name); found ***REMOVED***
			if g.callee ***REMOVED***
				if stash.obj != nil ***REMOVED***
					vm.push(stash.obj)
				***REMOVED*** else ***REMOVED***
					vm.push(_undefined)
				***REMOVED***
			***REMOVED***
			vm.push(v)
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***
	if g.callee ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	if stash != nil ***REMOVED***
		v := stash.getByIdx(idx)
		if v == nil ***REMOVED***
			panic(errAccessBeforeInit)
		***REMOVED***
		vm.push(v)
	***REMOVED***
end:
	vm.pc++
***REMOVED***

// scan dynamic stashes up to the given level (encoded as 8 most significant bits of idx), if not found
// return the indexed var binding value from stack
type loadMixedStack struct ***REMOVED***
	name   unistring.String
	idx    int
	level  uint8
	callee bool
***REMOVED***

// same as loadMixedStack, but the args have been moved to stash (therefore stack layout is different)
type loadMixedStack1 loadMixedStack

func (g *loadMixedStack) exec(vm *vm) ***REMOVED***
	stash := vm.stash
	name := g.name
	level := int(g.level)
	for i := 0; i < level; i++ ***REMOVED***
		if v, found := stash.getByName(name); found ***REMOVED***
			if g.callee ***REMOVED***
				if stash.obj != nil ***REMOVED***
					vm.push(stash.obj)
				***REMOVED*** else ***REMOVED***
					vm.push(_undefined)
				***REMOVED***
			***REMOVED***
			vm.push(v)
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***
	if g.callee ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	loadStack(g.idx).exec(vm)
	return
end:
	vm.pc++
***REMOVED***

func (g *loadMixedStack1) exec(vm *vm) ***REMOVED***
	stash := vm.stash
	name := g.name
	level := int(g.level)
	for i := 0; i < level; i++ ***REMOVED***
		if v, found := stash.getByName(name); found ***REMOVED***
			if g.callee ***REMOVED***
				if stash.obj != nil ***REMOVED***
					vm.push(stash.obj)
				***REMOVED*** else ***REMOVED***
					vm.push(_undefined)
				***REMOVED***
			***REMOVED***
			vm.push(v)
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***
	if g.callee ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	loadStack1(g.idx).exec(vm)
	return
end:
	vm.pc++
***REMOVED***

type loadMixedStackLex loadMixedStack

// same as loadMixedStackLex but when the arguments have been moved into stash
type loadMixedStack1Lex loadMixedStack

func (g *loadMixedStackLex) exec(vm *vm) ***REMOVED***
	stash := vm.stash
	name := g.name
	level := int(g.level)
	for i := 0; i < level; i++ ***REMOVED***
		if v, found := stash.getByName(name); found ***REMOVED***
			if g.callee ***REMOVED***
				if stash.obj != nil ***REMOVED***
					vm.push(stash.obj)
				***REMOVED*** else ***REMOVED***
					vm.push(_undefined)
				***REMOVED***
			***REMOVED***
			vm.push(v)
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***
	if g.callee ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	loadStackLex(g.idx).exec(vm)
	return
end:
	vm.pc++
***REMOVED***

func (g *loadMixedStack1Lex) exec(vm *vm) ***REMOVED***
	stash := vm.stash
	name := g.name
	level := int(g.level)
	for i := 0; i < level; i++ ***REMOVED***
		if v, found := stash.getByName(name); found ***REMOVED***
			if g.callee ***REMOVED***
				if stash.obj != nil ***REMOVED***
					vm.push(stash.obj)
				***REMOVED*** else ***REMOVED***
					vm.push(_undefined)
				***REMOVED***
			***REMOVED***
			vm.push(v)
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***
	if g.callee ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	loadStack1Lex(g.idx).exec(vm)
	return
end:
	vm.pc++
***REMOVED***

type resolveMixed struct ***REMOVED***
	name   unistring.String
	idx    uint32
	typ    varType
	strict bool
***REMOVED***

func newStashRef(typ varType, name unistring.String, v *[]Value, idx int) ref ***REMOVED***
	switch typ ***REMOVED***
	case varTypeVar:
		return &stashRef***REMOVED***
			n:   name,
			v:   v,
			idx: idx,
		***REMOVED***
	case varTypeLet:
		return &stashRefLex***REMOVED***
			stashRef: stashRef***REMOVED***
				n:   name,
				v:   v,
				idx: idx,
			***REMOVED***,
		***REMOVED***
	case varTypeConst:
		return &stashRefConst***REMOVED***
			stashRefLex: stashRefLex***REMOVED***
				stashRef: stashRef***REMOVED***
					n:   name,
					v:   v,
					idx: idx,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
	***REMOVED***
	panic("unsupported var type")
***REMOVED***

func (r *resolveMixed) exec(vm *vm) ***REMOVED***
	level := int(r.idx >> 24)
	idx := r.idx & 0x00FFFFFF
	stash := vm.stash
	var ref ref
	for i := 0; i < level; i++ ***REMOVED***
		ref = stash.getRefByName(r.name, r.strict)
		if ref != nil ***REMOVED***
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***

	if stash != nil ***REMOVED***
		ref = newStashRef(r.typ, r.name, &stash.values, int(idx))
		goto end
	***REMOVED***

	ref = &unresolvedRef***REMOVED***
		runtime: vm.r,
		name:    r.name,
	***REMOVED***

end:
	vm.refStack = append(vm.refStack, ref)
	vm.pc++
***REMOVED***

type resolveMixedStack struct ***REMOVED***
	name   unistring.String
	idx    int
	typ    varType
	level  uint8
	strict bool
***REMOVED***

type resolveMixedStack1 resolveMixedStack

func (r *resolveMixedStack) exec(vm *vm) ***REMOVED***
	level := int(r.level)
	stash := vm.stash
	var ref ref
	var idx int
	for i := 0; i < level; i++ ***REMOVED***
		ref = stash.getRefByName(r.name, r.strict)
		if ref != nil ***REMOVED***
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***

	if r.idx > 0 ***REMOVED***
		idx = vm.sb + vm.args + r.idx
	***REMOVED*** else ***REMOVED***
		idx = vm.sb + r.idx
	***REMOVED***

	ref = newStashRef(r.typ, r.name, (*[]Value)(&vm.stack), idx)

end:
	vm.refStack = append(vm.refStack, ref)
	vm.pc++
***REMOVED***

func (r *resolveMixedStack1) exec(vm *vm) ***REMOVED***
	level := int(r.level)
	stash := vm.stash
	var ref ref
	for i := 0; i < level; i++ ***REMOVED***
		ref = stash.getRefByName(r.name, r.strict)
		if ref != nil ***REMOVED***
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***

	ref = newStashRef(r.typ, r.name, (*[]Value)(&vm.stack), vm.sb+r.idx)

end:
	vm.refStack = append(vm.refStack, ref)
	vm.pc++
***REMOVED***

type _getValue struct***REMOVED******REMOVED***

var getValue _getValue

func (_getValue) exec(vm *vm) ***REMOVED***
	ref := vm.refStack[len(vm.refStack)-1]
	if v := ref.get(); v != nil ***REMOVED***
		vm.push(v)
	***REMOVED*** else ***REMOVED***
		vm.r.throwReferenceError(ref.refname())
		panic("Unreachable")
	***REMOVED***
	vm.pc++
***REMOVED***

type _putValue struct***REMOVED******REMOVED***

var putValue _putValue

func (_putValue) exec(vm *vm) ***REMOVED***
	l := len(vm.refStack) - 1
	ref := vm.refStack[l]
	vm.refStack[l] = nil
	vm.refStack = vm.refStack[:l]
	ref.set(vm.stack[vm.sp-1])
	vm.pc++
***REMOVED***

type _putValueP struct***REMOVED******REMOVED***

var putValueP _putValueP

func (_putValueP) exec(vm *vm) ***REMOVED***
	l := len(vm.refStack) - 1
	ref := vm.refStack[l]
	vm.refStack[l] = nil
	vm.refStack = vm.refStack[:l]
	ref.set(vm.stack[vm.sp-1])
	vm.sp--
	vm.pc++
***REMOVED***

type loadDynamic unistring.String

func (n loadDynamic) exec(vm *vm) ***REMOVED***
	name := unistring.String(n)
	var val Value
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if v, exists := stash.getByName(name); exists ***REMOVED***
			val = v
			break
		***REMOVED***
	***REMOVED***
	if val == nil ***REMOVED***
		val = vm.r.globalObject.self.getStr(name, nil)
		if val == nil ***REMOVED***
			vm.r.throwReferenceError(name)
		***REMOVED***
	***REMOVED***
	vm.push(val)
	vm.pc++
***REMOVED***

type loadDynamicRef unistring.String

func (n loadDynamicRef) exec(vm *vm) ***REMOVED***
	name := unistring.String(n)
	var val Value
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if v, exists := stash.getByName(name); exists ***REMOVED***
			val = v
			break
		***REMOVED***
	***REMOVED***
	if val == nil ***REMOVED***
		val = vm.r.globalObject.self.getStr(name, nil)
		if val == nil ***REMOVED***
			val = valueUnresolved***REMOVED***r: vm.r, ref: name***REMOVED***
		***REMOVED***
	***REMOVED***
	vm.push(val)
	vm.pc++
***REMOVED***

type loadDynamicCallee unistring.String

func (n loadDynamicCallee) exec(vm *vm) ***REMOVED***
	name := unistring.String(n)
	var val Value
	var callee *Object
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if v, exists := stash.getByName(name); exists ***REMOVED***
			callee = stash.obj
			val = v
			break
		***REMOVED***
	***REMOVED***
	if val == nil ***REMOVED***
		val = vm.r.globalObject.self.getStr(name, nil)
		if val == nil ***REMOVED***
			val = valueUnresolved***REMOVED***r: vm.r, ref: name***REMOVED***
		***REMOVED***
	***REMOVED***
	if callee != nil ***REMOVED***
		vm.push(callee)
	***REMOVED*** else ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	vm.push(val)
	vm.pc++
***REMOVED***

type _pop struct***REMOVED******REMOVED***

var pop _pop

func (_pop) exec(vm *vm) ***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

func (vm *vm) callEval(n int, strict bool) ***REMOVED***
	if vm.r.toObject(vm.stack[vm.sp-n-1]) == vm.r.global.Eval ***REMOVED***
		if n > 0 ***REMOVED***
			srcVal := vm.stack[vm.sp-n]
			if src, ok := srcVal.(valueString); ok ***REMOVED***
				var this Value
				if vm.sb >= 0 ***REMOVED***
					this = vm.stack[vm.sb]
				***REMOVED*** else ***REMOVED***
					this = vm.r.globalObject
				***REMOVED***
				ret := vm.r.eval(src, true, strict, this)
				vm.stack[vm.sp-n-2] = ret
			***REMOVED*** else ***REMOVED***
				vm.stack[vm.sp-n-2] = srcVal
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			vm.stack[vm.sp-n-2] = _undefined
		***REMOVED***

		vm.sp -= n + 1
		vm.pc++
	***REMOVED*** else ***REMOVED***
		call(n).exec(vm)
	***REMOVED***
***REMOVED***

type callEval uint32

func (numargs callEval) exec(vm *vm) ***REMOVED***
	vm.callEval(int(numargs), false)
***REMOVED***

type callEvalStrict uint32

func (numargs callEvalStrict) exec(vm *vm) ***REMOVED***
	vm.callEval(int(numargs), true)
***REMOVED***

type _boxThis struct***REMOVED******REMOVED***

var boxThis _boxThis

func (_boxThis) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sb]
	if v == _undefined || v == _null ***REMOVED***
		vm.stack[vm.sb] = vm.r.globalObject
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sb] = v.ToObject(vm.r)
	***REMOVED***
	vm.pc++
***REMOVED***

type call uint32

func (numargs call) exec(vm *vm) ***REMOVED***
	// this
	// callee
	// arg0
	// ...
	// arg<numargs-1>
	n := int(numargs)
	v := vm.stack[vm.sp-n-1] // callee
	obj := vm.toCallee(v)
repeat:
	switch f := obj.self.(type) ***REMOVED***
	case *funcObject:
		vm.pc++
		vm.pushCtx()
		vm.args = n
		vm.prg = f.prg
		vm.stash = f.stash
		vm.pc = 0
		vm.stack[vm.sp-n-1], vm.stack[vm.sp-n-2] = vm.stack[vm.sp-n-2], vm.stack[vm.sp-n-1]
		return
	case *nativeFuncObject:
		vm._nativeCall(f, n)
	case *boundFuncObject:
		vm._nativeCall(&f.nativeFuncObject, n)
	case *proxyObject:
		vm.pushCtx()
		vm.prg = nil
		vm.funcName = "proxy"
		ret := f.apply(FunctionCall***REMOVED***This: vm.stack[vm.sp-n-2], Arguments: vm.stack[vm.sp-n : vm.sp]***REMOVED***)
		if ret == nil ***REMOVED***
			ret = _undefined
		***REMOVED***
		vm.stack[vm.sp-n-2] = ret
		vm.popCtx()
		vm.sp -= n + 1
		vm.pc++
	case *lazyObject:
		obj.self = f.create(obj)
		goto repeat
	default:
		vm.r.typeErrorResult(true, "Not a function: %s", obj.toString())
	***REMOVED***
***REMOVED***

func (vm *vm) _nativeCall(f *nativeFuncObject, n int) ***REMOVED***
	if f.f != nil ***REMOVED***
		vm.pushCtx()
		vm.prg = nil
		vm.funcName = f.nameProp.get(nil).string()
		ret := f.f(FunctionCall***REMOVED***
			Arguments: vm.stack[vm.sp-n : vm.sp],
			This:      vm.stack[vm.sp-n-2],
		***REMOVED***)
		if ret == nil ***REMOVED***
			ret = _undefined
		***REMOVED***
		vm.stack[vm.sp-n-2] = ret
		vm.popCtx()
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-n-2] = _undefined
	***REMOVED***
	vm.sp -= n + 1
	vm.pc++
***REMOVED***

func (vm *vm) clearStack() ***REMOVED***
	sp := vm.sp
	stackTail := vm.stack[sp:]
	for i := range stackTail ***REMOVED***
		stackTail[i] = nil
	***REMOVED***
	vm.stack = vm.stack[:sp]
***REMOVED***

type enterBlock struct ***REMOVED***
	names     map[unistring.String]uint32
	stashSize uint32
	stackSize uint32
***REMOVED***

func (e *enterBlock) exec(vm *vm) ***REMOVED***
	if e.stashSize > 0 ***REMOVED***
		vm.newStash()
		vm.stash.values = make([]Value, e.stashSize)
		if len(e.names) > 0 ***REMOVED***
			vm.stash.names = e.names
		***REMOVED***
	***REMOVED***
	ss := int(e.stackSize)
	vm.stack.expand(vm.sp + ss - 1)
	vv := vm.stack[vm.sp : vm.sp+ss]
	for i := range vv ***REMOVED***
		vv[i] = nil
	***REMOVED***
	vm.sp += ss
	vm.pc++
***REMOVED***

type enterCatchBlock struct ***REMOVED***
	names     map[unistring.String]uint32
	stashSize uint32
	stackSize uint32
***REMOVED***

func (e *enterCatchBlock) exec(vm *vm) ***REMOVED***
	vm.newStash()
	vm.stash.values = make([]Value, e.stashSize)
	if len(e.names) > 0 ***REMOVED***
		vm.stash.names = e.names
	***REMOVED***
	vm.sp--
	vm.stash.values[0] = vm.stack[vm.sp]
	ss := int(e.stackSize)
	vm.stack.expand(vm.sp + ss - 1)
	vv := vm.stack[vm.sp : vm.sp+ss]
	for i := range vv ***REMOVED***
		vv[i] = nil
	***REMOVED***
	vm.sp += ss
	vm.pc++
***REMOVED***

type leaveBlock struct ***REMOVED***
	stackSize uint32
	popStash  bool
***REMOVED***

func (l *leaveBlock) exec(vm *vm) ***REMOVED***
	if l.popStash ***REMOVED***
		vm.stash = vm.stash.outer
	***REMOVED***
	if ss := l.stackSize; ss > 0 ***REMOVED***
		vm.sp -= int(ss)
	***REMOVED***
	vm.pc++
***REMOVED***

type enterFunc struct ***REMOVED***
	names       map[unistring.String]uint32
	stashSize   uint32
	stackSize   uint32
	numArgs     uint32
	argsToStash bool
	extensible  bool
***REMOVED***

func (e *enterFunc) exec(vm *vm) ***REMOVED***
	// Input stack:
	//
	// callee
	// this
	// arg0
	// ...
	// argN
	// <- sp

	// Output stack:
	//
	// this <- sb
	// <local stack vars...>
	// <- sp
	sp := vm.sp
	vm.sb = sp - vm.args - 1
	vm.newStash()
	stash := vm.stash
	stash.function = true
	stash.values = make([]Value, e.stashSize)
	if len(e.names) > 0 ***REMOVED***
		if e.extensible ***REMOVED***
			m := make(map[unistring.String]uint32, len(e.names))
			for name, idx := range e.names ***REMOVED***
				m[name] = idx
			***REMOVED***
			stash.names = m
		***REMOVED*** else ***REMOVED***
			stash.names = e.names
		***REMOVED***
	***REMOVED***

	ss := int(e.stackSize)
	ea := 0
	if e.argsToStash ***REMOVED***
		offset := vm.args - int(e.numArgs)
		copy(stash.values, vm.stack[sp-vm.args:sp])
		if offset > 0 ***REMOVED***
			vm.stash.extraArgs = make([]Value, offset)
			copy(stash.extraArgs, vm.stack[sp-offset:])
		***REMOVED*** else ***REMOVED***
			vv := stash.values[vm.args:e.numArgs]
			for i := range vv ***REMOVED***
				vv[i] = _undefined
			***REMOVED***
		***REMOVED***
		sp -= vm.args
	***REMOVED*** else ***REMOVED***
		d := int(e.numArgs) - vm.args
		if d > 0 ***REMOVED***
			ss += d
			ea = d
			vm.args = int(e.numArgs)
		***REMOVED***
	***REMOVED***
	vm.stack.expand(sp + ss - 1)
	if ea > 0 ***REMOVED***
		vv := vm.stack[sp : vm.sp+ea]
		for i := range vv ***REMOVED***
			vv[i] = _undefined
		***REMOVED***
	***REMOVED***
	vv := vm.stack[sp+ea : sp+ss]
	for i := range vv ***REMOVED***
		vv[i] = nil
	***REMOVED***
	vm.sp = sp + ss
	vm.pc++
***REMOVED***

type _ret struct***REMOVED******REMOVED***

var ret _ret

func (_ret) exec(vm *vm) ***REMOVED***
	// callee -3
	// this -2 <- sb
	// retval -1

	vm.stack[vm.sb-1] = vm.stack[vm.sp-1]
	vm.sp = vm.sb
	vm.popCtx()
	if vm.pc < 0 ***REMOVED***
		vm.halt = true
	***REMOVED***
***REMOVED***

type enterFuncStashless struct ***REMOVED***
	stackSize uint32
	args      uint32
***REMOVED***

func (e *enterFuncStashless) exec(vm *vm) ***REMOVED***
	sp := vm.sp
	vm.sb = sp - vm.args - 1
	d := int(e.args) - vm.args
	if d > 0 ***REMOVED***
		ss := sp + int(e.stackSize) + d
		vm.stack.expand(ss)
		vv := vm.stack[sp : sp+d]
		for i := range vv ***REMOVED***
			vv[i] = _undefined
		***REMOVED***
		vv = vm.stack[sp+d : ss]
		for i := range vv ***REMOVED***
			vv[i] = nil
		***REMOVED***
		vm.args = int(e.args)
		vm.sp = ss
	***REMOVED*** else ***REMOVED***
		if e.stackSize > 0 ***REMOVED***
			ss := sp + int(e.stackSize)
			vm.stack.expand(ss)
			vv := vm.stack[sp:ss]
			for i := range vv ***REMOVED***
				vv[i] = nil
			***REMOVED***
			vm.sp = ss
		***REMOVED***
	***REMOVED***
	vm.pc++
***REMOVED***

type newFunc struct ***REMOVED***
	prg    *Program
	name   unistring.String
	length uint32
	strict bool

	srcStart, srcEnd uint32
***REMOVED***

func (n *newFunc) exec(vm *vm) ***REMOVED***
	obj := vm.r.newFunc(n.name, int(n.length), n.strict)
	obj.prg = n.prg
	obj.stash = vm.stash
	obj.src = n.prg.src.Source()[n.srcStart:n.srcEnd]
	vm.push(obj.val)
	vm.pc++
***REMOVED***

func (vm *vm) alreadyDeclared(name unistring.String) Value ***REMOVED***
	return vm.r.newError(vm.r.global.SyntaxError, "Identifier '%s' has already been declared", name)
***REMOVED***

func (vm *vm) checkBindVarsGlobal(names []unistring.String) ***REMOVED***
	o := vm.r.globalObject.self
	sn := vm.r.global.stash.names
	if o, ok := o.(*baseObject); ok ***REMOVED***
		// shortcut
		for _, name := range names ***REMOVED***
			if !o.hasOwnPropertyStr(name) && !o.extensible ***REMOVED***
				panic(vm.r.NewTypeError("Cannot define global variable '%s', global object is not extensible", name))
			***REMOVED***
			if _, exists := sn[name]; exists ***REMOVED***
				panic(vm.alreadyDeclared(name))
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, name := range names ***REMOVED***
			if !o.hasOwnPropertyStr(name) && !o.isExtensible() ***REMOVED***
				panic(vm.r.NewTypeError("Cannot define global variable '%s', global object is not extensible", name))
			***REMOVED***
			if _, exists := sn[name]; exists ***REMOVED***
				panic(vm.alreadyDeclared(name))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (vm *vm) createGlobalVarBindings(names []unistring.String, d bool) ***REMOVED***
	globalVarNames := vm.r.global.varNames
	if globalVarNames == nil ***REMOVED***
		globalVarNames = make(map[unistring.String]struct***REMOVED******REMOVED***)
		vm.r.global.varNames = globalVarNames
	***REMOVED***
	o := vm.r.globalObject.self
	if o, ok := o.(*baseObject); ok ***REMOVED***
		for _, name := range names ***REMOVED***
			if !o.hasOwnPropertyStr(name) && o.extensible ***REMOVED***
				o._putProp(name, _undefined, true, true, d)
			***REMOVED***
			globalVarNames[name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var cf Flag
		if d ***REMOVED***
			cf = FLAG_TRUE
		***REMOVED*** else ***REMOVED***
			cf = FLAG_FALSE
		***REMOVED***
		for _, name := range names ***REMOVED***
			if !o.hasOwnPropertyStr(name) && o.isExtensible() ***REMOVED***
				o.defineOwnPropertyStr(name, PropertyDescriptor***REMOVED***
					Value:        _undefined,
					Writable:     FLAG_TRUE,
					Enumerable:   FLAG_TRUE,
					Configurable: cf,
				***REMOVED***, true)
				o.setOwnStr(name, _undefined, false)
			***REMOVED***
			globalVarNames[name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (vm *vm) createGlobalFuncBindings(names []unistring.String, d bool) ***REMOVED***
	globalVarNames := vm.r.global.varNames
	if globalVarNames == nil ***REMOVED***
		globalVarNames = make(map[unistring.String]struct***REMOVED******REMOVED***)
		vm.r.global.varNames = globalVarNames
	***REMOVED***
	o := vm.r.globalObject.self
	b := vm.sp - len(names)
	var shortcutObj *baseObject
	if o, ok := o.(*baseObject); ok ***REMOVED***
		shortcutObj = o
	***REMOVED***
	for i, name := range names ***REMOVED***
		var desc PropertyDescriptor
		prop := o.getOwnPropStr(name)
		desc.Value = vm.stack[b+i]
		if shortcutObj != nil && prop == nil && shortcutObj.extensible ***REMOVED***
			shortcutObj._putProp(name, desc.Value, true, true, d)
		***REMOVED*** else ***REMOVED***
			if prop, ok := prop.(*valueProperty); ok && !prop.configurable ***REMOVED***
				// no-op
			***REMOVED*** else ***REMOVED***
				desc.Writable = FLAG_TRUE
				desc.Enumerable = FLAG_TRUE
				if d ***REMOVED***
					desc.Configurable = FLAG_TRUE
				***REMOVED*** else ***REMOVED***
					desc.Configurable = FLAG_FALSE
				***REMOVED***
			***REMOVED***
			if shortcutObj != nil ***REMOVED***
				shortcutObj.defineOwnPropertyStr(name, desc, true)
			***REMOVED*** else ***REMOVED***
				o.defineOwnPropertyStr(name, desc, true)
				o.setOwnStr(name, desc.Value, false) // not a bug, see https://262.ecma-international.org/#sec-createglobalfunctionbinding
			***REMOVED***
		***REMOVED***
		globalVarNames[name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	vm.sp = b
***REMOVED***

func (vm *vm) checkBindFuncsGlobal(names []unistring.String) ***REMOVED***
	o := vm.r.globalObject.self
	sn := vm.r.global.stash.names
	for _, name := range names ***REMOVED***
		if _, exists := sn[name]; exists ***REMOVED***
			panic(vm.alreadyDeclared(name))
		***REMOVED***
		prop := o.getOwnPropStr(name)
		allowed := true
		switch prop := prop.(type) ***REMOVED***
		case nil:
			allowed = o.isExtensible()
		case *valueProperty:
			allowed = prop.configurable || prop.getterFunc == nil && prop.setterFunc == nil && prop.writable && prop.enumerable
		***REMOVED***
		if !allowed ***REMOVED***
			panic(vm.r.NewTypeError("Cannot redefine global function '%s'", name))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (vm *vm) checkBindLexGlobal(names []unistring.String) ***REMOVED***
	o := vm.r.globalObject.self
	s := &vm.r.global.stash
	for _, name := range names ***REMOVED***
		if _, exists := vm.r.global.varNames[name]; exists ***REMOVED***
			goto fail
		***REMOVED***
		if _, exists := s.names[name]; exists ***REMOVED***
			goto fail
		***REMOVED***
		if prop, ok := o.getOwnPropStr(name).(*valueProperty); ok && !prop.configurable ***REMOVED***
			goto fail
		***REMOVED***
		continue
	fail:
		panic(vm.alreadyDeclared(name))
	***REMOVED***
***REMOVED***

type bindVars struct ***REMOVED***
	names     []unistring.String
	deletable bool
***REMOVED***

func (d *bindVars) exec(vm *vm) ***REMOVED***
	var target *stash
	for _, name := range d.names ***REMOVED***
		for s := vm.stash; s != nil; s = s.outer ***REMOVED***
			if idx, exists := s.names[name]; exists && idx&maskVar == 0 ***REMOVED***
				panic(vm.alreadyDeclared(name))
			***REMOVED***
			if s.function ***REMOVED***
				target = s
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if target == nil ***REMOVED***
		target = vm.stash
	***REMOVED***
	deletable := d.deletable
	for _, name := range d.names ***REMOVED***
		target.createBinding(name, deletable)
	***REMOVED***
	vm.pc++
***REMOVED***

type bindGlobal struct ***REMOVED***
	vars, funcs, lets, consts []unistring.String

	deletable bool
***REMOVED***

func (b *bindGlobal) exec(vm *vm) ***REMOVED***
	vm.checkBindFuncsGlobal(b.funcs)
	vm.checkBindLexGlobal(b.lets)
	vm.checkBindLexGlobal(b.consts)
	vm.checkBindVarsGlobal(b.vars)

	s := &vm.r.global.stash
	for _, name := range b.lets ***REMOVED***
		s.createLexBinding(name, false)
	***REMOVED***
	for _, name := range b.consts ***REMOVED***
		s.createLexBinding(name, true)
	***REMOVED***
	vm.createGlobalFuncBindings(b.funcs, b.deletable)
	vm.createGlobalVarBindings(b.vars, b.deletable)
	vm.pc++
***REMOVED***

type jne int32

func (j jne) exec(vm *vm) ***REMOVED***
	vm.sp--
	if !vm.stack[vm.sp].ToBoolean() ***REMOVED***
		vm.pc += int(j)
	***REMOVED*** else ***REMOVED***
		vm.pc++
	***REMOVED***
***REMOVED***

type jeq int32

func (j jeq) exec(vm *vm) ***REMOVED***
	vm.sp--
	if vm.stack[vm.sp].ToBoolean() ***REMOVED***
		vm.pc += int(j)
	***REMOVED*** else ***REMOVED***
		vm.pc++
	***REMOVED***
***REMOVED***

type jeq1 int32

func (j jeq1) exec(vm *vm) ***REMOVED***
	if vm.stack[vm.sp-1].ToBoolean() ***REMOVED***
		vm.pc += int(j)
	***REMOVED*** else ***REMOVED***
		vm.pc++
	***REMOVED***
***REMOVED***

type jneq1 int32

func (j jneq1) exec(vm *vm) ***REMOVED***
	if !vm.stack[vm.sp-1].ToBoolean() ***REMOVED***
		vm.pc += int(j)
	***REMOVED*** else ***REMOVED***
		vm.pc++
	***REMOVED***
***REMOVED***

type _not struct***REMOVED******REMOVED***

var not _not

func (_not) exec(vm *vm) ***REMOVED***
	if vm.stack[vm.sp-1].ToBoolean() ***REMOVED***
		vm.stack[vm.sp-1] = valueFalse
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-1] = valueTrue
	***REMOVED***
	vm.pc++
***REMOVED***

func toPrimitiveNumber(v Value) Value ***REMOVED***
	if o, ok := v.(*Object); ok ***REMOVED***
		return o.toPrimitiveNumber()
	***REMOVED***
	return v
***REMOVED***

func toPrimitive(v Value) Value ***REMOVED***
	if o, ok := v.(*Object); ok ***REMOVED***
		return o.toPrimitive()
	***REMOVED***
	return v
***REMOVED***

func cmp(px, py Value) Value ***REMOVED***
	var ret bool
	var nx, ny float64

	if xs, ok := px.(valueString); ok ***REMOVED***
		if ys, ok := py.(valueString); ok ***REMOVED***
			ret = xs.compareTo(ys) < 0
			goto end
		***REMOVED***
	***REMOVED***

	if xi, ok := px.(valueInt); ok ***REMOVED***
		if yi, ok := py.(valueInt); ok ***REMOVED***
			ret = xi < yi
			goto end
		***REMOVED***
	***REMOVED***

	nx = px.ToFloat()
	ny = py.ToFloat()

	if math.IsNaN(nx) || math.IsNaN(ny) ***REMOVED***
		return _undefined
	***REMOVED***

	ret = nx < ny

end:
	if ret ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse

***REMOVED***

type _op_lt struct***REMOVED******REMOVED***

var op_lt _op_lt

func (_op_lt) exec(vm *vm) ***REMOVED***
	left := toPrimitiveNumber(vm.stack[vm.sp-2])
	right := toPrimitiveNumber(vm.stack[vm.sp-1])

	r := cmp(left, right)
	if r == _undefined ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = r
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type _op_lte struct***REMOVED******REMOVED***

var op_lte _op_lte

func (_op_lte) exec(vm *vm) ***REMOVED***
	left := toPrimitiveNumber(vm.stack[vm.sp-2])
	right := toPrimitiveNumber(vm.stack[vm.sp-1])

	r := cmp(right, left)
	if r == _undefined || r == valueTrue ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED***

	vm.sp--
	vm.pc++
***REMOVED***

type _op_gt struct***REMOVED******REMOVED***

var op_gt _op_gt

func (_op_gt) exec(vm *vm) ***REMOVED***
	left := toPrimitiveNumber(vm.stack[vm.sp-2])
	right := toPrimitiveNumber(vm.stack[vm.sp-1])

	r := cmp(right, left)
	if r == _undefined ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = r
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type _op_gte struct***REMOVED******REMOVED***

var op_gte _op_gte

func (_op_gte) exec(vm *vm) ***REMOVED***
	left := toPrimitiveNumber(vm.stack[vm.sp-2])
	right := toPrimitiveNumber(vm.stack[vm.sp-1])

	r := cmp(left, right)
	if r == _undefined || r == valueTrue ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED***

	vm.sp--
	vm.pc++
***REMOVED***

type _op_eq struct***REMOVED******REMOVED***

var op_eq _op_eq

func (_op_eq) exec(vm *vm) ***REMOVED***
	if vm.stack[vm.sp-2].Equals(vm.stack[vm.sp-1]) ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type _op_neq struct***REMOVED******REMOVED***

var op_neq _op_neq

func (_op_neq) exec(vm *vm) ***REMOVED***
	if vm.stack[vm.sp-2].Equals(vm.stack[vm.sp-1]) ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type _op_strict_eq struct***REMOVED******REMOVED***

var op_strict_eq _op_strict_eq

func (_op_strict_eq) exec(vm *vm) ***REMOVED***
	if vm.stack[vm.sp-2].StrictEquals(vm.stack[vm.sp-1]) ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type _op_strict_neq struct***REMOVED******REMOVED***

var op_strict_neq _op_strict_neq

func (_op_strict_neq) exec(vm *vm) ***REMOVED***
	if vm.stack[vm.sp-2].StrictEquals(vm.stack[vm.sp-1]) ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type _op_instanceof struct***REMOVED******REMOVED***

var op_instanceof _op_instanceof

func (_op_instanceof) exec(vm *vm) ***REMOVED***
	left := vm.stack[vm.sp-2]
	right := vm.r.toObject(vm.stack[vm.sp-1])

	if instanceOfOperator(left, right) ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED***

	vm.sp--
	vm.pc++
***REMOVED***

type _op_in struct***REMOVED******REMOVED***

var op_in _op_in

func (_op_in) exec(vm *vm) ***REMOVED***
	left := vm.stack[vm.sp-2]
	right := vm.r.toObject(vm.stack[vm.sp-1])

	if right.hasProperty(left) ***REMOVED***
		vm.stack[vm.sp-2] = valueTrue
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-2] = valueFalse
	***REMOVED***

	vm.sp--
	vm.pc++
***REMOVED***

type try struct ***REMOVED***
	catchOffset   int32
	finallyOffset int32
***REMOVED***

func (t try) exec(vm *vm) ***REMOVED***
	o := vm.pc
	vm.pc++
	ex := vm.runTry()
	if ex != nil && t.catchOffset > 0 ***REMOVED***
		// run the catch block (in try)
		vm.pc = o + int(t.catchOffset)
		// TODO: if ex.val is an Error, set the stack property
		vm.push(ex.val)
		ex = vm.runTry()
	***REMOVED***

	if t.finallyOffset > 0 ***REMOVED***
		pc := vm.pc
		// Run finally
		vm.pc = o + int(t.finallyOffset)
		vm.run()
		if vm.prg.code[vm.pc] == retFinally ***REMOVED***
			vm.pc = pc
		***REMOVED*** else ***REMOVED***
			// break or continue out of finally, dropping exception
			ex = nil
		***REMOVED***
	***REMOVED***

	vm.halt = false

	if ex != nil ***REMOVED***
		vm.pc = -1 // to prevent the current position from being captured in the stacktrace
		panic(ex)
	***REMOVED***
***REMOVED***

type _retFinally struct***REMOVED******REMOVED***

var retFinally _retFinally

func (_retFinally) exec(vm *vm) ***REMOVED***
	vm.pc++
***REMOVED***

type _throw struct***REMOVED******REMOVED***

var throw _throw

func (_throw) exec(vm *vm) ***REMOVED***
	panic(vm.stack[vm.sp-1])
***REMOVED***

type _new uint32

func (n _new) exec(vm *vm) ***REMOVED***
	sp := vm.sp - int(n)
	obj := vm.stack[sp-1]
	ctor := vm.r.toConstructor(obj)
	vm.stack[sp-1] = ctor(vm.stack[sp:vm.sp], nil)
	vm.sp = sp
	vm.pc++
***REMOVED***

type _loadNewTarget struct***REMOVED******REMOVED***

var loadNewTarget _loadNewTarget

func (_loadNewTarget) exec(vm *vm) ***REMOVED***
	if t := vm.newTarget; t != nil ***REMOVED***
		vm.push(t)
	***REMOVED*** else ***REMOVED***
		vm.push(_undefined)
	***REMOVED***
	vm.pc++
***REMOVED***

type _typeof struct***REMOVED******REMOVED***

var typeof _typeof

func (_typeof) exec(vm *vm) ***REMOVED***
	var r Value
	switch v := vm.stack[vm.sp-1].(type) ***REMOVED***
	case valueUndefined, valueUnresolved:
		r = stringUndefined
	case valueNull:
		r = stringObjectC
	case *Object:
	repeat:
		switch s := v.self.(type) ***REMOVED***
		case *funcObject, *nativeFuncObject, *boundFuncObject:
			r = stringFunction
		case *lazyObject:
			v.self = s.create(v)
			goto repeat
		default:
			r = stringObjectC
		***REMOVED***
	case valueBool:
		r = stringBoolean
	case valueString:
		r = stringString
	case valueInt, valueFloat:
		r = stringNumber
	case *Symbol:
		r = stringSymbol
	default:
		panic(fmt.Errorf("Unknown type: %T", v))
	***REMOVED***
	vm.stack[vm.sp-1] = r
	vm.pc++
***REMOVED***

type createArgs uint32

func (formalArgs createArgs) exec(vm *vm) ***REMOVED***
	v := &Object***REMOVED***runtime: vm.r***REMOVED***
	args := &argumentsObject***REMOVED******REMOVED***
	args.extensible = true
	args.prototype = vm.r.global.ObjectPrototype
	args.class = "Arguments"
	v.self = args
	args.val = v
	args.length = vm.args
	args.init()
	i := 0
	c := int(formalArgs)
	if vm.args < c ***REMOVED***
		c = vm.args
	***REMOVED***
	for ; i < c; i++ ***REMOVED***
		args._put(unistring.String(strconv.Itoa(i)), &mappedProperty***REMOVED***
			valueProperty: valueProperty***REMOVED***
				writable:     true,
				configurable: true,
				enumerable:   true,
			***REMOVED***,
			v: &vm.stash.values[i],
		***REMOVED***)
	***REMOVED***

	for _, v := range vm.stash.extraArgs ***REMOVED***
		args._put(unistring.String(strconv.Itoa(i)), v)
		i++
	***REMOVED***

	args._putProp("callee", vm.stack[vm.sb-1], true, false, true)
	args._putSym(SymIterator, valueProp(vm.r.global.arrayValues, true, false, true))
	vm.push(v)
	vm.pc++
***REMOVED***

type createArgsStrict uint32

func (formalArgs createArgsStrict) exec(vm *vm) ***REMOVED***
	args := vm.r.newBaseObject(vm.r.global.ObjectPrototype, "Arguments")
	i := 0
	c := int(formalArgs)
	if vm.args < c ***REMOVED***
		c = vm.args
	***REMOVED***
	for _, v := range vm.stash.values[:c] ***REMOVED***
		args._put(unistring.String(strconv.Itoa(i)), v)
		i++
	***REMOVED***

	for _, v := range vm.stash.extraArgs ***REMOVED***
		args._put(unistring.String(strconv.Itoa(i)), v)
		i++
	***REMOVED***

	args._putProp("length", intToValue(int64(vm.args)), true, false, true)
	args._put("callee", vm.r.global.throwerProperty)
	args._put("caller", vm.r.global.throwerProperty)
	args._putSym(SymIterator, valueProp(vm.r.global.arrayValues, true, false, true))
	vm.push(args.val)
	vm.pc++
***REMOVED***

type _enterWith struct***REMOVED******REMOVED***

var enterWith _enterWith

func (_enterWith) exec(vm *vm) ***REMOVED***
	vm.newStash()
	vm.stash.obj = vm.stack[vm.sp-1].ToObject(vm.r)
	vm.sp--
	vm.pc++
***REMOVED***

type _leaveWith struct***REMOVED******REMOVED***

var leaveWith _leaveWith

func (_leaveWith) exec(vm *vm) ***REMOVED***
	vm.stash = vm.stash.outer
	vm.pc++
***REMOVED***

func emptyIter() (propIterItem, iterNextFunc) ***REMOVED***
	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

type _enumerate struct***REMOVED******REMOVED***

var enumerate _enumerate

func (_enumerate) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-1]
	if v == _undefined || v == _null ***REMOVED***
		vm.iterStack = append(vm.iterStack, iterStackItem***REMOVED***f: emptyIter***REMOVED***)
	***REMOVED*** else ***REMOVED***
		vm.iterStack = append(vm.iterStack, iterStackItem***REMOVED***f: enumerateRecursive(v.ToObject(vm.r))***REMOVED***)
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type enumNext int32

func (jmp enumNext) exec(vm *vm) ***REMOVED***
	l := len(vm.iterStack) - 1
	item, n := vm.iterStack[l].f()
	if n != nil ***REMOVED***
		vm.iterStack[l].val = stringValueFromRaw(item.name)
		vm.iterStack[l].f = n
		vm.pc++
	***REMOVED*** else ***REMOVED***
		vm.pc += int(jmp)
	***REMOVED***
***REMOVED***

type _enumGet struct***REMOVED******REMOVED***

var enumGet _enumGet

func (_enumGet) exec(vm *vm) ***REMOVED***
	l := len(vm.iterStack) - 1
	vm.push(vm.iterStack[l].val)
	vm.pc++
***REMOVED***

type _enumPop struct***REMOVED******REMOVED***

var enumPop _enumPop

func (_enumPop) exec(vm *vm) ***REMOVED***
	l := len(vm.iterStack) - 1
	vm.iterStack[l] = iterStackItem***REMOVED******REMOVED***
	vm.iterStack = vm.iterStack[:l]
	vm.pc++
***REMOVED***

type _enumPopClose struct***REMOVED******REMOVED***

var enumPopClose _enumPopClose

func (_enumPopClose) exec(vm *vm) ***REMOVED***
	l := len(vm.iterStack) - 1
	item := vm.iterStack[l]
	vm.iterStack[l] = iterStackItem***REMOVED******REMOVED***
	vm.iterStack = vm.iterStack[:l]
	if iter := item.iter; iter != nil ***REMOVED***
		returnIter(iter)
	***REMOVED***
	vm.pc++
***REMOVED***

type _iterate struct***REMOVED******REMOVED***

var iterate _iterate

func (_iterate) exec(vm *vm) ***REMOVED***
	iter := vm.r.getIterator(vm.stack[vm.sp-1], nil)
	vm.iterStack = append(vm.iterStack, iterStackItem***REMOVED***iter: iter***REMOVED***)
	vm.sp--
	vm.pc++
***REMOVED***

type iterNext int32

func (jmp iterNext) exec(vm *vm) ***REMOVED***
	l := len(vm.iterStack) - 1
	iter := vm.iterStack[l].iter
	var res *Object
	var done bool
	var value Value
	ex := vm.try(func() ***REMOVED***
		res = vm.r.toObject(toMethod(iter.self.getStr("next", nil))(FunctionCall***REMOVED***This: iter***REMOVED***))
		done = nilSafe(res.self.getStr("done", nil)).ToBoolean()
		if !done ***REMOVED***
			value = nilSafe(res.self.getStr("value", nil))
			vm.iterStack[l].val = value
		***REMOVED***
	***REMOVED***)
	if ex == nil ***REMOVED***
		if done ***REMOVED***
			vm.pc += int(jmp)
		***REMOVED*** else ***REMOVED***
			vm.iterStack[l].val = value
			vm.pc++
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		l := len(vm.iterStack) - 1
		vm.iterStack[l] = iterStackItem***REMOVED******REMOVED***
		vm.iterStack = vm.iterStack[:l]
		panic(ex.val)
	***REMOVED***
***REMOVED***

type copyStash struct***REMOVED******REMOVED***

func (copyStash) exec(vm *vm) ***REMOVED***
	oldStash := vm.stash
	newStash := &stash***REMOVED***
		outer: oldStash.outer,
	***REMOVED***
	vm.stashAllocs++
	newStash.values = append([]Value(nil), oldStash.values...)
	vm.stash = newStash
	vm.pc++
***REMOVED***

type _throwAssignToConst struct***REMOVED******REMOVED***

var throwAssignToConst _throwAssignToConst

func (_throwAssignToConst) exec(vm *vm) ***REMOVED***
	panic(errAssignToConst)
***REMOVED***
