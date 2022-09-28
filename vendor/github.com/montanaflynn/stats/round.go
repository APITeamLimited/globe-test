package stats

import "math"

// Round a float to a specific decimal place or precision
func Round(input float64, places int) (rounded float64, err error) ***REMOVED***

	// If the float is not a number
	if math.IsNaN(input) ***REMOVED***
		return math.NaN(), NaNErr
	***REMOVED***

	// Find out the actual sign and correct the input for later
	sign := 1.0
	if input < 0 ***REMOVED***
		sign = -1
		input *= -1
	***REMOVED***

	// Use the places arg to get the amount of precision wanted
	precision := math.Pow(10, float64(places))

	// Find the decimal place we are looking to round
	digit := input * precision

	// Get the actual decimal number as a fraction to be compared
	_, decimal := math.Modf(digit)

	// If the decimal is less than .5 we round down otherwise up
	if decimal >= 0.5 ***REMOVED***
		rounded = math.Ceil(digit)
	***REMOVED*** else ***REMOVED***
		rounded = math.Floor(digit)
	***REMOVED***

	// Finally we do the math to actually create a rounded number
	return rounded / precision * sign, nil
***REMOVED***
