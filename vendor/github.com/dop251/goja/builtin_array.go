package goja

import (
	"math"
	"sort"
)

func (r *Runtime) newArray(prototype *Object) (a *arrayObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	a = &arrayObject***REMOVED******REMOVED***
	a.class = classArray
	a.val = v
	a.extensible = true
	v.self = a
	a.prototype = prototype
	a.init()
	return
***REMOVED***

func (r *Runtime) newArrayObject() *arrayObject ***REMOVED***
	return r.newArray(r.global.ArrayPrototype)
***REMOVED***

func setArrayValues(a *arrayObject, values []Value) *arrayObject ***REMOVED***
	a.values = values
	a.length = uint32(len(values))
	a.objCount = len(values)
	return a
***REMOVED***

func setArrayLength(a *arrayObject, l int64) *arrayObject ***REMOVED***
	a.setOwnStr("length", intToValue(l), true)
	return a
***REMOVED***

func arraySpeciesCreate(obj *Object, size int64) *Object ***REMOVED***
	if isArray(obj) ***REMOVED***
		v := obj.self.getStr("constructor", nil)
		if constructObj, ok := v.(*Object); ok ***REMOVED***
			v = constructObj.self.getSym(SymSpecies, nil)
			if v == _null ***REMOVED***
				v = nil
			***REMOVED***
		***REMOVED***

		if v != nil && v != _undefined ***REMOVED***
			constructObj, _ := v.(*Object)
			if constructObj != nil ***REMOVED***
				if constructor := constructObj.self.assertConstructor(); constructor != nil ***REMOVED***
					return constructor([]Value***REMOVED***intToValue(size)***REMOVED***, constructObj)
				***REMOVED***
			***REMOVED***
			panic(obj.runtime.NewTypeError("Species is not a constructor"))
		***REMOVED***
	***REMOVED***
	return obj.runtime.newArrayLength(size)
***REMOVED***

func max(a, b int64) int64 ***REMOVED***
	if a > b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func min(a, b int64) int64 ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func relToIdx(rel, l int64) int64 ***REMOVED***
	if rel >= 0 ***REMOVED***
		return min(rel, l)
	***REMOVED***
	return max(l+rel, 0)
***REMOVED***

func (r *Runtime) newArrayValues(values []Value) *Object ***REMOVED***
	return setArrayValues(r.newArrayObject(), values).val
***REMOVED***

func (r *Runtime) newArrayLength(l int64) *Object ***REMOVED***
	return setArrayLength(r.newArrayObject(), l).val
***REMOVED***

func (r *Runtime) builtin_newArray(args []Value, proto *Object) *Object ***REMOVED***
	l := len(args)
	if l == 1 ***REMOVED***
		if al, ok := args[0].(valueInt); ok ***REMOVED***
			return setArrayLength(r.newArray(proto), int64(al)).val
		***REMOVED*** else if f, ok := args[0].(valueFloat); ok ***REMOVED***
			al := int64(f)
			if float64(al) == float64(f) ***REMOVED***
				return r.newArrayLength(al)
			***REMOVED*** else ***REMOVED***
				panic(r.newError(r.global.RangeError, "Invalid array length"))
			***REMOVED***
		***REMOVED***
		return setArrayValues(r.newArray(proto), []Value***REMOVED***args[0]***REMOVED***).val
	***REMOVED*** else ***REMOVED***
		argsCopy := make([]Value, l)
		copy(argsCopy, args)
		return setArrayValues(r.newArray(proto), argsCopy).val
	***REMOVED***
***REMOVED***

func (r *Runtime) generic_push(obj *Object, call FunctionCall) Value ***REMOVED***
	l := toLength(obj.self.getStr("length", nil))
	nl := l + int64(len(call.Arguments))
	if nl >= maxInt ***REMOVED***
		r.typeErrorResult(true, "Invalid array length")
		panic("unreachable")
	***REMOVED***
	for i, arg := range call.Arguments ***REMOVED***
		obj.self.setOwnIdx(valueInt(l+int64(i)), arg, true)
	***REMOVED***
	n := valueInt(nl)
	obj.self.setOwnStr("length", n, true)
	return n
***REMOVED***

func (r *Runtime) arrayproto_push(call FunctionCall) Value ***REMOVED***
	obj := call.This.ToObject(r)
	return r.generic_push(obj, call)
***REMOVED***

func (r *Runtime) arrayproto_pop_generic(obj *Object) Value ***REMOVED***
	l := toLength(obj.self.getStr("length", nil))
	if l == 0 ***REMOVED***
		obj.self.setOwnStr("length", intToValue(0), true)
		return _undefined
	***REMOVED***
	idx := valueInt(l - 1)
	val := obj.self.getIdx(idx, nil)
	obj.self.deleteIdx(idx, true)
	obj.self.setOwnStr("length", idx, true)
	return val
***REMOVED***

func (r *Runtime) arrayproto_pop(call FunctionCall) Value ***REMOVED***
	obj := call.This.ToObject(r)
	if a, ok := obj.self.(*arrayObject); ok ***REMOVED***
		l := a.length
		if l > 0 ***REMOVED***
			var val Value
			l--
			if l < uint32(len(a.values)) ***REMOVED***
				val = a.values[l]
			***REMOVED***
			if val == nil ***REMOVED***
				// optimisation bail-out
				return r.arrayproto_pop_generic(obj)
			***REMOVED***
			if _, ok := val.(*valueProperty); ok ***REMOVED***
				// optimisation bail-out
				return r.arrayproto_pop_generic(obj)
			***REMOVED***
			//a._setLengthInt(l, false)
			a.values[l] = nil
			a.values = a.values[:l]
			a.length = l
			return val
		***REMOVED***
		return _undefined
	***REMOVED*** else ***REMOVED***
		return r.arrayproto_pop_generic(obj)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_join(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := int(toLength(o.self.getStr("length", nil)))
	var sep valueString
	if s := call.Argument(0); s != _undefined ***REMOVED***
		sep = s.toString()
	***REMOVED*** else ***REMOVED***
		sep = asciiString(",")
	***REMOVED***
	if l == 0 ***REMOVED***
		return stringEmpty
	***REMOVED***

	var buf valueStringBuilder

	element0 := o.self.getIdx(valueInt(0), nil)
	if element0 != nil && element0 != _undefined && element0 != _null ***REMOVED***
		buf.WriteString(element0.toString())
	***REMOVED***

	for i := 1; i < l; i++ ***REMOVED***
		buf.WriteString(sep)
		element := o.self.getIdx(valueInt(int64(i)), nil)
		if element != nil && element != _undefined && element != _null ***REMOVED***
			buf.WriteString(element.toString())
		***REMOVED***
	***REMOVED***

	return buf.String()
***REMOVED***

func (r *Runtime) arrayproto_toString(call FunctionCall) Value ***REMOVED***
	array := call.This.ToObject(r)
	f := array.self.getStr("join", nil)
	if fObj, ok := f.(*Object); ok ***REMOVED***
		if fcall, ok := fObj.self.assertCallable(); ok ***REMOVED***
			return fcall(FunctionCall***REMOVED***
				This: array,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return r.objectproto_toString(FunctionCall***REMOVED***
		This: array,
	***REMOVED***)
***REMOVED***

func (r *Runtime) writeItemLocaleString(item Value, buf *valueStringBuilder) ***REMOVED***
	if item != nil && item != _undefined && item != _null ***REMOVED***
		if f, ok := r.getVStr(item, "toLocaleString").(*Object); ok ***REMOVED***
			if c, ok := f.self.assertCallable(); ok ***REMOVED***
				strVal := c(FunctionCall***REMOVED***
					This: item,
				***REMOVED***)
				buf.WriteString(strVal.toString())
				return
			***REMOVED***
		***REMOVED***
		r.typeErrorResult(true, "Property 'toLocaleString' of object %s is not a function", item)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_toLocaleString(call FunctionCall) Value ***REMOVED***
	array := call.This.ToObject(r)
	var buf valueStringBuilder
	if a := r.checkStdArrayObj(array); a != nil ***REMOVED***
		for i, item := range a.values ***REMOVED***
			if i > 0 ***REMOVED***
				buf.WriteRune(',')
			***REMOVED***
			r.writeItemLocaleString(item, &buf)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		length := toLength(array.self.getStr("length", nil))
		for i := int64(0); i < length; i++ ***REMOVED***
			if i > 0 ***REMOVED***
				buf.WriteRune(',')
			***REMOVED***
			item := array.self.getIdx(valueInt(i), nil)
			r.writeItemLocaleString(item, &buf)
		***REMOVED***
	***REMOVED***

	return buf.String()
***REMOVED***

func isConcatSpreadable(obj *Object) bool ***REMOVED***
	spreadable := obj.self.getSym(SymIsConcatSpreadable, nil)
	if spreadable != nil && spreadable != _undefined ***REMOVED***
		return spreadable.ToBoolean()
	***REMOVED***
	return isArray(obj)
***REMOVED***

func (r *Runtime) arrayproto_concat_append(a *Object, item Value) ***REMOVED***
	aLength := toLength(a.self.getStr("length", nil))
	if obj, ok := item.(*Object); ok && isConcatSpreadable(obj) ***REMOVED***
		length := toLength(obj.self.getStr("length", nil))
		if aLength+length >= maxInt ***REMOVED***
			panic(r.NewTypeError("Invalid array length"))
		***REMOVED***
		for i := int64(0); i < length; i++ ***REMOVED***
			v := obj.self.getIdx(valueInt(i), nil)
			if v != nil ***REMOVED***
				createDataPropertyOrThrow(a, intToValue(aLength), v)
			***REMOVED***
			aLength++
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		createDataPropertyOrThrow(a, intToValue(aLength), item)
		aLength++
	***REMOVED***
	a.self.setOwnStr("length", intToValue(aLength), true)
***REMOVED***

func (r *Runtime) arrayproto_concat(call FunctionCall) Value ***REMOVED***
	obj := call.This.ToObject(r)
	a := arraySpeciesCreate(obj, 0)
	r.arrayproto_concat_append(a, call.This.ToObject(r))
	for _, item := range call.Arguments ***REMOVED***
		r.arrayproto_concat_append(a, item)
	***REMOVED***
	return a
***REMOVED***

func (r *Runtime) arrayproto_slice(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	start := relToIdx(call.Argument(0).ToInteger(), length)
	var end int64
	if endArg := call.Argument(1); endArg != _undefined ***REMOVED***
		end = endArg.ToInteger()
	***REMOVED*** else ***REMOVED***
		end = length
	***REMOVED***
	end = relToIdx(end, length)

	count := end - start
	if count < 0 ***REMOVED***
		count = 0
	***REMOVED***

	a := arraySpeciesCreate(o, count)
	if src := r.checkStdArrayObj(o); src != nil ***REMOVED***
		if dst, ok := a.self.(*arrayObject); ok ***REMOVED***
			values := make([]Value, count)
			copy(values, src.values[start:])
			setArrayValues(dst, values)
			return a
		***REMOVED***
	***REMOVED***

	n := int64(0)
	for start < end ***REMOVED***
		p := o.self.getIdx(valueInt(start), nil)
		if p != nil ***REMOVED***
			createDataPropertyOrThrow(a, valueInt(n), p)
		***REMOVED***
		start++
		n++
	***REMOVED***
	return a
***REMOVED***

func (r *Runtime) arrayproto_sort(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)

	var compareFn func(FunctionCall) Value
	arg := call.Argument(0)
	if arg != _undefined ***REMOVED***
		if arg, ok := call.Argument(0).(*Object); ok ***REMOVED***
			compareFn, _ = arg.self.assertCallable()
		***REMOVED***
		if compareFn == nil ***REMOVED***
			panic(r.NewTypeError("The comparison function must be either a function or undefined"))
		***REMOVED***
	***REMOVED***

	if r.checkStdArrayObj(o) != nil ***REMOVED***
		ctx := arraySortCtx***REMOVED***
			obj:     o.self,
			compare: compareFn,
		***REMOVED***

		sort.Stable(&ctx)
	***REMOVED*** else ***REMOVED***
		length := toLength(o.self.getStr("length", nil))
		a := make([]Value, 0, length)
		for i := int64(0); i < length; i++ ***REMOVED***
			idx := valueInt(i)
			if o.self.hasPropertyIdx(idx) ***REMOVED***
				a = append(a, nilSafe(o.self.getIdx(idx, nil)))
			***REMOVED***
		***REMOVED***
		ar := r.newArrayValues(a)
		ctx := arraySortCtx***REMOVED***
			obj:     ar.self,
			compare: compareFn,
		***REMOVED***

		sort.Stable(&ctx)
		for i := 0; i < len(a); i++ ***REMOVED***
			o.self.setOwnIdx(valueInt(i), a[i], true)
		***REMOVED***
		for i := int64(len(a)); i < length; i++ ***REMOVED***
			o.self.deleteIdx(valueInt(i), true)
		***REMOVED***
	***REMOVED***
	return o
***REMOVED***

func (r *Runtime) arrayproto_splice(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	actualStart := relToIdx(call.Argument(0).ToInteger(), length)
	var actualDeleteCount int64
	switch len(call.Arguments) ***REMOVED***
	case 0:
	case 1:
		actualDeleteCount = length - actualStart
	default:
		actualDeleteCount = min(max(call.Argument(1).ToInteger(), 0), length-actualStart)
	***REMOVED***
	a := arraySpeciesCreate(o, actualDeleteCount)
	itemCount := max(int64(len(call.Arguments)-2), 0)
	newLength := length - actualDeleteCount + itemCount
	if src := r.checkStdArrayObj(o); src != nil ***REMOVED***
		if dst, ok := a.self.(*arrayObject); ok ***REMOVED***
			values := make([]Value, actualDeleteCount)
			copy(values, src.values[actualStart:])
			setArrayValues(dst, values)
		***REMOVED*** else ***REMOVED***
			for k := int64(0); k < actualDeleteCount; k++ ***REMOVED***
				createDataPropertyOrThrow(a, intToValue(k), src.values[k+actualStart])
			***REMOVED***
			a.self.setOwnStr("length", intToValue(actualDeleteCount), true)
		***REMOVED***
		var values []Value
		if itemCount < actualDeleteCount ***REMOVED***
			values = src.values
			copy(values[actualStart+itemCount:], values[actualStart+actualDeleteCount:])
			tail := values[newLength:]
			for k := range tail ***REMOVED***
				tail[k] = nil
			***REMOVED***
			values = values[:newLength]
		***REMOVED*** else if itemCount > actualDeleteCount ***REMOVED***
			if int64(cap(src.values)) >= newLength ***REMOVED***
				values = src.values[:newLength]
				copy(values[actualStart+itemCount:], values[actualStart+actualDeleteCount:length])
			***REMOVED*** else ***REMOVED***
				values = make([]Value, newLength)
				copy(values, src.values[:actualStart])
				copy(values[actualStart+itemCount:], src.values[actualStart+actualDeleteCount:])
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			values = src.values
		***REMOVED***
		if itemCount > 0 ***REMOVED***
			copy(values[actualStart:], call.Arguments[2:])
		***REMOVED***
		src.values = values
		src.objCount = len(values)
	***REMOVED*** else ***REMOVED***
		for k := int64(0); k < actualDeleteCount; k++ ***REMOVED***
			from := valueInt(k + actualStart)
			if o.self.hasPropertyIdx(from) ***REMOVED***
				createDataPropertyOrThrow(a, valueInt(k), nilSafe(o.self.getIdx(from, nil)))
			***REMOVED***
		***REMOVED***

		if itemCount < actualDeleteCount ***REMOVED***
			for k := actualStart; k < length-actualDeleteCount; k++ ***REMOVED***
				from := valueInt(k + actualDeleteCount)
				to := valueInt(k + itemCount)
				if o.self.hasPropertyIdx(from) ***REMOVED***
					o.self.setOwnIdx(to, nilSafe(o.self.getIdx(from, nil)), true)
				***REMOVED*** else ***REMOVED***
					o.self.deleteIdx(to, true)
				***REMOVED***
			***REMOVED***

			for k := length; k > length-actualDeleteCount+itemCount; k-- ***REMOVED***
				o.self.deleteIdx(valueInt(k-1), true)
			***REMOVED***
		***REMOVED*** else if itemCount > actualDeleteCount ***REMOVED***
			for k := length - actualDeleteCount; k > actualStart; k-- ***REMOVED***
				from := valueInt(k + actualDeleteCount - 1)
				to := valueInt(k + itemCount - 1)
				if o.self.hasPropertyIdx(from) ***REMOVED***
					o.self.setOwnIdx(to, nilSafe(o.self.getIdx(from, nil)), true)
				***REMOVED*** else ***REMOVED***
					o.self.deleteIdx(to, true)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if itemCount > 0 ***REMOVED***
			for i, item := range call.Arguments[2:] ***REMOVED***
				o.self.setOwnIdx(valueInt(actualStart+int64(i)), item, true)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	o.self.setOwnStr("length", intToValue(newLength), true)

	return a
***REMOVED***

func (r *Runtime) arrayproto_unshift(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	argCount := int64(len(call.Arguments))
	newLen := intToValue(length + argCount)
	newSize := length + argCount
	if arr := r.checkStdArrayObj(o); arr != nil && newSize < math.MaxUint32 ***REMOVED***
		if int64(cap(arr.values)) >= newSize ***REMOVED***
			arr.values = arr.values[:newSize]
			copy(arr.values[argCount:], arr.values[:length])
		***REMOVED*** else ***REMOVED***
			values := make([]Value, newSize)
			copy(values[argCount:], arr.values)
			arr.values = values
		***REMOVED***
		copy(arr.values, call.Arguments)
		arr.objCount = int(arr.length)
	***REMOVED*** else ***REMOVED***
		for k := length - 1; k >= 0; k-- ***REMOVED***
			from := valueInt(k)
			to := valueInt(k + argCount)
			if o.self.hasPropertyIdx(from) ***REMOVED***
				o.self.setOwnIdx(to, nilSafe(o.self.getIdx(from, nil)), true)
			***REMOVED*** else ***REMOVED***
				o.self.deleteIdx(to, true)
			***REMOVED***
		***REMOVED***

		for k, arg := range call.Arguments ***REMOVED***
			o.self.setOwnIdx(valueInt(int64(k)), arg, true)
		***REMOVED***
	***REMOVED***

	o.self.setOwnStr("length", newLen, true)
	return newLen
***REMOVED***

func (r *Runtime) arrayproto_indexOf(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	if length == 0 ***REMOVED***
		return intToValue(-1)
	***REMOVED***

	n := call.Argument(1).ToInteger()
	if n >= length ***REMOVED***
		return intToValue(-1)
	***REMOVED***

	if n < 0 ***REMOVED***
		n = max(length+n, 0)
	***REMOVED***

	searchElement := call.Argument(0)

	if arr := r.checkStdArrayObj(o); arr != nil ***REMOVED***
		for i, val := range arr.values[n:] ***REMOVED***
			if searchElement.StrictEquals(val) ***REMOVED***
				return intToValue(n + int64(i))
			***REMOVED***
		***REMOVED***
		return intToValue(-1)
	***REMOVED***

	for ; n < length; n++ ***REMOVED***
		idx := valueInt(n)
		if o.self.hasPropertyIdx(idx) ***REMOVED***
			if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
				if searchElement.StrictEquals(val) ***REMOVED***
					return idx
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return intToValue(-1)
***REMOVED***

func (r *Runtime) arrayproto_includes(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	if length == 0 ***REMOVED***
		return valueFalse
	***REMOVED***

	n := call.Argument(1).ToInteger()
	if n >= length ***REMOVED***
		return valueFalse
	***REMOVED***

	if n < 0 ***REMOVED***
		n = max(length+n, 0)
	***REMOVED***

	searchElement := call.Argument(0)
	if searchElement == _negativeZero ***REMOVED***
		searchElement = _positiveZero
	***REMOVED***

	if arr := r.checkStdArrayObj(o); arr != nil ***REMOVED***
		for _, val := range arr.values[n:] ***REMOVED***
			if searchElement.SameAs(val) ***REMOVED***
				return valueTrue
			***REMOVED***
		***REMOVED***
		return valueFalse
	***REMOVED***

	for ; n < length; n++ ***REMOVED***
		idx := valueInt(n)
		val := nilSafe(o.self.getIdx(idx, nil))
		if searchElement.SameAs(val) ***REMOVED***
			return valueTrue
		***REMOVED***
	***REMOVED***

	return valueFalse
***REMOVED***

func (r *Runtime) arrayproto_lastIndexOf(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	if length == 0 ***REMOVED***
		return intToValue(-1)
	***REMOVED***

	var fromIndex int64

	if len(call.Arguments) < 2 ***REMOVED***
		fromIndex = length - 1
	***REMOVED*** else ***REMOVED***
		fromIndex = call.Argument(1).ToInteger()
		if fromIndex >= 0 ***REMOVED***
			fromIndex = min(fromIndex, length-1)
		***REMOVED*** else ***REMOVED***
			fromIndex += length
		***REMOVED***
	***REMOVED***

	searchElement := call.Argument(0)

	if arr := r.checkStdArrayObj(o); arr != nil ***REMOVED***
		vals := arr.values
		for k := fromIndex; k >= 0; k-- ***REMOVED***
			if v := vals[k]; v != nil && searchElement.StrictEquals(v) ***REMOVED***
				return intToValue(k)
			***REMOVED***
		***REMOVED***
		return intToValue(-1)
	***REMOVED***

	for k := fromIndex; k >= 0; k-- ***REMOVED***
		idx := valueInt(k)
		if o.self.hasPropertyIdx(idx) ***REMOVED***
			if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
				if searchElement.StrictEquals(val) ***REMOVED***
					return idx
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return intToValue(-1)
***REMOVED***

func (r *Runtime) arrayproto_every(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	callbackFn := r.toCallable(call.Argument(0))
	fc := FunctionCall***REMOVED***
		This:      call.Argument(1),
		Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
	***REMOVED***
	for k := int64(0); k < length; k++ ***REMOVED***
		idx := valueInt(k)
		if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
			fc.Arguments[0] = val
			fc.Arguments[1] = idx
			if !callbackFn(fc).ToBoolean() ***REMOVED***
				return valueFalse
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) arrayproto_some(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	callbackFn := r.toCallable(call.Argument(0))
	fc := FunctionCall***REMOVED***
		This:      call.Argument(1),
		Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
	***REMOVED***
	for k := int64(0); k < length; k++ ***REMOVED***
		idx := valueInt(k)
		if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
			fc.Arguments[0] = val
			fc.Arguments[1] = idx
			if callbackFn(fc).ToBoolean() ***REMOVED***
				return valueTrue
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) arrayproto_forEach(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	callbackFn := r.toCallable(call.Argument(0))
	fc := FunctionCall***REMOVED***
		This:      call.Argument(1),
		Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
	***REMOVED***
	for k := int64(0); k < length; k++ ***REMOVED***
		idx := valueInt(k)
		if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
			fc.Arguments[0] = val
			fc.Arguments[1] = idx
			callbackFn(fc)
		***REMOVED***
	***REMOVED***
	return _undefined
***REMOVED***

func (r *Runtime) arrayproto_map(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	callbackFn := r.toCallable(call.Argument(0))
	fc := FunctionCall***REMOVED***
		This:      call.Argument(1),
		Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
	***REMOVED***
	a := arraySpeciesCreate(o, length)
	if _, stdSrc := o.self.(*arrayObject); stdSrc ***REMOVED***
		if arr, ok := a.self.(*arrayObject); ok ***REMOVED***
			values := make([]Value, length)
			for k := int64(0); k < length; k++ ***REMOVED***
				idx := valueInt(k)
				if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
					fc.Arguments[0] = val
					fc.Arguments[1] = idx
					values[k] = callbackFn(fc)
				***REMOVED***
			***REMOVED***
			setArrayValues(arr, values)
			return a
		***REMOVED***
	***REMOVED***
	for k := int64(0); k < length; k++ ***REMOVED***
		idx := valueInt(k)
		if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
			fc.Arguments[0] = val
			fc.Arguments[1] = idx
			createDataPropertyOrThrow(a, idx, callbackFn(fc))
		***REMOVED***
	***REMOVED***
	return a
***REMOVED***

func (r *Runtime) arrayproto_filter(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		a := arraySpeciesCreate(o, 0)
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
		***REMOVED***
		if _, stdSrc := o.self.(*arrayObject); stdSrc ***REMOVED***
			if arr := r.checkStdArrayObj(a); arr != nil ***REMOVED***
				var values []Value
				for k := int64(0); k < length; k++ ***REMOVED***
					idx := valueInt(k)
					if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
						fc.Arguments[0] = val
						fc.Arguments[1] = idx
						if callbackFn(fc).ToBoolean() ***REMOVED***
							values = append(values, val)
						***REMOVED***
					***REMOVED***
				***REMOVED***
				setArrayValues(arr, values)
				return a
			***REMOVED***
		***REMOVED***

		to := int64(0)
		for k := int64(0); k < length; k++ ***REMOVED***
			idx := valueInt(k)
			if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
				fc.Arguments[0] = val
				fc.Arguments[1] = idx
				if callbackFn(fc).ToBoolean() ***REMOVED***
					createDataPropertyOrThrow(a, intToValue(to), val)
					to++
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return a
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	panic("unreachable")
***REMOVED***

func (r *Runtime) arrayproto_reduce(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      _undefined,
			Arguments: []Value***REMOVED***nil, nil, nil, o***REMOVED***,
		***REMOVED***

		var k int64

		if len(call.Arguments) >= 2 ***REMOVED***
			fc.Arguments[0] = call.Argument(1)
		***REMOVED*** else ***REMOVED***
			for ; k < length; k++ ***REMOVED***
				idx := valueInt(k)
				if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
					fc.Arguments[0] = val
					break
				***REMOVED***
			***REMOVED***
			if fc.Arguments[0] == nil ***REMOVED***
				r.typeErrorResult(true, "No initial value")
				panic("unreachable")
			***REMOVED***
			k++
		***REMOVED***

		for ; k < length; k++ ***REMOVED***
			idx := valueInt(k)
			if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
				fc.Arguments[1] = val
				fc.Arguments[2] = idx
				fc.Arguments[0] = callbackFn(fc)
			***REMOVED***
		***REMOVED***
		return fc.Arguments[0]
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	panic("unreachable")
***REMOVED***

func (r *Runtime) arrayproto_reduceRight(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length", nil))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      _undefined,
			Arguments: []Value***REMOVED***nil, nil, nil, o***REMOVED***,
		***REMOVED***

		k := length - 1

		if len(call.Arguments) >= 2 ***REMOVED***
			fc.Arguments[0] = call.Argument(1)
		***REMOVED*** else ***REMOVED***
			for ; k >= 0; k-- ***REMOVED***
				idx := valueInt(k)
				if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
					fc.Arguments[0] = val
					break
				***REMOVED***
			***REMOVED***
			if fc.Arguments[0] == nil ***REMOVED***
				r.typeErrorResult(true, "No initial value")
				panic("unreachable")
			***REMOVED***
			k--
		***REMOVED***

		for ; k >= 0; k-- ***REMOVED***
			idx := valueInt(k)
			if val := o.self.getIdx(idx, nil); val != nil ***REMOVED***
				fc.Arguments[1] = val
				fc.Arguments[2] = idx
				fc.Arguments[0] = callbackFn(fc)
			***REMOVED***
		***REMOVED***
		return fc.Arguments[0]
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	panic("unreachable")
***REMOVED***

func arrayproto_reverse_generic_step(o *Object, lower, upper int64) ***REMOVED***
	lowerP := valueInt(lower)
	upperP := valueInt(upper)
	lowerValue := o.self.getIdx(lowerP, nil)
	upperValue := o.self.getIdx(upperP, nil)
	if lowerValue != nil && upperValue != nil ***REMOVED***
		o.self.setOwnIdx(lowerP, upperValue, true)
		o.self.setOwnIdx(upperP, lowerValue, true)
	***REMOVED*** else if lowerValue == nil && upperValue != nil ***REMOVED***
		o.self.setOwnIdx(lowerP, upperValue, true)
		o.self.deleteIdx(upperP, true)
	***REMOVED*** else if lowerValue != nil && upperValue == nil ***REMOVED***
		o.self.deleteIdx(lowerP, true)
		o.self.setOwnIdx(upperP, lowerValue, true)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_reverse_generic(o *Object, start int64) ***REMOVED***
	l := toLength(o.self.getStr("length", nil))
	middle := l / 2
	for lower := start; lower != middle; lower++ ***REMOVED***
		arrayproto_reverse_generic_step(o, lower, l-lower-1)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_reverse(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	if a := r.checkStdArrayObj(o); a != nil ***REMOVED***
		l := len(a.values)
		middle := l / 2
		for lower := 0; lower != middle; lower++ ***REMOVED***
			upper := l - lower - 1
			a.values[lower], a.values[upper] = a.values[upper], a.values[lower]
		***REMOVED***
		//TODO: go arrays
	***REMOVED*** else ***REMOVED***
		r.arrayproto_reverse_generic(o, 0)
	***REMOVED***
	return o
***REMOVED***

func (r *Runtime) arrayproto_shift(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	if a := r.checkStdArrayObj(o); a != nil ***REMOVED***
		if len(a.values) == 0 ***REMOVED***
			return _undefined
		***REMOVED***
		first := a.values[0]
		copy(a.values, a.values[1:])
		a.values[len(a.values)-1] = nil
		a.values = a.values[:len(a.values)-1]
		a.length--
		return first
	***REMOVED***
	length := toLength(o.self.getStr("length", nil))
	if length == 0 ***REMOVED***
		o.self.setOwnStr("length", intToValue(0), true)
		return _undefined
	***REMOVED***
	first := o.self.getIdx(valueInt(0), nil)
	for i := int64(1); i < length; i++ ***REMOVED***
		idxFrom := valueInt(i)
		idxTo := valueInt(i - 1)
		if o.self.hasPropertyIdx(idxFrom) ***REMOVED***
			o.self.setOwnIdx(idxTo, nilSafe(o.self.getIdx(idxFrom, nil)), true)
		***REMOVED*** else ***REMOVED***
			o.self.deleteIdx(idxTo, true)
		***REMOVED***
	***REMOVED***

	lv := valueInt(length - 1)
	o.self.deleteIdx(lv, true)
	o.self.setOwnStr("length", lv, true)

	return first
***REMOVED***

func (r *Runtime) arrayproto_values(call FunctionCall) Value ***REMOVED***
	return r.createArrayIterator(call.This.ToObject(r), iterationKindValue)
***REMOVED***

func (r *Runtime) arrayproto_keys(call FunctionCall) Value ***REMOVED***
	return r.createArrayIterator(call.This.ToObject(r), iterationKindKey)
***REMOVED***

func (r *Runtime) arrayproto_copyWithin(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := toLength(o.self.getStr("length", nil))
	var relEnd, dir int64
	to := relToIdx(call.Argument(0).ToInteger(), l)
	from := relToIdx(call.Argument(1).ToInteger(), l)
	if end := call.Argument(2); end != _undefined ***REMOVED***
		relEnd = end.ToInteger()
	***REMOVED*** else ***REMOVED***
		relEnd = l
	***REMOVED***
	final := relToIdx(relEnd, l)
	count := min(final-from, l-to)
	if arr := r.checkStdArrayObj(o); arr != nil ***REMOVED***
		if count > 0 ***REMOVED***
			copy(arr.values[to:to+count], arr.values[from:from+count])
		***REMOVED***
		return o
	***REMOVED***
	if from < to && to < from+count ***REMOVED***
		dir = -1
		from = from + count - 1
		to = to + count - 1
	***REMOVED*** else ***REMOVED***
		dir = 1
	***REMOVED***
	for count > 0 ***REMOVED***
		if o.self.hasPropertyIdx(valueInt(from)) ***REMOVED***
			o.self.setOwnIdx(valueInt(to), nilSafe(o.self.getIdx(valueInt(from), nil)), true)
		***REMOVED*** else ***REMOVED***
			o.self.deleteIdx(valueInt(to), true)
		***REMOVED***
		from += dir
		to += dir
		count--
	***REMOVED***

	return o
***REMOVED***

func (r *Runtime) arrayproto_entries(call FunctionCall) Value ***REMOVED***
	return r.createArrayIterator(call.This.ToObject(r), iterationKindKeyValue)
***REMOVED***

func (r *Runtime) arrayproto_fill(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := toLength(o.self.getStr("length", nil))
	k := relToIdx(call.Argument(1).ToInteger(), l)
	var relEnd int64
	if endArg := call.Argument(2); endArg != _undefined ***REMOVED***
		relEnd = endArg.ToInteger()
	***REMOVED*** else ***REMOVED***
		relEnd = l
	***REMOVED***
	final := relToIdx(relEnd, l)
	value := call.Argument(0)
	if arr := r.checkStdArrayObj(o); arr != nil ***REMOVED***
		for ; k < final; k++ ***REMOVED***
			arr.values[k] = value
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for ; k < final; k++ ***REMOVED***
			o.self.setOwnIdx(valueInt(k), value, true)
		***REMOVED***
	***REMOVED***
	return o
***REMOVED***

func (r *Runtime) arrayproto_find(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := toLength(o.self.getStr("length", nil))
	predicate := r.toCallable(call.Argument(0))
	fc := FunctionCall***REMOVED***
		This:      call.Argument(1),
		Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
	***REMOVED***
	for k := int64(0); k < l; k++ ***REMOVED***
		idx := valueInt(k)
		kValue := o.self.getIdx(idx, nil)
		fc.Arguments[0], fc.Arguments[1] = kValue, idx
		if predicate(fc).ToBoolean() ***REMOVED***
			return kValue
		***REMOVED***
	***REMOVED***

	return _undefined
***REMOVED***

func (r *Runtime) arrayproto_findIndex(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := toLength(o.self.getStr("length", nil))
	predicate := r.toCallable(call.Argument(0))
	fc := FunctionCall***REMOVED***
		This:      call.Argument(1),
		Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
	***REMOVED***
	for k := int64(0); k < l; k++ ***REMOVED***
		idx := valueInt(k)
		kValue := o.self.getIdx(idx, nil)
		fc.Arguments[0], fc.Arguments[1] = kValue, idx
		if predicate(fc).ToBoolean() ***REMOVED***
			return idx
		***REMOVED***
	***REMOVED***

	return intToValue(-1)
***REMOVED***

func (r *Runtime) arrayproto_flat(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := toLength(o.self.getStr("length", nil))
	depthNum := int64(1)
	if len(call.Arguments) > 0 ***REMOVED***
		depthNum = call.Argument(0).ToInteger()
	***REMOVED***
	a := arraySpeciesCreate(o, 0)
	r.flattenIntoArray(a, o, l, 0, depthNum, nil, nil)
	return a
***REMOVED***

func (r *Runtime) flattenIntoArray(target, source *Object, sourceLen, start, depth int64, mapperFunction func(FunctionCall) Value, thisArg Value) int64 ***REMOVED***
	targetIndex, sourceIndex := start, int64(0)
	for sourceIndex < sourceLen ***REMOVED***
		p := intToValue(sourceIndex)
		if source.hasProperty(p.toString()) ***REMOVED***
			element := nilSafe(source.get(p, source))
			if mapperFunction != nil ***REMOVED***
				element = mapperFunction(FunctionCall***REMOVED***
					This:      thisArg,
					Arguments: []Value***REMOVED***element, p, source***REMOVED***,
				***REMOVED***)
			***REMOVED***
			var elementArray *Object
			if depth > 0 ***REMOVED***
				if elementObj, ok := element.(*Object); ok && isArray(elementObj) ***REMOVED***
					elementArray = elementObj
				***REMOVED***
			***REMOVED***
			if elementArray != nil ***REMOVED***
				elementLen := toLength(elementArray.self.getStr("length", nil))
				targetIndex = r.flattenIntoArray(target, elementArray, elementLen, targetIndex, depth-1, nil, nil)
			***REMOVED*** else ***REMOVED***
				if targetIndex >= maxInt-1 ***REMOVED***
					panic(r.NewTypeError("Invalid array length"))
				***REMOVED***
				createDataPropertyOrThrow(target, intToValue(targetIndex), element)
				targetIndex++
			***REMOVED***
		***REMOVED***
		sourceIndex++
	***REMOVED***
	return targetIndex
***REMOVED***

func (r *Runtime) arrayproto_flatMap(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := toLength(o.self.getStr("length", nil))
	callbackFn := r.toCallable(call.Argument(0))
	thisArg := Undefined()
	if len(call.Arguments) > 1 ***REMOVED***
		thisArg = call.Argument(1)
	***REMOVED***
	a := arraySpeciesCreate(o, 0)
	r.flattenIntoArray(a, o, l, 0, 1, callbackFn, thisArg)
	return a
***REMOVED***

func (r *Runtime) checkStdArrayObj(obj *Object) *arrayObject ***REMOVED***
	if arr, ok := obj.self.(*arrayObject); ok &&
		arr.propValueCount == 0 &&
		arr.length == uint32(len(arr.values)) &&
		uint32(arr.objCount) == arr.length ***REMOVED***

		return arr
	***REMOVED***

	return nil
***REMOVED***

func (r *Runtime) checkStdArray(v Value) *arrayObject ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		return r.checkStdArrayObj(obj)
	***REMOVED***

	return nil
***REMOVED***

func (r *Runtime) checkStdArrayIter(v Value) *arrayObject ***REMOVED***
	if arr := r.checkStdArray(v); arr != nil &&
		arr.getSym(SymIterator, nil) == r.global.arrayValues ***REMOVED***

		return arr
	***REMOVED***

	return nil
***REMOVED***

func (r *Runtime) array_from(call FunctionCall) Value ***REMOVED***
	var mapFn func(FunctionCall) Value
	if mapFnArg := call.Argument(1); mapFnArg != _undefined ***REMOVED***
		if mapFnObj, ok := mapFnArg.(*Object); ok ***REMOVED***
			if fn, ok := mapFnObj.self.assertCallable(); ok ***REMOVED***
				mapFn = fn
			***REMOVED***
		***REMOVED***
		if mapFn == nil ***REMOVED***
			panic(r.NewTypeError("%s is not a function", mapFnArg))
		***REMOVED***
	***REMOVED***
	t := call.Argument(2)
	items := call.Argument(0)
	if mapFn == nil && call.This == r.global.Array ***REMOVED*** // mapFn may mutate the array
		if arr := r.checkStdArrayIter(items); arr != nil ***REMOVED***
			items := make([]Value, len(arr.values))
			copy(items, arr.values)
			return r.newArrayValues(items)
		***REMOVED***
	***REMOVED***

	var ctor func(args []Value, newTarget *Object) *Object
	if call.This != r.global.Array ***REMOVED***
		if o, ok := call.This.(*Object); ok ***REMOVED***
			if c := o.self.assertConstructor(); c != nil ***REMOVED***
				ctor = c
			***REMOVED***
		***REMOVED***
	***REMOVED***
	var arr *Object
	if usingIterator := toMethod(r.getV(items, SymIterator)); usingIterator != nil ***REMOVED***
		if ctor != nil ***REMOVED***
			arr = ctor([]Value***REMOVED******REMOVED***, nil)
		***REMOVED*** else ***REMOVED***
			arr = r.newArrayValues(nil)
		***REMOVED***
		iter := r.getIterator(items, usingIterator)
		if mapFn == nil ***REMOVED***
			if a := r.checkStdArrayObj(arr); a != nil ***REMOVED***
				var values []Value
				r.iterate(iter, func(val Value) ***REMOVED***
					values = append(values, val)
				***REMOVED***)
				setArrayValues(a, values)
				return arr
			***REMOVED***
		***REMOVED***
		k := int64(0)
		r.iterate(iter, func(val Value) ***REMOVED***
			if mapFn != nil ***REMOVED***
				val = mapFn(FunctionCall***REMOVED***This: t, Arguments: []Value***REMOVED***val, intToValue(k)***REMOVED******REMOVED***)
			***REMOVED***
			createDataPropertyOrThrow(arr, intToValue(k), val)
			k++
		***REMOVED***)
		arr.self.setOwnStr("length", intToValue(k), true)
	***REMOVED*** else ***REMOVED***
		arrayLike := items.ToObject(r)
		l := toLength(arrayLike.self.getStr("length", nil))
		if ctor != nil ***REMOVED***
			arr = ctor([]Value***REMOVED***intToValue(l)***REMOVED***, nil)
		***REMOVED*** else ***REMOVED***
			arr = r.newArrayValues(nil)
		***REMOVED***
		if mapFn == nil ***REMOVED***
			if a := r.checkStdArrayObj(arr); a != nil ***REMOVED***
				values := make([]Value, l)
				for k := int64(0); k < l; k++ ***REMOVED***
					values[k] = nilSafe(arrayLike.self.getIdx(valueInt(k), nil))
				***REMOVED***
				setArrayValues(a, values)
				return arr
			***REMOVED***
		***REMOVED***
		for k := int64(0); k < l; k++ ***REMOVED***
			idx := valueInt(k)
			item := arrayLike.self.getIdx(idx, nil)
			if mapFn != nil ***REMOVED***
				item = mapFn(FunctionCall***REMOVED***This: t, Arguments: []Value***REMOVED***item, idx***REMOVED******REMOVED***)
			***REMOVED*** else ***REMOVED***
				item = nilSafe(item)
			***REMOVED***
			createDataPropertyOrThrow(arr, idx, item)
		***REMOVED***
		arr.self.setOwnStr("length", intToValue(l), true)
	***REMOVED***

	return arr
***REMOVED***

func (r *Runtime) array_isArray(call FunctionCall) Value ***REMOVED***
	if o, ok := call.Argument(0).(*Object); ok ***REMOVED***
		if isArray(o) ***REMOVED***
			return valueTrue
		***REMOVED***
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) array_of(call FunctionCall) Value ***REMOVED***
	var ctor func(args []Value, newTarget *Object) *Object
	if call.This != r.global.Array ***REMOVED***
		if o, ok := call.This.(*Object); ok ***REMOVED***
			if c := o.self.assertConstructor(); c != nil ***REMOVED***
				ctor = c
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if ctor == nil ***REMOVED***
		values := make([]Value, len(call.Arguments))
		copy(values, call.Arguments)
		return r.newArrayValues(values)
	***REMOVED***
	l := intToValue(int64(len(call.Arguments)))
	arr := ctor([]Value***REMOVED***l***REMOVED***, nil)
	for i, val := range call.Arguments ***REMOVED***
		createDataPropertyOrThrow(arr, intToValue(int64(i)), val)
	***REMOVED***
	arr.self.setOwnStr("length", l, true)
	return arr
***REMOVED***

func (r *Runtime) arrayIterProto_next(call FunctionCall) Value ***REMOVED***
	thisObj := r.toObject(call.This)
	if iter, ok := thisObj.self.(*arrayIterObject); ok ***REMOVED***
		return iter.next()
	***REMOVED***
	panic(r.NewTypeError("Method Array Iterator.prototype.next called on incompatible receiver %s", thisObj.String()))
***REMOVED***

func (r *Runtime) createArrayProto(val *Object) objectImpl ***REMOVED***
	o := &arrayObject***REMOVED***
		baseObject: baseObject***REMOVED***
			class:      classArray,
			val:        val,
			extensible: true,
			prototype:  r.global.ObjectPrototype,
		***REMOVED***,
	***REMOVED***
	o.init()

	o._putProp("constructor", r.global.Array, true, false, true)
	o._putProp("concat", r.newNativeFunc(r.arrayproto_concat, nil, "concat", nil, 1), true, false, true)
	o._putProp("copyWithin", r.newNativeFunc(r.arrayproto_copyWithin, nil, "copyWithin", nil, 2), true, false, true)
	o._putProp("entries", r.newNativeFunc(r.arrayproto_entries, nil, "entries", nil, 0), true, false, true)
	o._putProp("every", r.newNativeFunc(r.arrayproto_every, nil, "every", nil, 1), true, false, true)
	o._putProp("fill", r.newNativeFunc(r.arrayproto_fill, nil, "fill", nil, 1), true, false, true)
	o._putProp("filter", r.newNativeFunc(r.arrayproto_filter, nil, "filter", nil, 1), true, false, true)
	o._putProp("find", r.newNativeFunc(r.arrayproto_find, nil, "find", nil, 1), true, false, true)
	o._putProp("findIndex", r.newNativeFunc(r.arrayproto_findIndex, nil, "findIndex", nil, 1), true, false, true)
	o._putProp("flat", r.newNativeFunc(r.arrayproto_flat, nil, "flat", nil, 0), true, false, true)
	o._putProp("flatMap", r.newNativeFunc(r.arrayproto_flatMap, nil, "flatMap", nil, 1), true, false, true)
	o._putProp("forEach", r.newNativeFunc(r.arrayproto_forEach, nil, "forEach", nil, 1), true, false, true)
	o._putProp("includes", r.newNativeFunc(r.arrayproto_includes, nil, "includes", nil, 1), true, false, true)
	o._putProp("indexOf", r.newNativeFunc(r.arrayproto_indexOf, nil, "indexOf", nil, 1), true, false, true)
	o._putProp("join", r.newNativeFunc(r.arrayproto_join, nil, "join", nil, 1), true, false, true)
	o._putProp("keys", r.newNativeFunc(r.arrayproto_keys, nil, "keys", nil, 0), true, false, true)
	o._putProp("lastIndexOf", r.newNativeFunc(r.arrayproto_lastIndexOf, nil, "lastIndexOf", nil, 1), true, false, true)
	o._putProp("map", r.newNativeFunc(r.arrayproto_map, nil, "map", nil, 1), true, false, true)
	o._putProp("pop", r.newNativeFunc(r.arrayproto_pop, nil, "pop", nil, 0), true, false, true)
	o._putProp("push", r.newNativeFunc(r.arrayproto_push, nil, "push", nil, 1), true, false, true)
	o._putProp("reduce", r.newNativeFunc(r.arrayproto_reduce, nil, "reduce", nil, 1), true, false, true)
	o._putProp("reduceRight", r.newNativeFunc(r.arrayproto_reduceRight, nil, "reduceRight", nil, 1), true, false, true)
	o._putProp("reverse", r.newNativeFunc(r.arrayproto_reverse, nil, "reverse", nil, 0), true, false, true)
	o._putProp("shift", r.newNativeFunc(r.arrayproto_shift, nil, "shift", nil, 0), true, false, true)
	o._putProp("slice", r.newNativeFunc(r.arrayproto_slice, nil, "slice", nil, 2), true, false, true)
	o._putProp("some", r.newNativeFunc(r.arrayproto_some, nil, "some", nil, 1), true, false, true)
	o._putProp("sort", r.newNativeFunc(r.arrayproto_sort, nil, "sort", nil, 1), true, false, true)
	o._putProp("splice", r.newNativeFunc(r.arrayproto_splice, nil, "splice", nil, 2), true, false, true)
	o._putProp("toLocaleString", r.newNativeFunc(r.arrayproto_toLocaleString, nil, "toLocaleString", nil, 0), true, false, true)
	o._putProp("toString", r.global.arrayToString, true, false, true)
	o._putProp("unshift", r.newNativeFunc(r.arrayproto_unshift, nil, "unshift", nil, 1), true, false, true)
	o._putProp("values", r.global.arrayValues, true, false, true)

	o._putSym(SymIterator, valueProp(r.global.arrayValues, true, false, true))

	bl := r.newBaseObject(nil, classObject)
	bl.setOwnStr("copyWithin", valueTrue, true)
	bl.setOwnStr("entries", valueTrue, true)
	bl.setOwnStr("fill", valueTrue, true)
	bl.setOwnStr("find", valueTrue, true)
	bl.setOwnStr("findIndex", valueTrue, true)
	bl.setOwnStr("flat", valueTrue, true)
	bl.setOwnStr("flatMap", valueTrue, true)
	bl.setOwnStr("includes", valueTrue, true)
	bl.setOwnStr("keys", valueTrue, true)
	bl.setOwnStr("values", valueTrue, true)
	o._putSym(SymUnscopables, valueProp(bl.val, false, false, true))

	return o
***REMOVED***

func (r *Runtime) createArray(val *Object) objectImpl ***REMOVED***
	o := r.newNativeFuncConstructObj(val, r.builtin_newArray, "Array", r.global.ArrayPrototype, 1)
	o._putProp("from", r.newNativeFunc(r.array_from, nil, "from", nil, 1), true, false, true)
	o._putProp("isArray", r.newNativeFunc(r.array_isArray, nil, "isArray", nil, 1), true, false, true)
	o._putProp("of", r.newNativeFunc(r.array_of, nil, "of", nil, 0), true, false, true)
	r.putSpeciesReturnThis(o)

	return o
***REMOVED***

func (r *Runtime) createArrayIterProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.IteratorPrototype, classObject)

	o._putProp("next", r.newNativeFunc(r.arrayIterProto_next, nil, "next", nil, 0), true, false, true)
	o._putSym(SymToStringTag, valueProp(asciiString(classArrayIterator), false, false, true))

	return o
***REMOVED***

func (r *Runtime) initArray() ***REMOVED***
	r.global.arrayValues = r.newNativeFunc(r.arrayproto_values, nil, "values", nil, 0)
	r.global.arrayToString = r.newNativeFunc(r.arrayproto_toString, nil, "toString", nil, 0)

	r.global.ArrayIteratorPrototype = r.newLazyObject(r.createArrayIterProto)
	//r.global.ArrayPrototype = r.newArray(r.global.ObjectPrototype).val
	//o := r.global.ArrayPrototype.self
	r.global.ArrayPrototype = r.newLazyObject(r.createArrayProto)

	//r.global.Array = r.newNativeFuncConstruct(r.builtin_newArray, "Array", r.global.ArrayPrototype, 1)
	//o = r.global.Array.self
	//o._putProp("isArray", r.newNativeFunc(r.array_isArray, nil, "isArray", nil, 1), true, false, true)
	r.global.Array = r.newLazyObject(r.createArray)

	r.addToGlobal("Array", r.global.Array)
***REMOVED***

type sortable interface ***REMOVED***
	sortLen() int64
	sortGet(int64) Value
	swap(int64, int64)
***REMOVED***

type arraySortCtx struct ***REMOVED***
	obj     sortable
	compare func(FunctionCall) Value
***REMOVED***

func (a *arraySortCtx) sortCompare(x, y Value) int ***REMOVED***
	if x == nil && y == nil ***REMOVED***
		return 0
	***REMOVED***

	if x == nil ***REMOVED***
		return 1
	***REMOVED***

	if y == nil ***REMOVED***
		return -1
	***REMOVED***

	if x == _undefined && y == _undefined ***REMOVED***
		return 0
	***REMOVED***

	if x == _undefined ***REMOVED***
		return 1
	***REMOVED***

	if y == _undefined ***REMOVED***
		return -1
	***REMOVED***

	if a.compare != nil ***REMOVED***
		f := a.compare(FunctionCall***REMOVED***
			This:      _undefined,
			Arguments: []Value***REMOVED***x, y***REMOVED***,
		***REMOVED***).ToFloat()
		if f > 0 ***REMOVED***
			return 1
		***REMOVED***
		if f < 0 ***REMOVED***
			return -1
		***REMOVED***
		if math.Signbit(f) ***REMOVED***
			return -1
		***REMOVED***
		return 0
	***REMOVED***
	return x.toString().compareTo(y.toString())
***REMOVED***

// sort.Interface

func (a *arraySortCtx) Len() int ***REMOVED***
	return int(a.obj.sortLen())
***REMOVED***

func (a *arraySortCtx) Less(j, k int) bool ***REMOVED***
	return a.sortCompare(a.obj.sortGet(int64(j)), a.obj.sortGet(int64(k))) < 0
***REMOVED***

func (a *arraySortCtx) Swap(j, k int) ***REMOVED***
	a.obj.swap(int64(j), int64(k))
***REMOVED***
