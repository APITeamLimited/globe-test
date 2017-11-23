package validator

import "regexp"

const (
	alphaRegexString          = "^[a-zA-Z]+$"
	alphaNumericRegexString   = "^[a-zA-Z0-9]+$"
	numericRegexString        = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	numberRegexString         = "^[0-9]+$"
	hexadecimalRegexString    = "^[0-9a-fA-F]+$"
	hexcolorRegexString       = "^#(?:[0-9a-fA-F]***REMOVED***3***REMOVED***|[0-9a-fA-F]***REMOVED***6***REMOVED***)$"
	rgbRegexString            = "^rgb\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*\\)$"
	rgbaRegexString           = "^rgba\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	hslRegexString            = "^hsl\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*\\)$"
	hslaRegexString           = "^hsla\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0|[1-9]\\d?|100)%)\\s*,\\s*(?:(?:0.[1-9]*)|[01])\\s*\\)$"
	emailRegexString          = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`***REMOVED***\\|***REMOVED***~]|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`***REMOVED***\\|***REMOVED***~]|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])|(?:(?:[a-zA-Z]|\\d|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])*(?:[a-zA-Z]|\\d|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])))\\.)+(?:(?:[a-zA-Z]|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])|(?:(?:[a-zA-Z]|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])*(?:[a-zA-Z]|[\\x***REMOVED***00A0***REMOVED***-\\x***REMOVED***D7FF***REMOVED***\\x***REMOVED***F900***REMOVED***-\\x***REMOVED***FDCF***REMOVED***\\x***REMOVED***FDF0***REMOVED***-\\x***REMOVED***FFEF***REMOVED***])))\\.?$"
	base64RegexString         = "^(?:[A-Za-z0-9+\\/]***REMOVED***4***REMOVED***)*(?:[A-Za-z0-9+\\/]***REMOVED***2***REMOVED***==|[A-Za-z0-9+\\/]***REMOVED***3***REMOVED***=|[A-Za-z0-9+\\/]***REMOVED***4***REMOVED***)$"
	iSBN10RegexString         = "^(?:[0-9]***REMOVED***9***REMOVED***X|[0-9]***REMOVED***10***REMOVED***)$"
	iSBN13RegexString         = "^(?:(?:97(?:8|9))[0-9]***REMOVED***10***REMOVED***)$"
	uUID3RegexString          = "^[0-9a-f]***REMOVED***8***REMOVED***-[0-9a-f]***REMOVED***4***REMOVED***-3[0-9a-f]***REMOVED***3***REMOVED***-[0-9a-f]***REMOVED***4***REMOVED***-[0-9a-f]***REMOVED***12***REMOVED***$"
	uUID4RegexString          = "^[0-9a-f]***REMOVED***8***REMOVED***-[0-9a-f]***REMOVED***4***REMOVED***-4[0-9a-f]***REMOVED***3***REMOVED***-[89ab][0-9a-f]***REMOVED***3***REMOVED***-[0-9a-f]***REMOVED***12***REMOVED***$"
	uUID5RegexString          = "^[0-9a-f]***REMOVED***8***REMOVED***-[0-9a-f]***REMOVED***4***REMOVED***-5[0-9a-f]***REMOVED***3***REMOVED***-[89ab][0-9a-f]***REMOVED***3***REMOVED***-[0-9a-f]***REMOVED***12***REMOVED***$"
	uUIDRegexString           = "^[0-9a-f]***REMOVED***8***REMOVED***-[0-9a-f]***REMOVED***4***REMOVED***-[0-9a-f]***REMOVED***4***REMOVED***-[0-9a-f]***REMOVED***4***REMOVED***-[0-9a-f]***REMOVED***12***REMOVED***$"
	aSCIIRegexString          = "^[\x00-\x7F]*$"
	printableASCIIRegexString = "^[\x20-\x7E]*$"
	multibyteRegexString      = "[^\x00-\x7F]"
	dataURIRegexString        = "^data:.+\\/(.+);base64$"
	latitudeRegexString       = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	longitudeRegexString      = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	sSNRegexString            = `^\d***REMOVED***3***REMOVED***[- ]?\d***REMOVED***2***REMOVED***[- ]?\d***REMOVED***4***REMOVED***$`
)

var (
	alphaRegex          = regexp.MustCompile(alphaRegexString)
	alphaNumericRegex   = regexp.MustCompile(alphaNumericRegexString)
	numericRegex        = regexp.MustCompile(numericRegexString)
	numberRegex         = regexp.MustCompile(numberRegexString)
	hexadecimalRegex    = regexp.MustCompile(hexadecimalRegexString)
	hexcolorRegex       = regexp.MustCompile(hexcolorRegexString)
	rgbRegex            = regexp.MustCompile(rgbRegexString)
	rgbaRegex           = regexp.MustCompile(rgbaRegexString)
	hslRegex            = regexp.MustCompile(hslRegexString)
	hslaRegex           = regexp.MustCompile(hslaRegexString)
	emailRegex          = regexp.MustCompile(emailRegexString)
	base64Regex         = regexp.MustCompile(base64RegexString)
	iSBN10Regex         = regexp.MustCompile(iSBN10RegexString)
	iSBN13Regex         = regexp.MustCompile(iSBN13RegexString)
	uUID3Regex          = regexp.MustCompile(uUID3RegexString)
	uUID4Regex          = regexp.MustCompile(uUID4RegexString)
	uUID5Regex          = regexp.MustCompile(uUID5RegexString)
	uUIDRegex           = regexp.MustCompile(uUIDRegexString)
	aSCIIRegex          = regexp.MustCompile(aSCIIRegexString)
	printableASCIIRegex = regexp.MustCompile(printableASCIIRegexString)
	multibyteRegex      = regexp.MustCompile(multibyteRegexString)
	dataURIRegex        = regexp.MustCompile(dataURIRegexString)
	latitudeRegex       = regexp.MustCompile(latitudeRegexString)
	longitudeRegex      = regexp.MustCompile(longitudeRegexString)
	sSNRegex            = regexp.MustCompile(sSNRegexString)
)
