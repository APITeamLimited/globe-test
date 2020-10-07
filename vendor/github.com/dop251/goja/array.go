package goja

import (
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
	elementValue := ai.obj.self.getIdx(idxVal, nil)
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

func (a *arrayObject) _setLengthInt(l int64, throw bool) bool ***REMOVED***
	if l >= 0 && l <= math.MaxUint32 ***REMOVED***
		l := uint32(l)
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
	panic(a.val.runtime.newError(a.val.runtime.global.RangeError, "Invalid array length"))
***REMOVED***

func (a *arrayObject) setLengthInt(l int64, throw bool) bool ***REMOVED***
	if l == int64(a.length) ***REMOVED***
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
	if ok && l == int64(a.length) ***REMOVED***
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
	if i := strToIdx(name); i != math.MaxUint32 ***REMOVED***
		if i < uint32(len(a.values)) ***REMOVED***
			return a.values[i]
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

func (a *arrayObject) getLengthProp() Value ***REMOVED***
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
				if !a.setLengthInt(int64(idx)+1, throw) ***REMOVED***
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
	if idx := strToIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._setOwnIdx(idx, val, throw)
	***REMOVED*** else ***REMOVED***
		if name == "length" ***REMOVED***
			return a.setLength(val, throw)
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
	a   *arrayObject
	idx int
***REMOVED***

func (i *arrayPropIter) next() (propIterItem, iterNextFunc) ***REMOVED***
	for i.idx < len(i.a.values) ***REMOVED***
		name := unistring.String(strconv.Itoa(i.idx))
		prop := i.a.values[i.idx]
		i.idx++
		if prop != nil ***REMOVED***
			return propIterItem***REMOVED***name: name, value: prop***REMOVED***, i.next
		***REMOVED***
	***REMOVED***

	return i.a.baseObject.enumerateUnfiltered()()
***REMOVED***

func (a *arrayObject) enumerateUnfiltered() iterNextFunc ***REMOVED***
	return (&arrayPropIter***REMOVED***
		a: a,
	***REMOVED***).next
***REMOVED***

func (a *arrayObject) ownKeys(all bool, accum []Value) []Value ***REMOVED***
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
	return a.baseObject.ownKeys(all, accum)
***REMOVED***

func (a *arrayObject) hasOwnPropertyStr(name unistring.String) bool ***REMOVED***
	if idx := strToIdx(name); idx != math.MaxUint32 ***REMOVED***
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
					length:         uint32(a.length),
					propValueCount: a.propValueCount,
				***REMOVED***
				sa.setValues(a.values, a.objCount+1)
				sa.val.self = sa
				sa.init()
				sa.lengthProp.writable = a.lengthProp.writable
				return false
			***REMOVED*** else ***REMOVED***
				if bits.UintSize == 32 ***REMOVED***
					if targetLen >= math.MaxInt32 ***REMOVED***
						panic(a.val.runtime.NewTypeError("Array index overflows int"))
					***REMOVED***
				***REMOVED***
				tl := int(targetLen)
				// Use the same algorithm as in runtime.growSlice
				newcap := cap(a.values)
				doublecap := newcap + newcap
				if tl > doublecap ***REMOVED***
					newcap = tl
				***REMOVED*** else ***REMOVED***
					if len(a.values) < 1024 ***REMOVED***
						newcap = doublecap
					***REMOVED*** else ***REMOVED***
						for newcap < tl ***REMOVED***
							newcap += newcap / 4
						***REMOVED***
					***REMOVED***
				***REMOVED***
				newValues := make([]Value, tl, newcap)
				copy(newValues, a.values)
				a.values = newValues
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (r *Runtime) defineArrayLength(prop *valueProperty, descr PropertyDescriptor, setter func(Value, bool) bool, throw bool) bool ***REMOVED***
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

func (a *arrayObject) _defineIdxProperty(idx uint32, desc PropertyDescriptor, throw bool) bool ***REMOVED***
	var existing Value
	if idx < uint32(len(a.values)) ***REMOVED***
		existing = a.values[idx]
	***REMOVED***
	prop, ok := a.baseObject._defineOwnProperty(unistring.String(strconv.FormatUint(uint64(idx), 10)), existing, desc, throw)
	if ok ***REMOVED***
		if idx >= a.length ***REMOVED***
			if !a.setLengthInt(int64(idx)+1, throw) ***REMOVED***
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
			a.val.self.(*sparseArrayObject).add(uint32(idx), prop)
		***REMOVED***
	***REMOVED***
	return ok
***REMOVED***

func (a *arrayObject) defineOwnPropertyStr(name unistring.String, descr PropertyDescriptor, throw bool) bool ***REMOVED***
	if idx := strToIdx(name); idx != math.MaxUint32 ***REMOVED***
		return a._defineIdxProperty(idx, descr, throw)
	***REMOVED***
	if name == "length" ***REMOVED***
		return a.val.runtime.defineArrayLength(&a.lengthProp, descr, a.setLength, throw)
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
	if idx := strToIdx(name); idx != math.MaxUint32 ***REMOVED***
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
	if v, exists := ctx.get(a); exists ***REMOVED***
		return v
	***REMOVED***
	arr := make([]interface***REMOVED******REMOVED***, a.length)
	ctx.put(a, arr)
	for i, v := range a.values ***REMOVED***
		if v != nil ***REMOVED***
			arr[i] = exportValue(v, ctx)
		***REMOVED***
	***REMOVED***
	return arr
***REMOVED***

func (a *arrayObject) exportType() reflect.Type ***REMOVED***
	return reflectTypeArray
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

func strToIdx64(s unistring.String) int64 ***REMOVED***
	if s == "" ***REMOVED***
		return -1
	***REMOVED***
	l := len(s)
	if s[0] == '0' ***REMOVED***
		if l == 1 ***REMOVED***
			return 0
		***REMOVED***
		return -1
	***REMOVED***
	var n int64
	if l < 19 ***REMOVED***
		// guaranteed not to overflow
		for i := 0; i < len(s); i++ ***REMOVED***
			c := s[i]
			if c < '0' || c > '9' ***REMOVED***
				return -1
			***REMOVED***
			n = n*10 + int64(c-'0')
		***REMOVED***
		return n
	***REMOVED***
	if l > 19 ***REMOVED***
		// guaranteed to overflow
		return -1
	***REMOVED***
	c18 := s[18]
	if c18 < '0' || c18 > '9' ***REMOVED***
		return -1
	***REMOVED***
	for i := 0; i < 18; i++ ***REMOVED***
		c := s[i]
		if c < '0' || c > '9' ***REMOVED***
			return -1
		***REMOVED***
		n = n*10 + int64(c-'0')
	***REMOVED***
	if n >= math.MaxInt64/10+1 ***REMOVED***
		return -1
	***REMOVED***
	n *= 10
	n1 := n + int64(c18-'0')
	if n1 < n ***REMOVED***
		return -1
	***REMOVED***
	return n1
***REMOVED***

func strToIdx(s unistring.String) uint32 ***REMOVED***
	if s == "" ***REMOVED***
		return math.MaxUint32
	***REMOVED***
	l := len(s)
	if s[0] == '0' ***REMOVED***
		if l == 1 ***REMOVED***
			return 0
		***REMOVED***
		return math.MaxUint32
	***REMOVED***
	var n uint32
	if l < 10 ***REMOVED***
		// guaranteed not to overflow
		for i := 0; i < len(s); i++ ***REMOVED***
			c := s[i]
			if c < '0' || c > '9' ***REMOVED***
				return math.MaxUint32
			***REMOVED***
			n = n*10 + uint32(c-'0')
		***REMOVED***
		return n
	***REMOVED***
	if l > 10 ***REMOVED***
		// guaranteed to overflow
		return math.MaxUint32
	***REMOVED***
	c9 := s[9]
	if c9 < '0' || c9 > '9' ***REMOVED***
		return math.MaxUint32
	***REMOVED***
	for i := 0; i < 9; i++ ***REMOVED***
		c := s[i]
		if c < '0' || c > '9' ***REMOVED***
			return math.MaxUint32
		***REMOVED***
		n = n*10 + uint32(c-'0')
	***REMOVED***
	if n >= math.MaxUint32/10+1 ***REMOVED***
		return math.MaxUint32
	***REMOVED***
	n *= 10
	n1 := n + uint32(c9-'0')
	if n1 < n ***REMOVED***
		return math.MaxUint32
	***REMOVED***

	return n1
***REMOVED***

func strToGoIdx(s unistring.String) int ***REMOVED***
	if bits.UintSize == 64 ***REMOVED***
		return int(strToIdx64(s))
	***REMOVED***
	i := strToIdx(s)
	if i >= math.MaxInt32 ***REMOVED***
		return -1
	***REMOVED***
	return int(i)
***REMOVED***
