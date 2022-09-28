package goja

import (
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"sort"
	"strconv"

	"github.com/dop251/goja/unistring"
)

type sparseArrayItem struct ***REMOVED***
	idx   uint32
	value Value
***REMOVED***

type sparseArrayObject struct ***REMOVED***
	baseObject
	items          []sparseArrayItem
	length         uint32
	propValueCount int
	lengthProp     valueProperty
***REMOVED***

func (a *sparseArrayObject) findIdx(idx uint32) int ***REMOVED***
	return sort.Search(len(a.items), func(i int) bool ***REMOVED***
		return a.items[i].idx >= idx
	***REMOVED***)
***REMOVED***

func (a *sparseArrayObject) _setLengthInt(l uint32, throw bool) bool ***REMOVED***
	ret := true
	if l <= a.length ***REMOVED***
		if a.propValueCount > 0 ***REMOVED***
			// Slow path
			for i := len(a.items) - 1; i >= 0; i-- ***REMOVED***
				item := a.items[i]
				if item.idx <= l ***REMOVED***
					break
				***REMOVED***
				if prop, ok := item.value.(*valueProperty); ok ***REMOVED***
					if !prop.configurable ***REMOVED***
						l = item.idx + 1
						ret = false
						break
					***REMOVED***
					a.propValueCount--
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	idx := a.findIdx(l)

	aa := a.items[idx:]
	for i := range aa ***REMOVED***
		aa[i].value = nil
	***REMOVED***
	a.items = a.items[:idx]
	a.length = l
	if !ret ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "Cannot redefine property: length")
	***REMOVED***
	return ret
***REMOVED***

func (a *sparseArrayObject) setLengthInt(l uint32, throw bool) bool ***REMOVED***
	if l == a.length ***REMOVED***
		return true
	***REMOVED***
	if !a.lengthProp.writable ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	***REMOVED***
	return a._setLengthInt(l, throw)
***REMOVED***

func (a *sparseArrayObject) setLength(v uint32, throw bool) bool ***REMOVED***
	if !a.lengthProp.writable ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	***REMOVED***
	return a._setLengthInt(v, throw)
***REMOVED***

func (a *sparseArrayObject) _getIdx(idx uint32) Value ***REMOVED***
	i := a.findIdx(idx)
	if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
		return a.items[i].value
	***REMOVED***

	return nil
***REMOVED***

func (a *sparseArrayObject) getStr(name unistring.String, receiver Value) Value ***REMOVED***
	return a.getStrWithOwnProp(a.getOwnPropStr(name), name, receiver)
***REMOVED***

func (a *sparseArrayObject) getIdx(idx valueInt, receiver Value) Value ***REMOVED***
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

func (a *sparseArrayObject) getLengthProp() *valueProperty ***REMOVED***
	a.lengthProp.value = intToValue(int64(a.length))
	return &a.lengthProp
***REMOVED***

func (a *sparseArrayObject) getOwnPropStr(name unistring.String) Value ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._getIdx(idx)
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getOwnPropStr(name)
***REMOVED***

func (a *sparseArrayObject) getOwnPropIdx(idx valueInt) Value ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		return a._getIdx(idx)
	***REMOVED***
	return a.baseObject.getOwnPropStr(idx.string())
***REMOVED***

func (a *sparseArrayObject) add(idx uint32, val Value) ***REMOVED***
	i := a.findIdx(idx)
	a.items = append(a.items, sparseArrayItem***REMOVED******REMOVED***)
	copy(a.items[i+1:], a.items[i:])
	a.items[i] = sparseArrayItem***REMOVED***
		idx:   idx,
		value: val,
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) _setOwnIdx(idx uint32, val Value, throw bool) bool ***REMOVED***
	var prop Value
	i := a.findIdx(idx)
	if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
		prop = a.items[i].value
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
		***REMOVED***

		if idx >= a.length ***REMOVED***
			if !a.setLengthInt(idx+1, throw) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***

		if a.expand(idx) ***REMOVED***
			a.items = append(a.items, sparseArrayItem***REMOVED******REMOVED***)
			copy(a.items[i+1:], a.items[i:])
			a.items[i] = sparseArrayItem***REMOVED***
				idx:   idx,
				value: val,
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			ar := a.val.self.(*arrayObject)
			ar.values[idx] = val
			ar.objCount++
			return true
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if prop, ok := prop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				a.val.runtime.typeErrorResult(throw)
				return false
			***REMOVED***
			prop.set(a.val, val)
		***REMOVED*** else ***REMOVED***
			a.items[i].value = val
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (a *sparseArrayObject) setOwnStr(name unistring.String, val Value, throw bool) bool ***REMOVED***
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

func (a *sparseArrayObject) setOwnIdx(idx valueInt, val Value, throw bool) bool ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		return a._setOwnIdx(idx, val, throw)
	***REMOVED***

	return a.baseObject.setOwnStr(idx.string(), val, throw)
***REMOVED***

func (a *sparseArrayObject) setForeignStr(name unistring.String, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return a._setForeignStr(name, a.getOwnPropStr(name), val, receiver, throw)
***REMOVED***

func (a *sparseArrayObject) setForeignIdx(name valueInt, val, receiver Value, throw bool) (bool, bool) ***REMOVED***
	return a._setForeignIdx(name, a.getOwnPropIdx(name), val, receiver, throw)
***REMOVED***

type sparseArrayPropIter struct ***REMOVED***
	a   *sparseArrayObject
	idx int
***REMOVED***

func (i *sparseArrayPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.a.items) ***REMOVED***
		name := asciiString(strconv.Itoa(int(i.a.items[i.idx].idx)))
		prop := i.a.items[i.idx].value
		i.idx++
		if prop != nil ***REMOVED***
			return propIterItem***REMOVED***name: name, value: prop***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	return i.a.baseObject.iterateStringKeys()()
***REMOVED***

func (a *sparseArrayObject) iterateStringKeys() iterNextFunc ***REMOVED***
	return (&sparseArrayPropIter***REMOVED***
		a: a,
	***REMOVED***).next
***REMOVED***

func (a *sparseArrayObject) stringKeys(all bool, accum []Value) []Value ***REMOVED***
	if all ***REMOVED***
		for _, item := range a.items ***REMOVED***
			accum = append(accum, asciiString(strconv.FormatUint(uint64(item.idx), 10)))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, item := range a.items ***REMOVED***
			if prop, ok := item.value.(*valueProperty); ok && !prop.enumerable ***REMOVED***
				continue
			***REMOVED***
			accum = append(accum, asciiString(strconv.FormatUint(uint64(item.idx), 10)))
		***REMOVED***
	***REMOVED***

	return a.baseObject.stringKeys(all, accum)
***REMOVED***

func (a *sparseArrayObject) setValues(values []Value, objCount int) ***REMOVED***
	a.items = make([]sparseArrayItem, 0, objCount)
	for i, val := range values ***REMOVED***
		if val != nil ***REMOVED***
			a.items = append(a.items, sparseArrayItem***REMOVED***
				idx:   uint32(i),
				value: val,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		i := a.findIdx(idx)
		return i < len(a.items) && a.items[i].idx == idx
	***REMOVED*** else ***REMOVED***
		return a.baseObject.hasOwnPropertyStr(name)
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) hasOwnPropertyIdx(idx valueInt) bool ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		i := a.findIdx(idx)
		return i < len(a.items) && a.items[i].idx == idx
	***REMOVED***

	return a.baseObject.hasOwnPropertyStr(idx.string())
***REMOVED***

func (a *sparseArrayObject) expand(idx uint32) bool ***REMOVED***
	if l := len(a.items); l >= 1024 ***REMOVED***
		if ii := a.items[l-1].idx; ii > idx ***REMOVED***
			idx = ii
		***REMOVED***
		if (bits.UintSize == 64 || idx < math.MaxInt32) && int(idx)>>3 < l ***REMOVED***
			//log.Println("Switching sparse->standard")
			ar := &arrayObject***REMOVED***
				baseObject:     a.baseObject,
				length:         a.length,
				propValueCount: a.propValueCount,
			***REMOVED***
			ar.setValuesFromSparse(a.items, int(idx))
			ar.val.self = ar
			ar.lengthProp.writable = a.lengthProp.writable
			a._put("length", &ar.lengthProp)
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (a *sparseArrayObject) _defineIdxProperty(idx uint32, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	var existing Value
	i := a.findIdx(idx)
	if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
		existing = a.items[i].value
	***REMOVED***
	prop, ok := a.baseObject._defineOwnProperty(unistring.String(strconv.FormatUint(uint64(idx), 10)), existing, desc, throw)
	if ok ***REMOVED***
		if idx >= a.length ***REMOVED***
			if !a.setLengthInt(idx+1, throw) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		if i >= len(a.items) || a.items[i].idx != idx ***REMOVED***
			if a.expand(idx) ***REMOVED***
				a.items = append(a.items, sparseArrayItem***REMOVED******REMOVED***)
				copy(a.items[i+1:], a.items[i:])
				a.items[i] = sparseArrayItem***REMOVED***
					idx:   idx,
					value: prop,
				***REMOVED***
				if idx >= a.length ***REMOVED***
					a.length = idx + 1
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				a.val.self.(*arrayObject).values[idx] = prop
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			a.items[i].value = prop
		***REMOVED***
		if _, ok := prop.(*valueProperty); ok ***REMOVED***
			a.propValueCount++
		***REMOVED***
	***REMOVED***
	return ok
***REMOVED***

func (a *sparseArrayObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._defineIdxProperty(idx, descr, throw)
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.val.runtime.defineArrayLength(a.getLengthProp(), descr, a.setLength, throw)
	***REMOVED***
	return a.baseObject.defineOwnPropertyStr(name, descr, throw)
***REMOVED***

func (a *sparseArrayObject) defineOwnPropertyIdx(idx valueInt, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		return a._defineIdxProperty(idx, descr, throw)
	***REMOVED***
	return a.baseObject.defineOwnPropertyStr(idx.string(), descr, throw)
***REMOVED***

func (a *sparseArrayObject) _deleteIdxProp(idx uint32, throw bool) bool ***REMOVED***
	i := a.findIdx(idx)
	if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
		if p, ok := a.items[i].value.(*valueProperty); ok ***REMOVED***
			if !p.configurable ***REMOVED***
				a.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of %s", idx, a.val.toString())
				return false
			***REMOVED***
			a.propValueCount--
		***REMOVED***
		copy(a.items[i:], a.items[i+1:])
		a.items[len(a.items)-1].value = nil
		a.items = a.items[:len(a.items)-1]
	***REMOVED***
	return true
***REMOVED***

func (a *sparseArrayObject) deleteStr(name unistring.String, throw bool) bool ***REMOVED***
	if idx := strToArrayIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._deleteIdxProp(idx, throw)
	***REMOVED***
	return a.baseObject.deleteStr(name, throw)
***REMOVED***

func (a *sparseArrayObject) deleteIdx(idx valueInt, throw bool) bool ***REMOVED***
	if idx := toIdx(idx); idx != math.MaxUint32 ***REMOVED***
		return a._deleteIdxProp(idx, throw)
	***REMOVED***
	return a.baseObject.deleteStr(idx.string(), throw)
***REMOVED***

func (a *sparseArrayObject) sortLen() int ***REMOVED***
	if len(a.items) > 0 ***REMOVED***
		return toIntStrict(int64(a.items[len(a.items)-1].idx) + 1)
	***REMOVED***

	return 0
***REMOVED***

func (a *sparseArrayObject) export(ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	if v, exists := ctx.get(a.val); exists ***REMOVED***
		return v
	***REMOVED***
	arr := make([]interface***REMOVED******REMOVED***, a.length)
	ctx.put(a.val, arr)
	var prevIdx uint32
	for _, item := range a.items ***REMOVED***
		idx := item.idx
		for i := prevIdx; i < idx; i++ ***REMOVED***
			if a.prototype != nil ***REMOVED***
				if v := a.prototype.self.getIdx(valueInt(i), nil); v != nil ***REMOVED***
					arr[i] = exportValue(v, ctx)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		v := item.value
		if v != nil ***REMOVED***
			if prop, ok := v.(*valueProperty); ok ***REMOVED***
				v = prop.get(a.val)
			***REMOVED***
			arr[idx] = exportValue(v, ctx)
		***REMOVED***
		prevIdx = idx + 1
	***REMOVED***
	for i := prevIdx; i < a.length; i++ ***REMOVED***
		if a.prototype != nil ***REMOVED***
			if v := a.prototype.self.getIdx(valueInt(i), nil); v != nil ***REMOVED***
				arr[i] = exportValue(v, ctx)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return arr
***REMOVED***

func (a *sparseArrayObject) exportType() reflect.Type ***REMOVED***
	return reflectTypeArray
***REMOVED***

func (a *sparseArrayObject) exportToArrayOrSlice(dst reflect.Value, typ reflect.Type, ctx *objectExportCtx) error ***REMOVED***
	r := a.val.runtime
	if iter := a.getSym(SymIterator, nil); iter == r.global.arrayValues || iter == nil ***REMOVED***
		l := toIntStrict(int64(a.length))
		if typ.Kind() == reflect.Array ***REMOVED***
			if dst.Len() != l ***REMOVED***
				return fmt.Errorf("cannot convert an Array into an array, lengths mismatch (have %d, need %d)", l, dst.Len())
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			dst.Set(reflect.MakeSlice(typ, l, l))
		***REMOVED***
		ctx.putTyped(a.val, typ, dst.Interface())
		for _, item := range a.items ***REMOVED***
			val := item.value
			if p, ok := val.(*valueProperty); ok ***REMOVED***
				val = p.get(a.val)
			***REMOVED***
			idx := toIntStrict(int64(item.idx))
			if idx >= l ***REMOVED***
				break
			***REMOVED***
			err := r.toReflectValue(val, dst.Index(idx), ctx)
			if err != nil ***REMOVED***
				return fmt.Errorf("could not convert array element %v to %v at %d: %w", item.value, typ, idx, err)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	return a.baseObject.exportToArrayOrSlice(dst, typ, ctx)
***REMOVED***
