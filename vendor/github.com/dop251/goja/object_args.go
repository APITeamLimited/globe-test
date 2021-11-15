package goja

import "github.com/dop251/goja/unistring"

type argumentsObject struct ***REMOVED***
	baseObject
	length int
***REMOVED***

type mappedProperty struct ***REMOVED***
	valueProperty
	v *Value
***REMOVED***

func (a *argumentsObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	return a.getStrWithOwnProp(a.getOwnPropStr(name), name, receiver)
***REMOVED***

func (a *argumentsObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if mapped, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		if mapped.writable && mapped.enumerable && mapped.configurable ***REMOVED***
			return *mapped.v
		***REMOVED***
		return &valueProperty***REMOVED***
			value:        *mapped.v,
			writable:     mapped.writable,
			configurable: mapped.configurable,
			enumerable:   mapped.enumerable,
		***REMOVED***
	***REMOVED***

	return a.baseObject.getOwnPropStr(name)
***REMOVED***

func (a *argumentsObject) init() ***REMOVED***
	a.baseObject.init()
	a._putProp("length", intToValue(int64(a.length)), true, false, true)
***REMOVED***

func (a *argumentsObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if prop, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		if !prop.writable ***REMOVED***
			a.val.runtime.typeErrorResult(throw, "Property is not writable: %s", name)
			return false
		***REMOVED***
		*prop.v = val
		return true
	***REMOVED***
	return a.baseObject.setOwnStr(name, val, throw)
***REMOVED***

func (a *argumentsObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return a._setForeignStr(name, a.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

func (a *argumentsObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if prop, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		if !a.checkDeleteProp(name, &prop.valueProperty, throw) ***REMOVED***
			return false
		***REMOVED***
		a._delete(name)
		return true
	***REMOVED***

	return a.baseObject.deleteStr(name, throw)
***REMOVED***

type argumentsPropIter struct ***REMOVED***
	wrapped iterNextFunc
***REMOVED***

func (i *argumentsPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	var item propIterItem
	item, i.wrapped = i.wrapped()
	if i.wrapped == nil ***REMOVED***
		return propIterItem***REMOVED******REMOVED***, nil
	***REMOVED***
	if prop, ok := item.value.(*mappedProperty); ok ***REMOVED***
		item.value = *prop.v
	***REMOVED***
	return item, i.next
***REMOVED***

func (a *argumentsObject) iterateStringKeys() iterNextFunc ***REMOVED***
	return (&argumentsPropIter***REMOVED***
		wrapped: a.baseObject.iterateStringKeys(),
	***REMOVED***).next
***REMOVED***

func (a *argumentsObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if mapped, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		existing := &valueProperty***REMOVED***
			configurable: mapped.configurable,
			writable:     true,
			enumerable:   mapped.enumerable,
			value:        *mapped.v,
		***REMOVED***

		val, ok := a.baseObject._defineOwnProperty(name, existing, descr, throw)
		if !ok ***REMOVED***
			return false
		***REMOVED***

		if prop, ok := val.(*valueProperty); ok ***REMOVED***
			if !prop.accessor ***REMOVED***
				*mapped.v = prop.value
			***REMOVED***
			if prop.accessor || !prop.writable ***REMOVED***
				a._put(name, prop)
				return true
			***REMOVED***
			mapped.configurable = prop.configurable
			mapped.enumerable = prop.enumerable
		***REMOVED*** else ***REMOVED***
			*mapped.v = val
			mapped.configurable = true
			mapped.enumerable = true
		***REMOVED***

		return true
	***REMOVED***

	return a.baseObject.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (a *argumentsObject) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, exists := ctx.get(a); exists ***REMOVED***
		return v
	***REMOVED***
	arr := make([]interface***REMOVED******REMOVED***, a.length)
	ctx.put(a, arr)
	for i := range arr ***REMOVED***
		v := a.getIdx(valueInt(int64(i)), nil)
		if v != nil ***REMOVED***
			arr[i] = exportValue(v, ctx)
		***REMOVED***
	***REMOVED***
	return arr
***REMOVED***
