/*
Package ftoa provides ECMAScript-compliant floating point number conversion to string.

It contains code ported from Rhino (https://github.com/mozilla/rhino/blob/master/src/org/mozilla/javascript/DToA.java)
as well as from the original code by David M. Gay.

See LICENSE_LUCENE for the original copyright message and disclaimer.

*/
package ftoa

import (
	"math"
)

const (
	frac_mask = 0xfffff
	exp_shift = 20
	exp_msk1  = 0x100000

	exp_shiftL       = 52
	exp_mask_shifted = 0x7ff
	frac_maskL       = 0xfffffffffffff
	exp_msk1L        = 0x10000000000000
	exp_shift1       = 20
	exp_mask         = 0x7ff00000
	bias             = 1023
	p                = 53
	bndry_mask       = 0xfffff
	log2P            = 1
)

func lo0bits(x uint32) (k int) ***REMOVED***

	if (x & 7) != 0 ***REMOVED***
		if (x & 1) != 0 ***REMOVED***
			return 0
		***REMOVED***
		if (x & 2) != 0 ***REMOVED***
			return 1
		***REMOVED***
		return 2
	***REMOVED***
	if (x & 0xffff) == 0 ***REMOVED***
		k = 16
		x >>= 16
	***REMOVED***
	if (x & 0xff) == 0 ***REMOVED***
		k += 8
		x >>= 8
	***REMOVED***
	if (x & 0xf) == 0 ***REMOVED***
		k += 4
		x >>= 4
	***REMOVED***
	if (x & 0x3) == 0 ***REMOVED***
		k += 2
		x >>= 2
	***REMOVED***
	if (x & 1) == 0 ***REMOVED***
		k++
		x >>= 1
		if (x & 1) == 0 ***REMOVED***
			return 32
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func hi0bits(x uint32) (k int) ***REMOVED***

	if (x & 0xffff0000) == 0 ***REMOVED***
		k = 16
		x <<= 16
	***REMOVED***
	if (x & 0xff000000) == 0 ***REMOVED***
		k += 8
		x <<= 8
	***REMOVED***
	if (x & 0xf0000000) == 0 ***REMOVED***
		k += 4
		x <<= 4
	***REMOVED***
	if (x & 0xc0000000) == 0 ***REMOVED***
		k += 2
		x <<= 2
	***REMOVED***
	if (x & 0x80000000) == 0 ***REMOVED***
		k++
		if (x & 0x40000000) == 0 ***REMOVED***
			return 32
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func stuffBits(bits []byte, offset int, val uint32) ***REMOVED***
	bits[offset] = byte(val >> 24)
	bits[offset+1] = byte(val >> 16)
	bits[offset+2] = byte(val >> 8)
	bits[offset+3] = byte(val)
***REMOVED***

func d2b(d float64, b []byte) (e, bits int, dblBits []byte) ***REMOVED***
	dBits := math.Float64bits(d)
	d0 := uint32(dBits >> 32)
	d1 := uint32(dBits)

	z := d0 & frac_mask
	d0 &= 0x7fffffff /* clear sign bit, which we ignore */

	var de, k, i int
	if de = int(d0 >> exp_shift); de != 0 ***REMOVED***
		z |= exp_msk1
	***REMOVED***

	y := d1
	if y != 0 ***REMOVED***
		dblBits = b[:8]
		k = lo0bits(y)
		y >>= k
		if k != 0 ***REMOVED***
			stuffBits(dblBits, 4, y|z<<(32-k))
			z >>= k
		***REMOVED*** else ***REMOVED***
			stuffBits(dblBits, 4, y)
		***REMOVED***
		stuffBits(dblBits, 0, z)
		if z != 0 ***REMOVED***
			i = 2
		***REMOVED*** else ***REMOVED***
			i = 1
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		dblBits = b[:4]
		k = lo0bits(z)
		z >>= k
		stuffBits(dblBits, 0, z)
		k += 32
		i = 1
	***REMOVED***

	if de != 0 ***REMOVED***
		e = de - bias - (p - 1) + k
		bits = p - k
	***REMOVED*** else ***REMOVED***
		e = de - bias - (p - 1) + 1 + k
		bits = 32*i - hi0bits(z)
	***REMOVED***
	return
***REMOVED***
