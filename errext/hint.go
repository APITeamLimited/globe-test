/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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
		// No error, do nothing
		return nil
	***REMOVED***
	var oldhint HasHint
	if errors.As(err, &oldhint) ***REMOVED***
		// The given error already had a hint, wrap it
		hint = hint + " (" + oldhint.Hint() + ")"
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
	return wh.hint
***REMOVED***

var _ HasHint = withHint***REMOVED******REMOVED***
