package goja

import "reflect"

type lazyObject struct ***REMOVED***
	val    *Object
	create func(*Object) objectImpl
***REMOVED***

func (o *lazyObject) className() string ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.className()
***REMOVED***

func (o *lazyObject) get(n Value) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.get(n)
***REMOVED***

func (o *lazyObject) getProp(n Value) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getProp(n)
***REMOVED***

func (o *lazyObject) getPropStr(name string) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getPropStr(name)
***REMOVED***

func (o *lazyObject) getStr(name string) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getStr(name)
***REMOVED***

func (o *lazyObject) getOwnProp(name string) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getOwnProp(name)
***REMOVED***

func (o *lazyObject) put(n Value, val Value, throw bool) ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	obj.put(n, val, throw)
***REMOVED***

func (o *lazyObject) putStr(name string, val Value, throw bool) ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	obj.putStr(name, val, throw)
***REMOVED***

func (o *lazyObject) hasProperty(n Value) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasProperty(n)
***REMOVED***

func (o *lazyObject) hasPropertyStr(name string) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasPropertyStr(name)
***REMOVED***

func (o *lazyObject) hasOwnProperty(n Value) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasOwnProperty(n)
***REMOVED***

func (o *lazyObject) hasOwnPropertyStr(name string) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasOwnPropertyStr(name)
***REMOVED***

func (o *lazyObject) _putProp(name string, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj._putProp(name, value, writable, enumerable, configurable)
***REMOVED***

func (o *lazyObject) defineOwnProperty(name Value, descr propertyDescr, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.defineOwnProperty(name, descr, throw)
***REMOVED***

func (o *lazyObject) toPrimitiveNumber() Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.toPrimitiveNumber()
***REMOVED***

func (o *lazyObject) toPrimitiveString() Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.toPrimitiveString()
***REMOVED***

func (o *lazyObject) toPrimitive() Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.toPrimitive()
***REMOVED***

func (o *lazyObject) assertCallable() (call func(FunctionCall) Value, ok bool) ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.assertCallable()
***REMOVED***

func (o *lazyObject) deleteStr(name string, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.deleteStr(name, throw)
***REMOVED***

func (o *lazyObject) delete(name Value, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.delete(name, throw)
***REMOVED***

func (o *lazyObject) proto() *Object ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.proto()
***REMOVED***

func (o *lazyObject) hasInstance(v Value) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasInstance(v)
***REMOVED***

func (o *lazyObject) isExtensible() bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.isExtensible()
***REMOVED***

func (o *lazyObject) preventExtensions() ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	obj.preventExtensions()
***REMOVED***

func (o *lazyObject) enumerate(all, recusrive bool) iterNextFunc ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.enumerate(all, recusrive)
***REMOVED***

func (o *lazyObject) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj._enumerate(recursive)
***REMOVED***

func (o *lazyObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.export()
***REMOVED***

func (o *lazyObject) exportType() reflect.Type ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.exportType()
***REMOVED***

func (o *lazyObject) equal(other objectImpl) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.equal(other)
***REMOVED***

func (o *lazyObject) sortLen() int64 ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.sortLen()
***REMOVED***

func (o *lazyObject) sortGet(i int64) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.sortGet(i)
***REMOVED***

func (o *lazyObject) swap(i, j int64) ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	obj.swap(i, j)
***REMOVED***
