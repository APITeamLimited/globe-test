// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.13
// +build !go1.13

package errors

import "reflect"

// Is is a copy of Go 1.13's errors.Is for use with older Go versions.
func Is(err, target error) bool ***REMOVED***
	if target == nil ***REMOVED***
		return err == target
	***REMOVED***

	isComparable := reflect.TypeOf(target).Comparable()
	for ***REMOVED***
		if isComparable && err == target ***REMOVED***
			return true
		***REMOVED***
		if x, ok := err.(interface***REMOVED*** Is(error) bool ***REMOVED***); ok && x.Is(target) ***REMOVED***
			return true
		***REMOVED***
		if err = unwrap(err); err == nil ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
***REMOVED***

func unwrap(err error) error ***REMOVED***
	u, ok := err.(interface ***REMOVED***
		Unwrap() error
	***REMOVED***)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	return u.Unwrap()
***REMOVED***
