package goja

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/dop251/goja/unistring"
)

/*
DynamicObject is an interface representing a handler for a dynamic Object. Such an object can be created
using the Runtime.NewDynamicObject() method.

Note that Runtime.ToValue() does not have any special treatment for DynamicObject. The only way to create
a dynamic object is by using the Runtime.NewDynamicObject() method. This is done deliberately to avoid
silent code breaks when this interface changes.
*/
type DynamicObject interface ***REMOVED***
	// Get a property value for the key. May return nil if the property does not exist.
	Get(key string) Value
	// Set a property value for the key. Return true if success, false otherwise.
	Set(key string, val Value) bool
	// Has should return true if and only if the property exists.
	Has(key string) bool
	// Delete the property for the key. Returns true on success (note, that includes missing property).
	Delete(key string) bool
	// Keys returns a list of all existing property keys. There are no checks for duplicates or to make sure
	// that the order conforms to https://262.ecma-international.org/#sec-ordinaryownpropertykeys
	Keys() []string
***REMOVED***

/*
DynamicArray is an interface representing a handler for a dynamic array Object. Such an object can be created
using the Runtime.NewDynamicArray() method.

Any integer property key or a string property key that can be parsed into an int value (including negative
ones) is treated as an index and passed to the trap methods of the DynamicArray. Note this is different from
the regular ECMAScript arrays which only support positive indexes up to 2^32-1.

DynamicArray cannot be sparse, i.e. hasOwnProperty(num) will return true for num >= 0 && num < Len(). Deleting
such a property is equivalent to setting it to undefined. Note that this creates a slight peculiarity because
hasOwnProperty() will still return true, even after deletion.

Note that Runtime.ToValue() does not have any special treatment for DynamicArray. The only way to create
a dynamic array is by using the Runtime.NewDynamicArray() method. This is done deliberately to avoid
silent code breaks when this interface changes.
*/
type DynamicArray interface ***REMOVED***
	// Len returns the current array length.
	Len() int
	// Get an item at index idx. Note that idx may be any integer, negative or beyond the current length.
	Get(idx int) Value
	// Set an item at index idx. Note that idx may be any integer, negative or beyond the current length.
	// The expected behaviour when it's beyond length is that the array's length is increased to accommodate
	// the item. All elements in the 'new' section of the array should be zeroed.
	Set(idx int, val Value) bool
	// SetLen is called when the array's 'length' property is changed. If the length is increased all elements in the
	// 'new' section of the array should be zeroed.
	SetLen(int) bool
***REMOVED***

type baseDynamicObject struct ***REMOVED***
	val       *Object
	prototype *Object
***REMOVED***

type dynamicObject struct ***REMOVED***
	baseDynamicObject
	d DynamicObject
***REMOVED***

type dynamicArray struct ***REMOVED***
	baseDynamicObject
	a DynamicArray
***REMOVED***

/*
NewDynamicObject creates an Object backed by the provided DynamicObject handler.

All properties of this Object are Writable, Enumerable and Configurable data properties. Any attempt to define
a property that does not conform to this will fail.

The Object is always extensible and cannot be made non-extensible. Object.preventExtensions() will fail.

The Object's prototype is initially set to Object.prototype, but can be changed using regular mechanisms
(Object.SetPrototype() in Go or Object.setPrototypeOf() in JS).

The Object cannot have own Symbol properties, however its prototype can. If you need an iterator support for
example, you could create a regular object, set Symbol.iterator on that object and then use it as a
prototype. See TestDynamicObjectCustomProto for more details.

Export() returns the original DynamicObject.

This mechanism is similar to ECMAScript Proxy, however because all properties are enumerable and the object
is always extensible there is no need for invariant checks which removes the need to have a target object and
makes it a lot more efficient.
*/
func (r *Runtime) NewDynamicObject(d DynamicObject) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***
	o := &dynamicObject***REMOVED***
		d: d,
		baseDynamicObject: baseDynamicObject***REMOVED***
			val:       v,
			prototype: r.global.ObjectPrototype,
		***REMOVED***,
	***REMOVED***
	v.self = o
	return v
***REMOVED***

/*
NewDynamicArray creates an array Object backed by the provided DynamicArray handler.
It is similar to NewDynamicObject, the differences are:

- the Object is an array (i.e. Array.isArray() will return true and it will have the length property).

- the prototype will be initially set to Array.prototype.

- the Object cannot have any own string properties except for the 'length'.
*/
func (r *Runtime) NewDynamicArray(a DynamicArray) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***
	o := &dynamicArray***REMOVED***
		a: a,
		baseDynamicObject: baseDynamicObject***REMOVED***
			val:       v,
			prototype: r.global.ArrayPrototype,
		***REMOVED***,
	***REMOVED***
	v.self = o
	return v
***REMOVED***

func (*dynamicObject) sortLen() int64 ***REMOVED***
	return 0
***REMOVED***

func (*dynamicObject) sortGet(i int64) Value ***REMOVED***
	return nil
***REMOVED***

func (*dynamicObject) swap(i int64, i2 int64) ***REMOVED***
***REMOVED***

func (*dynamicObject) className() string ***REMOVED***
	return classObject
***REMOVED***

func (o *baseDynamicObject) getParentStr(p unistring.String, receiver Value) Value ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		if receiver == nil ***REMOVED***
			return proto.self.getStr(p, o.val)
		***REMOVED***
		return proto.self.getStr(p, receiver)
	***REMOVED***
	return nil
***REMOVED***

func (o *dynamicObject) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	prop := o.d.Get(p.String())
	if prop == nil ***REMOVED***
		return o.getParentStr(p, receiver)
	***REMOVED***
	return prop
***REMOVED***

func (o *baseDynamicObject) getParentIdx(p valueInt, receiver Value) Value ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		if receiver == nil ***REMOVED***
			return proto.self.getIdx(p, o.val)
		***REMOVED***
		return proto.self.getIdx(p, receiver)
	***REMOVED***
	return nil
***REMOVED***

func (o *dynamicObject) getIdx(p valueInt, receiver Value) Value ***REMOVED***
	prop := o.d.Get(p.String())
	if prop == nil ***REMOVED***
		return o.getParentIdx(p, receiver)
	***REMOVED***
	return prop
***REMOVED***

func (o *baseDynamicObject) getSym(p *Symbol, receiver Value) Value ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		if receiver == nil ***REMOVED***
			return proto.self.getSym(p, o.val)
		***REMOVED***
		return proto.self.getSym(p, receiver)
	***REMOVED***
	return nil
***REMOVED***

func (o *dynamicObject) getOwnPropStr(u unistring.String) Value ***REMOVED***
	return o.d.Get(u.String())
***REMOVED***

func (o *dynamicObject) getOwnPropIdx(v valueInt) Value ***REMOVED***
	return o.d.Get(v.String())
***REMOVED***

func (*baseDynamicObject) getOwnPropSym(*Symbol) Value ***REMOVED***
	return nil
***REMOVED***

func (o *dynamicObject) _set(prop string, v Value, throw bool) bool ***REMOVED***
	if o.d.Set(prop, v) ***REMOVED***
		return true
	***REMOVED***
	o.val.runtime.typeErrorResult(throw, "'Set' on a dynamic object returned false")
	return false
***REMOVED***

func (o *baseDynamicObject) _setSym(throw bool) ***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Dynamic objects do not support Symbol properties")
***REMOVED***

func (o *dynamicObject) setOwnStr(p unistring.String, v Value, throw bool) bool ***REMOVED***
	prop := p.String()
	if !o.d.Has(prop) ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			// we know it's foreign because prototype loops are not allowed
			if res, handled := proto.self.setForeignStr(p, v, o.val, throw); handled ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return o._set(prop, v, throw)
***REMOVED***

func (o *dynamicObject) setOwnIdx(p valueInt, v Value, throw bool) bool ***REMOVED***
	prop := p.String()
	if !o.d.Has(prop) ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			// we know it's foreign because prototype loops are not allowed
			if res, handled := proto.self.setForeignIdx(p, v, o.val, throw); handled ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return o._set(prop, v, throw)
***REMOVED***

func (o *baseDynamicObject) setOwnSym(s *Symbol, v Value, throw bool) bool ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		// we know it's foreign because prototype loops are not allowed
		if res, handled := proto.self.setForeignSym(s, v, o.val, throw); handled ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	o._setSym(throw)
	return false
***REMOVED***

func (o *baseDynamicObject) setParentForeignStr(p unistring.String, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		if receiver != proto ***REMOVED***
			return proto.self.setForeignStr(p, v, receiver, throw)
		***REMOVED***
		return proto.self.setOwnStr(p, v, throw), true
	***REMOVED***
	return false, false
***REMOVED***

func (o *dynamicObject) setForeignStr(p unistring.String, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	prop := p.String()
	if !o.d.Has(prop) ***REMOVED***
		return o.setParentForeignStr(p, v, receiver, throw)
	***REMOVED***
	return false, false
***REMOVED***

func (o *baseDynamicObject) setParentForeignIdx(p valueInt, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		if receiver != proto ***REMOVED***
			return proto.self.setForeignIdx(p, v, receiver, throw)
		***REMOVED***
		return proto.self.setOwnIdx(p, v, throw), true
	***REMOVED***
	return false, false
***REMOVED***

func (o *dynamicObject) setForeignIdx(p valueInt, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	prop := p.String()
	if !o.d.Has(prop) ***REMOVED***
		return o.setParentForeignIdx(p, v, receiver, throw)
	***REMOVED***
	return false, false
***REMOVED***

func (o *baseDynamicObject) setForeignSym(p *Symbol, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		if receiver != proto ***REMOVED***
			return proto.self.setForeignSym(p, v, receiver, throw)
		***REMOVED***
		return proto.self.setOwnSym(p, v, throw), true
	***REMOVED***
	return false, false
***REMOVED***

func (o *dynamicObject) hasPropertyStr(u unistring.String) bool ***REMOVED***
	if o.hasOwnPropertyStr(u) ***REMOVED***
		return true
	***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		return proto.self.hasPropertyStr(u)
	***REMOVED***
	return false
***REMOVED***

func (o *dynamicObject) hasPropertyIdx(idx valueInt) bool ***REMOVED***
	if o.hasOwnPropertyIdx(idx) ***REMOVED***
		return true
	***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		return proto.self.hasPropertyIdx(idx)
	***REMOVED***
	return false
***REMOVED***

func (o *baseDynamicObject) hasPropertySym(s *Symbol) bool ***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		return proto.self.hasPropertySym(s)
	***REMOVED***
	return false
***REMOVED***

func (o *dynamicObject) hasOwnPropertyStr(u unistring.String) bool ***REMOVED***
	return o.d.Has(u.String())
***REMOVED***

func (o *dynamicObject) hasOwnPropertyIdx(v valueInt) bool ***REMOVED***
	return o.d.Has(v.String())
***REMOVED***

func (*baseDynamicObject) hasOwnPropertySym(_ *Symbol) bool ***REMOVED***
	return false
***REMOVED***

func (o *baseDynamicObject) checkDynamicObjectPropertyDescr(name fmt.Stringer, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if descr.Getter != nil || descr.Setter != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Dynamic objects do not support accessor properties")
		return false
	***REMOVED***
	if descr.Writable == FLAG_FALSE ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Dynamic object field %q cannot be made read-only", name.String())
		return false
	***REMOVED***
	if descr.Enumerable == FLAG_FALSE ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Dynamic object field %q cannot be made non-enumerable", name.String())
		return false
	***REMOVED***
	if descr.Configurable == FLAG_FALSE ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Dynamic object field %q cannot be made non-configurable", name.String())
		return false
	***REMOVED***
	return true
***REMOVED***

func (o *dynamicObject) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	if o.checkDynamicObjectPropertyDescr(name, desc, throw) ***REMOVED***
		return o._set(name.String(), desc.Value, throw)
	***REMOVED***
	return false
***REMOVED***

func (o *dynamicObject) defineOwnPropertyIdx(name valueInt, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	if o.checkDynamicObjectPropertyDescr(name, desc, throw) ***REMOVED***
		return o._set(name.String(), desc.Value, throw)
	***REMOVED***
	return false
***REMOVED***

func (o *baseDynamicObject) defineOwnPropertySym(name *Symbol, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	o._setSym(throw)
	return false
***REMOVED***

func (o *dynamicObject) _delete(prop string, throw bool) bool ***REMOVED***
	if o.d.Delete(prop) ***REMOVED***
		return true
	***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Could not delete property %q of a dynamic object", prop)
	return false
***REMOVED***

func (o *dynamicObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	return o._delete(name.String(), throw)
***REMOVED***

func (o *dynamicObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	return o._delete(idx.String(), throw)
***REMOVED***

func (*baseDynamicObject) deleteSym(_ *Symbol, _ bool) bool ***REMOVED***
	return true
***REMOVED***

func (o *baseDynamicObject) toPrimitiveNumber() Value ***REMOVED***
	return o.val.genericToPrimitiveNumber()
***REMOVED***

func (o *baseDynamicObject) toPrimitiveString() Value ***REMOVED***
	return o.val.genericToPrimitiveString()
***REMOVED***

func (o *baseDynamicObject) toPrimitive() Value ***REMOVED***
	return o.val.genericToPrimitive()
***REMOVED***

func (o *baseDynamicObject) assertCallable() (call func(FunctionCall) Value, ok bool) ***REMOVED***
	return nil, false
***REMOVED***

func (*baseDynamicObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return nil
***REMOVED***

func (o *baseDynamicObject) proto() *Object ***REMOVED***
	return o.prototype
***REMOVED***

func (o *baseDynamicObject) setProto(proto *Object, throw bool) bool ***REMOVED***
	o.prototype = proto
	return true
***REMOVED***

func (o *baseDynamicObject) hasInstance(v Value) bool ***REMOVED***
	panic(o.val.runtime.NewTypeError("Expecting a function in instanceof check, but got a dynamic object"))
***REMOVED***

func (*baseDynamicObject) isExtensible() bool ***REMOVED***
	return true
***REMOVED***

func (o *baseDynamicObject) preventExtensions(throw bool) bool ***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Cannot make a dynamic object non-extensible")
	return false
***REMOVED***

type dynamicObjectPropIter struct ***REMOVED***
	o         *dynamicObject
	propNames []string
	idx       int
***REMOVED***

func (i *dynamicObjectPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.propNames) ***REMOVED***
		name := i.propNames[i.idx]
		i.idx++
		if i.o.d.Has(name) ***REMOVED***
			return propIterItem***REMOVED***name: unistring.NewFromString(name), enumerable: _ENUM_TRUE***REMOVED***, i.next
		***REMOVED***
	***REMOVED***
	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *dynamicObject) enumerateOwnKeys() iterNextFunc ***REMOVED***
	keys := o.d.Keys()
	return (&dynamicObjectPropIter***REMOVED***
		o:         o,
		propNames: keys,
	***REMOVED***).next
***REMOVED***

func (o *dynamicObject) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return o.d
***REMOVED***

func (o *dynamicObject) exportType() reflect.Type ***REMOVED***
	return reflect.TypeOf(o.d)
***REMOVED***

func (o *dynamicObject) equal(impl objectImpl) bool ***REMOVED***
	if other, ok := impl.(*dynamicObject); ok ***REMOVED***
		return o.d == other.d
	***REMOVED***
	return false
***REMOVED***

func (o *dynamicObject) ownKeys(all bool, accum []Value) []Value ***REMOVED***
	keys := o.d.Keys()
	if l := len(accum) + len(keys); l > cap(accum) ***REMOVED***
		oldAccum := accum
		accum = make([]Value, len(accum), l)
		copy(accum, oldAccum)
	***REMOVED***
	for _, key := range keys ***REMOVED***
		accum = append(accum, newStringValue(key))
	***REMOVED***
	return accum
***REMOVED***

func (*baseDynamicObject) ownSymbols(all bool, accum []Value) []Value ***REMOVED***
	return accum
***REMOVED***

func (o *dynamicObject) ownPropertyKeys(all bool, accum []Value) []Value ***REMOVED***
	return o.ownKeys(all, accum)
***REMOVED***

func (*baseDynamicObject) _putProp(name unistring.String, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	return nil
***REMOVED***

func (*baseDynamicObject) _putSym(s *Symbol, prop Value) ***REMOVED***
***REMOVED***

func (a *dynamicArray) sortLen() int64 ***REMOVED***
	return int64(a.a.Len())
***REMOVED***

func (a *dynamicArray) sortGet(i int64) Value ***REMOVED***
	return a.a.Get(int(i))
***REMOVED***

func (a *dynamicArray) swap(i int64, j int64) ***REMOVED***
	x := a.sortGet(i)
	y := a.sortGet(j)
	a.a.Set(int(i), y)
	a.a.Set(int(j), x)
***REMOVED***

func (a *dynamicArray) className() string ***REMOVED***
	return classArray
***REMOVED***

func (a *dynamicArray) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	if p == "length" ***REMOVED***
		return intToValue(int64(a.a.Len()))
	***REMOVED***
	if idx, ok := strPropToInt(p); ok ***REMOVED***
		return a.a.Get(idx)
	***REMOVED***
	return a.getParentStr(p, receiver)
***REMOVED***

func (a *dynamicArray) getIdx(p valueInt, receiver Value) Value ***REMOVED***
	if val := a.getOwnPropIdx(p); val != nil ***REMOVED***
		return val
	***REMOVED***
	return a.getParentIdx(p, receiver)
***REMOVED***

func (a *dynamicArray) getOwnPropStr(u unistring.String) Value ***REMOVED***
	if u == "length" ***REMOVED***
		return &valueProperty***REMOVED***
			value:    intToValue(int64(a.a.Len())),
			writable: true,
		***REMOVED***
	***REMOVED***
	if idx, ok := strPropToInt(u); ok ***REMOVED***
		return a.a.Get(idx)
	***REMOVED***
	return nil
***REMOVED***

func (a *dynamicArray) getOwnPropIdx(v valueInt) Value ***REMOVED***
	return a.a.Get(toIntStrict(int64(v)))
***REMOVED***

func (a *dynamicArray) _setLen(v Value, throw bool) bool ***REMOVED***
	if a.a.SetLen(toIntStrict(v.ToInteger())) ***REMOVED***
		return true
	***REMOVED***
	a.val.runtime.typeErrorResult(throw, "'SetLen' on a dynamic array returned false")
	return false
***REMOVED***

func (a *dynamicArray) setOwnStr(p unistring.String, v Value, throw bool) bool ***REMOVED***
	if p == "length" ***REMOVED***
		return a._setLen(v, throw)
	***REMOVED***
	if idx, ok := strPropToInt(p); ok ***REMOVED***
		return a._setIdx(idx, v, throw)
	***REMOVED***
	a.val.runtime.typeErrorResult(throw, "Cannot set property %q on a dynamic array", p.String())
	return false
***REMOVED***

func (a *dynamicArray) _setIdx(idx int, v Value, throw bool) bool ***REMOVED***
	if a.a.Set(idx, v) ***REMOVED***
		return true
	***REMOVED***
	a.val.runtime.typeErrorResult(throw, "'Set' on a dynamic array returned false")
	return false
***REMOVED***

func (a *dynamicArray) setOwnIdx(p valueInt, v Value, throw bool) bool ***REMOVED***
	return a._setIdx(toIntStrict(int64(p)), v, throw)
***REMOVED***

func (a *dynamicArray) setForeignStr(p unistring.String, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	return a.setParentForeignStr(p, v, receiver, throw)
***REMOVED***

func (a *dynamicArray) setForeignIdx(p valueInt, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	return a.setParentForeignIdx(p, v, receiver, throw)
***REMOVED***

func (a *dynamicArray) hasPropertyStr(u unistring.String) bool ***REMOVED***
	if a.hasOwnPropertyStr(u) ***REMOVED***
		return true
	***REMOVED***
	if proto := a.prototype; proto != nil ***REMOVED***
		return proto.self.hasPropertyStr(u)
	***REMOVED***
	return false
***REMOVED***

func (a *dynamicArray) hasPropertyIdx(idx valueInt) bool ***REMOVED***
	if a.hasOwnPropertyIdx(idx) ***REMOVED***
		return true
	***REMOVED***
	if proto := a.prototype; proto != nil ***REMOVED***
		return proto.self.hasPropertyIdx(idx)
	***REMOVED***
	return false
***REMOVED***

func (a *dynamicArray) _has(idx int) bool ***REMOVED***
	return idx >= 0 && idx < a.a.Len()
***REMOVED***

func (a *dynamicArray) hasOwnPropertyStr(u unistring.String) bool ***REMOVED***
	if u == "length" ***REMOVED***
		return true
	***REMOVED***
	if idx, ok := strPropToInt(u); ok ***REMOVED***
		return a._has(idx)
	***REMOVED***
	return false
***REMOVED***

func (a *dynamicArray) hasOwnPropertyIdx(v valueInt) bool ***REMOVED***
	return a._has(toIntStrict(int64(v)))
***REMOVED***

func (a *dynamicArray) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	if a.checkDynamicObjectPropertyDescr(name, desc, throw) ***REMOVED***
		if idx, ok := strPropToInt(name); ok ***REMOVED***
			return a._setIdx(idx, desc.Value, throw)
		***REMOVED***
		a.val.runtime.typeErrorResult(throw, "Cannot define property %q on a dynamic array", name.String())
	***REMOVED***
	return false
***REMOVED***

func (a *dynamicArray) defineOwnPropertyIdx(name valueInt, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	if a.checkDynamicObjectPropertyDescr(name, desc, throw) ***REMOVED***
		return a._setIdx(toIntStrict(int64(name)), desc.Value, throw)
	***REMOVED***
	return false
***REMOVED***

func (a *dynamicArray) _delete(idx int, throw bool) bool ***REMOVED***
	if a._has(idx) ***REMOVED***
		a._setIdx(idx, _undefined, throw)
	***REMOVED***
	return true
***REMOVED***

func (a *dynamicArray) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if idx, ok := strPropToInt(name); ok ***REMOVED***
		return a._delete(idx, throw)
	***REMOVED***
	if a.hasOwnPropertyStr(name) ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "Cannot delete property %q on a dynamic array", name.String())
		return false
	***REMOVED***
	return true
***REMOVED***

func (a *dynamicArray) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	return a._delete(toIntStrict(int64(idx)), throw)
***REMOVED***

type dynArrayPropIter struct ***REMOVED***
	a          DynamicArray
	idx, limit int
***REMOVED***

func (i *dynArrayPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.limit && i.idx < i.a.Len() ***REMOVED***
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: unistring.String(name), enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (a *dynamicArray) enumerateOwnKeys() iterNextFunc ***REMOVED***
	return (&dynArrayPropIter***REMOVED***
		a:     a.a,
		limit: a.a.Len(),
	***REMOVED***).next
***REMOVED***

func (a *dynamicArray) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return a.a
***REMOVED***

func (a *dynamicArray) exportType() reflect.Type ***REMOVED***
	return reflect.TypeOf(a.a)
***REMOVED***

func (a *dynamicArray) equal(impl objectImpl) bool ***REMOVED***
	if other, ok := impl.(*dynamicArray); ok ***REMOVED***
		return a == other
	***REMOVED***
	return false
***REMOVED***

func (a *dynamicArray) ownKeys(all bool, accum []Value) []Value ***REMOVED***
	al := a.a.Len()
	l := len(accum) + al
	if all ***REMOVED***
		l++
	***REMOVED***
	if l > cap(accum) ***REMOVED***
		oldAccum := accum
		accum = make([]Value, len(oldAccum), l)
		copy(accum, oldAccum)
	***REMOVED***
	for i := 0; i < al; i++ ***REMOVED***
		accum = append(accum, asciiString(strconv.Itoa(i)))
	***REMOVED***
	if all ***REMOVED***
		accum = append(accum, asciiString("length"))
	***REMOVED***
	return accum
***REMOVED***

func (a *dynamicArray) ownPropertyKeys(all bool, accum []Value) []Value ***REMOVED***
	return a.ownKeys(all, accum)
***REMOVED***
