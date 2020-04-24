package goja

import (
	"reflect"
	"strconv"
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
	o._setLen()
	o.baseObject._put("length", &o.lengthProp)
***REMOVED***

func (o *objectGoSliceReflect) _setLen() ***REMOVED***
	o.lengthProp.value = intToValue(int64(o.value.Len()))
***REMOVED***

func (o *objectGoSliceReflect) _has(n Value) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return idx < int64(o.value.Len())
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSliceReflect) _hasStr(name string) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		return idx < int64(o.value.Len())
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoSliceReflect) getIdx(idx int64) Value ***REMOVED***
	if idx < int64(o.value.Len()) ***REMOVED***
		return o.val.runtime.ToValue(o.value.Index(int(idx)).Interface())
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSliceReflect) _get(n Value) Value ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return o.getIdx(idx)
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSliceReflect) _getStr(name string) Value ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		return o.getIdx(idx)
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoSliceReflect) get(n Value) Value ***REMOVED***
	if v := o._get(n); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.get(n)
***REMOVED***

func (o *objectGoSliceReflect) getStr(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.getStr(name)
***REMOVED***

func (o *objectGoSliceReflect) getProp(n Value) Value ***REMOVED***
	if v := o._get(n); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.getProp(n)
***REMOVED***

func (o *objectGoSliceReflect) getPropStr(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.getPropStr(name)
***REMOVED***

func (o *objectGoSliceReflect) getOwnProp(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.getOwnProp(name)
***REMOVED***

func (o *objectGoSliceReflect) putIdx(idx int64, v Value, throw bool) ***REMOVED***
	if idx >= int64(o.value.Len()) ***REMOVED***
		if !o.sliceExtensible ***REMOVED***
			o.val.runtime.typeErrorResult(throw, "Cannot extend a Go unaddressable reflect slice")
			return
		***REMOVED***
		o.grow(int(idx + 1))
	***REMOVED***
	val, err := o.val.runtime.toReflectValue(v, o.value.Type().Elem())
	if err != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Go type conversion error: %v", err)
		return
	***REMOVED***
	o.value.Index(int(idx)).Set(val)
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
		o.value.SetLen(size)
	***REMOVED***
	o._setLen()
***REMOVED***

func (o *objectGoSliceReflect) put(n Value, val Value, throw bool) ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		o.putIdx(idx, val, throw)
		return
	***REMOVED***
	// TODO: length
	o.objectGoReflect.put(n, val, throw)
***REMOVED***

func (o *objectGoSliceReflect) putStr(name string, val Value, throw bool) ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		o.putIdx(idx, val, throw)
		return
	***REMOVED***
	if name == "length" ***REMOVED***
		o.baseObject.putStr(name, val, throw)
		return
	***REMOVED***
	o.objectGoReflect.putStr(name, val, throw)
***REMOVED***

func (o *objectGoSliceReflect) hasProperty(n Value) bool ***REMOVED***
	if o._has(n) ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect.hasProperty(n)
***REMOVED***

func (o *objectGoSliceReflect) hasPropertyStr(name string) bool ***REMOVED***
	if o._hasStr(name) ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect.hasOwnPropertyStr(name)
***REMOVED***

func (o *objectGoSliceReflect) hasOwnProperty(n Value) bool ***REMOVED***
	if o._has(n) ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect.hasOwnProperty(n)
***REMOVED***

func (o *objectGoSliceReflect) hasOwnPropertyStr(name string) bool ***REMOVED***
	if o._hasStr(name) ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect.hasOwnPropertyStr(name)
***REMOVED***

func (o *objectGoSliceReflect) _putProp(name string, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	o.putStr(name, value, false)
	return value
***REMOVED***

func (o *objectGoSliceReflect) defineOwnProperty(name Value, descr propertyDescr, throw bool) bool ***REMOVED***
	if !o.val.runtime.checkHostObjectPropertyDescr(name.String(), descr, throw) ***REMOVED***
		return false
	***REMOVED***
	o.put(name, descr.Value, throw)
	return true
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

func (o *objectGoSliceReflect) deleteStr(name string, throw bool) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 && idx < int64(o.value.Len()) ***REMOVED***
		o.value.Index(int(idx)).Set(reflect.Zero(o.value.Type().Elem()))
		return true
	***REMOVED***
	return o.objectGoReflect.deleteStr(name, throw)
***REMOVED***

func (o *objectGoSliceReflect) delete(name Value, throw bool) bool ***REMOVED***
	if idx := toIdx(name); idx >= 0 && idx < int64(o.value.Len()) ***REMOVED***
		o.value.Index(int(idx)).Set(reflect.Zero(o.value.Type().Elem()))
		return true
	***REMOVED***
	return o.objectGoReflect.delete(name, throw)
***REMOVED***

type gosliceReflectPropIter struct ***REMOVED***
	o          *objectGoSliceReflect
	recursive  bool
	idx, limit int
***REMOVED***

func (i *gosliceReflectPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.limit && i.idx < i.o.value.Len() ***REMOVED***
		name := strconv.Itoa(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: name, enumerable: _ENUM_TRUE***REMOVED***, i.next
	***REMOVED***

	if i.recursive ***REMOVED***
		return i.o.prototype.self._enumerate(i.recursive)()
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoSliceReflect) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next
***REMOVED***

func (o *objectGoSliceReflect) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	return (&gosliceReflectPropIter***REMOVED***
		o:         o,
		recursive: recursive,
		limit:     o.value.Len(),
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
	return o.get(intToValue(i))
***REMOVED***

func (o *objectGoSliceReflect) swap(i, j int64) ***REMOVED***
	ii := intToValue(i)
	jj := intToValue(j)
	x := o.get(ii)
	y := o.get(jj)

	o.put(ii, y, false)
	o.put(jj, x, false)
***REMOVED***
