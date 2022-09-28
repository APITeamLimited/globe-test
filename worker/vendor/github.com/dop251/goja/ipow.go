package goja

// ported from https://gist.github.com/orlp/3551590

var highest_bit_set = [256]byte***REMOVED***
	0, 1, 2, 2, 3, 3, 3, 3,
	4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 255, // anything past 63 is a guaranteed overflow with base > 1
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255,
***REMOVED***

func ipow(base, exp int64) (result int64) ***REMOVED***
	result = 1

	switch highest_bit_set[byte(exp)] ***REMOVED***
	case 255: // we use 255 as an overflow marker and return 0 on overflow/underflow
		if base == 1 ***REMOVED***
			return 1
		***REMOVED***

		if base == -1 ***REMOVED***
			return 1 - 2*(exp&1)
		***REMOVED***

		return 0
	case 6:
		if exp&1 != 0 ***REMOVED***
			result *= base
		***REMOVED***
		exp >>= 1
		base *= base
		fallthrough
	case 5:
		if exp&1 != 0 ***REMOVED***
			result *= base
		***REMOVED***
		exp >>= 1
		base *= base
		fallthrough
	case 4:
		if exp&1 != 0 ***REMOVED***
			result *= base
		***REMOVED***
		exp >>= 1
		base *= base
		fallthrough
	case 3:
		if exp&1 != 0 ***REMOVED***
			result *= base
		***REMOVED***
		exp >>= 1
		base *= base
		fallthrough
	case 2:
		if exp&1 != 0 ***REMOVED***
			result *= base
		***REMOVED***
		exp >>= 1
		base *= base
		fallthrough
	case 1:
		if exp&1 != 0 ***REMOVED***
			result *= base
		***REMOVED***
		fallthrough
	default:
		return result
	***REMOVED***
***REMOVED***
