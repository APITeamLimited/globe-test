// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.12

package impl

import "reflect"

type mapIter struct ***REMOVED***
	v    reflect.Value
	keys []reflect.Value
***REMOVED***

// mapRange provides a less-efficient equivalent to
// the Go 1.12 reflect.Value.MapRange method.
func mapRange(v reflect.Value) *mapIter ***REMOVED***
	return &mapIter***REMOVED***v: v***REMOVED***
***REMOVED***

func (i *mapIter) Next() bool ***REMOVED***
	if i.keys == nil ***REMOVED***
		i.keys = i.v.MapKeys()
	***REMOVED*** else ***REMOVED***
		i.keys = i.keys[1:]
	***REMOVED***
	return len(i.keys) > 0
***REMOVED***

func (i *mapIter) Key() reflect.Value ***REMOVED***
	return i.keys[0]
***REMOVED***

func (i *mapIter) Value() reflect.Value ***REMOVED***
	return i.v.MapIndex(i.keys[0])
***REMOVED***
