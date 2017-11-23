package goja

import (
	"fmt"
)

func (r *Runtime) builtin_Object(args []Value, proto *Object) *Object ***REMOVED***
	if len(args) > 0 ***REMOVED***
		arg := args[0]
		if arg != _undefined && arg != _null ***REMOVED***
			return arg.ToObject(r)
		***REMOVED***
	***REMOVED***
	return r.NewObject()
***REMOVED***

func (r *Runtime) object_getPrototypeOf(call FunctionCall) Value ***REMOVED***
	o := call.Argument(0).ToObject(r)
	p := o.self.proto()
	if p == nil ***REMOVED***
		return _null
	***REMOVED***
	return p
***REMOVED***

func (r *Runtime) object_getOwnPropertyDescriptor(call FunctionCall) Value ***REMOVED***
	obj := call.Argument(0).ToObject(r)
	propName := call.Argument(1).String()
	desc := obj.self.getOwnProp(propName)
	if desc == nil ***REMOVED***
		return _undefined
	***REMOVED***
	var writable, configurable, enumerable, accessor bool
	var get, set *Object
	var value Value
	if v, ok := desc.(*valueProperty); ok ***REMOVED***
		writable = v.writable
		configurable = v.configurable
		enumerable = v.enumerable
		accessor = v.accessor
		value = v.value
		get = v.getterFunc
		set = v.setterFunc
	***REMOVED*** else ***REMOVED***
		writable = true
		configurable = true
		enumerable = true
		value = desc
	***REMOVED***

	ret := r.NewObject()
	o := ret.self
	if !accessor ***REMOVED***
		o.putStr("value", value, false)
		o.putStr("writable", r.toBoolean(writable), false)
	***REMOVED*** else ***REMOVED***
		if get != nil ***REMOVED***
			o.putStr("get", get, false)
		***REMOVED*** else ***REMOVED***
			o.putStr("get", _undefined, false)
		***REMOVED***
		if set != nil ***REMOVED***
			o.putStr("set", set, false)
		***REMOVED*** else ***REMOVED***
			o.putStr("set", _undefined, false)
		***REMOVED***
	***REMOVED***
	o.putStr("enumerable", r.toBoolean(enumerable), false)
	o.putStr("configurable", r.toBoolean(configurable), false)

	return ret
***REMOVED***

func (r *Runtime) object_getOwnPropertyNames(call FunctionCall) Value ***REMOVED***
	// ES6
	obj := call.Argument(0).ToObject(r)
	// obj := r.toObject(call.Argument(0))

	var values []Value
	for item, f := obj.self.enumerate(true, false)(); f != nil; item, f = f() ***REMOVED***
		values = append(values, newStringValue(item.name))
	***REMOVED***
	return r.newArrayValues(values)
***REMOVED***

func (r *Runtime) toPropertyDescr(v Value) (ret propertyDescr) ***REMOVED***
	if o, ok := v.(*Object); ok ***REMOVED***
		descr := o.self

		ret.Value = descr.getStr("value")

		if p := descr.getStr("writable"); p != nil ***REMOVED***
			ret.Writable = ToFlag(p.ToBoolean())
		***REMOVED***
		if p := descr.getStr("enumerable"); p != nil ***REMOVED***
			ret.Enumerable = ToFlag(p.ToBoolean())
		***REMOVED***
		if p := descr.getStr("configurable"); p != nil ***REMOVED***
			ret.Configurable = ToFlag(p.ToBoolean())
		***REMOVED***

		ret.Getter = descr.getStr("get")
		ret.Setter = descr.getStr("set")

		if ret.Getter != nil && ret.Getter != _undefined ***REMOVED***
			if _, ok := r.toObject(ret.Getter).self.assertCallable(); !ok ***REMOVED***
				r.typeErrorResult(true, "getter must be a function")
			***REMOVED***
		***REMOVED***

		if ret.Setter != nil && ret.Setter != _undefined ***REMOVED***
			if _, ok := r.toObject(ret.Setter).self.assertCallable(); !ok ***REMOVED***
				r.typeErrorResult(true, "setter must be a function")
			***REMOVED***
		***REMOVED***

		if (ret.Getter != nil || ret.Setter != nil) && (ret.Value != nil || ret.Writable != FLAG_NOT_SET) ***REMOVED***
			r.typeErrorResult(true, "Invalid property descriptor. Cannot both specify accessors and a value or writable attribute")
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Property description must be an object: %s", v.String())
	***REMOVED***

	return
***REMOVED***

func (r *Runtime) _defineProperties(o *Object, p Value) ***REMOVED***
	type propItem struct ***REMOVED***
		name string
		prop propertyDescr
	***REMOVED***
	props := p.ToObject(r)
	var list []propItem
	for item, f := props.self.enumerate(false, false)(); f != nil; item, f = f() ***REMOVED***
		list = append(list, propItem***REMOVED***
			name: item.name,
			prop: r.toPropertyDescr(props.self.getStr(item.name)),
		***REMOVED***)
	***REMOVED***
	for _, prop := range list ***REMOVED***
		o.self.defineOwnProperty(newStringValue(prop.name), prop.prop, true)
	***REMOVED***
***REMOVED***

func (r *Runtime) object_create(call FunctionCall) Value ***REMOVED***
	var proto *Object
	if arg := call.Argument(0); arg != _null ***REMOVED***
		if o, ok := arg.(*Object); ok ***REMOVED***
			proto = o
		***REMOVED*** else ***REMOVED***
			r.typeErrorResult(true, "Object prototype may only be an Object or null: %s", arg.String())
		***REMOVED***
	***REMOVED***
	o := r.newBaseObject(proto, classObject).val

	if props := call.Argument(1); props != _undefined ***REMOVED***
		r._defineProperties(o, props)
	***REMOVED***

	return o
***REMOVED***

func (r *Runtime) object_defineProperty(call FunctionCall) (ret Value) ***REMOVED***
	if obj, ok := call.Argument(0).(*Object); ok ***REMOVED***
		descr := r.toPropertyDescr(call.Argument(2))
		obj.self.defineOwnProperty(call.Argument(1), descr, true)
		ret = call.Argument(0)
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Object.defineProperty called on non-object")
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) object_defineProperties(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.Argument(0))
	r._defineProperties(obj, call.Argument(1))
	return obj
***REMOVED***

func (r *Runtime) object_seal(call FunctionCall) Value ***REMOVED***
	// ES6
	arg := call.Argument(0)
	if obj, ok := arg.(*Object); ok ***REMOVED***
		descr := propertyDescr***REMOVED***
			Writable:     FLAG_TRUE,
			Enumerable:   FLAG_TRUE,
			Configurable: FLAG_FALSE,
		***REMOVED***
		for item, f := obj.self.enumerate(true, false)(); f != nil; item, f = f() ***REMOVED***
			v := obj.self.getOwnProp(item.name)
			if prop, ok := v.(*valueProperty); ok ***REMOVED***
				if !prop.configurable ***REMOVED***
					continue
				***REMOVED***
				prop.configurable = false
			***REMOVED*** else ***REMOVED***
				descr.Value = v
				obj.self.defineOwnProperty(newStringValue(item.name), descr, true)
				//obj.self._putProp(item.name, v, true, true, false)
			***REMOVED***
		***REMOVED***
		obj.self.preventExtensions()
		return obj
	***REMOVED***
	return arg
***REMOVED***

func (r *Runtime) object_freeze(call FunctionCall) Value ***REMOVED***
	arg := call.Argument(0)
	if obj, ok := arg.(*Object); ok ***REMOVED***
		descr := propertyDescr***REMOVED***
			Writable:     FLAG_FALSE,
			Enumerable:   FLAG_TRUE,
			Configurable: FLAG_FALSE,
		***REMOVED***
		for item, f := obj.self.enumerate(true, false)(); f != nil; item, f = f() ***REMOVED***
			v := obj.self.getOwnProp(item.name)
			if prop, ok := v.(*valueProperty); ok ***REMOVED***
				prop.configurable = false
				if prop.value != nil ***REMOVED***
					prop.writable = false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				descr.Value = v
				obj.self.defineOwnProperty(newStringValue(item.name), descr, true)
			***REMOVED***
		***REMOVED***
		obj.self.preventExtensions()
		return obj
	***REMOVED*** else ***REMOVED***
		// ES6 behavior
		return arg
	***REMOVED***
***REMOVED***

func (r *Runtime) object_preventExtensions(call FunctionCall) (ret Value) ***REMOVED***
	arg := call.Argument(0)
	if obj, ok := arg.(*Object); ok ***REMOVED***
		obj.self.preventExtensions()
		return obj
	***REMOVED***
	// ES6
	//r.typeErrorResult(true, "Object.preventExtensions called on non-object")
	//panic("Unreachable")
	return arg
***REMOVED***

func (r *Runtime) object_isSealed(call FunctionCall) Value ***REMOVED***
	if obj, ok := call.Argument(0).(*Object); ok ***REMOVED***
		if obj.self.isExtensible() ***REMOVED***
			return valueFalse
		***REMOVED***
		for item, f := obj.self.enumerate(true, false)(); f != nil; item, f = f() ***REMOVED***
			prop := obj.self.getOwnProp(item.name)
			if prop, ok := prop.(*valueProperty); ok ***REMOVED***
				if prop.configurable ***REMOVED***
					return valueFalse
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return valueFalse
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// ES6
		//r.typeErrorResult(true, "Object.isSealed called on non-object")
		return valueTrue
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) object_isFrozen(call FunctionCall) Value ***REMOVED***
	if obj, ok := call.Argument(0).(*Object); ok ***REMOVED***
		if obj.self.isExtensible() ***REMOVED***
			return valueFalse
		***REMOVED***
		for item, f := obj.self.enumerate(true, false)(); f != nil; item, f = f() ***REMOVED***
			prop := obj.self.getOwnProp(item.name)
			if prop, ok := prop.(*valueProperty); ok ***REMOVED***
				if prop.configurable || prop.value != nil && prop.writable ***REMOVED***
					return valueFalse
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return valueFalse
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// ES6
		//r.typeErrorResult(true, "Object.isFrozen called on non-object")
		return valueTrue
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) object_isExtensible(call FunctionCall) Value ***REMOVED***
	if obj, ok := call.Argument(0).(*Object); ok ***REMOVED***
		if obj.self.isExtensible() ***REMOVED***
			return valueTrue
		***REMOVED***
		return valueFalse
	***REMOVED*** else ***REMOVED***
		// ES6
		//r.typeErrorResult(true, "Object.isExtensible called on non-object")
		return valueFalse
	***REMOVED***
***REMOVED***

func (r *Runtime) object_keys(call FunctionCall) Value ***REMOVED***
	// ES6
	obj := call.Argument(0).ToObject(r)
	//if obj, ok := call.Argument(0).(*valueObject); ok ***REMOVED***
	var keys []Value
	for item, f := obj.self.enumerate(false, false)(); f != nil; item, f = f() ***REMOVED***
		keys = append(keys, newStringValue(item.name))
	***REMOVED***
	return r.newArrayValues(keys)
	//***REMOVED*** else ***REMOVED***
	//	r.typeErrorResult(true, "Object.keys called on non-object")
	//***REMOVED***
	//return nil
***REMOVED***

func (r *Runtime) objectproto_hasOwnProperty(call FunctionCall) Value ***REMOVED***
	p := call.Argument(0).String()
	o := call.This.ToObject(r)
	if o.self.hasOwnPropertyStr(p) ***REMOVED***
		return valueTrue
	***REMOVED*** else ***REMOVED***
		return valueFalse
	***REMOVED***
***REMOVED***

func (r *Runtime) objectproto_isPrototypeOf(call FunctionCall) Value ***REMOVED***
	if v, ok := call.Argument(0).(*Object); ok ***REMOVED***
		o := call.This.ToObject(r)
		for ***REMOVED***
			v = v.self.proto()
			if v == nil ***REMOVED***
				break
			***REMOVED***
			if v == o ***REMOVED***
				return valueTrue
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) objectproto_propertyIsEnumerable(call FunctionCall) Value ***REMOVED***
	p := call.Argument(0).ToString()
	o := call.This.ToObject(r)
	pv := o.self.getOwnProp(p.String())
	if pv == nil ***REMOVED***
		return valueFalse
	***REMOVED***
	if prop, ok := pv.(*valueProperty); ok ***REMOVED***
		if !prop.enumerable ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) objectproto_toString(call FunctionCall) Value ***REMOVED***
	switch o := call.This.(type) ***REMOVED***
	case valueNull:
		return stringObjectNull
	case valueUndefined:
		return stringObjectUndefined
	case *Object:
		return newStringValue(fmt.Sprintf("[object %s]", o.self.className()))
	default:
		obj := call.This.ToObject(r)
		return newStringValue(fmt.Sprintf("[object %s]", obj.self.className()))
	***REMOVED***
***REMOVED***

func (r *Runtime) objectproto_toLocaleString(call FunctionCall) Value ***REMOVED***
	return call.This.ToObject(r).ToString()
***REMOVED***

func (r *Runtime) objectproto_valueOf(call FunctionCall) Value ***REMOVED***
	return call.This.ToObject(r)
***REMOVED***

func (r *Runtime) initObject() ***REMOVED***
	o := r.global.ObjectPrototype.self
	o._putProp("toString", r.newNativeFunc(r.objectproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("toLocaleString", r.newNativeFunc(r.objectproto_toLocaleString, nil, "toLocaleString", nil, 0), true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.objectproto_valueOf, nil, "valueOf", nil, 0), true, false, true)
	o._putProp("hasOwnProperty", r.newNativeFunc(r.objectproto_hasOwnProperty, nil, "hasOwnProperty", nil, 1), true, false, true)
	o._putProp("isPrototypeOf", r.newNativeFunc(r.objectproto_isPrototypeOf, nil, "isPrototypeOf", nil, 1), true, false, true)
	o._putProp("propertyIsEnumerable", r.newNativeFunc(r.objectproto_propertyIsEnumerable, nil, "propertyIsEnumerable", nil, 1), true, false, true)

	r.global.Object = r.newNativeFuncConstruct(r.builtin_Object, classObject, r.global.ObjectPrototype, 1)
	o = r.global.Object.self
	o._putProp("defineProperty", r.newNativeFunc(r.object_defineProperty, nil, "defineProperty", nil, 3), true, false, true)
	o._putProp("defineProperties", r.newNativeFunc(r.object_defineProperties, nil, "defineProperties", nil, 2), true, false, true)
	o._putProp("getOwnPropertyDescriptor", r.newNativeFunc(r.object_getOwnPropertyDescriptor, nil, "getOwnPropertyDescriptor", nil, 2), true, false, true)
	o._putProp("getPrototypeOf", r.newNativeFunc(r.object_getPrototypeOf, nil, "getPrototypeOf", nil, 1), true, false, true)
	o._putProp("getOwnPropertyNames", r.newNativeFunc(r.object_getOwnPropertyNames, nil, "getOwnPropertyNames", nil, 1), true, false, true)
	o._putProp("create", r.newNativeFunc(r.object_create, nil, "create", nil, 2), true, false, true)
	o._putProp("seal", r.newNativeFunc(r.object_seal, nil, "seal", nil, 1), true, false, true)
	o._putProp("freeze", r.newNativeFunc(r.object_freeze, nil, "freeze", nil, 1), true, false, true)
	o._putProp("preventExtensions", r.newNativeFunc(r.object_preventExtensions, nil, "preventExtensions", nil, 1), true, false, true)
	o._putProp("isSealed", r.newNativeFunc(r.object_isSealed, nil, "isSealed", nil, 1), true, false, true)
	o._putProp("isFrozen", r.newNativeFunc(r.object_isFrozen, nil, "isFrozen", nil, 1), true, false, true)
	o._putProp("isExtensible", r.newNativeFunc(r.object_isExtensible, nil, "isExtensible", nil, 1), true, false, true)
	o._putProp("keys", r.newNativeFunc(r.object_keys, nil, "keys", nil, 1), true, false, true)

	r.addToGlobal("Object", r.global.Object)
***REMOVED***
