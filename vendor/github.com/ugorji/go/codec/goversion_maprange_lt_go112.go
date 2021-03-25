// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build !go1.12
// +build !go1.7 safe

package codec

import "reflect"

type mapIter struct ***REMOVED***
	m      reflect.Value
	keys   []reflect.Value
	j      int
	values bool
***REMOVED***

func (t *mapIter) ValidKV() (r bool) ***REMOVED***
	return true
***REMOVED***

func (t *mapIter) Next() (r bool) ***REMOVED***
	t.j++
	return t.j < len(t.keys)
***REMOVED***

func (t *mapIter) Key() reflect.Value ***REMOVED***
	return t.keys[t.j]
***REMOVED***

func (t *mapIter) Value() (r reflect.Value) ***REMOVED***
	if t.values ***REMOVED***
		return t.m.MapIndex(t.keys[t.j])
	***REMOVED***
	return
***REMOVED***

func (t *mapIter) Done() ***REMOVED******REMOVED***

func mapRange(t *mapIter, m, k, v reflect.Value, values bool) ***REMOVED***
	*t = mapIter***REMOVED***
		m:      m,
		keys:   m.MapKeys(),
		values: values,
		j:      -1,
	***REMOVED***
***REMOVED***
