package goja

import (
	"reflect"
	"strconv"

	"github.com/dop251/goja/unistring"
)

type objectGoArrayReflect struct ***REMOVED***
	objectGoReflect
	lengthProp valueProperty

	valueCache valueArrayCache

	putIdx func(idx int, v Value, throw bool) bool
***REMOVED***

type valueArrayCache []reflectValueWrapper

func (c *valueArrayCache) get(idx int) reflectValueWrapper ***REMOVED***
	if idx < len(*c) ***REMOVED***
		return (*c)[idx]
	***REMOVED***
	return nil
***REMOVED***

func (c *valueArrayCache) grow(newlen int) ***REMOVED***
	oldcap := cap(*c)
	if oldcap < newlen ***REMOVED***
		a := make([]reflectValueWrapper, newlen, growCap(newlen, len(*c), oldcap))
		copy(a, *c)
		*c = a
	***REMOVED*** else ***REMOVED***
		*c = (*c)[:newlen]
	***REMOVED***
***REMOVED***

func (c *valueArrayCache) put(idx int, w reflectValueWrapper) ***REMOVED***
	if len(*c) <= idx ***REMOVED***
		c.grow(idx + 1)
	***REMOVED***
	(*c)[idx] = w
***REMOVED***

func (c *valueArrayCache) shrink(newlen int) ***REMOVED***
	if len(*c) > newlen ***REMOVED***
		tail := (*c)[newlen:]
		for i, item := range tail ***REMOVED***
			if item != nil ***REMOVED***
				copyReflectValueWrapper(item)
				tail[i] = nil
			***REMOVED***
		***REMOVED***
		*c = (*c)[:newlen]
	***REMOVED***
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
	if v := o.valueCache.get(idx); v != nil ***REMOVED***
		return v.esValue()
	***REMOVED***

	v := o.value.Index(idx)

	res, w := o.elemToValue(v)
	if w != nil ***REMOVED***
		o.valueCache.put(idx, w)
	***REMOVED***

	return res
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
	cached := o.valueCache.get(idx)
	if cached != nil ***REMOVED***
		copyReflectValueWrapper(cached)
	***REMOVED***

	rv := o.value.Index(idx)
	err := o.val.runtime.toReflectValue(v, rv, &objectExportCtx***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		if cached != nil ***REMOVED***
			cached.setReflectValue(rv)
		***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Go type conversion error: %v", err)
		return false
	***REMOVED***
	if cached != nil ***REMOVED***
		o.valueCache[idx] = nil
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
		if cv := o.valueCache.get(idx); cv != nil ***REMOVED***
			copyReflectValueWrapper(cv)
			o.valueCache[idx] = nil
		***REMOVED***

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

func (o *objectGoArrayReflect) sortLen() int ***REMOVED***
	return o.value.Len()
***REMOVED***

func (o *objectGoArrayReflect) sortGet(i int) Value ***REMOVED***
	return o.getIdx(valueInt(i), nil)
***REMOVED***

func (o *objectGoArrayReflect) swap(i int, j int) ***REMOVED***
	vi := o.value.Index(i)
	vj := o.value.Index(j)
	tmp := reflect.New(o.value.Type().Elem()).Elem()
	tmp.Set(vi)
	vi.Set(vj)
	vj.Set(tmp)

	cachedI := o.valueCache.get(i)
	cachedJ := o.valueCache.get(j)
	if cachedI != nil ***REMOVED***
		cachedI.setReflectValue(vj)
		o.valueCache.put(j, cachedI)
	***REMOVED*** else ***REMOVED***
		if j < len(o.valueCache) ***REMOVED***
			o.valueCache[j] = nil
		***REMOVED***
	***REMOVED***

	if cachedJ != nil ***REMOVED***
		cachedJ.setReflectValue(vi)
		o.valueCache.put(i, cachedJ)
	***REMOVED*** else ***REMOVED***
		if i < len(o.valueCache) ***REMOVED***
			o.valueCache[i] = nil
		***REMOVED***
	***REMOVED***
***REMOVED***
