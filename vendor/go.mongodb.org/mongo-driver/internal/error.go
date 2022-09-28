// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal

import (
	"fmt"
)

// WrappedError represents an error that contains another error.
type WrappedError interface ***REMOVED***
	// Message gets the basic message of the error.
	Message() string
	// Inner gets the inner error if one exists.
	Inner() error
***REMOVED***

// RolledUpErrorMessage gets a flattened error message.
func RolledUpErrorMessage(err error) string ***REMOVED***
	if wrappedErr, ok := err.(WrappedError); ok ***REMOVED***
		inner := wrappedErr.Inner()
		if inner != nil ***REMOVED***
			return fmt.Sprintf("%s: %s", wrappedErr.Message(), RolledUpErrorMessage(inner))
		***REMOVED***

		return wrappedErr.Message()
	***REMOVED***

	return err.Error()
***REMOVED***

//UnwrapError attempts to unwrap the error down to its root cause.
func UnwrapError(err error) error ***REMOVED***

	switch tErr := err.(type) ***REMOVED***
	case WrappedError:
		return UnwrapError(tErr.Inner())
	case *multiError:
		return UnwrapError(tErr.errors[0])
	***REMOVED***

	return err
***REMOVED***

// WrapError wraps an error with a message.
func WrapError(inner error, message string) error ***REMOVED***
	return &wrappedError***REMOVED***message, inner***REMOVED***
***REMOVED***

// WrapErrorf wraps an error with a message.
func WrapErrorf(inner error, format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return &wrappedError***REMOVED***fmt.Sprintf(format, args...), inner***REMOVED***
***REMOVED***

// MultiError combines multiple errors into a single error. If there are no errors,
// nil is returned. If there is 1 error, it is returned. Otherwise, they are combined.
func MultiError(errors ...error) error ***REMOVED***

	// remove nils from the error list
	var nonNils []error
	for _, e := range errors ***REMOVED***
		if e != nil ***REMOVED***
			nonNils = append(nonNils, e)
		***REMOVED***
	***REMOVED***

	switch len(nonNils) ***REMOVED***
	case 0:
		return nil
	case 1:
		return nonNils[0]
	default:
		return &multiError***REMOVED***
			message: "multiple errors encountered",
			errors:  nonNils,
		***REMOVED***
	***REMOVED***
***REMOVED***

type multiError struct ***REMOVED***
	message string
	errors  []error
***REMOVED***

func (e *multiError) Message() string ***REMOVED***
	return e.message
***REMOVED***

func (e *multiError) Error() string ***REMOVED***
	result := e.message
	for _, e := range e.errors ***REMOVED***
		result += fmt.Sprintf("\n  %s", e)
	***REMOVED***
	return result
***REMOVED***

func (e *multiError) Errors() []error ***REMOVED***
	return e.errors
***REMOVED***

type wrappedError struct ***REMOVED***
	message string
	inner   error
***REMOVED***

func (e *wrappedError) Message() string ***REMOVED***
	return e.message
***REMOVED***

func (e *wrappedError) Error() string ***REMOVED***
	return RolledUpErrorMessage(e)
***REMOVED***

func (e *wrappedError) Inner() error ***REMOVED***
	return e.inner
***REMOVED***

func (e *wrappedError) Unwrap() error ***REMOVED***
	return e.inner
***REMOVED***
