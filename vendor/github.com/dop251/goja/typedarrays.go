package goja

import (
	"math"
	"math/bits"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/dop251/goja/unistring"
)

type byteOrder bool

const (
	bigEndian    byteOrder = false
	littleEndian byteOrder = true
)

var (
	nativeEndian byteOrder

	arrayBufferType = reflect.TypeOf(ArrayBuffer***REMOVED******REMOVED***)
)

type typedArrayObjectCtor func(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject

type arrayBufferObject struct ***REMOVED***
	baseObject
	detached bool
	data     []byte
***REMOVED***

// ArrayBuffer is a Go wrapper around ECMAScript ArrayBuffer. Calling Runtime.ToValue() on it
// returns the underlying ArrayBuffer. Calling Export() on an ECMAScript ArrayBuffer returns a wrapper.
// Use Runtime.NewArrayBuffer([]byte) to create one.
type ArrayBuffer struct ***REMOVED***
	buf *arrayBufferObject
***REMOVED***

type dataViewObject struct ***REMOVED***
	baseObject
	viewedArrayBuf      *arrayBufferObject
	byteLen, byteOffset int
***REMOVED***

type typedArray interface ***REMOVED***
	toRaw(Value) uint64
	get(idx int) Value
	set(idx int, value Value)
	getRaw(idx int) uint64
	setRaw(idx int, raw uint64)
	less(i, j int) bool
	swap(i, j int)
	typeMatch(v Value) bool
***REMOVED***

type uint8Array []uint8
type uint8ClampedArray []uint8
type int8Array []int8
type uint16Array []uint16
type int16Array []int16
type uint32Array []uint32
type int32Array []int32
type float32Array []float32
type float64Array []float64

type typedArrayObject struct ***REMOVED***
	baseObject
	viewedArrayBuf *arrayBufferObject
	defaultCtor    *Object
	length, offset int
	elemSize       int
	typedArray     typedArray
***REMOVED***

func (a ArrayBuffer) toValue(r *Runtime) Value ***REMOVED***
	if a.buf == nil ***REMOVED***
		return _null
	***REMOVED***
	v := a.buf.val
	if v.runtime != r ***REMOVED***
		panic(r.NewTypeError("Illegal runtime transition of an ArrayBuffer"))
	***REMOVED***
	return v
***REMOVED***

// Bytes returns the underlying []byte for this ArrayBuffer.
// For detached ArrayBuffers returns nil.
func (a ArrayBuffer) Bytes() []byte ***REMOVED***
	return a.buf.data
***REMOVED***

// Detach the ArrayBuffer. After this, the underlying []byte becomes unreferenced and any attempt
// to use this ArrayBuffer results in a TypeError.
// Returns false if it was already detached, true otherwise.
// Note, this method may only be called from the goroutine that 'owns' the Runtime, it may not
// be called concurrently.
func (a ArrayBuffer) Detach() bool ***REMOVED***
	if a.buf.detached ***REMOVED***
		return false
	***REMOVED***
	a.buf.detach()
	return true
***REMOVED***

// Detached returns true if the ArrayBuffer is detached.
func (a ArrayBuffer) Detached() bool ***REMOVED***
	return a.buf.detached
***REMOVED***

func (r *Runtime) NewArrayBuffer(data []byte) ArrayBuffer ***REMOVED***
	buf := r._newArrayBuffer(r.global.ArrayBufferPrototype, nil)
	buf.data = data
	return ArrayBuffer***REMOVED***
		buf: buf,
	***REMOVED***
***REMOVED***

func (a *uint8Array) get(idx int) Value ***REMOVED***
	return intToValue(int64((*a)[idx]))
***REMOVED***

func (a *uint8Array) getRaw(idx int) uint64 ***REMOVED***
	return uint64((*a)[idx])
***REMOVED***

func (a *uint8Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toUint8(value)
***REMOVED***

func (a *uint8Array) toRaw(v Value) uint64 ***REMOVED***
	return uint64(toUint8(v))
***REMOVED***

func (a *uint8Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = uint8(v)
***REMOVED***

func (a *uint8Array) less(i, j int) bool ***REMOVED***
	return (*a)[i] < (*a)[j]
***REMOVED***

func (a *uint8Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *uint8Array) typeMatch(v Value) bool ***REMOVED***
	if i, ok := v.(valueInt); ok ***REMOVED***
		return i >= 0 && i <= 255
	***REMOVED***
	return false
***REMOVED***

func (a *uint8ClampedArray) get(idx int) Value ***REMOVED***
	return intToValue(int64((*a)[idx]))
***REMOVED***

func (a *uint8ClampedArray) getRaw(idx int) uint64 ***REMOVED***
	return uint64((*a)[idx])
***REMOVED***

func (a *uint8ClampedArray) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toUint8Clamp(value)
***REMOVED***

func (a *uint8ClampedArray) toRaw(v Value) uint64 ***REMOVED***
	return uint64(toUint8Clamp(v))
***REMOVED***

func (a *uint8ClampedArray) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = uint8(v)
***REMOVED***

func (a *uint8ClampedArray) less(i, j int) bool ***REMOVED***
	return (*a)[i] < (*a)[j]
***REMOVED***

func (a *uint8ClampedArray) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *uint8ClampedArray) typeMatch(v Value) bool ***REMOVED***
	if i, ok := v.(valueInt); ok ***REMOVED***
		return i >= 0 && i <= 255
	***REMOVED***
	return false
***REMOVED***

func (a *int8Array) get(idx int) Value ***REMOVED***
	return intToValue(int64((*a)[idx]))
***REMOVED***

func (a *int8Array) getRaw(idx int) uint64 ***REMOVED***
	return uint64((*a)[idx])
***REMOVED***

func (a *int8Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toInt8(value)
***REMOVED***

func (a *int8Array) toRaw(v Value) uint64 ***REMOVED***
	return uint64(toInt8(v))
***REMOVED***

func (a *int8Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = int8(v)
***REMOVED***

func (a *int8Array) less(i, j int) bool ***REMOVED***
	return (*a)[i] < (*a)[j]
***REMOVED***

func (a *int8Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *int8Array) typeMatch(v Value) bool ***REMOVED***
	if i, ok := v.(valueInt); ok ***REMOVED***
		return i >= math.MinInt8 && i <= math.MaxInt8
	***REMOVED***
	return false
***REMOVED***

func (a *uint16Array) get(idx int) Value ***REMOVED***
	return intToValue(int64((*a)[idx]))
***REMOVED***

func (a *uint16Array) getRaw(idx int) uint64 ***REMOVED***
	return uint64((*a)[idx])
***REMOVED***

func (a *uint16Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toUint16(value)
***REMOVED***

func (a *uint16Array) toRaw(v Value) uint64 ***REMOVED***
	return uint64(toUint16(v))
***REMOVED***

func (a *uint16Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = uint16(v)
***REMOVED***

func (a *uint16Array) less(i, j int) bool ***REMOVED***
	return (*a)[i] < (*a)[j]
***REMOVED***

func (a *uint16Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *uint16Array) typeMatch(v Value) bool ***REMOVED***
	if i, ok := v.(valueInt); ok ***REMOVED***
		return i >= 0 && i <= math.MaxUint16
	***REMOVED***
	return false
***REMOVED***

func (a *int16Array) get(idx int) Value ***REMOVED***
	return intToValue(int64((*a)[idx]))
***REMOVED***

func (a *int16Array) getRaw(idx int) uint64 ***REMOVED***
	return uint64((*a)[idx])
***REMOVED***

func (a *int16Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toInt16(value)
***REMOVED***

func (a *int16Array) toRaw(v Value) uint64 ***REMOVED***
	return uint64(toInt16(v))
***REMOVED***

func (a *int16Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = int16(v)
***REMOVED***

func (a *int16Array) less(i, j int) bool ***REMOVED***
	return (*a)[i] < (*a)[j]
***REMOVED***

func (a *int16Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *int16Array) typeMatch(v Value) bool ***REMOVED***
	if i, ok := v.(valueInt); ok ***REMOVED***
		return i >= math.MinInt16 && i <= math.MaxInt16
	***REMOVED***
	return false
***REMOVED***

func (a *uint32Array) get(idx int) Value ***REMOVED***
	return intToValue(int64((*a)[idx]))
***REMOVED***

func (a *uint32Array) getRaw(idx int) uint64 ***REMOVED***
	return uint64((*a)[idx])
***REMOVED***

func (a *uint32Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toUint32(value)
***REMOVED***

func (a *uint32Array) toRaw(v Value) uint64 ***REMOVED***
	return uint64(toUint32(v))
***REMOVED***

func (a *uint32Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = uint32(v)
***REMOVED***

func (a *uint32Array) less(i, j int) bool ***REMOVED***
	return (*a)[i] < (*a)[j]
***REMOVED***

func (a *uint32Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *uint32Array) typeMatch(v Value) bool ***REMOVED***
	if i, ok := v.(valueInt); ok ***REMOVED***
		return i >= 0 && i <= math.MaxUint32
	***REMOVED***
	return false
***REMOVED***

func (a *int32Array) get(idx int) Value ***REMOVED***
	return intToValue(int64((*a)[idx]))
***REMOVED***

func (a *int32Array) getRaw(idx int) uint64 ***REMOVED***
	return uint64((*a)[idx])
***REMOVED***

func (a *int32Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toInt32(value)
***REMOVED***

func (a *int32Array) toRaw(v Value) uint64 ***REMOVED***
	return uint64(toInt32(v))
***REMOVED***

func (a *int32Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = int32(v)
***REMOVED***

func (a *int32Array) less(i, j int) bool ***REMOVED***
	return (*a)[i] < (*a)[j]
***REMOVED***

func (a *int32Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *int32Array) typeMatch(v Value) bool ***REMOVED***
	if i, ok := v.(valueInt); ok ***REMOVED***
		return i >= math.MinInt32 && i <= math.MaxInt32
	***REMOVED***
	return false
***REMOVED***

func (a *float32Array) get(idx int) Value ***REMOVED***
	return floatToValue(float64((*a)[idx]))
***REMOVED***

func (a *float32Array) getRaw(idx int) uint64 ***REMOVED***
	return uint64(math.Float32bits((*a)[idx]))
***REMOVED***

func (a *float32Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = toFloat32(value)
***REMOVED***

func (a *float32Array) toRaw(v Value) uint64 ***REMOVED***
	return uint64(math.Float32bits(toFloat32(v)))
***REMOVED***

func (a *float32Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = math.Float32frombits(uint32(v))
***REMOVED***

func typedFloatLess(x, y float64) bool ***REMOVED***
	xNan := math.IsNaN(x)
	yNan := math.IsNaN(y)
	if yNan ***REMOVED***
		return !xNan
	***REMOVED*** else if xNan ***REMOVED***
		return false
	***REMOVED***
	if x == 0 && y == 0 ***REMOVED*** // handle neg zero
		return math.Signbit(x)
	***REMOVED***
	return x < y
***REMOVED***

func (a *float32Array) less(i, j int) bool ***REMOVED***
	return typedFloatLess(float64((*a)[i]), float64((*a)[j]))
***REMOVED***

func (a *float32Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *float32Array) typeMatch(v Value) bool ***REMOVED***
	switch v.(type) ***REMOVED***
	case valueInt, valueFloat:
		return true
	***REMOVED***
	return false
***REMOVED***

func (a *float64Array) get(idx int) Value ***REMOVED***
	return floatToValue((*a)[idx])
***REMOVED***

func (a *float64Array) getRaw(idx int) uint64 ***REMOVED***
	return math.Float64bits((*a)[idx])
***REMOVED***

func (a *float64Array) set(idx int, value Value) ***REMOVED***
	(*a)[idx] = value.ToFloat()
***REMOVED***

func (a *float64Array) toRaw(v Value) uint64 ***REMOVED***
	return math.Float64bits(v.ToFloat())
***REMOVED***

func (a *float64Array) setRaw(idx int, v uint64) ***REMOVED***
	(*a)[idx] = math.Float64frombits(v)
***REMOVED***

func (a *float64Array) less(i, j int) bool ***REMOVED***
	return typedFloatLess((*a)[i], (*a)[j])
***REMOVED***

func (a *float64Array) swap(i, j int) ***REMOVED***
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
***REMOVED***

func (a *float64Array) typeMatch(v Value) bool ***REMOVED***
	switch v.(type) ***REMOVED***
	case valueInt, valueFloat:
		return true
	***REMOVED***
	return false
***REMOVED***

func (a *typedArrayObject) _getIdx(idx int) Value ***REMOVED***
	a.viewedArrayBuf.ensureNotDetached()
	if 0 <= idx && idx < a.length ***REMOVED***
		return a.typedArray.get(idx + a.offset)
	***REMOVED***
	return nil
***REMOVED***

func strToTAIdx(s unistring.String) (int, bool) ***REMOVED***
	i, err := strconv.ParseInt(string(s), 10, bits.UintSize)
	if err != nil ***REMOVED***
		return 0, false
	***REMOVED***
	return int(i), true
***REMOVED***

func (a *typedArrayObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if idx, ok := strToTAIdx(name); ok ***REMOVED***
		v := a._getIdx(idx)
		if v != nil ***REMOVED***
			return &valueProperty***REMOVED***
				value:      v,
				writable:   true,
				enumerable: true,
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	return a.baseObject.getOwnPropStr(name)
***REMOVED***

func (a *typedArrayObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	v := a._getIdx(toIntStrict(int64(idx)))
	if v != nil ***REMOVED***
		return &valueProperty***REMOVED***
			value:      v,
			writable:   true,
			enumerable: true,
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (a *typedArrayObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	if idx, ok := strToTAIdx(name); ok ***REMOVED***
		prop := a._getIdx(idx)
		if prop == nil ***REMOVED***
			if a.prototype != nil ***REMOVED***
				if receiver == nil ***REMOVED***
					return a.prototype.self.getStr(name, a.val)
				***REMOVED***
				return a.prototype.self.getStr(name, receiver)
			***REMOVED***
		***REMOVED***
		return prop
	***REMOVED***
	return a.baseObject.getStr(name, receiver)
***REMOVED***

func (a *typedArrayObject) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	prop := a._getIdx(toIntStrict(int64(idx)))
	if prop == nil ***REMOVED***
		if a.prototype != nil ***REMOVED***
			if receiver == nil ***REMOVED***
				return a.prototype.self.getIdx(idx, a.val)
			***REMOVED***
			return a.prototype.self.getIdx(idx, receiver)
		***REMOVED***
	***REMOVED***
	return prop
***REMOVED***

func (a *typedArrayObject) _putIdx(idx int, v Value, throw bool) bool ***REMOVED***
	v = v.ToNumber()
	a.viewedArrayBuf.ensureNotDetached()
	if idx >= 0 && idx < a.length ***REMOVED***
		a.typedArray.set(idx+a.offset, v)
		return true
	***REMOVED***
	// As far as I understand the specification this should throw, but neither V8 nor SpiderMonkey does
	return false
***REMOVED***

func (a *typedArrayObject) _hasIdx(idx int) bool ***REMOVED***
	a.viewedArrayBuf.ensureNotDetached()
	return idx >= 0 && idx < a.length
***REMOVED***

func (a *typedArrayObject) setOwnStr(p unistring.String, v Value, throw bool) bool ***REMOVED***
	if idx, ok := strToTAIdx(p); ok ***REMOVED***
		return a._putIdx(idx, v, throw)
	***REMOVED***
	return a.baseObject.setOwnStr(p, v, throw)
***REMOVED***

func (a *typedArrayObject) setOwnIdx(p valueInt, v Value, throw bool) bool ***REMOVED***
	return a._putIdx(toIntStrict(int64(p)), v, throw)
***REMOVED***

func (a *typedArrayObject) setForeignStr(p unistring.String, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	return a._setForeignStr(p, a.getOwnPropStr(p), v, receiver, throw)
***REMOVED***

func (a *typedArrayObject) setForeignIdx(p valueInt, v, receiver Value, throw bool) (res bool, handled bool) ***REMOVED***
	return a._setForeignIdx(p, trueValIfPresent(a.hasOwnPropertyIdx(p)), v, receiver, throw)
***REMOVED***

func (a *typedArrayObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if idx, ok := strToTAIdx(name); ok ***REMOVED***
		a.viewedArrayBuf.ensureNotDetached()
		return idx < a.length
	***REMOVED***

	return a.baseObject.hasOwnPropertyStr(name)
***REMOVED***

func (a *typedArrayObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	return a._hasIdx(toIntStrict(int64(idx)))
***REMOVED***

func (a *typedArrayObject) _defineIdxProperty(idx int, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	prop, ok := a._defineOwnProperty(unistring.String(strconv.Itoa(idx)), a.getOwnPropIdx(valueInt(idx)), desc, throw)
	if ok ***REMOVED***
		return a._putIdx(idx, prop, throw)
	***REMOVED***
	return ok
***REMOVED***

func (a *typedArrayObject) defineOwnPropertyStr(name unistring.String, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx, ok := strToTAIdx(name); ok ***REMOVED***
		return a._defineIdxProperty(idx, desc, throw)
	***REMOVED***
	return a.baseObject.defineOwnPropertyStr(name, desc, throw)
***REMOVED***

func (a *typedArrayObject) defineOwnPropertyIdx(name valueInt, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	return a._defineIdxProperty(toIntStrict(int64(name)), desc, throw)
***REMOVED***

func (a *typedArrayObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if idx, ok := strToTAIdx(name); ok ***REMOVED***
		if idx < a.length ***REMOVED***
			a.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of %s", idx, a.val.String())
		***REMOVED***
	***REMOVED***

	return a.baseObject.deleteStr(name, throw)
***REMOVED***

func (a *typedArrayObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	if idx >= 0 && int64(idx) < int64(a.length) ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of %s", idx, a.val.String())
	***REMOVED***

	return true
***REMOVED***

func (a *typedArrayObject) ownKeys(all bool, accum []Value) []Value ***REMOVED***
	if accum == nil ***REMOVED***
		accum = make([]Value, 0, a.length)
	***REMOVED***
	for i := 0; i < a.length; i++ ***REMOVED***
		accum = append(accum, asciiString(strconv.Itoa(i)))
	***REMOVED***
	return a.baseObject.ownKeys(all, accum)
***REMOVED***

type typedArrayPropIter struct ***REMOVED***
	a   *typedArrayObject
	idx int
***REMOVED***

func (i *typedArrayPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	if i.idx < i.a.length ***REMOVED***
		name := strconv.Itoa(i.idx)
		prop := i.a._getIdx(i.idx)
		i.idx++
		return propIterItem***REMOVED***name: unistring.String(name), value: prop***REMOVED***, i.next
	***REMOVED***

	return i.a.baseObject.enumerateUnfiltered()()
***REMOVED***

func (a *typedArrayObject) enumerateUnfiltered() iterNextFunc ***REMOVED***
	return (&typedArrayPropIter***REMOVED***
		a: a,
	***REMOVED***).next
***REMOVED***

func (r *Runtime) _newTypedArrayObject(buf *arrayBufferObject, offset, length, elemSize int, defCtor *Object, arr typedArray, proto *Object) *typedArrayObject ***REMOVED***
	o := &Object***REMOVED***runtime: r***REMOVED***
	a := &typedArrayObject***REMOVED***
		baseObject: baseObject***REMOVED***
			val:        o,
			class:      classObject,
			prototype:  proto,
			extensible: true,
		***REMOVED***,
		viewedArrayBuf: buf,
		offset:         offset,
		length:         length,
		elemSize:       elemSize,
		defaultCtor:    defCtor,
		typedArray:     arr,
	***REMOVED***
	o.self = a
	a.init()
	return a

***REMOVED***

func (r *Runtime) newUint8ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 1, r.global.Uint8Array, (*uint8Array)(&buf.data), proto)
***REMOVED***

func (r *Runtime) newUint8ClampedArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 1, r.global.Uint8ClampedArray, (*uint8ClampedArray)(&buf.data), proto)
***REMOVED***

func (r *Runtime) newInt8ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 1, r.global.Int8Array, (*int8Array)(unsafe.Pointer(&buf.data)), proto)
***REMOVED***

func (r *Runtime) newUint16ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 2, r.global.Uint16Array, (*uint16Array)(unsafe.Pointer(&buf.data)), proto)
***REMOVED***

func (r *Runtime) newInt16ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 2, r.global.Int16Array, (*int16Array)(unsafe.Pointer(&buf.data)), proto)
***REMOVED***

func (r *Runtime) newUint32ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 4, r.global.Uint32Array, (*uint32Array)(unsafe.Pointer(&buf.data)), proto)
***REMOVED***

func (r *Runtime) newInt32ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 4, r.global.Int32Array, (*int32Array)(unsafe.Pointer(&buf.data)), proto)
***REMOVED***

func (r *Runtime) newFloat32ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 4, r.global.Float32Array, (*float32Array)(unsafe.Pointer(&buf.data)), proto)
***REMOVED***

func (r *Runtime) newFloat64ArrayObject(buf *arrayBufferObject, offset, length int, proto *Object) *typedArrayObject ***REMOVED***
	return r._newTypedArrayObject(buf, offset, length, 8, r.global.Float64Array, (*float64Array)(unsafe.Pointer(&buf.data)), proto)
***REMOVED***

func (o *dataViewObject) getIdxAndByteOrder(idxVal, littleEndianVal Value, size int) (int, byteOrder) ***REMOVED***
	getIdx := o.val.runtime.toIndex(idxVal)
	o.viewedArrayBuf.ensureNotDetached()
	if getIdx+size > o.byteLen ***REMOVED***
		panic(o.val.runtime.newError(o.val.runtime.global.RangeError, "Index %d is out of bounds", getIdx))
	***REMOVED***
	getIdx += o.byteOffset
	var bo byteOrder
	if littleEndianVal != nil ***REMOVED***
		if littleEndianVal.ToBoolean() ***REMOVED***
			bo = littleEndian
		***REMOVED*** else ***REMOVED***
			bo = bigEndian
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		bo = nativeEndian
	***REMOVED***
	return getIdx, bo
***REMOVED***

func (o *arrayBufferObject) ensureNotDetached() ***REMOVED***
	if o.detached ***REMOVED***
		panic(o.val.runtime.NewTypeError("ArrayBuffer is detached"))
	***REMOVED***
***REMOVED***

func (o *arrayBufferObject) getFloat32(idx int, byteOrder byteOrder) float32 ***REMOVED***
	return math.Float32frombits(o.getUint32(idx, byteOrder))
***REMOVED***

func (o *arrayBufferObject) setFloat32(idx int, val float32, byteOrder byteOrder) ***REMOVED***
	o.setUint32(idx, math.Float32bits(val), byteOrder)
***REMOVED***

func (o *arrayBufferObject) getFloat64(idx int, byteOrder byteOrder) float64 ***REMOVED***
	return math.Float64frombits(o.getUint64(idx, byteOrder))
***REMOVED***

func (o *arrayBufferObject) setFloat64(idx int, val float64, byteOrder byteOrder) ***REMOVED***
	o.setUint64(idx, math.Float64bits(val), byteOrder)
***REMOVED***

func (o *arrayBufferObject) getUint64(idx int, byteOrder byteOrder) uint64 ***REMOVED***
	var b []byte
	if byteOrder == nativeEndian ***REMOVED***
		b = o.data[idx : idx+8]
	***REMOVED*** else ***REMOVED***
		b = make([]byte, 8)
		d := o.data[idx : idx+8]
		b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7] = d[7], d[6], d[5], d[4], d[3], d[2], d[1], d[0]
	***REMOVED***
	return *((*uint64)(unsafe.Pointer(&b[0])))
***REMOVED***

func (o *arrayBufferObject) setUint64(idx int, val uint64, byteOrder byteOrder) ***REMOVED***
	if byteOrder == nativeEndian ***REMOVED***
		*(*uint64)(unsafe.Pointer(&o.data[idx])) = val
	***REMOVED*** else ***REMOVED***
		b := (*[8]byte)(unsafe.Pointer(&val))
		d := o.data[idx : idx+8]
		d[0], d[1], d[2], d[3], d[4], d[5], d[6], d[7] = b[7], b[6], b[5], b[4], b[3], b[2], b[1], b[0]
	***REMOVED***
***REMOVED***

func (o *arrayBufferObject) getUint32(idx int, byteOrder byteOrder) uint32 ***REMOVED***
	var b []byte
	if byteOrder == nativeEndian ***REMOVED***
		b = o.data[idx : idx+4]
	***REMOVED*** else ***REMOVED***
		b = make([]byte, 4)
		d := o.data[idx : idx+4]
		b[0], b[1], b[2], b[3] = d[3], d[2], d[1], d[0]
	***REMOVED***
	return *((*uint32)(unsafe.Pointer(&b[0])))
***REMOVED***

func (o *arrayBufferObject) setUint32(idx int, val uint32, byteOrder byteOrder) ***REMOVED***
	if byteOrder == nativeEndian ***REMOVED***
		*(*uint32)(unsafe.Pointer(&o.data[idx])) = val
	***REMOVED*** else ***REMOVED***
		b := (*[4]byte)(unsafe.Pointer(&val))
		d := o.data[idx : idx+4]
		d[0], d[1], d[2], d[3] = b[3], b[2], b[1], b[0]
	***REMOVED***
***REMOVED***

func (o *arrayBufferObject) getUint16(idx int, byteOrder byteOrder) uint16 ***REMOVED***
	var b []byte
	if byteOrder == nativeEndian ***REMOVED***
		b = o.data[idx : idx+2]
	***REMOVED*** else ***REMOVED***
		b = make([]byte, 2)
		d := o.data[idx : idx+2]
		b[0], b[1] = d[1], d[0]
	***REMOVED***
	return *((*uint16)(unsafe.Pointer(&b[0])))
***REMOVED***

func (o *arrayBufferObject) setUint16(idx int, val uint16, byteOrder byteOrder) ***REMOVED***
	if byteOrder == nativeEndian ***REMOVED***
		*(*uint16)(unsafe.Pointer(&o.data[idx])) = val
	***REMOVED*** else ***REMOVED***
		b := (*[2]byte)(unsafe.Pointer(&val))
		d := o.data[idx : idx+2]
		d[0], d[1] = b[1], b[0]
	***REMOVED***
***REMOVED***

func (o *arrayBufferObject) getUint8(idx int) uint8 ***REMOVED***
	return o.data[idx]
***REMOVED***

func (o *arrayBufferObject) setUint8(idx int, val uint8) ***REMOVED***
	o.data[idx] = val
***REMOVED***

func (o *arrayBufferObject) getInt32(idx int, byteOrder byteOrder) int32 ***REMOVED***
	return int32(o.getUint32(idx, byteOrder))
***REMOVED***

func (o *arrayBufferObject) setInt32(idx int, val int32, byteOrder byteOrder) ***REMOVED***
	o.setUint32(idx, uint32(val), byteOrder)
***REMOVED***

func (o *arrayBufferObject) getInt16(idx int, byteOrder byteOrder) int16 ***REMOVED***
	return int16(o.getUint16(idx, byteOrder))
***REMOVED***

func (o *arrayBufferObject) setInt16(idx int, val int16, byteOrder byteOrder) ***REMOVED***
	o.setUint16(idx, uint16(val), byteOrder)
***REMOVED***

func (o *arrayBufferObject) getInt8(idx int) int8 ***REMOVED***
	return int8(o.data[idx])
***REMOVED***

func (o *arrayBufferObject) setInt8(idx int, val int8) ***REMOVED***
	o.setUint8(idx, uint8(val))
***REMOVED***

func (o *arrayBufferObject) detach() ***REMOVED***
	o.data = nil
	o.detached = true
***REMOVED***

func (o *arrayBufferObject) exportType() reflect.Type ***REMOVED***
	return arrayBufferType
***REMOVED***

func (o *arrayBufferObject) export(*objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	return ArrayBuffer***REMOVED***
		buf: o,
	***REMOVED***
***REMOVED***

func (r *Runtime) _newArrayBuffer(proto *Object, o *Object) *arrayBufferObject ***REMOVED***
	if o == nil ***REMOVED***
		o = &Object***REMOVED***runtime: r***REMOVED***
	***REMOVED***
	b := &arrayBufferObject***REMOVED***
		baseObject: baseObject***REMOVED***
			class:      classObject,
			val:        o,
			prototype:  proto,
			extensible: true,
		***REMOVED***,
	***REMOVED***
	o.self = b
	b.init()
	return b
***REMOVED***

func init() ***REMOVED***
	buf := [2]byte***REMOVED******REMOVED***
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xCAFE)

	switch buf ***REMOVED***
	case [2]byte***REMOVED***0xFE, 0xCA***REMOVED***:
		nativeEndian = littleEndian
	case [2]byte***REMOVED***0xCA, 0xFE***REMOVED***:
		nativeEndian = bigEndian
	default:
		panic("Could not determine native endianness.")
	***REMOVED***
***REMOVED***
