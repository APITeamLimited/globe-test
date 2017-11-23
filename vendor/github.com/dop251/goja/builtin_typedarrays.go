package goja

type objectArrayBuffer struct ***REMOVED***
	baseObject
	data []byte
***REMOVED***

func (o *objectArrayBuffer) export() interface***REMOVED******REMOVED*** ***REMOVED***
	return o.data
***REMOVED***

func (r *Runtime) _newArrayBuffer(proto *Object, o *Object) *objectArrayBuffer ***REMOVED***
	if o == nil ***REMOVED***
		o = &Object***REMOVED***runtime: r***REMOVED***
	***REMOVED***
	b := &objectArrayBuffer***REMOVED***
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

func (r *Runtime) builtin_ArrayBuffer(args []Value, proto *Object) *Object ***REMOVED***
	b := r._newArrayBuffer(proto, nil)
	if len(args) > 0 ***REMOVED***
		b.data = make([]byte, toLength(args[0]))
	***REMOVED***
	return b.val
***REMOVED***

func (r *Runtime) arrayBufferProto_getByteLength(call FunctionCall) Value ***REMOVED***
	o := r.toObject(call.This)
	if b, ok := o.self.(*objectArrayBuffer); ok ***REMOVED***
		return intToValue(int64(len(b.data)))
	***REMOVED***
	r.typeErrorResult(true, "Object is not ArrayBuffer: %s", o)
	panic("unreachable")
***REMOVED***

func (r *Runtime) arrayBufferProto_slice(call FunctionCall) Value ***REMOVED***
	o := r.toObject(call.This)
	if b, ok := o.self.(*objectArrayBuffer); ok ***REMOVED***
		l := int64(len(b.data))
		start := toLength(call.Argument(0))
		if start < 0 ***REMOVED***
			start = l + start
		***REMOVED***
		if start < 0 ***REMOVED***
			start = 0
		***REMOVED*** else if start > l ***REMOVED***
			start = l
		***REMOVED***
		var stop int64
		if arg := call.Argument(1); arg != _undefined ***REMOVED***
			stop = toLength(arg)
			if stop < 0 ***REMOVED***
				stop = int64(len(b.data)) + stop
			***REMOVED***
			if stop < 0 ***REMOVED***
				stop = 0
			***REMOVED*** else if stop > l ***REMOVED***
				stop = l
			***REMOVED***

		***REMOVED*** else ***REMOVED***
			stop = l
		***REMOVED***

		ret := r._newArrayBuffer(r.global.ArrayBufferPrototype, nil)

		if stop > start ***REMOVED***
			ret.data = b.data[start:stop]
		***REMOVED***

		return ret.val
	***REMOVED***
	r.typeErrorResult(true, "Object is not ArrayBuffer: %s", o)
	panic("unreachable")
***REMOVED***

func (r *Runtime) createArrayBufferProto(val *Object) objectImpl ***REMOVED***
	b := r._newArrayBuffer(r.global.Object, val)
	byteLengthProp := &valueProperty***REMOVED***
		accessor:     true,
		configurable: true,
		getterFunc:   r.newNativeFunc(r.arrayBufferProto_getByteLength, nil, "get byteLength", nil, 0),
	***REMOVED***
	b._put("byteLength", byteLengthProp)
	b._putProp("slice", r.newNativeFunc(r.arrayBufferProto_slice, nil, "slice", nil, 2), true, false, true)
	return b
***REMOVED***

func (r *Runtime) initTypedArrays() ***REMOVED***

	r.global.ArrayBufferPrototype = r.newLazyObject(r.createArrayBufferProto)

	r.global.ArrayBuffer = r.newNativeFuncConstruct(r.builtin_ArrayBuffer, "ArrayBuffer", r.global.ArrayBufferPrototype, 1)
	r.addToGlobal("ArrayBuffer", r.global.ArrayBuffer)
***REMOVED***
