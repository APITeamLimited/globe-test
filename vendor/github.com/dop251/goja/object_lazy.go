package goja

import (
	"reflect"

	"github.com/dop251/goja/unistring"
)

type lazyObject struct ***REMOVED***
	val    *Object
	create func(*Object) objectImpl
***REMOVED***

func (o *lazyObject) className() string ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.className()
***REMOVED***

func (o *lazyObject) getIdx(p valueInt, receiver Value) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getIdx(p, receiver)
***REMOVED***

func (o *lazyObject) getSym(p *Symbol, receiver Value) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getSym(p, receiver)
***REMOVED***

func (o *lazyObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getOwnPropIdx(idx)
***REMOVED***

func (o *lazyObject) getOwnPropSym(s *Symbol) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getOwnPropSym(s)
***REMOVED***

func (o *lazyObject) hasPropertyIdx(idx valueInt) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasPropertyIdx(idx)
***REMOVED***

func (o *lazyObject) hasPropertySym(s *Symbol) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasPropertySym(s)
***REMOVED***

func (o *lazyObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasOwnPropertyIdx(idx)
***REMOVED***

func (o *lazyObject) hasOwnPropertySym(s *Symbol) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasOwnPropertySym(s)
***REMOVED***

func (o *lazyObject) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.defineOwnPropertyStr(name, desc, throw)
***REMOVED***

func (o *lazyObject) defineOwnPropertyIdx(name valueInt, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.defineOwnPropertyIdx(name, desc, throw)
***REMOVED***

func (o *lazyObject) defineOwnPropertySym(name *Symbol, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.defineOwnPropertySym(name, desc, throw)
***REMOVED***

func (o *lazyObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.deleteIdx(idx, throw)
***REMOVED***

func (o *lazyObject) deleteSym(s *Symbol, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.deleteSym(s, throw)
***REMOVED***

func (o *lazyObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getStr(name, receiver)
***REMOVED***

func (o *lazyObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.getOwnPropStr(name)
***REMOVED***

func (o *lazyObject) setOwnStr(p unistring.String, v Value, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.setOwnStr(p, v, throw)
***REMOVED***

func (o *lazyObject) setOwnIdx(p valueInt, v Value, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.setOwnIdx(p, v, throw)
***REMOVED***

func (o *lazyObject) setOwnSym(p *Symbol, v Value, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.setOwnSym(p, v, throw)
***REMOVED***

func (o *lazyObject) setForeignStr(p unistring.String, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.setForeignStr(p, v, receiver, throw)
***REMOVED***

func (o *lazyObject) setForeignIdx(p valueInt, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.setForeignIdx(p, v, receiver, throw)
***REMOVED***

func (o *lazyObject) setForeignSym(p *Symbol, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.setForeignSym(p, v, receiver, throw)
***REMOVED***

func (o *lazyObject) hasPropertyStr(name unistring.String) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasPropertyStr(name)
***REMOVED***

func (o *lazyObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.hasOwnPropertyStr(name)
***REMOVED***

func (o *lazyObject) _putProp(unistring.String, Value, bool, bool, bool) Value ***REMOVED***
	panic("cannot use _putProp() in lazy object")
***REMOVED***

func (o *lazyObject) _putSym(*Symbol, Value) ***REMOVED***
	panic("cannot use _putSym() in lazy object")
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

func (o *lazyObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.assertConstructor()
***REMOVED***

func (o *lazyObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.deleteStr(name, throw)
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

func (o *lazyObject) preventExtensions(throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.preventExtensions(throw)
***REMOVED***

func (o *lazyObject) iterateStringKeys() iterNextFunc ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.iterateStringKeys()
***REMOVED***

func (o *lazyObject) iterateSymbols() iterNextFunc ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.iterateSymbols()
***REMOVED***

func (o *lazyObject) iterateKeys() iterNextFunc ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.iterateKeys()
***REMOVED***

func (o *lazyObject) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.export(ctx)
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

func (o *lazyObject) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.stringKeys(all, accum)
***REMOVED***

func (o *lazyObject) symbols(all bool, accum []Value) []Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.symbols(all, accum)
***REMOVED***

func (o *lazyObject) keys(all bool, accum []Value) []Value ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.keys(all, accum)
***REMOVED***

func (o *lazyObject) setProto(proto *Object, throw bool) bool ***REMOVED***
	obj := o.create(o.val)
	o.val.self = obj
	return obj.setProto(proto, throw)
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
