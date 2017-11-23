package humanize

import "strconv"

func stripTrailingZeros(s string) string ***REMOVED***
	offset := len(s) - 1
	for offset > 0 ***REMOVED***
		if s[offset] == '.' ***REMOVED***
			offset--
			break
		***REMOVED***
		if s[offset] != '0' ***REMOVED***
			break
		***REMOVED***
		offset--
	***REMOVED***
	return s[:offset+1]
***REMOVED***

// Ftoa converts a float to a string with no trailing zeros.
func Ftoa(num float64) string ***REMOVED***
	return stripTrailingZeros(strconv.FormatFloat(num, 'f', 6, 64))
***REMOVED***
