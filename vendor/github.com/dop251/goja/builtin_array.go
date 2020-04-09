package goja

import (
	"bytes"
	"math"
	"sort"
	"strings"
)

func (r *Runtime) builtin_newArray(args []Value, proto *Object) *Object ***REMOVED***
	l := len(args)
	if l == 1 ***REMOVED***
		if al, ok := args[0].assertInt(); ok ***REMOVED***
			return r.newArrayLength(al)
		***REMOVED*** else if f, ok := args[0].assertFloat(); ok ***REMOVED***
			al := int64(f)
			if float64(al) == f ***REMOVED***
				return r.newArrayLength(al)
			***REMOVED*** else ***REMOVED***
				panic(r.newError(r.global.RangeError, "Invalid array length"))
			***REMOVED***
		***REMOVED***
		return r.newArrayValues([]Value***REMOVED***args[0]***REMOVED***)
	***REMOVED*** else ***REMOVED***
		argsCopy := make([]Value, l)
		copy(argsCopy, args)
		return r.newArrayValues(argsCopy)
	***REMOVED***
***REMOVED***

func (r *Runtime) generic_push(obj *Object, call FunctionCall) Value ***REMOVED***
	l := toLength(obj.self.getStr("length"))
	nl := l + int64(len(call.Arguments))
	if nl >= maxInt ***REMOVED***
		r.typeErrorResult(true, "Invalid array length")
		panic("unreachable")
	***REMOVED***
	for i, arg := range call.Arguments ***REMOVED***
		obj.self.put(intToValue(l+int64(i)), arg, true)
	***REMOVED***
	n := intToValue(nl)
	obj.self.putStr("length", n, true)
	return n
***REMOVED***

func (r *Runtime) arrayproto_push(call FunctionCall) Value ***REMOVED***
	obj := call.This.ToObject(r)
	return r.generic_push(obj, call)
***REMOVED***

func (r *Runtime) arrayproto_pop_generic(obj *Object, call FunctionCall) Value ***REMOVED***
	l := toLength(obj.self.getStr("length"))
	if l == 0 ***REMOVED***
		obj.self.putStr("length", intToValue(0), true)
		return _undefined
	***REMOVED***
	idx := intToValue(l - 1)
	val := obj.self.get(idx)
	obj.self.delete(idx, true)
	obj.self.putStr("length", idx, true)
	return val
***REMOVED***

func (r *Runtime) arrayproto_pop(call FunctionCall) Value ***REMOVED***
	obj := call.This.ToObject(r)
	if a, ok := obj.self.(*arrayObject); ok ***REMOVED***
		l := a.length
		if l > 0 ***REMOVED***
			var val Value
			l--
			if l < int64(len(a.values)) ***REMOVED***
				val = a.values[l]
			***REMOVED***
			if val == nil ***REMOVED***
				// optimisation bail-out
				return r.arrayproto_pop_generic(obj, call)
			***REMOVED***
			if _, ok := val.(*valueProperty); ok ***REMOVED***
				// optimisation bail-out
				return r.arrayproto_pop_generic(obj, call)
			***REMOVED***
			//a._setLengthInt(l, false)
			a.values[l] = nil
			a.values = a.values[:l]
			a.length = l
			return val
		***REMOVED***
		return _undefined
	***REMOVED*** else ***REMOVED***
		return r.arrayproto_pop_generic(obj, call)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_join(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	l := int(toLength(o.self.getStr("length")))
	sep := ""
	if s := call.Argument(0); s != _undefined ***REMOVED***
		sep = s.String()
	***REMOVED*** else ***REMOVED***
		sep = ","
	***REMOVED***
	if l == 0 ***REMOVED***
		return stringEmpty
	***REMOVED***

	var buf bytes.Buffer

	element0 := o.self.get(intToValue(0))
	if element0 != nil && element0 != _undefined && element0 != _null ***REMOVED***
		buf.WriteString(element0.String())
	***REMOVED***

	for i := 1; i < l; i++ ***REMOVED***
		buf.WriteString(sep)
		element := o.self.get(intToValue(int64(i)))
		if element != nil && element != _undefined && element != _null ***REMOVED***
			buf.WriteString(element.String())
		***REMOVED***
	***REMOVED***

	return newStringValue(buf.String())
***REMOVED***

func (r *Runtime) arrayproto_toString(call FunctionCall) Value ***REMOVED***
	array := call.This.ToObject(r)
	f := array.self.getStr("join")
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

func (r *Runtime) writeItemLocaleString(item Value, buf *bytes.Buffer) ***REMOVED***
	if item != nil && item != _undefined && item != _null ***REMOVED***
		itemObj := item.ToObject(r)
		if f, ok := itemObj.self.getStr("toLocaleString").(*Object); ok ***REMOVED***
			if c, ok := f.self.assertCallable(); ok ***REMOVED***
				strVal := c(FunctionCall***REMOVED***
					This: itemObj,
				***REMOVED***)
				buf.WriteString(strVal.String())
				return
			***REMOVED***
		***REMOVED***
		r.typeErrorResult(true, "Property 'toLocaleString' of object %s is not a function", itemObj)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_toLocaleString_generic(obj *Object, start int64, buf *bytes.Buffer) Value ***REMOVED***
	length := toLength(obj.self.getStr("length"))
	for i := int64(start); i < length; i++ ***REMOVED***
		if i > 0 ***REMOVED***
			buf.WriteByte(',')
		***REMOVED***
		item := obj.self.get(intToValue(i))
		r.writeItemLocaleString(item, buf)
	***REMOVED***
	return newStringValue(buf.String())
***REMOVED***

func (r *Runtime) arrayproto_toLocaleString(call FunctionCall) Value ***REMOVED***
	array := call.This.ToObject(r)
	if a, ok := array.self.(*arrayObject); ok ***REMOVED***
		var buf bytes.Buffer
		for i := int64(0); i < a.length; i++ ***REMOVED***
			var item Value
			if i < int64(len(a.values)) ***REMOVED***
				item = a.values[i]
			***REMOVED***
			if item == nil ***REMOVED***
				return r.arrayproto_toLocaleString_generic(array, i, &buf)
			***REMOVED***
			if prop, ok := item.(*valueProperty); ok ***REMOVED***
				item = prop.get(array)
			***REMOVED***
			if i > 0 ***REMOVED***
				buf.WriteByte(',')
			***REMOVED***
			r.writeItemLocaleString(item, &buf)
		***REMOVED***
		return newStringValue(buf.String())
	***REMOVED*** else ***REMOVED***
		return r.arrayproto_toLocaleString_generic(array, 0, bytes.NewBuffer(nil))
	***REMOVED***

***REMOVED***

func (r *Runtime) arrayproto_concat_append(a *Object, item Value) ***REMOVED***
	descr := propertyDescr***REMOVED***
		Writable:     FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
		Configurable: FLAG_TRUE,
	***REMOVED***

	aLength := toLength(a.self.getStr("length"))
	if obj, ok := item.(*Object); ok ***REMOVED***
		if isArray(obj) ***REMOVED***
			length := toLength(obj.self.getStr("length"))
			for i := int64(0); i < length; i++ ***REMOVED***
				v := obj.self.get(intToValue(i))
				if v != nil ***REMOVED***
					descr.Value = v
					a.self.defineOwnProperty(intToValue(aLength), descr, false)
					aLength++
				***REMOVED*** else ***REMOVED***
					aLength++
					a.self.putStr("length", intToValue(aLength), false)
				***REMOVED***
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	descr.Value = item
	a.self.defineOwnProperty(intToValue(aLength), descr, false)
***REMOVED***

func (r *Runtime) arrayproto_concat(call FunctionCall) Value ***REMOVED***
	a := r.newArrayValues(nil)
	r.arrayproto_concat_append(a, call.This.ToObject(r))
	for _, item := range call.Arguments ***REMOVED***
		r.arrayproto_concat_append(a, item)
	***REMOVED***
	return a
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

func (r *Runtime) arrayproto_slice(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	start := call.Argument(0).ToInteger()
	if start < 0 ***REMOVED***
		start = max(length+start, 0)
	***REMOVED*** else ***REMOVED***
		start = min(start, length)
	***REMOVED***
	var end int64
	if endArg := call.Argument(1); endArg != _undefined ***REMOVED***
		end = endArg.ToInteger()
	***REMOVED*** else ***REMOVED***
		end = length
	***REMOVED***
	if end < 0 ***REMOVED***
		end = max(length+end, 0)
	***REMOVED*** else ***REMOVED***
		end = min(end, length)
	***REMOVED***

	count := end - start
	if count < 0 ***REMOVED***
		count = 0
	***REMOVED***
	a := r.newArrayLength(count)

	n := int64(0)
	descr := propertyDescr***REMOVED***
		Writable:     FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
		Configurable: FLAG_TRUE,
	***REMOVED***
	for start < end ***REMOVED***
		p := o.self.get(intToValue(start))
		if p != nil && p != _undefined ***REMOVED***
			descr.Value = p
			a.self.defineOwnProperty(intToValue(n), descr, false)
		***REMOVED***
		start++
		n++
	***REMOVED***
	return a
***REMOVED***

func (r *Runtime) arrayproto_sort(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)

	var compareFn func(FunctionCall) Value

	if arg, ok := call.Argument(0).(*Object); ok ***REMOVED***
		compareFn, _ = arg.self.assertCallable()
	***REMOVED***

	ctx := arraySortCtx***REMOVED***
		obj:     o.self,
		compare: compareFn,
	***REMOVED***

	sort.Sort(&ctx)
	return o
***REMOVED***

func (r *Runtime) arrayproto_splice(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	a := r.newArrayValues(nil)
	length := toLength(o.self.getStr("length"))
	relativeStart := call.Argument(0).ToInteger()
	var actualStart int64
	if relativeStart < 0 ***REMOVED***
		actualStart = max(length+relativeStart, 0)
	***REMOVED*** else ***REMOVED***
		actualStart = min(relativeStart, length)
	***REMOVED***

	actualDeleteCount := min(max(call.Argument(1).ToInteger(), 0), length-actualStart)

	for k := int64(0); k < actualDeleteCount; k++ ***REMOVED***
		from := intToValue(k + actualStart)
		if o.self.hasProperty(from) ***REMOVED***
			a.self.put(intToValue(k), o.self.get(from), false)
		***REMOVED***
	***REMOVED***

	itemCount := max(int64(len(call.Arguments)-2), 0)
	if itemCount < actualDeleteCount ***REMOVED***
		for k := actualStart; k < length-actualDeleteCount; k++ ***REMOVED***
			from := intToValue(k + actualDeleteCount)
			to := intToValue(k + itemCount)
			if o.self.hasProperty(from) ***REMOVED***
				o.self.put(to, o.self.get(from), true)
			***REMOVED*** else ***REMOVED***
				o.self.delete(to, true)
			***REMOVED***
		***REMOVED***

		for k := length; k > length-actualDeleteCount+itemCount; k-- ***REMOVED***
			o.self.delete(intToValue(k-1), true)
		***REMOVED***
	***REMOVED*** else if itemCount > actualDeleteCount ***REMOVED***
		for k := length - actualDeleteCount; k > actualStart; k-- ***REMOVED***
			from := intToValue(k + actualDeleteCount - 1)
			to := intToValue(k + itemCount - 1)
			if o.self.hasProperty(from) ***REMOVED***
				o.self.put(to, o.self.get(from), true)
			***REMOVED*** else ***REMOVED***
				o.self.delete(to, true)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if itemCount > 0 ***REMOVED***
		for i, item := range call.Arguments[2:] ***REMOVED***
			o.self.put(intToValue(actualStart+int64(i)), item, true)
		***REMOVED***
	***REMOVED***

	o.self.putStr("length", intToValue(length-actualDeleteCount+itemCount), true)

	return a
***REMOVED***

func (r *Runtime) arrayproto_unshift(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	argCount := int64(len(call.Arguments))
	for k := length - 1; k >= 0; k-- ***REMOVED***
		from := intToValue(k)
		to := intToValue(k + argCount)
		if o.self.hasProperty(from) ***REMOVED***
			o.self.put(to, o.self.get(from), true)
		***REMOVED*** else ***REMOVED***
			o.self.delete(to, true)
		***REMOVED***
	***REMOVED***

	for k, arg := range call.Arguments ***REMOVED***
		o.self.put(intToValue(int64(k)), arg, true)
	***REMOVED***

	newLen := intToValue(length + argCount)
	o.self.putStr("length", newLen, true)
	return newLen
***REMOVED***

func (r *Runtime) arrayproto_indexOf(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
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

	for ; n < length; n++ ***REMOVED***
		idx := intToValue(n)
		if val := o.self.get(idx); val != nil ***REMOVED***
			if searchElement.StrictEquals(val) ***REMOVED***
				return idx
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return intToValue(-1)
***REMOVED***

func (r *Runtime) arrayproto_lastIndexOf(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
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

	for k := fromIndex; k >= 0; k-- ***REMOVED***
		idx := intToValue(k)
		if val := o.self.get(idx); val != nil ***REMOVED***
			if searchElement.StrictEquals(val) ***REMOVED***
				return idx
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return intToValue(-1)
***REMOVED***

func (r *Runtime) arrayproto_every(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
		***REMOVED***
		for k := int64(0); k < length; k++ ***REMOVED***
			idx := intToValue(k)
			if val := o.self.get(idx); val != nil ***REMOVED***
				fc.Arguments[0] = val
				fc.Arguments[1] = idx
				if !callbackFn(fc).ToBoolean() ***REMOVED***
					return valueFalse
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	return valueTrue
***REMOVED***

func (r *Runtime) arrayproto_some(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
		***REMOVED***
		for k := int64(0); k < length; k++ ***REMOVED***
			idx := intToValue(k)
			if val := o.self.get(idx); val != nil ***REMOVED***
				fc.Arguments[0] = val
				fc.Arguments[1] = idx
				if callbackFn(fc).ToBoolean() ***REMOVED***
					return valueTrue
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) arrayproto_forEach(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
		***REMOVED***
		for k := int64(0); k < length; k++ ***REMOVED***
			idx := intToValue(k)
			if val := o.self.get(idx); val != nil ***REMOVED***
				fc.Arguments[0] = val
				fc.Arguments[1] = idx
				callbackFn(fc)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	return _undefined
***REMOVED***

func (r *Runtime) arrayproto_map(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
		***REMOVED***
		a := r.newArrayObject()
		a._setLengthInt(length, true)
		a.values = make([]Value, length)
		for k := int64(0); k < length; k++ ***REMOVED***
			idx := intToValue(k)
			if val := o.self.get(idx); val != nil ***REMOVED***
				fc.Arguments[0] = val
				fc.Arguments[1] = idx
				a.values[k] = callbackFn(fc)
				a.objCount++
			***REMOVED***
		***REMOVED***
		return a.val
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	panic("unreachable")
***REMOVED***

func (r *Runtime) arrayproto_filter(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	callbackFn := call.Argument(0).ToObject(r)
	if callbackFn, ok := callbackFn.self.assertCallable(); ok ***REMOVED***
		a := r.newArrayObject()
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, o***REMOVED***,
		***REMOVED***
		for k := int64(0); k < length; k++ ***REMOVED***
			idx := intToValue(k)
			if val := o.self.get(idx); val != nil ***REMOVED***
				fc.Arguments[0] = val
				fc.Arguments[1] = idx
				if callbackFn(fc).ToBoolean() ***REMOVED***
					a.values = append(a.values, val)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		a.length = int64(len(a.values))
		a.objCount = a.length
		return a.val
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "%s is not a function", call.Argument(0))
	***REMOVED***
	panic("unreachable")
***REMOVED***

func (r *Runtime) arrayproto_reduce(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
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
				idx := intToValue(k)
				if val := o.self.get(idx); val != nil ***REMOVED***
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
			idx := intToValue(k)
			if val := o.self.get(idx); val != nil ***REMOVED***
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
	length := toLength(o.self.getStr("length"))
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
				idx := intToValue(k)
				if val := o.self.get(idx); val != nil ***REMOVED***
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
			idx := intToValue(k)
			if val := o.self.get(idx); val != nil ***REMOVED***
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
	lowerP := intToValue(lower)
	upperP := intToValue(upper)
	lowerValue := o.self.get(lowerP)
	upperValue := o.self.get(upperP)
	if lowerValue != nil && upperValue != nil ***REMOVED***
		o.self.put(lowerP, upperValue, true)
		o.self.put(upperP, lowerValue, true)
	***REMOVED*** else if lowerValue == nil && upperValue != nil ***REMOVED***
		o.self.put(lowerP, upperValue, true)
		o.self.delete(upperP, true)
	***REMOVED*** else if lowerValue != nil && upperValue == nil ***REMOVED***
		o.self.delete(lowerP, true)
		o.self.put(upperP, lowerValue, true)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_reverse_generic(o *Object, start int64) ***REMOVED***
	l := toLength(o.self.getStr("length"))
	middle := l / 2
	for lower := start; lower != middle; lower++ ***REMOVED***
		arrayproto_reverse_generic_step(o, lower, l-lower-1)
	***REMOVED***
***REMOVED***

func (r *Runtime) arrayproto_reverse(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	if a, ok := o.self.(*arrayObject); ok ***REMOVED***
		l := a.length
		middle := l / 2
		al := int64(len(a.values))
		for lower := int64(0); lower != middle; lower++ ***REMOVED***
			upper := l - lower - 1
			var lowerValue, upperValue Value
			if upper >= al || lower >= al ***REMOVED***
				goto bailout
			***REMOVED***
			lowerValue = a.values[lower]
			if lowerValue == nil ***REMOVED***
				goto bailout
			***REMOVED***
			if _, ok := lowerValue.(*valueProperty); ok ***REMOVED***
				goto bailout
			***REMOVED***
			upperValue = a.values[upper]
			if upperValue == nil ***REMOVED***
				goto bailout
			***REMOVED***
			if _, ok := upperValue.(*valueProperty); ok ***REMOVED***
				goto bailout
			***REMOVED***

			a.values[lower], a.values[upper] = upperValue, lowerValue
			continue
		bailout:
			arrayproto_reverse_generic_step(o, lower, upper)
		***REMOVED***
		//TODO: go arrays
	***REMOVED*** else ***REMOVED***
		r.arrayproto_reverse_generic(o, 0)
	***REMOVED***
	return o
***REMOVED***

func (r *Runtime) arrayproto_shift(call FunctionCall) Value ***REMOVED***
	o := call.This.ToObject(r)
	length := toLength(o.self.getStr("length"))
	if length == 0 ***REMOVED***
		o.self.putStr("length", intToValue(0), true)
		return _undefined
	***REMOVED***
	first := o.self.get(intToValue(0))
	for i := int64(1); i < length; i++ ***REMOVED***
		v := o.self.get(intToValue(i))
		if v != nil && v != _undefined ***REMOVED***
			o.self.put(intToValue(i-1), v, true)
		***REMOVED*** else ***REMOVED***
			o.self.delete(intToValue(i-1), true)
		***REMOVED***
	***REMOVED***

	lv := intToValue(length - 1)
	o.self.delete(lv, true)
	o.self.putStr("length", lv, true)

	return first
***REMOVED***

func (r *Runtime) array_isArray(call FunctionCall) Value ***REMOVED***
	if o, ok := call.Argument(0).(*Object); ok ***REMOVED***
		if isArray(o) ***REMOVED***
			return valueTrue
		***REMOVED***
	***REMOVED***
	return valueFalse
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
	o._putProp("pop", r.newNativeFunc(r.arrayproto_pop, nil, "pop", nil, 0), true, false, true)
	o._putProp("push", r.newNativeFunc(r.arrayproto_push, nil, "push", nil, 1), true, false, true)
	o._putProp("join", r.newNativeFunc(r.arrayproto_join, nil, "join", nil, 1), true, false, true)
	o._putProp("toString", r.newNativeFunc(r.arrayproto_toString, nil, "toString", nil, 0), true, false, true)
	o._putProp("toLocaleString", r.newNativeFunc(r.arrayproto_toLocaleString, nil, "toLocaleString", nil, 0), true, false, true)
	o._putProp("concat", r.newNativeFunc(r.arrayproto_concat, nil, "concat", nil, 1), true, false, true)
	o._putProp("reverse", r.newNativeFunc(r.arrayproto_reverse, nil, "reverse", nil, 0), true, false, true)
	o._putProp("shift", r.newNativeFunc(r.arrayproto_shift, nil, "shift", nil, 0), true, false, true)
	o._putProp("slice", r.newNativeFunc(r.arrayproto_slice, nil, "slice", nil, 2), true, false, true)
	o._putProp("sort", r.newNativeFunc(r.arrayproto_sort, nil, "sort", nil, 1), true, false, true)
	o._putProp("splice", r.newNativeFunc(r.arrayproto_splice, nil, "splice", nil, 2), true, false, true)
	o._putProp("unshift", r.newNativeFunc(r.arrayproto_unshift, nil, "unshift", nil, 1), true, false, true)
	o._putProp("indexOf", r.newNativeFunc(r.arrayproto_indexOf, nil, "indexOf", nil, 1), true, false, true)
	o._putProp("lastIndexOf", r.newNativeFunc(r.arrayproto_lastIndexOf, nil, "lastIndexOf", nil, 1), true, false, true)
	o._putProp("every", r.newNativeFunc(r.arrayproto_every, nil, "every", nil, 1), true, false, true)
	o._putProp("some", r.newNativeFunc(r.arrayproto_some, nil, "some", nil, 1), true, false, true)
	o._putProp("forEach", r.newNativeFunc(r.arrayproto_forEach, nil, "forEach", nil, 1), true, false, true)
	o._putProp("map", r.newNativeFunc(r.arrayproto_map, nil, "map", nil, 1), true, false, true)
	o._putProp("filter", r.newNativeFunc(r.arrayproto_filter, nil, "filter", nil, 1), true, false, true)
	o._putProp("reduce", r.newNativeFunc(r.arrayproto_reduce, nil, "reduce", nil, 1), true, false, true)
	o._putProp("reduceRight", r.newNativeFunc(r.arrayproto_reduceRight, nil, "reduceRight", nil, 1), true, false, true)

	return o
***REMOVED***

func (r *Runtime) createArray(val *Object) objectImpl ***REMOVED***
	o := r.newNativeFuncConstructObj(val, r.builtin_newArray, "Array", r.global.ArrayPrototype, 1)
	o._putProp("isArray", r.newNativeFunc(r.array_isArray, nil, "isArray", nil, 1), true, false, true)
	return o
***REMOVED***

func (r *Runtime) initArray() ***REMOVED***
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

func (ctx *arraySortCtx) sortCompare(x, y Value) int ***REMOVED***
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

	if ctx.compare != nil ***REMOVED***
		f := ctx.compare(FunctionCall***REMOVED***
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
	return strings.Compare(x.String(), y.String())
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
