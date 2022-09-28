package goja

import (
	"fmt"
	"math"
	"sort"
	"unsafe"

	"github.com/dop251/goja/unistring"
)

type typedArraySortCtx struct ***REMOVED***
	ta           *typedArrayObject
	compare      func(FunctionCall) Value
	needValidate bool
***REMOVED***

func (ctx *typedArraySortCtx) Len() int ***REMOVED***
	return ctx.ta.length
***REMOVED***

func (ctx *typedArraySortCtx) Less(i, j int) bool ***REMOVED***
	if ctx.needValidate ***REMOVED***
		ctx.ta.viewedArrayBuf.ensureNotDetached(true)
		ctx.needValidate = false
	***REMOVED***
	offset := ctx.ta.offset
	if ctx.compare != nil ***REMOVED***
		x := ctx.ta.typedArray.get(offset + i)
		y := ctx.ta.typedArray.get(offset + j)
		res := ctx.compare(FunctionCall***REMOVED***
			This:      _undefined,
			Arguments: []Value***REMOVED***x, y***REMOVED***,
		***REMOVED***).ToNumber()
		ctx.needValidate = true
		if i, ok := res.(valueInt); ok ***REMOVED***
			return i < 0
		***REMOVED***
		f := res.ToFloat()
		if f < 0 ***REMOVED***
			return true
		***REMOVED***
		if f > 0 ***REMOVED***
			return false
		***REMOVED***
		if math.Signbit(f) ***REMOVED***
			return true
		***REMOVED***
		return false
	***REMOVED***

	return ctx.ta.typedArray.less(offset+i, offset+j)
***REMOVED***

func (ctx *typedArraySortCtx) Swap(i, j int) ***REMOVED***
	if ctx.needValidate ***REMOVED***
		ctx.ta.viewedArrayBuf.ensureNotDetached(true)
		ctx.needValidate = false
	***REMOVED***
	offset := ctx.ta.offset
	ctx.ta.typedArray.swap(offset+i, offset+j)
***REMOVED***

func allocByteSlice(size int) (b []byte) ***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			panic(rangeError(fmt.Sprintf("Buffer size is too large: %d", size)))
		***REMOVED***
	***REMOVED***()
	if size < 0 ***REMOVED***
		panic(rangeError(fmt.Sprintf("Invalid buffer size: %d", size)))
	***REMOVED***
	b = make([]byte, size)
	return
***REMOVED***

func (r *Runtime) builtin_newArrayBuffer(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("ArrayBuffer"))
	***REMOVED***
	b := r._newArrayBuffer(r.getPrototypeFromCtor(newTarget, r.global.ArrayBuffer, r.global.ArrayBufferPrototype), nil)
	if len(args) > 0 ***REMOVED***
		b.data = allocByteSlice(r.toIndex(args[0]))
	***REMOVED***
	return b.val
***REMOVED***

func (r *Runtime) arrayBufferProto_getByteLength(call FunctionCall) Value ***REMOVED***
	o := r.toObject(call.This)
	if b, ok := o.self.(*arrayBufferObject); ok ***REMOVED***
		if b.ensureNotDetached(false) ***REMOVED***
			return intToValue(int64(len(b.data)))
		***REMOVED***
		return intToValue(0)
	***REMOVED***
	panic(r.NewTypeError("Object is not ArrayBuffer: %s", o))
***REMOVED***

func (r *Runtime) arrayBufferProto_slice(call FunctionCall) Value ***REMOVED***
	o := r.toObject(call.This)
	if b, ok := o.self.(*arrayBufferObject); ok ***REMOVED***
		l := int64(len(b.data))
		start := relToIdx(call.Argument(0).ToInteger(), l)
		var stop int64
		if arg := call.Argument(1); arg != _undefined ***REMOVED***
			stop = arg.ToInteger()
		***REMOVED*** else ***REMOVED***
			stop = l
		***REMOVED***
		stop = relToIdx(stop, l)
		newLen := max(stop-start, 0)
		ret := r.speciesConstructor(o, r.global.ArrayBuffer)([]Value***REMOVED***intToValue(newLen)***REMOVED***, nil)
		if ab, ok := ret.self.(*arrayBufferObject); ok ***REMOVED***
			if newLen > 0 ***REMOVED***
				b.ensureNotDetached(true)
				if ret == o ***REMOVED***
					panic(r.NewTypeError("Species constructor returned the same ArrayBuffer"))
				***REMOVED***
				if int64(len(ab.data)) < newLen ***REMOVED***
					panic(r.NewTypeError("Species constructor returned an ArrayBuffer that is too small: %d", len(ab.data)))
				***REMOVED***
				ab.ensureNotDetached(true)
				copy(ab.data, b.data[start:stop])
			***REMOVED***
			return ret
		***REMOVED***
		panic(r.NewTypeError("Species constructor did not return an ArrayBuffer: %s", ret.String()))
	***REMOVED***
	panic(r.NewTypeError("Object is not ArrayBuffer: %s", o))
***REMOVED***

func (r *Runtime) arrayBuffer_isView(call FunctionCall) Value ***REMOVED***
	if o, ok := call.Argument(0).(*Object); ok ***REMOVED***
		if _, ok := o.self.(*dataViewObject); ok ***REMOVED***
			return valueTrue
		***REMOVED***
		if _, ok := o.self.(*typedArrayObject); ok ***REMOVED***
			return valueTrue
		***REMOVED***
	***REMOVED***
	return valueFalse
***REMOVED***

func (r *Runtime) newDataView(args []Value, newTarget *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("DataView"))
	***REMOVED***
	proto := r.getPrototypeFromCtor(newTarget, r.global.DataView, r.global.DataViewPrototype)
	var bufArg Value
	if len(args) > 0 ***REMOVED***
		bufArg = args[0]
	***REMOVED***
	var buffer *arrayBufferObject
	if o, ok := bufArg.(*Object); ok ***REMOVED***
		if b, ok := o.self.(*arrayBufferObject); ok ***REMOVED***
			buffer = b
		***REMOVED***
	***REMOVED***
	if buffer == nil ***REMOVED***
		panic(r.NewTypeError("First argument to DataView constructor must be an ArrayBuffer"))
	***REMOVED***
	var byteOffset, byteLen int
	if len(args) > 1 ***REMOVED***
		offsetArg := nilSafe(args[1])
		byteOffset = r.toIndex(offsetArg)
		buffer.ensureNotDetached(true)
		if byteOffset > len(buffer.data) ***REMOVED***
			panic(r.newError(r.global.RangeError, "Start offset %s is outside the bounds of the buffer", offsetArg.String()))
		***REMOVED***
	***REMOVED***
	if len(args) > 2 && args[2] != nil && args[2] != _undefined ***REMOVED***
		byteLen = r.toIndex(args[2])
		if byteOffset+byteLen > len(buffer.data) ***REMOVED***
			panic(r.newError(r.global.RangeError, "Invalid DataView length %d", byteLen))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		byteLen = len(buffer.data) - byteOffset
	***REMOVED***
	o := &Object***REMOVED***runtime: r***REMOVED***
	b := &dataViewObject***REMOVED***
		baseObject: baseObject***REMOVED***
			class:      classObject,
			val:        o,
			prototype:  proto,
			extensible: true,
		***REMOVED***,
		viewedArrayBuf: buffer,
		byteOffset:     byteOffset,
		byteLen:        byteLen,
	***REMOVED***
	o.self = b
	b.init()
	return o
***REMOVED***

func (r *Runtime) dataViewProto_getBuffer(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		return dv.viewedArrayBuf.val
	***REMOVED***
	panic(r.NewTypeError("Method get DataView.prototype.buffer called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getByteLen(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		dv.viewedArrayBuf.ensureNotDetached(true)
		return intToValue(int64(dv.byteLen))
	***REMOVED***
	panic(r.NewTypeError("Method get DataView.prototype.byteLength called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getByteOffset(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		dv.viewedArrayBuf.ensureNotDetached(true)
		return intToValue(int64(dv.byteOffset))
	***REMOVED***
	panic(r.NewTypeError("Method get DataView.prototype.byteOffset called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getFloat32(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		return floatToValue(float64(dv.viewedArrayBuf.getFloat32(dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 4))))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getFloat32 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getFloat64(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		return floatToValue(dv.viewedArrayBuf.getFloat64(dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 8)))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getFloat64 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getInt8(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idx, _ := dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 1)
		return intToValue(int64(dv.viewedArrayBuf.getInt8(idx)))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getInt8 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getInt16(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		return intToValue(int64(dv.viewedArrayBuf.getInt16(dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 2))))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getInt16 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getInt32(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		return intToValue(int64(dv.viewedArrayBuf.getInt32(dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 4))))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getInt32 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getUint8(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idx, _ := dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 1)
		return intToValue(int64(dv.viewedArrayBuf.getUint8(idx)))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getUint8 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getUint16(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		return intToValue(int64(dv.viewedArrayBuf.getUint16(dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 2))))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getUint16 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_getUint32(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		return intToValue(int64(dv.viewedArrayBuf.getUint32(dv.getIdxAndByteOrder(r.toIndex(call.Argument(0)), call.Argument(1), 4))))
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.getUint32 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setFloat32(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := toFloat32(call.Argument(1))
		idx, bo := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 4)
		dv.viewedArrayBuf.setFloat32(idx, val, bo)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setFloat32 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setFloat64(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := call.Argument(1).ToFloat()
		idx, bo := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 8)
		dv.viewedArrayBuf.setFloat64(idx, val, bo)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setFloat64 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setInt8(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := toInt8(call.Argument(1))
		idx, _ := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 1)
		dv.viewedArrayBuf.setInt8(idx, val)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setInt8 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setInt16(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := toInt16(call.Argument(1))
		idx, bo := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 2)
		dv.viewedArrayBuf.setInt16(idx, val, bo)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setInt16 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setInt32(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := toInt32(call.Argument(1))
		idx, bo := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 4)
		dv.viewedArrayBuf.setInt32(idx, val, bo)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setInt32 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setUint8(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := toUint8(call.Argument(1))
		idx, _ := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 1)
		dv.viewedArrayBuf.setUint8(idx, val)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setUint8 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setUint16(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := toUint16(call.Argument(1))
		idx, bo := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 2)
		dv.viewedArrayBuf.setUint16(idx, val, bo)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setUint16 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) dataViewProto_setUint32(call FunctionCall) Value ***REMOVED***
	if dv, ok := r.toObject(call.This).self.(*dataViewObject); ok ***REMOVED***
		idxVal := r.toIndex(call.Argument(0))
		val := toUint32(call.Argument(1))
		idx, bo := dv.getIdxAndByteOrder(idxVal, call.Argument(2), 4)
		dv.viewedArrayBuf.setUint32(idx, val, bo)
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method DataView.prototype.setUint32 called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_getBuffer(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		return ta.viewedArrayBuf.val
	***REMOVED***
	panic(r.NewTypeError("Method get TypedArray.prototype.buffer called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_getByteLen(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		if ta.viewedArrayBuf.data == nil ***REMOVED***
			return _positiveZero
		***REMOVED***
		return intToValue(int64(ta.length) * int64(ta.elemSize))
	***REMOVED***
	panic(r.NewTypeError("Method get TypedArray.prototype.byteLength called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_getLength(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		if ta.viewedArrayBuf.data == nil ***REMOVED***
			return _positiveZero
		***REMOVED***
		return intToValue(int64(ta.length))
	***REMOVED***
	panic(r.NewTypeError("Method get TypedArray.prototype.length called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_getByteOffset(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		if ta.viewedArrayBuf.data == nil ***REMOVED***
			return _positiveZero
		***REMOVED***
		return intToValue(int64(ta.offset) * int64(ta.elemSize))
	***REMOVED***
	panic(r.NewTypeError("Method get TypedArray.prototype.byteOffset called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_copyWithin(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		l := int64(ta.length)
		var relEnd int64
		to := toIntStrict(relToIdx(call.Argument(0).ToInteger(), l))
		from := toIntStrict(relToIdx(call.Argument(1).ToInteger(), l))
		if end := call.Argument(2); end != _undefined ***REMOVED***
			relEnd = end.ToInteger()
		***REMOVED*** else ***REMOVED***
			relEnd = l
		***REMOVED***
		final := toIntStrict(relToIdx(relEnd, l))
		data := ta.viewedArrayBuf.data
		offset := ta.offset
		elemSize := ta.elemSize
		if final > from ***REMOVED***
			ta.viewedArrayBuf.ensureNotDetached(true)
			copy(data[(offset+to)*elemSize:], data[(offset+from)*elemSize:(offset+final)*elemSize])
		***REMOVED***
		return call.This
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.copyWithin called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_entries(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		return r.createArrayIterator(ta.val, iterationKindKeyValue)
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.entries called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_every(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		callbackFn := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, call.This***REMOVED***,
		***REMOVED***
		for k := 0; k < ta.length; k++ ***REMOVED***
			if ta.isValidIntegerIndex(k) ***REMOVED***
				fc.Arguments[0] = ta.typedArray.get(ta.offset + k)
			***REMOVED*** else ***REMOVED***
				fc.Arguments[0] = _undefined
			***REMOVED***
			fc.Arguments[1] = intToValue(int64(k))
			if !callbackFn(fc).ToBoolean() ***REMOVED***
				return valueFalse
			***REMOVED***
		***REMOVED***
		return valueTrue

	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.every called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_fill(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		l := int64(ta.length)
		k := toIntStrict(relToIdx(call.Argument(1).ToInteger(), l))
		var relEnd int64
		if endArg := call.Argument(2); endArg != _undefined ***REMOVED***
			relEnd = endArg.ToInteger()
		***REMOVED*** else ***REMOVED***
			relEnd = l
		***REMOVED***
		final := toIntStrict(relToIdx(relEnd, l))
		value := ta.typedArray.toRaw(call.Argument(0))
		ta.viewedArrayBuf.ensureNotDetached(true)
		for ; k < final; k++ ***REMOVED***
			ta.typedArray.setRaw(ta.offset+k, value)
		***REMOVED***
		return call.This
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.fill called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_filter(call FunctionCall) Value ***REMOVED***
	o := r.toObject(call.This)
	if ta, ok := o.self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		callbackFn := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, call.This***REMOVED***,
		***REMOVED***
		buf := make([]byte, 0, ta.length*ta.elemSize)
		captured := 0
		rawVal := make([]byte, ta.elemSize)
		for k := 0; k < ta.length; k++ ***REMOVED***
			if ta.isValidIntegerIndex(k) ***REMOVED***
				fc.Arguments[0] = ta.typedArray.get(ta.offset + k)
				i := (ta.offset + k) * ta.elemSize
				copy(rawVal, ta.viewedArrayBuf.data[i:])
			***REMOVED*** else ***REMOVED***
				fc.Arguments[0] = _undefined
				for i := range rawVal ***REMOVED***
					rawVal[i] = 0
				***REMOVED***
			***REMOVED***
			fc.Arguments[1] = intToValue(int64(k))
			if callbackFn(fc).ToBoolean() ***REMOVED***
				buf = append(buf, rawVal...)
				captured++
			***REMOVED***
		***REMOVED***
		c := r.speciesConstructorObj(o, ta.defaultCtor)
		ab := r._newArrayBuffer(r.global.ArrayBufferPrototype, nil)
		ab.data = buf
		kept := r.toConstructor(ta.defaultCtor)([]Value***REMOVED***ab.val***REMOVED***, ta.defaultCtor)
		if c == ta.defaultCtor ***REMOVED***
			return kept
		***REMOVED*** else ***REMOVED***
			ret := r.typedArrayCreate(c, intToValue(int64(captured)))
			keptTa := kept.self.(*typedArrayObject)
			for i := 0; i < captured; i++ ***REMOVED***
				ret.typedArray.set(i, keptTa.typedArray.get(keptTa.offset+i))
			***REMOVED***
			return ret.val
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.filter called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_find(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		predicate := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, call.This***REMOVED***,
		***REMOVED***
		for k := 0; k < ta.length; k++ ***REMOVED***
			var val Value
			if ta.isValidIntegerIndex(k) ***REMOVED***
				val = ta.typedArray.get(ta.offset + k)
			***REMOVED***
			fc.Arguments[0] = val
			fc.Arguments[1] = intToValue(int64(k))
			if predicate(fc).ToBoolean() ***REMOVED***
				return val
			***REMOVED***
		***REMOVED***
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.find called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_findIndex(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		predicate := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, call.This***REMOVED***,
		***REMOVED***
		for k := 0; k < ta.length; k++ ***REMOVED***
			if ta.isValidIntegerIndex(k) ***REMOVED***
				fc.Arguments[0] = ta.typedArray.get(ta.offset + k)
			***REMOVED*** else ***REMOVED***
				fc.Arguments[0] = _undefined
			***REMOVED***
			fc.Arguments[1] = intToValue(int64(k))
			if predicate(fc).ToBoolean() ***REMOVED***
				return fc.Arguments[1]
			***REMOVED***
		***REMOVED***
		return intToValue(-1)
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.findIndex called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_forEach(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		callbackFn := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, call.This***REMOVED***,
		***REMOVED***
		for k := 0; k < ta.length; k++ ***REMOVED***
			var val Value
			if ta.isValidIntegerIndex(k) ***REMOVED***
				val = ta.typedArray.get(ta.offset + k)
			***REMOVED***
			fc.Arguments[0] = val
			fc.Arguments[1] = intToValue(int64(k))
			callbackFn(fc)
		***REMOVED***
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.forEach called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_includes(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		length := int64(ta.length)
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
		startIdx := toIntStrict(n)
		if !ta.viewedArrayBuf.ensureNotDetached(false) ***REMOVED***
			if searchElement == _undefined && startIdx < ta.length ***REMOVED***
				return valueTrue
			***REMOVED***
			return valueFalse
		***REMOVED***
		if ta.typedArray.typeMatch(searchElement) ***REMOVED***
			se := ta.typedArray.toRaw(searchElement)
			for k := startIdx; k < ta.length; k++ ***REMOVED***
				if ta.typedArray.getRaw(ta.offset+k) == se ***REMOVED***
					return valueTrue
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return valueFalse
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.includes called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_at(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		idx := call.Argument(0).ToInteger()
		length := int64(ta.length)
		if idx < 0 ***REMOVED***
			idx = length + idx
		***REMOVED***
		if idx >= length || idx < 0 ***REMOVED***
			return _undefined
		***REMOVED***
		if ta.viewedArrayBuf.ensureNotDetached(false) ***REMOVED***
			return ta.typedArray.get(ta.offset + int(idx))
		***REMOVED***
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.at called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_indexOf(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		length := int64(ta.length)
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

		if ta.viewedArrayBuf.ensureNotDetached(false) ***REMOVED***
			searchElement := call.Argument(0)
			if searchElement == _negativeZero ***REMOVED***
				searchElement = _positiveZero
			***REMOVED***
			if !IsNaN(searchElement) && ta.typedArray.typeMatch(searchElement) ***REMOVED***
				se := ta.typedArray.toRaw(searchElement)
				for k := toIntStrict(n); k < ta.length; k++ ***REMOVED***
					if ta.typedArray.getRaw(ta.offset+k) == se ***REMOVED***
						return intToValue(int64(k))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return intToValue(-1)
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.indexOf called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_join(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		s := call.Argument(0)
		var sep valueString
		if s != _undefined ***REMOVED***
			sep = s.toString()
		***REMOVED*** else ***REMOVED***
			sep = asciiString(",")
		***REMOVED***
		l := ta.length
		if l == 0 ***REMOVED***
			return stringEmpty
		***REMOVED***

		var buf valueStringBuilder

		var element0 Value
		if ta.isValidIntegerIndex(0) ***REMOVED***
			element0 = ta.typedArray.get(ta.offset + 0)
		***REMOVED***
		if element0 != nil && element0 != _undefined && element0 != _null ***REMOVED***
			buf.WriteString(element0.toString())
		***REMOVED***

		for i := 1; i < l; i++ ***REMOVED***
			buf.WriteString(sep)
			if ta.isValidIntegerIndex(i) ***REMOVED***
				element := ta.typedArray.get(ta.offset + i)
				if element != nil && element != _undefined && element != _null ***REMOVED***
					buf.WriteString(element.toString())
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return buf.String()
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.join called on incompatible receiver"))
***REMOVED***

func (r *Runtime) typedArrayProto_keys(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		return r.createArrayIterator(ta.val, iterationKindKey)
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.keys called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_lastIndexOf(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		length := int64(ta.length)
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
				if fromIndex < 0 ***REMOVED***
					fromIndex = -1 // prevent underflow in toIntStrict() on 32-bit platforms
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if ta.viewedArrayBuf.ensureNotDetached(false) ***REMOVED***
			searchElement := call.Argument(0)
			if searchElement == _negativeZero ***REMOVED***
				searchElement = _positiveZero
			***REMOVED***
			if !IsNaN(searchElement) && ta.typedArray.typeMatch(searchElement) ***REMOVED***
				se := ta.typedArray.toRaw(searchElement)
				for k := toIntStrict(fromIndex); k >= 0; k-- ***REMOVED***
					if ta.typedArray.getRaw(ta.offset+k) == se ***REMOVED***
						return intToValue(int64(k))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return intToValue(-1)
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.lastIndexOf called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_map(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		callbackFn := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, call.This***REMOVED***,
		***REMOVED***
		dst := r.typedArraySpeciesCreate(ta, []Value***REMOVED***intToValue(int64(ta.length))***REMOVED***)
		for i := 0; i < ta.length; i++ ***REMOVED***
			if ta.isValidIntegerIndex(i) ***REMOVED***
				fc.Arguments[0] = ta.typedArray.get(ta.offset + i)
			***REMOVED*** else ***REMOVED***
				fc.Arguments[0] = _undefined
			***REMOVED***
			fc.Arguments[1] = intToValue(int64(i))
			dst.typedArray.set(i, callbackFn(fc))
		***REMOVED***
		return dst.val
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.map called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_reduce(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		callbackFn := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      _undefined,
			Arguments: []Value***REMOVED***nil, nil, nil, call.This***REMOVED***,
		***REMOVED***
		k := 0
		if len(call.Arguments) >= 2 ***REMOVED***
			fc.Arguments[0] = call.Argument(1)
		***REMOVED*** else ***REMOVED***
			if ta.length > 0 ***REMOVED***
				fc.Arguments[0] = ta.typedArray.get(ta.offset + 0)
				k = 1
			***REMOVED***
		***REMOVED***
		if fc.Arguments[0] == nil ***REMOVED***
			panic(r.NewTypeError("Reduce of empty array with no initial value"))
		***REMOVED***
		for ; k < ta.length; k++ ***REMOVED***
			if ta.isValidIntegerIndex(k) ***REMOVED***
				fc.Arguments[1] = ta.typedArray.get(ta.offset + k)
			***REMOVED*** else ***REMOVED***
				fc.Arguments[1] = _undefined
			***REMOVED***
			idx := valueInt(k)
			fc.Arguments[2] = idx
			fc.Arguments[0] = callbackFn(fc)
		***REMOVED***
		return fc.Arguments[0]
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.reduce called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_reduceRight(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		callbackFn := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      _undefined,
			Arguments: []Value***REMOVED***nil, nil, nil, call.This***REMOVED***,
		***REMOVED***
		k := ta.length - 1
		if len(call.Arguments) >= 2 ***REMOVED***
			fc.Arguments[0] = call.Argument(1)
		***REMOVED*** else ***REMOVED***
			if k >= 0 ***REMOVED***
				fc.Arguments[0] = ta.typedArray.get(ta.offset + k)
				k--
			***REMOVED***
		***REMOVED***
		if fc.Arguments[0] == nil ***REMOVED***
			panic(r.NewTypeError("Reduce of empty array with no initial value"))
		***REMOVED***
		for ; k >= 0; k-- ***REMOVED***
			if ta.isValidIntegerIndex(k) ***REMOVED***
				fc.Arguments[1] = ta.typedArray.get(ta.offset + k)
			***REMOVED*** else ***REMOVED***
				fc.Arguments[1] = _undefined
			***REMOVED***
			idx := valueInt(k)
			fc.Arguments[2] = idx
			fc.Arguments[0] = callbackFn(fc)
		***REMOVED***
		return fc.Arguments[0]
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.reduceRight called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_reverse(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		l := ta.length
		middle := l / 2
		for lower := 0; lower != middle; lower++ ***REMOVED***
			upper := l - lower - 1
			ta.typedArray.swap(ta.offset+lower, ta.offset+upper)
		***REMOVED***

		return call.This
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.reverse called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_set(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		srcObj := call.Argument(0).ToObject(r)
		targetOffset := toIntStrict(call.Argument(1).ToInteger())
		if targetOffset < 0 ***REMOVED***
			panic(r.newError(r.global.RangeError, "offset should be >= 0"))
		***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		targetLen := ta.length
		if src, ok := srcObj.self.(*typedArrayObject); ok ***REMOVED***
			src.viewedArrayBuf.ensureNotDetached(true)
			srcLen := src.length
			if x := srcLen + targetOffset; x < 0 || x > targetLen ***REMOVED***
				panic(r.newError(r.global.RangeError, "Source is too large"))
			***REMOVED***
			if src.defaultCtor == ta.defaultCtor ***REMOVED***
				copy(ta.viewedArrayBuf.data[(ta.offset+targetOffset)*ta.elemSize:],
					src.viewedArrayBuf.data[src.offset*src.elemSize:(src.offset+srcLen)*src.elemSize])
			***REMOVED*** else ***REMOVED***
				curSrc := uintptr(unsafe.Pointer(&src.viewedArrayBuf.data[src.offset*src.elemSize]))
				endSrc := curSrc + uintptr(srcLen*src.elemSize)
				curDst := uintptr(unsafe.Pointer(&ta.viewedArrayBuf.data[(ta.offset+targetOffset)*ta.elemSize]))
				dstOffset := ta.offset + targetOffset
				srcOffset := src.offset
				if ta.elemSize == src.elemSize ***REMOVED***
					if curDst <= curSrc || curDst >= endSrc ***REMOVED***
						for i := 0; i < srcLen; i++ ***REMOVED***
							ta.typedArray.set(dstOffset+i, src.typedArray.get(srcOffset+i))
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						for i := srcLen - 1; i >= 0; i-- ***REMOVED***
							ta.typedArray.set(dstOffset+i, src.typedArray.get(srcOffset+i))
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					x := int(curDst-curSrc) / (src.elemSize - ta.elemSize)
					if x < 0 ***REMOVED***
						x = 0
					***REMOVED*** else if x > srcLen ***REMOVED***
						x = srcLen
					***REMOVED***
					if ta.elemSize < src.elemSize ***REMOVED***
						for i := x; i < srcLen; i++ ***REMOVED***
							ta.typedArray.set(dstOffset+i, src.typedArray.get(srcOffset+i))
						***REMOVED***
						for i := x - 1; i >= 0; i-- ***REMOVED***
							ta.typedArray.set(dstOffset+i, src.typedArray.get(srcOffset+i))
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						for i := 0; i < x; i++ ***REMOVED***
							ta.typedArray.set(dstOffset+i, src.typedArray.get(srcOffset+i))
						***REMOVED***
						for i := srcLen - 1; i >= x; i-- ***REMOVED***
							ta.typedArray.set(dstOffset+i, src.typedArray.get(srcOffset+i))
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			targetLen := ta.length
			srcLen := toIntStrict(toLength(srcObj.self.getStr("length", nil)))
			if x := srcLen + targetOffset; x < 0 || x > targetLen ***REMOVED***
				panic(r.newError(r.global.RangeError, "Source is too large"))
			***REMOVED***
			for i := 0; i < srcLen; i++ ***REMOVED***
				val := nilSafe(srcObj.self.getIdx(valueInt(i), nil))
				ta.viewedArrayBuf.ensureNotDetached(true)
				if ta.isValidIntegerIndex(i) ***REMOVED***
					ta.typedArray.set(targetOffset+i, val)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return _undefined
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.set called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_slice(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		length := int64(ta.length)
		start := toIntStrict(relToIdx(call.Argument(0).ToInteger(), length))
		var e int64
		if endArg := call.Argument(1); endArg != _undefined ***REMOVED***
			e = endArg.ToInteger()
		***REMOVED*** else ***REMOVED***
			e = length
		***REMOVED***
		end := toIntStrict(relToIdx(e, length))

		count := end - start
		if count < 0 ***REMOVED***
			count = 0
		***REMOVED***
		dst := r.typedArraySpeciesCreate(ta, []Value***REMOVED***intToValue(int64(count))***REMOVED***)
		if dst.defaultCtor == ta.defaultCtor ***REMOVED***
			if count > 0 ***REMOVED***
				ta.viewedArrayBuf.ensureNotDetached(true)
				offset := ta.offset
				elemSize := ta.elemSize
				copy(dst.viewedArrayBuf.data, ta.viewedArrayBuf.data[(offset+start)*elemSize:(offset+start+count)*elemSize])
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for i := 0; i < count; i++ ***REMOVED***
				ta.viewedArrayBuf.ensureNotDetached(true)
				dst.typedArray.set(i, ta.typedArray.get(ta.offset+start+i))
			***REMOVED***
		***REMOVED***
		return dst.val
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.slice called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_some(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		callbackFn := r.toCallable(call.Argument(0))
		fc := FunctionCall***REMOVED***
			This:      call.Argument(1),
			Arguments: []Value***REMOVED***nil, nil, call.This***REMOVED***,
		***REMOVED***
		for k := 0; k < ta.length; k++ ***REMOVED***
			if ta.isValidIntegerIndex(k) ***REMOVED***
				fc.Arguments[0] = ta.typedArray.get(ta.offset + k)
			***REMOVED*** else ***REMOVED***
				fc.Arguments[0] = _undefined
			***REMOVED***
			fc.Arguments[1] = intToValue(int64(k))
			if callbackFn(fc).ToBoolean() ***REMOVED***
				return valueTrue
			***REMOVED***
		***REMOVED***
		return valueFalse
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.some called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_sort(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		var compareFn func(FunctionCall) Value

		if arg := call.Argument(0); arg != _undefined ***REMOVED***
			compareFn = r.toCallable(arg)
		***REMOVED***

		ctx := typedArraySortCtx***REMOVED***
			ta:      ta,
			compare: compareFn,
		***REMOVED***

		sort.Stable(&ctx)
		return call.This
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.sort called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_subarray(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		l := int64(ta.length)
		beginIdx := relToIdx(call.Argument(0).ToInteger(), l)
		var relEnd int64
		if endArg := call.Argument(1); endArg != _undefined ***REMOVED***
			relEnd = endArg.ToInteger()
		***REMOVED*** else ***REMOVED***
			relEnd = l
		***REMOVED***
		endIdx := relToIdx(relEnd, l)
		newLen := max(endIdx-beginIdx, 0)
		return r.typedArraySpeciesCreate(ta, []Value***REMOVED***ta.viewedArrayBuf.val,
			intToValue((int64(ta.offset) + beginIdx) * int64(ta.elemSize)),
			intToValue(newLen),
		***REMOVED***).val
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.subarray called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_toLocaleString(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		length := ta.length
		var buf valueStringBuilder
		for i := 0; i < length; i++ ***REMOVED***
			ta.viewedArrayBuf.ensureNotDetached(true)
			if i > 0 ***REMOVED***
				buf.WriteRune(',')
			***REMOVED***
			item := ta.typedArray.get(ta.offset + i)
			r.writeItemLocaleString(item, &buf)
		***REMOVED***
		return buf.String()
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.toLocaleString called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_values(call FunctionCall) Value ***REMOVED***
	if ta, ok := r.toObject(call.This).self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		return r.createArrayIterator(ta.val, iterationKindValue)
	***REMOVED***
	panic(r.NewTypeError("Method TypedArray.prototype.values called on incompatible receiver %s", r.objectproto_toString(FunctionCall***REMOVED***This: call.This***REMOVED***)))
***REMOVED***

func (r *Runtime) typedArrayProto_toStringTag(call FunctionCall) Value ***REMOVED***
	if obj, ok := call.This.(*Object); ok ***REMOVED***
		if ta, ok := obj.self.(*typedArrayObject); ok ***REMOVED***
			return nilSafe(ta.defaultCtor.self.getStr("name", nil))
		***REMOVED***
	***REMOVED***

	return _undefined
***REMOVED***

func (r *Runtime) newTypedArray([]Value, *Object) *Object ***REMOVED***
	panic(r.NewTypeError("Abstract class TypedArray not directly constructable"))
***REMOVED***

func (r *Runtime) typedArray_from(call FunctionCall) Value ***REMOVED***
	c := r.toObject(call.This)
	var mapFc func(call FunctionCall) Value
	thisValue := call.Argument(2)
	if mapFn := call.Argument(1); mapFn != _undefined ***REMOVED***
		mapFc = r.toCallable(mapFn)
	***REMOVED***
	source := r.toObject(call.Argument(0))
	usingIter := toMethod(source.self.getSym(SymIterator, nil))
	if usingIter != nil ***REMOVED***
		values := r.iterableToList(source, usingIter)
		ta := r.typedArrayCreate(c, intToValue(int64(len(values))))
		if mapFc == nil ***REMOVED***
			for idx, val := range values ***REMOVED***
				ta.typedArray.set(idx, val)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			fc := FunctionCall***REMOVED***
				This:      thisValue,
				Arguments: []Value***REMOVED***nil, nil***REMOVED***,
			***REMOVED***
			for idx, val := range values ***REMOVED***
				fc.Arguments[0], fc.Arguments[1] = val, intToValue(int64(idx))
				val = mapFc(fc)
				ta.typedArray.set(idx, val)
			***REMOVED***
		***REMOVED***
		return ta.val
	***REMOVED***
	length := toIntStrict(toLength(source.self.getStr("length", nil)))
	ta := r.typedArrayCreate(c, intToValue(int64(length)))
	if mapFc == nil ***REMOVED***
		for i := 0; i < length; i++ ***REMOVED***
			ta.typedArray.set(i, nilSafe(source.self.getIdx(valueInt(i), nil)))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      thisValue,
			Arguments: []Value***REMOVED***nil, nil***REMOVED***,
		***REMOVED***
		for i := 0; i < length; i++ ***REMOVED***
			idx := valueInt(i)
			fc.Arguments[0], fc.Arguments[1] = source.self.getIdx(idx, nil), idx
			ta.typedArray.set(i, mapFc(fc))
		***REMOVED***
	***REMOVED***
	return ta.val
***REMOVED***

func (r *Runtime) typedArray_of(call FunctionCall) Value ***REMOVED***
	ta := r.typedArrayCreate(r.toObject(call.This), intToValue(int64(len(call.Arguments))))
	for i, val := range call.Arguments ***REMOVED***
		ta.typedArray.set(i, val)
	***REMOVED***
	return ta.val
***REMOVED***

func (r *Runtime) allocateTypedArray(newTarget *Object, length int, taCtor typedArrayObjectCtor, proto *Object) *typedArrayObject ***REMOVED***
	buf := r._newArrayBuffer(r.global.ArrayBufferPrototype, nil)
	ta := taCtor(buf, 0, length, r.getPrototypeFromCtor(newTarget, nil, proto))
	if length > 0 ***REMOVED***
		buf.data = allocByteSlice(length * ta.elemSize)
	***REMOVED***
	return ta
***REMOVED***

func (r *Runtime) typedArraySpeciesCreate(ta *typedArrayObject, args []Value) *typedArrayObject ***REMOVED***
	return r.typedArrayCreate(r.speciesConstructorObj(ta.val, ta.defaultCtor), args...)
***REMOVED***

func (r *Runtime) typedArrayCreate(ctor *Object, args ...Value) *typedArrayObject ***REMOVED***
	o := r.toConstructor(ctor)(args, ctor)
	if ta, ok := o.self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
		if len(args) == 1 ***REMOVED***
			if l, ok := args[0].(valueInt); ok ***REMOVED***
				if ta.length < int(l) ***REMOVED***
					panic(r.NewTypeError("Derived TypedArray constructor created an array which was too small"))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return ta
	***REMOVED***
	panic(r.NewTypeError("Invalid TypedArray: %s", o))
***REMOVED***

func (r *Runtime) typedArrayFrom(ctor, items *Object, mapFn, thisValue Value, taCtor typedArrayObjectCtor, proto *Object) *Object ***REMOVED***
	var mapFc func(call FunctionCall) Value
	if mapFn != nil ***REMOVED***
		mapFc = r.toCallable(mapFn)
		if thisValue == nil ***REMOVED***
			thisValue = _undefined
		***REMOVED***
	***REMOVED***
	usingIter := toMethod(items.self.getSym(SymIterator, nil))
	if usingIter != nil ***REMOVED***
		values := r.iterableToList(items, usingIter)
		ta := r.allocateTypedArray(ctor, len(values), taCtor, proto)
		if mapFc == nil ***REMOVED***
			for idx, val := range values ***REMOVED***
				ta.typedArray.set(idx, val)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			fc := FunctionCall***REMOVED***
				This:      thisValue,
				Arguments: []Value***REMOVED***nil, nil***REMOVED***,
			***REMOVED***
			for idx, val := range values ***REMOVED***
				fc.Arguments[0], fc.Arguments[1] = val, intToValue(int64(idx))
				val = mapFc(fc)
				ta.typedArray.set(idx, val)
			***REMOVED***
		***REMOVED***
		return ta.val
	***REMOVED***
	length := toIntStrict(toLength(items.self.getStr("length", nil)))
	ta := r.allocateTypedArray(ctor, length, taCtor, proto)
	if mapFc == nil ***REMOVED***
		for i := 0; i < length; i++ ***REMOVED***
			ta.typedArray.set(i, nilSafe(items.self.getIdx(valueInt(i), nil)))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fc := FunctionCall***REMOVED***
			This:      thisValue,
			Arguments: []Value***REMOVED***nil, nil***REMOVED***,
		***REMOVED***
		for i := 0; i < length; i++ ***REMOVED***
			idx := valueInt(i)
			fc.Arguments[0], fc.Arguments[1] = items.self.getIdx(idx, nil), idx
			ta.typedArray.set(i, mapFc(fc))
		***REMOVED***
	***REMOVED***
	return ta.val
***REMOVED***

func (r *Runtime) _newTypedArrayFromArrayBuffer(ab *arrayBufferObject, args []Value, newTarget *Object, taCtor typedArrayObjectCtor, proto *Object) *Object ***REMOVED***
	ta := taCtor(ab, 0, 0, r.getPrototypeFromCtor(newTarget, nil, proto))
	var byteOffset int
	if len(args) > 1 && args[1] != nil && args[1] != _undefined ***REMOVED***
		byteOffset = r.toIndex(args[1])
		if byteOffset%ta.elemSize != 0 ***REMOVED***
			panic(r.newError(r.global.RangeError, "Start offset of %s should be a multiple of %d", newTarget.self.getStr("name", nil), ta.elemSize))
		***REMOVED***
	***REMOVED***
	var length int
	if len(args) > 2 && args[2] != nil && args[2] != _undefined ***REMOVED***
		length = r.toIndex(args[2])
		ab.ensureNotDetached(true)
		if byteOffset+length*ta.elemSize > len(ab.data) ***REMOVED***
			panic(r.newError(r.global.RangeError, "Invalid typed array length: %d", length))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ab.ensureNotDetached(true)
		if len(ab.data)%ta.elemSize != 0 ***REMOVED***
			panic(r.newError(r.global.RangeError, "Byte length of %s should be a multiple of %d", newTarget.self.getStr("name", nil), ta.elemSize))
		***REMOVED***
		length = (len(ab.data) - byteOffset) / ta.elemSize
		if length < 0 ***REMOVED***
			panic(r.newError(r.global.RangeError, "Start offset %d is outside the bounds of the buffer", byteOffset))
		***REMOVED***
	***REMOVED***
	ta.offset = byteOffset / ta.elemSize
	ta.length = length
	return ta.val
***REMOVED***

func (r *Runtime) _newTypedArrayFromTypedArray(src *typedArrayObject, newTarget *Object, taCtor typedArrayObjectCtor, proto *Object) *Object ***REMOVED***
	dst := r.allocateTypedArray(newTarget, 0, taCtor, proto)
	src.viewedArrayBuf.ensureNotDetached(true)
	l := src.length

	dst.viewedArrayBuf.prototype = r.getPrototypeFromCtor(r.speciesConstructorObj(src.viewedArrayBuf.val, r.global.ArrayBuffer), r.global.ArrayBuffer, r.global.ArrayBufferPrototype)
	dst.viewedArrayBuf.data = allocByteSlice(toIntStrict(int64(l) * int64(dst.elemSize)))
	src.viewedArrayBuf.ensureNotDetached(true)
	if src.defaultCtor == dst.defaultCtor ***REMOVED***
		copy(dst.viewedArrayBuf.data, src.viewedArrayBuf.data[src.offset*src.elemSize:])
		dst.length = src.length
		return dst.val
	***REMOVED***
	dst.length = l
	for i := 0; i < l; i++ ***REMOVED***
		dst.typedArray.set(i, src.typedArray.get(src.offset+i))
	***REMOVED***
	return dst.val
***REMOVED***

func (r *Runtime) _newTypedArray(args []Value, newTarget *Object, taCtor typedArrayObjectCtor, proto *Object) *Object ***REMOVED***
	if newTarget == nil ***REMOVED***
		panic(r.needNew("TypedArray"))
	***REMOVED***
	if len(args) > 0 ***REMOVED***
		if obj, ok := args[0].(*Object); ok ***REMOVED***
			switch o := obj.self.(type) ***REMOVED***
			case *arrayBufferObject:
				return r._newTypedArrayFromArrayBuffer(o, args, newTarget, taCtor, proto)
			case *typedArrayObject:
				return r._newTypedArrayFromTypedArray(o, newTarget, taCtor, proto)
			default:
				return r.typedArrayFrom(newTarget, obj, nil, nil, taCtor, proto)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	var l int
	if len(args) > 0 ***REMOVED***
		if arg0 := args[0]; arg0 != nil ***REMOVED***
			l = r.toIndex(arg0)
		***REMOVED***
	***REMOVED***
	return r.allocateTypedArray(newTarget, l, taCtor, proto).val
***REMOVED***

func (r *Runtime) newUint8Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newUint8ArrayObject, proto)
***REMOVED***

func (r *Runtime) newUint8ClampedArray(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newUint8ClampedArrayObject, proto)
***REMOVED***

func (r *Runtime) newInt8Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newInt8ArrayObject, proto)
***REMOVED***

func (r *Runtime) newUint16Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newUint16ArrayObject, proto)
***REMOVED***

func (r *Runtime) newInt16Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newInt16ArrayObject, proto)
***REMOVED***

func (r *Runtime) newUint32Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newUint32ArrayObject, proto)
***REMOVED***

func (r *Runtime) newInt32Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newInt32ArrayObject, proto)
***REMOVED***

func (r *Runtime) newFloat32Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newFloat32ArrayObject, proto)
***REMOVED***

func (r *Runtime) newFloat64Array(args []Value, newTarget, proto *Object) *Object ***REMOVED***
	return r._newTypedArray(args, newTarget, r.newFloat64ArrayObject, proto)
***REMOVED***

func (r *Runtime) createArrayBufferProto(val *Object) objectImpl ***REMOVED***
	b := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)
	byteLengthProp := &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.arrayBufferProto_getByteLength, nil, "get byteLength", nil, 0),
	***REMOVED***
	b._put("byteLength", byteLengthProp)
	b._putProp("constructor", r.global.ArrayBuffer, true, false, true)
	b._putProp("slice", r.newNativeFunc(r.arrayBufferProto_slice, nil, "slice", nil, 2), true, false, true)
	b._putSym(SymToStringTag, valueProp(asciiString("ArrayBuffer"), false, false, true))
	return b
***REMOVED***

func (r *Runtime) createArrayBuffer(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.builtin_newArrayBuffer, r.global.ArrayBufferPrototype, "ArrayBuffer", 1)
	o._putProp("isView", r.newNativeFunc(r.arrayBuffer_isView, nil, "isView", nil, 1), true, false, true)
	r.putSpeciesReturnThis(o)

	return o
***REMOVED***

func (r *Runtime) createDataViewProto(val *Object) objectImpl ***REMOVED***
	b := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)
	b._put("buffer", &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.dataViewProto_getBuffer, nil, "get buffer", nil, 0),
	***REMOVED***)
	b._put("byteLength", &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.dataViewProto_getByteLen, nil, "get byteLength", nil, 0),
	***REMOVED***)
	b._put("byteOffset", &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.dataViewProto_getByteOffset, nil, "get byteOffset", nil, 0),
	***REMOVED***)
	b._putProp("constructor", r.global.DataView, true, false, true)
	b._putProp("getFloat32", r.newNativeFunc(r.dataViewProto_getFloat32, nil, "getFloat32", nil, 1), true, false, true)
	b._putProp("getFloat64", r.newNativeFunc(r.dataViewProto_getFloat64, nil, "getFloat64", nil, 1), true, false, true)
	b._putProp("getInt8", r.newNativeFunc(r.dataViewProto_getInt8, nil, "getInt8", nil, 1), true, false, true)
	b._putProp("getInt16", r.newNativeFunc(r.dataViewProto_getInt16, nil, "getInt16", nil, 1), true, false, true)
	b._putProp("getInt32", r.newNativeFunc(r.dataViewProto_getInt32, nil, "getInt32", nil, 1), true, false, true)
	b._putProp("getUint8", r.newNativeFunc(r.dataViewProto_getUint8, nil, "getUint8", nil, 1), true, false, true)
	b._putProp("getUint16", r.newNativeFunc(r.dataViewProto_getUint16, nil, "getUint16", nil, 1), true, false, true)
	b._putProp("getUint32", r.newNativeFunc(r.dataViewProto_getUint32, nil, "getUint32", nil, 1), true, false, true)
	b._putProp("setFloat32", r.newNativeFunc(r.dataViewProto_setFloat32, nil, "setFloat32", nil, 2), true, false, true)
	b._putProp("setFloat64", r.newNativeFunc(r.dataViewProto_setFloat64, nil, "setFloat64", nil, 2), true, false, true)
	b._putProp("setInt8", r.newNativeFunc(r.dataViewProto_setInt8, nil, "setInt8", nil, 2), true, false, true)
	b._putProp("setInt16", r.newNativeFunc(r.dataViewProto_setInt16, nil, "setInt16", nil, 2), true, false, true)
	b._putProp("setInt32", r.newNativeFunc(r.dataViewProto_setInt32, nil, "setInt32", nil, 2), true, false, true)
	b._putProp("setUint8", r.newNativeFunc(r.dataViewProto_setUint8, nil, "setUint8", nil, 2), true, false, true)
	b._putProp("setUint16", r.newNativeFunc(r.dataViewProto_setUint16, nil, "setUint16", nil, 2), true, false, true)
	b._putProp("setUint32", r.newNativeFunc(r.dataViewProto_setUint32, nil, "setUint32", nil, 2), true, false, true)
	b._putSym(SymToStringTag, valueProp(asciiString("DataView"), false, false, true))

	return b
***REMOVED***

func (r *Runtime) createDataView(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.newDataView, r.global.DataViewPrototype, "DataView", 1)
	return o
***REMOVED***

func (r *Runtime) createTypedArrayProto(val *Object) objectImpl ***REMOVED***
	b := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)
	b._put("buffer", &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.typedArrayProto_getBuffer, nil, "get buffer", nil, 0),
	***REMOVED***)
	b._put("byteLength", &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.typedArrayProto_getByteLen, nil, "get byteLength", nil, 0),
	***REMOVED***)
	b._put("byteOffset", &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.typedArrayProto_getByteOffset, nil, "get byteOffset", nil, 0),
	***REMOVED***)
	b._putProp("at", r.newNativeFunc(r.typedArrayProto_at, nil, "at", nil, 1), true, false, true)
	b._putProp("constructor", r.global.TypedArray, true, false, true)
	b._putProp("copyWithin", r.newNativeFunc(r.typedArrayProto_copyWithin, nil, "copyWithin", nil, 2), true, false, true)
	b._putProp("entries", r.newNativeFunc(r.typedArrayProto_entries, nil, "entries", nil, 0), true, false, true)
	b._putProp("every", r.newNativeFunc(r.typedArrayProto_every, nil, "every", nil, 1), true, false, true)
	b._putProp("fill", r.newNativeFunc(r.typedArrayProto_fill, nil, "fill", nil, 1), true, false, true)
	b._putProp("filter", r.newNativeFunc(r.typedArrayProto_filter, nil, "filter", nil, 1), true, false, true)
	b._putProp("find", r.newNativeFunc(r.typedArrayProto_find, nil, "find", nil, 1), true, false, true)
	b._putProp("findIndex", r.newNativeFunc(r.typedArrayProto_findIndex, nil, "findIndex", nil, 1), true, false, true)
	b._putProp("forEach", r.newNativeFunc(r.typedArrayProto_forEach, nil, "forEach", nil, 1), true, false, true)
	b._putProp("includes", r.newNativeFunc(r.typedArrayProto_includes, nil, "includes", nil, 1), true, false, true)
	b._putProp("indexOf", r.newNativeFunc(r.typedArrayProto_indexOf, nil, "indexOf", nil, 1), true, false, true)
	b._putProp("join", r.newNativeFunc(r.typedArrayProto_join, nil, "join", nil, 1), true, false, true)
	b._putProp("keys", r.newNativeFunc(r.typedArrayProto_keys, nil, "keys", nil, 0), true, false, true)
	b._putProp("lastIndexOf", r.newNativeFunc(r.typedArrayProto_lastIndexOf, nil, "lastIndexOf", nil, 1), true, false, true)
	b._put("length", &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.typedArrayProto_getLength, nil, "get length", nil, 0),
	***REMOVED***)
	b._putProp("map", r.newNativeFunc(r.typedArrayProto_map, nil, "map", nil, 1), true, false, true)
	b._putProp("reduce", r.newNativeFunc(r.typedArrayProto_reduce, nil, "reduce", nil, 1), true, false, true)
	b._putProp("reduceRight", r.newNativeFunc(r.typedArrayProto_reduceRight, nil, "reduceRight", nil, 1), true, false, true)
	b._putProp("reverse", r.newNativeFunc(r.typedArrayProto_reverse, nil, "reverse", nil, 0), true, false, true)
	b._putProp("set", r.newNativeFunc(r.typedArrayProto_set, nil, "set", nil, 1), true, false, true)
	b._putProp("slice", r.newNativeFunc(r.typedArrayProto_slice, nil, "slice", nil, 2), true, false, true)
	b._putProp("some", r.newNativeFunc(r.typedArrayProto_some, nil, "some", nil, 1), true, false, true)
	b._putProp("sort", r.newNativeFunc(r.typedArrayProto_sort, nil, "sort", nil, 1), true, false, true)
	b._putProp("subarray", r.newNativeFunc(r.typedArrayProto_subarray, nil, "subarray", nil, 2), true, false, true)
	b._putProp("toLocaleString", r.newNativeFunc(r.typedArrayProto_toLocaleString, nil, "toLocaleString", nil, 0), true, false, true)
	b._putProp("toString", r.global.arrayToString, true, false, true)
	valuesFunc := r.newNativeFunc(r.typedArrayProto_values, nil, "values", nil, 0)
	b._putProp("values", valuesFunc, true, false, true)
	b._putSym(SymIterator, valueProp(valuesFunc, true, false, true))
	b._putSym(SymToStringTag, &valueProperty***REMOVED***
		getterFunc:   r.newNativeFunc(r.typedArrayProto_toStringTag, nil, "get [Symbol.toStringTag]", nil, 0),
		accessor:     true,
		configurable: true,
	***REMOVED***)

	return b
***REMOVED***

func (r *Runtime) createTypedArray(val *Object) objectImpl ***REMOVED***
	o := r.newNativeConstructOnly(val, r.newTypedArray, r.global.TypedArrayPrototype, "TypedArray", 0)
	o._putProp("from", r.newNativeFunc(r.typedArray_from, nil, "from", nil, 1), true, false, true)
	o._putProp("of", r.newNativeFunc(r.typedArray_of, nil, "of", nil, 0), true, false, true)
	r.putSpeciesReturnThis(o)

	return o
***REMOVED***

func (r *Runtime) typedArrayCreator(ctor func(args []Value, newTarget, proto *Object) *Object, name unistring.String, bytesPerElement int) func(val *Object) objectImpl ***REMOVED***
	return func(val *Object) objectImpl ***REMOVED***
		p := r.newBaseObject(r.global.TypedArrayPrototype, classObject)
		o := r.newNativeConstructOnly(val, func(args []Value, newTarget *Object) *Object ***REMOVED***
			return ctor(args, newTarget, p.val)
		***REMOVED***, p.val, name, 3)

		p._putProp("constructor", o.val, true, false, true)

		o.prototype = r.global.TypedArray
		bpe := intToValue(int64(bytesPerElement))
		o._putProp("BYTES_PER_ELEMENT", bpe, false, false, false)
		p._putProp("BYTES_PER_ELEMENT", bpe, false, false, false)
		return o
	***REMOVED***
***REMOVED***

func (r *Runtime) initTypedArrays() ***REMOVED***

	r.global.ArrayBufferPrototype = r.newLazyObject(r.createArrayBufferProto)
	r.global.ArrayBuffer = r.newLazyObject(r.createArrayBuffer)
	r.addToGlobal("ArrayBuffer", r.global.ArrayBuffer)

	r.global.DataViewPrototype = r.newLazyObject(r.createDataViewProto)
	r.global.DataView = r.newLazyObject(r.createDataView)
	r.addToGlobal("DataView", r.global.DataView)

	r.global.TypedArrayPrototype = r.newLazyObject(r.createTypedArrayProto)
	r.global.TypedArray = r.newLazyObject(r.createTypedArray)

	r.global.Uint8Array = r.newLazyObject(r.typedArrayCreator(r.newUint8Array, "Uint8Array", 1))
	r.addToGlobal("Uint8Array", r.global.Uint8Array)

	r.global.Uint8ClampedArray = r.newLazyObject(r.typedArrayCreator(r.newUint8ClampedArray, "Uint8ClampedArray", 1))
	r.addToGlobal("Uint8ClampedArray", r.global.Uint8ClampedArray)

	r.global.Int8Array = r.newLazyObject(r.typedArrayCreator(r.newInt8Array, "Int8Array", 1))
	r.addToGlobal("Int8Array", r.global.Int8Array)

	r.global.Uint16Array = r.newLazyObject(r.typedArrayCreator(r.newUint16Array, "Uint16Array", 2))
	r.addToGlobal("Uint16Array", r.global.Uint16Array)

	r.global.Int16Array = r.newLazyObject(r.typedArrayCreator(r.newInt16Array, "Int16Array", 2))
	r.addToGlobal("Int16Array", r.global.Int16Array)

	r.global.Uint32Array = r.newLazyObject(r.typedArrayCreator(r.newUint32Array, "Uint32Array", 4))
	r.addToGlobal("Uint32Array", r.global.Uint32Array)

	r.global.Int32Array = r.newLazyObject(r.typedArrayCreator(r.newInt32Array, "Int32Array", 4))
	r.addToGlobal("Int32Array", r.global.Int32Array)

	r.global.Float32Array = r.newLazyObject(r.typedArrayCreator(r.newFloat32Array, "Float32Array", 4))
	r.addToGlobal("Float32Array", r.global.Float32Array)

	r.global.Float64Array = r.newLazyObject(r.typedArrayCreator(r.newFloat64Array, "Float64Array", 8))
	r.addToGlobal("Float64Array", r.global.Float64Array)
***REMOVED***
