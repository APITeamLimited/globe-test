package goja

func (r *Runtime) builtin_reflect_apply(call FunctionCall) Value ***REMOVED***
	return r.toCallable(call.Argument(0))(FunctionCall***REMOVED***
		This:      call.Argument(1),
		Arguments: r.createListFromArrayLike(call.Argument(2))***REMOVED***)
***REMOVED***

func (r *Runtime) toConstructor(v Value) func(args []Value, newTarget *Object) *Object ***REMOVED***
	if ctor := r.toObject(v).self.assertConstructor(); ctor != nil ***REMOVED***
		return ctor
	***REMOVED***
	panic(r.NewTypeError("Value is not a constructor"))
***REMOVED***

func (r *Runtime) builtin_reflect_construct(call FunctionCall) Value ***REMOVED***
	target := call.Argument(0)
	ctor := r.toConstructor(target)
	var newTarget Value
	if len(call.Arguments) > 2 ***REMOVED***
		newTarget = call.Argument(2)
		r.toConstructor(newTarget)
	***REMOVED*** else ***REMOVED***
		newTarget = target
	***REMOVED***
	return ctor(r.createListFromArrayLike(call.Argument(1)), r.toObject(newTarget))
***REMOVED***

func (r *Runtime) builtin_reflect_defineProperty(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	key := toPropertyKey(call.Argument(1))
	desc := r.toPropertyDescriptor(call.Argument(2))

	return r.toBoolean(target.defineOwnProperty(key, desc, false))
***REMOVED***

func (r *Runtime) builtin_reflect_deleteProperty(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	key := toPropertyKey(call.Argument(1))

	return r.toBoolean(target.delete(key, false))
***REMOVED***

func (r *Runtime) builtin_reflect_get(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	key := toPropertyKey(call.Argument(1))
	var receiver Value
	if len(call.Arguments) > 2 ***REMOVED***
		receiver = call.Arguments[2]
	***REMOVED***
	return target.get(key, receiver)
***REMOVED***

func (r *Runtime) builtin_reflect_getOwnPropertyDescriptor(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	key := toPropertyKey(call.Argument(1))
	return r.valuePropToDescriptorObject(target.getOwnProp(key))
***REMOVED***

func (r *Runtime) builtin_reflect_getPrototypeOf(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	if proto := target.self.proto(); proto != nil ***REMOVED***
		return proto
	***REMOVED***

	return _null
***REMOVED***

func (r *Runtime) builtin_reflect_has(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	key := toPropertyKey(call.Argument(1))
	return r.toBoolean(target.hasProperty(key))
***REMOVED***

func (r *Runtime) builtin_reflect_isExtensible(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	return r.toBoolean(target.self.isExtensible())
***REMOVED***

func (r *Runtime) builtin_reflect_ownKeys(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	return r.newArrayValues(target.self.keys(true, nil))
***REMOVED***

func (r *Runtime) builtin_reflect_preventExtensions(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	return r.toBoolean(target.self.preventExtensions(false))
***REMOVED***

func (r *Runtime) builtin_reflect_set(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	var receiver Value
	if len(call.Arguments) >= 4 ***REMOVED***
		receiver = call.Argument(3)
	***REMOVED*** else ***REMOVED***
		receiver = target
	***REMOVED***
	return r.toBoolean(target.set(call.Argument(1), call.Argument(2), receiver, false))
***REMOVED***

func (r *Runtime) builtin_reflect_setPrototypeOf(call FunctionCall) Value ***REMOVED***
	target := r.toObject(call.Argument(0))
	var proto *Object
	if arg := call.Argument(1); arg != _null ***REMOVED***
		proto = r.toObject(arg)
	***REMOVED***
	return r.toBoolean(target.self.setProto(proto, false))
***REMOVED***

func (r *Runtime) createReflect(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)

	o._putProp("apply", r.newNativeFunc(r.builtin_reflect_apply, nil, "apply", nil, 3), true, false, true)
	o._putProp("construct", r.newNativeFunc(r.builtin_reflect_construct, nil, "construct", nil, 2), true, false, true)
	o._putProp("defineProperty", r.newNativeFunc(r.builtin_reflect_defineProperty, nil, "defineProperty", nil, 3), true, false, true)
	o._putProp("deleteProperty", r.newNativeFunc(r.builtin_reflect_deleteProperty, nil, "deleteProperty", nil, 2), true, false, true)
	o._putProp("get", r.newNativeFunc(r.builtin_reflect_get, nil, "get", nil, 2), true, false, true)
	o._putProp("getOwnPropertyDescriptor", r.newNativeFunc(r.builtin_reflect_getOwnPropertyDescriptor, nil, "getOwnPropertyDescriptor", nil, 2), true, false, true)
	o._putProp("getPrototypeOf", r.newNativeFunc(r.builtin_reflect_getPrototypeOf, nil, "getPrototypeOf", nil, 1), true, false, true)
	o._putProp("has", r.newNativeFunc(r.builtin_reflect_has, nil, "has", nil, 2), true, false, true)
	o._putProp("isExtensible", r.newNativeFunc(r.builtin_reflect_isExtensible, nil, "isExtensible", nil, 1), true, false, true)
	o._putProp("ownKeys", r.newNativeFunc(r.builtin_reflect_ownKeys, nil, "ownKeys", nil, 1), true, false, true)
	o._putProp("preventExtensions", r.newNativeFunc(r.builtin_reflect_preventExtensions, nil, "preventExtensions", nil, 1), true, false, true)
	o._putProp("set", r.newNativeFunc(r.builtin_reflect_set, nil, "set", nil, 3), true, false, true)
	o._putProp("setPrototypeOf", r.newNativeFunc(r.builtin_reflect_setPrototypeOf, nil, "setPrototypeOf", nil, 2), true, false, true)

	return o
***REMOVED***

func (r *Runtime) initReflect() ***REMOVED***
	r.addToGlobal("Reflect", r.newLazyObject(r.createReflect))
***REMOVED***
