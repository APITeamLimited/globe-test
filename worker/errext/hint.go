package errext

import "errors"

// HasHint is a wrapper around an error with an attached user hint. These hints
// can be used to give extra human-readable information about the error,
// including suggestions on how the error can be fixed.
type HasHint interface ***REMOVED***
	error
	Hint() string
***REMOVED***

// WithHint is a helper that can attach a hint to the given error. If there is
// no error (i.e. the given error is nil), it won't do anything. If the given
// error already had a hint, this helper will wrap it so that the new hint is
// "new hint (old hint)".
func WithHint(err error, hint string) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil // No error, do nothing
	***REMOVED***
	return withHint***REMOVED***err, hint***REMOVED***
***REMOVED***

type withHint struct ***REMOVED***
	error
	hint string
***REMOVED***

func (wh withHint) Unwrap() error ***REMOVED***
	return wh.error
***REMOVED***

func (wh withHint) Hint() string ***REMOVED***
	hint := wh.hint
	var oldhint HasHint
	if errors.As(wh.error, &oldhint) ***REMOVED***
		// The given error already had a hint, wrap it
		hint = hint + " (" + oldhint.Hint() + ")"
	***REMOVED***

	return hint
***REMOVED***

var _ HasHint = withHint***REMOVED******REMOVED***
