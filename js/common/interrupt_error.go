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

package common

import (
	"errors"

	"github.com/dop251/goja"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
)

// InterruptError is an error that halts engine execution
type InterruptError struct ***REMOVED***
	Reason string
***REMOVED***

var _ errext.HasExitCode = &InterruptError***REMOVED******REMOVED***

// Error returns the reason of the interruption.
func (i *InterruptError) Error() string ***REMOVED***
	return i.Reason
***REMOVED***

// ExitCode returns the status code used when the k6 process exits.
func (i *InterruptError) ExitCode() errext.ExitCode ***REMOVED***
	return exitcodes.ScriptAborted
***REMOVED***

// AbortTest is the reason emitted when a test script calls test.abort()
const AbortTest = "test aborted"

// IsInterruptError returns true if err is *InterruptError.
func IsInterruptError(err error) bool ***REMOVED***
	if err == nil ***REMOVED***
		return false
	***REMOVED***
	var intErr *InterruptError
	return errors.As(err, &intErr)
***REMOVED***

// UnwrapGojaInterruptedError returns the internal error handled by goja.
func UnwrapGojaInterruptedError(err error) error ***REMOVED***
	var gojaErr *goja.InterruptedError
	if errors.As(err, &gojaErr) ***REMOVED***
		if e, ok := gojaErr.Value().(error); ok ***REMOVED***
			return e
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***
