package goja

import (
	"math"
	"reflect"
	"sort"
	"strconv"
)

type sparseArrayItem struct ***REMOVED***
	idx   int64
	value Value
***REMOVED***

type sparseArrayObject struct ***REMOVED***
	baseObject
	items          []sparseArrayItem
	length         int64
	propValueCount int
	lengthProp     valueProperty
***REMOVED***

func (a *sparseArrayObject) init() ***REMOVED***
	a.baseObject.init()
	a.lengthProp.writable = true

	a._put("length", &a.lengthProp)
***REMOVED***

func (a *sparseArrayObject) findIdx(idx int64) int ***REMOVED***
	return sort.Search(len(a.items), func(i int) bool ***REMOVED***
		return a.items[i].idx >= idx
	***REMOVED***)
***REMOVED***

func (a *sparseArrayObject) _setLengthInt(l int64, throw bool) bool ***REMOVED***
	if l >= 0 && l <= math.MaxUint32 ***REMOVED***
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
	panic(a.val.runtime.newError(a.val.runtime.global.RangeError, "Invalid array length"))
***REMOVED***

func (a *sparseArrayObject) setLengthInt(l int64, throw bool) bool ***REMOVED***
	if l == a.length ***REMOVED***
		return true
	***REMOVED***
	if !a.lengthProp.writable ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	***REMOVED***
	return a._setLengthInt(l, throw)
***REMOVED***

func (a *sparseArrayObject) setLength(v Value, throw bool) bool ***REMOVED***
	l, ok := toIntIgnoreNegZero(v)
	if ok && l == a.length ***REMOVED***
		return true
	***REMOVED***
	if !a.lengthProp.writable ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	***REMOVED***
	if ok ***REMOVED***
		return a._setLengthInt(l, throw)
	***REMOVED***
	panic(a.val.runtime.newError(a.val.runtime.global.RangeError, "Invalid array length"))
***REMOVED***

func (a *sparseArrayObject) getIdx(idx int64, origNameStr string, origName Value) (v Value) ***REMOVED***
	i := a.findIdx(idx)
	if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
		return a.items[i].value
	***REMOVED***

	if a.prototype != nil ***REMOVED***
		if origName != nil ***REMOVED***
			v = a.prototype.self.getProp(origName)
		***REMOVED*** else ***REMOVED***
			v = a.prototype.self.getPropStr(origNameStr)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (a *sparseArrayObject) getProp(n Value) Value ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return a.getIdx(idx, "", n)
	***REMOVED***

	if n.String() == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getProp(n)
***REMOVED***

func (a *sparseArrayObject) getLengthProp() Value ***REMOVED***
	a.lengthProp.value = intToValue(a.length)
	return &a.lengthProp
***REMOVED***

func (a *sparseArrayObject) getOwnProp(name string) Value ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		i := a.findIdx(idx)
		if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
			return a.items[i].value
		***REMOVED***
		return nil
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getOwnProp(name)
***REMOVED***

func (a *sparseArrayObject) getPropStr(name string) Value ***REMOVED***
	if i := strToIdx(name); i >= 0 ***REMOVED***
		return a.getIdx(i, name, nil)
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getPropStr(name)
***REMOVED***

func (a *sparseArrayObject) putIdx(idx int64, val Value, throw bool, origNameStr string, origName Value) ***REMOVED***
	var prop Value
	i := a.findIdx(idx)
	if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
		prop = a.items[i].value
	***REMOVED***

	if prop == nil ***REMOVED***
		if a.prototype != nil ***REMOVED***
			var pprop Value
			if origName != nil ***REMOVED***
				pprop = a.prototype.self.getProp(origName)
			***REMOVED*** else ***REMOVED***
				pprop = a.prototype.self.getPropStr(origNameStr)
			***REMOVED***
			if pprop, ok := pprop.(*valueProperty); ok ***REMOVED***
				if !pprop.isWritable() ***REMOVED***
					a.val.runtime.typeErrorResult(throw)
					return
				***REMOVED***
				if pprop.accessor ***REMOVED***
					pprop.set(a.val, val)
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if !a.extensible ***REMOVED***
			a.val.runtime.typeErrorResult(throw)
			return
		***REMOVED***

		if idx >= a.length ***REMOVED***
			if !a.setLengthInt(idx+1, throw) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***

		if a.expand() ***REMOVED***
			a.items = append(a.items, sparseArrayItem***REMOVED******REMOVED***)
			copy(a.items[i+1:], a.items[i:])
			a.items[i] = sparseArrayItem***REMOVED***
				idx:   idx,
				value: val,
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			a.val.self.(*arrayObject).putIdx(idx, val, throw, origNameStr, origName)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if prop, ok := prop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				a.val.runtime.typeErrorResult(throw)
				return
			***REMOVED***
			prop.set(a.val, val)
			return
		***REMOVED*** else ***REMOVED***
			a.items[i].value = val
		***REMOVED***
	***REMOVED***

***REMOVED***

func (a *sparseArrayObject) put(n Value, val Value, throw bool) ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		a.putIdx(idx, val, throw, "", n)
	***REMOVED*** else ***REMOVED***
		if n.String() == "length" ***REMOVED***
			a.setLength(val, throw)
		***REMOVED*** else ***REMOVED***
			a.baseObject.put(n, val, throw)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) putStr(name string, val Value, throw bool) ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		a.putIdx(idx, val, throw, name, nil)
	***REMOVED*** else ***REMOVED***
		if name == "length" ***REMOVED***
			a.setLength(val, throw)
		***REMOVED*** else ***REMOVED***
			a.baseObject.putStr(name, val, throw)
		***REMOVED***
	***REMOVED***
***REMOVED***

type sparseArrayPropIter struct ***REMOVED***
	a         *sparseArrayObject
	recursive bool
	idx       int
***REMOVED***

func (i *sparseArrayPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.a.items) ***REMOVED***
		name := strconv.Itoa(int(i.a.items[i.idx].idx))
		prop := i.a.items[i.idx].value
		i.idx++
		if prop != nil ***REMOVED***
			return propIterItem***REMOVED***name: name, value: prop***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	return i.a.baseObject._enumerate(i.recursive)()
***REMOVED***

func (a *sparseArrayObject) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	return (&sparseArrayPropIter***REMOVED***
		a:         a,
		recursive: recursive,
	***REMOVED***).next
***REMOVED***

func (a *sparseArrayObject) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: a._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next
***REMOVED***

func (a *sparseArrayObject) setValues(values []Value) ***REMOVED***
	a.items = nil
	for i, val := range values ***REMOVED***
		if val != nil ***REMOVED***
			a.items = append(a.items, sparseArrayItem***REMOVED***
				idx:   int64(i),
				value: val,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) hasOwnProperty(n Value) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		i := a.findIdx(idx)
		if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
			return a.items[i].value != _undefined
		***REMOVED***
		return false
	***REMOVED*** else ***REMOVED***
		return a.baseObject.hasOwnProperty(n)
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) hasOwnPropertyStr(name string) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		i := a.findIdx(idx)
		if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
			return a.items[i].value != _undefined
		***REMOVED***
		return false
	***REMOVED*** else ***REMOVED***
		return a.baseObject.hasOwnPropertyStr(name)
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) expand() bool ***REMOVED***
	if l := len(a.items); l >= 1024 ***REMOVED***
		if int(a.items[l-1].idx)/l < 8 ***REMOVED***
			//log.Println("Switching sparse->standard")
			ar := &arrayObject***REMOVED***
				baseObject:     a.baseObject,
				length:         a.length,
				propValueCount: a.propValueCount,
			***REMOVED***
			ar.setValuesFromSparse(a.items)
			ar.val.self = ar
			ar.init()
			ar.lengthProp.writable = a.lengthProp.writable
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (a *sparseArrayObject) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		var existing Value
		i := a.findIdx(idx)
		if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
			existing = a.items[i].value
		***REMOVED***
		prop, ok := a.baseObject._defineOwnProperty(n, existing, descr, throw)
		if ok ***REMOVED***
			if idx >= a.length ***REMOVED***
				if !a.setLengthInt(idx+1, throw) ***REMOVED***
					return false
				***REMOVED***
			***REMOVED***
			if i >= len(a.items) || a.items[i].idx != idx ***REMOVED***
				if a.expand() ***REMOVED***
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
					return a.val.self.defineOwnProperty(n, descr, throw)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				a.items[i].value = prop
			***REMOVED***
			if _, ok := prop.(*valueProperty); ok ***REMOVED***
				a.propValueCount++
			***REMOVED***
		***REMOVED***
		return ok
	***REMOVED*** else ***REMOVED***
		if n.String() == "length" ***REMOVED***
			return a.val.runtime.defineArrayLength(&a.lengthProp, descr, a.setLength, throw)
		***REMOVED***
		return a.baseObject.defineOwnProperty(n, descr, throw)
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) _deleteProp(idx int64, throw bool) bool ***REMOVED***
	i := a.findIdx(idx)
	if i < len(a.items) && a.items[i].idx == idx ***REMOVED***
		if p, ok := a.items[i].value.(*valueProperty); ok ***REMOVED***
			if !p.configurable ***REMOVED***
				a.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of %s", idx, a.val.ToString())
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

func (a *sparseArrayObject) delete(n Value, throw bool) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return a._deleteProp(idx, throw)
	***REMOVED***
	return a.baseObject.delete(n, throw)
***REMOVED***

func (a *sparseArrayObject) deleteStr(name string, throw bool) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		return a._deleteProp(idx, throw)
	***REMOVED***
	return a.baseObject.deleteStr(name, throw)
***REMOVED***

func (a *sparseArrayObject) sortLen() int64 ***REMOVED***
	if len(a.items) > 0 ***REMOVED***
		return a.items[len(a.items)-1].idx + 1
	***REMOVED***

	return 0
***REMOVED***

func (a *sparseArrayObject) sortGet(i int64) Value ***REMOVED***
	idx := a.findIdx(i)
	if idx < len(a.items) && a.items[idx].idx == i ***REMOVED***
		v := a.items[idx].value
		if p, ok := v.(*valueProperty); ok ***REMOVED***
			v = p.get(a.val)
		***REMOVED***
		return v
	***REMOVED***
	return nil
***REMOVED***

func (a *sparseArrayObject) swap(i, j int64) ***REMOVED***
	idxI := a.findIdx(i)
	idxJ := a.findIdx(j)

	if idxI < len(a.items) && a.items[idxI].idx == i && idxJ < len(a.items) && a.items[idxJ].idx == j ***REMOVED***
		a.items[idxI].value, a.items[idxJ].value = a.items[idxJ].value, a.items[idxI].value
	***REMOVED***
***REMOVED***

func (a *sparseArrayObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	arr := make([]interface***REMOVED******REMOVED***, a.length)
	for _, item := range a.items ***REMOVED***
		if item.value != nil ***REMOVED***
			arr[item.idx] = item.value.Export()
		***REMOVED***
	***REMOVED***
	return arr
***REMOVED***

func (a *sparseArrayObject) exportType() reflect.Type ***REMOVED***
	return reflectTypeArray
***REMOVED***
