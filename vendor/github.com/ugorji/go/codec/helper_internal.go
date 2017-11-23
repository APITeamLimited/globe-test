// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// All non-std package dependencies live in this file,
// so porting to different environment is easy (just update functions).

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

func panicValToErr(panicVal interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	if panicVal == nil ***REMOVED***
		return
	***REMOVED***
	// case nil
	switch xerr := panicVal.(type) ***REMOVED***
	case error:
		*err = xerr
	case string:
		*err = errors.New(xerr)
	default:
		*err = fmt.Errorf("%v", panicVal)
	***REMOVED***
	return
***REMOVED***

func hIsEmptyValue(v reflect.Value, deref, checkStruct bool) bool ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Invalid:
		return true
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		if deref ***REMOVED***
			if v.IsNil() ***REMOVED***
				return true
			***REMOVED***
			return hIsEmptyValue(v.Elem(), deref, checkStruct)
		***REMOVED*** else ***REMOVED***
			return v.IsNil()
		***REMOVED***
	case reflect.Struct:
		if !checkStruct ***REMOVED***
			return false
		***REMOVED***
		// return true if all fields are empty. else return false.
		// we cannot use equality check, because some fields may be maps/slices/etc
		// and consequently the structs are not comparable.
		// return v.Interface() == reflect.Zero(v.Type()).Interface()
		for i, n := 0, v.NumField(); i < n; i++ ***REMOVED***
			if !hIsEmptyValue(v.Field(i), deref, checkStruct) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func isEmptyValue(v reflect.Value, deref, checkStruct bool) bool ***REMOVED***
	return hIsEmptyValue(v, deref, checkStruct)
***REMOVED***

func pruneSignExt(v []byte, pos bool) (n int) ***REMOVED***
	if len(v) < 2 ***REMOVED***
	***REMOVED*** else if pos && v[0] == 0 ***REMOVED***
		for ; v[n] == 0 && n+1 < len(v) && (v[n+1]&(1<<7) == 0); n++ ***REMOVED***
		***REMOVED***
	***REMOVED*** else if !pos && v[0] == 0xff ***REMOVED***
		for ; v[n] == 0xff && n+1 < len(v) && (v[n+1]&(1<<7) != 0); n++ ***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func implementsIntf(typ, iTyp reflect.Type) (success bool, indir int8) ***REMOVED***
	if typ == nil ***REMOVED***
		return
	***REMOVED***
	rt := typ
	// The type might be a pointer and we need to keep
	// dereferencing to the base type until we find an implementation.
	for ***REMOVED***
		if rt.Implements(iTyp) ***REMOVED***
			return true, indir
		***REMOVED***
		if p := rt; p.Kind() == reflect.Ptr ***REMOVED***
			indir++
			if indir >= math.MaxInt8 ***REMOVED*** // insane number of indirections
				return false, 0
			***REMOVED***
			rt = p.Elem()
			continue
		***REMOVED***
		break
	***REMOVED***
	// No luck yet, but if this is a base type (non-pointer), the pointer might satisfy.
	if typ.Kind() != reflect.Ptr ***REMOVED***
		// Not a pointer, but does the pointer work?
		if reflect.PtrTo(typ).Implements(iTyp) ***REMOVED***
			return true, -1
		***REMOVED***
	***REMOVED***
	return false, 0
***REMOVED***

// validate that this function is correct ...
// culled from OGRE (Object-Oriented Graphics Rendering Engine)
// function: halfToFloatI (http://stderr.org/doc/ogre-doc/api/OgreBitwise_8h-source.html)
func halfFloatToFloatBits(yy uint16) (d uint32) ***REMOVED***
	y := uint32(yy)
	s := (y >> 15) & 0x01
	e := (y >> 10) & 0x1f
	m := y & 0x03ff

	if e == 0 ***REMOVED***
		if m == 0 ***REMOVED*** // plu or minus 0
			return s << 31
		***REMOVED*** else ***REMOVED*** // Denormalized number -- renormalize it
			for (m & 0x00000400) == 0 ***REMOVED***
				m <<= 1
				e -= 1
			***REMOVED***
			e += 1
			const zz uint32 = 0x0400
			m &= ^zz
		***REMOVED***
	***REMOVED*** else if e == 31 ***REMOVED***
		if m == 0 ***REMOVED*** // Inf
			return (s << 31) | 0x7f800000
		***REMOVED*** else ***REMOVED*** // NaN
			return (s << 31) | 0x7f800000 | (m << 13)
		***REMOVED***
	***REMOVED***
	e = e + (127 - 15)
	m = m << 13
	return (s << 31) | (e << 23) | m
***REMOVED***

// GrowCap will return a new capacity for a slice, given the following:
//   - oldCap: current capacity
//   - unit: in-memory size of an element
//   - num: number of elements to add
func growCap(oldCap, unit, num int) (newCap int) ***REMOVED***
	// appendslice logic (if cap < 1024, *2, else *1.25):
	//   leads to many copy calls, especially when copying bytes.
	//   bytes.Buffer model (2*cap + n): much better for bytes.
	// smarter way is to take the byte-size of the appended element(type) into account

	// maintain 3 thresholds:
	// t1: if cap <= t1, newcap = 2x
	// t2: if cap <= t2, newcap = 1.75x
	// t3: if cap <= t3, newcap = 1.5x
	//     else          newcap = 1.25x
	//
	// t1, t2, t3 >= 1024 always.
	// i.e. if unit size >= 16, then always do 2x or 1.25x (ie t1, t2, t3 are all same)
	//
	// With this, appending for bytes increase by:
	//    100% up to 4K
	//     75% up to 8K
	//     50% up to 16K
	//     25% beyond that

	// unit can be 0 e.g. for struct***REMOVED******REMOVED******REMOVED******REMOVED***; handle that appropriately
	var t1, t2, t3 int // thresholds
	if unit <= 1 ***REMOVED***
		t1, t2, t3 = 4*1024, 8*1024, 16*1024
	***REMOVED*** else if unit < 16 ***REMOVED***
		t3 = 16 / unit * 1024
		t1 = t3 * 1 / 4
		t2 = t3 * 2 / 4
	***REMOVED*** else ***REMOVED***
		t1, t2, t3 = 1024, 1024, 1024
	***REMOVED***

	var x int // temporary variable

	// x is multiplier here: one of 5, 6, 7 or 8; incr of 25%, 50%, 75% or 100% respectively
	if oldCap <= t1 ***REMOVED*** // [0,t1]
		x = 8
	***REMOVED*** else if oldCap > t3 ***REMOVED*** // (t3,infinity]
		x = 5
	***REMOVED*** else if oldCap <= t2 ***REMOVED*** // (t1,t2]
		x = 7
	***REMOVED*** else ***REMOVED*** // (t2,t3]
		x = 6
	***REMOVED***
	newCap = x * oldCap / 4

	if num > 0 ***REMOVED***
		newCap += num
	***REMOVED***

	// ensure newCap is a multiple of 64 (if it is > 64) or 16.
	if newCap > 64 ***REMOVED***
		if x = newCap % 64; x != 0 ***REMOVED***
			x = newCap / 64
			newCap = 64 * (x + 1)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if x = newCap % 16; x != 0 ***REMOVED***
			x = newCap / 16
			newCap = 16 * (x + 1)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func expandSliceValue(s reflect.Value, num int) reflect.Value ***REMOVED***
	if num <= 0 ***REMOVED***
		return s
	***REMOVED***
	l0 := s.Len()
	l1 := l0 + num // new slice length
	if l1 < l0 ***REMOVED***
		panic("ExpandSlice: slice overflow")
	***REMOVED***
	c0 := s.Cap()
	if l1 <= c0 ***REMOVED***
		return s.Slice(0, l1)
	***REMOVED***
	st := s.Type()
	c1 := growCap(c0, int(st.Elem().Size()), num)
	s2 := reflect.MakeSlice(st, l1, c1)
	// println("expandslicevalue: cap-old: ", c0, ", cap-new: ", c1, ", len-new: ", l1)
	reflect.Copy(s2, s)
	return s2
***REMOVED***
