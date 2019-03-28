/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package js

import (
	"context"
	"os"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/stats"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func newDevNullSampleChannel() chan stats.SampleContainer ***REMOVED***
	var ch = make(chan stats.SampleContainer, 100)
	go func() ***REMOVED***
		for range ch ***REMOVED***
		***REMOVED***
	***REMOVED***()
	return ch
***REMOVED***

func TestLoadOnceGlobalVars(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/C.js", []byte(`
		var globalVar;
		if (!globalVar) ***REMOVED***
			globalVar = Math.random();
		***REMOVED***
		export function C() ***REMOVED***
			return globalVar;
		***REMOVED***
	`), os.ModePerm))

	require.NoError(t, afero.WriteFile(fs, "/A.js", []byte(`
		import ***REMOVED*** C ***REMOVED*** from "./C.js";
		export function A() ***REMOVED***
			return C();
		***REMOVED***
	`), os.ModePerm))
	require.NoError(t, afero.WriteFile(fs, "/B.js", []byte(`
		import ***REMOVED*** C ***REMOVED*** from "./C.js";
		export function B() ***REMOVED***
			return C();
		***REMOVED***
	`), os.ModePerm))
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			import ***REMOVED*** A ***REMOVED*** from "./A.js";
			import ***REMOVED*** B ***REMOVED*** from "./B.js";

			export default function(data) ***REMOVED***
				if (A() === undefined) ***REMOVED***
					throw new Error("A() is undefined");
				***REMOVED***
				if (A() != B()) ***REMOVED***
					throw new Error("A() != B()    (" + A() + ") != (" + B() + ")");
				***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	arc.Files = make(map[string][]byte)
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			vu, err := r.NewVU(ch)
			require.NoError(t, err)
			err = vu.RunOnce(context.Background())
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLoadDoesntBreakHTTPGet(t *testing.T) ***REMOVED***
	// This test that functions such as http.get which require context still work if they are called
	// inside script that is imported

	tb := testutils.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/A.js", []byte(tb.Replacer.Replace(`
		import http from "k6/http";
		export function A() ***REMOVED***
			return http.get("HTTPBIN_URL/get");
		***REMOVED***
	`)), os.ModePerm))
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			import ***REMOVED*** A ***REMOVED*** from "./A.js";

			export default function(data) ***REMOVED***
				let resp = A();
				if (resp.status != 200) ***REMOVED***
					throw new Error("wrong status "+ resp.status);
				***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Hosts: tb.Dialer.Hosts***REMOVED***))
	arc := r1.MakeArchive()
	arc.Files = make(map[string][]byte)
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			vu, err := r.NewVU(ch)
			require.NoError(t, err)
			err = vu.RunOnce(context.Background())
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLoadGlobalVarsAreNotSharedBetweenVUs(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/A.js", []byte(`
		var globalVar = 0;
		export function A() ***REMOVED***
			globalVar += 1
			return globalVar;
		***REMOVED***
	`), os.ModePerm))
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			import ***REMOVED*** A ***REMOVED*** from "./A.js";

			export default function(data) ***REMOVED***
				var a = A();
				if (a == 1) ***REMOVED***
					a = 2;
				***REMOVED*** else ***REMOVED***
					throw new Error("wrong value of a " + a);
				***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	arc.Files = make(map[string][]byte)
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			vu, err := r.NewVU(ch)
			require.NoError(t, err)
			err = vu.RunOnce(context.Background())
			require.NoError(t, err)

			// run a second VU
			vu, err = r.NewVU(ch)
			require.NoError(t, err)
			err = vu.RunOnce(context.Background())
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLoadCycle(t *testing.T) ***REMOVED***
	// This is mostly the example from https://hacks.mozilla.org/2018/03/es-modules-a-cartoon-deep-dive/
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/counter.js", []byte(`
			let message = require("./main.js").message;
			exports.count = 5;
			export function a() ***REMOVED***
				return message;
			***REMOVED***
	`), os.ModePerm))

	require.NoError(t, afero.WriteFile(fs, "/main.js", []byte(`
			let counter = require("./counter.js");
			let count = counter.count;
			let a = counter.a;
			let message= "Eval complete";
			exports.message = message;

			export default function() ***REMOVED***
				if (count != 5) ***REMOVED***
					throw new Error("Wrong value of count "+ count);
				***REMOVED***
				let aMessage = a();
				if (aMessage != message) ***REMOVED***
					throw new Error("Wrong value of a() "+ aMessage);
				***REMOVED***
			***REMOVED***
	`), os.ModePerm))
	data, err := afero.ReadFile(fs, "/main.js")
	require.NoError(t, err)
	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/main.js",
		Data:     data,
	***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	arc.Files = make(map[string][]byte)
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			vu, err := r.NewVU(ch)
			require.NoError(t, err)
			err = vu.RunOnce(context.Background())
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***

***REMOVED***

func TestLoadCycleBinding(t *testing.T) ***REMOVED***
	// This is mostly the example from
	// http://2ality.com/2015/07/es6-module-exports.html#why-export-bindings
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/a.js", []byte(`
		import ***REMOVED***bar***REMOVED*** from './b.js';
		export function foo(a) ***REMOVED***
				if (a !== undefined) ***REMOVED***
					return "foo" + a;
				***REMOVED***
				return "foo" + bar(3);
		***REMOVED***
	`), os.ModePerm))

	require.NoError(t, afero.WriteFile(fs, "/b.js", []byte(`
		import ***REMOVED***foo***REMOVED*** from './a.js';
		export function bar(a) ***REMOVED***
				if (a !== undefined) ***REMOVED***
					return "bar" + a;
				***REMOVED***
				return "bar" + foo(5);
			***REMOVED***
	`), os.ModePerm))

	r1, err := New(&lib.SourceData***REMOVED***
		Filename: "/main.js",
		Data: []byte(`
			import ***REMOVED***foo***REMOVED*** from './a.js';
			import ***REMOVED***bar***REMOVED*** from './b.js';
			export default function() ***REMOVED***
				let fooMessage = foo();
				if (fooMessage != "foobar3") ***REMOVED***
					throw new Error("Wrong value of foo() "+ fooMessage);
				***REMOVED***
				let barMessage = bar();
				if (barMessage != "barfoo5") ***REMOVED***
					throw new Error("Wrong value of bar() "+ barMessage);
				***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	arc.Files = make(map[string][]byte)
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			vu, err := r.NewVU(ch)
			require.NoError(t, err)
			err = vu.RunOnce(context.Background())
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***
