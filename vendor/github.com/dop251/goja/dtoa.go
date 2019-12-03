package goja

// Ported from Rhino (https://github.com/mozilla/rhino/blob/master/src/org/mozilla/javascript/DToA.java)

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strconv"
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

	digits = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func lo0bits(x uint32) (k uint32) ***REMOVED***

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

func hi0bits(x uint32) (k uint32) ***REMOVED***

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

func d2b(d float64) (b *big.Int, e int32, bits uint32) ***REMOVED***
	dBits := math.Float64bits(d)
	d0 := uint32(dBits >> 32)
	d1 := uint32(dBits)

	z := d0 & frac_mask
	d0 &= 0x7fffffff /* clear sign bit, which we ignore */

	var de, k, i uint32
	var dbl_bits []byte
	if de = (d0 >> exp_shift); de != 0 ***REMOVED***
		z |= exp_msk1
	***REMOVED***

	y := d1
	if y != 0 ***REMOVED***
		dbl_bits = make([]byte, 8)
		k = lo0bits(y)
		y >>= k
		if k != 0 ***REMOVED***
			stuffBits(dbl_bits, 4, y|z<<(32-k))
			z >>= k
		***REMOVED*** else ***REMOVED***
			stuffBits(dbl_bits, 4, y)
		***REMOVED***
		stuffBits(dbl_bits, 0, z)
		if z != 0 ***REMOVED***
			i = 2
		***REMOVED*** else ***REMOVED***
			i = 1
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		dbl_bits = make([]byte, 4)
		k = lo0bits(z)
		z >>= k
		stuffBits(dbl_bits, 0, z)
		k += 32
		i = 1
	***REMOVED***

	if de != 0 ***REMOVED***
		e = int32(de - bias - (p - 1) + k)
		bits = p - k
	***REMOVED*** else ***REMOVED***
		e = int32(de - bias - (p - 1) + 1 + k)
		bits = 32*i - hi0bits(z)
	***REMOVED***
	b = (&big.Int***REMOVED******REMOVED***).SetBytes(dbl_bits)
	return
***REMOVED***

func dtobasestr(num float64, radix int) string ***REMOVED***
	var negative bool
	if num < 0 ***REMOVED***
		num = -num
		negative = true
	***REMOVED***

	dfloor := math.Floor(num)
	ldfloor := int64(dfloor)
	var intDigits string
	if dfloor == float64(ldfloor) ***REMOVED***
		if negative ***REMOVED***
			ldfloor = -ldfloor
		***REMOVED***
		intDigits = strconv.FormatInt(ldfloor, radix)
	***REMOVED*** else ***REMOVED***
		floorBits := math.Float64bits(num)
		exp := int(floorBits>>exp_shiftL) & exp_mask_shifted
		var mantissa int64
		if exp == 0 ***REMOVED***
			mantissa = int64((floorBits & frac_maskL) << 1)
		***REMOVED*** else ***REMOVED***
			mantissa = int64((floorBits & frac_maskL) | exp_msk1L)
		***REMOVED***

		if negative ***REMOVED***
			mantissa = -mantissa
		***REMOVED***
		exp -= 1075
		x := big.NewInt(mantissa)
		if exp > 0 ***REMOVED***
			x.Lsh(x, uint(exp))
		***REMOVED*** else if exp < 0 ***REMOVED***
			x.Rsh(x, uint(-exp))
		***REMOVED***
		intDigits = x.Text(radix)
	***REMOVED***

	if num == dfloor ***REMOVED***
		// No fraction part
		return intDigits
	***REMOVED*** else ***REMOVED***
		/* We have a fraction. */
		var buffer bytes.Buffer
		buffer.WriteString(intDigits)
		buffer.WriteByte('.')
		df := num - dfloor

		dBits := math.Float64bits(num)
		word0 := uint32(dBits >> 32)
		word1 := uint32(dBits)

		b, e, _ := d2b(df)
		//            JS_ASSERT(e < 0);
		/* At this point df = b * 2^e.  e must be less than zero because 0 < df < 1. */

		s2 := -int32((word0 >> exp_shift1) & (exp_mask >> exp_shift1))
		if s2 == 0 ***REMOVED***
			s2 = -1
		***REMOVED***
		s2 += bias + p
		/* 1/2^s2 = (nextDouble(d) - d)/2 */
		//            JS_ASSERT(-s2 < e);
		if -s2 >= e ***REMOVED***
			panic(fmt.Errorf("-s2 >= e: %d, %d", -s2, e))
		***REMOVED***
		mlo := big.NewInt(1)
		mhi := mlo
		if (word1 == 0) && ((word0 & bndry_mask) == 0) && ((word0 & (exp_mask & (exp_mask << 1))) != 0) ***REMOVED***
			/* The special case.  Here we want to be within a quarter of the last input
			   significant digit instead of one half of it when the output string's value is less than d.  */
			s2 += log2P
			mhi = big.NewInt(1 << log2P)
		***REMOVED***

		b.Lsh(b, uint(e+s2))
		s := big.NewInt(1)
		s.Lsh(s, uint(s2))
		/* At this point we have the following:
		 *   s = 2^s2;
		 *   1 > df = b/2^s2 > 0;
		 *   (d - prevDouble(d))/2 = mlo/2^s2;
		 *   (nextDouble(d) - d)/2 = mhi/2^s2. */
		bigBase := big.NewInt(int64(radix))

		done := false
		m := &big.Int***REMOVED******REMOVED***
		delta := &big.Int***REMOVED******REMOVED***
		for !done ***REMOVED***
			b.Mul(b, bigBase)
			b.DivMod(b, s, m)
			digit := byte(b.Int64())
			b, m = m, b
			mlo.Mul(mlo, bigBase)
			if mlo != mhi ***REMOVED***
				mhi.Mul(mhi, bigBase)
			***REMOVED***

			/* Do we yet have the shortest string that will round to d? */
			j := b.Cmp(mlo)
			/* j is b/2^s2 compared with mlo/2^s2. */

			delta.Sub(s, mhi)
			var j1 int
			if delta.Sign() <= 0 ***REMOVED***
				j1 = 1
			***REMOVED*** else ***REMOVED***
				j1 = b.Cmp(delta)
			***REMOVED***
			/* j1 is b/2^s2 compared with 1 - mhi/2^s2. */
			if j1 == 0 && (word1&1) == 0 ***REMOVED***
				if j > 0 ***REMOVED***
					digit++
				***REMOVED***
				done = true
			***REMOVED*** else if j < 0 || (j == 0 && ((word1 & 1) == 0)) ***REMOVED***
				if j1 > 0 ***REMOVED***
					/* Either dig or dig+1 would work here as the least significant digit.
					Use whichever would produce an output value closer to d. */
					b.Lsh(b, 1)
					j1 = b.Cmp(s)
					if j1 > 0 ***REMOVED*** /* The even test (|| (j1 == 0 && (digit & 1))) is not here because it messes up odd base output such as 3.5 in base 3.  */
						digit++
					***REMOVED***
				***REMOVED***
				done = true
			***REMOVED*** else if j1 > 0 ***REMOVED***
				digit++
				done = true
			***REMOVED***
			//                JS_ASSERT(digit < (uint32)base);
			buffer.WriteByte(digits[digit])
		***REMOVED***

		return buffer.String()
	***REMOVED***
***REMOVED***
