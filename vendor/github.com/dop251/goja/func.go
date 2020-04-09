package goja

import "reflect"

type baseFuncObject struct ***REMOVED***
	baseObject

	nameProp, lenProp valueProperty
***REMOVED***

type funcObject struct ***REMOVED***
	baseFuncObject

	stash *stash
	prg   *Program
	src   string
***REMOVED***

type nativeFuncObject struct ***REMOVED***
	baseFuncObject

	f         func(FunctionCall) Value
	construct func(args []Value) *Object
***REMOVED***

type boundFuncObject struct ***REMOVED***
	nativeFuncObject
***REMOVED***

func (f *nativeFuncObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	return f.f
***REMOVED***

func (f *nativeFuncObject) exportType() reflect.Type ***REMOVED***
	return reflect.TypeOf(f.f)
***REMOVED***

func (f *funcObject) _addProto(n string) Value ***REMOVED***
	if n == "prototype" ***REMOVED***
		if _, exists := f.values["prototype"]; !exists ***REMOVED***
			return f.addPrototype()
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (f *funcObject) getPropStr(name string) Value ***REMOVED***
	if v := f._addProto(name); v != nil ***REMOVED***
		return v
	***REMOVED***

	return f.baseObject.getPropStr(name)
***REMOVED***

func (f *funcObject) putStr(name string, val Value, throw bool) ***REMOVED***
	f._addProto(name)
	f.baseObject.putStr(name, val, throw)
***REMOVED***

func (f *funcObject) put(n Value, val Value, throw bool) ***REMOVED***
	f.putStr(n.String(), val, throw)
***REMOVED***

func (f *funcObject) deleteStr(name string, throw bool) bool ***REMOVED***
	f._addProto(name)
	return f.baseObject.deleteStr(name, throw)
***REMOVED***

func (f *funcObject) delete(n Value, throw bool) bool ***REMOVED***
	return f.deleteStr(n.String(), throw)
***REMOVED***

func (f *funcObject) addPrototype() Value ***REMOVED***
	proto := f.val.runtime.NewObject()
	proto.self._putProp("constructor", f.val, true, false, true)
	return f._putProp("prototype", proto, true, false, false)
***REMOVED***

func (f *funcObject) getProp(n Value) Value ***REMOVED***
	return f.getPropStr(n.String())
***REMOVED***

func (f *funcObject) hasOwnProperty(n Value) bool ***REMOVED***
	if r := f.baseObject.hasOwnProperty(n); r ***REMOVED***
		return true
	***REMOVED***

	name := n.String()
	if name == "prototype" ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (f *funcObject) hasOwnPropertyStr(name string) bool ***REMOVED***
	if r := f.baseObject.hasOwnPropertyStr(name); r ***REMOVED***
		return true
	***REMOVED***

	if name == "prototype" ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (f *funcObject) construct(args []Value) *Object ***REMOVED***
	proto := f.getStr("prototype")
	var protoObj *Object
	if p, ok := proto.(*Object); ok ***REMOVED***
		protoObj = p
	***REMOVED*** else ***REMOVED***
		protoObj = f.val.runtime.global.ObjectPrototype
	***REMOVED***
	obj := f.val.runtime.newBaseObject(protoObj, classObject).val
	ret := f.Call(FunctionCall***REMOVED***
		This:      obj,
		Arguments: args,
	***REMOVED***)

	if ret, ok := ret.(*Object); ok ***REMOVED***
		return ret
	***REMOVED***
	return obj
***REMOVED***

func (f *funcObject) Call(call FunctionCall) Value ***REMOVED***
	vm := f.val.runtime.vm
	pc := vm.pc

	vm.stack.expand(vm.sp + len(call.Arguments) + 1)
	vm.stack[vm.sp] = f.val
	vm.sp++
	if call.This != nil ***REMOVED***
		vm.stack[vm.sp] = call.This
	***REMOVED*** else ***REMOVED***
		vm.stack[vm.sp] = _undefined
	***REMOVED***
	vm.sp++
	for _, arg := range call.Arguments ***REMOVED***
		if arg != nil ***REMOVED***
			vm.stack[vm.sp] = arg
		***REMOVED*** else ***REMOVED***
			vm.stack[vm.sp] = _undefined
		***REMOVED***
		vm.sp++
	***REMOVED***

	vm.pc = -1
	vm.pushCtx()
	vm.args = len(call.Arguments)
	vm.prg = f.prg
	vm.stash = f.stash
	vm.pc = 0
	vm.run()
	vm.pc = pc
	vm.halt = false
	return vm.pop()
***REMOVED***

func (f *funcObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	return f.Call
***REMOVED***

func (f *funcObject) exportType() reflect.Type ***REMOVED***
	return reflect.TypeOf(f.Call)
***REMOVED***

func (f *funcObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	return f.Call, true
***REMOVED***

func (f *baseFuncObject) init(name string, length int) ***REMOVED***
	f.baseObject.init()

	f.nameProp.configurable = true
	f.nameProp.value = newStringValue(name)
	f._put("name", &f.nameProp)

	f.lenProp.configurable = true
	f.lenProp.value = valueInt(length)
	f._put("length", &f.lenProp)
***REMOVED***

func (f *baseFuncObject) hasInstance(v Value) bool ***REMOVED***
	if v, ok := v.(*Object); ok ***REMOVED***
		o := f.val.self.getStr("prototype")
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

func (f *nativeFuncObject) defaultConstruct(ccall func(ConstructorCall) *Object, args []Value) *Object ***REMOVED***
	proto := f.getStr("prototype")
	var protoObj *Object
	if p, ok := proto.(*Object); ok ***REMOVED***
		protoObj = p
	***REMOVED*** else ***REMOVED***
		protoObj = f.val.runtime.global.ObjectPrototype
	***REMOVED***
	obj := f.val.runtime.newBaseObject(protoObj, classObject).val
	ret := ccall(ConstructorCall***REMOVED***
		This:      obj,
		Arguments: args,
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

func (f *boundFuncObject) getProp(n Value) Value ***REMOVED***
	return f.getPropStr(n.String())
***REMOVED***

func (f *boundFuncObject) getPropStr(name string) Value ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		//f.runtime.typeErrorResult(true, "'caller' and 'arguments' are restricted function properties and cannot be accessed in this context.")
		return f.val.runtime.global.throwerProperty
	***REMOVED***
	return f.nativeFuncObject.getPropStr(name)
***REMOVED***

func (f *boundFuncObject) delete(n Value, throw bool) bool ***REMOVED***
	return f.deleteStr(n.String(), throw)
***REMOVED***

func (f *boundFuncObject) deleteStr(name string, throw bool) bool ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		return true
	***REMOVED***
	return f.nativeFuncObject.deleteStr(name, throw)
***REMOVED***

func (f *boundFuncObject) putStr(name string, val Value, throw bool) ***REMOVED***
	if name == "caller" || name == "arguments" ***REMOVED***
		f.val.runtime.typeErrorResult(true, "'caller' and 'arguments' are restricted function properties and cannot be accessed in this context.")
	***REMOVED***
	f.nativeFuncObject.putStr(name, val, throw)
***REMOVED***

func (f *boundFuncObject) put(n Value, val Value, throw bool) ***REMOVED***
	f.putStr(n.String(), val, throw)
***REMOVED***
