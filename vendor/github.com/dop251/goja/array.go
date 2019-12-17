package goja

import (
	"math"
	"reflect"
	"strconv"
)

type arrayObject struct ***REMOVED***
	baseObject
	values         []Value
	length         int64
	objCount       int64
	propValueCount int
	lengthProp     valueProperty
***REMOVED***

func (a *arrayObject) init() ***REMOVED***
	a.baseObject.init()
	a.lengthProp.writable = true

	a._put("length", &a.lengthProp)
***REMOVED***

func (a *arrayObject) _setLengthInt(l int64, throw bool) bool ***REMOVED***
	if l >= 0 && l <= math.MaxUint32 ***REMOVED***
		ret := true
		if l <= a.length ***REMOVED***
			if a.propValueCount > 0 ***REMOVED***
				// Slow path
				var s int64
				if a.length < int64(len(a.values)) ***REMOVED***
					s = a.length - 1
				***REMOVED*** else ***REMOVED***
					s = int64(len(a.values)) - 1
				***REMOVED***
				for i := s; i >= l; i-- ***REMOVED***
					if prop, ok := a.values[i].(*valueProperty); ok ***REMOVED***
						if !prop.configurable ***REMOVED***
							l = i + 1
							ret = false
							break
						***REMOVED***
						a.propValueCount--
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if l <= int64(len(a.values)) ***REMOVED***
			if l >= 16 && l < int64(cap(a.values))>>2 ***REMOVED***
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
	panic(a.val.runtime.newError(a.val.runtime.global.RangeError, "Invalid array length"))
***REMOVED***

func (a *arrayObject) setLengthInt(l int64, throw bool) bool ***REMOVED***
	if l == a.length ***REMOVED***
		return true
	***REMOVED***
	if !a.lengthProp.writable ***REMOVED***
		a.val.runtime.typeErrorResult(throw, "length is not writable")
		return false
	***REMOVED***
	return a._setLengthInt(l, throw)
***REMOVED***

func (a *arrayObject) setLength(v Value, throw bool) bool ***REMOVED***
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

func (a *arrayObject) getIdx(idx int64, origNameStr string, origName Value) (v Value) ***REMOVED***
	if idx >= 0 && idx < int64(len(a.values)) ***REMOVED***
		v = a.values[idx]
	***REMOVED***
	if v == nil && a.prototype != nil ***REMOVED***
		if origName != nil ***REMOVED***
			v = a.prototype.self.getProp(origName)
		***REMOVED*** else ***REMOVED***
			v = a.prototype.self.getPropStr(origNameStr)
		***REMOVED***
	***REMOVED***
	return
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

func toIdx(v Value) (idx int64) ***REMOVED***
	idx = -1
	if idxVal, ok1 := v.(valueInt); ok1 ***REMOVED***
		idx = int64(idxVal)
	***REMOVED*** else ***REMOVED***
		if i, err := strconv.ParseInt(v.String(), 10, 64); err == nil ***REMOVED***
			idx = i
		***REMOVED***
	***REMOVED***
	if idx >= 0 && idx < math.MaxUint32 ***REMOVED***
		return
	***REMOVED***
	return -1
***REMOVED***

func strToIdx(s string) (idx int64) ***REMOVED***
	idx = -1
	if i, err := strconv.ParseInt(s, 10, 64); err == nil ***REMOVED***
		idx = i
	***REMOVED***

	if idx >= 0 && idx < math.MaxUint32 ***REMOVED***
		return
	***REMOVED***
	return -1
***REMOVED***

func (a *arrayObject) getProp(n Value) Value ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return a.getIdx(idx, "", n)
	***REMOVED***

	if n.String() == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getProp(n)
***REMOVED***

func (a *arrayObject) getLengthProp() Value ***REMOVED***
	a.lengthProp.value = intToValue(a.length)
	return &a.lengthProp
***REMOVED***

func (a *arrayObject) getPropStr(name string) Value ***REMOVED***
	if i := strToIdx(name); i >= 0 ***REMOVED***
		return a.getIdx(i, name, nil)
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getPropStr(name)
***REMOVED***

func (a *arrayObject) getOwnProp(name string) Value ***REMOVED***
	if i := strToIdx(name); i >= 0 ***REMOVED***
		if i >= 0 && i < int64(len(a.values)) ***REMOVED***
			return a.values[i]
		***REMOVED***
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.getLengthProp()
	***REMOVED***
	return a.baseObject.getOwnProp(name)
***REMOVED***

func (a *arrayObject) putIdx(idx int64, val Value, throw bool, origNameStr string, origName Value) ***REMOVED***
	var prop Value
	if idx < int64(len(a.values)) ***REMOVED***
		prop = a.values[idx]
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
		if idx >= int64(len(a.values)) ***REMOVED***
			if !a.expand(idx) ***REMOVED***
				a.val.self.(*sparseArrayObject).putIdx(idx, val, throw, origNameStr, origName)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if prop, ok := prop.(*valueProperty); ok ***REMOVED***
			if !prop.isWritable() ***REMOVED***
				a.val.runtime.typeErrorResult(throw)
				return
			***REMOVED***
			prop.set(a.val, val)
			return
		***REMOVED***
	***REMOVED***

	a.values[idx] = val
	a.objCount++
***REMOVED***

func (a *arrayObject) put(n Value, val Value, throw bool) ***REMOVED***
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

func (a *arrayObject) putStr(name string, val Value, throw bool) ***REMOVED***
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

type arrayPropIter struct ***REMOVED***
	a         *arrayObject
	recursive bool
	idx       int
***REMOVED***

func (i *arrayPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.a.values) ***REMOVED***
		name := strconv.Itoa(i.idx)
		prop := i.a.values[i.idx]
		i.idx++
		if prop != nil ***REMOVED***
			return propIterItem***REMOVED***name: name, value: prop***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	return i.a.baseObject._enumerate(i.recursive)()
***REMOVED***

func (a *arrayObject) _enumerate(recursive bool) iterNextFunc ***REMOVED***
	return (&arrayPropIter***REMOVED***
		a:         a,
		recursive: recursive,
	***REMOVED***).next
***REMOVED***

func (a *arrayObject) enumerate(all, recursive bool) iterNextFunc ***REMOVED***
	return (&propFilterIter***REMOVED***
		wrapped: a._enumerate(recursive),
		all:     all,
		seen:    make(map[string]bool),
	***REMOVED***).next
***REMOVED***

func (a *arrayObject) hasOwnProperty(n Value) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return idx < int64(len(a.values)) && a.values[idx] != nil && a.values[idx] != _undefined
	***REMOVED*** else ***REMOVED***
		return a.baseObject.hasOwnProperty(n)
	***REMOVED***
***REMOVED***

func (a *arrayObject) hasOwnPropertyStr(name string) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		return idx < int64(len(a.values)) && a.values[idx] != nil && a.values[idx] != _undefined
	***REMOVED*** else ***REMOVED***
		return a.baseObject.hasOwnPropertyStr(name)
	***REMOVED***
***REMOVED***

func (a *arrayObject) expand(idx int64) bool ***REMOVED***
	targetLen := idx + 1
	if targetLen > int64(len(a.values)) ***REMOVED***
		if targetLen < int64(cap(a.values)) ***REMOVED***
			a.values = a.values[:targetLen]
		***REMOVED*** else ***REMOVED***
			if idx > 4096 && (a.objCount == 0 || idx/a.objCount > 10) ***REMOVED***
				//log.Println("Switching standard->sparse")
				sa := &sparseArrayObject***REMOVED***
					baseObject:     a.baseObject,
					length:         a.length,
					propValueCount: a.propValueCount,
				***REMOVED***
				sa.setValues(a.values)
				sa.val.self = sa
				sa.init()
				sa.lengthProp.writable = a.lengthProp.writable
				return false
			***REMOVED*** else ***REMOVED***
				// Use the same algorithm as in runtime.growSlice
				newcap := int64(cap(a.values))
				doublecap := newcap + newcap
				if targetLen > doublecap ***REMOVED***
					newcap = targetLen
				***REMOVED*** else ***REMOVED***
					if len(a.values) < 1024 ***REMOVED***
						newcap = doublecap
					***REMOVED*** else ***REMOVED***
						for newcap < targetLen ***REMOVED***
							newcap += newcap / 4
						***REMOVED***
					***REMOVED***
				***REMOVED***
				newValues := make([]Value, targetLen, newcap)
				copy(newValues, a.values)
				a.values = newValues
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (r *Runtime) defineArrayLength(prop *valueProperty, descr propertyDescr, setter func(Value, bool) bool, throw bool) bool ***REMOVED***
	ret := true

	if descr.Configurable == FLAG_TRUE || descr.Enumerable == FLAG_TRUE || descr.Getter != nil || descr.Setter != nil ***REMOVED***
		ret = false
		goto Reject
	***REMOVED***

	if newLen := descr.Value; newLen != nil ***REMOVED***
		ret = setter(newLen, false)
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

func (a *arrayObject) defineOwnProperty(n Value, descr propertyDescr, throw bool) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		var existing Value
		if idx < int64(len(a.values)) ***REMOVED***
			existing = a.values[idx]
		***REMOVED***
		prop, ok := a.baseObject._defineOwnProperty(n, existing, descr, throw)
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
				a.val.self.(*sparseArrayObject).putIdx(idx, prop, throw, "", nil)
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

func (a *arrayObject) _deleteProp(idx int64, throw bool) bool ***REMOVED***
	if idx < int64(len(a.values)) ***REMOVED***
		if v := a.values[idx]; v != nil ***REMOVED***
			if p, ok := v.(*valueProperty); ok ***REMOVED***
				if !p.configurable ***REMOVED***
					a.val.runtime.typeErrorResult(throw, "Cannot delete property '%d' of %s", idx, a.val.ToString())
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

func (a *arrayObject) delete(n Value, throw bool) bool ***REMOVED***
	if idx := toIdx(n); idx >= 0 ***REMOVED***
		return a._deleteProp(idx, throw)
	***REMOVED***
	return a.baseObject.delete(n, throw)
***REMOVED***

func (a *arrayObject) deleteStr(name string, throw bool) bool ***REMOVED***
	if idx := strToIdx(name); idx >= 0 ***REMOVED***
		return a._deleteProp(idx, throw)
	***REMOVED***
	return a.baseObject.deleteStr(name, throw)
***REMOVED***

func (a *arrayObject) export() interface***REMOVED******REMOVED*** ***REMOVED***
	arr := make([]interface***REMOVED******REMOVED***, a.length)
	for i, v := range a.values ***REMOVED***
		if v != nil ***REMOVED***
			arr[i] = v.Export()
		***REMOVED***
	***REMOVED***

	return arr
***REMOVED***

func (a *arrayObject) exportType() reflect.Type ***REMOVED***
	return reflectTypeArray
***REMOVED***

func (a *arrayObject) setValuesFromSparse(items []sparseArrayItem) ***REMOVED***
	a.values = make([]Value, int(items[len(items)-1].idx+1))
	for _, item := range items ***REMOVED***
		a.values[item.idx] = item.value
	***REMOVED***
	a.objCount = int64(len(items))
***REMOVED***
