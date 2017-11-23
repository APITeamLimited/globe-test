package goja

import (
	"reflect"
	"strconv"
)

type objectGoMapSimple struct ***REMOVED***
	baseObject
	data map[string]interface***REMOVED******REMOVED***
***REMOVED***

func (o *objectGoMapSimple) init() ***REMOVED***
	o.baseObject.init()
	o.prototype = o.val.runtime.global.ObjectPrototype
	o.class = classObject
	o.extensible = true
***REMOVED***

func (o *objectGoMapSimple) _get(n Value) Value ***REMOVED***
	return o._getStr(n.String())
***REMOVED***

func (o *objectGoMapSimple) _getStr(name string) Value ***REMOVED***
	v, exists := o.data[name]
	if !exists ***REMOVED***
		return nil
	***REMOVED***
	return o.val.runtime.ToValue(v)
***REMOVED***

func (o *objectGoMapSimple) get(n Value) Value ***REMOVED***
	return o.getStr(n.String())
***REMOVED***

func (o *objectGoMapSimple) getProp(n Value) Value ***REMOVED***
	return o.getPropStr(n.String())
***REMOVED***

func (o *objectGoMapSimple) getPropStr(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject.getPropStr(name)
***REMOVED***

func (o *objectGoMapSimple) getStr(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject._getStr(name)
***REMOVED***

func (o *objectGoMapSimple) getOwnProp(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject.getOwnProp(name)
***REMOVED***

func (o *objectGoMapSimple) put(n Value, val Value, throw bool) ***REMOVED***
	o.putStr(n.String(), val, throw)
***REMOVED***

func (o *objectGoMapSimple) _hasStr(name string) bool ***REMOVED***
	_, exists := o.data[name]
	return exists
***REMOVED***

func (o *objectGoMapSimple) _has(n Value) bool ***REMOVED***
	return o._hasStr(n.String())
***REMOVED***

func (o *objectGoMapSimple) putStr(name string, val Value, throw bool) ***REMOVED***
	if o.extensible || o._hasStr(name) ***REMOVED***
		o.data[name] = val.Export()
	***REMOVED*** else ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Host object is not extensible")
	***REMOVED***
***REMOVED***

func (o *objectGoMapSimple) hasProperty(n Value) bool ***REMOVED***
	if o._has(n) ***REMOVED***
		return true
	***REMOVED***
	return o.baseObject.hasProperty(n)
***REMOVED***

func (o *objectGoMapSimple) hasPropertyStr(name string) bool ***REMOVED***
	if o._hasStr(name) ***REMOVED***
		return true
	***REMOVED***
	return o.baseObject.hasOwnPropertyStr(name)
***REMOVED***

func (o *objectGoMapSimple) hasOwnProperty(n Value) bool ***REMOVED***
	return o._has(n)
***REMOVED***

func (o *objectGoMapSimple) hasOwnPropertyStr(name string) bool ***REMOVED***
	return o._hasStr(name)
***REMOVED***

func (o *objectGoMapSimple) _putProp(name string, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	o.putStr(name, value, false)
	return value
***REMOVED***

func (o *objectGoMapSimple) defineOwnProperty(name Value, descr propertyDescr, throw bool) bool ***REMOVED***
	if descr.Getter != nil || descr.Setter != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Host objects do not support accessor properties")
		return false
	***REMOVED***
	o.put(name, descr.Value, throw)
	return true
***REMOVED***

/*
func (o *objectGoMapSimple) toPrimitiveNumber() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoMapSimple) toPrimitiveString() Value ***REMOVED***
	return stringObjectObject
***REMOVED***

func (o *objectGoMapSimple) toPrimitive() Value ***REMOVED***
	return o.toPrimitiveString()
***REMOVED***

func (o *objectGoMapSimple) assertCallable() (call func(FunctionCall) Value, ok bool) ***REMOVED***
	return nil, false
***REMOVED***
*/

func (o *objectGoMapSimple) deleteStr(name string, throw bool) bool ***REMOVED***
	delete(o.data, name)
	return true
***REMOVED***

func (o *objectGoMapSimple) delete(name Value, throw bool) bool ***REMOVED***
	return o.deleteStr(name.String(), throw)
***REMOVED***

type gomapPropIter struct ***REMOVED***
	o         *objectGoMapSimple
	propNames []string
	recursive bool
	idx       int
***REMOVED***

func (i *gomapPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.propNames) ***REMOVED***
		name := i.propNames[i.idx]
		i.idx++
		if _, exists := i.o.data[name]; exists ***REMOVED***
			return propIterItem***REMOVED***name: name, enumerable: _ENUM_TRUE***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	if i.recursive ***REMOVED***
		return i.o.prototype.self._enumerate(true)()
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoMapSimple) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next
***REMOVED***

func (o *objectGoMapSimple) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	propNames := make([]string, len(o.data))
	i := 0
	for key, _ := range o.data ***REMOVED***
		propNames[i] = key
		i++
	***REMOVED***
	return (&gomapPropIter***REMOVED***
		o:         o,
		propNames: propNames,
		recursive: recursive,
	***REMOVED***).next
***REMOVED***

func (o *objectGoMapSimple) export() interface***REMOVED******REMOVED*** ***REMOVED***
	return o.data
***REMOVED***

func (o *objectGoMapSimple) exportType() reflect.Type ***REMOVED***
	return reflectTypeMap
***REMOVED***

func (o *objectGoMapSimple) equal(other objectImpl) bool ***REMOVED***
	if other, ok := other.(*objectGoMapSimple); ok ***REMOVED***
		return o == other
	***REMOVED***
	return false
***REMOVED***

func (o *objectGoMapSimple) sortLen() int64 ***REMOVED***
	return int64(len(o.data))
***REMOVED***

func (o *objectGoMapSimple) sortGet(i int64) Value ***REMOVED***
	return o.getStr(strconv.FormatInt(i, 10))
***REMOVED***

func (o *objectGoMapSimple) swap(i, j int64) ***REMOVED***
	ii := strconv.FormatInt(i, 10)
	jj := strconv.FormatInt(j, 10)
	x := o.getStr(ii)
	y := o.getStr(jj)

	o.putStr(ii, y, false)
	o.putStr(jj, x, false)
***REMOVED***
