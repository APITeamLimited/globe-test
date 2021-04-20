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
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/consts"
	"github.com/loadimpact/k6/lib/fsext"
	"github.com/loadimpact/k6/lib/testutils"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
)

const isWindows = runtime.GOOS == "windows"

func getSimpleBundle(tb testing.TB, filename, data string, opts ...interface***REMOVED******REMOVED***) (*Bundle, error) ***REMOVED***
	var (
		fs                        = afero.NewMemMapFs()
		rtOpts                    = lib.RuntimeOptions***REMOVED******REMOVED***
		logger logrus.FieldLogger = testutils.NewLogger(tb)
	)
	for _, o := range opts ***REMOVED***
		switch opt := o.(type) ***REMOVED***
		case afero.Fs:
			fs = opt
		case lib.RuntimeOptions:
			rtOpts = opt
		case logrus.FieldLogger:
			logger = opt
		***REMOVED***
	***REMOVED***
	return NewBundle(
		logger,
		&loader.SourceData***REMOVED***
			URL:  &url.URL***REMOVED***Path: filename, Scheme: "file"***REMOVED***,
			Data: []byte(data),
		***REMOVED***,
		map[string]afero.Fs***REMOVED***"file": fs, "https": afero.NewMemMapFs()***REMOVED***,
		rtOpts,
	)
***REMOVED***

func TestNewBundle(t *testing.T) ***REMOVED***
	t.Run("Blank", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", "")
		assert.EqualError(t, err, "no exported functions in script")
	***REMOVED***)
	t.Run("Invalid", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", "\x00")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "SyntaxError: file:///script.js: Unexpected character '\x00' (1:0)\n> 1 | \x00\n")
	***REMOVED***)
	t.Run("Error", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", `throw new Error("aaaa");`)
		exception := new(scriptException)
		assert.ErrorAs(t, err, &exception)
		assert.EqualError(t, err, "Error: aaaa\n\tat file:///script.js:1:7(3)\n")
	***REMOVED***)
	t.Run("InvalidExports", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", `exports = null`)
		assert.EqualError(t, err, "exports must be an object")
	***REMOVED***)
	t.Run("DefaultUndefined", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", `export default undefined;`)
		assert.EqualError(t, err, "no exported functions in script")
	***REMOVED***)
	t.Run("DefaultNull", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", `export default null;`)
		assert.EqualError(t, err, "no exported functions in script")
	***REMOVED***)
	t.Run("DefaultWrongType", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", `export default 12345;`)
		assert.EqualError(t, err, "no exported functions in script")
	***REMOVED***)
	t.Run("Minimal", func(t *testing.T) ***REMOVED***
		_, err := getSimpleBundle(t, "/script.js", `export default function() ***REMOVED******REMOVED***;`)
		assert.NoError(t, err)
	***REMOVED***)
	t.Run("stdin", func(t *testing.T) ***REMOVED***
		b, err := getSimpleBundle(t, "-", `export default function() ***REMOVED******REMOVED***;`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "file://-", b.Filename.String())
			assert.Equal(t, "file:///", b.BaseInitContext.pwd.String())
		***REMOVED***
	***REMOVED***)
	t.Run("CompatibilityMode", func(t *testing.T) ***REMOVED***
		t.Run("Extended/ok/global", func(t *testing.T) ***REMOVED***
			rtOpts := lib.RuntimeOptions***REMOVED***
				CompatibilityMode: null.StringFrom(lib.CompatibilityModeExtended.String()),
			***REMOVED***
			_, err := getSimpleBundle(t, "/script.js",
				`module.exports.default = function() ***REMOVED******REMOVED***
				if (global.Math != Math) ***REMOVED***
					throw new Error("global is not defined");
				***REMOVED***`, rtOpts)

			assert.NoError(t, err)
		***REMOVED***)
		t.Run("Base/ok/Minimal", func(t *testing.T) ***REMOVED***
			rtOpts := lib.RuntimeOptions***REMOVED***
				CompatibilityMode: null.StringFrom(lib.CompatibilityModeBase.String()),
			***REMOVED***
			_, err := getSimpleBundle(t, "/script.js",
				`module.exports.default = function() ***REMOVED******REMOVED***;`, rtOpts)
			assert.NoError(t, err)
		***REMOVED***)
		t.Run("Base/err", func(t *testing.T) ***REMOVED***
			testCases := []struct ***REMOVED***
				name       string
				compatMode string
				code       string
				expErr     string
			***REMOVED******REMOVED***
				***REMOVED***
					"InvalidCompat", "es1", `export default function() ***REMOVED******REMOVED***;`,
					`invalid compatibility mode "es1". Use: "extended", "base"`,
				***REMOVED***,
				// ES2015 modules are not supported
				***REMOVED***
					"Modules", "base", `export default function() ***REMOVED******REMOVED***;`,
					"file:///script.js: Line 1:1 Unexpected reserved word",
				***REMOVED***,
				// Arrow functions are not supported
				***REMOVED***
					"ArrowFuncs", "base",
					`module.exports.default = function() ***REMOVED******REMOVED***; () => ***REMOVED******REMOVED***;`,
					"file:///script.js: Line 1:42 Unexpected token ) (and 1 more errors)",
				***REMOVED***,
				// Promises are not supported
				***REMOVED***
					"Promise", "base",
					`module.exports.default = function() ***REMOVED******REMOVED***; new Promise(function(resolve, reject)***REMOVED******REMOVED***);`,
					"ReferenceError: Promise is not defined\n\tat file:///script.js:1:45(4)\n",
				***REMOVED***,
			***REMOVED***

			for _, tc := range testCases ***REMOVED***
				tc := tc
				t.Run(tc.name, func(t *testing.T) ***REMOVED***
					rtOpts := lib.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom(tc.compatMode)***REMOVED***
					_, err := getSimpleBundle(t, "/script.js", tc.code, rtOpts)
					assert.EqualError(t, err, tc.expErr)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Options", func(t *testing.T) ***REMOVED***
		t.Run("Empty", func(t *testing.T) ***REMOVED***
			_, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED******REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
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
					_, err := getSimpleBundle(t, "/script.js", fmt.Sprintf(`
						export let options = %s;
						export default function() ***REMOVED******REMOVED***;
					`, data.Expr))
					assert.EqualError(t, err, data.Error)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)

		t.Run("Paused", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					paused: true,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.BoolFrom(true), b.Options.Paused)
			***REMOVED***
		***REMOVED***)
		t.Run("VUs", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					vus: 100,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.VUs)
			***REMOVED***
		***REMOVED***)
		t.Run("Duration", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					duration: "10s",
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, types.NullDurationFrom(10*time.Second), b.Options.Duration)
			***REMOVED***
		***REMOVED***)
		t.Run("Iterations", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					iterations: 100,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(100), b.Options.Iterations)
			***REMOVED***
		***REMOVED***)
		t.Run("Stages", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					stages: [],
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Len(t, b.Options.Stages, 0)
			***REMOVED***

			t.Run("Empty", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle(t, "/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED******REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED******REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Target", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle(t, "/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***target: 10***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("Duration", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle(t, "/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***duration: "10s"***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("DurationAndTarget", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle(t, "/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***duration: "10s", target: 10***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 1) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
			t.Run("RampUpAndPlateau", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle(t, "/script.js", `
					export let options = ***REMOVED***
						stages: [
							***REMOVED***duration: "10s", target: 10***REMOVED***,
							***REMOVED***duration: "5s"***REMOVED***,
						],
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					if assert.Len(t, b.Options.Stages, 2) ***REMOVED***
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(10 * time.Second), Target: null.IntFrom(10)***REMOVED***, b.Options.Stages[0])
						assert.Equal(t, lib.Stage***REMOVED***Duration: types.NullDurationFrom(5 * time.Second)***REMOVED***, b.Options.Stages[1])
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
		t.Run("MaxRedirects", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					maxRedirects: 10,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, null.IntFrom(10), b.Options.MaxRedirects)
			***REMOVED***
		***REMOVED***)
		t.Run("InsecureSkipTLSVerify", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					insecureSkipTLSVerify: true,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
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

					b, err := getSimpleBundle(t, "/script.js", script)
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
				b, err := getSimpleBundle(t, "/script.js", `
					export let options = ***REMOVED***
						tlsVersion: ***REMOVED***
							min: "tls1.0",
							max: "tls1.2"
						***REMOVED***
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, b.Options.TLSVersion.Min, lib.TLSVersion(tls.VersionTLS10))
					assert.Equal(t, b.Options.TLSVersion.Max, lib.TLSVersion(tls.VersionTLS12))
				***REMOVED***
			***REMOVED***)
			t.Run("String", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle(t, "/script.js", `
					export let options = ***REMOVED***
						tlsVersion: "tls1.0"
					***REMOVED***;
					export default function() ***REMOVED******REMOVED***;
				`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, b.Options.TLSVersion.Min, lib.TLSVersion(tls.VersionTLS10))
					assert.Equal(t, b.Options.TLSVersion.Max, lib.TLSVersion(tls.VersionTLS10))
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
		t.Run("Thresholds", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					thresholds: ***REMOVED***
						http_req_duration: ["avg<100"],
					***REMOVED***,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`)
			if assert.NoError(t, err) ***REMOVED***
				if assert.Len(t, b.Options.Thresholds["http_req_duration"].Thresholds, 1) ***REMOVED***
					assert.Equal(t, "avg<100", b.Options.Thresholds["http_req_duration"].Thresholds[0].Source)
				***REMOVED***
			***REMOVED***
		***REMOVED***)

		t.Run("Unknown field", func(t *testing.T) ***REMOVED***
			logger := logrus.New()
			logger.SetLevel(logrus.InfoLevel)
			logger.Out = ioutil.Discard
			hook := testutils.SimpleLogrusHook***REMOVED***
				HookedLevels: []logrus.Level***REMOVED***logrus.WarnLevel, logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel***REMOVED***,
			***REMOVED***
			logger.AddHook(&hook)

			_, err := getSimpleBundle(t, "/script.js", `
				export let options = ***REMOVED***
					something: ***REMOVED***
						http_req_duration: ["avg<100"],
					***REMOVED***,
				***REMOVED***;
				export default function() ***REMOVED******REMOVED***;
			`, logger)
			require.NoError(t, err)
			entries := hook.Drain()
			require.Len(t, entries, 1)
			assert.Equal(t, logrus.WarnLevel, entries[0].Level)
			require.Contains(t, entries[0].Message, "There were unknown fields")
			require.Contains(t, entries[0].Data["error"].(error).Error(), "unknown field \"something\"")
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func getArchive(tb testing.TB, data string, rtOpts lib.RuntimeOptions) (*lib.Archive, error) ***REMOVED***
	b, err := getSimpleBundle(tb, "script.js", data, rtOpts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return b.makeArchive(), nil
***REMOVED***

func TestNewBundleFromArchive(t *testing.T) ***REMOVED***
	t.Parallel()

	es5Code := `module.exports.options = ***REMOVED*** vus: 12345 ***REMOVED***; module.exports.default = function() ***REMOVED*** return "hi!" ***REMOVED***;`
	es6Code := `export let options = ***REMOVED*** vus: 12345 ***REMOVED***; export default function() ***REMOVED*** return "hi!"; ***REMOVED***;`
	baseCompatModeRtOpts := lib.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom(lib.CompatibilityModeBase.String())***REMOVED***
	extCompatModeRtOpts := lib.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom(lib.CompatibilityModeExtended.String())***REMOVED***

	logger := testutils.NewLogger(t)
	checkBundle := func(t *testing.T, b *Bundle) ***REMOVED***
		assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, b.Options)
		bi, err := b.Instantiate(logger, 0)
		require.NoError(t, err)
		val, err := bi.exports[consts.DefaultFn](goja.Undefined())
		require.NoError(t, err)
		assert.Equal(t, "hi!", val.Export())
	***REMOVED***

	checkArchive := func(t *testing.T, arc *lib.Archive, rtOpts lib.RuntimeOptions, expError string) ***REMOVED***
		b, err := NewBundleFromArchive(logger, arc, rtOpts)
		if expError != "" ***REMOVED***
			require.Error(t, err)
			assert.Contains(t, err.Error(), expError)
		***REMOVED*** else ***REMOVED***
			require.NoError(t, err)
			checkBundle(t, b)
		***REMOVED***
	***REMOVED***

	t.Run("es6_script_default", func(t *testing.T) ***REMOVED***
		t.Parallel()
		arc, err := getArchive(t, es6Code, lib.RuntimeOptions***REMOVED******REMOVED***) // default options
		require.NoError(t, err)
		require.Equal(t, lib.CompatibilityModeExtended.String(), arc.CompatibilityMode)

		checkArchive(t, arc, lib.RuntimeOptions***REMOVED******REMOVED***, "") // default options
		checkArchive(t, arc, extCompatModeRtOpts, "")
		checkArchive(t, arc, baseCompatModeRtOpts, "Unexpected reserved word")
	***REMOVED***)

	t.Run("es6_script_explicit", func(t *testing.T) ***REMOVED***
		t.Parallel()
		arc, err := getArchive(t, es6Code, extCompatModeRtOpts)
		require.NoError(t, err)
		require.Equal(t, lib.CompatibilityModeExtended.String(), arc.CompatibilityMode)

		checkArchive(t, arc, lib.RuntimeOptions***REMOVED******REMOVED***, "")
		checkArchive(t, arc, extCompatModeRtOpts, "")
		checkArchive(t, arc, baseCompatModeRtOpts, "Unexpected reserved word")
	***REMOVED***)

	t.Run("es5_script_with_extended", func(t *testing.T) ***REMOVED***
		t.Parallel()
		arc, err := getArchive(t, es5Code, lib.RuntimeOptions***REMOVED******REMOVED***)
		require.NoError(t, err)
		require.Equal(t, lib.CompatibilityModeExtended.String(), arc.CompatibilityMode)

		checkArchive(t, arc, lib.RuntimeOptions***REMOVED******REMOVED***, "")
		checkArchive(t, arc, extCompatModeRtOpts, "")
		checkArchive(t, arc, baseCompatModeRtOpts, "")
	***REMOVED***)

	t.Run("es5_script", func(t *testing.T) ***REMOVED***
		t.Parallel()
		arc, err := getArchive(t, es5Code, baseCompatModeRtOpts)
		require.NoError(t, err)
		require.Equal(t, lib.CompatibilityModeBase.String(), arc.CompatibilityMode)

		checkArchive(t, arc, lib.RuntimeOptions***REMOVED******REMOVED***, "")
		checkArchive(t, arc, extCompatModeRtOpts, "")
		checkArchive(t, arc, baseCompatModeRtOpts, "")
	***REMOVED***)

	t.Run("es6_archive_with_wrong_compat_mode", func(t *testing.T) ***REMOVED***
		t.Parallel()
		arc, err := getArchive(t, es6Code, baseCompatModeRtOpts)
		require.Error(t, err)
		require.Nil(t, arc)
	***REMOVED***)

	t.Run("messed_up_archive", func(t *testing.T) ***REMOVED***
		t.Parallel()
		arc, err := getArchive(t, es6Code, extCompatModeRtOpts)
		require.NoError(t, err)
		arc.CompatibilityMode = "blah"                                           // intentionally break the archive
		checkArchive(t, arc, lib.RuntimeOptions***REMOVED******REMOVED***, "invalid compatibility mode") // fails when it uses the archive one
		checkArchive(t, arc, extCompatModeRtOpts, "")                            // works when I force the compat mode
		checkArchive(t, arc, baseCompatModeRtOpts, "Unexpected reserved word")   // failes because of ES6
	***REMOVED***)

	t.Run("script_options_dont_overwrite_metadata", func(t *testing.T) ***REMOVED***
		t.Parallel()
		code := `export let options = ***REMOVED*** vus: 12345 ***REMOVED***; export default function() ***REMOVED*** return options.vus; ***REMOVED***;`
		arc := &lib.Archive***REMOVED***
			Type:        "js",
			FilenameURL: &url.URL***REMOVED***Scheme: "file", Path: "/script"***REMOVED***,
			K6Version:   consts.Version,
			Data:        []byte(code),
			Options:     lib.Options***REMOVED***VUs: null.IntFrom(999)***REMOVED***,
			PwdURL:      &url.URL***REMOVED***Scheme: "file", Path: "/"***REMOVED***,
			Filesystems: nil,
		***REMOVED***
		b, err := NewBundleFromArchive(logger, arc, lib.RuntimeOptions***REMOVED******REMOVED***)
		require.NoError(t, err)
		bi, err := b.Instantiate(logger, 0)
		require.NoError(t, err)
		val, err := bi.exports[consts.DefaultFn](goja.Undefined())
		require.NoError(t, err)
		assert.Equal(t, int64(999), val.Export())
	***REMOVED***)
***REMOVED***

func TestOpen(t *testing.T) ***REMOVED***
	testCases := [...]struct ***REMOVED***
		name           string
		openPath       string
		pwd            string
		isError        bool
		isArchiveError bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:     "notOpeningUrls",
			openPath: "github.com",
			isError:  true,
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "simple",
			openPath: "file.txt",
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "simple with dot",
			openPath: "./file.txt",
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "simple with two dots",
			openPath: "../to/file.txt",
			pwd:      "/path/not",
		***REMOVED***,
		***REMOVED***
			name:     "fullpath",
			openPath: "/path/to/file.txt",
			pwd:      "/path/to",
		***REMOVED***,
		***REMOVED***
			name:     "fullpath2",
			openPath: "/path/to/file.txt",
			pwd:      "/path",
		***REMOVED***,
		***REMOVED***
			name:     "file is dir",
			openPath: "/path/to/",
			pwd:      "/path/to",
			isError:  true,
		***REMOVED***,
		***REMOVED***
			name:     "file is missing",
			openPath: "/path/to/missing.txt",
			isError:  true,
		***REMOVED***,
		***REMOVED***
			name:     "relative1",
			openPath: "to/file.txt",
			pwd:      "/path",
		***REMOVED***,
		***REMOVED***
			name:     "relative2",
			openPath: "./path/to/file.txt",
			pwd:      "/",
		***REMOVED***,
		***REMOVED***
			name:     "relative wonky",
			openPath: "../path/to/file.txt",
			pwd:      "/path",
		***REMOVED***,
		***REMOVED***
			name:     "empty open doesn't panic",
			openPath: "",
			pwd:      "/path",
			isError:  true,
		***REMOVED***,
	***REMOVED***
	fss := map[string]func() (afero.Fs, string, func())***REMOVED***
		"MemMapFS": func() (afero.Fs, string, func()) ***REMOVED***
			fs := afero.NewMemMapFs()
			require.NoError(t, fs.MkdirAll("/path/to", 0o755))
			require.NoError(t, afero.WriteFile(fs, "/path/to/file.txt", []byte(`hi`), 0o644))
			return fs, "", func() ***REMOVED******REMOVED***
		***REMOVED***,
		"OsFS": func() (afero.Fs, string, func()) ***REMOVED***
			prefix, err := ioutil.TempDir("", "k6_open_test")
			require.NoError(t, err)
			fs := afero.NewOsFs()
			filePath := filepath.Join(prefix, "/path/to/file.txt")
			require.NoError(t, fs.MkdirAll(filepath.Join(prefix, "/path/to"), 0o755))
			require.NoError(t, afero.WriteFile(fs, filePath, []byte(`hi`), 0o644))
			if isWindows ***REMOVED***
				fs = fsext.NewTrimFilePathSeparatorFs(fs)
			***REMOVED***
			return fs, prefix, func() ***REMOVED*** require.NoError(t, os.RemoveAll(prefix)) ***REMOVED***
		***REMOVED***,
	***REMOVED***

	logger := testutils.NewLogger(t)

	for name, fsInit := range fss ***REMOVED***
		fs, prefix, cleanUp := fsInit()
		defer cleanUp()
		fs = afero.NewReadOnlyFs(fs)
		t.Run(name, func(t *testing.T) ***REMOVED***
			for _, tCase := range testCases ***REMOVED***
				tCase := tCase

				testFunc := func(t *testing.T) ***REMOVED***
					openPath := tCase.openPath
					// if fullpath prepend prefix
					if openPath != "" && (openPath[0] == '/' || openPath[0] == '\\') ***REMOVED***
						openPath = filepath.Join(prefix, openPath)
					***REMOVED***
					if isWindows ***REMOVED***
						openPath = strings.Replace(openPath, `\`, `\\`, -1)
					***REMOVED***
					pwd := tCase.pwd
					if pwd == "" ***REMOVED***
						pwd = "/path/to/"
					***REMOVED***
					data := `
						export let file = open("` + openPath + `");
						export default function() ***REMOVED*** return file ***REMOVED***;`

					sourceBundle, err := getSimpleBundle(t, filepath.ToSlash(filepath.Join(prefix, pwd, "script.js")), data, fs)
					if tCase.isError ***REMOVED***
						assert.Error(t, err)
						return
					***REMOVED***
					require.NoError(t, err)

					arcBundle, err := NewBundleFromArchive(logger, sourceBundle.makeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)

					require.NoError(t, err)

					for source, b := range map[string]*Bundle***REMOVED***"source": sourceBundle, "archive": arcBundle***REMOVED*** ***REMOVED***
						b := b
						t.Run(source, func(t *testing.T) ***REMOVED***
							bi, err := b.Instantiate(logger, 0)
							require.NoError(t, err)
							v, err := bi.exports[consts.DefaultFn](goja.Undefined())
							require.NoError(t, err)
							assert.Equal(t, "hi", v.Export())
						***REMOVED***)
					***REMOVED***
				***REMOVED***

				t.Run(tCase.name, testFunc)
				if isWindows ***REMOVED***
					// windowsify the testcase
					tCase.openPath = strings.Replace(tCase.openPath, `/`, `\`, -1)
					tCase.pwd = strings.Replace(tCase.pwd, `/`, `\`, -1)
					t.Run(tCase.name+" with windows slash", testFunc)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBundleInstantiate(t *testing.T) ***REMOVED***
	b, err := getSimpleBundle(t, "/script.js", `
		export let options = ***REMOVED***
			vus: 5,
			teardownTimeout: '1s',
		***REMOVED***;
		let val = true;
		export default function() ***REMOVED*** return val; ***REMOVED***
	`)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	logger := testutils.NewLogger(t)

	bi, err := b.Instantiate(logger, 0)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	t.Run("Run", func(t *testing.T) ***REMOVED***
		v, err := bi.exports[consts.DefaultFn](goja.Undefined())
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, true, v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("SetAndRun", func(t *testing.T) ***REMOVED***
		bi.Runtime.Set("val", false)
		v, err := bi.exports[consts.DefaultFn](goja.Undefined())
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, false, v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Options", func(t *testing.T) ***REMOVED***
		// Ensure `options` properties are correctly marshalled
		jsOptions := bi.Runtime.Get("options").ToObject(bi.Runtime)
		vus := jsOptions.Get("vus").Export()
		assert.Equal(t, int64(5), vus)
		tdt := jsOptions.Get("teardownTimeout").Export()
		assert.Equal(t, "1s", tdt)

		// Ensure options propagate correctly from outside to the script
		optOrig := b.Options.VUs
		b.Options.VUs = null.IntFrom(10)
		bi2, err := b.Instantiate(logger, 0)
		assert.NoError(t, err)
		jsOptions = bi2.Runtime.Get("options").ToObject(bi2.Runtime)
		vus = jsOptions.Get("vus").Export()
		assert.Equal(t, int64(10), vus)
		b.Options.VUs = optOrig
	***REMOVED***)
***REMOVED***

func TestBundleEnv(t *testing.T) ***REMOVED***
	rtOpts := lib.RuntimeOptions***REMOVED***Env: map[string]string***REMOVED***
		"TEST_A": "1",
		"TEST_B": "",
	***REMOVED******REMOVED***
	data := `
		export default function() ***REMOVED***
			if (__ENV.TEST_A !== "1") ***REMOVED*** throw new Error("Invalid TEST_A: " + __ENV.TEST_A); ***REMOVED***
			if (__ENV.TEST_B !== "") ***REMOVED*** throw new Error("Invalid TEST_B: " + __ENV.TEST_B); ***REMOVED***
		***REMOVED***
	`
	b1, err := getSimpleBundle(t, "/script.js", data, rtOpts)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	logger := testutils.NewLogger(t)
	b2, err := NewBundleFromArchive(logger, b1.makeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	bundles := map[string]*Bundle***REMOVED***"Source": b1, "Archive": b2***REMOVED***
	for name, b := range bundles ***REMOVED***
		b := b
		t.Run(name, func(t *testing.T) ***REMOVED***
			assert.Equal(t, "1", b.RuntimeOptions.Env["TEST_A"])
			assert.Equal(t, "", b.RuntimeOptions.Env["TEST_B"])

			bi, err := b.Instantiate(logger, 0)
			if assert.NoError(t, err) ***REMOVED***
				_, err := bi.exports[consts.DefaultFn](goja.Undefined())
				assert.NoError(t, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBundleNotSharable(t *testing.T) ***REMOVED***
	data := `
		export default function() ***REMOVED***
			if (__ITER == 0) ***REMOVED***
				if (typeof __ENV.something !== "undefined") ***REMOVED***
					throw new Error("invalid something: " + __ENV.something + " should be undefined");
				***REMOVED***
				__ENV.something = __VU;
			***REMOVED*** else if (__ENV.something != __VU) ***REMOVED***
				throw new Error("invalid something: " + __ENV.something+ " should be "+ __VU);
			***REMOVED***
		***REMOVED***
	`
	b1, err := getSimpleBundle(t, "/script.js", data)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***
	logger := testutils.NewLogger(t)

	b2, err := NewBundleFromArchive(logger, b1.makeArchive(), lib.RuntimeOptions***REMOVED******REMOVED***)
	if !assert.NoError(t, err) ***REMOVED***
		return
	***REMOVED***

	bundles := map[string]*Bundle***REMOVED***"Source": b1, "Archive": b2***REMOVED***
	vus, iters := 10, 1000
	for name, b := range bundles ***REMOVED***
		b := b
		t.Run(name, func(t *testing.T) ***REMOVED***
			for i := 0; i < vus; i++ ***REMOVED***
				bi, err := b.Instantiate(logger, int64(i))
				require.NoError(t, err)
				for j := 0; j < iters; j++ ***REMOVED***
					bi.Runtime.Set("__ITER", j)
					_, err := bi.exports[consts.DefaultFn](goja.Undefined())
					assert.NoError(t, err)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestBundleMakeArchive(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		cm      lib.CompatibilityMode
		script  string
		exclaim string
	***REMOVED******REMOVED***
		***REMOVED***
			lib.CompatibilityModeExtended, `
				import exclaim from "./exclaim.js";
				export let options = ***REMOVED*** vus: 12345 ***REMOVED***;
				export let file = open("./file.txt");
				export default function() ***REMOVED*** return exclaim(file); ***REMOVED***;`,
			`export default function(s) ***REMOVED*** return s + "!" ***REMOVED***;`,
		***REMOVED***,
		***REMOVED***
			lib.CompatibilityModeBase, `
				var exclaim = require("./exclaim.js");
				module.exports.options = ***REMOVED*** vus: 12345 ***REMOVED***;
				module.exports.file = open("./file.txt");
				module.exports.default = function() ***REMOVED*** return exclaim(module.exports.file); ***REMOVED***;`,
			`module.exports.default = function(s) ***REMOVED*** return s + "!" ***REMOVED***;`,
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.cm.String(), func(t *testing.T) ***REMOVED***
			fs := afero.NewMemMapFs()
			_ = fs.MkdirAll("/path/to", 0o755)
			_ = afero.WriteFile(fs, "/path/to/file.txt", []byte(`hi`), 0o644)
			_ = afero.WriteFile(fs, "/path/to/exclaim.js", []byte(tc.exclaim), 0o644)

			rtOpts := lib.RuntimeOptions***REMOVED***CompatibilityMode: null.StringFrom(tc.cm.String())***REMOVED***
			b, err := getSimpleBundle(t, "/path/to/script.js", tc.script, fs, rtOpts)
			assert.NoError(t, err)

			arc := b.makeArchive()

			assert.Equal(t, "js", arc.Type)
			assert.Equal(t, lib.Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***, arc.Options)
			assert.Equal(t, "file:///path/to/script.js", arc.FilenameURL.String())
			assert.Equal(t, tc.script, string(arc.Data))
			assert.Equal(t, "file:///path/to/", arc.PwdURL.String())

			exclaimData, err := afero.ReadFile(arc.Filesystems["file"], "/path/to/exclaim.js")
			assert.NoError(t, err)
			assert.Equal(t, tc.exclaim, string(exclaimData))

			fileData, err := afero.ReadFile(arc.Filesystems["file"], "/path/to/file.txt")
			assert.NoError(t, err)
			assert.Equal(t, `hi`, string(fileData))
			assert.Equal(t, consts.Version, arc.K6Version)
			assert.Equal(t, tc.cm.String(), arc.CompatibilityMode)
		***REMOVED***)
	***REMOVED***
***REMOVED***
