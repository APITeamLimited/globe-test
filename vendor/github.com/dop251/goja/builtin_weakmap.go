package goja

type weakMap struct ***REMOVED***
	data map[uint64]Value
***REMOVED***

type weakMapObject struct ***REMOVED***
	baseObject
	m *weakMap
***REMOVED***

func newWeakMap() *weakMap ***REMOVED***
	return &weakMap***REMOVED***
		data: make(map[uint64]Value),
	***REMOVED***
***REMOVED***

func (wmo *weakMapObject) init() ***REMOVED***
	wmo.baseObject.init()
	wmo.m = newWeakMap()
***REMOVED***

func (wm *weakMap) removeId(id uint64) ***REMOVED***
	delete(wm.data, id)
***REMOVED***

func (wm *weakMap) set(key *Object, value Value) ***REMOVED***
	ref := key.getWeakRef()
	wm.data[ref.id] = value
	key.runtime.addWeakKey(ref.id, wm)
***REMOVED***

func (wm *weakMap) get(key *Object) Value ***REMOVED***
	ref := key.weakRef
	if ref == nil ***REMOVED***
		return nil
	***REMOVED***
	ret := wm.data[ref.id]
	return ret
***REMOVED***

func (wm *weakMap) remove(key *Object) bool ***REMOVED***
	ref := key.weakRef
	if ref == nil ***REMOVED***
		return false
	***REMOVED***
	_, exists := wm.data[ref.id]
	if exists ***REMOVED***
		delete(wm.data, ref.id)
		key.runtime.removeWeakKey(ref.id, wm)
	***REMOVED***
	return exists
***REMOVED***

func (wm *weakMap) has(key *Object) bool ***REMOVED***
	ref := key.weakRef
	if ref == nil ***REMOVED***
		return false
	***REMOVED***
	_, exists := wm.data[ref.id]
	return exists
***REMOVED***

func (r *Runtime) weakMapProto_delete(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	wmo, ok := thisObj.self.(*weakMapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method WeakMap.prototype.delete called on incompatible receiver %s", thisObj.String()))
	***REMOVED***
	key, ok := call.Argument(0).(*Object)
	if ok && wmo.m.remove(key) ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) weakMapProto_get(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	wmo, ok := thisObj.self.(*weakMapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method WeakMap.prototype.get called on incompatible receiver %s", thisObj.String()))
	***REMOVED***
	var res Value
	if key, ok := call.Argument(0).(*Object); ok ***REMOVED***
		res = wmo.m.get(key)
	***REMOVED***
	if res == nil ***REMOVED***
		return _undefined
	***REMOVED***
	return res
***REMOVED***

func (r *Runtime) weakMapProto_has(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	wmo, ok := thisObj.self.(*weakMapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method WeakMap.prototype.has called on incompatible receiver %s", thisObj.String()))
	***REMOVED***
	key, ok := call.Argument(0).(*Object)
	if ok && wmo.m.has(key) ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) weakMapProto_set(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	wmo, ok := thisObj.self.(*weakMapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method WeakMap.prototype.set called on incompatible receiver %s", thisObj.String()))
	***REMOVED***
	key := r.toObject(call.Argument(0))
	wmo.m.set(key, call.Argument(1))
	return call.This
***REMOVED***

func (r *Runtime) needNew(name string) *Object ***REMOVED***
	return r.NewTypeError("Constructor %s requires 'new'", name)
***REMOVED***

func (r *Runtime) getPrototypeFromCtor(newTarget, defCtor, defProto *Object) *Object ***REMOVED***
	if newTarget == defCtor ***REMOVED***
		return defProto
	***REMOVED***
	proto := newTarget.self.getStr("prototype", nil)
	if obj, ok := proto.(*Object); ok ***REMOVED***
		return obj
	***REMOVED***
	return defProto
***REMOVED***

func (r *Runtime) builtin_newWeakMap(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("WeakMap"))
	***REMOVED***
	proto := r.getPrototypeFromCtor(newTarget, r.global.WeakMap, r.global.WeakMapPrototype)
	o := &Object***REMOVED***runtime: r***REMOVED***

	wmo := &weakMapObject***REMOVED******REMOVED***
	wmo.class = classWeakMap
	wmo.val = o
	wmo.extensible = true
	o.self = wmo
	wmo.prototype = proto
	wmo.init()
	if len(args) > 0 ***REMOVED***
		if arg := args[0]; arg != nil && arg != _undefined && arg != _null ***REMOVED***
			adder := wmo.getStr("set", nil)
			iter := r.getIterator(arg, nil)
			i0 := valueInt(0)
			i1 := valueInt(1)
			if adder == r.global.weakMapAdder ***REMOVED***
				r.iterate(iter, func(item Value) ***REMOVED***
					itemObj := r.toObject(item)
					k := itemObj.self.getIdx(i0, nil)
					v := nilSafe(itemObj.self.getIdx(i1, nil))
					wmo.m.set(r.toObject(k), v)
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				adderFn := toMethod(adder)
				if adderFn == nil ***REMOVED***
					panic(r.NewTypeError("WeakMap.set in missing"))
				***REMOVED***
				r.iterate(iter, func(item Value) ***REMOVED***
					itemObj := r.toObject(item)
					k := itemObj.self.getIdx(i0, nil)
					v := itemObj.self.getIdx(i1, nil)
					adderFn(FunctionCall***REMOVED***This: o, Arguments: []Value***REMOVED***k, v***REMOVED******REMOVED***)
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return o
***REMOVED***

func (r *Runtime) createWeakMapProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)

	o._putProp("constructor", r.global.WeakMap, true, false, true)
	r.global.weakMapAdder = r.newNativeFunc(r.weakMapProto_set, nil, "set", nil, 2)
	o._putProp("set", r.global.weakMapAdder, true, false, true)
	o._putProp("delete", r.newNativeFunc(r.weakMapProto_delete, nil, "delete", nil, 1), true, false, true)
	o._putProp("has", r.newNativeFunc(r.weakMapProto_has, nil, "has", nil, 1), true, false, true)
	o._putProp("get", r.newNativeFunc(r.weakMapProto_get, nil, "get", nil, 1), true, false, true)

	o._putSym(symToStringTag, valueProp(asciiString(classWeakMap), false, false, true))

	return o
***REMOVED***

func (r *Runtime) createWeakMap(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.builtin_newWeakMap, r.global.WeakMapPrototype, "WeakMap", 0)

	return o
***REMOVED***

func (r *Runtime) initWeakMap() ***REMOVED***
	r.global.WeakMapPrototype = r.newLazyObject(r.createWeakMapProto)
	r.global.WeakMap = r.newLazyObject(r.createWeakMap)

	r.addToGlobal("WeakMap", r.global.WeakMap)
***REMOVED***
