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
	return r.newBaseObject(proto, classObject).val
***REMOVED***

func (r *Runtime) object_getPrototypeOf(call FunctionCall) Value ***REMOVED***
	o := call.Argument(0).ToObject(r)
	p := o.self.proto()
	if p == nil ***REMOVED***
		return _null
	***REMOVED***
	return p
***REMOVED***

func (r *Runtime) valuePropToDescriptorObject(desc Value) Value ***REMOVED***
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
	obj := ret.self
	if !accessor ***REMOVED***
		obj.setOwnStr("value", value, false)
		obj.setOwnStr("writable", r.toBoolean(writable), false)
	***REMOVED*** else ***REMOVED***
		if get != nil ***REMOVED***
			obj.setOwnStr("get", get, false)
		***REMOVED*** else ***REMOVED***
			obj.setOwnStr("get", _undefined, false)
		***REMOVED***
		if set != nil ***REMOVED***
			obj.setOwnStr("set", set, false)
		***REMOVED*** else ***REMOVED***
			obj.setOwnStr("set", _undefined, false)
		***REMOVED***
	***REMOVED***
	obj.setOwnStr("enumerable", r.toBoolean(enumerable), false)
	obj.setOwnStr("configurable", r.toBoolean(configurable), false)

	return ret
***REMOVED***

func (r *Runtime) object_getOwnPropertyDescriptor(call FunctionCall) Value ***REMOVED***
	o := call.Argument(0).ToObject(r)
	propName := toPropertyKey(call.Argument(1))
	return r.valuePropToDescriptorObject(o.getOwnProp(propName))
***REMOVED***

func (r *Runtime) object_getOwnPropertyNames(call FunctionCall) Value ***REMOVED***
	obj := call.Argument(0).ToObject(r)

	return r.newArrayValues(obj.self.ownKeys(true, nil))
***REMOVED***

func (r *Runtime) object_getOwnPropertySymbols(call FunctionCall) Value ***REMOVED***
	obj := call.Argument(0).ToObject(r)
	return r.newArrayValues(obj.self.ownSymbols(true, nil))
***REMOVED***

func (r *Runtime) toValueProp(v Value) *valueProperty ***REMOVED***
	if v == nil || v == _undefined ***REMOVED***
		return nil
	***REMOVED***
	obj := r.toObject(v)
	getter := obj.self.getStr("get", nil)
	setter := obj.self.getStr("set", nil)
	writable := obj.self.getStr("writable", nil)
	value := obj.self.getStr("value", nil)
	if (getter != nil || setter != nil) && (value != nil || writable != nil) ***REMOVED***
		r.typeErrorResult(true, "Invalid property descriptor. Cannot both specify accessors and a value or writable attribute")
	***REMOVED***

	ret := &valueProperty***REMOVED******REMOVED***
	if writable != nil && writable.ToBoolean() ***REMOVED***
		ret.writable = true
	***REMOVED***
	if e := obj.self.getStr("enumerable", nil); e != nil && e.ToBoolean() ***REMOVED***
		ret.enumerable = true
	***REMOVED***
	if c := obj.self.getStr("configurable", nil); c != nil && c.ToBoolean() ***REMOVED***
		ret.configurable = true
	***REMOVED***
	ret.value = value

	if getter != nil && getter != _undefined ***REMOVED***
		o := r.toObject(getter)
		if _, ok := o.self.assertCallable(); !ok ***REMOVED***
			r.typeErrorResult(true, "getter must be a function")
		***REMOVED***
		ret.getterFunc = o
	***REMOVED***

	if setter != nil && setter != _undefined ***REMOVED***
		o := r.toObject(v)
		if _, ok := o.self.assertCallable(); !ok ***REMOVED***
			r.typeErrorResult(true, "setter must be a function")
		***REMOVED***
		ret.setterFunc = o
	***REMOVED***

	if ret.getterFunc != nil || ret.setterFunc != nil ***REMOVED***
		ret.accessor = true
	***REMOVED***

	return ret
***REMOVED***

func (r *Runtime) toPropertyDescriptor(v Value) (ret PropertyDescriptor) ***REMOVED***
	if o, ok := v.(*Object); ok ***REMOVED***
		descr := o.self

		// Save the original descriptor for reference
		ret.jsDescriptor = o

		ret.Value = descr.getStr("value", nil)

		if p := descr.getStr("writable", nil); p != nil ***REMOVED***
			ret.Writable = ToFlag(p.ToBoolean())
		***REMOVED***
		if p := descr.getStr("enumerable", nil); p != nil ***REMOVED***
			ret.Enumerable = ToFlag(p.ToBoolean())
		***REMOVED***
		if p := descr.getStr("configurable", nil); p != nil ***REMOVED***
			ret.Configurable = ToFlag(p.ToBoolean())
		***REMOVED***

		ret.Getter = descr.getStr("get", nil)
		ret.Setter = descr.getStr("set", nil)

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
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Property description must be an object: %s", v.String())
	***REMOVED***

	return
***REMOVED***

func (r *Runtime) _defineProperties(o *Object, p Value) ***REMOVED***
	type propItem struct ***REMOVED***
		name Value
		prop PropertyDescriptor
	***REMOVED***
	props := p.ToObject(r)
	names := props.self.ownPropertyKeys(false, nil)
	list := make([]propItem, 0, len(names))
	for _, itemName := range names ***REMOVED***
		list = append(list, propItem***REMOVED***
			name: itemName,
			prop: r.toPropertyDescriptor(props.get(itemName, nil)),
		***REMOVED***)
	***REMOVED***
	for _, prop := range list ***REMOVED***
		o.defineOwnProperty(prop.name, prop.prop, true)
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
		descr := r.toPropertyDescriptor(call.Argument(2))
		obj.defineOwnProperty(toPropertyKey(call.Argument(1)), descr, true)
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
		descr := PropertyDescriptor***REMOVED***
			Writable:     FLAG_TRUE,
			Enumerable:   FLAG_TRUE,
			Configurable: FLAG_FALSE,
		***REMOVED***
		for _, key := range obj.self.ownPropertyKeys(true, nil) ***REMOVED***
			v := obj.getOwnProp(key)
			if prop, ok := v.(*valueProperty); ok ***REMOVED***
				if !prop.configurable ***REMOVED***
					continue
				***REMOVED***
				prop.configurable = false
			***REMOVED*** else ***REMOVED***
				descr.Value = v
				obj.defineOwnProperty(key, descr, true)
			***REMOVED***
		***REMOVED***
		obj.self.preventExtensions(false)
		return obj
	***REMOVED***
	return arg
***REMOVED***

func (r *Runtime) object_freeze(call FunctionCall) Value ***REMOVED***
	arg := call.Argument(0)
	if obj, ok := arg.(*Object); ok ***REMOVED***
		descr := PropertyDescriptor***REMOVED***
			Writable:     FLAG_FALSE,
			Enumerable:   FLAG_TRUE,
			Configurable: FLAG_FALSE,
		***REMOVED***
		for _, key := range obj.self.ownPropertyKeys(true, nil) ***REMOVED***
			v := obj.getOwnProp(key)
			if prop, ok := v.(*valueProperty); ok ***REMOVED***
				prop.configurable = false
				if prop.value != nil ***REMOVED***
					prop.writable = false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				descr.Value = v
				obj.defineOwnProperty(key, descr, true)
			***REMOVED***
		***REMOVED***
		obj.self.preventExtensions(false)
		return obj
	***REMOVED*** else ***REMOVED***
		// ES6 behavior
		return arg
	***REMOVED***
***REMOVED***

func (r *Runtime) object_preventExtensions(call FunctionCall) (ret Value) ***REMOVED***
	arg := call.Argument(0)
	if obj, ok := arg.(*Object); ok ***REMOVED***
		obj.self.preventExtensions(false)
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
		for _, key := range obj.self.ownPropertyKeys(true, nil) ***REMOVED***
			prop := obj.getOwnProp(key)
			if prop, ok := prop.(*valueProperty); ok ***REMOVED***
				if prop.configurable ***REMOVED***
					return valueFalse
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return valueFalse
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) object_isFrozen(call FunctionCall) Value ***REMOVED***
	if obj, ok := call.Argument(0).(*Object); ok ***REMOVED***
		if obj.self.isExtensible() ***REMOVED***
			return valueFalse
		***REMOVED***
		for _, key := range obj.self.ownPropertyKeys(true, nil) ***REMOVED***
			prop := obj.getOwnProp(key)
			if prop, ok := prop.(*valueProperty); ok ***REMOVED***
				if prop.configurable || prop.value != nil && prop.writable ***REMOVED***
					return valueFalse
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				return valueFalse
			***REMOVED***
		***REMOVED***
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
	obj := call.Argument(0).ToObject(r)

	return r.newArrayValues(obj.self.ownKeys(false, nil))
***REMOVED***

func (r *Runtime) objectproto_hasOwnProperty(call FunctionCall) Value ***REMOVED***
	p := toPropertyKey(call.Argument(0))
	o := call.This.ToObject(r)
	if o.hasOwnProperty(p) ***REMOVED***
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
	p := toPropertyKey(call.Argument(0))
	o := call.This.ToObject(r)
	pv := o.getOwnProp(p)
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
	default:
		obj := o.ToObject(r)
		var clsName string
		if isArray(obj) ***REMOVED***
			clsName = classArray
		***REMOVED*** else ***REMOVED***
			clsName = obj.self.className()
		***REMOVED***
		if tag := obj.self.getSym(SymToStringTag, nil); tag != nil ***REMOVED***
			if str, ok := tag.(valueString); ok ***REMOVED***
				clsName = str.String()
			***REMOVED***
		***REMOVED***
		return newStringValue(fmt.Sprintf("[object %s]", clsName))
	***REMOVED***
***REMOVED***

func (r *Runtime) objectproto_toLocaleString(call FunctionCall) Value ***REMOVED***
	toString := toMethod(r.getVStr(call.This, "toString"))
	return toString(FunctionCall***REMOVED***This: call.This***REMOVED***)
***REMOVED***

func (r *Runtime) objectproto_getProto(call FunctionCall) Value ***REMOVED***
	proto := call.This.ToObject(r).self.proto()
	if proto != nil ***REMOVED***
		return proto
	***REMOVED***
	return _null
***REMOVED***

func (r *Runtime) objectproto_setProto(call FunctionCall) Value ***REMOVED***
	o := call.This
	r.checkObjectCoercible(o)
	proto := r.toProto(call.Argument(0))
	if o, ok := o.(*Object); ok ***REMOVED***
		o.self.setProto(proto, true)
	***REMOVED***

	return _undefined
***REMOVED***

func (r *Runtime) objectproto_valueOf(call FunctionCall) Value ***REMOVED***
	return call.This.ToObject(r)
***REMOVED***

func (r *Runtime) object_assign(call FunctionCall) Value ***REMOVED***
	to := call.Argument(0).ToObject(r)
	if len(call.Arguments) > 1 ***REMOVED***
		for _, arg := range call.Arguments[1:] ***REMOVED***
			if arg != _undefined && arg != _null ***REMOVED***
				source := arg.ToObject(r)
				for _, key := range source.self.ownPropertyKeys(true, nil) ***REMOVED***
					p := source.getOwnProp(key)
					if p == nil ***REMOVED***
						continue
					***REMOVED***
					if v, ok := p.(*valueProperty); ok ***REMOVED***
						if !v.enumerable ***REMOVED***
							continue
						***REMOVED***
						p = v.get(source)
					***REMOVED***
					to.setOwn(key, p, true)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return to
***REMOVED***

func (r *Runtime) object_is(call FunctionCall) Value ***REMOVED***
	return r.toBoolean(call.Argument(0).SameAs(call.Argument(1)))
***REMOVED***

func (r *Runtime) toProto(proto Value) *Object ***REMOVED***
	if proto != _null ***REMOVED***
		if obj, ok := proto.(*Object); ok ***REMOVED***
			return obj
		***REMOVED*** else ***REMOVED***
			panic(r.NewTypeError("Object prototype may only be an Object or null: %s", proto))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *Runtime) object_setPrototypeOf(call FunctionCall) Value ***REMOVED***
	o := call.Argument(0)
	r.checkObjectCoercible(o)
	proto := r.toProto(call.Argument(1))
	if o, ok := o.(*Object); ok ***REMOVED***
		o.self.setProto(proto, true)
	***REMOVED***

	return o
***REMOVED***

func (r *Runtime) initObject() ***REMOVED***
	o := r.global.ObjectPrototype.self
	o._putProp("toString", r.newNativeFunc(r.objectproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("toLocaleString", r.newNativeFunc(r.objectproto_toLocaleString, nil, "toLocaleString", nil, 0), true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.objectproto_valueOf, nil, "valueOf", nil, 0), true, false, true)
	o._putProp("hasOwnProperty", r.newNativeFunc(r.objectproto_hasOwnProperty, nil, "hasOwnProperty", nil, 1), true, false, true)
	o._putProp("isPrototypeOf", r.newNativeFunc(r.objectproto_isPrototypeOf, nil, "isPrototypeOf", nil, 1), true, false, true)
	o._putProp("propertyIsEnumerable", r.newNativeFunc(r.objectproto_propertyIsEnumerable, nil, "propertyIsEnumerable", nil, 1), true, false, true)
	o.defineOwnPropertyStr(__proto__, PropertyDescriptor***REMOVED***
		Getter:       r.newNativeFunc(r.objectproto_getProto, nil, "get __proto__", nil, 0),
		Setter:       r.newNativeFunc(r.objectproto_setProto, nil, "set __proto__", nil, 1),
		Configurable: FLAG_TRUE,
	***REMOVED***, true)

	r.global.Object = r.newNativeFuncConstruct(r.builtin_Object, classObject, r.global.ObjectPrototype, 1)
	o = r.global.Object.self
	o._putProp("assign", r.newNativeFunc(r.object_assign, nil, "assign", nil, 2), true, false, true)
	o._putProp("defineProperty", r.newNativeFunc(r.object_defineProperty, nil, "defineProperty", nil, 3), true, false, true)
	o._putProp("defineProperties", r.newNativeFunc(r.object_defineProperties, nil, "defineProperties", nil, 2), true, false, true)
	o._putProp("getOwnPropertyDescriptor", r.newNativeFunc(r.object_getOwnPropertyDescriptor, nil, "getOwnPropertyDescriptor", nil, 2), true, false, true)
	o._putProp("getPrototypeOf", r.newNativeFunc(r.object_getPrototypeOf, nil, "getPrototypeOf", nil, 1), true, false, true)
	o._putProp("is", r.newNativeFunc(r.object_is, nil, "is", nil, 2), true, false, true)
	o._putProp("getOwnPropertyNames", r.newNativeFunc(r.object_getOwnPropertyNames, nil, "getOwnPropertyNames", nil, 1), true, false, true)
	o._putProp("getOwnPropertySymbols", r.newNativeFunc(r.object_getOwnPropertySymbols, nil, "getOwnPropertySymbols", nil, 1), true, false, true)
	o._putProp("create", r.newNativeFunc(r.object_create, nil, "create", nil, 2), true, false, true)
	o._putProp("seal", r.newNativeFunc(r.object_seal, nil, "seal", nil, 1), true, false, true)
	o._putProp("freeze", r.newNativeFunc(r.object_freeze, nil, "freeze", nil, 1), true, false, true)
	o._putProp("preventExtensions", r.newNativeFunc(r.object_preventExtensions, nil, "preventExtensions", nil, 1), true, false, true)
	o._putProp("isSealed", r.newNativeFunc(r.object_isSealed, nil, "isSealed", nil, 1), true, false, true)
	o._putProp("isFrozen", r.newNativeFunc(r.object_isFrozen, nil, "isFrozen", nil, 1), true, false, true)
	o._putProp("isExtensible", r.newNativeFunc(r.object_isExtensible, nil, "isExtensible", nil, 1), true, false, true)
	o._putProp("keys", r.newNativeFunc(r.object_keys, nil, "keys", nil, 1), true, false, true)
	o._putProp("setPrototypeOf", r.newNativeFunc(r.object_setPrototypeOf, nil, "setPrototypeOf", nil, 2), true, false, true)

	r.addToGlobal("Object", r.global.Object)
***REMOVED***
