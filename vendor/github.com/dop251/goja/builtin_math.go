package goja

import (
	"math"
	"math/bits"
)

func (r *Runtime) math_abs(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Abs(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_acos(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Acos(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_acosh(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Acosh(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_asin(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Asin(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_asinh(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Asinh(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_atan(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Atan(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_atanh(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Atanh(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_atan2(call FunctionCall) Value ***REMOVED***
	y := call.Argument(0).ToFloat()
	x := call.Argument(1).ToFloat()

	return floatToValue(math.Atan2(y, x))
***REMOVED***

func (r *Runtime) math_cbrt(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Cbrt(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_ceil(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Ceil(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_clz32(call FunctionCall) Value ***REMOVED***
	return intToValue(int64(bits.LeadingZeros32(toUint32(call.Argument(0)))))
***REMOVED***

func (r *Runtime) math_cos(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Cos(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_cosh(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Cosh(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_exp(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Exp(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_expm1(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Expm1(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_floor(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Floor(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_fround(call FunctionCall) Value ***REMOVED***
	return floatToValue(float64(float32(call.Argument(0).ToFloat())))
***REMOVED***

func (r *Runtime) math_hypot(call FunctionCall) Value ***REMOVED***
	var max float64
	var hasNaN bool
	absValues := make([]float64, 0, len(call.Arguments))
	for _, v := range call.Arguments ***REMOVED***
		arg := nilSafe(v).ToFloat()
		if math.IsNaN(arg) ***REMOVED***
			hasNaN = true
		***REMOVED*** else ***REMOVED***
			abs := math.Abs(arg)
			if abs > max ***REMOVED***
				max = abs
			***REMOVED***
			absValues = append(absValues, abs)
		***REMOVED***
	***REMOVED***
	if math.IsInf(max, 1) ***REMOVED***
		return _positiveInf
	***REMOVED***
	if hasNaN ***REMOVED***
		return _NaN
	***REMOVED***
	if max == 0 ***REMOVED***
		return _positiveZero
	***REMOVED***

	// Kahan summation to avoid rounding errors.
	// Normalize the numbers to the largest one to avoid overflow.
	var sum, compensation float64
	for _, n := range absValues ***REMOVED***
		n /= max
		summand := n*n - compensation
		preliminary := sum + summand
		compensation = (preliminary - sum) - summand
		sum = preliminary
	***REMOVED***
	return floatToValue(math.Sqrt(sum) * max)
***REMOVED***

func (r *Runtime) math_imul(call FunctionCall) Value ***REMOVED***
	x := toUint32(call.Argument(0))
	y := toUint32(call.Argument(1))
	return intToValue(int64(int32(x * y)))
***REMOVED***

func (r *Runtime) math_log(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Log(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_log1p(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Log1p(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_log10(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Log10(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_log2(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Log2(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_max(call FunctionCall) Value ***REMOVED***
	result := math.Inf(-1)
	args := call.Arguments
	for i, arg := range args ***REMOVED***
		n := nilSafe(arg).ToFloat()
		if math.IsNaN(n) ***REMOVED***
			args = args[i+1:]
			goto NaNLoop
		***REMOVED***
		result = math.Max(result, n)
	***REMOVED***

	return floatToValue(result)

NaNLoop:
	// All arguments still need to be coerced to number according to the specs.
	for _, arg := range args ***REMOVED***
		nilSafe(arg).ToFloat()
	***REMOVED***
	return _NaN
***REMOVED***

func (r *Runtime) math_min(call FunctionCall) Value ***REMOVED***
	result := math.Inf(1)
	args := call.Arguments
	for i, arg := range args ***REMOVED***
		n := nilSafe(arg).ToFloat()
		if math.IsNaN(n) ***REMOVED***
			args = args[i+1:]
			goto NaNLoop
		***REMOVED***
		result = math.Min(result, n)
	***REMOVED***

	return floatToValue(result)

NaNLoop:
	// All arguments still need to be coerced to number according to the specs.
	for _, arg := range args ***REMOVED***
		nilSafe(arg).ToFloat()
	***REMOVED***
	return _NaN
***REMOVED***

func (r *Runtime) math_pow(call FunctionCall) Value ***REMOVED***
	x := call.Argument(0)
	y := call.Argument(1)
	if x, ok := x.(valueInt); ok ***REMOVED***
		if y, ok := y.(valueInt); ok && y >= 0 && y < 64 ***REMOVED***
			if y == 0 ***REMOVED***
				return intToValue(1)
			***REMOVED***
			if x == 0 ***REMOVED***
				return intToValue(0)
			***REMOVED***
			ip := ipow(int64(x), int64(y))
			if ip != 0 ***REMOVED***
				return intToValue(ip)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	xf := x.ToFloat()
	yf := y.ToFloat()
	if math.Abs(xf) == 1 && math.IsInf(yf, 0) ***REMOVED***
		return _NaN
	***REMOVED***
	if xf == 1 && math.IsNaN(yf) ***REMOVED***
		return _NaN
	***REMOVED***
	return floatToValue(math.Pow(xf, yf))
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

func (r *Runtime) math_sign(call FunctionCall) Value ***REMOVED***
	arg := call.Argument(0)
	num := arg.ToFloat()
	if math.IsNaN(num) || num == 0 ***REMOVED*** // this will match -0 too
		return arg
	***REMOVED***
	if num > 0 ***REMOVED***
		return intToValue(1)
	***REMOVED***
	return intToValue(-1)
***REMOVED***

func (r *Runtime) math_sin(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Sin(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_sinh(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Sinh(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_sqrt(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Sqrt(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_tan(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Tan(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_tanh(call FunctionCall) Value ***REMOVED***
	return floatToValue(math.Tanh(call.Argument(0).ToFloat()))
***REMOVED***

func (r *Runtime) math_trunc(call FunctionCall) Value ***REMOVED***
	arg := call.Argument(0)
	if i, ok := arg.(valueInt); ok ***REMOVED***
		return i
	***REMOVED***
	return floatToValue(math.Trunc(arg.ToFloat()))
***REMOVED***

func (r *Runtime) createMath(val *Object) objectImpl ***REMOVED***
	m := &baseObject***REMOVED***
		class:      classMath,
		val:        val,
		extensible: true,
		prototype:  r.global.ObjectPrototype,
	***REMOVED***
	m.init()

	m._putProp("E", valueFloat(math.E), false, false, false)
	m._putProp("LN10", valueFloat(math.Ln10), false, false, false)
	m._putProp("LN2", valueFloat(math.Ln2), false, false, false)
	m._putProp("LOG10E", valueFloat(math.Log10E), false, false, false)
	m._putProp("LOG2E", valueFloat(math.Log2E), false, false, false)
	m._putProp("PI", valueFloat(math.Pi), false, false, false)
	m._putProp("SQRT1_2", valueFloat(sqrt1_2), false, false, false)
	m._putProp("SQRT2", valueFloat(math.Sqrt2), false, false, false)
	m._putSym(SymToStringTag, valueProp(asciiString(classMath), false, false, true))

	m._putProp("abs", r.newNativeFunc(r.math_abs, nil, "abs", nil, 1), true, false, true)
	m._putProp("acos", r.newNativeFunc(r.math_acos, nil, "acos", nil, 1), true, false, true)
	m._putProp("acosh", r.newNativeFunc(r.math_acosh, nil, "acosh", nil, 1), true, false, true)
	m._putProp("asin", r.newNativeFunc(r.math_asin, nil, "asin", nil, 1), true, false, true)
	m._putProp("asinh", r.newNativeFunc(r.math_asinh, nil, "asinh", nil, 1), true, false, true)
	m._putProp("atan", r.newNativeFunc(r.math_atan, nil, "atan", nil, 1), true, false, true)
	m._putProp("atanh", r.newNativeFunc(r.math_atanh, nil, "atanh", nil, 1), true, false, true)
	m._putProp("atan2", r.newNativeFunc(r.math_atan2, nil, "atan2", nil, 2), true, false, true)
	m._putProp("cbrt", r.newNativeFunc(r.math_cbrt, nil, "cbrt", nil, 1), true, false, true)
	m._putProp("ceil", r.newNativeFunc(r.math_ceil, nil, "ceil", nil, 1), true, false, true)
	m._putProp("clz32", r.newNativeFunc(r.math_clz32, nil, "clz32", nil, 1), true, false, true)
	m._putProp("cos", r.newNativeFunc(r.math_cos, nil, "cos", nil, 1), true, false, true)
	m._putProp("cosh", r.newNativeFunc(r.math_cosh, nil, "cosh", nil, 1), true, false, true)
	m._putProp("exp", r.newNativeFunc(r.math_exp, nil, "exp", nil, 1), true, false, true)
	m._putProp("expm1", r.newNativeFunc(r.math_expm1, nil, "expm1", nil, 1), true, false, true)
	m._putProp("floor", r.newNativeFunc(r.math_floor, nil, "floor", nil, 1), true, false, true)
	m._putProp("fround", r.newNativeFunc(r.math_fround, nil, "fround", nil, 1), true, false, true)
	m._putProp("hypot", r.newNativeFunc(r.math_hypot, nil, "hypot", nil, 2), true, false, true)
	m._putProp("imul", r.newNativeFunc(r.math_imul, nil, "imul", nil, 2), true, false, true)
	m._putProp("log", r.newNativeFunc(r.math_log, nil, "log", nil, 1), true, false, true)
	m._putProp("log1p", r.newNativeFunc(r.math_log1p, nil, "log1p", nil, 1), true, false, true)
	m._putProp("log10", r.newNativeFunc(r.math_log10, nil, "log10", nil, 1), true, false, true)
	m._putProp("log2", r.newNativeFunc(r.math_log2, nil, "log2", nil, 1), true, false, true)
	m._putProp("max", r.newNativeFunc(r.math_max, nil, "max", nil, 2), true, false, true)
	m._putProp("min", r.newNativeFunc(r.math_min, nil, "min", nil, 2), true, false, true)
	m._putProp("pow", r.newNativeFunc(r.math_pow, nil, "pow", nil, 2), true, false, true)
	m._putProp("random", r.newNativeFunc(r.math_random, nil, "random", nil, 0), true, false, true)
	m._putProp("round", r.newNativeFunc(r.math_round, nil, "round", nil, 1), true, false, true)
	m._putProp("sign", r.newNativeFunc(r.math_sign, nil, "sign", nil, 1), true, false, true)
	m._putProp("sin", r.newNativeFunc(r.math_sin, nil, "sin", nil, 1), true, false, true)
	m._putProp("sinh", r.newNativeFunc(r.math_sinh, nil, "sinh", nil, 1), true, false, true)
	m._putProp("sqrt", r.newNativeFunc(r.math_sqrt, nil, "sqrt", nil, 1), true, false, true)
	m._putProp("tan", r.newNativeFunc(r.math_tan, nil, "tan", nil, 1), true, false, true)
	m._putProp("tanh", r.newNativeFunc(r.math_tanh, nil, "tanh", nil, 1), true, false, true)
	m._putProp("trunc", r.newNativeFunc(r.math_trunc, nil, "trunc", nil, 1), true, false, true)

	return m
***REMOVED***

func (r *Runtime) initMath() ***REMOVED***
	r.addToGlobal("Math", r.newLazyObject(r.createMath))
***REMOVED***
