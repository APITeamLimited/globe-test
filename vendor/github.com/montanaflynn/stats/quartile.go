package stats

import "math"

// Quartiles holds the three quartile points
type Quartiles struct ***REMOVED***
	Q1 float64
	Q2 float64
	Q3 float64
***REMOVED***

// Quartile returns the three quartile points from a slice of data
func Quartile(input Float64Data) (Quartiles, error) ***REMOVED***

	il := input.Len()
	if il == 0 ***REMOVED***
		return Quartiles***REMOVED******REMOVED***, EmptyInput
	***REMOVED***

	// Start by sorting a copy of the slice
	copy := sortedCopy(input)

	// Find the cutoff places depeding on if
	// the input slice length is even or odd
	var c1 int
	var c2 int
	if il%2 == 0 ***REMOVED***
		c1 = il / 2
		c2 = il / 2
	***REMOVED*** else ***REMOVED***
		c1 = (il - 1) / 2
		c2 = c1 + 1
	***REMOVED***

	// Find the Medians with the cutoff points
	Q1, _ := Median(copy[:c1])
	Q2, _ := Median(copy)
	Q3, _ := Median(copy[c2:])

	return Quartiles***REMOVED***Q1, Q2, Q3***REMOVED***, nil

***REMOVED***

// InterQuartileRange finds the range between Q1 and Q3
func InterQuartileRange(input Float64Data) (float64, error) ***REMOVED***
	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***
	qs, _ := Quartile(input)
	iqr := qs.Q3 - qs.Q1
	return iqr, nil
***REMOVED***

// Midhinge finds the average of the first and third quartiles
func Midhinge(input Float64Data) (float64, error) ***REMOVED***
	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***
	qs, _ := Quartile(input)
	mh := (qs.Q1 + qs.Q3) / 2
	return mh, nil
***REMOVED***

// Trimean finds the average of the median and the midhinge
func Trimean(input Float64Data) (float64, error) ***REMOVED***
	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	c := sortedCopy(input)
	q, _ := Quartile(c)

	return (q.Q1 + (q.Q2 * 2) + q.Q3) / 4, nil
***REMOVED***
