package goja

import (
	"reflect"

	"github.com/dop251/goja/unistring"
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

func (o *objectGoMapSimple) _getStr(name string) Value ***REMOVED***
	v, exists := o.data[name]
	if !exists ***REMOVED***
		return nil
	***REMOVED***
	return o.val.runtime.ToValue(v)
***REMOVED***

func (o *objectGoMapSimple) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	if v := o._getStr(name.String()); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.baseObject.getStr(name, receiver)
***REMOVED***

func (o *objectGoMapSimple) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if v := o._getStr(name.String()); v != nil ***REMOVED***
		return v
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoMapSimple) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	n := name.String()
	if _, exists := o.data[n]; exists ***REMOVED***
		o.data[n] = val.Export()
		return true
	***REMOVED***
	if proto := o.prototype; proto != nil ***REMOVED***
		// we know it's foreign because prototype loops are not allowed
		if res, ok := proto.self.setForeignStr(name, val, o.val, throw); ok ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	// new property
	if !o.extensible ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "Cannot add property %s, object is not extensible", name)
		return false
	***REMOVED*** else ***REMOVED***
		o.data[n] = val.Export()
	***REMOVED***
	return true
***REMOVED***

func trueValIfPresent(present bool) Value ***REMOVED***
	if present ***REMOVED***
		return valueTrue
	***REMOVED***
	return nil
***REMOVED***

func (o *objectGoMapSimple) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return o._setForeignStr(name, trueValIfPresent(o._hasStr(name.String())), val, receiver, throw)
***REMOVED***

func (o *objectGoMapSimple) _hasStr(name string) bool ***REMOVED***
	_, exists := o.data[name]
	return exists
***REMOVED***

func (o *objectGoMapSimple) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	return o._hasStr(name.String())
***REMOVED***

func (o *objectGoMapSimple) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if !o.val.runtime.checkHostObjectPropertyDescr(name, descr, throw) ***REMOVED***
		return false
	***REMOVED***

	n := name.String()
	if o.extensible || o._hasStr(n) ***REMOVED***
		o.data[n] = descr.Value.Export()
		return true
	***REMOVED***

	o.val.runtime.typeErrorResult(throw, "Cannot define property %s, object is not extensible", n)
	return false
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

func (o *objectGoMapSimple) deleteStr(name unistring.String, _ bool) bool ***REMOVED***
	delete(o.data, name.String())
	return true
***REMOVED***

type gomapPropIter struct ***REMOVED***
	o         *objectGoMapSimple
	propNames []string
	idx       int
***REMOVED***

func (i *gomapPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.propNames) ***REMOVED***
		name := i.propNames[i.idx]
		i.idx++
		if _, exists := i.o.data[name]; exists ***REMOVED***
			return propIterItem***REMOVED***name: newStringValue(name), enumerable: _ENUM_TRUE***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoMapSimple) iterateStringKeys() iterNextFunc ***REMOVED***
	propNames := make([]string, len(o.data))
	i := 0
	for key := range o.data ***REMOVED***
		propNames[i] = key
		i++
	***REMOVED***

	return (&gomapPropIter***REMOVED***
		o:         o,
		propNames: propNames,
	***REMOVED***).next
***REMOVED***

func (o *objectGoMapSimple) stringKeys(_ bool, accum []Value) []Value ***REMOVED***
	// all own keys are enumerable
	for key := range o.data ***REMOVED***
		accum = append(accum, newStringValue(key))
	***REMOVED***
	return accum
***REMOVED***

func (o *objectGoMapSimple) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
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
