package goja

import (
	"math"

	"github.com/dop251/goja/ftoa"
)

func (r *Runtime) numberproto_valueOf(call FunctionCall) Value ***REMOVED***
	this := call.This
	if !isNumber(this) ***REMOVED***
		r.typeErrorResult(true, "Value is not a number")
	***REMOVED***
	switch t := this.(type) ***REMOVED***
	case valueInt, valueFloat:
		return this
	case *Object:
		if v, ok := t.self.(*primitiveValueObject); ok ***REMOVED***
			return v.pValue
		***REMOVED***
	***REMOVED***

	panic(r.NewTypeError("Number.prototype.valueOf is not generic"))
***REMOVED***

func isNumber(v Value) bool ***REMOVED***
	switch t := v.(type) ***REMOVED***
	case valueFloat, valueInt:
		return true
	case *Object:
		switch t := t.self.(type) ***REMOVED***
		case *primitiveValueObject:
			return isNumber(t.pValue)
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (r *Runtime) numberproto_toString(call FunctionCall) Value ***REMOVED***
	if !isNumber(call.This) ***REMOVED***
		r.typeErrorResult(true, "Value is not a number")
	***REMOVED***
	var radix int
	if arg := call.Argument(0); arg != _undefined ***REMOVED***
		radix = int(arg.ToInteger())
	***REMOVED*** else ***REMOVED***
		radix = 10
	***REMOVED***

	if radix < 2 || radix > 36 ***REMOVED***
		panic(r.newError(r.global.RangeError, "toString() radix argument must be between 2 and 36"))
	***REMOVED***

	num := call.This.ToFloat()

	if math.IsNaN(num) ***REMOVED***
		return stringNaN
	***REMOVED***

	if math.IsInf(num, 1) ***REMOVED***
		return stringInfinity
	***REMOVED***

	if math.IsInf(num, -1) ***REMOVED***
		return stringNegInfinity
	***REMOVED***

	if radix == 10 ***REMOVED***
		return asciiString(fToStr(num, ftoa.ModeStandard, 0))
	***REMOVED***

	return asciiString(ftoa.FToBaseStr(num, radix))
***REMOVED***

func (r *Runtime) numberproto_toFixed(call FunctionCall) Value ***REMOVED***
	num := r.toNumber(call.This).ToFloat()
	prec := call.Argument(0).ToInteger()

	if prec < 0 || prec > 100 ***REMOVED***
		panic(r.newError(r.global.RangeError, "toFixed() precision must be between 0 and 100"))
	***REMOVED***
	if math.IsNaN(num) ***REMOVED***
		return stringNaN
	***REMOVED***
	return asciiString(fToStr(num, ftoa.ModeFixed, int(prec)))
***REMOVED***

func (r *Runtime) numberproto_toExponential(call FunctionCall) Value ***REMOVED***
	num := r.toNumber(call.This).ToFloat()
	precVal := call.Argument(0)
	var prec int64
	if precVal == _undefined ***REMOVED***
		return asciiString(fToStr(num, ftoa.ModeStandardExponential, 0))
	***REMOVED*** else ***REMOVED***
		prec = precVal.ToInteger()
	***REMOVED***

	if math.IsNaN(num) ***REMOVED***
		return stringNaN
	***REMOVED***
	if math.IsInf(num, 1) ***REMOVED***
		return stringInfinity
	***REMOVED***
	if math.IsInf(num, -1) ***REMOVED***
		return stringNegInfinity
	***REMOVED***

	if prec < 0 || prec > 100 ***REMOVED***
		panic(r.newError(r.global.RangeError, "toExponential() precision must be between 0 and 100"))
	***REMOVED***

	return asciiString(fToStr(num, ftoa.ModeExponential, int(prec+1)))
***REMOVED***

func (r *Runtime) numberproto_toPrecision(call FunctionCall) Value ***REMOVED***
	numVal := r.toNumber(call.This)
	precVal := call.Argument(0)
	if precVal == _undefined ***REMOVED***
		return numVal.toString()
	***REMOVED***
	num := numVal.ToFloat()
	prec := precVal.ToInteger()

	if math.IsNaN(num) ***REMOVED***
		return stringNaN
	***REMOVED***
	if math.IsInf(num, 1) ***REMOVED***
		return stringInfinity
	***REMOVED***
	if math.IsInf(num, -1) ***REMOVED***
		return stringNegInfinity
	***REMOVED***
	if prec < 1 || prec > 100 ***REMOVED***
		panic(r.newError(r.global.RangeError, "toPrecision() precision must be between 1 and 100"))
	***REMOVED***

	return asciiString(fToStr(num, ftoa.ModePrecision, int(prec)))
***REMOVED***

func (r *Runtime) number_isFinite(call FunctionCall) Value ***REMOVED***
	switch arg := call.Argument(0).(type) ***REMOVED***
	case valueInt:
		return valueTrue
	case valueFloat:
		f := float64(arg)
		return r.toBoolean(!math.IsInf(f, 0) && !math.IsNaN(f))
	default:
		return valueFalse
	***REMOVED***
***REMOVED***

func (r *Runtime) number_isInteger(call FunctionCall) Value ***REMOVED***
	switch arg := call.Argument(0).(type) ***REMOVED***
	case valueInt:
		return valueTrue
	case valueFloat:
		f := float64(arg)
		return r.toBoolean(!math.IsNaN(f) && !math.IsInf(f, 0) && math.Floor(f) == f)
	default:
		return valueFalse
	***REMOVED***
***REMOVED***

func (r *Runtime) number_isNaN(call FunctionCall) Value ***REMOVED***
	if f, ok := call.Argument(0).(valueFloat); ok && math.IsNaN(float64(f)) ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) number_isSafeInteger(call FunctionCall) Value ***REMOVED***
	arg := call.Argument(0)
	if i, ok := arg.(valueInt); ok && i >= -(maxInt-1) && i <= maxInt-1 ***REMOVED***
		return valueTrue
	***REMOVED***
	if arg == _negativeZero ***REMOVED***
		return valueTrue
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) initNumber() ***REMOVED***
	r.global.NumberPrototype = r.newPrimitiveObject(valueInt(0), r.global.ObjectPrototype, classNumber)
	o := r.global.NumberPrototype.self
	o._putProp("toExponential", r.newNativeFunc(r.numberproto_toExponential, nil, "toExponential", nil, 1), true, false, true)
	o._putProp("toFixed", r.newNativeFunc(r.numberproto_toFixed, nil, "toFixed", nil, 1), true, false, true)
	o._putProp("toLocaleString", r.newNativeFunc(r.numberproto_toString, nil, "toLocaleString", nil, 0), true, false, true)
	o._putProp("toPrecision", r.newNativeFunc(r.numberproto_toPrecision, nil, "toPrecision", nil, 1), true, false, true)
	o._putProp("toString", r.newNativeFunc(r.numberproto_toString, nil, "toString", nil, 1), true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.numberproto_valueOf, nil, "valueOf", nil, 0), true, false, true)

	r.global.Number = r.newNativeFunc(r.builtin_Number, r.builtin_newNumber, "Number", r.global.NumberPrototype, 1)
	o = r.global.Number.self
	o._putProp("EPSILON", _epsilon, false, false, false)
	o._putProp("isFinite", r.newNativeFunc(r.number_isFinite, nil, "isFinite", nil, 1), true, false, true)
	o._putProp("isInteger", r.newNativeFunc(r.number_isInteger, nil, "isInteger", nil, 1), true, false, true)
	o._putProp("isNaN", r.newNativeFunc(r.number_isNaN, nil, "isNaN", nil, 1), true, false, true)
	o._putProp("isSafeInteger", r.newNativeFunc(r.number_isSafeInteger, nil, "isSafeInteger", nil, 1), true, false, true)
	o._putProp("MAX_SAFE_INTEGER", valueInt(maxInt-1), false, false, false)
	o._putProp("MIN_SAFE_INTEGER", valueInt(-(maxInt - 1)), false, false, false)
	o._putProp("MIN_VALUE", valueFloat(math.SmallestNonzeroFloat64), false, false, false)
	o._putProp("MAX_VALUE", valueFloat(math.MaxFloat64), false, false, false)
	o._putProp("NaN", _NaN, false, false, false)
	o._putProp("NEGATIVE_INFINITY", _negativeInf, false, false, false)
	o._putProp("parseFloat", r.Get("parseFloat"), true, false, true)
	o._putProp("parseInt", r.Get("parseInt"), true, false, true)
	o._putProp("POSITIVE_INFINITY", _positiveInf, false, false, false)
	r.addToGlobal("Number", r.global.Number)

***REMOVED***
