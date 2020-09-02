package goja

import (
	"reflect"

	"github.com/dop251/goja/unistring"
)

// Proxy is a Go wrapper around ECMAScript Proxy. Calling Runtime.ToValue() on it
// returns the underlying Proxy. Calling Export() on an ECMAScript Proxy returns a wrapper.
// Use Runtime.NewProxy() to create one.
type Proxy struct ***REMOVED***
	proxy *proxyObject
***REMOVED***

var (
	proxyType = reflect.TypeOf(Proxy***REMOVED******REMOVED***)
)

type proxyPropIter struct ***REMOVED***
	p     *proxyObject
	names []Value
	idx   int
***REMOVED***

func (i *proxyPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.names) ***REMOVED***
		name := i.names[i.idx]
		i.idx++
		if prop := i.p.val.getOwnProp(name); prop != nil ***REMOVED***
			return propIterItem***REMOVED***name: name.string(), value: prop***REMOVED***, i.next
		***REMOVED***
	***REMOVED***
	if proto := i.p.proto(); proto != nil ***REMOVED***
		return proto.self.enumerateUnfiltered()()
	***REMOVED***
	return propIterItem***REMOVED******REMOVED***, nil
***REMOVED***

func (r *Runtime) newProxyObject(target, handler, proto *Object) *proxyObject ***REMOVED***
	if p, ok := target.self.(*proxyObject); ok ***REMOVED***
		if p.handler == nil ***REMOVED***
			panic(r.NewTypeError("Cannot create proxy with a revoked proxy as target"))
		***REMOVED***
	***REMOVED***
	if p, ok := handler.self.(*proxyObject); ok ***REMOVED***
		if p.handler == nil ***REMOVED***
			panic(r.NewTypeError("Cannot create proxy with a revoked proxy as handler"))
		***REMOVED***
	***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***
	p := &proxyObject***REMOVED******REMOVED***
	v.self = p
	p.val = v
	p.class = classObject
	if proto == nil ***REMOVED***
		p.prototype = r.global.ObjectPrototype
	***REMOVED*** else ***REMOVED***
		p.prototype = proto
	***REMOVED***
	p.extensible = false
	p.init()
	p.target = target
	p.handler = handler
	if call, ok := target.self.assertCallable(); ok ***REMOVED***
		p.call = call
	***REMOVED***
	if ctor := target.self.assertConstructor(); ctor != nil ***REMOVED***
		p.ctor = ctor
	***REMOVED***
	return p
***REMOVED***

func (p Proxy) Revoke() ***REMOVED***
	p.proxy.revoke()
***REMOVED***

func (p Proxy) toValue(r *Runtime) Value ***REMOVED***
	if p.proxy == nil ***REMOVED***
		return _null
	***REMOVED***
	proxy := p.proxy.val
	if proxy.runtime != r ***REMOVED***
		panic(r.NewTypeError("Illegal runtime transition of a Proxy"))
	***REMOVED***
	return proxy
***REMOVED***

type proxyTrap string

const (
	proxy_trap_getPrototypeOf           = "getPrototypeOf"
	proxy_trap_setPrototypeOf           = "setPrototypeOf"
	proxy_trap_isExtensible             = "isExtensible"
	proxy_trap_preventExtensions        = "preventExtensions"
	proxy_trap_getOwnPropertyDescriptor = "getOwnPropertyDescriptor"
	proxy_trap_defineProperty           = "defineProperty"
	proxy_trap_has                      = "has"
	proxy_trap_get                      = "get"
	proxy_trap_set                      = "set"
	proxy_trap_deleteProperty           = "deleteProperty"
	proxy_trap_ownKeys                  = "ownKeys"
	proxy_trap_apply                    = "apply"
	proxy_trap_construct                = "construct"
)

func (p proxyTrap) String() (name string) ***REMOVED***
	return string(p)
***REMOVED***

type proxyObject struct ***REMOVED***
	baseObject
	target  *Object
	handler *Object
	call    func(FunctionCall) Value
	ctor    func(args []Value, newTarget *Object) *Object
***REMOVED***

func (p *proxyObject) proxyCall(trap proxyTrap, args ...Value) (Value, bool) ***REMOVED***
	r := p.val.runtime
	if p.handler == nil ***REMOVED***
		panic(r.NewTypeError("Proxy already revoked"))
	***REMOVED***

	if m := toMethod(r.getVStr(p.handler, unistring.String(trap.String()))); m != nil ***REMOVED***
		return m(FunctionCall***REMOVED***
			This:      p.handler,
			Arguments: args,
		***REMOVED***), true
	***REMOVED***

	return nil, false
***REMOVED***

func (p *proxyObject) proto() *Object ***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_getPrototypeOf, p.target); ok ***REMOVED***
		var handlerProto *Object
		if v != _null ***REMOVED***
			handlerProto = p.val.runtime.toObject(v)
		***REMOVED***
		if !p.target.self.isExtensible() && !p.__sameValue(handlerProto, p.target.self.proto()) ***REMOVED***
			panic(p.val.runtime.NewTypeError("'getPrototypeOf' on proxy: proxy target is non-extensible but the trap did not return its actual prototype"))
		***REMOVED***
		return handlerProto
	***REMOVED***

	return p.target.self.proto()
***REMOVED***

func (p *proxyObject) setProto(proto *Object, throw bool) bool ***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_setPrototypeOf, p.target, proto); ok ***REMOVED***
		if v.ToBoolean() ***REMOVED***
			if !p.target.self.isExtensible() && !p.__sameValue(proto, p.target.self.proto()) ***REMOVED***
				panic(p.val.runtime.NewTypeError("'setPrototypeOf' on proxy: trap returned truish for setting a new prototype on the non-extensible proxy target"))
			***REMOVED***
			return true
		***REMOVED*** else ***REMOVED***
			p.val.runtime.typeErrorResult(throw, "'setPrototypeOf' on proxy: trap returned falsish")
		***REMOVED***
	***REMOVED***

	return p.target.self.setProto(proto, throw)
***REMOVED***

func (p *proxyObject) isExtensible() bool ***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_isExtensible, p.target); ok ***REMOVED***
		booleanTrapResult := v.ToBoolean()
		if te := p.target.self.isExtensible(); booleanTrapResult != te ***REMOVED***
			panic(p.val.runtime.NewTypeError("'isExtensible' on proxy: trap result does not reflect extensibility of proxy target (which is '%v')", te))
		***REMOVED***
		return booleanTrapResult
	***REMOVED***

	return p.target.self.isExtensible()
***REMOVED***

func (p *proxyObject) preventExtensions(throw bool) bool ***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_preventExtensions, p.target); ok ***REMOVED***
		booleanTrapResult := v.ToBoolean()
		if !booleanTrapResult ***REMOVED***
			p.val.runtime.typeErrorResult(throw, "'preventExtensions' on proxy: trap returned falsish")
			return false
		***REMOVED***
		if te := p.target.self.isExtensible(); booleanTrapResult && te ***REMOVED***
			panic(p.val.runtime.NewTypeError("'preventExtensions' on proxy: trap returned truish but the proxy target is extensible"))
		***REMOVED***
	***REMOVED***

	return p.target.self.preventExtensions(throw)
***REMOVED***

func propToValueProp(v Value) *valueProperty ***REMOVED***
	if v == nil ***REMOVED***
		return nil
	***REMOVED***
	if v, ok := v.(*valueProperty); ok ***REMOVED***
		return v
	***REMOVED***
	return &valueProperty***REMOVED***
		value:        v,
		writable:     true,
		configurable: true,
		enumerable:   true,
	***REMOVED***
***REMOVED***

func (p *proxyObject) proxyDefineOwnProperty(name Value, descr PropertyDescriptor, throw bool) (bool, bool) ***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_defineProperty, p.target, name, descr.toValue(p.val.runtime)); ok ***REMOVED***
		booleanTrapResult := v.ToBoolean()
		if !booleanTrapResult ***REMOVED***
			p.val.runtime.typeErrorResult(throw, "'defineProperty' on proxy: trap returned falsish")
			return false, true
		***REMOVED***
		targetDesc := propToValueProp(p.target.getOwnProp(name))
		extensibleTarget := p.target.self.isExtensible()
		settingConfigFalse := descr.Configurable == FLAG_FALSE
		if targetDesc == nil ***REMOVED***
			if !extensibleTarget ***REMOVED***
				panic(p.val.runtime.NewTypeError())
			***REMOVED***
			if settingConfigFalse ***REMOVED***
				panic(p.val.runtime.NewTypeError())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if !p.__isCompatibleDescriptor(extensibleTarget, &descr, targetDesc) ***REMOVED***
				panic(p.val.runtime.NewTypeError())
			***REMOVED***
			if settingConfigFalse && targetDesc.configurable ***REMOVED***
				panic(p.val.runtime.NewTypeError())
			***REMOVED***
		***REMOVED***
		return booleanTrapResult, true
	***REMOVED***
	return false, false
***REMOVED***

func (p *proxyObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if v, ok := p.proxyDefineOwnProperty(stringValueFromRaw(name), descr, throw); ok ***REMOVED***
		return v
	***REMOVED***
	return p.target.self.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (p *proxyObject) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if v, ok := p.proxyDefineOwnProperty(idx, descr, throw); ok ***REMOVED***
		return v
	***REMOVED***
	return p.target.self.defineOwnPropertyIdx(idx, descr, throw)
***REMOVED***

func (p *proxyObject) defineOwnPropertySym(s *valueSymbol, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if v, ok := p.proxyDefineOwnProperty(s, descr, throw); ok ***REMOVED***
		return v
	***REMOVED***
	return p.target.self.defineOwnPropertySym(s, descr, throw)
***REMOVED***

func (p *proxyObject) proxyHas(name Value) (bool, bool) ***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_has, p.target, name); ok ***REMOVED***
		booleanTrapResult := v.ToBoolean()
		if !booleanTrapResult ***REMOVED***
			targetDesc := propToValueProp(p.target.getOwnProp(name))
			if targetDesc != nil ***REMOVED***
				if !targetDesc.configurable ***REMOVED***
					panic(p.val.runtime.NewTypeError("'has' on proxy: trap returned falsish for property '%s' which exists in the proxy target as non-configurable", name.String()))
				***REMOVED***
				if !p.target.self.isExtensible() ***REMOVED***
					panic(p.val.runtime.NewTypeError("'has' on proxy: trap returned falsish for property '%s' but the proxy target is not extensible", name.String()))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return booleanTrapResult, true
	***REMOVED***

	return false, false
***REMOVED***

func (p *proxyObject) hasPropertyStr(name unistring.String) bool ***REMOVED***
	if b, ok := p.proxyHas(stringValueFromRaw(name)); ok ***REMOVED***
		return b
	***REMOVED***

	return p.target.self.hasPropertyStr(name)
***REMOVED***

func (p *proxyObject) hasPropertyIdx(idx valueInt) bool ***REMOVED***
	if b, ok := p.proxyHas(idx); ok ***REMOVED***
		return b
	***REMOVED***

	return p.target.self.hasPropertyIdx(idx)
***REMOVED***

func (p *proxyObject) hasPropertySym(s *valueSymbol) bool ***REMOVED***
	if b, ok := p.proxyHas(s); ok ***REMOVED***
		return b
	***REMOVED***

	return p.target.self.hasPropertySym(s)
***REMOVED***

func (p *proxyObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	return p.getOwnPropStr(name) != nil
***REMOVED***

func (p *proxyObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	return p.getOwnPropIdx(idx) != nil
***REMOVED***

func (p *proxyObject) hasOwnPropertySym(s *valueSymbol) bool ***REMOVED***
	return p.getOwnPropSym(s) != nil
***REMOVED***

func (p *proxyObject) proxyGetOwnPropertyDescriptor(name Value) (Value, bool) ***REMOVED***
	target := p.target
	if v, ok := p.proxyCall(proxy_trap_getOwnPropertyDescriptor, target, name); ok ***REMOVED***
		r := p.val.runtime

		targetDesc := propToValueProp(target.getOwnProp(name))

		var trapResultObj *Object
		if v != nil && v != _undefined ***REMOVED***
			if obj, ok := v.(*Object); ok ***REMOVED***
				trapResultObj = obj
			***REMOVED*** else ***REMOVED***
				panic(r.NewTypeError("'getOwnPropertyDescriptor' on proxy: trap returned neither object nor undefined for property '%s'", name.String()))
			***REMOVED***
		***REMOVED***
		if trapResultObj == nil ***REMOVED***
			if targetDesc == nil ***REMOVED***
				return nil, true
			***REMOVED***
			if !targetDesc.configurable ***REMOVED***
				panic(r.NewTypeError())
			***REMOVED***
			if !target.self.isExtensible() ***REMOVED***
				panic(r.NewTypeError())
			***REMOVED***
			return nil, true
		***REMOVED***
		extensibleTarget := target.self.isExtensible()
		resultDesc := r.toPropertyDescriptor(trapResultObj)
		resultDesc.complete()
		if !p.__isCompatibleDescriptor(extensibleTarget, &resultDesc, targetDesc) ***REMOVED***
			panic(r.NewTypeError("'getOwnPropertyDescriptor' on proxy: trap returned descriptor for property '%s' that is incompatible with the existing property in the proxy target", name.String()))
		***REMOVED***

		if resultDesc.Configurable == FLAG_FALSE ***REMOVED***
			if targetDesc == nil ***REMOVED***
				panic(r.NewTypeError("'getOwnPropertyDescriptor' on proxy: trap reported non-configurability for property '%s' which is non-existent in the proxy target", name.String()))
			***REMOVED***

			if targetDesc.configurable ***REMOVED***
				panic(r.NewTypeError("'getOwnPropertyDescriptor' on proxy: trap reported non-configurability for property '%s' which is configurable in the proxy target", name.String()))
			***REMOVED***
		***REMOVED***

		if resultDesc.Writable == FLAG_TRUE && resultDesc.Configurable == FLAG_TRUE &&
			resultDesc.Enumerable == FLAG_TRUE ***REMOVED***
			return resultDesc.Value, true
		***REMOVED***
		return r.toValueProp(trapResultObj), true
	***REMOVED***

	return nil, false
***REMOVED***

func (p *proxyObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if v, ok := p.proxyGetOwnPropertyDescriptor(stringValueFromRaw(name)); ok ***REMOVED***
		return v
	***REMOVED***

	return p.target.self.getOwnPropStr(name)
***REMOVED***

func (p *proxyObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	if v, ok := p.proxyGetOwnPropertyDescriptor(idx.toString()); ok ***REMOVED***
		return v
	***REMOVED***

	return p.target.self.getOwnPropIdx(idx)
***REMOVED***

func (p *proxyObject) getOwnPropSym(s *valueSymbol) Value ***REMOVED***
	if v, ok := p.proxyGetOwnPropertyDescriptor(s); ok ***REMOVED***
		return v
	***REMOVED***

	return p.target.self.getOwnPropSym(s)
***REMOVED***

func (p *proxyObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	if v, ok := p.proxyGet(stringValueFromRaw(name), receiver); ok ***REMOVED***
		return v
	***REMOVED***
	return p.target.self.getStr(name, receiver)
***REMOVED***

func (p *proxyObject) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	if v, ok := p.proxyGet(idx.toString(), receiver); ok ***REMOVED***
		return v
	***REMOVED***
	return p.target.self.getIdx(idx, receiver)
***REMOVED***

func (p *proxyObject) getSym(s *valueSymbol, receiver Value) Value ***REMOVED***
	if v, ok := p.proxyGet(s, receiver); ok ***REMOVED***
		return v
	***REMOVED***
	return p.target.self.getSym(s, receiver)

***REMOVED***

func (p *proxyObject) proxyGet(name, receiver Value) (Value, bool) ***REMOVED***
	target := p.target
	if receiver == nil ***REMOVED***
		receiver = p.val
	***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_get, target, name, receiver); ok ***REMOVED***
		if targetDesc, ok := target.getOwnProp(name).(*valueProperty); ok ***REMOVED***
			if !targetDesc.accessor ***REMOVED***
				if !targetDesc.writable && !targetDesc.configurable && !v.SameAs(targetDesc.value) ***REMOVED***
					panic(p.val.runtime.NewTypeError("'get' on proxy: property '%s' is a read-only and non-configurable data property on the proxy target but the proxy did not return its actual value (expected '%s' but got '%s')", name.String(), nilSafe(targetDesc.value), ret))
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if !targetDesc.configurable && targetDesc.getterFunc == nil && v != _undefined ***REMOVED***
					panic(p.val.runtime.NewTypeError("'get' on proxy: property '%s' is a non-configurable accessor property on the proxy target and does not have a getter function, but the trap did not return 'undefined' (got '%s')", name.String(), ret))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return v, true
	***REMOVED***

	return nil, false
***REMOVED***

func (p *proxyObject) proxySet(name, value, receiver Value, throw bool) (bool, bool) ***REMOVED***
	target := p.target
	if v, ok := p.proxyCall(proxy_trap_set, target, name, value, receiver); ok ***REMOVED***
		if v.ToBoolean() ***REMOVED***
			if prop, ok := target.getOwnProp(name).(*valueProperty); ok ***REMOVED***
				if prop.accessor ***REMOVED***
					if !prop.configurable && prop.setterFunc == nil ***REMOVED***
						panic(p.val.runtime.NewTypeError("'set' on proxy: trap returned truish for property '%s' which exists in the proxy target as a non-configurable and non-writable accessor property without a setter", name.String()))
					***REMOVED***
				***REMOVED*** else if !prop.configurable && !prop.writable && !p.__sameValue(prop.value, value) ***REMOVED***
					panic(p.val.runtime.NewTypeError("'set' on proxy: trap returned truish for property '%s' which exists in the proxy target as a non-configurable and non-writable data property with a different value", name.String()))
				***REMOVED***
			***REMOVED***
			return true, true
		***REMOVED***
		if throw ***REMOVED***
			panic(p.val.runtime.NewTypeError("'set' on proxy: trap returned falsish for property '%s'", name.String()))
		***REMOVED***
		return false, true
	***REMOVED***

	return false, false
***REMOVED***

func (p *proxyObject) setOwnStr(name unistring.String, v Value, throw bool) bool ***REMOVED***
	if res, ok := p.proxySet(stringValueFromRaw(name), v, p.val, throw); ok ***REMOVED***
		return res
	***REMOVED***
	return p.target.setStr(name, v, p.val, throw)
***REMOVED***

func (p *proxyObject) setOwnIdx(idx valueInt, v Value, throw bool) bool ***REMOVED***
	if res, ok := p.proxySet(idx.toString(), v, p.val, throw); ok ***REMOVED***
		return res
	***REMOVED***
	return p.target.setIdx(idx, v, p.val, throw)
***REMOVED***

func (p *proxyObject) setOwnSym(s *valueSymbol, v Value, throw bool) bool ***REMOVED***
	if res, ok := p.proxySet(s, v, p.val, throw); ok ***REMOVED***
		return res
	***REMOVED***
	return p.target.setSym(s, v, p.val, throw)
***REMOVED***

func (p *proxyObject) setForeignStr(name unistring.String, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	if res, ok := p.proxySet(stringValueFromRaw(name), v, receiver, throw); ok ***REMOVED***
		return res, true
	***REMOVED***
	return p.target.setStr(name, v, receiver, throw), true
***REMOVED***

func (p *proxyObject) setForeignIdx(idx valueInt, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	if res, ok := p.proxySet(idx.toString(), v, receiver, throw); ok ***REMOVED***
		return res, true
	***REMOVED***
	return p.target.setIdx(idx, v, receiver, throw), true
***REMOVED***

func (p *proxyObject) setForeignSym(s *valueSymbol, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	if res, ok := p.proxySet(s, v, receiver, throw); ok ***REMOVED***
		return res, true
	***REMOVED***
	return p.target.setSym(s, v, receiver, throw), true
***REMOVED***

func (p *proxyObject) proxyDelete(n Value) (bool, bool) ***REMOVED***
	target := p.target
	if v, ok := p.proxyCall(proxy_trap_deleteProperty, target, n); ok ***REMOVED***
		if v.ToBoolean() ***REMOVED***
			if targetDesc, ok := target.getOwnProp(n).(*valueProperty); ok ***REMOVED***
				if !targetDesc.configurable ***REMOVED***
					panic(p.val.runtime.NewTypeError("'deleteProperty' on proxy: property '%s' is a non-configurable property but the trap returned truish", n.String()))
				***REMOVED***
			***REMOVED***
			return true, true
		***REMOVED***
		return false, true
	***REMOVED***
	return false, false
***REMOVED***

func (p *proxyObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if ret, ok := p.proxyDelete(stringValueFromRaw(name)); ok ***REMOVED***
		return ret
	***REMOVED***

	return p.target.self.deleteStr(name, throw)
***REMOVED***

func (p *proxyObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	if ret, ok := p.proxyDelete(idx.toString()); ok ***REMOVED***
		return ret
	***REMOVED***

	return p.target.self.deleteIdx(idx, throw)
***REMOVED***

func (p *proxyObject) deleteSym(s *valueSymbol, throw bool) bool ***REMOVED***
	if ret, ok := p.proxyDelete(s); ok ***REMOVED***
		return ret
	***REMOVED***

	return p.target.self.deleteSym(s, throw)
***REMOVED***

func (p *proxyObject) ownPropertyKeys(all bool, _ []Value) []Value ***REMOVED***
	if v, ok := p.proxyOwnKeys(); ok ***REMOVED***
		return v
	***REMOVED***
	return p.target.self.ownPropertyKeys(all, nil)
***REMOVED***

func (p *proxyObject) proxyOwnKeys() ([]Value, bool) ***REMOVED***
	target := p.target
	if v, ok := p.proxyCall(proxy_trap_ownKeys, p.target); ok ***REMOVED***
		keys := p.val.runtime.toObject(v)
		var keyList []Value
		keySet := make(map[Value]struct***REMOVED******REMOVED***)
		l := toLength(keys.self.getStr("length", nil))
		for k := int64(0); k < l; k++ ***REMOVED***
			item := keys.self.getIdx(valueInt(k), nil)
			if _, ok := item.(valueString); !ok ***REMOVED***
				if _, ok := item.(*valueSymbol); !ok ***REMOVED***
					panic(p.val.runtime.NewTypeError("%s is not a valid property name", item.String()))
				***REMOVED***
			***REMOVED***
			keyList = append(keyList, item)
			keySet[item] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
		ext := target.self.isExtensible()
		for _, itemName := range target.self.ownPropertyKeys(true, nil) ***REMOVED***
			if _, exists := keySet[itemName]; exists ***REMOVED***
				delete(keySet, itemName)
			***REMOVED*** else ***REMOVED***
				if !ext ***REMOVED***
					panic(p.val.runtime.NewTypeError("'ownKeys' on proxy: trap result did not include '%s'", itemName.String()))
				***REMOVED***
				prop := target.getOwnProp(itemName)
				if prop, ok := prop.(*valueProperty); ok && !prop.configurable ***REMOVED***
					panic(p.val.runtime.NewTypeError("'ownKeys' on proxy: trap result did not include non-configurable '%s'", itemName.String()))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if !ext && len(keyList) > 0 && len(keySet) > 0 ***REMOVED***
			panic(p.val.runtime.NewTypeError("'ownKeys' on proxy: trap returned extra keys but proxy target is non-extensible"))
		***REMOVED***

		return keyList, true
	***REMOVED***

	return nil, false
***REMOVED***

func (p *proxyObject) enumerateUnfiltered() iterNextFunc ***REMOVED***
	return (&proxyPropIter***REMOVED***
		p:     p,
		names: p.ownKeys(true, nil),
	***REMOVED***).next
***REMOVED***

func (p *proxyObject) assertCallable() (call func(FunctionCall) Value, ok bool) ***REMOVED***
	if p.call != nil ***REMOVED***
		return func(call FunctionCall) Value ***REMOVED***
			return p.apply(call)
		***REMOVED***, true
	***REMOVED***
	return nil, false
***REMOVED***

func (p *proxyObject) assertConstructor() func(args []Value, newTarget *Object) *Object ***REMOVED***
	if p.ctor != nil ***REMOVED***
		return p.construct
	***REMOVED***
	return nil
***REMOVED***

func (p *proxyObject) apply(call FunctionCall) Value ***REMOVED***
	if p.call == nil ***REMOVED***
		p.val.runtime.NewTypeError("proxy target is not a function")
	***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_apply, p.target, nilSafe(call.This), p.val.runtime.newArrayValues(call.Arguments)); ok ***REMOVED***
		return v
	***REMOVED***
	return p.call(call)
***REMOVED***

func (p *proxyObject) construct(args []Value, newTarget *Object) *Object ***REMOVED***
	if p.ctor == nil ***REMOVED***
		panic(p.val.runtime.NewTypeError("proxy target is not a constructor"))
	***REMOVED***
	if newTarget == nil ***REMOVED***
		newTarget = p.val
	***REMOVED***
	if v, ok := p.proxyCall(proxy_trap_construct, p.target, p.val.runtime.newArrayValues(args), newTarget); ok ***REMOVED***
		return p.val.runtime.toObject(v)
	***REMOVED***
	return p.ctor(args, newTarget)
***REMOVED***

func (p *proxyObject) __isCompatibleDescriptor(extensible bool, desc *PropertyDescriptor, current *valueProperty) bool ***REMOVED***
	if current == nil ***REMOVED***
		return extensible
	***REMOVED***

	/*if desc.Empty() ***REMOVED***
		return true
	***REMOVED****/

	/*if p.__isEquivalentDescriptor(desc, current) ***REMOVED***
		return true
	***REMOVED****/

	if !current.configurable ***REMOVED***
		if desc.Configurable == FLAG_TRUE ***REMOVED***
			return false
		***REMOVED***

		if desc.Enumerable != FLAG_NOT_SET && desc.Enumerable.Bool() != current.enumerable ***REMOVED***
			return false
		***REMOVED***

		if p.__isGenericDescriptor(desc) ***REMOVED***
			return true
		***REMOVED***

		if p.__isDataDescriptor(desc) != !current.accessor ***REMOVED***
			return desc.Configurable != FLAG_FALSE
		***REMOVED***

		if p.__isDataDescriptor(desc) && !current.accessor ***REMOVED***
			if !current.configurable ***REMOVED***
				if desc.Writable == FLAG_TRUE && !current.writable ***REMOVED***
					return false
				***REMOVED***
				if !current.writable ***REMOVED***
					if desc.Value != nil && !desc.Value.SameAs(current.value) ***REMOVED***
						return false
					***REMOVED***
				***REMOVED***
			***REMOVED***
			return true
		***REMOVED***
		if p.__isAccessorDescriptor(desc) && current.accessor ***REMOVED***
			if !current.configurable ***REMOVED***
				if desc.Setter != nil && desc.Setter.SameAs(current.setterFunc) ***REMOVED***
					return false
				***REMOVED***
				if desc.Getter != nil && desc.Getter.SameAs(current.getterFunc) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (p *proxyObject) __isAccessorDescriptor(desc *PropertyDescriptor) bool ***REMOVED***
	return desc.Setter != nil || desc.Getter != nil
***REMOVED***

func (p *proxyObject) __isDataDescriptor(desc *PropertyDescriptor) bool ***REMOVED***
	return desc.Value != nil || desc.Writable != FLAG_NOT_SET
***REMOVED***

func (p *proxyObject) __isGenericDescriptor(desc *PropertyDescriptor) bool ***REMOVED***
	return !p.__isAccessorDescriptor(desc) && !p.__isDataDescriptor(desc)
***REMOVED***

func (p *proxyObject) __sameValue(val1, val2 Value) bool ***REMOVED***
	if val1 == nil && val2 == nil ***REMOVED***
		return true
	***REMOVED***
	if val1 != nil ***REMOVED***
		return val1.SameAs(val2)
	***REMOVED***
	return false
***REMOVED***

func (p *proxyObject) filterKeys(vals []Value, all, symbols bool) []Value ***REMOVED***
	if !all ***REMOVED***
		k := 0
		for i, val := range vals ***REMOVED***
			var prop Value
			if symbols ***REMOVED***
				if s, ok := val.(*valueSymbol); ok ***REMOVED***
					prop = p.getOwnPropSym(s)
				***REMOVED*** else ***REMOVED***
					continue
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if _, ok := val.(*valueSymbol); !ok ***REMOVED***
					prop = p.getOwnPropStr(val.string())
				***REMOVED*** else ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			if prop == nil ***REMOVED***
				continue
			***REMOVED***
			if prop, ok := prop.(*valueProperty); ok && !prop.enumerable ***REMOVED***
				continue
			***REMOVED***
			if k != i ***REMOVED***
				vals[k] = vals[i]
			***REMOVED***
			k++
		***REMOVED***
		vals = vals[:k]
	***REMOVED*** else ***REMOVED***
		k := 0
		for i, val := range vals ***REMOVED***
			if _, ok := val.(*valueSymbol); ok != symbols ***REMOVED***
				continue
			***REMOVED***
			if k != i ***REMOVED***
				vals[k] = vals[i]
			***REMOVED***
			k++
		***REMOVED***
		vals = vals[:k]
	***REMOVED***
	return vals
***REMOVED***

func (p *proxyObject) ownKeys(all bool, _ []Value) []Value ***REMOVED*** // we can assume accum is empty
	if vals, ok := p.proxyOwnKeys(); ok ***REMOVED***
		return p.filterKeys(vals, all, false)
	***REMOVED***

	return p.target.self.ownKeys(all, nil)
***REMOVED***

func (p *proxyObject) ownSymbols(all bool, accum []Value) []Value ***REMOVED***
	if vals, ok := p.proxyOwnKeys(); ok ***REMOVED***
		res := p.filterKeys(vals, true, true)
		if accum == nil ***REMOVED***
			return res
		***REMOVED***
		accum = append(accum, res...)
		return accum
	***REMOVED***

	return p.target.self.ownSymbols(all, accum)
***REMOVED***

func (p *proxyObject) className() string ***REMOVED***
	if p.target == nil ***REMOVED***
		panic(p.val.runtime.NewTypeError("proxy has been revoked"))
	***REMOVED***
	if p.call != nil || p.ctor != nil ***REMOVED***
		return classFunction
	***REMOVED***
	return classObject
***REMOVED***

func (p *proxyObject) exportType() reflect.Type ***REMOVED***
	return proxyType
***REMOVED***

func (p *proxyObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return Proxy***REMOVED***
		proxy: p,
	***REMOVED***
***REMOVED***

func (p *proxyObject) revoke() ***REMOVED***
	p.handler = nil
	p.target = nil
***REMOVED***
