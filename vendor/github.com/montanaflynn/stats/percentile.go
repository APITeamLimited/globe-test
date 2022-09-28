package stats

import "math"

// Percentile finds the relative standing in a slice of floats
func Percentile(input Float64Data, percent float64) (percentile float64, err error) ***REMOVED***

	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	if percent <= 0 || percent > 100 ***REMOVED***
		return math.NaN(), BoundsErr
	***REMOVED***

	// Start by sorting a copy of the slice
	c := sortedCopy(input)

	// Multiply percent by length of input
	index := (percent / 100) * float64(len(c))

	// Check if the index is a whole number
	if index == float64(int64(index)) ***REMOVED***

		// Convert float to int
		i := int(index)

		// Find the value at the index
		percentile = c[i-1]

	***REMOVED*** else if index > 1 ***REMOVED***

		// Convert float to int via truncation
		i := int(index)

		// Find the average of the index and following values
		percentile, _ = Mean(Float64Data***REMOVED***c[i-1], c[i]***REMOVED***)

	***REMOVED*** else ***REMOVED***
		return math.NaN(), BoundsErr
	***REMOVED***

	return percentile, nil

***REMOVED***

// PercentileNearestRank finds the relative standing in a slice of floats using the Nearest Rank method
func PercentileNearestRank(input Float64Data, percent float64) (percentile float64, err error) ***REMOVED***

	// Find the length of items in the slice
	il := input.Len()

	// Return an error for empty slices
	if il == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Return error for less than 0 or greater than 100 percentages
	if percent < 0 || percent > 100 ***REMOVED***
		return math.NaN(), BoundsErr
	***REMOVED***

	// Start by sorting a copy of the slice
	c := sortedCopy(input)

	// Return the last item
	if percent == 100.0 ***REMOVED***
		return c[il-1], nil
	***REMOVED***

	// Find ordinal ranking
	or := int(math.Ceil(float64(il) * percent / 100))

	// Return the item that is in the place of the ordinal rank
	if or == 0 ***REMOVED***
		return c[0], nil
	***REMOVED***
	return c[or-1], nil

***REMOVED***
