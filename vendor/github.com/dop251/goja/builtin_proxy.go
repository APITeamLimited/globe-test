package goja

import (
	"github.com/dop251/goja/unistring"
)

type nativeProxyHandler struct ***REMOVED***
	handler *ProxyTrapConfig
***REMOVED***

func (h *nativeProxyHandler) getPrototypeOf(target *Object) (Value, bool) ***REMOVED***
	if trap := h.handler.GetPrototypeOf; trap != nil ***REMOVED***
		return trap(target), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) setPrototypeOf(target *Object, proto *Object) (bool, bool) ***REMOVED***
	if trap := h.handler.SetPrototypeOf; trap != nil ***REMOVED***
		return trap(target, proto), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) isExtensible(target *Object) (bool, bool) ***REMOVED***
	if trap := h.handler.IsExtensible; trap != nil ***REMOVED***
		return trap(target), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) preventExtensions(target *Object) (bool, bool) ***REMOVED***
	if trap := h.handler.PreventExtensions; trap != nil ***REMOVED***
		return trap(target), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) getOwnPropertyDescriptorStr(target *Object, prop unistring.String) (Value, bool) ***REMOVED***
	if trap := h.handler.GetOwnPropertyDescriptorIdx; trap != nil ***REMOVED***
		if idx, ok := strToInt(prop); ok ***REMOVED***
			desc := trap(target, idx)
			return desc.toValue(target.runtime), true
		***REMOVED***
	***REMOVED***
	if trap := h.handler.GetOwnPropertyDescriptor; trap != nil ***REMOVED***
		desc := trap(target, prop.String())
		return desc.toValue(target.runtime), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) getOwnPropertyDescriptorIdx(target *Object, prop valueInt) (Value, bool) ***REMOVED***
	if trap := h.handler.GetOwnPropertyDescriptorIdx; trap != nil ***REMOVED***
		desc := trap(target, toIntStrict(int64(prop)))
		return desc.toValue(target.runtime), true
	***REMOVED***
	if trap := h.handler.GetOwnPropertyDescriptor; trap != nil ***REMOVED***
		desc := trap(target, prop.String())
		return desc.toValue(target.runtime), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) getOwnPropertyDescriptorSym(target *Object, prop *Symbol) (Value, bool) ***REMOVED***
	if trap := h.handler.GetOwnPropertyDescriptorSym; trap != nil ***REMOVED***
		desc := trap(target, prop)
		return desc.toValue(target.runtime), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) definePropertyStr(target *Object, prop unistring.String, desc PropertyDescriptor) (bool, bool) ***REMOVED***
	if trap := h.handler.DefinePropertyIdx; trap != nil ***REMOVED***
		if idx, ok := strToInt(prop); ok ***REMOVED***
			return trap(target, idx, desc), true
		***REMOVED***
	***REMOVED***
	if trap := h.handler.DefineProperty; trap != nil ***REMOVED***
		return trap(target, prop.String(), desc), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) definePropertyIdx(target *Object, prop valueInt, desc PropertyDescriptor) (bool, bool) ***REMOVED***
	if trap := h.handler.DefinePropertyIdx; trap != nil ***REMOVED***
		return trap(target, toIntStrict(int64(prop)), desc), true
	***REMOVED***
	if trap := h.handler.DefineProperty; trap != nil ***REMOVED***
		return trap(target, prop.String(), desc), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) definePropertySym(target *Object, prop *Symbol, desc PropertyDescriptor) (bool, bool) ***REMOVED***
	if trap := h.handler.DefinePropertySym; trap != nil ***REMOVED***
		return trap(target, prop, desc), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) hasStr(target *Object, prop unistring.String) (bool, bool) ***REMOVED***
	if trap := h.handler.HasIdx; trap != nil ***REMOVED***
		if idx, ok := strToInt(prop); ok ***REMOVED***
			return trap(target, idx), true
		***REMOVED***
	***REMOVED***
	if trap := h.handler.Has; trap != nil ***REMOVED***
		return trap(target, prop.String()), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) hasIdx(target *Object, prop valueInt) (bool, bool) ***REMOVED***
	if trap := h.handler.HasIdx; trap != nil ***REMOVED***
		return trap(target, toIntStrict(int64(prop))), true
	***REMOVED***
	if trap := h.handler.Has; trap != nil ***REMOVED***
		return trap(target, prop.String()), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) hasSym(target *Object, prop *Symbol) (bool, bool) ***REMOVED***
	if trap := h.handler.HasSym; trap != nil ***REMOVED***
		return trap(target, prop), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) getStr(target *Object, prop unistring.String, receiver Value) (Value, bool) ***REMOVED***
	if trap := h.handler.GetIdx; trap != nil ***REMOVED***
		if idx, ok := strToInt(prop); ok ***REMOVED***
			return trap(target, idx, receiver), true
		***REMOVED***
	***REMOVED***
	if trap := h.handler.Get; trap != nil ***REMOVED***
		return trap(target, prop.String(), receiver), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) getIdx(target *Object, prop valueInt, receiver Value) (Value, bool) ***REMOVED***
	if trap := h.handler.GetIdx; trap != nil ***REMOVED***
		return trap(target, toIntStrict(int64(prop)), receiver), true
	***REMOVED***
	if trap := h.handler.Get; trap != nil ***REMOVED***
		return trap(target, prop.String(), receiver), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) getSym(target *Object, prop *Symbol, receiver Value) (Value, bool) ***REMOVED***
	if trap := h.handler.GetSym; trap != nil ***REMOVED***
		return trap(target, prop, receiver), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) setStr(target *Object, prop unistring.String, value Value, receiver Value) (bool, bool) ***REMOVED***
	if trap := h.handler.SetIdx; trap != nil ***REMOVED***
		if idx, ok := strToInt(prop); ok ***REMOVED***
			return trap(target, idx, value, receiver), true
		***REMOVED***
	***REMOVED***
	if trap := h.handler.Set; trap != nil ***REMOVED***
		return trap(target, prop.String(), value, receiver), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) setIdx(target *Object, prop valueInt, value Value, receiver Value) (bool, bool) ***REMOVED***
	if trap := h.handler.SetIdx; trap != nil ***REMOVED***
		return trap(target, toIntStrict(int64(prop)), value, receiver), true
	***REMOVED***
	if trap := h.handler.Set; trap != nil ***REMOVED***
		return trap(target, prop.String(), value, receiver), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) setSym(target *Object, prop *Symbol, value Value, receiver Value) (bool, bool) ***REMOVED***
	if trap := h.handler.SetSym; trap != nil ***REMOVED***
		return trap(target, prop, value, receiver), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) deleteStr(target *Object, prop unistring.String) (bool, bool) ***REMOVED***
	if trap := h.handler.DeletePropertyIdx; trap != nil ***REMOVED***
		if idx, ok := strToInt(prop); ok ***REMOVED***
			return trap(target, idx), true
		***REMOVED***
	***REMOVED***
	if trap := h.handler.DeleteProperty; trap != nil ***REMOVED***
		return trap(target, prop.String()), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) deleteIdx(target *Object, prop valueInt) (bool, bool) ***REMOVED***
	if trap := h.handler.DeletePropertyIdx; trap != nil ***REMOVED***
		return trap(target, toIntStrict(int64(prop))), true
	***REMOVED***
	if trap := h.handler.DeleteProperty; trap != nil ***REMOVED***
		return trap(target, prop.String()), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) deleteSym(target *Object, prop *Symbol) (bool, bool) ***REMOVED***
	if trap := h.handler.DeletePropertySym; trap != nil ***REMOVED***
		return trap(target, prop), true
	***REMOVED***
	return false, false
***REMOVED***

func (h *nativeProxyHandler) ownKeys(target *Object) (*Object, bool) ***REMOVED***
	if trap := h.handler.OwnKeys; trap != nil ***REMOVED***
		return trap(target), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) apply(target *Object, this Value, args []Value) (Value, bool) ***REMOVED***
	if trap := h.handler.Apply; trap != nil ***REMOVED***
		return trap(target, this, args), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) construct(target *Object, args []Value, newTarget *Object) (Value, bool) ***REMOVED***
	if trap := h.handler.Construct; trap != nil ***REMOVED***
		return trap(target, args, newTarget), true
	***REMOVED***
	return nil, false
***REMOVED***

func (h *nativeProxyHandler) toObject(runtime *Runtime) *Object ***REMOVED***
	return runtime.ToValue(h.handler).ToObject(runtime)
***REMOVED***

func (r *Runtime) newNativeProxyHandler(nativeHandler *ProxyTrapConfig) proxyHandler ***REMOVED***
	return &nativeProxyHandler***REMOVED***handler: nativeHandler***REMOVED***
***REMOVED***

// ProxyTrapConfig provides a simplified Go-friendly API for implementing Proxy traps.
// If an *Idx trap is defined it gets called for integer property keys, including negative ones. Note that
// this only includes string property keys that represent a canonical integer
// (i.e. "0", "123", but not "00", "01", " 1" or "-0").
// For efficiency strings representing integers exceeding 2^53 are not checked to see if they are canonical,
// i.e. the *Idx traps will receive "9007199254740993" as well as "9007199254740994", even though the former is not
// a canonical representation in ECMAScript (Number("9007199254740993") === 9007199254740992).
// See https://262.ecma-international.org/#sec-canonicalnumericindexstring
// If an *Idx trap is not set, the corresponding string one is used.
type ProxyTrapConfig struct ***REMOVED***
	// A trap for Object.getPrototypeOf, Reflect.getPrototypeOf, __proto__, Object.prototype.isPrototypeOf, instanceof
	GetPrototypeOf func(target *Object) (prototype *Object)

	// A trap for Object.setPrototypeOf, Reflect.setPrototypeOf
	SetPrototypeOf func(target *Object, prototype *Object) (success bool)

	// A trap for Object.isExtensible, Reflect.isExtensible
	IsExtensible func(target *Object) (success bool)

	// A trap for Object.preventExtensions, Reflect.preventExtensions
	PreventExtensions func(target *Object) (success bool)

	// A trap for Object.getOwnPropertyDescriptor, Reflect.getOwnPropertyDescriptor (string properties)
	GetOwnPropertyDescriptor func(target *Object, prop string) (propertyDescriptor PropertyDescriptor)

	// A trap for Object.getOwnPropertyDescriptor, Reflect.getOwnPropertyDescriptor (integer properties)
	GetOwnPropertyDescriptorIdx func(target *Object, prop int) (propertyDescriptor PropertyDescriptor)

	// A trap for Object.getOwnPropertyDescriptor, Reflect.getOwnPropertyDescriptor (Symbol properties)
	GetOwnPropertyDescriptorSym func(target *Object, prop *Symbol) (propertyDescriptor PropertyDescriptor)

	// A trap for Object.defineProperty, Reflect.defineProperty (string properties)
	DefineProperty func(target *Object, key string, propertyDescriptor PropertyDescriptor) (success bool)

	// A trap for Object.defineProperty, Reflect.defineProperty (integer properties)
	DefinePropertyIdx func(target *Object, key int, propertyDescriptor PropertyDescriptor) (success bool)

	// A trap for Object.defineProperty, Reflect.defineProperty (Symbol properties)
	DefinePropertySym func(target *Object, key *Symbol, propertyDescriptor PropertyDescriptor) (success bool)

	// A trap for the in operator, with operator, Reflect.has (string properties)
	Has func(target *Object, property string) (available bool)

	// A trap for the in operator, with operator, Reflect.has (integer properties)
	HasIdx func(target *Object, property int) (available bool)

	// A trap for the in operator, with operator, Reflect.has (Symbol properties)
	HasSym func(target *Object, property *Symbol) (available bool)

	// A trap for getting property values, Reflect.get (string properties)
	Get func(target *Object, property string, receiver Value) (value Value)

	// A trap for getting property values, Reflect.get (integer properties)
	GetIdx func(target *Object, property int, receiver Value) (value Value)

	// A trap for getting property values, Reflect.get (Symbol properties)
	GetSym func(target *Object, property *Symbol, receiver Value) (value Value)

	// A trap for setting property values, Reflect.set (string properties)
	Set func(target *Object, property string, value Value, receiver Value) (success bool)

	// A trap for setting property values, Reflect.set (integer properties)
	SetIdx func(target *Object, property int, value Value, receiver Value) (success bool)

	// A trap for setting property values, Reflect.set (Symbol properties)
	SetSym func(target *Object, property *Symbol, value Value, receiver Value) (success bool)

	// A trap for the delete operator, Reflect.deleteProperty (string properties)
	DeleteProperty func(target *Object, property string) (success bool)

	// A trap for the delete operator, Reflect.deleteProperty (integer properties)
	DeletePropertyIdx func(target *Object, property int) (success bool)

	// A trap for the delete operator, Reflect.deleteProperty (Symbol properties)
	DeletePropertySym func(target *Object, property *Symbol) (success bool)

	// A trap for Object.getOwnPropertyNames, Object.getOwnPropertySymbols, Object.keys, Reflect.ownKeys
	OwnKeys func(target *Object) (object *Object)

	// A trap for a function call, Function.prototype.apply, Function.prototype.call, Reflect.apply
	Apply func(target *Object, this Value, argumentsList []Value) (value Value)

	// A trap for the new operator, Reflect.construct
	Construct func(target *Object, argumentsList []Value, newTarget *Object) (value *Object)
***REMOVED***

func (r *Runtime) newProxy(args []Value, proto *Object) *Object ***REMOVED***
	if len(args) >= 2 ***REMOVED***
		if target, ok := args[0].(*Object); ok ***REMOVED***
			if proxyHandler, ok := args[1].(*Object); ok ***REMOVED***
				return r.newProxyObject(target, proxyHandler, proto).val
			***REMOVED***
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Cannot create proxy with a non-object as target or handler"))
***REMOVED***

func (r *Runtime) builtin_newProxy(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("Proxy"))
	***REMOVED***
	return r.newProxy(args, r.getPrototypeFromCtor(newTarget, r.global.Proxy, r.global.ObjectPrototype))
***REMOVED***

func (r *Runtime) NewProxy(target *Object, nativeHandler *ProxyTrapConfig) Proxy ***REMOVED***
	if p, ok := target.self.(*proxyObject); ok ***REMOVED***
		if p.handler == nil ***REMOVED***
			panic(r.NewTypeError("Cannot create proxy with a revoked proxy as target"))
		***REMOVED***
	***REMOVED***
	handler := r.newNativeProxyHandler(nativeHandler)
	proxy := r._newProxyObject(target, handler, nil)
	return Proxy***REMOVED***proxy: proxy***REMOVED***
***REMOVED***

func (r *Runtime) builtin_proxy_revocable(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) >= 2 ***REMOVED***
		if target, ok := call.Argument(0).(*Object); ok ***REMOVED***
			if proxyHandler, ok := call.Argument(1).(*Object); ok ***REMOVED***
				proxy := r.newProxyObject(target, proxyHandler, nil)
				revoke := r.newNativeFunc(func(FunctionCall) Value ***REMOVED***
					proxy.revoke()
					return _undefined
				***REMOVED***, nil, "", nil, 0)
				ret := r.NewObject()
				ret.self._putProp("proxy", proxy.val, true, true, true)
				ret.self._putProp("revoke", revoke, true, true, true)
				return ret
			***REMOVED***
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Cannot create proxy with a non-object as target or handler"))
***REMOVED***

func (r *Runtime) createProxy(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.builtin_newProxy, nil, "Proxy", 2)

	o._putProp("revocable", r.newNativeFunc(r.builtin_proxy_revocable, nil, "revocable", nil, 2), true, false, true)
	return o
***REMOVED***

func (r *Runtime) initProxy() ***REMOVED***
	r.global.Proxy = r.newLazyObject(r.createProxy)
	r.addToGlobal("Proxy", r.global.Proxy)
***REMOVED***
