package goja

import (
	"fmt"
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
	return r._newProxyObject(target, &jsProxyHandler***REMOVED***handler: handler***REMOVED***, proto)
***REMOVED***

func (r *Runtime) _newProxyObject(target *Object, handler proxyHandler, proto *Object) *proxyObject ***REMOVED***
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

func (p Proxy) Handler() *Object ***REMOVED***
	if handler := p.proxy.handler; handler != nil ***REMOVED***
		return handler.toObject(p.proxy.val.runtime)
	***REMOVED***
	return nil
***REMOVED***

func (p Proxy) Target() *Object ***REMOVED***
	return p.proxy.target
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

type proxyHandler interface ***REMOVED***
	getPrototypeOf(target *Object) (Value, bool)
	setPrototypeOf(target *Object, proto *Object) (bool, bool)
	isExtensible(target *Object) (bool, bool)
	preventExtensions(target *Object) (bool, bool)

	getOwnPropertyDescriptorStr(target *Object, prop unistring.String) (Value, bool)
	getOwnPropertyDescriptorIdx(target *Object, prop valueInt) (Value, bool)
	getOwnPropertyDescriptorSym(target *Object, prop *Symbol) (Value, bool)

	definePropertyStr(target *Object, prop unistring.String, desc PropertyDescriptor) (bool, bool)
	definePropertyIdx(target *Object, prop valueInt, desc PropertyDescriptor) (bool, bool)
	definePropertySym(target *Object, prop *Symbol, desc PropertyDescriptor) (bool, bool)

	hasStr(target *Object, prop unistring.String) (bool, bool)
	hasIdx(target *Object, prop valueInt) (bool, bool)
	hasSym(target *Object, prop *Symbol) (bool, bool)

	getStr(target *Object, prop unistring.String, receiver Value) (Value, bool)
	getIdx(target *Object, prop valueInt, receiver Value) (Value, bool)
	getSym(target *Object, prop *Symbol, receiver Value) (Value, bool)

	setStr(target *Object, prop unistring.String, value Value, receiver Value) (bool, bool)
	setIdx(target *Object, prop valueInt, value Value, receiver Value) (bool, bool)
	setSym(target *Object, prop *Symbol, value Value, receiver Value) (bool, bool)

	deleteStr(target *Object, prop unistring.String) (bool, bool)
	deleteIdx(target *Object, prop valueInt) (bool, bool)
	deleteSym(target *Object, prop *Symbol) (bool, bool)

	ownKeys(target *Object) (*Object, bool)
	apply(target *Object, this Value, args []Value) (Value, bool)
	construct(target *Object, args []Value, newTarget *Object) (Value, bool)

	toObject(*Runtime) *Object
***REMOVED***

type jsProxyHandler struct ***REMOVED***
	handler *Object
***REMOVED***

func (h *jsProxyHandler) toObject(*Runtime) *Object ***REMOVED***
	return h.handler
***REMOVED***

func (h *jsProxyHandler) proxyCall(trap proxyTrap, args ...Value) (Value, bool) ***REMOVED***
	r := h.handler.runtime

	if m := toMethod(r.getVStr(h.handler, unistring.String(trap.String()))); m != nil ***REMOVED***
		return m(FunctionCall***REMOVED***
			This:      h.handler,
			Arguments: args,
		***REMOVED***), true
	***REMOVED***

	return nil, false
***REMOVED***

func (h *jsProxyHandler) boolProxyCall(trap proxyTrap, args ...Value) (bool, bool) ***REMOVED***
	if v, ok := h.proxyCall(trap, args...); ok ***REMOVED***
		return v.ToBoolean(), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *jsProxyHandler) getPrototypeOf(target *Object) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_getPrototypeOf, target)
***REMOVED***

func (h *jsProxyHandler) setPrototypeOf(target *Object, proto *Object) (bool, bool) ***REMOVED***
	var protoVal Value
	if proto != nil ***REMOVED***
		protoVal = proto
	***REMOVED*** else ***REMOVED***
		protoVal = _null
	***REMOVED***
	return h.boolProxyCall(proxy_trap_setPrototypeOf, target, protoVal)
***REMOVED***

func (h *jsProxyHandler) isExtensible(target *Object) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_isExtensible, target)
***REMOVED***

func (h *jsProxyHandler) preventExtensions(target *Object) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_preventExtensions, target)
***REMOVED***

func (h *jsProxyHandler) getOwnPropertyDescriptorStr(target *Object, prop unistring.String) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_getOwnPropertyDescriptor, target, stringValueFromRaw(prop))
***REMOVED***

func (h *jsProxyHandler) getOwnPropertyDescriptorIdx(target *Object, prop valueInt) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_getOwnPropertyDescriptor, target, prop.toString())
***REMOVED***

func (h *jsProxyHandler) getOwnPropertyDescriptorSym(target *Object, prop *Symbol) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_getOwnPropertyDescriptor, target, prop)
***REMOVED***

func (h *jsProxyHandler) definePropertyStr(target *Object, prop unistring.String, desc PropertyDescriptor) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_defineProperty, target, stringValueFromRaw(prop), desc.toValue(h.handler.runtime))
***REMOVED***

func (h *jsProxyHandler) definePropertyIdx(target *Object, prop valueInt, desc PropertyDescriptor) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_defineProperty, target, prop.toString(), desc.toValue(h.handler.runtime))
***REMOVED***

func (h *jsProxyHandler) definePropertySym(target *Object, prop *Symbol, desc PropertyDescriptor) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_defineProperty, target, prop, desc.toValue(h.handler.runtime))
***REMOVED***

func (h *jsProxyHandler) hasStr(target *Object, prop unistring.String) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_has, target, stringValueFromRaw(prop))
***REMOVED***

func (h *jsProxyHandler) hasIdx(target *Object, prop valueInt) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_has, target, prop.toString())
***REMOVED***

func (h *jsProxyHandler) hasSym(target *Object, prop *Symbol) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_has, target, prop)
***REMOVED***

func (h *jsProxyHandler) getStr(target *Object, prop unistring.String, receiver Value) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_get, target, stringValueFromRaw(prop), receiver)
***REMOVED***

func (h *jsProxyHandler) getIdx(target *Object, prop valueInt, receiver Value) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_get, target, prop.toString(), receiver)
***REMOVED***

func (h *jsProxyHandler) getSym(target *Object, prop *Symbol, receiver Value) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_get, target, prop, receiver)
***REMOVED***

func (h *jsProxyHandler) setStr(target *Object, prop unistring.String, value Value, receiver Value) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_set, target, stringValueFromRaw(prop), value, receiver)
***REMOVED***

func (h *jsProxyHandler) setIdx(target *Object, prop valueInt, value Value, receiver Value) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_set, target, prop.toString(), value, receiver)
***REMOVED***

func (h *jsProxyHandler) setSym(target *Object, prop *Symbol, value Value, receiver Value) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_set, target, prop, value, receiver)
***REMOVED***

func (h *jsProxyHandler) deleteStr(target *Object, prop unistring.String) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_deleteProperty, target, stringValueFromRaw(prop))
***REMOVED***

func (h *jsProxyHandler) deleteIdx(target *Object, prop valueInt) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_deleteProperty, target, prop.toString())
***REMOVED***

func (h *jsProxyHandler) deleteSym(target *Object, prop *Symbol) (bool, bool) ***REMOVED***
	return h.boolProxyCall(proxy_trap_deleteProperty, target, prop)
***REMOVED***

func (h *jsProxyHandler) ownKeys(target *Object) (*Object, bool) ***REMOVED***
	if v, ok := h.proxyCall(proxy_trap_ownKeys, target); ok ***REMOVED***
		return h.handler.runtime.toObject(v), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *jsProxyHandler) apply(target *Object, this Value, args []Value) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_apply, target, this, h.handler.runtime.newArrayValues(args))
***REMOVED***

func (h *jsProxyHandler) construct(target *Object, args []Value, newTarget *Object) (Value, bool) ***REMOVED***
	return h.proxyCall(proxy_trap_construct, target, h.handler.runtime.newArrayValues(args), newTarget)
***REMOVED***

type proxyObject struct ***REMOVED***
	baseObject
	target  *Object
	handler proxyHandler
	call    func(FunctionCall) Value
	ctor    func(args []Value, newTarget *Object) *Object
***REMOVED***

func (p *proxyObject) checkHandler() proxyHandler ***REMOVED***
	r := p.val.runtime
	if handler := p.handler; handler != nil ***REMOVED***
		return handler
	***REMOVED***
	panic(r.NewTypeError("Proxy already revoked"))
***REMOVED***

func (p *proxyObject) proto() *Object ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().getPrototypeOf(target); ok ***REMOVED***
		var handlerProto *Object
		if v != _null ***REMOVED***
			handlerProto = p.val.runtime.toObject(v)
		***REMOVED***
		if !target.self.isExtensible() && !p.__sameValue(handlerProto, target.self.proto()) ***REMOVED***
			panic(p.val.runtime.NewTypeError("'getPrototypeOf' on proxy: proxy target is non-extensible but the trap did not return its actual prototype"))
		***REMOVED***
		return handlerProto
	***REMOVED***

	return target.self.proto()
***REMOVED***

func (p *proxyObject) setProto(proto *Object, throw bool) bool ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().setPrototypeOf(target, proto); ok ***REMOVED***
		if v ***REMOVED***
			if !target.self.isExtensible() && !p.__sameValue(proto, target.self.proto()) ***REMOVED***
				panic(p.val.runtime.NewTypeError("'setPrototypeOf' on proxy: trap returned truish for setting a new prototype on the non-extensible proxy target"))
			***REMOVED***
			return true
		***REMOVED*** else ***REMOVED***
			p.val.runtime.typeErrorResult(throw, "'setPrototypeOf' on proxy: trap returned falsish")
			return false
		***REMOVED***
	***REMOVED***

	return target.self.setProto(proto, throw)
***REMOVED***

func (p *proxyObject) isExtensible() bool ***REMOVED***
	target := p.target
	if booleanTrapResult, ok := p.checkHandler().isExtensible(p.target); ok ***REMOVED***
		if te := target.self.isExtensible(); booleanTrapResult != te ***REMOVED***
			panic(p.val.runtime.NewTypeError("'isExtensible' on proxy: trap result does not reflect extensibility of proxy target (which is '%v')", te))
		***REMOVED***
		return booleanTrapResult
	***REMOVED***

	return target.self.isExtensible()
***REMOVED***

func (p *proxyObject) preventExtensions(throw bool) bool ***REMOVED***
	target := p.target
	if booleanTrapResult, ok := p.checkHandler().preventExtensions(target); ok ***REMOVED***
		if !booleanTrapResult ***REMOVED***
			p.val.runtime.typeErrorResult(throw, "'preventExtensions' on proxy: trap returned falsish")
			return false
		***REMOVED***
		if te := target.self.isExtensible(); booleanTrapResult && te ***REMOVED***
			panic(p.val.runtime.NewTypeError("'preventExtensions' on proxy: trap returned truish but the proxy target is extensible"))
		***REMOVED***
	***REMOVED***

	return target.self.preventExtensions(throw)
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

func (p *proxyObject) proxyDefineOwnPropertyPreCheck(trapResult, throw bool) bool ***REMOVED***
	if !trapResult ***REMOVED***
		p.val.runtime.typeErrorResult(throw, "'defineProperty' on proxy: trap returned falsish")
		return false
	***REMOVED***
	return true
***REMOVED***

func (p *proxyObject) proxyDefineOwnPropertyPostCheck(prop Value, target *Object, descr PropertyDescriptor) ***REMOVED***
	targetDesc := propToValueProp(prop)
	extensibleTarget := target.self.isExtensible()
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
		if targetDesc.value != nil && !targetDesc.configurable && targetDesc.writable ***REMOVED***
			if descr.Writable == FLAG_FALSE ***REMOVED***
				panic(p.val.runtime.NewTypeError())
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *proxyObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	target := p.target
	if booleanTrapResult, ok := p.checkHandler().definePropertyStr(target, name, descr); ok ***REMOVED***
		if !p.proxyDefineOwnPropertyPreCheck(booleanTrapResult, throw) ***REMOVED***
			return false
		***REMOVED***
		p.proxyDefineOwnPropertyPostCheck(target.self.getOwnPropStr(name), target, descr)
		return true
	***REMOVED***
	return target.self.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (p *proxyObject) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	target := p.target
	if booleanTrapResult, ok := p.checkHandler().definePropertyIdx(target, idx, descr); ok ***REMOVED***
		if !p.proxyDefineOwnPropertyPreCheck(booleanTrapResult, throw) ***REMOVED***
			return false
		***REMOVED***
		p.proxyDefineOwnPropertyPostCheck(target.self.getOwnPropIdx(idx), target, descr)
		return true
	***REMOVED***

	return target.self.defineOwnPropertyIdx(idx, descr, throw)
***REMOVED***

func (p *proxyObject) defineOwnPropertySym(s *Symbol, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	target := p.target
	if booleanTrapResult, ok := p.checkHandler().definePropertySym(target, s, descr); ok ***REMOVED***
		if !p.proxyDefineOwnPropertyPreCheck(booleanTrapResult, throw) ***REMOVED***
			return false
		***REMOVED***
		p.proxyDefineOwnPropertyPostCheck(target.self.getOwnPropSym(s), target, descr)
		return true
	***REMOVED***

	return target.self.defineOwnPropertySym(s, descr, throw)
***REMOVED***

func (p *proxyObject) proxyHasChecks(targetProp Value, target *Object, name fmt.Stringer) ***REMOVED***
	targetDesc := propToValueProp(targetProp)
	if targetDesc != nil ***REMOVED***
		if !targetDesc.configurable ***REMOVED***
			panic(p.val.runtime.NewTypeError("'has' on proxy: trap returned falsish for property '%s' which exists in the proxy target as non-configurable", name.String()))
		***REMOVED***
		if !target.self.isExtensible() ***REMOVED***
			panic(p.val.runtime.NewTypeError("'has' on proxy: trap returned falsish for property '%s' but the proxy target is not extensible", name.String()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *proxyObject) hasPropertyStr(name unistring.String) bool ***REMOVED***
	target := p.target
	if b, ok := p.checkHandler().hasStr(target, name); ok ***REMOVED***
		if !b ***REMOVED***
			p.proxyHasChecks(target.self.getOwnPropStr(name), target, name)
		***REMOVED***
		return b
	***REMOVED***

	return target.self.hasPropertyStr(name)
***REMOVED***

func (p *proxyObject) hasPropertyIdx(idx valueInt) bool ***REMOVED***
	target := p.target
	if b, ok := p.checkHandler().hasIdx(target, idx); ok ***REMOVED***
		if !b ***REMOVED***
			p.proxyHasChecks(target.self.getOwnPropIdx(idx), target, idx)
		***REMOVED***
		return b
	***REMOVED***

	return target.self.hasPropertyIdx(idx)
***REMOVED***

func (p *proxyObject) hasPropertySym(s *Symbol) bool ***REMOVED***
	target := p.target
	if b, ok := p.checkHandler().hasSym(target, s); ok ***REMOVED***
		if !b ***REMOVED***
			p.proxyHasChecks(target.self.getOwnPropSym(s), target, s)
		***REMOVED***
		return b
	***REMOVED***

	return target.self.hasPropertySym(s)
***REMOVED***

func (p *proxyObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	return p.getOwnPropStr(name) != nil
***REMOVED***

func (p *proxyObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	return p.getOwnPropIdx(idx) != nil
***REMOVED***

func (p *proxyObject) hasOwnPropertySym(s *Symbol) bool ***REMOVED***
	return p.getOwnPropSym(s) != nil
***REMOVED***

func (p *proxyObject) proxyGetOwnPropertyDescriptor(targetProp Value, target *Object, trapResult Value, name fmt.Stringer) Value ***REMOVED***
	r := p.val.runtime
	targetDesc := propToValueProp(targetProp)
	var trapResultObj *Object
	if trapResult != nil && trapResult != _undefined ***REMOVED***
		if obj, ok := trapResult.(*Object); ok ***REMOVED***
			trapResultObj = obj
		***REMOVED*** else ***REMOVED***
			panic(r.NewTypeError("'getOwnPropertyDescriptor' on proxy: trap returned neither object nor undefined for property '%s'", name.String()))
		***REMOVED***
	***REMOVED***
	if trapResultObj == nil ***REMOVED***
		if targetDesc == nil ***REMOVED***
			return nil
		***REMOVED***
		if !targetDesc.configurable ***REMOVED***
			panic(r.NewTypeError())
		***REMOVED***
		if !target.self.isExtensible() ***REMOVED***
			panic(r.NewTypeError())
		***REMOVED***
		return nil
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

		if resultDesc.Writable == FLAG_FALSE && targetDesc.writable ***REMOVED***
			panic(r.NewTypeError("'getOwnPropertyDescriptor' on proxy: trap reported non-configurable and writable for property '%s' which is non-configurable, non-writable in the proxy target", name.String()))
		***REMOVED***
	***REMOVED***

	if resultDesc.Writable == FLAG_TRUE && resultDesc.Configurable == FLAG_TRUE &&
		resultDesc.Enumerable == FLAG_TRUE ***REMOVED***
		return resultDesc.Value
	***REMOVED***
	return r.toValueProp(trapResultObj)
***REMOVED***

func (p *proxyObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().getOwnPropertyDescriptorStr(target, name); ok ***REMOVED***
		return p.proxyGetOwnPropertyDescriptor(target.self.getOwnPropStr(name), target, v, name)
	***REMOVED***

	return target.self.getOwnPropStr(name)
***REMOVED***

func (p *proxyObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().getOwnPropertyDescriptorIdx(target, idx); ok ***REMOVED***
		return p.proxyGetOwnPropertyDescriptor(target.self.getOwnPropIdx(idx), target, v, idx)
	***REMOVED***

	return target.self.getOwnPropIdx(idx)
***REMOVED***

func (p *proxyObject) getOwnPropSym(s *Symbol) Value ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().getOwnPropertyDescriptorSym(target, s); ok ***REMOVED***
		return p.proxyGetOwnPropertyDescriptor(target.self.getOwnPropSym(s), target, v, s)
	***REMOVED***

	return target.self.getOwnPropSym(s)
***REMOVED***

func (p *proxyObject) proxyGetChecks(targetProp, trapResult Value, name fmt.Stringer) ***REMOVED***
	if targetDesc, ok := targetProp.(*valueProperty); ok ***REMOVED***
		if !targetDesc.accessor ***REMOVED***
			if !targetDesc.writable && !targetDesc.configurable && !trapResult.SameAs(targetDesc.value) ***REMOVED***
				panic(p.val.runtime.NewTypeError("'get' on proxy: property '%s' is a read-only and non-configurable data property on the proxy target but the proxy did not return its actual value (expected '%s' but got '%s')", name.String(), nilSafe(targetDesc.value), ret))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if !targetDesc.configurable && targetDesc.getterFunc == nil && trapResult != _undefined ***REMOVED***
				panic(p.val.runtime.NewTypeError("'get' on proxy: property '%s' is a non-configurable accessor property on the proxy target and does not have a getter function, but the trap did not return 'undefined' (got '%s')", name.String(), ret))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *proxyObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	target := p.target
	if receiver == nil ***REMOVED***
		receiver = p.val
	***REMOVED***
	if v, ok := p.checkHandler().getStr(target, name, receiver); ok ***REMOVED***
		p.proxyGetChecks(target.self.getOwnPropStr(name), v, name)
		return v
	***REMOVED***
	return target.self.getStr(name, receiver)
***REMOVED***

func (p *proxyObject) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	target := p.target
	if receiver == nil ***REMOVED***
		receiver = p.val
	***REMOVED***
	if v, ok := p.checkHandler().getIdx(target, idx, receiver); ok ***REMOVED***
		p.proxyGetChecks(target.self.getOwnPropIdx(idx), v, idx)
		return v
	***REMOVED***
	return target.self.getIdx(idx, receiver)
***REMOVED***

func (p *proxyObject) getSym(s *Symbol, receiver Value) Value ***REMOVED***
	target := p.target
	if receiver == nil ***REMOVED***
		receiver = p.val
	***REMOVED***
	if v, ok := p.checkHandler().getSym(target, s, receiver); ok ***REMOVED***
		p.proxyGetChecks(target.self.getOwnPropSym(s), v, s)
		return v
	***REMOVED***

	return target.self.getSym(s, receiver)
***REMOVED***

func (p *proxyObject) proxySetPreCheck(trapResult, throw bool, name fmt.Stringer) bool ***REMOVED***
	if !trapResult ***REMOVED***
		p.val.runtime.typeErrorResult(throw, "'set' on proxy: trap returned falsish for property '%s'", name.String())
	***REMOVED***
	return trapResult
***REMOVED***

func (p *proxyObject) proxySetPostCheck(targetProp, value Value, name fmt.Stringer) ***REMOVED***
	if prop, ok := targetProp.(*valueProperty); ok ***REMOVED***
		if prop.accessor ***REMOVED***
			if !prop.configurable && prop.setterFunc == nil ***REMOVED***
				panic(p.val.runtime.NewTypeError("'set' on proxy: trap returned truish for property '%s' which exists in the proxy target as a non-configurable and non-writable accessor property without a setter", name.String()))
			***REMOVED***
		***REMOVED*** else if !prop.configurable && !prop.writable && !p.__sameValue(prop.value, value) ***REMOVED***
			panic(p.val.runtime.NewTypeError("'set' on proxy: trap returned truish for property '%s' which exists in the proxy target as a non-configurable and non-writable data property with a different value", name.String()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *proxyObject) proxySetStr(name unistring.String, value, receiver Value, throw bool) bool ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().setStr(target, name, value, receiver); ok ***REMOVED***
		if p.proxySetPreCheck(v, throw, name) ***REMOVED***
			p.proxySetPostCheck(target.self.getOwnPropStr(name), value, name)
			return true
		***REMOVED***
		return false
	***REMOVED***
	return target.setStr(name, value, receiver, throw)
***REMOVED***

func (p *proxyObject) proxySetIdx(idx valueInt, value, receiver Value, throw bool) bool ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().setIdx(target, idx, value, receiver); ok ***REMOVED***
		if p.proxySetPreCheck(v, throw, idx) ***REMOVED***
			p.proxySetPostCheck(target.self.getOwnPropIdx(idx), value, idx)
			return true
		***REMOVED***
		return false
	***REMOVED***
	return target.setIdx(idx, value, receiver, throw)
***REMOVED***

func (p *proxyObject) proxySetSym(s *Symbol, value, receiver Value, throw bool) bool ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().setSym(target, s, value, receiver); ok ***REMOVED***
		if p.proxySetPreCheck(v, throw, s) ***REMOVED***
			p.proxySetPostCheck(target.self.getOwnPropSym(s), value, s)
			return true
		***REMOVED***
		return false
	***REMOVED***
	return target.setSym(s, value, receiver, throw)
***REMOVED***

func (p *proxyObject) setOwnStr(name unistring.String, v Value, throw bool) bool ***REMOVED***
	return p.proxySetStr(name, v, p.val, throw)
***REMOVED***

func (p *proxyObject) setOwnIdx(idx valueInt, v Value, throw bool) bool ***REMOVED***
	return p.proxySetIdx(idx, v, p.val, throw)
***REMOVED***

func (p *proxyObject) setOwnSym(s *Symbol, v Value, throw bool) bool ***REMOVED***
	return p.proxySetSym(s, v, p.val, throw)
***REMOVED***

func (p *proxyObject) setForeignStr(name unistring.String, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return p.proxySetStr(name, v, receiver, throw), true
***REMOVED***

func (p *proxyObject) setForeignIdx(idx valueInt, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return p.proxySetIdx(idx, v, receiver, throw), true
***REMOVED***

func (p *proxyObject) setForeignSym(s *Symbol, v, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return p.proxySetSym(s, v, receiver, throw), true
***REMOVED***

func (p *proxyObject) proxyDeleteCheck(trapResult bool, targetProp Value, name fmt.Stringer, target *Object) ***REMOVED***
	if trapResult ***REMOVED***
		if targetProp == nil ***REMOVED***
			return
		***REMOVED***
		if targetDesc, ok := targetProp.(*valueProperty); ok ***REMOVED***
			if !targetDesc.configurable ***REMOVED***
				panic(p.val.runtime.NewTypeError("'deleteProperty' on proxy: property '%s' is a non-configurable property but the trap returned truish", name.String()))
			***REMOVED***
		***REMOVED***
		if !target.self.isExtensible() ***REMOVED***
			panic(p.val.runtime.NewTypeError("'deleteProperty' on proxy: trap returned truish for property '%s' but the proxy target is non-extensible", name.String()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *proxyObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().deleteStr(target, name); ok ***REMOVED***
		p.proxyDeleteCheck(v, target.self.getOwnPropStr(name), name, target)
		return v
	***REMOVED***

	return target.self.deleteStr(name, throw)
***REMOVED***

func (p *proxyObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().deleteIdx(target, idx); ok ***REMOVED***
		p.proxyDeleteCheck(v, target.self.getOwnPropIdx(idx), idx, target)
		return v
	***REMOVED***

	return target.self.deleteIdx(idx, throw)
***REMOVED***

func (p *proxyObject) deleteSym(s *Symbol, throw bool) bool ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().deleteSym(target, s); ok ***REMOVED***
		p.proxyDeleteCheck(v, target.self.getOwnPropSym(s), s, target)
		return v
	***REMOVED***

	return target.self.deleteSym(s, throw)
***REMOVED***

func (p *proxyObject) ownPropertyKeys(all bool, _ []Value) []Value ***REMOVED***
	if v, ok := p.proxyOwnKeys(); ok ***REMOVED***
		if !all ***REMOVED***
			k := 0
			for i, key := range v ***REMOVED***
				prop := p.val.getOwnProp(key)
				if prop == nil ***REMOVED***
					continue
				***REMOVED***
				if prop, ok := prop.(*valueProperty); ok && !prop.enumerable ***REMOVED***
					continue
				***REMOVED***
				if k != i ***REMOVED***
					v[k] = v[i]
				***REMOVED***
				k++
			***REMOVED***
			v = v[:k]
		***REMOVED***
		return v
	***REMOVED***
	return p.target.self.ownPropertyKeys(all, nil)
***REMOVED***

func (p *proxyObject) proxyOwnKeys() ([]Value, bool) ***REMOVED***
	target := p.target
	if v, ok := p.checkHandler().ownKeys(target); ok ***REMOVED***
		keys := p.val.runtime.toObject(v)
		var keyList []Value
		keySet := make(map[Value]struct***REMOVED******REMOVED***)
		l := toLength(keys.self.getStr("length", nil))
		for k := int64(0); k < l; k++ ***REMOVED***
			item := keys.self.getIdx(valueInt(k), nil)
			if _, ok := item.(valueString); !ok ***REMOVED***
				if _, ok := item.(*Symbol); !ok ***REMOVED***
					panic(p.val.runtime.NewTypeError("%s is not a valid property name", item.String()))
				***REMOVED***
			***REMOVED***
			if _, exists := keySet[item]; exists ***REMOVED***
				panic(p.val.runtime.NewTypeError("'ownKeys' on proxy: trap returned duplicate entries"))
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

func (p *proxyObject) enumerateOwnKeys() iterNextFunc ***REMOVED***
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
		panic(p.val.runtime.NewTypeError("proxy target is not a function"))
	***REMOVED***
	if v, ok := p.checkHandler().apply(p.target, nilSafe(call.This), call.Arguments); ok ***REMOVED***
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
	if v, ok := p.checkHandler().construct(p.target, args, newTarget); ok ***REMOVED***
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

		if desc.IsGeneric() ***REMOVED***
			return true
		***REMOVED***

		if desc.IsData() != !current.accessor ***REMOVED***
			return desc.Configurable != FLAG_FALSE
		***REMOVED***

		if desc.IsData() && !current.accessor ***REMOVED***
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
		if desc.IsAccessor() && current.accessor ***REMOVED***
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
				if s, ok := val.(*Symbol); ok ***REMOVED***
					prop = p.getOwnPropSym(s)
				***REMOVED*** else ***REMOVED***
					continue
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if _, ok := val.(*Symbol); !ok ***REMOVED***
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
			if _, ok := val.(*Symbol); ok != symbols ***REMOVED***
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
		res := p.filterKeys(vals, all, true)
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
