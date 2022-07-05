package goja

import (
	"reflect"

	"github.com/dop251/goja/unistring"
)

type objectGoMapReflect struct ***REMOVED***
	objectGoReflect

	keyType, valueType reflect.Type
***REMOVED***

func (o *objectGoMapReflect) init() ***REMOVED***
	o.objectGoReflect.init()
	o.keyType = o.value.Type().Key()
	o.valueType = o.value.Type().Elem()
***REMOVED***

func (o *objectGoMapReflect) toKey(n Value, throw bool) reflect.Value ***REMOVED***
	key := reflect.New(o.keyType).Elem()
	err := o.val.runtime.toReflectValue(n, key, &objectExportCtx***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "map key conversion error: %v", err)
		return reflect.Value***REMOVED******REMOVED***
	***REMOVED***
	return key
***REMOVED***

func (o *objectGoMapReflect) strToKey(name string, throw bool) reflect.Value ***REMOVED***
	if o.keyType.Kind() == reflect.String ***REMOVED***
		return reflect.ValueOf(name).Convert(o.keyType)
	***REMOVED***
	return o.toKey(newStringValue(name), throw)
***REMOVED***

func (o *objectGoMapReflect) _get(n Value) Value ***REMOVED***
	key := o.toKey(n, false)
	if !key.IsValid() ***REMOVED***
		return nil
	***REMOVED***
	if v := o.value.MapIndex(key); v.IsValid() ***REMOVED***
		return o.val.runtime.ToValue(v.Interface())
	***REMOVED***

	return nil
***REMOVED***

func (o *objectGoMapReflect) _getStr(name string) Value ***REMOVED***
	key := o.strToKey(name, false)
	if !key.IsValid() ***REMOVED***
		return nil
	***REMOVED***
	if v := o.value.MapIndex(key); v.IsValid() ***REMOVED***
		return o.val.runtime.ToValue(v.Interface())
	***REMOVED***

	return nil
***REMOVED***

func (o *objectGoMapReflect) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	if v := o._getStr(name.String()); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.getStr(name, receiver)
***REMOVED***

func (o *objectGoMapReflect) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	if v := o._get(idx); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.getIdx(idx, receiver)
***REMOVED***

func (o *objectGoMapReflect) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if v := o._getStr(name.String()); v != nil ***REMOVED***
		return &valueProperty***REMOVED***
			value:      v,
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return o.objectGoReflect.getOwnPropStr(name)
***REMOVED***

func (o *objectGoMapReflect) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	if v := o._get(idx); v != nil ***REMOVED***
		return &valueProperty***REMOVED***
			value:      v,
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return o.objectGoReflect.getOwnPropStr(idx.string())
***REMOVED***

func (o *objectGoMapReflect) toValue(val Value, throw bool) (reflect.Value, bool) ***REMOVED***
	v := reflect.New(o.valueType).Elem()
	err := o.val.runtime.toReflectValue(val, v, &objectExportCtx***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "map value conversion error: %v", err)
		return reflect.Value***REMOVED******REMOVED***, false
	***REMOVED***

	return v, true
***REMOVED***

func (o *objectGoMapReflect) _put(key reflect.Value, val Value, throw bool) bool ***REMOVED***
	if key.IsValid() ***REMOVED***
		if o.extensible || o.value.MapIndex(key).IsValid() ***REMOVED***
			v, ok := o.toValue(val, throw)
			if !ok ***REMOVED***
				return false
			***REMOVED***
			o.value.SetMapIndex(key, v)
		***REMOVED*** else ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot set property %s, object is not extensible", key.String())
			return false
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoMapReflect) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	n := name.String()
	key := o.strToKey(n, false)
	if !key.IsValid() || !o.value.MapIndex(key).IsValid() ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			// we know it's foreign because prototype loops are not allowed
			if res, ok := proto.self.setForeignStr(name, val, o.val, throw); ok ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		// new property
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot add property %s, object is not extensible", n)
			return false
		***REMOVED*** else ***REMOVED***
			if throw && !key.IsValid() ***REMOVED***
				o.strToKey(n, true)
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	o._put(key, val, throw)
	return true
***REMOVED***

func (o *objectGoMapReflect) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	key := o.toKey(idx, false)
	if !key.IsValid() || !o.value.MapIndex(key).IsValid() ***REMOVED***
		if proto := o.prototype; proto != nil ***REMOVED***
			// we know it's foreign because prototype loops are not allowed
			if res, ok := proto.self.setForeignIdx(idx, val, o.val, throw); ok ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		// new property
		if !o.extensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot add property %d, object is not extensible", idx)
			return false
		***REMOVED*** else ***REMOVED***
			if throw && !key.IsValid() ***REMOVED***
				o.toKey(idx, true)
				return false
			***REMOVED***
		***REMOVED***
	***REMOVED***
	o._put(key, val, throw)
	return true
***REMOVED***

func (o *objectGoMapReflect) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignStr(name, trueValIfPresent(o.hasOwnPropertyStr(name)), val, receiver, throw)
***REMOVED***

func (o *objectGoMapReflect) setForeignIdx(idx valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignIdx(idx, trueValIfPresent(o.hasOwnPropertyIdx(idx)), val, receiver, throw)
***REMOVED***

func (o *objectGoMapReflect) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if !o.val.runtime.checkHostObjectPropertyDescr(name, descr, throw) ***REMOVED***
		return false
	***REMOVED***

	return o._put(o.strToKey(name.String(), throw), descr.Value, throw)
***REMOVED***

func (o *objectGoMapReflect) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if !o.val.runtime.checkHostObjectPropertyDescr(idx.string(), descr, throw) ***REMOVED***
		return false
	***REMOVED***

	return o._put(o.toKey(idx, throw), descr.Value, throw)
***REMOVED***

func (o *objectGoMapReflect) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	key := o.strToKey(name.String(), false)
	if key.IsValid() && o.value.MapIndex(key).IsValid() ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoMapReflect) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	key := o.toKey(idx, false)
	if key.IsValid() && o.value.MapIndex(key).IsValid() ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoMapReflect) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	key := o.strToKey(name.String(), throw)
	if !key.IsValid() ***REMOVED***
		return false
	***REMOVED***
	o.value.SetMapIndex(key, reflect.Value***REMOVED******REMOVED***)
	return true
***REMOVED***

func (o *objectGoMapReflect) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	key := o.toKey(idx, throw)
	if !key.IsValid() ***REMOVED***
		return false
	***REMOVED***
	o.value.SetMapIndex(key, reflect.Value***REMOVED******REMOVED***)
	return true
***REMOVED***

type gomapReflectPropIter struct ***REMOVED***
	o    *objectGoMapReflect
	keys []reflect.Value
	idx  int
***REMOVED***

func (i *gomapReflectPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.keys) ***REMOVED***
		key := i.keys[i.idx]
		v := i.o.value.MapIndex(key)
		i.idx++
		if v.IsValid() ***REMOVED***
			return propIterItem***REMOVED***name: newStringValue(key.String()), enumerable: _ENUM_TRUE***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoMapReflect) iterateStringKeys() iterNextFunc ***REMOVED***
	return (&gomapReflectPropIter***REMOVED***
		o:    o,
		keys: o.value.MapKeys(),
	***REMOVED***).next
***REMOVED***

func (o *objectGoMapReflect) stringKeys(_ bool, accum []Value) []Value ***REMOVED***
	// all own keys are enumerable
	for _, key := range o.value.MapKeys() ***REMOVED***
		accum = append(accum, newStringValue(key.String()))
	***REMOVED***

	return accum
***REMOVED***
