package mapstructure

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Error implements the error interface and can represents multiple
// errors that occur in the course of a single decode.
type Error struct ***REMOVED***
	Errors []string
***REMOVED***

func (e *Error) Error() string ***REMOVED***
	points := make([]string, len(e.Errors))
	for i, err := range e.Errors ***REMOVED***
		points[i] = fmt.Sprintf("* %s", err)
	***REMOVED***

	sort.Strings(points)
	return fmt.Sprintf(
		"%d error(s) decoding:\n\n%s",
		len(e.Errors), strings.Join(points, "\n"))
***REMOVED***

// WrappedErrors implements the errwrap.Wrapper interface to make this
// return value more useful with the errwrap and go-multierror libraries.
func (e *Error) WrappedErrors() []error ***REMOVED***
	if e == nil ***REMOVED***
		return nil
	***REMOVED***

	result := make([]error, len(e.Errors))
	for i, e := range e.Errors ***REMOVED***
		result[i] = errors.New(e)
	***REMOVED***

	return result
***REMOVED***

func appendErrors(errors []string, err error) []string ***REMOVED***
	switch e := err.(type) ***REMOVED***
	case *Error:
		return append(errors, e.Errors...)
	default:
		return append(errors, e.Error())
	***REMOVED***
***REMOVED***
