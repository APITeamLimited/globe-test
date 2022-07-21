package goja

import (
	"fmt"
	"math"
)

func (r *Runtime) builtin_Function(args []Value, proto *Object) *Object ***REMOVED***
	var sb valueStringBuilder
	sb.WriteString(asciiString("(function anonymous("))
	if len(args) > 1 ***REMOVED***
		ar := args[:len(args)-1]
		for i, arg := range ar ***REMOVED***
			sb.WriteString(arg.toString())
			if i < len(ar)-1 ***REMOVED***
				sb.WriteRune(',')
			***REMOVED***
		***REMOVED***
	***REMOVED***
	sb.WriteString(asciiString("\n) ***REMOVED***\n"))
	if len(args) > 0 ***REMOVED***
		sb.WriteString(args[len(args)-1].toString())
	***REMOVED***
	sb.WriteString(asciiString("\n***REMOVED***)"))

	ret := r.toObject(r.eval(sb.String(), false, false))
	ret.self.setProto(proto, true)
	return ret
***REMOVED***

func (r *Runtime) functionproto_toString(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
repeat:
	switch f := obj.self.(type) ***REMOVED***
	case *funcObject:
		return newStringValue(f.src)
	case *classFuncObject:
		return newStringValue(f.src)
	case *methodFuncObject:
		return newStringValue(f.src)
	case *arrowFuncObject:
		return newStringValue(f.src)
	case *nativeFuncObject:
		return newStringValue(fmt.Sprintf("function %s() ***REMOVED*** [native code] ***REMOVED***", nilSafe(f.getStr("name", nil)).toString()))
	case *boundFuncObject:
		return newStringValue(fmt.Sprintf("function %s() ***REMOVED*** [native code] ***REMOVED***", nilSafe(f.getStr("name", nil)).toString()))
	case *lazyObject:
		obj.self = f.create(obj)
		goto repeat
	case *proxyObject:
	repeat2:
		switch c := f.target.self.(type) ***REMOVED***
		case *classFuncObject, *methodFuncObject, *funcObject, *arrowFuncObject, *nativeFuncObject, *boundFuncObject:
			return asciiString("function () ***REMOVED*** [native code] ***REMOVED***")
		case *lazyObject:
			f.target.self = c.create(obj)
			goto repeat2
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Function.prototype.toString requires that 'this' be a Function"))
***REMOVED***

func (r *Runtime) functionproto_hasInstance(call FunctionCall) Value ***REMOVED***
	if o, ok := call.This.(*Object); ok ***REMOVED***
		if _, ok = o.self.assertCallable(); ok ***REMOVED***
			return r.toBoolean(o.self.hasInstance(call.Argument(0)))
		***REMOVED***
	***REMOVED***

	return valueFalse
***REMOVED***

func (r *Runtime) createListFromArrayLike(a Value) []Value ***REMOVED***
	o := r.toObject(a)
	if arr := r.checkStdArrayObj(o); arr != nil ***REMOVED***
		return arr.values
	***REMOVED***
	l := toLength(o.self.getStr("length", nil))
	res := make([]Value, 0, l)
	for k := int64(0); k < l; k++ ***REMOVED***
		res = append(res, nilSafe(o.self.getIdx(valueInt(k), nil)))
	***REMOVED***
	return res
***REMOVED***

func (r *Runtime) functionproto_apply(call FunctionCall) Value ***REMOVED***
	var args []Value
	if len(call.Arguments) >= 2 ***REMOVED***
		args = r.createListFromArrayLike(call.Arguments[1])
	***REMOVED***

	f := r.toCallable(call.This)
	return f(FunctionCall***REMOVED***
		This:      call.Argument(0),
		Arguments: args,
	***REMOVED***)
***REMOVED***

func (r *Runtime) functionproto_call(call FunctionCall) Value ***REMOVED***
	var args []Value
	if len(call.Arguments) > 0 ***REMOVED***
		args = call.Arguments[1:]
	***REMOVED***

	f := r.toCallable(call.This)
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

func (r *Runtime) boundConstruct(f *Object, target func([]Value, *Object) *Object, boundArgs []Value) func([]Value, *Object) *Object ***REMOVED***
	if target == nil ***REMOVED***
		return nil
	***REMOVED***
	var args []Value
	if len(boundArgs) > 1 ***REMOVED***
		args = make([]Value, len(boundArgs)-1)
		copy(args, boundArgs[1:])
	***REMOVED***
	return func(fargs []Value, newTarget *Object) *Object ***REMOVED***
		a := append(args, fargs...)
		if newTarget == f ***REMOVED***
			newTarget = nil
		***REMOVED***
		return target(a, newTarget)
	***REMOVED***
***REMOVED***

func (r *Runtime) functionproto_bind(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)

	fcall := r.toCallable(call.This)
	construct := obj.self.assertConstructor()

	var l = _positiveZero
	if obj.self.hasOwnPropertyStr("length") ***REMOVED***
		var li int64
		switch lenProp := nilSafe(obj.self.getStr("length", nil)).(type) ***REMOVED***
		case valueInt:
			li = lenProp.ToInteger()
		case valueFloat:
			switch lenProp ***REMOVED***
			case _positiveInf:
				l = lenProp
				goto lenNotInt
			case _negativeInf:
				goto lenNotInt
			case _negativeZero:
				// no-op, li == 0
			default:
				if !math.IsNaN(float64(lenProp)) ***REMOVED***
					li = int64(math.Abs(float64(lenProp)))
				***REMOVED*** // else li = 0
			***REMOVED***
		***REMOVED***
		if len(call.Arguments) > 1 ***REMOVED***
			li -= int64(len(call.Arguments)) - 1
		***REMOVED***
		if li < 0 ***REMOVED***
			li = 0
		***REMOVED***
		l = intToValue(li)
	***REMOVED***
lenNotInt:
	name := obj.self.getStr("name", nil)
	nameStr := stringBound_
	if s, ok := name.(valueString); ok ***REMOVED***
		nameStr = nameStr.concat(s)
	***REMOVED***

	v := &Object***REMOVED***runtime: r***REMOVED***
	ff := r.newNativeFuncAndConstruct(v, r.boundCallable(fcall, call.Arguments), r.boundConstruct(v, construct, call.Arguments), nil, nameStr.string(), l)
	bf := &boundFuncObject***REMOVED***
		nativeFuncObject: *ff,
		wrapped:          obj,
	***REMOVED***
	bf.prototype = obj.self.proto()
	v.self = bf

	return v
***REMOVED***

func (r *Runtime) initFunction() ***REMOVED***
	o := r.global.FunctionPrototype.self.(*nativeFuncObject)
	o.prototype = r.global.ObjectPrototype
	o._putProp("name", stringEmpty, false, false, true)
	o._putProp("apply", r.newNativeFunc(r.functionproto_apply, nil, "apply", nil, 2), true, false, true)
	o._putProp("bind", r.newNativeFunc(r.functionproto_bind, nil, "bind", nil, 1), true, false, true)
	o._putProp("call", r.newNativeFunc(r.functionproto_call, nil, "call", nil, 1), true, false, true)
	o._putProp("toString", r.newNativeFunc(r.functionproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putSym(SymHasInstance, valueProp(r.newNativeFunc(r.functionproto_hasInstance, nil, "[Symbol.hasInstance]", nil, 1), false, false, false))

	r.global.Function = r.newNativeFuncConstruct(r.builtin_Function, "Function", r.global.FunctionPrototype, 1)
	r.addToGlobal("Function", r.global.Function)
***REMOVED***
