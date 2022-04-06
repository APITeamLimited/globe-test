package goja

import (
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"strconv"

	"github.com/dop251/goja/unistring"
)

type arrayIterObject struct ***REMOVED***
	baseObject
	obj     *Object
	nextIdx int64
	kind    iterationKind
***REMOVED***

func (ai *arrayIterObject) next() Value ***REMOVED***
	if ai.obj == nil ***REMOVED***
		return ai.val.runtime.createIterResultObject(_undefined, true)
	***REMOVED***
	if ta, ok := ai.obj.self.(*typedArrayObject); ok ***REMOVED***
		ta.viewedArrayBuf.ensureNotDetached(true)
	***REMOVED***
	l := toLength(ai.obj.self.getStr("length", nil))
	index := ai.nextIdx
	if index >= l ***REMOVED***
		ai.obj = nil
		return ai.val.runtime.createIterResultObject(_undefined, true)
	***REMOVED***
	ai.nextIdx++
	idxVal := valueInt(index)
	if ai.kind == iterationKindKey ***REMOVED***
		return ai.val.runtime.createIterResultObject(idxVal, false)
	***REMOVED***
	elementValue := nilSafe(ai.obj.self.getIdx(idxVal, nil))
	var result Value
	if ai.kind == iterationKindValue ***REMOVED***
		result = elementValue
	***REMOVED*** else ***REMOVED***
		result = ai.val.runtime.newArrayValues([]Value***REMOVED***idxVal, elementValue***REMOVED***)
	***REMOVED***
	return ai.val.runtime.createIterResultObject(result, false)
***REMOVED***

func (r *Runtime) createArrayIterator(iterObj *Object, kind iterationKind) Value ***REMOVED***
	o := &Object***REMOVED***runtime: r***REMOVED***

	ai := &arrayIterObject***REMOVED***
		obj:  iterObj,
		kind: kind,
	***REMOVED***
	ai.class = classArrayIterator
	ai.val = o
	ai.extensible = true
	o.self = ai
	ai.prototype = r.global.ArrayIteratorPrototype
	ai.init()

	return o
***REMOVED***

type arrayObject struct ***REMOVED***
	baseObject
	values         []Value
	length         uint32
	objCount       int
	propValueCount int
	lengthProp     valueProperty
***REMOVED***

func (a *arrayObject) init() ***REMOVED***
	a.baseObject.init()
	a.lengthProp.writable = true

	a._put("length", &a.lengthProp)
***REMOVED***

func (a *arrayObject) _setLengthInt(l uint32, throw bool) bool ***REMOVED***
	ret := true
	if l <= a.length ***REMOVED***
		if a.propValueCount > 0 ***REMOVED***
			// Slow path
			for i := len(a.values) - 1; i >= int(l); i-- ***REMOVED***
				if prop, ok := a.values[i].(*valueProperty); ok ***REMOVED***
					if !prop.configurable ***REMOVED***
						l = uint32(i) + 1
						ret = false
						break
					***REMOVED***
					a.propValueCount--
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if l <= uint32(len(a.values)) ***REMOVED***
		if l >= 16 && l < uint32(cap(a.values))>>2 ***REMOVED***
			ar := make([]Value, l)
			copy(ar, a.values)
			a.values = ar
		***REMOVED*** else ***REMOVED***
			ar := a.values[l:len(a.values)]
			for i := range ar ***REMOVED***
				ar[i] = nil
			***REMOVED***
			a.values = a.values[:l]
		***REMOVED***
	***REMOVED***
	a.length = l
	if !ret ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "Cannot redefine property: length")
	***REMOVED***
	return ret
***REMOVED***

func (a *arrayObject) setLengthInt(l uint32, throw bool) bool ***REMOVED***
	if l == a.length ***REMOVED***
		return true
	***REMOVED***
	if !a.lengthProp.writable ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	***REMOVED***
	return a._setLengthInt(l, throw)
***REMOVED***

func (a *arrayObject) setLength(v uint32, throw bool) bool ***REMOVED***
	if !a.lengthProp.writable ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	***REMOVED***
	return a._setLengthInt(v, throw)
***REMOVED***

func (a *arrayObject) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
	prop := a.getOwnPropIdx(idx)
	if prop == nil ***REMOVED***
		if a.prototype != nil ***REMOVED***
			if receiver == nil ***REMOVED***
				return a.prototype.self.getIdx(idx, a.val)
			***REMOVED***
			return a.prototype.self.getIdx(idx, receiver)
		***REMOVED***
	***REMOVED***
	if prop, ok := prop.(*valueProperty); ok ***REMOVED***
		if receiver == nil ***REMOVED***
			return prop.get(a.val)
		***REMOVED***
		return prop.get(receiver)
	***REMOVED***
	return prop
***REMOVED***

func (a *arrayObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if len(a.values) > 0 ***REMOVED***
		if i := strToArrayIdx(name); i != math.MaxUint32 ***REMOVED***
			if i < uint32(len(a.values)) ***REMOVED***
				return a.values[i]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getOwnPropStr(name)
***REMOVED***

func (a *arrayObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	if i := toIdx(idx); i != math.MaxUint32 ***REMOVED***
		if i < uint32(len(a.values)) ***REMOVED***
			return a.values[i]
		***REMOVED***
		return nil
	***REMOVED***

	return a.baseObject.getOwnPropStr(idx.string())
***REMOVED***

func (a *arrayObject) sortLen() int64 ***REMOVED***
	return int64(len(a.values))
***REMOVED***

func (a *arrayObject) sortGet(i int64) Value ***REMOVED***
	v := a.values[i]
	if p, ok := v.(*valueProperty); ok ***REMOVED***
		v = p.get(a.val)
	***REMOVED***
	return v
***REMOVED***

func (a *arrayObject) swap(i, j int64) ***REMOVED***
	a.values[i], a.values[j] = a.values[j], a.values[i]
***REMOVED***

func (a *arrayObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	return a.getStrWithOwnProp(a.getOwnPropStr(name), name, receiver)
***REMOVED***

func (a *arrayObject) getLengthProp() *valueProperty ***REMOVED***
	a.lengthProp.value = intToValue(int64(a.length))
	return &a.lengthProp
***REMOVED***

func (a *arrayObject) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	if i := toIdx(idx); i != math.MaxUint32 ***REMOVED***
		return a._setOwnIdx(i, val, throw)
	***REMOVED*** else ***REMOVED***
		return a.baseObject.setOwnStr(idx.string(), val, throw)
	***REMOVED***
***REMOVED***

func (a *arrayObject) _setOwnIdx(idx uint32, val Value, throw bool) bool ***REMOVED***
	var prop Value
	if idx < uint32(len(a.values)) ***REMOVED***
		prop = a.values[idx]
	***REMOVED***

	if prop == nil ***REMOVED***
		if proto := a.prototype; proto != nil ***REMOVED***
			// we know it's foreign because prototype loops are not allowed
			if res, ok := proto.self.setForeignIdx(valueInt(idx), val, a.val, throw); ok ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
		// new property
		if !a.extensible ***REMOVED***
			a.val.runtime.typeErrorResult(throw, "Cannot add property %d, object is not extensible", idx)
			return false
		***REMOVED*** else ***REMOVED***
			if idx >= a.length ***REMOVED***
				if !a.setLengthInt(idx+1, throw) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if idx >= uint32(len(a.values)) ***REMOVED***
				if !a.expand(idx) ***REMOVED***
					a.val.self.(*sparseArrayObject).add(idx, val)
					return true
				***REMOVED***
			***REMOVED***
			a.objCount++
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if prop, ok := prop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				a.val.runtime.typeErrorResult(throw)
				return false
			***REMOVED***
			prop.set(a.val, val)
			return true
		***REMOVED***
	***REMOVED***
	a.values[idx] = val
	return true
***REMOVED***

func (a *arrayObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._setOwnIdx(idx, val, throw)
	***REMOVED*** else ***REMOVED***
		if name == "length" ***REMOVED***
			return a.setLength(a.val.runtime.toLengthUint32(val), throw)
		***REMOVED*** else ***REMOVED***
			return a.baseObject.setOwnStr(name, val, throw)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *arrayObject) setForeignIdx(idx valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return a._setForeignIdx(idx, a.getOwnPropIdx(idx), val, receiver, throw)
***REMOVED***

func (a *arrayObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return a._setForeignStr(name, a.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

type arrayPropIter struct ***REMOVED***
	a     *arrayObject
	limit int
	idx   int
***REMOVED***

func (i *arrayPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.a.values) && i.idx < i.limit ***REMOVED***
		name := asciiString(strconv.Itoa(i.idx))
		prop := i.a.values[i.idx]
		i.idx++
		if prop != nil ***REMOVED***
			return propIterItem***REMOVED***name: name, value: prop***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	return i.a.baseObject.iterateStringKeys()()
***REMOVED***

func (a *arrayObject) iterateStringKeys() iterNextFunc ***REMOVED***
	return (&arrayPropIter***REMOVED***
		a:     a,
		limit: len(a.values),
	***REMOVED***).next
***REMOVED***

func (a *arrayObject) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	for i, prop := range a.values ***REMOVED***
		name := strconv.Itoa(i)
		if prop != nil ***REMOVED***
			if !all ***REMOVED***
				if prop, ok := prop.(*valueProperty); ok && !prop.enumerable ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			accum = append(accum, asciiString(name))
		***REMOVED***
	***REMOVED***
	return a.baseObject.stringKeys(all, accum)
***REMOVED***

func (a *arrayObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		return idx < uint32(len(a.values)) && a.values[idx] != nil
	***REMOVED*** else ***REMOVED***
		return a.baseObject.hasOwnPropertyStr(name)
	***REMOVED***
***REMOVED***

func (a *arrayObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		return idx < uint32(len(a.values)) && a.values[idx] != nil
	***REMOVED***
	return a.baseObject.hasOwnPropertyStr(idx.string())
***REMOVED***

func (a *arrayObject) expand(idx uint32) bool ***REMOVED***
	targetLen := idx + 1
	if targetLen > uint32(len(a.values)) ***REMOVED***
		if targetLen < uint32(cap(a.values)) ***REMOVED***
			a.values = a.values[:targetLen]
		***REMOVED*** else ***REMOVED***
			if idx > 4096 && (a.objCount == 0 || idx/uint32(a.objCount) > 10) ***REMOVED***
				//log.Println("Switching standard->sparse")
				sa := &sparseArrayObject***REMOVED***
					baseObject:     a.baseObject,
					length:         a.length,
					propValueCount: a.propValueCount,
				***REMOVED***
				sa.setValues(a.values, a.objCount+1)
				sa.val.self = sa
				sa.lengthProp.writable = a.lengthProp.writable
				sa._put("length", &sa.lengthProp)
				return false
			***REMOVED*** else ***REMOVED***
				if bits.UintSize == 32 ***REMOVED***
					if targetLen >= math.MaxInt32 ***REMOVED***
						panic(a.val.runtime.NewTypeError("Array index overflows int"))
					***REMOVED***
				***REMOVED***
				tl := int(targetLen)
				newValues := make([]Value, tl, growCap(tl, len(a.values), cap(a.values)))
				copy(newValues, a.values)
				a.values = newValues
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (r *Runtime) defineArrayLength(prop *valueProperty, descr PropertyDescriptor, setter func(uint32, bool) bool, throw bool) bool ***REMOVED***
	var newLen uint32
	ret := true
	if descr.Value != nil ***REMOVED***
		newLen = r.toLengthUint32(descr.Value)
	***REMOVED***

	if descr.Configurable == FLAG_TRUE || descr.Enumerable == FLAG_TRUE || descr.Getter != nil || descr.Setter != nil ***REMOVED***
		ret = false
		goto Reject
	***REMOVED***

	if descr.Value != nil ***REMOVED***
		oldLen := uint32(prop.value.ToInteger())
		if oldLen != newLen ***REMOVED***
			ret = setter(newLen, false)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ret = true
	***REMOVED***

	if descr.Writable != FLAG_NOT_SET ***REMOVED***
		w := descr.Writable.Bool()
		if prop.writable ***REMOVED***
			prop.writable = w
		***REMOVED*** else ***REMOVED***
			if w ***REMOVED***
				ret = false
				goto Reject
			***REMOVED***
		***REMOVED***
	***REMOVED***

Reject:
	if !ret ***REMOVED***
		r.typeErrorResult(throw, "Cannot redefine property: length")
	***REMOVED***

	return ret
***REMOVED***

func (a *arrayObject) _defineIdxProperty(idx uint32, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	var existing Value
	if idx < uint32(len(a.values)) ***REMOVED***
		existing = a.values[idx]
	***REMOVED***
	prop, ok := a.baseObject._defineOwnProperty(unistring.String(strconv.FormatUint(uint64(idx), 10)), existing, desc, throw)
	if ok ***REMOVED***
		if idx >= a.length ***REMOVED***
			if !a.setLengthInt(idx+1, throw) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if a.expand(idx) ***REMOVED***
			a.values[idx] = prop
			a.objCount++
			if _, ok := prop.(*valueProperty); ok ***REMOVED***
				a.propValueCount++
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			a.val.self.(*sparseArrayObject).add(idx, prop)
		***REMOVED***
	***REMOVED***
	return ok
***REMOVED***

func (a *arrayObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._defineIdxProperty(idx, descr, throw)
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.val.runtime.defineArrayLength(a.getLengthProp(), descr, a.setLength, throw)
	***REMOVED***
	return a.baseObject.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (a *arrayObject) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		return a._defineIdxProperty(idx, descr, throw)
	***REMOVED***
	return a.baseObject.defineOwnPropertyStr(idx.string(), descr, throw)
***REMOVED***

func (a *arrayObject) _deleteIdxProp(idx uint32, throw bool) bool ***REMOVED***
	if idx < uint32(len(a.values)) ***REMOVED***
		if v := a.values[idx]; v != nil ***REMOVED***
			if p, ok := v.(*valueProperty); ok ***REMOVED***
				if !p.configurable ***REMOVED***
					a.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of %s", idx, a.val.toString())
					return false
				***REMOVED***
				a.propValueCount--
			***REMOVED***
			a.values[idx] = nil
			a.objCount--
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (a *arrayObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._deleteIdxProp(idx, throw)
	***REMOVED***
	return a.baseObject.deleteStr(name, throw)
***REMOVED***

func (a *arrayObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		return a._deleteIdxProp(idx, throw)
	***REMOVED***
	return a.baseObject.deleteStr(idx.string(), throw)
***REMOVED***

func (a *arrayObject) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, exists := ctx.get(a.val); exists ***REMOVED***
		return v
	***REMOVED***
	arr := make([]interface***REMOVED******REMOVED***, a.length)
	ctx.put(a.val, arr)
	if a.propValueCount == 0 && a.length == uint32(len(a.values)) && uint32(a.objCount) == a.length ***REMOVED***
		for i, v := range a.values ***REMOVED***
			if v != nil ***REMOVED***
				arr[i] = exportValue(v, ctx)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := uint32(0); i < a.length; i++ ***REMOVED***
			v := a.getIdx(valueInt(i), nil)
			if v != nil ***REMOVED***
				arr[i] = exportValue(v, ctx)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return arr
***REMOVED***

func (a *arrayObject) exportType() reflect.Type ***REMOVED***
	return reflectTypeArray
***REMOVED***

func (a *arrayObject) exportToArrayOrSlice(dst reflect.Value, typ reflect.Type, ctx *objectExportCtx) error ***REMOVED***
	r := a.val.runtime
	if iter := a.getSym(SymIterator, nil); iter == r.global.arrayValues || iter == nil ***REMOVED***
		l := toIntStrict(int64(a.length))
		if dst.Len() != l ***REMOVED***
			if typ.Kind() == reflect.Array ***REMOVED***
				return fmt.Errorf("cannot convert an Array into an array, lengths mismatch (have %d, need %d)", l, dst.Len())
			***REMOVED*** else ***REMOVED***
				dst.Set(reflect.MakeSlice(typ, l, l))
			***REMOVED***
		***REMOVED***
		ctx.putTyped(a.val, typ, dst.Interface())
		for i := 0; i < l; i++ ***REMOVED***
			if i >= len(a.values) ***REMOVED***
				break
			***REMOVED***
			val := a.values[i]
			if p, ok := val.(*valueProperty); ok ***REMOVED***
				val = p.get(a.val)
			***REMOVED***
			err := r.toReflectValue(val, dst.Index(i), ctx)
			if err != nil ***REMOVED***
				return fmt.Errorf("could not convert array element %v to %v at %d: %w", val, typ, i, err)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	return a.baseObject.exportToArrayOrSlice(dst, typ, ctx)
***REMOVED***

func (a *arrayObject) setValuesFromSparse(items []sparseArrayItem, newMaxIdx int) ***REMOVED***
	a.values = make([]Value, newMaxIdx+1)
	for _, item := range items ***REMOVED***
		a.values[item.idx] = item.value
	***REMOVED***
	a.objCount = len(items)
***REMOVED***

func toIdx(v valueInt) uint32 ***REMOVED***
	if v >= 0 && v < math.MaxUint32 ***REMOVED***
		return uint32(v)
	***REMOVED***
	return math.MaxUint32
***REMOVED***
