package humanize

import "strconv"

// Ordinal gives you the input number in a rank/ordinal format.
//
// Ordinal(3) -> 3rd
func Ordinal(x int) string ***REMOVED***
	suffix := "th"
	switch x % 10 ***REMOVED***
	case 1:
		if x%100 != 11 ***REMOVED***
			suffix = "st"
		***REMOVED***
	case 2:
		if x%100 != 12 ***REMOVED***
			suffix = "nd"
		***REMOVED***
	case 3:
		if x%100 != 13 ***REMOVED***
			suffix = "rd"
		***REMOVED***
	***REMOVED***
	return strconv.Itoa(x) + suffix
***REMOVED***
