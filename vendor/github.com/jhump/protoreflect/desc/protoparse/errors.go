package protoparse

import (
	"errors"
	"fmt"
)

// ErrInvalidSource is a sentinel error that is returned by calls to
// Parser.ParseFiles and Parser.ParseFilesButDoNotLink in the event that syntax
// or link errors are encountered, but the parser's configured ErrorReporter
// always returns nil.
var ErrInvalidSource = errors.New("parse failed: invalid proto source")

// ErrNoSyntax is a sentinel error that may be passed to a warning reporter.
// The error the reporter receives will be wrapped with source position that
// indicates the file that had no syntax statement.
var ErrNoSyntax = errors.New("no syntax specified; defaulting to proto2 syntax")

// ErrLookupImportAndProtoSet is the error returned if both LookupImport and LookupImportProto are set.
var ErrLookupImportAndProtoSet = errors.New("both LookupImport and LookupImportProto set")

// ErrorReporter is responsible for reporting the given error. If the reporter
// returns a non-nil error, parsing/linking will abort with that error. If the
// reporter returns nil, parsing will continue, allowing the parser to try to
// report as many syntax and/or link errors as it can find.
type ErrorReporter func(err ErrorWithPos) error

// WarningReporter is responsible for reporting the given warning. This is used
// for indicating non-error messages to the calling program for things that do
// not cause the parse to fail but are considered bad practice. Though they are
// just warnings, the details are supplied to the reporter via an error type.
type WarningReporter func(ErrorWithPos)

func defaultErrorReporter(err ErrorWithPos) error ***REMOVED***
	// abort parsing after first error encountered
	return err
***REMOVED***

type errorHandler struct ***REMOVED***
	errReporter  ErrorReporter
	errsReported int
	err          error

	warnReporter WarningReporter
***REMOVED***

func newErrorHandler(errRep ErrorReporter, warnRep WarningReporter) *errorHandler ***REMOVED***
	if errRep == nil ***REMOVED***
		errRep = defaultErrorReporter
	***REMOVED***
	return &errorHandler***REMOVED***
		errReporter:  errRep,
		warnReporter: warnRep,
	***REMOVED***
***REMOVED***

func (h *errorHandler) handleErrorWithPos(pos *SourcePos, format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	if h.err != nil ***REMOVED***
		return h.err
	***REMOVED***
	h.errsReported++
	err := h.errReporter(errorWithPos(pos, format, args...))
	h.err = err
	return err
***REMOVED***

func (h *errorHandler) handleError(err error) error ***REMOVED***
	if h.err != nil ***REMOVED***
		return h.err
	***REMOVED***
	if ewp, ok := err.(ErrorWithPos); ok ***REMOVED***
		h.errsReported++
		err = h.errReporter(ewp)
	***REMOVED***
	h.err = err
	return err
***REMOVED***

func (h *errorHandler) warn(pos *SourcePos, err error) ***REMOVED***
	if h.warnReporter != nil ***REMOVED***
		h.warnReporter(ErrorWithSourcePos***REMOVED***Pos: pos, Underlying: err***REMOVED***)
	***REMOVED***
***REMOVED***

func (h *errorHandler) getError() error ***REMOVED***
	if h.errsReported > 0 && h.err == nil ***REMOVED***
		return ErrInvalidSource
	***REMOVED***
	return h.err
***REMOVED***

// ErrorWithPos is an error about a proto source file that includes information
// about the location in the file that caused the error.
//
// The value of Error() will contain both the SourcePos and Underlying error.
// The value of Unwrap() will only be the Underlying error.
type ErrorWithPos interface ***REMOVED***
	error
	GetPosition() SourcePos
	Unwrap() error
***REMOVED***

// ErrorWithSourcePos is an error about a proto source file that includes
// information about the location in the file that caused the error.
//
// Errors that include source location information *might* be of this type.
// However, calling code that is trying to examine errors with location info
// should instead look for instances of the ErrorWithPos interface, which
// will find other kinds of errors. This type is only exported for backwards
// compatibility.
//
// SourcePos should always be set and never nil.
type ErrorWithSourcePos struct ***REMOVED***
	Underlying error
	Pos        *SourcePos
***REMOVED***

// Error implements the error interface
func (e ErrorWithSourcePos) Error() string ***REMOVED***
	sourcePos := e.GetPosition()
	return fmt.Sprintf("%s: %v", sourcePos, e.Underlying)
***REMOVED***

// GetPosition implements the ErrorWithPos interface, supplying a location in
// proto source that caused the error.
func (e ErrorWithSourcePos) GetPosition() SourcePos ***REMOVED***
	if e.Pos == nil ***REMOVED***
		return SourcePos***REMOVED***Filename: "<input>"***REMOVED***
	***REMOVED***
	return *e.Pos
***REMOVED***

// Unwrap implements the ErrorWithPos interface, supplying the underlying
// error. This error will not include location information.
func (e ErrorWithSourcePos) Unwrap() error ***REMOVED***
	return e.Underlying
***REMOVED***

var _ ErrorWithPos = ErrorWithSourcePos***REMOVED******REMOVED***

func errorWithPos(pos *SourcePos, format string, args ...interface***REMOVED******REMOVED***) ErrorWithPos ***REMOVED***
	return ErrorWithSourcePos***REMOVED***Pos: pos, Underlying: fmt.Errorf(format, args...)***REMOVED***
***REMOVED***

type errorWithFilename struct ***REMOVED***
	underlying error
	filename   string
***REMOVED***

func (e errorWithFilename) Error() string ***REMOVED***
	return fmt.Sprintf("%s: %v", e.filename, e.underlying)
***REMOVED***

func (e errorWithFilename) Unwrap() error ***REMOVED***
	return e.underlying
***REMOVED***

// ErrorUnusedImport may be passed to a warning reporter when an unused
// import is detected. The error the reporter receives will be wrapped
// with source position that indicates the file and line where the import
// statement appeared.
type ErrorUnusedImport interface ***REMOVED***
	error
	UnusedImport() string
***REMOVED***

type errUnusedImport string

func (e errUnusedImport) Error() string ***REMOVED***
	return fmt.Sprintf("import %q not used", string(e))
***REMOVED***

func (e errUnusedImport) UnusedImport() string ***REMOVED***
	return string(e)
***REMOVED***
