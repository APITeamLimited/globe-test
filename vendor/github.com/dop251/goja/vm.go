package goja

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	maxInt = 1 << 53
)

type valueStack []Value

type stash struct ***REMOVED***
	values    valueStack
	extraArgs valueStack
	names     map[string]uint32
	obj       objectImpl

	outer *stash
***REMOVED***

type context struct ***REMOVED***
	prg      *Program
	funcName string
	stash    *stash
	pc, sb   int
	args     int
***REMOVED***

type iterStackItem struct ***REMOVED***
	val Value
	f   iterNextFunc
***REMOVED***

type ref interface ***REMOVED***
	get() Value
	set(Value)
	refname() string
***REMOVED***

type stashRef struct ***REMOVED***
	v *Value
	n string
***REMOVED***

func (r stashRef) get() Value ***REMOVED***
	return *r.v
***REMOVED***

func (r *stashRef) set(v Value) ***REMOVED***
	*r.v = v
***REMOVED***

func (r *stashRef) refname() string ***REMOVED***
	return r.n
***REMOVED***

type objRef struct ***REMOVED***
	base   objectImpl
	name   string
	strict bool
***REMOVED***

func (r *objRef) get() Value ***REMOVED***
	return r.base.getStr(r.name)
***REMOVED***

func (r *objRef) set(v Value) ***REMOVED***
	r.base.putStr(r.name, v, r.strict)
***REMOVED***

func (r *objRef) refname() string ***REMOVED***
	return r.name
***REMOVED***

type unresolvedRef struct ***REMOVED***
	runtime *Runtime
	name    string
***REMOVED***

func (r *unresolvedRef) get() Value ***REMOVED***
	r.runtime.throwReferenceError(r.name)
	panic("Unreachable")
***REMOVED***

func (r *unresolvedRef) set(Value) ***REMOVED***
	r.get()
***REMOVED***

func (r *unresolvedRef) refname() string ***REMOVED***
	return r.name
***REMOVED***

type vm struct ***REMOVED***
	r            *Runtime
	prg          *Program
	funcName     string
	pc           int
	stack        valueStack
	sp, sb, args int

	stash     *stash
	callStack []context
	iterStack []iterStackItem
	refStack  []ref

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
	return valueFloat(float64(i))
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

func toInt(v Value) (int64, bool) ***REMOVED***
	num := v.ToNumber()
	if i, ok := num.assertInt(); ok ***REMOVED***
		return i, true
	***REMOVED***
	if f, ok := num.assertFloat(); ok ***REMOVED***
		if i, ok := floatToInt(f); ok ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

func toIntIgnoreNegZero(v Value) (int64, bool) ***REMOVED***
	num := v.ToNumber()
	if i, ok := num.assertInt(); ok ***REMOVED***
		return i, true
	***REMOVED***
	if f, ok := num.assertFloat(); ok ***REMOVED***
		if v == _negativeZero ***REMOVED***
			return 0, true
		***REMOVED***
		if i, ok := floatToInt(f); ok ***REMOVED***
			return i, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

func (s *valueStack) expand(idx int) ***REMOVED***
	if idx < len(*s) ***REMOVED***
		return
	***REMOVED***

	if idx < cap(*s) ***REMOVED***
		*s = (*s)[:idx+1]
	***REMOVED*** else ***REMOVED***
		n := make([]Value, idx+1, (idx+1)<<1)
		copy(n, *s)
		*s = n
	***REMOVED***
***REMOVED***

func (s *stash) put(name string, v Value) bool ***REMOVED***
	if s.obj != nil ***REMOVED***
		if found := s.obj.getStr(name); found != nil ***REMOVED***
			s.obj.putStr(name, v, false)
			return true
		***REMOVED***
		return false
	***REMOVED*** else ***REMOVED***
		if idx, found := s.names[name]; found ***REMOVED***
			s.values.expand(int(idx))
			s.values[idx] = v
			return true
		***REMOVED***
		return false
	***REMOVED***
***REMOVED***

func (s *stash) putByIdx(idx uint32, v Value) ***REMOVED***
	if s.obj != nil ***REMOVED***
		panic("Attempt to put by idx into an object scope")
	***REMOVED***
	s.values.expand(int(idx))
	s.values[idx] = v
***REMOVED***

func (s *stash) getByIdx(idx uint32) Value ***REMOVED***
	if int(idx) < len(s.values) ***REMOVED***
		return s.values[idx]
	***REMOVED***
	return _undefined
***REMOVED***

func (s *stash) getByName(name string, _ *vm) (v Value, exists bool) ***REMOVED***
	if s.obj != nil ***REMOVED***
		v = s.obj.getStr(name)
		if v == nil ***REMOVED***
			return nil, false
			//return valueUnresolved***REMOVED***r: vm.r, ref: name***REMOVED***, false
		***REMOVED***
		return v, true
	***REMOVED***
	if idx, exists := s.names[name]; exists ***REMOVED***
		return s.values[idx], true
	***REMOVED***
	return nil, false
	//return valueUnresolved***REMOVED***r: vm.r, ref: name***REMOVED***, false
***REMOVED***

func (s *stash) createBinding(name string) ***REMOVED***
	if s.names == nil ***REMOVED***
		s.names = make(map[string]uint32)
	***REMOVED***
	if _, exists := s.names[name]; !exists ***REMOVED***
		s.names[name] = uint32(len(s.names))
		s.values = append(s.values, _undefined)
	***REMOVED***
***REMOVED***

func (s *stash) deleteBinding(name string) bool ***REMOVED***
	if s.obj != nil ***REMOVED***
		return s.obj.deleteStr(name, false)
	***REMOVED***
	if idx, found := s.names[name]; found ***REMOVED***
		s.values[idx] = nil
		delete(s.names, name)
		return true
	***REMOVED***
	return false
***REMOVED***

func (vm *vm) newStash() ***REMOVED***
	vm.stash = &stash***REMOVED***
		outer: vm.stash,
	***REMOVED***
	vm.stashAllocs++
***REMOVED***

func (vm *vm) init() ***REMOVED***
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
		panic(v)
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

func (vm *vm) captureStack(stack []stackFrame, ctxOffset int) []stackFrame ***REMOVED***
	// Unroll the context stack
	stack = append(stack, stackFrame***REMOVED***prg: vm.prg, pc: vm.pc, funcName: vm.funcName***REMOVED***)
	for i := len(vm.callStack) - 1; i > ctxOffset-1; i-- ***REMOVED***
		if vm.callStack[i].pc != -1 ***REMOVED***
			stack = append(stack, stackFrame***REMOVED***prg: vm.callStack[i].prg, pc: vm.callStack[i].pc - 1, funcName: vm.callStack[i].funcName***REMOVED***)
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
			case *InterruptedError:
				x1.stack = vm.captureStack(x1.stack, ctxOffset)
				panic(x1)
			case *Exception:
				ex = x1
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
	ctx.prg = vm.prg
	ctx.funcName = vm.funcName
	ctx.stash = vm.stash
	ctx.pc = vm.pc
	ctx.sb = vm.sb
	ctx.args = vm.args
***REMOVED***

func (vm *vm) pushCtx() ***REMOVED***
	/*
		vm.ctxStack = append(vm.ctxStack, context***REMOVED***
			prg: vm.prg,
			stash: vm.stash,
			pc: vm.pc,
			sb: vm.sb,
			args: vm.args,
		***REMOVED***)*/
	vm.callStack = append(vm.callStack, context***REMOVED******REMOVED***)
	vm.saveCtx(&vm.callStack[len(vm.callStack)-1])
***REMOVED***

func (vm *vm) restoreCtx(ctx *context) ***REMOVED***
	vm.prg = ctx.prg
	vm.funcName = ctx.funcName
	vm.pc = ctx.pc
	vm.stash = ctx.stash
	vm.sb = ctx.sb
	vm.args = ctx.args
***REMOVED***

func (vm *vm) popCtx() ***REMOVED***
	l := len(vm.callStack) - 1
	vm.prg = vm.callStack[l].prg
	vm.callStack[l].prg = nil
	vm.funcName = vm.callStack[l].funcName
	vm.pc = vm.callStack[l].pc
	vm.stash = vm.callStack[l].stash
	vm.callStack[l].stash = nil
	vm.sb = vm.callStack[l].sb
	vm.args = vm.callStack[l].args

	vm.callStack = vm.callStack[:l]
***REMOVED***

func (r *Runtime) toObject(v Value, args ...interface***REMOVED******REMOVED***) *Object ***REMOVED***
	//r.checkResolveable(v)
	if obj, ok := v.(*Object); ok ***REMOVED***
		return obj
	***REMOVED***
	if len(args) > 0 ***REMOVED***
		panic(r.NewTypeError(args...))
	***REMOVED*** else ***REMOVED***
		panic(r.NewTypeError("Value is not an object: %s", v.String()))
	***REMOVED***
***REMOVED***

func (r *Runtime) toCallee(v Value) *Object ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		return obj
	***REMOVED***
	switch unresolved := v.(type) ***REMOVED***
	case valueUnresolved:
		unresolved.throw()
		panic("Unreachable")
	case memberUnresolved:
		r.typeErrorResult(true, "Object has no member '%s'", unresolved.ref)
		panic("Unreachable")
	***REMOVED***
	r.typeErrorResult(true, "Value is not an object: %s", v.ToString())
	panic("Unreachable")
***REMOVED***

type _newStash struct***REMOVED******REMOVED***

var newStash _newStash

func (_newStash) exec(vm *vm) ***REMOVED***
	vm.newStash()
	vm.pc++
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

type _loadGlobalObject struct***REMOVED******REMOVED***

var loadGlobalObject _loadGlobalObject

func (_loadGlobalObject) exec(vm *vm) ***REMOVED***
	vm.push(vm.r.globalObject)
	vm.pc++
***REMOVED***

type loadStack int

func (l loadStack) exec(vm *vm) ***REMOVED***
	// l < 0 -- arg<-l-1>
	// l > 0 -- var<l-1>
	// l == 0 -- this

	if l < 0 ***REMOVED***
		arg := int(-l)
		if arg > vm.args ***REMOVED***
			vm.push(_undefined)
		***REMOVED*** else ***REMOVED***
			vm.push(vm.stack[vm.sb+arg])
		***REMOVED***
	***REMOVED*** else if l > 0 ***REMOVED***
		vm.push(vm.stack[vm.sb+vm.args+int(l)])
	***REMOVED*** else ***REMOVED***
		vm.push(vm.stack[vm.sb])
	***REMOVED***
	vm.pc++
***REMOVED***

type _loadCallee struct***REMOVED******REMOVED***

var loadCallee _loadCallee

func (_loadCallee) exec(vm *vm) ***REMOVED***
	vm.push(vm.stack[vm.sb-1])
	vm.pc++
***REMOVED***

func (vm *vm) storeStack(s int) ***REMOVED***
	// l < 0 -- arg<-l-1>
	// l > 0 -- var<l-1>
	// l == 0 -- this

	if s < 0 ***REMOVED***
		vm.stack[vm.sb-s] = vm.stack[vm.sp-1]
	***REMOVED*** else if s > 0 ***REMOVED***
		vm.stack[vm.sb+vm.args+s] = vm.stack[vm.sp-1]
	***REMOVED*** else ***REMOVED***
		panic("Attempt to modify this")
	***REMOVED***
	vm.pc++
***REMOVED***

type storeStack int

func (s storeStack) exec(vm *vm) ***REMOVED***
	vm.storeStack(int(s))
***REMOVED***

type storeStackP int

func (s storeStackP) exec(vm *vm) ***REMOVED***
	vm.storeStack(int(s))
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
		left = o.self.toPrimitive()
	***REMOVED***

	if o, ok := right.(*Object); ok ***REMOVED***
		right = o.self.toPrimitive()
	***REMOVED***

	var ret Value

	leftString, isLeftString := left.assertString()
	rightString, isRightString := right.assertString()

	if isLeftString || isRightString ***REMOVED***
		if !isLeftString ***REMOVED***
			leftString = left.ToString()
		***REMOVED***
		if !isRightString ***REMOVED***
			rightString = right.ToString()
		***REMOVED***
		ret = leftString.concat(rightString)
	***REMOVED*** else ***REMOVED***
		if leftInt, ok := left.assertInt(); ok ***REMOVED***
			if rightInt, ok := right.assertInt(); ok ***REMOVED***
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

	if left, ok := left.assertInt(); ok ***REMOVED***
		if right, ok := right.assertInt(); ok ***REMOVED***
			result = intToValue(left - right)
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

	if left, ok := toInt(left); ok ***REMOVED***
		if right, ok := toInt(right); ok ***REMOVED***
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

	if leftInt, ok := toInt(left); ok ***REMOVED***
		if rightInt, ok := toInt(right); ok ***REMOVED***
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

	if i, ok := toInt(operand); ok ***REMOVED***
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

	if i, ok := toInt(v); ok ***REMOVED***
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

	if i, ok := toInt(v); ok ***REMOVED***
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
	right := toUInt32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left << (right & 0x1F)))
	vm.sp--
	vm.pc++
***REMOVED***

type _sar struct***REMOVED******REMOVED***

var sar _sar

func (_sar) exec(vm *vm) ***REMOVED***
	left := toInt32(vm.stack[vm.sp-2])
	right := toUInt32(vm.stack[vm.sp-1])
	vm.stack[vm.sp-2] = intToValue(int64(left >> (right & 0x1F)))
	vm.sp--
	vm.pc++
***REMOVED***

type _shr struct***REMOVED******REMOVED***

var shr _shr

func (_shr) exec(vm *vm) ***REMOVED***
	left := toUInt32(vm.stack[vm.sp-2])
	right := toUInt32(vm.stack[vm.sp-1])
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
	propName := vm.stack[vm.sp-2]
	val := vm.stack[vm.sp-1]

	obj.self.put(propName, val, false)

	vm.sp -= 2
	vm.stack[vm.sp-1] = val
	vm.pc++
***REMOVED***

type _setElemStrict struct***REMOVED******REMOVED***

var setElemStrict _setElemStrict

func (_setElemStrict) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-3])
	propName := vm.stack[vm.sp-2]
	val := vm.stack[vm.sp-1]

	obj.self.put(propName, val, true)

	vm.sp -= 2
	vm.stack[vm.sp-1] = val
	vm.pc++
***REMOVED***

type _deleteElem struct***REMOVED******REMOVED***

var deleteElem _deleteElem

func (_deleteElem) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	propName := vm.stack[vm.sp-1]
	if !obj.self.hasProperty(propName) || obj.self.delete(propName, false) ***REMOVED***
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
	propName := vm.stack[vm.sp-1]
	obj.self.delete(propName, true)
	vm.stack[vm.sp-2] = valueTrue
	vm.sp--
	vm.pc++
***REMOVED***

type deleteProp string

func (d deleteProp) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-1])
	if !obj.self.hasPropertyStr(string(d)) || obj.self.deleteStr(string(d), false) ***REMOVED***
		vm.stack[vm.sp-1] = valueTrue
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp-1] = valueFalse
	***REMOVED***
	vm.pc++
***REMOVED***

type deletePropStrict string

func (d deletePropStrict) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-1])
	obj.self.deleteStr(string(d), true)
	vm.stack[vm.sp-1] = valueTrue
	vm.pc++
***REMOVED***

type setProp string

func (p setProp) exec(vm *vm) ***REMOVED***
	val := vm.stack[vm.sp-1]
	vm.stack[vm.sp-2].ToObject(vm.r).self.putStr(string(p), val, false)
	vm.stack[vm.sp-2] = val
	vm.sp--
	vm.pc++
***REMOVED***

type setPropStrict string

func (p setPropStrict) exec(vm *vm) ***REMOVED***
	obj := vm.stack[vm.sp-2]
	val := vm.stack[vm.sp-1]

	obj1 := vm.r.toObject(obj)
	obj1.self.putStr(string(p), val, true)
	vm.stack[vm.sp-2] = val
	vm.sp--
	vm.pc++
***REMOVED***

type setProp1 string

func (p setProp1) exec(vm *vm) ***REMOVED***
	vm.r.toObject(vm.stack[vm.sp-2]).self._putProp(string(p), vm.stack[vm.sp-1], true, true, true)

	vm.sp--
	vm.pc++
***REMOVED***

type _setProto struct***REMOVED******REMOVED***

var setProto _setProto

func (_setProto) exec(vm *vm) ***REMOVED***
	vm.r.toObject(vm.stack[vm.sp-2]).self.putStr("__proto__", vm.stack[vm.sp-1], true)

	vm.sp--
	vm.pc++
***REMOVED***

type setPropGetter string

func (s setPropGetter) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	descr := propertyDescr***REMOVED***
		Getter:       val,
		Configurable: FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
	***REMOVED***

	obj.self.defineOwnProperty(newStringValue(string(s)), descr, false)

	vm.sp--
	vm.pc++
***REMOVED***

type setPropSetter string

func (s setPropSetter) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-2])
	val := vm.stack[vm.sp-1]

	descr := propertyDescr***REMOVED***
		Setter:       val,
		Configurable: FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
	***REMOVED***

	obj.self.defineOwnProperty(newStringValue(string(s)), descr, false)

	vm.sp--
	vm.pc++
***REMOVED***

type getProp string

func (g getProp) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-1]
	obj := v.baseObject(vm.r)
	if obj == nil ***REMOVED***
		panic(vm.r.NewTypeError("Cannot read property '%s' of undefined", g))
	***REMOVED***
	prop := obj.self.getPropStr(string(g))
	if prop1, ok := prop.(*valueProperty); ok ***REMOVED***
		vm.stack[vm.sp-1] = prop1.get(v)
	***REMOVED*** else ***REMOVED***
		if prop == nil ***REMOVED***
			prop = _undefined
		***REMOVED***
		vm.stack[vm.sp-1] = prop
	***REMOVED***

	vm.pc++
***REMOVED***

type getPropCallee string

func (g getPropCallee) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-1]
	obj := v.baseObject(vm.r)
	if obj == nil ***REMOVED***
		panic(vm.r.NewTypeError("Cannot read property '%s' of undefined", g))
	***REMOVED***
	prop := obj.self.getPropStr(string(g))
	if prop1, ok := prop.(*valueProperty); ok ***REMOVED***
		vm.stack[vm.sp-1] = prop1.get(v)
	***REMOVED*** else ***REMOVED***
		if prop == nil ***REMOVED***
			prop = memberUnresolved***REMOVED***valueUnresolved***REMOVED***r: vm.r, ref: string(g)***REMOVED******REMOVED***
		***REMOVED***
		vm.stack[vm.sp-1] = prop
	***REMOVED***

	vm.pc++
***REMOVED***

type _getElem struct***REMOVED******REMOVED***

var getElem _getElem

func (_getElem) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-2]
	obj := v.baseObject(vm.r)
	propName := vm.stack[vm.sp-1]
	if obj == nil ***REMOVED***
		panic(vm.r.NewTypeError("Cannot read property '%s' of undefined", propName.String()))
	***REMOVED***

	prop := obj.self.getProp(propName)
	if prop1, ok := prop.(*valueProperty); ok ***REMOVED***
		vm.stack[vm.sp-2] = prop1.get(v)
	***REMOVED*** else ***REMOVED***
		if prop == nil ***REMOVED***
			prop = _undefined
		***REMOVED***
		vm.stack[vm.sp-2] = prop
	***REMOVED***

	vm.sp--
	vm.pc++
***REMOVED***

type _getElemCallee struct***REMOVED******REMOVED***

var getElemCallee _getElemCallee

func (_getElemCallee) exec(vm *vm) ***REMOVED***
	v := vm.stack[vm.sp-2]
	obj := v.baseObject(vm.r)
	propName := vm.stack[vm.sp-1]
	if obj == nil ***REMOVED***
		vm.r.typeErrorResult(true, "Cannot read property '%s' of undefined", propName.String())
		panic("Unreachable")
	***REMOVED***

	prop := obj.self.getProp(propName)
	if prop1, ok := prop.(*valueProperty); ok ***REMOVED***
		vm.stack[vm.sp-2] = prop1.get(v)
	***REMOVED*** else ***REMOVED***
		if prop == nil ***REMOVED***
			prop = memberUnresolved***REMOVED***valueUnresolved***REMOVED***r: vm.r, ref: propName.String()***REMOVED******REMOVED***
		***REMOVED***
		vm.stack[vm.sp-2] = prop
	***REMOVED***

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

type newRegexp struct ***REMOVED***
	pattern regexpPattern
	src     valueString

	global, ignoreCase, multiline bool
***REMOVED***

func (n *newRegexp) exec(vm *vm) ***REMOVED***
	vm.push(vm.r.newRegExpp(n.pattern, n.src, n.global, n.ignoreCase, n.multiline, vm.r.global.RegExpPrototype))
	vm.pc++
***REMOVED***

func (vm *vm) setLocal(s int) ***REMOVED***
	v := vm.stack[vm.sp-1]
	level := s >> 24
	idx := uint32(s & 0x00FFFFFF)
	stash := vm.stash
	for i := 0; i < level; i++ ***REMOVED***
		stash = stash.outer
	***REMOVED***
	stash.putByIdx(idx, v)
	vm.pc++
***REMOVED***

type setLocal uint32

func (s setLocal) exec(vm *vm) ***REMOVED***
	vm.setLocal(int(s))
***REMOVED***

type setLocalP uint32

func (s setLocalP) exec(vm *vm) ***REMOVED***
	vm.setLocal(int(s))
	vm.sp--
***REMOVED***

type setVar struct ***REMOVED***
	name string
	idx  uint32
***REMOVED***

func (s setVar) exec(vm *vm) ***REMOVED***
	v := vm.peek()

	level := int(s.idx >> 24)
	idx := uint32(s.idx & 0x00FFFFFF)
	stash := vm.stash
	name := s.name
	for i := 0; i < level; i++ ***REMOVED***
		if stash.put(name, v) ***REMOVED***
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***

	if stash != nil ***REMOVED***
		stash.putByIdx(idx, v)
	***REMOVED*** else ***REMOVED***
		vm.r.globalObject.self.putStr(name, v, false)
	***REMOVED***

end:
	vm.pc++
***REMOVED***

type resolveVar1 string

func (s resolveVar1) exec(vm *vm) ***REMOVED***
	name := string(s)
	var ref ref
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if stash.obj != nil ***REMOVED***
			if stash.obj.hasPropertyStr(name) ***REMOVED***
				ref = &objRef***REMOVED***
					base: stash.obj,
					name: name,
				***REMOVED***
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if idx, exists := stash.names[name]; exists ***REMOVED***
				ref = &stashRef***REMOVED***
					v: &stash.values[idx],
				***REMOVED***
				goto end
			***REMOVED***
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

type deleteVar string

func (d deleteVar) exec(vm *vm) ***REMOVED***
	name := string(d)
	ret := true
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if stash.obj != nil ***REMOVED***
			if stash.obj.hasPropertyStr(name) ***REMOVED***
				ret = stash.obj.deleteStr(name, false)
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if _, exists := stash.names[name]; exists ***REMOVED***
				ret = false
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

type deleteGlobal string

func (d deleteGlobal) exec(vm *vm) ***REMOVED***
	name := string(d)
	var ret bool
	if vm.r.globalObject.self.hasPropertyStr(name) ***REMOVED***
		ret = vm.r.globalObject.self.deleteStr(name, false)
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

type resolveVar1Strict string

func (s resolveVar1Strict) exec(vm *vm) ***REMOVED***
	name := string(s)
	var ref ref
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if stash.obj != nil ***REMOVED***
			if stash.obj.hasPropertyStr(name) ***REMOVED***
				ref = &objRef***REMOVED***
					base:   stash.obj,
					name:   name,
					strict: true,
				***REMOVED***
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if idx, exists := stash.names[name]; exists ***REMOVED***
				ref = &stashRef***REMOVED***
					v: &stash.values[idx],
				***REMOVED***
				goto end
			***REMOVED***
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
		name:    string(s),
	***REMOVED***

end:
	vm.refStack = append(vm.refStack, ref)
	vm.pc++
***REMOVED***

type setGlobal string

func (s setGlobal) exec(vm *vm) ***REMOVED***
	v := vm.peek()

	vm.r.globalObject.self.putStr(string(s), v, false)
	vm.pc++
***REMOVED***

type setGlobalStrict string

func (s setGlobalStrict) exec(vm *vm) ***REMOVED***
	v := vm.peek()

	name := string(s)
	o := vm.r.globalObject.self
	if o.hasOwnPropertyStr(name) ***REMOVED***
		o.putStr(name, v, true)
	***REMOVED*** else ***REMOVED***
		vm.r.throwReferenceError(name)
	***REMOVED***
	vm.pc++
***REMOVED***

type getLocal uint32

func (g getLocal) exec(vm *vm) ***REMOVED***
	level := int(g >> 24)
	idx := uint32(g & 0x00FFFFFF)
	stash := vm.stash
	for i := 0; i < level; i++ ***REMOVED***
		stash = stash.outer
	***REMOVED***

	vm.push(stash.getByIdx(idx))
	vm.pc++
***REMOVED***

type getVar struct ***REMOVED***
	name string
	idx  uint32
	ref  bool
***REMOVED***

func (g getVar) exec(vm *vm) ***REMOVED***
	level := int(g.idx >> 24)
	idx := uint32(g.idx & 0x00FFFFFF)
	stash := vm.stash
	name := g.name
	for i := 0; i < level; i++ ***REMOVED***
		if v, found := stash.getByName(name, vm); found ***REMOVED***
			vm.push(v)
			goto end
		***REMOVED***
		stash = stash.outer
	***REMOVED***
	if stash != nil ***REMOVED***
		vm.push(stash.getByIdx(idx))
	***REMOVED*** else ***REMOVED***
		v := vm.r.globalObject.self.getStr(name)
		if v == nil ***REMOVED***
			if g.ref ***REMOVED***
				v = valueUnresolved***REMOVED***r: vm.r, ref: name***REMOVED***
			***REMOVED*** else ***REMOVED***
				vm.r.throwReferenceError(name)
			***REMOVED***
		***REMOVED***
		vm.push(v)
	***REMOVED***
end:
	vm.pc++
***REMOVED***

type resolveVar struct ***REMOVED***
	name   string
	idx    uint32
	strict bool
***REMOVED***

func (r resolveVar) exec(vm *vm) ***REMOVED***
	level := int(r.idx >> 24)
	idx := uint32(r.idx & 0x00FFFFFF)
	stash := vm.stash
	var ref ref
	for i := 0; i < level; i++ ***REMOVED***
		if stash.obj != nil ***REMOVED***
			if stash.obj.hasPropertyStr(r.name) ***REMOVED***
				ref = &objRef***REMOVED***
					base:   stash.obj,
					name:   r.name,
					strict: r.strict,
				***REMOVED***
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if idx, exists := stash.names[r.name]; exists ***REMOVED***
				ref = &stashRef***REMOVED***
					v: &stash.values[idx],
				***REMOVED***
				goto end
			***REMOVED***
		***REMOVED***
		stash = stash.outer
	***REMOVED***

	if stash != nil ***REMOVED***
		ref = &stashRef***REMOVED***
			v: &stash.values[idx],
		***REMOVED***
		goto end
	***REMOVED*** /*else ***REMOVED***
		if vm.r.globalObject.self.hasProperty(nameVal) ***REMOVED***
			ref = &objRef***REMOVED***
				base: vm.r.globalObject.self,
				name: r.name,
			***REMOVED***
			goto end
		***REMOVED***
	***REMOVED*** */

	ref = &unresolvedRef***REMOVED***
		runtime: vm.r,
		name:    r.name,
	***REMOVED***

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

type getVar1 string

func (n getVar1) exec(vm *vm) ***REMOVED***
	name := string(n)
	var val Value
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if v, exists := stash.getByName(name, vm); exists ***REMOVED***
			val = v
			break
		***REMOVED***
	***REMOVED***
	if val == nil ***REMOVED***
		val = vm.r.globalObject.self.getStr(name)
		if val == nil ***REMOVED***
			vm.r.throwReferenceError(name)
		***REMOVED***
	***REMOVED***
	vm.push(val)
	vm.pc++
***REMOVED***

type getVar1Callee string

func (n getVar1Callee) exec(vm *vm) ***REMOVED***
	name := string(n)
	var val Value
	for stash := vm.stash; stash != nil; stash = stash.outer ***REMOVED***
		if v, exists := stash.getByName(name, vm); exists ***REMOVED***
			val = v
			break
		***REMOVED***
	***REMOVED***
	if val == nil ***REMOVED***
		val = vm.r.globalObject.self.getStr(name)
		if val == nil ***REMOVED***
			val = valueUnresolved***REMOVED***r: vm.r, ref: name***REMOVED***
		***REMOVED***
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
			if src, ok := srcVal.assertString(); ok ***REMOVED***
				var this Value
				if vm.sb != 0 ***REMOVED***
					this = vm.stack[vm.sb]
				***REMOVED*** else ***REMOVED***
					this = vm.r.globalObject
				***REMOVED***
				ret := vm.r.eval(src.String(), true, strict, this)
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
	obj := vm.r.toCallee(v)
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
	case *lazyObject:
		obj.self = f.create(obj)
		goto repeat
	default:
		vm.r.typeErrorResult(true, "Not a function: %s", obj.ToString())
	***REMOVED***
***REMOVED***

func (vm *vm) _nativeCall(f *nativeFuncObject, n int) ***REMOVED***
	if f.f != nil ***REMOVED***
		vm.pushCtx()
		vm.prg = nil
		vm.funcName = f.nameProp.get(nil).String()
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
	stackTail := vm.stack[vm.sp:]
	for i := range stackTail ***REMOVED***
		stackTail[i] = nil
	***REMOVED***
	vm.stack = vm.stack[:vm.sp]
***REMOVED***

type enterFunc uint32

func (e enterFunc) exec(vm *vm) ***REMOVED***
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
	// <- sp

	vm.newStash()
	offset := vm.args - int(e)
	vm.stash.values = make([]Value, e)
	if offset > 0 ***REMOVED***
		copy(vm.stash.values, vm.stack[vm.sp-vm.args:])
		vm.stash.extraArgs = make([]Value, offset)
		copy(vm.stash.extraArgs, vm.stack[vm.sp-offset:])
	***REMOVED*** else ***REMOVED***
		copy(vm.stash.values, vm.stack[vm.sp-vm.args:])
		vv := vm.stash.values[vm.args:]
		for i := range vv ***REMOVED***
			vv[i] = _undefined
		***REMOVED***
	***REMOVED***
	vm.sp -= vm.args
	vm.sb = vm.sp - 1
	vm.pc++
***REMOVED***

type _ret struct***REMOVED******REMOVED***

var ret _ret

func (_ret) exec(vm *vm) ***REMOVED***
	// callee -3
	// this -2
	// retval -1

	vm.stack[vm.sp-3] = vm.stack[vm.sp-1]
	vm.sp -= 2
	vm.popCtx()
	if vm.pc < 0 ***REMOVED***
		vm.halt = true
	***REMOVED***
***REMOVED***

type enterFuncStashless struct ***REMOVED***
	stackSize uint32
	args      uint32
***REMOVED***

func (e enterFuncStashless) exec(vm *vm) ***REMOVED***
	vm.sb = vm.sp - vm.args - 1
	var ss int
	d := int(e.args) - vm.args
	if d > 0 ***REMOVED***
		ss = int(e.stackSize) + d
		vm.args = int(e.args)
	***REMOVED*** else ***REMOVED***
		ss = int(e.stackSize)
	***REMOVED***
	sp := vm.sp
	if ss > 0 ***REMOVED***
		vm.sp += int(ss)
		vm.stack.expand(vm.sp)
		s := vm.stack[sp:vm.sp]
		for i := range s ***REMOVED***
			s[i] = _undefined
		***REMOVED***
	***REMOVED***
	vm.pc++
***REMOVED***

type _retStashless struct***REMOVED******REMOVED***

var retStashless _retStashless

func (_retStashless) exec(vm *vm) ***REMOVED***
	retval := vm.stack[vm.sp-1]
	vm.sp = vm.sb
	vm.stack[vm.sp-1] = retval
	vm.popCtx()
	if vm.pc < 0 ***REMOVED***
		vm.halt = true
	***REMOVED***
***REMOVED***

type newFunc struct ***REMOVED***
	prg    *Program
	name   string
	length uint32
	strict bool

	srcStart, srcEnd uint32
***REMOVED***

func (n *newFunc) exec(vm *vm) ***REMOVED***
	obj := vm.r.newFunc(n.name, int(n.length), n.strict)
	obj.prg = n.prg
	obj.stash = vm.stash
	obj.src = n.prg.src.src[n.srcStart:n.srcEnd]
	vm.push(obj.val)
	vm.pc++
***REMOVED***

type bindName string

func (d bindName) exec(vm *vm) ***REMOVED***
	if vm.stash != nil ***REMOVED***
		vm.stash.createBinding(string(d))
	***REMOVED*** else ***REMOVED***
		vm.r.globalObject.self._putProp(string(d), _undefined, true, true, false)
	***REMOVED***
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
		return o.self.toPrimitiveNumber()
	***REMOVED***
	return v
***REMOVED***

func cmp(px, py Value) Value ***REMOVED***
	var ret bool
	var nx, ny float64

	if xs, ok := px.assertString(); ok ***REMOVED***
		if ys, ok := py.assertString(); ok ***REMOVED***
			ret = xs.compareTo(ys) < 0
			goto end
		***REMOVED***
	***REMOVED***

	if xi, ok := px.assertInt(); ok ***REMOVED***
		if yi, ok := py.assertInt(); ok ***REMOVED***
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

	if right.self.hasInstance(left) ***REMOVED***
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

	if right.self.hasProperty(left) ***REMOVED***
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
	dynamic       bool
***REMOVED***

func (t try) exec(vm *vm) ***REMOVED***
	o := vm.pc
	vm.pc++
	ex := vm.runTry()
	if ex != nil && t.catchOffset > 0 ***REMOVED***
		// run the catch block (in try)
		vm.pc = o + int(t.catchOffset)
		// TODO: if ex.val is an Error, set the stack property
		if t.dynamic ***REMOVED***
			vm.newStash()
			vm.stash.putByIdx(0, ex.val)
		***REMOVED*** else ***REMOVED***
			vm.push(ex.val)
		***REMOVED***
		ex = vm.runTry()
		if t.dynamic ***REMOVED***
			vm.stash = vm.stash.outer
		***REMOVED***
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
		panic(ex)
	***REMOVED***
***REMOVED***

type _retFinally struct***REMOVED******REMOVED***

var retFinally _retFinally

func (_retFinally) exec(vm *vm) ***REMOVED***
	vm.pc++
***REMOVED***

type enterCatch string

func (varName enterCatch) exec(vm *vm) ***REMOVED***
	vm.stash.names = map[string]uint32***REMOVED***
		string(varName): 0,
	***REMOVED***
	vm.pc++
***REMOVED***

type _throw struct***REMOVED******REMOVED***

var throw _throw

func (_throw) exec(vm *vm) ***REMOVED***
	panic(vm.stack[vm.sp-1])
***REMOVED***

type _new uint32

func (n _new) exec(vm *vm) ***REMOVED***
	obj := vm.r.toObject(vm.stack[vm.sp-1-int(n)])
repeat:
	switch f := obj.self.(type) ***REMOVED***
	case *funcObject:
		args := make([]Value, n)
		copy(args, vm.stack[vm.sp-int(n):])
		vm.sp -= int(n)
		vm.stack[vm.sp-1] = f.construct(args)
	case *nativeFuncObject:
		vm._nativeNew(f, int(n))
	case *boundFuncObject:
		vm._nativeNew(&f.nativeFuncObject, int(n))
	case *lazyObject:
		obj.self = f.create(obj)
		goto repeat
	default:
		vm.r.typeErrorResult(true, "Not a constructor")
	***REMOVED***

	vm.pc++
***REMOVED***

func (vm *vm) _nativeNew(f *nativeFuncObject, n int) ***REMOVED***
	if f.construct != nil ***REMOVED***
		args := make([]Value, n)
		copy(args, vm.stack[vm.sp-n:])
		vm.sp -= n
		vm.stack[vm.sp-1] = f.construct(args)
	***REMOVED*** else ***REMOVED***
		vm.r.typeErrorResult(true, "Not a constructor")
	***REMOVED***
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
		args._put(strconv.Itoa(i), &mappedProperty***REMOVED***
			valueProperty: valueProperty***REMOVED***
				writable:     true,
				configurable: true,
				enumerable:   true,
			***REMOVED***,
			v: &vm.stash.values[i],
		***REMOVED***)
	***REMOVED***

	for _, v := range vm.stash.extraArgs ***REMOVED***
		args._put(strconv.Itoa(i), v)
		i++
	***REMOVED***

	args._putProp("callee", vm.stack[vm.sb-1], true, false, true)
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
		args._put(strconv.Itoa(i), v)
		i++
	***REMOVED***

	for _, v := range vm.stash.extraArgs ***REMOVED***
		args._put(strconv.Itoa(i), v)
		i++
	***REMOVED***

	args._putProp("length", intToValue(int64(vm.args)), true, false, true)
	args._put("callee", vm.r.global.throwerProperty)
	args._put("caller", vm.r.global.throwerProperty)
	vm.push(args.val)
	vm.pc++
***REMOVED***

type _enterWith struct***REMOVED******REMOVED***

var enterWith _enterWith

func (_enterWith) exec(vm *vm) ***REMOVED***
	vm.newStash()
	vm.stash.obj = vm.stack[vm.sp-1].ToObject(vm.r).self
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
		vm.iterStack = append(vm.iterStack, iterStackItem***REMOVED***f: v.ToObject(vm.r).self.enumerate(false, true)***REMOVED***)
	***REMOVED***
	vm.sp--
	vm.pc++
***REMOVED***

type enumNext int32

func (jmp enumNext) exec(vm *vm) ***REMOVED***
	l := len(vm.iterStack) - 1
	item, n := vm.iterStack[l].f()
	if n != nil ***REMOVED***
		vm.iterStack[l].val = newStringValue(item.name)
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
