package goja

type setObject struct ***REMOVED***
	baseObject
	m *orderedMap
***REMOVED***

type setIterObject struct ***REMOVED***
	baseObject
	iter *orderedMapIter
	kind iterationKind
***REMOVED***

func (o *setIterObject) next() Value ***REMOVED***
	if o.iter == nil ***REMOVED***
		return o.val.runtime.createIterResultObject(_undefined, true)
	***REMOVED***

	entry := o.iter.next()
	if entry == nil ***REMOVED***
		o.iter = nil
		return o.val.runtime.createIterResultObject(_undefined, true)
	***REMOVED***

	var result Value
	switch o.kind ***REMOVED***
	case iterationKindValue:
		result = entry.key
	default:
		result = o.val.runtime.newArrayValues([]Value***REMOVED***entry.key, entry.key***REMOVED***)
	***REMOVED***

	return o.val.runtime.createIterResultObject(result, false)
***REMOVED***

func (so *setObject) init() ***REMOVED***
	so.baseObject.init()
	so.m = newOrderedMap(so.val.runtime.getHash())
***REMOVED***

func (r *Runtime) setProto_add(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	so, ok := thisObj.self.(*setObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Set.prototype.add called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	so.m.set(call.Argument(0), nil)
	return call.This
***REMOVED***

func (r *Runtime) setProto_clear(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	so, ok := thisObj.self.(*setObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Set.prototype.clear called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	so.m.clear()
	return _undefined
***REMOVED***

func (r *Runtime) setProto_delete(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	so, ok := thisObj.self.(*setObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Set.prototype.delete called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	return r.toBoolean(so.m.remove(call.Argument(0)))
***REMOVED***

func (r *Runtime) setProto_entries(call FunctionCall) Value ***REMOVED***
	return r.createSetIterator(call.This, iterationKindKeyValue)
***REMOVED***

func (r *Runtime) setProto_forEach(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	so, ok := thisObj.self.(*setObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Set.prototype.forEach called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***
	callbackFn, ok := r.toObject(call.Argument(0)).self.assertCallable()
	if !ok ***REMOVED***
		panic(r.NewTypeError("object is not a function %s"))
	***REMOVED***
	t := call.Argument(1)
	iter := so.m.newIter()
	for ***REMOVED***
		entry := iter.next()
		if entry == nil ***REMOVED***
			break
		***REMOVED***
		callbackFn(FunctionCall***REMOVED***This: t, Arguments: []Value***REMOVED***entry.key, entry.key, thisObj***REMOVED******REMOVED***)
	***REMOVED***

	return _undefined
***REMOVED***

func (r *Runtime) setProto_has(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	so, ok := thisObj.self.(*setObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Set.prototype.has called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	return r.toBoolean(so.m.has(call.Argument(0)))
***REMOVED***

func (r *Runtime) setProto_getSize(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	so, ok := thisObj.self.(*setObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method get Set.prototype.size called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	return intToValue(int64(so.m.size))
***REMOVED***

func (r *Runtime) setProto_values(call FunctionCall) Value ***REMOVED***
	return r.createSetIterator(call.This, iterationKindValue)
***REMOVED***

func (r *Runtime) builtin_newSet(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("Set"))
	***REMOVED***
	proto := r.getPrototypeFromCtor(newTarget, r.global.Set, r.global.SetPrototype)
	o := &Object***REMOVED***runtime: r***REMOVED***

	so := &setObject***REMOVED******REMOVED***
	so.class = classSet
	so.val = o
	so.extensible = true
	o.self = so
	so.prototype = proto
	so.init()
	if len(args) > 0 ***REMOVED***
		if arg := args[0]; arg != nil && arg != _undefined && arg != _null ***REMOVED***
			adder := so.getStr("add", nil)
			iter := r.getIterator(arg, nil)
			if adder == r.global.setAdder ***REMOVED***
				iter.iterate(func(item Value) ***REMOVED***
					so.m.set(item, nil)
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				adderFn := toMethod(adder)
				if adderFn == nil ***REMOVED***
					panic(r.NewTypeError("Set.add in missing"))
				***REMOVED***
				iter.iterate(func(item Value) ***REMOVED***
					adderFn(FunctionCall***REMOVED***This: o, Arguments: []Value***REMOVED***item***REMOVED******REMOVED***)
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return o
***REMOVED***

func (r *Runtime) createSetIterator(setValue Value, kind iterationKind) Value ***REMOVED***
	obj := r.toObject(setValue)
	setObj, ok := obj.self.(*setObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Object is not a Set"))
	***REMOVED***

	o := &Object***REMOVED***runtime: r***REMOVED***

	si := &setIterObject***REMOVED***
		iter: setObj.m.newIter(),
		kind: kind,
	***REMOVED***
	si.class = classSetIterator
	si.val = o
	si.extensible = true
	o.self = si
	si.prototype = r.global.SetIteratorPrototype
	si.init()

	return o
***REMOVED***

func (r *Runtime) setIterProto_next(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	if iter, ok := thisObj.self.(*setIterObject); ok ***REMOVED***
		return iter.next()
	***REMOVED***
	panic(r.NewTypeError("Method Set Iterator.prototype.next called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
***REMOVED***

func (r *Runtime) createSetProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)

	o._putProp("constructor", r.global.Set, true, false, true)
	r.global.setAdder = r.newNativeFunc(r.setProto_add, nil, "add", nil, 1)
	o._putProp("add", r.global.setAdder, true, false, true)

	o._putProp("clear", r.newNativeFunc(r.setProto_clear, nil, "clear", nil, 0), true, false, true)
	o._putProp("delete", r.newNativeFunc(r.setProto_delete, nil, "delete", nil, 1), true, false, true)
	o._putProp("forEach", r.newNativeFunc(r.setProto_forEach, nil, "forEach", nil, 1), true, false, true)
	o._putProp("has", r.newNativeFunc(r.setProto_has, nil, "has", nil, 1), true, false, true)
	o.setOwnStr("size", &valueProperty***REMOVED***
		getterFunc:   r.newNativeFunc(r.setProto_getSize, nil, "get size", nil, 0),
		accessor:     true,
		writable:     true,
		configurable: true,
	***REMOVED***, true)

	valuesFunc := r.newNativeFunc(r.setProto_values, nil, "values", nil, 0)
	o._putProp("values", valuesFunc, true, false, true)
	o._putProp("keys", valuesFunc, true, false, true)
	o._putProp("entries", r.newNativeFunc(r.setProto_entries, nil, "entries", nil, 0), true, false, true)
	o._putSym(SymIterator, valueProp(valuesFunc, true, false, true))
	o._putSym(SymToStringTag, valueProp(asciiString(classSet), false, false, true))

	return o
***REMOVED***

func (r *Runtime) createSet(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.builtin_newSet, r.global.SetPrototype, "Set", 0)
	r.putSpeciesReturnThis(o)

	return o
***REMOVED***

func (r *Runtime) createSetIterProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.IteratorPrototype, classObject)

	o._putProp("next", r.newNativeFunc(r.setIterProto_next, nil, "next", nil, 0), true, false, true)
	o._putSym(SymToStringTag, valueProp(asciiString(classSetIterator), false, false, true))

	return o
***REMOVED***

func (r *Runtime) initSet() ***REMOVED***
	r.global.SetIteratorPrototype = r.newLazyObject(r.createSetIterProto)

	r.global.SetPrototype = r.newLazyObject(r.createSetProto)
	r.global.Set = r.newLazyObject(r.createSet)

	r.addToGlobal("Set", r.global.Set)
***REMOVED***
