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
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/stats"
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
	var testCases = map[string]string***REMOVED***
		"module.exports": `
			var globalVar;
			if (!globalVar) ***REMOVED***
				globalVar = Math.random();
			***REMOVED***
			function C() ***REMOVED***
				return globalVar;
			***REMOVED***
			module.exports = ***REMOVED***
				C: C,
			***REMOVED***
		`,
		"direct export": `

			var globalVar;
			if (!globalVar) ***REMOVED***
				globalVar = Math.random();
			***REMOVED***
			export function C() ***REMOVED***
				return globalVar;
			***REMOVED***
		`,
	***REMOVED***
	for name, data := range testCases ***REMOVED***
		cData := data
		t.Run(name, func(t *testing.T) ***REMOVED***

			fs := afero.NewMemMapFs()
			require.NoError(t, afero.WriteFile(fs, "/C.js", []byte(cData), os.ModePerm))

			require.NoError(t, afero.WriteFile(fs, "/A.js", []byte(`
		import ***REMOVED*** C ***REMOVED*** from "./C.js";
		export function A() ***REMOVED***
			return C();
		***REMOVED***
	`), os.ModePerm))
			require.NoError(t, afero.WriteFile(fs, "/B.js", []byte(`
		var  c = require("./C.js");
		export function B() ***REMOVED***
			return c.C();
		***REMOVED***
	`), os.ModePerm))
			r1, err := getSimpleRunner("/script.js", `
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
		`, fs)
			require.NoError(t, err)

			arc := r1.MakeArchive()
			r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
			require.NoError(t, err)

			runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
			for name, r := range runners ***REMOVED***
				r := r
				t.Run(name, func(t *testing.T) ***REMOVED***
					ch := newDevNullSampleChannel()
					defer close(ch)
					initVU, err := r.NewVU(1, ch)

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
					require.NoError(t, err)
					err = vu.RunOnce()
					require.NoError(t, err)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLoadExportsIsUsableInModule(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/A.js", []byte(`
		export function A() ***REMOVED***
			return "A";
		***REMOVED***
		export function B() ***REMOVED***
			return exports.A() + "B";
		***REMOVED***
	`), os.ModePerm))
	r1, err := getSimpleRunner("/script.js", `
			import ***REMOVED*** A, B ***REMOVED*** from "./A.js";

			export default function(data) ***REMOVED***
				if (A() != "A") ***REMOVED***
					throw new Error("wrong value of A() " + A());
				***REMOVED***

				if (B() != "AB") ***REMOVED***
					throw new Error("wrong value of B() " + B());
				***REMOVED***
			***REMOVED***
		`, fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			initVU, err := r.NewVU(1, ch)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			require.NoError(t, err)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestLoadDoesntBreakHTTPGet(t *testing.T) ***REMOVED***
	// This test that functions such as http.get which require context still work if they are called
	// inside script that is imported

	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/A.js", []byte(tb.Replacer.Replace(`
		import http from "k6/http";
		export function A() ***REMOVED***
			return http.get("HTTPBIN_URL/get");
		***REMOVED***
	`)), os.ModePerm))
	r1, err := getSimpleRunner("/script.js", `
			import ***REMOVED*** A ***REMOVED*** from "./A.js";

			export default function(data) ***REMOVED***
				let resp = A();
				if (resp.status != 200) ***REMOVED***
					throw new Error("wrong status "+ resp.status);
				***REMOVED***
			***REMOVED***
		`, fs)
	require.NoError(t, err)

	require.NoError(t, r1.SetOptions(lib.Options***REMOVED***Hosts: tb.Dialer.Hosts***REMOVED***))
	arc := r1.MakeArchive()
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			initVU, err := r.NewVU(1, ch)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
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
	r1, err := getSimpleRunner("/script.js", `
			import ***REMOVED*** A ***REMOVED*** from "./A.js";

			export default function(data) ***REMOVED***
				var a = A();
				if (a == 1) ***REMOVED***
					a = 2;
				***REMOVED*** else ***REMOVED***
					throw new Error("wrong value of a " + a);
				***REMOVED***
			***REMOVED***
		`, fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			initVU, err := r.NewVU(1, ch)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)

			// run a second VU
			initVU, err = r.NewVU(2, ch)
			require.NoError(t, err)
			ctx, cancel = context.WithCancel(context.Background())
			defer cancel()
			vu = initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
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
	r1, err := getSimpleRunner("/main.js", string(data), fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			initVU, err := r.NewVU(1, ch)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
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

	r1, err := getSimpleRunner("/main.js", `
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
		`, fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			initVU, err := r.NewVU(1, ch)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBrowserified(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	//nolint: lll
	require.NoError(t, afero.WriteFile(fs, "/browserified.js", []byte(`
		(function(f)***REMOVED***if(typeof exports==="object"&&typeof module!=="undefined")***REMOVED***module.exports=f()***REMOVED***else if(typeof define==="function"&&define.amd)***REMOVED***define([],f)***REMOVED***else***REMOVED***var g;if(typeof window!=="undefined")***REMOVED***g=window***REMOVED***else if(typeof global!=="undefined")***REMOVED***g=global***REMOVED***else if(typeof self!=="undefined")***REMOVED***g=self***REMOVED***else***REMOVED***g=this***REMOVED***g.npmlibs = f()***REMOVED******REMOVED***)(function()***REMOVED***var define,module,exports;return (function()***REMOVED***function r(e,n,t)***REMOVED***function o(i,f)***REMOVED***if(!n[i])***REMOVED***if(!e[i])***REMOVED***var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a***REMOVED***var p=n[i]=***REMOVED***exports:***REMOVED******REMOVED******REMOVED***;e[i][0].call(p.exports,function(r)***REMOVED***var n=e[i][1][r];return o(n||r)***REMOVED***,p,p.exports,r,e,n,t)***REMOVED***return n[i].exports***REMOVED***for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o***REMOVED***return r***REMOVED***)()(***REMOVED***1:[function(require,module,exports)***REMOVED***
		module.exports.A = function () ***REMOVED***
			return "a";
		***REMOVED***

		***REMOVED***,***REMOVED******REMOVED***],2:[function(require,module,exports)***REMOVED***
		exports.B = function() ***REMOVED***
		return "b";
		***REMOVED***

		***REMOVED***,***REMOVED******REMOVED***],3:[function(require,module,exports)***REMOVED***
		exports.alpha = require('./a.js');
		exports.bravo = require('./b.js');

		***REMOVED***,***REMOVED***"./a.js":1,"./b.js":2***REMOVED***]***REMOVED***,***REMOVED******REMOVED***,[3])(3)
		***REMOVED***);
	`), os.ModePerm))

	r1, err := getSimpleRunner("/script.js", `
			import ***REMOVED***alpha, bravo ***REMOVED*** from "./browserified.js";

			export default function(data) ***REMOVED***
				if (alpha.A === undefined) ***REMOVED***
					throw new Error("alpha.A is undefined");
				***REMOVED***
				if (alpha.A() != "a") ***REMOVED***
					throw new Error("alpha.A() != 'a'    (" + alpha.A() + ") != 'a'");
				***REMOVED***

				if (bravo.B === undefined) ***REMOVED***
					throw new Error("bravo.B is undefined");
				***REMOVED***
				if (bravo.B() != "b") ***REMOVED***
					throw new Error("bravo.B() != 'b'    (" + bravo.B() + ") != 'b'");
				***REMOVED***
			***REMOVED***
		`, fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := make(chan stats.SampleContainer, 100)
			defer close(ch)
			initVU, err := r.NewVU(1, ch)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***
func TestLoadingUnexistingModuleDoesntPanic(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	data := `var b;
			try ***REMOVED***
				b = eval("require('buffer')");
			***REMOVED*** catch (err) ***REMOVED***
				b = "correct";
			***REMOVED***
			export default function() ***REMOVED***
				if (b != "correct") ***REMOVED***
					throw new Error("wrong b "+ JSON.stringify(b));
				***REMOVED***
			***REMOVED***`
	require.NoError(t, afero.WriteFile(fs, "/script.js", []byte(data), 0644))
	r1, err := getSimpleRunner("/script.js", data, fs)
	require.NoError(t, err)

	arc := r1.MakeArchive()
	var buf = &bytes.Buffer***REMOVED******REMOVED***
	require.NoError(t, arc.Write(buf))
	arc, err = lib.ReadArchive(buf)
	require.NoError(t, err)
	r2, err := NewFromArchive(arc, lib.RuntimeOptions***REMOVED******REMOVED***)
	require.NoError(t, err)

	runners := map[string]*Runner***REMOVED***"Source": r1, "Archive": r2***REMOVED***
	for name, r := range runners ***REMOVED***
		r := r
		t.Run(name, func(t *testing.T) ***REMOVED***
			ch := newDevNullSampleChannel()
			defer close(ch)
			initVU, err := r.NewVU(1, ch)
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: ctx***REMOVED***)
			err = vu.RunOnce()
			require.NoError(t, err)
		***REMOVED***)
	***REMOVED***
***REMOVED***
