// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

type valueRange struct ***REMOVED***
	value  uint16 // header: value:stride
	lo, hi byte   // header: lo:n
***REMOVED***

type sparseBlocks struct ***REMOVED***
	values []valueRange
	offset []uint16
***REMOVED***

var nfcSparse = sparseBlocks***REMOVED***
	values: nfcSparseValues[:],
	offset: nfcSparseOffset[:],
***REMOVED***

var nfkcSparse = sparseBlocks***REMOVED***
	values: nfkcSparseValues[:],
	offset: nfkcSparseOffset[:],
***REMOVED***

var (
	nfcData  = newNfcTrie(0)
	nfkcData = newNfkcTrie(0)
)

// lookupValue determines the type of block n and looks up the value for b.
// For n < t.cutoff, the block is a simple lookup table. Otherwise, the block
// is a list of ranges with an accompanying value. Given a matching range r,
// the value for b is by r.value + (b - r.lo) * stride.
func (t *sparseBlocks) lookup(n uint32, b byte) uint16 ***REMOVED***
	offset := t.offset[n]
	header := t.values[offset]
	lo := offset + 1
	hi := lo + uint16(header.lo)
	for lo < hi ***REMOVED***
		m := lo + (hi-lo)/2
		r := t.values[m]
		if r.lo <= b && b <= r.hi ***REMOVED***
			return r.value + uint16(b-r.lo)*header.value
		***REMOVED***
		if b < r.lo ***REMOVED***
			hi = m
		***REMOVED*** else ***REMOVED***
			lo = m + 1
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***
