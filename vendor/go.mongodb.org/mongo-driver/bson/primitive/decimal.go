// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on gopkg.in/mgo.v2/bson by Gustavo Niemeyer
// See THIRD-PARTY-NOTICES for original license terms.

package primitive

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// These constants are the maximum and minimum values for the exponent field in a decimal128 value.
const (
	MaxDecimal128Exp = 6111
	MinDecimal128Exp = -6176
)

// These errors are returned when an invalid value is parsed as a big.Int.
var (
	ErrParseNaN    = errors.New("cannot parse NaN as a *big.Int")
	ErrParseInf    = errors.New("cannot parse Infinity as a *big.Int")
	ErrParseNegInf = errors.New("cannot parse -Infinity as a *big.Int")
)

// Decimal128 holds decimal128 BSON values.
type Decimal128 struct ***REMOVED***
	h, l uint64
***REMOVED***

// NewDecimal128 creates a Decimal128 using the provide high and low uint64s.
func NewDecimal128(h, l uint64) Decimal128 ***REMOVED***
	return Decimal128***REMOVED***h: h, l: l***REMOVED***
***REMOVED***

// GetBytes returns the underlying bytes of the BSON decimal value as two uint64 values. The first
// contains the most first 8 bytes of the value and the second contains the latter.
func (d Decimal128) GetBytes() (uint64, uint64) ***REMOVED***
	return d.h, d.l
***REMOVED***

// String returns a string representation of the decimal value.
func (d Decimal128) String() string ***REMOVED***
	var posSign int      // positive sign
	var exp int          // exponent
	var high, low uint64 // significand high/low

	if d.h>>63&1 == 0 ***REMOVED***
		posSign = 1
	***REMOVED***

	switch d.h >> 58 & (1<<5 - 1) ***REMOVED***
	case 0x1F:
		return "NaN"
	case 0x1E:
		return "-Infinity"[posSign:]
	***REMOVED***

	low = d.l
	if d.h>>61&3 == 3 ***REMOVED***
		// Bits: 1*sign 2*ignored 14*exponent 111*significand.
		// Implicit 0b100 prefix in significand.
		exp = int(d.h >> 47 & (1<<14 - 1))
		//high = 4<<47 | d.h&(1<<47-1)
		// Spec says all of these values are out of range.
		high, low = 0, 0
	***REMOVED*** else ***REMOVED***
		// Bits: 1*sign 14*exponent 113*significand
		exp = int(d.h >> 49 & (1<<14 - 1))
		high = d.h & (1<<49 - 1)
	***REMOVED***
	exp += MinDecimal128Exp

	// Would be handled by the logic below, but that's trivial and common.
	if high == 0 && low == 0 && exp == 0 ***REMOVED***
		return "-0"[posSign:]
	***REMOVED***

	var repr [48]byte // Loop 5 times over 9 digits plus dot, negative sign, and leading zero.
	var last = len(repr)
	var i = len(repr)
	var dot = len(repr) + exp
	var rem uint32
Loop:
	for d9 := 0; d9 < 5; d9++ ***REMOVED***
		high, low, rem = divmod(high, low, 1e9)
		for d1 := 0; d1 < 9; d1++ ***REMOVED***
			// Handle "-0.0", "0.00123400", "-1.00E-6", "1.050E+3", etc.
			if i < len(repr) && (dot == i || low == 0 && high == 0 && rem > 0 && rem < 10 && (dot < i-6 || exp > 0)) ***REMOVED***
				exp += len(repr) - i
				i--
				repr[i] = '.'
				last = i - 1
				dot = len(repr) // Unmark.
			***REMOVED***
			c := '0' + byte(rem%10)
			rem /= 10
			i--
			repr[i] = c
			// Handle "0E+3", "1E+3", etc.
			if low == 0 && high == 0 && rem == 0 && i == len(repr)-1 && (dot < i-5 || exp > 0) ***REMOVED***
				last = i
				break Loop
			***REMOVED***
			if c != '0' ***REMOVED***
				last = i
			***REMOVED***
			// Break early. Works without it, but why.
			if dot > i && low == 0 && high == 0 && rem == 0 ***REMOVED***
				break Loop
			***REMOVED***
		***REMOVED***
	***REMOVED***
	repr[last-1] = '-'
	last--

	if exp > 0 ***REMOVED***
		return string(repr[last+posSign:]) + "E+" + strconv.Itoa(exp)
	***REMOVED***
	if exp < 0 ***REMOVED***
		return string(repr[last+posSign:]) + "E" + strconv.Itoa(exp)
	***REMOVED***
	return string(repr[last+posSign:])
***REMOVED***

// BigInt returns significand as big.Int and exponent, bi * 10 ^ exp.
func (d Decimal128) BigInt() (*big.Int, int, error) ***REMOVED***
	high, low := d.GetBytes()
	posSign := high>>63&1 == 0 // positive sign

	switch high >> 58 & (1<<5 - 1) ***REMOVED***
	case 0x1F:
		return nil, 0, ErrParseNaN
	case 0x1E:
		if posSign ***REMOVED***
			return nil, 0, ErrParseInf
		***REMOVED***
		return nil, 0, ErrParseNegInf
	***REMOVED***

	var exp int
	if high>>61&3 == 3 ***REMOVED***
		// Bits: 1*sign 2*ignored 14*exponent 111*significand.
		// Implicit 0b100 prefix in significand.
		exp = int(high >> 47 & (1<<14 - 1))
		//high = 4<<47 | d.h&(1<<47-1)
		// Spec says all of these values are out of range.
		high, low = 0, 0
	***REMOVED*** else ***REMOVED***
		// Bits: 1*sign 14*exponent 113*significand
		exp = int(high >> 49 & (1<<14 - 1))
		high = high & (1<<49 - 1)
	***REMOVED***
	exp += MinDecimal128Exp

	// Would be handled by the logic below, but that's trivial and common.
	if high == 0 && low == 0 && exp == 0 ***REMOVED***
		if posSign ***REMOVED***
			return new(big.Int), 0, nil
		***REMOVED***
		return new(big.Int), 0, nil
	***REMOVED***

	bi := big.NewInt(0)
	const host32bit = ^uint(0)>>32 == 0
	if host32bit ***REMOVED***
		bi.SetBits([]big.Word***REMOVED***big.Word(low), big.Word(low >> 32), big.Word(high), big.Word(high >> 32)***REMOVED***)
	***REMOVED*** else ***REMOVED***
		bi.SetBits([]big.Word***REMOVED***big.Word(low), big.Word(high)***REMOVED***)
	***REMOVED***

	if !posSign ***REMOVED***
		return bi.Neg(bi), exp, nil
	***REMOVED***
	return bi, exp, nil
***REMOVED***

// IsNaN returns whether d is NaN.
func (d Decimal128) IsNaN() bool ***REMOVED***
	return d.h>>58&(1<<5-1) == 0x1F
***REMOVED***

// IsInf returns:
//
//   +1 d == Infinity
//    0 other case
//   -1 d == -Infinity
//
func (d Decimal128) IsInf() int ***REMOVED***
	if d.h>>58&(1<<5-1) != 0x1E ***REMOVED***
		return 0
	***REMOVED***

	if d.h>>63&1 == 0 ***REMOVED***
		return 1
	***REMOVED***
	return -1
***REMOVED***

// IsZero returns true if d is the empty Decimal128.
func (d Decimal128) IsZero() bool ***REMOVED***
	return d.h == 0 && d.l == 0
***REMOVED***

// MarshalJSON returns Decimal128 as a string.
func (d Decimal128) MarshalJSON() ([]byte, error) ***REMOVED***
	return json.Marshal(d.String())
***REMOVED***

// UnmarshalJSON creates a primitive.Decimal128 from a JSON string, an extended JSON $numberDecimal value, or the string
// "null". If b is a JSON string or extended JSON value, d will have the value of that string, and if b is "null", d will
// be unchanged.
func (d *Decimal128) UnmarshalJSON(b []byte) error ***REMOVED***
	// Ignore "null" to keep parity with the standard library. Decoding a JSON null into a non-pointer Decimal128 field
	// will leave the field unchanged. For pointer values, encoding/json will set the pointer to nil and will not
	// enter the UnmarshalJSON hook.
	if string(b) == "null" ***REMOVED***
		return nil
	***REMOVED***

	var res interface***REMOVED******REMOVED***
	err := json.Unmarshal(b, &res)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	str, ok := res.(string)

	// Extended JSON
	if !ok ***REMOVED***
		m, ok := res.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return errors.New("not an extended JSON Decimal128: expected document")
		***REMOVED***
		d128, ok := m["$numberDecimal"]
		if !ok ***REMOVED***
			return errors.New("not an extended JSON Decimal128: expected key $numberDecimal")
		***REMOVED***
		str, ok = d128.(string)
		if !ok ***REMOVED***
			return errors.New("not an extended JSON Decimal128: expected decimal to be string")
		***REMOVED***
	***REMOVED***

	*d, err = ParseDecimal128(str)
	return err
***REMOVED***

func divmod(h, l uint64, div uint32) (qh, ql uint64, rem uint32) ***REMOVED***
	div64 := uint64(div)
	a := h >> 32
	aq := a / div64
	ar := a % div64
	b := ar<<32 + h&(1<<32-1)
	bq := b / div64
	br := b % div64
	c := br<<32 + l>>32
	cq := c / div64
	cr := c % div64
	d := cr<<32 + l&(1<<32-1)
	dq := d / div64
	dr := d % div64
	return (aq<<32 | bq), (cq<<32 | dq), uint32(dr)
***REMOVED***

var dNaN = Decimal128***REMOVED***0x1F << 58, 0***REMOVED***
var dPosInf = Decimal128***REMOVED***0x1E << 58, 0***REMOVED***
var dNegInf = Decimal128***REMOVED***0x3E << 58, 0***REMOVED***

func dErr(s string) (Decimal128, error) ***REMOVED***
	return dNaN, fmt.Errorf("cannot parse %q as a decimal128", s)
***REMOVED***

// match scientific notation number, example -10.15e-18
var normalNumber = regexp.MustCompile(`^(?P<int>[-+]?\d*)?(?:\.(?P<dec>\d*))?(?:[Ee](?P<exp>[-+]?\d+))?$`)

// ParseDecimal128 takes the given string and attempts to parse it into a valid
// Decimal128 value.
func ParseDecimal128(s string) (Decimal128, error) ***REMOVED***
	if s == "" ***REMOVED***
		return dErr(s)
	***REMOVED***

	matches := normalNumber.FindStringSubmatch(s)
	if len(matches) == 0 ***REMOVED***
		orig := s
		neg := s[0] == '-'
		if neg || s[0] == '+' ***REMOVED***
			s = s[1:]
		***REMOVED***

		if s == "NaN" || s == "nan" || strings.EqualFold(s, "nan") ***REMOVED***
			return dNaN, nil
		***REMOVED***
		if s == "Inf" || s == "inf" || strings.EqualFold(s, "inf") || strings.EqualFold(s, "infinity") ***REMOVED***
			if neg ***REMOVED***
				return dNegInf, nil
			***REMOVED***
			return dPosInf, nil
		***REMOVED***
		return dErr(orig)
	***REMOVED***

	intPart := matches[1]
	decPart := matches[2]
	expPart := matches[3]

	var err error
	exp := 0
	if expPart != "" ***REMOVED***
		exp, err = strconv.Atoi(expPart)
		if err != nil ***REMOVED***
			return dErr(s)
		***REMOVED***
	***REMOVED***
	if decPart != "" ***REMOVED***
		exp -= len(decPart)
	***REMOVED***

	if len(strings.Trim(intPart+decPart, "-0")) > 35 ***REMOVED***
		return dErr(s)
	***REMOVED***

	bi, ok := new(big.Int).SetString(intPart+decPart, 10)
	if !ok ***REMOVED***
		return dErr(s)
	***REMOVED***

	d, ok := ParseDecimal128FromBigInt(bi, exp)
	if !ok ***REMOVED***
		return dErr(s)
	***REMOVED***

	if bi.Sign() == 0 && s[0] == '-' ***REMOVED***
		d.h |= 1 << 63
	***REMOVED***

	return d, nil
***REMOVED***

var (
	ten  = big.NewInt(10)
	zero = new(big.Int)

	maxS, _ = new(big.Int).SetString("9999999999999999999999999999999999", 10)
)

// ParseDecimal128FromBigInt attempts to parse the given significand and exponent into a valid Decimal128 value.
func ParseDecimal128FromBigInt(bi *big.Int, exp int) (Decimal128, bool) ***REMOVED***
	//copy
	bi = new(big.Int).Set(bi)

	q := new(big.Int)
	r := new(big.Int)

	for bigIntCmpAbs(bi, maxS) == 1 ***REMOVED***
		bi, _ = q.QuoRem(bi, ten, r)
		if r.Cmp(zero) != 0 ***REMOVED***
			return Decimal128***REMOVED******REMOVED***, false
		***REMOVED***
		exp++
		if exp > MaxDecimal128Exp ***REMOVED***
			return Decimal128***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED***

	for exp < MinDecimal128Exp ***REMOVED***
		// Subnormal.
		bi, _ = q.QuoRem(bi, ten, r)
		if r.Cmp(zero) != 0 ***REMOVED***
			return Decimal128***REMOVED******REMOVED***, false
		***REMOVED***
		exp++
	***REMOVED***
	for exp > MaxDecimal128Exp ***REMOVED***
		// Clamped.
		bi.Mul(bi, ten)
		if bigIntCmpAbs(bi, maxS) == 1 ***REMOVED***
			return Decimal128***REMOVED******REMOVED***, false
		***REMOVED***
		exp--
	***REMOVED***

	b := bi.Bytes()
	var h, l uint64
	for i := 0; i < len(b); i++ ***REMOVED***
		if i < len(b)-8 ***REMOVED***
			h = h<<8 | uint64(b[i])
			continue
		***REMOVED***
		l = l<<8 | uint64(b[i])
	***REMOVED***

	h |= uint64(exp-MinDecimal128Exp) & uint64(1<<14-1) << 49
	if bi.Sign() == -1 ***REMOVED***
		h |= 1 << 63
	***REMOVED***

	return Decimal128***REMOVED***h: h, l: l***REMOVED***, true
***REMOVED***

// bigIntCmpAbs computes big.Int.Cmp(absoluteValue(x), absoluteValue(y)).
func bigIntCmpAbs(x, y *big.Int) int ***REMOVED***
	xAbs := bigIntAbsValue(x)
	yAbs := bigIntAbsValue(y)
	return xAbs.Cmp(yAbs)
***REMOVED***

// bigIntAbsValue returns a big.Int containing the absolute value of b.
// If b is already a non-negative number, it is returned without any changes or copies.
func bigIntAbsValue(b *big.Int) *big.Int ***REMOVED***
	if b.Sign() >= 0 ***REMOVED***
		return b // already positive
	***REMOVED***
	return new(big.Int).Abs(b)
***REMOVED***
