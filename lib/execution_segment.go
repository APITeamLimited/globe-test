/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

import (
	"encoding"
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// ExecutionSegment represents a (start, end] partition of the total execution
// work for a specific test. For example, if we want the split the execution of a
// test in 2 different parts, we can split it in two segments (0, 0.5] and (0,5, 1].
//
// We use rational numbers so it's easier to verify the correctness and easier to
// reason about portions of indivisible things, like VUs. This way, we can easily
// split a test in thirds (i.e. (0, 1/3], (1/3, 2/3], (2/3, 1]), without fearing
// that we'll lose a VU along the way...
//
// The most important part is that if work is split between multiple k6 instances,
// each k6 instance can precisely and reproducibly calculate its share of the work,
// just by knowing its own segment. There won't be a need to schedule the
// execution from a master node, or to even know how many other k6 instances are
// running!
type ExecutionSegment struct ***REMOVED***
	// 0 <= from < to <= 1
	from *big.Rat
	to   *big.Rat

	// derived, equals to-from, but pre-calculated here for speed
	length *big.Rat
***REMOVED***

// Ensure we implement those interfaces
var (
	_ encoding.TextUnmarshaler = &ExecutionSegment***REMOVED******REMOVED***
	_ fmt.Stringer             = &ExecutionSegment***REMOVED******REMOVED***
)

// Helpful "constants" so we don't initialize them in every function call
var (
	zeroRat, oneRat      = big.NewRat(0, 1), big.NewRat(1, 1) //nolint:gochecknoglobals
	oneBigInt, twoBigInt = big.NewInt(1), big.NewInt(2)       //nolint:gochecknoglobals
)

// NewExecutionSegment validates the supplied arguments (basically, that 0 <=
// from < to <= 1) and either returns an error, or it returns a
// fully-initialized and usable execution segment.
func NewExecutionSegment(from, to *big.Rat) (*ExecutionSegment, error) ***REMOVED***
	if from.Cmp(zeroRat) < 0 ***REMOVED***
		return nil, fmt.Errorf("segment start value should be at least 0 but was %s", from.FloatString(2))
	***REMOVED***
	if from.Cmp(to) >= 0 ***REMOVED***
		return nil, fmt.Errorf("segment start(%s) should be less than its end(%s)", from.FloatString(2), to.FloatString(2))
	***REMOVED***
	if to.Cmp(oneRat) > 0 ***REMOVED***
		return nil, fmt.Errorf("segment end value shouldn't be more than 1 but was %s", to.FloatString(2))
	***REMOVED***
	return &ExecutionSegment***REMOVED***
		from:   from,
		to:     to,
		length: new(big.Rat).Sub(to, from),
	***REMOVED***, nil
***REMOVED***

// stringToRat is a helper function that tries to convert a string to a rational
// number while allowing percentage, decimal, and fraction values.
func stringToRat(s string) (*big.Rat, error) ***REMOVED***
	if strings.HasSuffix(s, "%") ***REMOVED***
		num, ok := new(big.Int).SetString(strings.TrimSuffix(s, "%"), 10)
		if !ok ***REMOVED***
			return nil, fmt.Errorf("'%s' is not a valid percentage", s)
		***REMOVED***
		return new(big.Rat).SetFrac(num, big.NewInt(100)), nil
	***REMOVED***
	rat, ok := new(big.Rat).SetString(s)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("'%s' is not a valid percentage, decimal, fraction or interval value", s)
	***REMOVED***
	return rat, nil
***REMOVED***

// NewExecutionSegmentFromString validates the supplied string value and returns
// the newly created ExecutionSegment or and error from it.
//
// We are able to parse both single percentage/float/fraction values, and actual
// (from: to] segments. For the single values, we just treat them as the
// beginning segment - thus the execution segment can be used as a shortcut for
// quickly running an arbitrarily scaled-down version of a test.
//
// The parsing logic is that values with a colon, i.e. ':', are full segments:
//  `1/2:3/4`, `0.5:0.75`, `50%:75%`, and even `2/4:75%` should be (1/2, 3/4]
// And values without a colon are the end of a first segment:
//  `20%`, `0.2`,  and `1/5` should be converted to (0, 1/5]
// empty values should probably be treated as "1", i.e. the whole execution
func NewExecutionSegmentFromString(toStr string) (result *ExecutionSegment, err error) ***REMOVED***
	from := zeroRat
	if toStr == "" ***REMOVED***
		toStr = "1" // an empty string means a full 0:1 execution segment
	***REMOVED***
	if strings.ContainsRune(toStr, ':') ***REMOVED***
		fromToStr := strings.SplitN(toStr, ":", 2)
		toStr = fromToStr[1]
		if from, err = stringToRat(fromToStr[0]); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	to, err := stringToRat(toStr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewExecutionSegment(from, to)
***REMOVED***

// UnmarshalText implements the encoding.TextUnmarshaler interface, so that
// execution segments can be specified as CLI flags, environment variables, and
// JSON strings. It is a wrapper for the NewExecutionFromString() constructor.
func (es *ExecutionSegment) UnmarshalText(text []byte) (err error) ***REMOVED***
	segment, err := NewExecutionSegmentFromString(string(text))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*es = *segment
	return nil
***REMOVED***

func (es *ExecutionSegment) String() string ***REMOVED***
	if es == nil ***REMOVED***
		return "0:1"
	***REMOVED***
	return es.from.RatString() + ":" + es.to.RatString()
***REMOVED***

// MarshalText implements the encoding.TextMarshaler interface, so is used for
// text and JSON encoding of the execution segment.
func (es *ExecutionSegment) MarshalText() ([]byte, error) ***REMOVED***
	if es == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	return []byte(es.String()), nil
***REMOVED***

// FloatLength is a helper method for getting some more human-readable
// information about the execution segment.
func (es *ExecutionSegment) FloatLength() float64 ***REMOVED***
	if es == nil ***REMOVED***
		return 1.0
	***REMOVED***
	res, _ := es.length.Float64()
	return res
***REMOVED***

// Split evenly divides the execution segment into the specified number of
// equal consecutive execution sub-segments.
func (es *ExecutionSegment) Split(numParts int64) ([]*ExecutionSegment, error) ***REMOVED***
	if numParts < 1 ***REMOVED***
		return nil, fmt.Errorf("the number of parts should be at least 1, %d received", numParts)
	***REMOVED***

	from, to := zeroRat, oneRat
	if es != nil ***REMOVED***
		from, to = es.from, es.to
	***REMOVED***

	increment := new(big.Rat).Sub(to, from)
	increment.Denom().Mul(increment.Denom(), big.NewInt(numParts))

	results := make([]*ExecutionSegment, numParts)
	for i := int64(0); i < numParts; i++ ***REMOVED***
		segmentTo := new(big.Rat).Add(from, increment)
		segment, err := NewExecutionSegment(from, segmentTo)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		results[i] = segment
		from = segmentTo
	***REMOVED***

	if from.Cmp(to) != 0 ***REMOVED***
		return nil, fmt.Errorf("expected %s and %s to be equal", from, to)
	***REMOVED***

	return results, nil
***REMOVED***

// Equal returns true only if the two execution segments have the same from and
// to values.
func (es *ExecutionSegment) Equal(other *ExecutionSegment) bool ***REMOVED***
	if es == other ***REMOVED***
		return true
	***REMOVED***
	thisFrom, otherFrom, thisTo, otherTo := zeroRat, zeroRat, oneRat, oneRat
	if es != nil ***REMOVED***
		thisFrom, thisTo = es.from, es.to
	***REMOVED***
	if other != nil ***REMOVED***
		otherFrom, otherTo = other.from, other.to
	***REMOVED***
	return thisFrom.Cmp(otherFrom) == 0 && thisTo.Cmp(otherTo) == 0
***REMOVED***

// SubSegment returns a new execution sub-segment - if a is (1/2:1] and b is
// (0:1/2], then a.SubSegment(b) will return a new segment (1/2, 3/4].
//
// The basic formula for c = a.SubSegment(b) is:
//    c.from = a.from + b.from * (a.to - a.from)
//    c.to = c.from + (b.to - b.from) * (a.to - a.from)
func (es *ExecutionSegment) SubSegment(child *ExecutionSegment) *ExecutionSegment ***REMOVED***
	if child == nil ***REMOVED***
		return es // 100% sub-segment is the original segment
	***REMOVED***

	parentFrom, parentLength := zeroRat, oneRat
	if es != nil ***REMOVED***
		parentFrom, parentLength = es.from, es.length
	***REMOVED***

	resultFrom := new(big.Rat).Mul(parentLength, child.from)
	resultFrom.Add(resultFrom, parentFrom)

	resultLength := new(big.Rat).Mul(parentLength, child.length)
	return &ExecutionSegment***REMOVED***
		from:   resultFrom,
		length: resultLength,
		to:     new(big.Rat).Add(resultFrom, resultLength),
	***REMOVED***
***REMOVED***

// helper function for rounding (up) of rational numbers to big.Int values
func roundUp(rat *big.Rat) *big.Int ***REMOVED***
	quo, rem := new(big.Int).QuoRem(rat.Num(), rat.Denom(), new(big.Int))

	if rem.Mul(rem, twoBigInt).Cmp(rat.Denom()) >= 0 ***REMOVED***
		return quo.Add(quo, oneBigInt)
	***REMOVED***
	return quo
***REMOVED***

// Scale proportionally scales the supplied value, according to the execution
// segment's position and size of the work.
func (es *ExecutionSegment) Scale(value int64) int64 ***REMOVED***
	if es == nil ***REMOVED*** // no execution segment, i.e. 100%
		return value
	***REMOVED***
	// Instead of the first proposal that used remainders and floor:
	//    floor( (value * from) % 1 + value * length )
	// We're using an alternative approach with rounding that (hopefully) has
	// the same properties, but it's simpler and has better precision:
	//    round( (value * from) - round(value * from) + (value * (to - from)) )?
	// which reduces to:
	//    round( (value * to) - round(value * from) )?

	toValue := big.NewRat(value, 1)
	toValue.Mul(toValue, es.to)

	fromValue := big.NewRat(value, 1)
	fromValue.Mul(fromValue, es.from)

	toValue.Sub(toValue, new(big.Rat).SetFrac(roundUp(fromValue), oneBigInt))

	return roundUp(toValue).Int64()
***REMOVED***

// InPlaceScaleRat scales rational numbers in-place - it changes the passed
// argument (and also returns it, to allow for chaining, like many other big.Rat
// methods).
func (es *ExecutionSegment) InPlaceScaleRat(value *big.Rat) *big.Rat ***REMOVED***
	if es == nil ***REMOVED*** // no execution segment, i.e. 100%
		return value
	***REMOVED***
	return value.Mul(value, es.length)
***REMOVED***

// CopyScaleRat scales rational numbers without changing them - creates a new
// bit.Rat object and uses it for the calculation.
func (es *ExecutionSegment) CopyScaleRat(value *big.Rat) *big.Rat ***REMOVED***
	if es == nil ***REMOVED*** // no execution segment, i.e. 100%
		return value
	***REMOVED***
	return new(big.Rat).Mul(value, es.length)
***REMOVED***

// ExecutionSegmentSequence represents an ordered chain of execution segments,
// where the end of one segment is the beginning of the next. It can serialized
// as a comma-separated string of rational numbers "r1,r2,r3,...,rn", which
// represents the sequence (r1, r2], (r2, r3], (r3, r4], ..., (r***REMOVED***n-1***REMOVED***, rn].
// The empty value should be treated as if there is a single (0, 1] segment.
type ExecutionSegmentSequence []*ExecutionSegment

// NewExecutionSegmentSequence validates the that the supplied execution
// segments are non-overlapping and without gaps. It will return a new execution
// segment sequence if that is true, and an error if it's not.
func NewExecutionSegmentSequence(segments ...*ExecutionSegment) (ExecutionSegmentSequence, error) ***REMOVED***
	if len(segments) > 1 ***REMOVED***
		to := segments[0].to
		for i, segment := range segments[1:] ***REMOVED***
			if segment.from.Cmp(to) != 0 ***REMOVED***
				return nil, fmt.Errorf(
					"the start value %s of segment #%d should be equal to the end value of the previous one, but it is %s",
					segment.from, i+1, to,
				)
			***REMOVED***
			to = segment.to
		***REMOVED***
	***REMOVED***
	return ExecutionSegmentSequence(segments), nil
***REMOVED***

// NewExecutionSegmentSequenceFromString parses strings of the format
// "r1,r2,r3,...,rn", which represents the sequences like (r1, r2], (r2, r3],
// (r3, r4], ..., (r***REMOVED***n-1***REMOVED***, rn].
func NewExecutionSegmentSequenceFromString(strSeq string) (ExecutionSegmentSequence, error) ***REMOVED***
	if len(strSeq) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	points := strings.Split(strSeq, ",")
	if len(points) < 2 ***REMOVED***
		return nil, fmt.Errorf("at least 2 points are needed for an execution segment sequence, %d given", len(points))
	***REMOVED***
	var start *big.Rat

	segments := make([]*ExecutionSegment, 0, len(points)-1)
	for i, point := range points ***REMOVED***
		rat, err := stringToRat(point)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if i == 0 ***REMOVED***
			start = rat
			continue
		***REMOVED***

		segment, err := NewExecutionSegment(start, rat)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		segments = append(segments, segment)
		start = rat
	***REMOVED***

	return NewExecutionSegmentSequence(segments...)
***REMOVED***

// UnmarshalText implements the encoding.TextUnmarshaler interface, so that
// execution segment sequences can be specified as CLI flags, environment
// variables, and JSON strings.
func (ess *ExecutionSegmentSequence) UnmarshalText(text []byte) (err error) ***REMOVED***
	seq, err := NewExecutionSegmentSequenceFromString(string(text))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*ess = seq
	return nil
***REMOVED***

// MarshalText implements the encoding.TextMarshaler interface, so is used for
// text and JSON encoding of the execution segment sequences.
func (ess ExecutionSegmentSequence) MarshalText() ([]byte, error) ***REMOVED***
	return []byte(ess.String()), nil
***REMOVED***

// String just implements the fmt.Stringer interface, encoding the sequence of
// segments as "start1,end1,end2,end3,...,endn".
func (ess ExecutionSegmentSequence) String() string ***REMOVED***
	result := make([]string, 0, len(ess)+1)
	for i, s := range ess ***REMOVED***
		if i == 0 ***REMOVED***
			result = append(result, s.from.RatString())
		***REMOVED***
		result = append(result, s.to.RatString())
	***REMOVED***
	return strings.Join(result, ",")
***REMOVED***

// lowest common denominator based on https://rosettacode.org/wiki/Least_common_multiple#Go
func (ess ExecutionSegmentSequence) lcd() int64 ***REMOVED***
	var m, n, z big.Int
	z = *ess[0].length.Denom()
	for _, seg := range ess[1:] ***REMOVED***
		m = z
		n = *seg.length.Denom()
		z.Mul(z.Div(&m, z.GCD(nil, nil, &m, &n)), &n)
	***REMOVED***

	return z.Int64()
***REMOVED***

// GetStripedOffsets returns everything that you need in order to execute only
// the iterations that belong to the supplied segment...
//
// TODO: add a more detailed algorithm description
// TODO: basically https://docs.google.com/spreadsheets/d/1V_ivN2xuaMJIgOf1HkpOw1ex8QOhxp960itGGiRrNzo/edit
func (ess ExecutionSegmentSequence) GetStripedOffsets(segment *ExecutionSegment) (int64, []int64, int64, error) ***REMOVED***
	// TODO check if ExecutionSegmentSequence actually starts from 0 and goes to 1 ?
	var matchingSegment *ExecutionSegment
	for _, seg := range ess ***REMOVED***
		if seg.Equal(segment) ***REMOVED***
			matchingSegment = seg
			break
		***REMOVED***
	***REMOVED***
	if matchingSegment == nil ***REMOVED***
		return -1, nil, -1, fmt.Errorf("missing segment %s inside segment sequence %s", segment, ess)
	***REMOVED***

	ess = append([]*ExecutionSegment***REMOVED******REMOVED***, ess...)
	sort.SliceStable(ess,
		func(i, j int) bool ***REMOVED***
			// Yes this Less is actually More, but we want it sorted in descending order and the alternative is to reverse it after sort
			return ess[i].length.Cmp(ess[j].length) > 0
		***REMOVED***)

	var numerators = make([]int64, len(ess))
	var numeratorChanges = make([]*big.Rat, len(ess))
	var soonest = make([]*big.Rat, len(ess))
	var lcd = ess.lcd()

	for i := range ess ***REMOVED***
		numerators[i] = ess[i].length.Num().Int64() * (lcd / ess[i].length.Denom().Int64())
		soonest[i] = big.NewRat(0, 1)
		numeratorChanges[i] = big.NewRat(lcd, numerators[i])
	***REMOVED***

	start := int64(-1)
	var offsets []int64 // TODO we can preallocate this

	for i := int64(0); i < lcd; i++ ***REMOVED***
		var iRat = big.NewRat(i, 1)
		for index, value := range soonest ***REMOVED***
			if iRat.Cmp(value) >= 0 ***REMOVED***
				value.Add(value, numeratorChanges[index])
				if ess[index] == matchingSegment ***REMOVED***
					// TODO: this can be done for all segments and then we only get what we care about
					if start < 0 ***REMOVED***
						start = i
					***REMOVED*** else ***REMOVED***
						prev := start
						if len(offsets) > 0 ***REMOVED***
							prev = offsets[len(offsets)-1]
						***REMOVED***
						offsets = append(offsets, i-prev)
					***REMOVED***
					// we can stop iterating the moment len(offsets) + 1 = numerator and we have set start ...
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return start, offsets, lcd, nil
***REMOVED***
