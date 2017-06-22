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

package js

import (
	"crypto/tls"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestNewBundle(t *testing.T) ***REMOVED***
	t.Run("Blank", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(``),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte***REMOVED***0x00***REMOVED***,
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "Transform: SyntaxError: /script.js: Unexpected character '\x00' (1:0)\n> 1 | \x00\n    | ^ at <eval>:2:26853(114)")
	***REMOVED***)
	t.Run("Error", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(`throw new Error("aaaa");`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "Error: aaaa at /script.js:1:20(3)")
	***REMOVED***)
	t.Run("InvalidExports", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(`exports = null`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "exports must be an object")
	***REMOVED***)
	t.Run("DefaultUndefined", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default undefined;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultNull", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default null;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "script must export a default function")
	***REMOVED***)
	t.Run("DefaultWrongType", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data: []byte(`
				export default 12345;
			`),
		***REMOVED***, afero.NewMemMapFs())
		assert.EqualError(t, err, "default export must be a function")
	***REMOVED***)
	t.Run("Minimal", func(t *testing.T) ***REMOVED***
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
		***REMOVED***, afero.NewMemMapFs())
		assert.NoError(t, err)
	***REMOVED***)
	t.Run("stdin", func(t *testing.T) ***REMOVED***
		b, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "-",
			Data:     []byte(`export default function() ***REMOVED******REMOVED***;`),
		***REMOVED***, afero.NewMemMapFs())
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "-", b.Filename)
			assert.Equal(t, "/", b.BaseInitContext.pwd)
		***REMOVED***
	***REMOVED***)
	t.Run("Options", func(t *testing.T) ***REMOVED***
		t.Run("Empty", func(t *testing.T) ***REMOVED***
			_, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED******REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			invalidOptions := map[string]struct ***REMOVED***
				Expr, Error string
			***REMOVED******REMOVED***
				"Array":    ***REMOVED***`[]`, "json: cannot unmarshal array into Go value of type lib.Options"***REMOVED***,
				"Function": ***REMOVED***`function()***REMOVED******REMOVED***`, "json: unsupported type: func(goja.FunctionCall) goja.Value"***REMOVED***,
			***REMOVED***
			for name, data := range invalidOptions ***REMOVED***
				t.Run(name, func(t *testing.T) ***REMOVED***
					_, err := NewBundle(&lib.SourceData***REMOVED***
						Filename: "/script.js",
						Data: []byte(fmt.Sprintf(`
							export let options = %s;
							export default function() ***REMOVED******REMOVED***;
						`, data.Expr)),
					***REMOVED***, afero.NewMemMapFs())
					assert.EqualError(t, err, data.Error)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)

		t.Run("Paused", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						paused: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.Paused)
			***REMOVED***
		***REMOVED***)
		t.Run("VUs", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						vus: 100,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.VUs)
			***REMOVED***
		***REMOVED***)
		t.Run("VUsMax", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						vusMax: 100,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.VUsMax)
			***REMOVED***
		***REMOVED***)
		t.Run("Duration", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						duration: "10s",
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, lib.NullDurationFrom(10*time.Second), b.Options.Duration)
			***REMOVED***
		***REMOVED***)
		t.Run("Iterations", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						iterations: 100,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.Iterations)
			***REMOVED***
		***REMOVED***)
		t.Run("Stages", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						stages: [],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Len(t, b.Options.Stages, 0)
			***REMOVED***

			t.Run("Empty", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED******REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED******REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Target", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***target: 10***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Duration", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***duration: "10s"***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("DurationAndTarget", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***duration: "10s", target: 10***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("RampUpAndPlateau", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							stages: [
								***REMOVED***duration: "10s", target: 10***REMOVED***,
								***REMOVED***duration: "5s"***REMOVED***,
							],
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 2) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
						assert.Equal(t, lib.Stage***REMOVED***Duration: lib.NullDurationFrom(5 * time.Second)***REMOVED***, b.Options.Stages[1])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
		t.Run("Linger", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						linger: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.Linger)
			***REMOVED***
		***REMOVED***)
		t.Run("NoUsageReport", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						noUsageReport: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.NoUsageReport)
			***REMOVED***
		***REMOVED***)
		t.Run("MaxRedirects", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						maxRedirects: 10,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(10), b.Options.MaxRedirects)
			***REMOVED***
		***REMOVED***)
		t.Run("InsecureSkipTLSVerify", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						insecureSkipTLSVerify: true,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.InsecureSkipTLSVerify)
			***REMOVED***
		***REMOVED***)
		t.Run("TLSCipherSuites", func(t *testing.T) ***REMOVED***
			for suiteName, suiteID := range lib.SupportedTLSCipherSuites ***REMOVED***
				t.Run(suiteName, func(t *testing.T) ***REMOVED***
					script := `
					export let options = ***REMOVED***
						tlsCipherSuites: ["%s"]
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
					`
					script = fmt.Sprintf(script, suiteName)

					b, err := NewBundle(&lib.SourceData***REMOVED***
						Filename: "/script.js",
						Data:     []byte(script),
					***REMOVED***, afero.NewMemMapFs())
					if assert.NoError(t, err) ***REMOVED***
						if assert.Len(t, *b.Options.TLSCipherSuites, 1) ***REMOVED***
							assert.Equal(t, (*b.Options.TLSCipherSuites)[0], suiteID)
						***REMOVED***
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
		t.Run("TLSVersion", func(t *testing.T) ***REMOVED***
			t.Run("Object", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
						export let options = ***REMOVED***
							tlsVersion: ***REMOVED***
								min: "ssl3.0",
								max: "tls1.2"
							***REMOVED***
						***REMOVED***;
						export default function() ***REMOVED******REMOVED***;
					`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, b.Options.TLSVersion.Min, tls.VersionSSL30)
					assert.Equal(t, b.Options.TLSVersion.Max, tls.VersionTLS12)
				***REMOVED***
			***REMOVED***)
			t.Run("String", func(t *testing.T) ***REMOVED***
				b, err := NewBundle(&lib.SourceData***REMOVED***
					Filename: "/script.js",
					Data: []byte(`
					export let options = ***REMOVED***
						tlsVersion: "ssl3.0"
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
				***REMOVED***, afero.NewMemMapFs())
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, b.Options.TLSVersion.Min, tls.VersionSSL30)
					assert.Equal(t, b.Options.TLSVersion.Max, tls.VersionSSL30)
				***REMOVED***

			***REMOVED***)
		***REMOVED***)
		t.Run("Thresholds", func(t *testing.T) ***REMOVED***
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
					export let options = ***REMOVED***
						thresholds: ***REMOVED***
							http_req_duration: ["avg<100"],
						***REMOVED***,
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`),
			***REMOVED***, afero.NewMemMapFs())
			if assert.NoError(t, err) ***REMOVED***
				if assert.Len(t, b.Options.Thresholds["http_req_duration"].Thresholds, 1) ***REMOVED***
					assert.Equal(t, "avg<100", b.Options.Thresholds["http_req_duration"].Thresholds[0].Source)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestNewBundleFromArchive(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/path/to", 0755))
	assert.NoError(t, afero.WriteFile(fs, "/path/to/file.txt", []byte(`hi`), 0644))
	assert.NoError(t, afero.WriteFile(fs, "/path/to/exclaim.js", []byte(`export default function(s) ***REMOVED*** return s + "!" ***REMOVED***;`), 0644))

	b, err := NewBundle(&lib.SourceData***REMOVED***
		Filename: "/path/to/script.js",
		Data: []byte(`
			import exclaim from "./exclaim.js";
			export let options = ***REMOVED*** vus: 12345 ***REMOVED***;
			export let file = open("./file.txt");
			export default function() ***REMOVED*** return exclaim(file); ***REMOVED***;
		`),
	***REMOVED***, fs)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, b.Options)

	bi, err := b.Instantiate()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	v, err := bi.Default(goja.Undefined())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, "hi!", v.Export())

	arc := b.MakeArchive()
	assert.Equal(t, "js", arc.Type)
	assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, arc.Options)
	assert.Equal(t, "/path/to/script.js", arc.Filename)
	assert.Equal(t, "\"use strict\";Object.defineProperty(exports, \"__esModule\", ***REMOVED*** value: true ***REMOVED***);exports.file = exports.options = undefined;exports.default =\n\n\n\nfunction () ***REMOVED***return (0, _exclaim2.default)(file);***REMOVED***;var _exclaim = require(\"./exclaim.js\");var _exclaim2 = _interopRequireDefault(_exclaim);function _interopRequireDefault(obj) ***REMOVED***return obj && obj.__esModule ? obj : ***REMOVED*** default: obj ***REMOVED***;***REMOVED***var options = exports.options = ***REMOVED*** vus: 12345 ***REMOVED***;var file = exports.file = open(\"./file.txt\");;", string(arc.Data))
	assert.Equal(t, "/path/to", arc.Pwd)
	assert.Len(t, arc.Scripts, 1)
	assert.Equal(t, "\"use strict\";Object.defineProperty(exports, \"__esModule\", ***REMOVED*** value: true ***REMOVED***);exports.default = function (s) ***REMOVED***return s + \"!\";***REMOVED***;;", string(arc.Scripts["/path/to/exclaim.js"]))
	assert.Len(t, arc.Files, 1)
	assert.Equal(t, `hi`, string(arc.Files["/path/to/file.txt"]))

	b2, err := NewBundleFromArchive(arc)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, b2.Options)

	bi2, err := b.Instantiate()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	v2, err := bi2.Default(goja.Undefined())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	assert.Equal(t, "hi!", v2.Export())
***REMOVED***

func TestBundleInstantiate(t *testing.T) ***REMOVED***
	b, err := NewBundle(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
		let val = true;
		export default function() ***REMOVED*** return val; ***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	bi, err := b.Instantiate()
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	t.Run("Run", func(t *testing.T) ***REMOVED***
		v, err := bi.Default(goja.Undefined())
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, true, v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("SetAndRun", func(t *testing.T) ***REMOVED***
		bi.Runtime.Set("val", false)
		v, err := bi.Default(goja.Undefined())
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, false, v.Export())
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestBundleEnv(t *testing.T) ***REMOVED***
	assert.NoError(t, os.Setenv("TEST_A", "1"))
	assert.NoError(t, os.Setenv("TEST_B", ""))

	b1, err := NewBundle(&lib.SourceData***REMOVED***
		Filename: "/script.js",
		Data: []byte(`
			export default function() ***REMOVED***
				if (__ENV.TEST_A !== "1") ***REMOVED*** throw new Error("Invalid TEST_A: " + __ENV.TEST_A); ***REMOVED***
				if (__ENV.TEST_B !== "") ***REMOVED*** throw new Error("Invalid TEST_B: " + __ENV.TEST_B); ***REMOVED***
			***REMOVED***
		`),
	***REMOVED***, afero.NewMemMapFs())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	b2, err := NewBundleFromArchive(b1.MakeArchive())
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	bundles := map[string]*Bundle***REMOVED***"Source": b1, "Archive": b2***REMOVED***
	for name, b := range bundles ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, "1", b.Env["TEST_A"])
			assert.Equal(t, "", b.Env["TEST_B"])

			bi, err := b.Instantiate()
			if assert.NoError(t, err) ***REMOVED***
				_, err := bi.Default(goja.Undefined())
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
