package stats

import "math"

// Mean gets the average of a slice of numbers
func Mean(input Float64Data) (float64, error) ***REMOVED***

	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	sum, _ := input.Sum()

	return sum / float64(input.Len()), nil
***REMOVED***

// GeometricMean gets the geometric mean for a slice of numbers
func GeometricMean(input Float64Data) (float64, error) ***REMOVED***

	l := input.Len()
	if l == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Get the product of all the numbers
	var p float64
	for _, n := range input ***REMOVED***
		if p == 0 ***REMOVED***
			p = n
		***REMOVED*** else ***REMOVED***
			p *= n
		***REMOVED***
	***REMOVED***

	// Calculate the geometric mean
	return math.Pow(p, 1/float64(l)), nil
***REMOVED***

// HarmonicMean gets the harmonic mean for a slice of numbers
func HarmonicMean(input Float64Data) (float64, error) ***REMOVED***

	l := input.Len()
	if l == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Get the sum of all the numbers reciprocals and return an
	// error for values that cannot be included in harmonic mean
	var p float64
	for _, n := range input ***REMOVED***
		if n < 0 ***REMOVED***
			return math.NaN(), NegativeErr
		***REMOVED*** else if n == 0 ***REMOVED***
			return math.NaN(), ZeroErr
		***REMOVED***
		p += (1 / n)
	***REMOVED***

	return float64(l) / p, nil
***REMOVED***
