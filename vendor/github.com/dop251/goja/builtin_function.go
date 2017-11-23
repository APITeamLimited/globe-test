package goja

import (
	"fmt"
)

func (r *Runtime) builtin_Function(args []Value, proto *Object) *Object ***REMOVED***
	src := "(function anonymous("
	if len(args) > 1 ***REMOVED***
		for _, arg := range args[:len(args)-1] ***REMOVED***
			src += arg.String() + ","
		***REMOVED***
		src = src[:len(src)-1]
	***REMOVED***
	body := ""
	if len(args) > 0 ***REMOVED***
		body = args[len(args)-1].String()
	***REMOVED***
	src += ")***REMOVED***" + body + "***REMOVED***)"

	return r.toObject(r.eval(src, false, false, _undefined))
***REMOVED***

func (r *Runtime) functionproto_toString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
repeat:
	switch f := obj.self.(type) ***REMOVED***
	case *funcObject:
		return newStringValue(f.src)
	case *nativeFuncObject:
		return newStringValue(fmt.Sprintf("function %s() ***REMOVED*** [native code] ***REMOVED***", f.nameProp.get(call.This).ToString()))
	case *boundFuncObject:
		return newStringValue(fmt.Sprintf("function %s() ***REMOVED*** [native code] ***REMOVED***", f.nameProp.get(call.This).ToString()))
	case *lazyObject:
		obj.self = f.create(obj)
		goto repeat
	***REMOVED***

	r.typeErrorResult(true, "Object is not a function")
	return nil
***REMOVED***

func (r *Runtime) toValueArray(a Value) []Value ***REMOVED***
	obj := r.toObject(a)
	l := toUInt32(obj.self.getStr("length"))
	ret := make([]Value, l)
	for i := uint32(0); i < l; i++ ***REMOVED***
		ret[i] = obj.self.get(valueInt(i))
	***REMOVED***
	return ret
***REMOVED***

func (r *Runtime) functionproto_apply(call FunctionCall) Value ***REMOVED***
	f := r.toCallable(call.This)
	var args []Value
	if len(call.Arguments) >= 2 ***REMOVED***
		args = r.toValueArray(call.Arguments[1])
	***REMOVED***
	return f(FunctionCall***REMOVED***
		This:      call.Argument(0),
		Arguments: args,
	***REMOVED***)
***REMOVED***

func (r *Runtime) functionproto_call(call FunctionCall) Value ***REMOVED***
	f := r.toCallable(call.This)
	var args []Value
	if len(call.Arguments) > 0 ***REMOVED***
		args = call.Arguments[1:]
	***REMOVED***
	return f(FunctionCall***REMOVED***
		This:      call.Argument(0),
		Arguments: args,
	***REMOVED***)
***REMOVED***

func (r *Runtime) boundCallable(target func(FunctionCall) Value, boundArgs []Value) func(FunctionCall) Value ***REMOVED***
	var this Value
	var args []Value
	if len(boundArgs) > 0 ***REMOVED***
		this = boundArgs[0]
		args = make([]Value, len(boundArgs)-1)
		copy(args, boundArgs[1:])
	***REMOVED*** else ***REMOVED***
		this = _undefined
	***REMOVED***
	return func(call FunctionCall) Value ***REMOVED***
		a := append(args, call.Arguments...)
		return target(FunctionCall***REMOVED***
			This:      this,
			Arguments: a,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (r *Runtime) boundConstruct(target func([]Value) *Object, boundArgs []Value) func([]Value) *Object ***REMOVED***
	if target == nil ***REMOVED***
		return nil
	***REMOVED***
	var args []Value
	if len(boundArgs) > 1 ***REMOVED***
		args = make([]Value, len(boundArgs)-1)
		copy(args, boundArgs[1:])
	***REMOVED***
	return func(fargs []Value) *Object ***REMOVED***
		a := append(args, fargs...)
		copy(a, args)
		return target(a)
	***REMOVED***
***REMOVED***

func (r *Runtime) functionproto_bind(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	f := obj.self
	var fcall func(FunctionCall) Value
	var construct func([]Value) *Object
repeat:
	switch ff := f.(type) ***REMOVED***
	case *funcObject:
		fcall = ff.Call
		construct = ff.construct
	case *nativeFuncObject:
		fcall = ff.f
		construct = ff.construct
	case *boundFuncObject:
		f = &ff.nativeFuncObject
		goto repeat
	case *lazyObject:
		f = ff.create(obj)
		goto repeat
	default:
		r.typeErrorResult(true, "Value is not callable: %s", obj.ToString())
	***REMOVED***

	l := int(toUInt32(obj.self.getStr("length")))
	l -= len(call.Arguments) - 1
	if l < 0 ***REMOVED***
		l = 0
	***REMOVED***

	v := &Object***REMOVED***runtime: r***REMOVED***

	ff := r.newNativeFuncObj(v, r.boundCallable(fcall, call.Arguments), r.boundConstruct(construct, call.Arguments), "", nil, l)
	v.self = &boundFuncObject***REMOVED***
		nativeFuncObject: *ff,
	***REMOVED***

	//ret := r.newNativeFunc(r.boundCallable(f, call.Arguments), nil, "", nil, l)
	//o := ret.self
	//o.putStr("caller", r.global.throwerProperty, false)
	//o.putStr("arguments", r.global.throwerProperty, false)
	return v
***REMOVED***

func (r *Runtime) initFunction() ***REMOVED***
	o := r.global.FunctionPrototype.self
	o.(*nativeFuncObject).prototype = r.global.ObjectPrototype
	o._putProp("toString", r.newNativeFunc(r.functionproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("apply", r.newNativeFunc(r.functionproto_apply, nil, "apply", nil, 2), true, false, true)
	o._putProp("call", r.newNativeFunc(r.functionproto_call, nil, "call", nil, 1), true, false, true)
	o._putProp("bind", r.newNativeFunc(r.functionproto_bind, nil, "bind", nil, 1), true, false, true)

	r.global.Function = r.newNativeFuncConstruct(r.builtin_Function, "Function", r.global.FunctionPrototype, 1)
	r.addToGlobal("Function", r.global.Function)
***REMOVED***
