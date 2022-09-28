// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package set provides simple set data structures for uint64s.
package set

import "math/bits"

// int64s represents a set of integers within the range of 0..63.
type int64s uint64

func (bs *int64s) Len() int ***REMOVED***
	return bits.OnesCount64(uint64(*bs))
***REMOVED***
func (bs *int64s) Has(n uint64) bool ***REMOVED***
	return uint64(*bs)&(uint64(1)<<n) > 0
***REMOVED***
func (bs *int64s) Set(n uint64) ***REMOVED***
	*(*uint64)(bs) |= uint64(1) << n
***REMOVED***
func (bs *int64s) Clear(n uint64) ***REMOVED***
	*(*uint64)(bs) &^= uint64(1) << n
***REMOVED***

// Ints represents a set of integers within the range of 0..math.MaxUint64.
type Ints struct ***REMOVED***
	lo int64s
	hi map[uint64]struct***REMOVED******REMOVED***
***REMOVED***

func (bs *Ints) Len() int ***REMOVED***
	return bs.lo.Len() + len(bs.hi)
***REMOVED***
func (bs *Ints) Has(n uint64) bool ***REMOVED***
	if n < 64 ***REMOVED***
		return bs.lo.Has(n)
	***REMOVED***
	_, ok := bs.hi[n]
	return ok
***REMOVED***
func (bs *Ints) Set(n uint64) ***REMOVED***
	if n < 64 ***REMOVED***
		bs.lo.Set(n)
		return
	***REMOVED***
	if bs.hi == nil ***REMOVED***
		bs.hi = make(map[uint64]struct***REMOVED******REMOVED***)
	***REMOVED***
	bs.hi[n] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
***REMOVED***
func (bs *Ints) Clear(n uint64) ***REMOVED***
	if n < 64 ***REMOVED***
		bs.lo.Clear(n)
		return
	***REMOVED***
	delete(bs.hi, n)
***REMOVED***
