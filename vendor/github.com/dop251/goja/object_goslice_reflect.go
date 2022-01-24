package goja

import (
	"math"
	"math/bits"
	"reflect"

	"github.com/dop251/goja/unistring"
)

type objectGoSliceReflect struct ***REMOVED***
	objectGoArrayReflect
***REMOVED***

func (o *objectGoSliceReflect) init() ***REMOVED***
	o.objectGoArrayReflect._init()
	o.lengthProp.writable = true
	o.putIdx = o._putIdx
***REMOVED***

func (o *objectGoSliceReflect) _putIdx(idx int, v Value, throw bool) bool ***REMOVED***
	if idx >= o.value.Len() ***REMOVED***
		o.grow(idx + 1)
	***REMOVED***
	return o.objectGoArrayReflect._putIdx(idx, v, throw)
***REMOVED***

func (o *objectGoSliceReflect) grow(size int) ***REMOVED***
	oldcap := o.value.Cap()
	if oldcap < size ***REMOVED***
		n := reflect.MakeSlice(o.value.Type(), size, growCap(size, o.value.Len(), oldcap))
		reflect.Copy(n, o.value)
		o.value.Set(n)
	***REMOVED*** else ***REMOVED***
		tail := o.value.Slice(o.value.Len(), size)
		zero := reflect.Zero(o.value.Type().Elem())
		for i := 0; i < tail.Len(); i++ ***REMOVED***
			tail.Index(i).Set(zero)
		***REMOVED***
		o.value.SetLen(size)
	***REMOVED***
	o.updateLen()
***REMOVED***

func (o *objectGoSliceReflect) shrink(size int) ***REMOVED***
	tail := o.value.Slice(size, o.value.Len())
	zero := reflect.Zero(o.value.Type().Elem())
	for i := 0; i < tail.Len(); i++ ***REMOVED***
		tail.Index(i).Set(zero)
	***REMOVED***
	o.value.SetLen(size)
	o.updateLen()
***REMOVED***

func (o *objectGoSliceReflect) putLength(v uint32, throw bool) bool ***REMOVED***
	if bits.UintSize == 32 && v > math.MaxInt32 ***REMOVED***
		panic(rangeError("Integer value overflows 32-bit int"))
	***REMOVED***
	newLen := int(v)
	curLen := o.value.Len()
	if newLen > curLen ***REMOVED***
		o.grow(newLen)
	***REMOVED*** else if newLen < curLen ***REMOVED***
		o.shrink(newLen)
	***REMOVED***
	return true
***REMOVED***

func (o *objectGoSliceReflect) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if name == "length" ***REMOVED***
		return o.putLength(o.val.runtime.toLengthUint32(val), throw)
	***REMOVED***
	return o.objectGoArrayReflect.setOwnStr(name, val, throw)
***REMOVED***

func (o *objectGoSliceReflect) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if name == "length" ***REMOVED***
		return o.val.runtime.defineArrayLength(&o.lengthProp, descr, o.putLength, throw)
	***REMOVED***
	return o.objectGoArrayReflect.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (o *objectGoSliceReflect) equal(other objectImpl) bool ***REMOVED***
	if other, ok := other.(*objectGoSliceReflect); ok ***REMOVED***
		return o.value.Interface() == other.value.Interface()
	***REMOVED***
	return false
***REMOVED***
