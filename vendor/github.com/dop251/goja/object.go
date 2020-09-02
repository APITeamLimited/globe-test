package goja

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"

	"github.com/dop251/goja/unistring"
)

const (
	classObject   = "Object"
	classArray    = "Array"
	classWeakSet  = "WeakSet"
	classWeakMap  = "WeakMap"
	classMap      = "Map"
	classMath     = "Math"
	classSet      = "Set"
	classFunction = "Function"
	classNumber   = "Number"
	classString   = "String"
	classBoolean  = "Boolean"
	classError    = "Error"
	classRegExp   = "RegExp"
	classDate     = "Date"
	classJSON     = "JSON"

	classArrayIterator  = "Array Iterator"
	classMapIterator    = "Map Iterator"
	classSetIterator    = "Set Iterator"
	classStringIterator = "String Iterator"
)

var (
	hintDefault Value = asciiString("default")
	hintNumber  Value = asciiString("number")
	hintString  Value = asciiString("string")
)

type weakCollection interface ***REMOVED***
	removeId(uint64)
***REMOVED***

type weakCollections struct ***REMOVED***
	objId uint64
	colls []weakCollection
***REMOVED***

func (r *weakCollections) add(c weakCollection) ***REMOVED***
	for _, ec := range r.colls ***REMOVED***
		if ec == c ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	r.colls = append(r.colls, c)
***REMOVED***

func (r *weakCollections) id() uint64 ***REMOVED***
	return r.objId
***REMOVED***

func (r *weakCollections) remove(c weakCollection) ***REMOVED***
	if cap(r.colls) > 16 && cap(r.colls)>>2 > len(r.colls) ***REMOVED***
		// shrink
		colls := make([]weakCollection, 0, len(r.colls))
		for _, coll := range r.colls ***REMOVED***
			if coll != c ***REMOVED***
				colls = append(colls, coll)
			***REMOVED***
		***REMOVED***
		r.colls = colls
	***REMOVED*** else ***REMOVED***
		for i, coll := range r.colls ***REMOVED***
			if coll == c ***REMOVED***
				l := len(r.colls) - 1
				r.colls[i] = r.colls[l]
				r.colls[l] = nil
				r.colls = r.colls[:l]
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func finalizeObjectWeakRefs(r *weakCollections) ***REMOVED***
	id := r.id()
	for _, c := range r.colls ***REMOVED***
		c.removeId(id)
	***REMOVED***
	r.colls = nil
***REMOVED***

type Object struct ***REMOVED***
	id      uint64
	runtime *Runtime
	self    objectImpl

	// Contains references to all weak collections that contain this Object.
	// weakColls has a finalizer that removes the Object's id from all weak collections.
	// The id is the weakColls pointer value converted to uintptr.
	// Note, cannot set the finalizer on the *Object itself because it's a part of a
	// reference cycle.
	weakColls *weakCollections
***REMOVED***

type iterNextFunc func() (propIterItem, iterNextFunc)

type PropertyDescriptor struct ***REMOVED***
	jsDescriptor *Object

	Value Value

	Writable, Configurable, Enumerable Flag

	Getter, Setter Value
***REMOVED***

func (p *PropertyDescriptor) Empty() bool ***REMOVED***
	var empty PropertyDescriptor
	return *p == empty
***REMOVED***

func (p *PropertyDescriptor) toValue(r *Runtime) Value ***REMOVED***
	if p.jsDescriptor != nil ***REMOVED***
		return p.jsDescriptor
	***REMOVED***

	o := r.NewObject()
	s := o.self

	if p.Value != nil ***REMOVED***
		s._putProp("value", p.Value, true, true, true)
	***REMOVED***

	if p.Writable != FLAG_NOT_SET ***REMOVED***
		s._putProp("writable", valueBool(p.Writable.Bool()), true, true, true)
	***REMOVED***

	if p.Enumerable != FLAG_NOT_SET ***REMOVED***
		s._putProp("enumerable", valueBool(p.Enumerable.Bool()), true, true, true)
	***REMOVED***

	if p.Configurable != FLAG_NOT_SET ***REMOVED***
		s._putProp("configurable", valueBool(p.Configurable.Bool()), true, true, true)
	***REMOVED***

	if p.Getter != nil ***REMOVED***
		s._putProp("get", p.Getter, true, true, true)
	***REMOVED***
	if p.Setter != nil ***REMOVED***
		s._putProp("set", p.Setter, true, true, true)
	***REMOVED***

	return o
***REMOVED***

func (p *PropertyDescriptor) complete() ***REMOVED***
	if p.Getter == nil && p.Setter == nil ***REMOVED***
		if p.Value == nil ***REMOVED***
			p.Value = _undefined
		***REMOVED***
		if p.Writable == FLAG_NOT_SET ***REMOVED***
			p.Writable = FLAG_FALSE
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if p.Getter == nil ***REMOVED***
			p.Getter = _undefined
		***REMOVED***
		if p.Setter == nil ***REMOVED***
			p.Setter = _undefined
		***REMOVED***
	***REMOVED***
	if p.Enumerable == FLAG_NOT_SET ***REMOVED***
		p.Enumerable = FLAG_FALSE
	***REMOVED***
	if p.Configurable == FLAG_NOT_SET ***REMOVED***
		p.Configurable = FLAG_FALSE
	***REMOVED***
***REMOVED***

type objectExportCacheItem map[reflect.Type]interface***REMOVED******REMOVED***

type objectExportCtx struct ***REMOVED***
	cache map[objectImpl]interface***REMOVED******REMOVED***
***REMOVED***

type objectImpl interface ***REMOVED***
	sortable
	className() string
	getStr(p unistring.String, receiver Value) Value
	getIdx(p valueInt, receiver Value) Value
	getSym(p *valueSymbol, receiver Value) Value

	getOwnPropStr(unistring.String) Value
	getOwnPropIdx(valueInt) Value
	getOwnPropSym(*valueSymbol) Value

	setOwnStr(p unistring.String, v Value, throw bool) bool
	setOwnIdx(p valueInt, v Value, throw bool) bool
	setOwnSym(p *valueSymbol, v Value, throw bool) bool

	setForeignStr(p unistring.String, v, receiver Value, throw bool) (res bool, handled bool)
	setForeignIdx(p valueInt, v, receiver Value, throw bool) (res bool, handled bool)
	setForeignSym(p *valueSymbol, v, receiver Value, throw bool) (res bool, handled bool)

	hasPropertyStr(unistring.String) bool
	hasPropertyIdx(idx valueInt) bool
	hasPropertySym(s *valueSymbol) bool

	hasOwnPropertyStr(unistring.String) bool
	hasOwnPropertyIdx(valueInt) bool
	hasOwnPropertySym(s *valueSymbol) bool

	defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool
	defineOwnPropertyIdx(name valueInt, desc PropertyDescriptor, throw bool) bool
	defineOwnPropertySym(name *valueSymbol, desc PropertyDescriptor, throw bool) bool

	deleteStr(name unistring.String, throw bool) bool
	deleteIdx(idx valueInt, throw bool) bool
	deleteSym(s *valueSymbol, throw bool) bool

	toPrimitiveNumber() Value
	toPrimitiveString() Value
	toPrimitive() Value
	assertCallable() (call func(FunctionCall) Value, ok bool)
	assertConstructor() func(args []Value, newTarget *Object) *Object
	proto() *Object
	setProto(proto *Object, throw bool) bool
	hasInstance(v Value) bool
	isExtensible() bool
	preventExtensions(throw bool) bool
	enumerate() iterNextFunc
	enumerateUnfiltered() iterNextFunc
	export(ctx *objectExportCtx) interface***REMOVED******REMOVED***
	exportType() reflect.Type
	equal(objectImpl) bool
	ownKeys(all bool, accum []Value) []Value
	ownSymbols(all bool, accum []Value) []Value
	ownPropertyKeys(all bool, accum []Value) []Value

	_putProp(name unistring.String, value Value, writable, enumerable, configurable bool) Value
	_putSym(s *valueSymbol, prop Value)
***REMOVED***

type baseObject struct ***REMOVED***
	class      string
	val        *Object
	prototype  *Object
	extensible bool

	values    map[unistring.String]Value
	propNames []unistring.String

	lastSortedPropLen, idxPropCount int

	symValues *orderedMap
***REMOVED***

type guardedObject struct ***REMOVED***
	baseObject
	guardedProps map[unistring.String]struct***REMOVED******REMOVED***
***REMOVED***

type primitiveValueObject struct ***REMOVED***
	baseObject
	pValue Value
***REMOVED***

func (o *primitiveValueObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
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
	o.values = make(map[unistring.String]Value)
***REMOVED***

func (o *baseObject) className() string ***REMOVED***
	return o.class
***REMOVED***

func (o *baseObject) hasPropertyStr(name unistring.String) bool ***REMOVED***
	if o.val.self.hasOwnPropertyStr(name) ***REMOVED***
		return true
	***REMOVED***
	if o.prototype != nil ***REMOVED***
		return o.prototype.self.hasPropertyStr(name)
	***REMOVED***
	return false
***REMOVED***

func (o *baseObject) hasPropertyIdx(idx valueInt) bool ***REMOVED***
	return o.val.self.hasPropertyStr(idx.string())
***REMOVED***

func (o *baseObject) hasPropertySym(s *valueSymbol) bool ***REMOVED***
	if o.hasOwnPropertySym(s) ***REMOVED***
		return true
	***REMOVED***
	if o.prototype != nil ***REMOVED***
		return o.prototype.self.hasPropertySym(s)
	***REMOVED***
	return false
***REMOVED***

func (o *baseObject) getWithOwnProp(prop, p, receiver Value) Value ***REMOVED***
	if prop == nil && o.prototype != nil ***REMOVED***
		if receiver == nil ***REMOVED***
			return o.prototype.get(p, o.val)
		***REMOVED***
		return o.prototype.get(p, receiver)
	***REMOVED***
	if prop, ok := prop.(*valueProperty); ok ***REMOVED***
		if receiver == nil ***REMOVED***
			return prop.get(o.val)
		***REMOVED***
		return prop.get(receiver)
	***REMOVED***
	return prop
***REMOVED***

func (o *baseObject) getStrWithOwnProp(prop Value, name unistring.String, receiver Value) Value ***REMOVED***
	if prop == nil && o.prototype != nil ***REMOVED***
		if receiver == nil ***REMOVED***
			return o.prototype.self.getStr(name, o.val)
		***REMOVED***
		return o.prototype.self.getStr(name, receiver)
	***REMOVED***
	if prop, ok := prop.(*valueProperty); ok ***REMOVED***
		if receiver == nil ***REMOVED***
			return prop.get(o.val)
		***REMOVED***
		return prop.get(receiver)
	***REMOVED***
	return prop
***REMOVED***

func (o *baseObject) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	return o.val.self.getStr(idx.string(), receiver)
***REMOVED***

func (o *baseObject) getSym(s *valueSymbol, receiver Value) Value ***REMOVED***
	return o.getWithOwnProp(o.getOwnPropSym(s), s, receiver)
***REMOVED***

func (o *baseObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	prop := o.values[name]
	if prop == nil ***REMOVED***
		if o.prototype != nil ***REMOVED***
			if receiver == nil ***REMOVED***
				return o.prototype.self.getStr(name, o.val)
			***REMOVED***
			return o.prototype.self.getStr(name, receiver)
		***REMOVED***
	***REMOVED***
	if prop, ok := prop.(*valueProperty); ok ***REMOVED***
		if receiver == nil ***REMOVED***
			return prop.get(o.val)
		***REMOVED***
		return prop.get(receiver)
	***REMOVED***
	return prop
***REMOVED***

func (o *baseObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	return o.val.self.getOwnPropStr(idx.string())
***REMOVED***

func (o *baseObject) getOwnPropSym(s *valueSymbol) Value ***REMOVED***
	if o.symValues != nil ***REMOVED***
		return o.symValues.get(s)
	***REMOVED***
	return nil
***REMOVED***

func (o *baseObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	return o.values[name]
***REMOVED***

func (o *baseObject) checkDeleteProp(name unistring.String, prop *valueProperty, throw bool) bool ***REMOVED***
	if !prop.configurable ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Cannot delete property '%s' of %s", name, o.val.toString())
		return false
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) checkDelete(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if val, ok := val.(*valueProperty); ok ***REMOVED***
		return o.checkDeleteProp(name, val, throw)
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) _delete(name unistring.String) ***REMOVED***
	delete(o.values, name)
	for i, n := range o.propNames ***REMOVED***
		if n == name ***REMOVED***
			copy(o.propNames[i:], o.propNames[i+1:])
			o.propNames = o.propNames[:len(o.propNames)-1]
			if i < o.lastSortedPropLen ***REMOVED***
				o.lastSortedPropLen--
				if i < o.idxPropCount ***REMOVED***
					o.idxPropCount--
				***REMOVED***
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

func (o *baseObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	return o.val.self.deleteStr(idx.string(), throw)
***REMOVED***

func (o *baseObject) deleteSym(s *valueSymbol, throw bool) bool ***REMOVED***
	if o.symValues != nil ***REMOVED***
		if val := o.symValues.get(s); val != nil ***REMOVED***
			if !o.checkDelete(s.desc.string(), val, throw) ***REMOVED***
				return false
			***REMOVED***
			o.symValues.remove(s)
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if val, exists := o.values[name]; exists ***REMOVED***
		if !o.checkDelete(name, val, throw) ***REMOVED***
			return false
		***REMOVED***
		o._delete(name)
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) setProto(proto *Object, throw bool) bool ***REMOVED***
	current := o.prototype
	if current.SameAs(proto) ***REMOVED***
		return true
	***REMOVED***
	if !o.extensible ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "%s is not extensible", o.val)
		return false
	***REMOVED***
	for p := proto; p != nil; ***REMOVED***
		if p.SameAs(o.val) ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cyclic __proto__ value")
			return false
		***REMOVED***
		p = p.self.proto()
	***REMOVED***
	o.prototype = proto
	return true
***REMOVED***

func (o *baseObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	ownDesc := o.values[name]
	if ownDesc == nil ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			// we know it's foreign because prototype loops are not allowed
			if res, handled := proto.self.setForeignStr(name, val, o.val, throw); handled ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		// new property
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot add property %s, object is not extensible", name)
			return false
		***REMOVED*** else ***REMOVED***
			o.values[name] = val
			o.propNames = append(o.propNames, name)
		***REMOVED***
		return true
	***REMOVED***
	if prop, ok := ownDesc.(*valueProperty); ok ***REMOVED***
		if !prop.isWritable() ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
			return false
		***REMOVED*** else ***REMOVED***
			prop.set(o.val, val)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		o.values[name] = val
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	return o.val.self.setOwnStr(idx.string(), val, throw)
***REMOVED***

func (o *baseObject) setOwnSym(name *valueSymbol, val Value, throw bool) bool ***REMOVED***
	var ownDesc Value
	if o.symValues != nil ***REMOVED***
		ownDesc = o.symValues.get(name)
	***REMOVED***
	if ownDesc == nil ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			// we know it's foreign because prototype loops are not allowed
			if res, handled := proto.self.setForeignSym(name, val, o.val, throw); handled ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		// new property
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot add property %s, object is not extensible", name)
			return false
		***REMOVED*** else ***REMOVED***
			if o.symValues == nil ***REMOVED***
				o.symValues = newOrderedMap(nil)
			***REMOVED***
			o.symValues.set(name, val)
		***REMOVED***
		return true
	***REMOVED***
	if prop, ok := ownDesc.(*valueProperty); ok ***REMOVED***
		if !prop.isWritable() ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
			return false
		***REMOVED*** else ***REMOVED***
			prop.set(o.val, val)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		o.symValues.set(name, val)
	***REMOVED***
	return true
***REMOVED***

func (o *baseObject) _setForeignStr(name unistring.String, prop, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	if prop != nil ***REMOVED***
		if prop, ok := prop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				o.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
				return false, true
			***REMOVED***
			if prop.setterFunc != nil ***REMOVED***
				prop.set(receiver, val)
				return true, true
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			if receiver != proto ***REMOVED***
				return proto.self.setForeignStr(name, val, receiver, throw)
			***REMOVED***
			return proto.self.setOwnStr(name, val, throw), true
		***REMOVED***
	***REMOVED***
	return false, false
***REMOVED***

func (o *baseObject) _setForeignIdx(idx valueInt, prop, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	if prop != nil ***REMOVED***
		if prop, ok := prop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				o.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%d'", idx)
				return false, true
			***REMOVED***
			if prop.setterFunc != nil ***REMOVED***
				prop.set(receiver, val)
				return true, true
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			if receiver != proto ***REMOVED***
				return proto.self.setForeignIdx(idx, val, receiver, throw)
			***REMOVED***
			return proto.self.setOwnIdx(idx, val, throw), true
		***REMOVED***
	***REMOVED***
	return false, false
***REMOVED***

func (o *baseObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignStr(name, o.values[name], val, receiver, throw)
***REMOVED***

func (o *baseObject) setForeignIdx(name valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o.val.self.setForeignStr(name.string(), val, receiver, throw)
***REMOVED***

func (o *baseObject) setForeignSym(name *valueSymbol, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	var prop Value
	if o.symValues != nil ***REMOVED***
		prop = o.symValues.get(name)
	***REMOVED***
	if prop != nil ***REMOVED***
		if prop, ok := prop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				o.val.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
				return false, true
			***REMOVED***
			if prop.setterFunc != nil ***REMOVED***
				prop.set(receiver, val)
				return true, true
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			if receiver != o.val ***REMOVED***
				return proto.self.setForeignSym(name, val, receiver, throw)
			***REMOVED***
			return proto.self.setOwnSym(name, val, throw), true
		***REMOVED***
	***REMOVED***
	return false, false
***REMOVED***

func (o *baseObject) hasOwnPropertySym(s *valueSymbol) bool ***REMOVED***
	if o.symValues != nil ***REMOVED***
		return o.symValues.has(s)
	***REMOVED***
	return false
***REMOVED***

func (o *baseObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	_, exists := o.values[name]
	return exists
***REMOVED***

func (o *baseObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	return o.val.self.hasOwnPropertyStr(idx.string())
***REMOVED***

func (o *baseObject) _defineOwnProperty(name unistring.String, existingValue Value, descr PropertyDescriptor, throw bool) (val Value, ok bool) ***REMOVED***

	getterObj, _ := descr.Getter.(*Object)
	setterObj, _ := descr.Setter.(*Object)

	var existing *valueProperty

	if existingValue == nil ***REMOVED***
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot define property %s, object is not extensible", name)
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
	o.val.runtime.typeErrorResult(throw, "Cannot redefine property: %s", name)
	return nil, false

***REMOVED***

func (o *baseObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	existingVal := o.values[name]
	if v, ok := o._defineOwnProperty(name, existingVal, descr, throw); ok ***REMOVED***
		o.values[name] = v
		if existingVal == nil ***REMOVED***
			o.propNames = append(o.propNames, name)
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *baseObject) defineOwnPropertyIdx(idx valueInt, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	return o.val.self.defineOwnPropertyStr(idx.string(), desc, throw)
***REMOVED***

func (o *baseObject) defineOwnPropertySym(s *valueSymbol, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	var existingVal Value
	if o.symValues != nil ***REMOVED***
		existingVal = o.symValues.get(s)
	***REMOVED***
	if v, ok := o._defineOwnProperty(s.desc.string(), existingVal, descr, throw); ok ***REMOVED***
		if o.symValues == nil ***REMOVED***
			o.symValues = newOrderedMap(nil)
		***REMOVED***
		o.symValues.set(s, v)
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *baseObject) _put(name unistring.String, v Value) ***REMOVED***
	if _, exists := o.values[name]; !exists ***REMOVED***
		o.propNames = append(o.propNames, name)
	***REMOVED***

	o.values[name] = v
***REMOVED***

func valueProp(value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	if writable && enumerable && configurable ***REMOVED***
		return value
	***REMOVED***
	return &valueProperty***REMOVED***
		value:        value,
		writable:     writable,
		enumerable:   enumerable,
		configurable: configurable,
	***REMOVED***
***REMOVED***

func (o *baseObject) _putProp(name unistring.String, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	prop := valueProp(value, writable, enumerable, configurable)
	o._put(name, prop)
	return prop
***REMOVED***

func (o *baseObject) _putSym(s *valueSymbol, prop Value) ***REMOVED***
	if o.symValues == nil ***REMOVED***
		o.symValues = newOrderedMap(nil)
	***REMOVED***
	o.symValues.set(s, prop)
***REMOVED***

func (o *baseObject) tryPrimitive(methodName unistring.String) Value ***REMOVED***
	if method, ok := o.val.self.getStr(methodName, nil).(*Object); ok ***REMOVED***
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
	if v := o.tryPrimitive("valueOf"); v != nil ***REMOVED***
		return v
	***REMOVED***

	if v := o.tryPrimitive("toString"); v != nil ***REMOVED***
		return v
	***REMOVED***

	o.val.runtime.typeErrorResult(true, "Could not convert %v to primitive", o)
	return nil
***REMOVED***

func (o *Object) tryExoticToPrimitive(hint Value) Value ***REMOVED***
	exoticToPrimitive := toMethod(o.self.getSym(symToPrimitive, nil))
	if exoticToPrimitive != nil ***REMOVED***
		ret := exoticToPrimitive(FunctionCall***REMOVED***
			This:      o,
			Arguments: []Value***REMOVED***hint***REMOVED***,
		***REMOVED***)
		if _, fail := ret.(*Object); !fail ***REMOVED***
			return ret
		***REMOVED***
		panic(o.runtime.NewTypeError("Cannot convert object to primitive value"))
	***REMOVED***
	return nil
***REMOVED***

func (o *Object) toPrimitiveNumber() Value ***REMOVED***
	if v := o.tryExoticToPrimitive(hintNumber); v != nil ***REMOVED***
		return v
	***REMOVED***

	return o.self.toPrimitiveNumber()
***REMOVED***

func (o *Object) toPrimitiveString() Value ***REMOVED***
	if v := o.tryExoticToPrimitive(hintString); v != nil ***REMOVED***
		return v
	***REMOVED***

	return o.self.toPrimitiveString()
***REMOVED***

func (o *Object) toPrimitive() Value ***REMOVED***
	if v := o.tryExoticToPrimitive(hintDefault); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.self.toPrimitive()
***REMOVED***

func (o *baseObject) assertCallable() (func(FunctionCall) Value, bool) ***REMOVED***
	return nil, false
***REMOVED***

func (o *baseObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return nil
***REMOVED***

func (o *baseObject) proto() *Object ***REMOVED***
	return o.prototype
***REMOVED***

func (o *baseObject) isExtensible() bool ***REMOVED***
	return o.extensible
***REMOVED***

func (o *baseObject) preventExtensions(bool) bool ***REMOVED***
	o.extensible = false
	return true
***REMOVED***

func (o *baseObject) sortLen() int64 ***REMOVED***
	return toLength(o.val.self.getStr("length", nil))
***REMOVED***

func (o *baseObject) sortGet(i int64) Value ***REMOVED***
	return o.val.self.getIdx(valueInt(i), nil)
***REMOVED***

func (o *baseObject) swap(i, j int64) ***REMOVED***
	ii := valueInt(i)
	jj := valueInt(j)

	x := o.val.self.getIdx(ii, nil)
	y := o.val.self.getIdx(jj, nil)

	o.val.self.setOwnIdx(ii, y, false)
	o.val.self.setOwnIdx(jj, x, false)
***REMOVED***

func (o *baseObject) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, exists := ctx.get(o); exists ***REMOVED***
		return v
	***REMOVED***
	keys := o.ownKeys(false, nil)
	m := make(map[string]interface***REMOVED******REMOVED***, len(keys))
	ctx.put(o, m)
	for _, itemName := range keys ***REMOVED***
		itemNameStr := itemName.String()
		v := o.val.self.getStr(itemName.string(), nil)
		if v != nil ***REMOVED***
			m[itemNameStr] = exportValue(v, ctx)
		***REMOVED*** else ***REMOVED***
			m[itemNameStr] = nil
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
	name       unistring.String
	value      Value // set only when enumerable == _ENUM_UNKNOWN
	enumerable enumerableFlag
***REMOVED***

type objectPropIter struct ***REMOVED***
	o         *baseObject
	propNames []unistring.String
	idx       int
***REMOVED***

type propFilterIter struct ***REMOVED***
	wrapped iterNextFunc
	all     bool
	seen    map[unistring.String]bool
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

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *baseObject) enumerate() iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: o.val.self.enumerateUnfiltered(),
		seen:    make(map[unistring.String]bool),
	***REMOVED***).next
***REMOVED***

func (o *baseObject) ownIter() iterNextFunc ***REMOVED***
	if len(o.propNames) > o.lastSortedPropLen ***REMOVED***
		o.fixPropOrder()
	***REMOVED***
	propNames := make([]unistring.String, len(o.propNames))
	copy(propNames, o.propNames)
	return (&objectPropIter***REMOVED***
		o:         o,
		propNames: propNames,
	***REMOVED***).next
***REMOVED***

func (o *baseObject) recursiveIter(iter iterNextFunc) iterNextFunc ***REMOVED***
	return (&recursiveIter***REMOVED***
		o:       o,
		wrapped: iter,
	***REMOVED***).next
***REMOVED***

func (o *baseObject) enumerateUnfiltered() iterNextFunc ***REMOVED***
	return o.recursiveIter(o.ownIter())
***REMOVED***

type recursiveIter struct ***REMOVED***
	o       *baseObject
	wrapped iterNextFunc
***REMOVED***

func (iter *recursiveIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	item, next := iter.wrapped()
	if next != nil ***REMOVED***
		iter.wrapped = next
		return item, iter.next
	***REMOVED***
	if proto := iter.o.prototype; proto != nil ***REMOVED***
		return proto.self.enumerateUnfiltered()()
	***REMOVED***
	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *baseObject) equal(objectImpl) bool ***REMOVED***
	// Rely on parent reference comparison
	return false
***REMOVED***

// Reorder property names so that any integer properties are shifted to the beginning of the list
// in ascending order. This is to conform to ES6 9.1.12.
// Personally I think this requirement is strange. I can sort of understand where they are coming from,
// this way arrays can be specified just as objects with a 'magic' length property. However, I think
// it's safe to assume most devs don't use Objects to store integer properties. Therefore, performing
// property type checks when adding (and potentially looking up) properties would be unreasonable.
// Instead, we keep insertion order and only change it when (if) the properties get enumerated.
func (o *baseObject) fixPropOrder() ***REMOVED***
	names := o.propNames
	for i := o.lastSortedPropLen; i < len(names); i++ ***REMOVED***
		name := names[i]
		if idx := strToIdx(name); idx != math.MaxUint32 ***REMOVED***
			k := sort.Search(o.idxPropCount, func(j int) bool ***REMOVED***
				return strToIdx(names[j]) >= idx
			***REMOVED***)
			if k < i ***REMOVED***
				copy(names[k+1:i+1], names[k:i])
				names[k] = name
			***REMOVED***
			o.idxPropCount++
		***REMOVED***
	***REMOVED***
	o.lastSortedPropLen = len(names)
***REMOVED***

func (o *baseObject) ownKeys(all bool, keys []Value) []Value ***REMOVED***
	if len(o.propNames) > o.lastSortedPropLen ***REMOVED***
		o.fixPropOrder()
	***REMOVED***
	if all ***REMOVED***
		for _, k := range o.propNames ***REMOVED***
			keys = append(keys, stringValueFromRaw(k))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, k := range o.propNames ***REMOVED***
			prop := o.values[k]
			if prop, ok := prop.(*valueProperty); ok && !prop.enumerable ***REMOVED***
				continue
			***REMOVED***
			keys = append(keys, stringValueFromRaw(k))
		***REMOVED***
	***REMOVED***
	return keys
***REMOVED***

func (o *baseObject) ownSymbols(all bool, accum []Value) []Value ***REMOVED***
	if o.symValues != nil ***REMOVED***
		iter := o.symValues.newIter()
		if all ***REMOVED***
			for ***REMOVED***
				entry := iter.next()
				if entry == nil ***REMOVED***
					break
				***REMOVED***
				accum = append(accum, entry.key)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for ***REMOVED***
				entry := iter.next()
				if entry == nil ***REMOVED***
					break
				***REMOVED***
				if prop, ok := entry.value.(*valueProperty); ok ***REMOVED***
					if !prop.enumerable ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				accum = append(accum, entry.key)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return accum
***REMOVED***

func (o *baseObject) ownPropertyKeys(all bool, accum []Value) []Value ***REMOVED***
	return o.ownSymbols(all, o.val.self.ownKeys(all, accum))
***REMOVED***

func (o *baseObject) hasInstance(Value) bool ***REMOVED***
	panic(o.val.runtime.NewTypeError("Expecting a function in instanceof check, but got %s", o.val.toString()))
***REMOVED***

func toMethod(v Value) func(FunctionCall) Value ***REMOVED***
	if v == nil || IsUndefined(v) || IsNull(v) ***REMOVED***
		return nil
	***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		if call, ok := obj.self.assertCallable(); ok ***REMOVED***
			return call
		***REMOVED***
	***REMOVED***
	panic(typeError(fmt.Sprintf("%s is not a method", v.String())))
***REMOVED***

func instanceOfOperator(o Value, c *Object) bool ***REMOVED***
	if instOfHandler := toMethod(c.self.getSym(symHasInstance, c)); instOfHandler != nil ***REMOVED***
		return instOfHandler(FunctionCall***REMOVED***
			This:      c,
			Arguments: []Value***REMOVED***o***REMOVED***,
		***REMOVED***).ToBoolean()
	***REMOVED***

	return c.self.hasInstance(o)
***REMOVED***

func (o *Object) get(p Value, receiver Value) Value ***REMOVED***
	switch p := p.(type) ***REMOVED***
	case valueInt:
		return o.self.getIdx(p, receiver)
	case *valueSymbol:
		return o.self.getSym(p, receiver)
	default:
		return o.self.getStr(p.string(), receiver)
	***REMOVED***
***REMOVED***

func (o *Object) getOwnProp(p Value) Value ***REMOVED***
	switch p := p.(type) ***REMOVED***
	case valueInt:
		return o.self.getOwnPropIdx(p)
	case *valueSymbol:
		return o.self.getOwnPropSym(p)
	default:
		return o.self.getOwnPropStr(p.string())
	***REMOVED***
***REMOVED***

func (o *Object) hasOwnProperty(p Value) bool ***REMOVED***
	switch p := p.(type) ***REMOVED***
	case valueInt:
		return o.self.hasOwnPropertyIdx(p)
	case *valueSymbol:
		return o.self.hasOwnPropertySym(p)
	default:
		return o.self.hasOwnPropertyStr(p.string())
	***REMOVED***
***REMOVED***

func (o *Object) hasProperty(p Value) bool ***REMOVED***
	switch p := p.(type) ***REMOVED***
	case valueInt:
		return o.self.hasPropertyIdx(p)
	case *valueSymbol:
		return o.self.hasPropertySym(p)
	default:
		return o.self.hasPropertyStr(p.string())
	***REMOVED***
***REMOVED***

func (o *Object) setStr(name unistring.String, val, receiver Value, throw bool) bool ***REMOVED***
	if receiver == o ***REMOVED***
		return o.self.setOwnStr(name, val, throw)
	***REMOVED*** else ***REMOVED***
		if res, ok := o.self.setForeignStr(name, val, receiver, throw); !ok ***REMOVED***
			if robj, ok := receiver.(*Object); ok ***REMOVED***
				if prop := robj.self.getOwnPropStr(name); prop != nil ***REMOVED***
					if desc, ok := prop.(*valueProperty); ok ***REMOVED***
						if desc.accessor ***REMOVED***
							o.runtime.typeErrorResult(throw, "Receiver property %s is an accessor", name)
							return false
						***REMOVED***
						if !desc.writable ***REMOVED***
							o.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
							return false
						***REMOVED***
					***REMOVED***
					robj.self.defineOwnPropertyStr(name, PropertyDescriptor***REMOVED***Value: val***REMOVED***, throw)
				***REMOVED*** else ***REMOVED***
					robj.self.defineOwnPropertyStr(name, PropertyDescriptor***REMOVED***
						Value:        val,
						Writable:     FLAG_TRUE,
						Configurable: FLAG_TRUE,
						Enumerable:   FLAG_TRUE,
					***REMOVED***, throw)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				o.runtime.typeErrorResult(throw, "Receiver is not an object: %v", receiver)
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (o *Object) set(name Value, val, receiver Value, throw bool) bool ***REMOVED***
	switch name := name.(type) ***REMOVED***
	case valueInt:
		return o.setIdx(name, val, receiver, throw)
	case *valueSymbol:
		return o.setSym(name, val, receiver, throw)
	default:
		return o.setStr(name.string(), val, receiver, throw)
	***REMOVED***
***REMOVED***

func (o *Object) setOwn(name Value, val Value, throw bool) bool ***REMOVED***
	switch name := name.(type) ***REMOVED***
	case valueInt:
		return o.self.setOwnIdx(name, val, throw)
	case *valueSymbol:
		return o.self.setOwnSym(name, val, throw)
	default:
		return o.self.setOwnStr(name.string(), val, throw)
	***REMOVED***
***REMOVED***

func (o *Object) setIdx(name valueInt, val, receiver Value, throw bool) bool ***REMOVED***
	if receiver == o ***REMOVED***
		return o.self.setOwnIdx(name, val, throw)
	***REMOVED*** else ***REMOVED***
		if res, ok := o.self.setForeignIdx(name, val, receiver, throw); !ok ***REMOVED***
			if robj, ok := receiver.(*Object); ok ***REMOVED***
				if prop := robj.self.getOwnPropIdx(name); prop != nil ***REMOVED***
					if desc, ok := prop.(*valueProperty); ok ***REMOVED***
						if desc.accessor ***REMOVED***
							o.runtime.typeErrorResult(throw, "Receiver property %s is an accessor", name)
							return false
						***REMOVED***
						if !desc.writable ***REMOVED***
							o.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
							return false
						***REMOVED***
					***REMOVED***
					robj.self.defineOwnPropertyIdx(name, PropertyDescriptor***REMOVED***Value: val***REMOVED***, throw)
				***REMOVED*** else ***REMOVED***
					robj.self.defineOwnPropertyIdx(name, PropertyDescriptor***REMOVED***
						Value:        val,
						Writable:     FLAG_TRUE,
						Configurable: FLAG_TRUE,
						Enumerable:   FLAG_TRUE,
					***REMOVED***, throw)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				o.runtime.typeErrorResult(throw, "Receiver is not an object: %v", receiver)
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (o *Object) setSym(name *valueSymbol, val, receiver Value, throw bool) bool ***REMOVED***
	if receiver == o ***REMOVED***
		return o.self.setOwnSym(name, val, throw)
	***REMOVED*** else ***REMOVED***
		if res, ok := o.self.setForeignSym(name, val, receiver, throw); !ok ***REMOVED***
			if robj, ok := receiver.(*Object); ok ***REMOVED***
				if prop := robj.self.getOwnPropSym(name); prop != nil ***REMOVED***
					if desc, ok := prop.(*valueProperty); ok ***REMOVED***
						if desc.accessor ***REMOVED***
							o.runtime.typeErrorResult(throw, "Receiver property %s is an accessor", name)
							return false
						***REMOVED***
						if !desc.writable ***REMOVED***
							o.runtime.typeErrorResult(throw, "Cannot assign to read only property '%s'", name)
							return false
						***REMOVED***
					***REMOVED***
					robj.self.defineOwnPropertySym(name, PropertyDescriptor***REMOVED***Value: val***REMOVED***, throw)
				***REMOVED*** else ***REMOVED***
					robj.self.defineOwnPropertySym(name, PropertyDescriptor***REMOVED***
						Value:        val,
						Writable:     FLAG_TRUE,
						Configurable: FLAG_TRUE,
						Enumerable:   FLAG_TRUE,
					***REMOVED***, throw)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				o.runtime.typeErrorResult(throw, "Receiver is not an object: %v", receiver)
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (o *Object) delete(n Value, throw bool) bool ***REMOVED***
	switch n := n.(type) ***REMOVED***
	case valueInt:
		return o.self.deleteIdx(n, throw)
	case *valueSymbol:
		return o.self.deleteSym(n, throw)
	default:
		return o.self.deleteStr(n.string(), throw)
	***REMOVED***
***REMOVED***

func (o *Object) defineOwnProperty(n Value, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	switch n := n.(type) ***REMOVED***
	case valueInt:
		return o.self.defineOwnPropertyIdx(n, desc, throw)
	case *valueSymbol:
		return o.self.defineOwnPropertySym(n, desc, throw)
	default:
		return o.self.defineOwnPropertyStr(n.string(), desc, throw)
	***REMOVED***
***REMOVED***

func (o *Object) getWeakCollRefs() *weakCollections ***REMOVED***
	if o.weakColls == nil ***REMOVED***
		o.weakColls = &weakCollections***REMOVED***
			objId: o.getId(),
		***REMOVED***
		runtime.SetFinalizer(o.weakColls, finalizeObjectWeakRefs)
	***REMOVED***

	return o.weakColls
***REMOVED***

func (o *Object) getId() uint64 ***REMOVED***
	for o.id == 0 ***REMOVED***
		if o.runtime.hash == nil ***REMOVED***
			h := o.runtime.getHash()
			o.runtime.idSeq = h.Sum64()
		***REMOVED***
		o.id = o.runtime.idSeq
		o.runtime.idSeq++
	***REMOVED***
	return o.id
***REMOVED***

func (o *guardedObject) guard(props ...unistring.String) ***REMOVED***
	if o.guardedProps == nil ***REMOVED***
		o.guardedProps = make(map[unistring.String]struct***REMOVED******REMOVED***)
	***REMOVED***
	for _, p := range props ***REMOVED***
		o.guardedProps[p] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

func (o *guardedObject) check(p unistring.String) ***REMOVED***
	if _, exists := o.guardedProps[p]; exists ***REMOVED***
		o.val.self = &o.baseObject
	***REMOVED***
***REMOVED***

func (o *guardedObject) setOwnStr(p unistring.String, v Value, throw bool) bool ***REMOVED***
	res := o.baseObject.setOwnStr(p, v, throw)
	if res ***REMOVED***
		o.check(p)
	***REMOVED***
	return res
***REMOVED***

func (o *guardedObject) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	res := o.baseObject.defineOwnPropertyStr(name, desc, throw)
	if res ***REMOVED***
		o.check(name)
	***REMOVED***
	return res
***REMOVED***

func (o *guardedObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	res := o.baseObject.deleteStr(name, throw)
	if res ***REMOVED***
		o.check(name)
	***REMOVED***
	return res
***REMOVED***

func (ctx *objectExportCtx) get(key objectImpl) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	if v, exists := ctx.cache[key]; exists ***REMOVED***
		if item, ok := v.(objectExportCacheItem); ok ***REMOVED***
			r, exists := item[key.exportType()]
			return r, exists
		***REMOVED*** else ***REMOVED***
			return v, true
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func (ctx *objectExportCtx) getTyped(key objectImpl, typ reflect.Type) (interface***REMOVED******REMOVED***, bool) ***REMOVED***
	if v, exists := ctx.cache[key]; exists ***REMOVED***
		if item, ok := v.(objectExportCacheItem); ok ***REMOVED***
			r, exists := item[typ]
			return r, exists
		***REMOVED*** else ***REMOVED***
			if reflect.TypeOf(v) == typ ***REMOVED***
				return v, true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func (ctx *objectExportCtx) put(key objectImpl, value interface***REMOVED******REMOVED***) ***REMOVED***
	if ctx.cache == nil ***REMOVED***
		ctx.cache = make(map[objectImpl]interface***REMOVED******REMOVED***)
	***REMOVED***
	if item, ok := ctx.cache[key].(objectExportCacheItem); ok ***REMOVED***
		item[key.exportType()] = value
	***REMOVED*** else ***REMOVED***
		ctx.cache[key] = value
	***REMOVED***
***REMOVED***

func (ctx *objectExportCtx) putTyped(key objectImpl, typ reflect.Type, value interface***REMOVED******REMOVED***) ***REMOVED***
	if ctx.cache == nil ***REMOVED***
		ctx.cache = make(map[objectImpl]interface***REMOVED******REMOVED***)
	***REMOVED***
	v, exists := ctx.cache[key]
	if exists ***REMOVED***
		if item, ok := ctx.cache[key].(objectExportCacheItem); ok ***REMOVED***
			item[typ] = value
		***REMOVED*** else ***REMOVED***
			m := make(objectExportCacheItem, 2)
			m[key.exportType()] = v
			m[typ] = value
			ctx.cache[key] = m
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		m := make(objectExportCacheItem)
		m[typ] = value
		ctx.cache[key] = m
	***REMOVED***
***REMOVED***
