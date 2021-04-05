// Package errors provides simple error handling primitives.
//
// The traditional error handling idiom in Go is roughly akin to
//
//     if err != nil ***REMOVED***
//             return err
//     ***REMOVED***
//
// which when applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error by recording a stack trace at the point Wrap is called,
// together with the supplied message. For example
//
//     _, err := ioutil.ReadAll(r)
//     if err != nil ***REMOVED***
//             return errors.Wrap(err, "read failed")
//     ***REMOVED***
//
// If additional control is required, the errors.WithStack and
// errors.WithMessage functions destructure errors.Wrap into its component
// operations: annotating an error with a stack trace and with a message,
// respectively.
//
// Retrieving the cause of an error
//
// Using errors.Wrap constructs a stack of errors, adding context to the
// preceding error. Depending on the nature of the error it may be necessary
// to reverse the operation of errors.Wrap to retrieve the original error
// for inspection. Any error value which implements this interface
//
//     type causer interface ***REMOVED***
//             Cause() error
//     ***REMOVED***
//
// can be inspected by errors.Cause. errors.Cause will recursively retrieve
// the topmost error that does not implement causer, which is assumed to be
// the original cause. For example:
//
//     switch err := errors.Cause(err).(type) ***REMOVED***
//     case *MyError:
//             // handle specifically
//     default:
//             // unknown error
//     ***REMOVED***
//
// Although the causer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// Formatted printing of errors
//
// All error values returned from this package implement fmt.Formatter and can
// be formatted by the fmt package. The following verbs are supported:
//
//     %s    print the error. If the error has a Cause it will be
//           printed recursively.
//     %v    see %s
//     %+v   extended format. Each Frame of the error's StackTrace will
//           be printed in detail.
//
// Retrieving the stack trace of an error or wrapper
//
// New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
// invoked. This information can be retrieved with the following interface:
//
//     type stackTracer interface ***REMOVED***
//             StackTrace() errors.StackTrace
//     ***REMOVED***
//
// The returned errors.StackTrace type is defined as
//
//     type StackTrace []Frame
//
// The Frame type represents a call site in the stack trace. Frame supports
// the fmt.Formatter interface that can be used for printing information about
// the stack trace of this error. For example:
//
//     if err, ok := err.(stackTracer); ok ***REMOVED***
//             for _, f := range err.StackTrace() ***REMOVED***
//                     fmt.Printf("%+s:%d", f)
//             ***REMOVED***
//     ***REMOVED***
//
// Although the stackTracer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// See the documentation for Frame.Format for more details.
package errors

import (
	"fmt"
	"io"
)

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) error ***REMOVED***
	return &fundamental***REMOVED***
		msg:   message,
		stack: callers(),
	***REMOVED***
***REMOVED***

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return &fundamental***REMOVED***
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	***REMOVED***
***REMOVED***

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct ***REMOVED***
	msg string
	*stack
***REMOVED***

func (f *fundamental) Error() string ***REMOVED*** return f.msg ***REMOVED***

func (f *fundamental) Format(s fmt.State, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'v':
		if s.Flag('+') ***REMOVED***
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		***REMOVED***
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	***REMOVED***
***REMOVED***

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return &withStack***REMOVED***
		err,
		callers(),
	***REMOVED***
***REMOVED***

type withStack struct ***REMOVED***
	error
	*stack
***REMOVED***

func (w *withStack) Cause() error ***REMOVED*** return w.error ***REMOVED***

func (w *withStack) Format(s fmt.State, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'v':
		if s.Flag('+') ***REMOVED***
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		***REMOVED***
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	***REMOVED***
***REMOVED***

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	err = &withMessage***REMOVED***
		cause: err,
		msg:   message,
	***REMOVED***
	return &withStack***REMOVED***
		err,
		callers(),
	***REMOVED***
***REMOVED***

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	err = &withMessage***REMOVED***
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	***REMOVED***
	return &withStack***REMOVED***
		err,
		callers(),
	***REMOVED***
***REMOVED***

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return &withMessage***REMOVED***
		cause: err,
		msg:   message,
	***REMOVED***
***REMOVED***

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	return &withMessage***REMOVED***
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	***REMOVED***
***REMOVED***

type withMessage struct ***REMOVED***
	cause error
	msg   string
***REMOVED***

func (w *withMessage) Error() string ***REMOVED*** return w.msg + ": " + w.cause.Error() ***REMOVED***
func (w *withMessage) Cause() error  ***REMOVED*** return w.cause ***REMOVED***

func (w *withMessage) Format(s fmt.State, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'v':
		if s.Flag('+') ***REMOVED***
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		***REMOVED***
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	***REMOVED***
***REMOVED***

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface ***REMOVED***
//            Cause() error
//     ***REMOVED***
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error ***REMOVED***
	type causer interface ***REMOVED***
		Cause() error
	***REMOVED***

	for err != nil ***REMOVED***
		cause, ok := err.(causer)
		if !ok ***REMOVED***
			break
		***REMOVED***
		err = cause.Cause()
	***REMOVED***
	return err
***REMOVED***
