package goja

import (
	"reflect"
	"strconv"

	"github.com/dop251/goja/unistring"
)

type objectGoSliceReflect struct ***REMOVED***
	objectGoReflect
	lengthProp      valueProperty
	sliceExtensible bool
***REMOVED***

func (o *objectGoSliceReflect) init() ***REMOVED***
	o.objectGoReflect.init()
	o.class = classArray
	o.prototype = o.val.runtime.global.ArrayPrototype
	o.sliceExtensible = o.value.CanSet()
	o.lengthProp.writable = o.sliceExtensible
	o.updateLen()
	o.baseObject._put("length", &o.lengthProp)
***REMOVED***

func (o *objectGoSliceReflect) updateLen() ***REMOVED***
	o.lengthProp.value = intToValue(int64(o.value.Len()))
***REMOVED***

func (o *objectGoSliceReflect) _hasIdx(idx valueInt) bool ***REMOVED***
	if idx := int64(idx); idx >= 0 && idx < int64(o.value.Len()) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSliceReflect) _hasStr(name unistring.String) bool ***REMOVED***
	if idx := strToIdx64(name); idx >= 0 && idx < int64(o.value.Len()) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSliceReflect) _getIdx(idx int) Value ***REMOVED***
	v := o.value.Index(idx)
	if (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && v.IsNil() ***REMOVED***
		return _null
	***REMOVED***
	return o.val.runtime.ToValue(v.Interface())
***REMOVED***

func (o *objectGoSliceReflect) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	if idx := toInt(int64(idx)); idx >= 0 && idx < o.value.Len() ***REMOVED***
		return o._getIdx(idx)
	***REMOVED***
	return o.objectGoReflect.getStr(idx.string(), receiver)
***REMOVED***

func (o *objectGoSliceReflect) getStr(name unistring.String, receiver Value) Value ***REMOVED***
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

func (o *objectGoSliceReflect) getOwnPropStr(name unistring.String) Value ***REMOVED***
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

func (o *objectGoSliceReflect) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	if idx := toInt(int64(idx)); idx >= 0 && idx < o.value.Len() ***REMOVED***
		return &valueProperty***REMOVED***
			value:      o._getIdx(idx),
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSliceReflect) putIdx(idx int, v Value, throw bool) bool ***REMOVED***
	if idx >= o.value.Len() ***REMOVED***
		if !o.sliceExtensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot extend a Go unaddressable reflect slice")
			return false
		***REMOVED***
		o.grow(idx + 1)
	***REMOVED***
	val, err := o.val.runtime.toReflectValue(v, o.value.Type().Elem())
	if err != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Go type conversion error: %v", err)
		return false
	***REMOVED***
	o.value.Index(idx).Set(val)
	return true
***REMOVED***

func (o *objectGoSliceReflect) grow(size int) ***REMOVED***
	newcap := o.value.Cap()
	if newcap < size ***REMOVED***
		// Use the same algorithm as in runtime.growSlice
		doublecap := newcap + newcap
		if size > doublecap ***REMOVED***
			newcap = size
		***REMOVED*** else ***REMOVED***
			if o.value.Len() < 1024 ***REMOVED***
				newcap = doublecap
			***REMOVED*** else ***REMOVED***
				for newcap < size ***REMOVED***
					newcap += newcap / 4
				***REMOVED***
			***REMOVED***
		***REMOVED***

		n := reflect.MakeSlice(o.value.Type(), size, newcap)
		reflect.Copy(n, o.value)
		o.value.Set(n)
	***REMOVED*** else ***REMOVED***
		tail := o.value.Slice(o.value.Len(), size)
		zero := reflect.Zero(o.value.Type().Elem())
		for i := 0; i < tail.Len(); i++ ***REMOVED***
			tail.Index(i).Set(zero)
		***REMOVED***
		o.value.SetLen(size)
	***REMOVED***
	o.updateLen()
***REMOVED***

func (o *objectGoSliceReflect) shrink(size int) ***REMOVED***
	tail := o.value.Slice(size, o.value.Len())
	zero := reflect.Zero(o.value.Type().Elem())
	for i := 0; i < tail.Len(); i++ ***REMOVED***
		tail.Index(i).Set(zero)
	***REMOVED***
	o.value.SetLen(size)
	o.updateLen()
***REMOVED***

func (o *objectGoSliceReflect) putLength(v Value, throw bool) bool ***REMOVED***
	newLen := toInt(toLength(v))
	curLen := o.value.Len()
	if newLen > curLen ***REMOVED***
		if !o.sliceExtensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot extend Go slice")
			return false
		***REMOVED***
		o.grow(newLen)
	***REMOVED*** else if newLen < curLen ***REMOVED***
		if !o.sliceExtensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot shrink Go slice")
			return false
		***REMOVED***
		o.shrink(newLen)
	***REMOVED***
	return true
***REMOVED***

func (o *objectGoSliceReflect) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	if i := toInt(int64(idx)); i >= 0 ***REMOVED***
		if i >= o.value.Len() ***REMOVED***
			if res, ok := o._setForeignIdx(idx, nil, val, o.val, throw); ok ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		o.putIdx(i, val, throw)
	***REMOVED*** else ***REMOVED***
		name := idx.string()
		if res, ok := o._setForeignStr(name, nil, val, o.val, throw); !ok ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't set property '%s' on Go slice", name)
			return false
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (o *objectGoSliceReflect) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		if idx >= o.value.Len() ***REMOVED***
			if res, ok := o._setForeignStr(name, nil, val, o.val, throw); ok ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		o.putIdx(idx, val, throw)
	***REMOVED*** else ***REMOVED***
		if name == "length" ***REMOVED***
			return o.putLength(val, throw)
		***REMOVED***
		if res, ok := o._setForeignStr(name, nil, val, o.val, throw); !ok ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't set property '%s' on Go slice", name)
			return false
		***REMOVED*** else ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (o *objectGoSliceReflect) setForeignIdx(idx valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignIdx(idx, trueValIfPresent(o._hasIdx(idx)), val, receiver, throw)
***REMOVED***

func (o *objectGoSliceReflect) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignStr(name, trueValIfPresent(o._hasStr(name)), val, receiver, throw)
***REMOVED***

func (o *objectGoSliceReflect) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	return o._hasIdx(idx)
***REMOVED***

func (o *objectGoSliceReflect) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if o._hasStr(name) ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect._has(name.String())
***REMOVED***

func (o *objectGoSliceReflect) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if i := toInt(int64(idx)); i >= 0 ***REMOVED***
		if !o.val.runtime.checkHostObjectPropertyDescr(idx.string(), descr, throw) ***REMOVED***
			return false
		***REMOVED***
		val := descr.Value
		if val == nil ***REMOVED***
			val = _undefined
		***REMOVED***
		o.putIdx(i, val, throw)
		return true
	***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Cannot define property '%d' on a Go slice", idx)
	return false
***REMOVED***

func (o *objectGoSliceReflect) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		if !o.val.runtime.checkHostObjectPropertyDescr(name, descr, throw) ***REMOVED***
			return false
		***REMOVED***
		val := descr.Value
		if val == nil ***REMOVED***
			val = _undefined
		***REMOVED***
		o.putIdx(idx, val, throw)
		return true
	***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Cannot define property '%s' on a Go slice", name)
	return false
***REMOVED***

func (o *objectGoSliceReflect) toPrimitiveNumber() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoSliceReflect) toPrimitiveString() Value ***REMOVED***
	return o.val.runtime.arrayproto_join(FunctionCall***REMOVED***
		This: o.val,
	***REMOVED***)
***REMOVED***

func (o *objectGoSliceReflect) toPrimitive() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoSliceReflect) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if idx := strToIdx64(name); idx >= 0 ***REMOVED***
		if idx < int64(o.value.Len()) ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't delete from Go slice")
			return false
		***REMOVED***
		return true
	***REMOVED***

	return o.objectGoReflect.deleteStr(name, throw)
***REMOVED***

func (o *objectGoSliceReflect) deleteIdx(i valueInt, throw bool) bool ***REMOVED***
	idx := int64(i)
	if idx >= 0 ***REMOVED***
		if idx < int64(o.value.Len()) ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't delete from Go slice")
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

type gosliceReflectPropIter struct ***REMOVED***
	o          *objectGoSliceReflect
	idx, limit int
***REMOVED***

func (i *gosliceReflectPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.limit && i.idx < i.o.value.Len() ***REMOVED***
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: unistring.String(name), enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	return i.o.objectGoReflect.enumerateUnfiltered()()
***REMOVED***

func (o *objectGoSliceReflect) ownKeys(all bool, accum []Value) []Value ***REMOVED***
	for i := 0; i < o.value.Len(); i++ ***REMOVED***
		accum = append(accum, asciiString(strconv.Itoa(i)))
	***REMOVED***

	return o.objectGoReflect.ownKeys(all, accum)
***REMOVED***

func (o *objectGoSliceReflect) enumerateUnfiltered() iterNextFunc ***REMOVED***
	return (&gosliceReflectPropIter***REMOVED***
		o:     o,
		limit: o.value.Len(),
	***REMOVED***).next
***REMOVED***

func (o *objectGoSliceReflect) equal(other objectImpl) bool ***REMOVED***
	if other, ok := other.(*objectGoSliceReflect); ok ***REMOVED***
		return o.value.Interface() == other.value.Interface()
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSliceReflect) sortLen() int64 ***REMOVED***
	return int64(o.value.Len())
***REMOVED***

func (o *objectGoSliceReflect) sortGet(i int64) Value ***REMOVED***
	return o.getIdx(valueInt(i), nil)
***REMOVED***

func (o *objectGoSliceReflect) swap(i, j int64) ***REMOVED***
	ii := valueInt(i)
	jj := valueInt(j)
	x := o.getIdx(ii, nil)
	y := o.getIdx(jj, nil)

	o.setOwnIdx(ii, y, false)
	o.setOwnIdx(jj, x, false)
***REMOVED***
