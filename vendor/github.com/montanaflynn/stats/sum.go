package stats

import "math"

// Sum adds all the numbers of a slice together
func Sum(input Float64Data) (sum float64, err error) ***REMOVED***

	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Add em up
	for _, n := range input ***REMOVED***
		sum += n
	***REMOVED***

	return sum, nil
***REMOVED***
