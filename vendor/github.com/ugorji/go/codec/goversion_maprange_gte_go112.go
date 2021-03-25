// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build go1.12
// +build safe

package codec

import "reflect"

type mapIter struct ***REMOVED***
	t      *reflect.MapIter
	m      reflect.Value
	values bool
***REMOVED***

func (t *mapIter) ValidKV() (r bool) ***REMOVED***
	return true
***REMOVED***

func (t *mapIter) Next() (r bool) ***REMOVED***
	return t.t.Next()
***REMOVED***

func (t *mapIter) Key() reflect.Value ***REMOVED***
	return t.t.Key()
***REMOVED***

func (t *mapIter) Value() (r reflect.Value) ***REMOVED***
	if t.values ***REMOVED***
		return t.t.Value()
	***REMOVED***
	return
***REMOVED***

func (t *mapIter) Done() ***REMOVED******REMOVED***

func mapRange(t *mapIter, m, k, v reflect.Value, values bool) ***REMOVED***
	*t = mapIter***REMOVED***
		m:      m,
		t:      m.MapRange(),
		values: values,
	***REMOVED***
***REMOVED***
