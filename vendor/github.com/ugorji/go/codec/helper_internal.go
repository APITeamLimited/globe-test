// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// All non-std package dependencies live in this file,
// so porting to different environment is easy (just update functions).

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
		***REMOVED***
		// Denormalized number -- renormalize it
		for (m & 0x00000400) == 0 ***REMOVED***
			m <<= 1
			e -= 1
		***REMOVED***
		e += 1
		const zz uint32 = 0x0400
		m &= ^zz
	***REMOVED*** else if e == 31 ***REMOVED***
		if m == 0 ***REMOVED*** // Inf
			return (s << 31) | 0x7f800000
		***REMOVED***
		return (s << 31) | 0x7f800000 | (m << 13) // NaN
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
	if newCap <= oldCap ***REMOVED***
		newCap = oldCap + 1
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
