package goja

import "github.com/dop251/goja/unistring"

const propNameStack = "stack"

type errorObject struct ***REMOVED***
	baseObject
	stack          []StackFrame
	stackPropAdded bool
***REMOVED***

func (e *errorObject) formatStack() valueString ***REMOVED***
	var b valueStringBuilder
	if name := e.getStr("name", nil); name != nil ***REMOVED***
		b.WriteString(name.toString())
		b.WriteRune('\n')
	***REMOVED*** else ***REMOVED***
		b.WriteASCII("Error\n")
	***REMOVED***

	for _, frame := range e.stack ***REMOVED***
		b.WriteASCII("\tat ")
		frame.WriteToValueBuilder(&b)
		b.WriteRune('\n')
	***REMOVED***
	return b.String()
***REMOVED***

func (e *errorObject) addStackProp() Value ***REMOVED***
	if !e.stackPropAdded ***REMOVED***
		res := e._putProp(propNameStack, e.formatStack(), true, false, true)
		if len(e.propNames) > 1 ***REMOVED***
			// reorder property names to ensure 'stack' is the first one
			copy(e.propNames[1:], e.propNames)
			e.propNames[0] = propNameStack
		***REMOVED***
		e.stackPropAdded = true
		return res
	***REMOVED***
	return nil
***REMOVED***

func (e *errorObject) getStr(p unistring.String, receiver Value) Value ***REMOVED***
	return e.getStrWithOwnProp(e.getOwnPropStr(p), p, receiver)
***REMOVED***

func (e *errorObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	res := e.baseObject.getOwnPropStr(name)
	if res == nil && name == propNameStack ***REMOVED***
		return e.addStackProp()
	***REMOVED***

	return res
***REMOVED***

func (e *errorObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if name == propNameStack ***REMOVED***
		e.addStackProp()
	***REMOVED***
	return e.baseObject.setOwnStr(name, val, throw)
***REMOVED***

func (e *errorObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return e._setForeignStr(name, e.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

func (e *errorObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if name == propNameStack ***REMOVED***
		e.addStackProp()
	***REMOVED***
	return e.baseObject.deleteStr(name, throw)
***REMOVED***

func (e *errorObject) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	if name == propNameStack ***REMOVED***
		e.addStackProp()
	***REMOVED***
	return e.baseObject.defineOwnPropertyStr(name, desc, throw)
***REMOVED***

func (e *errorObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if e.baseObject.hasOwnPropertyStr(name) ***REMOVED***
		return true
	***REMOVED***

	return name == propNameStack && !e.stackPropAdded
***REMOVED***

func (e *errorObject) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	if all && !e.stackPropAdded ***REMOVED***
		accum = append(accum, asciiString(propNameStack))
	***REMOVED***
	return e.baseObject.stringKeys(all, accum)
***REMOVED***

func (e *errorObject) iterateStringKeys() iterNextFunc ***REMOVED***
	e.addStackProp()
	return e.baseObject.iterateStringKeys()
***REMOVED***

func (e *errorObject) init() ***REMOVED***
	e.baseObject.init()
	vm := e.val.runtime.vm
	e.stack = vm.captureStack(make([]StackFrame, 0, len(vm.callStack)+1), 0)
***REMOVED***

func (r *Runtime) newErrorObject(proto *Object, class string) *errorObject ***REMOVED***
	obj := &Object***REMOVED***runtime: r***REMOVED***
	o := &errorObject***REMOVED***
		baseObject: baseObject***REMOVED***
			class:      class,
			val:        obj,
			extensible: true,
			prototype:  proto,
		***REMOVED***,
	***REMOVED***
	obj.self = o
	o.init()
	return o
***REMOVED***

func (r *Runtime) builtin_Error(args []Value, proto *Object) *Object ***REMOVED***
	obj := r.newErrorObject(proto, classError)
	if len(args) > 0 && args[0] != _undefined ***REMOVED***
		obj._putProp("message", args[0], true, false, true)
	***REMOVED***
	return obj.val
***REMOVED***

func (r *Runtime) builtin_AggregateError(args []Value, proto *Object) *Object ***REMOVED***
	obj := r.newErrorObject(proto, classAggError)
	if len(args) > 1 && args[1] != nil && args[1] != _undefined ***REMOVED***
		obj._putProp("message", args[1].toString(), true, false, true)
	***REMOVED***
	var errors []Value
	if len(args) > 0 ***REMOVED***
		errors = r.iterableToList(args[0], nil)
	***REMOVED***
	obj._putProp("errors", r.newArrayValues(errors), true, false, true)

	return obj.val
***REMOVED***

func (r *Runtime) createErrorPrototype(name valueString) *Object ***REMOVED***
	o := r.newBaseObject(r.global.ErrorPrototype, classObject)
	o._putProp("message", stringEmpty, true, false, true)
	o._putProp("name", name, true, false, true)
	return o.val
***REMOVED***

func (r *Runtime) initErrors() ***REMOVED***
	r.global.ErrorPrototype = r.NewObject()
	o := r.global.ErrorPrototype.self
	o._putProp("message", stringEmpty, true, false, true)
	o._putProp("name", stringError, true, false, true)
	o._putProp("toString", r.newNativeFunc(r.error_toString, nil, "toString", nil, 0), true, false, true)

	r.global.Error = r.newNativeFuncConstruct(r.builtin_Error, "Error", r.global.ErrorPrototype, 1)
	r.addToGlobal("Error", r.global.Error)

	r.global.AggregateErrorPrototype = r.createErrorPrototype(stringAggregateError)
	r.global.AggregateError = r.newNativeFuncConstructProto(r.builtin_AggregateError, "AggregateError", r.global.AggregateErrorPrototype, r.global.Error, 2)
	r.addToGlobal("AggregateError", r.global.AggregateError)

	r.global.TypeErrorPrototype = r.createErrorPrototype(stringTypeError)

	r.global.TypeError = r.newNativeFuncConstructProto(r.builtin_Error, "TypeError", r.global.TypeErrorPrototype, r.global.Error, 1)
	r.addToGlobal("TypeError", r.global.TypeError)

	r.global.ReferenceErrorPrototype = r.createErrorPrototype(stringReferenceError)

	r.global.ReferenceError = r.newNativeFuncConstructProto(r.builtin_Error, "ReferenceError", r.global.ReferenceErrorPrototype, r.global.Error, 1)
	r.addToGlobal("ReferenceError", r.global.ReferenceError)

	r.global.SyntaxErrorPrototype = r.createErrorPrototype(stringSyntaxError)

	r.global.SyntaxError = r.newNativeFuncConstructProto(r.builtin_Error, "SyntaxError", r.global.SyntaxErrorPrototype, r.global.Error, 1)
	r.addToGlobal("SyntaxError", r.global.SyntaxError)

	r.global.RangeErrorPrototype = r.createErrorPrototype(stringRangeError)

	r.global.RangeError = r.newNativeFuncConstructProto(r.builtin_Error, "RangeError", r.global.RangeErrorPrototype, r.global.Error, 1)
	r.addToGlobal("RangeError", r.global.RangeError)

	r.global.EvalErrorPrototype = r.createErrorPrototype(stringEvalError)
	o = r.global.EvalErrorPrototype.self
	o._putProp("name", stringEvalError, true, false, true)

	r.global.EvalError = r.newNativeFuncConstructProto(r.builtin_Error, "EvalError", r.global.EvalErrorPrototype, r.global.Error, 1)
	r.addToGlobal("EvalError", r.global.EvalError)

	r.global.URIErrorPrototype = r.createErrorPrototype(stringURIError)

	r.global.URIError = r.newNativeFuncConstructProto(r.builtin_Error, "URIError", r.global.URIErrorPrototype, r.global.Error, 1)
	r.addToGlobal("URIError", r.global.URIError)

	r.global.GoErrorPrototype = r.createErrorPrototype(stringGoError)

	r.global.GoError = r.newNativeFuncConstructProto(r.builtin_Error, "GoError", r.global.GoErrorPrototype, r.global.Error, 1)
	r.addToGlobal("GoError", r.global.GoError)
***REMOVED***
