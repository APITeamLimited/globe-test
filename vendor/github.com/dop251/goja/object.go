package goja

import "reflect"

const (
	classObject   = "Object"
	classArray    = "Array"
	classFunction = "Function"
	classNumber   = "Number"
	classString   = "String"
	classBoolean  = "Boolean"
	classError    = "Error"
	classRegExp   = "RegExp"
	classDate     = "Date"
)

type Object struct ***REMOVED***
	runtime *Runtime
	self    objectImpl
***REMOVED***

type iterNextFunc func() (propIterItem, iterNextFunc)

type propertyDescr struct ***REMOVED***
	Value Value

	Writable, Configurable, Enumerable Flag

	Getter, Setter Value
***REMOVED***

type objectImpl interface ***REMOVED***
	sortable
	className() string
	get(Value) Value
	getProp(Value) Value
	getPropStr(string) Value
	getStr(string) Value
	getOwnProp(string) Value
	put(Value, Value, bool)
	putStr(string, Value, bool)
	hasProperty(Value) bool
	hasPropertyStr(string) bool
	hasOwnProperty(Value) bool
	hasOwnPropertyStr(string) bool
	_putProp(name string, value Value, writable, enumerable, configurable bool) Value
	defineOwnProperty(name Value, descr propertyDescr, throw bool) bool
	toPrimitiveNumber() Value
	toPrimitiveString() Value
	toPrimitive() Value
	assertCallable() (call func(FunctionCall) Value, ok bool)
	deleteStr(name string, throw bool) bool
	delete(name Value, throw bool) bool
	proto() *Object
	hasInstance(v Value) bool
	isExtensible() bool
	preventExtensions()
	enumerate(all, recusrive bool) iterNextFunc
	_enumerate(recursive bool) iterNextFunc
	export() interface***REMOVED******REMOVED***
	exportType() reflect.Type
	equal(objectImpl) bool
***REMOVED***

type baseObject struct ***REMOVED***
	class      string
	val        *Object
	prototype  *Object
	extensible bool

	values    map[string]Value
	propNames []string
***REMOVED***

type primitiveValueObject struct ***REMOVED***
	baseObject
	pValue Value
***REMOVED***

func (o *primitiveValueObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	return o.pValue.Export()
***REMOVED***

func (o *primitiveValueObject) exportType() reflect.Type ***REMOVED***
	return o.pValue.ExportType()
***REMOVED***

type FunctionCall struct ***REMOVED***
	This      Value
	Arguments []Value
***REMOVED***

type ConstructorCall struct ***REMOVED***
	This      *Object
	Arguments []Value
***REMOVED***

func (f FunctionCall) Argument(idx int) Value ***REMOVED***
	if idx < len(f.Arguments) ***REMOVED***
		return f.Arguments[idx]
	***REMOVED***
	return _undefined
***REMOVED***

func (f ConstructorCall) Argument(idx int) Value ***REMOVED***
	if idx < len(f.Arguments) ***REMOVED***
		return f.Arguments[idx]
	***REMOVED***
	return _undefined
***REMOVED***

func (o *baseObject) init() ***REMOVED***
	o.values = make(map[string]Value)
***REMOVED***

func (o *baseObject) className() string ***REMOVED***
	return o.class
***REMOVED***

func (o *baseObject) getPropStr(name string) Value ***REMOVED***
	if val := o.getOwnProp(name); val != nil ***REMOVED***
		return val
	***REMOVED***
	if o.prototype != nil ***REMOVED***
		return o.prototype.self.getPropStr(name)
	***REMOVED***
	return nil
***REMOVED***

func (o *baseObject) getProp(n Value) Value ***REMOVED***
	return o.val.self.getPropStr(n.String())
***REMOVED***

func (o *baseObject) hasProperty(n Value) bool ***REMOVED***
	return o.val.self.getProp(n) != nil
***REMOVED***

func (o *baseObject) hasPropertyStr(name string) bool ***REMOVED***
	return o.val.self.getPropStr(name) != nil
***REMOVED***

func (o *baseObject) _getStr(name string) Value ***REMOVED***
	p := o.getOwnProp(name)

	if p == nil && o.prototype != nil ***REMOVED***
		p = o.prototype.self.getPropStr(name)
	***REMOVED***

	if p, ok := p.(*valueProperty); ok ***REMOVED***
		return p.get(o.val)
	***REMOVED***

	return p
***REMOVED***

func (o *baseObject) getStr(name string) Value ***REMOVED***
	p := o.val.self.getPropStr(name)
	if p, ok := p.(*valueProperty); ok ***REMOVED***
		return p.get(o.val)
	***REMOVED***

	return p
***REMOVED***

func (o *baseObject) get(n Value) Value ***REMOVED***
	return o.getStr(n.String())
***REMOVED***

func (o *baseObject) checkDeleteProp(name string, prop *valueProperty, throw bool) bool ***REMOVED***
	if !prop.configurable ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Cannot delete property '%s' of %s", name, o.val.ToString())
		return false
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) checkDelete(name string, val Value, throw bool) bool ***REMOVED***
	if val, ok := val.(*valueProperty); ok ***REMOVED***
		return o.checkDeleteProp(name, val, throw)
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) _delete(name string) ***REMOVED***
	delete(o.values, name)
	for i, n := range o.propNames ***REMOVED***
		if n == name ***REMOVED***
			copy(o.propNames[i:], o.propNames[i+1:])
			o.propNames = o.propNames[:len(o.propNames)-1]
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (o *baseObject) deleteStr(name string, throw bool) bool ***REMOVED***
	if val, exists := o.values[name]; exists ***REMOVED***
		if !o.checkDelete(name, val, throw) ***REMOVED***
			return false
		***REMOVED***
		o._delete(name)
		return true
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) delete(n Value, throw bool) bool ***REMOVED***
	return o.deleteStr(n.String(), throw)
***REMOVED***

func (o *baseObject) put(n Value, val Value, throw bool) ***REMOVED***
	o.putStr(n.String(), val, throw)
***REMOVED***

func (o *baseObject) getOwnProp(name string) Value ***REMOVED***
	v := o.values[name]
	if v == nil && name == __proto__ ***REMOVED***
		return o.prototype
	***REMOVED***
	return v
***REMOVED***

func (o *baseObject) putStr(name string, val Value, throw bool) ***REMOVED***
	if v, exists := o.values[name]; exists ***REMOVED***
		if prop, ok := v.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				o.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
				return
			***REMOVED***
			prop.set(o.val, val)
			return
		***REMOVED***
		o.values[name] = val
		return
	***REMOVED***

	if name == __proto__ ***REMOVED***
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "%s is not extensible", o.val)
			return
		***REMOVED***
		if val == _undefined || val == _null ***REMOVED***
			o.prototype = nil
			return
		***REMOVED*** else ***REMOVED***
			if val, ok := val.(*Object); ok ***REMOVED***
				o.prototype = val
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	var pprop Value
	if proto := o.prototype; proto != nil ***REMOVED***
		pprop = proto.self.getPropStr(name)
	***REMOVED***

	if pprop != nil ***REMOVED***
		if prop, ok := pprop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				o.val.runtime.typeErrorResult(throw)
				return
			***REMOVED***
			if prop.accessor ***REMOVED***
				prop.set(o.val, val)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw)
			return
		***REMOVED***
	***REMOVED***

	o.values[name] = val
	o.propNames = append(o.propNames, name)
***REMOVED***

func (o *baseObject) hasOwnProperty(n Value) bool ***REMOVED***
	v := o.values[n.String()]
	return v != nil
***REMOVED***

func (o *baseObject) hasOwnPropertyStr(name string) bool ***REMOVED***
	v := o.values[name]
	return v != nil
***REMOVED***

func (o *baseObject) _defineOwnProperty(name, existingValue Value, descr propertyDescr, throw bool) (val Value, ok bool) ***REMOVED***

	getterObj, _ := descr.Getter.(*Object)
	setterObj, _ := descr.Setter.(*Object)

	var existing *valueProperty

	if existingValue == nil ***REMOVED***
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw)
			return nil, false
		***REMOVED***
		existing = &valueProperty***REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		if existing, ok = existingValue.(*valueProperty); !ok ***REMOVED***
			existing = &valueProperty***REMOVED***
				writable:     true,
				enumerable:   true,
				configurable: true,
				value:        existingValue,
			***REMOVED***
		***REMOVED***

		if !existing.configurable ***REMOVED***
			if descr.Configurable == FLAG_TRUE ***REMOVED***
				goto Reject
			***REMOVED***
			if descr.Enumerable != FLAG_NOT_SET && descr.Enumerable.Bool() != existing.enumerable ***REMOVED***
				goto Reject
			***REMOVED***
		***REMOVED***
		if existing.accessor && descr.Value != nil || !existing.accessor && (getterObj != nil || setterObj != nil) ***REMOVED***
			if !existing.configurable ***REMOVED***
				goto Reject
			***REMOVED***
		***REMOVED*** else if !existing.accessor ***REMOVED***
			if !existing.configurable ***REMOVED***
				if !existing.writable ***REMOVED***
					if descr.Writable == FLAG_TRUE ***REMOVED***
						goto Reject
					***REMOVED***
					if descr.Value != nil && !descr.Value.SameAs(existing.value) ***REMOVED***
						goto Reject
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if !existing.configurable ***REMOVED***
				if descr.Getter != nil && existing.getterFunc != getterObj || descr.Setter != nil && existing.setterFunc != setterObj ***REMOVED***
					goto Reject
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if descr.Writable == FLAG_TRUE && descr.Enumerable == FLAG_TRUE && descr.Configurable == FLAG_TRUE && descr.Value != nil ***REMOVED***
		return descr.Value, true
	***REMOVED***

	if descr.Writable != FLAG_NOT_SET ***REMOVED***
		existing.writable = descr.Writable.Bool()
	***REMOVED***
	if descr.Enumerable != FLAG_NOT_SET ***REMOVED***
		existing.enumerable = descr.Enumerable.Bool()
	***REMOVED***
	if descr.Configurable != FLAG_NOT_SET ***REMOVED***
		existing.configurable = descr.Configurable.Bool()
	***REMOVED***

	if descr.Value != nil ***REMOVED***
		existing.value = descr.Value
		existing.getterFunc = nil
		existing.setterFunc = nil
	***REMOVED***

	if descr.Value != nil || descr.Writable != FLAG_NOT_SET ***REMOVED***
		existing.accessor = false
	***REMOVED***

	if descr.Getter != nil ***REMOVED***
		existing.getterFunc = propGetter(o.val, descr.Getter, o.val.runtime)
		existing.value = nil
		existing.accessor = true
	***REMOVED***

	if descr.Setter != nil ***REMOVED***
		existing.setterFunc = propSetter(o.val, descr.Setter, o.val.runtime)
		existing.value = nil
		existing.accessor = true
	***REMOVED***

	if !existing.accessor && existing.value == nil ***REMOVED***
		existing.value = _undefined
	***REMOVED***

	return existing, true

Reject:
	o.val.runtime.typeErrorResult(throw, "Cannot redefine property: %s", name.ToString())
	return nil, false

***REMOVED***

func (o *baseObject) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool ***REMOVED***
	name := n.String()
	existingVal := o.values[name]
	if v, ok := o._defineOwnProperty(n, existingVal, descr, throw); ok ***REMOVED***
		o.values[name] = v
		if existingVal == nil ***REMOVED***
			o.propNames = append(o.propNames, name)
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *baseObject) _put(name string, v Value) ***REMOVED***
	if _, exists := o.values[name]; !exists ***REMOVED***
		o.propNames = append(o.propNames, name)
	***REMOVED***

	o.values[name] = v
***REMOVED***

func (o *baseObject) _putProp(name string, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	if writable && enumerable && configurable ***REMOVED***
		o._put(name, value)
		return value
	***REMOVED*** else ***REMOVED***
		p := &valueProperty***REMOVED***
			value:        value,
			writable:     writable,
			enumerable:   enumerable,
			configurable: configurable,
		***REMOVED***
		o._put(name, p)
		return p
	***REMOVED***
***REMOVED***

func (o *baseObject) tryPrimitive(methodName string) Value ***REMOVED***
	if method, ok := o.getStr(methodName).(*Object); ok ***REMOVED***
		if call, ok := method.self.assertCallable(); ok ***REMOVED***
			v := call(FunctionCall***REMOVED***
				This: o.val,
			***REMOVED***)
			if _, fail := v.(*Object); !fail ***REMOVED***
				return v
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o *baseObject) toPrimitiveNumber() Value ***REMOVED***
	if v := o.tryPrimitive("valueOf"); v != nil ***REMOVED***
		return v
	***REMOVED***

	if v := o.tryPrimitive("toString"); v != nil ***REMOVED***
		return v
	***REMOVED***

	o.val.runtime.typeErrorResult(true, "Could not convert %v to primitive", o)
	return nil
***REMOVED***

func (o *baseObject) toPrimitiveString() Value ***REMOVED***
	if v := o.tryPrimitive("toString"); v != nil ***REMOVED***
		return v
	***REMOVED***

	if v := o.tryPrimitive("valueOf"); v != nil ***REMOVED***
		return v
	***REMOVED***

	o.val.runtime.typeErrorResult(true, "Could not convert %v to primitive", o)
	return nil
***REMOVED***

func (o *baseObject) toPrimitive() Value ***REMOVED***
	return o.toPrimitiveNumber()
***REMOVED***

func (o *baseObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	return nil, false
***REMOVED***

func (o *baseObject) proto() *Object ***REMOVED***
	return o.prototype
***REMOVED***

func (o *baseObject) isExtensible() bool ***REMOVED***
	return o.extensible
***REMOVED***

func (o *baseObject) preventExtensions() ***REMOVED***
	o.extensible = false
***REMOVED***

func (o *baseObject) sortLen() int64 ***REMOVED***
	return toLength(o.val.self.getStr("length"))
***REMOVED***

func (o *baseObject) sortGet(i int64) Value ***REMOVED***
	return o.val.self.get(intToValue(i))
***REMOVED***

func (o *baseObject) swap(i, j int64) ***REMOVED***
	ii := intToValue(i)
	jj := intToValue(j)

	x := o.val.self.get(ii)
	y := o.val.self.get(jj)

	o.val.self.put(ii, y, false)
	o.val.self.put(jj, x, false)
***REMOVED***

func (o *baseObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	m := make(map[string]interface***REMOVED******REMOVED***)

	for item, f := o.enumerate(false, false)(); f != nil; item, f = f() ***REMOVED***
		v := item.value
		if v == nil ***REMOVED***
			v = o.getStr(item.name)
		***REMOVED***
		if v != nil ***REMOVED***
			m[item.name] = v.Export()
		***REMOVED*** else ***REMOVED***
			m[item.name] = nil
		***REMOVED***
	***REMOVED***
	return m
***REMOVED***

func (o *baseObject) exportType() reflect.Type ***REMOVED***
	return reflectTypeMap
***REMOVED***

type enumerableFlag int

const (
	_ENUM_UNKNOWN enumerableFlag = iota
	_ENUM_FALSE
	_ENUM_TRUE
)

type propIterItem struct ***REMOVED***
	name       string
	value      Value // set only when enumerable == _ENUM_UNKNOWN
	enumerable enumerableFlag
***REMOVED***

type objectPropIter struct ***REMOVED***
	o         *baseObject
	propNames []string
	recursive bool
	idx       int
***REMOVED***

type propFilterIter struct ***REMOVED***
	wrapped iterNextFunc
	all     bool
	seen    map[string]bool
***REMOVED***

func (i *propFilterIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for ***REMOVED***
		var item propIterItem
		item, i.wrapped = i.wrapped()
		if i.wrapped == nil ***REMOVED***
			return propIterItem***REMOVED******REMOVED***, nil
		***REMOVED***

		if !i.seen[item.name] ***REMOVED***
			i.seen[item.name] = true
			if !i.all ***REMOVED***
				if item.enumerable == _ENUM_FALSE ***REMOVED***
					continue
				***REMOVED***
				if item.enumerable == _ENUM_UNKNOWN ***REMOVED***
					if prop, ok := item.value.(*valueProperty); ok ***REMOVED***
						if !prop.enumerable ***REMOVED***
							continue
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
			return item, i.next
		***REMOVED***
	***REMOVED***
***REMOVED***

func (i *objectPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.propNames) ***REMOVED***
		name := i.propNames[i.idx]
		i.idx++
		prop := i.o.values[name]
		if prop != nil ***REMOVED***
			return propIterItem***REMOVED***name: name, value: prop***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	if i.recursive && i.o.prototype != nil ***REMOVED***
		return i.o.prototype.self._enumerate(i.recursive)()
	***REMOVED***
	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *baseObject) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	propNames := make([]string, len(o.propNames))
	copy(propNames, o.propNames)
	return (&objectPropIter***REMOVED***
		o:         o,
		propNames: propNames,
		recursive: recursive,
	***REMOVED***).next
***REMOVED***

func (o *baseObject) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next
***REMOVED***

func (o *baseObject) equal(other objectImpl) bool ***REMOVED***
	// Rely on parent reference comparison
	return false
***REMOVED***

func (o *baseObject) hasInstance(v Value) bool ***REMOVED***
	o.val.runtime.typeErrorResult(true, "Expecting a function in instanceof check, but got %s", o.val.ToString())
	panic("Unreachable")
***REMOVED***
