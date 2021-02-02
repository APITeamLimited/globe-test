/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
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

package data

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/stretchr/testify/require"
)

const makeArrayScript = `
var array = new data.SharedArray("shared",function() ***REMOVED***
    var n = 50;
    var arr = new Array(n);
    for (var i = 0 ; i <n; i++) ***REMOVED***
        arr[i] = ***REMOVED***value: "something" +i***REMOVED***;
    ***REMOVED***
	return arr;
***REMOVED***);
`

func newConfiguredRuntime(initEnv *common.InitEnvironment) (*goja.Runtime, error) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	ctx := common.WithInitEnv(context.Background(), initEnv)
	ctx = common.WithRuntime(ctx, rt)
	rt.Set("data", common.Bind(rt, new(data), &ctx))
	_, err := rt.RunString("var SharedArray = data.SharedArray;")

	return rt, err
***REMOVED***

func TestSharedArrayConstructorExceptions(t *testing.T) ***REMOVED***
	t.Parallel()
	initEnv := &common.InitEnvironment***REMOVED***
		SharedObjects: common.NewSharedObjects(),
	***REMOVED***
	rt, err := newConfiguredRuntime(initEnv)
	require.NoError(t, err)
	cases := map[string]struct ***REMOVED***
		code, err string
	***REMOVED******REMOVED***
		"returning string": ***REMOVED***
			code: `new SharedArray("wat", function() ***REMOVED***return "whatever"***REMOVED***);`,
			err:  "only arrays can be made into SharedArray",
		***REMOVED***,
		"empty name": ***REMOVED***
			code: `new SharedArray("", function() ***REMOVED***return []***REMOVED***);`,
			err:  "empty name provided to SharedArray's constructor",
		***REMOVED***,
		"function in the data": ***REMOVED***
			code: `
			var s = new SharedArray("wat2", function() ***REMOVED***return [***REMOVED***s: function() ***REMOVED******REMOVED******REMOVED***]***REMOVED***);
			if (s[0].s !== undefined) ***REMOVED***
				throw "s[0].s should be undefined"
			***REMOVED***
		`,
			err: "",
		***REMOVED***,
	***REMOVED***

	for name, testCase := range cases ***REMOVED***
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(testCase.code)
			if testCase.err == "" ***REMOVED***
				require.NoError(t, err)
				return // the t.Run
			***REMOVED***

			require.Error(t, err)
			exc := err.(*goja.Exception)
			require.Contains(t, exc.Error(), testCase.err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSharedArrayAnotherRuntimeExceptions(t *testing.T) ***REMOVED***
	t.Parallel()

	initEnv := &common.InitEnvironment***REMOVED***
		SharedObjects: common.NewSharedObjects(),
	***REMOVED***
	rt, err := newConfiguredRuntime(initEnv)
	require.NoError(t, err)
	_, err = rt.RunString(makeArrayScript)
	require.NoError(t, err)

	// create another Runtime with new ctx but keep the initEnv
	rt, err = newConfiguredRuntime(initEnv)
	require.NoError(t, err)
	_, err = rt.RunString(makeArrayScript)
	require.NoError(t, err)

	// use strict is required as otherwise just nothing happens
	cases := map[string]struct ***REMOVED***
		code, err string
	***REMOVED******REMOVED***
		"setting in for-of": ***REMOVED***
			code: `'use strict'; for (var v of array) ***REMOVED*** v.data = "bad"; ***REMOVED***`,
			err:  "Cannot add property data, object is not extensible",
		***REMOVED***,
		"setting from index": ***REMOVED***
			code: `'use strict'; array[2].data2 = "bad2"`,
			err:  "Cannot add property data2, object is not extensible",
		***REMOVED***,
		"setting property on the proxy": ***REMOVED***
			code: `'use strict'; array.something = "something"`,
			err:  "Host object field something cannot be made configurable",
		***REMOVED***,
		"setting index on the proxy": ***REMOVED***
			code: `'use strict'; array[2] = "something"`,
			err:  "Host object field 2 cannot be made configurable",
		***REMOVED***,
	***REMOVED***

	for name, testCase := range cases ***REMOVED***
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(testCase.code)
			if testCase.err == "" ***REMOVED***
				require.NoError(t, err)
				return // the t.Run
			***REMOVED***

			require.Error(t, err)
			exc := err.(*goja.Exception)
			require.Contains(t, exc.Error(), testCase.err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSharedArrayAnotherRuntimeWorking(t *testing.T) ***REMOVED***
	t.Parallel()

	initEnv := &common.InitEnvironment***REMOVED***
		SharedObjects: common.NewSharedObjects(),
	***REMOVED***
	rt, err := newConfiguredRuntime(initEnv)
	require.NoError(t, err)
	_, err = rt.RunString(makeArrayScript)
	require.NoError(t, err)

	// create another Runtime with new ctx but keep the initEnv
	rt, err = newConfiguredRuntime(initEnv)
	require.NoError(t, err)
	_, err = rt.RunString(makeArrayScript)
	require.NoError(t, err)

	_, err = rt.RunString(`
	if (array[2].value !== "something2") ***REMOVED***
		throw new Error("bad array[2]="+array[2].value);
	***REMOVED***
	if (array.length != 50) ***REMOVED***
		throw new Error("bad length " +array.length);
	***REMOVED***

	var i = 0;
	for (var v of array) ***REMOVED***
		if (v.value !== "something"+i) ***REMOVED***
			throw new Error("bad v.value="+v.value+" for i="+i);
		***REMOVED***
		i++;
	***REMOVED***
	`)
	require.NoError(t, err)
***REMOVED***
