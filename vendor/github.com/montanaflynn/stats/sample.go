package stats

import "math/rand"

// Sample returns sample from input with replacement or without
func Sample(input Float64Data, takenum int, replacement bool) ([]float64, error) ***REMOVED***

	if input.Len() == 0 ***REMOVED***
		return nil, EmptyInput
	***REMOVED***

	length := input.Len()
	if replacement ***REMOVED***

		result := Float64Data***REMOVED******REMOVED***
		rand.Seed(unixnano())

		// In every step, randomly take the num for
		for i := 0; i < takenum; i++ ***REMOVED***
			idx := rand.Intn(length)
			result = append(result, input[idx])
		***REMOVED***

		return result, nil

	***REMOVED*** else if !replacement && takenum <= length ***REMOVED***

		rand.Seed(unixnano())

		// Get permutation of number of indexies
		perm := rand.Perm(length)
		result := Float64Data***REMOVED******REMOVED***

		// Get element of input by permutated index
		for _, idx := range perm[0:takenum] ***REMOVED***
			result = append(result, input[idx])
		***REMOVED***

		return result, nil

	***REMOVED***

	return nil, BoundsErr
***REMOVED***
