package goja

import (
	"reflect"
)

var mapExportType = reflect.TypeOf([][2]interface***REMOVED******REMOVED******REMOVED******REMOVED***)

type mapObject struct ***REMOVED***
	baseObject
	m *orderedMap
***REMOVED***

type mapIterObject struct ***REMOVED***
	baseObject
	iter *orderedMapIter
	kind iterationKind
***REMOVED***

func (o *mapIterObject) next() Value ***REMOVED***
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
	case iterationKindKey:
		result = entry.key
	case iterationKindValue:
		result = entry.value
	default:
		result = o.val.runtime.newArrayValues([]Value***REMOVED***entry.key, entry.value***REMOVED***)
	***REMOVED***

	return o.val.runtime.createIterResultObject(result, false)
***REMOVED***

func (mo *mapObject) init() ***REMOVED***
	mo.baseObject.init()
	mo.m = newOrderedMap(mo.val.runtime.getHash())
***REMOVED***

func (mo *mapObject) exportType() reflect.Type ***REMOVED***
	return mapExportType
***REMOVED***

func (mo *mapObject) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	m := make([][2]interface***REMOVED******REMOVED***, mo.m.size)
	ctx.put(mo.val, m)

	iter := mo.m.newIter()
	for i := 0; i < len(m); i++ ***REMOVED***
		entry := iter.next()
		if entry == nil ***REMOVED***
			break
		***REMOVED***
		m[i][0] = exportValue(entry.key, ctx)
		m[i][1] = exportValue(entry.value, ctx)
	***REMOVED***

	return m
***REMOVED***

func (mo *mapObject) exportToMap(dst reflect.Value, typ reflect.Type, ctx *objectExportCtx) error ***REMOVED***
	if dst.IsNil() ***REMOVED***
		dst.Set(reflect.MakeMap(typ))
	***REMOVED***
	ctx.putTyped(mo.val, typ, dst.Interface())
	keyTyp := typ.Key()
	elemTyp := typ.Elem()
	iter := mo.m.newIter()
	r := mo.val.runtime
	for ***REMOVED***
		entry := iter.next()
		if entry == nil ***REMOVED***
			break
		***REMOVED***
		keyVal := reflect.New(keyTyp).Elem()
		err := r.toReflectValue(entry.key, keyVal, ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		elemVal := reflect.New(elemTyp).Elem()
		err = r.toReflectValue(entry.value, elemVal, ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		dst.SetMapIndex(keyVal, elemVal)
	***REMOVED***
	return nil
***REMOVED***

func (r *Runtime) mapProto_clear(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Map.prototype.clear called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	mo.m.clear()

	return _undefined
***REMOVED***

func (r *Runtime) mapProto_delete(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Map.prototype.delete called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	return r.toBoolean(mo.m.remove(call.Argument(0)))
***REMOVED***

func (r *Runtime) mapProto_get(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Map.prototype.get called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***

	return nilSafe(mo.m.get(call.Argument(0)))
***REMOVED***

func (r *Runtime) mapProto_has(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Map.prototype.has called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***
	if mo.m.has(call.Argument(0)) ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) mapProto_set(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Map.prototype.set called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***
	mo.m.set(call.Argument(0), call.Argument(1))
	return call.This
***REMOVED***

func (r *Runtime) mapProto_entries(call FunctionCall) Value ***REMOVED***
	return r.createMapIterator(call.This, iterationKindKeyValue)
***REMOVED***

func (r *Runtime) mapProto_forEach(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method Map.prototype.forEach called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***
	callbackFn, ok := r.toObject(call.Argument(0)).self.assertCallable()
	if !ok ***REMOVED***
		panic(r.NewTypeError("object is not a function %s"))
	***REMOVED***
	t := call.Argument(1)
	iter := mo.m.newIter()
	for ***REMOVED***
		entry := iter.next()
		if entry == nil ***REMOVED***
			break
		***REMOVED***
		callbackFn(FunctionCall***REMOVED***This: t, Arguments: []Value***REMOVED***entry.value, entry.key, thisObj***REMOVED******REMOVED***)
	***REMOVED***

	return _undefined
***REMOVED***

func (r *Runtime) mapProto_keys(call FunctionCall) Value ***REMOVED***
	return r.createMapIterator(call.This, iterationKindKey)
***REMOVED***

func (r *Runtime) mapProto_values(call FunctionCall) Value ***REMOVED***
	return r.createMapIterator(call.This, iterationKindValue)
***REMOVED***

func (r *Runtime) mapProto_getSize(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	mo, ok := thisObj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Method get Map.prototype.size called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
	***REMOVED***
	return intToValue(int64(mo.m.size))
***REMOVED***

func (r *Runtime) builtin_newMap(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("Map"))
	***REMOVED***
	proto := r.getPrototypeFromCtor(newTarget, r.global.Map, r.global.MapPrototype)
	o := &Object***REMOVED***runtime: r***REMOVED***

	mo := &mapObject***REMOVED******REMOVED***
	mo.class = classMap
	mo.val = o
	mo.extensible = true
	o.self = mo
	mo.prototype = proto
	mo.init()
	if len(args) > 0 ***REMOVED***
		if arg := args[0]; arg != nil && arg != _undefined && arg != _null ***REMOVED***
			adder := mo.getStr("set", nil)
			iter := r.getIterator(arg, nil)
			i0 := valueInt(0)
			i1 := valueInt(1)
			if adder == r.global.mapAdder ***REMOVED***
				iter.iterate(func(item Value) ***REMOVED***
					itemObj := r.toObject(item)
					k := nilSafe(itemObj.self.getIdx(i0, nil))
					v := nilSafe(itemObj.self.getIdx(i1, nil))
					mo.m.set(k, v)
				***REMOVED***)
			***REMOVED*** else ***REMOVED***
				adderFn := toMethod(adder)
				if adderFn == nil ***REMOVED***
					panic(r.NewTypeError("Map.set in missing"))
				***REMOVED***
				iter.iterate(func(item Value) ***REMOVED***
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

func (r *Runtime) createMapIterator(mapValue Value, kind iterationKind) Value ***REMOVED***
	obj := r.toObject(mapValue)
	mapObj, ok := obj.self.(*mapObject)
	if !ok ***REMOVED***
		panic(r.NewTypeError("Object is not a Map"))
	***REMOVED***

	o := &Object***REMOVED***runtime: r***REMOVED***

	mi := &mapIterObject***REMOVED***
		iter: mapObj.m.newIter(),
		kind: kind,
	***REMOVED***
	mi.class = classMapIterator
	mi.val = o
	mi.extensible = true
	o.self = mi
	mi.prototype = r.global.MapIteratorPrototype
	mi.init()

	return o
***REMOVED***

func (r *Runtime) mapIterProto_next(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	if iter, ok := thisObj.self.(*mapIterObject); ok ***REMOVED***
		return iter.next()
	***REMOVED***
	panic(r.NewTypeError("Method Map Iterator.prototype.next called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
***REMOVED***

func (r *Runtime) createMapProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)

	o._putProp("constructor", r.global.Map, true, false, true)
	o._putProp("clear", r.newNativeFunc(r.mapProto_clear, nil, "clear", nil, 0), true, false, true)
	r.global.mapAdder = r.newNativeFunc(r.mapProto_set, nil, "set", nil, 2)
	o._putProp("set", r.global.mapAdder, true, false, true)
	o._putProp("delete", r.newNativeFunc(r.mapProto_delete, nil, "delete", nil, 1), true, false, true)
	o._putProp("forEach", r.newNativeFunc(r.mapProto_forEach, nil, "forEach", nil, 1), true, false, true)
	o._putProp("has", r.newNativeFunc(r.mapProto_has, nil, "has", nil, 1), true, false, true)
	o._putProp("get", r.newNativeFunc(r.mapProto_get, nil, "get", nil, 1), true, false, true)
	o.setOwnStr("size", &valueProperty***REMOVED***
		getterFunc:   r.newNativeFunc(r.mapProto_getSize, nil, "get size", nil, 0),
		accessor:     true,
		writable:     true,
		configurable: true,
	***REMOVED***, true)
	o._putProp("keys", r.newNativeFunc(r.mapProto_keys, nil, "keys", nil, 0), true, false, true)
	o._putProp("values", r.newNativeFunc(r.mapProto_values, nil, "values", nil, 0), true, false, true)

	entriesFunc := r.newNativeFunc(r.mapProto_entries, nil, "entries", nil, 0)
	o._putProp("entries", entriesFunc, true, false, true)
	o._putSym(SymIterator, valueProp(entriesFunc, true, false, true))
	o._putSym(SymToStringTag, valueProp(asciiString(classMap), false, false, true))

	return o
***REMOVED***

func (r *Runtime) createMap(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.builtin_newMap, r.global.MapPrototype, "Map", 0)
	r.putSpeciesReturnThis(o)

	return o
***REMOVED***

func (r *Runtime) createMapIterProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.IteratorPrototype, classObject)

	o._putProp("next", r.newNativeFunc(r.mapIterProto_next, nil, "next", nil, 0), true, false, true)
	o._putSym(SymToStringTag, valueProp(asciiString(classMapIterator), false, false, true))

	return o
***REMOVED***

func (r *Runtime) initMap() ***REMOVED***
	r.global.MapIteratorPrototype = r.newLazyObject(r.createMapIterProto)

	r.global.MapPrototype = r.newLazyObject(r.createMapProto)
	r.global.Map = r.newLazyObject(r.createMap)

	r.addToGlobal("Map", r.global.Map)
***REMOVED***
