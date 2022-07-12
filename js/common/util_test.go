/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestThrow(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	fn1, ok := goja.AssertFunction(rt.ToValue(func() ***REMOVED*** Throw(rt, errors.New("aaaa")) ***REMOVED***))
	if assert.True(t, ok, "fn1 is invalid") ***REMOVED***
		_, err := fn1(goja.Undefined())
		assert.EqualError(t, err, "aaaa")

		fn2, ok := goja.AssertFunction(rt.ToValue(func() ***REMOVED*** Throw(rt, err) ***REMOVED***))
		if assert.True(t, ok, "fn1 is invalid") ***REMOVED***
			_, err := fn2(goja.Undefined())
			assert.EqualError(t, err, "aaaa")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestToBytes(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	b := []byte("hello")
	testCases := []struct ***REMOVED***
		in     interface***REMOVED******REMOVED***
		expOut []byte
		expErr string
	***REMOVED******REMOVED***
		***REMOVED***b, b, ""***REMOVED***,
		***REMOVED***"hello", b, ""***REMOVED***,
		***REMOVED***rt.NewArrayBuffer(b), b, ""***REMOVED***,
		***REMOVED***struct***REMOVED******REMOVED******REMOVED******REMOVED***, nil, "invalid type struct ***REMOVED******REMOVED***, expected string, []byte or ArrayBuffer"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%T", tc.in), func(t *testing.T) ***REMOVED***
			t.Parallel()
			out, err := ToBytes(tc.in)
			if tc.expErr != "" ***REMOVED***
				assert.EqualError(t, err, tc.expErr)
				return
			***REMOVED***
			assert.Equal(t, tc.expOut, out)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestToString(t *testing.T) ***REMOVED***
	t.Parallel()
	rt := goja.New()
	s := "hello"
	testCases := []struct ***REMOVED***
		in             interface***REMOVED******REMOVED***
		expOut, expErr string
	***REMOVED******REMOVED***
		***REMOVED***s, s, ""***REMOVED***,
		***REMOVED***"hello", s, ""***REMOVED***,
		***REMOVED***rt.NewArrayBuffer([]byte(s)), s, ""***REMOVED***,
		***REMOVED***struct***REMOVED******REMOVED******REMOVED******REMOVED***, "", "invalid type struct ***REMOVED******REMOVED***, expected string, []byte or ArrayBuffer"***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(fmt.Sprintf("%T", tc.in), func(t *testing.T) ***REMOVED***
			t.Parallel()
			out, err := ToString(tc.in)
			if tc.expErr != "" ***REMOVED***
				assert.EqualError(t, err, tc.expErr)
				return
			***REMOVED***
			assert.Equal(t, tc.expOut, out)
		***REMOVED***)
	***REMOVED***
***REMOVED***
