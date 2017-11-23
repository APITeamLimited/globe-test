package goja

import (
	"reflect"
	"strconv"
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
	o._setLen()
	o.baseObject._put("length", &o.lengthProp)
***REMOVED***

func (o *objectGoSlice) _setLen() ***REMOVED***
	o.lengthProp.value = intToValue(int64(len(*o.data)))
***REMOVED***

func (o *objectGoSlice) getIdx(idx int64) Value ***REMOVED***
	if idx < int64(len(*o.data)) ***REMOVED***
		return o.val.runtime.ToValue((*o.data)[idx])
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSlice) _get(n Value) Value ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return o.getIdx(idx)
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSlice) _getStr(name string) Value ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		return o.getIdx(idx)
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSlice) get(n Value) Value ***REMOVED***
	if v := o._get(n); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject._getStr(n.String())
***REMOVED***

func (o *objectGoSlice) getStr(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject._getStr(name)
***REMOVED***

func (o *objectGoSlice) getProp(n Value) Value ***REMOVED***
	if v := o._get(n); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject.getPropStr(n.String())
***REMOVED***

func (o *objectGoSlice) getPropStr(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject.getPropStr(name)
***REMOVED***

func (o *objectGoSlice) getOwnProp(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return &valueProperty***REMOVED***
			value:      v,
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return o.baseObject.getOwnProp(name)
***REMOVED***

func (o *objectGoSlice) grow(size int64) ***REMOVED***
	newcap := int64(cap(*o.data))
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
		*o.data = (*o.data)[:size]
	***REMOVED***
	o._setLen()
***REMOVED***

func (o *objectGoSlice) putIdx(idx int64, v Value, throw bool) ***REMOVED***
	if idx >= int64(len(*o.data)) ***REMOVED***
		if !o.sliceExtensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot extend Go slice")
			return
		***REMOVED***
		o.grow(idx + 1)
	***REMOVED***
	(*o.data)[idx] = v.Export()
***REMOVED***

func (o *objectGoSlice) put(n Value, val Value, throw bool) ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		o.putIdx(idx, val, throw)
		return
	***REMOVED***
	// TODO: length
	o.baseObject.put(n, val, throw)
***REMOVED***

func (o *objectGoSlice) putStr(name string, val Value, throw bool) ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		o.putIdx(idx, val, throw)
		return
	***REMOVED***
	// TODO: length
	o.baseObject.putStr(name, val, throw)
***REMOVED***

func (o *objectGoSlice) _has(n Value) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return idx < int64(len(*o.data))
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSlice) _hasStr(name string) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		return idx < int64(len(*o.data))
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSlice) hasProperty(n Value) bool ***REMOVED***
	if o._has(n) ***REMOVED***
		return true
	***REMOVED***
	return o.baseObject.hasProperty(n)
***REMOVED***

func (o *objectGoSlice) hasPropertyStr(name string) bool ***REMOVED***
	if o._hasStr(name) ***REMOVED***
		return true
	***REMOVED***
	return o.baseObject.hasPropertyStr(name)
***REMOVED***

func (o *objectGoSlice) hasOwnProperty(n Value) bool ***REMOVED***
	if o._has(n) ***REMOVED***
		return true
	***REMOVED***
	return o.baseObject.hasOwnProperty(n)
***REMOVED***

func (o *objectGoSlice) hasOwnPropertyStr(name string) bool ***REMOVED***
	if o._hasStr(name) ***REMOVED***
		return true
	***REMOVED***
	return o.baseObject.hasOwnPropertyStr(name)
***REMOVED***

func (o *objectGoSlice) _putProp(name string, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	o.putStr(name, value, false)
	return value
***REMOVED***

func (o *objectGoSlice) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		if !o.val.runtime.checkHostObjectPropertyDescr(n.String(), descr, throw) ***REMOVED***
			return false
		***REMOVED***
		val := descr.Value
		if val == nil ***REMOVED***
			val = _undefined
		***REMOVED***
		o.putIdx(idx, val, throw)
		return true
	***REMOVED***
	return o.baseObject.defineOwnProperty(n, descr, throw)
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

func (o *objectGoSlice) deleteStr(name string, throw bool) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 && idx < int64(len(*o.data)) ***REMOVED***
		(*o.data)[idx] = nil
		return true
	***REMOVED***
	return o.baseObject.deleteStr(name, throw)
***REMOVED***

func (o *objectGoSlice) delete(name Value, throw bool) bool ***REMOVED***
	if idx := toIdx(name); idx >= 0 && idx < int64(len(*o.data)) ***REMOVED***
		(*o.data)[idx] = nil
		return true
	***REMOVED***
	return o.baseObject.delete(name, throw)
***REMOVED***

type goslicePropIter struct ***REMOVED***
	o          *objectGoSlice
	recursive  bool
	idx, limit int
***REMOVED***

func (i *goslicePropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.limit && i.idx < len(*i.o.data) ***REMOVED***
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: name, enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	if i.recursive ***REMOVED***
		return i.o.prototype.self._enumerate(i.recursive)()
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoSlice) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next

***REMOVED***

func (o *objectGoSlice) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	return (&goslicePropIter***REMOVED***
		o:         o,
		recursive: recursive,
		limit:     len(*o.data),
	***REMOVED***).next
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
	return o.get(intToValue(i))
***REMOVED***

func (o *objectGoSlice) swap(i, j int64) ***REMOVED***
	ii := intToValue(i)
	jj := intToValue(j)
	x := o.get(ii)
	y := o.get(jj)

	o.put(ii, y, false)
	o.put(jj, x, false)
***REMOVED***
