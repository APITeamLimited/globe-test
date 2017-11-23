package goja

type argumentsObject struct ***REMOVED***
	baseObject
	length int
***REMOVED***

type mappedProperty struct ***REMOVED***
	valueProperty
	v *Value
***REMOVED***

func (a *argumentsObject) getPropStr(name string) Value ***REMOVED***
	if prop, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		return *prop.v
	***REMOVED***
	return a.baseObject.getPropStr(name)
***REMOVED***

func (a *argumentsObject) getProp(n Value) Value ***REMOVED***
	return a.getPropStr(n.String())
***REMOVED***

func (a *argumentsObject) init() ***REMOVED***
	a.baseObject.init()
	a._putProp("length", intToValue(int64(a.length)), true, false, true)
***REMOVED***

func (a *argumentsObject) put(n Value, val Value, throw bool) ***REMOVED***
	a.putStr(n.String(), val, throw)
***REMOVED***

func (a *argumentsObject) putStr(name string, val Value, throw bool) ***REMOVED***
	if prop, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		if !prop.writable ***REMOVED***
			a.val.runtime.typeErrorResult(throw, "Property is not writable: %s", name)
			return
		***REMOVED***
		*prop.v = val
		return
	***REMOVED***
	a.baseObject.putStr(name, val, throw)
***REMOVED***

func (a *argumentsObject) deleteStr(name string, throw bool) bool ***REMOVED***
	if prop, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		if !a.checkDeleteProp(name, &prop.valueProperty, throw) ***REMOVED***
			return false
		***REMOVED***
		a._delete(name)
		return true
	***REMOVED***

	return a.baseObject.deleteStr(name, throw)
***REMOVED***

func (a *argumentsObject) delete(n Value, throw bool) bool ***REMOVED***
	return a.deleteStr(n.String(), throw)
***REMOVED***

type argumentsPropIter1 struct ***REMOVED***
	a         *argumentsObject
	idx       int
	recursive bool
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

func (a *argumentsObject) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	return (&argumentsPropIter***REMOVED***
		wrapped: a.baseObject._enumerate(recursive),
	***REMOVED***).next

***REMOVED***

func (a *argumentsObject) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&argumentsPropIter***REMOVED***
		wrapped: a.baseObject.enumerate(all, recursive),
	***REMOVED***).next
***REMOVED***

func (a *argumentsObject) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool ***REMOVED***
	name := n.String()
	if mapped, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		existing := &valueProperty***REMOVED***
			configurable: mapped.configurable,
			writable:     true,
			enumerable:   mapped.enumerable,
			value:        mapped.get(a.val),
		***REMOVED***

		val, ok := a.baseObject._defineOwnProperty(n, existing, descr, throw)
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

	return a.baseObject.defineOwnProperty(n, descr, throw)
***REMOVED***

func (a *argumentsObject) getOwnProp(name string) Value ***REMOVED***
	if mapped, ok := a.values[name].(*mappedProperty); ok ***REMOVED***
		return *mapped.v
	***REMOVED***

	return a.baseObject.getOwnProp(name)
***REMOVED***

func (a *argumentsObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	arr := make([]interface***REMOVED******REMOVED***, a.length)
	for i, _ := range arr ***REMOVED***
		v := a.get(intToValue(int64(i)))
		if v != nil ***REMOVED***
			arr[i] = v.Export()
		***REMOVED***
	***REMOVED***
	return arr
***REMOVED***
