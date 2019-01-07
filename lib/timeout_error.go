package lib

// TimeoutError is used when somethings timeouts
type TimeoutError string

// NewTimeoutError returns a new TimeoutError reporting that timeout has happened at the provieded
// place
func NewTimeoutError(place string) TimeoutError ***REMOVED***
	return TimeoutError(place)
***REMOVED***

func (t TimeoutError) String() string ***REMOVED***
	return "Timeout during " + (string)(t)
***REMOVED***

func (t TimeoutError) Error() string ***REMOVED***
	return t.String()
***REMOVED***
