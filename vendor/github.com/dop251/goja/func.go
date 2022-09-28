package goja

import (
	"reflect"

	"github.com/dop251/goja/unistring"
)

type baseFuncObject struct ***REMOVED***
	baseObject

	lenProp valueProperty
***REMOVED***

type baseJsFuncObject struct ***REMOVED***
	baseFuncObject

	stash   *stash
	privEnv *privateEnv

	prg    *Program
	src    string
	strict bool
***REMOVED***

type funcObject struct ***REMOVED***
	baseJsFuncObject
***REMOVED***

type classFuncObject struct ***REMOVED***
	baseJsFuncObject
	initFields   *Program
	computedKeys []Value

	privateEnvType *privateEnvType
	privateMethods []Value

	derived bool
***REMOVED***

type methodFuncObject struct ***REMOVED***
	baseJsFuncObject
	homeObject *Object
***REMOVED***

type arrowFuncObject struct ***REMOVED***
	baseJsFuncObject
	funcObj   *Object
	newTarget Value
***REMOVED***

type nativeFuncObject struct ***REMOVED***
	baseFuncObject

	f         func(FunctionCall) Value
	construct func(args []Value, newTarget *Object) *Object
***REMOVED***

type boundFuncObject struct ***REMOVED***
	nativeFuncObject
	wrapped *Object
***REMOVED***

func (f *nativeFuncObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return f.f
***REMOVED***

func (f *funcObject) _addProto(n unistring.String) Value ***REMOVED***
	if n == "prototype" ***REMOVED***
		if _, exists := f.values[n]; !exists ***REMOVED***
			return f.addPrototype()
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (f *funcObject) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	return f.getStrWithOwnProp(f.getOwnPropStr(p), p, receiver)
***REMOVED***

func (f *funcObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if v := f._addProto(name); v != nil ***REMOVED***
		return v
	***REMOVED***

	return f.baseObject.getOwnPropStr(name)
***REMOVED***

func (f *funcObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	f._addProto(name)
	return f.baseObject.setOwnStr(name, val, throw)
***REMOVED***

func (f *funcObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return f._setForeignStr(name, f.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

func (f *funcObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	f._addProto(name)
	return f.baseObject.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (f *funcObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	f._addProto(name)
	return f.baseObject.deleteStr(name, throw)
***REMOVED***

func (f *funcObject) addPrototype() Value ***REMOVED***
	proto := f.val.runtime.NewObject()
	proto.self._putProp("constructor", f.val, true, false, true)
	return f._putProp("prototype", proto, true, false, false)
***REMOVED***

func (f *funcObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if f.baseObject.hasOwnPropertyStr(name) ***REMOVED***
		return true
	***REMOVED***

	if name == "prototype" ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (f *funcObject) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	if all ***REMOVED***
		if _, exists := f.values["prototype"]; !exists ***REMOVED***
			accum = append(accum, asciiString("prototype"))
		***REMOVED***
	***REMOVED***
	return f.baseFuncObject.stringKeys(all, accum)
***REMOVED***

func (f *funcObject) iterateStringKeys() iterNextFunc ***REMOVED***
	if _, exists := f.values["prototype"]; !exists ***REMOVED***
		f.addPrototype()
	***REMOVED***
	return f.baseFuncObject.iterateStringKeys()
***REMOVED***

func (f *baseFuncObject) createInstance(newTarget *Object) *Object ***REMOVED***
	r := f.val.runtime
	if newTarget == nil ***REMOVED***
		newTarget = f.val
	***REMOVED***
	proto := r.getPrototypeFromCtor(newTarget, nil, r.global.ObjectPrototype)

	return f.val.runtime.newBaseObject(proto, classObject).val
***REMOVED***

func (f *baseJsFuncObject) construct(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		newTarget = f.val
	***REMOVED***
	proto := newTarget.self.getStr("prototype", nil)
	var protoObj *Object
	if p, ok := proto.(*Object); ok ***REMOVED***
		protoObj = p
	***REMOVED*** else ***REMOVED***
		protoObj = f.val.runtime.global.ObjectPrototype
	***REMOVED***

	obj := f.val.runtime.newBaseObject(protoObj, classObject).val
	ret := f.call(FunctionCall***REMOVED***
		This:      obj,
		Arguments: args,
	***REMOVED***, newTarget)

	if ret, ok := ret.(*Object); ok ***REMOVED***
		return ret
	***REMOVED***
	return obj
***REMOVED***

func (f *classFuncObject) Call(FunctionCall) Value ***REMOVED***
	panic(f.val.runtime.NewTypeError("Class constructor cannot be invoked without 'new'"))
***REMOVED***

func (f *classFuncObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	return f.Call, true
***REMOVED***

func (f *classFuncObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return f.Call
***REMOVED***

func (f *classFuncObject) createInstance(args []Value, newTarget *Object) (instance *Object) ***REMOVED***
	if f.derived ***REMOVED***
		if ctor := f.prototype.self.assertConstructor(); ctor != nil ***REMOVED***
			instance = ctor(args, newTarget)
		***REMOVED*** else ***REMOVED***
			panic(f.val.runtime.NewTypeError("Super constructor is not a constructor"))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		instance = f.baseFuncObject.createInstance(newTarget)
	***REMOVED***
	return
***REMOVED***

func (f *classFuncObject) _initFields(instance *Object) ***REMOVED***
	if f.privateEnvType != nil ***REMOVED***
		penv := instance.self.getPrivateEnv(f.privateEnvType, true)
		penv.methods = f.privateMethods
	***REMOVED***
	if f.initFields != nil ***REMOVED***
		vm := f.val.runtime.vm
		vm.pushCtx()
		vm.prg = f.initFields
		vm.stash = f.stash
		vm.privEnv = f.privEnv
		vm.newTarget = nil

		// so that 'super' base could be correctly resolved (including from direct eval())
		vm.push(f.val)

		vm.sb = vm.sp
		vm.push(instance)
		vm.pc = 0
		vm.run()
		vm.popCtx()
		vm.sp -= 2
		vm.halt = false
	***REMOVED***
***REMOVED***

func (f *classFuncObject) construct(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		newTarget = f.val
	***REMOVED***
	if f.prg == nil ***REMOVED***
		instance := f.createInstance(args, newTarget)
		f._initFields(instance)
		return instance
	***REMOVED*** else ***REMOVED***
		var instance *Object
		var thisVal Value
		if !f.derived ***REMOVED***
			instance = f.createInstance(args, newTarget)
			f._initFields(instance)
			thisVal = instance
		***REMOVED***
		ret := f._call(args, newTarget, thisVal)

		if ret, ok := ret.(*Object); ok ***REMOVED***
			return ret
		***REMOVED***
		if f.derived ***REMOVED***
			r := f.val.runtime
			if ret != _undefined ***REMOVED***
				panic(r.NewTypeError("Derived constructors may only return object or undefined"))
			***REMOVED***
			if v := r.vm.stack[r.vm.sp+1]; v != nil ***REMOVED*** // using residual 'this' value (a bit hacky)
				instance = r.toObject(v)
			***REMOVED*** else ***REMOVED***
				panic(r.newError(r.global.ReferenceError, "Must call super constructor in derived class before returning from derived constructor"))
			***REMOVED***
		***REMOVED***
		return instance
	***REMOVED***
***REMOVED***

func (f *classFuncObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return f.construct
***REMOVED***

func (f *baseJsFuncObject) Call(call FunctionCall) Value ***REMOVED***
	return f.call(call, nil)
***REMOVED***

func (f *arrowFuncObject) Call(call FunctionCall) Value ***REMOVED***
	return f._call(call.Arguments, f.newTarget, nil)
***REMOVED***

func (f *baseJsFuncObject) _call(args []Value, newTarget, this Value) Value ***REMOVED***
	vm := f.val.runtime.vm

	vm.stack.expand(vm.sp + len(args) + 1)
	vm.stack[vm.sp] = f.val
	vm.sp++
	vm.stack[vm.sp] = this
	vm.sp++
	for _, arg := range args ***REMOVED***
		if arg != nil ***REMOVED***
			vm.stack[vm.sp] = arg
		***REMOVED*** else ***REMOVED***
			vm.stack[vm.sp] = _undefined
		***REMOVED***
		vm.sp++
	***REMOVED***

	pc := vm.pc
	if pc != -1 ***REMOVED***
		vm.pc++ // fake "return address" so that captureStack() records the correct call location
		vm.pushCtx()
		vm.callStack = append(vm.callStack, context***REMOVED***pc: -1***REMOVED***) // extra frame so that run() halts after ret
	***REMOVED*** else ***REMOVED***
		vm.pushCtx()
	***REMOVED***
	vm.args = len(args)
	vm.prg = f.prg
	vm.stash = f.stash
	vm.privEnv = f.privEnv
	vm.newTarget = newTarget
	vm.pc = 0
	vm.run()
	if pc != -1 ***REMOVED***
		vm.popCtx()
	***REMOVED***
	vm.pc = pc
	vm.halt = false
	return vm.pop()
***REMOVED***

func (f *baseJsFuncObject) call(call FunctionCall, newTarget Value) Value ***REMOVED***
	return f._call(call.Arguments, newTarget, nilSafe(call.This))
***REMOVED***

func (f *baseJsFuncObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return f.Call
***REMOVED***

func (f *baseFuncObject) exportType() reflect.Type ***REMOVED***
	return reflectTypeFunc
***REMOVED***

func (f *baseJsFuncObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	return f.Call, true
***REMOVED***

func (f *funcObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return f.construct
***REMOVED***

func (f *arrowFuncObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	return f.Call, true
***REMOVED***

func (f *arrowFuncObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return f.Call
***REMOVED***

func (f *baseFuncObject) init(name unistring.String, length Value) ***REMOVED***
	f.baseObject.init()

	f.lenProp.configurable = true
	f.lenProp.value = length
	f._put("length", &f.lenProp)

	f._putProp("name", stringValueFromRaw(name), false, false, true)
***REMOVED***

func (f *baseFuncObject) hasInstance(v Value) bool ***REMOVED***
	if v, ok := v.(*Object); ok ***REMOVED***
		o := f.val.self.getStr("prototype", nil)
		if o1, ok := o.(*Object); ok ***REMOVED***
			for ***REMOVED***
				v = v.self.proto()
				if v == nil ***REMOVED***
					return false
				***REMOVED***
				if o1 == v ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			f.val.runtime.typeErrorResult(true, "prototype is not an object")
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

func (f *nativeFuncObject) defaultConstruct(ccall func(ConstructorCall) *Object, args []Value, newTarget *Object) *Object ***REMOVED***
	obj := f.createInstance(newTarget)
	ret := ccall(ConstructorCall***REMOVED***
		This:      obj,
		Arguments: args,
		NewTarget: newTarget,
	***REMOVED***)

	if ret != nil ***REMOVED***
		return ret
	***REMOVED***
	return obj
***REMOVED***

func (f *nativeFuncObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	if f.f != nil ***REMOVED***
		return f.f, true
	***REMOVED***
	return nil, false
***REMOVED***

func (f *nativeFuncObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return f.construct
***REMOVED***

/*func (f *boundFuncObject) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	return f.getStrWithOwnProp(f.getOwnPropStr(p), p, receiver)
***REMOVED***

func (f *boundFuncObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		return f.val.runtime.global.throwerProperty
	***REMOVED***

	return f.nativeFuncObject.getOwnPropStr(name)
***REMOVED***

func (f *boundFuncObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		return true
	***REMOVED***
	return f.nativeFuncObject.deleteStr(name, throw)
***REMOVED***

func (f *boundFuncObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		panic(f.val.runtime.NewTypeError("'caller' and 'arguments' are restricted function properties and cannot be accessed in this context."))
	***REMOVED***
	return f.nativeFuncObject.setOwnStr(name, val, throw)
***REMOVED***

func (f *boundFuncObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return f._setForeignStr(name, f.getOwnPropStr(name), val, receiver, throw)
***REMOVED***
*/

func (f *boundFuncObject) hasInstance(v Value) bool ***REMOVED***
	return instanceOfOperator(v, f.wrapped)
***REMOVED***
