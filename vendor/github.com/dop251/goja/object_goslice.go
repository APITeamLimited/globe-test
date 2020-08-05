package goja

import (
	"reflect"
	"strconv"

	"github.com/dop251/goja/unistring"
)

type objectGoSlice struct ***REMOVED***
	baseObject
	data            *[]interface***REMOVED******REMOVED***
	lengthProp      valueProperty
	sliceExtensible bool
***REMOVED***

func (o *objectGoSlice) init() ***REMOVED***
	o.baseObject.init()
	o.class = classArray
	o.prototype = o.val.runtime.global.ArrayPrototype
	o.lengthProp.writable = o.sliceExtensible
	o.extensible = true
	o.updateLen()
	o.baseObject._put("length", &o.lengthProp)
***REMOVED***

func (o *objectGoSlice) updateLen() ***REMOVED***
	o.lengthProp.value = intToValue(int64(len(*o.data)))
***REMOVED***

func (o *objectGoSlice) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	var ownProp Value
	if idx := strToGoIdx(name); idx >= 0 && idx < len(*o.data) ***REMOVED***
		v := (*o.data)[idx]
		ownProp = o.val.runtime.ToValue(v)
	***REMOVED*** else if name == "length" ***REMOVED***
		ownProp = &o.lengthProp
	***REMOVED***

	return o.getStrWithOwnProp(ownProp, name, receiver)
***REMOVED***

func (o *objectGoSlice) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	if idx := int64(idx); idx >= 0 && idx < int64(len(*o.data)) ***REMOVED***
		v := (*o.data)[idx]
		return o.val.runtime.ToValue(v)
	***REMOVED***
	if o.prototype != nil ***REMOVED***
		if receiver == nil ***REMOVED***
			return o.prototype.self.getIdx(idx, o.val)
		***REMOVED***
		return o.prototype.self.getIdx(idx, receiver)
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSlice) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		if idx < len(*o.data) ***REMOVED***
			v := o.val.runtime.ToValue((*o.data)[idx])
			return &valueProperty***REMOVED***
				value:      v,
				writable:   true,
				enumerable: true,
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	if name == "length" ***REMOVED***
		return &o.lengthProp
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSlice) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	if idx := int64(idx); idx >= 0 && idx < int64(len(*o.data)) ***REMOVED***
		v := o.val.runtime.ToValue((*o.data)[idx])
		return &valueProperty***REMOVED***
			value:      v,
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSlice) grow(size int) ***REMOVED***
	newcap := cap(*o.data)
	if newcap < size ***REMOVED***
		// Use the same algorithm as in runtime.growSlice
		doublecap := newcap + newcap
		if size > doublecap ***REMOVED***
			newcap = size
		***REMOVED*** else ***REMOVED***
			if len(*o.data) < 1024 ***REMOVED***
				newcap = doublecap
			***REMOVED*** else ***REMOVED***
				for newcap < size ***REMOVED***
					newcap += newcap / 4
				***REMOVED***
			***REMOVED***
		***REMOVED***

		n := make([]interface***REMOVED******REMOVED***, size, newcap)
		copy(n, *o.data)
		*o.data = n
	***REMOVED*** else ***REMOVED***
		tail := (*o.data)[len(*o.data):size]
		for k := range tail ***REMOVED***
			tail[k] = nil
		***REMOVED***
		*o.data = (*o.data)[:size]
	***REMOVED***
	o.updateLen()
***REMOVED***

func (o *objectGoSlice) shrink(size int) ***REMOVED***
	tail := (*o.data)[size:]
	for k := range tail ***REMOVED***
		tail[k] = nil
	***REMOVED***
	*o.data = (*o.data)[:size]
	o.updateLen()
***REMOVED***

func (o *objectGoSlice) putIdx(idx int, v Value, throw bool) ***REMOVED***
	if idx >= len(*o.data) ***REMOVED***
		if !o.sliceExtensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot extend Go slice")
			return
		***REMOVED***
		o.grow(idx + 1)
	***REMOVED***
	(*o.data)[idx] = v.Export()
***REMOVED***

func (o *objectGoSlice) putLength(v Value, throw bool) bool ***REMOVED***
	newLen := toInt(toLength(v))
	curLen := len(*o.data)
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

func (o *objectGoSlice) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	if i := toInt(int64(idx)); i >= 0 ***REMOVED***
		if i >= len(*o.data) ***REMOVED***
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

func (o *objectGoSlice) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if idx := strToGoIdx(name); idx >= 0 ***REMOVED***
		if idx >= len(*o.data) ***REMOVED***
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

func (o *objectGoSlice) setForeignIdx(idx valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignIdx(idx, trueValIfPresent(o.hasOwnPropertyIdx(idx)), val, receiver, throw)
***REMOVED***

func (o *objectGoSlice) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignStr(name, trueValIfPresent(o.hasOwnPropertyStr(name)), val, receiver, throw)
***REMOVED***

func (o *objectGoSlice) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	if idx := int64(idx); idx >= 0 ***REMOVED***
		return idx < int64(len(*o.data))
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSlice) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if idx := strToIdx64(name); idx >= 0 ***REMOVED***
		return idx < int64(len(*o.data))
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSlice) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
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

func (o *objectGoSlice) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
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
	if name == "length" ***REMOVED***
		return o.val.runtime.defineArrayLength(&o.lengthProp, descr, o.putLength, throw)
	***REMOVED***
	o.val.runtime.typeErrorResult(throw, "Cannot define property '%s' on a Go slice", name)
	return false
***REMOVED***

func (o *objectGoSlice) toPrimitiveNumber() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoSlice) toPrimitiveString() Value ***REMOVED***
	return o.val.runtime.arrayproto_join(FunctionCall***REMOVED***
		This: o.val,
	***REMOVED***)
***REMOVED***

func (o *objectGoSlice) toPrimitive() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoSlice) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if idx := strToIdx64(name); idx >= 0 ***REMOVED***
		if idx < int64(len(*o.data)) ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't delete from Go slice")
			return false
		***REMOVED***
		return true
	***REMOVED***
	return o.baseObject.deleteStr(name, throw)
***REMOVED***

func (o *objectGoSlice) deleteIdx(i valueInt, throw bool) bool ***REMOVED***
	idx := int64(i)
	if idx >= 0 ***REMOVED***
		if idx < int64(len(*o.data)) ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Can't delete from Go slice")
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

type goslicePropIter struct ***REMOVED***
	o          *objectGoSlice
	idx, limit int
***REMOVED***

func (i *goslicePropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.limit && i.idx < len(*i.o.data) ***REMOVED***
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: unistring.String(name), enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoSlice) enumerateUnfiltered() iterNextFunc ***REMOVED***
	return o.recursiveIter((&goslicePropIter***REMOVED***
		o:     o,
		limit: len(*o.data),
	***REMOVED***).next)
***REMOVED***

func (o *objectGoSlice) ownKeys(_ bool, accum []Value) []Value ***REMOVED***
	for i := range *o.data ***REMOVED***
		accum = append(accum, asciiString(strconv.Itoa(i)))
	***REMOVED***

	return accum
***REMOVED***

func (o *objectGoSlice) export() interface***REMOVED******REMOVED*** ***REMOVED***
	return *o.data
***REMOVED***

func (o *objectGoSlice) exportType() reflect.Type ***REMOVED***
	return reflectTypeArray
***REMOVED***

func (o *objectGoSlice) equal(other objectImpl) bool ***REMOVED***
	if other, ok := other.(*objectGoSlice); ok ***REMOVED***
		return o.data == other.data
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSlice) sortLen() int64 ***REMOVED***
	return int64(len(*o.data))
***REMOVED***

func (o *objectGoSlice) sortGet(i int64) Value ***REMOVED***
	return o.getIdx(valueInt(i), nil)
***REMOVED***

func (o *objectGoSlice) swap(i, j int64) ***REMOVED***
	ii := valueInt(i)
	jj := valueInt(j)
	x := o.getIdx(ii, nil)
	y := o.getIdx(jj, nil)

	o.setOwnIdx(ii, y, false)
	o.setOwnIdx(jj, x, false)
***REMOVED***
