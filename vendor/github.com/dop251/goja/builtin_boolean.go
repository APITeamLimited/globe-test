package goja

func (r *Runtime) booleanproto_toString(call FunctionCall) Value ***REMOVED***
	var b bool
	switch o := call.This.(type) ***REMOVED***
	case valueBool:
		b = bool(o)
		goto success
	case *Object:
		if p, ok := o.self.(*primitiveValueObject); ok ***REMOVED***
			if b1, ok := p.pValue.(valueBool); ok ***REMOVED***
				b = bool(b1)
				goto success
			***REMOVED***
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Method Boolean.prototype.toString is called on incompatible receiver")

success:
	if b ***REMOVED***
		return stringTrue
	***REMOVED***
	return stringFalse
***REMOVED***

func (r *Runtime) booleanproto_valueOf(call FunctionCall) Value ***REMOVED***
	switch o := call.This.(type) ***REMOVED***
	case valueBool:
		return o
	case *Object:
		if p, ok := o.self.(*primitiveValueObject); ok ***REMOVED***
			if b, ok := p.pValue.(valueBool); ok ***REMOVED***
				return b
			***REMOVED***
		***REMOVED***
	***REMOVED***

	r.typeErrorResult(true, "Method Boolean.prototype.valueOf is called on incompatible receiver")
	return nil
***REMOVED***

func (r *Runtime) initBoolean() ***REMOVED***
	r.global.BooleanPrototype = r.newPrimitiveObject(valueFalse, r.global.ObjectPrototype, classBoolean)
	o := r.global.BooleanPrototype.self
	o._putProp("toString", r.newNativeFunc(r.booleanproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.booleanproto_valueOf, nil, "valueOf", nil, 0), true, false, true)

	r.global.Boolean = r.newNativeFunc(r.builtin_Boolean, r.builtin_newBoolean, "Boolean", r.global.BooleanPrototype, 1)
	r.addToGlobal("Boolean", r.global.Boolean)
***REMOVED***
