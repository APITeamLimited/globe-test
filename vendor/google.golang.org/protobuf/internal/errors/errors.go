// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errors implements functions to manipulate errors.
package errors

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/internal/detrand"
)

// Error is a sentinel matching all errors produced by this package.
var Error = errors.New("protobuf error")

// New formats a string according to the format specifier and arguments and
// returns an error that has a "proto" prefix.
func New(f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return &prefixError***REMOVED***s: format(f, x...)***REMOVED***
***REMOVED***

type prefixError struct***REMOVED*** s string ***REMOVED***

var prefix = func() string ***REMOVED***
	// Deliberately introduce instability into the error message string to
	// discourage users from performing error string comparisons.
	if detrand.Bool() ***REMOVED***
		return "proto: " // use non-breaking spaces (U+00a0)
	***REMOVED*** else ***REMOVED***
		return "proto: " // use regular spaces (U+0020)
	***REMOVED***
***REMOVED***()

func (e *prefixError) Error() string ***REMOVED***
	return prefix + e.s
***REMOVED***

func (e *prefixError) Unwrap() error ***REMOVED***
	return Error
***REMOVED***

// Wrap returns an error that has a "proto" prefix, the formatted string described
// by the format specifier and arguments, and a suffix of err. The error wraps err.
func Wrap(err error, f string, x ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return &wrapError***REMOVED***
		s:   format(f, x...),
		err: err,
	***REMOVED***
***REMOVED***

type wrapError struct ***REMOVED***
	s   string
	err error
***REMOVED***

func (e *wrapError) Error() string ***REMOVED***
	return format("%v%v: %v", prefix, e.s, e.err)
***REMOVED***

func (e *wrapError) Unwrap() error ***REMOVED***
	return e.err
***REMOVED***

func (e *wrapError) Is(target error) bool ***REMOVED***
	return target == Error
***REMOVED***

func format(f string, x ...interface***REMOVED******REMOVED***) string ***REMOVED***
	// avoid "proto: " prefix when chaining
	for i := 0; i < len(x); i++ ***REMOVED***
		switch e := x[i].(type) ***REMOVED***
		case *prefixError:
			x[i] = e.s
		case *wrapError:
			x[i] = format("%v: %v", e.s, e.err)
		***REMOVED***
	***REMOVED***
	return fmt.Sprintf(f, x...)
***REMOVED***

func InvalidUTF8(name string) error ***REMOVED***
	return New("field %v contains invalid UTF-8", name)
***REMOVED***

func RequiredNotSet(name string) error ***REMOVED***
	return New("required field %v not set", name)
***REMOVED***
