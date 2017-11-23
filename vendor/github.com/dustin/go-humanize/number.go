package humanize

/*
Slightly adapted from the source to fit go-humanize.

Author: https://github.com/gorhill
Source: https://gist.github.com/gorhill/5285193

*/

import (
	"math"
	"strconv"
)

var (
	renderFloatPrecisionMultipliers = [...]float64***REMOVED***
		1,
		10,
		100,
		1000,
		10000,
		100000,
		1000000,
		10000000,
		100000000,
		1000000000,
	***REMOVED***

	renderFloatPrecisionRounders = [...]float64***REMOVED***
		0.5,
		0.05,
		0.005,
		0.0005,
		0.00005,
		0.000005,
		0.0000005,
		0.00000005,
		0.000000005,
		0.0000000005,
	***REMOVED***
)

// FormatFloat produces a formatted number as string based on the following user-specified criteria:
// * thousands separator
// * decimal separator
// * decimal precision
//
// Usage: s := RenderFloat(format, n)
// The format parameter tells how to render the number n.
//
// See examples: http://play.golang.org/p/LXc1Ddm1lJ
//
// Examples of format strings, given n = 12345.6789:
// "#,###.##" => "12,345.67"
// "#,###." => "12,345"
// "#,###" => "12345,678"
// "#\u202F###,##" => "12 345,68"
// "#.###,###### => 12.345,678900
// "" (aka default format) => 12,345.67
//
// The highest precision allowed is 9 digits after the decimal symbol.
// There is also a version for integer number, FormatInteger(),
// which is convenient for calls within template.
func FormatFloat(format string, n float64) string ***REMOVED***
	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"
	if math.IsNaN(n) ***REMOVED***
		return "NaN"
	***REMOVED***
	if n > math.MaxFloat64 ***REMOVED***
		return "Infinity"
	***REMOVED***
	if n < -math.MaxFloat64 ***REMOVED***
		return "-Infinity"
	***REMOVED***

	// default format
	precision := 2
	decimalStr := "."
	thousandStr := ","
	positiveStr := ""
	negativeStr := "-"

	if len(format) > 0 ***REMOVED***
		format := []rune(format)

		// If there is an explicit format directive,
		// then default values are these:
		precision = 9
		thousandStr = ""

		// collect indices of meaningful formatting directives
		formatIndx := []int***REMOVED******REMOVED***
		for i, char := range format ***REMOVED***
			if char != '#' && char != '0' ***REMOVED***
				formatIndx = append(formatIndx, i)
			***REMOVED***
		***REMOVED***

		if len(formatIndx) > 0 ***REMOVED***
			// Directive at index 0:
			//   Must be a '+'
			//   Raise an error if not the case
			// index: 0123456789
			//        +0.000,000
			//        +000,000.0
			//        +0000.00
			//        +0000
			if formatIndx[0] == 0 ***REMOVED***
				if format[formatIndx[0]] != '+' ***REMOVED***
					panic("RenderFloat(): invalid positive sign directive")
				***REMOVED***
				positiveStr = "+"
				formatIndx = formatIndx[1:]
			***REMOVED***

			// Two directives:
			//   First is thousands separator
			//   Raise an error if not followed by 3-digit
			// 0123456789
			// 0.000,000
			// 000,000.00
			if len(formatIndx) == 2 ***REMOVED***
				if (formatIndx[1] - formatIndx[0]) != 4 ***REMOVED***
					panic("RenderFloat(): thousands separator directive must be followed by 3 digit-specifiers")
				***REMOVED***
				thousandStr = string(format[formatIndx[0]])
				formatIndx = formatIndx[1:]
			***REMOVED***

			// One directive:
			//   Directive is decimal separator
			//   The number of digit-specifier following the separator indicates wanted precision
			// 0123456789
			// 0.00
			// 000,0000
			if len(formatIndx) == 1 ***REMOVED***
				decimalStr = string(format[formatIndx[0]])
				precision = len(format) - formatIndx[0] - 1
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// generate sign part
	var signStr string
	if n >= 0.000000001 ***REMOVED***
		signStr = positiveStr
	***REMOVED*** else if n <= -0.000000001 ***REMOVED***
		signStr = negativeStr
		n = -n
	***REMOVED*** else ***REMOVED***
		signStr = ""
		n = 0.0
	***REMOVED***

	// split number into integer and fractional parts
	intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])

	// generate integer part string
	intStr := strconv.FormatInt(int64(intf), 10)

	// add thousand separator if required
	if len(thousandStr) > 0 ***REMOVED***
		for i := len(intStr); i > 3; ***REMOVED***
			i -= 3
			intStr = intStr[:i] + thousandStr + intStr[i:]
		***REMOVED***
	***REMOVED***

	// no fractional part, we can leave now
	if precision == 0 ***REMOVED***
		return signStr + intStr
	***REMOVED***

	// generate fractional part
	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	// may need padding
	if len(fracStr) < precision ***REMOVED***
		fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
	***REMOVED***

	return signStr + intStr + decimalStr + fracStr
***REMOVED***

// FormatInteger produces a formatted number as string.
// See FormatFloat.
func FormatInteger(format string, n int) string ***REMOVED***
	return FormatFloat(format, float64(n))
***REMOVED***
