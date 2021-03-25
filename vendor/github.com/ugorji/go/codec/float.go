// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "strconv"

// func parseFloat(b []byte, bitsize int) (f float64, err error) ***REMOVED***
// 	if bitsize == 32 ***REMOVED***
// 		return parseFloat32(b)
// 	***REMOVED*** else ***REMOVED***
// 		return parseFloat64(b)
// 	***REMOVED***
// ***REMOVED***

func parseFloat32(b []byte) (f float32, err error) ***REMOVED***
	return parseFloat32_custom(b)
	// return parseFloat32_strconv(b)
***REMOVED***

func parseFloat64(b []byte) (f float64, err error) ***REMOVED***
	return parseFloat64_custom(b)
	// return parseFloat64_strconv(b)
***REMOVED***

func parseFloat32_strconv(b []byte) (f float32, err error) ***REMOVED***
	f64, err := strconv.ParseFloat(stringView(b), 32)
	f = float32(f64)
	return
***REMOVED***

func parseFloat64_strconv(b []byte) (f float64, err error) ***REMOVED***
	return strconv.ParseFloat(stringView(b), 64)
***REMOVED***

// ------ parseFloat custom below --------

// We assume that a lot of floating point numbers in json files will be
// those that are handwritten, and with defined precision (in terms of number
// of digits after decimal point), etc.
//
// We further assume that this ones can be written in exact format.
//
// strconv.ParseFloat has some unnecessary overhead which we can do without
// for the common case:
//
//    - expensive char-by-char check to see if underscores are in right place
//    - testing for and skipping underscores
//    - check if the string matches ignorecase +/- inf, +/- infinity, nan
//    - support for base 16 (0xFFFF...)
//
// The functions below will try a fast-path for floats which can be decoded
// without any loss of precision, meaning they:
//
//    - fits within the significand bits of the 32-bits or 64-bits
//    - exponent fits within the exponent value
//    - there is no truncation (any extra numbers are all trailing zeros)
//
// To figure out what the values are for maxMantDigits, use this idea below:
//
// 2^23 =                 838 8608 (between 10^ 6 and 10^ 7) (significand bits of uint32)
// 2^32 =             42 9496 7296 (between 10^ 9 and 10^10) (full uint32)
// 2^52 =      4503 5996 2737 0496 (between 10^15 and 10^16) (significand bits of uint64)
// 2^64 = 1844 6744 0737 0955 1616 (between 10^19 and 10^20) (full uint64)
//
// Since we only allow for up to what can comfortably fit into the significand
// ignoring the exponent, and we only try to parse iff significand fits into the

// Exact powers of 10.
var float64pow10 = [...]float64***REMOVED***
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
***REMOVED***
var float32pow10 = [...]float32***REMOVED***1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10***REMOVED***

type floatinfo struct ***REMOVED***
	mantbits uint8
	expbits  uint8
	bias     int16

	exactPow10 int8 // Exact powers of ten are <= 10^N (32: 10, 64: 22)
	exactInts  int8 // Exact integers are <= 10^N

	maxMantDigits int8 // 10^19 fits in uint64, while 10^9 fits in uint32
***REMOVED***

var fi32 = floatinfo***REMOVED***23, 8, -127, 10, 7, 9***REMOVED***     // maxMantDigits = 9
var fi64 = floatinfo***REMOVED***52, 11, -1023, 22, 15, 19***REMOVED*** // maxMantDigits = 19

const fMax64 = 1e15
const fMax32 = 1e7

const fBase = 10

func parseFloatErr(b []byte) error ***REMOVED***
	return &strconv.NumError***REMOVED***
		Func: "ParseFloat",
		Err:  strconv.ErrSyntax,
		Num:  string(b),
	***REMOVED***
***REMOVED***

func parseFloat32_custom(b []byte) (f float32, err error) ***REMOVED***
	mantissa, exp, neg, trunc, bad, ok := readFloat(b, fi32)
	_ = trunc
	if bad ***REMOVED***
		return 0, parseFloatErr(b)
	***REMOVED***
	if ok ***REMOVED***
		// parseFloatDebug(b, 32, false, exp, trunc, ok)
		f = float32(mantissa)
		if neg ***REMOVED***
			f = -f
		***REMOVED***
		if exp != 0 ***REMOVED***
			indx := fExpIndx(exp)
			if exp < 0 ***REMOVED*** // int / 10^k
				f /= float32pow10[indx]
			***REMOVED*** else ***REMOVED*** // exp > 0
				if exp > fi32.exactPow10 ***REMOVED***
					f *= float32pow10[exp-fi32.exactPow10]
					if f < -fMax32 || f > fMax32 ***REMOVED*** // exponent too large - outside range
						goto FALLBACK
					***REMOVED***
					indx = uint8(fi32.exactPow10)
				***REMOVED***
				f *= float32pow10[indx]
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***
FALLBACK:
	// parseFloatDebug(b, 32, true, exp, trunc, ok)
	return parseFloat32_strconv(b)
***REMOVED***

func parseFloat64_custom(b []byte) (f float64, err error) ***REMOVED***
	mantissa, exp, neg, trunc, bad, ok := readFloat(b, fi64)
	_ = trunc
	if bad ***REMOVED***
		return 0, parseFloatErr(b)
	***REMOVED***
	if ok ***REMOVED***
		f = float64(mantissa)
		if neg ***REMOVED***
			f = -f
		***REMOVED***
		if exp != 0 ***REMOVED***
			indx := fExpIndx(exp)
			if exp < 0 ***REMOVED*** // int / 10^k
				f /= float64pow10[indx]
			***REMOVED*** else ***REMOVED*** // exp > 0
				if exp > fi64.exactPow10 ***REMOVED***
					f *= float64pow10[exp-fi64.exactPow10]
					if f < -fMax64 || f > fMax64 ***REMOVED*** // exponent too large - outside range
						goto FALLBACK
					***REMOVED***
					indx = uint8(fi64.exactPow10)
				***REMOVED***
				f *= float64pow10[indx]
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***
FALLBACK:
	return parseFloat64_strconv(b)
***REMOVED***

func fExpIndx(v int8) uint8 ***REMOVED***
	if v < 0 ***REMOVED***
		return uint8(-v)
	***REMOVED***
	return uint8(v)
***REMOVED***

func readFloat(s []byte, y floatinfo) (mantissa uint64, exp int8, neg, trunc, bad, ok bool) ***REMOVED***
	var i uint // make it uint, so that we eliminate bounds checking
	var slen = uint(len(s))
	if slen == 0 ***REMOVED***
		bad = true
		return
	***REMOVED***
	switch s[0] ***REMOVED***
	case '+':
		i++
	case '-':
		neg = true
		i++
	***REMOVED***

	// we considered punting early if string has length > maxMantDigits, but this doesn't account
	// for trailing 0's e.g. 700000000000000000000 can be encoded exactly as it is 7e20

	// var sawdot, sawdigits, sawexp bool
	var sawdot, sawexp bool
	var nd, ndMant, dp int8
L:
	for ; i < slen; i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '.':
			if sawdot ***REMOVED***
				bad = true
				return
			***REMOVED***
			sawdot = true
			dp = nd
		case '0':
			if nd == 0 ***REMOVED*** // ignore leading zeros
				dp--
				continue
			***REMOVED***
			nd++
			if ndMant < y.maxMantDigits ***REMOVED***
				// mantissa = (mantissa << 1) + (mantissa << 3)
				mantissa *= fBase
				ndMant++
			***REMOVED***
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// sawdigits = true
			nd++
			if ndMant < y.maxMantDigits ***REMOVED***
				// mantissa = (mantissa << 1) + (mantissa << 3) + uint64(s[i]-'0')
				mantissa = mantissa*fBase + uint64(s[i]-'0')
				// mantissa *= fBase
				// mantissa += uint64(s[i] - '0')
				ndMant++
			***REMOVED*** else ***REMOVED***
				trunc = true
				return // break L
			***REMOVED***
		case 'e', 'E':
			sawexp = true
			break L
		default:
			bad = true
			return
		***REMOVED***
	***REMOVED***
	// if !sawdigits ***REMOVED***
	// 	bad = true
	// 	return
	// ***REMOVED***
	if !sawdot ***REMOVED***
		dp = nd
	***REMOVED***

	if sawexp ***REMOVED***
		i++
		if i < slen ***REMOVED***
			var eneg bool
			if s[i] == '+' ***REMOVED***
				i++
			***REMOVED*** else if s[i] == '-' ***REMOVED***
				i++
				eneg = true
			***REMOVED***
			if i < slen ***REMOVED***
				// for exact match, exponent is 1 or 2 digits (float64: -22 to 37, float32: -1 to 17).
				// exit quick if exponent is more than 2 digits.
				if i+2 < slen ***REMOVED***
					return
				***REMOVED***

				var e int8

				if s[i] < '0' || s[i] > '9' ***REMOVED***
					bad = true
					return
				***REMOVED***
				e = e*fBase + int8(s[i]-'0') // (e << 1) + (e << 3) + int8(s[i]-'0')
				i++

				if i < slen ***REMOVED***
					if s[i] < '0' || s[i] > '9' ***REMOVED***
						bad = true
						return
					***REMOVED***
					e = e*fBase + int8(s[i]-'0') // (e << 1) + (e << 3) + int8(s[i]-'0')
					i++
				***REMOVED***

				if eneg ***REMOVED***
					dp -= e
				***REMOVED*** else ***REMOVED***
					dp += e
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if mantissa != 0 ***REMOVED***
		if mantissa>>y.mantbits != 0 ***REMOVED***
			return
		***REMOVED***
		exp = dp - ndMant
		if exp < -y.exactPow10 || exp > y.exactInts+y.exactPow10 ***REMOVED*** // cannot handle it
			return
		***REMOVED***
	***REMOVED***
	ok = true // && !trunc // if trunc=true, we return early (so here trunc=false)
	return
***REMOVED***

// fMul10ShiftU64

// func parseFloatDebug(b []byte, bitsize int, strconv bool, exp int8, trunc, ok bool) ***REMOVED***
// 	if strconv ***REMOVED***
// 		xdebugf("parseFloat%d: delegating: %s, exp: %d, trunc: %v, ok: %v", bitsize, b, exp, trunc, ok)
// 	***REMOVED*** else ***REMOVED***
// 		xdebug2f("parseFloat%d: attempting: %s, exp: %d, trunc: %v, ok: %v", bitsize, b, exp, trunc, ok)
// 	***REMOVED***
// ***REMOVED***
