package goja

import (
	"math"
)

func (r *Runtime) math_abs(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Abs(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_acos(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Acos(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_asin(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Asin(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_atan(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Atan(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_atan2(call FunctionCall) Value ***REMOVED***
	y := call.Argument(0).ToFloat()
	x := call.Argument(1).ToFloat()

	return floatToValue(math.Atan2(y, x))
***REMOVED***

func (r *Runtime) math_ceil(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Ceil(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_cos(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Cos(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_exp(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Exp(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_floor(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Floor(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_log(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Log(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_max(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) == 0 ***REMOVED***
		return _negativeInf
	***REMOVED***

	result := call.Arguments[0].ToFloat()
	if math.IsNaN(result) ***REMOVED***
		return _NaN
	***REMOVED***
	for _, arg := range call.Arguments[1:] ***REMOVED***
		f := arg.ToFloat()
		if math.IsNaN(f) ***REMOVED***
			return _NaN
		***REMOVED***
		result = math.Max(result, f)
	***REMOVED***
	return floatToValue(result)
***REMOVED***

func (r *Runtime) math_min(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) == 0 ***REMOVED***
		return _positiveInf
	***REMOVED***

	result := call.Arguments[0].ToFloat()
	if math.IsNaN(result) ***REMOVED***
		return _NaN
	***REMOVED***
	for _, arg := range call.Arguments[1:] ***REMOVED***
		f := arg.ToFloat()
		if math.IsNaN(f) ***REMOVED***
			return _NaN
		***REMOVED***
		result = math.Min(result, f)
	***REMOVED***
	return floatToValue(result)
***REMOVED***

func (r *Runtime) math_pow(call FunctionCall) Value ***REMOVED***
	x := call.Argument(0)
	y := call.Argument(1)
	if x, ok := x.assertInt(); ok ***REMOVED***
		if y, ok := y.assertInt(); ok && y >= 0 && y < 64 ***REMOVED***
			if y == 0 ***REMOVED***
				return intToValue(1)
			***REMOVED***
			if x == 0 ***REMOVED***
				return intToValue(0)
			***REMOVED***
			ip := ipow(x, y)
			if ip != 0 ***REMOVED***
				return intToValue(ip)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return floatToValue(math.Pow(x.ToFloat(), y.ToFloat()))
***REMOVED***

func (r *Runtime) math_random(call FunctionCall) Value ***REMOVED***
	return floatToValue(r.rand())
***REMOVED***

func (r *Runtime) math_round(call FunctionCall) Value ***REMOVED***
	f := call.Argument(0).ToFloat()
	if math.IsNaN(f) ***REMOVED***
		return _NaN
	***REMOVED***

	if f == 0 && math.Signbit(f) ***REMOVED***
		return _negativeZero
	***REMOVED***

	t := math.Trunc(f)

	if f >= 0 ***REMOVED***
		if f-t >= 0.5 ***REMOVED***
			return floatToValue(t + 1)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if t-f > 0.5 ***REMOVED***
			return floatToValue(t - 1)
		***REMOVED***
	***REMOVED***

	return floatToValue(t)
***REMOVED***

func (r *Runtime) math_sin(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Sin(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_sqrt(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Sqrt(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_tan(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Tan(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) createMath(val *Object) objectImpl ***REMOVED***
	m := &baseObject***REMOVED***
		class:      "Math",
		val:        val,
		extensible: true,
		prototype:  r.global.ObjectPrototype,
	***REMOVED***
	m.init()

	m._putProp("E", valueFloat(math.E), false, false, false)
	m._putProp("LN10", valueFloat(math.Ln10), false, false, false)
	m._putProp("LN2", valueFloat(math.Ln2), false, false, false)
	m._putProp("LOG2E", valueFloat(math.Log2E), false, false, false)
	m._putProp("LOG10E", valueFloat(math.Log10E), false, false, false)
	m._putProp("PI", valueFloat(math.Pi), false, false, false)
	m._putProp("SQRT1_2", valueFloat(sqrt1_2), false, false, false)
	m._putProp("SQRT2", valueFloat(math.Sqrt2), false, false, false)

	m._putProp("abs", r.newNativeFunc(r.math_abs, nil, "abs", nil, 1), true, false, true)
	m._putProp("acos", r.newNativeFunc(r.math_acos, nil, "acos", nil, 1), true, false, true)
	m._putProp("asin", r.newNativeFunc(r.math_asin, nil, "asin", nil, 1), true, false, true)
	m._putProp("atan", r.newNativeFunc(r.math_atan, nil, "atan", nil, 1), true, false, true)
	m._putProp("atan2", r.newNativeFunc(r.math_atan2, nil, "atan2", nil, 2), true, false, true)
	m._putProp("ceil", r.newNativeFunc(r.math_ceil, nil, "ceil", nil, 1), true, false, true)
	m._putProp("cos", r.newNativeFunc(r.math_cos, nil, "cos", nil, 1), true, false, true)
	m._putProp("exp", r.newNativeFunc(r.math_exp, nil, "exp", nil, 1), true, false, true)
	m._putProp("floor", r.newNativeFunc(r.math_floor, nil, "floor", nil, 1), true, false, true)
	m._putProp("log", r.newNativeFunc(r.math_log, nil, "log", nil, 1), true, false, true)
	m._putProp("max", r.newNativeFunc(r.math_max, nil, "max", nil, 2), true, false, true)
	m._putProp("min", r.newNativeFunc(r.math_min, nil, "min", nil, 2), true, false, true)
	m._putProp("pow", r.newNativeFunc(r.math_pow, nil, "pow", nil, 2), true, false, true)
	m._putProp("random", r.newNativeFunc(r.math_random, nil, "random", nil, 0), true, false, true)
	m._putProp("round", r.newNativeFunc(r.math_round, nil, "round", nil, 1), true, false, true)
	m._putProp("sin", r.newNativeFunc(r.math_sin, nil, "sin", nil, 1), true, false, true)
	m._putProp("sqrt", r.newNativeFunc(r.math_sqrt, nil, "sqrt", nil, 1), true, false, true)
	m._putProp("tan", r.newNativeFunc(r.math_tan, nil, "tan", nil, 1), true, false, true)

	return m
***REMOVED***

func (r *Runtime) initMath() ***REMOVED***
	r.addToGlobal("Math", r.newLazyObject(r.createMath))
***REMOVED***
