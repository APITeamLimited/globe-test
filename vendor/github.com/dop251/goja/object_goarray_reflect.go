package goja

import (
	"reflect"
	"strconv"

	"github.com/dop251/goja/unistring"
)

type objectGoArrayReflect struct ***REMOVED***
	objectGoReflect
	lengthProp valueProperty

	putIdx func(idx int, v Value, throw bool) bool
***REMOVED***

func (o *objectGoArrayReflect) _init() ***REMOVED***
	o.objectGoReflect.init()
	o.class = classArray
	o.prototype = o.val.runtime.global.ArrayPrototype
	o.updateLen()
	o.baseObject._put("length", &o.lengthProp)
***REMOVED***

func (o *objectGoArrayReflect) init() ***REMOVED***
	o._init()
	o.putIdx = o._putIdx
***REMOVED***

func (o *objectGoArrayReflect) updateLen() ***REMOVED***
	o.lengthProp.value = intToValue(int64(o.value.Len()))
***REMOVED***

func (o *objectGoArrayReflect) _hasIdx(idx valueInt) bool ***REMOVED***
	if idx := int64(idx); idx >= 0 && idx < int64(o.value.Len()) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoArrayReflect) _hasStr(name unistring.String) bool ***REMOVED***
	if idx := strToIdx64(name); idx >= 0 && idx < int64(o.value.Len()) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoArrayReflect) _getIdx(idx int) Value ***REMOVED***
	v := o.value.Index(idx)
	if (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && v.IsNil() ***REMOVED***
		return _null
	***REMOVED***
	return o.val.runtime.ToValue(v.Interface())
***REMOVED***

func (o *objectGoArrayReflect) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	if idx := toIntStrict(int64(idx)); idx >= 0 && idx < o.value.Len() ***REMOVED***
		return o._getIdx(idx)
	***REMOVED***
	return o.objectGoReflect.getStr(idx.string(), receiver)
***REMOVED***

func (o *objectGoArrayReflect) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	var ownProp Value
	if idx := strToGoIdx(name); idx >= 0 && idx < o.value.Len() ***REMOVED***
		ownProp = o._getIdx(idx)
	***REMOVED*** else if name == "length" ***REMOVED***
		ownProp = &o.lengthProp
	***REMOVED*** else ***REMOVED***
		ownProp = o.objectGoReflect.getOwnPropStr(name)
	***REMOVED***
	return o.getStrWithOwnProp(ownProp, name, receiver)
***REMOVED***

func (o *objectGoArrayReflect) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		if idx < o.value.Len() ***REMOVED***
			return &valueProperty***REMOVED***
				value:      o._getIdx(idx),
				writable:   true,
				enumerable: true,
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	if name == "length" ***REMOVED***
		return &o.lengthProp
	***REMOVED***
	return o.objectGoReflect.getOwnPropStr(name)
***REMOVED***

func (o *objectGoArrayReflect) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	if idx := toIntStrict(int64(idx)); idx >= 0 && idx < o.value.Len() ***REMOVED***
		return &valueProperty***REMOVED***
			value:      o._getIdx(idx),
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoArrayReflect) _putIdx(idx int, v Value, throw bool) bool ***REMOVED***
	err := o.val.runtime.toReflectValue(v, o.value.Index(idx), &objectExportCtx***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Go type conversion error: %v", err)
		return false
	***REMOVED***
	return true
***REMOVED***

func (o *objectGoArrayReflect) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	if i := toIntStrict(int64(idx)); i >= 0 ***REMOVED***
		if i >= o.value.Len() ***REMOVED***
			if res, ok := o._setForeignIdx(idx, nil, val, o.val, throw); ok ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		return o.putIdx(i, val, throw)
	***REMOVED*** else ***REMOVED***
		name := idx.string()
		if res, ok := o._setForeignStr(name, nil, val, o.val, throw); !ok ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't set property '%s' on Go slice", name)
			return false
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
***REMOVED***

func (o *objectGoArrayReflect) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		if idx >= o.value.Len() ***REMOVED***
			if res, ok := o._setForeignStr(name, nil, val, o.val, throw); ok ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		return o.putIdx(idx, val, throw)
	***REMOVED*** else ***REMOVED***
		if res, ok := o._setForeignStr(name, nil, val, o.val, throw); !ok ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't set property '%s' on Go slice", name)
			return false
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
***REMOVED***

func (o *objectGoArrayReflect) setForeignIdx(idx valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignIdx(idx, trueValIfPresent(o._hasIdx(idx)), val, receiver, throw)
***REMOVED***

func (o *objectGoArrayReflect) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignStr(name, trueValIfPresent(o.hasOwnPropertyStr(name)), val, receiver, throw)
***REMOVED***

func (o *objectGoArrayReflect) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	return o._hasIdx(idx)
***REMOVED***

func (o *objectGoArrayReflect) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if o._hasStr(name) || name == "length" ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect._has(name.String())
***REMOVED***

func (o *objectGoArrayReflect) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if i := toIntStrict(int64(idx)); i >= 0 ***REMOVED***
		if !o.val.runtime.checkHostObjectPropertyDescr(idx.string(), descr, throw) ***REMOVED***
			return false
		***REMOVED***
		val := descr.Value
		if val == nil ***REMOVED***
			val = _undefined
		***REMOVED***
		return o.putIdx(i, val, throw)
	***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Cannot define property '%d' on a Go slice", idx)
	return false
***REMOVED***

func (o *objectGoArrayReflect) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		if !o.val.runtime.checkHostObjectPropertyDescr(name, descr, throw) ***REMOVED***
			return false
		***REMOVED***
		val := descr.Value
		if val == nil ***REMOVED***
			val = _undefined
		***REMOVED***
		return o.putIdx(idx, val, throw)
	***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Cannot define property '%s' on a Go slice", name)
	return false
***REMOVED***

func (o *objectGoArrayReflect) toPrimitiveNumber() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoArrayReflect) toPrimitiveString() Value ***REMOVED***
	return o.val.runtime.arrayproto_join(FunctionCall***REMOVED***
		This: o.val,
	***REMOVED***)
***REMOVED***

func (o *objectGoArrayReflect) toPrimitive() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoArrayReflect) _deleteIdx(idx int) ***REMOVED***
	if idx < o.value.Len() ***REMOVED***
		o.value.Index(idx).Set(reflect.Zero(o.value.Type().Elem()))
	***REMOVED***
***REMOVED***

func (o *objectGoArrayReflect) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		o._deleteIdx(idx)
		return true
	***REMOVED***

	return o.objectGoReflect.deleteStr(name, throw)
***REMOVED***

func (o *objectGoArrayReflect) deleteIdx(i valueInt, throw bool) bool ***REMOVED***
	idx := toIntStrict(int64(i))
	if idx >= 0 ***REMOVED***
		o._deleteIdx(idx)
	***REMOVED***
	return true
***REMOVED***

type goArrayReflectPropIter struct ***REMOVED***
	o          *objectGoArrayReflect
	idx, limit int
***REMOVED***

func (i *goArrayReflectPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.limit && i.idx < i.o.value.Len() ***REMOVED***
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: asciiString(name), enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	return i.o.objectGoReflect.iterateStringKeys()()
***REMOVED***

func (o *objectGoArrayReflect) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	for i := 0; i < o.value.Len(); i++ ***REMOVED***
		accum = append(accum, asciiString(strconv.Itoa(i)))
	***REMOVED***

	return o.objectGoReflect.stringKeys(all, accum)
***REMOVED***

func (o *objectGoArrayReflect) iterateStringKeys() iterNextFunc ***REMOVED***
	return (&goArrayReflectPropIter***REMOVED***
		o:     o,
		limit: o.value.Len(),
	***REMOVED***).next
***REMOVED***

func (o *objectGoArrayReflect) equal(other objectImpl) bool ***REMOVED***
	if other, ok := other.(*objectGoArrayReflect); ok ***REMOVED***
		return o.value.Interface() == other.value.Interface()
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoArrayReflect) sortLen() int64 ***REMOVED***
	return int64(o.value.Len())
***REMOVED***

func (o *objectGoArrayReflect) sortGet(i int64) Value ***REMOVED***
	return o.getIdx(valueInt(i), nil)
***REMOVED***

func (o *objectGoArrayReflect) swap(i, j int64) ***REMOVED***
	ii := toIntStrict(i)
	jj := toIntStrict(j)
	x := o._getIdx(ii)
	y := o._getIdx(jj)

	o._putIdx(ii, y, false)
	o._putIdx(jj, x, false)
***REMOVED***
