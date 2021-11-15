package goja

type weakSetObject struct ***REMOVED***
	baseObject
	s weakMap
***REMOVED***

func (ws *weakSetObject) init() ***REMOVED***
	ws.baseObject.init()
	ws.s = weakMap(ws.val.runtime.genId())
***REMOVED***

func (r *Runtime) weakSetProto_add(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	wso, ok := thisObj.self.(*weakSetObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method WeakSet.prototype.add called on incompatible receiver %s", thisObj.String()))
	***REMOVED***
	wso.s.set(r.toObject(call.Argument(0)), nil)
	return call.This
***REMOVED***

func (r *Runtime) weakSetProto_delete(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	wso, ok := thisObj.self.(*weakSetObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method WeakSet.prototype.delete called on incompatible receiver %s", thisObj.String()))
	***REMOVED***
	obj, ok := call.Argument(0).(*Object)
	if ok && wso.s.remove(obj) ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) weakSetProto_has(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	wso, ok := thisObj.self.(*weakSetObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method WeakSet.prototype.has called on incompatible receiver %s", thisObj.String()))
	***REMOVED***
	obj, ok := call.Argument(0).(*Object)
	if ok && wso.s.has(obj) ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) populateWeakSetGeneric(s *Object, adderValue Value, iterable Value) ***REMOVED***
	adder := toMethod(adderValue)
	if adder == nil ***REMOVED***
		panic(r.NewTypeError("WeakSet.add is not set"))
	***REMOVED***
	iter := r.getIterator(iterable, nil)
	iter.iterate(func(val Value) ***REMOVED***
		adder(FunctionCall***REMOVED***This: s, Arguments: []Value***REMOVED***val***REMOVED******REMOVED***)
	***REMOVED***)
***REMOVED***

func (r *Runtime) builtin_newWeakSet(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("WeakSet"))
	***REMOVED***
	proto := r.getPrototypeFromCtor(newTarget, r.global.WeakSet, r.global.WeakSetPrototype)
	o := &Object***REMOVED***runtime: r***REMOVED***

	wso := &weakSetObject***REMOVED******REMOVED***
	wso.class = classWeakSet
	wso.val = o
	wso.extensible = true
	o.self = wso
	wso.prototype = proto
	wso.init()
	if len(args) > 0 ***REMOVED***
		if arg := args[0]; arg != nil && arg != _undefined && arg != _null ***REMOVED***
			adder := wso.getStr("add", nil)
			if adder == r.global.weakSetAdder ***REMOVED***
				if arr := r.checkStdArrayIter(arg); arr != nil ***REMOVED***
					for _, v := range arr.values ***REMOVED***
						wso.s.set(r.toObject(v), nil)
					***REMOVED***
					return o
				***REMOVED***
			***REMOVED***
			r.populateWeakSetGeneric(o, adder, arg)
		***REMOVED***
	***REMOVED***
	return o
***REMOVED***

func (r *Runtime) createWeakSetProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)

	o._putProp("constructor", r.global.WeakSet, true, false, true)
	r.global.weakSetAdder = r.newNativeFunc(r.weakSetProto_add, nil, "add", nil, 1)
	o._putProp("add", r.global.weakSetAdder, true, false, true)
	o._putProp("delete", r.newNativeFunc(r.weakSetProto_delete, nil, "delete", nil, 1), true, false, true)
	o._putProp("has", r.newNativeFunc(r.weakSetProto_has, nil, "has", nil, 1), true, false, true)

	o._putSym(SymToStringTag, valueProp(asciiString(classWeakSet), false, false, true))

	return o
***REMOVED***

func (r *Runtime) createWeakSet(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.builtin_newWeakSet, r.global.WeakSetPrototype, "WeakSet", 0)

	return o
***REMOVED***

func (r *Runtime) initWeakSet() ***REMOVED***
	r.global.WeakSetPrototype = r.newLazyObject(r.createWeakSetProto)
	r.global.WeakSet = r.newLazyObject(r.createWeakSet)

	r.addToGlobal("WeakSet", r.global.WeakSet)
***REMOVED***
