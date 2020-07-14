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
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringToES(t *testing.T, str string) *ExecutionSegment ***REMOVED***
	es := new(ExecutionSegment)
	require.NoError(t, es.UnmarshalText([]byte(str)))
	return es
***REMOVED***

func TestExecutionSegmentEquals(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("nil segment to full", func(t *testing.T) ***REMOVED***
		var nilEs *ExecutionSegment
		fullEs := stringToES(t, "0:1")
		require.True(t, nilEs.Equal(fullEs))
		require.True(t, fullEs.Equal(nilEs))
	***REMOVED***)

	t.Run("To it's self", func(t *testing.T) ***REMOVED***
		es := stringToES(t, "1/2:2/3")
		require.True(t, es.Equal(es))
	***REMOVED***)
***REMOVED***

func TestExecutionSegmentNew(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("from is below zero", func(t *testing.T) ***REMOVED***
		_, err := NewExecutionSegment(big.NewRat(-1, 1), big.NewRat(1, 1))
		require.Error(t, err)
	***REMOVED***)
	t.Run("to is more than 1", func(t *testing.T) ***REMOVED***
		_, err := NewExecutionSegment(big.NewRat(0, 1), big.NewRat(2, 1))
		require.Error(t, err)
	***REMOVED***)
	t.Run("from is smaller than to", func(t *testing.T) ***REMOVED***
		_, err := NewExecutionSegment(big.NewRat(1, 2), big.NewRat(1, 3))
		require.Error(t, err)
	***REMOVED***)

	t.Run("from is equal to 'to'", func(t *testing.T) ***REMOVED***
		_, err := NewExecutionSegment(big.NewRat(1, 2), big.NewRat(1, 2))
		require.Error(t, err)
	***REMOVED***)
	t.Run("ok", func(t *testing.T) ***REMOVED***
		_, err := NewExecutionSegment(big.NewRat(0, 1), big.NewRat(1, 1))
		require.NoError(t, err)
	***REMOVED***)
***REMOVED***

func TestExecutionSegmentUnmarshalText(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct ***REMOVED***
		input  string
		output *ExecutionSegment
		isErr  bool
	***REMOVED******REMOVED***
		***REMOVED***input: "0:1", output: &ExecutionSegment***REMOVED***from: zeroRat, to: oneRat***REMOVED******REMOVED***,
		***REMOVED***input: "0.5:0.75", output: &ExecutionSegment***REMOVED***from: big.NewRat(1, 2), to: big.NewRat(3, 4)***REMOVED******REMOVED***,
		***REMOVED***input: "1/2:3/4", output: &ExecutionSegment***REMOVED***from: big.NewRat(1, 2), to: big.NewRat(3, 4)***REMOVED******REMOVED***,
		***REMOVED***input: "50%:75%", output: &ExecutionSegment***REMOVED***from: big.NewRat(1, 2), to: big.NewRat(3, 4)***REMOVED******REMOVED***,
		***REMOVED***input: "2/4:75%", output: &ExecutionSegment***REMOVED***from: big.NewRat(1, 2), to: big.NewRat(3, 4)***REMOVED******REMOVED***,
		***REMOVED***input: "75%", output: &ExecutionSegment***REMOVED***from: zeroRat, to: big.NewRat(3, 4)***REMOVED******REMOVED***,
		***REMOVED***input: "125%", isErr: true***REMOVED***,
		***REMOVED***input: "1a5%", isErr: true***REMOVED***,
		***REMOVED***input: "1a5", isErr: true***REMOVED***,
		***REMOVED***input: "1a5%:2/3", isErr: true***REMOVED***,
		***REMOVED***input: "125%:250%", isErr: true***REMOVED***,
		***REMOVED***input: "55%:50%", isErr: true***REMOVED***,
		// TODO add more strange or not so strange cases
	***REMOVED***
	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(testCase.input, func(t *testing.T) ***REMOVED***
			es := new(ExecutionSegment)
			err := es.UnmarshalText([]byte(testCase.input))
			if testCase.isErr ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			require.True(t, es.Equal(testCase.output))

			// see if unmarshalling a stringified segment gets you back the same segment
			err = es.UnmarshalText([]byte(es.String()))
			require.NoError(t, err)
			require.True(t, es.Equal(testCase.output))
		***REMOVED***)
	***REMOVED***

	t.Run("Unmarshal nilSegment.String", func(t *testing.T) ***REMOVED***
		var nilEs *ExecutionSegment
		nilEsStr := nilEs.String()
		require.Equal(t, "0:1", nilEsStr)

		es := new(ExecutionSegment)
		err := es.UnmarshalText([]byte(nilEsStr))
		require.NoError(t, err)
		require.True(t, es.Equal(nilEs))
	***REMOVED***)
***REMOVED***

func TestExecutionSegmentSplit(t *testing.T) ***REMOVED***
	t.Parallel()

	var nilEs *ExecutionSegment
	_, err := nilEs.Split(-1)
	require.Error(t, err)

	_, err = nilEs.Split(0)
	require.Error(t, err)

	segments, err := nilEs.Split(1)
	require.NoError(t, err)
	require.Len(t, segments, 1)
	assert.Equal(t, "0:1", segments[0].String())

	segments, err = nilEs.Split(2)
	require.NoError(t, err)
	require.Len(t, segments, 2)
	assert.Equal(t, "0:1/2", segments[0].String())
	assert.Equal(t, "1/2:1", segments[1].String())

	segments, err = nilEs.Split(3)
	require.NoError(t, err)
	require.Len(t, segments, 3)
	assert.Equal(t, "0:1/3", segments[0].String())
	assert.Equal(t, "1/3:2/3", segments[1].String())
	assert.Equal(t, "2/3:1", segments[2].String())

	secondQuarter, err := NewExecutionSegment(big.NewRat(1, 4), big.NewRat(2, 4))
	require.NoError(t, err)

	segments, err = secondQuarter.Split(1)
	require.NoError(t, err)
	require.Len(t, segments, 1)
	assert.Equal(t, "1/4:1/2", segments[0].String())

	segments, err = secondQuarter.Split(2)
	require.NoError(t, err)
	require.Len(t, segments, 2)
	assert.Equal(t, "1/4:3/8", segments[0].String())
	assert.Equal(t, "3/8:1/2", segments[1].String())

	segments, err = secondQuarter.Split(3)
	require.NoError(t, err)
	require.Len(t, segments, 3)
	assert.Equal(t, "1/4:1/3", segments[0].String())
	assert.Equal(t, "1/3:5/12", segments[1].String())
	assert.Equal(t, "5/12:1/2", segments[2].String())

	segments, err = secondQuarter.Split(4)
	require.NoError(t, err)
	require.Len(t, segments, 4)
	assert.Equal(t, "1/4:5/16", segments[0].String())
	assert.Equal(t, "5/16:3/8", segments[1].String())
	assert.Equal(t, "3/8:7/16", segments[2].String())
	assert.Equal(t, "7/16:1/2", segments[3].String())
***REMOVED***

func TestExecutionSegmentFailures(t *testing.T) ***REMOVED***
	t.Parallel()
	es := new(ExecutionSegment)
	require.NoError(t, es.UnmarshalText([]byte("0:0.25")))
	require.Equal(t, int64(1), es.Scale(2))
	require.Equal(t, int64(1), es.Scale(3))

	require.NoError(t, es.UnmarshalText([]byte("0.25:0.5")))
	require.Equal(t, int64(0), es.Scale(2))
	require.Equal(t, int64(1), es.Scale(3))

	require.NoError(t, es.UnmarshalText([]byte("0.5:0.75")))
	require.Equal(t, int64(1), es.Scale(2))
	require.Equal(t, int64(0), es.Scale(3))

	require.NoError(t, es.UnmarshalText([]byte("0.75:1")))
	require.Equal(t, int64(0), es.Scale(2))
	require.Equal(t, int64(1), es.Scale(3))
***REMOVED***

func TestExecutionTupleScale(t *testing.T) ***REMOVED***
	t.Parallel()
	es := new(ExecutionSegment)
	require.NoError(t, es.UnmarshalText([]byte("0.5")))
	et, err := NewExecutionTuple(es, nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), et.ScaleInt64(2))
	require.Equal(t, int64(2), et.ScaleInt64(3))

	require.NoError(t, es.UnmarshalText([]byte("0.5:1.0")))
	et, err = NewExecutionTuple(es, nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), et.ScaleInt64(2))
	require.Equal(t, int64(1), et.ScaleInt64(3))

	ess, err := NewExecutionSegmentSequenceFromString("0,0.5,1")
	require.NoError(t, err)
	require.NoError(t, es.UnmarshalText([]byte("0.5")))
	et, err = NewExecutionTuple(es, &ess)
	require.NoError(t, err)
	require.Equal(t, int64(1), et.ScaleInt64(2))
	require.Equal(t, int64(2), et.ScaleInt64(3))

	require.NoError(t, es.UnmarshalText([]byte("0.5:1.0")))
	et, err = NewExecutionTuple(es, &ess)
	require.NoError(t, err)
	require.Equal(t, int64(1), et.ScaleInt64(2))
	require.Equal(t, int64(1), et.ScaleInt64(3))
***REMOVED***

func TestBigScale(t *testing.T) ***REMOVED***
	es := new(ExecutionSegment)
	ess, err := NewExecutionSegmentSequenceFromString("0,7/20,7/10,1")
	require.NoError(t, err)
	require.NoError(t, es.UnmarshalText([]byte("0:7/20")))
	et, err := NewExecutionTuple(es, &ess)
	require.NoError(t, err)
	require.Equal(t, int64(18), et.ScaleInt64(50))
***REMOVED***

func TestExecutionSegmentCopyScaleRat(t *testing.T) ***REMOVED***
	t.Parallel()
	es := new(ExecutionSegment)
	twoRat := big.NewRat(2, 1)
	threeRat := big.NewRat(3, 1)
	require.NoError(t, es.UnmarshalText([]byte("0.5")))
	require.Equal(t, oneRat, es.CopyScaleRat(twoRat))
	require.Equal(t, big.NewRat(3, 2), es.CopyScaleRat(threeRat))

	require.NoError(t, es.UnmarshalText([]byte("0.5:1.0")))
	require.Equal(t, oneRat, es.CopyScaleRat(twoRat))
	require.Equal(t, big.NewRat(3, 2), es.CopyScaleRat(threeRat))

	var nilEs *ExecutionSegment
	require.Equal(t, twoRat, nilEs.CopyScaleRat(twoRat))
	require.Equal(t, threeRat, nilEs.CopyScaleRat(threeRat))
***REMOVED***

func TestExecutionSegmentInPlaceScaleRat(t *testing.T) ***REMOVED***
	t.Parallel()
	es := new(ExecutionSegment)
	twoRat := big.NewRat(2, 1)
	threeRat := big.NewRat(3, 1)
	threeSecondsRat := big.NewRat(3, 2)
	require.NoError(t, es.UnmarshalText([]byte("0.5")))
	require.Equal(t, oneRat, es.InPlaceScaleRat(twoRat))
	require.Equal(t, oneRat, twoRat)
	require.Equal(t, threeSecondsRat, es.InPlaceScaleRat(threeRat))
	require.Equal(t, threeSecondsRat, threeRat)

	es = stringToES(t, "0.5:1.0")
	twoRat = big.NewRat(2, 1)
	threeRat = big.NewRat(3, 1)
	require.Equal(t, oneRat, es.InPlaceScaleRat(twoRat))
	require.Equal(t, oneRat, twoRat)
	require.Equal(t, threeSecondsRat, es.InPlaceScaleRat(threeRat))
	require.Equal(t, threeSecondsRat, threeRat)

	var nilEs *ExecutionSegment
	twoRat = big.NewRat(2, 1)
	threeRat = big.NewRat(3, 1)
	require.Equal(t, big.NewRat(2, 1), nilEs.InPlaceScaleRat(twoRat))
	require.Equal(t, big.NewRat(2, 1), twoRat)
	require.Equal(t, big.NewRat(3, 1), nilEs.InPlaceScaleRat(threeRat))
	require.Equal(t, big.NewRat(3, 1), threeRat)
***REMOVED***

func TestExecutionSegmentSubSegment(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct ***REMOVED***
		name              string
		base, sub, result *ExecutionSegment
	***REMOVED******REMOVED***
		// TODO add more strange or not so strange cases
		***REMOVED***
			name:   "nil base",
			base:   (*ExecutionSegment)(nil),
			sub:    stringToES(t, "0.2:0.3"),
			result: stringToES(t, "0.2:0.3"),
		***REMOVED***,

		***REMOVED***
			name:   "nil sub",
			base:   stringToES(t, "0.2:0.3"),
			sub:    (*ExecutionSegment)(nil),
			result: stringToES(t, "0.2:0.3"),
		***REMOVED***,
		***REMOVED***
			name:   "doc example",
			base:   stringToES(t, "1/2:1"),
			sub:    stringToES(t, "0:1/2"),
			result: stringToES(t, "1/2:3/4"),
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) ***REMOVED***
			require.Equal(t, testCase.result, testCase.base.SubSegment(testCase.sub))
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSplitBadSegment(t *testing.T) ***REMOVED***
	t.Parallel()
	es := &ExecutionSegment***REMOVED***from: oneRat, to: zeroRat***REMOVED***
	_, err := es.Split(5)
	require.Error(t, err)
***REMOVED***

func TestSegmentExecutionFloatLength(t *testing.T) ***REMOVED***
	t.Parallel()
	t.Run("nil has 1.0", func(t *testing.T) ***REMOVED***
		var nilEs *ExecutionSegment
		require.Equal(t, 1.0, nilEs.FloatLength())
	***REMOVED***)

	testCases := []struct ***REMOVED***
		es       *ExecutionSegment
		expected float64
	***REMOVED******REMOVED***
		// TODO add more strange or not so strange cases
		***REMOVED***
			es:       stringToES(t, "1/2:1"),
			expected: 0.5,
		***REMOVED***,
		***REMOVED***
			es:       stringToES(t, "1/3:1"),
			expected: 0.66666,
		***REMOVED***,

		***REMOVED***
			es:       stringToES(t, "0:1/2"),
			expected: 0.5,
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(testCase.es.String(), func(t *testing.T) ***REMOVED***
			require.InEpsilon(t, testCase.expected, testCase.es.FloatLength(), 0.001)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutionSegmentSequences(t *testing.T) ***REMOVED***
	t.Parallel()

	_, err := NewExecutionSegmentSequence(stringToES(t, "0:1/3"), stringToES(t, "1/2:1"))
	assert.Error(t, err)
***REMOVED***

func TestExecutionSegmentStringSequences(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct ***REMOVED***
		seq         string
		expSegments []string
		expError    bool
		canReverse  bool
		// TODO: checks for least common denominator and maybe striped partitioning
	***REMOVED******REMOVED***
		***REMOVED***seq: "", expSegments: nil***REMOVED***,
		***REMOVED***seq: "0.5", expError: true***REMOVED***,
		***REMOVED***seq: "1,1", expError: true***REMOVED***,
		***REMOVED***seq: "-0.5,1", expError: true***REMOVED***,
		***REMOVED***seq: "1/2,1/2", expError: true***REMOVED***,
		***REMOVED***seq: "1/2,1/3", expError: true***REMOVED***,
		***REMOVED***seq: "0,1,1/2", expError: true***REMOVED***,
		***REMOVED***seq: "0.5,1", expSegments: []string***REMOVED***"1/2:1"***REMOVED******REMOVED***,
		***REMOVED***seq: "1/2,1", expSegments: []string***REMOVED***"1/2:1"***REMOVED***, canReverse: true***REMOVED***,
		***REMOVED***seq: "1/3,2/3", expSegments: []string***REMOVED***"1/3:2/3"***REMOVED***, canReverse: true***REMOVED***,
		***REMOVED***seq: "0,1/3,2/3", expSegments: []string***REMOVED***"0:1/3", "1/3:2/3"***REMOVED***, canReverse: true***REMOVED***,
		***REMOVED***seq: "0,1/3,2/3,1", expSegments: []string***REMOVED***"0:1/3", "1/3:2/3", "2/3:1"***REMOVED***, canReverse: true***REMOVED***,
		***REMOVED***seq: "0.5,0.7", expSegments: []string***REMOVED***"1/2:7/10"***REMOVED******REMOVED***,
		***REMOVED***seq: "0.5,0.7,1", expSegments: []string***REMOVED***"1/2:7/10", "7/10:1"***REMOVED******REMOVED***,
		***REMOVED***seq: "0,1/13,2/13,1/3,1/2,3/4,1", expSegments: []string***REMOVED***
			"0:1/13", "1/13:2/13", "2/13:1/3", "1/3:1/2", "1/2:3/4", "3/4:1",
		***REMOVED***, canReverse: true***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.seq, func(t *testing.T) ***REMOVED***
			result, err := NewExecutionSegmentSequenceFromString(tc.seq)
			if tc.expError ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			require.Equal(t, len(tc.expSegments), len(result))
			for i, expStrSeg := range tc.expSegments ***REMOVED***
				expSeg, errl := NewExecutionSegmentFromString(expStrSeg)
				require.NoError(t, errl)
				assert.Truef(t, expSeg.Equal(result[i]), "Segment %d (%s) should be equal to %s", i, result[i], expSeg)
			***REMOVED***
			if tc.canReverse ***REMOVED***
				assert.Equal(t, result.String(), tc.seq)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// Return a randomly distributed sequence of n amount of
// execution segments whose length totals 1.
func generateRandomSequence(t testing.TB, n, m int64, r *rand.Rand) ExecutionSegmentSequence ***REMOVED***
	var err error
	ess := ExecutionSegmentSequence(make([]*ExecutionSegment, n))
	numerators := make([]int64, n)
	var denominator int64
	for i := int64(0); i < n; i++ ***REMOVED***
		numerators[i] = 1 + r.Int63n(m)
		denominator += numerators[i]
	***REMOVED***
	from := big.NewRat(0, 1)
	for i := int64(0); i < n; i++ ***REMOVED***
		to := new(big.Rat).Add(big.NewRat(numerators[i], denominator), from)
		ess[i], err = NewExecutionSegment(from, to)
		require.NoError(t, err)
		from = to
	***REMOVED***

	return ess
***REMOVED***

// Ensure that the sum of scaling all execution segments in
// the same sequence with scaling factor M results in M itself.
func TestExecutionSegmentScaleConsistency(t *testing.T) ***REMOVED***
	t.Parallel()

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	t.Logf("Random source seeded with %d\n", seed)

	const numTests = 10
	for i := 0; i < numTests; i++ ***REMOVED***
		scale := rand.Int31n(99) + 2
		seq := generateRandomSequence(t, r.Int63n(9)+2, 100, r)

		t.Run(fmt.Sprintf("%d_%s", scale, seq), func(t *testing.T) ***REMOVED***
			var total int64
			for _, segment := range seq ***REMOVED***
				total += segment.Scale(int64(scale))
			***REMOVED***
			assert.Equal(t, int64(scale), total)
		***REMOVED***)
	***REMOVED***
***REMOVED***

// Ensure that the sum of scaling all execution segments in
// the same sequence with scaling factor M results in M itself.
func TestExecutionTupleScaleConsistency(t *testing.T) ***REMOVED***
	t.Parallel()

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	t.Logf("Random source seeded with %d\n", seed)

	const numTests = 10
	for i := 0; i < numTests; i++ ***REMOVED***
		scale := rand.Int31n(99) + 2
		seq := generateRandomSequence(t, r.Int63n(9)+2, 200, r)

		et, err := NewExecutionTuple(seq[0], &seq)
		require.NoError(t, err)
		t.Run(fmt.Sprintf("%d_%s", scale, seq), func(t *testing.T) ***REMOVED***
			var total int64
			for i, segment := range seq ***REMOVED***
				assert.True(t, segment.Equal(et.Sequence.ExecutionSegmentSequence[i]))
				total += et.Sequence.ScaleInt64(i, int64(scale))
			***REMOVED***
			assert.Equal(t, int64(scale), total)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestExecutionSegmentScaleNoWobble(t *testing.T) ***REMOVED***
	t.Parallel()

	requireSegmentScaleGreater := func(t *testing.T, et *ExecutionTuple) ***REMOVED***
		var i, lastResult int64
		for i = 1; i < 1000; i++ ***REMOVED***
			result := et.ScaleInt64(i)
			require.True(t, result >= lastResult, "%d<%d", result, lastResult)
			lastResult = result
		***REMOVED***
	***REMOVED***

	// Baseline full segment test
	t.Run("0:1", func(t *testing.T) ***REMOVED***
		et, err := NewExecutionTuple(nil, nil)
		require.NoError(t, err)
		requireSegmentScaleGreater(t, et)
	***REMOVED***)

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	t.Logf("Random source seeded with %d\n", seed)

	// Random segments
	const numTests = 10
	for i := 0; i < numTests; i++ ***REMOVED***
		seq := generateRandomSequence(t, r.Int63n(9)+2, 100, r)

		es := seq[rand.Intn(len(seq))]

		et, err := NewExecutionTuple(seq[0], &seq)
		require.NoError(t, err)
		t.Run(es.String(), func(t *testing.T) ***REMOVED***
			requireSegmentScaleGreater(t, et)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestGetStripedOffsets(t *testing.T) ***REMOVED***
	t.Parallel()
	testCases := []struct ***REMOVED***
		seq     string
		seg     string
		start   int64
		offsets []int64
		lcd     int64
	***REMOVED******REMOVED***
		// full sequences
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0:0.3", start: 0, offsets: []int64***REMOVED***4, 3, 3***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0.3:0.5", start: 1, offsets: []int64***REMOVED***4, 6***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0.5:0.6", start: 2, offsets: []int64***REMOVED***10***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0.6:0.7", start: 3, offsets: []int64***REMOVED***10***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0.8:0.9", start: 8, offsets: []int64***REMOVED***10***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0.9:1", start: 9, offsets: []int64***REMOVED***10***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.2,0.5,0.6,0.7,0.8,0.9,1", seg: "0.9:1", start: 9, offsets: []int64***REMOVED***10***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.2,0.5,0.6,0.7,0.8,0.9,1", seg: "0:0.2", start: 1, offsets: []int64***REMOVED***4, 6***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.2,0.5,0.6,0.7,0.8,0.9,1", seg: "0.6:0.7", start: 3, offsets: []int64***REMOVED***10***REMOVED***, lcd: 10***REMOVED***,
		// not full sequences
		***REMOVED***seq: "0,0.2,0.5", seg: "0:0.2", start: 3, offsets: []int64***REMOVED***6, 4***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.2,0.5", seg: "0.2:0.5", start: 1, offsets: []int64***REMOVED***4, 2, 4***REMOVED***, lcd: 10***REMOVED***,
		***REMOVED***seq: "0,2/5,4/5", seg: "0:2/5", start: 0, offsets: []int64***REMOVED***3, 2***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seq: "0,2/5,4/5", seg: "2/5:4/5", start: 1, offsets: []int64***REMOVED***3, 2***REMOVED***, lcd: 5***REMOVED***,
		// no sequence
		***REMOVED***seg: "0:0.2", start: 1, offsets: []int64***REMOVED***5***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seg: "0:1/5", start: 1, offsets: []int64***REMOVED***5***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seg: "0:2/10", start: 1, offsets: []int64***REMOVED***5***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seg: "0:0.4", start: 1, offsets: []int64***REMOVED***2, 3***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seg: "0:2/5", start: 1, offsets: []int64***REMOVED***2, 3***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seg: "2/5:4/5", start: 1, offsets: []int64***REMOVED***3, 2***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seg: "0:4/10", start: 1, offsets: []int64***REMOVED***2, 3***REMOVED***, lcd: 5***REMOVED***,
		***REMOVED***seg: "1/10:5/10", start: 1, offsets: []int64***REMOVED***2, 2, 4, 2***REMOVED***, lcd: 10***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("seq:%s;segment:%s", tc.seq, tc.seg), func(t *testing.T) ***REMOVED***
			ess, err := NewExecutionSegmentSequenceFromString(tc.seq)
			require.NoError(t, err)
			segment, err := NewExecutionSegmentFromString(tc.seg)
			require.NoError(t, err)
			et, err := NewExecutionTuple(segment, &ess)
			require.NoError(t, err)

			start, offsets, lcd := et.GetStripedOffsets()

			assert.Equal(t, tc.start, start)
			assert.Equal(t, tc.offsets, offsets)
			assert.Equal(t, tc.lcd, lcd)

			ess2, err := NewExecutionSegmentSequenceFromString(tc.seq)
			require.NoError(t, err)
			assert.Equal(t, ess.String(), ess2.String())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSequenceLCD(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		seq string
		lcd int64
	***REMOVED******REMOVED***
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.1,0.5,0.6,0.7,0.8,0.9,1", lcd: 10***REMOVED***,
		***REMOVED***seq: "0,0.2,0.5,0.6,0.7,0.8,0.9,1", lcd: 10***REMOVED***,
		***REMOVED***seq: "0,1/3,5/6", lcd: 6***REMOVED***,
		***REMOVED***seq: "0,1/3,4/7", lcd: 21***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("seq:%s", tc.seq), func(t *testing.T) ***REMOVED***
			ess, err := NewExecutionSegmentSequenceFromString(tc.seq)
			require.NoError(t, err)
			require.Equal(t, tc.lcd, ess.LCD())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkGetStripedOffsets(b *testing.B) ***REMOVED***
	lengths := [...]int64***REMOVED***10, 100***REMOVED***
	const seed = 777
	r := rand.New(rand.NewSource(seed))

	for _, length := range lengths ***REMOVED***
		length := length
		b.Run(fmt.Sprintf("length%d,seed%d", length, seed), func(b *testing.B) ***REMOVED***
			sequence := generateRandomSequence(b, length, 100, r)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				segment := sequence[int(r.Int63())%len(sequence)]
				et, err := NewExecutionTuple(segment, &sequence)
				require.NoError(b, err)
				_, _, _ = et.GetStripedOffsets()
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkGetStripedOffsetsEven(b *testing.B) ***REMOVED***
	lengths := [...]int64***REMOVED***10, 100, 1000***REMOVED***
	generateSequence := func(n int64) ExecutionSegmentSequence ***REMOVED***
		var err error
		ess := ExecutionSegmentSequence(make([]*ExecutionSegment, n))
		numerators := make([]int64, n)
		var denominator int64
		for i := int64(0); i < n; i++ ***REMOVED***
			numerators[i] = 1 // nice and simple :)
			denominator += numerators[i]
		***REMOVED***
		ess[0], err = NewExecutionSegment(big.NewRat(0, 1), big.NewRat(numerators[0], denominator))
		require.NoError(b, err)
		for i := int64(1); i < n; i++ ***REMOVED***
			ess[i], err = NewExecutionSegment(ess[i-1].to, new(big.Rat).Add(big.NewRat(numerators[i], denominator), ess[i-1].to))
			require.NoError(b, err, "%d", i)
		***REMOVED***

		return ess
	***REMOVED***

	for _, length := range lengths ***REMOVED***
		length := length
		b.Run(fmt.Sprintf("length%d", length), func(b *testing.B) ***REMOVED***
			sequence := generateSequence(length)
			b.ResetTimer()
			for i := 0; i < b.N; i++ ***REMOVED***
				segment := sequence[111233%len(sequence)]
				et, err := NewExecutionTuple(segment, &sequence)
				require.NoError(b, err)
				_, _, _ = et.GetStripedOffsets()
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestGetNewExecutionTupleBesedOnValue(t *testing.T) ***REMOVED***
	t.Parallel()

	testCases := []struct ***REMOVED***
		seq      string
		seg      string
		value    int64
		expected string
	***REMOVED******REMOVED***
		// full sequences
		***REMOVED***seq: "0,1/3,2/3,1", seg: "0:1/3", value: 20, expected: "0,7/20,7/10,1"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("seq:%s;segment:%s", tc.seq, tc.seg), func(t *testing.T) ***REMOVED***
			ess, err := NewExecutionSegmentSequenceFromString(tc.seq)
			require.NoError(t, err)

			segment, err := NewExecutionSegmentFromString(tc.seg)
			require.NoError(t, err)

			et, err := NewExecutionTuple(segment, &ess)
			require.NoError(t, err)
			newET, err := et.GetNewExecutionTupleFromValue(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, newET.Sequence.String())
		***REMOVED***)
	***REMOVED***
***REMOVED***

func mustNewExecutionSegment(str string) *ExecutionSegment ***REMOVED***
	res, err := NewExecutionSegmentFromString(str)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return res
***REMOVED***

func mustNewExecutionSegmentSequence(str string) *ExecutionSegmentSequence ***REMOVED***
	res, err := NewExecutionSegmentSequenceFromString(str)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return &res
***REMOVED***

func TestNewExecutionTuple(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		seg           *ExecutionSegment
		seq           *ExecutionSegmentSequence
		scaleTests    map[int64]int64
		newScaleTests map[int64]map[int64]int64 // this is for after calling GetNewExecutionSegmentSequenceFromValue
	***REMOVED******REMOVED***
		***REMOVED***
			// both segment and sequence are nil
			scaleTests: map[int64]int64***REMOVED***
				50: 50,
				1:  1,
				0:  0,
			***REMOVED***,
			newScaleTests: map[int64]map[int64]int64***REMOVED***
				50: ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				1:  ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				0:  nil,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			seg: mustNewExecutionSegment("0:1"),
			// nil sequence
			scaleTests: map[int64]int64***REMOVED***
				50: 50,
				1:  1,
				0:  0,
			***REMOVED***,
			newScaleTests: map[int64]map[int64]int64***REMOVED***
				50: ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				1:  ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				0:  nil,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			seg: mustNewExecutionSegment("0:1"),
			seq: mustNewExecutionSegmentSequence("0,1"),
			scaleTests: map[int64]int64***REMOVED***
				50: 50,
				1:  1,
				0:  0,
			***REMOVED***,
			newScaleTests: map[int64]map[int64]int64***REMOVED***
				50: ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				1:  ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				0:  nil,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			seg: mustNewExecutionSegment("0:1"),
			seq: mustNewExecutionSegmentSequence(""),
			scaleTests: map[int64]int64***REMOVED***
				50: 50,
				1:  1,
				0:  0,
			***REMOVED***,
			newScaleTests: map[int64]map[int64]int64***REMOVED***
				50: ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				1:  ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				0:  nil,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			seg: mustNewExecutionSegment("0:1/3"),
			seq: mustNewExecutionSegmentSequence("0,1/3,2/3,1"),
			scaleTests: map[int64]int64***REMOVED***
				50: 17,
				3:  1,
				2:  1,
				1:  1,
				0:  0,
			***REMOVED***,
			newScaleTests: map[int64]map[int64]int64***REMOVED***
				50: ***REMOVED***50: 17, 1: 1, 0: 0***REMOVED***,
				20: ***REMOVED***50: 18, 1: 1, 0: 0***REMOVED***,
				3:  ***REMOVED***50: 17, 1: 1, 0: 0***REMOVED***,
				2:  ***REMOVED***50: 25, 1: 1, 0: 0***REMOVED***,
				1:  ***REMOVED***50: 50, 1: 1, 0: 0***REMOVED***,
				0:  nil,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			seg: mustNewExecutionSegment("1/3:2/3"),
			seq: mustNewExecutionSegmentSequence("0,1/3,2/3,1"),
			scaleTests: map[int64]int64***REMOVED***
				50: 17,
				3:  1,
				2:  1,
				1:  0,
				0:  0,
			***REMOVED***,
			newScaleTests: map[int64]map[int64]int64***REMOVED***
				50: ***REMOVED***50: 17, 1: 0, 0: 0***REMOVED***,
				20: ***REMOVED***50: 17, 1: 0, 0: 0***REMOVED***,
				3:  ***REMOVED***50: 17, 1: 0, 0: 0***REMOVED***,
				2:  ***REMOVED***50: 25, 1: 0, 0: 0***REMOVED***,
				1:  nil,
				0:  nil,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			seg: mustNewExecutionSegment("2/3:1"),
			seq: mustNewExecutionSegmentSequence("0,1/3,2/3,1"),
			scaleTests: map[int64]int64***REMOVED***
				50: 16,
				3:  1,
				2:  0,
				1:  0,
				0:  0,
			***REMOVED***,
			newScaleTests: map[int64]map[int64]int64***REMOVED***
				50: ***REMOVED***50: 16, 1: 0, 0: 0***REMOVED***,
				20: ***REMOVED***50: 15, 1: 0, 0: 0***REMOVED***,
				3:  ***REMOVED***50: 16, 1: 0, 0: 0***REMOVED***,
				2:  nil,
				1:  nil,
				0:  nil,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(fmt.Sprintf("seg:'%s',seq:'%s'", testCase.seg, testCase.seq), func(t *testing.T) ***REMOVED***
			et, err := NewExecutionTuple(testCase.seg, testCase.seq)
			require.NoError(t, err)

			for scaleValue, result := range testCase.scaleTests ***REMOVED***
				require.Equal(t, result, et.ScaleInt64(scaleValue), "%d->%d", scaleValue, result)
			***REMOVED***

			for value, newResult := range testCase.newScaleTests ***REMOVED***
				newET, err := et.GetNewExecutionTupleFromValue(value)
				if newResult == nil ***REMOVED***
					require.Error(t, err)
					continue
				***REMOVED***
				require.NoError(t, err)
				for scaleValue, result := range newResult ***REMOVED***
					require.Equal(t, result, newET.ScaleInt64(scaleValue),
						"GetNewExecutionTupleFromValue(%d)%d->%d", value, scaleValue, result)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func BenchmarkExecutionSegmentScale(b *testing.B) ***REMOVED***
	testCases := []struct ***REMOVED***
		seq string
		seg string
	***REMOVED******REMOVED***
		***REMOVED******REMOVED***,
		***REMOVED***seg: "0:1"***REMOVED***,
		***REMOVED***seq: "0,0.3,0.5,0.6,0.7,0.8,0.9,1", seg: "0:0.3"***REMOVED***,
		***REMOVED***seq: "0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1", seg: "0:0.1"***REMOVED***,
		***REMOVED***seg: "2/5:4/5"***REMOVED***,
		***REMOVED***seg: "2235/5213:4/5"***REMOVED***, // just wanted it to be ugly ;D
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		b.Run(fmt.Sprintf("seq:%s;segment:%s", tc.seq, tc.seg), func(b *testing.B) ***REMOVED***
			ess, err := NewExecutionSegmentSequenceFromString(tc.seq)
			require.NoError(b, err)
			segment, err := NewExecutionSegmentFromString(tc.seg)
			require.NoError(b, err)
			if tc.seg == "" ***REMOVED***
				segment = nil // specifically for the optimization
			***REMOVED***
			et, err := NewExecutionTuple(segment, &ess)
			require.NoError(b, err)
			for _, value := range []int64***REMOVED***5, 5523, 5000000, 67280421310721***REMOVED*** ***REMOVED***
				value := value
				b.Run(fmt.Sprintf("segment.Scale(%d)", value), func(b *testing.B) ***REMOVED***
					for i := 0; i < b.N; i++ ***REMOVED***
						segment.Scale(value)
					***REMOVED***
				***REMOVED***)

				b.Run(fmt.Sprintf("et.Scale(%d)", value), func(b *testing.B) ***REMOVED***
					for i := 0; i < b.N; i++ ***REMOVED***
						et, err = NewExecutionTuple(segment, &ess)
						require.NoError(b, err)
						et.ScaleInt64(value)
					***REMOVED***
				***REMOVED***)

				et.ScaleInt64(1) // precache
				b.Run(fmt.Sprintf("et.Scale(%d) prefilled", value), func(b *testing.B) ***REMOVED***
					for i := 0; i < b.N; i++ ***REMOVED***
						et.ScaleInt64(value)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

// TODO: test with randomized things