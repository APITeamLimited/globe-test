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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
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

func newConfiguredRuntime() (*goja.Runtime, error) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	m, ok := New().NewModuleInstance(
		&modulestest.VU***REMOVED***
			RuntimeField: rt,
			InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
			CtxField:     common.WithRuntime(context.Background(), rt),
			StateField:   nil,
		***REMOVED***,
	).(*Data)
	if !ok ***REMOVED***
		return rt, fmt.Errorf("not a Data module instance")
	***REMOVED***

	err := rt.Set("data", m.Exports().Named)
	if err != nil ***REMOVED***
		return rt, err //nolint:wrapcheck
	***REMOVED***
	_, err = rt.RunString("var SharedArray = data.SharedArray;")
	return rt, err //nolint:wrapcheck
***REMOVED***

func TestSharedArrayConstructorExceptions(t *testing.T) ***REMOVED***
	t.Parallel()
	rt, err := newConfiguredRuntime()
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
		"not a function": ***REMOVED***
			code: `var s = new SharedArray("wat3", "astring");`,
			err:  "a function is expected",
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

	rt, err := newConfiguredRuntime()
	require.NoError(t, err)
	_, err = rt.RunString(makeArrayScript)
	require.NoError(t, err)

	rt, err = newConfiguredRuntime()
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
		"setting property on the shared array": ***REMOVED***
			code: `'use strict'; array.something = "something"`,
			err:  `Cannot set property "something" on a dynamic array`,
		***REMOVED***,
		"setting index on the shared array": ***REMOVED***
			code: `'use strict'; array[2] = "something"`,
			err:  "SharedArray is immutable",
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

	rt := goja.New()
	vu := &modulestest.VU***REMOVED***
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment***REMOVED******REMOVED***,
		CtxField:     common.WithRuntime(context.Background(), rt),
		StateField:   nil,
	***REMOVED***
	m, ok := New().NewModuleInstance(vu).(*Data)
	require.True(t, ok)
	require.NoError(t, rt.Set("data", m.Exports().Named))

	_, err := rt.RunString(makeArrayScript)
	require.NoError(t, err)

	// create another Runtime with new ctx but keep the initEnv
	rt = goja.New()
	vu.RuntimeField = rt
	vu.CtxField = common.WithRuntime(context.Background(), rt)
	require.NoError(t, rt.Set("data", m.Exports().Named))

	_, err = rt.RunString(`var array = new data.SharedArray("shared", function() ***REMOVED***throw "wat";***REMOVED***);`)
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

	i = 0;
	array.forEach(function(v)***REMOVED***
		if (v.value !== "something"+i) ***REMOVED***
			throw new Error("bad v.value="+v.value+" for i="+i);
		***REMOVED***
		i++;
	***REMOVED***);


	`)
	require.NoError(t, err)
***REMOVED***

func TestSharedArrayRaceInInitialization(t *testing.T) ***REMOVED***
	t.Parallel()

	const instances = 10
	const repeats = 100
	for i := 0; i < repeats; i++ ***REMOVED***
		runtimes := make([]*goja.Runtime, instances)
		for j := 0; j < instances; j++ ***REMOVED***
			rt, err := newConfiguredRuntime()
			require.NoError(t, err)
			runtimes[j] = rt
		***REMOVED***
		var wg sync.WaitGroup
		for _, rt := range runtimes ***REMOVED***
			rt := rt
			wg.Add(1)
			go func() ***REMOVED***
				defer wg.Done()
				_, err := rt.RunString(`var array = new data.SharedArray("shared", function() ***REMOVED***return [1,2,3,4,5,6,7,8,9, 10]***REMOVED***);`)
				require.NoError(t, err)
			***REMOVED***()
		***REMOVED***
		ch := make(chan struct***REMOVED******REMOVED***)
		go func() ***REMOVED***
			wg.Wait()
			close(ch)
		***REMOVED***()

		select ***REMOVED***
		case <-ch:
			// everything is fine
		case <-time.After(time.Second * 10):
			t.Fatal("Took too long probably locked up")
		***REMOVED***
	***REMOVED***
***REMOVED***
