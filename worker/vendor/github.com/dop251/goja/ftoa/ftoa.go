package ftoa

import (
	"math"
	"math/big"
)

const (
	exp_11     = 0x3ff00000
	frac_mask1 = 0xfffff
	bletch     = 0x10
	quick_max  = 14
	int_max    = 14
)

var (
	tens = [...]float64***REMOVED***
		1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
		1e20, 1e21, 1e22,
	***REMOVED***

	bigtens = [...]float64***REMOVED***1e16, 1e32, 1e64, 1e128, 1e256***REMOVED***

	big5  = big.NewInt(5)
	big10 = big.NewInt(10)

	p05       = []*big.Int***REMOVED***big5, big.NewInt(25), big.NewInt(125)***REMOVED***
	pow5Cache [7]*big.Int

	dtoaModes = []int***REMOVED***
		ModeStandard:            0,
		ModeStandardExponential: 0,
		ModeFixed:               3,
		ModeExponential:         2,
		ModePrecision:           2,
	***REMOVED***
)

/*
d must be > 0 and must not be Inf

mode:
		0 ==> shortest string that yields d when read in
			and rounded to nearest.
		1 ==> like 0, but with Steele & White stopping rule;
			e.g. with IEEE P754 arithmetic , mode 0 gives
			1e23 whereas mode 1 gives 9.999999999999999e22.
		2 ==> max(1,ndigits) significant digits.  This gives a
			return value similar to that of ecvt, except
			that trailing zeros are suppressed.
		3 ==> through ndigits past the decimal point.  This
			gives a return value similar to that from fcvt,
			except that trailing zeros are suppressed, and
			ndigits can be negative.
		4,5 ==> similar to 2 and 3, respectively, but (in
			round-nearest mode) with the tests of mode 0 to
			possibly return a shorter string that rounds to d.
			With IEEE arithmetic and compilation with
			-DHonor_FLT_ROUNDS, modes 4 and 5 behave the same
			as modes 2 and 3 when FLT_ROUNDS != 1.
		6-9 ==> Debugging modes similar to mode - 4:  don't try
			fast floating-point estimate (if applicable).

		Values of mode other than 0-9 are treated as mode 0.
*/
func ftoa(d float64, mode int, biasUp bool, ndigits int, buf []byte) ([]byte, int) ***REMOVED***
	startPos := len(buf)
	dblBits := make([]byte, 0, 8)
	be, bbits, dblBits := d2b(d, dblBits)

	dBits := math.Float64bits(d)
	word0 := uint32(dBits >> 32)
	word1 := uint32(dBits)

	i := int((word0 >> exp_shift1) & (exp_mask >> exp_shift1))
	var d2 float64
	var denorm bool
	if i != 0 ***REMOVED***
		d2 = setWord0(d, (word0&frac_mask1)|exp_11)
		i -= bias
		denorm = false
	***REMOVED*** else ***REMOVED***
		/* d is denormalized */
		i = bbits + be + (bias + (p - 1) - 1)
		var x uint64
		if i > 32 ***REMOVED***
			x = uint64(word0)<<(64-i) | uint64(word1)>>(i-32)
		***REMOVED*** else ***REMOVED***
			x = uint64(word1) << (32 - i)
		***REMOVED***
		d2 = setWord0(float64(x), uint32((x>>32)-31*exp_mask))
		i -= (bias + (p - 1) - 1) + 1
		denorm = true
	***REMOVED***
	/* At this point d = f*2^i, where 1 <= f < 2.  d2 is an approximation of f. */
	ds := (d2-1.5)*0.289529654602168 + 0.1760912590558 + float64(i)*0.301029995663981
	k := int(ds)
	if ds < 0.0 && ds != float64(k) ***REMOVED***
		k-- /* want k = floor(ds) */
	***REMOVED***
	k_check := true
	if k >= 0 && k < len(tens) ***REMOVED***
		if d < tens[k] ***REMOVED***
			k--
		***REMOVED***
		k_check = false
	***REMOVED***
	/* At this point floor(log10(d)) <= k <= floor(log10(d))+1.
	   If k_check is zero, we're guaranteed that k = floor(log10(d)). */
	j := bbits - i - 1
	var b2, s2, b5, s5 int
	/* At this point d = b/2^j, where b is an odd integer. */
	if j >= 0 ***REMOVED***
		b2 = 0
		s2 = j
	***REMOVED*** else ***REMOVED***
		b2 = -j
		s2 = 0
	***REMOVED***
	if k >= 0 ***REMOVED***
		b5 = 0
		s5 = k
		s2 += k
	***REMOVED*** else ***REMOVED***
		b2 -= k
		b5 = -k
		s5 = 0
	***REMOVED***
	/* At this point d/10^k = (b * 2^b2 * 5^b5) / (2^s2 * 5^s5), where b is an odd integer,
	   b2 >= 0, b5 >= 0, s2 >= 0, and s5 >= 0. */
	if mode < 0 || mode > 9 ***REMOVED***
		mode = 0
	***REMOVED***
	try_quick := true
	if mode > 5 ***REMOVED***
		mode -= 4
		try_quick = false
	***REMOVED***
	leftright := true
	var ilim, ilim1 int
	switch mode ***REMOVED***
	case 0, 1:
		ilim, ilim1 = -1, -1
		ndigits = 0
	case 2:
		leftright = false
		fallthrough
	case 4:
		if ndigits <= 0 ***REMOVED***
			ndigits = 1
		***REMOVED***
		ilim, ilim1 = ndigits, ndigits
	case 3:
		leftright = false
		fallthrough
	case 5:
		i = ndigits + k + 1
		ilim = i
		ilim1 = i - 1
	***REMOVED***
	/* ilim is the maximum number of significant digits we want, based on k and ndigits. */
	/* ilim1 is the maximum number of significant digits we want, based on k and ndigits,
	   when it turns out that k was computed too high by one. */
	fast_failed := false
	if ilim >= 0 && ilim <= quick_max && try_quick ***REMOVED***

		/* Try to get by with floating-point arithmetic. */

		i = 0
		d2 = d
		k0 := k
		ilim0 := ilim
		ieps := 2 /* conservative */
		/* Divide d by 10^k, keeping track of the roundoff error and avoiding overflows. */
		if k > 0 ***REMOVED***
			ds = tens[k&0xf]
			j = k >> 4
			if (j & bletch) != 0 ***REMOVED***
				/* prevent overflows */
				j &= bletch - 1
				d /= bigtens[len(bigtens)-1]
				ieps++
			***REMOVED***
			for ; j != 0; i++ ***REMOVED***
				if (j & 1) != 0 ***REMOVED***
					ieps++
					ds *= bigtens[i]
				***REMOVED***
				j >>= 1
			***REMOVED***
			d /= ds
		***REMOVED*** else if j1 := -k; j1 != 0 ***REMOVED***
			d *= tens[j1&0xf]
			for j = j1 >> 4; j != 0; i++ ***REMOVED***
				if (j & 1) != 0 ***REMOVED***
					ieps++
					d *= bigtens[i]
				***REMOVED***
				j >>= 1
			***REMOVED***
		***REMOVED***
		/* Check that k was computed correctly. */
		if k_check && d < 1.0 && ilim > 0 ***REMOVED***
			if ilim1 <= 0 ***REMOVED***
				fast_failed = true
			***REMOVED*** else ***REMOVED***
				ilim = ilim1
				k--
				d *= 10.
				ieps++
			***REMOVED***
		***REMOVED***
		/* eps bounds the cumulative error. */
		eps := float64(ieps)*d + 7.0
		eps = setWord0(eps, _word0(eps)-(p-1)*exp_msk1)
		if ilim == 0 ***REMOVED***
			d -= 5.0
			if d > eps ***REMOVED***
				buf = append(buf, '1')
				k++
				return buf, k + 1
			***REMOVED***
			if d < -eps ***REMOVED***
				buf = append(buf, '0')
				return buf, 1
			***REMOVED***
			fast_failed = true
		***REMOVED***
		if !fast_failed ***REMOVED***
			fast_failed = true
			if leftright ***REMOVED***
				/* Use Steele & White method of only
				 * generating digits needed.
				 */
				eps = 0.5/tens[ilim-1] - eps
				for i = 0; ; ***REMOVED***
					l := int64(d)
					d -= float64(l)
					buf = append(buf, byte('0'+l))
					if d < eps ***REMOVED***
						return buf, k + 1
					***REMOVED***
					if 1.0-d < eps ***REMOVED***
						buf, k = bumpUp(buf, k)
						return buf, k + 1
					***REMOVED***
					i++
					if i >= ilim ***REMOVED***
						break
					***REMOVED***
					eps *= 10.0
					d *= 10.0
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				/* Generate ilim digits, then fix them up. */
				eps *= tens[ilim-1]
				for i = 1; ; i++ ***REMOVED***
					l := int64(d)
					d -= float64(l)
					buf = append(buf, byte('0'+l))
					if i == ilim ***REMOVED***
						if d > 0.5+eps ***REMOVED***
							buf, k = bumpUp(buf, k)
							return buf, k + 1
						***REMOVED*** else if d < 0.5-eps ***REMOVED***
							buf = stripTrailingZeroes(buf, startPos)
							return buf, k + 1
						***REMOVED***
						break
					***REMOVED***
					d *= 10.0
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if fast_failed ***REMOVED***
			buf = buf[:startPos]
			d = d2
			k = k0
			ilim = ilim0
		***REMOVED***
	***REMOVED***

	/* Do we have a "small" integer? */
	if be >= 0 && k <= int_max ***REMOVED***
		/* Yes. */
		ds = tens[k]
		if ndigits < 0 && ilim <= 0 ***REMOVED***
			if ilim < 0 || d < 5*ds || (!biasUp && d == 5*ds) ***REMOVED***
				buf = buf[:startPos]
				buf = append(buf, '0')
				return buf, 1
			***REMOVED***
			buf = append(buf, '1')
			k++
			return buf, k + 1
		***REMOVED***
		for i = 1; ; i++ ***REMOVED***
			l := int64(d / ds)
			d -= float64(l) * ds
			buf = append(buf, byte('0'+l))
			if i == ilim ***REMOVED***
				d += d
				if (d > ds) || (d == ds && (((l & 1) != 0) || biasUp)) ***REMOVED***
					buf, k = bumpUp(buf, k)
				***REMOVED***
				break
			***REMOVED***
			d *= 10.0
			if d == 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		return buf, k + 1
	***REMOVED***

	m2 := b2
	m5 := b5
	var mhi, mlo *big.Int
	if leftright ***REMOVED***
		if mode < 2 ***REMOVED***
			if denorm ***REMOVED***
				i = be + (bias + (p - 1) - 1 + 1)
			***REMOVED*** else ***REMOVED***
				i = 1 + p - bbits
			***REMOVED***
			/* i is 1 plus the number of trailing zero bits in d's significand. Thus,
			   (2^m2 * 5^m5) / (2^(s2+i) * 5^s5) = (1/2 lsb of d)/10^k. */
		***REMOVED*** else ***REMOVED***
			j = ilim - 1
			if m5 >= j ***REMOVED***
				m5 -= j
			***REMOVED*** else ***REMOVED***
				j -= m5
				s5 += j
				b5 += j
				m5 = 0
			***REMOVED***
			i = ilim
			if i < 0 ***REMOVED***
				m2 -= i
				i = 0
			***REMOVED***
			/* (2^m2 * 5^m5) / (2^(s2+i) * 5^s5) = (1/2 * 10^(1-ilim))/10^k. */
		***REMOVED***
		b2 += i
		s2 += i
		mhi = big.NewInt(1)
		/* (mhi * 2^m2 * 5^m5) / (2^s2 * 5^s5) = one-half of last printed (when mode >= 2) or
		   input (when mode < 2) significant digit, divided by 10^k. */
	***REMOVED***

	/* We still have d/10^k = (b * 2^b2 * 5^b5) / (2^s2 * 5^s5).  Reduce common factors in
	   b2, m2, and s2 without changing the equalities. */
	if m2 > 0 && s2 > 0 ***REMOVED***
		if m2 < s2 ***REMOVED***
			i = m2
		***REMOVED*** else ***REMOVED***
			i = s2
		***REMOVED***
		b2 -= i
		m2 -= i
		s2 -= i
	***REMOVED***

	b := new(big.Int).SetBytes(dblBits)
	/* Fold b5 into b and m5 into mhi. */
	if b5 > 0 ***REMOVED***
		if leftright ***REMOVED***
			if m5 > 0 ***REMOVED***
				pow5mult(mhi, m5)
				b.Mul(mhi, b)
			***REMOVED***
			j = b5 - m5
			if j != 0 ***REMOVED***
				pow5mult(b, j)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			pow5mult(b, b5)
		***REMOVED***
	***REMOVED***
	/* Now we have d/10^k = (b * 2^b2) / (2^s2 * 5^s5) and
	   (mhi * 2^m2) / (2^s2 * 5^s5) = one-half of last printed or input significant digit, divided by 10^k. */

	S := big.NewInt(1)
	if s5 > 0 ***REMOVED***
		pow5mult(S, s5)
	***REMOVED***
	/* Now we have d/10^k = (b * 2^b2) / (S * 2^s2) and
	   (mhi * 2^m2) / (S * 2^s2) = one-half of last printed or input significant digit, divided by 10^k. */

	/* Check for special case that d is a normalized power of 2. */
	spec_case := false
	if mode < 2 ***REMOVED***
		if (_word1(d) == 0) && ((_word0(d) & bndry_mask) == 0) &&
			((_word0(d) & (exp_mask & (exp_mask << 1))) != 0) ***REMOVED***
			/* The special case.  Here we want to be within a quarter of the last input
			   significant digit instead of one half of it when the decimal output string's value is less than d.  */
			b2 += log2P
			s2 += log2P
			spec_case = true
		***REMOVED***
	***REMOVED***

	/* Arrange for convenient computation of quotients:
	 * shift left if necessary so divisor has 4 leading 0 bits.
	 *
	 * Perhaps we should just compute leading 28 bits of S once
	 * and for all and pass them and a shift to quorem, so it
	 * can do shifts and ors to compute the numerator for q.
	 */
	var zz int
	if s5 != 0 ***REMOVED***
		S_bytes := S.Bytes()
		var S_hiWord uint32
		for idx := 0; idx < 4; idx++ ***REMOVED***
			S_hiWord = S_hiWord << 8
			if idx < len(S_bytes) ***REMOVED***
				S_hiWord |= uint32(S_bytes[idx])
			***REMOVED***
		***REMOVED***
		zz = 32 - hi0bits(S_hiWord)
	***REMOVED*** else ***REMOVED***
		zz = 1
	***REMOVED***
	i = (zz + s2) & 0x1f
	if i != 0 ***REMOVED***
		i = 32 - i
	***REMOVED***
	/* i is the number of leading zero bits in the most significant word of S*2^s2. */
	if i > 4 ***REMOVED***
		i -= 4
		b2 += i
		m2 += i
		s2 += i
	***REMOVED*** else if i < 4 ***REMOVED***
		i += 28
		b2 += i
		m2 += i
		s2 += i
	***REMOVED***
	/* Now S*2^s2 has exactly four leading zero bits in its most significant word. */
	if b2 > 0 ***REMOVED***
		b = b.Lsh(b, uint(b2))
	***REMOVED***
	if s2 > 0 ***REMOVED***
		S.Lsh(S, uint(s2))
	***REMOVED***
	/* Now we have d/10^k = b/S and
	   (mhi * 2^m2) / S = maximum acceptable error, divided by 10^k. */
	if k_check ***REMOVED***
		if b.Cmp(S) < 0 ***REMOVED***
			k--
			b.Mul(b, big10) /* we botched the k estimate */
			if leftright ***REMOVED***
				mhi.Mul(mhi, big10)
			***REMOVED***
			ilim = ilim1
		***REMOVED***
	***REMOVED***
	/* At this point 1 <= d/10^k = b/S < 10. */

	if ilim <= 0 && mode > 2 ***REMOVED***
		/* We're doing fixed-mode output and d is less than the minimum nonzero output in this mode.
		   Output either zero or the minimum nonzero output depending on which is closer to d. */
		if ilim >= 0 ***REMOVED***
			i = b.Cmp(S.Mul(S, big5))
		***REMOVED***
		if ilim < 0 || i < 0 || i == 0 && !biasUp ***REMOVED***
			/* Always emit at least one digit.  If the number appears to be zero
			   using the current mode, then emit one '0' digit and set decpt to 1. */
			buf = buf[:startPos]
			buf = append(buf, '0')
			return buf, 1
		***REMOVED***
		buf = append(buf, '1')
		k++
		return buf, k + 1
	***REMOVED***

	var dig byte
	if leftright ***REMOVED***
		if m2 > 0 ***REMOVED***
			mhi.Lsh(mhi, uint(m2))
		***REMOVED***

		/* Compute mlo -- check for special case
		 * that d is a normalized power of 2.
		 */

		mlo = mhi
		if spec_case ***REMOVED***
			mhi = mlo
			mhi = new(big.Int).Lsh(mhi, log2P)
		***REMOVED***
		/* mlo/S = maximum acceptable error, divided by 10^k, if the output is less than d. */
		/* mhi/S = maximum acceptable error, divided by 10^k, if the output is greater than d. */
		var z, delta big.Int
		for i = 1; ; i++ ***REMOVED***
			z.DivMod(b, S, b)
			dig = byte(z.Int64() + '0')
			/* Do we yet have the shortest decimal string
			 * that will round to d?
			 */
			j = b.Cmp(mlo)
			/* j is b/S compared with mlo/S. */
			delta.Sub(S, mhi)
			var j1 int
			if delta.Sign() <= 0 ***REMOVED***
				j1 = 1
			***REMOVED*** else ***REMOVED***
				j1 = b.Cmp(&delta)
			***REMOVED***
			/* j1 is b/S compared with 1 - mhi/S. */
			if (j1 == 0) && (mode == 0) && ((_word1(d) & 1) == 0) ***REMOVED***
				if dig == '9' ***REMOVED***
					var flag bool
					buf = append(buf, '9')
					if buf, flag = roundOff(buf, startPos); flag ***REMOVED***
						k++
						buf = append(buf, '1')
					***REMOVED***
					return buf, k + 1
				***REMOVED***
				if j > 0 ***REMOVED***
					dig++
				***REMOVED***
				buf = append(buf, dig)
				return buf, k + 1
			***REMOVED***
			if (j < 0) || ((j == 0) && (mode == 0) && ((_word1(d) & 1) == 0)) ***REMOVED***
				if j1 > 0 ***REMOVED***
					/* Either dig or dig+1 would work here as the least significant decimal digit.
					   Use whichever would produce a decimal value closer to d. */
					b.Lsh(b, 1)
					j1 = b.Cmp(S)
					if (j1 > 0) || (j1 == 0 && (((dig & 1) == 1) || biasUp)) ***REMOVED***
						dig++
						if dig == '9' ***REMOVED***
							buf = append(buf, '9')
							buf, flag := roundOff(buf, startPos)
							if flag ***REMOVED***
								k++
								buf = append(buf, '1')
							***REMOVED***
							return buf, k + 1
						***REMOVED***
					***REMOVED***
				***REMOVED***
				buf = append(buf, dig)
				return buf, k + 1
			***REMOVED***
			if j1 > 0 ***REMOVED***
				if dig == '9' ***REMOVED*** /* possible if i == 1 */
					buf = append(buf, '9')
					buf, flag := roundOff(buf, startPos)
					if flag ***REMOVED***
						k++
						buf = append(buf, '1')
					***REMOVED***
					return buf, k + 1
				***REMOVED***
				buf = append(buf, dig+1)
				return buf, k + 1
			***REMOVED***
			buf = append(buf, dig)
			if i == ilim ***REMOVED***
				break
			***REMOVED***
			b.Mul(b, big10)
			if mlo == mhi ***REMOVED***
				mhi.Mul(mhi, big10)
			***REMOVED*** else ***REMOVED***
				mlo.Mul(mlo, big10)
				mhi.Mul(mhi, big10)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		var z big.Int
		for i = 1; ; i++ ***REMOVED***
			z.DivMod(b, S, b)
			dig = byte(z.Int64() + '0')
			buf = append(buf, dig)
			if i >= ilim ***REMOVED***
				break
			***REMOVED***

			b.Mul(b, big10)
		***REMOVED***
	***REMOVED***
	/* Round off last digit */

	b.Lsh(b, 1)
	j = b.Cmp(S)
	if (j > 0) || (j == 0 && (((dig & 1) == 1) || biasUp)) ***REMOVED***
		var flag bool
		buf, flag = roundOff(buf, startPos)
		if flag ***REMOVED***
			k++
			buf = append(buf, '1')
			return buf, k + 1
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		buf = stripTrailingZeroes(buf, startPos)
	***REMOVED***

	return buf, k + 1
***REMOVED***

func bumpUp(buf []byte, k int) ([]byte, int) ***REMOVED***
	var lastCh byte
	stop := 0
	if len(buf) > 0 && buf[0] == '-' ***REMOVED***
		stop = 1
	***REMOVED***
	for ***REMOVED***
		lastCh = buf[len(buf)-1]
		buf = buf[:len(buf)-1]
		if lastCh != '9' ***REMOVED***
			break
		***REMOVED***
		if len(buf) == stop ***REMOVED***
			k++
			lastCh = '0'
			break
		***REMOVED***
	***REMOVED***
	buf = append(buf, lastCh+1)
	return buf, k
***REMOVED***

func setWord0(d float64, w uint32) float64 ***REMOVED***
	dBits := math.Float64bits(d)
	return math.Float64frombits(uint64(w)<<32 | dBits&0xffffffff)
***REMOVED***

func _word0(d float64) uint32 ***REMOVED***
	dBits := math.Float64bits(d)
	return uint32(dBits >> 32)
***REMOVED***

func _word1(d float64) uint32 ***REMOVED***
	dBits := math.Float64bits(d)
	return uint32(dBits)
***REMOVED***

func stripTrailingZeroes(buf []byte, startPos int) []byte ***REMOVED***
	bl := len(buf) - 1
	for bl >= startPos && buf[bl] == '0' ***REMOVED***
		bl--
	***REMOVED***
	return buf[:bl+1]
***REMOVED***

/* Set b = b * 5^k.  k must be nonnegative. */
func pow5mult(b *big.Int, k int) *big.Int ***REMOVED***
	if k < (1 << (len(pow5Cache) + 2)) ***REMOVED***
		i := k & 3
		if i != 0 ***REMOVED***
			b.Mul(b, p05[i-1])
		***REMOVED***
		k >>= 2
		i = 0
		for ***REMOVED***
			if k&1 != 0 ***REMOVED***
				b.Mul(b, pow5Cache[i])
			***REMOVED***
			k >>= 1
			if k == 0 ***REMOVED***
				break
			***REMOVED***
			i++
		***REMOVED***
		return b
	***REMOVED***
	return b.Mul(b, new(big.Int).Exp(big5, big.NewInt(int64(k)), nil))
***REMOVED***

func roundOff(buf []byte, startPos int) ([]byte, bool) ***REMOVED***
	i := len(buf)
	for i != startPos ***REMOVED***
		i--
		if buf[i] != '9' ***REMOVED***
			buf[i]++
			return buf[:i+1], false
		***REMOVED***
	***REMOVED***
	return buf[:startPos], true
***REMOVED***

func init() ***REMOVED***
	p := big.NewInt(625)
	pow5Cache[0] = p
	for i := 1; i < len(pow5Cache); i++ ***REMOVED***
		p = new(big.Int).Mul(p, p)
		pow5Cache[i] = p
	***REMOVED***
***REMOVED***
