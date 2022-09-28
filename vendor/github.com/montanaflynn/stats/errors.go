package stats

type statsErr struct ***REMOVED***
	err string
***REMOVED***

func (s statsErr) Error() string ***REMOVED***
	return s.err
***REMOVED***

// These are the package-wide error values.
// All error identification should use these values.
var (
	EmptyInput  = statsErr***REMOVED***"Input must not be empty."***REMOVED***
	SampleSize  = statsErr***REMOVED***"Samples number must be less than input length."***REMOVED***
	NaNErr      = statsErr***REMOVED***"Not a number"***REMOVED***
	NegativeErr = statsErr***REMOVED***"Slice must not contain negative values."***REMOVED***
	ZeroErr     = statsErr***REMOVED***"Slice must not contain zero values."***REMOVED***
	BoundsErr   = statsErr***REMOVED***"Input is outside of range."***REMOVED***
	SizeErr     = statsErr***REMOVED***"Slices must be the same length."***REMOVED***
	InfValue    = statsErr***REMOVED***"Value is infinite."***REMOVED***
)
