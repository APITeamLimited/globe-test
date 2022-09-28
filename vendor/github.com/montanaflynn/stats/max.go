package stats

import "math"

// Max finds the highest number in a slice
func Max(input Float64Data) (max float64, err error) ***REMOVED***

	// Return an error if there are no numbers
	if input.Len() == 0 ***REMOVED***
		return math.NaN(), EmptyInput
	***REMOVED***

	// Get the first value as the starting point
	max = input.Get(0)

	// Loop and replace higher values
	for i := 1; i < input.Len(); i++ ***REMOVED***
		if input.Get(i) > max ***REMOVED***
			max = input.Get(i)
		***REMOVED***
	***REMOVED***

	return max, nil
***REMOVED***
