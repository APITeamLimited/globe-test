package goja

import (
	"github.com/dop251/goja/unistring"
	"reflect"
)

type destructKeyedSource struct ***REMOVED***
	r        *Runtime
	wrapped  Value
	usedKeys map[Value]struct***REMOVED******REMOVED***
***REMOVED***

func newDestructKeyedSource(r *Runtime, wrapped Value) *destructKeyedSource ***REMOVED***
	return &destructKeyedSource***REMOVED***
		r:       r,
		wrapped: wrapped,
	***REMOVED***
***REMOVED***

func (r *Runtime) newDestructKeyedSource(wrapped Value) *Object ***REMOVED***
	return &Object***REMOVED***
		runtime: r,
		self:    newDestructKeyedSource(r, wrapped),
	***REMOVED***
***REMOVED***

func (d *destructKeyedSource) w() objectImpl ***REMOVED***
	return d.wrapped.ToObject(d.r).self
***REMOVED***

func (d *destructKeyedSource) recordKey(key Value) ***REMOVED***
	if d.usedKeys == nil ***REMOVED***
		d.usedKeys = make(map[Value]struct***REMOVED******REMOVED***)
	***REMOVED***
	d.usedKeys[key] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***

func (d *destructKeyedSource) sortLen() int64 ***REMOVED***
	return d.w().sortLen()
***REMOVED***

func (d *destructKeyedSource) sortGet(i int64) Value ***REMOVED***
	return d.w().sortGet(i)
***REMOVED***

func (d *destructKeyedSource) swap(i int64, i2 int64) ***REMOVED***
	d.w().swap(i, i2)
***REMOVED***

func (d *destructKeyedSource) className() string ***REMOVED***
	return d.w().className()
***REMOVED***

func (d *destructKeyedSource) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	d.recordKey(stringValueFromRaw(p))
	return d.w().getStr(p, receiver)
***REMOVED***

func (d *destructKeyedSource) getIdx(p valueInt, receiver Value) Value ***REMOVED***
	d.recordKey(p.toString())
	return d.w().getIdx(p, receiver)
***REMOVED***

func (d *destructKeyedSource) getSym(p *Symbol, receiver Value) Value ***REMOVED***
	d.recordKey(p)
	return d.w().getSym(p, receiver)
***REMOVED***

func (d *destructKeyedSource) getOwnPropStr(u unistring.String) Value ***REMOVED***
	d.recordKey(stringValueFromRaw(u))
	return d.w().getOwnPropStr(u)
***REMOVED***

func (d *destructKeyedSource) getOwnPropIdx(v valueInt) Value ***REMOVED***
	d.recordKey(v.toString())
	return d.w().getOwnPropIdx(v)
***REMOVED***

func (d *destructKeyedSource) getOwnPropSym(symbol *Symbol) Value ***REMOVED***
	d.recordKey(symbol)
	return d.w().getOwnPropSym(symbol)
***REMOVED***

func (d *destructKeyedSource) setOwnStr(p unistring.String, v Value, throw bool) bool ***REMOVED***
	return d.w().setOwnStr(p, v, throw)
***REMOVED***

func (d *destructKeyedSource) setOwnIdx(p valueInt, v Value, throw bool) bool ***REMOVED***
	return d.w().setOwnIdx(p, v, throw)
***REMOVED***

func (d *destructKeyedSource) setOwnSym(p *Symbol, v Value, throw bool) bool ***REMOVED***
	return d.w().setOwnSym(p, v, throw)
***REMOVED***

func (d *destructKeyedSource) setForeignStr(p unistring.String, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	return d.w().setForeignStr(p, v, receiver, throw)
***REMOVED***

func (d *destructKeyedSource) setForeignIdx(p valueInt, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	return d.w().setForeignIdx(p, v, receiver, throw)
***REMOVED***

func (d *destructKeyedSource) setForeignSym(p *Symbol, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	return d.w().setForeignSym(p, v, receiver, throw)
***REMOVED***

func (d *destructKeyedSource) hasPropertyStr(u unistring.String) bool ***REMOVED***
	return d.w().hasPropertyStr(u)
***REMOVED***

func (d *destructKeyedSource) hasPropertyIdx(idx valueInt) bool ***REMOVED***
	return d.w().hasPropertyIdx(idx)
***REMOVED***

func (d *destructKeyedSource) hasPropertySym(s *Symbol) bool ***REMOVED***
	return d.w().hasPropertySym(s)
***REMOVED***

func (d *destructKeyedSource) hasOwnPropertyStr(u unistring.String) bool ***REMOVED***
	return d.w().hasOwnPropertyStr(u)
***REMOVED***

func (d *destructKeyedSource) hasOwnPropertyIdx(v valueInt) bool ***REMOVED***
	return d.w().hasOwnPropertyIdx(v)
***REMOVED***

func (d *destructKeyedSource) hasOwnPropertySym(s *Symbol) bool ***REMOVED***
	return d.w().hasOwnPropertySym(s)
***REMOVED***

func (d *destructKeyedSource) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	return d.w().defineOwnPropertyStr(name, desc, throw)
***REMOVED***

func (d *destructKeyedSource) defineOwnPropertyIdx(name valueInt, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	return d.w().defineOwnPropertyIdx(name, desc, throw)
***REMOVED***

func (d *destructKeyedSource) defineOwnPropertySym(name *Symbol, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	return d.w().defineOwnPropertySym(name, desc, throw)
***REMOVED***

func (d *destructKeyedSource) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	return d.w().deleteStr(name, throw)
***REMOVED***

func (d *destructKeyedSource) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	return d.w().deleteIdx(idx, throw)
***REMOVED***

func (d *destructKeyedSource) deleteSym(s *Symbol, throw bool) bool ***REMOVED***
	return d.w().deleteSym(s, throw)
***REMOVED***

func (d *destructKeyedSource) toPrimitiveNumber() Value ***REMOVED***
	return d.w().toPrimitiveNumber()
***REMOVED***

func (d *destructKeyedSource) toPrimitiveString() Value ***REMOVED***
	return d.w().toPrimitiveString()
***REMOVED***

func (d *destructKeyedSource) toPrimitive() Value ***REMOVED***
	return d.w().toPrimitive()
***REMOVED***

func (d *destructKeyedSource) assertCallable() (call func(FunctionCall) Value, ok bool) ***REMOVED***
	return d.w().assertCallable()
***REMOVED***

func (d *destructKeyedSource) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	return d.w().assertConstructor()
***REMOVED***

func (d *destructKeyedSource) proto() *Object ***REMOVED***
	return d.w().proto()
***REMOVED***

func (d *destructKeyedSource) setProto(proto *Object, throw bool) bool ***REMOVED***
	return d.w().setProto(proto, throw)
***REMOVED***

func (d *destructKeyedSource) hasInstance(v Value) bool ***REMOVED***
	return d.w().hasInstance(v)
***REMOVED***

func (d *destructKeyedSource) isExtensible() bool ***REMOVED***
	return d.w().isExtensible()
***REMOVED***

func (d *destructKeyedSource) preventExtensions(throw bool) bool ***REMOVED***
	return d.w().preventExtensions(throw)
***REMOVED***

type destructKeyedSourceIter struct ***REMOVED***
	d       *destructKeyedSource
	wrapped iterNextFunc
***REMOVED***

func (i *destructKeyedSourceIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for ***REMOVED***
		item, next := i.wrapped()
		if next == nil ***REMOVED***
			return item, nil
		***REMOVED***
		i.wrapped = next
		if _, exists := i.d.usedKeys[item.name]; !exists ***REMOVED***
			return item, i.next
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *destructKeyedSource) iterateStringKeys() iterNextFunc ***REMOVED***
	return (&destructKeyedSourceIter***REMOVED***
		d:       d,
		wrapped: d.w().iterateStringKeys(),
	***REMOVED***).next
***REMOVED***

func (d *destructKeyedSource) iterateSymbols() iterNextFunc ***REMOVED***
	return (&destructKeyedSourceIter***REMOVED***
		d:       d,
		wrapped: d.w().iterateSymbols(),
	***REMOVED***).next
***REMOVED***

func (d *destructKeyedSource) iterateKeys() iterNextFunc ***REMOVED***
	return (&destructKeyedSourceIter***REMOVED***
		d:       d,
		wrapped: d.w().iterateKeys(),
	***REMOVED***).next
***REMOVED***

func (d *destructKeyedSource) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return d.w().export(ctx)
***REMOVED***

func (d *destructKeyedSource) exportType() reflect.Type ***REMOVED***
	return d.w().exportType()
***REMOVED***

func (d *destructKeyedSource) exportToMap(dst reflect.Value, typ reflect.Type, ctx *objectExportCtx) error ***REMOVED***
	return d.w().exportToMap(dst, typ, ctx)
***REMOVED***

func (d *destructKeyedSource) exportToArrayOrSlice(dst reflect.Value, typ reflect.Type, ctx *objectExportCtx) error ***REMOVED***
	return d.w().exportToArrayOrSlice(dst, typ, ctx)
***REMOVED***

func (d *destructKeyedSource) equal(impl objectImpl) bool ***REMOVED***
	return d.w().equal(impl)
***REMOVED***

func (d *destructKeyedSource) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	var next iterNextFunc
	if all ***REMOVED***
		next = d.iterateStringKeys()
	***REMOVED*** else ***REMOVED***
		next = (&enumerableIter***REMOVED***
			o:       d.wrapped.ToObject(d.r),
			wrapped: d.iterateStringKeys(),
		***REMOVED***).next
	***REMOVED***
	for item, next := next(); next != nil; item, next = next() ***REMOVED***
		accum = append(accum, item.name)
	***REMOVED***
	return accum
***REMOVED***

func (d *destructKeyedSource) filterUsedKeys(keys []Value) []Value ***REMOVED***
	k := 0
	for i, key := range keys ***REMOVED***
		if _, exists := d.usedKeys[key]; exists ***REMOVED***
			continue
		***REMOVED***
		if k != i ***REMOVED***
			keys[k] = key
		***REMOVED***
		k++
	***REMOVED***
	return keys[:k]
***REMOVED***

func (d *destructKeyedSource) symbols(all bool, accum []Value) []Value ***REMOVED***
	return d.filterUsedKeys(d.w().symbols(all, accum))
***REMOVED***

func (d *destructKeyedSource) keys(all bool, accum []Value) []Value ***REMOVED***
	return d.filterUsedKeys(d.w().keys(all, accum))
***REMOVED***

func (d *destructKeyedSource) _putProp(name unistring.String, value Value, writable, enumerable, configurable bool) Value ***REMOVED***
	return d.w()._putProp(name, value, writable, enumerable, configurable)
***REMOVED***

func (d *destructKeyedSource) _putSym(s *Symbol, prop Value) ***REMOVED***
	d.w()._putSym(s, prop)
***REMOVED***
