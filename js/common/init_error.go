package common

// InitContextError is an error that happened during the a test init context
type InitContextError string

// NewInitContextError returns a new InitContextError with the provided message
func NewInitContextError(msg string) InitContextError ***REMOVED***
	return (InitContextError)(msg)
***REMOVED***

func (i InitContextError) Error() string ***REMOVED***
	return (string)(i)
***REMOVED***

func (i InitContextError) String() string ***REMOVED***
	return (string)(i)
***REMOVED***
