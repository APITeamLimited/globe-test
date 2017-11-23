package goja

import "reflect"

type objectGoMapReflect struct ***REMOVED***
	objectGoReflect

	keyType, valueType reflect.Type
***REMOVED***

func (o *objectGoMapReflect) init() ***REMOVED***
	o.objectGoReflect.init()
	o.keyType = o.value.Type().Key()
	o.valueType = o.value.Type().Elem()
***REMOVED***

func (o *objectGoMapReflect) toKey(n Value) reflect.Value ***REMOVED***
	key, err := o.val.runtime.toReflectValue(n, o.keyType)
	if err != nil ***REMOVED***
		o.val.runtime.typeErrorResult(true, "map key conversion error: %v", err)
		panic("unreachable")
	***REMOVED***
	return key
***REMOVED***

func (o *objectGoMapReflect) strToKey(name string) reflect.Value ***REMOVED***
	if o.keyType.Kind() == reflect.String ***REMOVED***
		return reflect.ValueOf(name).Convert(o.keyType)
	***REMOVED***
	return o.toKey(newStringValue(name))
***REMOVED***

func (o *objectGoMapReflect) _get(n Value) Value ***REMOVED***
	if v := o.value.MapIndex(o.toKey(n)); v.IsValid() ***REMOVED***
		return o.val.runtime.ToValue(v.Interface())
	***REMOVED***

	return nil
***REMOVED***

func (o *objectGoMapReflect) _getStr(name string) Value ***REMOVED***
	if v := o.value.MapIndex(o.strToKey(name)); v.IsValid() ***REMOVED***
		return o.val.runtime.ToValue(v.Interface())
	***REMOVED***

	return nil
***REMOVED***

func (o *objectGoMapReflect) get(n Value) Value ***REMOVED***
	if v := o._get(n); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.get(n)
***REMOVED***

func (o *objectGoMapReflect) getStr(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return v
	***REMOVED***
	return o.objectGoReflect.getStr(name)
***REMOVED***

func (o *objectGoMapReflect) getProp(n Value) Value ***REMOVED***
	return o.get(n)
***REMOVED***

func (o *objectGoMapReflect) getPropStr(name string) Value ***REMOVED***
	return o.getStr(name)
***REMOVED***

func (o *objectGoMapReflect) getOwnProp(name string) Value ***REMOVED***
	if v := o._getStr(name); v != nil ***REMOVED***
		return &valueProperty***REMOVED***
			value:      v,
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return o.objectGoReflect.getOwnProp(name)
***REMOVED***

func (o *objectGoMapReflect) toValue(val Value, throw bool) (reflect.Value, bool) ***REMOVED***
	v, err := o.val.runtime.toReflectValue(val, o.valueType)
	if err != nil ***REMOVED***
		o.val.runtime.typeErrorResult(throw, "map value conversion error: %v", err)
		return reflect.Value***REMOVED******REMOVED***, false
	***REMOVED***

	return v, true
***REMOVED***

func (o *objectGoMapReflect) put(key, val Value, throw bool) ***REMOVED***
	k := o.toKey(key)
	v, ok := o.toValue(val, throw)
	if !ok ***REMOVED***
		return
	***REMOVED***
	o.value.SetMapIndex(k, v)
***REMOVED***

func (o *objectGoMapReflect) putStr(name string, val Value, throw bool) ***REMOVED***
	k := o.strToKey(name)
	v, ok := o.toValue(val, throw)
	if !ok ***REMOVED***
		return
	***REMOVED***
	o.value.SetMapIndex(k, v)
***REMOVED***

func (o *objectGoMapReflect) _putProp(name string, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	o.putStr(name, value, true)
	return value
***REMOVED***

func (o *objectGoMapReflect) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool ***REMOVED***
	name := n.String()
	if !o.val.runtime.checkHostObjectPropertyDescr(name, descr, throw) ***REMOVED***
		return false
	***REMOVED***

	o.put(n, descr.Value, throw)
	return true
***REMOVED***

func (o *objectGoMapReflect) hasOwnPropertyStr(name string) bool ***REMOVED***
	return o.value.MapIndex(o.strToKey(name)).IsValid()
***REMOVED***

func (o *objectGoMapReflect) hasOwnProperty(n Value) bool ***REMOVED***
	return o.value.MapIndex(o.toKey(n)).IsValid()
***REMOVED***

func (o *objectGoMapReflect) hasProperty(n Value) bool ***REMOVED***
	if o.hasOwnProperty(n) ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect.hasProperty(n)
***REMOVED***

func (o *objectGoMapReflect) hasPropertyStr(name string) bool ***REMOVED***
	if o.hasOwnPropertyStr(name) ***REMOVED***
		return true
	***REMOVED***
	return o.objectGoReflect.hasPropertyStr(name)
***REMOVED***

func (o *objectGoMapReflect) delete(n Value, throw bool) bool ***REMOVED***
	o.value.SetMapIndex(o.toKey(n), reflect.Value***REMOVED******REMOVED***)
	return true
***REMOVED***

func (o *objectGoMapReflect) deleteStr(name string, throw bool) bool ***REMOVED***
	o.value.SetMapIndex(o.strToKey(name), reflect.Value***REMOVED******REMOVED***)
	return true
***REMOVED***

type gomapReflectPropIter struct ***REMOVED***
	o         *objectGoMapReflect
	keys      []reflect.Value
	idx       int
	recursive bool
***REMOVED***

func (i *gomapReflectPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.keys) ***REMOVED***
		key := i.keys[i.idx]
		v := i.o.value.MapIndex(key)
		i.idx++
		if v.IsValid() ***REMOVED***
			return propIterItem***REMOVED***name: key.String(), enumerable: _ENUM_TRUE***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	if i.recursive ***REMOVED***
		return i.o.objectGoReflect._enumerate(true)()
	***REMOVED***

	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (o *objectGoMapReflect) _enumerate(recusrive bool) iterNextFunc ***REMOVED***
	r := &gomapReflectPropIter***REMOVED***
		o:         o,
		keys:      o.value.MapKeys(),
		recursive: recusrive,
	***REMOVED***
	return r.next
***REMOVED***

func (o *objectGoMapReflect) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: o._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next
***REMOVED***

func (o *objectGoMapReflect) equal(other objectImpl) bool ***REMOVED***
	if other, ok := other.(*objectGoMapReflect); ok ***REMOVED***
		return o.value.Interface() == other.value.Interface()
	***REMOVED***
	return false
***REMOVED***
