package humanize

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

var siPrefixTable = map[float64]string***REMOVED***
	-24: "y", // yocto
	-21: "z", // zepto
	-18: "a", // atto
	-15: "f", // femto
	-12: "p", // pico
	-9:  "n", // nano
	-6:  "µ", // micro
	-3:  "m", // milli
	0:   "",
	3:   "k", // kilo
	6:   "M", // mega
	9:   "G", // giga
	12:  "T", // tera
	15:  "P", // peta
	18:  "E", // exa
	21:  "Z", // zetta
	24:  "Y", // yotta
***REMOVED***

var revSIPrefixTable = revfmap(siPrefixTable)

// revfmap reverses the map and precomputes the power multiplier
func revfmap(in map[float64]string) map[string]float64 ***REMOVED***
	rv := map[string]float64***REMOVED******REMOVED***
	for k, v := range in ***REMOVED***
		rv[v] = math.Pow(10, k)
	***REMOVED***
	return rv
***REMOVED***

var riParseRegex *regexp.Regexp

func init() ***REMOVED***
	ri := `^([\-0-9.]+)\s?([`
	for _, v := range siPrefixTable ***REMOVED***
		ri += v
	***REMOVED***
	ri += `]?)(.*)`

	riParseRegex = regexp.MustCompile(ri)
***REMOVED***

// ComputeSI finds the most appropriate SI prefix for the given number
// and returns the prefix along with the value adjusted to be within
// that prefix.
//
// See also: SI, ParseSI.
//
// e.g. ComputeSI(2.2345e-12) -> (2.2345, "p")
func ComputeSI(input float64) (float64, string) ***REMOVED***
	if input == 0 ***REMOVED***
		return 0, ""
	***REMOVED***
	mag := math.Abs(input)
	exponent := math.Floor(logn(mag, 10))
	exponent = math.Floor(exponent/3) * 3

	value := mag / math.Pow(10, exponent)

	// Handle special case where value is exactly 1000.0
	// Should return 1 M instead of 1000 k
	if value == 1000.0 ***REMOVED***
		exponent += 3
		value = mag / math.Pow(10, exponent)
	***REMOVED***

	value = math.Copysign(value, input)

	prefix := siPrefixTable[exponent]
	return value, prefix
***REMOVED***

// SI returns a string with default formatting.
//
// SI uses Ftoa to format float value, removing trailing zeros.
//
// See also: ComputeSI, ParseSI.
//
// e.g. SI(1000000, "B") -> 1 MB
// e.g. SI(2.2345e-12, "F") -> 2.2345 pF
func SI(input float64, unit string) string ***REMOVED***
	value, prefix := ComputeSI(input)
	return Ftoa(value) + " " + prefix + unit
***REMOVED***

var errInvalid = errors.New("invalid input")

// ParseSI parses an SI string back into the number and unit.
//
// See also: SI, ComputeSI.
//
// e.g. ParseSI("2.2345 pF") -> (2.2345e-12, "F", nil)
func ParseSI(input string) (float64, string, error) ***REMOVED***
	found := riParseRegex.FindStringSubmatch(input)
	if len(found) != 4 ***REMOVED***
		return 0, "", errInvalid
	***REMOVED***
	mag := revSIPrefixTable[found[2]]
	unit := found[3]

	base, err := strconv.ParseFloat(found[1], 64)
	return base * mag, unit, err
***REMOVED***
