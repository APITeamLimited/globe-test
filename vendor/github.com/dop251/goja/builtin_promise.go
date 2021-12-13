package goja

import (
	"github.com/dop251/goja/unistring"
	"reflect"
)

type PromiseState int
type PromiseRejectionOperation int

type promiseReactionType int

const (
	PromiseStatePending PromiseState = iota
	PromiseStateFulfilled
	PromiseStateRejected
)

const (
	PromiseRejectionReject PromiseRejectionOperation = iota
	PromiseRejectionHandle
)

const (
	promiseReactionFulfill promiseReactionType = iota
	promiseReactionReject
)

type PromiseRejectionTracker func(p *Promise, operation PromiseRejectionOperation)

type jobCallback struct ***REMOVED***
	callback func(FunctionCall) Value
***REMOVED***

type promiseCapability struct ***REMOVED***
	promise               *Object
	resolveObj, rejectObj *Object
***REMOVED***

type promiseReaction struct ***REMOVED***
	capability *promiseCapability
	typ        promiseReactionType
	handler    *jobCallback
***REMOVED***

var typePromise = reflect.TypeOf((*Promise)(nil))

// Promise is a Go wrapper around ECMAScript Promise. Calling Runtime.ToValue() on it
// returns the underlying Object. Calling Export() on a Promise Object returns a Promise.
//
// Use Runtime.NewPromise() to create one. Calling Runtime.ToValue() on a zero object or nil returns null Value.
//
// WARNING: Instances of Promise are not goroutine-safe. See Runtime.NewPromise() for more details.
type Promise struct ***REMOVED***
	baseObject
	state            PromiseState
	result           Value
	fulfillReactions []*promiseReaction
	rejectReactions  []*promiseReaction
	handled          bool
***REMOVED***

func (p *Promise) State() PromiseState ***REMOVED***
	return p.state
***REMOVED***

func (p *Promise) Result() Value ***REMOVED***
	return p.result
***REMOVED***

func (p *Promise) toValue(r *Runtime) Value ***REMOVED***
	if p == nil || p.val == nil ***REMOVED***
		return _null
	***REMOVED***
	promise := p.val
	if promise.runtime != r ***REMOVED***
		panic(r.NewTypeError("Illegal runtime transition of a Promise"))
	***REMOVED***
	return promise
***REMOVED***

func (p *Promise) createResolvingFunctions() (resolve, reject *Object) ***REMOVED***
	r := p.val.runtime
	alreadyResolved := false
	return p.val.runtime.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
			if alreadyResolved ***REMOVED***
				return _undefined
			***REMOVED***
			alreadyResolved = true
			resolution := call.Argument(0)
			if resolution.SameAs(p.val) ***REMOVED***
				return p.reject(r.NewTypeError("Promise self-resolution"))
			***REMOVED***
			if obj, ok := resolution.(*Object); ok ***REMOVED***
				var thenAction Value
				ex := r.vm.try(func() ***REMOVED***
					thenAction = obj.self.getStr("then", nil)
				***REMOVED***)
				if ex != nil ***REMOVED***
					return p.reject(ex.val)
				***REMOVED***
				if call, ok := assertCallable(thenAction); ok ***REMOVED***
					job := r.newPromiseResolveThenableJob(p, resolution, &jobCallback***REMOVED***callback: call***REMOVED***)
					r.enqueuePromiseJob(job)
					return _undefined
				***REMOVED***
			***REMOVED***
			return p.fulfill(resolution)
		***REMOVED***, nil, "", nil, 1),
		p.val.runtime.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
			if alreadyResolved ***REMOVED***
				return _undefined
			***REMOVED***
			alreadyResolved = true
			reason := call.Argument(0)
			return p.reject(reason)
		***REMOVED***, nil, "", nil, 1)
***REMOVED***

func (p *Promise) reject(reason Value) Value ***REMOVED***
	reactions := p.rejectReactions
	p.result = reason
	p.fulfillReactions, p.rejectReactions = nil, nil
	p.state = PromiseStateRejected
	r := p.val.runtime
	if !p.handled ***REMOVED***
		r.trackPromiseRejection(p, PromiseRejectionReject)
	***REMOVED***
	r.triggerPromiseReactions(reactions, reason)
	return _undefined
***REMOVED***

func (p *Promise) fulfill(value Value) Value ***REMOVED***
	reactions := p.fulfillReactions
	p.result = value
	p.fulfillReactions, p.rejectReactions = nil, nil
	p.state = PromiseStateFulfilled
	p.val.runtime.triggerPromiseReactions(reactions, value)
	return _undefined
***REMOVED***

func (p *Promise) exportType() reflect.Type ***REMOVED***
	return typePromise
***REMOVED***

func (p *Promise) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return p
***REMOVED***

func (r *Runtime) newPromiseResolveThenableJob(p *Promise, thenable Value, then *jobCallback) func() ***REMOVED***
	return func() ***REMOVED***
		resolve, reject := p.createResolvingFunctions()
		ex := r.vm.try(func() ***REMOVED***
			r.callJobCallback(then, thenable, resolve, reject)
		***REMOVED***)
		if ex != nil ***REMOVED***
			if fn, ok := reject.self.assertCallable(); ok ***REMOVED***
				fn(FunctionCall***REMOVED***Arguments: []Value***REMOVED***ex.val***REMOVED******REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) enqueuePromiseJob(job func()) ***REMOVED***
	r.jobQueue = append(r.jobQueue, job)
***REMOVED***

func (r *Runtime) triggerPromiseReactions(reactions []*promiseReaction, argument Value) ***REMOVED***
	for _, reaction := range reactions ***REMOVED***
		r.enqueuePromiseJob(r.newPromiseReactionJob(reaction, argument))
	***REMOVED***
***REMOVED***

func (r *Runtime) newPromiseReactionJob(reaction *promiseReaction, argument Value) func() ***REMOVED***
	return func() ***REMOVED***
		var handlerResult Value
		fulfill := false
		if reaction.handler == nil ***REMOVED***
			handlerResult = argument
			if reaction.typ == promiseReactionFulfill ***REMOVED***
				fulfill = true
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			ex := r.vm.try(func() ***REMOVED***
				handlerResult = r.callJobCallback(reaction.handler, _undefined, argument)
				fulfill = true
			***REMOVED***)
			if ex != nil ***REMOVED***
				handlerResult = ex.val
			***REMOVED***
		***REMOVED***
		if reaction.capability != nil ***REMOVED***
			if fulfill ***REMOVED***
				reaction.capability.resolve(handlerResult)
			***REMOVED*** else ***REMOVED***
				reaction.capability.reject(handlerResult)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) newPromise(proto *Object) *Promise ***REMOVED***
	o := &Object***REMOVED***runtime: r***REMOVED***

	po := &Promise***REMOVED******REMOVED***
	po.class = classPromise
	po.val = o
	po.extensible = true
	o.self = po
	po.prototype = proto
	po.init()
	return po
***REMOVED***

func (r *Runtime) builtin_newPromise(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("Promise"))
	***REMOVED***
	var arg0 Value
	if len(args) > 0 ***REMOVED***
		arg0 = args[0]
	***REMOVED***
	executor := r.toCallable(arg0)

	proto := r.getPrototypeFromCtor(newTarget, r.global.Promise, r.global.PromisePrototype)
	po := r.newPromise(proto)

	resolve, reject := po.createResolvingFunctions()
	ex := r.vm.try(func() ***REMOVED***
		executor(FunctionCall***REMOVED***Arguments: []Value***REMOVED***resolve, reject***REMOVED******REMOVED***)
	***REMOVED***)
	if ex != nil ***REMOVED***
		if fn, ok := reject.self.assertCallable(); ok ***REMOVED***
			fn(FunctionCall***REMOVED***Arguments: []Value***REMOVED***ex.val***REMOVED******REMOVED***)
		***REMOVED***
	***REMOVED***
	return po.val
***REMOVED***

func (r *Runtime) promiseProto_then(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	if p, ok := thisObj.self.(*Promise); ok ***REMOVED***
		c := r.speciesConstructorObj(thisObj, r.global.Promise)
		resultCapability := r.newPromiseCapability(c)
		return r.performPromiseThen(p, call.Argument(0), call.Argument(1), resultCapability)
	***REMOVED***
	panic(r.NewTypeError("Method Promise.prototype.then called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: thisObj***REMOVED***)))
***REMOVED***

func (r *Runtime) newPromiseCapability(c *Object) *promiseCapability ***REMOVED***
	pcap := new(promiseCapability)
	if c == r.global.Promise ***REMOVED***
		p := r.newPromise(r.global.PromisePrototype)
		pcap.resolveObj, pcap.rejectObj = p.createResolvingFunctions()
		pcap.promise = p.val
	***REMOVED*** else ***REMOVED***
		var resolve, reject Value
		executor := r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
			if resolve != nil ***REMOVED***
				panic(r.NewTypeError("resolve is already set"))
			***REMOVED***
			if reject != nil ***REMOVED***
				panic(r.NewTypeError("reject is already set"))
			***REMOVED***
			if arg := call.Argument(0); arg != _undefined ***REMOVED***
				resolve = arg
			***REMOVED***
			if arg := call.Argument(1); arg != _undefined ***REMOVED***
				reject = arg
			***REMOVED***
			return nil
		***REMOVED***, nil, "", nil, 2)
		pcap.promise = r.toConstructor(c)([]Value***REMOVED***executor***REMOVED***, c)
		pcap.resolveObj = r.toObject(resolve)
		r.toCallable(pcap.resolveObj) // make sure it's callable
		pcap.rejectObj = r.toObject(reject)
		r.toCallable(pcap.rejectObj)
	***REMOVED***
	return pcap
***REMOVED***

func (r *Runtime) performPromiseThen(p *Promise, onFulfilled, onRejected Value, resultCapability *promiseCapability) Value ***REMOVED***
	var onFulfilledJobCallback, onRejectedJobCallback *jobCallback
	if f, ok := assertCallable(onFulfilled); ok ***REMOVED***
		onFulfilledJobCallback = &jobCallback***REMOVED***callback: f***REMOVED***
	***REMOVED***
	if f, ok := assertCallable(onRejected); ok ***REMOVED***
		onRejectedJobCallback = &jobCallback***REMOVED***callback: f***REMOVED***
	***REMOVED***
	fulfillReaction := &promiseReaction***REMOVED***
		capability: resultCapability,
		typ:        promiseReactionFulfill,
		handler:    onFulfilledJobCallback,
	***REMOVED***
	rejectReaction := &promiseReaction***REMOVED***
		capability: resultCapability,
		typ:        promiseReactionReject,
		handler:    onRejectedJobCallback,
	***REMOVED***
	switch p.state ***REMOVED***
	case PromiseStatePending:
		p.fulfillReactions = append(p.fulfillReactions, fulfillReaction)
		p.rejectReactions = append(p.rejectReactions, rejectReaction)
	case PromiseStateFulfilled:
		r.enqueuePromiseJob(r.newPromiseReactionJob(fulfillReaction, p.result))
	default:
		reason := p.result
		if !p.handled ***REMOVED***
			r.trackPromiseRejection(p, PromiseRejectionHandle)
		***REMOVED***
		r.enqueuePromiseJob(r.newPromiseReactionJob(rejectReaction, reason))
	***REMOVED***
	p.handled = true
	if resultCapability == nil ***REMOVED***
		return _undefined
	***REMOVED***
	return resultCapability.promise
***REMOVED***

func (r *Runtime) promiseProto_catch(call FunctionCall) Value ***REMOVED***
	return r.invoke(call.This, "then", _undefined, call.Argument(0))
***REMOVED***

func (r *Runtime) promiseResolve(c *Object, x Value) *Object ***REMOVED***
	if obj, ok := x.(*Object); ok ***REMOVED***
		xConstructor := nilSafe(obj.self.getStr("constructor", nil))
		if xConstructor.SameAs(c) ***REMOVED***
			return obj
		***REMOVED***
	***REMOVED***
	pcap := r.newPromiseCapability(c)
	pcap.resolve(x)
	return pcap.promise
***REMOVED***

func (r *Runtime) promiseProto_finally(call FunctionCall) Value ***REMOVED***
	promise := r.toObject(call.This)
	c := r.speciesConstructorObj(promise, r.global.Promise)
	onFinally := call.Argument(0)
	var thenFinally, catchFinally Value
	if onFinallyFn, ok := assertCallable(onFinally); !ok ***REMOVED***
		thenFinally, catchFinally = onFinally, onFinally
	***REMOVED*** else ***REMOVED***
		thenFinally = r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
			value := call.Argument(0)
			result := onFinallyFn(FunctionCall***REMOVED******REMOVED***)
			promise := r.promiseResolve(c, result)
			valueThunk := r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
				return value
			***REMOVED***, nil, "", nil, 0)
			return r.invoke(promise, "then", valueThunk)
		***REMOVED***, nil, "", nil, 1)

		catchFinally = r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
			reason := call.Argument(0)
			result := onFinallyFn(FunctionCall***REMOVED******REMOVED***)
			promise := r.promiseResolve(c, result)
			thrower := r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
				panic(reason)
			***REMOVED***, nil, "", nil, 0)
			return r.invoke(promise, "then", thrower)
		***REMOVED***, nil, "", nil, 1)
	***REMOVED***
	return r.invoke(promise, "then", thenFinally, catchFinally)
***REMOVED***

func (pcap *promiseCapability) resolve(result Value) ***REMOVED***
	pcap.promise.runtime.toCallable(pcap.resolveObj)(FunctionCall***REMOVED***Arguments: []Value***REMOVED***result***REMOVED******REMOVED***)
***REMOVED***

func (pcap *promiseCapability) reject(reason Value) ***REMOVED***
	pcap.promise.runtime.toCallable(pcap.rejectObj)(FunctionCall***REMOVED***Arguments: []Value***REMOVED***reason***REMOVED******REMOVED***)
***REMOVED***

func (pcap *promiseCapability) try(f func()) bool ***REMOVED***
	ex := pcap.promise.runtime.vm.try(f)
	if ex != nil ***REMOVED***
		pcap.reject(ex.val)
		return false
	***REMOVED***
	return true
***REMOVED***

func (r *Runtime) promise_all(call FunctionCall) Value ***REMOVED***
	c := r.toObject(call.This)
	pcap := r.newPromiseCapability(c)

	pcap.try(func() ***REMOVED***
		promiseResolve := r.toCallable(c.self.getStr("resolve", nil))
		iter := r.getIterator(call.Argument(0), nil)
		var values []Value
		remainingElementsCount := 1
		iter.iterate(func(nextValue Value) ***REMOVED***
			index := len(values)
			values = append(values, _undefined)
			nextPromise := promiseResolve(FunctionCall***REMOVED***This: c, Arguments: []Value***REMOVED***nextValue***REMOVED******REMOVED***)
			alreadyCalled := false
			onFulfilled := r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
				if alreadyCalled ***REMOVED***
					return _undefined
				***REMOVED***
				alreadyCalled = true
				values[index] = call.Argument(0)
				remainingElementsCount--
				if remainingElementsCount == 0 ***REMOVED***
					pcap.resolve(r.newArrayValues(values))
				***REMOVED***
				return _undefined
			***REMOVED***, nil, "", nil, 1)
			remainingElementsCount++
			r.invoke(nextPromise, "then", onFulfilled, pcap.rejectObj)
		***REMOVED***)
		remainingElementsCount--
		if remainingElementsCount == 0 ***REMOVED***
			pcap.resolve(r.newArrayValues(values))
		***REMOVED***
	***REMOVED***)
	return pcap.promise
***REMOVED***

func (r *Runtime) promise_allSettled(call FunctionCall) Value ***REMOVED***
	c := r.toObject(call.This)
	pcap := r.newPromiseCapability(c)

	pcap.try(func() ***REMOVED***
		promiseResolve := r.toCallable(c.self.getStr("resolve", nil))
		iter := r.getIterator(call.Argument(0), nil)
		var values []Value
		remainingElementsCount := 1
		iter.iterate(func(nextValue Value) ***REMOVED***
			index := len(values)
			values = append(values, _undefined)
			nextPromise := promiseResolve(FunctionCall***REMOVED***This: c, Arguments: []Value***REMOVED***nextValue***REMOVED******REMOVED***)
			alreadyCalled := false
			reaction := func(status Value, valueKey unistring.String) *Object ***REMOVED***
				return r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
					if alreadyCalled ***REMOVED***
						return _undefined
					***REMOVED***
					alreadyCalled = true
					obj := r.NewObject()
					obj.self._putProp("status", status, true, true, true)
					obj.self._putProp(valueKey, call.Argument(0), true, true, true)
					values[index] = obj
					remainingElementsCount--
					if remainingElementsCount == 0 ***REMOVED***
						pcap.resolve(r.newArrayValues(values))
					***REMOVED***
					return _undefined
				***REMOVED***, nil, "", nil, 1)
			***REMOVED***
			onFulfilled := reaction(asciiString("fulfilled"), "value")
			onRejected := reaction(asciiString("rejected"), "reason")
			remainingElementsCount++
			r.invoke(nextPromise, "then", onFulfilled, onRejected)
		***REMOVED***)
		remainingElementsCount--
		if remainingElementsCount == 0 ***REMOVED***
			pcap.resolve(r.newArrayValues(values))
		***REMOVED***
	***REMOVED***)
	return pcap.promise
***REMOVED***

func (r *Runtime) promise_any(call FunctionCall) Value ***REMOVED***
	c := r.toObject(call.This)
	pcap := r.newPromiseCapability(c)

	pcap.try(func() ***REMOVED***
		promiseResolve := r.toCallable(c.self.getStr("resolve", nil))
		iter := r.getIterator(call.Argument(0), nil)
		var errors []Value
		remainingElementsCount := 1
		iter.iterate(func(nextValue Value) ***REMOVED***
			index := len(errors)
			errors = append(errors, _undefined)
			nextPromise := promiseResolve(FunctionCall***REMOVED***This: c, Arguments: []Value***REMOVED***nextValue***REMOVED******REMOVED***)
			alreadyCalled := false
			onRejected := r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
				if alreadyCalled ***REMOVED***
					return _undefined
				***REMOVED***
				alreadyCalled = true
				errors[index] = call.Argument(0)
				remainingElementsCount--
				if remainingElementsCount == 0 ***REMOVED***
					_error := r.builtin_new(r.global.AggregateError, nil)
					_error.self._putProp("errors", r.newArrayValues(errors), true, false, true)
					pcap.reject(_error)
				***REMOVED***
				return _undefined
			***REMOVED***, nil, "", nil, 1)

			remainingElementsCount++
			r.invoke(nextPromise, "then", pcap.resolveObj, onRejected)
		***REMOVED***)
		remainingElementsCount--
		if remainingElementsCount == 0 ***REMOVED***
			_error := r.builtin_new(r.global.AggregateError, nil)
			_error.self._putProp("errors", r.newArrayValues(errors), true, false, true)
			pcap.reject(_error)
		***REMOVED***
	***REMOVED***)
	return pcap.promise
***REMOVED***

func (r *Runtime) promise_race(call FunctionCall) Value ***REMOVED***
	c := r.toObject(call.This)
	pcap := r.newPromiseCapability(c)

	pcap.try(func() ***REMOVED***
		promiseResolve := r.toCallable(c.self.getStr("resolve", nil))
		iter := r.getIterator(call.Argument(0), nil)
		iter.iterate(func(nextValue Value) ***REMOVED***
			nextPromise := promiseResolve(FunctionCall***REMOVED***This: c, Arguments: []Value***REMOVED***nextValue***REMOVED******REMOVED***)
			r.invoke(nextPromise, "then", pcap.resolveObj, pcap.rejectObj)
		***REMOVED***)
	***REMOVED***)
	return pcap.promise
***REMOVED***

func (r *Runtime) promise_reject(call FunctionCall) Value ***REMOVED***
	pcap := r.newPromiseCapability(r.toObject(call.This))
	pcap.reject(call.Argument(0))
	return pcap.promise
***REMOVED***

func (r *Runtime) promise_resolve(call FunctionCall) Value ***REMOVED***
	return r.promiseResolve(r.toObject(call.This), call.Argument(0))
***REMOVED***

func (r *Runtime) createPromiseProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)
	o._putProp("constructor", r.global.Promise, true, false, true)

	o._putProp("catch", r.newNativeFunc(r.promiseProto_catch, nil, "catch", nil, 1), true, false, true)
	o._putProp("finally", r.newNativeFunc(r.promiseProto_finally, nil, "finally", nil, 1), true, false, true)
	o._putProp("then", r.newNativeFunc(r.promiseProto_then, nil, "then", nil, 2), true, false, true)

	o._putSym(SymToStringTag, valueProp(asciiString(classPromise), false, false, true))

	return o
***REMOVED***

func (r *Runtime) createPromise(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.builtin_newPromise, r.global.PromisePrototype, "Promise", 1)

	o._putProp("all", r.newNativeFunc(r.promise_all, nil, "all", nil, 1), true, false, true)
	o._putProp("allSettled", r.newNativeFunc(r.promise_allSettled, nil, "allSettled", nil, 1), true, false, true)
	o._putProp("any", r.newNativeFunc(r.promise_any, nil, "any", nil, 1), true, false, true)
	o._putProp("race", r.newNativeFunc(r.promise_race, nil, "race", nil, 1), true, false, true)
	o._putProp("reject", r.newNativeFunc(r.promise_reject, nil, "reject", nil, 1), true, false, true)
	o._putProp("resolve", r.newNativeFunc(r.promise_resolve, nil, "resolve", nil, 1), true, false, true)

	r.putSpeciesReturnThis(o)

	return o
***REMOVED***

func (r *Runtime) initPromise() ***REMOVED***
	r.global.PromisePrototype = r.newLazyObject(r.createPromiseProto)
	r.global.Promise = r.newLazyObject(r.createPromise)

	r.addToGlobal("Promise", r.global.Promise)
***REMOVED***

func (r *Runtime) wrapPromiseReaction(fObj *Object) func(interface***REMOVED******REMOVED***) ***REMOVED***
	f, _ := AssertFunction(fObj)
	return func(x interface***REMOVED******REMOVED***) ***REMOVED***
		_, _ = f(nil, r.ToValue(x))
	***REMOVED***
***REMOVED***

// NewPromise creates and returns a Promise and resolving functions for it.
//
// WARNING: The returned values are not goroutine-safe and must not be called in parallel with VM running.
// In order to make use of this method you need an event loop such as the one in goja_nodejs (https://github.com/dop251/goja_nodejs)
// where it can be used like this:
//
//  loop := NewEventLoop()
//  loop.Start()
//  defer loop.Stop()
//  loop.RunOnLoop(func(vm *goja.Runtime) ***REMOVED***
//		p, resolve, _ := vm.NewPromise()
//		vm.Set("p", p)
//      go func() ***REMOVED***
//   		time.Sleep(500 * time.Millisecond)   // or perform any other blocking operation
//			loop.RunOnLoop(func(*goja.Runtime) ***REMOVED*** // resolve() must be called on the loop, cannot call it here
//				resolve(result)
//			***REMOVED***)
//		***REMOVED***()
//  ***REMOVED***
func (r *Runtime) NewPromise() (promise *Promise, resolve func(result interface***REMOVED******REMOVED***), reject func(reason interface***REMOVED******REMOVED***)) ***REMOVED***
	p := r.newPromise(r.global.PromisePrototype)
	resolveF, rejectF := p.createResolvingFunctions()
	return p, r.wrapPromiseReaction(resolveF), r.wrapPromiseReaction(rejectF)
***REMOVED***

// SetPromiseRejectionTracker registers a function that will be called in two scenarios: when a promise is rejected
// without any handlers (with operation argument set to PromiseRejectionReject), and when a handler is added to a
// rejected promise for the first time (with operation argument set to PromiseRejectionHandle).
//
// Setting a tracker replaces any existing one. Setting it to nil disables the functionality.
//
// See https://tc39.es/ecma262/#sec-host-promise-rejection-tracker for more details.
func (r *Runtime) SetPromiseRejectionTracker(tracker PromiseRejectionTracker) ***REMOVED***
	r.promiseRejectionTracker = tracker
***REMOVED***
